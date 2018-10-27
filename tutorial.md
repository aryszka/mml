# MML Tutorial

This is an early feature overview of an almost existing programming language.

## Code comment

```
// line comment

/* block comment */
```

## Numbers

An integer: `42`
A floating point number: `3.14`

## String

`"Hello, world!"`

## Boolean

`true` or `false`

## Defining a variable

`let x 42`

The equal sign is optional, could be also: `let x = 42`.

We can group multiple definitions:

```
let (
	x 42
	y 36
)
```

## Defining a mutable variable, and changing its value

```
let ~ x 42
x = 36
```

## List

Lists allow to group anything together.

`[1, 2, "and three"]`

We can access the values in a list by indexing it. This will return the second item in the list:

```
let (
	l      [1, 2, "and three"]
	second l[1]
)
```

We can take a slice of a list, which will be a new, usually smaller list:

- `list[1:2]` will contain only the second element of the original.
- `list[1:]` will contain all the elements of the original except of the first one.
- `list[:2]` will contain all the elements of the original except of the last one.
- `list[:]` will be a copy of the original.

We can use existing lists when constructing new ones:

`let newList [list..., "foo", otherList..., 4]`

In order to be able to replace the items of a list, we need to use a mutable list:

```
let mutableList ~[1, 2, "and three"]
mutableList[1] = 1
```

The list type is opaque, no algorithmic assumptions, expect acceptable or benchmark.

## String indexing and slicing

When we introduced strings above, we didn't mention indexing and slicing, while it is possible. Considering the
following definition of a string:

`let foo "foo"`

- `foo[1]` will result in `"o"`
- `foo[1:2]` will result in `"o"`, as well
- `foo[1:]` will result in `"oo"`
- `foo[:2]` will result in `"fo"`

Notice that `foo[1]` returns a string, too. Strings represent any kind of raw data, not only text.

## Structure

Structures allow to connect any value with any other value.

`{x: 42, y: 36}`

Accessing fields of a structure:

```
let (
	s {x: 42, y: 36}
	x s.x
)
```

Accessing fields by a dynamic key:

```
let (
	coords {x: 42, y: 36}
	coord  "x"
	x      coords[coord]
)
```

We can use existing structures when constructing new ones:

`let newCoords {coords..., y: 24, oldCoords: coords}`

We can use dynamic keys when constructing structures:

```
let (
	keyX   "x"
	keyY   "y"
	coords {[keyX]: 42, [keyY]: 36}
)
```

In order to be able to change fields of a structure, we need to use a mutable structure:

```
let coords ~{x: 42, y: 36}
coords.y = 24
```

## Operators

```
3 * (2 + 2)
-x
rich && windy ? sail : row
```

The rules are loosely based on the Go operators. There's no custom operators or overloading. The operands must
always have the same type except for the ternary operator (?:). We can't add an integer to a floating point
number.

Operator precedence follows the ones defined in Go. Controlling precedence is possible by grouping with parens.

## Function

`fn (x) 3 * (x + 2)`

Functions are first class values, every function is a 'lambda'. They can be passed around and assigned to
variables:

`let foo fn (x) 3 * (x + 2)`

However, the preferred way of writing the above is using the following shortcut:

`fn foo(x) 3 * (x + 2)`

Calling a function:

`foo(42)`

Just like with `let`, when defining multiple functions, it is possible to group the definitions:

```
fn (
	map(m, l)  fold(fn (c, a) [a..., m(c)], [], l)
	reverse(l) foldr(fn (c, a) [a..., c], [], l)
)
```

Functions can contain multiple statements when using a block:

```
fn foo(x) {
	let y 3 * (x + 2)
	return y
}
```

In these cases we need to use `return` in the block. Every function must have a return value.

The previously shown functions have only fixed arguments. It is possible to define functions that accept varying
number of arguments:

```
fn multiplyEach(by, ...numbers) map(fn (x) x * by, numbers)
let doubledNumbers multiplyEach(2, 1, 2, 3)
```

## Partial application

When a function call provides less arguments than the minimal number of arguments that the function can be
called with, the call returns a new, partially applied function that will expect only the missing arguments:

```
fn add(x, y) x + y

let (
	inc  add(1)
	nine inc(8)
)
```

## Chaining function calls

The following two statements are equivalent:

```
c(b(a(x)))
x -> a -> b -> c
```

Just as...

```
let (
	text io.readAll(os.stdin)
	data isError(text) ? text : json.parse(text)
)

return isError(data) ? log(data) : data
```

...is equivalent to:

`return os.stdin -> io.readAll -> errors.pass(json.parse) -> errors.only(log)`

