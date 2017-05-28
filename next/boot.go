package next

import (
	"fmt"
	"strconv"
)

var definitions = [][]string{{
	"chars", "space", "alias", " ",
}, {
	"chars", "tab", "alias", "\\t",
}, {
	"chars", "nl", "alias", "\\n",
}, {
	"chars", "backspace", "alias", "\\b",
}, {
	"chars", "formfeed", "alias", "\\f",
}, {
	"chars", "carryreturn", "alias", "\\r",
}, {
	"chars", "verticaltab", "alias", "\\v",
}, {
	"choice",
	"ws",
	"alias",
	"space",
	"tab",
	"nl",
	"backspace",
	"formfeed",
	"carryreturn",
	"verticaltab",
}, {
	"chars", "open-block-comment", "alias", "/*",
}, {
	"chars", "close-block-comment", "alias", "*/",
}, {
	"chars", "star", "alias", "*",
}, {
	"class", "not-slash", "alias", "^/",
}, {
	"class", "not-star", "alias", "^*",
}, {
	"chars", "double-slash", "alias", "//",
}, {
	"class", "not-nl", "alias", "^\\n",
}, {
	"sequence", "not-block-close", "alias", "star", "not-slash",
}, {
	"choice", "block-comment-char", "alias", "not-block-close", "not-star",
}, {
	"quantifier", "block-comment-body", "alias", "block-comment-char", "0", "-1",
}, {
	"sequence",
	"block-comment",
	"alias",
	"open-block-comment",
	"block-comment-body",
	"close-block-comment",
}, {
	"quantifier", "not-nls", "alias", "not-nl", "0", "-1",
}, {
	"sequence", "line-comment", "alias", "double-slash", "not-nls",
}, {
	"choice", "comment-segment", "alias", "block-comment", "line-comment",
}, {
	"quantifier", "wss", "alias", "ws", "0", "-1",
}, {
	"quantifier", "optional-nl", "alias", "nl", "0", "1",
}, {
	"choice",
	"ws-no-nl",
	"alias",
	"space",
	"tab",
	"backspace",
	"formfeed",
	"carryreturn",
	"verticaltab",
}, {
	"sequence",
	"continue-comment-segment",
	"alias",
	"ws-no-nl",
	"optional-nl",
	"ws-no-nl",
	"comment-segment",
}, {
	"quantifier", "continue-comment", "alias", "continue-comment-segment", "0", "-1",
}, {
	"sequence",
	"comment",
	"none",
	"comment-segment",
	"continue-comment",
}, {
	"choice", "wsc", "alias", "ws", "comment",
}, {
	"quantifier", "wscs", "alias", "wsc", "0", "-1",
}, {
	"anything", "anything", "alias",
}, {
	"chars", "any-char", "none", ".",
}, {
	"chars", "open-square", "alias", "[",
}, {
	"chars", "close-square", "alias", "]",
}, {
	"chars", "class-not", "none", "^",
}, {
	"chars", "dash", "alias", "-",
}, {
	"quantifier", "optional-class-not", "alias", "class-not", "0", "1",
}, {
	"class", "not-class-control", "alias", "^\\\\\\[\\]\\^\\-",
}, {
	"chars", "escape", "alias", "\\\\",
}, {
	"sequence", "escaped-char", "alias", "escape", "anything",
}, {
	"choice", "class-char", "none", "not-class-control", "escaped-char",
}, {
	"sequence", "char-range", "none", "class-char", "dash", "class-char",
}, {
	"choice", "char-or-range", "alias", "class-char", "char-range",
}, {
	"quantifier", "chars-or-ranges", "alias", "char-or-range", "0", "-1",
}, {
	"sequence", "char-class", "none", "open-square", "optional-class-not", "chars-or-ranges", "close-square",
}, {
	"chars", "double-quote", "alias", "\\\"",
}, {
	"class", "not-char-sequence-control", "alias", "^\\\\\\\"",
}, {
	"choice", "sequence-char", "none", "not-char-sequence-control", "escaped-char",
}, {
	"quantifier", "char-sequence-chars", "alias", "sequence-char", "0", "-1",
}, {
	"sequence", "char-sequence", "none", "double-quote", "char-sequence-chars", "double-quote",
}, {
	"choice", "terminal-element", "alias", "any-char", "char-class", "char-sequence",
}, {
	"quantifier", "terminal", "none", "terminal-element", "1", "-1",
}, {
	"class", "symbol-char", "alias", "^\\\\ \\n\\t\\b\\f\\r\\v\\b/.\\[\\]\\\"{}\\^+*?|():=;",
}, {
	"quantifier", "symbol-chars", "alias", "symbol-char", "1", "-1",
}, {
	"sequence", "symbol", "none", "symbol-chars",
}, {
	"chars", "open-paren", "alias", "(",
}, {
	"chars", "close-paren", "alias", ")",
}, {
	"sequence", "group", "alias", "open-paren", "wscs", "expression", "wscs", "close-paren",
}, {
	"chars", "open-brace", "alias", "{",
}, {
	"chars", "close-brace", "alias", "}",
}, {
	"class", "digit", "alias", "0-9",
}, {
	"quantifier", "count", "none", "digit", "1", "-1",
}, {
	"sequence", "count-quantifier", "none", "open-brace", "wscs", "count", "wscs", "close-brace",
}, {
	"chars", "comma", "alias", ",",
}, {
	"sequence",
	"range-quantifier",
	"none",
	"open-brace",
	"wscs",
	"count",
	"wscs",
	"comma",
	"wscs",
	"count",
	"close-brace",
}, {
	"chars", "one-or-more", "none", "+",
}, {
	"chars", "zero-or-more", "none", "*",
}, {
	"chars", "zero-or-one", "none", "?",
}, {
	"choice",
	"quantity",
	"alias",
	"count-quantifier",
	"range-quantifier",
	"one-or-more",
	"zero-or-more",
	"zero-or-one",
}, {
	"choice", "quantifiable", "alias", "terminal", "symbol", "group",
}, {
	"sequence", "quantifier", "none", "quantifiable", "wscs", "quantity",
}, {
	"choice", "item", "alias", "terminal", "symbol", "group", "quantifier",
}, {
	"sequence", "item-continue", "alias", "wscs", "item",
}, {
	"quantifier", "items-continue", "alias", "item-continue", "0", "-1",
}, {
	"sequence", "sequence", "none", "item", "items-continue",
}, {
	"choice", "element", "alias", "terminal", "symbol", "group", "quantifier", "sequence",
}, {
	"chars", "pipe", "alias", "|",
}, {
	"sequence", "element-continue", "alias", "wscs", "pipe", "wscs", "element",
}, {
	"quantifier", "elements-continue", "alias", "element-continue", "1", "-1",
}, {
	"sequence", "choice", "none", "element", "elements-continue",
}, {
	"choice",
	"expression",
	"alias",
	"terminal",
	"symbol",
	"group",
	"quantifier",
	"sequence",
	"choice",
}, {
	"chars", "alias", "none", "alias",
}, {
	"chars", "root", "none", "root",
}, {
	"choice", "flag-word", "alias", "alias", "root",
}, {
	"chars", "colon", "alias", ":",
}, {
	"sequence", "flag", "alias", "colon", "flag-word",
}, {
	"quantifier", "flags", "alias", "flag", "0", "-1",
}, {
	"chars", "equal", "alias", "=",
}, {
	"sequence", "definition", "none", "symbol", "flags", "wscs", "equal", "wscs", "expression",
}, {
	"chars", "semicolon", "alias", ";",
}, {
	"choice", "wsc-or-semicolon", "alias", "wsc", "semicolon",
}, {
	"quantifier", "wsc-or-semicolons", "alias", "wsc-or-semicolon", "0", "-1",
}, {
	"sequence",
	"subsequent-definition",
	"alias",
	"wscs",
	"semicolon",
	"wsc-or-semicolons",
	"definition",
}, {
	"quantifier",
	"subsequent-definitions",
	"alias",
	"subsequent-definition",
	"0",
	"-1",
}, {
	"sequence",
	"definitions",
	"alias",
	"definition",
	"subsequent-definitions",
}, {
	"quantifier",
	"opt-definitions",
	"alias",
	"definitions",
	"0",
	"1",
}, {
	"sequence",
	"document",
	"root",
	"wsc-or-semicolons",
	"opt-definitions",
	"wsc-or-semicolons",
}}

