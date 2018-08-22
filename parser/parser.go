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
	var p824 = choiceParser{id: 824, commit: 70, name: "ws", generalizations: []int{826}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{824, 826}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p824.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p825 = sequenceParser{id: 825, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{826}}
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
	p41.items = []parser{&p40, &p826, &p38}
	var p42 = sequenceParser{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p42.items = []parser{&p826, &p41}
	p43.items = []parser{&p826, &p41, &p42}
	p44.items = []parser{&p38, &p43}
	p825.items = []parser{&p44}
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
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{806, 113}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p806.options = []parser{&p805, &p14}
	var p807 = sequenceParser{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p807.items = []parser{&p826, &p806}
	p808.items = []parser{&p806, &p807}
	var p812 = sequenceParser{id: 812, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p794 = choiceParser{id: 794, commit: 66, name: "statement", generalizations: []int{458, 532}}
	var p187 = sequenceParser{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{794, 458, 532}}
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
	p184.items = []parser{&p826, &p14}
	p185.items = []parser{&p14, &p184}
	var p375 = choiceParser{id: 375, commit: 66, name: "expression", generalizations: []int{116, 784, 198, 567, 560, 794}}
	var p267 = choiceParser{id: 267, commit: 66, name: "primary-expression", generalizations: []int{116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p62 = choiceParser{id: 62, commit: 64, name: "int", generalizations: []int{267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p53 = sequenceParser{id: 53, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p52 = sequenceParser{id: 52, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p51 = charParser{id: 51, ranges: [][]rune{{49, 57}}}
	p52.items = []parser{&p51}
	var p46 = sequenceParser{id: 46, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p45 = charParser{id: 45, ranges: [][]rune{{48, 57}}}
	p46.items = []parser{&p45}
	p53.items = []parser{&p52, &p46}
	var p56 = sequenceParser{id: 56, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{62, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p55 = sequenceParser{id: 55, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p54 = charParser{id: 54, chars: []rune{48}}
	p55.items = []parser{&p54}
	var p48 = sequenceParser{id: 48, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p47 = charParser{id: 47, ranges: [][]rune{{48, 55}}}
	p48.items = []parser{&p47}
	p56.items = []parser{&p55, &p48}
	var p61 = sequenceParser{id: 61, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{62, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
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
	var p75 = choiceParser{id: 75, commit: 72, name: "float", generalizations: []int{267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p70 = sequenceParser{id: 70, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{75, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
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
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{75, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p72 = sequenceParser{id: 72, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p71 = charParser{id: 71, chars: []rune{46}}
	p72.items = []parser{&p71}
	p73.items = []parser{&p72, &p46, &p67}
	var p74 = sequenceParser{id: 74, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{75, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	p74.items = []parser{&p46, &p67}
	p75.options = []parser{&p70, &p73, &p74}
	var p88 = sequenceParser{id: 88, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 116, 141, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 742, 794}}
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
	var p100 = choiceParser{id: 100, commit: 66, name: "bool", generalizations: []int{267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p93 = sequenceParser{id: 93, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p89 = charParser{id: 89, chars: []rune{116}}
	var p90 = charParser{id: 90, chars: []rune{114}}
	var p91 = charParser{id: 91, chars: []rune{117}}
	var p92 = charParser{id: 92, chars: []rune{101}}
	p93.items = []parser{&p89, &p90, &p91, &p92}
	var p99 = sequenceParser{id: 99, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{100, 267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p94 = charParser{id: 94, chars: []rune{102}}
	var p95 = charParser{id: 95, chars: []rune{97}}
	var p96 = charParser{id: 96, chars: []rune{108}}
	var p97 = charParser{id: 97, chars: []rune{115}}
	var p98 = charParser{id: 98, chars: []rune{101}}
	p99.items = []parser{&p94, &p95, &p96, &p97, &p98}
	p100.options = []parser{&p93, &p99}
	var p478 = sequenceParser{id: 478, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 116, 784, 198, 375, 333, 334, 335, 336, 337, 338, 504, 567, 560, 794}}
	var p477 = sequenceParser{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p470 = charParser{id: 470, chars: []rune{114}}
	var p471 = charParser{id: 471, chars: []rune{101}}
	var p472 = charParser{id: 472, chars: []rune{99}}
	var p473 = charParser{id: 473, chars: []rune{101}}
	var p474 = charParser{id: 474, chars: []rune{105}}
	var p475 = charParser{id: 475, chars: []rune{118}}
	var p476 = charParser{id: 476, chars: []rune{101}}
	p477.items = []parser{&p470, &p471, &p472, &p473, &p474, &p475, &p476}
	p478.items = []parser{&p477, &p826, &p267}
	var p105 = sequenceParser{id: 105, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{267, 116, 141, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 733, 794}}
	var p102 = sequenceParser{id: 102, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p101 = charParser{id: 101, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p102.items = []parser{&p101}
	var p104 = sequenceParser{id: 104, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p103 = charParser{id: 103, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p104.items = []parser{&p103}
	p105.items = []parser{&p102, &p104}
	var p126 = sequenceParser{id: 126, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{116, 267, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
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
	p114.items = []parser{&p826, &p113}
	p115.items = []parser{&p113, &p114}
	var p120 = sequenceParser{id: 120, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p116 = choiceParser{id: 116, commit: 66, name: "list-item"}
	var p110 = sequenceParser{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{116, 149, 150}}
	var p109 = sequenceParser{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	var p108 = charParser{id: 108, chars: []rune{46}}
	p109.items = []parser{&p106, &p107, &p108}
	p110.items = []parser{&p267, &p826, &p109}
	p116.options = []parser{&p375, &p110}
	var p119 = sequenceParser{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p117.items = []parser{&p115, &p826, &p116}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p118.items = []parser{&p826, &p117}
	p119.items = []parser{&p826, &p117, &p118}
	p120.items = []parser{&p116, &p119}
	var p124 = sequenceParser{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p123 = charParser{id: 123, chars: []rune{93}}
	p124.items = []parser{&p123}
	p125.items = []parser{&p122, &p826, &p115, &p826, &p120, &p826, &p115, &p826, &p124}
	p126.items = []parser{&p125}
	var p131 = sequenceParser{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p128 = sequenceParser{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p127 = charParser{id: 127, chars: []rune{126}}
	p128.items = []parser{&p127}
	var p130 = sequenceParser{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p129.items = []parser{&p826, &p14}
	p130.items = []parser{&p826, &p14, &p129}
	p131.items = []parser{&p128, &p130, &p826, &p125}
	var p160 = sequenceParser{id: 160, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{267, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
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
	p136.items = []parser{&p826, &p14}
	p137.items = []parser{&p826, &p14, &p136}
	var p139 = sequenceParser{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p138.items = []parser{&p826, &p14}
	p139.items = []parser{&p826, &p14, &p138}
	var p135 = sequenceParser{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p134 = charParser{id: 134, chars: []rune{93}}
	p135.items = []parser{&p134}
	p140.items = []parser{&p133, &p137, &p826, &p375, &p139, &p826, &p135}
	p141.options = []parser{&p105, &p88, &p140}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p144.items = []parser{&p826, &p14}
	p145.items = []parser{&p826, &p14, &p144}
	var p143 = sequenceParser{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p142 = charParser{id: 142, chars: []rune{58}}
	p143.items = []parser{&p142}
	var p147 = sequenceParser{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p146.items = []parser{&p826, &p14}
	p147.items = []parser{&p826, &p14, &p146}
	p148.items = []parser{&p141, &p145, &p826, &p143, &p147, &p826, &p375}
	p149.options = []parser{&p148, &p110}
	var p153 = sequenceParser{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p150 = choiceParser{id: 150, commit: 2}
	p150.options = []parser{&p148, &p110}
	p151.items = []parser{&p115, &p826, &p150}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p152.items = []parser{&p826, &p151}
	p153.items = []parser{&p826, &p151, &p152}
	p154.items = []parser{&p149, &p153}
	var p158 = sequenceParser{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p157 = charParser{id: 157, chars: []rune{125}}
	p158.items = []parser{&p157}
	p159.items = []parser{&p156, &p826, &p115, &p826, &p154, &p826, &p115, &p826, &p158}
	p160.items = []parser{&p159}
	var p165 = sequenceParser{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 784, 198, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p162 = sequenceParser{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p161 = charParser{id: 161, chars: []rune{126}}
	p162.items = []parser{&p161}
	var p164 = sequenceParser{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p163.items = []parser{&p826, &p14}
	p164.items = []parser{&p826, &p14, &p163}
	p165.items = []parser{&p162, &p164, &p826, &p159}
	var p207 = sequenceParser{id: 207, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 198, 267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p204 = sequenceParser{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p202 = charParser{id: 202, chars: []rune{102}}
	var p203 = charParser{id: 203, chars: []rune{110}}
	p204.items = []parser{&p202, &p203}
	var p206 = sequenceParser{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p205 = sequenceParser{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p205.items = []parser{&p826, &p14}
	p206.items = []parser{&p826, &p14, &p205}
	var p201 = sequenceParser{id: 201, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p194 = sequenceParser{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p193 = charParser{id: 193, chars: []rune{40}}
	p194.items = []parser{&p193}
	var p169 = sequenceParser{id: 169, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p168 = sequenceParser{id: 168, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p166.items = []parser{&p115, &p826, &p105}
	var p167 = sequenceParser{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p167.items = []parser{&p826, &p166}
	p168.items = []parser{&p826, &p166, &p167}
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
	p174.items = []parser{&p826, &p14}
	p175.items = []parser{&p826, &p14, &p174}
	p176.items = []parser{&p173, &p175, &p826, &p105}
	p195.items = []parser{&p115, &p826, &p176}
	var p197 = sequenceParser{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p196 = charParser{id: 196, chars: []rune{41}}
	p197.items = []parser{&p196}
	var p200 = sequenceParser{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p199 = sequenceParser{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p199.items = []parser{&p826, &p14}
	p200.items = []parser{&p826, &p14, &p199}
	var p198 = choiceParser{id: 198, commit: 2}
	var p784 = choiceParser{id: 784, commit: 66, name: "simple-statement", generalizations: []int{198, 794}}
	var p503 = sequenceParser{id: 503, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 198, 504, 794}}
	var p502 = sequenceParser{id: 502, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p498 = charParser{id: 498, chars: []rune{115}}
	var p499 = charParser{id: 499, chars: []rune{101}}
	var p500 = charParser{id: 500, chars: []rune{110}}
	var p501 = charParser{id: 501, chars: []rune{100}}
	p502.items = []parser{&p498, &p499, &p500, &p501}
	p503.items = []parser{&p502, &p826, &p267, &p826, &p267}
	var p547 = sequenceParser{id: 547, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 198, 794}}
	var p544 = sequenceParser{id: 544, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p542 = charParser{id: 542, chars: []rune{103}}
	var p543 = charParser{id: 543, chars: []rune{111}}
	p544.items = []parser{&p542, &p543}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p545.items = []parser{&p826, &p14}
	p546.items = []parser{&p826, &p14, &p545}
	var p257 = sequenceParser{id: 257, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p254 = sequenceParser{id: 254, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p253 = charParser{id: 253, chars: []rune{40}}
	p254.items = []parser{&p253}
	var p256 = sequenceParser{id: 256, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p255 = charParser{id: 255, chars: []rune{41}}
	p256.items = []parser{&p255}
	p257.items = []parser{&p267, &p826, &p254, &p826, &p115, &p826, &p120, &p826, &p115, &p826, &p256}
	p547.items = []parser{&p544, &p546, &p826, &p257}
	var p556 = sequenceParser{id: 556, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 198, 794}}
	var p553 = sequenceParser{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p548 = charParser{id: 548, chars: []rune{100}}
	var p549 = charParser{id: 549, chars: []rune{101}}
	var p550 = charParser{id: 550, chars: []rune{102}}
	var p551 = charParser{id: 551, chars: []rune{101}}
	var p552 = charParser{id: 552, chars: []rune{114}}
	p553.items = []parser{&p548, &p549, &p550, &p551, &p552}
	var p555 = sequenceParser{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p554 = sequenceParser{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p554.items = []parser{&p826, &p14}
	p555.items = []parser{&p826, &p14, &p554}
	p556.items = []parser{&p553, &p555, &p826, &p257}
	var p621 = choiceParser{id: 621, commit: 64, name: "assignment", generalizations: []int{784, 198, 794}}
	var p601 = sequenceParser{id: 601, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{621, 784, 198, 794}}
	var p598 = sequenceParser{id: 598, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p595 = charParser{id: 595, chars: []rune{115}}
	var p596 = charParser{id: 596, chars: []rune{101}}
	var p597 = charParser{id: 597, chars: []rune{116}}
	p598.items = []parser{&p595, &p596, &p597}
	var p600 = sequenceParser{id: 600, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p599 = sequenceParser{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p599.items = []parser{&p826, &p14}
	p600.items = []parser{&p826, &p14, &p599}
	var p590 = sequenceParser{id: 590, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p585 = sequenceParser{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p585.items = []parser{&p826, &p14}
	p586.items = []parser{&p14, &p585}
	var p584 = sequenceParser{id: 584, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p583 = charParser{id: 583, chars: []rune{61}}
	p584.items = []parser{&p583}
	p587.items = []parser{&p586, &p826, &p584}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p588.items = []parser{&p826, &p14}
	p589.items = []parser{&p826, &p14, &p588}
	p590.items = []parser{&p267, &p826, &p587, &p589, &p826, &p375}
	p601.items = []parser{&p598, &p600, &p826, &p590}
	var p608 = sequenceParser{id: 608, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{621, 784, 198, 794}}
	var p605 = sequenceParser{id: 605, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p604 = sequenceParser{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p604.items = []parser{&p826, &p14}
	p605.items = []parser{&p826, &p14, &p604}
	var p603 = sequenceParser{id: 603, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p602 = charParser{id: 602, chars: []rune{61}}
	p603.items = []parser{&p602}
	var p607 = sequenceParser{id: 607, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p606 = sequenceParser{id: 606, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p606.items = []parser{&p826, &p14}
	p607.items = []parser{&p826, &p14, &p606}
	p608.items = []parser{&p267, &p605, &p826, &p603, &p607, &p826, &p375}
	var p620 = sequenceParser{id: 620, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{621, 784, 198, 794}}
	var p612 = sequenceParser{id: 612, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p609 = charParser{id: 609, chars: []rune{115}}
	var p610 = charParser{id: 610, chars: []rune{101}}
	var p611 = charParser{id: 611, chars: []rune{116}}
	p612.items = []parser{&p609, &p610, &p611}
	var p619 = sequenceParser{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p618 = sequenceParser{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p618.items = []parser{&p826, &p14}
	p619.items = []parser{&p826, &p14, &p618}
	var p614 = sequenceParser{id: 614, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p613 = charParser{id: 613, chars: []rune{40}}
	p614.items = []parser{&p613}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p594 = sequenceParser{id: 594, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p591 = sequenceParser{id: 591, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p591.items = []parser{&p115, &p826, &p590}
	var p592 = sequenceParser{id: 592, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p592.items = []parser{&p826, &p591}
	p593.items = []parser{&p826, &p591, &p592}
	p594.items = []parser{&p590, &p593}
	p615.items = []parser{&p115, &p826, &p594}
	var p617 = sequenceParser{id: 617, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p616 = charParser{id: 616, chars: []rune{41}}
	p617.items = []parser{&p616}
	p620.items = []parser{&p612, &p619, &p826, &p614, &p826, &p615, &p826, &p115, &p826, &p617}
	p621.options = []parser{&p601, &p608, &p620}
	var p793 = sequenceParser{id: 793, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{784, 198, 794}}
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
	p784.options = []parser{&p503, &p547, &p556, &p621, &p793, &p375}
	var p192 = sequenceParser{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{198}}
	var p189 = sequenceParser{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p188 = charParser{id: 188, chars: []rune{123}}
	p189.items = []parser{&p188}
	var p191 = sequenceParser{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p190 = charParser{id: 190, chars: []rune{125}}
	p191.items = []parser{&p190}
	p192.items = []parser{&p189, &p826, &p808, &p826, &p812, &p826, &p808, &p826, &p191}
	p198.options = []parser{&p784, &p192}
	p201.items = []parser{&p194, &p826, &p115, &p826, &p169, &p826, &p195, &p826, &p115, &p826, &p197, &p200, &p826, &p198}
	p207.items = []parser{&p204, &p206, &p826, &p201}
	var p217 = sequenceParser{id: 217, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p210 = sequenceParser{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p208 = charParser{id: 208, chars: []rune{102}}
	var p209 = charParser{id: 209, chars: []rune{110}}
	p210.items = []parser{&p208, &p209}
	var p214 = sequenceParser{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p213 = sequenceParser{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p213.items = []parser{&p826, &p14}
	p214.items = []parser{&p826, &p14, &p213}
	var p212 = sequenceParser{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p211 = charParser{id: 211, chars: []rune{126}}
	p212.items = []parser{&p211}
	var p216 = sequenceParser{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p215.items = []parser{&p826, &p14}
	p216.items = []parser{&p826, &p14, &p215}
	p217.items = []parser{&p210, &p214, &p826, &p212, &p216, &p826, &p201}
	var p245 = choiceParser{id: 245, commit: 64, name: "expression-indexer", generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p235 = sequenceParser{id: 235, commit: 66, name: "simple-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{245, 267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p228 = sequenceParser{id: 228, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p227 = charParser{id: 227, chars: []rune{91}}
	p228.items = []parser{&p227}
	var p232 = sequenceParser{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p231.items = []parser{&p826, &p14}
	p232.items = []parser{&p826, &p14, &p231}
	var p234 = sequenceParser{id: 234, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p233 = sequenceParser{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p233.items = []parser{&p826, &p14}
	p234.items = []parser{&p826, &p14, &p233}
	var p230 = sequenceParser{id: 230, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p229 = charParser{id: 229, chars: []rune{93}}
	p230.items = []parser{&p229}
	p235.items = []parser{&p267, &p826, &p228, &p232, &p826, &p375, &p234, &p826, &p230}
	var p244 = sequenceParser{id: 244, commit: 66, name: "range-indexer", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{245, 267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p237 = sequenceParser{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p236 = charParser{id: 236, chars: []rune{91}}
	p237.items = []parser{&p236}
	var p241 = sequenceParser{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p240 = sequenceParser{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p240.items = []parser{&p826, &p14}
	p241.items = []parser{&p826, &p14, &p240}
	var p226 = sequenceParser{id: 226, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{560, 566, 567}}
	var p218 = sequenceParser{id: 218, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p218.items = []parser{&p375}
	var p223 = sequenceParser{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p222 = sequenceParser{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p222.items = []parser{&p826, &p14}
	p223.items = []parser{&p826, &p14, &p222}
	var p221 = sequenceParser{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p220 = charParser{id: 220, chars: []rune{58}}
	p221.items = []parser{&p220}
	var p225 = sequenceParser{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p224.items = []parser{&p826, &p14}
	p225.items = []parser{&p826, &p14, &p224}
	var p219 = sequenceParser{id: 219, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p219.items = []parser{&p375}
	p226.items = []parser{&p218, &p223, &p826, &p221, &p225, &p826, &p219}
	var p243 = sequenceParser{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p242 = sequenceParser{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p242.items = []parser{&p826, &p14}
	p243.items = []parser{&p826, &p14, &p242}
	var p239 = sequenceParser{id: 239, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p238 = charParser{id: 238, chars: []rune{93}}
	p239.items = []parser{&p238}
	p244.items = []parser{&p267, &p826, &p237, &p241, &p826, &p226, &p243, &p826, &p239}
	p245.options = []parser{&p235, &p244}
	var p252 = sequenceParser{id: 252, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p248.items = []parser{&p826, &p14}
	p249.items = []parser{&p826, &p14, &p248}
	var p247 = sequenceParser{id: 247, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p246 = charParser{id: 246, chars: []rune{46}}
	p247.items = []parser{&p246}
	var p251 = sequenceParser{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p250.items = []parser{&p826, &p14}
	p251.items = []parser{&p826, &p14, &p250}
	p252.items = []parser{&p267, &p249, &p826, &p247, &p251, &p826, &p105}
	var p266 = sequenceParser{id: 266, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{267, 375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
	var p259 = sequenceParser{id: 259, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p258 = charParser{id: 258, chars: []rune{40}}
	p259.items = []parser{&p258}
	var p263 = sequenceParser{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p262 = sequenceParser{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p262.items = []parser{&p826, &p14}
	p263.items = []parser{&p826, &p14, &p262}
	var p265 = sequenceParser{id: 265, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p264 = sequenceParser{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p264.items = []parser{&p826, &p14}
	p265.items = []parser{&p826, &p14, &p264}
	var p261 = sequenceParser{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p260 = charParser{id: 260, chars: []rune{41}}
	p261.items = []parser{&p260}
	p266.items = []parser{&p259, &p263, &p826, &p375, &p265, &p826, &p261}
	p267.options = []parser{&p62, &p75, &p88, &p100, &p478, &p105, &p126, &p131, &p160, &p165, &p207, &p217, &p245, &p252, &p257, &p266}
	var p327 = sequenceParser{id: 327, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{375, 333, 334, 335, 336, 337, 338, 567, 560, 794}}
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
	p327.items = []parser{&p326, &p826, &p267}
	var p361 = choiceParser{id: 361, commit: 66, name: "binary-expression", generalizations: []int{375, 567, 560, 794}}
	var p341 = sequenceParser{id: 341, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 334, 335, 336, 337, 338, 375, 567, 560, 794}}
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
	p339.items = []parser{&p328, &p826, &p333}
	var p340 = sequenceParser{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p340.items = []parser{&p826, &p339}
	p341.items = []parser{&p333, &p826, &p339, &p340}
	var p344 = sequenceParser{id: 344, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 335, 336, 337, 338, 375, 567, 560, 794}}
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
	p342.items = []parser{&p329, &p826, &p334}
	var p343 = sequenceParser{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p343.items = []parser{&p826, &p342}
	p344.items = []parser{&p334, &p826, &p342, &p343}
	var p347 = sequenceParser{id: 347, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 336, 337, 338, 375, 567, 560, 794}}
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
	p345.items = []parser{&p330, &p826, &p335}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p346.items = []parser{&p826, &p345}
	p347.items = []parser{&p335, &p826, &p345, &p346}
	var p350 = sequenceParser{id: 350, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 337, 338, 375, 567, 560, 794}}
	var p336 = choiceParser{id: 336, commit: 66, name: "operand3", generalizations: []int{337, 338}}
	p336.options = []parser{&p335, &p347}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p331 = sequenceParser{id: 331, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p319 = sequenceParser{id: 319, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p317 = charParser{id: 317, chars: []rune{38}}
	var p318 = charParser{id: 318, chars: []rune{38}}
	p319.items = []parser{&p317, &p318}
	p331.items = []parser{&p319}
	p348.items = []parser{&p331, &p826, &p336}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p349.items = []parser{&p826, &p348}
	p350.items = []parser{&p336, &p826, &p348, &p349}
	var p353 = sequenceParser{id: 353, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 338, 375, 567, 560, 794}}
	var p337 = choiceParser{id: 337, commit: 66, name: "operand4", generalizations: []int{338}}
	p337.options = []parser{&p336, &p350}
	var p351 = sequenceParser{id: 351, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p332 = sequenceParser{id: 332, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p322 = sequenceParser{id: 322, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p320 = charParser{id: 320, chars: []rune{124}}
	var p321 = charParser{id: 321, chars: []rune{124}}
	p322.items = []parser{&p320, &p321}
	p332.items = []parser{&p322}
	p351.items = []parser{&p332, &p826, &p337}
	var p352 = sequenceParser{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p352.items = []parser{&p826, &p351}
	p353.items = []parser{&p337, &p826, &p351, &p352}
	var p360 = sequenceParser{id: 360, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{361, 375, 567, 560, 794}}
	var p338 = choiceParser{id: 338, commit: 66, name: "operand5"}
	p338.options = []parser{&p337, &p353}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p354.items = []parser{&p826, &p14}
	p355.items = []parser{&p14, &p354}
	var p325 = sequenceParser{id: 325, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p323 = charParser{id: 323, chars: []rune{45}}
	var p324 = charParser{id: 324, chars: []rune{62}}
	p325.items = []parser{&p323, &p324}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p356.items = []parser{&p826, &p14}
	p357.items = []parser{&p826, &p14, &p356}
	p358.items = []parser{&p355, &p826, &p325, &p357, &p826, &p338}
	var p359 = sequenceParser{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p359.items = []parser{&p826, &p358}
	p360.items = []parser{&p338, &p826, &p358, &p359}
	p361.options = []parser{&p341, &p344, &p347, &p350, &p353, &p360}
	var p374 = sequenceParser{id: 374, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{375, 567, 560, 794}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p366 = sequenceParser{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p366.items = []parser{&p826, &p14}
	p367.items = []parser{&p826, &p14, &p366}
	var p363 = sequenceParser{id: 363, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p362 = charParser{id: 362, chars: []rune{63}}
	p363.items = []parser{&p362}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p368.items = []parser{&p826, &p14}
	p369.items = []parser{&p826, &p14, &p368}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p370.items = []parser{&p826, &p14}
	p371.items = []parser{&p826, &p14, &p370}
	var p365 = sequenceParser{id: 365, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p364 = charParser{id: 364, chars: []rune{58}}
	p365.items = []parser{&p364}
	var p373 = sequenceParser{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p372.items = []parser{&p826, &p14}
	p373.items = []parser{&p826, &p14, &p372}
	p374.items = []parser{&p375, &p367, &p826, &p363, &p369, &p826, &p375, &p371, &p826, &p365, &p373, &p826, &p375}
	p375.options = []parser{&p267, &p327, &p361, &p374}
	p186.items = []parser{&p185, &p826, &p375}
	p187.items = []parser{&p183, &p826, &p186}
	var p412 = sequenceParser{id: 412, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{794, 458, 532}}
	var p378 = sequenceParser{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p376 = charParser{id: 376, chars: []rune{105}}
	var p377 = charParser{id: 377, chars: []rune{102}}
	p378.items = []parser{&p376, &p377}
	var p407 = sequenceParser{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p406 = sequenceParser{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p406.items = []parser{&p826, &p14}
	p407.items = []parser{&p826, &p14, &p406}
	var p409 = sequenceParser{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p408 = sequenceParser{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p408.items = []parser{&p826, &p14}
	p409.items = []parser{&p826, &p14, &p408}
	var p411 = sequenceParser{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p395 = sequenceParser{id: 395, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p387.items = []parser{&p826, &p14}
	p388.items = []parser{&p14, &p387}
	var p383 = sequenceParser{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p379 = charParser{id: 379, chars: []rune{101}}
	var p380 = charParser{id: 380, chars: []rune{108}}
	var p381 = charParser{id: 381, chars: []rune{115}}
	var p382 = charParser{id: 382, chars: []rune{101}}
	p383.items = []parser{&p379, &p380, &p381, &p382}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p389.items = []parser{&p826, &p14}
	p390.items = []parser{&p826, &p14, &p389}
	var p386 = sequenceParser{id: 386, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p384 = charParser{id: 384, chars: []rune{105}}
	var p385 = charParser{id: 385, chars: []rune{102}}
	p386.items = []parser{&p384, &p385}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p391.items = []parser{&p826, &p14}
	p392.items = []parser{&p826, &p14, &p391}
	var p394 = sequenceParser{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p393.items = []parser{&p826, &p14}
	p394.items = []parser{&p826, &p14, &p393}
	p395.items = []parser{&p388, &p826, &p383, &p390, &p826, &p386, &p392, &p826, &p375, &p394, &p826, &p192}
	var p410 = sequenceParser{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p410.items = []parser{&p826, &p395}
	p411.items = []parser{&p826, &p395, &p410}
	var p405 = sequenceParser{id: 405, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p402 = sequenceParser{id: 402, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p401 = sequenceParser{id: 401, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p401.items = []parser{&p826, &p14}
	p402.items = []parser{&p14, &p401}
	var p400 = sequenceParser{id: 400, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p396 = charParser{id: 396, chars: []rune{101}}
	var p397 = charParser{id: 397, chars: []rune{108}}
	var p398 = charParser{id: 398, chars: []rune{115}}
	var p399 = charParser{id: 399, chars: []rune{101}}
	p400.items = []parser{&p396, &p397, &p398, &p399}
	var p404 = sequenceParser{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p403 = sequenceParser{id: 403, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p403.items = []parser{&p826, &p14}
	p404.items = []parser{&p826, &p14, &p403}
	p405.items = []parser{&p402, &p826, &p400, &p404, &p826, &p192}
	p412.items = []parser{&p378, &p407, &p826, &p375, &p409, &p826, &p192, &p411, &p826, &p405}
	var p469 = sequenceParser{id: 469, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{458, 794, 532}}
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
	p465.items = []parser{&p826, &p14}
	p466.items = []parser{&p826, &p14, &p465}
	var p468 = sequenceParser{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p467 = sequenceParser{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p467.items = []parser{&p826, &p14}
	p468.items = []parser{&p826, &p14, &p467}
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
	p438.items = []parser{&p826, &p14}
	p439.items = []parser{&p826, &p14, &p438}
	var p441 = sequenceParser{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p440 = sequenceParser{id: 440, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p440.items = []parser{&p826, &p14}
	p441.items = []parser{&p826, &p14, &p440}
	var p437 = sequenceParser{id: 437, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p436 = charParser{id: 436, chars: []rune{58}}
	p437.items = []parser{&p436}
	p442.items = []parser{&p435, &p439, &p826, &p375, &p441, &p826, &p437}
	var p446 = sequenceParser{id: 446, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p444 = sequenceParser{id: 444, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p443 = charParser{id: 443, chars: []rune{59}}
	p444.items = []parser{&p443}
	var p445 = sequenceParser{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p445.items = []parser{&p826, &p444}
	p446.items = []parser{&p826, &p444, &p445}
	p447.items = []parser{&p442, &p446, &p826, &p794}
	var p430 = sequenceParser{id: 430, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{457, 458, 531, 532}}
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
	p423.items = []parser{&p826, &p14}
	p424.items = []parser{&p826, &p14, &p423}
	var p422 = sequenceParser{id: 422, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p421 = charParser{id: 421, chars: []rune{58}}
	p422.items = []parser{&p421}
	p425.items = []parser{&p420, &p424, &p826, &p422}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p427 = sequenceParser{id: 427, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p426 = charParser{id: 426, chars: []rune{59}}
	p427.items = []parser{&p426}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p428.items = []parser{&p826, &p427}
	p429.items = []parser{&p826, &p427, &p428}
	p430.items = []parser{&p425, &p429, &p826, &p794}
	p457.options = []parser{&p447, &p430}
	var p461 = sequenceParser{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p458 = choiceParser{id: 458, commit: 2}
	p458.options = []parser{&p447, &p430, &p794}
	p459.items = []parser{&p808, &p826, &p458}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p460.items = []parser{&p826, &p459}
	p461.items = []parser{&p826, &p459, &p460}
	p462.items = []parser{&p457, &p461}
	var p464 = sequenceParser{id: 464, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p463 = charParser{id: 463, chars: []rune{125}}
	p464.items = []parser{&p463}
	p469.items = []parser{&p454, &p466, &p826, &p375, &p468, &p826, &p456, &p826, &p808, &p826, &p462, &p826, &p808, &p826, &p464}
	var p541 = sequenceParser{id: 541, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{532, 794}}
	var p528 = sequenceParser{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p522 = charParser{id: 522, chars: []rune{115}}
	var p523 = charParser{id: 523, chars: []rune{101}}
	var p524 = charParser{id: 524, chars: []rune{108}}
	var p525 = charParser{id: 525, chars: []rune{101}}
	var p526 = charParser{id: 526, chars: []rune{99}}
	var p527 = charParser{id: 527, chars: []rune{116}}
	p528.items = []parser{&p522, &p523, &p524, &p525, &p526, &p527}
	var p540 = sequenceParser{id: 540, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p539 = sequenceParser{id: 539, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p539.items = []parser{&p826, &p14}
	p540.items = []parser{&p826, &p14, &p539}
	var p530 = sequenceParser{id: 530, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p529 = charParser{id: 529, chars: []rune{123}}
	p530.items = []parser{&p529}
	var p536 = sequenceParser{id: 536, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p531 = choiceParser{id: 531, commit: 2}
	var p521 = sequenceParser{id: 521, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{531, 532}}
	var p516 = sequenceParser{id: 516, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p509 = sequenceParser{id: 509, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p505 = charParser{id: 505, chars: []rune{99}}
	var p506 = charParser{id: 506, chars: []rune{97}}
	var p507 = charParser{id: 507, chars: []rune{115}}
	var p508 = charParser{id: 508, chars: []rune{101}}
	p509.items = []parser{&p505, &p506, &p507, &p508}
	var p513 = sequenceParser{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p512 = sequenceParser{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p512.items = []parser{&p826, &p14}
	p513.items = []parser{&p826, &p14, &p512}
	var p504 = choiceParser{id: 504, commit: 66, name: "communication"}
	var p497 = choiceParser{id: 497, commit: 66, name: "receive-statement", generalizations: []int{504}}
	var p487 = sequenceParser{id: 487, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{497, 504}}
	var p482 = sequenceParser{id: 482, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p479 = charParser{id: 479, chars: []rune{108}}
	var p480 = charParser{id: 480, chars: []rune{101}}
	var p481 = charParser{id: 481, chars: []rune{116}}
	p482.items = []parser{&p479, &p480, &p481}
	var p484 = sequenceParser{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p483 = sequenceParser{id: 483, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p483.items = []parser{&p826, &p14}
	p484.items = []parser{&p826, &p14, &p483}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p485.items = []parser{&p826, &p14}
	p486.items = []parser{&p826, &p14, &p485}
	p487.items = []parser{&p482, &p484, &p826, &p105, &p486, &p826, &p478}
	var p496 = sequenceParser{id: 496, commit: 64, name: "receive-assignment", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{497, 504}}
	var p491 = sequenceParser{id: 491, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p488 = charParser{id: 488, chars: []rune{115}}
	var p489 = charParser{id: 489, chars: []rune{101}}
	var p490 = charParser{id: 490, chars: []rune{116}}
	p491.items = []parser{&p488, &p489, &p490}
	var p493 = sequenceParser{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p492.items = []parser{&p826, &p14}
	p493.items = []parser{&p826, &p14, &p492}
	var p495 = sequenceParser{id: 495, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p494 = sequenceParser{id: 494, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p494.items = []parser{&p826, &p14}
	p495.items = []parser{&p826, &p14, &p494}
	p496.items = []parser{&p491, &p493, &p826, &p105, &p495, &p826, &p478}
	p497.options = []parser{&p487, &p496}
	p504.options = []parser{&p478, &p497, &p503}
	var p515 = sequenceParser{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p514 = sequenceParser{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p514.items = []parser{&p826, &p14}
	p515.items = []parser{&p826, &p14, &p514}
	var p511 = sequenceParser{id: 511, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p510 = charParser{id: 510, chars: []rune{58}}
	p511.items = []parser{&p510}
	p516.items = []parser{&p509, &p513, &p826, &p504, &p515, &p826, &p511}
	var p520 = sequenceParser{id: 520, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p518 = sequenceParser{id: 518, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p517 = charParser{id: 517, chars: []rune{59}}
	p518.items = []parser{&p517}
	var p519 = sequenceParser{id: 519, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p519.items = []parser{&p826, &p518}
	p520.items = []parser{&p826, &p518, &p519}
	p521.items = []parser{&p516, &p520, &p826, &p794}
	p531.options = []parser{&p521, &p430}
	var p535 = sequenceParser{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p533 = sequenceParser{id: 533, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p532 = choiceParser{id: 532, commit: 2}
	p532.options = []parser{&p521, &p430, &p794}
	p533.items = []parser{&p808, &p826, &p532}
	var p534 = sequenceParser{id: 534, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p534.items = []parser{&p826, &p533}
	p535.items = []parser{&p826, &p533, &p534}
	p536.items = []parser{&p531, &p535}
	var p538 = sequenceParser{id: 538, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p537 = charParser{id: 537, chars: []rune{125}}
	p538.items = []parser{&p537}
	p541.items = []parser{&p528, &p540, &p826, &p530, &p826, &p808, &p826, &p536, &p826, &p808, &p826, &p538}
	var p582 = sequenceParser{id: 582, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{794}}
	var p571 = sequenceParser{id: 571, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p568 = charParser{id: 568, chars: []rune{102}}
	var p569 = charParser{id: 569, chars: []rune{111}}
	var p570 = charParser{id: 570, chars: []rune{114}}
	p571.items = []parser{&p568, &p569, &p570}
	var p581 = choiceParser{id: 581, commit: 2}
	var p577 = sequenceParser{id: 577, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{581}}
	var p574 = sequenceParser{id: 574, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p572 = sequenceParser{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p572.items = []parser{&p826, &p14}
	p573.items = []parser{&p14, &p572}
	var p567 = choiceParser{id: 567, commit: 66, name: "loop-expression"}
	var p566 = choiceParser{id: 566, commit: 64, name: "range-over-expression", generalizations: []int{567}}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{566, 567}}
	var p562 = sequenceParser{id: 562, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p561 = sequenceParser{id: 561, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p561.items = []parser{&p826, &p14}
	p562.items = []parser{&p826, &p14, &p561}
	var p559 = sequenceParser{id: 559, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p557 = charParser{id: 557, chars: []rune{105}}
	var p558 = charParser{id: 558, chars: []rune{110}}
	p559.items = []parser{&p557, &p558}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p563 = sequenceParser{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p563.items = []parser{&p826, &p14}
	p564.items = []parser{&p826, &p14, &p563}
	var p560 = choiceParser{id: 560, commit: 2}
	p560.options = []parser{&p375, &p226}
	p565.items = []parser{&p105, &p562, &p826, &p559, &p564, &p826, &p560}
	p566.options = []parser{&p565, &p226}
	p567.options = []parser{&p375, &p566}
	p574.items = []parser{&p573, &p826, &p567}
	var p576 = sequenceParser{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p575 = sequenceParser{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p575.items = []parser{&p826, &p14}
	p576.items = []parser{&p826, &p14, &p575}
	p577.items = []parser{&p574, &p576, &p826, &p192}
	var p580 = sequenceParser{id: 580, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{581}}
	var p579 = sequenceParser{id: 579, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p578 = sequenceParser{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p578.items = []parser{&p826, &p14}
	p579.items = []parser{&p14, &p578}
	p580.items = []parser{&p579, &p826, &p192}
	p581.options = []parser{&p577, &p580}
	p582.items = []parser{&p571, &p826, &p581}
	var p730 = choiceParser{id: 730, commit: 66, name: "definition", generalizations: []int{794}}
	var p643 = sequenceParser{id: 643, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 794}}
	var p639 = sequenceParser{id: 639, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p636 = charParser{id: 636, chars: []rune{108}}
	var p637 = charParser{id: 637, chars: []rune{101}}
	var p638 = charParser{id: 638, chars: []rune{116}}
	p639.items = []parser{&p636, &p637, &p638}
	var p642 = sequenceParser{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p641 = sequenceParser{id: 641, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p641.items = []parser{&p826, &p14}
	p642.items = []parser{&p826, &p14, &p641}
	var p640 = choiceParser{id: 640, commit: 2}
	var p630 = sequenceParser{id: 630, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{640, 644, 645}}
	var p629 = sequenceParser{id: 629, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p626 = sequenceParser{id: 626, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p625 = sequenceParser{id: 625, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p624 = sequenceParser{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p624.items = []parser{&p826, &p14}
	p625.items = []parser{&p14, &p624}
	var p623 = sequenceParser{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p622 = charParser{id: 622, chars: []rune{61}}
	p623.items = []parser{&p622}
	p626.items = []parser{&p625, &p826, &p623}
	var p628 = sequenceParser{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p627 = sequenceParser{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p627.items = []parser{&p826, &p14}
	p628.items = []parser{&p826, &p14, &p627}
	p629.items = []parser{&p105, &p826, &p626, &p628, &p826, &p375}
	p630.items = []parser{&p629}
	var p635 = sequenceParser{id: 635, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{640, 644, 645}}
	var p632 = sequenceParser{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p631 = charParser{id: 631, chars: []rune{126}}
	p632.items = []parser{&p631}
	var p634 = sequenceParser{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p633 = sequenceParser{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p633.items = []parser{&p826, &p14}
	p634.items = []parser{&p826, &p14, &p633}
	p635.items = []parser{&p632, &p634, &p826, &p629}
	p640.options = []parser{&p630, &p635}
	p643.items = []parser{&p639, &p642, &p826, &p640}
	var p664 = sequenceParser{id: 664, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 794}}
	var p657 = sequenceParser{id: 657, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p654 = charParser{id: 654, chars: []rune{108}}
	var p655 = charParser{id: 655, chars: []rune{101}}
	var p656 = charParser{id: 656, chars: []rune{116}}
	p657.items = []parser{&p654, &p655, &p656}
	var p663 = sequenceParser{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p662 = sequenceParser{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p662.items = []parser{&p826, &p14}
	p663.items = []parser{&p826, &p14, &p662}
	var p659 = sequenceParser{id: 659, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p658 = charParser{id: 658, chars: []rune{40}}
	p659.items = []parser{&p658}
	var p649 = sequenceParser{id: 649, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p644 = choiceParser{id: 644, commit: 2}
	p644.options = []parser{&p630, &p635}
	var p648 = sequenceParser{id: 648, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p646 = sequenceParser{id: 646, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p645 = choiceParser{id: 645, commit: 2}
	p645.options = []parser{&p630, &p635}
	p646.items = []parser{&p115, &p826, &p645}
	var p647 = sequenceParser{id: 647, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p647.items = []parser{&p826, &p646}
	p648.items = []parser{&p826, &p646, &p647}
	p649.items = []parser{&p644, &p648}
	var p661 = sequenceParser{id: 661, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p660 = charParser{id: 660, chars: []rune{41}}
	p661.items = []parser{&p660}
	p664.items = []parser{&p657, &p663, &p826, &p659, &p826, &p115, &p826, &p649, &p826, &p115, &p826, &p661}
	var p679 = sequenceParser{id: 679, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 794}}
	var p668 = sequenceParser{id: 668, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p665 = charParser{id: 665, chars: []rune{108}}
	var p666 = charParser{id: 666, chars: []rune{101}}
	var p667 = charParser{id: 667, chars: []rune{116}}
	p668.items = []parser{&p665, &p666, &p667}
	var p676 = sequenceParser{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p675 = sequenceParser{id: 675, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p675.items = []parser{&p826, &p14}
	p676.items = []parser{&p826, &p14, &p675}
	var p670 = sequenceParser{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p669 = charParser{id: 669, chars: []rune{126}}
	p670.items = []parser{&p669}
	var p678 = sequenceParser{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p677 = sequenceParser{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p677.items = []parser{&p826, &p14}
	p678.items = []parser{&p826, &p14, &p677}
	var p672 = sequenceParser{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p671 = charParser{id: 671, chars: []rune{40}}
	p672.items = []parser{&p671}
	var p653 = sequenceParser{id: 653, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p652 = sequenceParser{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p650 = sequenceParser{id: 650, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p650.items = []parser{&p115, &p826, &p630}
	var p651 = sequenceParser{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p651.items = []parser{&p826, &p650}
	p652.items = []parser{&p826, &p650, &p651}
	p653.items = []parser{&p630, &p652}
	var p674 = sequenceParser{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p673 = charParser{id: 673, chars: []rune{41}}
	p674.items = []parser{&p673}
	p679.items = []parser{&p668, &p676, &p826, &p670, &p678, &p826, &p672, &p826, &p115, &p826, &p653, &p826, &p115, &p826, &p674}
	var p695 = sequenceParser{id: 695, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 794}}
	var p691 = sequenceParser{id: 691, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p689 = charParser{id: 689, chars: []rune{102}}
	var p690 = charParser{id: 690, chars: []rune{110}}
	p691.items = []parser{&p689, &p690}
	var p694 = sequenceParser{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p693 = sequenceParser{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p693.items = []parser{&p826, &p14}
	p694.items = []parser{&p826, &p14, &p693}
	var p692 = choiceParser{id: 692, commit: 2}
	var p683 = sequenceParser{id: 683, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{692, 700, 701}}
	var p682 = sequenceParser{id: 682, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p681 = sequenceParser{id: 681, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p680 = sequenceParser{id: 680, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p680.items = []parser{&p826, &p14}
	p681.items = []parser{&p826, &p14, &p680}
	p682.items = []parser{&p105, &p681, &p826, &p201}
	p683.items = []parser{&p682}
	var p688 = sequenceParser{id: 688, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{692, 700, 701}}
	var p685 = sequenceParser{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p684 = charParser{id: 684, chars: []rune{126}}
	p685.items = []parser{&p684}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p686.items = []parser{&p826, &p14}
	p687.items = []parser{&p826, &p14, &p686}
	p688.items = []parser{&p685, &p687, &p826, &p682}
	p692.options = []parser{&p683, &p688}
	p695.items = []parser{&p691, &p694, &p826, &p692}
	var p715 = sequenceParser{id: 715, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 794}}
	var p708 = sequenceParser{id: 708, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p706 = charParser{id: 706, chars: []rune{102}}
	var p707 = charParser{id: 707, chars: []rune{110}}
	p708.items = []parser{&p706, &p707}
	var p714 = sequenceParser{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p713 = sequenceParser{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p713.items = []parser{&p826, &p14}
	p714.items = []parser{&p826, &p14, &p713}
	var p710 = sequenceParser{id: 710, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p709 = charParser{id: 709, chars: []rune{40}}
	p710.items = []parser{&p709}
	var p705 = sequenceParser{id: 705, commit: 66, name: "mixed-function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p700 = choiceParser{id: 700, commit: 2}
	p700.options = []parser{&p683, &p688}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p702 = sequenceParser{id: 702, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p701 = choiceParser{id: 701, commit: 2}
	p701.options = []parser{&p683, &p688}
	p702.items = []parser{&p115, &p826, &p701}
	var p703 = sequenceParser{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p703.items = []parser{&p826, &p702}
	p704.items = []parser{&p826, &p702, &p703}
	p705.items = []parser{&p700, &p704}
	var p712 = sequenceParser{id: 712, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p711 = charParser{id: 711, chars: []rune{41}}
	p712.items = []parser{&p711}
	p715.items = []parser{&p708, &p714, &p826, &p710, &p826, &p115, &p826, &p705, &p826, &p115, &p826, &p712}
	var p729 = sequenceParser{id: 729, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{730, 794}}
	var p718 = sequenceParser{id: 718, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p716 = charParser{id: 716, chars: []rune{102}}
	var p717 = charParser{id: 717, chars: []rune{110}}
	p718.items = []parser{&p716, &p717}
	var p726 = sequenceParser{id: 726, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p725 = sequenceParser{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p725.items = []parser{&p826, &p14}
	p726.items = []parser{&p826, &p14, &p725}
	var p720 = sequenceParser{id: 720, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p719 = charParser{id: 719, chars: []rune{126}}
	p720.items = []parser{&p719}
	var p728 = sequenceParser{id: 728, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p727 = sequenceParser{id: 727, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p727.items = []parser{&p826, &p14}
	p728.items = []parser{&p826, &p14, &p727}
	var p722 = sequenceParser{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p721 = charParser{id: 721, chars: []rune{40}}
	p722.items = []parser{&p721}
	var p699 = sequenceParser{id: 699, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p698 = sequenceParser{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p696 = sequenceParser{id: 696, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p696.items = []parser{&p115, &p826, &p683}
	var p697 = sequenceParser{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p697.items = []parser{&p826, &p696}
	p698.items = []parser{&p826, &p696, &p697}
	p699.items = []parser{&p683, &p698}
	var p724 = sequenceParser{id: 724, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p723 = charParser{id: 723, chars: []rune{41}}
	p724.items = []parser{&p723}
	p729.items = []parser{&p718, &p726, &p826, &p720, &p728, &p826, &p722, &p826, &p115, &p826, &p699, &p826, &p115, &p826, &p724}
	p730.options = []parser{&p643, &p664, &p679, &p695, &p715, &p729}
	var p773 = choiceParser{id: 773, commit: 64, name: "require", generalizations: []int{794}}
	var p757 = sequenceParser{id: 757, commit: 66, name: "require-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{773, 794}}
	var p754 = sequenceParser{id: 754, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p747 = charParser{id: 747, chars: []rune{114}}
	var p748 = charParser{id: 748, chars: []rune{101}}
	var p749 = charParser{id: 749, chars: []rune{113}}
	var p750 = charParser{id: 750, chars: []rune{117}}
	var p751 = charParser{id: 751, chars: []rune{105}}
	var p752 = charParser{id: 752, chars: []rune{114}}
	var p753 = charParser{id: 753, chars: []rune{101}}
	p754.items = []parser{&p747, &p748, &p749, &p750, &p751, &p752, &p753}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p755.items = []parser{&p826, &p14}
	p756.items = []parser{&p826, &p14, &p755}
	var p742 = choiceParser{id: 742, commit: 64, name: "require-fact"}
	var p741 = sequenceParser{id: 741, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{742}}
	var p733 = choiceParser{id: 733, commit: 2}
	var p732 = sequenceParser{id: 732, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{733}}
	var p731 = charParser{id: 731, chars: []rune{46}}
	p732.items = []parser{&p731}
	p733.options = []parser{&p105, &p732}
	var p738 = sequenceParser{id: 738, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p736.items = []parser{&p826, &p14}
	p737.items = []parser{&p14, &p736}
	var p735 = sequenceParser{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p734 = charParser{id: 734, chars: []rune{61}}
	p735.items = []parser{&p734}
	p738.items = []parser{&p737, &p826, &p735}
	var p740 = sequenceParser{id: 740, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p739 = sequenceParser{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p739.items = []parser{&p826, &p14}
	p740.items = []parser{&p826, &p14, &p739}
	p741.items = []parser{&p733, &p826, &p738, &p740, &p826, &p88}
	p742.options = []parser{&p88, &p741}
	p757.items = []parser{&p754, &p756, &p826, &p742}
	var p772 = sequenceParser{id: 772, commit: 66, name: "require-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{773, 794}}
	var p765 = sequenceParser{id: 765, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p758 = charParser{id: 758, chars: []rune{114}}
	var p759 = charParser{id: 759, chars: []rune{101}}
	var p760 = charParser{id: 760, chars: []rune{113}}
	var p761 = charParser{id: 761, chars: []rune{117}}
	var p762 = charParser{id: 762, chars: []rune{105}}
	var p763 = charParser{id: 763, chars: []rune{114}}
	var p764 = charParser{id: 764, chars: []rune{101}}
	p765.items = []parser{&p758, &p759, &p760, &p761, &p762, &p763, &p764}
	var p771 = sequenceParser{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p770 = sequenceParser{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p770.items = []parser{&p826, &p14}
	p771.items = []parser{&p826, &p14, &p770}
	var p767 = sequenceParser{id: 767, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p766 = charParser{id: 766, chars: []rune{40}}
	p767.items = []parser{&p766}
	var p746 = sequenceParser{id: 746, commit: 66, name: "require-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p745 = sequenceParser{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p743 = sequenceParser{id: 743, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p743.items = []parser{&p115, &p826, &p742}
	var p744 = sequenceParser{id: 744, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p744.items = []parser{&p826, &p743}
	p745.items = []parser{&p826, &p743, &p744}
	p746.items = []parser{&p742, &p745}
	var p769 = sequenceParser{id: 769, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p768 = charParser{id: 768, chars: []rune{41}}
	p769.items = []parser{&p768}
	p772.items = []parser{&p765, &p771, &p826, &p767, &p826, &p115, &p826, &p746, &p826, &p115, &p826, &p769}
	p773.options = []parser{&p757, &p772}
	var p783 = sequenceParser{id: 783, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{794}}
	var p780 = sequenceParser{id: 780, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p774 = charParser{id: 774, chars: []rune{101}}
	var p775 = charParser{id: 775, chars: []rune{120}}
	var p776 = charParser{id: 776, chars: []rune{112}}
	var p777 = charParser{id: 777, chars: []rune{111}}
	var p778 = charParser{id: 778, chars: []rune{114}}
	var p779 = charParser{id: 779, chars: []rune{116}}
	p780.items = []parser{&p774, &p775, &p776, &p777, &p778, &p779}
	var p782 = sequenceParser{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p781 = sequenceParser{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p781.items = []parser{&p826, &p14}
	p782.items = []parser{&p826, &p14, &p781}
	p783.items = []parser{&p780, &p782, &p826, &p730}
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
	p794.options = []parser{&p187, &p412, &p469, &p541, &p582, &p730, &p773, &p783, &p803, &p784}
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
	b41.items = []builder{&b40, &b826, &b38}
	var b42 = sequenceBuilder{id: 42, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b42.items = []builder{&b826, &b41}
	b43.items = []builder{&b826, &b41, &b42}
	b44.items = []builder{&b38, &b43}
	b825.items = []builder{&b44}
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
	var b14 = sequenceBuilder{id: 14, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b13 = charBuilder{}
	b14.items = []builder{&b13}
	b806.options = []builder{&b805, &b14}
	var b807 = sequenceBuilder{id: 807, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b807.items = []builder{&b826, &b806}
	b808.items = []builder{&b806, &b807}
	var b812 = sequenceBuilder{id: 812, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b794 = choiceBuilder{id: 794, commit: 66}
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
	b184.items = []builder{&b826, &b14}
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
	var b478 = sequenceBuilder{id: 478, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b477 = sequenceBuilder{id: 477, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b470 = charBuilder{}
	var b471 = charBuilder{}
	var b472 = charBuilder{}
	var b473 = charBuilder{}
	var b474 = charBuilder{}
	var b475 = charBuilder{}
	var b476 = charBuilder{}
	b477.items = []builder{&b470, &b471, &b472, &b473, &b474, &b475, &b476}
	b478.items = []builder{&b477, &b826, &b267}
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
	b114.items = []builder{&b826, &b113}
	b115.items = []builder{&b113, &b114}
	var b120 = sequenceBuilder{id: 120, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b116 = choiceBuilder{id: 116, commit: 66}
	var b110 = sequenceBuilder{id: 110, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b109 = sequenceBuilder{id: 109, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	var b108 = charBuilder{}
	b109.items = []builder{&b106, &b107, &b108}
	b110.items = []builder{&b267, &b826, &b109}
	b116.options = []builder{&b375, &b110}
	var b119 = sequenceBuilder{id: 119, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b117.items = []builder{&b115, &b826, &b116}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b118.items = []builder{&b826, &b117}
	b119.items = []builder{&b826, &b117, &b118}
	b120.items = []builder{&b116, &b119}
	var b124 = sequenceBuilder{id: 124, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b123 = charBuilder{}
	b124.items = []builder{&b123}
	b125.items = []builder{&b122, &b826, &b115, &b826, &b120, &b826, &b115, &b826, &b124}
	b126.items = []builder{&b125}
	var b131 = sequenceBuilder{id: 131, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b128 = sequenceBuilder{id: 128, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b127 = charBuilder{}
	b128.items = []builder{&b127}
	var b130 = sequenceBuilder{id: 130, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b129.items = []builder{&b826, &b14}
	b130.items = []builder{&b826, &b14, &b129}
	b131.items = []builder{&b128, &b130, &b826, &b125}
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
	b136.items = []builder{&b826, &b14}
	b137.items = []builder{&b826, &b14, &b136}
	var b139 = sequenceBuilder{id: 139, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b138.items = []builder{&b826, &b14}
	b139.items = []builder{&b826, &b14, &b138}
	var b135 = sequenceBuilder{id: 135, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b134 = charBuilder{}
	b135.items = []builder{&b134}
	b140.items = []builder{&b133, &b137, &b826, &b375, &b139, &b826, &b135}
	b141.options = []builder{&b105, &b88, &b140}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b144.items = []builder{&b826, &b14}
	b145.items = []builder{&b826, &b14, &b144}
	var b143 = sequenceBuilder{id: 143, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b142 = charBuilder{}
	b143.items = []builder{&b142}
	var b147 = sequenceBuilder{id: 147, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b146.items = []builder{&b826, &b14}
	b147.items = []builder{&b826, &b14, &b146}
	b148.items = []builder{&b141, &b145, &b826, &b143, &b147, &b826, &b375}
	b149.options = []builder{&b148, &b110}
	var b153 = sequenceBuilder{id: 153, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b150 = choiceBuilder{id: 150, commit: 2}
	b150.options = []builder{&b148, &b110}
	b151.items = []builder{&b115, &b826, &b150}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b152.items = []builder{&b826, &b151}
	b153.items = []builder{&b826, &b151, &b152}
	b154.items = []builder{&b149, &b153}
	var b158 = sequenceBuilder{id: 158, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b157 = charBuilder{}
	b158.items = []builder{&b157}
	b159.items = []builder{&b156, &b826, &b115, &b826, &b154, &b826, &b115, &b826, &b158}
	b160.items = []builder{&b159}
	var b165 = sequenceBuilder{id: 165, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b162 = sequenceBuilder{id: 162, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b161 = charBuilder{}
	b162.items = []builder{&b161}
	var b164 = sequenceBuilder{id: 164, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b163.items = []builder{&b826, &b14}
	b164.items = []builder{&b826, &b14, &b163}
	b165.items = []builder{&b162, &b164, &b826, &b159}
	var b207 = sequenceBuilder{id: 207, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b204 = sequenceBuilder{id: 204, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b202 = charBuilder{}
	var b203 = charBuilder{}
	b204.items = []builder{&b202, &b203}
	var b206 = sequenceBuilder{id: 206, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b205 = sequenceBuilder{id: 205, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b205.items = []builder{&b826, &b14}
	b206.items = []builder{&b826, &b14, &b205}
	var b201 = sequenceBuilder{id: 201, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b194 = sequenceBuilder{id: 194, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b193 = charBuilder{}
	b194.items = []builder{&b193}
	var b169 = sequenceBuilder{id: 169, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b168 = sequenceBuilder{id: 168, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b166.items = []builder{&b115, &b826, &b105}
	var b167 = sequenceBuilder{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b167.items = []builder{&b826, &b166}
	b168.items = []builder{&b826, &b166, &b167}
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
	b174.items = []builder{&b826, &b14}
	b175.items = []builder{&b826, &b14, &b174}
	b176.items = []builder{&b173, &b175, &b826, &b105}
	b195.items = []builder{&b115, &b826, &b176}
	var b197 = sequenceBuilder{id: 197, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b196 = charBuilder{}
	b197.items = []builder{&b196}
	var b200 = sequenceBuilder{id: 200, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b199 = sequenceBuilder{id: 199, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b199.items = []builder{&b826, &b14}
	b200.items = []builder{&b826, &b14, &b199}
	var b198 = choiceBuilder{id: 198, commit: 2}
	var b784 = choiceBuilder{id: 784, commit: 66}
	var b503 = sequenceBuilder{id: 503, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b502 = sequenceBuilder{id: 502, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b498 = charBuilder{}
	var b499 = charBuilder{}
	var b500 = charBuilder{}
	var b501 = charBuilder{}
	b502.items = []builder{&b498, &b499, &b500, &b501}
	b503.items = []builder{&b502, &b826, &b267, &b826, &b267}
	var b547 = sequenceBuilder{id: 547, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b544 = sequenceBuilder{id: 544, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b542 = charBuilder{}
	var b543 = charBuilder{}
	b544.items = []builder{&b542, &b543}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b545.items = []builder{&b826, &b14}
	b546.items = []builder{&b826, &b14, &b545}
	var b257 = sequenceBuilder{id: 257, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b254 = sequenceBuilder{id: 254, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b253 = charBuilder{}
	b254.items = []builder{&b253}
	var b256 = sequenceBuilder{id: 256, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b255 = charBuilder{}
	b256.items = []builder{&b255}
	b257.items = []builder{&b267, &b826, &b254, &b826, &b115, &b826, &b120, &b826, &b115, &b826, &b256}
	b547.items = []builder{&b544, &b546, &b826, &b257}
	var b556 = sequenceBuilder{id: 556, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b553 = sequenceBuilder{id: 553, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b548 = charBuilder{}
	var b549 = charBuilder{}
	var b550 = charBuilder{}
	var b551 = charBuilder{}
	var b552 = charBuilder{}
	b553.items = []builder{&b548, &b549, &b550, &b551, &b552}
	var b555 = sequenceBuilder{id: 555, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b554 = sequenceBuilder{id: 554, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b554.items = []builder{&b826, &b14}
	b555.items = []builder{&b826, &b14, &b554}
	b556.items = []builder{&b553, &b555, &b826, &b257}
	var b621 = choiceBuilder{id: 621, commit: 64, name: "assignment"}
	var b601 = sequenceBuilder{id: 601, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b598 = sequenceBuilder{id: 598, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b595 = charBuilder{}
	var b596 = charBuilder{}
	var b597 = charBuilder{}
	b598.items = []builder{&b595, &b596, &b597}
	var b600 = sequenceBuilder{id: 600, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b599 = sequenceBuilder{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b599.items = []builder{&b826, &b14}
	b600.items = []builder{&b826, &b14, &b599}
	var b590 = sequenceBuilder{id: 590, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b585 = sequenceBuilder{id: 585, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b585.items = []builder{&b826, &b14}
	b586.items = []builder{&b14, &b585}
	var b584 = sequenceBuilder{id: 584, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b583 = charBuilder{}
	b584.items = []builder{&b583}
	b587.items = []builder{&b586, &b826, &b584}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b588.items = []builder{&b826, &b14}
	b589.items = []builder{&b826, &b14, &b588}
	b590.items = []builder{&b267, &b826, &b587, &b589, &b826, &b375}
	b601.items = []builder{&b598, &b600, &b826, &b590}
	var b608 = sequenceBuilder{id: 608, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b605 = sequenceBuilder{id: 605, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b604 = sequenceBuilder{id: 604, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b604.items = []builder{&b826, &b14}
	b605.items = []builder{&b826, &b14, &b604}
	var b603 = sequenceBuilder{id: 603, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b602 = charBuilder{}
	b603.items = []builder{&b602}
	var b607 = sequenceBuilder{id: 607, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b606 = sequenceBuilder{id: 606, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b606.items = []builder{&b826, &b14}
	b607.items = []builder{&b826, &b14, &b606}
	b608.items = []builder{&b267, &b605, &b826, &b603, &b607, &b826, &b375}
	var b620 = sequenceBuilder{id: 620, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b612 = sequenceBuilder{id: 612, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b609 = charBuilder{}
	var b610 = charBuilder{}
	var b611 = charBuilder{}
	b612.items = []builder{&b609, &b610, &b611}
	var b619 = sequenceBuilder{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b618 = sequenceBuilder{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b618.items = []builder{&b826, &b14}
	b619.items = []builder{&b826, &b14, &b618}
	var b614 = sequenceBuilder{id: 614, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b613 = charBuilder{}
	b614.items = []builder{&b613}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b594 = sequenceBuilder{id: 594, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b591 = sequenceBuilder{id: 591, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b591.items = []builder{&b115, &b826, &b590}
	var b592 = sequenceBuilder{id: 592, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b592.items = []builder{&b826, &b591}
	b593.items = []builder{&b826, &b591, &b592}
	b594.items = []builder{&b590, &b593}
	b615.items = []builder{&b115, &b826, &b594}
	var b617 = sequenceBuilder{id: 617, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b616 = charBuilder{}
	b617.items = []builder{&b616}
	b620.items = []builder{&b612, &b619, &b826, &b614, &b826, &b615, &b826, &b115, &b826, &b617}
	b621.options = []builder{&b601, &b608, &b620}
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
	b784.options = []builder{&b503, &b547, &b556, &b621, &b793, &b375}
	var b192 = sequenceBuilder{id: 192, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b189 = sequenceBuilder{id: 189, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b188 = charBuilder{}
	b189.items = []builder{&b188}
	var b191 = sequenceBuilder{id: 191, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b190 = charBuilder{}
	b191.items = []builder{&b190}
	b192.items = []builder{&b189, &b826, &b808, &b826, &b812, &b826, &b808, &b826, &b191}
	b198.options = []builder{&b784, &b192}
	b201.items = []builder{&b194, &b826, &b115, &b826, &b169, &b826, &b195, &b826, &b115, &b826, &b197, &b200, &b826, &b198}
	b207.items = []builder{&b204, &b206, &b826, &b201}
	var b217 = sequenceBuilder{id: 217, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b210 = sequenceBuilder{id: 210, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b208 = charBuilder{}
	var b209 = charBuilder{}
	b210.items = []builder{&b208, &b209}
	var b214 = sequenceBuilder{id: 214, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b213 = sequenceBuilder{id: 213, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b213.items = []builder{&b826, &b14}
	b214.items = []builder{&b826, &b14, &b213}
	var b212 = sequenceBuilder{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b211 = charBuilder{}
	b212.items = []builder{&b211}
	var b216 = sequenceBuilder{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b215.items = []builder{&b826, &b14}
	b216.items = []builder{&b826, &b14, &b215}
	b217.items = []builder{&b210, &b214, &b826, &b212, &b216, &b826, &b201}
	var b245 = choiceBuilder{id: 245, commit: 64, name: "expression-indexer"}
	var b235 = sequenceBuilder{id: 235, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b228 = sequenceBuilder{id: 228, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b227 = charBuilder{}
	b228.items = []builder{&b227}
	var b232 = sequenceBuilder{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b231.items = []builder{&b826, &b14}
	b232.items = []builder{&b826, &b14, &b231}
	var b234 = sequenceBuilder{id: 234, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b233 = sequenceBuilder{id: 233, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b233.items = []builder{&b826, &b14}
	b234.items = []builder{&b826, &b14, &b233}
	var b230 = sequenceBuilder{id: 230, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b229 = charBuilder{}
	b230.items = []builder{&b229}
	b235.items = []builder{&b267, &b826, &b228, &b232, &b826, &b375, &b234, &b826, &b230}
	var b244 = sequenceBuilder{id: 244, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b237 = sequenceBuilder{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b236 = charBuilder{}
	b237.items = []builder{&b236}
	var b241 = sequenceBuilder{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b240 = sequenceBuilder{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b240.items = []builder{&b826, &b14}
	b241.items = []builder{&b826, &b14, &b240}
	var b226 = sequenceBuilder{id: 226, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b218 = sequenceBuilder{id: 218, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b218.items = []builder{&b375}
	var b223 = sequenceBuilder{id: 223, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b222 = sequenceBuilder{id: 222, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b222.items = []builder{&b826, &b14}
	b223.items = []builder{&b826, &b14, &b222}
	var b221 = sequenceBuilder{id: 221, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b220 = charBuilder{}
	b221.items = []builder{&b220}
	var b225 = sequenceBuilder{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b224.items = []builder{&b826, &b14}
	b225.items = []builder{&b826, &b14, &b224}
	var b219 = sequenceBuilder{id: 219, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b219.items = []builder{&b375}
	b226.items = []builder{&b218, &b223, &b826, &b221, &b225, &b826, &b219}
	var b243 = sequenceBuilder{id: 243, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b242 = sequenceBuilder{id: 242, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b242.items = []builder{&b826, &b14}
	b243.items = []builder{&b826, &b14, &b242}
	var b239 = sequenceBuilder{id: 239, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b238 = charBuilder{}
	b239.items = []builder{&b238}
	b244.items = []builder{&b267, &b826, &b237, &b241, &b826, &b226, &b243, &b826, &b239}
	b245.options = []builder{&b235, &b244}
	var b252 = sequenceBuilder{id: 252, commit: 64, name: "symbol-indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b248.items = []builder{&b826, &b14}
	b249.items = []builder{&b826, &b14, &b248}
	var b247 = sequenceBuilder{id: 247, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b246 = charBuilder{}
	b247.items = []builder{&b246}
	var b251 = sequenceBuilder{id: 251, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b250.items = []builder{&b826, &b14}
	b251.items = []builder{&b826, &b14, &b250}
	b252.items = []builder{&b267, &b249, &b826, &b247, &b251, &b826, &b105}
	var b266 = sequenceBuilder{id: 266, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b259 = sequenceBuilder{id: 259, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b258 = charBuilder{}
	b259.items = []builder{&b258}
	var b263 = sequenceBuilder{id: 263, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b262 = sequenceBuilder{id: 262, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b262.items = []builder{&b826, &b14}
	b263.items = []builder{&b826, &b14, &b262}
	var b265 = sequenceBuilder{id: 265, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b264 = sequenceBuilder{id: 264, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b264.items = []builder{&b826, &b14}
	b265.items = []builder{&b826, &b14, &b264}
	var b261 = sequenceBuilder{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b260 = charBuilder{}
	b261.items = []builder{&b260}
	b266.items = []builder{&b259, &b263, &b826, &b375, &b265, &b826, &b261}
	b267.options = []builder{&b62, &b75, &b88, &b100, &b478, &b105, &b126, &b131, &b160, &b165, &b207, &b217, &b245, &b252, &b257, &b266}
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
	b327.items = []builder{&b326, &b826, &b267}
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
	b339.items = []builder{&b328, &b826, &b333}
	var b340 = sequenceBuilder{id: 340, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b340.items = []builder{&b826, &b339}
	b341.items = []builder{&b333, &b826, &b339, &b340}
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
	b342.items = []builder{&b329, &b826, &b334}
	var b343 = sequenceBuilder{id: 343, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b343.items = []builder{&b826, &b342}
	b344.items = []builder{&b334, &b826, &b342, &b343}
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
	b345.items = []builder{&b330, &b826, &b335}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b346.items = []builder{&b826, &b345}
	b347.items = []builder{&b335, &b826, &b345, &b346}
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
	b348.items = []builder{&b331, &b826, &b336}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b349.items = []builder{&b826, &b348}
	b350.items = []builder{&b336, &b826, &b348, &b349}
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
	b351.items = []builder{&b332, &b826, &b337}
	var b352 = sequenceBuilder{id: 352, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b352.items = []builder{&b826, &b351}
	b353.items = []builder{&b337, &b826, &b351, &b352}
	var b360 = sequenceBuilder{id: 360, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b338 = choiceBuilder{id: 338, commit: 66}
	b338.options = []builder{&b337, &b353}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b354.items = []builder{&b826, &b14}
	b355.items = []builder{&b14, &b354}
	var b325 = sequenceBuilder{id: 325, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b323 = charBuilder{}
	var b324 = charBuilder{}
	b325.items = []builder{&b323, &b324}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b356.items = []builder{&b826, &b14}
	b357.items = []builder{&b826, &b14, &b356}
	b358.items = []builder{&b355, &b826, &b325, &b357, &b826, &b338}
	var b359 = sequenceBuilder{id: 359, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b359.items = []builder{&b826, &b358}
	b360.items = []builder{&b338, &b826, &b358, &b359}
	b361.options = []builder{&b341, &b344, &b347, &b350, &b353, &b360}
	var b374 = sequenceBuilder{id: 374, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b366 = sequenceBuilder{id: 366, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b366.items = []builder{&b826, &b14}
	b367.items = []builder{&b826, &b14, &b366}
	var b363 = sequenceBuilder{id: 363, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b362 = charBuilder{}
	b363.items = []builder{&b362}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b368.items = []builder{&b826, &b14}
	b369.items = []builder{&b826, &b14, &b368}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b370.items = []builder{&b826, &b14}
	b371.items = []builder{&b826, &b14, &b370}
	var b365 = sequenceBuilder{id: 365, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b364 = charBuilder{}
	b365.items = []builder{&b364}
	var b373 = sequenceBuilder{id: 373, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b372.items = []builder{&b826, &b14}
	b373.items = []builder{&b826, &b14, &b372}
	b374.items = []builder{&b375, &b367, &b826, &b363, &b369, &b826, &b375, &b371, &b826, &b365, &b373, &b826, &b375}
	b375.options = []builder{&b267, &b327, &b361, &b374}
	b186.items = []builder{&b185, &b826, &b375}
	b187.items = []builder{&b183, &b826, &b186}
	var b412 = sequenceBuilder{id: 412, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b378 = sequenceBuilder{id: 378, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b376 = charBuilder{}
	var b377 = charBuilder{}
	b378.items = []builder{&b376, &b377}
	var b407 = sequenceBuilder{id: 407, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b406 = sequenceBuilder{id: 406, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b406.items = []builder{&b826, &b14}
	b407.items = []builder{&b826, &b14, &b406}
	var b409 = sequenceBuilder{id: 409, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b408 = sequenceBuilder{id: 408, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b408.items = []builder{&b826, &b14}
	b409.items = []builder{&b826, &b14, &b408}
	var b411 = sequenceBuilder{id: 411, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b395 = sequenceBuilder{id: 395, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b387.items = []builder{&b826, &b14}
	b388.items = []builder{&b14, &b387}
	var b383 = sequenceBuilder{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b379 = charBuilder{}
	var b380 = charBuilder{}
	var b381 = charBuilder{}
	var b382 = charBuilder{}
	b383.items = []builder{&b379, &b380, &b381, &b382}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b389.items = []builder{&b826, &b14}
	b390.items = []builder{&b826, &b14, &b389}
	var b386 = sequenceBuilder{id: 386, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b384 = charBuilder{}
	var b385 = charBuilder{}
	b386.items = []builder{&b384, &b385}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b391.items = []builder{&b826, &b14}
	b392.items = []builder{&b826, &b14, &b391}
	var b394 = sequenceBuilder{id: 394, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b393.items = []builder{&b826, &b14}
	b394.items = []builder{&b826, &b14, &b393}
	b395.items = []builder{&b388, &b826, &b383, &b390, &b826, &b386, &b392, &b826, &b375, &b394, &b826, &b192}
	var b410 = sequenceBuilder{id: 410, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b410.items = []builder{&b826, &b395}
	b411.items = []builder{&b826, &b395, &b410}
	var b405 = sequenceBuilder{id: 405, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b402 = sequenceBuilder{id: 402, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b401 = sequenceBuilder{id: 401, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b401.items = []builder{&b826, &b14}
	b402.items = []builder{&b14, &b401}
	var b400 = sequenceBuilder{id: 400, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b396 = charBuilder{}
	var b397 = charBuilder{}
	var b398 = charBuilder{}
	var b399 = charBuilder{}
	b400.items = []builder{&b396, &b397, &b398, &b399}
	var b404 = sequenceBuilder{id: 404, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b403 = sequenceBuilder{id: 403, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b403.items = []builder{&b826, &b14}
	b404.items = []builder{&b826, &b14, &b403}
	b405.items = []builder{&b402, &b826, &b400, &b404, &b826, &b192}
	b412.items = []builder{&b378, &b407, &b826, &b375, &b409, &b826, &b192, &b411, &b826, &b405}
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
	b465.items = []builder{&b826, &b14}
	b466.items = []builder{&b826, &b14, &b465}
	var b468 = sequenceBuilder{id: 468, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b467 = sequenceBuilder{id: 467, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b467.items = []builder{&b826, &b14}
	b468.items = []builder{&b826, &b14, &b467}
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
	b438.items = []builder{&b826, &b14}
	b439.items = []builder{&b826, &b14, &b438}
	var b441 = sequenceBuilder{id: 441, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b440 = sequenceBuilder{id: 440, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b440.items = []builder{&b826, &b14}
	b441.items = []builder{&b826, &b14, &b440}
	var b437 = sequenceBuilder{id: 437, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b436 = charBuilder{}
	b437.items = []builder{&b436}
	b442.items = []builder{&b435, &b439, &b826, &b375, &b441, &b826, &b437}
	var b446 = sequenceBuilder{id: 446, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b444 = sequenceBuilder{id: 444, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b443 = charBuilder{}
	b444.items = []builder{&b443}
	var b445 = sequenceBuilder{id: 445, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b445.items = []builder{&b826, &b444}
	b446.items = []builder{&b826, &b444, &b445}
	b447.items = []builder{&b442, &b446, &b826, &b794}
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
	b423.items = []builder{&b826, &b14}
	b424.items = []builder{&b826, &b14, &b423}
	var b422 = sequenceBuilder{id: 422, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b421 = charBuilder{}
	b422.items = []builder{&b421}
	b425.items = []builder{&b420, &b424, &b826, &b422}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b427 = sequenceBuilder{id: 427, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b426 = charBuilder{}
	b427.items = []builder{&b426}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b428.items = []builder{&b826, &b427}
	b429.items = []builder{&b826, &b427, &b428}
	b430.items = []builder{&b425, &b429, &b826, &b794}
	b457.options = []builder{&b447, &b430}
	var b461 = sequenceBuilder{id: 461, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b458 = choiceBuilder{id: 458, commit: 2}
	b458.options = []builder{&b447, &b430, &b794}
	b459.items = []builder{&b808, &b826, &b458}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b460.items = []builder{&b826, &b459}
	b461.items = []builder{&b826, &b459, &b460}
	b462.items = []builder{&b457, &b461}
	var b464 = sequenceBuilder{id: 464, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b463 = charBuilder{}
	b464.items = []builder{&b463}
	b469.items = []builder{&b454, &b466, &b826, &b375, &b468, &b826, &b456, &b826, &b808, &b826, &b462, &b826, &b808, &b826, &b464}
	var b541 = sequenceBuilder{id: 541, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b528 = sequenceBuilder{id: 528, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b522 = charBuilder{}
	var b523 = charBuilder{}
	var b524 = charBuilder{}
	var b525 = charBuilder{}
	var b526 = charBuilder{}
	var b527 = charBuilder{}
	b528.items = []builder{&b522, &b523, &b524, &b525, &b526, &b527}
	var b540 = sequenceBuilder{id: 540, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b539 = sequenceBuilder{id: 539, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b539.items = []builder{&b826, &b14}
	b540.items = []builder{&b826, &b14, &b539}
	var b530 = sequenceBuilder{id: 530, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b529 = charBuilder{}
	b530.items = []builder{&b529}
	var b536 = sequenceBuilder{id: 536, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b531 = choiceBuilder{id: 531, commit: 2}
	var b521 = sequenceBuilder{id: 521, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b516 = sequenceBuilder{id: 516, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b509 = sequenceBuilder{id: 509, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b505 = charBuilder{}
	var b506 = charBuilder{}
	var b507 = charBuilder{}
	var b508 = charBuilder{}
	b509.items = []builder{&b505, &b506, &b507, &b508}
	var b513 = sequenceBuilder{id: 513, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b512 = sequenceBuilder{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b512.items = []builder{&b826, &b14}
	b513.items = []builder{&b826, &b14, &b512}
	var b504 = choiceBuilder{id: 504, commit: 66}
	var b497 = choiceBuilder{id: 497, commit: 66}
	var b487 = sequenceBuilder{id: 487, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b482 = sequenceBuilder{id: 482, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b479 = charBuilder{}
	var b480 = charBuilder{}
	var b481 = charBuilder{}
	b482.items = []builder{&b479, &b480, &b481}
	var b484 = sequenceBuilder{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b483 = sequenceBuilder{id: 483, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b483.items = []builder{&b826, &b14}
	b484.items = []builder{&b826, &b14, &b483}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b485.items = []builder{&b826, &b14}
	b486.items = []builder{&b826, &b14, &b485}
	b487.items = []builder{&b482, &b484, &b826, &b105, &b486, &b826, &b478}
	var b496 = sequenceBuilder{id: 496, commit: 64, name: "receive-assignment", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b491 = sequenceBuilder{id: 491, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b488 = charBuilder{}
	var b489 = charBuilder{}
	var b490 = charBuilder{}
	b491.items = []builder{&b488, &b489, &b490}
	var b493 = sequenceBuilder{id: 493, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b492.items = []builder{&b826, &b14}
	b493.items = []builder{&b826, &b14, &b492}
	var b495 = sequenceBuilder{id: 495, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b494 = sequenceBuilder{id: 494, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b494.items = []builder{&b826, &b14}
	b495.items = []builder{&b826, &b14, &b494}
	b496.items = []builder{&b491, &b493, &b826, &b105, &b495, &b826, &b478}
	b497.options = []builder{&b487, &b496}
	b504.options = []builder{&b478, &b497, &b503}
	var b515 = sequenceBuilder{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b514 = sequenceBuilder{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b514.items = []builder{&b826, &b14}
	b515.items = []builder{&b826, &b14, &b514}
	var b511 = sequenceBuilder{id: 511, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b510 = charBuilder{}
	b511.items = []builder{&b510}
	b516.items = []builder{&b509, &b513, &b826, &b504, &b515, &b826, &b511}
	var b520 = sequenceBuilder{id: 520, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b518 = sequenceBuilder{id: 518, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b517 = charBuilder{}
	b518.items = []builder{&b517}
	var b519 = sequenceBuilder{id: 519, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b519.items = []builder{&b826, &b518}
	b520.items = []builder{&b826, &b518, &b519}
	b521.items = []builder{&b516, &b520, &b826, &b794}
	b531.options = []builder{&b521, &b430}
	var b535 = sequenceBuilder{id: 535, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b533 = sequenceBuilder{id: 533, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b532 = choiceBuilder{id: 532, commit: 2}
	b532.options = []builder{&b521, &b430, &b794}
	b533.items = []builder{&b808, &b826, &b532}
	var b534 = sequenceBuilder{id: 534, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b534.items = []builder{&b826, &b533}
	b535.items = []builder{&b826, &b533, &b534}
	b536.items = []builder{&b531, &b535}
	var b538 = sequenceBuilder{id: 538, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b537 = charBuilder{}
	b538.items = []builder{&b537}
	b541.items = []builder{&b528, &b540, &b826, &b530, &b826, &b808, &b826, &b536, &b826, &b808, &b826, &b538}
	var b582 = sequenceBuilder{id: 582, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b571 = sequenceBuilder{id: 571, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b568 = charBuilder{}
	var b569 = charBuilder{}
	var b570 = charBuilder{}
	b571.items = []builder{&b568, &b569, &b570}
	var b581 = choiceBuilder{id: 581, commit: 2}
	var b577 = sequenceBuilder{id: 577, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b574 = sequenceBuilder{id: 574, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b572 = sequenceBuilder{id: 572, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b572.items = []builder{&b826, &b14}
	b573.items = []builder{&b14, &b572}
	var b567 = choiceBuilder{id: 567, commit: 66}
	var b566 = choiceBuilder{id: 566, commit: 64, name: "range-over-expression"}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b562 = sequenceBuilder{id: 562, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b561 = sequenceBuilder{id: 561, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b561.items = []builder{&b826, &b14}
	b562.items = []builder{&b826, &b14, &b561}
	var b559 = sequenceBuilder{id: 559, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b557 = charBuilder{}
	var b558 = charBuilder{}
	b559.items = []builder{&b557, &b558}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b563 = sequenceBuilder{id: 563, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b563.items = []builder{&b826, &b14}
	b564.items = []builder{&b826, &b14, &b563}
	var b560 = choiceBuilder{id: 560, commit: 2}
	b560.options = []builder{&b375, &b226}
	b565.items = []builder{&b105, &b562, &b826, &b559, &b564, &b826, &b560}
	b566.options = []builder{&b565, &b226}
	b567.options = []builder{&b375, &b566}
	b574.items = []builder{&b573, &b826, &b567}
	var b576 = sequenceBuilder{id: 576, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b575 = sequenceBuilder{id: 575, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b575.items = []builder{&b826, &b14}
	b576.items = []builder{&b826, &b14, &b575}
	b577.items = []builder{&b574, &b576, &b826, &b192}
	var b580 = sequenceBuilder{id: 580, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b579 = sequenceBuilder{id: 579, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b578 = sequenceBuilder{id: 578, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b578.items = []builder{&b826, &b14}
	b579.items = []builder{&b14, &b578}
	b580.items = []builder{&b579, &b826, &b192}
	b581.options = []builder{&b577, &b580}
	b582.items = []builder{&b571, &b826, &b581}
	var b730 = choiceBuilder{id: 730, commit: 66}
	var b643 = sequenceBuilder{id: 643, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b639 = sequenceBuilder{id: 639, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b636 = charBuilder{}
	var b637 = charBuilder{}
	var b638 = charBuilder{}
	b639.items = []builder{&b636, &b637, &b638}
	var b642 = sequenceBuilder{id: 642, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b641 = sequenceBuilder{id: 641, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b641.items = []builder{&b826, &b14}
	b642.items = []builder{&b826, &b14, &b641}
	var b640 = choiceBuilder{id: 640, commit: 2}
	var b630 = sequenceBuilder{id: 630, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b629 = sequenceBuilder{id: 629, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b626 = sequenceBuilder{id: 626, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b625 = sequenceBuilder{id: 625, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b624 = sequenceBuilder{id: 624, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b624.items = []builder{&b826, &b14}
	b625.items = []builder{&b14, &b624}
	var b623 = sequenceBuilder{id: 623, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b622 = charBuilder{}
	b623.items = []builder{&b622}
	b626.items = []builder{&b625, &b826, &b623}
	var b628 = sequenceBuilder{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b627 = sequenceBuilder{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b627.items = []builder{&b826, &b14}
	b628.items = []builder{&b826, &b14, &b627}
	b629.items = []builder{&b105, &b826, &b626, &b628, &b826, &b375}
	b630.items = []builder{&b629}
	var b635 = sequenceBuilder{id: 635, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b632 = sequenceBuilder{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b631 = charBuilder{}
	b632.items = []builder{&b631}
	var b634 = sequenceBuilder{id: 634, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b633 = sequenceBuilder{id: 633, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b633.items = []builder{&b826, &b14}
	b634.items = []builder{&b826, &b14, &b633}
	b635.items = []builder{&b632, &b634, &b826, &b629}
	b640.options = []builder{&b630, &b635}
	b643.items = []builder{&b639, &b642, &b826, &b640}
	var b664 = sequenceBuilder{id: 664, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b657 = sequenceBuilder{id: 657, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b654 = charBuilder{}
	var b655 = charBuilder{}
	var b656 = charBuilder{}
	b657.items = []builder{&b654, &b655, &b656}
	var b663 = sequenceBuilder{id: 663, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b662 = sequenceBuilder{id: 662, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b662.items = []builder{&b826, &b14}
	b663.items = []builder{&b826, &b14, &b662}
	var b659 = sequenceBuilder{id: 659, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b658 = charBuilder{}
	b659.items = []builder{&b658}
	var b649 = sequenceBuilder{id: 649, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b644 = choiceBuilder{id: 644, commit: 2}
	b644.options = []builder{&b630, &b635}
	var b648 = sequenceBuilder{id: 648, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b646 = sequenceBuilder{id: 646, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b645 = choiceBuilder{id: 645, commit: 2}
	b645.options = []builder{&b630, &b635}
	b646.items = []builder{&b115, &b826, &b645}
	var b647 = sequenceBuilder{id: 647, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b647.items = []builder{&b826, &b646}
	b648.items = []builder{&b826, &b646, &b647}
	b649.items = []builder{&b644, &b648}
	var b661 = sequenceBuilder{id: 661, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b660 = charBuilder{}
	b661.items = []builder{&b660}
	b664.items = []builder{&b657, &b663, &b826, &b659, &b826, &b115, &b826, &b649, &b826, &b115, &b826, &b661}
	var b679 = sequenceBuilder{id: 679, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b668 = sequenceBuilder{id: 668, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b665 = charBuilder{}
	var b666 = charBuilder{}
	var b667 = charBuilder{}
	b668.items = []builder{&b665, &b666, &b667}
	var b676 = sequenceBuilder{id: 676, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b675 = sequenceBuilder{id: 675, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b675.items = []builder{&b826, &b14}
	b676.items = []builder{&b826, &b14, &b675}
	var b670 = sequenceBuilder{id: 670, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b669 = charBuilder{}
	b670.items = []builder{&b669}
	var b678 = sequenceBuilder{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b677 = sequenceBuilder{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b677.items = []builder{&b826, &b14}
	b678.items = []builder{&b826, &b14, &b677}
	var b672 = sequenceBuilder{id: 672, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b671 = charBuilder{}
	b672.items = []builder{&b671}
	var b653 = sequenceBuilder{id: 653, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b652 = sequenceBuilder{id: 652, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b650 = sequenceBuilder{id: 650, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b650.items = []builder{&b115, &b826, &b630}
	var b651 = sequenceBuilder{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b651.items = []builder{&b826, &b650}
	b652.items = []builder{&b826, &b650, &b651}
	b653.items = []builder{&b630, &b652}
	var b674 = sequenceBuilder{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b673 = charBuilder{}
	b674.items = []builder{&b673}
	b679.items = []builder{&b668, &b676, &b826, &b670, &b678, &b826, &b672, &b826, &b115, &b826, &b653, &b826, &b115, &b826, &b674}
	var b695 = sequenceBuilder{id: 695, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b691 = sequenceBuilder{id: 691, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b689 = charBuilder{}
	var b690 = charBuilder{}
	b691.items = []builder{&b689, &b690}
	var b694 = sequenceBuilder{id: 694, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b693 = sequenceBuilder{id: 693, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b693.items = []builder{&b826, &b14}
	b694.items = []builder{&b826, &b14, &b693}
	var b692 = choiceBuilder{id: 692, commit: 2}
	var b683 = sequenceBuilder{id: 683, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b682 = sequenceBuilder{id: 682, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b681 = sequenceBuilder{id: 681, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b680 = sequenceBuilder{id: 680, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b680.items = []builder{&b826, &b14}
	b681.items = []builder{&b826, &b14, &b680}
	b682.items = []builder{&b105, &b681, &b826, &b201}
	b683.items = []builder{&b682}
	var b688 = sequenceBuilder{id: 688, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b685 = sequenceBuilder{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b684 = charBuilder{}
	b685.items = []builder{&b684}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b686.items = []builder{&b826, &b14}
	b687.items = []builder{&b826, &b14, &b686}
	b688.items = []builder{&b685, &b687, &b826, &b682}
	b692.options = []builder{&b683, &b688}
	b695.items = []builder{&b691, &b694, &b826, &b692}
	var b715 = sequenceBuilder{id: 715, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b708 = sequenceBuilder{id: 708, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b706 = charBuilder{}
	var b707 = charBuilder{}
	b708.items = []builder{&b706, &b707}
	var b714 = sequenceBuilder{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b713 = sequenceBuilder{id: 713, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b713.items = []builder{&b826, &b14}
	b714.items = []builder{&b826, &b14, &b713}
	var b710 = sequenceBuilder{id: 710, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b709 = charBuilder{}
	b710.items = []builder{&b709}
	var b705 = sequenceBuilder{id: 705, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b700 = choiceBuilder{id: 700, commit: 2}
	b700.options = []builder{&b683, &b688}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b702 = sequenceBuilder{id: 702, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b701 = choiceBuilder{id: 701, commit: 2}
	b701.options = []builder{&b683, &b688}
	b702.items = []builder{&b115, &b826, &b701}
	var b703 = sequenceBuilder{id: 703, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b703.items = []builder{&b826, &b702}
	b704.items = []builder{&b826, &b702, &b703}
	b705.items = []builder{&b700, &b704}
	var b712 = sequenceBuilder{id: 712, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b711 = charBuilder{}
	b712.items = []builder{&b711}
	b715.items = []builder{&b708, &b714, &b826, &b710, &b826, &b115, &b826, &b705, &b826, &b115, &b826, &b712}
	var b729 = sequenceBuilder{id: 729, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b718 = sequenceBuilder{id: 718, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b716 = charBuilder{}
	var b717 = charBuilder{}
	b718.items = []builder{&b716, &b717}
	var b726 = sequenceBuilder{id: 726, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b725 = sequenceBuilder{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b725.items = []builder{&b826, &b14}
	b726.items = []builder{&b826, &b14, &b725}
	var b720 = sequenceBuilder{id: 720, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b719 = charBuilder{}
	b720.items = []builder{&b719}
	var b728 = sequenceBuilder{id: 728, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b727 = sequenceBuilder{id: 727, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b727.items = []builder{&b826, &b14}
	b728.items = []builder{&b826, &b14, &b727}
	var b722 = sequenceBuilder{id: 722, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b721 = charBuilder{}
	b722.items = []builder{&b721}
	var b699 = sequenceBuilder{id: 699, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b698 = sequenceBuilder{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b696 = sequenceBuilder{id: 696, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b696.items = []builder{&b115, &b826, &b683}
	var b697 = sequenceBuilder{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b697.items = []builder{&b826, &b696}
	b698.items = []builder{&b826, &b696, &b697}
	b699.items = []builder{&b683, &b698}
	var b724 = sequenceBuilder{id: 724, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b723 = charBuilder{}
	b724.items = []builder{&b723}
	b729.items = []builder{&b718, &b726, &b826, &b720, &b728, &b826, &b722, &b826, &b115, &b826, &b699, &b826, &b115, &b826, &b724}
	b730.options = []builder{&b643, &b664, &b679, &b695, &b715, &b729}
	var b773 = choiceBuilder{id: 773, commit: 64, name: "require"}
	var b757 = sequenceBuilder{id: 757, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b754 = sequenceBuilder{id: 754, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b747 = charBuilder{}
	var b748 = charBuilder{}
	var b749 = charBuilder{}
	var b750 = charBuilder{}
	var b751 = charBuilder{}
	var b752 = charBuilder{}
	var b753 = charBuilder{}
	b754.items = []builder{&b747, &b748, &b749, &b750, &b751, &b752, &b753}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b755.items = []builder{&b826, &b14}
	b756.items = []builder{&b826, &b14, &b755}
	var b742 = choiceBuilder{id: 742, commit: 64, name: "require-fact"}
	var b741 = sequenceBuilder{id: 741, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b733 = choiceBuilder{id: 733, commit: 2}
	var b732 = sequenceBuilder{id: 732, commit: 72, name: "require-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b731 = charBuilder{}
	b732.items = []builder{&b731}
	b733.options = []builder{&b105, &b732}
	var b738 = sequenceBuilder{id: 738, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b736.items = []builder{&b826, &b14}
	b737.items = []builder{&b14, &b736}
	var b735 = sequenceBuilder{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b734 = charBuilder{}
	b735.items = []builder{&b734}
	b738.items = []builder{&b737, &b826, &b735}
	var b740 = sequenceBuilder{id: 740, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b739 = sequenceBuilder{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b739.items = []builder{&b826, &b14}
	b740.items = []builder{&b826, &b14, &b739}
	b741.items = []builder{&b733, &b826, &b738, &b740, &b826, &b88}
	b742.options = []builder{&b88, &b741}
	b757.items = []builder{&b754, &b756, &b826, &b742}
	var b772 = sequenceBuilder{id: 772, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b765 = sequenceBuilder{id: 765, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b758 = charBuilder{}
	var b759 = charBuilder{}
	var b760 = charBuilder{}
	var b761 = charBuilder{}
	var b762 = charBuilder{}
	var b763 = charBuilder{}
	var b764 = charBuilder{}
	b765.items = []builder{&b758, &b759, &b760, &b761, &b762, &b763, &b764}
	var b771 = sequenceBuilder{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b770 = sequenceBuilder{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b770.items = []builder{&b826, &b14}
	b771.items = []builder{&b826, &b14, &b770}
	var b767 = sequenceBuilder{id: 767, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b766 = charBuilder{}
	b767.items = []builder{&b766}
	var b746 = sequenceBuilder{id: 746, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b745 = sequenceBuilder{id: 745, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b743 = sequenceBuilder{id: 743, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b743.items = []builder{&b115, &b826, &b742}
	var b744 = sequenceBuilder{id: 744, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b744.items = []builder{&b826, &b743}
	b745.items = []builder{&b826, &b743, &b744}
	b746.items = []builder{&b742, &b745}
	var b769 = sequenceBuilder{id: 769, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b768 = charBuilder{}
	b769.items = []builder{&b768}
	b772.items = []builder{&b765, &b771, &b826, &b767, &b826, &b115, &b826, &b746, &b826, &b115, &b826, &b769}
	b773.options = []builder{&b757, &b772}
	var b783 = sequenceBuilder{id: 783, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b780 = sequenceBuilder{id: 780, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b774 = charBuilder{}
	var b775 = charBuilder{}
	var b776 = charBuilder{}
	var b777 = charBuilder{}
	var b778 = charBuilder{}
	var b779 = charBuilder{}
	b780.items = []builder{&b774, &b775, &b776, &b777, &b778, &b779}
	var b782 = sequenceBuilder{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b781 = sequenceBuilder{id: 781, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b781.items = []builder{&b826, &b14}
	b782.items = []builder{&b826, &b14, &b781}
	b783.items = []builder{&b780, &b782, &b826, &b730}
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
	b794.options = []builder{&b187, &b412, &b469, &b541, &b582, &b730, &b773, &b783, &b803, &b784}
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
