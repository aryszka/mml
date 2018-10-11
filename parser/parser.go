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
	var p835 = sequenceParser{id: 835, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p833 = choiceParser{id: 833, commit: 2}
	var p831 = choiceParser{id: 831, commit: 70, name: "ws", generalizations: []int{833}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{831, 833}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{831, 833}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{831, 833}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{831, 833}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{831, 833}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{831, 833}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p831.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p832 = sequenceParser{id: 832, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{833}}
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
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{813, 111}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p39.items = []parser{&p14, &p833, &p38}
	var p40 = sequenceParser{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p40.items = []parser{&p833, &p39}
	p41.items = []parser{&p833, &p39, &p40}
	p42.items = []parser{&p38, &p41}
	p832.items = []parser{&p42}
	p833.options = []parser{&p831, &p832}
	var p834 = sequenceParser{id: 834, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p830 = sequenceParser{id: 830, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p827 = sequenceParser{id: 827, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p825 = charParser{id: 825, chars: []rune{35}}
	var p826 = charParser{id: 826, chars: []rune{33}}
	p827.items = []parser{&p825, &p826}
	var p824 = sequenceParser{id: 824, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p823 = sequenceParser{id: 823, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p821 = sequenceParser{id: 821, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p820 = charParser{id: 820, not: true, chars: []rune{10}}
	p821.items = []parser{&p820}
	var p822 = sequenceParser{id: 822, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p822.items = []parser{&p833, &p821}
	p823.items = []parser{&p821, &p822}
	p824.items = []parser{&p823}
	var p829 = sequenceParser{id: 829, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p828 = charParser{id: 828, chars: []rune{10}}
	p829.items = []parser{&p828}
	p830.items = []parser{&p827, &p833, &p824, &p833, &p829}
	var p815 = sequenceParser{id: 815, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p813 = choiceParser{id: 813, commit: 2}
	var p812 = sequenceParser{id: 812, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{813}}
	var p811 = charParser{id: 811, chars: []rune{59}}
	p812.items = []parser{&p811}
	p813.options = []parser{&p812, &p14}
	var p814 = sequenceParser{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p814.items = []parser{&p833, &p813}
	p815.items = []parser{&p813, &p814}
	var p819 = sequenceParser{id: 819, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p801 = choiceParser{id: 801, commit: 66, name: "statement", generalizations: []int{483, 547}}
	var p185 = sequenceParser{id: 185, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{801, 483, 547}}
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
	p182.items = []parser{&p833, &p14}
	p183.items = []parser{&p14, &p182}
	var p400 = choiceParser{id: 400, commit: 66, name: "expression", generalizations: []int{114, 791, 197, 582, 575, 801}}
	var p271 = choiceParser{id: 271, commit: 66, name: "primary-expression", generalizations: []int{114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p60 = choiceParser{id: 60, commit: 64, name: "int", generalizations: []int{271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p51 = sequenceParser{id: 51, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p50 = sequenceParser{id: 50, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p49 = charParser{id: 49, ranges: [][]rune{{49, 57}}}
	p50.items = []parser{&p49}
	var p44 = sequenceParser{id: 44, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p43 = charParser{id: 43, ranges: [][]rune{{48, 57}}}
	p44.items = []parser{&p43}
	p51.items = []parser{&p50, &p44}
	var p54 = sequenceParser{id: 54, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p53 = sequenceParser{id: 53, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p52 = charParser{id: 52, chars: []rune{48}}
	p53.items = []parser{&p52}
	var p46 = sequenceParser{id: 46, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 55}}}
	p46.items = []parser{&p45}
	p54.items = []parser{&p53, &p46}
	var p59 = sequenceParser{id: 59, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{60, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
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
	var p73 = choiceParser{id: 73, commit: 72, name: "float", generalizations: []int{271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p68 = sequenceParser{id: 68, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{73, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
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
	var p71 = sequenceParser{id: 71, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{73, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p70 = sequenceParser{id: 70, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p69 = charParser{id: 69, chars: []rune{46}}
	p70.items = []parser{&p69}
	p71.items = []parser{&p70, &p44, &p65}
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{73, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	p72.items = []parser{&p44, &p65}
	p73.options = []parser{&p68, &p71, &p72}
	var p86 = sequenceParser{id: 86, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 114, 139, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 757, 801}}
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
	var p98 = choiceParser{id: 98, commit: 66, name: "bool", generalizations: []int{271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p91 = sequenceParser{id: 91, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p87 = charParser{id: 87, chars: []rune{116}}
	var p88 = charParser{id: 88, chars: []rune{114}}
	var p89 = charParser{id: 89, chars: []rune{117}}
	var p90 = charParser{id: 90, chars: []rune{101}}
	p91.items = []parser{&p87, &p88, &p89, &p90}
	var p97 = sequenceParser{id: 97, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p92 = charParser{id: 92, chars: []rune{102}}
	var p93 = charParser{id: 93, chars: []rune{97}}
	var p94 = charParser{id: 94, chars: []rune{108}}
	var p95 = charParser{id: 95, chars: []rune{115}}
	var p96 = charParser{id: 96, chars: []rune{101}}
	p97.items = []parser{&p92, &p93, &p94, &p95, &p96}
	p98.options = []parser{&p91, &p97}
	var p515 = sequenceParser{id: 515, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 114, 791, 197, 400, 337, 338, 339, 340, 341, 392, 519, 582, 575, 801}}
	var p512 = sequenceParser{id: 512, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p505 = charParser{id: 505, chars: []rune{114}}
	var p506 = charParser{id: 506, chars: []rune{101}}
	var p507 = charParser{id: 507, chars: []rune{99}}
	var p508 = charParser{id: 508, chars: []rune{101}}
	var p509 = charParser{id: 509, chars: []rune{105}}
	var p510 = charParser{id: 510, chars: []rune{118}}
	var p511 = charParser{id: 511, chars: []rune{101}}
	p512.items = []parser{&p505, &p506, &p507, &p508, &p509, &p510, &p511}
	var p514 = sequenceParser{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p513 = sequenceParser{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p513.items = []parser{&p833, &p14}
	p514.items = []parser{&p833, &p14, &p513}
	p515.items = []parser{&p512, &p514, &p833, &p271}
	var p103 = sequenceParser{id: 103, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{271, 114, 139, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 748, 801}}
	var p100 = sequenceParser{id: 100, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p99 = charParser{id: 99, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p100.items = []parser{&p99}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p102.items = []parser{&p101}
	p103.items = []parser{&p100, &p102}
	var p124 = sequenceParser{id: 124, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{114, 271, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
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
	p112.items = []parser{&p833, &p111}
	p113.items = []parser{&p111, &p112}
	var p118 = sequenceParser{id: 118, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p114 = choiceParser{id: 114, commit: 66, name: "list-item"}
	var p108 = sequenceParser{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{114, 147, 148}}
	var p107 = sequenceParser{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p104 = charParser{id: 104, chars: []rune{46}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	p107.items = []parser{&p104, &p105, &p106}
	p108.items = []parser{&p271, &p833, &p107}
	p114.options = []parser{&p400, &p108}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p115 = sequenceParser{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p115.items = []parser{&p113, &p833, &p114}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p116.items = []parser{&p833, &p115}
	p117.items = []parser{&p833, &p115, &p116}
	p118.items = []parser{&p114, &p117}
	var p122 = sequenceParser{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p121 = charParser{id: 121, chars: []rune{93}}
	p122.items = []parser{&p121}
	p123.items = []parser{&p120, &p833, &p113, &p833, &p118, &p833, &p113, &p833, &p122}
	p124.items = []parser{&p123}
	var p129 = sequenceParser{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p126 = sequenceParser{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p125 = charParser{id: 125, chars: []rune{126}}
	p126.items = []parser{&p125}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p127 = sequenceParser{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p127.items = []parser{&p833, &p14}
	p128.items = []parser{&p833, &p14, &p127}
	p129.items = []parser{&p126, &p128, &p833, &p123}
	var p158 = sequenceParser{id: 158, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{271, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
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
	p134.items = []parser{&p833, &p14}
	p135.items = []parser{&p833, &p14, &p134}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p136 = sequenceParser{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p136.items = []parser{&p833, &p14}
	p137.items = []parser{&p833, &p14, &p136}
	var p133 = sequenceParser{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p132 = charParser{id: 132, chars: []rune{93}}
	p133.items = []parser{&p132}
	p138.items = []parser{&p131, &p135, &p833, &p400, &p137, &p833, &p133}
	p139.options = []parser{&p103, &p86, &p138}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p142 = sequenceParser{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p142.items = []parser{&p833, &p14}
	p143.items = []parser{&p833, &p14, &p142}
	var p141 = sequenceParser{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p140 = charParser{id: 140, chars: []rune{58}}
	p141.items = []parser{&p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p833, &p14}
	p145.items = []parser{&p833, &p14, &p144}
	p146.items = []parser{&p139, &p143, &p833, &p141, &p145, &p833, &p400}
	p147.options = []parser{&p146, &p108}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p149 = sequenceParser{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p148 = choiceParser{id: 148, commit: 2}
	p148.options = []parser{&p146, &p108}
	p149.items = []parser{&p113, &p833, &p148}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p150.items = []parser{&p833, &p149}
	p151.items = []parser{&p833, &p149, &p150}
	p152.items = []parser{&p147, &p151}
	var p156 = sequenceParser{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p155 = charParser{id: 155, chars: []rune{125}}
	p156.items = []parser{&p155}
	p157.items = []parser{&p154, &p833, &p113, &p833, &p152, &p833, &p113, &p833, &p156}
	p158.items = []parser{&p157}
	var p163 = sequenceParser{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 791, 197, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p160 = sequenceParser{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p159 = charParser{id: 159, chars: []rune{126}}
	p160.items = []parser{&p159}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p161 = sequenceParser{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p161.items = []parser{&p833, &p14}
	p162.items = []parser{&p833, &p14, &p161}
	p163.items = []parser{&p160, &p162, &p833, &p157}
	var p206 = sequenceParser{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{791, 197, 271, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p203 = sequenceParser{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p201 = charParser{id: 201, chars: []rune{102}}
	var p202 = charParser{id: 202, chars: []rune{110}}
	p203.items = []parser{&p201, &p202}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p204 = sequenceParser{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p204.items = []parser{&p833, &p14}
	p205.items = []parser{&p833, &p14, &p204}
	var p200 = sequenceParser{id: 200, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p192 = sequenceParser{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p191 = charParser{id: 191, chars: []rune{40}}
	p192.items = []parser{&p191}
	var p194 = choiceParser{id: 194, commit: 2}
	var p167 = sequenceParser{id: 167, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{194}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p164.items = []parser{&p113, &p833, &p103}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p165.items = []parser{&p833, &p164}
	p166.items = []parser{&p833, &p164, &p165}
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
	p172.items = []parser{&p833, &p14}
	p173.items = []parser{&p833, &p14, &p172}
	p174.items = []parser{&p171, &p173, &p833, &p103}
	p193.items = []parser{&p167, &p833, &p113, &p833, &p174}
	p194.options = []parser{&p167, &p193, &p174}
	var p196 = sequenceParser{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p195 = charParser{id: 195, chars: []rune{41}}
	p196.items = []parser{&p195}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p198 = sequenceParser{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p198.items = []parser{&p833, &p14}
	p199.items = []parser{&p833, &p14, &p198}
	var p197 = choiceParser{id: 197, commit: 2}
	var p791 = choiceParser{id: 791, commit: 66, name: "simple-statement", generalizations: []int{197, 801}}
	var p504 = sequenceParser{id: 504, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{791, 197, 519, 801}}
	var p499 = sequenceParser{id: 499, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p495 = charParser{id: 495, chars: []rune{115}}
	var p496 = charParser{id: 496, chars: []rune{101}}
	var p497 = charParser{id: 497, chars: []rune{110}}
	var p498 = charParser{id: 498, chars: []rune{100}}
	p499.items = []parser{&p495, &p496, &p497, &p498}
	var p501 = sequenceParser{id: 501, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p500 = sequenceParser{id: 500, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p500.items = []parser{&p833, &p14}
	p501.items = []parser{&p833, &p14, &p500}
	var p503 = sequenceParser{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p502 = sequenceParser{id: 502, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p502.items = []parser{&p833, &p14}
	p503.items = []parser{&p833, &p14, &p502}
	p504.items = []parser{&p499, &p501, &p833, &p271, &p503, &p833, &p271}
	var p562 = sequenceParser{id: 562, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{791, 197, 801}}
	var p559 = sequenceParser{id: 559, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p557 = charParser{id: 557, chars: []rune{103}}
	var p558 = charParser{id: 558, chars: []rune{111}}
	p559.items = []parser{&p557, &p558}
	var p561 = sequenceParser{id: 561, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p560 = sequenceParser{id: 560, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p560.items = []parser{&p833, &p14}
	p561.items = []parser{&p833, &p14, &p560}
	var p261 = sequenceParser{id: 261, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p258 = sequenceParser{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p257 = charParser{id: 257, chars: []rune{40}}
	p258.items = []parser{&p257}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{41}}
	p260.items = []parser{&p259}
	p261.items = []parser{&p271, &p833, &p258, &p833, &p113, &p833, &p118, &p833, &p113, &p833, &p260}
	p562.items = []parser{&p559, &p561, &p833, &p261}
	var p571 = sequenceParser{id: 571, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{791, 197, 801}}
	var p568 = sequenceParser{id: 568, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p563 = charParser{id: 563, chars: []rune{100}}
	var p564 = charParser{id: 564, chars: []rune{101}}
	var p565 = charParser{id: 565, chars: []rune{102}}
	var p566 = charParser{id: 566, chars: []rune{101}}
	var p567 = charParser{id: 567, chars: []rune{114}}
	p568.items = []parser{&p563, &p564, &p565, &p566, &p567}
	var p570 = sequenceParser{id: 570, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p569 = sequenceParser{id: 569, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p569.items = []parser{&p833, &p14}
	p570.items = []parser{&p833, &p14, &p569}
	p571.items = []parser{&p568, &p570, &p833, &p261}
	var p636 = choiceParser{id: 636, commit: 64, name: "assignment", generalizations: []int{791, 197, 801}}
	var p616 = sequenceParser{id: 616, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{636, 791, 197, 801}}
	var p613 = sequenceParser{id: 613, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p610 = charParser{id: 610, chars: []rune{115}}
	var p611 = charParser{id: 611, chars: []rune{101}}
	var p612 = charParser{id: 612, chars: []rune{116}}
	p613.items = []parser{&p610, &p611, &p612}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p614.items = []parser{&p833, &p14}
	p615.items = []parser{&p833, &p14, &p614}
	var p605 = sequenceParser{id: 605, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p602 = sequenceParser{id: 602, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p601 = sequenceParser{id: 601, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p600 = sequenceParser{id: 600, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p600.items = []parser{&p833, &p14}
	p601.items = []parser{&p14, &p600}
	var p599 = sequenceParser{id: 599, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p598 = charParser{id: 598, chars: []rune{61}}
	p599.items = []parser{&p598}
	p602.items = []parser{&p601, &p833, &p599}
	var p604 = sequenceParser{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p603 = sequenceParser{id: 603, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p603.items = []parser{&p833, &p14}
	p604.items = []parser{&p833, &p14, &p603}
	p605.items = []parser{&p271, &p833, &p602, &p604, &p833, &p400}
	p616.items = []parser{&p613, &p615, &p833, &p605}
	var p623 = sequenceParser{id: 623, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{636, 791, 197, 801}}
	var p620 = sequenceParser{id: 620, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p619 = sequenceParser{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p619.items = []parser{&p833, &p14}
	p620.items = []parser{&p833, &p14, &p619}
	var p618 = sequenceParser{id: 618, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p617 = charParser{id: 617, chars: []rune{61}}
	p618.items = []parser{&p617}
	var p622 = sequenceParser{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p621 = sequenceParser{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p621.items = []parser{&p833, &p14}
	p622.items = []parser{&p833, &p14, &p621}
	p623.items = []parser{&p271, &p620, &p833, &p618, &p622, &p833, &p400}
	var p635 = sequenceParser{id: 635, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{636, 791, 197, 801}}
	var p627 = sequenceParser{id: 627, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p624 = charParser{id: 624, chars: []rune{115}}
	var p625 = charParser{id: 625, chars: []rune{101}}
	var p626 = charParser{id: 626, chars: []rune{116}}
	p627.items = []parser{&p624, &p625, &p626}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p633 = sequenceParser{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p633.items = []parser{&p833, &p14}
	p634.items = []parser{&p833, &p14, &p633}
	var p629 = sequenceParser{id: 629, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p628 = charParser{id: 628, chars: []rune{40}}
	p629.items = []parser{&p628}
	var p630 = sequenceParser{id: 630, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p609 = sequenceParser{id: 609, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p608 = sequenceParser{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p606 = sequenceParser{id: 606, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p606.items = []parser{&p113, &p833, &p605}
	var p607 = sequenceParser{id: 607, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p607.items = []parser{&p833, &p606}
	p608.items = []parser{&p833, &p606, &p607}
	p609.items = []parser{&p605, &p608}
	p630.items = []parser{&p113, &p833, &p609}
	var p632 = sequenceParser{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p631 = charParser{id: 631, chars: []rune{41}}
	p632.items = []parser{&p631}
	p635.items = []parser{&p627, &p634, &p833, &p629, &p833, &p630, &p833, &p113, &p833, &p632}
	p636.options = []parser{&p616, &p623, &p635}
	var p800 = sequenceParser{id: 800, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{791, 197, 801}}
	var p793 = sequenceParser{id: 793, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p792 = charParser{id: 792, chars: []rune{40}}
	p793.items = []parser{&p792}
	var p797 = sequenceParser{id: 797, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p796 = sequenceParser{id: 796, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p796.items = []parser{&p833, &p14}
	p797.items = []parser{&p833, &p14, &p796}
	var p799 = sequenceParser{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p798 = sequenceParser{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p798.items = []parser{&p833, &p14}
	p799.items = []parser{&p833, &p14, &p798}
	var p795 = sequenceParser{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p794 = charParser{id: 794, chars: []rune{41}}
	p795.items = []parser{&p794}
	p800.items = []parser{&p793, &p797, &p833, &p791, &p799, &p833, &p795}
	p791.options = []parser{&p504, &p562, &p571, &p636, &p800, &p400}
	var p190 = sequenceParser{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{197}}
	var p187 = sequenceParser{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p186 = charParser{id: 186, chars: []rune{123}}
	p187.items = []parser{&p186}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{125}}
	p189.items = []parser{&p188}
	p190.items = []parser{&p187, &p833, &p815, &p833, &p819, &p833, &p815, &p833, &p189}
	p197.options = []parser{&p791, &p190}
	p200.items = []parser{&p192, &p833, &p113, &p833, &p194, &p833, &p113, &p833, &p196, &p199, &p833, &p197}
	p206.items = []parser{&p203, &p205, &p833, &p200}
	var p216 = sequenceParser{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p209 = sequenceParser{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p207 = charParser{id: 207, chars: []rune{102}}
	var p208 = charParser{id: 208, chars: []rune{110}}
	p209.items = []parser{&p207, &p208}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p212 = sequenceParser{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p212.items = []parser{&p833, &p14}
	p213.items = []parser{&p833, &p14, &p212}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p210 = charParser{id: 210, chars: []rune{126}}
	p211.items = []parser{&p210}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p214.items = []parser{&p833, &p14}
	p215.items = []parser{&p833, &p14, &p214}
	p216.items = []parser{&p209, &p213, &p833, &p211, &p215, &p833, &p200}
	var p256 = sequenceParser{id: 256, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p255 = sequenceParser{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p254 = sequenceParser{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p254.items = []parser{&p833, &p14}
	p255.items = []parser{&p833, &p14, &p254}
	var p253 = sequenceParser{id: 253, commit: 66, name: "index-list", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var p249 = choiceParser{id: 249, commit: 66, name: "index"}
	var p230 = sequenceParser{id: 230, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{249}}
	var p227 = sequenceParser{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p226 = charParser{id: 226, chars: []rune{46}}
	p227.items = []parser{&p226}
	var p229 = sequenceParser{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p228 = sequenceParser{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p228.items = []parser{&p833, &p14}
	p229.items = []parser{&p833, &p14, &p228}
	p230.items = []parser{&p227, &p229, &p833, &p103}
	var p239 = sequenceParser{id: 239, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{249}}
	var p232 = sequenceParser{id: 232, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p231 = charParser{id: 231, chars: []rune{91}}
	p232.items = []parser{&p231}
	var p236 = sequenceParser{id: 236, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p235 = sequenceParser{id: 235, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p235.items = []parser{&p833, &p14}
	p236.items = []parser{&p833, &p14, &p235}
	var p238 = sequenceParser{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p237 = sequenceParser{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p237.items = []parser{&p833, &p14}
	p238.items = []parser{&p833, &p14, &p237}
	var p234 = sequenceParser{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p233 = charParser{id: 233, chars: []rune{93}}
	p234.items = []parser{&p233}
	p239.items = []parser{&p232, &p236, &p833, &p400, &p238, &p833, &p234}
	var p248 = sequenceParser{id: 248, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{249}}
	var p241 = sequenceParser{id: 241, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p240 = charParser{id: 240, chars: []rune{91}}
	p241.items = []parser{&p240}
	var p245 = sequenceParser{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p244 = sequenceParser{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p244.items = []parser{&p833, &p14}
	p245.items = []parser{&p833, &p14, &p244}
	var p225 = sequenceParser{id: 225, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{575, 581, 582}}
	var p217 = sequenceParser{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p217.items = []parser{&p400}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p221 = sequenceParser{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p221.items = []parser{&p833, &p14}
	p222.items = []parser{&p833, &p14, &p221}
	var p220 = sequenceParser{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p219 = charParser{id: 219, chars: []rune{58}}
	p220.items = []parser{&p219}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p833, &p14}
	p224.items = []parser{&p833, &p14, &p223}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p400}
	p225.items = []parser{&p217, &p222, &p833, &p220, &p224, &p833, &p218}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p246 = sequenceParser{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p246.items = []parser{&p833, &p14}
	p247.items = []parser{&p833, &p14, &p246}
	var p243 = sequenceParser{id: 243, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p242 = charParser{id: 242, chars: []rune{93}}
	p243.items = []parser{&p242}
	p248.items = []parser{&p241, &p245, &p833, &p225, &p247, &p833, &p243}
	p249.options = []parser{&p230, &p239, &p248}
	var p252 = sequenceParser{id: 252, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p251 = sequenceParser{id: 251, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p250.items = []parser{&p833, &p14}
	p251.items = []parser{&p14, &p250}
	p252.items = []parser{&p251, &p833, &p249}
	p253.items = []parser{&p249, &p833, &p252}
	p256.items = []parser{&p271, &p255, &p833, &p253}
	var p270 = sequenceParser{id: 270, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{271, 400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p263 = sequenceParser{id: 263, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p262 = charParser{id: 262, chars: []rune{40}}
	p263.items = []parser{&p262}
	var p267 = sequenceParser{id: 267, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p266 = sequenceParser{id: 266, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p266.items = []parser{&p833, &p14}
	p267.items = []parser{&p833, &p14, &p266}
	var p269 = sequenceParser{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p268 = sequenceParser{id: 268, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p268.items = []parser{&p833, &p14}
	p269.items = []parser{&p833, &p14, &p268}
	var p265 = sequenceParser{id: 265, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p264 = charParser{id: 264, chars: []rune{41}}
	p265.items = []parser{&p264}
	p270.items = []parser{&p263, &p267, &p833, &p400, &p269, &p833, &p265}
	p271.options = []parser{&p60, &p73, &p86, &p98, &p515, &p103, &p124, &p129, &p158, &p163, &p206, &p216, &p256, &p261, &p270}
	var p331 = sequenceParser{id: 331, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{400, 337, 338, 339, 340, 341, 392, 582, 575, 801}}
	var p330 = choiceParser{id: 330, commit: 66, name: "unary-operator"}
	var p290 = sequenceParser{id: 290, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p289 = charParser{id: 289, chars: []rune{43}}
	p290.items = []parser{&p289}
	var p292 = sequenceParser{id: 292, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p291 = charParser{id: 291, chars: []rune{45}}
	p292.items = []parser{&p291}
	var p273 = sequenceParser{id: 273, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p272 = charParser{id: 272, chars: []rune{94}}
	p273.items = []parser{&p272}
	var p304 = sequenceParser{id: 304, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p303 = charParser{id: 303, chars: []rune{33}}
	p304.items = []parser{&p303}
	p330.options = []parser{&p290, &p292, &p273, &p304}
	p331.items = []parser{&p330, &p833, &p271}
	var p378 = choiceParser{id: 378, commit: 66, name: "binary-expression", generalizations: []int{400, 392, 582, 575, 801}}
	var p349 = sequenceParser{id: 349, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{378, 338, 339, 340, 341, 400, 392, 582, 575, 801}}
	var p337 = choiceParser{id: 337, commit: 66, name: "operand0", generalizations: []int{338, 339, 340, 341}}
	p337.options = []parser{&p271, &p331}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p344 = sequenceParser{id: 344, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p343.items = []parser{&p833, &p14}
	p344.items = []parser{&p14, &p343}
	var p332 = choiceParser{id: 332, commit: 66, name: "binary-op0"}
	var p275 = sequenceParser{id: 275, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p274 = charParser{id: 274, chars: []rune{38}}
	p275.items = []parser{&p274}
	var p282 = sequenceParser{id: 282, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{332}}
	var p280 = charParser{id: 280, chars: []rune{38}}
	var p281 = charParser{id: 281, chars: []rune{94}}
	p282.items = []parser{&p280, &p281}
	var p285 = sequenceParser{id: 285, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{332}}
	var p283 = charParser{id: 283, chars: []rune{60}}
	var p284 = charParser{id: 284, chars: []rune{60}}
	p285.items = []parser{&p283, &p284}
	var p288 = sequenceParser{id: 288, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{332}}
	var p286 = charParser{id: 286, chars: []rune{62}}
	var p287 = charParser{id: 287, chars: []rune{62}}
	p288.items = []parser{&p286, &p287}
	var p294 = sequenceParser{id: 294, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p293 = charParser{id: 293, chars: []rune{42}}
	p294.items = []parser{&p293}
	var p296 = sequenceParser{id: 296, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p295 = charParser{id: 295, chars: []rune{47}}
	p296.items = []parser{&p295}
	var p298 = sequenceParser{id: 298, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p297 = charParser{id: 297, chars: []rune{37}}
	p298.items = []parser{&p297}
	p332.options = []parser{&p275, &p282, &p285, &p288, &p294, &p296, &p298}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p345 = sequenceParser{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p345.items = []parser{&p833, &p14}
	p346.items = []parser{&p833, &p14, &p345}
	p347.items = []parser{&p344, &p833, &p332, &p346, &p833, &p337}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p348.items = []parser{&p833, &p347}
	p349.items = []parser{&p337, &p833, &p347, &p348}
	var p356 = sequenceParser{id: 356, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{378, 339, 340, 341, 400, 392, 582, 575, 801}}
	var p338 = choiceParser{id: 338, commit: 66, name: "operand1", generalizations: []int{339, 340, 341}}
	p338.options = []parser{&p337, &p349}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p351 = sequenceParser{id: 351, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p350.items = []parser{&p833, &p14}
	p351.items = []parser{&p14, &p350}
	var p333 = choiceParser{id: 333, commit: 66, name: "binary-op1"}
	var p277 = sequenceParser{id: 277, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p276 = charParser{id: 276, chars: []rune{124}}
	p277.items = []parser{&p276}
	var p279 = sequenceParser{id: 279, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p278 = charParser{id: 278, chars: []rune{94}}
	p279.items = []parser{&p278}
	var p300 = sequenceParser{id: 300, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p299 = charParser{id: 299, chars: []rune{43}}
	p300.items = []parser{&p299}
	var p302 = sequenceParser{id: 302, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p301 = charParser{id: 301, chars: []rune{45}}
	p302.items = []parser{&p301}
	p333.options = []parser{&p277, &p279, &p300, &p302}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p352.items = []parser{&p833, &p14}
	p353.items = []parser{&p833, &p14, &p352}
	p354.items = []parser{&p351, &p833, &p333, &p353, &p833, &p338}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p355.items = []parser{&p833, &p354}
	p356.items = []parser{&p338, &p833, &p354, &p355}
	var p363 = sequenceParser{id: 363, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{378, 340, 341, 400, 392, 582, 575, 801}}
	var p339 = choiceParser{id: 339, commit: 66, name: "operand2", generalizations: []int{340, 341}}
	p339.options = []parser{&p338, &p356}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p357.items = []parser{&p833, &p14}
	p358.items = []parser{&p14, &p357}
	var p334 = choiceParser{id: 334, commit: 66, name: "binary-op2"}
	var p307 = sequenceParser{id: 307, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p305 = charParser{id: 305, chars: []rune{61}}
	var p306 = charParser{id: 306, chars: []rune{61}}
	p307.items = []parser{&p305, &p306}
	var p310 = sequenceParser{id: 310, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p308 = charParser{id: 308, chars: []rune{33}}
	var p309 = charParser{id: 309, chars: []rune{61}}
	p310.items = []parser{&p308, &p309}
	var p312 = sequenceParser{id: 312, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p311 = charParser{id: 311, chars: []rune{60}}
	p312.items = []parser{&p311}
	var p315 = sequenceParser{id: 315, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p313 = charParser{id: 313, chars: []rune{60}}
	var p314 = charParser{id: 314, chars: []rune{61}}
	p315.items = []parser{&p313, &p314}
	var p317 = sequenceParser{id: 317, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p316 = charParser{id: 316, chars: []rune{62}}
	p317.items = []parser{&p316}
	var p320 = sequenceParser{id: 320, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p318 = charParser{id: 318, chars: []rune{62}}
	var p319 = charParser{id: 319, chars: []rune{61}}
	p320.items = []parser{&p318, &p319}
	p334.options = []parser{&p307, &p310, &p312, &p315, &p317, &p320}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p833, &p14}
	p360.items = []parser{&p833, &p14, &p359}
	p361.items = []parser{&p358, &p833, &p334, &p360, &p833, &p339}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p362.items = []parser{&p833, &p361}
	p363.items = []parser{&p339, &p833, &p361, &p362}
	var p370 = sequenceParser{id: 370, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{378, 341, 400, 392, 582, 575, 801}}
	var p340 = choiceParser{id: 340, commit: 66, name: "operand3", generalizations: []int{341}}
	p340.options = []parser{&p339, &p363}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p365 = sequenceParser{id: 365, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p364.items = []parser{&p833, &p14}
	p365.items = []parser{&p14, &p364}
	var p335 = sequenceParser{id: 335, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p323 = sequenceParser{id: 323, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p321 = charParser{id: 321, chars: []rune{38}}
	var p322 = charParser{id: 322, chars: []rune{38}}
	p323.items = []parser{&p321, &p322}
	p335.items = []parser{&p323}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p366.items = []parser{&p833, &p14}
	p367.items = []parser{&p833, &p14, &p366}
	p368.items = []parser{&p365, &p833, &p335, &p367, &p833, &p340}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p369.items = []parser{&p833, &p368}
	p370.items = []parser{&p340, &p833, &p368, &p369}
	var p377 = sequenceParser{id: 377, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{378, 400, 392, 582, 575, 801}}
	var p341 = choiceParser{id: 341, commit: 66, name: "operand4"}
	p341.options = []parser{&p340, &p370}
	var p375 = sequenceParser{id: 375, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p371.items = []parser{&p833, &p14}
	p372.items = []parser{&p14, &p371}
	var p336 = sequenceParser{id: 336, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p326 = sequenceParser{id: 326, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p324 = charParser{id: 324, chars: []rune{124}}
	var p325 = charParser{id: 325, chars: []rune{124}}
	p326.items = []parser{&p324, &p325}
	p336.items = []parser{&p326}
	var p374 = sequenceParser{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p373.items = []parser{&p833, &p14}
	p374.items = []parser{&p833, &p14, &p373}
	p375.items = []parser{&p372, &p833, &p336, &p374, &p833, &p341}
	var p376 = sequenceParser{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p376.items = []parser{&p833, &p375}
	p377.items = []parser{&p341, &p833, &p375, &p376}
	p378.options = []parser{&p349, &p356, &p363, &p370, &p377}
	var p391 = sequenceParser{id: 391, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{400, 392, 582, 575, 801}}
	var p384 = sequenceParser{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p383 = sequenceParser{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p383.items = []parser{&p833, &p14}
	p384.items = []parser{&p833, &p14, &p383}
	var p380 = sequenceParser{id: 380, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p379 = charParser{id: 379, chars: []rune{63}}
	p380.items = []parser{&p379}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p385.items = []parser{&p833, &p14}
	p386.items = []parser{&p833, &p14, &p385}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p387.items = []parser{&p833, &p14}
	p388.items = []parser{&p833, &p14, &p387}
	var p382 = sequenceParser{id: 382, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p381 = charParser{id: 381, chars: []rune{58}}
	p382.items = []parser{&p381}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p389.items = []parser{&p833, &p14}
	p390.items = []parser{&p833, &p14, &p389}
	p391.items = []parser{&p400, &p384, &p833, &p380, &p386, &p833, &p400, &p388, &p833, &p382, &p390, &p833, &p400}
	var p399 = sequenceParser{id: 399, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{400, 582, 575, 801}}
	var p392 = choiceParser{id: 392, commit: 66, name: "chainingOperand"}
	p392.options = []parser{&p271, &p331, &p378, &p391}
	var p397 = sequenceParser{id: 397, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p394 = sequenceParser{id: 394, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p393.items = []parser{&p833, &p14}
	p394.items = []parser{&p14, &p393}
	var p329 = sequenceParser{id: 329, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p327 = charParser{id: 327, chars: []rune{45}}
	var p328 = charParser{id: 328, chars: []rune{62}}
	p329.items = []parser{&p327, &p328}
	var p396 = sequenceParser{id: 396, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p395 = sequenceParser{id: 395, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p395.items = []parser{&p833, &p14}
	p396.items = []parser{&p833, &p14, &p395}
	p397.items = []parser{&p394, &p833, &p329, &p396, &p833, &p392}
	var p398 = sequenceParser{id: 398, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p398.items = []parser{&p833, &p397}
	p399.items = []parser{&p392, &p833, &p397, &p398}
	p400.options = []parser{&p271, &p331, &p378, &p391, &p399}
	p184.items = []parser{&p183, &p833, &p400}
	p185.items = []parser{&p181, &p833, &p184}
	var p437 = sequenceParser{id: 437, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{801, 483, 547}}
	var p403 = sequenceParser{id: 403, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p401 = charParser{id: 401, chars: []rune{105}}
	var p402 = charParser{id: 402, chars: []rune{102}}
	p403.items = []parser{&p401, &p402}
	var p432 = sequenceParser{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p431 = sequenceParser{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p431.items = []parser{&p833, &p14}
	p432.items = []parser{&p833, &p14, &p431}
	var p434 = sequenceParser{id: 434, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p433 = sequenceParser{id: 433, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p433.items = []parser{&p833, &p14}
	p434.items = []parser{&p833, &p14, &p433}
	var p436 = sequenceParser{id: 436, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p420 = sequenceParser{id: 420, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p413 = sequenceParser{id: 413, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p412 = sequenceParser{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p412.items = []parser{&p833, &p14}
	p413.items = []parser{&p14, &p412}
	var p408 = sequenceParser{id: 408, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p404 = charParser{id: 404, chars: []rune{101}}
	var p405 = charParser{id: 405, chars: []rune{108}}
	var p406 = charParser{id: 406, chars: []rune{115}}
	var p407 = charParser{id: 407, chars: []rune{101}}
	p408.items = []parser{&p404, &p405, &p406, &p407}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p414.items = []parser{&p833, &p14}
	p415.items = []parser{&p833, &p14, &p414}
	var p411 = sequenceParser{id: 411, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p409 = charParser{id: 409, chars: []rune{105}}
	var p410 = charParser{id: 410, chars: []rune{102}}
	p411.items = []parser{&p409, &p410}
	var p417 = sequenceParser{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p416 = sequenceParser{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p416.items = []parser{&p833, &p14}
	p417.items = []parser{&p833, &p14, &p416}
	var p419 = sequenceParser{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p418 = sequenceParser{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p418.items = []parser{&p833, &p14}
	p419.items = []parser{&p833, &p14, &p418}
	p420.items = []parser{&p413, &p833, &p408, &p415, &p833, &p411, &p417, &p833, &p400, &p419, &p833, &p190}
	var p435 = sequenceParser{id: 435, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p435.items = []parser{&p833, &p420}
	p436.items = []parser{&p833, &p420, &p435}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p426.items = []parser{&p833, &p14}
	p427.items = []parser{&p14, &p426}
	var p425 = sequenceParser{id: 425, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p421 = charParser{id: 421, chars: []rune{101}}
	var p422 = charParser{id: 422, chars: []rune{108}}
	var p423 = charParser{id: 423, chars: []rune{115}}
	var p424 = charParser{id: 424, chars: []rune{101}}
	p425.items = []parser{&p421, &p422, &p423, &p424}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p428.items = []parser{&p833, &p14}
	p429.items = []parser{&p833, &p14, &p428}
	p430.items = []parser{&p427, &p833, &p425, &p429, &p833, &p190}
	p437.items = []parser{&p403, &p432, &p833, &p400, &p434, &p833, &p190, &p436, &p833, &p430}
	var p494 = sequenceParser{id: 494, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{483, 801, 547}}
	var p479 = sequenceParser{id: 479, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p473 = charParser{id: 473, chars: []rune{115}}
	var p474 = charParser{id: 474, chars: []rune{119}}
	var p475 = charParser{id: 475, chars: []rune{105}}
	var p476 = charParser{id: 476, chars: []rune{116}}
	var p477 = charParser{id: 477, chars: []rune{99}}
	var p478 = charParser{id: 478, chars: []rune{104}}
	p479.items = []parser{&p473, &p474, &p475, &p476, &p477, &p478}
	var p491 = sequenceParser{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p490 = sequenceParser{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p490.items = []parser{&p833, &p14}
	p491.items = []parser{&p833, &p14, &p490}
	var p493 = sequenceParser{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p492.items = []parser{&p833, &p14}
	p493.items = []parser{&p833, &p14, &p492}
	var p481 = sequenceParser{id: 481, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p480 = charParser{id: 480, chars: []rune{123}}
	p481.items = []parser{&p480}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p482 = choiceParser{id: 482, commit: 2}
	var p472 = sequenceParser{id: 472, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{482, 483}}
	var p467 = sequenceParser{id: 467, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p460 = sequenceParser{id: 460, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p456 = charParser{id: 456, chars: []rune{99}}
	var p457 = charParser{id: 457, chars: []rune{97}}
	var p458 = charParser{id: 458, chars: []rune{115}}
	var p459 = charParser{id: 459, chars: []rune{101}}
	p460.items = []parser{&p456, &p457, &p458, &p459}
	var p464 = sequenceParser{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p463 = sequenceParser{id: 463, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p463.items = []parser{&p833, &p14}
	p464.items = []parser{&p833, &p14, &p463}
	var p466 = sequenceParser{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p465 = sequenceParser{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p465.items = []parser{&p833, &p14}
	p466.items = []parser{&p833, &p14, &p465}
	var p462 = sequenceParser{id: 462, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p461 = charParser{id: 461, chars: []rune{58}}
	p462.items = []parser{&p461}
	p467.items = []parser{&p460, &p464, &p833, &p400, &p466, &p833, &p462}
	var p471 = sequenceParser{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p469 = sequenceParser{id: 469, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p468 = charParser{id: 468, chars: []rune{59}}
	p469.items = []parser{&p468}
	var p470 = sequenceParser{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p470.items = []parser{&p833, &p469}
	p471.items = []parser{&p833, &p469, &p470}
	p472.items = []parser{&p467, &p471, &p833, &p801}
	var p455 = sequenceParser{id: 455, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{482, 483, 546, 547}}
	var p450 = sequenceParser{id: 450, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p445 = sequenceParser{id: 445, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p438 = charParser{id: 438, chars: []rune{100}}
	var p439 = charParser{id: 439, chars: []rune{101}}
	var p440 = charParser{id: 440, chars: []rune{102}}
	var p441 = charParser{id: 441, chars: []rune{97}}
	var p442 = charParser{id: 442, chars: []rune{117}}
	var p443 = charParser{id: 443, chars: []rune{108}}
	var p444 = charParser{id: 444, chars: []rune{116}}
	p445.items = []parser{&p438, &p439, &p440, &p441, &p442, &p443, &p444}
	var p449 = sequenceParser{id: 449, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p448 = sequenceParser{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p448.items = []parser{&p833, &p14}
	p449.items = []parser{&p833, &p14, &p448}
	var p447 = sequenceParser{id: 447, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p446 = charParser{id: 446, chars: []rune{58}}
	p447.items = []parser{&p446}
	p450.items = []parser{&p445, &p449, &p833, &p447}
	var p454 = sequenceParser{id: 454, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p452 = sequenceParser{id: 452, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p451 = charParser{id: 451, chars: []rune{59}}
	p452.items = []parser{&p451}
	var p453 = sequenceParser{id: 453, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p453.items = []parser{&p833, &p452}
	p454.items = []parser{&p833, &p452, &p453}
	p455.items = []parser{&p450, &p454, &p833, &p801}
	p482.options = []parser{&p472, &p455}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p484 = sequenceParser{id: 484, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p483 = choiceParser{id: 483, commit: 2}
	p483.options = []parser{&p472, &p455, &p801}
	p484.items = []parser{&p815, &p833, &p483}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p485.items = []parser{&p833, &p484}
	p486.items = []parser{&p833, &p484, &p485}
	p487.items = []parser{&p482, &p486}
	var p489 = sequenceParser{id: 489, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p488 = charParser{id: 488, chars: []rune{125}}
	p489.items = []parser{&p488}
	p494.items = []parser{&p479, &p491, &p833, &p400, &p493, &p833, &p481, &p833, &p815, &p833, &p487, &p833, &p815, &p833, &p489}
	var p556 = sequenceParser{id: 556, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{547, 801}}
	var p543 = sequenceParser{id: 543, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p537 = charParser{id: 537, chars: []rune{115}}
	var p538 = charParser{id: 538, chars: []rune{101}}
	var p539 = charParser{id: 539, chars: []rune{108}}
	var p540 = charParser{id: 540, chars: []rune{101}}
	var p541 = charParser{id: 541, chars: []rune{99}}
	var p542 = charParser{id: 542, chars: []rune{116}}
	p543.items = []parser{&p537, &p538, &p539, &p540, &p541, &p542}
	var p555 = sequenceParser{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p554 = sequenceParser{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p554.items = []parser{&p833, &p14}
	p555.items = []parser{&p833, &p14, &p554}
	var p545 = sequenceParser{id: 545, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p544 = charParser{id: 544, chars: []rune{123}}
	p545.items = []parser{&p544}
	var p551 = sequenceParser{id: 551, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p546 = choiceParser{id: 546, commit: 2}
	var p536 = sequenceParser{id: 536, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{546, 547}}
	var p531 = sequenceParser{id: 531, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p524 = sequenceParser{id: 524, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p520 = charParser{id: 520, chars: []rune{99}}
	var p521 = charParser{id: 521, chars: []rune{97}}
	var p522 = charParser{id: 522, chars: []rune{115}}
	var p523 = charParser{id: 523, chars: []rune{101}}
	p524.items = []parser{&p520, &p521, &p522, &p523}
	var p528 = sequenceParser{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p527 = sequenceParser{id: 527, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p527.items = []parser{&p833, &p14}
	p528.items = []parser{&p833, &p14, &p527}
	var p519 = choiceParser{id: 519, commit: 66, name: "communication"}
	var p518 = sequenceParser{id: 518, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{519}}
	var p517 = sequenceParser{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p516 = sequenceParser{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p516.items = []parser{&p833, &p14}
	p517.items = []parser{&p833, &p14, &p516}
	p518.items = []parser{&p103, &p517, &p833, &p515}
	p519.options = []parser{&p504, &p515, &p518}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p529.items = []parser{&p833, &p14}
	p530.items = []parser{&p833, &p14, &p529}
	var p526 = sequenceParser{id: 526, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p525 = charParser{id: 525, chars: []rune{58}}
	p526.items = []parser{&p525}
	p531.items = []parser{&p524, &p528, &p833, &p519, &p530, &p833, &p526}
	var p535 = sequenceParser{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p533 = sequenceParser{id: 533, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p532 = charParser{id: 532, chars: []rune{59}}
	p533.items = []parser{&p532}
	var p534 = sequenceParser{id: 534, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p534.items = []parser{&p833, &p533}
	p535.items = []parser{&p833, &p533, &p534}
	p536.items = []parser{&p531, &p535, &p833, &p801}
	p546.options = []parser{&p536, &p455}
	var p550 = sequenceParser{id: 550, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p548 = sequenceParser{id: 548, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p547 = choiceParser{id: 547, commit: 2}
	p547.options = []parser{&p536, &p455, &p801}
	p548.items = []parser{&p815, &p833, &p547}
	var p549 = sequenceParser{id: 549, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p549.items = []parser{&p833, &p548}
	p550.items = []parser{&p833, &p548, &p549}
	p551.items = []parser{&p546, &p550}
	var p553 = sequenceParser{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p552 = charParser{id: 552, chars: []rune{125}}
	p553.items = []parser{&p552}
	p556.items = []parser{&p543, &p555, &p833, &p545, &p833, &p815, &p833, &p551, &p833, &p815, &p833, &p553}
	var p597 = sequenceParser{id: 597, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{801}}
	var p586 = sequenceParser{id: 586, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p583 = charParser{id: 583, chars: []rune{102}}
	var p584 = charParser{id: 584, chars: []rune{111}}
	var p585 = charParser{id: 585, chars: []rune{114}}
	p586.items = []parser{&p583, &p584, &p585}
	var p596 = choiceParser{id: 596, commit: 2}
	var p592 = sequenceParser{id: 592, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{596}}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p587.items = []parser{&p833, &p14}
	p588.items = []parser{&p14, &p587}
	var p582 = choiceParser{id: 582, commit: 66, name: "loop-expression"}
	var p581 = choiceParser{id: 581, commit: 64, name: "range-over-expression", generalizations: []int{582}}
	var p580 = sequenceParser{id: 580, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{581, 582}}
	var p577 = sequenceParser{id: 577, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p576 = sequenceParser{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p576.items = []parser{&p833, &p14}
	p577.items = []parser{&p833, &p14, &p576}
	var p574 = sequenceParser{id: 574, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p572 = charParser{id: 572, chars: []rune{105}}
	var p573 = charParser{id: 573, chars: []rune{110}}
	p574.items = []parser{&p572, &p573}
	var p579 = sequenceParser{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p578 = sequenceParser{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p578.items = []parser{&p833, &p14}
	p579.items = []parser{&p833, &p14, &p578}
	var p575 = choiceParser{id: 575, commit: 2}
	p575.options = []parser{&p400, &p225}
	p580.items = []parser{&p103, &p577, &p833, &p574, &p579, &p833, &p575}
	p581.options = []parser{&p580, &p225}
	p582.options = []parser{&p400, &p581}
	p589.items = []parser{&p588, &p833, &p582}
	var p591 = sequenceParser{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p590.items = []parser{&p833, &p14}
	p591.items = []parser{&p833, &p14, &p590}
	p592.items = []parser{&p589, &p591, &p833, &p190}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{596}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p593.items = []parser{&p833, &p14}
	p594.items = []parser{&p14, &p593}
	p595.items = []parser{&p594, &p833, &p190}
	p596.options = []parser{&p592, &p595}
	p597.items = []parser{&p586, &p833, &p596}
	var p745 = choiceParser{id: 745, commit: 66, name: "definition", generalizations: []int{801}}
	var p658 = sequenceParser{id: 658, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{745, 801}}
	var p654 = sequenceParser{id: 654, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p651 = charParser{id: 651, chars: []rune{108}}
	var p652 = charParser{id: 652, chars: []rune{101}}
	var p653 = charParser{id: 653, chars: []rune{116}}
	p654.items = []parser{&p651, &p652, &p653}
	var p657 = sequenceParser{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p656 = sequenceParser{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p656.items = []parser{&p833, &p14}
	p657.items = []parser{&p833, &p14, &p656}
	var p655 = choiceParser{id: 655, commit: 2}
	var p645 = sequenceParser{id: 645, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{655, 659, 660}}
	var p644 = sequenceParser{id: 644, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p641 = sequenceParser{id: 641, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p640 = sequenceParser{id: 640, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p639 = sequenceParser{id: 639, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p639.items = []parser{&p833, &p14}
	p640.items = []parser{&p14, &p639}
	var p638 = sequenceParser{id: 638, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p637 = charParser{id: 637, chars: []rune{61}}
	p638.items = []parser{&p637}
	p641.items = []parser{&p640, &p833, &p638}
	var p643 = sequenceParser{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p642.items = []parser{&p833, &p14}
	p643.items = []parser{&p833, &p14, &p642}
	p644.items = []parser{&p103, &p833, &p641, &p643, &p833, &p400}
	p645.items = []parser{&p644}
	var p650 = sequenceParser{id: 650, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{655, 659, 660}}
	var p647 = sequenceParser{id: 647, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p646 = charParser{id: 646, chars: []rune{126}}
	p647.items = []parser{&p646}
	var p649 = sequenceParser{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p648 = sequenceParser{id: 648, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p648.items = []parser{&p833, &p14}
	p649.items = []parser{&p833, &p14, &p648}
	p650.items = []parser{&p647, &p649, &p833, &p644}
	p655.options = []parser{&p645, &p650}
	p658.items = []parser{&p654, &p657, &p833, &p655}
	var p679 = sequenceParser{id: 679, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{745, 801}}
	var p672 = sequenceParser{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p669 = charParser{id: 669, chars: []rune{108}}
	var p670 = charParser{id: 670, chars: []rune{101}}
	var p671 = charParser{id: 671, chars: []rune{116}}
	p672.items = []parser{&p669, &p670, &p671}
	var p678 = sequenceParser{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p677 = sequenceParser{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p677.items = []parser{&p833, &p14}
	p678.items = []parser{&p833, &p14, &p677}
	var p674 = sequenceParser{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p673 = charParser{id: 673, chars: []rune{40}}
	p674.items = []parser{&p673}
	var p664 = sequenceParser{id: 664, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p659 = choiceParser{id: 659, commit: 2}
	p659.options = []parser{&p645, &p650}
	var p663 = sequenceParser{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p661 = sequenceParser{id: 661, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p660 = choiceParser{id: 660, commit: 2}
	p660.options = []parser{&p645, &p650}
	p661.items = []parser{&p113, &p833, &p660}
	var p662 = sequenceParser{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p662.items = []parser{&p833, &p661}
	p663.items = []parser{&p833, &p661, &p662}
	p664.items = []parser{&p659, &p663}
	var p676 = sequenceParser{id: 676, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p675 = charParser{id: 675, chars: []rune{41}}
	p676.items = []parser{&p675}
	p679.items = []parser{&p672, &p678, &p833, &p674, &p833, &p113, &p833, &p664, &p833, &p113, &p833, &p676}
	var p694 = sequenceParser{id: 694, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{745, 801}}
	var p683 = sequenceParser{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{108}}
	var p681 = charParser{id: 681, chars: []rune{101}}
	var p682 = charParser{id: 682, chars: []rune{116}}
	p683.items = []parser{&p680, &p681, &p682}
	var p691 = sequenceParser{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p690 = sequenceParser{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p690.items = []parser{&p833, &p14}
	p691.items = []parser{&p833, &p14, &p690}
	var p685 = sequenceParser{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p684 = charParser{id: 684, chars: []rune{126}}
	p685.items = []parser{&p684}
	var p693 = sequenceParser{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p692 = sequenceParser{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p692.items = []parser{&p833, &p14}
	p693.items = []parser{&p833, &p14, &p692}
	var p687 = sequenceParser{id: 687, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p686 = charParser{id: 686, chars: []rune{40}}
	p687.items = []parser{&p686}
	var p668 = sequenceParser{id: 668, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p665 = sequenceParser{id: 665, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p665.items = []parser{&p113, &p833, &p645}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p666.items = []parser{&p833, &p665}
	p667.items = []parser{&p833, &p665, &p666}
	p668.items = []parser{&p645, &p667}
	var p689 = sequenceParser{id: 689, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p688 = charParser{id: 688, chars: []rune{41}}
	p689.items = []parser{&p688}
	p694.items = []parser{&p683, &p691, &p833, &p685, &p693, &p833, &p687, &p833, &p113, &p833, &p668, &p833, &p113, &p833, &p689}
	var p710 = sequenceParser{id: 710, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{745, 801}}
	var p706 = sequenceParser{id: 706, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p704 = charParser{id: 704, chars: []rune{102}}
	var p705 = charParser{id: 705, chars: []rune{110}}
	p706.items = []parser{&p704, &p705}
	var p709 = sequenceParser{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p708 = sequenceParser{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p708.items = []parser{&p833, &p14}
	p709.items = []parser{&p833, &p14, &p708}
	var p707 = choiceParser{id: 707, commit: 2}
	var p698 = sequenceParser{id: 698, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{707, 715, 716}}
	var p697 = sequenceParser{id: 697, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p695 = sequenceParser{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p695.items = []parser{&p833, &p14}
	p696.items = []parser{&p833, &p14, &p695}
	p697.items = []parser{&p103, &p696, &p833, &p200}
	p698.items = []parser{&p697}
	var p703 = sequenceParser{id: 703, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{707, 715, 716}}
	var p700 = sequenceParser{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p699 = charParser{id: 699, chars: []rune{126}}
	p700.items = []parser{&p699}
	var p702 = sequenceParser{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p701 = sequenceParser{id: 701, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p701.items = []parser{&p833, &p14}
	p702.items = []parser{&p833, &p14, &p701}
	p703.items = []parser{&p700, &p702, &p833, &p697}
	p707.options = []parser{&p698, &p703}
	p710.items = []parser{&p706, &p709, &p833, &p707}
	var p730 = sequenceParser{id: 730, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{745, 801}}
	var p723 = sequenceParser{id: 723, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p721 = charParser{id: 721, chars: []rune{102}}
	var p722 = charParser{id: 722, chars: []rune{110}}
	p723.items = []parser{&p721, &p722}
	var p729 = sequenceParser{id: 729, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p728 = sequenceParser{id: 728, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p728.items = []parser{&p833, &p14}
	p729.items = []parser{&p833, &p14, &p728}
	var p725 = sequenceParser{id: 725, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p724 = charParser{id: 724, chars: []rune{40}}
	p725.items = []parser{&p724}
	var p720 = sequenceParser{id: 720, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p715 = choiceParser{id: 715, commit: 2}
	p715.options = []parser{&p698, &p703}
	var p719 = sequenceParser{id: 719, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p717 = sequenceParser{id: 717, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p716 = choiceParser{id: 716, commit: 2}
	p716.options = []parser{&p698, &p703}
	p717.items = []parser{&p113, &p833, &p716}
	var p718 = sequenceParser{id: 718, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p718.items = []parser{&p833, &p717}
	p719.items = []parser{&p833, &p717, &p718}
	p720.items = []parser{&p715, &p719}
	var p727 = sequenceParser{id: 727, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p726 = charParser{id: 726, chars: []rune{41}}
	p727.items = []parser{&p726}
	p730.items = []parser{&p723, &p729, &p833, &p725, &p833, &p113, &p833, &p720, &p833, &p113, &p833, &p727}
	var p744 = sequenceParser{id: 744, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{745, 801}}
	var p733 = sequenceParser{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p731 = charParser{id: 731, chars: []rune{102}}
	var p732 = charParser{id: 732, chars: []rune{110}}
	p733.items = []parser{&p731, &p732}
	var p741 = sequenceParser{id: 741, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p740 = sequenceParser{id: 740, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p740.items = []parser{&p833, &p14}
	p741.items = []parser{&p833, &p14, &p740}
	var p735 = sequenceParser{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p734 = charParser{id: 734, chars: []rune{126}}
	p735.items = []parser{&p734}
	var p743 = sequenceParser{id: 743, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p742 = sequenceParser{id: 742, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p742.items = []parser{&p833, &p14}
	p743.items = []parser{&p833, &p14, &p742}
	var p737 = sequenceParser{id: 737, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p736 = charParser{id: 736, chars: []rune{40}}
	p737.items = []parser{&p736}
	var p714 = sequenceParser{id: 714, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p711 = sequenceParser{id: 711, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p711.items = []parser{&p113, &p833, &p698}
	var p712 = sequenceParser{id: 712, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p712.items = []parser{&p833, &p711}
	p713.items = []parser{&p833, &p711, &p712}
	p714.items = []parser{&p698, &p713}
	var p739 = sequenceParser{id: 739, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p738 = charParser{id: 738, chars: []rune{41}}
	p739.items = []parser{&p738}
	p744.items = []parser{&p733, &p741, &p833, &p735, &p743, &p833, &p737, &p833, &p113, &p833, &p714, &p833, &p113, &p833, &p739}
	p745.options = []parser{&p658, &p679, &p694, &p710, &p730, &p744}
	var p780 = choiceParser{id: 780, commit: 64, name: "use", generalizations: []int{801}}
	var p768 = sequenceParser{id: 768, commit: 66, name: "use-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{780, 801}}
	var p765 = sequenceParser{id: 765, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p762 = charParser{id: 762, chars: []rune{117}}
	var p763 = charParser{id: 763, chars: []rune{115}}
	var p764 = charParser{id: 764, chars: []rune{101}}
	p765.items = []parser{&p762, &p763, &p764}
	var p767 = sequenceParser{id: 767, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p766 = sequenceParser{id: 766, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p766.items = []parser{&p833, &p14}
	p767.items = []parser{&p833, &p14, &p766}
	var p757 = choiceParser{id: 757, commit: 64, name: "use-fact"}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{757}}
	var p748 = choiceParser{id: 748, commit: 2}
	var p747 = sequenceParser{id: 747, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{748}}
	var p746 = charParser{id: 746, chars: []rune{46}}
	p747.items = []parser{&p746}
	p748.options = []parser{&p103, &p747}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p751 = sequenceParser{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p751.items = []parser{&p833, &p14}
	p752.items = []parser{&p14, &p751}
	var p750 = sequenceParser{id: 750, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p749 = charParser{id: 749, chars: []rune{61}}
	p750.items = []parser{&p749}
	p753.items = []parser{&p752, &p833, &p750}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p754.items = []parser{&p833, &p14}
	p755.items = []parser{&p833, &p14, &p754}
	p756.items = []parser{&p748, &p833, &p753, &p755, &p833, &p86}
	p757.options = []parser{&p86, &p756}
	p768.items = []parser{&p765, &p767, &p833, &p757}
	var p779 = sequenceParser{id: 779, commit: 66, name: "use-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{780, 801}}
	var p772 = sequenceParser{id: 772, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p769 = charParser{id: 769, chars: []rune{117}}
	var p770 = charParser{id: 770, chars: []rune{115}}
	var p771 = charParser{id: 771, chars: []rune{101}}
	p772.items = []parser{&p769, &p770, &p771}
	var p778 = sequenceParser{id: 778, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p777 = sequenceParser{id: 777, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p777.items = []parser{&p833, &p14}
	p778.items = []parser{&p833, &p14, &p777}
	var p774 = sequenceParser{id: 774, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p773 = charParser{id: 773, chars: []rune{40}}
	p774.items = []parser{&p773}
	var p761 = sequenceParser{id: 761, commit: 66, name: "use-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p760 = sequenceParser{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p758 = sequenceParser{id: 758, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p758.items = []parser{&p113, &p833, &p757}
	var p759 = sequenceParser{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p759.items = []parser{&p833, &p758}
	p760.items = []parser{&p833, &p758, &p759}
	p761.items = []parser{&p757, &p760}
	var p776 = sequenceParser{id: 776, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p775 = charParser{id: 775, chars: []rune{41}}
	p776.items = []parser{&p775}
	p779.items = []parser{&p772, &p778, &p833, &p774, &p833, &p113, &p833, &p761, &p833, &p113, &p833, &p776}
	p780.options = []parser{&p768, &p779}
	var p790 = sequenceParser{id: 790, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{801}}
	var p787 = sequenceParser{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p781 = charParser{id: 781, chars: []rune{101}}
	var p782 = charParser{id: 782, chars: []rune{120}}
	var p783 = charParser{id: 783, chars: []rune{112}}
	var p784 = charParser{id: 784, chars: []rune{111}}
	var p785 = charParser{id: 785, chars: []rune{114}}
	var p786 = charParser{id: 786, chars: []rune{116}}
	p787.items = []parser{&p781, &p782, &p783, &p784, &p785, &p786}
	var p789 = sequenceParser{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p788 = sequenceParser{id: 788, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p788.items = []parser{&p833, &p14}
	p789.items = []parser{&p833, &p14, &p788}
	p790.items = []parser{&p787, &p789, &p833, &p745}
	var p810 = sequenceParser{id: 810, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{801}}
	var p803 = sequenceParser{id: 803, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p802 = charParser{id: 802, chars: []rune{40}}
	p803.items = []parser{&p802}
	var p807 = sequenceParser{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p806 = sequenceParser{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p806.items = []parser{&p833, &p14}
	p807.items = []parser{&p833, &p14, &p806}
	var p809 = sequenceParser{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p808 = sequenceParser{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p808.items = []parser{&p833, &p14}
	p809.items = []parser{&p833, &p14, &p808}
	var p805 = sequenceParser{id: 805, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p804 = charParser{id: 804, chars: []rune{41}}
	p805.items = []parser{&p804}
	p810.items = []parser{&p803, &p807, &p833, &p801, &p809, &p833, &p805}
	p801.options = []parser{&p185, &p437, &p494, &p556, &p597, &p745, &p780, &p790, &p810, &p791}
	var p818 = sequenceParser{id: 818, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p816 = sequenceParser{id: 816, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p816.items = []parser{&p815, &p833, &p801}
	var p817 = sequenceParser{id: 817, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p817.items = []parser{&p833, &p816}
	p818.items = []parser{&p833, &p816, &p817}
	p819.items = []parser{&p801, &p818}
	p834.items = []parser{&p830, &p833, &p815, &p833, &p819, &p833, &p815}
	p835.items = []parser{&p833, &p834, &p833}
	var b835 = sequenceBuilder{id: 835, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b833 = choiceBuilder{id: 833, commit: 2}
	var b831 = choiceBuilder{id: 831, commit: 70}
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
	b831.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b832 = sequenceBuilder{id: 832, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b39.items = []builder{&b14, &b833, &b38}
	var b40 = sequenceBuilder{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b40.items = []builder{&b833, &b39}
	b41.items = []builder{&b833, &b39, &b40}
	b42.items = []builder{&b38, &b41}
	b832.items = []builder{&b42}
	b833.options = []builder{&b831, &b832}
	var b834 = sequenceBuilder{id: 834, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b830 = sequenceBuilder{id: 830, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b827 = sequenceBuilder{id: 827, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b825 = charBuilder{}
	var b826 = charBuilder{}
	b827.items = []builder{&b825, &b826}
	var b824 = sequenceBuilder{id: 824, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b823 = sequenceBuilder{id: 823, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b821 = sequenceBuilder{id: 821, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b820 = charBuilder{}
	b821.items = []builder{&b820}
	var b822 = sequenceBuilder{id: 822, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b822.items = []builder{&b833, &b821}
	b823.items = []builder{&b821, &b822}
	b824.items = []builder{&b823}
	var b829 = sequenceBuilder{id: 829, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b828 = charBuilder{}
	b829.items = []builder{&b828}
	b830.items = []builder{&b827, &b833, &b824, &b833, &b829}
	var b815 = sequenceBuilder{id: 815, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b813 = choiceBuilder{id: 813, commit: 2}
	var b812 = sequenceBuilder{id: 812, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b811 = charBuilder{}
	b812.items = []builder{&b811}
	b813.options = []builder{&b812, &b14}
	var b814 = sequenceBuilder{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b814.items = []builder{&b833, &b813}
	b815.items = []builder{&b813, &b814}
	var b819 = sequenceBuilder{id: 819, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b801 = choiceBuilder{id: 801, commit: 66}
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
	b182.items = []builder{&b833, &b14}
	b183.items = []builder{&b14, &b182}
	var b400 = choiceBuilder{id: 400, commit: 66}
	var b271 = choiceBuilder{id: 271, commit: 66}
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
	var b515 = sequenceBuilder{id: 515, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b512 = sequenceBuilder{id: 512, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b505 = charBuilder{}
	var b506 = charBuilder{}
	var b507 = charBuilder{}
	var b508 = charBuilder{}
	var b509 = charBuilder{}
	var b510 = charBuilder{}
	var b511 = charBuilder{}
	b512.items = []builder{&b505, &b506, &b507, &b508, &b509, &b510, &b511}
	var b514 = sequenceBuilder{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b513 = sequenceBuilder{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b513.items = []builder{&b833, &b14}
	b514.items = []builder{&b833, &b14, &b513}
	b515.items = []builder{&b512, &b514, &b833, &b271}
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
	b112.items = []builder{&b833, &b111}
	b113.items = []builder{&b111, &b112}
	var b118 = sequenceBuilder{id: 118, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b114 = choiceBuilder{id: 114, commit: 66}
	var b108 = sequenceBuilder{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b107 = sequenceBuilder{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b104 = charBuilder{}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	b107.items = []builder{&b104, &b105, &b106}
	b108.items = []builder{&b271, &b833, &b107}
	b114.options = []builder{&b400, &b108}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b115 = sequenceBuilder{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b115.items = []builder{&b113, &b833, &b114}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b116.items = []builder{&b833, &b115}
	b117.items = []builder{&b833, &b115, &b116}
	b118.items = []builder{&b114, &b117}
	var b122 = sequenceBuilder{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b121 = charBuilder{}
	b122.items = []builder{&b121}
	b123.items = []builder{&b120, &b833, &b113, &b833, &b118, &b833, &b113, &b833, &b122}
	b124.items = []builder{&b123}
	var b129 = sequenceBuilder{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b126 = sequenceBuilder{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b125 = charBuilder{}
	b126.items = []builder{&b125}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b127 = sequenceBuilder{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b127.items = []builder{&b833, &b14}
	b128.items = []builder{&b833, &b14, &b127}
	b129.items = []builder{&b126, &b128, &b833, &b123}
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
	b134.items = []builder{&b833, &b14}
	b135.items = []builder{&b833, &b14, &b134}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b136 = sequenceBuilder{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b136.items = []builder{&b833, &b14}
	b137.items = []builder{&b833, &b14, &b136}
	var b133 = sequenceBuilder{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b132 = charBuilder{}
	b133.items = []builder{&b132}
	b138.items = []builder{&b131, &b135, &b833, &b400, &b137, &b833, &b133}
	b139.options = []builder{&b103, &b86, &b138}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b142 = sequenceBuilder{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b142.items = []builder{&b833, &b14}
	b143.items = []builder{&b833, &b14, &b142}
	var b141 = sequenceBuilder{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b140 = charBuilder{}
	b141.items = []builder{&b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b833, &b14}
	b145.items = []builder{&b833, &b14, &b144}
	b146.items = []builder{&b139, &b143, &b833, &b141, &b145, &b833, &b400}
	b147.options = []builder{&b146, &b108}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b149 = sequenceBuilder{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b148 = choiceBuilder{id: 148, commit: 2}
	b148.options = []builder{&b146, &b108}
	b149.items = []builder{&b113, &b833, &b148}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b150.items = []builder{&b833, &b149}
	b151.items = []builder{&b833, &b149, &b150}
	b152.items = []builder{&b147, &b151}
	var b156 = sequenceBuilder{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b155 = charBuilder{}
	b156.items = []builder{&b155}
	b157.items = []builder{&b154, &b833, &b113, &b833, &b152, &b833, &b113, &b833, &b156}
	b158.items = []builder{&b157}
	var b163 = sequenceBuilder{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b160 = sequenceBuilder{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b159 = charBuilder{}
	b160.items = []builder{&b159}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b161 = sequenceBuilder{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b161.items = []builder{&b833, &b14}
	b162.items = []builder{&b833, &b14, &b161}
	b163.items = []builder{&b160, &b162, &b833, &b157}
	var b206 = sequenceBuilder{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b203 = sequenceBuilder{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b201 = charBuilder{}
	var b202 = charBuilder{}
	b203.items = []builder{&b201, &b202}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b204 = sequenceBuilder{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b204.items = []builder{&b833, &b14}
	b205.items = []builder{&b833, &b14, &b204}
	var b200 = sequenceBuilder{id: 200, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b192 = sequenceBuilder{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b191 = charBuilder{}
	b192.items = []builder{&b191}
	var b194 = choiceBuilder{id: 194, commit: 2}
	var b167 = sequenceBuilder{id: 167, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b164.items = []builder{&b113, &b833, &b103}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b165.items = []builder{&b833, &b164}
	b166.items = []builder{&b833, &b164, &b165}
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
	b172.items = []builder{&b833, &b14}
	b173.items = []builder{&b833, &b14, &b172}
	b174.items = []builder{&b171, &b173, &b833, &b103}
	b193.items = []builder{&b167, &b833, &b113, &b833, &b174}
	b194.options = []builder{&b167, &b193, &b174}
	var b196 = sequenceBuilder{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b195 = charBuilder{}
	b196.items = []builder{&b195}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b198 = sequenceBuilder{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b198.items = []builder{&b833, &b14}
	b199.items = []builder{&b833, &b14, &b198}
	var b197 = choiceBuilder{id: 197, commit: 2}
	var b791 = choiceBuilder{id: 791, commit: 66}
	var b504 = sequenceBuilder{id: 504, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b499 = sequenceBuilder{id: 499, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b495 = charBuilder{}
	var b496 = charBuilder{}
	var b497 = charBuilder{}
	var b498 = charBuilder{}
	b499.items = []builder{&b495, &b496, &b497, &b498}
	var b501 = sequenceBuilder{id: 501, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b500 = sequenceBuilder{id: 500, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b500.items = []builder{&b833, &b14}
	b501.items = []builder{&b833, &b14, &b500}
	var b503 = sequenceBuilder{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b502 = sequenceBuilder{id: 502, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b502.items = []builder{&b833, &b14}
	b503.items = []builder{&b833, &b14, &b502}
	b504.items = []builder{&b499, &b501, &b833, &b271, &b503, &b833, &b271}
	var b562 = sequenceBuilder{id: 562, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b559 = sequenceBuilder{id: 559, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b557 = charBuilder{}
	var b558 = charBuilder{}
	b559.items = []builder{&b557, &b558}
	var b561 = sequenceBuilder{id: 561, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b560 = sequenceBuilder{id: 560, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b560.items = []builder{&b833, &b14}
	b561.items = []builder{&b833, &b14, &b560}
	var b261 = sequenceBuilder{id: 261, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b258 = sequenceBuilder{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b257 = charBuilder{}
	b258.items = []builder{&b257}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	b261.items = []builder{&b271, &b833, &b258, &b833, &b113, &b833, &b118, &b833, &b113, &b833, &b260}
	b562.items = []builder{&b559, &b561, &b833, &b261}
	var b571 = sequenceBuilder{id: 571, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b568 = sequenceBuilder{id: 568, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b563 = charBuilder{}
	var b564 = charBuilder{}
	var b565 = charBuilder{}
	var b566 = charBuilder{}
	var b567 = charBuilder{}
	b568.items = []builder{&b563, &b564, &b565, &b566, &b567}
	var b570 = sequenceBuilder{id: 570, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b569 = sequenceBuilder{id: 569, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b569.items = []builder{&b833, &b14}
	b570.items = []builder{&b833, &b14, &b569}
	b571.items = []builder{&b568, &b570, &b833, &b261}
	var b636 = choiceBuilder{id: 636, commit: 64, name: "assignment"}
	var b616 = sequenceBuilder{id: 616, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b613 = sequenceBuilder{id: 613, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b610 = charBuilder{}
	var b611 = charBuilder{}
	var b612 = charBuilder{}
	b613.items = []builder{&b610, &b611, &b612}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b614.items = []builder{&b833, &b14}
	b615.items = []builder{&b833, &b14, &b614}
	var b605 = sequenceBuilder{id: 605, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b602 = sequenceBuilder{id: 602, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b601 = sequenceBuilder{id: 601, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b600 = sequenceBuilder{id: 600, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b600.items = []builder{&b833, &b14}
	b601.items = []builder{&b14, &b600}
	var b599 = sequenceBuilder{id: 599, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b598 = charBuilder{}
	b599.items = []builder{&b598}
	b602.items = []builder{&b601, &b833, &b599}
	var b604 = sequenceBuilder{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b603 = sequenceBuilder{id: 603, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b603.items = []builder{&b833, &b14}
	b604.items = []builder{&b833, &b14, &b603}
	b605.items = []builder{&b271, &b833, &b602, &b604, &b833, &b400}
	b616.items = []builder{&b613, &b615, &b833, &b605}
	var b623 = sequenceBuilder{id: 623, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b620 = sequenceBuilder{id: 620, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b619 = sequenceBuilder{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b619.items = []builder{&b833, &b14}
	b620.items = []builder{&b833, &b14, &b619}
	var b618 = sequenceBuilder{id: 618, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b617 = charBuilder{}
	b618.items = []builder{&b617}
	var b622 = sequenceBuilder{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b621 = sequenceBuilder{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b621.items = []builder{&b833, &b14}
	b622.items = []builder{&b833, &b14, &b621}
	b623.items = []builder{&b271, &b620, &b833, &b618, &b622, &b833, &b400}
	var b635 = sequenceBuilder{id: 635, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b627 = sequenceBuilder{id: 627, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b624 = charBuilder{}
	var b625 = charBuilder{}
	var b626 = charBuilder{}
	b627.items = []builder{&b624, &b625, &b626}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b633 = sequenceBuilder{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b633.items = []builder{&b833, &b14}
	b634.items = []builder{&b833, &b14, &b633}
	var b629 = sequenceBuilder{id: 629, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b628 = charBuilder{}
	b629.items = []builder{&b628}
	var b630 = sequenceBuilder{id: 630, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b609 = sequenceBuilder{id: 609, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b608 = sequenceBuilder{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b606 = sequenceBuilder{id: 606, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b606.items = []builder{&b113, &b833, &b605}
	var b607 = sequenceBuilder{id: 607, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b607.items = []builder{&b833, &b606}
	b608.items = []builder{&b833, &b606, &b607}
	b609.items = []builder{&b605, &b608}
	b630.items = []builder{&b113, &b833, &b609}
	var b632 = sequenceBuilder{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b631 = charBuilder{}
	b632.items = []builder{&b631}
	b635.items = []builder{&b627, &b634, &b833, &b629, &b833, &b630, &b833, &b113, &b833, &b632}
	b636.options = []builder{&b616, &b623, &b635}
	var b800 = sequenceBuilder{id: 800, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b793 = sequenceBuilder{id: 793, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b792 = charBuilder{}
	b793.items = []builder{&b792}
	var b797 = sequenceBuilder{id: 797, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b796 = sequenceBuilder{id: 796, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b796.items = []builder{&b833, &b14}
	b797.items = []builder{&b833, &b14, &b796}
	var b799 = sequenceBuilder{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b798 = sequenceBuilder{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b798.items = []builder{&b833, &b14}
	b799.items = []builder{&b833, &b14, &b798}
	var b795 = sequenceBuilder{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b794 = charBuilder{}
	b795.items = []builder{&b794}
	b800.items = []builder{&b793, &b797, &b833, &b791, &b799, &b833, &b795}
	b791.options = []builder{&b504, &b562, &b571, &b636, &b800, &b400}
	var b190 = sequenceBuilder{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b187 = sequenceBuilder{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b186 = charBuilder{}
	b187.items = []builder{&b186}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	b190.items = []builder{&b187, &b833, &b815, &b833, &b819, &b833, &b815, &b833, &b189}
	b197.options = []builder{&b791, &b190}
	b200.items = []builder{&b192, &b833, &b113, &b833, &b194, &b833, &b113, &b833, &b196, &b199, &b833, &b197}
	b206.items = []builder{&b203, &b205, &b833, &b200}
	var b216 = sequenceBuilder{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b209 = sequenceBuilder{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b207 = charBuilder{}
	var b208 = charBuilder{}
	b209.items = []builder{&b207, &b208}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b212 = sequenceBuilder{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b212.items = []builder{&b833, &b14}
	b213.items = []builder{&b833, &b14, &b212}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b210 = charBuilder{}
	b211.items = []builder{&b210}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b214.items = []builder{&b833, &b14}
	b215.items = []builder{&b833, &b14, &b214}
	b216.items = []builder{&b209, &b213, &b833, &b211, &b215, &b833, &b200}
	var b256 = sequenceBuilder{id: 256, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b255 = sequenceBuilder{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b254 = sequenceBuilder{id: 254, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b254.items = []builder{&b833, &b14}
	b255.items = []builder{&b833, &b14, &b254}
	var b253 = sequenceBuilder{id: 253, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b249 = choiceBuilder{id: 249, commit: 66}
	var b230 = sequenceBuilder{id: 230, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b227 = sequenceBuilder{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b226 = charBuilder{}
	b227.items = []builder{&b226}
	var b229 = sequenceBuilder{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b228 = sequenceBuilder{id: 228, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b228.items = []builder{&b833, &b14}
	b229.items = []builder{&b833, &b14, &b228}
	b230.items = []builder{&b227, &b229, &b833, &b103}
	var b239 = sequenceBuilder{id: 239, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b232 = sequenceBuilder{id: 232, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b231 = charBuilder{}
	b232.items = []builder{&b231}
	var b236 = sequenceBuilder{id: 236, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b235 = sequenceBuilder{id: 235, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b235.items = []builder{&b833, &b14}
	b236.items = []builder{&b833, &b14, &b235}
	var b238 = sequenceBuilder{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b237 = sequenceBuilder{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b237.items = []builder{&b833, &b14}
	b238.items = []builder{&b833, &b14, &b237}
	var b234 = sequenceBuilder{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b233 = charBuilder{}
	b234.items = []builder{&b233}
	b239.items = []builder{&b232, &b236, &b833, &b400, &b238, &b833, &b234}
	var b248 = sequenceBuilder{id: 248, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b241 = sequenceBuilder{id: 241, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b240 = charBuilder{}
	b241.items = []builder{&b240}
	var b245 = sequenceBuilder{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b244 = sequenceBuilder{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b244.items = []builder{&b833, &b14}
	b245.items = []builder{&b833, &b14, &b244}
	var b225 = sequenceBuilder{id: 225, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b217.items = []builder{&b400}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b221 = sequenceBuilder{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b221.items = []builder{&b833, &b14}
	b222.items = []builder{&b833, &b14, &b221}
	var b220 = sequenceBuilder{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b219 = charBuilder{}
	b220.items = []builder{&b219}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b833, &b14}
	b224.items = []builder{&b833, &b14, &b223}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b400}
	b225.items = []builder{&b217, &b222, &b833, &b220, &b224, &b833, &b218}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b246 = sequenceBuilder{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b246.items = []builder{&b833, &b14}
	b247.items = []builder{&b833, &b14, &b246}
	var b243 = sequenceBuilder{id: 243, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b242 = charBuilder{}
	b243.items = []builder{&b242}
	b248.items = []builder{&b241, &b245, &b833, &b225, &b247, &b833, &b243}
	b249.options = []builder{&b230, &b239, &b248}
	var b252 = sequenceBuilder{id: 252, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b251 = sequenceBuilder{id: 251, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b250.items = []builder{&b833, &b14}
	b251.items = []builder{&b14, &b250}
	b252.items = []builder{&b251, &b833, &b249}
	b253.items = []builder{&b249, &b833, &b252}
	b256.items = []builder{&b271, &b255, &b833, &b253}
	var b270 = sequenceBuilder{id: 270, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b263 = sequenceBuilder{id: 263, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b262 = charBuilder{}
	b263.items = []builder{&b262}
	var b267 = sequenceBuilder{id: 267, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b266 = sequenceBuilder{id: 266, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b266.items = []builder{&b833, &b14}
	b267.items = []builder{&b833, &b14, &b266}
	var b269 = sequenceBuilder{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b268 = sequenceBuilder{id: 268, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b268.items = []builder{&b833, &b14}
	b269.items = []builder{&b833, &b14, &b268}
	var b265 = sequenceBuilder{id: 265, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b264 = charBuilder{}
	b265.items = []builder{&b264}
	b270.items = []builder{&b263, &b267, &b833, &b400, &b269, &b833, &b265}
	b271.options = []builder{&b60, &b73, &b86, &b98, &b515, &b103, &b124, &b129, &b158, &b163, &b206, &b216, &b256, &b261, &b270}
	var b331 = sequenceBuilder{id: 331, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b330 = choiceBuilder{id: 330, commit: 66}
	var b290 = sequenceBuilder{id: 290, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b289 = charBuilder{}
	b290.items = []builder{&b289}
	var b292 = sequenceBuilder{id: 292, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b291 = charBuilder{}
	b292.items = []builder{&b291}
	var b273 = sequenceBuilder{id: 273, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b272 = charBuilder{}
	b273.items = []builder{&b272}
	var b304 = sequenceBuilder{id: 304, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b303 = charBuilder{}
	b304.items = []builder{&b303}
	b330.options = []builder{&b290, &b292, &b273, &b304}
	b331.items = []builder{&b330, &b833, &b271}
	var b378 = choiceBuilder{id: 378, commit: 66}
	var b349 = sequenceBuilder{id: 349, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b337 = choiceBuilder{id: 337, commit: 66}
	b337.options = []builder{&b271, &b331}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b344 = sequenceBuilder{id: 344, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b343.items = []builder{&b833, &b14}
	b344.items = []builder{&b14, &b343}
	var b332 = choiceBuilder{id: 332, commit: 66}
	var b275 = sequenceBuilder{id: 275, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b274 = charBuilder{}
	b275.items = []builder{&b274}
	var b282 = sequenceBuilder{id: 282, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b280 = charBuilder{}
	var b281 = charBuilder{}
	b282.items = []builder{&b280, &b281}
	var b285 = sequenceBuilder{id: 285, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b283 = charBuilder{}
	var b284 = charBuilder{}
	b285.items = []builder{&b283, &b284}
	var b288 = sequenceBuilder{id: 288, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b286 = charBuilder{}
	var b287 = charBuilder{}
	b288.items = []builder{&b286, &b287}
	var b294 = sequenceBuilder{id: 294, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b293 = charBuilder{}
	b294.items = []builder{&b293}
	var b296 = sequenceBuilder{id: 296, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b295 = charBuilder{}
	b296.items = []builder{&b295}
	var b298 = sequenceBuilder{id: 298, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b297 = charBuilder{}
	b298.items = []builder{&b297}
	b332.options = []builder{&b275, &b282, &b285, &b288, &b294, &b296, &b298}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b345 = sequenceBuilder{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b345.items = []builder{&b833, &b14}
	b346.items = []builder{&b833, &b14, &b345}
	b347.items = []builder{&b344, &b833, &b332, &b346, &b833, &b337}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b348.items = []builder{&b833, &b347}
	b349.items = []builder{&b337, &b833, &b347, &b348}
	var b356 = sequenceBuilder{id: 356, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b338 = choiceBuilder{id: 338, commit: 66}
	b338.options = []builder{&b337, &b349}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b351 = sequenceBuilder{id: 351, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b350.items = []builder{&b833, &b14}
	b351.items = []builder{&b14, &b350}
	var b333 = choiceBuilder{id: 333, commit: 66}
	var b277 = sequenceBuilder{id: 277, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b276 = charBuilder{}
	b277.items = []builder{&b276}
	var b279 = sequenceBuilder{id: 279, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b278 = charBuilder{}
	b279.items = []builder{&b278}
	var b300 = sequenceBuilder{id: 300, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b299 = charBuilder{}
	b300.items = []builder{&b299}
	var b302 = sequenceBuilder{id: 302, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b301 = charBuilder{}
	b302.items = []builder{&b301}
	b333.options = []builder{&b277, &b279, &b300, &b302}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b352.items = []builder{&b833, &b14}
	b353.items = []builder{&b833, &b14, &b352}
	b354.items = []builder{&b351, &b833, &b333, &b353, &b833, &b338}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b355.items = []builder{&b833, &b354}
	b356.items = []builder{&b338, &b833, &b354, &b355}
	var b363 = sequenceBuilder{id: 363, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b339 = choiceBuilder{id: 339, commit: 66}
	b339.options = []builder{&b338, &b356}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b357.items = []builder{&b833, &b14}
	b358.items = []builder{&b14, &b357}
	var b334 = choiceBuilder{id: 334, commit: 66}
	var b307 = sequenceBuilder{id: 307, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b305 = charBuilder{}
	var b306 = charBuilder{}
	b307.items = []builder{&b305, &b306}
	var b310 = sequenceBuilder{id: 310, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b308 = charBuilder{}
	var b309 = charBuilder{}
	b310.items = []builder{&b308, &b309}
	var b312 = sequenceBuilder{id: 312, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b311 = charBuilder{}
	b312.items = []builder{&b311}
	var b315 = sequenceBuilder{id: 315, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b313 = charBuilder{}
	var b314 = charBuilder{}
	b315.items = []builder{&b313, &b314}
	var b317 = sequenceBuilder{id: 317, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b316 = charBuilder{}
	b317.items = []builder{&b316}
	var b320 = sequenceBuilder{id: 320, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b318 = charBuilder{}
	var b319 = charBuilder{}
	b320.items = []builder{&b318, &b319}
	b334.options = []builder{&b307, &b310, &b312, &b315, &b317, &b320}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b833, &b14}
	b360.items = []builder{&b833, &b14, &b359}
	b361.items = []builder{&b358, &b833, &b334, &b360, &b833, &b339}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b362.items = []builder{&b833, &b361}
	b363.items = []builder{&b339, &b833, &b361, &b362}
	var b370 = sequenceBuilder{id: 370, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b340 = choiceBuilder{id: 340, commit: 66}
	b340.options = []builder{&b339, &b363}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b365 = sequenceBuilder{id: 365, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b364.items = []builder{&b833, &b14}
	b365.items = []builder{&b14, &b364}
	var b335 = sequenceBuilder{id: 335, commit: 66, ranges: [][]int{{1, 1}}}
	var b323 = sequenceBuilder{id: 323, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b321 = charBuilder{}
	var b322 = charBuilder{}
	b323.items = []builder{&b321, &b322}
	b335.items = []builder{&b323}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b366.items = []builder{&b833, &b14}
	b367.items = []builder{&b833, &b14, &b366}
	b368.items = []builder{&b365, &b833, &b335, &b367, &b833, &b340}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b369.items = []builder{&b833, &b368}
	b370.items = []builder{&b340, &b833, &b368, &b369}
	var b377 = sequenceBuilder{id: 377, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b341 = choiceBuilder{id: 341, commit: 66}
	b341.options = []builder{&b340, &b370}
	var b375 = sequenceBuilder{id: 375, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b371.items = []builder{&b833, &b14}
	b372.items = []builder{&b14, &b371}
	var b336 = sequenceBuilder{id: 336, commit: 66, ranges: [][]int{{1, 1}}}
	var b326 = sequenceBuilder{id: 326, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b324 = charBuilder{}
	var b325 = charBuilder{}
	b326.items = []builder{&b324, &b325}
	b336.items = []builder{&b326}
	var b374 = sequenceBuilder{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b373.items = []builder{&b833, &b14}
	b374.items = []builder{&b833, &b14, &b373}
	b375.items = []builder{&b372, &b833, &b336, &b374, &b833, &b341}
	var b376 = sequenceBuilder{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b376.items = []builder{&b833, &b375}
	b377.items = []builder{&b341, &b833, &b375, &b376}
	b378.options = []builder{&b349, &b356, &b363, &b370, &b377}
	var b391 = sequenceBuilder{id: 391, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b384 = sequenceBuilder{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b383 = sequenceBuilder{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b383.items = []builder{&b833, &b14}
	b384.items = []builder{&b833, &b14, &b383}
	var b380 = sequenceBuilder{id: 380, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b379 = charBuilder{}
	b380.items = []builder{&b379}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b385.items = []builder{&b833, &b14}
	b386.items = []builder{&b833, &b14, &b385}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b387.items = []builder{&b833, &b14}
	b388.items = []builder{&b833, &b14, &b387}
	var b382 = sequenceBuilder{id: 382, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b381 = charBuilder{}
	b382.items = []builder{&b381}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b389.items = []builder{&b833, &b14}
	b390.items = []builder{&b833, &b14, &b389}
	b391.items = []builder{&b400, &b384, &b833, &b380, &b386, &b833, &b400, &b388, &b833, &b382, &b390, &b833, &b400}
	var b399 = sequenceBuilder{id: 399, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b392 = choiceBuilder{id: 392, commit: 66}
	b392.options = []builder{&b271, &b331, &b378, &b391}
	var b397 = sequenceBuilder{id: 397, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b394 = sequenceBuilder{id: 394, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b393.items = []builder{&b833, &b14}
	b394.items = []builder{&b14, &b393}
	var b329 = sequenceBuilder{id: 329, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b327 = charBuilder{}
	var b328 = charBuilder{}
	b329.items = []builder{&b327, &b328}
	var b396 = sequenceBuilder{id: 396, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b395 = sequenceBuilder{id: 395, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b395.items = []builder{&b833, &b14}
	b396.items = []builder{&b833, &b14, &b395}
	b397.items = []builder{&b394, &b833, &b329, &b396, &b833, &b392}
	var b398 = sequenceBuilder{id: 398, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b398.items = []builder{&b833, &b397}
	b399.items = []builder{&b392, &b833, &b397, &b398}
	b400.options = []builder{&b271, &b331, &b378, &b391, &b399}
	b184.items = []builder{&b183, &b833, &b400}
	b185.items = []builder{&b181, &b833, &b184}
	var b437 = sequenceBuilder{id: 437, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b403 = sequenceBuilder{id: 403, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b401 = charBuilder{}
	var b402 = charBuilder{}
	b403.items = []builder{&b401, &b402}
	var b432 = sequenceBuilder{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b431 = sequenceBuilder{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b431.items = []builder{&b833, &b14}
	b432.items = []builder{&b833, &b14, &b431}
	var b434 = sequenceBuilder{id: 434, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b433 = sequenceBuilder{id: 433, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b433.items = []builder{&b833, &b14}
	b434.items = []builder{&b833, &b14, &b433}
	var b436 = sequenceBuilder{id: 436, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b420 = sequenceBuilder{id: 420, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b413 = sequenceBuilder{id: 413, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b412 = sequenceBuilder{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b412.items = []builder{&b833, &b14}
	b413.items = []builder{&b14, &b412}
	var b408 = sequenceBuilder{id: 408, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b404 = charBuilder{}
	var b405 = charBuilder{}
	var b406 = charBuilder{}
	var b407 = charBuilder{}
	b408.items = []builder{&b404, &b405, &b406, &b407}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b414.items = []builder{&b833, &b14}
	b415.items = []builder{&b833, &b14, &b414}
	var b411 = sequenceBuilder{id: 411, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b409 = charBuilder{}
	var b410 = charBuilder{}
	b411.items = []builder{&b409, &b410}
	var b417 = sequenceBuilder{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b416 = sequenceBuilder{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b416.items = []builder{&b833, &b14}
	b417.items = []builder{&b833, &b14, &b416}
	var b419 = sequenceBuilder{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b418 = sequenceBuilder{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b418.items = []builder{&b833, &b14}
	b419.items = []builder{&b833, &b14, &b418}
	b420.items = []builder{&b413, &b833, &b408, &b415, &b833, &b411, &b417, &b833, &b400, &b419, &b833, &b190}
	var b435 = sequenceBuilder{id: 435, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b435.items = []builder{&b833, &b420}
	b436.items = []builder{&b833, &b420, &b435}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b426.items = []builder{&b833, &b14}
	b427.items = []builder{&b14, &b426}
	var b425 = sequenceBuilder{id: 425, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b421 = charBuilder{}
	var b422 = charBuilder{}
	var b423 = charBuilder{}
	var b424 = charBuilder{}
	b425.items = []builder{&b421, &b422, &b423, &b424}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b428.items = []builder{&b833, &b14}
	b429.items = []builder{&b833, &b14, &b428}
	b430.items = []builder{&b427, &b833, &b425, &b429, &b833, &b190}
	b437.items = []builder{&b403, &b432, &b833, &b400, &b434, &b833, &b190, &b436, &b833, &b430}
	var b494 = sequenceBuilder{id: 494, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b479 = sequenceBuilder{id: 479, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b473 = charBuilder{}
	var b474 = charBuilder{}
	var b475 = charBuilder{}
	var b476 = charBuilder{}
	var b477 = charBuilder{}
	var b478 = charBuilder{}
	b479.items = []builder{&b473, &b474, &b475, &b476, &b477, &b478}
	var b491 = sequenceBuilder{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b490 = sequenceBuilder{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b490.items = []builder{&b833, &b14}
	b491.items = []builder{&b833, &b14, &b490}
	var b493 = sequenceBuilder{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b492.items = []builder{&b833, &b14}
	b493.items = []builder{&b833, &b14, &b492}
	var b481 = sequenceBuilder{id: 481, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b480 = charBuilder{}
	b481.items = []builder{&b480}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b482 = choiceBuilder{id: 482, commit: 2}
	var b472 = sequenceBuilder{id: 472, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b467 = sequenceBuilder{id: 467, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b460 = sequenceBuilder{id: 460, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b456 = charBuilder{}
	var b457 = charBuilder{}
	var b458 = charBuilder{}
	var b459 = charBuilder{}
	b460.items = []builder{&b456, &b457, &b458, &b459}
	var b464 = sequenceBuilder{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b463 = sequenceBuilder{id: 463, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b463.items = []builder{&b833, &b14}
	b464.items = []builder{&b833, &b14, &b463}
	var b466 = sequenceBuilder{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b465 = sequenceBuilder{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b465.items = []builder{&b833, &b14}
	b466.items = []builder{&b833, &b14, &b465}
	var b462 = sequenceBuilder{id: 462, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b461 = charBuilder{}
	b462.items = []builder{&b461}
	b467.items = []builder{&b460, &b464, &b833, &b400, &b466, &b833, &b462}
	var b471 = sequenceBuilder{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b469 = sequenceBuilder{id: 469, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b468 = charBuilder{}
	b469.items = []builder{&b468}
	var b470 = sequenceBuilder{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b470.items = []builder{&b833, &b469}
	b471.items = []builder{&b833, &b469, &b470}
	b472.items = []builder{&b467, &b471, &b833, &b801}
	var b455 = sequenceBuilder{id: 455, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b450 = sequenceBuilder{id: 450, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b445 = sequenceBuilder{id: 445, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b438 = charBuilder{}
	var b439 = charBuilder{}
	var b440 = charBuilder{}
	var b441 = charBuilder{}
	var b442 = charBuilder{}
	var b443 = charBuilder{}
	var b444 = charBuilder{}
	b445.items = []builder{&b438, &b439, &b440, &b441, &b442, &b443, &b444}
	var b449 = sequenceBuilder{id: 449, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b448 = sequenceBuilder{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b448.items = []builder{&b833, &b14}
	b449.items = []builder{&b833, &b14, &b448}
	var b447 = sequenceBuilder{id: 447, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b446 = charBuilder{}
	b447.items = []builder{&b446}
	b450.items = []builder{&b445, &b449, &b833, &b447}
	var b454 = sequenceBuilder{id: 454, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b452 = sequenceBuilder{id: 452, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b451 = charBuilder{}
	b452.items = []builder{&b451}
	var b453 = sequenceBuilder{id: 453, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b453.items = []builder{&b833, &b452}
	b454.items = []builder{&b833, &b452, &b453}
	b455.items = []builder{&b450, &b454, &b833, &b801}
	b482.options = []builder{&b472, &b455}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b484 = sequenceBuilder{id: 484, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b483 = choiceBuilder{id: 483, commit: 2}
	b483.options = []builder{&b472, &b455, &b801}
	b484.items = []builder{&b815, &b833, &b483}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b485.items = []builder{&b833, &b484}
	b486.items = []builder{&b833, &b484, &b485}
	b487.items = []builder{&b482, &b486}
	var b489 = sequenceBuilder{id: 489, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b488 = charBuilder{}
	b489.items = []builder{&b488}
	b494.items = []builder{&b479, &b491, &b833, &b400, &b493, &b833, &b481, &b833, &b815, &b833, &b487, &b833, &b815, &b833, &b489}
	var b556 = sequenceBuilder{id: 556, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b543 = sequenceBuilder{id: 543, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b537 = charBuilder{}
	var b538 = charBuilder{}
	var b539 = charBuilder{}
	var b540 = charBuilder{}
	var b541 = charBuilder{}
	var b542 = charBuilder{}
	b543.items = []builder{&b537, &b538, &b539, &b540, &b541, &b542}
	var b555 = sequenceBuilder{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b554 = sequenceBuilder{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b554.items = []builder{&b833, &b14}
	b555.items = []builder{&b833, &b14, &b554}
	var b545 = sequenceBuilder{id: 545, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b544 = charBuilder{}
	b545.items = []builder{&b544}
	var b551 = sequenceBuilder{id: 551, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b546 = choiceBuilder{id: 546, commit: 2}
	var b536 = sequenceBuilder{id: 536, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b531 = sequenceBuilder{id: 531, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b524 = sequenceBuilder{id: 524, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b520 = charBuilder{}
	var b521 = charBuilder{}
	var b522 = charBuilder{}
	var b523 = charBuilder{}
	b524.items = []builder{&b520, &b521, &b522, &b523}
	var b528 = sequenceBuilder{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b527 = sequenceBuilder{id: 527, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b527.items = []builder{&b833, &b14}
	b528.items = []builder{&b833, &b14, &b527}
	var b519 = choiceBuilder{id: 519, commit: 66}
	var b518 = sequenceBuilder{id: 518, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b517 = sequenceBuilder{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b516 = sequenceBuilder{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b516.items = []builder{&b833, &b14}
	b517.items = []builder{&b833, &b14, &b516}
	b518.items = []builder{&b103, &b517, &b833, &b515}
	b519.options = []builder{&b504, &b515, &b518}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b529.items = []builder{&b833, &b14}
	b530.items = []builder{&b833, &b14, &b529}
	var b526 = sequenceBuilder{id: 526, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b525 = charBuilder{}
	b526.items = []builder{&b525}
	b531.items = []builder{&b524, &b528, &b833, &b519, &b530, &b833, &b526}
	var b535 = sequenceBuilder{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b533 = sequenceBuilder{id: 533, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b532 = charBuilder{}
	b533.items = []builder{&b532}
	var b534 = sequenceBuilder{id: 534, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b534.items = []builder{&b833, &b533}
	b535.items = []builder{&b833, &b533, &b534}
	b536.items = []builder{&b531, &b535, &b833, &b801}
	b546.options = []builder{&b536, &b455}
	var b550 = sequenceBuilder{id: 550, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b548 = sequenceBuilder{id: 548, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b547 = choiceBuilder{id: 547, commit: 2}
	b547.options = []builder{&b536, &b455, &b801}
	b548.items = []builder{&b815, &b833, &b547}
	var b549 = sequenceBuilder{id: 549, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b549.items = []builder{&b833, &b548}
	b550.items = []builder{&b833, &b548, &b549}
	b551.items = []builder{&b546, &b550}
	var b553 = sequenceBuilder{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b552 = charBuilder{}
	b553.items = []builder{&b552}
	b556.items = []builder{&b543, &b555, &b833, &b545, &b833, &b815, &b833, &b551, &b833, &b815, &b833, &b553}
	var b597 = sequenceBuilder{id: 597, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b586 = sequenceBuilder{id: 586, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b583 = charBuilder{}
	var b584 = charBuilder{}
	var b585 = charBuilder{}
	b586.items = []builder{&b583, &b584, &b585}
	var b596 = choiceBuilder{id: 596, commit: 2}
	var b592 = sequenceBuilder{id: 592, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b587.items = []builder{&b833, &b14}
	b588.items = []builder{&b14, &b587}
	var b582 = choiceBuilder{id: 582, commit: 66}
	var b581 = choiceBuilder{id: 581, commit: 64, name: "range-over-expression"}
	var b580 = sequenceBuilder{id: 580, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b577 = sequenceBuilder{id: 577, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b576 = sequenceBuilder{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b576.items = []builder{&b833, &b14}
	b577.items = []builder{&b833, &b14, &b576}
	var b574 = sequenceBuilder{id: 574, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b572 = charBuilder{}
	var b573 = charBuilder{}
	b574.items = []builder{&b572, &b573}
	var b579 = sequenceBuilder{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b578 = sequenceBuilder{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b578.items = []builder{&b833, &b14}
	b579.items = []builder{&b833, &b14, &b578}
	var b575 = choiceBuilder{id: 575, commit: 2}
	b575.options = []builder{&b400, &b225}
	b580.items = []builder{&b103, &b577, &b833, &b574, &b579, &b833, &b575}
	b581.options = []builder{&b580, &b225}
	b582.options = []builder{&b400, &b581}
	b589.items = []builder{&b588, &b833, &b582}
	var b591 = sequenceBuilder{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b590.items = []builder{&b833, &b14}
	b591.items = []builder{&b833, &b14, &b590}
	b592.items = []builder{&b589, &b591, &b833, &b190}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b593.items = []builder{&b833, &b14}
	b594.items = []builder{&b14, &b593}
	b595.items = []builder{&b594, &b833, &b190}
	b596.options = []builder{&b592, &b595}
	b597.items = []builder{&b586, &b833, &b596}
	var b745 = choiceBuilder{id: 745, commit: 66}
	var b658 = sequenceBuilder{id: 658, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b654 = sequenceBuilder{id: 654, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b651 = charBuilder{}
	var b652 = charBuilder{}
	var b653 = charBuilder{}
	b654.items = []builder{&b651, &b652, &b653}
	var b657 = sequenceBuilder{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b656 = sequenceBuilder{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b656.items = []builder{&b833, &b14}
	b657.items = []builder{&b833, &b14, &b656}
	var b655 = choiceBuilder{id: 655, commit: 2}
	var b645 = sequenceBuilder{id: 645, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b644 = sequenceBuilder{id: 644, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b641 = sequenceBuilder{id: 641, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b640 = sequenceBuilder{id: 640, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b639 = sequenceBuilder{id: 639, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b639.items = []builder{&b833, &b14}
	b640.items = []builder{&b14, &b639}
	var b638 = sequenceBuilder{id: 638, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b637 = charBuilder{}
	b638.items = []builder{&b637}
	b641.items = []builder{&b640, &b833, &b638}
	var b643 = sequenceBuilder{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b642.items = []builder{&b833, &b14}
	b643.items = []builder{&b833, &b14, &b642}
	b644.items = []builder{&b103, &b833, &b641, &b643, &b833, &b400}
	b645.items = []builder{&b644}
	var b650 = sequenceBuilder{id: 650, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b647 = sequenceBuilder{id: 647, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b646 = charBuilder{}
	b647.items = []builder{&b646}
	var b649 = sequenceBuilder{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b648 = sequenceBuilder{id: 648, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b648.items = []builder{&b833, &b14}
	b649.items = []builder{&b833, &b14, &b648}
	b650.items = []builder{&b647, &b649, &b833, &b644}
	b655.options = []builder{&b645, &b650}
	b658.items = []builder{&b654, &b657, &b833, &b655}
	var b679 = sequenceBuilder{id: 679, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b672 = sequenceBuilder{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b669 = charBuilder{}
	var b670 = charBuilder{}
	var b671 = charBuilder{}
	b672.items = []builder{&b669, &b670, &b671}
	var b678 = sequenceBuilder{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b677 = sequenceBuilder{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b677.items = []builder{&b833, &b14}
	b678.items = []builder{&b833, &b14, &b677}
	var b674 = sequenceBuilder{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b673 = charBuilder{}
	b674.items = []builder{&b673}
	var b664 = sequenceBuilder{id: 664, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b659 = choiceBuilder{id: 659, commit: 2}
	b659.options = []builder{&b645, &b650}
	var b663 = sequenceBuilder{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b661 = sequenceBuilder{id: 661, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b660 = choiceBuilder{id: 660, commit: 2}
	b660.options = []builder{&b645, &b650}
	b661.items = []builder{&b113, &b833, &b660}
	var b662 = sequenceBuilder{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b662.items = []builder{&b833, &b661}
	b663.items = []builder{&b833, &b661, &b662}
	b664.items = []builder{&b659, &b663}
	var b676 = sequenceBuilder{id: 676, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b675 = charBuilder{}
	b676.items = []builder{&b675}
	b679.items = []builder{&b672, &b678, &b833, &b674, &b833, &b113, &b833, &b664, &b833, &b113, &b833, &b676}
	var b694 = sequenceBuilder{id: 694, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b683 = sequenceBuilder{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	var b681 = charBuilder{}
	var b682 = charBuilder{}
	b683.items = []builder{&b680, &b681, &b682}
	var b691 = sequenceBuilder{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b690 = sequenceBuilder{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b690.items = []builder{&b833, &b14}
	b691.items = []builder{&b833, &b14, &b690}
	var b685 = sequenceBuilder{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b684 = charBuilder{}
	b685.items = []builder{&b684}
	var b693 = sequenceBuilder{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b692 = sequenceBuilder{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b692.items = []builder{&b833, &b14}
	b693.items = []builder{&b833, &b14, &b692}
	var b687 = sequenceBuilder{id: 687, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b686 = charBuilder{}
	b687.items = []builder{&b686}
	var b668 = sequenceBuilder{id: 668, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b665 = sequenceBuilder{id: 665, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b665.items = []builder{&b113, &b833, &b645}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b666.items = []builder{&b833, &b665}
	b667.items = []builder{&b833, &b665, &b666}
	b668.items = []builder{&b645, &b667}
	var b689 = sequenceBuilder{id: 689, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b688 = charBuilder{}
	b689.items = []builder{&b688}
	b694.items = []builder{&b683, &b691, &b833, &b685, &b693, &b833, &b687, &b833, &b113, &b833, &b668, &b833, &b113, &b833, &b689}
	var b710 = sequenceBuilder{id: 710, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b706 = sequenceBuilder{id: 706, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b704 = charBuilder{}
	var b705 = charBuilder{}
	b706.items = []builder{&b704, &b705}
	var b709 = sequenceBuilder{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b708 = sequenceBuilder{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b708.items = []builder{&b833, &b14}
	b709.items = []builder{&b833, &b14, &b708}
	var b707 = choiceBuilder{id: 707, commit: 2}
	var b698 = sequenceBuilder{id: 698, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b697 = sequenceBuilder{id: 697, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b695 = sequenceBuilder{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b695.items = []builder{&b833, &b14}
	b696.items = []builder{&b833, &b14, &b695}
	b697.items = []builder{&b103, &b696, &b833, &b200}
	b698.items = []builder{&b697}
	var b703 = sequenceBuilder{id: 703, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b700 = sequenceBuilder{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b699 = charBuilder{}
	b700.items = []builder{&b699}
	var b702 = sequenceBuilder{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b701 = sequenceBuilder{id: 701, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b701.items = []builder{&b833, &b14}
	b702.items = []builder{&b833, &b14, &b701}
	b703.items = []builder{&b700, &b702, &b833, &b697}
	b707.options = []builder{&b698, &b703}
	b710.items = []builder{&b706, &b709, &b833, &b707}
	var b730 = sequenceBuilder{id: 730, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b723 = sequenceBuilder{id: 723, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b721 = charBuilder{}
	var b722 = charBuilder{}
	b723.items = []builder{&b721, &b722}
	var b729 = sequenceBuilder{id: 729, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b728 = sequenceBuilder{id: 728, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b728.items = []builder{&b833, &b14}
	b729.items = []builder{&b833, &b14, &b728}
	var b725 = sequenceBuilder{id: 725, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b724 = charBuilder{}
	b725.items = []builder{&b724}
	var b720 = sequenceBuilder{id: 720, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b715 = choiceBuilder{id: 715, commit: 2}
	b715.options = []builder{&b698, &b703}
	var b719 = sequenceBuilder{id: 719, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b717 = sequenceBuilder{id: 717, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b716 = choiceBuilder{id: 716, commit: 2}
	b716.options = []builder{&b698, &b703}
	b717.items = []builder{&b113, &b833, &b716}
	var b718 = sequenceBuilder{id: 718, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b718.items = []builder{&b833, &b717}
	b719.items = []builder{&b833, &b717, &b718}
	b720.items = []builder{&b715, &b719}
	var b727 = sequenceBuilder{id: 727, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b726 = charBuilder{}
	b727.items = []builder{&b726}
	b730.items = []builder{&b723, &b729, &b833, &b725, &b833, &b113, &b833, &b720, &b833, &b113, &b833, &b727}
	var b744 = sequenceBuilder{id: 744, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b733 = sequenceBuilder{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b731 = charBuilder{}
	var b732 = charBuilder{}
	b733.items = []builder{&b731, &b732}
	var b741 = sequenceBuilder{id: 741, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b740 = sequenceBuilder{id: 740, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b740.items = []builder{&b833, &b14}
	b741.items = []builder{&b833, &b14, &b740}
	var b735 = sequenceBuilder{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b734 = charBuilder{}
	b735.items = []builder{&b734}
	var b743 = sequenceBuilder{id: 743, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b742 = sequenceBuilder{id: 742, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b742.items = []builder{&b833, &b14}
	b743.items = []builder{&b833, &b14, &b742}
	var b737 = sequenceBuilder{id: 737, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b736 = charBuilder{}
	b737.items = []builder{&b736}
	var b714 = sequenceBuilder{id: 714, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b711 = sequenceBuilder{id: 711, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b711.items = []builder{&b113, &b833, &b698}
	var b712 = sequenceBuilder{id: 712, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b712.items = []builder{&b833, &b711}
	b713.items = []builder{&b833, &b711, &b712}
	b714.items = []builder{&b698, &b713}
	var b739 = sequenceBuilder{id: 739, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b738 = charBuilder{}
	b739.items = []builder{&b738}
	b744.items = []builder{&b733, &b741, &b833, &b735, &b743, &b833, &b737, &b833, &b113, &b833, &b714, &b833, &b113, &b833, &b739}
	b745.options = []builder{&b658, &b679, &b694, &b710, &b730, &b744}
	var b780 = choiceBuilder{id: 780, commit: 64, name: "use"}
	var b768 = sequenceBuilder{id: 768, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b765 = sequenceBuilder{id: 765, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b762 = charBuilder{}
	var b763 = charBuilder{}
	var b764 = charBuilder{}
	b765.items = []builder{&b762, &b763, &b764}
	var b767 = sequenceBuilder{id: 767, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b766 = sequenceBuilder{id: 766, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b766.items = []builder{&b833, &b14}
	b767.items = []builder{&b833, &b14, &b766}
	var b757 = choiceBuilder{id: 757, commit: 64, name: "use-fact"}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b748 = choiceBuilder{id: 748, commit: 2}
	var b747 = sequenceBuilder{id: 747, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b746 = charBuilder{}
	b747.items = []builder{&b746}
	b748.options = []builder{&b103, &b747}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b751 = sequenceBuilder{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b751.items = []builder{&b833, &b14}
	b752.items = []builder{&b14, &b751}
	var b750 = sequenceBuilder{id: 750, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b749 = charBuilder{}
	b750.items = []builder{&b749}
	b753.items = []builder{&b752, &b833, &b750}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b754.items = []builder{&b833, &b14}
	b755.items = []builder{&b833, &b14, &b754}
	b756.items = []builder{&b748, &b833, &b753, &b755, &b833, &b86}
	b757.options = []builder{&b86, &b756}
	b768.items = []builder{&b765, &b767, &b833, &b757}
	var b779 = sequenceBuilder{id: 779, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b772 = sequenceBuilder{id: 772, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b769 = charBuilder{}
	var b770 = charBuilder{}
	var b771 = charBuilder{}
	b772.items = []builder{&b769, &b770, &b771}
	var b778 = sequenceBuilder{id: 778, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b777 = sequenceBuilder{id: 777, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b777.items = []builder{&b833, &b14}
	b778.items = []builder{&b833, &b14, &b777}
	var b774 = sequenceBuilder{id: 774, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b773 = charBuilder{}
	b774.items = []builder{&b773}
	var b761 = sequenceBuilder{id: 761, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b760 = sequenceBuilder{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b758 = sequenceBuilder{id: 758, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b758.items = []builder{&b113, &b833, &b757}
	var b759 = sequenceBuilder{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b759.items = []builder{&b833, &b758}
	b760.items = []builder{&b833, &b758, &b759}
	b761.items = []builder{&b757, &b760}
	var b776 = sequenceBuilder{id: 776, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b775 = charBuilder{}
	b776.items = []builder{&b775}
	b779.items = []builder{&b772, &b778, &b833, &b774, &b833, &b113, &b833, &b761, &b833, &b113, &b833, &b776}
	b780.options = []builder{&b768, &b779}
	var b790 = sequenceBuilder{id: 790, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b787 = sequenceBuilder{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b781 = charBuilder{}
	var b782 = charBuilder{}
	var b783 = charBuilder{}
	var b784 = charBuilder{}
	var b785 = charBuilder{}
	var b786 = charBuilder{}
	b787.items = []builder{&b781, &b782, &b783, &b784, &b785, &b786}
	var b789 = sequenceBuilder{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b788 = sequenceBuilder{id: 788, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b788.items = []builder{&b833, &b14}
	b789.items = []builder{&b833, &b14, &b788}
	b790.items = []builder{&b787, &b789, &b833, &b745}
	var b810 = sequenceBuilder{id: 810, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b803 = sequenceBuilder{id: 803, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b802 = charBuilder{}
	b803.items = []builder{&b802}
	var b807 = sequenceBuilder{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b806 = sequenceBuilder{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b806.items = []builder{&b833, &b14}
	b807.items = []builder{&b833, &b14, &b806}
	var b809 = sequenceBuilder{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b808 = sequenceBuilder{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b808.items = []builder{&b833, &b14}
	b809.items = []builder{&b833, &b14, &b808}
	var b805 = sequenceBuilder{id: 805, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b804 = charBuilder{}
	b805.items = []builder{&b804}
	b810.items = []builder{&b803, &b807, &b833, &b801, &b809, &b833, &b805}
	b801.options = []builder{&b185, &b437, &b494, &b556, &b597, &b745, &b780, &b790, &b810, &b791}
	var b818 = sequenceBuilder{id: 818, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b816 = sequenceBuilder{id: 816, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b816.items = []builder{&b815, &b833, &b801}
	var b817 = sequenceBuilder{id: 817, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b817.items = []builder{&b833, &b816}
	b818.items = []builder{&b833, &b816, &b817}
	b819.items = []builder{&b801, &b818}
	b834.items = []builder{&b830, &b833, &b815, &b833, &b819, &b833, &b815}
	b835.items = []builder{&b833, &b834, &b833}

	return parseInput(r, &p835, &b835)
}
