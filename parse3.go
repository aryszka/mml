package mml

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type node struct {
	token
	typ   string
	nodes []node
}

type parseResult struct {
	accepting bool
	valid     bool
	unparsed  []token
	node      node
}

type parser interface {
	parse(token) parseResult
	path() []string
	name() string
	out(...interface{})
}

type generator interface {
	canCreate(node, []string) bool
	create([]string, node, []string) parser
	name() string
	out(...interface{})
	member(string) bool
}

type baseParser struct {
	p []string
}

type baseGenerator struct {
	node string
}

type primitiveParser struct {
	baseParser
	accepting bool
	token     tokenType
	node      node
}

type primitiveGenerator struct {
	baseGenerator
	token tokenType
}

type optionalParser struct {
	baseParser
	optional         generator
	optionalAccepted bool
	init             node
	excluded         []string
	parser           parser
}

type sequenceParser struct {
	baseParser
	node          node
	init          node
	itemGenerator generator
	currentParser parser
	queue         []token
	excluded      []string
}

type groupParser struct {
	baseParser
	node          node
	init          node
	generators    []generator
	currentParser parser
	queue         []token
	excluded      []string
	accepted      []token
	itemAccepted  []token
}

type unionParser struct {
	baseParser
	currentParser    parser
	generators       []generator
	activeGenerators []generator
	node             node
	valid            bool
	queue            []token
	excluded         []string
	init             node
	hasAccepted      bool
}

type optionalGenerator struct {
	baseGenerator
	optional string
}

type sequenceGenerator struct {
	baseGenerator
	item string
}

type groupGenerator struct {
	baseGenerator
	items []string
}

type unionGenerator struct {
	baseGenerator
	union []string
}

var (
	isSep      func(node) bool
	postParse  = make(map[string]func(node) node)
	generators = make(map[string]generator)
	zeroNode   = node{}
)

func stringsContain(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}

func uniq(strs []string) []string {
	var strsu []string
	m := make(map[string]struct{})
	for _, s := range strs {
		if _, ok := m[s]; !ok {
			strsu = append(strsu, s)
			m[s] = struct{}{}
		}
	}

	return strsu
}

func setPostParse(p map[string]func(node) node) {
	for pi, pp := range p {
		postParse[pi] = pp
	}
}

func (n node) zero() bool { return n.typ == "" }

func (p *baseParser) path() []string { return p.p }

func (p *baseParser) name() string {
	path := p.path()
	if len(path) == 0 {
		return "empty-parser"
	}

	return path[len(path)-1]
}

func (p *baseParser) out(args ...interface{}) {
	log.Println(
		append(
			[]interface{}{strings.Join(p.path(), "/")},
			args...,
		)...,
	)
}

func (g *baseGenerator) name() string { return g.node }

func (g *baseGenerator) out(args ...interface{}) {
	log.Println(
		append(
			[]interface{}{g.name()},
			args...,
		)...,
	)
}

func newPrimitiveParser(path []string, name string, token tokenType, init node) *primitiveParser {
	p := &primitiveParser{}

	p.p = append(path, name)
	if init.zero() {
		p.accepting = true
		p.token = token
	} else {
		// p.out("initialized with node")
		p.node = init
	}

	return p
}

func (p *primitiveParser) parse(t token) parseResult {
	// p.out("parse", t)
	if !p.accepting || t.typ != p.token {
		// p.out("not accepting", p.node.typ == p.name())
		return parseResult{
			accepting: false,
			valid:     p.node.typ == p.name(),
			unparsed:  []token{t},
			node:      p.node,
		}
	}

	// p.out("accepting")
	p.node = node{typ: p.name(), token: t}
	p.accepting = false
	return parseResult{accepting: true}
}

func primitive(name string, token tokenType) {
	g := &primitiveGenerator{}
	g.node = name
	g.token = token
	generators[name] = g
}

func (g *primitiveGenerator) name() string { return g.node }

func (g *primitiveGenerator) canCreate(init node, excluded []string) bool {
	if stringsContain(excluded, g.name()) {
		return false
	}

	if !init.zero() && init.typ != g.name() {
		return false
	}

	return true
}

func (g *primitiveGenerator) create(path []string, init node, _ []string) parser {
	return newPrimitiveParser(path, g.name(), g.token, init)
}

func (g *primitiveGenerator) member(node string) bool {
	return node == g.name()
}

func newOptionalParser(path []string, name string, optional generator, init node, excluded []string) *optionalParser {
	p := &optionalParser{}
	p.optional = optional
	p.init = init
	p.excluded = append(excluded, name)
	p.p = append(path, name)
	return p
}

