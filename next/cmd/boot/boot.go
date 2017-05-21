package main

import "log"

func main() {
	s := next.NewSyntax(next.Options{})
	if err := define(s, defintions); err != nil {
		log.Fatalln(err)
	}

	if err := s.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := s.Generate(os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
