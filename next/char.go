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

func (d *charDefinition) member(n string, excluded []string) (bool, error) {
	return !stringsContain(excluded, d.name) && n == d.name, nil
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

func (g *charGenerator) nodeName() string                  { return g.name }
func (g *charGenerator) valid() bool                       { return g.isValid }
func (g *charGenerator) validate(Trace, []generator) error { return nil }

func (g *charGenerator) parser(t Trace, _ *Node) parser {
	return &charParser{
		name:  g.name,
		trace: t.Extend(g.name),
		value: g.value,
	}
}

func (p *charParser) nodeName() string { return p.name }

func (p *charParser) parse(c *context) {
	p.trace.Info("parsing", c.offset)

	if c.fillFromCache(p.name, nil) {
		p.trace.Info("found in cache")
		return
	}

	if t, ok := c.token(); ok && t == p.value {
		p.trace.Info("success", c.offset, t)
		c.success(newNode(p.name, Alias, c.offset, c.offset+1))
	} else {
		p.trace.Info("fail", c.offset)
		c.fail(p.name, c.offset, nil)
	}
}
