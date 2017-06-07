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
	case "doc":
		return Documentation
	case "root":
		return Root
	default:
		return None
	}
}

func checkBootDefinitionLength(d []string) error {
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

	return nil
}

func parseClass(c []rune) (not bool, chars []rune, ranges [][]rune, err error) {
	if c[0] == '^' {
		not = true
		c = c[1:]
	}

	for {
		if len(c) == 0 {
			return
		}

		var c0 rune
		c0, c = c[0], c[1:]
		switch c0 {
		case '[', ']', '^', '-':
			err = errInvalidDefinition
			return
		}

		if c0 == '\\' {
			if len(c) == 0 {
				err = errInvalidDefinition
				return
			}

			c0, c = unescapeChar(c[0]), c[1:]
		}

		if len(c) < 2 || c[0] != '-' {
			chars = append(chars, c0)
			continue
		}

		var c1 rune
		c1, c = c[1], c[2:]
		if c1 == '\\' {
			if len(c) == 0 {
				err = errInvalidDefinition
				return
			}

			c1, c = unescapeChar(c[0]), c[1:]
		}

		ranges = append(ranges, []rune{c0, c1})
	}
}

func defineAnything(s *Syntax, d []string) error {
	ct := stringToCommitType(d[2])
	return s.AnyChar(d[1], ct)
}

func defineClass(s *Syntax, d []string) error {
	ct := stringToCommitType(d[2])

	not, chars, ranges, err := parseClass([]rune(d[3]))
	if err != nil {
		return err
	}

	return s.Class(d[1], ct, not, chars, ranges)
}

func defineCharSequence(s *Syntax, d []string) error {
	ct := stringToCommitType(d[2])

	chars, err := unescape('\\', []rune{'"', '\\'}, []rune(d[3]))
	if err != nil {
		return err
	}

	return s.CharSequence(d[1], ct, chars)
}

func defineQuantifier(s *Syntax, d []string) error {
	ct := stringToCommitType(d[2])

	var (
		min, max int
		err      error
	)

	if min, err = strconv.Atoi(d[4]); err != nil {
		return err
	}

	if max, err = strconv.Atoi(d[5]); err != nil {
		return err
	}

	return s.Quantifier(d[1], ct, d[3], min, max)
}

func defineSequence(s *Syntax, d []string) error {
	ct := stringToCommitType(d[2])
	return s.Sequence(d[1], ct, d[3:]...)
}

func defineChoice(s *Syntax, d []string) error {
	ct := stringToCommitType(d[2])
	return s.Choice(d[1], ct, d[3:]...)
}

func defineBoot(s *Syntax, d []string) error {
	switch d[0] {
	case "anything":
		return defineAnything(s, d)
	case "class":
		return defineClass(s, d)
	case "chars":
		return defineCharSequence(s, d)
	case "quantifier":
		return defineQuantifier(s, d)
	case "sequence":
		return defineSequence(s, d)
	case "choice":
		return defineChoice(s, d)
	default:
		return errInvalidDefinition
	}
}

func defineAllBoot(s *Syntax, defs [][]string) error {
	for _, d := range defs {
		if err := defineBoot(s, d); err != nil {
			return err
		}
	}

	return nil
}

func initBoot(t Trace, definitions [][]string) (*Syntax, error) {
	s := NewSyntax(t)
	if err := defineAllBoot(s, definitions); err != nil {
		return nil, err
	}

	return s, s.Init()
}
