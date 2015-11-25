package pool

import (
)

type PoolElem interface{
	Close()
}

type Pool struct {
	callback func() PoolElem
	elems chan PoolElem
	maxSize int32
}

func NewPool(callback func() PoolElem, maxSize int32) Pool {
	return Pool{
		callback : callback,
		elems : make(chan PoolElem, maxSize),
		maxSize : maxSize,
	}
}

func (this Pool) Put(elem PoolElem) {
	select {
		case this.elems <- elem :
			break
		default :
			elem.Close()
	}
}

func (this Pool) Get() PoolElem {
	select {
		case e := <- this.elems :
			return e
		default :
			return this.callback()
	}
}