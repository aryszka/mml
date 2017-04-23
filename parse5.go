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
	traceOn
	traceDebug
)

type nodeType int

type typeList []nodeType

type node struct {
	token    *token
	nodeType nodeType
	typeName string
	nodes    []*node
	toks     []*token
}

type generatorResult struct {
	valid          bool
	parser         parser
	expectedLength int

	// hack for recursive generation
	preallocated bool
	required     bool
}

type parserResult struct {
	accepting bool
	valid     bool
	node      *node
	fromCache bool
	unparsed  *tokenStack
}

type parser interface {
	init(t trace, n *node)
	parse(*token) *parserResult // optimization: should not make memory allocations
}

type generator interface {
	create(nodeType, typeList) (*generatorResult, error)
	member(nodeType) (bool, error)
	nodeType() nodeType
}

type trace interface {
	extend(nodeType) trace
	outLevel(traceLevel, ...interface{})
	out(...interface{})
	debug(...interface{})
}

type registry interface {
	typeName(nodeType) string
	nodeType(typeName string) nodeType
	get(nodeType) (generator, bool)
	getParser(nodeType, nodeType, typeList) (*generatorResult, bool)
	setParser(nodeType, nodeType, typeList, *generatorResult)
	clearParsers()
	root() (generator, error)
	primitive(string, tokenType) error
	optional(string, string) error
	sequence(string, string) error
	group(string, ...string) error
	union(string, ...string) error
}

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

type cache struct {
	result *parserResult
	offset int
}

type parserRegistry struct {
	idSeed        nodeType
	typeIDs       map[string]nodeType
	typeNames     map[nodeType]string
	generators    map[nodeType]generator
	parsers       map[string]*generatorResult
	rootGenerator generator
}

type primitiveGenerator struct {
	typ       nodeType
	tokenType tokenType
	registry  registry
}

type primitiveParser struct {
	trace     trace
	registry  registry
	nodeType  nodeType
	tokenType tokenType
	result    *parserResult
	cache     *cache
}

type optionalGenerator struct {
	typ      nodeType
	optional nodeType
	registry registry
}

type optionalParser struct {
	trace          trace
	registry       registry
	nodeType       nodeType
	optional       parser
	initIsMember   bool
	initNode       *node
	result         *parserResult
	optionalResult *parserResult
	cacheToken     *token
	cache          *cache
	cacheChecked   bool
}

type sequenceGenerator struct {
	typ      nodeType
	item     nodeType
	registry registry
}

type sequenceParser struct {
	trace             trace
	registry          registry
	nodeType          nodeType
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
	cache             *cache
	skippingAfterDone bool
}

type groupGenerator struct {
	typ      nodeType
	items    []nodeType
	registry registry
}

type groupParser struct {
	trace             trace
	registry          registry
	parsers           []parser
	initParsers       []parser
	currentParser     parser
	initIsMember      []bool
	result            *parserResult
	tokenStack        *tokenStack
	cache             *cache
	initNode          *node
	skip              int
	parserIndex       int
	cacheChecked      bool
	initEvaluated     bool
	skippingAfterDone bool
	cacheToken        *token
	itemResult        *parserResult
}

type unionGenerator struct {
	typ          nodeType
	registry     registry
	elements     []generator
	elementTypes []nodeType
}

