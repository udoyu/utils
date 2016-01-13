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
	LOG_BASE_NAME = "all.log"
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
	//	fis, err := ReadDir(path)
	//	if nil != err {
	//		fmt.Printf("GetLogName|ReadDir %s|%s failed", path, err.Error())
	//		os.Exit(-1)
	//	}
	//	filename := ""
	//	if splitFlag {
	//		logfileindex = 0

	//		for _, fi := range fis {
	//			if fi.IsDir() {
	//				continue
	//			}
	//			logfileindex++
	//		}

	//		if logfileindex >= MAXLOGINDEX {
	//			logfileindex = 0
	//		}
	//		logfileindex++
	//		filename = path + fmt.Sprintf("%d.log", logfileindex)
	//	} else {
	//		filename = path + LOG_BASE_NAME
	//	}
	return path + LOG_BASE_NAME
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
	logfilepath = logbasepath 
	err := MakeDirAll(logfilepath)
	if nil != err {
		fmt.Printf("[simlog]LogInit|MakeDirAll logpath %s|%s failed\n", logfilepath, err.Error())
		os.Exit(-1)
	}

	logfilename = GetLogName(logfilepath)

	logfile, err = OpenAndCreateFile(logfilename, os.O_APPEND)
	if nil != err {
		fmt.Printf("log Start|open log file %s|%s\n", logfilename, err.Error())
		os.Exit(-1)
	}
	SimLogger = log.New(logfile, "\n", log.Ldate|log.Ltime|log.Llongfile)
	initlogfile()
	movelogdir()
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

func movelogdir() {
	fis, err := ReadDir(logbasepath)
	if nil != err {
		fmt.Printf("GetLogName|ReadDir %s|%s failed", logbasepath, err.Error())
		os.Exit(-1)
	}
	var name string
	var path string
	basepos := len(LOG_BASE_NAME+".")
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if fi.Name() == LOG_BASE_NAME {
			continue
		}
		name = fi.Name()
		//all.log.20160113
		if len(name) < basepos {
			continue
		}
		path = logbasepath + name[basepos:basepos+8] + "/"

		if err := MakeDirAll(path); err != nil {
			Error("movelogdir failed ", err.Error())
			continue
		}
		oldname := logbasepath + name
		newname := path + name
		if err := os.Rename(oldname, newname); err != nil {
			fmt.Printf("Rename %s -> %s failed err=%s", oldname, newname, err.Error())
			continue
		}
	}

}

func initlogfile() {
	movelogdir()
	//move all.log
	fis, err := ReadDir(logbasepath)
	if nil != err {
		fmt.Printf("GetLogName|ReadDir %s|%s failed", logbasepath, err.Error())
		os.Exit(-1)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if fi.Name() != LOG_BASE_NAME {
			continue
		}
		nowdate := time.Now()
		olddate := fi.ModTime()
		if nowdate.Day() == olddate.Day() && nowdate.Month() == olddate.Month() && nowdate.Year() == olddate.Year() {
			return
		}
		
		changelogfile()
	}
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
				//logfilepath = logbasepath + MakeLogPath(now) + "/"
				logfileindex = 0
				logfileindex++
				//logfilename = logfilepath + LOG_BASE_NAME
				
				changelogfile()
				//删除配置指定时间的日志文件
				movelogdir()
				removelogdir(MAXLOGDAY, now)
			}()
		}
	}

}

func changelogfile() {

	//创建新文件
	//logfilename = logfilepath + LOG_BASE_NAME
	//	if splitFlag {
	for i := 0; i < 1; i++ {
		//关闭上一个文件
		logfile.Close()
		now := time.Now()

		old := logfilepath + fmt.Sprintf("all.log.%04d%02d%02d-%02d%02d%02d.%d",
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), logfileindex)
		if err := os.Rename(logfilename, old); err != nil {
			fmt.Printf("Rename %s -> %s failed err=%s", logfilename, old, err.Error())
			break
		}
	}
	var err error
	logfile, err = OpenAndCreateFile(logfilename, os.O_TRUNC)
	if nil != err {
		fmt.Printf("ChangeLogPathOrFile|open log file %s|%s\n", logfilename, err.Error())
		return
	}
	SimLogger = log.New(logfile, "\n", log.Ldate|log.Ltime|log.Llongfile)
}
