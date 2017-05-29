// whitespace is ignored except for \n which is only ignored
// most of the time, but can serve as separator in:
// - list
// - struct
// - function args
// - statements
// - list, struct and function type constraints
ws:alias    = " " | "\b" | "\f" | "\r" | "\t" | "\v";
wsnl:alias  = ws | "\n";
wsc:alias   = ws | comment;
wscnl:alias = wsc | "\n";

// comments can be line or block comments
line-comment-content  = [^\n]*;
line-comment:alias    = "//" line-comment-content;
block-comment-content = ([^*] | "*" [^/])*;
block-comment:alias   = "/*" block-comment-content "*/";
comment-part:alias    = line-comment | block-comment;
comment               = comment-part (ws* "\n"? ws* comment-part);

decimal-digit:alias = [0-9];
octal-digit:alias   = [0-7];
hexa-digit:alias    = [0-9a-fA-F];

// interger examples: 42, 0666, 0xfff
decimal:alias = [1-9] decimal-digit*;
octal:alias   = "0" octal-digit*;
hexa:alias    = "0" [xX] hexa-digit+;
int           = decimal | octal | hexa;

// float examples: .0, 0., 3.14, 1E-12
exponent:alias = [eE] [+-]? decimal-digit+;
float          = decimal-digit+ "." decimal-digit* exponent?
               | "." decimal-digit+ exponent?
               | decimal-digit+ exponent;

