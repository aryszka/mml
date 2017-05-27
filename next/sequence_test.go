package next

import (
	"bytes"
	"testing"
)

func TestSequence(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "single",
		text: "a = a",
	}, {
		msg:  "multiple",
		text: "abc = a b c",
	}, {
		// TODO: try with choice here
		msg:  "combined",
		text: "a = a (b c)+ (d e? f{1, 3})",
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
