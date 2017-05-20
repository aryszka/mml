package next

import "fmt"

type Terminal struct {
	Chars    string
	Class    string
	Anything bool
}

func makeAnyCharDefinition(r *registry, name string, indexOffset int) []definition {
	name = fmt.Sprintf("%s:%d", name, indexOffset)
	return []definition{newAnyCharDefinition(r, name)}
}

func makeCharDefinitions(r *registry, name string, indexOffset int, chars string) []definition {
	c := []rune(chars)

	defs := make([]definition, len(c))
	for i, ci := range c {
		defs[i] = newCharDefinition(r, fmt.Sprintf("%s:%d", name, indexOffset+i), ci)
	}

	return defs
}

func makeClassDefinition(r *registry, name string, index int, class string) definition {
	name = fmt.Sprintf("%s:%d", name, index)
	c := []rune(class)

	var not bool
	if c[0] == '^' {
		not = true
		c = c[1:]
	}

	var (
		chars  []rune
		ranges [][]rune
	)

	for len(c) > 0 {
		if len(c) >= 3 && c[1] == '-' {
			ranges, c = append(ranges, []rune{c[0], c[2]}), c[3:]
			continue
		}

		chars, c = append(chars, c[0]), c[1:]
	}

	return newClassDefinition(r, name, not, chars, ranges)
}

func terminalDefinitions(r *registry, name string, t []Terminal) []definition {
	var defs []definition
	for i, ti := range t {
		if ti.Anything {
			defs = append(defs, makeAnyCharDefinition(r, name, i)...)
		} else if ti.Chars != "" {
			defs = append(defs, makeCharDefinitions(r, name, i, ti.Chars)...)
		} else {
			defs = append(defs, makeClassDefinition(r, name, i, ti.Class))
		}
	}

	return defs
}
