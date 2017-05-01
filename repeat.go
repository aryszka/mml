package mml

// TODO: find a good name for this

type repeatDefinition struct {
	name     string
	typ      nodeType
	registry *registry
	itemName string
	itemType nodeType
}

type repeatGenerator struct {
	name         string
	typ          nodeType
	isValid      bool
	first        generator
	rest         generator
	initIsMember bool
}

type repeatParser struct {
	name              string
	typ               nodeType
	trace             trace
	cache             *cache
	currentParser     parser
	skip              int
	skippingAfterDone bool
	result            *parserResult
	cacheChecked      bool
	initNode          *node
	initIsMember      bool
	tokenStack        *tokenStack
	initEvaluated     bool
	rest              generator
}

func newRepeat(
	r *registry,
	name string,
	nt nodeType,
	itemName string,
	itemType nodeType,
) *repeatDefinition {
	return &repeatDefinition{
		name:     name,
		typ:      nt,
		registry: r,
		itemName: itemName,
		itemType: itemType,
	}
}

func (d *repeatDefinition) typeName() string                { return d.name }
func (d *repeatDefinition) nodeType() nodeType              { return d.typ }
func (d *repeatDefinition) member(t nodeType) (bool, error) { return t == d.typ, nil }

func (d *repeatDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	g := &repeatGenerator{typ: d.typ, isValid: true, name: d.name}
	d.registry.setGenerator(d.typ, init, excluded, g)

	item, ok := d.registry.definition(d.itemType)
	if !ok {
		return nil, unspecifiedParser(d.itemName)
	}

	if excluded.contains(d.typ) {
		g.isValid = false
		return g, nil
	}

	first, err := item.generator(t, init, append(excluded, d.typ))
	if err != nil {
		return nil, err
	}

	rest, err := item.generator(t, 0, nil)
	if err != nil {
		return nil, err
	}

	if !rest.valid() {
		panic(requiredParserInvalid(d.itemName))
	}

	var initIsMember bool
	if init != 0 {
		if m, err := item.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	g.first = first
	g.rest = rest
	g.initIsMember = initIsMember
	return g, nil
}

func (g *repeatGenerator) typeName() string   { return g.name }
func (g *repeatGenerator) nodeType() nodeType { return g.typ }
func (g *repeatGenerator) valid() bool        { return g.isValid }

func (g *repeatGenerator) finalize(trace) error {
	if g.first != nil && !g.first.valid() {
		g.first = nil
	}

	if g.rest != nil && !g.rest.valid() {
		g.rest = nil
	}

	return nil
}

func (g *repeatGenerator) parser(t trace, c *cache, init *node) parser {
	t = t.extend(g.name)

	var currentParser parser
	if g.first != nil {
		currentParser = g.first.parser(t, c, init)
	}

	return &repeatParser{
		typ:           g.typ,
		name:          g.name,
		trace:         t,
		cache:         c,
		initNode:      init,
		rest:          g.rest,
		currentParser: currentParser,
		initIsMember:  g.initIsMember,
	}
}

func (p *repeatParser) typeName() string   { return p.name }
func (p *repeatParser) nodeType() nodeType { return p.typ }

func (p *repeatParser) parse(t *token) *parserResult {
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
					if !p.initIsMember {
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
			// the case when the item cannot accept the init node

			if p.tokenStack == nil {
				p.tokenStack = newTokenStack()
			}

			p.tokenStack.push(t)
			p.trace.info("sequence done, valid, no parser", p.result.node)
			p.skippingAfterDone = p.skip > 0

			if p.result == nil {
				p.result = &parserResult{}
			}

			if p.result.unparsed == nil {
				p.result.unparsed = newTokenStack()
			}

			p.result.accepting = p.skippingAfterDone
			p.result.valid = true
			p.result.unparsed.merge(p.tokenStack)

			// NOTE: this was not set in parse4
			// maybe every node should have a token
			if p.result.node.token == nil {
				if !p.result.unparsed.has() {
					panic(unexpectedResult(p.name))
				}

				p.result.node.token = p.result.unparsed.peek()
			}

			// NOTE: this was cached in parse4 only if there were nodes in the sequence
			var ct *token
			if p.result.node == nil {
				ct = p.result.unparsed.peek()
			} else {
				ct = p.result.node.token
			}

			p.cache.set(ct.offset, p.typ, p.result.node, p.result.valid)
			p.result.accepting = false
			return p.result
		}

		ir := p.currentParser.parse(t)
		if ir.accepting {
			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			if p.result == nil {
				p.result = &parserResult{}
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

		if ir.valid && ir.node != nil && len(ir.node.tokens) > 0 {
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
			if ir.fromCache {
				if p.tokenStack == nil {
					p.skip = 0
				} else {
					p.skip = p.tokenStack.findCachedNode(ir.node)
				}
			}

			p.currentParser = p.rest.parser(p.trace, p.cache, nil)

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.initIsMember && !p.initEvaluated {
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
			p.currentParser = p.rest.parser(p.trace, p.cache, nil)

			if p.tokenStack != nil && p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.result == nil {
			p.result = &parserResult{}
		}

		p.trace.info("sequence done, valid, item not valid", p.result.node)
		p.skippingAfterDone = p.skip > 0
		p.result.accepting = p.skippingAfterDone
		p.result.valid = true

		if p.tokenStack != nil {
			if p.result.unparsed == nil {
				p.result.unparsed = newTokenStack()
			}

			p.result.unparsed.merge(p.tokenStack)
		}

		// NOTE: this was not set in parse4
		// maybe every node should have a token
		if p.result.node == nil {
			p.result.node = &node{
				typ:  p.typ,
				name: p.name,
			}

			if p.result.unparsed == nil || !p.result.unparsed.has() {
				panic(unexpectedResult(p.name))
			}

			p.result.node.token = p.result.unparsed.peek()
		}

		// NOTE: this was cached in parse4 only if there were nodes in the sequence
		p.cache.set(p.result.node.token.offset, p.typ, p.result.node, p.result.valid)
		return p.result
	}
}
