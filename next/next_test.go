package next

import (
	"bytes"
	"testing"
)

type syntaxTest struct {
	msg    string
	syntax [][]string
	text   string
	node   *Node
	fail   bool
}

func checkNode(t *testing.T, left, right *Node) {
	if (left == nil) != (right == nil) {
		t.Error("nil reference doesn't match", left == nil, right == nil)
		return
	}

	if left == nil {
		return
	}

	if left.Name != right.Name {
		t.Error("name doesn't match", left.Name, right.Name)
		return
	}

	if left.from != right.from {
		t.Error("from doesn't match", left.Name, left.from, right.from)
		return
	}

	if left.to != right.to {
		t.Error("to doesn't match", left.Name, left.to, right.to)
		return
	}

	if len(left.Nodes) != len(right.Nodes) {
		t.Error("length doesn't match", left.Name, len(left.Nodes), len(right.Nodes))
		return
	}

	checkNodes(t, left.Nodes, right.Nodes)
}

func checkNodes(t *testing.T, left, right []*Node) {
	for len(left) > 0 {
		checkNode(t, left[0], right[0])
		if t.Failed() {
			return
		}

		left, right = left[1:], right[1:]
	}
}

func testSyntax(t *testing.T, st []syntaxTest) {
	traceLevel := TraceOff
	// traceLevel = TraceDebug

	for _, ti := range st {
		t.Run(ti.msg, func(t *testing.T) {
			s := NewSyntax(Options{Trace: NewTrace(traceLevel)})
			define(s, ti.syntax)

			if err := s.Init(); err != nil {
				t.Error(err)
				return
			}

			n, err := s.Parse(bytes.NewBufferString(ti.text))

			if ti.fail && err == nil {
				t.Error("failed to fail", n)
				return
			} else if !ti.fail && err != nil {
				t.Error(err)
				return
			} else if ti.fail {
				return
			}

			checkNode(t, n, ti.node)
		})
	}
}
