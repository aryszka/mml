ws:alias   = [ \b\f\r\t\v];
wsnl:alias = ws | "\n";

comment-line:alias = "#" [^\n]*;
comment            = comment-line (ws* "\n" ws* comment-line)*;

wsc:alias   = ws | comment-line;
wsnlc:alias = wsnl | comment-line;

quoted:alias        = "\"" ([^\\"] | "\\" .)* "\"";
symbol-non-ws:alias = ([^\\"\n=#.\[\] \b\f\r\t\v] | "\\" .)+;
symbol              = symbol-non-ws (ws* symbol-non-ws)* | quoted;

key-form:alias = symbol (ws* "." ws* symbol)*;
key            = key-form;
group-key      = (comment "\n" ws*)? "[" ws* key-form ws* "]";

value-chars:alias = ([^\\"\n=# \b\f\r\t\v] | "\\" .)+;
value             = value-chars (ws* value-chars)* | quoted;
key-val           = (comment "\n" ws*)? (key | key? ws* "=" ws* value?);

entry:alias = group-key | key-val;
doc:root    = (entry (ws* comment-line)? | wsnlc)*;

// TODO: not tested
// set as root for streaming:
single-entry = (entry (ws* comment-line)?
                | wsnlc* entry (ws* comment-line)?)
               [];
