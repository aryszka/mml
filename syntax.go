package mml

func init() {
	primitive("nl", nl)

	primitive("semicolon", semicolon)
	primitive("comma", comma)
	primitive("dot", dot)
	primitive("open-paren", openParen)
	primitive("close-paren", closeParen)
	primitive("open-brace", openBrace)
	primitive("close-brace", closeBrace)

	primitive("fn-word", fnWord)

	primitive("int", intToken)
	primitive("symbol", symbolToken)
	primitive("string", stringToken)

	sequence("nls", "nl")
	union("seq-sep", "nl", "semicolon")
	union("list-sep", "nl", "comma")
	group("spread", "dot", "dot", "dot")

	union("static-symbol", "symbol", "string")

	union("static-symbol-item", "static-symbol", "list-sep")
	sequence("static-symbol-sequence", "static-symbol-item")
	group("collect-symbol", "spread", "static-symbol")
	group("spread-expression", "expression", "spread")
	union("list-item", "expression", "spread-expression", "list-sep")
	sequence("list-sequence", "list-item")
	union("sequence-item", "statement", "seq-sep")
	sequence("statement-sequence", "sequence-item")

	optional("collect-argument", "collect-symbol")
	group("function-body", "open-brace", "statement-sequence", "close-brace")
	union("function-value", "expression", "function-body")
	group(
		"function-fact",
		"open-paren",
		"static-symbol-sequence",
		"collect-argument",
		"nls",
		"close-paren",
		"nls",
		"function-value",
	)
	group("function", "fn-word", "nls", "function-fact")

	group("function-call", "expression", "open-paren", "list-sequence", "close-paren")

	union(
		"expression",
		"int",
		"symbol",
		"function",
		"function-call",
	)

	// group("definition", "let", "nls", "static-symbol", "nls", "optional-eq", "nls", "expression")

	union(
		"statement",
		"expression",
	)

	union("document", "statement-sequence")

	isSep = func(n node) bool {
		switch n.typ {
		case "nl", "semicolon", "comma", "nls":
			return true
		default:
			return false
		}
	}

	setPostParse(map[string]func(node) node{
		"collect-symbol": func(n node) node {
			n.nodes = n.nodes[1:]
			return n
		},
		"spread-expression": func(n node) node {
			n.nodes = n.nodes[:1]
			return n
		},
		"function": func(n node) node {
			fact := n.nodes[1].nodes
			args := fact[1].nodes

			var value node
			if len(fact) == 5 {
				// when has varargs:
				args = append(args, fact[2])
				value = fact[4]
			} else {
				value = fact[3]
			}

			if value.typ == "function-body" {
				value = value.nodes[1]
			}

			n.nodes = append(args, value)
			return n
		},
		"function-call": func(n node) node {
			n.nodes = append(n.nodes[:1], n.nodes[2].nodes...)
			return n
		},
	})
}
