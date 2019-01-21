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
	var p827 = sequenceParser{id: 827, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p825 = choiceParser{id: 825, commit: 2}
	var p823 = choiceParser{id: 823, commit: 70, name: "ws", generalizations: []int{825, 15}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{823, 825, 15}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{823, 825, 15}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{823, 825, 15}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{823, 825, 15}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{823, 825, 15}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{823, 825, 15}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p823.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p824 = sequenceParser{id: 824, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825}}
	var p43 = sequenceParser{id: 43, commit: 66, name: "comment", ranges: [][]int{{1, 1}, {0, 1}}}
	var p39 = choiceParser{id: 39, commit: 66, name: "comment-part"}
	var p22 = sequenceParser{id: 22, commit: 74, name: "line-comment", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{39}}
	var p21 = sequenceParser{id: 21, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p19 = charParser{id: 19, chars: []rune{47}}
	var p20 = charParser{id: 20, chars: []rune{47}}
	p21.items = []parser{&p19, &p20}
	var p18 = sequenceParser{id: 18, commit: 72, name: "line-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var p17 = sequenceParser{id: 17, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p16 = charParser{id: 16, not: true, chars: []rune{10}}
	p17.items = []parser{&p16}
	p18.items = []parser{&p17}
	p22.items = []parser{&p21, &p18}
	var p38 = sequenceParser{id: 38, commit: 74, name: "block-comment", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{39}}
	var p34 = sequenceParser{id: 34, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p32 = charParser{id: 32, chars: []rune{47}}
	var p33 = charParser{id: 33, chars: []rune{42}}
	p34.items = []parser{&p32, &p33}
	var p31 = sequenceParser{id: 31, commit: 72, name: "block-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var p30 = choiceParser{id: 30, commit: 10}
	var p24 = sequenceParser{id: 24, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{30}}
	var p23 = charParser{id: 23, not: true, chars: []rune{42}}
	p24.items = []parser{&p23}
	var p29 = sequenceParser{id: 29, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{30}}
	var p26 = sequenceParser{id: 26, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p25 = charParser{id: 25, chars: []rune{42}}
	p26.items = []parser{&p25}
	var p28 = sequenceParser{id: 28, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p27 = charParser{id: 27, not: true, chars: []rune{47}}
	p28.items = []parser{&p27}
	p29.items = []parser{&p26, &p28}
	p30.options = []parser{&p24, &p29}
	p31.items = []parser{&p30}
	var p37 = sequenceParser{id: 37, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p35 = charParser{id: 35, chars: []rune{42}}
	var p36 = charParser{id: 36, chars: []rune{47}}
	p37.items = []parser{&p35, &p36}
	p38.items = []parser{&p34, &p31, &p37}
	p39.options = []parser{&p22, &p38}
	var p42 = sequenceParser{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p40 = sequenceParser{id: 40, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{805, 15, 112}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p40.items = []parser{&p14, &p825, &p39}
	var p41 = sequenceParser{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p41.items = []parser{&p825, &p40}
	p42.items = []parser{&p825, &p40, &p41}
	p43.items = []parser{&p39, &p42}
	p824.items = []parser{&p43}
	p825.options = []parser{&p823, &p824}
	var p826 = sequenceParser{id: 826, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p822 = sequenceParser{id: 822, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p819 = sequenceParser{id: 819, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p817 = charParser{id: 817, chars: []rune{35}}
	var p818 = charParser{id: 818, chars: []rune{33}}
	p819.items = []parser{&p817, &p818}
	var p816 = sequenceParser{id: 816, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p815 = sequenceParser{id: 815, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p813 = sequenceParser{id: 813, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p812 = charParser{id: 812, not: true, chars: []rune{10}}
	p813.items = []parser{&p812}
	var p814 = sequenceParser{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p814.items = []parser{&p825, &p813}
	p815.items = []parser{&p813, &p814}
	p816.items = []parser{&p815}
	var p821 = sequenceParser{id: 821, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p820 = charParser{id: 820, chars: []rune{10}}
	p821.items = []parser{&p820}
	p822.items = []parser{&p819, &p825, &p816, &p825, &p821}
	var p807 = sequenceParser{id: 807, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p805 = choiceParser{id: 805, commit: 2}
	var p804 = sequenceParser{id: 804, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{805}}
	var p803 = charParser{id: 803, chars: []rune{59}}
	p804.items = []parser{&p803}
	p805.options = []parser{&p804, &p14}
	var p806 = sequenceParser{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p806.items = []parser{&p825, &p805}
	p807.items = []parser{&p805, &p806}
	var p811 = sequenceParser{id: 811, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p793 = choiceParser{id: 793, commit: 66, name: "statement", generalizations: []int{480, 541}}
	var p186 = sequenceParser{id: 186, commit: 64, name: "return-value", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793, 480, 541}}
	var p183 = sequenceParser{id: 183, commit: 74, name: "return-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p182 = sequenceParser{id: 182, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p176 = charParser{id: 176, chars: []rune{114}}
	var p177 = charParser{id: 177, chars: []rune{101}}
	var p178 = charParser{id: 178, chars: []rune{116}}
	var p179 = charParser{id: 179, chars: []rune{117}}
	var p180 = charParser{id: 180, chars: []rune{114}}
	var p181 = charParser{id: 181, chars: []rune{110}}
	p182.items = []parser{&p176, &p177, &p178, &p179, &p180, &p181}
	var p15 = choiceParser{id: 15, commit: 66, name: "wsep"}
	p15.options = []parser{&p823, &p14}
	p183.items = []parser{&p182, &p15}
	var p185 = sequenceParser{id: 185, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p184 = sequenceParser{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p184.items = []parser{&p825, &p14}
	p185.items = []parser{&p825, &p14, &p184}
	var p401 = choiceParser{id: 401, commit: 66, name: "expression", generalizations: []int{115, 783, 198, 590, 583, 793}}
	var p272 = choiceParser{id: 272, commit: 66, name: "primary-expression", generalizations: []int{115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p61 = choiceParser{id: 61, commit: 64, name: "int", generalizations: []int{272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p52 = sequenceParser{id: 52, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{61, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p51 = sequenceParser{id: 51, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p50 = charParser{id: 50, ranges: [][]rune{{49, 57}}}
	p51.items = []parser{&p50}
	var p45 = sequenceParser{id: 45, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p44 = charParser{id: 44, ranges: [][]rune{{48, 57}}}
	p45.items = []parser{&p44}
	p52.items = []parser{&p51, &p45}
	var p55 = sequenceParser{id: 55, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{61, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p54 = sequenceParser{id: 54, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p53 = charParser{id: 53, chars: []rune{48}}
	p54.items = []parser{&p53}
	var p47 = sequenceParser{id: 47, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p46 = charParser{id: 46, ranges: [][]rune{{48, 55}}}
	p47.items = []parser{&p46}
	p55.items = []parser{&p54, &p47}
	var p60 = sequenceParser{id: 60, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{61, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p57 = sequenceParser{id: 57, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p56 = charParser{id: 56, chars: []rune{48}}
	p57.items = []parser{&p56}
	var p59 = sequenceParser{id: 59, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p58 = charParser{id: 58, chars: []rune{120, 88}}
	p59.items = []parser{&p58}
	var p49 = sequenceParser{id: 49, commit: 66, name: "hexa-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p48 = charParser{id: 48, ranges: [][]rune{{48, 57}, {97, 102}, {65, 70}}}
	p49.items = []parser{&p48}
	p60.items = []parser{&p57, &p59, &p49}
	p61.options = []parser{&p52, &p55, &p60}
	var p74 = choiceParser{id: 74, commit: 72, name: "float", generalizations: []int{272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p69 = sequenceParser{id: 69, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{74, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p68 = sequenceParser{id: 68, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p67 = charParser{id: 67, chars: []rune{46}}
	p68.items = []parser{&p67}
	var p66 = sequenceParser{id: 66, commit: 74, name: "exponent", ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var p63 = sequenceParser{id: 63, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p62 = charParser{id: 62, chars: []rune{101, 69}}
	p63.items = []parser{&p62}
	var p65 = sequenceParser{id: 65, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p64 = charParser{id: 64, chars: []rune{43, 45}}
	p65.items = []parser{&p64}
	p66.items = []parser{&p63, &p65, &p45}
	p69.items = []parser{&p45, &p68, &p45, &p66}
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{74, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p71 = sequenceParser{id: 71, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p70 = charParser{id: 70, chars: []rune{46}}
	p71.items = []parser{&p70}
	p72.items = []parser{&p71, &p45, &p66}
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{74, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	p73.items = []parser{&p45, &p66}
	p74.options = []parser{&p69, &p72, &p73}
	var p87 = sequenceParser{id: 87, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 115, 140, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 756, 793}}
	var p76 = sequenceParser{id: 76, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p75 = charParser{id: 75, chars: []rune{34}}
	p76.items = []parser{&p75}
	var p84 = choiceParser{id: 84, commit: 10}
	var p78 = sequenceParser{id: 78, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{84}}
	var p77 = charParser{id: 77, not: true, chars: []rune{92, 34}}
	p78.items = []parser{&p77}
	var p83 = sequenceParser{id: 83, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{84}}
	var p80 = sequenceParser{id: 80, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p79 = charParser{id: 79, chars: []rune{92}}
	p80.items = []parser{&p79}
	var p82 = sequenceParser{id: 82, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p81 = charParser{id: 81, not: true}
	p82.items = []parser{&p81}
	p83.items = []parser{&p80, &p82}
	p84.options = []parser{&p78, &p83}
	var p86 = sequenceParser{id: 86, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p85 = charParser{id: 85, chars: []rune{34}}
	p86.items = []parser{&p85}
	p87.items = []parser{&p76, &p84, &p86}
	var p99 = choiceParser{id: 99, commit: 66, name: "bool", generalizations: []int{272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p92 = sequenceParser{id: 92, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{99, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p88 = charParser{id: 88, chars: []rune{116}}
	var p89 = charParser{id: 89, chars: []rune{114}}
	var p90 = charParser{id: 90, chars: []rune{117}}
	var p91 = charParser{id: 91, chars: []rune{101}}
	p92.items = []parser{&p88, &p89, &p90, &p91}
	var p98 = sequenceParser{id: 98, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{99, 272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p93 = charParser{id: 93, chars: []rune{102}}
	var p94 = charParser{id: 94, chars: []rune{97}}
	var p95 = charParser{id: 95, chars: []rune{108}}
	var p96 = charParser{id: 96, chars: []rune{115}}
	var p97 = charParser{id: 97, chars: []rune{101}}
	p98.items = []parser{&p93, &p94, &p95, &p96, &p97}
	p99.options = []parser{&p92, &p98}
	var p514 = sequenceParser{id: 514, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 115, 783, 198, 401, 338, 339, 340, 341, 342, 393, 518, 590, 583, 793}}
	var p506 = sequenceParser{id: 506, commit: 74, name: "receive-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p505 = sequenceParser{id: 505, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p498 = charParser{id: 498, chars: []rune{114}}
	var p499 = charParser{id: 499, chars: []rune{101}}
	var p500 = charParser{id: 500, chars: []rune{99}}
	var p501 = charParser{id: 501, chars: []rune{101}}
	var p502 = charParser{id: 502, chars: []rune{105}}
	var p503 = charParser{id: 503, chars: []rune{118}}
	var p504 = charParser{id: 504, chars: []rune{101}}
	p505.items = []parser{&p498, &p499, &p500, &p501, &p502, &p503, &p504}
	p506.items = []parser{&p505, &p15}
	var p513 = sequenceParser{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p512 = sequenceParser{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p512.items = []parser{&p825, &p14}
	p513.items = []parser{&p825, &p14, &p512}
	p514.items = []parser{&p506, &p513, &p825, &p272}
	var p104 = sequenceParser{id: 104, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{272, 115, 140, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 747, 793}}
	var p101 = sequenceParser{id: 101, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p100 = charParser{id: 100, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p101.items = []parser{&p100}
	var p103 = sequenceParser{id: 103, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p102 = charParser{id: 102, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p103.items = []parser{&p102}
	p104.items = []parser{&p101, &p103}
	var p125 = sequenceParser{id: 125, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{115, 272, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p124 = sequenceParser{id: 124, commit: 66, name: "list-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p121 = sequenceParser{id: 121, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p120 = charParser{id: 120, chars: []rune{91}}
	p121.items = []parser{&p120}
	var p114 = sequenceParser{id: 114, commit: 66, name: "list-sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p112 = choiceParser{id: 112, commit: 2}
	var p111 = sequenceParser{id: 111, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{112}}
	var p110 = charParser{id: 110, chars: []rune{44}}
	p111.items = []parser{&p110}
	p112.options = []parser{&p14, &p111}
	var p113 = sequenceParser{id: 113, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p113.items = []parser{&p825, &p112}
	p114.items = []parser{&p112, &p113}
	var p119 = sequenceParser{id: 119, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p115 = choiceParser{id: 115, commit: 66, name: "list-item"}
	var p109 = sequenceParser{id: 109, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{115, 148, 149}}
	var p108 = sequenceParser{id: 108, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	p108.items = []parser{&p105, &p106, &p107}
	p109.items = []parser{&p272, &p825, &p108}
	p115.options = []parser{&p401, &p109}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p116.items = []parser{&p114, &p825, &p115}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p117.items = []parser{&p825, &p116}
	p118.items = []parser{&p825, &p116, &p117}
	p119.items = []parser{&p115, &p118}
	var p123 = sequenceParser{id: 123, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p122 = charParser{id: 122, chars: []rune{93}}
	p123.items = []parser{&p122}
	p124.items = []parser{&p121, &p825, &p114, &p825, &p119, &p825, &p114, &p825, &p123}
	p125.items = []parser{&p124}
	var p130 = sequenceParser{id: 130, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p127 = sequenceParser{id: 127, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p126 = charParser{id: 126, chars: []rune{126}}
	p127.items = []parser{&p126}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p128.items = []parser{&p825, &p14}
	p129.items = []parser{&p825, &p14, &p128}
	p130.items = []parser{&p127, &p129, &p825, &p124}
	var p159 = sequenceParser{id: 159, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{272, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p158 = sequenceParser{id: 158, commit: 66, name: "struct-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var p155 = sequenceParser{id: 155, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p154 = charParser{id: 154, chars: []rune{123}}
	p155.items = []parser{&p154}
	var p153 = sequenceParser{id: 153, commit: 66, name: "entry-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p148 = choiceParser{id: 148, commit: 2}
	var p147 = sequenceParser{id: 147, commit: 64, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{148, 149}}
	var p140 = choiceParser{id: 140, commit: 2}
	var p139 = sequenceParser{id: 139, commit: 64, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{140}}
	var p132 = sequenceParser{id: 132, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p131 = charParser{id: 131, chars: []rune{91}}
	p132.items = []parser{&p131}
	var p136 = sequenceParser{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p135 = sequenceParser{id: 135, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p135.items = []parser{&p825, &p14}
	p136.items = []parser{&p825, &p14, &p135}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p137.items = []parser{&p825, &p14}
	p138.items = []parser{&p825, &p14, &p137}
	var p134 = sequenceParser{id: 134, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p133 = charParser{id: 133, chars: []rune{93}}
	p134.items = []parser{&p133}
	p139.items = []parser{&p132, &p136, &p825, &p401, &p138, &p825, &p134}
	p140.options = []parser{&p104, &p87, &p139}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p143.items = []parser{&p825, &p14}
	p144.items = []parser{&p825, &p14, &p143}
	var p142 = sequenceParser{id: 142, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p141 = charParser{id: 141, chars: []rune{58}}
	p142.items = []parser{&p141}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p145.items = []parser{&p825, &p14}
	p146.items = []parser{&p825, &p14, &p145}
	p147.items = []parser{&p140, &p144, &p825, &p142, &p146, &p825, &p401}
	p148.options = []parser{&p147, &p109}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p149 = choiceParser{id: 149, commit: 2}
	p149.options = []parser{&p147, &p109}
	p150.items = []parser{&p114, &p825, &p149}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p151.items = []parser{&p825, &p150}
	p152.items = []parser{&p825, &p150, &p151}
	p153.items = []parser{&p148, &p152}
	var p157 = sequenceParser{id: 157, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p156 = charParser{id: 156, chars: []rune{125}}
	p157.items = []parser{&p156}
	p158.items = []parser{&p155, &p825, &p114, &p825, &p153, &p825, &p114, &p825, &p157}
	p159.items = []parser{&p158}
	var p164 = sequenceParser{id: 164, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 783, 198, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p161 = sequenceParser{id: 161, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p160 = charParser{id: 160, chars: []rune{126}}
	p161.items = []parser{&p160}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p162.items = []parser{&p825, &p14}
	p163.items = []parser{&p825, &p14, &p162}
	p164.items = []parser{&p161, &p163, &p825, &p158}
	var p207 = sequenceParser{id: 207, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 198, 272, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p204 = sequenceParser{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p202 = charParser{id: 202, chars: []rune{102}}
	var p203 = charParser{id: 203, chars: []rune{110}}
	p204.items = []parser{&p202, &p203}
	var p206 = sequenceParser{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p205.items = []parser{&p825, &p14}
	p206.items = []parser{&p825, &p14, &p205}
	var p201 = sequenceParser{id: 201, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p193 = sequenceParser{id: 193, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p192 = charParser{id: 192, chars: []rune{40}}
	p193.items = []parser{&p192}
	var p195 = choiceParser{id: 195, commit: 2}
	var p168 = sequenceParser{id: 168, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{195}}
	var p167 = sequenceParser{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p165.items = []parser{&p114, &p825, &p104}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p166.items = []parser{&p825, &p165}
	p167.items = []parser{&p825, &p165, &p166}
	p168.items = []parser{&p104, &p167}
	var p194 = sequenceParser{id: 194, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{195}}
	var p175 = sequenceParser{id: 175, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{195}}
	var p172 = sequenceParser{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p169 = charParser{id: 169, chars: []rune{46}}
	var p170 = charParser{id: 170, chars: []rune{46}}
	var p171 = charParser{id: 171, chars: []rune{46}}
	p172.items = []parser{&p169, &p170, &p171}
	var p174 = sequenceParser{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p173 = sequenceParser{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p173.items = []parser{&p825, &p14}
	p174.items = []parser{&p825, &p14, &p173}
	p175.items = []parser{&p172, &p174, &p825, &p104}
	p194.items = []parser{&p168, &p825, &p114, &p825, &p175}
	p195.options = []parser{&p168, &p194, &p175}
	var p197 = sequenceParser{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p196 = charParser{id: 196, chars: []rune{41}}
	p197.items = []parser{&p196}
	var p200 = sequenceParser{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p199.items = []parser{&p825, &p14}
	p200.items = []parser{&p825, &p14, &p199}
	var p198 = choiceParser{id: 198, commit: 2}
	var p783 = choiceParser{id: 783, commit: 66, name: "simple-statement", generalizations: []int{198, 793}}
	var p511 = sequenceParser{id: 511, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 198, 518, 793}}
	var p497 = sequenceParser{id: 497, commit: 74, name: "send-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p496 = sequenceParser{id: 496, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p492 = charParser{id: 492, chars: []rune{115}}
	var p493 = charParser{id: 493, chars: []rune{101}}
	var p494 = charParser{id: 494, chars: []rune{110}}
	var p495 = charParser{id: 495, chars: []rune{100}}
	p496.items = []parser{&p492, &p493, &p494, &p495}
	p497.items = []parser{&p496, &p15}
	var p508 = sequenceParser{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p507 = sequenceParser{id: 507, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p507.items = []parser{&p825, &p14}
	p508.items = []parser{&p825, &p14, &p507}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p509.items = []parser{&p825, &p14}
	p510.items = []parser{&p825, &p14, &p509}
	p511.items = []parser{&p497, &p508, &p825, &p272, &p510, &p825, &p272}
	var p564 = sequenceParser{id: 564, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 198, 793}}
	var p554 = sequenceParser{id: 554, commit: 74, name: "go-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p553 = sequenceParser{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p551 = charParser{id: 551, chars: []rune{103}}
	var p552 = charParser{id: 552, chars: []rune{111}}
	p553.items = []parser{&p551, &p552}
	p554.items = []parser{&p553, &p15}
	var p563 = sequenceParser{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p562 = sequenceParser{id: 562, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p562.items = []parser{&p825, &p14}
	p563.items = []parser{&p825, &p14, &p562}
	var p262 = sequenceParser{id: 262, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p259 = sequenceParser{id: 259, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p258 = charParser{id: 258, chars: []rune{40}}
	p259.items = []parser{&p258}
	var p261 = sequenceParser{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p260 = charParser{id: 260, chars: []rune{41}}
	p261.items = []parser{&p260}
	p262.items = []parser{&p272, &p825, &p259, &p825, &p114, &p825, &p119, &p825, &p114, &p825, &p261}
	p564.items = []parser{&p554, &p563, &p825, &p262}
	var p573 = sequenceParser{id: 573, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 198, 793}}
	var p570 = sequenceParser{id: 570, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p565 = charParser{id: 565, chars: []rune{100}}
	var p566 = charParser{id: 566, chars: []rune{101}}
	var p567 = charParser{id: 567, chars: []rune{102}}
	var p568 = charParser{id: 568, chars: []rune{101}}
	var p569 = charParser{id: 569, chars: []rune{114}}
	p570.items = []parser{&p565, &p566, &p567, &p568, &p569}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p571 = sequenceParser{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p571.items = []parser{&p825, &p14}
	p572.items = []parser{&p825, &p14, &p571}
	p573.items = []parser{&p570, &p572, &p825, &p262}
	var p637 = choiceParser{id: 637, commit: 64, name: "assignment", generalizations: []int{783, 198, 793}}
	var p621 = sequenceParser{id: 621, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{637, 783, 198, 793}}
	var p606 = sequenceParser{id: 606, commit: 74, name: "set-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p605 = sequenceParser{id: 605, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p602 = charParser{id: 602, chars: []rune{115}}
	var p603 = charParser{id: 603, chars: []rune{101}}
	var p604 = charParser{id: 604, chars: []rune{116}}
	p605.items = []parser{&p602, &p603, &p604}
	p606.items = []parser{&p605, &p15}
	var p620 = sequenceParser{id: 620, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p619 = sequenceParser{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p619.items = []parser{&p825, &p14}
	p620.items = []parser{&p825, &p14, &p619}
	var p614 = sequenceParser{id: 614, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p611 = sequenceParser{id: 611, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p610 = sequenceParser{id: 610, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p609 = sequenceParser{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p609.items = []parser{&p825, &p14}
	p610.items = []parser{&p14, &p609}
	var p608 = sequenceParser{id: 608, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p607 = charParser{id: 607, chars: []rune{61}}
	p608.items = []parser{&p607}
	p611.items = []parser{&p610, &p825, &p608}
	var p613 = sequenceParser{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p612 = sequenceParser{id: 612, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p612.items = []parser{&p825, &p14}
	p613.items = []parser{&p825, &p14, &p612}
	p614.items = []parser{&p272, &p825, &p611, &p613, &p825, &p401}
	p621.items = []parser{&p606, &p620, &p825, &p614}
	var p628 = sequenceParser{id: 628, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{637, 783, 198, 793}}
	var p625 = sequenceParser{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p624.items = []parser{&p825, &p14}
	p625.items = []parser{&p825, &p14, &p624}
	var p623 = sequenceParser{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p622 = charParser{id: 622, chars: []rune{61}}
	p623.items = []parser{&p622}
	var p627 = sequenceParser{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p626 = sequenceParser{id: 626, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p626.items = []parser{&p825, &p14}
	p627.items = []parser{&p825, &p14, &p626}
	p628.items = []parser{&p272, &p625, &p825, &p623, &p627, &p825, &p401}
	var p636 = sequenceParser{id: 636, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{637, 783, 198, 793}}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p634.items = []parser{&p825, &p14}
	p635.items = []parser{&p825, &p14, &p634}
	var p630 = sequenceParser{id: 630, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p629 = charParser{id: 629, chars: []rune{40}}
	p630.items = []parser{&p629}
	var p631 = sequenceParser{id: 631, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p618 = sequenceParser{id: 618, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p615.items = []parser{&p114, &p825, &p614}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p616.items = []parser{&p825, &p615}
	p617.items = []parser{&p825, &p615, &p616}
	p618.items = []parser{&p614, &p617}
	p631.items = []parser{&p114, &p825, &p618}
	var p633 = sequenceParser{id: 633, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p632 = charParser{id: 632, chars: []rune{41}}
	p633.items = []parser{&p632}
	p636.items = []parser{&p606, &p635, &p825, &p630, &p825, &p631, &p825, &p114, &p825, &p633}
	p637.options = []parser{&p621, &p628, &p636}
	var p792 = sequenceParser{id: 792, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 198, 793}}
	var p785 = sequenceParser{id: 785, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p784 = charParser{id: 784, chars: []rune{40}}
	p785.items = []parser{&p784}
	var p789 = sequenceParser{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p788 = sequenceParser{id: 788, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p788.items = []parser{&p825, &p14}
	p789.items = []parser{&p825, &p14, &p788}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p790.items = []parser{&p825, &p14}
	p791.items = []parser{&p825, &p14, &p790}
	var p787 = sequenceParser{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p786 = charParser{id: 786, chars: []rune{41}}
	p787.items = []parser{&p786}
	p792.items = []parser{&p785, &p789, &p825, &p783, &p791, &p825, &p787}
	p783.options = []parser{&p511, &p564, &p573, &p637, &p792, &p401}
	var p191 = sequenceParser{id: 191, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{198}}
	var p188 = sequenceParser{id: 188, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p187 = charParser{id: 187, chars: []rune{123}}
	p188.items = []parser{&p187}
	var p190 = sequenceParser{id: 190, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p189 = charParser{id: 189, chars: []rune{125}}
	p190.items = []parser{&p189}
	p191.items = []parser{&p188, &p825, &p807, &p825, &p811, &p825, &p807, &p825, &p190}
	p198.options = []parser{&p783, &p191}
	p201.items = []parser{&p193, &p825, &p114, &p825, &p195, &p825, &p114, &p825, &p197, &p200, &p825, &p198}
	p207.items = []parser{&p204, &p206, &p825, &p201}
	var p217 = sequenceParser{id: 217, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p210 = sequenceParser{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p208 = charParser{id: 208, chars: []rune{102}}
	var p209 = charParser{id: 209, chars: []rune{110}}
	p210.items = []parser{&p208, &p209}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p213.items = []parser{&p825, &p14}
	p214.items = []parser{&p825, &p14, &p213}
	var p212 = sequenceParser{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p211 = charParser{id: 211, chars: []rune{126}}
	p212.items = []parser{&p211}
	var p216 = sequenceParser{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p215.items = []parser{&p825, &p14}
	p216.items = []parser{&p825, &p14, &p215}
	p217.items = []parser{&p210, &p214, &p825, &p212, &p216, &p825, &p201}
	var p257 = sequenceParser{id: 257, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p256 = sequenceParser{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p255 = sequenceParser{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p255.items = []parser{&p825, &p14}
	p256.items = []parser{&p825, &p14, &p255}
	var p254 = sequenceParser{id: 254, commit: 66, name: "index-list", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var p250 = choiceParser{id: 250, commit: 66, name: "index"}
	var p231 = sequenceParser{id: 231, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{250}}
	var p228 = sequenceParser{id: 228, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p227 = charParser{id: 227, chars: []rune{46}}
	p228.items = []parser{&p227}
	var p230 = sequenceParser{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p229 = sequenceParser{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p229.items = []parser{&p825, &p14}
	p230.items = []parser{&p825, &p14, &p229}
	p231.items = []parser{&p228, &p230, &p825, &p104}
	var p240 = sequenceParser{id: 240, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{250}}
	var p233 = sequenceParser{id: 233, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p232 = charParser{id: 232, chars: []rune{91}}
	p233.items = []parser{&p232}
	var p237 = sequenceParser{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p236 = sequenceParser{id: 236, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p236.items = []parser{&p825, &p14}
	p237.items = []parser{&p825, &p14, &p236}
	var p239 = sequenceParser{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p238 = sequenceParser{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p238.items = []parser{&p825, &p14}
	p239.items = []parser{&p825, &p14, &p238}
	var p235 = sequenceParser{id: 235, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p234 = charParser{id: 234, chars: []rune{93}}
	p235.items = []parser{&p234}
	p240.items = []parser{&p233, &p237, &p825, &p401, &p239, &p825, &p235}
	var p249 = sequenceParser{id: 249, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{250}}
	var p242 = sequenceParser{id: 242, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p241 = charParser{id: 241, chars: []rune{91}}
	p242.items = []parser{&p241}
	var p246 = sequenceParser{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p245 = sequenceParser{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p245.items = []parser{&p825, &p14}
	p246.items = []parser{&p825, &p14, &p245}
	var p226 = sequenceParser{id: 226, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{583, 589, 590}}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p401}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p222.items = []parser{&p825, &p14}
	p223.items = []parser{&p825, &p14, &p222}
	var p221 = sequenceParser{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p220 = charParser{id: 220, chars: []rune{58}}
	p221.items = []parser{&p220}
	var p225 = sequenceParser{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p224.items = []parser{&p825, &p14}
	p225.items = []parser{&p825, &p14, &p224}
	var p219 = sequenceParser{id: 219, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p219.items = []parser{&p401}
	p226.items = []parser{&p218, &p223, &p825, &p221, &p225, &p825, &p219}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p247.items = []parser{&p825, &p14}
	p248.items = []parser{&p825, &p14, &p247}
	var p244 = sequenceParser{id: 244, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p243 = charParser{id: 243, chars: []rune{93}}
	p244.items = []parser{&p243}
	p249.items = []parser{&p242, &p246, &p825, &p226, &p248, &p825, &p244}
	p250.options = []parser{&p231, &p240, &p249}
	var p253 = sequenceParser{id: 253, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p252 = sequenceParser{id: 252, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p251 = sequenceParser{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p251.items = []parser{&p825, &p14}
	p252.items = []parser{&p14, &p251}
	p253.items = []parser{&p252, &p825, &p250}
	p254.items = []parser{&p250, &p825, &p253}
	p257.items = []parser{&p272, &p256, &p825, &p254}
	var p271 = sequenceParser{id: 271, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{272, 401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p264 = sequenceParser{id: 264, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p263 = charParser{id: 263, chars: []rune{40}}
	p264.items = []parser{&p263}
	var p268 = sequenceParser{id: 268, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p267 = sequenceParser{id: 267, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p267.items = []parser{&p825, &p14}
	p268.items = []parser{&p825, &p14, &p267}
	var p270 = sequenceParser{id: 270, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p269 = sequenceParser{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p269.items = []parser{&p825, &p14}
	p270.items = []parser{&p825, &p14, &p269}
	var p266 = sequenceParser{id: 266, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p265 = charParser{id: 265, chars: []rune{41}}
	p266.items = []parser{&p265}
	p271.items = []parser{&p264, &p268, &p825, &p401, &p270, &p825, &p266}
	p272.options = []parser{&p61, &p74, &p87, &p99, &p514, &p104, &p125, &p130, &p159, &p164, &p207, &p217, &p257, &p262, &p271}
	var p332 = sequenceParser{id: 332, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{401, 338, 339, 340, 341, 342, 393, 590, 583, 793}}
	var p331 = choiceParser{id: 331, commit: 66, name: "unary-operator"}
	var p291 = sequenceParser{id: 291, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{331}}
	var p290 = charParser{id: 290, chars: []rune{43}}
	p291.items = []parser{&p290}
	var p293 = sequenceParser{id: 293, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{331}}
	var p292 = charParser{id: 292, chars: []rune{45}}
	p293.items = []parser{&p292}
	var p274 = sequenceParser{id: 274, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{331}}
	var p273 = charParser{id: 273, chars: []rune{94}}
	p274.items = []parser{&p273}
	var p305 = sequenceParser{id: 305, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{331}}
	var p304 = charParser{id: 304, chars: []rune{33}}
	p305.items = []parser{&p304}
	p331.options = []parser{&p291, &p293, &p274, &p305}
	p332.items = []parser{&p331, &p825, &p272}
	var p379 = choiceParser{id: 379, commit: 66, name: "binary-expression", generalizations: []int{401, 393, 590, 583, 793}}
	var p350 = sequenceParser{id: 350, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{379, 339, 340, 341, 342, 401, 393, 590, 583, 793}}
	var p338 = choiceParser{id: 338, commit: 66, name: "operand0", generalizations: []int{339, 340, 341, 342}}
	p338.options = []parser{&p272, &p332}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p345 = sequenceParser{id: 345, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p344 = sequenceParser{id: 344, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p344.items = []parser{&p825, &p14}
	p345.items = []parser{&p14, &p344}
	var p333 = choiceParser{id: 333, commit: 66, name: "binary-op0"}
	var p276 = sequenceParser{id: 276, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p275 = charParser{id: 275, chars: []rune{38}}
	p276.items = []parser{&p275}
	var p283 = sequenceParser{id: 283, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{333}}
	var p281 = charParser{id: 281, chars: []rune{38}}
	var p282 = charParser{id: 282, chars: []rune{94}}
	p283.items = []parser{&p281, &p282}
	var p286 = sequenceParser{id: 286, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{333}}
	var p284 = charParser{id: 284, chars: []rune{60}}
	var p285 = charParser{id: 285, chars: []rune{60}}
	p286.items = []parser{&p284, &p285}
	var p289 = sequenceParser{id: 289, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{333}}
	var p287 = charParser{id: 287, chars: []rune{62}}
	var p288 = charParser{id: 288, chars: []rune{62}}
	p289.items = []parser{&p287, &p288}
	var p295 = sequenceParser{id: 295, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p294 = charParser{id: 294, chars: []rune{42}}
	p295.items = []parser{&p294}
	var p297 = sequenceParser{id: 297, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p296 = charParser{id: 296, chars: []rune{47}}
	p297.items = []parser{&p296}
	var p299 = sequenceParser{id: 299, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p298 = charParser{id: 298, chars: []rune{37}}
	p299.items = []parser{&p298}
	p333.options = []parser{&p276, &p283, &p286, &p289, &p295, &p297, &p299}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p346.items = []parser{&p825, &p14}
	p347.items = []parser{&p825, &p14, &p346}
	p348.items = []parser{&p345, &p825, &p333, &p347, &p825, &p338}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p349.items = []parser{&p825, &p348}
	p350.items = []parser{&p338, &p825, &p348, &p349}
	var p357 = sequenceParser{id: 357, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{379, 340, 341, 342, 401, 393, 590, 583, 793}}
	var p339 = choiceParser{id: 339, commit: 66, name: "operand1", generalizations: []int{340, 341, 342}}
	p339.options = []parser{&p338, &p350}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p351 = sequenceParser{id: 351, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p351.items = []parser{&p825, &p14}
	p352.items = []parser{&p14, &p351}
	var p334 = choiceParser{id: 334, commit: 66, name: "binary-op1"}
	var p278 = sequenceParser{id: 278, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p277 = charParser{id: 277, chars: []rune{124}}
	p278.items = []parser{&p277}
	var p280 = sequenceParser{id: 280, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p279 = charParser{id: 279, chars: []rune{94}}
	p280.items = []parser{&p279}
	var p301 = sequenceParser{id: 301, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p300 = charParser{id: 300, chars: []rune{43}}
	p301.items = []parser{&p300}
	var p303 = sequenceParser{id: 303, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p302 = charParser{id: 302, chars: []rune{45}}
	p303.items = []parser{&p302}
	p334.options = []parser{&p278, &p280, &p301, &p303}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p353.items = []parser{&p825, &p14}
	p354.items = []parser{&p825, &p14, &p353}
	p355.items = []parser{&p352, &p825, &p334, &p354, &p825, &p339}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p356.items = []parser{&p825, &p355}
	p357.items = []parser{&p339, &p825, &p355, &p356}
	var p364 = sequenceParser{id: 364, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{379, 341, 342, 401, 393, 590, 583, 793}}
	var p340 = choiceParser{id: 340, commit: 66, name: "operand2", generalizations: []int{341, 342}}
	p340.options = []parser{&p339, &p357}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p358.items = []parser{&p825, &p14}
	p359.items = []parser{&p14, &p358}
	var p335 = choiceParser{id: 335, commit: 66, name: "binary-op2"}
	var p308 = sequenceParser{id: 308, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p306 = charParser{id: 306, chars: []rune{61}}
	var p307 = charParser{id: 307, chars: []rune{61}}
	p308.items = []parser{&p306, &p307}
	var p311 = sequenceParser{id: 311, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p309 = charParser{id: 309, chars: []rune{33}}
	var p310 = charParser{id: 310, chars: []rune{61}}
	p311.items = []parser{&p309, &p310}
	var p313 = sequenceParser{id: 313, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p312 = charParser{id: 312, chars: []rune{60}}
	p313.items = []parser{&p312}
	var p316 = sequenceParser{id: 316, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p314 = charParser{id: 314, chars: []rune{60}}
	var p315 = charParser{id: 315, chars: []rune{61}}
	p316.items = []parser{&p314, &p315}
	var p318 = sequenceParser{id: 318, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p317 = charParser{id: 317, chars: []rune{62}}
	p318.items = []parser{&p317}
	var p321 = sequenceParser{id: 321, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p319 = charParser{id: 319, chars: []rune{62}}
	var p320 = charParser{id: 320, chars: []rune{61}}
	p321.items = []parser{&p319, &p320}
	p335.options = []parser{&p308, &p311, &p313, &p316, &p318, &p321}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p360.items = []parser{&p825, &p14}
	p361.items = []parser{&p825, &p14, &p360}
	p362.items = []parser{&p359, &p825, &p335, &p361, &p825, &p340}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p363.items = []parser{&p825, &p362}
	p364.items = []parser{&p340, &p825, &p362, &p363}
	var p371 = sequenceParser{id: 371, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{379, 342, 401, 393, 590, 583, 793}}
	var p341 = choiceParser{id: 341, commit: 66, name: "operand3", generalizations: []int{342}}
	p341.options = []parser{&p340, &p364}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p365 = sequenceParser{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p365.items = []parser{&p825, &p14}
	p366.items = []parser{&p14, &p365}
	var p336 = sequenceParser{id: 336, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p324 = sequenceParser{id: 324, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p322 = charParser{id: 322, chars: []rune{38}}
	var p323 = charParser{id: 323, chars: []rune{38}}
	p324.items = []parser{&p322, &p323}
	p336.items = []parser{&p324}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p367.items = []parser{&p825, &p14}
	p368.items = []parser{&p825, &p14, &p367}
	p369.items = []parser{&p366, &p825, &p336, &p368, &p825, &p341}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p370.items = []parser{&p825, &p369}
	p371.items = []parser{&p341, &p825, &p369, &p370}
	var p378 = sequenceParser{id: 378, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{379, 401, 393, 590, 583, 793}}
	var p342 = choiceParser{id: 342, commit: 66, name: "operand4"}
	p342.options = []parser{&p341, &p371}
	var p376 = sequenceParser{id: 376, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p372.items = []parser{&p825, &p14}
	p373.items = []parser{&p14, &p372}
	var p337 = sequenceParser{id: 337, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p327 = sequenceParser{id: 327, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p325 = charParser{id: 325, chars: []rune{124}}
	var p326 = charParser{id: 326, chars: []rune{124}}
	p327.items = []parser{&p325, &p326}
	p337.items = []parser{&p327}
	var p375 = sequenceParser{id: 375, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p374 = sequenceParser{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p374.items = []parser{&p825, &p14}
	p375.items = []parser{&p825, &p14, &p374}
	p376.items = []parser{&p373, &p825, &p337, &p375, &p825, &p342}
	var p377 = sequenceParser{id: 377, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p377.items = []parser{&p825, &p376}
	p378.items = []parser{&p342, &p825, &p376, &p377}
	p379.options = []parser{&p350, &p357, &p364, &p371, &p378}
	var p392 = sequenceParser{id: 392, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{401, 393, 590, 583, 793}}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p384 = sequenceParser{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p384.items = []parser{&p825, &p14}
	p385.items = []parser{&p825, &p14, &p384}
	var p381 = sequenceParser{id: 381, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p380 = charParser{id: 380, chars: []rune{63}}
	p381.items = []parser{&p380}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p386.items = []parser{&p825, &p14}
	p387.items = []parser{&p825, &p14, &p386}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p388.items = []parser{&p825, &p14}
	p389.items = []parser{&p825, &p14, &p388}
	var p383 = sequenceParser{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p382 = charParser{id: 382, chars: []rune{58}}
	p383.items = []parser{&p382}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p390.items = []parser{&p825, &p14}
	p391.items = []parser{&p825, &p14, &p390}
	p392.items = []parser{&p401, &p385, &p825, &p381, &p387, &p825, &p401, &p389, &p825, &p383, &p391, &p825, &p401}
	var p400 = sequenceParser{id: 400, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{401, 590, 583, 793}}
	var p393 = choiceParser{id: 393, commit: 66, name: "chainingOperand"}
	p393.options = []parser{&p272, &p332, &p379, &p392}
	var p398 = sequenceParser{id: 398, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p395 = sequenceParser{id: 395, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p394 = sequenceParser{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p394.items = []parser{&p825, &p14}
	p395.items = []parser{&p14, &p394}
	var p330 = sequenceParser{id: 330, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p328 = charParser{id: 328, chars: []rune{45}}
	var p329 = charParser{id: 329, chars: []rune{62}}
	p330.items = []parser{&p328, &p329}
	var p397 = sequenceParser{id: 397, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p396 = sequenceParser{id: 396, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p396.items = []parser{&p825, &p14}
	p397.items = []parser{&p825, &p14, &p396}
	p398.items = []parser{&p395, &p825, &p330, &p397, &p825, &p393}
	var p399 = sequenceParser{id: 399, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p399.items = []parser{&p825, &p398}
	p400.items = []parser{&p393, &p825, &p398, &p399}
	p401.options = []parser{&p272, &p332, &p379, &p392, &p400}
	p186.items = []parser{&p183, &p185, &p825, &p401}
	var p432 = sequenceParser{id: 432, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{793, 480, 541}}
	var p405 = sequenceParser{id: 405, commit: 74, name: "if-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p404 = sequenceParser{id: 404, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p402 = charParser{id: 402, chars: []rune{105}}
	var p403 = charParser{id: 403, chars: []rune{102}}
	p404.items = []parser{&p402, &p403}
	p405.items = []parser{&p404, &p15}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p426.items = []parser{&p825, &p14}
	p427.items = []parser{&p825, &p14, &p426}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p428.items = []parser{&p825, &p14}
	p429.items = []parser{&p825, &p14, &p428}
	var p431 = sequenceParser{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p420 = sequenceParser{id: 420, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p413 = sequenceParser{id: 413, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p412 = sequenceParser{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p412.items = []parser{&p825, &p14}
	p413.items = []parser{&p14, &p412}
	var p411 = sequenceParser{id: 411, commit: 74, name: "else-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p410 = sequenceParser{id: 410, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p406 = charParser{id: 406, chars: []rune{101}}
	var p407 = charParser{id: 407, chars: []rune{108}}
	var p408 = charParser{id: 408, chars: []rune{115}}
	var p409 = charParser{id: 409, chars: []rune{101}}
	p410.items = []parser{&p406, &p407, &p408, &p409}
	p411.items = []parser{&p410, &p15}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p414.items = []parser{&p825, &p14}
	p415.items = []parser{&p825, &p14, &p414}
	var p417 = sequenceParser{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p416 = sequenceParser{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p416.items = []parser{&p825, &p14}
	p417.items = []parser{&p825, &p14, &p416}
	var p419 = sequenceParser{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p418 = sequenceParser{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p418.items = []parser{&p825, &p14}
	p419.items = []parser{&p825, &p14, &p418}
	p420.items = []parser{&p413, &p825, &p411, &p415, &p825, &p405, &p417, &p825, &p401, &p419, &p825, &p191}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p430.items = []parser{&p825, &p420}
	p431.items = []parser{&p825, &p420, &p430}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p421 = sequenceParser{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p421.items = []parser{&p825, &p14}
	p422.items = []parser{&p14, &p421}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p423.items = []parser{&p825, &p14}
	p424.items = []parser{&p825, &p14, &p423}
	p425.items = []parser{&p422, &p825, &p411, &p424, &p825, &p191}
	p432.items = []parser{&p405, &p427, &p825, &p401, &p429, &p825, &p191, &p431, &p825, &p425}
	var p491 = sequenceParser{id: 491, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{480, 793, 541}}
	var p446 = sequenceParser{id: 446, commit: 74, name: "switch-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p445 = sequenceParser{id: 445, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p439 = charParser{id: 439, chars: []rune{115}}
	var p440 = charParser{id: 440, chars: []rune{119}}
	var p441 = charParser{id: 441, chars: []rune{105}}
	var p442 = charParser{id: 442, chars: []rune{116}}
	var p443 = charParser{id: 443, chars: []rune{99}}
	var p444 = charParser{id: 444, chars: []rune{104}}
	p445.items = []parser{&p439, &p440, &p441, &p442, &p443, &p444}
	p446.items = []parser{&p445, &p15}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p487.items = []parser{&p825, &p14}
	p488.items = []parser{&p825, &p14, &p487}
	var p490 = sequenceParser{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p489.items = []parser{&p825, &p14}
	p490.items = []parser{&p825, &p14, &p489}
	var p478 = sequenceParser{id: 478, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p477 = charParser{id: 477, chars: []rune{123}}
	p478.items = []parser{&p477}
	var p484 = sequenceParser{id: 484, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p479 = choiceParser{id: 479, commit: 2}
	var p476 = sequenceParser{id: 476, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{479, 480}}
	var p471 = sequenceParser{id: 471, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p438 = sequenceParser{id: 438, commit: 74, name: "case-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p437 = sequenceParser{id: 437, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p433 = charParser{id: 433, chars: []rune{99}}
	var p434 = charParser{id: 434, chars: []rune{97}}
	var p435 = charParser{id: 435, chars: []rune{115}}
	var p436 = charParser{id: 436, chars: []rune{101}}
	p437.items = []parser{&p433, &p434, &p435, &p436}
	p438.items = []parser{&p437, &p15}
	var p468 = sequenceParser{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p467 = sequenceParser{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p467.items = []parser{&p825, &p14}
	p468.items = []parser{&p825, &p14, &p467}
	var p470 = sequenceParser{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p469 = sequenceParser{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p469.items = []parser{&p825, &p14}
	p470.items = []parser{&p825, &p14, &p469}
	var p466 = sequenceParser{id: 466, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p465 = charParser{id: 465, chars: []rune{58}}
	p466.items = []parser{&p465}
	p471.items = []parser{&p438, &p468, &p825, &p401, &p470, &p825, &p466}
	var p475 = sequenceParser{id: 475, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p473 = sequenceParser{id: 473, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p472 = charParser{id: 472, chars: []rune{59}}
	p473.items = []parser{&p472}
	var p474 = sequenceParser{id: 474, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p474.items = []parser{&p825, &p473}
	p475.items = []parser{&p825, &p473, &p474}
	p476.items = []parser{&p471, &p475, &p825, &p793}
	var p464 = sequenceParser{id: 464, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{479, 480, 540, 541}}
	var p459 = sequenceParser{id: 459, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p454 = sequenceParser{id: 454, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p447 = charParser{id: 447, chars: []rune{100}}
	var p448 = charParser{id: 448, chars: []rune{101}}
	var p449 = charParser{id: 449, chars: []rune{102}}
	var p450 = charParser{id: 450, chars: []rune{97}}
	var p451 = charParser{id: 451, chars: []rune{117}}
	var p452 = charParser{id: 452, chars: []rune{108}}
	var p453 = charParser{id: 453, chars: []rune{116}}
	p454.items = []parser{&p447, &p448, &p449, &p450, &p451, &p452, &p453}
	var p458 = sequenceParser{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p457 = sequenceParser{id: 457, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p457.items = []parser{&p825, &p14}
	p458.items = []parser{&p825, &p14, &p457}
	var p456 = sequenceParser{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p455 = charParser{id: 455, chars: []rune{58}}
	p456.items = []parser{&p455}
	p459.items = []parser{&p454, &p458, &p825, &p456}
	var p463 = sequenceParser{id: 463, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p461 = sequenceParser{id: 461, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p460 = charParser{id: 460, chars: []rune{59}}
	p461.items = []parser{&p460}
	var p462 = sequenceParser{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p462.items = []parser{&p825, &p461}
	p463.items = []parser{&p825, &p461, &p462}
	p464.items = []parser{&p459, &p463, &p825, &p793}
	p479.options = []parser{&p476, &p464}
	var p483 = sequenceParser{id: 483, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p481 = sequenceParser{id: 481, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p480 = choiceParser{id: 480, commit: 2}
	p480.options = []parser{&p476, &p464, &p793}
	p481.items = []parser{&p807, &p825, &p480}
	var p482 = sequenceParser{id: 482, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p482.items = []parser{&p825, &p481}
	p483.items = []parser{&p825, &p481, &p482}
	p484.items = []parser{&p479, &p483}
	var p486 = sequenceParser{id: 486, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p485 = charParser{id: 485, chars: []rune{125}}
	p486.items = []parser{&p485}
	p491.items = []parser{&p446, &p488, &p825, &p401, &p490, &p825, &p478, &p825, &p807, &p825, &p484, &p825, &p807, &p825, &p486}
	var p550 = sequenceParser{id: 550, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{541, 793}}
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
	p548.items = []parser{&p825, &p14}
	p549.items = []parser{&p825, &p14, &p548}
	var p539 = sequenceParser{id: 539, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p538 = charParser{id: 538, chars: []rune{123}}
	p539.items = []parser{&p538}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p540 = choiceParser{id: 540, commit: 2}
	var p530 = sequenceParser{id: 530, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{540, 541}}
	var p525 = sequenceParser{id: 525, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p522 = sequenceParser{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p521 = sequenceParser{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p521.items = []parser{&p825, &p14}
	p522.items = []parser{&p825, &p14, &p521}
	var p518 = choiceParser{id: 518, commit: 66, name: "communication"}
	var p517 = sequenceParser{id: 517, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{518}}
	var p516 = sequenceParser{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p515 = sequenceParser{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p515.items = []parser{&p825, &p14}
	p516.items = []parser{&p825, &p14, &p515}
	p517.items = []parser{&p104, &p516, &p825, &p514}
	p518.options = []parser{&p511, &p514, &p517}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p523.items = []parser{&p825, &p14}
	p524.items = []parser{&p825, &p14, &p523}
	var p520 = sequenceParser{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p519 = charParser{id: 519, chars: []rune{58}}
	p520.items = []parser{&p519}
	p525.items = []parser{&p438, &p522, &p825, &p518, &p524, &p825, &p520}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p527 = sequenceParser{id: 527, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p526 = charParser{id: 526, chars: []rune{59}}
	p527.items = []parser{&p526}
	var p528 = sequenceParser{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p528.items = []parser{&p825, &p527}
	p529.items = []parser{&p825, &p527, &p528}
	p530.items = []parser{&p525, &p529, &p825, &p793}
	p540.options = []parser{&p530, &p464}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p542 = sequenceParser{id: 542, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p541 = choiceParser{id: 541, commit: 2}
	p541.options = []parser{&p530, &p464, &p793}
	p542.items = []parser{&p807, &p825, &p541}
	var p543 = sequenceParser{id: 543, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p543.items = []parser{&p825, &p542}
	p544.items = []parser{&p825, &p542, &p543}
	p545.items = []parser{&p540, &p544}
	var p547 = sequenceParser{id: 547, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p546 = charParser{id: 546, chars: []rune{125}}
	p547.items = []parser{&p546}
	p550.items = []parser{&p537, &p549, &p825, &p539, &p825, &p807, &p825, &p545, &p825, &p807, &p825, &p547}
	var p601 = sequenceParser{id: 601, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{793}}
	var p582 = sequenceParser{id: 582, commit: 74, name: "for-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p581 = sequenceParser{id: 581, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p578 = charParser{id: 578, chars: []rune{102}}
	var p579 = charParser{id: 579, chars: []rune{111}}
	var p580 = charParser{id: 580, chars: []rune{114}}
	p581.items = []parser{&p578, &p579, &p580}
	p582.items = []parser{&p581, &p15}
	var p600 = choiceParser{id: 600, commit: 2}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{600}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p592 = sequenceParser{id: 592, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p591 = sequenceParser{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p591.items = []parser{&p825, &p14}
	p592.items = []parser{&p14, &p591}
	var p590 = choiceParser{id: 590, commit: 66, name: "loop-expression"}
	var p589 = choiceParser{id: 589, commit: 64, name: "range-over-expression", generalizations: []int{590}}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{589, 590}}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p584.items = []parser{&p825, &p14}
	p585.items = []parser{&p825, &p14, &p584}
	var p577 = sequenceParser{id: 577, commit: 74, name: "in-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p576 = sequenceParser{id: 576, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p574 = charParser{id: 574, chars: []rune{105}}
	var p575 = charParser{id: 575, chars: []rune{110}}
	p576.items = []parser{&p574, &p575}
	p577.items = []parser{&p576, &p15}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p586.items = []parser{&p825, &p14}
	p587.items = []parser{&p825, &p14, &p586}
	var p583 = choiceParser{id: 583, commit: 2}
	p583.options = []parser{&p401, &p226}
	p588.items = []parser{&p104, &p585, &p825, &p577, &p587, &p825, &p583}
	p589.options = []parser{&p588, &p226}
	p590.options = []parser{&p401, &p589}
	p593.items = []parser{&p592, &p825, &p590}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p594.items = []parser{&p825, &p14}
	p595.items = []parser{&p825, &p14, &p594}
	p596.items = []parser{&p593, &p595, &p825, &p191}
	var p599 = sequenceParser{id: 599, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{600}}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p597.items = []parser{&p825, &p14}
	p598.items = []parser{&p14, &p597}
	p599.items = []parser{&p598, &p825, &p191}
	p600.options = []parser{&p596, &p599}
	p601.items = []parser{&p582, &p825, &p600}
	var p739 = choiceParser{id: 739, commit: 66, name: "definition", generalizations: []int{793}}
	var p660 = sequenceParser{id: 660, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 793}}
	var p642 = sequenceParser{id: 642, commit: 74, name: "let-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p641 = sequenceParser{id: 641, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p638 = charParser{id: 638, chars: []rune{108}}
	var p639 = charParser{id: 639, chars: []rune{101}}
	var p640 = charParser{id: 640, chars: []rune{116}}
	p641.items = []parser{&p638, &p639, &p640}
	p642.items = []parser{&p641, &p15}
	var p659 = sequenceParser{id: 659, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p658 = sequenceParser{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p658.items = []parser{&p825, &p14}
	p659.items = []parser{&p825, &p14, &p658}
	var p657 = choiceParser{id: 657, commit: 2}
	var p651 = sequenceParser{id: 651, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{657, 661, 662}}
	var p650 = sequenceParser{id: 650, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p647 = sequenceParser{id: 647, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p646 = sequenceParser{id: 646, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p645 = sequenceParser{id: 645, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p645.items = []parser{&p825, &p14}
	p646.items = []parser{&p14, &p645}
	var p644 = sequenceParser{id: 644, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p643 = charParser{id: 643, chars: []rune{61}}
	p644.items = []parser{&p643}
	p647.items = []parser{&p646, &p825, &p644}
	var p649 = sequenceParser{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p648 = sequenceParser{id: 648, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p648.items = []parser{&p825, &p14}
	p649.items = []parser{&p825, &p14, &p648}
	p650.items = []parser{&p104, &p825, &p647, &p649, &p825, &p401}
	p651.items = []parser{&p650}
	var p656 = sequenceParser{id: 656, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{657, 661, 662}}
	var p653 = sequenceParser{id: 653, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p652 = charParser{id: 652, chars: []rune{126}}
	p653.items = []parser{&p652}
	var p655 = sequenceParser{id: 655, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p654 = sequenceParser{id: 654, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p654.items = []parser{&p825, &p14}
	p655.items = []parser{&p825, &p14, &p654}
	p656.items = []parser{&p653, &p655, &p825, &p650}
	p657.options = []parser{&p651, &p656}
	p660.items = []parser{&p642, &p659, &p825, &p657}
	var p677 = sequenceParser{id: 677, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 793}}
	var p676 = sequenceParser{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p675 = sequenceParser{id: 675, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p675.items = []parser{&p825, &p14}
	p676.items = []parser{&p825, &p14, &p675}
	var p672 = sequenceParser{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p671 = charParser{id: 671, chars: []rune{40}}
	p672.items = []parser{&p671}
	var p666 = sequenceParser{id: 666, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p661 = choiceParser{id: 661, commit: 2}
	p661.options = []parser{&p651, &p656}
	var p665 = sequenceParser{id: 665, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p663 = sequenceParser{id: 663, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p662 = choiceParser{id: 662, commit: 2}
	p662.options = []parser{&p651, &p656}
	p663.items = []parser{&p114, &p825, &p662}
	var p664 = sequenceParser{id: 664, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p664.items = []parser{&p825, &p663}
	p665.items = []parser{&p825, &p663, &p664}
	p666.items = []parser{&p661, &p665}
	var p674 = sequenceParser{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p673 = charParser{id: 673, chars: []rune{41}}
	p674.items = []parser{&p673}
	p677.items = []parser{&p642, &p676, &p825, &p672, &p825, &p114, &p825, &p666, &p825, &p114, &p825, &p674}
	var p688 = sequenceParser{id: 688, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 793}}
	var p685 = sequenceParser{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p684 = sequenceParser{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p684.items = []parser{&p825, &p14}
	p685.items = []parser{&p825, &p14, &p684}
	var p679 = sequenceParser{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p678 = charParser{id: 678, chars: []rune{126}}
	p679.items = []parser{&p678}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p686.items = []parser{&p825, &p14}
	p687.items = []parser{&p825, &p14, &p686}
	var p681 = sequenceParser{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{40}}
	p681.items = []parser{&p680}
	var p670 = sequenceParser{id: 670, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p669 = sequenceParser{id: 669, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p667.items = []parser{&p114, &p825, &p651}
	var p668 = sequenceParser{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p668.items = []parser{&p825, &p667}
	p669.items = []parser{&p825, &p667, &p668}
	p670.items = []parser{&p651, &p669}
	var p683 = sequenceParser{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p682 = charParser{id: 682, chars: []rune{41}}
	p683.items = []parser{&p682}
	p688.items = []parser{&p642, &p685, &p825, &p679, &p687, &p825, &p681, &p825, &p114, &p825, &p670, &p825, &p114, &p825, &p683}
	var p704 = sequenceParser{id: 704, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 793}}
	var p700 = sequenceParser{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p698 = charParser{id: 698, chars: []rune{102}}
	var p699 = charParser{id: 699, chars: []rune{110}}
	p700.items = []parser{&p698, &p699}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p702 = sequenceParser{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p702.items = []parser{&p825, &p14}
	p703.items = []parser{&p825, &p14, &p702}
	var p701 = choiceParser{id: 701, commit: 2}
	var p692 = sequenceParser{id: 692, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{701, 709, 710}}
	var p691 = sequenceParser{id: 691, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p690 = sequenceParser{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p689 = sequenceParser{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p689.items = []parser{&p825, &p14}
	p690.items = []parser{&p825, &p14, &p689}
	p691.items = []parser{&p104, &p690, &p825, &p201}
	p692.items = []parser{&p691}
	var p697 = sequenceParser{id: 697, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{701, 709, 710}}
	var p694 = sequenceParser{id: 694, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p693 = charParser{id: 693, chars: []rune{126}}
	p694.items = []parser{&p693}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p695 = sequenceParser{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p695.items = []parser{&p825, &p14}
	p696.items = []parser{&p825, &p14, &p695}
	p697.items = []parser{&p694, &p696, &p825, &p691}
	p701.options = []parser{&p692, &p697}
	p704.items = []parser{&p700, &p703, &p825, &p701}
	var p724 = sequenceParser{id: 724, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 793}}
	var p717 = sequenceParser{id: 717, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p715 = charParser{id: 715, chars: []rune{102}}
	var p716 = charParser{id: 716, chars: []rune{110}}
	p717.items = []parser{&p715, &p716}
	var p723 = sequenceParser{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p722 = sequenceParser{id: 722, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p722.items = []parser{&p825, &p14}
	p723.items = []parser{&p825, &p14, &p722}
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
	p711.items = []parser{&p114, &p825, &p710}
	var p712 = sequenceParser{id: 712, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p712.items = []parser{&p825, &p711}
	p713.items = []parser{&p825, &p711, &p712}
	p714.items = []parser{&p709, &p713}
	var p721 = sequenceParser{id: 721, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p720 = charParser{id: 720, chars: []rune{41}}
	p721.items = []parser{&p720}
	p724.items = []parser{&p717, &p723, &p825, &p719, &p825, &p114, &p825, &p714, &p825, &p114, &p825, &p721}
	var p738 = sequenceParser{id: 738, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{739, 793}}
	var p727 = sequenceParser{id: 727, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p725 = charParser{id: 725, chars: []rune{102}}
	var p726 = charParser{id: 726, chars: []rune{110}}
	p727.items = []parser{&p725, &p726}
	var p735 = sequenceParser{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p734 = sequenceParser{id: 734, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p734.items = []parser{&p825, &p14}
	p735.items = []parser{&p825, &p14, &p734}
	var p729 = sequenceParser{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p728 = charParser{id: 728, chars: []rune{126}}
	p729.items = []parser{&p728}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p736.items = []parser{&p825, &p14}
	p737.items = []parser{&p825, &p14, &p736}
	var p731 = sequenceParser{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p730 = charParser{id: 730, chars: []rune{40}}
	p731.items = []parser{&p730}
	var p708 = sequenceParser{id: 708, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p707 = sequenceParser{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p705 = sequenceParser{id: 705, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p705.items = []parser{&p114, &p825, &p692}
	var p706 = sequenceParser{id: 706, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p706.items = []parser{&p825, &p705}
	p707.items = []parser{&p825, &p705, &p706}
	p708.items = []parser{&p692, &p707}
	var p733 = sequenceParser{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p732 = charParser{id: 732, chars: []rune{41}}
	p733.items = []parser{&p732}
	p738.items = []parser{&p727, &p735, &p825, &p729, &p737, &p825, &p731, &p825, &p114, &p825, &p708, &p825, &p114, &p825, &p733}
	p739.options = []parser{&p660, &p677, &p688, &p704, &p724, &p738}
	var p771 = choiceParser{id: 771, commit: 64, name: "use", generalizations: []int{793}}
	var p763 = sequenceParser{id: 763, commit: 66, name: "use-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{771, 793}}
	var p744 = sequenceParser{id: 744, commit: 74, name: "use-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p743 = sequenceParser{id: 743, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p740 = charParser{id: 740, chars: []rune{117}}
	var p741 = charParser{id: 741, chars: []rune{115}}
	var p742 = charParser{id: 742, chars: []rune{101}}
	p743.items = []parser{&p740, &p741, &p742}
	p744.items = []parser{&p743, &p15}
	var p762 = sequenceParser{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p761 = sequenceParser{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p761.items = []parser{&p825, &p14}
	p762.items = []parser{&p825, &p14, &p761}
	var p756 = choiceParser{id: 756, commit: 64, name: "use-fact"}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{756}}
	var p747 = choiceParser{id: 747, commit: 2}
	var p746 = sequenceParser{id: 746, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{747}}
	var p745 = charParser{id: 745, chars: []rune{46}}
	p746.items = []parser{&p745}
	p747.options = []parser{&p104, &p746}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p751 = sequenceParser{id: 751, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p750 = sequenceParser{id: 750, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p750.items = []parser{&p825, &p14}
	p751.items = []parser{&p14, &p750}
	var p749 = sequenceParser{id: 749, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p748 = charParser{id: 748, chars: []rune{61}}
	p749.items = []parser{&p748}
	p752.items = []parser{&p751, &p825, &p749}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p753.items = []parser{&p825, &p14}
	p754.items = []parser{&p825, &p14, &p753}
	p755.items = []parser{&p747, &p825, &p752, &p754, &p825, &p87}
	p756.options = []parser{&p87, &p755}
	p763.items = []parser{&p744, &p762, &p825, &p756}
	var p770 = sequenceParser{id: 770, commit: 66, name: "use-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{771, 793}}
	var p769 = sequenceParser{id: 769, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p768 = sequenceParser{id: 768, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p768.items = []parser{&p825, &p14}
	p769.items = []parser{&p825, &p14, &p768}
	var p765 = sequenceParser{id: 765, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p764 = charParser{id: 764, chars: []rune{40}}
	p765.items = []parser{&p764}
	var p760 = sequenceParser{id: 760, commit: 66, name: "use-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p759 = sequenceParser{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p757 = sequenceParser{id: 757, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p757.items = []parser{&p114, &p825, &p756}
	var p758 = sequenceParser{id: 758, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p758.items = []parser{&p825, &p757}
	p759.items = []parser{&p825, &p757, &p758}
	p760.items = []parser{&p756, &p759}
	var p767 = sequenceParser{id: 767, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p766 = charParser{id: 766, chars: []rune{41}}
	p767.items = []parser{&p766}
	p770.items = []parser{&p744, &p769, &p825, &p765, &p825, &p114, &p825, &p760, &p825, &p114, &p825, &p767}
	p771.options = []parser{&p763, &p770}
	var p782 = sequenceParser{id: 782, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793}}
	var p779 = sequenceParser{id: 779, commit: 74, name: "export-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p778 = sequenceParser{id: 778, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p772 = charParser{id: 772, chars: []rune{101}}
	var p773 = charParser{id: 773, chars: []rune{120}}
	var p774 = charParser{id: 774, chars: []rune{112}}
	var p775 = charParser{id: 775, chars: []rune{111}}
	var p776 = charParser{id: 776, chars: []rune{114}}
	var p777 = charParser{id: 777, chars: []rune{116}}
	p778.items = []parser{&p772, &p773, &p774, &p775, &p776, &p777}
	p779.items = []parser{&p778, &p15}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p780 = sequenceParser{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p780.items = []parser{&p825, &p14}
	p781.items = []parser{&p825, &p14, &p780}
	p782.items = []parser{&p779, &p781, &p825, &p739}
	var p802 = sequenceParser{id: 802, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{793}}
	var p795 = sequenceParser{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p794 = charParser{id: 794, chars: []rune{40}}
	p795.items = []parser{&p794}
	var p799 = sequenceParser{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p798 = sequenceParser{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p798.items = []parser{&p825, &p14}
	p799.items = []parser{&p825, &p14, &p798}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p800.items = []parser{&p825, &p14}
	p801.items = []parser{&p825, &p14, &p800}
	var p797 = sequenceParser{id: 797, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p796 = charParser{id: 796, chars: []rune{41}}
	p797.items = []parser{&p796}
	p802.items = []parser{&p795, &p799, &p825, &p793, &p801, &p825, &p797}
	p793.options = []parser{&p186, &p432, &p491, &p550, &p601, &p739, &p771, &p782, &p802, &p783}
	var p810 = sequenceParser{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p808 = sequenceParser{id: 808, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p808.items = []parser{&p807, &p825, &p793}
	var p809 = sequenceParser{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p809.items = []parser{&p825, &p808}
	p810.items = []parser{&p825, &p808, &p809}
	p811.items = []parser{&p793, &p810}
	p826.items = []parser{&p822, &p825, &p807, &p825, &p811, &p825, &p807}
	p827.items = []parser{&p825, &p826, &p825}
	var b827 = sequenceBuilder{id: 827, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b825 = choiceBuilder{id: 825, commit: 2}
	var b823 = choiceBuilder{id: 823, commit: 70}
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
	b823.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b824 = sequenceBuilder{id: 824, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
	var b43 = sequenceBuilder{id: 43, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b39 = choiceBuilder{id: 39, commit: 66}
	var b22 = sequenceBuilder{id: 22, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b21 = sequenceBuilder{id: 21, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b19 = charBuilder{}
	var b20 = charBuilder{}
	b21.items = []builder{&b19, &b20}
	var b18 = sequenceBuilder{id: 18, commit: 72, name: "line-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var b17 = sequenceBuilder{id: 17, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b16 = charBuilder{}
	b17.items = []builder{&b16}
	b18.items = []builder{&b17}
	b22.items = []builder{&b21, &b18}
	var b38 = sequenceBuilder{id: 38, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b34 = sequenceBuilder{id: 34, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b32 = charBuilder{}
	var b33 = charBuilder{}
	b34.items = []builder{&b32, &b33}
	var b31 = sequenceBuilder{id: 31, commit: 72, name: "block-comment-content", ranges: [][]int{{0, -1}, {0, -1}}}
	var b30 = choiceBuilder{id: 30, commit: 10}
	var b24 = sequenceBuilder{id: 24, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b23 = charBuilder{}
	b24.items = []builder{&b23}
	var b29 = sequenceBuilder{id: 29, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b26 = sequenceBuilder{id: 26, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b25 = charBuilder{}
	b26.items = []builder{&b25}
	var b28 = sequenceBuilder{id: 28, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b27 = charBuilder{}
	b28.items = []builder{&b27}
	b29.items = []builder{&b26, &b28}
	b30.options = []builder{&b24, &b29}
	b31.items = []builder{&b30}
	var b37 = sequenceBuilder{id: 37, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b35 = charBuilder{}
	var b36 = charBuilder{}
	b37.items = []builder{&b35, &b36}
	b38.items = []builder{&b34, &b31, &b37}
	b39.options = []builder{&b22, &b38}
	var b42 = sequenceBuilder{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b40 = sequenceBuilder{id: 40, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b14 = sequenceBuilder{id: 14, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b40.items = []builder{&b14, &b825, &b39}
	var b41 = sequenceBuilder{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b41.items = []builder{&b825, &b40}
	b42.items = []builder{&b825, &b40, &b41}
	b43.items = []builder{&b39, &b42}
	b824.items = []builder{&b43}
	b825.options = []builder{&b823, &b824}
	var b826 = sequenceBuilder{id: 826, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b822 = sequenceBuilder{id: 822, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b819 = sequenceBuilder{id: 819, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b817 = charBuilder{}
	var b818 = charBuilder{}
	b819.items = []builder{&b817, &b818}
	var b816 = sequenceBuilder{id: 816, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b815 = sequenceBuilder{id: 815, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b813 = sequenceBuilder{id: 813, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b812 = charBuilder{}
	b813.items = []builder{&b812}
	var b814 = sequenceBuilder{id: 814, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b814.items = []builder{&b825, &b813}
	b815.items = []builder{&b813, &b814}
	b816.items = []builder{&b815}
	var b821 = sequenceBuilder{id: 821, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b820 = charBuilder{}
	b821.items = []builder{&b820}
	b822.items = []builder{&b819, &b825, &b816, &b825, &b821}
	var b807 = sequenceBuilder{id: 807, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b805 = choiceBuilder{id: 805, commit: 2}
	var b804 = sequenceBuilder{id: 804, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b803 = charBuilder{}
	b804.items = []builder{&b803}
	b805.options = []builder{&b804, &b14}
	var b806 = sequenceBuilder{id: 806, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b806.items = []builder{&b825, &b805}
	b807.items = []builder{&b805, &b806}
	var b811 = sequenceBuilder{id: 811, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b793 = choiceBuilder{id: 793, commit: 66}
	var b186 = sequenceBuilder{id: 186, commit: 64, name: "return-value", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b183 = sequenceBuilder{id: 183, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b182 = sequenceBuilder{id: 182, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b176 = charBuilder{}
	var b177 = charBuilder{}
	var b178 = charBuilder{}
	var b179 = charBuilder{}
	var b180 = charBuilder{}
	var b181 = charBuilder{}
	b182.items = []builder{&b176, &b177, &b178, &b179, &b180, &b181}
	var b15 = choiceBuilder{id: 15, commit: 66}
	b15.options = []builder{&b823, &b14}
	b183.items = []builder{&b182, &b15}
	var b185 = sequenceBuilder{id: 185, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b184 = sequenceBuilder{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b184.items = []builder{&b825, &b14}
	b185.items = []builder{&b825, &b14, &b184}
	var b401 = choiceBuilder{id: 401, commit: 66}
	var b272 = choiceBuilder{id: 272, commit: 66}
	var b61 = choiceBuilder{id: 61, commit: 64, name: "int"}
	var b52 = sequenceBuilder{id: 52, commit: 74, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b51 = sequenceBuilder{id: 51, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b50 = charBuilder{}
	b51.items = []builder{&b50}
	var b45 = sequenceBuilder{id: 45, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b44 = charBuilder{}
	b45.items = []builder{&b44}
	b52.items = []builder{&b51, &b45}
	var b55 = sequenceBuilder{id: 55, commit: 74, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b54 = sequenceBuilder{id: 54, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b53 = charBuilder{}
	b54.items = []builder{&b53}
	var b47 = sequenceBuilder{id: 47, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b46 = charBuilder{}
	b47.items = []builder{&b46}
	b55.items = []builder{&b54, &b47}
	var b60 = sequenceBuilder{id: 60, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}}
	var b57 = sequenceBuilder{id: 57, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b56 = charBuilder{}
	b57.items = []builder{&b56}
	var b59 = sequenceBuilder{id: 59, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b58 = charBuilder{}
	b59.items = []builder{&b58}
	var b49 = sequenceBuilder{id: 49, commit: 66, allChars: true, ranges: [][]int{{1, 1}}}
	var b48 = charBuilder{}
	b49.items = []builder{&b48}
	b60.items = []builder{&b57, &b59, &b49}
	b61.options = []builder{&b52, &b55, &b60}
	var b74 = choiceBuilder{id: 74, commit: 72, name: "float"}
	var b69 = sequenceBuilder{id: 69, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}}
	var b68 = sequenceBuilder{id: 68, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b67 = charBuilder{}
	b68.items = []builder{&b67}
	var b66 = sequenceBuilder{id: 66, commit: 74, ranges: [][]int{{1, 1}, {0, 1}, {1, -1}, {1, 1}, {0, 1}, {1, -1}}}
	var b63 = sequenceBuilder{id: 63, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b62 = charBuilder{}
	b63.items = []builder{&b62}
	var b65 = sequenceBuilder{id: 65, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b64 = charBuilder{}
	b65.items = []builder{&b64}
	b66.items = []builder{&b63, &b65, &b45}
	b69.items = []builder{&b45, &b68, &b45, &b66}
	var b72 = sequenceBuilder{id: 72, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}}
	var b71 = sequenceBuilder{id: 71, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b70 = charBuilder{}
	b71.items = []builder{&b70}
	b72.items = []builder{&b71, &b45, &b66}
	var b73 = sequenceBuilder{id: 73, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}}
	b73.items = []builder{&b45, &b66}
	b74.options = []builder{&b69, &b72, &b73}
	var b87 = sequenceBuilder{id: 87, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}}
	var b76 = sequenceBuilder{id: 76, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b75 = charBuilder{}
	b76.items = []builder{&b75}
	var b84 = choiceBuilder{id: 84, commit: 10}
	var b78 = sequenceBuilder{id: 78, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b77 = charBuilder{}
	b78.items = []builder{&b77}
	var b83 = sequenceBuilder{id: 83, commit: 10, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b80 = sequenceBuilder{id: 80, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b79 = charBuilder{}
	b80.items = []builder{&b79}
	var b82 = sequenceBuilder{id: 82, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b81 = charBuilder{}
	b82.items = []builder{&b81}
	b83.items = []builder{&b80, &b82}
	b84.options = []builder{&b78, &b83}
	var b86 = sequenceBuilder{id: 86, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b85 = charBuilder{}
	b86.items = []builder{&b85}
	b87.items = []builder{&b76, &b84, &b86}
	var b99 = choiceBuilder{id: 99, commit: 66}
	var b92 = sequenceBuilder{id: 92, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b88 = charBuilder{}
	var b89 = charBuilder{}
	var b90 = charBuilder{}
	var b91 = charBuilder{}
	b92.items = []builder{&b88, &b89, &b90, &b91}
	var b98 = sequenceBuilder{id: 98, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b93 = charBuilder{}
	var b94 = charBuilder{}
	var b95 = charBuilder{}
	var b96 = charBuilder{}
	var b97 = charBuilder{}
	b98.items = []builder{&b93, &b94, &b95, &b96, &b97}
	b99.options = []builder{&b92, &b98}
	var b514 = sequenceBuilder{id: 514, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b506 = sequenceBuilder{id: 506, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b505 = sequenceBuilder{id: 505, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b498 = charBuilder{}
	var b499 = charBuilder{}
	var b500 = charBuilder{}
	var b501 = charBuilder{}
	var b502 = charBuilder{}
	var b503 = charBuilder{}
	var b504 = charBuilder{}
	b505.items = []builder{&b498, &b499, &b500, &b501, &b502, &b503, &b504}
	b506.items = []builder{&b505, &b15}
	var b513 = sequenceBuilder{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b512 = sequenceBuilder{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b512.items = []builder{&b825, &b14}
	b513.items = []builder{&b825, &b14, &b512}
	b514.items = []builder{&b506, &b513, &b825, &b272}
	var b104 = sequenceBuilder{id: 104, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b101 = sequenceBuilder{id: 101, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b100 = charBuilder{}
	b101.items = []builder{&b100}
	var b103 = sequenceBuilder{id: 103, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b102 = charBuilder{}
	b103.items = []builder{&b102}
	b104.items = []builder{&b101, &b103}
	var b125 = sequenceBuilder{id: 125, commit: 64, name: "list", ranges: [][]int{{1, 1}}}
	var b124 = sequenceBuilder{id: 124, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b121 = sequenceBuilder{id: 121, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b120 = charBuilder{}
	b121.items = []builder{&b120}
	var b114 = sequenceBuilder{id: 114, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b112 = choiceBuilder{id: 112, commit: 2}
	var b111 = sequenceBuilder{id: 111, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b110 = charBuilder{}
	b111.items = []builder{&b110}
	b112.options = []builder{&b14, &b111}
	var b113 = sequenceBuilder{id: 113, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b113.items = []builder{&b825, &b112}
	b114.items = []builder{&b112, &b113}
	var b119 = sequenceBuilder{id: 119, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b115 = choiceBuilder{id: 115, commit: 66}
	var b109 = sequenceBuilder{id: 109, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b108 = sequenceBuilder{id: 108, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	b108.items = []builder{&b105, &b106, &b107}
	b109.items = []builder{&b272, &b825, &b108}
	b115.options = []builder{&b401, &b109}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b116.items = []builder{&b114, &b825, &b115}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b117.items = []builder{&b825, &b116}
	b118.items = []builder{&b825, &b116, &b117}
	b119.items = []builder{&b115, &b118}
	var b123 = sequenceBuilder{id: 123, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b122 = charBuilder{}
	b123.items = []builder{&b122}
	b124.items = []builder{&b121, &b825, &b114, &b825, &b119, &b825, &b114, &b825, &b123}
	b125.items = []builder{&b124}
	var b130 = sequenceBuilder{id: 130, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b127 = sequenceBuilder{id: 127, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b126 = charBuilder{}
	b127.items = []builder{&b126}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b128.items = []builder{&b825, &b14}
	b129.items = []builder{&b825, &b14, &b128}
	b130.items = []builder{&b127, &b129, &b825, &b124}
	var b159 = sequenceBuilder{id: 159, commit: 64, name: "struct", ranges: [][]int{{1, 1}}}
	var b158 = sequenceBuilder{id: 158, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b155 = sequenceBuilder{id: 155, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b154 = charBuilder{}
	b155.items = []builder{&b154}
	var b153 = sequenceBuilder{id: 153, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b148 = choiceBuilder{id: 148, commit: 2}
	var b147 = sequenceBuilder{id: 147, commit: 64, name: "entry", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b140 = choiceBuilder{id: 140, commit: 2}
	var b139 = sequenceBuilder{id: 139, commit: 64, name: "expression-key", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b132 = sequenceBuilder{id: 132, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b131 = charBuilder{}
	b132.items = []builder{&b131}
	var b136 = sequenceBuilder{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b135 = sequenceBuilder{id: 135, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b135.items = []builder{&b825, &b14}
	b136.items = []builder{&b825, &b14, &b135}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b137.items = []builder{&b825, &b14}
	b138.items = []builder{&b825, &b14, &b137}
	var b134 = sequenceBuilder{id: 134, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b133 = charBuilder{}
	b134.items = []builder{&b133}
	b139.items = []builder{&b132, &b136, &b825, &b401, &b138, &b825, &b134}
	b140.options = []builder{&b104, &b87, &b139}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b143.items = []builder{&b825, &b14}
	b144.items = []builder{&b825, &b14, &b143}
	var b142 = sequenceBuilder{id: 142, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b141 = charBuilder{}
	b142.items = []builder{&b141}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b145.items = []builder{&b825, &b14}
	b146.items = []builder{&b825, &b14, &b145}
	b147.items = []builder{&b140, &b144, &b825, &b142, &b146, &b825, &b401}
	b148.options = []builder{&b147, &b109}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b149 = choiceBuilder{id: 149, commit: 2}
	b149.options = []builder{&b147, &b109}
	b150.items = []builder{&b114, &b825, &b149}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b151.items = []builder{&b825, &b150}
	b152.items = []builder{&b825, &b150, &b151}
	b153.items = []builder{&b148, &b152}
	var b157 = sequenceBuilder{id: 157, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b156 = charBuilder{}
	b157.items = []builder{&b156}
	b158.items = []builder{&b155, &b825, &b114, &b825, &b153, &b825, &b114, &b825, &b157}
	b159.items = []builder{&b158}
	var b164 = sequenceBuilder{id: 164, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b161 = sequenceBuilder{id: 161, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b160 = charBuilder{}
	b161.items = []builder{&b160}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b162.items = []builder{&b825, &b14}
	b163.items = []builder{&b825, &b14, &b162}
	b164.items = []builder{&b161, &b163, &b825, &b158}
	var b207 = sequenceBuilder{id: 207, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b204 = sequenceBuilder{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b202 = charBuilder{}
	var b203 = charBuilder{}
	b204.items = []builder{&b202, &b203}
	var b206 = sequenceBuilder{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b205.items = []builder{&b825, &b14}
	b206.items = []builder{&b825, &b14, &b205}
	var b201 = sequenceBuilder{id: 201, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b193 = sequenceBuilder{id: 193, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b192 = charBuilder{}
	b193.items = []builder{&b192}
	var b195 = choiceBuilder{id: 195, commit: 2}
	var b168 = sequenceBuilder{id: 168, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b167 = sequenceBuilder{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b165.items = []builder{&b114, &b825, &b104}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b166.items = []builder{&b825, &b165}
	b167.items = []builder{&b825, &b165, &b166}
	b168.items = []builder{&b104, &b167}
	var b194 = sequenceBuilder{id: 194, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b175 = sequenceBuilder{id: 175, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b172 = sequenceBuilder{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b169 = charBuilder{}
	var b170 = charBuilder{}
	var b171 = charBuilder{}
	b172.items = []builder{&b169, &b170, &b171}
	var b174 = sequenceBuilder{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b173 = sequenceBuilder{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b173.items = []builder{&b825, &b14}
	b174.items = []builder{&b825, &b14, &b173}
	b175.items = []builder{&b172, &b174, &b825, &b104}
	b194.items = []builder{&b168, &b825, &b114, &b825, &b175}
	b195.options = []builder{&b168, &b194, &b175}
	var b197 = sequenceBuilder{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b196 = charBuilder{}
	b197.items = []builder{&b196}
	var b200 = sequenceBuilder{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b199.items = []builder{&b825, &b14}
	b200.items = []builder{&b825, &b14, &b199}
	var b198 = choiceBuilder{id: 198, commit: 2}
	var b783 = choiceBuilder{id: 783, commit: 66}
	var b511 = sequenceBuilder{id: 511, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b497 = sequenceBuilder{id: 497, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b496 = sequenceBuilder{id: 496, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b492 = charBuilder{}
	var b493 = charBuilder{}
	var b494 = charBuilder{}
	var b495 = charBuilder{}
	b496.items = []builder{&b492, &b493, &b494, &b495}
	b497.items = []builder{&b496, &b15}
	var b508 = sequenceBuilder{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b507 = sequenceBuilder{id: 507, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b507.items = []builder{&b825, &b14}
	b508.items = []builder{&b825, &b14, &b507}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b509.items = []builder{&b825, &b14}
	b510.items = []builder{&b825, &b14, &b509}
	b511.items = []builder{&b497, &b508, &b825, &b272, &b510, &b825, &b272}
	var b564 = sequenceBuilder{id: 564, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b554 = sequenceBuilder{id: 554, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b553 = sequenceBuilder{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b551 = charBuilder{}
	var b552 = charBuilder{}
	b553.items = []builder{&b551, &b552}
	b554.items = []builder{&b553, &b15}
	var b563 = sequenceBuilder{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b562 = sequenceBuilder{id: 562, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b562.items = []builder{&b825, &b14}
	b563.items = []builder{&b825, &b14, &b562}
	var b262 = sequenceBuilder{id: 262, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b259 = sequenceBuilder{id: 259, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b258 = charBuilder{}
	b259.items = []builder{&b258}
	var b261 = sequenceBuilder{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b260 = charBuilder{}
	b261.items = []builder{&b260}
	b262.items = []builder{&b272, &b825, &b259, &b825, &b114, &b825, &b119, &b825, &b114, &b825, &b261}
	b564.items = []builder{&b554, &b563, &b825, &b262}
	var b573 = sequenceBuilder{id: 573, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b570 = sequenceBuilder{id: 570, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b565 = charBuilder{}
	var b566 = charBuilder{}
	var b567 = charBuilder{}
	var b568 = charBuilder{}
	var b569 = charBuilder{}
	b570.items = []builder{&b565, &b566, &b567, &b568, &b569}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b571 = sequenceBuilder{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b571.items = []builder{&b825, &b14}
	b572.items = []builder{&b825, &b14, &b571}
	b573.items = []builder{&b570, &b572, &b825, &b262}
	var b637 = choiceBuilder{id: 637, commit: 64, name: "assignment"}
	var b621 = sequenceBuilder{id: 621, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b606 = sequenceBuilder{id: 606, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b605 = sequenceBuilder{id: 605, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b602 = charBuilder{}
	var b603 = charBuilder{}
	var b604 = charBuilder{}
	b605.items = []builder{&b602, &b603, &b604}
	b606.items = []builder{&b605, &b15}
	var b620 = sequenceBuilder{id: 620, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b619 = sequenceBuilder{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b619.items = []builder{&b825, &b14}
	b620.items = []builder{&b825, &b14, &b619}
	var b614 = sequenceBuilder{id: 614, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b611 = sequenceBuilder{id: 611, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b610 = sequenceBuilder{id: 610, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b609 = sequenceBuilder{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b609.items = []builder{&b825, &b14}
	b610.items = []builder{&b14, &b609}
	var b608 = sequenceBuilder{id: 608, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b607 = charBuilder{}
	b608.items = []builder{&b607}
	b611.items = []builder{&b610, &b825, &b608}
	var b613 = sequenceBuilder{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b612 = sequenceBuilder{id: 612, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b612.items = []builder{&b825, &b14}
	b613.items = []builder{&b825, &b14, &b612}
	b614.items = []builder{&b272, &b825, &b611, &b613, &b825, &b401}
	b621.items = []builder{&b606, &b620, &b825, &b614}
	var b628 = sequenceBuilder{id: 628, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b625 = sequenceBuilder{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b624.items = []builder{&b825, &b14}
	b625.items = []builder{&b825, &b14, &b624}
	var b623 = sequenceBuilder{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b622 = charBuilder{}
	b623.items = []builder{&b622}
	var b627 = sequenceBuilder{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b626 = sequenceBuilder{id: 626, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b626.items = []builder{&b825, &b14}
	b627.items = []builder{&b825, &b14, &b626}
	b628.items = []builder{&b272, &b625, &b825, &b623, &b627, &b825, &b401}
	var b636 = sequenceBuilder{id: 636, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b634.items = []builder{&b825, &b14}
	b635.items = []builder{&b825, &b14, &b634}
	var b630 = sequenceBuilder{id: 630, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b629 = charBuilder{}
	b630.items = []builder{&b629}
	var b631 = sequenceBuilder{id: 631, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b618 = sequenceBuilder{id: 618, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b615.items = []builder{&b114, &b825, &b614}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b616.items = []builder{&b825, &b615}
	b617.items = []builder{&b825, &b615, &b616}
	b618.items = []builder{&b614, &b617}
	b631.items = []builder{&b114, &b825, &b618}
	var b633 = sequenceBuilder{id: 633, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b632 = charBuilder{}
	b633.items = []builder{&b632}
	b636.items = []builder{&b606, &b635, &b825, &b630, &b825, &b631, &b825, &b114, &b825, &b633}
	b637.options = []builder{&b621, &b628, &b636}
	var b792 = sequenceBuilder{id: 792, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b785 = sequenceBuilder{id: 785, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b784 = charBuilder{}
	b785.items = []builder{&b784}
	var b789 = sequenceBuilder{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b788 = sequenceBuilder{id: 788, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b788.items = []builder{&b825, &b14}
	b789.items = []builder{&b825, &b14, &b788}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b790.items = []builder{&b825, &b14}
	b791.items = []builder{&b825, &b14, &b790}
	var b787 = sequenceBuilder{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b786 = charBuilder{}
	b787.items = []builder{&b786}
	b792.items = []builder{&b785, &b789, &b825, &b783, &b791, &b825, &b787}
	b783.options = []builder{&b511, &b564, &b573, &b637, &b792, &b401}
	var b191 = sequenceBuilder{id: 191, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b188 = sequenceBuilder{id: 188, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b187 = charBuilder{}
	b188.items = []builder{&b187}
	var b190 = sequenceBuilder{id: 190, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b189 = charBuilder{}
	b190.items = []builder{&b189}
	b191.items = []builder{&b188, &b825, &b807, &b825, &b811, &b825, &b807, &b825, &b190}
	b198.options = []builder{&b783, &b191}
	b201.items = []builder{&b193, &b825, &b114, &b825, &b195, &b825, &b114, &b825, &b197, &b200, &b825, &b198}
	b207.items = []builder{&b204, &b206, &b825, &b201}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b210 = sequenceBuilder{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b208 = charBuilder{}
	var b209 = charBuilder{}
	b210.items = []builder{&b208, &b209}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b213.items = []builder{&b825, &b14}
	b214.items = []builder{&b825, &b14, &b213}
	var b212 = sequenceBuilder{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b211 = charBuilder{}
	b212.items = []builder{&b211}
	var b216 = sequenceBuilder{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b215.items = []builder{&b825, &b14}
	b216.items = []builder{&b825, &b14, &b215}
	b217.items = []builder{&b210, &b214, &b825, &b212, &b216, &b825, &b201}
	var b257 = sequenceBuilder{id: 257, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b256 = sequenceBuilder{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b255 = sequenceBuilder{id: 255, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b255.items = []builder{&b825, &b14}
	b256.items = []builder{&b825, &b14, &b255}
	var b254 = sequenceBuilder{id: 254, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b250 = choiceBuilder{id: 250, commit: 66}
	var b231 = sequenceBuilder{id: 231, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b228 = sequenceBuilder{id: 228, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b227 = charBuilder{}
	b228.items = []builder{&b227}
	var b230 = sequenceBuilder{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b229 = sequenceBuilder{id: 229, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b229.items = []builder{&b825, &b14}
	b230.items = []builder{&b825, &b14, &b229}
	b231.items = []builder{&b228, &b230, &b825, &b104}
	var b240 = sequenceBuilder{id: 240, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b233 = sequenceBuilder{id: 233, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b232 = charBuilder{}
	b233.items = []builder{&b232}
	var b237 = sequenceBuilder{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b236 = sequenceBuilder{id: 236, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b236.items = []builder{&b825, &b14}
	b237.items = []builder{&b825, &b14, &b236}
	var b239 = sequenceBuilder{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b238 = sequenceBuilder{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b238.items = []builder{&b825, &b14}
	b239.items = []builder{&b825, &b14, &b238}
	var b235 = sequenceBuilder{id: 235, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b234 = charBuilder{}
	b235.items = []builder{&b234}
	b240.items = []builder{&b233, &b237, &b825, &b401, &b239, &b825, &b235}
	var b249 = sequenceBuilder{id: 249, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b242 = sequenceBuilder{id: 242, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b241 = charBuilder{}
	b242.items = []builder{&b241}
	var b246 = sequenceBuilder{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b245 = sequenceBuilder{id: 245, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b245.items = []builder{&b825, &b14}
	b246.items = []builder{&b825, &b14, &b245}
	var b226 = sequenceBuilder{id: 226, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b401}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b222.items = []builder{&b825, &b14}
	b223.items = []builder{&b825, &b14, &b222}
	var b221 = sequenceBuilder{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b220 = charBuilder{}
	b221.items = []builder{&b220}
	var b225 = sequenceBuilder{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b224.items = []builder{&b825, &b14}
	b225.items = []builder{&b825, &b14, &b224}
	var b219 = sequenceBuilder{id: 219, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b219.items = []builder{&b401}
	b226.items = []builder{&b218, &b223, &b825, &b221, &b225, &b825, &b219}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b247.items = []builder{&b825, &b14}
	b248.items = []builder{&b825, &b14, &b247}
	var b244 = sequenceBuilder{id: 244, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b243 = charBuilder{}
	b244.items = []builder{&b243}
	b249.items = []builder{&b242, &b246, &b825, &b226, &b248, &b825, &b244}
	b250.options = []builder{&b231, &b240, &b249}
	var b253 = sequenceBuilder{id: 253, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b252 = sequenceBuilder{id: 252, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b251 = sequenceBuilder{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b251.items = []builder{&b825, &b14}
	b252.items = []builder{&b14, &b251}
	b253.items = []builder{&b252, &b825, &b250}
	b254.items = []builder{&b250, &b825, &b253}
	b257.items = []builder{&b272, &b256, &b825, &b254}
	var b271 = sequenceBuilder{id: 271, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b264 = sequenceBuilder{id: 264, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b263 = charBuilder{}
	b264.items = []builder{&b263}
	var b268 = sequenceBuilder{id: 268, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b267 = sequenceBuilder{id: 267, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b267.items = []builder{&b825, &b14}
	b268.items = []builder{&b825, &b14, &b267}
	var b270 = sequenceBuilder{id: 270, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b269 = sequenceBuilder{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b269.items = []builder{&b825, &b14}
	b270.items = []builder{&b825, &b14, &b269}
	var b266 = sequenceBuilder{id: 266, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b265 = charBuilder{}
	b266.items = []builder{&b265}
	b271.items = []builder{&b264, &b268, &b825, &b401, &b270, &b825, &b266}
	b272.options = []builder{&b61, &b74, &b87, &b99, &b514, &b104, &b125, &b130, &b159, &b164, &b207, &b217, &b257, &b262, &b271}
	var b332 = sequenceBuilder{id: 332, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b331 = choiceBuilder{id: 331, commit: 66}
	var b291 = sequenceBuilder{id: 291, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b290 = charBuilder{}
	b291.items = []builder{&b290}
	var b293 = sequenceBuilder{id: 293, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b292 = charBuilder{}
	b293.items = []builder{&b292}
	var b274 = sequenceBuilder{id: 274, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b273 = charBuilder{}
	b274.items = []builder{&b273}
	var b305 = sequenceBuilder{id: 305, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b304 = charBuilder{}
	b305.items = []builder{&b304}
	b331.options = []builder{&b291, &b293, &b274, &b305}
	b332.items = []builder{&b331, &b825, &b272}
	var b379 = choiceBuilder{id: 379, commit: 66}
	var b350 = sequenceBuilder{id: 350, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b338 = choiceBuilder{id: 338, commit: 66}
	b338.options = []builder{&b272, &b332}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b345 = sequenceBuilder{id: 345, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b344 = sequenceBuilder{id: 344, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b344.items = []builder{&b825, &b14}
	b345.items = []builder{&b14, &b344}
	var b333 = choiceBuilder{id: 333, commit: 66}
	var b276 = sequenceBuilder{id: 276, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b275 = charBuilder{}
	b276.items = []builder{&b275}
	var b283 = sequenceBuilder{id: 283, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b281 = charBuilder{}
	var b282 = charBuilder{}
	b283.items = []builder{&b281, &b282}
	var b286 = sequenceBuilder{id: 286, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b284 = charBuilder{}
	var b285 = charBuilder{}
	b286.items = []builder{&b284, &b285}
	var b289 = sequenceBuilder{id: 289, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b287 = charBuilder{}
	var b288 = charBuilder{}
	b289.items = []builder{&b287, &b288}
	var b295 = sequenceBuilder{id: 295, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b294 = charBuilder{}
	b295.items = []builder{&b294}
	var b297 = sequenceBuilder{id: 297, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b296 = charBuilder{}
	b297.items = []builder{&b296}
	var b299 = sequenceBuilder{id: 299, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b298 = charBuilder{}
	b299.items = []builder{&b298}
	b333.options = []builder{&b276, &b283, &b286, &b289, &b295, &b297, &b299}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b346.items = []builder{&b825, &b14}
	b347.items = []builder{&b825, &b14, &b346}
	b348.items = []builder{&b345, &b825, &b333, &b347, &b825, &b338}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b349.items = []builder{&b825, &b348}
	b350.items = []builder{&b338, &b825, &b348, &b349}
	var b357 = sequenceBuilder{id: 357, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b339 = choiceBuilder{id: 339, commit: 66}
	b339.options = []builder{&b338, &b350}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b351 = sequenceBuilder{id: 351, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b351.items = []builder{&b825, &b14}
	b352.items = []builder{&b14, &b351}
	var b334 = choiceBuilder{id: 334, commit: 66}
	var b278 = sequenceBuilder{id: 278, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b277 = charBuilder{}
	b278.items = []builder{&b277}
	var b280 = sequenceBuilder{id: 280, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b279 = charBuilder{}
	b280.items = []builder{&b279}
	var b301 = sequenceBuilder{id: 301, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b300 = charBuilder{}
	b301.items = []builder{&b300}
	var b303 = sequenceBuilder{id: 303, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b302 = charBuilder{}
	b303.items = []builder{&b302}
	b334.options = []builder{&b278, &b280, &b301, &b303}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b353.items = []builder{&b825, &b14}
	b354.items = []builder{&b825, &b14, &b353}
	b355.items = []builder{&b352, &b825, &b334, &b354, &b825, &b339}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b356.items = []builder{&b825, &b355}
	b357.items = []builder{&b339, &b825, &b355, &b356}
	var b364 = sequenceBuilder{id: 364, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b340 = choiceBuilder{id: 340, commit: 66}
	b340.options = []builder{&b339, &b357}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b358.items = []builder{&b825, &b14}
	b359.items = []builder{&b14, &b358}
	var b335 = choiceBuilder{id: 335, commit: 66}
	var b308 = sequenceBuilder{id: 308, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b306 = charBuilder{}
	var b307 = charBuilder{}
	b308.items = []builder{&b306, &b307}
	var b311 = sequenceBuilder{id: 311, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b309 = charBuilder{}
	var b310 = charBuilder{}
	b311.items = []builder{&b309, &b310}
	var b313 = sequenceBuilder{id: 313, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b312 = charBuilder{}
	b313.items = []builder{&b312}
	var b316 = sequenceBuilder{id: 316, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b314 = charBuilder{}
	var b315 = charBuilder{}
	b316.items = []builder{&b314, &b315}
	var b318 = sequenceBuilder{id: 318, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b317 = charBuilder{}
	b318.items = []builder{&b317}
	var b321 = sequenceBuilder{id: 321, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b319 = charBuilder{}
	var b320 = charBuilder{}
	b321.items = []builder{&b319, &b320}
	b335.options = []builder{&b308, &b311, &b313, &b316, &b318, &b321}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b360.items = []builder{&b825, &b14}
	b361.items = []builder{&b825, &b14, &b360}
	b362.items = []builder{&b359, &b825, &b335, &b361, &b825, &b340}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b363.items = []builder{&b825, &b362}
	b364.items = []builder{&b340, &b825, &b362, &b363}
	var b371 = sequenceBuilder{id: 371, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b341 = choiceBuilder{id: 341, commit: 66}
	b341.options = []builder{&b340, &b364}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b365 = sequenceBuilder{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b365.items = []builder{&b825, &b14}
	b366.items = []builder{&b14, &b365}
	var b336 = sequenceBuilder{id: 336, commit: 66, ranges: [][]int{{1, 1}}}
	var b324 = sequenceBuilder{id: 324, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b322 = charBuilder{}
	var b323 = charBuilder{}
	b324.items = []builder{&b322, &b323}
	b336.items = []builder{&b324}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b367.items = []builder{&b825, &b14}
	b368.items = []builder{&b825, &b14, &b367}
	b369.items = []builder{&b366, &b825, &b336, &b368, &b825, &b341}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b370.items = []builder{&b825, &b369}
	b371.items = []builder{&b341, &b825, &b369, &b370}
	var b378 = sequenceBuilder{id: 378, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b342 = choiceBuilder{id: 342, commit: 66}
	b342.options = []builder{&b341, &b371}
	var b376 = sequenceBuilder{id: 376, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b372.items = []builder{&b825, &b14}
	b373.items = []builder{&b14, &b372}
	var b337 = sequenceBuilder{id: 337, commit: 66, ranges: [][]int{{1, 1}}}
	var b327 = sequenceBuilder{id: 327, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b325 = charBuilder{}
	var b326 = charBuilder{}
	b327.items = []builder{&b325, &b326}
	b337.items = []builder{&b327}
	var b375 = sequenceBuilder{id: 375, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b374 = sequenceBuilder{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b374.items = []builder{&b825, &b14}
	b375.items = []builder{&b825, &b14, &b374}
	b376.items = []builder{&b373, &b825, &b337, &b375, &b825, &b342}
	var b377 = sequenceBuilder{id: 377, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b377.items = []builder{&b825, &b376}
	b378.items = []builder{&b342, &b825, &b376, &b377}
	b379.options = []builder{&b350, &b357, &b364, &b371, &b378}
	var b392 = sequenceBuilder{id: 392, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b384 = sequenceBuilder{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b384.items = []builder{&b825, &b14}
	b385.items = []builder{&b825, &b14, &b384}
	var b381 = sequenceBuilder{id: 381, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b380 = charBuilder{}
	b381.items = []builder{&b380}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b386.items = []builder{&b825, &b14}
	b387.items = []builder{&b825, &b14, &b386}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b388.items = []builder{&b825, &b14}
	b389.items = []builder{&b825, &b14, &b388}
	var b383 = sequenceBuilder{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b382 = charBuilder{}
	b383.items = []builder{&b382}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b390.items = []builder{&b825, &b14}
	b391.items = []builder{&b825, &b14, &b390}
	b392.items = []builder{&b401, &b385, &b825, &b381, &b387, &b825, &b401, &b389, &b825, &b383, &b391, &b825, &b401}
	var b400 = sequenceBuilder{id: 400, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b393 = choiceBuilder{id: 393, commit: 66}
	b393.options = []builder{&b272, &b332, &b379, &b392}
	var b398 = sequenceBuilder{id: 398, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b395 = sequenceBuilder{id: 395, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b394 = sequenceBuilder{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b394.items = []builder{&b825, &b14}
	b395.items = []builder{&b14, &b394}
	var b330 = sequenceBuilder{id: 330, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b328 = charBuilder{}
	var b329 = charBuilder{}
	b330.items = []builder{&b328, &b329}
	var b397 = sequenceBuilder{id: 397, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b396 = sequenceBuilder{id: 396, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b396.items = []builder{&b825, &b14}
	b397.items = []builder{&b825, &b14, &b396}
	b398.items = []builder{&b395, &b825, &b330, &b397, &b825, &b393}
	var b399 = sequenceBuilder{id: 399, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b399.items = []builder{&b825, &b398}
	b400.items = []builder{&b393, &b825, &b398, &b399}
	b401.options = []builder{&b272, &b332, &b379, &b392, &b400}
	b186.items = []builder{&b183, &b185, &b825, &b401}
	var b432 = sequenceBuilder{id: 432, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b405 = sequenceBuilder{id: 405, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b404 = sequenceBuilder{id: 404, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b402 = charBuilder{}
	var b403 = charBuilder{}
	b404.items = []builder{&b402, &b403}
	b405.items = []builder{&b404, &b15}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b426.items = []builder{&b825, &b14}
	b427.items = []builder{&b825, &b14, &b426}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b428.items = []builder{&b825, &b14}
	b429.items = []builder{&b825, &b14, &b428}
	var b431 = sequenceBuilder{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b420 = sequenceBuilder{id: 420, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b413 = sequenceBuilder{id: 413, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b412 = sequenceBuilder{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b412.items = []builder{&b825, &b14}
	b413.items = []builder{&b14, &b412}
	var b411 = sequenceBuilder{id: 411, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b410 = sequenceBuilder{id: 410, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b406 = charBuilder{}
	var b407 = charBuilder{}
	var b408 = charBuilder{}
	var b409 = charBuilder{}
	b410.items = []builder{&b406, &b407, &b408, &b409}
	b411.items = []builder{&b410, &b15}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b414.items = []builder{&b825, &b14}
	b415.items = []builder{&b825, &b14, &b414}
	var b417 = sequenceBuilder{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b416 = sequenceBuilder{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b416.items = []builder{&b825, &b14}
	b417.items = []builder{&b825, &b14, &b416}
	var b419 = sequenceBuilder{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b418 = sequenceBuilder{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b418.items = []builder{&b825, &b14}
	b419.items = []builder{&b825, &b14, &b418}
	b420.items = []builder{&b413, &b825, &b411, &b415, &b825, &b405, &b417, &b825, &b401, &b419, &b825, &b191}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b430.items = []builder{&b825, &b420}
	b431.items = []builder{&b825, &b420, &b430}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b421 = sequenceBuilder{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b421.items = []builder{&b825, &b14}
	b422.items = []builder{&b14, &b421}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b423.items = []builder{&b825, &b14}
	b424.items = []builder{&b825, &b14, &b423}
	b425.items = []builder{&b422, &b825, &b411, &b424, &b825, &b191}
	b432.items = []builder{&b405, &b427, &b825, &b401, &b429, &b825, &b191, &b431, &b825, &b425}
	var b491 = sequenceBuilder{id: 491, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b446 = sequenceBuilder{id: 446, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b445 = sequenceBuilder{id: 445, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b439 = charBuilder{}
	var b440 = charBuilder{}
	var b441 = charBuilder{}
	var b442 = charBuilder{}
	var b443 = charBuilder{}
	var b444 = charBuilder{}
	b445.items = []builder{&b439, &b440, &b441, &b442, &b443, &b444}
	b446.items = []builder{&b445, &b15}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b487.items = []builder{&b825, &b14}
	b488.items = []builder{&b825, &b14, &b487}
	var b490 = sequenceBuilder{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b489.items = []builder{&b825, &b14}
	b490.items = []builder{&b825, &b14, &b489}
	var b478 = sequenceBuilder{id: 478, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b477 = charBuilder{}
	b478.items = []builder{&b477}
	var b484 = sequenceBuilder{id: 484, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b479 = choiceBuilder{id: 479, commit: 2}
	var b476 = sequenceBuilder{id: 476, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b471 = sequenceBuilder{id: 471, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b438 = sequenceBuilder{id: 438, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b437 = sequenceBuilder{id: 437, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b433 = charBuilder{}
	var b434 = charBuilder{}
	var b435 = charBuilder{}
	var b436 = charBuilder{}
	b437.items = []builder{&b433, &b434, &b435, &b436}
	b438.items = []builder{&b437, &b15}
	var b468 = sequenceBuilder{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b467 = sequenceBuilder{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b467.items = []builder{&b825, &b14}
	b468.items = []builder{&b825, &b14, &b467}
	var b470 = sequenceBuilder{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b469 = sequenceBuilder{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b469.items = []builder{&b825, &b14}
	b470.items = []builder{&b825, &b14, &b469}
	var b466 = sequenceBuilder{id: 466, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b465 = charBuilder{}
	b466.items = []builder{&b465}
	b471.items = []builder{&b438, &b468, &b825, &b401, &b470, &b825, &b466}
	var b475 = sequenceBuilder{id: 475, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b473 = sequenceBuilder{id: 473, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b472 = charBuilder{}
	b473.items = []builder{&b472}
	var b474 = sequenceBuilder{id: 474, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b474.items = []builder{&b825, &b473}
	b475.items = []builder{&b825, &b473, &b474}
	b476.items = []builder{&b471, &b475, &b825, &b793}
	var b464 = sequenceBuilder{id: 464, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b459 = sequenceBuilder{id: 459, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b454 = sequenceBuilder{id: 454, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b447 = charBuilder{}
	var b448 = charBuilder{}
	var b449 = charBuilder{}
	var b450 = charBuilder{}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	b454.items = []builder{&b447, &b448, &b449, &b450, &b451, &b452, &b453}
	var b458 = sequenceBuilder{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b457 = sequenceBuilder{id: 457, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b457.items = []builder{&b825, &b14}
	b458.items = []builder{&b825, &b14, &b457}
	var b456 = sequenceBuilder{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b455 = charBuilder{}
	b456.items = []builder{&b455}
	b459.items = []builder{&b454, &b458, &b825, &b456}
	var b463 = sequenceBuilder{id: 463, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b461 = sequenceBuilder{id: 461, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b460 = charBuilder{}
	b461.items = []builder{&b460}
	var b462 = sequenceBuilder{id: 462, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b462.items = []builder{&b825, &b461}
	b463.items = []builder{&b825, &b461, &b462}
	b464.items = []builder{&b459, &b463, &b825, &b793}
	b479.options = []builder{&b476, &b464}
	var b483 = sequenceBuilder{id: 483, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b481 = sequenceBuilder{id: 481, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b480 = choiceBuilder{id: 480, commit: 2}
	b480.options = []builder{&b476, &b464, &b793}
	b481.items = []builder{&b807, &b825, &b480}
	var b482 = sequenceBuilder{id: 482, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b482.items = []builder{&b825, &b481}
	b483.items = []builder{&b825, &b481, &b482}
	b484.items = []builder{&b479, &b483}
	var b486 = sequenceBuilder{id: 486, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b485 = charBuilder{}
	b486.items = []builder{&b485}
	b491.items = []builder{&b446, &b488, &b825, &b401, &b490, &b825, &b478, &b825, &b807, &b825, &b484, &b825, &b807, &b825, &b486}
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
	b548.items = []builder{&b825, &b14}
	b549.items = []builder{&b825, &b14, &b548}
	var b539 = sequenceBuilder{id: 539, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b538 = charBuilder{}
	b539.items = []builder{&b538}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b540 = choiceBuilder{id: 540, commit: 2}
	var b530 = sequenceBuilder{id: 530, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b525 = sequenceBuilder{id: 525, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b522 = sequenceBuilder{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b521 = sequenceBuilder{id: 521, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b521.items = []builder{&b825, &b14}
	b522.items = []builder{&b825, &b14, &b521}
	var b518 = choiceBuilder{id: 518, commit: 66}
	var b517 = sequenceBuilder{id: 517, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b516 = sequenceBuilder{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b515 = sequenceBuilder{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b515.items = []builder{&b825, &b14}
	b516.items = []builder{&b825, &b14, &b515}
	b517.items = []builder{&b104, &b516, &b825, &b514}
	b518.options = []builder{&b511, &b514, &b517}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b523.items = []builder{&b825, &b14}
	b524.items = []builder{&b825, &b14, &b523}
	var b520 = sequenceBuilder{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b519 = charBuilder{}
	b520.items = []builder{&b519}
	b525.items = []builder{&b438, &b522, &b825, &b518, &b524, &b825, &b520}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b527 = sequenceBuilder{id: 527, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b526 = charBuilder{}
	b527.items = []builder{&b526}
	var b528 = sequenceBuilder{id: 528, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b528.items = []builder{&b825, &b527}
	b529.items = []builder{&b825, &b527, &b528}
	b530.items = []builder{&b525, &b529, &b825, &b793}
	b540.options = []builder{&b530, &b464}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b542 = sequenceBuilder{id: 542, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b541 = choiceBuilder{id: 541, commit: 2}
	b541.options = []builder{&b530, &b464, &b793}
	b542.items = []builder{&b807, &b825, &b541}
	var b543 = sequenceBuilder{id: 543, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b543.items = []builder{&b825, &b542}
	b544.items = []builder{&b825, &b542, &b543}
	b545.items = []builder{&b540, &b544}
	var b547 = sequenceBuilder{id: 547, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b546 = charBuilder{}
	b547.items = []builder{&b546}
	b550.items = []builder{&b537, &b549, &b825, &b539, &b825, &b807, &b825, &b545, &b825, &b807, &b825, &b547}
	var b601 = sequenceBuilder{id: 601, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b582 = sequenceBuilder{id: 582, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b581 = sequenceBuilder{id: 581, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b578 = charBuilder{}
	var b579 = charBuilder{}
	var b580 = charBuilder{}
	b581.items = []builder{&b578, &b579, &b580}
	b582.items = []builder{&b581, &b15}
	var b600 = choiceBuilder{id: 600, commit: 2}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b592 = sequenceBuilder{id: 592, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b591 = sequenceBuilder{id: 591, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b591.items = []builder{&b825, &b14}
	b592.items = []builder{&b14, &b591}
	var b590 = choiceBuilder{id: 590, commit: 66}
	var b589 = choiceBuilder{id: 589, commit: 64, name: "range-over-expression"}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b584.items = []builder{&b825, &b14}
	b585.items = []builder{&b825, &b14, &b584}
	var b577 = sequenceBuilder{id: 577, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b576 = sequenceBuilder{id: 576, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b574 = charBuilder{}
	var b575 = charBuilder{}
	b576.items = []builder{&b574, &b575}
	b577.items = []builder{&b576, &b15}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b586.items = []builder{&b825, &b14}
	b587.items = []builder{&b825, &b14, &b586}
	var b583 = choiceBuilder{id: 583, commit: 2}
	b583.options = []builder{&b401, &b226}
	b588.items = []builder{&b104, &b585, &b825, &b577, &b587, &b825, &b583}
	b589.options = []builder{&b588, &b226}
	b590.options = []builder{&b401, &b589}
	b593.items = []builder{&b592, &b825, &b590}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b594.items = []builder{&b825, &b14}
	b595.items = []builder{&b825, &b14, &b594}
	b596.items = []builder{&b593, &b595, &b825, &b191}
	var b599 = sequenceBuilder{id: 599, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b597.items = []builder{&b825, &b14}
	b598.items = []builder{&b14, &b597}
	b599.items = []builder{&b598, &b825, &b191}
	b600.options = []builder{&b596, &b599}
	b601.items = []builder{&b582, &b825, &b600}
	var b739 = choiceBuilder{id: 739, commit: 66}
	var b660 = sequenceBuilder{id: 660, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b642 = sequenceBuilder{id: 642, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b641 = sequenceBuilder{id: 641, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b638 = charBuilder{}
	var b639 = charBuilder{}
	var b640 = charBuilder{}
	b641.items = []builder{&b638, &b639, &b640}
	b642.items = []builder{&b641, &b15}
	var b659 = sequenceBuilder{id: 659, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b658 = sequenceBuilder{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b658.items = []builder{&b825, &b14}
	b659.items = []builder{&b825, &b14, &b658}
	var b657 = choiceBuilder{id: 657, commit: 2}
	var b651 = sequenceBuilder{id: 651, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b650 = sequenceBuilder{id: 650, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b647 = sequenceBuilder{id: 647, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b646 = sequenceBuilder{id: 646, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b645 = sequenceBuilder{id: 645, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b645.items = []builder{&b825, &b14}
	b646.items = []builder{&b14, &b645}
	var b644 = sequenceBuilder{id: 644, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b643 = charBuilder{}
	b644.items = []builder{&b643}
	b647.items = []builder{&b646, &b825, &b644}
	var b649 = sequenceBuilder{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b648 = sequenceBuilder{id: 648, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b648.items = []builder{&b825, &b14}
	b649.items = []builder{&b825, &b14, &b648}
	b650.items = []builder{&b104, &b825, &b647, &b649, &b825, &b401}
	b651.items = []builder{&b650}
	var b656 = sequenceBuilder{id: 656, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b653 = sequenceBuilder{id: 653, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b652 = charBuilder{}
	b653.items = []builder{&b652}
	var b655 = sequenceBuilder{id: 655, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b654 = sequenceBuilder{id: 654, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b654.items = []builder{&b825, &b14}
	b655.items = []builder{&b825, &b14, &b654}
	b656.items = []builder{&b653, &b655, &b825, &b650}
	b657.options = []builder{&b651, &b656}
	b660.items = []builder{&b642, &b659, &b825, &b657}
	var b677 = sequenceBuilder{id: 677, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b676 = sequenceBuilder{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b675 = sequenceBuilder{id: 675, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b675.items = []builder{&b825, &b14}
	b676.items = []builder{&b825, &b14, &b675}
	var b672 = sequenceBuilder{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b671 = charBuilder{}
	b672.items = []builder{&b671}
	var b666 = sequenceBuilder{id: 666, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b661 = choiceBuilder{id: 661, commit: 2}
	b661.options = []builder{&b651, &b656}
	var b665 = sequenceBuilder{id: 665, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b663 = sequenceBuilder{id: 663, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b662 = choiceBuilder{id: 662, commit: 2}
	b662.options = []builder{&b651, &b656}
	b663.items = []builder{&b114, &b825, &b662}
	var b664 = sequenceBuilder{id: 664, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b664.items = []builder{&b825, &b663}
	b665.items = []builder{&b825, &b663, &b664}
	b666.items = []builder{&b661, &b665}
	var b674 = sequenceBuilder{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b673 = charBuilder{}
	b674.items = []builder{&b673}
	b677.items = []builder{&b642, &b676, &b825, &b672, &b825, &b114, &b825, &b666, &b825, &b114, &b825, &b674}
	var b688 = sequenceBuilder{id: 688, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b685 = sequenceBuilder{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b684 = sequenceBuilder{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b684.items = []builder{&b825, &b14}
	b685.items = []builder{&b825, &b14, &b684}
	var b679 = sequenceBuilder{id: 679, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b678 = charBuilder{}
	b679.items = []builder{&b678}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b686.items = []builder{&b825, &b14}
	b687.items = []builder{&b825, &b14, &b686}
	var b681 = sequenceBuilder{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	b681.items = []builder{&b680}
	var b670 = sequenceBuilder{id: 670, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b669 = sequenceBuilder{id: 669, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b667.items = []builder{&b114, &b825, &b651}
	var b668 = sequenceBuilder{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b668.items = []builder{&b825, &b667}
	b669.items = []builder{&b825, &b667, &b668}
	b670.items = []builder{&b651, &b669}
	var b683 = sequenceBuilder{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b682 = charBuilder{}
	b683.items = []builder{&b682}
	b688.items = []builder{&b642, &b685, &b825, &b679, &b687, &b825, &b681, &b825, &b114, &b825, &b670, &b825, &b114, &b825, &b683}
	var b704 = sequenceBuilder{id: 704, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b700 = sequenceBuilder{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b698 = charBuilder{}
	var b699 = charBuilder{}
	b700.items = []builder{&b698, &b699}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b702 = sequenceBuilder{id: 702, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b702.items = []builder{&b825, &b14}
	b703.items = []builder{&b825, &b14, &b702}
	var b701 = choiceBuilder{id: 701, commit: 2}
	var b692 = sequenceBuilder{id: 692, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b691 = sequenceBuilder{id: 691, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b690 = sequenceBuilder{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b689 = sequenceBuilder{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b689.items = []builder{&b825, &b14}
	b690.items = []builder{&b825, &b14, &b689}
	b691.items = []builder{&b104, &b690, &b825, &b201}
	b692.items = []builder{&b691}
	var b697 = sequenceBuilder{id: 697, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b694 = sequenceBuilder{id: 694, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b693 = charBuilder{}
	b694.items = []builder{&b693}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b695 = sequenceBuilder{id: 695, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b695.items = []builder{&b825, &b14}
	b696.items = []builder{&b825, &b14, &b695}
	b697.items = []builder{&b694, &b696, &b825, &b691}
	b701.options = []builder{&b692, &b697}
	b704.items = []builder{&b700, &b703, &b825, &b701}
	var b724 = sequenceBuilder{id: 724, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b717 = sequenceBuilder{id: 717, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b715 = charBuilder{}
	var b716 = charBuilder{}
	b717.items = []builder{&b715, &b716}
	var b723 = sequenceBuilder{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b722 = sequenceBuilder{id: 722, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b722.items = []builder{&b825, &b14}
	b723.items = []builder{&b825, &b14, &b722}
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
	b711.items = []builder{&b114, &b825, &b710}
	var b712 = sequenceBuilder{id: 712, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b712.items = []builder{&b825, &b711}
	b713.items = []builder{&b825, &b711, &b712}
	b714.items = []builder{&b709, &b713}
	var b721 = sequenceBuilder{id: 721, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b720 = charBuilder{}
	b721.items = []builder{&b720}
	b724.items = []builder{&b717, &b723, &b825, &b719, &b825, &b114, &b825, &b714, &b825, &b114, &b825, &b721}
	var b738 = sequenceBuilder{id: 738, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b727 = sequenceBuilder{id: 727, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b725 = charBuilder{}
	var b726 = charBuilder{}
	b727.items = []builder{&b725, &b726}
	var b735 = sequenceBuilder{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b734 = sequenceBuilder{id: 734, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b734.items = []builder{&b825, &b14}
	b735.items = []builder{&b825, &b14, &b734}
	var b729 = sequenceBuilder{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b728 = charBuilder{}
	b729.items = []builder{&b728}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b736.items = []builder{&b825, &b14}
	b737.items = []builder{&b825, &b14, &b736}
	var b731 = sequenceBuilder{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b730 = charBuilder{}
	b731.items = []builder{&b730}
	var b708 = sequenceBuilder{id: 708, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b707 = sequenceBuilder{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b705 = sequenceBuilder{id: 705, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b705.items = []builder{&b114, &b825, &b692}
	var b706 = sequenceBuilder{id: 706, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b706.items = []builder{&b825, &b705}
	b707.items = []builder{&b825, &b705, &b706}
	b708.items = []builder{&b692, &b707}
	var b733 = sequenceBuilder{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b732 = charBuilder{}
	b733.items = []builder{&b732}
	b738.items = []builder{&b727, &b735, &b825, &b729, &b737, &b825, &b731, &b825, &b114, &b825, &b708, &b825, &b114, &b825, &b733}
	b739.options = []builder{&b660, &b677, &b688, &b704, &b724, &b738}
	var b771 = choiceBuilder{id: 771, commit: 64, name: "use"}
	var b763 = sequenceBuilder{id: 763, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b744 = sequenceBuilder{id: 744, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b743 = sequenceBuilder{id: 743, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b740 = charBuilder{}
	var b741 = charBuilder{}
	var b742 = charBuilder{}
	b743.items = []builder{&b740, &b741, &b742}
	b744.items = []builder{&b743, &b15}
	var b762 = sequenceBuilder{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b761 = sequenceBuilder{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b761.items = []builder{&b825, &b14}
	b762.items = []builder{&b825, &b14, &b761}
	var b756 = choiceBuilder{id: 756, commit: 64, name: "use-fact"}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b747 = choiceBuilder{id: 747, commit: 2}
	var b746 = sequenceBuilder{id: 746, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b745 = charBuilder{}
	b746.items = []builder{&b745}
	b747.options = []builder{&b104, &b746}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b751 = sequenceBuilder{id: 751, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b750 = sequenceBuilder{id: 750, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b750.items = []builder{&b825, &b14}
	b751.items = []builder{&b14, &b750}
	var b749 = sequenceBuilder{id: 749, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b748 = charBuilder{}
	b749.items = []builder{&b748}
	b752.items = []builder{&b751, &b825, &b749}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b753.items = []builder{&b825, &b14}
	b754.items = []builder{&b825, &b14, &b753}
	b755.items = []builder{&b747, &b825, &b752, &b754, &b825, &b87}
	b756.options = []builder{&b87, &b755}
	b763.items = []builder{&b744, &b762, &b825, &b756}
	var b770 = sequenceBuilder{id: 770, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b769 = sequenceBuilder{id: 769, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b768 = sequenceBuilder{id: 768, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b768.items = []builder{&b825, &b14}
	b769.items = []builder{&b825, &b14, &b768}
	var b765 = sequenceBuilder{id: 765, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b764 = charBuilder{}
	b765.items = []builder{&b764}
	var b760 = sequenceBuilder{id: 760, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b759 = sequenceBuilder{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b757 = sequenceBuilder{id: 757, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b757.items = []builder{&b114, &b825, &b756}
	var b758 = sequenceBuilder{id: 758, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b758.items = []builder{&b825, &b757}
	b759.items = []builder{&b825, &b757, &b758}
	b760.items = []builder{&b756, &b759}
	var b767 = sequenceBuilder{id: 767, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b766 = charBuilder{}
	b767.items = []builder{&b766}
	b770.items = []builder{&b744, &b769, &b825, &b765, &b825, &b114, &b825, &b760, &b825, &b114, &b825, &b767}
	b771.options = []builder{&b763, &b770}
	var b782 = sequenceBuilder{id: 782, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b779 = sequenceBuilder{id: 779, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b778 = sequenceBuilder{id: 778, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b772 = charBuilder{}
	var b773 = charBuilder{}
	var b774 = charBuilder{}
	var b775 = charBuilder{}
	var b776 = charBuilder{}
	var b777 = charBuilder{}
	b778.items = []builder{&b772, &b773, &b774, &b775, &b776, &b777}
	b779.items = []builder{&b778, &b15}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b780 = sequenceBuilder{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b780.items = []builder{&b825, &b14}
	b781.items = []builder{&b825, &b14, &b780}
	b782.items = []builder{&b779, &b781, &b825, &b739}
	var b802 = sequenceBuilder{id: 802, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b795 = sequenceBuilder{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b794 = charBuilder{}
	b795.items = []builder{&b794}
	var b799 = sequenceBuilder{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b798 = sequenceBuilder{id: 798, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b798.items = []builder{&b825, &b14}
	b799.items = []builder{&b825, &b14, &b798}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b800.items = []builder{&b825, &b14}
	b801.items = []builder{&b825, &b14, &b800}
	var b797 = sequenceBuilder{id: 797, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b796 = charBuilder{}
	b797.items = []builder{&b796}
	b802.items = []builder{&b795, &b799, &b825, &b793, &b801, &b825, &b797}
	b793.options = []builder{&b186, &b432, &b491, &b550, &b601, &b739, &b771, &b782, &b802, &b783}
	var b810 = sequenceBuilder{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b808 = sequenceBuilder{id: 808, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b808.items = []builder{&b807, &b825, &b793}
	var b809 = sequenceBuilder{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b809.items = []builder{&b825, &b808}
	b810.items = []builder{&b825, &b808, &b809}
	b811.items = []builder{&b793, &b810}
	b826.items = []builder{&b822, &b825, &b807, &b825, &b811, &b825, &b807}
	b827.items = []builder{&b825, &b826, &b825}

	return parseInput(r, &p827, &b827)
}
