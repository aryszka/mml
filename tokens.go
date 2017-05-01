package mml

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"text/scanner"
)

type tokenType int

const (
	noToken tokenType = iota

	comment // comment = "//" /[^\n]*/

	nl // nl = "\n"

	aliasWord // alias-word = "alias"
	andWord
	breakWord
	caseWord
	continueWord
	defaultWord
	deferWord
	elseWord
	exportWord
	falseWord
	fnWord
	forWord
	goWord
	ifWord
	importWord
	intWord
	inWord
	letWord
	setWord
	orWord
	panicWord
	receiveWord
	recoverWord
	sendWord
	switchWord
	symbolWord
	testWord
	trueWord
	typeWord

	openParen // open-paren = "("
	closeParen
	openSquare
	closeSquare
	openBrace
	closeBrace

	// as complex: channel = "<" expression ">"
	channel // channel = "<>" // TODO: this should be separated and bufferable

	colon
	comma
	communicate
	dot
	question
	singleEq
	spread
	tilde
	semicolon

	andNot // and-not = "&" "^"
	diff
	div
	doubleAnd
	doubleEq
	doubleOr
	greater
	greaterOrEq
	less
	lessOrEq
	mod
	mul
	not
	notEq
	power
	shiftLeft
	shiftRight
	singleAnd
	singleOr
	sum
	incOne
	decOne

	intToken    // int-token = /-?[1-9][0-9]*/
	stringToken // stringToken = "\"" /([^\\\"]|\\\\|\\\")*/ "\""
	boolToken   // bool-token = "bool"
	symbolToken // symbol-token = /[a-z_][a-z0-9_]*/

	eofTokenType
)

const (
	nlString = "\n"

	aliasWordString    = "alias"
	andWordString      = "and"
	breakWordString    = "break"
	caseWordString     = "case"
	continueWordString = "continue"
	defaultWordString  = "default"
	deferWordString    = "defer"
	elseWordString     = "else"
	exportWordString   = "export"
	falseWordString    = "false"
	fnWordString       = "fn"
	forWordString      = "for"
	goWordString       = "go"
	ifWordString       = "if"
	importWordString   = "import"
	intWordString      = "int"
	inWordString       = "in"
	letWordString      = "let"
	setWordString      = "set"
	orWordString       = "or"
	panicWordString    = "panic"
	receiveWordString  = "receive"
	recoverWordString  = "recover"
	sendWordString     = "send"
	switchWordString   = "switch"
	symbolWordString   = "symbol"
	testWordString     = "test"
	trueWordString     = "true"
	typeWordString     = "type"

	openParenString   = "("
	closeParenString  = ")"
	openSquareString  = "["
	closeSquareString = "]"
	openBraceString   = "{"
	closeBraceString  = "}"

	channelString     = "<>"
	colonString       = ":"
	commaString       = ","
	communicateString = "<-"
	dotString         = "."
	questionString    = "?"
	singleEqString    = "="
	spreadString      = "..."
	tildeString       = "~"
	semicolonString   = ";"

	andNotString      = "&~"
	diffString        = "-"
	divString         = "/"
	doubleAndString   = "&&"
	doubleEqString    = "=="
	doubleOrString    = "||"
	greaterString     = ">"
	greaterOrEqString = ">="
	lessString        = "<"
	lessOrEqString    = "<="
	modString         = "%"
	mulString         = "*"
	notString         = "!"
	notEqString       = "!="
	powerString       = "^"
	shiftLeftString   = "<<"
	shiftRightString  = ">>"
	singleAndString   = "&"
	singleOrString    = "|"
	sumString         = "+"
	incOneString      = "++"
	decOneString      = "--"
)

var (
	lineCommentTokenExpression  = regexp.MustCompile("^//")
	blockCommentTokenExpression = regexp.MustCompile("(?s)^/[*].*[*]/$")
	intTokenExpression          = regexp.MustCompile("^[0-9]+$")
	stringTokenExpression       = regexp.MustCompile("^\".*\"$")
	symbolTokenExpression       = regexp.MustCompile("^[a-zA-Z]\\w*$")
)

type token struct {
	typ                  tokenType
	value                string
	fileName             string
	offset, line, column int
	cache                *cache
}

type tokenReader struct {
	fileName string
	scanner  scanner.Scanner
	err      error
}

func (t token) String() string {
	if t.typ == eofTokenType {
		return "<eof>"
	}

	switch t.value {
	case "":
		return "<empty>"
	case "\n":
		return "<newline>"
	default:
		return t.value
	}
}

func newTokenReader(r io.Reader, fileName string) *tokenReader {
	tr := &tokenReader{fileName: fileName}
	tr.scanner.Init(r)
	tr.scanner.Error = tr.scannerError
	tr.scanner.Filename = fileName
	// tr.scanner.Mode &^= scanner.SkipComments
	tr.scanner.Whitespace &^= 1 << '\n'
	return tr
}

