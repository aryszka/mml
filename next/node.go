package next

type Node struct {
	Name   string
	Nodes  []*Node
	From   int
	To     int
	tokens []rune
}

func newNode(name string, from, to int) *Node {
	return &Node{
		Name: name,
		From: from,
		To:   to,
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

func (n *Node) applyTokens(t []rune) {
	n.tokens = t
	for _, ni := range n.Nodes {
		ni.applyTokens(t)
	}
}

func (n *Node) String() string {
	if n.From >= len(n.tokens) || n.To > len(n.tokens) {
		return n.Name + "incomplete:"
	}

	return n.Name + ":" + string(n.tokens[n.From:n.To])
}
