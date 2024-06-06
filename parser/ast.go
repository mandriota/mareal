package parser

import (
	"fmt"
	"strconv"
)

type Buff []*Node

func (b *Buff) Add(n ...*Node) {
	*b = append(*b, n...)
}

func (b *Buff) Sub(n *Node) error {
	if len(*b) == 0 {
		return fmt.Errorf("trying to pop from empty buffer")
	}
	
	*n = *(*b)[len(*b)-1]
	*b = (*b)[:len(*b)-1]
	return nil
}

type Node struct {
	Component Buff
	Token
}

func (n Node) String() (s string) {
	switch n.Typ {
	case TkStr, TkIdent:
		return n.Val.(string)
	case TkNum:
		return strconv.FormatFloat(n.Val.(float64), 'f', -1, 64)
	case TkArr:
		for _, el := range n.Component {
			s += el.String()
		}
	}

	return
}