(`log` prints all of its arguments to stderr and returns the last argument.)

## Effect

A program needs effects to be useful. Effects look and behave like functions, but must be marked with `~`.

```
fn~ printHello() {
	let message "Hello, world!"
	stdout(message)
}
```

`printHello` above is an effect because it calls `stdout`. `stdout` is a built-in effect (currently). Effects
don't need to have a return value.

Every function is an effect if any of the following conditions is true:

- accesses a mutable variable outside of its scope
- accesses a mutable list or structure defined outside of its scope
- contains channel communication
- calls other effects

Tip: try to use as few effects as possible, and try to concentrate them as close to the root of the program as
possible.

In every other way, effects are and behave just like functions.

(`log` is a special function, that is an effect but the compiler doesn't consider it as such. It's the only
special function and it is not possible to define similar ones.)

## If

```
if foo {
	return bar
} else if baz {
	return qux
} else {
	return quux
}
```

The blocks of the if statements have their own scope.

## Switch

```
switch x {
case 42:
	println("fourtytwo")
case 36:
	println("thirtysix")
default:
	println("not sure")
}
```

or equivalently:

```
switch {
case x == 42:
	println("fourtytwo")
case x == 36:
	println("thirtysix")
default:
	println("not sure")
}
```

The cases in the switch have their own scope. There's no need to use break, there's no fall through.

## Loop

Without condition:

```
for {
	if shallQuit() {
		break
	}

	println("foo")
}
```

The counterpart of `break` is `continue`. `return` also stops the loop, and returns from the enclosing function.

Loop with condition:

```
for !shallQuit() {
	println("foo")
}
```

Iterating over a list's items:

```
for item in list {
	println(item)
}
```

Iterating over a structure's values is the same:

```
for value in struct {
	println(value)
}
```

Iterating over the keys of a structure is iterating over the list of its keys:

```
for key in keys(struct) {
	println(struct[key])
}
```

Iterating over a range of numbers:

```
for number in 36:42 {
	println(number)
}
```

Partial ranges are accepted:

```
for number in :42 {
	println(number)
}
```

While the above example will count from 0 to 42, the following one only stops when explicitly told so:

```
for number in 42: {
	if number > 99 {
		break
	}

	println(number)
}
```

Iterating over the indexes of a list:

```
for i in :len(list) {
	println(list[i])
}
```

## Goroutine

This feature is borrowed from Go:

`go concurrentJob()`

## Channel

This feature is borrowed from Go, with some limitations. The syntax is also slightly different:

```
let c chan()
go (fn~ () send c 42)()
println("should be fourtytwo:", receive c)
```

Buffered channels are initialized with `bufchan`:

`let c bufchan(2)`

Closing a channel:

`close(c)`

Receiving from a closed channel will result in errClosed.

Looping over a channel will receive values from it until it's closed:

```
for value in c {
	if value == errClosed {
		break
	}

	println("received:", value)
}
```

## Select

This feature is borrowed from Go:

```
select {
case m receive messages:
	processMessage(m)
	bugged = bugged + 1
case send reports format("bugged %d times", [bugged]):
case receive stop:
	return
default:
	doWork()
}
```

## Defer

This feature is borrowed from Go.

```
fn sortIt(ctx, l) {
	let span tracing.childOf(ctx.span, "sort_it")
	defer tracing.finish(span)
	return sort(comparePrio, l)
}
```

## Error

Some functions can return errors. Errors are distinct types. They can be created with the `error` function.

`error("invalid arguments")`

They can be checked with the `isError` function.

## Panic/Recover

Considering the scope of the problems where MML tries to provide value, these features may not be included in
the first couple of revisions of the language or their future may be kept pending. They behave like in Go, except
for a syntax difference in case of `recover`:

```
fn foo() {
	defer recover(fn (err) log(err))
}
```

## Use

MML code is organized into modules. When a module requires the functionality of another module, it can import it
with the `use` keyword.

```
use "strings"
println(strings.join(", ", [1, 2, 3]))
```

It is possible to set a custom symbol for the imported modules:

```
use s "strings"
println(s.join(", ", [1, 2, 3]))
```

It is possible to import a module's members inline:

```
use . "strings"
println(join(", ", [1, 2, 3]))
```

`use` statements can be grouped, too:

```
use (
	   "strings"
	.  "ints"
	ht "net/http"
)
```

When importing a module, the top level of the imported module's statements is executed if it is imported for the first time
during the lifecycle of the program. If the top level statements of the imported module contain calls to
effects, then the use statement has to be marked with `~`.

`use ~ "config"`

