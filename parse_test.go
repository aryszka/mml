package mml

import (
	"io"
	"strings"
	"testing"
)

func compareNodes(a, b []node) bool {
	if len(a) != len(b) {
		// println("length wrong", len(a), len(b))
		return false
	}

	// println("length ok", len(a))

	for i, ai := range a {
		if ai.typ != b[i].typ {
			// println("type wrong", ai.typ.String(), b[i].typ.String())
			return false
		}

		// println("type ok", ai.typ.String())

		if ai.token.value != b[i].token.value {
			// println("token wrong", ai.token.value, b[i].token.value)
			return false
		}

		// println("token ok", ai.token.value)

		if !compareNodes(ai.nodes, b[i].nodes) {
			return false
		}
	}

	return true
}

func TestParse(t *testing.T) {
	for _, ti := range []struct {
		msg   string
		code  string
		nodes []node
		fail  bool
	}{{
		msg:  "empty doc",
		code: "",
	}, {
		msg:   "int",
		code:  `1`,
		nodes: []node{{typ: intNode, token: token{value: "1"}}},
	}, {
		msg:   "string",
		code:  `"foo"`,
		nodes: []node{{typ: stringNode, token: token{value: "foo"}}},
	}, {
		msg:   "string with escape",
		code:  "\"foo\\n\\\"bar\\\"\"",
		nodes: []node{{typ: stringNode, token: token{value: "foo\n\"bar\""}}},
	}, {
		msg:  "channel",
		code: "<>",
		nodes: []node{{
			typ:   channelNode,
			token: token{value: "<"},
		}},
	}, {
		msg:  "symbol",
		code: "foo",
		nodes: []node{{
			typ:   symbolNode,
			token: token{value: "foo"},
		}},
	}, {
		msg:  "dynamic symbol",
		code: "symbol(foo)",
		nodes: []node{{
			typ:   dynamicSymbolNode,
			token: token{value: "symbol"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "foo"},
			}},
		}},
	}, {
		msg:  "true",
		code: "true",
		nodes: []node{{
			typ:   trueNode,
			token: token{value: "true"},
		}},
	}, {
		msg:  "false",
		code: "false",
		nodes: []node{{
			typ:   falseNode,
			token: token{value: "false"},
		}},
	}, {
		msg:  "empty list",
		code: `[]`,
		nodes: []node{{
			typ:   listNode,
			token: token{value: "["},
		}},
	}, {
		msg:  "list with single item",
		code: `[2]`,
		nodes: []node{{
			typ:   listNode,
			token: token{value: "["},
			nodes: []node{{
				typ:   intNode,
				token: token{value: "2"},
			}},
		}},
	}, {
		msg:  "list with multiple items",
		code: `[1, 2, 3]`,
		nodes: []node{{
			typ:   listNode,
			token: token{value: "["},
			nodes: []node{{
				typ:   intNode,
				token: token{value: "1"},
			}, {
				typ:   intNode,
				token: token{value: "2"},
			}, {
				typ:   intNode,
				token: token{value: "3"},
			}},
		}},
	}, {
		msg:  "list with different items and list",
		code: `[1, a, 3, ["foo", symbol(bar), [6]]]`,
		nodes: []node{{
			typ:   listNode,
			token: token{value: "["},
			nodes: []node{{
				typ:   intNode,
				token: token{value: "1"},
			}, {
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   intNode,
				token: token{value: "3"},
			}, {
				typ:   listNode,
				token: token{value: "["},
				nodes: []node{{
					typ:   stringNode,
					token: token{value: "foo"},
				}, {
					typ:   dynamicSymbolNode,
					token: token{value: "symbol"},
					nodes: []node{{
						typ:   symbolNode,
						token: token{value: "bar"},
					}},
				}, {
					typ:   listNode,
					token: token{value: "["},
					nodes: []node{{
						typ:   intNode,
						token: token{value: "6"},
					}},
				}},
			}},
		}},
	}, {
		msg:  "comment in list",
		code: `[1, /* comment */ 2]`,
		nodes: []node{{
			typ:   listNode,
			token: token{value: "["},
			nodes: []node{{
				typ:   intNode,
				token: token{value: "1"},
			}, {
				typ:   intNode,
				token: token{value: "2"},
			}},
		}},
	}, {
		msg: "new line in list",
		code: `[1
			, 2]`,
		nodes: []node{{
			typ:   listNode,
			token: token{value: "["},
			nodes: []node{{
				typ:   intNode,
				token: token{value: "1"},
			}, {
				typ:   intNode,
				token: token{value: "2"},
			}},
		}},
	}, {
		msg:  "mutable list",
		code: `~[1, 2, 3]`,
		nodes: []node{{
			typ:   mutableListNode,
			token: token{value: "~"},
			nodes: []node{{
				typ:   intNode,
				token: token{value: "1"},
			}, {
				typ:   intNode,
				token: token{value: "2"},
			}, {
				typ:   intNode,
				token: token{value: "3"},
			}},
		}},
	}, {
		msg:  "empty structure",
		code: `{}`,
		nodes: []node{{
			typ:   structureNode,
			token: token{value: "{"},
		}},
	}, {
		msg:  "structure",
		code: `{foo: "bar"}`,
		nodes: []node{{
			typ:   structureNode,
			token: token{value: "{"},
			nodes: []node{{
				typ:   structureDefinitionNode,
				token: token{value: "foo"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "foo"},
				}, {
					typ:   stringNode,
					token: token{value: "bar"},
				}},
			}},
		}},
	}, {
		msg:  "structure with multiple items",
		code: `{foo: "bar", baz: {qux: "quux"}}`,
		nodes: []node{{
			typ:   structureNode,
			token: token{value: "{"},
			nodes: []node{{
				typ:   structureDefinitionNode,
				token: token{value: "foo"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "foo"},
				}, {
					typ:   stringNode,
					token: token{value: "bar"},
				}},
			}, {
				typ:   structureDefinitionNode,
				token: token{value: "baz"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "baz"},
				}, {
					typ:   structureNode,
					token: token{value: "{"},
					nodes: []node{{
						typ:   structureDefinitionNode,
						token: token{value: "qux"},
						nodes: []node{{
							typ:   symbolNode,
							token: token{value: "qux"},
						}, {
							typ:   stringNode,
							token: token{value: "quux"},
						}},
					}},
				}},
			}},
		}},
	}, {
		msg:  "mutable structure",
		code: `~{foo: "bar"}`,
		nodes: []node{{
			typ:   mutableStructureNode,
			token: token{value: "~"},
			nodes: []node{{
				typ:   structureDefinitionNode,
				token: token{value: "foo"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "foo"},
				}, {
					typ:   stringNode,
					token: token{value: "bar"},
				}},
			}},
		}},
	}, {
		msg:  "empty and",
		code: `and()`,
		nodes: []node{{
			typ:   andExpressionNode,
			token: token{value: "and"},
		}},
	}, {
		msg:  "and",
		code: `and(a, b, c)`,
		nodes: []node{{
			typ:   andExpressionNode,
			token: token{value: "and"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}, {
				typ:   symbolNode,
				token: token{value: "c"},
			}},
		}},
	}, {
		msg:  "empty or",
		code: `or()`,
		nodes: []node{{
			typ:   orExpressionNode,
			token: token{value: "or"},
		}},
	}, {
		msg:  "or",
		code: `or(a, b, c)`,
		nodes: []node{{
			typ:   orExpressionNode,
			token: token{value: "or"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}, {
				typ:   symbolNode,
				token: token{value: "c"},
			}},
		}},
	}, {
		msg:  "empty function",
		code: `fn () {;}`,
		nodes: []node{{
			typ:   functionNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   statementSequenceNode,
				token: token{value: ";"},
			}},
		}},
	}, {
		msg:  "function returning empty object",
		code: `fn () {}`,
		nodes: []node{{
			typ:   functionNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   structureNode,
				token: token{value: "{"},
			}},
		}},
	}, {
		msg:  "identity",
		code: `fn (x) x`,
		nodes: []node{{
			typ:   functionNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "x"},
			}, {
				typ:   symbolNode,
				token: token{value: "x"},
			}},
		}},
	}, {
		msg:  "list identity",
		code: `fn (...l) l`,
		nodes: []node{{
			typ:   functionNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   collectSymbolNode,
				token: token{value: "."},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "l"},
				}},
			}, {
				typ:   symbolNode,
				token: token{value: "l"},
			}},
		}},
	}, {
		msg:  "simple function",
		code: `fn (a, b) { a; b }`,
		nodes: []node{{
			typ:   functionNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}, {
				typ:   statementSequenceNode,
				token: token{value: "a"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "a"},
				}, {
					typ:   symbolNode,
					token: token{value: "b"},
				}},
			}},
		}},
	}, {
		msg:  "function with collect",
		code: `fn (a, b, ...c) { a; b; c }`,
		nodes: []node{{
			typ:   functionNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}, {
				typ:   collectSymbolNode,
				token: token{value: "."},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "c"},
				}},
			}, {
				typ:   statementSequenceNode,
				token: token{value: "a"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "a"},
				}, {
					typ:   symbolNode,
					token: token{value: "b"},
				}, {
					typ:   symbolNode,
					token: token{value: "c"},
				}},
			}},
		}},
	}, {
		msg:  "function effect",
		code: `fn~ (a, b, ...c) { a; b; c }`,
		nodes: []node{{
			typ:   functionEffectNode,
			token: token{value: "fn"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}, {
				typ:   collectSymbolNode,
				token: token{value: "."},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "c"},
				}},
			}, {
				typ:   statementSequenceNode,
				token: token{value: "a"},
				nodes: []node{{
					typ:   symbolNode,
					token: token{value: "a"},
				}, {
					typ:   symbolNode,
					token: token{value: "b"},
				}, {
					typ:   symbolNode,
					token: token{value: "c"},
				}},
			}},
		}},
	}, {
		msg:  "symbol query",
		code: `a.b`,
		nodes: []node{{
			typ:   symbolQueryNode,
			token: token{value: "a"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}},
		}},
	}, {
		msg:  "expression query",
		code: `a[b]`,
		nodes: []node{{
			typ:   expressionQueryNode,
			token: token{value: "a"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   symbolNode,
				token: token{value: "b"},
			}},
		}},
	}, {
		msg:  "expression query, infinite range",
		code: `a[:]`,
		nodes: []node{{
			typ:   expressionQueryNode,
			token: token{value: "a"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   rangeExpressionNode,
				token: token{value: ":"},
				nodes: []node{{}, {}},
			}},
		}},
	}, {
		msg:  "expression query, lower limit",
		code: `a[3:]`,
		nodes: []node{{
			typ:   expressionQueryNode,
			token: token{value: "a"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   rangeExpressionNode,
				token: token{value: "3"},
				nodes: []node{{
					typ:   intNode,
					token: token{value: "3"},
				}, {}},
			}},
		}},
	}, {
		msg:  "expression query, upper limit",
		code: `a[:3]`,
		nodes: []node{{
			typ:   expressionQueryNode,
			token: token{value: "a"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   rangeExpressionNode,
				token: token{value: ":"},
				nodes: []node{{}, {
					typ:   intNode,
					token: token{value: "3"},
				}},
			}},
		}},
	}, {
		msg:  "expression query, range",
		code: `a[3:42]`,
		nodes: []node{{
			typ:   expressionQueryNode,
			token: token{value: "a"},
			nodes: []node{{
				typ:   symbolNode,
				token: token{value: "a"},
			}, {
				typ:   rangeExpressionNode,
				token: token{value: "3"},
				nodes: []node{{
					typ:   intNode,
					token: token{value: "3"},
				}, {
					typ:   intNode,
					token: token{value: "42"},
				}},
			}},
		}},
	}, {
		msg:  "fail",
		code: `[`,
		fail: true,
	}, {
		msg: "document sequence",
		code: `1; 2
			3;
			4`,
		nodes: []node{{
			typ:   intNode,
			token: token{value: "1"},
		}, {
			typ:   intNode,
			token: token{value: "2"},
		}, {
			typ:   intNode,
			token: token{value: "3"},
		}, {
			typ:   intNode,
			token: token{value: "4"},
		}},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			r := strings.NewReader(ti.code)
			n, err := parse(r, "test")
			if ti.fail && (err == nil || err == io.EOF) {
				t.Error("failed to fail")
				return
			}

			if ti.fail {
				return
			}

			if err != nil && err != io.EOF {
				t.Error(err)
				return
			}

			if !compareNodes(n, ti.nodes) {
				t.Error("invalid parse result")
			}
		})
	}
}
