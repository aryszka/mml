package next

import (
	"log"
	"os"
	"time"
)

var definitions = [][]string{{
	"chars", "space", " ",
}, {
	"chars", "tab", "\\t",
}, {
	"chars", "nl", "\\n",
}, {
	"chars", "backspace", "\\b",
}, {
	"chars", "formfeed", "\\f",
}, {
	"chars", "carryreturn", "\\r",
}, {
	"chars", "verticaltab", "\\v",
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
	"chars", "open-block-comment", "/*",
}, {
	"chars", "close-block-comment", "*/",
}, {
	"chars", "star", "*",
}, {
	"class", "not-slash", "^/",
}, {
	"class", "not-star", "^*",
}, {
	"chars", "double-slash", "//",
}, {
	"class", "not-nl", "^\\n",
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
	"sequence",
	"continue-comment-segment",
	"alias",
	"ws",
	"optional-nl",
	"ws",
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
	"anything", "anything",
}, {
	"chars", "any-char", ".",
}, {
	"chars", "open-square", "[",
}, {
	"chars", "close-square", "]",
}, {
	"chars", "not", "^",
}, {
	"chars", "dash", "-",
}, {
	"quantifier", "optional-not", "alias", "not", "0", "1",
}, {
	"class", "not-class-control", "^\\\\\\[\\]\\^\\-",
}, {
	"chars", "escape", "\\\\",
}, {
	"sequence", "escaped-char", "alias", "escape", "anything",
}, {
	"choice", "class-char", "alias", "not-class-control", "escaped-char",
}, {
	"sequence", "char-range", "alias", "class-char", "dash", "class-char",
}, {
	"choice", "char-or-range", "alias", "class-char", "char-range",
}, {
	"quantifier", "chars-or-ranges", "alias", "char-or-range", "0", "-1",
}, {
	"sequence", "char-class", "none", "open-square", "optional-not", "chars-or-ranges", "close-square",
}, {
	"chars", "double-quote", "\\\"",
}, {
	"class", "not-char-sequence-control", "^\\\\\\\"",
}, {
	"choice", "char-sequence-char", "alias", "not-char-sequence-control", "escaped-char",
}, {
	"quantifier", "char-sequence-chars", "alias", "char-sequence-char", "0", "-1",
}, {
	"sequence", "char-sequence", "none", "double-quote", "char-sequence-chars", "double-quote",
}, {
	"choice", "terminal", "none", "any-char", "char-class", "char-sequence",
}, {
	"class", "symbol-char", "^\\\\ \\t\\b\\f\\r\\v\\b/.\\\"\\[\\]\\^?*|():=;",
}, {
	"quantifier", "symbol-chars", "alias", "symbol-char", "1", "-1",
}, {
	"sequence", "symbol", "none", "symbol-chars",
}, {
	"chars", "open-paren", "(",
}, {
	"chars", "close-paren", ")",
}, {
	"sequence", "group", "alias", "open-paren", "wscs", "expression", "wscs", "close-paren",
}, {
	"chars", "open-brace", "{",
}, {
	"chars", "close-brace", "}",
}, {
	"class", "digit", "0-9",
}, {
	"quantifier", "digits", "alias", "digit", "1", "-1",
}, {
	"sequence", "count-quantifier", "none", "open-brace", "wscs", "digits", "wscs", "close-brace",
}, {
	"chars", "comma", ",",
}, {
	"sequence",
	"range-quantifier",
	"none",
	"open-brace",
	"wscs",
	"digits",
	"wscs",
	"comma",
	"wscs",
	"digits",
	"close-brace",
}, {
	"chars", "one-or-more", "+",
}, {
	"chars", "zero-or-more", "*",
}, {
	"chars", "zero-or-one", "?",
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
	"choice", "item", "none", "terminal", "symbol", "group", "quantifier",
}, {
	"sequence", "item-continue", "alias", "wscs", "item",
}, {
	"quantifier", "items-continue", "alias", "item-continue", "0", "-1",
}, {
	"sequence", "sequence", "none", "item", "items-continue",
}, {
	"choice",
	"expression",
	"none",
	"terminal",
	"symbol",
	"group",
	"quantifier",
	"sequence",
}, {
	"chars", "alias-word", "alias",
}, {
	"chars", "root-word", "root",
}, {
	"choice", "flag-word", "alias", "alias-word", "root-word",
}, {
	"chars", "colon", ":",
}, {
	"sequence", "flag", "alias", "colon", "flag-word",
}, {
	"quantifier", "flags", "alias", "flag", "0", "-1",
}, {
	"chars", "equal", "=",
}, {
	"sequence", "definition", "none", "symbol", "flags", "wscs", "equal", "wscs", "expression",
}, {
	"chars", "semicolon", ";",
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

func boot() {
	s, err := defineSyntax()
	if err != nil {
		log.Fatalln(err)
	}

	def, err := os.Open("syntax.p")
	if err != nil {
		log.Fatalln(err)
	}

	now := time.Now()
	_, err = s.Parse(def)
	log.Println(time.Since(now))
	if err != nil {
		log.Fatalln(err)
	}
}