It is a good practice to avoid effect calls on the top level of broadly used modules.

## Export

Only those definitions can be accessed in the imported modules that are exported:

`export fn foo() "foo"`

Grouped definitions can be exported, as well:

```
export fn (
	foo() "foo"
	bar() "bar"
)
```

## Testing

`test` is a special syntax that is considered only during the test phase:

```
fn inc(x) x + 1

test "inc" {
	test "basic" {
		test(inc(2) == 3)
		test(inc(-1) == 0)
	}

	use ints
	test("overflow", inc(ints.max) == ints.min)
}
```

## Commas and semicolons

Semicolons separate statements on the top level of a module or in a block:

`fn foo() "foo"; let f foo(); println(f)`

Commas separate:

- items of a list: `[1, 2, 3]`
- entries of a structure: `{a: 1, b: 2, c: 3}`
- parameters in a function: `fn (a, b, c) a * b * c`
- arguments in a function call: `foo(1, 2, 3)`
- definitions in definition groups: `let (a 1, b 2, c 3)`
- modules in use groups: `use ("strings", i "ints", . "lang")`

Both semicolons and commas can be replaced by new lines. E.g a function call may look like this:

```
foo(
	1
	2
	3
)
```

There is an ambiguity between a function returning an empty structure or an effect not having any statements:

- Function returning an empty structure: `fn empty() {}`
- Effect not having statements: `fn~ noopEffect() {;}`

## Lisp notation

Any MML code has an equivalent Lisp style notation. This notation is close to the Scheme flavor of Lisp. The two
representations can be converted from and to each other without changing the meaning of the code.

MML:

```
fn fold(f, i, l) len(l) == 0 ?
	[] :
	fold(f, f(l[0], i), l[1:])
```

MMLS:

```
(def (fold f i l) (if (eq (len l) 0)
  []
  (fold f (f (car l) i) (cdr l))))
```

This portability is used as a design constraint, and may not stay maintained forever.

## Built-ins

The following built-in functions are currently available:

- `len`: length of a string, list, structure, channel
- `keys`: keys of a structure
- `format`: formatted string in the style of Go's `fmt.Sprintf`
- `stdin`: reads a string from the standard input of provided length, if negative then reads to EOF, can return
  an error
- `stdout`: writes a string to the standard output, can return an error
- `stderr`: writes a string to the standard error, can return an error
- `string`: the string representation of the input argument
- `has`: true if the provided structure has the provided key
- `chan`: creates a channel
- `bufchan`: creates a buffered channel
- `close`: closes a channel
- `isBool`: true if the argument is a boolean
- `isInt`: true if the argument is an integer
- `isFloat`: true if the argument is a floating point number
- `isString`: true if the argument is a string
- `isError`: true if the argument is an error
- `error`: creates an error
- `panic`: panic in Go style
- `open`: opens a file for reading, can return an error
- `create`: creates a file for writing, can return an error
- `closeFile`: closes a file
- `args`: returns the startup arguments of the program
- `parseAST`: parses text into a raw AST with MML's syntax
- `parseInt`: parses an integer
- `parseFloat`: parses a floating point number

Many of these built-in functions will be migrated to the standard library.

## Standard library

MML currently has the following standard library modules:

- errors
- ints
- list
- log
- strings

Most of the functions in the standard library are also accessible through the bundled 'lang' module.

## Package management

MML won't have its own package management system. It will rely on either Nix or Guix, and in addition, it will
put best effort into supporting packages that are installed in a standard Unix way.

## The compiler

The compiler tansforms MML into Go or JavaScript code. The compiler does the following transient checks before
generating its output:

- every symbol is defined
- every definition is used or exported
- every function parameter is used
- every function (not effect) has a return value
- every execution path of an effect has a return value or none of them have
- every return value is used
- functions not marked as effects don't have effects
- no function parameters are changed
- only those variables are changed that are marked mutable
- only the items of mutable lists are changed
- only the values of mutable structures are changed
- every function and operator is passed only such arguments whose type the function or operator can accept
- no list index or slice range is used that is not guaranteed to fall within the length of the list
- no structure is referenced with a key that is not guaranteed to be available in the structure

The built-in functions `len`, `has` and the type checking functions, e.g. `isInt`, play a special role during
the compile time type check.

## Interpreter and REPL

MML code can be executed without compilation, it even supports shebang: `#! /usr/bin/mml`. In this case no Go or
JS installation is required. The compile time checks are applied before running the program in interpreter
mode, too.

In REPL mode, the special builtin `delete` can be used to clear definitions of the top level scope. `delete` is
only available in the REPL.
