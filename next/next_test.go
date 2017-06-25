package next

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

type testItem struct {
	msg            string
	text           string
	fail           bool
	node           *Node
	nodes          []*Node
	ignorePosition bool
}

func testSyntaxReader(r io.Reader, traceLevel int) (*Syntax, error) {
	trace := NewTrace(0)

	b, err := bootSyntax(trace)
	if err != nil {
		return nil, err
	}

	doc, err := b.Parse(r)
	if err != nil {
		return nil, err
	}

	trace = NewTrace(traceLevel)
	s := NewSyntax(trace)
	if err := define(s, doc); err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	return s, nil
}

func testSyntaxString(s string, traceLevel int) (*Syntax, error) {
	return testSyntaxReader(bytes.NewBufferString(s), traceLevel)
}

func testSyntax(file string, traceLevel int) (*Syntax, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return testSyntaxReader(f, traceLevel)
}

func checkNodesPosition(t *testing.T, left, right []*Node, position bool) {
	if len(left) != len(right) {
		t.Error("length doesn't match", len(left), len(right))
		return
	}

	for len(left) > 0 {
		checkNodePosition(t, left[0], right[0], position)
		if t.Failed() {
			return
		}

		left, right = left[1:], right[1:]
	}
}

func checkNodePosition(t *testing.T, left, right *Node, position bool) {
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

	if position && left.from != right.from {
		t.Error("from doesn't match", left.Name, left.from, right.from)
		return
	}

	if position && left.to != right.to {
		t.Error("to doesn't match", left.Name, left.to, right.to)
		return
	}

	if len(left.Nodes) != len(right.Nodes) {
		t.Error("length doesn't match", left.Name, len(left.Nodes), len(right.Nodes))
		t.Log(left)
		t.Log(right)
		for {
			if len(left.Nodes) > 0 {
				t.Log("<", left.Nodes[0])
				left.Nodes = left.Nodes[1:]
			}

			if len(right.Nodes) > 0 {
				t.Log(">", right.Nodes[0])
				right.Nodes = right.Nodes[1:]
			}

			if len(left.Nodes) == 0 && len(right.Nodes) == 0 {
				break
			}
		}
		return
	}

	checkNodesPosition(t, left.Nodes, right.Nodes, position)
}

func checkNodes(t *testing.T, left, right []*Node) {
	checkNodesPosition(t, left, right, true)
}

func checkNode(t *testing.T, left, right *Node) {
	checkNodePosition(t, left, right, true)
}

func checkNodesIgnorePosition(t *testing.T, left, right []*Node) {
	checkNodesPosition(t, left, right, false)
}

func checkNodeIgnorePosition(t *testing.T, left, right *Node) {
	checkNodePosition(t, left, right, false)
}

func testReaderTrace(t *testing.T, r io.Reader, rootName string, traceLevel int, tests []testItem) {
	s, err := testSyntaxReader(r, traceLevel)
	if err != nil {
		t.Error(err)
		return
	}

	start := time.Now()
	defer func() { t.Log("\ntotal duration", time.Since(start)) }()

	for _, ti := range tests {
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

			t.Log(n)

			cn := checkNode
			if ti.ignorePosition {
				cn = checkNodeIgnorePosition
			}

			if ti.node != nil {
				cn(t, n, ti.node)
			} else {
				cn(t, n, &Node{
					Name:  rootName,
					from:  0,
					to:    len(ti.text),
					Nodes: ti.nodes,
				})
			}
		})
	}
}

func testStringTrace(t *testing.T, s string, traceLevel int, tests []testItem) {
	testReaderTrace(t, bytes.NewBufferString(s), "", traceLevel, tests)
}

func testString(t *testing.T, s string, tests []testItem) {
	testStringTrace(t, s, 0, tests)
}

func testTrace(t *testing.T, file, rootName string, traceLevel int, tests []testItem) {
	f, err := os.Open(file)
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()
	testReaderTrace(t, f, rootName, traceLevel, tests)
}

