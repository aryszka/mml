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

	if left.name != right.name {
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
			name:  "int",
			token: &token{value: "42"},
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
			name:  "int",
			token: &token{value: "42"},
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
	}, {
		msg:       "int repetition, optional",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"repetition", "int-repetition", "int"},
			{"optional", "optional-int-repetition", "int-repetition"},
		},
		text: "1 2 3",
		node: &node{
			name:  "int-repetition",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "1"},
			}, {
				name:  "int",
				token: &token{value: "2"},
			}, {
				name:  "int",
				token: &token{value: "3"},
			}},
		},
	}, {
		msg:       "int repetition, optional, empty",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"repetition", "int-repetition", "int"},
			{"optional", "optional-int-repetition", "int-repetition"},
		},
		node: &node{
			name:  "int-repetition",
			token: &token{},
		},
	}, {
		msg:       "empty repetition",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"repetition", "int-repetition", "int"}},
		node: &node{
			name:  "int-repetition",
			token: &token{},
		},
	}, {
		msg:       "repetition with a single item",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"repetition", "int-repetition", "int"}},
		text:      "42",
		node: &node{
			name:  "int-repetition",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "42"},
			}},
		},
	}, {
		msg:       "repetition with multiple items",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"repetition", "int-repetition", "int"}},
		text:      "1 2 3",
		node: &node{
			name:  "int-repetition",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "1"},
			}, {
				name:  "int",
				token: &token{value: "2"},
			}, {
				name:  "int",
				token: &token{value: "3"},
			}},
		},
	}, {
		msg:       "repetition with optional item",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{"repetition", "optional-int-repetition", "optional-int"},
		},
		text: "42",
		node: &node{
			name:  "optional-int-repetition",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "42"},
			}},
		},
	}, {
		msg:       "repetition with multiple optional items",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{"repetition", "optional-int-repetition", "optional-int"},
		},
		text: "1 2 3",
		node: &node{
			name:  "optional-int-repetition",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "1"},
			}, {
				name:  "int",
				token: &token{value: "2"},
			}, {
				name:  "int",
				token: &token{value: "3"},
			}},
		},
	}, {
		msg:       "sequence with single int",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"sequence", "int-sequence", "int"}},
		text:      "42",
		node: &node{
			name:  "int-sequence",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "42"},
			}},
		},
	}, {
		msg:       "sequence with single optional int",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{"sequence", "optional-int-sequence", "optional-int"},
		},
		text: "42",
		node: &node{
			name:  "optional-int-sequence",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "42"},
			}},
		},
	}, {
		msg:       "sequence with single int, not int",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"sequence", "int-sequence", "int"}},
		text:      "\"foo\"",
		fail:      true,
	}, {
		msg:       "sequence with multiple ints",
		primitive: [][]interface{}{{"int", intToken}},
		complex:   [][]string{{"sequence", "int-sequence", "int", "int", "int"}},
		text:      "1 2 3",
		node: &node{
			name:  "int-sequence",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "1"},
			}, {
				name:  "int",
				token: &token{value: "2"},
			}, {
				name:  "int",
				token: &token{value: "3"},
			}},
		},
	}, {
		msg: "sequence with optional item",
		primitive: [][]interface{}{
			{"int", intToken},
			{"string", stringToken},
		},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{"sequence", "sequence-with-optional", "optional-int", "string"},
		},
		text: "42 \"foo\"",
		node: &node{
			name:  "sequence-with-optional",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "42"},
			}, {
				name:  "string",
				token: &token{value: "\"foo\""},
			}},
		},
	}, {
		msg: "sequence with optional item, missing",
		primitive: [][]interface{}{
			{"int", intToken},
			{"string", stringToken},
		},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{"sequence", "sequence-with-optional", "optional-int", "string"},
		},
		text: "\"foo\"",
		node: &node{
			name:  "sequence-with-optional",
			token: &token{value: "\"foo\""},
			nodes: []*node{{
				name:  "string",
				token: &token{value: "\"foo\""},
			}},
		},
	}, {
		msg:       "sequence with only optional, empty",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{
				"sequence",
				"sequence-with-only-optional",
				"optional-int",
				"optional-int",
				"optional-int",
			},
		},
		node: &node{
			name: "sequence-with-only-optional",
		},
	}, {
		msg:       "sequence with only optional, less",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{
				"sequence",
				"sequence-with-only-optional",
				"optional-int",
				"optional-int",
				"optional-int",
			},
		},
		text: "1 2",
		node: &node{
			name:  "sequence-with-only-optional",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "1"},
			}, {
				name:  "int",
				token: &token{value: "2"},
			}},
		},
	}, {
		msg:       "sequence with only optional, exact",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{
				"sequence",
				"sequence-with-only-optional",
				"optional-int",
				"optional-int",
				"optional-int",
			},
		},
		text: "1 2 3",
		node: &node{
			name:  "sequence-with-only-optional",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "1"},
			}, {
				name:  "int",
				token: &token{value: "2"},
			}, {
				name:  "int",
				token: &token{value: "3"},
			}},
		},
	}, {
		msg:       "sequence with only optional, more",
		primitive: [][]interface{}{{"int", intToken}},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{
				"sequence",
				"sequence-with-only-optional",
				"optional-int",
				"optional-int",
				"optional-int",
			},
		},
		text: "1 2 3 4",
		fail: true,
	}, {
		msg: "choice of int and string",
		primitive: [][]interface{}{
			{"int", intToken},
			{"string", stringToken},
		},
		complex: [][]string{
			{"choice", "int-or-string", "int", "string"},
		},
		text: "\"foo\"",
		node: &node{
			name:  "string",
			token: &token{value: "\"foo\""},
		},
	}, {
		msg: "choice of int and sequence with optional int",
		primitive: [][]interface{}{
			{"int", intToken},
			{"string", stringToken},
		},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{"sequence", "sequence-with-optional", "optional-int", "string"},
			{"choice", "int-or-sequence-with-optional", "int", "sequence-with-optional"},
		},
		text: "42 \"foo\"",
		node: &node{
			name:  "sequence-with-optional",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int",
				token: &token{value: "42"},
			}, {
				name:  "string",
				token: &token{value: "\"foo\""},
			}},
		},
	}, {
		msg: "choice of int and sequence with optional int, token fall through",
		primitive: [][]interface{}{
			{"int", intToken},
			{"string", stringToken},
		},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{
				"sequence",
				"sequence-with-optional",
				"optional-int",
				"optional-int",
				"string",
				"string",
			},
			{"choice", "int-or-sequence-with-optional", "int", "sequence-with-optional"},
		},
		text: "\"foo\" \"bar\"",
		node: &node{
			name:  "sequence-with-optional",
			token: &token{value: "\"foo\""},
			nodes: []*node{{
				name:  "string",
				token: &token{value: "\"foo\""},
			}, {
				name:  "string",
				token: &token{value: "\"bar\""},
			}},
		},
	}, {
		msg: "choice of int and sequence with optional int, init fall through",
		primitive: [][]interface{}{
			{"int", intToken},
			{"string", stringToken},
		},
		complex: [][]string{
			{"optional", "optional-int", "int"},
			{
				"sequence",
				"sequence-with-optional",
				"optional-int",
				"optional-int",
				"string",
				"string",
			},
			{"choice", "int-or-sequence-with-optional", "int", "sequence-with-optional"},
		},
		text: "\"foo\" \"bar\"",
		node: &node{
			name:  "sequence-with-optional",
			token: &token{value: "\"foo\""},
			nodes: []*node{{
				name:  "string",
				token: &token{value: "\"foo\""},
			}, {
				name:  "string",
				token: &token{value: "\"bar\""},
			}},
		},
	}, {
		msg: "expression inside expression",
		primitive: [][]interface{}{
			{"symbol", symbolToken},
			{"symbol-word", symbolWord},
			{"open-paren", openParen},
			{"close-paren", closeParen},
		},
		complex: [][]string{
			{"sequence", "function-call", "expression", "open-paren", "expression", "close-paren"},
			{"sequence", "dynamic-symbol", "symbol-word", "open-paren", "expression", "close-paren"},
			{"choice", "expression", "symbol", "function-call", "dynamic-symbol"},
		},
		text: "symbol(f(a))",
		node: &node{
			name:  "dynamic-symbol",
			token: &token{value: "symbol"},
			nodes: []*node{{
				name:  "symbol-word",
				token: &token{value: "symbol"},
			}, {
				name:  "open-paren",
				token: &token{value: "("},
			}, {
				name:  "function-call",
				token: &token{value: "f"},
				nodes: []*node{{
					name:  "symbol",
					token: &token{value: "f"},
				}, {
					name:  "open-paren",
					token: &token{value: "("},
				}, {
					name:  "symbol",
					token: &token{value: "a"},
				}, {
					name:  "close-paren",
					token: &token{value: ")"},
				}},
			}, {
				name:  "close-paren",
				token: &token{value: ")"},
			}},
		},
	}, {
		msg: "chained symbol query",
		primitive: [][]interface{}{
			{"symbol", symbolToken},
			{"dot", dot},
		},
		complex: [][]string{
			{"sequence", "symbol-query", "expression", "dot", "symbol"},
			{"choice", "expression", "symbol", "symbol-query"},
		},
		text: "a.b.c",
		node: &node{
			name:  "symbol-query",
			token: &token{value: "a"},
			nodes: []*node{{
				name:  "symbol-query",
				token: &token{value: "a"},
				nodes: []*node{{
					name:  "symbol",
					token: &token{value: "a"},
				}, {
					name:  "dot",
					token: &token{value: "."},
				}, {
					name:  "symbol",
					token: &token{value: "b"},
				}},
			}, {
				name:  "dot",
				token: &token{value: "."},
			}, {
				name:  "symbol",
				token: &token{value: "c"},
			}},
		},
	}, {
		msg: "repetition in repetition",
		primitive: [][]interface{}{
			{"int", intToken},
		},
		complex: [][]string{
			{"repetition", "int-repetition", "int"},
			{"repetition", "repetition-in-repetition", "int-repetition"},
		},
		text: "42",
		node: &node{
			name:  "repetition-in-repetition",
			token: &token{value: "42"},
			nodes: []*node{{
				name:  "int-repetition",
				token: &token{value: "42"},
				nodes: []*node{{
					name:  "int",
					token: &token{value: "42"},
				}},
			}},
		},
	}, {
		msg: "reproduce repetition endless loop",
		primitive: [][]interface{}{
			{"nl", nl},
			{"colon", colon},
			{"switch-word", switchWord},
			{"case-word", caseWord},
			{"default-word", defaultWord},
			{"open-brace", openBrace},
			{"close-brace", closeBrace},
			{"symbol", symbolToken},
		},
		complex: [][]string{
			{"repetition", "nls", "nl"},
			{"choice", "match-expression", "expression"},
			{"sequence", "switch-clause", "case-word", "match-expression", "colon",
				"statement-repetition"},
			{"repetition", "switch-clause-repetition", "switch-clause"},
			{"sequence", "default-clause", "default-word", "colon", "nls", "statement-repetition"},
			{"choice", "seq-sep", "nl"},
			{"choice", "statement-repetition-item", "expression", "seq-sep"},
			{"repetition", "statement-repetition", "statement-repetition-item"},
			{
				"sequence",
				"switch-conditional",
				"switch-word",
				"nls",
				"open-brace",
				"nls",
				"switch-clause-repetition",
				"nls",
				"default-clause",
				"nls",
				"switch-clause-repetition",
				"nls",
				"close-brace",
			},
			{"choice", "expression", "symbol", "switch-conditional"},
			{"repetition", "document", "statement-repetition"},
		},
		text: `switch {
				default: a
			}`,
		node: &node{
			name:  "document",
			token: &token{value: "switch"},
			nodes: []*node{{
				name:  "statement-repetition",
				token: &token{value: "switch"},
				nodes: []*node{{
					name:  "switch-conditional",
					token: &token{value: "switch"},
					nodes: []*node{{
						name:  "switch-word",
						token: &token{value: "switch"},
					}, {
						name:  "nls",
						token: &token{value: "{"},
					}, {
						name:  "open-brace",
						token: &token{value: "{"},
					}, {
						name:  "nls",
						token: &token{value: "\n"},
						nodes: []*node{{
							name:  "nl",
							token: &token{value: "\n"},
						}},
					}, {
						name:  "switch-clause-repetition",
						token: &token{value: "default"},
					}, {
						name:  "nls",
						token: &token{value: "default"},
					}, {
						name:  "default-clause",
						token: &token{value: "default"},
						nodes: []*node{{
							name:  "default-word",
							token: &token{value: "default"},
						}, {
							name:  "colon",
							token: &token{value: ":"},
						}, {
							name:  "nls",
							token: &token{value: "a"},
						}, {
							name:  "statement-repetition",
							token: &token{value: "a"},
							nodes: []*node{{
								name:  "symbol",
								token: &token{value: "a"},
							}, {
								name:  "nl",
								token: &token{value: "\n"},
							}},
						}},
					}, {
						name:  "nls",
						token: &token{value: "}"},
					}, {
						name:  "switch-clause-repetition",
						token: &token{value: "}"},
					}, {
						name:  "nls",
						token: &token{value: "}"},
					}, {
						name:  "close-brace",
						token: &token{value: "}"},
					}},
				}},
			}},
		},
	}, {
		msg: "newline in sequence",
		primitive: [][]interface{}{
			{"nl", nl},
			{"colon", colon},
			{"switch-word", switchWord},
			{"case-word", caseWord},
			{"default-word", defaultWord},
			{"open-brace", openBrace},
			{"close-brace", closeBrace},
			{"symbol", symbolToken},
		},
		complex: [][]string{
			{"repetition", "nls", "nl"},
			{"choice", "match-expression", "expression"},
			{"sequence", "switch-clause", "case-word", "match-expression", "colon",
				"statement-repetition"},
			{"repetition", "switch-clause-repetition", "switch-clause"},
			{"sequence", "default-clause", "default-word", "colon", "nls", "statement-repetition"},
			{"choice", "seq-sep", "nl"},
			{"choice", "statement-repetition-item", "statement", "seq-sep"},
			{"repetition", "statement-repetition", "statement-repetition-item"},
			{
				"sequence",
				"switch-conditional",
				"switch-word",
				"nls",
				"open-brace",
				"nls",
				"switch-clause-repetition",
				"nls",
				"default-clause",
				"nls",
				"switch-clause-repetition",
				"nls",
				"close-brace",
			},
			{"choice", "conditional", "switch-conditional"},
			{"choice", "expression", "symbol", "conditional"},
			{"choice", "statement", "expression"},
			{"choice", "document", "statement-repetition"},
		},
		text: `switch {
				default: a
			}`,
		node: &node{
			name:  "statement-repetition",
			token: &token{value: "switch"},
			nodes: []*node{{
				name:  "switch-conditional",
				token: &token{value: "switch"},
				nodes: []*node{{
					name:  "switch-word",
					token: &token{value: "switch"},
				}, {
					name:  "nls",
					token: &token{value: "{"},
				}, {
					name:  "open-brace",
					token: &token{value: "{"},
				}, {
					name:  "nls",
					token: &token{value: "\n"},
					nodes: []*node{{
						name:  "nl",
						token: &token{value: "\n"},
					}},
				}, {
					name:  "switch-clause-repetition",
					token: &token{value: "default"},
				}, {
					name:  "nls",
					token: &token{value: "default"},
				}, {
					name:  "default-clause",
					token: &token{value: "default"},
					nodes: []*node{{
						name:  "default-word",
						token: &token{value: "default"},
					}, {
						name:  "colon",
						token: &token{value: ":"},
					}, {
						name:  "nls",
						token: &token{value: "a"},
					}, {
						name:  "statement-repetition",
						token: &token{value: "a"},
						nodes: []*node{{
							name:  "symbol",
							token: &token{value: "a"},
						}, {
							name:  "nl",
							token: &token{value: "\n"},
						}},
					}},
				}, {
					name:  "nls",
					token: &token{value: "}"},
				}, {
					name:  "switch-clause-repetition",
					token: &token{value: "}"},
				}, {
					name:  "nls",
					token: &token{value: "}"},
				}, {
					name:  "close-brace",
					token: &token{value: "}"},
				}},
			}},
		},
	}, {
		msg: "recursive repetition",
		primitive: [][]interface{}{
			{"int", intToken},
		},
		// ints = int | int ints
		complex: [][]string{
			{"sequence", "int-sequence", "ints", "int"},
			{"choice", "ints", "int", "int-sequence"},
		},
		text: "1 2 3",
		node: &node{
			name:  "int-sequence",
			token: &token{value: "1"},
			nodes: []*node{{
				name:  "int-sequence",
				token: &token{value: "1"},
				nodes: []*node{{
					name:  "int",
					token: &token{value: "1"},
				}, {
					name:  "int",
					token: &token{value: "2"},
				}},
			}, {
				name:  "int",
				token: &token{value: "3"},
			}},
		},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			var l traceLevel
			// l = traceDebug
			trace := newTrace(l)
			s := withTrace(trace)

			err := s.defineSyntax(ti.primitive, ti.complex)
			if err != nil {
				t.Error(err)
				return
			}

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
