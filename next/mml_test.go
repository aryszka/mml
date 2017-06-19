package next

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestMML(t *testing.T) {
	trace := NewTrace(0)

	b, err := bootSyntax(trace)
	if err != nil {
		t.Error(err)
		return
	}

	mml, err := os.Open("mml.p")
	if err != nil {
		t.Error(err)
		return
	}

	defer mml.Close()

	mmlDoc, err := b.Parse(mml)
	if err != nil {
		t.Error(err)
		return
	}

	// trace = NewTrace(1)
	s := NewSyntax(trace)
	if err := define(s, mmlDoc); err != nil {
		t.Error(err)
		return
	}

	if err := s.Init(); err != nil {
		t.Error(err)
		return
	}

	start := time.Now()
	defer func() { t.Log("\nTestMML, total duration", time.Since(start)) }()
	for _, ti := range []struct {
		msg            string
		text           string
		fail           bool
		node           *Node
		nodes          []*Node
		ignorePosition bool
	}{{
		msg:  "empty",
		node: &Node{Name: "mml"},
	}, {
		msg:  "single line comment",
		text: "// foo bar baz",
		nodes: []*Node{{
			Name: "comment",
			to:   14,
			Nodes: []*Node{{
				Name: "line-comment-content",
				from: 2,
				to:   14,
			}},
		}},
	}, {
		msg:  "multiple line comments",
		text: "// foo bar\n// baz qux",
		nodes: []*Node{{
			Name: "comment",
			to:   21,
			Nodes: []*Node{{
				Name: "line-comment-content",
				from: 2,
				to:   10,
			}, {
				Name: "line-comment-content",
				from: 13,
				to:   21,
			}},
		}},
	}, {
		msg:  "block comment",
		text: "/* foo bar baz */",
		nodes: []*Node{{
			Name: "comment",
			to:   17,
			Nodes: []*Node{{
				Name: "block-comment-content",
				from: 2,
				to:   15,
			}},
		}},
	}, {
		msg:  "block comments",
		text: "/* foo bar */\n/* baz qux */",
		nodes: []*Node{{
			Name: "comment",
			to:   27,
			Nodes: []*Node{{
				Name: "block-comment-content",
				from: 2,
				to:   11,
			}, {
				Name: "block-comment-content",
				from: 16,
				to:   25,
			}},
		}},
	}, {
		msg:  "mixed comments",
		text: "// foo\n/* bar */\n// baz",
		nodes: []*Node{{
			Name: "comment",
			to:   23,
			Nodes: []*Node{{
				Name: "line-comment-content",
				from: 2,
				to:   6,
			}, {
				Name: "block-comment-content",
				from: 9,
				to:   14,
			}, {
				Name: "line-comment-content",
				from: 19,
				to:   23,
			}},
		}},
	}, {
		msg:  "int",
		text: "42",
		nodes: []*Node{{
			Name: "int",
			to:   2,
		}},
	}, {
		msg:  "ints",
		text: "1; 2; 3",
		nodes: []*Node{{
			Name: "int",
			to:   1,
		}, {
			Name: "int",
			from: 3,
			to:   4,
		}, {
			Name: "int",
			from: 6,
			to:   7,
		}},
	}, {
		msg:  "int, octal",
		text: "052",
		nodes: []*Node{{
			Name: "int",
			to:   3,
		}},
	}, {
		msg:  "int, hexa",
		text: "0x2a",
		nodes: []*Node{{
			Name: "int",
			to:   4,
		}},
	}, {
		msg:  "float, 0.",
		text: "0.",
		nodes: []*Node{{
			Name: "float",
			to:   2,
		}},
	}, {
		msg:  "float, 72.40",
		text: "72.40",
		nodes: []*Node{{
			Name: "float",
			to:   5,
		}},
	}, {
		msg:  "float, 072.40",
		text: "072.40",
		nodes: []*Node{{
			Name: "float",
			to:   6,
		}},
	}, {
		msg:  "float, 2.71828",
		text: "2.71828",
		nodes: []*Node{{
			Name: "float",
			to:   7,
		}},
	}, {
		msg:  "float, 6.67428e-11",
		text: "6.67428e-11",
		nodes: []*Node{{
			Name: "float",
			to:   11,
		}},
	}, {
		msg:  "float, 1E6",
		text: "1E6",
		nodes: []*Node{{
			Name: "float",
			to:   3,
		}},
	}, {
		msg:  "float, .25",
		text: ".25",
		nodes: []*Node{{
			Name: "float",
			to:   3,
		}},
	}, {
		msg:  "float, .12345E+5",
		text: ".12345E+5",
		nodes: []*Node{{
			Name: "float",
			to:   9,
		}},
	}, {
		msg:  "string, empty",
		text: "\"\"",
		nodes: []*Node{{
			Name: "string",
			to:   2,
		}},
	}, {
		msg:  "string",
		text: "\"foo\"",
		nodes: []*Node{{
			Name: "string",
			to:   5,
		}},
	}, {
		msg:  "string, with new line",
		text: "\"foo\nbar\"",
		nodes: []*Node{{
			Name: "string",
			to:   9,
		}},
	}, {
		msg:  "string, with escaped new line",
		text: "\"foo\\nbar\"",
		nodes: []*Node{{
			Name: "string",
			to:   10,
		}},
	}, {
		msg:  "string, with quotes",
		text: "\"foo \\\"bar\\\" baz\"",
		nodes: []*Node{{
			Name: "string",
			to:   17,
		}},
	}, {
		msg:  "bool, true",
		text: "true",
		nodes: []*Node{{
			Name: "true",
			to:   4,
		}},
	}, {
		msg:  "bool, false",
		text: "false",
		nodes: []*Node{{
			Name: "false",
			to:   5,
		}},
	}, {
		msg:  "symbol",
		text: "foo",
		nodes: []*Node{{
			Name: "symbol",
			to:   3,
		}},
	}, {
		msg:  "dynamic-symbol",
		text: "symbol(a)",
		nodes: []*Node{{
			Name: "dynamic-symbol",
			to:   9,
			Nodes: []*Node{{
				Name: "symbol",
				from: 7,
				to:   8,
			}},
		}},
	}, {
		msg:  "empty list",
		text: "[]",
		nodes: []*Node{{
			Name: "list",
			to:   2,
		}},
	}, {
		msg:  "list",
		text: "[a, b, c]",
		nodes: []*Node{{
			Name: "list",
			to:   9,
			Nodes: []*Node{{
				Name: "symbol",
				from: 1,
				to:   2,
			}, {
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 7,
				to:   8,
			}},
		}},
	}, {
		msg: "list, new lines",
		text: `[
			a
			b
			c
		]`,
		nodes: []*Node{{
			Name: "list",
			to:   20,
			Nodes: []*Node{{
				Name: "symbol",
				from: 5,
				to:   6,
			}, {
				Name: "symbol",
				from: 10,
				to:   11,
			}, {
				Name: "symbol",
				from: 15,
				to:   16,
			}},
		}},
	}, {
		msg:  "list, complex",
		text: "[a, b, c..., [d, e], [f, [g]]...]",
		nodes: []*Node{{
			Name: "list",
			to:   33,
			Nodes: []*Node{{
				Name: "symbol",
				from: 1,
				to:   2,
			}, {
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "spread-expression",
				from: 7,
				to:   11,
				Nodes: []*Node{{
					Name: "symbol",
					from: 7,
					to:   8,
				}},
			}, {
				Name: "list",
				from: 13,
				to:   19,
				Nodes: []*Node{{
					Name: "symbol",
					from: 14,
					to:   15,
				}, {
					Name: "symbol",
					from: 17,
					to:   18,
				}},
			}, {
				Name: "spread-expression",
				from: 21,
				to:   32,
				Nodes: []*Node{{
					Name: "list",
					from: 21,
					to:   29,
					Nodes: []*Node{{
						Name: "symbol",
						from: 22,
						to:   23,
					}, {
						Name: "list",
						from: 25,
						to:   28,
						Nodes: []*Node{{
							Name: "symbol",
							from: 26,
							to:   27,
						}},
					}},
				}},
			}},
		}},
	}, {
		msg:  "mutable list",
		text: "~[a, b, c]",
		nodes: []*Node{{
			Name: "mutable-list",
			to:   10,
			Nodes: []*Node{{
				Name: "symbol",
				from: 2,
				to:   3,
			}, {
				Name: "symbol",
				from: 5,
				to:   6,
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}},
		}},
	}, {
		msg:  "empty struct",
		text: "{}",
		nodes: []*Node{{
			Name: "struct",
			to:   2,
		}},
	}, {
		msg:  "struct",
		text: "{foo: 1, \"bar\": 2, symbol(baz): 3, [qux]: 4}",
		nodes: []*Node{{
			Name: "struct",
			to:   44,
			Nodes: []*Node{{
				Name: "entry",
				from: 1,
				to:   7,
				Nodes: []*Node{{
					Name: "symbol",
					from: 1,
					to:   4,
				}, {
					Name: "int",
					from: 6,
					to:   7,
				}},
			}, {
				Name: "entry",
				from: 9,
				to:   17,
				Nodes: []*Node{{
					Name: "string",
					from: 9,
					to:   14,
				}, {
					Name: "int",
					from: 16,
					to:   17,
				}},
			}, {
				Name: "entry",
				from: 19,
				to:   33,
				Nodes: []*Node{{
					Name: "dynamic-symbol",
					from: 19,
					to:   30,
					Nodes: []*Node{{
						Name: "symbol",
						from: 26,
						to:   29,
					}},
				}, {
					Name: "int",
					from: 32,
					to:   33,
				}},
			}, {
				Name: "entry",
				from: 35,
				to:   43,
				Nodes: []*Node{{
					Name: "indexer-symbol",
					from: 35,
					to:   40,
					Nodes: []*Node{{
						Name: "symbol",
						from: 36,
						to:   39,
					}},
				}, {
					Name: "int",
					from: 42,
					to:   43,
				}},
			}},
		}},
	}, {
		msg:  "struct, complex",
		text: "{foo: 1, {bar: 2}..., {baz: {}}...}",
		nodes: []*Node{{
			Name: "struct",
			to:   35,
			Nodes: []*Node{{
				Name: "entry",
				from: 1,
				to:   7,
				Nodes: []*Node{{
					Name: "symbol",
					from: 1,
					to:   4,
				}, {
					Name: "int",
					from: 6,
					to:   7,
				}},
			}, {
				Name: "spread-expression",
				from: 9,
				to:   20,
				Nodes: []*Node{{
					Name: "struct",
					from: 9,
					to:   17,
					Nodes: []*Node{{
						Name: "entry",
						from: 10,
						to:   16,
						Nodes: []*Node{{
							Name: "symbol",
							from: 10,
							to:   13,
						}, {
							Name: "int",
							from: 15,
							to:   16,
						}},
					}},
				}},
			}, {
				Name: "spread-expression",
				from: 22,
				to:   34,
				Nodes: []*Node{{
					Name: "struct",
					from: 22,
					to:   31,
					Nodes: []*Node{{
						Name: "entry",
						from: 23,
						to:   30,
						Nodes: []*Node{{
							Name: "symbol",
							from: 23,
							to:   26,
						}, {
							Name: "struct",
							from: 28,
							to:   30,
						}},
					}},
				}},
			}},
		}},
	}, {
		msg:  "struct with indexer key",
		text: "{[a]: b}",
		nodes: []*Node{{
			Name: "struct",
			to:   8,
			Nodes: []*Node{{
				Name: "entry",
				from: 1,
				to:   7,
				Nodes: []*Node{{
					Name: "indexer-symbol",
					from: 1,
					to:   4,
					Nodes: []*Node{{
						Name: "symbol",
						from: 2,
						to:   3,
					}},
				}, {
					Name: "symbol",
					from: 6,
					to:   7,
				}},
			}},
		}},
	}, {
		msg:  "mutable struct",
		text: "~{foo: 1}",
		nodes: []*Node{{
			Name: "mutable-struct",
			to:   9,
			Nodes: []*Node{{
				Name: "entry",
				from: 2,
				to:   8,
				Nodes: []*Node{{
					Name: "symbol",
					from: 2,
					to:   5,
				}, {
					Name: "int",
					from: 7,
					to:   8,
				}},
			}},
		}},
	}, {
		msg:  "channel",
		text: "<>",
		nodes: []*Node{{
			Name: "channel",
			to:   2,
		}},
	}, {
		msg:  "buffered channel",
		text: "<42>",
		nodes: []*Node{{
			Name: "channel",
			to:   4,
			Nodes: []*Node{{
				Name: "int",
				from: 1,
				to:   3,
			}},
		}},
	}, {
		msg:  "and expression",
		text: "and(a, b, c)",
		nodes: []*Node{{
			Name: "function-application",
			to:   12,
			Nodes: []*Node{{
				Name: "symbol",
				to:   3,
			}, {
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "symbol",
				from: 10,
				to:   11,
			}},
		}},
	}, {
		msg:  "or expression",
		text: "or(a, b, c)",
		nodes: []*Node{{
			Name: "function-application",
			to:   11,
			Nodes: []*Node{{
				Name: "symbol",
				to:   2,
			}, {
				Name: "symbol",
				from: 3,
				to:   4,
			}, {
				Name: "symbol",
				from: 6,
				to:   7,
			}, {
				Name: "symbol",
				from: 9,
				to:   10,
			}},
		}},
	}, {
		msg:  "function",
		text: "fn () 42",
		nodes: []*Node{{
			Name: "function",
			to:   8,
			Nodes: []*Node{{
				Name: "int",
				from: 6,
				to:   8,
			}},
		}},
	}, {
		msg:  "function, noop",
		text: "fn () {;}",
		nodes: []*Node{{
			Name: "function",
			to:   9,
			Nodes: []*Node{{
				Name: "block",
				from: 6,
				to:   9,
			}},
		}},
	}, {
		msg:  "function with args",
		text: "fn (a, b, c) [a, b, c]",
		nodes: []*Node{{
			Name: "function",
			to:   22,
			Nodes: []*Node{{
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "symbol",
				from: 10,
				to:   11,
			}, {
				Name: "list",
				from: 13,
				to:   22,
				Nodes: []*Node{{
					Name: "symbol",
					from: 14,
					to:   15,
				}, {
					Name: "symbol",
					from: 17,
					to:   18,
				}, {
					Name: "symbol",
					from: 20,
					to:   21,
				}},
			}},
		}},
	}, {
		msg: "function with args in new lines",
		text: `fn (
			a
			b
			c
		) [a, b, c]`,
		nodes: []*Node{{
			Name: "function",
			to:   33,
			Nodes: []*Node{{
				Name: "symbol",
				from: 8,
				to:   9,
			}, {
				Name: "symbol",
				from: 13,
				to:   14,
			}, {
				Name: "symbol",
				from: 18,
				to:   19,
			}, {
				Name: "list",
				from: 24,
				to:   33,
				Nodes: []*Node{{
					Name: "symbol",
					from: 25,
					to:   26,
				}, {
					Name: "symbol",
					from: 28,
					to:   29,
				}, {
					Name: "symbol",
					from: 31,
					to:   32,
				}},
			}},
		}},
	}, {
		msg:  "function with spread arg",
		text: "fn (a, b, ...c) [a, b, c]",
		nodes: []*Node{{
			Name: "function",
			to:   25,
			Nodes: []*Node{{
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "collect-symbol",
				from: 10,
				to:   14,
				Nodes: []*Node{{
					Name: "symbol",
					from: 13,
					to:   14,
				}},
			}, {
				Name: "list",
				from: 16,
				to:   25,
				Nodes: []*Node{{
					Name: "symbol",
					from: 17,
					to:   18,
				}, {
					Name: "symbol",
					from: 20,
					to:   21,
				}, {
					Name: "symbol",
					from: 23,
					to:   24,
				}},
			}},
		}},
	}, {
		msg:  "effect",
		text: "fn ~ () 42",
		nodes: []*Node{{
			Name: "effect",
			to:   10,
			Nodes: []*Node{{
				Name: "int",
				from: 8,
				to:   10,
			}},
		}},
	}, {
		msg:  "indexer",
		text: "a[42]",
		nodes: []*Node{{
			Name: "indexer",
			to:   5,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "int",
				from: 2,
				to:   4,
			}},
		}},
	}, {
		msg:  "range indexer",
		text: "a[3:9]",
		nodes: []*Node{{
			Name: "indexer",
			to:   6,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "range-from",
				from: 2,
				to:   3,
				Nodes: []*Node{{
					Name: "int",
					from: 2,
					to:   3,
				}},
			}, {
				Name: "range-to",
				from: 4,
				to:   5,
				Nodes: []*Node{{
					Name: "int",
					from: 4,
					to:   5,
				}},
			}},
		}},
	}, {
		msg:  "range indexer, lower unbound",
		text: "a[:9]",
		nodes: []*Node{{
			Name: "indexer",
			to:   5,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "range-to",
				from: 3,
				to:   4,
				Nodes: []*Node{{
					Name: "int",
					from: 3,
					to:   4,
				}},
			}},
		}},
	}, {
		msg:  "range indexer, upper unbound",
		text: "a[3:]",
		nodes: []*Node{{
			Name: "indexer",
			to:   5,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "range-from",
				from: 2,
				to:   3,
				Nodes: []*Node{{
					Name: "int",
					from: 2,
					to:   3,
				}},
			}},
		}},
	}, {
		msg:  "indexer, chained",
		text: "a[b][c][d]",
		nodes: []*Node{{
			Name: "indexer",
			to:   10,
			Nodes: []*Node{{
				Name: "indexer",
				to:   7,
				Nodes: []*Node{{
					Name: "indexer",
					to:   4,
					Nodes: []*Node{{
						Name: "symbol",
						to:   1,
					}, {
						Name: "symbol",
						from: 2,
						to:   3,
					}},
				}, {
					Name: "symbol",
					from: 5,
					to:   6,
				}},
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}},
		}},
	}, {
		msg:  "symbol indexer",
		text: "a.b",
		nodes: []*Node{{
			Name: "indexer",
			to:   3,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 2,
				to:   3,
			}},
		}},
	}, {
		msg:  "symbol indexer, with string",
		text: "a.\"b\"",
		nodes: []*Node{{
			Name: "indexer",
			to:   5,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "string",
				from: 2,
				to:   5,
			}},
		}},
	}, {
		msg:  "symbol indexer, with dynamic symbol",
		text: "a.symbol(b)",
		nodes: []*Node{{
			Name: "indexer",
			to:   11,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "dynamic-symbol",
				from: 2,
				to:   11,
				Nodes: []*Node{{
					Name: "symbol",
					from: 9,
					to:   10,
				}},
			}},
		}},
	}, {
		msg:  "chained symbol indexer",
		text: "a.b.c.d",
		nodes: []*Node{{
			Name: "indexer",
			to:   7,
			Nodes: []*Node{{
				Name: "indexer",
				to:   5,
				Nodes: []*Node{{
					Name: "indexer",
					to:   3,
					Nodes: []*Node{{
						Name: "symbol",
						to:   1,
					}, {
						Name: "symbol",
						from: 2,
						to:   3,
					}},
				}, {
					Name: "symbol",
					from: 4,
					to:   5,
				}},
			}, {
				Name: "symbol",
				from: 6,
				to:   7,
			}},
		}},
	}, {
		msg:  "chained symbol indexer on new line",
		text: "a\n.b\n.c",
		nodes: []*Node{{
			Name: "indexer",
			to:   7,
			Nodes: []*Node{{
				Name: "indexer",
				to:   4,
				Nodes: []*Node{{
					Name: "symbol",
					to:   1,
				}, {
					Name: "symbol",
					from: 3,
					to:   4,
				}},
			}, {
				Name: "symbol",
				from: 6,
				to:   7,
			}},
		}},
	}, {
		msg:  "chained symbol indexer on new line after dot",
		text: "a.\nb.\nc",
		nodes: []*Node{{
			Name: "indexer",
			to:   7,
			Nodes: []*Node{{
				Name: "indexer",
				to:   4,
				Nodes: []*Node{{
					Name: "symbol",
					to:   1,
				}, {
					Name: "symbol",
					from: 3,
					to:   4,
				}},
			}, {
				Name: "symbol",
				from: 6,
				to:   7,
			}},
		}},
	}, {
		msg:  "function application",
		text: "f()",
		nodes: []*Node{{
			Name: "function-application",
			to:   3,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}},
		}},
	}, {
		msg:  "function application, single arg",
		text: "f(a)",
		nodes: []*Node{{
			Name: "function-application",
			to:   4,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 2,
				to:   3,
			}},
		}},
	}, {
		msg:  "function application, multiple args",
		text: "f(a, b, c)",
		nodes: []*Node{{
			Name: "function-application",
			to:   10,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 2,
				to:   3,
			}, {
				Name: "symbol",
				from: 5,
				to:   6,
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}},
		}},
	}, {
		msg:  "function application, multiple args, new line",
		text: "f(a\nb\nc\n)",
		nodes: []*Node{{
			Name: "function-application",
			to:   9,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 2,
				to:   3,
			}, {
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 6,
				to:   7,
			}},
		}},
	}, {
		msg:  "function application, spread",
		text: "f(a, b..., c, d...)",
		nodes: []*Node{{
			Name: "function-application",
			to:   19,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 2,
				to:   3,
			}, {
				Name: "spread-expression",
				from: 5,
				to:   9,
				Nodes: []*Node{{
					Name: "symbol",
					from: 5,
					to:   6,
				}},
			}, {
				Name: "symbol",
				from: 11,
				to:   12,
			}, {
				Name: "spread-expression",
				from: 14,
				to:   18,
				Nodes: []*Node{{
					Name: "symbol",
					from: 14,
					to:   15,
				}},
			}},
		}},
	}, {
		msg:  "chained function application",
		text: "f(a)(b)(c)",
		nodes: []*Node{{
			Name: "function-application",
			to:   10,
			Nodes: []*Node{{
				Name: "function-application",
				to:   7,
				Nodes: []*Node{{
					Name: "function-application",
					to:   4,
					Nodes: []*Node{{
						Name: "symbol",
						to:   1,
					}, {
						Name: "symbol",
						from: 2,
						to:   3,
					}},
				}, {
					Name: "symbol",
					from: 5,
					to:   6,
				}},
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}},
		}},
	}, {
		msg:  "embedded function application",
		text: "f(g(h(a)))",
		nodes: []*Node{{
			Name: "function-application",
			to:   10,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "function-application",
				from: 2,
				to:   9,
				Nodes: []*Node{{
					Name: "symbol",
					from: 2,
					to:   3,
				}, {
					Name: "function-application",
					from: 4,
					to:   8,
					Nodes: []*Node{{
						Name: "symbol",
						from: 4,
						to:   5,
					}, {
						Name: "symbol",
						from: 6,
						to:   7,
					}},
				}},
			}},
		}},
	}, {
		msg:  "if",
		text: "if a { b() }",
		nodes: []*Node{{
			Name: "if",
			to:   12,
			Nodes: []*Node{{
				Name: "symbol",
				from: 3,
				to:   4,
			}, {
				Name: "block",
				from: 5,
				to:   12,
				Nodes: []*Node{{
					Name: "function-application",
					from: 7,
					to:   10,
					Nodes: []*Node{{
						Name: "symbol",
						from: 7,
						to:   8,
					}},
				}},
			}},
		}},
	}, {
		msg:  "if, else",
		text: "if a { b } else { c }",
		nodes: []*Node{{
			Name: "if",
			to:   21,
			Nodes: []*Node{{
				Name: "symbol",
				from: 3,
				to:   4,
			}, {
				Name: "block",
				from: 5,
				to:   10,
				Nodes: []*Node{{
					Name: "symbol",
					from: 7,
					to:   8,
				}},
			}, {
				Name: "block",
				from: 16,
				to:   21,
				Nodes: []*Node{{
					Name: "symbol",
					from: 18,
					to:   19,
				}},
			}},
		}},
	}, {
		msg: "if, else if, else if, else",
		text: `
			if a { b }
			else if c { d }
			else if e { f }
			else { g }
		`,
		nodes: []*Node{{
			Name: "if",
			from: 4,
			to:   66,
			Nodes: []*Node{{
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "block",
				from: 9,
				to:   14,
				Nodes: []*Node{{
					Name: "symbol",
					from: 11,
					to:   12,
				}},
			}, {
				Name: "symbol",
				from: 26,
				to:   27,
			}, {
				Name: "block",
				from: 28,
				to:   33,
				Nodes: []*Node{{
					Name: "symbol",
					from: 30,
					to:   31,
				}},
			}, {
				Name: "symbol",
				from: 45,
				to:   46,
			}, {
				Name: "block",
				from: 47,
				to:   52,
				Nodes: []*Node{{
					Name: "symbol",
					from: 49,
					to:   50,
				}},
			}, {
				Name: "block",
				from: 61,
				to:   66,
				Nodes: []*Node{{
					Name: "symbol",
					from: 63,
					to:   64,
				}},
			}},
		}},
	}, {
		msg:  "switch, empty",
		text: "switch {default:}",
		nodes: []*Node{{
			Name: "switch",
			to:   17,
			Nodes: []*Node{{
				Name: "default",
				from: 8,
				to:   16,
			}},
		}},
	}, {
		msg:  "switch, single case",
		text: "switch a {case b: c}",
		nodes: []*Node{{
			Name: "switch",
			to:   20,
			Nodes: []*Node{{
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "case",
				from: 10,
				to:   17,
				Nodes: []*Node{{
					Name: "symbol",
					from: 15,
					to:   16,
				}},
			}, {
				Name: "symbol",
				from: 18,
				to:   19,
			}},
		}},
	}, {
		msg:  "switch",
		text: "switch a {case b: c; case d: e; default: f}",
		nodes: []*Node{{
			Name: "switch",
			to:   43,
			Nodes: []*Node{{
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "case",
				from: 10,
				to:   17,
				Nodes: []*Node{{
					Name: "symbol",
					from: 15,
					to:   16,
				}},
			}, {
				Name: "symbol",
				from: 18,
				to:   19,
			}, {
				Name: "case",
				from: 21,
				to:   28,
				Nodes: []*Node{{
					Name: "symbol",
					from: 26,
					to:   27,
				}},
			}, {
				Name: "symbol",
				from: 29,
				to:   30,
			}, {
				Name: "default",
				from: 32,
				to:   40,
			}, {
				Name: "symbol",
				from: 41,
				to:   42,
			}},
		}},
	}, {
		msg: "switch, all new lines",
		text: `switch
			a
			{
			case
			b
			:
			c
			case
			d
			:
			e
			default
			:
			f
		}`,
		nodes: []*Node{{
			Name: "switch",
			to:   87,
			Nodes: []*Node{{
				Name: "symbol",
				from: 10,
				to:   11,
			}, {
				Name: "case",
				from: 20,
				to:   34,
				Nodes: []*Node{{
					Name: "symbol",
					from: 28,
					to:   29,
				}},
			}, {
				Name: "symbol",
				from: 38,
				to:   39,
			}, {
				Name: "case",
				from: 43,
				to:   57,
				Nodes: []*Node{{
					Name: "symbol",
					from: 51,
					to:   52,
				}},
			}, {
				Name: "symbol",
				from: 61,
				to:   62,
			}, {
				Name: "default",
				from: 66,
				to:   78,
			}, {
				Name: "symbol",
				from: 82,
				to:   83,
			}},
		}},
	}, {
		msg:  "match expression, empty",
		text: "match a {}",
		nodes: []*Node{{
			Name: "match",
			to:   10,
			Nodes: []*Node{{
				Name: "symbol",
				from: 6,
				to:   7,
			}},
		}},
	}, {
		msg: "match expression",
		text: `match a {
			case [first, ...rest]: first
		}`,
		nodes: []*Node{{
			Name: "match",
			to:   45,
			Nodes: []*Node{{
				Name: "symbol",
				from: 6,
				to:   7,
			}, {
				Name: "match-case",
				from: 13,
				to:   35,
				Nodes: []*Node{{
					Name: "list-type",
					from: 18,
					to:   34,
					Nodes: []*Node{{
						Name: "list-destructure-type",
						from: 19,
						to:   33,
						Nodes: []*Node{{
							Name: "destructure-item",
							from: 19,
							to:   24,
							Nodes: []*Node{{
								Name: "symbol",
								from: 19,
								to:   24,
							}},
						}, {
							Name: "collect-destructure-item",
							from: 26,
							to:   33,
							Nodes: []*Node{{
								Name: "destructure-item",
								from: 29,
								to:   33,
								Nodes: []*Node{{
									Name: "symbol",
									from: 29,
									to:   33,
								}},
							}},
						}},
					}},
				}},
			}, {
				Name: "symbol",
				from: 36,
				to:   41,
			}},
		}},
	}, {
		msg: "match expression, multiple cases",
		text: `match a {
			case [0]: []
			case [2:]: a[2:]
			default: error("invalid length")
		}`,
		nodes: []*Node{{
			Name: "match",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "match-case",
				Nodes: []*Node{{
					Name: "list-type",
					Nodes: []*Node{{
						Name: "items-type",
						Nodes: []*Node{{
							Name: "items-quantifier",
							Nodes: []*Node{{
								Name: "int",
							}},
						}},
					}},
				}},
			}, {
				Name: "list",
			}, {
				Name: "match-case",
				Nodes: []*Node{{
					Name: "list-type",
					Nodes: []*Node{{
						Name: "items-type",
						Nodes: []*Node{{
							Name: "items-quantifier",
							Nodes: []*Node{{
								Name: "static-range-from",
								Nodes: []*Node{{
									Name: "int",
								}},
							}},
						}},
					}},
				}},
			}, {
				Name: "indexer",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "range-from",
					Nodes: []*Node{{
						Name: "int",
					}},
				}},
			}, {
				Name: "default",
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "string",
				}},
			}},
		}},
		ignorePosition: true,
	}, {
		msg: "match function",
		text: `match a {
			case fn () int: a()
			default: 42
		}`,
		nodes: []*Node{{
			Name: "match",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "match-case",
				Nodes: []*Node{{
					Name: "function-type",
					Nodes: []*Node{{
						Name: "int-type",
					}},
				}},
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "default",
			}, {
				Name: "int",
			}},
		}},
		ignorePosition: true,
	}, {
		msg: "match expression, combined",
		text: `match a {
			case [fn (int)]: a[0]()
			default: 42
		}`,
		nodes: []*Node{{
			Name: "match",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "match-case",
				Nodes: []*Node{{
					Name: "list-type",
					Nodes: []*Node{{
						Name: "items-type",
						Nodes: []*Node{{
							Name: "function-type",
							Nodes: []*Node{{
								Name: "arg-type",
								Nodes: []*Node{{
									Name: "int-type",
								}},
							}},
						}},
					}},
				}},
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "indexer",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "int",
					}},
				}},
			}, {
				Name: "default",
			}, {
				Name: "int",
			}},
		}},
		ignorePosition: true,
	}, {
		msg: "match expression, complex",
		text: `match a {
				case [first T int|string, op fn ([T, int, ...T]) int, ...rest T]:
					op([first, now(), rest...])
				default:
					error("invalid list")
			}`,
		nodes: []*Node{{
			Name: "match",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "match-case",
				Nodes: []*Node{{
					Name: "list-match",
					Nodes: []*Node{{
						Name: "list-destructure-match",
						Nodes: []*Node{{
							Name: "destructure-match-item",
							Nodes: []*Node{{
								Name: "symbol",
							}, {
								Name: "symbol",
							}, {
								Name: "int-type",
							}, {
								Name: "string-type",
							}},
						}, {
							Name: "destructure-match-item",
							Nodes: []*Node{{
								Name: "symbol",
							}, {
								Name: "function-type",
								Nodes: []*Node{{
									Name: "arg-type",
									Nodes: []*Node{{
										Name: "list-type",
										Nodes: []*Node{{
											Name: "list-destructure-type",
											Nodes: []*Node{{
												Name: "destructure-item",
												Nodes: []*Node{{
													Name: "symbol",
												}},
											}, {
												Name: "destructure-item",
												Nodes: []*Node{{
													Name: "int-type",
												}},
											}, {
												Name: "collect-destructure-item",
												Nodes: []*Node{{
													Name: "destructure-item",
													Nodes: []*Node{{
														Name: "symbol",
													}},
												}},
											}},
										}},
									}},
								}, {
									Name: "int-type",
								}},
							}},
						}, {
							Name: "collect-destructure-match-item",
							Nodes: []*Node{{
								Name: "destructure-match-item",
								Nodes: []*Node{{
									Name: "symbol",
								}, {
									Name: "symbol",
								}},
							}},
						}},
					}},
				}},
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "list",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "function-application",
						Nodes: []*Node{{
							Name: "symbol",
						}},
					}, {
						Name: "spread-expression",
						Nodes: []*Node{{
							Name: "symbol",
						}},
					}},
				}},
			}, {
				Name: "default",
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "string",
				}},
			}},
		}},
		ignorePosition: true,
	}, {
		msg:  "receive op",
		text: "<-chan",
		nodes: []*Node{{
			Name: "receive-op",
			Nodes: []*Node{{
				Name: "symbol",
			}},
		}},
		ignorePosition: true,
	}, {
		msg:  "send op",
		text: "chan <- a",
		nodes: []*Node{{
			Name: "send",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "symbol",
			}},
		}},
		ignorePosition: true,
	}, {
		msg: "select, empty",
		text: `select {
		}`,
		nodes: []*Node{{
			Name: "select",
			to:   12,
		}},
	}, {
		msg: "select",
		text: `select {
			case let a <-r: s <- a
			case s <- f(): g()
			default: h()
		}`,
		nodes: []*Node{{
			Name: "select",
			Nodes: []*Node{{
				Name: "select-case",
				Nodes: []*Node{{
					Name: "receive-definition",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "receive-op",
						Nodes: []*Node{{
							Name: "symbol",
						}},
					}},
				}},
			}, {
				Name: "send",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "symbol",
				}},
			}, {
				Name: "select-case",
				Nodes: []*Node{{
					Name: "send",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "function-application",
						Nodes: []*Node{{
							Name: "symbol",
						}},
					}},
				}},
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "default",
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}},
		}},
		ignorePosition: true,
	}, {
		msg: "select, call",
		text: `select {
			case let a receive(r): f()
			case send(s, g()): h()
			default: i()
		}`,
		nodes: []*Node{{
			Name: "select",
			Nodes: []*Node{{
				Name: "select-case",
				Nodes: []*Node{{
					Name: "receive-definition",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "receive-call",
						Nodes: []*Node{{
							Name: "symbol",
						}},
					}},
				}},
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "select-case",
				Nodes: []*Node{{
					Name: "send",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "function-application",
						Nodes: []*Node{{
							Name: "symbol",
						}},
					}},
				}},
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "default",
			}, {
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}},
		}},
		ignorePosition: true,
	}, {
		msg:  "go",
		text: "go f()",
		nodes: []*Node{{
			Name: "go",
			Nodes: []*Node{{
				Name: "function-application",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}},
		}},
		ignorePosition: true,
	}, {
		msg:  "ternary expression",
		text: "a ? b : c",
		nodes: []*Node{{
			Name: "ternary-expression",
			to:   9,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}},
		}},
	}, {
		msg:  "multiple ternary expressions, consequence",
		text: "a ? b ? c : d : e",
		nodes: []*Node{{
			Name: "ternary-expression",
			to:   17,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "ternary-expression",
				from: 4,
				to:   13,
				Nodes: []*Node{{
					Name: "symbol",
					from: 4,
					to:   5,
				}, {
					Name: "symbol",
					from: 8,
					to:   9,
				}, {
					Name: "symbol",
					from: 12,
					to:   13,
				}},
			}, {
				Name: "symbol",
				from: 16,
				to:   17,
			}},
		}},
	}, {
		msg:  "multiple ternary expressions, alternative",
		text: "a ? b : c ? d : e",
		nodes: []*Node{{
			Name: "ternary-expression",
			to:   17,
			Nodes: []*Node{{
				Name: "symbol",
				to:   1,
			}, {
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "ternary-expression",
				from: 8,
				to:   17,
				Nodes: []*Node{{
					Name: "symbol",
					from: 8,
					to:   9,
				}, {
					Name: "symbol",
					from: 12,
					to:   13,
				}, {
					Name: "symbol",
					from: 16,
					to:   17,
				}},
			}},
		}},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			n, err := s.Parse(bytes.NewBufferString(ti.text))

			if ti.fail && err == nil {
				t.Error("failed to fail")
				return
			} else if !ti.fail && err != nil {
				t.Error(err)
				return
			} else if ti.fail {
				return
			}

			t.Log(n)

			cn := checkNode
			if ti.ignorePosition {
				cn = checkNodeIgnorePosition
			}

			if ti.node != nil {
				cn(t, n, ti.node)
			} else {
				cn(t, n, &Node{
					Name:  "mml",
					from:  0,
					to:    len(ti.text),
					Nodes: ti.nodes,
				})
			}
		})
	}
}
