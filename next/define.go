package next

import (
	"errors"
	"strconv"
)

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
		if len(d) < 3 {
			return errInvalidDefinition
		}

		switch d[0] {
		case "chars", "class":
			if len(d) < 4 {
				return errInvalidDefinition
			}
		case "quantifier":
			if len(d) != 6 {
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
			ct := stringToCommitType(d[2])
			err = s.Terminal(d[1], ct, Terminal{Anything: true})
		case "chars":
			ct := stringToCommitType(d[2])
			err = s.Terminal(d[1], ct, Terminal{Chars: d[3]})
		case "class":
			ct := stringToCommitType(d[2])
			err = s.Terminal(d[1], ct, Terminal{Class: d[3]})
		case "quantifier":
			ct := stringToCommitType(d[2])

			var min, max int
			min, err = strconv.Atoi(d[4])
			if err == nil {
				max, err = strconv.Atoi(d[5])
				if err == nil {
					err = s.Quantifier(d[1], ct, d[3], min, max)
				}
			}
		case "sequence":
			ct := stringToCommitType(d[2])
			err = s.Sequence(d[1], ct, d[3:]...)
		case "choice":
			ct := stringToCommitType(d[2])
			err = s.Choice(d[1], ct, d[3:]...)
		default:
			err = errInvalidDefinition
		}

		if err != nil {
			return err
		}
	}

	return nil
}
