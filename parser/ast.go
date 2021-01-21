package parser

type buff []*Node
type Node struct {
	Component buff
	*Token
}

func (b *buff) Add(n ...*Node) {
	*b = append(*b, n...)
}
