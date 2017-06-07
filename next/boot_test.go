package next

import (
	"os"
	"testing"
)

func TestBoot(t *testing.T) {
	var trace Trace
	// trace = NewTrace(1)

	b, err := initBoot(trace, bootDefinitions)
	if err != nil {
		t.Error(err)
		return
	}

	f, err := os.Open("syntax.p")
	if err != nil {
		t.Error(err)
		return
	}

	s := NewSyntax(trace)
	if err = s.read(b, f); err != nil {
		t.Error(err)
		return
	}
}
