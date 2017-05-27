package next

type charDefinition struct {
	name     string
	any      bool
	not      bool
	chars    []rune
	ranges   [][]rune
	registry *registry
}

type charGenerator struct {
	name   string
	id     int
	isVoid bool
	any    bool
	not    bool
	chars  []rune
	ranges [][]rune
}

type charParser struct {
	name   string
	genID  int
	trace  Trace
	any    bool
	not    bool
	chars  []rune
	ranges [][]rune
}

func newChar(r *registry, name string, any, not bool, chars []rune, ranges [][]rune) *charDefinition {
	return &charDefinition{
		name:     name,
		any:      any,
		not:      not,
		chars:    chars,
		ranges:   ranges,
		registry: r,
	}
}

func (d *charDefinition) nodeName() string { return d.name }

func (d *charDefinition) generator(_ Trace, init string, excluded []string) (generator, bool, error) {
	if stringsContain(excluded, d.name) {
		return nil, false, nil
	}

	if init != "" {
		return nil, false, nil
	}

	id := d.registry.genID(d.name, init, excluded)
	if g, ok := d.registry.generator(id); ok {
		return g, true, nil
	}

	g := &charGenerator{
		name:   d.name,
		id:     id,
		any:    d.any,
		not:    d.not,
		chars:  d.chars,
		ranges: d.ranges,
	}

	d.registry.setGenerator(id, g)
	return g, true, nil
}

func (g *charGenerator) nodeName() string { return g.name }
func (g *charGenerator) void() bool       { return false }
func (g *charGenerator) finalize(Trace)   {}

func (g *charGenerator) parser(t Trace, _ *Node) parser {
	return &charParser{
		name:   g.name,
		genID:  g.id,
		trace:  t.Extend(g.name),
		any:    g.any,
		not:    g.not,
		chars:  g.chars,
		ranges: g.ranges,
	}
}

func (p *charParser) nodeName() string { return p.name }

func (p *charParser) match(t rune) bool {
	if p.any {
		return true
	}

	for _, ci := range p.chars {
		if ci == t {
			return !p.not
		}
	}

	for _, ri := range p.ranges {
		if t >= ri[0] && t <= ri[1] {
			return !p.not
		}
	}

	return p.not
}

func (p *charParser) parse(c *context) {
	p.trace.Info("parsing char", c.offset)
	if t, ok := c.token(); ok && p.match(t) {
		p.trace.Info("success", c.offset, string([]rune{t}))
		c.success(p.genID, newNode(p.name, Alias, c.offset, c.offset+1))
		return
	} else {
		p.trace.Info("fail", c.offset, string(t))
		c.fail(p.genID, c.offset)
		return
	}
}
