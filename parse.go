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
	primitive(colonNode, colon)
	union(seqSep, nlNode, semicolonNode)
	primitive(commaNode, comma)
	union(listSep, nlNode, commaNode)
	primitive(tildeNode, tilde)
	primitive(openSquareNode, openSquare)
	primitive(closeSquareNode, closeSquare)
	primitive(greaterNode, greater)
	primitive(lessNode, less)
	primitive(openParenNode, openParen)
	primitive(closeParenNode, closeParen)
	primitive(openBraceNode, openBrace)
	primitive(closeBraceNode, closeBrace)
	primitive(dotNode, dot)
	group(spreadNode, dotNode, dotNode, dotNode)

	primitive(symbolWordNode, symbolWord)
	primitive(trueNode, trueWord)
	primitive(falseNode, falseWord)
	primitive(andNode, andWord)
	primitive(orNode, orWord)
	primitive(fnNode, fnWord)
	primitive(switchNode, switchWord)
	primitive(caseNode, caseWord)
	primitive(defaultNode, defaultWord)

	primitive(intNode, intToken)
	primitive(stringNode, stringToken)
	primitive(symbolNode, symbolToken)

	union(boolNode, trueNode, falseNode)
	group(channelNode, lessNode, greaterNode)

	union(staticSymbolNode, symbolNode, stringNode)
	group(dynamicSymbolNode, symbolWordNode, openParenNode, nls, expressionNode, nls, closeParenNode)
	union(symbolExpressionNode, staticSymbolNode, dynamicSymbolNode)

	group(spreadExpressionNode, expressionNode, spreadNode)

	union(listItemNode, expressionNode, spreadExpressionNode, listSep)
	sequence(listSequenceNode, listItemNode)
	group(listNode, openSquareNode, listSequenceNode, closeSquareNode)

	group(mutableListNode, tildeNode, listNode)

	group(structureDefinitionNode, symbolExpressionNode, nls, colonNode, nls, expressionNode)
	union(structureItemNode, structureDefinitionNode, spreadExpressionNode, listSep)
	sequence(structureSequenceNode, structureItemNode)
	group(structureNode, openBraceNode, structureSequenceNode, closeBraceNode)

	group(mutableStructureNode, tildeNode, structureNode)

	union(expressionItemNode, expressionNode, listSep)
	sequence(expressionSequenceNode, expressionItemNode)

	group(andExpressionNode, andNode, openParenNode, expressionSequenceNode, closeParenNode)
	group(orExpressionNode, orNode, openParenNode, expressionSequenceNode, closeParenNode)

	union(staticSymbolItemNode, staticSymbolNode, listSep)
	sequence(staticSymbolSequenceNode, staticSymbolItemNode)
	group(collectSymbolNode, spreadNode, staticSymbolNode)
	optional(collectArgumentNode, collectSymbolNode)
	group(functionBodyNode, openBraceNode, statementSequenceNode, closeBraceNode)
	union(functionValueNode, expressionNode, functionBodyNode)
	group(
		functionFactNode,
		openParenNode,
		staticSymbolSequenceNode,
		collectArgumentNode,
		nls,
		closeParenNode,
		nls,
		functionValueNode,
	)
	group(functionNode, fnNode, nls, functionFactNode)
	group(functionEffectNode, fnNode, tildeNode, nls, functionFactNode)

	group(symbolQueryNode, expressionNode /* due to conflicts with statement sequences nls, */, dotNode, nls, symbolExpressionNode)
	optional(optionalExpressionNode, expressionNode)
	group(rangeExpressionNode, optionalExpressionNode /* nls, */, colonNode, nls, optionalExpressionNode)
	union(queryExpressionNode, expressionNode, rangeExpressionNode)
	group(
		expressionQueryNode,
		expressionNode,
		openSquareNode,
		nls,
		queryExpressionNode,
		nls,
		closeSquareNode,
	)
	union(queryNode, symbolQueryNode, expressionQueryNode)

	union(matchExpressionNode, expressionNode)
	group(switchClauseNode, caseNode, nls, matchExpressionNode, nls, colonNode, statementSequenceNode)
	union(switchClauseSequenceItemNode, switchClauseNode, nlNode)
	sequence(switchClauseSequenceNode, switchClauseSequenceItemNode)
	group(defaultClauseNode, defaultNode, nls, colonNode, statementSequenceNode)
	optional(optionalDefaultClauseNode, defaultClauseNode)
	group(
		switchConditionalNode,
		switchNode,
		nls,
		openBraceNode,
		switchClauseSequenceNode,
		optionalDefaultClauseNode,
		switchClauseSequenceNode,
		closeBraceNode,
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
		mutableListNode,
		structureNode,
		mutableStructureNode,
		andExpressionNode,
		orExpressionNode,
		functionNode,
		functionEffectNode,
		queryNode,
		switchConditionalNode,
	)

	union(statementNode, expressionNode)

	union(sequenceItemNode, statementNode, seqSep)
	sequence(statementSequenceNode, sequenceItemNode)
	union(documentNode, statementSequenceNode)
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
	colonNode
	commaNode
	seqSep
	listSep
	tildeNode
	spreadNode
	openParenNode
	closeParenNode
	openSquareNode
	closeSquareNode
	greaterNode
	lessNode
	openBraceNode
	closeBraceNode
	dotNode

	symbolWordNode
	symbolNode
	intNode
	stringNode
	channelNode
	boolNode
	trueNode
	falseNode
	andNode
	orNode
	fnNode
	switchNode
	caseNode
	defaultNode

	staticSymbolNode
	dynamicSymbolNode
	expressionNode
	symbolExpressionNode
	spreadExpressionNode
	listItemNode
	listSequenceNode
	listNode
	mutableListNode
	structureDefinitionNode
	structureItemNode
	structureSequenceNode
	structureNode
	mutableStructureNode
	expressionItemNode
	expressionSequenceNode
	andExpressionNode
	orExpressionNode
	staticSymbolItemNode
	staticSymbolSequenceNode
	collectSymbolNode
	collectArgumentNode
	functionBodyNode
	functionValueNode
	functionFactNode
	functionNode
	functionEffectNode
	symbolQueryNode
	optionalExpressionNode
	rangeExpressionNode
	queryExpressionNode
	expressionQueryNode
	queryNode
	matchExpressionNode
	switchClauseNode
	defaultClauseNode
	switchConditionalNode
	switchClauseSequenceNode
	optionalDefaultClauseNode
	switchClauseSequenceItemNode

	statementNode

	sequenceItemNode
	statementSequenceNode
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
	current          parser
	parsers          []nodeType
	currentlyParsing []nodeType
}

