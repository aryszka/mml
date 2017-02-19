package mml

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"strings"
// )
//
// type node struct {
// 	token
// 	typ   string
// 	nodes []node
// }
//
// type parseError struct {
// 	nodeType  string
// 	tokenType tokenType
// 	token     token
// }
//
// type parserInstance interface {
// 	accept(token) bool
// 	valid() bool
// 	node() node
// 	error() error
// }
//
// type parser struct {
// 	name    string
// 	create  func([]string, []node, []string) (parserInstance, bool)
// 	members []string
// }
//
// type parserBase struct {
// 	v bool
// 	n node
// 	e error
// 	name string
// 	path []string
// }
//
// type primitiveParser struct {
// 	parserBase
// }
//
// type optionalParser struct {
// 	parserBase
// 	optional *parser
// 	instance parserInstance
// }
//
// type sequenceParser struct {
// 	parserBase
// 	itemParser *parser
// 	current    parserInstance
// 	excludes   []string
// }
//
// type groupParser struct {
// 	parserBase
// 	name     string
// 	current  parserInstance
// 	parsers  []*parser
// 	excludes []string
// }
//
// type unionParser struct {
// 	parserBase
// 	name     string
// 	current  parserInstance
// 	parsers  []*parser
// 	excludes []string
// 	accepted []token
// 	queue    []token
// }
//
// var (
// 	isSep     func(node) bool
// 	postParse = make(map[string]func(node) node)
// 	parsers   = make(map[string]*parser)
// )
//
// func perror(nodeType string, tt tokenType, t token) error {
// 	return &parseError{nodeType: nodeType, tokenType: tt, token: t}
// }
//
// func (pe *parseError) Error() string {
// 	format := "%s:%d:%d: error when parsing %v, unexpected %v"
// 	args := []interface{}{
// 		pe.token.fileName, pe.token.line, pe.token.column, pe.nodeType, pe.token,
// 	}
//
// 	if pe.tokenType != noToken {
// 		format += ", expecting: %v"
// 		args = append(args, pe.tokenType)
// 	}
//
// 	return fmt.Sprintf(format, args...)
// }
//
// func (p *parser) member(m string) bool {
// 	for _, mi := range p.members {
// 		if mi == m {
// 			return true
// 		}
// 	}
//
// 	return false
// }
//
// func setPostParse(p map[string]func(node) node) {
// 	for pi, pp := range p {
// 		postParse[pi] = pp
// 	}
// }
//
// func (p *parser) String() string {
// 	return p.name
// }
//
// func (p *parserBase) valid() bool  { return p.v }
// func (p *parserBase) node() node   { return p.n }
// func (p *parserBase) error() error { return p.e }
//
// func (p *parserBase) out(args ...interface{}) {
// 	return
// 	log.Println(append([]interface{}{strings.Join(p.path, ":")}, args...)...)
// }
//
// func newPrimitiveParser(path []string, nodeType string, t tokenType) *primitiveParser {
// 	p := &primitiveParser{}
// 	p.n = node{typ: nodeType, token: token{typ: t}}
// 	p.path = append(path, nodeType)
// 	return p
// }
//
// func (p *primitiveParser) accept(t token) bool {
// 	p.out("accepting", t.value)
// 	if p.v || p.e != nil {
// 		return false
// 	}
//
// 	if t.typ == p.n.token.typ {
// 		p.n.token = t
// 		p.v = true
// 	} else {
// 		p.e = perror(p.n.typ, p.n.token.typ, t)
// 		p.v = false
// 	}
//
// 	return p.v
// }
//
// func newOptionalParser(path []string, name string, optional *parser, excludes []string) *optionalParser {
// 	p := &optionalParser{}
// 	p.n = node{}
// 	p.path = append(path, name)
// 	p.instance, _ = optional.create(p.path, nil, excludes)
// 	return p
// }
//
// func (p *optionalParser) accept(t token) bool {
// 	p.out("accepting", t.value)
// 	if p.instance.accept(t) {
// 		if p.n.typ == "" {
// 			p.n.typ = p.instance.node().typ
// 		}
//
// 		return true
// 	}
//
// 	if p.n.typ == "" {
// 		p.v = true
// 		return false
// 	}
//
// 	p.v = p.instance.valid()
// 	p.n = p.instance.node()
// 	p.e = p.instance.error()
// 	return false
// }
//
// func newSequenceParser(path []string, nodeType string, itemParser *parser, excludes []string) *sequenceParser {
// 	p := &sequenceParser{excludes: excludes}
// 	p.n = node{typ: nodeType}
// 	p.path = append(path, nodeType)
// 	p.itemParser = itemParser
// 	return p
// }
//
// func (p *sequenceParser) accept(t token) bool {
// 	p.out("accepting", t.value)
// 	current := p.current
// 	if current == nil {
// 		var excludes []string
// 		if len(p.n.nodes) == 0 {
// 			excludes = p.excludes
// 		}
//
// 		current, _ = p.itemParser.create(p.path, nil, excludes)
// 	}
//
// 	if current.accept(t) {
// 		p.current = current
// 		return true
// 	}
//
// 	if p.current == nil {
// 		p.v = true
// 		return false
// 	}
//
// 	if !current.valid() {
// 		p.e = current.error()
// 		return false
// 	}
//
// 	n := current.node()
// 	p.n.nodes = append(p.n.nodes, n)
// 	if p.n.token.typ == noToken {
// 		p.n.token = n.token
// 	}
//
// 	p.current = nil
// 	return p.accept(t)
// }
//
// func newGroupParser(path []string, nodeType string, parsers []*parser, excludes []string) *groupParser {
// 	p := &groupParser{excludes: excludes}
// 	p.n = node{typ: nodeType}
// 	p.path = append(path, nodeType)
// 	p.parsers = parsers
// 	return p
// }
//
// func (p *groupParser) accept(t token) bool {
// 	p.out("accepting", t.value)
// 	if p.current == nil {
// 		if len(p.parsers) == 0 {
// 			p.v = true
// 			return false
// 		}
//
// 		var excludes []string
// 		if len(p.n.nodes) == 0 {
// 			excludes = p.excludes
// 		}
//
// 		p.current, _ = p.parsers[0].create(p.path, nil, excludes)
// 		p.parsers = p.parsers[1:]
// 	}
//
// 	if p.current.accept(t) {
// 		return true
// 	}
//
// 	if !p.current.valid() {
// 		p.e = p.current.error()
// 		return false
// 	}
//
// 	n := p.current.node()
// 	if n.typ != "" {
// 		p.n.nodes = append(p.n.nodes, n)
// 		if p.n.token.typ == noToken {
// 			p.n.token = n.token
// 		}
// 	}
//
// 	p.current = nil
// 	return p.accept(t)
// }
//
// func newUnionParser(path []string, name string, parsers []*parser, excludes []string) *unionParser {
// 	p := &unionParser{excludes: excludes}
// 	p.parsers = parsers
// 	p.path = append(path, name)
// 	return p
// }
//
// func (p *unionParser) accept(t token) bool {
// 	p.out("accepting", t.value)
// 	// take current
// 	if p.current == nil {
// 		for {
// 			// if still some to try
// 			if len(p.parsers) == 0 {
// 				return false
// 			}
//
// 			// initial node for group of sequence
// 			var init []node
// 			if p.n.typ != "" {
// 				init = []node{p.n}
// 			}
//
// 			// skip if excluded
// 			if exclude(p.excludes, p.parsers[0].name) {
// 				p.parsers = p.parsers[1:]
// 			} else {
// 				cp := p.parsers[0]
// 				// extend excludes with current
// 				excludes := append(p.excludes, cp.name)
//
// 				// find ok parser
// 				current, ok := cp.create(p.path, init, excludes)
// 				p.parsers = p.parsers[1:]
//
// 				if ok {
// 					p.current = current
// 					break
// 				}
// 			}
// 		}
// 	}
//
// 	// try current
// 	if p.current.accept(t) {
// 		// save accepted token
// 		p.accepted = append(p.accepted, t)
//
// 		// if has saved tokens, use them for the next
// 		if len(p.queue) > 0 {
// 			t, p.queue = p.queue[0], p.queue[1:]
// 			return p.accept(t)
// 		}
//
// 		return true
// 	}
//
// 	if !p.current.valid() {
// 		if p.n.typ == "" {
// 			p.v = false
// 			p.e = p.current.error()
// 		}
//
// 		// enqueue accepted
// 		p.current = nil
// 		if len(p.accepted) > 0 {
// 			t, p.queue, p.accepted = p.accepted[0], append(p.accepted[1:], t), nil
// 		}
//
// 		// if not valid, try with next parser
// 		return p.accept(t)
// 	}
//
// 	// if valid, store node and try next
// 	p.n = p.current.node()
// 	p.v = true
// 	p.current = nil
// 	// p.accepted, p.queue = p.queue, nil
//
// 	// put back the current into the accepted parsers
// 	p.parsers = append([]*parser{parsers[p.n.typ]}, p.parsers...)
//
// 	return p.accept(t)
// }
//
// func exclude(excludes []string, typ string) bool {
// 	for _, e := range excludes {
// 		if e == typ {
// 			return true
// 		}
// 	}
//
// 	return false
// }
//
// func primitive(name string, t tokenType) {
// 	p := &parser{name: name}
// 	p.create = func(path []string, init []node, _ []string) (parserInstance, bool) {
// 		if len(init) > 0 {
// 			return nil, false
// 		}
//
// 		return newPrimitiveParser(path, name, t), true
// 	}
//
// 	parsers[name] = p
// }
//
// func optional(name, typ string) {
// 	p := &parser{name: name}
// 	p.create = func(path []string, init []node, excludes []string) (parserInstance, bool) {
// 		if len(init) > 0 {
// 			return nil, false
// 		}
//
// 		return newOptionalParser(path, name, parsers[typ], excludes), true
// 	}
//
// 	parsers[name] = p
// }
//
// func sequence(name, typ string) {
// 	p := &parser{name: name}
// 	p.create = func(path []string, init []node, excludes []string) (parserInstance, bool) {
// 		itemParser := parsers[typ]
// 		for _, i := range init {
// 			if !itemParser.member(i.typ) || exclude(excludes, i.typ) {
// 				return nil, false
// 			}
// 		}
//
// 		pi := newSequenceParser(path, name, itemParser, excludes)
// 		pi.n.nodes = init
// 		return pi, true
// 	}
//
// 	parsers[name] = p
// }
//
// func group(name string, types ...string) {
// 	p := &parser{name: name}
// 	p.create = func(path []string, init []node, excludes []string) (parserInstance, bool) {
// 		types := types[:]
// 		if len(init) > len(types) {
// 			return nil, false
// 		}
//
// 		for _, i := range init {
// 			itemParser := parsers[types[0]]
// 			if !itemParser.member(i.typ) {
// 				return nil, false
// 			}
//
// 			types = types[1:]
// 		}
//
// 		itemParsers := make([]*parser, len(types))
// 		for i, t := range types {
// 			itemParsers[i] = parsers[t]
// 		}
//
// 		// log.Println("group parser created", name, len(itemParsers), len(excludes))
// 		pi := newGroupParser(path, name, itemParsers, excludes)
// 		if len(init) > 0 {
// 			pi.n.nodes = init
// 			pi.n.token = init[0].token
// 		}
//
// 		pi.name = name
// 		return pi, true
// 	}
//
// 	parsers[name] = p
// }
//
// func union(name string, types ...string) {
// 	p := &parser{name: name, members: types}
// 	p.create = func(path []string, init []node, excludes []string) (parserInstance, bool) {
// 		if len(init) > 1 {
// 			return nil, false
// 		}
//
// 		itemParsers := make([]*parser, len(types))
// 		for i, t := range types {
// 			itemParsers[i] = parsers[t]
// 		}
//
// 		pi := newUnionParser(path, name, itemParsers, excludes)
// 		if len(init) > 0 {
// 			pi.n = init[0]
// 		}
//
// 		pi.name = name
// 		return pi, true
// 	}
//
// 	parsers[name] = p
// }
//
// func dropSeps(n []node) []node {
// 	if isSep == nil {
// 		return n
// 	}
//
// 	nn := make([]node, 0, len(n))
// 	for _, ni := range n {
// 		if !isSep(ni) {
// 			nn = append(nn, ni)
// 		}
// 	}
//
// 	return nn
// }
//
// func postParseNode(n node) node {
// 	n.nodes = postParseNodes(n.nodes)
// 	if pp, ok := postParse[n.typ]; ok {
// 		n = pp(n)
// 	}
//
// 	return n
// }
//
// func postParseNodes(n []node) []node {
// 	n = dropSeps(n)
// 	for i, ni := range n {
// 		n[i] = postParseNode(ni)
// 	}
//
// 	return n
// }
//
// func parse(p *parser, r *tokenReader) (node, error) {
// 	pi, _ := p.create(nil, nil, nil)
// 	for {
// 		t, err := r.next()
// 		if err != nil && err != io.EOF {
// 			return pi.node(), err
// 		}
//
// 		if err == io.EOF {
// 			pi.accept(token{})
// 			n := pi.node()
// 			if len(n.nodes) > 1 {
// 			}
// 			n = postParseNode(n)
// 			if len(n.nodes) > 1 {
// 			}
// 			return n, nil
// 		}
//
// 		// log.Println("root accepting", t.value)
// 		if !pi.accept(t) {
// 			return pi.node(), perror(p.String(), noToken, t)
// 		}
// 	}
// }
