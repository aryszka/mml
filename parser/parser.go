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
	offset        int
	readOffset    int
	consumed      int
	failOffset    int
	failingParser parser
	readErr       error
	eof           bool
	results       *results
	tokens        []rune
	matchLast     bool
}

func newContext(r io.RuneReader) *context {
	return &context{reader: r, results: &results{}, failOffset: -1}
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
func parseInput(r io.Reader, p parser, b builder) (*Node, error) {
	c := newContext(bufio.NewReader(r))
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
	var p831 = sequenceParser{id: 831, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p829 = choiceParser{id: 829, commit: 2}
	var p827 = choiceParser{id: 827, commit: 70, name: "ws", generalizations: []int{829}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827, 829}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827, 829}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827, 829}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827, 829}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827, 829}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827, 829}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p827.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p828 = sequenceParser{id: 828, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{829}}
	var p42 = sequenceParser{id: 42, commit: 66, name: "comment", ranges: [][]int{{1, 1}, {0, 1}}}
	var p38 = choiceParser{id: 38, commit: 66, name: "comment-part"}
	var p21 = sequenceParser{id: 21, commit: 74, name: "line-comment", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{38}}
	var p20 = sequenceParser{id: 20, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p18 = charParser{id: 18, chars: []rune{47}}
	var p19 = charParser{id: 19, chars: []rune{47}}
	p20.items = []parser{&p18, &p19}
	var p17 = sequenceParser{id: 17, commit: 72, name: "line-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var p16 = sequenceParser{id: 16, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p15 = charParser{id: 15, not: true, chars: []rune{10}}
	p16.items = []parser{&p15}
	p17.items = []parser{&p16}
	p21.items = []parser{&p20, &p17}
	var p37 = sequenceParser{id: 37, commit: 74, name: "block-comment", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{38}}
	var p33 = sequenceParser{id: 33, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p31 = charParser{id: 31, chars: []rune{47}}
	var p32 = charParser{id: 32, chars: []rune{42}}
	p33.items = []parser{&p31, &p32}
	var p30 = sequenceParser{id: 30, commit: 72, name: "block-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var p29 = choiceParser{id: 29, commit: 10}
	var p23 = sequenceParser{id: 23, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{29}}
	var p22 = charParser{id: 22, not: true, chars: []rune{42}}
	p23.items = []parser{&p22}
	var p28 = sequenceParser{id: 28, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{29}}
	var p25 = sequenceParser{id: 25, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p24 = charParser{id: 24, chars: []rune{42}}
	p25.items = []parser{&p24}
	var p27 = sequenceParser{id: 27, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p26 = charParser{id: 26, not: true, chars: []rune{47}}
	p27.items = []parser{&p26}
	p28.items = []parser{&p25, &p27}
	p29.options = []parser{&p23, &p28}
	p30.items = []parser{&p29}
	var p36 = sequenceParser{id: 36, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p34 = charParser{id: 34, chars: []rune{42}}
	var p35 = charParser{id: 35, chars: []rune{47}}
	p36.items = []parser{&p34, &p35}
	p37.items = []parser{&p33, &p30, &p36}
	p38.options = []parser{&p21, &p37}
	var p41 = sequenceParser{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p39 = sequenceParser{id: 39, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{809, 111}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p39.items = []parser{&p14, &p829, &p38}
	var p40 = sequenceParser{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p40.items = []parser{&p829, &p39}
	p41.items = []parser{&p829, &p39, &p40}
	p42.items = []parser{&p38, &p41}
	p828.items = []parser{&p42}
	p829.options = []parser{&p827, &p828}
	var p830 = sequenceParser{id: 830, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p826 = sequenceParser{id: 826, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p823 = sequenceParser{id: 823, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p821 = charParser{id: 821, chars: []rune{35}}
	var p822 = charParser{id: 822, chars: []rune{33}}
	p823.items = []parser{&p821, &p822}
	var p820 = sequenceParser{id: 820, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p819 = sequenceParser{id: 819, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p817 = sequenceParser{id: 817, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p816 = charParser{id: 816, not: true, chars: []rune{10}}
	p817.items = []parser{&p816}
	var p818 = sequenceParser{id: 818, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p818.items = []parser{&p829, &p817}
	p819.items = []parser{&p817, &p818}
	p820.items = []parser{&p819}
	var p825 = sequenceParser{id: 825, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p824 = charParser{id: 824, chars: []rune{10}}
	p825.items = []parser{&p824}
	p826.items = []parser{&p823, &p829, &p820, &p829, &p825}
	var p811 = sequenceParser{id: 811, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p809 = choiceParser{id: 809, commit: 2}
	var p808 = sequenceParser{id: 808, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{809}}
	var p807 = charParser{id: 807, chars: []rune{59}}
	p808.items = []parser{&p807}
	p809.options = []parser{&p808, &p14}
	var p810 = sequenceParser{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p810.items = []parser{&p829, &p809}
	p811.items = []parser{&p809, &p810}
	var p815 = sequenceParser{id: 815, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p797 = choiceParser{id: 797, commit: 66, name: "statement", generalizations: []int{479, 543}}
	var p185 = sequenceParser{id: 185, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{797, 479, 543}}
	var p181 = sequenceParser{id: 181, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p175 = charParser{id: 175, chars: []rune{114}}
	var p176 = charParser{id: 176, chars: []rune{101}}
	var p177 = charParser{id: 177, chars: []rune{116}}
	var p178 = charParser{id: 178, chars: []rune{117}}
	var p179 = charParser{id: 179, chars: []rune{114}}
	var p180 = charParser{id: 180, chars: []rune{110}}
	p181.items = []parser{&p175, &p176, &p177, &p178, &p179, &p180}
	var p184 = sequenceParser{id: 184, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p183 = sequenceParser{id: 183, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p182 = sequenceParser{id: 182, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p182.items = []parser{&p829, &p14}
	p183.items = []parser{&p14, &p182}
	var p396 = choiceParser{id: 396, commit: 66, name: "expression", generalizations: []int{114, 787, 197, 578, 571, 797}}
	var p267 = choiceParser{id: 267, commit: 66, name: "primary-expression", generalizations: []int{114, 787, 197, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p247 = choiceParser{id: 247, commit: 66, name: "subprimary-expression", generalizations: []int{114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p60 = choiceParser{id: 60, commit: 64, name: "int", generalizations: []int{247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p51 = sequenceParser{id: 51, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p50 = sequenceParser{id: 50, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p49 = charParser{id: 49, ranges: [][]rune{{49, 57}}}
	p50.items = []parser{&p49}
	var p44 = sequenceParser{id: 44, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p43 = charParser{id: 43, ranges: [][]rune{{48, 57}}}
	p44.items = []parser{&p43}
	p51.items = []parser{&p50, &p44}
	var p54 = sequenceParser{id: 54, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p53 = sequenceParser{id: 53, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p52 = charParser{id: 52, chars: []rune{48}}
	p53.items = []parser{&p52}
	var p46 = sequenceParser{id: 46, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 55}}}
	p46.items = []parser{&p45}
	p54.items = []parser{&p53, &p46}
	var p59 = sequenceParser{id: 59, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{60, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p56 = sequenceParser{id: 56, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p55 = charParser{id: 55, chars: []rune{48}}
	p56.items = []parser{&p55}
	var p58 = sequenceParser{id: 58, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p57 = charParser{id: 57, chars: []rune{120, 88}}
	p58.items = []parser{&p57}
	var p48 = sequenceParser{id: 48, commit: 66, name: "hexa-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p47 = charParser{id: 47, ranges: [][]rune{{48, 57}, {97, 102}, {65, 70}}}
	p48.items = []parser{&p47}
	p59.items = []parser{&p56, &p58, &p48}
	p60.options = []parser{&p51, &p54, &p59}
	var p73 = choiceParser{id: 73, commit: 72, name: "float", generalizations: []int{247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p68 = sequenceParser{id: 68, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{73, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p67 = sequenceParser{id: 67, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p66 = charParser{id: 66, chars: []rune{46}}
	p67.items = []parser{&p66}
	var p65 = sequenceParser{id: 65, commit: 74, name: "exponent", ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var p62 = sequenceParser{id: 62, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p61 = charParser{id: 61, chars: []rune{101, 69}}
	p62.items = []parser{&p61}
	var p64 = sequenceParser{id: 64, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p63 = charParser{id: 63, chars: []rune{43, 45}}
	p64.items = []parser{&p63}
	p65.items = []parser{&p62, &p64, &p44}
	p68.items = []parser{&p44, &p67, &p44, &p65}
	var p71 = sequenceParser{id: 71, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{73, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p70 = sequenceParser{id: 70, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p69 = charParser{id: 69, chars: []rune{46}}
	p70.items = []parser{&p69}
	p71.items = []parser{&p70, &p44, &p65}
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{73, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	p72.items = []parser{&p44, &p65}
	p73.options = []parser{&p68, &p71, &p72}
	var p86 = sequenceParser{id: 86, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 114, 139, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 753, 797}}
	var p75 = sequenceParser{id: 75, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p74 = charParser{id: 74, chars: []rune{34}}
	p75.items = []parser{&p74}
	var p83 = choiceParser{id: 83, commit: 10}
	var p77 = sequenceParser{id: 77, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{83}}
	var p76 = charParser{id: 76, not: true, chars: []rune{92, 34}}
	p77.items = []parser{&p76}
	var p82 = sequenceParser{id: 82, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{83}}
	var p79 = sequenceParser{id: 79, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p78 = charParser{id: 78, chars: []rune{92}}
	p79.items = []parser{&p78}
	var p81 = sequenceParser{id: 81, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p80 = charParser{id: 80, not: true}
	p81.items = []parser{&p80}
	p82.items = []parser{&p79, &p81}
	p83.options = []parser{&p77, &p82}
	var p85 = sequenceParser{id: 85, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p84 = charParser{id: 84, chars: []rune{34}}
	p85.items = []parser{&p84}
	p86.items = []parser{&p75, &p83, &p85}
	var p98 = choiceParser{id: 98, commit: 66, name: "bool", generalizations: []int{247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p91 = sequenceParser{id: 91, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p87 = charParser{id: 87, chars: []rune{116}}
	var p88 = charParser{id: 88, chars: []rune{114}}
	var p89 = charParser{id: 89, chars: []rune{117}}
	var p90 = charParser{id: 90, chars: []rune{101}}
	p91.items = []parser{&p87, &p88, &p89, &p90}
	var p97 = sequenceParser{id: 97, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p92 = charParser{id: 92, chars: []rune{102}}
	var p93 = charParser{id: 93, chars: []rune{97}}
	var p94 = charParser{id: 94, chars: []rune{108}}
	var p95 = charParser{id: 95, chars: []rune{115}}
	var p96 = charParser{id: 96, chars: []rune{101}}
	p97.items = []parser{&p92, &p93, &p94, &p95, &p96}
	p98.options = []parser{&p91, &p97}
	var p511 = sequenceParser{id: 511, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 114, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 515, 578, 571, 797}}
	var p508 = sequenceParser{id: 508, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p501 = charParser{id: 501, chars: []rune{114}}
	var p502 = charParser{id: 502, chars: []rune{101}}
	var p503 = charParser{id: 503, chars: []rune{99}}
	var p504 = charParser{id: 504, chars: []rune{101}}
	var p505 = charParser{id: 505, chars: []rune{105}}
	var p506 = charParser{id: 506, chars: []rune{118}}
	var p507 = charParser{id: 507, chars: []rune{101}}
	p508.items = []parser{&p501, &p502, &p503, &p504, &p505, &p506, &p507}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p509.items = []parser{&p829, &p14}
	p510.items = []parser{&p829, &p14, &p509}
	p511.items = []parser{&p508, &p510, &p829, &p267}
	var p103 = sequenceParser{id: 103, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{247, 114, 139, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 744, 797}}
	var p100 = sequenceParser{id: 100, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p99 = charParser{id: 99, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p100.items = []parser{&p99}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p102.items = []parser{&p101}
	p103.items = []parser{&p100, &p102}
	var p124 = sequenceParser{id: 124, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{114, 247, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p123 = sequenceParser{id: 123, commit: 66, name: "list-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p120 = sequenceParser{id: 120, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p119 = charParser{id: 119, chars: []rune{91}}
	p120.items = []parser{&p119}
	var p113 = sequenceParser{id: 113, commit: 66, name: "list-sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p111 = choiceParser{id: 111, commit: 2}
	var p110 = sequenceParser{id: 110, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{111}}
	var p109 = charParser{id: 109, chars: []rune{44}}
	p110.items = []parser{&p109}
	p111.options = []parser{&p14, &p110}
	var p112 = sequenceParser{id: 112, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p112.items = []parser{&p829, &p111}
	p113.items = []parser{&p111, &p112}
	var p118 = sequenceParser{id: 118, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p114 = choiceParser{id: 114, commit: 66, name: "list-item"}
	var p108 = sequenceParser{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{114, 147, 148}}
	var p107 = sequenceParser{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p104 = charParser{id: 104, chars: []rune{46}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	p107.items = []parser{&p104, &p105, &p106}
	p108.items = []parser{&p267, &p829, &p107}
	p114.options = []parser{&p396, &p108}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p115 = sequenceParser{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p115.items = []parser{&p113, &p829, &p114}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p116.items = []parser{&p829, &p115}
	p117.items = []parser{&p829, &p115, &p116}
	p118.items = []parser{&p114, &p117}
	var p122 = sequenceParser{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p121 = charParser{id: 121, chars: []rune{93}}
	p122.items = []parser{&p121}
	p123.items = []parser{&p120, &p829, &p113, &p829, &p118, &p829, &p113, &p829, &p122}
	p124.items = []parser{&p123}
	var p129 = sequenceParser{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p126 = sequenceParser{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p125 = charParser{id: 125, chars: []rune{126}}
	p126.items = []parser{&p125}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p127 = sequenceParser{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p127.items = []parser{&p829, &p14}
	p128.items = []parser{&p829, &p14, &p127}
	p129.items = []parser{&p126, &p128, &p829, &p123}
	var p158 = sequenceParser{id: 158, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{247, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p157 = sequenceParser{id: 157, commit: 66, name: "struct-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p154 = sequenceParser{id: 154, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p153 = charParser{id: 153, chars: []rune{123}}
	p154.items = []parser{&p153}
	var p152 = sequenceParser{id: 152, commit: 66, name: "entry-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p147 = choiceParser{id: 147, commit: 2}
	var p146 = sequenceParser{id: 146, commit: 64, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{147, 148}}
	var p139 = choiceParser{id: 139, commit: 2}
	var p138 = sequenceParser{id: 138, commit: 64, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{139}}
	var p131 = sequenceParser{id: 131, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p130 = charParser{id: 130, chars: []rune{91}}
	p131.items = []parser{&p130}
	var p135 = sequenceParser{id: 135, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p134 = sequenceParser{id: 134, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p134.items = []parser{&p829, &p14}
	p135.items = []parser{&p829, &p14, &p134}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p136 = sequenceParser{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p136.items = []parser{&p829, &p14}
	p137.items = []parser{&p829, &p14, &p136}
	var p133 = sequenceParser{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p132 = charParser{id: 132, chars: []rune{93}}
	p133.items = []parser{&p132}
	p138.items = []parser{&p131, &p135, &p829, &p396, &p137, &p829, &p133}
	p139.options = []parser{&p103, &p86, &p138}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p142 = sequenceParser{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p142.items = []parser{&p829, &p14}
	p143.items = []parser{&p829, &p14, &p142}
	var p141 = sequenceParser{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p140 = charParser{id: 140, chars: []rune{58}}
	p141.items = []parser{&p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p829, &p14}
	p145.items = []parser{&p829, &p14, &p144}
	p146.items = []parser{&p139, &p143, &p829, &p141, &p145, &p829, &p396}
	p147.options = []parser{&p146, &p108}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p149 = sequenceParser{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p148 = choiceParser{id: 148, commit: 2}
	p148.options = []parser{&p146, &p108}
	p149.items = []parser{&p113, &p829, &p148}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p150.items = []parser{&p829, &p149}
	p151.items = []parser{&p829, &p149, &p150}
	p152.items = []parser{&p147, &p151}
	var p156 = sequenceParser{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p155 = charParser{id: 155, chars: []rune{125}}
	p156.items = []parser{&p155}
	p157.items = []parser{&p154, &p829, &p113, &p829, &p152, &p829, &p113, &p829, &p156}
	p158.items = []parser{&p157}
	var p163 = sequenceParser{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 787, 197, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p160 = sequenceParser{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p159 = charParser{id: 159, chars: []rune{126}}
	p160.items = []parser{&p159}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p161 = sequenceParser{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p161.items = []parser{&p829, &p14}
	p162.items = []parser{&p829, &p14, &p161}
	p163.items = []parser{&p160, &p162, &p829, &p157}
	var p206 = sequenceParser{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{787, 197, 247, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p203 = sequenceParser{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p201 = charParser{id: 201, chars: []rune{102}}
	var p202 = charParser{id: 202, chars: []rune{110}}
	p203.items = []parser{&p201, &p202}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p204 = sequenceParser{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p204.items = []parser{&p829, &p14}
	p205.items = []parser{&p829, &p14, &p204}
	var p200 = sequenceParser{id: 200, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p192 = sequenceParser{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p191 = charParser{id: 191, chars: []rune{40}}
	p192.items = []parser{&p191}
	var p194 = choiceParser{id: 194, commit: 2}
	var p167 = sequenceParser{id: 167, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{194}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p164.items = []parser{&p113, &p829, &p103}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p165.items = []parser{&p829, &p164}
	p166.items = []parser{&p829, &p164, &p165}
	p167.items = []parser{&p103, &p166}
	var p193 = sequenceParser{id: 193, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{194}}
	var p174 = sequenceParser{id: 174, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{194}}
	var p171 = sequenceParser{id: 171, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p168 = charParser{id: 168, chars: []rune{46}}
	var p169 = charParser{id: 169, chars: []rune{46}}
	var p170 = charParser{id: 170, chars: []rune{46}}
	p171.items = []parser{&p168, &p169, &p170}
	var p173 = sequenceParser{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p172 = sequenceParser{id: 172, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p172.items = []parser{&p829, &p14}
	p173.items = []parser{&p829, &p14, &p172}
	p174.items = []parser{&p171, &p173, &p829, &p103}
	p193.items = []parser{&p167, &p829, &p113, &p829, &p174}
	p194.options = []parser{&p167, &p193, &p174}
	var p196 = sequenceParser{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p195 = charParser{id: 195, chars: []rune{41}}
	p196.items = []parser{&p195}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p198 = sequenceParser{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p198.items = []parser{&p829, &p14}
	p199.items = []parser{&p829, &p14, &p198}
	var p197 = choiceParser{id: 197, commit: 2}
	var p787 = choiceParser{id: 787, commit: 66, name: "simple-statement", generalizations: []int{197, 797}}
	var p500 = sequenceParser{id: 500, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{787, 197, 515, 797}}
	var p495 = sequenceParser{id: 495, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p491 = charParser{id: 491, chars: []rune{115}}
	var p492 = charParser{id: 492, chars: []rune{101}}
	var p493 = charParser{id: 493, chars: []rune{110}}
	var p494 = charParser{id: 494, chars: []rune{100}}
	p495.items = []parser{&p491, &p492, &p493, &p494}
	var p497 = sequenceParser{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p496 = sequenceParser{id: 496, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p496.items = []parser{&p829, &p14}
	p497.items = []parser{&p829, &p14, &p496}
	var p499 = sequenceParser{id: 499, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p498 = sequenceParser{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p498.items = []parser{&p829, &p14}
	p499.items = []parser{&p829, &p14, &p498}
	p500.items = []parser{&p495, &p497, &p829, &p267, &p499, &p829, &p267}
	var p558 = sequenceParser{id: 558, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{787, 197, 797}}
	var p555 = sequenceParser{id: 555, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p553 = charParser{id: 553, chars: []rune{103}}
	var p554 = charParser{id: 554, chars: []rune{111}}
	p555.items = []parser{&p553, &p554}
	var p557 = sequenceParser{id: 557, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p556 = sequenceParser{id: 556, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p556.items = []parser{&p829, &p14}
	p557.items = []parser{&p829, &p14, &p556}
	var p237 = sequenceParser{id: 237, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p234 = sequenceParser{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p233 = charParser{id: 233, chars: []rune{40}}
	p234.items = []parser{&p233}
	var p236 = sequenceParser{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p235 = charParser{id: 235, chars: []rune{41}}
	p236.items = []parser{&p235}
	p237.items = []parser{&p267, &p829, &p234, &p829, &p113, &p829, &p118, &p829, &p113, &p829, &p236}
	p558.items = []parser{&p555, &p557, &p829, &p237}
	var p567 = sequenceParser{id: 567, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{787, 197, 797}}
	var p564 = sequenceParser{id: 564, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p559 = charParser{id: 559, chars: []rune{100}}
	var p560 = charParser{id: 560, chars: []rune{101}}
	var p561 = charParser{id: 561, chars: []rune{102}}
	var p562 = charParser{id: 562, chars: []rune{101}}
	var p563 = charParser{id: 563, chars: []rune{114}}
	p564.items = []parser{&p559, &p560, &p561, &p562, &p563}
	var p566 = sequenceParser{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p565.items = []parser{&p829, &p14}
	p566.items = []parser{&p829, &p14, &p565}
	p567.items = []parser{&p564, &p566, &p829, &p237}
	var p632 = choiceParser{id: 632, commit: 64, name: "assignment", generalizations: []int{787, 197, 797}}
	var p612 = sequenceParser{id: 612, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{632, 787, 197, 797}}
	var p609 = sequenceParser{id: 609, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p606 = charParser{id: 606, chars: []rune{115}}
	var p607 = charParser{id: 607, chars: []rune{101}}
	var p608 = charParser{id: 608, chars: []rune{116}}
	p609.items = []parser{&p606, &p607, &p608}
	var p611 = sequenceParser{id: 611, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p610 = sequenceParser{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p610.items = []parser{&p829, &p14}
	p611.items = []parser{&p829, &p14, &p610}
	var p601 = sequenceParser{id: 601, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p596.items = []parser{&p829, &p14}
	p597.items = []parser{&p14, &p596}
	var p595 = sequenceParser{id: 595, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p594 = charParser{id: 594, chars: []rune{61}}
	p595.items = []parser{&p594}
	p598.items = []parser{&p597, &p829, &p595}
	var p600 = sequenceParser{id: 600, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p599 = sequenceParser{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p599.items = []parser{&p829, &p14}
	p600.items = []parser{&p829, &p14, &p599}
	p601.items = []parser{&p267, &p829, &p598, &p600, &p829, &p396}
	p612.items = []parser{&p609, &p611, &p829, &p601}
	var p619 = sequenceParser{id: 619, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{632, 787, 197, 797}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p615.items = []parser{&p829, &p14}
	p616.items = []parser{&p829, &p14, &p615}
	var p614 = sequenceParser{id: 614, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p613 = charParser{id: 613, chars: []rune{61}}
	p614.items = []parser{&p613}
	var p618 = sequenceParser{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p617.items = []parser{&p829, &p14}
	p618.items = []parser{&p829, &p14, &p617}
	p619.items = []parser{&p267, &p616, &p829, &p614, &p618, &p829, &p396}
	var p631 = sequenceParser{id: 631, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{632, 787, 197, 797}}
	var p623 = sequenceParser{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p620 = charParser{id: 620, chars: []rune{115}}
	var p621 = charParser{id: 621, chars: []rune{101}}
	var p622 = charParser{id: 622, chars: []rune{116}}
	p623.items = []parser{&p620, &p621, &p622}
	var p630 = sequenceParser{id: 630, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p629 = sequenceParser{id: 629, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p629.items = []parser{&p829, &p14}
	p630.items = []parser{&p829, &p14, &p629}
	var p625 = sequenceParser{id: 625, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p624 = charParser{id: 624, chars: []rune{40}}
	p625.items = []parser{&p624}
	var p626 = sequenceParser{id: 626, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p605 = sequenceParser{id: 605, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p604 = sequenceParser{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p602 = sequenceParser{id: 602, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p602.items = []parser{&p113, &p829, &p601}
	var p603 = sequenceParser{id: 603, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p603.items = []parser{&p829, &p602}
	p604.items = []parser{&p829, &p602, &p603}
	p605.items = []parser{&p601, &p604}
	p626.items = []parser{&p113, &p829, &p605}
	var p628 = sequenceParser{id: 628, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p627 = charParser{id: 627, chars: []rune{41}}
	p628.items = []parser{&p627}
	p631.items = []parser{&p623, &p630, &p829, &p625, &p829, &p626, &p829, &p113, &p829, &p628}
	p632.options = []parser{&p612, &p619, &p631}
	var p796 = sequenceParser{id: 796, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{787, 197, 797}}
	var p789 = sequenceParser{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p788 = charParser{id: 788, chars: []rune{40}}
	p789.items = []parser{&p788}
	var p793 = sequenceParser{id: 793, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p792.items = []parser{&p829, &p14}
	p793.items = []parser{&p829, &p14, &p792}
	var p795 = sequenceParser{id: 795, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p794 = sequenceParser{id: 794, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p794.items = []parser{&p829, &p14}
	p795.items = []parser{&p829, &p14, &p794}
	var p791 = sequenceParser{id: 791, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p790 = charParser{id: 790, chars: []rune{41}}
	p791.items = []parser{&p790}
	p796.items = []parser{&p789, &p793, &p829, &p787, &p795, &p829, &p791}
	p787.options = []parser{&p500, &p558, &p567, &p632, &p796, &p396}
	var p190 = sequenceParser{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{197}}
	var p187 = sequenceParser{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p186 = charParser{id: 186, chars: []rune{123}}
	p187.items = []parser{&p186}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{125}}
	p189.items = []parser{&p188}
	p190.items = []parser{&p187, &p829, &p811, &p829, &p815, &p829, &p811, &p829, &p189}
	p197.options = []parser{&p787, &p190}
	p200.items = []parser{&p192, &p829, &p113, &p829, &p194, &p829, &p113, &p829, &p196, &p199, &p829, &p197}
	p206.items = []parser{&p203, &p205, &p829, &p200}
	var p216 = sequenceParser{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p209 = sequenceParser{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p207 = charParser{id: 207, chars: []rune{102}}
	var p208 = charParser{id: 208, chars: []rune{110}}
	p209.items = []parser{&p207, &p208}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p212 = sequenceParser{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p212.items = []parser{&p829, &p14}
	p213.items = []parser{&p829, &p14, &p212}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p210 = charParser{id: 210, chars: []rune{126}}
	p211.items = []parser{&p210}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p214.items = []parser{&p829, &p14}
	p215.items = []parser{&p829, &p14, &p214}
	p216.items = []parser{&p209, &p213, &p829, &p211, &p215, &p829, &p200}
	var p232 = sequenceParser{id: 232, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p229 = sequenceParser{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p228 = sequenceParser{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p228.items = []parser{&p829, &p14}
	p229.items = []parser{&p829, &p14, &p228}
	var p227 = sequenceParser{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p226 = charParser{id: 226, chars: []rune{46}}
	p227.items = []parser{&p226}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p230 = sequenceParser{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p230.items = []parser{&p829, &p14}
	p231.items = []parser{&p829, &p14, &p230}
	p232.items = []parser{&p267, &p229, &p829, &p227, &p231, &p829, &p103}
	var p246 = sequenceParser{id: 246, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{247, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p239 = sequenceParser{id: 239, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p238 = charParser{id: 238, chars: []rune{40}}
	p239.items = []parser{&p238}
	var p243 = sequenceParser{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p242 = sequenceParser{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p242.items = []parser{&p829, &p14}
	p243.items = []parser{&p829, &p14, &p242}
	var p245 = sequenceParser{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p244 = sequenceParser{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p244.items = []parser{&p829, &p14}
	p245.items = []parser{&p829, &p14, &p244}
	var p241 = sequenceParser{id: 241, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p240 = charParser{id: 240, chars: []rune{41}}
	p241.items = []parser{&p240}
	p246.items = []parser{&p239, &p243, &p829, &p396, &p245, &p829, &p241}
	p247.options = []parser{&p60, &p73, &p86, &p98, &p511, &p103, &p124, &p129, &p158, &p163, &p206, &p216, &p232, &p237, &p246}
	var p266 = choiceParser{id: 266, commit: 64, name: "expression-indexer", generalizations: []int{267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p256 = sequenceParser{id: 256, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p249 = sequenceParser{id: 249, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p248 = charParser{id: 248, chars: []rune{91}}
	p249.items = []parser{&p248}
	var p253 = sequenceParser{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p252 = sequenceParser{id: 252, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p252.items = []parser{&p829, &p14}
	p253.items = []parser{&p829, &p14, &p252}
	var p255 = sequenceParser{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p254 = sequenceParser{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p254.items = []parser{&p829, &p14}
	p255.items = []parser{&p829, &p14, &p254}
	var p251 = sequenceParser{id: 251, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p250 = charParser{id: 250, chars: []rune{93}}
	p251.items = []parser{&p250}
	p256.items = []parser{&p247, &p829, &p249, &p253, &p829, &p396, &p255, &p829, &p251}
	var p265 = sequenceParser{id: 265, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 267, 396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p258 = sequenceParser{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p257 = charParser{id: 257, chars: []rune{91}}
	p258.items = []parser{&p257}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p261 = sequenceParser{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p261.items = []parser{&p829, &p14}
	p262.items = []parser{&p829, &p14, &p261}
	var p225 = sequenceParser{id: 225, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{571, 577, 578}}
	var p217 = sequenceParser{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p217.items = []parser{&p396}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p221 = sequenceParser{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p221.items = []parser{&p829, &p14}
	p222.items = []parser{&p829, &p14, &p221}
	var p220 = sequenceParser{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p219 = charParser{id: 219, chars: []rune{58}}
	p220.items = []parser{&p219}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p829, &p14}
	p224.items = []parser{&p829, &p14, &p223}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p396}
	p225.items = []parser{&p217, &p222, &p829, &p220, &p224, &p829, &p218}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p263.items = []parser{&p829, &p14}
	p264.items = []parser{&p829, &p14, &p263}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{93}}
	p260.items = []parser{&p259}
	p265.items = []parser{&p247, &p829, &p258, &p262, &p829, &p225, &p264, &p829, &p260}
	p266.options = []parser{&p256, &p265}
	p267.options = []parser{&p247, &p266}
	var p327 = sequenceParser{id: 327, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{396, 333, 334, 335, 336, 337, 388, 578, 571, 797}}
	var p326 = choiceParser{id: 326, commit: 66, name: "unary-operator"}
	var p286 = sequenceParser{id: 286, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{326}}
	var p285 = charParser{id: 285, chars: []rune{43}}
	p286.items = []parser{&p285}
	var p288 = sequenceParser{id: 288, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{326}}
	var p287 = charParser{id: 287, chars: []rune{45}}
	p288.items = []parser{&p287}
	var p269 = sequenceParser{id: 269, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{326}}
	var p268 = charParser{id: 268, chars: []rune{94}}
	p269.items = []parser{&p268}
	var p300 = sequenceParser{id: 300, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{326}}
	var p299 = charParser{id: 299, chars: []rune{33}}
	p300.items = []parser{&p299}
	p326.options = []parser{&p286, &p288, &p269, &p300}
	p327.items = []parser{&p326, &p829, &p267}
	var p374 = choiceParser{id: 374, commit: 66, name: "binary-expression", generalizations: []int{396, 388, 578, 571, 797}}
	var p345 = sequenceParser{id: 345, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 334, 335, 336, 337, 396, 388, 578, 571, 797}}
	var p333 = choiceParser{id: 333, commit: 66, name: "operand0", generalizations: []int{334, 335, 336, 337}}
	p333.options = []parser{&p267, &p327}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p340 = sequenceParser{id: 340, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p339 = sequenceParser{id: 339, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p339.items = []parser{&p829, &p14}
	p340.items = []parser{&p14, &p339}
	var p328 = choiceParser{id: 328, commit: 66, name: "binary-op0"}
	var p271 = sequenceParser{id: 271, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p270 = charParser{id: 270, chars: []rune{38}}
	p271.items = []parser{&p270}
	var p278 = sequenceParser{id: 278, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{328}}
	var p276 = charParser{id: 276, chars: []rune{38}}
	var p277 = charParser{id: 277, chars: []rune{94}}
	p278.items = []parser{&p276, &p277}
	var p281 = sequenceParser{id: 281, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{328}}
	var p279 = charParser{id: 279, chars: []rune{60}}
	var p280 = charParser{id: 280, chars: []rune{60}}
	p281.items = []parser{&p279, &p280}
	var p284 = sequenceParser{id: 284, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{328}}
	var p282 = charParser{id: 282, chars: []rune{62}}
	var p283 = charParser{id: 283, chars: []rune{62}}
	p284.items = []parser{&p282, &p283}
	var p290 = sequenceParser{id: 290, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p289 = charParser{id: 289, chars: []rune{42}}
	p290.items = []parser{&p289}
	var p292 = sequenceParser{id: 292, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p291 = charParser{id: 291, chars: []rune{47}}
	p292.items = []parser{&p291}
	var p294 = sequenceParser{id: 294, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p293 = charParser{id: 293, chars: []rune{37}}
	p294.items = []parser{&p293}
	p328.options = []parser{&p271, &p278, &p281, &p284, &p290, &p292, &p294}
	var p342 = sequenceParser{id: 342, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p341 = sequenceParser{id: 341, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p341.items = []parser{&p829, &p14}
	p342.items = []parser{&p829, &p14, &p341}
	p343.items = []parser{&p340, &p829, &p328, &p342, &p829, &p333}
	var p344 = sequenceParser{id: 344, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p344.items = []parser{&p829, &p343}
	p345.items = []parser{&p333, &p829, &p343, &p344}
	var p352 = sequenceParser{id: 352, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 335, 336, 337, 396, 388, 578, 571, 797}}
	var p334 = choiceParser{id: 334, commit: 66, name: "operand1", generalizations: []int{335, 336, 337}}
	p334.options = []parser{&p333, &p345}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p346.items = []parser{&p829, &p14}
	p347.items = []parser{&p14, &p346}
	var p329 = choiceParser{id: 329, commit: 66, name: "binary-op1"}
	var p273 = sequenceParser{id: 273, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p272 = charParser{id: 272, chars: []rune{124}}
	p273.items = []parser{&p272}
	var p275 = sequenceParser{id: 275, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p274 = charParser{id: 274, chars: []rune{94}}
	p275.items = []parser{&p274}
	var p296 = sequenceParser{id: 296, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p295 = charParser{id: 295, chars: []rune{43}}
	p296.items = []parser{&p295}
	var p298 = sequenceParser{id: 298, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p297 = charParser{id: 297, chars: []rune{45}}
	p298.items = []parser{&p297}
	p329.options = []parser{&p273, &p275, &p296, &p298}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p348.items = []parser{&p829, &p14}
	p349.items = []parser{&p829, &p14, &p348}
	p350.items = []parser{&p347, &p829, &p329, &p349, &p829, &p334}
	var p351 = sequenceParser{id: 351, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p351.items = []parser{&p829, &p350}
	p352.items = []parser{&p334, &p829, &p350, &p351}
	var p359 = sequenceParser{id: 359, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 336, 337, 396, 388, 578, 571, 797}}
	var p335 = choiceParser{id: 335, commit: 66, name: "operand2", generalizations: []int{336, 337}}
	p335.options = []parser{&p334, &p352}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p353.items = []parser{&p829, &p14}
	p354.items = []parser{&p14, &p353}
	var p330 = choiceParser{id: 330, commit: 66, name: "binary-op2"}
	var p303 = sequenceParser{id: 303, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{330}}
	var p301 = charParser{id: 301, chars: []rune{61}}
	var p302 = charParser{id: 302, chars: []rune{61}}
	p303.items = []parser{&p301, &p302}
	var p306 = sequenceParser{id: 306, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{330}}
	var p304 = charParser{id: 304, chars: []rune{33}}
	var p305 = charParser{id: 305, chars: []rune{61}}
	p306.items = []parser{&p304, &p305}
	var p308 = sequenceParser{id: 308, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p307 = charParser{id: 307, chars: []rune{60}}
	p308.items = []parser{&p307}
	var p311 = sequenceParser{id: 311, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{330}}
	var p309 = charParser{id: 309, chars: []rune{60}}
	var p310 = charParser{id: 310, chars: []rune{61}}
	p311.items = []parser{&p309, &p310}
	var p313 = sequenceParser{id: 313, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p312 = charParser{id: 312, chars: []rune{62}}
	p313.items = []parser{&p312}
	var p316 = sequenceParser{id: 316, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{330}}
	var p314 = charParser{id: 314, chars: []rune{62}}
	var p315 = charParser{id: 315, chars: []rune{61}}
	p316.items = []parser{&p314, &p315}
	p330.options = []parser{&p303, &p306, &p308, &p311, &p313, &p316}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p355.items = []parser{&p829, &p14}
	p356.items = []parser{&p829, &p14, &p355}
	p357.items = []parser{&p354, &p829, &p330, &p356, &p829, &p335}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p358.items = []parser{&p829, &p357}
	p359.items = []parser{&p335, &p829, &p357, &p358}
	var p366 = sequenceParser{id: 366, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 337, 396, 388, 578, 571, 797}}
	var p336 = choiceParser{id: 336, commit: 66, name: "operand3", generalizations: []int{337}}
	p336.options = []parser{&p335, &p359}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p360.items = []parser{&p829, &p14}
	p361.items = []parser{&p14, &p360}
	var p331 = sequenceParser{id: 331, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p319 = sequenceParser{id: 319, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p317 = charParser{id: 317, chars: []rune{38}}
	var p318 = charParser{id: 318, chars: []rune{38}}
	p319.items = []parser{&p317, &p318}
	p331.items = []parser{&p319}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p362.items = []parser{&p829, &p14}
	p363.items = []parser{&p829, &p14, &p362}
	p364.items = []parser{&p361, &p829, &p331, &p363, &p829, &p336}
	var p365 = sequenceParser{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p365.items = []parser{&p829, &p364}
	p366.items = []parser{&p336, &p829, &p364, &p365}
	var p373 = sequenceParser{id: 373, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 396, 388, 578, 571, 797}}
	var p337 = choiceParser{id: 337, commit: 66, name: "operand4"}
	p337.options = []parser{&p336, &p366}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p367.items = []parser{&p829, &p14}
	p368.items = []parser{&p14, &p367}
	var p332 = sequenceParser{id: 332, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p322 = sequenceParser{id: 322, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p320 = charParser{id: 320, chars: []rune{124}}
	var p321 = charParser{id: 321, chars: []rune{124}}
	p322.items = []parser{&p320, &p321}
	p332.items = []parser{&p322}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p369.items = []parser{&p829, &p14}
	p370.items = []parser{&p829, &p14, &p369}
	p371.items = []parser{&p368, &p829, &p332, &p370, &p829, &p337}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p372.items = []parser{&p829, &p371}
	p373.items = []parser{&p337, &p829, &p371, &p372}
	p374.options = []parser{&p345, &p352, &p359, &p366, &p373}
	var p387 = sequenceParser{id: 387, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{396, 388, 578, 571, 797}}
	var p380 = sequenceParser{id: 380, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p379 = sequenceParser{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p379.items = []parser{&p829, &p14}
	p380.items = []parser{&p829, &p14, &p379}
	var p376 = sequenceParser{id: 376, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p375 = charParser{id: 375, chars: []rune{63}}
	p376.items = []parser{&p375}
	var p382 = sequenceParser{id: 382, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p381 = sequenceParser{id: 381, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p381.items = []parser{&p829, &p14}
	p382.items = []parser{&p829, &p14, &p381}
	var p384 = sequenceParser{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p383 = sequenceParser{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p383.items = []parser{&p829, &p14}
	p384.items = []parser{&p829, &p14, &p383}
	var p378 = sequenceParser{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p377 = charParser{id: 377, chars: []rune{58}}
	p378.items = []parser{&p377}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p385.items = []parser{&p829, &p14}
	p386.items = []parser{&p829, &p14, &p385}
	p387.items = []parser{&p396, &p380, &p829, &p376, &p382, &p829, &p396, &p384, &p829, &p378, &p386, &p829, &p396}
	var p395 = sequenceParser{id: 395, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{396, 578, 571, 797}}
	var p388 = choiceParser{id: 388, commit: 66, name: "chainingOperand"}
	p388.options = []parser{&p267, &p327, &p374, &p387}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p389.items = []parser{&p829, &p14}
	p390.items = []parser{&p14, &p389}
	var p325 = sequenceParser{id: 325, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p323 = charParser{id: 323, chars: []rune{45}}
	var p324 = charParser{id: 324, chars: []rune{62}}
	p325.items = []parser{&p323, &p324}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p391.items = []parser{&p829, &p14}
	p392.items = []parser{&p829, &p14, &p391}
	p393.items = []parser{&p390, &p829, &p325, &p392, &p829, &p388}
	var p394 = sequenceParser{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p394.items = []parser{&p829, &p393}
	p395.items = []parser{&p388, &p829, &p393, &p394}
	p396.options = []parser{&p267, &p327, &p374, &p387, &p395}
	p184.items = []parser{&p183, &p829, &p396}
	p185.items = []parser{&p181, &p829, &p184}
	var p433 = sequenceParser{id: 433, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{797, 479, 543}}
	var p399 = sequenceParser{id: 399, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p397 = charParser{id: 397, chars: []rune{105}}
	var p398 = charParser{id: 398, chars: []rune{102}}
	p399.items = []parser{&p397, &p398}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p427.items = []parser{&p829, &p14}
	p428.items = []parser{&p829, &p14, &p427}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p429.items = []parser{&p829, &p14}
	p430.items = []parser{&p829, &p14, &p429}
	var p432 = sequenceParser{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p416 = sequenceParser{id: 416, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p409 = sequenceParser{id: 409, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p408.items = []parser{&p829, &p14}
	p409.items = []parser{&p14, &p408}
	var p404 = sequenceParser{id: 404, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p400 = charParser{id: 400, chars: []rune{101}}
	var p401 = charParser{id: 401, chars: []rune{108}}
	var p402 = charParser{id: 402, chars: []rune{115}}
	var p403 = charParser{id: 403, chars: []rune{101}}
	p404.items = []parser{&p400, &p401, &p402, &p403}
	var p411 = sequenceParser{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p410 = sequenceParser{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p410.items = []parser{&p829, &p14}
	p411.items = []parser{&p829, &p14, &p410}
	var p407 = sequenceParser{id: 407, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p405 = charParser{id: 405, chars: []rune{105}}
	var p406 = charParser{id: 406, chars: []rune{102}}
	p407.items = []parser{&p405, &p406}
	var p413 = sequenceParser{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p412 = sequenceParser{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p412.items = []parser{&p829, &p14}
	p413.items = []parser{&p829, &p14, &p412}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p414.items = []parser{&p829, &p14}
	p415.items = []parser{&p829, &p14, &p414}
	p416.items = []parser{&p409, &p829, &p404, &p411, &p829, &p407, &p413, &p829, &p396, &p415, &p829, &p190}
	var p431 = sequenceParser{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p431.items = []parser{&p829, &p416}
	p432.items = []parser{&p829, &p416, &p431}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p422.items = []parser{&p829, &p14}
	p423.items = []parser{&p14, &p422}
	var p421 = sequenceParser{id: 421, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p417 = charParser{id: 417, chars: []rune{101}}
	var p418 = charParser{id: 418, chars: []rune{108}}
	var p419 = charParser{id: 419, chars: []rune{115}}
	var p420 = charParser{id: 420, chars: []rune{101}}
	p421.items = []parser{&p417, &p418, &p419, &p420}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p424.items = []parser{&p829, &p14}
	p425.items = []parser{&p829, &p14, &p424}
	p426.items = []parser{&p423, &p829, &p421, &p425, &p829, &p190}
	p433.items = []parser{&p399, &p428, &p829, &p396, &p430, &p829, &p190, &p432, &p829, &p426}
	var p490 = sequenceParser{id: 490, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{479, 797, 543}}
	var p475 = sequenceParser{id: 475, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p469 = charParser{id: 469, chars: []rune{115}}
	var p470 = charParser{id: 470, chars: []rune{119}}
	var p471 = charParser{id: 471, chars: []rune{105}}
	var p472 = charParser{id: 472, chars: []rune{116}}
	var p473 = charParser{id: 473, chars: []rune{99}}
	var p474 = charParser{id: 474, chars: []rune{104}}
	p475.items = []parser{&p469, &p470, &p471, &p472, &p473, &p474}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p486.items = []parser{&p829, &p14}
	p487.items = []parser{&p829, &p14, &p486}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p488.items = []parser{&p829, &p14}
	p489.items = []parser{&p829, &p14, &p488}
	var p477 = sequenceParser{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p476 = charParser{id: 476, chars: []rune{123}}
	p477.items = []parser{&p476}
	var p483 = sequenceParser{id: 483, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p478 = choiceParser{id: 478, commit: 2}
	var p468 = sequenceParser{id: 468, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{478, 479}}
	var p463 = sequenceParser{id: 463, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p456 = sequenceParser{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p452 = charParser{id: 452, chars: []rune{99}}
	var p453 = charParser{id: 453, chars: []rune{97}}
	var p454 = charParser{id: 454, chars: []rune{115}}
	var p455 = charParser{id: 455, chars: []rune{101}}
	p456.items = []parser{&p452, &p453, &p454, &p455}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p459.items = []parser{&p829, &p14}
	p460.items = []parser{&p829, &p14, &p459}
	var p462 = sequenceParser{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p461 = sequenceParser{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p461.items = []parser{&p829, &p14}
	p462.items = []parser{&p829, &p14, &p461}
	var p458 = sequenceParser{id: 458, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p457 = charParser{id: 457, chars: []rune{58}}
	p458.items = []parser{&p457}
	p463.items = []parser{&p456, &p460, &p829, &p396, &p462, &p829, &p458}
	var p467 = sequenceParser{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p465 = sequenceParser{id: 465, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p464 = charParser{id: 464, chars: []rune{59}}
	p465.items = []parser{&p464}
	var p466 = sequenceParser{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p466.items = []parser{&p829, &p465}
	p467.items = []parser{&p829, &p465, &p466}
	p468.items = []parser{&p463, &p467, &p829, &p797}
	var p451 = sequenceParser{id: 451, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{478, 479, 542, 543}}
	var p446 = sequenceParser{id: 446, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p441 = sequenceParser{id: 441, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p434 = charParser{id: 434, chars: []rune{100}}
	var p435 = charParser{id: 435, chars: []rune{101}}
	var p436 = charParser{id: 436, chars: []rune{102}}
	var p437 = charParser{id: 437, chars: []rune{97}}
	var p438 = charParser{id: 438, chars: []rune{117}}
	var p439 = charParser{id: 439, chars: []rune{108}}
	var p440 = charParser{id: 440, chars: []rune{116}}
	p441.items = []parser{&p434, &p435, &p436, &p437, &p438, &p439, &p440}
	var p445 = sequenceParser{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p444 = sequenceParser{id: 444, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p444.items = []parser{&p829, &p14}
	p445.items = []parser{&p829, &p14, &p444}
	var p443 = sequenceParser{id: 443, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p442 = charParser{id: 442, chars: []rune{58}}
	p443.items = []parser{&p442}
	p446.items = []parser{&p441, &p445, &p829, &p443}
	var p450 = sequenceParser{id: 450, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p448 = sequenceParser{id: 448, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p447 = charParser{id: 447, chars: []rune{59}}
	p448.items = []parser{&p447}
	var p449 = sequenceParser{id: 449, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p449.items = []parser{&p829, &p448}
	p450.items = []parser{&p829, &p448, &p449}
	p451.items = []parser{&p446, &p450, &p829, &p797}
	p478.options = []parser{&p468, &p451}
	var p482 = sequenceParser{id: 482, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p480 = sequenceParser{id: 480, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p479 = choiceParser{id: 479, commit: 2}
	p479.options = []parser{&p468, &p451, &p797}
	p480.items = []parser{&p811, &p829, &p479}
	var p481 = sequenceParser{id: 481, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p481.items = []parser{&p829, &p480}
	p482.items = []parser{&p829, &p480, &p481}
	p483.items = []parser{&p478, &p482}
	var p485 = sequenceParser{id: 485, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p484 = charParser{id: 484, chars: []rune{125}}
	p485.items = []parser{&p484}
	p490.items = []parser{&p475, &p487, &p829, &p396, &p489, &p829, &p477, &p829, &p811, &p829, &p483, &p829, &p811, &p829, &p485}
	var p552 = sequenceParser{id: 552, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{543, 797}}
	var p539 = sequenceParser{id: 539, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p533 = charParser{id: 533, chars: []rune{115}}
	var p534 = charParser{id: 534, chars: []rune{101}}
	var p535 = charParser{id: 535, chars: []rune{108}}
	var p536 = charParser{id: 536, chars: []rune{101}}
	var p537 = charParser{id: 537, chars: []rune{99}}
	var p538 = charParser{id: 538, chars: []rune{116}}
	p539.items = []parser{&p533, &p534, &p535, &p536, &p537, &p538}
	var p551 = sequenceParser{id: 551, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p550 = sequenceParser{id: 550, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p550.items = []parser{&p829, &p14}
	p551.items = []parser{&p829, &p14, &p550}
	var p541 = sequenceParser{id: 541, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p540 = charParser{id: 540, chars: []rune{123}}
	p541.items = []parser{&p540}
	var p547 = sequenceParser{id: 547, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p542 = choiceParser{id: 542, commit: 2}
	var p532 = sequenceParser{id: 532, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{542, 543}}
	var p527 = sequenceParser{id: 527, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p520 = sequenceParser{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p516 = charParser{id: 516, chars: []rune{99}}
	var p517 = charParser{id: 517, chars: []rune{97}}
	var p518 = charParser{id: 518, chars: []rune{115}}
	var p519 = charParser{id: 519, chars: []rune{101}}
	p520.items = []parser{&p516, &p517, &p518, &p519}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p523.items = []parser{&p829, &p14}
	p524.items = []parser{&p829, &p14, &p523}
	var p515 = choiceParser{id: 515, commit: 66, name: "communication"}
	var p514 = sequenceParser{id: 514, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{515}}
	var p513 = sequenceParser{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p512 = sequenceParser{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p512.items = []parser{&p829, &p14}
	p513.items = []parser{&p829, &p14, &p512}
	p514.items = []parser{&p103, &p513, &p829, &p511}
	p515.options = []parser{&p500, &p511, &p514}
	var p526 = sequenceParser{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p525.items = []parser{&p829, &p14}
	p526.items = []parser{&p829, &p14, &p525}
	var p522 = sequenceParser{id: 522, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p521 = charParser{id: 521, chars: []rune{58}}
	p522.items = []parser{&p521}
	p527.items = []parser{&p520, &p524, &p829, &p515, &p526, &p829, &p522}
	var p531 = sequenceParser{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p529 = sequenceParser{id: 529, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p528 = charParser{id: 528, chars: []rune{59}}
	p529.items = []parser{&p528}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p530.items = []parser{&p829, &p529}
	p531.items = []parser{&p829, &p529, &p530}
	p532.items = []parser{&p527, &p531, &p829, &p797}
	p542.options = []parser{&p532, &p451}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p543 = choiceParser{id: 543, commit: 2}
	p543.options = []parser{&p532, &p451, &p797}
	p544.items = []parser{&p811, &p829, &p543}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p545.items = []parser{&p829, &p544}
	p546.items = []parser{&p829, &p544, &p545}
	p547.items = []parser{&p542, &p546}
	var p549 = sequenceParser{id: 549, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p548 = charParser{id: 548, chars: []rune{125}}
	p549.items = []parser{&p548}
	p552.items = []parser{&p539, &p551, &p829, &p541, &p829, &p811, &p829, &p547, &p829, &p811, &p829, &p549}
	var p593 = sequenceParser{id: 593, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{797}}
	var p582 = sequenceParser{id: 582, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p579 = charParser{id: 579, chars: []rune{102}}
	var p580 = charParser{id: 580, chars: []rune{111}}
	var p581 = charParser{id: 581, chars: []rune{114}}
	p582.items = []parser{&p579, &p580, &p581}
	var p592 = choiceParser{id: 592, commit: 2}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{592}}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p583 = sequenceParser{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p583.items = []parser{&p829, &p14}
	p584.items = []parser{&p14, &p583}
	var p578 = choiceParser{id: 578, commit: 66, name: "loop-expression"}
	var p577 = choiceParser{id: 577, commit: 64, name: "range-over-expression", generalizations: []int{578}}
	var p576 = sequenceParser{id: 576, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{577, 578}}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p572.items = []parser{&p829, &p14}
	p573.items = []parser{&p829, &p14, &p572}
	var p570 = sequenceParser{id: 570, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p568 = charParser{id: 568, chars: []rune{105}}
	var p569 = charParser{id: 569, chars: []rune{110}}
	p570.items = []parser{&p568, &p569}
	var p575 = sequenceParser{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p574 = sequenceParser{id: 574, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p574.items = []parser{&p829, &p14}
	p575.items = []parser{&p829, &p14, &p574}
	var p571 = choiceParser{id: 571, commit: 2}
	p571.options = []parser{&p396, &p225}
	p576.items = []parser{&p103, &p573, &p829, &p570, &p575, &p829, &p571}
	p577.options = []parser{&p576, &p225}
	p578.options = []parser{&p396, &p577}
	p585.items = []parser{&p584, &p829, &p578}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p586.items = []parser{&p829, &p14}
	p587.items = []parser{&p829, &p14, &p586}
	p588.items = []parser{&p585, &p587, &p829, &p190}
	var p591 = sequenceParser{id: 591, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{592}}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p589.items = []parser{&p829, &p14}
	p590.items = []parser{&p14, &p589}
	p591.items = []parser{&p590, &p829, &p190}
	p592.options = []parser{&p588, &p591}
	p593.items = []parser{&p582, &p829, &p592}
	var p741 = choiceParser{id: 741, commit: 66, name: "definition", generalizations: []int{797}}
	var p654 = sequenceParser{id: 654, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 797}}
	var p650 = sequenceParser{id: 650, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p647 = charParser{id: 647, chars: []rune{108}}
	var p648 = charParser{id: 648, chars: []rune{101}}
	var p649 = charParser{id: 649, chars: []rune{116}}
	p650.items = []parser{&p647, &p648, &p649}
	var p653 = sequenceParser{id: 653, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p652 = sequenceParser{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p652.items = []parser{&p829, &p14}
	p653.items = []parser{&p829, &p14, &p652}
	var p651 = choiceParser{id: 651, commit: 2}
	var p641 = sequenceParser{id: 641, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{651, 655, 656}}
	var p640 = sequenceParser{id: 640, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p637 = sequenceParser{id: 637, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p635.items = []parser{&p829, &p14}
	p636.items = []parser{&p14, &p635}
	var p634 = sequenceParser{id: 634, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p633 = charParser{id: 633, chars: []rune{61}}
	p634.items = []parser{&p633}
	p637.items = []parser{&p636, &p829, &p634}
	var p639 = sequenceParser{id: 639, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p638 = sequenceParser{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p638.items = []parser{&p829, &p14}
	p639.items = []parser{&p829, &p14, &p638}
	p640.items = []parser{&p103, &p829, &p637, &p639, &p829, &p396}
	p641.items = []parser{&p640}
	var p646 = sequenceParser{id: 646, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{651, 655, 656}}
	var p643 = sequenceParser{id: 643, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p642 = charParser{id: 642, chars: []rune{126}}
	p643.items = []parser{&p642}
	var p645 = sequenceParser{id: 645, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p644 = sequenceParser{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p644.items = []parser{&p829, &p14}
	p645.items = []parser{&p829, &p14, &p644}
	p646.items = []parser{&p643, &p645, &p829, &p640}
	p651.options = []parser{&p641, &p646}
	p654.items = []parser{&p650, &p653, &p829, &p651}
	var p675 = sequenceParser{id: 675, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 797}}
	var p668 = sequenceParser{id: 668, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p665 = charParser{id: 665, chars: []rune{108}}
	var p666 = charParser{id: 666, chars: []rune{101}}
	var p667 = charParser{id: 667, chars: []rune{116}}
	p668.items = []parser{&p665, &p666, &p667}
	var p674 = sequenceParser{id: 674, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p673 = sequenceParser{id: 673, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p673.items = []parser{&p829, &p14}
	p674.items = []parser{&p829, &p14, &p673}
	var p670 = sequenceParser{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p669 = charParser{id: 669, chars: []rune{40}}
	p670.items = []parser{&p669}
	var p660 = sequenceParser{id: 660, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p655 = choiceParser{id: 655, commit: 2}
	p655.options = []parser{&p641, &p646}
	var p659 = sequenceParser{id: 659, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p657 = sequenceParser{id: 657, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p656 = choiceParser{id: 656, commit: 2}
	p656.options = []parser{&p641, &p646}
	p657.items = []parser{&p113, &p829, &p656}
	var p658 = sequenceParser{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p658.items = []parser{&p829, &p657}
	p659.items = []parser{&p829, &p657, &p658}
	p660.items = []parser{&p655, &p659}
	var p672 = sequenceParser{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p671 = charParser{id: 671, chars: []rune{41}}
	p672.items = []parser{&p671}
	p675.items = []parser{&p668, &p674, &p829, &p670, &p829, &p113, &p829, &p660, &p829, &p113, &p829, &p672}
	var p690 = sequenceParser{id: 690, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 797}}
	var p679 = sequenceParser{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p676 = charParser{id: 676, chars: []rune{108}}
	var p677 = charParser{id: 677, chars: []rune{101}}
	var p678 = charParser{id: 678, chars: []rune{116}}
	p679.items = []parser{&p676, &p677, &p678}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p686.items = []parser{&p829, &p14}
	p687.items = []parser{&p829, &p14, &p686}
	var p681 = sequenceParser{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{126}}
	p681.items = []parser{&p680}
	var p689 = sequenceParser{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p688.items = []parser{&p829, &p14}
	p689.items = []parser{&p829, &p14, &p688}
	var p683 = sequenceParser{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p682 = charParser{id: 682, chars: []rune{40}}
	p683.items = []parser{&p682}
	var p664 = sequenceParser{id: 664, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p663 = sequenceParser{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p661 = sequenceParser{id: 661, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p661.items = []parser{&p113, &p829, &p641}
	var p662 = sequenceParser{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p662.items = []parser{&p829, &p661}
	p663.items = []parser{&p829, &p661, &p662}
	p664.items = []parser{&p641, &p663}
	var p685 = sequenceParser{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p684 = charParser{id: 684, chars: []rune{41}}
	p685.items = []parser{&p684}
	p690.items = []parser{&p679, &p687, &p829, &p681, &p689, &p829, &p683, &p829, &p113, &p829, &p664, &p829, &p113, &p829, &p685}
	var p706 = sequenceParser{id: 706, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 797}}
	var p702 = sequenceParser{id: 702, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p700 = charParser{id: 700, chars: []rune{102}}
	var p701 = charParser{id: 701, chars: []rune{110}}
	p702.items = []parser{&p700, &p701}
	var p705 = sequenceParser{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p704.items = []parser{&p829, &p14}
	p705.items = []parser{&p829, &p14, &p704}
	var p703 = choiceParser{id: 703, commit: 2}
	var p694 = sequenceParser{id: 694, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{703, 711, 712}}
	var p693 = sequenceParser{id: 693, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p692 = sequenceParser{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p691 = sequenceParser{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p691.items = []parser{&p829, &p14}
	p692.items = []parser{&p829, &p14, &p691}
	p693.items = []parser{&p103, &p692, &p829, &p200}
	p694.items = []parser{&p693}
	var p699 = sequenceParser{id: 699, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{703, 711, 712}}
	var p696 = sequenceParser{id: 696, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p695 = charParser{id: 695, chars: []rune{126}}
	p696.items = []parser{&p695}
	var p698 = sequenceParser{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p697 = sequenceParser{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p697.items = []parser{&p829, &p14}
	p698.items = []parser{&p829, &p14, &p697}
	p699.items = []parser{&p696, &p698, &p829, &p693}
	p703.options = []parser{&p694, &p699}
	p706.items = []parser{&p702, &p705, &p829, &p703}
	var p726 = sequenceParser{id: 726, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 797}}
	var p719 = sequenceParser{id: 719, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p717 = charParser{id: 717, chars: []rune{102}}
	var p718 = charParser{id: 718, chars: []rune{110}}
	p719.items = []parser{&p717, &p718}
	var p725 = sequenceParser{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p724 = sequenceParser{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p724.items = []parser{&p829, &p14}
	p725.items = []parser{&p829, &p14, &p724}
	var p721 = sequenceParser{id: 721, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p720 = charParser{id: 720, chars: []rune{40}}
	p721.items = []parser{&p720}
	var p716 = sequenceParser{id: 716, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p711 = choiceParser{id: 711, commit: 2}
	p711.options = []parser{&p694, &p699}
	var p715 = sequenceParser{id: 715, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p712 = choiceParser{id: 712, commit: 2}
	p712.options = []parser{&p694, &p699}
	p713.items = []parser{&p113, &p829, &p712}
	var p714 = sequenceParser{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p714.items = []parser{&p829, &p713}
	p715.items = []parser{&p829, &p713, &p714}
	p716.items = []parser{&p711, &p715}
	var p723 = sequenceParser{id: 723, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p722 = charParser{id: 722, chars: []rune{41}}
	p723.items = []parser{&p722}
	p726.items = []parser{&p719, &p725, &p829, &p721, &p829, &p113, &p829, &p716, &p829, &p113, &p829, &p723}
	var p740 = sequenceParser{id: 740, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 797}}
	var p729 = sequenceParser{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p727 = charParser{id: 727, chars: []rune{102}}
	var p728 = charParser{id: 728, chars: []rune{110}}
	p729.items = []parser{&p727, &p728}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p736.items = []parser{&p829, &p14}
	p737.items = []parser{&p829, &p14, &p736}
	var p731 = sequenceParser{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p730 = charParser{id: 730, chars: []rune{126}}
	p731.items = []parser{&p730}
	var p739 = sequenceParser{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p738 = sequenceParser{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p738.items = []parser{&p829, &p14}
	p739.items = []parser{&p829, &p14, &p738}
	var p733 = sequenceParser{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p732 = charParser{id: 732, chars: []rune{40}}
	p733.items = []parser{&p732}
	var p710 = sequenceParser{id: 710, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p709 = sequenceParser{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p707 = sequenceParser{id: 707, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p707.items = []parser{&p113, &p829, &p694}
	var p708 = sequenceParser{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p708.items = []parser{&p829, &p707}
	p709.items = []parser{&p829, &p707, &p708}
	p710.items = []parser{&p694, &p709}
	var p735 = sequenceParser{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p734 = charParser{id: 734, chars: []rune{41}}
	p735.items = []parser{&p734}
	p740.items = []parser{&p729, &p737, &p829, &p731, &p739, &p829, &p733, &p829, &p113, &p829, &p710, &p829, &p113, &p829, &p735}
	p741.options = []parser{&p654, &p675, &p690, &p706, &p726, &p740}
	var p776 = choiceParser{id: 776, commit: 64, name: "use", generalizations: []int{797}}
	var p764 = sequenceParser{id: 764, commit: 66, name: "use-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{776, 797}}
	var p761 = sequenceParser{id: 761, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p758 = charParser{id: 758, chars: []rune{117}}
	var p759 = charParser{id: 759, chars: []rune{115}}
	var p760 = charParser{id: 760, chars: []rune{101}}
	p761.items = []parser{&p758, &p759, &p760}
	var p763 = sequenceParser{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p762 = sequenceParser{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p762.items = []parser{&p829, &p14}
	p763.items = []parser{&p829, &p14, &p762}
	var p753 = choiceParser{id: 753, commit: 64, name: "use-fact"}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{753}}
	var p744 = choiceParser{id: 744, commit: 2}
	var p743 = sequenceParser{id: 743, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{744}}
	var p742 = charParser{id: 742, chars: []rune{46}}
	p743.items = []parser{&p742}
	p744.options = []parser{&p103, &p743}
	var p749 = sequenceParser{id: 749, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p748 = sequenceParser{id: 748, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p747 = sequenceParser{id: 747, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p747.items = []parser{&p829, &p14}
	p748.items = []parser{&p14, &p747}
	var p746 = sequenceParser{id: 746, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p745 = charParser{id: 745, chars: []rune{61}}
	p746.items = []parser{&p745}
	p749.items = []parser{&p748, &p829, &p746}
	var p751 = sequenceParser{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p750 = sequenceParser{id: 750, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p750.items = []parser{&p829, &p14}
	p751.items = []parser{&p829, &p14, &p750}
	p752.items = []parser{&p744, &p829, &p749, &p751, &p829, &p86}
	p753.options = []parser{&p86, &p752}
	p764.items = []parser{&p761, &p763, &p829, &p753}
	var p775 = sequenceParser{id: 775, commit: 66, name: "use-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{776, 797}}
	var p768 = sequenceParser{id: 768, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p765 = charParser{id: 765, chars: []rune{117}}
	var p766 = charParser{id: 766, chars: []rune{115}}
	var p767 = charParser{id: 767, chars: []rune{101}}
	p768.items = []parser{&p765, &p766, &p767}
	var p774 = sequenceParser{id: 774, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p773 = sequenceParser{id: 773, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p773.items = []parser{&p829, &p14}
	p774.items = []parser{&p829, &p14, &p773}
	var p770 = sequenceParser{id: 770, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p769 = charParser{id: 769, chars: []rune{40}}
	p770.items = []parser{&p769}
	var p757 = sequenceParser{id: 757, commit: 66, name: "use-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p754.items = []parser{&p113, &p829, &p753}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p755.items = []parser{&p829, &p754}
	p756.items = []parser{&p829, &p754, &p755}
	p757.items = []parser{&p753, &p756}
	var p772 = sequenceParser{id: 772, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p771 = charParser{id: 771, chars: []rune{41}}
	p772.items = []parser{&p771}
	p775.items = []parser{&p768, &p774, &p829, &p770, &p829, &p113, &p829, &p757, &p829, &p113, &p829, &p772}
	p776.options = []parser{&p764, &p775}
	var p786 = sequenceParser{id: 786, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{797}}
	var p783 = sequenceParser{id: 783, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p777 = charParser{id: 777, chars: []rune{101}}
	var p778 = charParser{id: 778, chars: []rune{120}}
	var p779 = charParser{id: 779, chars: []rune{112}}
	var p780 = charParser{id: 780, chars: []rune{111}}
	var p781 = charParser{id: 781, chars: []rune{114}}
	var p782 = charParser{id: 782, chars: []rune{116}}
	p783.items = []parser{&p777, &p778, &p779, &p780, &p781, &p782}
	var p785 = sequenceParser{id: 785, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p784 = sequenceParser{id: 784, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p784.items = []parser{&p829, &p14}
	p785.items = []parser{&p829, &p14, &p784}
	p786.items = []parser{&p783, &p785, &p829, &p741}
	var p806 = sequenceParser{id: 806, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{797}}
	var p799 = sequenceParser{id: 799, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p798 = charParser{id: 798, chars: []rune{40}}
	p799.items = []parser{&p798}
	var p803 = sequenceParser{id: 803, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p802 = sequenceParser{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p802.items = []parser{&p829, &p14}
	p803.items = []parser{&p829, &p14, &p802}
	var p805 = sequenceParser{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p804 = sequenceParser{id: 804, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p804.items = []parser{&p829, &p14}
	p805.items = []parser{&p829, &p14, &p804}
	var p801 = sequenceParser{id: 801, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p800 = charParser{id: 800, chars: []rune{41}}
	p801.items = []parser{&p800}
	p806.items = []parser{&p799, &p803, &p829, &p797, &p805, &p829, &p801}
	p797.options = []parser{&p185, &p433, &p490, &p552, &p593, &p741, &p776, &p786, &p806, &p787}
	var p814 = sequenceParser{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p812 = sequenceParser{id: 812, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p812.items = []parser{&p811, &p829, &p797}
	var p813 = sequenceParser{id: 813, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p813.items = []parser{&p829, &p812}
	p814.items = []parser{&p829, &p812, &p813}
	p815.items = []parser{&p797, &p814}
	p830.items = []parser{&p826, &p829, &p811, &p829, &p815, &p829, &p811}
	p831.items = []parser{&p829, &p830, &p829}
	var b831 = sequenceBuilder{id: 831, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b829 = choiceBuilder{id: 829, commit: 2}
	var b827 = choiceBuilder{id: 827, commit: 70}
	var b2 = sequenceBuilder{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b1 = charBuilder{}
	b2.items = []builder{&b1}
	var b4 = sequenceBuilder{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b3 = charBuilder{}
	b4.items = []builder{&b3}
	var b6 = sequenceBuilder{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b5 = charBuilder{}
	b6.items = []builder{&b5}
	var b8 = sequenceBuilder{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b7 = charBuilder{}
	b8.items = []builder{&b7}
	var b10 = sequenceBuilder{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b9 = charBuilder{}
	b10.items = []builder{&b9}
	var b12 = sequenceBuilder{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b11 = charBuilder{}
	b12.items = []builder{&b11}
	b827.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b828 = sequenceBuilder{id: 828, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
	var b42 = sequenceBuilder{id: 42, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b38 = choiceBuilder{id: 38, commit: 66}
	var b21 = sequenceBuilder{id: 21, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b20 = sequenceBuilder{id: 20, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b18 = charBuilder{}
	var b19 = charBuilder{}
	b20.items = []builder{&b18, &b19}
	var b17 = sequenceBuilder{id: 17, commit: 72, name: "line-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var b16 = sequenceBuilder{id: 16, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b15 = charBuilder{}
	b16.items = []builder{&b15}
	b17.items = []builder{&b16}
	b21.items = []builder{&b20, &b17}
	var b37 = sequenceBuilder{id: 37, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b33 = sequenceBuilder{id: 33, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b31 = charBuilder{}
	var b32 = charBuilder{}
	b33.items = []builder{&b31, &b32}
	var b30 = sequenceBuilder{id: 30, commit: 72, name: "block-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var b29 = choiceBuilder{id: 29, commit: 10}
	var b23 = sequenceBuilder{id: 23, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b22 = charBuilder{}
	b23.items = []builder{&b22}
	var b28 = sequenceBuilder{id: 28, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b25 = sequenceBuilder{id: 25, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b24 = charBuilder{}
	b25.items = []builder{&b24}
	var b27 = sequenceBuilder{id: 27, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b26 = charBuilder{}
	b27.items = []builder{&b26}
	b28.items = []builder{&b25, &b27}
	b29.options = []builder{&b23, &b28}
	b30.items = []builder{&b29}
	var b36 = sequenceBuilder{id: 36, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b34 = charBuilder{}
	var b35 = charBuilder{}
	b36.items = []builder{&b34, &b35}
	b37.items = []builder{&b33, &b30, &b36}
	b38.options = []builder{&b21, &b37}
	var b41 = sequenceBuilder{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b39 = sequenceBuilder{id: 39, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b14 = sequenceBuilder{id: 14, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b39.items = []builder{&b14, &b829, &b38}
	var b40 = sequenceBuilder{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b40.items = []builder{&b829, &b39}
	b41.items = []builder{&b829, &b39, &b40}
	b42.items = []builder{&b38, &b41}
	b828.items = []builder{&b42}
	b829.options = []builder{&b827, &b828}
	var b830 = sequenceBuilder{id: 830, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b826 = sequenceBuilder{id: 826, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b823 = sequenceBuilder{id: 823, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b821 = charBuilder{}
	var b822 = charBuilder{}
	b823.items = []builder{&b821, &b822}
	var b820 = sequenceBuilder{id: 820, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b819 = sequenceBuilder{id: 819, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b817 = sequenceBuilder{id: 817, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b816 = charBuilder{}
	b817.items = []builder{&b816}
	var b818 = sequenceBuilder{id: 818, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b818.items = []builder{&b829, &b817}
	b819.items = []builder{&b817, &b818}
	b820.items = []builder{&b819}
	var b825 = sequenceBuilder{id: 825, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b824 = charBuilder{}
	b825.items = []builder{&b824}
	b826.items = []builder{&b823, &b829, &b820, &b829, &b825}
	var b811 = sequenceBuilder{id: 811, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b809 = choiceBuilder{id: 809, commit: 2}
	var b808 = sequenceBuilder{id: 808, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b807 = charBuilder{}
	b808.items = []builder{&b807}
	b809.options = []builder{&b808, &b14}
	var b810 = sequenceBuilder{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b810.items = []builder{&b829, &b809}
	b811.items = []builder{&b809, &b810}
	var b815 = sequenceBuilder{id: 815, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b797 = choiceBuilder{id: 797, commit: 66}
	var b185 = sequenceBuilder{id: 185, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b181 = sequenceBuilder{id: 181, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b175 = charBuilder{}
	var b176 = charBuilder{}
	var b177 = charBuilder{}
	var b178 = charBuilder{}
	var b179 = charBuilder{}
	var b180 = charBuilder{}
	b181.items = []builder{&b175, &b176, &b177, &b178, &b179, &b180}
	var b184 = sequenceBuilder{id: 184, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b183 = sequenceBuilder{id: 183, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b182 = sequenceBuilder{id: 182, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b182.items = []builder{&b829, &b14}
	b183.items = []builder{&b14, &b182}
	var b396 = choiceBuilder{id: 396, commit: 66}
	var b267 = choiceBuilder{id: 267, commit: 66}
	var b247 = choiceBuilder{id: 247, commit: 66}
	var b60 = choiceBuilder{id: 60, commit: 64, name: "int"}
	var b51 = sequenceBuilder{id: 51, commit: 74, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b50 = sequenceBuilder{id: 50, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b49 = charBuilder{}
	b50.items = []builder{&b49}
	var b44 = sequenceBuilder{id: 44, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b43 = charBuilder{}
	b44.items = []builder{&b43}
	b51.items = []builder{&b50, &b44}
	var b54 = sequenceBuilder{id: 54, commit: 74, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b53 = sequenceBuilder{id: 53, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b52 = charBuilder{}
	b53.items = []builder{&b52}
	var b46 = sequenceBuilder{id: 46, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b45 = charBuilder{}
	b46.items = []builder{&b45}
	b54.items = []builder{&b53, &b46}
	var b59 = sequenceBuilder{id: 59, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}}
	var b56 = sequenceBuilder{id: 56, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b55 = charBuilder{}
	b56.items = []builder{&b55}
	var b58 = sequenceBuilder{id: 58, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b57 = charBuilder{}
	b58.items = []builder{&b57}
	var b48 = sequenceBuilder{id: 48, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b47 = charBuilder{}
	b48.items = []builder{&b47}
	b59.items = []builder{&b56, &b58, &b48}
	b60.options = []builder{&b51, &b54, &b59}
	var b73 = choiceBuilder{id: 73, commit: 72, name: "float"}
	var b68 = sequenceBuilder{id: 68, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}}
	var b67 = sequenceBuilder{id: 67, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b66 = charBuilder{}
	b67.items = []builder{&b66}
	var b65 = sequenceBuilder{id: 65, commit: 74, ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var b62 = sequenceBuilder{id: 62, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b61 = charBuilder{}
	b62.items = []builder{&b61}
	var b64 = sequenceBuilder{id: 64, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b63 = charBuilder{}
	b64.items = []builder{&b63}
	b65.items = []builder{&b62, &b64, &b44}
	b68.items = []builder{&b44, &b67, &b44, &b65}
	var b71 = sequenceBuilder{id: 71, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}}
	var b70 = sequenceBuilder{id: 70, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b69 = charBuilder{}
	b70.items = []builder{&b69}
	b71.items = []builder{&b70, &b44, &b65}
	var b72 = sequenceBuilder{id: 72, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}}
	b72.items = []builder{&b44, &b65}
	b73.options = []builder{&b68, &b71, &b72}
	var b86 = sequenceBuilder{id: 86, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}}
	var b75 = sequenceBuilder{id: 75, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b74 = charBuilder{}
	b75.items = []builder{&b74}
	var b83 = choiceBuilder{id: 83, commit: 10}
	var b77 = sequenceBuilder{id: 77, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b76 = charBuilder{}
	b77.items = []builder{&b76}
	var b82 = sequenceBuilder{id: 82, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b79 = sequenceBuilder{id: 79, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b78 = charBuilder{}
	b79.items = []builder{&b78}
	var b81 = sequenceBuilder{id: 81, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b80 = charBuilder{}
	b81.items = []builder{&b80}
	b82.items = []builder{&b79, &b81}
	b83.options = []builder{&b77, &b82}
	var b85 = sequenceBuilder{id: 85, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b84 = charBuilder{}
	b85.items = []builder{&b84}
	b86.items = []builder{&b75, &b83, &b85}
	var b98 = choiceBuilder{id: 98, commit: 66}
	var b91 = sequenceBuilder{id: 91, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b87 = charBuilder{}
	var b88 = charBuilder{}
	var b89 = charBuilder{}
	var b90 = charBuilder{}
	b91.items = []builder{&b87, &b88, &b89, &b90}
	var b97 = sequenceBuilder{id: 97, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b92 = charBuilder{}
	var b93 = charBuilder{}
	var b94 = charBuilder{}
	var b95 = charBuilder{}
	var b96 = charBuilder{}
	b97.items = []builder{&b92, &b93, &b94, &b95, &b96}
	b98.options = []builder{&b91, &b97}
	var b511 = sequenceBuilder{id: 511, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b508 = sequenceBuilder{id: 508, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b501 = charBuilder{}
	var b502 = charBuilder{}
	var b503 = charBuilder{}
	var b504 = charBuilder{}
	var b505 = charBuilder{}
	var b506 = charBuilder{}
	var b507 = charBuilder{}
	b508.items = []builder{&b501, &b502, &b503, &b504, &b505, &b506, &b507}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b509.items = []builder{&b829, &b14}
	b510.items = []builder{&b829, &b14, &b509}
	b511.items = []builder{&b508, &b510, &b829, &b267}
	var b103 = sequenceBuilder{id: 103, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b100 = sequenceBuilder{id: 100, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b99 = charBuilder{}
	b100.items = []builder{&b99}
	var b102 = sequenceBuilder{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b101 = charBuilder{}
	b102.items = []builder{&b101}
	b103.items = []builder{&b100, &b102}
	var b124 = sequenceBuilder{id: 124, commit: 64, name: "list", ranges: [][]int{{1, 1}}}
	var b123 = sequenceBuilder{id: 123, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b120 = sequenceBuilder{id: 120, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b119 = charBuilder{}
	b120.items = []builder{&b119}
	var b113 = sequenceBuilder{id: 113, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b111 = choiceBuilder{id: 111, commit: 2}
	var b110 = sequenceBuilder{id: 110, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b109 = charBuilder{}
	b110.items = []builder{&b109}
	b111.options = []builder{&b14, &b110}
	var b112 = sequenceBuilder{id: 112, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b112.items = []builder{&b829, &b111}
	b113.items = []builder{&b111, &b112}
	var b118 = sequenceBuilder{id: 118, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b114 = choiceBuilder{id: 114, commit: 66}
	var b108 = sequenceBuilder{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b107 = sequenceBuilder{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b104 = charBuilder{}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	b107.items = []builder{&b104, &b105, &b106}
	b108.items = []builder{&b267, &b829, &b107}
	b114.options = []builder{&b396, &b108}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b115 = sequenceBuilder{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b115.items = []builder{&b113, &b829, &b114}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b116.items = []builder{&b829, &b115}
	b117.items = []builder{&b829, &b115, &b116}
	b118.items = []builder{&b114, &b117}
	var b122 = sequenceBuilder{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b121 = charBuilder{}
	b122.items = []builder{&b121}
	b123.items = []builder{&b120, &b829, &b113, &b829, &b118, &b829, &b113, &b829, &b122}
	b124.items = []builder{&b123}
	var b129 = sequenceBuilder{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b126 = sequenceBuilder{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b125 = charBuilder{}
	b126.items = []builder{&b125}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b127 = sequenceBuilder{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b127.items = []builder{&b829, &b14}
	b128.items = []builder{&b829, &b14, &b127}
	b129.items = []builder{&b126, &b128, &b829, &b123}
	var b158 = sequenceBuilder{id: 158, commit: 64, name: "struct", ranges: [][]int{{1, 1}}}
	var b157 = sequenceBuilder{id: 157, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b154 = sequenceBuilder{id: 154, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b153 = charBuilder{}
	b154.items = []builder{&b153}
	var b152 = sequenceBuilder{id: 152, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b147 = choiceBuilder{id: 147, commit: 2}
	var b146 = sequenceBuilder{id: 146, commit: 64, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b139 = choiceBuilder{id: 139, commit: 2}
	var b138 = sequenceBuilder{id: 138, commit: 64, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b131 = sequenceBuilder{id: 131, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b130 = charBuilder{}
	b131.items = []builder{&b130}
	var b135 = sequenceBuilder{id: 135, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b134 = sequenceBuilder{id: 134, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b134.items = []builder{&b829, &b14}
	b135.items = []builder{&b829, &b14, &b134}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b136 = sequenceBuilder{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b136.items = []builder{&b829, &b14}
	b137.items = []builder{&b829, &b14, &b136}
	var b133 = sequenceBuilder{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b132 = charBuilder{}
	b133.items = []builder{&b132}
	b138.items = []builder{&b131, &b135, &b829, &b396, &b137, &b829, &b133}
	b139.options = []builder{&b103, &b86, &b138}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b142 = sequenceBuilder{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b142.items = []builder{&b829, &b14}
	b143.items = []builder{&b829, &b14, &b142}
	var b141 = sequenceBuilder{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b140 = charBuilder{}
	b141.items = []builder{&b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b829, &b14}
	b145.items = []builder{&b829, &b14, &b144}
	b146.items = []builder{&b139, &b143, &b829, &b141, &b145, &b829, &b396}
	b147.options = []builder{&b146, &b108}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b149 = sequenceBuilder{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b148 = choiceBuilder{id: 148, commit: 2}
	b148.options = []builder{&b146, &b108}
	b149.items = []builder{&b113, &b829, &b148}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b150.items = []builder{&b829, &b149}
	b151.items = []builder{&b829, &b149, &b150}
	b152.items = []builder{&b147, &b151}
	var b156 = sequenceBuilder{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b155 = charBuilder{}
	b156.items = []builder{&b155}
	b157.items = []builder{&b154, &b829, &b113, &b829, &b152, &b829, &b113, &b829, &b156}
	b158.items = []builder{&b157}
	var b163 = sequenceBuilder{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b160 = sequenceBuilder{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b159 = charBuilder{}
	b160.items = []builder{&b159}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b161 = sequenceBuilder{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b161.items = []builder{&b829, &b14}
	b162.items = []builder{&b829, &b14, &b161}
	b163.items = []builder{&b160, &b162, &b829, &b157}
	var b206 = sequenceBuilder{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b203 = sequenceBuilder{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b201 = charBuilder{}
	var b202 = charBuilder{}
	b203.items = []builder{&b201, &b202}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b204 = sequenceBuilder{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b204.items = []builder{&b829, &b14}
	b205.items = []builder{&b829, &b14, &b204}
	var b200 = sequenceBuilder{id: 200, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b192 = sequenceBuilder{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b191 = charBuilder{}
	b192.items = []builder{&b191}
	var b194 = choiceBuilder{id: 194, commit: 2}
	var b167 = sequenceBuilder{id: 167, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b164.items = []builder{&b113, &b829, &b103}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b165.items = []builder{&b829, &b164}
	b166.items = []builder{&b829, &b164, &b165}
	b167.items = []builder{&b103, &b166}
	var b193 = sequenceBuilder{id: 193, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b174 = sequenceBuilder{id: 174, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b171 = sequenceBuilder{id: 171, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b168 = charBuilder{}
	var b169 = charBuilder{}
	var b170 = charBuilder{}
	b171.items = []builder{&b168, &b169, &b170}
	var b173 = sequenceBuilder{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b172 = sequenceBuilder{id: 172, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b172.items = []builder{&b829, &b14}
	b173.items = []builder{&b829, &b14, &b172}
	b174.items = []builder{&b171, &b173, &b829, &b103}
	b193.items = []builder{&b167, &b829, &b113, &b829, &b174}
	b194.options = []builder{&b167, &b193, &b174}
	var b196 = sequenceBuilder{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b195 = charBuilder{}
	b196.items = []builder{&b195}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b198 = sequenceBuilder{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b198.items = []builder{&b829, &b14}
	b199.items = []builder{&b829, &b14, &b198}
	var b197 = choiceBuilder{id: 197, commit: 2}
	var b787 = choiceBuilder{id: 787, commit: 66}
	var b500 = sequenceBuilder{id: 500, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b495 = sequenceBuilder{id: 495, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b491 = charBuilder{}
	var b492 = charBuilder{}
	var b493 = charBuilder{}
	var b494 = charBuilder{}
	b495.items = []builder{&b491, &b492, &b493, &b494}
	var b497 = sequenceBuilder{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b496 = sequenceBuilder{id: 496, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b496.items = []builder{&b829, &b14}
	b497.items = []builder{&b829, &b14, &b496}
	var b499 = sequenceBuilder{id: 499, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b498 = sequenceBuilder{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b498.items = []builder{&b829, &b14}
	b499.items = []builder{&b829, &b14, &b498}
	b500.items = []builder{&b495, &b497, &b829, &b267, &b499, &b829, &b267}
	var b558 = sequenceBuilder{id: 558, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b555 = sequenceBuilder{id: 555, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b553 = charBuilder{}
	var b554 = charBuilder{}
	b555.items = []builder{&b553, &b554}
	var b557 = sequenceBuilder{id: 557, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b556 = sequenceBuilder{id: 556, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b556.items = []builder{&b829, &b14}
	b557.items = []builder{&b829, &b14, &b556}
	var b237 = sequenceBuilder{id: 237, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b234 = sequenceBuilder{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b233 = charBuilder{}
	b234.items = []builder{&b233}
	var b236 = sequenceBuilder{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b235 = charBuilder{}
	b236.items = []builder{&b235}
	b237.items = []builder{&b267, &b829, &b234, &b829, &b113, &b829, &b118, &b829, &b113, &b829, &b236}
	b558.items = []builder{&b555, &b557, &b829, &b237}
	var b567 = sequenceBuilder{id: 567, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b564 = sequenceBuilder{id: 564, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b559 = charBuilder{}
	var b560 = charBuilder{}
	var b561 = charBuilder{}
	var b562 = charBuilder{}
	var b563 = charBuilder{}
	b564.items = []builder{&b559, &b560, &b561, &b562, &b563}
	var b566 = sequenceBuilder{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b565.items = []builder{&b829, &b14}
	b566.items = []builder{&b829, &b14, &b565}
	b567.items = []builder{&b564, &b566, &b829, &b237}
	var b632 = choiceBuilder{id: 632, commit: 64, name: "assignment"}
	var b612 = sequenceBuilder{id: 612, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b609 = sequenceBuilder{id: 609, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b606 = charBuilder{}
	var b607 = charBuilder{}
	var b608 = charBuilder{}
	b609.items = []builder{&b606, &b607, &b608}
	var b611 = sequenceBuilder{id: 611, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b610 = sequenceBuilder{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b610.items = []builder{&b829, &b14}
	b611.items = []builder{&b829, &b14, &b610}
	var b601 = sequenceBuilder{id: 601, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b596.items = []builder{&b829, &b14}
	b597.items = []builder{&b14, &b596}
	var b595 = sequenceBuilder{id: 595, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b594 = charBuilder{}
	b595.items = []builder{&b594}
	b598.items = []builder{&b597, &b829, &b595}
	var b600 = sequenceBuilder{id: 600, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b599 = sequenceBuilder{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b599.items = []builder{&b829, &b14}
	b600.items = []builder{&b829, &b14, &b599}
	b601.items = []builder{&b267, &b829, &b598, &b600, &b829, &b396}
	b612.items = []builder{&b609, &b611, &b829, &b601}
	var b619 = sequenceBuilder{id: 619, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b615.items = []builder{&b829, &b14}
	b616.items = []builder{&b829, &b14, &b615}
	var b614 = sequenceBuilder{id: 614, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b613 = charBuilder{}
	b614.items = []builder{&b613}
	var b618 = sequenceBuilder{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b617.items = []builder{&b829, &b14}
	b618.items = []builder{&b829, &b14, &b617}
	b619.items = []builder{&b267, &b616, &b829, &b614, &b618, &b829, &b396}
	var b631 = sequenceBuilder{id: 631, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b623 = sequenceBuilder{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b620 = charBuilder{}
	var b621 = charBuilder{}
	var b622 = charBuilder{}
	b623.items = []builder{&b620, &b621, &b622}
	var b630 = sequenceBuilder{id: 630, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b629 = sequenceBuilder{id: 629, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b629.items = []builder{&b829, &b14}
	b630.items = []builder{&b829, &b14, &b629}
	var b625 = sequenceBuilder{id: 625, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b624 = charBuilder{}
	b625.items = []builder{&b624}
	var b626 = sequenceBuilder{id: 626, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b605 = sequenceBuilder{id: 605, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b604 = sequenceBuilder{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b602 = sequenceBuilder{id: 602, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b602.items = []builder{&b113, &b829, &b601}
	var b603 = sequenceBuilder{id: 603, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b603.items = []builder{&b829, &b602}
	b604.items = []builder{&b829, &b602, &b603}
	b605.items = []builder{&b601, &b604}
	b626.items = []builder{&b113, &b829, &b605}
	var b628 = sequenceBuilder{id: 628, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b627 = charBuilder{}
	b628.items = []builder{&b627}
	b631.items = []builder{&b623, &b630, &b829, &b625, &b829, &b626, &b829, &b113, &b829, &b628}
	b632.options = []builder{&b612, &b619, &b631}
	var b796 = sequenceBuilder{id: 796, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b789 = sequenceBuilder{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b788 = charBuilder{}
	b789.items = []builder{&b788}
	var b793 = sequenceBuilder{id: 793, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b792.items = []builder{&b829, &b14}
	b793.items = []builder{&b829, &b14, &b792}
	var b795 = sequenceBuilder{id: 795, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b794 = sequenceBuilder{id: 794, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b794.items = []builder{&b829, &b14}
	b795.items = []builder{&b829, &b14, &b794}
	var b791 = sequenceBuilder{id: 791, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b790 = charBuilder{}
	b791.items = []builder{&b790}
	b796.items = []builder{&b789, &b793, &b829, &b787, &b795, &b829, &b791}
	b787.options = []builder{&b500, &b558, &b567, &b632, &b796, &b396}
	var b190 = sequenceBuilder{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b187 = sequenceBuilder{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b186 = charBuilder{}
	b187.items = []builder{&b186}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	b190.items = []builder{&b187, &b829, &b811, &b829, &b815, &b829, &b811, &b829, &b189}
	b197.options = []builder{&b787, &b190}
	b200.items = []builder{&b192, &b829, &b113, &b829, &b194, &b829, &b113, &b829, &b196, &b199, &b829, &b197}
	b206.items = []builder{&b203, &b205, &b829, &b200}
	var b216 = sequenceBuilder{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b209 = sequenceBuilder{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b207 = charBuilder{}
	var b208 = charBuilder{}
	b209.items = []builder{&b207, &b208}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b212 = sequenceBuilder{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b212.items = []builder{&b829, &b14}
	b213.items = []builder{&b829, &b14, &b212}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b210 = charBuilder{}
	b211.items = []builder{&b210}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b214.items = []builder{&b829, &b14}
	b215.items = []builder{&b829, &b14, &b214}
	b216.items = []builder{&b209, &b213, &b829, &b211, &b215, &b829, &b200}
	var b232 = sequenceBuilder{id: 232, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b229 = sequenceBuilder{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b228 = sequenceBuilder{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b228.items = []builder{&b829, &b14}
	b229.items = []builder{&b829, &b14, &b228}
	var b227 = sequenceBuilder{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b226 = charBuilder{}
	b227.items = []builder{&b226}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b230 = sequenceBuilder{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b230.items = []builder{&b829, &b14}
	b231.items = []builder{&b829, &b14, &b230}
	b232.items = []builder{&b267, &b229, &b829, &b227, &b231, &b829, &b103}
	var b246 = sequenceBuilder{id: 246, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b239 = sequenceBuilder{id: 239, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b238 = charBuilder{}
	b239.items = []builder{&b238}
	var b243 = sequenceBuilder{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b242 = sequenceBuilder{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b242.items = []builder{&b829, &b14}
	b243.items = []builder{&b829, &b14, &b242}
	var b245 = sequenceBuilder{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b244 = sequenceBuilder{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b244.items = []builder{&b829, &b14}
	b245.items = []builder{&b829, &b14, &b244}
	var b241 = sequenceBuilder{id: 241, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b240 = charBuilder{}
	b241.items = []builder{&b240}
	b246.items = []builder{&b239, &b243, &b829, &b396, &b245, &b829, &b241}
	b247.options = []builder{&b60, &b73, &b86, &b98, &b511, &b103, &b124, &b129, &b158, &b163, &b206, &b216, &b232, &b237, &b246}
	var b266 = choiceBuilder{id: 266, commit: 64, name: "expression-indexer"}
	var b256 = sequenceBuilder{id: 256, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b249 = sequenceBuilder{id: 249, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b248 = charBuilder{}
	b249.items = []builder{&b248}
	var b253 = sequenceBuilder{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b252 = sequenceBuilder{id: 252, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b252.items = []builder{&b829, &b14}
	b253.items = []builder{&b829, &b14, &b252}
	var b255 = sequenceBuilder{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b254 = sequenceBuilder{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b254.items = []builder{&b829, &b14}
	b255.items = []builder{&b829, &b14, &b254}
	var b251 = sequenceBuilder{id: 251, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b250 = charBuilder{}
	b251.items = []builder{&b250}
	b256.items = []builder{&b247, &b829, &b249, &b253, &b829, &b396, &b255, &b829, &b251}
	var b265 = sequenceBuilder{id: 265, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b258 = sequenceBuilder{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b257 = charBuilder{}
	b258.items = []builder{&b257}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b261 = sequenceBuilder{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b261.items = []builder{&b829, &b14}
	b262.items = []builder{&b829, &b14, &b261}
	var b225 = sequenceBuilder{id: 225, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b217.items = []builder{&b396}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b221 = sequenceBuilder{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b221.items = []builder{&b829, &b14}
	b222.items = []builder{&b829, &b14, &b221}
	var b220 = sequenceBuilder{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b219 = charBuilder{}
	b220.items = []builder{&b219}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b829, &b14}
	b224.items = []builder{&b829, &b14, &b223}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b396}
	b225.items = []builder{&b217, &b222, &b829, &b220, &b224, &b829, &b218}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b263.items = []builder{&b829, &b14}
	b264.items = []builder{&b829, &b14, &b263}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	b265.items = []builder{&b247, &b829, &b258, &b262, &b829, &b225, &b264, &b829, &b260}
	b266.options = []builder{&b256, &b265}
	b267.options = []builder{&b247, &b266}
	var b327 = sequenceBuilder{id: 327, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b326 = choiceBuilder{id: 326, commit: 66}
	var b286 = sequenceBuilder{id: 286, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b285 = charBuilder{}
	b286.items = []builder{&b285}
	var b288 = sequenceBuilder{id: 288, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b287 = charBuilder{}
	b288.items = []builder{&b287}
	var b269 = sequenceBuilder{id: 269, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b268 = charBuilder{}
	b269.items = []builder{&b268}
	var b300 = sequenceBuilder{id: 300, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b299 = charBuilder{}
	b300.items = []builder{&b299}
	b326.options = []builder{&b286, &b288, &b269, &b300}
	b327.items = []builder{&b326, &b829, &b267}
	var b374 = choiceBuilder{id: 374, commit: 66}
	var b345 = sequenceBuilder{id: 345, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b333 = choiceBuilder{id: 333, commit: 66}
	b333.options = []builder{&b267, &b327}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b340 = sequenceBuilder{id: 340, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b339 = sequenceBuilder{id: 339, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b339.items = []builder{&b829, &b14}
	b340.items = []builder{&b14, &b339}
	var b328 = choiceBuilder{id: 328, commit: 66}
	var b271 = sequenceBuilder{id: 271, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b270 = charBuilder{}
	b271.items = []builder{&b270}
	var b278 = sequenceBuilder{id: 278, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b276 = charBuilder{}
	var b277 = charBuilder{}
	b278.items = []builder{&b276, &b277}
	var b281 = sequenceBuilder{id: 281, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b279 = charBuilder{}
	var b280 = charBuilder{}
	b281.items = []builder{&b279, &b280}
	var b284 = sequenceBuilder{id: 284, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b282 = charBuilder{}
	var b283 = charBuilder{}
	b284.items = []builder{&b282, &b283}
	var b290 = sequenceBuilder{id: 290, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b289 = charBuilder{}
	b290.items = []builder{&b289}
	var b292 = sequenceBuilder{id: 292, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b291 = charBuilder{}
	b292.items = []builder{&b291}
	var b294 = sequenceBuilder{id: 294, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b293 = charBuilder{}
	b294.items = []builder{&b293}
	b328.options = []builder{&b271, &b278, &b281, &b284, &b290, &b292, &b294}
	var b342 = sequenceBuilder{id: 342, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b341 = sequenceBuilder{id: 341, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b341.items = []builder{&b829, &b14}
	b342.items = []builder{&b829, &b14, &b341}
	b343.items = []builder{&b340, &b829, &b328, &b342, &b829, &b333}
	var b344 = sequenceBuilder{id: 344, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b344.items = []builder{&b829, &b343}
	b345.items = []builder{&b333, &b829, &b343, &b344}
	var b352 = sequenceBuilder{id: 352, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b334 = choiceBuilder{id: 334, commit: 66}
	b334.options = []builder{&b333, &b345}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b346.items = []builder{&b829, &b14}
	b347.items = []builder{&b14, &b346}
	var b329 = choiceBuilder{id: 329, commit: 66}
	var b273 = sequenceBuilder{id: 273, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b272 = charBuilder{}
	b273.items = []builder{&b272}
	var b275 = sequenceBuilder{id: 275, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b274 = charBuilder{}
	b275.items = []builder{&b274}
	var b296 = sequenceBuilder{id: 296, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b295 = charBuilder{}
	b296.items = []builder{&b295}
	var b298 = sequenceBuilder{id: 298, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b297 = charBuilder{}
	b298.items = []builder{&b297}
	b329.options = []builder{&b273, &b275, &b296, &b298}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b348.items = []builder{&b829, &b14}
	b349.items = []builder{&b829, &b14, &b348}
	b350.items = []builder{&b347, &b829, &b329, &b349, &b829, &b334}
	var b351 = sequenceBuilder{id: 351, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b351.items = []builder{&b829, &b350}
	b352.items = []builder{&b334, &b829, &b350, &b351}
	var b359 = sequenceBuilder{id: 359, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b335 = choiceBuilder{id: 335, commit: 66}
	b335.options = []builder{&b334, &b352}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b353.items = []builder{&b829, &b14}
	b354.items = []builder{&b14, &b353}
	var b330 = choiceBuilder{id: 330, commit: 66}
	var b303 = sequenceBuilder{id: 303, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b301 = charBuilder{}
	var b302 = charBuilder{}
	b303.items = []builder{&b301, &b302}
	var b306 = sequenceBuilder{id: 306, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b304 = charBuilder{}
	var b305 = charBuilder{}
	b306.items = []builder{&b304, &b305}
	var b308 = sequenceBuilder{id: 308, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b307 = charBuilder{}
	b308.items = []builder{&b307}
	var b311 = sequenceBuilder{id: 311, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b309 = charBuilder{}
	var b310 = charBuilder{}
	b311.items = []builder{&b309, &b310}
	var b313 = sequenceBuilder{id: 313, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b312 = charBuilder{}
	b313.items = []builder{&b312}
	var b316 = sequenceBuilder{id: 316, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b314 = charBuilder{}
	var b315 = charBuilder{}
	b316.items = []builder{&b314, &b315}
	b330.options = []builder{&b303, &b306, &b308, &b311, &b313, &b316}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b355.items = []builder{&b829, &b14}
	b356.items = []builder{&b829, &b14, &b355}
	b357.items = []builder{&b354, &b829, &b330, &b356, &b829, &b335}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b358.items = []builder{&b829, &b357}
	b359.items = []builder{&b335, &b829, &b357, &b358}
	var b366 = sequenceBuilder{id: 366, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b336 = choiceBuilder{id: 336, commit: 66}
	b336.options = []builder{&b335, &b359}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b360.items = []builder{&b829, &b14}
	b361.items = []builder{&b14, &b360}
	var b331 = sequenceBuilder{id: 331, commit: 66, ranges: [][]int{{1, 1}}}
	var b319 = sequenceBuilder{id: 319, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b317 = charBuilder{}
	var b318 = charBuilder{}
	b319.items = []builder{&b317, &b318}
	b331.items = []builder{&b319}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b362.items = []builder{&b829, &b14}
	b363.items = []builder{&b829, &b14, &b362}
	b364.items = []builder{&b361, &b829, &b331, &b363, &b829, &b336}
	var b365 = sequenceBuilder{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b365.items = []builder{&b829, &b364}
	b366.items = []builder{&b336, &b829, &b364, &b365}
	var b373 = sequenceBuilder{id: 373, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b337 = choiceBuilder{id: 337, commit: 66}
	b337.options = []builder{&b336, &b366}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b367.items = []builder{&b829, &b14}
	b368.items = []builder{&b14, &b367}
	var b332 = sequenceBuilder{id: 332, commit: 66, ranges: [][]int{{1, 1}}}
	var b322 = sequenceBuilder{id: 322, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b320 = charBuilder{}
	var b321 = charBuilder{}
	b322.items = []builder{&b320, &b321}
	b332.items = []builder{&b322}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b369.items = []builder{&b829, &b14}
	b370.items = []builder{&b829, &b14, &b369}
	b371.items = []builder{&b368, &b829, &b332, &b370, &b829, &b337}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b372.items = []builder{&b829, &b371}
	b373.items = []builder{&b337, &b829, &b371, &b372}
	b374.options = []builder{&b345, &b352, &b359, &b366, &b373}
	var b387 = sequenceBuilder{id: 387, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b380 = sequenceBuilder{id: 380, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b379 = sequenceBuilder{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b379.items = []builder{&b829, &b14}
	b380.items = []builder{&b829, &b14, &b379}
	var b376 = sequenceBuilder{id: 376, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b375 = charBuilder{}
	b376.items = []builder{&b375}
	var b382 = sequenceBuilder{id: 382, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b381 = sequenceBuilder{id: 381, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b381.items = []builder{&b829, &b14}
	b382.items = []builder{&b829, &b14, &b381}
	var b384 = sequenceBuilder{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b383 = sequenceBuilder{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b383.items = []builder{&b829, &b14}
	b384.items = []builder{&b829, &b14, &b383}
	var b378 = sequenceBuilder{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b377 = charBuilder{}
	b378.items = []builder{&b377}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b385.items = []builder{&b829, &b14}
	b386.items = []builder{&b829, &b14, &b385}
	b387.items = []builder{&b396, &b380, &b829, &b376, &b382, &b829, &b396, &b384, &b829, &b378, &b386, &b829, &b396}
	var b395 = sequenceBuilder{id: 395, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b388 = choiceBuilder{id: 388, commit: 66}
	b388.options = []builder{&b267, &b327, &b374, &b387}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b389.items = []builder{&b829, &b14}
	b390.items = []builder{&b14, &b389}
	var b325 = sequenceBuilder{id: 325, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b323 = charBuilder{}
	var b324 = charBuilder{}
	b325.items = []builder{&b323, &b324}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b391.items = []builder{&b829, &b14}
	b392.items = []builder{&b829, &b14, &b391}
	b393.items = []builder{&b390, &b829, &b325, &b392, &b829, &b388}
	var b394 = sequenceBuilder{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b394.items = []builder{&b829, &b393}
	b395.items = []builder{&b388, &b829, &b393, &b394}
	b396.options = []builder{&b267, &b327, &b374, &b387, &b395}
	b184.items = []builder{&b183, &b829, &b396}
	b185.items = []builder{&b181, &b829, &b184}
	var b433 = sequenceBuilder{id: 433, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b399 = sequenceBuilder{id: 399, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b397 = charBuilder{}
	var b398 = charBuilder{}
	b399.items = []builder{&b397, &b398}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b427.items = []builder{&b829, &b14}
	b428.items = []builder{&b829, &b14, &b427}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b429.items = []builder{&b829, &b14}
	b430.items = []builder{&b829, &b14, &b429}
	var b432 = sequenceBuilder{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b416 = sequenceBuilder{id: 416, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b409 = sequenceBuilder{id: 409, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b408.items = []builder{&b829, &b14}
	b409.items = []builder{&b14, &b408}
	var b404 = sequenceBuilder{id: 404, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b400 = charBuilder{}
	var b401 = charBuilder{}
	var b402 = charBuilder{}
	var b403 = charBuilder{}
	b404.items = []builder{&b400, &b401, &b402, &b403}
	var b411 = sequenceBuilder{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b410 = sequenceBuilder{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b410.items = []builder{&b829, &b14}
	b411.items = []builder{&b829, &b14, &b410}
	var b407 = sequenceBuilder{id: 407, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b405 = charBuilder{}
	var b406 = charBuilder{}
	b407.items = []builder{&b405, &b406}
	var b413 = sequenceBuilder{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b412 = sequenceBuilder{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b412.items = []builder{&b829, &b14}
	b413.items = []builder{&b829, &b14, &b412}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b414.items = []builder{&b829, &b14}
	b415.items = []builder{&b829, &b14, &b414}
	b416.items = []builder{&b409, &b829, &b404, &b411, &b829, &b407, &b413, &b829, &b396, &b415, &b829, &b190}
	var b431 = sequenceBuilder{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b431.items = []builder{&b829, &b416}
	b432.items = []builder{&b829, &b416, &b431}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b422.items = []builder{&b829, &b14}
	b423.items = []builder{&b14, &b422}
	var b421 = sequenceBuilder{id: 421, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b417 = charBuilder{}
	var b418 = charBuilder{}
	var b419 = charBuilder{}
	var b420 = charBuilder{}
	b421.items = []builder{&b417, &b418, &b419, &b420}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b424.items = []builder{&b829, &b14}
	b425.items = []builder{&b829, &b14, &b424}
	b426.items = []builder{&b423, &b829, &b421, &b425, &b829, &b190}
	b433.items = []builder{&b399, &b428, &b829, &b396, &b430, &b829, &b190, &b432, &b829, &b426}
	var b490 = sequenceBuilder{id: 490, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b475 = sequenceBuilder{id: 475, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b469 = charBuilder{}
	var b470 = charBuilder{}
	var b471 = charBuilder{}
	var b472 = charBuilder{}
	var b473 = charBuilder{}
	var b474 = charBuilder{}
	b475.items = []builder{&b469, &b470, &b471, &b472, &b473, &b474}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b486.items = []builder{&b829, &b14}
	b487.items = []builder{&b829, &b14, &b486}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b488.items = []builder{&b829, &b14}
	b489.items = []builder{&b829, &b14, &b488}
	var b477 = sequenceBuilder{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b476 = charBuilder{}
	b477.items = []builder{&b476}
	var b483 = sequenceBuilder{id: 483, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b478 = choiceBuilder{id: 478, commit: 2}
	var b468 = sequenceBuilder{id: 468, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b463 = sequenceBuilder{id: 463, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b456 = sequenceBuilder{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	var b454 = charBuilder{}
	var b455 = charBuilder{}
	b456.items = []builder{&b452, &b453, &b454, &b455}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b459.items = []builder{&b829, &b14}
	b460.items = []builder{&b829, &b14, &b459}
	var b462 = sequenceBuilder{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b461 = sequenceBuilder{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b461.items = []builder{&b829, &b14}
	b462.items = []builder{&b829, &b14, &b461}
	var b458 = sequenceBuilder{id: 458, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b457 = charBuilder{}
	b458.items = []builder{&b457}
	b463.items = []builder{&b456, &b460, &b829, &b396, &b462, &b829, &b458}
	var b467 = sequenceBuilder{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b465 = sequenceBuilder{id: 465, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b464 = charBuilder{}
	b465.items = []builder{&b464}
	var b466 = sequenceBuilder{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b466.items = []builder{&b829, &b465}
	b467.items = []builder{&b829, &b465, &b466}
	b468.items = []builder{&b463, &b467, &b829, &b797}
	var b451 = sequenceBuilder{id: 451, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b446 = sequenceBuilder{id: 446, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b441 = sequenceBuilder{id: 441, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b434 = charBuilder{}
	var b435 = charBuilder{}
	var b436 = charBuilder{}
	var b437 = charBuilder{}
	var b438 = charBuilder{}
	var b439 = charBuilder{}
	var b440 = charBuilder{}
	b441.items = []builder{&b434, &b435, &b436, &b437, &b438, &b439, &b440}
	var b445 = sequenceBuilder{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b444 = sequenceBuilder{id: 444, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b444.items = []builder{&b829, &b14}
	b445.items = []builder{&b829, &b14, &b444}
	var b443 = sequenceBuilder{id: 443, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b442 = charBuilder{}
	b443.items = []builder{&b442}
	b446.items = []builder{&b441, &b445, &b829, &b443}
	var b450 = sequenceBuilder{id: 450, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b448 = sequenceBuilder{id: 448, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b447 = charBuilder{}
	b448.items = []builder{&b447}
	var b449 = sequenceBuilder{id: 449, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b449.items = []builder{&b829, &b448}
	b450.items = []builder{&b829, &b448, &b449}
	b451.items = []builder{&b446, &b450, &b829, &b797}
	b478.options = []builder{&b468, &b451}
	var b482 = sequenceBuilder{id: 482, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b480 = sequenceBuilder{id: 480, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b479 = choiceBuilder{id: 479, commit: 2}
	b479.options = []builder{&b468, &b451, &b797}
	b480.items = []builder{&b811, &b829, &b479}
	var b481 = sequenceBuilder{id: 481, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b481.items = []builder{&b829, &b480}
	b482.items = []builder{&b829, &b480, &b481}
	b483.items = []builder{&b478, &b482}
	var b485 = sequenceBuilder{id: 485, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b484 = charBuilder{}
	b485.items = []builder{&b484}
	b490.items = []builder{&b475, &b487, &b829, &b396, &b489, &b829, &b477, &b829, &b811, &b829, &b483, &b829, &b811, &b829, &b485}
	var b552 = sequenceBuilder{id: 552, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b539 = sequenceBuilder{id: 539, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b533 = charBuilder{}
	var b534 = charBuilder{}
	var b535 = charBuilder{}
	var b536 = charBuilder{}
	var b537 = charBuilder{}
	var b538 = charBuilder{}
	b539.items = []builder{&b533, &b534, &b535, &b536, &b537, &b538}
	var b551 = sequenceBuilder{id: 551, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b550 = sequenceBuilder{id: 550, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b550.items = []builder{&b829, &b14}
	b551.items = []builder{&b829, &b14, &b550}
	var b541 = sequenceBuilder{id: 541, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b540 = charBuilder{}
	b541.items = []builder{&b540}
	var b547 = sequenceBuilder{id: 547, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b542 = choiceBuilder{id: 542, commit: 2}
	var b532 = sequenceBuilder{id: 532, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b527 = sequenceBuilder{id: 527, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b520 = sequenceBuilder{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b516 = charBuilder{}
	var b517 = charBuilder{}
	var b518 = charBuilder{}
	var b519 = charBuilder{}
	b520.items = []builder{&b516, &b517, &b518, &b519}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b523.items = []builder{&b829, &b14}
	b524.items = []builder{&b829, &b14, &b523}
	var b515 = choiceBuilder{id: 515, commit: 66}
	var b514 = sequenceBuilder{id: 514, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b513 = sequenceBuilder{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b512 = sequenceBuilder{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b512.items = []builder{&b829, &b14}
	b513.items = []builder{&b829, &b14, &b512}
	b514.items = []builder{&b103, &b513, &b829, &b511}
	b515.options = []builder{&b500, &b511, &b514}
	var b526 = sequenceBuilder{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b525.items = []builder{&b829, &b14}
	b526.items = []builder{&b829, &b14, &b525}
	var b522 = sequenceBuilder{id: 522, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b521 = charBuilder{}
	b522.items = []builder{&b521}
	b527.items = []builder{&b520, &b524, &b829, &b515, &b526, &b829, &b522}
	var b531 = sequenceBuilder{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b529 = sequenceBuilder{id: 529, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b528 = charBuilder{}
	b529.items = []builder{&b528}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b530.items = []builder{&b829, &b529}
	b531.items = []builder{&b829, &b529, &b530}
	b532.items = []builder{&b527, &b531, &b829, &b797}
	b542.options = []builder{&b532, &b451}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b543 = choiceBuilder{id: 543, commit: 2}
	b543.options = []builder{&b532, &b451, &b797}
	b544.items = []builder{&b811, &b829, &b543}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b545.items = []builder{&b829, &b544}
	b546.items = []builder{&b829, &b544, &b545}
	b547.items = []builder{&b542, &b546}
	var b549 = sequenceBuilder{id: 549, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b548 = charBuilder{}
	b549.items = []builder{&b548}
	b552.items = []builder{&b539, &b551, &b829, &b541, &b829, &b811, &b829, &b547, &b829, &b811, &b829, &b549}
	var b593 = sequenceBuilder{id: 593, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b582 = sequenceBuilder{id: 582, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b579 = charBuilder{}
	var b580 = charBuilder{}
	var b581 = charBuilder{}
	b582.items = []builder{&b579, &b580, &b581}
	var b592 = choiceBuilder{id: 592, commit: 2}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b583 = sequenceBuilder{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b583.items = []builder{&b829, &b14}
	b584.items = []builder{&b14, &b583}
	var b578 = choiceBuilder{id: 578, commit: 66}
	var b577 = choiceBuilder{id: 577, commit: 64, name: "range-over-expression"}
	var b576 = sequenceBuilder{id: 576, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b572.items = []builder{&b829, &b14}
	b573.items = []builder{&b829, &b14, &b572}
	var b570 = sequenceBuilder{id: 570, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b568 = charBuilder{}
	var b569 = charBuilder{}
	b570.items = []builder{&b568, &b569}
	var b575 = sequenceBuilder{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b574 = sequenceBuilder{id: 574, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b574.items = []builder{&b829, &b14}
	b575.items = []builder{&b829, &b14, &b574}
	var b571 = choiceBuilder{id: 571, commit: 2}
	b571.options = []builder{&b396, &b225}
	b576.items = []builder{&b103, &b573, &b829, &b570, &b575, &b829, &b571}
	b577.options = []builder{&b576, &b225}
	b578.options = []builder{&b396, &b577}
	b585.items = []builder{&b584, &b829, &b578}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b586.items = []builder{&b829, &b14}
	b587.items = []builder{&b829, &b14, &b586}
	b588.items = []builder{&b585, &b587, &b829, &b190}
	var b591 = sequenceBuilder{id: 591, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b589.items = []builder{&b829, &b14}
	b590.items = []builder{&b14, &b589}
	b591.items = []builder{&b590, &b829, &b190}
	b592.options = []builder{&b588, &b591}
	b593.items = []builder{&b582, &b829, &b592}
	var b741 = choiceBuilder{id: 741, commit: 66}
	var b654 = sequenceBuilder{id: 654, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b650 = sequenceBuilder{id: 650, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b647 = charBuilder{}
	var b648 = charBuilder{}
	var b649 = charBuilder{}
	b650.items = []builder{&b647, &b648, &b649}
	var b653 = sequenceBuilder{id: 653, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b652 = sequenceBuilder{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b652.items = []builder{&b829, &b14}
	b653.items = []builder{&b829, &b14, &b652}
	var b651 = choiceBuilder{id: 651, commit: 2}
	var b641 = sequenceBuilder{id: 641, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b640 = sequenceBuilder{id: 640, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b637 = sequenceBuilder{id: 637, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b635.items = []builder{&b829, &b14}
	b636.items = []builder{&b14, &b635}
	var b634 = sequenceBuilder{id: 634, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b633 = charBuilder{}
	b634.items = []builder{&b633}
	b637.items = []builder{&b636, &b829, &b634}
	var b639 = sequenceBuilder{id: 639, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b638 = sequenceBuilder{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b638.items = []builder{&b829, &b14}
	b639.items = []builder{&b829, &b14, &b638}
	b640.items = []builder{&b103, &b829, &b637, &b639, &b829, &b396}
	b641.items = []builder{&b640}
	var b646 = sequenceBuilder{id: 646, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b643 = sequenceBuilder{id: 643, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b642 = charBuilder{}
	b643.items = []builder{&b642}
	var b645 = sequenceBuilder{id: 645, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b644 = sequenceBuilder{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b644.items = []builder{&b829, &b14}
	b645.items = []builder{&b829, &b14, &b644}
	b646.items = []builder{&b643, &b645, &b829, &b640}
	b651.options = []builder{&b641, &b646}
	b654.items = []builder{&b650, &b653, &b829, &b651}
	var b675 = sequenceBuilder{id: 675, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b668 = sequenceBuilder{id: 668, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b665 = charBuilder{}
	var b666 = charBuilder{}
	var b667 = charBuilder{}
	b668.items = []builder{&b665, &b666, &b667}
	var b674 = sequenceBuilder{id: 674, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b673 = sequenceBuilder{id: 673, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b673.items = []builder{&b829, &b14}
	b674.items = []builder{&b829, &b14, &b673}
	var b670 = sequenceBuilder{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b669 = charBuilder{}
	b670.items = []builder{&b669}
	var b660 = sequenceBuilder{id: 660, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b655 = choiceBuilder{id: 655, commit: 2}
	b655.options = []builder{&b641, &b646}
	var b659 = sequenceBuilder{id: 659, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b657 = sequenceBuilder{id: 657, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b656 = choiceBuilder{id: 656, commit: 2}
	b656.options = []builder{&b641, &b646}
	b657.items = []builder{&b113, &b829, &b656}
	var b658 = sequenceBuilder{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b658.items = []builder{&b829, &b657}
	b659.items = []builder{&b829, &b657, &b658}
	b660.items = []builder{&b655, &b659}
	var b672 = sequenceBuilder{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b671 = charBuilder{}
	b672.items = []builder{&b671}
	b675.items = []builder{&b668, &b674, &b829, &b670, &b829, &b113, &b829, &b660, &b829, &b113, &b829, &b672}
	var b690 = sequenceBuilder{id: 690, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b679 = sequenceBuilder{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b676 = charBuilder{}
	var b677 = charBuilder{}
	var b678 = charBuilder{}
	b679.items = []builder{&b676, &b677, &b678}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b686.items = []builder{&b829, &b14}
	b687.items = []builder{&b829, &b14, &b686}
	var b681 = sequenceBuilder{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	b681.items = []builder{&b680}
	var b689 = sequenceBuilder{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b688.items = []builder{&b829, &b14}
	b689.items = []builder{&b829, &b14, &b688}
	var b683 = sequenceBuilder{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b682 = charBuilder{}
	b683.items = []builder{&b682}
	var b664 = sequenceBuilder{id: 664, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b663 = sequenceBuilder{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b661 = sequenceBuilder{id: 661, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b661.items = []builder{&b113, &b829, &b641}
	var b662 = sequenceBuilder{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b662.items = []builder{&b829, &b661}
	b663.items = []builder{&b829, &b661, &b662}
	b664.items = []builder{&b641, &b663}
	var b685 = sequenceBuilder{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b684 = charBuilder{}
	b685.items = []builder{&b684}
	b690.items = []builder{&b679, &b687, &b829, &b681, &b689, &b829, &b683, &b829, &b113, &b829, &b664, &b829, &b113, &b829, &b685}
	var b706 = sequenceBuilder{id: 706, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b702 = sequenceBuilder{id: 702, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b700 = charBuilder{}
	var b701 = charBuilder{}
	b702.items = []builder{&b700, &b701}
	var b705 = sequenceBuilder{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b704.items = []builder{&b829, &b14}
	b705.items = []builder{&b829, &b14, &b704}
	var b703 = choiceBuilder{id: 703, commit: 2}
	var b694 = sequenceBuilder{id: 694, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b693 = sequenceBuilder{id: 693, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b692 = sequenceBuilder{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b691 = sequenceBuilder{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b691.items = []builder{&b829, &b14}
	b692.items = []builder{&b829, &b14, &b691}
	b693.items = []builder{&b103, &b692, &b829, &b200}
	b694.items = []builder{&b693}
	var b699 = sequenceBuilder{id: 699, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b696 = sequenceBuilder{id: 696, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b695 = charBuilder{}
	b696.items = []builder{&b695}
	var b698 = sequenceBuilder{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b697 = sequenceBuilder{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b697.items = []builder{&b829, &b14}
	b698.items = []builder{&b829, &b14, &b697}
	b699.items = []builder{&b696, &b698, &b829, &b693}
	b703.options = []builder{&b694, &b699}
	b706.items = []builder{&b702, &b705, &b829, &b703}
	var b726 = sequenceBuilder{id: 726, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b719 = sequenceBuilder{id: 719, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b717 = charBuilder{}
	var b718 = charBuilder{}
	b719.items = []builder{&b717, &b718}
	var b725 = sequenceBuilder{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b724 = sequenceBuilder{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b724.items = []builder{&b829, &b14}
	b725.items = []builder{&b829, &b14, &b724}
	var b721 = sequenceBuilder{id: 721, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b720 = charBuilder{}
	b721.items = []builder{&b720}
	var b716 = sequenceBuilder{id: 716, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b711 = choiceBuilder{id: 711, commit: 2}
	b711.options = []builder{&b694, &b699}
	var b715 = sequenceBuilder{id: 715, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b712 = choiceBuilder{id: 712, commit: 2}
	b712.options = []builder{&b694, &b699}
	b713.items = []builder{&b113, &b829, &b712}
	var b714 = sequenceBuilder{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b714.items = []builder{&b829, &b713}
	b715.items = []builder{&b829, &b713, &b714}
	b716.items = []builder{&b711, &b715}
	var b723 = sequenceBuilder{id: 723, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b722 = charBuilder{}
	b723.items = []builder{&b722}
	b726.items = []builder{&b719, &b725, &b829, &b721, &b829, &b113, &b829, &b716, &b829, &b113, &b829, &b723}
	var b740 = sequenceBuilder{id: 740, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b729 = sequenceBuilder{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b727 = charBuilder{}
	var b728 = charBuilder{}
	b729.items = []builder{&b727, &b728}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b736.items = []builder{&b829, &b14}
	b737.items = []builder{&b829, &b14, &b736}
	var b731 = sequenceBuilder{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b730 = charBuilder{}
	b731.items = []builder{&b730}
	var b739 = sequenceBuilder{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b738 = sequenceBuilder{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b738.items = []builder{&b829, &b14}
	b739.items = []builder{&b829, &b14, &b738}
	var b733 = sequenceBuilder{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b732 = charBuilder{}
	b733.items = []builder{&b732}
	var b710 = sequenceBuilder{id: 710, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b709 = sequenceBuilder{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b707 = sequenceBuilder{id: 707, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b707.items = []builder{&b113, &b829, &b694}
	var b708 = sequenceBuilder{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b708.items = []builder{&b829, &b707}
	b709.items = []builder{&b829, &b707, &b708}
	b710.items = []builder{&b694, &b709}
	var b735 = sequenceBuilder{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b734 = charBuilder{}
	b735.items = []builder{&b734}
	b740.items = []builder{&b729, &b737, &b829, &b731, &b739, &b829, &b733, &b829, &b113, &b829, &b710, &b829, &b113, &b829, &b735}
	b741.options = []builder{&b654, &b675, &b690, &b706, &b726, &b740}
	var b776 = choiceBuilder{id: 776, commit: 64, name: "use"}
	var b764 = sequenceBuilder{id: 764, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b761 = sequenceBuilder{id: 761, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b758 = charBuilder{}
	var b759 = charBuilder{}
	var b760 = charBuilder{}
	b761.items = []builder{&b758, &b759, &b760}
	var b763 = sequenceBuilder{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b762 = sequenceBuilder{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b762.items = []builder{&b829, &b14}
	b763.items = []builder{&b829, &b14, &b762}
	var b753 = choiceBuilder{id: 753, commit: 64, name: "use-fact"}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b744 = choiceBuilder{id: 744, commit: 2}
	var b743 = sequenceBuilder{id: 743, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b742 = charBuilder{}
	b743.items = []builder{&b742}
	b744.options = []builder{&b103, &b743}
	var b749 = sequenceBuilder{id: 749, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b748 = sequenceBuilder{id: 748, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b747 = sequenceBuilder{id: 747, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b747.items = []builder{&b829, &b14}
	b748.items = []builder{&b14, &b747}
	var b746 = sequenceBuilder{id: 746, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b745 = charBuilder{}
	b746.items = []builder{&b745}
	b749.items = []builder{&b748, &b829, &b746}
	var b751 = sequenceBuilder{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b750 = sequenceBuilder{id: 750, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b750.items = []builder{&b829, &b14}
	b751.items = []builder{&b829, &b14, &b750}
	b752.items = []builder{&b744, &b829, &b749, &b751, &b829, &b86}
	b753.options = []builder{&b86, &b752}
	b764.items = []builder{&b761, &b763, &b829, &b753}
	var b775 = sequenceBuilder{id: 775, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b768 = sequenceBuilder{id: 768, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b765 = charBuilder{}
	var b766 = charBuilder{}
	var b767 = charBuilder{}
	b768.items = []builder{&b765, &b766, &b767}
	var b774 = sequenceBuilder{id: 774, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b773 = sequenceBuilder{id: 773, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b773.items = []builder{&b829, &b14}
	b774.items = []builder{&b829, &b14, &b773}
	var b770 = sequenceBuilder{id: 770, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b769 = charBuilder{}
	b770.items = []builder{&b769}
	var b757 = sequenceBuilder{id: 757, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b754.items = []builder{&b113, &b829, &b753}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b755.items = []builder{&b829, &b754}
	b756.items = []builder{&b829, &b754, &b755}
	b757.items = []builder{&b753, &b756}
	var b772 = sequenceBuilder{id: 772, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b771 = charBuilder{}
	b772.items = []builder{&b771}
	b775.items = []builder{&b768, &b774, &b829, &b770, &b829, &b113, &b829, &b757, &b829, &b113, &b829, &b772}
	b776.options = []builder{&b764, &b775}
	var b786 = sequenceBuilder{id: 786, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b783 = sequenceBuilder{id: 783, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b777 = charBuilder{}
	var b778 = charBuilder{}
	var b779 = charBuilder{}
	var b780 = charBuilder{}
	var b781 = charBuilder{}
	var b782 = charBuilder{}
	b783.items = []builder{&b777, &b778, &b779, &b780, &b781, &b782}
	var b785 = sequenceBuilder{id: 785, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b784 = sequenceBuilder{id: 784, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b784.items = []builder{&b829, &b14}
	b785.items = []builder{&b829, &b14, &b784}
	b786.items = []builder{&b783, &b785, &b829, &b741}
	var b806 = sequenceBuilder{id: 806, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b799 = sequenceBuilder{id: 799, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b798 = charBuilder{}
	b799.items = []builder{&b798}
	var b803 = sequenceBuilder{id: 803, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b802 = sequenceBuilder{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b802.items = []builder{&b829, &b14}
	b803.items = []builder{&b829, &b14, &b802}
	var b805 = sequenceBuilder{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b804 = sequenceBuilder{id: 804, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b804.items = []builder{&b829, &b14}
	b805.items = []builder{&b829, &b14, &b804}
	var b801 = sequenceBuilder{id: 801, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b800 = charBuilder{}
	b801.items = []builder{&b800}
	b806.items = []builder{&b799, &b803, &b829, &b797, &b805, &b829, &b801}
	b797.options = []builder{&b185, &b433, &b490, &b552, &b593, &b741, &b776, &b786, &b806, &b787}
	var b814 = sequenceBuilder{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b812 = sequenceBuilder{id: 812, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b812.items = []builder{&b811, &b829, &b797}
	var b813 = sequenceBuilder{id: 813, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b813.items = []builder{&b829, &b812}
	b814.items = []builder{&b829, &b812, &b813}
	b815.items = []builder{&b797, &b814}
	b830.items = []builder{&b826, &b829, &b811, &b829, &b815, &b829, &b811}
	b831.items = []builder{&b829, &b830, &b829}

	return parseInput(r, &p831, &b831)
}
