package router

import (
    "sync"
)

type HostInterface interface {
    GetObj() interface{}
    PutObj(v interface{}) 
    Init()
    Status() bool
    ChangeStatus()
    Size() int64
    Empty() bool
}

type Host struct {
    pool *sync.Pool
    StatusParam *LockParam
}

type HostGroup struct {
    rw_lock *sync.RWMutex
    Hosts []HostInterface
}

func (this *Host) GetObj() interface{} {
    return this.pool.Get()
}

func (this *Host) PubObj(obj interface{}) {
    this.pool.Put(obj)
}

func (this *Host) Init() {
    this.pool = new(sync.Pool)
    this.StatusParam = NewLockParam(false) 
}

func (this *Host) Status() bool {
    return this.StatusParam.Value().(bool)
}

func (this *Host) ChangeStatus() {
    this.StatusParam.SetValue(!this.Status())
}

func NewHostGroup() *HostGroup {
    v := &HostGroup{rw_lock:new(sync.RWMutex)}
    return v
}

func (this *HostGroup) AddHost(v HostInterface) {
    this.rw_lock.Lock()
    defer this.rw_lock.Unlock()
    this.Hosts = append(this.Hosts, v)
}

func (this *HostGroup) GetHost() HostInterface {
    this.rw_lock.RLock()
    defer this.rw_lock.RUnlock()
    for _,h := range this.Hosts {
        if h.Status() && !h.Empty() {
            return h
        }
    }
    return nil
}

