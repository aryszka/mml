package mml

import "errors"

type tokenCache struct {
	match   *intSet
	noMatch *intSet
	nodes   []*node // TODO: potential optimization can be to use a balanced binary tree
}

type cache struct {
	tokens []*tokenCache // TODO: try with pointers, too
}

var errDamagedCache = errors.New("damaged token/node cache")

func (c *cache) get(offset int, t nodeType) (*node, bool, bool) {
	if len(c.tokens) <= offset {
		return nil, false, false
	}

	tc := c.tokens[offset]

	if tc.noMatch.has(t) {
		return nil, false, true
	}

	if tc.match.has(t) {
		for _, n := range tc.nodes {
			if n.typ == t {
				return n, true, true
			}
		}

		panic(errDamagedCache)
	}

	return nil, false, false
}

func (c *cache) set(offset int, t nodeType, n *node, match bool) {
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
		tc = &tokenCache{
			match:   &intSet{},
			noMatch: &intSet{},
		}
		c.tokens[offset] = tc
	}

	if !match {
		// common use case leaked in. The reason is that this check is required in all current use
		// cases. E.g. there can be a group cached already, which can be the first item of a longer
		// group, and not matching the longer group should not overwrite the cached match of the
		// shorter.
		if tc.match.has(t) {
			return
		}

		tc.noMatch.set(t)
		c.tokens[offset] = tc
		return

		// TODO: there was a missing return here
	}

	if tc.match.has(t) {
		for i, ni := range tc.nodes {
			if ni.typ == t {
				tc.nodes[i] = n
				c.tokens[offset] = tc
				return
			}
		}

		panic(errDamagedCache)
	}

	tc.match.set(t)
	tc.nodes = append(tc.nodes, n)
}

func (c *cache) clear() {
	c.tokens = nil
}
