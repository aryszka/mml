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
	name     string
	id       int
	commit   CommitType
	items    []builder
	ranges   [][]int
	allChars bool
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
	} else {
		if c.results.pending(c.offset, b.id) {
			return nil, false
		}
		c.results.markPending(c.offset, b.id)
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
	name    string
	id      int
	commit  CommitType
	options []builder
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
	} else {
		if c.results.pending(c.offset, b.id) {
			return nil, false
		}
		c.results.markPending(c.offset, b.id)
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
	var p844 = sequenceParser{id: 844, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p842 = choiceParser{id: 842, commit: 2}
	var p840 = choiceParser{id: 840, commit: 70, name: "ws", generalizations: []int{842}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{840, 842}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{840, 842}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{840, 842}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{840, 842}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{840, 842}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{840, 842}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p840.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p841 = sequenceParser{id: 841, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{842}}
	var p44 = sequenceParser{id: 44, commit: 66, name: "comment", ranges: [][]int{{1, 1}, {0, 1}}}
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
	var p43 = sequenceParser{id: 43, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p41 = sequenceParser{id: 41, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p40 = sequenceParser{id: 40, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p39 = charParser{id: 39, chars: []rune{10}}
	p40.items = []parser{&p39}
	p41.items = []parser{&p40, &p842, &p38}
	var p42 = sequenceParser{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p42.items = []parser{&p842, &p41}
	p43.items = []parser{&p842, &p41, &p42}
	p44.items = []parser{&p38, &p43}
	p841.items = []parser{&p44}
	p842.options = []parser{&p840, &p841}
	var p843 = sequenceParser{id: 843, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p839 = sequenceParser{id: 839, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p836 = sequenceParser{id: 836, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p834 = charParser{id: 834, chars: []rune{35}}
	var p835 = charParser{id: 835, chars: []rune{33}}
	p836.items = []parser{&p834, &p835}
	var p833 = sequenceParser{id: 833, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p832 = sequenceParser{id: 832, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p830 = sequenceParser{id: 830, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p829 = charParser{id: 829, not: true, chars: []rune{10}}
	p830.items = []parser{&p829}
	var p831 = sequenceParser{id: 831, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p831.items = []parser{&p842, &p830}
	p832.items = []parser{&p830, &p831}
	p833.items = []parser{&p832}
	var p838 = sequenceParser{id: 838, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p837 = charParser{id: 837, chars: []rune{10}}
	p838.items = []parser{&p837}
	p839.items = []parser{&p836, &p842, &p833, &p842, &p838}
	var p824 = sequenceParser{id: 824, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p822 = choiceParser{id: 822, commit: 2}
	var p821 = sequenceParser{id: 821, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{822}}
	var p820 = charParser{id: 820, chars: []rune{59}}
	p821.items = []parser{&p820}
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{822, 113}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p822.options = []parser{&p821, &p14}
	var p823 = sequenceParser{id: 823, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p823.items = []parser{&p842, &p822}
	p824.items = []parser{&p822, &p823}
	var p828 = sequenceParser{id: 828, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p810 = choiceParser{id: 810, commit: 66, name: "statement", generalizations: []int{471, 538}}
	var p800 = choiceParser{id: 800, commit: 66, name: "simple-statement", generalizations: []int{211, 810, 471, 538}}
	var p388 = choiceParser{id: 388, commit: 66, name: "expression", generalizations: []int{116, 211, 800, 810, 471, 538, 588, 581}}
	var p280 = choiceParser{id: 280, commit: 66, name: "primary-expression", generalizations: []int{116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p62 = choiceParser{id: 62, commit: 64, name: "int", generalizations: []int{280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p53 = sequenceParser{id: 53, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p52 = sequenceParser{id: 52, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p51 = charParser{id: 51, ranges: [][]rune{{49, 57}}}
	p52.items = []parser{&p51}
	var p46 = sequenceParser{id: 46, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 57}}}
	p46.items = []parser{&p45}
	p53.items = []parser{&p52, &p46}
	var p56 = sequenceParser{id: 56, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p55 = sequenceParser{id: 55, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p54 = charParser{id: 54, chars: []rune{48}}
	p55.items = []parser{&p54}
	var p48 = sequenceParser{id: 48, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p47 = charParser{id: 47, ranges: [][]rune{{48, 55}}}
	p48.items = []parser{&p47}
	p56.items = []parser{&p55, &p48}
	var p61 = sequenceParser{id: 61, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{62, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p58 = sequenceParser{id: 58, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p57 = charParser{id: 57, chars: []rune{48}}
	p58.items = []parser{&p57}
	var p60 = sequenceParser{id: 60, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p59 = charParser{id: 59, chars: []rune{120, 88}}
	p60.items = []parser{&p59}
	var p50 = sequenceParser{id: 50, commit: 66, name: "hexa-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p49 = charParser{id: 49, ranges: [][]rune{{48, 57}, {97, 102}, {65, 70}}}
	p50.items = []parser{&p49}
	p61.items = []parser{&p58, &p60, &p50}
	p62.options = []parser{&p53, &p56, &p61}
	var p75 = choiceParser{id: 75, commit: 72, name: "float", generalizations: []int{280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p70 = sequenceParser{id: 70, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{75, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p69 = sequenceParser{id: 69, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p68 = charParser{id: 68, chars: []rune{46}}
	p69.items = []parser{&p68}
	var p67 = sequenceParser{id: 67, commit: 74, name: "exponent", ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var p64 = sequenceParser{id: 64, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p63 = charParser{id: 63, chars: []rune{101, 69}}
	p64.items = []parser{&p63}
	var p66 = sequenceParser{id: 66, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p65 = charParser{id: 65, chars: []rune{43, 45}}
	p66.items = []parser{&p65}
	p67.items = []parser{&p64, &p66, &p46}
	p70.items = []parser{&p46, &p69, &p46, &p67}
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{75, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p72 = sequenceParser{id: 72, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p71 = charParser{id: 71, chars: []rune{46}}
	p72.items = []parser{&p71}
	p73.items = []parser{&p72, &p46, &p67}
	var p74 = sequenceParser{id: 74, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{75, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	p74.items = []parser{&p46, &p67}
	p75.options = []parser{&p70, &p73, &p74}
	var p88 = sequenceParser{id: 88, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 116, 141, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581, 758}}
	var p77 = sequenceParser{id: 77, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p76 = charParser{id: 76, chars: []rune{34}}
	p77.items = []parser{&p76}
	var p85 = choiceParser{id: 85, commit: 10}
	var p79 = sequenceParser{id: 79, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{85}}
	var p78 = charParser{id: 78, not: true, chars: []rune{92, 34}}
	p79.items = []parser{&p78}
	var p84 = sequenceParser{id: 84, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{85}}
	var p81 = sequenceParser{id: 81, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p80 = charParser{id: 80, chars: []rune{92}}
	p81.items = []parser{&p80}
	var p83 = sequenceParser{id: 83, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p82 = charParser{id: 82, not: true}
	p83.items = []parser{&p82}
	p84.items = []parser{&p81, &p83}
	p85.options = []parser{&p79, &p84}
	var p87 = sequenceParser{id: 87, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p86 = charParser{id: 86, chars: []rune{34}}
	p87.items = []parser{&p86}
	p88.items = []parser{&p77, &p85, &p87}
	var p100 = choiceParser{id: 100, commit: 66, name: "bool", generalizations: []int{280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p93 = sequenceParser{id: 93, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p89 = charParser{id: 89, chars: []rune{116}}
	var p90 = charParser{id: 90, chars: []rune{114}}
	var p91 = charParser{id: 91, chars: []rune{117}}
	var p92 = charParser{id: 92, chars: []rune{101}}
	p93.items = []parser{&p89, &p90, &p91, &p92}
	var p99 = sequenceParser{id: 99, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 280, 116, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p94 = charParser{id: 94, chars: []rune{102}}
	var p95 = charParser{id: 95, chars: []rune{97}}
	var p96 = charParser{id: 96, chars: []rune{108}}
	var p97 = charParser{id: 97, chars: []rune{115}}
	var p98 = charParser{id: 98, chars: []rune{101}}
	p99.items = []parser{&p94, &p95, &p96, &p97, &p98}
	p100.options = []parser{&p93, &p99}
	var p105 = sequenceParser{id: 105, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{280, 116, 141, 211, 388, 346, 347, 348, 349, 350, 351, 800, 599, 810, 471, 538, 588, 581, 749}}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p102.items = []parser{&p101}
	var p104 = sequenceParser{id: 104, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p103 = charParser{id: 103, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p104.items = []parser{&p103}
	p105.items = []parser{&p102, &p104}
	var p126 = sequenceParser{id: 126, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{116, 280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p125 = sequenceParser{id: 125, commit: 66, name: "list-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p122 = sequenceParser{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p121 = charParser{id: 121, chars: []rune{91}}
	p122.items = []parser{&p121}
	var p115 = sequenceParser{id: 115, commit: 66, name: "list-sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p113 = choiceParser{id: 113, commit: 2}
	var p112 = sequenceParser{id: 112, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{113}}
	var p111 = charParser{id: 111, chars: []rune{44}}
	p112.items = []parser{&p111}
	p113.options = []parser{&p14, &p112}
	var p114 = sequenceParser{id: 114, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p114.items = []parser{&p842, &p113}
	p115.items = []parser{&p113, &p114}
	var p120 = sequenceParser{id: 120, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p116 = choiceParser{id: 116, commit: 66, name: "list-item"}
	var p110 = sequenceParser{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{116, 149, 150}}
	var p109 = sequenceParser{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	var p108 = charParser{id: 108, chars: []rune{46}}
	p109.items = []parser{&p106, &p107, &p108}
	p110.items = []parser{&p280, &p842, &p109}
	p116.options = []parser{&p388, &p110}
	var p119 = sequenceParser{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p117.items = []parser{&p115, &p842, &p116}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p118.items = []parser{&p842, &p117}
	p119.items = []parser{&p842, &p117, &p118}
	p120.items = []parser{&p116, &p119}
	var p124 = sequenceParser{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p123 = charParser{id: 123, chars: []rune{93}}
	p124.items = []parser{&p123}
	p125.items = []parser{&p122, &p842, &p115, &p842, &p120, &p842, &p115, &p842, &p124}
	p126.items = []parser{&p125}
	var p131 = sequenceParser{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p128 = sequenceParser{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p127 = charParser{id: 127, chars: []rune{126}}
	p128.items = []parser{&p127}
	var p130 = sequenceParser{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p129.items = []parser{&p842, &p14}
	p130.items = []parser{&p842, &p14, &p129}
	p131.items = []parser{&p128, &p130, &p842, &p125}
	var p160 = sequenceParser{id: 160, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p159 = sequenceParser{id: 159, commit: 66, name: "struct-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p156 = sequenceParser{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p155 = charParser{id: 155, chars: []rune{123}}
	p156.items = []parser{&p155}
	var p154 = sequenceParser{id: 154, commit: 66, name: "entry-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p149 = choiceParser{id: 149, commit: 2}
	var p148 = sequenceParser{id: 148, commit: 64, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{149, 150}}
	var p141 = choiceParser{id: 141, commit: 2}
	var p140 = sequenceParser{id: 140, commit: 64, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{141}}
	var p133 = sequenceParser{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p132 = charParser{id: 132, chars: []rune{91}}
	p133.items = []parser{&p132}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p136 = sequenceParser{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p136.items = []parser{&p842, &p14}
	p137.items = []parser{&p842, &p14, &p136}
	var p139 = sequenceParser{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p138.items = []parser{&p842, &p14}
	p139.items = []parser{&p842, &p14, &p138}
	var p135 = sequenceParser{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p134 = charParser{id: 134, chars: []rune{93}}
	p135.items = []parser{&p134}
	p140.items = []parser{&p133, &p137, &p842, &p388, &p139, &p842, &p135}
	p141.options = []parser{&p105, &p88, &p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p842, &p14}
	p145.items = []parser{&p842, &p14, &p144}
	var p143 = sequenceParser{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p142 = charParser{id: 142, chars: []rune{58}}
	p143.items = []parser{&p142}
	var p147 = sequenceParser{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p146.items = []parser{&p842, &p14}
	p147.items = []parser{&p842, &p14, &p146}
	p148.items = []parser{&p141, &p145, &p842, &p143, &p147, &p842, &p388}
	p149.options = []parser{&p148, &p110}
	var p153 = sequenceParser{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p150 = choiceParser{id: 150, commit: 2}
	p150.options = []parser{&p148, &p110}
	p151.items = []parser{&p115, &p842, &p150}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p152.items = []parser{&p842, &p151}
	p153.items = []parser{&p842, &p151, &p152}
	p154.items = []parser{&p149, &p153}
	var p158 = sequenceParser{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p157 = charParser{id: 157, chars: []rune{125}}
	p158.items = []parser{&p157}
	p159.items = []parser{&p156, &p842, &p115, &p842, &p154, &p842, &p115, &p842, &p158}
	p160.items = []parser{&p159}
	var p165 = sequenceParser{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p162 = sequenceParser{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p161 = charParser{id: 161, chars: []rune{126}}
	p162.items = []parser{&p161}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p163.items = []parser{&p842, &p14}
	p164.items = []parser{&p842, &p14, &p163}
	p165.items = []parser{&p162, &p164, &p842, &p159}
	var p178 = choiceParser{id: 178, commit: 64, name: "channel", generalizations: []int{280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p168 = sequenceParser{id: 168, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{178, 280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p166 = charParser{id: 166, chars: []rune{60}}
	var p167 = charParser{id: 167, chars: []rune{62}}
	p168.items = []parser{&p166, &p167}
	var p177 = sequenceParser{id: 177, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{178, 280, 211, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p170 = sequenceParser{id: 170, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p169 = charParser{id: 169, chars: []rune{60}}
	p170.items = []parser{&p169}
	var p174 = sequenceParser{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p173 = sequenceParser{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p173.items = []parser{&p842, &p14}
	p174.items = []parser{&p842, &p14, &p173}
	var p176 = sequenceParser{id: 176, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p175 = sequenceParser{id: 175, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p175.items = []parser{&p842, &p14}
	p176.items = []parser{&p842, &p14, &p175}
	var p172 = sequenceParser{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p171 = charParser{id: 171, chars: []rune{62}}
	p172.items = []parser{&p171}
	p177.items = []parser{&p170, &p174, &p842, &p388, &p176, &p842, &p172}
	p178.options = []parser{&p168, &p177}
	var p220 = sequenceParser{id: 220, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{211, 280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p217 = sequenceParser{id: 217, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p215 = charParser{id: 215, chars: []rune{102}}
	var p216 = charParser{id: 216, chars: []rune{110}}
	p217.items = []parser{&p215, &p216}
	var p219 = sequenceParser{id: 219, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p218 = sequenceParser{id: 218, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p218.items = []parser{&p842, &p14}
	p219.items = []parser{&p842, &p14, &p218}
	var p214 = sequenceParser{id: 214, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p207 = sequenceParser{id: 207, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p206 = charParser{id: 206, chars: []rune{40}}
	p207.items = []parser{&p206}
	var p182 = sequenceParser{id: 182, commit: 66, name: "argument-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p181 = sequenceParser{id: 181, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p179 = sequenceParser{id: 179, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p179.items = []parser{&p115, &p842, &p105}
	var p180 = sequenceParser{id: 180, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p180.items = []parser{&p842, &p179}
	p181.items = []parser{&p842, &p179, &p180}
	p182.items = []parser{&p105, &p181}
	var p208 = sequenceParser{id: 208, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p189 = sequenceParser{id: 189, commit: 64, name: "collect-argument", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p186 = sequenceParser{id: 186, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p183 = charParser{id: 183, chars: []rune{46}}
	var p184 = charParser{id: 184, chars: []rune{46}}
	var p185 = charParser{id: 185, chars: []rune{46}}
	p186.items = []parser{&p183, &p184, &p185}
	var p188 = sequenceParser{id: 188, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p187 = sequenceParser{id: 187, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p187.items = []parser{&p842, &p14}
	p188.items = []parser{&p842, &p14, &p187}
	p189.items = []parser{&p186, &p188, &p842, &p105}
	p208.items = []parser{&p115, &p842, &p189}
	var p210 = sequenceParser{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p209 = charParser{id: 209, chars: []rune{41}}
	p210.items = []parser{&p209}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p212 = sequenceParser{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p212.items = []parser{&p842, &p14}
	p213.items = []parser{&p842, &p14, &p212}
	var p211 = choiceParser{id: 211, commit: 2}
	var p205 = sequenceParser{id: 205, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{211}}
	var p202 = sequenceParser{id: 202, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p201 = charParser{id: 201, chars: []rune{123}}
	p202.items = []parser{&p201}
	var p204 = sequenceParser{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p203 = charParser{id: 203, chars: []rune{125}}
	p204.items = []parser{&p203}
	p205.items = []parser{&p202, &p842, &p824, &p842, &p828, &p842, &p824, &p842, &p204}
	p211.options = []parser{&p800, &p205}
	p214.items = []parser{&p207, &p842, &p115, &p842, &p182, &p842, &p208, &p842, &p115, &p842, &p210, &p213, &p842, &p211}
	p220.items = []parser{&p217, &p219, &p842, &p214}
	var p230 = sequenceParser{id: 230, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p223 = sequenceParser{id: 223, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p221 = charParser{id: 221, chars: []rune{102}}
	var p222 = charParser{id: 222, chars: []rune{110}}
	p223.items = []parser{&p221, &p222}
	var p227 = sequenceParser{id: 227, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p226 = sequenceParser{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p226.items = []parser{&p842, &p14}
	p227.items = []parser{&p842, &p14, &p226}
	var p225 = sequenceParser{id: 225, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p224 = charParser{id: 224, chars: []rune{126}}
	p225.items = []parser{&p224}
	var p229 = sequenceParser{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p228 = sequenceParser{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p228.items = []parser{&p842, &p14}
	p229.items = []parser{&p842, &p14, &p228}
	p230.items = []parser{&p223, &p227, &p842, &p225, &p229, &p842, &p214}
	var p258 = choiceParser{id: 258, commit: 64, name: "expression-indexer", generalizations: []int{280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p248 = sequenceParser{id: 248, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{258, 280, 388, 346, 347, 348, 349, 350, 351, 800, 599, 810, 471, 538, 588, 581}}
	var p241 = sequenceParser{id: 241, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p240 = charParser{id: 240, chars: []rune{91}}
	p241.items = []parser{&p240}
	var p245 = sequenceParser{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p244 = sequenceParser{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p244.items = []parser{&p842, &p14}
	p245.items = []parser{&p842, &p14, &p244}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p246 = sequenceParser{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p246.items = []parser{&p842, &p14}
	p247.items = []parser{&p842, &p14, &p246}
	var p243 = sequenceParser{id: 243, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p242 = charParser{id: 242, chars: []rune{93}}
	p243.items = []parser{&p242}
	p248.items = []parser{&p280, &p842, &p241, &p245, &p842, &p388, &p247, &p842, &p243}
	var p257 = sequenceParser{id: 257, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{258, 280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p250 = sequenceParser{id: 250, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p249 = charParser{id: 249, chars: []rune{91}}
	p250.items = []parser{&p249}
	var p254 = sequenceParser{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p253 = sequenceParser{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p253.items = []parser{&p842, &p14}
	p254.items = []parser{&p842, &p14, &p253}
	var p239 = sequenceParser{id: 239, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{581, 587, 588}}
	var p231 = sequenceParser{id: 231, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p231.items = []parser{&p388}
	var p236 = sequenceParser{id: 236, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p235 = sequenceParser{id: 235, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p235.items = []parser{&p842, &p14}
	p236.items = []parser{&p842, &p14, &p235}
	var p234 = sequenceParser{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p233 = charParser{id: 233, chars: []rune{58}}
	p234.items = []parser{&p233}
	var p238 = sequenceParser{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p237 = sequenceParser{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p237.items = []parser{&p842, &p14}
	p238.items = []parser{&p842, &p14, &p237}
	var p232 = sequenceParser{id: 232, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p232.items = []parser{&p388}
	p239.items = []parser{&p231, &p236, &p842, &p234, &p238, &p842, &p232}
	var p256 = sequenceParser{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p255 = sequenceParser{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p255.items = []parser{&p842, &p14}
	p256.items = []parser{&p842, &p14, &p255}
	var p252 = sequenceParser{id: 252, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p251 = charParser{id: 251, chars: []rune{93}}
	p252.items = []parser{&p251}
	p257.items = []parser{&p280, &p842, &p250, &p254, &p842, &p239, &p256, &p842, &p252}
	p258.options = []parser{&p248, &p257}
	var p265 = sequenceParser{id: 265, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 388, 346, 347, 348, 349, 350, 351, 800, 599, 810, 471, 538, 588, 581}}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p261 = sequenceParser{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p261.items = []parser{&p842, &p14}
	p262.items = []parser{&p842, &p14, &p261}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{46}}
	p260.items = []parser{&p259}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p263.items = []parser{&p842, &p14}
	p264.items = []parser{&p842, &p14, &p263}
	p265.items = []parser{&p280, &p262, &p842, &p260, &p264, &p842, &p105}
	var p270 = sequenceParser{id: 270, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p267 = sequenceParser{id: 267, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p266 = charParser{id: 266, chars: []rune{40}}
	p267.items = []parser{&p266}
	var p269 = sequenceParser{id: 269, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p268 = charParser{id: 268, chars: []rune{41}}
	p269.items = []parser{&p268}
	p270.items = []parser{&p280, &p842, &p267, &p842, &p115, &p842, &p120, &p842, &p115, &p842, &p269}
	var p486 = sequenceParser{id: 486, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 510, 538, 588, 581}}
	var p485 = sequenceParser{id: 485, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p483 = charParser{id: 483, chars: []rune{60}}
	var p484 = charParser{id: 484, chars: []rune{45}}
	p485.items = []parser{&p483, &p484}
	p486.items = []parser{&p485, &p842, &p280}
	var p279 = sequenceParser{id: 279, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{280, 388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p272 = sequenceParser{id: 272, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p271 = charParser{id: 271, chars: []rune{40}}
	p272.items = []parser{&p271}
	var p276 = sequenceParser{id: 276, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p275 = sequenceParser{id: 275, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p275.items = []parser{&p842, &p14}
	p276.items = []parser{&p842, &p14, &p275}
	var p278 = sequenceParser{id: 278, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p277 = sequenceParser{id: 277, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p277.items = []parser{&p842, &p14}
	p278.items = []parser{&p842, &p14, &p277}
	var p274 = sequenceParser{id: 274, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p273 = charParser{id: 273, chars: []rune{41}}
	p274.items = []parser{&p273}
	p279.items = []parser{&p272, &p276, &p842, &p388, &p278, &p842, &p274}
	p280.options = []parser{&p62, &p75, &p88, &p100, &p105, &p126, &p131, &p160, &p165, &p178, &p220, &p230, &p258, &p265, &p270, &p486, &p279}
	var p340 = sequenceParser{id: 340, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{388, 346, 347, 348, 349, 350, 351, 800, 810, 471, 538, 588, 581}}
	var p339 = choiceParser{id: 339, commit: 66, name: "unary-operator"}
	var p299 = sequenceParser{id: 299, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{339}}
	var p298 = charParser{id: 298, chars: []rune{43}}
	p299.items = []parser{&p298}
	var p301 = sequenceParser{id: 301, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{339}}
	var p300 = charParser{id: 300, chars: []rune{45}}
	p301.items = []parser{&p300}
	var p282 = sequenceParser{id: 282, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{339}}
	var p281 = charParser{id: 281, chars: []rune{94}}
	p282.items = []parser{&p281}
	p339.options = []parser{&p299, &p301, &p282}
	p340.items = []parser{&p339, &p842, &p280}
	var p374 = choiceParser{id: 374, commit: 66, name: "binary-expression", generalizations: []int{388, 800, 810, 471, 538, 588, 581}}
	var p354 = sequenceParser{id: 354, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 347, 348, 349, 350, 351, 388, 800, 810, 471, 538, 588, 581}}
	var p346 = choiceParser{id: 346, commit: 66, name: "operand0", generalizations: []int{347, 348, 349, 350, 351}}
	p346.options = []parser{&p280, &p340}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p341 = choiceParser{id: 341, commit: 66, name: "binary-op0"}
	var p284 = sequenceParser{id: 284, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{341}}
	var p283 = charParser{id: 283, chars: []rune{38}}
	p284.items = []parser{&p283}
	var p291 = sequenceParser{id: 291, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{341}}
	var p289 = charParser{id: 289, chars: []rune{38}}
	var p290 = charParser{id: 290, chars: []rune{94}}
	p291.items = []parser{&p289, &p290}
	var p294 = sequenceParser{id: 294, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{341}}
	var p292 = charParser{id: 292, chars: []rune{60}}
	var p293 = charParser{id: 293, chars: []rune{60}}
	p294.items = []parser{&p292, &p293}
	var p297 = sequenceParser{id: 297, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{341}}
	var p295 = charParser{id: 295, chars: []rune{62}}
	var p296 = charParser{id: 296, chars: []rune{62}}
	p297.items = []parser{&p295, &p296}
	var p303 = sequenceParser{id: 303, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{341}}
	var p302 = charParser{id: 302, chars: []rune{42}}
	p303.items = []parser{&p302}
	var p305 = sequenceParser{id: 305, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{341}}
	var p304 = charParser{id: 304, chars: []rune{47}}
	p305.items = []parser{&p304}
	var p307 = sequenceParser{id: 307, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{341}}
	var p306 = charParser{id: 306, chars: []rune{37}}
	p307.items = []parser{&p306}
	p341.options = []parser{&p284, &p291, &p294, &p297, &p303, &p305, &p307}
	p352.items = []parser{&p341, &p842, &p346}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p353.items = []parser{&p842, &p352}
	p354.items = []parser{&p346, &p842, &p352, &p353}
	var p357 = sequenceParser{id: 357, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 348, 349, 350, 351, 388, 800, 810, 471, 538, 588, 581}}
	var p347 = choiceParser{id: 347, commit: 66, name: "operand1", generalizations: []int{348, 349, 350, 351}}
	p347.options = []parser{&p346, &p354}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p342 = choiceParser{id: 342, commit: 66, name: "binary-op1"}
	var p286 = sequenceParser{id: 286, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{342}}
	var p285 = charParser{id: 285, chars: []rune{124}}
	p286.items = []parser{&p285}
	var p288 = sequenceParser{id: 288, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{342}}
	var p287 = charParser{id: 287, chars: []rune{94}}
	p288.items = []parser{&p287}
	var p309 = sequenceParser{id: 309, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{342}}
	var p308 = charParser{id: 308, chars: []rune{43}}
	p309.items = []parser{&p308}
	var p311 = sequenceParser{id: 311, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{342}}
	var p310 = charParser{id: 310, chars: []rune{45}}
	p311.items = []parser{&p310}
	p342.options = []parser{&p286, &p288, &p309, &p311}
	p355.items = []parser{&p342, &p842, &p347}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p356.items = []parser{&p842, &p355}
	p357.items = []parser{&p347, &p842, &p355, &p356}
	var p360 = sequenceParser{id: 360, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 349, 350, 351, 388, 800, 810, 471, 538, 588, 581}}
	var p348 = choiceParser{id: 348, commit: 66, name: "operand2", generalizations: []int{349, 350, 351}}
	p348.options = []parser{&p347, &p357}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p343 = choiceParser{id: 343, commit: 66, name: "binary-op2"}
	var p316 = sequenceParser{id: 316, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{343}}
	var p314 = charParser{id: 314, chars: []rune{61}}
	var p315 = charParser{id: 315, chars: []rune{61}}
	p316.items = []parser{&p314, &p315}
	var p319 = sequenceParser{id: 319, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{343}}
	var p317 = charParser{id: 317, chars: []rune{33}}
	var p318 = charParser{id: 318, chars: []rune{61}}
	p319.items = []parser{&p317, &p318}
	var p321 = sequenceParser{id: 321, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{343}}
	var p320 = charParser{id: 320, chars: []rune{60}}
	p321.items = []parser{&p320}
	var p324 = sequenceParser{id: 324, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{343}}
	var p322 = charParser{id: 322, chars: []rune{60}}
	var p323 = charParser{id: 323, chars: []rune{61}}
	p324.items = []parser{&p322, &p323}
	var p326 = sequenceParser{id: 326, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{343}}
	var p325 = charParser{id: 325, chars: []rune{62}}
	p326.items = []parser{&p325}
	var p329 = sequenceParser{id: 329, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{343}}
	var p327 = charParser{id: 327, chars: []rune{62}}
	var p328 = charParser{id: 328, chars: []rune{61}}
	p329.items = []parser{&p327, &p328}
	p343.options = []parser{&p316, &p319, &p321, &p324, &p326, &p329}
	p358.items = []parser{&p343, &p842, &p348}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p842, &p358}
	p360.items = []parser{&p348, &p842, &p358, &p359}
	var p363 = sequenceParser{id: 363, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 350, 351, 388, 800, 810, 471, 538, 588, 581}}
	var p349 = choiceParser{id: 349, commit: 66, name: "operand3", generalizations: []int{350, 351}}
	p349.options = []parser{&p348, &p360}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p344 = sequenceParser{id: 344, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p332 = sequenceParser{id: 332, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p330 = charParser{id: 330, chars: []rune{38}}
	var p331 = charParser{id: 331, chars: []rune{38}}
	p332.items = []parser{&p330, &p331}
	p344.items = []parser{&p332}
	p361.items = []parser{&p344, &p842, &p349}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p362.items = []parser{&p842, &p361}
	p363.items = []parser{&p349, &p842, &p361, &p362}
	var p366 = sequenceParser{id: 366, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 351, 388, 800, 810, 471, 538, 588, 581}}
	var p350 = choiceParser{id: 350, commit: 66, name: "operand4", generalizations: []int{351}}
	p350.options = []parser{&p349, &p363}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p345 = sequenceParser{id: 345, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p335 = sequenceParser{id: 335, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p333 = charParser{id: 333, chars: []rune{124}}
	var p334 = charParser{id: 334, chars: []rune{124}}
	p335.items = []parser{&p333, &p334}
	p345.items = []parser{&p335}
	p364.items = []parser{&p345, &p842, &p350}
	var p365 = sequenceParser{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p365.items = []parser{&p842, &p364}
	p366.items = []parser{&p350, &p842, &p364, &p365}
	var p373 = sequenceParser{id: 373, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{374, 388, 800, 810, 471, 538, 588, 581}}
	var p351 = choiceParser{id: 351, commit: 66, name: "operand5"}
	p351.options = []parser{&p350, &p366}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p367.items = []parser{&p842, &p14}
	p368.items = []parser{&p14, &p367}
	var p338 = sequenceParser{id: 338, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p336 = charParser{id: 336, chars: []rune{45}}
	var p337 = charParser{id: 337, chars: []rune{62}}
	p338.items = []parser{&p336, &p337}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p369.items = []parser{&p842, &p14}
	p370.items = []parser{&p842, &p14, &p369}
	p371.items = []parser{&p368, &p842, &p338, &p370, &p842, &p351}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p372.items = []parser{&p842, &p371}
	p373.items = []parser{&p351, &p842, &p371, &p372}
	p374.options = []parser{&p354, &p357, &p360, &p363, &p366, &p373}
	var p387 = sequenceParser{id: 387, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{388, 800, 810, 471, 538, 588, 581}}
	var p380 = sequenceParser{id: 380, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p379 = sequenceParser{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p379.items = []parser{&p842, &p14}
	p380.items = []parser{&p842, &p14, &p379}
	var p376 = sequenceParser{id: 376, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p375 = charParser{id: 375, chars: []rune{63}}
	p376.items = []parser{&p375}
	var p382 = sequenceParser{id: 382, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p381 = sequenceParser{id: 381, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p381.items = []parser{&p842, &p14}
	p382.items = []parser{&p842, &p14, &p381}
	var p384 = sequenceParser{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p383 = sequenceParser{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p383.items = []parser{&p842, &p14}
	p384.items = []parser{&p842, &p14, &p383}
	var p378 = sequenceParser{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p377 = charParser{id: 377, chars: []rune{58}}
	p378.items = []parser{&p377}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p385.items = []parser{&p842, &p14}
	p386.items = []parser{&p842, &p14, &p385}
	p387.items = []parser{&p388, &p380, &p842, &p376, &p382, &p842, &p388, &p384, &p842, &p378, &p386, &p842, &p388}
	p388.options = []parser{&p280, &p340, &p374, &p387}
	var p509 = sequenceParser{id: 509, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{800, 810, 471, 510, 538}}
	var p508 = sequenceParser{id: 508, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p506 = charParser{id: 506, chars: []rune{60}}
	var p507 = charParser{id: 507, chars: []rune{45}}
	p508.items = []parser{&p506, &p507}
	p509.items = []parser{&p280, &p842, &p508, &p842, &p280}
	var p553 = sequenceParser{id: 553, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{800, 810, 471, 538}}
	var p550 = sequenceParser{id: 550, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p548 = charParser{id: 548, chars: []rune{103}}
	var p549 = charParser{id: 549, chars: []rune{111}}
	p550.items = []parser{&p548, &p549}
	var p552 = sequenceParser{id: 552, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p551 = sequenceParser{id: 551, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p551.items = []parser{&p842, &p14}
	p552.items = []parser{&p842, &p14, &p551}
	p553.items = []parser{&p550, &p552, &p842, &p270}
	var p562 = sequenceParser{id: 562, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{800, 810, 471, 538}}
	var p559 = sequenceParser{id: 559, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p554 = charParser{id: 554, chars: []rune{100}}
	var p555 = charParser{id: 555, chars: []rune{101}}
	var p556 = charParser{id: 556, chars: []rune{102}}
	var p557 = charParser{id: 557, chars: []rune{101}}
	var p558 = charParser{id: 558, chars: []rune{114}}
	p559.items = []parser{&p554, &p555, &p556, &p557, &p558}
	var p561 = sequenceParser{id: 561, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p560 = sequenceParser{id: 560, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p560.items = []parser{&p842, &p14}
	p561.items = []parser{&p842, &p14, &p560}
	p562.items = []parser{&p559, &p561, &p842, &p270}
	var p637 = choiceParser{id: 637, commit: 64, name: "assignment", generalizations: []int{800, 810, 471, 538}}
	var p618 = sequenceParser{id: 618, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{637, 800, 810, 471, 538}}
	var p615 = sequenceParser{id: 615, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p612 = charParser{id: 612, chars: []rune{115}}
	var p613 = charParser{id: 613, chars: []rune{101}}
	var p614 = charParser{id: 614, chars: []rune{116}}
	p615.items = []parser{&p612, &p613, &p614}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p616.items = []parser{&p842, &p14}
	p617.items = []parser{&p842, &p14, &p616}
	var p607 = sequenceParser{id: 607, commit: 64, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p599 = choiceParser{id: 599, commit: 66, name: "assignable"}
	p599.options = []parser{&p105, &p248, &p265}
	var p604 = sequenceParser{id: 604, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p603 = sequenceParser{id: 603, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p602 = sequenceParser{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p602.items = []parser{&p842, &p14}
	p603.items = []parser{&p14, &p602}
	var p601 = sequenceParser{id: 601, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p600 = charParser{id: 600, chars: []rune{61}}
	p601.items = []parser{&p600}
	p604.items = []parser{&p603, &p842, &p601}
	var p606 = sequenceParser{id: 606, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p605 = sequenceParser{id: 605, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p605.items = []parser{&p842, &p14}
	p606.items = []parser{&p842, &p14, &p605}
	p607.items = []parser{&p599, &p842, &p604, &p606, &p842, &p388}
	p618.items = []parser{&p615, &p617, &p842, &p607}
	var p625 = sequenceParser{id: 625, commit: 66, name: "assign-equal", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{637, 800, 810, 471, 538}}
	var p622 = sequenceParser{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p621 = sequenceParser{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p621.items = []parser{&p842, &p14}
	p622.items = []parser{&p842, &p14, &p621}
	var p620 = sequenceParser{id: 620, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p619 = charParser{id: 619, chars: []rune{61}}
	p620.items = []parser{&p619}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p623 = sequenceParser{id: 623, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p623.items = []parser{&p842, &p14}
	p624.items = []parser{&p842, &p14, &p623}
	p625.items = []parser{&p599, &p622, &p842, &p620, &p624, &p842, &p388}
	var p636 = sequenceParser{id: 636, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{637, 800, 810, 471, 538}}
	var p629 = sequenceParser{id: 629, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p626 = charParser{id: 626, chars: []rune{115}}
	var p627 = charParser{id: 627, chars: []rune{101}}
	var p628 = charParser{id: 628, chars: []rune{116}}
	p629.items = []parser{&p626, &p627, &p628}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p634.items = []parser{&p842, &p14}
	p635.items = []parser{&p842, &p14, &p634}
	var p631 = sequenceParser{id: 631, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p630 = charParser{id: 630, chars: []rune{40}}
	p631.items = []parser{&p630}
	var p611 = sequenceParser{id: 611, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p610 = sequenceParser{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p608 = sequenceParser{id: 608, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p608.items = []parser{&p115, &p842, &p607}
	var p609 = sequenceParser{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p609.items = []parser{&p842, &p608}
	p610.items = []parser{&p842, &p608, &p609}
	p611.items = []parser{&p607, &p610}
	var p633 = sequenceParser{id: 633, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p632 = charParser{id: 632, chars: []rune{41}}
	p633.items = []parser{&p632}
	p636.items = []parser{&p629, &p635, &p842, &p631, &p842, &p115, &p842, &p611, &p842, &p115, &p842, &p633}
	p637.options = []parser{&p618, &p625, &p636}
	var p809 = sequenceParser{id: 809, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{800, 810, 471, 538}}
	var p802 = sequenceParser{id: 802, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p801 = charParser{id: 801, chars: []rune{40}}
	p802.items = []parser{&p801}
	var p806 = sequenceParser{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p805 = sequenceParser{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p805.items = []parser{&p842, &p14}
	p806.items = []parser{&p842, &p14, &p805}
	var p808 = sequenceParser{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p807 = sequenceParser{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p807.items = []parser{&p842, &p14}
	p808.items = []parser{&p842, &p14, &p807}
	var p804 = sequenceParser{id: 804, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p803 = charParser{id: 803, chars: []rune{41}}
	p804.items = []parser{&p803}
	p809.items = []parser{&p802, &p806, &p842, &p800, &p808, &p842, &p804}
	p800.options = []parser{&p388, &p509, &p553, &p562, &p637, &p809}
	var p200 = sequenceParser{id: 200, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{810, 471, 538}}
	var p196 = sequenceParser{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p190 = charParser{id: 190, chars: []rune{114}}
	var p191 = charParser{id: 191, chars: []rune{101}}
	var p192 = charParser{id: 192, chars: []rune{116}}
	var p193 = charParser{id: 193, chars: []rune{117}}
	var p194 = charParser{id: 194, chars: []rune{114}}
	var p195 = charParser{id: 195, chars: []rune{110}}
	p196.items = []parser{&p190, &p191, &p192, &p193, &p194, &p195}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p198 = sequenceParser{id: 198, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p197 = sequenceParser{id: 197, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p197.items = []parser{&p842, &p14}
	p198.items = []parser{&p14, &p197}
	p199.items = []parser{&p198, &p842, &p388}
	p200.items = []parser{&p196, &p842, &p199}
	var p425 = sequenceParser{id: 425, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{810, 471, 538}}
	var p391 = sequenceParser{id: 391, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p389 = charParser{id: 389, chars: []rune{105}}
	var p390 = charParser{id: 390, chars: []rune{102}}
	p391.items = []parser{&p389, &p390}
	var p420 = sequenceParser{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p419 = sequenceParser{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p419.items = []parser{&p842, &p14}
	p420.items = []parser{&p842, &p14, &p419}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p421 = sequenceParser{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p421.items = []parser{&p842, &p14}
	p422.items = []parser{&p842, &p14, &p421}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p401 = sequenceParser{id: 401, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p400 = sequenceParser{id: 400, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p400.items = []parser{&p842, &p14}
	p401.items = []parser{&p14, &p400}
	var p396 = sequenceParser{id: 396, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p392 = charParser{id: 392, chars: []rune{101}}
	var p393 = charParser{id: 393, chars: []rune{108}}
	var p394 = charParser{id: 394, chars: []rune{115}}
	var p395 = charParser{id: 395, chars: []rune{101}}
	p396.items = []parser{&p392, &p393, &p394, &p395}
	var p403 = sequenceParser{id: 403, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p402 = sequenceParser{id: 402, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p402.items = []parser{&p842, &p14}
	p403.items = []parser{&p842, &p14, &p402}
	var p399 = sequenceParser{id: 399, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p397 = charParser{id: 397, chars: []rune{105}}
	var p398 = charParser{id: 398, chars: []rune{102}}
	p399.items = []parser{&p397, &p398}
	var p405 = sequenceParser{id: 405, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p404 = sequenceParser{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p404.items = []parser{&p842, &p14}
	p405.items = []parser{&p842, &p14, &p404}
	var p407 = sequenceParser{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p406 = sequenceParser{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p406.items = []parser{&p842, &p14}
	p407.items = []parser{&p842, &p14, &p406}
	p408.items = []parser{&p401, &p842, &p396, &p403, &p842, &p399, &p405, &p842, &p388, &p407, &p842, &p205}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p423.items = []parser{&p842, &p408}
	p424.items = []parser{&p842, &p408, &p423}
	var p418 = sequenceParser{id: 418, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p414.items = []parser{&p842, &p14}
	p415.items = []parser{&p14, &p414}
	var p413 = sequenceParser{id: 413, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p409 = charParser{id: 409, chars: []rune{101}}
	var p410 = charParser{id: 410, chars: []rune{108}}
	var p411 = charParser{id: 411, chars: []rune{115}}
	var p412 = charParser{id: 412, chars: []rune{101}}
	p413.items = []parser{&p409, &p410, &p411, &p412}
	var p417 = sequenceParser{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p416 = sequenceParser{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p416.items = []parser{&p842, &p14}
	p417.items = []parser{&p842, &p14, &p416}
	p418.items = []parser{&p415, &p842, &p413, &p417, &p842, &p205}
	p425.items = []parser{&p391, &p420, &p842, &p388, &p422, &p842, &p205, &p424, &p842, &p418}
	var p482 = sequenceParser{id: 482, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{471, 810, 538}}
	var p467 = sequenceParser{id: 467, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p461 = charParser{id: 461, chars: []rune{115}}
	var p462 = charParser{id: 462, chars: []rune{119}}
	var p463 = charParser{id: 463, chars: []rune{105}}
	var p464 = charParser{id: 464, chars: []rune{116}}
	var p465 = charParser{id: 465, chars: []rune{99}}
	var p466 = charParser{id: 466, chars: []rune{104}}
	p467.items = []parser{&p461, &p462, &p463, &p464, &p465, &p466}
	var p479 = sequenceParser{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p478 = sequenceParser{id: 478, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p478.items = []parser{&p842, &p14}
	p479.items = []parser{&p842, &p14, &p478}
	var p481 = sequenceParser{id: 481, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p480 = sequenceParser{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p480.items = []parser{&p842, &p14}
	p481.items = []parser{&p842, &p14, &p480}
	var p469 = sequenceParser{id: 469, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p468 = charParser{id: 468, chars: []rune{123}}
	p469.items = []parser{&p468}
	var p475 = sequenceParser{id: 475, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p470 = choiceParser{id: 470, commit: 2}
	var p460 = sequenceParser{id: 460, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{470, 471}}
	var p455 = sequenceParser{id: 455, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p448 = sequenceParser{id: 448, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p444 = charParser{id: 444, chars: []rune{99}}
	var p445 = charParser{id: 445, chars: []rune{97}}
	var p446 = charParser{id: 446, chars: []rune{115}}
	var p447 = charParser{id: 447, chars: []rune{101}}
	p448.items = []parser{&p444, &p445, &p446, &p447}
	var p452 = sequenceParser{id: 452, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p451 = sequenceParser{id: 451, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p451.items = []parser{&p842, &p14}
	p452.items = []parser{&p842, &p14, &p451}
	var p454 = sequenceParser{id: 454, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p453 = sequenceParser{id: 453, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p453.items = []parser{&p842, &p14}
	p454.items = []parser{&p842, &p14, &p453}
	var p450 = sequenceParser{id: 450, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p449 = charParser{id: 449, chars: []rune{58}}
	p450.items = []parser{&p449}
	p455.items = []parser{&p448, &p452, &p842, &p388, &p454, &p842, &p450}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p457 = sequenceParser{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p456 = charParser{id: 456, chars: []rune{59}}
	p457.items = []parser{&p456}
	var p458 = sequenceParser{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p458.items = []parser{&p842, &p457}
	p459.items = []parser{&p842, &p457, &p458}
	p460.items = []parser{&p455, &p459, &p842, &p810}
	var p443 = sequenceParser{id: 443, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{470, 471, 537, 538}}
	var p438 = sequenceParser{id: 438, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p433 = sequenceParser{id: 433, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p426 = charParser{id: 426, chars: []rune{100}}
	var p427 = charParser{id: 427, chars: []rune{101}}
	var p428 = charParser{id: 428, chars: []rune{102}}
	var p429 = charParser{id: 429, chars: []rune{97}}
	var p430 = charParser{id: 430, chars: []rune{117}}
	var p431 = charParser{id: 431, chars: []rune{108}}
	var p432 = charParser{id: 432, chars: []rune{116}}
	p433.items = []parser{&p426, &p427, &p428, &p429, &p430, &p431, &p432}
	var p437 = sequenceParser{id: 437, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p436 = sequenceParser{id: 436, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p436.items = []parser{&p842, &p14}
	p437.items = []parser{&p842, &p14, &p436}
	var p435 = sequenceParser{id: 435, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p434 = charParser{id: 434, chars: []rune{58}}
	p435.items = []parser{&p434}
	p438.items = []parser{&p433, &p437, &p842, &p435}
	var p442 = sequenceParser{id: 442, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p440 = sequenceParser{id: 440, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p439 = charParser{id: 439, chars: []rune{59}}
	p440.items = []parser{&p439}
	var p441 = sequenceParser{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p441.items = []parser{&p842, &p440}
	p442.items = []parser{&p842, &p440, &p441}
	p443.items = []parser{&p438, &p442, &p842, &p810}
	p470.options = []parser{&p460, &p443}
	var p474 = sequenceParser{id: 474, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p472 = sequenceParser{id: 472, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p471 = choiceParser{id: 471, commit: 2}
	p471.options = []parser{&p460, &p443, &p810}
	p472.items = []parser{&p824, &p842, &p471}
	var p473 = sequenceParser{id: 473, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p473.items = []parser{&p842, &p472}
	p474.items = []parser{&p842, &p472, &p473}
	p475.items = []parser{&p470, &p474}
	var p477 = sequenceParser{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p476 = charParser{id: 476, chars: []rune{125}}
	p477.items = []parser{&p476}
	p482.items = []parser{&p467, &p479, &p842, &p388, &p481, &p842, &p469, &p842, &p824, &p842, &p475, &p842, &p824, &p842, &p477}
	var p547 = sequenceParser{id: 547, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{538, 810}}
	var p534 = sequenceParser{id: 534, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p528 = charParser{id: 528, chars: []rune{115}}
	var p529 = charParser{id: 529, chars: []rune{101}}
	var p530 = charParser{id: 530, chars: []rune{108}}
	var p531 = charParser{id: 531, chars: []rune{101}}
	var p532 = charParser{id: 532, chars: []rune{99}}
	var p533 = charParser{id: 533, chars: []rune{116}}
	p534.items = []parser{&p528, &p529, &p530, &p531, &p532, &p533}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p545.items = []parser{&p842, &p14}
	p546.items = []parser{&p842, &p14, &p545}
	var p536 = sequenceParser{id: 536, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p535 = charParser{id: 535, chars: []rune{123}}
	p536.items = []parser{&p535}
	var p542 = sequenceParser{id: 542, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p537 = choiceParser{id: 537, commit: 2}
	var p527 = sequenceParser{id: 527, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{537, 538}}
	var p522 = sequenceParser{id: 522, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p515 = sequenceParser{id: 515, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p511 = charParser{id: 511, chars: []rune{99}}
	var p512 = charParser{id: 512, chars: []rune{97}}
	var p513 = charParser{id: 513, chars: []rune{115}}
	var p514 = charParser{id: 514, chars: []rune{101}}
	p515.items = []parser{&p511, &p512, &p513, &p514}
	var p519 = sequenceParser{id: 519, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p518 = sequenceParser{id: 518, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p518.items = []parser{&p842, &p14}
	p519.items = []parser{&p842, &p14, &p518}
	var p510 = choiceParser{id: 510, commit: 66, name: "communication"}
	var p505 = choiceParser{id: 505, commit: 66, name: "receive-statement", generalizations: []int{510}}
	var p495 = sequenceParser{id: 495, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{505, 510}}
	var p490 = sequenceParser{id: 490, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p487 = charParser{id: 487, chars: []rune{108}}
	var p488 = charParser{id: 488, chars: []rune{101}}
	var p489 = charParser{id: 489, chars: []rune{116}}
	p490.items = []parser{&p487, &p488, &p489}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p491 = sequenceParser{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p491.items = []parser{&p842, &p14}
	p492.items = []parser{&p842, &p14, &p491}
	var p494 = sequenceParser{id: 494, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p493 = sequenceParser{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p493.items = []parser{&p842, &p14}
	p494.items = []parser{&p842, &p14, &p493}
	p495.items = []parser{&p490, &p492, &p842, &p105, &p494, &p842, &p486}
	var p504 = sequenceParser{id: 504, commit: 64, name: "receive-assignment", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{505, 510}}
	var p499 = sequenceParser{id: 499, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p496 = charParser{id: 496, chars: []rune{115}}
	var p497 = charParser{id: 497, chars: []rune{101}}
	var p498 = charParser{id: 498, chars: []rune{116}}
	p499.items = []parser{&p496, &p497, &p498}
	var p501 = sequenceParser{id: 501, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p500 = sequenceParser{id: 500, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p500.items = []parser{&p842, &p14}
	p501.items = []parser{&p842, &p14, &p500}
	var p503 = sequenceParser{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p502 = sequenceParser{id: 502, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p502.items = []parser{&p842, &p14}
	p503.items = []parser{&p842, &p14, &p502}
	p504.items = []parser{&p499, &p501, &p842, &p105, &p503, &p842, &p486}
	p505.options = []parser{&p495, &p504}
	p510.options = []parser{&p486, &p505, &p509}
	var p521 = sequenceParser{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p520 = sequenceParser{id: 520, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p520.items = []parser{&p842, &p14}
	p521.items = []parser{&p842, &p14, &p520}
	var p517 = sequenceParser{id: 517, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p516 = charParser{id: 516, chars: []rune{58}}
	p517.items = []parser{&p516}
	p522.items = []parser{&p515, &p519, &p842, &p510, &p521, &p842, &p517}
	var p526 = sequenceParser{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p524 = sequenceParser{id: 524, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p523 = charParser{id: 523, chars: []rune{59}}
	p524.items = []parser{&p523}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p525.items = []parser{&p842, &p524}
	p526.items = []parser{&p842, &p524, &p525}
	p527.items = []parser{&p522, &p526, &p842, &p810}
	p537.options = []parser{&p527, &p443}
	var p541 = sequenceParser{id: 541, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p539 = sequenceParser{id: 539, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p538 = choiceParser{id: 538, commit: 2}
	p538.options = []parser{&p527, &p443, &p810}
	p539.items = []parser{&p824, &p842, &p538}
	var p540 = sequenceParser{id: 540, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p540.items = []parser{&p842, &p539}
	p541.items = []parser{&p842, &p539, &p540}
	p542.items = []parser{&p537, &p541}
	var p544 = sequenceParser{id: 544, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p543 = charParser{id: 543, chars: []rune{125}}
	p544.items = []parser{&p543}
	p547.items = []parser{&p534, &p546, &p842, &p536, &p842, &p824, &p842, &p542, &p842, &p824, &p842, &p544}
	var p568 = sequenceParser{id: 568, commit: 72, name: "break", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{810}}
	var p563 = charParser{id: 563, chars: []rune{98}}
	var p564 = charParser{id: 564, chars: []rune{114}}
	var p565 = charParser{id: 565, chars: []rune{101}}
	var p566 = charParser{id: 566, chars: []rune{97}}
	var p567 = charParser{id: 567, chars: []rune{107}}
	p568.items = []parser{&p563, &p564, &p565, &p566, &p567}
	var p577 = sequenceParser{id: 577, commit: 72, name: "continue", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{810}}
	var p569 = charParser{id: 569, chars: []rune{99}}
	var p570 = charParser{id: 570, chars: []rune{111}}
	var p571 = charParser{id: 571, chars: []rune{110}}
	var p572 = charParser{id: 572, chars: []rune{116}}
	var p573 = charParser{id: 573, chars: []rune{105}}
	var p574 = charParser{id: 574, chars: []rune{110}}
	var p575 = charParser{id: 575, chars: []rune{117}}
	var p576 = charParser{id: 576, chars: []rune{101}}
	p577.items = []parser{&p569, &p570, &p571, &p572, &p573, &p574, &p575, &p576}
	var p598 = sequenceParser{id: 598, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{810}}
	var p592 = sequenceParser{id: 592, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p589 = charParser{id: 589, chars: []rune{102}}
	var p590 = charParser{id: 590, chars: []rune{111}}
	var p591 = charParser{id: 591, chars: []rune{114}}
	p592.items = []parser{&p589, &p590, &p591}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p593.items = []parser{&p842, &p14}
	p594.items = []parser{&p14, &p593}
	var p588 = choiceParser{id: 588, commit: 66, name: "loop-expression"}
	var p587 = choiceParser{id: 587, commit: 64, name: "range-over-expression", generalizations: []int{588}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{587, 588}}
	var p583 = sequenceParser{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p582 = sequenceParser{id: 582, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p582.items = []parser{&p842, &p14}
	p583.items = []parser{&p842, &p14, &p582}
	var p580 = sequenceParser{id: 580, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p578 = charParser{id: 578, chars: []rune{105}}
	var p579 = charParser{id: 579, chars: []rune{110}}
	p580.items = []parser{&p578, &p579}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p584.items = []parser{&p842, &p14}
	p585.items = []parser{&p842, &p14, &p584}
	var p581 = choiceParser{id: 581, commit: 2}
	p581.options = []parser{&p388, &p239}
	p586.items = []parser{&p105, &p583, &p842, &p580, &p585, &p842, &p581}
	p587.options = []parser{&p586, &p239}
	p588.options = []parser{&p388, &p587}
	p595.items = []parser{&p594, &p842, &p588}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p596.items = []parser{&p842, &p14}
	p597.items = []parser{&p842, &p14, &p596}
	p598.items = []parser{&p592, &p842, &p595, &p597, &p842, &p205}
	var p746 = choiceParser{id: 746, commit: 66, name: "definition", generalizations: []int{810}}
	var p659 = sequenceParser{id: 659, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{746, 810}}
	var p655 = sequenceParser{id: 655, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p652 = charParser{id: 652, chars: []rune{108}}
	var p653 = charParser{id: 653, chars: []rune{101}}
	var p654 = charParser{id: 654, chars: []rune{116}}
	p655.items = []parser{&p652, &p653, &p654}
	var p658 = sequenceParser{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p657 = sequenceParser{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p657.items = []parser{&p842, &p14}
	p658.items = []parser{&p842, &p14, &p657}
	var p656 = choiceParser{id: 656, commit: 2}
	var p646 = sequenceParser{id: 646, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{656, 660, 661}}
	var p645 = sequenceParser{id: 645, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p641 = sequenceParser{id: 641, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p640 = sequenceParser{id: 640, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p640.items = []parser{&p842, &p14}
	p641.items = []parser{&p14, &p640}
	var p639 = sequenceParser{id: 639, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p638 = charParser{id: 638, chars: []rune{61}}
	p639.items = []parser{&p638}
	p642.items = []parser{&p641, &p842, &p639}
	var p644 = sequenceParser{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p643 = sequenceParser{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p643.items = []parser{&p842, &p14}
	p644.items = []parser{&p842, &p14, &p643}
	p645.items = []parser{&p105, &p842, &p642, &p644, &p842, &p388}
	p646.items = []parser{&p645}
	var p651 = sequenceParser{id: 651, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{656, 660, 661}}
	var p648 = sequenceParser{id: 648, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p647 = charParser{id: 647, chars: []rune{126}}
	p648.items = []parser{&p647}
	var p650 = sequenceParser{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p649 = sequenceParser{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p649.items = []parser{&p842, &p14}
	p650.items = []parser{&p842, &p14, &p649}
	p651.items = []parser{&p648, &p650, &p842, &p645}
	p656.options = []parser{&p646, &p651}
	p659.items = []parser{&p655, &p658, &p842, &p656}
	var p680 = sequenceParser{id: 680, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{746, 810}}
	var p673 = sequenceParser{id: 673, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p670 = charParser{id: 670, chars: []rune{108}}
	var p671 = charParser{id: 671, chars: []rune{101}}
	var p672 = charParser{id: 672, chars: []rune{116}}
	p673.items = []parser{&p670, &p671, &p672}
	var p679 = sequenceParser{id: 679, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p678 = sequenceParser{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p678.items = []parser{&p842, &p14}
	p679.items = []parser{&p842, &p14, &p678}
	var p675 = sequenceParser{id: 675, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p674 = charParser{id: 674, chars: []rune{40}}
	p675.items = []parser{&p674}
	var p665 = sequenceParser{id: 665, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p660 = choiceParser{id: 660, commit: 2}
	p660.options = []parser{&p646, &p651}
	var p664 = sequenceParser{id: 664, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p662 = sequenceParser{id: 662, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p661 = choiceParser{id: 661, commit: 2}
	p661.options = []parser{&p646, &p651}
	p662.items = []parser{&p115, &p842, &p661}
	var p663 = sequenceParser{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p663.items = []parser{&p842, &p662}
	p664.items = []parser{&p842, &p662, &p663}
	p665.items = []parser{&p660, &p664}
	var p677 = sequenceParser{id: 677, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p676 = charParser{id: 676, chars: []rune{41}}
	p677.items = []parser{&p676}
	p680.items = []parser{&p673, &p679, &p842, &p675, &p842, &p115, &p842, &p665, &p842, &p115, &p842, &p677}
	var p695 = sequenceParser{id: 695, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{746, 810}}
	var p684 = sequenceParser{id: 684, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p681 = charParser{id: 681, chars: []rune{108}}
	var p682 = charParser{id: 682, chars: []rune{101}}
	var p683 = charParser{id: 683, chars: []rune{116}}
	p684.items = []parser{&p681, &p682, &p683}
	var p692 = sequenceParser{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p691 = sequenceParser{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p691.items = []parser{&p842, &p14}
	p692.items = []parser{&p842, &p14, &p691}
	var p686 = sequenceParser{id: 686, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p685 = charParser{id: 685, chars: []rune{126}}
	p686.items = []parser{&p685}
	var p694 = sequenceParser{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p693 = sequenceParser{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p693.items = []parser{&p842, &p14}
	p694.items = []parser{&p842, &p14, &p693}
	var p688 = sequenceParser{id: 688, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p687 = charParser{id: 687, chars: []rune{40}}
	p688.items = []parser{&p687}
	var p669 = sequenceParser{id: 669, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p668 = sequenceParser{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p666.items = []parser{&p115, &p842, &p646}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p667.items = []parser{&p842, &p666}
	p668.items = []parser{&p842, &p666, &p667}
	p669.items = []parser{&p646, &p668}
	var p690 = sequenceParser{id: 690, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p689 = charParser{id: 689, chars: []rune{41}}
	p690.items = []parser{&p689}
	p695.items = []parser{&p684, &p692, &p842, &p686, &p694, &p842, &p688, &p842, &p115, &p842, &p669, &p842, &p115, &p842, &p690}
	var p711 = sequenceParser{id: 711, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{746, 810}}
	var p707 = sequenceParser{id: 707, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p705 = charParser{id: 705, chars: []rune{102}}
	var p706 = charParser{id: 706, chars: []rune{110}}
	p707.items = []parser{&p705, &p706}
	var p710 = sequenceParser{id: 710, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p709 = sequenceParser{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p709.items = []parser{&p842, &p14}
	p710.items = []parser{&p842, &p14, &p709}
	var p708 = choiceParser{id: 708, commit: 2}
	var p699 = sequenceParser{id: 699, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{708, 716, 717}}
	var p698 = sequenceParser{id: 698, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p697 = sequenceParser{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p696.items = []parser{&p842, &p14}
	p697.items = []parser{&p842, &p14, &p696}
	p698.items = []parser{&p105, &p697, &p842, &p214}
	p699.items = []parser{&p698}
	var p704 = sequenceParser{id: 704, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{708, 716, 717}}
	var p701 = sequenceParser{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p700 = charParser{id: 700, chars: []rune{126}}
	p701.items = []parser{&p700}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p702 = sequenceParser{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p702.items = []parser{&p842, &p14}
	p703.items = []parser{&p842, &p14, &p702}
	p704.items = []parser{&p701, &p703, &p842, &p698}
	p708.options = []parser{&p699, &p704}
	p711.items = []parser{&p707, &p710, &p842, &p708}
	var p731 = sequenceParser{id: 731, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{746, 810}}
	var p724 = sequenceParser{id: 724, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p722 = charParser{id: 722, chars: []rune{102}}
	var p723 = charParser{id: 723, chars: []rune{110}}
	p724.items = []parser{&p722, &p723}
	var p730 = sequenceParser{id: 730, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p729 = sequenceParser{id: 729, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p729.items = []parser{&p842, &p14}
	p730.items = []parser{&p842, &p14, &p729}
	var p726 = sequenceParser{id: 726, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p725 = charParser{id: 725, chars: []rune{40}}
	p726.items = []parser{&p725}
	var p721 = sequenceParser{id: 721, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p716 = choiceParser{id: 716, commit: 2}
	p716.options = []parser{&p699, &p704}
	var p720 = sequenceParser{id: 720, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p718 = sequenceParser{id: 718, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p717 = choiceParser{id: 717, commit: 2}
	p717.options = []parser{&p699, &p704}
	p718.items = []parser{&p115, &p842, &p717}
	var p719 = sequenceParser{id: 719, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p719.items = []parser{&p842, &p718}
	p720.items = []parser{&p842, &p718, &p719}
	p721.items = []parser{&p716, &p720}
	var p728 = sequenceParser{id: 728, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p727 = charParser{id: 727, chars: []rune{41}}
	p728.items = []parser{&p727}
	p731.items = []parser{&p724, &p730, &p842, &p726, &p842, &p115, &p842, &p721, &p842, &p115, &p842, &p728}
	var p745 = sequenceParser{id: 745, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{746, 810}}
	var p734 = sequenceParser{id: 734, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p732 = charParser{id: 732, chars: []rune{102}}
	var p733 = charParser{id: 733, chars: []rune{110}}
	p734.items = []parser{&p732, &p733}
	var p742 = sequenceParser{id: 742, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p741 = sequenceParser{id: 741, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p741.items = []parser{&p842, &p14}
	p742.items = []parser{&p842, &p14, &p741}
	var p736 = sequenceParser{id: 736, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p735 = charParser{id: 735, chars: []rune{126}}
	p736.items = []parser{&p735}
	var p744 = sequenceParser{id: 744, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p743 = sequenceParser{id: 743, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p743.items = []parser{&p842, &p14}
	p744.items = []parser{&p842, &p14, &p743}
	var p738 = sequenceParser{id: 738, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p737 = charParser{id: 737, chars: []rune{40}}
	p738.items = []parser{&p737}
	var p715 = sequenceParser{id: 715, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p714 = sequenceParser{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p712 = sequenceParser{id: 712, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p712.items = []parser{&p115, &p842, &p699}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p713.items = []parser{&p842, &p712}
	p714.items = []parser{&p842, &p712, &p713}
	p715.items = []parser{&p699, &p714}
	var p740 = sequenceParser{id: 740, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p739 = charParser{id: 739, chars: []rune{41}}
	p740.items = []parser{&p739}
	p745.items = []parser{&p734, &p742, &p842, &p736, &p744, &p842, &p738, &p842, &p115, &p842, &p715, &p842, &p115, &p842, &p740}
	p746.options = []parser{&p659, &p680, &p695, &p711, &p731, &p745}
	var p789 = choiceParser{id: 789, commit: 64, name: "require", generalizations: []int{810}}
	var p773 = sequenceParser{id: 773, commit: 66, name: "require-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{789, 810}}
	var p770 = sequenceParser{id: 770, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p763 = charParser{id: 763, chars: []rune{114}}
	var p764 = charParser{id: 764, chars: []rune{101}}
	var p765 = charParser{id: 765, chars: []rune{113}}
	var p766 = charParser{id: 766, chars: []rune{117}}
	var p767 = charParser{id: 767, chars: []rune{105}}
	var p768 = charParser{id: 768, chars: []rune{114}}
	var p769 = charParser{id: 769, chars: []rune{101}}
	p770.items = []parser{&p763, &p764, &p765, &p766, &p767, &p768, &p769}
	var p772 = sequenceParser{id: 772, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p771 = sequenceParser{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p771.items = []parser{&p842, &p14}
	p772.items = []parser{&p842, &p14, &p771}
	var p758 = choiceParser{id: 758, commit: 64, name: "require-fact"}
	var p757 = sequenceParser{id: 757, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{758}}
	var p749 = choiceParser{id: 749, commit: 2}
	var p748 = sequenceParser{id: 748, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{749}}
	var p747 = charParser{id: 747, chars: []rune{46}}
	p748.items = []parser{&p747}
	p749.options = []parser{&p105, &p748}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p752.items = []parser{&p842, &p14}
	p753.items = []parser{&p14, &p752}
	var p751 = sequenceParser{id: 751, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p750 = charParser{id: 750, chars: []rune{61}}
	p751.items = []parser{&p750}
	p754.items = []parser{&p753, &p842, &p751}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p755.items = []parser{&p842, &p14}
	p756.items = []parser{&p842, &p14, &p755}
	p757.items = []parser{&p749, &p842, &p754, &p756, &p842, &p88}
	p758.options = []parser{&p88, &p757}
	p773.items = []parser{&p770, &p772, &p842, &p758}
	var p788 = sequenceParser{id: 788, commit: 66, name: "require-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{789, 810}}
	var p781 = sequenceParser{id: 781, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p774 = charParser{id: 774, chars: []rune{114}}
	var p775 = charParser{id: 775, chars: []rune{101}}
	var p776 = charParser{id: 776, chars: []rune{113}}
	var p777 = charParser{id: 777, chars: []rune{117}}
	var p778 = charParser{id: 778, chars: []rune{105}}
	var p779 = charParser{id: 779, chars: []rune{114}}
	var p780 = charParser{id: 780, chars: []rune{101}}
	p781.items = []parser{&p774, &p775, &p776, &p777, &p778, &p779, &p780}
	var p787 = sequenceParser{id: 787, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p786 = sequenceParser{id: 786, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p786.items = []parser{&p842, &p14}
	p787.items = []parser{&p842, &p14, &p786}
	var p783 = sequenceParser{id: 783, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p782 = charParser{id: 782, chars: []rune{40}}
	p783.items = []parser{&p782}
	var p762 = sequenceParser{id: 762, commit: 66, name: "require-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p761 = sequenceParser{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p759 = sequenceParser{id: 759, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p759.items = []parser{&p115, &p842, &p758}
	var p760 = sequenceParser{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p760.items = []parser{&p842, &p759}
	p761.items = []parser{&p842, &p759, &p760}
	p762.items = []parser{&p758, &p761}
	var p785 = sequenceParser{id: 785, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p784 = charParser{id: 784, chars: []rune{41}}
	p785.items = []parser{&p784}
	p788.items = []parser{&p781, &p787, &p842, &p783, &p842, &p115, &p842, &p762, &p842, &p115, &p842, &p785}
	p789.options = []parser{&p773, &p788}
	var p799 = sequenceParser{id: 799, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{810}}
	var p796 = sequenceParser{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p790 = charParser{id: 790, chars: []rune{101}}
	var p791 = charParser{id: 791, chars: []rune{120}}
	var p792 = charParser{id: 792, chars: []rune{112}}
	var p793 = charParser{id: 793, chars: []rune{111}}
	var p794 = charParser{id: 794, chars: []rune{114}}
	var p795 = charParser{id: 795, chars: []rune{116}}
	p796.items = []parser{&p790, &p791, &p792, &p793, &p794, &p795}
	var p798 = sequenceParser{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p797 = sequenceParser{id: 797, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p797.items = []parser{&p842, &p14}
	p798.items = []parser{&p842, &p14, &p797}
	p799.items = []parser{&p796, &p798, &p842, &p746}
	var p819 = sequenceParser{id: 819, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{810}}
	var p812 = sequenceParser{id: 812, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p811 = charParser{id: 811, chars: []rune{40}}
	p812.items = []parser{&p811}
	var p816 = sequenceParser{id: 816, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p815 = sequenceParser{id: 815, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p815.items = []parser{&p842, &p14}
	p816.items = []parser{&p842, &p14, &p815}
	var p818 = sequenceParser{id: 818, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p817 = sequenceParser{id: 817, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p817.items = []parser{&p842, &p14}
	p818.items = []parser{&p842, &p14, &p817}
	var p814 = sequenceParser{id: 814, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p813 = charParser{id: 813, chars: []rune{41}}
	p814.items = []parser{&p813}
	p819.items = []parser{&p812, &p816, &p842, &p810, &p818, &p842, &p814}
	p810.options = []parser{&p800, &p200, &p425, &p482, &p547, &p568, &p577, &p598, &p746, &p789, &p799, &p819}
	var p827 = sequenceParser{id: 827, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p825 = sequenceParser{id: 825, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p825.items = []parser{&p824, &p842, &p810}
	var p826 = sequenceParser{id: 826, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p826.items = []parser{&p842, &p825}
	p827.items = []parser{&p842, &p825, &p826}
	p828.items = []parser{&p810, &p827}
	p843.items = []parser{&p839, &p842, &p824, &p842, &p828, &p842, &p824}
	p844.items = []parser{&p842, &p843, &p842}
	var b844 = sequenceBuilder{id: 844, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b842 = choiceBuilder{id: 842, commit: 2}
	var b840 = choiceBuilder{id: 840, commit: 70}
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
	b840.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b841 = sequenceBuilder{id: 841, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
	var b44 = sequenceBuilder{id: 44, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
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
	var b43 = sequenceBuilder{id: 43, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b41 = sequenceBuilder{id: 41, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b40 = sequenceBuilder{id: 40, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b39 = charBuilder{}
	b40.items = []builder{&b39}
	b41.items = []builder{&b40, &b842, &b38}
	var b42 = sequenceBuilder{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b42.items = []builder{&b842, &b41}
	b43.items = []builder{&b842, &b41, &b42}
	b44.items = []builder{&b38, &b43}
	b841.items = []builder{&b44}
	b842.options = []builder{&b840, &b841}
	var b843 = sequenceBuilder{id: 843, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b839 = sequenceBuilder{id: 839, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b836 = sequenceBuilder{id: 836, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b834 = charBuilder{}
	var b835 = charBuilder{}
	b836.items = []builder{&b834, &b835}
	var b833 = sequenceBuilder{id: 833, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b832 = sequenceBuilder{id: 832, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b830 = sequenceBuilder{id: 830, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b829 = charBuilder{}
	b830.items = []builder{&b829}
	var b831 = sequenceBuilder{id: 831, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b831.items = []builder{&b842, &b830}
	b832.items = []builder{&b830, &b831}
	b833.items = []builder{&b832}
	var b838 = sequenceBuilder{id: 838, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b837 = charBuilder{}
	b838.items = []builder{&b837}
	b839.items = []builder{&b836, &b842, &b833, &b842, &b838}
	var b824 = sequenceBuilder{id: 824, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b822 = choiceBuilder{id: 822, commit: 2}
	var b821 = sequenceBuilder{id: 821, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b820 = charBuilder{}
	b821.items = []builder{&b820}
	var b14 = sequenceBuilder{id: 14, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b822.options = []builder{&b821, &b14}
	var b823 = sequenceBuilder{id: 823, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b823.items = []builder{&b842, &b822}
	b824.items = []builder{&b822, &b823}
	var b828 = sequenceBuilder{id: 828, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b810 = choiceBuilder{id: 810, commit: 66}
	var b800 = choiceBuilder{id: 800, commit: 66}
	var b388 = choiceBuilder{id: 388, commit: 66}
	var b280 = choiceBuilder{id: 280, commit: 66}
	var b62 = choiceBuilder{id: 62, commit: 64, name: "int"}
	var b53 = sequenceBuilder{id: 53, commit: 74, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b52 = sequenceBuilder{id: 52, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b51 = charBuilder{}
	b52.items = []builder{&b51}
	var b46 = sequenceBuilder{id: 46, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b45 = charBuilder{}
	b46.items = []builder{&b45}
	b53.items = []builder{&b52, &b46}
	var b56 = sequenceBuilder{id: 56, commit: 74, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b55 = sequenceBuilder{id: 55, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b54 = charBuilder{}
	b55.items = []builder{&b54}
	var b48 = sequenceBuilder{id: 48, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b47 = charBuilder{}
	b48.items = []builder{&b47}
	b56.items = []builder{&b55, &b48}
	var b61 = sequenceBuilder{id: 61, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}}
	var b58 = sequenceBuilder{id: 58, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b57 = charBuilder{}
	b58.items = []builder{&b57}
	var b60 = sequenceBuilder{id: 60, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b59 = charBuilder{}
	b60.items = []builder{&b59}
	var b50 = sequenceBuilder{id: 50, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b49 = charBuilder{}
	b50.items = []builder{&b49}
	b61.items = []builder{&b58, &b60, &b50}
	b62.options = []builder{&b53, &b56, &b61}
	var b75 = choiceBuilder{id: 75, commit: 72, name: "float"}
	var b70 = sequenceBuilder{id: 70, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}}
	var b69 = sequenceBuilder{id: 69, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b68 = charBuilder{}
	b69.items = []builder{&b68}
	var b67 = sequenceBuilder{id: 67, commit: 74, ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var b64 = sequenceBuilder{id: 64, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b63 = charBuilder{}
	b64.items = []builder{&b63}
	var b66 = sequenceBuilder{id: 66, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b65 = charBuilder{}
	b66.items = []builder{&b65}
	b67.items = []builder{&b64, &b66, &b46}
	b70.items = []builder{&b46, &b69, &b46, &b67}
	var b73 = sequenceBuilder{id: 73, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}}
	var b72 = sequenceBuilder{id: 72, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b71 = charBuilder{}
	b72.items = []builder{&b71}
	b73.items = []builder{&b72, &b46, &b67}
	var b74 = sequenceBuilder{id: 74, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}}
	b74.items = []builder{&b46, &b67}
	b75.options = []builder{&b70, &b73, &b74}
	var b88 = sequenceBuilder{id: 88, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}}
	var b77 = sequenceBuilder{id: 77, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b76 = charBuilder{}
	b77.items = []builder{&b76}
	var b85 = choiceBuilder{id: 85, commit: 10}
	var b79 = sequenceBuilder{id: 79, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b78 = charBuilder{}
	b79.items = []builder{&b78}
	var b84 = sequenceBuilder{id: 84, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b81 = sequenceBuilder{id: 81, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b80 = charBuilder{}
	b81.items = []builder{&b80}
	var b83 = sequenceBuilder{id: 83, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b82 = charBuilder{}
	b83.items = []builder{&b82}
	b84.items = []builder{&b81, &b83}
	b85.options = []builder{&b79, &b84}
	var b87 = sequenceBuilder{id: 87, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b86 = charBuilder{}
	b87.items = []builder{&b86}
	b88.items = []builder{&b77, &b85, &b87}
	var b100 = choiceBuilder{id: 100, commit: 66}
	var b93 = sequenceBuilder{id: 93, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b89 = charBuilder{}
	var b90 = charBuilder{}
	var b91 = charBuilder{}
	var b92 = charBuilder{}
	b93.items = []builder{&b89, &b90, &b91, &b92}
	var b99 = sequenceBuilder{id: 99, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b94 = charBuilder{}
	var b95 = charBuilder{}
	var b96 = charBuilder{}
	var b97 = charBuilder{}
	var b98 = charBuilder{}
	b99.items = []builder{&b94, &b95, &b96, &b97, &b98}
	b100.options = []builder{&b93, &b99}
	var b105 = sequenceBuilder{id: 105, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b102 = sequenceBuilder{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b101 = charBuilder{}
	b102.items = []builder{&b101}
	var b104 = sequenceBuilder{id: 104, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b103 = charBuilder{}
	b104.items = []builder{&b103}
	b105.items = []builder{&b102, &b104}
	var b126 = sequenceBuilder{id: 126, commit: 64, name: "list", ranges: [][]int{{1, 1}}}
	var b125 = sequenceBuilder{id: 125, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b122 = sequenceBuilder{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b121 = charBuilder{}
	b122.items = []builder{&b121}
	var b115 = sequenceBuilder{id: 115, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b113 = choiceBuilder{id: 113, commit: 2}
	var b112 = sequenceBuilder{id: 112, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b111 = charBuilder{}
	b112.items = []builder{&b111}
	b113.options = []builder{&b14, &b112}
	var b114 = sequenceBuilder{id: 114, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b114.items = []builder{&b842, &b113}
	b115.items = []builder{&b113, &b114}
	var b120 = sequenceBuilder{id: 120, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b116 = choiceBuilder{id: 116, commit: 66}
	var b110 = sequenceBuilder{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b109 = sequenceBuilder{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	var b108 = charBuilder{}
	b109.items = []builder{&b106, &b107, &b108}
	b110.items = []builder{&b280, &b842, &b109}
	b116.options = []builder{&b388, &b110}
	var b119 = sequenceBuilder{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b117.items = []builder{&b115, &b842, &b116}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b118.items = []builder{&b842, &b117}
	b119.items = []builder{&b842, &b117, &b118}
	b120.items = []builder{&b116, &b119}
	var b124 = sequenceBuilder{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b123 = charBuilder{}
	b124.items = []builder{&b123}
	b125.items = []builder{&b122, &b842, &b115, &b842, &b120, &b842, &b115, &b842, &b124}
	b126.items = []builder{&b125}
	var b131 = sequenceBuilder{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b128 = sequenceBuilder{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b127 = charBuilder{}
	b128.items = []builder{&b127}
	var b130 = sequenceBuilder{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b129.items = []builder{&b842, &b14}
	b130.items = []builder{&b842, &b14, &b129}
	b131.items = []builder{&b128, &b130, &b842, &b125}
	var b160 = sequenceBuilder{id: 160, commit: 64, name: "struct", ranges: [][]int{{1, 1}}}
	var b159 = sequenceBuilder{id: 159, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b156 = sequenceBuilder{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b155 = charBuilder{}
	b156.items = []builder{&b155}
	var b154 = sequenceBuilder{id: 154, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b149 = choiceBuilder{id: 149, commit: 2}
	var b148 = sequenceBuilder{id: 148, commit: 64, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b141 = choiceBuilder{id: 141, commit: 2}
	var b140 = sequenceBuilder{id: 140, commit: 64, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b133 = sequenceBuilder{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b132 = charBuilder{}
	b133.items = []builder{&b132}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b136 = sequenceBuilder{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b136.items = []builder{&b842, &b14}
	b137.items = []builder{&b842, &b14, &b136}
	var b139 = sequenceBuilder{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b138.items = []builder{&b842, &b14}
	b139.items = []builder{&b842, &b14, &b138}
	var b135 = sequenceBuilder{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b134 = charBuilder{}
	b135.items = []builder{&b134}
	b140.items = []builder{&b133, &b137, &b842, &b388, &b139, &b842, &b135}
	b141.options = []builder{&b105, &b88, &b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b842, &b14}
	b145.items = []builder{&b842, &b14, &b144}
	var b143 = sequenceBuilder{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b142 = charBuilder{}
	b143.items = []builder{&b142}
	var b147 = sequenceBuilder{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b146.items = []builder{&b842, &b14}
	b147.items = []builder{&b842, &b14, &b146}
	b148.items = []builder{&b141, &b145, &b842, &b143, &b147, &b842, &b388}
	b149.options = []builder{&b148, &b110}
	var b153 = sequenceBuilder{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b150 = choiceBuilder{id: 150, commit: 2}
	b150.options = []builder{&b148, &b110}
	b151.items = []builder{&b115, &b842, &b150}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b152.items = []builder{&b842, &b151}
	b153.items = []builder{&b842, &b151, &b152}
	b154.items = []builder{&b149, &b153}
	var b158 = sequenceBuilder{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b157 = charBuilder{}
	b158.items = []builder{&b157}
	b159.items = []builder{&b156, &b842, &b115, &b842, &b154, &b842, &b115, &b842, &b158}
	b160.items = []builder{&b159}
	var b165 = sequenceBuilder{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b162 = sequenceBuilder{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b161 = charBuilder{}
	b162.items = []builder{&b161}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b163.items = []builder{&b842, &b14}
	b164.items = []builder{&b842, &b14, &b163}
	b165.items = []builder{&b162, &b164, &b842, &b159}
	var b178 = choiceBuilder{id: 178, commit: 64, name: "channel"}
	var b168 = sequenceBuilder{id: 168, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b166 = charBuilder{}
	var b167 = charBuilder{}
	b168.items = []builder{&b166, &b167}
	var b177 = sequenceBuilder{id: 177, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b170 = sequenceBuilder{id: 170, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b169 = charBuilder{}
	b170.items = []builder{&b169}
	var b174 = sequenceBuilder{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b173 = sequenceBuilder{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b173.items = []builder{&b842, &b14}
	b174.items = []builder{&b842, &b14, &b173}
	var b176 = sequenceBuilder{id: 176, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b175 = sequenceBuilder{id: 175, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b175.items = []builder{&b842, &b14}
	b176.items = []builder{&b842, &b14, &b175}
	var b172 = sequenceBuilder{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b171 = charBuilder{}
	b172.items = []builder{&b171}
	b177.items = []builder{&b170, &b174, &b842, &b388, &b176, &b842, &b172}
	b178.options = []builder{&b168, &b177}
	var b220 = sequenceBuilder{id: 220, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b217 = sequenceBuilder{id: 217, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b215 = charBuilder{}
	var b216 = charBuilder{}
	b217.items = []builder{&b215, &b216}
	var b219 = sequenceBuilder{id: 219, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b218 = sequenceBuilder{id: 218, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b218.items = []builder{&b842, &b14}
	b219.items = []builder{&b842, &b14, &b218}
	var b214 = sequenceBuilder{id: 214, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b207 = sequenceBuilder{id: 207, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b206 = charBuilder{}
	b207.items = []builder{&b206}
	var b182 = sequenceBuilder{id: 182, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b181 = sequenceBuilder{id: 181, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b179 = sequenceBuilder{id: 179, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b179.items = []builder{&b115, &b842, &b105}
	var b180 = sequenceBuilder{id: 180, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b180.items = []builder{&b842, &b179}
	b181.items = []builder{&b842, &b179, &b180}
	b182.items = []builder{&b105, &b181}
	var b208 = sequenceBuilder{id: 208, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b189 = sequenceBuilder{id: 189, commit: 64, name: "collect-argument", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b186 = sequenceBuilder{id: 186, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b183 = charBuilder{}
	var b184 = charBuilder{}
	var b185 = charBuilder{}
	b186.items = []builder{&b183, &b184, &b185}
	var b188 = sequenceBuilder{id: 188, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b187 = sequenceBuilder{id: 187, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b187.items = []builder{&b842, &b14}
	b188.items = []builder{&b842, &b14, &b187}
	b189.items = []builder{&b186, &b188, &b842, &b105}
	b208.items = []builder{&b115, &b842, &b189}
	var b210 = sequenceBuilder{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b209 = charBuilder{}
	b210.items = []builder{&b209}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b212 = sequenceBuilder{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b212.items = []builder{&b842, &b14}
	b213.items = []builder{&b842, &b14, &b212}
	var b211 = choiceBuilder{id: 211, commit: 2}
	var b205 = sequenceBuilder{id: 205, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b202 = sequenceBuilder{id: 202, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b201 = charBuilder{}
	b202.items = []builder{&b201}
	var b204 = sequenceBuilder{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b203 = charBuilder{}
	b204.items = []builder{&b203}
	b205.items = []builder{&b202, &b842, &b824, &b842, &b828, &b842, &b824, &b842, &b204}
	b211.options = []builder{&b800, &b205}
	b214.items = []builder{&b207, &b842, &b115, &b842, &b182, &b842, &b208, &b842, &b115, &b842, &b210, &b213, &b842, &b211}
	b220.items = []builder{&b217, &b219, &b842, &b214}
	var b230 = sequenceBuilder{id: 230, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b223 = sequenceBuilder{id: 223, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b221 = charBuilder{}
	var b222 = charBuilder{}
	b223.items = []builder{&b221, &b222}
	var b227 = sequenceBuilder{id: 227, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b226 = sequenceBuilder{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b226.items = []builder{&b842, &b14}
	b227.items = []builder{&b842, &b14, &b226}
	var b225 = sequenceBuilder{id: 225, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b224 = charBuilder{}
	b225.items = []builder{&b224}
	var b229 = sequenceBuilder{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b228 = sequenceBuilder{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b228.items = []builder{&b842, &b14}
	b229.items = []builder{&b842, &b14, &b228}
	b230.items = []builder{&b223, &b227, &b842, &b225, &b229, &b842, &b214}
	var b258 = choiceBuilder{id: 258, commit: 64, name: "expression-indexer"}
	var b248 = sequenceBuilder{id: 248, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b241 = sequenceBuilder{id: 241, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b240 = charBuilder{}
	b241.items = []builder{&b240}
	var b245 = sequenceBuilder{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b244 = sequenceBuilder{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b244.items = []builder{&b842, &b14}
	b245.items = []builder{&b842, &b14, &b244}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b246 = sequenceBuilder{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b246.items = []builder{&b842, &b14}
	b247.items = []builder{&b842, &b14, &b246}
	var b243 = sequenceBuilder{id: 243, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b242 = charBuilder{}
	b243.items = []builder{&b242}
	b248.items = []builder{&b280, &b842, &b241, &b245, &b842, &b388, &b247, &b842, &b243}
	var b257 = sequenceBuilder{id: 257, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b250 = sequenceBuilder{id: 250, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b249 = charBuilder{}
	b250.items = []builder{&b249}
	var b254 = sequenceBuilder{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b253 = sequenceBuilder{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b253.items = []builder{&b842, &b14}
	b254.items = []builder{&b842, &b14, &b253}
	var b239 = sequenceBuilder{id: 239, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b231 = sequenceBuilder{id: 231, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b231.items = []builder{&b388}
	var b236 = sequenceBuilder{id: 236, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b235 = sequenceBuilder{id: 235, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b235.items = []builder{&b842, &b14}
	b236.items = []builder{&b842, &b14, &b235}
	var b234 = sequenceBuilder{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b233 = charBuilder{}
	b234.items = []builder{&b233}
	var b238 = sequenceBuilder{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b237 = sequenceBuilder{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b237.items = []builder{&b842, &b14}
	b238.items = []builder{&b842, &b14, &b237}
	var b232 = sequenceBuilder{id: 232, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b232.items = []builder{&b388}
	b239.items = []builder{&b231, &b236, &b842, &b234, &b238, &b842, &b232}
	var b256 = sequenceBuilder{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b255 = sequenceBuilder{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b255.items = []builder{&b842, &b14}
	b256.items = []builder{&b842, &b14, &b255}
	var b252 = sequenceBuilder{id: 252, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b251 = charBuilder{}
	b252.items = []builder{&b251}
	b257.items = []builder{&b280, &b842, &b250, &b254, &b842, &b239, &b256, &b842, &b252}
	b258.options = []builder{&b248, &b257}
	var b265 = sequenceBuilder{id: 265, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b261 = sequenceBuilder{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b261.items = []builder{&b842, &b14}
	b262.items = []builder{&b842, &b14, &b261}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b263.items = []builder{&b842, &b14}
	b264.items = []builder{&b842, &b14, &b263}
	b265.items = []builder{&b280, &b262, &b842, &b260, &b264, &b842, &b105}
	var b270 = sequenceBuilder{id: 270, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b267 = sequenceBuilder{id: 267, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b266 = charBuilder{}
	b267.items = []builder{&b266}
	var b269 = sequenceBuilder{id: 269, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b268 = charBuilder{}
	b269.items = []builder{&b268}
	b270.items = []builder{&b280, &b842, &b267, &b842, &b115, &b842, &b120, &b842, &b115, &b842, &b269}
	var b486 = sequenceBuilder{id: 486, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b485 = sequenceBuilder{id: 485, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b483 = charBuilder{}
	var b484 = charBuilder{}
	b485.items = []builder{&b483, &b484}
	b486.items = []builder{&b485, &b842, &b280}
	var b279 = sequenceBuilder{id: 279, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b272 = sequenceBuilder{id: 272, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b271 = charBuilder{}
	b272.items = []builder{&b271}
	var b276 = sequenceBuilder{id: 276, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b275 = sequenceBuilder{id: 275, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b275.items = []builder{&b842, &b14}
	b276.items = []builder{&b842, &b14, &b275}
	var b278 = sequenceBuilder{id: 278, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b277 = sequenceBuilder{id: 277, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b277.items = []builder{&b842, &b14}
	b278.items = []builder{&b842, &b14, &b277}
	var b274 = sequenceBuilder{id: 274, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b273 = charBuilder{}
	b274.items = []builder{&b273}
	b279.items = []builder{&b272, &b276, &b842, &b388, &b278, &b842, &b274}
	b280.options = []builder{&b62, &b75, &b88, &b100, &b105, &b126, &b131, &b160, &b165, &b178, &b220, &b230, &b258, &b265, &b270, &b486, &b279}
	var b340 = sequenceBuilder{id: 340, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b339 = choiceBuilder{id: 339, commit: 66}
	var b299 = sequenceBuilder{id: 299, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b298 = charBuilder{}
	b299.items = []builder{&b298}
	var b301 = sequenceBuilder{id: 301, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b300 = charBuilder{}
	b301.items = []builder{&b300}
	var b282 = sequenceBuilder{id: 282, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b281 = charBuilder{}
	b282.items = []builder{&b281}
	b339.options = []builder{&b299, &b301, &b282}
	b340.items = []builder{&b339, &b842, &b280}
	var b374 = choiceBuilder{id: 374, commit: 66}
	var b354 = sequenceBuilder{id: 354, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b346 = choiceBuilder{id: 346, commit: 66}
	b346.options = []builder{&b280, &b340}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b341 = choiceBuilder{id: 341, commit: 66}
	var b284 = sequenceBuilder{id: 284, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b283 = charBuilder{}
	b284.items = []builder{&b283}
	var b291 = sequenceBuilder{id: 291, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b289 = charBuilder{}
	var b290 = charBuilder{}
	b291.items = []builder{&b289, &b290}
	var b294 = sequenceBuilder{id: 294, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b292 = charBuilder{}
	var b293 = charBuilder{}
	b294.items = []builder{&b292, &b293}
	var b297 = sequenceBuilder{id: 297, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b295 = charBuilder{}
	var b296 = charBuilder{}
	b297.items = []builder{&b295, &b296}
	var b303 = sequenceBuilder{id: 303, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b302 = charBuilder{}
	b303.items = []builder{&b302}
	var b305 = sequenceBuilder{id: 305, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b304 = charBuilder{}
	b305.items = []builder{&b304}
	var b307 = sequenceBuilder{id: 307, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b306 = charBuilder{}
	b307.items = []builder{&b306}
	b341.options = []builder{&b284, &b291, &b294, &b297, &b303, &b305, &b307}
	b352.items = []builder{&b341, &b842, &b346}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b353.items = []builder{&b842, &b352}
	b354.items = []builder{&b346, &b842, &b352, &b353}
	var b357 = sequenceBuilder{id: 357, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b347 = choiceBuilder{id: 347, commit: 66}
	b347.options = []builder{&b346, &b354}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b342 = choiceBuilder{id: 342, commit: 66}
	var b286 = sequenceBuilder{id: 286, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b285 = charBuilder{}
	b286.items = []builder{&b285}
	var b288 = sequenceBuilder{id: 288, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b287 = charBuilder{}
	b288.items = []builder{&b287}
	var b309 = sequenceBuilder{id: 309, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b308 = charBuilder{}
	b309.items = []builder{&b308}
	var b311 = sequenceBuilder{id: 311, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b310 = charBuilder{}
	b311.items = []builder{&b310}
	b342.options = []builder{&b286, &b288, &b309, &b311}
	b355.items = []builder{&b342, &b842, &b347}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b356.items = []builder{&b842, &b355}
	b357.items = []builder{&b347, &b842, &b355, &b356}
	var b360 = sequenceBuilder{id: 360, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b348 = choiceBuilder{id: 348, commit: 66}
	b348.options = []builder{&b347, &b357}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b343 = choiceBuilder{id: 343, commit: 66}
	var b316 = sequenceBuilder{id: 316, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b314 = charBuilder{}
	var b315 = charBuilder{}
	b316.items = []builder{&b314, &b315}
	var b319 = sequenceBuilder{id: 319, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b317 = charBuilder{}
	var b318 = charBuilder{}
	b319.items = []builder{&b317, &b318}
	var b321 = sequenceBuilder{id: 321, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b320 = charBuilder{}
	b321.items = []builder{&b320}
	var b324 = sequenceBuilder{id: 324, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b322 = charBuilder{}
	var b323 = charBuilder{}
	b324.items = []builder{&b322, &b323}
	var b326 = sequenceBuilder{id: 326, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b325 = charBuilder{}
	b326.items = []builder{&b325}
	var b329 = sequenceBuilder{id: 329, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b327 = charBuilder{}
	var b328 = charBuilder{}
	b329.items = []builder{&b327, &b328}
	b343.options = []builder{&b316, &b319, &b321, &b324, &b326, &b329}
	b358.items = []builder{&b343, &b842, &b348}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b842, &b358}
	b360.items = []builder{&b348, &b842, &b358, &b359}
	var b363 = sequenceBuilder{id: 363, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b349 = choiceBuilder{id: 349, commit: 66}
	b349.options = []builder{&b348, &b360}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b344 = sequenceBuilder{id: 344, commit: 66, ranges: [][]int{{1, 1}}}
	var b332 = sequenceBuilder{id: 332, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b330 = charBuilder{}
	var b331 = charBuilder{}
	b332.items = []builder{&b330, &b331}
	b344.items = []builder{&b332}
	b361.items = []builder{&b344, &b842, &b349}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b362.items = []builder{&b842, &b361}
	b363.items = []builder{&b349, &b842, &b361, &b362}
	var b366 = sequenceBuilder{id: 366, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b350 = choiceBuilder{id: 350, commit: 66}
	b350.options = []builder{&b349, &b363}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b345 = sequenceBuilder{id: 345, commit: 66, ranges: [][]int{{1, 1}}}
	var b335 = sequenceBuilder{id: 335, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b333 = charBuilder{}
	var b334 = charBuilder{}
	b335.items = []builder{&b333, &b334}
	b345.items = []builder{&b335}
	b364.items = []builder{&b345, &b842, &b350}
	var b365 = sequenceBuilder{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b365.items = []builder{&b842, &b364}
	b366.items = []builder{&b350, &b842, &b364, &b365}
	var b373 = sequenceBuilder{id: 373, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b351 = choiceBuilder{id: 351, commit: 66}
	b351.options = []builder{&b350, &b366}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b367.items = []builder{&b842, &b14}
	b368.items = []builder{&b14, &b367}
	var b338 = sequenceBuilder{id: 338, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b336 = charBuilder{}
	var b337 = charBuilder{}
	b338.items = []builder{&b336, &b337}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b369.items = []builder{&b842, &b14}
	b370.items = []builder{&b842, &b14, &b369}
	b371.items = []builder{&b368, &b842, &b338, &b370, &b842, &b351}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b372.items = []builder{&b842, &b371}
	b373.items = []builder{&b351, &b842, &b371, &b372}
	b374.options = []builder{&b354, &b357, &b360, &b363, &b366, &b373}
	var b387 = sequenceBuilder{id: 387, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b380 = sequenceBuilder{id: 380, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b379 = sequenceBuilder{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b379.items = []builder{&b842, &b14}
	b380.items = []builder{&b842, &b14, &b379}
	var b376 = sequenceBuilder{id: 376, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b375 = charBuilder{}
	b376.items = []builder{&b375}
	var b382 = sequenceBuilder{id: 382, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b381 = sequenceBuilder{id: 381, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b381.items = []builder{&b842, &b14}
	b382.items = []builder{&b842, &b14, &b381}
	var b384 = sequenceBuilder{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b383 = sequenceBuilder{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b383.items = []builder{&b842, &b14}
	b384.items = []builder{&b842, &b14, &b383}
	var b378 = sequenceBuilder{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b377 = charBuilder{}
	b378.items = []builder{&b377}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b385.items = []builder{&b842, &b14}
	b386.items = []builder{&b842, &b14, &b385}
	b387.items = []builder{&b388, &b380, &b842, &b376, &b382, &b842, &b388, &b384, &b842, &b378, &b386, &b842, &b388}
	b388.options = []builder{&b280, &b340, &b374, &b387}
	var b509 = sequenceBuilder{id: 509, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b508 = sequenceBuilder{id: 508, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b506 = charBuilder{}
	var b507 = charBuilder{}
	b508.items = []builder{&b506, &b507}
	b509.items = []builder{&b280, &b842, &b508, &b842, &b280}
	var b553 = sequenceBuilder{id: 553, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b550 = sequenceBuilder{id: 550, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b548 = charBuilder{}
	var b549 = charBuilder{}
	b550.items = []builder{&b548, &b549}
	var b552 = sequenceBuilder{id: 552, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b551 = sequenceBuilder{id: 551, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b551.items = []builder{&b842, &b14}
	b552.items = []builder{&b842, &b14, &b551}
	b553.items = []builder{&b550, &b552, &b842, &b270}
	var b562 = sequenceBuilder{id: 562, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b559 = sequenceBuilder{id: 559, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b554 = charBuilder{}
	var b555 = charBuilder{}
	var b556 = charBuilder{}
	var b557 = charBuilder{}
	var b558 = charBuilder{}
	b559.items = []builder{&b554, &b555, &b556, &b557, &b558}
	var b561 = sequenceBuilder{id: 561, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b560 = sequenceBuilder{id: 560, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b560.items = []builder{&b842, &b14}
	b561.items = []builder{&b842, &b14, &b560}
	b562.items = []builder{&b559, &b561, &b842, &b270}
	var b637 = choiceBuilder{id: 637, commit: 64, name: "assignment"}
	var b618 = sequenceBuilder{id: 618, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b615 = sequenceBuilder{id: 615, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b612 = charBuilder{}
	var b613 = charBuilder{}
	var b614 = charBuilder{}
	b615.items = []builder{&b612, &b613, &b614}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b616.items = []builder{&b842, &b14}
	b617.items = []builder{&b842, &b14, &b616}
	var b607 = sequenceBuilder{id: 607, commit: 64, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b599 = choiceBuilder{id: 599, commit: 66}
	b599.options = []builder{&b105, &b248, &b265}
	var b604 = sequenceBuilder{id: 604, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b603 = sequenceBuilder{id: 603, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b602 = sequenceBuilder{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b602.items = []builder{&b842, &b14}
	b603.items = []builder{&b14, &b602}
	var b601 = sequenceBuilder{id: 601, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b600 = charBuilder{}
	b601.items = []builder{&b600}
	b604.items = []builder{&b603, &b842, &b601}
	var b606 = sequenceBuilder{id: 606, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b605 = sequenceBuilder{id: 605, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b605.items = []builder{&b842, &b14}
	b606.items = []builder{&b842, &b14, &b605}
	b607.items = []builder{&b599, &b842, &b604, &b606, &b842, &b388}
	b618.items = []builder{&b615, &b617, &b842, &b607}
	var b625 = sequenceBuilder{id: 625, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b622 = sequenceBuilder{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b621 = sequenceBuilder{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b621.items = []builder{&b842, &b14}
	b622.items = []builder{&b842, &b14, &b621}
	var b620 = sequenceBuilder{id: 620, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b619 = charBuilder{}
	b620.items = []builder{&b619}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b623 = sequenceBuilder{id: 623, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b623.items = []builder{&b842, &b14}
	b624.items = []builder{&b842, &b14, &b623}
	b625.items = []builder{&b599, &b622, &b842, &b620, &b624, &b842, &b388}
	var b636 = sequenceBuilder{id: 636, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b629 = sequenceBuilder{id: 629, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b626 = charBuilder{}
	var b627 = charBuilder{}
	var b628 = charBuilder{}
	b629.items = []builder{&b626, &b627, &b628}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b634.items = []builder{&b842, &b14}
	b635.items = []builder{&b842, &b14, &b634}
	var b631 = sequenceBuilder{id: 631, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b630 = charBuilder{}
	b631.items = []builder{&b630}
	var b611 = sequenceBuilder{id: 611, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b610 = sequenceBuilder{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b608 = sequenceBuilder{id: 608, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b608.items = []builder{&b115, &b842, &b607}
	var b609 = sequenceBuilder{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b609.items = []builder{&b842, &b608}
	b610.items = []builder{&b842, &b608, &b609}
	b611.items = []builder{&b607, &b610}
	var b633 = sequenceBuilder{id: 633, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b632 = charBuilder{}
	b633.items = []builder{&b632}
	b636.items = []builder{&b629, &b635, &b842, &b631, &b842, &b115, &b842, &b611, &b842, &b115, &b842, &b633}
	b637.options = []builder{&b618, &b625, &b636}
	var b809 = sequenceBuilder{id: 809, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b802 = sequenceBuilder{id: 802, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b801 = charBuilder{}
	b802.items = []builder{&b801}
	var b806 = sequenceBuilder{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b805 = sequenceBuilder{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b805.items = []builder{&b842, &b14}
	b806.items = []builder{&b842, &b14, &b805}
	var b808 = sequenceBuilder{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b807 = sequenceBuilder{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b807.items = []builder{&b842, &b14}
	b808.items = []builder{&b842, &b14, &b807}
	var b804 = sequenceBuilder{id: 804, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b803 = charBuilder{}
	b804.items = []builder{&b803}
	b809.items = []builder{&b802, &b806, &b842, &b800, &b808, &b842, &b804}
	b800.options = []builder{&b388, &b509, &b553, &b562, &b637, &b809}
	var b200 = sequenceBuilder{id: 200, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b196 = sequenceBuilder{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b190 = charBuilder{}
	var b191 = charBuilder{}
	var b192 = charBuilder{}
	var b193 = charBuilder{}
	var b194 = charBuilder{}
	var b195 = charBuilder{}
	b196.items = []builder{&b190, &b191, &b192, &b193, &b194, &b195}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b198 = sequenceBuilder{id: 198, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b197 = sequenceBuilder{id: 197, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b197.items = []builder{&b842, &b14}
	b198.items = []builder{&b14, &b197}
	b199.items = []builder{&b198, &b842, &b388}
	b200.items = []builder{&b196, &b842, &b199}
	var b425 = sequenceBuilder{id: 425, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b391 = sequenceBuilder{id: 391, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b389 = charBuilder{}
	var b390 = charBuilder{}
	b391.items = []builder{&b389, &b390}
	var b420 = sequenceBuilder{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b419 = sequenceBuilder{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b419.items = []builder{&b842, &b14}
	b420.items = []builder{&b842, &b14, &b419}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b421 = sequenceBuilder{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b421.items = []builder{&b842, &b14}
	b422.items = []builder{&b842, &b14, &b421}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b401 = sequenceBuilder{id: 401, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b400 = sequenceBuilder{id: 400, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b400.items = []builder{&b842, &b14}
	b401.items = []builder{&b14, &b400}
	var b396 = sequenceBuilder{id: 396, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b392 = charBuilder{}
	var b393 = charBuilder{}
	var b394 = charBuilder{}
	var b395 = charBuilder{}
	b396.items = []builder{&b392, &b393, &b394, &b395}
	var b403 = sequenceBuilder{id: 403, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b402 = sequenceBuilder{id: 402, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b402.items = []builder{&b842, &b14}
	b403.items = []builder{&b842, &b14, &b402}
	var b399 = sequenceBuilder{id: 399, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b397 = charBuilder{}
	var b398 = charBuilder{}
	b399.items = []builder{&b397, &b398}
	var b405 = sequenceBuilder{id: 405, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b404 = sequenceBuilder{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b404.items = []builder{&b842, &b14}
	b405.items = []builder{&b842, &b14, &b404}
	var b407 = sequenceBuilder{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b406 = sequenceBuilder{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b406.items = []builder{&b842, &b14}
	b407.items = []builder{&b842, &b14, &b406}
	b408.items = []builder{&b401, &b842, &b396, &b403, &b842, &b399, &b405, &b842, &b388, &b407, &b842, &b205}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b423.items = []builder{&b842, &b408}
	b424.items = []builder{&b842, &b408, &b423}
	var b418 = sequenceBuilder{id: 418, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b414.items = []builder{&b842, &b14}
	b415.items = []builder{&b14, &b414}
	var b413 = sequenceBuilder{id: 413, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b409 = charBuilder{}
	var b410 = charBuilder{}
	var b411 = charBuilder{}
	var b412 = charBuilder{}
	b413.items = []builder{&b409, &b410, &b411, &b412}
	var b417 = sequenceBuilder{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b416 = sequenceBuilder{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b416.items = []builder{&b842, &b14}
	b417.items = []builder{&b842, &b14, &b416}
	b418.items = []builder{&b415, &b842, &b413, &b417, &b842, &b205}
	b425.items = []builder{&b391, &b420, &b842, &b388, &b422, &b842, &b205, &b424, &b842, &b418}
	var b482 = sequenceBuilder{id: 482, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b467 = sequenceBuilder{id: 467, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b461 = charBuilder{}
	var b462 = charBuilder{}
	var b463 = charBuilder{}
	var b464 = charBuilder{}
	var b465 = charBuilder{}
	var b466 = charBuilder{}
	b467.items = []builder{&b461, &b462, &b463, &b464, &b465, &b466}
	var b479 = sequenceBuilder{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b478 = sequenceBuilder{id: 478, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b478.items = []builder{&b842, &b14}
	b479.items = []builder{&b842, &b14, &b478}
	var b481 = sequenceBuilder{id: 481, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b480 = sequenceBuilder{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b480.items = []builder{&b842, &b14}
	b481.items = []builder{&b842, &b14, &b480}
	var b469 = sequenceBuilder{id: 469, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b468 = charBuilder{}
	b469.items = []builder{&b468}
	var b475 = sequenceBuilder{id: 475, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b470 = choiceBuilder{id: 470, commit: 2}
	var b460 = sequenceBuilder{id: 460, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b455 = sequenceBuilder{id: 455, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b448 = sequenceBuilder{id: 448, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b444 = charBuilder{}
	var b445 = charBuilder{}
	var b446 = charBuilder{}
	var b447 = charBuilder{}
	b448.items = []builder{&b444, &b445, &b446, &b447}
	var b452 = sequenceBuilder{id: 452, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b451 = sequenceBuilder{id: 451, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b451.items = []builder{&b842, &b14}
	b452.items = []builder{&b842, &b14, &b451}
	var b454 = sequenceBuilder{id: 454, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b453 = sequenceBuilder{id: 453, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b453.items = []builder{&b842, &b14}
	b454.items = []builder{&b842, &b14, &b453}
	var b450 = sequenceBuilder{id: 450, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b449 = charBuilder{}
	b450.items = []builder{&b449}
	b455.items = []builder{&b448, &b452, &b842, &b388, &b454, &b842, &b450}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b457 = sequenceBuilder{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b456 = charBuilder{}
	b457.items = []builder{&b456}
	var b458 = sequenceBuilder{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b458.items = []builder{&b842, &b457}
	b459.items = []builder{&b842, &b457, &b458}
	b460.items = []builder{&b455, &b459, &b842, &b810}
	var b443 = sequenceBuilder{id: 443, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b438 = sequenceBuilder{id: 438, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b433 = sequenceBuilder{id: 433, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b426 = charBuilder{}
	var b427 = charBuilder{}
	var b428 = charBuilder{}
	var b429 = charBuilder{}
	var b430 = charBuilder{}
	var b431 = charBuilder{}
	var b432 = charBuilder{}
	b433.items = []builder{&b426, &b427, &b428, &b429, &b430, &b431, &b432}
	var b437 = sequenceBuilder{id: 437, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b436 = sequenceBuilder{id: 436, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b436.items = []builder{&b842, &b14}
	b437.items = []builder{&b842, &b14, &b436}
	var b435 = sequenceBuilder{id: 435, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b434 = charBuilder{}
	b435.items = []builder{&b434}
	b438.items = []builder{&b433, &b437, &b842, &b435}
	var b442 = sequenceBuilder{id: 442, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b440 = sequenceBuilder{id: 440, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b439 = charBuilder{}
	b440.items = []builder{&b439}
	var b441 = sequenceBuilder{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b441.items = []builder{&b842, &b440}
	b442.items = []builder{&b842, &b440, &b441}
	b443.items = []builder{&b438, &b442, &b842, &b810}
	b470.options = []builder{&b460, &b443}
	var b474 = sequenceBuilder{id: 474, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b472 = sequenceBuilder{id: 472, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b471 = choiceBuilder{id: 471, commit: 2}
	b471.options = []builder{&b460, &b443, &b810}
	b472.items = []builder{&b824, &b842, &b471}
	var b473 = sequenceBuilder{id: 473, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b473.items = []builder{&b842, &b472}
	b474.items = []builder{&b842, &b472, &b473}
	b475.items = []builder{&b470, &b474}
	var b477 = sequenceBuilder{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b476 = charBuilder{}
	b477.items = []builder{&b476}
	b482.items = []builder{&b467, &b479, &b842, &b388, &b481, &b842, &b469, &b842, &b824, &b842, &b475, &b842, &b824, &b842, &b477}
	var b547 = sequenceBuilder{id: 547, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b534 = sequenceBuilder{id: 534, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b528 = charBuilder{}
	var b529 = charBuilder{}
	var b530 = charBuilder{}
	var b531 = charBuilder{}
	var b532 = charBuilder{}
	var b533 = charBuilder{}
	b534.items = []builder{&b528, &b529, &b530, &b531, &b532, &b533}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b545.items = []builder{&b842, &b14}
	b546.items = []builder{&b842, &b14, &b545}
	var b536 = sequenceBuilder{id: 536, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b535 = charBuilder{}
	b536.items = []builder{&b535}
	var b542 = sequenceBuilder{id: 542, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b537 = choiceBuilder{id: 537, commit: 2}
	var b527 = sequenceBuilder{id: 527, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b522 = sequenceBuilder{id: 522, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b515 = sequenceBuilder{id: 515, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b511 = charBuilder{}
	var b512 = charBuilder{}
	var b513 = charBuilder{}
	var b514 = charBuilder{}
	b515.items = []builder{&b511, &b512, &b513, &b514}
	var b519 = sequenceBuilder{id: 519, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b518 = sequenceBuilder{id: 518, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b518.items = []builder{&b842, &b14}
	b519.items = []builder{&b842, &b14, &b518}
	var b510 = choiceBuilder{id: 510, commit: 66}
	var b505 = choiceBuilder{id: 505, commit: 66}
	var b495 = sequenceBuilder{id: 495, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b490 = sequenceBuilder{id: 490, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b487 = charBuilder{}
	var b488 = charBuilder{}
	var b489 = charBuilder{}
	b490.items = []builder{&b487, &b488, &b489}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b491 = sequenceBuilder{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b491.items = []builder{&b842, &b14}
	b492.items = []builder{&b842, &b14, &b491}
	var b494 = sequenceBuilder{id: 494, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b493 = sequenceBuilder{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b493.items = []builder{&b842, &b14}
	b494.items = []builder{&b842, &b14, &b493}
	b495.items = []builder{&b490, &b492, &b842, &b105, &b494, &b842, &b486}
	var b504 = sequenceBuilder{id: 504, commit: 64, name: "receive-assignment", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b499 = sequenceBuilder{id: 499, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b496 = charBuilder{}
	var b497 = charBuilder{}
	var b498 = charBuilder{}
	b499.items = []builder{&b496, &b497, &b498}
	var b501 = sequenceBuilder{id: 501, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b500 = sequenceBuilder{id: 500, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b500.items = []builder{&b842, &b14}
	b501.items = []builder{&b842, &b14, &b500}
	var b503 = sequenceBuilder{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b502 = sequenceBuilder{id: 502, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b502.items = []builder{&b842, &b14}
	b503.items = []builder{&b842, &b14, &b502}
	b504.items = []builder{&b499, &b501, &b842, &b105, &b503, &b842, &b486}
	b505.options = []builder{&b495, &b504}
	b510.options = []builder{&b486, &b505, &b509}
	var b521 = sequenceBuilder{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b520 = sequenceBuilder{id: 520, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b520.items = []builder{&b842, &b14}
	b521.items = []builder{&b842, &b14, &b520}
	var b517 = sequenceBuilder{id: 517, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b516 = charBuilder{}
	b517.items = []builder{&b516}
	b522.items = []builder{&b515, &b519, &b842, &b510, &b521, &b842, &b517}
	var b526 = sequenceBuilder{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b524 = sequenceBuilder{id: 524, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b523 = charBuilder{}
	b524.items = []builder{&b523}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b525.items = []builder{&b842, &b524}
	b526.items = []builder{&b842, &b524, &b525}
	b527.items = []builder{&b522, &b526, &b842, &b810}
	b537.options = []builder{&b527, &b443}
	var b541 = sequenceBuilder{id: 541, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b539 = sequenceBuilder{id: 539, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b538 = choiceBuilder{id: 538, commit: 2}
	b538.options = []builder{&b527, &b443, &b810}
	b539.items = []builder{&b824, &b842, &b538}
	var b540 = sequenceBuilder{id: 540, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b540.items = []builder{&b842, &b539}
	b541.items = []builder{&b842, &b539, &b540}
	b542.items = []builder{&b537, &b541}
	var b544 = sequenceBuilder{id: 544, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b543 = charBuilder{}
	b544.items = []builder{&b543}
	b547.items = []builder{&b534, &b546, &b842, &b536, &b842, &b824, &b842, &b542, &b842, &b824, &b842, &b544}
	var b568 = sequenceBuilder{id: 568, commit: 72, name: "break", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b563 = charBuilder{}
	var b564 = charBuilder{}
	var b565 = charBuilder{}
	var b566 = charBuilder{}
	var b567 = charBuilder{}
	b568.items = []builder{&b563, &b564, &b565, &b566, &b567}
	var b577 = sequenceBuilder{id: 577, commit: 72, name: "continue", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b569 = charBuilder{}
	var b570 = charBuilder{}
	var b571 = charBuilder{}
	var b572 = charBuilder{}
	var b573 = charBuilder{}
	var b574 = charBuilder{}
	var b575 = charBuilder{}
	var b576 = charBuilder{}
	b577.items = []builder{&b569, &b570, &b571, &b572, &b573, &b574, &b575, &b576}
	var b598 = sequenceBuilder{id: 598, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b592 = sequenceBuilder{id: 592, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b589 = charBuilder{}
	var b590 = charBuilder{}
	var b591 = charBuilder{}
	b592.items = []builder{&b589, &b590, &b591}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b593.items = []builder{&b842, &b14}
	b594.items = []builder{&b14, &b593}
	var b588 = choiceBuilder{id: 588, commit: 66}
	var b587 = choiceBuilder{id: 587, commit: 64, name: "range-over-expression"}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b583 = sequenceBuilder{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b582 = sequenceBuilder{id: 582, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b582.items = []builder{&b842, &b14}
	b583.items = []builder{&b842, &b14, &b582}
	var b580 = sequenceBuilder{id: 580, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b578 = charBuilder{}
	var b579 = charBuilder{}
	b580.items = []builder{&b578, &b579}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b584.items = []builder{&b842, &b14}
	b585.items = []builder{&b842, &b14, &b584}
	var b581 = choiceBuilder{id: 581, commit: 2}
	b581.options = []builder{&b388, &b239}
	b586.items = []builder{&b105, &b583, &b842, &b580, &b585, &b842, &b581}
	b587.options = []builder{&b586, &b239}
	b588.options = []builder{&b388, &b587}
	b595.items = []builder{&b594, &b842, &b588}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b596.items = []builder{&b842, &b14}
	b597.items = []builder{&b842, &b14, &b596}
	b598.items = []builder{&b592, &b842, &b595, &b597, &b842, &b205}
	var b746 = choiceBuilder{id: 746, commit: 66}
	var b659 = sequenceBuilder{id: 659, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b655 = sequenceBuilder{id: 655, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b652 = charBuilder{}
	var b653 = charBuilder{}
	var b654 = charBuilder{}
	b655.items = []builder{&b652, &b653, &b654}
	var b658 = sequenceBuilder{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b657 = sequenceBuilder{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b657.items = []builder{&b842, &b14}
	b658.items = []builder{&b842, &b14, &b657}
	var b656 = choiceBuilder{id: 656, commit: 2}
	var b646 = sequenceBuilder{id: 646, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b645 = sequenceBuilder{id: 645, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b641 = sequenceBuilder{id: 641, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b640 = sequenceBuilder{id: 640, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b640.items = []builder{&b842, &b14}
	b641.items = []builder{&b14, &b640}
	var b639 = sequenceBuilder{id: 639, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b638 = charBuilder{}
	b639.items = []builder{&b638}
	b642.items = []builder{&b641, &b842, &b639}
	var b644 = sequenceBuilder{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b643 = sequenceBuilder{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b643.items = []builder{&b842, &b14}
	b644.items = []builder{&b842, &b14, &b643}
	b645.items = []builder{&b105, &b842, &b642, &b644, &b842, &b388}
	b646.items = []builder{&b645}
	var b651 = sequenceBuilder{id: 651, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b648 = sequenceBuilder{id: 648, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b647 = charBuilder{}
	b648.items = []builder{&b647}
	var b650 = sequenceBuilder{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b649 = sequenceBuilder{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b649.items = []builder{&b842, &b14}
	b650.items = []builder{&b842, &b14, &b649}
	b651.items = []builder{&b648, &b650, &b842, &b645}
	b656.options = []builder{&b646, &b651}
	b659.items = []builder{&b655, &b658, &b842, &b656}
	var b680 = sequenceBuilder{id: 680, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b673 = sequenceBuilder{id: 673, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b670 = charBuilder{}
	var b671 = charBuilder{}
	var b672 = charBuilder{}
	b673.items = []builder{&b670, &b671, &b672}
	var b679 = sequenceBuilder{id: 679, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b678 = sequenceBuilder{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b678.items = []builder{&b842, &b14}
	b679.items = []builder{&b842, &b14, &b678}
	var b675 = sequenceBuilder{id: 675, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b674 = charBuilder{}
	b675.items = []builder{&b674}
	var b665 = sequenceBuilder{id: 665, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b660 = choiceBuilder{id: 660, commit: 2}
	b660.options = []builder{&b646, &b651}
	var b664 = sequenceBuilder{id: 664, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b662 = sequenceBuilder{id: 662, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b661 = choiceBuilder{id: 661, commit: 2}
	b661.options = []builder{&b646, &b651}
	b662.items = []builder{&b115, &b842, &b661}
	var b663 = sequenceBuilder{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b663.items = []builder{&b842, &b662}
	b664.items = []builder{&b842, &b662, &b663}
	b665.items = []builder{&b660, &b664}
	var b677 = sequenceBuilder{id: 677, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b676 = charBuilder{}
	b677.items = []builder{&b676}
	b680.items = []builder{&b673, &b679, &b842, &b675, &b842, &b115, &b842, &b665, &b842, &b115, &b842, &b677}
	var b695 = sequenceBuilder{id: 695, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b684 = sequenceBuilder{id: 684, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b681 = charBuilder{}
	var b682 = charBuilder{}
	var b683 = charBuilder{}
	b684.items = []builder{&b681, &b682, &b683}
	var b692 = sequenceBuilder{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b691 = sequenceBuilder{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b691.items = []builder{&b842, &b14}
	b692.items = []builder{&b842, &b14, &b691}
	var b686 = sequenceBuilder{id: 686, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b685 = charBuilder{}
	b686.items = []builder{&b685}
	var b694 = sequenceBuilder{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b693 = sequenceBuilder{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b693.items = []builder{&b842, &b14}
	b694.items = []builder{&b842, &b14, &b693}
	var b688 = sequenceBuilder{id: 688, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b687 = charBuilder{}
	b688.items = []builder{&b687}
	var b669 = sequenceBuilder{id: 669, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b668 = sequenceBuilder{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b666.items = []builder{&b115, &b842, &b646}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b667.items = []builder{&b842, &b666}
	b668.items = []builder{&b842, &b666, &b667}
	b669.items = []builder{&b646, &b668}
	var b690 = sequenceBuilder{id: 690, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b689 = charBuilder{}
	b690.items = []builder{&b689}
	b695.items = []builder{&b684, &b692, &b842, &b686, &b694, &b842, &b688, &b842, &b115, &b842, &b669, &b842, &b115, &b842, &b690}
	var b711 = sequenceBuilder{id: 711, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b707 = sequenceBuilder{id: 707, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b705 = charBuilder{}
	var b706 = charBuilder{}
	b707.items = []builder{&b705, &b706}
	var b710 = sequenceBuilder{id: 710, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b709 = sequenceBuilder{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b709.items = []builder{&b842, &b14}
	b710.items = []builder{&b842, &b14, &b709}
	var b708 = choiceBuilder{id: 708, commit: 2}
	var b699 = sequenceBuilder{id: 699, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b698 = sequenceBuilder{id: 698, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b697 = sequenceBuilder{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b696.items = []builder{&b842, &b14}
	b697.items = []builder{&b842, &b14, &b696}
	b698.items = []builder{&b105, &b697, &b842, &b214}
	b699.items = []builder{&b698}
	var b704 = sequenceBuilder{id: 704, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b701 = sequenceBuilder{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b700 = charBuilder{}
	b701.items = []builder{&b700}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b702 = sequenceBuilder{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b702.items = []builder{&b842, &b14}
	b703.items = []builder{&b842, &b14, &b702}
	b704.items = []builder{&b701, &b703, &b842, &b698}
	b708.options = []builder{&b699, &b704}
	b711.items = []builder{&b707, &b710, &b842, &b708}
	var b731 = sequenceBuilder{id: 731, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b724 = sequenceBuilder{id: 724, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b722 = charBuilder{}
	var b723 = charBuilder{}
	b724.items = []builder{&b722, &b723}
	var b730 = sequenceBuilder{id: 730, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b729 = sequenceBuilder{id: 729, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b729.items = []builder{&b842, &b14}
	b730.items = []builder{&b842, &b14, &b729}
	var b726 = sequenceBuilder{id: 726, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b725 = charBuilder{}
	b726.items = []builder{&b725}
	var b721 = sequenceBuilder{id: 721, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b716 = choiceBuilder{id: 716, commit: 2}
	b716.options = []builder{&b699, &b704}
	var b720 = sequenceBuilder{id: 720, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b718 = sequenceBuilder{id: 718, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b717 = choiceBuilder{id: 717, commit: 2}
	b717.options = []builder{&b699, &b704}
	b718.items = []builder{&b115, &b842, &b717}
	var b719 = sequenceBuilder{id: 719, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b719.items = []builder{&b842, &b718}
	b720.items = []builder{&b842, &b718, &b719}
	b721.items = []builder{&b716, &b720}
	var b728 = sequenceBuilder{id: 728, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b727 = charBuilder{}
	b728.items = []builder{&b727}
	b731.items = []builder{&b724, &b730, &b842, &b726, &b842, &b115, &b842, &b721, &b842, &b115, &b842, &b728}
	var b745 = sequenceBuilder{id: 745, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b734 = sequenceBuilder{id: 734, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b732 = charBuilder{}
	var b733 = charBuilder{}
	b734.items = []builder{&b732, &b733}
	var b742 = sequenceBuilder{id: 742, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b741 = sequenceBuilder{id: 741, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b741.items = []builder{&b842, &b14}
	b742.items = []builder{&b842, &b14, &b741}
	var b736 = sequenceBuilder{id: 736, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b735 = charBuilder{}
	b736.items = []builder{&b735}
	var b744 = sequenceBuilder{id: 744, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b743 = sequenceBuilder{id: 743, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b743.items = []builder{&b842, &b14}
	b744.items = []builder{&b842, &b14, &b743}
	var b738 = sequenceBuilder{id: 738, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b737 = charBuilder{}
	b738.items = []builder{&b737}
	var b715 = sequenceBuilder{id: 715, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b714 = sequenceBuilder{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b712 = sequenceBuilder{id: 712, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b712.items = []builder{&b115, &b842, &b699}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b713.items = []builder{&b842, &b712}
	b714.items = []builder{&b842, &b712, &b713}
	b715.items = []builder{&b699, &b714}
	var b740 = sequenceBuilder{id: 740, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b739 = charBuilder{}
	b740.items = []builder{&b739}
	b745.items = []builder{&b734, &b742, &b842, &b736, &b744, &b842, &b738, &b842, &b115, &b842, &b715, &b842, &b115, &b842, &b740}
	b746.options = []builder{&b659, &b680, &b695, &b711, &b731, &b745}
	var b789 = choiceBuilder{id: 789, commit: 64, name: "require"}
	var b773 = sequenceBuilder{id: 773, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b770 = sequenceBuilder{id: 770, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b763 = charBuilder{}
	var b764 = charBuilder{}
	var b765 = charBuilder{}
	var b766 = charBuilder{}
	var b767 = charBuilder{}
	var b768 = charBuilder{}
	var b769 = charBuilder{}
	b770.items = []builder{&b763, &b764, &b765, &b766, &b767, &b768, &b769}
	var b772 = sequenceBuilder{id: 772, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b771 = sequenceBuilder{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b771.items = []builder{&b842, &b14}
	b772.items = []builder{&b842, &b14, &b771}
	var b758 = choiceBuilder{id: 758, commit: 64, name: "require-fact"}
	var b757 = sequenceBuilder{id: 757, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b749 = choiceBuilder{id: 749, commit: 2}
	var b748 = sequenceBuilder{id: 748, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b747 = charBuilder{}
	b748.items = []builder{&b747}
	b749.options = []builder{&b105, &b748}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b752.items = []builder{&b842, &b14}
	b753.items = []builder{&b14, &b752}
	var b751 = sequenceBuilder{id: 751, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b750 = charBuilder{}
	b751.items = []builder{&b750}
	b754.items = []builder{&b753, &b842, &b751}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b755.items = []builder{&b842, &b14}
	b756.items = []builder{&b842, &b14, &b755}
	b757.items = []builder{&b749, &b842, &b754, &b756, &b842, &b88}
	b758.options = []builder{&b88, &b757}
	b773.items = []builder{&b770, &b772, &b842, &b758}
	var b788 = sequenceBuilder{id: 788, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b781 = sequenceBuilder{id: 781, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b774 = charBuilder{}
	var b775 = charBuilder{}
	var b776 = charBuilder{}
	var b777 = charBuilder{}
	var b778 = charBuilder{}
	var b779 = charBuilder{}
	var b780 = charBuilder{}
	b781.items = []builder{&b774, &b775, &b776, &b777, &b778, &b779, &b780}
	var b787 = sequenceBuilder{id: 787, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b786 = sequenceBuilder{id: 786, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b786.items = []builder{&b842, &b14}
	b787.items = []builder{&b842, &b14, &b786}
	var b783 = sequenceBuilder{id: 783, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b782 = charBuilder{}
	b783.items = []builder{&b782}
	var b762 = sequenceBuilder{id: 762, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b761 = sequenceBuilder{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b759 = sequenceBuilder{id: 759, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b759.items = []builder{&b115, &b842, &b758}
	var b760 = sequenceBuilder{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b760.items = []builder{&b842, &b759}
	b761.items = []builder{&b842, &b759, &b760}
	b762.items = []builder{&b758, &b761}
	var b785 = sequenceBuilder{id: 785, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b784 = charBuilder{}
	b785.items = []builder{&b784}
	b788.items = []builder{&b781, &b787, &b842, &b783, &b842, &b115, &b842, &b762, &b842, &b115, &b842, &b785}
	b789.options = []builder{&b773, &b788}
	var b799 = sequenceBuilder{id: 799, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b796 = sequenceBuilder{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b790 = charBuilder{}
	var b791 = charBuilder{}
	var b792 = charBuilder{}
	var b793 = charBuilder{}
	var b794 = charBuilder{}
	var b795 = charBuilder{}
	b796.items = []builder{&b790, &b791, &b792, &b793, &b794, &b795}
	var b798 = sequenceBuilder{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b797 = sequenceBuilder{id: 797, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b797.items = []builder{&b842, &b14}
	b798.items = []builder{&b842, &b14, &b797}
	b799.items = []builder{&b796, &b798, &b842, &b746}
	var b819 = sequenceBuilder{id: 819, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b812 = sequenceBuilder{id: 812, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b811 = charBuilder{}
	b812.items = []builder{&b811}
	var b816 = sequenceBuilder{id: 816, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b815 = sequenceBuilder{id: 815, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b815.items = []builder{&b842, &b14}
	b816.items = []builder{&b842, &b14, &b815}
	var b818 = sequenceBuilder{id: 818, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b817 = sequenceBuilder{id: 817, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b817.items = []builder{&b842, &b14}
	b818.items = []builder{&b842, &b14, &b817}
	var b814 = sequenceBuilder{id: 814, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b813 = charBuilder{}
	b814.items = []builder{&b813}
	b819.items = []builder{&b812, &b816, &b842, &b810, &b818, &b842, &b814}
	b810.options = []builder{&b800, &b200, &b425, &b482, &b547, &b568, &b577, &b598, &b746, &b789, &b799, &b819}
	var b827 = sequenceBuilder{id: 827, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b825 = sequenceBuilder{id: 825, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b825.items = []builder{&b824, &b842, &b810}
	var b826 = sequenceBuilder{id: 826, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b826.items = []builder{&b842, &b825}
	b827.items = []builder{&b842, &b825, &b826}
	b828.items = []builder{&b810, &b827}
	b843.items = []builder{&b839, &b842, &b824, &b842, &b828, &b842, &b824}
	b844.items = []builder{&b842, &b843, &b842}

	return parseInput(r, &p844, &b844)
}
