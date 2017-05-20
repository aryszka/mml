package next

type choiceDefinition struct {
	name     string
	registry *registry
	commit   CommitType
	elements []string
}

type choiceGenerator struct {
	name         string
	isValid      bool
	commit       CommitType
	generators   [][]generator
	initIsMember bool
}

type choiceParser struct {
	name         string
	trace        Trace
	init         *Node
	node         *Node
	generators   [][]generator
	initIsMember bool
	initIndex    int
	parserIndex  int
	valid        bool
}

func newChoice(r *registry, name string, ct CommitType, elements []string) *choiceDefinition {
	return &choiceDefinition{
		name:     name,
		registry: r,
		commit:   ct,
		elements: elements,
	}
}

func (d *choiceDefinition) nodeName() string { return d.name }

func (d *choiceDefinition) member(n string, excluded []string) (bool, error) {
	if stringsContain(excluded, d.name) {
		return false, nil
	}

	if n == d.name {
		return true, nil
	}

	defs, err := d.registry.findDefinitions(d.elements)
	if err != nil {
		return false, err
	}

	excluded = append(excluded, d.name)
	for _, di := range defs {
		if m, err := di.member(n, excluded); m || err != nil {
			return m, err
		}
	}

	return false, nil
}

func (d *choiceDefinition) generator(t Trace, init string, excluded []string) (generator, error) {
	t = t.Extend(d.name)

	if g, ok := d.registry.generator(d.name, init, excluded); ok {
		return g, nil
	}

	g := &choiceGenerator{
		name:    d.name,
		isValid: true,
		commit:  d.commit,
	}

	d.registry.setGenerator(d.name, init, excluded, g)

	elements, err := d.registry.findDefinitions(d.elements)
	if err != nil {
		return nil, err
	}

	generators := make([][]generator, len(elements)+1)
	for i, it := range append([]string{init}, d.elements...) {
		g := make([]generator, len(elements))
		for j, e := range elements {
			ge, err := e.generator(t, it, excluded)
			if err != nil {
				return nil, err
			}

			g[j] = ge
		}

		generators[i] = g
	}

	var initIsMember bool
	if init != "" {
		if m, err := d.member(init, nil); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	if !initIsMember && (len(generators[0]) == 0 || !generators[0][0].valid()) {
		g.isValid = false
		return g, nil
	}

	g.generators = generators
	g.initIsMember = initIsMember
	return g, nil
}

func (g *choiceGenerator) nodeName() string { return g.name }
func (g *choiceGenerator) valid() bool      { return g.isValid }

func (g *choiceGenerator) validate(t Trace, excluded []generator) error {
	t = t.Extend(g.name)

	if !g.isValid {
		return nil
	}

	if generatorsContain(excluded, g) {
		return nil
	}

	excluded = append(excluded, g)
	for i := 0; i < len(g.generators); i++ {
		for j := 0; j < len(g.generators[i]); j++ {
			if err := g.generators[i][j].validate(t, excluded); err != nil {
				return err
			}
		}
	}

	if !g.initIsMember && !g.generators[0][0].valid() {
		g.isValid = false
	}

	return nil
}

func (g *choiceGenerator) parser(t Trace, init *Node) parser {
	return &choiceParser{
		name:         g.name,
		trace:        t.Extend(g.name),
		init:         init,
		node:         newNode(g.name, g.commit, 0, 0),
		generators:   g.generators,
		initIsMember: g.initIsMember,
	}
}

func (p *choiceParser) nodeName() string { return p.name }

func (p *choiceParser) stepInit() {
	p.initIndex = p.parserIndex + 1
	p.parserIndex = 0
}

func (p *choiceParser) stepParser() {
	p.parserIndex++
}

func (p *choiceParser) nextParser() (parser, bool, bool) {
	if p.initIndex == len(p.generators) {
		return nil, false, false
	}

	if p.initIndex == 0 && p.initIsMember {
		return nil, true, true
	}

	for {
		if p.parserIndex == len(p.generators[p.initIndex]) {
			return nil, false, false
		}

		if p.generators[p.initIndex][p.parserIndex].valid() {
			var init *Node
			if len(p.node.Nodes) > 0 {
				init = p.node.Nodes[0]
			}

			return p.generators[p.initIndex][p.parserIndex].parser(p.trace, init), false, true
		}

		p.stepParser()
	}
}

func (p *choiceParser) appendNode(n *Node) {
	p.node.clear()
	p.node.appendNode(n)
	p.stepInit()
	p.valid = true
}

func (p *choiceParser) parse(c *context) {
	p.trace.Info("parsing", c.offset)

	if c.fillFromCache(p.name, p.init) {
		return
	}

	c.initRange(p.node, p.init)
	for {
		elementParser, member, ok := p.nextParser()
		if !ok {
			if p.valid {
				c.success(p.node)
				return
			}

			c.fail(p.name, p.node.from, p.init)
			return
		}

		if member {
			p.appendNode(p.init)
			continue
		}

		elementParser.parse(c)
		if c.valid && c.node.len() > p.node.len() {
			p.appendNode(c.node)
			continue
		}

		p.stepParser()
	}
}
