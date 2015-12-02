package signal

import (
	"syscall"
)

//+build drwin

func Kill(pid int, sid syscall.Signal) error {
	return nil
}
