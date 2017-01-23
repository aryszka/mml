package mml

import (
	"fmt"
	"io"
	"os"
)

func syntax() {
	primitive(nlNode, nl)
	sequence(nls, nlNode)
	primitive(semicolonNode, semicolon)
	union(sep, nlNode, semicolonNode)
	sequence(sepSequence, sep)
	group(seps, sep, sepSequence)

	primitive(symbolWordNode, symbolWord)
	primitive(spreadNode, spread)
	primitive(commaNode, comma)
	primitive(trueNode, trueWord)
	primitive(falseNode, falseWord)
	primitive(openSquareNode, openSquare)
	primitive(closeSquareNode, closeSquare)
	primitive(greaterNode, greater)
	primitive(lessNode, less)
	primitive(openParenNode, openParen)
	primitive(closeParenNode, closeParen)

	primitive(intNode, intToken)
	primitive(stringNode, stringToken)
	primitive(symbolNode, symbolToken)

	union(boolNode, trueNode, falseNode)
	group(channelNode, lessNode, greaterNode)

	union(staticSymbolNode, symbolNode, stringNode)
	group(
		dynamicSymbolNode,
		symbolWordNode,
		nls,
		openParenNode,
		nls,
		expressionNode,
		nls,
		closeParenNode,
	)
	union(symbolExpressionNode, staticSymbolNode, dynamicSymbolNode)

	group(spreadExpressionNode, spreadNode, nls, expressionNode)
	union(listItemNode, expressionNode, spreadExpressionNode)
	group(listCommaItemNode, nls, commaNode, nls, listItemNode)
	sequence(listCommaSequenceNode, listCommaItemNode)
	group(listSequenceNode, listItemNode, listCommaSequenceNode)
	optional(optionalListSequenceNode, listSequenceNode)
	group(
		listNode,
		openSquareNode,
		nls,
		optionalListSequenceNode,
		nls,
		closeSquareNode,
	)

	union(
		expressionNode,
		intNode,
		stringNode,
		channelNode,
		symbolNode,
		dynamicSymbolNode,
		boolNode,
		listNode,
	)

	union(statementNode, expressionNode)

	group(sequenceSepItemNode, sep, nls, statementNode)
	sequence(sequenceSepSequenceNode, sequenceSepItemNode)
	group(sequenceSequenceNode, statementNode, sequenceSepSequenceNode)
	optional(optionalSequenceSequenceNode, sequenceSequenceNode)
	group(documentNode, nls, optionalSequenceSequenceNode, nls)
}

func init() {
	syntax()
}

type nodeType int

const (
	noNode nodeType = iota

	nlNode
	nls
	semicolonNode
	sep
	sepSequence
	seps

	symbolWordNode
	symbolNode
	spreadNode
	intNode
	stringNode
	channelNode
	boolNode
	commaNode
	trueNode
	falseNode
	openParenNode
	closeParenNode
	openSquareNode
	closeSquareNode
	greaterNode
	lessNode

	staticSymbolNode
	dynamicSymbolNode
	expressionNode
	symbolExpressionNode
	spreadExpressionNode
	listItemNode
	listCommaItemNode
	listCommaSequenceNode
	listSequenceNode
	optionalListSequenceNode
	listNode
	statementNode

	sequenceSepItemNode
	sequenceSepSequenceNode
	sequenceSequenceNode
	optionalSequenceSequenceNode
	documentNode
)

type parseError struct {
	nodeType  nodeType
	tokenType tokenType
	token     token
}

type node struct {
	token
	typ   nodeType
	nodes []node
}

type parser interface {

	// when returned false, should not be called anymore
	accept(t token) bool

	valid() bool
	node() node
	error() error
}

type parserBase struct {
	e error
	n node
	v bool
}

type primitiveParser struct {
	parserBase
}

type optionalParser struct {
	parserBase
	optional nodeType
	parser   parser
}

type unionParser struct {
	parserBase
	parsers []parser
}

type groupParser struct {
	parserBase
	current parser
	parsers []nodeType
}

type sequenceParser struct {
	parserBase
	itemType nodeType
	parser   parser
}

var parsers = make(map[nodeType]func() parser)

