package httplib

import (
    "net/http" 
    "github.com/astaxie/beego/httplib" 
)

func Transmit(method, trans_url string, r *http.Request) (*httplib.BeegoHttpRequest, error) {
    var b *httplib.BeegoHttpRequest
    urlString := trans_url
    if len(r.URL.RawQuery) != 0 {
        urlString += "?" + r.URL.RawQuery
    }
    
    switch method {
        case "GET" : b = httplib.Get(urlString)
        case "POST" : b = httplib.Post(urlString)
        case "PUT" : b = httplib.Put(urlString)
        case "DELETE" : b = httplib.Delete(urlString)
        case "HEAD" : b = httplib.Head(urlString)
    }
    if r.Form == nil {
        r.ParseForm()
    }
    for k, v := range r.Form {
        b.Param(k, v[0])
    }
    for k,v := range r.Header{
            for _, s := range v {
                b.Header(k, s)
            }
    }
    resp, err := b.Response()
    if err == nil {
        r.Header = resp.Header
    }
    return b, err
}

