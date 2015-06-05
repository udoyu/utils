package mybeego

import (
	"github.com/astaxie/beego"
	"github.com/udoyu/utils/beego/httplib"
	"net/http"
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

func (this *Controller) Transmit(url string) (*httplib.BeegoHttpRequest, error) {
	r := this.Ctx.Request
	var b *httplib.BeegoHttpRequest
	urlString := url
	if len(r.URL.RawQuery) != 0 {
		urlString += "?" + r.URL.RawQuery
	}

	switch r.Method {
	case "GET":
		b = httplib.Get(urlString)
	case "POST":
		b = httplib.Post(urlString)
	case "PUT":
		b = httplib.Put(urlString)
	case "DELETE":
		b = httplib.Delete(urlString)
	case "HEAD":
		b = httplib.Head(urlString)
	}
	if r.Form != nil && len(r.Form) != 0 {
		for k, v := range r.Form {
			b.Param(k, v[0])
		}
	} else {
		b.Request().Body = r.Body
	}
	header := this.GetSession("trans_session")
	if header != nil {
		cookies := header.([]*http.Cookie)
		req := b.Request()
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
	err := b.Do()
	cookies := b.Response().Cookies()
	if len(cookies) > 0 {
		this.SetSession("trans_session", cookies)
	}
	return b, err

}
