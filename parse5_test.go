package mml

import (
	"bytes"
	"testing"
)

func checkTokens(left, right *token) bool {
	if (left == nil) != (right == nil) {
		return false
	}

	if left == nil {
		return true
	}

	return left.value == right.value
}

func checkNodes(left, right *node) bool {
	if (left == nil) != (right == nil) {
		return false
	}

	if left == nil {
		return true
	}

	if left.typeName != right.typeName {
		return false
	}

	if !checkTokens(left.token, right.token) {
		return false
	}

	if len(left.nodes) != len(right.nodes) {
		return false
	}

	for i, n := range left.nodes {
		if !checkNodes(n, right.nodes[i]) {
			return false
		}
	}

	return true
}

func def(f ...func(s *syntax) error) func(s *syntax) error {
	return func(s *syntax) error {
		for _, fi := range f {
			if err := fi(s); err != nil {
				return err
			}
		}

		return nil
	}
}

func TestParse(t *testing.T) {
	for _, ti := range []struct {
		msg    string
		syntax func(s *syntax) error
		text   string
		node   *node
		fail   bool
	}{{
		msg:    "int",
		syntax: func(s *syntax) error { return s.primitive("int", intToken) },
		text:   "42",
		node: &node{
			typeName: "int",
			token:    &token{value: "42"},
		},
	}, {
		msg:    "int, with empty input",
		syntax: func(s *syntax) error { return s.primitive("int", intToken) },
		fail:   true,
	}, {
		msg: "optional int",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
		),
		text: "42",
		node: &node{
			typeName: "int",
			token:    &token{value: "42"},
		},
	}, {
		msg: "optional int, empty",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
		),
		node: zeroNode,
	}, {
		msg: "optional int, not int",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
		),
		text: "\"foo\"",
		fail: true,
	}, {
		msg: "int sequence, optional",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.sequence("int-sequence", "int") },
			func(s *syntax) error { return s.optional("optional-int-sequence", "int-sequence") },
		),
		text: "1 2 3",
		node: &node{
			typeName: "int-sequence",
			token:    &token{value: "1"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "1"},
			}, {
				typeName: "int",
				token:    &token{value: "2"},
			}, {
				typeName: "int",
				token:    &token{value: "3"},
			}},
		},
	}, {
		msg: "int sequence, optional, empty",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.sequence("int-sequence", "int") },
			func(s *syntax) error { return s.optional("optional-int-sequence", "int-sequence") },
		),
		node: &node{
			typeName: "int-sequence",
			token:    eofToken,
		},
	}, {
		msg: "empty sequence",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.sequence("int-sequence", "int") },
		),
		node: &node{
			typeName: "int-sequence",
			token:    eofToken,
		},
	}, {
		msg: "sequence with a single item",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.sequence("int-sequence", "int") },
		),
		text: "42",
		node: &node{
			typeName: "int-sequence",
			token:    &token{value: "42"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "42"},
			}},
		},
	}, {
		msg: "sequence with multiple items",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.sequence("int-sequence", "int") },
		),
		text: "1 2 3",
		node: &node{
			typeName: "int-sequence",
			token:    &token{value: "1"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "1"},
			}, {
				typeName: "int",
				token:    &token{value: "2"},
			}, {
				typeName: "int",
				token:    &token{value: "3"},
			}},
		},
	}, {
		msg: "sequence with optional item",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error { return s.sequence("int-sequence", "optional-int") },
		),
		text: "42",
		node: &node{
			typeName: "int-sequence",
			token:    &token{value: "42"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "42"},
			}},
		},
	}, {
		msg: "sequence with multiple optional items",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error { return s.sequence("int-sequence", "optional-int") },
		),
		text: "1 2 3",
		node: &node{
			typeName: "int-sequence",
			token:    &token{value: "1"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "1"},
			}, {
				typeName: "int",
				token:    &token{value: "2"},
			}, {
				typeName: "int",
				token:    &token{value: "3"},
			}},
		},
	}, {
		msg: "group with single int",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.group("int-group", "int") },
		),
		text: "42",
		node: &node{
			typeName: "int-group",
			token:    &token{value: "42"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "42"},
			}},
		},
	}, {
		msg: "group with single optional int",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error { return s.group("int-group", "optional-int") },
		),
		text: "42",
		node: &node{
			typeName: "int-group",
			token:    &token{value: "42"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "42"},
			}},
		},
	}, {
		msg: "group with single int, not int",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.group("int-group", "int") },
		),
		text: "\"foo\"",
		fail: true,
	}, {
		msg: "group with multiple ints",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.group("int-group", "int", "int", "int") },
		),
		text: "1 2 3",
		node: &node{
			typeName: "int-group",
			token:    &token{value: "1"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "1"},
			}, {
				typeName: "int",
				token:    &token{value: "2"},
			}, {
				typeName: "int",
				token:    &token{value: "3"},
			}},
		},
	}, {
		msg: "group with optional item",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.primitive("string", stringToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error {
				return s.group("group-with-optional", "optional-int", "string")
			},
		),
		text: "42 \"foo\"",
		node: &node{
			typeName: "group-with-optional",
			token:    &token{value: "42"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "42"},
			}, {
				typeName: "string",
				token:    &token{value: "\"foo\""},
			}},
		},
	}, {
		msg: "group with optional item, missing",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.primitive("string", stringToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error {
				return s.group("group-with-optional", "optional-int", "string")
			},
		),
		text: "\"foo\"",
		node: &node{
			typeName: "group-with-optional",
			token:    &token{value: "\"foo\""},
			nodes: []*node{{
				typeName: "string",
				token:    &token{value: "\"foo\""},
			}},
		},
	}, {
		msg: "union of int and string",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.primitive("string", stringToken) },
			func(s *syntax) error { return s.union("int-or-string", "int", "string") },
		),
		text: "\"foo\"",
		node: &node{
			typeName: "string",
			token:    &token{value: "\"foo\""},
		},
	}, {
		msg: "union of int and group with optional int",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.primitive("string", stringToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error { return s.group("group-with-optional", "optional-int", "string") },
			func(s *syntax) error {
				return s.union("int-or-group-with-optional", "int", "group-with-optional")
			},
		),
		text: "42 \"foo\"",
		node: &node{
			typeName: "group-with-optional",
			token:    &token{value: "42"},
			nodes: []*node{{
				typeName: "int",
				token:    &token{value: "42"},
			}, {
				typeName: "string",
				token:    &token{value: "\"foo\""},
			}},
		},
	}, {
		msg: "union of int and group with optional int, token fall through",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.primitive("string", stringToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error {
				return s.group(
					"group-with-optional",
					"optional-int",
					"optional-int",
					"string",
					"string",
				)
			},
			func(s *syntax) error {
				return s.union("int-or-group-with-optional", "int", "group-with-optional")
			},
		),
		text: "\"foo\" \"bar\"",
		node: &node{
			typeName: "group-with-optional",
			token:    &token{value: "\"foo\""},
			nodes: []*node{{
				typeName: "string",
				token:    &token{value: "\"foo\""},
			}, {
				typeName: "string",
				token:    &token{value: "\"bar\""},
			}},
		},
	}, {
		msg: "union of int and group with optional int, init fall through",
		syntax: def(
			func(s *syntax) error { return s.primitive("int", intToken) },
			func(s *syntax) error { return s.primitive("string", stringToken) },
			func(s *syntax) error { return s.optional("optional-int", "int") },
			func(s *syntax) error {
				return s.group(
					"group-with-optional",
					"optional-int",
					"optional-int",
					"string",
					"string",
				)
			},
			func(s *syntax) error {
				return s.union("int-or-group-with-optional", "string", "group-with-optional")
			},
		),
		text: "\"foo\" \"bar\"",
		node: &node{
			typeName: "group-with-optional",
			token:    &token{value: "\"foo\""},
			nodes: []*node{{
				typeName: "string",
				token:    &token{value: "\"foo\""},
			}, {
				typeName: "string",
				token:    &token{value: "\"bar\""},
			}},
		},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			s := newSyntax()
			s.traceLevel = traceDebug
			if err := ti.syntax(s); err != nil {
				t.Error(err)
				return
			}

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

			if !checkNodes(n, ti.node) {
				t.Error("failed to match nodes", n, ti.node)
			}
		})
	}
}
