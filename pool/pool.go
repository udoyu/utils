package pool

import (
    "sync"
    "sync/atomic"
)

type Pool struct {
    pool *sync.Pool
    size int64
}

func (this *Pool) Get() interface{} {
    v = this.pool.Get()
    if v != nil {
        atomic.AddInt64(&this.size, int64(-1))
    }
    return v
}

func (this *Pool) Put(v interface{}) {
    atomic.AddInt64(&this.size, int64(1))
    this.pool.Put(v)
}

func (this *Pool) Size() int64 {
    return this.size
}

func (this *Pool) Empty() bool {
    return this.size == 0
}

