package mml

import (
	"reflect"
	"testing"
)

func testEnvEq(left, right *env) bool {
	return true
}

func testEq(left, right interface{}) bool {
	switch lt := left.(type) {
	case list:
		rt := right.(list)

		if lt.mutable != rt.mutable {
			return false
		}

		if len(lt.values) != len(rt.values) {
			return false
		}

		for i := range lt.values {
			if lt.values[i] != rt.values[i] {
				return false
			}
		}
	case structure:
		rt := right.(structure)

		if lt.mutable != rt.mutable {
			return false
		}

		if len(lt.values) != len(rt.values) {
			return false
		}

		for k := range lt.values {
			if lt.values[k] != rt.values[k] {
				return false
			}
		}
	case chan interface{}:
		rt := right.(chan interface{})
		return reflect.ValueOf(lt).Cap() == reflect.ValueOf(rt).Cap()
	case statementList:
		rt := right.(statementList)
		if len(lt.statements) != len(rt.statements) {
			return false
		}

		for i := range lt.statements {
			if !testEq(lt.statements[i], rt.statements[i]) {
				return false
			}
		}
	case function:
		rt := right.(function)

		if lt.effect != rt.effect {
			return false
		}

		if len(lt.params) != len(rt.params) {
			return false
		}

		for i := range lt.params {
			if lt.params[i] != rt.params[i] {
				return false
			}
		}

		if lt.collectParam != rt.collectParam {
			return false
		}

		if !testEq(lt.statement, rt.statement) {
			return false
		}

		if !testEnvEq(lt.env, rt.env) {
			return false
		}
	case functionApplication:
		rt := right.(functionApplication)

		if !testEq(lt.function, rt.function) {
			return false
		}

		if len(lt.args) != len(rt.args) {
			return false
		}

		for i := range lt.args {
			if !testEq(lt.args[i], rt.args[i]) {
				return false
			}
		}
	default:
		return left == right
	}

	return true
}

func testEvalStatement(text string, value interface{}) func(*testing.T) {
	return func(t *testing.T) {
		c, err := parseStatement(text)
		if err != nil {
			t.Fatal(err)
		}

		v, err := eval(newEnv(), c)
		if err != nil {
			t.Fatal(err)
		}

		if !testEq(v, value) {
			t.Errorf("got: %v, expected: %v", v, value)
		}
	}
}

