package config

import (
	"reflect"
	"github.com/udoyu/utils"
)

/*
<root>
   <node>
		<n1>
		k1 = 1
		k2 = b
		</n1>
		<n2>
		k1 = 1
		k2 = b
		</n2>
   </node>
</root>
 */

type XNodeInterface interface{
	Parse(string) error
}

type XNodeParse struct {
	*Parser
}

func NewXNodeParse(xi XNodeInterface) *XNodeParse {
	return &XNodeParse {
		Parser : NewParser(newXParserI(xi)),
	}
}

type xParserI struct {
	xi reflect.Value
}

func newXParserI(xi XNodeInterface) xParserI {
	return xParserI{
		xi : reflect.ValueOf(xi),
	}
}

func (this xParserI) IsKey(str string) (key string, isEnd, ok bool) {
	if len(str) > 2 && str[0] == '<' && str[len(str)-1] == '>' {
		ok = true
		if str[1] != '/' {
			key = str[1 : len(str)-1]
			
		} else {
			key = str[2 : len(str)-1]
			isEnd = true
		}
	}
	return key, isEnd, ok
}

func (this xParserI) Parse(node *Node, str string) error {
	if str == "" {
		return nil
	} 
	if node.Value() == nil {
		node.SetValue(utils.New(this.xi))
	}
	return node.Value().(XNodeInterface).Parse(str)
}