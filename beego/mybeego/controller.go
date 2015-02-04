package mybeego

import (
	"github.com/astaxie/beego"
)

type Controller struct {
     beego.Controller
}

func (this *Controller) Query(key string) string {
    if this.Ctx.Input.Request.MultipartForm == nil {
        this.Ctx.Input.ParseFormOrMulitForm(beego.MaxMemory)
    }
    return this.Ctx.Input.Query(key)
}

func (this *Controller) FormValue(key string) string {
    return this.Ctx.Input.Query(key)
}
