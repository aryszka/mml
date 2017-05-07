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
	name           string
	typ            nodeType
	trace          trace
	cache          *cache
	generators     []generator
	initGenerators []generator
	initIsMember   []bool
	skip           int
	done           bool
	result         *parserResult
	cacheChecked   bool
	initNode       *node
	itemIndex      int
	currentParser  parser
	initEvaluated  bool
	tokenStack     *tokenStack
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

func (d *sequenceDefinition) typeName() string   { return d.name }
func (d *sequenceDefinition) nodeType() nodeType { return d.typ }

func (d *sequenceDefinition) member(t nodeType, excluded typeList) (bool, error) {
	return !excluded.contains(t) && t == d.typ, nil
}

func (d *sequenceDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	if len(d.itemTypes) == 0 {
		return nil, sequenceWithoutItems(d.name)
	}

	items, err := d.registry.findDefinitions(d.itemTypes)
	if err != nil {
		return nil, err
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

	var (
		generators     []generator
		initGenerators []generator
		initIsMember   []bool
	)

	excluded = append(excluded, d.typ)

	// generators := make([]generator, len(items))
	// for i, item := range items {
	// 	if i == 0 && init != 0 {
	// 		continue
	// 	}

	// 	var x typeList
	// 	if i == 0 {
	// 		x = excluded
	// 	}

	// 	gi, err := item.generator(t, 0, x)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if !gi.valid() {
	// 		return nil, sequenceItemParserInvalid(d.name)
	// 	}

	// 	generators[i] = gi
	// }

	// var (
	// 	initGenerators []generator
	// 	initIsMember []bool
	// )

	// if init != 0 {
	// 	for i, item := range items {
	// 		gi, err := item.generator(t, init, excluded)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		m, err := item.member(init)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		if !m && !gi.valid() {
	// 			g.isValid = false
	// 			return g, nil
	// 		}

	// 		initGenerators = append(initGenerators, gi)
	// 		initIsMember = append(initIsMember, m)
	// 	}
	// }

	// TODO: is it handled when the item generator becomes invalid after the generator was created?

	// maybe three lists are needed

	// TODO: always use the excluded for the first, in repetition, too

	hasInit := init != 0
	for i, di := range items {
		var x typeList
		if i == 0 || hasInit {
			x = excluded
		}

		if i > 0 || !hasInit {
			withoutInit, err := di.generator(t, 0, x)
			if err != nil {
				return nil, err
			}

			if !withoutInit.valid() {
				t.debug(d.name, i)
				return nil, sequenceItemParserInvalid(d.name)
			}

			generators = append(generators, withoutInit)
		} else {
			generators = append(generators, nil)
		}

		if hasInit {
			withInit, err := di.generator(t, init, x)
			if err != nil {
				return nil, err
			}

			m, err := di.member(init, excluded)
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
				hasInit = false
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
		result:         &parserResult{},
	}
}

func (p *sequenceParser) typeName() string   { return p.name }
func (p *sequenceParser) nodeType() nodeType { return p.typ }

// TODO: what happens if there is an init node but it's optional and returns valid but it's not consumed

func (p *sequenceParser) parse(t *token) *parserResult {
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
				false,
				len(p.initIsMember) > p.itemIndex && p.initIsMember[p.itemIndex],
			) {
				p.trace.info("found in cache, valid:", p.result.valid)
				p.done = true
				return p.result
			}
		}

		if p.currentParser == nil {
			p.trace.debug(
				"taking next parser",
				p.initNode == nil,
				p.initEvaluated,
				p.itemIndex < len(p.initGenerators),
			)
			if p.initNode == nil || p.initEvaluated {
				p.currentParser = p.generators[p.itemIndex].parser(p.trace, p.cache, nil)
			} else if p.itemIndex < len(p.initGenerators) {
				if p.initGenerators[p.itemIndex] != nil {
					p.currentParser = p.initGenerators[p.itemIndex].parser(
						p.trace,
						p.cache,
						p.initNode,
					)
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

		if ir != nil && ir.valid && ir.node != nil {
			p.initEvaluated = true
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
				p.trace.info("sequence done, no more parsers, valid", p.result.node)
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

				p.done = true

				p.result.accepting = p.skip > 0
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
			if p.itemIndex == len(p.generators) {
				p.trace.info("sequence done, zero item, valid", p.result.node)
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

				p.done = true
				p.result.accepting = p.skip > 0
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

			if p.result.node == nil {
				p.result.node = &node{
					name: p.name,
					typ:  p.typ,
				}
			}

			p.result.node.append(p.initNode)

			if p.itemIndex == len(p.generators) {
				p.trace.info("sequence done, init item, valid", p.result.node)
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
				p.done = false
				return p.result
			}

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.trace.info("sequence done, invalid")

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
		p.done = false
		return p.result
	}
}
