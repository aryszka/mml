package mml

func toMML(typ string, a ...interface{}) structure {
	var s structure

	s.values = make(map[string]interface{})
	s.values["type"] = typ
	for i := 0; i < len(a); i += 2 {
		s.values[a[i].(string)] = a[i+1]
	}

	return s
}

func mapMML(c []interface{}) list {
	var l list
	for _, ci := range c {
		l.values = append(l.values, codeMML(ci))
	}

	return l
}

func commentMML(comment) structure       { return toMML("comment") }
func primitiveMML(p primitive) structure { return toMML("primitive", "value", p) }
func symbolMML(s symbol) structure       { return toMML("symbol", "name", s.name) }
func spreadMML(s spread) structure       { return toMML("spread", "value", codeMML(s.value)) }

func listMML(l list) structure {
	return toMML("list", "mutable", l.mutable, "values", mapMML(l.values))
}

func expressionKeyMML(k expressionKey) structure {
	return toMML("expression-key", "value", codeMML(k.value))
}

func entryMML(e entry) structure {
	return toMML("entry", "key", codeMML(e.key), "value", codeMML(e.value))
}

func structureMML(s structure) structure {
	return toMML("structure", "mutable", s.mutable, "entries", mapMML(s.entries))
}

func retMML(r ret) structure {
	return toMML("ret", "value", codeMML(r.value))
}

func statementListMML(l statementList) structure {
	return toMML("statement-list", "statements", mapMML(l.statements))
}

func mapStringsMML(s []string) list {
	var i []interface{}
	for _, si := range s {
		i = append(i, si)
	}

	return mapMML(i)
}

func functionMML(f function) structure {
	return toMML("function",
		"primitive", f.primitive,
		"effect", f.effect,
		"params", mapStringsMML(f.params),
		"collectParam", f.collectParam,
		"statement", codeMML(f.statement),
	)
}

func rangeExpressionMML(e rangeExpression) structure {
	m := toMML("range-expression")

	if e.from != nil {
		m.values["from"] = codeMML(e.from)
	}

	if e.to != nil {
		m.values["to"] = codeMML(e.to)
	}

	return m
}

func indexerMML(i indexer) structure {
	return toMML("indexer", "expression", codeMML(i.expression), "index", codeMML(i.index))
}

func functionApplicationMML(a functionApplication) structure {
	return toMML("function-application", "function", codeMML(a.function), "args", mapMML(a.args))
}

func unaryOperatorMML(o unaryOperator) int   { return int(o) }
func binaryOperatorMML(o binaryOperator) int { return int(o) }

func unaryMML(u unary) structure {
	return toMML("unary", "op", codeMML(u.op), "arg", codeMML(u.arg))
}

func binaryMML(b binary) structure {
	return toMML("binary",
		"op", codeMML(b.op),
		"left", codeMML(b.left),
		"right", codeMML(b.right),
	)
}

func condMML(c cond) structure {
	m := toMML("cond",
		"condition", codeMML(c.condition),
		"consequent", codeMML(c.consequent),
	)

	if c.alternative != nil {
		m.values["alternative"] = codeMML(c.alternative)
	}

	return m
}

func switchCaseMML(c switchCase) structure {
	return toMML("switch-case", "expression", codeMML(c.expression), "body", codeMML(c.body))
}

func mapSwitchCasesMML(c []switchCase) list {
	var i []interface{}

	for _, ci := range c {
		i = append(i, ci)
	}

	return mapMML(i)
}

func switchStatementMML(s switchStatement) structure {
	m := toMML("switch-statement",
		"cases", mapSwitchCasesMML(s.cases),
		"defaultStatements", codeMML(s.defaultStatements),
	)

	if s.expression != nil {
		m.values["expression"] = codeMML(s.expression)
	}

	return m
}

func controlStatementMML(s controlStatement) structure {
	return toMML("control-statement", "type", int(s))
}

func rangeOverMML(r rangeOver) structure {
	m := toMML("range-over")

	if r.symbol != "" {
		m.values["symbol"] = r.symbol
	}

	if r.expression != nil {
		m.values["expression"] = r.expression
	}

	return m
}

