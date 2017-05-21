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
