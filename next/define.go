package next

import "strconv"

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

func dropComments(n *Node) *Node {
	ncc := *n
	nc := &ncc

	nc.Nodes = nil
	for _, ni := range n.Nodes {
		if ni.Name == "comment" {
			continue
		}

		nc.Nodes = append(nc.Nodes, dropComments(ni))
	}

	return nc
}

func flagsToCommitType(n []*Node) CommitType {
	var ct CommitType
	for _, ni := range n {
		switch ni.Name {
		case "alias":
			ct |= Alias
		case "doc":
			ct |= Documentation
		case "root":
			ct |= Root
		}
	}

	return ct
}

func toRune(c string) rune {
	return []rune(c)[0]
}

func nodeChar(n *Node) rune {
	s := n.Text()
	if s[0] == '\\' {
		return unescapeChar(toRune(s[1:]))
	}

	return toRune(s)
}

func defineMembers(s *Syntax, name string, n ...*Node) ([]string, error) {
	var refs []string
	for i, ni := range n {
		nmi := childName(name, i)
		switch ni.Name {
		case "symbol":
			refs = append(refs, ni.Text())
		default:
			refs = append(refs, nmi)
			if err := defineExpression(s, nmi, Alias, ni); err != nil {
				return nil, err
			}
		}
	}

	return refs, nil
}

func defineClass(s *Syntax, name string, ct CommitType, n []*Node) error {
	var (
		not    bool
		chars  []rune
		ranges [][]rune
	)

	if n[0].Name == "class-not" {
		not, n = true, n[1:]
	}

	for _, c := range n {
		switch c.Name {
		case "class-char":
			chars = append(chars, nodeChar(c))
		case "char-range":
			ranges = append(ranges, []rune{nodeChar(c.Nodes[0]), nodeChar(c.Nodes[1])})
		}
	}

	return s.Class(name, ct, not, chars, ranges)
}

func defineCharSequence(s *Syntax, name string, ct CommitType, charNodes []*Node) error {
	var chars []rune
	for _, ci := range charNodes {
		chars = append(chars, nodeChar(ci))
	}

	return s.CharSequence(name, ct, chars)
}

func defineQuantifier(s *Syntax, name string, ct CommitType, n *Node, q *Node) error {
	refs, err := defineMembers(s, name, n)
	if err != nil {
		return err
	}

	var min, max int
	switch q.Name {
	case "count-quantifier":
		min, err = strconv.Atoi(q.Nodes[0].Text())
		if err != nil {
			return err
		}

		max = min
	case "range-quantifier":
		min, err = strconv.Atoi(q.Nodes[0].Text())
		if err != nil {
			return err
		}

		max, err = strconv.Atoi(q.Nodes[1].Text())
		if err != nil {
			return err
		}
	case "one-or-more":
		min, max = 1, -1
	case "zero-or-more":
		min, max = 0, -1
	case "zero-or-one":
		min, max = 0, 1
	}

	return s.Quantifier(name, ct, refs[0], min, max)
}

func defineSequence(s *Syntax, name string, ct CommitType, n ...*Node) error {
	refs, err := defineMembers(s, name, n...)
	if err != nil {
		return err
	}

	return s.Sequence(name, ct, refs...)
}

func defineChoice(s *Syntax, name string, ct CommitType, n ...*Node) error {
	refs, err := defineMembers(s, name, n...)
	if err != nil {
		return err
	}

	return s.Choice(name, ct, refs...)
}

func defineExpression(s *Syntax, name string, ct CommitType, expression *Node) error {
	var err error
	switch expression.Name {
	case "any-char":
		err = s.AnyChar(name, ct)
	case "char-class":
		err = defineClass(s, name, ct, expression.Nodes)
	case "char-sequence":
		err = defineCharSequence(s, name, ct, expression.Nodes)
	case "symbol":
		err = defineSequence(s, name, ct, expression)
	case "quantifier":
		err = defineQuantifier(s, name, ct, expression.Nodes[0], expression.Nodes[1])
	case "sequence":
		err = defineSequence(s, name, ct, expression.Nodes...)
	case "choice":
		err = defineChoice(s, name, ct, expression.Nodes...)
	}

	return err
}

func defineDefinition(s *Syntax, n *Node) error {
	return defineExpression(
		s,
		n.Nodes[0].Text(),
		flagsToCommitType(n.Nodes[1:len(n.Nodes)-1]),
		n.Nodes[len(n.Nodes)-1],
	)
}

func define(s *Syntax, n *Node) error {
	if n.Name != "syntax" {
		return ErrInvalidSyntax
	}

	n = dropComments(n)

	for _, ni := range n.Nodes {
		if err := defineDefinition(s, ni); err != nil {
			return err
		}
	}

	return nil
}
