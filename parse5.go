package mml

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type traceLevel int

const (
	traceOff traceLevel = iota
	traceWarn
	traceOn
	traceDebug
)

type nodeType uint64

type typeList []nodeType

type node struct {
	token    *token
	nodeType nodeType
	typeName string
	nodes    []*node
	toks     []*token
}

type parserResult struct {
	accepting bool
	valid     bool
	node      *node
	fromCache bool
	unparsed  *tokenStack
}

type parser interface {
	reset()
	check(trace) error
	instance(trace, *node) parser // optimization: keep some predefined ones
	parse(*token) *parserResult   // optimization: should not make memory allocations
	valid() bool
	expectedLength() int
	nodeType() nodeType
}

type generator interface {
	create(trace, nodeType, typeList) (parser, error)
	member(nodeType) (bool, error)
	nodeType() nodeType
}

type trace interface {
	extend(nodeType) trace
	outLevel(traceLevel, ...interface{})
	out(...interface{})
	debug(...interface{})
	warn(...interface{})
}

type registry interface {
	typeName(nodeType) string
	nodeType(typeName string) nodeType
	get(nodeType) (generator, bool)
	getParser(nodeType, nodeType, typeList) (parser, bool)
	setParser(nodeType, nodeType, typeList, parser) // TODO: need to use this as a prototype
	init() error
	finalize(trace) error
	reset()
	rootGenerator() (generator, error) // TODO: temporary
	root() (parser, error)
	primitive(string, tokenType) error
	optional(string, string) error
	sequence(string, string) error
	group(string, ...string) error
	union(string, ...string) error
}

// TODO: there are appends/merges beyond the prealloacted

type tokenStack struct {
	trace      trace
	stack      []*token
	nextLength int
	need       int
	token      *token
	tokenIndex int
	skip       int
}

type parserTrace struct {
	registry registry
	level    traceLevel
	path     string
}

type parserRegistry struct {
	idSeed     nodeType
	typeIDs    map[string]nodeType
	typeNames  map[nodeType]string
	generators map[nodeType]generator
	parsers    map[string]parser
	rootGen    generator
}

type primitiveGenerator struct {
	typ       nodeType
	tokenType tokenType
	registry  registry
}

type primitiveParser struct {
	trace     trace
	registry  registry
	typ       nodeType
	typeName  string
	tokenType tokenType
	result    *parserResult
	isValid   bool
	length    int
}

type optionalGenerator struct {
	typ      nodeType
	optional nodeType
	registry registry
}

type optionalParser struct {
	trace            trace
	registry         registry
	typ              nodeType
	typeName         string
	optional         parser
	initIsMember     bool
	initNode         *node
	result           *parserResult
	optionalResult   *parserResult
	cacheToken       *token
	cacheChecked     bool
	isValid          bool
	length           int
	optionalInstance parser
}

type sequenceGenerator struct {
	typ      nodeType
	item     nodeType
	registry registry
}

type sequenceParser struct {
	trace             trace
	registry          registry
	typ               nodeType
	typeName          string
	first             parser
	rest              parser
	initIsMember      bool
	initNode          *node
	result            *parserResult
	skip              int
	currentParser     parser
	cacheChecked      bool
	cacheToken        *token
	itemResult        *parserResult
	tokenStack        *tokenStack
	initEvaluated     bool
	skippingAfterDone bool
	isValid           bool
	itemLength        int
}

type groupGenerator struct {
	typ      nodeType
	items    []nodeType
	registry registry
}

// TODO: verify that type name is used on all nodes

type groupParser struct {
	trace             trace
	registry          registry
	parsers           []parser
	initParsers       []parser
	currentParser     parser
	initIsMember      []bool
	result            *parserResult
	tokenStack        *tokenStack
	initNode          *node
	skip              int
	parserIndex       int
	cacheChecked      bool
	initEvaluated     bool
	skippingAfterDone bool
	cacheToken        *token
	itemResult        *parserResult
	isValid           bool
	length            int
	initType          nodeType
	typ               nodeType
	typeName          string
}

type unionGenerator struct {
	typ          nodeType
	registry     registry
	elements     []generator
	elementTypes []nodeType
}

type unionParser struct {
	trace             trace
	registry          registry
	typ               nodeType
	typeName          string
	tokenStack        *tokenStack
	parsers           [][]parser
	initIsMember      bool
	initNode          *node
	result            *parserResult
	initTypeIndex     int
	parserIndex       int
	skip              int
	skippingAfterDone bool
	cacheChecked      bool
	cacheToken        *token
	itemResult        *parserResult
	isValid           bool
	length            int
	currentParser     parser
}

type syntax struct {
	registry   registry
	traceLevel traceLevel
}

var (
	errNoParsersDefined     = errors.New("no parser defined")
	errFailedToCreateParser = errors.New("failed to create parser")
	errUnexpectedEOF        = errors.New("unexpected EOF")

	zeroNode = &node{}

	// TODO: can optimize with special handling of eof
	eofToken = &token{offset: -1, value: "<eof>", cache: newCache()}
)

func unexpectedInitNode(typeName, initTypeName string) error {
	return fmt.Errorf("unexpected init node: %s, %s", typeName, initTypeName)
}

func unspecifiedParser(typeName string) error {
	return fmt.Errorf("unspecified parser: %s", typeName)
}

func requiredParserInvalid(typeName string) error {
	return fmt.Errorf("required parser invalid: %s", typeName)
}

func duplicateNodeType(nodeType string) error {
	return fmt.Errorf("duplicate node type definition in syntax: %s", nodeType)
}

func unexpectedResult(nodeType string) error {
	return fmt.Errorf("unexpected parse result: %s", nodeType)
}

func optionalContainingSelf(nodeType string) error {
	return fmt.Errorf("optional containing self: %s", nodeType)
}

func sequenceContainingSelf(nodeType string) error {
	return fmt.Errorf("sequence containing self: %s", nodeType)
}

