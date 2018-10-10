package main

import (
	"log"
	"os"

	"github.com/aryszka/mml"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	if err := mml.EvalInput(f); err != nil {
		log.Fatalln(err)
	}
}
