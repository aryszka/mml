package next

import "fmt"

type sequenceDefinition struct {
	name     string
	items    []string
	registry *registry
}

type sequenceGenerator struct {
	name         string
	isValid      bool
	init         string
	initial      []generator
	rest         []generator
	initIsMember []bool
}

type sequenceParser struct {
	name         string
	trace        Trace
	init         *Node
	initial      []generator
	rest         []generator
	initIsMember []bool
	node         *Node
	itemParser   parser
}

func sequenceWithoutItems(name string) error {
	return fmt.Errorf("sequence without items: %s", name)
}

func invalidSequenceItem(name, itemName string) error {
	return fmt.Errorf("invalid sequence item %s/%s", name, itemName)
}

func newSequence(r *registry, name string, items []string) *sequenceDefinition {
	return &sequenceDefinition{
		name:     name,
		items:    items,
		registry: r,
	}
}

func (d *sequenceDefinition) nodeName() string                 { return d.name }
func (d *sequenceDefinition) member(name string) (bool, error) { return name == d.name, nil }

func (d *sequenceDefinition) generator(t Trace, init string, excluded []string) (generator, error) {
	t = t.Extend(d.name)

	if g, ok := d.registry.generator(d.name, init, excluded); ok {
		return g, nil
	}

	if len(d.items) == 0 {
		return nil, sequenceWithoutItems(d.name)
	}

	items, err := d.registry.findDefinitions(d.items)
	if err != nil {
		return nil, err
	}

	g := &sequenceGenerator{
		name:    d.name,
		isValid: true,
		init:    init,
	}

	d.registry.setGenerator(d.name, init, excluded, g)
	if stringsContain(excluded, d.name) {
		g.isValid = false
		return g, nil
	}

	excluded = append(excluded, d.name)
	g.initial = make([]generator, len(items))
	g.rest = make([]generator, len(items))
	g.initIsMember = make([]bool, len(items))
	for i, item := range items {
		gi, err := item.generator(t, init, excluded)
		if err != nil {
			return nil, err
		}

		g.initial[i] = gi

		gi, err = item.generator(t, "", nil)
		if err != nil {
			return nil, err
		}

		if !gi.valid() {
			return nil, invalidSequenceItem(d.name, item.nodeName())
		}

		m, err := item.member(init)
		if err != nil {
			return nil, err
		}

		g.rest[i] = gi
		g.initIsMember[i] = m
	}

	return g, nil
}

func (g *sequenceGenerator) nodeName() string { return g.name }
func (g *sequenceGenerator) valid() bool      { return g.isValid }

// TODO: for the sake of the generated code, better not to keep the invalid generators

func (g *sequenceGenerator) validate(Trace, []string) error {
	if !g.isValid {
		return nil
	}

	var i int
	for i = 0; i < len(g.initial); i++ {
		if g.initial[i].valid() {
			continue
		}

		if i == 0 {
			g.isValid = false
			return nil
		}

		break
	}

	for j := 1; j < len(g.rest); j++ {
		if g.rest[j].valid() {
			continue
		}

		if j >= i {
			g.isValid = false
			return nil
		}
	}

	return nil
}

func (g *sequenceGenerator) parser(t Trace, init *Node) parser {
	return &sequenceParser{
		name:         g.name,
		trace:        t.Extend(g.name),
		node:         &Node{Name: g.name, From: -1, To: -1},
		init:         init,
		initial:      g.initial,
		rest:         g.rest,
		initIsMember: g.initIsMember,
	}
}

func (p *sequenceParser) nodeName() string { return p.name }

func (p *sequenceParser) useInitial() bool {
	if len(p.node.Nodes) == 0 {
		return true
	}

	if p.init == nil {
		return false
	}

	return !p.node.startsWith(p.init)
}

func (p *sequenceParser) nextParser() (bool, bool) {
	switch {
	case p.useInitial() && p.initial[0].valid():
		p.itemParser = p.initial[0].parser(p.trace, p.init)

	case p.useInitial() && p.initIsMember[0]:
		p.node.appendNode(p.init)
		return true, false

	case !p.useInitial() && p.rest[0].valid():
		p.itemParser = p.rest[0].parser(p.trace, nil)

	default:
		return false, false
	}

	p.initial = p.initial[1:]
	p.rest = p.rest[1:]
	p.initIsMember = p.initIsMember[1:]
	return false, true
}

func (p *sequenceParser) parsed() int {
	parsed := p.node.To - p.node.From
	if p.init != nil {
		ri := p.init.To - p.init.From
		if ri <= parsed {
			parsed -= ri
		}
	}

	return parsed
}

func (p *sequenceParser) parse(c *context) {
	if c.fillFromCache(p.name, p.init) {
		return
	}

	for {
		if len(p.initial) == 0 {
			c.succeed(p.node)
			return
		}

		if next, ok := p.nextParser(); !ok {
			c.moveOffset(-p.parsed())
			c.fail(p.name)
			return
		} else if next {
			continue
		}

		p.itemParser.parse(c)
		if c.valid {
			p.node.appendNode(c.node)
			continue
		}

		c.moveOffset(-p.parsed())
		c.fail(p.name)
		return
	}
}
