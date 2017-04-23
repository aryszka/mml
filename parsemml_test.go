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
		text: "1 2\n3;4 ;\n 5",
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
		}, {
			typeName: "semicolon",
			token:    &token{value: ";"},
		}, {
			typeName: "int",
			token:    &token{value: "4"},
		}, {
			typeName: "semicolon",
			token:    &token{value: ";"},
		}, {
			typeName: "nl",
			token:    &token{value: "\n"},
		}, {
			typeName: "int",
			token:    &token{value: "5"},
		}},
	}, {
		msg:  "string",
		text: "\"foo\"",
		nodes: []*node{{
			typeName: "string",
			token:    &token{value: "\"foo\""},
		}},
	}, {
		msg:  "symbol",
		text: "foo",
		nodes: []*node{{
			typeName: "symbol",
			token:    &token{value: "foo"},
		}},
	}, {
		msg:  "dynamic symbol",
		text: "symbol(a)",
		nodes: []*node{{
			typeName: "dynamic-symbol",
			token:    &token{value: "symbol"},
			nodes: []*node{{
				typeName: "symbol-word",
				token:    &token{value: "symbol"},
			}, {
				typeName: "open-paren",
				token:    &token{value: "("},
			}, {
				typeName: "nls",
				token:    &token{value: "a"},
			}, {
				typeName: "symbol",
				token:    &token{value: "a"},
			}, {
				typeName: "nls",
				token:    &token{value: ")"},
			}, {
				typeName: "close-paren",
				token:    &token{value: ")"},
			}},
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
				return
			}

			for i, ni := range n.nodes {
				if !checkNodes(ni, ti.nodes[i]) {
					t.Error("failed to match nodes", ni, ti.nodes[i])
				}
			}
		})
	}
}
