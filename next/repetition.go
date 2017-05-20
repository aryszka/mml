package next

type repetitionDefinition struct {
	name     string
	item     string
	registry *registry
	commit   CommitType
}

type repetitionGenerator struct {
	name         string
	isValid      bool
	initial      generator
	rest         generator
	initIsMember bool
	commit       CommitType
}

type repetitionParser struct {
	name         string
	trace        Trace
	node         *Node
	init         *Node
	initial      generator
	rest         generator
	initIsMember bool
	itemParser   parser
}

func newRepetition(r *registry, name string, ct CommitType, item string) *repetitionDefinition {
	return &repetitionDefinition{
		name:     name,
		item:     item,
		registry: r,
		commit:   ct,
	}
}

func (d *repetitionDefinition) nodeName() string { return d.name }

func (d *repetitionDefinition) generator(t Trace, init string, excluded []string) (generator, error) {
	t = t.Extend(d.name)

	if g, ok := d.registry.generator(d.name, init, excluded); ok {
		return g, nil
	}

	item, err := d.registry.findDefinition(d.item)
	if err != nil {
		return nil, err
	}

	g := &repetitionGenerator{
		name:    d.name,
		isValid: true,
		commit:  d.commit,
	}

	d.registry.setGenerator(d.name, init, excluded, g)
	if stringsContain(excluded, d.name) {
		g.isValid = false
		return g, nil
	}

	excluded = append(excluded, d.name)

	initial, err := item.generator(t, init, excluded)
	if err != nil {
		return nil, err
	}

	rest, err := item.generator(t, "", []string{d.name})
	if err != nil {
		return nil, err
	}

	var initIsMember bool
	if init != "" {
		if m, err := item.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	g.initial = initial
	g.rest = rest
	g.initIsMember = initIsMember

	return g, nil
}

func (d *repetitionDefinition) member(n string) (bool, error) {
	return n == d.name, nil
}

func (g *repetitionGenerator) nodeName() string { return g.name }
func (g *repetitionGenerator) valid() bool      { return g.isValid }

func (g *repetitionGenerator) validate(Trace, []string) error {
	if !g.isValid {
		return nil
	}

	if !g.initial.valid() {
		g.isValid = false
	}

	return nil
}

func (g *repetitionGenerator) parser(t Trace, init *Node) parser {
	return &repetitionParser{
		name:         g.name,
		trace:        t.Extend(g.name),
		node:         newNode(g.name, g.commit, 0, 0),
		init:         init,
		initial:      g.initial,
		rest:         g.rest,
		initIsMember: g.initIsMember,
	}
}

func (p *repetitionParser) nodeName() string { return p.name }

func (p *repetitionParser) nextParser() (bool, bool) {
	switch {
	case useInitial(p.node, p.init) && p.initial.valid():
		p.itemParser = p.initial.parser(p.trace, p.init)
		return false, true

	case !useInitial(p.node, p.init) && p.rest.valid():
		p.itemParser = p.rest.parser(p.trace, p.init)
		return false, true

	case useInitial(p.node, p.init) && p.initIsMember:
		return true, true

	default:
		return false, false
	}
}

func (p *repetitionParser) parse(c *context) {
	p.trace.Info("parsing")

	if c.fillFromCache(p.name, p.init) {
		return
	}

	p.node.from = c.offset
	p.node.to = p.node.from

	for {
		if member, ok := p.nextParser(); !ok {
			c.success(p.node)
			return
		} else if member {
			p.node.appendNode(p.init)
			continue
		}

		p.itemParser.parse(c)
		if c.valid {
			p.node.appendNode(c.node)
			continue
		}

		c.success(p.node)
		return
	}
}
