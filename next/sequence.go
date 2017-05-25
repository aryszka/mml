package next

import "fmt"

type sequenceDefinition struct {
	name     string
	items    []string
	registry *registry
	commit   CommitType
}

type sequenceGenerator struct {
	name         string
	id           int
	isValid      bool
	initial      []generator
	rest         []generator
	initIsMember []bool
	commit       CommitType
}

type sequenceParser struct {
	name         string
	genID        int
	trace        Trace
	init         *Node
	initial      []generator
	rest         []generator
	initIsMember []bool
	node         *Node
	initConsumed bool
}

func sequenceWithoutItems(name string) error {
	return fmt.Errorf("sequence without items: %s", name)
}

func invalidSequenceItem(name, itemName string) error {
	return fmt.Errorf("invalid sequence item %s/%s", name, itemName)
}

func newSequence(r *registry, name string, ct CommitType, items []string) *sequenceDefinition {
	return &sequenceDefinition{
		name:     name,
		items:    items,
		registry: r,
		commit:   ct,
	}
}

func (d *sequenceDefinition) nodeName() string { return d.name }

func (d *sequenceDefinition) generator(t Trace, init string, excluded []string) (generator, bool, error) {
	t = t.Extend(d.name)

	if stringsContain(excluded, d.name) {
		return nil, false, nil
	}

	id := d.registry.genID(d.name, init, excluded)
	if g, ok := d.registry.generator(id); ok {
		return g, true, nil
	}

	// TODO: standardize where these checks happen
	if len(d.items) == 0 {
		return nil, false, sequenceWithoutItems(d.name)
	}

	items, err := d.registry.findDefinitions(d.items)
	if err != nil {
		return nil, false, err
	}

	initial := make([]generator, len(items))
	rest := make([]generator, len(items))
	excluded = append(excluded, d.name)
	for i, item := range items {
		g, ok, err := item.generator(t, init, excluded)
		if !ok && i == 0 || err != nil {
			return nil, false, err
		}

		if ok {
			initial[i] = g
		}

		if i == 0 {
			continue
		}

		g, ok, err = item.generator(t, "", nil)
		if err != nil {
			return nil, false, err
		}

		if !ok {
			return nil, false, invalidSequenceItem(d.name, item.nodeName())
		}

		rest[i] = g
	}

	g := &sequenceGenerator{
		name:    d.name,
		id:      id,
		commit:  d.commit,
		initial: initial,
		rest:    rest,
	}

	d.registry.setGenerator(id, g)
	return g, true, nil
}

func (g *sequenceGenerator) nodeName() string { return g.name }

func (g *sequenceGenerator) parser(t Trace, init *Node) parser {
	return &sequenceParser{
		name:    g.name,
		genID:   g.id,
		trace:   t.Extend(g.name),
		node:    newNode(g.name, g.commit, 0, 0),
		init:    init,
		initial: g.initial,
		rest:    g.rest,
	}
}

func (p *sequenceParser) nodeName() string { return p.name }

func (p *sequenceParser) nextParser() (parser, bool) {
	var gen generator
	if p.initConsumed {
		gen = p.rest[0]
	} else {
		gen = p.initial[0]
	}

	p.initial, p.rest = p.initial[1:], p.rest[1:]
	if gen == nil {
		return nil, false
	}

	return gen.parser(p.trace, nil), true
}

func (p *sequenceParser) parse(c *context) {
	if c.fillFromCache(p.genID, p.init) {
		p.trace.Info("found in cache", c.match)
		return
	}

	c.initNode(p.node, p.init)
	for {
		p.trace.Info("parsing", c.offset)

		if len(p.initial) == 0 {
			p.trace.Info("success")
			c.success(p.genID, p.node)
			return
		}

		itemParser, ok := p.nextParser()
		if !ok {
			p.trace.Info("fail, no parser")
			c.fail(p.genID, p.node.from)
			return
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

		p.trace.Info("fail, no match")
		c.fail(p.genID, p.node.from)
		return
	}
}
