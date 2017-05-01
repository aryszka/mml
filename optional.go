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
	optional generator
	initIsMember bool
	isValid bool
}

type optionalParser struct {
	name  string
	typ   nodeType
	trace trace
	optional parser
	initNode *node
	initIsMember bool
}

func unspecifiedParser(typeName string) error {
	return fmt.Errorf("unspecified parser: %s", typeName)
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

func (d *optionalDefinition) typeName() string {
	return d.name
}

func (d *optionalDefinition) nodeType() nodeType {
	return d.typ
}

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

	g := &optionalGenerator{
		typ:      d.typ,
		name: d.name,
		isValid:  true,
	}

	d.registry.setGenerator(d.typ, init, excluded, g)

	optional, ok := d.registry.definition(d.optionalType)
	if !ok {
		return nil, unspecifiedParser(d.optionalName)
	}

	if m, err := optional.member(d.typ); err != nil {
		return nil, err
	} else if m {
		return nil, optionalContainingSelf(d.name)
	}

	if excluded.contains(d.typ) {
		g.isValid = false
		return g, nil
	}

	excluded = append(excluded, d.typ)
	optGenerator, err := optional.generator(t, init, excluded)
	if err != nil {
		return nil, err
	}

	var initIsMember bool
	if init != 0 {
		if m, err := optional.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	g.optional = optGenerator
	g.initIsMember = initIsMember
	return g, nil
}

func (g *optionalGenerator) typeName() string {
	return g.name
}

func (g *optionalGenerator) nodeType() nodeType {
	return g.typ
}

func (g *optionalGenerator) valid() bool {
	return g.isValid
}

func (g *optionalGenerator) finalize(trace) error {
	if g.optional != nil && !g.optional.valid() {
		g.optional = nil
	}

	return nil
}

func (g *optionalGenerator) parser(t trace, init *node) parser {
	t = t.extend(g.name)

	var op parser
	if g.optional != nil {
		op = g.optional.parser(t, init)
	}

	return &optionalParser{
		name: g.name,
		typ: g.typ,
		trace: t,
		initNode: init,
		optional: op,
		initIsMember: g.initIsMember,
	}
}

func (p *optionalParser) typeName() string {
	return p.name
}

func (p *optionalParser) nodeType() nodeType {
	return p.typ
}

func (p *optionalParser) parse(t *token) *parserResult {
	p.trace.info("parsing", t)

	r := &parserResult{}

	var or *parserResult
	if p.optional != nil {
		or = p.optional.parse(t)
		if or.accepting {
			r.accepting = true
			return r
		}
	}

	r.accepting = false
	r.valid = true
	if or == nil {
		if r.unparsed == nil {
			r.unparsed = newTokenStack()
		}

		r.unparsed.push(t)
	} else if or.unparsed != nil && or.unparsed.has() {
		if r.unparsed == nil {
			r.unparsed = newTokenStack()
		}

		r.unparsed.merge(or.unparsed)
	}

	if or != nil && or.valid {
		p.trace.info("parse done, valid:", r.valid)
		r.node = or.node
		r.fromCache = or.fromCache
	} else if p.initIsMember {
		r.node = p.initNode
		p.trace.info("init node is a member, valid")
		r.fromCache = false
	} else {
		r.node = nil
		p.trace.info("missing optional, valid")
	}

	r.accepting = false
	return r
}