func itemParserCannotBeCreated(nodeType string) error {
	return fmt.Errorf("item parser cannot be created: %s", nodeType)
}

func groupWithoutItems(nodeType string) error {
	return fmt.Errorf("group without items: %s", nodeType)
}

func groupItemParserNotFound(nodeType string) error {
	return fmt.Errorf("group item parser not found: %s", nodeType)
}

func unionWithoutElements(nodeType string) error {
	return fmt.Errorf("union without elements: %s", nodeType)
}

func unexpectedToken(nodeType string, t *token) error {
	return fmt.Errorf("unexpected token: %v, %v", nodeType, t)
}

func (l typeList) contains(t nodeType) bool {
	for _, ti := range l {
		if ti == t {
			return true
		}
	}

	return false
}

func (n *node) tokens() []*token {
	return n.toks
}

func (n *node) len() int {
	return len(n.toks)
}

func (n *node) setToken(t *token) {
	// only for primitive:
	n.token = t
	n.toks[0] = t
}

func (n *node) append(na *node) {
	if na == zeroNode {
		return
	}

	n.nodes = append(n.nodes, na)
	n.toks = append(n.toks, na.tokens()...)
	if len(n.nodes) == 1 && len(n.toks) > 0 {
		n.token = n.toks[0]
	}
}

func (n *node) clearNodes() {
	n.nodes = nil
	n.toks = nil
}

func (n *node) String() string {
	nc := make([]string, len(n.nodes))
	for i, ni := range n.nodes {
		nc[i] = ni.String()
	}

	return fmt.Sprintf("{%s:%v:[%s]}", n.typeName, n.token, strings.Join(nc, ", "))
}

// TODO: fix the expected length

func newTokenStack(expectedSize int) *tokenStack {
	return &tokenStack{
		stack: make([]*token, 0, expectedSize),
	}
}

func (s *tokenStack) append(t *token) {
	if len(s.stack) == cap(s.stack) {
		s.trace.warn("token stack exceeded expected size on append", s.stack, t)
	}

	s.stack = append(s.stack, t)
}

func (s *tokenStack) merge(from *tokenStack) {
	s.need = len(s.stack) + len(from.stack) - cap(s.stack)
	if s.need > 0 {
		s.trace.warn("token stack exceeded expected size on merge", s.stack, from.stack)

		s.stack = s.stack[:cap(s.stack)]
		for s.need > 0 {
			s.stack = append(s.stack, nil)
			s.need--
		}
	} else {
		s.stack = s.stack[:len(s.stack)+len(from.stack)]
	}

	copy(s.stack[len(s.stack)-len(from.stack):], from.stack)
}

func (s *tokenStack) mergeTokens(t []*token) {
	for len(t) > 0 {
		s.append(t[len(t)-1])
		t = t[:len(t)-1]
	}
}

func (s *tokenStack) has() bool {
	return len(s.stack) > 0
}

func (s *tokenStack) peek() *token {
	return s.stack[len(s.stack)-1]
}

func (s *tokenStack) pop() *token {
	s.token, s.stack = s.stack[len(s.stack)-1], s.stack[:len(s.stack)-1]
	return s.token
}

func (s *tokenStack) drop(n int) {
	s.nextLength = len(s.stack) - n
	if s.nextLength < 0 {
		s.trace.warn("stack dropping more tokens than it contains")
		s.nextLength = 0
	}

	s.stack = s.stack[:s.nextLength]
}

func (s *tokenStack) clear() {
	s.drop(len(s.stack))
}

func (s *tokenStack) findCachedNode(n *node) int {
	for s.tokenIndex, s.token = range n.tokens() {
		if s.token != s.peek() {
			continue
		}

		if n.len()-s.tokenIndex > len(s.stack) {
			s.skip = n.len() - s.tokenIndex - len(s.stack)
			s.clear()
		} else {
			s.skip = 0
			s.drop(n.len() - s.tokenIndex)
		}

		return s.skip
	}

	return 0
}

func (s *tokenStack) setTrace(t trace) {
	s.trace = t
}

func newTrace(l traceLevel, r registry) *parserTrace {
	return &parserTrace{
		registry: r,
		level:    l,
	}
}

func (t *parserTrace) extend(n nodeType) trace {
	// TODO: save this operation of trace level is off
	return &parserTrace{
		registry: t.registry,
		level:    t.level,
		path:     t.path + "/" + t.registry.typeName(n),
	}
}

func (t *parserTrace) outLevel(l traceLevel, a ...interface{}) {
	if l > t.level {
		return
	}

	log.Println(append([]interface{}{t.path}, a...)...)
}

func (t *parserTrace) out(a ...interface{}) {
	t.outLevel(traceOn, a...)
}

func (t *parserTrace) warn(a ...interface{}) {
	t.outLevel(traceWarn, a...)
}

func (t *parserTrace) debug(a ...interface{}) {
	t.outLevel(traceDebug, a...)
}

// TODO: this cache doesn't help. If the init node matters, then it is worth looking into the cache during the
// instantiation, before all the allocations. The matches can be stored on the token in a list, while the
// non-matches can be stored as id sets also on the token. Maybe it is enough to cache the union elements. Or we
// can check in an id set whether something exists in the cache, and only look it up when it does. The lookup
// also can be in a balanced tree.

// TODO: check the cache before instantiation

func newRegistry() *parserRegistry {
	return &parserRegistry{
		idSeed:     1,
		typeIDs:    make(map[string]nodeType),
		typeNames:  make(map[nodeType]string),
		generators: make(map[nodeType]generator),
		parsers:    make(map[string]parser),
	}
}

func (r *parserRegistry) nodeType(typeName string) nodeType {
	t, ok := r.typeIDs[typeName]
	if ok {
		return t
	}

	t = r.idSeed
	r.idSeed++
	r.typeIDs[typeName] = t
	r.typeNames[t] = typeName
	// log.Println("registered type", typeName, t)
	return t
}

