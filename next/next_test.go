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

func checkNode(left, right *Node) bool {
	if (left == nil) != (right == nil) {
		return false
	}

	if left == nil {
		return true
	}

	if left.Name != right.Name {
		return false
	}

	if left.from != right.from {
		return false
	}

	if left.to != right.to {
		return false
	}

	return checkNodes(left.Nodes, right.Nodes)
}

func checkNodes(left, right []*Node) bool {
	if len(left) != len(right) {
		return false
	}

	for len(left) > 0 {
		if !checkNode(left[0], right[0]) {
			return false
		}

		left, right = left[1:], right[1:]
	}

	return true
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

			// if !checkNode(n, ti.node) {
			// 	t.Error("node doesn't match", n)
			// }
		})
	}
}
