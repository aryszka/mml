package mml

// import (
// 	"bytes"
// 	"testing"
// )
//
// func compareNode(t *testing.T, a, b node) {
// 	if a.typ != b.typ {
// 		t.Fatal("invalid node type", a.typ, b.typ)
// 	}
//
// 	if a.token.value != b.token.value {
// 		t.Fatal("invalid token value", a.token.value, b.token.value)
// 	}
//
// 	compareNodes(t, a.nodes, b.nodes)
// }
//
// func compareNodes(t *testing.T, a, b []node) {
// 	if len(a) != len(b) {
// 		for _, ai := range a {
// 			t.Log(ai)
// 		}
//
// 		t.Fatal("invalid node length", len(a), len(b))
// 		return
// 	}
//
// 	for i, ai := range a {
// 		compareNode(t, ai, b[i])
// 	}
// }
//
// func TestParse(t *testing.T) {
// 	for _, ti := range []struct {
// 		msg   string
// 		code  string
// 		nodes []node
// 		fail  bool
// 	}{{
// 		msg:   "int",
// 		code:  "42",
// 		nodes: []node{{typ: "int", token: token{value: "42"}}},
// 	}, {
// 		msg:  "multiple ints",
// 		code: "1 2; 3",
// 		nodes: []node{{
// 			typ:   "int",
// 			token: token{value: "1"},
// 		}, {
// 			typ:   "int",
// 			token: token{value: "2"},
// 		}, {
// 			typ:   "int",
// 			token: token{value: "3"},
// 		}},
// 	}, {
// 		msg:  "string",
// 		code: "\"abc\"",
// 		nodes: []node{{
// 			typ:   "string",
// 			token: token{value: "\"abc\""},
// 		}},
// 	}, {
// 		msg:  "symbol",
// 		code: "a",
// 		nodes: []node{{
// 			typ:   "symbol",
// 			token: token{value: "a"},
// 		}},
// 	}, {
// 		msg:  "dynamic symbol",
// 		code: "symbol(f(a))",
// 		nodes: []node{{
// 			typ:   "dynamic-symbol",
// 			token: token{value: "symbol"},
// 			nodes: []node{{
// 				typ:   "function-call",
// 				token: token{value: "f"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "f"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "bool",
// 		code: "true false",
// 		nodes: []node{{
// 			typ:   "true",
// 			token: token{value: "true"},
// 		}, {
// 			typ:   "false",
// 			token: token{value: "false"},
// 		}},
// 	}, {
// 		msg:  "empty list",
// 		code: "[]",
// 		nodes: []node{{
// 			typ:   "list",
// 			token: token{value: "["},
// 		}},
// 	}, {
// 		msg:  "list",
// 		code: "[1, 2, f(a), [3, 4, []]]",
// 		nodes: []node{{
// 			typ:   "list",
// 			token: token{value: "["},
// 			nodes: []node{{
// 				typ:   "int",
// 				token: token{value: "1"},
// 			}, {
// 				typ:   "int",
// 				token: token{value: "2"},
// 			}, {
// 				typ:   "function-call",
// 				token: token{value: "f"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "f"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}, {
// 				typ:   "list",
// 				token: token{value: "["},
// 				nodes: []node{{
// 					typ:   "int",
// 					token: token{value: "3"},
// 				}, {
// 					typ:   "int",
// 					token: token{value: "4"},
// 				}, {
// 					typ:   "list",
// 					token: token{value: "["},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "mutable-list",
// 		code: "~[1, 2, f(a), [3, 4, ~[]]]",
// 		nodes: []node{{
// 			typ:   "mutable-list",
// 			token: token{value: "~"},
// 			nodes: []node{{
// 				typ:   "int",
// 				token: token{value: "1"},
// 			}, {
// 				typ:   "int",
// 				token: token{value: "2"},
// 			}, {
// 				typ:   "function-call",
// 				token: token{value: "f"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "f"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}, {
// 				typ:   "list",
// 				token: token{value: "["},
// 				nodes: []node{{
// 					typ:   "int",
// 					token: token{value: "3"},
// 				}, {
// 					typ:   "int",
// 					token: token{value: "4"},
// 				}, {
// 					typ:   "mutable-list",
// 					token: token{value: "~"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "empty structure",
// 		code: "{}",
// 		nodes: []node{{
// 			typ:   "structure",
// 			token: token{value: "{"},
// 		}},
// 	}, {
// 		msg:  "structure",
// 		code: "{a: 1, b: 2, ...c, d: {e: 3, f: {}}}",
// 		nodes: []node{{
// 			typ:   "structure",
// 			token: token{value: "{"},
// 			nodes: []node{{
// 				typ:   "structure-definition",
// 				token: token{value: "a"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}, {
// 					typ:   "int",
// 					token: token{value: "1"},
// 				}},
// 			}, {
// 				typ:   "structure-definition",
// 				token: token{value: "b"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "b"},
// 				}, {
// 					typ:   "int",
// 					token: token{value: "2"},
// 				}},
// 			}, {
// 				typ:   "spread-expression",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}, {
// 				typ:   "structure-definition",
// 				token: token{value: "d"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "d"},
// 				}, {
// 					typ:   "structure",
// 					token: token{value: "{"},
// 					nodes: []node{{
// 						typ:   "structure-definition",
// 						token: token{value: "e"},
// 						nodes: []node{{
// 							typ:   "symbol",
// 							token: token{value: "e"},
// 						}, {
// 							typ:   "int",
// 							token: token{value: "3"},
// 						}},
// 					}, {
// 						typ:   "structure-definition",
// 						token: token{value: "f"},
// 						nodes: []node{{
// 							typ:   "symbol",
// 							token: token{value: "f"},
// 						}, {
// 							typ:   "structure",
// 							token: token{value: "{"},
// 						}},
// 					}},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "mutable structure",
// 		code: "~{a: 1, b: 2, ...c, d: {e: 3, f: ~{}}}",
// 		nodes: []node{{
// 			typ:   "mutable-structure",
// 			token: token{value: "~"},
// 			nodes: []node{{
// 				typ:   "structure-definition",
// 				token: token{value: "a"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}, {
// 					typ:   "int",
// 					token: token{value: "1"},
// 				}},
// 			}, {
// 				typ:   "structure-definition",
// 				token: token{value: "b"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "b"},
// 				}, {
// 					typ:   "int",
// 					token: token{value: "2"},
// 				}},
// 			}, {
// 				typ:   "spread-expression",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}, {
// 				typ:   "structure-definition",
// 				token: token{value: "d"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "d"},
// 				}, {
// 					typ:   "structure",
// 					token: token{value: "{"},
// 					nodes: []node{{
// 						typ:   "structure-definition",
// 						token: token{value: "e"},
// 						nodes: []node{{
// 							typ:   "symbol",
// 							token: token{value: "e"},
// 						}, {
// 							typ:   "int",
// 							token: token{value: "3"},
// 						}},
// 					}, {
// 						typ:   "structure-definition",
// 						token: token{value: "f"},
// 						nodes: []node{{
// 							typ:   "symbol",
// 							token: token{value: "f"},
// 						}, {
// 							typ:   "mutable-structure",
// 							token: token{value: "~"},
// 						}},
// 					}},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "symbol query",
// 		code: "a.b",
// 		nodes: []node{{
// 			typ:   "symbol-query",
// 			token: token{value: "a"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "chained symbol query",
// 		code: "a.b.c",
// 		nodes: []node{{
// 			typ:   "symbol-query",
// 			token: token{value: "a"},
// 			nodes: []node{{
// 				typ:   "symbol-query",
// 				token: token{value: "a"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "b"},
// 				}},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "c"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "noop function",
// 		code: "fn () {;}",
// 		nodes: []node{{
// 			typ:   "function",
// 			token: token{value: "fn"},
// 			nodes: []node{{
// 				typ:   "statement-sequence",
// 				token: token{value: ";"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "simple function",
// 		code: "fn () 3",
// 		nodes: []node{{
// 			typ:   "function",
// 			token: token{value: "fn"},
// 			nodes: []node{{
// 				typ:   "int",
// 				token: token{value: "3"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "identity",
// 		code: "fn (x) x",
// 		nodes: []node{{
// 			typ:   "function",
// 			token: token{value: "fn"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "x"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "x"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "list",
// 		code: "fn (...x) x",
// 		nodes: []node{{
// 			typ:   "function",
// 			token: token{value: "fn"},
// 			nodes: []node{{
// 				typ:   "collect-symbol",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "x"},
// 				}},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "x"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "function",
// 		code: "fn (a, b, ...c) { c }",
// 		nodes: []node{{
// 			typ:   "function",
// 			token: token{value: "fn"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}, {
// 				typ:   "collect-symbol",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}, {
// 				typ:   "statement-sequence",
// 				token: token{value: "c"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "function with sequence",
// 		code: "fn (a, b, ...c) { a(b); c }",
// 		nodes: []node{{
// 			typ:   "function",
// 			token: token{value: "fn"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}, {
// 				typ:   "collect-symbol",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}, {
// 				typ:   "statement-sequence",
// 				token: token{value: "a"},
// 				nodes: []node{{
// 					typ:   "function-call",
// 					token: token{value: "a"},
// 					nodes: []node{{
// 						typ:   "symbol",
// 						token: token{value: "a"},
// 					}, {
// 						typ:   "symbol",
// 						token: token{value: "b"},
// 					}},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "function call",
// 		code: "f(a)",
// 		nodes: []node{{
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "f"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "chained function call",
// 		code: "f(a)(b)",
// 		nodes: []node{{
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "function-call",
// 				token: token{value: "f"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "f"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "chained function call, whitespace",
// 		code: "f(a) (b)",
// 		nodes: []node{{
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "function-call",
// 				token: token{value: "f"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "f"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "function call argument",
// 		code: "f(g(a))",
// 		nodes: []node{{
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "f"},
// 			}, {
// 				typ:   "function-call",
// 				token: token{value: "g"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "g"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "function call sequence",
// 		code: "f(a) f(b)g(a)",
// 		nodes: []node{{
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "f"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}},
// 		}, {
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "f"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}},
// 		}, {
// 			typ:   "function-call",
// 			token: token{value: "g"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "g"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "function call with multiple arguments",
// 		code: "f(...a, b, ...c)",
// 		nodes: []node{{
// 			typ:   "function-call",
// 			token: token{value: "f"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "f"},
// 			}, {
// 				typ:   "spread-expression",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}, {
// 				typ:   "spread-expression",
// 				token: token{value: "."},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "switch conditional with default only",
// 		code: "switch{default: 42}",
// 		nodes: []node{{
// 			typ:   "switch-conditional",
// 			token: token{value: "switch"},
// 			nodes: []node{{
// 				typ:   "default-clause",
// 				token: token{value: "default"},
// 				nodes: []node{{
// 					typ:   "int",
// 					token: token{value: "42"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg: "switch conditional with cases",
// 		code: `
// 			switch {
// 				case a: b
// 				default: x
// 				case c: d
// 			}`,
// 		nodes: []node{{
// 			typ:   "switch-conditional",
// 			token: token{value: "switch"},
// 			nodes: []node{{
// 				typ:   "switch-clause",
// 				token: token{value: "case"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "a"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "b"},
// 				}},
// 			}, {
// 				typ:   "default-clause",
// 				token: token{value: "default"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "x"},
// 				}},
// 			}, {
// 				typ:   "switch-clause",
// 				token: token{value: "case"},
// 				nodes: []node{{
// 					typ:   "symbol",
// 					token: token{value: "c"},
// 				}, {
// 					typ:   "symbol",
// 					token: token{value: "d"},
// 				}},
// 			}},
// 		}},
// 	}, {
// 		msg:  "definition",
// 		code: "let a b",
// 		nodes: []node{{
// 			typ:   "value-definition",
// 			token: token{value: "let"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}},
// 		}},
// 	}, {
// 		msg:  "mutable definition",
// 		code: "let ~ a b",
// 		nodes: []node{{
// 			typ:   "mutable-value-definition",
// 			token: token{value: "let"},
// 			nodes: []node{{
// 				typ:   "symbol",
// 				token: token{value: "a"},
// 			}, {
// 				typ:   "symbol",
// 				token: token{value: "b"},
// 			}},
// 		}},
// 	}} {
// 		t.Run(ti.msg, func(t *testing.T) {
// 			cache = tokenCache{tokens: make(map[token]tokenCacheItem)}
// 			r := newTokenReader(bytes.NewBufferString(ti.code), "<test>")
// 			n, err := parse(generators["document"], r)
// 			if err == nil && ti.fail {
// 				t.Fatal("failed to fail")
// 			} else if err != nil && !ti.fail {
// 				t.Fatal(err)
// 			} else if err != nil {
// 				return
// 			}
//
// 			compareNodes(t, n.nodes, ti.nodes)
// 		})
// 	}
// }
