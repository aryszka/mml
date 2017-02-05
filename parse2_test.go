package mml

import (
	"bytes"
	"testing"
)

func compareNode(t *testing.T, a, b node) {
	if a.typ != b.typ {
		t.Fatal("invalid node type", a.typ, b.typ)
	}

	if a.token.value != b.token.value {
		t.Fatal("invalid token value", a.token.value, b.token.value)
	}

	compareNodes(t, a.nodes, b.nodes)
}

func compareNodes(t *testing.T, a, b []node) {
	if len(a) != len(b) {
		t.Fatal("invalid node length", len(a), len(b))
		return
	}

	for i, ai := range a {
		compareNode(t, ai, b[i])
	}
}

func TestParse(t *testing.T) {
	for _, ti := range []struct {
		msg   string
		code  string
		nodes []node
		fail  bool
	}{{
		msg:   "int",
		code:  "42",
		nodes: []node{{typ: "int", token: token{value: "42"}}},
	}, {
		msg:  "multiple ints",
		code: "1 2; 3",
		nodes: []node{{
			typ:   "int",
			token: token{value: "1"},
		}, {
			typ:   "int",
			token: token{value: "2"},
		}, {
			typ:   "int",
			token: token{value: "3"},
		}},
	}, {
		msg:  "symbol",
		code: "a",
		nodes: []node{{
			typ:   "symbol",
			token: token{value: "a"},
		}},
	}, {
		msg:  "string",
		code: "\"abc\"",
		nodes: []node{{
			typ:   "string",
			token: token{value: "\"abc\""},
		}},
	}, {
		msg:  "noop function",
		code: "fn () {;}",
		nodes: []node{{
			typ:   "function",
			token: token{value: "fn"},
			nodes: []node{{
				typ:   "statement-sequence",
				token: token{value: ";"},
			}},
		}},
	}, {
		msg:  "simple function",
		code: "fn () 3",
		nodes: []node{{
			typ:   "function",
			token: token{value: "fn"},
			nodes: []node{{
				typ:   "int",
				token: token{value: "3"},
			}},
		}},
	}, {
		msg:  "identity",
		code: "fn (x) x",
		nodes: []node{{
			typ:   "function",
			token: token{value: "fn"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "x"},
			}, {
				typ:   "symbol",
				token: token{value: "x"},
			}},
		}},
	}, {
		msg:  "list",
		code: "fn (...x) x",
		nodes: []node{{
			typ:   "function",
			token: token{value: "fn"},
			nodes: []node{{
				typ:   "collect-symbol",
				token: token{value: "."},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "x"},
				}},
			}, {
				typ:   "symbol",
				token: token{value: "x"},
			}},
		}},
	}, {
		msg:  "function",
		code: "fn (a, b, ...c) { a(b); c }",
		nodes: []node{{
			typ:   "function",
			token: token{value: "fn"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "a"},
			}, {
				typ:   "symbol",
				token: token{value: "b"},
			}, {
				typ:   "collect-symbol",
				token: token{value: "."},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "c"},
				}},
			}, {
				typ:   "statement-sequence",
				token: token{value: "a"},
				nodes: []node{{
					typ:   "function-call",
					token: token{value: "a"},
					nodes: []node{{
						typ:   "symbol",
						token: token{value: "a"},
					}, {
						typ:   "symbol",
						token: token{value: "b"},
					}},
				}, {
					typ:   "symbol",
					token: token{value: "c"},
				}},
			}},
		}},
	}, {
		msg:  "function call",
		code: "f(a)",
		nodes: []node{{
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "f"},
			}, {
				typ:   "symbol",
				token: token{value: "a"},
			}},
		}},
	}, {
		msg:  "chained function call",
		code: "f(a)(b)",
		nodes: []node{{
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "function-call",
				token: token{value: "f"},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "f"},
				}, {
					typ:   "symbol",
					token: token{value: "a"},
				}},
			}, {
				typ:   "symbol",
				token: token{value: "b"},
			}},
		}},
	}, {
		msg:  "chained function call, whitespace",
		code: "f(a) (b)",
		nodes: []node{{
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "function-call",
				token: token{value: "f"},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "f"},
				}, {
					typ:   "symbol",
					token: token{value: "a"},
				}},
			}, {
				typ:   "symbol",
				token: token{value: "b"},
			}},
		}},
	}, {
		msg:  "function call argument",
		code: "f(g(a))",
		nodes: []node{{
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "f"},
			}, {
				typ:   "function-call",
				token: token{value: "g"},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "g"},
				}, {
					typ:   "symbol",
					token: token{value: "a"},
				}},
			}},
		}},
	}, {
		msg:  "function call sequence",
		code: "f(a) f(b)g(a)",
		nodes: []node{{
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "f"},
			}, {
				typ:   "symbol",
				token: token{value: "a"},
			}},
		}, {
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "f"},
			}, {
				typ:   "symbol",
				token: token{value: "b"},
			}},
		}, {
			typ:   "function-call",
			token: token{value: "g"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "g"},
			}, {
				typ:   "symbol",
				token: token{value: "a"},
			}},
		}},
	}, {
		msg:  "function call with multiple arguments",
		code: "f(a..., b, c...)",
		nodes: []node{{
			typ:   "function-call",
			token: token{value: "f"},
			nodes: []node{{
				typ:   "symbol",
				token: token{value: "f"},
			}, {
				typ:   "spread-expression",
				token: token{value: "a"},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "a"},
				}},
			}, {
				typ:   "symbol",
				token: token{value: "b"},
			}, {
				typ:   "spread-expression",
				token: token{value: "c"},
				nodes: []node{{
					typ:   "symbol",
					token: token{value: "c"},
				}},
			}},
		}},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			r := newTokenReader(bytes.NewBufferString(ti.code), "<test>")
			n, err := parse(parsers["document"], r)
			if err == nil && ti.fail {
				t.Error("failed to fail")
			} else if err != nil && !ti.fail {
				t.Error(err)
			} else if err != nil {
				return
			}

			compareNodes(t, n.nodes, ti.nodes)
		})
	}
}
