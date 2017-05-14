package next

type cacheItem struct {
	name string
	node *Node
}

type tokenCache struct {
	match   []*cacheItem // TODO: potential optimization can be to use a balanced binary tree
	noMatch []string
}

type cache struct {
	tokens []*tokenCache // TODO: try with pointers, too
}

func (c *cache) get(offset int, name string) (*Node, bool, bool) {
	if len(c.tokens) <= offset {
		return nil, false, false
	}

	tc := c.tokens[offset]
	if tc == nil {
		return nil, false, false
	}

	for _, i := range tc.noMatch {
		if i == name {
			return nil, false, true
		}
	}

	for _, i := range tc.match {
		if i.name == name {
			return i.node, true, true
		}
	}

	return nil, false, false
}

func (c *cache) set(offset int, name string, n *Node) {
	if len(c.tokens) <= offset {
		if cap(c.tokens) > offset {
			c.tokens = c.tokens[:offset+1]
		} else {
			c.tokens = c.tokens[:cap(c.tokens)]
			for len(c.tokens) <= offset {
				c.tokens = append(c.tokens, nil)
			}
		}
	}

	tc := c.tokens[offset]
	if tc == nil {
		tc = &tokenCache{}
		c.tokens[offset] = tc
	}

	if n == nil {
		for _, i := range tc.match {
			if i.name == name {
				return
			}
		}

		for _, i := range tc.noMatch {
			if i == name {
				return
			}
		}

		tc.noMatch = append(tc.noMatch, name)
		return
	}

	for _, i := range tc.match {
		if i.name == name {
			i.node = n
			return
		}
	}

	tc.match = append(tc.match, &cacheItem{
		name: name,
		node: n,
	})
}

func (c *cache) clear() {
	c.tokens = nil
}
