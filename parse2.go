package mml

import (
	"fmt"
	"io"
	"log"
)

type node struct {
	token
	typ   string
	nodes []node
}

type parseError struct {
	nodeType  string
	tokenType tokenType
	token     token
}

type parserInstance interface {
	accept(token) bool
	valid() bool
	node() node
	error() error
}

type parser struct {
	name    string
	create  func([]node, []string) (parserInstance, bool)
	members []string
}

type parserBase struct {
	v bool
	n node
	e error
}

type primitiveParser struct {
	parserBase
}

type optionalParser struct {
	parserBase
	optional *parser
	instance parserInstance
}

type sequenceParser struct {
	parserBase
	itemParser *parser
	current    parserInstance
	excludes   []string
}

type groupParser struct {
	parserBase
	current  parserInstance
	parsers  []*parser
	excludes []string
}

type unionParser struct {
	parserBase
	name     string
	current  parserInstance
	parsers  []*parser
	excludes []string
	accepted []token
	queue    []token
}

var (
	isSep     func(node) bool
	postParse = make(map[string]func(node) node)
	parsers   = make(map[string]*parser)
)

func perror(nodeType string, tt tokenType, t token) error {
	return &parseError{nodeType: nodeType, tokenType: tt, token: t}
}

func (pe *parseError) Error() string {
	format := "%s:%d:%d: error when parsing %v, unexpected %v"
	args := []interface{}{
		pe.token.fileName, pe.token.line, pe.token.column, pe.nodeType, pe.token,
	}

	if pe.tokenType != noToken {
		format += ", expecting: %v"
		args = append(args, pe.tokenType)
	}

	return fmt.Sprintf(format, args...)
}

func (p *parser) member(m string) bool {
	for _, mi := range p.members {
		if mi == m {
			return true
		}
	}

	return false
}

func setPostParse(p map[string]func(node) node) {
	for pi, pp := range p {
		postParse[pi] = pp
	}
}

func (p *parser) String() string {
	return p.name
}

func (p *parserBase) valid() bool  { return p.v }
func (p *parserBase) node() node   { return p.n }
func (p *parserBase) error() error { return p.e }

func newPrimitiveParser(nodeType string, t tokenType) *primitiveParser {
	p := &primitiveParser{}
	p.n = node{typ: nodeType, token: token{typ: t}}
	return p
}

func (p *primitiveParser) accept(t token) bool {
	if p.v || p.e != nil {
		return false
	}

	if t.typ == p.n.token.typ {
		p.n.token = t
		p.v = true
	} else {
		p.e = perror(p.n.typ, p.n.token.typ, t)
		p.v = false
	}

	return p.v
}

func newOptionalParser(optional *parser, excludes []string) *optionalParser {
	p := &optionalParser{}
	p.n = node{}
	p.instance, _ = optional.create(nil, excludes)
	return p
}

func (p *optionalParser) accept(t token) bool {
	if p.instance.accept(t) {
		if p.n.typ == "" {
			p.n.typ = p.instance.node().typ
		}

		return true
	}

	if p.n.typ == "" {
		p.v = true
		return false
	}

	p.v = p.instance.valid()
	p.n = p.instance.node()
	p.e = p.instance.error()
	return false
}

func newSequenceParser(nodeType string, itemParser *parser, excludes []string) *sequenceParser {
	p := &sequenceParser{excludes: excludes}
	p.n = node{typ: nodeType}
	p.itemParser = itemParser
	return p
}

func (p *sequenceParser) accept(t token) bool {
	current := p.current
	if current == nil {
		var excludes []string
		if len(p.n.nodes) == 0 {
			excludes = p.excludes
		}

		current, _ = p.itemParser.create(nil, excludes)
	}

	if current.accept(t) {
		p.current = current
		return true
	}

	if p.current == nil {
		p.v = true
		return false
	}

	if !current.valid() {
		p.e = current.error()
		return false
	}

	n := current.node()
	p.n.nodes = append(p.n.nodes, n)
	if p.n.token.typ == noToken {
		p.n.token = n.token
	}

	p.current = nil
	return p.accept(t)
}

func newGroupParser(nodeType string, parsers []*parser, excludes []string) *groupParser {
	p := &groupParser{excludes: excludes}
	p.n = node{typ: nodeType}
	p.parsers = parsers
	return p
}

func (p *groupParser) accept(t token) bool {
	if p.current == nil {
		if len(p.parsers) == 0 {
			p.v = true
			return false
		}

		var excludes []string
		if len(p.n.nodes) == 0 {
			excludes = p.excludes
		}

		p.current, _ = p.parsers[0].create(nil, excludes)
		p.parsers = p.parsers[1:]
	}

	if p.current.accept(t) {
		return true
	}

	if !p.current.valid() {
		p.e = p.current.error()
		return false
	}

	n := p.current.node()
	if n.typ != "" {
		p.n.nodes = append(p.n.nodes, n)
		if p.n.token.typ == noToken {
			p.n.token = n.token
		}
	}

	p.current = nil
	return p.accept(t)
}