func (p *optionalParser) parse(t token) parseResult {
	// p.out("parsing", t)
	if p.parser == nil {
		p.parser = p.optional.create(p.path(), p.init, p.excluded)
	}

	result := p.parser.parse(t)
	if !result.accepting {
		result.valid = true
	}

	return result
}

func optional(name, optional string) {
	g := &optionalGenerator{}
	g.node = name
	g.optional = optional
	generators[name] = g
}

func (g *optionalGenerator) name() string { return g.node }

func (g *optionalGenerator) canCreate(init node, excluded []string) bool {
	optional, ok := generators[g.optional]
	if !ok {
		panic("generator not found:" + g.optional)
	}

	if stringsContain(excluded, g.name()) {
		return false
	}

	return optional.canCreate(init, excluded)
}

func (g *optionalGenerator) create(path []string, init node, excluded []string) parser {
	return newOptionalParser(path, g.name(), generators[g.optional], init, excluded)
}

func (g *optionalGenerator) member(node string) bool {
	return node == g.name() || node == g.optional
}

func newSequenceParser(path []string, name string, itemGenerator generator, init node, excluded []string) *sequenceParser {
	p := &sequenceParser{}
	p.node = node{typ: name}
	p.init = init
	p.itemGenerator = itemGenerator
	p.p = append(path, name)
	p.excluded = excluded
	return p
}

func (p *sequenceParser) parse(t token) parseResult {
	// p.out("parsing", t)
	if p.currentParser == nil {
		if len(p.node.nodes) == 0 {
			p.currentParser = p.itemGenerator.create(p.path(), p.init, p.excluded)
		} else {
			p.currentParser = p.itemGenerator.create(p.path(), zeroNode, []string{p.name()})
		}
	}

	itemResult := p.currentParser.parse(t)
	if itemResult.accepting {
		if len(p.queue) > 0 {
			t, p.queue = p.queue[0], p.queue[1:]
			// p.out("accepting from queue")
			return p.parse(t)
		}

		// p.out("accepting")
		return parseResult{accepting: true}
	}

	// p.out("item parse done")

	if !itemResult.valid {
		// p.out("invalid")
		return parseResult{
			valid: true,
			node:  p.node,
			unparsed: append(
				itemResult.unparsed,
				p.queue...,
			),
		}
	}

	// p.out("valid")

	if !itemResult.node.zero() {
		// p.out("has node")
		if len(p.node.nodes) == 0 {
			p.node.token = itemResult.node.token
		}

		// p.out("appending node", itemResult.node.typ)
		p.node.nodes = append(p.node.nodes, itemResult.node)
	}

	p.currentParser = nil
	p.queue = append(itemResult.unparsed, p.queue...)
	t, p.queue = p.queue[0], p.queue[1:]

	// p.out("next from queue")
	return p.parse(t)
}

func sequence(name string, item string) {
	g := &sequenceGenerator{}
	g.node = name
	g.item = item
	generators[name] = g
}

func (g *sequenceGenerator) name() string { return g.node }

func (g *sequenceGenerator) canCreate(init node, excluded []string) bool {
	gen, ok := generators[g.item]
	if !ok {
		panic("generator not found: " + g.item)
	}

	if stringsContain(excluded, g.name()) {
		return false
	}

	return gen.canCreate(init, append(excluded, g.name()))
}

func (g *sequenceGenerator) create(path []string, init node, excluded []string) parser {
	return newSequenceParser(path, g.node, generators[g.item], init, append(excluded, g.name()))
}

func (g *sequenceGenerator) member(node string) bool {
	return node == g.name()
}

func newGroupParser(path []string, name string, generators []generator, init node, excluded []string) *groupParser {
	p := &groupParser{}
	p.node = node{typ: name}
	p.init = init
	p.generators = generators
	p.p = append(path, name)
	p.excluded = excluded
	if !p.init.zero() {
		// p.out("initialized with node")
	}

	return p
}

