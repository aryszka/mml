package mml

import "fmt"

type choiceDefinition struct {
	name         string
	typ          nodeType
	registry     *registry
	elementNames []string
	elementTypes []nodeType
	elements     []definition
}

type choiceGenerator struct {
	name         string
	typ          nodeType
	isValid      bool
	generators   [][]generator
	initIsMember bool
}

type choiceParser struct {
	name          string
	typ           nodeType
	trace         trace
	cache         *cache
	generators    [][]generator
	initIsMember  bool
	initNode      *node
	result        *parserResult
	currentParser parser
	skip          int
	done          bool
	cacheChecked  bool
	tokenStack    *tokenStack
	elementIndex  int
	initTypeIndex int
}

func choiceWithoutElements(nodeType string) error {
	return fmt.Errorf("choice without elements: %s", nodeType)
}

func newChoice(
	r *registry,
	name string,
	nt nodeType,
	elementNames []string,
	elementTypes []nodeType,
) *choiceDefinition {
	return &choiceDefinition{
		name:         name,
		typ:          nt,
		registry:     r,
		elementNames: elementNames,
		elementTypes: elementTypes,
	}
}

func (d *choiceDefinition) typeName() string   { return d.name }
func (d *choiceDefinition) nodeType() nodeType { return d.typ }

func (d *choiceDefinition) expand(ignore typeList) ([]definition, error) {
	if ignore.contains(d.typ) {
		return nil, nil
	}

	var definitions []definition
	for _, et := range d.elementTypes {
		ed, err := d.registry.findDefinition(et)
		if err != nil {
			return nil, err
		}

		if xd, ok := ed.(expander); ok {
			edx, err := xd.expand(append(ignore, d.typ))
			if err != nil {
				return nil, err
			}

			definitions = append(definitions, edx...)
		} else if !ignore.contains(et) {
			definitions = append(definitions, ed)
		}
	}

	return definitions, nil
}

func (d *choiceDefinition) checkExpand() error {
	if len(d.elements) > 0 {
		return nil
	}

	elements, err := d.expand(nil)
	if err != nil {
		return err
	}

	if len(elements) == 0 {
		return choiceWithoutElements(d.name)
	}

	d.elements = elements
	return nil
}

func (d *choiceDefinition) member(t nodeType, excluded typeList) (bool, error) {
	if err := d.checkExpand(); err != nil {
		return false, err
	}

	for _, e := range d.elements {
		if m, err := e.member(t, excluded); m || err != nil {
			return m, err
		}
	}

	return false, nil
}

