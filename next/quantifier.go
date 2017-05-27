package next

type quantifierDefinition struct {
	name     string
	item     string
	min      int
	max      int
	registry *registry
	commit   CommitType
}

type quantifierGenerator struct {
	name     string
	id       int
	item     string
	min      int
	max      int
	commit   CommitType
	first    generator
	restInit generator
	itemName string
	rest     generator
}

type quantifierParser struct {
	name     string
	genID    int
	trace    Trace
	min      int
	max      int
	first    generator
	restInit generator
	rest     generator
	node     *Node
	itemName string
	init     *Node
}

func newQuantifier(r *registry, name string, ct CommitType, item string, min, max int) *quantifierDefinition {
	return &quantifierDefinition{
		name:     name,
		item:     item,
		min:      min,
		max:      max,
		registry: r,
		commit:   ct,
	}
}

func (d *quantifierDefinition) nodeName() string { return d.name }

func (d *quantifierDefinition) generator(t Trace, init string, excluded []string) (generator, bool, error) {
	t = t.Extend(d.name)

	if stringsContain(excluded, d.name) {
		return nil, false, nil
	}

	id := d.registry.genID(d.name, init, excluded)
	if g, ok := d.registry.generator(id); ok {
		return g, true, nil
	}

	item, err := d.registry.findDefinition(d.item)
	if err != nil {
		return nil, false, err
	}

	excluded = append(excluded, d.name)
	first, ok, err := item.generator(t, init, excluded)
	if !ok || err != nil {
		return nil, false, err
	}

	excluded = []string{d.name}

	restInit, ok, err := item.generator(t, init, excluded)
	if err != nil {
		return nil, false, err
	} else if !ok {
		restInit = nil
	}

	rest, ok, err := item.generator(t, "", excluded)
	if err != nil {
		return nil, false, err
	} else if !ok {
		rest = nil
	}

	g := &quantifierGenerator{
		name:     d.name,
		id:       id,
		min:      d.min,
		max:      d.max,
		commit:   d.commit,
		first:    first,
		restInit: restInit,
		rest:     rest,
		itemName: item.nodeName(),
	}

	d.registry.setGenerator(id, g)
	return g, true, nil
}

func (g *quantifierGenerator) nodeName() string { return g.name }

func (g *quantifierGenerator) parser(t Trace, init *Node) parser {
	return &quantifierParser{
		name:     g.name,
		genID:    g.id,
		trace:    t.Extend(g.name),
		min:      g.min,
		max:      g.max,
		node:     newNode(g.name, g.commit, 0, 0),
		init:     init,
		first:    g.first,
		restInit: g.restInit,
		rest:     g.rest,
		itemName: g.itemName,
	}
}

func (p *quantifierParser) nodeName() string { return p.name }

func (p *quantifierParser) parse(c *context) {
	if c.fillFromCache(p.genID, p.init) {
		p.trace.Info("found in cache", c.match)
		return
	}

	c.initNode(p.node, p.init)
	for {
		p.trace.Info("parsing", c.offset)

		if p.max >= 0 && len(p.node.Nodes) == p.max {
			p.trace.Info("success, max")
			c.success(p.genID, p.node)
			return
		}

		var itemParser parser
		if len(p.node.Nodes) == 0 && p.first != nil {
			itemParser = p.first.parser(p.trace, p.init)
		} else if p.node.len() == 0 && p.restInit != nil {
			itemParser = p.restInit.parser(p.trace, p.init)
		} else if p.rest != nil {
			itemParser = p.rest.parser(p.trace, nil)
		}

		if itemParser != nil {
			itemParser.parse(c)
			if c.match && c.node.len() > 0 {
				p.node.appendNode(c.node)
				continue
			}
		}

		if p.init != nil && p.node.len() == 0 && p.init.Name == p.itemName {
			p.node.appendNode(p.init)
			continue
		}

		if itemParser != nil && c.match && len(p.node.Nodes) < p.min {
			p.node.appendNode(c.node)
			continue
		}

		if len(p.node.Nodes) < p.min {
			p.trace.Info("fail, short")
			c.fail(p.genID, c.node.from)
			return
		}

		p.trace.Info("success, next item invalid")
		c.success(p.genID, p.node)
		return
	}
}
