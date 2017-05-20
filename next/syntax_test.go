package next

import "testing"

func TestOptional(t *testing.T) {
	testSyntax(t, []syntaxTest{{
		msg: "fake rsc check, executes in time due to hungriness",
		syntax: [][]string{{
			"chars", "a", "a",
		}, {
			"anything", "any-char",
		}, {
			"optional", "opt-a", "alias", "a",
		}, {
			"repetition", "anything", "alias", "any-char",
		}, {
			"sequence", "match", "none",
			"opt-a", "opt-a", "opt-a", "a", "a", "a",
			"anything",
		}},
		text: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		node: &Node{
			Name: "match",
			from: 0,
			to:   36,
		},
	}})

	// TODO: could test the complexity of f()()()()()(), too
}
