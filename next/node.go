package next

type Node struct {
	from, to int
}

func newNode(name string, ct CommitType, from, to int) *Node {
	return &Node{}
}

func (n *Node) tokenLength() int {
	return 0
}

func (n *Node) nodeLength() int {
	return 0
}

func (n *Node) append(c *Node) {}

func (n *Node) clear() {}
