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
	tl = TraceDebug
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
