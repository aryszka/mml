package next

type choiceDefinition struct {
	name     string
	registry *registry
	commit   CommitType
	elements []string
}

type choiceGenerator struct {
	name         string
	id           int
	isValid      bool
	commit       CommitType
	generators   [][]generator
	initIsMember bool
}

type choiceParser struct {
	name         string
	genID        int
	trace        Trace
	commit       CommitType
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

func (d *choiceDefinition) generator(t Trace, init string, excluded []string) (generator, bool, error) {
	t = t.Extend(d.name)

	if stringsContain(excluded, d.name) {
		return nil, false, nil
	}

	id := d.registry.genID(d.name, init, excluded)
	if g, ok := d.registry.generator(id); ok {
		return g, true, nil
	}

	elements, err := d.registry.findDefinitions(d.elements)
	if err != nil {
		return nil, false, err
	}

	// TODO: seems like getting a parser when it shouldn't
	// - see terminal tests

	excluded = append(excluded, d.name)
	generators := make([][]generator, len(d.elements)+2)
	for i, it := range append([]string{init}, append(d.elements, d.name)...) {
		g := make([]generator, len(elements))
		for j, e := range elements {
			ge, ok, err := e.generator(t, it, excluded)
			if err != nil {
				return nil, false, err
			}

			if ok {
				g[j] = ge
			}
		}

		generators[i] = g
	}

	if generators[0][0] == nil {
		return nil, false, nil
	}

	g := &choiceGenerator{
		name:       d.name,
		id:         id,
		isValid:    true,
		commit:     d.commit,
		generators: generators,
	}

	d.registry.setGenerator(id, g)
	return g, true, nil
}

func (g *choiceGenerator) nodeName() string { return g.name }

func (g *choiceGenerator) parser(t Trace, init *Node) parser {
	return &choiceParser{
		name:       g.name,
		genID:      g.id,
		trace:      t.Extend(g.name),
		commit:     g.commit,
		init:       init,
		node:       newNode(g.name, g.commit, 0, 0),
		generators: g.generators,
	}
}

func (p *choiceParser) nodeName() string { return p.name }

func (p *choiceParser) parse(c *context) {
	if c.fillFromCache(p.genID, p.init) {
		p.trace.Info("found in cache", c.match)
		return
	}

	c.initNode(p.node, p.init)
	matchIndex := -1
	var (
		initIndex    int
		elementIndex int
		match        bool
	)

	for {
		p.trace.Info("parsing", c.offset)

		if elementIndex == len(p.generators[initIndex]) {
			if matchIndex >= 0 {
				initIndex, elementIndex, matchIndex = matchIndex+1, 0, -1
			} else if match && initIndex < len(p.generators)-1 {
				initIndex, elementIndex, matchIndex = len(p.generators)-1, 0, -1
			} else if match {
				p.trace.Info("success")
				c.success(p.genID, p.node)
				return
			} else {
				p.trace.Info("fail")
				c.fail(p.genID, p.node.from)
				return
			}
		}

		if p.generators[initIndex][elementIndex] == nil {
			elementIndex++
			continue
		}

		init := p.init
		if initIndex > 0 && initIndex == len(p.generators)-1 {
			init = p.node
		} else if initIndex > 0 {
			init = p.node.Nodes[0]
		}

		if init == nil {
			c.offset = p.node.from
		} else {
			c.offset = init.to
		}

		p.generators[initIndex][elementIndex].parser(p.trace, init).parse(c)

		if c.match && (len(p.node.Nodes) == 0 || c.node.len() > p.node.len()) {
			match = true
			matchIndex = elementIndex

			if initIndex == len(p.generators)-1 {
				p.node = newNode(p.name, p.commit, 0, 0)
			} else {
				p.node.clear()
			}

			p.node.appendNode(c.node)
		}

		if !c.match && init != nil && init.Name == p.name && len(p.node.Nodes) == 0 {
			match = true
			matchIndex = elementIndex
			p.node.appendNode(c.node)
		}

		elementIndex++
		continue
	}
}