type sequenceParser struct {
	parserBase
	itemType         nodeType
	parser           parser
	currentlyParsing []nodeType
}

var parsers = make(map[nodeType]func([]nodeType) parser)

func (nt nodeType) String() string {
	switch nt {
	case nlNode:
		return "newline"
	case nls:
		return "nls"
	case semicolonNode:
		return "semicolon"
	case colonNode:
		return "colon"
	case commaNode:
		return "comma"
	case tildeNode:
		return "tilde"
	case spreadNode:
		return "spread"
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
	case openBraceNode:
		return "openBrace"
	case closeBraceNode:
		return "closeBrace"
	case dotNode:
		return "dot"

	case symbolWordNode:
		return "symbolWord"
	case symbolNode:
		return "symbol"
	case intNode:
		return "int"
	case stringNode:
		return "string"
	case channelNode:
		return "channel"
	case boolNode:
		return "bool"
	case trueNode:
		return "true"
	case falseNode:
		return "false"
	case andNode:
		return "and"
	case orNode:
		return "or"
	case fnNode:
		return "fn"
	case switchNode:
		return "switch"
	case caseNode:
		return "case"
	case defaultNode:
		return "default"

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
	case listSequenceNode:
		return "listSequence"
	case listNode:
		return "list"
	case mutableListNode:
		return "mutable-list"
	case structureNode:
		return "structure"
	case mutableStructureNode:
		return "mutable-structure"
	case expressionItemNode:
		return "expressionItem"
	case expressionSequenceNode:
		return "expressionSequence"
	case andExpressionNode:
		return "andExpression"
	case orExpressionNode:
		return "orExpression"
	case staticSymbolItemNode:
		return "staticSymbolItem"
	case staticSymbolSequenceNode:
		return "staticSymbolSequence"
	case collectSymbolNode:
		return "collectSymbol"
	case collectArgumentNode:
		return "collectArgument"
	case functionBodyNode:
		return "functionBody"
	case functionValueNode:
		return "functionValue"
	case functionFactNode:
		return "functionFact"
	case functionNode:
		return "function"
	case functionEffectNode:
		return "effect-function"
	case statementSequenceNode:
		return "statementSequence"
	case symbolQueryNode:
		return "symbolQuery"
	case optionalExpressionNode:
		return "optionalExpression"
	case rangeExpressionNode:
		return "rangeExpression"
	case queryExpressionNode:
		return "queryExpression"
	case expressionQueryNode:
		return "expressionQuery"
	case queryNode:
		return "query"
	case matchExpressionNode:
		return "matchExpression"
	case switchClauseNode:
		return "switchClause"
	case defaultClauseNode:
		return "defaultClause"
	case switchConditionalNode:
		return "switchConditional"
	case switchClauseSequenceNode:
		return "switchClauseSequence"
	case optionalDefaultClauseNode:
		return "optionalDefaultClause"
	case switchClauseSequenceItemNode:
		return "switchClauseSequenceItem"

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

func newOptionalParser(nt nodeType, currentlyParsing []nodeType) parser {
	p := &optionalParser{}
	p.n = node{typ: noNode}
	p.optional = nt
	p.parser = parsers[nt](currentlyParsing)
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

func newUnionParser(nt nodeType, nts []nodeType, currentlyParsing []nodeType) parser {
	p := &unionParser{}
	p.n = node{typ: nt}
	currentlyParsing = append(currentlyParsing, nt)
	for _, nti := range nts {
		var found bool
		for _, cnti := range currentlyParsing {
			if cnti == nti {
				found = true
				break
			}
		}

		if !found {
			p.parsers = append(p.parsers, parsers[nti](currentlyParsing))
		}
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

func newGroupParser(nt nodeType, nts []nodeType, currentlyParsing []nodeType) parser {
	p := &groupParser{}
	p.n = node{typ: nt}
	p.parsers = nts
	p.currentlyParsing = append(currentlyParsing, nt)
	return p
}

func (p *groupParser) accept(t token) bool {
	if p.current == nil {
		if len(p.parsers) == 0 {
			p.v = true
			return false
		}

		p.current, p.parsers = parsers[p.parsers[0]](p.currentlyParsing), p.parsers[1:]
		p.currentlyParsing = nil
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

func newSequenceParser(nt nodeType, itemType nodeType, currentlyParsing []nodeType) parser {
	p := &sequenceParser{}
	p.n = node{typ: nt}
	p.itemType = itemType
	p.currentlyParsing = currentlyParsing
	return p
}

func (p *sequenceParser) accept(t token) bool {
	current := p.parser
	if current == nil {
		current = parsers[p.itemType](p.currentlyParsing)
		p.currentlyParsing = nil
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
	parsers[nt] = func(currentlyParsing []nodeType) parser { return newPrimitiveParser(nt, tt) }
}

func optional(nt, ont nodeType) {
	parsers[nt] = func(currentlyParsing []nodeType) parser { return newOptionalParser(ont, currentlyParsing) }
}

func union(nt nodeType, nts ...nodeType) {
	parsers[nt] = func(currentlyParsing []nodeType) parser { return newUnionParser(nt, nts, currentlyParsing) }
}

func group(nt nodeType, nts ...nodeType) {
	parsers[nt] = func(currentlyParsing []nodeType) parser { return newGroupParser(nt, nts, currentlyParsing) }
}

func sequence(nt nodeType, itemType nodeType) {
	parsers[nt] = func(currentlyParsing []nodeType) parser { return newSequenceParser(nt, itemType, currentlyParsing) }
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

func dropSeps(n []node) []node {
	nn := make([]node, 0, len(n))
	for _, ni := range n {
		switch ni.typ {
		case nlNode, semicolonNode, commaNode, nls:
		default:
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
	n.nodes = dropSeps(n.nodes)
	n.nodes = []node{postParseNode(n.nodes[2])}
	return n
}

func postParseList(n node) node {
	n.nodes = n.nodes[1].nodes
	n.nodes = dropSeps(n.nodes)
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseMutableList(n node) node {
	l := postParseList(n.nodes[1])
	n.nodes = l.nodes
	return n
}

func postParseStructureDefinition(n node) node {
	n.nodes = dropSeps(n.nodes)
	n.nodes = append(n.nodes[0:1], n.nodes[2])
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseStructure(n node) node {
	n.nodes = n.nodes[1].nodes
	n.nodes = dropSeps(n.nodes)
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseMutableStructure(n node) node {
	s := postParseStructure(n.nodes[1])
	n.nodes = s.nodes
	return n
}

func postParseExpressionSequence(n node) node {
	n.nodes = dropSeps(n.nodes)
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseFunctionCall(n node) node {
	n.nodes = dropSeps(n.nodes)
	seq := n.nodes[2]
	seq = postParseExpressionSequence(seq)
	n.nodes = seq.nodes
	return n
}

func postParseStaticSymbolSequence(n node) node {
	n.nodes = dropSeps(n.nodes)
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseCollectSymbolNode(n node) node {
	n.nodes = n.nodes[1:]
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseFunctionFact(n node) node {
	n.nodes = dropSeps(n.nodes)

	argsFact := n.nodes[:len(n.nodes)-1]
	fixedArgs := argsFact[1]
	fixedArgs = postParseStaticSymbolSequence(fixedArgs)
	args := fixedArgs.nodes
	if len(argsFact) == 4 {
		collectArg := argsFact[2]
		collectArg = postParseCollectSymbolNode(collectArg)
		args = append(args, collectArg)
	}

	value := n.nodes[len(n.nodes)-1]
	if value.typ == functionBodyNode {
		value = value.nodes[1]
	}

	value = postParseNode(value)

	n.nodes = append(args, value)
	return n
}

func postParseFunction(n node) node {
	n.nodes = dropSeps(n.nodes)
	f := n.nodes[1]
	f = postParseFunctionFact(f)
	n.nodes = f.nodes
	return n
}

func postParseFunctionEffect(n node) node {
	n.nodes = dropSeps(n.nodes)
	f := n.nodes[2]
	f = postParseFunctionFact(f)
	n.nodes = f.nodes
	return n
}

func postParseRangeExpression(n node) node {
	n.nodes = dropSeps(n.nodes)

	if len(n.nodes) == 1 {
		n.nodes = make([]node, 2)
		return n
	} else if n.nodes[0].typ == colonNode {
		n.nodes[0] = node{}
	} else if len(n.nodes) == 2 {
		n.nodes[1] = node{}
	} else {
		n.nodes[1] = n.nodes[2]
		n.nodes = n.nodes[:2]
	}

	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseQuery(n node) node {
	n.nodes = dropSeps(n.nodes)
	n.nodes = append(n.nodes[:1], n.nodes[2])
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseSwitchClause(n node) node {
	n.nodes = dropSeps(n.nodes)
	seq := n.nodes[3].nodes
	seq = dropSeps(seq)
	n.nodes = append(n.nodes[1:2], seq...)
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseDefaultClause(n node) node {
	n.nodes = dropSeps(n.nodes)
	n.nodes = n.nodes[2].nodes
	n.nodes = dropSeps(n.nodes)
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseSwitch(n node) node {
	// union(matchExpressionNode, expressionNode)
	// group(switchClauseNode, caseNode, nls, matchExpressionNode, nls, colonNode, statementSequenceNode)
	// union(switchClauseSequenceItemNode, switchClauseNode, nlNode)
	// sequence(switchClauseSequenceNode, switchClauseSequenceItemNode)
	// group(defaultClauseNode, defaultNode, nls, colonNode, statementSequenceNode)
	// optional(optionalDefaultClauseNode, defaultClauseNode)
	// group(
	// 	switchConditionalNode,
	// 	switchNode,
	// 	nls,
	// 	openBraceNode,
	// 	switchClauseSequenceNode,
	// 	optionalDefaultClauseNode,
	// 	switchClauseSequenceNode,
	// 	closeBraceNode,
	// )

	n.nodes = dropSeps(n.nodes)

	clauses := n.nodes[2].nodes
	clauses = dropSeps(clauses)
	if len(n.nodes) == 6 {
		defaultClause := n.nodes[3]
		trailingClauses := n.nodes[4].nodes
		trailingClauses = dropSeps(trailingClauses)
		clauses = append(append(clauses, trailingClauses...), defaultClause)
	} else {
		trailingClauses := n.nodes[3].nodes
		trailingClauses = dropSeps(trailingClauses)
		clauses = append(clauses, trailingClauses...)
	}

	n.nodes = clauses
	n.nodes = postParseNodes(n.nodes)
	return n
}

func postParseDocument(n node) node {
	n.nodes = dropSeps(n.nodes)
	n.nodes = postParseNodes(n.nodes)
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
	case mutableListNode:
		return postParseMutableList(n)
	case structureDefinitionNode:
		return postParseStructureDefinition(n)
	case structureNode:
		return postParseStructure(n)
	case mutableStructureNode:
		return postParseMutableStructure(n)
	case andExpressionNode, orExpressionNode:
		return postParseFunctionCall(n)
	case functionNode:
		return postParseFunction(n)
	case functionEffectNode:
		return postParseFunctionEffect(n)
	case symbolQueryNode, expressionQueryNode:
		return postParseQuery(n)
	case rangeExpressionNode:
		return postParseRangeExpression(n)
	case switchClauseNode:
		return postParseSwitchClause(n)
	case defaultClauseNode:
		return postParseDefaultClause(n)
	case switchConditionalNode:
		return postParseSwitch(n)
	case statementSequenceNode:
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
	p := parsers[documentNode](nil)

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
