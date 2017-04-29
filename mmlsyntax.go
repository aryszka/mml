package mml

func newMMLSyntax() (*syntax, error) {
	primitive := [][]interface{}{
		{"nl", nl},
		{"dot", dot},
		{"comma", comma},
		{"colon", colon},
		{"semicolon", semicolon},
		{"tilde", tilde},

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

		{"int", intToken},
		{"string", stringToken},
		{"symbol", symbolToken},
	}

	complex := [][]string{
		{"union", "seq-sep", "nl", "semicolon"},
		{"sequence", "nls", "nl"},

		{"union", "bool", "true", "false"},

		{
			"group",
			"dynamic-symbol",
			"symbol-word",
			"open-paren",
			"nls",
			"expression",
			"nls",
			"close-paren",
		},

		{"union", "static-symbol", "symbol", "string"},
		{"union", "symbol-expression", "static-symbol", "dynamic-symbol"},

		{"union", "list-sep", "nl", "comma"},
		{"group", "spread", "dot", "dot", "dot"},
		{"group", "spread-expression", "spread", "expression"},
		{"union", "list-item", "expression", "spread-expression", "list-sep"},
		{"sequence", "list-sequence", "list-item"},
		{"group", "list", "open-square", "list-sequence", "close-square"},
		{"group", "mutable-list", "tilde", "open-square", "list-sequence", "close-square"},

		{"group", "structure-definition", "symbol-expression", "nls", "colon", "nls", "expression"},
		{"union", "structure-item", "structure-definition", "spread-expression", "list-sep"},
		{"sequence", "structure-sequence", "structure-item"},
		{"group", "structure", "open-brace", "structure-sequence", "close-brace"},
		{"group", "mutable-structure", "tilde", "open-brace", "structure-sequence", "close-brace"},

		{"union", "static-symbol-item", "static-symbol", "list-sep"},
		{"sequence", "static-symbol-sequence", "static-symbol-item"},
		{"group", "collect-symbol", "spread", "static-symbol"},
		{"optional", "collect-argument", "collect-symbol"},
		{"group", "function-body", "open-brace", "statement-sequence", "close-brace"},
		{"union", "function-value", "expression", "function-body"},
		{
			"group",
			"function-fact",
			"open-paren",
			"static-symbol-sequence",
			"collect-argument",
			"nls",
			"close-paren",
			"nls",
			"function-value",
			// TODO: function-value could be simply an expression if there was a sequence as an
			// expression
		},

		{"group", "function", "fn-word", "function-fact"},

		{"group", "symbol-query", "expression", "dot", "symbol-expression"},
		{"union", "query", "symbol-query"},

		{"group", "function-call", "expression", "open-paren", "list-sequence", "close-paren"},

		{"union", "match-expression", "expression"},
		{"group", "switch-clause", "case-word", "match-expression", "colon", "statement-sequence"},
		{"sequence", "switch-clause-sequence", "switch-clause"},
		{"group", "default-clause", "default-word", "colon", "statement-sequence"},
		{
			"group",
			"switch-conditional",
			"switch-word",
			"nls",
			"open-brace",
			"nls",
			"switch-clause-sequence",
			"nls",
			"default-clause",
			"nls",
			"switch-clause-sequence",
			"nls",
			"close-brace",
		},
		{"union", "conditional", "switch-conditional"},

		{
			"union",
			"expression",
			"int",
			"string",
			"symbol",
			"bool",
			"dynamic-symbol",
			"function",
			"list",
			"mutable-list",
			"structure",
			"mutable-structure",
			"query",
			"function-call",
			"conditional",
		},

		{"union", "statement", "expression"},
		{"union", "statement-sequence-item", "statement", "seq-sep"},
		{"sequence", "statement-sequence", "statement-sequence-item"},
		{"union", "document", "statement-sequence"},
	}

	return defineSyntax(primitive, complex)
}
