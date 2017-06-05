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

func newSequenceDefinition(name string, ct CommitType, items []string) *sequenceDefinition {
	return &sequenceDefinition{
		name:   name,
		commit: ct,
		items:  items,
	}
}

func (d *sequenceDefinition) parser(r *registry) (parser, error) {
	p, ok := r.parser(d.name)
	if ok {
		return p, nil
	}

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

	p = &sequenceParser{
		name:   d.name,
		commit: d.commit,
		items:  items,
	}

	r.setParser(p)
	return p, nil
}

func (p *sequenceParser) parse(t Trace, c *context, excluded []string) {
	t = t.Extend(p.name)
	t.Println0("parsing sequence", c.offset)

	if stringsContain(excluded, p.name) {
		t.Println0("excluded")
		c.fail(p.name, c.offset)
		return
	}

	if p.commit&Documentation != 0 {
		t.Println0("fail, doc")
		c.fail(p.name, c.offset)
		return
	}

	if c.checkCache(p.name) {
		t.Println0("found in cache")
		return
	}

	excluded = append(excluded, p.name)
	items := p.items
	node := newNode(p.name, p.commit, c.offset, c.offset)

	for len(items) > 0 {
		if node.tokenLength() > 0 {
			excluded = nil
		}

		items[0].parse(t, c, excluded)
		items = items[1:]

		if !c.match {
			t.Println0("fail, item failed")
			c.fail(p.name, node.from)
			return
		}

		node.append(c.node)
	}

	t.Println0("success, items parsed")
	c.success(node)
}
