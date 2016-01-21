package simlog

import (
	"fmt"
	"time"
	"os"
	"testing"
)

func TestMakeDirAll(t *testing.T) {
	dir := "TestMakeDirAll"
	if err := MakeDirAll(dir); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	st, e := file.Stat()
	if e != nil {
		t.Fatal(e)
	}
	if !st.IsDir() {
		t.Fatal("MakeDirAll not dir")
	}
	os.RemoveAll(dir)
}

func TestOpenAndCreateFile1(t *testing.T) {
	f := "TestOpenAndCreateFile1.log"
	file, err := OpenAndCreateFile(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	
	os.Remove(f)
}

func TestOpenAndCreateFile2(t *testing.T) {
	f := "TestOpenAndCreateFile2.log"
	file, err := OpenAndCreateFile(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	file, err = OpenAndCreateFile(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	os.Remove(f)
}

func TestReadDir(t *testing.T) {
	dir := "TestReadDir"
	if err := MakeDirAll(dir); err != nil {
		t.Fatal(err)
	}
	files, err := ReadDir(dir)
	if err != nil {
		t.Error(err)
	}
	if len(files) != 0 {
		t.Error("ReadDir failed")
	}
	os.RemoveAll(dir)
}

func TestReadDir1(t *testing.T) {
	dir := "TestReadDir1"
	if err := MakeDirAll(dir); err != nil {
		t.Fatal(err)
	}
	f := dir + "/test.log"
	file, err := OpenAndCreateFile(f, 0)
	if err != nil {
		t.Error(err)
	}
	file.Close()
	files, e := ReadDir(dir)
	if e != nil {
		t.Error(e)
	}
	if len(files) != 1 {
		t.Error("ReadDir failed|", len(files))
	}
	
	os.RemoveAll(dir)
}

func TestGetFileSize(t *testing.T) {
	f := "TestGetFileSize.log"
	file, err := OpenAndCreateFile(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	
	size, e := GetFileSize(file)
	if e != nil || size != 0 {
		t.Error("GetFileSize failed|size=", size, "|err=",e)
	}
	n, e1 := file.WriteString("12345")
	if e1 != nil {
		t.Error("file.WriteString failed|err=", e)
	}
	size, e = GetFileSize(file)
	if e != nil || size != int64(n) {
		t.Error("GetFileSize failed|size=", size, "|n=", n, "|err=", e)
	}
	file.Close()
	os.Remove(f)
}

func TestLogPathName(t *testing.T) {
	now := time.Now()
	format := "%04d%02d%02d"
	rightName := fmt.Sprintf(format, now.Year(), now.Month(), now.Day())
	pathName := LogPathName(now)
	
	if pathName != rightName {
		t.Error("pathName=", pathName, " |rightName=", rightName)
	}
}

func TestLogPathName1(t *testing.T) {
	now := time.Now()
	format := "test%04d%02d%02d_02"
	rightName := fmt.Sprintf(format, now.Year(), now.Month(), now.Day())
	pathName := LogPathName(now, 2, "test")
	if pathName != rightName {
		t.Error("pathName=", pathName, " |rightName=", rightName)
	}
}

