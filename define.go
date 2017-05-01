package mml

func (s *syntax) defineSyntax(primitive [][]interface{}, complex [][]string) error {
	for _, p := range primitive {
		if err := s.primitive(p[0].(string), p[1].(tokenType)); err != nil {
			return err
		}
	}

	for _, c := range complex {
		var err error
		switch c[0] {
		case "optional":
			err = s.optional(c[1], c[2])
		case "repetition":
			err = s.repetition(c[1], c[2])
		case "sequence":
			err = s.sequence(c[1], c[2:]...)
		case "choice":
			err = s.choice(c[1], c[2:]...)
		default:
			panic("invalid parser type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}
