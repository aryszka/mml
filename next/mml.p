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
string = "\"" ([^\\\"] | "\\" .)* "\"";

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
dynamic-symbol          = "symbol" wsc* "(" wscnl* expression wscnl* ")";
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
list-sep:alias        = wsc* ("," | "\n") (wscnl | ",")*;
list-item:alias       = expression | spread-expression;
expression-list:alias = list-item (list-sep list-item)*;

// list example: [1, 2, 3]
// lists can be constructed with other lists: [l1..., l2...]
list-fact:alias = "[" (wscnl | ",")* expression-list? (wscnl | ",")* "]";
list            = list-fact;
mutable-list    = "~" wscnl* list-fact;

indexer-symbol    = "[" wscnl* expression wscnl* "]";
entry             = (symbol-expression | indexer-symbol) wscnl* ":" wscnl* expression;
entry-list:alias  = (entry | spread-expression) (list-sep (entry | spread-expression))*;
struct-fact:alias = "{" (wscnl | ",")* entry-list? (wscnl | ",")* "}";
struct            = struct-fact;
mutable-struct    = "~" wscnl* struct-fact;

channel = "<>" | "<" wscnl* int wscnl* ">";

and-expression:doc = "and" wsc* "(" (wscnl | ",")* expression-list? (wscnl | ",")* ")";
or-expression:doc  = "or" wsc* "(" (wscnl | ",")* expression-list? (wscnl | ",")* ")";

// TODO: use collect
argument-list:alias = static-symbol (list-sep static-symbol)*;
collect-symbol      = "..." wscnl* static-symbol;
function-fact:alias = "(" (wscnl | ",")*
                      argument-list?
                      (wscnl | ",")*
                      collect-symbol?
                      (wscnl | ",")* ")" wscnl*
                      expression;
function            = "fn" wscnl* function-fact;
effect              = "fn" wscnl* "~" wscnl* function-fact;

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
range-expression:alias   = range-from? wscnl* ":" wscnl* range-to?;
indexer-expression:alias = expression | range-expression;
expression-indexer:alias = primary-expression wsc* "[" wscnl* indexer-expression wscnl* "]";
symbol-indexer:alias     = primary-expression wscnl* "." wscnl* symbol-expression; // TODO: test with a float on a new line
indexer                  = expression-indexer | symbol-indexer;

// TODO: implement doc flag and test and() and or()
function-application = primary-expression wsc* "(" (wscnl | ",")* expression-list? (wscnl | ",")* ")";

// if = "if" wscnl* expression wscnl* block
//      (wscnl* "else" wscnl* "if" wscnl* expression wscnl* block)*
//      (wscnl* "else" wscnl* block)?;
// 
// case-sep:alias = wsc* "\n" (wsc | "\n")*;
// default        = "default" wscnl* ":" (wscnl | ";")* statements? (wscnl | ";")*;
// case           = "case" wscnl* expression wscnl* ":" (wscnl | ";")* statements? (wscnl | ";")*;
// cases:alias    = case (case-sep case);
// switch         = "switch" wscnl* expression? wscnl* "{" wscnl*
//                  cases? (case-sep default)? (case-sep cases)?
//                  wscnl* "}";
// 
// pattern-case        = "case" wscnl* type-expression? wscnl* ":" (wscnl | ";")* statements? (wscnl | ";")*;
// pattern-cases:alias = pattern-case (case-sep pattern-case);
// match               = "match" wscnl* expression wscnl* "{" wscnl*
//                       pattern-cases? (case-sep default)? (case-sep pattern-cases)?
//                       wscnl* "}";

// conditional:alias = if
//                     // | switch
//                     // | match
// 	            ;

// receive-call:alias       = "receive" wsc* "(" (wscnl | ",")* expression (wscnl | ",")* ")";
// receive-call-op:alias    = "<-" wsc* primary-expression;
// receive-call-group:alias = "(" wscnl* receive-expression wscnl* ")";
// receive-expression       = receive-call | receive-call-op | receive-call-group;
// receive-assign-capture:alias    = assignable wscnl* ("=" wscnl*)? receive-expression;
// receive-assign-capture-grouped:alias = "(" wscnl*
//                                        assignable wscnl* ("=" wscnl*)? receive-expression
//                                        wscnl* ")";
// receive-assignment       = "set" wscnl* (receive-assign-capture | receive-assign-capture-grouped);
// receive-assignment-equal = assignable wscnl* "=" wscnl*
// receive-capture = symbol-expression wscnl* ("=" wscnl*)? expression;
// receive-capture-grouped = "(" wscnl* symbol-expression wscnl* ("=" wscnl*)? expression wscnl* ")";
// receive-definition       = "let" wscnl* (receive-capture | receive-capture-grouped);
// receive-mutable-definition = "let" wcnl* "~" wscnl* (receive-capture | receive-capture-grouped);
// receive-statement        = receive-assignment | receive-definition;
// send-call:alias          = "send" wsc* "(" (wscnl | ",")* expression list-sep expression (wscnl | ",")* ")";
// send-call-op:alias       = primary-expression wsc* "<-" wsc* expression;
// send-call-group:alias    = "(" wscnl* send wscnl* ")";
// send                     = send-call | send-call-op | send-call-group;
// close                    = "close" wsc* "(" (wscnl | ",")* expression (wscnl | ",")* ")";
// communication:alias      = receive-expression | receive-statement | send | communication-group;
// communication-group:alias = "(" wscnl* communication wscnl* ")";
// select-case              = "case" wscnl* communication wscnl* ":" (wscnl | ";")* statements? (wscnl | ";")*;
// select-cases:alias       = select-case (case-sep select-case)*;
// select                   = "select" wscnl* "{"
//                            select-cases? (case-sep default)? (case-sep select-cases)?
//                            wscnl* "}";
// go                       = "go" wscnl* function-application;
// 
// require-expression = "require" wscnl* string;
// 
// panic   = "panic" wsc* "(" (wscnl | ",")* expression (wscnl | ",")* ")";
// recover = "recover" wsc* "(" (wscnl | ",")* ")";
// 
// block = "{" (wscnl | ";")* statements? (wscnl | ";")* "}";
// expression-group = "(" wscnl* expression wscnl* ")";

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
                         // | conditional // pseudo-expression
                         // | receive-call
                         // | select // pseudo-expression
                         // | require-expression
                         // | recover
                         // | block // pseudo-expression
                         // | expression-group;
                         ;

// plus                 = "+";
// minus                = "-";
// logical-not          = "!";
// binary-not           = "^";
// unary-operator:alias = plus | minus | logical-not | binary-not;
// unary-expression     = unary-operator wsc* primary-expression | receive-call-op;
// 
// mul        = "*";
// div        = "/";
// mod        = "%";
// lshift     = "<<";
// rshift     = ">>";
// binary-and = "&";
// and-not    = "&^";
// 
// add       = "+";
// sub       = "-";
// binary-or = "|";
// xor       = "^";
// 
// eq            = "==";
// not-eq        = "!=";
// less          = "<";
// less-or-eq    = "<=";
// greater       = ">";
// greater-or-eq = ">=";
// 
// logical-and = "&&";
// logical-or  = "||";
// 
// chain = "->";
// 
// binary-op0:alias = mul | div | mod | lshift | rshift | binary-and | and-not;
// binary-op1:alias = add | sub | binary-or | xor;
// binary-op2:alias = eq | not-eq | less | less-or-eq | greater | greater-or-eq;
// binary-op3:alias = logical-and;
// binary-op4:alias = logical-or;
// binary-op5:alias = chain;
// 
// operand0:alias = primary-expression | unary-expression;
// operand1:alias = operand0 | binary0;
// operand2:alias = operand1 | binary1;
// operand3:alias = operand2 | binary2;
// operand4:alias = operand3 | binary3;
// operand5:alias = operand4 | binary4;
// 
// binary0 = operand0 wsc* binary-op0 wsc* operand0;
// binary1 = operand1 wsc* binary-op1 wsc* operand1;
// binary2 = operand2 wsc* binary-op2 wsc* operand2;
// binary3 = operand3 wsc* binary-op3 wsc* operand3;
// binary4 = operand4 wsc* binary-op4 wsc* operand4;
// binary5 = operand5 wsc* binary-op5 wsc* operand5;
// 
// binary-expression:alias = binary0 | binary1 | binary2 | binary3 | binary4 | binary5;

// TODO: this cannot be a primary expression
ternary-item:alias = primary-expression; // | unary-expression | binary-expression;
ternary-expression = ternary-item wscnl* "?" wscnl* ternary-item wscnl* ":" wscnl* ternary-item;

expression:alias = primary-expression
                 // | unary-expression
                 // | binary-expression;
		 | ternary-expression
                 ;

// TODO: code()
// TODO: observability

// break              = "break";
// continue           = "continue";
// loop-control:alias = break | continue;
// 
// in-expression   = static-symbol wscnl* "in" wscnl* (expression | range-expression);
// loop-expression = expression | in-expression;
// loop            = "for" wscnl* (loop-expression wscnl*)? block;
// 
// assignable:alias     = symbol-expression | indexer;
// assign-capture:alias = assignable wscnl* ("=" wscnl*)? expression;
// assign-set:alias     = "set" wscnl* assign-capture;
// assign-equal:alias   = assignable wscnl* "=" wscnl* expression;
// assign-captures:alias = assign-capture (list-sep assign-capture)*;
// assign-group:alias         = "set" wscnl* "(" (wscnl | ",")* assign-captures? (wscnl | ",")* ")";
// assignment           = assign-set | assign-equal | assign-group;
// 
// value-capture:alias = symbol-expression wscnl* ("=" wscnl*)? expression;
// value-definition = "let" wscnl* value-capture;
// mutable-definition = "~" "let" wscnl* value-capture;
// mutable-capture:alias = "~" wscnl* symbol-expression wscnl* ("=" wscnl*)? expression;
// value-captures:alias           = value-capture (list-sep value-capture)*;
// mixed-captures:alias     = (value-capture | mutable-capture) (list-sep (value-capture | mutable-capture))*;
// value-definition-group   = "let" wscnl* "(" (wscnl | ",")* mixed-captures? (wscnl | ",")* ")";
// mutable-definition-group = "let" wscnl* "~" wscnl* "(" (wscnl | ",")* captures? (wscnl | ",")* ")";
// 
// function-definition-fact:alias        = static-symbol wscnl* function-fact;
// effect-definition-fact:alias          = "~" wscnl* function-definition-fact;
// function-definition                   = "fn" wscnl* function-definition-fact;
// effect-definition                     = "fn" wscnl* effect-definition-fact;
// function-definition-facts:alias       = function-definition-fact (list-sep function-definition-fact)*;
// mixed-function-definition-facts:alias = (function-definition-fact | effect-definition-fact)
//                                         (list-sep (function-definition-fact | effect-definition-fact))*;
// function-definition-group             = "fn" wscnl* "(" (wscnl | ",")*
//                                         mixed-function-definition-facts?
//                                         (wscnl | ",")* ")";
// effect-definition-group               = "fn" wscnl* "~" wscnl* "(" (wscnl | ",")*
//                                         function-definition-facts
//                                         (wscnl | ",")* ")";
// 
// definition = value-definition
//            | mutable-definition
//            | value-definition-group
//            | mutable-definition-group
//            | function-definition
//            | effect-definition
//            | function-definition-group
//            | effect-definition-group;
// 
// int-type    = "int";
// float-type  = "float";
// string-type = "string";
// bool-type   = "bool";
// error-type  = "error";
// 
// primitive-type:alias = int-type
//                      | float-type
//                      | string-type
//                      | bool-type
//                      | error-type;
// 
// collect-type-expression = (static-symbol wscnl*)? "..." wscnl* type-expression;
// item-types:alias        = type-expression (list-sep type-expression)*;
// list-type-fact:alias    = "[" (wscnl | ",")*
//                           item-types?
//                           (wscnl | ",")*
//                           collect-type-expression? // TODO: can this be in other position, too?
//                           (wscnl | ",")*
//                           (":" wscnl* int (":" wscnl* int)?)?
//                           wscnl* "]";
// list-type               = list-type-fact;
// mutable-list-type       = "~" wscnl* list-type-fact;
// 
// entry-type             = static-symbol wscnl* ":" wscnl* type-expression;
// entry-types:alias      = entry-type (list-sep entry-type)*;
// struct-type-fact:alias = "{" (wscnl | ",")* entry-types? (wscnl | ",")* "}";
// struct-type            = struct-type-fact;
// mutable-struct-type    = "~" wscnl* struct-type-fact;
// 
// function-type-fact:alias = "(" (wscnl | ",")*
//                            item-types?
//                            (wscnl | ",")*
//                            collect-type-expression?
//                            (wscnl | ",")* ")" wscnl*
//                            type-expression?;
// function-type            = "fn" wscnl* function-type-fact;
// effect-type              = "fn" wscnl* "~" wscnl* function-type-fact;
// 
// type-fact:alias = primitive-type
//                 | list-type
//                 | mutable-list-type
//                 | struct-type
//                 | mutable-struct-type
//                 | function-type
//                 | effect-type;
// 
// type-union = type-fact (wscnl* "|" wscnl* type-fact)*;
// 
// type-expression = (static-symbol wscnl*)? type-union | static-symbol;
// type-constraint = "type" wscnl* static-symbol wscnl* type-union;
// type-alias      = "type" wscnl* "alias" wscnl* static-symbol wscnl* type-union;
// 
// require-statement = "require" wscnl* static-symbol wscnl* ("=" wscnl*)? string;
// 
// statement-group = "(" wscnl* statement wscnl* ")";

statement:alias = expression
                // | loop-control
                // | loop
                // | send
                // | close
                // | panic
                // | assignment
                // | definition
                // | type-constraint
                // | type-alias
                // | require-statement
                // | statement-group;
                ;

shebang-command  = [^\n]*;
shebang          = "#!" shebang-command "\n";
sep:alias        = wsc* (";" | "\n") (wscnl | ";")*;
statements:alias = statement (sep statement)*;
mml:root         = shebang? (wscnl | ";")* statements? (wscnl | ";")*;
