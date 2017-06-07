package next

type quantifierDefinition struct {
	name     string
	commit   CommitType
	min, max int
	item     string
}

type quantifierParser struct {
	name     string
	commit   CommitType
	min, max int
	item     parser
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

func (d *quantifierDefinition) parser(r *registry) (parser, error) {
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
		item, err = itemDefinition.parser(r)
		if err != nil {
			return nil, err
		}
	}

	qp.item = item
	return qp, nil
}

func (d *quantifierDefinition) commitType() CommitType {
	return d.commit
}

func (p *quantifierParser) nodeName() string { return p.name }

func (p *quantifierParser) parse(t Trace, c *context, excluded []string) {
	t = t.Extend(p.name)
	t.Out1("parsing quantifier", c.offset)

	if p.commit&Documentation != 0 {
		t.Out1("fail, doc")
		c.fail(c.offset)
		return
	}

	if stringsContain(excluded, p.name) {
		t.Out1("excluded")
		c.fail(c.offset)
		return
	}

	var node *Node
	if m, ok := c.fromCache(p.name); ok {
		t.Out1("found in cache, match:", m)

		if !m {
			return
		}

		node = c.node
		if m, _ = c.fromCache(p.item.nodeName()); !m || c.node.tokenLength() <= node.tokenLength() {
			t.Out1("no matching item found in cache")
			t.Out2("offset before storing node", c.offset)
			c.success(node)
			t.Out2("offset after storing node", c.offset)
			return
		}

		node.clear()
		node.append(c.node)
	}

	node = newNode(p.name, p.commit, c.offset, c.offset)
	excluded = append(excluded, p.name)

	for {
		if p.max >= 0 && node.nodeLength() == p.max {
			t.Out1("success, max reached")
			c.cache.set(node.from, p.name, node)
			c.success(node)
			return
		}

		if node.tokenLength() > 0 {
			excluded = nil
		}

		p.item.parse(t, c, excluded)
		if !c.match {
			if node.nodeLength() >= p.min {
				t.Out1("success, no more match")
				c.cache.set(node.from, p.name, node)
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