// string example: "Hello, world!"
// only \ and " need to be escaped, e.g. allows new lines
// common escaped chars get unescaped, the rest gets unescaped to themselves
string = "\"" ([^\\\"] | "\\" .)* "\"";

true       = "true";
false      = "false";
bool:alias = true | false;

// symbols normally can have only \w chars: fooBar_baz
// basic symbols cannot start with a digit
// some positions allow strings to be used as symbols, e.g: let "123" 123
// when this is not possible, dynamic symbols need to be used, but they are
// not allowed in every case, e.g: {symbol(foo()): "bar"}
symbol                  = [a-zA-Z_][a-zA-Z_0-9]*;
static-symbol:alias     = symbol | string;
dynamic-symbol          = "symbol" wscnl* "(" wscnl* expression wscnl* ")";
spread-symbol           = static-symbol "...";
dynamic-spread-symbol   = dynamic-symbol "...";
symbol-expression:alias = static-symbol | dynamic-symbol;
spread-expression:alias = spread-symbol | dynamic-spread-symbol;

// list items are separated by comma or new line (or both)
list-sep:alias        = wsc* ("," | "\n") (wscnl | ",")*;
list-item:alias       = expression | spread-expression;
expression-list:alias = list-item (list-sep list-item)*;

// list example: [1, 2, 3]
// lists can be constructed with other lists: [l1..., l2...]
list-fact:alias = "[" (wscnl | ",")* expression-list? (wscnl | ",")* "]";
list            = list-fact;
mutable-list    = "~" wscnl* list-fact;

entry             = symbol-expression wscnl* ":" wscnl* expression | spread-expression;
entry-list:alias  = entry (list-sep entry)*;
struct-fact:alias = "{" (wscnl | ",")* entry-list? (wscnl | ",")* "}";
struct            = struct-fact;
mutable-struct    = "~" wscnl* struct-fact;

channel = "<>" | "<" wscnl* int wscnl* ">";

and-expression = "and" wscnl* "(" wscnl* expression-list? wscnl* ")";
or-expression  = "or" wscnl* "(" wscnl* expression-list? wscnl* ")";

argument:alias      = static-symbol;
argument-list:alias = argument (list-sep argument)*;
function-fact:alias = "(" (wscnl | ",")*
                      argument-list?
		      (wscnl | ",")*
		      spreac-symbol?
		      (wscnl | ",")* ")" wscnl*
		      expression;
function            = "fn" wscnl* function-fact;
effect              = "fn" wscnl* "~" wscnl* function-fact;

range-expression         = expression? wscnl* ":" wscnl* expression?; // not sure which one is which
indexer-expression:alias = expression | range-expression;
expression-indexer:alias = expression wsc* "[" wscnl* indexer-expression wscnl* "]";
symbol-indexer           = expression wscnl* "." wscnl* symbol-expression;
indexer:alias            = expression-indexer | symbol-indexer;

function-call = expression wsc* "(" wscnl* expression-list? wscnl* ")";

tertiary-if = expression wscnl* "?" wscnl* expression wscnl* ":" wscnl* expression;

if = "if" wscnl* expression wscnl* block
     (wscnl* "else" wscnl* "if" wscnl* expression wscnl* block)*
     (wscnl* "else" wscnl* block)?;

case-sep:alias = wsc* "\n" (wsc | "\n")*;
default        = "default" wscnl* ":" wscnl* statements?;
case           = "case" wscnl* expression wscnl* ":" statements?;
cases:alias    = case (case-sep case);
switch         = "switch" wscnl* expression? wscnl* "{" wscnl*
                 cases? case-sep default? case-sep cases?
                 wscnl* "}";

int-type    = "int";
float-type  = "float";
string-type = "string";
bool-type   = "bool";

primitive-type:alias = int-type | float-type | string-type | bool-type | error-type;

collect-type-expression = (static-symbol wscnl*)? "..." wscnl* type-expression;
item-types:alias        = type-expression (list-sep type-expression)*;
list-type-fact:alias    = "[" (wscnl | ",")*
                          item-types?
			  (wscnl | ",")*
			  collect-type-expression?
			  (wscnl | ",")*
			  (":" wscnl* int (":" wscnl* int)?)?
			  wscnl* "]";
list-type               = list-type-fact;
mutable-list-type       = "~" wscnl* list-type-fact;

entry-type             = static-symbol wscnl* ":" wscnl* type-expression;
entry-types:alias      = entry-type (list-sep entry-type)*;
struct-type-fact:alias = "{" wscnl* entry-types? wscnl* "}";
struct-type            = struct-type-fact;
mutable-struct-type    = "~" wscnl* struct-type-fact;

function-type-fact:alias = "(" (wscnl | ",")*
                           item-types?
			   (wscnl | ",")*
			   collect-type-expression?
			   (wscnl | ",")* ")" wscnl*
                           type-expression?;
function-type            = "fn" wscnl* function-type-fact;
effect-type              = "fn" wscnl* "~" wscnl* function-type-fact;

type-fact:alias = primitive-type
                | list-type
                | mutable-list-type
                | struct-type
                | mutable-struct-type
                | function-type
                | effect-type;

type-union = type-fact (wscnl* "|" wscnl* type-fact)*;

type-expression = (static-symbol wscnl*)? (type-fact | type-union)
                | static-symbol;

pattern-case        = "case" wscnl* type-expression? wscnl* ":" statements?;
pattern-cases:alias = pattern-case (case-sep pattern-case);
match               = "match" wscnl* expression wscnl* "{" wscnl*
                      pattern-cases? case-sep default? case-sep pattern-cases?
                      wscnl* "}";

conditional = tertiary-if
            | if
            | switch
            | match;

plus                 = "+";
minus                = "-";
logical-not          = "!";
binary-not           = "^";
receive-operator     = "<-";
unary-operator:alias = plus | minus | logical-not | binary-not | receive-operator;
unary-expression     = unary-operator wscnl* primary-expression;

mul        = "*";
div        = "/";
mod        = "%";
lshift     = "<<";
rshift     = ">>";
binary-and = "&";
and-not    = "&^";

add       = "+";
sub       = "-";
binary-or = "|";
xor       = "^";

eq            = "==";
not-eq        = "!=";
less          = "<";
less-or-eq    = "<=";
greater       = ">";
greater-or-eq = ">=";

logical-and = "&&";
logical-or  = "||";

chain = "->";

binary-op0:alias = mul | div | mod | lshift | rshift | binary-and | and-not;
binary-op1:alias = add | sub | binary-or | xor;
binary-op2:alias = eq | not-eq | less | less-or-eq | greater | greater-or-eq;
binary-op3:alias = logical-and;
binary-op4:alias = logical-or;
binary-op5:alias = chain;

operand0:alias = primary-expression | unary-expression;
operand1:alias = operand0 | binary0;
operand2:alias = operand1 | binary1;
operand3:alias = operand2 | binary2;
operand4:alias = operand3 | binary3;
operand5:alias = operand4 | binary4;

binary0 = operand0 binary-op0 operand0;
binary1 = operand1 binary-op1 operand1;
binary2 = operand2 binary-op2 operand2;
binary3 = operand3 binary-op3 operand3;
binary4 = operand4 binary-op4 operand4;
binary5 = operand5 binary-op5 operand5;

binary-expression:alias = binary0 | binary1 | binary2 | binary3 | binary4 | binary5;

primary-expression:alias = int
                         | float
                         | string
                         | bool
                         | symbol
                         | dynamic-symbol
                         | list
                         | mutable-list
                         | struct
                         | mutable-struct
                         | channel
                         | and-expression
                         | or-expression
                         | function
                         | effect
                         | indexer
                         | function-call
                         | conditional
                         | send-call
                         | receive-call
                         | close
                         | select
                         | require-expression
                         | block
                         | expression-group;

expression = primary-expression
           | unary-expression
           | binary-expression;

statement = expression
          | assignment
          | loop-control
          | loop
          | panic-call
          | recover-call
          | definition
          | type-constraint
          | type-alias
          | require-statement;

shebang-command  = [^\n]*
shebang          = "#!" shebang-command "\n"
sep:alias        = wsc* (";" | "\n") (wscnl | ";")*;
statements:alias = statement (sep statement)*;
mml:root         = shebang? (wscnl | ";")* statements? (wscnl | ";")*;
