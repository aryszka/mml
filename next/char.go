package next

type charDefinition struct {
	name     string
	value    rune
	registry *registry
}

// TODO: merge gen and parser

type charGenerator struct {
	name    string
	value   rune
	isValid bool
}

type charParser struct {
	name  string
	trace Trace
	value rune
}

func newCharDefinition(r *registry, name string, value rune) *charDefinition {
	return &charDefinition{
		name:     name,
		value:    value,
		registry: r,
	}
}

func (d *charDefinition) nodeName() string {
	return d.name
}

func (d *charDefinition) member(name string) (bool, error) {
	return name == d.name, nil
}

func (d *charDefinition) generator(_ Trace, init string, excluded []string) (generator, error) {
	if g, ok := d.registry.generator(d.name, init, excluded); ok {
		return g, nil
	}

	g := &charGenerator{
		name:    d.name,
		isValid: !stringsContain(excluded, d.name) && init == "",
		value:   d.value,
	}

	d.registry.setGenerator(d.name, init, excluded, g)
	return g, nil
}

func (g *charGenerator) nodeName() string               { return g.name }
func (g *charGenerator) valid() bool                    { return g.isValid }
func (g *charGenerator) validate(Trace, []string) error { return nil }

func (g *charGenerator) parser(t Trace, _ *Node) parser {
	return &charParser{
		name:  g.name,
		trace: t.Extend(g.name),
		value: g.value,
	}
}

func (p *charParser) nodeName() string { return p.name }

func (p *charParser) parse(c *context) {
	if c.fillFromCache(p.name, nil) {
		return
	}

	if t, ok := c.token(); ok && t == p.value {
		c.succeed(newNode(p.name, Alias, c.offset, c.offset+1))
		c.offset += 1
	} else {
		c.fail(p.name)
	}
}
