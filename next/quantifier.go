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

func (d *quantifierDefinition) parser(r *registry) (parser, error) {
	p, ok := r.parser(d.name)
	if ok {
		return p, nil
	}

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

	p = &quantifierParser{
		name:   d.name,
		commit: d.commit,
		min:    d.min,
		max:    d.max,
		item:   item,
	}

	r.setParser(p)
	return p, nil
}

func (p *quantifierParser) parse(t Trace, c *context, excluded []string) {
	t = t.Extend(p.name)
	t.Println0("parsing quantifier", c.offset)

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
	node := newNode(p.name, p.commit, c.offset, c.offset)

	for {
		if p.max >= 0 && node.nodeLength() == p.max {
			t.Println0("success, max reached")
			c.success(node)
			return
		}

		if node.tokenLength() > 0 {
			excluded = nil
		}

		p.item.parse(t, c, excluded)
		if !c.match {
			if node.nodeLength() >= p.min {
				t.Println0("success, no more match")
				c.success(node)
			} else {
				t.Println0("fail, min not reached")
				c.fail(p.name, node.from)
			}

			return
		}

		node.append(c.node)
	}
}
