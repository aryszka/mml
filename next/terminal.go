package next

import (
	"errors"
	"fmt"
)

type Terminal struct {
	Anything bool
	Class    string
	Chars    string
}

var ErrInvalidTerminal = errors.New("invalid terminal")

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

func runesContain(rs []rune, r rune) bool {
	for _, ri := range rs {
		if ri == r {
			return true
		}
	}

	return false
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
			return nil, ErrInvalidTerminal
		default:
			unescaped = append(unescaped, ci)
		}
	}

	if escaped {
		return nil, ErrInvalidTerminal
	}

	return unescaped, nil
}

func makeAnyCharDefinition(r *registry, name string, indexOffset int) definition {
	name = fmt.Sprintf("%s:%d", name, indexOffset)
	return newChar(r, name, true, false, nil, nil)
}

func makeCharDefinitions(r *registry, name string, indexOffset int, chars string) ([]definition, error) {
	c, err := unescape('\\', []rune{'"', '\\'}, []rune(chars))
	if err != nil {
		return nil, err
	}

	defs := make([]definition, len(c))
	for i, ci := range c {
		defs[i] = newChar(
			r,
			fmt.Sprintf("%s:%d", name, indexOffset+i),
			false,
			false,
			[]rune{ci},
			nil,
		)
	}

	return defs, nil
}

func parseCharClass(c []rune) (not bool, chars []rune, ranges [][]rune, err error) {
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
			err = ErrInvalidTerminal
			return
		}

		if c0 == '\\' {
			if len(c) == 0 {
				err = ErrInvalidTerminal
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
				err = ErrInvalidTerminal
				return
			}

			c1, c = unescapeChar(c[0]), c[1:]
		}

		ranges = append(ranges, []rune{c0, c1})
	}
}

func makeClassDefinition(r *registry, name string, index int, class string) (definition, error) {
	not, chars, ranges, err := parseCharClass([]rune(class))
	if err != nil {
		return nil, err
	}

	name = fmt.Sprintf("%s:%d", name, index)
	return newChar(r, name, false, not, chars, ranges), nil
}

func terminalDefinitions(r *registry, name string, t []Terminal) ([]definition, error) {
	var defs []definition
	for i, ti := range t {
		if ti.Anything {
			defs = append(defs, makeAnyCharDefinition(r, name, i))
		} else if ti.Chars != "" {
			di, err := makeCharDefinitions(r, name, i, ti.Chars)
			if err != nil {
				return nil, err
			}

			defs = append(defs, di...)
		} else {
			di, err := makeClassDefinition(r, name, i, ti.Class)
			if err != nil {
				return nil, err
			}

			defs = append(defs, di)
		}
	}

	return defs, nil
}
