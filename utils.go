package utils

import(
	"reflect"
)

//s := New(reflect.ValueOf("a")).(string)
//s = "hello"
//fmt.Println(s)
func New(v reflect.Value) interface{} {
        if v.Kind() == reflect.Ptr {
                return New(v.Elem())
        }
        r := reflect.New(v.Type())
        return r.Elem().Interface()
}