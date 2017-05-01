package mml

import "fmt"

type primitiveDefinition struct {
	name      string
	typ       nodeType
	registry  *registry
	tokenType tokenType
}

type primitiveGenerator struct {
	name      string
	typ       nodeType
	isValid   bool
	tokenType tokenType
}

type primitiveParser struct {
	name      string
	typ       nodeType
	trace     trace
	tokenType tokenType
}

func unexpectedInitNode(typeName, initTypeName string) error {
	return fmt.Errorf("unexpected init node: %s, %s", typeName, initTypeName)
}

func newPrimitive(r *registry, name string, nt nodeType, tt tokenType) *primitiveDefinition {
	return &primitiveDefinition{
		name:      name,
		typ:       nt,
		registry:  r,
		tokenType: tt,
	}
}

func (d *primitiveDefinition) typeName() string                { return d.name }
func (d *primitiveDefinition) nodeType() nodeType              { return d.typ }
func (d *primitiveDefinition) member(t nodeType) (bool, error) { return t == d.typ, nil }

func (d *primitiveDefinition) generator(_ trace, init nodeType, excluded typeList) (generator, error) {
	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	g := &primitiveGenerator{
		name:      d.name,
		typ:       d.typ,
		isValid:   !excluded.contains(d.typ) && init == 0,
		tokenType: d.tokenType,
	}

	d.registry.setGenerator(d.typ, init, excluded, g)
	return g, nil
}

func (g *primitiveGenerator) typeName() string     { return g.name }
func (g *primitiveGenerator) nodeType() nodeType   { return g.typ }
func (g *primitiveGenerator) valid() bool          { return g.isValid }
func (g *primitiveGenerator) finalize(trace) error { return nil }

func (g *primitiveGenerator) parser(t trace, _ *cache, init *node) parser {
	if init != nil {
		panic(unexpectedInitNode(g.name, init.name))
	}

	return &primitiveParser{
		name:      g.name,
		typ:       g.typ,
		trace:     t.extend(g.name),
		tokenType: g.tokenType,
	}
}

func (p *primitiveParser) typeName() string   { return p.name }
func (p *primitiveParser) nodeType() nodeType { return p.typ }

func (p *primitiveParser) parse(t *token) *parserResult {
	p.trace.info("parsing", t)

	if t.typ == p.tokenType {
		p.trace.info("valid")
		return &parserResult{
			valid: true,
			node: &node{
				name:   p.name,
				typ:    p.typ,
				token:  t,
				tokens: []*token{t},
			},
		}
	}

	p.trace.info("invalid")
	up := newTokenStack()
	up.push(t)
	return &parserResult{
		unparsed: up,
	}
}
