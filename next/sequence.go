package next

type sequenceDefinition struct {
	name   string
	commit CommitType
	items  []string
}

type sequenceParser struct {
	name   string
	commit CommitType
	items  []parser
}

func newSequence(name string, ct CommitType, items []string) *sequenceDefinition {
	return &sequenceDefinition{
		name:   name,
		commit: ct,
		items:  items,
	}
}

func (d *sequenceDefinition) nodeName() string { return d.name }

func (d *sequenceDefinition) parser(r *registry) (parser, error) {
	p, ok := r.parser(d.name)
	if ok {
		return p, nil
	}

	sp := &sequenceParser{
		name:   d.name,
		commit: d.commit,
	}

	r.setParser(sp)

	var items []parser
	for _, i := range d.items {
		item, ok := r.parser(i)
		if ok {
			items = append(items, item)
			continue
		}

		itemDefinition, ok := r.definition(i)
		if !ok {
			return nil, parserNotFound(i)
		}

		item, err := itemDefinition.parser(r)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	sp.items = items
	return sp, nil
}

func (d *sequenceDefinition) commitType() CommitType {
	return d.commit
}

func (p *sequenceParser) nodeName() string { return p.name }

func (p *sequenceParser) parse(t Trace, c *context, excluded []string) {
	t = t.Extend(p.name)
	t.Out1("parsing sequence", c.offset)

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

		if !m || len(p.items) == 0 {
			return
		}

		node = c.node
		if m, _ = c.fromCache(p.items[0].nodeName()); !m || c.node.tokenLength() <= node.tokenLength() {
			t.Out1("no matching item found in cache")
			t.Out2("offset before storing node", c.offset)
			c.success(node)
			t.Out2("offset after storing node", c.offset)
			return
		}

		node.clear()
		node.append(c.node)
	}

	excluded = append(excluded, p.name)
	items := p.items
	node = newNode(p.name, p.commit, c.offset, c.offset)

	for len(items) > 0 {
		if node.tokenLength() > 0 {
			excluded = nil
		}

		items[0].parse(t, c, excluded)
		items = items[1:]

		if !c.match {
			t.Out1("fail, item failed")
			c.cache.set(node.from, p.name, nil)
			c.fail(node.from)
			return
		}

		node.append(c.node)
	}

	t.Out1("success, items parsed")
	c.cache.set(node.from, p.name, node)
	c.success(node)
}
