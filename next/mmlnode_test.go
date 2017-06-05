package next

import (
	"bytes"
	"os"
	"testing"
)

func TestMMLNodes(t *testing.T) {
	sb, err := defineSyntax()
	if err != nil {
		t.Error(err)
		return
	}

	def, err := os.Open("syntax.p")
	if err != nil {
		t.Error(err)
		return
	}

	defer def.Close()

	n, err := sb.Parse(def)
	if err != nil {
		t.Error("error parsing syntax", err)
		return
	}

	ss, err := defineDocument(n)
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

	n, err = ss.Parse(mml)
	if err != nil {
		t.Error("error parsing mml syntax", err)
		return
	}

	tl := TraceOff
	// tl = TraceDebug
	s, err := defineDocumentTrace(n, tl)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Init()
	if err != nil {
		t.Error(err)
		return
	}

	for _, ti := range []struct {
		msg   string
		text  string
		fail  bool
		node  *Node
		nodes []*Node
	}{{
		msg:  "empty",
		node: &Node{Name: "mml"},
	}, {
		msg:  "single line comment",
		text: "// foo bar baz",
		nodes: []*Node{{
			Name: "comment",
			from: 0,
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
			from: 0,
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
			from: 0,
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
			from: 0,
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
			from: 0,
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
			from: 0,
			to:   2,
		}},
	}, {
		msg:  "ints",
		text: "1; 2; 3",
		nodes: []*Node{{
			Name: "int",
			from: 0,
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
			from: 0,
			to:   3,
		}},
	}, {
		msg:  "int, hexa",
		text: "0x2a",
		nodes: []*Node{{
			Name: "int",
			from: 0,
			to:   4,
		}},
	}, {
		msg:  "float, 0.",
		text: "0.",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   2,
		}},
	}, {
		msg:  "float, 72.40",
		text: "72.40",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   5,
		}},
	}, {
		msg:  "float, 072.40",
		text: "072.40",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   6,
		}},
	}, {
		msg:  "float, 2.71828",
		text: "2.71828",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   7,
		}},
	}, {
		msg:  "float, 6.67428e-11",
		text: "6.67428e-11",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   11,
		}},
	}, {
		msg:  "float, 1E6",
		text: "1E6",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   3,
		}},
	}, {
		msg:  "float, .25",
		text: ".25",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   3,
		}},
	}, {
		msg:  "float, .12345E+5",
		text: ".12345E+5",
		nodes: []*Node{{
			Name: "float",
			from: 0,
			to:   9,
		}},
	}, {
		msg:  "string, empty",
		text: "\"\"",
		nodes: []*Node{{
			Name: "string",
			from: 0,
			to:   2,
		}},
	}, {
		msg:  "string",
		text: "\"foo\"",
		nodes: []*Node{{
			Name: "string",
			from: 0,
			to:   5,
		}},
	}, {
		msg:  "string, with new line",
		text: "\"foo\nbar\"",
		nodes: []*Node{{
			Name: "string",
			from: 0,
			to:   9,
		}},
	}, {
		msg:  "string, with escaped new line",
		text: "\"foo\\nbar\"",
		nodes: []*Node{{
			Name: "string",
			from: 0,
			to:   10,
		}},
	}, {
		msg:  "string, with quotes",
		text: "\"foo \\\"bar\\\" baz\"",
		nodes: []*Node{{
			Name: "string",
			from: 0,
			to:   17,
		}},
	}, {
		msg:  "bool, true",
		text: "true",
		nodes: []*Node{{
			Name: "true",
			from: 0,
			to:   4,
		}},
	}, {
		msg:  "bool, false",
		text: "false",
		nodes: []*Node{{
			Name: "false",
			from: 0,
			to:   5,
		}},
	}, {
		msg:  "symbol",
		text: "foo",
		nodes: []*Node{{
			Name: "symbol",
			from: 0,
			to:   3,
		}},
	}, {
		msg:  "dynamic-symbol",
		text: "symbol(a)",
		nodes: []*Node{{
			Name: "dynamic-symbol",
			from: 0,
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
			from: 0,
			to:   2,
		}},
	}, {
		msg:  "list",
		text: "[a, b, c]",
		nodes: []*Node{{
			Name: "list",
			from: 0,
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
			from: 0,
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
			from: 0,
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
			from: 0,
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
			from: 0,
			to:   2,
		}},
	}, {
		msg:  "struct",
		text: "{foo: 1, \"bar\": 2, symbol(baz): 3, [qux]: 4}",
		nodes: []*Node{{
			Name: "struct",
			from: 0,
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
			from: 0,
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
		msg:  "mutable struct",
		text: "~{foo: 1}",
		nodes: []*Node{{
			Name: "mutable-struct",
			from: 0,
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
			from: 0,
			to:   2,
		}},
	}, {
		msg:  "buffered channel",
		text: "<42>",
		nodes: []*Node{{
			Name: "channel",
			from: 0,
			to:   4,
			Nodes: []*Node{{
				Name: "int",
				from: 1,
				to:   3,
			}},
		}},
	}, {
		msg:  "function",
		text: "fn () 42",
		nodes: []*Node{{
			Name: "function",
			from: 0,
			to:   8,
			Nodes: []*Node{{
				Name: "int",
				from: 6,
				to:   8,
			}},
		}},
	}, {
		msg:  "function with args",
		text: "fn (a, b, c) [a, b, c]",
		nodes: []*Node{{
			Name: "function",
			from: 0,
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
			from: 0,
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
		msg:  "function with collect arg",
		text: "fn (a, b, ...c) [a, b, c]",
		nodes: []*Node{{
			Name: "function",
			from: 0,
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
			from: 0,
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
			from: 0,
			to:   5,
			Nodes: []*Node{{
				Name: "symbol",
				from: 0,
				to:   1,
			}, {
				Name: "int",
				from: 2,
				to:   4,
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

			if ti.node != nil {
				checkNode(t, n, ti.node)
			} else {
				checkNode(t, n, &Node{
					Name:  "mml",
					from:  0,
					to:    len(ti.text),
					Nodes: ti.nodes,
				})
			}
		})
	}
}