func (d *choiceDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	g := &choiceGenerator{
		name:    d.name,
		typ:     d.typ,
		isValid: true,
	}
	d.registry.setGenerator(d.typ, init, excluded, g)

	if err := d.checkExpand(); err != nil {
		return nil, err
	}

	expandedTypes := make([]nodeType, len(d.elements))
	for i, e := range d.elements {
		expandedTypes[i] = e.nodeType()
	}

	generators := make([][]generator, len(d.elements)+1)
	for i, it := range append([]nodeType{init}, expandedTypes...) {
		g := make([]generator, len(d.elements))
		for j, e := range d.elements {
			ge, err := e.generator(t, it, excluded)
			if err != nil {
				return nil, err
			}

			if ge.valid() {
				g[j] = ge
			}
		}

		generators[i] = g
	}

	var initIsMember bool
	if init != 0 {
		if m, err := d.member(init, excluded); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	if !initIsMember && (len(generators[0]) == 0 || generators[0][0] == nil) {
		g.isValid = false
		return g, nil
	}

	g.generators = generators
	g.initIsMember = initIsMember
	return g, nil
}

func (g *choiceGenerator) typeName() string   { return g.name }
func (g *choiceGenerator) nodeType() nodeType { return g.typ }
func (g *choiceGenerator) valid() bool        { return g.isValid }

func (g *choiceGenerator) finalize(trace) error {
	var hasValid bool

	for _, gs := range g.generators {
		for i, gg := range gs {
			if gg != nil && gg.valid() {
				hasValid = true
				continue
			}

			gs[i] = nil
		}
	}

	g.isValid = hasValid || g.initIsMember
	return nil
}

func (g *choiceGenerator) parser(t trace, c *cache, init *node) parser {
	t = t.extend(g.name)

	result := &parserResult{}
	if g.initIsMember {
		result.node = init
	}

	var currentParser parser
	if g.generators[0][0] != nil {
		currentParser = g.generators[0][0].parser(t, c, init)
	}

	return &choiceParser{
		name:          g.name,
		typ:           g.typ,
		trace:         t,
		cache:         c,
		generators:    g.generators,
		initIsMember:  g.initIsMember,
		initNode:      init,
		result:        result,
		currentParser: currentParser,
	}
}

func (p *choiceParser) typeName() string   { return p.name }
func (p *choiceParser) nodeType() nodeType { return p.typ }

// TODO: reconsider choice caching

// TODO: don't use the previous node as init if it was empty

func (p *choiceParser) parse(t *token) *parserResult {
parseLoop:
	for {
		traceToken(p.trace, t, p.initNode, p.result)

		var accepting, ret bool
		if p.skip, accepting, ret = checkSkip(p.skip, p.done); ret {
			p.result.accepting = accepting
			return p.result
		}

		if !p.cacheChecked {
			p.cacheChecked = true

			if p.result.fillFromCache(
				p.cache,
				p.typ,
				t,
				p.initNode,
				p.initIsMember,
				false,
			) {
				p.trace.info("found in cache, valid:", p.result.valid)
				p.done = true
				return p.result
			}
		}

		var er *parserResult
		if p.currentParser != nil {
			er = p.currentParser.parse(t)
			if er.accepting {
				if p.tokenStack != nil && p.tokenStack.has() {
					t = p.tokenStack.pop()
					continue parseLoop
				}

				p.result.accepting = true
				return p.result
			}

			if er.unparsed != nil {
				if p.tokenStack == nil {
					p.tokenStack = newTokenStack()
				}

				p.tokenStack.merge(er.unparsed)
			}
		} else {
			if p.tokenStack == nil {
				p.tokenStack = newTokenStack()
			}

			p.tokenStack.push(t)
		}

		if er == nil || !er.valid {
			for {
				p.elementIndex++

				if p.elementIndex == len(p.generators[p.initTypeIndex]) {
					break
				}

				if p.generators[p.initTypeIndex][p.elementIndex] != nil {
					break
				}
			}

			if p.elementIndex == len(p.generators[p.initTypeIndex]) {
				p.trace.info("done, valid:", p.result.node != nil, p.result.node)
				p.result.accepting = false
				p.result.valid = p.result.node != nil

				var ct *token
				if p.result.node != nil {
					ct = p.result.node.token
				} else if p.initNode != nil {
					ct = p.initNode.token
				} else {
					if p.tokenStack == nil || !p.tokenStack.has() {
						panic(unexpectedResult(p.name))
					}

					ct = p.tokenStack.peek()
				}

				p.cache.set(ct.offset, p.typ, p.result.node, p.result.valid)

				if p.tokenStack != nil {
					if p.result.unparsed == nil {
						p.result.unparsed = newTokenStack()
					}

					p.result.unparsed.merge(p.tokenStack)
				}

				p.result.accepting = false
				p.done = true
				return p.result
			}

			p.currentParser = p.generators[p.initTypeIndex][p.elementIndex].parser(p.trace, p.cache, p.result.node)

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		// TODO: test valid optional as the only match in a choice
		if p.result.node == nil ||
			er != nil && er.node != nil &&
				len(p.result.node.tokens) < len(er.node.tokens) {

			p.result.node = er.node

			p.initTypeIndex = p.elementIndex + 1
			if p.initTypeIndex < len(p.generators) {
				p.elementIndex = 0
				for {
					if p.generators[p.initTypeIndex][p.elementIndex] != nil {
						break
					}

					p.elementIndex++

					if p.elementIndex == len(p.generators[p.initTypeIndex]) {
						break
					}
				}
			}

			if er.fromCache {
				if er.node != nil && p.tokenStack != nil {
					p.skip = p.tokenStack.findCachedNode(er.node)
				} else {
					p.skip = 0
				}
			}
		} else {
			for {
				p.elementIndex++

				if p.elementIndex == len(p.generators[p.initTypeIndex]) {
					break
				}

				if p.generators[p.initTypeIndex][p.elementIndex] != nil {
					break
				}
			}
		}

		if p.initTypeIndex < len(p.generators) && p.elementIndex < len(p.generators[p.initTypeIndex]) {
			p.currentParser = p.generators[p.initTypeIndex][p.elementIndex].parser(p.trace, p.cache, p.result.node)

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.result.accepting = false
		p.result.valid = p.result.node != nil
		p.trace.info("done, valid:", p.result.valid, p.result.node)

		var ct *token
		if p.result.node != nil {
			ct = p.result.node.token
		} else if p.initNode != nil {
			ct = p.initNode.token
		} else {
			if p.tokenStack == nil || !p.tokenStack.has() {
				panic(unexpectedResult(p.name))
			}

			ct = p.tokenStack.peek()
		}

		p.cache.set(ct.offset, p.typ, p.result.node, p.result.valid)

		if p.tokenStack != nil {
			if p.result.unparsed == nil {
				p.result.unparsed = newTokenStack()
			}

			p.result.unparsed.merge(p.tokenStack)
		}

		p.result.accepting = false
		p.done = true
		if p.skip > 0 {
			p.result.accepting = true
			return p.result
		}

		return p.result
	}
}
