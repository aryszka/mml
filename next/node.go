package next

import "fmt"

type CommitType int

const (
	None  CommitType = 0
	Alias CommitType = 1 << iota
)

type Node struct {
	Name       string
	Nodes      []*Node
	from       int
	to         int
	tokens     []rune
	commitType CommitType
}

func newNode(name string, ct CommitType, from, to int) *Node {
	return &Node{
		Name:       name,
		from:       from,
		to:         to,
		commitType: ct,
	}
}

func (n *Node) startsWith(p *Node) bool {
	// if n == p {
	// 	return true
	// }

	for _, ni := range n.Nodes {
		if ni.startsWith(p) {
			return true
		}

		if ni.to != n.from {
			return false
		}
	}

	return false
}

func (n *Node) appendRange(from, to int) {
	if n.from == 0 && n.to == 0 {
		n.from = from
	}

	n.to = to
}

func (n *Node) appendNode(p *Node) {
	n.Nodes = append(n.Nodes, p)
	n.appendRange(p.from, p.to)
}

func (n *Node) clear() {
	n.from = 0
	n.to = 0
	n.Nodes = nil
}

func (n *Node) commit() {
	return

	var nodes []*Node
	for _, ni := range n.Nodes {
		ni.commit()
		if ni.commitType&Alias != 0 {
			nodes = append(nodes, ni.Nodes...)
		} else {
			nodes = append(nodes, ni)
		}
	}

	n.Nodes = nodes
}

func (n *Node) applyTokens(t []rune) {
	n.tokens = t
	for _, ni := range n.Nodes {
		ni.applyTokens(t)
	}
}

func (n *Node) len() int {
	return n.to - n.from
}

func (n *Node) String() string {
	if n.from >= len(n.tokens) || n.to > len(n.tokens) {
		return n.Name + "incomplete:"
	}

	return fmt.Sprintf("%s:%d:%d:%s", n.Name, n.from, n.to, string(n.tokens[n.from:n.to]))
}
