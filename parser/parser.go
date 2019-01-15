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
	var p829 = sequenceParser{id: 829, commit: 32, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p827 = choiceParser{id: 827, commit: 2}
	var p825 = choiceParser{id: 825, commit: 70, name: "ws", generalizations: []int{827, 15}}
	var p2 = sequenceParser{id: 2, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825, 827, 15}}
	var p1 = charParser{id: 1, chars: []rune{32}}
	p2.items = []parser{&p1}
	var p4 = sequenceParser{id: 4, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825, 827, 15}}
	var p3 = charParser{id: 3, chars: []rune{8}}
	p4.items = []parser{&p3}
	var p6 = sequenceParser{id: 6, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825, 827, 15}}
	var p5 = charParser{id: 5, chars: []rune{12}}
	p6.items = []parser{&p5}
	var p8 = sequenceParser{id: 8, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825, 827, 15}}
	var p7 = charParser{id: 7, chars: []rune{13}}
	p8.items = []parser{&p7}
	var p10 = sequenceParser{id: 10, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825, 827, 15}}
	var p9 = charParser{id: 9, chars: []rune{9}}
	p10.items = []parser{&p9}
	var p12 = sequenceParser{id: 12, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{825, 827, 15}}
	var p11 = charParser{id: 11, chars: []rune{11}}
	p12.items = []parser{&p11}
	p825.options = []parser{&p2, &p4, &p6, &p8, &p10, &p12}
	var p826 = sequenceParser{id: 826, commit: 70, name: "wsc", ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{827}}
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
	var p14 = sequenceParser{id: 14, commit: 74, name: "nl", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{807, 15, 112}}
	var p13 = charParser{id: 13, chars: []rune{10}}
	p14.items = []parser{&p13}
	p40.items = []parser{&p14, &p827, &p39}
	var p41 = sequenceParser{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p41.items = []parser{&p827, &p40}
	p42.items = []parser{&p827, &p40, &p41}
	p43.items = []parser{&p39, &p42}
	p826.items = []parser{&p43}
	p827.options = []parser{&p825, &p826}
	var p828 = sequenceParser{id: 828, commit: 66, name: "mml:wsroot", ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var p824 = sequenceParser{id: 824, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var p821 = sequenceParser{id: 821, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p819 = charParser{id: 819, chars: []rune{35}}
	var p820 = charParser{id: 820, chars: []rune{33}}
	p821.items = []parser{&p819, &p820}
	var p818 = sequenceParser{id: 818, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var p817 = sequenceParser{id: 817, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p815 = sequenceParser{id: 815, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var p814 = charParser{id: 814, not: true, chars: []rune{10}}
	p815.items = []parser{&p814}
	var p816 = sequenceParser{id: 816, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p816.items = []parser{&p827, &p815}
	p817.items = []parser{&p815, &p816}
	p818.items = []parser{&p817}
	var p823 = sequenceParser{id: 823, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p822 = charParser{id: 822, chars: []rune{10}}
	p823.items = []parser{&p822}
	p824.items = []parser{&p821, &p827, &p818, &p827, &p823}
	var p809 = sequenceParser{id: 809, commit: 66, name: "sep", ranges: [][]int{{1, 1}, {0, -1}}}
	var p807 = choiceParser{id: 807, commit: 2}
	var p806 = sequenceParser{id: 806, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{807}}
	var p805 = charParser{id: 805, chars: []rune{59}}
	p806.items = []parser{&p805}
	p807.options = []parser{&p806, &p14}
	var p808 = sequenceParser{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p808.items = []parser{&p827, &p807}
	p809.items = []parser{&p807, &p808}
	var p813 = sequenceParser{id: 813, commit: 66, name: "statement-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p795 = choiceParser{id: 795, commit: 66, name: "statement", generalizations: []int{482, 543}}
	var p187 = sequenceParser{id: 187, commit: 64, name: "return", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}, generalizations: []int{795, 482, 543}}
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
	p15.options = []parser{&p825, &p14}
	p183.items = []parser{&p182, &p15}
	var p186 = sequenceParser{id: 186, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p185 = sequenceParser{id: 185, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p184 = sequenceParser{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p184.items = []parser{&p827, &p14}
	p185.items = []parser{&p14, &p184}
	var p403 = choiceParser{id: 403, commit: 66, name: "expression", generalizations: []int{115, 785, 200, 592, 585, 795}}
	var p274 = choiceParser{id: 274, commit: 66, name: "primary-expression", generalizations: []int{115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p61 = choiceParser{id: 61, commit: 64, name: "int", generalizations: []int{274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p52 = sequenceParser{id: 52, commit: 74, name: "decimal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{61, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p51 = sequenceParser{id: 51, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p50 = charParser{id: 50, ranges: [][]rune{{49, 57}}}
	p51.items = []parser{&p50}
	var p45 = sequenceParser{id: 45, commit: 66, name: "decimal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p44 = charParser{id: 44, ranges: [][]rune{{48, 57}}}
	p45.items = []parser{&p44}
	p52.items = []parser{&p51, &p45}
	var p55 = sequenceParser{id: 55, commit: 74, name: "octal", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{61, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p54 = sequenceParser{id: 54, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p53 = charParser{id: 53, chars: []rune{48}}
	p54.items = []parser{&p53}
	var p47 = sequenceParser{id: 47, commit: 66, name: "octal-digit", allChars: true, ranges: [][]int{{1, 1}}}
	var p46 = charParser{id: 46, ranges: [][]rune{{48, 55}}}
	p47.items = []parser{&p46}
	p55.items = []parser{&p54, &p47}
	var p60 = sequenceParser{id: 60, commit: 74, name: "hexa", ranges: [][]int{{1, 1}, {1, 1}, {1, -1}, {1, 1}, {1, 1}, {1, -1}}, generalizations: []int{61, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
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
	var p74 = choiceParser{id: 74, commit: 72, name: "float", generalizations: []int{274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p69 = sequenceParser{id: 69, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {0, -1}, {0, 1}}, generalizations: []int{74, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
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
	var p72 = sequenceParser{id: 72, commit: 10, ranges: [][]int{{1, 1}, {1, -1}, {0, 1}, {1, 1}, {1, -1}, {0, 1}}, generalizations: []int{74, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p71 = sequenceParser{id: 71, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p70 = charParser{id: 70, chars: []rune{46}}
	p71.items = []parser{&p70}
	p72.items = []parser{&p71, &p45, &p66}
	var p73 = sequenceParser{id: 73, commit: 10, ranges: [][]int{{1, -1}, {1, 1}, {1, -1}, {1, 1}}, generalizations: []int{74, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	p73.items = []parser{&p45, &p66}
	p74.options = []parser{&p69, &p72, &p73}
	var p87 = sequenceParser{id: 87, commit: 72, name: "string", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 115, 140, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 758, 795}}
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
	var p99 = choiceParser{id: 99, commit: 66, name: "bool", generalizations: []int{274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p92 = sequenceParser{id: 92, commit: 72, name: "true", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{99, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p88 = charParser{id: 88, chars: []rune{116}}
	var p89 = charParser{id: 89, chars: []rune{114}}
	var p90 = charParser{id: 90, chars: []rune{117}}
	var p91 = charParser{id: 91, chars: []rune{101}}
	p92.items = []parser{&p88, &p89, &p90, &p91}
	var p98 = sequenceParser{id: 98, commit: 72, name: "false", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{99, 274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p93 = charParser{id: 93, chars: []rune{102}}
	var p94 = charParser{id: 94, chars: []rune{97}}
	var p95 = charParser{id: 95, chars: []rune{108}}
	var p96 = charParser{id: 96, chars: []rune{115}}
	var p97 = charParser{id: 97, chars: []rune{101}}
	p98.items = []parser{&p93, &p94, &p95, &p96, &p97}
	p99.options = []parser{&p92, &p98}
	var p516 = sequenceParser{id: 516, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 115, 785, 200, 403, 340, 341, 342, 343, 344, 395, 520, 592, 585, 795}}
	var p508 = sequenceParser{id: 508, commit: 74, name: "receive-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p507 = sequenceParser{id: 507, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p500 = charParser{id: 500, chars: []rune{114}}
	var p501 = charParser{id: 501, chars: []rune{101}}
	var p502 = charParser{id: 502, chars: []rune{99}}
	var p503 = charParser{id: 503, chars: []rune{101}}
	var p504 = charParser{id: 504, chars: []rune{105}}
	var p505 = charParser{id: 505, chars: []rune{118}}
	var p506 = charParser{id: 506, chars: []rune{101}}
	p507.items = []parser{&p500, &p501, &p502, &p503, &p504, &p505, &p506}
	p508.items = []parser{&p507, &p15}
	var p515 = sequenceParser{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p514 = sequenceParser{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p514.items = []parser{&p827, &p14}
	p515.items = []parser{&p827, &p14, &p514}
	p516.items = []parser{&p508, &p515, &p827, &p274}
	var p104 = sequenceParser{id: 104, commit: 72, name: "symbol", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{274, 115, 140, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 749, 795}}
	var p101 = sequenceParser{id: 101, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p100 = charParser{id: 100, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}}}
	p101.items = []parser{&p100}
	var p103 = sequenceParser{id: 103, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p102 = charParser{id: 102, chars: []rune{95}, ranges: [][]rune{{97, 122}, {65, 90}, {48, 57}}}
	p103.items = []parser{&p102}
	p104.items = []parser{&p101, &p103}
	var p125 = sequenceParser{id: 125, commit: 64, name: "list", ranges: [][]int{{1, 1}}, generalizations: []int{115, 274, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
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
	p113.items = []parser{&p827, &p112}
	p114.items = []parser{&p112, &p113}
	var p119 = sequenceParser{id: 119, commit: 66, name: "expression-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p115 = choiceParser{id: 115, commit: 66, name: "list-item"}
	var p109 = sequenceParser{id: 109, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{115, 148, 149}}
	var p108 = sequenceParser{id: 108, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p105 = charParser{id: 105, chars: []rune{46}}
	var p106 = charParser{id: 106, chars: []rune{46}}
	var p107 = charParser{id: 107, chars: []rune{46}}
	p108.items = []parser{&p105, &p106, &p107}
	p109.items = []parser{&p274, &p827, &p108}
	p115.options = []parser{&p403, &p109}
	var p118 = sequenceParser{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p116 = sequenceParser{id: 116, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p116.items = []parser{&p114, &p827, &p115}
	var p117 = sequenceParser{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p117.items = []parser{&p827, &p116}
	p118.items = []parser{&p827, &p116, &p117}
	p119.items = []parser{&p115, &p118}
	var p123 = sequenceParser{id: 123, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p122 = charParser{id: 122, chars: []rune{93}}
	p123.items = []parser{&p122}
	p124.items = []parser{&p121, &p827, &p114, &p827, &p119, &p827, &p114, &p827, &p123}
	p125.items = []parser{&p124}
	var p130 = sequenceParser{id: 130, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p127 = sequenceParser{id: 127, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p126 = charParser{id: 126, chars: []rune{126}}
	p127.items = []parser{&p126}
	var p129 = sequenceParser{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p128 = sequenceParser{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p128.items = []parser{&p827, &p14}
	p129.items = []parser{&p827, &p14, &p128}
	p130.items = []parser{&p127, &p129, &p827, &p124}
	var p159 = sequenceParser{id: 159, commit: 64, name: "struct", ranges: [][]int{{1, 1}}, generalizations: []int{274, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
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
	p135.items = []parser{&p827, &p14}
	p136.items = []parser{&p827, &p14, &p135}
	var p138 = sequenceParser{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p137 = sequenceParser{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p137.items = []parser{&p827, &p14}
	p138.items = []parser{&p827, &p14, &p137}
	var p134 = sequenceParser{id: 134, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p133 = charParser{id: 133, chars: []rune{93}}
	p134.items = []parser{&p133}
	p139.items = []parser{&p132, &p136, &p827, &p403, &p138, &p827, &p134}
	p140.options = []parser{&p104, &p87, &p139}
	var p144 = sequenceParser{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p143 = sequenceParser{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p143.items = []parser{&p827, &p14}
	p144.items = []parser{&p827, &p14, &p143}
	var p142 = sequenceParser{id: 142, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p141 = charParser{id: 141, chars: []rune{58}}
	p142.items = []parser{&p141}
	var p146 = sequenceParser{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p145 = sequenceParser{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p145.items = []parser{&p827, &p14}
	p146.items = []parser{&p827, &p14, &p145}
	p147.items = []parser{&p140, &p144, &p827, &p142, &p146, &p827, &p403}
	p148.options = []parser{&p147, &p109}
	var p152 = sequenceParser{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p150 = sequenceParser{id: 150, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p149 = choiceParser{id: 149, commit: 2}
	p149.options = []parser{&p147, &p109}
	p150.items = []parser{&p114, &p827, &p149}
	var p151 = sequenceParser{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p151.items = []parser{&p827, &p150}
	p152.items = []parser{&p827, &p150, &p151}
	p153.items = []parser{&p148, &p152}
	var p157 = sequenceParser{id: 157, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p156 = charParser{id: 156, chars: []rune{125}}
	p157.items = []parser{&p156}
	p158.items = []parser{&p155, &p827, &p114, &p827, &p153, &p827, &p114, &p827, &p157}
	p159.items = []parser{&p158}
	var p164 = sequenceParser{id: 164, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 785, 200, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p161 = sequenceParser{id: 161, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p160 = charParser{id: 160, chars: []rune{126}}
	p161.items = []parser{&p160}
	var p163 = sequenceParser{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p162 = sequenceParser{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p162.items = []parser{&p827, &p14}
	p163.items = []parser{&p827, &p14, &p162}
	p164.items = []parser{&p161, &p163, &p827, &p158}
	var p209 = sequenceParser{id: 209, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785, 200, 274, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p206 = sequenceParser{id: 206, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p204 = charParser{id: 204, chars: []rune{102}}
	var p205 = charParser{id: 205, chars: []rune{110}}
	p206.items = []parser{&p204, &p205}
	var p208 = sequenceParser{id: 208, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p207 = sequenceParser{id: 207, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p207.items = []parser{&p827, &p14}
	p208.items = []parser{&p827, &p14, &p207}
	var p203 = sequenceParser{id: 203, commit: 66, name: "function-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p195 = sequenceParser{id: 195, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p194 = charParser{id: 194, chars: []rune{40}}
	p195.items = []parser{&p194}
	var p197 = choiceParser{id: 197, commit: 2}
	var p168 = sequenceParser{id: 168, commit: 66, name: "parameter-list", ranges: [][]int{{1, 1}, {0, 1}}, generalizations: []int{197}}
	var p167 = sequenceParser{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p165 = sequenceParser{id: 165, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p165.items = []parser{&p114, &p827, &p104}
	var p166 = sequenceParser{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p166.items = []parser{&p827, &p165}
	p167.items = []parser{&p827, &p165, &p166}
	p168.items = []parser{&p104, &p167}
	var p196 = sequenceParser{id: 196, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}, generalizations: []int{197}}
	var p175 = sequenceParser{id: 175, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{197}}
	var p172 = sequenceParser{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p169 = charParser{id: 169, chars: []rune{46}}
	var p170 = charParser{id: 170, chars: []rune{46}}
	var p171 = charParser{id: 171, chars: []rune{46}}
	p172.items = []parser{&p169, &p170, &p171}
	var p174 = sequenceParser{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p173 = sequenceParser{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p173.items = []parser{&p827, &p14}
	p174.items = []parser{&p827, &p14, &p173}
	p175.items = []parser{&p172, &p174, &p827, &p104}
	p196.items = []parser{&p168, &p827, &p114, &p827, &p175}
	p197.options = []parser{&p168, &p196, &p175}
	var p199 = sequenceParser{id: 199, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p198 = charParser{id: 198, chars: []rune{41}}
	p199.items = []parser{&p198}
	var p202 = sequenceParser{id: 202, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p201 = sequenceParser{id: 201, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p201.items = []parser{&p827, &p14}
	p202.items = []parser{&p827, &p14, &p201}
	var p200 = choiceParser{id: 200, commit: 2}
	var p785 = choiceParser{id: 785, commit: 66, name: "simple-statement", generalizations: []int{200, 795}}
	var p513 = sequenceParser{id: 513, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785, 200, 520, 795}}
	var p499 = sequenceParser{id: 499, commit: 74, name: "send-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p498 = sequenceParser{id: 498, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p494 = charParser{id: 494, chars: []rune{115}}
	var p495 = charParser{id: 495, chars: []rune{101}}
	var p496 = charParser{id: 496, chars: []rune{110}}
	var p497 = charParser{id: 497, chars: []rune{100}}
	p498.items = []parser{&p494, &p495, &p496, &p497}
	p499.items = []parser{&p498, &p15}
	var p510 = sequenceParser{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p509 = sequenceParser{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p509.items = []parser{&p827, &p14}
	p510.items = []parser{&p827, &p14, &p509}
	var p512 = sequenceParser{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p511 = sequenceParser{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p511.items = []parser{&p827, &p14}
	p512.items = []parser{&p827, &p14, &p511}
	p513.items = []parser{&p499, &p510, &p827, &p274, &p512, &p827, &p274}
	var p566 = sequenceParser{id: 566, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785, 200, 795}}
	var p556 = sequenceParser{id: 556, commit: 74, name: "go-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p555 = sequenceParser{id: 555, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p553 = charParser{id: 553, chars: []rune{103}}
	var p554 = charParser{id: 554, chars: []rune{111}}
	p555.items = []parser{&p553, &p554}
	p556.items = []parser{&p555, &p15}
	var p565 = sequenceParser{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p564 = sequenceParser{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p564.items = []parser{&p827, &p14}
	p565.items = []parser{&p827, &p14, &p564}
	var p264 = sequenceParser{id: 264, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p261 = sequenceParser{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p260 = charParser{id: 260, chars: []rune{40}}
	p261.items = []parser{&p260}
	var p263 = sequenceParser{id: 263, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p262 = charParser{id: 262, chars: []rune{41}}
	p263.items = []parser{&p262}
	p264.items = []parser{&p274, &p827, &p261, &p827, &p114, &p827, &p119, &p827, &p114, &p827, &p263}
	p566.items = []parser{&p556, &p565, &p827, &p264}
	var p575 = sequenceParser{id: 575, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785, 200, 795}}
	var p572 = sequenceParser{id: 572, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p567 = charParser{id: 567, chars: []rune{100}}
	var p568 = charParser{id: 568, chars: []rune{101}}
	var p569 = charParser{id: 569, chars: []rune{102}}
	var p570 = charParser{id: 570, chars: []rune{101}}
	var p571 = charParser{id: 571, chars: []rune{114}}
	p572.items = []parser{&p567, &p568, &p569, &p570, &p571}
	var p574 = sequenceParser{id: 574, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p573 = sequenceParser{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p573.items = []parser{&p827, &p14}
	p574.items = []parser{&p827, &p14, &p573}
	p575.items = []parser{&p572, &p574, &p827, &p264}
	var p639 = choiceParser{id: 639, commit: 64, name: "assignment", generalizations: []int{785, 200, 795}}
	var p623 = sequenceParser{id: 623, commit: 66, name: "assign-set", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{639, 785, 200, 795}}
	var p608 = sequenceParser{id: 608, commit: 74, name: "set-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p607 = sequenceParser{id: 607, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p604 = charParser{id: 604, chars: []rune{115}}
	var p605 = charParser{id: 605, chars: []rune{101}}
	var p606 = charParser{id: 606, chars: []rune{116}}
	p607.items = []parser{&p604, &p605, &p606}
	p608.items = []parser{&p607, &p15}
	var p622 = sequenceParser{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p621 = sequenceParser{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p621.items = []parser{&p827, &p14}
	p622.items = []parser{&p827, &p14, &p621}
	var p616 = sequenceParser{id: 616, commit: 66, name: "assign-capture", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p613 = sequenceParser{id: 613, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p612 = sequenceParser{id: 612, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p611 = sequenceParser{id: 611, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p611.items = []parser{&p827, &p14}
	p612.items = []parser{&p14, &p611}
	var p610 = sequenceParser{id: 610, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p609 = charParser{id: 609, chars: []rune{61}}
	p610.items = []parser{&p609}
	p613.items = []parser{&p612, &p827, &p610}
	var p615 = sequenceParser{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p614 = sequenceParser{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p614.items = []parser{&p827, &p14}
	p615.items = []parser{&p827, &p14, &p614}
	p616.items = []parser{&p274, &p827, &p613, &p615, &p827, &p403}
	p623.items = []parser{&p608, &p622, &p827, &p616}
	var p630 = sequenceParser{id: 630, commit: 66, name: "assign-eq", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{639, 785, 200, 795}}
	var p627 = sequenceParser{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p626 = sequenceParser{id: 626, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p626.items = []parser{&p827, &p14}
	p627.items = []parser{&p827, &p14, &p626}
	var p625 = sequenceParser{id: 625, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p624 = charParser{id: 624, chars: []rune{61}}
	p625.items = []parser{&p624}
	var p629 = sequenceParser{id: 629, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p628 = sequenceParser{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p628.items = []parser{&p827, &p14}
	p629.items = []parser{&p827, &p14, &p628}
	p630.items = []parser{&p274, &p627, &p827, &p625, &p629, &p827, &p403}
	var p638 = sequenceParser{id: 638, commit: 66, name: "assign-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{639, 785, 200, 795}}
	var p637 = sequenceParser{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p636 = sequenceParser{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p636.items = []parser{&p827, &p14}
	p637.items = []parser{&p827, &p14, &p636}
	var p632 = sequenceParser{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p631 = charParser{id: 631, chars: []rune{40}}
	p632.items = []parser{&p631}
	var p633 = sequenceParser{id: 633, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p620 = sequenceParser{id: 620, commit: 66, name: "assign-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p619 = sequenceParser{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p617 = sequenceParser{id: 617, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p617.items = []parser{&p114, &p827, &p616}
	var p618 = sequenceParser{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p618.items = []parser{&p827, &p617}
	p619.items = []parser{&p827, &p617, &p618}
	p620.items = []parser{&p616, &p619}
	p633.items = []parser{&p114, &p827, &p620}
	var p635 = sequenceParser{id: 635, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p634 = charParser{id: 634, chars: []rune{41}}
	p635.items = []parser{&p634}
	p638.items = []parser{&p608, &p637, &p827, &p632, &p827, &p633, &p827, &p114, &p827, &p635}
	p639.options = []parser{&p623, &p630, &p638}
	var p794 = sequenceParser{id: 794, commit: 66, name: "simple-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{785, 200, 795}}
	var p787 = sequenceParser{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p786 = charParser{id: 786, chars: []rune{40}}
	p787.items = []parser{&p786}
	var p791 = sequenceParser{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p790 = sequenceParser{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p790.items = []parser{&p827, &p14}
	p791.items = []parser{&p827, &p14, &p790}
	var p793 = sequenceParser{id: 793, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p792 = sequenceParser{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p792.items = []parser{&p827, &p14}
	p793.items = []parser{&p827, &p14, &p792}
	var p789 = sequenceParser{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p788 = charParser{id: 788, chars: []rune{41}}
	p789.items = []parser{&p788}
	p794.items = []parser{&p787, &p791, &p827, &p785, &p793, &p827, &p789}
	p785.options = []parser{&p513, &p566, &p575, &p639, &p794, &p403}
	var p193 = sequenceParser{id: 193, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{200}}
	var p190 = sequenceParser{id: 190, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p189 = charParser{id: 189, chars: []rune{123}}
	p190.items = []parser{&p189}
	var p192 = sequenceParser{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p191 = charParser{id: 191, chars: []rune{125}}
	p192.items = []parser{&p191}
	p193.items = []parser{&p190, &p827, &p809, &p827, &p813, &p827, &p809, &p827, &p192}
	p200.options = []parser{&p785, &p193}
	p203.items = []parser{&p195, &p827, &p114, &p827, &p197, &p827, &p114, &p827, &p199, &p202, &p827, &p200}
	p209.items = []parser{&p206, &p208, &p827, &p203}
	var p219 = sequenceParser{id: 219, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p212 = sequenceParser{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p210 = charParser{id: 210, chars: []rune{102}}
	var p211 = charParser{id: 211, chars: []rune{110}}
	p212.items = []parser{&p210, &p211}
	var p216 = sequenceParser{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p215 = sequenceParser{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p215.items = []parser{&p827, &p14}
	p216.items = []parser{&p827, &p14, &p215}
	var p214 = sequenceParser{id: 214, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p213 = charParser{id: 213, chars: []rune{126}}
	p214.items = []parser{&p213}
	var p218 = sequenceParser{id: 218, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p217 = sequenceParser{id: 217, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p217.items = []parser{&p827, &p14}
	p218.items = []parser{&p827, &p14, &p217}
	p219.items = []parser{&p212, &p216, &p827, &p214, &p218, &p827, &p203}
	var p259 = sequenceParser{id: 259, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p258 = sequenceParser{id: 258, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p257 = sequenceParser{id: 257, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p257.items = []parser{&p827, &p14}
	p258.items = []parser{&p827, &p14, &p257}
	var p256 = sequenceParser{id: 256, commit: 66, name: "index-list", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var p252 = choiceParser{id: 252, commit: 66, name: "index"}
	var p233 = sequenceParser{id: 233, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{252}}
	var p230 = sequenceParser{id: 230, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p229 = charParser{id: 229, chars: []rune{46}}
	p230.items = []parser{&p229}
	var p232 = sequenceParser{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p231 = sequenceParser{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p231.items = []parser{&p827, &p14}
	p232.items = []parser{&p827, &p14, &p231}
	p233.items = []parser{&p230, &p232, &p827, &p104}
	var p242 = sequenceParser{id: 242, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{252}}
	var p235 = sequenceParser{id: 235, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p234 = charParser{id: 234, chars: []rune{91}}
	p235.items = []parser{&p234}
	var p239 = sequenceParser{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p238 = sequenceParser{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p238.items = []parser{&p827, &p14}
	p239.items = []parser{&p827, &p14, &p238}
	var p241 = sequenceParser{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p240 = sequenceParser{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p240.items = []parser{&p827, &p14}
	p241.items = []parser{&p827, &p14, &p240}
	var p237 = sequenceParser{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p236 = charParser{id: 236, chars: []rune{93}}
	p237.items = []parser{&p236}
	p242.items = []parser{&p235, &p239, &p827, &p403, &p241, &p827, &p237}
	var p251 = sequenceParser{id: 251, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{252}}
	var p244 = sequenceParser{id: 244, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p243 = charParser{id: 243, chars: []rune{91}}
	p244.items = []parser{&p243}
	var p248 = sequenceParser{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p247 = sequenceParser{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p247.items = []parser{&p827, &p14}
	p248.items = []parser{&p827, &p14, &p247}
	var p228 = sequenceParser{id: 228, commit: 66, name: "range", ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{585, 591, 592}}
	var p220 = sequenceParser{id: 220, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	p220.items = []parser{&p403}
	var p225 = sequenceParser{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p224 = sequenceParser{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p224.items = []parser{&p827, &p14}
	p225.items = []parser{&p827, &p14, &p224}
	var p223 = sequenceParser{id: 223, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p222 = charParser{id: 222, chars: []rune{58}}
	p223.items = []parser{&p222}
	var p227 = sequenceParser{id: 227, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p226 = sequenceParser{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p226.items = []parser{&p827, &p14}
	p227.items = []parser{&p827, &p14, &p226}
	var p221 = sequenceParser{id: 221, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	p221.items = []parser{&p403}
	p228.items = []parser{&p220, &p225, &p827, &p223, &p227, &p827, &p221}
	var p250 = sequenceParser{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p249 = sequenceParser{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p249.items = []parser{&p827, &p14}
	p250.items = []parser{&p827, &p14, &p249}
	var p246 = sequenceParser{id: 246, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p245 = charParser{id: 245, chars: []rune{93}}
	p246.items = []parser{&p245}
	p251.items = []parser{&p244, &p248, &p827, &p228, &p250, &p827, &p246}
	p252.options = []parser{&p233, &p242, &p251}
	var p255 = sequenceParser{id: 255, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p254 = sequenceParser{id: 254, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p253 = sequenceParser{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p253.items = []parser{&p827, &p14}
	p254.items = []parser{&p14, &p253}
	p255.items = []parser{&p254, &p827, &p252}
	p256.items = []parser{&p252, &p827, &p255}
	p259.items = []parser{&p274, &p258, &p827, &p256}
	var p273 = sequenceParser{id: 273, commit: 66, name: "expression-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{274, 403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p266 = sequenceParser{id: 266, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p265 = charParser{id: 265, chars: []rune{40}}
	p266.items = []parser{&p265}
	var p270 = sequenceParser{id: 270, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p269 = sequenceParser{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p269.items = []parser{&p827, &p14}
	p270.items = []parser{&p827, &p14, &p269}
	var p272 = sequenceParser{id: 272, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p271 = sequenceParser{id: 271, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p271.items = []parser{&p827, &p14}
	p272.items = []parser{&p827, &p14, &p271}
	var p268 = sequenceParser{id: 268, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p267 = charParser{id: 267, chars: []rune{41}}
	p268.items = []parser{&p267}
	p273.items = []parser{&p266, &p270, &p827, &p403, &p272, &p827, &p268}
	p274.options = []parser{&p61, &p74, &p87, &p99, &p516, &p104, &p125, &p130, &p159, &p164, &p209, &p219, &p259, &p264, &p273}
	var p334 = sequenceParser{id: 334, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{403, 340, 341, 342, 343, 344, 395, 592, 585, 795}}
	var p333 = choiceParser{id: 333, commit: 66, name: "unary-operator"}
	var p293 = sequenceParser{id: 293, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p292 = charParser{id: 292, chars: []rune{43}}
	p293.items = []parser{&p292}
	var p295 = sequenceParser{id: 295, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p294 = charParser{id: 294, chars: []rune{45}}
	p295.items = []parser{&p294}
	var p276 = sequenceParser{id: 276, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p275 = charParser{id: 275, chars: []rune{94}}
	p276.items = []parser{&p275}
	var p307 = sequenceParser{id: 307, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{333}}
	var p306 = charParser{id: 306, chars: []rune{33}}
	p307.items = []parser{&p306}
	p333.options = []parser{&p293, &p295, &p276, &p307}
	p334.items = []parser{&p333, &p827, &p274}
	var p381 = choiceParser{id: 381, commit: 66, name: "binary-expression", generalizations: []int{403, 395, 592, 585, 795}}
	var p352 = sequenceParser{id: 352, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{381, 341, 342, 343, 344, 403, 395, 592, 585, 795}}
	var p340 = choiceParser{id: 340, commit: 66, name: "operand0", generalizations: []int{341, 342, 343, 344}}
	p340.options = []parser{&p274, &p334}
	var p350 = sequenceParser{id: 350, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p347 = sequenceParser{id: 347, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p346 = sequenceParser{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p346.items = []parser{&p827, &p14}
	p347.items = []parser{&p14, &p346}
	var p335 = choiceParser{id: 335, commit: 66, name: "binary-op0"}
	var p278 = sequenceParser{id: 278, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p277 = charParser{id: 277, chars: []rune{38}}
	p278.items = []parser{&p277}
	var p285 = sequenceParser{id: 285, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p283 = charParser{id: 283, chars: []rune{38}}
	var p284 = charParser{id: 284, chars: []rune{94}}
	p285.items = []parser{&p283, &p284}
	var p288 = sequenceParser{id: 288, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p286 = charParser{id: 286, chars: []rune{60}}
	var p287 = charParser{id: 287, chars: []rune{60}}
	p288.items = []parser{&p286, &p287}
	var p291 = sequenceParser{id: 291, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{335}}
	var p289 = charParser{id: 289, chars: []rune{62}}
	var p290 = charParser{id: 290, chars: []rune{62}}
	p291.items = []parser{&p289, &p290}
	var p297 = sequenceParser{id: 297, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p296 = charParser{id: 296, chars: []rune{42}}
	p297.items = []parser{&p296}
	var p299 = sequenceParser{id: 299, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p298 = charParser{id: 298, chars: []rune{47}}
	p299.items = []parser{&p298}
	var p301 = sequenceParser{id: 301, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{335}}
	var p300 = charParser{id: 300, chars: []rune{37}}
	p301.items = []parser{&p300}
	p335.options = []parser{&p278, &p285, &p288, &p291, &p297, &p299, &p301}
	var p349 = sequenceParser{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p348 = sequenceParser{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p348.items = []parser{&p827, &p14}
	p349.items = []parser{&p827, &p14, &p348}
	p350.items = []parser{&p347, &p827, &p335, &p349, &p827, &p340}
	var p351 = sequenceParser{id: 351, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p351.items = []parser{&p827, &p350}
	p352.items = []parser{&p340, &p827, &p350, &p351}
	var p359 = sequenceParser{id: 359, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{381, 342, 343, 344, 403, 395, 592, 585, 795}}
	var p341 = choiceParser{id: 341, commit: 66, name: "operand1", generalizations: []int{342, 343, 344}}
	p341.options = []parser{&p340, &p352}
	var p357 = sequenceParser{id: 357, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p354 = sequenceParser{id: 354, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p353 = sequenceParser{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p353.items = []parser{&p827, &p14}
	p354.items = []parser{&p14, &p353}
	var p336 = choiceParser{id: 336, commit: 66, name: "binary-op1"}
	var p280 = sequenceParser{id: 280, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{336}}
	var p279 = charParser{id: 279, chars: []rune{124}}
	p280.items = []parser{&p279}
	var p282 = sequenceParser{id: 282, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{336}}
	var p281 = charParser{id: 281, chars: []rune{94}}
	p282.items = []parser{&p281}
	var p303 = sequenceParser{id: 303, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{336}}
	var p302 = charParser{id: 302, chars: []rune{43}}
	p303.items = []parser{&p302}
	var p305 = sequenceParser{id: 305, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{336}}
	var p304 = charParser{id: 304, chars: []rune{45}}
	p305.items = []parser{&p304}
	p336.options = []parser{&p280, &p282, &p303, &p305}
	var p356 = sequenceParser{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p355 = sequenceParser{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p355.items = []parser{&p827, &p14}
	p356.items = []parser{&p827, &p14, &p355}
	p357.items = []parser{&p354, &p827, &p336, &p356, &p827, &p341}
	var p358 = sequenceParser{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p358.items = []parser{&p827, &p357}
	p359.items = []parser{&p341, &p827, &p357, &p358}
	var p366 = sequenceParser{id: 366, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{381, 343, 344, 403, 395, 592, 585, 795}}
	var p342 = choiceParser{id: 342, commit: 66, name: "operand2", generalizations: []int{343, 344}}
	p342.options = []parser{&p341, &p359}
	var p364 = sequenceParser{id: 364, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p361 = sequenceParser{id: 361, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p360 = sequenceParser{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p360.items = []parser{&p827, &p14}
	p361.items = []parser{&p14, &p360}
	var p337 = choiceParser{id: 337, commit: 66, name: "binary-op2"}
	var p310 = sequenceParser{id: 310, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{337}}
	var p308 = charParser{id: 308, chars: []rune{61}}
	var p309 = charParser{id: 309, chars: []rune{61}}
	p310.items = []parser{&p308, &p309}
	var p313 = sequenceParser{id: 313, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{337}}
	var p311 = charParser{id: 311, chars: []rune{33}}
	var p312 = charParser{id: 312, chars: []rune{61}}
	p313.items = []parser{&p311, &p312}
	var p315 = sequenceParser{id: 315, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{337}}
	var p314 = charParser{id: 314, chars: []rune{60}}
	p315.items = []parser{&p314}
	var p318 = sequenceParser{id: 318, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{337}}
	var p316 = charParser{id: 316, chars: []rune{60}}
	var p317 = charParser{id: 317, chars: []rune{61}}
	p318.items = []parser{&p316, &p317}
	var p320 = sequenceParser{id: 320, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{337}}
	var p319 = charParser{id: 319, chars: []rune{62}}
	p320.items = []parser{&p319}
	var p323 = sequenceParser{id: 323, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}, generalizations: []int{337}}
	var p321 = charParser{id: 321, chars: []rune{62}}
	var p322 = charParser{id: 322, chars: []rune{61}}
	p323.items = []parser{&p321, &p322}
	p337.options = []parser{&p310, &p313, &p315, &p318, &p320, &p323}
	var p363 = sequenceParser{id: 363, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p362 = sequenceParser{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p362.items = []parser{&p827, &p14}
	p363.items = []parser{&p827, &p14, &p362}
	p364.items = []parser{&p361, &p827, &p337, &p363, &p827, &p342}
	var p365 = sequenceParser{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p365.items = []parser{&p827, &p364}
	p366.items = []parser{&p342, &p827, &p364, &p365}
	var p373 = sequenceParser{id: 373, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{381, 344, 403, 395, 592, 585, 795}}
	var p343 = choiceParser{id: 343, commit: 66, name: "operand3", generalizations: []int{344}}
	p343.options = []parser{&p342, &p366}
	var p371 = sequenceParser{id: 371, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p368 = sequenceParser{id: 368, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p367 = sequenceParser{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p367.items = []parser{&p827, &p14}
	p368.items = []parser{&p14, &p367}
	var p338 = sequenceParser{id: 338, commit: 66, name: "binary-op3", ranges: [][]int{{1, 1}}}
	var p326 = sequenceParser{id: 326, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p324 = charParser{id: 324, chars: []rune{38}}
	var p325 = charParser{id: 325, chars: []rune{38}}
	p326.items = []parser{&p324, &p325}
	p338.items = []parser{&p326}
	var p370 = sequenceParser{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p369 = sequenceParser{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p369.items = []parser{&p827, &p14}
	p370.items = []parser{&p827, &p14, &p369}
	p371.items = []parser{&p368, &p827, &p338, &p370, &p827, &p343}
	var p372 = sequenceParser{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p372.items = []parser{&p827, &p371}
	p373.items = []parser{&p343, &p827, &p371, &p372}
	var p380 = sequenceParser{id: 380, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{381, 403, 395, 592, 585, 795}}
	var p344 = choiceParser{id: 344, commit: 66, name: "operand4"}
	p344.options = []parser{&p343, &p373}
	var p378 = sequenceParser{id: 378, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p375 = sequenceParser{id: 375, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p374 = sequenceParser{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p374.items = []parser{&p827, &p14}
	p375.items = []parser{&p14, &p374}
	var p339 = sequenceParser{id: 339, commit: 66, name: "binary-op4", ranges: [][]int{{1, 1}}}
	var p329 = sequenceParser{id: 329, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p327 = charParser{id: 327, chars: []rune{124}}
	var p328 = charParser{id: 328, chars: []rune{124}}
	p329.items = []parser{&p327, &p328}
	p339.items = []parser{&p329}
	var p377 = sequenceParser{id: 377, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p376 = sequenceParser{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p376.items = []parser{&p827, &p14}
	p377.items = []parser{&p827, &p14, &p376}
	p378.items = []parser{&p375, &p827, &p339, &p377, &p827, &p344}
	var p379 = sequenceParser{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p379.items = []parser{&p827, &p378}
	p380.items = []parser{&p344, &p827, &p378, &p379}
	p381.options = []parser{&p352, &p359, &p366, &p373, &p380}
	var p394 = sequenceParser{id: 394, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{403, 395, 592, 585, 795}}
	var p387 = sequenceParser{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p386 = sequenceParser{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p386.items = []parser{&p827, &p14}
	p387.items = []parser{&p827, &p14, &p386}
	var p383 = sequenceParser{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p382 = charParser{id: 382, chars: []rune{63}}
	p383.items = []parser{&p382}
	var p389 = sequenceParser{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p388 = sequenceParser{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p388.items = []parser{&p827, &p14}
	p389.items = []parser{&p827, &p14, &p388}
	var p391 = sequenceParser{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p390 = sequenceParser{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p390.items = []parser{&p827, &p14}
	p391.items = []parser{&p827, &p14, &p390}
	var p385 = sequenceParser{id: 385, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p384 = charParser{id: 384, chars: []rune{58}}
	p385.items = []parser{&p384}
	var p393 = sequenceParser{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p392 = sequenceParser{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p392.items = []parser{&p827, &p14}
	p393.items = []parser{&p827, &p14, &p392}
	p394.items = []parser{&p403, &p387, &p827, &p383, &p389, &p827, &p403, &p391, &p827, &p385, &p393, &p827, &p403}
	var p402 = sequenceParser{id: 402, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}, generalizations: []int{403, 592, 585, 795}}
	var p395 = choiceParser{id: 395, commit: 66, name: "chainingOperand"}
	p395.options = []parser{&p274, &p334, &p381, &p394}
	var p400 = sequenceParser{id: 400, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p397 = sequenceParser{id: 397, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p396 = sequenceParser{id: 396, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p396.items = []parser{&p827, &p14}
	p397.items = []parser{&p14, &p396}
	var p332 = sequenceParser{id: 332, commit: 74, name: "chain", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p330 = charParser{id: 330, chars: []rune{45}}
	var p331 = charParser{id: 331, chars: []rune{62}}
	p332.items = []parser{&p330, &p331}
	var p399 = sequenceParser{id: 399, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p398 = sequenceParser{id: 398, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p398.items = []parser{&p827, &p14}
	p399.items = []parser{&p827, &p14, &p398}
	p400.items = []parser{&p397, &p827, &p332, &p399, &p827, &p395}
	var p401 = sequenceParser{id: 401, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p401.items = []parser{&p827, &p400}
	p402.items = []parser{&p395, &p827, &p400, &p401}
	p403.options = []parser{&p274, &p334, &p381, &p394, &p402}
	p186.items = []parser{&p185, &p827, &p403}
	p187.items = []parser{&p183, &p827, &p186}
	var p188 = sequenceParser{id: 188, commit: 64, name: "return-statement", ranges: [][]int{{1, 1}}, generalizations: []int{795, 482, 543}}
	p188.items = []parser{&p183}
	var p434 = sequenceParser{id: 434, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{795, 482, 543}}
	var p407 = sequenceParser{id: 407, commit: 74, name: "if-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p406 = sequenceParser{id: 406, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p404 = charParser{id: 404, chars: []rune{105}}
	var p405 = charParser{id: 405, chars: []rune{102}}
	p406.items = []parser{&p404, &p405}
	p407.items = []parser{&p406, &p15}
	var p429 = sequenceParser{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p428 = sequenceParser{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p428.items = []parser{&p827, &p14}
	p429.items = []parser{&p827, &p14, &p428}
	var p431 = sequenceParser{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p430 = sequenceParser{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p430.items = []parser{&p827, &p14}
	p431.items = []parser{&p827, &p14, &p430}
	var p433 = sequenceParser{id: 433, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p422 = sequenceParser{id: 422, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p415 = sequenceParser{id: 415, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p414 = sequenceParser{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p414.items = []parser{&p827, &p14}
	p415.items = []parser{&p14, &p414}
	var p413 = sequenceParser{id: 413, commit: 74, name: "else-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p412 = sequenceParser{id: 412, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p408 = charParser{id: 408, chars: []rune{101}}
	var p409 = charParser{id: 409, chars: []rune{108}}
	var p410 = charParser{id: 410, chars: []rune{115}}
	var p411 = charParser{id: 411, chars: []rune{101}}
	p412.items = []parser{&p408, &p409, &p410, &p411}
	p413.items = []parser{&p412, &p15}
	var p417 = sequenceParser{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p416 = sequenceParser{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p416.items = []parser{&p827, &p14}
	p417.items = []parser{&p827, &p14, &p416}
	var p419 = sequenceParser{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p418 = sequenceParser{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p418.items = []parser{&p827, &p14}
	p419.items = []parser{&p827, &p14, &p418}
	var p421 = sequenceParser{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p420 = sequenceParser{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p420.items = []parser{&p827, &p14}
	p421.items = []parser{&p827, &p14, &p420}
	p422.items = []parser{&p415, &p827, &p413, &p417, &p827, &p407, &p419, &p827, &p403, &p421, &p827, &p193}
	var p432 = sequenceParser{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p432.items = []parser{&p827, &p422}
	p433.items = []parser{&p827, &p422, &p432}
	var p427 = sequenceParser{id: 427, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p424 = sequenceParser{id: 424, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p423 = sequenceParser{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p423.items = []parser{&p827, &p14}
	p424.items = []parser{&p14, &p423}
	var p426 = sequenceParser{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p425 = sequenceParser{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p425.items = []parser{&p827, &p14}
	p426.items = []parser{&p827, &p14, &p425}
	p427.items = []parser{&p424, &p827, &p413, &p426, &p827, &p193}
	p434.items = []parser{&p407, &p429, &p827, &p403, &p431, &p827, &p193, &p433, &p827, &p427}
	var p493 = sequenceParser{id: 493, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{482, 795, 543}}
	var p448 = sequenceParser{id: 448, commit: 74, name: "switch-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p447 = sequenceParser{id: 447, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p441 = charParser{id: 441, chars: []rune{115}}
	var p442 = charParser{id: 442, chars: []rune{119}}
	var p443 = charParser{id: 443, chars: []rune{105}}
	var p444 = charParser{id: 444, chars: []rune{116}}
	var p445 = charParser{id: 445, chars: []rune{99}}
	var p446 = charParser{id: 446, chars: []rune{104}}
	p447.items = []parser{&p441, &p442, &p443, &p444, &p445, &p446}
	p448.items = []parser{&p447, &p15}
	var p490 = sequenceParser{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p489 = sequenceParser{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p489.items = []parser{&p827, &p14}
	p490.items = []parser{&p827, &p14, &p489}
	var p492 = sequenceParser{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p491 = sequenceParser{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p491.items = []parser{&p827, &p14}
	p492.items = []parser{&p827, &p14, &p491}
	var p480 = sequenceParser{id: 480, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p479 = charParser{id: 479, chars: []rune{123}}
	p480.items = []parser{&p479}
	var p486 = sequenceParser{id: 486, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p481 = choiceParser{id: 481, commit: 2}
	var p478 = sequenceParser{id: 478, commit: 66, name: "case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{481, 482}}
	var p473 = sequenceParser{id: 473, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p440 = sequenceParser{id: 440, commit: 74, name: "case-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p439 = sequenceParser{id: 439, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p435 = charParser{id: 435, chars: []rune{99}}
	var p436 = charParser{id: 436, chars: []rune{97}}
	var p437 = charParser{id: 437, chars: []rune{115}}
	var p438 = charParser{id: 438, chars: []rune{101}}
	p439.items = []parser{&p435, &p436, &p437, &p438}
	p440.items = []parser{&p439, &p15}
	var p470 = sequenceParser{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p469 = sequenceParser{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p469.items = []parser{&p827, &p14}
	p470.items = []parser{&p827, &p14, &p469}
	var p472 = sequenceParser{id: 472, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p471 = sequenceParser{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p471.items = []parser{&p827, &p14}
	p472.items = []parser{&p827, &p14, &p471}
	var p468 = sequenceParser{id: 468, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p467 = charParser{id: 467, chars: []rune{58}}
	p468.items = []parser{&p467}
	p473.items = []parser{&p440, &p470, &p827, &p403, &p472, &p827, &p468}
	var p477 = sequenceParser{id: 477, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p475 = sequenceParser{id: 475, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p474 = charParser{id: 474, chars: []rune{59}}
	p475.items = []parser{&p474}
	var p476 = sequenceParser{id: 476, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p476.items = []parser{&p827, &p475}
	p477.items = []parser{&p827, &p475, &p476}
	p478.items = []parser{&p473, &p477, &p827, &p795}
	var p466 = sequenceParser{id: 466, commit: 66, name: "default-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{481, 482, 542, 543}}
	var p461 = sequenceParser{id: 461, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p456 = sequenceParser{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p449 = charParser{id: 449, chars: []rune{100}}
	var p450 = charParser{id: 450, chars: []rune{101}}
	var p451 = charParser{id: 451, chars: []rune{102}}
	var p452 = charParser{id: 452, chars: []rune{97}}
	var p453 = charParser{id: 453, chars: []rune{117}}
	var p454 = charParser{id: 454, chars: []rune{108}}
	var p455 = charParser{id: 455, chars: []rune{116}}
	p456.items = []parser{&p449, &p450, &p451, &p452, &p453, &p454, &p455}
	var p460 = sequenceParser{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p459 = sequenceParser{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p459.items = []parser{&p827, &p14}
	p460.items = []parser{&p827, &p14, &p459}
	var p458 = sequenceParser{id: 458, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p457 = charParser{id: 457, chars: []rune{58}}
	p458.items = []parser{&p457}
	p461.items = []parser{&p456, &p460, &p827, &p458}
	var p465 = sequenceParser{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p463 = sequenceParser{id: 463, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p462 = charParser{id: 462, chars: []rune{59}}
	p463.items = []parser{&p462}
	var p464 = sequenceParser{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p464.items = []parser{&p827, &p463}
	p465.items = []parser{&p827, &p463, &p464}
	p466.items = []parser{&p461, &p465, &p827, &p795}
	p481.options = []parser{&p478, &p466}
	var p485 = sequenceParser{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p483 = sequenceParser{id: 483, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p482 = choiceParser{id: 482, commit: 2}
	p482.options = []parser{&p478, &p466, &p795}
	p483.items = []parser{&p809, &p827, &p482}
	var p484 = sequenceParser{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p484.items = []parser{&p827, &p483}
	p485.items = []parser{&p827, &p483, &p484}
	p486.items = []parser{&p481, &p485}
	var p488 = sequenceParser{id: 488, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p487 = charParser{id: 487, chars: []rune{125}}
	p488.items = []parser{&p487}
	p493.items = []parser{&p448, &p490, &p827, &p403, &p492, &p827, &p480, &p827, &p809, &p827, &p486, &p827, &p809, &p827, &p488}
	var p552 = sequenceParser{id: 552, commit: 64, name: "select", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{543, 795}}
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
	p550.items = []parser{&p827, &p14}
	p551.items = []parser{&p827, &p14, &p550}
	var p541 = sequenceParser{id: 541, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p540 = charParser{id: 540, chars: []rune{123}}
	p541.items = []parser{&p540}
	var p547 = sequenceParser{id: 547, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var p542 = choiceParser{id: 542, commit: 2}
	var p532 = sequenceParser{id: 532, commit: 66, name: "select-case-line", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}, generalizations: []int{542, 543}}
	var p527 = sequenceParser{id: 527, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p524 = sequenceParser{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p523 = sequenceParser{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p523.items = []parser{&p827, &p14}
	p524.items = []parser{&p827, &p14, &p523}
	var p520 = choiceParser{id: 520, commit: 66, name: "communication"}
	var p519 = sequenceParser{id: 519, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{520}}
	var p518 = sequenceParser{id: 518, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p517 = sequenceParser{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p517.items = []parser{&p827, &p14}
	p518.items = []parser{&p827, &p14, &p517}
	p519.items = []parser{&p104, &p518, &p827, &p516}
	p520.options = []parser{&p513, &p516, &p519}
	var p526 = sequenceParser{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p525 = sequenceParser{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p525.items = []parser{&p827, &p14}
	p526.items = []parser{&p827, &p14, &p525}
	var p522 = sequenceParser{id: 522, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p521 = charParser{id: 521, chars: []rune{58}}
	p522.items = []parser{&p521}
	p527.items = []parser{&p440, &p524, &p827, &p520, &p526, &p827, &p522}
	var p531 = sequenceParser{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p529 = sequenceParser{id: 529, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p528 = charParser{id: 528, chars: []rune{59}}
	p529.items = []parser{&p528}
	var p530 = sequenceParser{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p530.items = []parser{&p827, &p529}
	p531.items = []parser{&p827, &p529, &p530}
	p532.items = []parser{&p527, &p531, &p827, &p795}
	p542.options = []parser{&p532, &p466}
	var p546 = sequenceParser{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p544 = sequenceParser{id: 544, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p543 = choiceParser{id: 543, commit: 2}
	p543.options = []parser{&p532, &p466, &p795}
	p544.items = []parser{&p809, &p827, &p543}
	var p545 = sequenceParser{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p545.items = []parser{&p827, &p544}
	p546.items = []parser{&p827, &p544, &p545}
	p547.items = []parser{&p542, &p546}
	var p549 = sequenceParser{id: 549, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p548 = charParser{id: 548, chars: []rune{125}}
	p549.items = []parser{&p548}
	p552.items = []parser{&p539, &p551, &p827, &p541, &p827, &p809, &p827, &p547, &p827, &p809, &p827, &p549}
	var p603 = sequenceParser{id: 603, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}, generalizations: []int{795}}
	var p584 = sequenceParser{id: 584, commit: 74, name: "for-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p583 = sequenceParser{id: 583, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p580 = charParser{id: 580, chars: []rune{102}}
	var p581 = charParser{id: 581, chars: []rune{111}}
	var p582 = charParser{id: 582, chars: []rune{114}}
	p583.items = []parser{&p580, &p581, &p582}
	p584.items = []parser{&p583, &p15}
	var p602 = choiceParser{id: 602, commit: 2}
	var p598 = sequenceParser{id: 598, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{602}}
	var p595 = sequenceParser{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p594 = sequenceParser{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p593 = sequenceParser{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p593.items = []parser{&p827, &p14}
	p594.items = []parser{&p14, &p593}
	var p592 = choiceParser{id: 592, commit: 66, name: "loop-expression"}
	var p591 = choiceParser{id: 591, commit: 64, name: "range-over-expression", generalizations: []int{592}}
	var p590 = sequenceParser{id: 590, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{591, 592}}
	var p587 = sequenceParser{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p586 = sequenceParser{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p586.items = []parser{&p827, &p14}
	p587.items = []parser{&p827, &p14, &p586}
	var p579 = sequenceParser{id: 579, commit: 74, name: "in-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p578 = sequenceParser{id: 578, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p576 = charParser{id: 576, chars: []rune{105}}
	var p577 = charParser{id: 577, chars: []rune{110}}
	p578.items = []parser{&p576, &p577}
	p579.items = []parser{&p578, &p15}
	var p589 = sequenceParser{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p588 = sequenceParser{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p588.items = []parser{&p827, &p14}
	p589.items = []parser{&p827, &p14, &p588}
	var p585 = choiceParser{id: 585, commit: 2}
	p585.options = []parser{&p403, &p228}
	p590.items = []parser{&p104, &p587, &p827, &p579, &p589, &p827, &p585}
	p591.options = []parser{&p590, &p228}
	p592.options = []parser{&p403, &p591}
	p595.items = []parser{&p594, &p827, &p592}
	var p597 = sequenceParser{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p596 = sequenceParser{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p596.items = []parser{&p827, &p14}
	p597.items = []parser{&p827, &p14, &p596}
	p598.items = []parser{&p595, &p597, &p827, &p193}
	var p601 = sequenceParser{id: 601, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}, generalizations: []int{602}}
	var p600 = sequenceParser{id: 600, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p599 = sequenceParser{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p599.items = []parser{&p827, &p14}
	p600.items = []parser{&p14, &p599}
	p601.items = []parser{&p600, &p827, &p193}
	p602.options = []parser{&p598, &p601}
	p603.items = []parser{&p584, &p827, &p602}
	var p741 = choiceParser{id: 741, commit: 66, name: "definition", generalizations: []int{795}}
	var p662 = sequenceParser{id: 662, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 795}}
	var p644 = sequenceParser{id: 644, commit: 74, name: "let-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p643 = sequenceParser{id: 643, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p640 = charParser{id: 640, chars: []rune{108}}
	var p641 = charParser{id: 641, chars: []rune{101}}
	var p642 = charParser{id: 642, chars: []rune{116}}
	p643.items = []parser{&p640, &p641, &p642}
	p644.items = []parser{&p643, &p15}
	var p661 = sequenceParser{id: 661, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p660 = sequenceParser{id: 660, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p660.items = []parser{&p827, &p14}
	p661.items = []parser{&p827, &p14, &p660}
	var p659 = choiceParser{id: 659, commit: 2}
	var p653 = sequenceParser{id: 653, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}, generalizations: []int{659, 663, 664}}
	var p652 = sequenceParser{id: 652, commit: 66, name: "value-capture-fact", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p649 = sequenceParser{id: 649, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p648 = sequenceParser{id: 648, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p647 = sequenceParser{id: 647, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p647.items = []parser{&p827, &p14}
	p648.items = []parser{&p14, &p647}
	var p646 = sequenceParser{id: 646, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p645 = charParser{id: 645, chars: []rune{61}}
	p646.items = []parser{&p645}
	p649.items = []parser{&p648, &p827, &p646}
	var p651 = sequenceParser{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p650 = sequenceParser{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p650.items = []parser{&p827, &p14}
	p651.items = []parser{&p827, &p14, &p650}
	p652.items = []parser{&p104, &p827, &p649, &p651, &p827, &p403}
	p653.items = []parser{&p652}
	var p658 = sequenceParser{id: 658, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{659, 663, 664}}
	var p655 = sequenceParser{id: 655, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p654 = charParser{id: 654, chars: []rune{126}}
	p655.items = []parser{&p654}
	var p657 = sequenceParser{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p656 = sequenceParser{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p656.items = []parser{&p827, &p14}
	p657.items = []parser{&p827, &p14, &p656}
	p658.items = []parser{&p655, &p657, &p827, &p652}
	p659.options = []parser{&p653, &p658}
	p662.items = []parser{&p644, &p661, &p827, &p659}
	var p679 = sequenceParser{id: 679, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 795}}
	var p678 = sequenceParser{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p677 = sequenceParser{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p677.items = []parser{&p827, &p14}
	p678.items = []parser{&p827, &p14, &p677}
	var p674 = sequenceParser{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p673 = charParser{id: 673, chars: []rune{40}}
	p674.items = []parser{&p673}
	var p668 = sequenceParser{id: 668, commit: 66, name: "mixed-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p663 = choiceParser{id: 663, commit: 2}
	p663.options = []parser{&p653, &p658}
	var p667 = sequenceParser{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p665 = sequenceParser{id: 665, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var p664 = choiceParser{id: 664, commit: 2}
	p664.options = []parser{&p653, &p658}
	p665.items = []parser{&p114, &p827, &p664}
	var p666 = sequenceParser{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p666.items = []parser{&p827, &p665}
	p667.items = []parser{&p827, &p665, &p666}
	p668.items = []parser{&p663, &p667}
	var p676 = sequenceParser{id: 676, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p675 = charParser{id: 675, chars: []rune{41}}
	p676.items = []parser{&p675}
	p679.items = []parser{&p644, &p678, &p827, &p674, &p827, &p114, &p827, &p668, &p827, &p114, &p827, &p676}
	var p690 = sequenceParser{id: 690, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 795}}
	var p687 = sequenceParser{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p686 = sequenceParser{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p686.items = []parser{&p827, &p14}
	p687.items = []parser{&p827, &p14, &p686}
	var p681 = sequenceParser{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p680 = charParser{id: 680, chars: []rune{126}}
	p681.items = []parser{&p680}
	var p689 = sequenceParser{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p688 = sequenceParser{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p688.items = []parser{&p827, &p14}
	p689.items = []parser{&p827, &p14, &p688}
	var p683 = sequenceParser{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p682 = charParser{id: 682, chars: []rune{40}}
	p683.items = []parser{&p682}
	var p672 = sequenceParser{id: 672, commit: 66, name: "value-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p671 = sequenceParser{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p669 = sequenceParser{id: 669, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p669.items = []parser{&p114, &p827, &p653}
	var p670 = sequenceParser{id: 670, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p670.items = []parser{&p827, &p669}
	p671.items = []parser{&p827, &p669, &p670}
	p672.items = []parser{&p653, &p671}
	var p685 = sequenceParser{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p684 = charParser{id: 684, chars: []rune{41}}
	p685.items = []parser{&p684}
	p690.items = []parser{&p644, &p687, &p827, &p681, &p689, &p827, &p683, &p827, &p114, &p827, &p672, &p827, &p114, &p827, &p685}
	var p706 = sequenceParser{id: 706, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 795}}
	var p702 = sequenceParser{id: 702, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p700 = charParser{id: 700, chars: []rune{102}}
	var p701 = charParser{id: 701, chars: []rune{110}}
	p702.items = []parser{&p700, &p701}
	var p705 = sequenceParser{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p704 = sequenceParser{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p704.items = []parser{&p827, &p14}
	p705.items = []parser{&p827, &p14, &p704}
	var p703 = choiceParser{id: 703, commit: 2}
	var p694 = sequenceParser{id: 694, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}, generalizations: []int{703, 711, 712}}
	var p693 = sequenceParser{id: 693, commit: 66, name: "function-definition-fact", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var p692 = sequenceParser{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p691 = sequenceParser{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p691.items = []parser{&p827, &p14}
	p692.items = []parser{&p827, &p14, &p691}
	p693.items = []parser{&p104, &p692, &p827, &p203}
	p694.items = []parser{&p693}
	var p699 = sequenceParser{id: 699, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{703, 711, 712}}
	var p696 = sequenceParser{id: 696, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p695 = charParser{id: 695, chars: []rune{126}}
	p696.items = []parser{&p695}
	var p698 = sequenceParser{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p697 = sequenceParser{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p697.items = []parser{&p827, &p14}
	p698.items = []parser{&p827, &p14, &p697}
	p699.items = []parser{&p696, &p698, &p827, &p693}
	p703.options = []parser{&p694, &p699}
	p706.items = []parser{&p702, &p705, &p827, &p703}
	var p726 = sequenceParser{id: 726, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 795}}
	var p719 = sequenceParser{id: 719, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p717 = charParser{id: 717, chars: []rune{102}}
	var p718 = charParser{id: 718, chars: []rune{110}}
	p719.items = []parser{&p717, &p718}
	var p725 = sequenceParser{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p724 = sequenceParser{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p724.items = []parser{&p827, &p14}
	p725.items = []parser{&p827, &p14, &p724}
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
	p713.items = []parser{&p114, &p827, &p712}
	var p714 = sequenceParser{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p714.items = []parser{&p827, &p713}
	p715.items = []parser{&p827, &p713, &p714}
	p716.items = []parser{&p711, &p715}
	var p723 = sequenceParser{id: 723, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p722 = charParser{id: 722, chars: []rune{41}}
	p723.items = []parser{&p722}
	p726.items = []parser{&p719, &p725, &p827, &p721, &p827, &p114, &p827, &p716, &p827, &p114, &p827, &p723}
	var p740 = sequenceParser{id: 740, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{741, 795}}
	var p729 = sequenceParser{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p727 = charParser{id: 727, chars: []rune{102}}
	var p728 = charParser{id: 728, chars: []rune{110}}
	p729.items = []parser{&p727, &p728}
	var p737 = sequenceParser{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p736 = sequenceParser{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p736.items = []parser{&p827, &p14}
	p737.items = []parser{&p827, &p14, &p736}
	var p731 = sequenceParser{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p730 = charParser{id: 730, chars: []rune{126}}
	p731.items = []parser{&p730}
	var p739 = sequenceParser{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p738 = sequenceParser{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p738.items = []parser{&p827, &p14}
	p739.items = []parser{&p827, &p14, &p738}
	var p733 = sequenceParser{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p732 = charParser{id: 732, chars: []rune{40}}
	p733.items = []parser{&p732}
	var p710 = sequenceParser{id: 710, commit: 66, name: "function-capture-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p709 = sequenceParser{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p707 = sequenceParser{id: 707, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p707.items = []parser{&p114, &p827, &p694}
	var p708 = sequenceParser{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p708.items = []parser{&p827, &p707}
	p709.items = []parser{&p827, &p707, &p708}
	p710.items = []parser{&p694, &p709}
	var p735 = sequenceParser{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p734 = charParser{id: 734, chars: []rune{41}}
	p735.items = []parser{&p734}
	p740.items = []parser{&p729, &p737, &p827, &p731, &p739, &p827, &p733, &p827, &p114, &p827, &p710, &p827, &p114, &p827, &p735}
	p741.options = []parser{&p662, &p679, &p690, &p706, &p726, &p740}
	var p773 = choiceParser{id: 773, commit: 64, name: "use", generalizations: []int{795}}
	var p765 = sequenceParser{id: 765, commit: 66, name: "use-statement", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{773, 795}}
	var p746 = sequenceParser{id: 746, commit: 74, name: "use-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p745 = sequenceParser{id: 745, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p742 = charParser{id: 742, chars: []rune{117}}
	var p743 = charParser{id: 743, chars: []rune{115}}
	var p744 = charParser{id: 744, chars: []rune{101}}
	p745.items = []parser{&p742, &p743, &p744}
	p746.items = []parser{&p745, &p15}
	var p764 = sequenceParser{id: 764, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p763 = sequenceParser{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p763.items = []parser{&p827, &p14}
	p764.items = []parser{&p827, &p14, &p763}
	var p758 = choiceParser{id: 758, commit: 64, name: "use-fact"}
	var p757 = sequenceParser{id: 757, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{758}}
	var p749 = choiceParser{id: 749, commit: 2}
	var p748 = sequenceParser{id: 748, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}, generalizations: []int{749}}
	var p747 = charParser{id: 747, chars: []rune{46}}
	p748.items = []parser{&p747}
	p749.options = []parser{&p104, &p748}
	var p754 = sequenceParser{id: 754, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var p753 = sequenceParser{id: 753, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var p752 = sequenceParser{id: 752, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p752.items = []parser{&p827, &p14}
	p753.items = []parser{&p14, &p752}
	var p751 = sequenceParser{id: 751, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p750 = charParser{id: 750, chars: []rune{61}}
	p751.items = []parser{&p750}
	p754.items = []parser{&p753, &p827, &p751}
	var p756 = sequenceParser{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p755 = sequenceParser{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p755.items = []parser{&p827, &p14}
	p756.items = []parser{&p827, &p14, &p755}
	p757.items = []parser{&p749, &p827, &p754, &p756, &p827, &p87}
	p758.options = []parser{&p87, &p757}
	p765.items = []parser{&p746, &p764, &p827, &p758}
	var p772 = sequenceParser{id: 772, commit: 66, name: "use-statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{773, 795}}
	var p771 = sequenceParser{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p770 = sequenceParser{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p770.items = []parser{&p827, &p14}
	p771.items = []parser{&p827, &p14, &p770}
	var p767 = sequenceParser{id: 767, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p766 = charParser{id: 766, chars: []rune{40}}
	p767.items = []parser{&p766}
	var p762 = sequenceParser{id: 762, commit: 66, name: "use-fact-list", ranges: [][]int{{1, 1}, {0, 1}}}
	var p761 = sequenceParser{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p759 = sequenceParser{id: 759, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p759.items = []parser{&p114, &p827, &p758}
	var p760 = sequenceParser{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p760.items = []parser{&p827, &p759}
	p761.items = []parser{&p827, &p759, &p760}
	p762.items = []parser{&p758, &p761}
	var p769 = sequenceParser{id: 769, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p768 = charParser{id: 768, chars: []rune{41}}
	p769.items = []parser{&p768}
	p772.items = []parser{&p746, &p771, &p827, &p767, &p827, &p114, &p827, &p762, &p827, &p114, &p827, &p769}
	p773.options = []parser{&p765, &p772}
	var p784 = sequenceParser{id: 784, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{795}}
	var p781 = sequenceParser{id: 781, commit: 74, name: "export-mod", ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p780 = sequenceParser{id: 780, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var p774 = charParser{id: 774, chars: []rune{101}}
	var p775 = charParser{id: 775, chars: []rune{120}}
	var p776 = charParser{id: 776, chars: []rune{112}}
	var p777 = charParser{id: 777, chars: []rune{111}}
	var p778 = charParser{id: 778, chars: []rune{114}}
	var p779 = charParser{id: 779, chars: []rune{116}}
	p780.items = []parser{&p774, &p775, &p776, &p777, &p778, &p779}
	p781.items = []parser{&p780, &p15}
	var p783 = sequenceParser{id: 783, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p782 = sequenceParser{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p782.items = []parser{&p827, &p14}
	p783.items = []parser{&p827, &p14, &p782}
	p784.items = []parser{&p781, &p783, &p827, &p741}
	var p804 = sequenceParser{id: 804, commit: 66, name: "statement-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}, generalizations: []int{795}}
	var p797 = sequenceParser{id: 797, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p796 = charParser{id: 796, chars: []rune{40}}
	p797.items = []parser{&p796}
	var p801 = sequenceParser{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p800 = sequenceParser{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p800.items = []parser{&p827, &p14}
	p801.items = []parser{&p827, &p14, &p800}
	var p803 = sequenceParser{id: 803, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p802 = sequenceParser{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p802.items = []parser{&p827, &p14}
	p803.items = []parser{&p827, &p14, &p802}
	var p799 = sequenceParser{id: 799, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var p798 = charParser{id: 798, chars: []rune{41}}
	p799.items = []parser{&p798}
	p804.items = []parser{&p797, &p801, &p827, &p795, &p803, &p827, &p799}
	p795.options = []parser{&p187, &p188, &p434, &p493, &p552, &p603, &p741, &p773, &p784, &p804, &p785}
	var p812 = sequenceParser{id: 812, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var p810 = sequenceParser{id: 810, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	p810.items = []parser{&p809, &p827, &p795}
	var p811 = sequenceParser{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	p811.items = []parser{&p827, &p810}
	p812.items = []parser{&p827, &p810, &p811}
	p813.items = []parser{&p795, &p812}
	p828.items = []parser{&p824, &p827, &p809, &p827, &p813, &p827, &p809}
	p829.items = []parser{&p827, &p828, &p827}
	var b829 = sequenceBuilder{id: 829, commit: 32, name: "mml", ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b827 = choiceBuilder{id: 827, commit: 2}
	var b825 = choiceBuilder{id: 825, commit: 70}
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
	b825.options = []builder{&b2, &b4, &b6, &b8, &b10, &b12}
	var b826 = sequenceBuilder{id: 826, commit: 70, ranges: [][]int{{1, 1}, {1, 1}}}
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
	b40.items = []builder{&b14, &b827, &b39}
	var b41 = sequenceBuilder{id: 41, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b41.items = []builder{&b827, &b40}
	b42.items = []builder{&b827, &b40, &b41}
	b43.items = []builder{&b39, &b42}
	b826.items = []builder{&b43}
	b827.options = []builder{&b825, &b826}
	var b828 = sequenceBuilder{id: 828, commit: 66, ranges: [][]int{{0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}}}
	var b824 = sequenceBuilder{id: 824, commit: 64, name: "shebang", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b821 = sequenceBuilder{id: 821, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b819 = charBuilder{}
	var b820 = charBuilder{}
	b821.items = []builder{&b819, &b820}
	var b818 = sequenceBuilder{id: 818, commit: 64, name: "shebang-command", ranges: [][]int{{0, 1}}}
	var b817 = sequenceBuilder{id: 817, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b815 = sequenceBuilder{id: 815, commit: 2, allChars: true, ranges: [][]int{{1, 1}}}
	var b814 = charBuilder{}
	b815.items = []builder{&b814}
	var b816 = sequenceBuilder{id: 816, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b816.items = []builder{&b827, &b815}
	b817.items = []builder{&b815, &b816}
	b818.items = []builder{&b817}
	var b823 = sequenceBuilder{id: 823, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b822 = charBuilder{}
	b823.items = []builder{&b822}
	b824.items = []builder{&b821, &b827, &b818, &b827, &b823}
	var b809 = sequenceBuilder{id: 809, commit: 66, ranges: [][]int{{1, 1}, {0, -1}}}
	var b807 = choiceBuilder{id: 807, commit: 2}
	var b806 = sequenceBuilder{id: 806, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b805 = charBuilder{}
	b806.items = []builder{&b805}
	b807.options = []builder{&b806, &b14}
	var b808 = sequenceBuilder{id: 808, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b808.items = []builder{&b827, &b807}
	b809.items = []builder{&b807, &b808}
	var b813 = sequenceBuilder{id: 813, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b795 = choiceBuilder{id: 795, commit: 66}
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
	b15.options = []builder{&b825, &b14}
	b183.items = []builder{&b182, &b15}
	var b186 = sequenceBuilder{id: 186, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b185 = sequenceBuilder{id: 185, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b184 = sequenceBuilder{id: 184, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b184.items = []builder{&b827, &b14}
	b185.items = []builder{&b14, &b184}
	var b403 = choiceBuilder{id: 403, commit: 66}
	var b274 = choiceBuilder{id: 274, commit: 66}
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
	var b516 = sequenceBuilder{id: 516, commit: 64, name: "receive", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b508 = sequenceBuilder{id: 508, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b507 = sequenceBuilder{id: 507, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b500 = charBuilder{}
	var b501 = charBuilder{}
	var b502 = charBuilder{}
	var b503 = charBuilder{}
	var b504 = charBuilder{}
	var b505 = charBuilder{}
	var b506 = charBuilder{}
	b507.items = []builder{&b500, &b501, &b502, &b503, &b504, &b505, &b506}
	b508.items = []builder{&b507, &b15}
	var b515 = sequenceBuilder{id: 515, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b514 = sequenceBuilder{id: 514, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b514.items = []builder{&b827, &b14}
	b515.items = []builder{&b827, &b14, &b514}
	b516.items = []builder{&b508, &b515, &b827, &b274}
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
	b113.items = []builder{&b827, &b112}
	b114.items = []builder{&b112, &b113}
	var b119 = sequenceBuilder{id: 119, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b115 = choiceBuilder{id: 115, commit: 66}
	var b109 = sequenceBuilder{id: 109, commit: 64, name: "spread-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b108 = sequenceBuilder{id: 108, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b105 = charBuilder{}
	var b106 = charBuilder{}
	var b107 = charBuilder{}
	b108.items = []builder{&b105, &b106, &b107}
	b109.items = []builder{&b274, &b827, &b108}
	b115.options = []builder{&b403, &b109}
	var b118 = sequenceBuilder{id: 118, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b116 = sequenceBuilder{id: 116, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b116.items = []builder{&b114, &b827, &b115}
	var b117 = sequenceBuilder{id: 117, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b117.items = []builder{&b827, &b116}
	b118.items = []builder{&b827, &b116, &b117}
	b119.items = []builder{&b115, &b118}
	var b123 = sequenceBuilder{id: 123, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b122 = charBuilder{}
	b123.items = []builder{&b122}
	b124.items = []builder{&b121, &b827, &b114, &b827, &b119, &b827, &b114, &b827, &b123}
	b125.items = []builder{&b124}
	var b130 = sequenceBuilder{id: 130, commit: 64, name: "mutable-list", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b127 = sequenceBuilder{id: 127, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b126 = charBuilder{}
	b127.items = []builder{&b126}
	var b129 = sequenceBuilder{id: 129, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b128 = sequenceBuilder{id: 128, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b128.items = []builder{&b827, &b14}
	b129.items = []builder{&b827, &b14, &b128}
	b130.items = []builder{&b127, &b129, &b827, &b124}
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
	b135.items = []builder{&b827, &b14}
	b136.items = []builder{&b827, &b14, &b135}
	var b138 = sequenceBuilder{id: 138, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b137 = sequenceBuilder{id: 137, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b137.items = []builder{&b827, &b14}
	b138.items = []builder{&b827, &b14, &b137}
	var b134 = sequenceBuilder{id: 134, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b133 = charBuilder{}
	b134.items = []builder{&b133}
	b139.items = []builder{&b132, &b136, &b827, &b403, &b138, &b827, &b134}
	b140.options = []builder{&b104, &b87, &b139}
	var b144 = sequenceBuilder{id: 144, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b143 = sequenceBuilder{id: 143, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b143.items = []builder{&b827, &b14}
	b144.items = []builder{&b827, &b14, &b143}
	var b142 = sequenceBuilder{id: 142, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b141 = charBuilder{}
	b142.items = []builder{&b141}
	var b146 = sequenceBuilder{id: 146, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b145 = sequenceBuilder{id: 145, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b145.items = []builder{&b827, &b14}
	b146.items = []builder{&b827, &b14, &b145}
	b147.items = []builder{&b140, &b144, &b827, &b142, &b146, &b827, &b403}
	b148.options = []builder{&b147, &b109}
	var b152 = sequenceBuilder{id: 152, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b150 = sequenceBuilder{id: 150, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b149 = choiceBuilder{id: 149, commit: 2}
	b149.options = []builder{&b147, &b109}
	b150.items = []builder{&b114, &b827, &b149}
	var b151 = sequenceBuilder{id: 151, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b151.items = []builder{&b827, &b150}
	b152.items = []builder{&b827, &b150, &b151}
	b153.items = []builder{&b148, &b152}
	var b157 = sequenceBuilder{id: 157, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b156 = charBuilder{}
	b157.items = []builder{&b156}
	b158.items = []builder{&b155, &b827, &b114, &b827, &b153, &b827, &b114, &b827, &b157}
	b159.items = []builder{&b158}
	var b164 = sequenceBuilder{id: 164, commit: 64, name: "mutable-struct", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b161 = sequenceBuilder{id: 161, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b160 = charBuilder{}
	b161.items = []builder{&b160}
	var b163 = sequenceBuilder{id: 163, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b162 = sequenceBuilder{id: 162, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b162.items = []builder{&b827, &b14}
	b163.items = []builder{&b827, &b14, &b162}
	b164.items = []builder{&b161, &b163, &b827, &b158}
	var b209 = sequenceBuilder{id: 209, commit: 64, name: "function", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b206 = sequenceBuilder{id: 206, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b204 = charBuilder{}
	var b205 = charBuilder{}
	b206.items = []builder{&b204, &b205}
	var b208 = sequenceBuilder{id: 208, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b207 = sequenceBuilder{id: 207, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b207.items = []builder{&b827, &b14}
	b208.items = []builder{&b827, &b14, &b207}
	var b203 = sequenceBuilder{id: 203, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b195 = sequenceBuilder{id: 195, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b194 = charBuilder{}
	b195.items = []builder{&b194}
	var b197 = choiceBuilder{id: 197, commit: 2}
	var b168 = sequenceBuilder{id: 168, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b167 = sequenceBuilder{id: 167, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b165 = sequenceBuilder{id: 165, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b165.items = []builder{&b114, &b827, &b104}
	var b166 = sequenceBuilder{id: 166, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b166.items = []builder{&b827, &b165}
	b167.items = []builder{&b827, &b165, &b166}
	b168.items = []builder{&b104, &b167}
	var b196 = sequenceBuilder{id: 196, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {1, 1}}}
	var b175 = sequenceBuilder{id: 175, commit: 64, name: "collect-parameter", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b172 = sequenceBuilder{id: 172, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b169 = charBuilder{}
	var b170 = charBuilder{}
	var b171 = charBuilder{}
	b172.items = []builder{&b169, &b170, &b171}
	var b174 = sequenceBuilder{id: 174, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b173 = sequenceBuilder{id: 173, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b173.items = []builder{&b827, &b14}
	b174.items = []builder{&b827, &b14, &b173}
	b175.items = []builder{&b172, &b174, &b827, &b104}
	b196.items = []builder{&b168, &b827, &b114, &b827, &b175}
	b197.options = []builder{&b168, &b196, &b175}
	var b199 = sequenceBuilder{id: 199, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b198 = charBuilder{}
	b199.items = []builder{&b198}
	var b202 = sequenceBuilder{id: 202, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b201 = sequenceBuilder{id: 201, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b201.items = []builder{&b827, &b14}
	b202.items = []builder{&b827, &b14, &b201}
	var b200 = choiceBuilder{id: 200, commit: 2}
	var b785 = choiceBuilder{id: 785, commit: 66}
	var b513 = sequenceBuilder{id: 513, commit: 64, name: "send", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b499 = sequenceBuilder{id: 499, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b498 = sequenceBuilder{id: 498, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b494 = charBuilder{}
	var b495 = charBuilder{}
	var b496 = charBuilder{}
	var b497 = charBuilder{}
	b498.items = []builder{&b494, &b495, &b496, &b497}
	b499.items = []builder{&b498, &b15}
	var b510 = sequenceBuilder{id: 510, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b509 = sequenceBuilder{id: 509, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b509.items = []builder{&b827, &b14}
	b510.items = []builder{&b827, &b14, &b509}
	var b512 = sequenceBuilder{id: 512, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b511 = sequenceBuilder{id: 511, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b511.items = []builder{&b827, &b14}
	b512.items = []builder{&b827, &b14, &b511}
	b513.items = []builder{&b499, &b510, &b827, &b274, &b512, &b827, &b274}
	var b566 = sequenceBuilder{id: 566, commit: 64, name: "go", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b556 = sequenceBuilder{id: 556, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b555 = sequenceBuilder{id: 555, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b553 = charBuilder{}
	var b554 = charBuilder{}
	b555.items = []builder{&b553, &b554}
	b556.items = []builder{&b555, &b15}
	var b565 = sequenceBuilder{id: 565, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b564 = sequenceBuilder{id: 564, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b564.items = []builder{&b827, &b14}
	b565.items = []builder{&b827, &b14, &b564}
	var b264 = sequenceBuilder{id: 264, commit: 64, name: "function-application", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b261 = sequenceBuilder{id: 261, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b260 = charBuilder{}
	b261.items = []builder{&b260}
	var b263 = sequenceBuilder{id: 263, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b262 = charBuilder{}
	b263.items = []builder{&b262}
	b264.items = []builder{&b274, &b827, &b261, &b827, &b114, &b827, &b119, &b827, &b114, &b827, &b263}
	b566.items = []builder{&b556, &b565, &b827, &b264}
	var b575 = sequenceBuilder{id: 575, commit: 64, name: "defer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b572 = sequenceBuilder{id: 572, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b567 = charBuilder{}
	var b568 = charBuilder{}
	var b569 = charBuilder{}
	var b570 = charBuilder{}
	var b571 = charBuilder{}
	b572.items = []builder{&b567, &b568, &b569, &b570, &b571}
	var b574 = sequenceBuilder{id: 574, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b573 = sequenceBuilder{id: 573, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b573.items = []builder{&b827, &b14}
	b574.items = []builder{&b827, &b14, &b573}
	b575.items = []builder{&b572, &b574, &b827, &b264}
	var b639 = choiceBuilder{id: 639, commit: 64, name: "assignment"}
	var b623 = sequenceBuilder{id: 623, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b608 = sequenceBuilder{id: 608, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b607 = sequenceBuilder{id: 607, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b604 = charBuilder{}
	var b605 = charBuilder{}
	var b606 = charBuilder{}
	b607.items = []builder{&b604, &b605, &b606}
	b608.items = []builder{&b607, &b15}
	var b622 = sequenceBuilder{id: 622, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b621 = sequenceBuilder{id: 621, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b621.items = []builder{&b827, &b14}
	b622.items = []builder{&b827, &b14, &b621}
	var b616 = sequenceBuilder{id: 616, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b613 = sequenceBuilder{id: 613, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b612 = sequenceBuilder{id: 612, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b611 = sequenceBuilder{id: 611, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b611.items = []builder{&b827, &b14}
	b612.items = []builder{&b14, &b611}
	var b610 = sequenceBuilder{id: 610, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b609 = charBuilder{}
	b610.items = []builder{&b609}
	b613.items = []builder{&b612, &b827, &b610}
	var b615 = sequenceBuilder{id: 615, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b614 = sequenceBuilder{id: 614, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b614.items = []builder{&b827, &b14}
	b615.items = []builder{&b827, &b14, &b614}
	b616.items = []builder{&b274, &b827, &b613, &b615, &b827, &b403}
	b623.items = []builder{&b608, &b622, &b827, &b616}
	var b630 = sequenceBuilder{id: 630, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b627 = sequenceBuilder{id: 627, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b626 = sequenceBuilder{id: 626, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b626.items = []builder{&b827, &b14}
	b627.items = []builder{&b827, &b14, &b626}
	var b625 = sequenceBuilder{id: 625, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b624 = charBuilder{}
	b625.items = []builder{&b624}
	var b629 = sequenceBuilder{id: 629, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b628 = sequenceBuilder{id: 628, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b628.items = []builder{&b827, &b14}
	b629.items = []builder{&b827, &b14, &b628}
	b630.items = []builder{&b274, &b627, &b827, &b625, &b629, &b827, &b403}
	var b638 = sequenceBuilder{id: 638, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b637 = sequenceBuilder{id: 637, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b636 = sequenceBuilder{id: 636, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b636.items = []builder{&b827, &b14}
	b637.items = []builder{&b827, &b14, &b636}
	var b632 = sequenceBuilder{id: 632, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b631 = charBuilder{}
	b632.items = []builder{&b631}
	var b633 = sequenceBuilder{id: 633, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b620 = sequenceBuilder{id: 620, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b619 = sequenceBuilder{id: 619, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b617 = sequenceBuilder{id: 617, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b617.items = []builder{&b114, &b827, &b616}
	var b618 = sequenceBuilder{id: 618, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b618.items = []builder{&b827, &b617}
	b619.items = []builder{&b827, &b617, &b618}
	b620.items = []builder{&b616, &b619}
	b633.items = []builder{&b114, &b827, &b620}
	var b635 = sequenceBuilder{id: 635, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b634 = charBuilder{}
	b635.items = []builder{&b634}
	b638.items = []builder{&b608, &b637, &b827, &b632, &b827, &b633, &b827, &b114, &b827, &b635}
	b639.options = []builder{&b623, &b630, &b638}
	var b794 = sequenceBuilder{id: 794, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b787 = sequenceBuilder{id: 787, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b786 = charBuilder{}
	b787.items = []builder{&b786}
	var b791 = sequenceBuilder{id: 791, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b790 = sequenceBuilder{id: 790, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b790.items = []builder{&b827, &b14}
	b791.items = []builder{&b827, &b14, &b790}
	var b793 = sequenceBuilder{id: 793, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b792 = sequenceBuilder{id: 792, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b792.items = []builder{&b827, &b14}
	b793.items = []builder{&b827, &b14, &b792}
	var b789 = sequenceBuilder{id: 789, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b788 = charBuilder{}
	b789.items = []builder{&b788}
	b794.items = []builder{&b787, &b791, &b827, &b785, &b793, &b827, &b789}
	b785.options = []builder{&b513, &b566, &b575, &b639, &b794, &b403}
	var b193 = sequenceBuilder{id: 193, commit: 64, name: "block", ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b190 = sequenceBuilder{id: 190, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b189 = charBuilder{}
	b190.items = []builder{&b189}
	var b192 = sequenceBuilder{id: 192, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b191 = charBuilder{}
	b192.items = []builder{&b191}
	b193.items = []builder{&b190, &b827, &b809, &b827, &b813, &b827, &b809, &b827, &b192}
	b200.options = []builder{&b785, &b193}
	b203.items = []builder{&b195, &b827, &b114, &b827, &b197, &b827, &b114, &b827, &b199, &b202, &b827, &b200}
	b209.items = []builder{&b206, &b208, &b827, &b203}
	var b219 = sequenceBuilder{id: 219, commit: 64, name: "effect", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b212 = sequenceBuilder{id: 212, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b210 = charBuilder{}
	var b211 = charBuilder{}
	b212.items = []builder{&b210, &b211}
	var b216 = sequenceBuilder{id: 216, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b215 = sequenceBuilder{id: 215, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b215.items = []builder{&b827, &b14}
	b216.items = []builder{&b827, &b14, &b215}
	var b214 = sequenceBuilder{id: 214, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b213 = charBuilder{}
	b214.items = []builder{&b213}
	var b218 = sequenceBuilder{id: 218, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b217 = sequenceBuilder{id: 217, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b217.items = []builder{&b827, &b14}
	b218.items = []builder{&b827, &b14, &b217}
	b219.items = []builder{&b212, &b216, &b827, &b214, &b218, &b827, &b203}
	var b259 = sequenceBuilder{id: 259, commit: 64, name: "indexer", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b258 = sequenceBuilder{id: 258, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b257 = sequenceBuilder{id: 257, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b257.items = []builder{&b827, &b14}
	b258.items = []builder{&b827, &b14, &b257}
	var b256 = sequenceBuilder{id: 256, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}}}
	var b252 = choiceBuilder{id: 252, commit: 66}
	var b233 = sequenceBuilder{id: 233, commit: 64, name: "symbol-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b230 = sequenceBuilder{id: 230, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b229 = charBuilder{}
	b230.items = []builder{&b229}
	var b232 = sequenceBuilder{id: 232, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b231 = sequenceBuilder{id: 231, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b231.items = []builder{&b827, &b14}
	b232.items = []builder{&b827, &b14, &b231}
	b233.items = []builder{&b230, &b232, &b827, &b104}
	var b242 = sequenceBuilder{id: 242, commit: 64, name: "expression-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b235 = sequenceBuilder{id: 235, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b234 = charBuilder{}
	b235.items = []builder{&b234}
	var b239 = sequenceBuilder{id: 239, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b238 = sequenceBuilder{id: 238, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b238.items = []builder{&b827, &b14}
	b239.items = []builder{&b827, &b14, &b238}
	var b241 = sequenceBuilder{id: 241, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b240 = sequenceBuilder{id: 240, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b240.items = []builder{&b827, &b14}
	b241.items = []builder{&b827, &b14, &b240}
	var b237 = sequenceBuilder{id: 237, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b236 = charBuilder{}
	b237.items = []builder{&b236}
	b242.items = []builder{&b235, &b239, &b827, &b403, &b241, &b827, &b237}
	var b251 = sequenceBuilder{id: 251, commit: 64, name: "range-index", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b244 = sequenceBuilder{id: 244, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b243 = charBuilder{}
	b244.items = []builder{&b243}
	var b248 = sequenceBuilder{id: 248, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b247 = sequenceBuilder{id: 247, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b247.items = []builder{&b827, &b14}
	b248.items = []builder{&b827, &b14, &b247}
	var b228 = sequenceBuilder{id: 228, commit: 66, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b220 = sequenceBuilder{id: 220, commit: 64, name: "range-from", ranges: [][]int{{1, 1}}}
	b220.items = []builder{&b403}
	var b225 = sequenceBuilder{id: 225, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b224 = sequenceBuilder{id: 224, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b224.items = []builder{&b827, &b14}
	b225.items = []builder{&b827, &b14, &b224}
	var b223 = sequenceBuilder{id: 223, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b222 = charBuilder{}
	b223.items = []builder{&b222}
	var b227 = sequenceBuilder{id: 227, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b226 = sequenceBuilder{id: 226, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b226.items = []builder{&b827, &b14}
	b227.items = []builder{&b827, &b14, &b226}
	var b221 = sequenceBuilder{id: 221, commit: 64, name: "range-to", ranges: [][]int{{1, 1}}}
	b221.items = []builder{&b403}
	b228.items = []builder{&b220, &b225, &b827, &b223, &b227, &b827, &b221}
	var b250 = sequenceBuilder{id: 250, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b249 = sequenceBuilder{id: 249, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b249.items = []builder{&b827, &b14}
	b250.items = []builder{&b827, &b14, &b249}
	var b246 = sequenceBuilder{id: 246, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b245 = charBuilder{}
	b246.items = []builder{&b245}
	b251.items = []builder{&b244, &b248, &b827, &b228, &b250, &b827, &b246}
	b252.options = []builder{&b233, &b242, &b251}
	var b255 = sequenceBuilder{id: 255, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b254 = sequenceBuilder{id: 254, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b253 = sequenceBuilder{id: 253, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b253.items = []builder{&b827, &b14}
	b254.items = []builder{&b14, &b253}
	b255.items = []builder{&b254, &b827, &b252}
	b256.items = []builder{&b252, &b827, &b255}
	b259.items = []builder{&b274, &b258, &b827, &b256}
	var b273 = sequenceBuilder{id: 273, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b266 = sequenceBuilder{id: 266, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b265 = charBuilder{}
	b266.items = []builder{&b265}
	var b270 = sequenceBuilder{id: 270, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b269 = sequenceBuilder{id: 269, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b269.items = []builder{&b827, &b14}
	b270.items = []builder{&b827, &b14, &b269}
	var b272 = sequenceBuilder{id: 272, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b271 = sequenceBuilder{id: 271, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b271.items = []builder{&b827, &b14}
	b272.items = []builder{&b827, &b14, &b271}
	var b268 = sequenceBuilder{id: 268, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b267 = charBuilder{}
	b268.items = []builder{&b267}
	b273.items = []builder{&b266, &b270, &b827, &b403, &b272, &b827, &b268}
	b274.options = []builder{&b61, &b74, &b87, &b99, &b516, &b104, &b125, &b130, &b159, &b164, &b209, &b219, &b259, &b264, &b273}
	var b334 = sequenceBuilder{id: 334, commit: 64, name: "unary-expression", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b333 = choiceBuilder{id: 333, commit: 66}
	var b293 = sequenceBuilder{id: 293, commit: 72, name: "plus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b292 = charBuilder{}
	b293.items = []builder{&b292}
	var b295 = sequenceBuilder{id: 295, commit: 72, name: "minus", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b294 = charBuilder{}
	b295.items = []builder{&b294}
	var b276 = sequenceBuilder{id: 276, commit: 72, name: "binary-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b275 = charBuilder{}
	b276.items = []builder{&b275}
	var b307 = sequenceBuilder{id: 307, commit: 72, name: "logical-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b306 = charBuilder{}
	b307.items = []builder{&b306}
	b333.options = []builder{&b293, &b295, &b276, &b307}
	b334.items = []builder{&b333, &b827, &b274}
	var b381 = choiceBuilder{id: 381, commit: 66}
	var b352 = sequenceBuilder{id: 352, commit: 64, name: "binary0", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b340 = choiceBuilder{id: 340, commit: 66}
	b340.options = []builder{&b274, &b334}
	var b350 = sequenceBuilder{id: 350, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b347 = sequenceBuilder{id: 347, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b346 = sequenceBuilder{id: 346, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b346.items = []builder{&b827, &b14}
	b347.items = []builder{&b14, &b346}
	var b335 = choiceBuilder{id: 335, commit: 66}
	var b278 = sequenceBuilder{id: 278, commit: 72, name: "binary-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b277 = charBuilder{}
	b278.items = []builder{&b277}
	var b285 = sequenceBuilder{id: 285, commit: 72, name: "and-not", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b283 = charBuilder{}
	var b284 = charBuilder{}
	b285.items = []builder{&b283, &b284}
	var b288 = sequenceBuilder{id: 288, commit: 72, name: "lshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b286 = charBuilder{}
	var b287 = charBuilder{}
	b288.items = []builder{&b286, &b287}
	var b291 = sequenceBuilder{id: 291, commit: 72, name: "rshift", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b289 = charBuilder{}
	var b290 = charBuilder{}
	b291.items = []builder{&b289, &b290}
	var b297 = sequenceBuilder{id: 297, commit: 72, name: "mul", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b296 = charBuilder{}
	b297.items = []builder{&b296}
	var b299 = sequenceBuilder{id: 299, commit: 72, name: "div", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b298 = charBuilder{}
	b299.items = []builder{&b298}
	var b301 = sequenceBuilder{id: 301, commit: 72, name: "mod", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b300 = charBuilder{}
	b301.items = []builder{&b300}
	b335.options = []builder{&b278, &b285, &b288, &b291, &b297, &b299, &b301}
	var b349 = sequenceBuilder{id: 349, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b348 = sequenceBuilder{id: 348, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b348.items = []builder{&b827, &b14}
	b349.items = []builder{&b827, &b14, &b348}
	b350.items = []builder{&b347, &b827, &b335, &b349, &b827, &b340}
	var b351 = sequenceBuilder{id: 351, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b351.items = []builder{&b827, &b350}
	b352.items = []builder{&b340, &b827, &b350, &b351}
	var b359 = sequenceBuilder{id: 359, commit: 64, name: "binary1", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b341 = choiceBuilder{id: 341, commit: 66}
	b341.options = []builder{&b340, &b352}
	var b357 = sequenceBuilder{id: 357, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b354 = sequenceBuilder{id: 354, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b353 = sequenceBuilder{id: 353, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b353.items = []builder{&b827, &b14}
	b354.items = []builder{&b14, &b353}
	var b336 = choiceBuilder{id: 336, commit: 66}
	var b280 = sequenceBuilder{id: 280, commit: 72, name: "binary-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b279 = charBuilder{}
	b280.items = []builder{&b279}
	var b282 = sequenceBuilder{id: 282, commit: 72, name: "xor", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b281 = charBuilder{}
	b282.items = []builder{&b281}
	var b303 = sequenceBuilder{id: 303, commit: 72, name: "add", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b302 = charBuilder{}
	b303.items = []builder{&b302}
	var b305 = sequenceBuilder{id: 305, commit: 72, name: "sub", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b304 = charBuilder{}
	b305.items = []builder{&b304}
	b336.options = []builder{&b280, &b282, &b303, &b305}
	var b356 = sequenceBuilder{id: 356, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b355 = sequenceBuilder{id: 355, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b355.items = []builder{&b827, &b14}
	b356.items = []builder{&b827, &b14, &b355}
	b357.items = []builder{&b354, &b827, &b336, &b356, &b827, &b341}
	var b358 = sequenceBuilder{id: 358, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b358.items = []builder{&b827, &b357}
	b359.items = []builder{&b341, &b827, &b357, &b358}
	var b366 = sequenceBuilder{id: 366, commit: 64, name: "binary2", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b342 = choiceBuilder{id: 342, commit: 66}
	b342.options = []builder{&b341, &b359}
	var b364 = sequenceBuilder{id: 364, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b361 = sequenceBuilder{id: 361, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b360 = sequenceBuilder{id: 360, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b360.items = []builder{&b827, &b14}
	b361.items = []builder{&b14, &b360}
	var b337 = choiceBuilder{id: 337, commit: 66}
	var b310 = sequenceBuilder{id: 310, commit: 72, name: "eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b308 = charBuilder{}
	var b309 = charBuilder{}
	b310.items = []builder{&b308, &b309}
	var b313 = sequenceBuilder{id: 313, commit: 72, name: "not-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b311 = charBuilder{}
	var b312 = charBuilder{}
	b313.items = []builder{&b311, &b312}
	var b315 = sequenceBuilder{id: 315, commit: 72, name: "less", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b314 = charBuilder{}
	b315.items = []builder{&b314}
	var b318 = sequenceBuilder{id: 318, commit: 72, name: "less-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b316 = charBuilder{}
	var b317 = charBuilder{}
	b318.items = []builder{&b316, &b317}
	var b320 = sequenceBuilder{id: 320, commit: 72, name: "greater", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b319 = charBuilder{}
	b320.items = []builder{&b319}
	var b323 = sequenceBuilder{id: 323, commit: 72, name: "greater-or-eq", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b321 = charBuilder{}
	var b322 = charBuilder{}
	b323.items = []builder{&b321, &b322}
	b337.options = []builder{&b310, &b313, &b315, &b318, &b320, &b323}
	var b363 = sequenceBuilder{id: 363, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b362 = sequenceBuilder{id: 362, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b362.items = []builder{&b827, &b14}
	b363.items = []builder{&b827, &b14, &b362}
	b364.items = []builder{&b361, &b827, &b337, &b363, &b827, &b342}
	var b365 = sequenceBuilder{id: 365, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b365.items = []builder{&b827, &b364}
	b366.items = []builder{&b342, &b827, &b364, &b365}
	var b373 = sequenceBuilder{id: 373, commit: 64, name: "binary3", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b343 = choiceBuilder{id: 343, commit: 66}
	b343.options = []builder{&b342, &b366}
	var b371 = sequenceBuilder{id: 371, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b368 = sequenceBuilder{id: 368, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b367 = sequenceBuilder{id: 367, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b367.items = []builder{&b827, &b14}
	b368.items = []builder{&b14, &b367}
	var b338 = sequenceBuilder{id: 338, commit: 66, ranges: [][]int{{1, 1}}}
	var b326 = sequenceBuilder{id: 326, commit: 72, name: "logical-and", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b324 = charBuilder{}
	var b325 = charBuilder{}
	b326.items = []builder{&b324, &b325}
	b338.items = []builder{&b326}
	var b370 = sequenceBuilder{id: 370, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b369 = sequenceBuilder{id: 369, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b369.items = []builder{&b827, &b14}
	b370.items = []builder{&b827, &b14, &b369}
	b371.items = []builder{&b368, &b827, &b338, &b370, &b827, &b343}
	var b372 = sequenceBuilder{id: 372, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b372.items = []builder{&b827, &b371}
	b373.items = []builder{&b343, &b827, &b371, &b372}
	var b380 = sequenceBuilder{id: 380, commit: 64, name: "binary4", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b344 = choiceBuilder{id: 344, commit: 66}
	b344.options = []builder{&b343, &b373}
	var b378 = sequenceBuilder{id: 378, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b375 = sequenceBuilder{id: 375, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b374 = sequenceBuilder{id: 374, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b374.items = []builder{&b827, &b14}
	b375.items = []builder{&b14, &b374}
	var b339 = sequenceBuilder{id: 339, commit: 66, ranges: [][]int{{1, 1}}}
	var b329 = sequenceBuilder{id: 329, commit: 72, name: "logical-or", allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b327 = charBuilder{}
	var b328 = charBuilder{}
	b329.items = []builder{&b327, &b328}
	b339.items = []builder{&b329}
	var b377 = sequenceBuilder{id: 377, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b376 = sequenceBuilder{id: 376, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b376.items = []builder{&b827, &b14}
	b377.items = []builder{&b827, &b14, &b376}
	b378.items = []builder{&b375, &b827, &b339, &b377, &b827, &b344}
	var b379 = sequenceBuilder{id: 379, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b379.items = []builder{&b827, &b378}
	b380.items = []builder{&b344, &b827, &b378, &b379}
	b381.options = []builder{&b352, &b359, &b366, &b373, &b380}
	var b394 = sequenceBuilder{id: 394, commit: 64, name: "ternary-expression", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b387 = sequenceBuilder{id: 387, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b386 = sequenceBuilder{id: 386, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b386.items = []builder{&b827, &b14}
	b387.items = []builder{&b827, &b14, &b386}
	var b383 = sequenceBuilder{id: 383, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b382 = charBuilder{}
	b383.items = []builder{&b382}
	var b389 = sequenceBuilder{id: 389, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b388 = sequenceBuilder{id: 388, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b388.items = []builder{&b827, &b14}
	b389.items = []builder{&b827, &b14, &b388}
	var b391 = sequenceBuilder{id: 391, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b390 = sequenceBuilder{id: 390, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b390.items = []builder{&b827, &b14}
	b391.items = []builder{&b827, &b14, &b390}
	var b385 = sequenceBuilder{id: 385, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b384 = charBuilder{}
	b385.items = []builder{&b384}
	var b393 = sequenceBuilder{id: 393, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b392 = sequenceBuilder{id: 392, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b392.items = []builder{&b827, &b14}
	b393.items = []builder{&b827, &b14, &b392}
	b394.items = []builder{&b403, &b387, &b827, &b383, &b389, &b827, &b403, &b391, &b827, &b385, &b393, &b827, &b403}
	var b402 = sequenceBuilder{id: 402, commit: 64, name: "chaining", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}, {0, -1}}}
	var b395 = choiceBuilder{id: 395, commit: 66}
	b395.options = []builder{&b274, &b334, &b381, &b394}
	var b400 = sequenceBuilder{id: 400, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b397 = sequenceBuilder{id: 397, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b396 = sequenceBuilder{id: 396, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b396.items = []builder{&b827, &b14}
	b397.items = []builder{&b14, &b396}
	var b332 = sequenceBuilder{id: 332, commit: 74, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b330 = charBuilder{}
	var b331 = charBuilder{}
	b332.items = []builder{&b330, &b331}
	var b399 = sequenceBuilder{id: 399, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b398 = sequenceBuilder{id: 398, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b398.items = []builder{&b827, &b14}
	b399.items = []builder{&b827, &b14, &b398}
	b400.items = []builder{&b397, &b827, &b332, &b399, &b827, &b395}
	var b401 = sequenceBuilder{id: 401, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b401.items = []builder{&b827, &b400}
	b402.items = []builder{&b395, &b827, &b400, &b401}
	b403.options = []builder{&b274, &b334, &b381, &b394, &b402}
	b186.items = []builder{&b185, &b827, &b403}
	b187.items = []builder{&b183, &b827, &b186}
	var b188 = sequenceBuilder{id: 188, commit: 64, name: "return-statement", ranges: [][]int{{1, 1}}}
	b188.items = []builder{&b183}
	var b434 = sequenceBuilder{id: 434, commit: 64, name: "if", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b407 = sequenceBuilder{id: 407, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b406 = sequenceBuilder{id: 406, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b404 = charBuilder{}
	var b405 = charBuilder{}
	b406.items = []builder{&b404, &b405}
	b407.items = []builder{&b406, &b15}
	var b429 = sequenceBuilder{id: 429, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b428 = sequenceBuilder{id: 428, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b428.items = []builder{&b827, &b14}
	b429.items = []builder{&b827, &b14, &b428}
	var b431 = sequenceBuilder{id: 431, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b430 = sequenceBuilder{id: 430, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b430.items = []builder{&b827, &b14}
	b431.items = []builder{&b827, &b14, &b430}
	var b433 = sequenceBuilder{id: 433, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b422 = sequenceBuilder{id: 422, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b415 = sequenceBuilder{id: 415, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b414 = sequenceBuilder{id: 414, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b414.items = []builder{&b827, &b14}
	b415.items = []builder{&b14, &b414}
	var b413 = sequenceBuilder{id: 413, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b412 = sequenceBuilder{id: 412, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b408 = charBuilder{}
	var b409 = charBuilder{}
	var b410 = charBuilder{}
	var b411 = charBuilder{}
	b412.items = []builder{&b408, &b409, &b410, &b411}
	b413.items = []builder{&b412, &b15}
	var b417 = sequenceBuilder{id: 417, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b416 = sequenceBuilder{id: 416, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b416.items = []builder{&b827, &b14}
	b417.items = []builder{&b827, &b14, &b416}
	var b419 = sequenceBuilder{id: 419, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b418 = sequenceBuilder{id: 418, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b418.items = []builder{&b827, &b14}
	b419.items = []builder{&b827, &b14, &b418}
	var b421 = sequenceBuilder{id: 421, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b420 = sequenceBuilder{id: 420, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b420.items = []builder{&b827, &b14}
	b421.items = []builder{&b827, &b14, &b420}
	b422.items = []builder{&b415, &b827, &b413, &b417, &b827, &b407, &b419, &b827, &b403, &b421, &b827, &b193}
	var b432 = sequenceBuilder{id: 432, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b432.items = []builder{&b827, &b422}
	b433.items = []builder{&b827, &b422, &b432}
	var b427 = sequenceBuilder{id: 427, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b424 = sequenceBuilder{id: 424, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b423 = sequenceBuilder{id: 423, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b423.items = []builder{&b827, &b14}
	b424.items = []builder{&b14, &b423}
	var b426 = sequenceBuilder{id: 426, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b425 = sequenceBuilder{id: 425, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b425.items = []builder{&b827, &b14}
	b426.items = []builder{&b827, &b14, &b425}
	b427.items = []builder{&b424, &b827, &b413, &b426, &b827, &b193}
	b434.items = []builder{&b407, &b429, &b827, &b403, &b431, &b827, &b193, &b433, &b827, &b427}
	var b493 = sequenceBuilder{id: 493, commit: 64, name: "switch", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b448 = sequenceBuilder{id: 448, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b447 = sequenceBuilder{id: 447, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b441 = charBuilder{}
	var b442 = charBuilder{}
	var b443 = charBuilder{}
	var b444 = charBuilder{}
	var b445 = charBuilder{}
	var b446 = charBuilder{}
	b447.items = []builder{&b441, &b442, &b443, &b444, &b445, &b446}
	b448.items = []builder{&b447, &b15}
	var b490 = sequenceBuilder{id: 490, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b489 = sequenceBuilder{id: 489, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b489.items = []builder{&b827, &b14}
	b490.items = []builder{&b827, &b14, &b489}
	var b492 = sequenceBuilder{id: 492, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b491 = sequenceBuilder{id: 491, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b491.items = []builder{&b827, &b14}
	b492.items = []builder{&b827, &b14, &b491}
	var b480 = sequenceBuilder{id: 480, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b479 = charBuilder{}
	b480.items = []builder{&b479}
	var b486 = sequenceBuilder{id: 486, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b481 = choiceBuilder{id: 481, commit: 2}
	var b478 = sequenceBuilder{id: 478, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b473 = sequenceBuilder{id: 473, commit: 64, name: "case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b440 = sequenceBuilder{id: 440, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b439 = sequenceBuilder{id: 439, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b435 = charBuilder{}
	var b436 = charBuilder{}
	var b437 = charBuilder{}
	var b438 = charBuilder{}
	b439.items = []builder{&b435, &b436, &b437, &b438}
	b440.items = []builder{&b439, &b15}
	var b470 = sequenceBuilder{id: 470, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b469 = sequenceBuilder{id: 469, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b469.items = []builder{&b827, &b14}
	b470.items = []builder{&b827, &b14, &b469}
	var b472 = sequenceBuilder{id: 472, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b471 = sequenceBuilder{id: 471, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b471.items = []builder{&b827, &b14}
	b472.items = []builder{&b827, &b14, &b471}
	var b468 = sequenceBuilder{id: 468, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b467 = charBuilder{}
	b468.items = []builder{&b467}
	b473.items = []builder{&b440, &b470, &b827, &b403, &b472, &b827, &b468}
	var b477 = sequenceBuilder{id: 477, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b475 = sequenceBuilder{id: 475, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b474 = charBuilder{}
	b475.items = []builder{&b474}
	var b476 = sequenceBuilder{id: 476, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b476.items = []builder{&b827, &b475}
	b477.items = []builder{&b827, &b475, &b476}
	b478.items = []builder{&b473, &b477, &b827, &b795}
	var b466 = sequenceBuilder{id: 466, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b461 = sequenceBuilder{id: 461, commit: 64, name: "default", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b456 = sequenceBuilder{id: 456, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b449 = charBuilder{}
	var b450 = charBuilder{}
	var b451 = charBuilder{}
	var b452 = charBuilder{}
	var b453 = charBuilder{}
	var b454 = charBuilder{}
	var b455 = charBuilder{}
	b456.items = []builder{&b449, &b450, &b451, &b452, &b453, &b454, &b455}
	var b460 = sequenceBuilder{id: 460, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b459 = sequenceBuilder{id: 459, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b459.items = []builder{&b827, &b14}
	b460.items = []builder{&b827, &b14, &b459}
	var b458 = sequenceBuilder{id: 458, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b457 = charBuilder{}
	b458.items = []builder{&b457}
	b461.items = []builder{&b456, &b460, &b827, &b458}
	var b465 = sequenceBuilder{id: 465, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b463 = sequenceBuilder{id: 463, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b462 = charBuilder{}
	b463.items = []builder{&b462}
	var b464 = sequenceBuilder{id: 464, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b464.items = []builder{&b827, &b463}
	b465.items = []builder{&b827, &b463, &b464}
	b466.items = []builder{&b461, &b465, &b827, &b795}
	b481.options = []builder{&b478, &b466}
	var b485 = sequenceBuilder{id: 485, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b483 = sequenceBuilder{id: 483, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b482 = choiceBuilder{id: 482, commit: 2}
	b482.options = []builder{&b478, &b466, &b795}
	b483.items = []builder{&b809, &b827, &b482}
	var b484 = sequenceBuilder{id: 484, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b484.items = []builder{&b827, &b483}
	b485.items = []builder{&b827, &b483, &b484}
	b486.items = []builder{&b481, &b485}
	var b488 = sequenceBuilder{id: 488, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b487 = charBuilder{}
	b488.items = []builder{&b487}
	b493.items = []builder{&b448, &b490, &b827, &b403, &b492, &b827, &b480, &b827, &b809, &b827, &b486, &b827, &b809, &b827, &b488}
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
	b550.items = []builder{&b827, &b14}
	b551.items = []builder{&b827, &b14, &b550}
	var b541 = sequenceBuilder{id: 541, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b540 = charBuilder{}
	b541.items = []builder{&b540}
	var b547 = sequenceBuilder{id: 547, commit: 2, ranges: [][]int{{1, 1}, {0, 1}}}
	var b542 = choiceBuilder{id: 542, commit: 2}
	var b532 = sequenceBuilder{id: 532, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {0, 1}}}
	var b527 = sequenceBuilder{id: 527, commit: 64, name: "select-case", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b524 = sequenceBuilder{id: 524, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b523 = sequenceBuilder{id: 523, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b523.items = []builder{&b827, &b14}
	b524.items = []builder{&b827, &b14, &b523}
	var b520 = choiceBuilder{id: 520, commit: 66}
	var b519 = sequenceBuilder{id: 519, commit: 64, name: "receive-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b518 = sequenceBuilder{id: 518, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b517 = sequenceBuilder{id: 517, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b517.items = []builder{&b827, &b14}
	b518.items = []builder{&b827, &b14, &b517}
	b519.items = []builder{&b104, &b518, &b827, &b516}
	b520.options = []builder{&b513, &b516, &b519}
	var b526 = sequenceBuilder{id: 526, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b525 = sequenceBuilder{id: 525, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b525.items = []builder{&b827, &b14}
	b526.items = []builder{&b827, &b14, &b525}
	var b522 = sequenceBuilder{id: 522, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b521 = charBuilder{}
	b522.items = []builder{&b521}
	b527.items = []builder{&b440, &b524, &b827, &b520, &b526, &b827, &b522}
	var b531 = sequenceBuilder{id: 531, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b529 = sequenceBuilder{id: 529, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b528 = charBuilder{}
	b529.items = []builder{&b528}
	var b530 = sequenceBuilder{id: 530, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b530.items = []builder{&b827, &b529}
	b531.items = []builder{&b827, &b529, &b530}
	b532.items = []builder{&b527, &b531, &b827, &b795}
	b542.options = []builder{&b532, &b466}
	var b546 = sequenceBuilder{id: 546, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b544 = sequenceBuilder{id: 544, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b543 = choiceBuilder{id: 543, commit: 2}
	b543.options = []builder{&b532, &b466, &b795}
	b544.items = []builder{&b809, &b827, &b543}
	var b545 = sequenceBuilder{id: 545, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b545.items = []builder{&b827, &b544}
	b546.items = []builder{&b827, &b544, &b545}
	b547.items = []builder{&b542, &b546}
	var b549 = sequenceBuilder{id: 549, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b548 = charBuilder{}
	b549.items = []builder{&b548}
	b552.items = []builder{&b539, &b551, &b827, &b541, &b827, &b809, &b827, &b547, &b827, &b809, &b827, &b549}
	var b603 = sequenceBuilder{id: 603, commit: 64, name: "loop", ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b584 = sequenceBuilder{id: 584, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b583 = sequenceBuilder{id: 583, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b580 = charBuilder{}
	var b581 = charBuilder{}
	var b582 = charBuilder{}
	b583.items = []builder{&b580, &b581, &b582}
	b584.items = []builder{&b583, &b15}
	var b602 = choiceBuilder{id: 602, commit: 2}
	var b598 = sequenceBuilder{id: 598, commit: 2, ranges: [][]int{{0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b595 = sequenceBuilder{id: 595, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b594 = sequenceBuilder{id: 594, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b593 = sequenceBuilder{id: 593, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b593.items = []builder{&b827, &b14}
	b594.items = []builder{&b14, &b593}
	var b592 = choiceBuilder{id: 592, commit: 66}
	var b591 = choiceBuilder{id: 591, commit: 64, name: "range-over-expression"}
	var b590 = sequenceBuilder{id: 590, commit: 2, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b587 = sequenceBuilder{id: 587, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b586 = sequenceBuilder{id: 586, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b586.items = []builder{&b827, &b14}
	b587.items = []builder{&b827, &b14, &b586}
	var b579 = sequenceBuilder{id: 579, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b578 = sequenceBuilder{id: 578, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b576 = charBuilder{}
	var b577 = charBuilder{}
	b578.items = []builder{&b576, &b577}
	b579.items = []builder{&b578, &b15}
	var b589 = sequenceBuilder{id: 589, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b588 = sequenceBuilder{id: 588, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b588.items = []builder{&b827, &b14}
	b589.items = []builder{&b827, &b14, &b588}
	var b585 = choiceBuilder{id: 585, commit: 2}
	b585.options = []builder{&b403, &b228}
	b590.items = []builder{&b104, &b587, &b827, &b579, &b589, &b827, &b585}
	b591.options = []builder{&b590, &b228}
	b592.options = []builder{&b403, &b591}
	b595.items = []builder{&b594, &b827, &b592}
	var b597 = sequenceBuilder{id: 597, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b596 = sequenceBuilder{id: 596, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b596.items = []builder{&b827, &b14}
	b597.items = []builder{&b827, &b14, &b596}
	b598.items = []builder{&b595, &b597, &b827, &b193}
	var b601 = sequenceBuilder{id: 601, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b600 = sequenceBuilder{id: 600, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b599 = sequenceBuilder{id: 599, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b599.items = []builder{&b827, &b14}
	b600.items = []builder{&b14, &b599}
	b601.items = []builder{&b600, &b827, &b193}
	b602.options = []builder{&b598, &b601}
	b603.items = []builder{&b584, &b827, &b602}
	var b741 = choiceBuilder{id: 741, commit: 66}
	var b662 = sequenceBuilder{id: 662, commit: 64, name: "value-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b644 = sequenceBuilder{id: 644, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b643 = sequenceBuilder{id: 643, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b640 = charBuilder{}
	var b641 = charBuilder{}
	var b642 = charBuilder{}
	b643.items = []builder{&b640, &b641, &b642}
	b644.items = []builder{&b643, &b15}
	var b661 = sequenceBuilder{id: 661, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b660 = sequenceBuilder{id: 660, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b660.items = []builder{&b827, &b14}
	b661.items = []builder{&b827, &b14, &b660}
	var b659 = choiceBuilder{id: 659, commit: 2}
	var b653 = sequenceBuilder{id: 653, commit: 64, name: "value-capture", ranges: [][]int{{1, 1}}}
	var b652 = sequenceBuilder{id: 652, commit: 66, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b649 = sequenceBuilder{id: 649, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b648 = sequenceBuilder{id: 648, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b647 = sequenceBuilder{id: 647, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b647.items = []builder{&b827, &b14}
	b648.items = []builder{&b14, &b647}
	var b646 = sequenceBuilder{id: 646, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b645 = charBuilder{}
	b646.items = []builder{&b645}
	b649.items = []builder{&b648, &b827, &b646}
	var b651 = sequenceBuilder{id: 651, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b650 = sequenceBuilder{id: 650, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b650.items = []builder{&b827, &b14}
	b651.items = []builder{&b827, &b14, &b650}
	b652.items = []builder{&b104, &b827, &b649, &b651, &b827, &b403}
	b653.items = []builder{&b652}
	var b658 = sequenceBuilder{id: 658, commit: 64, name: "mutable-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b655 = sequenceBuilder{id: 655, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b654 = charBuilder{}
	b655.items = []builder{&b654}
	var b657 = sequenceBuilder{id: 657, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b656 = sequenceBuilder{id: 656, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b656.items = []builder{&b827, &b14}
	b657.items = []builder{&b827, &b14, &b656}
	b658.items = []builder{&b655, &b657, &b827, &b652}
	b659.options = []builder{&b653, &b658}
	b662.items = []builder{&b644, &b661, &b827, &b659}
	var b679 = sequenceBuilder{id: 679, commit: 64, name: "value-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b678 = sequenceBuilder{id: 678, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b677 = sequenceBuilder{id: 677, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b677.items = []builder{&b827, &b14}
	b678.items = []builder{&b827, &b14, &b677}
	var b674 = sequenceBuilder{id: 674, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b673 = charBuilder{}
	b674.items = []builder{&b673}
	var b668 = sequenceBuilder{id: 668, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b663 = choiceBuilder{id: 663, commit: 2}
	b663.options = []builder{&b653, &b658}
	var b667 = sequenceBuilder{id: 667, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b665 = sequenceBuilder{id: 665, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	var b664 = choiceBuilder{id: 664, commit: 2}
	b664.options = []builder{&b653, &b658}
	b665.items = []builder{&b114, &b827, &b664}
	var b666 = sequenceBuilder{id: 666, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b666.items = []builder{&b827, &b665}
	b667.items = []builder{&b827, &b665, &b666}
	b668.items = []builder{&b663, &b667}
	var b676 = sequenceBuilder{id: 676, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b675 = charBuilder{}
	b676.items = []builder{&b675}
	b679.items = []builder{&b644, &b678, &b827, &b674, &b827, &b114, &b827, &b668, &b827, &b114, &b827, &b676}
	var b690 = sequenceBuilder{id: 690, commit: 64, name: "mutable-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b687 = sequenceBuilder{id: 687, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b686 = sequenceBuilder{id: 686, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b686.items = []builder{&b827, &b14}
	b687.items = []builder{&b827, &b14, &b686}
	var b681 = sequenceBuilder{id: 681, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b680 = charBuilder{}
	b681.items = []builder{&b680}
	var b689 = sequenceBuilder{id: 689, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b688 = sequenceBuilder{id: 688, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b688.items = []builder{&b827, &b14}
	b689.items = []builder{&b827, &b14, &b688}
	var b683 = sequenceBuilder{id: 683, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b682 = charBuilder{}
	b683.items = []builder{&b682}
	var b672 = sequenceBuilder{id: 672, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b671 = sequenceBuilder{id: 671, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b669 = sequenceBuilder{id: 669, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b669.items = []builder{&b114, &b827, &b653}
	var b670 = sequenceBuilder{id: 670, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b670.items = []builder{&b827, &b669}
	b671.items = []builder{&b827, &b669, &b670}
	b672.items = []builder{&b653, &b671}
	var b685 = sequenceBuilder{id: 685, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b684 = charBuilder{}
	b685.items = []builder{&b684}
	b690.items = []builder{&b644, &b687, &b827, &b681, &b689, &b827, &b683, &b827, &b114, &b827, &b672, &b827, &b114, &b827, &b685}
	var b706 = sequenceBuilder{id: 706, commit: 64, name: "function-definition", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b702 = sequenceBuilder{id: 702, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b700 = charBuilder{}
	var b701 = charBuilder{}
	b702.items = []builder{&b700, &b701}
	var b705 = sequenceBuilder{id: 705, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b704 = sequenceBuilder{id: 704, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b704.items = []builder{&b827, &b14}
	b705.items = []builder{&b827, &b14, &b704}
	var b703 = choiceBuilder{id: 703, commit: 2}
	var b694 = sequenceBuilder{id: 694, commit: 64, name: "function-capture", ranges: [][]int{{1, 1}}}
	var b693 = sequenceBuilder{id: 693, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b692 = sequenceBuilder{id: 692, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b691 = sequenceBuilder{id: 691, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b691.items = []builder{&b827, &b14}
	b692.items = []builder{&b827, &b14, &b691}
	b693.items = []builder{&b104, &b692, &b827, &b203}
	b694.items = []builder{&b693}
	var b699 = sequenceBuilder{id: 699, commit: 64, name: "effect-capture", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b696 = sequenceBuilder{id: 696, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b695 = charBuilder{}
	b696.items = []builder{&b695}
	var b698 = sequenceBuilder{id: 698, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b697 = sequenceBuilder{id: 697, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b697.items = []builder{&b827, &b14}
	b698.items = []builder{&b827, &b14, &b697}
	b699.items = []builder{&b696, &b698, &b827, &b693}
	b703.options = []builder{&b694, &b699}
	b706.items = []builder{&b702, &b705, &b827, &b703}
	var b726 = sequenceBuilder{id: 726, commit: 64, name: "function-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b719 = sequenceBuilder{id: 719, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b717 = charBuilder{}
	var b718 = charBuilder{}
	b719.items = []builder{&b717, &b718}
	var b725 = sequenceBuilder{id: 725, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b724 = sequenceBuilder{id: 724, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b724.items = []builder{&b827, &b14}
	b725.items = []builder{&b827, &b14, &b724}
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
	b713.items = []builder{&b114, &b827, &b712}
	var b714 = sequenceBuilder{id: 714, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b714.items = []builder{&b827, &b713}
	b715.items = []builder{&b827, &b713, &b714}
	b716.items = []builder{&b711, &b715}
	var b723 = sequenceBuilder{id: 723, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b722 = charBuilder{}
	b723.items = []builder{&b722}
	b726.items = []builder{&b719, &b725, &b827, &b721, &b827, &b114, &b827, &b716, &b827, &b114, &b827, &b723}
	var b740 = sequenceBuilder{id: 740, commit: 64, name: "effect-definition-group", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b729 = sequenceBuilder{id: 729, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b727 = charBuilder{}
	var b728 = charBuilder{}
	b729.items = []builder{&b727, &b728}
	var b737 = sequenceBuilder{id: 737, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b736 = sequenceBuilder{id: 736, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b736.items = []builder{&b827, &b14}
	b737.items = []builder{&b827, &b14, &b736}
	var b731 = sequenceBuilder{id: 731, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b730 = charBuilder{}
	b731.items = []builder{&b730}
	var b739 = sequenceBuilder{id: 739, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b738 = sequenceBuilder{id: 738, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b738.items = []builder{&b827, &b14}
	b739.items = []builder{&b827, &b14, &b738}
	var b733 = sequenceBuilder{id: 733, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b732 = charBuilder{}
	b733.items = []builder{&b732}
	var b710 = sequenceBuilder{id: 710, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b709 = sequenceBuilder{id: 709, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b707 = sequenceBuilder{id: 707, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b707.items = []builder{&b114, &b827, &b694}
	var b708 = sequenceBuilder{id: 708, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b708.items = []builder{&b827, &b707}
	b709.items = []builder{&b827, &b707, &b708}
	b710.items = []builder{&b694, &b709}
	var b735 = sequenceBuilder{id: 735, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b734 = charBuilder{}
	b735.items = []builder{&b734}
	b740.items = []builder{&b729, &b737, &b827, &b731, &b739, &b827, &b733, &b827, &b114, &b827, &b710, &b827, &b114, &b827, &b735}
	b741.options = []builder{&b662, &b679, &b690, &b706, &b726, &b740}
	var b773 = choiceBuilder{id: 773, commit: 64, name: "use"}
	var b765 = sequenceBuilder{id: 765, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b746 = sequenceBuilder{id: 746, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b745 = sequenceBuilder{id: 745, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b742 = charBuilder{}
	var b743 = charBuilder{}
	var b744 = charBuilder{}
	b745.items = []builder{&b742, &b743, &b744}
	b746.items = []builder{&b745, &b15}
	var b764 = sequenceBuilder{id: 764, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b763 = sequenceBuilder{id: 763, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b763.items = []builder{&b827, &b14}
	b764.items = []builder{&b827, &b14, &b763}
	var b758 = choiceBuilder{id: 758, commit: 64, name: "use-fact"}
	var b757 = sequenceBuilder{id: 757, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {0, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b749 = choiceBuilder{id: 749, commit: 2}
	var b748 = sequenceBuilder{id: 748, commit: 72, name: "use-inline", allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b747 = charBuilder{}
	b748.items = []builder{&b747}
	b749.options = []builder{&b104, &b748}
	var b754 = sequenceBuilder{id: 754, commit: 2, ranges: [][]int{{0, 1}, {0, -1}, {1, 1}}}
	var b753 = sequenceBuilder{id: 753, commit: 2, ranges: [][]int{{1, 1}, {0, -1}}}
	var b752 = sequenceBuilder{id: 752, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b752.items = []builder{&b827, &b14}
	b753.items = []builder{&b14, &b752}
	var b751 = sequenceBuilder{id: 751, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b750 = charBuilder{}
	b751.items = []builder{&b750}
	b754.items = []builder{&b753, &b827, &b751}
	var b756 = sequenceBuilder{id: 756, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b755 = sequenceBuilder{id: 755, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b755.items = []builder{&b827, &b14}
	b756.items = []builder{&b827, &b14, &b755}
	b757.items = []builder{&b749, &b827, &b754, &b756, &b827, &b87}
	b758.options = []builder{&b87, &b757}
	b765.items = []builder{&b746, &b764, &b827, &b758}
	var b772 = sequenceBuilder{id: 772, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {0, 1}, {0, -1}, {1, 1}}}
	var b771 = sequenceBuilder{id: 771, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b770 = sequenceBuilder{id: 770, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b770.items = []builder{&b827, &b14}
	b771.items = []builder{&b827, &b14, &b770}
	var b767 = sequenceBuilder{id: 767, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b766 = charBuilder{}
	b767.items = []builder{&b766}
	var b762 = sequenceBuilder{id: 762, commit: 66, ranges: [][]int{{1, 1}, {0, 1}}}
	var b761 = sequenceBuilder{id: 761, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b759 = sequenceBuilder{id: 759, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b759.items = []builder{&b114, &b827, &b758}
	var b760 = sequenceBuilder{id: 760, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b760.items = []builder{&b827, &b759}
	b761.items = []builder{&b827, &b759, &b760}
	b762.items = []builder{&b758, &b761}
	var b769 = sequenceBuilder{id: 769, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b768 = charBuilder{}
	b769.items = []builder{&b768}
	b772.items = []builder{&b746, &b771, &b827, &b767, &b827, &b114, &b827, &b762, &b827, &b114, &b827, &b769}
	b773.options = []builder{&b765, &b772}
	var b784 = sequenceBuilder{id: 784, commit: 64, name: "export", ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b781 = sequenceBuilder{id: 781, commit: 74, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b780 = sequenceBuilder{id: 780, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	var b774 = charBuilder{}
	var b775 = charBuilder{}
	var b776 = charBuilder{}
	var b777 = charBuilder{}
	var b778 = charBuilder{}
	var b779 = charBuilder{}
	b780.items = []builder{&b774, &b775, &b776, &b777, &b778, &b779}
	b781.items = []builder{&b780, &b15}
	var b783 = sequenceBuilder{id: 783, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b782 = sequenceBuilder{id: 782, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b782.items = []builder{&b827, &b14}
	b783.items = []builder{&b827, &b14, &b782}
	b784.items = []builder{&b781, &b783, &b827, &b741}
	var b804 = sequenceBuilder{id: 804, commit: 66, ranges: [][]int{{1, 1}, {0, 1}, {0, -1}, {1, 1}, {0, 1}, {0, -1}, {1, 1}}}
	var b797 = sequenceBuilder{id: 797, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b796 = charBuilder{}
	b797.items = []builder{&b796}
	var b801 = sequenceBuilder{id: 801, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b800 = sequenceBuilder{id: 800, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b800.items = []builder{&b827, &b14}
	b801.items = []builder{&b827, &b14, &b800}
	var b803 = sequenceBuilder{id: 803, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b802 = sequenceBuilder{id: 802, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b802.items = []builder{&b827, &b14}
	b803.items = []builder{&b827, &b14, &b802}
	var b799 = sequenceBuilder{id: 799, commit: 10, allChars: true, ranges: [][]int{{1, 1}, {1, 1}}}
	var b798 = charBuilder{}
	b799.items = []builder{&b798}
	b804.items = []builder{&b797, &b801, &b827, &b795, &b803, &b827, &b799}
	b795.options = []builder{&b187, &b188, &b434, &b493, &b552, &b603, &b741, &b773, &b784, &b804, &b785}
	var b812 = sequenceBuilder{id: 812, commit: 2, ranges: [][]int{{0, -1}, {1, 1}, {0, -1}}}
	var b810 = sequenceBuilder{id: 810, commit: 2, ranges: [][]int{{1, 1}, {0, -1}, {1, 1}}}
	b810.items = []builder{&b809, &b827, &b795}
	var b811 = sequenceBuilder{id: 811, commit: 2, ranges: [][]int{{0, -1}, {1, 1}}}
	b811.items = []builder{&b827, &b810}
	b812.items = []builder{&b827, &b810, &b811}
	b813.items = []builder{&b795, &b812}
	b828.items = []builder{&b824, &b827, &b809, &b827, &b813, &b827, &b809}
	b829.items = []builder{&b827, &b828, &b827}

	return parseInput(r, &p829, &b829)
}
