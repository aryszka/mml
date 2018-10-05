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
	var p838 = sequenceParser{id: 838, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p836 = choiceParser{id: 836, commit: 2}
	var p834 = choiceParser{id: 834, commit: 70, name: "ws", generalizations: []int{836}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{834, 836}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{834, 836}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{834, 836}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{834, 836}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{834, 836}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{834, 836}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p834.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p835 = sequenceParser{id: 835, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{836}}
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
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{816, 111}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p39.items = []parser{&p14, &p836, &p38}
	var p40 = sequenceParser{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p40.items = []parser{&p836, &p39}
	p41.items = []parser{&p836, &p39, &p40}
	p42.items = []parser{&p38, &p41}
	p835.items = []parser{&p42}
	p836.options = []parser{&p834, &p835}
	var p837 = sequenceParser{id: 837, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p833 = sequenceParser{id: 833, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p830 = sequenceParser{id: 830, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p828 = charParser{id: 828, chars: []rune{35}}
	var p829 = charParser{id: 829, chars: []rune{33}}
	p830.items = []parser{&p828, &p829}
	var p827 = sequenceParser{id: 827, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p826 = sequenceParser{id: 826, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p824 = sequenceParser{id: 824, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p823 = charParser{id: 823, not: true, chars: []rune{10}}
	p824.items = []parser{&p823}
	var p825 = sequenceParser{id: 825, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p825.items = []parser{&p836, &p824}
	p826.items = []parser{&p824, &p825}
	p827.items = []parser{&p826}
	var p832 = sequenceParser{id: 832, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p831 = charParser{id: 831, chars: []rune{10}}
	p832.items = []parser{&p831}
	p833.items = []parser{&p830, &p836, &p827, &p836, &p832}
	var p818 = sequenceParser{id: 818, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p816 = choiceParser{id: 816, commit: 2}
	var p815 = sequenceParser{id: 815, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{816}}
	var p814 = charParser{id: 814, chars: []rune{59}}
	p815.items = []parser{&p814}
	p816.options = []parser{&p815, &p14}
	var p817 = sequenceParser{id: 817, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p817.items = []parser{&p836, &p816}
	p818.items = []parser{&p816, &p817}
	var p822 = sequenceParser{id: 822, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p804 = choiceParser{id: 804, commit: 66, name: "statement", generalizations: []int{478, 542}}
	var p185 = sequenceParser{id: 185, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{804, 478, 542}}
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
	p182.items = []parser{&p836, &p14}
	p183.items = []parser{&p14, &p182}
	var p395 = choiceParser{id: 395, commit: 66, name: "expression", generalizations: []int{114, 794, 197, 577, 570, 804}}
	var p266 = choiceParser{id: 266, commit: 66, name: "primary-expression", generalizations: []int{114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p60 = choiceParser{id: 60, commit: 64, name: "int", generalizations: []int{266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p51 = sequenceParser{id: 51, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p50 = sequenceParser{id: 50, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p49 = charParser{id: 49, ranges: [][]rune{{49, 57}}}
	p50.items = []parser{&p49}
	var p44 = sequenceParser{id: 44, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p43 = charParser{id: 43, ranges: [][]rune{{48, 57}}}
	p44.items = []parser{&p43}
	p51.items = []parser{&p50, &p44}
	var p54 = sequenceParser{id: 54, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{60, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p53 = sequenceParser{id: 53, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p52 = charParser{id: 52, chars: []rune{48}}
	p53.items = []parser{&p52}
	var p46 = sequenceParser{id: 46, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 55}}}
	p46.items = []parser{&p45}
	p54.items = []parser{&p53, &p46}
	var p59 = sequenceParser{id: 59, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{60, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
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
	var p73 = choiceParser{id: 73, commit: 72, name: "float", generalizations: []int{266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p68 = sequenceParser{id: 68, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{73, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
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
	var p71 = sequenceParser{id: 71, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{73, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p70 = sequenceParser{id: 70, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p69 = charParser{id: 69, chars: []rune{46}}
	p70.items = []parser{&p69}
	p71.items = []parser{&p70, &p44, &p65}
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{73, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	p72.items = []parser{&p44, &p65}
	p73.options = []parser{&p68, &p71, &p72}
	var p86 = sequenceParser{id: 86, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 114, 139, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 752, 804}}
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
	var p98 = choiceParser{id: 98, commit: 66, name: "bool", generalizations: []int{266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p91 = sequenceParser{id: 91, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p87 = charParser{id: 87, chars: []rune{116}}
	var p88 = charParser{id: 88, chars: []rune{114}}
	var p89 = charParser{id: 89, chars: []rune{117}}
	var p90 = charParser{id: 90, chars: []rune{101}}
	p91.items = []parser{&p87, &p88, &p89, &p90}
	var p97 = sequenceParser{id: 97, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{98, 266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p92 = charParser{id: 92, chars: []rune{102}}
	var p93 = charParser{id: 93, chars: []rune{97}}
	var p94 = charParser{id: 94, chars: []rune{108}}
	var p95 = charParser{id: 95, chars: []rune{115}}
	var p96 = charParser{id: 96, chars: []rune{101}}
	p97.items = []parser{&p92, &p93, &p94, &p95, &p96}
	p98.options = []parser{&p91, &p97}
	var p510 = sequenceParser{id: 510, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 114, 794, 197, 395, 332, 333, 334, 335, 336, 387, 514, 577, 570, 804}}
	var p507 = sequenceParser{id: 507, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p500 = charParser{id: 500, chars: []rune{114}}
	var p501 = charParser{id: 501, chars: []rune{101}}
	var p502 = charParser{id: 502, chars: []rune{99}}
	var p503 = charParser{id: 503, chars: []rune{101}}
	var p504 = charParser{id: 504, chars: []rune{105}}
	var p505 = charParser{id: 505, chars: []rune{118}}
	var p506 = charParser{id: 506, chars: []rune{101}}
	p507.items = []parser{&p500, &p501, &p502, &p503, &p504, &p505, &p506}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p508 = sequenceParser{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p508.items = []parser{&p836, &p14}
	p509.items = []parser{&p836, &p14, &p508}
	p510.items = []parser{&p507, &p509, &p836, &p266}
	var p103 = sequenceParser{id: 103, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{266, 114, 139, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 743, 804}}
	var p100 = sequenceParser{id: 100, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p99 = charParser{id: 99, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p100.items = []parser{&p99}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p102.items = []parser{&p101}
	p103.items = []parser{&p100, &p102}
	var p124 = sequenceParser{id: 124, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{114, 266, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
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
	p112.items = []parser{&p836, &p111}
	p113.items = []parser{&p111, &p112}
	var p118 = sequenceParser{id: 118, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p114 = choiceParser{id: 114, commit: 66, name: "list-item"}
	var p108 = sequenceParser{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{114, 147, 148}}
	var p107 = sequenceParser{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p104 = charParser{id: 104, chars: []rune{46}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	p107.items = []parser{&p104, &p105, &p106}
	p108.items = []parser{&p266, &p836, &p107}
	p114.options = []parser{&p395, &p108}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p115 = sequenceParser{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p115.items = []parser{&p113, &p836, &p114}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p116.items = []parser{&p836, &p115}
	p117.items = []parser{&p836, &p115, &p116}
	p118.items = []parser{&p114, &p117}
	var p122 = sequenceParser{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p121 = charParser{id: 121, chars: []rune{93}}
	p122.items = []parser{&p121}
	p123.items = []parser{&p120, &p836, &p113, &p836, &p118, &p836, &p113, &p836, &p122}
	p124.items = []parser{&p123}
	var p129 = sequenceParser{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p126 = sequenceParser{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p125 = charParser{id: 125, chars: []rune{126}}
	p126.items = []parser{&p125}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p127 = sequenceParser{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p127.items = []parser{&p836, &p14}
	p128.items = []parser{&p836, &p14, &p127}
	p129.items = []parser{&p126, &p128, &p836, &p123}
	var p158 = sequenceParser{id: 158, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{266, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
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
	p134.items = []parser{&p836, &p14}
	p135.items = []parser{&p836, &p14, &p134}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p136 = sequenceParser{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p136.items = []parser{&p836, &p14}
	p137.items = []parser{&p836, &p14, &p136}
	var p133 = sequenceParser{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p132 = charParser{id: 132, chars: []rune{93}}
	p133.items = []parser{&p132}
	p138.items = []parser{&p131, &p135, &p836, &p395, &p137, &p836, &p133}
	p139.options = []parser{&p103, &p86, &p138}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p142 = sequenceParser{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p142.items = []parser{&p836, &p14}
	p143.items = []parser{&p836, &p14, &p142}
	var p141 = sequenceParser{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p140 = charParser{id: 140, chars: []rune{58}}
	p141.items = []parser{&p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p836, &p14}
	p145.items = []parser{&p836, &p14, &p144}
	p146.items = []parser{&p139, &p143, &p836, &p141, &p145, &p836, &p395}
	p147.options = []parser{&p146, &p108}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p149 = sequenceParser{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p148 = choiceParser{id: 148, commit: 2}
	p148.options = []parser{&p146, &p108}
	p149.items = []parser{&p113, &p836, &p148}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p150.items = []parser{&p836, &p149}
	p151.items = []parser{&p836, &p149, &p150}
	p152.items = []parser{&p147, &p151}
	var p156 = sequenceParser{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p155 = charParser{id: 155, chars: []rune{125}}
	p156.items = []parser{&p155}
	p157.items = []parser{&p154, &p836, &p113, &p836, &p152, &p836, &p113, &p836, &p156}
	p158.items = []parser{&p157}
	var p163 = sequenceParser{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 794, 197, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p160 = sequenceParser{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p159 = charParser{id: 159, chars: []rune{126}}
	p160.items = []parser{&p159}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p161 = sequenceParser{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p161.items = []parser{&p836, &p14}
	p162.items = []parser{&p836, &p14, &p161}
	p163.items = []parser{&p160, &p162, &p836, &p157}
	var p206 = sequenceParser{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 197, 266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p203 = sequenceParser{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p201 = charParser{id: 201, chars: []rune{102}}
	var p202 = charParser{id: 202, chars: []rune{110}}
	p203.items = []parser{&p201, &p202}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p204 = sequenceParser{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p204.items = []parser{&p836, &p14}
	p205.items = []parser{&p836, &p14, &p204}
	var p200 = sequenceParser{id: 200, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p192 = sequenceParser{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p191 = charParser{id: 191, chars: []rune{40}}
	p192.items = []parser{&p191}
	var p194 = choiceParser{id: 194, commit: 2}
	var p167 = sequenceParser{id: 167, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{194}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p164.items = []parser{&p113, &p836, &p103}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p165.items = []parser{&p836, &p164}
	p166.items = []parser{&p836, &p164, &p165}
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
	p172.items = []parser{&p836, &p14}
	p173.items = []parser{&p836, &p14, &p172}
	p174.items = []parser{&p171, &p173, &p836, &p103}
	p193.items = []parser{&p167, &p836, &p113, &p836, &p174}
	p194.options = []parser{&p167, &p193, &p174}
	var p196 = sequenceParser{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p195 = charParser{id: 195, chars: []rune{41}}
	p196.items = []parser{&p195}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p198 = sequenceParser{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p198.items = []parser{&p836, &p14}
	p199.items = []parser{&p836, &p14, &p198}
	var p197 = choiceParser{id: 197, commit: 2}
	var p794 = choiceParser{id: 794, commit: 66, name: "simple-statement", generalizations: []int{197, 804}}
	var p499 = sequenceParser{id: 499, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 197, 514, 804}}
	var p494 = sequenceParser{id: 494, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p490 = charParser{id: 490, chars: []rune{115}}
	var p491 = charParser{id: 491, chars: []rune{101}}
	var p492 = charParser{id: 492, chars: []rune{110}}
	var p493 = charParser{id: 493, chars: []rune{100}}
	p494.items = []parser{&p490, &p491, &p492, &p493}
	var p496 = sequenceParser{id: 496, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p495 = sequenceParser{id: 495, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p495.items = []parser{&p836, &p14}
	p496.items = []parser{&p836, &p14, &p495}
	var p498 = sequenceParser{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p497 = sequenceParser{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p497.items = []parser{&p836, &p14}
	p498.items = []parser{&p836, &p14, &p497}
	p499.items = []parser{&p494, &p496, &p836, &p266, &p498, &p836, &p266}
	var p557 = sequenceParser{id: 557, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 197, 804}}
	var p554 = sequenceParser{id: 554, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p552 = charParser{id: 552, chars: []rune{103}}
	var p553 = charParser{id: 553, chars: []rune{111}}
	p554.items = []parser{&p552, &p553}
	var p556 = sequenceParser{id: 556, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p555 = sequenceParser{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p555.items = []parser{&p836, &p14}
	p556.items = []parser{&p836, &p14, &p555}
	var p256 = sequenceParser{id: 256, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p253 = sequenceParser{id: 253, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p252 = charParser{id: 252, chars: []rune{40}}
	p253.items = []parser{&p252}
	var p255 = sequenceParser{id: 255, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p254 = charParser{id: 254, chars: []rune{41}}
	p255.items = []parser{&p254}
	p256.items = []parser{&p266, &p836, &p253, &p836, &p113, &p836, &p118, &p836, &p113, &p836, &p255}
	p557.items = []parser{&p554, &p556, &p836, &p256}
	var p566 = sequenceParser{id: 566, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 197, 804}}
	var p563 = sequenceParser{id: 563, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p558 = charParser{id: 558, chars: []rune{100}}
	var p559 = charParser{id: 559, chars: []rune{101}}
	var p560 = charParser{id: 560, chars: []rune{102}}
	var p561 = charParser{id: 561, chars: []rune{101}}
	var p562 = charParser{id: 562, chars: []rune{114}}
	p563.items = []parser{&p558, &p559, &p560, &p561, &p562}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p564.items = []parser{&p836, &p14}
	p565.items = []parser{&p836, &p14, &p564}
	p566.items = []parser{&p563, &p565, &p836, &p256}
	var p631 = choiceParser{id: 631, commit: 64, name: "assignment", generalizations: []int{794, 197, 804}}
	var p611 = sequenceParser{id: 611, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{631, 794, 197, 804}}
	var p608 = sequenceParser{id: 608, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p605 = charParser{id: 605, chars: []rune{115}}
	var p606 = charParser{id: 606, chars: []rune{101}}
	var p607 = charParser{id: 607, chars: []rune{116}}
	p608.items = []parser{&p605, &p606, &p607}
	var p610 = sequenceParser{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p609 = sequenceParser{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p609.items = []parser{&p836, &p14}
	p610.items = []parser{&p836, &p14, &p609}
	var p600 = sequenceParser{id: 600, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p595.items = []parser{&p836, &p14}
	p596.items = []parser{&p14, &p595}
	var p594 = sequenceParser{id: 594, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p593 = charParser{id: 593, chars: []rune{61}}
	p594.items = []parser{&p593}
	p597.items = []parser{&p596, &p836, &p594}
	var p599 = sequenceParser{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p598.items = []parser{&p836, &p14}
	p599.items = []parser{&p836, &p14, &p598}
	p600.items = []parser{&p266, &p836, &p597, &p599, &p836, &p395}
	p611.items = []parser{&p608, &p610, &p836, &p600}
	var p618 = sequenceParser{id: 618, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{631, 794, 197, 804}}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p614.items = []parser{&p836, &p14}
	p615.items = []parser{&p836, &p14, &p614}
	var p613 = sequenceParser{id: 613, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p612 = charParser{id: 612, chars: []rune{61}}
	p613.items = []parser{&p612}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p616.items = []parser{&p836, &p14}
	p617.items = []parser{&p836, &p14, &p616}
	p618.items = []parser{&p266, &p615, &p836, &p613, &p617, &p836, &p395}
	var p630 = sequenceParser{id: 630, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{631, 794, 197, 804}}
	var p622 = sequenceParser{id: 622, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p619 = charParser{id: 619, chars: []rune{115}}
	var p620 = charParser{id: 620, chars: []rune{101}}
	var p621 = charParser{id: 621, chars: []rune{116}}
	p622.items = []parser{&p619, &p620, &p621}
	var p629 = sequenceParser{id: 629, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p628 = sequenceParser{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p628.items = []parser{&p836, &p14}
	p629.items = []parser{&p836, &p14, &p628}
	var p624 = sequenceParser{id: 624, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p623 = charParser{id: 623, chars: []rune{40}}
	p624.items = []parser{&p623}
	var p625 = sequenceParser{id: 625, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p604 = sequenceParser{id: 604, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p603 = sequenceParser{id: 603, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p601 = sequenceParser{id: 601, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p601.items = []parser{&p113, &p836, &p600}
	var p602 = sequenceParser{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p602.items = []parser{&p836, &p601}
	p603.items = []parser{&p836, &p601, &p602}
	p604.items = []parser{&p600, &p603}
	p625.items = []parser{&p113, &p836, &p604}
	var p627 = sequenceParser{id: 627, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p626 = charParser{id: 626, chars: []rune{41}}
	p627.items = []parser{&p626}
	p630.items = []parser{&p622, &p629, &p836, &p624, &p836, &p625, &p836, &p113, &p836, &p627}
	p631.options = []parser{&p611, &p618, &p630}
	var p803 = sequenceParser{id: 803, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794, 197, 804}}
	var p796 = sequenceParser{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p795 = charParser{id: 795, chars: []rune{40}}
	p796.items = []parser{&p795}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p799 = sequenceParser{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p799.items = []parser{&p836, &p14}
	p800.items = []parser{&p836, &p14, &p799}
	var p802 = sequenceParser{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p801.items = []parser{&p836, &p14}
	p802.items = []parser{&p836, &p14, &p801}
	var p798 = sequenceParser{id: 798, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p797 = charParser{id: 797, chars: []rune{41}}
	p798.items = []parser{&p797}
	p803.items = []parser{&p796, &p800, &p836, &p794, &p802, &p836, &p798}
	p794.options = []parser{&p499, &p557, &p566, &p631, &p803, &p395}
	var p190 = sequenceParser{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{197}}
	var p187 = sequenceParser{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p186 = charParser{id: 186, chars: []rune{123}}
	p187.items = []parser{&p186}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{125}}
	p189.items = []parser{&p188}
	p190.items = []parser{&p187, &p836, &p818, &p836, &p822, &p836, &p818, &p836, &p189}
	p197.options = []parser{&p794, &p190}
	p200.items = []parser{&p192, &p836, &p113, &p836, &p194, &p836, &p113, &p836, &p196, &p199, &p836, &p197}
	p206.items = []parser{&p203, &p205, &p836, &p200}
	var p216 = sequenceParser{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p209 = sequenceParser{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p207 = charParser{id: 207, chars: []rune{102}}
	var p208 = charParser{id: 208, chars: []rune{110}}
	p209.items = []parser{&p207, &p208}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p212 = sequenceParser{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p212.items = []parser{&p836, &p14}
	p213.items = []parser{&p836, &p14, &p212}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p210 = charParser{id: 210, chars: []rune{126}}
	p211.items = []parser{&p210}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p214.items = []parser{&p836, &p14}
	p215.items = []parser{&p836, &p14, &p214}
	p216.items = []parser{&p209, &p213, &p836, &p211, &p215, &p836, &p200}
	var p244 = choiceParser{id: 244, commit: 64, name: "expression-indexer", generalizations: []int{266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p234 = sequenceParser{id: 234, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{244, 266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p227 = sequenceParser{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p226 = charParser{id: 226, chars: []rune{91}}
	p227.items = []parser{&p226}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p230 = sequenceParser{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p230.items = []parser{&p836, &p14}
	p231.items = []parser{&p836, &p14, &p230}
	var p233 = sequenceParser{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p232 = sequenceParser{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p232.items = []parser{&p836, &p14}
	p233.items = []parser{&p836, &p14, &p232}
	var p229 = sequenceParser{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p228 = charParser{id: 228, chars: []rune{93}}
	p229.items = []parser{&p228}
	p234.items = []parser{&p266, &p836, &p227, &p231, &p836, &p395, &p233, &p836, &p229}
	var p243 = sequenceParser{id: 243, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{244, 266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p236 = sequenceParser{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p235 = charParser{id: 235, chars: []rune{91}}
	p236.items = []parser{&p235}
	var p240 = sequenceParser{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p239 = sequenceParser{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p239.items = []parser{&p836, &p14}
	p240.items = []parser{&p836, &p14, &p239}
	var p225 = sequenceParser{id: 225, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{570, 576, 577}}
	var p217 = sequenceParser{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p217.items = []parser{&p395}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p221 = sequenceParser{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p221.items = []parser{&p836, &p14}
	p222.items = []parser{&p836, &p14, &p221}
	var p220 = sequenceParser{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p219 = charParser{id: 219, chars: []rune{58}}
	p220.items = []parser{&p219}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p836, &p14}
	p224.items = []parser{&p836, &p14, &p223}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p395}
	p225.items = []parser{&p217, &p222, &p836, &p220, &p224, &p836, &p218}
	var p242 = sequenceParser{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p241 = sequenceParser{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p241.items = []parser{&p836, &p14}
	p242.items = []parser{&p836, &p14, &p241}
	var p238 = sequenceParser{id: 238, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p237 = charParser{id: 237, chars: []rune{93}}
	p238.items = []parser{&p237}
	p243.items = []parser{&p266, &p836, &p236, &p240, &p836, &p225, &p242, &p836, &p238}
	p244.options = []parser{&p234, &p243}
	var p251 = sequenceParser{id: 251, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p247.items = []parser{&p836, &p14}
	p248.items = []parser{&p836, &p14, &p247}
	var p246 = sequenceParser{id: 246, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p245 = charParser{id: 245, chars: []rune{46}}
	p246.items = []parser{&p245}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p249.items = []parser{&p836, &p14}
	p250.items = []parser{&p836, &p14, &p249}
	p251.items = []parser{&p266, &p248, &p836, &p246, &p250, &p836, &p103}
	var p265 = sequenceParser{id: 265, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{266, 395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
	var p258 = sequenceParser{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p257 = charParser{id: 257, chars: []rune{40}}
	p258.items = []parser{&p257}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p261 = sequenceParser{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p261.items = []parser{&p836, &p14}
	p262.items = []parser{&p836, &p14, &p261}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p263.items = []parser{&p836, &p14}
	p264.items = []parser{&p836, &p14, &p263}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{41}}
	p260.items = []parser{&p259}
	p265.items = []parser{&p258, &p262, &p836, &p395, &p264, &p836, &p260}
	p266.options = []parser{&p60, &p73, &p86, &p98, &p510, &p103, &p124, &p129, &p158, &p163, &p206, &p216, &p244, &p251, &p256, &p265}
	var p326 = sequenceParser{id: 326, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{395, 332, 333, 334, 335, 336, 387, 577, 570, 804}}
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
	p326.items = []parser{&p325, &p836, &p266}
	var p373 = choiceParser{id: 373, commit: 66, name: "binary-expression", generalizations: []int{395, 387, 577, 570, 804}}
	var p344 = sequenceParser{id: 344, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 333, 334, 335, 336, 395, 387, 577, 570, 804}}
	var p332 = choiceParser{id: 332, commit: 66, name: "operand0", generalizations: []int{333, 334, 335, 336}}
	p332.options = []parser{&p266, &p326}
	var p342 = sequenceParser{id: 342, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p339 = sequenceParser{id: 339, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p338 = sequenceParser{id: 338, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p338.items = []parser{&p836, &p14}
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
	p340.items = []parser{&p836, &p14}
	p341.items = []parser{&p836, &p14, &p340}
	p342.items = []parser{&p339, &p836, &p327, &p341, &p836, &p332}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p343.items = []parser{&p836, &p342}
	p344.items = []parser{&p332, &p836, &p342, &p343}
	var p351 = sequenceParser{id: 351, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 334, 335, 336, 395, 387, 577, 570, 804}}
	var p333 = choiceParser{id: 333, commit: 66, name: "operand1", generalizations: []int{334, 335, 336}}
	p333.options = []parser{&p332, &p344}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p345 = sequenceParser{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p345.items = []parser{&p836, &p14}
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
	p347.items = []parser{&p836, &p14}
	p348.items = []parser{&p836, &p14, &p347}
	p349.items = []parser{&p346, &p836, &p328, &p348, &p836, &p333}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p350.items = []parser{&p836, &p349}
	p351.items = []parser{&p333, &p836, &p349, &p350}
	var p358 = sequenceParser{id: 358, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 335, 336, 395, 387, 577, 570, 804}}
	var p334 = choiceParser{id: 334, commit: 66, name: "operand2", generalizations: []int{335, 336}}
	p334.options = []parser{&p333, &p351}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p352.items = []parser{&p836, &p14}
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
	p354.items = []parser{&p836, &p14}
	p355.items = []parser{&p836, &p14, &p354}
	p356.items = []parser{&p353, &p836, &p329, &p355, &p836, &p334}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p357.items = []parser{&p836, &p356}
	p358.items = []parser{&p334, &p836, &p356, &p357}
	var p365 = sequenceParser{id: 365, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 336, 395, 387, 577, 570, 804}}
	var p335 = choiceParser{id: 335, commit: 66, name: "operand3", generalizations: []int{336}}
	p335.options = []parser{&p334, &p358}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p836, &p14}
	p360.items = []parser{&p14, &p359}
	var p330 = sequenceParser{id: 330, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p318 = sequenceParser{id: 318, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p316 = charParser{id: 316, chars: []rune{38}}
	var p317 = charParser{id: 317, chars: []rune{38}}
	p318.items = []parser{&p316, &p317}
	p330.items = []parser{&p318}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p361.items = []parser{&p836, &p14}
	p362.items = []parser{&p836, &p14, &p361}
	p363.items = []parser{&p360, &p836, &p330, &p362, &p836, &p335}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p364.items = []parser{&p836, &p363}
	p365.items = []parser{&p335, &p836, &p363, &p364}
	var p372 = sequenceParser{id: 372, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{373, 395, 387, 577, 570, 804}}
	var p336 = choiceParser{id: 336, commit: 66, name: "operand4"}
	p336.options = []parser{&p335, &p365}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p366.items = []parser{&p836, &p14}
	p367.items = []parser{&p14, &p366}
	var p331 = sequenceParser{id: 331, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p321 = sequenceParser{id: 321, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p319 = charParser{id: 319, chars: []rune{124}}
	var p320 = charParser{id: 320, chars: []rune{124}}
	p321.items = []parser{&p319, &p320}
	p331.items = []parser{&p321}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p368.items = []parser{&p836, &p14}
	p369.items = []parser{&p836, &p14, &p368}
	p370.items = []parser{&p367, &p836, &p331, &p369, &p836, &p336}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p371.items = []parser{&p836, &p370}
	p372.items = []parser{&p336, &p836, &p370, &p371}
	p373.options = []parser{&p344, &p351, &p358, &p365, &p372}
	var p386 = sequenceParser{id: 386, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{395, 387, 577, 570, 804}}
	var p379 = sequenceParser{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p378 = sequenceParser{id: 378, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p378.items = []parser{&p836, &p14}
	p379.items = []parser{&p836, &p14, &p378}
	var p375 = sequenceParser{id: 375, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p374 = charParser{id: 374, chars: []rune{63}}
	p375.items = []parser{&p374}
	var p381 = sequenceParser{id: 381, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p380 = sequenceParser{id: 380, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p380.items = []parser{&p836, &p14}
	p381.items = []parser{&p836, &p14, &p380}
	var p383 = sequenceParser{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p382 = sequenceParser{id: 382, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p382.items = []parser{&p836, &p14}
	p383.items = []parser{&p836, &p14, &p382}
	var p377 = sequenceParser{id: 377, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p376 = charParser{id: 376, chars: []rune{58}}
	p377.items = []parser{&p376}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p384 = sequenceParser{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p384.items = []parser{&p836, &p14}
	p385.items = []parser{&p836, &p14, &p384}
	p386.items = []parser{&p395, &p379, &p836, &p375, &p381, &p836, &p395, &p383, &p836, &p377, &p385, &p836, &p395}
	var p394 = sequenceParser{id: 394, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{395, 577, 570, 804}}
	var p387 = choiceParser{id: 387, commit: 66, name: "chainingOperand"}
	p387.options = []parser{&p266, &p326, &p373, &p386}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p388.items = []parser{&p836, &p14}
	p389.items = []parser{&p14, &p388}
	var p324 = sequenceParser{id: 324, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p322 = charParser{id: 322, chars: []rune{45}}
	var p323 = charParser{id: 323, chars: []rune{62}}
	p324.items = []parser{&p322, &p323}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p390.items = []parser{&p836, &p14}
	p391.items = []parser{&p836, &p14, &p390}
	p392.items = []parser{&p389, &p836, &p324, &p391, &p836, &p387}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p393.items = []parser{&p836, &p392}
	p394.items = []parser{&p387, &p836, &p392, &p393}
	p395.options = []parser{&p266, &p326, &p373, &p386, &p394}
	p184.items = []parser{&p183, &p836, &p395}
	p185.items = []parser{&p181, &p836, &p184}
	var p432 = sequenceParser{id: 432, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{804, 478, 542}}
	var p398 = sequenceParser{id: 398, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p396 = charParser{id: 396, chars: []rune{105}}
	var p397 = charParser{id: 397, chars: []rune{102}}
	p398.items = []parser{&p396, &p397}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p426.items = []parser{&p836, &p14}
	p427.items = []parser{&p836, &p14, &p426}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p428.items = []parser{&p836, &p14}
	p429.items = []parser{&p836, &p14, &p428}
	var p431 = sequenceParser{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p407 = sequenceParser{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p407.items = []parser{&p836, &p14}
	p408.items = []parser{&p14, &p407}
	var p403 = sequenceParser{id: 403, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p399 = charParser{id: 399, chars: []rune{101}}
	var p400 = charParser{id: 400, chars: []rune{108}}
	var p401 = charParser{id: 401, chars: []rune{115}}
	var p402 = charParser{id: 402, chars: []rune{101}}
	p403.items = []parser{&p399, &p400, &p401, &p402}
	var p410 = sequenceParser{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p409 = sequenceParser{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p409.items = []parser{&p836, &p14}
	p410.items = []parser{&p836, &p14, &p409}
	var p406 = sequenceParser{id: 406, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p404 = charParser{id: 404, chars: []rune{105}}
	var p405 = charParser{id: 405, chars: []rune{102}}
	p406.items = []parser{&p404, &p405}
	var p412 = sequenceParser{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p411 = sequenceParser{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p411.items = []parser{&p836, &p14}
	p412.items = []parser{&p836, &p14, &p411}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p413 = sequenceParser{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p413.items = []parser{&p836, &p14}
	p414.items = []parser{&p836, &p14, &p413}
	p415.items = []parser{&p408, &p836, &p403, &p410, &p836, &p406, &p412, &p836, &p395, &p414, &p836, &p190}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p430.items = []parser{&p836, &p415}
	p431.items = []parser{&p836, &p415, &p430}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p421 = sequenceParser{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p421.items = []parser{&p836, &p14}
	p422.items = []parser{&p14, &p421}
	var p420 = sequenceParser{id: 420, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p416 = charParser{id: 416, chars: []rune{101}}
	var p417 = charParser{id: 417, chars: []rune{108}}
	var p418 = charParser{id: 418, chars: []rune{115}}
	var p419 = charParser{id: 419, chars: []rune{101}}
	p420.items = []parser{&p416, &p417, &p418, &p419}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p423.items = []parser{&p836, &p14}
	p424.items = []parser{&p836, &p14, &p423}
	p425.items = []parser{&p422, &p836, &p420, &p424, &p836, &p190}
	p432.items = []parser{&p398, &p427, &p836, &p395, &p429, &p836, &p190, &p431, &p836, &p425}
	var p489 = sequenceParser{id: 489, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{478, 804, 542}}
	var p474 = sequenceParser{id: 474, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p468 = charParser{id: 468, chars: []rune{115}}
	var p469 = charParser{id: 469, chars: []rune{119}}
	var p470 = charParser{id: 470, chars: []rune{105}}
	var p471 = charParser{id: 471, chars: []rune{116}}
	var p472 = charParser{id: 472, chars: []rune{99}}
	var p473 = charParser{id: 473, chars: []rune{104}}
	p474.items = []parser{&p468, &p469, &p470, &p471, &p472, &p473}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p485.items = []parser{&p836, &p14}
	p486.items = []parser{&p836, &p14, &p485}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p487.items = []parser{&p836, &p14}
	p488.items = []parser{&p836, &p14, &p487}
	var p476 = sequenceParser{id: 476, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p475 = charParser{id: 475, chars: []rune{123}}
	p476.items = []parser{&p475}
	var p482 = sequenceParser{id: 482, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p477 = choiceParser{id: 477, commit: 2}
	var p467 = sequenceParser{id: 467, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{477, 478}}
	var p462 = sequenceParser{id: 462, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p455 = sequenceParser{id: 455, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p451 = charParser{id: 451, chars: []rune{99}}
	var p452 = charParser{id: 452, chars: []rune{97}}
	var p453 = charParser{id: 453, chars: []rune{115}}
	var p454 = charParser{id: 454, chars: []rune{101}}
	p455.items = []parser{&p451, &p452, &p453, &p454}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p458 = sequenceParser{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p458.items = []parser{&p836, &p14}
	p459.items = []parser{&p836, &p14, &p458}
	var p461 = sequenceParser{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p460.items = []parser{&p836, &p14}
	p461.items = []parser{&p836, &p14, &p460}
	var p457 = sequenceParser{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p456 = charParser{id: 456, chars: []rune{58}}
	p457.items = []parser{&p456}
	p462.items = []parser{&p455, &p459, &p836, &p395, &p461, &p836, &p457}
	var p466 = sequenceParser{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p464 = sequenceParser{id: 464, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p463 = charParser{id: 463, chars: []rune{59}}
	p464.items = []parser{&p463}
	var p465 = sequenceParser{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p465.items = []parser{&p836, &p464}
	p466.items = []parser{&p836, &p464, &p465}
	p467.items = []parser{&p462, &p466, &p836, &p804}
	var p450 = sequenceParser{id: 450, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{477, 478, 541, 542}}
	var p445 = sequenceParser{id: 445, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p440 = sequenceParser{id: 440, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p433 = charParser{id: 433, chars: []rune{100}}
	var p434 = charParser{id: 434, chars: []rune{101}}
	var p435 = charParser{id: 435, chars: []rune{102}}
	var p436 = charParser{id: 436, chars: []rune{97}}
	var p437 = charParser{id: 437, chars: []rune{117}}
	var p438 = charParser{id: 438, chars: []rune{108}}
	var p439 = charParser{id: 439, chars: []rune{116}}
	p440.items = []parser{&p433, &p434, &p435, &p436, &p437, &p438, &p439}
	var p444 = sequenceParser{id: 444, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p443 = sequenceParser{id: 443, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p443.items = []parser{&p836, &p14}
	p444.items = []parser{&p836, &p14, &p443}
	var p442 = sequenceParser{id: 442, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p441 = charParser{id: 441, chars: []rune{58}}
	p442.items = []parser{&p441}
	p445.items = []parser{&p440, &p444, &p836, &p442}
	var p449 = sequenceParser{id: 449, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p447 = sequenceParser{id: 447, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p446 = charParser{id: 446, chars: []rune{59}}
	p447.items = []parser{&p446}
	var p448 = sequenceParser{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p448.items = []parser{&p836, &p447}
	p449.items = []parser{&p836, &p447, &p448}
	p450.items = []parser{&p445, &p449, &p836, &p804}
	p477.options = []parser{&p467, &p450}
	var p481 = sequenceParser{id: 481, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p479 = sequenceParser{id: 479, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p478 = choiceParser{id: 478, commit: 2}
	p478.options = []parser{&p467, &p450, &p804}
	p479.items = []parser{&p818, &p836, &p478}
	var p480 = sequenceParser{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p480.items = []parser{&p836, &p479}
	p481.items = []parser{&p836, &p479, &p480}
	p482.items = []parser{&p477, &p481}
	var p484 = sequenceParser{id: 484, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p483 = charParser{id: 483, chars: []rune{125}}
	p484.items = []parser{&p483}
	p489.items = []parser{&p474, &p486, &p836, &p395, &p488, &p836, &p476, &p836, &p818, &p836, &p482, &p836, &p818, &p836, &p484}
	var p551 = sequenceParser{id: 551, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{542, 804}}
	var p538 = sequenceParser{id: 538, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p532 = charParser{id: 532, chars: []rune{115}}
	var p533 = charParser{id: 533, chars: []rune{101}}
	var p534 = charParser{id: 534, chars: []rune{108}}
	var p535 = charParser{id: 535, chars: []rune{101}}
	var p536 = charParser{id: 536, chars: []rune{99}}
	var p537 = charParser{id: 537, chars: []rune{116}}
	p538.items = []parser{&p532, &p533, &p534, &p535, &p536, &p537}
	var p550 = sequenceParser{id: 550, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p549 = sequenceParser{id: 549, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p549.items = []parser{&p836, &p14}
	p550.items = []parser{&p836, &p14, &p549}
	var p540 = sequenceParser{id: 540, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p539 = charParser{id: 539, chars: []rune{123}}
	p540.items = []parser{&p539}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p541 = choiceParser{id: 541, commit: 2}
	var p531 = sequenceParser{id: 531, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{541, 542}}
	var p526 = sequenceParser{id: 526, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p519 = sequenceParser{id: 519, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p515 = charParser{id: 515, chars: []rune{99}}
	var p516 = charParser{id: 516, chars: []rune{97}}
	var p517 = charParser{id: 517, chars: []rune{115}}
	var p518 = charParser{id: 518, chars: []rune{101}}
	p519.items = []parser{&p515, &p516, &p517, &p518}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p522 = sequenceParser{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p522.items = []parser{&p836, &p14}
	p523.items = []parser{&p836, &p14, &p522}
	var p514 = choiceParser{id: 514, commit: 66, name: "communication"}
	var p513 = sequenceParser{id: 513, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{514}}
	var p512 = sequenceParser{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p511 = sequenceParser{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p511.items = []parser{&p836, &p14}
	p512.items = []parser{&p836, &p14, &p511}
	p513.items = []parser{&p103, &p512, &p836, &p510}
	p514.options = []parser{&p499, &p510, &p513}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p524.items = []parser{&p836, &p14}
	p525.items = []parser{&p836, &p14, &p524}
	var p521 = sequenceParser{id: 521, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p520 = charParser{id: 520, chars: []rune{58}}
	p521.items = []parser{&p520}
	p526.items = []parser{&p519, &p523, &p836, &p514, &p525, &p836, &p521}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p528 = sequenceParser{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p527 = charParser{id: 527, chars: []rune{59}}
	p528.items = []parser{&p527}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p529.items = []parser{&p836, &p528}
	p530.items = []parser{&p836, &p528, &p529}
	p531.items = []parser{&p526, &p530, &p836, &p804}
	p541.options = []parser{&p531, &p450}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p543 = sequenceParser{id: 543, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p542 = choiceParser{id: 542, commit: 2}
	p542.options = []parser{&p531, &p450, &p804}
	p543.items = []parser{&p818, &p836, &p542}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p544.items = []parser{&p836, &p543}
	p545.items = []parser{&p836, &p543, &p544}
	p546.items = []parser{&p541, &p545}
	var p548 = sequenceParser{id: 548, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p547 = charParser{id: 547, chars: []rune{125}}
	p548.items = []parser{&p547}
	p551.items = []parser{&p538, &p550, &p836, &p540, &p836, &p818, &p836, &p546, &p836, &p818, &p836, &p548}
	var p592 = sequenceParser{id: 592, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{804}}
	var p581 = sequenceParser{id: 581, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p578 = charParser{id: 578, chars: []rune{102}}
	var p579 = charParser{id: 579, chars: []rune{111}}
	var p580 = charParser{id: 580, chars: []rune{114}}
	p581.items = []parser{&p578, &p579, &p580}
	var p591 = choiceParser{id: 591, commit: 2}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{591}}
	var p584 = sequenceParser{id: 584, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p583 = sequenceParser{id: 583, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p582 = sequenceParser{id: 582, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p582.items = []parser{&p836, &p14}
	p583.items = []parser{&p14, &p582}
	var p577 = choiceParser{id: 577, commit: 66, name: "loop-expression"}
	var p576 = choiceParser{id: 576, commit: 64, name: "range-over-expression", generalizations: []int{577}}
	var p575 = sequenceParser{id: 575, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{576, 577}}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p571 = sequenceParser{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p571.items = []parser{&p836, &p14}
	p572.items = []parser{&p836, &p14, &p571}
	var p569 = sequenceParser{id: 569, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p567 = charParser{id: 567, chars: []rune{105}}
	var p568 = charParser{id: 568, chars: []rune{110}}
	p569.items = []parser{&p567, &p568}
	var p574 = sequenceParser{id: 574, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p573.items = []parser{&p836, &p14}
	p574.items = []parser{&p836, &p14, &p573}
	var p570 = choiceParser{id: 570, commit: 2}
	p570.options = []parser{&p395, &p225}
	p575.items = []parser{&p103, &p572, &p836, &p569, &p574, &p836, &p570}
	p576.options = []parser{&p575, &p225}
	p577.options = []parser{&p395, &p576}
	p584.items = []parser{&p583, &p836, &p577}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p585.items = []parser{&p836, &p14}
	p586.items = []parser{&p836, &p14, &p585}
	p587.items = []parser{&p584, &p586, &p836, &p190}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{591}}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p588.items = []parser{&p836, &p14}
	p589.items = []parser{&p14, &p588}
	p590.items = []parser{&p589, &p836, &p190}
	p591.options = []parser{&p587, &p590}
	p592.items = []parser{&p581, &p836, &p591}
	var p740 = choiceParser{id: 740, commit: 66, name: "definition", generalizations: []int{804}}
	var p653 = sequenceParser{id: 653, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 804}}
	var p649 = sequenceParser{id: 649, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p646 = charParser{id: 646, chars: []rune{108}}
	var p647 = charParser{id: 647, chars: []rune{101}}
	var p648 = charParser{id: 648, chars: []rune{116}}
	p649.items = []parser{&p646, &p647, &p648}
	var p652 = sequenceParser{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p651 = sequenceParser{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p651.items = []parser{&p836, &p14}
	p652.items = []parser{&p836, &p14, &p651}
	var p650 = choiceParser{id: 650, commit: 2}
	var p640 = sequenceParser{id: 640, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{650, 654, 655}}
	var p639 = sequenceParser{id: 639, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p634.items = []parser{&p836, &p14}
	p635.items = []parser{&p14, &p634}
	var p633 = sequenceParser{id: 633, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p632 = charParser{id: 632, chars: []rune{61}}
	p633.items = []parser{&p632}
	p636.items = []parser{&p635, &p836, &p633}
	var p638 = sequenceParser{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p637 = sequenceParser{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p637.items = []parser{&p836, &p14}
	p638.items = []parser{&p836, &p14, &p637}
	p639.items = []parser{&p103, &p836, &p636, &p638, &p836, &p395}
	p640.items = []parser{&p639}
	var p645 = sequenceParser{id: 645, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{650, 654, 655}}
	var p642 = sequenceParser{id: 642, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p641 = charParser{id: 641, chars: []rune{126}}
	p642.items = []parser{&p641}
	var p644 = sequenceParser{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p643 = sequenceParser{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p643.items = []parser{&p836, &p14}
	p644.items = []parser{&p836, &p14, &p643}
	p645.items = []parser{&p642, &p644, &p836, &p639}
	p650.options = []parser{&p640, &p645}
	p653.items = []parser{&p649, &p652, &p836, &p650}
	var p674 = sequenceParser{id: 674, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 804}}
	var p667 = sequenceParser{id: 667, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p664 = charParser{id: 664, chars: []rune{108}}
	var p665 = charParser{id: 665, chars: []rune{101}}
	var p666 = charParser{id: 666, chars: []rune{116}}
	p667.items = []parser{&p664, &p665, &p666}
	var p673 = sequenceParser{id: 673, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p672 = sequenceParser{id: 672, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p672.items = []parser{&p836, &p14}
	p673.items = []parser{&p836, &p14, &p672}
	var p669 = sequenceParser{id: 669, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p668 = charParser{id: 668, chars: []rune{40}}
	p669.items = []parser{&p668}
	var p659 = sequenceParser{id: 659, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p654 = choiceParser{id: 654, commit: 2}
	p654.options = []parser{&p640, &p645}
	var p658 = sequenceParser{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p656 = sequenceParser{id: 656, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p655 = choiceParser{id: 655, commit: 2}
	p655.options = []parser{&p640, &p645}
	p656.items = []parser{&p113, &p836, &p655}
	var p657 = sequenceParser{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p657.items = []parser{&p836, &p656}
	p658.items = []parser{&p836, &p656, &p657}
	p659.items = []parser{&p654, &p658}
	var p671 = sequenceParser{id: 671, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p670 = charParser{id: 670, chars: []rune{41}}
	p671.items = []parser{&p670}
	p674.items = []parser{&p667, &p673, &p836, &p669, &p836, &p113, &p836, &p659, &p836, &p113, &p836, &p671}
	var p689 = sequenceParser{id: 689, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 804}}
	var p678 = sequenceParser{id: 678, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p675 = charParser{id: 675, chars: []rune{108}}
	var p676 = charParser{id: 676, chars: []rune{101}}
	var p677 = charParser{id: 677, chars: []rune{116}}
	p678.items = []parser{&p675, &p676, &p677}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p685 = sequenceParser{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p685.items = []parser{&p836, &p14}
	p686.items = []parser{&p836, &p14, &p685}
	var p680 = sequenceParser{id: 680, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p679 = charParser{id: 679, chars: []rune{126}}
	p680.items = []parser{&p679}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p687.items = []parser{&p836, &p14}
	p688.items = []parser{&p836, &p14, &p687}
	var p682 = sequenceParser{id: 682, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p681 = charParser{id: 681, chars: []rune{40}}
	p682.items = []parser{&p681}
	var p663 = sequenceParser{id: 663, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p662 = sequenceParser{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p660 = sequenceParser{id: 660, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p660.items = []parser{&p113, &p836, &p640}
	var p661 = sequenceParser{id: 661, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p661.items = []parser{&p836, &p660}
	p662.items = []parser{&p836, &p660, &p661}
	p663.items = []parser{&p640, &p662}
	var p684 = sequenceParser{id: 684, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p683 = charParser{id: 683, chars: []rune{41}}
	p684.items = []parser{&p683}
	p689.items = []parser{&p678, &p686, &p836, &p680, &p688, &p836, &p682, &p836, &p113, &p836, &p663, &p836, &p113, &p836, &p684}
	var p705 = sequenceParser{id: 705, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 804}}
	var p701 = sequenceParser{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p699 = charParser{id: 699, chars: []rune{102}}
	var p700 = charParser{id: 700, chars: []rune{110}}
	p701.items = []parser{&p699, &p700}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p703.items = []parser{&p836, &p14}
	p704.items = []parser{&p836, &p14, &p703}
	var p702 = choiceParser{id: 702, commit: 2}
	var p693 = sequenceParser{id: 693, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{702, 710, 711}}
	var p692 = sequenceParser{id: 692, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p691 = sequenceParser{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p690 = sequenceParser{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p690.items = []parser{&p836, &p14}
	p691.items = []parser{&p836, &p14, &p690}
	p692.items = []parser{&p103, &p691, &p836, &p200}
	p693.items = []parser{&p692}
	var p698 = sequenceParser{id: 698, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{702, 710, 711}}
	var p695 = sequenceParser{id: 695, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p694 = charParser{id: 694, chars: []rune{126}}
	p695.items = []parser{&p694}
	var p697 = sequenceParser{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p696.items = []parser{&p836, &p14}
	p697.items = []parser{&p836, &p14, &p696}
	p698.items = []parser{&p695, &p697, &p836, &p692}
	p702.options = []parser{&p693, &p698}
	p705.items = []parser{&p701, &p704, &p836, &p702}
	var p725 = sequenceParser{id: 725, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 804}}
	var p718 = sequenceParser{id: 718, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p716 = charParser{id: 716, chars: []rune{102}}
	var p717 = charParser{id: 717, chars: []rune{110}}
	p718.items = []parser{&p716, &p717}
	var p724 = sequenceParser{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p723 = sequenceParser{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p723.items = []parser{&p836, &p14}
	p724.items = []parser{&p836, &p14, &p723}
	var p720 = sequenceParser{id: 720, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p719 = charParser{id: 719, chars: []rune{40}}
	p720.items = []parser{&p719}
	var p715 = sequenceParser{id: 715, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p710 = choiceParser{id: 710, commit: 2}
	p710.options = []parser{&p693, &p698}
	var p714 = sequenceParser{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p712 = sequenceParser{id: 712, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p711 = choiceParser{id: 711, commit: 2}
	p711.options = []parser{&p693, &p698}
	p712.items = []parser{&p113, &p836, &p711}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p713.items = []parser{&p836, &p712}
	p714.items = []parser{&p836, &p712, &p713}
	p715.items = []parser{&p710, &p714}
	var p722 = sequenceParser{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p721 = charParser{id: 721, chars: []rune{41}}
	p722.items = []parser{&p721}
	p725.items = []parser{&p718, &p724, &p836, &p720, &p836, &p113, &p836, &p715, &p836, &p113, &p836, &p722}
	var p739 = sequenceParser{id: 739, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 804}}
	var p728 = sequenceParser{id: 728, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p726 = charParser{id: 726, chars: []rune{102}}
	var p727 = charParser{id: 727, chars: []rune{110}}
	p728.items = []parser{&p726, &p727}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p735 = sequenceParser{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p735.items = []parser{&p836, &p14}
	p736.items = []parser{&p836, &p14, &p735}
	var p730 = sequenceParser{id: 730, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p729 = charParser{id: 729, chars: []rune{126}}
	p730.items = []parser{&p729}
	var p738 = sequenceParser{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p737.items = []parser{&p836, &p14}
	p738.items = []parser{&p836, &p14, &p737}
	var p732 = sequenceParser{id: 732, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p731 = charParser{id: 731, chars: []rune{40}}
	p732.items = []parser{&p731}
	var p709 = sequenceParser{id: 709, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p708 = sequenceParser{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p706 = sequenceParser{id: 706, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p706.items = []parser{&p113, &p836, &p693}
	var p707 = sequenceParser{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p707.items = []parser{&p836, &p706}
	p708.items = []parser{&p836, &p706, &p707}
	p709.items = []parser{&p693, &p708}
	var p734 = sequenceParser{id: 734, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p733 = charParser{id: 733, chars: []rune{41}}
	p734.items = []parser{&p733}
	p739.items = []parser{&p728, &p736, &p836, &p730, &p738, &p836, &p732, &p836, &p113, &p836, &p709, &p836, &p113, &p836, &p734}
	p740.options = []parser{&p653, &p674, &p689, &p705, &p725, &p739}
	var p783 = choiceParser{id: 783, commit: 64, name: "require", generalizations: []int{804}}
	var p767 = sequenceParser{id: 767, commit: 66, name: "require-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 804}}
	var p764 = sequenceParser{id: 764, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p757 = charParser{id: 757, chars: []rune{114}}
	var p758 = charParser{id: 758, chars: []rune{101}}
	var p759 = charParser{id: 759, chars: []rune{113}}
	var p760 = charParser{id: 760, chars: []rune{117}}
	var p761 = charParser{id: 761, chars: []rune{105}}
	var p762 = charParser{id: 762, chars: []rune{114}}
	var p763 = charParser{id: 763, chars: []rune{101}}
	p764.items = []parser{&p757, &p758, &p759, &p760, &p761, &p762, &p763}
	var p766 = sequenceParser{id: 766, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p765 = sequenceParser{id: 765, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p765.items = []parser{&p836, &p14}
	p766.items = []parser{&p836, &p14, &p765}
	var p752 = choiceParser{id: 752, commit: 64, name: "require-fact"}
	var p751 = sequenceParser{id: 751, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{752}}
	var p743 = choiceParser{id: 743, commit: 2}
	var p742 = sequenceParser{id: 742, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{743}}
	var p741 = charParser{id: 741, chars: []rune{46}}
	p742.items = []parser{&p741}
	p743.options = []parser{&p103, &p742}
	var p748 = sequenceParser{id: 748, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p747 = sequenceParser{id: 747, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p746 = sequenceParser{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p746.items = []parser{&p836, &p14}
	p747.items = []parser{&p14, &p746}
	var p745 = sequenceParser{id: 745, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p744 = charParser{id: 744, chars: []rune{61}}
	p745.items = []parser{&p744}
	p748.items = []parser{&p747, &p836, &p745}
	var p750 = sequenceParser{id: 750, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p749 = sequenceParser{id: 749, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p749.items = []parser{&p836, &p14}
	p750.items = []parser{&p836, &p14, &p749}
	p751.items = []parser{&p743, &p836, &p748, &p750, &p836, &p86}
	p752.options = []parser{&p86, &p751}
	p767.items = []parser{&p764, &p766, &p836, &p752}
	var p782 = sequenceParser{id: 782, commit: 66, name: "require-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{783, 804}}
	var p775 = sequenceParser{id: 775, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p768 = charParser{id: 768, chars: []rune{114}}
	var p769 = charParser{id: 769, chars: []rune{101}}
	var p770 = charParser{id: 770, chars: []rune{113}}
	var p771 = charParser{id: 771, chars: []rune{117}}
	var p772 = charParser{id: 772, chars: []rune{105}}
	var p773 = charParser{id: 773, chars: []rune{114}}
	var p774 = charParser{id: 774, chars: []rune{101}}
	p775.items = []parser{&p768, &p769, &p770, &p771, &p772, &p773, &p774}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p780 = sequenceParser{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p780.items = []parser{&p836, &p14}
	p781.items = []parser{&p836, &p14, &p780}
	var p777 = sequenceParser{id: 777, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p776 = charParser{id: 776, chars: []rune{40}}
	p777.items = []parser{&p776}
	var p756 = sequenceParser{id: 756, commit: 66, name: "require-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p753.items = []parser{&p113, &p836, &p752}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p754.items = []parser{&p836, &p753}
	p755.items = []parser{&p836, &p753, &p754}
	p756.items = []parser{&p752, &p755}
	var p779 = sequenceParser{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p778 = charParser{id: 778, chars: []rune{41}}
	p779.items = []parser{&p778}
	p782.items = []parser{&p775, &p781, &p836, &p777, &p836, &p113, &p836, &p756, &p836, &p113, &p836, &p779}
	p783.options = []parser{&p767, &p782}
	var p793 = sequenceParser{id: 793, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{804}}
	var p790 = sequenceParser{id: 790, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p784 = charParser{id: 784, chars: []rune{101}}
	var p785 = charParser{id: 785, chars: []rune{120}}
	var p786 = charParser{id: 786, chars: []rune{112}}
	var p787 = charParser{id: 787, chars: []rune{111}}
	var p788 = charParser{id: 788, chars: []rune{114}}
	var p789 = charParser{id: 789, chars: []rune{116}}
	p790.items = []parser{&p784, &p785, &p786, &p787, &p788, &p789}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p791.items = []parser{&p836, &p14}
	p792.items = []parser{&p836, &p14, &p791}
	p793.items = []parser{&p790, &p792, &p836, &p740}
	var p813 = sequenceParser{id: 813, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{804}}
	var p806 = sequenceParser{id: 806, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p805 = charParser{id: 805, chars: []rune{40}}
	p806.items = []parser{&p805}
	var p810 = sequenceParser{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p809 = sequenceParser{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p809.items = []parser{&p836, &p14}
	p810.items = []parser{&p836, &p14, &p809}
	var p812 = sequenceParser{id: 812, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p811 = sequenceParser{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p811.items = []parser{&p836, &p14}
	p812.items = []parser{&p836, &p14, &p811}
	var p808 = sequenceParser{id: 808, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p807 = charParser{id: 807, chars: []rune{41}}
	p808.items = []parser{&p807}
	p813.items = []parser{&p806, &p810, &p836, &p804, &p812, &p836, &p808}
	p804.options = []parser{&p185, &p432, &p489, &p551, &p592, &p740, &p783, &p793, &p813, &p794}
	var p821 = sequenceParser{id: 821, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p819 = sequenceParser{id: 819, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p819.items = []parser{&p818, &p836, &p804}
	var p820 = sequenceParser{id: 820, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p820.items = []parser{&p836, &p819}
	p821.items = []parser{&p836, &p819, &p820}
	p822.items = []parser{&p804, &p821}
	p837.items = []parser{&p833, &p836, &p818, &p836, &p822, &p836, &p818}
	p838.items = []parser{&p836, &p837, &p836}
	var b838 = sequenceBuilder{id: 838, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b836 = choiceBuilder{id: 836, commit: 2}
	var b834 = choiceBuilder{id: 834, commit: 70}
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
	b834.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b835 = sequenceBuilder{id: 835, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b39.items = []builder{&b14, &b836, &b38}
	var b40 = sequenceBuilder{id: 40, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b40.items = []builder{&b836, &b39}
	b41.items = []builder{&b836, &b39, &b40}
	b42.items = []builder{&b38, &b41}
	b835.items = []builder{&b42}
	b836.options = []builder{&b834, &b835}
	var b837 = sequenceBuilder{id: 837, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b833 = sequenceBuilder{id: 833, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b830 = sequenceBuilder{id: 830, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b828 = charBuilder{}
	var b829 = charBuilder{}
	b830.items = []builder{&b828, &b829}
	var b827 = sequenceBuilder{id: 827, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b826 = sequenceBuilder{id: 826, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b824 = sequenceBuilder{id: 824, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b823 = charBuilder{}
	b824.items = []builder{&b823}
	var b825 = sequenceBuilder{id: 825, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b825.items = []builder{&b836, &b824}
	b826.items = []builder{&b824, &b825}
	b827.items = []builder{&b826}
	var b832 = sequenceBuilder{id: 832, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b831 = charBuilder{}
	b832.items = []builder{&b831}
	b833.items = []builder{&b830, &b836, &b827, &b836, &b832}
	var b818 = sequenceBuilder{id: 818, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b816 = choiceBuilder{id: 816, commit: 2}
	var b815 = sequenceBuilder{id: 815, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b814 = charBuilder{}
	b815.items = []builder{&b814}
	b816.options = []builder{&b815, &b14}
	var b817 = sequenceBuilder{id: 817, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b817.items = []builder{&b836, &b816}
	b818.items = []builder{&b816, &b817}
	var b822 = sequenceBuilder{id: 822, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b804 = choiceBuilder{id: 804, commit: 66}
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
	b182.items = []builder{&b836, &b14}
	b183.items = []builder{&b14, &b182}
	var b395 = choiceBuilder{id: 395, commit: 66}
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
	var b510 = sequenceBuilder{id: 510, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b507 = sequenceBuilder{id: 507, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b500 = charBuilder{}
	var b501 = charBuilder{}
	var b502 = charBuilder{}
	var b503 = charBuilder{}
	var b504 = charBuilder{}
	var b505 = charBuilder{}
	var b506 = charBuilder{}
	b507.items = []builder{&b500, &b501, &b502, &b503, &b504, &b505, &b506}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b508 = sequenceBuilder{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b508.items = []builder{&b836, &b14}
	b509.items = []builder{&b836, &b14, &b508}
	b510.items = []builder{&b507, &b509, &b836, &b266}
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
	b112.items = []builder{&b836, &b111}
	b113.items = []builder{&b111, &b112}
	var b118 = sequenceBuilder{id: 118, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b114 = choiceBuilder{id: 114, commit: 66}
	var b108 = sequenceBuilder{id: 108, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b107 = sequenceBuilder{id: 107, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b104 = charBuilder{}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	b107.items = []builder{&b104, &b105, &b106}
	b108.items = []builder{&b266, &b836, &b107}
	b114.options = []builder{&b395, &b108}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b115 = sequenceBuilder{id: 115, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b115.items = []builder{&b113, &b836, &b114}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b116.items = []builder{&b836, &b115}
	b117.items = []builder{&b836, &b115, &b116}
	b118.items = []builder{&b114, &b117}
	var b122 = sequenceBuilder{id: 122, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b121 = charBuilder{}
	b122.items = []builder{&b121}
	b123.items = []builder{&b120, &b836, &b113, &b836, &b118, &b836, &b113, &b836, &b122}
	b124.items = []builder{&b123}
	var b129 = sequenceBuilder{id: 129, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b126 = sequenceBuilder{id: 126, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b125 = charBuilder{}
	b126.items = []builder{&b125}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b127 = sequenceBuilder{id: 127, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b127.items = []builder{&b836, &b14}
	b128.items = []builder{&b836, &b14, &b127}
	b129.items = []builder{&b126, &b128, &b836, &b123}
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
	b134.items = []builder{&b836, &b14}
	b135.items = []builder{&b836, &b14, &b134}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b136 = sequenceBuilder{id: 136, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b136.items = []builder{&b836, &b14}
	b137.items = []builder{&b836, &b14, &b136}
	var b133 = sequenceBuilder{id: 133, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b132 = charBuilder{}
	b133.items = []builder{&b132}
	b138.items = []builder{&b131, &b135, &b836, &b395, &b137, &b836, &b133}
	b139.options = []builder{&b103, &b86, &b138}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b142 = sequenceBuilder{id: 142, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b142.items = []builder{&b836, &b14}
	b143.items = []builder{&b836, &b14, &b142}
	var b141 = sequenceBuilder{id: 141, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b140 = charBuilder{}
	b141.items = []builder{&b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b836, &b14}
	b145.items = []builder{&b836, &b14, &b144}
	b146.items = []builder{&b139, &b143, &b836, &b141, &b145, &b836, &b395}
	b147.options = []builder{&b146, &b108}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b149 = sequenceBuilder{id: 149, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b148 = choiceBuilder{id: 148, commit: 2}
	b148.options = []builder{&b146, &b108}
	b149.items = []builder{&b113, &b836, &b148}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b150.items = []builder{&b836, &b149}
	b151.items = []builder{&b836, &b149, &b150}
	b152.items = []builder{&b147, &b151}
	var b156 = sequenceBuilder{id: 156, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b155 = charBuilder{}
	b156.items = []builder{&b155}
	b157.items = []builder{&b154, &b836, &b113, &b836, &b152, &b836, &b113, &b836, &b156}
	b158.items = []builder{&b157}
	var b163 = sequenceBuilder{id: 163, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b160 = sequenceBuilder{id: 160, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b159 = charBuilder{}
	b160.items = []builder{&b159}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b161 = sequenceBuilder{id: 161, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b161.items = []builder{&b836, &b14}
	b162.items = []builder{&b836, &b14, &b161}
	b163.items = []builder{&b160, &b162, &b836, &b157}
	var b206 = sequenceBuilder{id: 206, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b203 = sequenceBuilder{id: 203, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b201 = charBuilder{}
	var b202 = charBuilder{}
	b203.items = []builder{&b201, &b202}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b204 = sequenceBuilder{id: 204, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b204.items = []builder{&b836, &b14}
	b205.items = []builder{&b836, &b14, &b204}
	var b200 = sequenceBuilder{id: 200, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b192 = sequenceBuilder{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b191 = charBuilder{}
	b192.items = []builder{&b191}
	var b194 = choiceBuilder{id: 194, commit: 2}
	var b167 = sequenceBuilder{id: 167, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b164.items = []builder{&b113, &b836, &b103}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b165.items = []builder{&b836, &b164}
	b166.items = []builder{&b836, &b164, &b165}
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
	b172.items = []builder{&b836, &b14}
	b173.items = []builder{&b836, &b14, &b172}
	b174.items = []builder{&b171, &b173, &b836, &b103}
	b193.items = []builder{&b167, &b836, &b113, &b836, &b174}
	b194.options = []builder{&b167, &b193, &b174}
	var b196 = sequenceBuilder{id: 196, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b195 = charBuilder{}
	b196.items = []builder{&b195}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b198 = sequenceBuilder{id: 198, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b198.items = []builder{&b836, &b14}
	b199.items = []builder{&b836, &b14, &b198}
	var b197 = choiceBuilder{id: 197, commit: 2}
	var b794 = choiceBuilder{id: 794, commit: 66}
	var b499 = sequenceBuilder{id: 499, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b494 = sequenceBuilder{id: 494, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b490 = charBuilder{}
	var b491 = charBuilder{}
	var b492 = charBuilder{}
	var b493 = charBuilder{}
	b494.items = []builder{&b490, &b491, &b492, &b493}
	var b496 = sequenceBuilder{id: 496, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b495 = sequenceBuilder{id: 495, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b495.items = []builder{&b836, &b14}
	b496.items = []builder{&b836, &b14, &b495}
	var b498 = sequenceBuilder{id: 498, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b497 = sequenceBuilder{id: 497, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b497.items = []builder{&b836, &b14}
	b498.items = []builder{&b836, &b14, &b497}
	b499.items = []builder{&b494, &b496, &b836, &b266, &b498, &b836, &b266}
	var b557 = sequenceBuilder{id: 557, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b554 = sequenceBuilder{id: 554, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b552 = charBuilder{}
	var b553 = charBuilder{}
	b554.items = []builder{&b552, &b553}
	var b556 = sequenceBuilder{id: 556, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b555 = sequenceBuilder{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b555.items = []builder{&b836, &b14}
	b556.items = []builder{&b836, &b14, &b555}
	var b256 = sequenceBuilder{id: 256, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b253 = sequenceBuilder{id: 253, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b252 = charBuilder{}
	b253.items = []builder{&b252}
	var b255 = sequenceBuilder{id: 255, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b254 = charBuilder{}
	b255.items = []builder{&b254}
	b256.items = []builder{&b266, &b836, &b253, &b836, &b113, &b836, &b118, &b836, &b113, &b836, &b255}
	b557.items = []builder{&b554, &b556, &b836, &b256}
	var b566 = sequenceBuilder{id: 566, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b563 = sequenceBuilder{id: 563, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b558 = charBuilder{}
	var b559 = charBuilder{}
	var b560 = charBuilder{}
	var b561 = charBuilder{}
	var b562 = charBuilder{}
	b563.items = []builder{&b558, &b559, &b560, &b561, &b562}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b564.items = []builder{&b836, &b14}
	b565.items = []builder{&b836, &b14, &b564}
	b566.items = []builder{&b563, &b565, &b836, &b256}
	var b631 = choiceBuilder{id: 631, commit: 64, name: "assignment"}
	var b611 = sequenceBuilder{id: 611, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b608 = sequenceBuilder{id: 608, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b605 = charBuilder{}
	var b606 = charBuilder{}
	var b607 = charBuilder{}
	b608.items = []builder{&b605, &b606, &b607}
	var b610 = sequenceBuilder{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b609 = sequenceBuilder{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b609.items = []builder{&b836, &b14}
	b610.items = []builder{&b836, &b14, &b609}
	var b600 = sequenceBuilder{id: 600, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b595.items = []builder{&b836, &b14}
	b596.items = []builder{&b14, &b595}
	var b594 = sequenceBuilder{id: 594, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b593 = charBuilder{}
	b594.items = []builder{&b593}
	b597.items = []builder{&b596, &b836, &b594}
	var b599 = sequenceBuilder{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b598.items = []builder{&b836, &b14}
	b599.items = []builder{&b836, &b14, &b598}
	b600.items = []builder{&b266, &b836, &b597, &b599, &b836, &b395}
	b611.items = []builder{&b608, &b610, &b836, &b600}
	var b618 = sequenceBuilder{id: 618, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b614.items = []builder{&b836, &b14}
	b615.items = []builder{&b836, &b14, &b614}
	var b613 = sequenceBuilder{id: 613, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b612 = charBuilder{}
	b613.items = []builder{&b612}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b616.items = []builder{&b836, &b14}
	b617.items = []builder{&b836, &b14, &b616}
	b618.items = []builder{&b266, &b615, &b836, &b613, &b617, &b836, &b395}
	var b630 = sequenceBuilder{id: 630, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b622 = sequenceBuilder{id: 622, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b619 = charBuilder{}
	var b620 = charBuilder{}
	var b621 = charBuilder{}
	b622.items = []builder{&b619, &b620, &b621}
	var b629 = sequenceBuilder{id: 629, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b628 = sequenceBuilder{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b628.items = []builder{&b836, &b14}
	b629.items = []builder{&b836, &b14, &b628}
	var b624 = sequenceBuilder{id: 624, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b623 = charBuilder{}
	b624.items = []builder{&b623}
	var b625 = sequenceBuilder{id: 625, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b604 = sequenceBuilder{id: 604, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b603 = sequenceBuilder{id: 603, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b601 = sequenceBuilder{id: 601, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b601.items = []builder{&b113, &b836, &b600}
	var b602 = sequenceBuilder{id: 602, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b602.items = []builder{&b836, &b601}
	b603.items = []builder{&b836, &b601, &b602}
	b604.items = []builder{&b600, &b603}
	b625.items = []builder{&b113, &b836, &b604}
	var b627 = sequenceBuilder{id: 627, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b626 = charBuilder{}
	b627.items = []builder{&b626}
	b630.items = []builder{&b622, &b629, &b836, &b624, &b836, &b625, &b836, &b113, &b836, &b627}
	b631.options = []builder{&b611, &b618, &b630}
	var b803 = sequenceBuilder{id: 803, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b796 = sequenceBuilder{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b795 = charBuilder{}
	b796.items = []builder{&b795}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b799 = sequenceBuilder{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b799.items = []builder{&b836, &b14}
	b800.items = []builder{&b836, &b14, &b799}
	var b802 = sequenceBuilder{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b801.items = []builder{&b836, &b14}
	b802.items = []builder{&b836, &b14, &b801}
	var b798 = sequenceBuilder{id: 798, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b797 = charBuilder{}
	b798.items = []builder{&b797}
	b803.items = []builder{&b796, &b800, &b836, &b794, &b802, &b836, &b798}
	b794.options = []builder{&b499, &b557, &b566, &b631, &b803, &b395}
	var b190 = sequenceBuilder{id: 190, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b187 = sequenceBuilder{id: 187, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b186 = charBuilder{}
	b187.items = []builder{&b186}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	b190.items = []builder{&b187, &b836, &b818, &b836, &b822, &b836, &b818, &b836, &b189}
	b197.options = []builder{&b794, &b190}
	b200.items = []builder{&b192, &b836, &b113, &b836, &b194, &b836, &b113, &b836, &b196, &b199, &b836, &b197}
	b206.items = []builder{&b203, &b205, &b836, &b200}
	var b216 = sequenceBuilder{id: 216, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b209 = sequenceBuilder{id: 209, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b207 = charBuilder{}
	var b208 = charBuilder{}
	b209.items = []builder{&b207, &b208}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b212 = sequenceBuilder{id: 212, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b212.items = []builder{&b836, &b14}
	b213.items = []builder{&b836, &b14, &b212}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b210 = charBuilder{}
	b211.items = []builder{&b210}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b214.items = []builder{&b836, &b14}
	b215.items = []builder{&b836, &b14, &b214}
	b216.items = []builder{&b209, &b213, &b836, &b211, &b215, &b836, &b200}
	var b244 = choiceBuilder{id: 244, commit: 64, name: "expression-indexer"}
	var b234 = sequenceBuilder{id: 234, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b227 = sequenceBuilder{id: 227, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b226 = charBuilder{}
	b227.items = []builder{&b226}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b230 = sequenceBuilder{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b230.items = []builder{&b836, &b14}
	b231.items = []builder{&b836, &b14, &b230}
	var b233 = sequenceBuilder{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b232 = sequenceBuilder{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b232.items = []builder{&b836, &b14}
	b233.items = []builder{&b836, &b14, &b232}
	var b229 = sequenceBuilder{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b228 = charBuilder{}
	b229.items = []builder{&b228}
	b234.items = []builder{&b266, &b836, &b227, &b231, &b836, &b395, &b233, &b836, &b229}
	var b243 = sequenceBuilder{id: 243, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b236 = sequenceBuilder{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b235 = charBuilder{}
	b236.items = []builder{&b235}
	var b240 = sequenceBuilder{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b239 = sequenceBuilder{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b239.items = []builder{&b836, &b14}
	b240.items = []builder{&b836, &b14, &b239}
	var b225 = sequenceBuilder{id: 225, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b217.items = []builder{&b395}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b221 = sequenceBuilder{id: 221, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b221.items = []builder{&b836, &b14}
	b222.items = []builder{&b836, &b14, &b221}
	var b220 = sequenceBuilder{id: 220, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b219 = charBuilder{}
	b220.items = []builder{&b219}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b836, &b14}
	b224.items = []builder{&b836, &b14, &b223}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b395}
	b225.items = []builder{&b217, &b222, &b836, &b220, &b224, &b836, &b218}
	var b242 = sequenceBuilder{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b241 = sequenceBuilder{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b241.items = []builder{&b836, &b14}
	b242.items = []builder{&b836, &b14, &b241}
	var b238 = sequenceBuilder{id: 238, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b237 = charBuilder{}
	b238.items = []builder{&b237}
	b243.items = []builder{&b266, &b836, &b236, &b240, &b836, &b225, &b242, &b836, &b238}
	b244.options = []builder{&b234, &b243}
	var b251 = sequenceBuilder{id: 251, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b247.items = []builder{&b836, &b14}
	b248.items = []builder{&b836, &b14, &b247}
	var b246 = sequenceBuilder{id: 246, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b245 = charBuilder{}
	b246.items = []builder{&b245}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b249.items = []builder{&b836, &b14}
	b250.items = []builder{&b836, &b14, &b249}
	b251.items = []builder{&b266, &b248, &b836, &b246, &b250, &b836, &b103}
	var b265 = sequenceBuilder{id: 265, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b258 = sequenceBuilder{id: 258, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b257 = charBuilder{}
	b258.items = []builder{&b257}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b261 = sequenceBuilder{id: 261, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b261.items = []builder{&b836, &b14}
	b262.items = []builder{&b836, &b14, &b261}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b263.items = []builder{&b836, &b14}
	b264.items = []builder{&b836, &b14, &b263}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	b265.items = []builder{&b258, &b262, &b836, &b395, &b264, &b836, &b260}
	b266.options = []builder{&b60, &b73, &b86, &b98, &b510, &b103, &b124, &b129, &b158, &b163, &b206, &b216, &b244, &b251, &b256, &b265}
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
	b326.items = []builder{&b325, &b836, &b266}
	var b373 = choiceBuilder{id: 373, commit: 66}
	var b344 = sequenceBuilder{id: 344, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b332 = choiceBuilder{id: 332, commit: 66}
	b332.options = []builder{&b266, &b326}
	var b342 = sequenceBuilder{id: 342, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b339 = sequenceBuilder{id: 339, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b338 = sequenceBuilder{id: 338, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b338.items = []builder{&b836, &b14}
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
	b340.items = []builder{&b836, &b14}
	b341.items = []builder{&b836, &b14, &b340}
	b342.items = []builder{&b339, &b836, &b327, &b341, &b836, &b332}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b343.items = []builder{&b836, &b342}
	b344.items = []builder{&b332, &b836, &b342, &b343}
	var b351 = sequenceBuilder{id: 351, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b333 = choiceBuilder{id: 333, commit: 66}
	b333.options = []builder{&b332, &b344}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b345 = sequenceBuilder{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b345.items = []builder{&b836, &b14}
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
	b347.items = []builder{&b836, &b14}
	b348.items = []builder{&b836, &b14, &b347}
	b349.items = []builder{&b346, &b836, &b328, &b348, &b836, &b333}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b350.items = []builder{&b836, &b349}
	b351.items = []builder{&b333, &b836, &b349, &b350}
	var b358 = sequenceBuilder{id: 358, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b334 = choiceBuilder{id: 334, commit: 66}
	b334.options = []builder{&b333, &b351}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b352.items = []builder{&b836, &b14}
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
	b354.items = []builder{&b836, &b14}
	b355.items = []builder{&b836, &b14, &b354}
	b356.items = []builder{&b353, &b836, &b329, &b355, &b836, &b334}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b357.items = []builder{&b836, &b356}
	b358.items = []builder{&b334, &b836, &b356, &b357}
	var b365 = sequenceBuilder{id: 365, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b335 = choiceBuilder{id: 335, commit: 66}
	b335.options = []builder{&b334, &b358}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b836, &b14}
	b360.items = []builder{&b14, &b359}
	var b330 = sequenceBuilder{id: 330, commit: 66, ranges: [][]int{{1, 1}}}
	var b318 = sequenceBuilder{id: 318, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b316 = charBuilder{}
	var b317 = charBuilder{}
	b318.items = []builder{&b316, &b317}
	b330.items = []builder{&b318}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b361.items = []builder{&b836, &b14}
	b362.items = []builder{&b836, &b14, &b361}
	b363.items = []builder{&b360, &b836, &b330, &b362, &b836, &b335}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b364.items = []builder{&b836, &b363}
	b365.items = []builder{&b335, &b836, &b363, &b364}
	var b372 = sequenceBuilder{id: 372, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b336 = choiceBuilder{id: 336, commit: 66}
	b336.options = []builder{&b335, &b365}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b366.items = []builder{&b836, &b14}
	b367.items = []builder{&b14, &b366}
	var b331 = sequenceBuilder{id: 331, commit: 66, ranges: [][]int{{1, 1}}}
	var b321 = sequenceBuilder{id: 321, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b319 = charBuilder{}
	var b320 = charBuilder{}
	b321.items = []builder{&b319, &b320}
	b331.items = []builder{&b321}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b368.items = []builder{&b836, &b14}
	b369.items = []builder{&b836, &b14, &b368}
	b370.items = []builder{&b367, &b836, &b331, &b369, &b836, &b336}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b371.items = []builder{&b836, &b370}
	b372.items = []builder{&b336, &b836, &b370, &b371}
	b373.options = []builder{&b344, &b351, &b358, &b365, &b372}
	var b386 = sequenceBuilder{id: 386, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b379 = sequenceBuilder{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b378 = sequenceBuilder{id: 378, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b378.items = []builder{&b836, &b14}
	b379.items = []builder{&b836, &b14, &b378}
	var b375 = sequenceBuilder{id: 375, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b374 = charBuilder{}
	b375.items = []builder{&b374}
	var b381 = sequenceBuilder{id: 381, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b380 = sequenceBuilder{id: 380, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b380.items = []builder{&b836, &b14}
	b381.items = []builder{&b836, &b14, &b380}
	var b383 = sequenceBuilder{id: 383, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b382 = sequenceBuilder{id: 382, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b382.items = []builder{&b836, &b14}
	b383.items = []builder{&b836, &b14, &b382}
	var b377 = sequenceBuilder{id: 377, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b376 = charBuilder{}
	b377.items = []builder{&b376}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b384 = sequenceBuilder{id: 384, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b384.items = []builder{&b836, &b14}
	b385.items = []builder{&b836, &b14, &b384}
	b386.items = []builder{&b395, &b379, &b836, &b375, &b381, &b836, &b395, &b383, &b836, &b377, &b385, &b836, &b395}
	var b394 = sequenceBuilder{id: 394, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b387 = choiceBuilder{id: 387, commit: 66}
	b387.options = []builder{&b266, &b326, &b373, &b386}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b388.items = []builder{&b836, &b14}
	b389.items = []builder{&b14, &b388}
	var b324 = sequenceBuilder{id: 324, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b322 = charBuilder{}
	var b323 = charBuilder{}
	b324.items = []builder{&b322, &b323}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b390.items = []builder{&b836, &b14}
	b391.items = []builder{&b836, &b14, &b390}
	b392.items = []builder{&b389, &b836, &b324, &b391, &b836, &b387}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b393.items = []builder{&b836, &b392}
	b394.items = []builder{&b387, &b836, &b392, &b393}
	b395.options = []builder{&b266, &b326, &b373, &b386, &b394}
	b184.items = []builder{&b183, &b836, &b395}
	b185.items = []builder{&b181, &b836, &b184}
	var b432 = sequenceBuilder{id: 432, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b398 = sequenceBuilder{id: 398, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b396 = charBuilder{}
	var b397 = charBuilder{}
	b398.items = []builder{&b396, &b397}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b426.items = []builder{&b836, &b14}
	b427.items = []builder{&b836, &b14, &b426}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b428.items = []builder{&b836, &b14}
	b429.items = []builder{&b836, &b14, &b428}
	var b431 = sequenceBuilder{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b407 = sequenceBuilder{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b407.items = []builder{&b836, &b14}
	b408.items = []builder{&b14, &b407}
	var b403 = sequenceBuilder{id: 403, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b399 = charBuilder{}
	var b400 = charBuilder{}
	var b401 = charBuilder{}
	var b402 = charBuilder{}
	b403.items = []builder{&b399, &b400, &b401, &b402}
	var b410 = sequenceBuilder{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b409 = sequenceBuilder{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b409.items = []builder{&b836, &b14}
	b410.items = []builder{&b836, &b14, &b409}
	var b406 = sequenceBuilder{id: 406, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b404 = charBuilder{}
	var b405 = charBuilder{}
	b406.items = []builder{&b404, &b405}
	var b412 = sequenceBuilder{id: 412, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b411 = sequenceBuilder{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b411.items = []builder{&b836, &b14}
	b412.items = []builder{&b836, &b14, &b411}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b413 = sequenceBuilder{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b413.items = []builder{&b836, &b14}
	b414.items = []builder{&b836, &b14, &b413}
	b415.items = []builder{&b408, &b836, &b403, &b410, &b836, &b406, &b412, &b836, &b395, &b414, &b836, &b190}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b430.items = []builder{&b836, &b415}
	b431.items = []builder{&b836, &b415, &b430}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b421 = sequenceBuilder{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b421.items = []builder{&b836, &b14}
	b422.items = []builder{&b14, &b421}
	var b420 = sequenceBuilder{id: 420, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b416 = charBuilder{}
	var b417 = charBuilder{}
	var b418 = charBuilder{}
	var b419 = charBuilder{}
	b420.items = []builder{&b416, &b417, &b418, &b419}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b423.items = []builder{&b836, &b14}
	b424.items = []builder{&b836, &b14, &b423}
	b425.items = []builder{&b422, &b836, &b420, &b424, &b836, &b190}
	b432.items = []builder{&b398, &b427, &b836, &b395, &b429, &b836, &b190, &b431, &b836, &b425}
	var b489 = sequenceBuilder{id: 489, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b474 = sequenceBuilder{id: 474, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b468 = charBuilder{}
	var b469 = charBuilder{}
	var b470 = charBuilder{}
	var b471 = charBuilder{}
	var b472 = charBuilder{}
	var b473 = charBuilder{}
	b474.items = []builder{&b468, &b469, &b470, &b471, &b472, &b473}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b485.items = []builder{&b836, &b14}
	b486.items = []builder{&b836, &b14, &b485}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b487.items = []builder{&b836, &b14}
	b488.items = []builder{&b836, &b14, &b487}
	var b476 = sequenceBuilder{id: 476, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b475 = charBuilder{}
	b476.items = []builder{&b475}
	var b482 = sequenceBuilder{id: 482, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b477 = choiceBuilder{id: 477, commit: 2}
	var b467 = sequenceBuilder{id: 467, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b462 = sequenceBuilder{id: 462, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b455 = sequenceBuilder{id: 455, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	var b454 = charBuilder{}
	b455.items = []builder{&b451, &b452, &b453, &b454}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b458 = sequenceBuilder{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b458.items = []builder{&b836, &b14}
	b459.items = []builder{&b836, &b14, &b458}
	var b461 = sequenceBuilder{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b460.items = []builder{&b836, &b14}
	b461.items = []builder{&b836, &b14, &b460}
	var b457 = sequenceBuilder{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b456 = charBuilder{}
	b457.items = []builder{&b456}
	b462.items = []builder{&b455, &b459, &b836, &b395, &b461, &b836, &b457}
	var b466 = sequenceBuilder{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b464 = sequenceBuilder{id: 464, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b463 = charBuilder{}
	b464.items = []builder{&b463}
	var b465 = sequenceBuilder{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b465.items = []builder{&b836, &b464}
	b466.items = []builder{&b836, &b464, &b465}
	b467.items = []builder{&b462, &b466, &b836, &b804}
	var b450 = sequenceBuilder{id: 450, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b445 = sequenceBuilder{id: 445, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b440 = sequenceBuilder{id: 440, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b433 = charBuilder{}
	var b434 = charBuilder{}
	var b435 = charBuilder{}
	var b436 = charBuilder{}
	var b437 = charBuilder{}
	var b438 = charBuilder{}
	var b439 = charBuilder{}
	b440.items = []builder{&b433, &b434, &b435, &b436, &b437, &b438, &b439}
	var b444 = sequenceBuilder{id: 444, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b443 = sequenceBuilder{id: 443, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b443.items = []builder{&b836, &b14}
	b444.items = []builder{&b836, &b14, &b443}
	var b442 = sequenceBuilder{id: 442, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b441 = charBuilder{}
	b442.items = []builder{&b441}
	b445.items = []builder{&b440, &b444, &b836, &b442}
	var b449 = sequenceBuilder{id: 449, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b447 = sequenceBuilder{id: 447, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b446 = charBuilder{}
	b447.items = []builder{&b446}
	var b448 = sequenceBuilder{id: 448, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b448.items = []builder{&b836, &b447}
	b449.items = []builder{&b836, &b447, &b448}
	b450.items = []builder{&b445, &b449, &b836, &b804}
	b477.options = []builder{&b467, &b450}
	var b481 = sequenceBuilder{id: 481, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b479 = sequenceBuilder{id: 479, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b478 = choiceBuilder{id: 478, commit: 2}
	b478.options = []builder{&b467, &b450, &b804}
	b479.items = []builder{&b818, &b836, &b478}
	var b480 = sequenceBuilder{id: 480, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b480.items = []builder{&b836, &b479}
	b481.items = []builder{&b836, &b479, &b480}
	b482.items = []builder{&b477, &b481}
	var b484 = sequenceBuilder{id: 484, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b483 = charBuilder{}
	b484.items = []builder{&b483}
	b489.items = []builder{&b474, &b486, &b836, &b395, &b488, &b836, &b476, &b836, &b818, &b836, &b482, &b836, &b818, &b836, &b484}
	var b551 = sequenceBuilder{id: 551, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b538 = sequenceBuilder{id: 538, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b532 = charBuilder{}
	var b533 = charBuilder{}
	var b534 = charBuilder{}
	var b535 = charBuilder{}
	var b536 = charBuilder{}
	var b537 = charBuilder{}
	b538.items = []builder{&b532, &b533, &b534, &b535, &b536, &b537}
	var b550 = sequenceBuilder{id: 550, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b549 = sequenceBuilder{id: 549, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b549.items = []builder{&b836, &b14}
	b550.items = []builder{&b836, &b14, &b549}
	var b540 = sequenceBuilder{id: 540, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b539 = charBuilder{}
	b540.items = []builder{&b539}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b541 = choiceBuilder{id: 541, commit: 2}
	var b531 = sequenceBuilder{id: 531, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b526 = sequenceBuilder{id: 526, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b519 = sequenceBuilder{id: 519, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b515 = charBuilder{}
	var b516 = charBuilder{}
	var b517 = charBuilder{}
	var b518 = charBuilder{}
	b519.items = []builder{&b515, &b516, &b517, &b518}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b522 = sequenceBuilder{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b522.items = []builder{&b836, &b14}
	b523.items = []builder{&b836, &b14, &b522}
	var b514 = choiceBuilder{id: 514, commit: 66}
	var b513 = sequenceBuilder{id: 513, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b512 = sequenceBuilder{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b511 = sequenceBuilder{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b511.items = []builder{&b836, &b14}
	b512.items = []builder{&b836, &b14, &b511}
	b513.items = []builder{&b103, &b512, &b836, &b510}
	b514.options = []builder{&b499, &b510, &b513}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b524.items = []builder{&b836, &b14}
	b525.items = []builder{&b836, &b14, &b524}
	var b521 = sequenceBuilder{id: 521, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b520 = charBuilder{}
	b521.items = []builder{&b520}
	b526.items = []builder{&b519, &b523, &b836, &b514, &b525, &b836, &b521}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b528 = sequenceBuilder{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b527 = charBuilder{}
	b528.items = []builder{&b527}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b529.items = []builder{&b836, &b528}
	b530.items = []builder{&b836, &b528, &b529}
	b531.items = []builder{&b526, &b530, &b836, &b804}
	b541.options = []builder{&b531, &b450}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b543 = sequenceBuilder{id: 543, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b542 = choiceBuilder{id: 542, commit: 2}
	b542.options = []builder{&b531, &b450, &b804}
	b543.items = []builder{&b818, &b836, &b542}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b544.items = []builder{&b836, &b543}
	b545.items = []builder{&b836, &b543, &b544}
	b546.items = []builder{&b541, &b545}
	var b548 = sequenceBuilder{id: 548, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b547 = charBuilder{}
	b548.items = []builder{&b547}
	b551.items = []builder{&b538, &b550, &b836, &b540, &b836, &b818, &b836, &b546, &b836, &b818, &b836, &b548}
	var b592 = sequenceBuilder{id: 592, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b581 = sequenceBuilder{id: 581, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b578 = charBuilder{}
	var b579 = charBuilder{}
	var b580 = charBuilder{}
	b581.items = []builder{&b578, &b579, &b580}
	var b591 = choiceBuilder{id: 591, commit: 2}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b584 = sequenceBuilder{id: 584, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b583 = sequenceBuilder{id: 583, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b582 = sequenceBuilder{id: 582, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b582.items = []builder{&b836, &b14}
	b583.items = []builder{&b14, &b582}
	var b577 = choiceBuilder{id: 577, commit: 66}
	var b576 = choiceBuilder{id: 576, commit: 64, name: "range-over-expression"}
	var b575 = sequenceBuilder{id: 575, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b571 = sequenceBuilder{id: 571, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b571.items = []builder{&b836, &b14}
	b572.items = []builder{&b836, &b14, &b571}
	var b569 = sequenceBuilder{id: 569, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b567 = charBuilder{}
	var b568 = charBuilder{}
	b569.items = []builder{&b567, &b568}
	var b574 = sequenceBuilder{id: 574, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b573.items = []builder{&b836, &b14}
	b574.items = []builder{&b836, &b14, &b573}
	var b570 = choiceBuilder{id: 570, commit: 2}
	b570.options = []builder{&b395, &b225}
	b575.items = []builder{&b103, &b572, &b836, &b569, &b574, &b836, &b570}
	b576.options = []builder{&b575, &b225}
	b577.options = []builder{&b395, &b576}
	b584.items = []builder{&b583, &b836, &b577}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b585.items = []builder{&b836, &b14}
	b586.items = []builder{&b836, &b14, &b585}
	b587.items = []builder{&b584, &b586, &b836, &b190}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b588.items = []builder{&b836, &b14}
	b589.items = []builder{&b14, &b588}
	b590.items = []builder{&b589, &b836, &b190}
	b591.options = []builder{&b587, &b590}
	b592.items = []builder{&b581, &b836, &b591}
	var b740 = choiceBuilder{id: 740, commit: 66}
	var b653 = sequenceBuilder{id: 653, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b649 = sequenceBuilder{id: 649, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b646 = charBuilder{}
	var b647 = charBuilder{}
	var b648 = charBuilder{}
	b649.items = []builder{&b646, &b647, &b648}
	var b652 = sequenceBuilder{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b651 = sequenceBuilder{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b651.items = []builder{&b836, &b14}
	b652.items = []builder{&b836, &b14, &b651}
	var b650 = choiceBuilder{id: 650, commit: 2}
	var b640 = sequenceBuilder{id: 640, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b639 = sequenceBuilder{id: 639, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b634.items = []builder{&b836, &b14}
	b635.items = []builder{&b14, &b634}
	var b633 = sequenceBuilder{id: 633, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b632 = charBuilder{}
	b633.items = []builder{&b632}
	b636.items = []builder{&b635, &b836, &b633}
	var b638 = sequenceBuilder{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b637 = sequenceBuilder{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b637.items = []builder{&b836, &b14}
	b638.items = []builder{&b836, &b14, &b637}
	b639.items = []builder{&b103, &b836, &b636, &b638, &b836, &b395}
	b640.items = []builder{&b639}
	var b645 = sequenceBuilder{id: 645, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b642 = sequenceBuilder{id: 642, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b641 = charBuilder{}
	b642.items = []builder{&b641}
	var b644 = sequenceBuilder{id: 644, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b643 = sequenceBuilder{id: 643, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b643.items = []builder{&b836, &b14}
	b644.items = []builder{&b836, &b14, &b643}
	b645.items = []builder{&b642, &b644, &b836, &b639}
	b650.options = []builder{&b640, &b645}
	b653.items = []builder{&b649, &b652, &b836, &b650}
	var b674 = sequenceBuilder{id: 674, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b667 = sequenceBuilder{id: 667, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b664 = charBuilder{}
	var b665 = charBuilder{}
	var b666 = charBuilder{}
	b667.items = []builder{&b664, &b665, &b666}
	var b673 = sequenceBuilder{id: 673, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b672 = sequenceBuilder{id: 672, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b672.items = []builder{&b836, &b14}
	b673.items = []builder{&b836, &b14, &b672}
	var b669 = sequenceBuilder{id: 669, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b668 = charBuilder{}
	b669.items = []builder{&b668}
	var b659 = sequenceBuilder{id: 659, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b654 = choiceBuilder{id: 654, commit: 2}
	b654.options = []builder{&b640, &b645}
	var b658 = sequenceBuilder{id: 658, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b656 = sequenceBuilder{id: 656, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b655 = choiceBuilder{id: 655, commit: 2}
	b655.options = []builder{&b640, &b645}
	b656.items = []builder{&b113, &b836, &b655}
	var b657 = sequenceBuilder{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b657.items = []builder{&b836, &b656}
	b658.items = []builder{&b836, &b656, &b657}
	b659.items = []builder{&b654, &b658}
	var b671 = sequenceBuilder{id: 671, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b670 = charBuilder{}
	b671.items = []builder{&b670}
	b674.items = []builder{&b667, &b673, &b836, &b669, &b836, &b113, &b836, &b659, &b836, &b113, &b836, &b671}
	var b689 = sequenceBuilder{id: 689, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b678 = sequenceBuilder{id: 678, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b675 = charBuilder{}
	var b676 = charBuilder{}
	var b677 = charBuilder{}
	b678.items = []builder{&b675, &b676, &b677}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b685 = sequenceBuilder{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b685.items = []builder{&b836, &b14}
	b686.items = []builder{&b836, &b14, &b685}
	var b680 = sequenceBuilder{id: 680, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b679 = charBuilder{}
	b680.items = []builder{&b679}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b687.items = []builder{&b836, &b14}
	b688.items = []builder{&b836, &b14, &b687}
	var b682 = sequenceBuilder{id: 682, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b681 = charBuilder{}
	b682.items = []builder{&b681}
	var b663 = sequenceBuilder{id: 663, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b662 = sequenceBuilder{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b660 = sequenceBuilder{id: 660, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b660.items = []builder{&b113, &b836, &b640}
	var b661 = sequenceBuilder{id: 661, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b661.items = []builder{&b836, &b660}
	b662.items = []builder{&b836, &b660, &b661}
	b663.items = []builder{&b640, &b662}
	var b684 = sequenceBuilder{id: 684, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b683 = charBuilder{}
	b684.items = []builder{&b683}
	b689.items = []builder{&b678, &b686, &b836, &b680, &b688, &b836, &b682, &b836, &b113, &b836, &b663, &b836, &b113, &b836, &b684}
	var b705 = sequenceBuilder{id: 705, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b701 = sequenceBuilder{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b699 = charBuilder{}
	var b700 = charBuilder{}
	b701.items = []builder{&b699, &b700}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b703.items = []builder{&b836, &b14}
	b704.items = []builder{&b836, &b14, &b703}
	var b702 = choiceBuilder{id: 702, commit: 2}
	var b693 = sequenceBuilder{id: 693, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b692 = sequenceBuilder{id: 692, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b691 = sequenceBuilder{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b690 = sequenceBuilder{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b690.items = []builder{&b836, &b14}
	b691.items = []builder{&b836, &b14, &b690}
	b692.items = []builder{&b103, &b691, &b836, &b200}
	b693.items = []builder{&b692}
	var b698 = sequenceBuilder{id: 698, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b695 = sequenceBuilder{id: 695, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b694 = charBuilder{}
	b695.items = []builder{&b694}
	var b697 = sequenceBuilder{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b696.items = []builder{&b836, &b14}
	b697.items = []builder{&b836, &b14, &b696}
	b698.items = []builder{&b695, &b697, &b836, &b692}
	b702.options = []builder{&b693, &b698}
	b705.items = []builder{&b701, &b704, &b836, &b702}
	var b725 = sequenceBuilder{id: 725, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b718 = sequenceBuilder{id: 718, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b716 = charBuilder{}
	var b717 = charBuilder{}
	b718.items = []builder{&b716, &b717}
	var b724 = sequenceBuilder{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b723 = sequenceBuilder{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b723.items = []builder{&b836, &b14}
	b724.items = []builder{&b836, &b14, &b723}
	var b720 = sequenceBuilder{id: 720, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b719 = charBuilder{}
	b720.items = []builder{&b719}
	var b715 = sequenceBuilder{id: 715, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b710 = choiceBuilder{id: 710, commit: 2}
	b710.options = []builder{&b693, &b698}
	var b714 = sequenceBuilder{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b712 = sequenceBuilder{id: 712, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b711 = choiceBuilder{id: 711, commit: 2}
	b711.options = []builder{&b693, &b698}
	b712.items = []builder{&b113, &b836, &b711}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b713.items = []builder{&b836, &b712}
	b714.items = []builder{&b836, &b712, &b713}
	b715.items = []builder{&b710, &b714}
	var b722 = sequenceBuilder{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b721 = charBuilder{}
	b722.items = []builder{&b721}
	b725.items = []builder{&b718, &b724, &b836, &b720, &b836, &b113, &b836, &b715, &b836, &b113, &b836, &b722}
	var b739 = sequenceBuilder{id: 739, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b728 = sequenceBuilder{id: 728, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b726 = charBuilder{}
	var b727 = charBuilder{}
	b728.items = []builder{&b726, &b727}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b735 = sequenceBuilder{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b735.items = []builder{&b836, &b14}
	b736.items = []builder{&b836, &b14, &b735}
	var b730 = sequenceBuilder{id: 730, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b729 = charBuilder{}
	b730.items = []builder{&b729}
	var b738 = sequenceBuilder{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b737.items = []builder{&b836, &b14}
	b738.items = []builder{&b836, &b14, &b737}
	var b732 = sequenceBuilder{id: 732, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b731 = charBuilder{}
	b732.items = []builder{&b731}
	var b709 = sequenceBuilder{id: 709, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b708 = sequenceBuilder{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b706 = sequenceBuilder{id: 706, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b706.items = []builder{&b113, &b836, &b693}
	var b707 = sequenceBuilder{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b707.items = []builder{&b836, &b706}
	b708.items = []builder{&b836, &b706, &b707}
	b709.items = []builder{&b693, &b708}
	var b734 = sequenceBuilder{id: 734, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b733 = charBuilder{}
	b734.items = []builder{&b733}
	b739.items = []builder{&b728, &b736, &b836, &b730, &b738, &b836, &b732, &b836, &b113, &b836, &b709, &b836, &b113, &b836, &b734}
	b740.options = []builder{&b653, &b674, &b689, &b705, &b725, &b739}
	var b783 = choiceBuilder{id: 783, commit: 64, name: "require"}
	var b767 = sequenceBuilder{id: 767, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b764 = sequenceBuilder{id: 764, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b757 = charBuilder{}
	var b758 = charBuilder{}
	var b759 = charBuilder{}
	var b760 = charBuilder{}
	var b761 = charBuilder{}
	var b762 = charBuilder{}
	var b763 = charBuilder{}
	b764.items = []builder{&b757, &b758, &b759, &b760, &b761, &b762, &b763}
	var b766 = sequenceBuilder{id: 766, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b765 = sequenceBuilder{id: 765, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b765.items = []builder{&b836, &b14}
	b766.items = []builder{&b836, &b14, &b765}
	var b752 = choiceBuilder{id: 752, commit: 64, name: "require-fact"}
	var b751 = sequenceBuilder{id: 751, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b743 = choiceBuilder{id: 743, commit: 2}
	var b742 = sequenceBuilder{id: 742, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b741 = charBuilder{}
	b742.items = []builder{&b741}
	b743.options = []builder{&b103, &b742}
	var b748 = sequenceBuilder{id: 748, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b747 = sequenceBuilder{id: 747, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b746 = sequenceBuilder{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b746.items = []builder{&b836, &b14}
	b747.items = []builder{&b14, &b746}
	var b745 = sequenceBuilder{id: 745, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b744 = charBuilder{}
	b745.items = []builder{&b744}
	b748.items = []builder{&b747, &b836, &b745}
	var b750 = sequenceBuilder{id: 750, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b749 = sequenceBuilder{id: 749, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b749.items = []builder{&b836, &b14}
	b750.items = []builder{&b836, &b14, &b749}
	b751.items = []builder{&b743, &b836, &b748, &b750, &b836, &b86}
	b752.options = []builder{&b86, &b751}
	b767.items = []builder{&b764, &b766, &b836, &b752}
	var b782 = sequenceBuilder{id: 782, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b775 = sequenceBuilder{id: 775, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b768 = charBuilder{}
	var b769 = charBuilder{}
	var b770 = charBuilder{}
	var b771 = charBuilder{}
	var b772 = charBuilder{}
	var b773 = charBuilder{}
	var b774 = charBuilder{}
	b775.items = []builder{&b768, &b769, &b770, &b771, &b772, &b773, &b774}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b780 = sequenceBuilder{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b780.items = []builder{&b836, &b14}
	b781.items = []builder{&b836, &b14, &b780}
	var b777 = sequenceBuilder{id: 777, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b776 = charBuilder{}
	b777.items = []builder{&b776}
	var b756 = sequenceBuilder{id: 756, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b753.items = []builder{&b113, &b836, &b752}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b754.items = []builder{&b836, &b753}
	b755.items = []builder{&b836, &b753, &b754}
	b756.items = []builder{&b752, &b755}
	var b779 = sequenceBuilder{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b778 = charBuilder{}
	b779.items = []builder{&b778}
	b782.items = []builder{&b775, &b781, &b836, &b777, &b836, &b113, &b836, &b756, &b836, &b113, &b836, &b779}
	b783.options = []builder{&b767, &b782}
	var b793 = sequenceBuilder{id: 793, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b790 = sequenceBuilder{id: 790, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b784 = charBuilder{}
	var b785 = charBuilder{}
	var b786 = charBuilder{}
	var b787 = charBuilder{}
	var b788 = charBuilder{}
	var b789 = charBuilder{}
	b790.items = []builder{&b784, &b785, &b786, &b787, &b788, &b789}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b791.items = []builder{&b836, &b14}
	b792.items = []builder{&b836, &b14, &b791}
	b793.items = []builder{&b790, &b792, &b836, &b740}
	var b813 = sequenceBuilder{id: 813, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b806 = sequenceBuilder{id: 806, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b805 = charBuilder{}
	b806.items = []builder{&b805}
	var b810 = sequenceBuilder{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b809 = sequenceBuilder{id: 809, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b809.items = []builder{&b836, &b14}
	b810.items = []builder{&b836, &b14, &b809}
	var b812 = sequenceBuilder{id: 812, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b811 = sequenceBuilder{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b811.items = []builder{&b836, &b14}
	b812.items = []builder{&b836, &b14, &b811}
	var b808 = sequenceBuilder{id: 808, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b807 = charBuilder{}
	b808.items = []builder{&b807}
	b813.items = []builder{&b806, &b810, &b836, &b804, &b812, &b836, &b808}
	b804.options = []builder{&b185, &b432, &b489, &b551, &b592, &b740, &b783, &b793, &b813, &b794}
	var b821 = sequenceBuilder{id: 821, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b819 = sequenceBuilder{id: 819, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b819.items = []builder{&b818, &b836, &b804}
	var b820 = sequenceBuilder{id: 820, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b820.items = []builder{&b836, &b819}
	b821.items = []builder{&b836, &b819, &b820}
	b822.items = []builder{&b804, &b821}
	b837.items = []builder{&b833, &b836, &b818, &b836, &b822, &b836, &b818}
	b838.items = []builder{&b836, &b837, &b836}

	return parseInput(r, &p838, &b838)
}
