package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type ParseInterface interface {
	IsKey(string) (key string, isEnd, ok bool) //isEnd表示
	Parse(*Node, string) error
}

type Parser struct {
	node *Node
	pi   ParseInterface
}

func NewParser(pi ParseInterface) *Parser {
	return &Parser{
		node: NewNode("", nil),
		pi:   pi,
	}
}

func (this *Parser) Node() *Node {
	return this.node
}

func (this *Parser) ParseString(str string) error {
	if key, isEnd, ok := this.pi.IsKey(str); ok {
		if !isEnd {
			if key == "" {
				return errors.New("Error 1001:key is empty")
			}
			if this.node == nil {
				this.node = NewNode(key, nil)
			} else {
				this.node = this.node.AddChild(key, NewNode(key, nil))
			}
		} else {
			if this.node.name != key {
				return errors.New("Error 1002:key is wrong")
			} else {
				this.node = this.node.root
			}
		}
		return nil
	} else {
		return this.pi.Parse(this.node, str)
	}

	return nil
}

func (this *Parser) LoadFile(filename string) error {
	fd, err := os.Open(filename)
	if nil != err {
		return err
	}
	defer fd.Close()
	buf := bufio.NewReader(fd)
	flag := true
	for flag {
		line, err := buf.ReadString('\n')
		if io.EOF == err {
			flag = false
		}

		line = strings.TrimLeft(line, " ")
		line = strings.TrimLeft(line, "\t")
		line = strings.TrimRight(line, " ")
		line = strings.TrimRight(line, "\t")
		if len(line) < 2 || '#' == line[0] || '\n' == line[0] || '\r' == line[0] {
			continue
		}

		length := len(line)
		if line[length-2] == '\r' {
			length -= 1
			line = line[:length-1]
		}
		err = this.ParseString(line)
		if err != nil {
			return err
		}
	}
	return nil
}
