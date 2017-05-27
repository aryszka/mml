package next

import (
	"bytes"
	"testing"
)

func TestGroup(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "empty group",
		text: "a = ()",
		fail: true,
	}, {
		msg:  "group",
		text: "a = (a*)",
	}, {
		msg:  "group in group",
		text: "a = ((a*))",
		// TODO: test with sequence
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
