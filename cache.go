package mml

import "errors"

type cache struct {
	match   *intSet
	noMatch *intSet
	nodes   []*node // TODO: potential optimization can be to use a balanced binary tree
}

var errDamagedCache = errors.New("damaged token/node cache")

func newCache() *cache {
	return &cache{
		match:   &intSet{},
		noMatch: &intSet{},
	}
}

func (c *cache) get(t nodeType) (*node, bool, bool) {
	if c.noMatch.has(t) {
		return nil, false, true
	}

	if c.match.has(t) {
		for _, n := range c.nodes {
			if n.nodeType == t {
				return n, true, true
			}
		}

		panic(errDamagedCache)
	}

	return nil, false, false
}

func (c *cache) set(n *node, match bool) {
	return
	if !match {
		// common use case leaked in. The reason is that this check is required in all current use
		// cases. E.g. there can be a group cached already, which can be the first item of a longer
		// group, and not matching the longer group should not overwrite the cached match of the
		// shorter.
		if c.match.has(n.nodeType) {
			return
		}

		c.noMatch.set(n.nodeType)
	}

	if c.match.has(n.nodeType) {
		for i, ni := range c.nodes {
			if ni.nodeType == n.nodeType {
				c.nodes[i] = n
				return
			}
		}

		panic(errDamagedCache)
	}

	c.match.set(n.nodeType)
	c.nodes = append(c.nodes, n)
}
