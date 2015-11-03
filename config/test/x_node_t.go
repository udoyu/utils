package main

import (
	"fmt"
	"strings"
	"github.com/udoyu/utils/config"
)

/*
<root>
   <node>
		<n1>
		k1 = 1
		k2 = b
		</n1>
		<n2>
		k1 = 2
		k2 = c
		</n2>
   </node>
</root>
*/

type Node struct {
	key string
	value string
}

type Nodes []Node

func (this *Nodes) Parse(str string) error {
	fmt.Println(str)
	kv := strings.Split(str, "=")
	if len(kv) < 2 {return nil}
	*this = append(*this, Node{key:kv[0], value:kv[1]})
	return nil
}

func main () {
	xp := config.NewXNodeParse(&Nodes{})
	xp.LoadFile("xn.conf")
	node := xp.Node()
	node = node.Search("root", "node")
	if node != nil {
		//取n1节点
		n1 := node.Child("n1")
		if n1 != nil {
			fmt.Println(n1.Value())
		}
		
		//遍历node所有节点
		for k, v := range node.AllChild() {
			fmt.Println(k, v.Value())
		}
	}
}