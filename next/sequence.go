package next

import "fmt"

type sequenceDefinition struct {
	name     string
	items    []string
	registry *registry
	commit   CommitType
}

type sequenceGenerator struct {
	name      string
	id        int
	isVoid    bool
	isValid   bool
	first     generator
	restInit  []generator
	rest      []generator
	firstName string
	restNames []string
	commit    CommitType
	initName  string // TODO: the same problem in the quantifier
}

type sequenceParser struct {
	name      string
	genID     int
	trace     Trace
	init      *Node
	first     generator
	restInit  []generator
	rest      []generator
	firstName string
	restNames []string
	node      *Node
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

	g := &sequenceGenerator{
		name:     d.name,
		id:       id,
		commit:   d.commit,
		initName: init,
	}

	d.registry.setGenerator(id, g)

	firstName := d.items[0]
	first, ok, err := items[0].generator(t, init, append(excluded, d.name))
	if err != nil {
		return nil, false, err
	} else if !ok { // TODO: should check if init can be used
		first = nil
	}

	items = items[1:]

	restNames := make([]string, len(items))
	restInit := make([]generator, len(items))
	rest := make([]generator, len(items))
	for i, item := range items {
		restNames[i] = item.nodeName()

		g, ok, err := item.generator(t, init, nil)
		if err != nil {
			return nil, false, err
		}

		if ok {
			restInit[i] = g
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

	g.first = first
	g.restInit = restInit
	g.rest = rest
	g.firstName = firstName
	g.restNames = restNames
	g.initName = init

	return g, true, nil
}

func (g *sequenceGenerator) nodeName() string { return g.name }
func (g *sequenceGenerator) void() bool       { return g.isVoid }

func (g *sequenceGenerator) finalize(t Trace) {
	t = t.Extend(g.name)

	canUseInit := g.initName == g.firstName || stringsContain(g.restNames, g.initName)
	if g.first == nil && !canUseInit || g.first != nil && g.first.void() && !canUseInit {
		g.isVoid = true
		return
	}

	var restInitVoid bool
	for i := range g.rest {
		if g.restInit[i] == nil {
			restInitVoid = true
		} else if g.restInit[i].void() {
			g.restInit[i] = nil
			restInitVoid = true
		}

		if g.rest[i] != nil && g.rest[i].void() {
			g.rest[i] = nil
		}

		if restInitVoid && g.rest[i] == nil && !canUseInit {
			g.isVoid = true
			return
		}
	}
}

func (g *sequenceGenerator) parser(t Trace, init *Node) parser {
	return &sequenceParser{
		name:      g.name,
		genID:     g.id,
		trace:     t.Extend(g.name),
		node:      newNode(g.name, g.commit, 0, 0),
		init:      init,
		first:     g.first,
		restInit:  g.restInit,
		rest:      g.rest,
		firstName: g.firstName,
		restNames: g.restNames,
	}
}

func (p *sequenceParser) nodeName() string { return p.name }

func (p *sequenceParser) nextParser() parser {
	if len(p.node.Nodes) == 0 {
		if p.first == nil {
			return nil
		}

		return p.first.parser(p.trace, p.init)
	}

	var (
		rest generator
		init *Node
	)

	if p.node.len() == 0 {
		rest = p.restInit[0]
		init = p.init
	} else {
		rest = p.rest[0]
	}

	p.restInit, p.rest, p.restNames = p.restInit[1:], p.rest[1:], p.restNames[1:]

	if rest == nil {
		return nil
	}

	return rest.parser(p.trace, init)
}

func (p *sequenceParser) parse(c *context) {
	if c.fillFromCache(p.genID, p.init) {
		p.trace.Info("found in cache", c.match)
		return
	}

	c.initNode(p.node, p.init)
	for {
		p.trace.Info("parsing sequence", c.offset)

		if len(p.node.Nodes) > 0 && len(p.rest) == 0 {
			p.trace.Info("success")
			c.success(p.genID, p.node)
			return
		}

		itemParser := p.nextParser()
		if itemParser != nil {
			itemParser.parse(c)
			if c.match && c.node.len() > 0 {
				p.node.appendNode(c.node)
				continue
			}
		}

		if p.init != nil && p.node.len() == 0 &&
			(len(p.node.Nodes) == 0 && p.init.Name == p.firstName ||
				len(p.node.Nodes) > 0 && p.init.Name == p.restNames[0]) {

			p.node.appendNode(p.init)
			continue
		}

		if itemParser != nil && c.match {
			p.node.appendNode(c.node)
			continue
		}

		p.trace.Info("fail, no match")
		c.fail(p.genID, p.node.from)
		return
	}
}