func (r *parserRegistry) typeName(t nodeType) string {
	return r.typeNames[t]
}

func (r *parserRegistry) get(t nodeType) (generator, bool) {
	g, ok := r.generators[t]
	return g, ok
}

func parserKey(t nodeType, init nodeType, excluded typeList) string {
	// or just use a hash?

	s := make([]string, len(excluded)+2)
	for i, ni := range append([]nodeType{t, init}, excluded...) {
		s[i] = fmt.Sprint(ni)
	}

	return strings.Join(s, "_")
}

func (r *parserRegistry) getParser(t nodeType, init nodeType, excluded typeList) (parser, bool) {
	p, ok := r.parsers[parserKey(t, init, excluded)]
	return p, ok
}

func (r *parserRegistry) setParser(t nodeType, init nodeType, excluded typeList, p parser) {
	r.parsers[parserKey(t, init, excluded)] = p
}

func (r *parserRegistry) init() error {
	return nil
}

func (r *parserRegistry) reset() {
	for _, p := range r.parsers {
		p.reset()
	}
}

func (r *parserRegistry) finalize(t trace) error {
	var done bool
	for !done {
		done = true
		for k, p := range r.parsers {
			if !p.valid() {
				delete(r.parsers, k)
				done = false
				continue
			}

			if err := p.check(t); err != nil {
				return err
			}

			if !p.valid() {
				delete(r.parsers, k)
				done = false
				continue
			}
		}
	}

	return nil
}

func (r *parserRegistry) rootGenerator() (generator, error) {
	if r.rootGenerator == nil {
		return nil, errNoParsersDefined
	}

	return r.rootGen, nil
}

func (r *parserRegistry) root() (parser, error) {
	return nil, nil
}

func (r *parserRegistry) register(nt nodeType, g generator) error {
	if _, exists := r.generators[nt]; exists {
		return duplicateNodeType(r.typeNames[nt])
	}

	r.generators[nt] = g
	r.rootGen = g // the last one is the root
	return nil
}

func (r *parserRegistry) primitive(typeName string, t tokenType) error {
	g := &primitiveGenerator{
		typ:       r.nodeType(typeName),
		tokenType: t,
		registry:  r,
	}

	return r.register(g.typ, g)
}

func (r *parserRegistry) optional(typeName string, optional string) error {
	g := &optionalGenerator{
		typ:      r.nodeType(typeName),
		optional: r.nodeType(optional),
		registry: r,
	}

	return r.register(g.typ, g)
}

func (r *parserRegistry) sequence(typeName string, itemType string) error {
	g := &sequenceGenerator{
		typ:      r.nodeType(typeName),
		item:     r.nodeType(itemType),
		registry: r,
	}

	return r.register(g.typ, g)
}

func (r *parserRegistry) group(typeName string, itemTypes ...string) error {
	items := make([]nodeType, len(itemTypes))
	for i, t := range itemTypes {
		items[i] = r.nodeType(t)
	}

	g := &groupGenerator{
		typ:      r.nodeType(typeName),
		items:    items,
		registry: r,
	}

	return r.register(g.typ, g)
}

func (r *parserRegistry) union(typeName string, elementTypes ...string) error {
	elements := make([]nodeType, len(elementTypes))
	for i, t := range elementTypes {
		elements[i] = r.nodeType(t)
	}

	g := &unionGenerator{
		typ:          r.nodeType(typeName),
		elementTypes: elements,
		registry:     r,
	}

	return r.register(g.typ, g)
}

func (g *primitiveGenerator) create(t trace, init nodeType, excluded typeList) (parser, error) {
	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &primitiveParser{
		typ:       g.typ,
		typeName:  g.registry.typeName(g.typ),
		registry:  g.registry,
		isValid:   !excluded.contains(g.typ) && init == 0,
		tokenType: g.tokenType,
		length:    1,
	}
	g.registry.setParser(g.typ, init, excluded, p)
	return p, nil
}

func (g *primitiveGenerator) member(t nodeType) (bool, error) {
	return t == g.typ, nil
}

func (g *primitiveGenerator) nodeType() nodeType {
	return g.typ
}

func (p *primitiveParser) check(trace) error {
	return nil
}

func (p *primitiveParser) reset() {
}

// TODO: try lazy node creation, or even better clone on init only when it was valid. Watch for references to
// the result.node for types or whatever

func (p *primitiveParser) instance(t trace, n *node) parser {
	if n != zeroNode {
		panic(unexpectedInitNode(p.registry.typeName(p.typ), n.typeName))
	}

	i := *p
	t = t.extend(p.typ)
	up := newTokenStack(1)
	up.setTrace(t)
	i.trace = t
	i.result = &parserResult{
		node: &node{
			nodeType: p.typ,
			typeName: p.typeName,
			toks:     []*token{nil},
		},
		unparsed: up,
	}
	return &i
}

// TODO: primitive is sometimes not cached. It is ok for the primitive, but not for the others.

func (p *primitiveParser) parse(t *token) *parserResult {
	p.trace.out("parsing", t)

	if t.typ == p.tokenType {
		p.result.valid = true
		p.result.node.setToken(t)
		p.trace.out("valid", p.result.node)
	} else {
		p.trace.out("invalid")
		p.result.valid = false
		p.result.unparsed.append(t)
	}

	return p.result
}

func (p *primitiveParser) valid() bool {
	return p.isValid
}

func (p *primitiveParser) expectedLength() int {
	return p.length
}

func (p *primitiveParser) nodeType() nodeType {
	return p.typ
}

