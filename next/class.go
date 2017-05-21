package next

type classDefinition struct {
	name     string
	anything bool
	not      bool
	chars    []rune
	ranges   [][]rune
	registry *registry
}

type classGenerator struct {
	name     string
	isValid  bool
	anything bool
	not      bool
	chars    []rune
	ranges   [][]rune
}

type classParser struct {
	name     string
	trace    Trace
	anything bool
	not      bool
	chars    []rune
	ranges   [][]rune
}

func newAnyCharDefinition(r *registry, name string) *classDefinition {
	return &classDefinition{
		name:     name,
		anything: true,
		registry: r,
	}
}

func newClassDefinition(r *registry, name string, not bool, chars []rune, ranges [][]rune) *classDefinition {
	return &classDefinition{
		name:     name,
		not:      not,
		chars:    chars,
		ranges:   ranges,
		registry: r,
	}
}

func (d *classDefinition) nodeName() string { return d.name }

func (d *classDefinition) member(n string, excluded []string) (bool, error) {
	return !stringsContain(excluded, d.name) && n == d.name, nil
}

func (d *classDefinition) generator(_ Trace, init string, excluded []string) (generator, error) {
	if g, ok := d.registry.generator(d.name, init, excluded); ok {
		return g, nil
	}

	g := &classGenerator{
		name:     d.name,
		isValid:  !stringsContain(excluded, d.name) && init == "",
		anything: d.anything,
		not:      d.not,
		chars:    d.chars,
		ranges:   d.ranges,
	}

	d.registry.setGenerator(d.name, init, excluded, g)
	return g, nil
}

func (g *classGenerator) nodeName() string                  { return g.name }
func (g *classGenerator) valid() bool                       { return g.isValid }
func (g *classGenerator) validate(Trace, []generator) error { return nil }

func (g *classGenerator) parser(t Trace, _ *Node) parser {
	return &classParser{
		name:     g.name,
		trace:    t.Extend(g.name),
		anything: g.anything,
		not:      g.not,
		chars:    g.chars,
		ranges:   g.ranges,
	}
}

func (p *classParser) nodeName() string { return p.name }

func (p *classParser) match(t rune) bool {
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

func (p *classParser) parse(c *context) {
	p.trace.Info("parsing", c.offset)

	if c.fillFromCache(p.name, nil) {
		p.trace.Info("found in cache")
		return
	}

	if t, ok := c.token(); ok && (p.anything || p.match(t)) {
		p.trace.Info("success", c.offset, string([]rune{t}))
		c.success(newNode(p.name, Alias, c.offset, c.offset+1))
	} else {
		p.trace.Info("fail", c.offset)
		c.fail(p.name, c.offset, nil)
	}
}
