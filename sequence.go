package mml

import "fmt"

type sequenceDefinition struct {
	name      string
	typ       nodeType
	registry  *registry
	itemNames []string
	itemTypes []nodeType
}

type sequenceGenerator struct {
	name           string
	typ            nodeType
	isValid        bool
	initType       nodeType
	generators     []generator
	initGenerators []generator
	initIsMember   []bool
}

type sequenceParser struct {
	name              string
	typ               nodeType
	trace             trace
	cache             *cache
	generators        []generator
	initGenerators    []generator
	initIsMember      []bool
	skip              int
	skippingAfterDone bool
	result            *parserResult
	cacheChecked      bool
	initNode          *node
	itemIndex         int
	currentParser     parser
	initEvaluated     bool
	tokenStack        *tokenStack
}

func sequenceWithoutItems(nodeType string) error {
	return fmt.Errorf("sequence without items: %s", nodeType)
}

func sequenceItemParserInvalid(nodeType string) error {
	return fmt.Errorf("sequence item parser invalid: %s", nodeType)
}

func newSequence(
	r *registry,
	name string,
	nt nodeType,
	itemNames []string,
	itemTypes []nodeType,
) *sequenceDefinition {
	return &sequenceDefinition{
		name:      name,
		typ:       nt,
		registry:  r,
		itemNames: itemNames,
		itemTypes: itemTypes,
	}
}

func (d *sequenceDefinition) typeName() string                { return d.name }
func (d *sequenceDefinition) nodeType() nodeType              { return d.typ }
func (d *sequenceDefinition) member(t nodeType) (bool, error) { return t == d.typ, nil }

func (d *sequenceDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	if len(d.itemTypes) == 0 {
		return nil, sequenceWithoutItems(d.name)
	}

	g := &sequenceGenerator{
		typ:      d.typ,
		name:     d.name,
		isValid:  true,
		initType: init,
	}
	d.registry.setGenerator(d.typ, init, excluded, g)

	if excluded.contains(d.typ) {
		g.isValid = false
		return g, nil
	}

	items := make([]definition, len(d.itemTypes))
	for i, it := range d.itemTypes {
		di, ok := d.registry.definition(it)
		if !ok {
			return nil, unspecifiedParser(d.itemNames[i])
		}

		items[i] = di
	}

	excluded = append(excluded, d.typ)

	var (
		generators     []generator
		initGenerators []generator
		initIsMember   []bool
	)

	initType := init
	for i, di := range items {
		var x typeList
		if i == 0 || initType != 0 {
			x = excluded
		}

		if i > 0 || initType == 0 {
			withoutInit, err := di.generator(t, 0, x)
			if err != nil {
				return nil, err
			}

			if !withoutInit.valid() {
				return nil, sequenceItemParserInvalid(d.name)
			}

			generators = append(generators, withoutInit)
		} else {
			generators = append(generators, nil)
		}

		if initType != 0 {
			withInit, err := di.generator(t, initType, x)
			if err != nil {
				return nil, err
			}

			m, err := di.member(initType)
			if err != nil {
				return nil, err
			}

			if !m && !withInit.valid() {
				g.isValid = false
				return g, nil
			}

			// needs a nil check in the parser
			initGenerators = append(initGenerators, withInit)
			initIsMember = append(initIsMember, m)

			if withInit.valid() || m {
				initType = 0
			}
		}
	}

	g.generators = generators
	g.initGenerators = initGenerators
	g.initIsMember = initIsMember
	return g, nil
}

func (g *sequenceGenerator) typeName() string   { return g.name }
func (g *sequenceGenerator) nodeType() nodeType { return g.typ }
func (g *sequenceGenerator) valid() bool        { return g.isValid }

func (g *sequenceGenerator) finalize(trace) error {
	for _, gi := range g.generators {
		if gi != nil && !gi.valid() {
			return sequenceItemParserInvalid(g.name)
		}
	}

	if g.initType != 0 {
		var foundInvalid bool
		for i, gi := range g.initGenerators {
			if foundInvalid {
				g.initGenerators[i] = nil
				continue
			}

			if gi.valid() {
				continue
			}

			if i == 0 {
				g.isValid = false
				return nil
			}

			g.initGenerators[i] = nil
			foundInvalid = true
		}
	}

	return nil
}

func (g *sequenceGenerator) parser(t trace, c *cache, init *node) parser {
	t = t.extend(g.name)

	return &sequenceParser{
		typ:            g.typ,
		name:           g.name,
		trace:          t,
		cache:          c,
		generators:     g.generators,
		initGenerators: g.initGenerators,
		initIsMember:   g.initIsMember,
		initNode:       init,
	}
}

func (p *sequenceParser) typeName() string   { return p.name }
func (p *sequenceParser) nodeType() nodeType { return p.typ }

