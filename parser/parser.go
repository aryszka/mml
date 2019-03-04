/*
This file was generated with treerack (https://github.com/aryszka/treerack).

The contents of this file fall under different licenses.

The code between the "// head" and "// eo head" lines falls under the same
license as the source code of treerack (https://github.com/aryszka/treerack),
unless explicitly stated otherwise, if treerack's license allows changing the
license of this source code.

Treerack's license: MIT https://opensource.org/licenses/MIT
where YEAR=2017, COPYRIGHT HOLDER=Arpad Ryszka (arpad.ryszka@gmail.com)

The rest of the content of this file falls under the same license as the one
that the user of treerack generating this file declares for it, or it is
unlicensed.
*/

package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type charParser struct {
	name   string
	id     int
	not    bool
	chars  []rune
	ranges [][]rune
}
type charBuilder struct {
	name string
	id   int
}

func (p *charParser) nodeName() string {
	return p.name
}
func (p *charParser) nodeID() int {
	return p.id
}
func (p *charParser) commitType() CommitType {
	return Alias
}
func matchChar(chars []rune, ranges [][]rune, not bool, char rune) bool {
	for _, ci := range chars {
		if ci == char {
			return !not
		}
	}
	for _, ri := range ranges {
		if char >= ri[0] && char <= ri[1] {
			return !not
		}
	}
	return not
}
func (p *charParser) match(t rune) bool {
	return matchChar(p.chars, p.ranges, p.not, t)
}
func (p *charParser) parse(c *context) {
	if tok, ok := c.token(); !ok || !p.match(tok) {
		if c.offset > c.failOffset {
			c.failOffset = c.offset
			c.failingParser = nil
		}
		c.fail(c.offset)
		return
	}
	c.success(c.offset + 1)
}
func (b *charBuilder) nodeName() string {
	return b.name
}
func (b *charBuilder) nodeID() int {
	return b.id
}
func (b *charBuilder) build(c *context) ([]*Node, bool) {
	return nil, false
}

type sequenceParser struct {
	name            string
	id              int
	commit          CommitType
	items           []parser
	ranges          [][]int
	generalizations []int
	allChars        bool
}
type sequenceBuilder struct {
	name            string
	id              int
	commit          CommitType
	items           []builder
	ranges          [][]int
	generalizations []int
	allChars        bool
}

func (p *sequenceParser) nodeName() string {
	return p.name
}
func (p *sequenceParser) nodeID() int {
	return p.id
}
func (p *sequenceParser) commitType() CommitType {
	return p.commit
}
func (p *sequenceParser) parse(c *context) {
	if !p.allChars {
		if c.results.pending(c.offset, p.id) {
			c.fail(c.offset)
			return
		}
		c.results.markPending(c.offset, p.id)
	}
	var (
		currentCount int
		parsed       bool
	)
	itemIndex := 0
	from := c.offset
	to := c.offset
	for itemIndex < len(p.items) {
		p.items[itemIndex].parse(c)
		if !c.matchLast {
			if currentCount >= p.ranges[itemIndex][0] {
				itemIndex++
				currentCount = 0
				continue
			}
			c.offset = from
			if c.fromResults(p) {
				if to > c.failOffset {
					c.failOffset = -1
					c.failingParser = nil
				}
				if !p.allChars {
					c.results.unmarkPending(from, p.id)
				}
				return
			}
			if c.failingParser == nil && p.commit&userDefined != 0 && p.commit&Whitespace == 0 && p.commit&FailPass == 0 {
				c.failingParser = p
			}
			c.fail(from)
			if !p.allChars {
				c.results.unmarkPending(from, p.id)
			}
			return
		}
		parsed = c.offset > to
		if parsed {
			currentCount++
		}
		to = c.offset
		if !parsed || p.ranges[itemIndex][1] > 0 && currentCount == p.ranges[itemIndex][1] {
			itemIndex++
			currentCount = 0
		}
	}
	if p.commit&NoKeyword != 0 && c.isKeyword(from, to) {
		if c.failingParser == nil && p.commit&userDefined != 0 && p.commit&Whitespace == 0 && p.commit&FailPass == 0 {
			c.failingParser = p
		}
		c.fail(from)
		if !p.allChars {
			c.results.unmarkPending(from, p.id)
		}
		return
	}
	for _, g := range p.generalizations {
		if c.results.pending(from, g) {
			c.results.setMatch(from, g, to)
		}
	}
	if to > c.failOffset {
		c.failOffset = -1
		c.failingParser = nil
	}
	c.results.setMatch(from, p.id, to)
	c.success(to)
	if !p.allChars {
		c.results.unmarkPending(from, p.id)
	}
}
func (b *sequenceBuilder) nodeName() string {
	return b.name
}
func (b *sequenceBuilder) nodeID() int {
	return b.id
}
func (b *sequenceBuilder) build(c *context) ([]*Node, bool) {
	to, ok := c.results.longestMatch(c.offset, b.id)
	if !ok {
		return nil, false
	}
	from := c.offset
	parsed := to > from
	if b.allChars {
		c.offset = to
		if b.commit&Alias != 0 {
			return nil, true
		}
		return []*Node{{Name: b.name, From: from, To: to, tokens: c.tokens}}, true
	} else if parsed {
		c.results.dropMatchTo(c.offset, b.id, to)
		for _, g := range b.generalizations {
			c.results.dropMatchTo(c.offset, g, to)
		}
	} else {
		if c.results.pending(c.offset, b.id) {
			return nil, false
		}
		c.results.markPending(c.offset, b.id)
		for _, g := range b.generalizations {
			c.results.markPending(c.offset, g)
		}
	}
	var (
		itemIndex    int
		currentCount int
		nodes        []*Node
	)
	for itemIndex < len(b.items) {
		itemFrom := c.offset
		n, ok := b.items[itemIndex].build(c)
		if !ok {
			itemIndex++
			currentCount = 0
			continue
		}
		if c.offset > itemFrom {
			nodes = append(nodes, n...)
			currentCount++
			if b.ranges[itemIndex][1] > 0 && currentCount == b.ranges[itemIndex][1] {
				itemIndex++
				currentCount = 0
			}
			continue
		}
		if currentCount < b.ranges[itemIndex][0] {
			for i := 0; i < b.ranges[itemIndex][0]-currentCount; i++ {
				nodes = append(nodes, n...)
			}
		}
		itemIndex++
		currentCount = 0
	}
	if !parsed {
		c.results.unmarkPending(from, b.id)
		for _, g := range b.generalizations {
			c.results.unmarkPending(from, g)
		}
	}
	if b.commit&Alias != 0 {
		return nodes, true
	}
	return []*Node{{Name: b.name, From: from, To: to, Nodes: nodes, tokens: c.tokens}}, true
}

type choiceParser struct {
	name            string
	id              int
	commit          CommitType
	options         []parser
	generalizations []int
}
type choiceBuilder struct {
	name            string
	id              int
	commit          CommitType
	options         []builder
	generalizations []int
}

func (p *choiceParser) nodeName() string {
	return p.name
}
func (p *choiceParser) nodeID() int {
	return p.id
}
func (p *choiceParser) commitType() CommitType {
	return p.commit
}
func (p *choiceParser) parse(c *context) {
	if c.fromResults(p) {
		return
	}
	if c.results.pending(c.offset, p.id) {
		c.fail(c.offset)
		return
	}
	c.results.markPending(c.offset, p.id)
	var (
		match         bool
		optionIndex   int
		foundMatch    bool
		failingParser parser
	)
	from := c.offset
	to := c.offset
	initialFailOffset := c.failOffset
	initialFailingParser := c.failingParser
	failOffset := initialFailOffset
	for {
		foundMatch = false
		optionIndex = 0
		for optionIndex < len(p.options) {
			p.options[optionIndex].parse(c)
			optionIndex++
			if !c.matchLast {
				if c.failOffset > failOffset {
					failOffset = c.failOffset
					failingParser = c.failingParser
				}
			}
			if !c.matchLast || match && c.offset <= to {
				c.offset = from
				continue
			}
			match = true
			foundMatch = true
			to = c.offset
			c.offset = from
			c.results.setMatch(from, p.id, to)
		}
		if !foundMatch {
			break
		}
	}
	if match {
		if p.commit&NoKeyword != 0 && c.isKeyword(from, to) {
			if c.failingParser == nil && p.commit&userDefined != 0 && p.commit&Whitespace == 0 && p.commit&FailPass == 0 {
				c.failingParser = p
			}
			c.fail(from)
			c.results.unmarkPending(from, p.id)
			return
		}
		if failOffset > to {
			c.failOffset = failOffset
			c.failingParser = failingParser
		} else if to > initialFailOffset {
			c.failOffset = -1
			c.failingParser = nil
		} else {
			c.failOffset = initialFailOffset
			c.failingParser = initialFailingParser
		}
		c.success(to)
		c.results.unmarkPending(from, p.id)
		return
	}
	if failOffset > initialFailOffset {
		c.failOffset = failOffset
		c.failingParser = failingParser
		if c.failingParser == nil && p.commitType()&userDefined != 0 && p.commitType()&Whitespace == 0 && p.commitType()&FailPass == 0 {
			c.failingParser = p
		}
	}
	c.results.setNoMatch(from, p.id)
	c.fail(from)
	c.results.unmarkPending(from, p.id)
}
func (b *choiceBuilder) nodeName() string {
	return b.name
}
func (b *choiceBuilder) nodeID() int {
	return b.id
}
func (b *choiceBuilder) build(c *context) ([]*Node, bool) {
	to, ok := c.results.longestMatch(c.offset, b.id)
	if !ok {
		return nil, false
	}
	from := c.offset
	parsed := to > from
	if parsed {
		c.results.dropMatchTo(c.offset, b.id, to)
		for _, g := range b.generalizations {
			c.results.dropMatchTo(c.offset, g, to)
		}
	} else {
		if c.results.pending(c.offset, b.id) {
			return nil, false
		}
		c.results.markPending(c.offset, b.id)
		for _, g := range b.generalizations {
			c.results.markPending(c.offset, g)
		}
	}
	var option builder
	for _, o := range b.options {
		if c.results.hasMatchTo(c.offset, o.nodeID(), to) {
			option = o
			break
		}
	}
	n, _ := option.build(c)
	if !parsed {
		c.results.unmarkPending(from, b.id)
		for _, g := range b.generalizations {
			c.results.unmarkPending(from, g)
		}
	}
	if b.commit&Alias != 0 {
		return n, true
	}
	return []*Node{{Name: b.name, From: from, To: to, Nodes: n, tokens: c.tokens}}, true
}

type idSet struct{ ids []uint }

func divModBits(id int) (int, int) {
	return id / strconv.IntSize, id % strconv.IntSize
}
func (s *idSet) set(id int) {
	d, m := divModBits(id)
	if d >= len(s.ids) {
		if d < cap(s.ids) {
			s.ids = s.ids[:d+1]
		} else {
			s.ids = s.ids[:cap(s.ids)]
			for i := cap(s.ids); i <= d; i++ {
				s.ids = append(s.ids, 0)
			}
		}
	}
	s.ids[d] |= 1 << uint(m)
}
func (s *idSet) unset(id int) {
	d, m := divModBits(id)
	if d >= len(s.ids) {
		return
	}
	s.ids[d] &^= 1 << uint(m)
}
func (s *idSet) has(id int) bool {
	d, m := divModBits(id)
	if d >= len(s.ids) {
		return false
	}
	return s.ids[d]&(1<<uint(m)) != 0
}

type results struct {
	noMatch   []*idSet
	match     [][]int
	isPending [][]int
}

func ensureOffsetInts(ints [][]int, offset int) [][]int {
	if len(ints) > offset {
		return ints
	}
	if cap(ints) > offset {
		ints = ints[:offset+1]
		return ints
	}
	ints = ints[:cap(ints)]
	for i := len(ints); i <= offset; i++ {
		ints = append(ints, nil)
	}
	return ints
}
func ensureOffsetIDs(ids []*idSet, offset int) []*idSet {
	if len(ids) > offset {
		return ids
	}
	if cap(ids) > offset {
		ids = ids[:offset+1]
		return ids
	}
	ids = ids[:cap(ids)]
	for i := len(ids); i <= offset; i++ {
		ids = append(ids, nil)
	}
	return ids
}
func (r *results) setMatch(offset, id, to int) {
	r.match = ensureOffsetInts(r.match, offset)
	for i := 0; i < len(r.match[offset]); i += 2 {
		if r.match[offset][i] != id || r.match[offset][i+1] != to {
			continue
		}
		return
	}
	r.match[offset] = append(r.match[offset], id, to)
}
func (r *results) setNoMatch(offset, id int) {
	if len(r.match) > offset {
		for i := 0; i < len(r.match[offset]); i += 2 {
			if r.match[offset][i] != id {
				continue
			}
			return
		}
	}
	r.noMatch = ensureOffsetIDs(r.noMatch, offset)
	if r.noMatch[offset] == nil {
		r.noMatch[offset] = &idSet{}
	}
	r.noMatch[offset].set(id)
}
func (r *results) hasMatchTo(offset, id, to int) bool {
	if len(r.match) <= offset {
		return false
	}
	for i := 0; i < len(r.match[offset]); i += 2 {
		if r.match[offset][i] != id {
			continue
		}
		if r.match[offset][i+1] == to {
			return true
		}
	}
	return false
}
func (r *results) longestMatch(offset, id int) (int, bool) {
	if len(r.match) <= offset {
		return 0, false
	}
	var found bool
	to := -1
	for i := 0; i < len(r.match[offset]); i += 2 {
		if r.match[offset][i] != id {
			continue
		}
		if r.match[offset][i+1] > to {
			to = r.match[offset][i+1]
		}
		found = true
	}
	return to, found
}
func (r *results) longestResult(offset, id int) (int, bool, bool) {
	if len(r.noMatch) > offset && r.noMatch[offset] != nil && r.noMatch[offset].has(id) {
		return 0, false, true
	}
	to, ok := r.longestMatch(offset, id)
	return to, ok, ok
}
func (r *results) dropMatchTo(offset, id, to int) {
	for i := 0; i < len(r.match[offset]); i += 2 {
		if r.match[offset][i] != id {
			continue
		}
		if r.match[offset][i+1] == to {
			r.match[offset][i] = -1
			return
		}
	}
}
func (r *results) resetPending() {
	r.isPending = nil
}
func (r *results) pending(offset, id int) bool {
	if len(r.isPending) <= id {
		return false
	}
	for i := range r.isPending[id] {
		if r.isPending[id][i] == offset {
			return true
		}
	}
	return false
}
func (r *results) markPending(offset, id int) {
	r.isPending = ensureOffsetInts(r.isPending, id)
	for i := range r.isPending[id] {
		if r.isPending[id][i] == -1 {
			r.isPending[id][i] = offset
			return
		}
	}
	r.isPending[id] = append(r.isPending[id], offset)
}
func (r *results) unmarkPending(offset, id int) {
	for i := range r.isPending[id] {
		if r.isPending[id][i] == offset {
			r.isPending[id][i] = -1
			break
		}
	}
}

type context struct {
	reader        io.RuneReader
	keywords      []parser
	offset        int
	readOffset    int
	consumed      int
	offsetLimit   int
	failOffset    int
	failingParser parser
	readErr       error
	eof           bool
	results       *results
	tokens        []rune
	matchLast     bool
}

func newContext(r io.RuneReader, keywords []parser) *context {
	return &context{reader: r, keywords: keywords, results: &results{}, offsetLimit: -1, failOffset: -1}
}
func (c *context) read() bool {
	if c.eof || c.readErr != nil {
		return false
	}
	token, n, err := c.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			if n == 0 {
				c.eof = true
				return false
			}
		} else {
			c.readErr = err
			return false
		}
	}
	c.readOffset++
	if token == unicode.ReplacementChar {
		c.readErr = ErrInvalidUnicodeCharacter
		return false
	}
	c.tokens = append(c.tokens, token)
	return true
}
func (c *context) token() (rune, bool) {
	if c.offset == c.offsetLimit {
		return 0, false
	}
	if c.offset == c.readOffset {
		if !c.read() {
			return 0, false
		}
	}
	return c.tokens[c.offset], true
}
func (c *context) fromResults(p parser) bool {
	to, m, ok := c.results.longestResult(c.offset, p.nodeID())
	if !ok {
		return false
	}
	if m {
		c.success(to)
	} else {
		c.fail(c.offset)
	}
	return true
}
func (c *context) isKeyword(from, to int) bool {
	ol := c.offsetLimit
	c.offsetLimit = to
	defer func() {
		c.offsetLimit = ol
	}()
	for _, kw := range c.keywords {
		c.offset = from
		kw.parse(c)
		if c.matchLast && c.offset == to {
			return true
		}
	}
	return false
}
func (c *context) success(to int) {
	c.offset = to
	c.matchLast = true
	if to > c.consumed {
		c.consumed = to
	}
}
func (c *context) fail(offset int) {
	c.offset = offset
	c.matchLast = false
}
func findLine(tokens []rune, offset int) (line, column int) {
	tokens = tokens[:offset]
	for i := range tokens {
		column++
		if tokens[i] == '\n' {
			column = 0
			line++
		}
	}
	return
}
func (c *context) parseError(p parser) error {
	definition := p.nodeName()
	flagIndex := strings.Index(definition, ":")
	if flagIndex > 0 {
		definition = definition[:flagIndex]
	}
	if c.failingParser == nil {
		c.failOffset = c.consumed
	}
	line, col := findLine(c.tokens, c.failOffset)
	return &ParseError{Offset: c.failOffset, Line: line, Column: col, Definition: definition}
}
func (c *context) finalizeParse(root parser) error {
	fp := c.failingParser
	if fp == nil {
		fp = root
	}
	to, match, found := c.results.longestResult(0, root.nodeID())
	if !found || !match || found && match && to < c.readOffset {
		return c.parseError(fp)
	}
	c.read()
	if c.eof {
		return nil
	}
	if c.readErr != nil {
		return c.readErr
	}
	return c.parseError(root)
}

type Node struct {
	Name     string
	Nodes    []*Node
	From, To int
	tokens   []rune
}

func (n *Node) Tokens() []rune {
	return n.tokens
}
func (n *Node) String() string {
	return fmt.Sprintf("%s:%d:%d:%s", n.Name, n.From, n.To, n.Text())
}
func (n *Node) Text() string {
	return string(n.Tokens()[n.From:n.To])
}

type CommitType int

const (
	None  CommitType = 0
	Alias CommitType = 1 << iota
	Whitespace
	NoWhitespace
	Keyword
	NoKeyword
	FailPass
	Root
	userDefined
)

type formatFlags int

const (
	formatNone   formatFlags = 0
	formatPretty formatFlags = 1 << iota
	formatIncludeComments
)

type ParseError struct {
	Input      string
	Offset     int
	Line       int
	Column     int
	Definition string
}
type parser interface {
	nodeName() string
	nodeID() int
	commitType() CommitType
	parse(*context)
}
type builder interface {
	nodeName() string
	nodeID() int
	build(*context) ([]*Node, bool)
}

var ErrInvalidUnicodeCharacter = errors.New("invalid unicode character")

func (pe *ParseError) Error() string {
	return fmt.Sprintf("%s:%d:%d:parse failed, parsing: %s", pe.Input, pe.Line+1, pe.Column+1, pe.Definition)
}
func parseInput(r io.Reader, p parser, b builder, kw []parser) (*Node, error) {
	c := newContext(bufio.NewReader(r), kw)
	p.parse(c)
	if c.readErr != nil {
		return nil, c.readErr
	}
	if err := c.finalizeParse(p); err != nil {
		if perr, ok := err.(*ParseError); ok {
			perr.Input = "<input>"
		}
		return nil, err
	}
	c.offset = 0
	c.results.resetPending()
	n, _ := b.build(c)
	return n[0], nil
}