func test(t *testing.T, file, rootName string, tests []testItem) {
	testTrace(t, file, rootName, 0, tests)
}

func TestRecursion(t *testing.T) {
	testString(
		t,
		`A = "a" | A "a"`,
		[]testItem{{
			msg:  "recursion in choice, right, left, commit",
			text: "aaa",
			node: &Node{
				Name: "A",
				Nodes: []*Node{{
					Name: "A",
					Nodes: []*Node{{
						Name: "A",
					}},
				}},
			},
			ignorePosition: true,
		}},
	)

	testString(
		t,
		`A = "a" | "a" A`,
		[]testItem{{
			msg:  "recursion in choice, right, right, commit",
			text: "aaa",
			node: &Node{
				Name: "A",
				Nodes: []*Node{{
					Name: "A",
					Nodes: []*Node{{
						Name: "A",
					}},
				}},
			},
			ignorePosition: true,
		}},
	)

	testString(
		t,
		`A = "a" A | "a"`,
		[]testItem{{
			msg:  "recursion in choice, left, right, commit",
			text: "aaa",
			node: &Node{
				Name: "A",
				Nodes: []*Node{{
					Name: "A",
					Nodes: []*Node{{
						Name: "A",
					}},
				}},
			},
			ignorePosition: true,
		}},
	)

	testString(
		t,
		`A = A "a" | "a"`,
		[]testItem{{
			msg:  "recursion in choice, left, left, commit",
			text: "aaa",
			node: &Node{
				Name: "A",
				Nodes: []*Node{{
					Name: "A",
					Nodes: []*Node{{
						Name: "A",
					}},
				}},
			},
			ignorePosition: true,
		}},
	)

	testString(
		t,
		`A':alias = "a" | A' "a"; A = A'`,
		[]testItem{{
			msg:  "recursion in choice, right, left, alias",
			text: "aaa",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}},
	)

	testString(
		t,
		`A':alias = "a" | "a" A'; A = A'`,
		[]testItem{{
			msg:  "recursion in choice, right, right, alias",
			text: "aaa",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}},
	)

	testString(
		t,
		`A':alias = "a" A' | "a"; A = A'`,
		[]testItem{{
			msg:  "recursion in choice, left, right, alias",
			text: "aaa",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}},
	)

	testString(
		t,
		`A':alias = A' "a" | "a"; A = A'`,
		[]testItem{{
			msg:  "recursion in choice, left, left, alias",
			text: "aaa",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}},
	)
}

func TestSequence(t *testing.T) {
	testString(
		t,
		`AB = "a" | "a"? "a"? "b" "b"`,
		[]testItem{{
			msg:  "sequence with optional items",
			text: "abb",
			node: &Node{
				Name: "AB",
				to:   3,
			},
		}, {
			msg:  "sequence with optional items, none",
			text: "bb",
			node: &Node{
				Name: "AB",
				to:   2,
			},
		}},
	)

	testString(
		t,
		`A = "a" | (A?)*`,
		[]testItem{{
			msg:  "sequence in choice with redundant quantifier",
			text: "aaa",
			node: &Node{
				Name: "A",
				Nodes: []*Node{{
					Name: "A",
				}, {
					Name: "A",
				}, {
					Name: "A",
				}},
			},
			ignorePosition: true,
		}},
	)

	testString(
		t,
		`A = ("a"*)*`,
		[]testItem{{
			msg:  "sequence with redundant quantifier",
			text: "aaa",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}},
	)
}

