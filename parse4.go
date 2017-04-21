/*
[generator]
- created by the syntax definition
- the syntax definition accepts the node type and, in case of complex parsers, a list of other node types
- creates a parser instance for a concrete node
- before creating, it can tell whether it can create a parser with the provided init node and excluded node
  types
- create accepts an optional init node and a list of excluded nodes
- can tell whether a node type is a member of the type of the generator
- can create and create return an error when a referenced generator is not defined by the syntax

[parser]
- tries parsing a concrete node
- returns whether it accepts more tokens or it is done
- when done it can have a valid or invalid result
- when done it returns the unparsed tokens
- when done with a valid result, returns the parsed node
- on init, it accepts an init node
- on init, it accepts a set of excluded node types
- if it has a valid result for the given token position in the cache, returns that
- if it knows that it cannot be valid for the given token position, returns invalid
- when the parse was successful for the given token position, caches the result
- when the parse was unsuccessful for the given token position, caches the result, but only if there was no
  successful parse before, because this can happen when parsing further
- parse returns an error when it detects that the syntax is invalid
- calling parse after inidicating done, causes a panic
- returns how many tokens were taken from the cache in addition to the provided ones minus the unparsed
- in complex parsers, the can create check of the items must happen before checking the cache
- every complex parser sets the token of the returned node to the token of the first item node
- token position for the cache considers empty sequences

[init item]
- a node that is already parsed at the given token position and can be used as the initial segment of a more
  complex node

[excluded node types]
- nodes that are being tried at the given token position on a higher level in the parser tree

[parse result]
- tells whether the parser is accepting more tokens
- if not, tells whether the node is valid, and contains the parsed node and the unparsed tokens
- if an item was taken from the cache, it tells how many tokens more were taken from the cache than accepted

[primitive generator]
- cannot create a parser when its node type is excluded
- cannot create a parser when it is supplied with an init node, and the init node is of a different type
- creates a primitive parser with its name, expected token type and init node

[primitive parser]
- on init, expects the node type, and either the token type or an init node
- when it has an init node, it is automatically valid, and doesn't accept more tokens
- when it doesn't have an init node, it accepts a single token of the provided token type
- possible valid results: the node of the token
- it always returns its own node type
- no need to cache it

[optional generator]
- returns an error when the optional generator is not defined in the syntax
- returns false if it is excluded
- returns the result of the optional generator otherwise
- extends the excluded with itself
- cannot contain itself

[optional parser]
- on init, expects the node type, the generator of the optional node, an optional init node and the list of
  excluded nodes
- on init, it adds itself to the excluded list
- always returns a valid result
- when the parse of the optional node failed, returns a valid result with a zero node and all the tokens passed
  in
- when the parse of the optional node succeeded, returns the result of the optional parser
- it never returns its own node type
- if the result is empty, the first unparsed token is used as the node token
- possible valid results: the optional node or a zero node

[sequence generator]
- returns an error when the item generator is not defined by the syntax
- can create returns false if the sequence is excluded
- extends the excluded with itself
- returns true if the init item is a member type and is not excluded
- returns the result of the item generator
- cannot contain itself

[sequence parser]
- on init, expects the node type, an item generator, an optional init node and the list of excluded nodes
- the init node is considered an item
- always returns a valid result
- when the parse of an item failed, returns the existing items
- when the parse of an item succeeded, stores it, queues the unparsed tokens, and tries to parse the next item
- the init item is only used with the first item
- when an item from the cache has more read ahead than tokens in the queue, it ignores the right amount of
  tokens before continuing with the next item
- when there is an init item, it's token is used to check whether there is a cached result
- in case of the first item, it uses the excluded types and init node to initialize the item, for the rest it
  uses only itself as the excluded type and the zero node
- the unparsed tokens are stored in a queue, returned as unparsed when done
- possible valid results: an empty node, or a node with item nodes
- it always returns its own node type
- parses only the first node with the init item
- if there is an init item, and the parse of the first node fails, and the init item is a member of the item,
  then it is added as the first node
- if the result is empty, the first unparsed token is used as the node token
- it returns also if zero
- TODO: what if the init item can be an element? try to do the same as in the group

[group generator]
- returns an error if any of the items is not defined by the syntax
- returns false if it is excluded
- returns an error if it doesn't have items
- extends the excluded with itself
- returns true if the first item returns true
- returns true if it has an init item and it can be the first item
- can contain itself

[group parser]
- on init, it expects the node type, the generators of the group items, an optional init node and a list of the
  excluded node types
- the init node is considered the first item or the init node of the first item
- it always uses the next generator for the next item. When there are no more generators, the parse is
  successful
- the unparsed tokens are stored in a queue, returned as unparsed when done
- when creating the parser of the first item, it passes in the init item and the excluded types. For the rest of
  the items, no init item and no excluded types are passed in.
- if the parse of the first item failed, it checks if the init item can be used as the first item, and if yes,
  continues with next item, otherwise it fails
- if the parse of an item fails, it fails
- on failure, it returns the tokens of the parsed items, the unparsed tokens and the tokens in the queue
- if the parse of an item succeeds, it appends the node to its nodes, and continues with the next item
- possible valid results: the group node with the non-zero items
- it always returns its own node type

[union generator]
- it expands the unions in the union for the actual items
- returns an error if any of the items is not defined in the syntax
- returns an error if it doesn't have items
- can contain itself, but it's ignored
- returns true if any of the generators return true
- returns the generators that return true

[union parser]
- on init, it expects the node type, an optional init node, the element generators and a list of the excluded
  node types
- the init node is considered an element or an init node to the elements
- possible valid results: the node of the matching element, can be zero from optional
- when the element parsing hasn't started, it tries to find a generator that accepts the current init item and
  the set of excluded types and parses the item with that
- when the element parser failed, tries the next generator
- when an element parser succeeded tries all the generators again, for a result that consumes more tokens than
  the last successful element
- it never returns its own node type
- TODO: should be able to use the init node as an element

[errors]
- TODO: errors coming from invalid syntax specification
- TODO: errors coming from invalid syntax

[tracing]
*/
package mml

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"strings"
// )
//
// type idSet struct {
// 	buckets []uint64
// }
//
// type node struct {
// 	*token
// 	typ    string
// 	typeID uint64
// 	nodes  []*node
// }
//
// type generator interface {
// 	setTypeID(uint64)
// 	create(t trace, init uint64, excludedTypes *idSet) (parserRegistry, bool, error)
// 	member(typeID uint64) (bool, error)
// 	cacheMembers() error
// }
//
// type parserRegistry interface {
// 	get(init *node) (parser, bool)
// }
//
// type parseResult struct {
// 	accepting bool
// 	valid     bool
// 	unparsed  []*token
// 	fromCache bool
// 	node      *node
// }
//
// type parser interface {
// 	initialize(init *node)
// 	parse(t *token) (parseResult, error)
// }
//
// type parserReg struct {
// 	parsers []parser
// }
//
// type cacheItem struct {
// 	match   []*node
// 	noMatch *idSet
// }
//
// type tokenCache struct {
// 	tokens []*cacheItem
// }
//
// type traceLevel int
//
// const (
// 	traceOff traceLevel = iota
// 	traceOn
// 	traceDebug
// )
//
// type trace struct {
// 	level      traceLevel
// 	pathString string
// }
//
// type baseParser struct {
// 	trace         trace
// 	typeID        uint64
// 	init          *node
// 	excludedTypes *idSet
// 	done          bool
// 	skip          int
// 	typeName      string
// }
//
// type backtrackingParser struct {
// 	baseParser
// 	queue         []*token
// 	initEvaluated bool
// 	cacheChecked  bool
// 	node          *node
// }
//
// type collectionParser struct {
// 	backtrackingParser
// 	parser         parser
// 	firstGenerator generator
// }
//
// type primitiveGenerator struct {
// 	nodeType  string
// 	typeID    uint64
// 	tokenType tokenType
// 	parser    *primitiveParser
// 	typeName  string
// }
//
// type primitiveParser struct {
// 	baseParser
// 	token tokenType
// }
//
// type optionalGenerator struct {
// 	nodeType          string
// 	typeID            uint64
// 	optional          uint64
// 	optionalGenerator generator
// 	members           *idSet
// }
//
// type optionalParser struct {
// 	baseParser
// 	optional       generator
// 	optionalParser parserRegistry
// }
//
// type sequenceGenerator struct {
// 	nodeType string
// 	typeID   uint64
// 	itemType uint64
// 	item     generator
// }
//
// type sequenceParser struct {
// 	collectionParser
// 	generator generator
// }
//
// type groupGenerator struct {
// 	nodeType       string
// 	typeID         uint64
// 	itemTypes      []uint64
// 	itemGenerators []generator
// 	first          generator
// }
//
// type groupParser struct {
// 	collectionParser
// 	items     []generator
// 	firstItem generator
// }
//
// type unionGenerator struct {
// 	nodeType     string
// 	typeID       uint64
// 	elementTypes []uint64
// 	expanded     []generator
// 	members      *idSet
// }
//
// type unionParser struct {
// 	backtrackingParser
// 	elements       []generator
// 	activeElements []generator
// 	parser         parser
// 	valid          bool
// }
//
// var (
// 	generators              = make(map[uint64]generator)
// 	generatorsByName        = make(map[string]generator)
// 	typeNames               = make(map[uint64]string)
// 	typeIDs                 = make(map[string]uint64)
// 	isSep                   func(*node) bool
// 	postParse                     = make(map[string]func(*node) *node)
// 	cache                         = &tokenCache{}
// 	eofToken                      = &token{offset: -1}
// 	zeroNode                *node = nil
// 	voidResult                    = parseResult{}
// 	errInvalidRootGenerator       = errors.New("invalid root generator")
// 	errUnexpectedEOF              = errors.New("unexpected EOF")
// 	typeIDFeed              uint64
// )
//
// func init() {
// 	registerType("zero")
// }
//
// func (r *parserReg) get(init *node) (parser, bool) {
// 	var id int
// 	if init != nil {
// 		id = int(init.typeID)
// 	}
//
// 	if len(r.parsers) <= id {
// 		return nil, false
// 	}
//
// 	p := r.parsers[id]
// 	return p, p != nil
// }
//
// func cacheMembers() {
// 	for _, g := range generators {
// 		if err := g.cacheMembers(); err != nil {
// 			panic(err)
// 		}
// 	}
// }
//
// func stringsContain(strc []string, stri ...string) bool {
// 	for _, sc := range strc {
// 		for _, si := range stri {
// 			if si == sc {
// 				return true
// 			}
// 		}
// 	}
//
// 	return false
// }
//
// func intsContain(ic []int, i ...int) bool {
// 	for _, ici := range ic {
// 		for _, ii := range i {
// 			if ii == ici {
// 				return true
// 			}
// 		}
// 	}
//
// 	return false
// }
//
// func (s *idSet) add(id uint64) {
// 	bucket, flag := id/64, id%64
// 	need := int(bucket + 1)
// 	if len(s.buckets) < need {
// 		if cap(s.buckets) >= need {
// 			s.buckets = s.buckets[:need]
// 		} else {
// 			s.buckets = s.buckets[:cap(s.buckets)]
// 			for len(s.buckets) < need {
// 				s.buckets = append(s.buckets, 0)
// 			}
// 		}
// 	}
//
// 	s.buckets[bucket] |= 1 << flag
// }
//
// func (s *idSet) clone() *idSet {
// 	buckets := make([]uint64, len(s.buckets), cap(s.buckets))
// 	copy(buckets, s.buckets)
// 	return &idSet{buckets: buckets}
// }
//
// func (s *idSet) extend(ids ...uint64) *idSet {
// 	var se *idSet
// 	if s == nil {
// 		se = &idSet{}
// 	} else {
// 		se = s.clone()
// 	}
//
// 	for _, id := range ids {
// 		se.add(id)
// 	}
//
// 	return se
// }
//
// func (s *idSet) contains(id uint64) bool {
// 	if s == nil {
// 		return false
// 	}
//
// 	bucket, flag := id/64, id%64
// 	if uint64(len(s.buckets)) <= bucket {
// 		return false
// 	}
//
// 	return s.buckets[bucket]&(1<<flag) != 0
// }
//
// func (c *tokenCache) getMatch(t *token, id uint64) (*node, bool) {
// 	if t.offset < 0 {
// 		return nil, false
// 	}
//
// 	offset := t.offset
//
// 	if len(c.tokens) <= offset || c.tokens[offset] == nil {
// 		return nil, false
// 	}
//
// 	ci := c.tokens[offset]
// 	if ci.noMatch.contains(id) {
// 		return nil, false
// 	}
//
// 	index := int(id)
// 	if len(ci.match) <= index || ci.match[index] == nil {
// 		return nil, false
// 	}
//
// 	return ci.match[index], true
// }
//
// func (c *tokenCache) hasNoMatch(t *token, id uint64) bool {
// 	if t.offset < 0 {
// 		return false
// 	}
//
// 	offset := t.offset
//
// 	if len(c.tokens) <= offset || c.tokens[offset] == nil {
// 		return false
// 	}
//
// 	ci := c.tokens[offset]
// 	return ci.noMatch.contains(id)
// }
//
// func (c *tokenCache) setMatch(t *token, id uint64, n *node) {
// 	if t.offset < 0 {
// 		return
// 	}
//
// 	offset := t.offset
//
// 	need := offset + 1
// 	if len(c.tokens) < need {
// 		if cap(c.tokens) >= need {
// 			c.tokens = c.tokens[:need]
// 		} else {
// 			c.tokens = c.tokens[:cap(c.tokens)]
// 			for len(c.tokens) < need {
// 				c.tokens = append(c.tokens, nil)
// 			}
// 		}
// 	}
//
// 	ci := c.tokens[offset]
// 	if ci == nil {
// 		ci = &cacheItem{noMatch: &idSet{}}
// 		c.tokens[offset] = ci
// 	}
//
// 	need = int(id) + 1
// 	if len(ci.match) < need {
// 		if cap(ci.match) >= need {
// 			ci.match = ci.match[:need]
// 		} else {
// 			ci.match = ci.match[:cap(ci.match)]
// 			for len(ci.match) < need {
// 				ci.match = append(ci.match, nil)
// 			}
// 		}
// 	}
//
// 	ci.match[id] = n
// }
//
// func (c *tokenCache) setNoMatch(t *token, id uint64) {
// 	if t.offset < 0 {
// 		return
// 	}
//
// 	offset := t.offset
//
// 	need := offset + 1
// 	if len(c.tokens) < need {
// 		if cap(c.tokens) >= need {
// 			c.tokens = c.tokens[:need]
// 		} else {
// 			c.tokens = c.tokens[:cap(c.tokens)]
// 			for len(c.tokens) < need {
// 				c.tokens = append(c.tokens, nil)
// 			}
// 		}
// 	}
//
// 	ci := c.tokens[offset]
// 	if ci == nil {
// 		ci = &cacheItem{noMatch: &idSet{}}
// 		c.tokens[offset] = ci
// 	}
//
// 	// reasonable to leak this check over to here:
// 	// a shorter variant may already have been parsed
// 	index := int(id)
// 	if len(ci.match) >= index+1 && ci.match[index] != nil {
// 		return
// 	}
//
// 	ci.noMatch.add(id)
// }
//
// func newTrace(l traceLevel) trace {
// 	return trace{level: l}
// }
//
// func (t trace) extend(nodeType string) trace {
// 	et := trace{
// 		level: t.level,
// 	}
//
// 	if et.level > traceOff {
// 		et.pathString = t.pathString + "/" + nodeType
// 	}
//
// 	return et
// }
//
// func (t trace) outLevel(l traceLevel, a ...interface{}) {
// 	if l > t.level {
// 		return
// 	}
//
// 	if t.pathString == "" {
// 		log.Println(a...)
// 		return
// 	}
//
// 	log.Println(append([]interface{}{t.pathString}, a...)...)
// }
//
// func (t trace) out(a ...interface{}) {
// 	t.outLevel(traceOn, a...)
// }
//
// func (t trace) debug(a ...interface{}) {
// 	t.outLevel(traceDebug, a...)
// }
//
// func (n *node) zero() bool {
// 	return n == nil
// }
//
// func (n *node) tokens() []*token {
// 	if n.zero() {
// 		return nil
// 	}
//
// 	if len(n.nodes) == 0 {
// 		if n.token == nil || n.token.typ == noToken {
// 			return nil
// 		}
//
// 		return []*token{n.token}
// 	}
//
// 	var t []*token
// 	for _, ni := range n.nodes {
// 		t = append(t, ni.tokens()...)
// 	}
//
// 	return t
// }
//
// func (n *node) length() int {
// 	return len(n.tokens())
// }
//
// func (n *node) String() string {
// 	var nodes []string
// 	for _, ni := range n.nodes {
// 		nodes = append(nodes, ni.String())
// 	}
//
// 	return fmt.Sprintf("node:%s:%v(%s)", n.typ, n.token, strings.Join(nodes, ", "))
// }
//
// func registerType(nodeType string) uint64 {
// 	if id, ok := typeIDs[nodeType]; ok {
// 		return id
// 	}
//
// 	id := typeIDFeed
// 	typeIDFeed++
// 	typeIDs[nodeType] = id
// 	typeNames[id] = nodeType
// 	return id
// }
//
// func register(nodeType string, g generator) generator {
// 	id := registerType(nodeType)
// 	g.setTypeID(id)
// 	generators[id] = g
// 	generatorsByName[nodeType] = g
// 	return g
// }
//
// func unexpectedToken(nodeType string, t *token) error {
// 	return fmt.Errorf("unexpected token: %v, %v", nodeType, t)
// }
//
// func unspecifiedParser(typeID uint64) error {
// 	return fmt.Errorf("unspecified parser: %s", typeNames[typeID])
// }
//
// func optionalContainingSelf(nodeType string) error {
// 	return fmt.Errorf("optional containing self: %s", nodeType)
// }
//
// func sequenceContainingSelf(nodeType string) error {
// 	return fmt.Errorf("sequence containing self: %s", nodeType)
// }
//
// func unexpectedResult(nodeType string) error {
// 	return fmt.Errorf("unexpected parse result: %s", nodeType)
// }
//
// func groupWithoutItems(nodeType string) error {
// 	return fmt.Errorf("group without items: %s", nodeType)
// }
//
// func unionWithoutElements(nodeType string) error {
// 	return fmt.Errorf("union without elements: %s", nodeType)
// }
//
// func invalidParserState(nodeType string) error {
// 	return fmt.Errorf("invalid parser state: %s", nodeType)
// }
//
// func primitiveParserNotReady(typ string) error {
// 	return fmt.Errorf("primitive parser not ready: %s", typ)
// }
//
// func (p *baseParser) checkDone(currentToken *token) {
// 	if p.done {
// 		panic(unexpectedToken(typeNames[p.typeID], currentToken))
// 	}
// }
//
// func (p *baseParser) checkSkip() (parseResult, bool) {
// 	if p.skip == 0 {
// 		return voidResult, false
// 	}
//
// 	p.skip--
// 	return parseResult{accepting: true}, true
// }
//
// func (p *backtrackingParser) unparsed(t ...*token) parseResult {
// 	p.trace.out("returning unparsed", len(t), len(p.queue), t, p.queue)
// 	return parseResult{unparsed: append(t, p.queue...)}
// }
//
// func (p *backtrackingParser) abort(err error, unparsed ...*token) (parseResult, error) {
// 	p.trace.out("aborting", unparsed)
// 	return p.unparsed(unparsed...), err
// }
//
// func (p *backtrackingParser) checkCache(t *token) (parseResult, bool) {
// 	// should not get here when parsing from the queue
//
// 	ct := t
// 	if !p.init.zero() {
// 		ct = p.init.token
// 	}
//
// 	if cache.hasNoMatch(ct, p.typeID) {
// 		p.trace.out("no match identified in cache")
// 		return p.unparsed(t), true
// 	}
//
// 	if n, ok := cache.getMatch(ct, p.typeID); ok {
// 		p.trace.out("cached match", ct, p.typeID, "the init:", p.init, n)
// 		return parseResult{
// 			valid:     true,
// 			node:      n,
// 			fromCache: true,
// 			unparsed:  []*token{t},
// 		}, true
// 	}
//
// 	return voidResult, false
// }
//
// func (p *collectionParser) appendNode(n *node) {
// 	if n.zero() {
// 		return
// 	}
//
// 	p.node.nodes = append(p.node.nodes, n)
// 	if len(p.node.nodes) == 1 {
// 		p.node.token = n.token
// 	}
// }
//
// // TODO: is this really required
// func (p *collectionParser) appendInitIfMember() (bool, error) {
// 	if p.init.zero() {
// 		return false, nil
// 	}
//
// 	if m, err := p.firstGenerator.member(p.init.typeID); !m || err != nil {
// 		return m, err
// 	}
//
// 	p.appendNode(p.init)
// 	return true, nil
// }
//
// func (p *collectionParser) appendParsedItem(n *node, fromCache bool) {
// 	p.appendNode(n)
// 	if !fromCache {
// 		return
// 	}
//
// 	t := n.tokens()
// 	for i, ti := range t {
// 		if ti == p.queue[0] {
// 			c := len(t) - i
// 			if c > len(p.queue) {
// 				p.queue, p.skip = nil, c-len(p.queue)
// 			} else {
// 				p.queue = p.queue[len(t):]
// 			}
//
// 			break
// 		}
// 	}
// }
//
// func (p *backtrackingParser) parseNextToken(parser parser) (parseResult, error) {
// 	if len(p.queue) > 0 {
// 		var t *token
// 		t, p.queue = p.queue[0], p.queue[1:]
// 		return parser.parse(t)
// 	}
//
// 	return parseResult{accepting: true}, nil
// }
//
// func primitive(nodeType string, token tokenType) generator {
// 	return register(nodeType, &primitiveGenerator{
// 		nodeType:  nodeType,
// 		typeID:    registerType(nodeType),
// 		typeName:  nodeType,
// 		tokenType: token,
// 		parser: &primitiveParser{
// 			baseParser: baseParser{
// 				typeName: nodeType,
// 			},
// 			token: token,
// 		},
// 	})
// }
//
// func (g *primitiveGenerator) setTypeID(id uint64) {
// 	g.typeID = id
// 	g.parser.typeID = id
// }
//
// func (g *primitiveGenerator) create(t trace, init uint64, excludedTypes *idSet) (parserRegistry, bool, error) {
// 	if excludedTypes.contains(g.typeID) {
// 		return nil, false, nil
// 	}
//
// 	if init != 0 {
// 		return nil, false, nil
// 	}
//
// 	g.parser.trace = t.extend(g.nodeType)
// 	g.parser.done = false
//
// 	return &parserReg{parsers: []parser{g.parser}}, true, nil
// }
//
// func (g *primitiveGenerator) member(typeID uint64) (bool, error) {
// 	return typeID == g.typeID, nil
// }
//
// func (g *primitiveGenerator) cacheMembers() error {
// 	return nil
// }
//
// func (p *primitiveParser) parse(t *token) (parseResult, error) {
// 	p.trace.out("parsing", t)
//
// 	p.checkDone(t)
// 	p.done = true
//
// 	if t.typ != p.token {
// 		p.trace.out("invalid token")
// 		return parseResult{
// 			unparsed: []*token{t},
// 		}, nil
// 	}
//
// 	p.trace.out("valid token")
// 	n := &node{typeID: p.typeID, typ: p.typeName, token: t}
// 	return parseResult{
// 		valid: true,
// 		node:  n,
// 	}, nil
// }
//
// func optional(nodeType string, optionalType string) generator {
// 	return register(nodeType, &optionalGenerator{
// 		nodeType: nodeType,
// 		typeID:   registerType(nodeType),
// 		optional: registerType(optionalType),
// 	})
// }
//
// func (g *optionalGenerator) setTypeID(id uint64) { g.typeID = id }
//
// func (g *optionalGenerator) create(t trace, init uint64, excludedTypes *idSet) (parserRegistry, bool, error) {
// 	optional := g.optionalGenerator
//
// 	if m, err := optional.member(g.typeID); err != nil {
// 		return nil, false, err
// 	} else if m {
// 		return nil, false, optionalContainingSelf(g.nodeType)
// 	}
//
// 	if excludedTypes.contains(g.typeID) {
// 		return nil, false, nil
// 	}
//
// 	optionalParser, ok, err := optional.create(t, init, excludedTypes.extend(g.typeID))
// 	if !ok || err != nil {
// 		return nil, false, err
// 	}
//
// 	parser := newOptionalParser(
// 		t.extend(g.nodeType),
// 		g.typeID,
// 		optional,
// 		optionalParser,
// 		init,
// 		excludedTypes.extend(g.typeID),
// 	)
//
// 	parsers := make([]parser, int(init) - 1)
// 	parsers[int(init) - 1] = parser
// 	return parserReg{parsers: parsers}, true, nil
// }
//
// func (g *optionalGenerator) member(typeID uint64) (bool, error) {
// 	return g.members.contains(typeID), nil
// }
//
// func (g *optionalGenerator) cacheMembers() error {
// 	if g.members != nil {
// 		return nil
// 	}
//
// 	optional, ok := generators[g.optional]
// 	g.optionalGenerator = generators[g.optional]
// 	if !ok {
// 		panic(unspecifiedParser(g.optional))
// 	}
//
// 	g.members = &idSet{}
// 	for id, gi := range generators {
// 		if err := gi.cacheMembers(); err != nil {
// 			return err
// 		}
//
// 		if m, err := optional.member(id); err != nil {
// 			return err
// 		} else if m {
// 			g.members.add(id)
// 		}
// 	}
//
// 	g.members.add(g.typeID)
// 	return nil
// }
//
// func newOptionalParser(
// 	t trace,
// 	typeID uint64,
// 	optional generator,
// 	optionalParser parserRegistry,
// 	init uint64,
// 	excludedTypes *idSet,
// ) *optionalParser {
// 	return &optionalParser{
// 		baseParser: baseParser{
// 			trace:         t,
// 			init:          init,
// 			excludedTypes: excludedTypes.extend(typeID),
// 			typeID:        typeID,
// 		},
// 		optional: optional,
// 		optionalParser: parser,
// 	}
// }
//
// func (p *optionalParser) parse(t *token) (parseResult, error) {
// 	p.trace.out("parsing", t)
// 	p.checkDone(t)
//
// 	if p.optionalParser == nil {
// 		if ok, err := p.optional.canCreate(p.init, p.excludedTypes); !ok || err != nil {
// 			p.trace.out("cannot create optional", p.init)
// 			p.done = true
// 			r := parseResult{unparsed: []*token{t}}
//
// 			if !p.init.zero() {
// 				if m, err := p.optional.member(p.init.typeID); err != nil {
// 					return r, err
// 				} else if m {
// 					p.trace.out("init is a member")
// 					r.node = p.init
// 					r.valid = true
// 					return r, nil
// 				}
// 			}
//
// 			return r, err
// 		}
//
// 		optional, err := p.optional.create(p.trace, p.init, p.excludedTypes)
// 		if err != nil {
// 			p.trace.out("failed to create optional")
// 			p.done = true
// 			return parseResult{unparsed: []*token{t}}, err
// 		}
//
// 		p.optionalParser = optional
// 	}
//
// 	ct := t
// 	if !p.init.zero() {
// 		ct = p.init.token
// 	}
//
// 	if cache.hasNoMatch(ct, p.typeID) {
// 		p.trace.out("cached mismatch")
// 		p.done = true
// 		return parseResult{unparsed: []*token{t}}, nil
// 	}
//
// 	if cn, ok := cache.getMatch(ct, p.typeID); ok {
// 		p.trace.out("cached match")
// 		p.done = true
// 		return parseResult{
// 			valid:     true,
// 			node:      cn,
// 			unparsed:  []*token{t},
// 			fromCache: true,
// 		}, nil
// 	}
//
// 	r, err := p.optionalParser.parse(t)
// 	if err != nil {
// 		p.trace.out("failed to parse optional")
// 		p.done = true
// 		return parseResult{unparsed: []*token{t}}, err
// 	}
//
// 	if r.accepting {
// 		return r, nil
// 	}
//
// 	p.trace.out("optional done, parsed:", r.valid)
// 	p.done = true
//
// 	if r.node.zero() {
// 		if len(r.unparsed) == 0 {
// 			panic(unexpectedResult(typeNames[p.typeID]))
// 		}
//
// 		ct = r.unparsed[0]
// 	} else {
// 		ct = r.node.token
// 	}
//
// 	cache.setMatch(ct, p.typeID, r.node)
// 	r.valid = true
// 	return r, nil
// }
//
// func sequence(nodeType, itemType string) generator {
// 	return register(nodeType, &sequenceGenerator{
// 		nodeType: nodeType,
// 		typeID:   registerType(nodeType),
// 		itemType: registerType(itemType),
// 	})
// }
//
// func (g *sequenceGenerator) setTypeID(id uint64) { g.typeID = id }
//
// func (g *sequenceGenerator) canCreate(init *node, excludedTypes *idSet) (bool, error) {
// 	item, ok := generators[g.itemType]
// 	if !ok {
// 		return false, unspecifiedParser(g.itemType)
// 	}
//
// 	if m, err := item.member(g.typeID); err != nil {
// 		return false, err
// 	} else if m {
// 		return false, sequenceContainingSelf(g.nodeType)
// 	}
//
// 	if excludedTypes.contains(g.typeID) {
// 		return false, nil
// 	}
//
// 	excludedTypes = excludedTypes.extend(g.typeID)
//
// 	if !init.zero() {
// 		if m, err := item.member(init.typeID); err != nil {
// 			return false, err
// 		} else if m && !excludedTypes.contains(init.typeID) {
// 			return true, nil
// 		}
// 	}
//
// 	return item.canCreate(init, excludedTypes)
// }
//
// func (g *sequenceGenerator) create(t trace, init *node, excludedTypes *idSet) (parser, error) {
// 	item := g.item
// 	return newSequenceParser(
// 		t.extend(g.nodeType),
// 		g.typeID,
// 		g.nodeType,
// 		item,
// 		init,
// 		excludedTypes.extend(g.typeID),
// 	), nil
// }
//
// func (g *sequenceGenerator) member(typeID uint64) (bool, error) {
// 	return typeID == g.typeID, nil
// }
//
// func (g *sequenceGenerator) cacheMembers() error {
// 	g.item = generators[g.itemType]
// 	return nil
// }
//
// func newSequenceParser(
// 	t trace,
// 	typeID uint64,
// 	typeName string,
// 	item generator,
// 	init *node,
// 	excludedTypes *idSet,
// ) *sequenceParser {
// 	return &sequenceParser{
// 		collectionParser: collectionParser{
// 			backtrackingParser: backtrackingParser{
// 				baseParser: baseParser{
// 					trace:         t,
// 					init:          init,
// 					excludedTypes: excludedTypes,
// 					typeID:        typeID,
// 				},
// 				node: &node{typeID: typeID, typ: typeName},
// 			},
// 			firstGenerator: item,
// 		},
// 		generator: item,
// 	}
// }
//
// func (p *sequenceParser) nextParser() (parser, bool, error) {
// 	var (
// 		init     *node
// 		excluded *idSet
// 	)
//
// 	if len(p.node.nodes) > 0 {
// 		excluded = &idSet{}
// 		excluded.add(p.typeID)
// 	} else {
// 		init = p.init
// 		excluded = p.excludedTypes
// 	}
//
// 	if ok, err := p.generator.canCreate(init, excluded); !ok || err != nil {
// 		return nil, ok, err
// 	}
//
// 	parser, err := p.generator.create(p.trace, init, excluded)
// 	return parser, err == nil, err
// }
//
// func (p *sequenceParser) parse(t *token) (parseResult, error) {
// parseLoop:
// 	for {
// 		p.trace.out("parsing", t)
//
// 		p.checkDone(t)
// 		if r, ok := p.checkSkip(); ok {
// 			return r, nil
// 		}
//
// 		if p.parser == nil {
// 			parser, ok, err := p.nextParser()
// 			if !ok || err != nil {
// 				p.trace.out("failed to create next item parser")
// 				p.done = true
// 				return p.abort(err, t)
// 			}
//
// 			p.parser = parser
// 		}
//
// 		if !p.cacheChecked {
// 			p.cacheChecked = true
// 			if r, ok := p.checkCache(t); ok {
// 				p.done = true
// 				return r, nil
// 			}
// 		}
//
// 		r, err := p.parser.parse(t)
// 		if err != nil {
// 			p.trace.out("failed to parse item")
// 			p.done = true
// 			return p.abort(err, t)
// 		}
//
// 		if r.accepting {
// 			if len(p.queue) > 0 {
// 				t, p.queue = p.queue[0], p.queue[1:]
// 				continue parseLoop
// 			}
//
// 			return parseResult{accepting: true}, nil
// 		}
//
// 		p.parser = nil
// 		p.queue = append(r.unparsed, p.queue...)
//
// 		if r.valid && !r.node.zero() {
// 			p.appendParsedItem(r.node, r.fromCache)
// 			if len(p.queue) > 0 {
// 				t, p.queue = p.queue[0], p.queue[1:]
// 				continue parseLoop
// 			}
//
// 			return parseResult{accepting: true}, nil
// 		}
//
// 		if !p.initEvaluated {
// 			p.initEvaluated = true
// 			if ok, err := p.appendInitIfMember(); err != nil {
// 				p.trace.out("failed to check init item membership")
// 				p.done = true
// 				return p.abort(err)
// 			} else if ok {
// 				if len(p.queue) > 0 {
// 					t, p.queue = p.queue[0], p.queue[1:]
// 					continue parseLoop
// 				}
//
// 				return parseResult{accepting: true}, nil
// 			}
// 		}
//
// 		p.trace.out("parse done", p.node, p.node.nodes)
// 		p.done = true
//
// 		if p.node.token != nil {
// 			cache.setMatch(p.node.token, p.typeID, p.node)
// 		}
//
// 		return parseResult{
// 			valid:    true,
// 			unparsed: p.queue,
// 			node:     p.node,
// 		}, nil
// 	}
// }
//
// func group(nodeType string, itemTypes ...string) generator {
// 	itemTypeIDs := make([]uint64, len(itemTypes))
// 	for i, t := range itemTypes {
// 		itemTypeIDs[i] = registerType(t)
// 	}
//
// 	return register(nodeType, &groupGenerator{
// 		nodeType:  nodeType,
// 		itemTypes: itemTypeIDs,
// 		typeID:    registerType(nodeType),
// 	})
// }
//
// func (g *groupGenerator) setTypeID(id uint64) { g.typeID = id }
//
// func (g *groupGenerator) getItemGenerators() ([]generator, error) {
// 	ig := make([]generator, len(g.itemTypes))
// 	for i, it := range g.itemTypes {
// 		g, ok := generators[it]
// 		if !ok {
// 			return nil, unspecifiedParser(it)
// 		}
//
// 		ig[i] = g
// 	}
//
// 	return ig, nil
// }
//
// func (g *groupGenerator) canCreate(init *node, excludedTypes *idSet) (bool, error) {
// 	if len(g.itemTypes) == 0 {
// 		return false, groupWithoutItems(g.nodeType)
// 	}
//
// 	if excludedTypes.contains(g.typeID) {
// 		return false, nil
// 	}
//
// 	first := g.first
//
// 	if ok, err := first.canCreate(init, excludedTypes.extend(g.typeID)); ok || err != nil {
// 		return ok, err
// 	}
//
// 	if ok, err := first.member(init.typeID); ok || err != nil {
// 		return ok, err
// 	}
//
// 	return false, nil
// }
//
// func (g *groupGenerator) create(t trace, init *node, excludedTypes *idSet) (parser, error) {
// 	if len(g.itemGenerators) == 0 {
// 		ig, err := g.getItemGenerators()
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		g.itemGenerators = ig
// 	}
//
// 	return newGroupParser(
// 		t.extend(g.nodeType),
// 		g.typeID,
// 		g.nodeType,
// 		g.itemGenerators,
// 		init,
// 		excludedTypes,
// 	), nil
// }
//
// func (g *groupGenerator) member(typeID uint64) (bool, error) {
// 	return typeID == g.typeID, nil
// }
//
// func (g *groupGenerator) cacheMembers() error {
// 	g.first = generators[g.itemTypes[0]]
// 	return nil
// }
//
// func newGroupParser(
// 	t trace,
// 	typeID uint64,
// 	typeName string,
// 	items []generator,
// 	init *node,
// 	excludedTypes *idSet,
// ) *groupParser {
// 	t.out("create, excluded:", excludedTypes)
// 	return &groupParser{
// 		collectionParser: collectionParser{
// 			backtrackingParser: backtrackingParser{
// 				baseParser: baseParser{
// 					trace:         t,
// 					typeID:        typeID,
// 					init:          init,
// 					excludedTypes: excludedTypes,
// 				},
// 				node: &node{typeID: typeID, typ: typeName},
// 			},
// 			firstGenerator: items[0],
// 		},
// 		items:     items,
// 		firstItem: items[0],
// 	}
// }
//
// func (p *groupParser) nextParser() (parser, bool, error) {
// 	var item generator
// 	item, p.items = p.items[0], p.items[1:]
//
// 	var (
// 		init     *node
// 		excluded *idSet
// 	)
//
// 	if len(p.node.nodes) == 0 {
// 		init = p.init
// 		excluded = p.excludedTypes.extend(p.typeID)
// 	}
//
// 	if ok, err := item.canCreate(init, excluded); !ok || err != nil {
// 		return nil, ok, err
// 	}
//
// 	parser, err := item.create(p.trace, init, excluded)
// 	return parser, err == nil, err
// }
//
// func (p *groupParser) parseOrDone() (parseResult, error) {
// 	p.trace.out("parse done")
// 	p.done = true
// 	p.trace.out("caching group", p.node, p.node.token, p.node.token.offset)
// 	cache.setMatch(p.node.token, p.typeID, p.node)
// 	return parseResult{
// 		valid:    true,
// 		node:     p.node,
// 		unparsed: p.queue,
// 	}, nil
// }
//
// func (p *groupParser) parse(t *token) (parseResult, error) {
// parseLoop:
// 	for {
// 		p.trace.out("parsing", t, p.node, p.node.typ, p.init)
// 		p.checkDone(t)
//
// 		if r, ok := p.checkSkip(); ok {
// 			return r, nil
// 		}
//
// 		if p.parser == nil {
// 			if parser, ok, err := p.nextParser(); err != nil {
// 				p.trace.out("failed to create next item parser")
// 				p.done = true
// 				return p.abort(err, t)
// 			} else if !ok {
// 				panic("this should not happen")
// 			} else {
// 				p.parser = parser
// 			}
// 		}
//
// 		// this prevents checking if the init can be the first item
// 		if !p.cacheChecked {
// 			p.cacheChecked = true
// 			if r, ok := p.checkCache(t); ok {
// 				if !r.valid {
// 					p.trace.out("group from cache, unparsed:", len(r.unparsed))
// 					p.done = true
// 					return r, nil
// 				}
//
// 				if p.init.zero() {
// 					p.trace.out("group from cache, unparsed:", len(r.unparsed))
// 					p.done = true
// 					return r, nil
// 				}
// 			}
// 		}
//
// 		r, err := p.parser.parse(t)
// 		if err != nil {
// 			p.trace.out("failed to parse item")
// 			p.done = true
// 			return p.abort(err, t)
// 		}
//
// 		if r.accepting {
// 			if len(p.queue) > 0 {
// 				t, p.queue = p.queue[0], p.queue[1:]
// 				continue parseLoop
// 			}
//
// 			return parseResult{accepting: true}, nil
// 		}
//
// 		p.parser = nil
// 		p.queue = append(r.unparsed, p.queue...)
// 		p.trace.out("item parse done, queue:", p.queue, r.valid, r.node)
//
// 		if r.valid {
// 			p.trace.out("continue group", r.valid, r.node.zero())
// 			p.initEvaluated = true // only used for the first item
// 			p.appendParsedItem(r.node, r.fromCache)
// 			p.trace.out("checking parse or done", p.queue, p.node)
// 			if len(p.items) > 0 {
// 				p.trace.out("expects more items", len(p.node.nodes), p.node, len(p.items), "queue:", p.queue)
// 				if len(p.queue) > 0 {
// 					t, p.queue = p.queue[0], p.queue[1:]
// 					continue parseLoop
// 				}
//
// 				return parseResult{accepting: true}, nil
// 			}
//
// 			return p.parseOrDone()
// 		}
//
// 		if !p.initEvaluated {
// 			p.initEvaluated = true
// 			if ok, err := p.appendInitIfMember(); err != nil {
// 				p.trace.out("failed to check init item membership")
// 				p.done = true
// 				return p.abort(err)
// 			} else if ok {
// 				p.trace.out("continuing with init")
// 				if len(p.items) > 0 {
// 					p.trace.out("expects more items", len(p.node.nodes), p.node, len(p.items), "queue:", p.queue)
// 					if len(p.queue) > 0 {
// 						t, p.queue = p.queue[0], p.queue[1:]
// 						continue parseLoop
// 					}
//
// 					return parseResult{accepting: true}, nil
// 				}
//
// 				return p.parseOrDone()
// 			}
// 		}
//
// 		if r.valid {
// 			if len(p.items) > 0 {
// 				p.trace.out("expects more items", len(p.node.nodes), p.node, len(p.items), "queue:", p.queue)
// 				if len(p.queue) > 0 {
// 					t, p.queue = p.queue[0], p.queue[1:]
// 					continue parseLoop
// 				}
//
// 				return parseResult{accepting: true}, nil
// 			}
//
// 			return p.parseOrDone()
// 		}
//
// 		p.trace.out("group invalid item")
//
// 		var ct *token
// 		if p.node.zero() {
// 			ct = p.queue[0]
// 		} else {
// 			ct = p.node.token
// 		}
//
// 		if ct != nil {
// 			cache.setNoMatch(ct, p.typeID)
// 		}
//
// 		p.done = true
// 		p.trace.out("returning from end of group", p.node.tokens(), p.queue)
// 		// var pn func(n *node)
// 		// pn = func(n *node) {
// 		// 	p.trace.out(n, n.tokens())
// 		// 	for _, ni := range n.nodes {
// 		// 		pn(ni)
// 		// 	}
// 		// }
// 		// pn(p.node)
//
// 		unparsed := p.node.tokens()
// 		if p.init.length() > len(unparsed) {
// 			unparsed = nil
// 		} else {
// 			unparsed = unparsed[p.init.length():]
// 		}
// 		unparsed = append(unparsed, p.queue...)
//
// 		return parseResult{unparsed: unparsed}, nil
// 	}
// }
//
// func union(nodeType string, elementTypes ...string) generator {
// 	elementTypeIDs := make([]uint64, len(elementTypes))
// 	for i, t := range elementTypes {
// 		elementTypeIDs[i] = registerType(t)
// 	}
//
// 	return register(nodeType, &unionGenerator{
// 		nodeType:     nodeType,
// 		elementTypes: elementTypeIDs,
// 		typeID:       registerType(nodeType),
// 	})
// }
//
// func (g *unionGenerator) setTypeID(id uint64) { g.typeID = id }
//
// func (g *unionGenerator) expand(skip *idSet) ([]generator, error) {
// 	if skip.contains(g.typeID) {
// 		return nil, nil
// 	}
//
// 	var expanded []generator
// 	for _, et := range g.elementTypes {
// 		eg, ok := generators[et]
// 		if !ok {
// 			return nil, unspecifiedParser(et)
// 		}
//
// 		if ug, ok := eg.(*unionGenerator); ok {
// 			ugx, err := ug.expand(skip.extend(g.typeID))
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			expanded = append(expanded, ugx...)
// 		} else if !skip.contains(et) {
// 			expanded = append(expanded, eg)
// 		}
// 	}
//
// 	return expanded, nil
// }
//
// func (g *unionGenerator) canCreate(init *node, excludedTypes *idSet) (bool, error) {
// 	for _, g := range g.expanded {
// 		if ok, err := g.canCreate(init, excludedTypes); ok || err != nil {
// 			return ok, err
// 		}
// 	}
//
// 	return g.member(init.typeID)
// }
//
// func (g *unionGenerator) create(t trace, init *node, excludedTypes *idSet) (parser, error) {
// 	var gen []generator
// 	for _, g := range g.expanded {
// 		if ok, err := g.canCreate(init, excludedTypes); err != nil {
// 			return nil, err
// 		} else if ok {
// 			gen = append(gen, g)
// 		}
// 	}
//
// 	var n *node
// 	if !init.zero() {
// 		if ok, err := g.member(init.typeID); err != nil {
// 			return nil, err
// 		} else if ok {
// 			n = init
// 		}
// 	}
//
// 	return newUnionParser(
// 		t.extend(g.nodeType),
// 		g.typeID,
// 		gen,
// 		n,
// 		init,
// 		excludedTypes,
// 	), nil
// }
//
// func (g *unionGenerator) member(typeID uint64) (bool, error) {
// 	return g.members.contains(typeID), nil
// }
//
// func (g *unionGenerator) cacheMembers() error {
// 	expanded, err := g.expand(nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	if len(expanded) == 0 {
// 		return unionWithoutElements(g.nodeType)
// 	}
//
// 	g.expanded = expanded
//
// 	if g.members != nil {
// 		return nil
// 	}
//
// 	g.members = &idSet{}
// 	for _, gi := range g.expanded {
// 		var id uint64
// 		switch git := gi.(type) {
// 		case *primitiveGenerator:
// 			id = git.typeID
// 		case *optionalGenerator:
// 			id = git.typeID
// 		case *sequenceGenerator:
// 			id = git.typeID
// 		case *groupGenerator:
// 			id = git.typeID
// 		case *unionGenerator:
// 			id = git.typeID
// 		default:
// 			panic(unspecifiedParser(0 /* a lie, TODO */))
// 		}
//
// 		if m, err := gi.member(id); err != nil {
// 			return err
// 		} else if m {
// 			g.members.add(id)
// 		}
// 	}
//
// 	return nil
// }
//
// func newUnionParser(
// 	t trace,
// 	typeID uint64,
// 	elements []generator,
// 	node *node,
// 	init *node,
// 	excludedTypes *idSet,
// ) *unionParser {
// 	return &unionParser{
// 		backtrackingParser: backtrackingParser{
// 			baseParser: baseParser{
// 				trace:         t,
// 				typeID:        typeID,
// 				init:          init,
// 				excludedTypes: excludedTypes,
// 			},
// 			node: node,
// 		},
// 		elements:       elements,
// 		activeElements: elements,
// 		valid:          !node.zero(),
// 	}
// }
//
// func dropSeps(n []*node) []*node {
// 	if isSep == nil {
// 		return n
// 	}
//
// 	nn := make([]*node, 0, len(n))
// 	for _, ni := range n {
// 		if !isSep(ni) {
// 			nn = append(nn, ni)
// 		}
// 	}
//
// 	return nn
// }
//
// func postParseNode(n *node) *node {
// 	n.nodes = postParseNodes(n.nodes)
// 	if pp, ok := postParse[n.typ]; ok {
// 		n = pp(n)
// 	}
//
// 	return n
// }
//
// func postParseNodes(n []*node) []*node {
// 	n = dropSeps(n)
// 	for i, ni := range n {
// 		n[i] = postParseNode(ni)
// 	}
//
// 	return n
// }
//
// func (p *unionParser) cacheKey(t ...*token) *token {
// 	if p.node.zero() {
// 		if len(t) > 0 {
// 			return t[0]
// 		} else if len(p.queue) > 0 {
// 			return p.queue[0]
// 		} else {
// 			panic(invalidParserState(typeNames[p.typeID]))
// 		}
// 	}
//
// 	return p.node.token
// }
//
// func (p *unionParser) setDone(t ...*token) parseResult {
// 	ct := p.cacheKey(t...)
// 	if p.valid {
// 		p.trace.out("union parse success", p.node)
// 		cache.setMatch(ct, p.typeID, p.node)
// 	} else {
// 		p.trace.out("parse failed", t)
// 		cache.setNoMatch(ct, p.typeID)
// 	}
//
// 	r := p.unparsed(t...)
// 	r.valid = p.valid
// 	r.node = p.node
// 	return r
// }
//
// func (p *unionParser) parse(t *token) (parseResult, error) {
// parseLoop:
// 	for {
// 		p.trace.out("parsing", t)
//
// 		p.checkDone(t)
// 		if r, ok := p.checkSkip(); ok {
// 			return r, nil
// 		}
//
// 		// it's a combo
// 		for p.parser == nil {
// 			if len(p.activeElements) == 0 {
// 				p.done = true
// 				p.trace.out("normal done")
// 				return p.setDone(t), nil
// 			}
//
// 			var element generator
// 			element, p.activeElements = p.activeElements[0], p.activeElements[1:]
//
// 			init := p.init
// 			if !p.node.zero() {
// 				init = p.node
// 			}
//
// 			ok, err := element.canCreate(init, p.excludedTypes)
// 			if err != nil {
// 				p.done = true
// 				return p.abort(err, t)
// 			}
//
// 			if !ok {
// 				continue
// 			}
//
// 			parser, err := element.create(p.trace, init, p.excludedTypes)
// 			if err != nil {
// 				p.done = true
// 				return p.abort(err, t)
// 			}
//
// 			p.parser = parser
// 		}
//
// 		// if !p.cacheChecked {
// 		// 	p.cacheChecked = true
// 		// 	if r, ok := p.checkCache(t); ok {
// 		// 		p.done = true
// 		// 		return r, nil
// 		// 	}
// 		// }
//
// 		r, err := p.parser.parse(t)
// 		if err != nil {
// 			p.done = true
// 			return p.abort(err, t)
// 		}
//
// 		if r.accepting {
// 			if len(p.queue) > 0 {
// 				t, p.queue = p.queue[0], p.queue[1:]
// 				continue parseLoop
// 			}
//
// 			return parseResult{accepting: true}, nil
// 		}
//
// 		p.parser = nil
// 		p.queue = append(r.unparsed, p.queue...)
// 		p.trace.out("parser returned", r.unparsed, p.queue, r.fromCache)
//
// 		if !r.valid {
// 			// if len(p.activeElements) == 0 {
// 			// 	p.done = true
// 			// 	return p.setDone(), nil
// 			// }
//
// 			if len(p.queue) > 0 {
// 				t, p.queue = p.queue[0], p.queue[1:]
// 				continue parseLoop
// 			}
//
// 			return parseResult{accepting: true}, nil
// 		}
//
// 		p.trace.out("union successful")
//
// 		if !p.valid || r.node.length() > p.node.length() {
// 			// TODO: the union cache is more complicated
// 			// TODO: the init item cache can be complicated in other case, too
//
// 			if r.fromCache {
// 				if r.node.length() > len(p.queue) {
// 					p.queue, p.skip = nil, r.node.length()-len(p.queue)
// 				} else {
// 					p.queue = p.queue[r.node.length():]
// 				}
// 			}
//
// 			p.trace.out("reset union")
// 			p.node = r.node
// 			p.valid = true
// 			// ct := p.cacheKey()
// 			// cache.setMatch(ct, p.node.typ, p.node)
// 			p.activeElements = p.elements
// 		}
//
// 		// need to do the skip:
// 		// if len(p.activeElements) == 0 {
// 		// 	p.done = true
// 		// 	p.trace.out("no more elements to try", len(p.queue))
// 		// 	return p.setDone(), nil
// 		// }
//
// 		if len(p.queue) > 0 {
// 			t, p.queue = p.queue[0], p.queue[1:]
// 			continue parseLoop
// 		}
//
// 		return parseResult{accepting: true}, nil
// 	}
// }
//
// func setPostParse(p map[string]func(*node) *node) {
// 	for pi, pp := range p {
// 		postParse[pi] = pp
// 	}
// }
//
// func parse(l traceLevel, g generator, r *tokenReader) (*node, error) {
// 	if ok, err := g.canCreate(zeroNode, nil); err != nil {
// 		return zeroNode, err
// 	} else if !ok {
// 		return zeroNode, errInvalidRootGenerator
// 	}
//
// 	trace := newTrace(l)
// 	p, err := g.create(trace, zeroNode, nil)
// 	if err != nil {
// 		return zeroNode, err
// 	}
//
// 	last := parseResult{accepting: true}
// 	for {
// 		t, err := r.next()
// 		if err != nil && err != io.EOF {
// 			return zeroNode, err
// 		}
//
// 		if !last.accepting {
// 			if err != io.EOF {
// 				return zeroNode, unexpectedToken("root", &t)
// 			}
//
// 			return last.node, nil
// 		}
//
// 		if err == io.EOF {
// 			last, err = p.parse(eofToken)
// 			if err != nil {
// 				return zeroNode, err
// 			}
//
// 			if !last.valid {
// 				trace.out("last not valid")
// 				return zeroNode, errUnexpectedEOF
// 			}
//
// 			if len(last.unparsed) != 1 || last.unparsed[0] != eofToken {
// 				trace.out("unexpected unparsed", len(last.unparsed), last.unparsed)
// 				return zeroNode, errUnexpectedEOF
// 			}
//
// 			trace.out("parsed", last.node)
// 			return postParseNode(last.node), nil
// 		}
//
// 		last, err = p.parse(&t)
// 		if err != nil {
// 			return zeroNode, err
// 		}
//
// 		if !last.accepting {
// 			if !last.valid {
// 				return zeroNode, unexpectedToken("root", &t)
// 			}
//
// 			if len(last.unparsed) > 0 {
// 				return zeroNode, unexpectedToken("root", last.unparsed[0])
// 			}
// 		}
// 	}
// }
