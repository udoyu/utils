package httplib

import (
    "net/http" 
)

func Transmit(trans_url string, r *http.Request, w http.ResponseWriter) (*BeegoHttpRequest, error) {
    var b *BeegoHttpRequest
    urlString := trans_url
    if len(r.URL.RawQuery) != 0 {
        urlString += "?" + r.URL.RawQuery
    }
    
    switch r.Method {
        case "GET" : b = Get(urlString)
        case "POST" : b = Post(urlString)
        case "PUT" : b = Put(urlString)
        case "DELETE" : b = Delete(urlString)
        case "HEAD" : b = Head(urlString)
    }
    if r.Form != nil && len(r.Form) != 0{
        for k, v := range r.Form {
            b.Param(k, v[0])
        }
    } else {
            b.Request().Body = r.Body
    }
    b.Request().Header = r.Header
    
    err := b.Do()
    if err == nil {
        header := b.Response().Header
        for k,v := range header {
            for _,s := range v {
                w.Header().Add(k, s)
            }
        }
    }
    return b, err
}

