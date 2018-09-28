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
	var p818 = sequenceParser{id: 818, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p816 = choiceParser{id: 816, commit: 2}
	var p814 = choiceParser{id: 814, commit: 70, name: "ws", generalizations: []int{816}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{814, 816}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{814, 816}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{814, 816}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{814, 816}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{814, 816}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{814, 816}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p814.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p815 = sequenceParser{id: 815, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{816}}
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
	p41.items = []parser{&p40, &p816, &p38}
	var p42 = sequenceParser{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p42.items = []parser{&p816, &p41}
	p43.items = []parser{&p816, &p41, &p42}
	p44.items = []parser{&p38, &p43}
	p815.items = []parser{&p44}
	p816.options = []parser{&p814, &p815}
	var p817 = sequenceParser{id: 817, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p813 = sequenceParser{id: 813, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p810 = sequenceParser{id: 810, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p808 = charParser{id: 808, chars: []rune{35}}
	var p809 = charParser{id: 809, chars: []rune{33}}
	p810.items = []parser{&p808, &p809}
	var p807 = sequenceParser{id: 807, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p806 = sequenceParser{id: 806, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p804 = sequenceParser{id: 804, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p803 = charParser{id: 803, not: true, chars: []rune{10}}
	p804.items = []parser{&p803}
	var p805 = sequenceParser{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p805.items = []parser{&p816, &p804}
	p806.items = []parser{&p804, &p805}
	p807.items = []parser{&p806}
	var p812 = sequenceParser{id: 812, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p811 = charParser{id: 811, chars: []rune{10}}
	p812.items = []parser{&p811}
	p813.items = []parser{&p810, &p816, &p807, &p816, &p812}
	var p798 = sequenceParser{id: 798, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p796 = choiceParser{id: 796, commit: 2}
	var p795 = sequenceParser{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{796}}
	var p794 = charParser{id: 794, chars: []rune{59}}
	p795.items = []parser{&p794}
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{796, 113}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p796.options = []parser{&p795, &p14}
	var p797 = sequenceParser{id: 797, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p797.items = []parser{&p816, &p796}
	p798.items = []parser{&p796, &p797}
	var p802 = sequenceParser{id: 802, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p784 = choiceParser{id: 784, commit: 66, name: "statement", generalizations: []int{458, 522}}
	var p187 = sequenceParser{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{784, 458, 522}}
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
	p184.items = []parser{&p816, &p14}
	p185.items = []parser{&p14, &p184}
	var p375 = choiceParser{id: 375, commit: 66, name: "expression", generalizations: []int{116, 774, 198, 557, 550, 784}}
	var p267 = choiceParser{id: 267, commit: 66, name: "primary-expression", generalizations: []int{116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p62 = choiceParser{id: 62, commit: 64, name: "int", generalizations: []int{267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p53 = sequenceParser{id: 53, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p52 = sequenceParser{id: 52, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p51 = charParser{id: 51, ranges: [][]rune{{49, 57}}}
	p52.items = []parser{&p51}
	var p46 = sequenceParser{id: 46, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 57}}}
	p46.items = []parser{&p45}
	p53.items = []parser{&p52, &p46}
	var p56 = sequenceParser{id: 56, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p55 = sequenceParser{id: 55, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p54 = charParser{id: 54, chars: []rune{48}}
	p55.items = []parser{&p54}
	var p48 = sequenceParser{id: 48, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p47 = charParser{id: 47, ranges: [][]rune{{48, 55}}}
	p48.items = []parser{&p47}
	p56.items = []parser{&p55, &p48}
	var p61 = sequenceParser{id: 61, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{62, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
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
	var p75 = choiceParser{id: 75, commit: 72, name: "float", generalizations: []int{267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p70 = sequenceParser{id: 70, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{75, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
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
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{75, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p72 = sequenceParser{id: 72, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p71 = charParser{id: 71, chars: []rune{46}}
	p72.items = []parser{&p71}
	p73.items = []parser{&p72, &p46, &p67}
	var p74 = sequenceParser{id: 74, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{75, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	p74.items = []parser{&p46, &p67}
	p75.options = []parser{&p70, &p73, &p74}
	var p88 = sequenceParser{id: 88, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 116, 141, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 732, 784}}
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
	var p100 = choiceParser{id: 100, commit: 66, name: "bool", generalizations: []int{267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p93 = sequenceParser{id: 93, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p89 = charParser{id: 89, chars: []rune{116}}
	var p90 = charParser{id: 90, chars: []rune{114}}
	var p91 = charParser{id: 91, chars: []rune{117}}
	var p92 = charParser{id: 92, chars: []rune{101}}
	p93.items = []parser{&p89, &p90, &p91, &p92}
	var p99 = sequenceParser{id: 99, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p94 = charParser{id: 94, chars: []rune{102}}
	var p95 = charParser{id: 95, chars: []rune{97}}
	var p96 = charParser{id: 96, chars: []rune{108}}
	var p97 = charParser{id: 97, chars: []rune{115}}
	var p98 = charParser{id: 98, chars: []rune{101}}
	p99.items = []parser{&p94, &p95, &p96, &p97, &p98}
	p100.options = []parser{&p93, &p99}
	var p480 = sequenceParser{id: 480, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 116, 774, 198, 375, 333, 334, 335, 336, 337, 338, 494, 557, 550, 784}}
	var p477 = sequenceParser{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p470 = charParser{id: 470, chars: []rune{114}}
	var p471 = charParser{id: 471, chars: []rune{101}}
	var p472 = charParser{id: 472, chars: []rune{99}}
	var p473 = charParser{id: 473, chars: []rune{101}}
	var p474 = charParser{id: 474, chars: []rune{105}}
	var p475 = charParser{id: 475, chars: []rune{118}}
	var p476 = charParser{id: 476, chars: []rune{101}}
	p477.items = []parser{&p470, &p471, &p472, &p473, &p474, &p475, &p476}
	var p479 = sequenceParser{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p478 = sequenceParser{id: 478, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p478.items = []parser{&p816, &p14}
	p479.items = []parser{&p816, &p14, &p478}
	p480.items = []parser{&p477, &p479, &p816, &p267}
	var p105 = sequenceParser{id: 105, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{267, 116, 141, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 723, 784}}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p102.items = []parser{&p101}
	var p104 = sequenceParser{id: 104, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p103 = charParser{id: 103, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p104.items = []parser{&p103}
	p105.items = []parser{&p102, &p104}
	var p126 = sequenceParser{id: 126, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{116, 267, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
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
	p114.items = []parser{&p816, &p113}
	p115.items = []parser{&p113, &p114}
	var p120 = sequenceParser{id: 120, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p116 = choiceParser{id: 116, commit: 66, name: "list-item"}
	var p110 = sequenceParser{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{116, 149, 150}}
	var p109 = sequenceParser{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	var p108 = charParser{id: 108, chars: []rune{46}}
	p109.items = []parser{&p106, &p107, &p108}
	p110.items = []parser{&p267, &p816, &p109}
	p116.options = []parser{&p375, &p110}
	var p119 = sequenceParser{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p117.items = []parser{&p115, &p816, &p116}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p118.items = []parser{&p816, &p117}
	p119.items = []parser{&p816, &p117, &p118}
	p120.items = []parser{&p116, &p119}
	var p124 = sequenceParser{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p123 = charParser{id: 123, chars: []rune{93}}
	p124.items = []parser{&p123}
	p125.items = []parser{&p122, &p816, &p115, &p816, &p120, &p816, &p115, &p816, &p124}
	p126.items = []parser{&p125}
	var p131 = sequenceParser{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p128 = sequenceParser{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p127 = charParser{id: 127, chars: []rune{126}}
	p128.items = []parser{&p127}
	var p130 = sequenceParser{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p129.items = []parser{&p816, &p14}
	p130.items = []parser{&p816, &p14, &p129}
	p131.items = []parser{&p128, &p130, &p816, &p125}
	var p160 = sequenceParser{id: 160, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{267, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
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
	p136.items = []parser{&p816, &p14}
	p137.items = []parser{&p816, &p14, &p136}
	var p139 = sequenceParser{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p138.items = []parser{&p816, &p14}
	p139.items = []parser{&p816, &p14, &p138}
	var p135 = sequenceParser{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p134 = charParser{id: 134, chars: []rune{93}}
	p135.items = []parser{&p134}
	p140.items = []parser{&p133, &p137, &p816, &p375, &p139, &p816, &p135}
	p141.options = []parser{&p105, &p88, &p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p816, &p14}
	p145.items = []parser{&p816, &p14, &p144}
	var p143 = sequenceParser{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p142 = charParser{id: 142, chars: []rune{58}}
	p143.items = []parser{&p142}
	var p147 = sequenceParser{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p146.items = []parser{&p816, &p14}
	p147.items = []parser{&p816, &p14, &p146}
	p148.items = []parser{&p141, &p145, &p816, &p143, &p147, &p816, &p375}
	p149.options = []parser{&p148, &p110}
	var p153 = sequenceParser{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p150 = choiceParser{id: 150, commit: 2}
	p150.options = []parser{&p148, &p110}
	p151.items = []parser{&p115, &p816, &p150}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p152.items = []parser{&p816, &p151}
	p153.items = []parser{&p816, &p151, &p152}
	p154.items = []parser{&p149, &p153}
	var p158 = sequenceParser{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p157 = charParser{id: 157, chars: []rune{125}}
	p158.items = []parser{&p157}
	p159.items = []parser{&p156, &p816, &p115, &p816, &p154, &p816, &p115, &p816, &p158}
	p160.items = []parser{&p159}
	var p165 = sequenceParser{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 774, 198, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p162 = sequenceParser{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p161 = charParser{id: 161, chars: []rune{126}}
	p162.items = []parser{&p161}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p163.items = []parser{&p816, &p14}
	p164.items = []parser{&p816, &p14, &p163}
	p165.items = []parser{&p162, &p164, &p816, &p159}
	var p207 = sequenceParser{id: 207, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{774, 198, 267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p204 = sequenceParser{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p202 = charParser{id: 202, chars: []rune{102}}
	var p203 = charParser{id: 203, chars: []rune{110}}
	p204.items = []parser{&p202, &p203}
	var p206 = sequenceParser{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p205.items = []parser{&p816, &p14}
	p206.items = []parser{&p816, &p14, &p205}
	var p201 = sequenceParser{id: 201, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p194 = sequenceParser{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p193 = charParser{id: 193, chars: []rune{40}}
	p194.items = []parser{&p193}
	var p169 = sequenceParser{id: 169, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p168 = sequenceParser{id: 168, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p166.items = []parser{&p115, &p816, &p105}
	var p167 = sequenceParser{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p167.items = []parser{&p816, &p166}
	p168.items = []parser{&p816, &p166, &p167}
	p169.items = []parser{&p105, &p168}
	var p195 = sequenceParser{id: 195, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p176 = sequenceParser{id: 176, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p173 = sequenceParser{id: 173, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p170 = charParser{id: 170, chars: []rune{46}}
	var p171 = charParser{id: 171, chars: []rune{46}}
	var p172 = charParser{id: 172, chars: []rune{46}}
	p173.items = []parser{&p170, &p171, &p172}
	var p175 = sequenceParser{id: 175, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p174 = sequenceParser{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p174.items = []parser{&p816, &p14}
	p175.items = []parser{&p816, &p14, &p174}
	p176.items = []parser{&p173, &p175, &p816, &p105}
	p195.items = []parser{&p115, &p816, &p176}
	var p197 = sequenceParser{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p196 = charParser{id: 196, chars: []rune{41}}
	p197.items = []parser{&p196}
	var p200 = sequenceParser{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p199.items = []parser{&p816, &p14}
	p200.items = []parser{&p816, &p14, &p199}
	var p198 = choiceParser{id: 198, commit: 2}
	var p774 = choiceParser{id: 774, commit: 66, name: "simple-statement", generalizations: []int{198, 784}}
	var p490 = sequenceParser{id: 490, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{774, 198, 494, 784}}
	var p485 = sequenceParser{id: 485, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p481 = charParser{id: 481, chars: []rune{115}}
	var p482 = charParser{id: 482, chars: []rune{101}}
	var p483 = charParser{id: 483, chars: []rune{110}}
	var p484 = charParser{id: 484, chars: []rune{100}}
	p485.items = []parser{&p481, &p482, &p483, &p484}
	var p487 = sequenceParser{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p486.items = []parser{&p816, &p14}
	p487.items = []parser{&p816, &p14, &p486}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p488.items = []parser{&p816, &p14}
	p489.items = []parser{&p816, &p14, &p488}
	p490.items = []parser{&p485, &p487, &p816, &p267, &p489, &p816, &p267}
	var p537 = sequenceParser{id: 537, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{774, 198, 784}}
	var p534 = sequenceParser{id: 534, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p532 = charParser{id: 532, chars: []rune{103}}
	var p533 = charParser{id: 533, chars: []rune{111}}
	p534.items = []parser{&p532, &p533}
	var p536 = sequenceParser{id: 536, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p535 = sequenceParser{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p535.items = []parser{&p816, &p14}
	p536.items = []parser{&p816, &p14, &p535}
	var p257 = sequenceParser{id: 257, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p254 = sequenceParser{id: 254, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p253 = charParser{id: 253, chars: []rune{40}}
	p254.items = []parser{&p253}
	var p256 = sequenceParser{id: 256, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p255 = charParser{id: 255, chars: []rune{41}}
	p256.items = []parser{&p255}
	p257.items = []parser{&p267, &p816, &p254, &p816, &p115, &p816, &p120, &p816, &p115, &p816, &p256}
	p537.items = []parser{&p534, &p536, &p816, &p257}
	var p546 = sequenceParser{id: 546, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{774, 198, 784}}
	var p543 = sequenceParser{id: 543, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p538 = charParser{id: 538, chars: []rune{100}}
	var p539 = charParser{id: 539, chars: []rune{101}}
	var p540 = charParser{id: 540, chars: []rune{102}}
	var p541 = charParser{id: 541, chars: []rune{101}}
	var p542 = charParser{id: 542, chars: []rune{114}}
	p543.items = []parser{&p538, &p539, &p540, &p541, &p542}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p544.items = []parser{&p816, &p14}
	p545.items = []parser{&p816, &p14, &p544}
	p546.items = []parser{&p543, &p545, &p816, &p257}
	var p611 = choiceParser{id: 611, commit: 64, name: "assignment", generalizations: []int{774, 198, 784}}
	var p591 = sequenceParser{id: 591, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{611, 774, 198, 784}}
	var p588 = sequenceParser{id: 588, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p585 = charParser{id: 585, chars: []rune{115}}
	var p586 = charParser{id: 586, chars: []rune{101}}
	var p587 = charParser{id: 587, chars: []rune{116}}
	p588.items = []parser{&p585, &p586, &p587}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p589.items = []parser{&p816, &p14}
	p590.items = []parser{&p816, &p14, &p589}
	var p580 = sequenceParser{id: 580, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p577 = sequenceParser{id: 577, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p576 = sequenceParser{id: 576, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p575 = sequenceParser{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p575.items = []parser{&p816, &p14}
	p576.items = []parser{&p14, &p575}
	var p574 = sequenceParser{id: 574, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p573 = charParser{id: 573, chars: []rune{61}}
	p574.items = []parser{&p573}
	p577.items = []parser{&p576, &p816, &p574}
	var p579 = sequenceParser{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p578 = sequenceParser{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p578.items = []parser{&p816, &p14}
	p579.items = []parser{&p816, &p14, &p578}
	p580.items = []parser{&p267, &p816, &p577, &p579, &p816, &p375}
	p591.items = []parser{&p588, &p590, &p816, &p580}
	var p598 = sequenceParser{id: 598, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{611, 774, 198, 784}}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p594.items = []parser{&p816, &p14}
	p595.items = []parser{&p816, &p14, &p594}
	var p593 = sequenceParser{id: 593, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p592 = charParser{id: 592, chars: []rune{61}}
	p593.items = []parser{&p592}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p596.items = []parser{&p816, &p14}
	p597.items = []parser{&p816, &p14, &p596}
	p598.items = []parser{&p267, &p595, &p816, &p593, &p597, &p816, &p375}
	var p610 = sequenceParser{id: 610, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{611, 774, 198, 784}}
	var p602 = sequenceParser{id: 602, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p599 = charParser{id: 599, chars: []rune{115}}
	var p600 = charParser{id: 600, chars: []rune{101}}
	var p601 = charParser{id: 601, chars: []rune{116}}
	p602.items = []parser{&p599, &p600, &p601}
	var p609 = sequenceParser{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p608 = sequenceParser{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p608.items = []parser{&p816, &p14}
	p609.items = []parser{&p816, &p14, &p608}
	var p604 = sequenceParser{id: 604, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p603 = charParser{id: 603, chars: []rune{40}}
	p604.items = []parser{&p603}
	var p605 = sequenceParser{id: 605, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p584 = sequenceParser{id: 584, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p583 = sequenceParser{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p581 = sequenceParser{id: 581, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p581.items = []parser{&p115, &p816, &p580}
	var p582 = sequenceParser{id: 582, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p582.items = []parser{&p816, &p581}
	p583.items = []parser{&p816, &p581, &p582}
	p584.items = []parser{&p580, &p583}
	p605.items = []parser{&p115, &p816, &p584}
	var p607 = sequenceParser{id: 607, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p606 = charParser{id: 606, chars: []rune{41}}
	p607.items = []parser{&p606}
	p610.items = []parser{&p602, &p609, &p816, &p604, &p816, &p605, &p816, &p115, &p816, &p607}
	p611.options = []parser{&p591, &p598, &p610}
	var p783 = sequenceParser{id: 783, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{774, 198, 784}}
	var p776 = sequenceParser{id: 776, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p775 = charParser{id: 775, chars: []rune{40}}
	p776.items = []parser{&p775}
	var p780 = sequenceParser{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p779 = sequenceParser{id: 779, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p779.items = []parser{&p816, &p14}
	p780.items = []parser{&p816, &p14, &p779}
	var p782 = sequenceParser{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p781.items = []parser{&p816, &p14}
	p782.items = []parser{&p816, &p14, &p781}
	var p778 = sequenceParser{id: 778, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p777 = charParser{id: 777, chars: []rune{41}}
	p778.items = []parser{&p777}
	p783.items = []parser{&p776, &p780, &p816, &p774, &p782, &p816, &p778}
	p774.options = []parser{&p490, &p537, &p546, &p611, &p783, &p375}
	var p192 = sequenceParser{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{198}}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{123}}
	p189.items = []parser{&p188}
	var p191 = sequenceParser{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p190 = charParser{id: 190, chars: []rune{125}}
	p191.items = []parser{&p190}
	p192.items = []parser{&p189, &p816, &p798, &p816, &p802, &p816, &p798, &p816, &p191}
	p198.options = []parser{&p774, &p192}
	p201.items = []parser{&p194, &p816, &p115, &p816, &p169, &p816, &p195, &p816, &p115, &p816, &p197, &p200, &p816, &p198}
	p207.items = []parser{&p204, &p206, &p816, &p201}
	var p217 = sequenceParser{id: 217, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p210 = sequenceParser{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p208 = charParser{id: 208, chars: []rune{102}}
	var p209 = charParser{id: 209, chars: []rune{110}}
	p210.items = []parser{&p208, &p209}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p213.items = []parser{&p816, &p14}
	p214.items = []parser{&p816, &p14, &p213}
	var p212 = sequenceParser{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p211 = charParser{id: 211, chars: []rune{126}}
	p212.items = []parser{&p211}
	var p216 = sequenceParser{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p215.items = []parser{&p816, &p14}
	p216.items = []parser{&p816, &p14, &p215}
	p217.items = []parser{&p210, &p214, &p816, &p212, &p216, &p816, &p201}
	var p245 = choiceParser{id: 245, commit: 64, name: "expression-indexer", generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p235 = sequenceParser{id: 235, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{245, 267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p228 = sequenceParser{id: 228, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p227 = charParser{id: 227, chars: []rune{91}}
	p228.items = []parser{&p227}
	var p232 = sequenceParser{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p231.items = []parser{&p816, &p14}
	p232.items = []parser{&p816, &p14, &p231}
	var p234 = sequenceParser{id: 234, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p233 = sequenceParser{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p233.items = []parser{&p816, &p14}
	p234.items = []parser{&p816, &p14, &p233}
	var p230 = sequenceParser{id: 230, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p229 = charParser{id: 229, chars: []rune{93}}
	p230.items = []parser{&p229}
	p235.items = []parser{&p267, &p816, &p228, &p232, &p816, &p375, &p234, &p816, &p230}
	var p244 = sequenceParser{id: 244, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{245, 267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p237 = sequenceParser{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p236 = charParser{id: 236, chars: []rune{91}}
	p237.items = []parser{&p236}
	var p241 = sequenceParser{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p240 = sequenceParser{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p240.items = []parser{&p816, &p14}
	p241.items = []parser{&p816, &p14, &p240}
	var p226 = sequenceParser{id: 226, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{550, 556, 557}}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p375}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p222.items = []parser{&p816, &p14}
	p223.items = []parser{&p816, &p14, &p222}
	var p221 = sequenceParser{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p220 = charParser{id: 220, chars: []rune{58}}
	p221.items = []parser{&p220}
	var p225 = sequenceParser{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p224.items = []parser{&p816, &p14}
	p225.items = []parser{&p816, &p14, &p224}
	var p219 = sequenceParser{id: 219, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p219.items = []parser{&p375}
	p226.items = []parser{&p218, &p223, &p816, &p221, &p225, &p816, &p219}
	var p243 = sequenceParser{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p242 = sequenceParser{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p242.items = []parser{&p816, &p14}
	p243.items = []parser{&p816, &p14, &p242}
	var p239 = sequenceParser{id: 239, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p238 = charParser{id: 238, chars: []rune{93}}
	p239.items = []parser{&p238}
	p244.items = []parser{&p267, &p816, &p237, &p241, &p816, &p226, &p243, &p816, &p239}
	p245.options = []parser{&p235, &p244}
	var p252 = sequenceParser{id: 252, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p248.items = []parser{&p816, &p14}
	p249.items = []parser{&p816, &p14, &p248}
	var p247 = sequenceParser{id: 247, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p246 = charParser{id: 246, chars: []rune{46}}
	p247.items = []parser{&p246}
	var p251 = sequenceParser{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p250.items = []parser{&p816, &p14}
	p251.items = []parser{&p816, &p14, &p250}
	p252.items = []parser{&p267, &p249, &p816, &p247, &p251, &p816, &p105}
	var p266 = sequenceParser{id: 266, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
	var p259 = sequenceParser{id: 259, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p258 = charParser{id: 258, chars: []rune{40}}
	p259.items = []parser{&p258}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p262.items = []parser{&p816, &p14}
	p263.items = []parser{&p816, &p14, &p262}
	var p265 = sequenceParser{id: 265, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p264.items = []parser{&p816, &p14}
	p265.items = []parser{&p816, &p14, &p264}
	var p261 = sequenceParser{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p260 = charParser{id: 260, chars: []rune{41}}
	p261.items = []parser{&p260}
	p266.items = []parser{&p259, &p263, &p816, &p375, &p265, &p816, &p261}
	p267.options = []parser{&p62, &p75, &p88, &p100, &p480, &p105, &p126, &p131, &p160, &p165, &p207, &p217, &p245, &p252, &p257, &p266}
	var p327 = sequenceParser{id: 327, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{375, 333, 334, 335, 336, 337, 338, 557, 550, 784}}
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
	p327.items = []parser{&p326, &p816, &p267}
	var p361 = choiceParser{id: 361, commit: 66, name: "binary-expression", generalizations: []int{375, 557, 550, 784}}
	var p341 = sequenceParser{id: 341, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 334, 335, 336, 337, 338, 375, 557, 550, 784}}
	var p333 = choiceParser{id: 333, commit: 66, name: "operand0", generalizations: []int{334, 335, 336, 337, 338}}
	p333.options = []parser{&p267, &p327}
	var p339 = sequenceParser{id: 339, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
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
	p339.items = []parser{&p328, &p816, &p333}
	var p340 = sequenceParser{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p340.items = []parser{&p816, &p339}
	p341.items = []parser{&p333, &p816, &p339, &p340}
	var p344 = sequenceParser{id: 344, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 335, 336, 337, 338, 375, 557, 550, 784}}
	var p334 = choiceParser{id: 334, commit: 66, name: "operand1", generalizations: []int{335, 336, 337, 338}}
	p334.options = []parser{&p333, &p341}
	var p342 = sequenceParser{id: 342, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
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
	p342.items = []parser{&p329, &p816, &p334}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p343.items = []parser{&p816, &p342}
	p344.items = []parser{&p334, &p816, &p342, &p343}
	var p347 = sequenceParser{id: 347, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 336, 337, 338, 375, 557, 550, 784}}
	var p335 = choiceParser{id: 335, commit: 66, name: "operand2", generalizations: []int{336, 337, 338}}
	p335.options = []parser{&p334, &p344}
	var p345 = sequenceParser{id: 345, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
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
	p345.items = []parser{&p330, &p816, &p335}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p346.items = []parser{&p816, &p345}
	p347.items = []parser{&p335, &p816, &p345, &p346}
	var p350 = sequenceParser{id: 350, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 337, 338, 375, 557, 550, 784}}
	var p336 = choiceParser{id: 336, commit: 66, name: "operand3", generalizations: []int{337, 338}}
	p336.options = []parser{&p335, &p347}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p331 = sequenceParser{id: 331, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p319 = sequenceParser{id: 319, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p317 = charParser{id: 317, chars: []rune{38}}
	var p318 = charParser{id: 318, chars: []rune{38}}
	p319.items = []parser{&p317, &p318}
	p331.items = []parser{&p319}
	p348.items = []parser{&p331, &p816, &p336}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p349.items = []parser{&p816, &p348}
	p350.items = []parser{&p336, &p816, &p348, &p349}
	var p353 = sequenceParser{id: 353, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 338, 375, 557, 550, 784}}
	var p337 = choiceParser{id: 337, commit: 66, name: "operand4", generalizations: []int{338}}
	p337.options = []parser{&p336, &p350}
	var p351 = sequenceParser{id: 351, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p332 = sequenceParser{id: 332, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p322 = sequenceParser{id: 322, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p320 = charParser{id: 320, chars: []rune{124}}
	var p321 = charParser{id: 321, chars: []rune{124}}
	p322.items = []parser{&p320, &p321}
	p332.items = []parser{&p322}
	p351.items = []parser{&p332, &p816, &p337}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p352.items = []parser{&p816, &p351}
	p353.items = []parser{&p337, &p816, &p351, &p352}
	var p360 = sequenceParser{id: 360, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 375, 557, 550, 784}}
	var p338 = choiceParser{id: 338, commit: 66, name: "operand5"}
	p338.options = []parser{&p337, &p353}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p354.items = []parser{&p816, &p14}
	p355.items = []parser{&p14, &p354}
	var p325 = sequenceParser{id: 325, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p323 = charParser{id: 323, chars: []rune{45}}
	var p324 = charParser{id: 324, chars: []rune{62}}
	p325.items = []parser{&p323, &p324}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p356.items = []parser{&p816, &p14}
	p357.items = []parser{&p816, &p14, &p356}
	p358.items = []parser{&p355, &p816, &p325, &p357, &p816, &p338}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p816, &p358}
	p360.items = []parser{&p338, &p816, &p358, &p359}
	p361.options = []parser{&p341, &p344, &p347, &p350, &p353, &p360}
	var p374 = sequenceParser{id: 374, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{375, 557, 550, 784}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p366.items = []parser{&p816, &p14}
	p367.items = []parser{&p816, &p14, &p366}
	var p363 = sequenceParser{id: 363, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p362 = charParser{id: 362, chars: []rune{63}}
	p363.items = []parser{&p362}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p368.items = []parser{&p816, &p14}
	p369.items = []parser{&p816, &p14, &p368}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p370.items = []parser{&p816, &p14}
	p371.items = []parser{&p816, &p14, &p370}
	var p365 = sequenceParser{id: 365, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p364 = charParser{id: 364, chars: []rune{58}}
	p365.items = []parser{&p364}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p372.items = []parser{&p816, &p14}
	p373.items = []parser{&p816, &p14, &p372}
	p374.items = []parser{&p375, &p367, &p816, &p363, &p369, &p816, &p375, &p371, &p816, &p365, &p373, &p816, &p375}
	p375.options = []parser{&p267, &p327, &p361, &p374}
	p186.items = []parser{&p185, &p816, &p375}
	p187.items = []parser{&p183, &p816, &p186}
	var p412 = sequenceParser{id: 412, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{784, 458, 522}}
	var p378 = sequenceParser{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p376 = charParser{id: 376, chars: []rune{105}}
	var p377 = charParser{id: 377, chars: []rune{102}}
	p378.items = []parser{&p376, &p377}
	var p407 = sequenceParser{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p406 = sequenceParser{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p406.items = []parser{&p816, &p14}
	p407.items = []parser{&p816, &p14, &p406}
	var p409 = sequenceParser{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p408.items = []parser{&p816, &p14}
	p409.items = []parser{&p816, &p14, &p408}
	var p411 = sequenceParser{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p395 = sequenceParser{id: 395, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p387.items = []parser{&p816, &p14}
	p388.items = []parser{&p14, &p387}
	var p383 = sequenceParser{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p379 = charParser{id: 379, chars: []rune{101}}
	var p380 = charParser{id: 380, chars: []rune{108}}
	var p381 = charParser{id: 381, chars: []rune{115}}
	var p382 = charParser{id: 382, chars: []rune{101}}
	p383.items = []parser{&p379, &p380, &p381, &p382}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p389.items = []parser{&p816, &p14}
	p390.items = []parser{&p816, &p14, &p389}
	var p386 = sequenceParser{id: 386, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p384 = charParser{id: 384, chars: []rune{105}}
	var p385 = charParser{id: 385, chars: []rune{102}}
	p386.items = []parser{&p384, &p385}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p391.items = []parser{&p816, &p14}
	p392.items = []parser{&p816, &p14, &p391}
	var p394 = sequenceParser{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p393.items = []parser{&p816, &p14}
	p394.items = []parser{&p816, &p14, &p393}
	p395.items = []parser{&p388, &p816, &p383, &p390, &p816, &p386, &p392, &p816, &p375, &p394, &p816, &p192}
	var p410 = sequenceParser{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p410.items = []parser{&p816, &p395}
	p411.items = []parser{&p816, &p395, &p410}
	var p405 = sequenceParser{id: 405, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p402 = sequenceParser{id: 402, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p401 = sequenceParser{id: 401, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p401.items = []parser{&p816, &p14}
	p402.items = []parser{&p14, &p401}
	var p400 = sequenceParser{id: 400, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p396 = charParser{id: 396, chars: []rune{101}}
	var p397 = charParser{id: 397, chars: []rune{108}}
	var p398 = charParser{id: 398, chars: []rune{115}}
	var p399 = charParser{id: 399, chars: []rune{101}}
	p400.items = []parser{&p396, &p397, &p398, &p399}
	var p404 = sequenceParser{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p403 = sequenceParser{id: 403, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p403.items = []parser{&p816, &p14}
	p404.items = []parser{&p816, &p14, &p403}
	p405.items = []parser{&p402, &p816, &p400, &p404, &p816, &p192}
	p412.items = []parser{&p378, &p407, &p816, &p375, &p409, &p816, &p192, &p411, &p816, &p405}
	var p469 = sequenceParser{id: 469, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{458, 784, 522}}
	var p454 = sequenceParser{id: 454, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p448 = charParser{id: 448, chars: []rune{115}}
	var p449 = charParser{id: 449, chars: []rune{119}}
	var p450 = charParser{id: 450, chars: []rune{105}}
	var p451 = charParser{id: 451, chars: []rune{116}}
	var p452 = charParser{id: 452, chars: []rune{99}}
	var p453 = charParser{id: 453, chars: []rune{104}}
	p454.items = []parser{&p448, &p449, &p450, &p451, &p452, &p453}
	var p466 = sequenceParser{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p465 = sequenceParser{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p465.items = []parser{&p816, &p14}
	p466.items = []parser{&p816, &p14, &p465}
	var p468 = sequenceParser{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p467 = sequenceParser{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p467.items = []parser{&p816, &p14}
	p468.items = []parser{&p816, &p14, &p467}
	var p456 = sequenceParser{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p455 = charParser{id: 455, chars: []rune{123}}
	p456.items = []parser{&p455}
	var p462 = sequenceParser{id: 462, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p457 = choiceParser{id: 457, commit: 2}
	var p447 = sequenceParser{id: 447, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{457, 458}}
	var p442 = sequenceParser{id: 442, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p435 = sequenceParser{id: 435, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p431 = charParser{id: 431, chars: []rune{99}}
	var p432 = charParser{id: 432, chars: []rune{97}}
	var p433 = charParser{id: 433, chars: []rune{115}}
	var p434 = charParser{id: 434, chars: []rune{101}}
	p435.items = []parser{&p431, &p432, &p433, &p434}
	var p439 = sequenceParser{id: 439, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p438 = sequenceParser{id: 438, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p438.items = []parser{&p816, &p14}
	p439.items = []parser{&p816, &p14, &p438}
	var p441 = sequenceParser{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p440 = sequenceParser{id: 440, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p440.items = []parser{&p816, &p14}
	p441.items = []parser{&p816, &p14, &p440}
	var p437 = sequenceParser{id: 437, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p436 = charParser{id: 436, chars: []rune{58}}
	p437.items = []parser{&p436}
	p442.items = []parser{&p435, &p439, &p816, &p375, &p441, &p816, &p437}
	var p446 = sequenceParser{id: 446, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p444 = sequenceParser{id: 444, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p443 = charParser{id: 443, chars: []rune{59}}
	p444.items = []parser{&p443}
	var p445 = sequenceParser{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p445.items = []parser{&p816, &p444}
	p446.items = []parser{&p816, &p444, &p445}
	p447.items = []parser{&p442, &p446, &p816, &p784}
	var p430 = sequenceParser{id: 430, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{457, 458, 521, 522}}
	var p425 = sequenceParser{id: 425, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p420 = sequenceParser{id: 420, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p413 = charParser{id: 413, chars: []rune{100}}
	var p414 = charParser{id: 414, chars: []rune{101}}
	var p415 = charParser{id: 415, chars: []rune{102}}
	var p416 = charParser{id: 416, chars: []rune{97}}
	var p417 = charParser{id: 417, chars: []rune{117}}
	var p418 = charParser{id: 418, chars: []rune{108}}
	var p419 = charParser{id: 419, chars: []rune{116}}
	p420.items = []parser{&p413, &p414, &p415, &p416, &p417, &p418, &p419}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p423.items = []parser{&p816, &p14}
	p424.items = []parser{&p816, &p14, &p423}
	var p422 = sequenceParser{id: 422, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p421 = charParser{id: 421, chars: []rune{58}}
	p422.items = []parser{&p421}
	p425.items = []parser{&p420, &p424, &p816, &p422}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p427 = sequenceParser{id: 427, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p426 = charParser{id: 426, chars: []rune{59}}
	p427.items = []parser{&p426}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p428.items = []parser{&p816, &p427}
	p429.items = []parser{&p816, &p427, &p428}
	p430.items = []parser{&p425, &p429, &p816, &p784}
	p457.options = []parser{&p447, &p430}
	var p461 = sequenceParser{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p458 = choiceParser{id: 458, commit: 2}
	p458.options = []parser{&p447, &p430, &p784}
	p459.items = []parser{&p798, &p816, &p458}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p460.items = []parser{&p816, &p459}
	p461.items = []parser{&p816, &p459, &p460}
	p462.items = []parser{&p457, &p461}
	var p464 = sequenceParser{id: 464, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p463 = charParser{id: 463, chars: []rune{125}}
	p464.items = []parser{&p463}
	p469.items = []parser{&p454, &p466, &p816, &p375, &p468, &p816, &p456, &p816, &p798, &p816, &p462, &p816, &p798, &p816, &p464}
	var p531 = sequenceParser{id: 531, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{522, 784}}
	var p518 = sequenceParser{id: 518, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p512 = charParser{id: 512, chars: []rune{115}}
	var p513 = charParser{id: 513, chars: []rune{101}}
	var p514 = charParser{id: 514, chars: []rune{108}}
	var p515 = charParser{id: 515, chars: []rune{101}}
	var p516 = charParser{id: 516, chars: []rune{99}}
	var p517 = charParser{id: 517, chars: []rune{116}}
	p518.items = []parser{&p512, &p513, &p514, &p515, &p516, &p517}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p529.items = []parser{&p816, &p14}
	p530.items = []parser{&p816, &p14, &p529}
	var p520 = sequenceParser{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p519 = charParser{id: 519, chars: []rune{123}}
	p520.items = []parser{&p519}
	var p526 = sequenceParser{id: 526, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p521 = choiceParser{id: 521, commit: 2}
	var p511 = sequenceParser{id: 511, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{521, 522}}
	var p506 = sequenceParser{id: 506, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p499 = sequenceParser{id: 499, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p495 = charParser{id: 495, chars: []rune{99}}
	var p496 = charParser{id: 496, chars: []rune{97}}
	var p497 = charParser{id: 497, chars: []rune{115}}
	var p498 = charParser{id: 498, chars: []rune{101}}
	p499.items = []parser{&p495, &p496, &p497, &p498}
	var p503 = sequenceParser{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p502 = sequenceParser{id: 502, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p502.items = []parser{&p816, &p14}
	p503.items = []parser{&p816, &p14, &p502}
	var p494 = choiceParser{id: 494, commit: 66, name: "communication"}
	var p493 = sequenceParser{id: 493, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{494}}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p491 = sequenceParser{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p491.items = []parser{&p816, &p14}
	p492.items = []parser{&p816, &p14, &p491}
	p493.items = []parser{&p105, &p492, &p816, &p480}
	p494.options = []parser{&p480, &p493, &p490}
	var p505 = sequenceParser{id: 505, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p504 = sequenceParser{id: 504, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p504.items = []parser{&p816, &p14}
	p505.items = []parser{&p816, &p14, &p504}
	var p501 = sequenceParser{id: 501, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p500 = charParser{id: 500, chars: []rune{58}}
	p501.items = []parser{&p500}
	p506.items = []parser{&p499, &p503, &p816, &p494, &p505, &p816, &p501}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p508 = sequenceParser{id: 508, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p507 = charParser{id: 507, chars: []rune{59}}
	p508.items = []parser{&p507}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p509.items = []parser{&p816, &p508}
	p510.items = []parser{&p816, &p508, &p509}
	p511.items = []parser{&p506, &p510, &p816, &p784}
	p521.options = []parser{&p511, &p430}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p522 = choiceParser{id: 522, commit: 2}
	p522.options = []parser{&p511, &p430, &p784}
	p523.items = []parser{&p798, &p816, &p522}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p524.items = []parser{&p816, &p523}
	p525.items = []parser{&p816, &p523, &p524}
	p526.items = []parser{&p521, &p525}
	var p528 = sequenceParser{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p527 = charParser{id: 527, chars: []rune{125}}
	p528.items = []parser{&p527}
	p531.items = []parser{&p518, &p530, &p816, &p520, &p816, &p798, &p816, &p526, &p816, &p798, &p816, &p528}
	var p572 = sequenceParser{id: 572, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{784}}
	var p561 = sequenceParser{id: 561, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p558 = charParser{id: 558, chars: []rune{102}}
	var p559 = charParser{id: 559, chars: []rune{111}}
	var p560 = charParser{id: 560, chars: []rune{114}}
	p561.items = []parser{&p558, &p559, &p560}
	var p571 = choiceParser{id: 571, commit: 2}
	var p567 = sequenceParser{id: 567, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{571}}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p563 = sequenceParser{id: 563, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p562 = sequenceParser{id: 562, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p562.items = []parser{&p816, &p14}
	p563.items = []parser{&p14, &p562}
	var p557 = choiceParser{id: 557, commit: 66, name: "loop-expression"}
	var p556 = choiceParser{id: 556, commit: 64, name: "range-over-expression", generalizations: []int{557}}
	var p555 = sequenceParser{id: 555, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{556, 557}}
	var p552 = sequenceParser{id: 552, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p551 = sequenceParser{id: 551, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p551.items = []parser{&p816, &p14}
	p552.items = []parser{&p816, &p14, &p551}
	var p549 = sequenceParser{id: 549, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p547 = charParser{id: 547, chars: []rune{105}}
	var p548 = charParser{id: 548, chars: []rune{110}}
	p549.items = []parser{&p547, &p548}
	var p554 = sequenceParser{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p553 = sequenceParser{id: 553, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p553.items = []parser{&p816, &p14}
	p554.items = []parser{&p816, &p14, &p553}
	var p550 = choiceParser{id: 550, commit: 2}
	p550.options = []parser{&p375, &p226}
	p555.items = []parser{&p105, &p552, &p816, &p549, &p554, &p816, &p550}
	p556.options = []parser{&p555, &p226}
	p557.options = []parser{&p375, &p556}
	p564.items = []parser{&p563, &p816, &p557}
	var p566 = sequenceParser{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p565.items = []parser{&p816, &p14}
	p566.items = []parser{&p816, &p14, &p565}
	p567.items = []parser{&p564, &p566, &p816, &p192}
	var p570 = sequenceParser{id: 570, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{571}}
	var p569 = sequenceParser{id: 569, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p568 = sequenceParser{id: 568, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p568.items = []parser{&p816, &p14}
	p569.items = []parser{&p14, &p568}
	p570.items = []parser{&p569, &p816, &p192}
	p571.options = []parser{&p567, &p570}
	p572.items = []parser{&p561, &p816, &p571}
	var p720 = choiceParser{id: 720, commit: 66, name: "definition", generalizations: []int{784}}
	var p633 = sequenceParser{id: 633, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{720, 784}}
	var p629 = sequenceParser{id: 629, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p626 = charParser{id: 626, chars: []rune{108}}
	var p627 = charParser{id: 627, chars: []rune{101}}
	var p628 = charParser{id: 628, chars: []rune{116}}
	p629.items = []parser{&p626, &p627, &p628}
	var p632 = sequenceParser{id: 632, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p631 = sequenceParser{id: 631, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p631.items = []parser{&p816, &p14}
	p632.items = []parser{&p816, &p14, &p631}
	var p630 = choiceParser{id: 630, commit: 2}
	var p620 = sequenceParser{id: 620, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{630, 634, 635}}
	var p619 = sequenceParser{id: 619, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p614.items = []parser{&p816, &p14}
	p615.items = []parser{&p14, &p614}
	var p613 = sequenceParser{id: 613, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p612 = charParser{id: 612, chars: []rune{61}}
	p613.items = []parser{&p612}
	p616.items = []parser{&p615, &p816, &p613}
	var p618 = sequenceParser{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p617.items = []parser{&p816, &p14}
	p618.items = []parser{&p816, &p14, &p617}
	p619.items = []parser{&p105, &p816, &p616, &p618, &p816, &p375}
	p620.items = []parser{&p619}
	var p625 = sequenceParser{id: 625, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{630, 634, 635}}
	var p622 = sequenceParser{id: 622, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p621 = charParser{id: 621, chars: []rune{126}}
	p622.items = []parser{&p621}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p623 = sequenceParser{id: 623, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p623.items = []parser{&p816, &p14}
	p624.items = []parser{&p816, &p14, &p623}
	p625.items = []parser{&p622, &p624, &p816, &p619}
	p630.options = []parser{&p620, &p625}
	p633.items = []parser{&p629, &p632, &p816, &p630}
	var p654 = sequenceParser{id: 654, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{720, 784}}
	var p647 = sequenceParser{id: 647, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p644 = charParser{id: 644, chars: []rune{108}}
	var p645 = charParser{id: 645, chars: []rune{101}}
	var p646 = charParser{id: 646, chars: []rune{116}}
	p647.items = []parser{&p644, &p645, &p646}
	var p653 = sequenceParser{id: 653, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p652 = sequenceParser{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p652.items = []parser{&p816, &p14}
	p653.items = []parser{&p816, &p14, &p652}
	var p649 = sequenceParser{id: 649, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p648 = charParser{id: 648, chars: []rune{40}}
	p649.items = []parser{&p648}
	var p639 = sequenceParser{id: 639, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p634 = choiceParser{id: 634, commit: 2}
	p634.options = []parser{&p620, &p625}
	var p638 = sequenceParser{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p635 = choiceParser{id: 635, commit: 2}
	p635.options = []parser{&p620, &p625}
	p636.items = []parser{&p115, &p816, &p635}
	var p637 = sequenceParser{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p637.items = []parser{&p816, &p636}
	p638.items = []parser{&p816, &p636, &p637}
	p639.items = []parser{&p634, &p638}
	var p651 = sequenceParser{id: 651, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p650 = charParser{id: 650, chars: []rune{41}}
	p651.items = []parser{&p650}
	p654.items = []parser{&p647, &p653, &p816, &p649, &p816, &p115, &p816, &p639, &p816, &p115, &p816, &p651}
	var p669 = sequenceParser{id: 669, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{720, 784}}
	var p658 = sequenceParser{id: 658, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p655 = charParser{id: 655, chars: []rune{108}}
	var p656 = charParser{id: 656, chars: []rune{101}}
	var p657 = charParser{id: 657, chars: []rune{116}}
	p658.items = []parser{&p655, &p656, &p657}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p665 = sequenceParser{id: 665, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p665.items = []parser{&p816, &p14}
	p666.items = []parser{&p816, &p14, &p665}
	var p660 = sequenceParser{id: 660, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p659 = charParser{id: 659, chars: []rune{126}}
	p660.items = []parser{&p659}
	var p668 = sequenceParser{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p667.items = []parser{&p816, &p14}
	p668.items = []parser{&p816, &p14, &p667}
	var p662 = sequenceParser{id: 662, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p661 = charParser{id: 661, chars: []rune{40}}
	p662.items = []parser{&p661}
	var p643 = sequenceParser{id: 643, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p640 = sequenceParser{id: 640, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p640.items = []parser{&p115, &p816, &p620}
	var p641 = sequenceParser{id: 641, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p641.items = []parser{&p816, &p640}
	p642.items = []parser{&p816, &p640, &p641}
	p643.items = []parser{&p620, &p642}
	var p664 = sequenceParser{id: 664, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p663 = charParser{id: 663, chars: []rune{41}}
	p664.items = []parser{&p663}
	p669.items = []parser{&p658, &p666, &p816, &p660, &p668, &p816, &p662, &p816, &p115, &p816, &p643, &p816, &p115, &p816, &p664}
	var p685 = sequenceParser{id: 685, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{720, 784}}
	var p681 = sequenceParser{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p679 = charParser{id: 679, chars: []rune{102}}
	var p680 = charParser{id: 680, chars: []rune{110}}
	p681.items = []parser{&p679, &p680}
	var p684 = sequenceParser{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p683 = sequenceParser{id: 683, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p683.items = []parser{&p816, &p14}
	p684.items = []parser{&p816, &p14, &p683}
	var p682 = choiceParser{id: 682, commit: 2}
	var p673 = sequenceParser{id: 673, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{682, 690, 691}}
	var p672 = sequenceParser{id: 672, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p671 = sequenceParser{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p670 = sequenceParser{id: 670, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p670.items = []parser{&p816, &p14}
	p671.items = []parser{&p816, &p14, &p670}
	p672.items = []parser{&p105, &p671, &p816, &p201}
	p673.items = []parser{&p672}
	var p678 = sequenceParser{id: 678, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{682, 690, 691}}
	var p675 = sequenceParser{id: 675, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p674 = charParser{id: 674, chars: []rune{126}}
	p675.items = []parser{&p674}
	var p677 = sequenceParser{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p676 = sequenceParser{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p676.items = []parser{&p816, &p14}
	p677.items = []parser{&p816, &p14, &p676}
	p678.items = []parser{&p675, &p677, &p816, &p672}
	p682.options = []parser{&p673, &p678}
	p685.items = []parser{&p681, &p684, &p816, &p682}
	var p705 = sequenceParser{id: 705, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{720, 784}}
	var p698 = sequenceParser{id: 698, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p696 = charParser{id: 696, chars: []rune{102}}
	var p697 = charParser{id: 697, chars: []rune{110}}
	p698.items = []parser{&p696, &p697}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p703.items = []parser{&p816, &p14}
	p704.items = []parser{&p816, &p14, &p703}
	var p700 = sequenceParser{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p699 = charParser{id: 699, chars: []rune{40}}
	p700.items = []parser{&p699}
	var p695 = sequenceParser{id: 695, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p690 = choiceParser{id: 690, commit: 2}
	p690.options = []parser{&p673, &p678}
	var p694 = sequenceParser{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p692 = sequenceParser{id: 692, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p691 = choiceParser{id: 691, commit: 2}
	p691.options = []parser{&p673, &p678}
	p692.items = []parser{&p115, &p816, &p691}
	var p693 = sequenceParser{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p693.items = []parser{&p816, &p692}
	p694.items = []parser{&p816, &p692, &p693}
	p695.items = []parser{&p690, &p694}
	var p702 = sequenceParser{id: 702, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p701 = charParser{id: 701, chars: []rune{41}}
	p702.items = []parser{&p701}
	p705.items = []parser{&p698, &p704, &p816, &p700, &p816, &p115, &p816, &p695, &p816, &p115, &p816, &p702}
	var p719 = sequenceParser{id: 719, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{720, 784}}
	var p708 = sequenceParser{id: 708, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p706 = charParser{id: 706, chars: []rune{102}}
	var p707 = charParser{id: 707, chars: []rune{110}}
	p708.items = []parser{&p706, &p707}
	var p716 = sequenceParser{id: 716, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p715 = sequenceParser{id: 715, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p715.items = []parser{&p816, &p14}
	p716.items = []parser{&p816, &p14, &p715}
	var p710 = sequenceParser{id: 710, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p709 = charParser{id: 709, chars: []rune{126}}
	p710.items = []parser{&p709}
	var p718 = sequenceParser{id: 718, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p717 = sequenceParser{id: 717, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p717.items = []parser{&p816, &p14}
	p718.items = []parser{&p816, &p14, &p717}
	var p712 = sequenceParser{id: 712, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p711 = charParser{id: 711, chars: []rune{40}}
	p712.items = []parser{&p711}
	var p689 = sequenceParser{id: 689, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p686.items = []parser{&p115, &p816, &p673}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p687.items = []parser{&p816, &p686}
	p688.items = []parser{&p816, &p686, &p687}
	p689.items = []parser{&p673, &p688}
	var p714 = sequenceParser{id: 714, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p713 = charParser{id: 713, chars: []rune{41}}
	p714.items = []parser{&p713}
	p719.items = []parser{&p708, &p716, &p816, &p710, &p718, &p816, &p712, &p816, &p115, &p816, &p689, &p816, &p115, &p816, &p714}
	p720.options = []parser{&p633, &p654, &p669, &p685, &p705, &p719}
	var p763 = choiceParser{id: 763, commit: 64, name: "require", generalizations: []int{784}}
	var p747 = sequenceParser{id: 747, commit: 66, name: "require-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{763, 784}}
	var p744 = sequenceParser{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p737 = charParser{id: 737, chars: []rune{114}}
	var p738 = charParser{id: 738, chars: []rune{101}}
	var p739 = charParser{id: 739, chars: []rune{113}}
	var p740 = charParser{id: 740, chars: []rune{117}}
	var p741 = charParser{id: 741, chars: []rune{105}}
	var p742 = charParser{id: 742, chars: []rune{114}}
	var p743 = charParser{id: 743, chars: []rune{101}}
	p744.items = []parser{&p737, &p738, &p739, &p740, &p741, &p742, &p743}
	var p746 = sequenceParser{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p745 = sequenceParser{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p745.items = []parser{&p816, &p14}
	p746.items = []parser{&p816, &p14, &p745}
	var p732 = choiceParser{id: 732, commit: 64, name: "require-fact"}
	var p731 = sequenceParser{id: 731, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{732}}
	var p723 = choiceParser{id: 723, commit: 2}
	var p722 = sequenceParser{id: 722, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{723}}
	var p721 = charParser{id: 721, chars: []rune{46}}
	p722.items = []parser{&p721}
	p723.options = []parser{&p105, &p722}
	var p728 = sequenceParser{id: 728, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p727 = sequenceParser{id: 727, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p726 = sequenceParser{id: 726, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p726.items = []parser{&p816, &p14}
	p727.items = []parser{&p14, &p726}
	var p725 = sequenceParser{id: 725, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p724 = charParser{id: 724, chars: []rune{61}}
	p725.items = []parser{&p724}
	p728.items = []parser{&p727, &p816, &p725}
	var p730 = sequenceParser{id: 730, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p729 = sequenceParser{id: 729, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p729.items = []parser{&p816, &p14}
	p730.items = []parser{&p816, &p14, &p729}
	p731.items = []parser{&p723, &p816, &p728, &p730, &p816, &p88}
	p732.options = []parser{&p88, &p731}
	p747.items = []parser{&p744, &p746, &p816, &p732}
	var p762 = sequenceParser{id: 762, commit: 66, name: "require-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{763, 784}}
	var p755 = sequenceParser{id: 755, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p748 = charParser{id: 748, chars: []rune{114}}
	var p749 = charParser{id: 749, chars: []rune{101}}
	var p750 = charParser{id: 750, chars: []rune{113}}
	var p751 = charParser{id: 751, chars: []rune{117}}
	var p752 = charParser{id: 752, chars: []rune{105}}
	var p753 = charParser{id: 753, chars: []rune{114}}
	var p754 = charParser{id: 754, chars: []rune{101}}
	p755.items = []parser{&p748, &p749, &p750, &p751, &p752, &p753, &p754}
	var p761 = sequenceParser{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p760 = sequenceParser{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p760.items = []parser{&p816, &p14}
	p761.items = []parser{&p816, &p14, &p760}
	var p757 = sequenceParser{id: 757, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p756 = charParser{id: 756, chars: []rune{40}}
	p757.items = []parser{&p756}
	var p736 = sequenceParser{id: 736, commit: 66, name: "require-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p735 = sequenceParser{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p733 = sequenceParser{id: 733, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p733.items = []parser{&p115, &p816, &p732}
	var p734 = sequenceParser{id: 734, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p734.items = []parser{&p816, &p733}
	p735.items = []parser{&p816, &p733, &p734}
	p736.items = []parser{&p732, &p735}
	var p759 = sequenceParser{id: 759, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p758 = charParser{id: 758, chars: []rune{41}}
	p759.items = []parser{&p758}
	p762.items = []parser{&p755, &p761, &p816, &p757, &p816, &p115, &p816, &p736, &p816, &p115, &p816, &p759}
	p763.options = []parser{&p747, &p762}
	var p773 = sequenceParser{id: 773, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784}}
	var p770 = sequenceParser{id: 770, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p764 = charParser{id: 764, chars: []rune{101}}
	var p765 = charParser{id: 765, chars: []rune{120}}
	var p766 = charParser{id: 766, chars: []rune{112}}
	var p767 = charParser{id: 767, chars: []rune{111}}
	var p768 = charParser{id: 768, chars: []rune{114}}
	var p769 = charParser{id: 769, chars: []rune{116}}
	p770.items = []parser{&p764, &p765, &p766, &p767, &p768, &p769}
	var p772 = sequenceParser{id: 772, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p771 = sequenceParser{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p771.items = []parser{&p816, &p14}
	p772.items = []parser{&p816, &p14, &p771}
	p773.items = []parser{&p770, &p772, &p816, &p720}
	var p793 = sequenceParser{id: 793, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784}}
	var p786 = sequenceParser{id: 786, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p785 = charParser{id: 785, chars: []rune{40}}
	p786.items = []parser{&p785}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p789 = sequenceParser{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p789.items = []parser{&p816, &p14}
	p790.items = []parser{&p816, &p14, &p789}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p791.items = []parser{&p816, &p14}
	p792.items = []parser{&p816, &p14, &p791}
	var p788 = sequenceParser{id: 788, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p787 = charParser{id: 787, chars: []rune{41}}
	p788.items = []parser{&p787}
	p793.items = []parser{&p786, &p790, &p816, &p784, &p792, &p816, &p788}
	p784.options = []parser{&p187, &p412, &p469, &p531, &p572, &p720, &p763, &p773, &p793, &p774}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p799 = sequenceParser{id: 799, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p799.items = []parser{&p798, &p816, &p784}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p800.items = []parser{&p816, &p799}
	p801.items = []parser{&p816, &p799, &p800}
	p802.items = []parser{&p784, &p801}
	p817.items = []parser{&p813, &p816, &p798, &p816, &p802, &p816, &p798}
	p818.items = []parser{&p816, &p817, &p816}
	var b818 = sequenceBuilder{id: 818, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b816 = choiceBuilder{id: 816, commit: 2}
	var b814 = choiceBuilder{id: 814, commit: 70}
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
	b814.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b815 = sequenceBuilder{id: 815, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b41.items = []builder{&b40, &b816, &b38}
	var b42 = sequenceBuilder{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b42.items = []builder{&b816, &b41}
	b43.items = []builder{&b816, &b41, &b42}
	b44.items = []builder{&b38, &b43}
	b815.items = []builder{&b44}
	b816.options = []builder{&b814, &b815}
	var b817 = sequenceBuilder{id: 817, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b813 = sequenceBuilder{id: 813, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b810 = sequenceBuilder{id: 810, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b808 = charBuilder{}
	var b809 = charBuilder{}
	b810.items = []builder{&b808, &b809}
	var b807 = sequenceBuilder{id: 807, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b806 = sequenceBuilder{id: 806, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b804 = sequenceBuilder{id: 804, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b803 = charBuilder{}
	b804.items = []builder{&b803}
	var b805 = sequenceBuilder{id: 805, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b805.items = []builder{&b816, &b804}
	b806.items = []builder{&b804, &b805}
	b807.items = []builder{&b806}
	var b812 = sequenceBuilder{id: 812, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b811 = charBuilder{}
	b812.items = []builder{&b811}
	b813.items = []builder{&b810, &b816, &b807, &b816, &b812}
	var b798 = sequenceBuilder{id: 798, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b796 = choiceBuilder{id: 796, commit: 2}
	var b795 = sequenceBuilder{id: 795, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b794 = charBuilder{}
	b795.items = []builder{&b794}
	var b14 = sequenceBuilder{id: 14, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b796.options = []builder{&b795, &b14}
	var b797 = sequenceBuilder{id: 797, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b797.items = []builder{&b816, &b796}
	b798.items = []builder{&b796, &b797}
	var b802 = sequenceBuilder{id: 802, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b784 = choiceBuilder{id: 784, commit: 66}
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
	b184.items = []builder{&b816, &b14}
	b185.items = []builder{&b14, &b184}
	var b375 = choiceBuilder{id: 375, commit: 66}
	var b267 = choiceBuilder{id: 267, commit: 66}
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
	var b480 = sequenceBuilder{id: 480, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b477 = sequenceBuilder{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b470 = charBuilder{}
	var b471 = charBuilder{}
	var b472 = charBuilder{}
	var b473 = charBuilder{}
	var b474 = charBuilder{}
	var b475 = charBuilder{}
	var b476 = charBuilder{}
	b477.items = []builder{&b470, &b471, &b472, &b473, &b474, &b475, &b476}
	var b479 = sequenceBuilder{id: 479, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b478 = sequenceBuilder{id: 478, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b478.items = []builder{&b816, &b14}
	b479.items = []builder{&b816, &b14, &b478}
	b480.items = []builder{&b477, &b479, &b816, &b267}
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
	b114.items = []builder{&b816, &b113}
	b115.items = []builder{&b113, &b114}
	var b120 = sequenceBuilder{id: 120, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b116 = choiceBuilder{id: 116, commit: 66}
	var b110 = sequenceBuilder{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b109 = sequenceBuilder{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	var b108 = charBuilder{}
	b109.items = []builder{&b106, &b107, &b108}
	b110.items = []builder{&b267, &b816, &b109}
	b116.options = []builder{&b375, &b110}
	var b119 = sequenceBuilder{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b117.items = []builder{&b115, &b816, &b116}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b118.items = []builder{&b816, &b117}
	b119.items = []builder{&b816, &b117, &b118}
	b120.items = []builder{&b116, &b119}
	var b124 = sequenceBuilder{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b123 = charBuilder{}
	b124.items = []builder{&b123}
	b125.items = []builder{&b122, &b816, &b115, &b816, &b120, &b816, &b115, &b816, &b124}
	b126.items = []builder{&b125}
	var b131 = sequenceBuilder{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b128 = sequenceBuilder{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b127 = charBuilder{}
	b128.items = []builder{&b127}
	var b130 = sequenceBuilder{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b129.items = []builder{&b816, &b14}
	b130.items = []builder{&b816, &b14, &b129}
	b131.items = []builder{&b128, &b130, &b816, &b125}
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
	b136.items = []builder{&b816, &b14}
	b137.items = []builder{&b816, &b14, &b136}
	var b139 = sequenceBuilder{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b138.items = []builder{&b816, &b14}
	b139.items = []builder{&b816, &b14, &b138}
	var b135 = sequenceBuilder{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b134 = charBuilder{}
	b135.items = []builder{&b134}
	b140.items = []builder{&b133, &b137, &b816, &b375, &b139, &b816, &b135}
	b141.options = []builder{&b105, &b88, &b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b816, &b14}
	b145.items = []builder{&b816, &b14, &b144}
	var b143 = sequenceBuilder{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b142 = charBuilder{}
	b143.items = []builder{&b142}
	var b147 = sequenceBuilder{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b146.items = []builder{&b816, &b14}
	b147.items = []builder{&b816, &b14, &b146}
	b148.items = []builder{&b141, &b145, &b816, &b143, &b147, &b816, &b375}
	b149.options = []builder{&b148, &b110}
	var b153 = sequenceBuilder{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b150 = choiceBuilder{id: 150, commit: 2}
	b150.options = []builder{&b148, &b110}
	b151.items = []builder{&b115, &b816, &b150}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b152.items = []builder{&b816, &b151}
	b153.items = []builder{&b816, &b151, &b152}
	b154.items = []builder{&b149, &b153}
	var b158 = sequenceBuilder{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b157 = charBuilder{}
	b158.items = []builder{&b157}
	b159.items = []builder{&b156, &b816, &b115, &b816, &b154, &b816, &b115, &b816, &b158}
	b160.items = []builder{&b159}
	var b165 = sequenceBuilder{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b162 = sequenceBuilder{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b161 = charBuilder{}
	b162.items = []builder{&b161}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b163.items = []builder{&b816, &b14}
	b164.items = []builder{&b816, &b14, &b163}
	b165.items = []builder{&b162, &b164, &b816, &b159}
	var b207 = sequenceBuilder{id: 207, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b204 = sequenceBuilder{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b202 = charBuilder{}
	var b203 = charBuilder{}
	b204.items = []builder{&b202, &b203}
	var b206 = sequenceBuilder{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b205.items = []builder{&b816, &b14}
	b206.items = []builder{&b816, &b14, &b205}
	var b201 = sequenceBuilder{id: 201, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b194 = sequenceBuilder{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b193 = charBuilder{}
	b194.items = []builder{&b193}
	var b169 = sequenceBuilder{id: 169, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b168 = sequenceBuilder{id: 168, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b166.items = []builder{&b115, &b816, &b105}
	var b167 = sequenceBuilder{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b167.items = []builder{&b816, &b166}
	b168.items = []builder{&b816, &b166, &b167}
	b169.items = []builder{&b105, &b168}
	var b195 = sequenceBuilder{id: 195, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b176 = sequenceBuilder{id: 176, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b173 = sequenceBuilder{id: 173, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b170 = charBuilder{}
	var b171 = charBuilder{}
	var b172 = charBuilder{}
	b173.items = []builder{&b170, &b171, &b172}
	var b175 = sequenceBuilder{id: 175, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b174 = sequenceBuilder{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b174.items = []builder{&b816, &b14}
	b175.items = []builder{&b816, &b14, &b174}
	b176.items = []builder{&b173, &b175, &b816, &b105}
	b195.items = []builder{&b115, &b816, &b176}
	var b197 = sequenceBuilder{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b196 = charBuilder{}
	b197.items = []builder{&b196}
	var b200 = sequenceBuilder{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b199.items = []builder{&b816, &b14}
	b200.items = []builder{&b816, &b14, &b199}
	var b198 = choiceBuilder{id: 198, commit: 2}
	var b774 = choiceBuilder{id: 774, commit: 66}
	var b490 = sequenceBuilder{id: 490, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b485 = sequenceBuilder{id: 485, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b481 = charBuilder{}
	var b482 = charBuilder{}
	var b483 = charBuilder{}
	var b484 = charBuilder{}
	b485.items = []builder{&b481, &b482, &b483, &b484}
	var b487 = sequenceBuilder{id: 487, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b486.items = []builder{&b816, &b14}
	b487.items = []builder{&b816, &b14, &b486}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b488.items = []builder{&b816, &b14}
	b489.items = []builder{&b816, &b14, &b488}
	b490.items = []builder{&b485, &b487, &b816, &b267, &b489, &b816, &b267}
	var b537 = sequenceBuilder{id: 537, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b534 = sequenceBuilder{id: 534, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b532 = charBuilder{}
	var b533 = charBuilder{}
	b534.items = []builder{&b532, &b533}
	var b536 = sequenceBuilder{id: 536, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b535 = sequenceBuilder{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b535.items = []builder{&b816, &b14}
	b536.items = []builder{&b816, &b14, &b535}
	var b257 = sequenceBuilder{id: 257, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b254 = sequenceBuilder{id: 254, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b253 = charBuilder{}
	b254.items = []builder{&b253}
	var b256 = sequenceBuilder{id: 256, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b255 = charBuilder{}
	b256.items = []builder{&b255}
	b257.items = []builder{&b267, &b816, &b254, &b816, &b115, &b816, &b120, &b816, &b115, &b816, &b256}
	b537.items = []builder{&b534, &b536, &b816, &b257}
	var b546 = sequenceBuilder{id: 546, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b543 = sequenceBuilder{id: 543, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b538 = charBuilder{}
	var b539 = charBuilder{}
	var b540 = charBuilder{}
	var b541 = charBuilder{}
	var b542 = charBuilder{}
	b543.items = []builder{&b538, &b539, &b540, &b541, &b542}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b544.items = []builder{&b816, &b14}
	b545.items = []builder{&b816, &b14, &b544}
	b546.items = []builder{&b543, &b545, &b816, &b257}
	var b611 = choiceBuilder{id: 611, commit: 64, name: "assignment"}
	var b591 = sequenceBuilder{id: 591, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b588 = sequenceBuilder{id: 588, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b585 = charBuilder{}
	var b586 = charBuilder{}
	var b587 = charBuilder{}
	b588.items = []builder{&b585, &b586, &b587}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b589.items = []builder{&b816, &b14}
	b590.items = []builder{&b816, &b14, &b589}
	var b580 = sequenceBuilder{id: 580, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b577 = sequenceBuilder{id: 577, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b576 = sequenceBuilder{id: 576, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b575 = sequenceBuilder{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b575.items = []builder{&b816, &b14}
	b576.items = []builder{&b14, &b575}
	var b574 = sequenceBuilder{id: 574, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b573 = charBuilder{}
	b574.items = []builder{&b573}
	b577.items = []builder{&b576, &b816, &b574}
	var b579 = sequenceBuilder{id: 579, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b578 = sequenceBuilder{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b578.items = []builder{&b816, &b14}
	b579.items = []builder{&b816, &b14, &b578}
	b580.items = []builder{&b267, &b816, &b577, &b579, &b816, &b375}
	b591.items = []builder{&b588, &b590, &b816, &b580}
	var b598 = sequenceBuilder{id: 598, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b594.items = []builder{&b816, &b14}
	b595.items = []builder{&b816, &b14, &b594}
	var b593 = sequenceBuilder{id: 593, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b592 = charBuilder{}
	b593.items = []builder{&b592}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b596.items = []builder{&b816, &b14}
	b597.items = []builder{&b816, &b14, &b596}
	b598.items = []builder{&b267, &b595, &b816, &b593, &b597, &b816, &b375}
	var b610 = sequenceBuilder{id: 610, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b602 = sequenceBuilder{id: 602, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b599 = charBuilder{}
	var b600 = charBuilder{}
	var b601 = charBuilder{}
	b602.items = []builder{&b599, &b600, &b601}
	var b609 = sequenceBuilder{id: 609, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b608 = sequenceBuilder{id: 608, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b608.items = []builder{&b816, &b14}
	b609.items = []builder{&b816, &b14, &b608}
	var b604 = sequenceBuilder{id: 604, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b603 = charBuilder{}
	b604.items = []builder{&b603}
	var b605 = sequenceBuilder{id: 605, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b584 = sequenceBuilder{id: 584, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b583 = sequenceBuilder{id: 583, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b581 = sequenceBuilder{id: 581, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b581.items = []builder{&b115, &b816, &b580}
	var b582 = sequenceBuilder{id: 582, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b582.items = []builder{&b816, &b581}
	b583.items = []builder{&b816, &b581, &b582}
	b584.items = []builder{&b580, &b583}
	b605.items = []builder{&b115, &b816, &b584}
	var b607 = sequenceBuilder{id: 607, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b606 = charBuilder{}
	b607.items = []builder{&b606}
	b610.items = []builder{&b602, &b609, &b816, &b604, &b816, &b605, &b816, &b115, &b816, &b607}
	b611.options = []builder{&b591, &b598, &b610}
	var b783 = sequenceBuilder{id: 783, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b776 = sequenceBuilder{id: 776, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b775 = charBuilder{}
	b776.items = []builder{&b775}
	var b780 = sequenceBuilder{id: 780, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b779 = sequenceBuilder{id: 779, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b779.items = []builder{&b816, &b14}
	b780.items = []builder{&b816, &b14, &b779}
	var b782 = sequenceBuilder{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b781.items = []builder{&b816, &b14}
	b782.items = []builder{&b816, &b14, &b781}
	var b778 = sequenceBuilder{id: 778, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b777 = charBuilder{}
	b778.items = []builder{&b777}
	b783.items = []builder{&b776, &b780, &b816, &b774, &b782, &b816, &b778}
	b774.options = []builder{&b490, &b537, &b546, &b611, &b783, &b375}
	var b192 = sequenceBuilder{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	var b191 = sequenceBuilder{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b190 = charBuilder{}
	b191.items = []builder{&b190}
	b192.items = []builder{&b189, &b816, &b798, &b816, &b802, &b816, &b798, &b816, &b191}
	b198.options = []builder{&b774, &b192}
	b201.items = []builder{&b194, &b816, &b115, &b816, &b169, &b816, &b195, &b816, &b115, &b816, &b197, &b200, &b816, &b198}
	b207.items = []builder{&b204, &b206, &b816, &b201}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b210 = sequenceBuilder{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b208 = charBuilder{}
	var b209 = charBuilder{}
	b210.items = []builder{&b208, &b209}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b213.items = []builder{&b816, &b14}
	b214.items = []builder{&b816, &b14, &b213}
	var b212 = sequenceBuilder{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b211 = charBuilder{}
	b212.items = []builder{&b211}
	var b216 = sequenceBuilder{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b215.items = []builder{&b816, &b14}
	b216.items = []builder{&b816, &b14, &b215}
	b217.items = []builder{&b210, &b214, &b816, &b212, &b216, &b816, &b201}
	var b245 = choiceBuilder{id: 245, commit: 64, name: "expression-indexer"}
	var b235 = sequenceBuilder{id: 235, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b228 = sequenceBuilder{id: 228, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b227 = charBuilder{}
	b228.items = []builder{&b227}
	var b232 = sequenceBuilder{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b231.items = []builder{&b816, &b14}
	b232.items = []builder{&b816, &b14, &b231}
	var b234 = sequenceBuilder{id: 234, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b233 = sequenceBuilder{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b233.items = []builder{&b816, &b14}
	b234.items = []builder{&b816, &b14, &b233}
	var b230 = sequenceBuilder{id: 230, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b229 = charBuilder{}
	b230.items = []builder{&b229}
	b235.items = []builder{&b267, &b816, &b228, &b232, &b816, &b375, &b234, &b816, &b230}
	var b244 = sequenceBuilder{id: 244, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b237 = sequenceBuilder{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b236 = charBuilder{}
	b237.items = []builder{&b236}
	var b241 = sequenceBuilder{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b240 = sequenceBuilder{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b240.items = []builder{&b816, &b14}
	b241.items = []builder{&b816, &b14, &b240}
	var b226 = sequenceBuilder{id: 226, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b375}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b222.items = []builder{&b816, &b14}
	b223.items = []builder{&b816, &b14, &b222}
	var b221 = sequenceBuilder{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b220 = charBuilder{}
	b221.items = []builder{&b220}
	var b225 = sequenceBuilder{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b224.items = []builder{&b816, &b14}
	b225.items = []builder{&b816, &b14, &b224}
	var b219 = sequenceBuilder{id: 219, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b219.items = []builder{&b375}
	b226.items = []builder{&b218, &b223, &b816, &b221, &b225, &b816, &b219}
	var b243 = sequenceBuilder{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b242 = sequenceBuilder{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b242.items = []builder{&b816, &b14}
	b243.items = []builder{&b816, &b14, &b242}
	var b239 = sequenceBuilder{id: 239, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b238 = charBuilder{}
	b239.items = []builder{&b238}
	b244.items = []builder{&b267, &b816, &b237, &b241, &b816, &b226, &b243, &b816, &b239}
	b245.options = []builder{&b235, &b244}
	var b252 = sequenceBuilder{id: 252, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b248.items = []builder{&b816, &b14}
	b249.items = []builder{&b816, &b14, &b248}
	var b247 = sequenceBuilder{id: 247, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b246 = charBuilder{}
	b247.items = []builder{&b246}
	var b251 = sequenceBuilder{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b250.items = []builder{&b816, &b14}
	b251.items = []builder{&b816, &b14, &b250}
	b252.items = []builder{&b267, &b249, &b816, &b247, &b251, &b816, &b105}
	var b266 = sequenceBuilder{id: 266, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b259 = sequenceBuilder{id: 259, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b258 = charBuilder{}
	b259.items = []builder{&b258}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b262.items = []builder{&b816, &b14}
	b263.items = []builder{&b816, &b14, &b262}
	var b265 = sequenceBuilder{id: 265, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b264.items = []builder{&b816, &b14}
	b265.items = []builder{&b816, &b14, &b264}
	var b261 = sequenceBuilder{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b260 = charBuilder{}
	b261.items = []builder{&b260}
	b266.items = []builder{&b259, &b263, &b816, &b375, &b265, &b816, &b261}
	b267.options = []builder{&b62, &b75, &b88, &b100, &b480, &b105, &b126, &b131, &b160, &b165, &b207, &b217, &b245, &b252, &b257, &b266}
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
	b327.items = []builder{&b326, &b816, &b267}
	var b361 = choiceBuilder{id: 361, commit: 66}
	var b341 = sequenceBuilder{id: 341, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b333 = choiceBuilder{id: 333, commit: 66}
	b333.options = []builder{&b267, &b327}
	var b339 = sequenceBuilder{id: 339, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
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
	b339.items = []builder{&b328, &b816, &b333}
	var b340 = sequenceBuilder{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b340.items = []builder{&b816, &b339}
	b341.items = []builder{&b333, &b816, &b339, &b340}
	var b344 = sequenceBuilder{id: 344, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b334 = choiceBuilder{id: 334, commit: 66}
	b334.options = []builder{&b333, &b341}
	var b342 = sequenceBuilder{id: 342, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
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
	b342.items = []builder{&b329, &b816, &b334}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b343.items = []builder{&b816, &b342}
	b344.items = []builder{&b334, &b816, &b342, &b343}
	var b347 = sequenceBuilder{id: 347, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b335 = choiceBuilder{id: 335, commit: 66}
	b335.options = []builder{&b334, &b344}
	var b345 = sequenceBuilder{id: 345, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
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
	b345.items = []builder{&b330, &b816, &b335}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b346.items = []builder{&b816, &b345}
	b347.items = []builder{&b335, &b816, &b345, &b346}
	var b350 = sequenceBuilder{id: 350, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b336 = choiceBuilder{id: 336, commit: 66}
	b336.options = []builder{&b335, &b347}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b331 = sequenceBuilder{id: 331, commit: 66, ranges: [][]int{{1, 1}}}
	var b319 = sequenceBuilder{id: 319, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b317 = charBuilder{}
	var b318 = charBuilder{}
	b319.items = []builder{&b317, &b318}
	b331.items = []builder{&b319}
	b348.items = []builder{&b331, &b816, &b336}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b349.items = []builder{&b816, &b348}
	b350.items = []builder{&b336, &b816, &b348, &b349}
	var b353 = sequenceBuilder{id: 353, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b337 = choiceBuilder{id: 337, commit: 66}
	b337.options = []builder{&b336, &b350}
	var b351 = sequenceBuilder{id: 351, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b332 = sequenceBuilder{id: 332, commit: 66, ranges: [][]int{{1, 1}}}
	var b322 = sequenceBuilder{id: 322, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b320 = charBuilder{}
	var b321 = charBuilder{}
	b322.items = []builder{&b320, &b321}
	b332.items = []builder{&b322}
	b351.items = []builder{&b332, &b816, &b337}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b352.items = []builder{&b816, &b351}
	b353.items = []builder{&b337, &b816, &b351, &b352}
	var b360 = sequenceBuilder{id: 360, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b338 = choiceBuilder{id: 338, commit: 66}
	b338.options = []builder{&b337, &b353}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b354.items = []builder{&b816, &b14}
	b355.items = []builder{&b14, &b354}
	var b325 = sequenceBuilder{id: 325, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b323 = charBuilder{}
	var b324 = charBuilder{}
	b325.items = []builder{&b323, &b324}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b356.items = []builder{&b816, &b14}
	b357.items = []builder{&b816, &b14, &b356}
	b358.items = []builder{&b355, &b816, &b325, &b357, &b816, &b338}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b816, &b358}
	b360.items = []builder{&b338, &b816, &b358, &b359}
	b361.options = []builder{&b341, &b344, &b347, &b350, &b353, &b360}
	var b374 = sequenceBuilder{id: 374, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b366.items = []builder{&b816, &b14}
	b367.items = []builder{&b816, &b14, &b366}
	var b363 = sequenceBuilder{id: 363, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b362 = charBuilder{}
	b363.items = []builder{&b362}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b368.items = []builder{&b816, &b14}
	b369.items = []builder{&b816, &b14, &b368}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b370.items = []builder{&b816, &b14}
	b371.items = []builder{&b816, &b14, &b370}
	var b365 = sequenceBuilder{id: 365, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b364 = charBuilder{}
	b365.items = []builder{&b364}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b372.items = []builder{&b816, &b14}
	b373.items = []builder{&b816, &b14, &b372}
	b374.items = []builder{&b375, &b367, &b816, &b363, &b369, &b816, &b375, &b371, &b816, &b365, &b373, &b816, &b375}
	b375.options = []builder{&b267, &b327, &b361, &b374}
	b186.items = []builder{&b185, &b816, &b375}
	b187.items = []builder{&b183, &b816, &b186}
	var b412 = sequenceBuilder{id: 412, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b378 = sequenceBuilder{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b376 = charBuilder{}
	var b377 = charBuilder{}
	b378.items = []builder{&b376, &b377}
	var b407 = sequenceBuilder{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b406 = sequenceBuilder{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b406.items = []builder{&b816, &b14}
	b407.items = []builder{&b816, &b14, &b406}
	var b409 = sequenceBuilder{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b408.items = []builder{&b816, &b14}
	b409.items = []builder{&b816, &b14, &b408}
	var b411 = sequenceBuilder{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b395 = sequenceBuilder{id: 395, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b387.items = []builder{&b816, &b14}
	b388.items = []builder{&b14, &b387}
	var b383 = sequenceBuilder{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b379 = charBuilder{}
	var b380 = charBuilder{}
	var b381 = charBuilder{}
	var b382 = charBuilder{}
	b383.items = []builder{&b379, &b380, &b381, &b382}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b389.items = []builder{&b816, &b14}
	b390.items = []builder{&b816, &b14, &b389}
	var b386 = sequenceBuilder{id: 386, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b384 = charBuilder{}
	var b385 = charBuilder{}
	b386.items = []builder{&b384, &b385}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b391.items = []builder{&b816, &b14}
	b392.items = []builder{&b816, &b14, &b391}
	var b394 = sequenceBuilder{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b393.items = []builder{&b816, &b14}
	b394.items = []builder{&b816, &b14, &b393}
	b395.items = []builder{&b388, &b816, &b383, &b390, &b816, &b386, &b392, &b816, &b375, &b394, &b816, &b192}
	var b410 = sequenceBuilder{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b410.items = []builder{&b816, &b395}
	b411.items = []builder{&b816, &b395, &b410}
	var b405 = sequenceBuilder{id: 405, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b402 = sequenceBuilder{id: 402, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b401 = sequenceBuilder{id: 401, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b401.items = []builder{&b816, &b14}
	b402.items = []builder{&b14, &b401}
	var b400 = sequenceBuilder{id: 400, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b396 = charBuilder{}
	var b397 = charBuilder{}
	var b398 = charBuilder{}
	var b399 = charBuilder{}
	b400.items = []builder{&b396, &b397, &b398, &b399}
	var b404 = sequenceBuilder{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b403 = sequenceBuilder{id: 403, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b403.items = []builder{&b816, &b14}
	b404.items = []builder{&b816, &b14, &b403}
	b405.items = []builder{&b402, &b816, &b400, &b404, &b816, &b192}
	b412.items = []builder{&b378, &b407, &b816, &b375, &b409, &b816, &b192, &b411, &b816, &b405}
	var b469 = sequenceBuilder{id: 469, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b454 = sequenceBuilder{id: 454, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b448 = charBuilder{}
	var b449 = charBuilder{}
	var b450 = charBuilder{}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	b454.items = []builder{&b448, &b449, &b450, &b451, &b452, &b453}
	var b466 = sequenceBuilder{id: 466, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b465 = sequenceBuilder{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b465.items = []builder{&b816, &b14}
	b466.items = []builder{&b816, &b14, &b465}
	var b468 = sequenceBuilder{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b467 = sequenceBuilder{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b467.items = []builder{&b816, &b14}
	b468.items = []builder{&b816, &b14, &b467}
	var b456 = sequenceBuilder{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b455 = charBuilder{}
	b456.items = []builder{&b455}
	var b462 = sequenceBuilder{id: 462, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b457 = choiceBuilder{id: 457, commit: 2}
	var b447 = sequenceBuilder{id: 447, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b442 = sequenceBuilder{id: 442, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b435 = sequenceBuilder{id: 435, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b431 = charBuilder{}
	var b432 = charBuilder{}
	var b433 = charBuilder{}
	var b434 = charBuilder{}
	b435.items = []builder{&b431, &b432, &b433, &b434}
	var b439 = sequenceBuilder{id: 439, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b438 = sequenceBuilder{id: 438, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b438.items = []builder{&b816, &b14}
	b439.items = []builder{&b816, &b14, &b438}
	var b441 = sequenceBuilder{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b440 = sequenceBuilder{id: 440, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b440.items = []builder{&b816, &b14}
	b441.items = []builder{&b816, &b14, &b440}
	var b437 = sequenceBuilder{id: 437, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b436 = charBuilder{}
	b437.items = []builder{&b436}
	b442.items = []builder{&b435, &b439, &b816, &b375, &b441, &b816, &b437}
	var b446 = sequenceBuilder{id: 446, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b444 = sequenceBuilder{id: 444, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b443 = charBuilder{}
	b444.items = []builder{&b443}
	var b445 = sequenceBuilder{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b445.items = []builder{&b816, &b444}
	b446.items = []builder{&b816, &b444, &b445}
	b447.items = []builder{&b442, &b446, &b816, &b784}
	var b430 = sequenceBuilder{id: 430, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b425 = sequenceBuilder{id: 425, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b420 = sequenceBuilder{id: 420, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b413 = charBuilder{}
	var b414 = charBuilder{}
	var b415 = charBuilder{}
	var b416 = charBuilder{}
	var b417 = charBuilder{}
	var b418 = charBuilder{}
	var b419 = charBuilder{}
	b420.items = []builder{&b413, &b414, &b415, &b416, &b417, &b418, &b419}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b423.items = []builder{&b816, &b14}
	b424.items = []builder{&b816, &b14, &b423}
	var b422 = sequenceBuilder{id: 422, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b421 = charBuilder{}
	b422.items = []builder{&b421}
	b425.items = []builder{&b420, &b424, &b816, &b422}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b427 = sequenceBuilder{id: 427, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b426 = charBuilder{}
	b427.items = []builder{&b426}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b428.items = []builder{&b816, &b427}
	b429.items = []builder{&b816, &b427, &b428}
	b430.items = []builder{&b425, &b429, &b816, &b784}
	b457.options = []builder{&b447, &b430}
	var b461 = sequenceBuilder{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b458 = choiceBuilder{id: 458, commit: 2}
	b458.options = []builder{&b447, &b430, &b784}
	b459.items = []builder{&b798, &b816, &b458}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b460.items = []builder{&b816, &b459}
	b461.items = []builder{&b816, &b459, &b460}
	b462.items = []builder{&b457, &b461}
	var b464 = sequenceBuilder{id: 464, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b463 = charBuilder{}
	b464.items = []builder{&b463}
	b469.items = []builder{&b454, &b466, &b816, &b375, &b468, &b816, &b456, &b816, &b798, &b816, &b462, &b816, &b798, &b816, &b464}
	var b531 = sequenceBuilder{id: 531, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b518 = sequenceBuilder{id: 518, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b512 = charBuilder{}
	var b513 = charBuilder{}
	var b514 = charBuilder{}
	var b515 = charBuilder{}
	var b516 = charBuilder{}
	var b517 = charBuilder{}
	b518.items = []builder{&b512, &b513, &b514, &b515, &b516, &b517}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b529.items = []builder{&b816, &b14}
	b530.items = []builder{&b816, &b14, &b529}
	var b520 = sequenceBuilder{id: 520, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b519 = charBuilder{}
	b520.items = []builder{&b519}
	var b526 = sequenceBuilder{id: 526, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b521 = choiceBuilder{id: 521, commit: 2}
	var b511 = sequenceBuilder{id: 511, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b506 = sequenceBuilder{id: 506, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b499 = sequenceBuilder{id: 499, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b495 = charBuilder{}
	var b496 = charBuilder{}
	var b497 = charBuilder{}
	var b498 = charBuilder{}
	b499.items = []builder{&b495, &b496, &b497, &b498}
	var b503 = sequenceBuilder{id: 503, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b502 = sequenceBuilder{id: 502, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b502.items = []builder{&b816, &b14}
	b503.items = []builder{&b816, &b14, &b502}
	var b494 = choiceBuilder{id: 494, commit: 66}
	var b493 = sequenceBuilder{id: 493, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b491 = sequenceBuilder{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b491.items = []builder{&b816, &b14}
	b492.items = []builder{&b816, &b14, &b491}
	b493.items = []builder{&b105, &b492, &b816, &b480}
	b494.options = []builder{&b480, &b493, &b490}
	var b505 = sequenceBuilder{id: 505, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b504 = sequenceBuilder{id: 504, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b504.items = []builder{&b816, &b14}
	b505.items = []builder{&b816, &b14, &b504}
	var b501 = sequenceBuilder{id: 501, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b500 = charBuilder{}
	b501.items = []builder{&b500}
	b506.items = []builder{&b499, &b503, &b816, &b494, &b505, &b816, &b501}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b508 = sequenceBuilder{id: 508, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b507 = charBuilder{}
	b508.items = []builder{&b507}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b509.items = []builder{&b816, &b508}
	b510.items = []builder{&b816, &b508, &b509}
	b511.items = []builder{&b506, &b510, &b816, &b784}
	b521.options = []builder{&b511, &b430}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b522 = choiceBuilder{id: 522, commit: 2}
	b522.options = []builder{&b511, &b430, &b784}
	b523.items = []builder{&b798, &b816, &b522}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b524.items = []builder{&b816, &b523}
	b525.items = []builder{&b816, &b523, &b524}
	b526.items = []builder{&b521, &b525}
	var b528 = sequenceBuilder{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b527 = charBuilder{}
	b528.items = []builder{&b527}
	b531.items = []builder{&b518, &b530, &b816, &b520, &b816, &b798, &b816, &b526, &b816, &b798, &b816, &b528}
	var b572 = sequenceBuilder{id: 572, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b561 = sequenceBuilder{id: 561, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b558 = charBuilder{}
	var b559 = charBuilder{}
	var b560 = charBuilder{}
	b561.items = []builder{&b558, &b559, &b560}
	var b571 = choiceBuilder{id: 571, commit: 2}
	var b567 = sequenceBuilder{id: 567, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b563 = sequenceBuilder{id: 563, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b562 = sequenceBuilder{id: 562, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b562.items = []builder{&b816, &b14}
	b563.items = []builder{&b14, &b562}
	var b557 = choiceBuilder{id: 557, commit: 66}
	var b556 = choiceBuilder{id: 556, commit: 64, name: "range-over-expression"}
	var b555 = sequenceBuilder{id: 555, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b552 = sequenceBuilder{id: 552, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b551 = sequenceBuilder{id: 551, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b551.items = []builder{&b816, &b14}
	b552.items = []builder{&b816, &b14, &b551}
	var b549 = sequenceBuilder{id: 549, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b547 = charBuilder{}
	var b548 = charBuilder{}
	b549.items = []builder{&b547, &b548}
	var b554 = sequenceBuilder{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b553 = sequenceBuilder{id: 553, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b553.items = []builder{&b816, &b14}
	b554.items = []builder{&b816, &b14, &b553}
	var b550 = choiceBuilder{id: 550, commit: 2}
	b550.options = []builder{&b375, &b226}
	b555.items = []builder{&b105, &b552, &b816, &b549, &b554, &b816, &b550}
	b556.options = []builder{&b555, &b226}
	b557.options = []builder{&b375, &b556}
	b564.items = []builder{&b563, &b816, &b557}
	var b566 = sequenceBuilder{id: 566, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b565.items = []builder{&b816, &b14}
	b566.items = []builder{&b816, &b14, &b565}
	b567.items = []builder{&b564, &b566, &b816, &b192}
	var b570 = sequenceBuilder{id: 570, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b569 = sequenceBuilder{id: 569, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b568 = sequenceBuilder{id: 568, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b568.items = []builder{&b816, &b14}
	b569.items = []builder{&b14, &b568}
	b570.items = []builder{&b569, &b816, &b192}
	b571.options = []builder{&b567, &b570}
	b572.items = []builder{&b561, &b816, &b571}
	var b720 = choiceBuilder{id: 720, commit: 66}
	var b633 = sequenceBuilder{id: 633, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b629 = sequenceBuilder{id: 629, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b626 = charBuilder{}
	var b627 = charBuilder{}
	var b628 = charBuilder{}
	b629.items = []builder{&b626, &b627, &b628}
	var b632 = sequenceBuilder{id: 632, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b631 = sequenceBuilder{id: 631, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b631.items = []builder{&b816, &b14}
	b632.items = []builder{&b816, &b14, &b631}
	var b630 = choiceBuilder{id: 630, commit: 2}
	var b620 = sequenceBuilder{id: 620, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b619 = sequenceBuilder{id: 619, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b614.items = []builder{&b816, &b14}
	b615.items = []builder{&b14, &b614}
	var b613 = sequenceBuilder{id: 613, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b612 = charBuilder{}
	b613.items = []builder{&b612}
	b616.items = []builder{&b615, &b816, &b613}
	var b618 = sequenceBuilder{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b617.items = []builder{&b816, &b14}
	b618.items = []builder{&b816, &b14, &b617}
	b619.items = []builder{&b105, &b816, &b616, &b618, &b816, &b375}
	b620.items = []builder{&b619}
	var b625 = sequenceBuilder{id: 625, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b622 = sequenceBuilder{id: 622, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b621 = charBuilder{}
	b622.items = []builder{&b621}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b623 = sequenceBuilder{id: 623, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b623.items = []builder{&b816, &b14}
	b624.items = []builder{&b816, &b14, &b623}
	b625.items = []builder{&b622, &b624, &b816, &b619}
	b630.options = []builder{&b620, &b625}
	b633.items = []builder{&b629, &b632, &b816, &b630}
	var b654 = sequenceBuilder{id: 654, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b647 = sequenceBuilder{id: 647, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b644 = charBuilder{}
	var b645 = charBuilder{}
	var b646 = charBuilder{}
	b647.items = []builder{&b644, &b645, &b646}
	var b653 = sequenceBuilder{id: 653, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b652 = sequenceBuilder{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b652.items = []builder{&b816, &b14}
	b653.items = []builder{&b816, &b14, &b652}
	var b649 = sequenceBuilder{id: 649, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b648 = charBuilder{}
	b649.items = []builder{&b648}
	var b639 = sequenceBuilder{id: 639, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b634 = choiceBuilder{id: 634, commit: 2}
	b634.options = []builder{&b620, &b625}
	var b638 = sequenceBuilder{id: 638, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b635 = choiceBuilder{id: 635, commit: 2}
	b635.options = []builder{&b620, &b625}
	b636.items = []builder{&b115, &b816, &b635}
	var b637 = sequenceBuilder{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b637.items = []builder{&b816, &b636}
	b638.items = []builder{&b816, &b636, &b637}
	b639.items = []builder{&b634, &b638}
	var b651 = sequenceBuilder{id: 651, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b650 = charBuilder{}
	b651.items = []builder{&b650}
	b654.items = []builder{&b647, &b653, &b816, &b649, &b816, &b115, &b816, &b639, &b816, &b115, &b816, &b651}
	var b669 = sequenceBuilder{id: 669, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b658 = sequenceBuilder{id: 658, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b655 = charBuilder{}
	var b656 = charBuilder{}
	var b657 = charBuilder{}
	b658.items = []builder{&b655, &b656, &b657}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b665 = sequenceBuilder{id: 665, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b665.items = []builder{&b816, &b14}
	b666.items = []builder{&b816, &b14, &b665}
	var b660 = sequenceBuilder{id: 660, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b659 = charBuilder{}
	b660.items = []builder{&b659}
	var b668 = sequenceBuilder{id: 668, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b667.items = []builder{&b816, &b14}
	b668.items = []builder{&b816, &b14, &b667}
	var b662 = sequenceBuilder{id: 662, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b661 = charBuilder{}
	b662.items = []builder{&b661}
	var b643 = sequenceBuilder{id: 643, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b640 = sequenceBuilder{id: 640, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b640.items = []builder{&b115, &b816, &b620}
	var b641 = sequenceBuilder{id: 641, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b641.items = []builder{&b816, &b640}
	b642.items = []builder{&b816, &b640, &b641}
	b643.items = []builder{&b620, &b642}
	var b664 = sequenceBuilder{id: 664, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b663 = charBuilder{}
	b664.items = []builder{&b663}
	b669.items = []builder{&b658, &b666, &b816, &b660, &b668, &b816, &b662, &b816, &b115, &b816, &b643, &b816, &b115, &b816, &b664}
	var b685 = sequenceBuilder{id: 685, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b681 = sequenceBuilder{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b679 = charBuilder{}
	var b680 = charBuilder{}
	b681.items = []builder{&b679, &b680}
	var b684 = sequenceBuilder{id: 684, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b683 = sequenceBuilder{id: 683, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b683.items = []builder{&b816, &b14}
	b684.items = []builder{&b816, &b14, &b683}
	var b682 = choiceBuilder{id: 682, commit: 2}
	var b673 = sequenceBuilder{id: 673, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b672 = sequenceBuilder{id: 672, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b671 = sequenceBuilder{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b670 = sequenceBuilder{id: 670, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b670.items = []builder{&b816, &b14}
	b671.items = []builder{&b816, &b14, &b670}
	b672.items = []builder{&b105, &b671, &b816, &b201}
	b673.items = []builder{&b672}
	var b678 = sequenceBuilder{id: 678, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b675 = sequenceBuilder{id: 675, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b674 = charBuilder{}
	b675.items = []builder{&b674}
	var b677 = sequenceBuilder{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b676 = sequenceBuilder{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b676.items = []builder{&b816, &b14}
	b677.items = []builder{&b816, &b14, &b676}
	b678.items = []builder{&b675, &b677, &b816, &b672}
	b682.options = []builder{&b673, &b678}
	b685.items = []builder{&b681, &b684, &b816, &b682}
	var b705 = sequenceBuilder{id: 705, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b698 = sequenceBuilder{id: 698, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b696 = charBuilder{}
	var b697 = charBuilder{}
	b698.items = []builder{&b696, &b697}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b703.items = []builder{&b816, &b14}
	b704.items = []builder{&b816, &b14, &b703}
	var b700 = sequenceBuilder{id: 700, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b699 = charBuilder{}
	b700.items = []builder{&b699}
	var b695 = sequenceBuilder{id: 695, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b690 = choiceBuilder{id: 690, commit: 2}
	b690.options = []builder{&b673, &b678}
	var b694 = sequenceBuilder{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b692 = sequenceBuilder{id: 692, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b691 = choiceBuilder{id: 691, commit: 2}
	b691.options = []builder{&b673, &b678}
	b692.items = []builder{&b115, &b816, &b691}
	var b693 = sequenceBuilder{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b693.items = []builder{&b816, &b692}
	b694.items = []builder{&b816, &b692, &b693}
	b695.items = []builder{&b690, &b694}
	var b702 = sequenceBuilder{id: 702, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b701 = charBuilder{}
	b702.items = []builder{&b701}
	b705.items = []builder{&b698, &b704, &b816, &b700, &b816, &b115, &b816, &b695, &b816, &b115, &b816, &b702}
	var b719 = sequenceBuilder{id: 719, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b708 = sequenceBuilder{id: 708, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b706 = charBuilder{}
	var b707 = charBuilder{}
	b708.items = []builder{&b706, &b707}
	var b716 = sequenceBuilder{id: 716, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b715 = sequenceBuilder{id: 715, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b715.items = []builder{&b816, &b14}
	b716.items = []builder{&b816, &b14, &b715}
	var b710 = sequenceBuilder{id: 710, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b709 = charBuilder{}
	b710.items = []builder{&b709}
	var b718 = sequenceBuilder{id: 718, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b717 = sequenceBuilder{id: 717, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b717.items = []builder{&b816, &b14}
	b718.items = []builder{&b816, &b14, &b717}
	var b712 = sequenceBuilder{id: 712, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b711 = charBuilder{}
	b712.items = []builder{&b711}
	var b689 = sequenceBuilder{id: 689, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b686.items = []builder{&b115, &b816, &b673}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b687.items = []builder{&b816, &b686}
	b688.items = []builder{&b816, &b686, &b687}
	b689.items = []builder{&b673, &b688}
	var b714 = sequenceBuilder{id: 714, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b713 = charBuilder{}
	b714.items = []builder{&b713}
	b719.items = []builder{&b708, &b716, &b816, &b710, &b718, &b816, &b712, &b816, &b115, &b816, &b689, &b816, &b115, &b816, &b714}
	b720.options = []builder{&b633, &b654, &b669, &b685, &b705, &b719}
	var b763 = choiceBuilder{id: 763, commit: 64, name: "require"}
	var b747 = sequenceBuilder{id: 747, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b744 = sequenceBuilder{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b737 = charBuilder{}
	var b738 = charBuilder{}
	var b739 = charBuilder{}
	var b740 = charBuilder{}
	var b741 = charBuilder{}
	var b742 = charBuilder{}
	var b743 = charBuilder{}
	b744.items = []builder{&b737, &b738, &b739, &b740, &b741, &b742, &b743}
	var b746 = sequenceBuilder{id: 746, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b745 = sequenceBuilder{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b745.items = []builder{&b816, &b14}
	b746.items = []builder{&b816, &b14, &b745}
	var b732 = choiceBuilder{id: 732, commit: 64, name: "require-fact"}
	var b731 = sequenceBuilder{id: 731, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b723 = choiceBuilder{id: 723, commit: 2}
	var b722 = sequenceBuilder{id: 722, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b721 = charBuilder{}
	b722.items = []builder{&b721}
	b723.options = []builder{&b105, &b722}
	var b728 = sequenceBuilder{id: 728, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b727 = sequenceBuilder{id: 727, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b726 = sequenceBuilder{id: 726, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b726.items = []builder{&b816, &b14}
	b727.items = []builder{&b14, &b726}
	var b725 = sequenceBuilder{id: 725, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b724 = charBuilder{}
	b725.items = []builder{&b724}
	b728.items = []builder{&b727, &b816, &b725}
	var b730 = sequenceBuilder{id: 730, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b729 = sequenceBuilder{id: 729, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b729.items = []builder{&b816, &b14}
	b730.items = []builder{&b816, &b14, &b729}
	b731.items = []builder{&b723, &b816, &b728, &b730, &b816, &b88}
	b732.options = []builder{&b88, &b731}
	b747.items = []builder{&b744, &b746, &b816, &b732}
	var b762 = sequenceBuilder{id: 762, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b755 = sequenceBuilder{id: 755, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b748 = charBuilder{}
	var b749 = charBuilder{}
	var b750 = charBuilder{}
	var b751 = charBuilder{}
	var b752 = charBuilder{}
	var b753 = charBuilder{}
	var b754 = charBuilder{}
	b755.items = []builder{&b748, &b749, &b750, &b751, &b752, &b753, &b754}
	var b761 = sequenceBuilder{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b760 = sequenceBuilder{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b760.items = []builder{&b816, &b14}
	b761.items = []builder{&b816, &b14, &b760}
	var b757 = sequenceBuilder{id: 757, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b756 = charBuilder{}
	b757.items = []builder{&b756}
	var b736 = sequenceBuilder{id: 736, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b735 = sequenceBuilder{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b733 = sequenceBuilder{id: 733, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b733.items = []builder{&b115, &b816, &b732}
	var b734 = sequenceBuilder{id: 734, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b734.items = []builder{&b816, &b733}
	b735.items = []builder{&b816, &b733, &b734}
	b736.items = []builder{&b732, &b735}
	var b759 = sequenceBuilder{id: 759, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b758 = charBuilder{}
	b759.items = []builder{&b758}
	b762.items = []builder{&b755, &b761, &b816, &b757, &b816, &b115, &b816, &b736, &b816, &b115, &b816, &b759}
	b763.options = []builder{&b747, &b762}
	var b773 = sequenceBuilder{id: 773, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b770 = sequenceBuilder{id: 770, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b764 = charBuilder{}
	var b765 = charBuilder{}
	var b766 = charBuilder{}
	var b767 = charBuilder{}
	var b768 = charBuilder{}
	var b769 = charBuilder{}
	b770.items = []builder{&b764, &b765, &b766, &b767, &b768, &b769}
	var b772 = sequenceBuilder{id: 772, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b771 = sequenceBuilder{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b771.items = []builder{&b816, &b14}
	b772.items = []builder{&b816, &b14, &b771}
	b773.items = []builder{&b770, &b772, &b816, &b720}
	var b793 = sequenceBuilder{id: 793, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b786 = sequenceBuilder{id: 786, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b785 = charBuilder{}
	b786.items = []builder{&b785}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b789 = sequenceBuilder{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b789.items = []builder{&b816, &b14}
	b790.items = []builder{&b816, &b14, &b789}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b791.items = []builder{&b816, &b14}
	b792.items = []builder{&b816, &b14, &b791}
	var b788 = sequenceBuilder{id: 788, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b787 = charBuilder{}
	b788.items = []builder{&b787}
	b793.items = []builder{&b786, &b790, &b816, &b784, &b792, &b816, &b788}
	b784.options = []builder{&b187, &b412, &b469, &b531, &b572, &b720, &b763, &b773, &b793, &b774}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b799 = sequenceBuilder{id: 799, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b799.items = []builder{&b798, &b816, &b784}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b800.items = []builder{&b816, &b799}
	b801.items = []builder{&b816, &b799, &b800}
	b802.items = []builder{&b784, &b801}
	b817.items = []builder{&b813, &b816, &b798, &b816, &b802, &b816, &b798}
	b818.items = []builder{&b816, &b817, &b816}

	return parseInput(r, &p818, &b818)
}