func newUnionParser(parsers []*parser, excludes []string) *unionParser {
	p := &unionParser{excludes: excludes}
	p.parsers = parsers
	return p
}

func (p *unionParser) accept(t token) bool {
	// need a token queue

	out := func(a ...interface{}) {
		return
		typ := p.n.typ
		if typ == "" {
			typ = "unset node"
		}

		current := "unset parser"
		if p.current != nil {
			current = p.current.node().typ
		}

		log.Println(append([]interface{}{"union:", p.name, typ, current}, a...)...)
	}

	out("accepting", t.value)

	if p.current == nil {
		out("finding parser")
		for {
			if len(p.parsers) == 0 {
				out("no parser found")
				return false
			}

			var init []node
			if p.n.typ != "" {
				init = []node{p.n}
			}

			if exclude(p.excludes, p.parsers[0].name) {
				p.parsers = p.parsers[1:]
			} else {
				cp := p.parsers[0]
				excludes := append(p.excludes, cp.name)
				current, ok := cp.create(init, excludes)
				p.parsers = p.parsers[1:]

				if ok {
					p.current = current
					out("found next parser", cp.name)
					break
				}
			}
		}
	}

	out("calling current accept")
	if p.current.accept(t) {
		out("accepted")
		p.accepted = append(p.accepted, t)
		if len(p.queue) > 0 {
			t, p.queue = p.queue[0], p.queue[1:]
			return p.accept(t)
		}

		return true
	}

	if !p.current.valid() {
		out("not valid")
		if p.n.typ == "" {
			p.v = false
			p.e = p.current.error()
		}

		p.current = nil
		out("trying with next parser")
		if len(p.accepted) > 0 {
			out("has accepted")
			t, p.queue, p.accepted = p.accepted[0], append(p.accepted[1:], t), nil
		}

		return p.accept(t)
	}

	out("valid")

	p.n = p.current.node()
	p.v = true
	p.current = nil
	p.parsers = append([]*parser{parsers[p.n.typ]}, p.parsers...)

	out("trying with the same set of parsers")
	return p.accept(t)
}

func exclude(excludes []string, typ string) bool {
	for _, e := range excludes {
		if e == typ {
			return true
		}
	}

	return false
}

func primitive(name string, t tokenType) {
	p := &parser{name: name}
	p.create = func(init []node, _ []string) (parserInstance, bool) {
		if len(init) > 0 {
			return nil, false
		}

		return newPrimitiveParser(name, t), true
	}

	parsers[name] = p
}

func optional(name, typ string) {
	p := &parser{name: name}
	p.create = func(init []node, excludes []string) (parserInstance, bool) {
		if len(init) > 0 {
			return nil, false
		}

		return newOptionalParser(parsers[typ], excludes), true
	}

	parsers[name] = p
}

func sequence(name, typ string) {
	p := &parser{name: name}
	p.create = func(init []node, excludes []string) (parserInstance, bool) {
		itemParser := parsers[typ]
		for _, i := range init {
			if !itemParser.member(i.typ) || exclude(excludes, i.typ) {
				return nil, false
			}
		}

		pi := newSequenceParser(name, itemParser, excludes)
		pi.n.nodes = init
		return pi, true
	}

	parsers[name] = p
}

func group(name string, types ...string) {
	p := &parser{name: name}
	p.create = func(init []node, excludes []string) (parserInstance, bool) {
		types := types[:]
		if len(init) > len(types) {
			return nil, false
		}

		for _, i := range init {
			itemParser := parsers[types[0]]
			if !itemParser.member(i.typ) {
				return nil, false
			}

			types = types[1:]
		}

		itemParsers := make([]*parser, len(types))
		for i, t := range types {
			itemParsers[i] = parsers[t]
		}

		pi := newGroupParser(name, itemParsers, excludes)
		if len(init) > 0 {
			pi.n.nodes = init
			pi.n.token = init[0].token
		}

		return pi, true
	}

	parsers[name] = p
}

func union(name string, types ...string) {
	p := &parser{name: name, members: types}
	p.create = func(init []node, excludes []string) (parserInstance, bool) {
		if len(init) > 0 {
			return nil, false
		}

		itemParsers := make([]*parser, len(types))
		for i, t := range types {
			itemParsers[i] = parsers[t]
		}

		pi := newUnionParser(itemParsers, excludes)
		pi.name = name
		return pi, true
	}

	parsers[name] = p
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

func parse(p *parser, r *tokenReader) (node, error) {
	pi, _ := p.create(nil, nil)
	for {
		t, err := r.next()
		if err != nil && err != io.EOF {
			return pi.node(), err
		}

		if err == io.EOF {
			pi.accept(token{})
			n := pi.node()
			if len(n.nodes) > 1 {
			}
			n = postParseNode(n)
			if len(n.nodes) > 1 {
			}
			return n, nil
		}

		if !pi.accept(t) {
			return pi.node(), perror(p.String(), noToken, t)
		}
	}
}
