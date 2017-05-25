ws:alias  = " " | "\t" | "\n" | "\b" | "\f" | "\r" | "\v";
wsc:alias = ws | comment;

block-comment:alias   = "/*" ("*" [^/] | [^*])* "*/";
line-comment:alias    = "//" [^\n]*;
comment-segment:alias = line-comment | block-comment;
comment               = comment-segment (ws* "\n"? ws* comment-segment)*;

any-char = ".";

class-char:alias = [^\[\]\^\-] | "\\" .;
char-range:alias = class-char "-" class-char;
char-class       = "[" "^"? (class-char | char-range)* "]";

char-sequence = "\"" ([^\\\"] | "\\" .) "\"";

terminal = any-char | char-class | char-sequence;

symbol = [^\\ \n\t\b\f\r\v/.\[\]\"{}+*?|();]+;

count-quantifier = "{" wsc* [0-9]* wsc* "}";
range-quantifier = "{" wsc* [0-9]* wsc* "," wsc* [0-9]* wsc* "}";
one-or-more      = "+";
zero-or-more     = "*";
zero-or-one      = "?";
quantifier       = count-quantifier
                 | range-quantifier
		 | one-or-more
		 | zero-or-more
		 | zero-or-one;

quantified = (terminal | symbol | group) wsc* quantifier;

item     = terminal | symbol | group | quantified;
sequence = item (wsc* item)*;

element = terminal | symbol | group | quantified | sequence;
choice  = element (wsc* "|" wsc* element);

group = "(" wsc* expression wsc* ")";

expression = terminal
           | symbol
           | group
	   | quantified
	   | sequence
	   | choice;

flag       = "alias" | "root";
definition = symbol (":" flag)* wsc* "=" wsc* expression;

definitions:alias = definition (wsc* ";" (wsc | ";")* definition)*
document:root = (wsc | ";")* definitions? (wsc | ";")*;
