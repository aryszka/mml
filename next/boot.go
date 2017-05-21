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
	"repetition", "block-comment-body", "alias", "block-comment-char",
}, {
	"sequence",
	"block-comment",
	"alias",
	"open-block-comment",
	"block-comment-body",
	"close-block-comment",
}, {
	"repetition", "not-nls", "alias", "not-nl",
}, {
	"sequence", "line-comment", "alias", "double-slash", "not-nls",
}, {
	"choice", "comment-atom", "alias", "block-comment", "line-comment",
}, {
	"repetition", "wss", "alias", "ws",
}, {
	"optional", "optional-nl", "alias", "nl",
}, {
	"sequence",
	"continue-comment-atom",
	"alias",
	"ws",
	"optional-nl",
	"ws",
	"comment-atom",
}, {
	"repetition", "continue-comment", "alias", "continue-comment-atom",
}, {
	"sequence",
	"comment",
	"none",
	"comment-atom",
	"continue-comment",
}, {
	"choice", "wsc", "alias", "ws", "comment",
}, {
	"repetition", "wscs", "alias", "wsc",
}, {
	"class", "symbol-char", "^\\\\ \\t\\b\\f\\r\\v\\b/.\\\"\\[\\]\\^?*|():=;",
}, {
	"repetition", "symbol-chars", "alias", "symbol-char",
}, {
	"sequence", "symbol", "none", "symbol-char", "symbol-chars",
}, {
	"choice", "primitive", "alias", "symbol",
}, {
	"choice", "expression", "none", "primitive",
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
	"repetition", "flags", "alias", "flag",
}, {
	"chars", "equal", "=",
}, {
	"sequence", "definition", "none", "symbol", "flags", "wscs", "equal", "wscs", "expression",
}, {
	"chars", "semicolon", ";",
}, {
	"choice", "wsc-or-semicolon", "alias", "wsc", "semicolon",
}, {
	"repetition", "wsc-or-semicolons", "alias", "wsc-or-semicolon",
}, {
	"sequence",
	"subsequent-definition",
	"alias",
	"wscs",
	"semicolon",
	"wsc-or-semicolons",
	"definition",
}, {
	"repetition",
	"subsequent-definitions",
	"alias",
	"subsequent-definition",
}, {
	"sequence",
	"definitions",
	"alias",
	"definition",
	"subsequent-definitions",
}, {
	"optional",
	"opt-definitions",
	"alias",
	"definitions",
}, {
	"sequence",
	"document",
	"none",
	"wsc-or-semicolons",
	"opt-definitions",
	"wsc-or-semicolons",
}}

func defineSyntax() (*Syntax, error) {
	s := NewSyntax(Options{Trace: NewTrace(TraceOff)})
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
