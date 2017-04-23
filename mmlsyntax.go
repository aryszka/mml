package mml

func newMMLSyntax() (*syntax, error) {
	primitive := [][]interface{}{
		{"nl", nl},
		{"semicolon", semicolon},
		{"open-paren", openParen},
		{"close-paren", closeParen},
		{"int", intToken},
		{"string", stringToken},
		{"symbol", symbolToken},
		{"symbol-word", symbolWord},
	}

	complex := [][]string{
		{"union", "seq-sep", "nl", "semicolon"},
		{"sequence", "nls", "nl"},

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

		{"union", "expression", "int", "string", "symbol", "dynamic-symbol"},
		{"union", "statement", "expression"},
		{"union", "statement-sequence-item", "statement", "seq-sep"},
		{"sequence", "statement-sequence", "statement-sequence-item"},
		{"union", "document", "statement-sequence"},
	}

	s := newSyntax()

	for _, p := range primitive {
		if err := s.primitive(p[0].(string), p[1].(tokenType)); err != nil {
			return nil, err
		}
	}

	for _, c := range complex {
		var err error
		switch c[0] {
		case "optional":
			err = s.optional(c[1], c[2])
		case "sequence":
			err = s.sequence(c[1], c[2])
		case "group":
			err = s.group(c[1], c[2:]...)
		case "union":
			err = s.union(c[1], c[2:]...)
		default:
			panic("invalid parser type")
		}

		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
