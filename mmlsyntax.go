package mml

func (s *syntax) newMMLSyntax() error {
	primitive := [][]interface{}{
		{"nl", nl},
		{"dot", dot},
		{"comma", comma},
		{"colon", colon},
		{"semicolon", semicolon},
		{"tilde", tilde},
		{"single-eq", singleEq},

		{"open-paren", openParen},
		{"close-paren", closeParen},
		{"open-square", openSquare},
		{"close-square", closeSquare},
		{"open-brace", openBrace},
		{"close-brace", closeBrace},

		{"symbol-word", symbolWord},
		{"true", trueWord},
		{"false", falseWord},
		{"fn-word", fnWord},
		{"switch-word", switchWord},
		{"case-word", caseWord},
		{"default-word", defaultWord},
		{"let-word", letWord},
		{"if-word", ifWord},
		{"else-word", elseWord},
		{"set-word", setWord},

		{"int", intToken},
		{"string", stringToken},
		{"symbol", symbolToken},
	}

	complex := [][]string{
		{"choice", "seq-sep", "nl", "semicolon"},
		{"repeat", "nls", "nl"},

		{"choice", "bool", "true", "false"},

		{
			"sequence",
			"dynamic-symbol",
			"symbol-word",
			"open-paren",
			"nls",
			"expression",
			"nls",
			"close-paren",
		},

		{"choice", "static-symbol", "symbol", "string"},
		{"choice", "symbol-expression", "static-symbol", "dynamic-symbol"},

		{"choice", "list-sep", "nl", "comma"},
		{"sequence", "spread", "dot", "dot", "dot"},
		{"sequence", "spread-expression", "spread", "expression"},
		{"choice", "list-item", "expression", "spread-expression", "list-sep"},
		{"repeat", "list-repeat", "list-item"},
		{"sequence", "list", "open-square", "list-repeat", "close-square"},
		{"sequence", "mutable-list", "tilde", "open-square", "list-repeat", "close-square"},

		{"sequence", "structure-definition", "symbol-expression", "nls", "colon", "nls", "expression"},
		{"choice", "structure-item", "structure-definition", "spread-expression", "list-sep"},
		{"repeat", "structure-repeat", "structure-item"},
		{"sequence", "structure", "open-brace", "structure-repeat", "close-brace"},
		{"sequence", "mutable-structure", "tilde", "open-brace", "structure-repeat", "close-brace"},

		{"choice", "static-symbol-item", "static-symbol", "list-sep"},
		{"repeat", "static-symbol-repeat", "static-symbol-item"},
		{"sequence", "collect-symbol", "spread", "static-symbol"},
		{"optional", "collect-argument", "collect-symbol"},
		{"sequence", "function-body", "open-brace", "statement-repeat", "close-brace"},
		{"choice", "function-value", "expression", "function-body"},
		{
			"sequence",
			"function-fact",
			"open-paren",
			"static-symbol-repeat",
			"collect-argument",
			"nls",
			"close-paren",
			"nls",
			"function-value",
			// TODO: function-value could be simply an expression if there was a repeat as an
			// expression
		},

		{"sequence", "function", "fn-word", "function-fact"},
		{"sequence", "effect", "fn-word", "tilde", "function-fact"},

		{"sequence", "symbol-query", "expression", "dot", "symbol-expression"},
		{"optional", "optional-expression", "expression"},
		{
			"sequence",
			"range-expression",
			"optional-expression",
			"nls",
			"colon",
			"nls",
			"optional-expression",
		},
		{"choice", "query-expression", "expression", "range-expression"},
		{
			"sequence",
			"expression-query",
			"expression",
			"open-square",
			"nls",
			"query-expression",
			"nls",
			"close-square",
		},
		{"choice", "query", "symbol-query", "expression-query"},

		{"sequence", "function-call", "expression", "open-paren", "list-repeat", "close-paren"},

		{"choice", "match-expression", "expression"},
		{"sequence", "switch-clause", "case-word", "match-expression", "colon", "statement-repeat"},
		{"repeat", "switch-clause-repeat", "switch-clause"},
		{"sequence", "default-clause", "default-word", "colon", "statement-repeat"},
		{
			"sequence",
			"switch-conditional",
			"switch-word",
			"nls",
			"open-brace",
			"nls",
			"switch-clause-repeat",
			"nls",
			"default-clause",
			"nls",
			"switch-clause-repeat",
			"nls",
			"close-brace",
		},
		{
			"sequence",
			"if-conditional", // TODO: test
			"if-word",
			"nls",
			"match-expression",
			"nls",
			"open-brace",
			"nls",
			"statement-repeat",
			"nls",
			"close-brace",
			"nls",
			"else-word",
			"nls",
			"open-brace",
			"nls",
			"statement-repeat",
			"nls",
			"close-brace",
		},
		{"choice", "conditional", "switch-conditional", "if-conditional"},

		{
			"choice",
			"expression",
			"int",
			"string",
			"symbol",
			"bool",
			"dynamic-symbol",
			"function",
			"effect",
			"list",
			"mutable-list",
			"structure",
			"mutable-structure",
			"query",
			"function-call",
			"conditional",
		},

		{"optional", "optional-single-eq", "single-eq"},

		{
			"sequence",
			"definition-item",
			"symbol-expression",
			"nls",
			"optional-single-eq",
			"nls",
			"expression",
		},

		{
			"sequence",
			"value-definition",
			"let-word",
			"nls",
			"definition-item",
		},

		{
			"sequence",
			"mutable-value-definition",
			"let-word",
			"nls",
			"tilde",
			"nls",
			"definition-item",
		},

		{
			"sequence",
			"value-assignment", // TODO: test
			"set-word",
			"nls",
			"definition-item",
		},

		{"choice", "value-definition-repeat-item", "definition-item", "list-sep"},
		{"repeat", "value-definition-repeat", "value-definition-repeat-item"},

		{
			"sequence",
			"value-definition-sequence",
			"let-word",
			"nls",
			"open-paren",
			"value-definition-repeat",
			"close-paren",
		},

		{
			"sequence",
			"mutable-value-definition-sequence",
			"let-word",
			"nls",
			"tilde",
			"nls",
			"open-paren",
			"value-definition-repeat",
			"close-paren",
		},

		{"sequence", "function-definition", "fn-word", "nls", "symbol-expression", "nls", "function-fact"},
		{"sequence", "effect-definition", "fn-word", "nls", "tilde", "nls", "symbol-expression", "nls", "function-fact"},

		{
			"choice",
			"definition",
			"value-definition",                  // TODO: test
			"mutable-value-definition",          // TODO: test
			"value-definition-sequence",         // TODO: test
			"mutable-value-definition-sequence", // TODO: test
			"function-definition",               // TODO: test
			"effect-definition",                 // TODO: test
		},

		{"choice", "statement", "expression", "definition", "value-assignment"},
		{"choice", "statement-repeat-item", "statement", "seq-sep"},
		{"repeat", "statement-repeat", "statement-repeat-item"},
		{"choice", "document", "statement-repeat"},
	}

	return s.defineSyntax(primitive, complex)
}
