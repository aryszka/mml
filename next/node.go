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

func (n *Node) startsWith(p *Node) bool {
	if n == p {
		return true
	}

	for _, ni := range n.Nodes {
		if ni.startsWith(p) {
			return true
		}

		if ni.To != n.From {
			return false
		}
	}

	return false
}

func (n *Node) appendRange(from, to int) {
	if n.From < 0 {
		n.From = from
	}

	n.To = to
}

func (n *Node) appendNode(p *Node) {
	n.Nodes = append(n.Nodes, p)
	n.appendRange(p.From, p.To)
}
