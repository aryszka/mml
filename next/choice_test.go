package next

import (
	"bytes"
	"testing"
)

func TestChoice(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "choice",
		text: "a = a | b",
	}, {
		msg:  "multiple",
		text: "abcd = a | b c | d",
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
