package next

type quantifierDefinition struct {
	name     string
	item     string
	min      int
	max      int
	registry *registry
	commit   CommitType
}

type quantifierGenerator struct {
	name     string
	id       int
	isVoid   bool
	item     string
	min      int
	max      int
	commit   CommitType
	first    generator
	restInit generator
	itemName string
	rest     generator
	initName string
}

type quantifierParser struct {
	name     string
	genID    int
	trace    Trace
	min      int
	max      int
	first    generator
	restInit generator
	rest     generator
	node     *Node
	itemName string
	init     *Node
}

func newQuantifier(r *registry, name string, ct CommitType, item string, min, max int) *quantifierDefinition {
	return &quantifierDefinition{
		name:     name,
		item:     item,
		min:      min,
		max:      max,
		registry: r,
		commit:   ct,
	}
}

func (d *quantifierDefinition) nodeName() string { return d.name }

func (d *quantifierDefinition) generator(t Trace, init string, excluded []string) (generator, bool, error) {
	// TODO: maybe all validation from here should be moved to validate

	t = t.Extend(d.name)

	if stringsContain(excluded, d.name) {
		return nil, false, nil
	}

	id := d.registry.genID(d.name, init, excluded)
	if g, ok := d.registry.generator(id); ok {
		return g, true, nil
	}

	item, err := d.registry.findDefinition(d.item)
	if err != nil {
		return nil, false, err
	}

	g := &quantifierGenerator{
		name:     d.name,
		id:       id,
		min:      d.min,
		max:      d.max,
		commit:   d.commit,
		itemName: item.nodeName(),
		initName: init,
	}

	d.registry.setGenerator(id, g)

	excluded = append(excluded, d.name)
	first, ok, err := item.generator(t, init, excluded)
	if err != nil {
		return g, false, err
	} else if !ok {
		first = nil
	}

	excluded = []string{d.name}

	var restInit generator
	if init != "" {
		rig, ok, err := item.generator(t, init, excluded)
		if err != nil {
			return nil, false, err
		} else if ok {
			restInit = rig
		}
	}

	rest, ok, err := item.generator(t, "", excluded)
	if err != nil {
		return nil, false, err
	} else if !ok {
		rest = nil
	}

	g.first = first
	g.restInit = restInit
	g.rest = rest

	return g, true, nil
}

func (g *quantifierGenerator) nodeName() string { return g.name }
func (g *quantifierGenerator) void() bool       { return g.isVoid }

func (g *quantifierGenerator) finalize(t Trace) {
	t.Extend(g.name)

	canUseInit := g.initName == g.itemName

	if g.first != nil && g.first.void() {
		g.first = nil
	}

	if g.restInit != nil && g.restInit.void() {
		g.restInit = nil
	}

	if g.rest != nil && g.rest.void() {
		g.rest = nil
	}

	g.isVoid = g.first == nil || g.min > 1 && (g.restInit == nil && g.rest == nil) && !canUseInit
}

func (g *quantifierGenerator) parser(t Trace, init *Node) parser {
	return &quantifierParser{
		name:     g.name,
		genID:    g.id,
		trace:    t.Extend(g.name),
		min:      g.min,
		max:      g.max,
		node:     newNode(g.name, g.commit, 0, 0),
		init:     init,
		first:    g.first,
		restInit: g.restInit,
		rest:     g.rest,
		itemName: g.itemName,
	}
}

func (p *quantifierParser) nodeName() string { return p.name }

func (p *quantifierParser) parse(c *context) {
	if c.fillFromCache(p.genID, p.init) {
		p.trace.Info("found in cache", c.match)
		return
	}

	c.initNode(p.node, p.init)
	for {
		p.trace.Info("parsing", c.offset)

		if p.max >= 0 && len(p.node.Nodes) == p.max {
			p.trace.Info("success, max")
			c.success(p.genID, p.node)
			return
		}

		var itemParser parser
		if len(p.node.Nodes) == 0 {
			if p.first != nil {
				itemParser = p.first.parser(p.trace, p.init)
			}
		} else if p.restInit != nil && p.node.len() == 0 {
			itemParser = p.restInit.parser(p.trace, p.init)
		} else if p.rest != nil && len(p.node.Nodes) > 0 && (p.node.len() > 0 || p.restInit == nil) {
			itemParser = p.rest.parser(p.trace, nil)
		}

		if itemParser != nil {
			itemParser.parse(c)
			if c.match && c.node.len() > 0 {
				p.node.appendNode(c.node)
				continue
			}
		}

		if p.init != nil && p.node.len() == 0 && p.init.Name == p.itemName {
			p.node.appendNode(p.init)
			continue
		}

		// TODO: why this condition?
		if itemParser != nil && c.match && len(p.node.Nodes) < p.min {
			p.node.appendNode(c.node)
			continue
		}

		if len(p.node.Nodes) < p.min {
			p.trace.Info("fail, short")
			c.fail(p.genID, p.node.from)
			return
		}

		p.trace.Info("success, next item invalid")
		c.success(p.genID, p.node)
		return
	}
}
