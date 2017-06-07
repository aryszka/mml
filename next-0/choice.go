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
	isVoid       bool
	commit       CommitType
	generators   [][]generator
	elementNames [][]string
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
	elementNames [][]string
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

	g := &choiceGenerator{
		name:   d.name,
		id:     id,
		commit: d.commit,
	}

	d.registry.setGenerator(id, g)

	excluded = append(excluded, d.name)
	elementNames := make([][]string, len(d.elements)+2)
	generators := make([][]generator, len(d.elements)+2)
	for i, it := range append([]string{init}, append(d.elements, d.name)...) {
		n := make([]string, len(elements))
		g := make([]generator, len(elements))

		for j, e := range elements {
			n[j] = e.nodeName()

			ge, ok, err := e.generator(t, it, excluded)
			if err != nil {
				return nil, false, err
			}

			// if d.name == "primary-expression" {
			// 	t.Debug("element", e.nodeName())
			// }

			if ok {
				g[j] = ge
			}
		}

		elementNames[i] = n
		generators[i] = g
	}

	g.generators = generators
	g.elementNames = elementNames

	return g, true, nil
}

func (g *choiceGenerator) nodeName() string { return g.name }
func (g *choiceGenerator) void() bool       { return g.isVoid }

func (g *choiceGenerator) finalize(t Trace) {
	t = t.Extend(g.name)

	for i := range g.generators {
		var hasOne bool
		for j := range g.generators[i] {
			if g.generators[i][j] != nil {
				if g.generators[i][j].void() {
					if g.generators[i][j].nodeName() == "indexer" {
						t.Debug("removed", i, j)
					}

					g.generators[i][j] = nil
				} else {
					if g.generators[i][j].nodeName() == "indexer" {
						t.Debug("left", i, j)
					}

					hasOne = true
				}
			}
		}

		if i == 0 && !hasOne {
			g.isVoid = true
			return
		}
	}
}

func (g *choiceGenerator) parser(t Trace, init *Node) parser {
	return &choiceParser{
		name:         g.name,
		genID:        g.id,
		trace:        t.Extend(g.name),
		commit:       g.commit,
		init:         init,
		node:         newNode(g.name, g.commit, 0, 0),
		generators:   g.generators,
		elementNames: g.elementNames,
	}
}

func (p *choiceParser) nodeName() string { return p.name }

func (p *choiceParser) appendNode(n *Node, initIndex int) {
	if initIndex == len(p.generators)-1 {
		p.node = newNode(p.name, p.commit, 0, 0)
	} else {
		p.node.clear()
	}

	p.node.appendNode(n)
}

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
		p.trace.Info("parsing", c.offset, initIndex, elementIndex)

		if elementIndex == len(p.generators[initIndex]) {
			if matchIndex >= 0 && initIndex < len(p.generators)-1 {
				initIndex, elementIndex, matchIndex = len(p.generators)-1, 0, -1
			} else if matchIndex >= 0 {
				initIndex, elementIndex, matchIndex = matchIndex+1, 0, -1
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

		var init *Node
		if initIndex == len(p.generators)-1 {
			init = p.node
		} else if initIndex > 0 {
			init = p.node.Nodes[0]
		} else {
			init = p.init
		}

		var elementParser parser
		p.trace.Debug("checking parser", p.generators[initIndex][elementIndex] != nil)
		if p.generators[initIndex][elementIndex] != nil {
			elementParser = p.generators[initIndex][elementIndex].parser(p.trace, init)
		}

		if elementParser != nil {
			if init == nil {
				c.offset = p.node.from
			} else {
				c.offset = init.to
			}

			elementParser.parse(c)
			if c.match && c.node.len() > p.node.len() {
				match = true
				matchIndex = elementIndex
				p.appendNode(c.node, initIndex)

				elementIndex++
				continue
			}
		}

		if init != nil &&
			p.node.len() == 0 &&
			init.Name == p.elementNames[initIndex][elementIndex] {

			match = true
			matchIndex = elementIndex
			p.appendNode(init, initIndex)

			elementIndex++
			continue
		}

		if elementParser != nil && c.match && len(p.node.Nodes) == 0 {
			match = true
			matchIndex = elementIndex
			p.appendNode(c.node, initIndex)
		}

		elementIndex++
		continue
	}
}