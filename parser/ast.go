package parser

type buff []*Node
type Node struct {
	Component buff
	*Token
}

func (b *buff) Add(n ...*Node) {
	*b = append(*b, n...)
}

func (b *buff) Sub(n ...*Node) {
	if len(n) <= len(*b) {
		for i, el := range n {
			*el = *(*b)[len(*b)-len(n)+i]
		}
		*b = (*b)[:len(n)-1]
	}
}