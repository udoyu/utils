package router
import (
    "sync"
)

type RouterInterface {
    func GetObj(rtype string, index int) interface{}
    func Init(v interface{}) int
}

//map[rtype]index[*hosts]
type HostGroupIndexMap map[int][*HostGroup]
type HostGroupTypeIndexMap map[string]HostGroupIndexMap

type Router struct {
    hostGroupTypeIndex HostGroupTypeIndexMap
}

func (this *Router) Init() {
    this.hostGroupTypeIndex = make(HostGroupTypeIndexMap)
}

func (this *Router) getHostGroup(rtype string, index int) *HostGroup {
    m,ok := this.hostGroupTypeIndex[rtype]
    if ok {
        v,ok1 := m[index]
        if ok1 {
            return v
        }
    }
    return nil
}

func (this *Router) getHost(rtype string, index int) HostInterface {
    hg := this.getHostGroup(rtype, index)
    if hg != nil {
        return hg.GetHost()
    }
    return nil
}

func (this *Router) GetObj(rtype string, index int) interface{} {
    v := this.getHost(rtype, index)
    if v != nil {
        return v.GetObj()
    }
    return nil
}