type unionParser struct {
	trace             trace
	cache             *cache
	registry          registry
	nodeType          nodeType
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
	eofToken = &token{offset: -1, value: "<eof>"}
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
	if len(n.nodes) == 1 {
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

func newTokenStack(expectedSize int) *tokenStack {
	return &tokenStack{
		stack: make([]*token, 0, expectedSize),
	}
}

func (s *tokenStack) append(t *token) {
	if len(s.stack) == cap(s.stack) {
		s.trace.debug("token stack exceeded expected size on append")
	}

	s.stack = append(s.stack, t)
}

func (s *tokenStack) merge(from *tokenStack) {
	s.need = len(s.stack) + len(from.stack) - cap(s.stack)
	if s.need > 0 {
		s.trace.debug("token stack exceeded expected size on merge")

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
		s.trace.debug("stack dropping more tokens than it contains")
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

func (t *parserTrace) debug(a ...interface{}) {
	t.outLevel(traceDebug, a...)
}

func (c *cache) set(offset int, r *parserResult) {
	c.offset = offset
	c.result = r
}

func (c *cache) has(offset int) bool {
	if offset != c.offset {
		return false
	}

	return c.result != nil
}

func (c *cache) get() *parserResult {
	c.result.fromCache = true
	return c.result
}

func newRegistry() *parserRegistry {
	return &parserRegistry{
		idSeed:     1,
		typeIDs:    make(map[string]nodeType),
		typeNames:  make(map[nodeType]string),
		generators: make(map[nodeType]generator),
		parsers:    make(map[string]*generatorResult),
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

func (r *parserRegistry) getParser(t nodeType, init nodeType, excluded typeList) (*generatorResult, bool) {
	p, ok := r.parsers[parserKey(t, init, excluded)]
	return p, ok
}

func (r *parserRegistry) setParser(t nodeType, init nodeType, excluded typeList, p *generatorResult) {
	r.parsers[parserKey(t, init, excluded)] = p
}

func (r *parserRegistry) clearParsers() {
	r.parsers = make(map[string]*generatorResult)
}

func (r *parserRegistry) root() (generator, error) {
	if r.rootGenerator == nil {
		return nil, errNoParsersDefined
	}

	return r.rootGenerator, nil
}

func (r *parserRegistry) register(nt nodeType, g generator) error {
	if _, exists := r.generators[nt]; exists {
		return duplicateNodeType(r.typeNames[nt])
	}

	r.generators[nt] = g
	r.rootGenerator = g // the last one is the root
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

func (g *primitiveGenerator) create(init nodeType, excluded typeList) (*generatorResult, error) {
	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &generatorResult{parser: &primitiveParser{nodeType: g.typ, registry: g.registry}, preallocated: true}
	g.registry.setParser(g.typ, init, excluded, p)

	if excluded.contains(g.typ) || init != 0 {
		if p.required {
			return nil, requiredParserInvalid(g.registry.typeName(g.typ))
		}

		p.preallocated = false
		p.parser = nil
		return p, nil
	}

	n := &node{
		nodeType: g.typ,
		typeName: g.registry.typeName(g.typ),
	}

	p.valid = true
	p.parser.(*primitiveParser).nodeType = g.typ
	p.parser.(*primitiveParser).tokenType = g.tokenType
	p.parser.(*primitiveParser).result = &parserResult{
		node:     n,
		unparsed: newTokenStack(1),
	}
	p.parser.(*primitiveParser).cache = &cache{}
	p.expectedLength = 1
	p.preallocated = false
	return p, nil
}

func (g *primitiveGenerator) member(t nodeType) (bool, error) {
	return t == g.typ, nil
}

func (g *primitiveGenerator) nodeType() nodeType {
	return g.typ
}

// TODO: try lazy node creation

func (p *primitiveParser) init(t trace, n *node) {
	if n != zeroNode {
		panic(unexpectedInitNode(p.registry.typeName(p.nodeType), n.typeName))
	}

	p.trace = t.extend(p.result.node.nodeType)
	p.result.node = &node{
		nodeType: p.result.node.nodeType,
		typeName: p.result.node.typeName,
		toks:     []*token{nil},
	}

	p.result.unparsed.clear()
	p.result.unparsed.setTrace(p.trace)
}

func (p *primitiveParser) parse(t *token) *parserResult {
	p.trace.out("parsing", t)

	if p.cache.has(t.offset) {
		p.trace.out("found in cache, valid:", p.result.valid)
		p.result = p.cache.get()
		p.result.unparsed.append(t)
		return p.result
	}

	if t.typ == p.tokenType {
		p.trace.out("valid")
		p.result.valid = true
		p.result.node.setToken(t)
	} else {
		p.trace.out("invalid")
		p.result.valid = false
		p.result.unparsed.append(t)
	}

	p.cache.set(t.offset, p.result)
	return p.result
}

// TODO: the group should hanlde when an optional item in the beginning doesn't accept the init node that
// otherwise can be good for a later node

func (g *optionalGenerator) create(init nodeType, excluded typeList) (*generatorResult, error) {
	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &generatorResult{parser: &optionalParser{}, preallocated: true}
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
		if p.required {
			return nil, requiredParserInvalid(g.registry.typeName(g.typ))
		}

		p.preallocated = false
		p.parser = nil
		return p, nil
	}

	excluded = append(excluded, g.typ)
	optParser, err := optional.create(init, excluded)
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

	expectedUnparsedLength := optParser.expectedLength
	if expectedUnparsedLength == 0 {
		expectedUnparsedLength = 1
	}

	p.valid = true
	p.parser.(*optionalParser).registry = g.registry
	p.parser.(*optionalParser).nodeType = g.typ
	p.parser.(*optionalParser).optional = optParser.parser
	p.parser.(*optionalParser).initIsMember = initIsMember
	p.parser.(*optionalParser).result = &parserResult{
		node:     zeroNode,
		unparsed: newTokenStack(expectedUnparsedLength),
	}
	p.parser.(*optionalParser).cache = &cache{}
	p.expectedLength = optParser.expectedLength
	p.preallocated = false
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

func (p *optionalParser) init(t trace, n *node) {
	p.trace = t.extend(p.nodeType)
	p.initNode = n
	p.result.node = zeroNode
	p.result.unparsed.clear()
	p.result.unparsed.setTrace(p.trace)
	p.cacheChecked = false
	if p.optional != nil {
		p.optional.init(p.trace, n)
	}
}

func (p *optionalParser) parse(t *token) *parserResult {
	p.trace.out("parsing", t)

	if !p.cacheChecked {
		p.cacheChecked = true

		p.cacheToken = t
		if p.initNode != zeroNode {
			p.cacheToken = p.initNode.token
		}

		if p.cache.has(p.cacheToken.offset) {
			p.trace.out("found in cache, valid:", p.result.valid)
			p.result = p.cache.get()
			p.result.unparsed.append(t)
			return p.result
		}
	}

	if p.optional != nil {
		p.optionalResult = p.optional.parse(t)
		if p.optionalResult.accepting {
			p.result.accepting = true
			return p.result
		}
	}

	p.result.accepting = false
	p.result.valid = true
	if p.optional == nil {
		p.result.unparsed.append(t)
	} else {
		p.result.unparsed.merge(p.optionalResult.unparsed)
	}

	if p.optional != nil && p.optionalResult.valid {
		p.trace.out("parse done, valid:", p.result.valid)
		p.result.node = p.optionalResult.node
		p.result.fromCache = p.optionalResult.fromCache
	} else if p.initIsMember {
		p.trace.out("init node is a member, valid")
		p.result.node = p.initNode
		p.result.fromCache = false
	} else {
		p.result.node = zeroNode
		p.trace.out("missing optional, valid", p.initIsMember)
	}

	p.cacheToken = p.result.node.token
	if p.result.node == zeroNode {
		if !p.result.unparsed.has() {
			panic(unexpectedResult(p.registry.typeName(p.nodeType)))
		}

		p.cacheToken = p.result.unparsed.peek()
	}

	p.cache.set(p.cacheToken.offset, p.result)
	return p.result
}

func (g *sequenceGenerator) create(init nodeType, excluded typeList) (*generatorResult, error) {
	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &generatorResult{parser: &sequenceParser{}, preallocated: true}
	g.registry.setParser(g.typ, init, excluded, p)

	item, ok := g.registry.get(g.item)
	if !ok {
		return nil, unspecifiedParser(g.registry.typeName(g.item))
	}

	if m, err := g.member(g.item); err != nil {
		return nil, err
	} else if m {
		return nil, sequenceContainingSelf(g.registry.typeName(g.typ))
	}

	if excluded.contains(g.typ) {
		if p.required {
			return nil, requiredParserInvalid(g.registry.typeName(g.typ))
		}

		p.preallocated = false
		p.parser = nil
		return p, nil
	}

	allExcluded := append(excluded, g.typ)
	selfExcluded := typeList{g.typ}

	first, err := item.create(init, allExcluded)
	if err != nil {
		return nil, err
	}

	rest, err := item.create(0, selfExcluded)
	if err != nil {
		return nil, err
	}

	if !rest.valid && !rest.preallocated {
		panic(itemParserCannotBeCreated(g.registry.typeName(g.typ)))
	}

	if rest.preallocated {
		rest.required = true
	}

	var initIsMember bool
	if init != 0 {
		if m, err := item.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	if !first.valid && !first.preallocated && !initIsMember {
		if p.required {
			return nil, requiredParserInvalid(g.registry.typeName(g.typ))
		}

		p.preallocated = false
		p.parser = nil
		return p, nil
	}

	if first.preallocated && !initIsMember {
		first.required = true
	}

	expectedLength := first.expectedLength
	if rest.expectedLength > expectedLength {
		expectedLength = rest.expectedLength
	}

	p.valid = true
	p.parser.(*sequenceParser).registry = g.registry
	p.parser.(*sequenceParser).first = first.parser
	p.parser.(*sequenceParser).rest = rest.parser
	p.parser.(*sequenceParser).initIsMember = initIsMember
	p.parser.(*sequenceParser).result = &parserResult{
		node: &node{
			nodeType: g.typ,
			typeName: g.registry.typeName(g.typ),
		},
		unparsed: newTokenStack(expectedLength),
	}
	p.parser.(*sequenceParser).tokenStack = newTokenStack(expectedLength)
	p.parser.(*sequenceParser).cache = &cache{}
	p.expectedLength = expectedLength
	p.preallocated = false
	return p, nil
}

func (g *sequenceGenerator) member(t nodeType) (bool, error) {
	return t == g.typ, nil
}

func (g *sequenceGenerator) nodeType() nodeType {
	return g.typ
}

func (p *sequenceParser) init(t trace, n *node) {
	p.trace = t.extend(p.result.node.nodeType)
	p.initNode = n
	p.cacheChecked = false
	p.currentParser = p.first
	p.currentParser.init(p.trace, n)
	p.result.node = &node{
		nodeType: p.result.node.nodeType,
		typeName: p.result.node.typeName,
	}
	p.skip = 0
	p.result.unparsed.clear()
	p.result.unparsed.setTrace(p.trace)
	p.tokenStack.clear()
	p.tokenStack.setTrace(p.trace)
	p.initEvaluated = false
	p.skippingAfterDone = false
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

			if p.cache.has(p.cacheToken.offset) {
				p.trace.out("found in cache, valid:", p.result.valid)
				p.result = p.cache.get()
				p.result.unparsed.append(t)
				return p.result
			}
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
		p.currentParser = p.rest

		if p.itemResult.valid && p.itemResult.node != zeroNode {
			p.initEvaluated = true
			p.result.node.append(p.itemResult.node)
			if p.itemResult.fromCache {
				p.skip = p.tokenStack.findCachedNode(p.itemResult.node)
			}

			p.currentParser.init(p.trace, zeroNode)

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

			p.currentParser.init(p.trace, zeroNode)

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.trace.out("parse done, valid")
		p.skippingAfterDone = p.skip > 0
		p.result.accepting = p.skippingAfterDone
		p.result.valid = true
		p.result.unparsed.merge(p.tokenStack)

		// NOTE: this was not set in parse4
		// maybe every node should have a token
		if p.result.node.token == nil {
			if !p.result.unparsed.has() {
				panic(unexpectedResult(p.registry.typeName(p.nodeType)))
			}

			p.result.node.token = p.result.unparsed.peek()
		}

		// NOTE: this was cached in parse4 only if there were nodes in the sequence
		p.cacheToken = p.result.node.token
		if p.cacheToken == nil {
			p.cacheToken = p.result.unparsed.peek()
		}

		p.cache.set(p.cacheToken.offset, p.result)
		return p.result
	}
}

// TODO: if there is an init, then cannot return valid in group when there is an init unless everyting is
// optional. Can check the expected length.

func (g *groupGenerator) create(init nodeType, excluded typeList) (*generatorResult, error) {
	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &generatorResult{parser: &groupParser{}, preallocated: true}
	g.registry.setParser(g.typ, init, excluded, p)

	if len(g.items) == 0 {
		return nil, groupWithoutItems(g.registry.typeName(g.typ))
	}

	if excluded.contains(g.typ) {
		if p.required {
			return nil, requiredParserInvalid(g.registry.typeName(g.typ))
		}

		p.preallocated = false
		p.parser = nil
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
		parsers        []parser
		initParsers    []parser
		initIsMember   []bool
		expectedLength int
	)

	initType := init
	for i, gi := range items {
		var x typeList
		if i == 0 || initType != 0 {
			x = excluded
		}

		var lengthWithout, lengthWith int

		if i > 0 || initType == 0 {
			withoutInit, err := gi.create(0, x)
			if err != nil {
				return nil, err
			}

			if !withoutInit.valid && !withoutInit.preallocated {
				panic(groupItemParserNotFound(g.registry.typeName(g.typ)))
			}

			if withoutInit.preallocated {
				withoutInit.required = true
			}

			parsers = append(parsers, withoutInit.parser)
			lengthWithout = withoutInit.expectedLength
		} else {
			parsers = append(parsers, nil)
		}

		if initType != 0 {
			withInit, err := gi.create(initType, x)
			if err != nil {
				return nil, err
			}

			m, err := gi.member(initType)
			if err != nil {
				return nil, err
			}

			if !m && !withInit.valid && !withInit.preallocated {
				if p.required {
					return nil, requiredParserInvalid(g.registry.typeName(g.typ))
				}

				p.preallocated = false
				p.parser = nil
				return p, nil
			}

			if !m && withInit.preallocated {
				withInit.required = true
			}

			// needs a nil check in the parser
			initParsers = append(initParsers, withInit.parser)
			initIsMember = append(initIsMember, m)

			if withInit.valid {
				lengthWith = withInit.expectedLength
			}

			if withInit.valid && withInit.expectedLength > 0 || m {
				initType = 0
			}
		}

		if lengthWith > lengthWithout {
			expectedLength += lengthWith
		} else {
			expectedLength += lengthWithout
		}
	}

	p.valid = true
	p.parser.(*groupParser).registry = g.registry
	p.parser.(*groupParser).parsers = parsers
	p.parser.(*groupParser).initParsers = initParsers
	p.parser.(*groupParser).initIsMember = initIsMember
	p.parser.(*groupParser).result = &parserResult{
		node: &node{
			nodeType: g.typ,
			typeName: g.registry.typeName(g.typ),
		},
		unparsed: newTokenStack(expectedLength),
	}
	p.parser.(*groupParser).tokenStack = newTokenStack(expectedLength)
	p.parser.(*groupParser).cache = &cache{}
	p.expectedLength = expectedLength
	p.preallocated = false
	return p, nil
}

func (g *groupGenerator) member(t nodeType) (bool, error) {
	return t == g.typ, nil
}

func (g *groupGenerator) nodeType() nodeType {
	return g.typ
}

func (p *groupParser) init(t trace, n *node) {
	p.trace = t.extend(p.result.node.nodeType)
	p.initNode = n
	p.result.node = &node{
		nodeType: p.result.node.nodeType,
		typeName: p.result.node.typeName,
		nodes:    make([]*node, 0, len(p.parsers)),
		toks:     make([]*token, 0, cap(p.result.node.toks)),
	}
	p.result.unparsed.clear()
	p.result.unparsed.setTrace(p.trace)
	p.tokenStack.clear()
	p.tokenStack.setTrace(p.trace)
	p.currentParser = nil
	p.parserIndex = 0
	p.cacheChecked = false
	p.initEvaluated = false
	p.result.node.nodes = p.result.node.nodes[:0]
	p.skip = 0
	p.skippingAfterDone = false
}

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

		if !p.cacheChecked {
			p.cacheChecked = true

			p.cacheToken = t
			if p.initNode != zeroNode {
				p.cacheToken = p.initNode.token
			}

			if p.cache.has(t.offset) {
				p.trace.out("found in cache, valid:", p.result.valid)
				p.result = p.cache.get()
				p.result.unparsed.append(t)
				return p.result
			}
		}

		if p.currentParser == nil {
			if p.initNode == zeroNode || p.initEvaluated {
				p.currentParser = p.parsers[p.parserIndex]
				p.currentParser.init(p.trace, zeroNode)
			} else if p.parserIndex < len(p.initParsers) {
				p.currentParser = p.initParsers[p.parserIndex]
				if p.currentParser != nil {
					p.currentParser.init(p.trace, p.initNode)
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
				p.trace.out("group done, valid")
				p.result.valid = true

				p.cacheToken = p.result.node.token
				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.registry.typeName(p.result.node.nodeType)))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cache.set(p.cacheToken.offset, p.result)

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
				p.trace.out("group done, valid")
				p.result.valid = true

				p.cacheToken = p.result.node.token
				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.registry.typeName(p.result.node.nodeType)))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cache.set(p.cacheToken.offset, p.result)

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

		p.cacheToken = p.result.node.token
		if p.cacheToken == nil {
			if !p.tokenStack.has() {
				panic(unexpectedResult(p.registry.typeName(p.result.node.nodeType)))
			}

			p.cacheToken = p.tokenStack.peek()
		}

		p.cache.set(p.cacheToken.offset, p.result)

		p.result.unparsed.merge(p.tokenStack)
		if p.result.node.len() > p.initNode.len() {
			p.result.unparsed.mergeTokens(p.result.node.tokens()[p.initNode.len():])
		}

		return p.result
	}
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

func (g *unionGenerator) create(init nodeType, excluded typeList) (*generatorResult, error) {
	if p, ok := g.registry.getParser(g.typ, init, excluded); ok {
		return p, nil
	}

	p := &generatorResult{parser: &unionParser{}, preallocated: true}
	g.registry.setParser(g.typ, init, excluded, p)

	if err := g.checkExpand(); err != nil {
		return nil, err
	}

	expandedTypes := make([]nodeType, len(g.elements))
	for i, e := range g.elements {
		expandedTypes[i] = e.nodeType()
	}

	var expectedLength int
	parsers := make([][]parser, len(g.elements)+1) // TODO: check the index shift
	for i, it := range append([]nodeType{init}, expandedTypes...) {
		p := make([]parser, len(g.elements))
		for j, e := range g.elements {
			gr, err := e.create(it, excluded)
			if err != nil {
				return nil, err
			}

			if gr.valid || gr.preallocated {
				if gr.preallocated {
					gr.required = true
				}

				p[j] = gr.parser
				if gr.expectedLength > expectedLength {
					expectedLength = gr.expectedLength
				}
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
		if p.required {
			return nil, requiredParserInvalid(g.registry.typeName(g.typ))
		}

		p.preallocated = false
		p.parser = nil
		return p, nil
	}

	p.valid = true
	p.parser.(*unionParser).registry = g.registry
	p.parser.(*unionParser).nodeType = g.typ
	p.parser.(*unionParser).parsers = parsers
	p.parser.(*unionParser).initIsMember = initIsMember
	p.parser.(*unionParser).result = &parserResult{
		unparsed: newTokenStack(expectedLength),
		node:     zeroNode,
	}
	p.parser.(*unionParser).tokenStack = newTokenStack(expectedLength)
	p.parser.(*unionParser).cache = &cache{}
	p.expectedLength = expectedLength
	p.preallocated = false
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

func (p *unionParser) init(t trace, n *node) {
	p.trace = t.extend(p.nodeType)
	p.initNode = n
	if p.initIsMember {
		p.result.node = n
	} else {
		p.result.node = zeroNode
	}

	p.tokenStack.clear()
	p.tokenStack.setTrace(p.trace)
	p.result.unparsed.clear()
	p.result.unparsed.setTrace(p.trace)
	p.initTypeIndex = 0
	p.parserIndex = 0
	p.parsers[0][0].init(p.trace, n)
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

			if p.cache.has(p.cacheToken.offset) {
				p.trace.out("found in cache, valid:", p.result.valid)
				p.result = p.cache.get()
				p.result.unparsed.append(t)
				p.result.accepting = false
				return p.result
			}
		}

		p.itemResult = p.parsers[p.initTypeIndex][p.parserIndex].parse(t)
		if p.itemResult.accepting {
			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		p.tokenStack.merge(p.itemResult.unparsed)

		if !p.itemResult.valid {
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
				p.trace.out("done, valid:", p.result.node != zeroNode)
				p.result.accepting = false
				p.result.valid = p.result.node != zeroNode

				p.cacheToken = p.result.node.token

				if p.cacheToken == nil {
					p.cacheToken = p.initNode.token
				}

				if p.cacheToken == nil {
					if !p.tokenStack.has() {
						panic(unexpectedResult(p.registry.typeName(p.result.node.nodeType)))
					}

					p.cacheToken = p.tokenStack.peek()
				}

				p.cache.set(p.cacheToken.offset, p.result)

				p.result.unparsed.merge(p.tokenStack)
				return p.result
			}

			p.parsers[p.initTypeIndex][p.parserIndex].init(p.trace, p.result.node)

			if p.tokenStack.has() {
				t = p.tokenStack.pop()
				continue parseLoop
			}

			p.result.accepting = true
			return p.result
		}

		if p.result.node == zeroNode || p.result.node.len() < p.itemResult.node.len() {
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
			p.parsers[p.initTypeIndex][p.parserIndex].init(p.trace, p.result.node)

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
				panic(unexpectedResult(p.registry.typeName(p.result.node.nodeType)))
			}

			p.cacheToken = p.tokenStack.peek()
		}

		p.cache.set(p.cacheToken.offset, p.result)

		p.result.unparsed.merge(p.tokenStack)

		if p.skip > 0 {
			p.skippingAfterDone = true
			p.result.accepting = true
			return p.result
		}

		return p.result
	}
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
	s.registry.clearParsers()

	root, err := s.registry.root()
	if err != nil {
		return zeroNode, err
	}

	trace := newTrace(s.traceLevel, s.registry)

	gr, err := root.create(0, nil)
	if err != nil {
		return zeroNode, err
	}

	if !gr.valid {
		panic(errFailedToCreateParser)
	}

	parser := gr.parser
	parser.init(trace, zeroNode)

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

			return last.node, nil
		}

		if err == io.EOF {
			if !last.accepting {
			}

			last = parser.parse(eofToken)

			if !last.valid {
				return zeroNode, errUnexpectedEOF
			}

			if !last.unparsed.has() || last.unparsed.peek() != eofToken {
				return zeroNode, errUnexpectedEOF
			}

			return last.node, nil
		}

		last = parser.parse(&t)
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