func (p *sequenceParser) parse(t *token) *parserResult {
parseLoop:
	for {
		p.trace.info("parsing", t)

		if p.skip > 0 {
			p.skip--

			if p.result == nil {
				p.result = &parserResult{}
			}

			p.result.accepting = true
			return p.result
		}

		if p.skippingAfterDone {
			if p.result == nil {
				p.result = &parserResult{}
			}

			if p.result.unparsed == nil {
				p.result.unparsed = newTokenStack()
			}

			p.result.accepting = false
			p.result.unparsed.push(t)
			return p.result
		}

		// need to check membership before the cache check

		if !p.cacheChecked {
			p.cacheChecked = true

			ct := t
			if p.initNode != nil {
				ct = p.initNode.token
			}

			if n, m, ok := p.cache.get(ct.offset, p.typ); ok {
				if p.result == nil {
					p.result = &parserResult{}
				}

				if p.result.unparsed == nil {
					p.result.unparsed = newTokenStack()
				}

				if m {
					if len(p.initIsMember) > p.itemIndex && !p.initIsMember[p.itemIndex] {
						p.result.valid = true
						p.result.node = n
						p.result.unparsed.push(t)
						p.result.fromCache = true
						p.result.accepting = false
						p.trace.info("found in cache, valid:", p.result.valid, p.result.node)
						return p.result
					}
				} else {
					p.result.valid = false
					p.result.unparsed.push(t)
					p.result.fromCache = true
					p.result.accepting = false
					p.trace.info("found in cache, valid:", p.result.valid, p.result.node)
					return p.result
				}
			}
		}

		if p.currentParser == nil {
			if p.initNode == nil || p.initEvaluated {
				p.trace.debug(p.generators[p.itemIndex] == nil)
				p.currentParser = p.generators[p.itemIndex].parser(p.trace, p.cache, nil)
			} else if p.itemIndex < len(p.initGenerators) {
				if p.initGenerators[p.itemIndex] != nil {
					p.currentParser = p.initGenerators[p.itemIndex].parser(p.trace, p.cache, p.initNode)
				}
			}
		}

		// can be nil with init, only to check membership:
		var ir *parserResult
		if p.currentParser == nil {
			if p.tokenStack == nil {
				p.tokenStack = newTokenStack()
			}

			p.tokenStack.push(t)
		} else {
			ir = p.currentParser.parse(t)
			if ir.accepting {
				if p.result == nil {
					p.result = &parserResult{}
				}

				p.result.accepting = true
				if p.tokenStack != nil && p.tokenStack.has() {
					t = p.tokenStack.pop()
					continue parseLoop
				}

				p.result.accepting = true
				return p.result
			}

			if ir.unparsed != nil {
				if p.tokenStack == nil {
					p.tokenStack = newTokenStack()
				}

				p.tokenStack.merge(ir.unparsed)
			}

			p.currentParser = nil
		}

		p.itemIndex++

		// TODO: when can the result node be nil?
		if ir != nil && ir.valid && ir.node != nil {
			p.initEvaluated = true
			if p.result == nil {
				p.result = &parserResult{}
			}

			if p.result.node == nil {
				p.result.node = &node{
					name: p.name,
					typ:  p.typ,
				}
			}

			p.result.node.append(ir.node)
			if ir.fromCache && p.tokenStack != nil {
				p.skip = p.tokenStack.findCachedNode(ir.node)
			} else {
				p.skip = 0
			}

			if p.itemIndex == len(p.generators) {
				p.trace.info("group done, no more parsers, valid", p.result.node)
				p.result.valid = true

				ct := p.result.node.token
				if ct == nil {
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

				p.skippingAfterDone = p.skip > 0
				p.result.accepting = p.skippingAfterDone
				return p.result
			}

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if ir != nil && ir.valid {
			if p.result == nil {
				p.result = &parserResult{}
			}

			if p.itemIndex == len(p.generators) {
				p.trace.info("group done, zero item, valid", p.result.node)
				p.result.valid = true

				if p.result.node == nil {
					p.result.node = &node{
						name: p.name,
						typ:  p.typ,
					}
				}

				ct := p.result.node.token
				if ct == nil {
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

				p.skippingAfterDone = p.skip > 0
				p.result.accepting = p.skippingAfterDone
				return p.result
			}

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.initNode != nil && !p.initEvaluated && len(p.initIsMember) >= p.itemIndex && p.initIsMember[p.itemIndex-1] {
			p.initEvaluated = true

			if p.result == nil {
				p.result = &parserResult{}
			}

			if p.result.node == nil {
				p.result.node = &node{
					name: p.name,
					typ:  p.typ,
				}
			}

			p.result.node.append(p.initNode)

			if p.itemIndex == len(p.generators) {
				p.trace.info("group done, init item, valid", p.result.node)
				p.result.valid = true

				ct := p.result.node.token
				if ct == nil {
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
				return p.result
			}

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.trace.info("group done, invalid")

		if p.result == nil {
			p.result = &parserResult{}
		}

		p.result.valid = false
		p.result.accepting = false

		var ct *token
		if p.result.node == nil || p.result.node.token == nil {
			if p.tokenStack == nil || !p.tokenStack.has() {
				panic(unexpectedResult(p.name))
			}

			ct = p.tokenStack.peek()
		} else {
			ct = p.result.node.token
		}

		p.cache.set(ct.offset, p.typ, p.result.node, p.result.valid)

		if p.tokenStack != nil {
			if p.result.unparsed == nil {
				p.result.unparsed = newTokenStack()
			}

			p.result.unparsed.merge(p.tokenStack)
		}

		if p.result.node != nil && (p.initNode == nil || len(p.result.node.tokens) > len(p.initNode.tokens)) {
			if p.result.unparsed == nil {
				p.result.unparsed = newTokenStack()
			}

			var i int
			if p.initNode != nil {
				i = len(p.initNode.tokens)
			}

			p.result.unparsed.mergeTokens(p.result.node.tokens[i:])
		}

		p.result.accepting = false
		return p.result
	}
}
