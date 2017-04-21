package mml

func init() {
	// primitive("nl", nl)

	// primitive("semicolon", semicolon)
	// primitive("comma", comma)
	// primitive("dot", dot)
	// primitive("tilde", tilde)
	// primitive("colon", colon)

	// primitive("open-paren", openParen)
	// primitive("close-paren", closeParen)
	// primitive("open-square", openSquare)
	// primitive("close-square", closeSquare)
	// primitive("open-brace", openBrace)
	// primitive("close-brace", closeBrace)

	// primitive("fn-word", fnWord)
	// primitive("symbol-word", symbolWord)
	// primitive("true", trueWord)
	// primitive("false", falseWord)
	// primitive("switch-word", switchWord)
	// primitive("case-word", caseWord)
	// primitive("default-word", defaultWord)
	// primitive("let-word", letWord)
	// primitive("if-word", ifWord)
	// primitive("else-word", elseWord)

	// primitive("single-eq", singleEq)

	// primitive("int", intToken)
	// primitive("symbol", symbolToken)
	// primitive("string", stringToken)

	// union("bool", "true", "false")

	// sequence("nls", "nl")
	// union("seq-sep", "nl", "semicolon")
	// union("list-sep", "nl", "comma")
	// group("spread", "dot", "dot", "dot")

	// union("static-symbol", "symbol", "string")
	// group("dynamic-symbol", "symbol-word", "open-paren", "nls", "expression", "nls", "close-paren")
	// union("symbol-expression", "static-symbol", "dynamic-symbol")

	// union("static-symbol-item", "static-symbol", "list-sep")
	// sequence("static-symbol-sequence", "static-symbol-item")
	// group("collect-symbol", "spread", "static-symbol")
	// group("spread-expression", "spread", "expression") // we can turn this around once having a single token for ...
	// union("list-item", "expression", "spread-expression", "list-sep")
	// sequence("list-sequence", "list-item")
	// union("sequence-item", "statement", "seq-sep")
	// sequence("statement-sequence", "sequence-item")
	// group("structure-definition", "symbol-expression", "nls", "colon", "nls", "expression")
	// union("structure-item", "structure-definition", "spread-expression", "list-sep")
	// sequence("structure-sequence", "structure-item")

	// group("list", "open-square", "list-sequence", "close-square")
	// group("mutable-list", "tilde", "list")

	// group("structure", "open-brace", "structure-sequence", "close-brace")
	// group("mutable-structure", "tilde", "structure")

	// optional("collect-argument", "collect-symbol")
	// group("function-body", "open-brace", "statement-sequence", "close-brace")
	// union("function-value", "expression", "function-body")
	// group(
	// 	"function-fact",
	// 	"open-paren",
	// 	"static-symbol-sequence",
	// 	"collect-argument",
	// 	"nls",
	// 	"close-paren",
	// 	"nls",
	// 	"function-value",
	// )
	// group("function", "fn-word", "nls", "function-fact")

	// group("effect-function", "fn-word", "nls", "tilde", "nls", "function-fact")

	// group("symbol-query", "expression", "nls", "dot", "nls", "symbol-expression")
	// optional("optional-expression", "expression")
	// group("range-expression", "optional-expression", "nls", "colon", "nls", "optional-expression")
	// union("query-expression", "expression", "range-expression")
	// group("expression-query", "expression", "open-square", "nls", "query-expression", "nls", "close-square")
	// union("query", "symbol-query", "expression-query")

	// group("function-call", "expression", "open-paren", "list-sequence", "close-paren")

	// union("match-expression", "expression")
	// group("switch-clause", "case-word", "match-expression", "colon", "statement-sequence")
	// group("default-clause", "default-word", "colon", "statement-sequence")
	// sequence("switch-clause-sequence", "switch-clause")
	// group(
	// 	"switch-conditional",
	// 	"switch-word",
	// 	"nls",
	// 	"open-brace",
	// 	"nls",
	// 	"switch-clause-sequence",
	// 	"nls",
	// 	"default-clause",
	// 	"nls",
	// 	"switch-clause-sequence",
	// 	"nls",
	// 	"close-brace",
	// )
	// group(
	// 	"if-conditional",
	// 	"if-word",
	// 	"nls",
	// 	"expression",
	// 	"nls",
	// 	"open-brace",
	// 	"nls",
	// 	"statement-sequence",
	// 	"nls",
	// 	"close-brace",
	// 	"nls",
	// 	"else-word",
	// 	"nls",
	// 	"open-brace",
	// 	"nls",
	// 	"statement-sequence",
	// 	"nls",
	// 	"close-brace",
	// )
	// union(
	// 	"conditional",
	// 	"switch-conditional",
	// 	"if-conditional",
	// )

	// union(
	// 	"expression",
	// 	"int",
	// 	"string",
	// 	"symbol",
	// 	"dynamic-symbol",
	// 	"bool",
	// 	"list",
	// 	"mutable-list",
	// 	"structure",
	// 	"mutable-structure",
	// 	"function",
	// 	"effect-function",
	// 	"query",
	// 	"function-call",
	// 	"conditional",
	// )

	// optional("optional-single-eq", "single-eq")
	// group("definition-item", "static-symbol", "nls", "optional-single-eq", "nls", "expression")
	// group("value-definition", "let-word", "nls", "definition-item")
	// group("mutable-value-definition", "let-word", "nls", "tilde", "nls", "definition-item")
	// union("definition-item-separator", "comma", "nl")
	// union("value-definition-sequence-item", "definition-item", "definition-item-separator")
	// sequence("value-definition-sequence", "value-definition-sequence-item")

	// group(
	// 	"value-definition-group",
	// 	"let-word",
	// 	"nls",
	// 	"open-paren",
	// 	"nls",
	// 	"value-definition-sequence",
	// 	"close-paren",
	// )

	// group(
	// 	"mutable-value-definition-group",
	// 	"let-word",
	// 	"nls",
	// 	"tilde",
	// 	"nls",
	// 	"open-paren",
	// 	"nls",
	// 	"value-definition-sequence",
	// 	"close-paren",
	// )

	// group("function-definition", "fn-word", "nls", "symbol-expression", "nls", "function-fact")
	// group(
	// 	"effect-function-definition",
	// 	"fn-word",
	// 	"nls",
	// 	"tilde",
	// 	"nls",
	// 	"symbol-expression",
	// 	"nls",
	// 	"function-fact",
	// )

	// union("definition",
	// 	"value-definition",
	// 	"mutable-value-definition",
	// 	"function-definition",
	// 	"effect-function-definition",
	// 	"value-definition-group",
	// 	"mutable-value-definition-group",
	// )

	// union(
	// 	"statement",
	// 	"expression",
	// 	"definition",
	// )

	// union("document", "statement-sequence")

	// cacheMembers()

	// isSep = func(n *node) bool {
	// 	switch n.typ {
	// 	case "nl", "semicolon", "comma", "nls":
	// 		return true
	// 	default:
	// 		return false
	// 	}
	// }

	// setPostParse(map[string]func(*node) *node{
	// 	"dynamic-symbol": func(n *node) *node {
	// 		n.nodes = n.nodes[2:3]
	// 		return n
	// 	},

	// 	"collect-symbol": func(n *node) *node {
	// 		n.nodes = n.nodes[1:]
	// 		return n
	// 	},

	// 	"spread-expression": func(n *node) *node {
	// 		n.nodes = n.nodes[1:]
	// 		return n
	// 	},

	// 	"list": func(n *node) *node {
	// 		n.nodes = n.nodes[1].nodes
	// 		return n
	// 	},

	// 	"mutable-list": func(n *node) *node {
	// 		n.nodes = n.nodes[1].nodes
	// 		return n
	// 	},

	// 	"structure-definition": func(n *node) *node {
	// 		n.nodes = append(n.nodes[:1], n.nodes[2])
	// 		return n
	// 	},

	// 	"structure": func(n *node) *node {
	// 		n.nodes = n.nodes[1].nodes
	// 		return n
	// 	},

	// 	"mutable-structure": func(n *node) *node {
	// 		n.nodes = n.nodes[1].nodes
	// 		return n
	// 	},

	// 	"function": func(n *node) *node {
	// 		fact := n.nodes[1].nodes
	// 		args := fact[1].nodes

	// 		var value *node
	// 		if len(fact) == 5 {
	// 			// when has varargs:
	// 			args = append(args, fact[2])
	// 			value = fact[4]
	// 		} else {
	// 			value = fact[3]
	// 		}

	// 		if value.typ == "function-body" {
	// 			value = value.nodes[1]
	// 		}

	// 		n.nodes = append(args, value)
	// 		return n
	// 	},

	// 	"effect-function": func(n *node) *node {
	// 		fact := n.nodes[2].nodes
	// 		args := fact[2].nodes

	// 		var value *node
	// 		if len(fact) == 5 {
	// 			// when has varargs:
	// 			args = append(args, fact[2])
	// 			value = fact[4]
	// 		} else {
	// 			value = fact[3]
	// 		}

	// 		if value.typ == "function-body" {
	// 			value = value.nodes[1]
	// 		}

	// 		n.nodes = append(args, value)
	// 		return n
	// 	},

	// 	"symbol-query": func(n *node) *node {
	// 		n.nodes = append(n.nodes[:1], n.nodes[2])
	// 		return n
	// 	},

	// 	"range-expression": func(n *node) *node {
	// 		if len(n.nodes) == 1 {
	// 			n.nodes = make([]*node, 2)
	// 			return n
	// 		}

	// 		if n.nodes[0].typ == "colon" {
	// 			n.nodes = []*node{{}, n.nodes[1]}
	// 			return n
	// 		}

	// 		n.nodes = append(n.nodes[:1], n.nodes[2:]...)
	// 		return n
	// 	},

	// 	"expression-query": func(n *node) *node {
	// 		n.nodes = append(n.nodes[:1], n.nodes[2])
	// 		return n
	// 	},

	// 	"function-call": func(n *node) *node {
	// 		n.nodes = append(n.nodes[:1], n.nodes[2].nodes...)
	// 		return n
	// 	},

	// 	"switch-clause": func(n *node) *node {
	// 		n.nodes = append(n.nodes[1:2], n.nodes[3].nodes...)
	// 		return n
	// 	},

	// 	"default-clause": func(n *node) *node {
	// 		n.nodes = n.nodes[2].nodes
	// 		return n
	// 	},

	// 	"switch-conditional": func(n *node) *node {
	// 		n.nodes = n.nodes[2 : len(n.nodes)-1]

	// 		var nodes []*node
	// 		for _, ni := range n.nodes {
	// 			switch ni.typ {
	// 			case "switch-clause-sequence":
	// 				if len(ni.nodes) > 0 {
	// 					nodes = append(nodes, ni.nodes...)
	// 				}
	// 			case "default-clause":
	// 				nodes = append(nodes, ni)
	// 			}
	// 		}

	// 		n.nodes = nodes
	// 		return n
	// 	},

	// 	"if-conditional": func(n *node) *node {
	// 		return n
	// 	},

	// 	"definition-item": func(n *node) *node {
	// 		if len(n.nodes) == 2 {
	// 			return n
	// 		}

	// 		n.nodes = append(n.nodes[:1], n.nodes[2])
	// 		return n
	// 	},

	// 	"value-definition": func(n *node) *node {
	// 		n.nodes = n.nodes[1].nodes
	// 		return n
	// 	},

	// 	"mutable-value-definition": func(n *node) *node {
	// 		n.nodes = n.nodes[2].nodes
	// 		return n
	// 	},

	// 	"function-definition": func(n *node) *node {
	// 		n.nodes = n.nodes[1:]
	// 		return n
	// 	},

	// 	"effect-function-definition": func(n *node) *node {
	// 		n.nodes = n.nodes[2:]
	// 		return n
	// 	},
	// })
}
