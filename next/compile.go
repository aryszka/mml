package next

func runesContain(rs []rune, r rune) bool {
	for _, ri := range rs {
		if ri == r {
			return true
		}
	}

	return false
}

func unescapeChar(c rune) rune {
	switch c {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'r':
		return '\r'
	case 'v':
		return '\v'
	default:
		return c
	}
}

func unescape(escape rune, banned []rune, chars []rune) ([]rune, error) {
	var (
		unescaped []rune
		escaped   bool
	)

	for _, ci := range chars {
		if escaped {
			unescaped = append(unescaped, unescapeChar(ci))
			escaped = false
			continue
		}

		switch {
		case ci == escape:
			escaped = true
		case runesContain(banned, ci):
			return nil, ErrInvalidCharacter
		default:
			unescaped = append(unescaped, ci)
		}
	}

	if escaped {
		return nil, ErrInvalidCharacter
	}

	return unescaped, nil
}

func compile(s *Syntax, n *Node) error {
	return nil
}
