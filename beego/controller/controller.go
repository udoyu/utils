package controller

import (
//        "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
        "github.com/xiying/xytool/simini"
	//"github.com/astaxie/beego/context"
)

type Controller struct {
	beego.Controller
        handler CtrHandler
}

type CtrHandler interface {
	Handler(*Controller) bool
        AddHandler(path string, h func(*Controller)int)
        DelHandler(path string)
}

func NewController(h CtrHandler) *Controller{
    c := &Controller{}
    c.handler = h
    return c
}

var globalSessions *session.Manager

func InitSession(ini *simini.SimIni) error {
        storeType := ini.GetStringVal("session", "type")
        storeConf := ini.GetStringVal("session", "conf")
        s,e := session.NewManager(storeType, storeConf)
        if e != nil {
            beego.Error(e.Error())
        }
        globalSessions = s
	go globalSessions.GC()
        return nil
}

func GlobalSession() *session.Manager {
    return globalSessions
}

func SetStaticPath(k, v string) {
    beego.SetStaticPath(k,v)
}

func Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
    return beego.Router(rootpath, c, mappingMethods...)
}

func Run(params ...string) {
    beego.Run(params...)
}


//func AddCtrHandler(s string, h CtrHandler) {
//	globalClientMap[s] = h
//}
//
//func DelCtrHandler(s string) {
//	delete(globalClientMap, s)
//}

func (this *Controller) Get() {
        this.handler.Handler(this)
}

func (this *Controller) Path() string {
    return this.Ctx.Request.URL.Path
}

func (this *Controller) FormValue(key string) string {
	return this.Ctx.Request.FormValue(key)
}

func (this *Controller) WriteString(s string) {
	this.Ctx.WriteString(s)
}

func (this *Controller) Write(s []byte) {
        this.Ctx.ResponseWriter.Write(s)
}

func (this *Controller) SessionStart() (session session.SessionStore, err error) {
	return globalSessions.SessionStart(this.Ctx.ResponseWriter,
		this.Ctx.Request)
}

func (this *Controller) SessionRelease(session session.SessionStore) {
	session.SessionRelease(this.Ctx.ResponseWriter)
}

