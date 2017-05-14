package mml

import "fmt"

type optionalDefinition struct {
	name         string
	typ          nodeType
	registry     *registry
	optionalName string
	optionalType nodeType
}

type optionalGenerator struct {
	name         string
	typ          nodeType
	optionalName string
	optionalType nodeType
	optional     generator
	isValid      bool
	initIsMember bool
}

type optionalParser struct {
	name         string
	typ          nodeType
	trace        trace
	cache        *cache
	optional     parser
	initNode     *node
	initIsMember bool
	result       *parserResult
}

// TODO: the empty is a bad idea, everything should have its own node

func optionalContainingSelf(nodeType string) error {
	return fmt.Errorf("optional containing self: %s", nodeType)
}

func newOptional(
	r *registry,
	name string,
	nt nodeType,
	optionalName string,
	optionalType nodeType,
) *optionalDefinition {
	return &optionalDefinition{
		name:         name,
		typ:          nt,
		registry:     r,
		optionalName: optionalName,
		optionalType: optionalType,
	}
}

func (d *optionalDefinition) typeName() string   { return d.name }
func (d *optionalDefinition) nodeType() nodeType { return d.typ }

// func (d *optionalDefinition) expand(ignore typeList) ([]definition, error) {
// 	optional, err := d.registry.findDefinition(d.optionalType)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var defs []definition
// 	xd, ok := optional.(expander);
// 	if ok {
// 		defs, err = xd.expand(append(ignore, d.typ))
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		defs = []definition{optional}
// 	}
//
// 	defs = append(defs, emptyParser)
// 	return defs, nil
// }

func (d *optionalDefinition) member(t nodeType) (bool, error) {
	optional, err := d.registry.findDefinition(d.optionalType)
	if err != nil {
		return false, err
	}

	if t == d.typ {
		return true, nil
	}

	return optional.member(t)
}

// TODO: for the test "optional expression as expression, empty input", introduce the 'nothing' parser?

func (d *optionalDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	optional, err := d.registry.findDefinition(d.optionalType)
	if err != nil {
		return nil, err
	}

	// TODO: why not?
	// the problem is when it contains only itself
	// it makes no sense. This check should have failed.
	// what happens if it contains itself through a repetition?
	// maybe the same check should go into the repetition, too
	// or just still why not, and only check that every rule expands to terminals
	// if m, err := optional.member(d.typ); err != nil {
	// 	return nil, err
	// } else if m {
	// 	return nil, optionalContainingSelf(d.name)
	// }

	var initIsMember bool
	if init != 0 {
		if m, err := optional.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	g := &optionalGenerator{
		typ:          d.typ,
		name:         d.name,
		isValid:      true,
		initIsMember: initIsMember,
	}

	d.registry.setGenerator(d.typ, init, excluded, g)

	if excluded.contains(d.typ) {
		g.isValid = false
		return g, nil
	}

	optGenerator, err := optional.generator(t, init, excluded)
	if err != nil {
		return nil, err
	}

	g.optional = optGenerator
	return g, nil
}

func (g *optionalGenerator) typeName() string   { return g.name }
func (g *optionalGenerator) nodeType() nodeType { return g.typ }
func (g *optionalGenerator) valid() bool        { return g.isValid }

func (g *optionalGenerator) finalize(trace) error {
	if g.optional != nil && !g.optional.valid() {
		g.optional = nil
	}

	return nil
}

func (g *optionalGenerator) parser(t trace, c *cache, init *node) parser {
	t = t.extend(g.name)

	var op parser
	if g.optional != nil {
		op = g.optional.parser(t, c, init)
	}

	return &optionalParser{
		typ:          g.typ,
		name:         g.name,
		trace:        t,
		cache:        c,
		initIsMember: g.initIsMember,
		initNode:     init,
		optional:     op,
		result: &parserResult{
			valid: true,
		},
	}
}

func (p *optionalParser) typeName() string   { return p.name }
func (p *optionalParser) nodeType() nodeType { return p.typ }

// TODO: fix the cache so that it can be used in the optional and the choice

func (p *optionalParser) parse(t *token) *parserResult {
	traceToken(p.trace, t, p.initNode, p.result)

	if p.result.fillFromCache(
		p.cache,
		p.typ,
		t,
		p.initNode,
		p.initIsMember,
		false,
	) {
		p.trace.info("found in cache, valid:", p.result.valid)
		return p.result
	}

	var or *parserResult
	if p.optional != nil {
		or = p.optional.parse(t)
		if or.accepting {
			p.result.accepting = true
			return p.result
		}
	}

	p.result.accepting = false
	p.result.valid = true

	if or == nil {
		if p.result.unparsed == nil {
			p.result.unparsed = newTokenStack()
		}

		p.result.unparsed.push(t)
	} else if or.unparsed != nil && or.unparsed.has() {
		if p.result.unparsed == nil {
			p.result.unparsed = newTokenStack()
		}

		p.result.unparsed.merge(or.unparsed)
	}

	if or != nil && or.valid {
		p.trace.info("parse done, valid")
		p.result.node = or.node
		p.result.fromCache = or.fromCache
	} else if p.initIsMember {
		p.trace.info("init node is a member, valid")
		p.result.node = p.initNode
		p.result.fromCache = false
	} else {
		p.trace.info("empty optional, valid")
		p.result.node = nil
		p.result.fromCache = false
	}

	var ct *token
	if p.result.node != nil {
		ct = p.result.node.token
	} else if p.initNode != nil {
		ct = p.initNode.token
	} else {
		if p.result.unparsed == nil || !p.result.unparsed.has() {
			panic(unexpectedResult(p.name))
		}

		ct = p.result.unparsed.peek()
	}

	p.cache.set(ct.offset, p.typ, p.result.node, true)
	return p.result
}
