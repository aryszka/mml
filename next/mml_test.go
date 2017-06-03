package next

import (
	"testing"
	"os"
)

func TestMML(t *testing.T) {
	s, err := defineSyntax()
	if err != nil {
		t.Fatal(err)
	}

	def, err := os.Open("syntax.p")
	if err != nil {
		t.Fatal(err)
	}

	defer def.Close()

	n, err := s.Parse(def)
	if err != nil {
		t.Fatal(err)
	}

	st, err := defineDocument(n)
	if err != nil {
		t.Fatal(err)
	}

	mml, err := os.Open("mml.p")
	if err != nil {
		t.Fatal(err)
	}

	defer mml.Close()

	n, err = st.Parse(mml)
	if err != nil {
		t.Fatal(err)
	}

	mmlst, err := defineDocument(n)
	if err != nil {
		t.Fatal(err)
	}

	tokens, err := os.Open("tokens.mml")
	if err != nil {
		t.Fatal(err)
	}

	defer tokens.Close()

	n, err = mmlst.Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(n)
	for _, ni := range n.Nodes {
		t.Log(ni)
	}
}