func loopMML(l loop) structure {
	m := toMML("loop", "body", codeMML(l.body))

	if l.expression != nil {
		m.values["expression"] = codeMML(l.expression)
	}

	return m
}

func definitionMML(d definition) structure {
	return toMML("definition",
		"mutable", d.mutable,
		"symbol", d.symbol,
		"expression", codeMML(d.expression),
	)
}

func mapDefinitionsMML(d []definition) list {
	var i []interface{}

	for _, di := range d {
		i = append(i, di)
	}

	return mapMML(i)
}

func definitionListMML(d definitionList) structure {
	return toMML("definition-list", "definitions", mapDefinitionsMML(d.definitions))
}

func assignMML(a assign) structure {
	return toMML("assign", "capture", codeMML(a.capture), "value", codeMML(a.value))
}

func mapAssignsMML(a []assign) list {
	var i []interface{}

	for _, ai := range a {
		i = append(i, ai)
	}

	return mapMML(i)
}

func assignListMML(a assignList) structure {
	return toMML("assign-list", "assignments", mapAssignsMML(a.assignments))
}

func sendMML(s send) structure {
	return toMML("send", "channel", codeMML(s.channel), "value", codeMML(s.value))
}

func receiveMML(r receive) structure { return toMML("receive", "channel", codeMML(r.channel)) }
func goMML(g goStatement) structure  { return toMML("go", "application", codeMML(g.application)) }

func deferredMML(d deferred) structure {
	return toMML("deferred", "function", codeMML(d.function), "args", mapMML(d.args))
}

func deferMML(d deferStatement) structure {
	return toMML("defer", "application", codeMML(d.application))
}

func selectCaseMML(c selectCase) structure {
	return toMML("select-case", "expression", codeMML(c.expression), "body", c.body)
}

func mapSelectCasesMML(c []selectCase) list {
	var i []interface{}

	for _, ci := range c {
		i = append(i, ci)
	}

	return mapMML(i)
}

func selectMML(s selectStatement) structure {
	return toMML("select",
		"cases", mapSelectCasesMML(s.cases),
		"hasDefault", s.hasDefault,
		"defaultStatements", codeMML(s.defaultStatements),
	)
}

func moduleMML(m module) structure { return toMML("module", "statements", mapMML(m.statements)) }

func codeMML(c interface{}) interface{} {
	switch ct := c.(type) {
	case comment:
		return commentMML(ct)
	case int, float64, string, bool:
		return c
	case primitive:
		return primitiveMML(ct)
	case symbol:
		return symbolMML(ct)
	case spread:
		return spreadMML(ct)
	case list:
		return listMML(ct)
	case expressionKey:
		return expressionKeyMML(ct)
	case entry:
		return entryMML(ct)
	case structure:
		return structureMML(ct)
	case ret:
		return retMML(ct)
	case statementList:
		return statementListMML(ct)
	case function:
		return functionMML(ct)
	case rangeExpression:
		return rangeExpressionMML(ct)
	case indexer:
		return indexerMML(ct)
	case functionApplication:
		return functionApplicationMML(ct)
	case unaryOperator:
		return unaryOperatorMML(ct)
	case unary:
		return unaryMML(ct)
	case binaryOperator:
		return binaryOperatorMML(ct)
	case binary:
		return binaryMML(ct)
	case cond:
		return condMML(ct)
	case switchCase:
		return switchCaseMML(ct)
	case switchStatement:
		return switchStatementMML(ct)
	case controlStatement:
		return controlStatementMML(ct)
	case rangeOver:
		return rangeOverMML(ct)
	case loop:
		return loopMML(ct)
	case definition:
		return definitionMML(ct)
	case definitionList:
		return definitionListMML(ct)
	case assign:
		return assignMML(ct)
	case assignList:
		return assignListMML(ct)
	case send:
		return sendMML(ct)
	case receive:
		return receiveMML(ct)
	case goStatement:
		return goMML(ct)
	case deferred:
		return deferredMML(ct)
	case deferStatement:
		return deferMML(ct)
	case selectCase:
		return selectCaseMML(ct)
	case selectStatement:
		return selectMML(ct)
	case module:
		return moduleMML(ct)
	default:
		panic(errUnsupportedCode)
	}
}
