package mml

// TODO: find a good name for this

type repetitionDefinition struct {
	name     string
	typ      nodeType
	registry *registry
	itemName string
	itemType nodeType
}

type repetitionGenerator struct {
	name         string
	typ          nodeType
	isValid      bool
	first        generator
	rest         generator
	initIsMember bool
}

type repetitionParser struct {
	name          string
	typ           nodeType
	trace         trace
	cache         *cache
	currentParser parser
	rest          generator
	result        *parserResult
	tokenStack    *tokenStack
	cacheChecked  bool
	initNode      *node
	initIsMember  bool
	initEvaluated bool
	skip          int
	done          bool // TODO: rename this field
}

func newRepetition(
	r *registry,
	name string,
	nt nodeType,
	itemName string,
	itemType nodeType,
) *repetitionDefinition {
	return &repetitionDefinition{
		name:     name,
		typ:      nt,
		registry: r,
		itemName: itemName,
		itemType: itemType,
	}
}

func (d *repetitionDefinition) typeName() string   { return d.name }
func (d *repetitionDefinition) nodeType() nodeType { return d.typ }

func (d *repetitionDefinition) member(t nodeType) (bool, error) {
	return t == d.typ, nil
}

func (d *repetitionDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	item, err := d.registry.findDefinition(d.itemType)
	if err != nil {
		return nil, err
	}

	g := &repetitionGenerator{
		typ:     d.typ,
		isValid: true,
		name:    d.name,
	}

	d.registry.setGenerator(d.typ, init, excluded, g)
	if excluded.contains(d.typ) {
		g.isValid = false
		return g, nil
	}

	excluded = append(excluded, d.typ)

	var initIsMember bool
	if init != 0 {
		if m, err := item.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	// TODO: revert the membership bullshit

	g.initIsMember = initIsMember

	first, err := item.generator(t, init, excluded)
	if err != nil {
		return nil, err
	}

	rest, err := item.generator(t, 0, typeList{d.typ})
	if err != nil {
		return nil, err
	}

	if !rest.valid() {
		panic(requiredParserInvalid(d.itemName))
	}

	g.first = first
	g.rest = rest

	return g, nil
}

func (g *repetitionGenerator) typeName() string   { return g.name }
func (g *repetitionGenerator) nodeType() nodeType { return g.typ }
func (g *repetitionGenerator) valid() bool        { return g.isValid }

func (g *repetitionGenerator) finalize(trace) error {
	if g.first != nil && !g.first.valid() {
		g.first = nil
	}

	if g.rest != nil && !g.rest.valid() {
		g.rest = nil
	}

	return nil
}

func (g *repetitionGenerator) parser(t trace, c *cache, init *node) parser {
	t = t.extend(g.name)

	var currentParser parser
	if g.first != nil {
		currentParser = g.first.parser(t, c, init)
	}

	return &repetitionParser{
		typ:           g.typ,
		name:          g.name,
		trace:         t,
		cache:         c,
		initNode:      init,
		rest:          g.rest,
		currentParser: currentParser,
		initIsMember:  g.initIsMember,
		result: &parserResult{
			valid: true,
		},
	}
}

func (p *repetitionParser) typeName() string   { return p.name }
func (p *repetitionParser) nodeType() nodeType { return p.typ }

// TODO: check if the node token is ensured in every such case where the node doesn't have any tokens

func (p *repetitionParser) parse(t *token) *parserResult {
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
				p.initIsMember,
			) {
				p.trace.info("found in cache, valid:", p.result.valid)
				p.done = true
				return p.result
			}
		}

		if p.currentParser == nil {
			p.initEvaluated = true
			if p.initIsMember {
				p.result.ensureNode(p.name, p.typ)
				p.result.node.append(p.initNode)
				p.currentParser = p.rest.parser(p.trace, p.cache, nil)
				continue parseLoop
			}

			p.result.ensureNode(p.name, p.typ)
			p.result.node.token = t
			p.result.unparsed = newTokenStack()
			p.result.unparsed.push(t)

			p.trace.info("repetition done, valid, no parser")
			p.done = true
			p.cache.set(t.offset, p.typ, p.result.node, p.result.valid)
			return p.result
		}

		ir := p.currentParser.parse(t)
		if ir.accepting {
			if st, ok := p.tokenStack.popIfAny(); ok {
				t = st
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.tokenStack = mergeStack(p.tokenStack, ir.unparsed)
		if ir.valid && ir.node != nil && len(ir.node.tokens) > 0 {
			p.initEvaluated = true

			p.result.ensureNode(p.name, p.typ)
			p.result.node.append(ir.node)
			if ir.fromCache {
				p.skip = p.tokenStack.findCachedNode(ir.node)
			}

			p.currentParser = p.rest.parser(p.trace, p.cache, nil)
			if st, ok := p.tokenStack.popIfAny(); ok {
				t = st
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.initIsMember && !p.initEvaluated {
			p.initEvaluated = true

			p.result.ensureNode(p.name, p.typ)
			p.result.node.append(p.initNode)
			p.currentParser = p.rest.parser(p.trace, p.cache, nil)

			if st, ok := p.tokenStack.popIfAny(); ok {
				t = st
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.result.accepting = p.skip > 0
		p.result.valid = true
		p.result.mergeStack(p.tokenStack)

		p.result.ensureNode(p.name, p.typ)
		if p.result.node.token == nil {
			p.result.assertUnparsed(p.name)
			p.result.node.token = p.result.unparsed.peek()
		}

		p.trace.info("repetition done, valid")
		p.done = true
		p.cache.set(p.result.node.token.offset, p.typ, p.result.node, p.result.valid)
		return p.result
	}
}
