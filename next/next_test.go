package next

import (
	"bytes"
	"testing"
	"time"
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

func stringToCommitType(s string) CommitType {
	switch s {
	case "alias":
		return Alias
	default:
		return None
	}
}

func testSyntax(t *testing.T, st []syntaxTest) {
	// traceLevel := TraceDebug
	traceLevel := TraceOff

	for _, ti := range st {
		t.Run(ti.msg, func(t *testing.T) {
			s := NewSyntax(Options{Trace: NewTrace(traceLevel)})

			for _, d := range ti.syntax {
				if len(d) < 2 {
					t.Error("invalid syntax definition")
					return
				}

				switch d[0] {
				case "chars", "class":
					if len(d) < 3 {
						t.Error("invalid syntax definition")
						return
					}
				case "optional", "repetition":
					if len(d) != 4 {
						t.Error("invalid syntax definition")
						return
					}
				case "sequence", "choice":
					if len(d) < 4 {
						t.Error("invalid syntax definition")
						return
					}
				}

				var err error
				switch d[0] {
				case "anything":
					err = s.Terminal(d[1], Terminal{Anything: true})
				case "chars":
					ts := make([]Terminal, len(d)-2)
					for i, di := range d[2:] {
						ts[i] = Terminal{Chars: di}
					}

					err = s.Terminal(d[1], ts...)
				case "class":
					ts := make([]Terminal, len(d)-2)
					for i, di := range d[2:] {
						ts[i] = Terminal{Class: di}
					}

					err = s.Terminal(d[1], Terminal{Class: d[2]})
				case "optional":
					ct := stringToCommitType(d[2])
					err = s.Optional(d[1], ct, d[3])
				case "repetition":
					ct := stringToCommitType(d[2])
					err = s.Repetition(d[1], ct, d[3])
				case "sequence":
					ct := stringToCommitType(d[2])
					err = s.Sequence(d[1], ct, d[3:]...)
				case "choice":
					ct := stringToCommitType(d[2])
					err = s.Choice(d[1], ct, d[3:]...)
				}

				if err != nil {
					t.Error(err)
					return
				}
			}

			if err := s.Init(); err != nil {
				t.Error(err)
				return
			}

			start := time.Now()
			n, err := s.Parse(bytes.NewBufferString(ti.text))
			t.Log("parse time", time.Now().Sub(start))

			if ti.fail && err == nil {
				t.Error("failed to fail", n)
				return
			} else if !ti.fail && err != nil {
				t.Error(err)
				return
			} else if ti.fail {
				return
			}

			if !checkNode(n, ti.node) {
				t.Error("node doesn't match", n)
			}
		})
	}
}
