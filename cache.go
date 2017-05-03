package mml

import "errors"

type cacheItem struct {
	typ  nodeType
	node *node
}

type tokenCache struct {
	match   *intSet
	noMatch *intSet
	items   []*cacheItem // TODO: potential optimization can be to use a balanced binary tree
}

type cache struct {
	tokens []*tokenCache // TODO: try with pointers, too
}

var errDamagedCache = errors.New("damaged token/node cache")

// TODO: reconsider using values instead of pointers

func (c *cache) get(offset int, t nodeType) (*node, bool, bool) {
	if len(c.tokens) <= offset {
		return nil, false, false
	}

	tc := c.tokens[offset]
	if tc == nil {
		return nil, false, false
	}

	if tc.noMatch.has(t) {
		return nil, false, true
	}

	if tc.match.has(t) {
		for _, i := range tc.items {
			if i.typ == t {
				return i.node, true, true
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
		return

		// TODO: there was a missing return here
	}

	if tc.match.has(t) {
		for _, ii := range tc.items {
			if ii.typ == t {
				ii.node = n
				return
			}
		}

		panic(errDamagedCache)
	}

	tc.match.set(t)
	tc.items = append(tc.items, &cacheItem{
		typ:  t,
		node: n,
	})
}

func (c *cache) clear() {
	c.tokens = nil
}
