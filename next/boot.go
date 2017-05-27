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
	"class", "symbol-char", "^\\\\ \\t\\b\\f\\r\\v\\b/.\\\"\\[\\]\\^?*|():=;",
}, {
	"quantifier", "symbol-chars", "alias", "symbol-char", "1", "-1",
}, {
	"sequence", "symbol", "none", "symbol-chars",
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
	"choice", "expression", "none", "symbol", "terminal",
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
