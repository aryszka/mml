package next

import (
	"bytes"
	"testing"
)

func TestQuantifier(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "count",
		text: "a = \"a\"{3}",
	}, {
		msg:  "range",
		text: "a = \"a\"{3, 9}",
	}, {
		msg:  "one or more",
		text: "a = \"a\"+",
	}, {
		msg:  "zero or more",
		text: "a = \"a\"*",
	}, {
		msg:  "zero or one",
		text: "a = \"a\"?",
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
