package next

import (
	"bytes"
	"testing"
)

func TestTerminal(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "any char",
		text: "a-char = .",
	}, {
		msg:  "class",
		text: "digit = [0-9]",
	}, {
		msg:  "class, with multiple ranges",
		text: "digit = [0-9a-z]",
	}, {
		msg:  "class, with multiple ranges and chars",
		text: "digit = [0-9a-z_]",
	}, {
		msg:  "class, with multiple ranges, chars and escaped chars",
		text: "digit = [0-9a-z_\\-\\\\\\]]",
	}, {
		msg:  "char sequence",
		text: "foo = \"foo\"",
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