func TestQuantifiers(t *testing.T) {
	testString(
		t,
		`A = "a" "b"{0} "a"`,
		[]testItem{{
			msg:  "zero",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero, fail",
			text: "aba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{1} "a"`,
		[]testItem{{
			msg:  "one, missing",
			text: "aa",
			fail: true,
		}, {
			msg:  "one",
			text: "aba",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}, {
			msg:  "one, too much",
			text: "abba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{3} "a"`,
		[]testItem{{
			msg:  "three, missing",
			text: "abba",
			fail: true,
		}, {
			msg:  "three",
			text: "abbba",
			node: &Node{
				Name: "A",
				to:   5,
			},
		}, {
			msg:  "three, too much",
			text: "abbbba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{0,1} "a"`,
		[]testItem{{
			msg:  "zero or one explicit, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or one explicit",
			text: "aba",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}, {
			msg:  "zero or one explicit, too much",
			text: "abba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{,1} "a"`,
		[]testItem{{
			msg:  "zero or one explicit, omit zero, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or one explicit, omit zero",
			text: "aba",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}, {
			msg:  "zero or one explicit, omit zero, too much",
			text: "abba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"? "a"`,
		[]testItem{{
			msg:  "zero or one explicit, shortcut, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or one explicit, shortcut",
			text: "aba",
			node: &Node{
				Name: "A",
				to:   3,
			},
		}, {
			msg:  "zero or one explicit, shortcut, too much",
			text: "abba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{0,3} "a"`,
		[]testItem{{
			msg:  "zero or three, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or three",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}, {
			msg:  "zero or three",
			text: "abbba",
			node: &Node{
				Name: "A",
				to:   5,
			},
		}, {
			msg:  "zero or three, too much",
			text: "abbbba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{,3} "a"`,
		[]testItem{{
			msg:  "zero or three, omit zero, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or three, omit zero",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}, {
			msg:  "zero or three, omit zero",
			text: "abbba",
			node: &Node{
				Name: "A",
				to:   5,
			},
		}, {
			msg:  "zero or three, omit zero, too much",
			text: "abbbba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{1,3} "a"`,
		[]testItem{{
			msg:  "one or three, missing",
			text: "aa",
			fail: true,
		}, {
			msg:  "one or three",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}, {
			msg:  "one or three",
			text: "abbba",
			node: &Node{
				Name: "A",
				to:   5,
			},
		}, {
			msg:  "one or three, too much",
			text: "abbbba",
			fail: true,
		}},
	)

	testString(
		t,
		`A = "a" "b"{3,5} "a"`,
		[]testItem{{
			msg:  "three or five, missing",
			text: "abba",
			fail: true,
		}, {
			msg:  "three or five",
			text: "abbbba",
			node: &Node{
				Name: "A",
				to:   6,
			},
		}, {
			msg:  "three or five",
			text: "abbbbba",
			node: &Node{
				Name: "A",
				to:   7,
			},
		}, {
			msg:  "three or five, too much",
			text: "abbbbbba",
			fail: true,
		}},
	)

	testStringTrace(
		t,
		`A = "a" "b"{0,} "a"`,
		1,
		[]testItem{{
			msg:  "zero or more, explicit, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or more, explicit",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}},
	)

	testStringTrace(
		t,
		`A = "a" "b"* "a"`,
		1,
		[]testItem{{
			msg:  "zero or more, shortcut, missing",
			text: "aa",
			node: &Node{
				Name: "A",
				to:   2,
			},
		}, {
			msg:  "zero or more, shortcut",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}},
	)

	testStringTrace(
		t,
		`A = "a" "b"{1,} "a"`,
		1,
		[]testItem{{
			msg:  "one or more, explicit, missing",
			text: "aa",
			fail: true,
		}, {
			msg:  "one or more, explicit",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}},
	)

	testStringTrace(
		t,
		`A = "a" "b"+ "a"`,
		1,
		[]testItem{{
			msg:  "one or more, shortcut, missing",
			text: "aa",
			fail: true,
		}, {
			msg:  "one or more, shortcut",
			text: "abba",
			node: &Node{
				Name: "A",
				to:   4,
			},
		}},
	)

	testStringTrace(
		t,
		`A = "a" "b"{3,} "a"`,
		1,
		[]testItem{{
			msg:  "three or more, explicit, missing",
			text: "abba",
			fail: true,
		}, {
			msg:  "three or more, explicit",
			text: "abbbba",
			node: &Node{
				Name: "A",
				to:   6,
			},
		}},
	)
}
