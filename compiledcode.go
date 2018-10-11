package mml

func toCompiled(typ string, a ...interface{}) map[string]interface{} {
	s := make(map[string]interface{})
	s["type"] = typ
	for i := 0; i < len(a); i += 2 {
		s[a[i].(string)] = a[i+1]
	}

	return s
}

func mapCompiled(c []interface{}) []interface{} {
	var l []interface{}
	for _, ci := range c {
		l = append(l, codeCompiled(ci))
	}

	return l
}

func commentCompiled(comment) map[string]interface{}       { return toCompiled("comment") }
func primitiveCompiled(p primitive) map[string]interface{} { return toCompiled("primitive", "value", p) }
func symbolCompiled(s symbol) map[string]interface{}       { return toCompiled("symbol", "name", s.name) }
func spreadCompiled(s spread) map[string]interface{} {
	return toCompiled("spread", "value", codeCompiled(s.value))
}

func listCompiled(l list) map[string]interface{} {
	return toCompiled("list", "mutable", l.mutable, "values", mapCompiled(l.values))
}

func expressionKeyCompiled(k expressionKey) map[string]interface{} {
	return toCompiled("expression-key", "value", codeCompiled(k.value))
}

func entryCompiled(e entry) map[string]interface{} {
	return toCompiled("entry", "key", codeCompiled(e.key), "value", codeCompiled(e.value))
}

func structureCompiled(s structure) map[string]interface{} {
	return toCompiled("structure", "mutable", s.mutable, "entries", mapCompiled(s.entries))
}

func retCompiled(r ret) map[string]interface{} {
	return toCompiled("ret", "value", codeCompiled(r.value))
}

func statementListCompiled(l statementList) map[string]interface{} {
	return toCompiled("statement-list", "statements", mapCompiled(l.statements))
}

func mapStringsCompiled(s []string) []interface{} {
	var i []interface{}
	for _, si := range s {
		i = append(i, si)
	}

	return mapCompiled(i)
}

func functionCompiled(f function) map[string]interface{} {
	return toCompiled("function",
		"primitive", f.primitive,
		"effect", f.effect,
		"params", mapStringsCompiled(f.params),
		"collectParam", f.collectParam,
		"statement", codeCompiled(f.statement),
	)
}

func rangeExpressionCompiled(e rangeExpression) map[string]interface{} {
	m := toCompiled("range-expression")

	if e.from != nil {
		m["from"] = codeCompiled(e.from)
	}

	if e.to != nil {
		m["to"] = codeCompiled(e.to)
	}

	return m
}

func indexerCompiled(i indexer) map[string]interface{} {
	return toCompiled("indexer", "expression", codeCompiled(i.expression), "index", codeCompiled(i.index))
}

func functionApplicationCompiled(a functionApplication) map[string]interface{} {
	return toCompiled("function-application", "function", codeCompiled(a.function), "args", mapCompiled(a.args))
}

func unaryOperatorCompiled(o unaryOperator) int   { return int(o) }
func binaryOperatorCompiled(o binaryOperator) int { return int(o) }

func unaryCompiled(u unary) map[string]interface{} {
	return toCompiled("unary", "op", codeCompiled(u.op), "arg", codeCompiled(u.arg))
}

func binaryCompiled(b binary) map[string]interface{} {
	return toCompiled("binary",
		"op", codeCompiled(b.op),
		"left", codeCompiled(b.left),
		"right", codeCompiled(b.right),
	)
}

func condCompiled(c cond) map[string]interface{} {
	m := toCompiled("cond",
		"condition", codeCompiled(c.condition),
		"consequent", codeCompiled(c.consequent),
		"ternary", codeCompiled(c.ternary),
	)

	if c.alternative != nil {
		m["alternative"] = codeCompiled(c.alternative)
	}

	return m
}

func switchCaseCompiled(c switchCase) map[string]interface{} {
	return toCompiled("switch-case", "expression", codeCompiled(c.expression), "body", codeCompiled(c.body))
}

func mapSwitchCasesCompiled(c []switchCase) []interface{} {
	var i []interface{}
	for _, ci := range c {
		i = append(i, ci)
	}

	return mapCompiled(i)
}

func switchStatementCompiled(s switchStatement) map[string]interface{} {
	m := toCompiled("switch-statement",
		"cases", mapSwitchCasesCompiled(s.cases),
		"defaultStatements", codeCompiled(s.defaultStatements),
	)

	if s.expression != nil {
		m["expression"] = codeCompiled(s.expression)
	}

	return m
}

func controlStatementCompiled(s controlStatement) map[string]interface{} {
	return toCompiled("control-statement", "control", int(s))
}

func rangeOverCompiled(r rangeOver) map[string]interface{} {
	m := toCompiled("range-over")

	if r.symbol != "" {
		m["symbol"] = codeCompiled(r.symbol)
	}

	if r.expression != nil {
		m["expression"] = codeCompiled(r.expression)
	}

	return m
}

func loopCompiled(l loop) map[string]interface{} {
	m := toCompiled("loop", "body", codeCompiled(l.body))

	if l.expression != nil {
		m["expression"] = codeCompiled(l.expression)
	}

	return m
}

