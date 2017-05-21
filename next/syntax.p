/*
Syntax for parsing itself.
*/

// ws:alias  = " " | "\t" | "\n" | "\b" | "\f" | "\r" | "\v";
// wsc:alias = ws | comment;
// 
// block-comment:alias = "/*" ("*" [^/] | [^*])* "*/";
// line-comment:alias  = "//" [^\n]*;
// comment-atom:alias  = line-comment | block-comment;
// comment             = comment-atom (ws* "\n"? ws* comment-atom)*;
// 
// any-char:alias = ".";
// 
// char-sequence = "\"" ([^\\\"] | "\\" .) "\"";
// 
// class-char:alias = [^\[\]\^\-] | "\\" .;
// char-range:alias = class-char "-" class-char;
// char-class       = "[" "^"? (class-char | char-range)* "]";
// 
// terminal = any-char | char-sequence | char-class;
// 
// symbol-char:alias = [^\\ \t\b\f\r\v\b/.\"\[\]\^?*|():=;];
// symbol            = symbol-char symbol-char*;
// 
// optional    = primitive wsc "?";
// repetition  = primitive wsc "*";
// sequence    = primitive (wsc* primitive)*;
// choice      = (complex (wsc* "|" wsc* complex);
// group:alias = "(" wsc* expression wsc* ")";
// 
// primitive:alias = terminal
//                 | symbol
//                 | group;
// 
// copmlex:alias = optional
//               | repetition
//               | sequence;
// 
// expression = primitive
//            | complex
//            | choice;
// 
// flag       = "alias" | "root";
// definition = symbol (":" flag)* wsc* "=" wsc* expression;
// 
// definitions:alias = definition (wsc* ";" (wsc | ";")* definition)*
// document:root = (wsc | ";")* definitions? (wsc | ";")*;
