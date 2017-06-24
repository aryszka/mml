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
wsnlc:alias = wsc | "\n";

// comments can be line or block comments
line-comment-content  = [^\n]*;
line-comment:alias    = "//" line-comment-content;
block-comment-content = ([^*] | "*" [^/])*;
block-comment:alias   = "/*" block-comment-content "*/";
comment-part:alias    = line-comment | block-comment;
comment               = comment-part (ws* "\n"? ws* comment-part)*;

decimal-digit:alias = [0-9];
octal-digit:alias   = [0-7];
hexa-digit:alias    = [0-9a-fA-F];

// interger examples: 42, 0666, 0xfff
decimal:alias = [1-9] decimal-digit*;
octal:alias   = "0" octal-digit*;
hexa:alias    = "0" [xX] hexa-digit+;
int           = decimal | octal | hexa;

// float examples: .0, 0., 3.14, 1E-12
exponent:alias = [eE] [+\-]? decimal-digit+;
float          = decimal-digit+ "." decimal-digit* exponent?
               | "." decimal-digit+ exponent?
               | decimal-digit+ exponent;

// string example: "Hello, world!"
// only \ and " need to be escaped, e.g. allows new lines
// common escaped chars get unescaped, the rest gets unescaped to themselves
string = "\"" ([^\\"] | "\\" .)* "\"";

true       = "true";
false      = "false";
bool:alias = true | false;

// symbols normally can have only \w chars: fooBar_baz
// basic symbols cannot start with a digit
// some positions allow strings to be used as symbols, e.g: let "123" 123
// when this is not possible, dynamic symbols need to be used, but they are
// not allowed in every case, e.g: {symbol(foo()): "bar"}
// TODO: needs decision log for dynamic symbol
// TODO: exclude keywords
//
// dynamic symbol decision log:
// - every value is equatable
// - structs can act as hashtables (optimization is transparent)
// - in structs, must differentiate between symbol and value of a symbol when used as a key
// - js style [a] would be enough for the structs
// - the variables in a scope are like fields in a struct
// - [a] would be ambigous with the list as an expression
// - a logical loophole is closed with symbol(a)
// - dynamic-symbols need to be handled differently in match expressions and type expressions
symbol                  = [a-zA-Z_][a-zA-Z_0-9]*;
static-symbol:alias     = symbol | string;
dynamic-symbol          = "symbol" wsc* "(" wsnlc* expression wsnlc* ")";
symbol-expression:alias = static-symbol | dynamic-symbol;

// TODO: what happens when a dynamic symbol gets exported?

// list items are separated by comma or new line (or both)
/*
        []
        [a, b, c]
        [
                a
                b
                c
        ]
        [1, 2, a..., [b, c], [d, [e]]...]
*/
spread-expression     = primary-expression wsc* "...";
list-sep:alias        = wsc* ("," | "\n") (wsnlc | ",")*;
list-item:alias       = expression | spread-expression;
expression-list:alias = list-item (list-sep list-item)*;

// list example: [1, 2, 3]
// lists can be constructed with other lists: [l1..., l2...]
list-fact:alias = "[" (wsnlc | ",")* expression-list? (wsnlc | ",")* "]";
list            = list-fact;
mutable-list    = "~" wsnlc* list-fact;

indexer-symbol    = "[" wsnlc* expression wsnlc* "]";
entry             = (symbol-expression | indexer-symbol) wsnlc* ":" wsnlc* expression;
entry-list:alias  = (entry | spread-expression) (list-sep (entry | spread-expression))*;
struct-fact:alias = "{" (wsnlc | ",")* entry-list? (wsnlc | ",")* "}";
struct            = struct-fact;
mutable-struct    = "~" wsnlc* struct-fact;

channel = "<>" | "<" wsnlc* int wsnlc* ">";

and-expression:doc = "and" wsc* "(" (wsnlc | ",")* expression-list? (wsnlc | ",")* ")";
or-expression:doc  = "or" wsc* "(" (wsnlc | ",")* expression-list? (wsnlc | ",")* ")";

// TODO: use collect
argument-list:alias = static-symbol (list-sep static-symbol)*;
collect-symbol      = "..." wsnlc* static-symbol;
function-fact:alias = "(" (wsnlc | ",")*
                      argument-list?
                      (wsnlc | ",")*
                      collect-symbol?
                      (wsnlc | ",")* ")" wsnlc*
                      expression;
function            = "fn" wsnlc* function-fact; // can it ever cause a conflict with call and grouping?
effect              = "fn" wsnlc* "~" wsnlc* function-fact;

/*
a[42]
a[3:9]
a[:9]
a[3:]
a[b][c][d]
a.foo
a."foo"
a.symbol(foo)
*/
range-from               = expression;
range-to                 = expression;
range-expression:alias   = range-from? wsnlc* ":" wsnlc* range-to?;
indexer-expression:alias = expression | range-expression;
expression-indexer:alias = primary-expression wsc* "[" wsnlc* indexer-expression wsnlc* "]";
symbol-indexer:alias     = primary-expression wsnlc* "." wsnlc* symbol-expression; // TODO: test with a float on a new line
indexer                  = expression-indexer | symbol-indexer;

function-application = primary-expression wsc* "(" (wsnlc | ",")* expression-list? (wsnlc | ",")* ")";

if = "if" wsnlc* expression wsnlc* block
     (wsnlc* "else" wsnlc* "if" wsnlc* expression wsnlc* block)*
     (wsnlc* "else" wsnlc* block)?;

default            = "default" wsnlc* ":";
default-line:alias = default (wsnlc | ";")* statement?;
case               = "case" wsnlc* expression wsnlc* ":";
case-line:alias    = case (wsnlc | ";")* statement?;
switch             = "switch" wsnlc* expression? wsnlc* "{" (wsnlc | ";")*
                     ((case-line | default-line) (sep (case-line | default-line | statement))*)?
                     (wsnlc | ";")* "}";
// TODO: empty case not handled

int-type    = "int";
float-type  = "float";
string-type = "string";
bool-type   = "bool";
error-type  = "error";

primitive-type:alias = int-type
                     | float-type
                     | string-type
                     | bool-type
                     | error-type;

type-alias-name:alias = static-symbol;

static-range-from             = int;
static-range-to               = int;
static-range-expression:alias = static-range-from? wsnlc* ":" wsnlc* static-range-to?;
items-quantifier              = int | static-range-expression;
// TODO: maybe this can be confusing with matching constants. Shall we support matching constants, values?

items-type = items-quantifier
           | type-set (wsnlc* ":" wsnlc* items-quantifier)?
           | static-symbol wsnlc* type-set (wsnlc* ":" wsnlc* items-quantifier)?;

destructure-item = type-set | static-symbol wsnlc* type-set;

collect-destructure-item = "..." wsnlc* destructure-item?
                           (wsnlc* ":" items-quantifier)?;
list-destructure-type    = destructure-item
                           (list-sep destructure-item)*
                           (list-sep collect-destructure-item)?
                         | collect-destructure-item;
list-type-fact:alias     = "[" (wsnlc | ",")*
                           (items-type | list-destructure-type)?
                           (wsnlc | ",")* "]";
list-type                = list-type-fact;
mutable-list-type        = "~" wsnlc* list-type-fact;

destructure-match-item = match-set
                       | static-symbol wsnlc* match-set
                       | static-symbol wsnlc* static-symbol wsnlc* match-set;

collect-destructure-match-item = "..." wsnlc* destructure-match-item?
                           (wsnlc* ":" items-quantifier)?;
list-destructure-match   = destructure-match-item
                           (list-sep destructure-match-item)*
                           (list-sep collect-destructure-match-item)?
                         | collect-destructure-match-item;
list-match-fact:alias    = "[" (wsnlc | ",")*
                           (list-destructure-match | items-type)?
                           (wsnlc | ",")* "]";
list-match               = list-match-fact;
mutable-list-match       = "~" wsnlc* list-match;

entry-type             = static-symbol (wsnlc* ":" wsnlc* destructure-item)?;
entry-types:alias      = entry-type (list-sep entry-type)*;
struct-type-fact:alias = "{" (wsnlc | ",")* entry-types? (wsnlc | ",")* "}";
struct-type            = struct-type-fact;
mutable-struct-type    = "~" wsnlc* struct-type-fact;

entry-match             = static-symbol (wsnlc* ":" wsnlc* destructure-match-item)?;
entry-matches:alias     = entry-match (list-sep entry-match)*;
struct-match-fact:alias = "{" (wsnlc | ",")* entry-matches?  (wsnlc | ",")* "}";
struct-match            = struct-match-fact;
mutable-struct-match    = "~" wsnlc* struct-match-fact;

arg-type                 = type-set | static-symbol wsnlc* type-set;
args-type:alias          = arg-type (list-sep arg-type)*;
function-type-fact:alias = "(" wsnlc* args-type?  wsnlc* ")"
                            (wsc* (type-set | static-symbol wsc* type-set))?;
function-type            = "fn" wsnlc* function-type-fact;
effect-type              = "fn" wsnlc* "~" wsnlc* function-type-fact;

// TODO: heavy naming crime

receive-direction = "receive";
send-direction    = "send";
channel-type      = "<" wsnlc*
                    (receive-direction | send-direction)? wsnlc*
                    destructure-item?
                    wsnlc* ">";

type-fact-group:alias = "(" wsnlc* type-fact wsnlc* ")";
type-fact:alias = primitive-type
                | type-alias-name
                | list-type
                | mutable-list-type
                | struct-type
                | mutable-struct-type
                | function-type
                | effect-type
                | channel-type
                | type-fact-group;

type-set:alias        = type-fact (wsnlc* "|" wsnlc* type-fact)*;
type-expression:alias = type-set | static-symbol wsc* type-set;

match-fact:alias = list-match
                 | mutable-list-match
                 | struct-match
                 | mutable-struct-match;

match-set:alias        = type-set | match-fact;
match-expression:alias = match-set | static-symbol wsc* match-set;

match-case               = "case" wsnlc* match-expression wsnlc* ":";
match-case-line:alias    = match-case (wsnlc | ";")* statement?;
match                    = "match" wsnlc* expression wsnlc* "{" (wsnlc | ";")*
                           ((match-case-line | default-line)
                           (sep (match-case-line | default-line | statement))*)?
                           (wsnlc | ";")* "}";

conditional:alias = if
                  | switch
                  | match;

receive-call                    = "receive" wsc* "(" (wsnlc | ",")* expression (wsnlc | ",")* ")";
receive-op                      = "<-" wsc* primary-expression;
receive-expression-group:alias  = "(" wsnlc* receive-expression wsnlc* ")";
receive-expression:alias        = receive-call | receive-op | receive-expression-group;

receive-assign-capture:alias = assignable wsnlc* ("=" wsnlc*)? receive-expression;
receive-assignment           = "set" wsnlc* receive-assign-capture;
receive-assignment-equal     = assignable wsnlc* "=" wsnlc* receive-expression;
receive-capture:alias        = symbol-expression wsnlc* ("=" wsnlc*)? receive-expression;
receive-definition           = "let" wsnlc* receive-capture;
receive-mutable-definition   = "let" wcnl* "~" wsnlc* receive-capture;
receive-statement:alias      = receive-assignment | receive-definition;

send-call:alias       = "send" wsc* "(" (wsnlc | ",")* expression list-sep expression (wsnlc | ",")* ")";
send-op:alias         = primary-expression wsc* "<-" wsc* expression;
send-call-group:alias = "(" wsnlc* send wsnlc* ")";
send                  = send-call | send-op | send-call-group;

close = "close" wsc* "(" (wsnlc | ",")* expression (wsnlc | ",")* ")";

communication-group:alias = "(" wsnlc* communication wsnlc* ")";
communication:alias       = receive-expression | receive-statement | send | communication-group;

select-case            = "case" wsnlc* communication wsnlc* ":";
select-case-line:alias = select-case (wsnlc | ";")* statement?;
select                 = "select" wsnlc* "{" (wsnlc | ";")*
                         ((select-case-line | default-line)
                          (sep (select-case-line | default-line | statement))*)?
                         (wsnlc | ";")* "}";

go = "go" wsnlc* function-application;

/*
require . = "mml/foo"
require bar = "mml/foo"
require . "mml/foo"
require bar "mml/foo"
require "mml/foo"
require (
        . = "mml/foo"
        bar = "mml/foo"
        . "mml/foo"
        bar "mml/foo"
        "mml/foo"
)
require ()
*/
require-inline                = ".";
require-fact                  = string
                              | (static-symbol | require-inline) (wsnlc* "=")? wsnlc* string;
require-facts:alias           = require-fact (list-sep require-fact)*;
require-statement:alias       = "require" wsnlc* require-fact;
require-statement-group:alias = "require" wsc* "(" (wsnlc | ",")*
                                require-facts?
                                (wsnlc | ",")* ")";
require                       = require-statement | require-statement-group;

panic   = "panic" wsc* "(" (wsnlc | ",")* expression (wsnlc | ",")* ")";
recover = "recover" wsc* "(" (wsnlc | ",")* ")";

block                  = "{" (wsnlc | ";")* statements? (wsnlc | ";")* "}";
expression-group:alias = "(" wsnlc* expression wsnlc* ")";

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
                         | and-expression // only documentation
                         | or-expression // only documentation
                         | function
                         | effect
                         | indexer
                         | function-application // pseudo-expression
                         | conditional // pseudo-expression
                         | receive-call
                         | select // pseudo-expression
                         | recover
                         | block // pseudo-expression
                         | expression-group;

plus                 = "+";
minus                = "-";
logical-not          = "!";
binary-not           = "^";
unary-operator:alias = plus | minus | logical-not | binary-not;
unary-expression = unary-operator wsc* primary-expression | receive-op;

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

binary0 = operand0 wsc* binary-op0 wsc* operand0;
binary1 = operand1 wsc* binary-op1 wsc* operand1;
binary2 = operand2 wsc* binary-op2 wsc* operand2;
binary3 = operand3 wsc* binary-op3 wsc* operand3;
binary4 = operand4 wsc* binary-op4 wsc* operand4;
binary5 = operand5 wsc* binary-op5 wsc* operand5;

binary-expression:alias = binary0 | binary1 | binary2 | binary3 | binary4 | binary5;

ternary-expression = expression wsnlc* "?" wsnlc* expression wsnlc* ":" wsnlc* expression;

expression:alias = primary-expression
                 | unary-expression
                 | binary-expression
                 | ternary-expression;

// TODO: code()
// TODO: observability

break              = "break";
continue           = "continue";
loop-control:alias = break | continue;

in-expression   = static-symbol wsnlc* "in" wsnlc* (expression | range-expression);
loop-expression = expression | in-expression;
loop            = "for" wsnlc* (block | loop-expression wsnlc* block);

/*
a = b
set c = d
set e f
set (
        g = h
        i j
)
*/
assignable:alias      = symbol-expression | indexer;
assign-capture        = assignable wsnlc* ("=" wsnlc*)? expression;
assign-set:alias      = "set" wsnlc* assign-capture;
assign-equal          = assignable wsnlc* "=" wsnlc* expression;
assign-captures:alias = assign-capture (list-sep assign-capture)*;
assign-group:alias    = "set" wsnlc* "(" (wsnlc | ",")* assign-captures? (wsnlc | ",")* ")";
assignment            = assign-set | assign-equal | assign-group;

/*
let a = b
let c d
let ~ e = f
let ~ g h
let (
        i = j
        k l
        ~ m = n
        ~ o p
)
let ~ (
        q = r
        s t
)
*/
value-capture-fact:alias = symbol-expression wsnlc* ("=" wsnlc*)? expression;
value-capture            = value-capture-fact;
mutable-capture          = "~" wsnlc* value-capture-fact;
value-definition         = "let" wsnlc* (value-capture | mutable-capture);
value-captures:alias     = value-capture (list-sep value-capture)*;
mixed-captures:alias     = (value-capture | mutable-capture) (list-sep (value-capture | mutable-capture))*;
value-definition-group   = "let" wsnlc* "(" (wsnlc | ",")* mixed-captures? (wsnlc | ",")* ")";
mutable-definition-group = "let" wsnlc* "~" wsnlc* "(" (wsnlc | ",")* value-captures? (wsnlc | ",")* ")";

/*
fn a() b
fn ~ c() d
fn (
        e() f
        ~ g() h
)
fn ~ (
        i()
        j()
)
*/
function-definition-fact:alias = static-symbol wsnlc* function-fact;
function-capture               = function-definition-fact;
effect-capture                 = "~" wsnlc* function-definition-fact;
function-definition            = "fn" wsnlc* (function-capture | effect-capture);
function-captures:alias        = function-capture (list-sep function-capture)*;
mixed-function-captures:alias  = (function-capture | effect-capture)
                                 (list-sep (function-capture | effect-capture))*;
function-definition-group      = "fn" wsnlc* "(" (wsnlc | ",")*
                                 mixed-function-captures?
                                 (wsnlc | ",")* ")";
effect-definition-group        = "fn" wsnlc* "~" wsnlc* "(" (wsnlc | ",")*
                                 function-captures?
                                 (wsnlc | ",")* ")";

definition:alias = value-definition
                 | value-definition-group
                 | mutable-definition-group
                 | function-definition
                 | function-definition-group
                 | effect-definition-group;

// TODO: cannot do:
// type alias a int|fn () string|error
// needs grouping of type-set

type-alias      = "type" wsnlc* "alias" wsnlc* static-symbol wsnlc* type-set;
type-constraint = "type" wsnlc* static-symbol wsnlc* type-set;

statement-group:alias = "(" wsnlc* statement wsnlc* ")";

statement:alias = send
                | close
                | panic
                | require
                | loop-control
                | go
                | loop
                | assignment
                | definition
                | expression
                | type-alias
                | type-constraint
                | statement-group;

shebang-command  = [^\n]*;
shebang          = "#!" shebang-command "\n";
sep:alias        = wsc* (";" | "\n") (wsnlc | ";")*;
statements:alias = statement (sep statement)*;
mml:root         = shebang? (wsnlc | ";")* statements? (wsnlc | ";")*;