func definitionCompiled(d definition) map[string]interface{} {
	return toCompiled("definition",
		"mutable", d.mutable,
		"exported", d.exported,
		"symbol", d.symbol,
		"expression", codeCompiled(d.expression),
	)
}

func mapDefinitionsCompiled(d []definition) []interface{} {
	var i []interface{}
	for _, di := range d {
		i = append(i, di)
	}

	return mapCompiled(i)
}

func definitionListCompiled(d definitionList) map[string]interface{} {
	return toCompiled("definition-list", "definitions", mapDefinitionsCompiled(d.definitions))
}

func assignCompiled(a assign) map[string]interface{} {
	return toCompiled("assign", "capture", codeCompiled(a.capture), "value", codeCompiled(a.value))
}

func mapAssignsCompiled(a []assign) []interface{} {
	var i []interface{}
	for _, ai := range a {
		i = append(i, ai)
	}

	return mapCompiled(i)
}

func assignListCompiled(a assignList) map[string]interface{} {
	return toCompiled("assign-list", "assignments", mapAssignsCompiled(a.assignments))
}

func sendCompiled(s send) map[string]interface{} {
	return toCompiled("send", "channel", codeCompiled(s.channel), "value", codeCompiled(s.value))
}

func receiveCompiled(r receive) map[string]interface{} {
	return toCompiled("receive", "channel", codeCompiled(r.channel))
}
func goCompiled(g goStatement) map[string]interface{} {
	return toCompiled("go", "application", codeCompiled(g.application))
}

func deferredCompiled(d deferred) map[string]interface{} {
	return toCompiled("deferred", "function", codeCompiled(d.function), "args", mapCompiled(d.args))
}

func deferCompiled(d deferStatement) map[string]interface{} {
	return toCompiled("defer", "application", codeCompiled(d.application))
}

func selectCaseCompiled(c selectCase) map[string]interface{} {
	return toCompiled("select-case", "expression", codeCompiled(c.expression), "body", codeCompiled(c.body))
}

func mapSelectCasesCompiled(c []selectCase) []interface{} {
	var i []interface{}
	for _, ci := range c {
		i = append(i, ci)
	}

	return mapCompiled(i)
}

func selectCompiled(s selectStatement) map[string]interface{} {
	return toCompiled("select",
		"cases", mapSelectCasesCompiled(s.cases),
		"hasDefault", s.hasDefault,
		"defaultStatements", codeCompiled(s.defaultStatements),
	)
}

func moduleCompiled(m module) map[string]interface{} {
	return toCompiled("module", "statements", mapCompiled(m.statements))
}

func useCompiled(u use) map[string]interface{} {
	return toCompiled("use", "path", codeCompiled(u.path), "capture", codeCompiled(u.capture))
}

func mapUsesCompiled(u []use) []interface{} {
	var i []interface{}
	for _, ui := range u {
		i = append(i, ui)
	}

	return mapCompiled(i)
}

func useListCompiled(u useList) map[string]interface{} {
	return toCompiled("use-list",
		"uses", mapUsesCompiled(u.uses),
	)
}

func codeCompiled(c interface{}) interface{} {
	switch ct := c.(type) {
	case comment:
		return commentCompiled(ct)
	case int, float64, string, bool:
		return c
	case primitive:
		return primitiveCompiled(ct)
	case symbol:
		return symbolCompiled(ct)
	case spread:
		return spreadCompiled(ct)
	case list:
		return listCompiled(ct)
	case expressionKey:
		return expressionKeyCompiled(ct)
	case entry:
		return entryCompiled(ct)
	case structure:
		return structureCompiled(ct)
	case ret:
		return retCompiled(ct)
	case statementList:
		return statementListCompiled(ct)
	case function:
		return functionCompiled(ct)
	case rangeExpression:
		return rangeExpressionCompiled(ct)
	case indexer:
		return indexerCompiled(ct)
	case functionApplication:
		return functionApplicationCompiled(ct)
	case unaryOperator:
		return unaryOperatorCompiled(ct)
	case unary:
		return unaryCompiled(ct)
	case binaryOperator:
		return binaryOperatorCompiled(ct)
	case binary:
		return binaryCompiled(ct)
	case cond:
		return condCompiled(ct)
	case switchCase:
		return switchCaseCompiled(ct)
	case switchStatement:
		return switchStatementCompiled(ct)
	case controlStatement:
		return controlStatementCompiled(ct)
	case rangeOver:
		return rangeOverCompiled(ct)
	case loop:
		return loopCompiled(ct)
	case definition:
		return definitionCompiled(ct)
	case definitionList:
		return definitionListCompiled(ct)
	case assign:
		return assignCompiled(ct)
	case assignList:
		return assignListCompiled(ct)
	case send:
		return sendCompiled(ct)
	case receive:
		return receiveCompiled(ct)
	case goStatement:
		return goCompiled(ct)
	case deferred:
		return deferredCompiled(ct)
	case deferStatement:
		return deferCompiled(ct)
	case selectCase:
		return selectCaseCompiled(ct)
	case selectStatement:
		return selectCompiled(ct)
	case module:
		return moduleCompiled(ct)
	case use:
		return useCompiled(ct)
	case useList:
		return useListCompiled(ct)
	default:
		panic(errUnsupportedCode)
	}
}
