package router

import (
    "sync"
)

type HostInterface interface {
    func GetObj() interface{}
    func PutObj(v interface{}) 
    func Init()
    func Status() bool
    func ChangeStatus()
    func Size() int64
    func Empty() bool
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
    this.pool.Pub(obj)
}

func (this *Host) Init() {
    this.pool = new(Pool)
    this.StatusParam = NewLockParam(false) 
}

func (this *Host) Status() bool {
    return this.StatusParam.Value().(bool)
}

func (this *Host) ChangeStatus() {
    return this.StatusParam.SetValue(!this.Status())
}

func NewHostGroup() *HostGroup {
    v := &HostGroup{rw_lock:new(sync.RWMutex)}
    return v
}

func (this *HostGroup) AddHost(v HostInterface) {
    this.rw_lock.WLock()
    defer this.rw_lock.WUnlock()
    hosts = append(this.Hosts, v)
    this.Hosts = hosts
}

func (this *HostGroup) GetHost() HostInterface {
    this.rw_lock.RLock()
    defer this.rw_lock.RUnock()
    for _,h := range this.Hosts {
        if h.Status() && !h.Empty() {
            return h
        }
    }
    return nil
}

