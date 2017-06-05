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
