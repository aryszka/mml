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
		if m, err := item.member(init, nil); err != nil {
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

func (d *repetitionDefinition) member(n string, excluded []string) (bool, error) {
	if n == d.item {
		return true, nil
	}

	if stringsContain(excluded, d.item) {
		return false, nil
	}

	item, err := d.registry.findDefinition(d.item)
	if err != nil {
		return false, err
	}

	return item.member(n, append(excluded, d.name))
}

func (g *repetitionGenerator) nodeName() string { return g.name }
func (g *repetitionGenerator) valid() bool      { return g.isValid }

func (g *repetitionGenerator) validate(t Trace, excluded []generator) error {
	t = t.Extend(g.name)

	if !g.isValid {
		return nil
	}

	if generatorsContain(excluded, g) {
		return nil
	}

	excluded = append(excluded, g)

	if err := g.initial.validate(t, excluded); err != nil {
		return err
	}

	if err := g.rest.validate(t, excluded); err != nil {
		return err
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

func (p *repetitionParser) nextParser() (parser, bool, bool) {
	switch {
	case useInitial(p.node, p.init) && p.initial.valid():
		return p.initial.parser(p.trace, p.init), false, true

	case !useInitial(p.node, p.init) && p.rest.valid():
		return p.rest.parser(p.trace, p.init), false, true

	case useInitial(p.node, p.init) && p.initIsMember:
		return nil, true, true

	default:
		return nil, false, false
	}
}

func (p *repetitionParser) parse(c *context) {
	if c.fillFromCache(p.name, p.init) {
		p.trace.Info("found in cache", c.offset)
		return
	}

	c.initRange(p.node, p.init)
	for {
		p.trace.Info("parsing", c.offset)

		itemParser, member, ok := p.nextParser()
		if !ok {
			c.success(p.node)
			return
		} else if member {
			p.node.appendNode(p.init)
			continue
		}

		itemParser.parse(c)
		if c.valid {
			p.node.appendNode(c.node)
			continue
		}

		c.success(p.node)
		return
	}
}
