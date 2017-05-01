package mml

import (
	"bytes"
	"testing"
)

func checkTokens(left, right *token) bool {
	if (left == nil) != (right == nil) {
		return false
	}

	if left == nil {
		return true
	}

	return left.value == right.value
}

func checkNodes(left, right *node) bool {
	if (left == nil) != (right == nil) {
		return false
	}

	if left == nil {
		return true
	}

	if left.typ != right.typ {
		return false
	}

	if !checkTokens(left.token, right.token) {
		return false
	}

	if len(left.nodes) != len(right.nodes) {
		return false
	}

	for i, n := range left.nodes {
		if !checkNodes(n, right.nodes[i]) {
			return false
		}
	}

	return true
}

func def(f ...func(s *syntax) error) func(s *syntax) error {
	return func(s *syntax) error {
		for _, fi := range f {
			if err := fi(s); err != nil {
				return err
			}
		}

		return nil
	}
}

func TestParse(t *testing.T) {
	for _, ti := range []struct {
		msg       string
		primitive [][]interface{}
		complex   [][]string
		text      string
		node      *node
		fail      bool
	}{{
		msg:       "int",
		primitive: [][]interface{}{{"int", intToken}},
		text:      "42",
		node: &node{
			name: "int",
			token:    &token{value: "42"},
		},
	}, {
		msg:       "int, with empty input",
		primitive: [][]interface{}{{"int", intToken}},
		fail:      true,
	}, {
		msg:       "optional int",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"optional", "optional-int", "int"}},
		text:      "42",
		node: &node{
			name: "int",
			token:    &token{value: "42"},
		},
	}, {
		msg:       "optional int, empty",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"optional", "optional-int", "int"}},
	}, {
		msg:       "optional int, not int",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"optional", "optional-int", "int"}},
		text:      "\"foo\"",
		fail:      true,
	// }, {
	// 	msg:       "int sequence, optional",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"sequence", "int-sequence", "int"},
	// 		{"optional", "optional-int-sequence", "int-sequence"},
	// 	},
	// 	text: "1 2 3",
	// 	node: &node{
	// 		name: "int-sequence",
	// 		token:    &token{value: "1"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "1"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "2"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "3"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "int sequence, optional, empty",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"sequence", "int-sequence", "int"},
	// 		{"optional", "optional-int-sequence", "int-sequence"},
	// 	},
	// 	node: &node{
	// 		name: "int-sequence",
	// 		token:    &token{},
	// 	},
	// }, {
	// 	msg:       "empty sequence",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex:   [][]string{{"sequence", "int-sequence", "int"}},
	// 	node: &node{
	// 		name: "int-sequence",
	// 		token:    &token{},
	// 	},
	// }, {
	// 	msg:       "sequence with a single item",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex:   [][]string{{"sequence", "int-sequence", "int"}},
	// 	text:      "42",
	// 	node: &node{
	// 		name: "int-sequence",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "42"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "sequence with multiple items",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex:   [][]string{{"sequence", "int-sequence", "int"}},
	// 	text:      "1 2 3",
	// 	node: &node{
	// 		name: "int-sequence",
	// 		token:    &token{value: "1"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "1"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "2"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "3"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "sequence with optional item",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{"sequence", "optional-int-sequence", "optional-int"},
	// 	},
	// 	text: "42",
	// 	node: &node{
	// 		name: "optional-int-sequence",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "42"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "sequence with multiple optional items",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{"sequence", "optional-int-sequence", "optional-int"},
	// 	},
	// 	text: "1 2 3",
	// 	node: &node{
	// 		name: "optional-int-sequence",
	// 		token:    &token{value: "1"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "1"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "2"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "3"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "group with single int",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex:   [][]string{{"group", "int-group", "int"}},
	// 	text:      "42",
	// 	node: &node{
	// 		name: "int-group",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "42"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "group with single optional int",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{"group", "optional-int-group", "optional-int"},
	// 	},
	// 	text: "42",
	// 	node: &node{
	// 		name: "optional-int-group",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "42"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "group with single int, not int",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex:   [][]string{{"group", "int-group", "int"}},
	// 	text:      "\"foo\"",
	// 	fail:      true,
	// }, {
	// 	msg:       "group with multiple ints",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex:   [][]string{{"group", "int-group", "int", "int", "int"}},
	// 	text:      "1 2 3",
	// 	node: &node{
	// 		name: "int-group",
	// 		token:    &token{value: "1"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "1"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "2"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "3"},
	// 		}},
	// 	},
	// }, {
	// 	msg: "group with optional item",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 		{"string", stringToken},
	// 	},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{"group", "group-with-optional", "optional-int", "string"},
	// 	},
	// 	text: "42 \"foo\"",
	// 	node: &node{
	// 		name: "group-with-optional",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "42"},
	// 		}, {
	// 			name: "string",
	// 			token:    &token{value: "\"foo\""},
	// 		}},
	// 	},
	// }, {
	// 	msg: "group with optional item, missing",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 		{"string", stringToken},
	// 	},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{"group", "group-with-optional", "optional-int", "string"},
	// 	},
	// 	text: "\"foo\"",
	// 	node: &node{
	// 		name: "group-with-optional",
	// 		token:    &token{value: "\"foo\""},
	// 		nodes: []*node{{
	// 			name: "string",
	// 			token:    &token{value: "\"foo\""},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "group with only optional, empty",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{
	// 			"group",
	// 			"group-with-only-optional",
	// 			"optional-int",
	// 			"optional-int",
	// 			"optional-int",
	// 		},
	// 	},
	// 	node: &node{
	// 		name: "group-with-only-optional",
	// 	},
	// }, {
	// 	msg:       "group with only optional, less",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{
	// 			"group",
	// 			"group-with-only-optional",
	// 			"optional-int",
	// 			"optional-int",
	// 			"optional-int",
	// 		},
	// 	},
	// 	text: "1 2",
	// 	node: &node{
	// 		name: "group-with-only-optional",
	// 		token:    &token{value: "1"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "1"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "2"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "group with only optional, exact",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{
	// 			"group",
	// 			"group-with-only-optional",
	// 			"optional-int",
	// 			"optional-int",
	// 			"optional-int",
	// 		},
	// 	},
	// 	text: "1 2 3",
	// 	node: &node{
	// 		name: "group-with-only-optional",
	// 		token:    &token{value: "1"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "1"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "2"},
	// 		}, {
	// 			name: "int",
	// 			token:    &token{value: "3"},
	// 		}},
	// 	},
	// }, {
	// 	msg:       "group with only optional, more",
	// 	primitive: [][]interface{}{{"int", intToken}},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{
	// 			"group",
	// 			"group-with-only-optional",
	// 			"optional-int",
	// 			"optional-int",
	// 			"optional-int",
	// 		},
	// 	},
	// 	text: "1 2 3 4",
	// 	fail: true,
	// }, {
	// 	msg: "union of int and string",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 		{"string", stringToken},
	// 	},
	// 	complex: [][]string{
	// 		{"union", "int-or-string", "int", "string"},
	// 	},
	// 	text: "\"foo\"",
	// 	node: &node{
	// 		name: "string",
	// 		token:    &token{value: "\"foo\""},
	// 	},
	// }, {
	// 	msg: "union of int and group with optional int",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 		{"string", stringToken},
	// 	},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{"group", "group-with-optional", "optional-int", "string"},
	// 		{"union", "int-or-group-with-optional", "int", "group-with-optional"},
	// 	},
	// 	text: "42 \"foo\"",
	// 	node: &node{
	// 		name: "group-with-optional",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int",
	// 			token:    &token{value: "42"},
	// 		}, {
	// 			name: "string",
	// 			token:    &token{value: "\"foo\""},
	// 		}},
	// 	},
	// }, {
	// 	msg: "union of int and group with optional int, token fall through",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 		{"string", stringToken},
	// 	},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{
	// 			"group",
	// 			"group-with-optional",
	// 			"optional-int",
	// 			"optional-int",
	// 			"string",
	// 			"string",
	// 		},
	// 		{"union", "int-or-group-with-optional", "int", "group-with-optional"},
	// 	},
	// 	text: "\"foo\" \"bar\"",
	// 	node: &node{
	// 		name: "group-with-optional",
	// 		token:    &token{value: "\"foo\""},
	// 		nodes: []*node{{
	// 			name: "string",
	// 			token:    &token{value: "\"foo\""},
	// 		}, {
	// 			name: "string",
	// 			token:    &token{value: "\"bar\""},
	// 		}},
	// 	},
	// }, {
	// 	msg: "union of int and group with optional int, init fall through",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 		{"string", stringToken},
	// 	},
	// 	complex: [][]string{
	// 		{"optional", "optional-int", "int"},
	// 		{
	// 			"group",
	// 			"group-with-optional",
	// 			"optional-int",
	// 			"optional-int",
	// 			"string",
	// 			"string",
	// 		},
	// 		{"union", "int-or-group-with-optional", "int", "group-with-optional"},
	// 	},
	// 	text: "\"foo\" \"bar\"",
	// 	node: &node{
	// 		name: "group-with-optional",
	// 		token:    &token{value: "\"foo\""},
	// 		nodes: []*node{{
	// 			name: "string",
	// 			token:    &token{value: "\"foo\""},
	// 		}, {
	// 			name: "string",
	// 			token:    &token{value: "\"bar\""},
	// 		}},
	// 	},
	// }, {
	// 	msg: "expression inside expression",
	// 	primitive: [][]interface{}{
	// 		{"symbol", symbolToken},
	// 		{"symbol-word", symbolWord},
	// 		{"open-paren", openParen},
	// 		{"close-paren", closeParen},
	// 	},
	// 	complex: [][]string{
	// 		{"group", "function-call", "expression", "open-paren", "expression", "close-paren"},
	// 		{"group", "dynamic-symbol", "symbol-word", "open-paren", "expression", "close-paren"},
	// 		{"union", "expression", "symbol", "function-call", "dynamic-symbol"},
	// 	},
	// 	text: "symbol(f(a))",
	// 	node: &node{
	// 		name: "dynamic-symbol",
	// 		token:    &token{value: "symbol"},
	// 		nodes: []*node{{
	// 			name: "symbol-word",
	// 			token:    &token{value: "symbol"},
	// 		}, {
	// 			name: "open-paren",
	// 			token:    &token{value: "("},
	// 		}, {
	// 			name: "function-call",
	// 			token:    &token{value: "f"},
	// 			nodes: []*node{{
	// 				name: "symbol",
	// 				token:    &token{value: "f"},
	// 			}, {
	// 				name: "open-paren",
	// 				token:    &token{value: "("},
	// 			}, {
	// 				name: "symbol",
	// 				token:    &token{value: "a"},
	// 			}, {
	// 				name: "close-paren",
	// 				token:    &token{value: ")"},
	// 			}},
	// 		}, {
	// 			name: "close-paren",
	// 			token:    &token{value: ")"},
	// 		}},
	// 	},
	// }, {
	// 	msg: "chained symbol query",
	// 	primitive: [][]interface{}{
	// 		{"symbol", symbolToken},
	// 		{"dot", dot},
	// 	},
	// 	complex: [][]string{
	// 		{"group", "symbol-query", "expression", "dot", "symbol"},
	// 		{"union", "expression", "symbol", "symbol-query"},
	// 	},
	// 	text: "a.b.c",
	// 	node: &node{
	// 		name: "symbol-query",
	// 		token:    &token{value: "a"},
	// 		nodes: []*node{{
	// 			name: "symbol-query",
	// 			token:    &token{value: "a"},
	// 			nodes: []*node{{
	// 				name: "symbol",
	// 				token:    &token{value: "a"},
	// 			}, {
	// 				name: "dot",
	// 				token:    &token{value: "."},
	// 			}, {
	// 				name: "symbol",
	// 				token:    &token{value: "b"},
	// 			}},
	// 		}, {
	// 			name: "dot",
	// 			token:    &token{value: "."},
	// 		}, {
	// 			name: "symbol",
	// 			token:    &token{value: "c"},
	// 		}},
	// 	},
	// }, {
	// 	msg: "sequence in sequence",
	// 	primitive: [][]interface{}{
	// 		{"int", intToken},
	// 	},
	// 	complex: [][]string{
	// 		{"sequence", "int-sequence", "int"},
	// 		{"sequence", "sequence-in-sequence", "int-sequence"},
	// 	},
	// 	text: "42",
	// 	node: &node{
	// 		name: "sequence-in-sequence",
	// 		token:    &token{value: "42"},
	// 		nodes: []*node{{
	// 			name: "int-sequence",
	// 			token:    &token{value: "42"},
	// 			nodes: []*node{{
	// 				name: "int",
	// 				token:    &token{value: "42"},
	// 			}},
	// 		}},
	// 	},
	// }, {
	// 	msg: "reproduce sequence endless loop",
	// 	primitive: [][]interface{}{
	// 		{"nl", nl},
	// 		{"colon", colon},
	// 		{"switch-word", switchWord},
	// 		{"case-word", caseWord},
	// 		{"default-word", defaultWord},
	// 		{"open-brace", openBrace},
	// 		{"close-brace", closeBrace},
	// 		{"symbol", symbolToken},
	// 	},
	// 	complex: [][]string{
	// 		{"sequence", "nls", "nl"},
	// 		{"union", "match-expression", "expression"},
	// 		{"group", "switch-clause", "case-word", "match-expression", "colon", "statement-sequence"},
	// 		{"sequence", "switch-clause-sequence", "switch-clause"},
	// 		{"group", "default-clause", "default-word", "colon", "nls", "statement-sequence"},
	// 		{"union", "seq-sep", "nl"},
	// 		{"union", "statement-sequence-item", "expression", "seq-sep"},
	// 		{"sequence", "statement-sequence", "statement-sequence-item"},
	// 		{
	// 			"group",
	// 			"switch-conditional",
	// 			"switch-word",
	// 			"nls",
	// 			"open-brace",
	// 			"nls",
	// 			"switch-clause-sequence",
	// 			"nls",
	// 			"default-clause",
	// 			"nls",
	// 			"switch-clause-sequence",
	// 			"nls",
	// 			"close-brace",
	// 		},
	// 		{"union", "expression", "symbol", "switch-conditional"},
	// 		{"sequence", "document", "statement-sequence"},
	// 	},
	// 	text: `switch {
	// 		default: a
	// 	}`,
	// 	node: &node{
	// 		name: "document",
	// 		token:    &token{value: "switch"},
	// 		nodes: []*node{{
	// 			name: "statement-sequence",
	// 			token:    &token{value: "switch"},
	// 			nodes: []*node{{
	// 				name: "switch-conditional",
	// 				token:    &token{value: "switch"},
	// 				nodes: []*node{{
	// 					name: "switch-word",
	// 					token:    &token{value: "switch"},
	// 				}, {
	// 					name: "nls",
	// 					token:    &token{value: "{"},
	// 				}, {
	// 					name: "open-brace",
	// 					token:    &token{value: "{"},
	// 				}, {
	// 					name: "nls",
	// 					token:    &token{value: "\n"},
	// 					nodes: []*node{{
	// 						name: "nl",
	// 						token:    &token{value: "\n"},
	// 					}},
	// 				}, {
	// 					name: "switch-clause-sequence",
	// 					token:    &token{value: "default"},
	// 				}, {
	// 					name: "nls",
	// 					token:    &token{value: "default"},
	// 				}, {
	// 					name: "default-clause",
	// 					token:    &token{value: "default"},
	// 					nodes: []*node{{
	// 						name: "default-word",
	// 						token:    &token{value: "default"},
	// 					}, {
	// 						name: "colon",
	// 						token:    &token{value: ":"},
	// 					}, {
	// 						name: "nls",
	// 						token:    &token{value: "a"},
	// 					}, {
	// 						name: "statement-sequence",
	// 						token:    &token{value: "a"},
	// 						nodes: []*node{{
	// 							name: "symbol",
	// 							token:    &token{value: "a"},
	// 						}, {
	// 							name: "nl",
	// 							token:    &token{value: "\n"},
	// 						}},
	// 					}},
	// 				}, {
	// 					name: "nls",
	// 					token:    &token{value: "}"},
	// 				}, {
	// 					name: "switch-clause-sequence",
	// 					token:    &token{value: "}"},
	// 				}, {
	// 					name: "nls",
	// 					token:    &token{value: "}"},
	// 				}, {
	// 					name: "close-brace",
	// 					token:    &token{value: "}"},
	// 				}},
	// 			}},
	// 		}},
	// 	},
	// }, {
	// 	msg: "newline in group",
	// 	primitive: [][]interface{}{
	// 		{"nl", nl},
	// 		{"colon", colon},
	// 		{"switch-word", switchWord},
	// 		{"case-word", caseWord},
	// 		{"default-word", defaultWord},
	// 		{"open-brace", openBrace},
	// 		{"close-brace", closeBrace},
	// 		{"symbol", symbolToken},
	// 	},
	// 	complex: [][]string{
	// 		{"sequence", "nls", "nl"},
	// 		{"union", "match-expression", "expression"},
	// 		{"group", "switch-clause", "case-word", "match-expression", "colon", "statement-sequence"},
	// 		{"sequence", "switch-clause-sequence", "switch-clause"},
	// 		{"group", "default-clause", "default-word", "colon", "nls", "statement-sequence"},
	// 		{"union", "seq-sep", "nl"},
	// 		{"union", "statement-sequence-item", "statement", "seq-sep"},
	// 		{"sequence", "statement-sequence", "statement-sequence-item"},
	// 		{
	// 			"group",
	// 			"switch-conditional",
	// 			"switch-word",
	// 			"nls",
	// 			"open-brace",
	// 			"nls",
	// 			"switch-clause-sequence",
	// 			"nls",
	// 			"default-clause",
	// 			"nls",
	// 			"switch-clause-sequence",
	// 			"nls",
	// 			"close-brace",
	// 		},
	// 		{"union", "conditional", "switch-conditional"},
	// 		{"union", "expression", "symbol", "conditional"},
	// 		{"union", "statement", "expression"},
	// 		{"union", "document", "statement-sequence"},
	// 	},
	// 	text: `switch {
	// 		default: a
	// 	}`,
	// 	node: &node{
	// 		name: "statement-sequence",
	// 		token:    &token{value: "switch"},
	// 		nodes: []*node{{
	// 			name: "switch-conditional",
	// 			token:    &token{value: "switch"},
	// 			nodes: []*node{{
	// 				name: "switch-word",
	// 				token:    &token{value: "switch"},
	// 			}, {
	// 				name: "nls",
	// 				token:    &token{value: "{"},
	// 			}, {
	// 				name: "open-brace",
	// 				token:    &token{value: "{"},
	// 			}, {
	// 				name: "nls",
	// 				token:    &token{value: "\n"},
	// 				nodes: []*node{{
	// 					name: "nl",
	// 					token:    &token{value: "\n"},
	// 				}},
	// 			}, {
	// 				name: "switch-clause-sequence",
	// 				token:    &token{value: "default"},
	// 			}, {
	// 				name: "nls",
	// 				token:    &token{value: "default"},
	// 			}, {
	// 				name: "default-clause",
	// 				token:    &token{value: "default"},
	// 				nodes: []*node{{
	// 					name: "default-word",
	// 					token:    &token{value: "default"},
	// 				}, {
	// 					name: "colon",
	// 					token:    &token{value: ":"},
	// 				}, {
	// 					name: "nls",
	// 					token:    &token{value: "a"},
	// 				}, {
	// 					name: "statement-sequence",
	// 					token:    &token{value: "a"},
	// 					nodes: []*node{{
	// 						name: "symbol",
	// 						token:    &token{value: "a"},
	// 					}, {
	// 						name: "nl",
	// 						token:    &token{value: "\n"},
	// 					}},
	// 				}},
	// 			}, {
	// 				name: "nls",
	// 				token:    &token{value: "}"},
	// 			}, {
	// 				name: "switch-clause-sequence",
	// 				token:    &token{value: "}"},
	// 			}, {
	// 				name: "nls",
	// 				token:    &token{value: "}"},
	// 			}, {
	// 				name: "close-brace",
	// 				token:    &token{value: "}"},
	// 			}},
	// 		}},
	// 	},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			var l traceLevel
			trace := newTrace(l)
			s := withTrace(trace)

			err := s.defineSyntax(ti.primitive, ti.complex)
			if err != nil {
				t.Error(err)
				return
			}

			// s.traceLevel = traceDebug

			b := bytes.NewBufferString(ti.text)

			n, err := s.parse(b, "test")
			if !ti.fail && err != nil {
				t.Error(err)
				return
			} else if ti.fail && err == nil {
				t.Error("failed to fail")
				return
			}

			if ti.fail {
				return
			}

			if !checkNodes(n, ti.node) {
				t.Error("failed to match nodes")
				t.Log(n)
				t.Log(ti.node)
			}
		})
	}
}
