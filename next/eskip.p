/*
Eskip routing configuration format for Skipper: https://github.com/zalando/skipper
*/

// TODO: definition with comment, doc = comment, or just replace comment

eskip:root = (expression | definitions)?;

comment-line:alias = "//" [^\n]*;
space:alias        = [ \b\f\r\t\v];
comment:alias      = comment-line (space* "\n" space* comment-line)*;

wsc:alias = [ \b\f\n\r\t\v] | comment;

decimal-digit:alias = [0-9];
octal-digit:alias   = [0-7];
hexa-digit:alias    = [0-9a-fA-F];

decimal:alias = [1-9] decimal-digit*;
octal:alias   = "0" octal-digit*;
hexa:alias    = "0" [xX] hexa-digit+;
int           = decimal | octal | hexa;

exponent:alias = [eE] [+\-]? decimal-digit+;
float          = decimal-digit+ "." decimal-digit* exponent?
               | "." decimal-digit+ exponent?
               | decimal-digit+ exponent;

number:alias = "-"? (int | float);

string = "\"" ([^\\"] | "\\" .)* "\"";
regexp = "/" ([^\\/] | "\\" .)* "/";
symbol = [a-zA-Z_] [a-zA-z0-9_]*;

arg:alias  = number | string | regexp;
args:alias = arg (wsc* "," wsc* arg)*;
term:alias = symbol wsc* "(" wsc* args? wsc* ")";

predicate        = term;
predicates:alias = "*" | predicate (wsc* "&&" wsc* predicate)*;

filter        = term;
filters:alias = filter (wsc* "->" wsc* filter)*;

address:alias = string;
shunt         = "<shunt>";
loopback      = "<loopback>";
backend:alias = address | shunt | loopback;

expression = predicates (wsc* "->" wsc* filters)? wsc* "->" wsc* backend;

id:alias   = symbol;
definition = id wsc* ":" wsc* expression;

free-sep:alias    = (wsc | ";");
sep:alias         = wsc* ";" free-sep*;
definitions:alias = free-sep* definition (sep definition)* free-sep*;
