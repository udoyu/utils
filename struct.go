package utils

import (
	"reflect"
	"strconv"
)

func String(s string) *string {
	return &s
}

func Int32(i int32) *int32 {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func Uint32(i uint32) *uint32 {
	return &i
}

func Uint64(i uint64) *uint64 {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

type QueryInterface interface {
	Query(t reflect.StructField) string
}

func StructTravelSet(v interface{}, qi QueryInterface) interface{} {
	return StructTravel(v, qi, StructTravelFunc)
}

/*
example:
type QueryI struct {
	req http.Request
}

func (this QueryI) Query(t reflect.StructField) string {
	form := utils.StringCutRightExp(t.Tag.Get("json"), ",", 1)
	return this.FormValue(form)
}

 */





func StructTravel(v interface{},
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

func StructTravelFunc(v reflect.Value, t reflect.StructField, qi QueryInterface) {
	if v.Kind() == reflect.Ptr {
		if v.NumMethod() == 0 {
			if str := qi.Query(t); len(str) != 0 && v.CanSet() {
				//form := utils.StringCutRightExp(t.Tag.Get("json"), ",", 1)
				switch v.Type() {
				case STRING_PTR:
					v.Set(reflect.ValueOf(String(str)))
				case INT32_PTR:
					v.Set(reflect.ValueOf(Int32(int32(Atoi(str)))))
				case INT64_PTR:
					v.Set(reflect.ValueOf(Int64(Atol(str))))
				case UINT32_PTR:
					v.Set(reflect.ValueOf(Uint32(uint32(Atoi(str)))))
				case UINT64_PTR:
					v.Set(reflect.ValueOf(Uint64(uint64(Atol(str)))))
				case BOOL_PTR:
					if str == "true" {
						v.Set(reflect.ValueOf(Bool(true)))
					} else {
						v.Set(reflect.ValueOf(Bool(false)))
					}
				}
			}
		}
	} else {
		if v.NumMethod() == 0 {
			if str := qi.Query(t); len(str) != 0 && v.CanSet() {
				//form := utils.StringCutRightExp(t.Tag.Get("json"), ",", 1)
				switch v.Type() {
				case STRING:
					v.Set(reflect.ValueOf(str))
				case INT32:
					v.Set(reflect.ValueOf(int32(Atoi(str))))
				case INT64:
					v.Set(reflect.ValueOf(Atol(str)))
				case UINT32:
					v.Set(reflect.ValueOf(uint32(Atoi(str))))
				case UINT64:
					v.Set(reflect.ValueOf(uint64(Atol(str))))
				case BOOL:
					if str == "true" {
						v.Set(reflect.ValueOf(true))
					} else {
						v.Set(reflect.ValueOf(false))
					}
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

func Atoul(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

func Atou(s string) uint {
	return uint(Atoul(s))
}

var (
	STRING_PTR = reflect.TypeOf(String(""))
	INT32_PTR  = reflect.TypeOf(Int32(int32(0)))
	INT64_PTR  = reflect.TypeOf(Int64(int64(0)))
	UINT32_PTR = reflect.TypeOf(Uint32(uint32(0)))
	UINT64_PTR = reflect.TypeOf(Uint64(uint64(0)))
	BOOL_PTR   = reflect.TypeOf(Bool(false))
	
	STRING = reflect.TypeOf("")
    INT32  = reflect.TypeOf(int32(0))
    INT64  = reflect.TypeOf(int64(0))
    UINT32 = reflect.TypeOf(uint32(0))
    UINT64 = reflect.TypeOf(uint64(0))
    BOOL   = reflect.TypeOf(false)
)