func (g *optionalGenerator) create(t trace, init nodeType, excluded typeList) (parser, error) {
	t = t.extend(g.typ)

	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &optionalParser{typ: g.typ, isValid: true, typeName: g.registry.typeName(g.typ)}
	g.registry.setParser(g.typ, init, excluded, p)

	optional, ok := g.registry.get(g.optional)
	if !ok {
		return nil, unspecifiedParser(g.registry.typeName(g.optional))
	}

	if m, err := optional.member(g.typ); err != nil {
		return nil, err
	} else if m {
		return nil, optionalContainingSelf(g.registry.typeName(g.optional))
	}

	if excluded.contains(g.typ) {
		p.isValid = false
		return p, nil
	}

	excluded = append(excluded, g.typ)
	optParser, err := optional.create(t, init, excluded)
	if err != nil {
		return nil, err
	}

	var initIsMember bool
	if init != 0 {
		if m, err := g.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	length := optParser.expectedLength()
	if length == 0 {
		length = 1
	}

	p.registry = g.registry
	p.optional = optParser
	p.initIsMember = initIsMember
	p.length = optParser.expectedLength()
	return p, nil
}

func (g *optionalGenerator) member(t nodeType) (bool, error) {
	optional, ok := g.registry.get(g.optional)
	if !ok {
		return false, unspecifiedParser(g.registry.typeName(g.optional))
	}

	return optional.member(t)
}

func (g *optionalGenerator) nodeType() nodeType {
	return g.typ
}

func (p *optionalParser) check(trace) error {
	if p.optional != nil && !p.optional.valid() {
		p.optional = nil
	}

	return nil
}

func (p *optionalParser) reset() {
}

// TODO: review allocations in init, checking with benchmarks if they make any difference

func (p *optionalParser) instance(t trace, n *node) parser {
	t = t.extend(p.typ)
	up := newTokenStack(p.length)
	up.setTrace(t)

	var oi parser
	if p.optional != nil {
		oi = p.optional.instance(t, n)
	}

	i := *p
	i.trace = t
	i.initNode = n
	i.result = &parserResult{
		node:     zeroNode,
		unparsed: up,
	}
	i.optionalInstance = oi
	return &i
}

func (p *optionalParser) parse(t *token) *parserResult {
	p.trace.out("parsing", t)

	// p.cacheChecked = true

	p.cacheToken = t
	if p.initNode != zeroNode {
		p.cacheToken = p.initNode.token
	}

	if n, m, ok := p.cacheToken.cache.get(p.typ); ok {
		if m {
			if !p.initIsMember {
				p.result.valid = true
				p.result.node = n
				p.result.unparsed.append(t)
				p.result.fromCache = true
				p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
			}
		} else {
			p.result.valid = false
			p.result.unparsed.append(t)
			p.result.fromCache = true
			p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
		}

		return p.result
	}

	if p.optionalInstance != nil {
		p.optionalResult = p.optionalInstance.parse(t)
		if p.optionalResult.accepting {
			p.result.accepting = true
			return p.result
		}
	}

	p.result.accepting = false
	p.result.valid = true
	if p.optionalInstance == nil {
		p.result.unparsed.append(t)
	} else {
		p.result.unparsed.merge(p.optionalResult.unparsed)
	}

	if p.optionalInstance != nil && p.optionalResult.valid {
		p.trace.out("parse done, valid:", p.result.valid, p.result.node)
		p.result.node = p.optionalResult.node
		p.result.fromCache = p.optionalResult.fromCache
	} else if p.initIsMember {
		p.trace.out("init node is a member, valid", p.result.node)
		p.result.node = p.initNode
		p.result.fromCache = false
	} else {
		p.result.node = zeroNode
		p.trace.out("missing optional, valid", p.initIsMember)
	}

	p.cacheToken = p.result.node.token
	if p.result.node == zeroNode {
		if !p.result.unparsed.has() {
			panic(unexpectedResult(p.registry.typeName(p.typ)))
		}

		p.cacheToken = p.result.unparsed.peek()
	}

	p.cacheToken.cache.set(p.result.node, p.result.valid)
	return p.result
}

func (p *optionalParser) valid() bool {
	return p.isValid
}

func (p *optionalParser) expectedLength() int {
	return p.length
}

func (p *optionalParser) nodeType() nodeType {
	return p.typ
}

func (g *sequenceGenerator) create(t trace, init nodeType, excluded typeList) (parser, error) {
	t = t.extend(g.typ)

	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &sequenceParser{typ: g.typ, isValid: true, typeName: g.registry.typeName(g.typ)}
	g.registry.setParser(g.typ, init, excluded, p)

	item, ok := g.registry.get(g.item)
	if !ok {
		return nil, unspecifiedParser(g.registry.typeName(g.item))
	}

	if excluded.contains(g.typ) {
		p.isValid = false
		return p, nil
	}

	first, err := item.create(t, init, append(excluded, g.typ))
	if err != nil {
		return nil, err
	}

	rest, err := item.create(t, 0, nil)
	if err != nil {
		return nil, err
	}

	// TODO: panic on invalid rest (?)

	var initIsMember bool
	if init != 0 {
		if m, err := item.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	p.registry = g.registry
	p.first = first
	p.rest = rest
	p.initIsMember = initIsMember
	return p, nil
}

func (g *sequenceGenerator) member(t nodeType) (bool, error) {
	return t == g.typ, nil
}

func (g *sequenceGenerator) nodeType() nodeType {
	return g.typ
}

func (p *sequenceParser) check(trace) error {
	if p.first != nil && !p.first.valid() {
		p.first = nil
	}

	if p.rest != nil && !p.rest.valid() {
		p.rest = nil
	}

	var length int
	if p.first != nil {
		length = p.first.expectedLength()
	}

	if p.rest != nil && p.rest.expectedLength() > length {
		length = p.rest.expectedLength()
	}

	p.itemLength = length
	return nil
}

func (p *sequenceParser) reset() {
}

func (p *sequenceParser) instance(t trace, n *node) parser {
	t = t.extend(p.typ)
	up := newTokenStack(p.itemLength)
	up.setTrace(t)
	ts := newTokenStack(p.itemLength)
	ts.setTrace(t)

	i := *p
	i.trace = t
	i.initNode = n

	if p.first != nil {
		i.currentParser = p.first.instance(t, n)
	}

	i.result = &parserResult{
		node: &node{
			nodeType: p.typ,
			typeName: p.typeName,
		},
		unparsed: up,
	}

	i.tokenStack = ts
	return &i
}

func (p *sequenceParser) parse(t *token) *parserResult {
parseLoop:
	for {
		p.trace.out("parsing", t)

		if p.skip > 0 {
			p.skip--
			p.result.accepting = true
			return p.result
		}

		if p.skippingAfterDone {
			p.result.accepting = false
			return p.result
		}

		if !p.cacheChecked {
			p.cacheChecked = true

			p.cacheToken = t
			if p.initNode != zeroNode {
				p.cacheToken = p.initNode.token
			}

			if n, m, ok := p.cacheToken.cache.get(p.typ); ok {
				if m {
					if !p.initIsMember {
						p.result.valid = true
						p.result.node = n
						p.result.unparsed.append(t)
						p.result.fromCache = true
						p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
						return p.result
					}
				} else {
					p.result.valid = false
					p.result.unparsed.append(t)
					p.result.fromCache = true
					p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
					return p.result
				}
			}
		}

		if p.currentParser == nil {
			p.tokenStack.append(t)
			p.trace.out("parse done, valid", p.result.node)
			p.skippingAfterDone = p.skip > 0
			p.result.accepting = p.skippingAfterDone
			p.result.valid = true
			p.result.unparsed.merge(p.tokenStack)

			// NOTE: this was not set in parse4
			// maybe every node should have a token
			if p.result.node.token == nil {
				if !p.result.unparsed.has() {
					panic(unexpectedResult(p.registry.typeName(p.typ)))
				}

				p.result.node.token = p.result.unparsed.peek()
			}

			// NOTE: this was cached in parse4 only if there were nodes in the sequence
			p.cacheToken = p.result.node.token
			if p.cacheToken == nil {
				p.cacheToken = p.result.unparsed.peek()
			}

			p.cacheToken.cache.set(p.result.node, p.result.valid)
			return p.result
		}

		p.itemResult = p.currentParser.parse(t)
		if p.itemResult.accepting {
			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.tokenStack.merge(p.itemResult.unparsed)

		if p.itemResult.valid && p.itemResult.node != zeroNode {
			p.initEvaluated = true
			p.result.node.append(p.itemResult.node)
			if p.itemResult.fromCache {
				p.skip = p.tokenStack.findCachedNode(p.itemResult.node)
			}

			p.currentParser = p.rest.instance(p.trace, zeroNode)

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.initIsMember && !p.initEvaluated {
			p.initEvaluated = true
			p.result.node.append(p.initNode)

			p.currentParser = p.rest.instance(p.trace, zeroNode)

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.trace.out("parse done, valid", p.result.node)
		p.skippingAfterDone = p.skip > 0
		p.result.accepting = p.skippingAfterDone
		p.result.valid = true
		p.result.unparsed.merge(p.tokenStack)

		// NOTE: this was not set in parse4
		// maybe every node should have a token
		if p.result.node.token == nil {
			if !p.result.unparsed.has() {
				panic(unexpectedResult(p.registry.typeName(p.typ)))
			}

			p.result.node.token = p.result.unparsed.peek()
		}

		// NOTE: this was cached in parse4 only if there were nodes in the sequence
		p.cacheToken = p.result.node.token
		if p.cacheToken == nil {
			p.cacheToken = p.result.unparsed.peek()
		}

		p.cacheToken.cache.set(p.result.node, p.result.valid)
		return p.result
	}
}

func (p *sequenceParser) valid() bool {
	return p.isValid
}

func (p *sequenceParser) expectedLength() int {
	return p.itemLength
}

func (p *sequenceParser) nodeType() nodeType {
	return p.typ
}

func (g *groupGenerator) create(t trace, init nodeType, excluded typeList) (parser, error) {
	t = t.extend(g.typ)

	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &groupParser{isValid: true, initType: init, typ: g.typ, typeName: g.registry.typeName(g.typ)}
	g.registry.setParser(g.typ, init, excluded, p)

	if len(g.items) == 0 {
		return nil, groupWithoutItems(g.registry.typeName(g.typ))
	}

	if excluded.contains(g.typ) {
		p.isValid = false
		return p, nil
	}

	items := make([]generator, len(g.items))
	for i, it := range g.items {
		gi, ok := g.registry.get(it)
		if !ok {
			return nil, unspecifiedParser(g.registry.typeName(it))
		}

		items[i] = gi
	}

	excluded = append(excluded, g.typ)

	var (
		parsers      []parser
		initParsers  []parser
		initIsMember []bool
	)

	initType := init
	for i, gi := range items {
		var x typeList
		if i == 0 || initType != 0 {
			x = excluded
		}

		if i > 0 || initType == 0 {
			withoutInit, err := gi.create(t, 0, x)
			if err != nil {
				return nil, err
			}

			if !withoutInit.valid() {
				return nil, groupItemParserNotFound(g.registry.typeName(g.typ))
			}

			parsers = append(parsers, withoutInit)
		} else {
			parsers = append(parsers, nil)
		}

		if initType != 0 {
			withInit, err := gi.create(t, initType, x)
			if err != nil {
				return nil, err
			}

			m, err := gi.member(initType)
			if err != nil {
				return nil, err
			}

			if !m && !withInit.valid() {
				p.isValid = false
				return p, nil
			}

			// needs a nil check in the parser
			initParsers = append(initParsers, withInit)
			initIsMember = append(initIsMember, m)

			if withInit.valid() && withInit.expectedLength() > 0 || m {
				initType = 0
			}
		}
	}

	p.registry = g.registry
	p.parsers = parsers
	p.initParsers = initParsers
	p.initIsMember = initIsMember
	return p, nil
}

func (g *groupGenerator) member(t nodeType) (bool, error) {
	return t == g.typ, nil
}

func (g *groupGenerator) nodeType() nodeType {
	return g.typ
}

func (p *groupParser) check(trace) error {
	for _, pi := range p.parsers {
		if pi != nil && !pi.valid() {
			return groupItemParserNotFound(p.typeName)
		}
	}

	if p.initType == 0 {
		return nil
	}

	var foundInvalid bool
	for i, pi := range p.initParsers {
		if foundInvalid {
			p.initParsers[i] = nil
			continue
		}

		if pi.valid() {
			continue
		}

		if i == 0 {
			p.isValid = false
			return nil
		}

		p.initParsers[i] = nil
		foundInvalid = true
	}

	var length int
	parsers := p.parsers
	initParsers := p.initParsers
	for {
		if len(parsers) == 0 && len(initParsers) == 0 {
			break
		}

		var without int
		if len(parsers) > 0 {
			if parsers[0] != nil {
				without = parsers[0].expectedLength()
			}

			parsers = parsers[1:]
		}

		var with int
		if len(initParsers) > 0 {
			if initParsers[0] != nil {
				with = initParsers[0].expectedLength()
			}

			initParsers = initParsers[1:]
		}

		if with > without {
			length += with
		} else {
			length += without
		}
	}

	p.length = length
	return nil
}

func (p *groupParser) reset() {
}

func (p *groupParser) instance(t trace, n *node) parser {
	t = t.extend(p.typ)
	up := newTokenStack(p.length)
	up.setTrace(t)
	ts := newTokenStack(p.length)
	ts.setTrace(t)

	i := *p
	i.trace = t
	i.initNode = n
	i.result = &parserResult{
		node: &node{
			nodeType: p.typ,
			typeName: p.typeName,
			nodes:    make([]*node, 0, len(p.parsers)),
			toks:     nil, // TODO: this can be preallocated as the sum of tokens of the items
		},
		unparsed: up,
	}
	i.tokenStack = ts
	return &i
}

// TODO: test when subsequent required init parser is nil

func (p *groupParser) parse(t *token) *parserResult {
parseLoop:
	for {
		p.trace.out("parsing", t)

		if p.skip > 0 {
			p.skip--
			p.result.accepting = true
			return p.result
		}

		if p.skippingAfterDone {
			p.result.accepting = false
			return p.result
		}

		// need to check membership before the cache check

		if !p.cacheChecked {
			p.cacheChecked = true

			p.cacheToken = t
			if p.initNode != zeroNode {
				p.cacheToken = p.initNode.token
			}

			if n, m, ok := p.cacheToken.cache.get(p.typ); ok {
				if m {
					if len(p.initIsMember) > p.parserIndex && !p.initIsMember[p.parserIndex] {
						p.result.valid = true
						p.result.node = n
						p.result.unparsed.append(t)
						p.result.fromCache = true
						p.result.accepting = false
						p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
						return p.result
					}
				} else {
					p.result.valid = false
					p.result.unparsed.append(t)
					p.result.fromCache = true
					p.result.accepting = false
					p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
					return p.result
				}
			}
		}

		if p.currentParser == nil {
			if p.initNode == zeroNode || p.initEvaluated {
				p.currentParser = p.parsers[p.parserIndex].instance(p.trace, zeroNode)
			} else if p.parserIndex < len(p.initParsers) {
				if p.initParsers[p.parserIndex] != nil {
					p.currentParser = p.initParsers[p.parserIndex].instance(p.trace, p.initNode)
				}
			}
		}

		// can be nil with init, only to check membership:
		if p.currentParser == nil {
			p.tokenStack.append(t)
			p.itemResult = nil
		} else {
			p.itemResult = p.currentParser.parse(t)
			if p.itemResult.accepting {
				p.result.accepting = true
				if p.tokenStack.has() {
					t = p.tokenStack.pop()
					continue parseLoop
				}

				p.result.accepting = true
				return p.result
			}

			p.tokenStack.merge(p.itemResult.unparsed)
			p.currentParser = nil
		}

		p.parserIndex++

		if p.itemResult != nil && p.itemResult.valid && p.itemResult.node != zeroNode {
			p.initEvaluated = true
			p.result.node.append(p.itemResult.node)
			if p.itemResult.fromCache {
				p.skip = p.tokenStack.findCachedNode(p.itemResult.node)
			}

			if p.parserIndex == len(p.parsers) {
				p.trace.out("group done, valid", p.result.node)
				p.result.valid = true

				p.cacheToken = p.result.node.token
				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.typeName))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cacheToken.cache.set(p.result.node, p.result.valid)

				p.result.unparsed.merge(p.tokenStack)
				p.skippingAfterDone = p.skip > 0
				p.result.accepting = p.skippingAfterDone
				return p.result
			}

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.itemResult != nil && p.itemResult.valid {
			if p.parserIndex == len(p.parsers) {
				p.trace.out("group done, valid", p.result.node)
				p.result.valid = true

				p.cacheToken = p.result.node.token
				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.typeName))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cacheToken.cache.set(p.result.node, p.result.valid)

				p.result.unparsed.merge(p.tokenStack)
				p.skippingAfterDone = p.skip > 0
				p.result.accepting = p.skippingAfterDone
				return p.result
			}

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.initNode != nil && !p.initEvaluated && len(p.initIsMember) >= p.parserIndex && p.initIsMember[p.parserIndex-1] {
			p.initEvaluated = true

			p.result.node.append(p.initNode)

			if p.parserIndex == len(p.parsers) {
				p.trace.out("group done, valid", p.result.node)
				p.result.valid = true

				p.cacheToken = p.result.node.token
				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.typeName))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cacheToken.cache.set(p.result.node, p.result.valid)

				p.result.unparsed.merge(p.tokenStack)
				p.result.accepting = false
				return p.result
			}

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.trace.out("group done, invalid")

		p.result.valid = false
		p.result.accepting = false

		p.cacheToken = p.result.node.token
		if p.cacheToken == nil {
			if !p.tokenStack.has() {
				panic(unexpectedResult(p.typeName))
			}

			p.cacheToken = p.tokenStack.peek()
		}

		p.cacheToken.cache.set(p.result.node, p.result.valid)

		p.result.unparsed.merge(p.tokenStack)
		if p.result.node.len() > p.initNode.len() {
			p.result.unparsed.mergeTokens(p.result.node.tokens()[p.initNode.len():])
		}

		return p.result
	}
}

func (p *groupParser) valid() bool {
	return p.isValid
}

func (p *groupParser) expectedLength() int {
	return p.length
}

func (p *groupParser) nodeType() nodeType {
	return p.typ
}

func (g *unionGenerator) expand(ignore typeList) ([]generator, error) {
	if ignore.contains(g.typ) {
		return nil, nil
	}

	var generators []generator
	for _, et := range g.elementTypes {
		eg, ok := g.registry.get(et)
		if !ok {
			return nil, unspecifiedParser(g.registry.typeName(et))
		}

		if ug, ok := eg.(*unionGenerator); ok {
			ugx, err := ug.expand(append(ignore, g.typ))
			if err != nil {
				return nil, err
			}

			generators = append(generators, ugx...)
		} else if !ignore.contains(et) {
			generators = append(generators, eg)
		}
	}

	return generators, nil
}

func (g *unionGenerator) checkExpand() error {
	if len(g.elements) > 0 {
		return nil
	}

	elements, err := g.expand(nil)
	if err != nil {
		return err
	}

	if len(elements) == 0 {
		return unionWithoutElements(g.registry.typeName(g.typ))
	}

	g.elements = elements
	return nil
}

func (g *unionGenerator) create(t trace, init nodeType, excluded typeList) (parser, error) {
	t = t.extend(g.typ)

	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &unionParser{typ: g.typ, isValid: true, typeName: g.registry.typeName(g.typ)}
	g.registry.setParser(g.typ, init, excluded, p)

	if err := g.checkExpand(); err != nil {
		return nil, err
	}

	expandedTypes := make([]nodeType, len(g.elements))
	for i, e := range g.elements {
		expandedTypes[i] = e.nodeType()
	}

	parsers := make([][]parser, len(g.elements)+1)
	for i, it := range append([]nodeType{init}, expandedTypes...) {
		p := make([]parser, len(g.elements))
		for j, e := range g.elements {
			pe, err := e.create(t, it, excluded)
			if err != nil {
				return nil, err
			}

			if pe.valid() {
				p[j] = pe
			}
		}

		parsers[i] = p
	}

	var initIsMember bool
	if init != 0 {
		if m, err := g.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	if !initIsMember && (len(parsers[0]) == 0 || parsers[0][0] == nil) {
		p.isValid = false
		return p, nil
	}

	p.registry = g.registry
	p.typ = g.typ
	p.parsers = parsers
	p.initIsMember = initIsMember
	return p, nil
}

func (g *unionGenerator) member(t nodeType) (bool, error) {
	if err := g.checkExpand(); err != nil {
		return false, err
	}

	for _, e := range g.elements {
		if m, err := e.member(t); m || err != nil {
			return m, err
		}
	}

	return false, nil
}

func (g *unionGenerator) nodeType() nodeType {
	return g.typ
}

func (p *unionParser) check(t trace) error {
	var (
		hasValid bool
		length   int
	)

	for _, ps := range p.parsers {
		for i, pp := range ps {
			if pp != nil && pp.valid() {
				hasValid = true
				if pp.expectedLength() > length {
					length = pp.expectedLength()
				}

				continue
			}

			ps[i] = nil
		}
	}

	p.isValid = hasValid || p.initIsMember
	p.length = length

	return nil
}

func (p *unionParser) reset() {
}

func (p *unionParser) instance(t trace, n *node) parser {
	t = t.extend(p.typ)
	up := newTokenStack(p.length)
	up.setTrace(t)
	ts := newTokenStack(p.length)
	ts.setTrace(t)

	i := *p
	i.trace = t
	i.initNode = n

	i.result = &parserResult{
		unparsed: up,
	}

	if p.initIsMember {
		i.result.node = n
	} else {
		i.result.node = zeroNode
	}

	if p.parsers[0][0] != nil {
		i.currentParser = p.parsers[0][0].instance(t, n)
	}

	i.tokenStack = ts
	return &i
}

func (p *unionParser) parse(t *token) *parserResult {
parseLoop:
	for {
		p.trace.out("parsing", t)

		if p.skip > 0 {
			p.skip--
			p.result.accepting = true
			return p.result
		}

		if p.skippingAfterDone {
			p.result.accepting = false
			return p.result
		}

		// parse4 not using cache in union
		// maybe the new caching style helps
		if !p.cacheChecked {
			p.cacheChecked = true

			p.cacheToken = t
			if p.initNode != zeroNode {
				p.cacheToken = p.initNode.token
			}

			if n, m, ok := p.cacheToken.cache.get(p.typ); ok {
				if m {
					if !p.initIsMember {
						p.result.valid = true
						p.result.node = n
						p.result.unparsed.append(t)
						p.result.fromCache = true
						p.result.accepting = false
						p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
						return p.result
					}
				} else {
					p.result.valid = false
					p.result.unparsed.append(t)
					p.result.fromCache = true
					p.result.accepting = false
					p.trace.out("found in cache, valid:", p.result.valid, p.result.node)
					return p.result
				}
			}
		}

		if p.currentParser != nil {
			p.itemResult = p.currentParser.parse(t)
			if p.itemResult.accepting {
				if p.tokenStack.has() {
					t = p.tokenStack.pop()
					continue parseLoop
				}

				p.result.accepting = true
				return p.result
			}

			p.tokenStack.merge(p.itemResult.unparsed)
		} else {
			p.tokenStack.append(t)
		}

		if p.itemResult == nil || !p.itemResult.valid {
			for {
				p.parserIndex++

				if p.parserIndex == len(p.parsers[p.initTypeIndex]) {
					break
				}

				if p.parsers[p.initTypeIndex][p.parserIndex] != nil {
					break
				}
			}

			if p.parserIndex == len(p.parsers[p.initTypeIndex]) {
				p.trace.out("done, valid:", p.result.node != zeroNode, p.result.node)
				p.result.accepting = false
				p.result.valid = p.result.node != zeroNode

				p.cacheToken = p.result.node.token

				if p.cacheToken == nil {
					p.cacheToken = p.initNode.token
				}

				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.typeName))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cacheToken.cache.set(p.result.node, p.result.valid)

				p.result.unparsed.merge(p.tokenStack)
				return p.result
			}

			p.currentParser = p.parsers[p.initTypeIndex][p.parserIndex].instance(p.trace, p.result.node)

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.result.node == zeroNode || p.itemResult != nil && p.result.node.len() < p.itemResult.node.len() {
			p.result.node = p.itemResult.node

			p.initTypeIndex = p.parserIndex + 1
			if p.initTypeIndex < len(p.parsers) {
				p.parserIndex = 0
				for {
					if p.parsers[p.initTypeIndex][p.parserIndex] != nil {
						break
					}

					p.parserIndex++

					if p.parserIndex == len(p.parsers[p.initTypeIndex]) {
						break
					}
				}
			}

			if p.itemResult.fromCache {
				p.skip = p.tokenStack.findCachedNode(p.itemResult.node)
			}
		} else {
			for {
				p.parserIndex++

				if p.parserIndex == len(p.parsers[p.initTypeIndex]) {
					break
				}

				if p.parsers[p.initTypeIndex][p.parserIndex] != nil {
					break
				}
			}
		}

		if p.initTypeIndex < len(p.parsers) && p.parserIndex < len(p.parsers[p.initTypeIndex]) {
			p.currentParser = p.parsers[p.initTypeIndex][p.parserIndex].instance(p.trace, p.result.node)

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.result.accepting = false
		p.result.valid = p.result.node != zeroNode

		p.cacheToken = p.result.node.token

		if p.cacheToken == nil {
			p.cacheToken = p.initNode.token
		}

		if p.cacheToken == nil {
			if !p.tokenStack.has() {
				panic(unexpectedResult(p.typeName))
			}

			p.cacheToken = p.tokenStack.peek()
		}

		p.cacheToken.cache.set(p.result.node, p.result.valid)

		p.result.unparsed.merge(p.tokenStack)

		if p.skip > 0 {
			p.skippingAfterDone = true
			p.result.accepting = true
			return p.result
		}

		return p.result
	}
}

func (p *unionParser) valid() bool {
	return p.isValid
}

func (p *unionParser) expectedLength() int {
	return p.length
}

func (p *unionParser) nodeType() nodeType {
	return p.typ
}

func newSyntax() *syntax {
	return &syntax{
		registry: newRegistry(),
	}
}

func (s *syntax) primitive(typeName string, t tokenType) error {
	return s.registry.primitive(typeName, t)
}

func (s *syntax) optional(typeName string, optional string) error {
	return s.registry.optional(typeName, optional)
}

func (s *syntax) sequence(typeName string, item string) error {
	return s.registry.sequence(typeName, item)
}

func (s *syntax) group(typeName string, items ...string) error {
	return s.registry.group(typeName, items...)
}

func (s *syntax) union(typeName string, elements ...string) error {
	return s.registry.union(typeName, elements...)
}

func (s *syntax) parse(r *tokenReader) (*node, error) {
	eofToken.cache = newCache()
	s.registry.reset()

	root, err := s.registry.rootGenerator()
	if err != nil {
		return zeroNode, err
	}

	trace := newTrace(s.traceLevel, s.registry)

	parser, err := root.create(trace, 0, nil)
	if err != nil {
		return zeroNode, err
	}

	if err := s.registry.finalize(trace); err != nil {
		return zeroNode, err
	}

	if !parser.valid() {
		return nil, errFailedToCreateParser
	}

	parserInstance := parser.instance(trace, zeroNode)

	last := &parserResult{accepting: true, node: zeroNode}
	for {
		t, err := r.next()
		if err != nil && err != io.EOF {
			return zeroNode, err
		}

		if !last.accepting {
			if err != io.EOF {
				return zeroNode, unexpectedToken("root", &t)
			}

			// TODO: should we check here if it is valid?
			if !last.valid {
				return zeroNode, errUnexpectedEOF
			}

			return last.node, nil
		}

		if err == io.EOF {
			last = parserInstance.parse(eofToken)

			if !last.valid {
				return zeroNode, errUnexpectedEOF
			}

			if !last.unparsed.has() || last.unparsed.peek() != eofToken {
				return zeroNode, errUnexpectedEOF
			}

			return last.node, nil
		}

		t.cache = newCache()
		last = parserInstance.parse(&t)
		if !last.accepting {
			if !last.valid {
				return zeroNode, unexpectedToken("root", &t)
			}

			if last.unparsed.has() {
				return zeroNode, unexpectedToken("root", last.unparsed.peek())
			}
		}
	}
}

func defineSyntax(primitive [][]interface{}, complex [][]string) (*syntax, error) {
	s := newSyntax()

	for _, p := range primitive {
		if err := s.primitive(p[0].(string), p[1].(tokenType)); err != nil {
			return nil, err
		}
	}

	for _, c := range complex {
		var err error
		switch c[0] {
		case "optional":
			err = s.optional(c[1], c[2])
		case "sequence":
			err = s.sequence(c[1], c[2])
		case "group":
			err = s.group(c[1], c[2:]...)
		case "union":
			err = s.union(c[1], c[2:]...)
		default:
			panic("invalid parser type")
		}

		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
