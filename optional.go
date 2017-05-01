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
	instance     *optionalParser
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

func (d *optionalDefinition) member(t nodeType) (bool, error) {
	optional, ok := d.registry.definition(d.optionalType)
	if !ok {
		return false, unspecifiedParser(d.optionalName)
	}

	return optional.member(t)
}

func (d *optionalDefinition) generator(t trace, init nodeType, excluded typeList) (generator, error) {
	t = t.extend(d.name)

	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	optional, ok := d.registry.definition(d.optionalType)
	if !ok {
		return nil, unspecifiedParser(d.optionalName)
	}

	if m, err := optional.member(d.typ); err != nil {
		return nil, err
	} else if m {
		return nil, optionalContainingSelf(d.name)
	}

	var initIsMember bool
	if init != 0 {
		if m, err := optional.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	g := &optionalGenerator{
		typ:     d.typ,
		name:    d.name,
		isValid: true,
		instance: &optionalParser{
			typ:          d.typ,
			name:         d.name,
			initIsMember: initIsMember,
			result: &parserResult{
				valid: true,
			},
		},
	}

	d.registry.setGenerator(d.typ, init, excluded, g)

	if excluded.contains(d.typ) {
		g.isValid = false
		return g, nil
	}

	excluded = append(excluded, d.typ)
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

	g.instance.trace = t
	g.instance.cache = c
	g.instance.initNode = init
	g.instance.optional = op

	if g.instance.result.unparsed != nil {
		g.instance.result.unparsed.clear()
	}

	return g.instance
}

func (p *optionalParser) typeName() string   { return p.name }
func (p *optionalParser) nodeType() nodeType { return p.typ }

func (p *optionalParser) parse(t *token) *parserResult {
	p.trace.info("parsing", t)

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
		p.trace.info("missing optional, valid")
		p.result.node = nil
		p.result.fromCache = false
	}

	return p.result
}
