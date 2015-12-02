package simlog

// from beego

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	DayDur   = time.Hour * 24   //一天的duration
	TimerDur = time.Minute * 10 //用于定时判断日期的时间间隔
)

var (
	initFlag     bool        //是否已经初始化
	splitFlag    bool        //是否根据文件大小切割
	logbasepath  string      //日志根目录，传参设定
	MAXFILESIZE  int64       //每个文件最大大小，传参设定
	MAXLOGINDEX  int         //每天最大文件个数，传参设定
	MAXLOGDAY    int         //保存文件最大天数，传参设定
	MAXLOGCNT    int         //日志判断周期，即生成MAXLOGCNT条日志后开始判断当前文件是否超过设定大小
	logdate      time.Time   //当天的日期，用来判断日期是否改变
	logfilepath  string      //日志文件路径，由logbasepath+当前日期得到
	logfileindex int         //日志文件当前序号
	logfilename  string      //日志文件名，由logfilepath+logfileindex.log得到
	logfile      *os.File    //日志文件
	logfilelock  *sync.Mutex //日志锁，在改变全局变量时需要用到
	logcnt       int         //当前日志条数，用来和日志判断周期比较，避免频繁判断文件大小
	// logLevel controls the global log level used by the logger.
)

func getdate(t time.Time) int {
	tdate, _ := strconv.Atoi(fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day()))
	return tdate
}

func GetLogName(path string) string {
	fis, err := ReadDir(path)
	if nil != err {
		fmt.Printf("[nettao]GetLogName|ReadDir %s|%s failed", path, err.Error())
		os.Exit(-1)
	}
	filename := ""
	if splitFlag {
		logfileindex = 0

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			logfileindex++
		}

		if logfileindex >= MAXLOGINDEX {
			logfileindex = 0
		}
		logfileindex++
		filename = path + fmt.Sprintf("%d.log", logfileindex)
	} else {
		filename = path + "all.log"
	}
	return filename
}

func logInit(path string, maxday int, loglevel Level) {
	if path == logbasepath && maxday == MAXLOGDAY && loglevel == loglevel {
		return
	}
	logfilelock = new(sync.Mutex)
	logfilelock.Lock()
	defer logfilelock.Unlock()
	now := time.Now()

	MAXLOGDAY = maxday
	SetLogLevel(loglevel)
	logdate = now
	logbasepath = path + "/"
	logfilepath = logbasepath + MakeLogPath(now) + "/"
	err := MakeDirAll(logfilepath)
	if nil != err {
		fmt.Printf("[nettao]LogInit|MakeDirAll logpath %s|%s failed\n", logfilepath, err.Error())
		os.Exit(-1)
	}

	logfilename = GetLogName(logfilepath)

	logfile, err = OpenAndCreateFile(logfilename, os.O_APPEND)
	if nil != err {
		fmt.Printf("[nettao]Start|open log file %s|%s\n", logfilename, err.Error())
		os.Exit(-1)
	}
	SimLogger = log.New(logfile, "\n", log.Ldate|log.Ltime|log.Llongfile)
	removelogdir(MAXLOGDAY, now)
	if !initFlag {
		initFlag = true
		go changelogdate()
	}
}

func setLogSplit(maxsize, maxindex int) {
	splitFlag = true
	MAXFILESIZE = int64(maxsize * 1024 * 1024)
	MAXLOGCNT = int(MAXFILESIZE / 1024)
	MAXLOGINDEX = maxindex
}

func logClose() {
	if nil != logfile {
		logfile.Close()
	}
}

func changelogindex(addsize int) {
	if filesize, _ := GetFileSize(logfile); filesize+int64(500+addsize) >= MAXFILESIZE {
		if logfileindex >= MAXLOGINDEX {
			logfileindex = 0
		}
		logfileindex++
		changelogfile()
	}

}

func removelogdir(daynum int, now time.Time) {
	tmpdate := now.Add(-DayDur * time.Duration(daynum))
	tmpdir := logbasepath + MakeLogPath(tmpdate) + "/"
	os.RemoveAll(tmpdir)
}

func changelogdate() {
	logTimer := time.NewTicker(TimerDur)
	for {
		select {
		case <-logTimer.C:
			func() {
				logfilelock.Lock()
				defer logfilelock.Unlock()
				now := time.Now()
				if now.Day() == logdate.Day() && now.Month() == logdate.Month() && now.Year() == logdate.Year() {
					return
				}
				//如果日期改变
				//将全局信息重置
				logdate = now
				logfilepath = logbasepath + MakeLogPath(now) + "/"
				logfileindex = 0
				logfileindex++
				err := MakeDirAll(logfilepath)
				if nil != err {
					fmt.Printf("[nettao]ChangeLogPathOrFile|MakeDirAll %s|%s \n", logfilepath, err.Error())
					os.Exit(-1)
				}
				//删除配置指定时间的日志文件
				removelogdir(MAXLOGDAY, now)
				changelogfile()
			}()
		}
	}

}

func changelogfile() {
	//关闭上一个文件
	logfile.Close()
	//创建新文件
	if splitFlag {
		logfilename = logfilepath + fmt.Sprintf("%d.log", logfileindex)
	} else {
		logfilename = logfilepath + "all.log"
	}
	var err error
	logfile, err = OpenAndCreateFile(logfilename, os.O_TRUNC)
	if nil != err {
		fmt.Printf("[nettao]ChangeLogPathOrFile|open log file %s|%s\n", logfilename, err.Error())
		os.Exit(-1)
	}
	SimLogger = log.New(logfile, "\n", log.Ldate|log.Ltime|log.Llongfile)
}