func Parse(r io.Reader) (*Node, error) {
	var p48 = sequenceParser{id: 48, commit: 280, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{209, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p44 = charParser{id: 44, chars: []rune{116}}
	var p45 = charParser{id: 45, chars: []rune{114}}
	var p46 = charParser{id: 46, chars: []rune{117}}
	var p47 = charParser{id: 47, chars: []rune{101}}
	p48.items = []parser{&p44, &p45, &p46, &p47}
	var p54 = sequenceParser{id: 54, commit: 280, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{209, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p49 = charParser{id: 49, chars: []rune{102}}
	var p50 = charParser{id: 50, chars: []rune{97}}
	var p51 = charParser{id: 51, chars: []rune{108}}
	var p52 = charParser{id: 52, chars: []rune{115}}
	var p53 = charParser{id: 53, chars: []rune{101}}
	p54.items = []parser{&p49, &p50, &p51, &p52, &p53}
	var p61 = sequenceParser{id: 61, commit: 282, name: "return", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{289, 798}}
	var p55 = charParser{id: 55, chars: []rune{114}}
	var p56 = charParser{id: 56, chars: []rune{101}}
	var p57 = charParser{id: 57, chars: []rune{116}}
	var p58 = charParser{id: 58, chars: []rune{117}}
	var p59 = charParser{id: 59, chars: []rune{114}}
	var p60 = charParser{id: 60, chars: []rune{110}}
	p61.items = []parser{&p55, &p56, &p57, &p58, &p59, &p60}
	var p67 = sequenceParser{id: 67, commit: 282, name: "check", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p62 = charParser{id: 62, chars: []rune{99}}
	var p63 = charParser{id: 63, chars: []rune{104}}
	var p64 = charParser{id: 64, chars: []rune{101}}
	var p65 = charParser{id: 65, chars: []rune{99}}
	var p66 = charParser{id: 66, chars: []rune{107}}
	p67.items = []parser{&p62, &p63, &p64, &p65, &p66}
	var p70 = sequenceParser{id: 70, commit: 282, name: "fn", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p68 = charParser{id: 68, chars: []rune{102}}
	var p69 = charParser{id: 69, chars: []rune{110}}
	p70.items = []parser{&p68, &p69}
	var p73 = sequenceParser{id: 73, commit: 282, name: "if", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p71 = charParser{id: 71, chars: []rune{105}}
	var p72 = charParser{id: 72, chars: []rune{102}}
	p73.items = []parser{&p71, &p72}
	var p78 = sequenceParser{id: 78, commit: 282, name: "else", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p74 = charParser{id: 74, chars: []rune{101}}
	var p75 = charParser{id: 75, chars: []rune{108}}
	var p76 = charParser{id: 76, chars: []rune{115}}
	var p77 = charParser{id: 77, chars: []rune{101}}
	p78.items = []parser{&p74, &p75, &p76, &p77}
	var p85 = sequenceParser{id: 85, commit: 282, name: "switch", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p79 = charParser{id: 79, chars: []rune{115}}
	var p80 = charParser{id: 80, chars: []rune{119}}
	var p81 = charParser{id: 81, chars: []rune{105}}
	var p82 = charParser{id: 82, chars: []rune{116}}
	var p83 = charParser{id: 83, chars: []rune{99}}
	var p84 = charParser{id: 84, chars: []rune{104}}
	p85.items = []parser{&p79, &p80, &p81, &p82, &p83, &p84}
	var p90 = sequenceParser{id: 90, commit: 282, name: "case", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p86 = charParser{id: 86, chars: []rune{99}}
	var p87 = charParser{id: 87, chars: []rune{97}}
	var p88 = charParser{id: 88, chars: []rune{115}}
	var p89 = charParser{id: 89, chars: []rune{101}}
	p90.items = []parser{&p86, &p87, &p88, &p89}
	var p98 = sequenceParser{id: 98, commit: 282, name: "default", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p91 = charParser{id: 91, chars: []rune{100}}
	var p92 = charParser{id: 92, chars: []rune{101}}
	var p93 = charParser{id: 93, chars: []rune{102}}
	var p94 = charParser{id: 94, chars: []rune{97}}
	var p95 = charParser{id: 95, chars: []rune{117}}
	var p96 = charParser{id: 96, chars: []rune{108}}
	var p97 = charParser{id: 97, chars: []rune{116}}
	p98.items = []parser{&p91, &p92, &p93, &p94, &p95, &p96, &p97}
	var p103 = sequenceParser{id: 103, commit: 282, name: "send", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p99 = charParser{id: 99, chars: []rune{115}}
	var p100 = charParser{id: 100, chars: []rune{101}}
	var p101 = charParser{id: 101, chars: []rune{110}}
	var p102 = charParser{id: 102, chars: []rune{100}}
	p103.items = []parser{&p99, &p100, &p101, &p102}
	var p111 = sequenceParser{id: 111, commit: 282, name: "receive", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p104 = charParser{id: 104, chars: []rune{114}}
	var p105 = charParser{id: 105, chars: []rune{101}}
	var p106 = charParser{id: 106, chars: []rune{99}}
	var p107 = charParser{id: 107, chars: []rune{101}}
	var p108 = charParser{id: 108, chars: []rune{105}}
	var p109 = charParser{id: 109, chars: []rune{118}}
	var p110 = charParser{id: 110, chars: []rune{101}}
	p111.items = []parser{&p104, &p105, &p106, &p107, &p108, &p109, &p110}
	var p118 = sequenceParser{id: 118, commit: 282, name: "select", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p112 = charParser{id: 112, chars: []rune{115}}
	var p113 = charParser{id: 113, chars: []rune{101}}
	var p114 = charParser{id: 114, chars: []rune{108}}
	var p115 = charParser{id: 115, chars: []rune{101}}
	var p116 = charParser{id: 116, chars: []rune{99}}
	var p117 = charParser{id: 117, chars: []rune{116}}
	p118.items = []parser{&p112, &p113, &p114, &p115, &p116, &p117}
	var p121 = sequenceParser{id: 121, commit: 282, name: "go", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p119 = charParser{id: 119, chars: []rune{103}}
	var p120 = charParser{id: 120, chars: []rune{111}}
	p121.items = []parser{&p119, &p120}
	var p127 = sequenceParser{id: 127, commit: 282, name: "defer", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p122 = charParser{id: 122, chars: []rune{100}}
	var p123 = charParser{id: 123, chars: []rune{101}}
	var p124 = charParser{id: 124, chars: []rune{102}}
	var p125 = charParser{id: 125, chars: []rune{101}}
	var p126 = charParser{id: 126, chars: []rune{114}}
	p127.items = []parser{&p122, &p123, &p124, &p125, &p126}
	var p130 = sequenceParser{id: 130, commit: 282, name: "in", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p128 = charParser{id: 128, chars: []rune{105}}
	var p129 = charParser{id: 129, chars: []rune{110}}
	p130.items = []parser{&p128, &p129}
	var p134 = sequenceParser{id: 134, commit: 282, name: "for", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p131 = charParser{id: 131, chars: []rune{102}}
	var p132 = charParser{id: 132, chars: []rune{111}}
	var p133 = charParser{id: 133, chars: []rune{114}}
	p134.items = []parser{&p131, &p132, &p133}
	var p140 = sequenceParser{id: 140, commit: 280, name: "break", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{610, 798}}
	var p135 = charParser{id: 135, chars: []rune{98}}
	var p136 = charParser{id: 136, chars: []rune{114}}
	var p137 = charParser{id: 137, chars: []rune{101}}
	var p138 = charParser{id: 138, chars: []rune{97}}
	var p139 = charParser{id: 139, chars: []rune{107}}
	p140.items = []parser{&p135, &p136, &p137, &p138, &p139}
	var p149 = sequenceParser{id: 149, commit: 280, name: "continue", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{610, 798}}
	var p141 = charParser{id: 141, chars: []rune{99}}
	var p142 = charParser{id: 142, chars: []rune{111}}
	var p143 = charParser{id: 143, chars: []rune{110}}
	var p144 = charParser{id: 144, chars: []rune{116}}
	var p145 = charParser{id: 145, chars: []rune{105}}
	var p146 = charParser{id: 146, chars: []rune{110}}
	var p147 = charParser{id: 147, chars: []rune{117}}
	var p148 = charParser{id: 148, chars: []rune{101}}
	p149.items = []parser{&p141, &p142, &p143, &p144, &p145, &p146, &p147, &p148}
	var p153 = sequenceParser{id: 153, commit: 282, name: "let", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p150 = charParser{id: 150, chars: []rune{108}}
	var p151 = charParser{id: 151, chars: []rune{101}}
	var p152 = charParser{id: 152, chars: []rune{116}}
	p153.items = []parser{&p150, &p151, &p152}
	var p160 = sequenceParser{id: 160, commit: 282, name: "export", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p154 = charParser{id: 154, chars: []rune{101}}
	var p155 = charParser{id: 155, chars: []rune{120}}
	var p156 = charParser{id: 156, chars: []rune{112}}
	var p157 = charParser{id: 157, chars: []rune{111}}
	var p158 = charParser{id: 158, chars: []rune{114}}
	var p159 = charParser{id: 159, chars: []rune{116}}
	p160.items = []parser{&p154, &p155, &p156, &p157, &p158, &p159}
	var p164 = sequenceParser{id: 164, commit: 282, name: "use", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p161 = charParser{id: 161, chars: []rune{117}}
	var p162 = charParser{id: 162, chars: []rune{115}}
	var p163 = charParser{id: 163, chars: []rune{101}}
	p164.items = []parser{&p161, &p162, &p163}
	var p830 = sequenceParser{id: 830, commit: 128, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p828 = choiceParser{id: 828, commit: 2}
	var p826 = choiceParser{id: 826, commit: 262, name: "ws", generalizations: []int{828}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p826.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p827 = sequenceParser{id: 827, commit: 262, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{828}}
	var p43 = choiceParser{id: 43, commit: 258, name: "comment"}
	var p26 = sequenceParser{id: 26, commit: 256, name: "line-comment", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{43}}
	var p22 = sequenceParser{id: 22, commit: 266, name: "comment-line", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p21 = sequenceParser{id: 21, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p19 = charParser{id: 19, chars: []rune{47}}
	var p20 = charParser{id: 20, chars: []rune{47}}
	p21.items = []parser{&p19, &p20}
	var p18 = sequenceParser{id: 18, commit: 264, name: "line-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var p17 = sequenceParser{id: 17, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p16 = charParser{id: 16, not: true, chars: []rune{10}}
	p17.items = []parser{&p16}
	p18.items = []parser{&p17}
	p22.items = []parser{&p21, &p18}
	var p25 = sequenceParser{id: 25, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p23 = sequenceParser{id: 23, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p14 = sequenceParser{id: 14, commit: 266, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{810, 222}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p23.items = []parser{&p14, &p828, &p22}
	var p24 = sequenceParser{id: 24, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p24.items = []parser{&p828, &p23}
	p25.items = []parser{&p828, &p23, &p24}
	p26.items = []parser{&p22, &p25}
	var p42 = sequenceParser{id: 42, commit: 264, name: "block-comment", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{43}}
	var p38 = sequenceParser{id: 38, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p36 = charParser{id: 36, chars: []rune{47}}
	var p37 = charParser{id: 37, chars: []rune{42}}
	p38.items = []parser{&p36, &p37}
	var p35 = sequenceParser{id: 35, commit: 264, name: "block-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var p34 = choiceParser{id: 34, commit: 10}
	var p28 = sequenceParser{id: 28, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{34}}
	var p27 = charParser{id: 27, not: true, chars: []rune{42}}
	p28.items = []parser{&p27}
	var p33 = sequenceParser{id: 33, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{34}}
	var p30 = sequenceParser{id: 30, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p29 = charParser{id: 29, chars: []rune{42}}
	p30.items = []parser{&p29}
	var p32 = sequenceParser{id: 32, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p31 = charParser{id: 31, not: true, chars: []rune{47}}
	p32.items = []parser{&p31}
	p33.items = []parser{&p30, &p32}
	p34.options = []parser{&p28, &p33}
	p35.items = []parser{&p34}
	var p41 = sequenceParser{id: 41, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p39 = charParser{id: 39, chars: []rune{42}}
	var p40 = charParser{id: 40, chars: []rune{47}}
	p41.items = []parser{&p39, &p40}
	p42.items = []parser{&p38, &p35, &p41}
	p43.options = []parser{&p26, &p42}
	p827.items = []parser{&p43}
	p828.options = []parser{&p826, &p827}
	var p829 = sequenceParser{id: 829, commit: 258, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p825 = sequenceParser{id: 825, commit: 256, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p824 = sequenceParser{id: 824, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p822 = charParser{id: 822, chars: []rune{35}}
	var p823 = charParser{id: 823, chars: []rune{33}}
	p824.items = []parser{&p822, &p823}
	var p821 = sequenceParser{id: 821, commit: 256, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p820 = sequenceParser{id: 820, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p818 = sequenceParser{id: 818, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p817 = charParser{id: 817, not: true, chars: []rune{10}}
	p818.items = []parser{&p817}
	var p819 = sequenceParser{id: 819, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p819.items = []parser{&p828, &p818}
	p820.items = []parser{&p818, &p819}
	p821.items = []parser{&p820}
	p825.items = []parser{&p824, &p828, &p821, &p828, &p14}
	var p812 = sequenceParser{id: 812, commit: 258, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p810 = choiceParser{id: 810, commit: 2}
	var p809 = sequenceParser{id: 809, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{810}}
	var p808 = charParser{id: 808, chars: []rune{59}}
	p809.items = []parser{&p808}
	p810.options = []parser{&p809, &p14}
	var p811 = sequenceParser{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p811.items = []parser{&p828, &p810}
	p812.items = []parser{&p810, &p811}
	var p816 = sequenceParser{id: 816, commit: 258, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p798 = choiceParser{id: 798, commit: 258, name: "statement"}
	var p289 = choiceParser{id: 289, commit: 256, name: "ret", generalizations: []int{798}}
	var p288 = sequenceParser{id: 288, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{289, 798}}
	var p287 = sequenceParser{id: 287, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p286 = sequenceParser{id: 286, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p286.items = []parser{&p828, &p14}
	p287.items = []parser{&p828, &p14, &p286}
	var p502 = choiceParser{id: 502, commit: 258, name: "expression", generalizations: []int{225, 314, 618, 611}}
	var p373 = choiceParser{id: 373, commit: 258, name: "primary-expression", generalizations: []int{225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p182 = choiceParser{id: 182, commit: 256, name: "int", generalizations: []int{373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p173 = sequenceParser{id: 173, commit: 266, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{182, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p172 = sequenceParser{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p171 = charParser{id: 171, ranges: [][]rune{{49, 57}}}
	p172.items = []parser{&p171}
	var p166 = sequenceParser{id: 166, commit: 258, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p165 = charParser{id: 165, ranges: [][]rune{{48, 57}}}
	p166.items = []parser{&p165}
	p173.items = []parser{&p172, &p166}
	var p176 = sequenceParser{id: 176, commit: 266, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{182, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p175 = sequenceParser{id: 175, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p174 = charParser{id: 174, chars: []rune{48}}
	p175.items = []parser{&p174}
	var p168 = sequenceParser{id: 168, commit: 258, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p167 = charParser{id: 167, ranges: [][]rune{{48, 55}}}
	p168.items = []parser{&p167}
	p176.items = []parser{&p175, &p168}
	var p181 = sequenceParser{id: 181, commit: 266, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{182, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p178 = sequenceParser{id: 178, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p177 = charParser{id: 177, chars: []rune{48}}
	p178.items = []parser{&p177}
	var p180 = sequenceParser{id: 180, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p179 = charParser{id: 179, chars: []rune{120, 88}}
	p180.items = []parser{&p179}
	var p170 = sequenceParser{id: 170, commit: 258, name: "hexa-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p169 = charParser{id: 169, ranges: [][]rune{{48, 57}, {97, 102}, {65, 70}}}
	p170.items = []parser{&p169}
	p181.items = []parser{&p178, &p180, &p170}
	p182.options = []parser{&p173, &p176, &p181}
	var p195 = choiceParser{id: 195, commit: 264, name: "float", generalizations: []int{373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p190 = sequenceParser{id: 190, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{195, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{46}}
	p189.items = []parser{&p188}
	var p187 = sequenceParser{id: 187, commit: 266, name: "exponent", ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var p184 = sequenceParser{id: 184, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p183 = charParser{id: 183, chars: []rune{101, 69}}
	p184.items = []parser{&p183}
	var p186 = sequenceParser{id: 186, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p185 = charParser{id: 185, chars: []rune{43, 45}}
	p186.items = []parser{&p185}
	p187.items = []parser{&p184, &p186, &p166}
	p190.items = []parser{&p166, &p189, &p166, &p187}
	var p193 = sequenceParser{id: 193, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{195, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p192 = sequenceParser{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p191 = charParser{id: 191, chars: []rune{46}}
	p192.items = []parser{&p191}
	p193.items = []parser{&p192, &p166, &p187}
	var p194 = sequenceParser{id: 194, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{195, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	p194.items = []parser{&p166, &p187}
	p195.options = []parser{&p190, &p193, &p194}
	var p208 = sequenceParser{id: 208, commit: 264, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 225, 250, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611, 742, 772, 766, 767}}
	var p197 = sequenceParser{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p196 = charParser{id: 196, chars: []rune{34}}
	p197.items = []parser{&p196}
	var p205 = choiceParser{id: 205, commit: 10}
	var p199 = sequenceParser{id: 199, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{205}}
	var p198 = charParser{id: 198, not: true, chars: []rune{92, 34}}
	p199.items = []parser{&p198}
	var p204 = sequenceParser{id: 204, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{205}}
	var p201 = sequenceParser{id: 201, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p200 = charParser{id: 200, chars: []rune{92}}
	p201.items = []parser{&p200}
	var p203 = sequenceParser{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p202 = charParser{id: 202, not: true}
	p203.items = []parser{&p202}
	p204.items = []parser{&p201, &p203}
	p205.options = []parser{&p199, &p204}
	var p207 = sequenceParser{id: 207, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p206 = charParser{id: 206, chars: []rune{34}}
	p207.items = []parser{&p206}
	p208.items = []parser{&p197, &p205, &p207}
	var p209 = choiceParser{id: 209, commit: 258, name: "bool", generalizations: []int{373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	p209.options = []parser{&p48, &p54}
	var p214 = sequenceParser{id: 214, commit: 296, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 225, 250, 314, 630, 502, 439, 440, 441, 442, 443, 494, 618, 611, 733, 748}}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p210 = charParser{id: 210, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p211.items = []parser{&p210}
	var p213 = sequenceParser{id: 213, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p212 = charParser{id: 212, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p213.items = []parser{&p212}
	p214.items = []parser{&p211, &p213}
	var p235 = sequenceParser{id: 235, commit: 256, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{225, 373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p234 = sequenceParser{id: 234, commit: 258, name: "list-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p231 = sequenceParser{id: 231, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p230 = charParser{id: 230, chars: []rune{91}}
	p231.items = []parser{&p230}
	var p224 = sequenceParser{id: 224, commit: 258, name: "list-sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p222 = choiceParser{id: 222, commit: 2}
	var p221 = sequenceParser{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{222}}
	var p220 = charParser{id: 220, chars: []rune{44}}
	p221.items = []parser{&p220}
	p222.options = []parser{&p14, &p221}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p828, &p222}
	p224.items = []parser{&p222, &p223}
	var p229 = sequenceParser{id: 229, commit: 258, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p225 = choiceParser{id: 225, commit: 258, name: "list-item"}
	var p219 = sequenceParser{id: 219, commit: 256, name: "spread", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{225, 258, 259}}
	var p218 = sequenceParser{id: 218, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p215 = charParser{id: 215, chars: []rune{46}}
	var p216 = charParser{id: 216, chars: []rune{46}}
	var p217 = charParser{id: 217, chars: []rune{46}}
	p218.items = []parser{&p215, &p216, &p217}
	p219.items = []parser{&p373, &p828, &p218}
	p225.options = []parser{&p502, &p219}
	var p228 = sequenceParser{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p226 = sequenceParser{id: 226, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p226.items = []parser{&p224, &p828, &p225}
	var p227 = sequenceParser{id: 227, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p227.items = []parser{&p828, &p226}
	p228.items = []parser{&p828, &p226, &p227}
	p229.items = []parser{&p225, &p228}
	var p233 = sequenceParser{id: 233, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p232 = charParser{id: 232, chars: []rune{93}}
	p233.items = []parser{&p232}
	p234.items = []parser{&p231, &p828, &p224, &p828, &p229, &p828, &p224, &p828, &p233}
	p235.items = []parser{&p234}
	var p240 = sequenceParser{id: 240, commit: 256, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p237 = sequenceParser{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p236 = charParser{id: 236, chars: []rune{126}}
	p237.items = []parser{&p236}
	var p239 = sequenceParser{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p238 = sequenceParser{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p238.items = []parser{&p828, &p14}
	p239.items = []parser{&p828, &p14, &p238}
	p240.items = []parser{&p237, &p239, &p828, &p234}
	var p269 = sequenceParser{id: 269, commit: 256, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p268 = sequenceParser{id: 268, commit: 258, name: "struct-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p265 = sequenceParser{id: 265, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p264 = charParser{id: 264, chars: []rune{123}}
	p265.items = []parser{&p264}
	var p263 = sequenceParser{id: 263, commit: 258, name: "entry-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p258 = choiceParser{id: 258, commit: 2}
	var p257 = sequenceParser{id: 257, commit: 256, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{258, 259}}
	var p250 = choiceParser{id: 250, commit: 2}
	var p249 = sequenceParser{id: 249, commit: 256, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{250}}
	var p242 = sequenceParser{id: 242, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p241 = charParser{id: 241, chars: []rune{91}}
	p242.items = []parser{&p241}
	var p246 = sequenceParser{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p245 = sequenceParser{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p245.items = []parser{&p828, &p14}
	p246.items = []parser{&p828, &p14, &p245}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p247.items = []parser{&p828, &p14}
	p248.items = []parser{&p828, &p14, &p247}
	var p244 = sequenceParser{id: 244, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p243 = charParser{id: 243, chars: []rune{93}}
	p244.items = []parser{&p243}
	p249.items = []parser{&p242, &p246, &p828, &p502, &p248, &p828, &p244}
	p250.options = []parser{&p214, &p208, &p249}
	var p254 = sequenceParser{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p253 = sequenceParser{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p253.items = []parser{&p828, &p14}
	p254.items = []parser{&p828, &p14, &p253}
	var p252 = sequenceParser{id: 252, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p251 = charParser{id: 251, chars: []rune{58}}
	p252.items = []parser{&p251}
	var p256 = sequenceParser{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p255 = sequenceParser{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p255.items = []parser{&p828, &p14}
	p256.items = []parser{&p828, &p14, &p255}
	p257.items = []parser{&p250, &p254, &p828, &p252, &p256, &p828, &p502}
	p258.options = []parser{&p257, &p219}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p260 = sequenceParser{id: 260, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p259 = choiceParser{id: 259, commit: 2}
	p259.options = []parser{&p257, &p219}
	p260.items = []parser{&p224, &p828, &p259}
	var p261 = sequenceParser{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p261.items = []parser{&p828, &p260}
	p262.items = []parser{&p828, &p260, &p261}
	p263.items = []parser{&p258, &p262}
	var p267 = sequenceParser{id: 267, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p266 = charParser{id: 266, chars: []rune{125}}
	p267.items = []parser{&p266}
	p268.items = []parser{&p265, &p828, &p224, &p828, &p263, &p828, &p224, &p828, &p267}
	p269.items = []parser{&p268}
	var p274 = sequenceParser{id: 274, commit: 256, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p271 = sequenceParser{id: 271, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p270 = charParser{id: 270, chars: []rune{126}}
	p271.items = []parser{&p270}
	var p273 = sequenceParser{id: 273, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p272 = sequenceParser{id: 272, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p272.items = []parser{&p828, &p14}
	p273.items = []parser{&p828, &p14, &p272}
	p274.items = []parser{&p271, &p273, &p828, &p268}
	var p320 = sequenceParser{id: 320, commit: 256, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{314, 373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p319 = sequenceParser{id: 319, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p318 = sequenceParser{id: 318, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p318.items = []parser{&p828, &p14}
	p319.items = []parser{&p828, &p14, &p318}
	var p317 = sequenceParser{id: 317, commit: 258, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p309 = sequenceParser{id: 309, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p308 = charParser{id: 308, chars: []rune{40}}
	p309.items = []parser{&p308}
	var p311 = choiceParser{id: 311, commit: 2}
	var p278 = sequenceParser{id: 278, commit: 258, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{311}}
	var p277 = sequenceParser{id: 277, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p275 = sequenceParser{id: 275, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p275.items = []parser{&p224, &p828, &p214}
	var p276 = sequenceParser{id: 276, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p276.items = []parser{&p828, &p275}
	p277.items = []parser{&p828, &p275, &p276}
	p278.items = []parser{&p214, &p277}
	var p310 = sequenceParser{id: 310, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{311}}
	var p285 = sequenceParser{id: 285, commit: 256, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{311}}
	var p282 = sequenceParser{id: 282, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p279 = charParser{id: 279, chars: []rune{46}}
	var p280 = charParser{id: 280, chars: []rune{46}}
	var p281 = charParser{id: 281, chars: []rune{46}}
	p282.items = []parser{&p279, &p280, &p281}
	var p284 = sequenceParser{id: 284, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p283 = sequenceParser{id: 283, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p283.items = []parser{&p828, &p14}
	p284.items = []parser{&p828, &p14, &p283}
	p285.items = []parser{&p282, &p284, &p828, &p214}
	p310.items = []parser{&p278, &p828, &p224, &p828, &p285}
	p311.options = []parser{&p278, &p310, &p285}
	var p313 = sequenceParser{id: 313, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p312 = charParser{id: 312, chars: []rune{41}}
	p313.items = []parser{&p312}
	var p316 = sequenceParser{id: 316, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p315 = sequenceParser{id: 315, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p315.items = []parser{&p828, &p14}
	p316.items = []parser{&p828, &p14, &p315}
	var p314 = choiceParser{id: 314, commit: 2}
	var p293 = choiceParser{id: 293, commit: 258, name: "simple-statement", generalizations: []int{314, 798}}
	var p574 = sequenceParser{id: 574, commit: 256, name: "send-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798, 581}}
	var p571 = sequenceParser{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p570 = sequenceParser{id: 570, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p570.items = []parser{&p828, &p14}
	p571.items = []parser{&p828, &p14, &p570}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p572.items = []parser{&p828, &p14}
	p573.items = []parser{&p828, &p14, &p572}
	p574.items = []parser{&p103, &p571, &p828, &p373, &p573, &p828, &p373}
	var p606 = sequenceParser{id: 606, commit: 256, name: "go-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var p605 = sequenceParser{id: 605, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p604 = sequenceParser{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p604.items = []parser{&p828, &p14}
	p605.items = []parser{&p828, &p14, &p604}
	var p372 = sequenceParser{id: 372, commit: 256, name: "application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 798, 618, 611}}
	var p369 = sequenceParser{id: 369, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p368 = charParser{id: 368, chars: []rune{40}}
	p369.items = []parser{&p368}
	var p371 = sequenceParser{id: 371, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p370 = charParser{id: 370, chars: []rune{41}}
	p371.items = []parser{&p370}
	p372.items = []parser{&p373, &p828, &p369, &p828, &p224, &p828, &p229, &p828, &p224, &p828, &p371}
	p606.items = []parser{&p121, &p605, &p828, &p372}
	var p609 = sequenceParser{id: 609, commit: 256, name: "defer-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var p608 = sequenceParser{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p607 = sequenceParser{id: 607, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p607.items = []parser{&p828, &p14}
	p608.items = []parser{&p828, &p14, &p607}
	p609.items = []parser{&p127, &p608, &p828, &p372}
	var p637 = sequenceParser{id: 637, commit: 256, name: "assign", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var p630 = choiceParser{id: 630, commit: 2}
	var p367 = sequenceParser{id: 367, commit: 256, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{630, 373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p365 = sequenceParser{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p365.items = []parser{&p828, &p14}
	p366.items = []parser{&p828, &p14, &p365}
	var p364 = sequenceParser{id: 364, commit: 258, name: "index-list", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var p360 = choiceParser{id: 360, commit: 258, name: "index"}
	var p341 = sequenceParser{id: 341, commit: 256, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{360}}
	var p338 = sequenceParser{id: 338, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p337 = charParser{id: 337, chars: []rune{46}}
	p338.items = []parser{&p337}
	var p340 = sequenceParser{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p339 = sequenceParser{id: 339, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p339.items = []parser{&p828, &p14}
	p340.items = []parser{&p828, &p14, &p339}
	p341.items = []parser{&p338, &p340, &p828, &p214}
	var p350 = sequenceParser{id: 350, commit: 256, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{360}}
	var p343 = sequenceParser{id: 343, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p342 = charParser{id: 342, chars: []rune{91}}
	p343.items = []parser{&p342}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p346.items = []parser{&p828, &p14}
	p347.items = []parser{&p828, &p14, &p346}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p348.items = []parser{&p828, &p14}
	p349.items = []parser{&p828, &p14, &p348}
	var p345 = sequenceParser{id: 345, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p344 = charParser{id: 344, chars: []rune{93}}
	p345.items = []parser{&p344}
	p350.items = []parser{&p343, &p347, &p828, &p502, &p349, &p828, &p345}
	var p359 = sequenceParser{id: 359, commit: 256, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{360}}
	var p352 = sequenceParser{id: 352, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p351 = charParser{id: 351, chars: []rune{91}}
	p352.items = []parser{&p351}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p355.items = []parser{&p828, &p14}
	p356.items = []parser{&p828, &p14, &p355}
	var p336 = sequenceParser{id: 336, commit: 258, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{611, 617, 618}}
	var p328 = sequenceParser{id: 328, commit: 256, name: "range-from", ranges: [][]int{{1, 1}}}
	p328.items = []parser{&p502}
	var p333 = sequenceParser{id: 333, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p332 = sequenceParser{id: 332, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p332.items = []parser{&p828, &p14}
	p333.items = []parser{&p828, &p14, &p332}
	var p331 = sequenceParser{id: 331, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p330 = charParser{id: 330, chars: []rune{58}}
	p331.items = []parser{&p330}
	var p335 = sequenceParser{id: 335, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p334 = sequenceParser{id: 334, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p334.items = []parser{&p828, &p14}
	p335.items = []parser{&p828, &p14, &p334}
	var p329 = sequenceParser{id: 329, commit: 256, name: "range-to", ranges: [][]int{{1, 1}}}
	p329.items = []parser{&p502}
	p336.items = []parser{&p328, &p333, &p828, &p331, &p335, &p828, &p329}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p357.items = []parser{&p828, &p14}
	p358.items = []parser{&p828, &p14, &p357}
	var p354 = sequenceParser{id: 354, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p353 = charParser{id: 353, chars: []rune{93}}
	p354.items = []parser{&p353}
	p359.items = []parser{&p352, &p356, &p828, &p336, &p358, &p828, &p354}
	p360.options = []parser{&p341, &p350, &p359}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p361.items = []parser{&p828, &p14}
	p362.items = []parser{&p14, &p361}
	p363.items = []parser{&p362, &p828, &p360}
	p364.items = []parser{&p360, &p828, &p363}
	p367.items = []parser{&p373, &p366, &p828, &p364}
	p630.options = []parser{&p214, &p367}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p633 = sequenceParser{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p633.items = []parser{&p828, &p14}
	p634.items = []parser{&p828, &p14, &p633}
	var p632 = sequenceParser{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p631 = charParser{id: 631, chars: []rune{61}}
	p632.items = []parser{&p631}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p635.items = []parser{&p828, &p14}
	p636.items = []parser{&p828, &p14, &p635}
	p637.items = []parser{&p630, &p634, &p828, &p632, &p636, &p828, &p502}
	var p302 = sequenceParser{id: 302, commit: 258, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var p295 = sequenceParser{id: 295, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p294 = charParser{id: 294, chars: []rune{40}}
	p295.items = []parser{&p294}
	var p299 = sequenceParser{id: 299, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p298 = sequenceParser{id: 298, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p298.items = []parser{&p828, &p14}
	p299.items = []parser{&p828, &p14, &p298}
	var p301 = sequenceParser{id: 301, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p300 = sequenceParser{id: 300, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p300.items = []parser{&p828, &p14}
	p301.items = []parser{&p828, &p14, &p300}
	var p297 = sequenceParser{id: 297, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p296 = charParser{id: 296, chars: []rune{41}}
	p297.items = []parser{&p296}
	p302.items = []parser{&p295, &p299, &p828, &p293, &p301, &p828, &p297}
	p293.options = []parser{&p574, &p606, &p609, &p637, &p302}
	var p307 = sequenceParser{id: 307, commit: 256, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{314}}
	var p304 = sequenceParser{id: 304, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p303 = charParser{id: 303, chars: []rune{123}}
	p304.items = []parser{&p303}
	var p306 = sequenceParser{id: 306, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p305 = charParser{id: 305, chars: []rune{125}}
	p306.items = []parser{&p305}
	p307.items = []parser{&p304, &p828, &p812, &p828, &p816, &p828, &p812, &p828, &p306}
	p314.options = []parser{&p502, &p293, &p307}
	p317.items = []parser{&p309, &p828, &p224, &p828, &p311, &p828, &p224, &p828, &p313, &p316, &p828, &p314}
	p320.items = []parser{&p70, &p319, &p828, &p317}
	var p327 = sequenceParser{id: 327, commit: 256, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p324 = sequenceParser{id: 324, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p323 = sequenceParser{id: 323, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p323.items = []parser{&p828, &p14}
	p324.items = []parser{&p828, &p14, &p323}
	var p322 = sequenceParser{id: 322, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p321 = charParser{id: 321, chars: []rune{126}}
	p322.items = []parser{&p321}
	var p326 = sequenceParser{id: 326, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p325 = sequenceParser{id: 325, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p325.items = []parser{&p828, &p14}
	p326.items = []parser{&p828, &p14, &p325}
	p327.items = []parser{&p70, &p324, &p828, &p322, &p326, &p828, &p317}
	var p577 = sequenceParser{id: 577, commit: 256, name: "receive-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 581, 618, 611}}
	var p576 = sequenceParser{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p575 = sequenceParser{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p575.items = []parser{&p828, &p14}
	p576.items = []parser{&p828, &p14, &p575}
	p577.items = []parser{&p111, &p576, &p828, &p373}
	var p511 = sequenceParser{id: 511, commit: 258, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p504 = sequenceParser{id: 504, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p503 = charParser{id: 503, chars: []rune{40}}
	p504.items = []parser{&p503}
	var p508 = sequenceParser{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p507 = sequenceParser{id: 507, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p507.items = []parser{&p828, &p14}
	p508.items = []parser{&p828, &p14, &p507}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p509.items = []parser{&p828, &p14}
	p510.items = []parser{&p828, &p14, &p509}
	var p506 = sequenceParser{id: 506, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p505 = charParser{id: 505, chars: []rune{41}}
	p506.items = []parser{&p505}
	p511.items = []parser{&p504, &p508, &p828, &p502, &p510, &p828, &p506}
	p373.options = []parser{&p182, &p195, &p208, &p209, &p214, &p235, &p240, &p269, &p274, &p320, &p327, &p367, &p372, &p577, &p511}
	var p433 = sequenceParser{id: 433, commit: 256, name: "unary", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var p432 = choiceParser{id: 432, commit: 258, name: "unary-operator"}
	var p392 = sequenceParser{id: 392, commit: 264, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var p391 = charParser{id: 391, chars: []rune{43}}
	p392.items = []parser{&p391}
	var p394 = sequenceParser{id: 394, commit: 264, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var p393 = charParser{id: 393, chars: []rune{45}}
	p394.items = []parser{&p393}
	var p375 = sequenceParser{id: 375, commit: 264, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var p374 = charParser{id: 374, chars: []rune{94}}
	p375.items = []parser{&p374}
	var p406 = sequenceParser{id: 406, commit: 264, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var p405 = charParser{id: 405, chars: []rune{33}}
	p406.items = []parser{&p405}
	p432.options = []parser{&p392, &p394, &p375, &p406}
	p433.items = []parser{&p432, &p828, &p373}
	var p480 = choiceParser{id: 480, commit: 258, name: "binary", generalizations: []int{502, 494, 618, 611}}
	var p451 = sequenceParser{id: 451, commit: 256, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 440, 441, 442, 443, 502, 494, 618, 611}}
	var p439 = choiceParser{id: 439, commit: 258, name: "operand0", generalizations: []int{440, 441, 442, 443}}
	p439.options = []parser{&p373, &p433}
	var p449 = sequenceParser{id: 449, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p446 = sequenceParser{id: 446, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p445 = sequenceParser{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p445.items = []parser{&p828, &p14}
	p446.items = []parser{&p14, &p445}
	var p434 = choiceParser{id: 434, commit: 258, name: "binary-op0"}
	var p377 = sequenceParser{id: 377, commit: 264, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var p376 = charParser{id: 376, chars: []rune{38}}
	p377.items = []parser{&p376}
	var p384 = sequenceParser{id: 384, commit: 264, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{434}}
	var p382 = charParser{id: 382, chars: []rune{38}}
	var p383 = charParser{id: 383, chars: []rune{94}}
	p384.items = []parser{&p382, &p383}
	var p387 = sequenceParser{id: 387, commit: 264, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{434}}
	var p385 = charParser{id: 385, chars: []rune{60}}
	var p386 = charParser{id: 386, chars: []rune{60}}
	p387.items = []parser{&p385, &p386}
	var p390 = sequenceParser{id: 390, commit: 264, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{434}}
	var p388 = charParser{id: 388, chars: []rune{62}}
	var p389 = charParser{id: 389, chars: []rune{62}}
	p390.items = []parser{&p388, &p389}
	var p396 = sequenceParser{id: 396, commit: 264, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var p395 = charParser{id: 395, chars: []rune{42}}
	p396.items = []parser{&p395}
	var p398 = sequenceParser{id: 398, commit: 264, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var p397 = charParser{id: 397, chars: []rune{47}}
	p398.items = []parser{&p397}
	var p400 = sequenceParser{id: 400, commit: 264, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var p399 = charParser{id: 399, chars: []rune{37}}
	p400.items = []parser{&p399}
	p434.options = []parser{&p377, &p384, &p387, &p390, &p396, &p398, &p400}
	var p448 = sequenceParser{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p447 = sequenceParser{id: 447, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p447.items = []parser{&p828, &p14}
	p448.items = []parser{&p828, &p14, &p447}
	p449.items = []parser{&p446, &p828, &p434, &p448, &p828, &p439}
	var p450 = sequenceParser{id: 450, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p450.items = []parser{&p828, &p449}
	p451.items = []parser{&p439, &p828, &p449, &p450}
	var p458 = sequenceParser{id: 458, commit: 256, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 441, 442, 443, 502, 494, 618, 611}}
	var p440 = choiceParser{id: 440, commit: 258, name: "operand1", generalizations: []int{441, 442, 443}}
	p440.options = []parser{&p439, &p451}
	var p456 = sequenceParser{id: 456, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p453 = sequenceParser{id: 453, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p452 = sequenceParser{id: 452, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p452.items = []parser{&p828, &p14}
	p453.items = []parser{&p14, &p452}
	var p435 = choiceParser{id: 435, commit: 258, name: "binary-op1"}
	var p379 = sequenceParser{id: 379, commit: 264, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var p378 = charParser{id: 378, chars: []rune{124}}
	p379.items = []parser{&p378}
	var p381 = sequenceParser{id: 381, commit: 264, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var p380 = charParser{id: 380, chars: []rune{94}}
	p381.items = []parser{&p380}
	var p402 = sequenceParser{id: 402, commit: 264, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var p401 = charParser{id: 401, chars: []rune{43}}
	p402.items = []parser{&p401}
	var p404 = sequenceParser{id: 404, commit: 264, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var p403 = charParser{id: 403, chars: []rune{45}}
	p404.items = []parser{&p403}
	p435.options = []parser{&p379, &p381, &p402, &p404}
	var p455 = sequenceParser{id: 455, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p454 = sequenceParser{id: 454, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p454.items = []parser{&p828, &p14}
	p455.items = []parser{&p828, &p14, &p454}
	p456.items = []parser{&p453, &p828, &p435, &p455, &p828, &p440}
	var p457 = sequenceParser{id: 457, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p457.items = []parser{&p828, &p456}
	p458.items = []parser{&p440, &p828, &p456, &p457}
	var p465 = sequenceParser{id: 465, commit: 256, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 442, 443, 502, 494, 618, 611}}
	var p441 = choiceParser{id: 441, commit: 258, name: "operand2", generalizations: []int{442, 443}}
	p441.options = []parser{&p440, &p458}
	var p463 = sequenceParser{id: 463, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p459.items = []parser{&p828, &p14}
	p460.items = []parser{&p14, &p459}
	var p436 = choiceParser{id: 436, commit: 258, name: "binary-op2"}
	var p409 = sequenceParser{id: 409, commit: 264, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var p407 = charParser{id: 407, chars: []rune{61}}
	var p408 = charParser{id: 408, chars: []rune{61}}
	p409.items = []parser{&p407, &p408}
	var p412 = sequenceParser{id: 412, commit: 264, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var p410 = charParser{id: 410, chars: []rune{33}}
	var p411 = charParser{id: 411, chars: []rune{61}}
	p412.items = []parser{&p410, &p411}
	var p414 = sequenceParser{id: 414, commit: 264, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{436}}
	var p413 = charParser{id: 413, chars: []rune{60}}
	p414.items = []parser{&p413}
	var p417 = sequenceParser{id: 417, commit: 264, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var p415 = charParser{id: 415, chars: []rune{60}}
	var p416 = charParser{id: 416, chars: []rune{61}}
	p417.items = []parser{&p415, &p416}
	var p419 = sequenceParser{id: 419, commit: 264, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{436}}
	var p418 = charParser{id: 418, chars: []rune{62}}
	p419.items = []parser{&p418}
	var p422 = sequenceParser{id: 422, commit: 264, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var p420 = charParser{id: 420, chars: []rune{62}}
	var p421 = charParser{id: 421, chars: []rune{61}}
	p422.items = []parser{&p420, &p421}
	p436.options = []parser{&p409, &p412, &p414, &p417, &p419, &p422}
	var p462 = sequenceParser{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p461 = sequenceParser{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p461.items = []parser{&p828, &p14}
	p462.items = []parser{&p828, &p14, &p461}
	p463.items = []parser{&p460, &p828, &p436, &p462, &p828, &p441}
	var p464 = sequenceParser{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p464.items = []parser{&p828, &p463}
	p465.items = []parser{&p441, &p828, &p463, &p464}
	var p472 = sequenceParser{id: 472, commit: 256, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 443, 502, 494, 618, 611}}
	var p442 = choiceParser{id: 442, commit: 258, name: "operand3", generalizations: []int{443}}
	p442.options = []parser{&p441, &p465}
	var p470 = sequenceParser{id: 470, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p467 = sequenceParser{id: 467, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p466 = sequenceParser{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p466.items = []parser{&p828, &p14}
	p467.items = []parser{&p14, &p466}
	var p437 = sequenceParser{id: 437, commit: 258, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p425 = sequenceParser{id: 425, commit: 264, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p423 = charParser{id: 423, chars: []rune{38}}
	var p424 = charParser{id: 424, chars: []rune{38}}
	p425.items = []parser{&p423, &p424}
	p437.items = []parser{&p425}
	var p469 = sequenceParser{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p468 = sequenceParser{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p468.items = []parser{&p828, &p14}
	p469.items = []parser{&p828, &p14, &p468}
	p470.items = []parser{&p467, &p828, &p437, &p469, &p828, &p442}
	var p471 = sequenceParser{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p471.items = []parser{&p828, &p470}
	p472.items = []parser{&p442, &p828, &p470, &p471}
	var p479 = sequenceParser{id: 479, commit: 256, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 502, 494, 618, 611}}
	var p443 = choiceParser{id: 443, commit: 258, name: "operand4"}
	p443.options = []parser{&p442, &p472}
	var p477 = sequenceParser{id: 477, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p474 = sequenceParser{id: 474, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p473 = sequenceParser{id: 473, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p473.items = []parser{&p828, &p14}
	p474.items = []parser{&p14, &p473}
	var p438 = sequenceParser{id: 438, commit: 258, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p428 = sequenceParser{id: 428, commit: 264, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p426 = charParser{id: 426, chars: []rune{124}}
	var p427 = charParser{id: 427, chars: []rune{124}}
	p428.items = []parser{&p426, &p427}
	p438.items = []parser{&p428}
	var p476 = sequenceParser{id: 476, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p475 = sequenceParser{id: 475, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p475.items = []parser{&p828, &p14}
	p476.items = []parser{&p828, &p14, &p475}
	p477.items = []parser{&p474, &p828, &p438, &p476, &p828, &p443}
	var p478 = sequenceParser{id: 478, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p478.items = []parser{&p828, &p477}
	p479.items = []parser{&p443, &p828, &p477, &p478}
	p480.options = []parser{&p451, &p458, &p465, &p472, &p479}
	var p493 = sequenceParser{id: 493, commit: 256, name: "ternary", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{502, 494, 618, 611}}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p485.items = []parser{&p828, &p14}
	p486.items = []parser{&p828, &p14, &p485}
	var p482 = sequenceParser{id: 482, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p481 = charParser{id: 481, chars: []rune{63}}
	p482.items = []parser{&p481}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p487.items = []parser{&p828, &p14}
	p488.items = []parser{&p828, &p14, &p487}
	var p490 = sequenceParser{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p489.items = []parser{&p828, &p14}
	p490.items = []parser{&p828, &p14, &p489}
	var p484 = sequenceParser{id: 484, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p483 = charParser{id: 483, chars: []rune{58}}
	p484.items = []parser{&p483}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p491 = sequenceParser{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p491.items = []parser{&p828, &p14}
	p492.items = []parser{&p828, &p14, &p491}
	p493.items = []parser{&p502, &p486, &p828, &p482, &p488, &p828, &p502, &p490, &p828, &p484, &p492, &p828, &p502}
	var p501 = sequenceParser{id: 501, commit: 256, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{502, 798, 618, 611}}
	var p494 = choiceParser{id: 494, commit: 258, name: "chainingOperand"}
	p494.options = []parser{&p373, &p433, &p480, &p493}
	var p499 = sequenceParser{id: 499, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p496 = sequenceParser{id: 496, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p495 = sequenceParser{id: 495, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p495.items = []parser{&p828, &p14}
	p496.items = []parser{&p14, &p495}
	var p431 = sequenceParser{id: 431, commit: 266, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p429 = charParser{id: 429, chars: []rune{45}}
	var p430 = charParser{id: 430, chars: []rune{62}}
	p431.items = []parser{&p429, &p430}
	var p498 = sequenceParser{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p497 = sequenceParser{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p497.items = []parser{&p828, &p14}
	p498.items = []parser{&p828, &p14, &p497}
	p499.items = []parser{&p496, &p828, &p431, &p498, &p828, &p494}
	var p500 = sequenceParser{id: 500, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p500.items = []parser{&p828, &p499}
	p501.items = []parser{&p494, &p828, &p499, &p500}
	p502.options = []parser{&p373, &p433, &p480, &p493, &p501}
	p288.items = []parser{&p61, &p287, &p828, &p502}
	p289.options = []parser{&p61, &p288}
	var p292 = sequenceParser{id: 292, commit: 256, name: "check-ret", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var p291 = sequenceParser{id: 291, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p290 = sequenceParser{id: 290, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p290.items = []parser{&p828, &p14}
	p291.items = []parser{&p828, &p14, &p290}
	p292.items = []parser{&p67, &p291, &p828, &p502}
	var p532 = sequenceParser{id: 532, commit: 256, name: "if-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{798}}
	var p527 = sequenceParser{id: 527, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p526 = sequenceParser{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p526.items = []parser{&p828, &p14}
	p527.items = []parser{&p828, &p14, &p526}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p528 = sequenceParser{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p528.items = []parser{&p828, &p14}
	p529.items = []parser{&p828, &p14, &p528}
	var p531 = sequenceParser{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p520 = sequenceParser{id: 520, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p513 = sequenceParser{id: 513, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p512 = sequenceParser{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p512.items = []parser{&p828, &p14}
	p513.items = []parser{&p14, &p512}
	var p515 = sequenceParser{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p514 = sequenceParser{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p514.items = []parser{&p828, &p14}
	p515.items = []parser{&p828, &p14, &p514}
	var p517 = sequenceParser{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p516 = sequenceParser{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p516.items = []parser{&p828, &p14}
	p517.items = []parser{&p828, &p14, &p516}
	var p519 = sequenceParser{id: 519, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p518 = sequenceParser{id: 518, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p518.items = []parser{&p828, &p14}
	p519.items = []parser{&p828, &p14, &p518}
	p520.items = []parser{&p513, &p828, &p78, &p515, &p828, &p73, &p517, &p828, &p502, &p519, &p828, &p307}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p530.items = []parser{&p828, &p520}
	p531.items = []parser{&p828, &p520, &p530}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p522 = sequenceParser{id: 522, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p521 = sequenceParser{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p521.items = []parser{&p828, &p14}
	p522.items = []parser{&p14, &p521}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p523.items = []parser{&p828, &p14}
	p524.items = []parser{&p828, &p14, &p523}
	p525.items = []parser{&p522, &p828, &p78, &p524, &p828, &p307}
	p532.items = []parser{&p73, &p527, &p828, &p502, &p529, &p828, &p307, &p531, &p828, &p525}
	var p569 = choiceParser{id: 569, commit: 256, name: "switch-statement", generalizations: []int{798}}
	var p559 = sequenceParser{id: 559, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{569, 798}}
	var p558 = sequenceParser{id: 558, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p557 = sequenceParser{id: 557, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p557.items = []parser{&p828, &p14}
	p558.items = []parser{&p828, &p14, &p557}
	var p554 = sequenceParser{id: 554, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p553 = charParser{id: 553, chars: []rune{123}}
	p554.items = []parser{&p553}
	var p552 = choiceParser{id: 552, commit: 258, name: "cases"}
	var p548 = sequenceParser{id: 548, commit: 258, name: "case-blocks", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{552}}
	var p539 = sequenceParser{id: 539, commit: 256, name: "case-block", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p536 = sequenceParser{id: 536, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p535 = sequenceParser{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p535.items = []parser{&p828, &p14}
	p536.items = []parser{&p828, &p14, &p535}
	var p538 = sequenceParser{id: 538, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p537 = sequenceParser{id: 537, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p537.items = []parser{&p828, &p14}
	p538.items = []parser{&p828, &p14, &p537}
	var p534 = sequenceParser{id: 534, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p533 = charParser{id: 533, chars: []rune{58}}
	p534.items = []parser{&p533}
	p539.items = []parser{&p90, &p536, &p828, &p502, &p538, &p828, &p534, &p828, &p812, &p828, &p816}
	var p547 = sequenceParser{id: 547, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p545.items = []parser{&p812, &p828, &p539}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p546.items = []parser{&p828, &p545}
	p547.items = []parser{&p828, &p545, &p546}
	p548.items = []parser{&p539, &p547}
	var p544 = sequenceParser{id: 544, commit: 256, name: "default-block", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{552, 596}}
	var p543 = sequenceParser{id: 543, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p542 = sequenceParser{id: 542, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p542.items = []parser{&p828, &p14}
	p543.items = []parser{&p828, &p14, &p542}
	var p541 = sequenceParser{id: 541, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p540 = charParser{id: 540, chars: []rune{58}}
	p541.items = []parser{&p540}
	p544.items = []parser{&p98, &p543, &p828, &p541, &p828, &p812, &p828, &p816}
	var p551 = sequenceParser{id: 551, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{552}}
	var p549 = sequenceParser{id: 549, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p549.items = []parser{&p548, &p828, &p812}
	var p550 = sequenceParser{id: 550, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p550.items = []parser{&p812, &p828, &p548}
	p551.items = []parser{&p549, &p828, &p544, &p828, &p550}
	p552.options = []parser{&p548, &p544, &p551}
	var p556 = sequenceParser{id: 556, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p555 = charParser{id: 555, chars: []rune{125}}
	p556.items = []parser{&p555}
	p559.items = []parser{&p85, &p558, &p828, &p554, &p828, &p812, &p828, &p552, &p828, &p812, &p828, &p556}
	var p568 = sequenceParser{id: 568, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{569, 798}}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p564.items = []parser{&p828, &p14}
	p565.items = []parser{&p828, &p14, &p564}
	var p567 = sequenceParser{id: 567, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p566 = sequenceParser{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p566.items = []parser{&p828, &p14}
	p567.items = []parser{&p828, &p14, &p566}
	var p561 = sequenceParser{id: 561, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p560 = charParser{id: 560, chars: []rune{123}}
	p561.items = []parser{&p560}
	var p563 = sequenceParser{id: 563, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p562 = charParser{id: 562, chars: []rune{125}}
	p563.items = []parser{&p562}
	p568.items = []parser{&p85, &p565, &p828, &p502, &p567, &p828, &p561, &p828, &p812, &p828, &p552, &p828, &p812, &p828, &p563}
	p569.options = []parser{&p559, &p568}
	var p603 = sequenceParser{id: 603, commit: 256, name: "select-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var p602 = sequenceParser{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p601 = sequenceParser{id: 601, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p601.items = []parser{&p828, &p14}
	p602.items = []parser{&p828, &p14, &p601}
	var p598 = sequenceParser{id: 598, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p597 = charParser{id: 597, chars: []rune{123}}
	p598.items = []parser{&p597}
	var p596 = choiceParser{id: 596, commit: 258, name: "select-cases"}
	var p592 = sequenceParser{id: 592, commit: 258, name: "select-case-blocks", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{596}}
	var p588 = sequenceParser{id: 588, commit: 256, name: "select-case-block", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p584.items = []parser{&p828, &p14}
	p585.items = []parser{&p828, &p14, &p584}
	var p581 = choiceParser{id: 581, commit: 258, name: "communication"}
	var p580 = sequenceParser{id: 580, commit: 256, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{581}}
	var p579 = sequenceParser{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p578 = sequenceParser{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p578.items = []parser{&p828, &p14}
	p579.items = []parser{&p828, &p14, &p578}
	p580.items = []parser{&p214, &p579, &p828, &p577}
	p581.options = []parser{&p574, &p577, &p580}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p586.items = []parser{&p828, &p14}
	p587.items = []parser{&p828, &p14, &p586}
	var p583 = sequenceParser{id: 583, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p582 = charParser{id: 582, chars: []rune{58}}
	p583.items = []parser{&p582}
	p588.items = []parser{&p90, &p585, &p828, &p581, &p587, &p828, &p583, &p828, &p812, &p828, &p816}
	var p591 = sequenceParser{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p589.items = []parser{&p812, &p828, &p588}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p590.items = []parser{&p828, &p589}
	p591.items = []parser{&p828, &p589, &p590}
	p592.items = []parser{&p588, &p591}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{596}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p593.items = []parser{&p592, &p828, &p812}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p594.items = []parser{&p812, &p828, &p592}
	p595.items = []parser{&p593, &p828, &p544, &p828, &p594}
	p596.options = []parser{&p592, &p544, &p595}
	var p600 = sequenceParser{id: 600, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p599 = charParser{id: 599, chars: []rune{125}}
	p600.items = []parser{&p599}
	p603.items = []parser{&p118, &p602, &p828, &p598, &p828, &p812, &p828, &p596, &p828, &p812, &p828, &p600}
	var p610 = choiceParser{id: 610, commit: 258, name: "loop-control", generalizations: []int{798}}
	p610.options = []parser{&p140, &p149}
	var p629 = sequenceParser{id: 629, commit: 256, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var p628 = choiceParser{id: 628, commit: 2}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{628}}
	var p621 = sequenceParser{id: 621, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p620 = sequenceParser{id: 620, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p619 = sequenceParser{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p619.items = []parser{&p828, &p14}
	p620.items = []parser{&p14, &p619}
	var p618 = choiceParser{id: 618, commit: 258, name: "loop-expression"}
	var p617 = choiceParser{id: 617, commit: 256, name: "range-over", generalizations: []int{618}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{617, 618}}
	var p613 = sequenceParser{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p612 = sequenceParser{id: 612, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p612.items = []parser{&p828, &p14}
	p613.items = []parser{&p828, &p14, &p612}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p614.items = []parser{&p828, &p14}
	p615.items = []parser{&p828, &p14, &p614}
	var p611 = choiceParser{id: 611, commit: 2}
	p611.options = []parser{&p502, &p336}
	p616.items = []parser{&p214, &p613, &p828, &p130, &p615, &p828, &p611}
	p617.options = []parser{&p616, &p336}
	p618.options = []parser{&p502, &p617}
	p621.items = []parser{&p620, &p828, &p618}
	var p623 = sequenceParser{id: 623, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p622 = sequenceParser{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p622.items = []parser{&p828, &p14}
	p623.items = []parser{&p828, &p14, &p622}
	p624.items = []parser{&p621, &p623, &p828, &p307}
	var p627 = sequenceParser{id: 627, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{628}}
	var p626 = sequenceParser{id: 626, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p625 = sequenceParser{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p625.items = []parser{&p828, &p14}
	p626.items = []parser{&p14, &p625}
	p627.items = []parser{&p626, &p828, &p307}
	p628.options = []parser{&p624, &p627}
	p629.items = []parser{&p134, &p828, &p628}
	var p730 = choiceParser{id: 730, commit: 258, name: "definition", generalizations: []int{798}}
	var p657 = sequenceParser{id: 657, commit: 256, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var p653 = sequenceParser{id: 653, commit: 266, name: "docsLet", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	var p638 = sequenceParser{id: 638, commit: 256, name: "docs", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	p638.items = []parser{&p43, &p828, &p14}
	p653.items = []parser{&p638, &p153}
	var p656 = sequenceParser{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p655 = sequenceParser{id: 655, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p655.items = []parser{&p828, &p14}
	p656.items = []parser{&p828, &p14, &p655}
	var p654 = choiceParser{id: 654, commit: 2}
	var p647 = sequenceParser{id: 647, commit: 256, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{654, 659}}
	var p646 = sequenceParser{id: 646, commit: 258, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p643 = sequenceParser{id: 643, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p641 = sequenceParser{id: 641, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p641.items = []parser{&p828, &p14}
	p642.items = []parser{&p14, &p641}
	var p640 = sequenceParser{id: 640, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p639 = charParser{id: 639, chars: []rune{61}}
	p640.items = []parser{&p639}
	p643.items = []parser{&p642, &p828, &p640}
	var p645 = sequenceParser{id: 645, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p644 = sequenceParser{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p644.items = []parser{&p828, &p14}
	p645.items = []parser{&p828, &p14, &p644}
	p646.items = []parser{&p214, &p828, &p643, &p645, &p828, &p502}
	p647.items = []parser{&p646}
	var p652 = sequenceParser{id: 652, commit: 256, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{654, 659}}
	var p649 = sequenceParser{id: 649, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p648 = charParser{id: 648, chars: []rune{126}}
	p649.items = []parser{&p648}
	var p651 = sequenceParser{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p650 = sequenceParser{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p650.items = []parser{&p828, &p14}
	p651.items = []parser{&p828, &p14, &p650}
	p652.items = []parser{&p649, &p651, &p828, &p646}
	p654.options = []parser{&p647, &p652}
	p657.items = []parser{&p653, &p656, &p828, &p654}
	var p675 = sequenceParser{id: 675, commit: 256, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var p674 = sequenceParser{id: 674, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p673 = sequenceParser{id: 673, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p673.items = []parser{&p828, &p14}
	p674.items = []parser{&p828, &p14, &p673}
	var p670 = sequenceParser{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p669 = charParser{id: 669, chars: []rune{40}}
	p670.items = []parser{&p669}
	var p668 = sequenceParser{id: 668, commit: 258, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p660 = sequenceParser{id: 660, commit: 264, name: "docs-mixed-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	var p659 = choiceParser{id: 659, commit: 10}
	p659.options = []parser{&p647, &p652}
	p660.items = []parser{&p638, &p659}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p665 = sequenceParser{id: 665, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p665.items = []parser{&p224, &p828, &p660}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p666.items = []parser{&p828, &p665}
	p667.items = []parser{&p828, &p665, &p666}
	p668.items = []parser{&p660, &p667}
	var p672 = sequenceParser{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p671 = charParser{id: 671, chars: []rune{41}}
	p672.items = []parser{&p671}
	p675.items = []parser{&p153, &p674, &p828, &p670, &p828, &p224, &p828, &p668, &p828, &p224, &p828, &p672}
	var p686 = sequenceParser{id: 686, commit: 256, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var p683 = sequenceParser{id: 683, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p682 = sequenceParser{id: 682, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p682.items = []parser{&p828, &p14}
	p683.items = []parser{&p828, &p14, &p682}
	var p677 = sequenceParser{id: 677, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p676 = charParser{id: 676, chars: []rune{126}}
	p677.items = []parser{&p676}
	var p685 = sequenceParser{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p684 = sequenceParser{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p684.items = []parser{&p828, &p14}
	p685.items = []parser{&p828, &p14, &p684}
	var p679 = sequenceParser{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p678 = charParser{id: 678, chars: []rune{40}}
	p679.items = []parser{&p678}
	var p664 = sequenceParser{id: 664, commit: 258, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p658 = sequenceParser{id: 658, commit: 264, name: "docs-value-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	p658.items = []parser{&p638, &p647}
	var p663 = sequenceParser{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p661 = sequenceParser{id: 661, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p661.items = []parser{&p224, &p828, &p658}
	var p662 = sequenceParser{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p662.items = []parser{&p828, &p661}
	p663.items = []parser{&p828, &p661, &p662}
	p664.items = []parser{&p658, &p663}
	var p681 = sequenceParser{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{41}}
	p681.items = []parser{&p680}
	p686.items = []parser{&p153, &p683, &p828, &p677, &p685, &p828, &p679, &p828, &p224, &p828, &p664, &p828, &p224, &p828, &p681}
	var p700 = sequenceParser{id: 700, commit: 256, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var p696 = sequenceParser{id: 696, commit: 266, name: "docsFn", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	p696.items = []parser{&p638, &p70}
	var p699 = sequenceParser{id: 699, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p698 = sequenceParser{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p698.items = []parser{&p828, &p14}
	p699.items = []parser{&p828, &p14, &p698}
	var p697 = choiceParser{id: 697, commit: 2}
	var p690 = sequenceParser{id: 690, commit: 256, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{697, 702}}
	var p689 = sequenceParser{id: 689, commit: 258, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p687.items = []parser{&p828, &p14}
	p688.items = []parser{&p828, &p14, &p687}
	p689.items = []parser{&p214, &p688, &p828, &p317}
	p690.items = []parser{&p689}
	var p695 = sequenceParser{id: 695, commit: 256, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{697, 702}}
	var p692 = sequenceParser{id: 692, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p691 = charParser{id: 691, chars: []rune{126}}
	p692.items = []parser{&p691}
	var p694 = sequenceParser{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p693 = sequenceParser{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p693.items = []parser{&p828, &p14}
	p694.items = []parser{&p828, &p14, &p693}
	p695.items = []parser{&p692, &p694, &p828, &p689}
	p697.options = []parser{&p690, &p695}
	p700.items = []parser{&p696, &p699, &p828, &p697}
	var p718 = sequenceParser{id: 718, commit: 256, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var p717 = sequenceParser{id: 717, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p716 = sequenceParser{id: 716, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p716.items = []parser{&p828, &p14}
	p717.items = []parser{&p828, &p14, &p716}
	var p713 = sequenceParser{id: 713, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p712 = charParser{id: 712, chars: []rune{40}}
	p713.items = []parser{&p712}
	var p711 = sequenceParser{id: 711, commit: 258, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p703 = sequenceParser{id: 703, commit: 264, name: "docs-mixed-function-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	var p702 = choiceParser{id: 702, commit: 10}
	p702.options = []parser{&p690, &p695}
	p703.items = []parser{&p638, &p702}
	var p710 = sequenceParser{id: 710, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p708 = sequenceParser{id: 708, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p708.items = []parser{&p224, &p828, &p703}
	var p709 = sequenceParser{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p709.items = []parser{&p828, &p708}
	p710.items = []parser{&p828, &p708, &p709}
	p711.items = []parser{&p703, &p710}
	var p715 = sequenceParser{id: 715, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p714 = charParser{id: 714, chars: []rune{41}}
	p715.items = []parser{&p714}
	p718.items = []parser{&p70, &p717, &p828, &p713, &p828, &p224, &p828, &p711, &p828, &p224, &p828, &p715}
	var p729 = sequenceParser{id: 729, commit: 256, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var p726 = sequenceParser{id: 726, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p725 = sequenceParser{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p725.items = []parser{&p828, &p14}
	p726.items = []parser{&p828, &p14, &p725}
	var p720 = sequenceParser{id: 720, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p719 = charParser{id: 719, chars: []rune{126}}
	p720.items = []parser{&p719}
	var p728 = sequenceParser{id: 728, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p727 = sequenceParser{id: 727, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p727.items = []parser{&p828, &p14}
	p728.items = []parser{&p828, &p14, &p727}
	var p722 = sequenceParser{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p721 = charParser{id: 721, chars: []rune{40}}
	p722.items = []parser{&p721}
	var p707 = sequenceParser{id: 707, commit: 258, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p701 = sequenceParser{id: 701, commit: 264, name: "docs-function-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	p701.items = []parser{&p638, &p690}
	var p706 = sequenceParser{id: 706, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p704.items = []parser{&p224, &p828, &p701}
	var p705 = sequenceParser{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p705.items = []parser{&p828, &p704}
	p706.items = []parser{&p828, &p704, &p705}
	p707.items = []parser{&p701, &p706}
	var p724 = sequenceParser{id: 724, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p723 = charParser{id: 723, chars: []rune{41}}
	p724.items = []parser{&p723}
	p729.items = []parser{&p70, &p726, &p828, &p720, &p728, &p828, &p722, &p828, &p224, &p828, &p707, &p828, &p224, &p828, &p724}
	p730.options = []parser{&p657, &p675, &p686, &p700, &p718, &p729}
	var p797 = sequenceParser{id: 797, commit: 256, name: "export-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var p796 = sequenceParser{id: 796, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p795 = sequenceParser{id: 795, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p795.items = []parser{&p828, &p14}
	p796.items = []parser{&p828, &p14, &p795}
	p797.items = []parser{&p160, &p796, &p828, &p730}
	var p794 = choiceParser{id: 794, commit: 256, name: "use-modules", generalizations: []int{798}}
	var p775 = sequenceParser{id: 775, commit: 258, name: "use-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 798}}
	var p774 = sequenceParser{id: 774, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p773 = sequenceParser{id: 773, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p773.items = []parser{&p828, &p14}
	p774.items = []parser{&p828, &p14, &p773}
	var p772 = choiceParser{id: 772, commit: 2}
	var p742 = choiceParser{id: 742, commit: 256, name: "use-fact", generalizations: []int{772, 766, 767}}
	var p741 = sequenceParser{id: 741, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{742, 772, 766, 767}}
	var p733 = choiceParser{id: 733, commit: 2}
	var p732 = sequenceParser{id: 732, commit: 264, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{733, 748}}
	var p731 = charParser{id: 731, chars: []rune{46}}
	p732.items = []parser{&p731}
	p733.options = []parser{&p214, &p732}
	var p738 = sequenceParser{id: 738, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p736.items = []parser{&p828, &p14}
	p737.items = []parser{&p14, &p736}
	var p735 = sequenceParser{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p734 = charParser{id: 734, chars: []rune{61}}
	p735.items = []parser{&p734}
	p738.items = []parser{&p737, &p828, &p735}
	var p740 = sequenceParser{id: 740, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p739 = sequenceParser{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p739.items = []parser{&p828, &p14}
	p740.items = []parser{&p828, &p14, &p739}
	p741.items = []parser{&p733, &p828, &p738, &p740, &p828, &p208}
	p742.options = []parser{&p208, &p741}
	var p761 = choiceParser{id: 761, commit: 256, name: "use-effect", generalizations: []int{772, 766, 767}}
	var p747 = sequenceParser{id: 747, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{761, 772, 766, 767}}
	var p744 = sequenceParser{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p743 = charParser{id: 743, chars: []rune{126}}
	p744.items = []parser{&p743}
	var p746 = sequenceParser{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p745 = sequenceParser{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p745.items = []parser{&p828, &p14}
	p746.items = []parser{&p828, &p14, &p745}
	p747.items = []parser{&p744, &p746, &p828, &p208}
	var p760 = sequenceParser{id: 760, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{761, 772, 766, 767}}
	var p748 = choiceParser{id: 748, commit: 2}
	p748.options = []parser{&p214, &p732}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p751 = sequenceParser{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p751.items = []parser{&p828, &p14}
	p752.items = []parser{&p14, &p751}
	var p750 = sequenceParser{id: 750, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p749 = charParser{id: 749, chars: []rune{61}}
	p750.items = []parser{&p749}
	p753.items = []parser{&p752, &p828, &p750}
	var p757 = sequenceParser{id: 757, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p756.items = []parser{&p828, &p14}
	p757.items = []parser{&p828, &p14, &p756}
	var p755 = sequenceParser{id: 755, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p754 = charParser{id: 754, chars: []rune{126}}
	p755.items = []parser{&p754}
	var p759 = sequenceParser{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p758 = sequenceParser{id: 758, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p758.items = []parser{&p828, &p14}
	p759.items = []parser{&p828, &p14, &p758}
	p760.items = []parser{&p748, &p828, &p753, &p757, &p828, &p755, &p759, &p828, &p208}
	p761.options = []parser{&p747, &p760}
	p772.options = []parser{&p742, &p761}
	p775.items = []parser{&p164, &p774, &p828, &p772}
	var p782 = sequenceParser{id: 782, commit: 258, name: "use-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 798}}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p780 = sequenceParser{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p780.items = []parser{&p828, &p14}
	p781.items = []parser{&p828, &p14, &p780}
	var p777 = sequenceParser{id: 777, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p776 = charParser{id: 776, chars: []rune{40}}
	p777.items = []parser{&p776}
	var p771 = sequenceParser{id: 771, commit: 258, name: "use-mixed-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p766 = choiceParser{id: 766, commit: 2}
	p766.options = []parser{&p742, &p761}
	var p770 = sequenceParser{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p768 = sequenceParser{id: 768, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p767 = choiceParser{id: 767, commit: 2}
	p767.options = []parser{&p742, &p761}
	p768.items = []parser{&p224, &p828, &p767}
	var p769 = sequenceParser{id: 769, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p769.items = []parser{&p828, &p768}
	p770.items = []parser{&p828, &p768, &p769}
	p771.items = []parser{&p766, &p770}
	var p779 = sequenceParser{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p778 = charParser{id: 778, chars: []rune{41}}
	p779.items = []parser{&p778}
	p782.items = []parser{&p164, &p781, &p828, &p777, &p828, &p224, &p828, &p771, &p828, &p224, &p828, &p779}
	var p793 = sequenceParser{id: 793, commit: 258, name: "use-effect-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 798}}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p789 = sequenceParser{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p789.items = []parser{&p828, &p14}
	p790.items = []parser{&p828, &p14, &p789}
	var p784 = sequenceParser{id: 784, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p783 = charParser{id: 783, chars: []rune{126}}
	p784.items = []parser{&p783}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p791.items = []parser{&p828, &p14}
	p792.items = []parser{&p828, &p14, &p791}
	var p786 = sequenceParser{id: 786, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p785 = charParser{id: 785, chars: []rune{40}}
	p786.items = []parser{&p785}
	var p765 = sequenceParser{id: 765, commit: 258, name: "use-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p764 = sequenceParser{id: 764, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p762 = sequenceParser{id: 762, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p762.items = []parser{&p224, &p828, &p742}
	var p763 = sequenceParser{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p763.items = []parser{&p828, &p762}
	p764.items = []parser{&p828, &p762, &p763}
	p765.items = []parser{&p742, &p764}
	var p788 = sequenceParser{id: 788, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p787 = charParser{id: 787, chars: []rune{41}}
	p788.items = []parser{&p787}
	p793.items = []parser{&p164, &p790, &p828, &p784, &p792, &p828, &p786, &p828, &p224, &p828, &p765, &p828, &p224, &p828, &p788}
	p794.options = []parser{&p775, &p782, &p793}
	var p807 = sequenceParser{id: 807, commit: 258, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var p800 = sequenceParser{id: 800, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p799 = charParser{id: 799, chars: []rune{40}}
	p800.items = []parser{&p799}
	var p804 = sequenceParser{id: 804, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p803 = sequenceParser{id: 803, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p803.items = []parser{&p828, &p14}
	p804.items = []parser{&p828, &p14, &p803}
	var p806 = sequenceParser{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p805 = sequenceParser{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p805.items = []parser{&p828, &p14}
	p806.items = []parser{&p828, &p14, &p805}
	var p802 = sequenceParser{id: 802, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p801 = charParser{id: 801, chars: []rune{41}}
	p802.items = []parser{&p801}
	p807.items = []parser{&p800, &p804, &p828, &p798, &p806, &p828, &p802}
	p798.options = []parser{&p289, &p292, &p372, &p501, &p532, &p569, &p574, &p603, &p606, &p609, &p610, &p629, &p730, &p797, &p794, &p807, &p293}
	var p815 = sequenceParser{id: 815, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p813 = sequenceParser{id: 813, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p813.items = []parser{&p812, &p828, &p798}
	var p814 = sequenceParser{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p814.items = []parser{&p828, &p813}
	p815.items = []parser{&p828, &p813, &p814}
	p816.items = []parser{&p798, &p815}
	p829.items = []parser{&p825, &p828, &p812, &p828, &p816, &p828, &p812}
	p830.items = []parser{&p828, &p829, &p828}
	var b830 = sequenceBuilder{id: 830, commit: 128, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b828 = choiceBuilder{id: 828, commit: 2}
	var b826 = choiceBuilder{id: 826, commit: 262, generalizations: []int{828}}
	var b2 = sequenceBuilder{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var b1 = charBuilder{}
	b2.items = []builder{&b1}
	var b4 = sequenceBuilder{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var b3 = charBuilder{}
	b4.items = []builder{&b3}
	var b6 = sequenceBuilder{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var b5 = charBuilder{}
	b6.items = []builder{&b5}
	var b8 = sequenceBuilder{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var b7 = charBuilder{}
	b8.items = []builder{&b7}
	var b10 = sequenceBuilder{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var b9 = charBuilder{}
	b10.items = []builder{&b9}
	var b12 = sequenceBuilder{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826, 828}}
	var b11 = charBuilder{}
	b12.items = []builder{&b11}
	b826.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b827 = sequenceBuilder{id: 827, commit: 262, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{828}}
	var b43 = choiceBuilder{id: 43, commit: 258}
	var b26 = sequenceBuilder{id: 26, commit: 256, name: "line-comment", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{43}}
	var b22 = sequenceBuilder{id: 22, commit: 266, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b21 = sequenceBuilder{id: 21, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b19 = charBuilder{}
	var b20 = charBuilder{}
	b21.items = []builder{&b19, &b20}
	var b18 = sequenceBuilder{id: 18, commit: 264, name: "line-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var b17 = sequenceBuilder{id: 17, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b16 = charBuilder{}
	b17.items = []builder{&b16}
	b18.items = []builder{&b17}
	b22.items = []builder{&b21, &b18}
	var b25 = sequenceBuilder{id: 25, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b23 = sequenceBuilder{id: 23, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b14 = sequenceBuilder{id: 14, commit: 266, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{810, 222}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b23.items = []builder{&b14, &b828, &b22}
	var b24 = sequenceBuilder{id: 24, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b24.items = []builder{&b828, &b23}
	b25.items = []builder{&b828, &b23, &b24}
	b26.items = []builder{&b22, &b25}
	var b42 = sequenceBuilder{id: 42, commit: 264, name: "block-comment", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{43}}
	var b38 = sequenceBuilder{id: 38, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b36 = charBuilder{}
	var b37 = charBuilder{}
	b38.items = []builder{&b36, &b37}
	var b35 = sequenceBuilder{id: 35, commit: 264, name: "block-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var b34 = choiceBuilder{id: 34, commit: 10}
	var b28 = sequenceBuilder{id: 28, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{34}}
	var b27 = charBuilder{}
	b28.items = []builder{&b27}
	var b33 = sequenceBuilder{id: 33, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{34}}
	var b30 = sequenceBuilder{id: 30, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b29 = charBuilder{}
	b30.items = []builder{&b29}
	var b32 = sequenceBuilder{id: 32, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b31 = charBuilder{}
	b32.items = []builder{&b31}
	b33.items = []builder{&b30, &b32}
	b34.options = []builder{&b28, &b33}
	b35.items = []builder{&b34}
	var b41 = sequenceBuilder{id: 41, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b39 = charBuilder{}
	var b40 = charBuilder{}
	b41.items = []builder{&b39, &b40}
	b42.items = []builder{&b38, &b35, &b41}
	b43.options = []builder{&b26, &b42}
	b827.items = []builder{&b43}
	b828.options = []builder{&b826, &b827}
	var b829 = sequenceBuilder{id: 829, commit: 258, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b825 = sequenceBuilder{id: 825, commit: 256, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b824 = sequenceBuilder{id: 824, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b822 = charBuilder{}
	var b823 = charBuilder{}
	b824.items = []builder{&b822, &b823}
	var b821 = sequenceBuilder{id: 821, commit: 256, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b820 = sequenceBuilder{id: 820, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b818 = sequenceBuilder{id: 818, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b817 = charBuilder{}
	b818.items = []builder{&b817}
	var b819 = sequenceBuilder{id: 819, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b819.items = []builder{&b828, &b818}
	b820.items = []builder{&b818, &b819}
	b821.items = []builder{&b820}
	b825.items = []builder{&b824, &b828, &b821, &b828, &b14}
	var b812 = sequenceBuilder{id: 812, commit: 258, ranges: [][]int{{1, 1}, {0, -1}}}
	var b810 = choiceBuilder{id: 810, commit: 2}
	var b809 = sequenceBuilder{id: 809, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{810}}
	var b808 = charBuilder{}
	b809.items = []builder{&b808}
	b810.options = []builder{&b809, &b14}
	var b811 = sequenceBuilder{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b811.items = []builder{&b828, &b810}
	b812.items = []builder{&b810, &b811}
	var b816 = sequenceBuilder{id: 816, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b798 = choiceBuilder{id: 798, commit: 258}
	var b289 = choiceBuilder{id: 289, commit: 256, name: "ret", generalizations: []int{798}}
	var b61 = sequenceBuilder{id: 61, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{289, 798}}
	var b55 = charBuilder{}
	var b56 = charBuilder{}
	var b57 = charBuilder{}
	var b58 = charBuilder{}
	var b59 = charBuilder{}
	var b60 = charBuilder{}
	b61.items = []builder{&b55, &b56, &b57, &b58, &b59, &b60}
	var b288 = sequenceBuilder{id: 288, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{289, 798}}
	var b287 = sequenceBuilder{id: 287, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b286 = sequenceBuilder{id: 286, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b286.items = []builder{&b828, &b14}
	b287.items = []builder{&b828, &b14, &b286}
	var b502 = choiceBuilder{id: 502, commit: 258, generalizations: []int{225, 314, 618, 611}}
	var b373 = choiceBuilder{id: 373, commit: 258, generalizations: []int{225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b182 = choiceBuilder{id: 182, commit: 256, name: "int", generalizations: []int{373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b173 = sequenceBuilder{id: 173, commit: 266, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{182, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b172 = sequenceBuilder{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b171 = charBuilder{}
	b172.items = []builder{&b171}
	var b166 = sequenceBuilder{id: 166, commit: 258, allChars: true, ranges: [][]int{{1, 1}}}
	var b165 = charBuilder{}
	b166.items = []builder{&b165}
	b173.items = []builder{&b172, &b166}
	var b176 = sequenceBuilder{id: 176, commit: 266, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{182, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b175 = sequenceBuilder{id: 175, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b174 = charBuilder{}
	b175.items = []builder{&b174}
	var b168 = sequenceBuilder{id: 168, commit: 258, allChars: true, ranges: [][]int{{1, 1}}}
	var b167 = charBuilder{}
	b168.items = []builder{&b167}
	b176.items = []builder{&b175, &b168}
	var b181 = sequenceBuilder{id: 181, commit: 266, ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{182, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b178 = sequenceBuilder{id: 178, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b177 = charBuilder{}
	b178.items = []builder{&b177}
	var b180 = sequenceBuilder{id: 180, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b179 = charBuilder{}
	b180.items = []builder{&b179}
	var b170 = sequenceBuilder{id: 170, commit: 258, allChars: true, ranges: [][]int{{1, 1}}}
	var b169 = charBuilder{}
	b170.items = []builder{&b169}
	b181.items = []builder{&b178, &b180, &b170}
	b182.options = []builder{&b173, &b176, &b181}
	var b195 = choiceBuilder{id: 195, commit: 264, name: "float", generalizations: []int{373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b190 = sequenceBuilder{id: 190, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{195, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	var b187 = sequenceBuilder{id: 187, commit: 266, ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var b184 = sequenceBuilder{id: 184, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b183 = charBuilder{}
	b184.items = []builder{&b183}
	var b186 = sequenceBuilder{id: 186, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b185 = charBuilder{}
	b186.items = []builder{&b185}
	b187.items = []builder{&b184, &b186, &b166}
	b190.items = []builder{&b166, &b189, &b166, &b187}
	var b193 = sequenceBuilder{id: 193, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{195, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b192 = sequenceBuilder{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b191 = charBuilder{}
	b192.items = []builder{&b191}
	b193.items = []builder{&b192, &b166, &b187}
	var b194 = sequenceBuilder{id: 194, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{195, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	b194.items = []builder{&b166, &b187}
	b195.options = []builder{&b190, &b193, &b194}
	var b208 = sequenceBuilder{id: 208, commit: 264, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 225, 250, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611, 742, 772, 766, 767}}
	var b197 = sequenceBuilder{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b196 = charBuilder{}
	b197.items = []builder{&b196}
	var b205 = choiceBuilder{id: 205, commit: 10}
	var b199 = sequenceBuilder{id: 199, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{205}}
	var b198 = charBuilder{}
	b199.items = []builder{&b198}
	var b204 = sequenceBuilder{id: 204, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{205}}
	var b201 = sequenceBuilder{id: 201, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b200 = charBuilder{}
	b201.items = []builder{&b200}
	var b203 = sequenceBuilder{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b202 = charBuilder{}
	b203.items = []builder{&b202}
	b204.items = []builder{&b201, &b203}
	b205.options = []builder{&b199, &b204}
	var b207 = sequenceBuilder{id: 207, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b206 = charBuilder{}
	b207.items = []builder{&b206}
	b208.items = []builder{&b197, &b205, &b207}
	var b209 = choiceBuilder{id: 209, commit: 258, generalizations: []int{373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b48 = sequenceBuilder{id: 48, commit: 280, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{209, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b44 = charBuilder{}
	var b45 = charBuilder{}
	var b46 = charBuilder{}
	var b47 = charBuilder{}
	b48.items = []builder{&b44, &b45, &b46, &b47}
	var b54 = sequenceBuilder{id: 54, commit: 280, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{209, 373, 225, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b49 = charBuilder{}
	var b50 = charBuilder{}
	var b51 = charBuilder{}
	var b52 = charBuilder{}
	var b53 = charBuilder{}
	b54.items = []builder{&b49, &b50, &b51, &b52, &b53}
	b209.options = []builder{&b48, &b54}
	var b214 = sequenceBuilder{id: 214, commit: 296, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 225, 250, 314, 630, 502, 439, 440, 441, 442, 443, 494, 618, 611, 733, 748}}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b210 = charBuilder{}
	b211.items = []builder{&b210}
	var b213 = sequenceBuilder{id: 213, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b212 = charBuilder{}
	b213.items = []builder{&b212}
	b214.items = []builder{&b211, &b213}
	var b235 = sequenceBuilder{id: 235, commit: 256, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{225, 373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b234 = sequenceBuilder{id: 234, commit: 258, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b231 = sequenceBuilder{id: 231, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b230 = charBuilder{}
	b231.items = []builder{&b230}
	var b224 = sequenceBuilder{id: 224, commit: 258, ranges: [][]int{{1, 1}, {0, -1}}}
	var b222 = choiceBuilder{id: 222, commit: 2}
	var b221 = sequenceBuilder{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{222}}
	var b220 = charBuilder{}
	b221.items = []builder{&b220}
	b222.options = []builder{&b14, &b221}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b828, &b222}
	b224.items = []builder{&b222, &b223}
	var b229 = sequenceBuilder{id: 229, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b225 = choiceBuilder{id: 225, commit: 258}
	var b219 = sequenceBuilder{id: 219, commit: 256, name: "spread", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{225, 258, 259}}
	var b218 = sequenceBuilder{id: 218, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b215 = charBuilder{}
	var b216 = charBuilder{}
	var b217 = charBuilder{}
	b218.items = []builder{&b215, &b216, &b217}
	b219.items = []builder{&b373, &b828, &b218}
	b225.options = []builder{&b502, &b219}
	var b228 = sequenceBuilder{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b226 = sequenceBuilder{id: 226, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b226.items = []builder{&b224, &b828, &b225}
	var b227 = sequenceBuilder{id: 227, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b227.items = []builder{&b828, &b226}
	b228.items = []builder{&b828, &b226, &b227}
	b229.items = []builder{&b225, &b228}
	var b233 = sequenceBuilder{id: 233, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b232 = charBuilder{}
	b233.items = []builder{&b232}
	b234.items = []builder{&b231, &b828, &b224, &b828, &b229, &b828, &b224, &b828, &b233}
	b235.items = []builder{&b234}
	var b240 = sequenceBuilder{id: 240, commit: 256, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b237 = sequenceBuilder{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b236 = charBuilder{}
	b237.items = []builder{&b236}
	var b239 = sequenceBuilder{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b238 = sequenceBuilder{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b238.items = []builder{&b828, &b14}
	b239.items = []builder{&b828, &b14, &b238}
	b240.items = []builder{&b237, &b239, &b828, &b234}
	var b269 = sequenceBuilder{id: 269, commit: 256, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b268 = sequenceBuilder{id: 268, commit: 258, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b265 = sequenceBuilder{id: 265, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b264 = charBuilder{}
	b265.items = []builder{&b264}
	var b263 = sequenceBuilder{id: 263, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b258 = choiceBuilder{id: 258, commit: 2}
	var b257 = sequenceBuilder{id: 257, commit: 256, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{258, 259}}
	var b250 = choiceBuilder{id: 250, commit: 2}
	var b249 = sequenceBuilder{id: 249, commit: 256, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{250}}
	var b242 = sequenceBuilder{id: 242, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b241 = charBuilder{}
	b242.items = []builder{&b241}
	var b246 = sequenceBuilder{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b245 = sequenceBuilder{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b245.items = []builder{&b828, &b14}
	b246.items = []builder{&b828, &b14, &b245}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b247.items = []builder{&b828, &b14}
	b248.items = []builder{&b828, &b14, &b247}
	var b244 = sequenceBuilder{id: 244, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b243 = charBuilder{}
	b244.items = []builder{&b243}
	b249.items = []builder{&b242, &b246, &b828, &b502, &b248, &b828, &b244}
	b250.options = []builder{&b214, &b208, &b249}
	var b254 = sequenceBuilder{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b253 = sequenceBuilder{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b253.items = []builder{&b828, &b14}
	b254.items = []builder{&b828, &b14, &b253}
	var b252 = sequenceBuilder{id: 252, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b251 = charBuilder{}
	b252.items = []builder{&b251}
	var b256 = sequenceBuilder{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b255 = sequenceBuilder{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b255.items = []builder{&b828, &b14}
	b256.items = []builder{&b828, &b14, &b255}
	b257.items = []builder{&b250, &b254, &b828, &b252, &b256, &b828, &b502}
	b258.options = []builder{&b257, &b219}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b260 = sequenceBuilder{id: 260, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b259 = choiceBuilder{id: 259, commit: 2}
	b259.options = []builder{&b257, &b219}
	b260.items = []builder{&b224, &b828, &b259}
	var b261 = sequenceBuilder{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b261.items = []builder{&b828, &b260}
	b262.items = []builder{&b828, &b260, &b261}
	b263.items = []builder{&b258, &b262}
	var b267 = sequenceBuilder{id: 267, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b266 = charBuilder{}
	b267.items = []builder{&b266}
	b268.items = []builder{&b265, &b828, &b224, &b828, &b263, &b828, &b224, &b828, &b267}
	b269.items = []builder{&b268}
	var b274 = sequenceBuilder{id: 274, commit: 256, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 314, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b271 = sequenceBuilder{id: 271, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b270 = charBuilder{}
	b271.items = []builder{&b270}
	var b273 = sequenceBuilder{id: 273, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b272 = sequenceBuilder{id: 272, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b272.items = []builder{&b828, &b14}
	b273.items = []builder{&b828, &b14, &b272}
	b274.items = []builder{&b271, &b273, &b828, &b268}
	var b320 = sequenceBuilder{id: 320, commit: 256, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{314, 373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b70 = sequenceBuilder{id: 70, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b68 = charBuilder{}
	var b69 = charBuilder{}
	b70.items = []builder{&b68, &b69}
	var b319 = sequenceBuilder{id: 319, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b318 = sequenceBuilder{id: 318, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b318.items = []builder{&b828, &b14}
	b319.items = []builder{&b828, &b14, &b318}
	var b317 = sequenceBuilder{id: 317, commit: 258, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b309 = sequenceBuilder{id: 309, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b308 = charBuilder{}
	b309.items = []builder{&b308}
	var b311 = choiceBuilder{id: 311, commit: 2}
	var b278 = sequenceBuilder{id: 278, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{311}}
	var b277 = sequenceBuilder{id: 277, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b275 = sequenceBuilder{id: 275, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b275.items = []builder{&b224, &b828, &b214}
	var b276 = sequenceBuilder{id: 276, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b276.items = []builder{&b828, &b275}
	b277.items = []builder{&b828, &b275, &b276}
	b278.items = []builder{&b214, &b277}
	var b310 = sequenceBuilder{id: 310, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{311}}
	var b285 = sequenceBuilder{id: 285, commit: 256, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{311}}
	var b282 = sequenceBuilder{id: 282, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b279 = charBuilder{}
	var b280 = charBuilder{}
	var b281 = charBuilder{}
	b282.items = []builder{&b279, &b280, &b281}
	var b284 = sequenceBuilder{id: 284, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b283 = sequenceBuilder{id: 283, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b283.items = []builder{&b828, &b14}
	b284.items = []builder{&b828, &b14, &b283}
	b285.items = []builder{&b282, &b284, &b828, &b214}
	b310.items = []builder{&b278, &b828, &b224, &b828, &b285}
	b311.options = []builder{&b278, &b310, &b285}
	var b313 = sequenceBuilder{id: 313, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b312 = charBuilder{}
	b313.items = []builder{&b312}
	var b316 = sequenceBuilder{id: 316, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b315 = sequenceBuilder{id: 315, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b315.items = []builder{&b828, &b14}
	b316.items = []builder{&b828, &b14, &b315}
	var b314 = choiceBuilder{id: 314, commit: 2}
	var b293 = choiceBuilder{id: 293, commit: 258, generalizations: []int{314, 798}}
	var b574 = sequenceBuilder{id: 574, commit: 256, name: "send-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798, 581}}
	var b103 = sequenceBuilder{id: 103, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b99 = charBuilder{}
	var b100 = charBuilder{}
	var b101 = charBuilder{}
	var b102 = charBuilder{}
	b103.items = []builder{&b99, &b100, &b101, &b102}
	var b571 = sequenceBuilder{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b570 = sequenceBuilder{id: 570, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b570.items = []builder{&b828, &b14}
	b571.items = []builder{&b828, &b14, &b570}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b572.items = []builder{&b828, &b14}
	b573.items = []builder{&b828, &b14, &b572}
	b574.items = []builder{&b103, &b571, &b828, &b373, &b573, &b828, &b373}
	var b606 = sequenceBuilder{id: 606, commit: 256, name: "go-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var b121 = sequenceBuilder{id: 121, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b119 = charBuilder{}
	var b120 = charBuilder{}
	b121.items = []builder{&b119, &b120}
	var b605 = sequenceBuilder{id: 605, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b604 = sequenceBuilder{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b604.items = []builder{&b828, &b14}
	b605.items = []builder{&b828, &b14, &b604}
	var b372 = sequenceBuilder{id: 372, commit: 256, name: "application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 798, 618, 611}}
	var b369 = sequenceBuilder{id: 369, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b368 = charBuilder{}
	b369.items = []builder{&b368}
	var b371 = sequenceBuilder{id: 371, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b370 = charBuilder{}
	b371.items = []builder{&b370}
	b372.items = []builder{&b373, &b828, &b369, &b828, &b224, &b828, &b229, &b828, &b224, &b828, &b371}
	b606.items = []builder{&b121, &b605, &b828, &b372}
	var b609 = sequenceBuilder{id: 609, commit: 256, name: "defer-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var b127 = sequenceBuilder{id: 127, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b122 = charBuilder{}
	var b123 = charBuilder{}
	var b124 = charBuilder{}
	var b125 = charBuilder{}
	var b126 = charBuilder{}
	b127.items = []builder{&b122, &b123, &b124, &b125, &b126}
	var b608 = sequenceBuilder{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b607 = sequenceBuilder{id: 607, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b607.items = []builder{&b828, &b14}
	b608.items = []builder{&b828, &b14, &b607}
	b609.items = []builder{&b127, &b608, &b828, &b372}
	var b637 = sequenceBuilder{id: 637, commit: 256, name: "assign", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var b630 = choiceBuilder{id: 630, commit: 2}
	var b367 = sequenceBuilder{id: 367, commit: 256, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{630, 373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b365 = sequenceBuilder{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b365.items = []builder{&b828, &b14}
	b366.items = []builder{&b828, &b14, &b365}
	var b364 = sequenceBuilder{id: 364, commit: 258, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b360 = choiceBuilder{id: 360, commit: 258}
	var b341 = sequenceBuilder{id: 341, commit: 256, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{360}}
	var b338 = sequenceBuilder{id: 338, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b337 = charBuilder{}
	b338.items = []builder{&b337}
	var b340 = sequenceBuilder{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b339 = sequenceBuilder{id: 339, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b339.items = []builder{&b828, &b14}
	b340.items = []builder{&b828, &b14, &b339}
	b341.items = []builder{&b338, &b340, &b828, &b214}
	var b350 = sequenceBuilder{id: 350, commit: 256, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{360}}
	var b343 = sequenceBuilder{id: 343, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b342 = charBuilder{}
	b343.items = []builder{&b342}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b346.items = []builder{&b828, &b14}
	b347.items = []builder{&b828, &b14, &b346}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b348.items = []builder{&b828, &b14}
	b349.items = []builder{&b828, &b14, &b348}
	var b345 = sequenceBuilder{id: 345, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b344 = charBuilder{}
	b345.items = []builder{&b344}
	b350.items = []builder{&b343, &b347, &b828, &b502, &b349, &b828, &b345}
	var b359 = sequenceBuilder{id: 359, commit: 256, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{360}}
	var b352 = sequenceBuilder{id: 352, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b351 = charBuilder{}
	b352.items = []builder{&b351}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b355.items = []builder{&b828, &b14}
	b356.items = []builder{&b828, &b14, &b355}
	var b336 = sequenceBuilder{id: 336, commit: 258, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{611, 617, 618}}
	var b328 = sequenceBuilder{id: 328, commit: 256, name: "range-from", ranges: [][]int{{1, 1}}}
	b328.items = []builder{&b502}
	var b333 = sequenceBuilder{id: 333, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b332 = sequenceBuilder{id: 332, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b332.items = []builder{&b828, &b14}
	b333.items = []builder{&b828, &b14, &b332}
	var b331 = sequenceBuilder{id: 331, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b330 = charBuilder{}
	b331.items = []builder{&b330}
	var b335 = sequenceBuilder{id: 335, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b334 = sequenceBuilder{id: 334, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b334.items = []builder{&b828, &b14}
	b335.items = []builder{&b828, &b14, &b334}
	var b329 = sequenceBuilder{id: 329, commit: 256, name: "range-to", ranges: [][]int{{1, 1}}}
	b329.items = []builder{&b502}
	b336.items = []builder{&b328, &b333, &b828, &b331, &b335, &b828, &b329}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b357.items = []builder{&b828, &b14}
	b358.items = []builder{&b828, &b14, &b357}
	var b354 = sequenceBuilder{id: 354, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b353 = charBuilder{}
	b354.items = []builder{&b353}
	b359.items = []builder{&b352, &b356, &b828, &b336, &b358, &b828, &b354}
	b360.options = []builder{&b341, &b350, &b359}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b361.items = []builder{&b828, &b14}
	b362.items = []builder{&b14, &b361}
	b363.items = []builder{&b362, &b828, &b360}
	b364.items = []builder{&b360, &b828, &b363}
	b367.items = []builder{&b373, &b366, &b828, &b364}
	b630.options = []builder{&b214, &b367}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b633 = sequenceBuilder{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b633.items = []builder{&b828, &b14}
	b634.items = []builder{&b828, &b14, &b633}
	var b632 = sequenceBuilder{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b631 = charBuilder{}
	b632.items = []builder{&b631}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b635.items = []builder{&b828, &b14}
	b636.items = []builder{&b828, &b14, &b635}
	b637.items = []builder{&b630, &b634, &b828, &b632, &b636, &b828, &b502}
	var b302 = sequenceBuilder{id: 302, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{293, 314, 798}}
	var b295 = sequenceBuilder{id: 295, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b294 = charBuilder{}
	b295.items = []builder{&b294}
	var b299 = sequenceBuilder{id: 299, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b298 = sequenceBuilder{id: 298, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b298.items = []builder{&b828, &b14}
	b299.items = []builder{&b828, &b14, &b298}
	var b301 = sequenceBuilder{id: 301, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b300 = sequenceBuilder{id: 300, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b300.items = []builder{&b828, &b14}
	b301.items = []builder{&b828, &b14, &b300}
	var b297 = sequenceBuilder{id: 297, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b296 = charBuilder{}
	b297.items = []builder{&b296}
	b302.items = []builder{&b295, &b299, &b828, &b293, &b301, &b828, &b297}
	b293.options = []builder{&b574, &b606, &b609, &b637, &b302}
	var b307 = sequenceBuilder{id: 307, commit: 256, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{314}}
	var b304 = sequenceBuilder{id: 304, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b303 = charBuilder{}
	b304.items = []builder{&b303}
	var b306 = sequenceBuilder{id: 306, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b305 = charBuilder{}
	b306.items = []builder{&b305}
	b307.items = []builder{&b304, &b828, &b812, &b828, &b816, &b828, &b812, &b828, &b306}
	b314.options = []builder{&b502, &b293, &b307}
	b317.items = []builder{&b309, &b828, &b224, &b828, &b311, &b828, &b224, &b828, &b313, &b316, &b828, &b314}
	b320.items = []builder{&b70, &b319, &b828, &b317}
	var b327 = sequenceBuilder{id: 327, commit: 256, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b324 = sequenceBuilder{id: 324, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b323 = sequenceBuilder{id: 323, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b323.items = []builder{&b828, &b14}
	b324.items = []builder{&b828, &b14, &b323}
	var b322 = sequenceBuilder{id: 322, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b321 = charBuilder{}
	b322.items = []builder{&b321}
	var b326 = sequenceBuilder{id: 326, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b325 = sequenceBuilder{id: 325, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b325.items = []builder{&b828, &b14}
	b326.items = []builder{&b828, &b14, &b325}
	b327.items = []builder{&b70, &b324, &b828, &b322, &b326, &b828, &b317}
	var b577 = sequenceBuilder{id: 577, commit: 256, name: "receive-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 581, 618, 611}}
	var b111 = sequenceBuilder{id: 111, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b104 = charBuilder{}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	var b108 = charBuilder{}
	var b109 = charBuilder{}
	var b110 = charBuilder{}
	b111.items = []builder{&b104, &b105, &b106, &b107, &b108, &b109, &b110}
	var b576 = sequenceBuilder{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b575 = sequenceBuilder{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b575.items = []builder{&b828, &b14}
	b576.items = []builder{&b828, &b14, &b575}
	b577.items = []builder{&b111, &b576, &b828, &b373}
	var b511 = sequenceBuilder{id: 511, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{373, 502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b504 = sequenceBuilder{id: 504, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b503 = charBuilder{}
	b504.items = []builder{&b503}
	var b508 = sequenceBuilder{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b507 = sequenceBuilder{id: 507, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b507.items = []builder{&b828, &b14}
	b508.items = []builder{&b828, &b14, &b507}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b509.items = []builder{&b828, &b14}
	b510.items = []builder{&b828, &b14, &b509}
	var b506 = sequenceBuilder{id: 506, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b505 = charBuilder{}
	b506.items = []builder{&b505}
	b511.items = []builder{&b504, &b508, &b828, &b502, &b510, &b828, &b506}
	b373.options = []builder{&b182, &b195, &b208, &b209, &b214, &b235, &b240, &b269, &b274, &b320, &b327, &b367, &b372, &b577, &b511}
	var b433 = sequenceBuilder{id: 433, commit: 256, name: "unary", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{502, 439, 440, 441, 442, 443, 494, 618, 611}}
	var b432 = choiceBuilder{id: 432, commit: 258}
	var b392 = sequenceBuilder{id: 392, commit: 264, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var b391 = charBuilder{}
	b392.items = []builder{&b391}
	var b394 = sequenceBuilder{id: 394, commit: 264, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var b393 = charBuilder{}
	b394.items = []builder{&b393}
	var b375 = sequenceBuilder{id: 375, commit: 264, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var b374 = charBuilder{}
	b375.items = []builder{&b374}
	var b406 = sequenceBuilder{id: 406, commit: 264, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{432}}
	var b405 = charBuilder{}
	b406.items = []builder{&b405}
	b432.options = []builder{&b392, &b394, &b375, &b406}
	b433.items = []builder{&b432, &b828, &b373}
	var b480 = choiceBuilder{id: 480, commit: 258, generalizations: []int{502, 494, 618, 611}}
	var b451 = sequenceBuilder{id: 451, commit: 256, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 440, 441, 442, 443, 502, 494, 618, 611}}
	var b439 = choiceBuilder{id: 439, commit: 258, generalizations: []int{440, 441, 442, 443}}
	b439.options = []builder{&b373, &b433}
	var b449 = sequenceBuilder{id: 449, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b446 = sequenceBuilder{id: 446, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b445 = sequenceBuilder{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b445.items = []builder{&b828, &b14}
	b446.items = []builder{&b14, &b445}
	var b434 = choiceBuilder{id: 434, commit: 258}
	var b377 = sequenceBuilder{id: 377, commit: 264, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var b376 = charBuilder{}
	b377.items = []builder{&b376}
	var b384 = sequenceBuilder{id: 384, commit: 264, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{434}}
	var b382 = charBuilder{}
	var b383 = charBuilder{}
	b384.items = []builder{&b382, &b383}
	var b387 = sequenceBuilder{id: 387, commit: 264, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{434}}
	var b385 = charBuilder{}
	var b386 = charBuilder{}
	b387.items = []builder{&b385, &b386}
	var b390 = sequenceBuilder{id: 390, commit: 264, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{434}}
	var b388 = charBuilder{}
	var b389 = charBuilder{}
	b390.items = []builder{&b388, &b389}
	var b396 = sequenceBuilder{id: 396, commit: 264, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var b395 = charBuilder{}
	b396.items = []builder{&b395}
	var b398 = sequenceBuilder{id: 398, commit: 264, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var b397 = charBuilder{}
	b398.items = []builder{&b397}
	var b400 = sequenceBuilder{id: 400, commit: 264, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{434}}
	var b399 = charBuilder{}
	b400.items = []builder{&b399}
	b434.options = []builder{&b377, &b384, &b387, &b390, &b396, &b398, &b400}
	var b448 = sequenceBuilder{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b447 = sequenceBuilder{id: 447, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b447.items = []builder{&b828, &b14}
	b448.items = []builder{&b828, &b14, &b447}
	b449.items = []builder{&b446, &b828, &b434, &b448, &b828, &b439}
	var b450 = sequenceBuilder{id: 450, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b450.items = []builder{&b828, &b449}
	b451.items = []builder{&b439, &b828, &b449, &b450}
	var b458 = sequenceBuilder{id: 458, commit: 256, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 441, 442, 443, 502, 494, 618, 611}}
	var b440 = choiceBuilder{id: 440, commit: 258, generalizations: []int{441, 442, 443}}
	b440.options = []builder{&b439, &b451}
	var b456 = sequenceBuilder{id: 456, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b453 = sequenceBuilder{id: 453, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b452 = sequenceBuilder{id: 452, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b452.items = []builder{&b828, &b14}
	b453.items = []builder{&b14, &b452}
	var b435 = choiceBuilder{id: 435, commit: 258}
	var b379 = sequenceBuilder{id: 379, commit: 264, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var b378 = charBuilder{}
	b379.items = []builder{&b378}
	var b381 = sequenceBuilder{id: 381, commit: 264, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var b380 = charBuilder{}
	b381.items = []builder{&b380}
	var b402 = sequenceBuilder{id: 402, commit: 264, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var b401 = charBuilder{}
	b402.items = []builder{&b401}
	var b404 = sequenceBuilder{id: 404, commit: 264, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{435}}
	var b403 = charBuilder{}
	b404.items = []builder{&b403}
	b435.options = []builder{&b379, &b381, &b402, &b404}
	var b455 = sequenceBuilder{id: 455, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b454 = sequenceBuilder{id: 454, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b454.items = []builder{&b828, &b14}
	b455.items = []builder{&b828, &b14, &b454}
	b456.items = []builder{&b453, &b828, &b435, &b455, &b828, &b440}
	var b457 = sequenceBuilder{id: 457, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b457.items = []builder{&b828, &b456}
	b458.items = []builder{&b440, &b828, &b456, &b457}
	var b465 = sequenceBuilder{id: 465, commit: 256, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 442, 443, 502, 494, 618, 611}}
	var b441 = choiceBuilder{id: 441, commit: 258, generalizations: []int{442, 443}}
	b441.options = []builder{&b440, &b458}
	var b463 = sequenceBuilder{id: 463, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b459.items = []builder{&b828, &b14}
	b460.items = []builder{&b14, &b459}
	var b436 = choiceBuilder{id: 436, commit: 258}
	var b409 = sequenceBuilder{id: 409, commit: 264, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var b407 = charBuilder{}
	var b408 = charBuilder{}
	b409.items = []builder{&b407, &b408}
	var b412 = sequenceBuilder{id: 412, commit: 264, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var b410 = charBuilder{}
	var b411 = charBuilder{}
	b412.items = []builder{&b410, &b411}
	var b414 = sequenceBuilder{id: 414, commit: 264, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{436}}
	var b413 = charBuilder{}
	b414.items = []builder{&b413}
	var b417 = sequenceBuilder{id: 417, commit: 264, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var b415 = charBuilder{}
	var b416 = charBuilder{}
	b417.items = []builder{&b415, &b416}
	var b419 = sequenceBuilder{id: 419, commit: 264, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{436}}
	var b418 = charBuilder{}
	b419.items = []builder{&b418}
	var b422 = sequenceBuilder{id: 422, commit: 264, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{436}}
	var b420 = charBuilder{}
	var b421 = charBuilder{}
	b422.items = []builder{&b420, &b421}
	b436.options = []builder{&b409, &b412, &b414, &b417, &b419, &b422}
	var b462 = sequenceBuilder{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b461 = sequenceBuilder{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b461.items = []builder{&b828, &b14}
	b462.items = []builder{&b828, &b14, &b461}
	b463.items = []builder{&b460, &b828, &b436, &b462, &b828, &b441}
	var b464 = sequenceBuilder{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b464.items = []builder{&b828, &b463}
	b465.items = []builder{&b441, &b828, &b463, &b464}
	var b472 = sequenceBuilder{id: 472, commit: 256, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 443, 502, 494, 618, 611}}
	var b442 = choiceBuilder{id: 442, commit: 258, generalizations: []int{443}}
	b442.options = []builder{&b441, &b465}
	var b470 = sequenceBuilder{id: 470, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b467 = sequenceBuilder{id: 467, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b466 = sequenceBuilder{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b466.items = []builder{&b828, &b14}
	b467.items = []builder{&b14, &b466}
	var b437 = sequenceBuilder{id: 437, commit: 258, ranges: [][]int{{1, 1}}}
	var b425 = sequenceBuilder{id: 425, commit: 264, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b423 = charBuilder{}
	var b424 = charBuilder{}
	b425.items = []builder{&b423, &b424}
	b437.items = []builder{&b425}
	var b469 = sequenceBuilder{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b468 = sequenceBuilder{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b468.items = []builder{&b828, &b14}
	b469.items = []builder{&b828, &b14, &b468}
	b470.items = []builder{&b467, &b828, &b437, &b469, &b828, &b442}
	var b471 = sequenceBuilder{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b471.items = []builder{&b828, &b470}
	b472.items = []builder{&b442, &b828, &b470, &b471}
	var b479 = sequenceBuilder{id: 479, commit: 256, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{480, 502, 494, 618, 611}}
	var b443 = choiceBuilder{id: 443, commit: 258}
	b443.options = []builder{&b442, &b472}
	var b477 = sequenceBuilder{id: 477, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b474 = sequenceBuilder{id: 474, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b473 = sequenceBuilder{id: 473, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b473.items = []builder{&b828, &b14}
	b474.items = []builder{&b14, &b473}
	var b438 = sequenceBuilder{id: 438, commit: 258, ranges: [][]int{{1, 1}}}
	var b428 = sequenceBuilder{id: 428, commit: 264, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b426 = charBuilder{}
	var b427 = charBuilder{}
	b428.items = []builder{&b426, &b427}
	b438.items = []builder{&b428}
	var b476 = sequenceBuilder{id: 476, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b475 = sequenceBuilder{id: 475, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b475.items = []builder{&b828, &b14}
	b476.items = []builder{&b828, &b14, &b475}
	b477.items = []builder{&b474, &b828, &b438, &b476, &b828, &b443}
	var b478 = sequenceBuilder{id: 478, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b478.items = []builder{&b828, &b477}
	b479.items = []builder{&b443, &b828, &b477, &b478}
	b480.options = []builder{&b451, &b458, &b465, &b472, &b479}
	var b493 = sequenceBuilder{id: 493, commit: 256, name: "ternary", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{502, 494, 618, 611}}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b485.items = []builder{&b828, &b14}
	b486.items = []builder{&b828, &b14, &b485}
	var b482 = sequenceBuilder{id: 482, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b481 = charBuilder{}
	b482.items = []builder{&b481}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b487.items = []builder{&b828, &b14}
	b488.items = []builder{&b828, &b14, &b487}
	var b490 = sequenceBuilder{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b489.items = []builder{&b828, &b14}
	b490.items = []builder{&b828, &b14, &b489}
	var b484 = sequenceBuilder{id: 484, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b483 = charBuilder{}
	b484.items = []builder{&b483}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b491 = sequenceBuilder{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b491.items = []builder{&b828, &b14}
	b492.items = []builder{&b828, &b14, &b491}
	b493.items = []builder{&b502, &b486, &b828, &b482, &b488, &b828, &b502, &b490, &b828, &b484, &b492, &b828, &b502}
	var b501 = sequenceBuilder{id: 501, commit: 256, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{502, 798, 618, 611}}
	var b494 = choiceBuilder{id: 494, commit: 258}
	b494.options = []builder{&b373, &b433, &b480, &b493}
	var b499 = sequenceBuilder{id: 499, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b496 = sequenceBuilder{id: 496, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b495 = sequenceBuilder{id: 495, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b495.items = []builder{&b828, &b14}
	b496.items = []builder{&b14, &b495}
	var b431 = sequenceBuilder{id: 431, commit: 266, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b429 = charBuilder{}
	var b430 = charBuilder{}
	b431.items = []builder{&b429, &b430}
	var b498 = sequenceBuilder{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b497 = sequenceBuilder{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b497.items = []builder{&b828, &b14}
	b498.items = []builder{&b828, &b14, &b497}
	b499.items = []builder{&b496, &b828, &b431, &b498, &b828, &b494}
	var b500 = sequenceBuilder{id: 500, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b500.items = []builder{&b828, &b499}
	b501.items = []builder{&b494, &b828, &b499, &b500}
	b502.options = []builder{&b373, &b433, &b480, &b493, &b501}
	b288.items = []builder{&b61, &b287, &b828, &b502}
	b289.options = []builder{&b61, &b288}
	var b292 = sequenceBuilder{id: 292, commit: 256, name: "check-ret", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var b67 = sequenceBuilder{id: 67, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b62 = charBuilder{}
	var b63 = charBuilder{}
	var b64 = charBuilder{}
	var b65 = charBuilder{}
	var b66 = charBuilder{}
	b67.items = []builder{&b62, &b63, &b64, &b65, &b66}
	var b291 = sequenceBuilder{id: 291, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b290 = sequenceBuilder{id: 290, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b290.items = []builder{&b828, &b14}
	b291.items = []builder{&b828, &b14, &b290}
	b292.items = []builder{&b67, &b291, &b828, &b502}
	var b532 = sequenceBuilder{id: 532, commit: 256, name: "if-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{798}}
	var b73 = sequenceBuilder{id: 73, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b71 = charBuilder{}
	var b72 = charBuilder{}
	b73.items = []builder{&b71, &b72}
	var b527 = sequenceBuilder{id: 527, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b526 = sequenceBuilder{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b526.items = []builder{&b828, &b14}
	b527.items = []builder{&b828, &b14, &b526}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b528 = sequenceBuilder{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b528.items = []builder{&b828, &b14}
	b529.items = []builder{&b828, &b14, &b528}
	var b531 = sequenceBuilder{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b520 = sequenceBuilder{id: 520, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b513 = sequenceBuilder{id: 513, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b512 = sequenceBuilder{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b512.items = []builder{&b828, &b14}
	b513.items = []builder{&b14, &b512}
	var b78 = sequenceBuilder{id: 78, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b74 = charBuilder{}
	var b75 = charBuilder{}
	var b76 = charBuilder{}
	var b77 = charBuilder{}
	b78.items = []builder{&b74, &b75, &b76, &b77}
	var b515 = sequenceBuilder{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b514 = sequenceBuilder{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b514.items = []builder{&b828, &b14}
	b515.items = []builder{&b828, &b14, &b514}
	var b517 = sequenceBuilder{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b516 = sequenceBuilder{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b516.items = []builder{&b828, &b14}
	b517.items = []builder{&b828, &b14, &b516}
	var b519 = sequenceBuilder{id: 519, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b518 = sequenceBuilder{id: 518, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b518.items = []builder{&b828, &b14}
	b519.items = []builder{&b828, &b14, &b518}
	b520.items = []builder{&b513, &b828, &b78, &b515, &b828, &b73, &b517, &b828, &b502, &b519, &b828, &b307}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b530.items = []builder{&b828, &b520}
	b531.items = []builder{&b828, &b520, &b530}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b522 = sequenceBuilder{id: 522, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b521 = sequenceBuilder{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b521.items = []builder{&b828, &b14}
	b522.items = []builder{&b14, &b521}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b523.items = []builder{&b828, &b14}
	b524.items = []builder{&b828, &b14, &b523}
	b525.items = []builder{&b522, &b828, &b78, &b524, &b828, &b307}
	b532.items = []builder{&b73, &b527, &b828, &b502, &b529, &b828, &b307, &b531, &b828, &b525}
	var b569 = choiceBuilder{id: 569, commit: 256, name: "switch-statement", generalizations: []int{798}}
	var b559 = sequenceBuilder{id: 559, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{569, 798}}
	var b85 = sequenceBuilder{id: 85, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b79 = charBuilder{}
	var b80 = charBuilder{}
	var b81 = charBuilder{}
	var b82 = charBuilder{}
	var b83 = charBuilder{}
	var b84 = charBuilder{}
	b85.items = []builder{&b79, &b80, &b81, &b82, &b83, &b84}
	var b558 = sequenceBuilder{id: 558, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b557 = sequenceBuilder{id: 557, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b557.items = []builder{&b828, &b14}
	b558.items = []builder{&b828, &b14, &b557}
	var b554 = sequenceBuilder{id: 554, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b553 = charBuilder{}
	b554.items = []builder{&b553}
	var b552 = choiceBuilder{id: 552, commit: 258}
	var b548 = sequenceBuilder{id: 548, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{552}}
	var b539 = sequenceBuilder{id: 539, commit: 256, name: "case-block", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b90 = sequenceBuilder{id: 90, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b86 = charBuilder{}
	var b87 = charBuilder{}
	var b88 = charBuilder{}
	var b89 = charBuilder{}
	b90.items = []builder{&b86, &b87, &b88, &b89}
	var b536 = sequenceBuilder{id: 536, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b535 = sequenceBuilder{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b535.items = []builder{&b828, &b14}
	b536.items = []builder{&b828, &b14, &b535}
	var b538 = sequenceBuilder{id: 538, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b537 = sequenceBuilder{id: 537, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b537.items = []builder{&b828, &b14}
	b538.items = []builder{&b828, &b14, &b537}
	var b534 = sequenceBuilder{id: 534, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b533 = charBuilder{}
	b534.items = []builder{&b533}
	b539.items = []builder{&b90, &b536, &b828, &b502, &b538, &b828, &b534, &b828, &b812, &b828, &b816}
	var b547 = sequenceBuilder{id: 547, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b545.items = []builder{&b812, &b828, &b539}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b546.items = []builder{&b828, &b545}
	b547.items = []builder{&b828, &b545, &b546}
	b548.items = []builder{&b539, &b547}
	var b544 = sequenceBuilder{id: 544, commit: 256, name: "default-block", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{552, 596}}
	var b98 = sequenceBuilder{id: 98, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b91 = charBuilder{}
	var b92 = charBuilder{}
	var b93 = charBuilder{}
	var b94 = charBuilder{}
	var b95 = charBuilder{}
	var b96 = charBuilder{}
	var b97 = charBuilder{}
	b98.items = []builder{&b91, &b92, &b93, &b94, &b95, &b96, &b97}
	var b543 = sequenceBuilder{id: 543, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b542 = sequenceBuilder{id: 542, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b542.items = []builder{&b828, &b14}
	b543.items = []builder{&b828, &b14, &b542}
	var b541 = sequenceBuilder{id: 541, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b540 = charBuilder{}
	b541.items = []builder{&b540}
	b544.items = []builder{&b98, &b543, &b828, &b541, &b828, &b812, &b828, &b816}
	var b551 = sequenceBuilder{id: 551, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{552}}
	var b549 = sequenceBuilder{id: 549, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b549.items = []builder{&b548, &b828, &b812}
	var b550 = sequenceBuilder{id: 550, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b550.items = []builder{&b812, &b828, &b548}
	b551.items = []builder{&b549, &b828, &b544, &b828, &b550}
	b552.options = []builder{&b548, &b544, &b551}
	var b556 = sequenceBuilder{id: 556, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b555 = charBuilder{}
	b556.items = []builder{&b555}
	b559.items = []builder{&b85, &b558, &b828, &b554, &b828, &b812, &b828, &b552, &b828, &b812, &b828, &b556}
	var b568 = sequenceBuilder{id: 568, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{569, 798}}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b564.items = []builder{&b828, &b14}
	b565.items = []builder{&b828, &b14, &b564}
	var b567 = sequenceBuilder{id: 567, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b566 = sequenceBuilder{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b566.items = []builder{&b828, &b14}
	b567.items = []builder{&b828, &b14, &b566}
	var b561 = sequenceBuilder{id: 561, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b560 = charBuilder{}
	b561.items = []builder{&b560}
	var b563 = sequenceBuilder{id: 563, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b562 = charBuilder{}
	b563.items = []builder{&b562}
	b568.items = []builder{&b85, &b565, &b828, &b502, &b567, &b828, &b561, &b828, &b812, &b828, &b552, &b828, &b812, &b828, &b563}
	b569.options = []builder{&b559, &b568}
	var b603 = sequenceBuilder{id: 603, commit: 256, name: "select-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var b118 = sequenceBuilder{id: 118, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b112 = charBuilder{}
	var b113 = charBuilder{}
	var b114 = charBuilder{}
	var b115 = charBuilder{}
	var b116 = charBuilder{}
	var b117 = charBuilder{}
	b118.items = []builder{&b112, &b113, &b114, &b115, &b116, &b117}
	var b602 = sequenceBuilder{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b601 = sequenceBuilder{id: 601, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b601.items = []builder{&b828, &b14}
	b602.items = []builder{&b828, &b14, &b601}
	var b598 = sequenceBuilder{id: 598, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b597 = charBuilder{}
	b598.items = []builder{&b597}
	var b596 = choiceBuilder{id: 596, commit: 258}
	var b592 = sequenceBuilder{id: 592, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{596}}
	var b588 = sequenceBuilder{id: 588, commit: 256, name: "select-case-block", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b584.items = []builder{&b828, &b14}
	b585.items = []builder{&b828, &b14, &b584}
	var b581 = choiceBuilder{id: 581, commit: 258}
	var b580 = sequenceBuilder{id: 580, commit: 256, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{581}}
	var b579 = sequenceBuilder{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b578 = sequenceBuilder{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b578.items = []builder{&b828, &b14}
	b579.items = []builder{&b828, &b14, &b578}
	b580.items = []builder{&b214, &b579, &b828, &b577}
	b581.options = []builder{&b574, &b577, &b580}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b586.items = []builder{&b828, &b14}
	b587.items = []builder{&b828, &b14, &b586}
	var b583 = sequenceBuilder{id: 583, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b582 = charBuilder{}
	b583.items = []builder{&b582}
	b588.items = []builder{&b90, &b585, &b828, &b581, &b587, &b828, &b583, &b828, &b812, &b828, &b816}
	var b591 = sequenceBuilder{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b589.items = []builder{&b812, &b828, &b588}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b590.items = []builder{&b828, &b589}
	b591.items = []builder{&b828, &b589, &b590}
	b592.items = []builder{&b588, &b591}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{596}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b593.items = []builder{&b592, &b828, &b812}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b594.items = []builder{&b812, &b828, &b592}
	b595.items = []builder{&b593, &b828, &b544, &b828, &b594}
	b596.options = []builder{&b592, &b544, &b595}
	var b600 = sequenceBuilder{id: 600, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b599 = charBuilder{}
	b600.items = []builder{&b599}
	b603.items = []builder{&b118, &b602, &b828, &b598, &b828, &b812, &b828, &b596, &b828, &b812, &b828, &b600}
	var b610 = choiceBuilder{id: 610, commit: 258, generalizations: []int{798}}
	var b140 = sequenceBuilder{id: 140, commit: 280, name: "break", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{610, 798}}
	var b135 = charBuilder{}
	var b136 = charBuilder{}
	var b137 = charBuilder{}
	var b138 = charBuilder{}
	var b139 = charBuilder{}
	b140.items = []builder{&b135, &b136, &b137, &b138, &b139}
	var b149 = sequenceBuilder{id: 149, commit: 280, name: "continue", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{610, 798}}
	var b141 = charBuilder{}
	var b142 = charBuilder{}
	var b143 = charBuilder{}
	var b144 = charBuilder{}
	var b145 = charBuilder{}
	var b146 = charBuilder{}
	var b147 = charBuilder{}
	var b148 = charBuilder{}
	b149.items = []builder{&b141, &b142, &b143, &b144, &b145, &b146, &b147, &b148}
	b610.options = []builder{&b140, &b149}
	var b629 = sequenceBuilder{id: 629, commit: 256, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var b134 = sequenceBuilder{id: 134, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b131 = charBuilder{}
	var b132 = charBuilder{}
	var b133 = charBuilder{}
	b134.items = []builder{&b131, &b132, &b133}
	var b628 = choiceBuilder{id: 628, commit: 2}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{628}}
	var b621 = sequenceBuilder{id: 621, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b620 = sequenceBuilder{id: 620, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b619 = sequenceBuilder{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b619.items = []builder{&b828, &b14}
	b620.items = []builder{&b14, &b619}
	var b618 = choiceBuilder{id: 618, commit: 258}
	var b617 = choiceBuilder{id: 617, commit: 256, name: "range-over", generalizations: []int{618}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{617, 618}}
	var b613 = sequenceBuilder{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b612 = sequenceBuilder{id: 612, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b612.items = []builder{&b828, &b14}
	b613.items = []builder{&b828, &b14, &b612}
	var b130 = sequenceBuilder{id: 130, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b128 = charBuilder{}
	var b129 = charBuilder{}
	b130.items = []builder{&b128, &b129}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b614.items = []builder{&b828, &b14}
	b615.items = []builder{&b828, &b14, &b614}
	var b611 = choiceBuilder{id: 611, commit: 2}
	b611.options = []builder{&b502, &b336}
	b616.items = []builder{&b214, &b613, &b828, &b130, &b615, &b828, &b611}
	b617.options = []builder{&b616, &b336}
	b618.options = []builder{&b502, &b617}
	b621.items = []builder{&b620, &b828, &b618}
	var b623 = sequenceBuilder{id: 623, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b622 = sequenceBuilder{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b622.items = []builder{&b828, &b14}
	b623.items = []builder{&b828, &b14, &b622}
	b624.items = []builder{&b621, &b623, &b828, &b307}
	var b627 = sequenceBuilder{id: 627, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{628}}
	var b626 = sequenceBuilder{id: 626, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b625 = sequenceBuilder{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b625.items = []builder{&b828, &b14}
	b626.items = []builder{&b14, &b625}
	b627.items = []builder{&b626, &b828, &b307}
	b628.options = []builder{&b624, &b627}
	b629.items = []builder{&b134, &b828, &b628}
	var b730 = choiceBuilder{id: 730, commit: 258, generalizations: []int{798}}
	var b657 = sequenceBuilder{id: 657, commit: 256, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var b653 = sequenceBuilder{id: 653, commit: 266, ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	var b638 = sequenceBuilder{id: 638, commit: 256, name: "docs", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	b638.items = []builder{&b43, &b828, &b14}
	var b153 = sequenceBuilder{id: 153, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b150 = charBuilder{}
	var b151 = charBuilder{}
	var b152 = charBuilder{}
	b153.items = []builder{&b150, &b151, &b152}
	b653.items = []builder{&b638, &b153}
	var b656 = sequenceBuilder{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b655 = sequenceBuilder{id: 655, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b655.items = []builder{&b828, &b14}
	b656.items = []builder{&b828, &b14, &b655}
	var b654 = choiceBuilder{id: 654, commit: 2}
	var b647 = sequenceBuilder{id: 647, commit: 256, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{654, 659}}
	var b646 = sequenceBuilder{id: 646, commit: 258, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b643 = sequenceBuilder{id: 643, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b641 = sequenceBuilder{id: 641, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b641.items = []builder{&b828, &b14}
	b642.items = []builder{&b14, &b641}
	var b640 = sequenceBuilder{id: 640, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b639 = charBuilder{}
	b640.items = []builder{&b639}
	b643.items = []builder{&b642, &b828, &b640}
	var b645 = sequenceBuilder{id: 645, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b644 = sequenceBuilder{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b644.items = []builder{&b828, &b14}
	b645.items = []builder{&b828, &b14, &b644}
	b646.items = []builder{&b214, &b828, &b643, &b645, &b828, &b502}
	b647.items = []builder{&b646}
	var b652 = sequenceBuilder{id: 652, commit: 256, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{654, 659}}
	var b649 = sequenceBuilder{id: 649, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b648 = charBuilder{}
	b649.items = []builder{&b648}
	var b651 = sequenceBuilder{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b650 = sequenceBuilder{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b650.items = []builder{&b828, &b14}
	b651.items = []builder{&b828, &b14, &b650}
	b652.items = []builder{&b649, &b651, &b828, &b646}
	b654.options = []builder{&b647, &b652}
	b657.items = []builder{&b653, &b656, &b828, &b654}
	var b675 = sequenceBuilder{id: 675, commit: 256, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var b674 = sequenceBuilder{id: 674, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b673 = sequenceBuilder{id: 673, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b673.items = []builder{&b828, &b14}
	b674.items = []builder{&b828, &b14, &b673}
	var b670 = sequenceBuilder{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b669 = charBuilder{}
	b670.items = []builder{&b669}
	var b668 = sequenceBuilder{id: 668, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b660 = sequenceBuilder{id: 660, commit: 264, name: "docs-mixed-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	var b659 = choiceBuilder{id: 659, commit: 10}
	b659.options = []builder{&b647, &b652}
	b660.items = []builder{&b638, &b659}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b665 = sequenceBuilder{id: 665, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b665.items = []builder{&b224, &b828, &b660}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b666.items = []builder{&b828, &b665}
	b667.items = []builder{&b828, &b665, &b666}
	b668.items = []builder{&b660, &b667}
	var b672 = sequenceBuilder{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b671 = charBuilder{}
	b672.items = []builder{&b671}
	b675.items = []builder{&b153, &b674, &b828, &b670, &b828, &b224, &b828, &b668, &b828, &b224, &b828, &b672}
	var b686 = sequenceBuilder{id: 686, commit: 256, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var b683 = sequenceBuilder{id: 683, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b682 = sequenceBuilder{id: 682, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b682.items = []builder{&b828, &b14}
	b683.items = []builder{&b828, &b14, &b682}
	var b677 = sequenceBuilder{id: 677, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b676 = charBuilder{}
	b677.items = []builder{&b676}
	var b685 = sequenceBuilder{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b684 = sequenceBuilder{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b684.items = []builder{&b828, &b14}
	b685.items = []builder{&b828, &b14, &b684}
	var b679 = sequenceBuilder{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b678 = charBuilder{}
	b679.items = []builder{&b678}
	var b664 = sequenceBuilder{id: 664, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b658 = sequenceBuilder{id: 658, commit: 264, name: "docs-value-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	b658.items = []builder{&b638, &b647}
	var b663 = sequenceBuilder{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b661 = sequenceBuilder{id: 661, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b661.items = []builder{&b224, &b828, &b658}
	var b662 = sequenceBuilder{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b662.items = []builder{&b828, &b661}
	b663.items = []builder{&b828, &b661, &b662}
	b664.items = []builder{&b658, &b663}
	var b681 = sequenceBuilder{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	b681.items = []builder{&b680}
	b686.items = []builder{&b153, &b683, &b828, &b677, &b685, &b828, &b679, &b828, &b224, &b828, &b664, &b828, &b224, &b828, &b681}
	var b700 = sequenceBuilder{id: 700, commit: 256, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var b696 = sequenceBuilder{id: 696, commit: 266, ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	b696.items = []builder{&b638, &b70}
	var b699 = sequenceBuilder{id: 699, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b698 = sequenceBuilder{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b698.items = []builder{&b828, &b14}
	b699.items = []builder{&b828, &b14, &b698}
	var b697 = choiceBuilder{id: 697, commit: 2}
	var b690 = sequenceBuilder{id: 690, commit: 256, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{697, 702}}
	var b689 = sequenceBuilder{id: 689, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b687.items = []builder{&b828, &b14}
	b688.items = []builder{&b828, &b14, &b687}
	b689.items = []builder{&b214, &b688, &b828, &b317}
	b690.items = []builder{&b689}
	var b695 = sequenceBuilder{id: 695, commit: 256, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{697, 702}}
	var b692 = sequenceBuilder{id: 692, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b691 = charBuilder{}
	b692.items = []builder{&b691}
	var b694 = sequenceBuilder{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b693 = sequenceBuilder{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b693.items = []builder{&b828, &b14}
	b694.items = []builder{&b828, &b14, &b693}
	b695.items = []builder{&b692, &b694, &b828, &b689}
	b697.options = []builder{&b690, &b695}
	b700.items = []builder{&b696, &b699, &b828, &b697}
	var b718 = sequenceBuilder{id: 718, commit: 256, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var b717 = sequenceBuilder{id: 717, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b716 = sequenceBuilder{id: 716, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b716.items = []builder{&b828, &b14}
	b717.items = []builder{&b828, &b14, &b716}
	var b713 = sequenceBuilder{id: 713, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b712 = charBuilder{}
	b713.items = []builder{&b712}
	var b711 = sequenceBuilder{id: 711, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b703 = sequenceBuilder{id: 703, commit: 264, name: "docs-mixed-function-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	var b702 = choiceBuilder{id: 702, commit: 10}
	b702.options = []builder{&b690, &b695}
	b703.items = []builder{&b638, &b702}
	var b710 = sequenceBuilder{id: 710, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b708 = sequenceBuilder{id: 708, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b708.items = []builder{&b224, &b828, &b703}
	var b709 = sequenceBuilder{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b709.items = []builder{&b828, &b708}
	b710.items = []builder{&b828, &b708, &b709}
	b711.items = []builder{&b703, &b710}
	var b715 = sequenceBuilder{id: 715, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b714 = charBuilder{}
	b715.items = []builder{&b714}
	b718.items = []builder{&b70, &b717, &b828, &b713, &b828, &b224, &b828, &b711, &b828, &b224, &b828, &b715}
	var b729 = sequenceBuilder{id: 729, commit: 256, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 798}}
	var b726 = sequenceBuilder{id: 726, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b725 = sequenceBuilder{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b725.items = []builder{&b828, &b14}
	b726.items = []builder{&b828, &b14, &b725}
	var b720 = sequenceBuilder{id: 720, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b719 = charBuilder{}
	b720.items = []builder{&b719}
	var b728 = sequenceBuilder{id: 728, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b727 = sequenceBuilder{id: 727, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b727.items = []builder{&b828, &b14}
	b728.items = []builder{&b828, &b14, &b727}
	var b722 = sequenceBuilder{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b721 = charBuilder{}
	b722.items = []builder{&b721}
	var b707 = sequenceBuilder{id: 707, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b701 = sequenceBuilder{id: 701, commit: 264, name: "docs-function-capture", ranges: [][]int{{0, 1}, {1, 1}, {0, 1}, {1, 1}}}
	b701.items = []builder{&b638, &b690}
	var b706 = sequenceBuilder{id: 706, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b704.items = []builder{&b224, &b828, &b701}
	var b705 = sequenceBuilder{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b705.items = []builder{&b828, &b704}
	b706.items = []builder{&b828, &b704, &b705}
	b707.items = []builder{&b701, &b706}
	var b724 = sequenceBuilder{id: 724, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b723 = charBuilder{}
	b724.items = []builder{&b723}
	b729.items = []builder{&b70, &b726, &b828, &b720, &b728, &b828, &b722, &b828, &b224, &b828, &b707, &b828, &b224, &b828, &b724}
	b730.options = []builder{&b657, &b675, &b686, &b700, &b718, &b729}
	var b797 = sequenceBuilder{id: 797, commit: 256, name: "export-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var b160 = sequenceBuilder{id: 160, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b154 = charBuilder{}
	var b155 = charBuilder{}
	var b156 = charBuilder{}
	var b157 = charBuilder{}
	var b158 = charBuilder{}
	var b159 = charBuilder{}
	b160.items = []builder{&b154, &b155, &b156, &b157, &b158, &b159}
	var b796 = sequenceBuilder{id: 796, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b795 = sequenceBuilder{id: 795, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b795.items = []builder{&b828, &b14}
	b796.items = []builder{&b828, &b14, &b795}
	b797.items = []builder{&b160, &b796, &b828, &b730}
	var b794 = choiceBuilder{id: 794, commit: 256, name: "use-modules", generalizations: []int{798}}
	var b775 = sequenceBuilder{id: 775, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 798}}
	var b164 = sequenceBuilder{id: 164, commit: 282, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b161 = charBuilder{}
	var b162 = charBuilder{}
	var b163 = charBuilder{}
	b164.items = []builder{&b161, &b162, &b163}
	var b774 = sequenceBuilder{id: 774, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b773 = sequenceBuilder{id: 773, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b773.items = []builder{&b828, &b14}
	b774.items = []builder{&b828, &b14, &b773}
	var b772 = choiceBuilder{id: 772, commit: 2}
	var b742 = choiceBuilder{id: 742, commit: 256, name: "use-fact", generalizations: []int{772, 766, 767}}
	var b741 = sequenceBuilder{id: 741, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{742, 772, 766, 767}}
	var b733 = choiceBuilder{id: 733, commit: 2}
	var b732 = sequenceBuilder{id: 732, commit: 264, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{733, 748}}
	var b731 = charBuilder{}
	b732.items = []builder{&b731}
	b733.options = []builder{&b214, &b732}
	var b738 = sequenceBuilder{id: 738, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b736.items = []builder{&b828, &b14}
	b737.items = []builder{&b14, &b736}
	var b735 = sequenceBuilder{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b734 = charBuilder{}
	b735.items = []builder{&b734}
	b738.items = []builder{&b737, &b828, &b735}
	var b740 = sequenceBuilder{id: 740, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b739 = sequenceBuilder{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b739.items = []builder{&b828, &b14}
	b740.items = []builder{&b828, &b14, &b739}
	b741.items = []builder{&b733, &b828, &b738, &b740, &b828, &b208}
	b742.options = []builder{&b208, &b741}
	var b761 = choiceBuilder{id: 761, commit: 256, name: "use-effect", generalizations: []int{772, 766, 767}}
	var b747 = sequenceBuilder{id: 747, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{761, 772, 766, 767}}
	var b744 = sequenceBuilder{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b743 = charBuilder{}
	b744.items = []builder{&b743}
	var b746 = sequenceBuilder{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b745 = sequenceBuilder{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b745.items = []builder{&b828, &b14}
	b746.items = []builder{&b828, &b14, &b745}
	b747.items = []builder{&b744, &b746, &b828, &b208}
	var b760 = sequenceBuilder{id: 760, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{761, 772, 766, 767}}
	var b748 = choiceBuilder{id: 748, commit: 2}
	b748.options = []builder{&b214, &b732}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b751 = sequenceBuilder{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b751.items = []builder{&b828, &b14}
	b752.items = []builder{&b14, &b751}
	var b750 = sequenceBuilder{id: 750, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b749 = charBuilder{}
	b750.items = []builder{&b749}
	b753.items = []builder{&b752, &b828, &b750}
	var b757 = sequenceBuilder{id: 757, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b756.items = []builder{&b828, &b14}
	b757.items = []builder{&b828, &b14, &b756}
	var b755 = sequenceBuilder{id: 755, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b754 = charBuilder{}
	b755.items = []builder{&b754}
	var b759 = sequenceBuilder{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b758 = sequenceBuilder{id: 758, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b758.items = []builder{&b828, &b14}
	b759.items = []builder{&b828, &b14, &b758}
	b760.items = []builder{&b748, &b828, &b753, &b757, &b828, &b755, &b759, &b828, &b208}
	b761.options = []builder{&b747, &b760}
	b772.options = []builder{&b742, &b761}
	b775.items = []builder{&b164, &b774, &b828, &b772}
	var b782 = sequenceBuilder{id: 782, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 798}}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b780 = sequenceBuilder{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b780.items = []builder{&b828, &b14}
	b781.items = []builder{&b828, &b14, &b780}
	var b777 = sequenceBuilder{id: 777, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b776 = charBuilder{}
	b777.items = []builder{&b776}
	var b771 = sequenceBuilder{id: 771, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b766 = choiceBuilder{id: 766, commit: 2}
	b766.options = []builder{&b742, &b761}
	var b770 = sequenceBuilder{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b768 = sequenceBuilder{id: 768, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b767 = choiceBuilder{id: 767, commit: 2}
	b767.options = []builder{&b742, &b761}
	b768.items = []builder{&b224, &b828, &b767}
	var b769 = sequenceBuilder{id: 769, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b769.items = []builder{&b828, &b768}
	b770.items = []builder{&b828, &b768, &b769}
	b771.items = []builder{&b766, &b770}
	var b779 = sequenceBuilder{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b778 = charBuilder{}
	b779.items = []builder{&b778}
	b782.items = []builder{&b164, &b781, &b828, &b777, &b828, &b224, &b828, &b771, &b828, &b224, &b828, &b779}
	var b793 = sequenceBuilder{id: 793, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 798}}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b789 = sequenceBuilder{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b789.items = []builder{&b828, &b14}
	b790.items = []builder{&b828, &b14, &b789}
	var b784 = sequenceBuilder{id: 784, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b783 = charBuilder{}
	b784.items = []builder{&b783}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b791.items = []builder{&b828, &b14}
	b792.items = []builder{&b828, &b14, &b791}
	var b786 = sequenceBuilder{id: 786, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b785 = charBuilder{}
	b786.items = []builder{&b785}
	var b765 = sequenceBuilder{id: 765, commit: 258, ranges: [][]int{{1, 1}, {0, 1}}}
	var b764 = sequenceBuilder{id: 764, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b762 = sequenceBuilder{id: 762, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b762.items = []builder{&b224, &b828, &b742}
	var b763 = sequenceBuilder{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b763.items = []builder{&b828, &b762}
	b764.items = []builder{&b828, &b762, &b763}
	b765.items = []builder{&b742, &b764}
	var b788 = sequenceBuilder{id: 788, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b787 = charBuilder{}
	b788.items = []builder{&b787}
	b793.items = []builder{&b164, &b790, &b828, &b784, &b792, &b828, &b786, &b828, &b224, &b828, &b765, &b828, &b224, &b828, &b788}
	b794.options = []builder{&b775, &b782, &b793}
	var b807 = sequenceBuilder{id: 807, commit: 258, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{798}}
	var b800 = sequenceBuilder{id: 800, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b799 = charBuilder{}
	b800.items = []builder{&b799}
	var b804 = sequenceBuilder{id: 804, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b803 = sequenceBuilder{id: 803, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b803.items = []builder{&b828, &b14}
	b804.items = []builder{&b828, &b14, &b803}
	var b806 = sequenceBuilder{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b805 = sequenceBuilder{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b805.items = []builder{&b828, &b14}
	b806.items = []builder{&b828, &b14, &b805}
	var b802 = sequenceBuilder{id: 802, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b801 = charBuilder{}
	b802.items = []builder{&b801}
	b807.items = []builder{&b800, &b804, &b828, &b798, &b806, &b828, &b802}
	b798.options = []builder{&b289, &b292, &b372, &b501, &b532, &b569, &b574, &b603, &b606, &b609, &b610, &b629, &b730, &b797, &b794, &b807, &b293}
	var b815 = sequenceBuilder{id: 815, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b813 = sequenceBuilder{id: 813, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b813.items = []builder{&b812, &b828, &b798}
	var b814 = sequenceBuilder{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b814.items = []builder{&b828, &b813}
	b815.items = []builder{&b828, &b813, &b814}
	b816.items = []builder{&b798, &b815}
	b829.items = []builder{&b825, &b828, &b812, &b828, &b816, &b828, &b812}
	b830.items = []builder{&b828, &b829, &b828}

	var keywords = []parser{&p48, &p54, &p61, &p67, &p70, &p73, &p78, &p85, &p90, &p98, &p103, &p111, &p118, &p121, &p127, &p130, &p134, &p140, &p149, &p153, &p160, &p164}

	return parseInput(r, &p830, &b830, keywords)
}
