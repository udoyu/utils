package utils

import(
	"reflect"
)

//s := "a"
//ps := New(reflect.ValueOf(&a)).(*string)
//*ps = "hello"
//fmt.Println(*ps)
func New(v reflect.Value) interface{} {
        r := reflect.New(v.Type())
        if v.Kind() == reflect.Ptr {
                r.Elem().Set(reflect.New(v.Elem().Type()))
        }
        return r.Elem().Interface()
}
