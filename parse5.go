package mml

import (
	"errors"
	"fmt"
	"log"
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
	nodes    []*node
	toks     []*token
}

type generatorResult struct {
	valid          bool
	parser         parser
	expectedLength int
}

type parseResult struct {
	accepting bool
	valid     bool
	node      *node
	fromCache bool
	unparsed  *tokenStack
}

type parser interface {
	init(n *node)
	parse(*token) *parseResult // optimization: should not make memory allocations
}

type generator interface {
	create(trace, nodeType, typeList) (*generatorResult, error)
	member(nodeType) (bool, error)
}

type trace interface {
	extend(nodeType) trace
	outLevel(traceLevel, ...interface{})
	out(...interface{})
	debug(...interface{})
}

type registry interface {
	typeName(n nodeType) string
	nodeType(typeName string) nodeType
	get(nodeType) (generator, bool)
	root() (generator, error)
	primitive(string, tokenType) error
	optional(string, string) error
	sequence(string, string) error
	group(string, ...string) error
	union(string, ...string) error
}

// TODO: use this for the unparsed
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
	result *parseResult
	offset int
}

type parserRegistry struct {
	idSeed        nodeType
	typeIDs       map[string]nodeType
	typeNames     map[nodeType]string
	generators    map[nodeType]generator
	rootGenerator generator
}

type primitiveGenerator struct {
	nodeType  nodeType
	tokenType tokenType
}

type primitiveParser struct {
	trace     trace
	nodeType  nodeType
	tokenType tokenType
	result    *parseResult
	cache     *cache
}

type optionalGenerator struct {
	nodeType nodeType
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
	result         *parseResult
	optionalResult *parseResult
	cacheToken     *token
	cache          *cache
	cacheChecked   bool
}

type sequenceGenerator struct {
	nodeType nodeType
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
	result            *parseResult
	skip              int
	currentParser     parser
	cacheChecked      bool
	cacheToken        *token
	itemResult        *parseResult
	tokenStack        *tokenStack
	initEvaluated     bool
	cache             *cache
	skippingAfterDone bool
}

type groupGenerator struct {
	nodeType nodeType
	items    []nodeType
	registry registry
}

type groupParser struct {
	trace             trace
	registry          registry
	parsers           []parser
	currentParser     parser
	initIsMember      bool
	result            *parseResult
	tokenStack        *tokenStack
	cache             *cache
	initNode          *node
	skip              int
	parserIndex       int
	cacheChecked      bool
	initEvaluated     bool
	skippingAfterDone bool
	cacheToken        *token
	itemResult        *parseResult
}

var (
	errUnexpectedInitNode = errors.New("unexpected init node")
	errNoParsersDefined   = errors.New("no parser defined")

	zeroNode = &node{}
)

