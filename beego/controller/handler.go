package controller

import (
	"github.com/astaxie/beego"
)

type ControllerFunc func(ctl beego.ControllerInterface)
type ControllerHandlerMap map[string]ControllerFunc

type ControllerHandlerInterface interface {
	Add(path string, h ControllerFunc)
	Delete(path string)
	Default(v beego.ControllerInterface)
	Do(path string, v beego.ControllerInterface)
}

type ControllerHandler struct {
	HandlerFunc ControllerHandlerMap
}

func (this *ControllerHandler) Add(path string, h ControllerFunc) {
	this.HandlerFunc[path] = h
}
func (this *ControllerHandler) Delete(path string) {
	delete(this.HandlerFunc, path)
}
func (this *ControllerHandler) Do(path string, ctl beego.ControllerInterface) {
	if h, ok := this.HandlerFunc[path]; ok {
		h(ctl)
	} else {
		this.Default(ctl)
	}
}
func (this *ControllerHandler) Default(ctl beego.ControllerInterface) {
	ctl.(*beego.Controller).Ctx.WriteString("")
}
