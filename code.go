package mml

type comment struct{}

type primitive func(*env) (interface{}, error)

type symbol struct {
	name string
}

type spread struct {
	value interface{}
}

type list struct {
	mutable bool
	values  []interface{}
}

type expressionKey struct {
	value interface{}
}

type entry struct {
	key, value interface{}
}

type structure struct {
	mutable bool
	entries []interface{}
	values  map[string]interface{}
}

type ret struct {
	value interface{}
}

type statementList struct {
	statements []interface{}
}

type function struct {
	primitive    primitive
	effect       bool
	params       []string
	args         []interface{}
	collectParam string
	statement    interface{}
	env          *env
}

type rangeExpression struct {
	from, to interface{}
}

type indexer struct {
	expression interface{}
	index      interface{}
}

type functionApplication struct {
	function interface{}
	args     []interface{}
}

type unaryOperator int

const (
	binaryNot unaryOperator = iota
	plus
	minus
	logicalNot
)

type unary struct {
	op  unaryOperator
	arg interface{}
}

type binaryOperator int

const (
	binaryAnd binaryOperator = iota
	binaryOr
	xor
	andNot
	lshift
	rshift
	mul
	div
	mod
	add
	sub
	eq
	notEq
	less
	lessOrEq
	greater
	greaterOrEq
	logicalAnd
	logicalOr
)

type binary struct {
	op          binaryOperator
	left, right interface{}
}

type cond struct {
	condition, consequent, alternative interface{}
	ternary                            bool
}

type switchCase struct {
	expression interface{}
	body       statementList
}

// TODO: needs a hasDefault because then it is impossible whether it has a return value or not
type switchStatement struct {
	expression        interface{}
	cases             []switchCase
	defaultStatements statementList
}

type controlStatement int

const (
	breakStatement controlStatement = iota
	continueStatement
)

type rangeOver struct {
	symbol     string
	expression interface{}
}

type loop struct {
	expression interface{}
	body       statementList
}

type definition struct {
	mutable, exported bool
	symbol            string
	expression        interface{}
}

type definitionList struct {
	definitions []definition
}

type assign struct {
	capture interface{}
	value   interface{}
}

type assignList struct {
	assignments []assign
}

type send struct {
	channel interface{}
	value   interface{}
}

type receive struct {
	channel interface{}
}

type goStatement struct {
	application functionApplication
}

type deferred struct {
	function function
	args     []interface{}
}

type deferStatement struct {
	application functionApplication
}

type selectCase struct {
	expression interface{}
	body       statementList
}

type selectStatement struct {
	cases             []selectCase
	hasDefault        bool
	defaultStatements statementList
}

type use struct {
	path string
}

type useList struct {
	uses []use
}

type module struct {
	text       string
	shebang    string
	statements []interface{}
}
