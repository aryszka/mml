package next

type Node struct {
	Name string
	Nodes []*Node
	From int
	To int
}

func newNode(name string, from, to int) *Node {
	return &Node{
		Name: name,
		From: from,
		To: to,
	}
}

func (n *Node) startsWith(*Node) bool {
	panic(ErrNotImplemented)
}

func (n *Node) appendNode(p *Node) {
	panic(ErrNotImplemented)
}

// func (n *Node) appendNode(p *Node) {
// 	n.Nodes = append(n.Nodes, p)
// 	n.appendToken(p.Tokens...)
// }
// 
// func (n *Node) appendToken(t ...*Token) {
// 	for _, ti := range t {
// 		n.Tokens = append(n.Tokens, ti)
// 		if n.Reference == nil {
// 			n.Reference = ti
// 		}
// 	}
// }
// 
// func (n *Node) startsWith(p *Node) bool {
// 	if n == p {
// 		return true
// 	}
// 
// 	if len(n.Nodes) == 0 {
// 		return false
// 	}
// 
// 	for _, ni := range n.Nodes {
// 		if ni.startsWith(p) {
// 			return true
// 		}
// 
// 		if len(ni.Tokens) > 0 {
// 			return false
// 		}
// 	}
// 
// 	return false
// }
