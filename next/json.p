// JSON (http://www.json.org)
ws:alias    = [ \b\f\n\r\t];
true        = "true";
false       = "false";
null        = "null";
string      = "\"" ([^\\"\b\f\n\r\t] | "\\" (["\\/bfnrt] | "u" [0-9a-f]{4}))* "\"";
number      = "-"? ("0" | [1-9][0-9]*) ("." [0-9]+)? ([eE] [+\-]? [0-9]+)?;
entry       = string ws* ":" ws* value;
object      = "{" ws* (entry (ws* "," ws* entry)*)? ws* "}";
array       = "[" ws* (value (ws* "," ws* value)*)? ws* "]";
value:alias = true | false | null | string | number | object | array;
json        = value;

// TODO: value should be an alias but test it first like this
