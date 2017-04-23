package mml

import (
	"bytes"
	"testing"
)

func TestParseMML(t *testing.T) {
	s, err := newMMLSyntax()
	if err != nil {
		t.Error(err)
		return
	}

	s.traceLevel = traceDebug

	for _, ti := range []struct {
		msg   string
		text  string
		nodes []*node
		fail  bool
	}{{
		msg: "empty document",
	}, {
		msg:  "single int",
		text: "42",
		nodes: []*node{{
			typeName: "int",
			token:    &token{value: "42"},
		}},
	}, {
		msg:  "multiple ints",
		text: "1 2\n3",
		nodes: []*node{{
			typeName: "int",
			token:    &token{value: "1"},
		}, {
			typeName: "int",
			token:    &token{value: "2"},
		}, {
			typeName: "nl",
			token:    &token{value: "\n"},
		}, {
			typeName: "int",
			token:    &token{value: "3"},
		}},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			b := bytes.NewBufferString(ti.text)
			r := newTokenReader(b, "<test>")

			n, err := s.parse(r)
			if !ti.fail && err != nil {
				t.Error(err)
				return
			} else if ti.fail && err == nil {
				t.Error("failed to fail")
				return
			}

			if ti.fail {
				return
			}

			if n.typeName != "statement-sequence" {
				t.Error("invalid root node type", n.typeName, "statement-sequence")
				return
			}

			if len(n.nodes) != len(ti.nodes) {
				t.Error("invalid number of nodes", len(n.nodes), len(ti.nodes))
				return
			}

			if len(n.nodes) == 0 && n.token != eofToken || len(n.nodes) > 0 && n.token != n.nodes[0].token {
				t.Error("invalid document token", n.token)
			}

			for i, ni := range n.nodes {
				if !checkNodes(ni, ti.nodes[i]) {
					t.Error("failed to match nodes", n, ti.nodes[i])
				}
			}
		})
	}
}
