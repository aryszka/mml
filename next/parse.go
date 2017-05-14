package next

type definition interface {
	nodeName() string
	member(string) (bool, error)
	generator(trace, string, []string) (generator, error)

	// TODO: try to do this during validation of the generators
	// terminates([]string) bool
}

type generator interface {
	nodeName() string
	valid() bool
	validate(trace, []string) error
	parser(trace, *Node) parser
}

type parser interface {
	nodeName() string
	parse(*parserContext)
}

type parserContext struct {
	cache *cache
	token *Token
	accepting bool
	valid     bool
	node      *Node
	unparsed  *stack
	fromCache bool
}

func stringsContain(ss []string, s string) bool {
	for _, si := range ss {
		if si == s {
			return true
		}
	}

	return false
}

func (c *parserContext) nextToken() rune {
	panic(ErrNotImplemented)
}

func (c *parserContext) offset() int {
	panic(ErrNotImplemented)
}

func (c *parserContext) succeed(*Node) {
	panic(ErrNotImplemented)
}

func (c *parserContext) fail(string) {
	panic(ErrNotImplemented)
}

func (c *parserContext) fillFromCache(string, *Node, *Node) bool {
	panic(ErrNotImplemented)
}

func (c *parserContext) reverse(int) {
	panic(ErrNotImplemented)
}

// func (c *parserContext) initNext() {
// 	c.accepting = false
// 	c.valid = false
// 	c.node = nil
// 	c.fromCache = false
// 	c.unparsed.clear()
// }
// 
// func (c *parserContext) fillFromCache(name string, init, itemInit *Node) bool {
// 	t := c.token
// 
// 	if init != nil {
// 		t = init.Reference
// 	}
// 	
// 	if itemInit != nil {
// 		t = itemInit.Reference
// 	}
// 
// 	n, m, ok := c.cache.get(t.Offset, name)
// 	if !ok {
// 		return false
// 	}
// 
// 	if init != nil && n != init {
// 		return false
// 	}
// 
// 	if itemInit != nil && (len(n.Nodes) == 0 || n.Nodes[0] != itemInit) {
// 		return false
// 	}
// 
// 	c.valid = m
// 	c.node = n
// 	c.fromCache = true
// 	c.unparsed.push(c.token)
// 	return true
// }
// 
// func (c *parserContext) succeed(name string, n *Node, s *stack, accepting bool) {
// 	c.valid = true
// 	c.node = n
// 	c.accepting = accepting
// 
// 	if s != nil {
// 		c.unparsed.merge(s)
// 	}
// 
// 	c.cache.set(c.node.Reference.Offset, name, c.node)
// }
// 
// func (c *parserContext) fail(name string, s *stack, accepting bool) {
// 	c.accepting = accepting
// 
// 	if s != nil {
// 		c.unparsed.merge(s)
// 	}
// 
// 	c.unparsed.push(c.token)
// 
// 	// cache, reference, node nil
// }
// 
// func checkSkip(skip int, done bool) (int, bool, bool) {
// 	if skip == 0 {
// 		return 0, false, false
// 	}
// 
// 	skip--
// 	if skip > 0 || !done {
// 		return skip, true, true
// 	}
// 
// 	if done {
// 		return 0, false, true
// 	}
// 
// 	return 0, false, false
// }

// could be a reader
// could contain a stack of the whole read tokens and the rest just could use references to the tokens
// skip could become a single operation and there was no need for the stack