func TestEval(t *testing.T) {
	t.Run("empty module, whitespace", func(t *testing.T) {
		m, err := parseModule("  \t ")
		if err != nil {
			t.Fatal(err)
		}

		err = evalModule(newEnv(), m)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("integer", func(t *testing.T) {
		t.Run("octal", testEvalStatement("052", 42))
		t.Run("decimal", testEvalStatement("42", 42))
		t.Run("hexa", testEvalStatement("0x2a", 42))
	})

	t.Run("float", func(t *testing.T) {
		t.Run("float", testEvalStatement("42.42", 42.42))
		t.Run("leading zero", testEvalStatement("0.42", .42))
		t.Run("no leading zero", testEvalStatement(".42", .42))
		t.Run("exponent lower", testEvalStatement("42e85", 42e85))
		t.Run("exponent upper", testEvalStatement("42E85", 42e85))
		t.Run("exponent +", testEvalStatement("42e+85", 42e85))
		t.Run("exponent -", testEvalStatement("42e-85", 42e-85))
		t.Run("exponent leading zero", testEvalStatement("0.42e-85", .42e-85))
		t.Run("exponent no leading zero", testEvalStatement(".42e-85", .42e-85))
	})

	t.Run("string", func(t *testing.T) {
		t.Run("string", testEvalStatement(`"Hello, world!"`, "Hello, world!"))
		t.Run("string multiline", testEvalStatement("\"Hello,\nworld!\"", "Hello,\nworld!"))
		t.Run("string escaped multiline", testEvalStatement("\"Hello,\\nworld!\"", "Hello,\nworld!"))
		t.Run("string escaped", testEvalStatement("\"\\\"Hello, world!\\\"\"", `"Hello, world!"`))
	})

	t.Run("bool", func(t *testing.T) {
		t.Run("true", testEvalStatement("true", true))
		t.Run("false", testEvalStatement("false", false))
	})

	t.Run("list", func(t *testing.T) {
		t.Run("empty", testEvalStatement("[]", list{}))
		t.Run("simple values", testEvalStatement("[1, 2, 3]", list{values: []interface{}{1, 2, 3}}))
		t.Run("spread", testEvalStatement(
			"[1, 2, [3, 4]..., 5]",
			list{values: []interface{}{1, 2, 3, 4, 5}},
		))
	})

	t.Run("mutable list", func(t *testing.T) {
		t.Run("empty", testEvalStatement("~[]", list{mutable: true}))
		t.Run("simple values", testEvalStatement(
			"~[1, 2, 3]",
			list{mutable: true, values: []interface{}{1, 2, 3}},
		))
		t.Run("spread", testEvalStatement(
			"~[1, 2, [3, 4]..., 5]",
			list{mutable: true, values: []interface{}{1, 2, 3, 4, 5}},
		))
	})

	t.Run("struct", func(t *testing.T) {
		t.Run("empty", testEvalStatement("{}", structure{}))
		t.Run("simple", testEvalStatement(
			"{foo: 1, bar: 2, baz: 3}",
			structure{values: map[string]interface{}{
				"foo": 1,
				"bar": 2,
				"baz": 3,
			}},
		))
		t.Run("string key", testEvalStatement(
			`{"foo": 42}`,
			structure{values: map[string]interface{}{
				"foo": 42,
			}},
		))
		t.Run("expression key", testEvalStatement(
			`{["foo"]: 42}`,
			structure{values: map[string]interface{}{
				"foo": 42,
			}},
		))
		t.Run("spread", testEvalStatement(
			"{foo: 1, bar: 2, {baz: 3, qux: 4}..., quux: 5}",
			structure{values: map[string]interface{}{
				"foo":  1,
				"bar":  2,
				"baz":  3,
				"qux":  4,
				"quux": 5,
			}},
		))
	})

	t.Run("mutable struct", func(t *testing.T) {
		t.Run("empty", testEvalStatement("~{}", structure{mutable: true}))
		t.Run("simple", testEvalStatement(
			"~{foo: 1, bar: 2, baz: 3}",
			structure{
				mutable: true,
				values: map[string]interface{}{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
			},
		))
		t.Run("string key", testEvalStatement(
			`~{"foo": 42}`,
			structure{
				mutable: true,
				values: map[string]interface{}{
					"foo": 42,
				},
			},
		))
		t.Run("expression key", testEvalStatement(
			`~{["foo"]: 42}`,
			structure{
				mutable: true,
				values: map[string]interface{}{
					"foo": 42,
				},
			},
		))
		t.Run("spread", testEvalStatement(
			"~{foo: 1, bar: 2, {baz: 3, qux: 4}..., quux: 5}",
			structure{
				mutable: true,
				values: map[string]interface{}{
					"foo":  1,
					"bar":  2,
					"baz":  3,
					"qux":  4,
					"quux": 5,
				},
			},
		))
	})

	t.Run("channel", func(t *testing.T) {
		t.Run("unbuffered", testEvalStatement("<>", make(chan interface{})))
		t.Run("buffered", testEvalStatement("<42>", make(chan interface{}, 42)))
	})

	t.Run("function", func(t *testing.T) {
		t.Run("no params, void", testEvalStatement(
			"fn () {;}",
			function{statement: statementList{}},
		))
		t.Run("no params, expression", testEvalStatement(
			"fn () 42",
			function{statement: 42},
		))
		t.Run("no params, block", testEvalStatement(
			"fn () { return 42 }",
			function{statement: statementList{statements: []interface{}{ret{value: 42}}}},
		))
		t.Run("void", testEvalStatement(
			"fn (foo, bar, baz) {;}",
			function{
				statement: statementList{},
				params:    []string{"foo", "bar", "baz"},
			},
		))
		t.Run("expression", testEvalStatement(
			"fn (foo, bar, baz) 42",
			function{
				params:    []string{"foo", "bar", "baz"},
				statement: 42,
			},
		))
		t.Run("block", testEvalStatement(
			"fn (foo, bar, baz) { return 42 }",
			function{
				params:    []string{"foo", "bar", "baz"},
				statement: statementList{statements: []interface{}{ret{value: 42}}},
			},
		))
		t.Run("collect, void", testEvalStatement(
			"fn (foo, bar, baz, ...qux) {;}",
			function{
				statement:    statementList{},
				params:       []string{"foo", "bar", "baz"},
				collectParam: "qux",
			},
		))
		t.Run("collect, expression", testEvalStatement(
			"fn (foo, bar, baz, ...qux) 42",
			function{
				params:       []string{"foo", "bar", "baz"},
				collectParam: "qux",
				statement:    42,
			},
		))
		t.Run("collect, block", testEvalStatement(
			"fn (foo, bar, baz, ...qux) { return 42 }",
			function{
				params:       []string{"foo", "bar", "baz"},
				collectParam: "qux",
				statement:    statementList{statements: []interface{}{ret{value: 42}}},
			},
		))
		t.Run("struct", testEvalStatement(
			"fn () {}",
			function{
				statement: structure{},
			},
		))
	})

	t.Run("effect", func(t *testing.T) {
		t.Run("no params, void", testEvalStatement(
			"fn~ () {;}",
			function{statement: statementList{}, effect: true},
		))
		t.Run("no params, expression", testEvalStatement(
			"fn~ () 42",
			function{effect: true, statement: 42},
		))
		t.Run("no params, block", testEvalStatement(
			"fn~ () { return 42 }",
			function{
				effect:    true,
				statement: statementList{statements: []interface{}{ret{value: 42}}},
			},
		))
		t.Run("void", testEvalStatement(
			"fn~ (foo, bar, baz) {;}",
			function{
				statement: statementList{},
				effect:    true,
				params:    []string{"foo", "bar", "baz"},
			},
		))
		t.Run("expression", testEvalStatement(
			"fn~ (foo, bar, baz) 42",
			function{
				effect:    true,
				params:    []string{"foo", "bar", "baz"},
				statement: 42,
			},
		))
		t.Run("block", testEvalStatement(
			"fn~ (foo, bar, baz) { return 42 }",
			function{
				effect:    true,
				params:    []string{"foo", "bar", "baz"},
				statement: statementList{statements: []interface{}{ret{value: 42}}},
			},
		))
		t.Run("collect, void", testEvalStatement(
			"fn~ (foo, bar, baz, ...qux) {;}",
			function{
				statement:    statementList{},
				effect:       true,
				params:       []string{"foo", "bar", "baz"},
				collectParam: "qux",
			},
		))
		t.Run("collect, expression", testEvalStatement(
			"fn~ (foo, bar, baz, ...qux) 42",
			function{
				effect:       true,
				params:       []string{"foo", "bar", "baz"},
				collectParam: "qux",
				statement:    42,
			},
		))
		t.Run("collect, block", testEvalStatement(
			"fn~ (foo, bar, baz, ...qux) { return 42 }",
			function{
				effect:       true,
				params:       []string{"foo", "bar", "baz"},
				collectParam: "qux",
				statement:    statementList{statements: []interface{}{ret{value: 42}}},
			},
		))
		t.Run("struct", testEvalStatement(
			"fn~ () {}",
			function{
				effect:    true,
				statement: structure{},
			},
		))
	})

	t.Run("indexer expression", func(t *testing.T) {
		t.Run("expression indexer", func(t *testing.T) {
			t.Run("list", func(t *testing.T) {
				t.Run("index", testEvalStatement("[1, 2, 3][1]", 2))
				t.Run("range", func(t *testing.T) {
					t.Run("full", testEvalStatement(
						"[1, 2, 3, 4, 5][2:4]",
						list{values: []interface{}{3, 4}},
					))
					t.Run("from", testEvalStatement(
						"[1, 2, 3, 4, 5][2:]",
						list{values: []interface{}{3, 4, 5}},
					))
					t.Run("to", testEvalStatement(
						"[1, 2, 3, 4, 5][:4]",
						list{values: []interface{}{1, 2, 3, 4}},
					))
					t.Run("reslice", testEvalStatement(
						"[1, 2, 3, 4, 5][:]",
						list{values: []interface{}{1, 2, 3, 4, 5}},
					))
				})
			})
			t.Run("struct", testEvalStatement(`{foo: 1, bar: 2, baz: 3}["bar"]`, 2))
		})
		t.Run("symbol indexer", testEvalStatement("{foo: 1, bar: 2, baz: 3}.bar", 2))
	})

	t.Run("function application", func(t *testing.T) {
		t.Run("void", testEvalStatement("fn () {;}()", nil))
		t.Run("partial", testEvalStatement(
			"fn (a, b, c) {;}(42)",
			function{
				params:    []string{"a", "b", "c"},
				args:      []interface{}{42},
				statement: statementList{},
			},
		))
		t.Run("simple value", testEvalStatement("(fn () 42)()", 42))
		t.Run("argument", testEvalStatement("(fn (a, b, c) a)(1, 2, 3)", 1))
		t.Run("collect", testEvalStatement(
			"(fn (a, b, ...c) c)(1, 2, 3, 4, 5)",
			list{values: []interface{}{3, 4, 5}},
		))
		t.Run("spread", testEvalStatement(
			"(fn (a, b, c) c)([1, 2, 3]...)",
			3,
		))
		t.Run("collect and spread", testEvalStatement(
			"(fn (a, b, ...c) c)(1, [2, 3, 4]..., 5)",
			list{values: []interface{}{3, 4, 5}},
		))
		t.Run("statement list", func(t *testing.T) {
			t.Run("empty, void", testEvalStatement("fn () {;}()", nil))
			t.Run("partial", testEvalStatement(
				"fn (a, b, c) { a; b; return c }(42)",
				function{
					params: []string{"a", "b", "c"},
					args:   []interface{}{42},
					statement: statementList{statements: []interface{}{
						symbol{name: "a"},
						symbol{name: "b"},
						ret{value: symbol{name: "c"}},
					}},
				},
			))
			t.Run("void", testEvalStatement("fn (a, b, c) { a; b; c }(1, 2, 3)", nil))
			t.Run("simple value", testEvalStatement("fn (a, b, c) { a; b; c; return 42 }(1, 2, 3)", 42))
			t.Run("argument", testEvalStatement("fn (a, b, c) { a; b; return c }(1, 2, 3)", 3))
			t.Run("collect", testEvalStatement(
				"(fn (a, b, ...c) { a; b; return c })(1, 2, 3, 4, 5)",
				list{values: []interface{}{3, 4, 5}},
			))
			t.Run("spread", testEvalStatement(
				"(fn (a, b, c) { a; b; return c })([1, 2, 3]...)",
				3,
			))
			t.Run("collect and spread", testEvalStatement(
				"(fn (a, b, ...c) { a; b; return c })(1, [2, 3, 4]..., 5)",
				list{values: []interface{}{3, 4, 5}},
			))
		})
	})

	t.Run("operator", func(t *testing.T) {
		t.Run("unary", testEvalStatement("^2", -3))
		t.Run("binary", testEvalStatement("3 & 7", 3))
		t.Run("binary, chained", testEvalStatement("3 & 7 & 5", 1))
		t.Run("precedence", testEvalStatement("3 * 4 - 2 * 5", 2))
		t.Run("grouping", testEvalStatement("3 * (4 - 2) * 5", 30))
	})

	t.Run("function chaining", testEvalStatement(
		`[1, 2, 3]
		 -> fn (l) l
		 -> fn (l) (fn (a, b, c) c)(l...)
		 -> fn (x) x`,
		3,
	))

	t.Run("ternay expression", func(t *testing.T) {
		t.Run("true", testEvalStatement("true ? 42 : 36", 42))
		t.Run("false", testEvalStatement("false ? 42 : 36", 36))
	})

	t.Run("if", func(t *testing.T) {
		t.Run("simple", testEvalStatement(
			`fn () {
				if true {
					return 42
				}

				return 36
			}()`,
			42,
		))
		t.Run("false", testEvalStatement(
			`fn () {
				if false {
					return 42
				}

				return 36
			}()`,
			36,
		))
		t.Run("else", testEvalStatement(
			`fn () {
				if false {
					return 42
				} else {
					return 36
				}
			}()`,
			36,
		))
		t.Run("else if", testEvalStatement(
			`fn () {
				if false {
					return 42
				} else if true {
					return 36
				} else {
					return 24
				}
			}()`,
			36,
		))
		t.Run("else if, false", testEvalStatement(
			`fn () {
				if false {
					return 42
				} else if false {
					return 36
				}

				return 24
			}()`,
			24,
		))
	})

	t.Run("switch", func(t *testing.T) {
		t.Run("empty", testEvalStatement(
			"fn () { switch { default: }; return 42 }()",
			42,
		))
		t.Run("default only", testEvalStatement(
			"fn () { switch { default: return 42 }; return 36 }()",
			42,
		))
		t.Run("cases only", testEvalStatement(
			"fn () { switch { case false: return 42; case true: return 36; } }()",
			36,
		))
		t.Run("cases and default", testEvalStatement(
			`fn () {
				switch {
				case false: return 42
				case true: return 36
				default: return 24
				}
			}()`,
			36,
		))
		t.Run("cases and default, choose default", testEvalStatement(
			`fn () {
				switch {
				case false: return 42
				case false: return 36
				default: return 24
				}
			}()`,
			24,
		))
		t.Run("expression", func(t *testing.T) {
			t.Run("empty", testEvalStatement(
				"fn () { switch 1 { default: }; return 42 }()",
				42,
			))
			t.Run("default only", testEvalStatement(
				"fn () { switch 1 { default: return 42 }; return 36 }()",
				42,
			))
			t.Run("cases only", testEvalStatement(
				"fn () { switch 1 { case 1: return 42; case 2: return 36; } }()",
				42,
			))
			t.Run("cases and default, choose default", testEvalStatement(
				`fn () {
					switch 0 {
					case 1: return 42
					case 2: return 36
					default: return 24
					}
				}()`,
				24,
			))
		})
	})

	t.Run("loop", func(t *testing.T) {
		t.Run("simple", testEvalStatement("fn () { for { return 42 } }()", 42))
		t.Run("true", testEvalStatement("fn () { for true { return 42 }; return 36 }()", 42))
		t.Run("false", testEvalStatement("fn () { for false { return 42 }; return 36 }()", 36))
		t.Run("range", func(t *testing.T) {
			t.Run("range only", func(t *testing.T) {
				t.Run("closed", testEvalStatement(
					`fn () {
						for 0:3 {
							return 42
						}

						return 36
					}()`,
					42,
				))
				t.Run("from 0", testEvalStatement(
					`fn () {
						for : {
							return 42
						}

						return 36
					}()`,
					42,
				))
				t.Run("from", testEvalStatement(
					`fn () {
						for 1: {
							return 42
						}

						return 36
					}()`,
					42,
				))
				t.Run("to", testEvalStatement(
					`fn () {
						for :2 {
							return 42
						}

						return 36
					}()`,
					42,
				))
			})
			t.Run("named", func(t *testing.T) {
				t.Run("closed", testEvalStatement(
					`fn () {
						for a in 0:3 {
							if a == 2 {
								return a
							}
						}

						return 42
					}()`,
					2,
				))
				t.Run("from 0", testEvalStatement(
					`fn () {
						for a in : {
							if a == 2 {
								return a
							}
						}

						return 42
					}()`,
					2,
				))
				t.Run("from", testEvalStatement(
					`fn () {
						for a in 1: {
							if a == 2 {
								return a
							}
						}

						return 42
					}()`,
					2,
				))
				t.Run("to", testEvalStatement(
					`fn () {
						for a in :3 {
							if a == 2 {
								return a
							}
						}

						return 42
					}()`,
					2,
				))
			})
		})
		t.Run("range list", testEvalStatement(
			`fn () {
				for a in [0, 1, 2] {
					if a == 2 {
						return a
					}
				}

				return 42
			}()`,
			2,
		))
		t.Run("range struct", testEvalStatement(
			`fn () {
				for a in {foo: 0, bar: 1, baz: 2} {
					if a == 2 {
						return a
					}
				}

				return 42
			}()`,
			2,
		))
		t.Run("break", testEvalStatement(
			`fn () {
				for a in [0, 1, 2] {
					if a > 0 {
						break
					}

					if a > 1 {
						return a
					}
				}

				return 42
			}()`,
			42,
		))
		t.Run("continue", testEvalStatement(
			`fn () {
				for a in [0, 1, 2] {
					if a > 0 {
						continue
					}

					if a > 1 {
						return a
					}
				}

				return 42
			}()`,
			42,
		))
	})

	t.Run("define", func(t *testing.T) {
		t.Run("value", func(t *testing.T) {
			t.Run("immutable", testEvalStatement("fn () { let a 42; return a }()", 42))
			t.Run("mutable", testEvalStatement("fn () { let a 42; return a }()", 42))
		})
		t.Run("value group", func(t *testing.T) {
			t.Run("immutable", testEvalStatement(
				`fn () {
					let (
						a 42
						b 36
					)

					return b - a
				}()`,
				-6,
			))
			t.Run("mutable", testEvalStatement(
				`fn () {
					let ~ (
						a 42
						b 36
					)

					return b - a
				}()`,
				-6,
			))
			t.Run("mixed", testEvalStatement(
				`fn () {
					let (
						  a 42
						~ b 36
					)

					return b - a
				}()`,
				-6,
			))
		})
		t.Run("function", func(t *testing.T) {
			t.Run("pure", testEvalStatement(
				"fn () { fn square(x) x * x; return square }()",
				function{
					params: []string{"x"},
					statement: binary{
						op:    mul,
						left:  symbol{name: "x"},
						right: symbol{name: "x"},
					},
				},
			))
			t.Run("mutable", testEvalStatement(
				"fn () { fn ~ printSquare(x) print(x * x); return printSquare }()",
				function{
					effect: true,
					params: []string{"x"},
					statement: functionApplication{
						function: symbol{name: "print"},
						args: []interface{}{
							binary{
								op:    mul,
								left:  symbol{name: "x"},
								right: symbol{name: "x"},
							},
						},
					},
				},
			))
		})
		t.Run("function group", func(t *testing.T) {
			t.Run("pure", testEvalStatement(
				`fn () {
					fn (
						a () 42,
						b (x) 36,
						c (x) 24,
					)
					
					return c
				}()`,
				function{
					params:    []string{"x"},
					statement: 24,
				},
			))
			t.Run("effect", testEvalStatement(
				`fn () {
					fn ~ (
						a () 42,
						b (x) 36,
						c (x) 24,
					)
					
					return c
				}()`,
				function{
					effect:    true,
					params:    []string{"x"},
					statement: 24,
				},
			))
			t.Run("mixed", testEvalStatement(
				`fn () {
					fn (
						  a () 42,
						~ b (x) 36,
						~ c (x) 24,
					)
					
					return c
				}()`,
				function{
					effect:    true,
					params:    []string{"x"},
					statement: 24,
				},
			))
		})
	})

	t.Run("assign", func(t *testing.T) {
		t.Run("symbol", testEvalStatement(
			"fn () { let ~ a 42; set a 36; return a }()",
			36,
		))
		t.Run("indexer", testEvalStatement(
			"fn () { let a ~[1, 2, 3]; set a[1] 42; return a[1] }()",
			42,
		))
		t.Run("indexer, eq", testEvalStatement(
			"fn () { let a ~[1, 2, 3]; a[1] = 42; return a[1] }()",
			42,
		))
		t.Run("struct key", testEvalStatement(
			`fn () { let a ~{foo: 1, bar: 2, baz: 3}; set a["foo"] 42; return a.foo }()`,
			42,
		))
		t.Run("symbol indexer", testEvalStatement(
			"fn () { let a ~{foo: 1, bar: 2, baz: 3}; set a.foo 42; return a.foo }()",
			42,
		))
		t.Run("group", testEvalStatement(
			`fn () {
				let (
					~ a 42
					  b "foo"
				)

				set (
					a    ~{foo: 42, b: 36}
					a[b] 36
					a.b  24
				)

				return [a.foo, a.b]
			}()`,
			list{values: []interface{}{36, 24}},
		))
	})
}
