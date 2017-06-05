package next

import (
	"bytes"
	"os"
	"testing"
)

func TestNodes(t *testing.T) {
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
		t.Error(err)
		return
	}

	s, err := defineDocument(n)
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
		node: &Node{Name: "document"},
	}, {
		msg:  "single line comment",
		text: "// foo bar baz",
		nodes: []*Node{{
			Name: "comment",
			from: 0,
			to:   14,
		}},
	}, {
		msg:  "multiple line comments",
		text: "// foo bar\n// baz qux",
		nodes: []*Node{{
			Name: "comment",
			from: 0,
			to:   21,
		}},
	}, {
		msg:  "block comment",
		text: "/* foo bar baz */",
		nodes: []*Node{{
			Name: "comment",
			from: 0,
			to:   17,
		}},
	}, {
		msg:  "block comments",
		text: "/* foo bar */\n/* baz qux */",
		nodes: []*Node{{
			Name: "comment",
			from: 0,
			to:   27,
		}},
	}, {
		msg:  "mixed comments",
		text: "// foo\n/* bar */\n// baz",
		nodes: []*Node{{
			Name: "comment",
			from: 0,
			to:   23,
		}},
	}, {
		msg:  "any char",
		text: "any = .",
		nodes: []*Node{{
			Name: "definition",
			from: 0,
			to:   7,
			Nodes: []*Node{{
				Name: "symbol",
				from: 0,
				to:   3,
			}, {
				Name: "any-char",
				from: 6,
				to:   7,
			}},
		}},
	}, {
		msg:  "char class",
		text: "char-class = [^abci-k]",
		nodes: []*Node{{
			Name: "definition",
			from: 0,
			to:   22,
			Nodes: []*Node{{
				Name: "symbol",
				from: 0,
				to:   10,
			}, {
				Name: "char-class",
				from: 13,
				to:   22,
				Nodes: []*Node{{
					Name: "class-not",
					from: 14,
					to:   15,
				}, {
					Name: "class-char",
					from: 15,
					to:   16,
				}, {
					Name: "class-char",
					from: 16,
					to:   17,
				}, {
					Name: "class-char",
					from: 17,
					to:   18,
				}, {
					Name: "char-range",
					from: 18,
					to:   21,
					Nodes: []*Node{{
						Name: "class-char",
						from: 18,
						to:   19,
					}, {
						Name: "class-char",
						from: 20,
						to:   21,
					}},
				}},
			}},
		}},
	}, {
		msg:  "char sequence",
		text: "char-sequence = \"foo\"",
		nodes: []*Node{{
			Name: "definition",
			from: 0,
			to:   21,
			Nodes: []*Node{{
				Name: "symbol",
				from: 0,
				to:   13,
			}, {
				Name: "char-sequence",
				from: 16,
				to:   21,
				Nodes: []*Node{{
					Name: "sequence-char",
					from: 17,
					to:   18,
				}, {
					Name: "sequence-char",
					from: 18,
					to:   19,
				}, {
					Name: "sequence-char",
					from: 19,
					to:   20,
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

			if ti.node != nil {
				checkNode(t, n, ti.node)
			} else {
				checkNodes(t, n.Nodes, ti.nodes)
			}
		})
	}
}