func (nt nodeType) String() string {
	switch nt {
	case nlNode:
		return "newline"
	case nls:
		return "nls"
	case semicolonNode:
		return "semicolon"
	case sep:
		return "separator"
	case sepSequence:
		return "separator-sequence"
	case seps:
		return "separators"

	case symbolWordNode:
		return "symbolWord"
	case symbolNode:
		return "symbol"
	case spreadNode:
		return "spread"
	case intNode:
		return "int"
	case stringNode:
		return "string"
	case channelNode:
		return "channel"
	case boolNode:
		return "bool"
	case commaNode:
		return "comma"
	case trueNode:
		return "true"
	case falseNode:
		return "false"
	case openParenNode:
		return "openParen"
	case closeParenNode:
		return "closeParen"
	case openSquareNode:
		return "openSquare"
	case closeSquareNode:
		return "closeSquare"
	case greaterNode:
		return "greater"
	case lessNode:
		return "less"

	case staticSymbolNode:
		return "staticSymbol"
	case dynamicSymbolNode:
		return "dynamicSymbol"
	case expressionNode:
		return "expression"
	case symbolExpressionNode:
		return "symbolExpression"
	case spreadExpressionNode:
		return "spreadExpression"
	case listItemNode:
		return "listItem"
	case listCommaItemNode:
		return "listCommaItem"
	case listCommaSequenceNode:
		return "listCommaSequence"
	case listSequenceNode:
		return "listSequence"
	case optionalListSequenceNode:
		return "optionalListSequence"
	case listNode:
		return "list"
	case statementNode:
		return "statement"
	case documentNode:
		return "document"

	default:
		return "not-a-node"
	}
}

func perror(nt nodeType, tt tokenType, t token) error {
	return &parseError{nodeType: nt, tokenType: tt, token: t}
}

func (pe *parseError) Error() string {
	format := "%s:%d:%d: error when parsing %v, unexpected %v"
	args := []interface{}{
		pe.token.fileName, pe.token.line, pe.token.column,
		pe.nodeType, pe.token,
	}

	if pe.tokenType != noToken {
		format += ", expecting: %v"
		args = append(args, pe.tokenType)
	}

	return fmt.Sprintf(format, args...)
}

func (p *parserBase) valid() bool  { return p.v }
func (p *parserBase) node() node   { return p.n }
func (p *parserBase) error() error { return p.e }

func newPrimitiveParser(nt nodeType, tt tokenType) parser {
	p := &primitiveParser{}
	p.n = node{typ: nt, token: token{typ: tt}}
	return p
}

func (p *primitiveParser) accept(t token) bool {
	if p.v || p.e != nil {
		return false
	}

	if t.typ == p.n.token.typ {
		p.n.token = t
		p.v = true
	} else {
		p.e = perror(p.n.typ, p.n.token.typ, t)
		p.v = false
	}

	return p.v
}

func newOptionalParser(nt nodeType) parser {
	p := &optionalParser{}
	p.n = node{typ: noNode}
	p.optional = nt
	p.parser = parsers[nt]()
	return p
}

func (p *optionalParser) accept(t token) bool {
	if p.parser.accept(t) {
		if p.n.typ == noNode {
			p.n.typ = p.optional
		}

		return true
	}

	if p.n.typ == noNode {
		p.v = true
		return false
	}

	p.v = p.parser.valid()
	p.n = p.parser.node()
	p.e = p.parser.error()
	return false
}

func newUnionParser(nt nodeType, nts ...nodeType) parser {
	p := &unionParser{}
	p.n = node{typ: nt}
	for _, nti := range nts {
		p.parsers = append(p.parsers, parsers[nti]())
	}

	return p
}

func (p *unionParser) accept(t token) bool {
	var accepting []parser
	for _, pi := range p.parsers {
		if pi.accept(t) {
			accepting = append(accepting, pi)
		}
	}

	if len(accepting) > 0 {
		p.parsers = accepting
		return true
	}

	for _, pi := range p.parsers {
		if pi.valid() {
			p.n = pi.node()
			p.v = true
			return false
		}
	}

	p.e = perror(p.n.typ, noToken, t)
	return false
}

func newGroupParser(nt nodeType, nts ...nodeType) parser {
	p := &groupParser{}
	p.n = node{typ: nt}
	p.parsers = nts
	return p
}

func (p *groupParser) accept(t token) bool {
	if p.current == nil {
		if len(p.parsers) == 0 {
			p.v = true
			return false
		}

		p.current, p.parsers = parsers[p.parsers[0]](), p.parsers[1:]
	}

	if p.current.accept(t) {
		return true
	}

	if !p.current.valid() {
		p.e = p.current.error()
		return false
	}

	n := p.current.node()
	if n.typ != noNode {
		p.n.nodes = append(p.n.nodes, n)
		if p.n.token.typ == noToken {
			p.n.token = n.token
		}
	}

	p.current = nil
	return p.accept(t)
}

func newSequenceParser(nt nodeType, itemType nodeType) parser {
	p := &sequenceParser{}
	p.n = node{typ: nt}
	p.itemType = itemType
	return p
}

