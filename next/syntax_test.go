package next

import "testing"

func TestBoot(t *testing.T) {
	t.Skip()
	boot()
}

func TestOptional(t *testing.T) {
	testSyntax(t, []syntaxTest{{
		msg: "fake rsc check, executes in time due to hungriness",
		syntax: [][]string{{
			"chars", "a", "alias", "a",
		}, {
			"anything", "any-char", "alias",
		}, {
			"quantifier", "opt-a", "alias", "a", "0", "1",
		}, {
			"quantifier", "anything", "alias", "any-char", "0", "-1",
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
