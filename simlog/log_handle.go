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

func getdate(t time.Time) int {
	tdate, _ := strconv.Atoi(fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day()))
	return tdate
}

type LogHandler struct {
	BaseName    string
	MaxDataSize int
	MaxFileSize int64 //每个文件最大大小，传参设定
	MaxLogIndex int   //每天最大文件个数，传参设定
	MaxDay      int   //保存文件最大天数，传参设定
	MaxSize     int   //日志判断周期，即生成MAXLOGCNT条日志后开始判断当前文件是否超过设定大小
	logger      *log.Logger
	level       Level      //日志级别
	splitFlag   bool       //是否根据文件大小切割
	logPath     string     //日志根目录，传参设定
	logDate     time.Time  //当天的日期，用来判断日期是否改变
	filePath    string     //日志文件路径，由logbasepath+当前日期得到
	fileIndex   int        //日志文件当前序号
	fileName    string     //日志文件名，由logfilepath+logfileindex.log得到
	file        *os.File   //日志文件
	lock        sync.Mutex //日志锁，在改变全局变量时需要用到
	size        int        //当前日志条数，用来和日志判断周期比较，避免频繁判断文件大小
}

func (this *LogHandler) OutPut(level Level, v ...interface{}) {
	if this.level < level {
		str := fmt.Sprint(v...)
		size := len(str)
		if size > this.MaxDataSize {
			size = this.MaxDataSize
		}
		if this.logger != nil {

			this.lock.Lock()
			this.logger.Output(4, level.String()+str[:size])
			this.size += size
			this.logSplit()
			this.lock.Unlock()
		} else {
			log.Println(level.String() + str[:size])
		}
	}
}

func (this *LogHandler) Trace(v ...interface{}) {
	this.OutPut(LevelTrace, v...)
}

func (this *LogHandler) Debug(v ...interface{}) {
	this.OutPut(LevelDebug, v...)
}

func (this *LogHandler) Info(v ...interface{}) {
	this.OutPut(LevelInfo, v...)
}

func (this *LogHandler) Warn(v ...interface{}) {
	this.OutPut(LevelWarning, v...)
}

func (this *LogHandler) Error(v ...interface{}) {
	this.OutPut(LevelError, v...)
}

func (this *LogHandler) Critical(v ...interface{}) {
	this.OutPut(LevelCritical, v...)
}

func (this *LogHandler) SetLevel(l Level) {
	this.level = l
}

func (this *LogHandler) GetLevel() Level {
	return this.level
}

func (this *LogHandler) SetLogSplit(maxsize, maxindex int) {
	this.splitFlag = true
	this.MaxFileSize = int64(maxsize * 1024 * 1024)
	this.MaxSize = int(this.MaxFileSize / 1024)
	this.MaxLogIndex = maxindex
}

func (this *LogHandler) logSplit() {
	if this.splitFlag {
		if this.size > this.MaxSize {
			this.changelogindex()
			this.size = 0
		}
	}
}

func (this *LogHandler) changelogindex() {
	if filesize, _ := GetFileSize(this.file); filesize >= this.MaxFileSize {
		if this.fileIndex >= this.MaxLogIndex {
			this.fileIndex = 0
		}
		this.fileIndex++
		this.changelogfile(time.Now())
	}
}

func (this *LogHandler) Init(path string, maxday int, loglevel Level) {
	if path == this.logPath && maxday == this.MaxDay && loglevel == this.level {
		return
	}
	this.MaxDataSize = 4096
	now := time.Now()
	this.lock.Lock()
	defer this.lock.Unlock()
	this.MaxDay = maxday
	this.level = loglevel
	this.logDate = now
	this.logPath = path + "/"
	this.filePath = this.logPath
	err := MakeDirAll(this.filePath)
	if nil != err {
		fmt.Printf("[simlog]LogInit|MakeDirAll logpath %s|%s failed\n", this.filePath, err.Error())
		os.Exit(-1)
	}
	if this.BaseName == "" {
		this.BaseName = "all.log"
	}
	this.fileName = this.filePath + this.BaseName

	this.file, err = OpenAndCreateFile(this.fileName, os.O_APPEND)
	if nil != err {
		fmt.Printf("log Start|open log file %s|%s\n", this.fileName, err.Error())
		os.Exit(-1)
	}
	this.logger = log.New(this.file, "\n", log.Ldate|log.Ltime|log.Llongfile)

	if this.logger == nil {
		this.initlogfile()
		this.movelogdir()
		this.removelogdir(this.MaxDay, now)
		go this.changelogdate()
	}
}

