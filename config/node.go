package config

type Node struct {
	name  string
	value interface{}
	root  *Node
	child map[string]*Node
}

func NewNode(name string, data interface{}) *Node {
	return &Node{
		name:  name,
		value: data,
		child: make(map[string]*Node),
	}
}
func (this Node) Root() *Node {
	return this.root
}
func (this Node) AllChild() map[string]*Node {
	return this.child
}
func (this Node) Child(key string) *Node {
	n, ok := this.child[key]
	if !ok {
		return nil
	}
	return n
}
func (this Node) Search(keys ...string) *Node {
	keyarr := []string(keys)
	if len(keyarr) == 1 {
		return this.Child(keyarr[0])
	} else if len(keyarr) > 1 {
		c := this.Child(keyarr[0])
		if c == nil {
			return nil
		}
		return c.Search(keyarr[1:]...)
	} else {
		return nil
	}
	return nil
}
func (this *Node) Value() interface{} {
	return this.value
}
func (this *Node) SetRoot(root *Node) {
	this.root = root
}
func (this *Node) AddChild(key string, c *Node) *Node {
	c.root = this
	this.child[key] = c
	return c
}
func (this *Node) SetValue(value interface{}) {
	this.value = value
}
