package simlog

import (
	"os"
	"testing"
)

func BenchmarkPrintln(t *testing.B) {
	path := "test"
	Init(path, 1, 1)
	SetSplit(100, 0)
	buf := make([]byte, 1024)
	str := string(buf)
	for i := 0; i < t.N; i++ {
		Debug(str)
	}
	Close()
	os.RemoveAll(path)
}
