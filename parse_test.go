package mml

import (
	"io"
	"strings"
	"testing"
)

func compareNodes(a, b []node) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ai := range a {
		if ai.typ != b[i].typ {
			return false
		}

		if ai.token.value != b[i].token.value {
			return false
		}

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
		msg: "mutable list",
		code: `~[1, 2, 3]`,
		nodes: []node{{
			typ: mutableListNode,
			token: token{value: "~"},
			nodes: []node{{
				typ: intNode,
				token: token{value: "1"},
			}, {
				typ: intNode,
				token: token{value: "2"},
			}, {
				typ: intNode,
				token: token{value: "3"},
			}},
		}},
	}, {
		msg:  "fail",
		code: `[`,
		fail: true,
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
			}

			if !compareNodes(n, ti.nodes) {
				t.Error("invalid parse result")
			}
		})
	}
}
