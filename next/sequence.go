package next

type sequenceDefinition struct {
	name   string
	commit CommitType
	items  []string
}

type sequenceParser struct {
	name      string
	commit    CommitType
	items     []parser
	including []parser
}

func newSequence(name string, ct CommitType, items []string) *sequenceDefinition {
	return &sequenceDefinition{
		name:   name,
		commit: ct,
		items:  items,
	}
}

func (d *sequenceDefinition) nodeName() string { return d.name }

func (d *sequenceDefinition) parser(r *registry, path []string) (parser, error) {
	if stringsContain(path, d.name) {
		panic(errCannotIncludeParsers)
	}

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

		item, err := itemDefinition.parser(r, path)
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

func (p *sequenceParser) setIncludedBy(i parser, path []string) {
	if stringsContain(path, p.name) {
		panic(errCannotIncludeParsers)
	}

	p.including = append(p.including, i)
}

func (p *sequenceParser) cacheIncluded(*context, *Node) {
	panic(errCannotIncludeParsers)
}

/*
should be possible to parse:

a = "0"
b = "1"
c = a* e b
d = a | c
e = b | d

input: 111
*/

func (p *sequenceParser) parse(t Trace, c *context) {
	t = t.Extend(p.name)
	t.Out1("parsing sequence", c.offset)

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

	items := p.items
	node := newNode(p.name, p.commit, c.offset, c.offset)

	for len(items) > 0 {
		t.Out2("next sequence item")
		// n, m, ok := c.cache.get(c.offset, items[0].nodeName())
		m, ok := c.fromCache(items[0].nodeName())
		if ok {
			t.Out1("sequence item found in cache, match:", m)
			if m {
				t.Out2("sequence item from cache:", c.node.Name, len(c.node.Nodes), c.node.from)
				node.append(c.node)
				items = items[1:]
				continue
			}

			c.cache.set(node.from, p.name, nil)
			c.fail(node.from)
			return
		}

		items[0].parse(t, c)
		items = items[1:]

		if !c.match {
			t.Out1("fail, item failed")
			c.cache.set(node.from, p.name, nil)
			c.fail(node.from)
			return
		}

		t.Out2("appending sequence item", c.node.Name, len(c.node.Nodes))
		node.append(c.node)
	}

	t.Out1("success, items parsed")
	t.Out2("nodes", node.nodeLength())
	if node.Name == "group" {
		t.Out2("caching group", node.from, node.Nodes[2].Name, node.Nodes[2].nodeLength())
	}

	c.cache.set(node.from, p.name, node)
	for _, i := range p.including {
		i.cacheIncluded(c, node)
	}

	t.Out2("caching sequence and included by done")
	c.success(node)
}
