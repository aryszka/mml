package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/aryszka/mml"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(b)

	for i := 0; i < 1; i++ {
		r.Seek(0, 0)
		if err := mml.Compile(r, os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