func defineSyntax() (*Syntax, error) {
	l := TraceOff
	// l = TraceDebug
	s := NewSyntax(Options{Trace: NewTrace(l)})
	if err := define(s, definitions); err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	return s, nil
}

func flagsToCommitType(n []*Node) CommitType {
	var ct CommitType
	for _, ni := range n {
		switch ni.Name {
		case "alias":
			ct |= Alias
		case "root":
			ct |= Root
		}
	}

	return ct
}

func childName(name string, childIndex int) string {
	return fmt.Sprintf("%s:%d", name, childIndex)
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

func toRune(c string) rune {
	return []rune(c)[0]
}

func singleChar(n *Node) rune {
	s := n.Text()
	if s[0] == '\\' {
		return unescapeChar(toRune(s[1:]))
	}

	return toRune(s)
}

func defineClass(s *Syntax, name string, n []*Node) error {
	var (
		not    bool
		chars  []rune
		ranges [][]rune
	)

	if n[0].len() > 0 {
		not = true
	}

	for _, c := range n[1:] {
		switch c.Name {
		case "class-char":
			chars = append(chars, singleChar(c))
		case "class-range":
			ranges = append(ranges, []rune{singleChar(c.Nodes[0]), singleChar(c.Nodes[1])})
		}
	}

	return s.Class(name, not, chars, ranges)
}

func defineCharSequence(s *Syntax, name string, chars []*Node) ([]string, error) {
	var refs []string
	for i, c := range chars {
		char := singleChar(c)
		ref := childName(name, i)
		if err := s.Char(ref, char); err != nil {
			return nil, err
		}

		refs = append(refs, ref)
	}

	return refs, nil
}

func defineTerminalChars(s *Syntax, name string, c *Node) (refs []string, err error) {
	switch c.Name {
	case "any-char":
		err = s.AnyChar(name)
		refs = []string{name}
	case "char-class":
		err = defineClass(s, name, c.Nodes)
		refs = []string{name}
	case "char-sequence":
		refs, err = defineCharSequence(s, name, c.Nodes)
	}

	return
}

func defineTerminal(s *Syntax, name string, ct CommitType, t *Node) error {
	var refs []string
	for i, c := range t.Nodes {
		ref := childName(name, i)
		if crefs, err := defineTerminalChars(s, ref, c); err != nil {
			return err
		} else {
			refs = append(refs, crefs...)
		}
	}

	return s.Sequence(name, ct, refs...)
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
	case "terminal":
		err = defineTerminal(s, name, ct, expression)
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

func documentDefinition(s *Syntax, n *Node) error {
	return defineExpression(
		s,
		n.Nodes[0].Text(),
		flagsToCommitType(n.Nodes[1:len(n.Nodes)-1]),
		n.Nodes[len(n.Nodes)-1],
	)
}

func defineDocument(n *Node) (*Syntax, error) {
	if n.Name != "document" {
		return nil, ErrInvalidSyntax
	}

	s := NewSyntax(Options{Trace: NewTrace(TraceOff)})
	for _, ni := range n.Nodes {
		switch ni.Name {
		case "comment":
			continue
		case "definition":
			if err := documentDefinition(s, ni); err != nil {
				return nil, err
			}
		}
	}

	return s, nil
}
