package next

import "testing"

var stringSyntax = [][]string{{
	"chars", "quote", "\"",
}, {
	"class", "not-escaped", "^\\\"",
}, {
	"chars", "escape", "\\",
}, {
	"anything", "anything",
}, {
	"sequence", "escaped", "alias", "escape", "anything",
}, {
	"choice", "string-char", "alias", "escaped", "not-escaped",
}, {
	"repetition", "string-body", "alias", "string-char",
}, {
	"sequence", "string", "none", "quote", "string-body", "quote",
}}

func TestString(t *testing.T) {
	testSyntax(t, []syntaxTest{{
		msg:    "empty",
		text:   `""`,
		syntax: stringSyntax,
		node: &Node{
			Name: "string",
			from: 0,
			to:   2,
		},
	}, {
		msg:    "simple string",
		text:   `"foo bar baz"`,
		syntax: stringSyntax,
		node: &Node{
			Name: "string",
			from: 0,
			to:   13,
		},
	}, {
		msg:    "string with escape",
		text:   "\"fo\\o bar baz\"",
		syntax: stringSyntax,
		node: &Node{
			Name: "string",
			from: 0,
			to:   14,
		},
	}, {
		msg:    "string with quote",
		text:   "\"foo \\\"bar\\\" baz\"",
		syntax: stringSyntax,
		node: &Node{
			Name: "string",
			from: 0,
			to:   17,
		},
	}, {
		msg:    "string with backslash",
		text:   "\"foo \\\\bar baz\"",
		syntax: stringSyntax,
		node: &Node{
			Name: "string",
			from: 0,
			to:   15,
		},
	}})
}
