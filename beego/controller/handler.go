package controller
import (
    "github.com/astaxie/beego"
)

type ControllerFunc func(ctl beego.ControllerInterface)
type ControllerHandlerMap map[string]ControllerFunc

type ControllerHandlerInterface interface{
    func Add(path string, h ControllerFunc)
    func Delete(path string)
    func Default(v beego.ControllerInterface)
    func Do(path string, v beego.ControllerHandlerInterface)
}

type ControllerHandler struct {
    HandlerFunc ControllerHandlerMap
}
type (this *ControllerHandler) Add(path string, h ControllerFunc) {
    this.HandlerFunc[path] = h
}
func (this *ControllerHandler) Delete(path string) {
    delete(this.HandlerFunc, path)
}
func (this *ControllerHandler) Do(path string, ctl beego.ControllerInterface) {
    if h,ok := this.HandlerFunc[path]
    if ok {
        h(ctl)
    } else {
        this.Default(ctl)
    }
}
func (this *ControllerHandler) Default(ctl beego.ControllerInterface) {
    ctl.Ctx.WriteString("")
}