func (p *sequenceParser) accept(t token) bool {
	current := p.parser
	if current == nil {
		current = parsers[p.itemType]()
	}

	if current.accept(t) {
		p.parser = current
		return true
	}

	if p.parser == nil {
		p.v = true
		return false
	}

	if !current.valid() {
		p.e = current.error()
		return false
	}

	n := current.node()
	p.n.nodes = append(p.n.nodes, n)
	if p.n.token.typ == noToken {
		p.n.token = n.token
	}

	p.parser = nil
	return p.accept(t)
}

func primitive(nt nodeType, tt tokenType) {
	parsers[nt] = func() parser { return newPrimitiveParser(nt, tt) }
}

func optional(nt, ont nodeType) {
	parsers[nt] = func() parser { return newOptionalParser(ont) }
}

func union(nt nodeType, nts ...nodeType) {
	parsers[nt] = func() parser { return newUnionParser(nt, nts...) }
}

func group(nt nodeType, nts ...nodeType) {
	parsers[nt] = func() parser { return newGroupParser(nt, nts...) }
}

func sequence(nt nodeType, itemType nodeType) {
	parsers[nt] = func() parser { return newSequenceParser(nt, itemType) }
}

func unquoteString(s string) string {
	var (
		r       []byte
		escaped bool
	)

	for _, c := range []byte(s)[1 : len(s)-1] {
		if escaped {
			escaped = false
			switch c {
			case 'a':
				c = '\a'
			case 'b':
				c = '\b'
			case 'f':
				c = '\f'
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			case 'v':
				c = '\v'
			case '\\':
				c = '\\'
			}

			r = append(r, c)
		} else {
			switch c {
			case '\\':
				escaped = true
			default:
				r = append(r, c)
			}
		}
	}

	return string(r)
}

func dropNls(n []node) []node {
	nn := make([]node, 0, len(n))
	for _, ni := range n {
		if ni.typ != nls {
			nn = append(nn, ni)
		}
	}

	return nn
}

func dropSeps(n []node) []node {
	nn := make([]node, 0, len(n))
	for _, ni := range n {
		if ni.typ != semicolonNode && ni.typ != nlNode {
			nn = append(nn, ni)
		}
	}

	return nn
}

func postParseString(n node) node {
	n.token.value = unquoteString(n.token.value)
	return n
}

func postParseChannel(n node) node {
	n.nodes = nil
	return n
}

func postParseDynamicSymbol(n node) node {
	n.nodes = dropNls(n.nodes)
	n.nodes = postParseNodes(n.nodes[2:3])
	return n
}

func postParseList(n node) node {
	n.nodes = n.nodes[1 : len(n.nodes)-1]
	n.nodes = dropNls(n.nodes)

	if len(n.nodes) > 0 {
		seq := n.nodes[0]
		first := seq.nodes[0]
		commaSeq := seq.nodes[1].nodes

		n.nodes = make([]node, 1+len(commaSeq))
		n.nodes[0] = first

		for i, ni := range commaSeq {
			ni.nodes = dropNls(ni.nodes)
			n.nodes[1+i] = ni.nodes[1]
		}
	}

	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseDocument(n node) node {
	n.nodes = dropNls(n.nodes)
	if len(n.nodes) == 0 {
		return n
	}

	sseq := n.nodes[0].nodes
	first := sseq[0]
	following := sseq[1].nodes
	rest := make([]node, len(following))
	for i, ni := range following {
		ni.nodes = dropNls(ni.nodes)
		ni.nodes = dropSeps(ni.nodes)
		rest[i] = ni.nodes[0]
	}

	n.nodes = postParseNodes(append([]node{first}, rest...))
	return n
}

func postParseNode(n node) node {
	switch n.typ {
	case stringNode:
		return postParseString(n)
	case channelNode:
		return postParseChannel(n)
	case dynamicSymbolNode:
		return postParseDynamicSymbol(n)
	case listNode:
		return postParseList(n)
	case documentNode:
		return postParseDocument(n)
	default:
		return n
	}
}

func postParseNodes(n []node) []node {
	nn := make([]node, len(n))
	for i, ni := range n {
		nn[i] = postParseNode(ni)
	}

	return nn
}

func parse(r io.Reader, source string) ([]node, error) {
	tr := newTokenReader(r, source)
	p := parsers[documentNode]()

	for {
		t, err := tr.next()
		if err != nil && err != io.EOF {
			return postParseNodes(p.node().nodes), err
		}

		if !p.accept(t) {
			perr := p.error()
			if t.typ != eofToken && perr == nil {
				return nil, perror(documentNode, noToken, t)
			} else if perr != nil {
				err = perr
				return nil, perr
			}

			return postParseNode(p.node()).nodes, nil
		}
	}
}

func parseFile(fileName string) ([]node, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return parse(f, fileName)
}

func parseInput(r io.Reader) ([]node, error) {
	return parse(r, "<input>")
}