func (p *groupParser) parse(t token) parseResult {
	// p.out("parsing", t, p.queue)
	if p.currentParser == nil {
		if len(p.generators) == 0 {
			// p.out("done")
			// p.out("returning", append([]token{t}, p.queue...))
			return parseResult{
				valid:    true,
				node:     p.node,
				unparsed: append([]token{t}, p.queue...),
			}
		}

		if len(p.node.nodes) == 0 {
			p.currentParser = p.generators[0].create(
				p.path(),
				p.init,
				p.excluded,
			)
		} else {
			p.currentParser = p.generators[0].create(p.path(), zeroNode, nil)
		}

		p.generators = p.generators[1:]
	}

	itemResult := p.currentParser.parse(t)

	if itemResult.accepting {
		p.itemAccepted = append(p.itemAccepted, t)
		if len(p.queue) > 0 {
			t, p.queue = p.queue[0], p.queue[1:]
			// p.out("accepting from queue")
			// p.out("same item, accepted", len(p.itemAccepted), len(p.accepted))
			return p.parse(t)
		}

		// p.out("accepting")
		return parseResult{accepting: true}
	}

	if !itemResult.valid {
		p.itemAccepted = nil
		if len(p.node.nodes) == 0 && !p.init.zero() &&
			generators[p.currentParser.name()].member(p.init.typ) {

			// p.out("init item as node")
			p.node.token = p.init.token
			p.node.nodes = append(p.node.nodes, p.init)
			p.currentParser = nil
			p.queue = append(itemResult.unparsed, p.queue...)
			t, p.queue = p.queue[0], p.queue[1:]
			// p.out("invalid, accepted", len(p.itemAccepted), len(p.accepted))
			return p.parse(t)
		}

		// p.out("invalid")
		// p.out(
		// 	"returning rather",
		// 	p.accepted,
		// 	itemResult.unparsed,
		// 	p.queue,
		// )
		return parseResult{
			unparsed: append(
				p.accepted,
				append(
					itemResult.unparsed,
					p.queue...,
				)...,
			),
		}
	}

	if !itemResult.node.zero() {
		if len(p.node.nodes) == 0 {
			p.node.token = itemResult.node.token
		}

		p.node.nodes = append(p.node.nodes, itemResult.node)
	}

	p.itemAccepted = p.itemAccepted[0 : len(p.itemAccepted)-len(itemResult.unparsed)+1]

	p.currentParser = nil
	// p.out(
	// 	"adding to accepted",
	// 	p.accepted,
	// 	p.itemAccepted,
	// 	itemResult.valid,
	// 	itemResult.node.zero(),
	// 	itemResult.unparsed,
	// )
	p.accepted = append(p.accepted, p.itemAccepted...)
	p.itemAccepted = nil
	p.queue = append(itemResult.unparsed, p.queue...)
	t, p.queue = p.queue[0], p.queue[1:]

	// p.out("next from queue")
	// p.out("valid, accepted", len(p.itemAccepted), len(p.accepted))
	return p.parse(t)
}

func group(name string, items ...string) {
	g := &groupGenerator{}
	g.node = name
	g.items = items
	generators[name] = g
}

func (g *groupGenerator) name() string { return g.node }

func (g *groupGenerator) canCreate(init node, excluded []string) bool {
	if stringsContain(excluded, g.name()) {
		return false
	}

	for _, gi := range g.items {
		if _, ok := generators[gi]; !ok {
			panic("generator not found: " + gi)
		}
	}

	if len(g.items) == 0 {
		return false
	}

	first := g.items[0]
	if generators[first].canCreate(init, append(excluded, g.name())) {
		return true
	}

	if !init.zero() && generators[first].member(init.typ) {
		return true
	}

	return false
}

func (g *groupGenerator) create(path []string, init node, excluded []string) parser {
	gens := make([]generator, len(g.items))
	for i, item := range g.items {
		gens[i] = generators[item]
	}

	return newGroupParser(path, g.node, gens, init, append(excluded, g.name()))
}

func (g *groupGenerator) member(node string) bool {
	return node == g.name()
}

func newUnionParser(path []string, name string, init node, generators []generator, excluded []string) *unionParser {
	p := &unionParser{}
	p.p = append(path, name)
	p.node = init
	p.generators = generators
	p.activeGenerators = generators
	p.excluded = append(excluded, name)

	gs := make([]string, len(p.generators))
	for i, gi := range p.generators {
		gs[i] = gi.name()
	}

	// p.out("created", name, gs, p.excluded)
	return p
}

