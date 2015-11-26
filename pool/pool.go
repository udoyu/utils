package pool

import (
	"sync/atomic"
)

type PoolElem interface{
	Close()
}

type Pool struct {
	callback func() PoolElem
	elems chan PoolElem
	maxIdle int32
	maxActive int32
	curActive int32
}

func NewPool(callback func() PoolElem, maxIdle, maxActive int32) Pool {
	return Pool{
		callback : callback,
		elems : make(chan PoolElem, maxIdle),
		maxIdle : maxIdle,
		maxActive : maxActive,
	}
}

func (this *Pool) Put(elem PoolElem) {
	select {
		case this.elems <- elem :
			break
		default :
			atomic.AddInt32(&this.curActive, -1)
			elem.Close()
	}
}

func (this *Pool) Get() PoolElem {
	select {
		case e := <- this.elems :
			return e
		default :
			ca := atomic.AddInt32(&this.curActive, 1)
			if ca < this.maxActive {
				return this.callback()
			}
	}
	return nil
}