package pool

import (
	"reflect"
)

func New(v reflect.Value) interface{} {
        r := reflect.New(v.Type())
        if v.Kind() == reflect.Ptr {
                r.Elem().Set(reflect.New(v.Elem().Type()))
        }
        return r.Elem().Interface()
}

type PoolElem interface{
	Close()
}

type Pool struct {
	elemValue reflect.Value
	elems chan PoolElem
	maxSize int32
}

func NewPool(elem PoolElem, maxSize int32) *Pool {
	return &Pool{
		elemValue : reflect.ValueOf(elem),
		elems : make(chan PoolElem, maxSize),
		maxSize : maxSize,
	}
}

func (this *Pool) Put(elem PoolElem) {
	select {
		case this.elems <- elem :
			break
		default :
			elem.Close()
	}
}

func (this *Pool) Get() PoolElem {
	select {
		case e := <- this.elems :
			return e
		default :
			return New(this.elemValue).(PoolElem)
	}
}