package next

type optionalDefinition struct {
	name     string
	registry *registry
	commit   CommitType
	optional string
}

type optionalGenerator struct {
	name         string
	commit       CommitType
	isValid      bool
	optional     generator
	initIsMember bool
}

type optionalParser struct {
	name         string
	trace        Trace
	initIsMember bool
	node         *Node
	init         *Node
	optional     generator
}

func newOptional(r *registry, name string, ct CommitType, optional string) *optionalDefinition {
	return &optionalDefinition{
		name:     name,
		registry: r,
		commit:   ct,
		optional: optional,
	}
}

func (d *optionalDefinition) nodeName() string { return d.name }

func (d *optionalDefinition) generator(t Trace, init string, excluded []string) (generator, error) {
	t = t.Extend(d.name)

	if g, ok := d.registry.generator(d.name, init, excluded); ok {
		return g, nil
	}

	optional, err := d.registry.findDefinition(d.optional)
	if err != nil {
		return nil, err
	}

	g := &optionalGenerator{
		name:    d.name,
		isValid: true,
		commit:  d.commit,
	}

	d.registry.setGenerator(d.name, init, excluded, g)

	if stringsContain(excluded, d.name) {
		g.isValid = false
		return g, nil
	}

	var initIsMember bool
	if init != "" {
		if m, err := optional.member(init, nil); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	optGenerator, err := optional.generator(t, init, excluded)
	if err != nil {
		return nil, err
	}

	g.optional = optGenerator
	g.initIsMember = initIsMember
	return g, nil
}

func (d *optionalDefinition) member(n string, excluded []string) (bool, error) {
	if stringsContain(excluded, d.optional) {
		return false, nil
	}

	optional, err := d.registry.findDefinition(d.optional)
	if err != nil {
		return false, err
	}

	return optional.member(n, append(excluded, d.name))
}

func (g *optionalGenerator) nodeName() string { return g.name }
func (g *optionalGenerator) valid() bool      { return g.isValid }

func (g *optionalGenerator) validate(t Trace, excluded []generator) error {
	t = t.Extend(g.name)

	if !g.isValid {
		return nil
	}

	if generatorsContain(excluded, g) {
		return nil
	}

	if err := g.optional.validate(t, append(excluded, g)); err != nil {
		return err
	}

	if g.optional != nil && !g.optional.valid() {
		g.isValid = false
	}

	return nil
}

func (g *optionalGenerator) parser(t Trace, init *Node) parser {
	return &optionalParser{
		name:         g.name,
		trace:        t.Extend(g.name),
		initIsMember: g.initIsMember,
		init:         init,
		optional:     g.optional,
		node:         newNode(g.name, g.commit, 0, 0),
	}
}

func (p *optionalParser) nodeName() string { return p.name }

func (p *optionalParser) parse(c *context) {
	p.trace.Info("parsing", c.offset)

	if c.fillFromCache(p.name, p.init) {
		return
	}

	c.initRange(p.node, p.init)
	if p.optional == nil || !p.optional.valid() {
		c.success(p.node)
		return
	}

	optional := p.optional.parser(p.trace, p.init)
	optional.parse(c)
	if c.valid {
		p.node.appendNode(c.node)
	}

	c.success(p.node)
}