// token or error, no token when eof
func (tr *tokenReader) next() (token, error) {
	var t token

	if tr.err != nil {
		return t, tr.err
	}

	st := tr.scanner.Scan()
	if tr.err != nil {
		return t, tr.err
	}

	if st == scanner.EOF {
		tr.err = io.EOF
		return token{
			typ:      eofTokenType, // this is not needed in this formType
			value:    "eof",
			fileName: tr.fileName,
			line:     tr.scanner.Line,
			column:   tr.scanner.Column,
		}, io.EOF
	}

	t.value = tr.scanner.TokenText()
	t.fileName = tr.fileName
	t.offset = tr.scanner.Offset
	t.line = tr.scanner.Line
	t.column = tr.scanner.Column

	switch {
	case t.value == nlString:
		t.typ = nl

	case t.value == aliasWordString:
		t.typ = aliasWord
	case t.value == andWordString:
		t.typ = andWord
	case t.value == breakWordString:
		t.typ = breakWord
	case t.value == caseWordString:
		t.typ = caseWord
	case t.value == continueWordString:
		t.typ = continueWord
	case t.value == defaultWordString:
		t.typ = defaultWord
	case t.value == deferWordString:
		t.typ = deferWord
	case t.value == elseWordString:
		t.typ = elseWord
	case t.value == exportWordString:
		t.typ = exportWord
	case t.value == falseWordString:
		t.typ = falseWord
	case t.value == fnWordString:
		t.typ = fnWord
	case t.value == forWordString:
		t.typ = forWord
	case t.value == goWordString:
		t.typ = goWord
	case t.value == ifWordString:
		t.typ = ifWord
	case t.value == importWordString:
		t.typ = importWord
	case t.value == intWordString:
		t.typ = intWord
	case t.value == inWordString:
		t.typ = inWord
	case t.value == letWordString:
		t.typ = letWord
	case t.value == setWordString:
		t.typ = setWord
	// case t.value == orWordString:
	// 	t.typ = orWord
	case t.value == panicWordString:
		t.typ = panicWord
	case t.value == receiveWordString:
		t.typ = receiveWord
	case t.value == recoverWordString:
		t.typ = recoverWord
	case t.value == sendWordString:
		t.typ = sendWord
	case t.value == switchWordString:
		t.typ = switchWord
	case t.value == symbolWordString:
		t.typ = symbolWord
	case t.value == testWordString:
		t.typ = testWord
	case t.value == trueWordString:
		t.typ = trueWord
	case t.value == typeWordString:
		t.typ = typeWord

	case t.value == openParenString:
		t.typ = openParen
	case t.value == closeParenString:
		t.typ = closeParen
	case t.value == openSquareString:
		t.typ = openSquare
	case t.value == closeSquareString:
		t.typ = closeSquare
	case t.value == openBraceString:
		t.typ = openBrace
	case t.value == closeBraceString:
		t.typ = closeBrace

	case t.value == channelString:
		t.typ = channel
	case t.value == colonString:
		t.typ = colon
	case t.value == commaString:
		t.typ = comma
	case t.value == communicateString:
		t.typ = communicate
	case t.value == dotString:
		t.typ = dot
	case t.value == questionString:
		t.typ = question
	case t.value == singleEqString:
		t.typ = singleEq
	case t.value == spreadString:
		t.typ = spread
	case t.value == tildeString:
		t.typ = tilde
	case t.value == semicolonString:
		t.typ = semicolon

	case t.value == andNotString:
		t.typ = andNot
	case t.value == diffString:
		t.typ = diff
	case t.value == divString:
		t.typ = div
	case t.value == doubleAndString:
		t.typ = doubleAnd
	case t.value == doubleEqString:
		t.typ = doubleEq
	case t.value == doubleOrString:
		t.typ = doubleOr
	case t.value == greaterString:
		t.typ = greater
	case t.value == greaterOrEqString:
		t.typ = greaterOrEq
	case t.value == lessString:
		t.typ = less
	case t.value == lessOrEqString:
		t.typ = lessOrEq
	case t.value == modString:
		t.typ = mod
	case t.value == mulString:
		t.typ = mul
	case t.value == notString:
		t.typ = not
	case t.value == notEqString:
		t.typ = notEq
	case t.value == powerString:
		t.typ = power
	case t.value == shiftLeftString:
		t.typ = shiftLeft
	case t.value == shiftRightString:
		t.typ = shiftRight
	case t.value == singleAndString:
		t.typ = singleAnd
	case t.value == singleOrString:
		t.typ = singleOr
	case t.value == sumString:
		t.typ = sum
	case t.value == incOneString:
		t.typ = incOne
	case t.value == decOneString:
		t.typ = decOne

	case lineCommentTokenExpression.MatchString(t.value):
		t.typ = comment
	case blockCommentTokenExpression.MatchString(t.value):
		t.typ = comment
	case intTokenExpression.MatchString(t.value):
		t.typ = intToken
	case stringTokenExpression.MatchString(t.value):
		t.typ = stringToken
	case symbolTokenExpression.MatchString(t.value):
		t.typ = symbolToken

	default:
		return t, fmt.Errorf("%s:%d:%d: unknown token: '%s'",
			tr.scanner.Filename, t.line, t.column, t.value)
	}

	return t, nil
}

func (tr *tokenReader) scannerError(_ *scanner.Scanner, msg string) {
	tr.err = errors.New(msg)
}
