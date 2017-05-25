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
	name    string
	id      int
	item    string
	min     int
	max     int
	commit  CommitType
	initial generator
	rest    generator
}

type quantifierParser struct {
	name         string
	genID        int
	trace        Trace
	min          int
	max          int
	initial      generator
	rest         generator
	node         *Node
	init         *Node
	initConsumed bool
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

	initial, ok, err := item.generator(t, init, append(excluded, d.name))
	if !ok || err != nil {
		return nil, false, err
	}

	rest, ok, err := item.generator(t, "", []string{d.name})
	if !ok || err != nil {
		return nil, false, err
	}

	g := &quantifierGenerator{
		name:    d.name,
		id:      id,
		min:     d.min,
		max:     d.max,
		commit:  d.commit,
		initial: initial,
		rest:    rest,
	}

	d.registry.setGenerator(id, g)
	return g, true, nil
}

func (g *quantifierGenerator) nodeName() string { return g.name }

func (g *quantifierGenerator) parser(t Trace, init *Node) parser {
	return &quantifierParser{
		name:    g.name,
		genID:   g.id,
		trace:   t.Extend(g.name),
		min:     g.min,
		max:     g.max,
		node:    newNode(g.name, g.commit, 0, 0),
		init:    init,
		initial: g.initial,
		rest:    g.rest,
	}
}

func (p *quantifierParser) nodeName() string { return p.name }

func (p *quantifierParser) nextParser() (parser parser, member bool) {
	useInitial := len(p.node.Nodes) == 0 || (p.init != nil && !p.initConsumed)
	switch {
	case useInitial:
		parser = p.initial.parser(p.trace, p.init)

	case useInitial && p.init.Name == p.name:
		// this should happen if the initial parser failed
		member = true

	default:
		parser = p.rest.parser(p.trace, p.init)
	}

	return
}

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
		if p.init == nil || p.initConsumed {
			itemParser = p.rest.parser(p.trace, p.init)
		} else {
			itemParser = p.initial.parser(p.trace, nil)
		}

		itemParser.parse(c)
		if c.match {
			p.node.appendNode(c.node)
			if len(c.node.Nodes) > 0 {
				p.initConsumed = true
			}

			continue
		}

		if p.init != nil && !p.initConsumed {
			p.node.appendNode(p.init)
			p.initConsumed = true
			continue
		}

		if len(p.node.Nodes) < p.min {
			p.trace.Info("fail")
			c.fail(p.genID, c.node.from)
			return
		}

		p.trace.Info("success, item invalid")
		c.success(p.genID, p.node)
		return
	}
}
