package router

import (
	"sync"
)

type LockParam struct {
	rw_lock *sync.RWMutex
	v       interface{}
}

func NewLockParam(v interface{}) *LockParam {
	l := &LockParam{rw_lock: new(sync.RWMutex), v: v}
	return l
}

func (this *LockParam) Value() interface{} {
	this.rw_lock.RLock()
	defer this.rw_lock.RUnlock()
	return this.v
}

func (this *LockParam) SetValue(v interface{}) {
	this.rw_lock.Lock()
	defer this.rw_lock.Unlock()
	this.v = v
}
