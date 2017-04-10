package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/aryszka/mml"
)

func main() {
	f, err := os.Create("cpu.out")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	if err := mml.Compile(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