func (p *unionParser) parse(t token) parseResult {
	// p.out("parsing", t)
	if p.currentParser == nil {
		// p.out("excluded", p.excluded)
		for {
			if len(p.activeGenerators) == 0 {
				// p.out("finished union, valid:", p.valid)
				return parseResult{
					node:  p.node,
					valid: p.valid,
					unparsed: append(
						[]token{t},
						p.queue...,
					),
				}
			}

			var g generator
			g, p.activeGenerators = p.activeGenerators[0], p.activeGenerators[1:]
			// p.out("looking for generator", g.name())
			if g.canCreate(p.node, p.excluded) {
				p.currentParser = g.create(p.path(), p.node, p.excluded)
				break
			}
		}
	}

	// p.out("call to parse")
	elementResult := p.currentParser.parse(t)

	if elementResult.accepting {
		p.hasAccepted = true
		// p.out("accepting")
		if len(p.queue) > 0 {
			// p.out("from queue", p.queue)
			t, p.queue = p.queue[0], p.queue[1:]
			// p.out("queue set after accept", p.queue)
			return p.parse(t)
		}

		return parseResult{accepting: true}
	}

	// p.out("element parse done")

	p.currentParser = nil

	if !elementResult.valid {
		// p.out("invalid union parse", p.valid, elementResult.unparsed, p.queue)
		p.queue = append(elementResult.unparsed, p.queue...)
		// p.out("queue set after invalid", p.queue)
		if len(p.queue) > 0 {
			t, p.queue = p.queue[0], p.queue[1:]
			// p.out("queue set after taken on invalid", p.queue)
			return p.parse(t)
		}

		return parseResult{accepting: true}
	}

	// p.out("valid")

	if !p.valid || p.hasAccepted {
		// p.out("setting valid")
		p.valid = true
		p.node = elementResult.node
		p.activeGenerators = p.generators
		p.hasAccepted = false
	}

	// p.out("a valid union parse", p.valid, elementResult.unparsed, p.queue)
	p.queue = append(elementResult.unparsed, p.queue...)
	// p.out("queue set after valid", p.queue)
	if len(p.queue) == 0 {
		// p.out("next from outside")
		return parseResult{accepting: true}
	}

	t, p.queue = p.queue[0], p.queue[1:]
	// p.out("queue set after taken on valid", p.queue)
	// p.out("next from queue")
	return p.parse(t)
}

func union(node string, union ...string) {
	g := &unionGenerator{}
	g.node = node
	g.union = union
	generators[node] = g
}

func (g *unionGenerator) name() string { return g.node }

func (g *unionGenerator) expand(path []string) []string {
	if stringsContain(path, g.name()) {
		panic("union expansion loop")
	}

	var expanded []string
	for _, name := range g.union {
		gi, ok := generators[name]
		if !ok {
			panic("generator not found")
		}

		if u, ok := gi.(*unionGenerator); ok {
			expanded = append(expanded, u.expand(append(path, g.name()))...)
		} else {
			expanded = append(expanded, name)
		}
	}

	return uniq(expanded)
}

func (g *unionGenerator) canCreate(init node, excluded []string) bool {
	if len(g.union) == 0 {
		return false
	}

	expanded := g.expand(nil)
	for _, element := range expanded {
		if generators[element].canCreate(init, excluded) {
			return true
		}
	}

	return false
}

func (g *unionGenerator) create(path []string, init node, excluded []string) parser {
	expanded := g.expand(nil)

	var gens []generator
	for _, element := range expanded {
		gen := generators[element]
		if gen.canCreate(init, excluded) {
			gens = append(gens, gen)
		}
	}

	return newUnionParser(path, g.node, init, gens, excluded)
}

func (g *unionGenerator) member(node string) bool {
	expanded := g.expand(nil)
	for _, gi := range expanded {
		if generators[gi].member(node) {
			return true
		}
	}

	return false
}

func dropSeps(n []node) []node {
	if isSep == nil {
		return n
	}

	nn := make([]node, 0, len(n))
	for _, ni := range n {
		if !isSep(ni) {
			nn = append(nn, ni)
		}
	}

	return nn
}

func postParseNode(n node) node {
	n.nodes = postParseNodes(n.nodes)
	if pp, ok := postParse[n.typ]; ok {
		n = pp(n)
	}

	return n
}

func postParseNodes(n []node) []node {
	n = dropSeps(n)
	for i, ni := range n {
		n[i] = postParseNode(ni)
	}

	return n
}

func parse(p generator, r *tokenReader) (node, error) {
	gi := p.create(nil, zeroNode, nil)
	for {
		t, err := r.next()
		if err != nil && err != io.EOF {
			return node{}, err
		}

		if err == io.EOF {
			// gi.out("accepting after eof", token{})
			result := gi.parse(token{})
			if len(result.unparsed) != 1 && result.unparsed[0].typ != noToken {
				// println("unparsed length", len(result.unparsed))
				if len(result.unparsed) > 0 {
					// println(result.unparsed[0].value)
				}
				return node{}, errors.New("unexpected EOF")
			}

			// println("root post-parsing", result.node.typ, len(result.node.nodes))
			// for _, ni := range result.node.nodes {
			// 	println(ni.typ)
			// }

			n := postParseNode(result.node)
			// println("root returning", n.typ, len(n.nodes))
			return n, nil
		}

		// println("root accepting", t.value, t.value == "", len(t.value))
		// gi.out("accepting", t, t.value == "")
		result := gi.parse(t)
		if len(result.unparsed) > 0 {
			// println("unparsed", len(result.unparsed), t.value)
			// for _, up := range result.unparsed {
			// 	println(up.value)
			// }
			return node{}, fmt.Errorf("unexpected token: %v", t)
		}
	}
}
