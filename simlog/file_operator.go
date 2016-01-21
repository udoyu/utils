package simlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func MakeDirAll(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

//flag设置追加os.O_Append或清零os.O_TRUNC
func OpenAndCreateFile(filename string, flag int) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|flag, os.ModePerm)
}


func ReadDir(dir string) (fi []os.FileInfo, err error) {
	file, err := os.Open(dir)
	if nil != err {
		return fi, err
	}
	defer file.Close()
	fi, err = file.Readdir(0)
	return fi, err
}

func GetFileSize(file *os.File) (int64, error) {
	fi, err := file.Stat()
	if nil != err {
		return -1, err
	}
	return fi.Size(), nil
}

func ClearFile(filename string) error {
	return ioutil.WriteFile(filename, []byte(""), os.ModePerm)
}

//时间， 索引号， 前缀
func LogPathName(curtime time.Time, v ...interface{}) string {
	format := "%s%04d%02d%02d"
	prefix := ""
	index := -1
	ok := false
	if len(v) > 0 {
		if index,ok = v[0].(int);ok {
			format += "_%02d"
		}
	}
	if len(v) > 1 {
		if prefix,ok = v[1].(string);ok {
		}
	}
	if -1 != index {
		return fmt.Sprintf(format, prefix, curtime.Year(), curtime.Month(), curtime.Day(), index)
	} else {
		return fmt.Sprintf(format, prefix, curtime.Year(), curtime.Month(), curtime.Day())
	}
}
