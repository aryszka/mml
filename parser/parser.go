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
	var p819 = sequenceParser{id: 819, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p817 = choiceParser{id: 817, commit: 2}
	var p815 = choiceParser{id: 815, commit: 70, name: "ws", generalizations: []int{817}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 817}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 817}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 817}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 817}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 817}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{815, 817}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p815.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p816 = sequenceParser{id: 816, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{817}}
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
	p41.items = []parser{&p40, &p817, &p38}
	var p42 = sequenceParser{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p42.items = []parser{&p817, &p41}
	p43.items = []parser{&p817, &p41, &p42}
	p44.items = []parser{&p38, &p43}
	p816.items = []parser{&p44}
	p817.options = []parser{&p815, &p816}
	var p818 = sequenceParser{id: 818, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p814 = sequenceParser{id: 814, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p811 = sequenceParser{id: 811, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p809 = charParser{id: 809, chars: []rune{35}}
	var p810 = charParser{id: 810, chars: []rune{33}}
	p811.items = []parser{&p809, &p810}
	var p808 = sequenceParser{id: 808, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p807 = sequenceParser{id: 807, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p805 = sequenceParser{id: 805, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p804 = charParser{id: 804, not: true, chars: []rune{10}}
	p805.items = []parser{&p804}
	var p806 = sequenceParser{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p806.items = []parser{&p817, &p805}
	p807.items = []parser{&p805, &p806}
	p808.items = []parser{&p807}
	var p813 = sequenceParser{id: 813, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p812 = charParser{id: 812, chars: []rune{10}}
	p813.items = []parser{&p812}
	p814.items = []parser{&p811, &p817, &p808, &p817, &p813}
	var p799 = sequenceParser{id: 799, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p797 = choiceParser{id: 797, commit: 2}
	var p796 = sequenceParser{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{797}}
	var p795 = charParser{id: 795, chars: []rune{59}}
	p796.items = []parser{&p795}
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{797, 113}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p797.options = []parser{&p796, &p14}
	var p798 = sequenceParser{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p798.items = []parser{&p817, &p797}
	p799.items = []parser{&p797, &p798}
	var p803 = sequenceParser{id: 803, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p785 = choiceParser{id: 785, commit: 66, name: "statement", generalizations: []int{459, 523}}
	var p187 = sequenceParser{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{785, 459, 523}}
	var p183 = sequenceParser{id: 183, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p177 = charParser{id: 177, chars: []rune{114}}
	var p178 = charParser{id: 178, chars: []rune{101}}
	var p179 = charParser{id: 179, chars: []rune{116}}
	var p180 = charParser{id: 180, chars: []rune{117}}
	var p181 = charParser{id: 181, chars: []rune{114}}
	var p182 = charParser{id: 182, chars: []rune{110}}
	p183.items = []parser{&p177, &p178, &p179, &p180, &p181, &p182}
	var p186 = sequenceParser{id: 186, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p185 = sequenceParser{id: 185, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p184 = sequenceParser{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p184.items = []parser{&p817, &p14}
	p185.items = []parser{&p14, &p184}
	var p376 = choiceParser{id: 376, commit: 66, name: "expression", generalizations: []int{116, 775, 199, 558, 551, 785}}
	var p268 = choiceParser{id: 268, commit: 66, name: "primary-expression", generalizations: []int{116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p62 = choiceParser{id: 62, commit: 64, name: "int", generalizations: []int{268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p53 = sequenceParser{id: 53, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p52 = sequenceParser{id: 52, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p51 = charParser{id: 51, ranges: [][]rune{{49, 57}}}
	p52.items = []parser{&p51}
	var p46 = sequenceParser{id: 46, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 57}}}
	p46.items = []parser{&p45}
	p53.items = []parser{&p52, &p46}
	var p56 = sequenceParser{id: 56, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p55 = sequenceParser{id: 55, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p54 = charParser{id: 54, chars: []rune{48}}
	p55.items = []parser{&p54}
	var p48 = sequenceParser{id: 48, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p47 = charParser{id: 47, ranges: [][]rune{{48, 55}}}
	p48.items = []parser{&p47}
	p56.items = []parser{&p55, &p48}
	var p61 = sequenceParser{id: 61, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{62, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
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
	var p75 = choiceParser{id: 75, commit: 72, name: "float", generalizations: []int{268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p70 = sequenceParser{id: 70, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{75, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
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
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{75, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p72 = sequenceParser{id: 72, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p71 = charParser{id: 71, chars: []rune{46}}
	p72.items = []parser{&p71}
	p73.items = []parser{&p72, &p46, &p67}
	var p74 = sequenceParser{id: 74, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{75, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	p74.items = []parser{&p46, &p67}
	p75.options = []parser{&p70, &p73, &p74}
	var p88 = sequenceParser{id: 88, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 116, 141, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 733, 785}}
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
	var p100 = choiceParser{id: 100, commit: 66, name: "bool", generalizations: []int{268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p93 = sequenceParser{id: 93, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p89 = charParser{id: 89, chars: []rune{116}}
	var p90 = charParser{id: 90, chars: []rune{114}}
	var p91 = charParser{id: 91, chars: []rune{117}}
	var p92 = charParser{id: 92, chars: []rune{101}}
	p93.items = []parser{&p89, &p90, &p91, &p92}
	var p99 = sequenceParser{id: 99, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p94 = charParser{id: 94, chars: []rune{102}}
	var p95 = charParser{id: 95, chars: []rune{97}}
	var p96 = charParser{id: 96, chars: []rune{108}}
	var p97 = charParser{id: 97, chars: []rune{115}}
	var p98 = charParser{id: 98, chars: []rune{101}}
	p99.items = []parser{&p94, &p95, &p96, &p97, &p98}
	p100.options = []parser{&p93, &p99}
	var p481 = sequenceParser{id: 481, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 116, 775, 199, 376, 334, 335, 336, 337, 338, 339, 495, 558, 551, 785}}
	var p478 = sequenceParser{id: 478, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p471 = charParser{id: 471, chars: []rune{114}}
	var p472 = charParser{id: 472, chars: []rune{101}}
	var p473 = charParser{id: 473, chars: []rune{99}}
	var p474 = charParser{id: 474, chars: []rune{101}}
	var p475 = charParser{id: 475, chars: []rune{105}}
	var p476 = charParser{id: 476, chars: []rune{118}}
	var p477 = charParser{id: 477, chars: []rune{101}}
	p478.items = []parser{&p471, &p472, &p473, &p474, &p475, &p476, &p477}
	var p480 = sequenceParser{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p479 = sequenceParser{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p479.items = []parser{&p817, &p14}
	p480.items = []parser{&p817, &p14, &p479}
	p481.items = []parser{&p478, &p480, &p817, &p268}
	var p105 = sequenceParser{id: 105, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{268, 116, 141, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 724, 785}}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p102.items = []parser{&p101}
	var p104 = sequenceParser{id: 104, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p103 = charParser{id: 103, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p104.items = []parser{&p103}
	p105.items = []parser{&p102, &p104}
	var p126 = sequenceParser{id: 126, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{116, 268, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
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
	p114.items = []parser{&p817, &p113}
	p115.items = []parser{&p113, &p114}
	var p120 = sequenceParser{id: 120, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p116 = choiceParser{id: 116, commit: 66, name: "list-item"}
	var p110 = sequenceParser{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{116, 149, 150}}
	var p109 = sequenceParser{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	var p108 = charParser{id: 108, chars: []rune{46}}
	p109.items = []parser{&p106, &p107, &p108}
	p110.items = []parser{&p268, &p817, &p109}
	p116.options = []parser{&p376, &p110}
	var p119 = sequenceParser{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p117.items = []parser{&p115, &p817, &p116}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p118.items = []parser{&p817, &p117}
	p119.items = []parser{&p817, &p117, &p118}
	p120.items = []parser{&p116, &p119}
	var p124 = sequenceParser{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p123 = charParser{id: 123, chars: []rune{93}}
	p124.items = []parser{&p123}
	p125.items = []parser{&p122, &p817, &p115, &p817, &p120, &p817, &p115, &p817, &p124}
	p126.items = []parser{&p125}
	var p131 = sequenceParser{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p128 = sequenceParser{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p127 = charParser{id: 127, chars: []rune{126}}
	p128.items = []parser{&p127}
	var p130 = sequenceParser{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p129.items = []parser{&p817, &p14}
	p130.items = []parser{&p817, &p14, &p129}
	p131.items = []parser{&p128, &p130, &p817, &p125}
	var p160 = sequenceParser{id: 160, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{268, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
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
	p136.items = []parser{&p817, &p14}
	p137.items = []parser{&p817, &p14, &p136}
	var p139 = sequenceParser{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p138.items = []parser{&p817, &p14}
	p139.items = []parser{&p817, &p14, &p138}
	var p135 = sequenceParser{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p134 = charParser{id: 134, chars: []rune{93}}
	p135.items = []parser{&p134}
	p140.items = []parser{&p133, &p137, &p817, &p376, &p139, &p817, &p135}
	p141.options = []parser{&p105, &p88, &p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p817, &p14}
	p145.items = []parser{&p817, &p14, &p144}
	var p143 = sequenceParser{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p142 = charParser{id: 142, chars: []rune{58}}
	p143.items = []parser{&p142}
	var p147 = sequenceParser{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p146.items = []parser{&p817, &p14}
	p147.items = []parser{&p817, &p14, &p146}
	p148.items = []parser{&p141, &p145, &p817, &p143, &p147, &p817, &p376}
	p149.options = []parser{&p148, &p110}
	var p153 = sequenceParser{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p150 = choiceParser{id: 150, commit: 2}
	p150.options = []parser{&p148, &p110}
	p151.items = []parser{&p115, &p817, &p150}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p152.items = []parser{&p817, &p151}
	p153.items = []parser{&p817, &p151, &p152}
	p154.items = []parser{&p149, &p153}
	var p158 = sequenceParser{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p157 = charParser{id: 157, chars: []rune{125}}
	p158.items = []parser{&p157}
	p159.items = []parser{&p156, &p817, &p115, &p817, &p154, &p817, &p115, &p817, &p158}
	p160.items = []parser{&p159}
	var p165 = sequenceParser{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 775, 199, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p162 = sequenceParser{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p161 = charParser{id: 161, chars: []rune{126}}
	p162.items = []parser{&p161}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p163.items = []parser{&p817, &p14}
	p164.items = []parser{&p817, &p14, &p163}
	p165.items = []parser{&p162, &p164, &p817, &p159}
	var p208 = sequenceParser{id: 208, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{775, 199, 268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p205 = sequenceParser{id: 205, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p203 = charParser{id: 203, chars: []rune{102}}
	var p204 = charParser{id: 204, chars: []rune{110}}
	p205.items = []parser{&p203, &p204}
	var p207 = sequenceParser{id: 207, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p206 = sequenceParser{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p206.items = []parser{&p817, &p14}
	p207.items = []parser{&p817, &p14, &p206}
	var p202 = sequenceParser{id: 202, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p194 = sequenceParser{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p193 = charParser{id: 193, chars: []rune{40}}
	p194.items = []parser{&p193}
	var p196 = choiceParser{id: 196, commit: 2}
	var p169 = sequenceParser{id: 169, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{196}}
	var p168 = sequenceParser{id: 168, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p166.items = []parser{&p115, &p817, &p105}
	var p167 = sequenceParser{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p167.items = []parser{&p817, &p166}
	p168.items = []parser{&p817, &p166, &p167}
	p169.items = []parser{&p105, &p168}
	var p195 = sequenceParser{id: 195, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{196}}
	var p176 = sequenceParser{id: 176, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{196}}
	var p173 = sequenceParser{id: 173, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p170 = charParser{id: 170, chars: []rune{46}}
	var p171 = charParser{id: 171, chars: []rune{46}}
	var p172 = charParser{id: 172, chars: []rune{46}}
	p173.items = []parser{&p170, &p171, &p172}
	var p175 = sequenceParser{id: 175, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p174 = sequenceParser{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p174.items = []parser{&p817, &p14}
	p175.items = []parser{&p817, &p14, &p174}
	p176.items = []parser{&p173, &p175, &p817, &p105}
	p195.items = []parser{&p169, &p817, &p115, &p817, &p176}
	p196.options = []parser{&p169, &p195, &p176}
	var p198 = sequenceParser{id: 198, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p197 = charParser{id: 197, chars: []rune{41}}
	p198.items = []parser{&p197}
	var p201 = sequenceParser{id: 201, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p200 = sequenceParser{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p200.items = []parser{&p817, &p14}
	p201.items = []parser{&p817, &p14, &p200}
	var p199 = choiceParser{id: 199, commit: 2}
	var p775 = choiceParser{id: 775, commit: 66, name: "simple-statement", generalizations: []int{199, 785}}
	var p491 = sequenceParser{id: 491, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{775, 199, 495, 785}}
	var p486 = sequenceParser{id: 486, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p482 = charParser{id: 482, chars: []rune{115}}
	var p483 = charParser{id: 483, chars: []rune{101}}
	var p484 = charParser{id: 484, chars: []rune{110}}
	var p485 = charParser{id: 485, chars: []rune{100}}
	p486.items = []parser{&p482, &p483, &p484, &p485}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p487.items = []parser{&p817, &p14}
	p488.items = []parser{&p817, &p14, &p487}
	var p490 = sequenceParser{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p489.items = []parser{&p817, &p14}
	p490.items = []parser{&p817, &p14, &p489}
	p491.items = []parser{&p486, &p488, &p817, &p268, &p490, &p817, &p268}
	var p538 = sequenceParser{id: 538, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{775, 199, 785}}
	var p535 = sequenceParser{id: 535, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p533 = charParser{id: 533, chars: []rune{103}}
	var p534 = charParser{id: 534, chars: []rune{111}}
	p535.items = []parser{&p533, &p534}
	var p537 = sequenceParser{id: 537, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p536 = sequenceParser{id: 536, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p536.items = []parser{&p817, &p14}
	p537.items = []parser{&p817, &p14, &p536}
	var p258 = sequenceParser{id: 258, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p255 = sequenceParser{id: 255, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p254 = charParser{id: 254, chars: []rune{40}}
	p255.items = []parser{&p254}
	var p257 = sequenceParser{id: 257, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p256 = charParser{id: 256, chars: []rune{41}}
	p257.items = []parser{&p256}
	p258.items = []parser{&p268, &p817, &p255, &p817, &p115, &p817, &p120, &p817, &p115, &p817, &p257}
	p538.items = []parser{&p535, &p537, &p817, &p258}
	var p547 = sequenceParser{id: 547, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{775, 199, 785}}
	var p544 = sequenceParser{id: 544, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p539 = charParser{id: 539, chars: []rune{100}}
	var p540 = charParser{id: 540, chars: []rune{101}}
	var p541 = charParser{id: 541, chars: []rune{102}}
	var p542 = charParser{id: 542, chars: []rune{101}}
	var p543 = charParser{id: 543, chars: []rune{114}}
	p544.items = []parser{&p539, &p540, &p541, &p542, &p543}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p545.items = []parser{&p817, &p14}
	p546.items = []parser{&p817, &p14, &p545}
	p547.items = []parser{&p544, &p546, &p817, &p258}
	var p612 = choiceParser{id: 612, commit: 64, name: "assignment", generalizations: []int{775, 199, 785}}
	var p592 = sequenceParser{id: 592, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{612, 775, 199, 785}}
	var p589 = sequenceParser{id: 589, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p586 = charParser{id: 586, chars: []rune{115}}
	var p587 = charParser{id: 587, chars: []rune{101}}
	var p588 = charParser{id: 588, chars: []rune{116}}
	p589.items = []parser{&p586, &p587, &p588}
	var p591 = sequenceParser{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p590.items = []parser{&p817, &p14}
	p591.items = []parser{&p817, &p14, &p590}
	var p581 = sequenceParser{id: 581, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p578 = sequenceParser{id: 578, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p577 = sequenceParser{id: 577, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p576 = sequenceParser{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p576.items = []parser{&p817, &p14}
	p577.items = []parser{&p14, &p576}
	var p575 = sequenceParser{id: 575, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p574 = charParser{id: 574, chars: []rune{61}}
	p575.items = []parser{&p574}
	p578.items = []parser{&p577, &p817, &p575}
	var p580 = sequenceParser{id: 580, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p579 = sequenceParser{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p579.items = []parser{&p817, &p14}
	p580.items = []parser{&p817, &p14, &p579}
	p581.items = []parser{&p268, &p817, &p578, &p580, &p817, &p376}
	p592.items = []parser{&p589, &p591, &p817, &p581}
	var p599 = sequenceParser{id: 599, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{612, 775, 199, 785}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p595.items = []parser{&p817, &p14}
	p596.items = []parser{&p817, &p14, &p595}
	var p594 = sequenceParser{id: 594, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p593 = charParser{id: 593, chars: []rune{61}}
	p594.items = []parser{&p593}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p597.items = []parser{&p817, &p14}
	p598.items = []parser{&p817, &p14, &p597}
	p599.items = []parser{&p268, &p596, &p817, &p594, &p598, &p817, &p376}
	var p611 = sequenceParser{id: 611, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{612, 775, 199, 785}}
	var p603 = sequenceParser{id: 603, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p600 = charParser{id: 600, chars: []rune{115}}
	var p601 = charParser{id: 601, chars: []rune{101}}
	var p602 = charParser{id: 602, chars: []rune{116}}
	p603.items = []parser{&p600, &p601, &p602}
	var p610 = sequenceParser{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p609 = sequenceParser{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p609.items = []parser{&p817, &p14}
	p610.items = []parser{&p817, &p14, &p609}
	var p605 = sequenceParser{id: 605, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p604 = charParser{id: 604, chars: []rune{40}}
	p605.items = []parser{&p604}
	var p606 = sequenceParser{id: 606, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p585 = sequenceParser{id: 585, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p582 = sequenceParser{id: 582, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p582.items = []parser{&p115, &p817, &p581}
	var p583 = sequenceParser{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p583.items = []parser{&p817, &p582}
	p584.items = []parser{&p817, &p582, &p583}
	p585.items = []parser{&p581, &p584}
	p606.items = []parser{&p115, &p817, &p585}
	var p608 = sequenceParser{id: 608, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p607 = charParser{id: 607, chars: []rune{41}}
	p608.items = []parser{&p607}
	p611.items = []parser{&p603, &p610, &p817, &p605, &p817, &p606, &p817, &p115, &p817, &p608}
	p612.options = []parser{&p592, &p599, &p611}
	var p784 = sequenceParser{id: 784, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{775, 199, 785}}
	var p777 = sequenceParser{id: 777, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p776 = charParser{id: 776, chars: []rune{40}}
	p777.items = []parser{&p776}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p780 = sequenceParser{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p780.items = []parser{&p817, &p14}
	p781.items = []parser{&p817, &p14, &p780}
	var p783 = sequenceParser{id: 783, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p782 = sequenceParser{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p782.items = []parser{&p817, &p14}
	p783.items = []parser{&p817, &p14, &p782}
	var p779 = sequenceParser{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p778 = charParser{id: 778, chars: []rune{41}}
	p779.items = []parser{&p778}
	p784.items = []parser{&p777, &p781, &p817, &p775, &p783, &p817, &p779}
	p775.options = []parser{&p491, &p538, &p547, &p612, &p784, &p376}
	var p192 = sequenceParser{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{199}}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{123}}
	p189.items = []parser{&p188}
	var p191 = sequenceParser{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p190 = charParser{id: 190, chars: []rune{125}}
	p191.items = []parser{&p190}
	p192.items = []parser{&p189, &p817, &p799, &p817, &p803, &p817, &p799, &p817, &p191}
	p199.options = []parser{&p775, &p192}
	p202.items = []parser{&p194, &p817, &p115, &p817, &p196, &p817, &p115, &p817, &p198, &p201, &p817, &p199}
	p208.items = []parser{&p205, &p207, &p817, &p202}
	var p218 = sequenceParser{id: 218, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p209 = charParser{id: 209, chars: []rune{102}}
	var p210 = charParser{id: 210, chars: []rune{110}}
	p211.items = []parser{&p209, &p210}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p214.items = []parser{&p817, &p14}
	p215.items = []parser{&p817, &p14, &p214}
	var p213 = sequenceParser{id: 213, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p212 = charParser{id: 212, chars: []rune{126}}
	p213.items = []parser{&p212}
	var p217 = sequenceParser{id: 217, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p216 = sequenceParser{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p216.items = []parser{&p817, &p14}
	p217.items = []parser{&p817, &p14, &p216}
	p218.items = []parser{&p211, &p215, &p817, &p213, &p217, &p817, &p202}
	var p246 = choiceParser{id: 246, commit: 64, name: "expression-indexer", generalizations: []int{268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p236 = sequenceParser{id: 236, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{246, 268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p229 = sequenceParser{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p228 = charParser{id: 228, chars: []rune{91}}
	p229.items = []parser{&p228}
	var p233 = sequenceParser{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p232 = sequenceParser{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p232.items = []parser{&p817, &p14}
	p233.items = []parser{&p817, &p14, &p232}
	var p235 = sequenceParser{id: 235, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p234 = sequenceParser{id: 234, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p234.items = []parser{&p817, &p14}
	p235.items = []parser{&p817, &p14, &p234}
	var p231 = sequenceParser{id: 231, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p230 = charParser{id: 230, chars: []rune{93}}
	p231.items = []parser{&p230}
	p236.items = []parser{&p268, &p817, &p229, &p233, &p817, &p376, &p235, &p817, &p231}
	var p245 = sequenceParser{id: 245, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{246, 268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p238 = sequenceParser{id: 238, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p237 = charParser{id: 237, chars: []rune{91}}
	p238.items = []parser{&p237}
	var p242 = sequenceParser{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p241 = sequenceParser{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p241.items = []parser{&p817, &p14}
	p242.items = []parser{&p817, &p14, &p241}
	var p227 = sequenceParser{id: 227, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{551, 557, 558}}
	var p219 = sequenceParser{id: 219, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p219.items = []parser{&p376}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p817, &p14}
	p224.items = []parser{&p817, &p14, &p223}
	var p222 = sequenceParser{id: 222, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p221 = charParser{id: 221, chars: []rune{58}}
	p222.items = []parser{&p221}
	var p226 = sequenceParser{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p225 = sequenceParser{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p225.items = []parser{&p817, &p14}
	p226.items = []parser{&p817, &p14, &p225}
	var p220 = sequenceParser{id: 220, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p220.items = []parser{&p376}
	p227.items = []parser{&p219, &p224, &p817, &p222, &p226, &p817, &p220}
	var p244 = sequenceParser{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p243 = sequenceParser{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p243.items = []parser{&p817, &p14}
	p244.items = []parser{&p817, &p14, &p243}
	var p240 = sequenceParser{id: 240, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p239 = charParser{id: 239, chars: []rune{93}}
	p240.items = []parser{&p239}
	p245.items = []parser{&p268, &p817, &p238, &p242, &p817, &p227, &p244, &p817, &p240}
	p246.options = []parser{&p236, &p245}
	var p253 = sequenceParser{id: 253, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p249.items = []parser{&p817, &p14}
	p250.items = []parser{&p817, &p14, &p249}
	var p248 = sequenceParser{id: 248, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p247 = charParser{id: 247, chars: []rune{46}}
	p248.items = []parser{&p247}
	var p252 = sequenceParser{id: 252, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p251 = sequenceParser{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p251.items = []parser{&p817, &p14}
	p252.items = []parser{&p817, &p14, &p251}
	p253.items = []parser{&p268, &p250, &p817, &p248, &p252, &p817, &p105}
	var p267 = sequenceParser{id: 267, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{268, 376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{40}}
	p260.items = []parser{&p259}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p263.items = []parser{&p817, &p14}
	p264.items = []parser{&p817, &p14, &p263}
	var p266 = sequenceParser{id: 266, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p265 = sequenceParser{id: 265, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p265.items = []parser{&p817, &p14}
	p266.items = []parser{&p817, &p14, &p265}
	var p262 = sequenceParser{id: 262, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p261 = charParser{id: 261, chars: []rune{41}}
	p262.items = []parser{&p261}
	p267.items = []parser{&p260, &p264, &p817, &p376, &p266, &p817, &p262}
	p268.options = []parser{&p62, &p75, &p88, &p100, &p481, &p105, &p126, &p131, &p160, &p165, &p208, &p218, &p246, &p253, &p258, &p267}
	var p328 = sequenceParser{id: 328, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{376, 334, 335, 336, 337, 338, 339, 558, 551, 785}}
	var p327 = choiceParser{id: 327, commit: 66, name: "unary-operator"}
	var p287 = sequenceParser{id: 287, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p286 = charParser{id: 286, chars: []rune{43}}
	p287.items = []parser{&p286}
	var p289 = sequenceParser{id: 289, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p288 = charParser{id: 288, chars: []rune{45}}
	p289.items = []parser{&p288}
	var p270 = sequenceParser{id: 270, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p269 = charParser{id: 269, chars: []rune{94}}
	p270.items = []parser{&p269}
	var p301 = sequenceParser{id: 301, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{327}}
	var p300 = charParser{id: 300, chars: []rune{33}}
	p301.items = []parser{&p300}
	p327.options = []parser{&p287, &p289, &p270, &p301}
	p328.items = []parser{&p327, &p817, &p268}
	var p362 = choiceParser{id: 362, commit: 66, name: "binary-expression", generalizations: []int{376, 558, 551, 785}}
	var p342 = sequenceParser{id: 342, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{362, 335, 336, 337, 338, 339, 376, 558, 551, 785}}
	var p334 = choiceParser{id: 334, commit: 66, name: "operand0", generalizations: []int{335, 336, 337, 338, 339}}
	p334.options = []parser{&p268, &p328}
	var p340 = sequenceParser{id: 340, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p329 = choiceParser{id: 329, commit: 66, name: "binary-op0"}
	var p272 = sequenceParser{id: 272, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p271 = charParser{id: 271, chars: []rune{38}}
	p272.items = []parser{&p271}
	var p279 = sequenceParser{id: 279, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p277 = charParser{id: 277, chars: []rune{38}}
	var p278 = charParser{id: 278, chars: []rune{94}}
	p279.items = []parser{&p277, &p278}
	var p282 = sequenceParser{id: 282, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p280 = charParser{id: 280, chars: []rune{60}}
	var p281 = charParser{id: 281, chars: []rune{60}}
	p282.items = []parser{&p280, &p281}
	var p285 = sequenceParser{id: 285, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{329}}
	var p283 = charParser{id: 283, chars: []rune{62}}
	var p284 = charParser{id: 284, chars: []rune{62}}
	p285.items = []parser{&p283, &p284}
	var p291 = sequenceParser{id: 291, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p290 = charParser{id: 290, chars: []rune{42}}
	p291.items = []parser{&p290}
	var p293 = sequenceParser{id: 293, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p292 = charParser{id: 292, chars: []rune{47}}
	p293.items = []parser{&p292}
	var p295 = sequenceParser{id: 295, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{329}}
	var p294 = charParser{id: 294, chars: []rune{37}}
	p295.items = []parser{&p294}
	p329.options = []parser{&p272, &p279, &p282, &p285, &p291, &p293, &p295}
	p340.items = []parser{&p329, &p817, &p334}
	var p341 = sequenceParser{id: 341, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p341.items = []parser{&p817, &p340}
	p342.items = []parser{&p334, &p817, &p340, &p341}
	var p345 = sequenceParser{id: 345, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{362, 336, 337, 338, 339, 376, 558, 551, 785}}
	var p335 = choiceParser{id: 335, commit: 66, name: "operand1", generalizations: []int{336, 337, 338, 339}}
	p335.options = []parser{&p334, &p342}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p330 = choiceParser{id: 330, commit: 66, name: "binary-op1"}
	var p274 = sequenceParser{id: 274, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p273 = charParser{id: 273, chars: []rune{124}}
	p274.items = []parser{&p273}
	var p276 = sequenceParser{id: 276, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p275 = charParser{id: 275, chars: []rune{94}}
	p276.items = []parser{&p275}
	var p297 = sequenceParser{id: 297, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p296 = charParser{id: 296, chars: []rune{43}}
	p297.items = []parser{&p296}
	var p299 = sequenceParser{id: 299, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{330}}
	var p298 = charParser{id: 298, chars: []rune{45}}
	p299.items = []parser{&p298}
	p330.options = []parser{&p274, &p276, &p297, &p299}
	p343.items = []parser{&p330, &p817, &p335}
	var p344 = sequenceParser{id: 344, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p344.items = []parser{&p817, &p343}
	p345.items = []parser{&p335, &p817, &p343, &p344}
	var p348 = sequenceParser{id: 348, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{362, 337, 338, 339, 376, 558, 551, 785}}
	var p336 = choiceParser{id: 336, commit: 66, name: "operand2", generalizations: []int{337, 338, 339}}
	p336.options = []parser{&p335, &p345}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p331 = choiceParser{id: 331, commit: 66, name: "binary-op2"}
	var p304 = sequenceParser{id: 304, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{331}}
	var p302 = charParser{id: 302, chars: []rune{61}}
	var p303 = charParser{id: 303, chars: []rune{61}}
	p304.items = []parser{&p302, &p303}
	var p307 = sequenceParser{id: 307, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{331}}
	var p305 = charParser{id: 305, chars: []rune{33}}
	var p306 = charParser{id: 306, chars: []rune{61}}
	p307.items = []parser{&p305, &p306}
	var p309 = sequenceParser{id: 309, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{331}}
	var p308 = charParser{id: 308, chars: []rune{60}}
	p309.items = []parser{&p308}
	var p312 = sequenceParser{id: 312, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{331}}
	var p310 = charParser{id: 310, chars: []rune{60}}
	var p311 = charParser{id: 311, chars: []rune{61}}
	p312.items = []parser{&p310, &p311}
	var p314 = sequenceParser{id: 314, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{331}}
	var p313 = charParser{id: 313, chars: []rune{62}}
	p314.items = []parser{&p313}
	var p317 = sequenceParser{id: 317, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{331}}
	var p315 = charParser{id: 315, chars: []rune{62}}
	var p316 = charParser{id: 316, chars: []rune{61}}
	p317.items = []parser{&p315, &p316}
	p331.options = []parser{&p304, &p307, &p309, &p312, &p314, &p317}
	p346.items = []parser{&p331, &p817, &p336}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p347.items = []parser{&p817, &p346}
	p348.items = []parser{&p336, &p817, &p346, &p347}
	var p351 = sequenceParser{id: 351, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{362, 338, 339, 376, 558, 551, 785}}
	var p337 = choiceParser{id: 337, commit: 66, name: "operand3", generalizations: []int{338, 339}}
	p337.options = []parser{&p336, &p348}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p332 = sequenceParser{id: 332, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p320 = sequenceParser{id: 320, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p318 = charParser{id: 318, chars: []rune{38}}
	var p319 = charParser{id: 319, chars: []rune{38}}
	p320.items = []parser{&p318, &p319}
	p332.items = []parser{&p320}
	p349.items = []parser{&p332, &p817, &p337}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p350.items = []parser{&p817, &p349}
	p351.items = []parser{&p337, &p817, &p349, &p350}
	var p354 = sequenceParser{id: 354, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{362, 339, 376, 558, 551, 785}}
	var p338 = choiceParser{id: 338, commit: 66, name: "operand4", generalizations: []int{339}}
	p338.options = []parser{&p337, &p351}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p333 = sequenceParser{id: 333, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p323 = sequenceParser{id: 323, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p321 = charParser{id: 321, chars: []rune{124}}
	var p322 = charParser{id: 322, chars: []rune{124}}
	p323.items = []parser{&p321, &p322}
	p333.items = []parser{&p323}
	p352.items = []parser{&p333, &p817, &p338}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p353.items = []parser{&p817, &p352}
	p354.items = []parser{&p338, &p817, &p352, &p353}
	var p361 = sequenceParser{id: 361, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{362, 376, 558, 551, 785}}
	var p339 = choiceParser{id: 339, commit: 66, name: "operand5"}
	p339.options = []parser{&p338, &p354}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p355.items = []parser{&p817, &p14}
	p356.items = []parser{&p14, &p355}
	var p326 = sequenceParser{id: 326, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p324 = charParser{id: 324, chars: []rune{45}}
	var p325 = charParser{id: 325, chars: []rune{62}}
	p326.items = []parser{&p324, &p325}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p357.items = []parser{&p817, &p14}
	p358.items = []parser{&p817, &p14, &p357}
	p359.items = []parser{&p356, &p817, &p326, &p358, &p817, &p339}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p360.items = []parser{&p817, &p359}
	p361.items = []parser{&p339, &p817, &p359, &p360}
	p362.options = []parser{&p342, &p345, &p348, &p351, &p354, &p361}
	var p375 = sequenceParser{id: 375, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{376, 558, 551, 785}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p367.items = []parser{&p817, &p14}
	p368.items = []parser{&p817, &p14, &p367}
	var p364 = sequenceParser{id: 364, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p363 = charParser{id: 363, chars: []rune{63}}
	p364.items = []parser{&p363}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p369.items = []parser{&p817, &p14}
	p370.items = []parser{&p817, &p14, &p369}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p371.items = []parser{&p817, &p14}
	p372.items = []parser{&p817, &p14, &p371}
	var p366 = sequenceParser{id: 366, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p365 = charParser{id: 365, chars: []rune{58}}
	p366.items = []parser{&p365}
	var p374 = sequenceParser{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p373.items = []parser{&p817, &p14}
	p374.items = []parser{&p817, &p14, &p373}
	p375.items = []parser{&p376, &p368, &p817, &p364, &p370, &p817, &p376, &p372, &p817, &p366, &p374, &p817, &p376}
	p376.options = []parser{&p268, &p328, &p362, &p375}
	p186.items = []parser{&p185, &p817, &p376}
	p187.items = []parser{&p183, &p817, &p186}
	var p413 = sequenceParser{id: 413, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{785, 459, 523}}
	var p379 = sequenceParser{id: 379, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p377 = charParser{id: 377, chars: []rune{105}}
	var p378 = charParser{id: 378, chars: []rune{102}}
	p379.items = []parser{&p377, &p378}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p407 = sequenceParser{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p407.items = []parser{&p817, &p14}
	p408.items = []parser{&p817, &p14, &p407}
	var p410 = sequenceParser{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p409 = sequenceParser{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p409.items = []parser{&p817, &p14}
	p410.items = []parser{&p817, &p14, &p409}
	var p412 = sequenceParser{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p396 = sequenceParser{id: 396, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p388.items = []parser{&p817, &p14}
	p389.items = []parser{&p14, &p388}
	var p384 = sequenceParser{id: 384, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p380 = charParser{id: 380, chars: []rune{101}}
	var p381 = charParser{id: 381, chars: []rune{108}}
	var p382 = charParser{id: 382, chars: []rune{115}}
	var p383 = charParser{id: 383, chars: []rune{101}}
	p384.items = []parser{&p380, &p381, &p382, &p383}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p390.items = []parser{&p817, &p14}
	p391.items = []parser{&p817, &p14, &p390}
	var p387 = sequenceParser{id: 387, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p385 = charParser{id: 385, chars: []rune{105}}
	var p386 = charParser{id: 386, chars: []rune{102}}
	p387.items = []parser{&p385, &p386}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p392.items = []parser{&p817, &p14}
	p393.items = []parser{&p817, &p14, &p392}
	var p395 = sequenceParser{id: 395, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p394 = sequenceParser{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p394.items = []parser{&p817, &p14}
	p395.items = []parser{&p817, &p14, &p394}
	p396.items = []parser{&p389, &p817, &p384, &p391, &p817, &p387, &p393, &p817, &p376, &p395, &p817, &p192}
	var p411 = sequenceParser{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p411.items = []parser{&p817, &p396}
	p412.items = []parser{&p817, &p396, &p411}
	var p406 = sequenceParser{id: 406, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p403 = sequenceParser{id: 403, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p402 = sequenceParser{id: 402, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p402.items = []parser{&p817, &p14}
	p403.items = []parser{&p14, &p402}
	var p401 = sequenceParser{id: 401, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p397 = charParser{id: 397, chars: []rune{101}}
	var p398 = charParser{id: 398, chars: []rune{108}}
	var p399 = charParser{id: 399, chars: []rune{115}}
	var p400 = charParser{id: 400, chars: []rune{101}}
	p401.items = []parser{&p397, &p398, &p399, &p400}
	var p405 = sequenceParser{id: 405, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p404 = sequenceParser{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p404.items = []parser{&p817, &p14}
	p405.items = []parser{&p817, &p14, &p404}
	p406.items = []parser{&p403, &p817, &p401, &p405, &p817, &p192}
	p413.items = []parser{&p379, &p408, &p817, &p376, &p410, &p817, &p192, &p412, &p817, &p406}
	var p470 = sequenceParser{id: 470, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{459, 785, 523}}
	var p455 = sequenceParser{id: 455, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p449 = charParser{id: 449, chars: []rune{115}}
	var p450 = charParser{id: 450, chars: []rune{119}}
	var p451 = charParser{id: 451, chars: []rune{105}}
	var p452 = charParser{id: 452, chars: []rune{116}}
	var p453 = charParser{id: 453, chars: []rune{99}}
	var p454 = charParser{id: 454, chars: []rune{104}}
	p455.items = []parser{&p449, &p450, &p451, &p452, &p453, &p454}
	var p467 = sequenceParser{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p466 = sequenceParser{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p466.items = []parser{&p817, &p14}
	p467.items = []parser{&p817, &p14, &p466}
	var p469 = sequenceParser{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p468 = sequenceParser{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p468.items = []parser{&p817, &p14}
	p469.items = []parser{&p817, &p14, &p468}
	var p457 = sequenceParser{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p456 = charParser{id: 456, chars: []rune{123}}
	p457.items = []parser{&p456}
	var p463 = sequenceParser{id: 463, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p458 = choiceParser{id: 458, commit: 2}
	var p448 = sequenceParser{id: 448, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{458, 459}}
	var p443 = sequenceParser{id: 443, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p436 = sequenceParser{id: 436, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p432 = charParser{id: 432, chars: []rune{99}}
	var p433 = charParser{id: 433, chars: []rune{97}}
	var p434 = charParser{id: 434, chars: []rune{115}}
	var p435 = charParser{id: 435, chars: []rune{101}}
	p436.items = []parser{&p432, &p433, &p434, &p435}
	var p440 = sequenceParser{id: 440, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p439 = sequenceParser{id: 439, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p439.items = []parser{&p817, &p14}
	p440.items = []parser{&p817, &p14, &p439}
	var p442 = sequenceParser{id: 442, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p441 = sequenceParser{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p441.items = []parser{&p817, &p14}
	p442.items = []parser{&p817, &p14, &p441}
	var p438 = sequenceParser{id: 438, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p437 = charParser{id: 437, chars: []rune{58}}
	p438.items = []parser{&p437}
	p443.items = []parser{&p436, &p440, &p817, &p376, &p442, &p817, &p438}
	var p447 = sequenceParser{id: 447, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p445 = sequenceParser{id: 445, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p444 = charParser{id: 444, chars: []rune{59}}
	p445.items = []parser{&p444}
	var p446 = sequenceParser{id: 446, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p446.items = []parser{&p817, &p445}
	p447.items = []parser{&p817, &p445, &p446}
	p448.items = []parser{&p443, &p447, &p817, &p785}
	var p431 = sequenceParser{id: 431, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{458, 459, 522, 523}}
	var p426 = sequenceParser{id: 426, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p421 = sequenceParser{id: 421, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p414 = charParser{id: 414, chars: []rune{100}}
	var p415 = charParser{id: 415, chars: []rune{101}}
	var p416 = charParser{id: 416, chars: []rune{102}}
	var p417 = charParser{id: 417, chars: []rune{97}}
	var p418 = charParser{id: 418, chars: []rune{117}}
	var p419 = charParser{id: 419, chars: []rune{108}}
	var p420 = charParser{id: 420, chars: []rune{116}}
	p421.items = []parser{&p414, &p415, &p416, &p417, &p418, &p419, &p420}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p424.items = []parser{&p817, &p14}
	p425.items = []parser{&p817, &p14, &p424}
	var p423 = sequenceParser{id: 423, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p422 = charParser{id: 422, chars: []rune{58}}
	p423.items = []parser{&p422}
	p426.items = []parser{&p421, &p425, &p817, &p423}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p428 = sequenceParser{id: 428, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p427 = charParser{id: 427, chars: []rune{59}}
	p428.items = []parser{&p427}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p429.items = []parser{&p817, &p428}
	p430.items = []parser{&p817, &p428, &p429}
	p431.items = []parser{&p426, &p430, &p817, &p785}
	p458.options = []parser{&p448, &p431}
	var p462 = sequenceParser{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p459 = choiceParser{id: 459, commit: 2}
	p459.options = []parser{&p448, &p431, &p785}
	p460.items = []parser{&p799, &p817, &p459}
	var p461 = sequenceParser{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p461.items = []parser{&p817, &p460}
	p462.items = []parser{&p817, &p460, &p461}
	p463.items = []parser{&p458, &p462}
	var p465 = sequenceParser{id: 465, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p464 = charParser{id: 464, chars: []rune{125}}
	p465.items = []parser{&p464}
	p470.items = []parser{&p455, &p467, &p817, &p376, &p469, &p817, &p457, &p817, &p799, &p817, &p463, &p817, &p799, &p817, &p465}
	var p532 = sequenceParser{id: 532, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{523, 785}}
	var p519 = sequenceParser{id: 519, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p513 = charParser{id: 513, chars: []rune{115}}
	var p514 = charParser{id: 514, chars: []rune{101}}
	var p515 = charParser{id: 515, chars: []rune{108}}
	var p516 = charParser{id: 516, chars: []rune{101}}
	var p517 = charParser{id: 517, chars: []rune{99}}
	var p518 = charParser{id: 518, chars: []rune{116}}
	p519.items = []parser{&p513, &p514, &p515, &p516, &p517, &p518}
	var p531 = sequenceParser{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p530.items = []parser{&p817, &p14}
	p531.items = []parser{&p817, &p14, &p530}
	var p521 = sequenceParser{id: 521, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p520 = charParser{id: 520, chars: []rune{123}}
	p521.items = []parser{&p520}
	var p527 = sequenceParser{id: 527, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p522 = choiceParser{id: 522, commit: 2}
	var p512 = sequenceParser{id: 512, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{522, 523}}
	var p507 = sequenceParser{id: 507, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p500 = sequenceParser{id: 500, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p496 = charParser{id: 496, chars: []rune{99}}
	var p497 = charParser{id: 497, chars: []rune{97}}
	var p498 = charParser{id: 498, chars: []rune{115}}
	var p499 = charParser{id: 499, chars: []rune{101}}
	p500.items = []parser{&p496, &p497, &p498, &p499}
	var p504 = sequenceParser{id: 504, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p503 = sequenceParser{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p503.items = []parser{&p817, &p14}
	p504.items = []parser{&p817, &p14, &p503}
	var p495 = choiceParser{id: 495, commit: 66, name: "communication"}
	var p494 = sequenceParser{id: 494, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{495}}
	var p493 = sequenceParser{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p492.items = []parser{&p817, &p14}
	p493.items = []parser{&p817, &p14, &p492}
	p494.items = []parser{&p105, &p493, &p817, &p481}
	p495.options = []parser{&p481, &p494, &p491}
	var p506 = sequenceParser{id: 506, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p505 = sequenceParser{id: 505, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p505.items = []parser{&p817, &p14}
	p506.items = []parser{&p817, &p14, &p505}
	var p502 = sequenceParser{id: 502, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p501 = charParser{id: 501, chars: []rune{58}}
	p502.items = []parser{&p501}
	p507.items = []parser{&p500, &p504, &p817, &p495, &p506, &p817, &p502}
	var p511 = sequenceParser{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p509 = sequenceParser{id: 509, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p508 = charParser{id: 508, chars: []rune{59}}
	p509.items = []parser{&p508}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p510.items = []parser{&p817, &p509}
	p511.items = []parser{&p817, &p509, &p510}
	p512.items = []parser{&p507, &p511, &p817, &p785}
	p522.options = []parser{&p512, &p431}
	var p526 = sequenceParser{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p523 = choiceParser{id: 523, commit: 2}
	p523.options = []parser{&p512, &p431, &p785}
	p524.items = []parser{&p799, &p817, &p523}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p525.items = []parser{&p817, &p524}
	p526.items = []parser{&p817, &p524, &p525}
	p527.items = []parser{&p522, &p526}
	var p529 = sequenceParser{id: 529, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p528 = charParser{id: 528, chars: []rune{125}}
	p529.items = []parser{&p528}
	p532.items = []parser{&p519, &p531, &p817, &p521, &p817, &p799, &p817, &p527, &p817, &p799, &p817, &p529}
	var p573 = sequenceParser{id: 573, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{785}}
	var p562 = sequenceParser{id: 562, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p559 = charParser{id: 559, chars: []rune{102}}
	var p560 = charParser{id: 560, chars: []rune{111}}
	var p561 = charParser{id: 561, chars: []rune{114}}
	p562.items = []parser{&p559, &p560, &p561}
	var p572 = choiceParser{id: 572, commit: 2}
	var p568 = sequenceParser{id: 568, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{572}}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p563 = sequenceParser{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p563.items = []parser{&p817, &p14}
	p564.items = []parser{&p14, &p563}
	var p558 = choiceParser{id: 558, commit: 66, name: "loop-expression"}
	var p557 = choiceParser{id: 557, commit: 64, name: "range-over-expression", generalizations: []int{558}}
	var p556 = sequenceParser{id: 556, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{557, 558}}
	var p553 = sequenceParser{id: 553, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p552 = sequenceParser{id: 552, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p552.items = []parser{&p817, &p14}
	p553.items = []parser{&p817, &p14, &p552}
	var p550 = sequenceParser{id: 550, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p548 = charParser{id: 548, chars: []rune{105}}
	var p549 = charParser{id: 549, chars: []rune{110}}
	p550.items = []parser{&p548, &p549}
	var p555 = sequenceParser{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p554 = sequenceParser{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p554.items = []parser{&p817, &p14}
	p555.items = []parser{&p817, &p14, &p554}
	var p551 = choiceParser{id: 551, commit: 2}
	p551.options = []parser{&p376, &p227}
	p556.items = []parser{&p105, &p553, &p817, &p550, &p555, &p817, &p551}
	p557.options = []parser{&p556, &p227}
	p558.options = []parser{&p376, &p557}
	p565.items = []parser{&p564, &p817, &p558}
	var p567 = sequenceParser{id: 567, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p566 = sequenceParser{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p566.items = []parser{&p817, &p14}
	p567.items = []parser{&p817, &p14, &p566}
	p568.items = []parser{&p565, &p567, &p817, &p192}
	var p571 = sequenceParser{id: 571, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{572}}
	var p570 = sequenceParser{id: 570, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p569 = sequenceParser{id: 569, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p569.items = []parser{&p817, &p14}
	p570.items = []parser{&p14, &p569}
	p571.items = []parser{&p570, &p817, &p192}
	p572.options = []parser{&p568, &p571}
	p573.items = []parser{&p562, &p817, &p572}
	var p721 = choiceParser{id: 721, commit: 66, name: "definition", generalizations: []int{785}}
	var p634 = sequenceParser{id: 634, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{721, 785}}
	var p630 = sequenceParser{id: 630, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p627 = charParser{id: 627, chars: []rune{108}}
	var p628 = charParser{id: 628, chars: []rune{101}}
	var p629 = charParser{id: 629, chars: []rune{116}}
	p630.items = []parser{&p627, &p628, &p629}
	var p633 = sequenceParser{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p632 = sequenceParser{id: 632, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p632.items = []parser{&p817, &p14}
	p633.items = []parser{&p817, &p14, &p632}
	var p631 = choiceParser{id: 631, commit: 2}
	var p621 = sequenceParser{id: 621, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{631, 635, 636}}
	var p620 = sequenceParser{id: 620, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p615.items = []parser{&p817, &p14}
	p616.items = []parser{&p14, &p615}
	var p614 = sequenceParser{id: 614, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p613 = charParser{id: 613, chars: []rune{61}}
	p614.items = []parser{&p613}
	p617.items = []parser{&p616, &p817, &p614}
	var p619 = sequenceParser{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p618 = sequenceParser{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p618.items = []parser{&p817, &p14}
	p619.items = []parser{&p817, &p14, &p618}
	p620.items = []parser{&p105, &p817, &p617, &p619, &p817, &p376}
	p621.items = []parser{&p620}
	var p626 = sequenceParser{id: 626, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{631, 635, 636}}
	var p623 = sequenceParser{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p622 = charParser{id: 622, chars: []rune{126}}
	p623.items = []parser{&p622}
	var p625 = sequenceParser{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p624.items = []parser{&p817, &p14}
	p625.items = []parser{&p817, &p14, &p624}
	p626.items = []parser{&p623, &p625, &p817, &p620}
	p631.options = []parser{&p621, &p626}
	p634.items = []parser{&p630, &p633, &p817, &p631}
	var p655 = sequenceParser{id: 655, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{721, 785}}
	var p648 = sequenceParser{id: 648, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p645 = charParser{id: 645, chars: []rune{108}}
	var p646 = charParser{id: 646, chars: []rune{101}}
	var p647 = charParser{id: 647, chars: []rune{116}}
	p648.items = []parser{&p645, &p646, &p647}
	var p654 = sequenceParser{id: 654, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p653 = sequenceParser{id: 653, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p653.items = []parser{&p817, &p14}
	p654.items = []parser{&p817, &p14, &p653}
	var p650 = sequenceParser{id: 650, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p649 = charParser{id: 649, chars: []rune{40}}
	p650.items = []parser{&p649}
	var p640 = sequenceParser{id: 640, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p635 = choiceParser{id: 635, commit: 2}
	p635.options = []parser{&p621, &p626}
	var p639 = sequenceParser{id: 639, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p637 = sequenceParser{id: 637, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p636 = choiceParser{id: 636, commit: 2}
	p636.options = []parser{&p621, &p626}
	p637.items = []parser{&p115, &p817, &p636}
	var p638 = sequenceParser{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p638.items = []parser{&p817, &p637}
	p639.items = []parser{&p817, &p637, &p638}
	p640.items = []parser{&p635, &p639}
	var p652 = sequenceParser{id: 652, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p651 = charParser{id: 651, chars: []rune{41}}
	p652.items = []parser{&p651}
	p655.items = []parser{&p648, &p654, &p817, &p650, &p817, &p115, &p817, &p640, &p817, &p115, &p817, &p652}
	var p670 = sequenceParser{id: 670, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{721, 785}}
	var p659 = sequenceParser{id: 659, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p656 = charParser{id: 656, chars: []rune{108}}
	var p657 = charParser{id: 657, chars: []rune{101}}
	var p658 = charParser{id: 658, chars: []rune{116}}
	p659.items = []parser{&p656, &p657, &p658}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p666.items = []parser{&p817, &p14}
	p667.items = []parser{&p817, &p14, &p666}
	var p661 = sequenceParser{id: 661, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p660 = charParser{id: 660, chars: []rune{126}}
	p661.items = []parser{&p660}
	var p669 = sequenceParser{id: 669, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p668 = sequenceParser{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p668.items = []parser{&p817, &p14}
	p669.items = []parser{&p817, &p14, &p668}
	var p663 = sequenceParser{id: 663, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p662 = charParser{id: 662, chars: []rune{40}}
	p663.items = []parser{&p662}
	var p644 = sequenceParser{id: 644, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p643 = sequenceParser{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p641 = sequenceParser{id: 641, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p641.items = []parser{&p115, &p817, &p621}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p642.items = []parser{&p817, &p641}
	p643.items = []parser{&p817, &p641, &p642}
	p644.items = []parser{&p621, &p643}
	var p665 = sequenceParser{id: 665, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p664 = charParser{id: 664, chars: []rune{41}}
	p665.items = []parser{&p664}
	p670.items = []parser{&p659, &p667, &p817, &p661, &p669, &p817, &p663, &p817, &p115, &p817, &p644, &p817, &p115, &p817, &p665}
	var p686 = sequenceParser{id: 686, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{721, 785}}
	var p682 = sequenceParser{id: 682, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{102}}
	var p681 = charParser{id: 681, chars: []rune{110}}
	p682.items = []parser{&p680, &p681}
	var p685 = sequenceParser{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p684 = sequenceParser{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p684.items = []parser{&p817, &p14}
	p685.items = []parser{&p817, &p14, &p684}
	var p683 = choiceParser{id: 683, commit: 2}
	var p674 = sequenceParser{id: 674, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{683, 691, 692}}
	var p673 = sequenceParser{id: 673, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p672 = sequenceParser{id: 672, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p671 = sequenceParser{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p671.items = []parser{&p817, &p14}
	p672.items = []parser{&p817, &p14, &p671}
	p673.items = []parser{&p105, &p672, &p817, &p202}
	p674.items = []parser{&p673}
	var p679 = sequenceParser{id: 679, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{683, 691, 692}}
	var p676 = sequenceParser{id: 676, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p675 = charParser{id: 675, chars: []rune{126}}
	p676.items = []parser{&p675}
	var p678 = sequenceParser{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p677 = sequenceParser{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p677.items = []parser{&p817, &p14}
	p678.items = []parser{&p817, &p14, &p677}
	p679.items = []parser{&p676, &p678, &p817, &p673}
	p683.options = []parser{&p674, &p679}
	p686.items = []parser{&p682, &p685, &p817, &p683}
	var p706 = sequenceParser{id: 706, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{721, 785}}
	var p699 = sequenceParser{id: 699, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p697 = charParser{id: 697, chars: []rune{102}}
	var p698 = charParser{id: 698, chars: []rune{110}}
	p699.items = []parser{&p697, &p698}
	var p705 = sequenceParser{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p704.items = []parser{&p817, &p14}
	p705.items = []parser{&p817, &p14, &p704}
	var p701 = sequenceParser{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p700 = charParser{id: 700, chars: []rune{40}}
	p701.items = []parser{&p700}
	var p696 = sequenceParser{id: 696, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p691 = choiceParser{id: 691, commit: 2}
	p691.options = []parser{&p674, &p679}
	var p695 = sequenceParser{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p693 = sequenceParser{id: 693, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p692 = choiceParser{id: 692, commit: 2}
	p692.options = []parser{&p674, &p679}
	p693.items = []parser{&p115, &p817, &p692}
	var p694 = sequenceParser{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p694.items = []parser{&p817, &p693}
	p695.items = []parser{&p817, &p693, &p694}
	p696.items = []parser{&p691, &p695}
	var p703 = sequenceParser{id: 703, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p702 = charParser{id: 702, chars: []rune{41}}
	p703.items = []parser{&p702}
	p706.items = []parser{&p699, &p705, &p817, &p701, &p817, &p115, &p817, &p696, &p817, &p115, &p817, &p703}
	var p720 = sequenceParser{id: 720, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{721, 785}}
	var p709 = sequenceParser{id: 709, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p707 = charParser{id: 707, chars: []rune{102}}
	var p708 = charParser{id: 708, chars: []rune{110}}
	p709.items = []parser{&p707, &p708}
	var p717 = sequenceParser{id: 717, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p716 = sequenceParser{id: 716, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p716.items = []parser{&p817, &p14}
	p717.items = []parser{&p817, &p14, &p716}
	var p711 = sequenceParser{id: 711, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p710 = charParser{id: 710, chars: []rune{126}}
	p711.items = []parser{&p710}
	var p719 = sequenceParser{id: 719, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p718 = sequenceParser{id: 718, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p718.items = []parser{&p817, &p14}
	p719.items = []parser{&p817, &p14, &p718}
	var p713 = sequenceParser{id: 713, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p712 = charParser{id: 712, chars: []rune{40}}
	p713.items = []parser{&p712}
	var p690 = sequenceParser{id: 690, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p689 = sequenceParser{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p687.items = []parser{&p115, &p817, &p674}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p688.items = []parser{&p817, &p687}
	p689.items = []parser{&p817, &p687, &p688}
	p690.items = []parser{&p674, &p689}
	var p715 = sequenceParser{id: 715, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p714 = charParser{id: 714, chars: []rune{41}}
	p715.items = []parser{&p714}
	p720.items = []parser{&p709, &p717, &p817, &p711, &p719, &p817, &p713, &p817, &p115, &p817, &p690, &p817, &p115, &p817, &p715}
	p721.options = []parser{&p634, &p655, &p670, &p686, &p706, &p720}
	var p764 = choiceParser{id: 764, commit: 64, name: "require", generalizations: []int{785}}
	var p748 = sequenceParser{id: 748, commit: 66, name: "require-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{764, 785}}
	var p745 = sequenceParser{id: 745, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p738 = charParser{id: 738, chars: []rune{114}}
	var p739 = charParser{id: 739, chars: []rune{101}}
	var p740 = charParser{id: 740, chars: []rune{113}}
	var p741 = charParser{id: 741, chars: []rune{117}}
	var p742 = charParser{id: 742, chars: []rune{105}}
	var p743 = charParser{id: 743, chars: []rune{114}}
	var p744 = charParser{id: 744, chars: []rune{101}}
	p745.items = []parser{&p738, &p739, &p740, &p741, &p742, &p743, &p744}
	var p747 = sequenceParser{id: 747, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p746 = sequenceParser{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p746.items = []parser{&p817, &p14}
	p747.items = []parser{&p817, &p14, &p746}
	var p733 = choiceParser{id: 733, commit: 64, name: "require-fact"}
	var p732 = sequenceParser{id: 732, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{733}}
	var p724 = choiceParser{id: 724, commit: 2}
	var p723 = sequenceParser{id: 723, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{724}}
	var p722 = charParser{id: 722, chars: []rune{46}}
	p723.items = []parser{&p722}
	p724.options = []parser{&p105, &p723}
	var p729 = sequenceParser{id: 729, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p728 = sequenceParser{id: 728, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p727 = sequenceParser{id: 727, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p727.items = []parser{&p817, &p14}
	p728.items = []parser{&p14, &p727}
	var p726 = sequenceParser{id: 726, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p725 = charParser{id: 725, chars: []rune{61}}
	p726.items = []parser{&p725}
	p729.items = []parser{&p728, &p817, &p726}
	var p731 = sequenceParser{id: 731, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p730 = sequenceParser{id: 730, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p730.items = []parser{&p817, &p14}
	p731.items = []parser{&p817, &p14, &p730}
	p732.items = []parser{&p724, &p817, &p729, &p731, &p817, &p88}
	p733.options = []parser{&p88, &p732}
	p748.items = []parser{&p745, &p747, &p817, &p733}
	var p763 = sequenceParser{id: 763, commit: 66, name: "require-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{764, 785}}
	var p756 = sequenceParser{id: 756, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p749 = charParser{id: 749, chars: []rune{114}}
	var p750 = charParser{id: 750, chars: []rune{101}}
	var p751 = charParser{id: 751, chars: []rune{113}}
	var p752 = charParser{id: 752, chars: []rune{117}}
	var p753 = charParser{id: 753, chars: []rune{105}}
	var p754 = charParser{id: 754, chars: []rune{114}}
	var p755 = charParser{id: 755, chars: []rune{101}}
	p756.items = []parser{&p749, &p750, &p751, &p752, &p753, &p754, &p755}
	var p762 = sequenceParser{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p761 = sequenceParser{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p761.items = []parser{&p817, &p14}
	p762.items = []parser{&p817, &p14, &p761}
	var p758 = sequenceParser{id: 758, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p757 = charParser{id: 757, chars: []rune{40}}
	p758.items = []parser{&p757}
	var p737 = sequenceParser{id: 737, commit: 66, name: "require-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p734 = sequenceParser{id: 734, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p734.items = []parser{&p115, &p817, &p733}
	var p735 = sequenceParser{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p735.items = []parser{&p817, &p734}
	p736.items = []parser{&p817, &p734, &p735}
	p737.items = []parser{&p733, &p736}
	var p760 = sequenceParser{id: 760, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p759 = charParser{id: 759, chars: []rune{41}}
	p760.items = []parser{&p759}
	p763.items = []parser{&p756, &p762, &p817, &p758, &p817, &p115, &p817, &p737, &p817, &p115, &p817, &p760}
	p764.options = []parser{&p748, &p763}
	var p774 = sequenceParser{id: 774, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785}}
	var p771 = sequenceParser{id: 771, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p765 = charParser{id: 765, chars: []rune{101}}
	var p766 = charParser{id: 766, chars: []rune{120}}
	var p767 = charParser{id: 767, chars: []rune{112}}
	var p768 = charParser{id: 768, chars: []rune{111}}
	var p769 = charParser{id: 769, chars: []rune{114}}
	var p770 = charParser{id: 770, chars: []rune{116}}
	p771.items = []parser{&p765, &p766, &p767, &p768, &p769, &p770}
	var p773 = sequenceParser{id: 773, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p772 = sequenceParser{id: 772, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p772.items = []parser{&p817, &p14}
	p773.items = []parser{&p817, &p14, &p772}
	p774.items = []parser{&p771, &p773, &p817, &p721}
	var p794 = sequenceParser{id: 794, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785}}
	var p787 = sequenceParser{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p786 = charParser{id: 786, chars: []rune{40}}
	p787.items = []parser{&p786}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p790.items = []parser{&p817, &p14}
	p791.items = []parser{&p817, &p14, &p790}
	var p793 = sequenceParser{id: 793, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p792.items = []parser{&p817, &p14}
	p793.items = []parser{&p817, &p14, &p792}
	var p789 = sequenceParser{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p788 = charParser{id: 788, chars: []rune{41}}
	p789.items = []parser{&p788}
	p794.items = []parser{&p787, &p791, &p817, &p785, &p793, &p817, &p789}
	p785.options = []parser{&p187, &p413, &p470, &p532, &p573, &p721, &p764, &p774, &p794, &p775}
	var p802 = sequenceParser{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p800.items = []parser{&p799, &p817, &p785}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p801.items = []parser{&p817, &p800}
	p802.items = []parser{&p817, &p800, &p801}
	p803.items = []parser{&p785, &p802}
	p818.items = []parser{&p814, &p817, &p799, &p817, &p803, &p817, &p799}
	p819.items = []parser{&p817, &p818, &p817}
	var b819 = sequenceBuilder{id: 819, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b817 = choiceBuilder{id: 817, commit: 2}
	var b815 = choiceBuilder{id: 815, commit: 70}
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
	b815.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b816 = sequenceBuilder{id: 816, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b41.items = []builder{&b40, &b817, &b38}
	var b42 = sequenceBuilder{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b42.items = []builder{&b817, &b41}
	b43.items = []builder{&b817, &b41, &b42}
	b44.items = []builder{&b38, &b43}
	b816.items = []builder{&b44}
	b817.options = []builder{&b815, &b816}
	var b818 = sequenceBuilder{id: 818, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b814 = sequenceBuilder{id: 814, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b811 = sequenceBuilder{id: 811, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b809 = charBuilder{}
	var b810 = charBuilder{}
	b811.items = []builder{&b809, &b810}
	var b808 = sequenceBuilder{id: 808, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b807 = sequenceBuilder{id: 807, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b805 = sequenceBuilder{id: 805, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b804 = charBuilder{}
	b805.items = []builder{&b804}
	var b806 = sequenceBuilder{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b806.items = []builder{&b817, &b805}
	b807.items = []builder{&b805, &b806}
	b808.items = []builder{&b807}
	var b813 = sequenceBuilder{id: 813, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b812 = charBuilder{}
	b813.items = []builder{&b812}
	b814.items = []builder{&b811, &b817, &b808, &b817, &b813}
	var b799 = sequenceBuilder{id: 799, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b797 = choiceBuilder{id: 797, commit: 2}
	var b796 = sequenceBuilder{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b795 = charBuilder{}
	b796.items = []builder{&b795}
	var b14 = sequenceBuilder{id: 14, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b797.options = []builder{&b796, &b14}
	var b798 = sequenceBuilder{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b798.items = []builder{&b817, &b797}
	b799.items = []builder{&b797, &b798}
	var b803 = sequenceBuilder{id: 803, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b785 = choiceBuilder{id: 785, commit: 66}
	var b187 = sequenceBuilder{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b183 = sequenceBuilder{id: 183, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b177 = charBuilder{}
	var b178 = charBuilder{}
	var b179 = charBuilder{}
	var b180 = charBuilder{}
	var b181 = charBuilder{}
	var b182 = charBuilder{}
	b183.items = []builder{&b177, &b178, &b179, &b180, &b181, &b182}
	var b186 = sequenceBuilder{id: 186, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b185 = sequenceBuilder{id: 185, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b184 = sequenceBuilder{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b184.items = []builder{&b817, &b14}
	b185.items = []builder{&b14, &b184}
	var b376 = choiceBuilder{id: 376, commit: 66}
	var b268 = choiceBuilder{id: 268, commit: 66}
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
	var b481 = sequenceBuilder{id: 481, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b478 = sequenceBuilder{id: 478, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b471 = charBuilder{}
	var b472 = charBuilder{}
	var b473 = charBuilder{}
	var b474 = charBuilder{}
	var b475 = charBuilder{}
	var b476 = charBuilder{}
	var b477 = charBuilder{}
	b478.items = []builder{&b471, &b472, &b473, &b474, &b475, &b476, &b477}
	var b480 = sequenceBuilder{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b479 = sequenceBuilder{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b479.items = []builder{&b817, &b14}
	b480.items = []builder{&b817, &b14, &b479}
	b481.items = []builder{&b478, &b480, &b817, &b268}
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
	b114.items = []builder{&b817, &b113}
	b115.items = []builder{&b113, &b114}
	var b120 = sequenceBuilder{id: 120, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b116 = choiceBuilder{id: 116, commit: 66}
	var b110 = sequenceBuilder{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b109 = sequenceBuilder{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	var b108 = charBuilder{}
	b109.items = []builder{&b106, &b107, &b108}
	b110.items = []builder{&b268, &b817, &b109}
	b116.options = []builder{&b376, &b110}
	var b119 = sequenceBuilder{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b117.items = []builder{&b115, &b817, &b116}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b118.items = []builder{&b817, &b117}
	b119.items = []builder{&b817, &b117, &b118}
	b120.items = []builder{&b116, &b119}
	var b124 = sequenceBuilder{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b123 = charBuilder{}
	b124.items = []builder{&b123}
	b125.items = []builder{&b122, &b817, &b115, &b817, &b120, &b817, &b115, &b817, &b124}
	b126.items = []builder{&b125}
	var b131 = sequenceBuilder{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b128 = sequenceBuilder{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b127 = charBuilder{}
	b128.items = []builder{&b127}
	var b130 = sequenceBuilder{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b129.items = []builder{&b817, &b14}
	b130.items = []builder{&b817, &b14, &b129}
	b131.items = []builder{&b128, &b130, &b817, &b125}
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
	b136.items = []builder{&b817, &b14}
	b137.items = []builder{&b817, &b14, &b136}
	var b139 = sequenceBuilder{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b138.items = []builder{&b817, &b14}
	b139.items = []builder{&b817, &b14, &b138}
	var b135 = sequenceBuilder{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b134 = charBuilder{}
	b135.items = []builder{&b134}
	b140.items = []builder{&b133, &b137, &b817, &b376, &b139, &b817, &b135}
	b141.options = []builder{&b105, &b88, &b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b817, &b14}
	b145.items = []builder{&b817, &b14, &b144}
	var b143 = sequenceBuilder{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b142 = charBuilder{}
	b143.items = []builder{&b142}
	var b147 = sequenceBuilder{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b146.items = []builder{&b817, &b14}
	b147.items = []builder{&b817, &b14, &b146}
	b148.items = []builder{&b141, &b145, &b817, &b143, &b147, &b817, &b376}
	b149.options = []builder{&b148, &b110}
	var b153 = sequenceBuilder{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b150 = choiceBuilder{id: 150, commit: 2}
	b150.options = []builder{&b148, &b110}
	b151.items = []builder{&b115, &b817, &b150}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b152.items = []builder{&b817, &b151}
	b153.items = []builder{&b817, &b151, &b152}
	b154.items = []builder{&b149, &b153}
	var b158 = sequenceBuilder{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b157 = charBuilder{}
	b158.items = []builder{&b157}
	b159.items = []builder{&b156, &b817, &b115, &b817, &b154, &b817, &b115, &b817, &b158}
	b160.items = []builder{&b159}
	var b165 = sequenceBuilder{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b162 = sequenceBuilder{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b161 = charBuilder{}
	b162.items = []builder{&b161}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b163.items = []builder{&b817, &b14}
	b164.items = []builder{&b817, &b14, &b163}
	b165.items = []builder{&b162, &b164, &b817, &b159}
	var b208 = sequenceBuilder{id: 208, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b205 = sequenceBuilder{id: 205, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b203 = charBuilder{}
	var b204 = charBuilder{}
	b205.items = []builder{&b203, &b204}
	var b207 = sequenceBuilder{id: 207, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b206 = sequenceBuilder{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b206.items = []builder{&b817, &b14}
	b207.items = []builder{&b817, &b14, &b206}
	var b202 = sequenceBuilder{id: 202, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b194 = sequenceBuilder{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b193 = charBuilder{}
	b194.items = []builder{&b193}
	var b196 = choiceBuilder{id: 196, commit: 2}
	var b169 = sequenceBuilder{id: 169, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b168 = sequenceBuilder{id: 168, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b166.items = []builder{&b115, &b817, &b105}
	var b167 = sequenceBuilder{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b167.items = []builder{&b817, &b166}
	b168.items = []builder{&b817, &b166, &b167}
	b169.items = []builder{&b105, &b168}
	var b195 = sequenceBuilder{id: 195, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b176 = sequenceBuilder{id: 176, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b173 = sequenceBuilder{id: 173, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b170 = charBuilder{}
	var b171 = charBuilder{}
	var b172 = charBuilder{}
	b173.items = []builder{&b170, &b171, &b172}
	var b175 = sequenceBuilder{id: 175, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b174 = sequenceBuilder{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b174.items = []builder{&b817, &b14}
	b175.items = []builder{&b817, &b14, &b174}
	b176.items = []builder{&b173, &b175, &b817, &b105}
	b195.items = []builder{&b169, &b817, &b115, &b817, &b176}
	b196.options = []builder{&b169, &b195, &b176}
	var b198 = sequenceBuilder{id: 198, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b197 = charBuilder{}
	b198.items = []builder{&b197}
	var b201 = sequenceBuilder{id: 201, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b200 = sequenceBuilder{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b200.items = []builder{&b817, &b14}
	b201.items = []builder{&b817, &b14, &b200}
	var b199 = choiceBuilder{id: 199, commit: 2}
	var b775 = choiceBuilder{id: 775, commit: 66}
	var b491 = sequenceBuilder{id: 491, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b486 = sequenceBuilder{id: 486, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b482 = charBuilder{}
	var b483 = charBuilder{}
	var b484 = charBuilder{}
	var b485 = charBuilder{}
	b486.items = []builder{&b482, &b483, &b484, &b485}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b487.items = []builder{&b817, &b14}
	b488.items = []builder{&b817, &b14, &b487}
	var b490 = sequenceBuilder{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b489.items = []builder{&b817, &b14}
	b490.items = []builder{&b817, &b14, &b489}
	b491.items = []builder{&b486, &b488, &b817, &b268, &b490, &b817, &b268}
	var b538 = sequenceBuilder{id: 538, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b535 = sequenceBuilder{id: 535, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b533 = charBuilder{}
	var b534 = charBuilder{}
	b535.items = []builder{&b533, &b534}
	var b537 = sequenceBuilder{id: 537, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b536 = sequenceBuilder{id: 536, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b536.items = []builder{&b817, &b14}
	b537.items = []builder{&b817, &b14, &b536}
	var b258 = sequenceBuilder{id: 258, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b255 = sequenceBuilder{id: 255, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b254 = charBuilder{}
	b255.items = []builder{&b254}
	var b257 = sequenceBuilder{id: 257, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b256 = charBuilder{}
	b257.items = []builder{&b256}
	b258.items = []builder{&b268, &b817, &b255, &b817, &b115, &b817, &b120, &b817, &b115, &b817, &b257}
	b538.items = []builder{&b535, &b537, &b817, &b258}
	var b547 = sequenceBuilder{id: 547, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b544 = sequenceBuilder{id: 544, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b539 = charBuilder{}
	var b540 = charBuilder{}
	var b541 = charBuilder{}
	var b542 = charBuilder{}
	var b543 = charBuilder{}
	b544.items = []builder{&b539, &b540, &b541, &b542, &b543}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b545.items = []builder{&b817, &b14}
	b546.items = []builder{&b817, &b14, &b545}
	b547.items = []builder{&b544, &b546, &b817, &b258}
	var b612 = choiceBuilder{id: 612, commit: 64, name: "assignment"}
	var b592 = sequenceBuilder{id: 592, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b589 = sequenceBuilder{id: 589, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b586 = charBuilder{}
	var b587 = charBuilder{}
	var b588 = charBuilder{}
	b589.items = []builder{&b586, &b587, &b588}
	var b591 = sequenceBuilder{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b590.items = []builder{&b817, &b14}
	b591.items = []builder{&b817, &b14, &b590}
	var b581 = sequenceBuilder{id: 581, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b578 = sequenceBuilder{id: 578, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b577 = sequenceBuilder{id: 577, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b576 = sequenceBuilder{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b576.items = []builder{&b817, &b14}
	b577.items = []builder{&b14, &b576}
	var b575 = sequenceBuilder{id: 575, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b574 = charBuilder{}
	b575.items = []builder{&b574}
	b578.items = []builder{&b577, &b817, &b575}
	var b580 = sequenceBuilder{id: 580, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b579 = sequenceBuilder{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b579.items = []builder{&b817, &b14}
	b580.items = []builder{&b817, &b14, &b579}
	b581.items = []builder{&b268, &b817, &b578, &b580, &b817, &b376}
	b592.items = []builder{&b589, &b591, &b817, &b581}
	var b599 = sequenceBuilder{id: 599, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b595.items = []builder{&b817, &b14}
	b596.items = []builder{&b817, &b14, &b595}
	var b594 = sequenceBuilder{id: 594, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b593 = charBuilder{}
	b594.items = []builder{&b593}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b597.items = []builder{&b817, &b14}
	b598.items = []builder{&b817, &b14, &b597}
	b599.items = []builder{&b268, &b596, &b817, &b594, &b598, &b817, &b376}
	var b611 = sequenceBuilder{id: 611, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b603 = sequenceBuilder{id: 603, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b600 = charBuilder{}
	var b601 = charBuilder{}
	var b602 = charBuilder{}
	b603.items = []builder{&b600, &b601, &b602}
	var b610 = sequenceBuilder{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b609 = sequenceBuilder{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b609.items = []builder{&b817, &b14}
	b610.items = []builder{&b817, &b14, &b609}
	var b605 = sequenceBuilder{id: 605, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b604 = charBuilder{}
	b605.items = []builder{&b604}
	var b606 = sequenceBuilder{id: 606, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b585 = sequenceBuilder{id: 585, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b582 = sequenceBuilder{id: 582, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b582.items = []builder{&b115, &b817, &b581}
	var b583 = sequenceBuilder{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b583.items = []builder{&b817, &b582}
	b584.items = []builder{&b817, &b582, &b583}
	b585.items = []builder{&b581, &b584}
	b606.items = []builder{&b115, &b817, &b585}
	var b608 = sequenceBuilder{id: 608, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b607 = charBuilder{}
	b608.items = []builder{&b607}
	b611.items = []builder{&b603, &b610, &b817, &b605, &b817, &b606, &b817, &b115, &b817, &b608}
	b612.options = []builder{&b592, &b599, &b611}
	var b784 = sequenceBuilder{id: 784, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b777 = sequenceBuilder{id: 777, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b776 = charBuilder{}
	b777.items = []builder{&b776}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b780 = sequenceBuilder{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b780.items = []builder{&b817, &b14}
	b781.items = []builder{&b817, &b14, &b780}
	var b783 = sequenceBuilder{id: 783, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b782 = sequenceBuilder{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b782.items = []builder{&b817, &b14}
	b783.items = []builder{&b817, &b14, &b782}
	var b779 = sequenceBuilder{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b778 = charBuilder{}
	b779.items = []builder{&b778}
	b784.items = []builder{&b777, &b781, &b817, &b775, &b783, &b817, &b779}
	b775.options = []builder{&b491, &b538, &b547, &b612, &b784, &b376}
	var b192 = sequenceBuilder{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	var b191 = sequenceBuilder{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b190 = charBuilder{}
	b191.items = []builder{&b190}
	b192.items = []builder{&b189, &b817, &b799, &b817, &b803, &b817, &b799, &b817, &b191}
	b199.options = []builder{&b775, &b192}
	b202.items = []builder{&b194, &b817, &b115, &b817, &b196, &b817, &b115, &b817, &b198, &b201, &b817, &b199}
	b208.items = []builder{&b205, &b207, &b817, &b202}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b209 = charBuilder{}
	var b210 = charBuilder{}
	b211.items = []builder{&b209, &b210}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b214.items = []builder{&b817, &b14}
	b215.items = []builder{&b817, &b14, &b214}
	var b213 = sequenceBuilder{id: 213, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b212 = charBuilder{}
	b213.items = []builder{&b212}
	var b217 = sequenceBuilder{id: 217, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b216 = sequenceBuilder{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b216.items = []builder{&b817, &b14}
	b217.items = []builder{&b817, &b14, &b216}
	b218.items = []builder{&b211, &b215, &b817, &b213, &b217, &b817, &b202}
	var b246 = choiceBuilder{id: 246, commit: 64, name: "expression-indexer"}
	var b236 = sequenceBuilder{id: 236, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b229 = sequenceBuilder{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b228 = charBuilder{}
	b229.items = []builder{&b228}
	var b233 = sequenceBuilder{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b232 = sequenceBuilder{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b232.items = []builder{&b817, &b14}
	b233.items = []builder{&b817, &b14, &b232}
	var b235 = sequenceBuilder{id: 235, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b234 = sequenceBuilder{id: 234, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b234.items = []builder{&b817, &b14}
	b235.items = []builder{&b817, &b14, &b234}
	var b231 = sequenceBuilder{id: 231, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b230 = charBuilder{}
	b231.items = []builder{&b230}
	b236.items = []builder{&b268, &b817, &b229, &b233, &b817, &b376, &b235, &b817, &b231}
	var b245 = sequenceBuilder{id: 245, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b238 = sequenceBuilder{id: 238, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b237 = charBuilder{}
	b238.items = []builder{&b237}
	var b242 = sequenceBuilder{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b241 = sequenceBuilder{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b241.items = []builder{&b817, &b14}
	b242.items = []builder{&b817, &b14, &b241}
	var b227 = sequenceBuilder{id: 227, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b219 = sequenceBuilder{id: 219, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b219.items = []builder{&b376}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b817, &b14}
	b224.items = []builder{&b817, &b14, &b223}
	var b222 = sequenceBuilder{id: 222, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b221 = charBuilder{}
	b222.items = []builder{&b221}
	var b226 = sequenceBuilder{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b225 = sequenceBuilder{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b225.items = []builder{&b817, &b14}
	b226.items = []builder{&b817, &b14, &b225}
	var b220 = sequenceBuilder{id: 220, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b220.items = []builder{&b376}
	b227.items = []builder{&b219, &b224, &b817, &b222, &b226, &b817, &b220}
	var b244 = sequenceBuilder{id: 244, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b243 = sequenceBuilder{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b243.items = []builder{&b817, &b14}
	b244.items = []builder{&b817, &b14, &b243}
	var b240 = sequenceBuilder{id: 240, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b239 = charBuilder{}
	b240.items = []builder{&b239}
	b245.items = []builder{&b268, &b817, &b238, &b242, &b817, &b227, &b244, &b817, &b240}
	b246.options = []builder{&b236, &b245}
	var b253 = sequenceBuilder{id: 253, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b249.items = []builder{&b817, &b14}
	b250.items = []builder{&b817, &b14, &b249}
	var b248 = sequenceBuilder{id: 248, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b247 = charBuilder{}
	b248.items = []builder{&b247}
	var b252 = sequenceBuilder{id: 252, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b251 = sequenceBuilder{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b251.items = []builder{&b817, &b14}
	b252.items = []builder{&b817, &b14, &b251}
	b253.items = []builder{&b268, &b250, &b817, &b248, &b252, &b817, &b105}
	var b267 = sequenceBuilder{id: 267, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b263.items = []builder{&b817, &b14}
	b264.items = []builder{&b817, &b14, &b263}
	var b266 = sequenceBuilder{id: 266, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b265 = sequenceBuilder{id: 265, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b265.items = []builder{&b817, &b14}
	b266.items = []builder{&b817, &b14, &b265}
	var b262 = sequenceBuilder{id: 262, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b261 = charBuilder{}
	b262.items = []builder{&b261}
	b267.items = []builder{&b260, &b264, &b817, &b376, &b266, &b817, &b262}
	b268.options = []builder{&b62, &b75, &b88, &b100, &b481, &b105, &b126, &b131, &b160, &b165, &b208, &b218, &b246, &b253, &b258, &b267}
	var b328 = sequenceBuilder{id: 328, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b327 = choiceBuilder{id: 327, commit: 66}
	var b287 = sequenceBuilder{id: 287, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b286 = charBuilder{}
	b287.items = []builder{&b286}
	var b289 = sequenceBuilder{id: 289, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b288 = charBuilder{}
	b289.items = []builder{&b288}
	var b270 = sequenceBuilder{id: 270, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b269 = charBuilder{}
	b270.items = []builder{&b269}
	var b301 = sequenceBuilder{id: 301, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b300 = charBuilder{}
	b301.items = []builder{&b300}
	b327.options = []builder{&b287, &b289, &b270, &b301}
	b328.items = []builder{&b327, &b817, &b268}
	var b362 = choiceBuilder{id: 362, commit: 66}
	var b342 = sequenceBuilder{id: 342, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b334 = choiceBuilder{id: 334, commit: 66}
	b334.options = []builder{&b268, &b328}
	var b340 = sequenceBuilder{id: 340, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b329 = choiceBuilder{id: 329, commit: 66}
	var b272 = sequenceBuilder{id: 272, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b271 = charBuilder{}
	b272.items = []builder{&b271}
	var b279 = sequenceBuilder{id: 279, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b277 = charBuilder{}
	var b278 = charBuilder{}
	b279.items = []builder{&b277, &b278}
	var b282 = sequenceBuilder{id: 282, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b280 = charBuilder{}
	var b281 = charBuilder{}
	b282.items = []builder{&b280, &b281}
	var b285 = sequenceBuilder{id: 285, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b283 = charBuilder{}
	var b284 = charBuilder{}
	b285.items = []builder{&b283, &b284}
	var b291 = sequenceBuilder{id: 291, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b290 = charBuilder{}
	b291.items = []builder{&b290}
	var b293 = sequenceBuilder{id: 293, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b292 = charBuilder{}
	b293.items = []builder{&b292}
	var b295 = sequenceBuilder{id: 295, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b294 = charBuilder{}
	b295.items = []builder{&b294}
	b329.options = []builder{&b272, &b279, &b282, &b285, &b291, &b293, &b295}
	b340.items = []builder{&b329, &b817, &b334}
	var b341 = sequenceBuilder{id: 341, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b341.items = []builder{&b817, &b340}
	b342.items = []builder{&b334, &b817, &b340, &b341}
	var b345 = sequenceBuilder{id: 345, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b335 = choiceBuilder{id: 335, commit: 66}
	b335.options = []builder{&b334, &b342}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b330 = choiceBuilder{id: 330, commit: 66}
	var b274 = sequenceBuilder{id: 274, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b273 = charBuilder{}
	b274.items = []builder{&b273}
	var b276 = sequenceBuilder{id: 276, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b275 = charBuilder{}
	b276.items = []builder{&b275}
	var b297 = sequenceBuilder{id: 297, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b296 = charBuilder{}
	b297.items = []builder{&b296}
	var b299 = sequenceBuilder{id: 299, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b298 = charBuilder{}
	b299.items = []builder{&b298}
	b330.options = []builder{&b274, &b276, &b297, &b299}
	b343.items = []builder{&b330, &b817, &b335}
	var b344 = sequenceBuilder{id: 344, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b344.items = []builder{&b817, &b343}
	b345.items = []builder{&b335, &b817, &b343, &b344}
	var b348 = sequenceBuilder{id: 348, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b336 = choiceBuilder{id: 336, commit: 66}
	b336.options = []builder{&b335, &b345}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b331 = choiceBuilder{id: 331, commit: 66}
	var b304 = sequenceBuilder{id: 304, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b302 = charBuilder{}
	var b303 = charBuilder{}
	b304.items = []builder{&b302, &b303}
	var b307 = sequenceBuilder{id: 307, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b305 = charBuilder{}
	var b306 = charBuilder{}
	b307.items = []builder{&b305, &b306}
	var b309 = sequenceBuilder{id: 309, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b308 = charBuilder{}
	b309.items = []builder{&b308}
	var b312 = sequenceBuilder{id: 312, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b310 = charBuilder{}
	var b311 = charBuilder{}
	b312.items = []builder{&b310, &b311}
	var b314 = sequenceBuilder{id: 314, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b313 = charBuilder{}
	b314.items = []builder{&b313}
	var b317 = sequenceBuilder{id: 317, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b315 = charBuilder{}
	var b316 = charBuilder{}
	b317.items = []builder{&b315, &b316}
	b331.options = []builder{&b304, &b307, &b309, &b312, &b314, &b317}
	b346.items = []builder{&b331, &b817, &b336}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b347.items = []builder{&b817, &b346}
	b348.items = []builder{&b336, &b817, &b346, &b347}
	var b351 = sequenceBuilder{id: 351, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b337 = choiceBuilder{id: 337, commit: 66}
	b337.options = []builder{&b336, &b348}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b332 = sequenceBuilder{id: 332, commit: 66, ranges: [][]int{{1, 1}}}
	var b320 = sequenceBuilder{id: 320, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b318 = charBuilder{}
	var b319 = charBuilder{}
	b320.items = []builder{&b318, &b319}
	b332.items = []builder{&b320}
	b349.items = []builder{&b332, &b817, &b337}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b350.items = []builder{&b817, &b349}
	b351.items = []builder{&b337, &b817, &b349, &b350}
	var b354 = sequenceBuilder{id: 354, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b338 = choiceBuilder{id: 338, commit: 66}
	b338.options = []builder{&b337, &b351}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b333 = sequenceBuilder{id: 333, commit: 66, ranges: [][]int{{1, 1}}}
	var b323 = sequenceBuilder{id: 323, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b321 = charBuilder{}
	var b322 = charBuilder{}
	b323.items = []builder{&b321, &b322}
	b333.items = []builder{&b323}
	b352.items = []builder{&b333, &b817, &b338}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b353.items = []builder{&b817, &b352}
	b354.items = []builder{&b338, &b817, &b352, &b353}
	var b361 = sequenceBuilder{id: 361, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b339 = choiceBuilder{id: 339, commit: 66}
	b339.options = []builder{&b338, &b354}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b355.items = []builder{&b817, &b14}
	b356.items = []builder{&b14, &b355}
	var b326 = sequenceBuilder{id: 326, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b324 = charBuilder{}
	var b325 = charBuilder{}
	b326.items = []builder{&b324, &b325}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b357.items = []builder{&b817, &b14}
	b358.items = []builder{&b817, &b14, &b357}
	b359.items = []builder{&b356, &b817, &b326, &b358, &b817, &b339}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b360.items = []builder{&b817, &b359}
	b361.items = []builder{&b339, &b817, &b359, &b360}
	b362.options = []builder{&b342, &b345, &b348, &b351, &b354, &b361}
	var b375 = sequenceBuilder{id: 375, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b367.items = []builder{&b817, &b14}
	b368.items = []builder{&b817, &b14, &b367}
	var b364 = sequenceBuilder{id: 364, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b363 = charBuilder{}
	b364.items = []builder{&b363}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b369.items = []builder{&b817, &b14}
	b370.items = []builder{&b817, &b14, &b369}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b371.items = []builder{&b817, &b14}
	b372.items = []builder{&b817, &b14, &b371}
	var b366 = sequenceBuilder{id: 366, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b365 = charBuilder{}
	b366.items = []builder{&b365}
	var b374 = sequenceBuilder{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b373.items = []builder{&b817, &b14}
	b374.items = []builder{&b817, &b14, &b373}
	b375.items = []builder{&b376, &b368, &b817, &b364, &b370, &b817, &b376, &b372, &b817, &b366, &b374, &b817, &b376}
	b376.options = []builder{&b268, &b328, &b362, &b375}
	b186.items = []builder{&b185, &b817, &b376}
	b187.items = []builder{&b183, &b817, &b186}
	var b413 = sequenceBuilder{id: 413, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b379 = sequenceBuilder{id: 379, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b377 = charBuilder{}
	var b378 = charBuilder{}
	b379.items = []builder{&b377, &b378}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b407 = sequenceBuilder{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b407.items = []builder{&b817, &b14}
	b408.items = []builder{&b817, &b14, &b407}
	var b410 = sequenceBuilder{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b409 = sequenceBuilder{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b409.items = []builder{&b817, &b14}
	b410.items = []builder{&b817, &b14, &b409}
	var b412 = sequenceBuilder{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b396 = sequenceBuilder{id: 396, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b388.items = []builder{&b817, &b14}
	b389.items = []builder{&b14, &b388}
	var b384 = sequenceBuilder{id: 384, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b380 = charBuilder{}
	var b381 = charBuilder{}
	var b382 = charBuilder{}
	var b383 = charBuilder{}
	b384.items = []builder{&b380, &b381, &b382, &b383}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b390.items = []builder{&b817, &b14}
	b391.items = []builder{&b817, &b14, &b390}
	var b387 = sequenceBuilder{id: 387, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b385 = charBuilder{}
	var b386 = charBuilder{}
	b387.items = []builder{&b385, &b386}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b392.items = []builder{&b817, &b14}
	b393.items = []builder{&b817, &b14, &b392}
	var b395 = sequenceBuilder{id: 395, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b394 = sequenceBuilder{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b394.items = []builder{&b817, &b14}
	b395.items = []builder{&b817, &b14, &b394}
	b396.items = []builder{&b389, &b817, &b384, &b391, &b817, &b387, &b393, &b817, &b376, &b395, &b817, &b192}
	var b411 = sequenceBuilder{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b411.items = []builder{&b817, &b396}
	b412.items = []builder{&b817, &b396, &b411}
	var b406 = sequenceBuilder{id: 406, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b403 = sequenceBuilder{id: 403, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b402 = sequenceBuilder{id: 402, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b402.items = []builder{&b817, &b14}
	b403.items = []builder{&b14, &b402}
	var b401 = sequenceBuilder{id: 401, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b397 = charBuilder{}
	var b398 = charBuilder{}
	var b399 = charBuilder{}
	var b400 = charBuilder{}
	b401.items = []builder{&b397, &b398, &b399, &b400}
	var b405 = sequenceBuilder{id: 405, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b404 = sequenceBuilder{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b404.items = []builder{&b817, &b14}
	b405.items = []builder{&b817, &b14, &b404}
	b406.items = []builder{&b403, &b817, &b401, &b405, &b817, &b192}
	b413.items = []builder{&b379, &b408, &b817, &b376, &b410, &b817, &b192, &b412, &b817, &b406}
	var b470 = sequenceBuilder{id: 470, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b455 = sequenceBuilder{id: 455, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b449 = charBuilder{}
	var b450 = charBuilder{}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	var b454 = charBuilder{}
	b455.items = []builder{&b449, &b450, &b451, &b452, &b453, &b454}
	var b467 = sequenceBuilder{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b466 = sequenceBuilder{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b466.items = []builder{&b817, &b14}
	b467.items = []builder{&b817, &b14, &b466}
	var b469 = sequenceBuilder{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b468 = sequenceBuilder{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b468.items = []builder{&b817, &b14}
	b469.items = []builder{&b817, &b14, &b468}
	var b457 = sequenceBuilder{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b456 = charBuilder{}
	b457.items = []builder{&b456}
	var b463 = sequenceBuilder{id: 463, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b458 = choiceBuilder{id: 458, commit: 2}
	var b448 = sequenceBuilder{id: 448, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b443 = sequenceBuilder{id: 443, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b436 = sequenceBuilder{id: 436, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b432 = charBuilder{}
	var b433 = charBuilder{}
	var b434 = charBuilder{}
	var b435 = charBuilder{}
	b436.items = []builder{&b432, &b433, &b434, &b435}
	var b440 = sequenceBuilder{id: 440, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b439 = sequenceBuilder{id: 439, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b439.items = []builder{&b817, &b14}
	b440.items = []builder{&b817, &b14, &b439}
	var b442 = sequenceBuilder{id: 442, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b441 = sequenceBuilder{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b441.items = []builder{&b817, &b14}
	b442.items = []builder{&b817, &b14, &b441}
	var b438 = sequenceBuilder{id: 438, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b437 = charBuilder{}
	b438.items = []builder{&b437}
	b443.items = []builder{&b436, &b440, &b817, &b376, &b442, &b817, &b438}
	var b447 = sequenceBuilder{id: 447, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b445 = sequenceBuilder{id: 445, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b444 = charBuilder{}
	b445.items = []builder{&b444}
	var b446 = sequenceBuilder{id: 446, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b446.items = []builder{&b817, &b445}
	b447.items = []builder{&b817, &b445, &b446}
	b448.items = []builder{&b443, &b447, &b817, &b785}
	var b431 = sequenceBuilder{id: 431, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b426 = sequenceBuilder{id: 426, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b421 = sequenceBuilder{id: 421, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b414 = charBuilder{}
	var b415 = charBuilder{}
	var b416 = charBuilder{}
	var b417 = charBuilder{}
	var b418 = charBuilder{}
	var b419 = charBuilder{}
	var b420 = charBuilder{}
	b421.items = []builder{&b414, &b415, &b416, &b417, &b418, &b419, &b420}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b424.items = []builder{&b817, &b14}
	b425.items = []builder{&b817, &b14, &b424}
	var b423 = sequenceBuilder{id: 423, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b422 = charBuilder{}
	b423.items = []builder{&b422}
	b426.items = []builder{&b421, &b425, &b817, &b423}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b428 = sequenceBuilder{id: 428, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b427 = charBuilder{}
	b428.items = []builder{&b427}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b429.items = []builder{&b817, &b428}
	b430.items = []builder{&b817, &b428, &b429}
	b431.items = []builder{&b426, &b430, &b817, &b785}
	b458.options = []builder{&b448, &b431}
	var b462 = sequenceBuilder{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b459 = choiceBuilder{id: 459, commit: 2}
	b459.options = []builder{&b448, &b431, &b785}
	b460.items = []builder{&b799, &b817, &b459}
	var b461 = sequenceBuilder{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b461.items = []builder{&b817, &b460}
	b462.items = []builder{&b817, &b460, &b461}
	b463.items = []builder{&b458, &b462}
	var b465 = sequenceBuilder{id: 465, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b464 = charBuilder{}
	b465.items = []builder{&b464}
	b470.items = []builder{&b455, &b467, &b817, &b376, &b469, &b817, &b457, &b817, &b799, &b817, &b463, &b817, &b799, &b817, &b465}
	var b532 = sequenceBuilder{id: 532, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b519 = sequenceBuilder{id: 519, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b513 = charBuilder{}
	var b514 = charBuilder{}
	var b515 = charBuilder{}
	var b516 = charBuilder{}
	var b517 = charBuilder{}
	var b518 = charBuilder{}
	b519.items = []builder{&b513, &b514, &b515, &b516, &b517, &b518}
	var b531 = sequenceBuilder{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b530.items = []builder{&b817, &b14}
	b531.items = []builder{&b817, &b14, &b530}
	var b521 = sequenceBuilder{id: 521, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b520 = charBuilder{}
	b521.items = []builder{&b520}
	var b527 = sequenceBuilder{id: 527, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b522 = choiceBuilder{id: 522, commit: 2}
	var b512 = sequenceBuilder{id: 512, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b507 = sequenceBuilder{id: 507, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b500 = sequenceBuilder{id: 500, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b496 = charBuilder{}
	var b497 = charBuilder{}
	var b498 = charBuilder{}
	var b499 = charBuilder{}
	b500.items = []builder{&b496, &b497, &b498, &b499}
	var b504 = sequenceBuilder{id: 504, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b503 = sequenceBuilder{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b503.items = []builder{&b817, &b14}
	b504.items = []builder{&b817, &b14, &b503}
	var b495 = choiceBuilder{id: 495, commit: 66}
	var b494 = sequenceBuilder{id: 494, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b493 = sequenceBuilder{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b492.items = []builder{&b817, &b14}
	b493.items = []builder{&b817, &b14, &b492}
	b494.items = []builder{&b105, &b493, &b817, &b481}
	b495.options = []builder{&b481, &b494, &b491}
	var b506 = sequenceBuilder{id: 506, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b505 = sequenceBuilder{id: 505, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b505.items = []builder{&b817, &b14}
	b506.items = []builder{&b817, &b14, &b505}
	var b502 = sequenceBuilder{id: 502, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b501 = charBuilder{}
	b502.items = []builder{&b501}
	b507.items = []builder{&b500, &b504, &b817, &b495, &b506, &b817, &b502}
	var b511 = sequenceBuilder{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b509 = sequenceBuilder{id: 509, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b508 = charBuilder{}
	b509.items = []builder{&b508}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b510.items = []builder{&b817, &b509}
	b511.items = []builder{&b817, &b509, &b510}
	b512.items = []builder{&b507, &b511, &b817, &b785}
	b522.options = []builder{&b512, &b431}
	var b526 = sequenceBuilder{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b523 = choiceBuilder{id: 523, commit: 2}
	b523.options = []builder{&b512, &b431, &b785}
	b524.items = []builder{&b799, &b817, &b523}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b525.items = []builder{&b817, &b524}
	b526.items = []builder{&b817, &b524, &b525}
	b527.items = []builder{&b522, &b526}
	var b529 = sequenceBuilder{id: 529, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b528 = charBuilder{}
	b529.items = []builder{&b528}
	b532.items = []builder{&b519, &b531, &b817, &b521, &b817, &b799, &b817, &b527, &b817, &b799, &b817, &b529}
	var b573 = sequenceBuilder{id: 573, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b562 = sequenceBuilder{id: 562, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b559 = charBuilder{}
	var b560 = charBuilder{}
	var b561 = charBuilder{}
	b562.items = []builder{&b559, &b560, &b561}
	var b572 = choiceBuilder{id: 572, commit: 2}
	var b568 = sequenceBuilder{id: 568, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b563 = sequenceBuilder{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b563.items = []builder{&b817, &b14}
	b564.items = []builder{&b14, &b563}
	var b558 = choiceBuilder{id: 558, commit: 66}
	var b557 = choiceBuilder{id: 557, commit: 64, name: "range-over-expression"}
	var b556 = sequenceBuilder{id: 556, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b553 = sequenceBuilder{id: 553, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b552 = sequenceBuilder{id: 552, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b552.items = []builder{&b817, &b14}
	b553.items = []builder{&b817, &b14, &b552}
	var b550 = sequenceBuilder{id: 550, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b548 = charBuilder{}
	var b549 = charBuilder{}
	b550.items = []builder{&b548, &b549}
	var b555 = sequenceBuilder{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b554 = sequenceBuilder{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b554.items = []builder{&b817, &b14}
	b555.items = []builder{&b817, &b14, &b554}
	var b551 = choiceBuilder{id: 551, commit: 2}
	b551.options = []builder{&b376, &b227}
	b556.items = []builder{&b105, &b553, &b817, &b550, &b555, &b817, &b551}
	b557.options = []builder{&b556, &b227}
	b558.options = []builder{&b376, &b557}
	b565.items = []builder{&b564, &b817, &b558}
	var b567 = sequenceBuilder{id: 567, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b566 = sequenceBuilder{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b566.items = []builder{&b817, &b14}
	b567.items = []builder{&b817, &b14, &b566}
	b568.items = []builder{&b565, &b567, &b817, &b192}
	var b571 = sequenceBuilder{id: 571, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b570 = sequenceBuilder{id: 570, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b569 = sequenceBuilder{id: 569, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b569.items = []builder{&b817, &b14}
	b570.items = []builder{&b14, &b569}
	b571.items = []builder{&b570, &b817, &b192}
	b572.options = []builder{&b568, &b571}
	b573.items = []builder{&b562, &b817, &b572}
	var b721 = choiceBuilder{id: 721, commit: 66}
	var b634 = sequenceBuilder{id: 634, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b630 = sequenceBuilder{id: 630, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b627 = charBuilder{}
	var b628 = charBuilder{}
	var b629 = charBuilder{}
	b630.items = []builder{&b627, &b628, &b629}
	var b633 = sequenceBuilder{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b632 = sequenceBuilder{id: 632, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b632.items = []builder{&b817, &b14}
	b633.items = []builder{&b817, &b14, &b632}
	var b631 = choiceBuilder{id: 631, commit: 2}
	var b621 = sequenceBuilder{id: 621, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b620 = sequenceBuilder{id: 620, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b615.items = []builder{&b817, &b14}
	b616.items = []builder{&b14, &b615}
	var b614 = sequenceBuilder{id: 614, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b613 = charBuilder{}
	b614.items = []builder{&b613}
	b617.items = []builder{&b616, &b817, &b614}
	var b619 = sequenceBuilder{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b618 = sequenceBuilder{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b618.items = []builder{&b817, &b14}
	b619.items = []builder{&b817, &b14, &b618}
	b620.items = []builder{&b105, &b817, &b617, &b619, &b817, &b376}
	b621.items = []builder{&b620}
	var b626 = sequenceBuilder{id: 626, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b623 = sequenceBuilder{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b622 = charBuilder{}
	b623.items = []builder{&b622}
	var b625 = sequenceBuilder{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b624.items = []builder{&b817, &b14}
	b625.items = []builder{&b817, &b14, &b624}
	b626.items = []builder{&b623, &b625, &b817, &b620}
	b631.options = []builder{&b621, &b626}
	b634.items = []builder{&b630, &b633, &b817, &b631}
	var b655 = sequenceBuilder{id: 655, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b648 = sequenceBuilder{id: 648, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b645 = charBuilder{}
	var b646 = charBuilder{}
	var b647 = charBuilder{}
	b648.items = []builder{&b645, &b646, &b647}
	var b654 = sequenceBuilder{id: 654, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b653 = sequenceBuilder{id: 653, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b653.items = []builder{&b817, &b14}
	b654.items = []builder{&b817, &b14, &b653}
	var b650 = sequenceBuilder{id: 650, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b649 = charBuilder{}
	b650.items = []builder{&b649}
	var b640 = sequenceBuilder{id: 640, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b635 = choiceBuilder{id: 635, commit: 2}
	b635.options = []builder{&b621, &b626}
	var b639 = sequenceBuilder{id: 639, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b637 = sequenceBuilder{id: 637, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b636 = choiceBuilder{id: 636, commit: 2}
	b636.options = []builder{&b621, &b626}
	b637.items = []builder{&b115, &b817, &b636}
	var b638 = sequenceBuilder{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b638.items = []builder{&b817, &b637}
	b639.items = []builder{&b817, &b637, &b638}
	b640.items = []builder{&b635, &b639}
	var b652 = sequenceBuilder{id: 652, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b651 = charBuilder{}
	b652.items = []builder{&b651}
	b655.items = []builder{&b648, &b654, &b817, &b650, &b817, &b115, &b817, &b640, &b817, &b115, &b817, &b652}
	var b670 = sequenceBuilder{id: 670, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b659 = sequenceBuilder{id: 659, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b656 = charBuilder{}
	var b657 = charBuilder{}
	var b658 = charBuilder{}
	b659.items = []builder{&b656, &b657, &b658}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b666.items = []builder{&b817, &b14}
	b667.items = []builder{&b817, &b14, &b666}
	var b661 = sequenceBuilder{id: 661, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b660 = charBuilder{}
	b661.items = []builder{&b660}
	var b669 = sequenceBuilder{id: 669, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b668 = sequenceBuilder{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b668.items = []builder{&b817, &b14}
	b669.items = []builder{&b817, &b14, &b668}
	var b663 = sequenceBuilder{id: 663, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b662 = charBuilder{}
	b663.items = []builder{&b662}
	var b644 = sequenceBuilder{id: 644, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b643 = sequenceBuilder{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b641 = sequenceBuilder{id: 641, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b641.items = []builder{&b115, &b817, &b621}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b642.items = []builder{&b817, &b641}
	b643.items = []builder{&b817, &b641, &b642}
	b644.items = []builder{&b621, &b643}
	var b665 = sequenceBuilder{id: 665, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b664 = charBuilder{}
	b665.items = []builder{&b664}
	b670.items = []builder{&b659, &b667, &b817, &b661, &b669, &b817, &b663, &b817, &b115, &b817, &b644, &b817, &b115, &b817, &b665}
	var b686 = sequenceBuilder{id: 686, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b682 = sequenceBuilder{id: 682, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	var b681 = charBuilder{}
	b682.items = []builder{&b680, &b681}
	var b685 = sequenceBuilder{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b684 = sequenceBuilder{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b684.items = []builder{&b817, &b14}
	b685.items = []builder{&b817, &b14, &b684}
	var b683 = choiceBuilder{id: 683, commit: 2}
	var b674 = sequenceBuilder{id: 674, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b673 = sequenceBuilder{id: 673, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b672 = sequenceBuilder{id: 672, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b671 = sequenceBuilder{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b671.items = []builder{&b817, &b14}
	b672.items = []builder{&b817, &b14, &b671}
	b673.items = []builder{&b105, &b672, &b817, &b202}
	b674.items = []builder{&b673}
	var b679 = sequenceBuilder{id: 679, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b676 = sequenceBuilder{id: 676, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b675 = charBuilder{}
	b676.items = []builder{&b675}
	var b678 = sequenceBuilder{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b677 = sequenceBuilder{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b677.items = []builder{&b817, &b14}
	b678.items = []builder{&b817, &b14, &b677}
	b679.items = []builder{&b676, &b678, &b817, &b673}
	b683.options = []builder{&b674, &b679}
	b686.items = []builder{&b682, &b685, &b817, &b683}
	var b706 = sequenceBuilder{id: 706, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b699 = sequenceBuilder{id: 699, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b697 = charBuilder{}
	var b698 = charBuilder{}
	b699.items = []builder{&b697, &b698}
	var b705 = sequenceBuilder{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b704.items = []builder{&b817, &b14}
	b705.items = []builder{&b817, &b14, &b704}
	var b701 = sequenceBuilder{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b700 = charBuilder{}
	b701.items = []builder{&b700}
	var b696 = sequenceBuilder{id: 696, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b691 = choiceBuilder{id: 691, commit: 2}
	b691.options = []builder{&b674, &b679}
	var b695 = sequenceBuilder{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b693 = sequenceBuilder{id: 693, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b692 = choiceBuilder{id: 692, commit: 2}
	b692.options = []builder{&b674, &b679}
	b693.items = []builder{&b115, &b817, &b692}
	var b694 = sequenceBuilder{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b694.items = []builder{&b817, &b693}
	b695.items = []builder{&b817, &b693, &b694}
	b696.items = []builder{&b691, &b695}
	var b703 = sequenceBuilder{id: 703, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b702 = charBuilder{}
	b703.items = []builder{&b702}
	b706.items = []builder{&b699, &b705, &b817, &b701, &b817, &b115, &b817, &b696, &b817, &b115, &b817, &b703}
	var b720 = sequenceBuilder{id: 720, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b709 = sequenceBuilder{id: 709, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b707 = charBuilder{}
	var b708 = charBuilder{}
	b709.items = []builder{&b707, &b708}
	var b717 = sequenceBuilder{id: 717, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b716 = sequenceBuilder{id: 716, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b716.items = []builder{&b817, &b14}
	b717.items = []builder{&b817, &b14, &b716}
	var b711 = sequenceBuilder{id: 711, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b710 = charBuilder{}
	b711.items = []builder{&b710}
	var b719 = sequenceBuilder{id: 719, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b718 = sequenceBuilder{id: 718, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b718.items = []builder{&b817, &b14}
	b719.items = []builder{&b817, &b14, &b718}
	var b713 = sequenceBuilder{id: 713, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b712 = charBuilder{}
	b713.items = []builder{&b712}
	var b690 = sequenceBuilder{id: 690, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b689 = sequenceBuilder{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b687.items = []builder{&b115, &b817, &b674}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b688.items = []builder{&b817, &b687}
	b689.items = []builder{&b817, &b687, &b688}
	b690.items = []builder{&b674, &b689}
	var b715 = sequenceBuilder{id: 715, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b714 = charBuilder{}
	b715.items = []builder{&b714}
	b720.items = []builder{&b709, &b717, &b817, &b711, &b719, &b817, &b713, &b817, &b115, &b817, &b690, &b817, &b115, &b817, &b715}
	b721.options = []builder{&b634, &b655, &b670, &b686, &b706, &b720}
	var b764 = choiceBuilder{id: 764, commit: 64, name: "require"}
	var b748 = sequenceBuilder{id: 748, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b745 = sequenceBuilder{id: 745, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b738 = charBuilder{}
	var b739 = charBuilder{}
	var b740 = charBuilder{}
	var b741 = charBuilder{}
	var b742 = charBuilder{}
	var b743 = charBuilder{}
	var b744 = charBuilder{}
	b745.items = []builder{&b738, &b739, &b740, &b741, &b742, &b743, &b744}
	var b747 = sequenceBuilder{id: 747, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b746 = sequenceBuilder{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b746.items = []builder{&b817, &b14}
	b747.items = []builder{&b817, &b14, &b746}
	var b733 = choiceBuilder{id: 733, commit: 64, name: "require-fact"}
	var b732 = sequenceBuilder{id: 732, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b724 = choiceBuilder{id: 724, commit: 2}
	var b723 = sequenceBuilder{id: 723, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b722 = charBuilder{}
	b723.items = []builder{&b722}
	b724.options = []builder{&b105, &b723}
	var b729 = sequenceBuilder{id: 729, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b728 = sequenceBuilder{id: 728, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b727 = sequenceBuilder{id: 727, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b727.items = []builder{&b817, &b14}
	b728.items = []builder{&b14, &b727}
	var b726 = sequenceBuilder{id: 726, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b725 = charBuilder{}
	b726.items = []builder{&b725}
	b729.items = []builder{&b728, &b817, &b726}
	var b731 = sequenceBuilder{id: 731, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b730 = sequenceBuilder{id: 730, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b730.items = []builder{&b817, &b14}
	b731.items = []builder{&b817, &b14, &b730}
	b732.items = []builder{&b724, &b817, &b729, &b731, &b817, &b88}
	b733.options = []builder{&b88, &b732}
	b748.items = []builder{&b745, &b747, &b817, &b733}
	var b763 = sequenceBuilder{id: 763, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b756 = sequenceBuilder{id: 756, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b749 = charBuilder{}
	var b750 = charBuilder{}
	var b751 = charBuilder{}
	var b752 = charBuilder{}
	var b753 = charBuilder{}
	var b754 = charBuilder{}
	var b755 = charBuilder{}
	b756.items = []builder{&b749, &b750, &b751, &b752, &b753, &b754, &b755}
	var b762 = sequenceBuilder{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b761 = sequenceBuilder{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b761.items = []builder{&b817, &b14}
	b762.items = []builder{&b817, &b14, &b761}
	var b758 = sequenceBuilder{id: 758, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b757 = charBuilder{}
	b758.items = []builder{&b757}
	var b737 = sequenceBuilder{id: 737, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b734 = sequenceBuilder{id: 734, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b734.items = []builder{&b115, &b817, &b733}
	var b735 = sequenceBuilder{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b735.items = []builder{&b817, &b734}
	b736.items = []builder{&b817, &b734, &b735}
	b737.items = []builder{&b733, &b736}
	var b760 = sequenceBuilder{id: 760, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b759 = charBuilder{}
	b760.items = []builder{&b759}
	b763.items = []builder{&b756, &b762, &b817, &b758, &b817, &b115, &b817, &b737, &b817, &b115, &b817, &b760}
	b764.options = []builder{&b748, &b763}
	var b774 = sequenceBuilder{id: 774, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b771 = sequenceBuilder{id: 771, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b765 = charBuilder{}
	var b766 = charBuilder{}
	var b767 = charBuilder{}
	var b768 = charBuilder{}
	var b769 = charBuilder{}
	var b770 = charBuilder{}
	b771.items = []builder{&b765, &b766, &b767, &b768, &b769, &b770}
	var b773 = sequenceBuilder{id: 773, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b772 = sequenceBuilder{id: 772, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b772.items = []builder{&b817, &b14}
	b773.items = []builder{&b817, &b14, &b772}
	b774.items = []builder{&b771, &b773, &b817, &b721}
	var b794 = sequenceBuilder{id: 794, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b787 = sequenceBuilder{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b786 = charBuilder{}
	b787.items = []builder{&b786}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b790.items = []builder{&b817, &b14}
	b791.items = []builder{&b817, &b14, &b790}
	var b793 = sequenceBuilder{id: 793, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b792.items = []builder{&b817, &b14}
	b793.items = []builder{&b817, &b14, &b792}
	var b789 = sequenceBuilder{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b788 = charBuilder{}
	b789.items = []builder{&b788}
	b794.items = []builder{&b787, &b791, &b817, &b785, &b793, &b817, &b789}
	b785.options = []builder{&b187, &b413, &b470, &b532, &b573, &b721, &b764, &b774, &b794, &b775}
	var b802 = sequenceBuilder{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b800.items = []builder{&b799, &b817, &b785}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b801.items = []builder{&b817, &b800}
	b802.items = []builder{&b817, &b800, &b801}
	b803.items = []builder{&b785, &b802}
	b818.items = []builder{&b814, &b817, &b799, &b817, &b803, &b817, &b799}
	b819.items = []builder{&b817, &b818, &b817}

	return parseInput(r, &p819, &b819)
}
