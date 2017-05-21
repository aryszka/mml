package next

import (
	"bytes"
	"testing"
)

func TestComment(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg: "no comment",
	}, {
		msg:  "simple comment",
		text: "// some comment",
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			s, err := defineSyntax()
			if err != nil {
				t.Error(err)
				return
			}

			_, err = s.Parse(bytes.NewBufferString(ti.text))
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
}
