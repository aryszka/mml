package next

import (
	"bytes"
	"testing"
)

func TestSymbol(t *testing.T) {
	testSyntax(t, []syntaxTest{{
		msg: "word ignored",
		syntax: [][]string{
			{"chars", "foo-word-chars", "foo"},
		},
		text: "foo",
	}, {
		msg: "word",
		syntax: [][]string{
			{"chars", "foo-word-chars", "foo"},
			{"sequence", "foo-word", "none", "foo-word-chars"},
		},
		text: "foo",
		node: &Node{
			Name: "foo-word",
			from: 0,
			to:   3,
		},
	}, {
		msg:    "word, no match",
		syntax: [][]string{{"chars", "foo-word", "foo"}},
		text:   "bar",
		fail:   true,
	}, {
		msg:    "word, no match, last",
		syntax: [][]string{{"chars", "bar-word", "bar"}},
		text:   "baz",
		fail:   true,
	}, {
		msg:    "char class, ignored",
		syntax: [][]string{{"class", "a", "a-z"}},
		text:   "a",
	}, {
		msg: "char class",
		syntax: [][]string{
			{"class", "lowercase-chars", "a-z"},
			{"sequence", "lowercase", "none", "lowercase-chars"},
		},
		text: "a",
		node: &Node{
			Name: "lowercase",
			from: 0,
			to:   1,
		},
	}, {
		msg:    "char class, fail",
		syntax: [][]string{{"class", "a", "a-z"}},
		text:   "A",
		fail:   true,
	}, {
		msg: "symbol",
		syntax: [][]string{
			{"class", "letter", "a-z"},
			{"class", "symbol-char", "a-zA-Z0-9_"},
			{"repetition", "symbol-chars", "alias", "symbol-char"},
			{"sequence", "symbol", "none", "letter", "symbol-chars"},
		},
		text: "fooBar",
		node: &Node{
			Name: "symbol",
			from: 0,
			to:   6,
		},
	}})
}

func TestSymbolSyntax(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "simple",
		text: "foo = bar",
	}, {
		msg:  "with flags",
		text: "foo:alias = bar",
	}, {
		msg:  "with invalid flag",
		text: "foo:bar = baz",
		fail: true,
	}, {
		msg:  "with invalid char",
		text: "f[o]o = bar",
		fail: true,
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			s, err := defineSyntax()
			if err != nil {
				t.Error(err)
				return
			}

			_, err = s.Parse(bytes.NewBufferString(ti.text))
			if ti.fail && err == nil {
				t.Error("failed to fail")
			} else if !ti.fail && err != nil {
				t.Error(err)
			}
		})
	}
}
