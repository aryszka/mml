package mml

type parserResult struct {
	accepting bool
	valid     bool
	node      *node
	fromCache bool
	unparsed  *tokenStack
}

func (r *parserResult) fillFromCache(
	c *cache,
	typ nodeType,
	t *token,
	init *node,
	initIsMember bool,
	initIsItemMember bool,
) bool {
	ct := t
	if init != nil {
		ct = init.token
	}

	n, m, ok := c.get(ct.offset, typ)
	if !ok {
		return false
	}

	if !m {
		if r.unparsed == nil {
			r.unparsed = newTokenStack()
		}

		r.valid = false
		r.unparsed.push(t)
		r.fromCache = true
		r.accepting = false
		return true
	}

	// optional and choice:
	if initIsMember && init != nil && n != init {
		return false
	}

	// repetition and sequence:
	if initIsItemMember && init != nil && (len(n.nodes) == 0 || n.nodes[0] != init) {
		return false
	}

	if r.unparsed == nil {
		r.unparsed = newTokenStack()
	}

	r.valid = true
	r.node = n
	r.unparsed.push(t)
	r.fromCache = true
	r.accepting = false
	return true
}

func (r *parserResult) ensureNode(name string, typ nodeType) {
	if r.node != nil {
		return
	}

	r.node = &node{
		name: name,
		typ:  typ,
	}
}

func (r *parserResult) ensureUnparsed() {
	if r.unparsed == nil {
		r.unparsed = newTokenStack()
	}
}

func (r *parserResult) assertUnparsed(name string) {
	if r.unparsed == nil || !r.unparsed.has() {
		panic(unexpectedResult(name))
	}
}

func (r *parserResult) mergeStack(s *tokenStack) {
	if s == nil {
		return
	}

	r.ensureUnparsed()
	r.unparsed.merge(s)
}

func checkSkip(skip int, done bool) (int, bool, bool) {
	if skip == 0 {
		return 0, false, false
	}

	skip--
	if skip > 0 || !done {
		return skip, true, true
	}

	if done {
		return 0, false, true
	}

	return 0, false, false
}
