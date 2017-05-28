ws:alias  = " " | "\t" | "\n" | "\b" | "\f" | "\r" | "\v";
wsc:alias = ws | comment;

block-comment:alias   = "/*" ("*" [^/] | [^*])* "*/";
line-comment:alias    = "//" [^\n]*;
comment-segment:alias = line-comment | block-comment;
comment               = comment-segment (ws* "\n"? ws* comment-segment)*;

any-char:alias = ".";

class-char:alias = [^\\\[\]\^\-] | "\\" .;
char-range:alias = class-char "-" class-char;
char-class:alias = "[" "^"? (class-char | char-range)* "]";

char-sequence:alias = "\"" ([^\\\"] | "\\" .)* "\"";

terminal = (any-char | char-class | char-sequence)+;

symbol = [^\\ \n\t\b\f\r\v/.\[\]\"{}+*?|();]+;

group:alias = "(" wsc* expression wsc* ")";

count-quantifier = "{" wsc* [0-9]+ wsc* "}";
range-quantifier = "{" wsc* [0-9]+ wsc* "," wsc* [0-9]+ wsc* "}";
one-or-more      = "+";
zero-or-more     = "*";
zero-or-one      = "?";
quantity:alias   = count-quantifier
                 | range-quantifier
		 | one-or-more
		 | zero-or-more
		 | zero-or-one;

quantifier = (terminal | symbol | group) wsc* quantity;

item:alias = terminal | symbol | group | quantifier;
sequence   = item (wsc* item)*;

element:alias = terminal | symbol | group | quantifier | sequence;
choice        = element (wsc* "|" wsc* element)+;

expression = terminal
           | symbol
           | group
	   | quantifier
	   | sequence
	   | choice;

flag       = "alias" | "root";
definition = symbol (":" flag)* wsc* "=" wsc* expression;

definitions:alias = definition (wsc* ";" (wsc | ";")* definition)*
document:root     = (wsc | ";")* definitions? (wsc | ";")*;