func unspecifiedParser(typeName string) error {
	return fmt.Errorf("unspecified parser: %s", typeName)
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

func groupWithoutItems(nodeType string) error {
	return fmt.Errorf("group without items: %s", nodeType)
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
	if len(n.toks) == 1 {
		n.token = n.toks[0]
	}
}

func (n *node) clearNodes() {
	n.nodes = nil
	n.toks = nil
}

func newTokenStack(t trace, expectedSize int) *tokenStack {
	return &tokenStack{
		trace: t,
		stack: make([]*token, 0, expectedSize),
	}
}

func (s *tokenStack) append(t *token) {
	if len(s.stack) == cap(s.stack) {
		s.trace.debug("token stack exceeded expected size")
	}

	s.stack = append(s.stack, t)
}

func (s *tokenStack) merge(from *tokenStack) {
	s.need = len(s.stack) + len(from.stack) - cap(s.stack)
	if s.need > 0 {
		s.trace.debug("token stack exceeded expected size")

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

func newTrace(r registry, l traceLevel, root nodeType) *parserTrace {
	return &parserTrace{
		registry: r,
		level:    l,
		path:     "/" + r.typeName(root),
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

	log.Printf("%s: ", t.path)
	log.Println(a...)
}

func (t *parserTrace) out(a ...interface{}) {
	t.outLevel(traceOn, a...)
}

func (t *parserTrace) debug(a ...interface{}) {
	t.outLevel(traceDebug, a...)
}

func (c *cache) set(offset int, r *parseResult) {
	r.fromCache = true
	c.offset = offset
	c.result = r
}

func (c *cache) has(offset int) bool {
	if offset != c.offset {
		return false
	}

	return c.result != nil
}

func (c *cache) get() *parseResult {
	return c.result
}

func newRegistry() *parserRegistry {
	return &parserRegistry{
		typeIDs:    make(map[string]nodeType),
		typeNames:  make(map[nodeType]string),
		generators: make(map[nodeType]generator),
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
	return t
}

func (r *parserRegistry) typeName(t nodeType) string {
	return r.typeNames[t]
}

func (r *parserRegistry) get(t nodeType) (generator, bool) {
	g, ok := r.generators[t]
	return g, ok
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
		nodeType:  r.nodeType(typeName),
		tokenType: t,
	}

	return r.register(g.nodeType, g)
}

func (r *parserRegistry) optional(typeName string, optional string) error {
	g := &optionalGenerator{
		nodeType: r.nodeType(typeName),
		optional: r.nodeType(optional),
		registry: r,
	}

	return r.register(g.nodeType, g)
}

func (r *parserRegistry) sequence(typeName string, itemType string) error {
	g := &sequenceGenerator{
		nodeType: r.nodeType(typeName),
		item:     r.nodeType(itemType),
		registry: r,
	}

	return r.register(g.nodeType, g)
}

func (r *parserRegistry) group(typeName string, itemTypes ...string) error {
	items := make([]nodeType, len(itemTypes))
	for i, t := range itemTypes {
		items[i] = r.nodeType(t)
	}

	g := &groupGenerator{
		nodeType: r.nodeType(typeName),
		items:    items,
		registry: r,
	}

	return r.register(g.nodeType, g)
}

func (r *parserRegistry) union(typeName string, elementTypes ...string) error {
	return nil
}

func (g *primitiveGenerator) create(t trace, init nodeType, excluded typeList) (*generatorResult, error) {
	if excluded.contains(g.nodeType) || init != 0 {
		return &generatorResult{}, nil
	}

	t = t.extend(g.nodeType)
	n := &node{
		nodeType: g.nodeType,
		toks:     []*token{nil},
	}

	return &generatorResult{
		valid: true,
		parser: &primitiveParser{
			trace:     t,
			nodeType:  g.nodeType,
			tokenType: g.tokenType,
			result: &parseResult{
				node:     n,
				unparsed: newTokenStack(t, 1),
			},
			cache: &cache{},
		},
		expectedLength: 1,
	}, nil
}

func (g *primitiveGenerator) member(t nodeType) (bool, error) {
	return t == g.nodeType, nil
}

func (p *primitiveParser) init(n *node) {
	if n != nil {
		panic(errUnexpectedInitNode)
	}

	p.result.unparsed.clear()
}

func (p *primitiveParser) parse(t *token) *parseResult {
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
		p.result.unparsed.append(t)
	}

	p.cache.set(t.offset, p.result)
	return p.result
}

func (g *optionalGenerator) create(t trace, init nodeType, excluded typeList) (*generatorResult, error) {
	optional, ok := g.registry.get(g.optional)
	if !ok {
		return nil, unspecifiedParser(g.registry.typeName(g.optional))
	}

	if m, err := optional.member(g.nodeType); err != nil {
		return nil, err
	} else if m {
		return nil, optionalContainingSelf(g.registry.typeName(g.optional))
	}

	if excluded.contains(g.nodeType) {
		return &generatorResult{}, nil
	}

	t = t.extend(g.nodeType)
	excluded = append(excluded, g.nodeType)
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

	return &generatorResult{
		valid: true,
		parser: &optionalParser{
			trace:        t,
			registry:     g.registry,
			nodeType:     g.nodeType,
			optional:     optParser.parser,
			initIsMember: initIsMember,
			result: &parseResult{
				node:     zeroNode,
				unparsed: newTokenStack(t, optParser.expectedLength),
			},
			cache: &cache{},
		},
		expectedLength: optParser.expectedLength,
	}, nil
}

func (g *optionalGenerator) member(t nodeType) (bool, error) {
	optional, ok := g.registry.get(g.optional)
	if !ok {
		return false, unspecifiedParser(g.registry.typeName(g.optional))
	}

	return optional.member(t)
}

func (p *optionalParser) init(n *node) {
	p.initNode = n
	p.result.unparsed.clear()
	p.cacheChecked = false
	p.optional.init(n)
}

func (p *optionalParser) parse(t *token) *parseResult {
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

	p.optionalResult = p.optional.parse(t)
	if p.optionalResult.accepting {
		p.result.accepting = true
		return p.result
	}

	p.result.accepting = false
	p.result.valid = true
	p.result.unparsed.merge(p.optionalResult.unparsed)

	if p.optionalResult.valid {
		p.trace.out("parse done, valid:", p.result.valid)
		p.result.node = p.optionalResult.node
		p.result.fromCache = p.optionalResult.fromCache
	} else if p.initIsMember {
		p.trace.out("init node is a member, valid")
		p.result.node = p.initNode
		p.result.fromCache = false
	} else {
		p.result.node = zeroNode
		p.trace.out("missing optional, valid")
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

func (g *sequenceGenerator) create(t trace, init nodeType, excluded typeList) (*generatorResult, error) {
	item, ok := g.registry.get(g.item)
	if !ok {
		return nil, unspecifiedParser(g.registry.typeName(g.item))
	}

	if m, err := g.member(g.item); err != nil {
		return nil, err
	} else if m {
		return nil, sequenceContainingSelf(g.registry.typeName(g.nodeType))
	}

	if excluded.contains(g.nodeType) {
		return &generatorResult{}, nil
	}

	t = t.extend(g.nodeType)
	allExcluded := append(excluded, g.nodeType)
	selfExcluded := typeList{g.nodeType}

	first, err := item.create(t, init, allExcluded)
	if err != nil {
		return nil, err
	}

	rest, err := item.create(t, 0, selfExcluded)
	if err != nil {
		return nil, err
	}

	var initIsMember bool
	if init != 0 {
		if m, err := item.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	if !first.valid && !initIsMember {
		return &generatorResult{}, nil
	}

	expectedLength := first.expectedLength
	if rest.expectedLength > expectedLength {
		expectedLength = rest.expectedLength
	}

	return &generatorResult{
		parser: &sequenceParser{
			trace:        t,
			registry:     g.registry,
			first:        first.parser,
			rest:         rest.parser,
			initIsMember: initIsMember,
			result: &parseResult{
				node: &node{
					nodeType: g.nodeType,
				},
				unparsed: newTokenStack(t, expectedLength),
			},
			tokenStack: newTokenStack(t, expectedLength),
			cache:      &cache{},
		},
		expectedLength: expectedLength,
	}, nil
}

func (g *sequenceGenerator) member(t nodeType) (bool, error) {
	return t == g.nodeType, nil
}

func (p *sequenceParser) init(n *node) {
	p.initNode = n
	p.cacheChecked = false
	p.currentParser = p.first
	p.currentParser.init(n)
	p.result.node.nodes = nil
	p.skip = 0
	p.result.unparsed.clear()
	p.tokenStack.clear()
	p.initEvaluated = false
	p.skippingAfterDone = false
}

func (p *sequenceParser) parse(t *token) *parseResult {
	p.trace.out("parsing", t)

parseLoop:
	for {
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
		p.currentParser.init(nil)

		if p.itemResult.valid && p.itemResult.node != zeroNode {
			p.initEvaluated = true
			p.result.node.append(p.itemResult.node)
			if p.itemResult.fromCache {
				p.skip = p.tokenStack.findCachedNode(p.itemResult.node)
			}

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

func (g *groupGenerator) create(t trace, init nodeType, excluded typeList) (*generatorResult, error) {
	if len(g.items) == 0 {
		return nil, groupWithoutItems(g.registry.typeName(g.nodeType))
	}

	if excluded.contains(g.nodeType) {
		return &generatorResult{}, nil
	}

	items := make([]generator, len(g.items))
	for i, it := range g.items {
		gi, ok := g.registry.get(it)
		if !ok {
			return nil, unspecifiedParser(g.registry.typeName(it))
		}

		items[i] = gi
	}

	t = t.extend(g.nodeType)
	excluded = append(excluded, g.nodeType)

	first, err := items[0].create(t, init, excluded)
	if err != nil {
		return nil, err
	}

	expectedLength := first.expectedLength

	parsers := make([]parser, len(items))
	parsers[0] = first.parser
	for i, ig := range items[1:] {
		gr, err := ig.create(t, 0, nil)
		if err != nil {
			return nil, err
		}

		parsers[i+1] = gr.parser
		expectedLength += gr.expectedLength
	}

	var initIsMember bool
	if init != 0 {
		if m, err := g.member(init); err != nil {
			return nil, err
		} else {
			initIsMember = m
		}
	}

	if !first.valid && !initIsMember {
		return &generatorResult{}, nil
	}

	return &generatorResult{
		valid: true,
		parser: &groupParser{
			trace:        t,
			registry:     g.registry,
			parsers:      parsers,
			initIsMember: initIsMember,
			result: &parseResult{
				node: &node{
					nodeType: g.nodeType,
				},
				unparsed: newTokenStack(t, expectedLength),
			},
			tokenStack: newTokenStack(t, expectedLength),
			cache:      &cache{},
		},
		expectedLength: expectedLength,
	}, nil
}

func (g *groupGenerator) member(t nodeType) (bool, error) {
	return t == g.nodeType, nil
}

func (p *groupParser) init(n *node) {
	p.initNode = n
	p.result.unparsed.clear()
	p.tokenStack.clear()
	p.currentParser = nil
	p.parserIndex = 0
	p.cacheChecked = false
	p.initEvaluated = false
	p.result.node.nodes = nil
	p.skip = 0
	p.skippingAfterDone = false
}

func (p *groupParser) parse(t *token) *parseResult {
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

			// NOTE: there was an init node check here in parse4
		}

		if p.currentParser == nil {
			p.currentParser = p.parsers[p.parserIndex]
			if p.parserIndex == 0 {
				p.currentParser.init(p.initNode)
			} else {
				p.currentParser.init(nil)
			}

			p.parserIndex++
		}

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

		p.currentParser = nil
		p.tokenStack.merge(p.itemResult.unparsed)

		if p.itemResult.valid {
			p.initEvaluated = true
			p.result.node.append(p.itemResult.node)
			if p.itemResult.fromCache {
				p.skip = p.tokenStack.findCachedNode(p.itemResult.node)
				// TODO: skip should be always discarded
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

		if !p.initEvaluated {
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
