package utils

type SetCallback func(v interface{}, vs ...interface{})

type SetInterface interface {
	Insert(v interface{})
	Has(v interface{}) bool
	Remove(v interface{})
	Range(callback SetCallback, vs ...interface{})
	Size() int
	ToSlice() []interface{}
	Wake() //防止读锁一直占用，而写锁被卡住
}

type Set struct {
	v map[interface{}]struct{}
}

func NewSet() *Set {
	s := &Set{
		v: make(map[interface{}]struct{}),
	}
	return s
}

func (this *Set) Insert(v interface{}) {
	this.v[v] = struct{}{}
}

func (this *Set) Has(v interface{}) bool {
	_, ok := this.v[v]
	return ok
}

func (this *Set) Remove(v interface{}) {
	delete(this.v, v)
}

func (this *Set) Range(callback SetCallback, vs ...interface{}) {
	for k, _ := range this.v {
		callback(k, vs...)
	}
}

func (this *Set) Size() int {
	return len(this.v)
}

func (this *Set) ToSlice() []interface{} {
	arr := make([]interface{}, 0, len(this.v))
	for k, _ := range this.v {
		arr = append(arr, k)
	}
	return arr
}

func (this *Set) Wake() {

}
