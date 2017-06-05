package next

import (
	"log"
	"os"
	"testing"
	"time"
)

func boot(t *testing.T) {
	s, err := defineSyntax()
	if err != nil {
		log.Fatalln(err)
	}

	def, err := os.Open("syntax.p")
	if err != nil {
		log.Fatalln(err)
	}

	defer def.Close()

	start := time.Now()
	n, err := s.Parse(def)
	log.Println(time.Since(start))
	if err != nil {
		log.Fatalln(err)
	}

	st, err := defineDocument(n)
	if err != nil {
		log.Fatalln("document:", err)
	}

	if _, err := def.Seek(0, 0); err != nil {
		log.Fatalln(err)
	}

	start = time.Now()
	nt, err := st.Parse(def)
	log.Println(time.Since(start))
	if err != nil {
		log.Fatalln(err)
	}

	checkNode(t, nt, n)
	// for {
	// 	if len(nt.Nodes) > 0 {
	// 		t.Log("<", nt.Nodes[0].Name, nt.Nodes[0])
	// 		nt.Nodes = nt.Nodes[1:]
	// 	}

	// 	if len(n.Nodes) > 0 {
	// 		t.Log(">", n.Nodes[0].Name, n.Nodes[0])
	// 		n.Nodes = n.Nodes[1:]
	// 	}

	// 	if len(nt.Nodes) == 0 && len(n.Nodes) == 0 {
	// 		break
	// 	}
	// }

	stt, err := defineDocument(nt)
	if err != nil {
		log.Fatalln("document:", err)
	}

	if _, err := def.Seek(0, 0); err != nil {
		log.Fatalln(err)
	}

	start = time.Now()
	ntt, err := stt.Parse(def)
	log.Println(time.Since(start))
	if err != nil {
		log.Fatalln(err)
	}

	checkNode(t, ntt, nt)
}

func TestBoot(t *testing.T) {
	boot(t)
}

func TestOptional(t *testing.T) {
	testSyntax(t, []syntaxTest{{
		msg: "fake rsc check, executes in time due to hungriness",
		syntax: [][]string{{
			"chars", "a", "alias", "a",
		}, {
			"anything", "any-char", "alias",
		}, {
			"quantifier", "opt-a", "alias", "a", "0", "1",
		}, {
			"quantifier", "anything", "alias", "any-char", "0", "-1",
		}, {
			"sequence", "match", "none",
			"opt-a", "opt-a", "opt-a", "a", "a", "a",
			"anything",
		}},
		text: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		node: &Node{
			Name: "match",
			from: 0,
			to:   36,
		},
	}})

	// TODO: could test the complexity of f()()()()()(), too
}