func (this *LogHandler) Close() {
	if this.file != nil {
		this.lock.Lock()
		defer this.lock.Unlock()
		this.logger = nil
		this.file.Close()
	}
}

func (this *LogHandler) movelogdir() {
	fis, err := ReadDir(this.logPath)
	if nil != err {
		fmt.Printf("GetLogName|ReadDir %s|%s failed", this.logPath, err.Error())
		os.Exit(-1)
	}
	var name string
	var path string
	now := time.Now()
	nowstr := fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day())
	basepos := len(this.BaseName + ".")
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if fi.Name() == this.BaseName {
			continue
		}
		name = fi.Name()
		//all.log.20160113
		if len(name) < basepos || name[:len(this.BaseName)] != this.BaseName {
			continue
		}
		if len(name) < basepos+8 || name[basepos:basepos+8] == nowstr {
			continue
		}
		path = this.logPath + name[basepos:basepos+8] + "/"

		if err := MakeDirAll(path); err != nil {
			this.Error("movelogdir failed ", err.Error())
			continue
		}
		oldname := this.logPath + name
		newname := path + name
		if err := os.Rename(oldname, newname); err != nil {
			fmt.Printf("Rename %s -> %s failed err=%s", oldname, newname, err.Error())
			continue
		}
	}
}

func (this *LogHandler) changelogfile(date time.Time) {
	//创建新文件
	//logfilename = logfilepath + LOG_BASE_NAME
	//	if splitFlag {
	for i := 0; i < 1; i++ {
		//关闭上一个文件
		this.file.Close()
		old := this.filePath + fmt.Sprintf("all.log.%04d%02d%02d-%02d%02d%02d.%d",
			date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), this.fileIndex)
		if err := os.Rename(this.fileName, old); err != nil {
			fmt.Printf("Rename %s -> %s failed err=%s", this.fileName, old, err.Error())
			break
		}
	}
	var err error
	this.file, err = OpenAndCreateFile(this.fileName, os.O_TRUNC)
	if nil != err {
		fmt.Printf("ChangeLogPathOrFile|open log file %s|%s\n", this.fileName, err.Error())
		return
	}
	this.size = 0
	this.logger = log.New(this.file, "\n", log.Ldate|log.Ltime|log.Llongfile)
}

func (this *LogHandler) initlogfile() {
	this.movelogdir()
	fis, err := ReadDir(this.logPath)
	if nil != err {
		fmt.Printf("GetLogName|ReadDir %s|%s failed", this.logPath, err.Error())
		os.Exit(-1)
	}
	nowdate := time.Now()
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if fi.Name() != this.BaseName {
			continue
		}

		olddate := fi.ModTime()
		if nowdate.Day() == olddate.Day() &&
			nowdate.Month() == olddate.Month() &&
			nowdate.Year() == olddate.Year() {
			return
		}

		this.changelogfile(olddate)
		break
	}
}

func (this *LogHandler) changelogdate() {
	logTimer := time.NewTicker(TimerDur)
	for {
		select {
		case <-logTimer.C:
			func() {
				this.lock.Lock()
				defer this.lock.Unlock()
				now := time.Now()
				if now.Day() == this.logDate.Day() &&
					now.Month() == this.logDate.Month() &&
					now.Year() == this.logDate.Year() {
					return
				}
				//如果日期改变
				//将全局信息重置
				this.logDate = now
				//logfilepath = logbasepath + MakeLogPath(now) + "/"
				this.fileIndex = 0
				this.fileIndex++
				//logfilename = logfilepath + LOG_BASE_NAME
				st, e := this.file.Stat()
				if e != nil {
					this.Error("logfile.Stat failed|err=", e.Error())
					this.changelogfile(now)
				} else {
					this.changelogfile(st.ModTime())
				}
				//删除配置指定时间的日志文件
				this.movelogdir()
				this.removelogdir(this.MaxDay, now)
			}()
		}
	}
}

func (this *LogHandler) removelogdir(daynum int, now time.Time) {
	tmpdate := now.Add(-DayDur * time.Duration(daynum))
	tmpdir := this.logPath + LogPathName(tmpdate) + "/"
	os.RemoveAll(tmpdir)
}
