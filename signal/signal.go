package signal

import (
	"os/signal"
	"os"
	"syscall"
)
var (
	SIGUSR1 = syscall.Signal(10)
	SIGUSR2 = syscall.Signal(12)
)

type SignalHandler struct {
	signalMap map[os.Signal]func()
	signals chan os.Signal
}

func NewSignalHandler() SignalHandler {
	return SignalHandler{
		signalMap : make(map[os.Signal]func()),
		signals : make(chan os.Signal),
	}
}

func (this SignalHandler) Listen() {
	for k, _ := range this.signalMap {
		signal.Notify(this.signals, k)
	}
	for sig := range this.signals {
		callback, ok := this.signalMap[sig]
		if ok {
			callback()
		}
	}
}

func (this SignalHandler) Register(sig int, callback func()) {
	this.signalMap[syscall.Signal(sig)] = callback
}

var (
	signalHandler SignalHandler = NewSignalHandler()
)

func Register(sig int, callback func()) {
	signalHandler.Register(sig, callback)
}

func Listen() {
	signalHandler.Listen()
}

