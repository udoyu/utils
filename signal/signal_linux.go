package signal

import (
	"syscall"
)

//+build linux
func init() {
	SIGUSR1 = syscall.SIGUSR1
	SIGUSR2 = syscall.SIGUSR2
}
func Kill(pid int, sig syscall.Signal) error {
	return syscall.Kill(pid, sig)
}
