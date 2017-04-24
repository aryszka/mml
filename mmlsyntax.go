package mml

func newMMLSyntax() (*syntax, error) {
	primitive := [][]interface{}{
		{"nl", nl},
		{"dot", dot},
		{"comma", comma},
		{"semicolon", semicolon},
		{"open-paren", openParen},
		{"close-paren", closeParen},
		{"open-square", openSquare},
		{"close-square", closeSquare},
		{"int", intToken},
		{"string", stringToken},
		{"symbol", symbolToken},
		{"symbol-word", symbolWord},
		{"true", trueWord},
		{"false", falseWord},
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

		{"union", "list-sep", "nl", "comma"},
		{"group", "spread", "dot", "dot", "dot"},
		{"group", "spread-expression", "spread", "expression"},
		{"union", "list-item", "expression", "spread-expression", "list-sep"},
		{"sequence", "list-sequence", "list-item"},
		{"group", "list", "open-square", "list-sequence", "close-square"},
		{"group", "function-call", "expression", "open-paren", "list-sequence", "close-paren"},

		{
			"union",
			"expression",
			"int",
			"string",
			"symbol",
			"bool",
			"dynamic-symbol",
			"list",
			"function-call",
		},

		{"union", "statement", "expression"},
		{"union", "statement-sequence-item", "statement", "seq-sep"},
		{"sequence", "statement-sequence", "statement-sequence-item"},
		{"union", "document", "statement-sequence"},
	}

	return defineSyntax(primitive, complex)
}
