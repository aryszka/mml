package next

type quantifierDefinition struct {
	name     string
	commit   CommitType
	min, max int
	item     string
}

type quantifierParser struct {
	name       string
	commit     CommitType
	min, max   int
	item       parser
	includedBy []parser
}

func newQuantifier(name string, ct CommitType, item string, min, max int) *quantifierDefinition {
	return &quantifierDefinition{
		name:   name,
		commit: ct,
		min:    min,
		max:    max,
		item:   item,
	}
}

func (d *quantifierDefinition) nodeName() string { return d.name }

func (d *quantifierDefinition) parser(r *registry, path []string) (parser, error) {
	if stringsContain(path, d.name) {
		panic(errCannotIncludeParsers)
	}

	p, ok := r.parser(d.name)
	if ok {
		return p, nil
	}

	qp := &quantifierParser{
		name:   d.name,
		commit: d.commit,
		min:    d.min,
		max:    d.max,
	}

	r.setParser(qp)

	item, ok := r.parser(d.item)
	if !ok {
		itemDefinition, ok := r.definition(d.item)
		if !ok {
			return nil, parserNotFound(d.item)
		}

		var err error
		item, err = itemDefinition.parser(r, path)
		if err != nil {
			return nil, err
		}
	}

	qp.item = item
	return qp, nil
}

func (d *quantifierDefinition) commitType() CommitType { return d.commit }
func (p *quantifierParser) nodeName() string           { return p.name }

// TODO: merge the quantifier into the sequence
// DOC: sequences are hungry and are not revisited, a*a cannot match anything.
// DOC: how to match a tailing a? (..)*a | .(..)*a

func (p *quantifierParser) setIncludedBy(i parser, path []string) {
	if stringsContain(path, p.name) {
		panic(errCannotIncludeParsers)
	}

	p.includedBy = append(p.includedBy, i)
}

func (p *quantifierParser) cacheIncluded(*context, *Node) {
	panic(errCannotIncludeParsers)
}

func (p *quantifierParser) parse(t Trace, c *context) {
	t = t.Extend(p.name)
	t.Out1("parsing quantifier", c.offset)

	if p.commit&Documentation != 0 {
		t.Out1("fail, doc")
		c.fail(c.offset)
		return
	}

	if c.excluded(c.offset, p.name) {
		t.Out1("excluded")
		c.fail(c.offset)
		return
	}

	c.exclude(c.offset, p.name)
	defer c.include(c.offset, p.name)

	node := newNode(p.name, p.commit, c.offset, c.offset)

	// this way of checking the cache definitely needs the testing of the russ cox form
	for {
		if p.max >= 0 && node.nodeLength() == p.max {
			t.Out1("success, max reached")
			c.cache.set(node.from, p.name, node)
			for _, i := range p.includedBy {
				i.cacheIncluded(c, node)
			}

			c.success(node)
			return
		}

		t.Out2("next quantifier item")

		// n, m, ok := c.cache.get(c.offset, p.item.nodeName())
		m, ok := c.fromCache(p.item.nodeName())
		if ok {
			t.Out1("quantifier item found in cache, match:", m, c.offset, c.node.tokenLength())
			if m {
				node.append(c.node)
				if c.node.tokenLength() > 0 {
					t.Out2("taking next after cached found")
					continue
				}
			}

			if node.nodeLength() >= p.min {
				t.Out1("success, no more match")
				c.cache.set(node.from, p.name, node)
				for _, i := range p.includedBy {
					i.cacheIncluded(c, node)
				}

				c.success(node)
			} else {
				t.Out1("fail, min not reached")
				c.cache.set(node.from, p.name, nil)
				c.fail(node.from)
			}

			return
		}

		p.item.parse(t, c)
		if !c.match || c.node.tokenLength() == 0 {
			if node.nodeLength() >= p.min {
				t.Out1("success, no more match")
				c.cache.set(node.from, p.name, node)
				for _, i := range p.includedBy {
					i.cacheIncluded(c, node)
				}

				c.success(node)
			} else {
				t.Out1("fail, min not reached")
				c.cache.set(node.from, p.name, nil)
				c.fail(node.from)
			}

			return
		}

		node.append(c.node)
	}
}
