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
	instance  *primitiveParser
}

type primitiveParser struct {
	name      string
	typ       nodeType
	trace     trace
	cache     *cache
	tokenType tokenType
	success   *parserResult
	fail      *parserResult
}

func unexpectedInitNode(typeName, initTypeName string) error {
	return fmt.Errorf("unexpected init node: %s, %s", typeName, initTypeName)
}

// TODO: cache primitive due to node allocation. Worth checking the cache everywhere before the parser instance
// is created.

func newPrimitive(r *registry, name string, nt nodeType, tt tokenType) *primitiveDefinition {
	return &primitiveDefinition{
		name:      name,
		typ:       nt,
		registry:  r,
		tokenType: tt,
	}
}

func (d *primitiveDefinition) typeName() string   { return d.name }
func (d *primitiveDefinition) nodeType() nodeType { return d.typ }

func (d *primitiveDefinition) member(t nodeType, excluded typeList) (bool, error) {
	return !excluded.contains(t) && t == d.typ, nil
}

func (d *primitiveDefinition) generator(_ trace, init nodeType, excluded typeList) (generator, error) {
	if g, ok := d.registry.generator(d.typ, init, excluded); ok {
		return g, nil
	}

	g := &primitiveGenerator{
		name:      d.name,
		typ:       d.typ,
		isValid:   !excluded.contains(d.typ) && init == 0,
		tokenType: d.tokenType,
		instance: &primitiveParser{
			name:      d.name,
			typ:       d.typ,
			tokenType: d.tokenType,
			success: &parserResult{
				valid: true,
			},
			fail: &parserResult{
				unparsed: withLength(1),
			},
		},
	}

	d.registry.setGenerator(d.typ, init, excluded, g)
	return g, nil
}

func (g *primitiveGenerator) typeName() string     { return g.name }
func (g *primitiveGenerator) nodeType() nodeType   { return g.typ }
func (g *primitiveGenerator) valid() bool          { return g.isValid }
func (g *primitiveGenerator) finalize(trace) error { return nil }

// TODO: can we always instantiate the parser when the token is already there? It would help checking the cache
// before unnecessary instantiation

func (g *primitiveGenerator) parser(t trace, c *cache, init *node) parser {
	if init != nil {
		panic(unexpectedInitNode(g.name, init.name))
	}

	g.instance.trace = t.extend(g.name)
	g.instance.cache = c
	g.instance.fail.unparsed.clear()
	g.instance.success.fromCache = false
	return g.instance
}

func (p *primitiveParser) typeName() string   { return p.name }
func (p *primitiveParser) nodeType() nodeType { return p.typ }

func (p *primitiveParser) parse(t *token) *parserResult {
	p.trace.info("parsing", t)

	if n, m, ok := p.cache.get(t.offset, p.typ); ok {
		if m {
			p.trace.info("found in cache, valid")
			p.success.node = n
			p.success.fromCache = true
			p.success.unparsed.push(t)
			return p.success
		} else {
			p.trace.info("found in cache, invalid")
			p.fail.unparsed.push(t)
			return p.fail
		}
	}

	if t.typ == p.tokenType {
		p.trace.info("valid")

		p.success.node = &node{
			name:   p.name,
			typ:    p.typ,
			token:  t,
			tokens: []*token{t},
		}

		return p.success
	}

	p.trace.info("invalid")
	p.fail.unparsed.push(t)
	return p.fail
}
