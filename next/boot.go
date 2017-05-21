package next

import (
	"log"
	"os"
	"time"
)

var definitions = [][]string{{
	"chars", "space", " ",
}, {
	"chars", "tab", "\t",
}, {
	"chars", "nl", "\n",
}, {
	"chars", "backspace", "\b",
}, {
	"chars", "formfeed", "\f",
}, {
	"chars", "carryreturn", "\r",
}, {
	"chars", "verticaltab", "\v",
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
	"class", "not-nl", "^\n",
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
	"chars", "semicolon", ";",
}, {
	"choice", "wsc-or-semicolon", "alias", "wsc", "semicolon",
}, {
	"repetition", "wsc-or-semicolons", "alias", "wsc-or-semicolon",
}, {
	"sequence",
	"document",
	"none",
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
