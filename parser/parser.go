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
	var p828 = sequenceParser{id: 828, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p826 = choiceParser{id: 826, commit: 2}
	var p824 = choiceParser{id: 824, commit: 70, name: "ws", generalizations: []int{826, 15}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826, 15}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826, 15}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826, 15}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826, 15}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826, 15}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826, 15}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p824.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p825 = sequenceParser{id: 825, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826}}
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
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{806, 15, 112}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p40.items = []parser{&p14, &p826, &p39}
	var p41 = sequenceParser{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p41.items = []parser{&p826, &p40}
	p42.items = []parser{&p826, &p40, &p41}
	p43.items = []parser{&p39, &p42}
	p825.items = []parser{&p43}
	p826.options = []parser{&p824, &p825}
	var p827 = sequenceParser{id: 827, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p823 = sequenceParser{id: 823, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p820 = sequenceParser{id: 820, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p818 = charParser{id: 818, chars: []rune{35}}
	var p819 = charParser{id: 819, chars: []rune{33}}
	p820.items = []parser{&p818, &p819}
	var p817 = sequenceParser{id: 817, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p816 = sequenceParser{id: 816, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p814 = sequenceParser{id: 814, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p813 = charParser{id: 813, not: true, chars: []rune{10}}
	p814.items = []parser{&p813}
	var p815 = sequenceParser{id: 815, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p815.items = []parser{&p826, &p814}
	p816.items = []parser{&p814, &p815}
	p817.items = []parser{&p816}
	var p822 = sequenceParser{id: 822, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p821 = charParser{id: 821, chars: []rune{10}}
	p822.items = []parser{&p821}
	p823.items = []parser{&p820, &p826, &p817, &p826, &p822}
	var p808 = sequenceParser{id: 808, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p806 = choiceParser{id: 806, commit: 2}
	var p805 = sequenceParser{id: 805, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{806}}
	var p804 = charParser{id: 804, chars: []rune{59}}
	p805.items = []parser{&p804}
	p806.options = []parser{&p805, &p14}
	var p807 = sequenceParser{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p807.items = []parser{&p826, &p806}
	p808.items = []parser{&p806, &p807}
	var p812 = sequenceParser{id: 812, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p794 = choiceParser{id: 794, commit: 66, name: "statement", generalizations: []int{481, 542}}
	var p187 = sequenceParser{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{794, 481, 542}}
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
	p15.options = []parser{&p824, &p14}
	p183.items = []parser{&p182, &p15}
	var p186 = sequenceParser{id: 186, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p185 = sequenceParser{id: 185, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p184 = sequenceParser{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p184.items = []parser{&p826, &p14}
	p185.items = []parser{&p14, &p184}
	var p402 = choiceParser{id: 402, commit: 66, name: "expression", generalizations: []int{115, 784, 199, 591, 584, 794}}
	var p273 = choiceParser{id: 273, commit: 66, name: "primary-expression", generalizations: []int{115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p61 = choiceParser{id: 61, commit: 64, name: "int", generalizations: []int{273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p52 = sequenceParser{id: 52, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{61, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p51 = sequenceParser{id: 51, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p50 = charParser{id: 50, ranges: [][]rune{{49, 57}}}
	p51.items = []parser{&p50}
	var p45 = sequenceParser{id: 45, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p44 = charParser{id: 44, ranges: [][]rune{{48, 57}}}
	p45.items = []parser{&p44}
	p52.items = []parser{&p51, &p45}
	var p55 = sequenceParser{id: 55, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{61, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p54 = sequenceParser{id: 54, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p53 = charParser{id: 53, chars: []rune{48}}
	p54.items = []parser{&p53}
	var p47 = sequenceParser{id: 47, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p46 = charParser{id: 46, ranges: [][]rune{{48, 55}}}
	p47.items = []parser{&p46}
	p55.items = []parser{&p54, &p47}
	var p60 = sequenceParser{id: 60, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{61, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
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
	var p74 = choiceParser{id: 74, commit: 72, name: "float", generalizations: []int{273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p69 = sequenceParser{id: 69, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{74, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
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
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{74, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p71 = sequenceParser{id: 71, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p70 = charParser{id: 70, chars: []rune{46}}
	p71.items = []parser{&p70}
	p72.items = []parser{&p71, &p45, &p66}
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{74, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	p73.items = []parser{&p45, &p66}
	p74.options = []parser{&p69, &p72, &p73}
	var p87 = sequenceParser{id: 87, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 115, 140, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 757, 794}}
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
	var p99 = choiceParser{id: 99, commit: 66, name: "bool", generalizations: []int{273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p92 = sequenceParser{id: 92, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{99, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p88 = charParser{id: 88, chars: []rune{116}}
	var p89 = charParser{id: 89, chars: []rune{114}}
	var p90 = charParser{id: 90, chars: []rune{117}}
	var p91 = charParser{id: 91, chars: []rune{101}}
	p92.items = []parser{&p88, &p89, &p90, &p91}
	var p98 = sequenceParser{id: 98, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{99, 273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p93 = charParser{id: 93, chars: []rune{102}}
	var p94 = charParser{id: 94, chars: []rune{97}}
	var p95 = charParser{id: 95, chars: []rune{108}}
	var p96 = charParser{id: 96, chars: []rune{115}}
	var p97 = charParser{id: 97, chars: []rune{101}}
	p98.items = []parser{&p93, &p94, &p95, &p96, &p97}
	p99.options = []parser{&p92, &p98}
	var p515 = sequenceParser{id: 515, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 115, 784, 199, 402, 339, 340, 341, 342, 343, 394, 519, 591, 584, 794}}
	var p507 = sequenceParser{id: 507, commit: 74, name: "receive-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p506 = sequenceParser{id: 506, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p499 = charParser{id: 499, chars: []rune{114}}
	var p500 = charParser{id: 500, chars: []rune{101}}
	var p501 = charParser{id: 501, chars: []rune{99}}
	var p502 = charParser{id: 502, chars: []rune{101}}
	var p503 = charParser{id: 503, chars: []rune{105}}
	var p504 = charParser{id: 504, chars: []rune{118}}
	var p505 = charParser{id: 505, chars: []rune{101}}
	p506.items = []parser{&p499, &p500, &p501, &p502, &p503, &p504, &p505}
	p507.items = []parser{&p506, &p15}
	var p514 = sequenceParser{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p513 = sequenceParser{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p513.items = []parser{&p826, &p14}
	p514.items = []parser{&p826, &p14, &p513}
	p515.items = []parser{&p507, &p514, &p826, &p273}
	var p104 = sequenceParser{id: 104, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{273, 115, 140, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 748, 794}}
	var p101 = sequenceParser{id: 101, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p100 = charParser{id: 100, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p101.items = []parser{&p100}
	var p103 = sequenceParser{id: 103, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p102 = charParser{id: 102, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p103.items = []parser{&p102}
	p104.items = []parser{&p101, &p103}
	var p125 = sequenceParser{id: 125, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{115, 273, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
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
	p113.items = []parser{&p826, &p112}
	p114.items = []parser{&p112, &p113}
	var p119 = sequenceParser{id: 119, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p115 = choiceParser{id: 115, commit: 66, name: "list-item"}
	var p109 = sequenceParser{id: 109, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{115, 148, 149}}
	var p108 = sequenceParser{id: 108, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	p108.items = []parser{&p105, &p106, &p107}
	p109.items = []parser{&p273, &p826, &p108}
	p115.options = []parser{&p402, &p109}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p116.items = []parser{&p114, &p826, &p115}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p117.items = []parser{&p826, &p116}
	p118.items = []parser{&p826, &p116, &p117}
	p119.items = []parser{&p115, &p118}
	var p123 = sequenceParser{id: 123, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p122 = charParser{id: 122, chars: []rune{93}}
	p123.items = []parser{&p122}
	p124.items = []parser{&p121, &p826, &p114, &p826, &p119, &p826, &p114, &p826, &p123}
	p125.items = []parser{&p124}
	var p130 = sequenceParser{id: 130, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p127 = sequenceParser{id: 127, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p126 = charParser{id: 126, chars: []rune{126}}
	p127.items = []parser{&p126}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p128.items = []parser{&p826, &p14}
	p129.items = []parser{&p826, &p14, &p128}
	p130.items = []parser{&p127, &p129, &p826, &p124}
	var p159 = sequenceParser{id: 159, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{273, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
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
	p135.items = []parser{&p826, &p14}
	p136.items = []parser{&p826, &p14, &p135}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p137.items = []parser{&p826, &p14}
	p138.items = []parser{&p826, &p14, &p137}
	var p134 = sequenceParser{id: 134, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p133 = charParser{id: 133, chars: []rune{93}}
	p134.items = []parser{&p133}
	p139.items = []parser{&p132, &p136, &p826, &p402, &p138, &p826, &p134}
	p140.options = []parser{&p104, &p87, &p139}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p143.items = []parser{&p826, &p14}
	p144.items = []parser{&p826, &p14, &p143}
	var p142 = sequenceParser{id: 142, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p141 = charParser{id: 141, chars: []rune{58}}
	p142.items = []parser{&p141}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p145.items = []parser{&p826, &p14}
	p146.items = []parser{&p826, &p14, &p145}
	p147.items = []parser{&p140, &p144, &p826, &p142, &p146, &p826, &p402}
	p148.options = []parser{&p147, &p109}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p149 = choiceParser{id: 149, commit: 2}
	p149.options = []parser{&p147, &p109}
	p150.items = []parser{&p114, &p826, &p149}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p151.items = []parser{&p826, &p150}
	p152.items = []parser{&p826, &p150, &p151}
	p153.items = []parser{&p148, &p152}
	var p157 = sequenceParser{id: 157, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p156 = charParser{id: 156, chars: []rune{125}}
	p157.items = []parser{&p156}
	p158.items = []parser{&p155, &p826, &p114, &p826, &p153, &p826, &p114, &p826, &p157}
	p159.items = []parser{&p158}
	var p164 = sequenceParser{id: 164, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 784, 199, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p161 = sequenceParser{id: 161, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p160 = charParser{id: 160, chars: []rune{126}}
	p161.items = []parser{&p160}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p162.items = []parser{&p826, &p14}
	p163.items = []parser{&p826, &p14, &p162}
	p164.items = []parser{&p161, &p163, &p826, &p158}
	var p208 = sequenceParser{id: 208, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 199, 273, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p205 = sequenceParser{id: 205, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p203 = charParser{id: 203, chars: []rune{102}}
	var p204 = charParser{id: 204, chars: []rune{110}}
	p205.items = []parser{&p203, &p204}
	var p207 = sequenceParser{id: 207, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p206 = sequenceParser{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p206.items = []parser{&p826, &p14}
	p207.items = []parser{&p826, &p14, &p206}
	var p202 = sequenceParser{id: 202, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p194 = sequenceParser{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p193 = charParser{id: 193, chars: []rune{40}}
	p194.items = []parser{&p193}
	var p196 = choiceParser{id: 196, commit: 2}
	var p168 = sequenceParser{id: 168, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{196}}
	var p167 = sequenceParser{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p165.items = []parser{&p114, &p826, &p104}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p166.items = []parser{&p826, &p165}
	p167.items = []parser{&p826, &p165, &p166}
	p168.items = []parser{&p104, &p167}
	var p195 = sequenceParser{id: 195, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{196}}
	var p175 = sequenceParser{id: 175, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{196}}
	var p172 = sequenceParser{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p169 = charParser{id: 169, chars: []rune{46}}
	var p170 = charParser{id: 170, chars: []rune{46}}
	var p171 = charParser{id: 171, chars: []rune{46}}
	p172.items = []parser{&p169, &p170, &p171}
	var p174 = sequenceParser{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p173 = sequenceParser{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p173.items = []parser{&p826, &p14}
	p174.items = []parser{&p826, &p14, &p173}
	p175.items = []parser{&p172, &p174, &p826, &p104}
	p195.items = []parser{&p168, &p826, &p114, &p826, &p175}
	p196.options = []parser{&p168, &p195, &p175}
	var p198 = sequenceParser{id: 198, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p197 = charParser{id: 197, chars: []rune{41}}
	p198.items = []parser{&p197}
	var p201 = sequenceParser{id: 201, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p200 = sequenceParser{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p200.items = []parser{&p826, &p14}
	p201.items = []parser{&p826, &p14, &p200}
	var p199 = choiceParser{id: 199, commit: 2}
	var p784 = choiceParser{id: 784, commit: 66, name: "simple-statement", generalizations: []int{199, 794}}
	var p512 = sequenceParser{id: 512, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 199, 519, 794}}
	var p498 = sequenceParser{id: 498, commit: 74, name: "send-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p497 = sequenceParser{id: 497, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p493 = charParser{id: 493, chars: []rune{115}}
	var p494 = charParser{id: 494, chars: []rune{101}}
	var p495 = charParser{id: 495, chars: []rune{110}}
	var p496 = charParser{id: 496, chars: []rune{100}}
	p497.items = []parser{&p493, &p494, &p495, &p496}
	p498.items = []parser{&p497, &p15}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p508 = sequenceParser{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p508.items = []parser{&p826, &p14}
	p509.items = []parser{&p826, &p14, &p508}
	var p511 = sequenceParser{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p510.items = []parser{&p826, &p14}
	p511.items = []parser{&p826, &p14, &p510}
	p512.items = []parser{&p498, &p509, &p826, &p273, &p511, &p826, &p273}
	var p565 = sequenceParser{id: 565, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 199, 794}}
	var p555 = sequenceParser{id: 555, commit: 74, name: "go-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p554 = sequenceParser{id: 554, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p552 = charParser{id: 552, chars: []rune{103}}
	var p553 = charParser{id: 553, chars: []rune{111}}
	p554.items = []parser{&p552, &p553}
	p555.items = []parser{&p554, &p15}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p563 = sequenceParser{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p563.items = []parser{&p826, &p14}
	p564.items = []parser{&p826, &p14, &p563}
	var p263 = sequenceParser{id: 263, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p260 = sequenceParser{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p259 = charParser{id: 259, chars: []rune{40}}
	p260.items = []parser{&p259}
	var p262 = sequenceParser{id: 262, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p261 = charParser{id: 261, chars: []rune{41}}
	p262.items = []parser{&p261}
	p263.items = []parser{&p273, &p826, &p260, &p826, &p114, &p826, &p119, &p826, &p114, &p826, &p262}
	p565.items = []parser{&p555, &p564, &p826, &p263}
	var p574 = sequenceParser{id: 574, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 199, 794}}
	var p571 = sequenceParser{id: 571, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p566 = charParser{id: 566, chars: []rune{100}}
	var p567 = charParser{id: 567, chars: []rune{101}}
	var p568 = charParser{id: 568, chars: []rune{102}}
	var p569 = charParser{id: 569, chars: []rune{101}}
	var p570 = charParser{id: 570, chars: []rune{114}}
	p571.items = []parser{&p566, &p567, &p568, &p569, &p570}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p572.items = []parser{&p826, &p14}
	p573.items = []parser{&p826, &p14, &p572}
	p574.items = []parser{&p571, &p573, &p826, &p263}
	var p638 = choiceParser{id: 638, commit: 64, name: "assignment", generalizations: []int{784, 199, 794}}
	var p622 = sequenceParser{id: 622, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{638, 784, 199, 794}}
	var p607 = sequenceParser{id: 607, commit: 74, name: "set-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p606 = sequenceParser{id: 606, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p603 = charParser{id: 603, chars: []rune{115}}
	var p604 = charParser{id: 604, chars: []rune{101}}
	var p605 = charParser{id: 605, chars: []rune{116}}
	p606.items = []parser{&p603, &p604, &p605}
	p607.items = []parser{&p606, &p15}
	var p621 = sequenceParser{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p620 = sequenceParser{id: 620, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p620.items = []parser{&p826, &p14}
	p621.items = []parser{&p826, &p14, &p620}
	var p615 = sequenceParser{id: 615, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p612 = sequenceParser{id: 612, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p611 = sequenceParser{id: 611, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p610 = sequenceParser{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p610.items = []parser{&p826, &p14}
	p611.items = []parser{&p14, &p610}
	var p609 = sequenceParser{id: 609, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p608 = charParser{id: 608, chars: []rune{61}}
	p609.items = []parser{&p608}
	p612.items = []parser{&p611, &p826, &p609}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p613 = sequenceParser{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p613.items = []parser{&p826, &p14}
	p614.items = []parser{&p826, &p14, &p613}
	p615.items = []parser{&p273, &p826, &p612, &p614, &p826, &p402}
	p622.items = []parser{&p607, &p621, &p826, &p615}
	var p629 = sequenceParser{id: 629, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{638, 784, 199, 794}}
	var p626 = sequenceParser{id: 626, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p625 = sequenceParser{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p625.items = []parser{&p826, &p14}
	p626.items = []parser{&p826, &p14, &p625}
	var p624 = sequenceParser{id: 624, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p623 = charParser{id: 623, chars: []rune{61}}
	p624.items = []parser{&p623}
	var p628 = sequenceParser{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p627 = sequenceParser{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p627.items = []parser{&p826, &p14}
	p628.items = []parser{&p826, &p14, &p627}
	p629.items = []parser{&p273, &p626, &p826, &p624, &p628, &p826, &p402}
	var p637 = sequenceParser{id: 637, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{638, 784, 199, 794}}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p635 = sequenceParser{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p635.items = []parser{&p826, &p14}
	p636.items = []parser{&p826, &p14, &p635}
	var p631 = sequenceParser{id: 631, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p630 = charParser{id: 630, chars: []rune{40}}
	p631.items = []parser{&p630}
	var p632 = sequenceParser{id: 632, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p619 = sequenceParser{id: 619, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p618 = sequenceParser{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p616 = sequenceParser{id: 616, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p616.items = []parser{&p114, &p826, &p615}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p617.items = []parser{&p826, &p616}
	p618.items = []parser{&p826, &p616, &p617}
	p619.items = []parser{&p615, &p618}
	p632.items = []parser{&p114, &p826, &p619}
	var p634 = sequenceParser{id: 634, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p633 = charParser{id: 633, chars: []rune{41}}
	p634.items = []parser{&p633}
	p637.items = []parser{&p607, &p636, &p826, &p631, &p826, &p632, &p826, &p114, &p826, &p634}
	p638.options = []parser{&p622, &p629, &p637}
	var p793 = sequenceParser{id: 793, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 199, 794}}
	var p786 = sequenceParser{id: 786, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p785 = charParser{id: 785, chars: []rune{40}}
	p786.items = []parser{&p785}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p789 = sequenceParser{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p789.items = []parser{&p826, &p14}
	p790.items = []parser{&p826, &p14, &p789}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p791.items = []parser{&p826, &p14}
	p792.items = []parser{&p826, &p14, &p791}
	var p788 = sequenceParser{id: 788, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p787 = charParser{id: 787, chars: []rune{41}}
	p788.items = []parser{&p787}
	p793.items = []parser{&p786, &p790, &p826, &p784, &p792, &p826, &p788}
	p784.options = []parser{&p512, &p565, &p574, &p638, &p793, &p402}
	var p192 = sequenceParser{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{199}}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{123}}
	p189.items = []parser{&p188}
	var p191 = sequenceParser{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p190 = charParser{id: 190, chars: []rune{125}}
	p191.items = []parser{&p190}
	p192.items = []parser{&p189, &p826, &p808, &p826, &p812, &p826, &p808, &p826, &p191}
	p199.options = []parser{&p784, &p192}
	p202.items = []parser{&p194, &p826, &p114, &p826, &p196, &p826, &p114, &p826, &p198, &p201, &p826, &p199}
	p208.items = []parser{&p205, &p207, &p826, &p202}
	var p218 = sequenceParser{id: 218, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p211 = sequenceParser{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p209 = charParser{id: 209, chars: []rune{102}}
	var p210 = charParser{id: 210, chars: []rune{110}}
	p211.items = []parser{&p209, &p210}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p214.items = []parser{&p826, &p14}
	p215.items = []parser{&p826, &p14, &p214}
	var p213 = sequenceParser{id: 213, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p212 = charParser{id: 212, chars: []rune{126}}
	p213.items = []parser{&p212}
	var p217 = sequenceParser{id: 217, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p216 = sequenceParser{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p216.items = []parser{&p826, &p14}
	p217.items = []parser{&p826, &p14, &p216}
	p218.items = []parser{&p211, &p215, &p826, &p213, &p217, &p826, &p202}
	var p258 = sequenceParser{id: 258, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p257 = sequenceParser{id: 257, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p256 = sequenceParser{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p256.items = []parser{&p826, &p14}
	p257.items = []parser{&p826, &p14, &p256}
	var p255 = sequenceParser{id: 255, commit: 66, name: "index-list", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var p251 = choiceParser{id: 251, commit: 66, name: "index"}
	var p232 = sequenceParser{id: 232, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{251}}
	var p229 = sequenceParser{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p228 = charParser{id: 228, chars: []rune{46}}
	p229.items = []parser{&p228}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p230 = sequenceParser{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p230.items = []parser{&p826, &p14}
	p231.items = []parser{&p826, &p14, &p230}
	p232.items = []parser{&p229, &p231, &p826, &p104}
	var p241 = sequenceParser{id: 241, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{251}}
	var p234 = sequenceParser{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p233 = charParser{id: 233, chars: []rune{91}}
	p234.items = []parser{&p233}
	var p238 = sequenceParser{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p237 = sequenceParser{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p237.items = []parser{&p826, &p14}
	p238.items = []parser{&p826, &p14, &p237}
	var p240 = sequenceParser{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p239 = sequenceParser{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p239.items = []parser{&p826, &p14}
	p240.items = []parser{&p826, &p14, &p239}
	var p236 = sequenceParser{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p235 = charParser{id: 235, chars: []rune{93}}
	p236.items = []parser{&p235}
	p241.items = []parser{&p234, &p238, &p826, &p402, &p240, &p826, &p236}
	var p250 = sequenceParser{id: 250, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{251}}
	var p243 = sequenceParser{id: 243, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p242 = charParser{id: 242, chars: []rune{91}}
	p243.items = []parser{&p242}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p246 = sequenceParser{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p246.items = []parser{&p826, &p14}
	p247.items = []parser{&p826, &p14, &p246}
	var p227 = sequenceParser{id: 227, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{584, 590, 591}}
	var p219 = sequenceParser{id: 219, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p219.items = []parser{&p402}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p223.items = []parser{&p826, &p14}
	p224.items = []parser{&p826, &p14, &p223}
	var p222 = sequenceParser{id: 222, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p221 = charParser{id: 221, chars: []rune{58}}
	p222.items = []parser{&p221}
	var p226 = sequenceParser{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p225 = sequenceParser{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p225.items = []parser{&p826, &p14}
	p226.items = []parser{&p826, &p14, &p225}
	var p220 = sequenceParser{id: 220, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p220.items = []parser{&p402}
	p227.items = []parser{&p219, &p224, &p826, &p222, &p226, &p826, &p220}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p248.items = []parser{&p826, &p14}
	p249.items = []parser{&p826, &p14, &p248}
	var p245 = sequenceParser{id: 245, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p244 = charParser{id: 244, chars: []rune{93}}
	p245.items = []parser{&p244}
	p250.items = []parser{&p243, &p247, &p826, &p227, &p249, &p826, &p245}
	p251.options = []parser{&p232, &p241, &p250}
	var p254 = sequenceParser{id: 254, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p253 = sequenceParser{id: 253, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p252 = sequenceParser{id: 252, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p252.items = []parser{&p826, &p14}
	p253.items = []parser{&p14, &p252}
	p254.items = []parser{&p253, &p826, &p251}
	p255.items = []parser{&p251, &p826, &p254}
	p258.items = []parser{&p273, &p257, &p826, &p255}
	var p272 = sequenceParser{id: 272, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{273, 402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p265 = sequenceParser{id: 265, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p264 = charParser{id: 264, chars: []rune{40}}
	p265.items = []parser{&p264}
	var p269 = sequenceParser{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p268 = sequenceParser{id: 268, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p268.items = []parser{&p826, &p14}
	p269.items = []parser{&p826, &p14, &p268}
	var p271 = sequenceParser{id: 271, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p270 = sequenceParser{id: 270, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p270.items = []parser{&p826, &p14}
	p271.items = []parser{&p826, &p14, &p270}
	var p267 = sequenceParser{id: 267, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p266 = charParser{id: 266, chars: []rune{41}}
	p267.items = []parser{&p266}
	p272.items = []parser{&p265, &p269, &p826, &p402, &p271, &p826, &p267}
	p273.options = []parser{&p61, &p74, &p87, &p99, &p515, &p104, &p125, &p130, &p159, &p164, &p208, &p218, &p258, &p263, &p272}
	var p333 = sequenceParser{id: 333, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{402, 339, 340, 341, 342, 343, 394, 591, 584, 794}}
	var p332 = choiceParser{id: 332, commit: 66, name: "unary-operator"}
	var p292 = sequenceParser{id: 292, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p291 = charParser{id: 291, chars: []rune{43}}
	p292.items = []parser{&p291}
	var p294 = sequenceParser{id: 294, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p293 = charParser{id: 293, chars: []rune{45}}
	p294.items = []parser{&p293}
	var p275 = sequenceParser{id: 275, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p274 = charParser{id: 274, chars: []rune{94}}
	p275.items = []parser{&p274}
	var p306 = sequenceParser{id: 306, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{332}}
	var p305 = charParser{id: 305, chars: []rune{33}}
	p306.items = []parser{&p305}
	p332.options = []parser{&p292, &p294, &p275, &p306}
	p333.items = []parser{&p332, &p826, &p273}
	var p380 = choiceParser{id: 380, commit: 66, name: "binary-expression", generalizations: []int{402, 394, 591, 584, 794}}
	var p351 = sequenceParser{id: 351, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 340, 341, 342, 343, 402, 394, 591, 584, 794}}
	var p339 = choiceParser{id: 339, commit: 66, name: "operand0", generalizations: []int{340, 341, 342, 343}}
	p339.options = []parser{&p273, &p333}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p345 = sequenceParser{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p345.items = []parser{&p826, &p14}
	p346.items = []parser{&p14, &p345}
	var p334 = choiceParser{id: 334, commit: 66, name: "binary-op0"}
	var p277 = sequenceParser{id: 277, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p276 = charParser{id: 276, chars: []rune{38}}
	p277.items = []parser{&p276}
	var p284 = sequenceParser{id: 284, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p282 = charParser{id: 282, chars: []rune{38}}
	var p283 = charParser{id: 283, chars: []rune{94}}
	p284.items = []parser{&p282, &p283}
	var p287 = sequenceParser{id: 287, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p285 = charParser{id: 285, chars: []rune{60}}
	var p286 = charParser{id: 286, chars: []rune{60}}
	p287.items = []parser{&p285, &p286}
	var p290 = sequenceParser{id: 290, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{334}}
	var p288 = charParser{id: 288, chars: []rune{62}}
	var p289 = charParser{id: 289, chars: []rune{62}}
	p290.items = []parser{&p288, &p289}
	var p296 = sequenceParser{id: 296, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p295 = charParser{id: 295, chars: []rune{42}}
	p296.items = []parser{&p295}
	var p298 = sequenceParser{id: 298, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p297 = charParser{id: 297, chars: []rune{47}}
	p298.items = []parser{&p297}
	var p300 = sequenceParser{id: 300, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{334}}
	var p299 = charParser{id: 299, chars: []rune{37}}
	p300.items = []parser{&p299}
	p334.options = []parser{&p277, &p284, &p287, &p290, &p296, &p298, &p300}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p347.items = []parser{&p826, &p14}
	p348.items = []parser{&p826, &p14, &p347}
	p349.items = []parser{&p346, &p826, &p334, &p348, &p826, &p339}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p350.items = []parser{&p826, &p349}
	p351.items = []parser{&p339, &p826, &p349, &p350}
	var p358 = sequenceParser{id: 358, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 341, 342, 343, 402, 394, 591, 584, 794}}
	var p340 = choiceParser{id: 340, commit: 66, name: "operand1", generalizations: []int{341, 342, 343}}
	p340.options = []parser{&p339, &p351}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p352.items = []parser{&p826, &p14}
	p353.items = []parser{&p14, &p352}
	var p335 = choiceParser{id: 335, commit: 66, name: "binary-op1"}
	var p279 = sequenceParser{id: 279, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p278 = charParser{id: 278, chars: []rune{124}}
	p279.items = []parser{&p278}
	var p281 = sequenceParser{id: 281, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p280 = charParser{id: 280, chars: []rune{94}}
	p281.items = []parser{&p280}
	var p302 = sequenceParser{id: 302, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p301 = charParser{id: 301, chars: []rune{43}}
	p302.items = []parser{&p301}
	var p304 = sequenceParser{id: 304, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p303 = charParser{id: 303, chars: []rune{45}}
	p304.items = []parser{&p303}
	p335.options = []parser{&p279, &p281, &p302, &p304}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p354.items = []parser{&p826, &p14}
	p355.items = []parser{&p826, &p14, &p354}
	p356.items = []parser{&p353, &p826, &p335, &p355, &p826, &p340}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p357.items = []parser{&p826, &p356}
	p358.items = []parser{&p340, &p826, &p356, &p357}
	var p365 = sequenceParser{id: 365, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 342, 343, 402, 394, 591, 584, 794}}
	var p341 = choiceParser{id: 341, commit: 66, name: "operand2", generalizations: []int{342, 343}}
	p341.options = []parser{&p340, &p358}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p826, &p14}
	p360.items = []parser{&p14, &p359}
	var p336 = choiceParser{id: 336, commit: 66, name: "binary-op2"}
	var p309 = sequenceParser{id: 309, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{336}}
	var p307 = charParser{id: 307, chars: []rune{61}}
	var p308 = charParser{id: 308, chars: []rune{61}}
	p309.items = []parser{&p307, &p308}
	var p312 = sequenceParser{id: 312, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{336}}
	var p310 = charParser{id: 310, chars: []rune{33}}
	var p311 = charParser{id: 311, chars: []rune{61}}
	p312.items = []parser{&p310, &p311}
	var p314 = sequenceParser{id: 314, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{336}}
	var p313 = charParser{id: 313, chars: []rune{60}}
	p314.items = []parser{&p313}
	var p317 = sequenceParser{id: 317, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{336}}
	var p315 = charParser{id: 315, chars: []rune{60}}
	var p316 = charParser{id: 316, chars: []rune{61}}
	p317.items = []parser{&p315, &p316}
	var p319 = sequenceParser{id: 319, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{336}}
	var p318 = charParser{id: 318, chars: []rune{62}}
	p319.items = []parser{&p318}
	var p322 = sequenceParser{id: 322, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{336}}
	var p320 = charParser{id: 320, chars: []rune{62}}
	var p321 = charParser{id: 321, chars: []rune{61}}
	p322.items = []parser{&p320, &p321}
	p336.options = []parser{&p309, &p312, &p314, &p317, &p319, &p322}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p361.items = []parser{&p826, &p14}
	p362.items = []parser{&p826, &p14, &p361}
	p363.items = []parser{&p360, &p826, &p336, &p362, &p826, &p341}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p364.items = []parser{&p826, &p363}
	p365.items = []parser{&p341, &p826, &p363, &p364}
	var p372 = sequenceParser{id: 372, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 343, 402, 394, 591, 584, 794}}
	var p342 = choiceParser{id: 342, commit: 66, name: "operand3", generalizations: []int{343}}
	p342.options = []parser{&p341, &p365}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p366.items = []parser{&p826, &p14}
	p367.items = []parser{&p14, &p366}
	var p337 = sequenceParser{id: 337, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p325 = sequenceParser{id: 325, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p323 = charParser{id: 323, chars: []rune{38}}
	var p324 = charParser{id: 324, chars: []rune{38}}
	p325.items = []parser{&p323, &p324}
	p337.items = []parser{&p325}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p368.items = []parser{&p826, &p14}
	p369.items = []parser{&p826, &p14, &p368}
	p370.items = []parser{&p367, &p826, &p337, &p369, &p826, &p342}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p371.items = []parser{&p826, &p370}
	p372.items = []parser{&p342, &p826, &p370, &p371}
	var p379 = sequenceParser{id: 379, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{380, 402, 394, 591, 584, 794}}
	var p343 = choiceParser{id: 343, commit: 66, name: "operand4"}
	p343.options = []parser{&p342, &p372}
	var p377 = sequenceParser{id: 377, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p374 = sequenceParser{id: 374, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p373.items = []parser{&p826, &p14}
	p374.items = []parser{&p14, &p373}
	var p338 = sequenceParser{id: 338, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p328 = sequenceParser{id: 328, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p326 = charParser{id: 326, chars: []rune{124}}
	var p327 = charParser{id: 327, chars: []rune{124}}
	p328.items = []parser{&p326, &p327}
	p338.items = []parser{&p328}
	var p376 = sequenceParser{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p375 = sequenceParser{id: 375, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p375.items = []parser{&p826, &p14}
	p376.items = []parser{&p826, &p14, &p375}
	p377.items = []parser{&p374, &p826, &p338, &p376, &p826, &p343}
	var p378 = sequenceParser{id: 378, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p378.items = []parser{&p826, &p377}
	p379.items = []parser{&p343, &p826, &p377, &p378}
	p380.options = []parser{&p351, &p358, &p365, &p372, &p379}
	var p393 = sequenceParser{id: 393, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{402, 394, 591, 584, 794}}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p385 = sequenceParser{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p385.items = []parser{&p826, &p14}
	p386.items = []parser{&p826, &p14, &p385}
	var p382 = sequenceParser{id: 382, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p381 = charParser{id: 381, chars: []rune{63}}
	p382.items = []parser{&p381}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p387.items = []parser{&p826, &p14}
	p388.items = []parser{&p826, &p14, &p387}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p389.items = []parser{&p826, &p14}
	p390.items = []parser{&p826, &p14, &p389}
	var p384 = sequenceParser{id: 384, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p383 = charParser{id: 383, chars: []rune{58}}
	p384.items = []parser{&p383}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p391.items = []parser{&p826, &p14}
	p392.items = []parser{&p826, &p14, &p391}
	p393.items = []parser{&p402, &p386, &p826, &p382, &p388, &p826, &p402, &p390, &p826, &p384, &p392, &p826, &p402}
	var p401 = sequenceParser{id: 401, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{402, 591, 584, 794}}
	var p394 = choiceParser{id: 394, commit: 66, name: "chainingOperand"}
	p394.options = []parser{&p273, &p333, &p380, &p393}
	var p399 = sequenceParser{id: 399, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p396 = sequenceParser{id: 396, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p395 = sequenceParser{id: 395, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p395.items = []parser{&p826, &p14}
	p396.items = []parser{&p14, &p395}
	var p331 = sequenceParser{id: 331, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p329 = charParser{id: 329, chars: []rune{45}}
	var p330 = charParser{id: 330, chars: []rune{62}}
	p331.items = []parser{&p329, &p330}
	var p398 = sequenceParser{id: 398, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p397 = sequenceParser{id: 397, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p397.items = []parser{&p826, &p14}
	p398.items = []parser{&p826, &p14, &p397}
	p399.items = []parser{&p396, &p826, &p331, &p398, &p826, &p394}
	var p400 = sequenceParser{id: 400, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p400.items = []parser{&p826, &p399}
	p401.items = []parser{&p394, &p826, &p399, &p400}
	p402.options = []parser{&p273, &p333, &p380, &p393, &p401}
	p186.items = []parser{&p185, &p826, &p402}
	p187.items = []parser{&p183, &p826, &p186}
	var p433 = sequenceParser{id: 433, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{794, 481, 542}}
	var p406 = sequenceParser{id: 406, commit: 74, name: "if-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p405 = sequenceParser{id: 405, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p403 = charParser{id: 403, chars: []rune{105}}
	var p404 = charParser{id: 404, chars: []rune{102}}
	p405.items = []parser{&p403, &p404}
	p406.items = []parser{&p405, &p15}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p427.items = []parser{&p826, &p14}
	p428.items = []parser{&p826, &p14, &p427}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p429.items = []parser{&p826, &p14}
	p430.items = []parser{&p826, &p14, &p429}
	var p432 = sequenceParser{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p421 = sequenceParser{id: 421, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p413 = sequenceParser{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p413.items = []parser{&p826, &p14}
	p414.items = []parser{&p14, &p413}
	var p412 = sequenceParser{id: 412, commit: 74, name: "else-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p411 = sequenceParser{id: 411, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p407 = charParser{id: 407, chars: []rune{101}}
	var p408 = charParser{id: 408, chars: []rune{108}}
	var p409 = charParser{id: 409, chars: []rune{115}}
	var p410 = charParser{id: 410, chars: []rune{101}}
	p411.items = []parser{&p407, &p408, &p409, &p410}
	p412.items = []parser{&p411, &p15}
	var p416 = sequenceParser{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p415.items = []parser{&p826, &p14}
	p416.items = []parser{&p826, &p14, &p415}
	var p418 = sequenceParser{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p417 = sequenceParser{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p417.items = []parser{&p826, &p14}
	p418.items = []parser{&p826, &p14, &p417}
	var p420 = sequenceParser{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p419 = sequenceParser{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p419.items = []parser{&p826, &p14}
	p420.items = []parser{&p826, &p14, &p419}
	p421.items = []parser{&p414, &p826, &p412, &p416, &p826, &p406, &p418, &p826, &p402, &p420, &p826, &p192}
	var p431 = sequenceParser{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p431.items = []parser{&p826, &p421}
	p432.items = []parser{&p826, &p421, &p431}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p422.items = []parser{&p826, &p14}
	p423.items = []parser{&p14, &p422}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p424.items = []parser{&p826, &p14}
	p425.items = []parser{&p826, &p14, &p424}
	p426.items = []parser{&p423, &p826, &p412, &p425, &p826, &p192}
	p433.items = []parser{&p406, &p428, &p826, &p402, &p430, &p826, &p192, &p432, &p826, &p426}
	var p492 = sequenceParser{id: 492, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{481, 794, 542}}
	var p447 = sequenceParser{id: 447, commit: 74, name: "switch-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p446 = sequenceParser{id: 446, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p440 = charParser{id: 440, chars: []rune{115}}
	var p441 = charParser{id: 441, chars: []rune{119}}
	var p442 = charParser{id: 442, chars: []rune{105}}
	var p443 = charParser{id: 443, chars: []rune{116}}
	var p444 = charParser{id: 444, chars: []rune{99}}
	var p445 = charParser{id: 445, chars: []rune{104}}
	p446.items = []parser{&p440, &p441, &p442, &p443, &p444, &p445}
	p447.items = []parser{&p446, &p15}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p488 = sequenceParser{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p488.items = []parser{&p826, &p14}
	p489.items = []parser{&p826, &p14, &p488}
	var p491 = sequenceParser{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p490 = sequenceParser{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p490.items = []parser{&p826, &p14}
	p491.items = []parser{&p826, &p14, &p490}
	var p479 = sequenceParser{id: 479, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p478 = charParser{id: 478, chars: []rune{123}}
	p479.items = []parser{&p478}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p480 = choiceParser{id: 480, commit: 2}
	var p477 = sequenceParser{id: 477, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{480, 481}}
	var p472 = sequenceParser{id: 472, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p439 = sequenceParser{id: 439, commit: 74, name: "case-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p438 = sequenceParser{id: 438, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p434 = charParser{id: 434, chars: []rune{99}}
	var p435 = charParser{id: 435, chars: []rune{97}}
	var p436 = charParser{id: 436, chars: []rune{115}}
	var p437 = charParser{id: 437, chars: []rune{101}}
	p438.items = []parser{&p434, &p435, &p436, &p437}
	p439.items = []parser{&p438, &p15}
	var p469 = sequenceParser{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p468 = sequenceParser{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p468.items = []parser{&p826, &p14}
	p469.items = []parser{&p826, &p14, &p468}
	var p471 = sequenceParser{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p470 = sequenceParser{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p470.items = []parser{&p826, &p14}
	p471.items = []parser{&p826, &p14, &p470}
	var p467 = sequenceParser{id: 467, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p466 = charParser{id: 466, chars: []rune{58}}
	p467.items = []parser{&p466}
	p472.items = []parser{&p439, &p469, &p826, &p402, &p471, &p826, &p467}
	var p476 = sequenceParser{id: 476, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p474 = sequenceParser{id: 474, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p473 = charParser{id: 473, chars: []rune{59}}
	p474.items = []parser{&p473}
	var p475 = sequenceParser{id: 475, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p475.items = []parser{&p826, &p474}
	p476.items = []parser{&p826, &p474, &p475}
	p477.items = []parser{&p472, &p476, &p826, &p794}
	var p465 = sequenceParser{id: 465, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{480, 481, 541, 542}}
	var p460 = sequenceParser{id: 460, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p455 = sequenceParser{id: 455, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p448 = charParser{id: 448, chars: []rune{100}}
	var p449 = charParser{id: 449, chars: []rune{101}}
	var p450 = charParser{id: 450, chars: []rune{102}}
	var p451 = charParser{id: 451, chars: []rune{97}}
	var p452 = charParser{id: 452, chars: []rune{117}}
	var p453 = charParser{id: 453, chars: []rune{108}}
	var p454 = charParser{id: 454, chars: []rune{116}}
	p455.items = []parser{&p448, &p449, &p450, &p451, &p452, &p453, &p454}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p458 = sequenceParser{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p458.items = []parser{&p826, &p14}
	p459.items = []parser{&p826, &p14, &p458}
	var p457 = sequenceParser{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p456 = charParser{id: 456, chars: []rune{58}}
	p457.items = []parser{&p456}
	p460.items = []parser{&p455, &p459, &p826, &p457}
	var p464 = sequenceParser{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p462 = sequenceParser{id: 462, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p461 = charParser{id: 461, chars: []rune{59}}
	p462.items = []parser{&p461}
	var p463 = sequenceParser{id: 463, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p463.items = []parser{&p826, &p462}
	p464.items = []parser{&p826, &p462, &p463}
	p465.items = []parser{&p460, &p464, &p826, &p794}
	p480.options = []parser{&p477, &p465}
	var p484 = sequenceParser{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p482 = sequenceParser{id: 482, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p481 = choiceParser{id: 481, commit: 2}
	p481.options = []parser{&p477, &p465, &p794}
	p482.items = []parser{&p808, &p826, &p481}
	var p483 = sequenceParser{id: 483, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p483.items = []parser{&p826, &p482}
	p484.items = []parser{&p826, &p482, &p483}
	p485.items = []parser{&p480, &p484}
	var p487 = sequenceParser{id: 487, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p486 = charParser{id: 486, chars: []rune{125}}
	p487.items = []parser{&p486}
	p492.items = []parser{&p447, &p489, &p826, &p402, &p491, &p826, &p479, &p826, &p808, &p826, &p485, &p826, &p808, &p826, &p487}
	var p551 = sequenceParser{id: 551, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{542, 794}}
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
	p549.items = []parser{&p826, &p14}
	p550.items = []parser{&p826, &p14, &p549}
	var p540 = sequenceParser{id: 540, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p539 = charParser{id: 539, chars: []rune{123}}
	p540.items = []parser{&p539}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p541 = choiceParser{id: 541, commit: 2}
	var p531 = sequenceParser{id: 531, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{541, 542}}
	var p526 = sequenceParser{id: 526, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p522 = sequenceParser{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p522.items = []parser{&p826, &p14}
	p523.items = []parser{&p826, &p14, &p522}
	var p519 = choiceParser{id: 519, commit: 66, name: "communication"}
	var p518 = sequenceParser{id: 518, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{519}}
	var p517 = sequenceParser{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p516 = sequenceParser{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p516.items = []parser{&p826, &p14}
	p517.items = []parser{&p826, &p14, &p516}
	p518.items = []parser{&p104, &p517, &p826, &p515}
	p519.options = []parser{&p512, &p515, &p518}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p524.items = []parser{&p826, &p14}
	p525.items = []parser{&p826, &p14, &p524}
	var p521 = sequenceParser{id: 521, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p520 = charParser{id: 520, chars: []rune{58}}
	p521.items = []parser{&p520}
	p526.items = []parser{&p439, &p523, &p826, &p519, &p525, &p826, &p521}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p528 = sequenceParser{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p527 = charParser{id: 527, chars: []rune{59}}
	p528.items = []parser{&p527}
	var p529 = sequenceParser{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p529.items = []parser{&p826, &p528}
	p530.items = []parser{&p826, &p528, &p529}
	p531.items = []parser{&p526, &p530, &p826, &p794}
	p541.options = []parser{&p531, &p465}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p543 = sequenceParser{id: 543, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p542 = choiceParser{id: 542, commit: 2}
	p542.options = []parser{&p531, &p465, &p794}
	p543.items = []parser{&p808, &p826, &p542}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p544.items = []parser{&p826, &p543}
	p545.items = []parser{&p826, &p543, &p544}
	p546.items = []parser{&p541, &p545}
	var p548 = sequenceParser{id: 548, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p547 = charParser{id: 547, chars: []rune{125}}
	p548.items = []parser{&p547}
	p551.items = []parser{&p538, &p550, &p826, &p540, &p826, &p808, &p826, &p546, &p826, &p808, &p826, &p548}
	var p602 = sequenceParser{id: 602, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{794}}
	var p583 = sequenceParser{id: 583, commit: 74, name: "for-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p582 = sequenceParser{id: 582, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p579 = charParser{id: 579, chars: []rune{102}}
	var p580 = charParser{id: 580, chars: []rune{111}}
	var p581 = charParser{id: 581, chars: []rune{114}}
	p582.items = []parser{&p579, &p580, &p581}
	p583.items = []parser{&p582, &p15}
	var p601 = choiceParser{id: 601, commit: 2}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{601}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p592 = sequenceParser{id: 592, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p592.items = []parser{&p826, &p14}
	p593.items = []parser{&p14, &p592}
	var p591 = choiceParser{id: 591, commit: 66, name: "loop-expression"}
	var p590 = choiceParser{id: 590, commit: 64, name: "range-over-expression", generalizations: []int{591}}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{590, 591}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p585.items = []parser{&p826, &p14}
	p586.items = []parser{&p826, &p14, &p585}
	var p578 = sequenceParser{id: 578, commit: 74, name: "in-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p577 = sequenceParser{id: 577, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p575 = charParser{id: 575, chars: []rune{105}}
	var p576 = charParser{id: 576, chars: []rune{110}}
	p577.items = []parser{&p575, &p576}
	p578.items = []parser{&p577, &p15}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p587.items = []parser{&p826, &p14}
	p588.items = []parser{&p826, &p14, &p587}
	var p584 = choiceParser{id: 584, commit: 2}
	p584.options = []parser{&p402, &p227}
	p589.items = []parser{&p104, &p586, &p826, &p578, &p588, &p826, &p584}
	p590.options = []parser{&p589, &p227}
	p591.options = []parser{&p402, &p590}
	p594.items = []parser{&p593, &p826, &p591}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p595.items = []parser{&p826, &p14}
	p596.items = []parser{&p826, &p14, &p595}
	p597.items = []parser{&p594, &p596, &p826, &p192}
	var p600 = sequenceParser{id: 600, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{601}}
	var p599 = sequenceParser{id: 599, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p598.items = []parser{&p826, &p14}
	p599.items = []parser{&p14, &p598}
	p600.items = []parser{&p599, &p826, &p192}
	p601.options = []parser{&p597, &p600}
	p602.items = []parser{&p583, &p826, &p601}
	var p740 = choiceParser{id: 740, commit: 66, name: "definition", generalizations: []int{794}}
	var p661 = sequenceParser{id: 661, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 794}}
	var p643 = sequenceParser{id: 643, commit: 74, name: "let-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p642 = sequenceParser{id: 642, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p639 = charParser{id: 639, chars: []rune{108}}
	var p640 = charParser{id: 640, chars: []rune{101}}
	var p641 = charParser{id: 641, chars: []rune{116}}
	p642.items = []parser{&p639, &p640, &p641}
	p643.items = []parser{&p642, &p15}
	var p660 = sequenceParser{id: 660, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p659 = sequenceParser{id: 659, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p659.items = []parser{&p826, &p14}
	p660.items = []parser{&p826, &p14, &p659}
	var p658 = choiceParser{id: 658, commit: 2}
	var p652 = sequenceParser{id: 652, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{658, 662, 663}}
	var p651 = sequenceParser{id: 651, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p648 = sequenceParser{id: 648, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p647 = sequenceParser{id: 647, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p646 = sequenceParser{id: 646, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p646.items = []parser{&p826, &p14}
	p647.items = []parser{&p14, &p646}
	var p645 = sequenceParser{id: 645, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p644 = charParser{id: 644, chars: []rune{61}}
	p645.items = []parser{&p644}
	p648.items = []parser{&p647, &p826, &p645}
	var p650 = sequenceParser{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p649 = sequenceParser{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p649.items = []parser{&p826, &p14}
	p650.items = []parser{&p826, &p14, &p649}
	p651.items = []parser{&p104, &p826, &p648, &p650, &p826, &p402}
	p652.items = []parser{&p651}
	var p657 = sequenceParser{id: 657, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{658, 662, 663}}
	var p654 = sequenceParser{id: 654, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p653 = charParser{id: 653, chars: []rune{126}}
	p654.items = []parser{&p653}
	var p656 = sequenceParser{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p655 = sequenceParser{id: 655, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p655.items = []parser{&p826, &p14}
	p656.items = []parser{&p826, &p14, &p655}
	p657.items = []parser{&p654, &p656, &p826, &p651}
	p658.options = []parser{&p652, &p657}
	p661.items = []parser{&p643, &p660, &p826, &p658}
	var p678 = sequenceParser{id: 678, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 794}}
	var p677 = sequenceParser{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p676 = sequenceParser{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p676.items = []parser{&p826, &p14}
	p677.items = []parser{&p826, &p14, &p676}
	var p673 = sequenceParser{id: 673, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p672 = charParser{id: 672, chars: []rune{40}}
	p673.items = []parser{&p672}
	var p667 = sequenceParser{id: 667, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p662 = choiceParser{id: 662, commit: 2}
	p662.options = []parser{&p652, &p657}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p664 = sequenceParser{id: 664, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p663 = choiceParser{id: 663, commit: 2}
	p663.options = []parser{&p652, &p657}
	p664.items = []parser{&p114, &p826, &p663}
	var p665 = sequenceParser{id: 665, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p665.items = []parser{&p826, &p664}
	p666.items = []parser{&p826, &p664, &p665}
	p667.items = []parser{&p662, &p666}
	var p675 = sequenceParser{id: 675, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p674 = charParser{id: 674, chars: []rune{41}}
	p675.items = []parser{&p674}
	p678.items = []parser{&p643, &p677, &p826, &p673, &p826, &p114, &p826, &p667, &p826, &p114, &p826, &p675}
	var p689 = sequenceParser{id: 689, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 794}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p685 = sequenceParser{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p685.items = []parser{&p826, &p14}
	p686.items = []parser{&p826, &p14, &p685}
	var p680 = sequenceParser{id: 680, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p679 = charParser{id: 679, chars: []rune{126}}
	p680.items = []parser{&p679}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p687.items = []parser{&p826, &p14}
	p688.items = []parser{&p826, &p14, &p687}
	var p682 = sequenceParser{id: 682, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p681 = charParser{id: 681, chars: []rune{40}}
	p682.items = []parser{&p681}
	var p671 = sequenceParser{id: 671, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p670 = sequenceParser{id: 670, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p668 = sequenceParser{id: 668, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p668.items = []parser{&p114, &p826, &p652}
	var p669 = sequenceParser{id: 669, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p669.items = []parser{&p826, &p668}
	p670.items = []parser{&p826, &p668, &p669}
	p671.items = []parser{&p652, &p670}
	var p684 = sequenceParser{id: 684, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p683 = charParser{id: 683, chars: []rune{41}}
	p684.items = []parser{&p683}
	p689.items = []parser{&p643, &p686, &p826, &p680, &p688, &p826, &p682, &p826, &p114, &p826, &p671, &p826, &p114, &p826, &p684}
	var p705 = sequenceParser{id: 705, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 794}}
	var p701 = sequenceParser{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p699 = charParser{id: 699, chars: []rune{102}}
	var p700 = charParser{id: 700, chars: []rune{110}}
	p701.items = []parser{&p699, &p700}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p703.items = []parser{&p826, &p14}
	p704.items = []parser{&p826, &p14, &p703}
	var p702 = choiceParser{id: 702, commit: 2}
	var p693 = sequenceParser{id: 693, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{702, 710, 711}}
	var p692 = sequenceParser{id: 692, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p691 = sequenceParser{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p690 = sequenceParser{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p690.items = []parser{&p826, &p14}
	p691.items = []parser{&p826, &p14, &p690}
	p692.items = []parser{&p104, &p691, &p826, &p202}
	p693.items = []parser{&p692}
	var p698 = sequenceParser{id: 698, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{702, 710, 711}}
	var p695 = sequenceParser{id: 695, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p694 = charParser{id: 694, chars: []rune{126}}
	p695.items = []parser{&p694}
	var p697 = sequenceParser{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p696.items = []parser{&p826, &p14}
	p697.items = []parser{&p826, &p14, &p696}
	p698.items = []parser{&p695, &p697, &p826, &p692}
	p702.options = []parser{&p693, &p698}
	p705.items = []parser{&p701, &p704, &p826, &p702}
	var p725 = sequenceParser{id: 725, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 794}}
	var p718 = sequenceParser{id: 718, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p716 = charParser{id: 716, chars: []rune{102}}
	var p717 = charParser{id: 717, chars: []rune{110}}
	p718.items = []parser{&p716, &p717}
	var p724 = sequenceParser{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p723 = sequenceParser{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p723.items = []parser{&p826, &p14}
	p724.items = []parser{&p826, &p14, &p723}
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
	p712.items = []parser{&p114, &p826, &p711}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p713.items = []parser{&p826, &p712}
	p714.items = []parser{&p826, &p712, &p713}
	p715.items = []parser{&p710, &p714}
	var p722 = sequenceParser{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p721 = charParser{id: 721, chars: []rune{41}}
	p722.items = []parser{&p721}
	p725.items = []parser{&p718, &p724, &p826, &p720, &p826, &p114, &p826, &p715, &p826, &p114, &p826, &p722}
	var p739 = sequenceParser{id: 739, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{740, 794}}
	var p728 = sequenceParser{id: 728, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p726 = charParser{id: 726, chars: []rune{102}}
	var p727 = charParser{id: 727, chars: []rune{110}}
	p728.items = []parser{&p726, &p727}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p735 = sequenceParser{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p735.items = []parser{&p826, &p14}
	p736.items = []parser{&p826, &p14, &p735}
	var p730 = sequenceParser{id: 730, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p729 = charParser{id: 729, chars: []rune{126}}
	p730.items = []parser{&p729}
	var p738 = sequenceParser{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p737.items = []parser{&p826, &p14}
	p738.items = []parser{&p826, &p14, &p737}
	var p732 = sequenceParser{id: 732, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p731 = charParser{id: 731, chars: []rune{40}}
	p732.items = []parser{&p731}
	var p709 = sequenceParser{id: 709, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p708 = sequenceParser{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p706 = sequenceParser{id: 706, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p706.items = []parser{&p114, &p826, &p693}
	var p707 = sequenceParser{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p707.items = []parser{&p826, &p706}
	p708.items = []parser{&p826, &p706, &p707}
	p709.items = []parser{&p693, &p708}
	var p734 = sequenceParser{id: 734, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p733 = charParser{id: 733, chars: []rune{41}}
	p734.items = []parser{&p733}
	p739.items = []parser{&p728, &p736, &p826, &p730, &p738, &p826, &p732, &p826, &p114, &p826, &p709, &p826, &p114, &p826, &p734}
	p740.options = []parser{&p661, &p678, &p689, &p705, &p725, &p739}
	var p772 = choiceParser{id: 772, commit: 64, name: "use", generalizations: []int{794}}
	var p764 = sequenceParser{id: 764, commit: 66, name: "use-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{772, 794}}
	var p745 = sequenceParser{id: 745, commit: 74, name: "use-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p744 = sequenceParser{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p741 = charParser{id: 741, chars: []rune{117}}
	var p742 = charParser{id: 742, chars: []rune{115}}
	var p743 = charParser{id: 743, chars: []rune{101}}
	p744.items = []parser{&p741, &p742, &p743}
	p745.items = []parser{&p744, &p15}
	var p763 = sequenceParser{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p762 = sequenceParser{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p762.items = []parser{&p826, &p14}
	p763.items = []parser{&p826, &p14, &p762}
	var p757 = choiceParser{id: 757, commit: 64, name: "use-fact"}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{757}}
	var p748 = choiceParser{id: 748, commit: 2}
	var p747 = sequenceParser{id: 747, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{748}}
	var p746 = charParser{id: 746, chars: []rune{46}}
	p747.items = []parser{&p746}
	p748.options = []parser{&p104, &p747}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p751 = sequenceParser{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p751.items = []parser{&p826, &p14}
	p752.items = []parser{&p14, &p751}
	var p750 = sequenceParser{id: 750, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p749 = charParser{id: 749, chars: []rune{61}}
	p750.items = []parser{&p749}
	p753.items = []parser{&p752, &p826, &p750}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p754.items = []parser{&p826, &p14}
	p755.items = []parser{&p826, &p14, &p754}
	p756.items = []parser{&p748, &p826, &p753, &p755, &p826, &p87}
	p757.options = []parser{&p87, &p756}
	p764.items = []parser{&p745, &p763, &p826, &p757}
	var p771 = sequenceParser{id: 771, commit: 66, name: "use-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{772, 794}}
	var p770 = sequenceParser{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p769 = sequenceParser{id: 769, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p769.items = []parser{&p826, &p14}
	p770.items = []parser{&p826, &p14, &p769}
	var p766 = sequenceParser{id: 766, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p765 = charParser{id: 765, chars: []rune{40}}
	p766.items = []parser{&p765}
	var p761 = sequenceParser{id: 761, commit: 66, name: "use-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p760 = sequenceParser{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p758 = sequenceParser{id: 758, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p758.items = []parser{&p114, &p826, &p757}
	var p759 = sequenceParser{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p759.items = []parser{&p826, &p758}
	p760.items = []parser{&p826, &p758, &p759}
	p761.items = []parser{&p757, &p760}
	var p768 = sequenceParser{id: 768, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p767 = charParser{id: 767, chars: []rune{41}}
	p768.items = []parser{&p767}
	p771.items = []parser{&p745, &p770, &p826, &p766, &p826, &p114, &p826, &p761, &p826, &p114, &p826, &p768}
	p772.options = []parser{&p764, &p771}
	var p783 = sequenceParser{id: 783, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794}}
	var p780 = sequenceParser{id: 780, commit: 74, name: "export-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p779 = sequenceParser{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p773 = charParser{id: 773, chars: []rune{101}}
	var p774 = charParser{id: 774, chars: []rune{120}}
	var p775 = charParser{id: 775, chars: []rune{112}}
	var p776 = charParser{id: 776, chars: []rune{111}}
	var p777 = charParser{id: 777, chars: []rune{114}}
	var p778 = charParser{id: 778, chars: []rune{116}}
	p779.items = []parser{&p773, &p774, &p775, &p776, &p777, &p778}
	p780.items = []parser{&p779, &p15}
	var p782 = sequenceParser{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p781.items = []parser{&p826, &p14}
	p782.items = []parser{&p826, &p14, &p781}
	p783.items = []parser{&p780, &p782, &p826, &p740}
	var p803 = sequenceParser{id: 803, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794}}
	var p796 = sequenceParser{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p795 = charParser{id: 795, chars: []rune{40}}
	p796.items = []parser{&p795}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p799 = sequenceParser{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p799.items = []parser{&p826, &p14}
	p800.items = []parser{&p826, &p14, &p799}
	var p802 = sequenceParser{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p801.items = []parser{&p826, &p14}
	p802.items = []parser{&p826, &p14, &p801}
	var p798 = sequenceParser{id: 798, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p797 = charParser{id: 797, chars: []rune{41}}
	p798.items = []parser{&p797}
	p803.items = []parser{&p796, &p800, &p826, &p794, &p802, &p826, &p798}
	p794.options = []parser{&p187, &p433, &p492, &p551, &p602, &p740, &p772, &p783, &p803, &p784}
	var p811 = sequenceParser{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p809 = sequenceParser{id: 809, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p809.items = []parser{&p808, &p826, &p794}
	var p810 = sequenceParser{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p810.items = []parser{&p826, &p809}
	p811.items = []parser{&p826, &p809, &p810}
	p812.items = []parser{&p794, &p811}
	p827.items = []parser{&p823, &p826, &p808, &p826, &p812, &p826, &p808}
	p828.items = []parser{&p826, &p827, &p826}
	var b828 = sequenceBuilder{id: 828, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b826 = choiceBuilder{id: 826, commit: 2}
	var b824 = choiceBuilder{id: 824, commit: 70}
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
	b824.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b825 = sequenceBuilder{id: 825, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b40.items = []builder{&b14, &b826, &b39}
	var b41 = sequenceBuilder{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b41.items = []builder{&b826, &b40}
	b42.items = []builder{&b826, &b40, &b41}
	b43.items = []builder{&b39, &b42}
	b825.items = []builder{&b43}
	b826.options = []builder{&b824, &b825}
	var b827 = sequenceBuilder{id: 827, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b823 = sequenceBuilder{id: 823, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b820 = sequenceBuilder{id: 820, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b818 = charBuilder{}
	var b819 = charBuilder{}
	b820.items = []builder{&b818, &b819}
	var b817 = sequenceBuilder{id: 817, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b816 = sequenceBuilder{id: 816, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b814 = sequenceBuilder{id: 814, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b813 = charBuilder{}
	b814.items = []builder{&b813}
	var b815 = sequenceBuilder{id: 815, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b815.items = []builder{&b826, &b814}
	b816.items = []builder{&b814, &b815}
	b817.items = []builder{&b816}
	var b822 = sequenceBuilder{id: 822, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b821 = charBuilder{}
	b822.items = []builder{&b821}
	b823.items = []builder{&b820, &b826, &b817, &b826, &b822}
	var b808 = sequenceBuilder{id: 808, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b806 = choiceBuilder{id: 806, commit: 2}
	var b805 = sequenceBuilder{id: 805, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b804 = charBuilder{}
	b805.items = []builder{&b804}
	b806.options = []builder{&b805, &b14}
	var b807 = sequenceBuilder{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b807.items = []builder{&b826, &b806}
	b808.items = []builder{&b806, &b807}
	var b812 = sequenceBuilder{id: 812, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b794 = choiceBuilder{id: 794, commit: 66}
	var b187 = sequenceBuilder{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
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
	b15.options = []builder{&b824, &b14}
	b183.items = []builder{&b182, &b15}
	var b186 = sequenceBuilder{id: 186, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b185 = sequenceBuilder{id: 185, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b184 = sequenceBuilder{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b184.items = []builder{&b826, &b14}
	b185.items = []builder{&b14, &b184}
	var b402 = choiceBuilder{id: 402, commit: 66}
	var b273 = choiceBuilder{id: 273, commit: 66}
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
	var b515 = sequenceBuilder{id: 515, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b507 = sequenceBuilder{id: 507, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b506 = sequenceBuilder{id: 506, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b499 = charBuilder{}
	var b500 = charBuilder{}
	var b501 = charBuilder{}
	var b502 = charBuilder{}
	var b503 = charBuilder{}
	var b504 = charBuilder{}
	var b505 = charBuilder{}
	b506.items = []builder{&b499, &b500, &b501, &b502, &b503, &b504, &b505}
	b507.items = []builder{&b506, &b15}
	var b514 = sequenceBuilder{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b513 = sequenceBuilder{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b513.items = []builder{&b826, &b14}
	b514.items = []builder{&b826, &b14, &b513}
	b515.items = []builder{&b507, &b514, &b826, &b273}
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
	b113.items = []builder{&b826, &b112}
	b114.items = []builder{&b112, &b113}
	var b119 = sequenceBuilder{id: 119, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b115 = choiceBuilder{id: 115, commit: 66}
	var b109 = sequenceBuilder{id: 109, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b108 = sequenceBuilder{id: 108, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	b108.items = []builder{&b105, &b106, &b107}
	b109.items = []builder{&b273, &b826, &b108}
	b115.options = []builder{&b402, &b109}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b116.items = []builder{&b114, &b826, &b115}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b117.items = []builder{&b826, &b116}
	b118.items = []builder{&b826, &b116, &b117}
	b119.items = []builder{&b115, &b118}
	var b123 = sequenceBuilder{id: 123, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b122 = charBuilder{}
	b123.items = []builder{&b122}
	b124.items = []builder{&b121, &b826, &b114, &b826, &b119, &b826, &b114, &b826, &b123}
	b125.items = []builder{&b124}
	var b130 = sequenceBuilder{id: 130, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b127 = sequenceBuilder{id: 127, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b126 = charBuilder{}
	b127.items = []builder{&b126}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b128.items = []builder{&b826, &b14}
	b129.items = []builder{&b826, &b14, &b128}
	b130.items = []builder{&b127, &b129, &b826, &b124}
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
	b135.items = []builder{&b826, &b14}
	b136.items = []builder{&b826, &b14, &b135}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b137.items = []builder{&b826, &b14}
	b138.items = []builder{&b826, &b14, &b137}
	var b134 = sequenceBuilder{id: 134, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b133 = charBuilder{}
	b134.items = []builder{&b133}
	b139.items = []builder{&b132, &b136, &b826, &b402, &b138, &b826, &b134}
	b140.options = []builder{&b104, &b87, &b139}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b143.items = []builder{&b826, &b14}
	b144.items = []builder{&b826, &b14, &b143}
	var b142 = sequenceBuilder{id: 142, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b141 = charBuilder{}
	b142.items = []builder{&b141}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b145.items = []builder{&b826, &b14}
	b146.items = []builder{&b826, &b14, &b145}
	b147.items = []builder{&b140, &b144, &b826, &b142, &b146, &b826, &b402}
	b148.options = []builder{&b147, &b109}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b149 = choiceBuilder{id: 149, commit: 2}
	b149.options = []builder{&b147, &b109}
	b150.items = []builder{&b114, &b826, &b149}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b151.items = []builder{&b826, &b150}
	b152.items = []builder{&b826, &b150, &b151}
	b153.items = []builder{&b148, &b152}
	var b157 = sequenceBuilder{id: 157, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b156 = charBuilder{}
	b157.items = []builder{&b156}
	b158.items = []builder{&b155, &b826, &b114, &b826, &b153, &b826, &b114, &b826, &b157}
	b159.items = []builder{&b158}
	var b164 = sequenceBuilder{id: 164, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b161 = sequenceBuilder{id: 161, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b160 = charBuilder{}
	b161.items = []builder{&b160}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b162.items = []builder{&b826, &b14}
	b163.items = []builder{&b826, &b14, &b162}
	b164.items = []builder{&b161, &b163, &b826, &b158}
	var b208 = sequenceBuilder{id: 208, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b205 = sequenceBuilder{id: 205, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b203 = charBuilder{}
	var b204 = charBuilder{}
	b205.items = []builder{&b203, &b204}
	var b207 = sequenceBuilder{id: 207, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b206 = sequenceBuilder{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b206.items = []builder{&b826, &b14}
	b207.items = []builder{&b826, &b14, &b206}
	var b202 = sequenceBuilder{id: 202, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b194 = sequenceBuilder{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b193 = charBuilder{}
	b194.items = []builder{&b193}
	var b196 = choiceBuilder{id: 196, commit: 2}
	var b168 = sequenceBuilder{id: 168, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b167 = sequenceBuilder{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b165.items = []builder{&b114, &b826, &b104}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b166.items = []builder{&b826, &b165}
	b167.items = []builder{&b826, &b165, &b166}
	b168.items = []builder{&b104, &b167}
	var b195 = sequenceBuilder{id: 195, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b175 = sequenceBuilder{id: 175, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b172 = sequenceBuilder{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b169 = charBuilder{}
	var b170 = charBuilder{}
	var b171 = charBuilder{}
	b172.items = []builder{&b169, &b170, &b171}
	var b174 = sequenceBuilder{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b173 = sequenceBuilder{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b173.items = []builder{&b826, &b14}
	b174.items = []builder{&b826, &b14, &b173}
	b175.items = []builder{&b172, &b174, &b826, &b104}
	b195.items = []builder{&b168, &b826, &b114, &b826, &b175}
	b196.options = []builder{&b168, &b195, &b175}
	var b198 = sequenceBuilder{id: 198, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b197 = charBuilder{}
	b198.items = []builder{&b197}
	var b201 = sequenceBuilder{id: 201, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b200 = sequenceBuilder{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b200.items = []builder{&b826, &b14}
	b201.items = []builder{&b826, &b14, &b200}
	var b199 = choiceBuilder{id: 199, commit: 2}
	var b784 = choiceBuilder{id: 784, commit: 66}
	var b512 = sequenceBuilder{id: 512, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b498 = sequenceBuilder{id: 498, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b497 = sequenceBuilder{id: 497, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b493 = charBuilder{}
	var b494 = charBuilder{}
	var b495 = charBuilder{}
	var b496 = charBuilder{}
	b497.items = []builder{&b493, &b494, &b495, &b496}
	b498.items = []builder{&b497, &b15}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b508 = sequenceBuilder{id: 508, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b508.items = []builder{&b826, &b14}
	b509.items = []builder{&b826, &b14, &b508}
	var b511 = sequenceBuilder{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b510.items = []builder{&b826, &b14}
	b511.items = []builder{&b826, &b14, &b510}
	b512.items = []builder{&b498, &b509, &b826, &b273, &b511, &b826, &b273}
	var b565 = sequenceBuilder{id: 565, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b555 = sequenceBuilder{id: 555, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b554 = sequenceBuilder{id: 554, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b552 = charBuilder{}
	var b553 = charBuilder{}
	b554.items = []builder{&b552, &b553}
	b555.items = []builder{&b554, &b15}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b563 = sequenceBuilder{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b563.items = []builder{&b826, &b14}
	b564.items = []builder{&b826, &b14, &b563}
	var b263 = sequenceBuilder{id: 263, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b260 = sequenceBuilder{id: 260, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b259 = charBuilder{}
	b260.items = []builder{&b259}
	var b262 = sequenceBuilder{id: 262, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b261 = charBuilder{}
	b262.items = []builder{&b261}
	b263.items = []builder{&b273, &b826, &b260, &b826, &b114, &b826, &b119, &b826, &b114, &b826, &b262}
	b565.items = []builder{&b555, &b564, &b826, &b263}
	var b574 = sequenceBuilder{id: 574, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b571 = sequenceBuilder{id: 571, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b566 = charBuilder{}
	var b567 = charBuilder{}
	var b568 = charBuilder{}
	var b569 = charBuilder{}
	var b570 = charBuilder{}
	b571.items = []builder{&b566, &b567, &b568, &b569, &b570}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b572.items = []builder{&b826, &b14}
	b573.items = []builder{&b826, &b14, &b572}
	b574.items = []builder{&b571, &b573, &b826, &b263}
	var b638 = choiceBuilder{id: 638, commit: 64, name: "assignment"}
	var b622 = sequenceBuilder{id: 622, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b607 = sequenceBuilder{id: 607, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b606 = sequenceBuilder{id: 606, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b603 = charBuilder{}
	var b604 = charBuilder{}
	var b605 = charBuilder{}
	b606.items = []builder{&b603, &b604, &b605}
	b607.items = []builder{&b606, &b15}
	var b621 = sequenceBuilder{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b620 = sequenceBuilder{id: 620, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b620.items = []builder{&b826, &b14}
	b621.items = []builder{&b826, &b14, &b620}
	var b615 = sequenceBuilder{id: 615, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b612 = sequenceBuilder{id: 612, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b611 = sequenceBuilder{id: 611, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b610 = sequenceBuilder{id: 610, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b610.items = []builder{&b826, &b14}
	b611.items = []builder{&b14, &b610}
	var b609 = sequenceBuilder{id: 609, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b608 = charBuilder{}
	b609.items = []builder{&b608}
	b612.items = []builder{&b611, &b826, &b609}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b613 = sequenceBuilder{id: 613, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b613.items = []builder{&b826, &b14}
	b614.items = []builder{&b826, &b14, &b613}
	b615.items = []builder{&b273, &b826, &b612, &b614, &b826, &b402}
	b622.items = []builder{&b607, &b621, &b826, &b615}
	var b629 = sequenceBuilder{id: 629, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b626 = sequenceBuilder{id: 626, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b625 = sequenceBuilder{id: 625, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b625.items = []builder{&b826, &b14}
	b626.items = []builder{&b826, &b14, &b625}
	var b624 = sequenceBuilder{id: 624, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b623 = charBuilder{}
	b624.items = []builder{&b623}
	var b628 = sequenceBuilder{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b627 = sequenceBuilder{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b627.items = []builder{&b826, &b14}
	b628.items = []builder{&b826, &b14, &b627}
	b629.items = []builder{&b273, &b626, &b826, &b624, &b628, &b826, &b402}
	var b637 = sequenceBuilder{id: 637, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b635 = sequenceBuilder{id: 635, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b635.items = []builder{&b826, &b14}
	b636.items = []builder{&b826, &b14, &b635}
	var b631 = sequenceBuilder{id: 631, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b630 = charBuilder{}
	b631.items = []builder{&b630}
	var b632 = sequenceBuilder{id: 632, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b619 = sequenceBuilder{id: 619, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b618 = sequenceBuilder{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b616 = sequenceBuilder{id: 616, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b616.items = []builder{&b114, &b826, &b615}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b617.items = []builder{&b826, &b616}
	b618.items = []builder{&b826, &b616, &b617}
	b619.items = []builder{&b615, &b618}
	b632.items = []builder{&b114, &b826, &b619}
	var b634 = sequenceBuilder{id: 634, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b633 = charBuilder{}
	b634.items = []builder{&b633}
	b637.items = []builder{&b607, &b636, &b826, &b631, &b826, &b632, &b826, &b114, &b826, &b634}
	b638.options = []builder{&b622, &b629, &b637}
	var b793 = sequenceBuilder{id: 793, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b786 = sequenceBuilder{id: 786, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b785 = charBuilder{}
	b786.items = []builder{&b785}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b789 = sequenceBuilder{id: 789, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b789.items = []builder{&b826, &b14}
	b790.items = []builder{&b826, &b14, &b789}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b791.items = []builder{&b826, &b14}
	b792.items = []builder{&b826, &b14, &b791}
	var b788 = sequenceBuilder{id: 788, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b787 = charBuilder{}
	b788.items = []builder{&b787}
	b793.items = []builder{&b786, &b790, &b826, &b784, &b792, &b826, &b788}
	b784.options = []builder{&b512, &b565, &b574, &b638, &b793, &b402}
	var b192 = sequenceBuilder{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	var b191 = sequenceBuilder{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b190 = charBuilder{}
	b191.items = []builder{&b190}
	b192.items = []builder{&b189, &b826, &b808, &b826, &b812, &b826, &b808, &b826, &b191}
	b199.options = []builder{&b784, &b192}
	b202.items = []builder{&b194, &b826, &b114, &b826, &b196, &b826, &b114, &b826, &b198, &b201, &b826, &b199}
	b208.items = []builder{&b205, &b207, &b826, &b202}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b211 = sequenceBuilder{id: 211, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b209 = charBuilder{}
	var b210 = charBuilder{}
	b211.items = []builder{&b209, &b210}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b214.items = []builder{&b826, &b14}
	b215.items = []builder{&b826, &b14, &b214}
	var b213 = sequenceBuilder{id: 213, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b212 = charBuilder{}
	b213.items = []builder{&b212}
	var b217 = sequenceBuilder{id: 217, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b216 = sequenceBuilder{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b216.items = []builder{&b826, &b14}
	b217.items = []builder{&b826, &b14, &b216}
	b218.items = []builder{&b211, &b215, &b826, &b213, &b217, &b826, &b202}
	var b258 = sequenceBuilder{id: 258, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b257 = sequenceBuilder{id: 257, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b256 = sequenceBuilder{id: 256, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b256.items = []builder{&b826, &b14}
	b257.items = []builder{&b826, &b14, &b256}
	var b255 = sequenceBuilder{id: 255, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b251 = choiceBuilder{id: 251, commit: 66}
	var b232 = sequenceBuilder{id: 232, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b229 = sequenceBuilder{id: 229, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b228 = charBuilder{}
	b229.items = []builder{&b228}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b230 = sequenceBuilder{id: 230, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b230.items = []builder{&b826, &b14}
	b231.items = []builder{&b826, &b14, &b230}
	b232.items = []builder{&b229, &b231, &b826, &b104}
	var b241 = sequenceBuilder{id: 241, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b234 = sequenceBuilder{id: 234, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b233 = charBuilder{}
	b234.items = []builder{&b233}
	var b238 = sequenceBuilder{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b237 = sequenceBuilder{id: 237, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b237.items = []builder{&b826, &b14}
	b238.items = []builder{&b826, &b14, &b237}
	var b240 = sequenceBuilder{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b239 = sequenceBuilder{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b239.items = []builder{&b826, &b14}
	b240.items = []builder{&b826, &b14, &b239}
	var b236 = sequenceBuilder{id: 236, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b235 = charBuilder{}
	b236.items = []builder{&b235}
	b241.items = []builder{&b234, &b238, &b826, &b402, &b240, &b826, &b236}
	var b250 = sequenceBuilder{id: 250, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b243 = sequenceBuilder{id: 243, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b242 = charBuilder{}
	b243.items = []builder{&b242}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b246 = sequenceBuilder{id: 246, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b246.items = []builder{&b826, &b14}
	b247.items = []builder{&b826, &b14, &b246}
	var b227 = sequenceBuilder{id: 227, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b219 = sequenceBuilder{id: 219, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b219.items = []builder{&b402}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b223.items = []builder{&b826, &b14}
	b224.items = []builder{&b826, &b14, &b223}
	var b222 = sequenceBuilder{id: 222, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b221 = charBuilder{}
	b222.items = []builder{&b221}
	var b226 = sequenceBuilder{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b225 = sequenceBuilder{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b225.items = []builder{&b826, &b14}
	b226.items = []builder{&b826, &b14, &b225}
	var b220 = sequenceBuilder{id: 220, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b220.items = []builder{&b402}
	b227.items = []builder{&b219, &b224, &b826, &b222, &b226, &b826, &b220}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b248.items = []builder{&b826, &b14}
	b249.items = []builder{&b826, &b14, &b248}
	var b245 = sequenceBuilder{id: 245, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b244 = charBuilder{}
	b245.items = []builder{&b244}
	b250.items = []builder{&b243, &b247, &b826, &b227, &b249, &b826, &b245}
	b251.options = []builder{&b232, &b241, &b250}
	var b254 = sequenceBuilder{id: 254, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b253 = sequenceBuilder{id: 253, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b252 = sequenceBuilder{id: 252, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b252.items = []builder{&b826, &b14}
	b253.items = []builder{&b14, &b252}
	b254.items = []builder{&b253, &b826, &b251}
	b255.items = []builder{&b251, &b826, &b254}
	b258.items = []builder{&b273, &b257, &b826, &b255}
	var b272 = sequenceBuilder{id: 272, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b265 = sequenceBuilder{id: 265, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b264 = charBuilder{}
	b265.items = []builder{&b264}
	var b269 = sequenceBuilder{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b268 = sequenceBuilder{id: 268, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b268.items = []builder{&b826, &b14}
	b269.items = []builder{&b826, &b14, &b268}
	var b271 = sequenceBuilder{id: 271, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b270 = sequenceBuilder{id: 270, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b270.items = []builder{&b826, &b14}
	b271.items = []builder{&b826, &b14, &b270}
	var b267 = sequenceBuilder{id: 267, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b266 = charBuilder{}
	b267.items = []builder{&b266}
	b272.items = []builder{&b265, &b269, &b826, &b402, &b271, &b826, &b267}
	b273.options = []builder{&b61, &b74, &b87, &b99, &b515, &b104, &b125, &b130, &b159, &b164, &b208, &b218, &b258, &b263, &b272}
	var b333 = sequenceBuilder{id: 333, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b332 = choiceBuilder{id: 332, commit: 66}
	var b292 = sequenceBuilder{id: 292, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b291 = charBuilder{}
	b292.items = []builder{&b291}
	var b294 = sequenceBuilder{id: 294, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b293 = charBuilder{}
	b294.items = []builder{&b293}
	var b275 = sequenceBuilder{id: 275, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b274 = charBuilder{}
	b275.items = []builder{&b274}
	var b306 = sequenceBuilder{id: 306, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b305 = charBuilder{}
	b306.items = []builder{&b305}
	b332.options = []builder{&b292, &b294, &b275, &b306}
	b333.items = []builder{&b332, &b826, &b273}
	var b380 = choiceBuilder{id: 380, commit: 66}
	var b351 = sequenceBuilder{id: 351, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b339 = choiceBuilder{id: 339, commit: 66}
	b339.options = []builder{&b273, &b333}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b345 = sequenceBuilder{id: 345, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b345.items = []builder{&b826, &b14}
	b346.items = []builder{&b14, &b345}
	var b334 = choiceBuilder{id: 334, commit: 66}
	var b277 = sequenceBuilder{id: 277, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b276 = charBuilder{}
	b277.items = []builder{&b276}
	var b284 = sequenceBuilder{id: 284, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b282 = charBuilder{}
	var b283 = charBuilder{}
	b284.items = []builder{&b282, &b283}
	var b287 = sequenceBuilder{id: 287, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b285 = charBuilder{}
	var b286 = charBuilder{}
	b287.items = []builder{&b285, &b286}
	var b290 = sequenceBuilder{id: 290, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b288 = charBuilder{}
	var b289 = charBuilder{}
	b290.items = []builder{&b288, &b289}
	var b296 = sequenceBuilder{id: 296, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b295 = charBuilder{}
	b296.items = []builder{&b295}
	var b298 = sequenceBuilder{id: 298, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b297 = charBuilder{}
	b298.items = []builder{&b297}
	var b300 = sequenceBuilder{id: 300, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b299 = charBuilder{}
	b300.items = []builder{&b299}
	b334.options = []builder{&b277, &b284, &b287, &b290, &b296, &b298, &b300}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b347.items = []builder{&b826, &b14}
	b348.items = []builder{&b826, &b14, &b347}
	b349.items = []builder{&b346, &b826, &b334, &b348, &b826, &b339}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b350.items = []builder{&b826, &b349}
	b351.items = []builder{&b339, &b826, &b349, &b350}
	var b358 = sequenceBuilder{id: 358, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b340 = choiceBuilder{id: 340, commit: 66}
	b340.options = []builder{&b339, &b351}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b352.items = []builder{&b826, &b14}
	b353.items = []builder{&b14, &b352}
	var b335 = choiceBuilder{id: 335, commit: 66}
	var b279 = sequenceBuilder{id: 279, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b278 = charBuilder{}
	b279.items = []builder{&b278}
	var b281 = sequenceBuilder{id: 281, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b280 = charBuilder{}
	b281.items = []builder{&b280}
	var b302 = sequenceBuilder{id: 302, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b301 = charBuilder{}
	b302.items = []builder{&b301}
	var b304 = sequenceBuilder{id: 304, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b303 = charBuilder{}
	b304.items = []builder{&b303}
	b335.options = []builder{&b279, &b281, &b302, &b304}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b354.items = []builder{&b826, &b14}
	b355.items = []builder{&b826, &b14, &b354}
	b356.items = []builder{&b353, &b826, &b335, &b355, &b826, &b340}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b357.items = []builder{&b826, &b356}
	b358.items = []builder{&b340, &b826, &b356, &b357}
	var b365 = sequenceBuilder{id: 365, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b341 = choiceBuilder{id: 341, commit: 66}
	b341.options = []builder{&b340, &b358}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b826, &b14}
	b360.items = []builder{&b14, &b359}
	var b336 = choiceBuilder{id: 336, commit: 66}
	var b309 = sequenceBuilder{id: 309, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b307 = charBuilder{}
	var b308 = charBuilder{}
	b309.items = []builder{&b307, &b308}
	var b312 = sequenceBuilder{id: 312, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b310 = charBuilder{}
	var b311 = charBuilder{}
	b312.items = []builder{&b310, &b311}
	var b314 = sequenceBuilder{id: 314, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b313 = charBuilder{}
	b314.items = []builder{&b313}
	var b317 = sequenceBuilder{id: 317, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b315 = charBuilder{}
	var b316 = charBuilder{}
	b317.items = []builder{&b315, &b316}
	var b319 = sequenceBuilder{id: 319, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b318 = charBuilder{}
	b319.items = []builder{&b318}
	var b322 = sequenceBuilder{id: 322, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b320 = charBuilder{}
	var b321 = charBuilder{}
	b322.items = []builder{&b320, &b321}
	b336.options = []builder{&b309, &b312, &b314, &b317, &b319, &b322}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b361.items = []builder{&b826, &b14}
	b362.items = []builder{&b826, &b14, &b361}
	b363.items = []builder{&b360, &b826, &b336, &b362, &b826, &b341}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b364.items = []builder{&b826, &b363}
	b365.items = []builder{&b341, &b826, &b363, &b364}
	var b372 = sequenceBuilder{id: 372, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b342 = choiceBuilder{id: 342, commit: 66}
	b342.options = []builder{&b341, &b365}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b366.items = []builder{&b826, &b14}
	b367.items = []builder{&b14, &b366}
	var b337 = sequenceBuilder{id: 337, commit: 66, ranges: [][]int{{1, 1}}}
	var b325 = sequenceBuilder{id: 325, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b323 = charBuilder{}
	var b324 = charBuilder{}
	b325.items = []builder{&b323, &b324}
	b337.items = []builder{&b325}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b368.items = []builder{&b826, &b14}
	b369.items = []builder{&b826, &b14, &b368}
	b370.items = []builder{&b367, &b826, &b337, &b369, &b826, &b342}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b371.items = []builder{&b826, &b370}
	b372.items = []builder{&b342, &b826, &b370, &b371}
	var b379 = sequenceBuilder{id: 379, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b343 = choiceBuilder{id: 343, commit: 66}
	b343.options = []builder{&b342, &b372}
	var b377 = sequenceBuilder{id: 377, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b374 = sequenceBuilder{id: 374, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b373.items = []builder{&b826, &b14}
	b374.items = []builder{&b14, &b373}
	var b338 = sequenceBuilder{id: 338, commit: 66, ranges: [][]int{{1, 1}}}
	var b328 = sequenceBuilder{id: 328, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b326 = charBuilder{}
	var b327 = charBuilder{}
	b328.items = []builder{&b326, &b327}
	b338.items = []builder{&b328}
	var b376 = sequenceBuilder{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b375 = sequenceBuilder{id: 375, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b375.items = []builder{&b826, &b14}
	b376.items = []builder{&b826, &b14, &b375}
	b377.items = []builder{&b374, &b826, &b338, &b376, &b826, &b343}
	var b378 = sequenceBuilder{id: 378, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b378.items = []builder{&b826, &b377}
	b379.items = []builder{&b343, &b826, &b377, &b378}
	b380.options = []builder{&b351, &b358, &b365, &b372, &b379}
	var b393 = sequenceBuilder{id: 393, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b385 = sequenceBuilder{id: 385, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b385.items = []builder{&b826, &b14}
	b386.items = []builder{&b826, &b14, &b385}
	var b382 = sequenceBuilder{id: 382, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b381 = charBuilder{}
	b382.items = []builder{&b381}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b387.items = []builder{&b826, &b14}
	b388.items = []builder{&b826, &b14, &b387}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b389.items = []builder{&b826, &b14}
	b390.items = []builder{&b826, &b14, &b389}
	var b384 = sequenceBuilder{id: 384, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b383 = charBuilder{}
	b384.items = []builder{&b383}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b391.items = []builder{&b826, &b14}
	b392.items = []builder{&b826, &b14, &b391}
	b393.items = []builder{&b402, &b386, &b826, &b382, &b388, &b826, &b402, &b390, &b826, &b384, &b392, &b826, &b402}
	var b401 = sequenceBuilder{id: 401, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b394 = choiceBuilder{id: 394, commit: 66}
	b394.options = []builder{&b273, &b333, &b380, &b393}
	var b399 = sequenceBuilder{id: 399, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b396 = sequenceBuilder{id: 396, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b395 = sequenceBuilder{id: 395, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b395.items = []builder{&b826, &b14}
	b396.items = []builder{&b14, &b395}
	var b331 = sequenceBuilder{id: 331, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b329 = charBuilder{}
	var b330 = charBuilder{}
	b331.items = []builder{&b329, &b330}
	var b398 = sequenceBuilder{id: 398, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b397 = sequenceBuilder{id: 397, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b397.items = []builder{&b826, &b14}
	b398.items = []builder{&b826, &b14, &b397}
	b399.items = []builder{&b396, &b826, &b331, &b398, &b826, &b394}
	var b400 = sequenceBuilder{id: 400, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b400.items = []builder{&b826, &b399}
	b401.items = []builder{&b394, &b826, &b399, &b400}
	b402.options = []builder{&b273, &b333, &b380, &b393, &b401}
	b186.items = []builder{&b185, &b826, &b402}
	b187.items = []builder{&b183, &b826, &b186}
	var b433 = sequenceBuilder{id: 433, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b406 = sequenceBuilder{id: 406, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b405 = sequenceBuilder{id: 405, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b403 = charBuilder{}
	var b404 = charBuilder{}
	b405.items = []builder{&b403, &b404}
	b406.items = []builder{&b405, &b15}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b427.items = []builder{&b826, &b14}
	b428.items = []builder{&b826, &b14, &b427}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b429.items = []builder{&b826, &b14}
	b430.items = []builder{&b826, &b14, &b429}
	var b432 = sequenceBuilder{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b421 = sequenceBuilder{id: 421, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b413 = sequenceBuilder{id: 413, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b413.items = []builder{&b826, &b14}
	b414.items = []builder{&b14, &b413}
	var b412 = sequenceBuilder{id: 412, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b411 = sequenceBuilder{id: 411, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b407 = charBuilder{}
	var b408 = charBuilder{}
	var b409 = charBuilder{}
	var b410 = charBuilder{}
	b411.items = []builder{&b407, &b408, &b409, &b410}
	b412.items = []builder{&b411, &b15}
	var b416 = sequenceBuilder{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b415.items = []builder{&b826, &b14}
	b416.items = []builder{&b826, &b14, &b415}
	var b418 = sequenceBuilder{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b417 = sequenceBuilder{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b417.items = []builder{&b826, &b14}
	b418.items = []builder{&b826, &b14, &b417}
	var b420 = sequenceBuilder{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b419 = sequenceBuilder{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b419.items = []builder{&b826, &b14}
	b420.items = []builder{&b826, &b14, &b419}
	b421.items = []builder{&b414, &b826, &b412, &b416, &b826, &b406, &b418, &b826, &b402, &b420, &b826, &b192}
	var b431 = sequenceBuilder{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b431.items = []builder{&b826, &b421}
	b432.items = []builder{&b826, &b421, &b431}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b422.items = []builder{&b826, &b14}
	b423.items = []builder{&b14, &b422}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b424.items = []builder{&b826, &b14}
	b425.items = []builder{&b826, &b14, &b424}
	b426.items = []builder{&b423, &b826, &b412, &b425, &b826, &b192}
	b433.items = []builder{&b406, &b428, &b826, &b402, &b430, &b826, &b192, &b432, &b826, &b426}
	var b492 = sequenceBuilder{id: 492, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b447 = sequenceBuilder{id: 447, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b446 = sequenceBuilder{id: 446, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b440 = charBuilder{}
	var b441 = charBuilder{}
	var b442 = charBuilder{}
	var b443 = charBuilder{}
	var b444 = charBuilder{}
	var b445 = charBuilder{}
	b446.items = []builder{&b440, &b441, &b442, &b443, &b444, &b445}
	b447.items = []builder{&b446, &b15}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b488 = sequenceBuilder{id: 488, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b488.items = []builder{&b826, &b14}
	b489.items = []builder{&b826, &b14, &b488}
	var b491 = sequenceBuilder{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b490 = sequenceBuilder{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b490.items = []builder{&b826, &b14}
	b491.items = []builder{&b826, &b14, &b490}
	var b479 = sequenceBuilder{id: 479, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b478 = charBuilder{}
	b479.items = []builder{&b478}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b480 = choiceBuilder{id: 480, commit: 2}
	var b477 = sequenceBuilder{id: 477, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b472 = sequenceBuilder{id: 472, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b439 = sequenceBuilder{id: 439, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b438 = sequenceBuilder{id: 438, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b434 = charBuilder{}
	var b435 = charBuilder{}
	var b436 = charBuilder{}
	var b437 = charBuilder{}
	b438.items = []builder{&b434, &b435, &b436, &b437}
	b439.items = []builder{&b438, &b15}
	var b469 = sequenceBuilder{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b468 = sequenceBuilder{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b468.items = []builder{&b826, &b14}
	b469.items = []builder{&b826, &b14, &b468}
	var b471 = sequenceBuilder{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b470 = sequenceBuilder{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b470.items = []builder{&b826, &b14}
	b471.items = []builder{&b826, &b14, &b470}
	var b467 = sequenceBuilder{id: 467, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b466 = charBuilder{}
	b467.items = []builder{&b466}
	b472.items = []builder{&b439, &b469, &b826, &b402, &b471, &b826, &b467}
	var b476 = sequenceBuilder{id: 476, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b474 = sequenceBuilder{id: 474, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b473 = charBuilder{}
	b474.items = []builder{&b473}
	var b475 = sequenceBuilder{id: 475, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b475.items = []builder{&b826, &b474}
	b476.items = []builder{&b826, &b474, &b475}
	b477.items = []builder{&b472, &b476, &b826, &b794}
	var b465 = sequenceBuilder{id: 465, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b460 = sequenceBuilder{id: 460, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b455 = sequenceBuilder{id: 455, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b448 = charBuilder{}
	var b449 = charBuilder{}
	var b450 = charBuilder{}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	var b454 = charBuilder{}
	b455.items = []builder{&b448, &b449, &b450, &b451, &b452, &b453, &b454}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b458 = sequenceBuilder{id: 458, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b458.items = []builder{&b826, &b14}
	b459.items = []builder{&b826, &b14, &b458}
	var b457 = sequenceBuilder{id: 457, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b456 = charBuilder{}
	b457.items = []builder{&b456}
	b460.items = []builder{&b455, &b459, &b826, &b457}
	var b464 = sequenceBuilder{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b462 = sequenceBuilder{id: 462, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b461 = charBuilder{}
	b462.items = []builder{&b461}
	var b463 = sequenceBuilder{id: 463, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b463.items = []builder{&b826, &b462}
	b464.items = []builder{&b826, &b462, &b463}
	b465.items = []builder{&b460, &b464, &b826, &b794}
	b480.options = []builder{&b477, &b465}
	var b484 = sequenceBuilder{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b482 = sequenceBuilder{id: 482, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b481 = choiceBuilder{id: 481, commit: 2}
	b481.options = []builder{&b477, &b465, &b794}
	b482.items = []builder{&b808, &b826, &b481}
	var b483 = sequenceBuilder{id: 483, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b483.items = []builder{&b826, &b482}
	b484.items = []builder{&b826, &b482, &b483}
	b485.items = []builder{&b480, &b484}
	var b487 = sequenceBuilder{id: 487, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b486 = charBuilder{}
	b487.items = []builder{&b486}
	b492.items = []builder{&b447, &b489, &b826, &b402, &b491, &b826, &b479, &b826, &b808, &b826, &b485, &b826, &b808, &b826, &b487}
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
	b549.items = []builder{&b826, &b14}
	b550.items = []builder{&b826, &b14, &b549}
	var b540 = sequenceBuilder{id: 540, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b539 = charBuilder{}
	b540.items = []builder{&b539}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b541 = choiceBuilder{id: 541, commit: 2}
	var b531 = sequenceBuilder{id: 531, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b526 = sequenceBuilder{id: 526, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b522 = sequenceBuilder{id: 522, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b522.items = []builder{&b826, &b14}
	b523.items = []builder{&b826, &b14, &b522}
	var b519 = choiceBuilder{id: 519, commit: 66}
	var b518 = sequenceBuilder{id: 518, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b517 = sequenceBuilder{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b516 = sequenceBuilder{id: 516, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b516.items = []builder{&b826, &b14}
	b517.items = []builder{&b826, &b14, &b516}
	b518.items = []builder{&b104, &b517, &b826, &b515}
	b519.options = []builder{&b512, &b515, &b518}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b524.items = []builder{&b826, &b14}
	b525.items = []builder{&b826, &b14, &b524}
	var b521 = sequenceBuilder{id: 521, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b520 = charBuilder{}
	b521.items = []builder{&b520}
	b526.items = []builder{&b439, &b523, &b826, &b519, &b525, &b826, &b521}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b528 = sequenceBuilder{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b527 = charBuilder{}
	b528.items = []builder{&b527}
	var b529 = sequenceBuilder{id: 529, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b529.items = []builder{&b826, &b528}
	b530.items = []builder{&b826, &b528, &b529}
	b531.items = []builder{&b526, &b530, &b826, &b794}
	b541.options = []builder{&b531, &b465}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b543 = sequenceBuilder{id: 543, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b542 = choiceBuilder{id: 542, commit: 2}
	b542.options = []builder{&b531, &b465, &b794}
	b543.items = []builder{&b808, &b826, &b542}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b544.items = []builder{&b826, &b543}
	b545.items = []builder{&b826, &b543, &b544}
	b546.items = []builder{&b541, &b545}
	var b548 = sequenceBuilder{id: 548, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b547 = charBuilder{}
	b548.items = []builder{&b547}
	b551.items = []builder{&b538, &b550, &b826, &b540, &b826, &b808, &b826, &b546, &b826, &b808, &b826, &b548}
	var b602 = sequenceBuilder{id: 602, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b583 = sequenceBuilder{id: 583, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b582 = sequenceBuilder{id: 582, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b579 = charBuilder{}
	var b580 = charBuilder{}
	var b581 = charBuilder{}
	b582.items = []builder{&b579, &b580, &b581}
	b583.items = []builder{&b582, &b15}
	var b601 = choiceBuilder{id: 601, commit: 2}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b592 = sequenceBuilder{id: 592, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b592.items = []builder{&b826, &b14}
	b593.items = []builder{&b14, &b592}
	var b591 = choiceBuilder{id: 591, commit: 66}
	var b590 = choiceBuilder{id: 590, commit: 64, name: "range-over-expression"}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b585.items = []builder{&b826, &b14}
	b586.items = []builder{&b826, &b14, &b585}
	var b578 = sequenceBuilder{id: 578, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b577 = sequenceBuilder{id: 577, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b575 = charBuilder{}
	var b576 = charBuilder{}
	b577.items = []builder{&b575, &b576}
	b578.items = []builder{&b577, &b15}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b587.items = []builder{&b826, &b14}
	b588.items = []builder{&b826, &b14, &b587}
	var b584 = choiceBuilder{id: 584, commit: 2}
	b584.options = []builder{&b402, &b227}
	b589.items = []builder{&b104, &b586, &b826, &b578, &b588, &b826, &b584}
	b590.options = []builder{&b589, &b227}
	b591.options = []builder{&b402, &b590}
	b594.items = []builder{&b593, &b826, &b591}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b595.items = []builder{&b826, &b14}
	b596.items = []builder{&b826, &b14, &b595}
	b597.items = []builder{&b594, &b596, &b826, &b192}
	var b600 = sequenceBuilder{id: 600, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b599 = sequenceBuilder{id: 599, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b598.items = []builder{&b826, &b14}
	b599.items = []builder{&b14, &b598}
	b600.items = []builder{&b599, &b826, &b192}
	b601.options = []builder{&b597, &b600}
	b602.items = []builder{&b583, &b826, &b601}
	var b740 = choiceBuilder{id: 740, commit: 66}
	var b661 = sequenceBuilder{id: 661, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b643 = sequenceBuilder{id: 643, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b642 = sequenceBuilder{id: 642, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b639 = charBuilder{}
	var b640 = charBuilder{}
	var b641 = charBuilder{}
	b642.items = []builder{&b639, &b640, &b641}
	b643.items = []builder{&b642, &b15}
	var b660 = sequenceBuilder{id: 660, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b659 = sequenceBuilder{id: 659, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b659.items = []builder{&b826, &b14}
	b660.items = []builder{&b826, &b14, &b659}
	var b658 = choiceBuilder{id: 658, commit: 2}
	var b652 = sequenceBuilder{id: 652, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b651 = sequenceBuilder{id: 651, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b648 = sequenceBuilder{id: 648, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b647 = sequenceBuilder{id: 647, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b646 = sequenceBuilder{id: 646, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b646.items = []builder{&b826, &b14}
	b647.items = []builder{&b14, &b646}
	var b645 = sequenceBuilder{id: 645, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b644 = charBuilder{}
	b645.items = []builder{&b644}
	b648.items = []builder{&b647, &b826, &b645}
	var b650 = sequenceBuilder{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b649 = sequenceBuilder{id: 649, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b649.items = []builder{&b826, &b14}
	b650.items = []builder{&b826, &b14, &b649}
	b651.items = []builder{&b104, &b826, &b648, &b650, &b826, &b402}
	b652.items = []builder{&b651}
	var b657 = sequenceBuilder{id: 657, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b654 = sequenceBuilder{id: 654, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b653 = charBuilder{}
	b654.items = []builder{&b653}
	var b656 = sequenceBuilder{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b655 = sequenceBuilder{id: 655, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b655.items = []builder{&b826, &b14}
	b656.items = []builder{&b826, &b14, &b655}
	b657.items = []builder{&b654, &b656, &b826, &b651}
	b658.options = []builder{&b652, &b657}
	b661.items = []builder{&b643, &b660, &b826, &b658}
	var b678 = sequenceBuilder{id: 678, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b677 = sequenceBuilder{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b676 = sequenceBuilder{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b676.items = []builder{&b826, &b14}
	b677.items = []builder{&b826, &b14, &b676}
	var b673 = sequenceBuilder{id: 673, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b672 = charBuilder{}
	b673.items = []builder{&b672}
	var b667 = sequenceBuilder{id: 667, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b662 = choiceBuilder{id: 662, commit: 2}
	b662.options = []builder{&b652, &b657}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b664 = sequenceBuilder{id: 664, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b663 = choiceBuilder{id: 663, commit: 2}
	b663.options = []builder{&b652, &b657}
	b664.items = []builder{&b114, &b826, &b663}
	var b665 = sequenceBuilder{id: 665, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b665.items = []builder{&b826, &b664}
	b666.items = []builder{&b826, &b664, &b665}
	b667.items = []builder{&b662, &b666}
	var b675 = sequenceBuilder{id: 675, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b674 = charBuilder{}
	b675.items = []builder{&b674}
	b678.items = []builder{&b643, &b677, &b826, &b673, &b826, &b114, &b826, &b667, &b826, &b114, &b826, &b675}
	var b689 = sequenceBuilder{id: 689, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b685 = sequenceBuilder{id: 685, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b685.items = []builder{&b826, &b14}
	b686.items = []builder{&b826, &b14, &b685}
	var b680 = sequenceBuilder{id: 680, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b679 = charBuilder{}
	b680.items = []builder{&b679}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b687.items = []builder{&b826, &b14}
	b688.items = []builder{&b826, &b14, &b687}
	var b682 = sequenceBuilder{id: 682, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b681 = charBuilder{}
	b682.items = []builder{&b681}
	var b671 = sequenceBuilder{id: 671, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b670 = sequenceBuilder{id: 670, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b668 = sequenceBuilder{id: 668, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b668.items = []builder{&b114, &b826, &b652}
	var b669 = sequenceBuilder{id: 669, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b669.items = []builder{&b826, &b668}
	b670.items = []builder{&b826, &b668, &b669}
	b671.items = []builder{&b652, &b670}
	var b684 = sequenceBuilder{id: 684, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b683 = charBuilder{}
	b684.items = []builder{&b683}
	b689.items = []builder{&b643, &b686, &b826, &b680, &b688, &b826, &b682, &b826, &b114, &b826, &b671, &b826, &b114, &b826, &b684}
	var b705 = sequenceBuilder{id: 705, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b701 = sequenceBuilder{id: 701, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b699 = charBuilder{}
	var b700 = charBuilder{}
	b701.items = []builder{&b699, &b700}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b703.items = []builder{&b826, &b14}
	b704.items = []builder{&b826, &b14, &b703}
	var b702 = choiceBuilder{id: 702, commit: 2}
	var b693 = sequenceBuilder{id: 693, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b692 = sequenceBuilder{id: 692, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b691 = sequenceBuilder{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b690 = sequenceBuilder{id: 690, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b690.items = []builder{&b826, &b14}
	b691.items = []builder{&b826, &b14, &b690}
	b692.items = []builder{&b104, &b691, &b826, &b202}
	b693.items = []builder{&b692}
	var b698 = sequenceBuilder{id: 698, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b695 = sequenceBuilder{id: 695, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b694 = charBuilder{}
	b695.items = []builder{&b694}
	var b697 = sequenceBuilder{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b696.items = []builder{&b826, &b14}
	b697.items = []builder{&b826, &b14, &b696}
	b698.items = []builder{&b695, &b697, &b826, &b692}
	b702.options = []builder{&b693, &b698}
	b705.items = []builder{&b701, &b704, &b826, &b702}
	var b725 = sequenceBuilder{id: 725, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b718 = sequenceBuilder{id: 718, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b716 = charBuilder{}
	var b717 = charBuilder{}
	b718.items = []builder{&b716, &b717}
	var b724 = sequenceBuilder{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b723 = sequenceBuilder{id: 723, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b723.items = []builder{&b826, &b14}
	b724.items = []builder{&b826, &b14, &b723}
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
	b712.items = []builder{&b114, &b826, &b711}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b713.items = []builder{&b826, &b712}
	b714.items = []builder{&b826, &b712, &b713}
	b715.items = []builder{&b710, &b714}
	var b722 = sequenceBuilder{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b721 = charBuilder{}
	b722.items = []builder{&b721}
	b725.items = []builder{&b718, &b724, &b826, &b720, &b826, &b114, &b826, &b715, &b826, &b114, &b826, &b722}
	var b739 = sequenceBuilder{id: 739, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b728 = sequenceBuilder{id: 728, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b726 = charBuilder{}
	var b727 = charBuilder{}
	b728.items = []builder{&b726, &b727}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b735 = sequenceBuilder{id: 735, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b735.items = []builder{&b826, &b14}
	b736.items = []builder{&b826, &b14, &b735}
	var b730 = sequenceBuilder{id: 730, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b729 = charBuilder{}
	b730.items = []builder{&b729}
	var b738 = sequenceBuilder{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b737.items = []builder{&b826, &b14}
	b738.items = []builder{&b826, &b14, &b737}
	var b732 = sequenceBuilder{id: 732, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b731 = charBuilder{}
	b732.items = []builder{&b731}
	var b709 = sequenceBuilder{id: 709, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b708 = sequenceBuilder{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b706 = sequenceBuilder{id: 706, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b706.items = []builder{&b114, &b826, &b693}
	var b707 = sequenceBuilder{id: 707, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b707.items = []builder{&b826, &b706}
	b708.items = []builder{&b826, &b706, &b707}
	b709.items = []builder{&b693, &b708}
	var b734 = sequenceBuilder{id: 734, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b733 = charBuilder{}
	b734.items = []builder{&b733}
	b739.items = []builder{&b728, &b736, &b826, &b730, &b738, &b826, &b732, &b826, &b114, &b826, &b709, &b826, &b114, &b826, &b734}
	b740.options = []builder{&b661, &b678, &b689, &b705, &b725, &b739}
	var b772 = choiceBuilder{id: 772, commit: 64, name: "use"}
	var b764 = sequenceBuilder{id: 764, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b745 = sequenceBuilder{id: 745, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b744 = sequenceBuilder{id: 744, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b741 = charBuilder{}
	var b742 = charBuilder{}
	var b743 = charBuilder{}
	b744.items = []builder{&b741, &b742, &b743}
	b745.items = []builder{&b744, &b15}
	var b763 = sequenceBuilder{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b762 = sequenceBuilder{id: 762, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b762.items = []builder{&b826, &b14}
	b763.items = []builder{&b826, &b14, &b762}
	var b757 = choiceBuilder{id: 757, commit: 64, name: "use-fact"}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b748 = choiceBuilder{id: 748, commit: 2}
	var b747 = sequenceBuilder{id: 747, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b746 = charBuilder{}
	b747.items = []builder{&b746}
	b748.options = []builder{&b104, &b747}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b751 = sequenceBuilder{id: 751, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b751.items = []builder{&b826, &b14}
	b752.items = []builder{&b14, &b751}
	var b750 = sequenceBuilder{id: 750, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b749 = charBuilder{}
	b750.items = []builder{&b749}
	b753.items = []builder{&b752, &b826, &b750}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b754.items = []builder{&b826, &b14}
	b755.items = []builder{&b826, &b14, &b754}
	b756.items = []builder{&b748, &b826, &b753, &b755, &b826, &b87}
	b757.options = []builder{&b87, &b756}
	b764.items = []builder{&b745, &b763, &b826, &b757}
	var b771 = sequenceBuilder{id: 771, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b770 = sequenceBuilder{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b769 = sequenceBuilder{id: 769, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b769.items = []builder{&b826, &b14}
	b770.items = []builder{&b826, &b14, &b769}
	var b766 = sequenceBuilder{id: 766, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b765 = charBuilder{}
	b766.items = []builder{&b765}
	var b761 = sequenceBuilder{id: 761, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b760 = sequenceBuilder{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b758 = sequenceBuilder{id: 758, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b758.items = []builder{&b114, &b826, &b757}
	var b759 = sequenceBuilder{id: 759, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b759.items = []builder{&b826, &b758}
	b760.items = []builder{&b826, &b758, &b759}
	b761.items = []builder{&b757, &b760}
	var b768 = sequenceBuilder{id: 768, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b767 = charBuilder{}
	b768.items = []builder{&b767}
	b771.items = []builder{&b745, &b770, &b826, &b766, &b826, &b114, &b826, &b761, &b826, &b114, &b826, &b768}
	b772.options = []builder{&b764, &b771}
	var b783 = sequenceBuilder{id: 783, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b780 = sequenceBuilder{id: 780, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b779 = sequenceBuilder{id: 779, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b773 = charBuilder{}
	var b774 = charBuilder{}
	var b775 = charBuilder{}
	var b776 = charBuilder{}
	var b777 = charBuilder{}
	var b778 = charBuilder{}
	b779.items = []builder{&b773, &b774, &b775, &b776, &b777, &b778}
	b780.items = []builder{&b779, &b15}
	var b782 = sequenceBuilder{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b781.items = []builder{&b826, &b14}
	b782.items = []builder{&b826, &b14, &b781}
	b783.items = []builder{&b780, &b782, &b826, &b740}
	var b803 = sequenceBuilder{id: 803, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b796 = sequenceBuilder{id: 796, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b795 = charBuilder{}
	b796.items = []builder{&b795}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b799 = sequenceBuilder{id: 799, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b799.items = []builder{&b826, &b14}
	b800.items = []builder{&b826, &b14, &b799}
	var b802 = sequenceBuilder{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b801.items = []builder{&b826, &b14}
	b802.items = []builder{&b826, &b14, &b801}
	var b798 = sequenceBuilder{id: 798, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b797 = charBuilder{}
	b798.items = []builder{&b797}
	b803.items = []builder{&b796, &b800, &b826, &b794, &b802, &b826, &b798}
	b794.options = []builder{&b187, &b433, &b492, &b551, &b602, &b740, &b772, &b783, &b803, &b784}
	var b811 = sequenceBuilder{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b809 = sequenceBuilder{id: 809, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b809.items = []builder{&b808, &b826, &b794}
	var b810 = sequenceBuilder{id: 810, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b810.items = []builder{&b826, &b809}
	b811.items = []builder{&b826, &b809, &b810}
	b812.items = []builder{&b794, &b811}
	b827.items = []builder{&b823, &b826, &b808, &b826, &b812, &b826, &b808}
	b828.items = []builder{&b826, &b827, &b826}

	return parseInput(r, &p828, &b828)
}
