package utils

import (
    "reflect"
    "strconv"
)

func Atol(s string) int64 {
        i, _ := strconv.ParseInt(s, 10, 64)
        return i
}

func Atoi(s string) int {
        i, _ := strconv.Atoi(s)
        return i
}

type QueryInterface interface {
    Query(t reflect.StructField) string
}

var (
    STRING_PTR = reflect.TypeOf(proto.String(""))
    INT32_PTR = reflect.TypeOf(proto.Int32(int32(0)))
    INT64_PTR = reflect.TypeOf(proto.Int64(int64(0)))
    UINT32_PTR = reflect.TypeOf(proto.Uint32(uint32(0)))
    UINT64_PTR = reflect.TypeOf(proto.Uint64(uint64(0)))
    BOOL_PTR = reflect.TypeOf(proto.Bool(false))
)

func StructTravelFunc (v reflect.Value, t reflect.StructField, qi QueryInterface) {
       if str:=qi.Query(t); len(str)!=0 {
           //form := utils.StringCutRightExp(t.Tag.Get("json"), ",", 1)
           switch v.Type() {
               case STRING_PTR : v.Elem().SetString(str)
               case INT32_PTR : if v.Elem().CanSet() {v.Elem().SetInt(int64(Atoi(str)))}
               case INT64_PTR : if v.Elem().CanSet() {v.Elem().SetInt(Atol(str))}
               case UINT32_PTR : if v.Elem().CanSet() {v.Elem().SetUint(uint64(Atoi(str)))}
               case UINT64_PTR : if v.Elem().CanSet() {v.Elem().SetUint(uint64(Atol(str)))}
               case BOOL_PTR : if v.Elem().CanSet() {
                                   if str == "true" {
                                       v.Elem().SetBool(true)
                                   } else {
                                       v.Elem().SetBool(false)
                                   }
                               }
           }
       }

}

func StructTravel (v interface{}, 
                   qi QueryInterface, 
                   f func(v reflect.Value, 
                          t reflect.StructField, 
                          qi QueryInterface)) interface{} {
    values := reflect.ValueOf(v).Elem()
    vtypes := reflect.TypeOf(v).Elem()
    for i := 0; i < values.NumField(); i++ {
        f(values.Field(i), vtypes.Field(i), qi)
    }
    return v
}
