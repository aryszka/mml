package next

import (
	"bytes"
	"testing"
)

func TestChoice(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "choice",
		text: "a = a | b",
		node: &Node{
			Name: "choice",
			from: 4,
			to:   9,
			Nodes: []*Node{{
				Name: "symbol",
				from: 4,
				to:   5,
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}},
		},
	}, {
		msg:  "multiple",
		text: "abcd = a | b c | d",
		node: &Node{
			Name: "choice",
			from: 7,
			to:   18,
			Nodes: []*Node{{
				Name: "symbol",
				from: 7,
				to:   8,
			}, {
				Name: "sequence",
				from: 11,
				to:   14,
				Nodes: []*Node{{
					Name: "symbol",
					from: 11,
					to:   12,
				}, {
					Name: "symbol",
					from: 13,
					to:   14,
				}},
			}, {
				Name: "symbol",
				from: 17,
				to:   18,
			}},
		},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			s, err := defineSyntax()
			if err != nil {
				t.Error(err)
				return
			}

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

			checkNode(t, n.Nodes[0].Nodes[1], ti.node)
		})
	}
}
