package utils

import (
	"sync"
)

type SafeSet struct {
	set  *Set
	lock *sync.RWMutex
}

func NewSafeSet() *SafeSet {
	s := &SafeSet{
		set:  NewSet(),
		lock: new(sync.RWMutex),
	}
	return s
}

func (this *SafeSet) Insert(v interface{}) {
	this.lock.Lock()
	this.set.Insert(v)
	this.lock.Unlock()
}

func (this *SafeSet) Has(v interface{}) bool {
	this.lock.RLock()
	b := this.set.Has(v)
	this.lock.RUnlock()
	return b
}

func (this *SafeSet) Remove(v interface{}) {
	this.lock.Lock()
	this.set.Remove(v)
	this.lock.Unlock()
}

func (this *SafeSet) Range(callback SetCallback, vs ...interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	this.set.Range(callback, vs...)
}

func (this *SafeSet) Size() int {
	this.lock.RLock()
	size := this.set.Size()
	this.lock.RUnlock()
	return size
}

func (this *SafeSet) ToSlice() []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.set.ToSlice()
}
