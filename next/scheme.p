// TODO: comment

ws:alias         = [ \b\f\n\r\t\v];
comment:alias    = ";" [^\n]*;
wsc:alias        = ws | comment;
number           = "-"? ("0" | [1-9][0-9]*) ("." [0-9]+)? ([eE] [+\-]? [0-9]+)?;
string           = "\"" ([^\\"] | "\\" .)* "\"";
symbol           = ([^\\ \n\t\b\f\r\v\"()\[\]#] | "\\" .)+;
list-form:alias  = "(" wsc* (expression wsc*)* ")"
                 | "[" wsc* (expression wsc*)* "]";
list             = list-form;
vector           = "#" list-form;
expression:alias = number | string | symbol | list;
scheme           = wsc* (expression wsc*)*;
