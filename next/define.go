package next

import "errors"

var errInvalidDefinition = errors.New("invalid syntax definition")

func stringToCommitType(s string) CommitType {
	switch s {
	case "alias":
		return Alias
	default:
		return None
	}
}

func define(s *Syntax, defs [][]string) error {
	for _, d := range defs {
		if len(d) < 2 {
			return errInvalidDefinition
		}

		switch d[0] {
		case "chars", "class":
			if len(d) < 3 {
				return errInvalidDefinition
			}
		case "optional", "repetition":
			if len(d) != 4 {
				return errInvalidDefinition
			}
		case "sequence", "choice":
			if len(d) < 4 {
				return errInvalidDefinition
			}
		}

		var err error
		switch d[0] {
		case "anything":
			err = s.Terminal(d[1], Terminal{Anything: true})
		case "chars":
			ts := make([]Terminal, len(d)-2)
			for i, di := range d[2:] {
				ts[i] = Terminal{Chars: di}
			}

			err = s.Terminal(d[1], ts...)
		case "class":
			ts := make([]Terminal, len(d)-2)
			for i, di := range d[2:] {
				ts[i] = Terminal{Class: di}
			}

			err = s.Terminal(d[1], Terminal{Class: d[2]})
		case "optional":
			ct := stringToCommitType(d[2])
			err = s.Optional(d[1], ct, d[3])
		case "repetition":
			ct := stringToCommitType(d[2])
			err = s.Repetition(d[1], ct, d[3])
		case "sequence":
			ct := stringToCommitType(d[2])
			err = s.Sequence(d[1], ct, d[3:]...)
		case "choice":
			ct := stringToCommitType(d[2])
			err = s.Choice(d[1], ct, d[3:]...)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
