package utils

import (
    "reflect"
    "strconv"
)

type QueryInterface interface {
    Query(t reflect.StructField) string
}

func StructTravelSet(v interface{}, qi QueryInterface) interface{} {
    return StructTravel(v, qi, StructTravelFunc)
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

func StructTravelFunc (v reflect.Value, t reflect.StructField, qi QueryInterface) {
       if v.Kind() == reflect.Ptr {
           v = v.Elem()
           if v.Kind() == reflect.Struct {
               StructTravel(v, qi, StructTravelFunc)
           } else if str:=qi.Query(t); len(str)!=0 && v.CanSet() {
               //form := utils.StringCutRightExp(t.Tag.Get("json"), ",", 1)
               switch v.Kind() {
                   case reflect.String: v.SetString(str)
                   case reflect.Int32 : v.SetInt(int64(Atoi(str)))
                   case reflect.Int64 : v.SetInt(Atol(str))
                   case reflect.Uint32 : v.SetUint(uint64(Atoi(str)))
                   case reflect.Uint64 : v.SetUint(uint64(Atol(str)))
                   case reflect.Bool : if str == "true" {
                                           v.Elem().SetBool(true)
                                       } else {
                                           v.Elem().SetBool(false)
                                       }
               }
           }
       }
}

func Atol(s string) int64 {
        i, _ := strconv.ParseInt(s, 10, 64)
        return i
}

func Atoi(s string) int {
        i, _ := strconv.Atoi(s)
        return i
}

