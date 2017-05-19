package next

import "testing"

func TestSymbol(t *testing.T) {
	testSyntax(t, []syntaxTest{{
		msg:    "word",
		syntax: [][]string{{"chars", "foo-word", "foo"}},
		text:   "foo",
		node: &Node{
			Name: "foo-word",
			From: 0,
			To:   3,
			Nodes: []*Node{{
				Name: "foo-word:0",
				From: 0,
				To:   1,
			}, {
				Name: "foo-word:1",
				From: 1,
				To:   2,
			}, {
				Name: "foo-word:2",
				From: 2,
				To:   3,
			}},
		},
	}, {
		msg:    "word, no match",
		syntax: [][]string{{"chars", "foo-word", "foo"}},
		text:   "bar",
		fail:   true,
	}, {
		msg:    "letter",
		syntax: [][]string{{"class", "a", "a-z"}},
		text:   "a",
		node: &Node{
			Name: "a",
			From: 0,
			To:   1,
			Nodes: []*Node{{
				Name: "a:0",
				From: 0,
				To:   1,
			}},
		},
	}, {
		msg:    "letter, fail",
		syntax: [][]string{{"class", "a", "a-z"}},
		text:   "A",
		fail:   true,
	}})
}
