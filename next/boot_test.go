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

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error(err)
		return
	}

	n0, err := s.Parse(f)
	if err != nil {
		t.Error(err)
		return
	}

	s0 := NewSyntax(trace)
	if err := compile(s0, n0); err != nil {
		t.Error(err)
		return
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error(err)
		return
	}

	n1, err := s.Parse(f)
	if err != nil {
		t.Error(err)
		return
	}

	checkNode(t, n1, n0)
}
