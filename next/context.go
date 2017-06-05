package next

type context struct {
	offset int
	match  bool
	node   *Node
}

func (c *context) token() (rune, bool) {
	return 0, false
}

func (c *context) checkCache(name string) bool {
	return false
}

func (c *context) success(n *Node) {
}

func (c *context) fail(name string, offset int) {
}
