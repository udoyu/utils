package controller
import (
    "github.com/astaxie/beego"
)

type ControllerFunc func(ctl beego.ControllerInterface)
type ControllerHandler map[string]ControllerFunc

type ControllerClientInterface interface{
    func AddHandler(path string, h ControllerFunc)
    func DelHandler(path string)
    func Handler(beego.ControllerClientInterface)
}

type ControllerClient struct {
    HandlerFunc ControllerHandler
}

type (this *ControllerClient)AddHandler(path string, h ControllerFunc) {
    this.HandlerFunc[path] = h
}
type (this *ControllerClient)Handler(ctl beego.ControllerInterface) {
    if h,ok := this.HandlerFunc[ctl.Ctx.Input.Url()]
    if ok {
        h(ctl)
    }
}
func (this *ControllerClient)DelHandler(path string) {
    delete(this.HandlerFunc, path)
}

