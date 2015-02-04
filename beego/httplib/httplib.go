package httplib

import (
    "net/http" 
)

func Transmit(trans_url string, r *http.Request) (*BeegoHttpRequest, error) {
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
    if r.Form == nil {
        r.ParseForm()
    }
    for k, v := range r.Form {
        b.Param(k, v[0])
    }
    b.Request().Header = r.Header
    err := b.Do()
    if err == nil {
        r.Header = b.Request().Header
    }
    return b, err
}

