package next

import (
	"testing"
	"log"
	"os"
)

func TestMML(*testing.T) {
	s, err := defineSyntax()
	if err != nil {
		log.Fatalln(err)
	}

	def, err := os.Open("syntax.p")
	if err != nil {
		log.Fatalln(err)
	}

	defer def.Close()

	n, err := s.Parse(def)
	if err != nil {
		log.Fatalln(err)
	}

	st, err := defineDocument(n)
	if err != nil {
		log.Fatalln(err)
	}

	mml, err := os.Open("mml.p")
	if err != nil {
		log.Fatalln(err)
	}

	defer mml.Close()

	_, err = st.Parse(mml)
	if err != nil {
		log.Fatalln(err)
	}
}
