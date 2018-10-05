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
	var p837 = sequenceParser{id: 837, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p835 = choiceParser{id: 835, commit: 2}
	var p833 = choiceParser{id: 833, commit: 70, name: "ws", generalizations: []int{835}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833, 835}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833, 835}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833, 835}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833, 835}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833, 835}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833, 835}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p833.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p834 = sequenceParser{id: 834, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{835}}
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
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 111}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p39.items = []parser{&p14, &p835, &p38}
	var p40 = sequenceParser{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p40.items = []parser{&p835, &p39}
	p41.items = []parser{&p835, &p39, &p40}
	p42.items = []parser{&p38, &p41}
	p834.items = []parser{&p42}
	p835.options = []parser{&p833, &p834}
	var p836 = sequenceParser{id: 836, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p832 = sequenceParser{id: 832, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p829 = sequenceParser{id: 829, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p827 = charParser{id: 827, chars: []rune{35}}
	var p828 = charParser{id: 828, chars: []rune{33}}
	p829.items = []parser{&p827, &p828}
	var p826 = sequenceParser{id: 826, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p825 = sequenceParser{id: 825, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p823 = sequenceParser{id: 823, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p822 = charParser{id: 822, not: true, chars: []rune{10}}
	p823.items = []parser{&p822}
	var p824 = sequenceParser{id: 824, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p824.items = []parser{&p835, &p823}
	p825.items = []parser{&p823, &p824}
	p826.items = []parser{&p825}
	var p831 = sequenceParser{id: 831, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p830 = charParser{id: 830, chars: []rune{10}}
	p831.items = []parser{&p830}
	p832.items = []parser{&p829, &p835, &p826, &p835, &p831}
	var p817 = sequenceParser{id: 817, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p815 = choiceParser{id: 815, commit: 2}
	var p814 = sequenceParser{id: 814, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815}}
	var p813 = charParser{id: 813, chars: []rune{59}}
	p814.items = []parser{&p813}
	p815.options = []parser{&p814, &p14}
	var p816 = sequenceParser{id: 816, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p816.items = []parser{&p835, &p815}
	p817.items = []parser{&p815, &p816}
	var p821 = sequenceParser{id: 821, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p803 = choiceParser{id: 803, commit: 66, name: "statement", generalizations: []int{477, 541}}
	var p185 = sequenceParser{id: 185, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{803, 477, 541}}
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
	p182.items = []parser{&p835, &p14}
	p183.items = []parser{&p14, &p182}
	var p394 = choiceParser{id: 394, commit: 66, name: "expression", generalizations: []int{114, 793, 197, 576, 569, 803}}
	var p266 = choiceParser{id: 266, commit: 66, name: "primary-expression", generalizations: []int{114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p60 = choiceParser{id: 60, commit: 64, name: "int", generalizations: []int{266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p51 = sequenceParser{id: 51, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p50 = sequenceParser{id: 50, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p49 = charParser{id: 49, ranges: [][]rune{{49, 57}}}
	p50.items = []parser{&p49}
	var p44 = sequenceParser{id: 44, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p43 = charParser{id: 43, ranges: [][]rune{{48, 57}}}
	p44.items = []parser{&p43}
	p51.items = []parser{&p50, &p44}
	var p54 = sequenceParser{id: 54, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p53 = sequenceParser{id: 53, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p52 = charParser{id: 52, chars: []rune{48}}
	p53.items = []parser{&p52}
	var p46 = sequenceParser{id: 46, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 55}}}
	p46.items = []parser{&p45}
	p54.items = []parser{&p53, &p46}
	var p59 = sequenceParser{id: 59, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{60, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
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
	var p73 = choiceParser{id: 73, commit: 72, name: "float", generalizations: []int{266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p68 = sequenceParser{id: 68, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{73, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
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
	var p71 = sequenceParser{id: 71, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{73, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p70 = sequenceParser{id: 70, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p69 = charParser{id: 69, chars: []rune{46}}
	p70.items = []parser{&p69}
	p71.items = []parser{&p70, &p44, &p65}
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{73, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	p72.items = []parser{&p44, &p65}
	p73.options = []parser{&p68, &p71, &p72}
	var p86 = sequenceParser{id: 86, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 114, 139, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 751, 803}}
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
	var p98 = choiceParser{id: 98, commit: 66, name: "bool", generalizations: []int{266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p91 = sequenceParser{id: 91, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p87 = charParser{id: 87, chars: []rune{116}}
	var p88 = charParser{id: 88, chars: []rune{114}}
	var p89 = charParser{id: 89, chars: []rune{117}}
	var p90 = charParser{id: 90, chars: []rune{101}}
	p91.items = []parser{&p87, &p88, &p89, &p90}
	var p97 = sequenceParser{id: 97, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p92 = charParser{id: 92, chars: []rune{102}}
	var p93 = charParser{id: 93, chars: []rune{97}}
	var p94 = charParser{id: 94, chars: []rune{108}}
	var p95 = charParser{id: 95, chars: []rune{115}}
	var p96 = charParser{id: 96, chars: []rune{101}}
	p97.items = []parser{&p92, &p93, &p94, &p95, &p96}
	p98.options = []parser{&p91, &p97}
	var p499 = sequenceParser{id: 499, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 114, 793, 197, 394, 332, 333, 334, 335, 336, 337, 513, 576, 569, 803}}
	var p496 = sequenceParser{id: 496, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p489 = charParser{id: 489, chars: []rune{114}}
	var p490 = charParser{id: 490, chars: []rune{101}}
	var p491 = charParser{id: 491, chars: []rune{99}}
	var p492 = charParser{id: 492, chars: []rune{101}}
	var p493 = charParser{id: 493, chars: []rune{105}}
	var p494 = charParser{id: 494, chars: []rune{118}}
	var p495 = charParser{id: 495, chars: []rune{101}}
	p496.items = []parser{&p489, &p490, &p491, &p492, &p493, &p494, &p495}
	var p498 = sequenceParser{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p497 = sequenceParser{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p497.items = []parser{&p835, &p14}
	p498.items = []parser{&p835, &p14, &p497}
	p499.items = []parser{&p496, &p498, &p835, &p266}
	var p103 = sequenceParser{id: 103, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{266, 114, 139, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 742, 803}}
	var p100 = sequenceParser{id: 100, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p99 = charParser{id: 99, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p100.items = []parser{&p99}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p102.items = []parser{&p101}
	p103.items = []parser{&p100, &p102}
	var p124 = sequenceParser{id: 124, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{114, 266, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
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
	p112.items = []parser{&p835, &p111}
	p113.items = []parser{&p111, &p112}
	var p118 = sequenceParser{id: 118, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p114 = choiceParser{id: 114, commit: 66, name: "list-item"}
	var p108 = sequenceParser{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{114, 147, 148}}
	var p107 = sequenceParser{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p104 = charParser{id: 104, chars: []rune{46}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	p107.items = []parser{&p104, &p105, &p106}
	p108.items = []parser{&p266, &p835, &p107}
	p114.options = []parser{&p394, &p108}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p115 = sequenceParser{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p115.items = []parser{&p113, &p835, &p114}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p116.items = []parser{&p835, &p115}
	p117.items = []parser{&p835, &p115, &p116}
	p118.items = []parser{&p114, &p117}
	var p122 = sequenceParser{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p121 = charParser{id: 121, chars: []rune{93}}
	p122.items = []parser{&p121}
	p123.items = []parser{&p120, &p835, &p113, &p835, &p118, &p835, &p113, &p835, &p122}
	p124.items = []parser{&p123}
	var p129 = sequenceParser{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p126 = sequenceParser{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p125 = charParser{id: 125, chars: []rune{126}}
	p126.items = []parser{&p125}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p127 = sequenceParser{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p127.items = []parser{&p835, &p14}
	p128.items = []parser{&p835, &p14, &p127}
	p129.items = []parser{&p126, &p128, &p835, &p123}
	var p158 = sequenceParser{id: 158, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{266, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
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
	p134.items = []parser{&p835, &p14}
	p135.items = []parser{&p835, &p14, &p134}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p136 = sequenceParser{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p136.items = []parser{&p835, &p14}
	p137.items = []parser{&p835, &p14, &p136}
	var p133 = sequenceParser{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p132 = charParser{id: 132, chars: []rune{93}}
	p133.items = []parser{&p132}
	p138.items = []parser{&p131, &p135, &p835, &p394, &p137, &p835, &p133}
	p139.options = []parser{&p103, &p86, &p138}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p142 = sequenceParser{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p142.items = []parser{&p835, &p14}
	p143.items = []parser{&p835, &p14, &p142}
	var p141 = sequenceParser{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p140 = charParser{id: 140, chars: []rune{58}}
	p141.items = []parser{&p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p835, &p14}
	p145.items = []parser{&p835, &p14, &p144}
	p146.items = []parser{&p139, &p143, &p835, &p141, &p145, &p835, &p394}
	p147.options = []parser{&p146, &p108}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p149 = sequenceParser{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p148 = choiceParser{id: 148, commit: 2}
	p148.options = []parser{&p146, &p108}
	p149.items = []parser{&p113, &p835, &p148}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p150.items = []parser{&p835, &p149}
	p151.items = []parser{&p835, &p149, &p150}
	p152.items = []parser{&p147, &p151}
	var p156 = sequenceParser{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p155 = charParser{id: 155, chars: []rune{125}}
	p156.items = []parser{&p155}
	p157.items = []parser{&p154, &p835, &p113, &p835, &p152, &p835, &p113, &p835, &p156}
	p158.items = []parser{&p157}
	var p163 = sequenceParser{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 793, 197, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p160 = sequenceParser{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p159 = charParser{id: 159, chars: []rune{126}}
	p160.items = []parser{&p159}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p161 = sequenceParser{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p161.items = []parser{&p835, &p14}
	p162.items = []parser{&p835, &p14, &p161}
	p163.items = []parser{&p160, &p162, &p835, &p157}
	var p206 = sequenceParser{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793, 197, 266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p203 = sequenceParser{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p201 = charParser{id: 201, chars: []rune{102}}
	var p202 = charParser{id: 202, chars: []rune{110}}
	p203.items = []parser{&p201, &p202}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p204 = sequenceParser{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p204.items = []parser{&p835, &p14}
	p205.items = []parser{&p835, &p14, &p204}
	var p200 = sequenceParser{id: 200, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p192 = sequenceParser{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p191 = charParser{id: 191, chars: []rune{40}}
	p192.items = []parser{&p191}
	var p194 = choiceParser{id: 194, commit: 2}
	var p167 = sequenceParser{id: 167, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{194}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p164.items = []parser{&p113, &p835, &p103}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p165.items = []parser{&p835, &p164}
	p166.items = []parser{&p835, &p164, &p165}
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
	p172.items = []parser{&p835, &p14}
	p173.items = []parser{&p835, &p14, &p172}
	p174.items = []parser{&p171, &p173, &p835, &p103}
	p193.items = []parser{&p167, &p835, &p113, &p835, &p174}
	p194.options = []parser{&p167, &p193, &p174}
	var p196 = sequenceParser{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p195 = charParser{id: 195, chars: []rune{41}}
	p196.items = []parser{&p195}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p198 = sequenceParser{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p198.items = []parser{&p835, &p14}
	p199.items = []parser{&p835, &p14, &p198}
	var p197 = choiceParser{id: 197, commit: 2}
	var p793 = choiceParser{id: 793, commit: 66, name: "simple-statement", generalizations: []int{197, 803}}
	var p509 = sequenceParser{id: 509, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793, 197, 513, 803}}
	var p504 = sequenceParser{id: 504, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p500 = charParser{id: 500, chars: []rune{115}}
	var p501 = charParser{id: 501, chars: []rune{101}}
	var p502 = charParser{id: 502, chars: []rune{110}}
	var p503 = charParser{id: 503, chars: []rune{100}}
	p504.items = []parser{&p500, &p501, &p502, &p503}
	var p506 = sequenceParser{id: 506, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p505 = sequenceParser{id: 505, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p505.items = []parser{&p835, &p14}
	p506.items = []parser{&p835, &p14, &p505}
	var p508 = sequenceParser{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p507 = sequenceParser{id: 507, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p507.items = []parser{&p835, &p14}
	p508.items = []parser{&p835, &p14, &p507}
	p509.items = []parser{&p504, &p506, &p835, &p266, &p508, &p835, &p266}
	var p556 = sequenceParser{id: 556, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793, 197, 803}}
	var p553 = sequenceParser{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p551 = charParser{id: 551, chars: []rune{103}}
	var p552 = charParser{id: 552, chars: []rune{111}}
	p553.items = []parser{&p551, &p552}
	var p555 = sequenceParser{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p554 = sequenceParser{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p554.items = []parser{&p835, &p14}
	p555.items = []parser{&p835, &p14, &p554}
	var p256 = sequenceParser{id: 256, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p253 = sequenceParser{id: 253, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p252 = charParser{id: 252, chars: []rune{40}}
	p253.items = []parser{&p252}
	var p255 = sequenceParser{id: 255, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p254 = charParser{id: 254, chars: []rune{41}}
	p255.items = []parser{&p254}
	p256.items = []parser{&p266, &p835, &p253, &p835, &p113, &p835, &p118, &p835, &p113, &p835, &p255}
	p556.items = []parser{&p553, &p555, &p835, &p256}
	var p565 = sequenceParser{id: 565, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793, 197, 803}}
	var p562 = sequenceParser{id: 562, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p557 = charParser{id: 557, chars: []rune{100}}
	var p558 = charParser{id: 558, chars: []rune{101}}
	var p559 = charParser{id: 559, chars: []rune{102}}
	var p560 = charParser{id: 560, chars: []rune{101}}
	var p561 = charParser{id: 561, chars: []rune{114}}
	p562.items = []parser{&p557, &p558, &p559, &p560, &p561}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p563 = sequenceParser{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p563.items = []parser{&p835, &p14}
	p564.items = []parser{&p835, &p14, &p563}
	p565.items = []parser{&p562, &p564, &p835, &p256}
	var p630 = choiceParser{id: 630, commit: 64, name: "assignment", generalizations: []int{793, 197, 803}}
	var p610 = sequenceParser{id: 610, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{630, 793, 197, 803}}
	var p607 = sequenceParser{id: 607, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p604 = charParser{id: 604, chars: []rune{115}}
	var p605 = charParser{id: 605, chars: []rune{101}}
	var p606 = charParser{id: 606, chars: []rune{116}}
	p607.items = []parser{&p604, &p605, &p606}
	var p609 = sequenceParser{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p608 = sequenceParser{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p608.items = []parser{&p835, &p14}
	p609.items = []parser{&p835, &p14, &p608}
	var p599 = sequenceParser{id: 599, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p594.items = []parser{&p835, &p14}
	p595.items = []parser{&p14, &p594}
	var p593 = sequenceParser{id: 593, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p592 = charParser{id: 592, chars: []rune{61}}
	p593.items = []parser{&p592}
	p596.items = []parser{&p595, &p835, &p593}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p597.items = []parser{&p835, &p14}
	p598.items = []parser{&p835, &p14, &p597}
	p599.items = []parser{&p266, &p835, &p596, &p598, &p835, &p394}
	p610.items = []parser{&p607, &p609, &p835, &p599}
	var p617 = sequenceParser{id: 617, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{630, 793, 197, 803}}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p613 = sequenceParser{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p613.items = []parser{&p835, &p14}
	p614.items = []parser{&p835, &p14, &p613}
	var p612 = sequenceParser{id: 612, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p611 = charParser{id: 611, chars: []rune{61}}
	p612.items = []parser{&p611}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p615.items = []parser{&p835, &p14}
	p616.items = []parser{&p835, &p14, &p615}
	p617.items = []parser{&p266, &p614, &p835, &p612, &p616, &p835, &p394}
	var p629 = sequenceParser{id: 629, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{630, 793, 197, 803}}
	var p621 = sequenceParser{id: 621, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p618 = charParser{id: 618, chars: []rune{115}}
	var p619 = charParser{id: 619, chars: []rune{101}}
	var p620 = charParser{id: 620, chars: []rune{116}}
	p621.items = []parser{&p618, &p619, &p620}
	var p628 = sequenceParser{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p627 = sequenceParser{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p627.items = []parser{&p835, &p14}
	p628.items = []parser{&p835, &p14, &p627}
	var p623 = sequenceParser{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p622 = charParser{id: 622, chars: []rune{40}}
	p623.items = []parser{&p622}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p603 = sequenceParser{id: 603, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p602 = sequenceParser{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p600 = sequenceParser{id: 600, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p600.items = []parser{&p113, &p835, &p599}
	var p601 = sequenceParser{id: 601, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p601.items = []parser{&p835, &p600}
	p602.items = []parser{&p835, &p600, &p601}
	p603.items = []parser{&p599, &p602}
	p624.items = []parser{&p113, &p835, &p603}
	var p626 = sequenceParser{id: 626, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p625 = charParser{id: 625, chars: []rune{41}}
	p626.items = []parser{&p625}
	p629.items = []parser{&p621, &p628, &p835, &p623, &p835, &p624, &p835, &p113, &p835, &p626}
	p630.options = []parser{&p610, &p617, &p629}
	var p802 = sequenceParser{id: 802, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793, 197, 803}}
	var p795 = sequenceParser{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p794 = charParser{id: 794, chars: []rune{40}}
	p795.items = []parser{&p794}
	var p799 = sequenceParser{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p798 = sequenceParser{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p798.items = []parser{&p835, &p14}
	p799.items = []parser{&p835, &p14, &p798}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p800.items = []parser{&p835, &p14}
	p801.items = []parser{&p835, &p14, &p800}
	var p797 = sequenceParser{id: 797, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p796 = charParser{id: 796, chars: []rune{41}}
	p797.items = []parser{&p796}
	p802.items = []parser{&p795, &p799, &p835, &p793, &p801, &p835, &p797}
	p793.options = []parser{&p509, &p556, &p565, &p630, &p802, &p394}
	var p190 = sequenceParser{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{197}}
	var p187 = sequenceParser{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p186 = charParser{id: 186, chars: []rune{123}}
	p187.items = []parser{&p186}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{125}}
	p189.items = []parser{&p188}
	p190.items = []parser{&p187, &p835, &p817, &p835, &p821, &p835, &p817, &p835, &p189}
	p197.options = []parser{&p793, &p190}
	p200.items = []parser{&p192, &p835, &p113, &p835, &p194, &p835, &p113, &p835, &p196, &p199, &p835, &p197}
	p206.items = []parser{&p203, &p205, &p835, &p200}
	var p216 = sequenceParser{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p209 = sequenceParser{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p207 = charParser{id: 207, chars: []rune{102}}
	var p208 = charParser{id: 208, chars: []rune{110}}
	p209.items = []parser{&p207, &p208}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p212 = sequenceParser{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p212.items = []parser{&p835, &p14}
	p213.items = []parser{&p835, &p14, &p212}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p210 = charParser{id: 210, chars: []rune{126}}
	p211.items = []parser{&p210}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p214.items = []parser{&p835, &p14}
	p215.items = []parser{&p835, &p14, &p214}
	p216.items = []parser{&p209, &p213, &p835, &p211, &p215, &p835, &p200}
	var p244 = choiceParser{id: 244, commit: 64, name: "expression-indexer", generalizations: []int{266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p234 = sequenceParser{id: 234, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{244, 266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p227 = sequenceParser{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p226 = charParser{id: 226, chars: []rune{91}}
	p227.items = []parser{&p226}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p230 = sequenceParser{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p230.items = []parser{&p835, &p14}
	p231.items = []parser{&p835, &p14, &p230}
	var p233 = sequenceParser{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p232 = sequenceParser{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p232.items = []parser{&p835, &p14}
	p233.items = []parser{&p835, &p14, &p232}
	var p229 = sequenceParser{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p228 = charParser{id: 228, chars: []rune{93}}
	p229.items = []parser{&p228}
	p234.items = []parser{&p266, &p835, &p227, &p231, &p835, &p394, &p233, &p835, &p229}
	var p243 = sequenceParser{id: 243, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{244, 266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p236 = sequenceParser{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p235 = charParser{id: 235, chars: []rune{91}}
	p236.items = []parser{&p235}
	var p240 = sequenceParser{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p239 = sequenceParser{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p239.items = []parser{&p835, &p14}
	p240.items = []parser{&p835, &p14, &p239}
	var p225 = sequenceParser{id: 225, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{569, 575, 576}}
	var p217 = sequenceParser{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p217.items = []parser{&p394}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p221 = sequenceParser{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p221.items = []parser{&p835, &p14}
	p222.items = []parser{&p835, &p14, &p221}
	var p220 = sequenceParser{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p219 = charParser{id: 219, chars: []rune{58}}
	p220.items = []parser{&p219}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p835, &p14}
	p224.items = []parser{&p835, &p14, &p223}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p394}
	p225.items = []parser{&p217, &p222, &p835, &p220, &p224, &p835, &p218}
	var p242 = sequenceParser{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p241 = sequenceParser{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p241.items = []parser{&p835, &p14}
	p242.items = []parser{&p835, &p14, &p241}
	var p238 = sequenceParser{id: 238, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p237 = charParser{id: 237, chars: []rune{93}}
	p238.items = []parser{&p237}
	p243.items = []parser{&p266, &p835, &p236, &p240, &p835, &p225, &p242, &p835, &p238}
	p244.options = []parser{&p234, &p243}
	var p251 = sequenceParser{id: 251, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p247.items = []parser{&p835, &p14}
	p248.items = []parser{&p835, &p14, &p247}
	var p246 = sequenceParser{id: 246, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p245 = charParser{id: 245, chars: []rune{46}}
	p246.items = []parser{&p245}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p249.items = []parser{&p835, &p14}
	p250.items = []parser{&p835, &p14, &p249}
	p251.items = []parser{&p266, &p248, &p835, &p246, &p250, &p835, &p103}
	var p265 = sequenceParser{id: 265, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p258 = sequenceParser{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p257 = charParser{id: 257, chars: []rune{40}}
	p258.items = []parser{&p257}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p261 = sequenceParser{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p261.items = []parser{&p835, &p14}
	p262.items = []parser{&p835, &p14, &p261}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p263.items = []parser{&p835, &p14}
	p264.items = []parser{&p835, &p14, &p263}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{41}}
	p260.items = []parser{&p259}
	p265.items = []parser{&p258, &p262, &p835, &p394, &p264, &p835, &p260}
	p266.options = []parser{&p60, &p73, &p86, &p98, &p499, &p103, &p124, &p129, &p158, &p163, &p206, &p216, &p244, &p251, &p256, &p265}
	var p326 = sequenceParser{id: 326, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{394, 332, 333, 334, 335, 336, 337, 576, 569, 803}}
	var p325 = choiceParser{id: 325, commit: 66, name: "unary-operator"}
	var p285 = sequenceParser{id: 285, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{325}}
	var p284 = charParser{id: 284, chars: []rune{43}}
	p285.items = []parser{&p284}
	var p287 = sequenceParser{id: 287, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{325}}
	var p286 = charParser{id: 286, chars: []rune{45}}
	p287.items = []parser{&p286}
	var p268 = sequenceParser{id: 268, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{325}}
	var p267 = charParser{id: 267, chars: []rune{94}}
	p268.items = []parser{&p267}
	var p299 = sequenceParser{id: 299, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{325}}
	var p298 = charParser{id: 298, chars: []rune{33}}
	p299.items = []parser{&p298}
	p325.options = []parser{&p285, &p287, &p268, &p299}
	p326.items = []parser{&p325, &p835, &p266}
	var p380 = choiceParser{id: 380, commit: 66, name: "binary-expression", generalizations: []int{394, 576, 569, 803}}
	var p344 = sequenceParser{id: 344, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 333, 334, 335, 336, 337, 394, 576, 569, 803}}
	var p332 = choiceParser{id: 332, commit: 66, name: "operand0", generalizations: []int{333, 334, 335, 336, 337}}
	p332.options = []parser{&p266, &p326}
	var p342 = sequenceParser{id: 342, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p339 = sequenceParser{id: 339, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p338 = sequenceParser{id: 338, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p338.items = []parser{&p835, &p14}
	p339.items = []parser{&p14, &p338}
	var p327 = choiceParser{id: 327, commit: 66, name: "binary-op0"}
	var p270 = sequenceParser{id: 270, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p269 = charParser{id: 269, chars: []rune{38}}
	p270.items = []parser{&p269}
	var p277 = sequenceParser{id: 277, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{327}}
	var p275 = charParser{id: 275, chars: []rune{38}}
	var p276 = charParser{id: 276, chars: []rune{94}}
	p277.items = []parser{&p275, &p276}
	var p280 = sequenceParser{id: 280, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{327}}
	var p278 = charParser{id: 278, chars: []rune{60}}
	var p279 = charParser{id: 279, chars: []rune{60}}
	p280.items = []parser{&p278, &p279}
	var p283 = sequenceParser{id: 283, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{327}}
	var p281 = charParser{id: 281, chars: []rune{62}}
	var p282 = charParser{id: 282, chars: []rune{62}}
	p283.items = []parser{&p281, &p282}
	var p289 = sequenceParser{id: 289, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p288 = charParser{id: 288, chars: []rune{42}}
	p289.items = []parser{&p288}
	var p291 = sequenceParser{id: 291, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p290 = charParser{id: 290, chars: []rune{47}}
	p291.items = []parser{&p290}
	var p293 = sequenceParser{id: 293, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p292 = charParser{id: 292, chars: []rune{37}}
	p293.items = []parser{&p292}
	p327.options = []parser{&p270, &p277, &p280, &p283, &p289, &p291, &p293}
	var p341 = sequenceParser{id: 341, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p340 = sequenceParser{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p340.items = []parser{&p835, &p14}
	p341.items = []parser{&p835, &p14, &p340}
	p342.items = []parser{&p339, &p835, &p327, &p341, &p835, &p332}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p343.items = []parser{&p835, &p342}
	p344.items = []parser{&p332, &p835, &p342, &p343}
	var p351 = sequenceParser{id: 351, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 334, 335, 336, 337, 394, 576, 569, 803}}
	var p333 = choiceParser{id: 333, commit: 66, name: "operand1", generalizations: []int{334, 335, 336, 337}}
	p333.options = []parser{&p332, &p344}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p345 = sequenceParser{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p345.items = []parser{&p835, &p14}
	p346.items = []parser{&p14, &p345}
	var p328 = choiceParser{id: 328, commit: 66, name: "binary-op1"}
	var p272 = sequenceParser{id: 272, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p271 = charParser{id: 271, chars: []rune{124}}
	p272.items = []parser{&p271}
	var p274 = sequenceParser{id: 274, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p273 = charParser{id: 273, chars: []rune{94}}
	p274.items = []parser{&p273}
	var p295 = sequenceParser{id: 295, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p294 = charParser{id: 294, chars: []rune{43}}
	p295.items = []parser{&p294}
	var p297 = sequenceParser{id: 297, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{328}}
	var p296 = charParser{id: 296, chars: []rune{45}}
	p297.items = []parser{&p296}
	p328.options = []parser{&p272, &p274, &p295, &p297}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p347.items = []parser{&p835, &p14}
	p348.items = []parser{&p835, &p14, &p347}
	p349.items = []parser{&p346, &p835, &p328, &p348, &p835, &p333}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p350.items = []parser{&p835, &p349}
	p351.items = []parser{&p333, &p835, &p349, &p350}
	var p358 = sequenceParser{id: 358, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 335, 336, 337, 394, 576, 569, 803}}
	var p334 = choiceParser{id: 334, commit: 66, name: "operand2", generalizations: []int{335, 336, 337}}
	p334.options = []parser{&p333, &p351}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p352.items = []parser{&p835, &p14}
	p353.items = []parser{&p14, &p352}
	var p329 = choiceParser{id: 329, commit: 66, name: "binary-op2"}
	var p302 = sequenceParser{id: 302, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p300 = charParser{id: 300, chars: []rune{61}}
	var p301 = charParser{id: 301, chars: []rune{61}}
	p302.items = []parser{&p300, &p301}
	var p305 = sequenceParser{id: 305, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p303 = charParser{id: 303, chars: []rune{33}}
	var p304 = charParser{id: 304, chars: []rune{61}}
	p305.items = []parser{&p303, &p304}
	var p307 = sequenceParser{id: 307, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p306 = charParser{id: 306, chars: []rune{60}}
	p307.items = []parser{&p306}
	var p310 = sequenceParser{id: 310, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p308 = charParser{id: 308, chars: []rune{60}}
	var p309 = charParser{id: 309, chars: []rune{61}}
	p310.items = []parser{&p308, &p309}
	var p312 = sequenceParser{id: 312, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p311 = charParser{id: 311, chars: []rune{62}}
	p312.items = []parser{&p311}
	var p315 = sequenceParser{id: 315, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p313 = charParser{id: 313, chars: []rune{62}}
	var p314 = charParser{id: 314, chars: []rune{61}}
	p315.items = []parser{&p313, &p314}
	p329.options = []parser{&p302, &p305, &p307, &p310, &p312, &p315}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p354.items = []parser{&p835, &p14}
	p355.items = []parser{&p835, &p14, &p354}
	p356.items = []parser{&p353, &p835, &p329, &p355, &p835, &p334}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p357.items = []parser{&p835, &p356}
	p358.items = []parser{&p334, &p835, &p356, &p357}
	var p365 = sequenceParser{id: 365, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 336, 337, 394, 576, 569, 803}}
	var p335 = choiceParser{id: 335, commit: 66, name: "operand3", generalizations: []int{336, 337}}
	p335.options = []parser{&p334, &p358}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p835, &p14}
	p360.items = []parser{&p14, &p359}
	var p330 = sequenceParser{id: 330, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p318 = sequenceParser{id: 318, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p316 = charParser{id: 316, chars: []rune{38}}
	var p317 = charParser{id: 317, chars: []rune{38}}
	p318.items = []parser{&p316, &p317}
	p330.items = []parser{&p318}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p361.items = []parser{&p835, &p14}
	p362.items = []parser{&p835, &p14, &p361}
	p363.items = []parser{&p360, &p835, &p330, &p362, &p835, &p335}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p364.items = []parser{&p835, &p363}
	p365.items = []parser{&p335, &p835, &p363, &p364}
	var p372 = sequenceParser{id: 372, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 337, 394, 576, 569, 803}}
	var p336 = choiceParser{id: 336, commit: 66, name: "operand4", generalizations: []int{337}}
	p336.options = []parser{&p335, &p365}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p366.items = []parser{&p835, &p14}
	p367.items = []parser{&p14, &p366}
	var p331 = sequenceParser{id: 331, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p321 = sequenceParser{id: 321, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p319 = charParser{id: 319, chars: []rune{124}}
	var p320 = charParser{id: 320, chars: []rune{124}}
	p321.items = []parser{&p319, &p320}
	p331.items = []parser{&p321}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p368.items = []parser{&p835, &p14}
	p369.items = []parser{&p835, &p14, &p368}
	p370.items = []parser{&p367, &p835, &p331, &p369, &p835, &p336}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p371.items = []parser{&p835, &p370}
	p372.items = []parser{&p336, &p835, &p370, &p371}
	var p379 = sequenceParser{id: 379, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 394, 576, 569, 803}}
	var p337 = choiceParser{id: 337, commit: 66, name: "operand5"}
	p337.options = []parser{&p336, &p372}
	var p377 = sequenceParser{id: 377, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p374 = sequenceParser{id: 374, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p373.items = []parser{&p835, &p14}
	p374.items = []parser{&p14, &p373}
	var p324 = sequenceParser{id: 324, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p322 = charParser{id: 322, chars: []rune{45}}
	var p323 = charParser{id: 323, chars: []rune{62}}
	p324.items = []parser{&p322, &p323}
	var p376 = sequenceParser{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p375 = sequenceParser{id: 375, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p375.items = []parser{&p835, &p14}
	p376.items = []parser{&p835, &p14, &p375}
	p377.items = []parser{&p374, &p835, &p324, &p376, &p835, &p337}
	var p378 = sequenceParser{id: 378, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p378.items = []parser{&p835, &p377}
	p379.items = []parser{&p337, &p835, &p377, &p378}
	p380.options = []parser{&p344, &p351, &p358, &p365, &p372, &p379}
	var p393 = sequenceParser{id: 393, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{394, 576, 569, 803}}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p385.items = []parser{&p835, &p14}
	p386.items = []parser{&p835, &p14, &p385}
	var p382 = sequenceParser{id: 382, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p381 = charParser{id: 381, chars: []rune{63}}
	p382.items = []parser{&p381}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p387.items = []parser{&p835, &p14}
	p388.items = []parser{&p835, &p14, &p387}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p389.items = []parser{&p835, &p14}
	p390.items = []parser{&p835, &p14, &p389}
	var p384 = sequenceParser{id: 384, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p383 = charParser{id: 383, chars: []rune{58}}
	p384.items = []parser{&p383}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p391.items = []parser{&p835, &p14}
	p392.items = []parser{&p835, &p14, &p391}
	p393.items = []parser{&p394, &p386, &p835, &p382, &p388, &p835, &p394, &p390, &p835, &p384, &p392, &p835, &p394}
	p394.options = []parser{&p266, &p326, &p380, &p393}
	p184.items = []parser{&p183, &p835, &p394}
	p185.items = []parser{&p181, &p835, &p184}
	var p431 = sequenceParser{id: 431, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{803, 477, 541}}
	var p397 = sequenceParser{id: 397, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p395 = charParser{id: 395, chars: []rune{105}}
	var p396 = charParser{id: 396, chars: []rune{102}}
	p397.items = []parser{&p395, &p396}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p425.items = []parser{&p835, &p14}
	p426.items = []parser{&p835, &p14, &p425}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p427.items = []parser{&p835, &p14}
	p428.items = []parser{&p835, &p14, &p427}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p407 = sequenceParser{id: 407, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p406 = sequenceParser{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p406.items = []parser{&p835, &p14}
	p407.items = []parser{&p14, &p406}
	var p402 = sequenceParser{id: 402, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p398 = charParser{id: 398, chars: []rune{101}}
	var p399 = charParser{id: 399, chars: []rune{108}}
	var p400 = charParser{id: 400, chars: []rune{115}}
	var p401 = charParser{id: 401, chars: []rune{101}}
	p402.items = []parser{&p398, &p399, &p400, &p401}
	var p409 = sequenceParser{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p408.items = []parser{&p835, &p14}
	p409.items = []parser{&p835, &p14, &p408}
	var p405 = sequenceParser{id: 405, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p403 = charParser{id: 403, chars: []rune{105}}
	var p404 = charParser{id: 404, chars: []rune{102}}
	p405.items = []parser{&p403, &p404}
	var p411 = sequenceParser{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p410 = sequenceParser{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p410.items = []parser{&p835, &p14}
	p411.items = []parser{&p835, &p14, &p410}
	var p413 = sequenceParser{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p412 = sequenceParser{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p412.items = []parser{&p835, &p14}
	p413.items = []parser{&p835, &p14, &p412}
	p414.items = []parser{&p407, &p835, &p402, &p409, &p835, &p405, &p411, &p835, &p394, &p413, &p835, &p190}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p429.items = []parser{&p835, &p414}
	p430.items = []parser{&p835, &p414, &p429}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p421 = sequenceParser{id: 421, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p420 = sequenceParser{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p420.items = []parser{&p835, &p14}
	p421.items = []parser{&p14, &p420}
	var p419 = sequenceParser{id: 419, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p415 = charParser{id: 415, chars: []rune{101}}
	var p416 = charParser{id: 416, chars: []rune{108}}
	var p417 = charParser{id: 417, chars: []rune{115}}
	var p418 = charParser{id: 418, chars: []rune{101}}
	p419.items = []parser{&p415, &p416, &p417, &p418}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p422.items = []parser{&p835, &p14}
	p423.items = []parser{&p835, &p14, &p422}
	p424.items = []parser{&p421, &p835, &p419, &p423, &p835, &p190}
	p431.items = []parser{&p397, &p426, &p835, &p394, &p428, &p835, &p190, &p430, &p835, &p424}
	var p488 = sequenceParser{id: 488, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{477, 803, 541}}
	var p473 = sequenceParser{id: 473, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p467 = charParser{id: 467, chars: []rune{115}}
	var p468 = charParser{id: 468, chars: []rune{119}}
	var p469 = charParser{id: 469, chars: []rune{105}}
	var p470 = charParser{id: 470, chars: []rune{116}}
	var p471 = charParser{id: 471, chars: []rune{99}}
	var p472 = charParser{id: 472, chars: []rune{104}}
	p473.items = []parser{&p467, &p468, &p469, &p470, &p471, &p472}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p484 = sequenceParser{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p484.items = []parser{&p835, &p14}
	p485.items = []parser{&p835, &p14, &p484}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p486.items = []parser{&p835, &p14}
	p487.items = []parser{&p835, &p14, &p486}
	var p475 = sequenceParser{id: 475, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p474 = charParser{id: 474, chars: []rune{123}}
	p475.items = []parser{&p474}
	var p481 = sequenceParser{id: 481, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p476 = choiceParser{id: 476, commit: 2}
	var p466 = sequenceParser{id: 466, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{476, 477}}
	var p461 = sequenceParser{id: 461, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p454 = sequenceParser{id: 454, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p450 = charParser{id: 450, chars: []rune{99}}
	var p451 = charParser{id: 451, chars: []rune{97}}
	var p452 = charParser{id: 452, chars: []rune{115}}
	var p453 = charParser{id: 453, chars: []rune{101}}
	p454.items = []parser{&p450, &p451, &p452, &p453}
	var p458 = sequenceParser{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p457 = sequenceParser{id: 457, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p457.items = []parser{&p835, &p14}
	p458.items = []parser{&p835, &p14, &p457}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p459.items = []parser{&p835, &p14}
	p460.items = []parser{&p835, &p14, &p459}
	var p456 = sequenceParser{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p455 = charParser{id: 455, chars: []rune{58}}
	p456.items = []parser{&p455}
	p461.items = []parser{&p454, &p458, &p835, &p394, &p460, &p835, &p456}
	var p465 = sequenceParser{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p463 = sequenceParser{id: 463, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p462 = charParser{id: 462, chars: []rune{59}}
	p463.items = []parser{&p462}
	var p464 = sequenceParser{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p464.items = []parser{&p835, &p463}
	p465.items = []parser{&p835, &p463, &p464}
	p466.items = []parser{&p461, &p465, &p835, &p803}
	var p449 = sequenceParser{id: 449, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{476, 477, 540, 541}}
	var p444 = sequenceParser{id: 444, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p439 = sequenceParser{id: 439, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p432 = charParser{id: 432, chars: []rune{100}}
	var p433 = charParser{id: 433, chars: []rune{101}}
	var p434 = charParser{id: 434, chars: []rune{102}}
	var p435 = charParser{id: 435, chars: []rune{97}}
	var p436 = charParser{id: 436, chars: []rune{117}}
	var p437 = charParser{id: 437, chars: []rune{108}}
	var p438 = charParser{id: 438, chars: []rune{116}}
	p439.items = []parser{&p432, &p433, &p434, &p435, &p436, &p437, &p438}
	var p443 = sequenceParser{id: 443, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p442 = sequenceParser{id: 442, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p442.items = []parser{&p835, &p14}
	p443.items = []parser{&p835, &p14, &p442}
	var p441 = sequenceParser{id: 441, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p440 = charParser{id: 440, chars: []rune{58}}
	p441.items = []parser{&p440}
	p444.items = []parser{&p439, &p443, &p835, &p441}
	var p448 = sequenceParser{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p446 = sequenceParser{id: 446, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p445 = charParser{id: 445, chars: []rune{59}}
	p446.items = []parser{&p445}
	var p447 = sequenceParser{id: 447, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p447.items = []parser{&p835, &p446}
	p448.items = []parser{&p835, &p446, &p447}
	p449.items = []parser{&p444, &p448, &p835, &p803}
	p476.options = []parser{&p466, &p449}
	var p480 = sequenceParser{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p478 = sequenceParser{id: 478, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p477 = choiceParser{id: 477, commit: 2}
	p477.options = []parser{&p466, &p449, &p803}
	p478.items = []parser{&p817, &p835, &p477}
	var p479 = sequenceParser{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p479.items = []parser{&p835, &p478}
	p480.items = []parser{&p835, &p478, &p479}
	p481.items = []parser{&p476, &p480}
	var p483 = sequenceParser{id: 483, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p482 = charParser{id: 482, chars: []rune{125}}
	p483.items = []parser{&p482}
	p488.items = []parser{&p473, &p485, &p835, &p394, &p487, &p835, &p475, &p835, &p817, &p835, &p481, &p835, &p817, &p835, &p483}
	var p550 = sequenceParser{id: 550, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{541, 803}}
	var p537 = sequenceParser{id: 537, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p531 = charParser{id: 531, chars: []rune{115}}
	var p532 = charParser{id: 532, chars: []rune{101}}
	var p533 = charParser{id: 533, chars: []rune{108}}
	var p534 = charParser{id: 534, chars: []rune{101}}
	var p535 = charParser{id: 535, chars: []rune{99}}
	var p536 = charParser{id: 536, chars: []rune{116}}
	p537.items = []parser{&p531, &p532, &p533, &p534, &p535, &p536}
	var p549 = sequenceParser{id: 549, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p548 = sequenceParser{id: 548, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p548.items = []parser{&p835, &p14}
	p549.items = []parser{&p835, &p14, &p548}
	var p539 = sequenceParser{id: 539, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p538 = charParser{id: 538, chars: []rune{123}}
	p539.items = []parser{&p538}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p540 = choiceParser{id: 540, commit: 2}
	var p530 = sequenceParser{id: 530, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{540, 541}}
	var p525 = sequenceParser{id: 525, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p518 = sequenceParser{id: 518, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p514 = charParser{id: 514, chars: []rune{99}}
	var p515 = charParser{id: 515, chars: []rune{97}}
	var p516 = charParser{id: 516, chars: []rune{115}}
	var p517 = charParser{id: 517, chars: []rune{101}}
	p518.items = []parser{&p514, &p515, &p516, &p517}
	var p522 = sequenceParser{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p521 = sequenceParser{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p521.items = []parser{&p835, &p14}
	p522.items = []parser{&p835, &p14, &p521}
	var p513 = choiceParser{id: 513, commit: 66, name: "communication"}
	var p512 = sequenceParser{id: 512, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{513}}
	var p511 = sequenceParser{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p510.items = []parser{&p835, &p14}
	p511.items = []parser{&p835, &p14, &p510}
	p512.items = []parser{&p103, &p511, &p835, &p499}
	p513.options = []parser{&p499, &p512, &p509}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p523.items = []parser{&p835, &p14}
	p524.items = []parser{&p835, &p14, &p523}
	var p520 = sequenceParser{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p519 = charParser{id: 519, chars: []rune{58}}
	p520.items = []parser{&p519}
	p525.items = []parser{&p518, &p522, &p835, &p513, &p524, &p835, &p520}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p527 = sequenceParser{id: 527, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p526 = charParser{id: 526, chars: []rune{59}}
	p527.items = []parser{&p526}
	var p528 = sequenceParser{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p528.items = []parser{&p835, &p527}
	p529.items = []parser{&p835, &p527, &p528}
	p530.items = []parser{&p525, &p529, &p835, &p803}
	p540.options = []parser{&p530, &p449}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p542 = sequenceParser{id: 542, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p541 = choiceParser{id: 541, commit: 2}
	p541.options = []parser{&p530, &p449, &p803}
	p542.items = []parser{&p817, &p835, &p541}
	var p543 = sequenceParser{id: 543, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p543.items = []parser{&p835, &p542}
	p544.items = []parser{&p835, &p542, &p543}
	p545.items = []parser{&p540, &p544}
	var p547 = sequenceParser{id: 547, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p546 = charParser{id: 546, chars: []rune{125}}
	p547.items = []parser{&p546}
	p550.items = []parser{&p537, &p549, &p835, &p539, &p835, &p817, &p835, &p545, &p835, &p817, &p835, &p547}
	var p591 = sequenceParser{id: 591, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{803}}
	var p580 = sequenceParser{id: 580, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p577 = charParser{id: 577, chars: []rune{102}}
	var p578 = charParser{id: 578, chars: []rune{111}}
	var p579 = charParser{id: 579, chars: []rune{114}}
	p580.items = []parser{&p577, &p578, &p579}
	var p590 = choiceParser{id: 590, commit: 2}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{590}}
	var p583 = sequenceParser{id: 583, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p582 = sequenceParser{id: 582, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p581 = sequenceParser{id: 581, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p581.items = []parser{&p835, &p14}
	p582.items = []parser{&p14, &p581}
	var p576 = choiceParser{id: 576, commit: 66, name: "loop-expression"}
	var p575 = choiceParser{id: 575, commit: 64, name: "range-over-expression", generalizations: []int{576}}
	var p574 = sequenceParser{id: 574, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{575, 576}}
	var p571 = sequenceParser{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p570 = sequenceParser{id: 570, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p570.items = []parser{&p835, &p14}
	p571.items = []parser{&p835, &p14, &p570}
	var p568 = sequenceParser{id: 568, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p566 = charParser{id: 566, chars: []rune{105}}
	var p567 = charParser{id: 567, chars: []rune{110}}
	p568.items = []parser{&p566, &p567}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p572.items = []parser{&p835, &p14}
	p573.items = []parser{&p835, &p14, &p572}
	var p569 = choiceParser{id: 569, commit: 2}
	p569.options = []parser{&p394, &p225}
	p574.items = []parser{&p103, &p571, &p835, &p568, &p573, &p835, &p569}
	p575.options = []parser{&p574, &p225}
	p576.options = []parser{&p394, &p575}
	p583.items = []parser{&p582, &p835, &p576}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p584.items = []parser{&p835, &p14}
	p585.items = []parser{&p835, &p14, &p584}
	p586.items = []parser{&p583, &p585, &p835, &p190}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{590}}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p587.items = []parser{&p835, &p14}
	p588.items = []parser{&p14, &p587}
	p589.items = []parser{&p588, &p835, &p190}
	p590.options = []parser{&p586, &p589}
	p591.items = []parser{&p580, &p835, &p590}
	var p739 = choiceParser{id: 739, commit: 66, name: "definition", generalizations: []int{803}}
	var p652 = sequenceParser{id: 652, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 803}}
	var p648 = sequenceParser{id: 648, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p645 = charParser{id: 645, chars: []rune{108}}
	var p646 = charParser{id: 646, chars: []rune{101}}
	var p647 = charParser{id: 647, chars: []rune{116}}
	p648.items = []parser{&p645, &p646, &p647}
	var p651 = sequenceParser{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p650 = sequenceParser{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p650.items = []parser{&p835, &p14}
	p651.items = []parser{&p835, &p14, &p650}
	var p649 = choiceParser{id: 649, commit: 2}
	var p639 = sequenceParser{id: 639, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{649, 653, 654}}
	var p638 = sequenceParser{id: 638, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p633 = sequenceParser{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p633.items = []parser{&p835, &p14}
	p634.items = []parser{&p14, &p633}
	var p632 = sequenceParser{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p631 = charParser{id: 631, chars: []rune{61}}
	p632.items = []parser{&p631}
	p635.items = []parser{&p634, &p835, &p632}
	var p637 = sequenceParser{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p636.items = []parser{&p835, &p14}
	p637.items = []parser{&p835, &p14, &p636}
	p638.items = []parser{&p103, &p835, &p635, &p637, &p835, &p394}
	p639.items = []parser{&p638}
	var p644 = sequenceParser{id: 644, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{649, 653, 654}}
	var p641 = sequenceParser{id: 641, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p640 = charParser{id: 640, chars: []rune{126}}
	p641.items = []parser{&p640}
	var p643 = sequenceParser{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p642.items = []parser{&p835, &p14}
	p643.items = []parser{&p835, &p14, &p642}
	p644.items = []parser{&p641, &p643, &p835, &p638}
	p649.options = []parser{&p639, &p644}
	p652.items = []parser{&p648, &p651, &p835, &p649}
	var p673 = sequenceParser{id: 673, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 803}}
	var p666 = sequenceParser{id: 666, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p663 = charParser{id: 663, chars: []rune{108}}
	var p664 = charParser{id: 664, chars: []rune{101}}
	var p665 = charParser{id: 665, chars: []rune{116}}
	p666.items = []parser{&p663, &p664, &p665}
	var p672 = sequenceParser{id: 672, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p671 = sequenceParser{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p671.items = []parser{&p835, &p14}
	p672.items = []parser{&p835, &p14, &p671}
	var p668 = sequenceParser{id: 668, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p667 = charParser{id: 667, chars: []rune{40}}
	p668.items = []parser{&p667}
	var p658 = sequenceParser{id: 658, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p653 = choiceParser{id: 653, commit: 2}
	p653.options = []parser{&p639, &p644}
	var p657 = sequenceParser{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p655 = sequenceParser{id: 655, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p654 = choiceParser{id: 654, commit: 2}
	p654.options = []parser{&p639, &p644}
	p655.items = []parser{&p113, &p835, &p654}
	var p656 = sequenceParser{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p656.items = []parser{&p835, &p655}
	p657.items = []parser{&p835, &p655, &p656}
	p658.items = []parser{&p653, &p657}
	var p670 = sequenceParser{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p669 = charParser{id: 669, chars: []rune{41}}
	p670.items = []parser{&p669}
	p673.items = []parser{&p666, &p672, &p835, &p668, &p835, &p113, &p835, &p658, &p835, &p113, &p835, &p670}
	var p688 = sequenceParser{id: 688, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 803}}
	var p677 = sequenceParser{id: 677, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p674 = charParser{id: 674, chars: []rune{108}}
	var p675 = charParser{id: 675, chars: []rune{101}}
	var p676 = charParser{id: 676, chars: []rune{116}}
	p677.items = []parser{&p674, &p675, &p676}
	var p685 = sequenceParser{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p684 = sequenceParser{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p684.items = []parser{&p835, &p14}
	p685.items = []parser{&p835, &p14, &p684}
	var p679 = sequenceParser{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p678 = charParser{id: 678, chars: []rune{126}}
	p679.items = []parser{&p678}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p686.items = []parser{&p835, &p14}
	p687.items = []parser{&p835, &p14, &p686}
	var p681 = sequenceParser{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{40}}
	p681.items = []parser{&p680}
	var p662 = sequenceParser{id: 662, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p661 = sequenceParser{id: 661, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p659 = sequenceParser{id: 659, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p659.items = []parser{&p113, &p835, &p639}
	var p660 = sequenceParser{id: 660, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p660.items = []parser{&p835, &p659}
	p661.items = []parser{&p835, &p659, &p660}
	p662.items = []parser{&p639, &p661}
	var p683 = sequenceParser{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p682 = charParser{id: 682, chars: []rune{41}}
	p683.items = []parser{&p682}
	p688.items = []parser{&p677, &p685, &p835, &p679, &p687, &p835, &p681, &p835, &p113, &p835, &p662, &p835, &p113, &p835, &p683}
	var p704 = sequenceParser{id: 704, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 803}}
	var p700 = sequenceParser{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p698 = charParser{id: 698, chars: []rune{102}}
	var p699 = charParser{id: 699, chars: []rune{110}}
	p700.items = []parser{&p698, &p699}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p702 = sequenceParser{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p702.items = []parser{&p835, &p14}
	p703.items = []parser{&p835, &p14, &p702}
	var p701 = choiceParser{id: 701, commit: 2}
	var p692 = sequenceParser{id: 692, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{701, 709, 710}}
	var p691 = sequenceParser{id: 691, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p690 = sequenceParser{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p689 = sequenceParser{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p689.items = []parser{&p835, &p14}
	p690.items = []parser{&p835, &p14, &p689}
	p691.items = []parser{&p103, &p690, &p835, &p200}
	p692.items = []parser{&p691}
	var p697 = sequenceParser{id: 697, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{701, 709, 710}}
	var p694 = sequenceParser{id: 694, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p693 = charParser{id: 693, chars: []rune{126}}
	p694.items = []parser{&p693}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p695 = sequenceParser{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p695.items = []parser{&p835, &p14}
	p696.items = []parser{&p835, &p14, &p695}
	p697.items = []parser{&p694, &p696, &p835, &p691}
	p701.options = []parser{&p692, &p697}
	p704.items = []parser{&p700, &p703, &p835, &p701}
	var p724 = sequenceParser{id: 724, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 803}}
	var p717 = sequenceParser{id: 717, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p715 = charParser{id: 715, chars: []rune{102}}
	var p716 = charParser{id: 716, chars: []rune{110}}
	p717.items = []parser{&p715, &p716}
	var p723 = sequenceParser{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p722 = sequenceParser{id: 722, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p722.items = []parser{&p835, &p14}
	p723.items = []parser{&p835, &p14, &p722}
	var p719 = sequenceParser{id: 719, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p718 = charParser{id: 718, chars: []rune{40}}
	p719.items = []parser{&p718}
	var p714 = sequenceParser{id: 714, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p709 = choiceParser{id: 709, commit: 2}
	p709.options = []parser{&p692, &p697}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p711 = sequenceParser{id: 711, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p710 = choiceParser{id: 710, commit: 2}
	p710.options = []parser{&p692, &p697}
	p711.items = []parser{&p113, &p835, &p710}
	var p712 = sequenceParser{id: 712, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p712.items = []parser{&p835, &p711}
	p713.items = []parser{&p835, &p711, &p712}
	p714.items = []parser{&p709, &p713}
	var p721 = sequenceParser{id: 721, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p720 = charParser{id: 720, chars: []rune{41}}
	p721.items = []parser{&p720}
	p724.items = []parser{&p717, &p723, &p835, &p719, &p835, &p113, &p835, &p714, &p835, &p113, &p835, &p721}
	var p738 = sequenceParser{id: 738, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 803}}
	var p727 = sequenceParser{id: 727, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p725 = charParser{id: 725, chars: []rune{102}}
	var p726 = charParser{id: 726, chars: []rune{110}}
	p727.items = []parser{&p725, &p726}
	var p735 = sequenceParser{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p734 = sequenceParser{id: 734, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p734.items = []parser{&p835, &p14}
	p735.items = []parser{&p835, &p14, &p734}
	var p729 = sequenceParser{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p728 = charParser{id: 728, chars: []rune{126}}
	p729.items = []parser{&p728}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p736.items = []parser{&p835, &p14}
	p737.items = []parser{&p835, &p14, &p736}
	var p731 = sequenceParser{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p730 = charParser{id: 730, chars: []rune{40}}
	p731.items = []parser{&p730}
	var p708 = sequenceParser{id: 708, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p707 = sequenceParser{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p705 = sequenceParser{id: 705, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p705.items = []parser{&p113, &p835, &p692}
	var p706 = sequenceParser{id: 706, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p706.items = []parser{&p835, &p705}
	p707.items = []parser{&p835, &p705, &p706}
	p708.items = []parser{&p692, &p707}
	var p733 = sequenceParser{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p732 = charParser{id: 732, chars: []rune{41}}
	p733.items = []parser{&p732}
	p738.items = []parser{&p727, &p735, &p835, &p729, &p737, &p835, &p731, &p835, &p113, &p835, &p708, &p835, &p113, &p835, &p733}
	p739.options = []parser{&p652, &p673, &p688, &p704, &p724, &p738}
	var p782 = choiceParser{id: 782, commit: 64, name: "require", generalizations: []int{803}}
	var p766 = sequenceParser{id: 766, commit: 66, name: "require-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{782, 803}}
	var p763 = sequenceParser{id: 763, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p756 = charParser{id: 756, chars: []rune{114}}
	var p757 = charParser{id: 757, chars: []rune{101}}
	var p758 = charParser{id: 758, chars: []rune{113}}
	var p759 = charParser{id: 759, chars: []rune{117}}
	var p760 = charParser{id: 760, chars: []rune{105}}
	var p761 = charParser{id: 761, chars: []rune{114}}
	var p762 = charParser{id: 762, chars: []rune{101}}
	p763.items = []parser{&p756, &p757, &p758, &p759, &p760, &p761, &p762}
	var p765 = sequenceParser{id: 765, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p764 = sequenceParser{id: 764, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p764.items = []parser{&p835, &p14}
	p765.items = []parser{&p835, &p14, &p764}
	var p751 = choiceParser{id: 751, commit: 64, name: "require-fact"}
	var p750 = sequenceParser{id: 750, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{751}}
	var p742 = choiceParser{id: 742, commit: 2}
	var p741 = sequenceParser{id: 741, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{742}}
	var p740 = charParser{id: 740, chars: []rune{46}}
	p741.items = []parser{&p740}
	p742.options = []parser{&p103, &p741}
	var p747 = sequenceParser{id: 747, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p746 = sequenceParser{id: 746, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p745 = sequenceParser{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p745.items = []parser{&p835, &p14}
	p746.items = []parser{&p14, &p745}
	var p744 = sequenceParser{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p743 = charParser{id: 743, chars: []rune{61}}
	p744.items = []parser{&p743}
	p747.items = []parser{&p746, &p835, &p744}
	var p749 = sequenceParser{id: 749, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p748 = sequenceParser{id: 748, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p748.items = []parser{&p835, &p14}
	p749.items = []parser{&p835, &p14, &p748}
	p750.items = []parser{&p742, &p835, &p747, &p749, &p835, &p86}
	p751.options = []parser{&p86, &p750}
	p766.items = []parser{&p763, &p765, &p835, &p751}
	var p781 = sequenceParser{id: 781, commit: 66, name: "require-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{782, 803}}
	var p774 = sequenceParser{id: 774, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p767 = charParser{id: 767, chars: []rune{114}}
	var p768 = charParser{id: 768, chars: []rune{101}}
	var p769 = charParser{id: 769, chars: []rune{113}}
	var p770 = charParser{id: 770, chars: []rune{117}}
	var p771 = charParser{id: 771, chars: []rune{105}}
	var p772 = charParser{id: 772, chars: []rune{114}}
	var p773 = charParser{id: 773, chars: []rune{101}}
	p774.items = []parser{&p767, &p768, &p769, &p770, &p771, &p772, &p773}
	var p780 = sequenceParser{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p779 = sequenceParser{id: 779, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p779.items = []parser{&p835, &p14}
	p780.items = []parser{&p835, &p14, &p779}
	var p776 = sequenceParser{id: 776, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p775 = charParser{id: 775, chars: []rune{40}}
	p776.items = []parser{&p775}
	var p755 = sequenceParser{id: 755, commit: 66, name: "require-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p752.items = []parser{&p113, &p835, &p751}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p753.items = []parser{&p835, &p752}
	p754.items = []parser{&p835, &p752, &p753}
	p755.items = []parser{&p751, &p754}
	var p778 = sequenceParser{id: 778, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p777 = charParser{id: 777, chars: []rune{41}}
	p778.items = []parser{&p777}
	p781.items = []parser{&p774, &p780, &p835, &p776, &p835, &p113, &p835, &p755, &p835, &p113, &p835, &p778}
	p782.options = []parser{&p766, &p781}
	var p792 = sequenceParser{id: 792, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{803}}
	var p789 = sequenceParser{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p783 = charParser{id: 783, chars: []rune{101}}
	var p784 = charParser{id: 784, chars: []rune{120}}
	var p785 = charParser{id: 785, chars: []rune{112}}
	var p786 = charParser{id: 786, chars: []rune{111}}
	var p787 = charParser{id: 787, chars: []rune{114}}
	var p788 = charParser{id: 788, chars: []rune{116}}
	p789.items = []parser{&p783, &p784, &p785, &p786, &p787, &p788}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p790.items = []parser{&p835, &p14}
	p791.items = []parser{&p835, &p14, &p790}
	p792.items = []parser{&p789, &p791, &p835, &p739}
	var p812 = sequenceParser{id: 812, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{803}}
	var p805 = sequenceParser{id: 805, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p804 = charParser{id: 804, chars: []rune{40}}
	p805.items = []parser{&p804}
	var p809 = sequenceParser{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p808 = sequenceParser{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p808.items = []parser{&p835, &p14}
	p809.items = []parser{&p835, &p14, &p808}
	var p811 = sequenceParser{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p810 = sequenceParser{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p810.items = []parser{&p835, &p14}
	p811.items = []parser{&p835, &p14, &p810}
	var p807 = sequenceParser{id: 807, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p806 = charParser{id: 806, chars: []rune{41}}
	p807.items = []parser{&p806}
	p812.items = []parser{&p805, &p809, &p835, &p803, &p811, &p835, &p807}
	p803.options = []parser{&p185, &p431, &p488, &p550, &p591, &p739, &p782, &p792, &p812, &p793}
	var p820 = sequenceParser{id: 820, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p818 = sequenceParser{id: 818, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p818.items = []parser{&p817, &p835, &p803}
	var p819 = sequenceParser{id: 819, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p819.items = []parser{&p835, &p818}
	p820.items = []parser{&p835, &p818, &p819}
	p821.items = []parser{&p803, &p820}
	p836.items = []parser{&p832, &p835, &p817, &p835, &p821, &p835, &p817}
	p837.items = []parser{&p835, &p836, &p835}
	var b837 = sequenceBuilder{id: 837, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b835 = choiceBuilder{id: 835, commit: 2}
	var b833 = choiceBuilder{id: 833, commit: 70}
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
	b833.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b834 = sequenceBuilder{id: 834, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b39.items = []builder{&b14, &b835, &b38}
	var b40 = sequenceBuilder{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b40.items = []builder{&b835, &b39}
	b41.items = []builder{&b835, &b39, &b40}
	b42.items = []builder{&b38, &b41}
	b834.items = []builder{&b42}
	b835.options = []builder{&b833, &b834}
	var b836 = sequenceBuilder{id: 836, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b832 = sequenceBuilder{id: 832, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b829 = sequenceBuilder{id: 829, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b827 = charBuilder{}
	var b828 = charBuilder{}
	b829.items = []builder{&b827, &b828}
	var b826 = sequenceBuilder{id: 826, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b825 = sequenceBuilder{id: 825, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b823 = sequenceBuilder{id: 823, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b822 = charBuilder{}
	b823.items = []builder{&b822}
	var b824 = sequenceBuilder{id: 824, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b824.items = []builder{&b835, &b823}
	b825.items = []builder{&b823, &b824}
	b826.items = []builder{&b825}
	var b831 = sequenceBuilder{id: 831, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b830 = charBuilder{}
	b831.items = []builder{&b830}
	b832.items = []builder{&b829, &b835, &b826, &b835, &b831}
	var b817 = sequenceBuilder{id: 817, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b815 = choiceBuilder{id: 815, commit: 2}
	var b814 = sequenceBuilder{id: 814, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b813 = charBuilder{}
	b814.items = []builder{&b813}
	b815.options = []builder{&b814, &b14}
	var b816 = sequenceBuilder{id: 816, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b816.items = []builder{&b835, &b815}
	b817.items = []builder{&b815, &b816}
	var b821 = sequenceBuilder{id: 821, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b803 = choiceBuilder{id: 803, commit: 66}
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
	b182.items = []builder{&b835, &b14}
	b183.items = []builder{&b14, &b182}
	var b394 = choiceBuilder{id: 394, commit: 66}
	var b266 = choiceBuilder{id: 266, commit: 66}
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
	var b499 = sequenceBuilder{id: 499, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b496 = sequenceBuilder{id: 496, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b489 = charBuilder{}
	var b490 = charBuilder{}
	var b491 = charBuilder{}
	var b492 = charBuilder{}
	var b493 = charBuilder{}
	var b494 = charBuilder{}
	var b495 = charBuilder{}
	b496.items = []builder{&b489, &b490, &b491, &b492, &b493, &b494, &b495}
	var b498 = sequenceBuilder{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b497 = sequenceBuilder{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b497.items = []builder{&b835, &b14}
	b498.items = []builder{&b835, &b14, &b497}
	b499.items = []builder{&b496, &b498, &b835, &b266}
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
	b112.items = []builder{&b835, &b111}
	b113.items = []builder{&b111, &b112}
	var b118 = sequenceBuilder{id: 118, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b114 = choiceBuilder{id: 114, commit: 66}
	var b108 = sequenceBuilder{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b107 = sequenceBuilder{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b104 = charBuilder{}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	b107.items = []builder{&b104, &b105, &b106}
	b108.items = []builder{&b266, &b835, &b107}
	b114.options = []builder{&b394, &b108}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b115 = sequenceBuilder{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b115.items = []builder{&b113, &b835, &b114}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b116.items = []builder{&b835, &b115}
	b117.items = []builder{&b835, &b115, &b116}
	b118.items = []builder{&b114, &b117}
	var b122 = sequenceBuilder{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b121 = charBuilder{}
	b122.items = []builder{&b121}
	b123.items = []builder{&b120, &b835, &b113, &b835, &b118, &b835, &b113, &b835, &b122}
	b124.items = []builder{&b123}
	var b129 = sequenceBuilder{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b126 = sequenceBuilder{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b125 = charBuilder{}
	b126.items = []builder{&b125}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b127 = sequenceBuilder{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b127.items = []builder{&b835, &b14}
	b128.items = []builder{&b835, &b14, &b127}
	b129.items = []builder{&b126, &b128, &b835, &b123}
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
	b134.items = []builder{&b835, &b14}
	b135.items = []builder{&b835, &b14, &b134}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b136 = sequenceBuilder{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b136.items = []builder{&b835, &b14}
	b137.items = []builder{&b835, &b14, &b136}
	var b133 = sequenceBuilder{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b132 = charBuilder{}
	b133.items = []builder{&b132}
	b138.items = []builder{&b131, &b135, &b835, &b394, &b137, &b835, &b133}
	b139.options = []builder{&b103, &b86, &b138}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b142 = sequenceBuilder{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b142.items = []builder{&b835, &b14}
	b143.items = []builder{&b835, &b14, &b142}
	var b141 = sequenceBuilder{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b140 = charBuilder{}
	b141.items = []builder{&b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b835, &b14}
	b145.items = []builder{&b835, &b14, &b144}
	b146.items = []builder{&b139, &b143, &b835, &b141, &b145, &b835, &b394}
	b147.options = []builder{&b146, &b108}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b149 = sequenceBuilder{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b148 = choiceBuilder{id: 148, commit: 2}
	b148.options = []builder{&b146, &b108}
	b149.items = []builder{&b113, &b835, &b148}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b150.items = []builder{&b835, &b149}
	b151.items = []builder{&b835, &b149, &b150}
	b152.items = []builder{&b147, &b151}
	var b156 = sequenceBuilder{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b155 = charBuilder{}
	b156.items = []builder{&b155}
	b157.items = []builder{&b154, &b835, &b113, &b835, &b152, &b835, &b113, &b835, &b156}
	b158.items = []builder{&b157}
	var b163 = sequenceBuilder{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b160 = sequenceBuilder{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b159 = charBuilder{}
	b160.items = []builder{&b159}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b161 = sequenceBuilder{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b161.items = []builder{&b835, &b14}
	b162.items = []builder{&b835, &b14, &b161}
	b163.items = []builder{&b160, &b162, &b835, &b157}
	var b206 = sequenceBuilder{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b203 = sequenceBuilder{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b201 = charBuilder{}
	var b202 = charBuilder{}
	b203.items = []builder{&b201, &b202}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b204 = sequenceBuilder{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b204.items = []builder{&b835, &b14}
	b205.items = []builder{&b835, &b14, &b204}
	var b200 = sequenceBuilder{id: 200, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b192 = sequenceBuilder{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b191 = charBuilder{}
	b192.items = []builder{&b191}
	var b194 = choiceBuilder{id: 194, commit: 2}
	var b167 = sequenceBuilder{id: 167, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b164.items = []builder{&b113, &b835, &b103}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b165.items = []builder{&b835, &b164}
	b166.items = []builder{&b835, &b164, &b165}
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
	b172.items = []builder{&b835, &b14}
	b173.items = []builder{&b835, &b14, &b172}
	b174.items = []builder{&b171, &b173, &b835, &b103}
	b193.items = []builder{&b167, &b835, &b113, &b835, &b174}
	b194.options = []builder{&b167, &b193, &b174}
	var b196 = sequenceBuilder{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b195 = charBuilder{}
	b196.items = []builder{&b195}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b198 = sequenceBuilder{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b198.items = []builder{&b835, &b14}
	b199.items = []builder{&b835, &b14, &b198}
	var b197 = choiceBuilder{id: 197, commit: 2}
	var b793 = choiceBuilder{id: 793, commit: 66}
	var b509 = sequenceBuilder{id: 509, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b504 = sequenceBuilder{id: 504, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b500 = charBuilder{}
	var b501 = charBuilder{}
	var b502 = charBuilder{}
	var b503 = charBuilder{}
	b504.items = []builder{&b500, &b501, &b502, &b503}
	var b506 = sequenceBuilder{id: 506, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b505 = sequenceBuilder{id: 505, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b505.items = []builder{&b835, &b14}
	b506.items = []builder{&b835, &b14, &b505}
	var b508 = sequenceBuilder{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b507 = sequenceBuilder{id: 507, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b507.items = []builder{&b835, &b14}
	b508.items = []builder{&b835, &b14, &b507}
	b509.items = []builder{&b504, &b506, &b835, &b266, &b508, &b835, &b266}
	var b556 = sequenceBuilder{id: 556, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b553 = sequenceBuilder{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b551 = charBuilder{}
	var b552 = charBuilder{}
	b553.items = []builder{&b551, &b552}
	var b555 = sequenceBuilder{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b554 = sequenceBuilder{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b554.items = []builder{&b835, &b14}
	b555.items = []builder{&b835, &b14, &b554}
	var b256 = sequenceBuilder{id: 256, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b253 = sequenceBuilder{id: 253, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b252 = charBuilder{}
	b253.items = []builder{&b252}
	var b255 = sequenceBuilder{id: 255, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b254 = charBuilder{}
	b255.items = []builder{&b254}
	b256.items = []builder{&b266, &b835, &b253, &b835, &b113, &b835, &b118, &b835, &b113, &b835, &b255}
	b556.items = []builder{&b553, &b555, &b835, &b256}
	var b565 = sequenceBuilder{id: 565, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b562 = sequenceBuilder{id: 562, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b557 = charBuilder{}
	var b558 = charBuilder{}
	var b559 = charBuilder{}
	var b560 = charBuilder{}
	var b561 = charBuilder{}
	b562.items = []builder{&b557, &b558, &b559, &b560, &b561}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b563 = sequenceBuilder{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b563.items = []builder{&b835, &b14}
	b564.items = []builder{&b835, &b14, &b563}
	b565.items = []builder{&b562, &b564, &b835, &b256}
	var b630 = choiceBuilder{id: 630, commit: 64, name: "assignment"}
	var b610 = sequenceBuilder{id: 610, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b607 = sequenceBuilder{id: 607, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b604 = charBuilder{}
	var b605 = charBuilder{}
	var b606 = charBuilder{}
	b607.items = []builder{&b604, &b605, &b606}
	var b609 = sequenceBuilder{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b608 = sequenceBuilder{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b608.items = []builder{&b835, &b14}
	b609.items = []builder{&b835, &b14, &b608}
	var b599 = sequenceBuilder{id: 599, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b594.items = []builder{&b835, &b14}
	b595.items = []builder{&b14, &b594}
	var b593 = sequenceBuilder{id: 593, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b592 = charBuilder{}
	b593.items = []builder{&b592}
	b596.items = []builder{&b595, &b835, &b593}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b597.items = []builder{&b835, &b14}
	b598.items = []builder{&b835, &b14, &b597}
	b599.items = []builder{&b266, &b835, &b596, &b598, &b835, &b394}
	b610.items = []builder{&b607, &b609, &b835, &b599}
	var b617 = sequenceBuilder{id: 617, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b613 = sequenceBuilder{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b613.items = []builder{&b835, &b14}
	b614.items = []builder{&b835, &b14, &b613}
	var b612 = sequenceBuilder{id: 612, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b611 = charBuilder{}
	b612.items = []builder{&b611}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b615.items = []builder{&b835, &b14}
	b616.items = []builder{&b835, &b14, &b615}
	b617.items = []builder{&b266, &b614, &b835, &b612, &b616, &b835, &b394}
	var b629 = sequenceBuilder{id: 629, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b621 = sequenceBuilder{id: 621, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b618 = charBuilder{}
	var b619 = charBuilder{}
	var b620 = charBuilder{}
	b621.items = []builder{&b618, &b619, &b620}
	var b628 = sequenceBuilder{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b627 = sequenceBuilder{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b627.items = []builder{&b835, &b14}
	b628.items = []builder{&b835, &b14, &b627}
	var b623 = sequenceBuilder{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b622 = charBuilder{}
	b623.items = []builder{&b622}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b603 = sequenceBuilder{id: 603, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b602 = sequenceBuilder{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b600 = sequenceBuilder{id: 600, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b600.items = []builder{&b113, &b835, &b599}
	var b601 = sequenceBuilder{id: 601, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b601.items = []builder{&b835, &b600}
	b602.items = []builder{&b835, &b600, &b601}
	b603.items = []builder{&b599, &b602}
	b624.items = []builder{&b113, &b835, &b603}
	var b626 = sequenceBuilder{id: 626, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b625 = charBuilder{}
	b626.items = []builder{&b625}
	b629.items = []builder{&b621, &b628, &b835, &b623, &b835, &b624, &b835, &b113, &b835, &b626}
	b630.options = []builder{&b610, &b617, &b629}
	var b802 = sequenceBuilder{id: 802, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b795 = sequenceBuilder{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b794 = charBuilder{}
	b795.items = []builder{&b794}
	var b799 = sequenceBuilder{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b798 = sequenceBuilder{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b798.items = []builder{&b835, &b14}
	b799.items = []builder{&b835, &b14, &b798}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b800.items = []builder{&b835, &b14}
	b801.items = []builder{&b835, &b14, &b800}
	var b797 = sequenceBuilder{id: 797, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b796 = charBuilder{}
	b797.items = []builder{&b796}
	b802.items = []builder{&b795, &b799, &b835, &b793, &b801, &b835, &b797}
	b793.options = []builder{&b509, &b556, &b565, &b630, &b802, &b394}
	var b190 = sequenceBuilder{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b187 = sequenceBuilder{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b186 = charBuilder{}
	b187.items = []builder{&b186}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	b190.items = []builder{&b187, &b835, &b817, &b835, &b821, &b835, &b817, &b835, &b189}
	b197.options = []builder{&b793, &b190}
	b200.items = []builder{&b192, &b835, &b113, &b835, &b194, &b835, &b113, &b835, &b196, &b199, &b835, &b197}
	b206.items = []builder{&b203, &b205, &b835, &b200}
	var b216 = sequenceBuilder{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b209 = sequenceBuilder{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b207 = charBuilder{}
	var b208 = charBuilder{}
	b209.items = []builder{&b207, &b208}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b212 = sequenceBuilder{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b212.items = []builder{&b835, &b14}
	b213.items = []builder{&b835, &b14, &b212}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b210 = charBuilder{}
	b211.items = []builder{&b210}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b214.items = []builder{&b835, &b14}
	b215.items = []builder{&b835, &b14, &b214}
	b216.items = []builder{&b209, &b213, &b835, &b211, &b215, &b835, &b200}
	var b244 = choiceBuilder{id: 244, commit: 64, name: "expression-indexer"}
	var b234 = sequenceBuilder{id: 234, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b227 = sequenceBuilder{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b226 = charBuilder{}
	b227.items = []builder{&b226}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b230 = sequenceBuilder{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b230.items = []builder{&b835, &b14}
	b231.items = []builder{&b835, &b14, &b230}
	var b233 = sequenceBuilder{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b232 = sequenceBuilder{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b232.items = []builder{&b835, &b14}
	b233.items = []builder{&b835, &b14, &b232}
	var b229 = sequenceBuilder{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b228 = charBuilder{}
	b229.items = []builder{&b228}
	b234.items = []builder{&b266, &b835, &b227, &b231, &b835, &b394, &b233, &b835, &b229}
	var b243 = sequenceBuilder{id: 243, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b236 = sequenceBuilder{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b235 = charBuilder{}
	b236.items = []builder{&b235}
	var b240 = sequenceBuilder{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b239 = sequenceBuilder{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b239.items = []builder{&b835, &b14}
	b240.items = []builder{&b835, &b14, &b239}
	var b225 = sequenceBuilder{id: 225, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b217.items = []builder{&b394}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b221 = sequenceBuilder{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b221.items = []builder{&b835, &b14}
	b222.items = []builder{&b835, &b14, &b221}
	var b220 = sequenceBuilder{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b219 = charBuilder{}
	b220.items = []builder{&b219}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b835, &b14}
	b224.items = []builder{&b835, &b14, &b223}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b394}
	b225.items = []builder{&b217, &b222, &b835, &b220, &b224, &b835, &b218}
	var b242 = sequenceBuilder{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b241 = sequenceBuilder{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b241.items = []builder{&b835, &b14}
	b242.items = []builder{&b835, &b14, &b241}
	var b238 = sequenceBuilder{id: 238, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b237 = charBuilder{}
	b238.items = []builder{&b237}
	b243.items = []builder{&b266, &b835, &b236, &b240, &b835, &b225, &b242, &b835, &b238}
	b244.options = []builder{&b234, &b243}
	var b251 = sequenceBuilder{id: 251, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b247.items = []builder{&b835, &b14}
	b248.items = []builder{&b835, &b14, &b247}
	var b246 = sequenceBuilder{id: 246, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b245 = charBuilder{}
	b246.items = []builder{&b245}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b249.items = []builder{&b835, &b14}
	b250.items = []builder{&b835, &b14, &b249}
	b251.items = []builder{&b266, &b248, &b835, &b246, &b250, &b835, &b103}
	var b265 = sequenceBuilder{id: 265, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b258 = sequenceBuilder{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b257 = charBuilder{}
	b258.items = []builder{&b257}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b261 = sequenceBuilder{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b261.items = []builder{&b835, &b14}
	b262.items = []builder{&b835, &b14, &b261}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b263.items = []builder{&b835, &b14}
	b264.items = []builder{&b835, &b14, &b263}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	b265.items = []builder{&b258, &b262, &b835, &b394, &b264, &b835, &b260}
	b266.options = []builder{&b60, &b73, &b86, &b98, &b499, &b103, &b124, &b129, &b158, &b163, &b206, &b216, &b244, &b251, &b256, &b265}
	var b326 = sequenceBuilder{id: 326, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b325 = choiceBuilder{id: 325, commit: 66}
	var b285 = sequenceBuilder{id: 285, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b284 = charBuilder{}
	b285.items = []builder{&b284}
	var b287 = sequenceBuilder{id: 287, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b286 = charBuilder{}
	b287.items = []builder{&b286}
	var b268 = sequenceBuilder{id: 268, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b267 = charBuilder{}
	b268.items = []builder{&b267}
	var b299 = sequenceBuilder{id: 299, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b298 = charBuilder{}
	b299.items = []builder{&b298}
	b325.options = []builder{&b285, &b287, &b268, &b299}
	b326.items = []builder{&b325, &b835, &b266}
	var b380 = choiceBuilder{id: 380, commit: 66}
	var b344 = sequenceBuilder{id: 344, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b332 = choiceBuilder{id: 332, commit: 66}
	b332.options = []builder{&b266, &b326}
	var b342 = sequenceBuilder{id: 342, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b339 = sequenceBuilder{id: 339, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b338 = sequenceBuilder{id: 338, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b338.items = []builder{&b835, &b14}
	b339.items = []builder{&b14, &b338}
	var b327 = choiceBuilder{id: 327, commit: 66}
	var b270 = sequenceBuilder{id: 270, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b269 = charBuilder{}
	b270.items = []builder{&b269}
	var b277 = sequenceBuilder{id: 277, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b275 = charBuilder{}
	var b276 = charBuilder{}
	b277.items = []builder{&b275, &b276}
	var b280 = sequenceBuilder{id: 280, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b278 = charBuilder{}
	var b279 = charBuilder{}
	b280.items = []builder{&b278, &b279}
	var b283 = sequenceBuilder{id: 283, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b281 = charBuilder{}
	var b282 = charBuilder{}
	b283.items = []builder{&b281, &b282}
	var b289 = sequenceBuilder{id: 289, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b288 = charBuilder{}
	b289.items = []builder{&b288}
	var b291 = sequenceBuilder{id: 291, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b290 = charBuilder{}
	b291.items = []builder{&b290}
	var b293 = sequenceBuilder{id: 293, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b292 = charBuilder{}
	b293.items = []builder{&b292}
	b327.options = []builder{&b270, &b277, &b280, &b283, &b289, &b291, &b293}
	var b341 = sequenceBuilder{id: 341, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b340 = sequenceBuilder{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b340.items = []builder{&b835, &b14}
	b341.items = []builder{&b835, &b14, &b340}
	b342.items = []builder{&b339, &b835, &b327, &b341, &b835, &b332}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b343.items = []builder{&b835, &b342}
	b344.items = []builder{&b332, &b835, &b342, &b343}
	var b351 = sequenceBuilder{id: 351, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b333 = choiceBuilder{id: 333, commit: 66}
	b333.options = []builder{&b332, &b344}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b345 = sequenceBuilder{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b345.items = []builder{&b835, &b14}
	b346.items = []builder{&b14, &b345}
	var b328 = choiceBuilder{id: 328, commit: 66}
	var b272 = sequenceBuilder{id: 272, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b271 = charBuilder{}
	b272.items = []builder{&b271}
	var b274 = sequenceBuilder{id: 274, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b273 = charBuilder{}
	b274.items = []builder{&b273}
	var b295 = sequenceBuilder{id: 295, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b294 = charBuilder{}
	b295.items = []builder{&b294}
	var b297 = sequenceBuilder{id: 297, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b296 = charBuilder{}
	b297.items = []builder{&b296}
	b328.options = []builder{&b272, &b274, &b295, &b297}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b347.items = []builder{&b835, &b14}
	b348.items = []builder{&b835, &b14, &b347}
	b349.items = []builder{&b346, &b835, &b328, &b348, &b835, &b333}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b350.items = []builder{&b835, &b349}
	b351.items = []builder{&b333, &b835, &b349, &b350}
	var b358 = sequenceBuilder{id: 358, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b334 = choiceBuilder{id: 334, commit: 66}
	b334.options = []builder{&b333, &b351}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b352.items = []builder{&b835, &b14}
	b353.items = []builder{&b14, &b352}
	var b329 = choiceBuilder{id: 329, commit: 66}
	var b302 = sequenceBuilder{id: 302, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b300 = charBuilder{}
	var b301 = charBuilder{}
	b302.items = []builder{&b300, &b301}
	var b305 = sequenceBuilder{id: 305, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b303 = charBuilder{}
	var b304 = charBuilder{}
	b305.items = []builder{&b303, &b304}
	var b307 = sequenceBuilder{id: 307, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b306 = charBuilder{}
	b307.items = []builder{&b306}
	var b310 = sequenceBuilder{id: 310, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b308 = charBuilder{}
	var b309 = charBuilder{}
	b310.items = []builder{&b308, &b309}
	var b312 = sequenceBuilder{id: 312, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b311 = charBuilder{}
	b312.items = []builder{&b311}
	var b315 = sequenceBuilder{id: 315, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b313 = charBuilder{}
	var b314 = charBuilder{}
	b315.items = []builder{&b313, &b314}
	b329.options = []builder{&b302, &b305, &b307, &b310, &b312, &b315}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b354.items = []builder{&b835, &b14}
	b355.items = []builder{&b835, &b14, &b354}
	b356.items = []builder{&b353, &b835, &b329, &b355, &b835, &b334}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b357.items = []builder{&b835, &b356}
	b358.items = []builder{&b334, &b835, &b356, &b357}
	var b365 = sequenceBuilder{id: 365, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b335 = choiceBuilder{id: 335, commit: 66}
	b335.options = []builder{&b334, &b358}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b835, &b14}
	b360.items = []builder{&b14, &b359}
	var b330 = sequenceBuilder{id: 330, commit: 66, ranges: [][]int{{1, 1}}}
	var b318 = sequenceBuilder{id: 318, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b316 = charBuilder{}
	var b317 = charBuilder{}
	b318.items = []builder{&b316, &b317}
	b330.items = []builder{&b318}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b361.items = []builder{&b835, &b14}
	b362.items = []builder{&b835, &b14, &b361}
	b363.items = []builder{&b360, &b835, &b330, &b362, &b835, &b335}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b364.items = []builder{&b835, &b363}
	b365.items = []builder{&b335, &b835, &b363, &b364}
	var b372 = sequenceBuilder{id: 372, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b336 = choiceBuilder{id: 336, commit: 66}
	b336.options = []builder{&b335, &b365}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b366.items = []builder{&b835, &b14}
	b367.items = []builder{&b14, &b366}
	var b331 = sequenceBuilder{id: 331, commit: 66, ranges: [][]int{{1, 1}}}
	var b321 = sequenceBuilder{id: 321, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b319 = charBuilder{}
	var b320 = charBuilder{}
	b321.items = []builder{&b319, &b320}
	b331.items = []builder{&b321}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b368.items = []builder{&b835, &b14}
	b369.items = []builder{&b835, &b14, &b368}
	b370.items = []builder{&b367, &b835, &b331, &b369, &b835, &b336}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b371.items = []builder{&b835, &b370}
	b372.items = []builder{&b336, &b835, &b370, &b371}
	var b379 = sequenceBuilder{id: 379, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b337 = choiceBuilder{id: 337, commit: 66}
	b337.options = []builder{&b336, &b372}
	var b377 = sequenceBuilder{id: 377, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b374 = sequenceBuilder{id: 374, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b373.items = []builder{&b835, &b14}
	b374.items = []builder{&b14, &b373}
	var b324 = sequenceBuilder{id: 324, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b322 = charBuilder{}
	var b323 = charBuilder{}
	b324.items = []builder{&b322, &b323}
	var b376 = sequenceBuilder{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b375 = sequenceBuilder{id: 375, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b375.items = []builder{&b835, &b14}
	b376.items = []builder{&b835, &b14, &b375}
	b377.items = []builder{&b374, &b835, &b324, &b376, &b835, &b337}
	var b378 = sequenceBuilder{id: 378, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b378.items = []builder{&b835, &b377}
	b379.items = []builder{&b337, &b835, &b377, &b378}
	b380.options = []builder{&b344, &b351, &b358, &b365, &b372, &b379}
	var b393 = sequenceBuilder{id: 393, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b385.items = []builder{&b835, &b14}
	b386.items = []builder{&b835, &b14, &b385}
	var b382 = sequenceBuilder{id: 382, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b381 = charBuilder{}
	b382.items = []builder{&b381}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b387.items = []builder{&b835, &b14}
	b388.items = []builder{&b835, &b14, &b387}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b389.items = []builder{&b835, &b14}
	b390.items = []builder{&b835, &b14, &b389}
	var b384 = sequenceBuilder{id: 384, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b383 = charBuilder{}
	b384.items = []builder{&b383}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b391.items = []builder{&b835, &b14}
	b392.items = []builder{&b835, &b14, &b391}
	b393.items = []builder{&b394, &b386, &b835, &b382, &b388, &b835, &b394, &b390, &b835, &b384, &b392, &b835, &b394}
	b394.options = []builder{&b266, &b326, &b380, &b393}
	b184.items = []builder{&b183, &b835, &b394}
	b185.items = []builder{&b181, &b835, &b184}
	var b431 = sequenceBuilder{id: 431, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b397 = sequenceBuilder{id: 397, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b395 = charBuilder{}
	var b396 = charBuilder{}
	b397.items = []builder{&b395, &b396}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b425.items = []builder{&b835, &b14}
	b426.items = []builder{&b835, &b14, &b425}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b427.items = []builder{&b835, &b14}
	b428.items = []builder{&b835, &b14, &b427}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b407 = sequenceBuilder{id: 407, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b406 = sequenceBuilder{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b406.items = []builder{&b835, &b14}
	b407.items = []builder{&b14, &b406}
	var b402 = sequenceBuilder{id: 402, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b398 = charBuilder{}
	var b399 = charBuilder{}
	var b400 = charBuilder{}
	var b401 = charBuilder{}
	b402.items = []builder{&b398, &b399, &b400, &b401}
	var b409 = sequenceBuilder{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b408.items = []builder{&b835, &b14}
	b409.items = []builder{&b835, &b14, &b408}
	var b405 = sequenceBuilder{id: 405, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b403 = charBuilder{}
	var b404 = charBuilder{}
	b405.items = []builder{&b403, &b404}
	var b411 = sequenceBuilder{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b410 = sequenceBuilder{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b410.items = []builder{&b835, &b14}
	b411.items = []builder{&b835, &b14, &b410}
	var b413 = sequenceBuilder{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b412 = sequenceBuilder{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b412.items = []builder{&b835, &b14}
	b413.items = []builder{&b835, &b14, &b412}
	b414.items = []builder{&b407, &b835, &b402, &b409, &b835, &b405, &b411, &b835, &b394, &b413, &b835, &b190}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b429.items = []builder{&b835, &b414}
	b430.items = []builder{&b835, &b414, &b429}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b421 = sequenceBuilder{id: 421, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b420 = sequenceBuilder{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b420.items = []builder{&b835, &b14}
	b421.items = []builder{&b14, &b420}
	var b419 = sequenceBuilder{id: 419, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b415 = charBuilder{}
	var b416 = charBuilder{}
	var b417 = charBuilder{}
	var b418 = charBuilder{}
	b419.items = []builder{&b415, &b416, &b417, &b418}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b422.items = []builder{&b835, &b14}
	b423.items = []builder{&b835, &b14, &b422}
	b424.items = []builder{&b421, &b835, &b419, &b423, &b835, &b190}
	b431.items = []builder{&b397, &b426, &b835, &b394, &b428, &b835, &b190, &b430, &b835, &b424}
	var b488 = sequenceBuilder{id: 488, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b473 = sequenceBuilder{id: 473, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b467 = charBuilder{}
	var b468 = charBuilder{}
	var b469 = charBuilder{}
	var b470 = charBuilder{}
	var b471 = charBuilder{}
	var b472 = charBuilder{}
	b473.items = []builder{&b467, &b468, &b469, &b470, &b471, &b472}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b484 = sequenceBuilder{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b484.items = []builder{&b835, &b14}
	b485.items = []builder{&b835, &b14, &b484}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b486.items = []builder{&b835, &b14}
	b487.items = []builder{&b835, &b14, &b486}
	var b475 = sequenceBuilder{id: 475, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b474 = charBuilder{}
	b475.items = []builder{&b474}
	var b481 = sequenceBuilder{id: 481, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b476 = choiceBuilder{id: 476, commit: 2}
	var b466 = sequenceBuilder{id: 466, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b461 = sequenceBuilder{id: 461, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b454 = sequenceBuilder{id: 454, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b450 = charBuilder{}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	b454.items = []builder{&b450, &b451, &b452, &b453}
	var b458 = sequenceBuilder{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b457 = sequenceBuilder{id: 457, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b457.items = []builder{&b835, &b14}
	b458.items = []builder{&b835, &b14, &b457}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b459.items = []builder{&b835, &b14}
	b460.items = []builder{&b835, &b14, &b459}
	var b456 = sequenceBuilder{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b455 = charBuilder{}
	b456.items = []builder{&b455}
	b461.items = []builder{&b454, &b458, &b835, &b394, &b460, &b835, &b456}
	var b465 = sequenceBuilder{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b463 = sequenceBuilder{id: 463, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b462 = charBuilder{}
	b463.items = []builder{&b462}
	var b464 = sequenceBuilder{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b464.items = []builder{&b835, &b463}
	b465.items = []builder{&b835, &b463, &b464}
	b466.items = []builder{&b461, &b465, &b835, &b803}
	var b449 = sequenceBuilder{id: 449, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b444 = sequenceBuilder{id: 444, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b439 = sequenceBuilder{id: 439, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b432 = charBuilder{}
	var b433 = charBuilder{}
	var b434 = charBuilder{}
	var b435 = charBuilder{}
	var b436 = charBuilder{}
	var b437 = charBuilder{}
	var b438 = charBuilder{}
	b439.items = []builder{&b432, &b433, &b434, &b435, &b436, &b437, &b438}
	var b443 = sequenceBuilder{id: 443, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b442 = sequenceBuilder{id: 442, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b442.items = []builder{&b835, &b14}
	b443.items = []builder{&b835, &b14, &b442}
	var b441 = sequenceBuilder{id: 441, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b440 = charBuilder{}
	b441.items = []builder{&b440}
	b444.items = []builder{&b439, &b443, &b835, &b441}
	var b448 = sequenceBuilder{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b446 = sequenceBuilder{id: 446, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b445 = charBuilder{}
	b446.items = []builder{&b445}
	var b447 = sequenceBuilder{id: 447, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b447.items = []builder{&b835, &b446}
	b448.items = []builder{&b835, &b446, &b447}
	b449.items = []builder{&b444, &b448, &b835, &b803}
	b476.options = []builder{&b466, &b449}
	var b480 = sequenceBuilder{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b478 = sequenceBuilder{id: 478, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b477 = choiceBuilder{id: 477, commit: 2}
	b477.options = []builder{&b466, &b449, &b803}
	b478.items = []builder{&b817, &b835, &b477}
	var b479 = sequenceBuilder{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b479.items = []builder{&b835, &b478}
	b480.items = []builder{&b835, &b478, &b479}
	b481.items = []builder{&b476, &b480}
	var b483 = sequenceBuilder{id: 483, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b482 = charBuilder{}
	b483.items = []builder{&b482}
	b488.items = []builder{&b473, &b485, &b835, &b394, &b487, &b835, &b475, &b835, &b817, &b835, &b481, &b835, &b817, &b835, &b483}
	var b550 = sequenceBuilder{id: 550, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b537 = sequenceBuilder{id: 537, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b531 = charBuilder{}
	var b532 = charBuilder{}
	var b533 = charBuilder{}
	var b534 = charBuilder{}
	var b535 = charBuilder{}
	var b536 = charBuilder{}
	b537.items = []builder{&b531, &b532, &b533, &b534, &b535, &b536}
	var b549 = sequenceBuilder{id: 549, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b548 = sequenceBuilder{id: 548, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b548.items = []builder{&b835, &b14}
	b549.items = []builder{&b835, &b14, &b548}
	var b539 = sequenceBuilder{id: 539, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b538 = charBuilder{}
	b539.items = []builder{&b538}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b540 = choiceBuilder{id: 540, commit: 2}
	var b530 = sequenceBuilder{id: 530, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b525 = sequenceBuilder{id: 525, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b518 = sequenceBuilder{id: 518, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b514 = charBuilder{}
	var b515 = charBuilder{}
	var b516 = charBuilder{}
	var b517 = charBuilder{}
	b518.items = []builder{&b514, &b515, &b516, &b517}
	var b522 = sequenceBuilder{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b521 = sequenceBuilder{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b521.items = []builder{&b835, &b14}
	b522.items = []builder{&b835, &b14, &b521}
	var b513 = choiceBuilder{id: 513, commit: 66}
	var b512 = sequenceBuilder{id: 512, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b511 = sequenceBuilder{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b510.items = []builder{&b835, &b14}
	b511.items = []builder{&b835, &b14, &b510}
	b512.items = []builder{&b103, &b511, &b835, &b499}
	b513.options = []builder{&b499, &b512, &b509}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b523.items = []builder{&b835, &b14}
	b524.items = []builder{&b835, &b14, &b523}
	var b520 = sequenceBuilder{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b519 = charBuilder{}
	b520.items = []builder{&b519}
	b525.items = []builder{&b518, &b522, &b835, &b513, &b524, &b835, &b520}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b527 = sequenceBuilder{id: 527, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b526 = charBuilder{}
	b527.items = []builder{&b526}
	var b528 = sequenceBuilder{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b528.items = []builder{&b835, &b527}
	b529.items = []builder{&b835, &b527, &b528}
	b530.items = []builder{&b525, &b529, &b835, &b803}
	b540.options = []builder{&b530, &b449}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b542 = sequenceBuilder{id: 542, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b541 = choiceBuilder{id: 541, commit: 2}
	b541.options = []builder{&b530, &b449, &b803}
	b542.items = []builder{&b817, &b835, &b541}
	var b543 = sequenceBuilder{id: 543, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b543.items = []builder{&b835, &b542}
	b544.items = []builder{&b835, &b542, &b543}
	b545.items = []builder{&b540, &b544}
	var b547 = sequenceBuilder{id: 547, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b546 = charBuilder{}
	b547.items = []builder{&b546}
	b550.items = []builder{&b537, &b549, &b835, &b539, &b835, &b817, &b835, &b545, &b835, &b817, &b835, &b547}
	var b591 = sequenceBuilder{id: 591, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b580 = sequenceBuilder{id: 580, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b577 = charBuilder{}
	var b578 = charBuilder{}
	var b579 = charBuilder{}
	b580.items = []builder{&b577, &b578, &b579}
	var b590 = choiceBuilder{id: 590, commit: 2}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b583 = sequenceBuilder{id: 583, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b582 = sequenceBuilder{id: 582, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b581 = sequenceBuilder{id: 581, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b581.items = []builder{&b835, &b14}
	b582.items = []builder{&b14, &b581}
	var b576 = choiceBuilder{id: 576, commit: 66}
	var b575 = choiceBuilder{id: 575, commit: 64, name: "range-over-expression"}
	var b574 = sequenceBuilder{id: 574, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b571 = sequenceBuilder{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b570 = sequenceBuilder{id: 570, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b570.items = []builder{&b835, &b14}
	b571.items = []builder{&b835, &b14, &b570}
	var b568 = sequenceBuilder{id: 568, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b566 = charBuilder{}
	var b567 = charBuilder{}
	b568.items = []builder{&b566, &b567}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b572.items = []builder{&b835, &b14}
	b573.items = []builder{&b835, &b14, &b572}
	var b569 = choiceBuilder{id: 569, commit: 2}
	b569.options = []builder{&b394, &b225}
	b574.items = []builder{&b103, &b571, &b835, &b568, &b573, &b835, &b569}
	b575.options = []builder{&b574, &b225}
	b576.options = []builder{&b394, &b575}
	b583.items = []builder{&b582, &b835, &b576}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b584.items = []builder{&b835, &b14}
	b585.items = []builder{&b835, &b14, &b584}
	b586.items = []builder{&b583, &b585, &b835, &b190}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b587.items = []builder{&b835, &b14}
	b588.items = []builder{&b14, &b587}
	b589.items = []builder{&b588, &b835, &b190}
	b590.options = []builder{&b586, &b589}
	b591.items = []builder{&b580, &b835, &b590}
	var b739 = choiceBuilder{id: 739, commit: 66}
	var b652 = sequenceBuilder{id: 652, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b648 = sequenceBuilder{id: 648, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b645 = charBuilder{}
	var b646 = charBuilder{}
	var b647 = charBuilder{}
	b648.items = []builder{&b645, &b646, &b647}
	var b651 = sequenceBuilder{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b650 = sequenceBuilder{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b650.items = []builder{&b835, &b14}
	b651.items = []builder{&b835, &b14, &b650}
	var b649 = choiceBuilder{id: 649, commit: 2}
	var b639 = sequenceBuilder{id: 639, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b638 = sequenceBuilder{id: 638, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b633 = sequenceBuilder{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b633.items = []builder{&b835, &b14}
	b634.items = []builder{&b14, &b633}
	var b632 = sequenceBuilder{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b631 = charBuilder{}
	b632.items = []builder{&b631}
	b635.items = []builder{&b634, &b835, &b632}
	var b637 = sequenceBuilder{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b636.items = []builder{&b835, &b14}
	b637.items = []builder{&b835, &b14, &b636}
	b638.items = []builder{&b103, &b835, &b635, &b637, &b835, &b394}
	b639.items = []builder{&b638}
	var b644 = sequenceBuilder{id: 644, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b641 = sequenceBuilder{id: 641, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b640 = charBuilder{}
	b641.items = []builder{&b640}
	var b643 = sequenceBuilder{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b642.items = []builder{&b835, &b14}
	b643.items = []builder{&b835, &b14, &b642}
	b644.items = []builder{&b641, &b643, &b835, &b638}
	b649.options = []builder{&b639, &b644}
	b652.items = []builder{&b648, &b651, &b835, &b649}
	var b673 = sequenceBuilder{id: 673, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b666 = sequenceBuilder{id: 666, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b663 = charBuilder{}
	var b664 = charBuilder{}
	var b665 = charBuilder{}
	b666.items = []builder{&b663, &b664, &b665}
	var b672 = sequenceBuilder{id: 672, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b671 = sequenceBuilder{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b671.items = []builder{&b835, &b14}
	b672.items = []builder{&b835, &b14, &b671}
	var b668 = sequenceBuilder{id: 668, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b667 = charBuilder{}
	b668.items = []builder{&b667}
	var b658 = sequenceBuilder{id: 658, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b653 = choiceBuilder{id: 653, commit: 2}
	b653.options = []builder{&b639, &b644}
	var b657 = sequenceBuilder{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b655 = sequenceBuilder{id: 655, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b654 = choiceBuilder{id: 654, commit: 2}
	b654.options = []builder{&b639, &b644}
	b655.items = []builder{&b113, &b835, &b654}
	var b656 = sequenceBuilder{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b656.items = []builder{&b835, &b655}
	b657.items = []builder{&b835, &b655, &b656}
	b658.items = []builder{&b653, &b657}
	var b670 = sequenceBuilder{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b669 = charBuilder{}
	b670.items = []builder{&b669}
	b673.items = []builder{&b666, &b672, &b835, &b668, &b835, &b113, &b835, &b658, &b835, &b113, &b835, &b670}
	var b688 = sequenceBuilder{id: 688, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b677 = sequenceBuilder{id: 677, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b674 = charBuilder{}
	var b675 = charBuilder{}
	var b676 = charBuilder{}
	b677.items = []builder{&b674, &b675, &b676}
	var b685 = sequenceBuilder{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b684 = sequenceBuilder{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b684.items = []builder{&b835, &b14}
	b685.items = []builder{&b835, &b14, &b684}
	var b679 = sequenceBuilder{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b678 = charBuilder{}
	b679.items = []builder{&b678}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b686.items = []builder{&b835, &b14}
	b687.items = []builder{&b835, &b14, &b686}
	var b681 = sequenceBuilder{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	b681.items = []builder{&b680}
	var b662 = sequenceBuilder{id: 662, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b661 = sequenceBuilder{id: 661, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b659 = sequenceBuilder{id: 659, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b659.items = []builder{&b113, &b835, &b639}
	var b660 = sequenceBuilder{id: 660, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b660.items = []builder{&b835, &b659}
	b661.items = []builder{&b835, &b659, &b660}
	b662.items = []builder{&b639, &b661}
	var b683 = sequenceBuilder{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b682 = charBuilder{}
	b683.items = []builder{&b682}
	b688.items = []builder{&b677, &b685, &b835, &b679, &b687, &b835, &b681, &b835, &b113, &b835, &b662, &b835, &b113, &b835, &b683}
	var b704 = sequenceBuilder{id: 704, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b700 = sequenceBuilder{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b698 = charBuilder{}
	var b699 = charBuilder{}
	b700.items = []builder{&b698, &b699}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b702 = sequenceBuilder{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b702.items = []builder{&b835, &b14}
	b703.items = []builder{&b835, &b14, &b702}
	var b701 = choiceBuilder{id: 701, commit: 2}
	var b692 = sequenceBuilder{id: 692, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b691 = sequenceBuilder{id: 691, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b690 = sequenceBuilder{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b689 = sequenceBuilder{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b689.items = []builder{&b835, &b14}
	b690.items = []builder{&b835, &b14, &b689}
	b691.items = []builder{&b103, &b690, &b835, &b200}
	b692.items = []builder{&b691}
	var b697 = sequenceBuilder{id: 697, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b694 = sequenceBuilder{id: 694, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b693 = charBuilder{}
	b694.items = []builder{&b693}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b695 = sequenceBuilder{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b695.items = []builder{&b835, &b14}
	b696.items = []builder{&b835, &b14, &b695}
	b697.items = []builder{&b694, &b696, &b835, &b691}
	b701.options = []builder{&b692, &b697}
	b704.items = []builder{&b700, &b703, &b835, &b701}
	var b724 = sequenceBuilder{id: 724, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b717 = sequenceBuilder{id: 717, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b715 = charBuilder{}
	var b716 = charBuilder{}
	b717.items = []builder{&b715, &b716}
	var b723 = sequenceBuilder{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b722 = sequenceBuilder{id: 722, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b722.items = []builder{&b835, &b14}
	b723.items = []builder{&b835, &b14, &b722}
	var b719 = sequenceBuilder{id: 719, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b718 = charBuilder{}
	b719.items = []builder{&b718}
	var b714 = sequenceBuilder{id: 714, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b709 = choiceBuilder{id: 709, commit: 2}
	b709.options = []builder{&b692, &b697}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b711 = sequenceBuilder{id: 711, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b710 = choiceBuilder{id: 710, commit: 2}
	b710.options = []builder{&b692, &b697}
	b711.items = []builder{&b113, &b835, &b710}
	var b712 = sequenceBuilder{id: 712, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b712.items = []builder{&b835, &b711}
	b713.items = []builder{&b835, &b711, &b712}
	b714.items = []builder{&b709, &b713}
	var b721 = sequenceBuilder{id: 721, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b720 = charBuilder{}
	b721.items = []builder{&b720}
	b724.items = []builder{&b717, &b723, &b835, &b719, &b835, &b113, &b835, &b714, &b835, &b113, &b835, &b721}
	var b738 = sequenceBuilder{id: 738, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b727 = sequenceBuilder{id: 727, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b725 = charBuilder{}
	var b726 = charBuilder{}
	b727.items = []builder{&b725, &b726}
	var b735 = sequenceBuilder{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b734 = sequenceBuilder{id: 734, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b734.items = []builder{&b835, &b14}
	b735.items = []builder{&b835, &b14, &b734}
	var b729 = sequenceBuilder{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b728 = charBuilder{}
	b729.items = []builder{&b728}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b736.items = []builder{&b835, &b14}
	b737.items = []builder{&b835, &b14, &b736}
	var b731 = sequenceBuilder{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b730 = charBuilder{}
	b731.items = []builder{&b730}
	var b708 = sequenceBuilder{id: 708, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b707 = sequenceBuilder{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b705 = sequenceBuilder{id: 705, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b705.items = []builder{&b113, &b835, &b692}
	var b706 = sequenceBuilder{id: 706, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b706.items = []builder{&b835, &b705}
	b707.items = []builder{&b835, &b705, &b706}
	b708.items = []builder{&b692, &b707}
	var b733 = sequenceBuilder{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b732 = charBuilder{}
	b733.items = []builder{&b732}
	b738.items = []builder{&b727, &b735, &b835, &b729, &b737, &b835, &b731, &b835, &b113, &b835, &b708, &b835, &b113, &b835, &b733}
	b739.options = []builder{&b652, &b673, &b688, &b704, &b724, &b738}
	var b782 = choiceBuilder{id: 782, commit: 64, name: "require"}
	var b766 = sequenceBuilder{id: 766, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b763 = sequenceBuilder{id: 763, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b756 = charBuilder{}
	var b757 = charBuilder{}
	var b758 = charBuilder{}
	var b759 = charBuilder{}
	var b760 = charBuilder{}
	var b761 = charBuilder{}
	var b762 = charBuilder{}
	b763.items = []builder{&b756, &b757, &b758, &b759, &b760, &b761, &b762}
	var b765 = sequenceBuilder{id: 765, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b764 = sequenceBuilder{id: 764, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b764.items = []builder{&b835, &b14}
	b765.items = []builder{&b835, &b14, &b764}
	var b751 = choiceBuilder{id: 751, commit: 64, name: "require-fact"}
	var b750 = sequenceBuilder{id: 750, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b742 = choiceBuilder{id: 742, commit: 2}
	var b741 = sequenceBuilder{id: 741, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b740 = charBuilder{}
	b741.items = []builder{&b740}
	b742.options = []builder{&b103, &b741}
	var b747 = sequenceBuilder{id: 747, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b746 = sequenceBuilder{id: 746, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b745 = sequenceBuilder{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b745.items = []builder{&b835, &b14}
	b746.items = []builder{&b14, &b745}
	var b744 = sequenceBuilder{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b743 = charBuilder{}
	b744.items = []builder{&b743}
	b747.items = []builder{&b746, &b835, &b744}
	var b749 = sequenceBuilder{id: 749, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b748 = sequenceBuilder{id: 748, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b748.items = []builder{&b835, &b14}
	b749.items = []builder{&b835, &b14, &b748}
	b750.items = []builder{&b742, &b835, &b747, &b749, &b835, &b86}
	b751.options = []builder{&b86, &b750}
	b766.items = []builder{&b763, &b765, &b835, &b751}
	var b781 = sequenceBuilder{id: 781, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b774 = sequenceBuilder{id: 774, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b767 = charBuilder{}
	var b768 = charBuilder{}
	var b769 = charBuilder{}
	var b770 = charBuilder{}
	var b771 = charBuilder{}
	var b772 = charBuilder{}
	var b773 = charBuilder{}
	b774.items = []builder{&b767, &b768, &b769, &b770, &b771, &b772, &b773}
	var b780 = sequenceBuilder{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b779 = sequenceBuilder{id: 779, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b779.items = []builder{&b835, &b14}
	b780.items = []builder{&b835, &b14, &b779}
	var b776 = sequenceBuilder{id: 776, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b775 = charBuilder{}
	b776.items = []builder{&b775}
	var b755 = sequenceBuilder{id: 755, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b752.items = []builder{&b113, &b835, &b751}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b753.items = []builder{&b835, &b752}
	b754.items = []builder{&b835, &b752, &b753}
	b755.items = []builder{&b751, &b754}
	var b778 = sequenceBuilder{id: 778, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b777 = charBuilder{}
	b778.items = []builder{&b777}
	b781.items = []builder{&b774, &b780, &b835, &b776, &b835, &b113, &b835, &b755, &b835, &b113, &b835, &b778}
	b782.options = []builder{&b766, &b781}
	var b792 = sequenceBuilder{id: 792, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b789 = sequenceBuilder{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b783 = charBuilder{}
	var b784 = charBuilder{}
	var b785 = charBuilder{}
	var b786 = charBuilder{}
	var b787 = charBuilder{}
	var b788 = charBuilder{}
	b789.items = []builder{&b783, &b784, &b785, &b786, &b787, &b788}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b790.items = []builder{&b835, &b14}
	b791.items = []builder{&b835, &b14, &b790}
	b792.items = []builder{&b789, &b791, &b835, &b739}
	var b812 = sequenceBuilder{id: 812, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b805 = sequenceBuilder{id: 805, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b804 = charBuilder{}
	b805.items = []builder{&b804}
	var b809 = sequenceBuilder{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b808 = sequenceBuilder{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b808.items = []builder{&b835, &b14}
	b809.items = []builder{&b835, &b14, &b808}
	var b811 = sequenceBuilder{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b810 = sequenceBuilder{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b810.items = []builder{&b835, &b14}
	b811.items = []builder{&b835, &b14, &b810}
	var b807 = sequenceBuilder{id: 807, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b806 = charBuilder{}
	b807.items = []builder{&b806}
	b812.items = []builder{&b805, &b809, &b835, &b803, &b811, &b835, &b807}
	b803.options = []builder{&b185, &b431, &b488, &b550, &b591, &b739, &b782, &b792, &b812, &b793}
	var b820 = sequenceBuilder{id: 820, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b818 = sequenceBuilder{id: 818, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b818.items = []builder{&b817, &b835, &b803}
	var b819 = sequenceBuilder{id: 819, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b819.items = []builder{&b835, &b818}
	b820.items = []builder{&b835, &b818, &b819}
	b821.items = []builder{&b803, &b820}
	b836.items = []builder{&b832, &b835, &b817, &b835, &b821, &b835, &b817}
	b837.items = []builder{&b835, &b836, &b835}

	return parseInput(r, &p837, &b837)
}
