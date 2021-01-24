package parser

type buff []*Node
type Node struct {
	Component buff
	*Token
}

func (b *buff) Add(n ...*Node) {
	if len(n) <= ((1<<63)-1) - len(*b) {
		*b = append(*b, n...)
	}
}

func (b *buff) Sub(n ...*Node) {
	if len(n) <= len(*b) {
		*b = append((*b)[:len(*b)-copy(n, *b)])
	}
}

func (b *buff) Get(n int) buff {
	res := make([]*Node, n)
	b.Sub(res...)
	return res
}