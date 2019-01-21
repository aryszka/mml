// Generated code
package main

import "github.com/aryszka/mml"

var _args interface{} = mml.Args
var _bool interface{} = mml.Bool
var _close interface{} = mml.Close
var _error interface{} = mml.Error
var _exit interface{} = mml.Exit
var _float interface{} = mml.Float
var _format interface{} = mml.Format
var _has interface{} = mml.Has
var _int interface{} = mml.Int
var _isBool interface{} = mml.IsBool
var _isChannel interface{} = mml.IsChannel
var _isError interface{} = mml.IsError
var _isFloat interface{} = mml.IsFloat
var _isFunction interface{} = mml.IsFunction
var _isInt interface{} = mml.IsInt
var _isList interface{} = mml.IsList
var _isString interface{} = mml.IsString
var _isStruct interface{} = mml.IsStruct
var _keys interface{} = mml.Keys
var _len interface{} = mml.Len
var _open interface{} = mml.Open
var _panic interface{} = mml.Panic
var _parseAST interface{} = mml.ParseAST
var _parseFloat interface{} = mml.ParseFloat
var _parseInt interface{} = mml.ParseInt
var _stderr interface{} = mml.Stderr
var _stdin interface{} = mml.Stdin
var _stdout interface{} = mml.Stdout
var _string interface{} = mml.String

func init() {
	var modulePath string
	modulePath = "main.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _printValidationErrors interface{}
		var _validateDefinitions interface{}
		var _compileModuleCode interface{}
		var _compileModules interface{}
		var _modules interface{}
		var _validation interface{}
		var _builtins interface{}
		var _code interface{}
		var _parse interface{}
		var _definitions interface{}
		var _snippets interface{}
		var _compile interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _onlyErr interface{}
		var _passErr interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _is interface{}
		mml.Nop(_printValidationErrors, _validateDefinitions, _compileModuleCode, _compileModules, _modules, _validation, _builtins, _code, _parse, _definitions, _snippets, _compile, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _bind, _identity, _eq, _any, _function, _channel, _type, _listOf, _structOf, _range, _or, _and, _predicate, _is)
		var __lang = mml.Modules.Use("lang.mml")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_onlyErr = __lang.Values["onlyErr"]
		_passErr = __lang.Values["passErr"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_predicate = __lang.Values["predicate"]
		_is = __lang.Values["is"]
		_code = mml.Modules.Use("code.mml")
		_parse = mml.Modules.Use("parse.mml")
		_definitions = mml.Modules.Use("definitions.mml")
		_snippets = mml.Modules.Use("snippets.mml")
		_compile = mml.Modules.Use("compile.mml")
		_printValidationErrors = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0]
				var _errors = a[1]

				mml.Nop(_m, _errors)

				mml.Nop()
				_log.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s:", mml.Ref(_m, "path"))}).Values))}).Values)
				for _, _e := range _errors.(*mml.List).Values {

					mml.Nop()
					_log.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _e)}).Values)
				}
				return nil
			},
			FixedArgs: 2,
		}
		_validateDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _modules = a[0]

				mml.Nop(_modules)
				var _hasErrors interface{}
				mml.Nop(_hasErrors)
				_hasErrors = false
				for _, _m := range _modules.(*mml.List).Values {
					var _errors interface{}
					mml.Nop(_errors)
					_errors = mml.Ref(_definitions, "validate").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m)}).Values)
					c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _errors)}).Values), 0)
					if c.(bool) {
						mml.Nop()
						_hasErrors = true
						_printValidationErrors.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m, _errors)}).Values)
					}
				}
				c = _hasErrors
				if c.(bool) {
					mml.Nop()
					return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "undefined reference(s) found")}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_compileModuleCode = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _moduleCode = a[0]

				mml.Nop(_moduleCode)

				mml.Nop()
				_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "modulePath = \"%s\"", mml.Ref(_moduleCode, "path"))}).Values))}).Values)
				_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "moduleHead"))}).Values)
				_onlyErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _log)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _passErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _stdout)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_compile, "do").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _moduleCode)}).Values))}).Values))}).Values)
				_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "moduleFooter"))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_compileModules = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0]

				mml.Nop(_m)

				mml.Nop()
				for _, _mi := range _m.(*mml.List).Values {

					mml.Nop()
					_compileModuleCode.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mi)}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		c = mml.BinaryOp(13, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _args)}).Values), 2)
		if c.(bool) {
			mml.Nop()
			_stderr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "usage: mml source_file")}).Values)
			_exit.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values)
		}
		_modules = mml.Ref(_parse, "modules").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_args, 1))}).Values)
		c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _modules)}).Values)
		if c.(bool) {
			mml.Nop()
			_panic.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _modules)}).Values)
		}
		_validation = _validateDefinitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _modules)}).Values)
		c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _validation)}).Values)
		if c.(bool) {
			mml.Nop()
			_panic.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _validation)}).Values)
		}
		_builtins = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _k = a[0]

				mml.Nop(_k)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var _%s interface{} = mml.%s", _k, mml.Ref(mml.Ref(_code, "builtin"), _k))}).Values)
			},
			FixedArgs: 1,
		})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _sort.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _left = a[0]
				var _right = a[1]

				mml.Nop(_left, _right)
				return mml.BinaryOp(13, _left, _right)
			},
			FixedArgs: 2,
		})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "builtin"))}).Values))}).Values))}).Values))}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "head"))}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _builtins)}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "initHead"))}).Values)
		_compileModules.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _modules)}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "initFooter"))}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "mainHead"))}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_args, 1))}).Values)
		_stdout.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_snippets, "mainFooter"))}).Values)
		return exports
	})
	modulePath = "lang.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _onlyErr interface{}
		var _passErr interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _is interface{}
		var _logger interface{}
		var _list interface{}
		var _strings interface{}
		var _ints interface{}
		var _errors interface{}
		var _functions interface{}
		var _match interface{}
		mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _bind, _identity, _eq, _any, _function, _channel, _type, _listOf, _structOf, _range, _or, _and, _predicate, _is, _logger, _list, _strings, _ints, _errors, _functions, _match)
		_list = mml.Modules.Use("list.mml")
		_strings = mml.Modules.Use("strings.mml")
		_ints = mml.Modules.Use("ints.mml")
		_logger = mml.Modules.Use("log.mml")
		_errors = mml.Modules.Use("errors.mml")
		_functions = mml.Modules.Use("functions.mml")
		_match = mml.Modules.Use("match.mml")
		_fold = mml.Ref(_list, "fold")
		exports["fold"] = _fold
		_foldr = mml.Ref(_list, "foldr")
		exports["foldr"] = _foldr
		_map = mml.Ref(_list, "map")
		exports["map"] = _map
		_filter = mml.Ref(_list, "filter")
		exports["filter"] = _filter
		_contains = mml.Ref(_list, "contains")
		exports["contains"] = _contains
		_sort = mml.Ref(_list, "sort")
		exports["sort"] = _sort
		_flat = mml.Ref(_list, "flat")
		exports["flat"] = _flat
		_uniq = mml.Ref(_list, "uniq")
		exports["uniq"] = _uniq
		_every = mml.Ref(_list, "every")
		exports["every"] = _every
		_some = mml.Ref(_list, "some")
		exports["some"] = _some
		_join = mml.Ref(_strings, "join")
		exports["join"] = _join
		_joins = mml.Ref(_strings, "joins")
		exports["joins"] = _joins
		_formats = mml.Ref(_strings, "formats")
		exports["formats"] = _formats
		_enum = mml.Ref(_ints, "enum")
		exports["enum"] = _enum
		_log = mml.Ref(_logger, "println")
		exports["log"] = _log
		_onlyErr = mml.Ref(_errors, "only")
		exports["onlyErr"] = _onlyErr
		_passErr = mml.Ref(_errors, "pass")
		exports["passErr"] = _passErr
		_bind = mml.Ref(_functions, "bind")
		exports["bind"] = _bind
		_identity = mml.Ref(_functions, "identity")
		exports["identity"] = _identity
		_eq = mml.Ref(_functions, "eq")
		exports["eq"] = _eq
		_any = mml.Ref(_match, "any")
		exports["any"] = _any
		_function = mml.Ref(_match, "function")
		exports["function"] = _function
		_channel = mml.Ref(_match, "channel")
		exports["channel"] = _channel
		_type = mml.Ref(_match, "type")
		exports["type"] = _type
		_listOf = mml.Ref(_match, "listOf")
		exports["listOf"] = _listOf
		_structOf = mml.Ref(_match, "structOf")
		exports["structOf"] = _structOf
		_range = mml.Ref(_match, "range")
		exports["range"] = _range
		_or = mml.Ref(_match, "or")
		exports["or"] = _or
		_and = mml.Ref(_match, "and")
		exports["and"] = _and
		_predicate = mml.Ref(_match, "predicate")
		exports["predicate"] = _predicate
		_is = mml.Ref(_match, "is")
		exports["is"] = _is
		return exports
	})
	modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _first interface{}
		var _contains interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _sort interface{}
		mml.Nop(_fold, _foldr, _map, _filter, _first, _contains, _flat, _uniq, _every, _some, _sort)
		_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _i = a[1]
				var _l = a[2]

				mml.Nop(_f, _i, _l)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values), 0)
					if c.(bool) {
						return _i
					} else {
						return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, 0), _i)}).Values), mml.RefRange(_l, 1, nil))}).Values)
					}
				}()
			},
			FixedArgs: 3,
		}
		exports["fold"] = _fold
		_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _i = a[1]
				var _l = a[2]

				mml.Nop(_f, _i, _l)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values), 0)
					if c.(bool) {
						return _i
					} else {
						return _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 3,
		}
		exports["foldr"] = _foldr
		_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0]
				var _l = a[1]

				mml.Nop(_m, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]

						mml.Nop(_c, _r)
						return &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _m.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values))}
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["map"] = _map
		_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]
				var _l = a[1]

				mml.Nop(_p, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]

						mml.Nop(_c, _r)
						return func() interface{} {
							c = _p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
							if c.(bool) {
								return &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _c)}
							} else {
								return _r
							}
						}()
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["filter"] = _filter
		_first = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]
				var _l = a[1]

				mml.Nop(_p, _l)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values), 0)
					if c.(bool) {
						return &mml.List{Values: []interface{}{}}
					} else {
						return func() interface{} {
							c = _p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, 0))}).Values)
							if c.(bool) {
								return _l
							} else {
								return _first.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _p, mml.RefRange(_l, 1, nil))}).Values)
							}
						}()
					}
				}()
			},
			FixedArgs: 2,
		}
		exports["first"] = _first
		_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0]
				var _l = a[1]

				mml.Nop(_i, _l)
				return mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _first.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ii = a[0]

						mml.Nop(_ii)
						return mml.BinaryOp(11, _ii, _i)
					},
					FixedArgs: 1,
				}, _l)}).Values))}).Values), 0)
			},
			FixedArgs: 2,
		}
		exports["contains"] = _contains
		_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _result = a[1]

						mml.Nop(_c, _result)
						return &mml.List{Values: append(append([]interface{}{}, _result.(*mml.List).Values...), _c.(*mml.List).Values...)}
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 1,
		}
		exports["flat"] = _flat
		_uniq = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _eq = a[0]
				var _l = a[1]

				mml.Nop(_eq, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _u = a[1]

						mml.Nop(_c, _u)
						return func() interface{} {
							c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
								F: func(a []interface{}) interface{} {
									var c interface{}
									mml.Nop(c)
									var _i = a[0]

									mml.Nop(_i)
									return _eq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _i, _c)}).Values)
								},
								FixedArgs: 1,
							}, _u)}).Values))}).Values), 0)
							if c.(bool) {
								return &mml.List{Values: append(append([]interface{}{}, _u.(*mml.List).Values...), _c)}
							} else {
								return _u
							}
						}()
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["uniq"] = _uniq
		_every = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]
				var _l = a[1]

				mml.Nop(_p, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _r = a[1]

						mml.Nop(_i, _r)
						return (_r.(bool) && _p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _i)}).Values).(bool))
					},
					FixedArgs: 2,
				}, true, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["every"] = _every
		_some = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]
				var _l = a[1]

				mml.Nop(_p, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _r = a[1]

						mml.Nop(_i, _r)
						return (_r.(bool) || _p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _i)}).Values).(bool))
					},
					FixedArgs: 2,
				}, false, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["some"] = _some
		_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0]
				var _l = a[1]

				mml.Nop(_less, _l)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values), 0)
					if c.(bool) {
						return &mml.List{Values: []interface{}{}}
					} else {
						return &mml.List{Values: append(append(append([]interface{}{}, _sort.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _less)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _i = a[0]

								mml.Nop(_i)
								return !_less.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, 0), _i)}).Values).(bool)
							},
							FixedArgs: 1,
						})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_l, 1, nil))}).Values))}).Values).(*mml.List).Values...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _less)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _less.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, 0))}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_l, 1, nil))}).Values))}).Values).(*mml.List).Values...)}
					}
				}()
			},
			FixedArgs: 2,
		}
		exports["sort"] = _sort
		return exports
	})
	modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _firstOr interface{}
		var _join interface{}
		var _joins interface{}
		var _joinTwo interface{}
		var _formats interface{}
		var _formatOne interface{}
		var _escape interface{}
		var _unescape interface{}
		mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape)
		_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0]
				var _l = a[1]

				mml.Nop(_v, _l)
				return func() interface{} {
					c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values), 0)
					if c.(bool) {
						return mml.Ref(_l, 0)
					} else {
						return _v
					}
				}()
			},
			FixedArgs: 2,
		}
		_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _s = a[1]

				mml.Nop(_j, _s)
				return func() interface{} {
					c = mml.BinaryOp(13, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values), 2)
					if c.(bool) {
						return _firstOr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "", _s)}).Values)
					} else {
						return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _j, mml.RefRange(_s, 1, nil))}).Values))
					}
				}()
			},
			FixedArgs: 2,
		}
		exports["join"] = _join
		_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _s interface{}
				_s = &mml.List{a[1:]}

				mml.Nop(_j, _s)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _j, _s)}).Values)
			},
			FixedArgs: 1,
		}
		exports["joins"] = _joins
		_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _left = a[1]
				var _right = a[2]

				mml.Nop(_j, _left, _right)
				return _joins.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _j, _left, _right)}).Values)
			},
			FixedArgs: 3,
		}
		exports["joinTwo"] = _joinTwo
		_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a interface{}
				_a = &mml.List{a[1:]}

				mml.Nop(_f, _a)
				return _format.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _a)}).Values)
			},
			FixedArgs: 1,
		}
		exports["formats"] = _formats
		_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a = a[1]

				mml.Nop(_f, _a)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _a)}).Values)
			},
			FixedArgs: 2,
		}
		exports["formatOne"] = _formatOne
		_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				var _first interface{}
				mml.Nop(_first)
				c = mml.BinaryOp(11, _s, "")
				if c.(bool) {
					mml.Nop()
					return ""
				}
				_first = mml.Ref(_s, 0)
				switch _first {
				case "\b":

					mml.Nop()
					_first = "\\b"
				case "\f":

					mml.Nop()
					_first = "\\f"
				case "\n":

					mml.Nop()
					_first = "\\n"
				case "\r":

					mml.Nop()
					_first = "\\r"
				case "\t":

					mml.Nop()
					_first = "\\t"
				case "\v":

					mml.Nop()
					_first = "\\v"
				case "\"":

					mml.Nop()
					_first = "\\\""
				case "\\":

					mml.Nop()
					_first = "\\\\"
				}
				return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_s, 1, nil))}).Values))
				return nil
			},
			FixedArgs: 1,
		}
		exports["escape"] = _escape
		_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				var _esc interface{}
				var _r interface{}
				mml.Nop(_esc, _r)
				_esc = false
				_r = &mml.List{Values: []interface{}{}}
				for _i := 0; _i < _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values).(int); _i++ {
					var _c interface{}
					mml.Nop(_c)
					_c = mml.Ref(_s, _i)
					c = _esc
					if c.(bool) {
						mml.Nop()
						switch _c {
						case "b":

							mml.Nop()
							_c = "\b"
						case "f":

							mml.Nop()
							_c = "\f"
						case "n":

							mml.Nop()
							_c = "\n"
						case "r":

							mml.Nop()
							_c = "\r"
						case "t":

							mml.Nop()
							_c = "\t"
						case "v":

							mml.Nop()
							_c = "\v"
						}
						_r = &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _c)}
						_esc = false
						continue
					}
					c = mml.BinaryOp(11, _c, "\\")
					if c.(bool) {
						mml.Nop()
						_esc = true
						continue
					}
					_r = &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _c)}
				}
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "", _r)}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		exports["unescape"] = _unescape
		return exports
	})
	modulePath = "ints.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _counter interface{}
		var _enum interface{}
		var _max interface{}
		var _min interface{}
		mml.Nop(_counter, _enum, _max, _min)
		_counter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				var _c interface{}
				mml.Nop(_c)
				_c = mml.UnaryOp(2, 1)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()

						mml.Nop()
						_c = mml.BinaryOp(9, _c, 1)
						return _c
						return nil
					},
					FixedArgs: 0,
				}
				return nil
			},
			FixedArgs: 0,
		}
		exports["counter"] = _counter
		_enum = _counter
		exports["enum"] = _enum
		_max = 9000
		exports["max"] = _max
		_min = mml.UnaryOp(2, 9000)
		exports["min"] = _min
		return exports
	})
	modulePath = "log.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _println interface{}
		var _list interface{}
		var _strings interface{}
		mml.Nop(_println, _list, _strings)
		_list = mml.Modules.Use("list.mml")
		_strings = mml.Modules.Use("strings.mml")
		_println = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a interface{}
				_a = &mml.List{a[0:]}

				mml.Nop(_a)

				mml.Nop()
				_stderr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, " ")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_list, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _string)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values))}).Values))}).Values)
				_stderr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\n")}).Values)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values), 0)
					if c.(bool) {
						return ""
					} else {
						return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values), 1))
					}
				}()
				return nil
			},
			FixedArgs: 0,
		}
		exports["println"] = _println
		return exports
	})
	modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _ifErr interface{}
		var _not interface{}
		var _yes interface{}
		var _pass interface{}
		var _only interface{}
		var _any interface{}
		var _list interface{}
		mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list)
		_list = mml.Modules.Use("list.mml")
		_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0]
				var _f = a[1]

				mml.Nop(_mod, _f)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _a = a[0]

						mml.Nop(_a)
						return func() interface{} {
							c = _mod.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values))}).Values)
							if c.(bool) {
								return _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values)
							} else {
								return _a
							}
						}()
					},
					FixedArgs: 1,
				}
			},
			FixedArgs: 2,
		}
		_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0]

				mml.Nop(_x)
				return !_x.(bool)
			},
			FixedArgs: 1,
		}
		_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0]

				mml.Nop(_x)
				return _x
			},
			FixedArgs: 1,
		}
		_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				return _ifErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _not, _f)}).Values)
			},
			FixedArgs: 1,
		}
		exports["pass"] = _pass
		_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				return _ifErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _yes, _f)}).Values)
			},
			FixedArgs: 1,
		}
		exports["only"] = _only
		_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return mml.Ref(_list, "fold").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]

						mml.Nop(_c, _r)
						return func() interface{} {
							c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r)}).Values)
							if c.(bool) {
								return _r
							} else {
								return func() interface{} {
									c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
									if c.(bool) {
										return _c
									} else {
										return &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _c)}
									}
								}()
							}
						}()
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 1,
		}
		exports["any"] = _any
		return exports
	})
	modulePath = "functions.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		mml.Nop(_bind, _identity, _eq)
		_bind = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _args interface{}
				_args = &mml.List{a[1:]}

				mml.Nop(_f, _args)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _nextArgs interface{}
						_nextArgs = &mml.List{a[0:]}

						mml.Nop(_nextArgs)
						var _a interface{}
						mml.Nop(_a)
						_a = &mml.List{Values: append(append([]interface{}{}, _args.(*mml.List).Values...), _nextArgs.(*mml.List).Values...)}
						return _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a.(*mml.List).Values...)}).Values)
						return nil
					},
					FixedArgs: 0,
				}
			},
			FixedArgs: 1,
		}
		exports["bind"] = _bind
		_identity = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0]

				mml.Nop(_x)
				return _x
			},
			FixedArgs: 1,
		}
		exports["identity"] = _identity
		_eq = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x interface{}
				_x = &mml.List{a[0:]}

				mml.Nop(_x)
				return (mml.BinaryOp(13, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _x)}).Values), 2).(bool) || (mml.BinaryOp(11, mml.Ref(_x, 0), mml.Ref(_x, 1)).(bool) && _eq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_x, 1, nil).(*mml.List).Values...)}).Values).(bool)))
			},
			FixedArgs: 0,
		}
		exports["eq"] = _eq
		return exports
	})
	modulePath = "match.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _token interface{}
		var _none interface{}
		var _integer interface{}
		var _floating interface{}
		var _stringType interface{}
		var _boolean interface{}
		var _errorType interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _type interface{}
		var _defineType interface{}
		var _intRangeType interface{}
		var _floatRangeType interface{}
		var _isRange interface{}
		var _isNaturalRange interface{}
		var _defineRange interface{}
		var _intRange interface{}
		var _floatRange interface{}
		var _stringRangeType interface{}
		var _stringRange interface{}
		var _listType interface{}
		var _listRange interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _unionType interface{}
		var _intersectType interface{}
		var _predicateType interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _isSimpleType interface{}
		var _isComplexType interface{}
		var _isType interface{}
		var _typesMatch interface{}
		var _primitives interface{}
		var _matchPrimitive interface{}
		var _matchInt interface{}
		var _matchFloat interface{}
		var _matchString interface{}
		var _matchToList interface{}
		var _matchToListType interface{}
		var _matchList interface{}
		var _matchStruct interface{}
		var _matchUnion interface{}
		var _matchIntersection interface{}
		var _matchOne interface{}
		var _is interface{}
		var _ints interface{}
		var _floats interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _first interface{}
		var _contains interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _sort interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		mml.Nop(_token, _none, _integer, _floating, _stringType, _boolean, _errorType, _any, _function, _channel, _type, _defineType, _intRangeType, _floatRangeType, _isRange, _isNaturalRange, _defineRange, _intRange, _floatRange, _stringRangeType, _stringRange, _listType, _listRange, _listOf, _structOf, _range, _unionType, _intersectType, _predicateType, _or, _and, _predicate, _isSimpleType, _isComplexType, _isType, _typesMatch, _primitives, _matchPrimitive, _matchInt, _matchFloat, _matchString, _matchToList, _matchToListType, _matchList, _matchStruct, _matchUnion, _matchIntersection, _matchOne, _is, _ints, _floats, _fold, _foldr, _map, _filter, _first, _contains, _flat, _uniq, _every, _some, _sort, _bind, _identity, _eq)
		var __list = mml.Modules.Use("list.mml")
		_fold = __list.Values["fold"]
		_foldr = __list.Values["foldr"]
		_map = __list.Values["map"]
		_filter = __list.Values["filter"]
		_first = __list.Values["first"]
		_contains = __list.Values["contains"]
		_flat = __list.Values["flat"]
		_uniq = __list.Values["uniq"]
		_every = __list.Values["every"]
		_some = __list.Values["some"]
		_sort = __list.Values["sort"]
		var __functions = mml.Modules.Use("functions.mml")
		_bind = __functions.Values["bind"]
		_identity = __functions.Values["identity"]
		_eq = __functions.Values["eq"]
		_ints = mml.Modules.Use("ints.mml")
		_floats = mml.Modules.Use("floats.mml")
		_token = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _token
			},
			FixedArgs: 0,
		}
		_none = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _none
			},
			FixedArgs: 0,
		}
		_integer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _integer
			},
			FixedArgs: 0,
		}
		_floating = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _floating
			},
			FixedArgs: 0,
		}
		_stringType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _stringType
			},
			FixedArgs: 0,
		}
		_boolean = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _boolean
			},
			FixedArgs: 0,
		}
		_errorType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _errorType
			},
			FixedArgs: 0,
		}
		_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _any
			},
			FixedArgs: 0,
		}
		exports["any"] = _any
		_function = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _function
			},
			FixedArgs: 0,
		}
		exports["function"] = _function
		_channel = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return _channel
			},
			FixedArgs: 0,
		}
		exports["channel"] = _channel
		_type = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _t = a[0]

				mml.Nop(_t)

				mml.Nop()
				switch _t {
				case _int:

					mml.Nop()
					return _integer
				case _float:

					mml.Nop()
					return _floating
				case _string:

					mml.Nop()
					return _stringType
				case _bool:

					mml.Nop()
					return _boolean
				case _error:

					mml.Nop()
					return _errorType
				default:

					mml.Nop()
					return _t
				}
				return nil
			},
			FixedArgs: 1,
		}
		exports["type"] = _type
		_defineType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _name = a[0]

				mml.Nop(_name)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["token"] = _token
					s.Values["type"] = _name
					return s
				}()
			},
			FixedArgs: 1,
		}
		_intRangeType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "int-range")}).Values)
		_floatRangeType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "float-range")}).Values)
		_isRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ofType = a[0]
				var _min = a[1]
				var _max = a[2]

				mml.Nop(_ofType, _min, _max)
				return ((_ofType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min)}).Values).(bool) && _ofType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _max)}).Values).(bool)) && mml.BinaryOp(14, _min, _max).(bool))
			},
			FixedArgs: 3,
		}
		_isNaturalRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _min = a[0]
				var _max = a[1]

				mml.Nop(_min, _max)
				return (_isRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isInt, _min, _max)}).Values).(bool) && mml.BinaryOp(16, _min, 0).(bool))
			},
			FixedArgs: 2,
		}
		_defineRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _base = a[0]
				var _validate = a[1]
				var _min = a[2]
				var _max = a[3]

				mml.Nop(_base, _validate, _min, _max)
				return func() interface{} {
					c = _validate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min, _max)}).Values)
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _base.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["min"] = _min
							s.Values["max"] = _max
							return s
						}()
					} else {
						return _none
					}
				}()
			},
			FixedArgs: 4,
		}
		_intRange = _defineRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _intRangeType, _isRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isInt)}).Values))}).Values)
		_floatRange = _defineRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _floatRangeType, _isRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isFloat)}).Values))}).Values)
		_stringRangeType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "string")}).Values)
		_stringRange = _defineRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _stringRangeType, _isNaturalRange)}).Values)
		_listType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "list")}).Values)
		_listRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _item = a[0]
				var _min = a[1]
				var _max = a[2]

				mml.Nop(_item, _min, _max)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _defineRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _isNaturalRange, _min, _max)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["item"] = _item
					return s
				}()
			},
			FixedArgs: 3,
		}
		_listOf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _item = a[0]

				mml.Nop(_item)
				return _listRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _type.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _item)}).Values), 0, mml.Ref(_ints, "max"))}).Values)
			},
			FixedArgs: 1,
		}
		exports["listOf"] = _listOf
		_structOf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }(), _s)}).Values)
					if c.(bool) {
						return _s
					} else {
						return _none
					}
				}()
			},
			FixedArgs: 1,
		}
		exports["structOf"] = _structOf
		_range = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _min = a[1]
				var _max = a[2]

				mml.Nop(_match, _min, _max)
				var _m interface{}
				mml.Nop(_m)
				_m = _type.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values)
				switch {
				case mml.BinaryOp(11, _m, _integer):

					mml.Nop()
					return _intRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min, _max)}).Values)
				case mml.BinaryOp(11, _m, _floating):

					mml.Nop()
					return _floatRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min, _max)}).Values)
				case mml.BinaryOp(11, _m, _stringType):

					mml.Nop()
					return _stringRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min, _max)}).Values)
				case _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _m)}).Values):

					mml.Nop()
					return _listRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_m, "item"), _min, _max)}).Values)
				default:

					mml.Nop()
					return _none
				}
				return nil
			},
			FixedArgs: 3,
		}
		exports["range"] = _range
		_unionType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "union")}).Values)
		_intersectType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "intersection")}).Values)
		_predicateType = _defineType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "predicate")}).Values)
		_or = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _matches interface{}
				_matches = &mml.List{a[0:]}

				mml.Nop(_matches)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _unionType.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["matches"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _type, _matches)}).Values)
					return s
				}()
			},
			FixedArgs: 0,
		}
		exports["or"] = _or
		_and = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _matches interface{}
				_matches = &mml.List{a[0:]}

				mml.Nop(_matches)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _intersectType.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["matches"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _type, _matches)}).Values)
					return s
				}()
			},
			FixedArgs: 0,
		}
		exports["and"] = _and
		_predicate = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]

				mml.Nop(_p)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _predicateType.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["check"] = _p
					return s
				}()
			},
			FixedArgs: 1,
		}
		exports["predicate"] = _predicate
		_isSimpleType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _t = a[0]

				mml.Nop(_t)
				return _some.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _bind.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _eq, _t)}).Values), &mml.List{Values: append([]interface{}{}, _integer, _floating, _stringType, _boolean, _function, _errorType, _channel)})}).Values)
			},
			FixedArgs: 1,
		}
		_isComplexType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _t = a[0]

				mml.Nop(_t)
				return ((((_isStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _t)}).Values).(bool) && _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "token", _t)}).Values).(bool)) && mml.BinaryOp(11, mml.Ref(_t, "token"), _token).(bool)) && _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _t)}).Values).(bool)) && _isString.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_t, "type"))}).Values).(bool))
			},
			FixedArgs: 1,
		}
		_isType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _t = a[0]

				mml.Nop(_t)
				return (_isSimpleType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _t)}).Values).(bool) || _isComplexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _t)}).Values).(bool))
			},
			FixedArgs: 1,
		}
		_typesMatch = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _type = a[0]
				var _value = a[1]

				mml.Nop(_type, _value)
				return ((_isComplexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _type)}).Values).(bool) && _isComplexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool)) && mml.BinaryOp(11, mml.Ref(_type, "type"), mml.Ref(_value, "type")).(bool))
			},
			FixedArgs: 2,
		}
		_primitives = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["int"] = func() interface{} {
				s := &mml.Struct{Values: make(map[string]interface{})}
				s.Values["checkValue"] = _isInt
				s.Values["type"] = _integer
				s.Values["rangeType"] = _intRangeType
				s.Values["rangeValue"] = _identity
				return s
			}()
			s.Values["float"] = func() interface{} {
				s := &mml.Struct{Values: make(map[string]interface{})}
				s.Values["checkValue"] = _isFloat
				s.Values["type"] = _floating
				s.Values["rangeType"] = _floatRangeType
				s.Values["rangeValue"] = _identity
				return s
			}()
			s.Values["string"] = func() interface{} {
				s := &mml.Struct{Values: make(map[string]interface{})}
				s.Values["checkValue"] = _isString
				s.Values["type"] = _stringType
				s.Values["rangeType"] = _stringRangeType
				s.Values["rangeValue"] = _len
				return s
			}()
			return s
		}()
		_matchPrimitive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _def = a[0]
				var _match = a[1]
				var _value = a[2]

				mml.Nop(_def, _match, _value)

				mml.Nop()
				switch {
				case !mml.Ref(_def, "checkValue").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool):

					mml.Nop()
					return false
				case mml.BinaryOp(11, _match, mml.Ref(_def, "type")):

					mml.Nop()
					return true
				case _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_def, "rangeType"), _match)}).Values):
					var _rv interface{}
					mml.Nop(_rv)
					_rv = mml.Ref(_def, "rangeValue").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
					return (mml.BinaryOp(16, _rv, mml.Ref(_match, "min")).(bool) && mml.BinaryOp(14, _rv, mml.Ref(_match, "max")).(bool))
				default:

					mml.Nop()
					return false
				}
				return nil
			},
			FixedArgs: 3,
		}
		_matchInt = _matchPrimitive.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_primitives, "int"))}).Values)
		_matchFloat = _matchPrimitive.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_primitives, "float"))}).Values)
		_matchString = _matchPrimitive.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_primitives, "string"))}).Values)
		_matchToList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)

				mml.Nop()
				c = mml.BinaryOp(13, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values), _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values))
				if c.(bool) {
					mml.Nop()
					return false
				}
				for _i := 0; _i < _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values).(int); _i++ {

					mml.Nop()
					c = !_is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_match, _i), mml.Ref(_value, _i))}).Values).(bool)
					if c.(bool) {
						mml.Nop()
						return false
					}
				}
				return true
				return nil
			},
			FixedArgs: 2,
		}
		_matchToListType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)
				return ((mml.BinaryOp(16, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values), mml.Ref(_match, "min")).(bool) && mml.BinaryOp(14, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values), mml.Ref(_match, "max")).(bool)) && _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_match, "item"))}).Values), _value)}).Values).(bool))
			},
			FixedArgs: 2,
		}
		_matchList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)

				mml.Nop()
				switch {
				case !_isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool):

					mml.Nop()
					return false
				case _isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values):

					mml.Nop()
					return _matchToList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _match)}).Values):

					mml.Nop()
					return _matchToListType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				default:

					mml.Nop()
					return false
				}
				return nil
			},
			FixedArgs: 2,
		}
		_matchStruct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)
				return (_isStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool) && _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _key = a[0]

						mml.Nop(_key)
						return (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _key, _value)}).Values).(bool) && _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_match, _key), mml.Ref(_value, _key))}).Values).(bool))
					},
					FixedArgs: 1,
				}, _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values))}).Values).(bool))
			},
			FixedArgs: 2,
		}
		_matchUnion = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)
				return _some.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]

						mml.Nop(_m)
						return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m, _value)}).Values)
					},
					FixedArgs: 1,
				}, mml.Ref(_match, "matches"))}).Values)
			},
			FixedArgs: 2,
		}
		_matchIntersection = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)
				return _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]

						mml.Nop(_m)
						return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m, _value)}).Values)
					},
					FixedArgs: 1,
				}, mml.Ref(_match, "matches"))}).Values)
			},
			FixedArgs: 2,
		}
		_matchOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]

				mml.Nop(_match, _value)

				mml.Nop()
				switch {
				case mml.BinaryOp(11, _match, _none):

					mml.Nop()
					return false
				case mml.BinaryOp(11, _match, _any):

					mml.Nop()
					return true
				case _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _predicateType, _match)}).Values):

					mml.Nop()
					return mml.Ref(_match, "check").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
				case (_isType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool) && mml.BinaryOp(12, _value, _function).(bool)):

					mml.Nop()
					return false
				case mml.BinaryOp(11, _match, _value):

					mml.Nop()
					return true
				case (mml.BinaryOp(11, _match, _integer).(bool) || _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _intRangeType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case (mml.BinaryOp(11, _match, _floating).(bool) || _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _floatRangeType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchFloat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case (mml.BinaryOp(11, _match, _stringType).(bool) || _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _stringRangeType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchString.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case mml.BinaryOp(11, _match, _boolean):

					mml.Nop()
					return _isBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
				case mml.BinaryOp(11, _match, _function):

					mml.Nop()
					return _isFunction.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
				case mml.BinaryOp(11, _match, _channel):

					mml.Nop()
					return _isChannel.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
				case mml.BinaryOp(11, _match, _errorType):

					mml.Nop()
					return _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
				case (_isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values).(bool) || _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _unionType, _match)}).Values):

					mml.Nop()
					return _matchUnion.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _typesMatch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _intersectType, _match)}).Values):

					mml.Nop()
					return _matchIntersection.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _isStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values):

					mml.Nop()
					return _matchStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				default:

					mml.Nop()
					return false
				}
				return nil
			},
			FixedArgs: 2,
		}
		_is = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _values interface{}
				_values = &mml.List{a[1:]}

				mml.Nop(_match, _values)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _values)}).Values), 0)
					if c.(bool) {
						return _bind.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is, _match)}).Values)
					} else {
						return _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _matchOne.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values), _values)}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		exports["is"] = _is
		return exports
	})
	modulePath = "floats.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _min interface{}
		var _max interface{}
		mml.Nop(_min, _max)
		_min = mml.UnaryOp(2, 9000)
		exports["min"] = _min
		_max = 9000
		exports["max"] = _max
		return exports
	})
	modulePath = "code.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _controlStatement interface{}
		var _breakControl interface{}
		var _continueControl interface{}
		var _unaryOp interface{}
		var _binaryNot interface{}
		var _plus interface{}
		var _minus interface{}
		var _logicalNot interface{}
		var _binaryOp interface{}
		var _binaryAnd interface{}
		var _binaryOr interface{}
		var _xor interface{}
		var _andNot interface{}
		var _lshift interface{}
		var _rshift interface{}
		var _mul interface{}
		var _div interface{}
		var _mod interface{}
		var _add interface{}
		var _sub interface{}
		var _equals interface{}
		var _notEq interface{}
		var _less interface{}
		var _lessOrEq interface{}
		var _greater interface{}
		var _greaterOrEq interface{}
		var _logicalAnd interface{}
		var _logicalOr interface{}
		var _builtin interface{}
		var _flattenedStatements interface{}
		var _getModuleName interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _onlyErr interface{}
		var _passErr interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _is interface{}
		mml.Nop(_controlStatement, _breakControl, _continueControl, _unaryOp, _binaryNot, _plus, _minus, _logicalNot, _binaryOp, _binaryAnd, _binaryOr, _xor, _andNot, _lshift, _rshift, _mul, _div, _mod, _add, _sub, _equals, _notEq, _less, _lessOrEq, _greater, _greaterOrEq, _logicalAnd, _logicalOr, _builtin, _flattenedStatements, _getModuleName, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _bind, _identity, _eq, _any, _function, _channel, _type, _listOf, _structOf, _range, _or, _and, _predicate, _is)
		var __lang = mml.Modules.Use("lang.mml")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_onlyErr = __lang.Values["onlyErr"]
		_passErr = __lang.Values["passErr"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_predicate = __lang.Values["predicate"]
		_is = __lang.Values["is"]
		_controlStatement = _enum.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["controlStatement"] = _controlStatement
		_breakControl = _controlStatement.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["breakControl"] = _breakControl
		_continueControl = _controlStatement.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["continueControl"] = _continueControl
		_unaryOp = _enum.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["unaryOp"] = _unaryOp
		_binaryNot = _unaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["binaryNot"] = _binaryNot
		_plus = _unaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["plus"] = _plus
		_minus = _unaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["minus"] = _minus
		_logicalNot = _unaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["logicalNot"] = _logicalNot
		_binaryOp = _enum.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["binaryOp"] = _binaryOp
		_binaryAnd = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["binaryAnd"] = _binaryAnd
		_binaryOr = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["binaryOr"] = _binaryOr
		_xor = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["xor"] = _xor
		_andNot = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["andNot"] = _andNot
		_lshift = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["lshift"] = _lshift
		_rshift = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["rshift"] = _rshift
		_mul = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["mul"] = _mul
		_div = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["div"] = _div
		_mod = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["mod"] = _mod
		_add = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["add"] = _add
		_sub = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["sub"] = _sub
		_equals = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["equals"] = _equals
		_notEq = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["notEq"] = _notEq
		_less = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["less"] = _less
		_lessOrEq = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["lessOrEq"] = _lessOrEq
		_greater = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["greater"] = _greater
		_greaterOrEq = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["greaterOrEq"] = _greaterOrEq
		_logicalAnd = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["logicalAnd"] = _logicalAnd
		_logicalOr = _binaryOp.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
		exports["logicalOr"] = _logicalOr
		_builtin = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["len"] = "Len"
			s.Values["isError"] = "IsError"
			s.Values["keys"] = "Keys"
			s.Values["format"] = "Format"
			s.Values["stdin"] = "Stdin"
			s.Values["stdout"] = "Stdout"
			s.Values["stderr"] = "Stderr"
			s.Values["int"] = "Int"
			s.Values["float"] = "Float"
			s.Values["string"] = "String"
			s.Values["bool"] = "Bool"
			s.Values["has"] = "Has"
			s.Values["isBool"] = "IsBool"
			s.Values["isInt"] = "IsInt"
			s.Values["isFloat"] = "IsFloat"
			s.Values["isString"] = "IsString"
			s.Values["isList"] = "IsList"
			s.Values["isStruct"] = "IsStruct"
			s.Values["isFunction"] = "IsFunction"
			s.Values["isChannel"] = "IsChannel"
			s.Values["exit"] = "Exit"
			s.Values["error"] = "Error"
			s.Values["panic"] = "Panic"
			s.Values["open"] = "Open"
			s.Values["close"] = "Close"
			s.Values["args"] = "Args"
			s.Values["parseAST"] = "ParseAST"
			s.Values["parseInt"] = "ParseInt"
			s.Values["parseFloat"] = "ParseFloat"
			return s
		}()
		exports["builtin"] = _builtin
		_flattenedStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _itemType = a[0]
				var _listType = a[1]
				var _listProp = a[2]
				var _statements = a[3]

				mml.Nop(_itemType, _listType, _listProp, _statements)
				var _type interface{}
				var _toList interface{}
				mml.Nop(_type, _toList)
				_type = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _s)}).Values).(bool) && _contains.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "type"), &mml.List{Values: append([]interface{}{}, _itemType, _listType)})}).Values).(bool))
					},
					FixedArgs: 1,
				}
				_toList = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return func() interface{} {
							c = mml.BinaryOp(11, mml.Ref(_s, "type"), _itemType)
							if c.(bool) {
								return &mml.List{Values: append([]interface{}{}, _s)}
							} else {
								return mml.Ref(_s, _listProp)
							}
						}()
					},
					FixedArgs: 1,
				}
				return _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _toList)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _type)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _statements)}).Values))}).Values))}).Values)
				return nil
			},
			FixedArgs: 4,
		}
		exports["flattenedStatements"] = _flattenedStatements
		_getModuleName = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]

				mml.Nop(_path)
				return _path
			},
			FixedArgs: 1,
		}
		exports["getModuleName"] = _getModuleName
		return exports
	})
	modulePath = "parse.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _parseString interface{}
		var _spread interface{}
		var _expressionList interface{}
		var _list interface{}
		var _mutableList interface{}
		var _expressionKey interface{}
		var _entry interface{}
		var _struct interface{}
		var _mutableStruct interface{}
		var _statementList interface{}
		var _function interface{}
		var _effect interface{}
		var _symbolIndex interface{}
		var _expressionIndex interface{}
		var _indexer interface{}
		var _mutableCapture interface{}
		var _valueDefinition interface{}
		var _functionDefinition interface{}
		var _assign interface{}
		var _parseSend interface{}
		var _parseReceive interface{}
		var _parseGo interface{}
		var _parseDefer interface{}
		var _receiveDefinition interface{}
		var _symbol interface{}
		var _ret interface{}
		var _functionFact interface{}
		var _range interface{}
		var _rangeIndex interface{}
		var _indexerNodes interface{}
		var _application interface{}
		var _unary interface{}
		var _binary interface{}
		var _chaining interface{}
		var _ternary interface{}
		var _parseIf interface{}
		var _parseSwitch interface{}
		var _rangeOver interface{}
		var _loop interface{}
		var _valueCapture interface{}
		var _definitions interface{}
		var _mutableDefinitions interface{}
		var _functionCapture interface{}
		var _effectCapture interface{}
		var _effectDefinitions interface{}
		var _assignCaptures interface{}
		var _parseSelect interface{}
		var _parseExport interface{}
		var _useFact interface{}
		var _parseUse interface{}
		var _parse interface{}
		var _parseFile interface{}
		var _findExportNames interface{}
		var _parseModule interface{}
		var _modules interface{}
		var _code interface{}
		var _strings interface{}
		var _errors interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _onlyErr interface{}
		var _passErr interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _channel interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _is interface{}
		mml.Nop(_parseString, _spread, _expressionList, _list, _mutableList, _expressionKey, _entry, _struct, _mutableStruct, _statementList, _function, _effect, _symbolIndex, _expressionIndex, _indexer, _mutableCapture, _valueDefinition, _functionDefinition, _assign, _parseSend, _parseReceive, _parseGo, _parseDefer, _receiveDefinition, _symbol, _ret, _functionFact, _range, _rangeIndex, _indexerNodes, _application, _unary, _binary, _chaining, _ternary, _parseIf, _parseSwitch, _rangeOver, _loop, _valueCapture, _definitions, _mutableDefinitions, _functionCapture, _effectCapture, _effectDefinitions, _assignCaptures, _parseSelect, _parseExport, _useFact, _parseUse, _parse, _parseFile, _findExportNames, _parseModule, _modules, _code, _strings, _errors, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _bind, _identity, _eq, _any, _function, _channel, _type, _listOf, _structOf, _range, _or, _and, _predicate, _is)
		var __lang = mml.Modules.Use("lang.mml")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_onlyErr = __lang.Values["onlyErr"]
		_passErr = __lang.Values["passErr"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_predicate = __lang.Values["predicate"]
		_is = __lang.Values["is"]
		_code = mml.Modules.Use("code.mml")
		_strings = mml.Modules.Use("strings.mml")
		_errors = mml.Modules.Use("errors.mml")
		_parseString = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return mml.Ref(_strings, "unescape").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "text"), 1, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "text"))}).Values), 1)))}).Values)
			},
			FixedArgs: 1,
		}
		_spread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "spread"
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_expressionList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0]

				mml.Nop(_nodes)
				return _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, _nodes)}).Values)
			},
			FixedArgs: 1,
		}
		_list = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "list"
					s.Values["values"] = _expressionList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
					s.Values["mutable"] = false
					return s
				}()
			},
			FixedArgs: 1,
		}
		_mutableList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["mutable"] = true
					return s
				}()
			},
			FixedArgs: 1,
		}
		_expressionKey = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "expression-key"
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_entry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "entry"
					s.Values["key"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_struct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "struct"
					s.Values["entries"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					s.Values["mutable"] = false
					return s
				}()
			},
			FixedArgs: 1,
		}
		_mutableStruct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _struct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["mutable"] = true
					return s
				}()
			},
			FixedArgs: 1,
		}
		_statementList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "statement-list"
					s.Values["statements"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_function = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return _functionFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
			},
			FixedArgs: 1,
		}
		_effect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _function.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["effect"] = true
					return s
				}()
			},
			FixedArgs: 1,
		}
		_symbolIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
			},
			FixedArgs: 1,
		}
		_expressionIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
			},
			FixedArgs: 1,
		}
		_indexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return _indexerNodes.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
			},
			FixedArgs: 1,
		}
		_mutableCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _valueCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["mutable"] = true
					return s
				}()
			},
			FixedArgs: 1,
		}
		_valueDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
			},
			FixedArgs: 1,
		}
		_functionDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
			},
			FixedArgs: 1,
		}
		_assign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "assign-list"
					s.Values["assignments"] = _assignCaptures.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_parseSend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "send"
					s.Values["channel"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_parseReceive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "receive"
					s.Values["channel"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_parseGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "go"
					s.Values["application"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_parseDefer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "defer"
					s.Values["application"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_receiveDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return _valueCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
			},
			FixedArgs: 1,
		}
		_symbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)

				mml.Nop()
				switch mml.Ref(_ast, "text") {
				case "break":

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "control-statement"
						s.Values["control"] = mml.Ref(_code, "breakControl")
						return s
					}()
				case "continue":

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "control-statement"
						s.Values["control"] = mml.Ref(_code, "continueControl")
						return s
					}()
				case "return":

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "ret"
						return s
					}()
				default:

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "symbol"
						s.Values["name"] = mml.Ref(_ast, "text")
						return s
					}()
				}
				return nil
			},
			FixedArgs: 1,
		}
		_ret = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 0)
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "ret"
							return s
						}()
					} else {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "ret"
							s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
							return s
						}()
					}
				}()
			},
			FixedArgs: 1,
		}
		_functionFact = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0]

				mml.Nop(_nodes)
				var _last interface{}
				var _params interface{}
				var _lastParam interface{}
				var _hasCollectParam interface{}
				var _fixedParams interface{}
				mml.Nop(_last, _params, _lastParam, _hasCollectParam, _fixedParams)
				_last = mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 1)
				_params = mml.RefRange(_nodes, nil, _last)
				_lastParam = mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values), 1)
				_hasCollectParam = (mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values), 0).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_params, _lastParam), "name"), "collect-parameter").(bool))
				_fixedParams = func() interface{} {
					c = _hasCollectParam
					if c.(bool) {
						return mml.RefRange(_params, nil, _lastParam)
					} else {
						return _params
					}
				}()
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "function"
					s.Values["params"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _p = a[0]

							mml.Nop(_p)
							return mml.Ref(_p, "name")
						},
						FixedArgs: 1,
					})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fixedParams)}).Values))}).Values)
					s.Values["collectParam"] = func() interface{} {
						c = _hasCollectParam
						if c.(bool) {
							return mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(mml.Ref(_params, _lastParam), "nodes"), 0))}).Values), "name")
						} else {
							return ""
						}
					}()
					s.Values["statement"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, _last))}).Values)
					s.Values["effect"] = false
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_range = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _v interface{}
				mml.Nop(_v)
				_v = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				return func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_ast, "name"), "range-from")
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "range-expression"
							s.Values["from"] = _v
							return s
						}()
					} else {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "range-expression"
							s.Values["to"] = _v
							return s
						}()
					}
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_rangeIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _r interface{}
				mml.Nop(_r)
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 0)
				if c.(bool) {
					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "range-expression"
						return s
					}()
				}
				_r = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)
				if c.(bool) {
					mml.Nop()
					return _r
				}
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _r.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["to"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values), "to")
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_indexerNodes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]

				mml.Nop(_n)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "indexer"
					s.Values["expression"] = func() interface{} {
						c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n)}).Values), 2)
						if c.(bool) {
							return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, 0))}).Values)
						} else {
							return _indexerNodes.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_n, nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n)}).Values), 1)))}).Values)
						}
					}()
					s.Values["index"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n)}).Values), 1)))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_application = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "function-application"
					s.Values["function"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["args"] = _expressionList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_unary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _op interface{}
				mml.Nop(_op)
				_op = mml.Ref(_code, "binaryNot")
				switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name") {
				case "plus":

					mml.Nop()
					_op = mml.Ref(_code, "plus")
				case "minus":

					mml.Nop()
					_op = mml.Ref(_code, "minus")
				case "logical-not":

					mml.Nop()
					_op = mml.Ref(_code, "logicalNot")
				}
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "unary"
					s.Values["op"] = _op
					s.Values["arg"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_binary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _op interface{}
				mml.Nop(_op)
				_op = mml.Ref(_code, "binaryAnd")
				switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 2)), "name") {
				case "xor":

					mml.Nop()
					_op = mml.Ref(_code, "xor")
				case "and-not":

					mml.Nop()
					_op = mml.Ref(_code, "andNot")
				case "lshift":

					mml.Nop()
					_op = mml.Ref(_code, "lshift")
				case "rshift":

					mml.Nop()
					_op = mml.Ref(_code, "rshift")
				case "mul":

					mml.Nop()
					_op = mml.Ref(_code, "mul")
				case "div":

					mml.Nop()
					_op = mml.Ref(_code, "div")
				case "mod":

					mml.Nop()
					_op = mml.Ref(_code, "mod")
				case "add":

					mml.Nop()
					_op = mml.Ref(_code, "add")
				case "sub":

					mml.Nop()
					_op = mml.Ref(_code, "sub")
				case "eq":

					mml.Nop()
					_op = mml.Ref(_code, "equals")
				case "not-eq":

					mml.Nop()
					_op = mml.Ref(_code, "notEq")
				case "less":

					mml.Nop()
					_op = mml.Ref(_code, "less")
				case "less-or-eq":

					mml.Nop()
					_op = mml.Ref(_code, "lessOrEq")
				case "greater":

					mml.Nop()
					_op = mml.Ref(_code, "greater")
				case "greater-or-eq":

					mml.Nop()
					_op = mml.Ref(_code, "greaterOrEq")
				case "logical-and":

					mml.Nop()
					_op = mml.Ref(_code, "logicalAnd")
				case "logical-or":

					mml.Nop()
					_op = mml.Ref(_code, "logicalOr")
				}
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "binary"
					s.Values["op"] = _op
					s.Values["left"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 3)
						if c.(bool) {
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								s.Values["name"] = mml.Ref(_ast, "name")
								s.Values["nodes"] = mml.RefRange(mml.Ref(_ast, "nodes"), nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 2))
								return s
							}()
						} else {
							return mml.Ref(mml.Ref(_ast, "nodes"), 0)
						}
					}())}).Values)
					s.Values["right"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)))}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_chaining = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _a interface{}
				var _n interface{}
				mml.Nop(_a, _n)
				_a = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				_n = mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil)
				for {

					mml.Nop()
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n)}).Values), 0)
					if c.(bool) {
						mml.Nop()
						return _a
					}
					_a = func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "function-application"
						s.Values["function"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, 0))}).Values)
						s.Values["args"] = &mml.List{Values: append([]interface{}{}, _a)}
						return s
					}()
					_n = mml.RefRange(_n, 1, nil)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_ternary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "cond"
					s.Values["condition"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["consequent"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					s.Values["alternative"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 2))}).Values)
					s.Values["ternary"] = true
					return s
				}()
			},
			FixedArgs: 1,
		}
		_parseIf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _cond interface{}
				var _alternative interface{}
				mml.Nop(_cond, _alternative)
				_cond = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "cond"
					s.Values["ternary"] = false
					s.Values["condition"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["consequent"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}()
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 2)
				if c.(bool) {
					mml.Nop()
					return _cond
				}
				_alternative = func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 3)
					if c.(bool) {
						return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 2))}).Values)
					} else {
						return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _ast.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["nodes"] = mml.RefRange(mml.Ref(_ast, "nodes"), 2, nil)
							return s
						}())}).Values)
					}
				}()
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _cond.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["alternative"] = _alternative
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_parseSwitch = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _hasExpression interface{}
				var _expression interface{}
				var _nodes interface{}
				var _groupLines interface{}
				var _cases interface{}
				var _lines interface{}
				var _s interface{}
				mml.Nop(_hasExpression, _expression, _nodes, _groupLines, _cases, _lines, _s)
				_hasExpression = (mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "case").(bool) && mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "default").(bool))
				_expression = func() interface{} {
					c = _hasExpression
					if c.(bool) {
						return mml.Ref(mml.Ref(_ast, "nodes"), 0)
					} else {
						return func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }()
					}
				}()
				_nodes = func() interface{} {
					c = _hasExpression
					if c.(bool) {
						return mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil)
					} else {
						return mml.Ref(_ast, "nodes")
					}
				}()
				_groupLines = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						var _isDefault interface{}
						var _current interface{}
						var _cases interface{}
						var _defaults interface{}
						mml.Nop(_isDefault, _current, _cases, _defaults)
						_isDefault = false
						_current = &mml.List{Values: []interface{}{}}
						_cases = &mml.List{Values: []interface{}{}}
						_defaults = &mml.List{Values: []interface{}{}}
						for _, _n := range _nodes.(*mml.List).Values {

							mml.Nop()
							switch mml.Ref(_n, "name") {
							case "case":

								mml.Nop()
								c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _current)}).Values), 0)
								if c.(bool) {
									mml.Nop()
									c = _isDefault
									if c.(bool) {
										mml.Nop()
										_defaults = _current
									} else {
										mml.Nop()
										_cases = &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _current)}
									}
								}
								_current = &mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_n, "nodes"), 0))}
								_isDefault = false
							case "default":

								mml.Nop()
								c = (mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _current)}).Values), 0).(bool) && !_isDefault.(bool))
								if c.(bool) {
									mml.Nop()
									_cases = &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _current)}
								}
								_current = &mml.List{Values: []interface{}{}}
								_isDefault = true
							default:

								mml.Nop()
								_current = &mml.List{Values: append(append([]interface{}{}, _current.(*mml.List).Values...), _n)}
							}
						}
						c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _current)}).Values), 0)
						if c.(bool) {
							mml.Nop()
							c = _isDefault
							if c.(bool) {
								mml.Nop()
								_defaults = _current
							} else {
								mml.Nop()
								_cases = &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _current)}
							}
						}
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["cases"] = _cases
							s.Values["defaults"] = _defaults
							return s
						}()
						return nil
					},
					FixedArgs: 0,
				}
				_cases = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)

						mml.Nop()
						return _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _c = a[0]

								mml.Nop(_c)
								return func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["type"] = "switch-case"
									s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, 0))}).Values)
									s.Values["body"] = func() interface{} {
										s := &mml.Struct{Values: make(map[string]interface{})}
										s.Values["type"] = "statement-list"
										s.Values["statements"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.RefRange(_c, 1, nil))}).Values)
										return s
									}()
									return s
								}()
							},
							FixedArgs: 1,
						}, _c)}).Values)
						return nil
					},
					FixedArgs: 1,
				}
				_lines = _groupLines.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				_s = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "switch-statement"
					s.Values["cases"] = _cases.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lines, "cases"))}).Values)
					s.Values["defaultStatements"] = func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "statement-list"
						s.Values["statements"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_lines, "defaults"))}).Values)
						return s
					}()
					return s
				}()
				return func() interface{} {
					c = _hasExpression
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _s.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _expression)}).Values)
							return s
						}()
					} else {
						return _s
					}
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_rangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _expression interface{}
				mml.Nop(_expression)
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 0)
				if c.(bool) {
					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "range-over"
						return s
					}()
				}
				c = (mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "symbol").(bool))
				if c.(bool) {
					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "range-over"
						s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
						return s
					}()
				}
				_expression = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _nodes = a[0]

						mml.Nop(_nodes)
						var _exp interface{}
						mml.Nop(_exp)
						_exp = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, 0))}).Values)
						c = ((!_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _exp)}).Values).(bool) || mml.BinaryOp(12, mml.Ref(_exp, "type"), "range-expression").(bool)) || mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 1).(bool))
						if c.(bool) {
							mml.Nop()
							return _exp
						}
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _exp.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["to"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, 1))}).Values), "to")
							return s
						}()
						return nil
					},
					FixedArgs: 1,
				}
				c = mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "symbol")
				if c.(bool) {
					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "range-over"
						s.Values["expression"] = _expression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
						return s
					}()
				}
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "range-over"
					s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
					s.Values["expression"] = _expression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil))}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_loop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _loop interface{}
				var _expression interface{}
				var _emptyRange interface{}
				mml.Nop(_loop, _expression, _emptyRange)
				_loop = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "loop"
					return s
				}()
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)
				if c.(bool) {
					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						func() {
							sp := _loop.(*mml.Struct)
							for k, v := range sp.Values {
								s.Values[k] = v
							}
						}()
						s.Values["body"] = _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
						return s
					}()
				}
				_expression = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				_emptyRange = (((_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _expression)}).Values).(bool) && mml.BinaryOp(11, mml.Ref(_expression, "type"), "range-over").(bool)) && !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "symbol", _expression)}).Values).(bool)) && !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _expression)}).Values).(bool))
				return func() interface{} {
					c = _emptyRange
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _loop.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["body"] = _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
							return s
						}()
					} else {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _loop.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["expression"] = _expression
							s.Values["body"] = _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
							return s
						}()
					}
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_valueCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "definition"
					s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
					s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					s.Values["mutable"] = false
					s.Values["exported"] = false
					return s
				}()
			},
			FixedArgs: 1,
		}
		_definitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "definition-list"
					s.Values["definitions"] = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _c = a[0]

							mml.Nop(_c)
							return (!_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _c)}).Values).(bool) || mml.BinaryOp(12, mml.Ref(_c, "type"), "comment").(bool))
						},
						FixedArgs: 1,
					})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_mutableDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _dl interface{}
				mml.Nop(_dl)
				_dl = _definitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _dl.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["definitions"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _d = a[0]

							mml.Nop(_d)
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _d.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["mutable"] = true
								return s
							}()
						},
						FixedArgs: 1,
					})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_dl, "definitions"))}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_functionCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "definition"
					s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
					s.Values["expression"] = _functionFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil))}).Values)
					s.Values["mutable"] = false
					s.Values["exported"] = false
					return s
				}()
			},
			FixedArgs: 1,
		}
		_effectCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _f interface{}
				mml.Nop(_f)
				_f = _functionCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _f.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["expression"] = func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						func() {
							sp := mml.Ref(_f, "expression").(*mml.Struct)
							for k, v := range sp.Values {
								s.Values[k] = v
							}
						}()
						s.Values["effect"] = true
						return s
					}()
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_effectDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _dl interface{}
				mml.Nop(_dl)
				_dl = _definitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _dl.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["definitions"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _d = a[0]

							mml.Nop(_d)
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _d.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["effect"] = true
								return s
							}()
						},
						FixedArgs: 1,
					}, mml.Ref(_dl, "definitions"))}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_assignCaptures = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0]

				mml.Nop(_nodes)

				mml.Nop()
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 0)
				if c.(bool) {
					mml.Nop()
					return &mml.List{Values: []interface{}{}}
				}
				return &mml.List{Values: append(append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "assign"
					s.Values["capture"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, 0))}).Values)
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, 1))}).Values)
					return s
				}()), _assignCaptures.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_nodes, 2, nil))}).Values).(*mml.List).Values...)}
				return nil
			},
			FixedArgs: 1,
		}
		_parseSelect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _nodes interface{}
				var _groupLines interface{}
				var _cases interface{}
				var _lines interface{}
				mml.Nop(_nodes, _groupLines, _cases, _lines)
				_nodes = mml.Ref(_ast, "nodes")
				_groupLines = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						var _isDefault interface{}
						var _hasDefault interface{}
						var _current interface{}
						var _cases interface{}
						var _defaults interface{}
						mml.Nop(_isDefault, _hasDefault, _current, _cases, _defaults)
						_isDefault = false
						_hasDefault = false
						_current = &mml.List{Values: []interface{}{}}
						_cases = &mml.List{Values: []interface{}{}}
						_defaults = &mml.List{Values: []interface{}{}}
						for _, _n := range _nodes.(*mml.List).Values {

							mml.Nop()
							switch mml.Ref(_n, "name") {
							case "case":

								mml.Nop()
								c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _current)}).Values), 0)
								if c.(bool) {
									mml.Nop()
									c = _isDefault
									if c.(bool) {
										mml.Nop()
										_defaults = _current
									} else {
										mml.Nop()
										_cases = &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _current)}
									}
								}
								_current = &mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_n, "nodes"), 0))}
								_isDefault = false
							case "default":

								mml.Nop()
								c = (mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _current)}).Values), 0).(bool) && !_isDefault.(bool))
								if c.(bool) {
									mml.Nop()
									_cases = &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _current)}
								}
								_current = &mml.List{Values: []interface{}{}}
								_isDefault = true
								_hasDefault = true
							default:

								mml.Nop()
								_current = &mml.List{Values: append(append([]interface{}{}, _current.(*mml.List).Values...), _n)}
							}
						}
						c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _current)}).Values), 0)
						if c.(bool) {
							mml.Nop()
							c = _isDefault
							if c.(bool) {
								mml.Nop()
								_defaults = _current
							} else {
								mml.Nop()
								_cases = &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _current)}
							}
						}
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["cases"] = _cases
							s.Values["defaults"] = _defaults
							s.Values["hasDefault"] = _hasDefault
							return s
						}()
						return nil
					},
					FixedArgs: 0,
				}
				_cases = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)

						mml.Nop()
						return _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _c = a[0]

								mml.Nop(_c)
								return func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["type"] = "select-case"
									s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, 0))}).Values)
									s.Values["body"] = func() interface{} {
										s := &mml.Struct{Values: make(map[string]interface{})}
										s.Values["type"] = "statement-list"
										s.Values["statements"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.RefRange(_c, 1, nil))}).Values)
										return s
									}()
									return s
								}()
							},
							FixedArgs: 1,
						}, _c)}).Values)
						return nil
					},
					FixedArgs: 1,
				}
				_lines = _groupLines.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["cases"] = _cases.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lines, "cases"))}).Values)
					s.Values["defaultStatements"] = func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "statement-list"
						s.Values["statements"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_lines, "defaults"))}).Values)
						return s
					}()
					s.Values["hasDefault"] = mml.Ref(_lines, "hasDefault")
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_parseExport = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _d interface{}
				mml.Nop(_d)
				_d = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "definition-list"
					s.Values["definitions"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _d = a[0]

							mml.Nop(_d)
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _d.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["exported"] = true
								return s
							}()
						},
						FixedArgs: 1,
					})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						c = mml.BinaryOp(11, mml.Ref(_d, "type"), "definition")
						if c.(bool) {
							return &mml.List{Values: append([]interface{}{}, _d)}
						} else {
							return mml.Ref(_d, "definitions")
						}
					}())}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_useFact = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				var _capture interface{}
				var _path interface{}
				mml.Nop(_capture, _path)
				_capture = ""
				_path = ""
				switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name") {
				case "use-inline":

					mml.Nop()
					_capture = "."
					_path = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
				case "symbol":

					mml.Nop()
					_capture = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
					_path = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
				default:

					mml.Nop()
					_path = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				}
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "use"
					s.Values["capture"] = _capture
					s.Values["path"] = _path
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_parseUse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "use-list"
					s.Values["uses"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_parse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]

				mml.Nop(_ast)

				mml.Nop()
				switch mml.Ref(_ast, "name") {
				case "line-comment-content":

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "comment"
						return s
					}()
				case "int":

					mml.Nop()
					return _parseInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "text"))}).Values)
				case "float":

					mml.Nop()
					return _parseFloat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "text"))}).Values)
				case "string":

					mml.Nop()
					return _parseString.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "true":

					mml.Nop()
					return true
				case "false":

					mml.Nop()
					return false
				case "symbol":

					mml.Nop()
					return _symbol.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "spread-expression":

					mml.Nop()
					return _spread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "list":

					mml.Nop()
					return _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "mutable-list":

					mml.Nop()
					return _mutableList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "expression-key":

					mml.Nop()
					return _expressionKey.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "entry":

					mml.Nop()
					return _entry.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "struct":

					mml.Nop()
					return _struct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "mutable-struct":

					mml.Nop()
					return _mutableStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "return-value":

					mml.Nop()
					return _ret.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "return":

					mml.Nop()
					return _ret.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "block":

					mml.Nop()
					return _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "function":

					mml.Nop()
					return _function.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "effect":

					mml.Nop()
					return _effect.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "range-from":

					mml.Nop()
					return _range.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "range-to":

					mml.Nop()
					return _range.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "symbol-index":

					mml.Nop()
					return _symbolIndex.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "expression-index":

					mml.Nop()
					return _expressionIndex.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "range-index":

					mml.Nop()
					return _rangeIndex.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "indexer":

					mml.Nop()
					return _indexer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "function-application":

					mml.Nop()
					return _application.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "unary-expression":

					mml.Nop()
					return _unary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "binary0":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "binary1":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "binary2":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "binary3":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "binary4":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "chaining":

					mml.Nop()
					return _chaining.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "ternary-expression":

					mml.Nop()
					return _ternary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "if":

					mml.Nop()
					return _parseIf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "switch":

					mml.Nop()
					return _parseSwitch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "range-over-expression":

					mml.Nop()
					return _rangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "loop":

					mml.Nop()
					return _loop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "value-capture":

					mml.Nop()
					return _valueCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "mutable-capture":

					mml.Nop()
					return _mutableCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "value-definition":

					mml.Nop()
					return _valueDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "value-definition-group":

					mml.Nop()
					return _definitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "mutable-definition-group":

					mml.Nop()
					return _mutableDefinitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "function-capture":

					mml.Nop()
					return _functionCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "effect-capture":

					mml.Nop()
					return _effectCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "function-definition":

					mml.Nop()
					return _functionDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "function-definition-group":

					mml.Nop()
					return _definitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "effect-definition-group":

					mml.Nop()
					return _effectDefinitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "assignment":

					mml.Nop()
					return _assign.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "send":

					mml.Nop()
					return _parseSend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "receive":

					mml.Nop()
					return _parseReceive.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "go":

					mml.Nop()
					return _parseGo.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "defer":

					mml.Nop()
					return _parseDefer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "receive-definition":

					mml.Nop()
					return _receiveDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "select":

					mml.Nop()
					return _parseSelect.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "export":

					mml.Nop()
					return _parseExport.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "use-fact":

					mml.Nop()
					return _useFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "use":

					mml.Nop()
					return _parseUse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				default:

					mml.Nop()
					return _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_parseFile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]

				mml.Nop(_path)
				var _in interface{}
				var _ast interface{}
				mml.Nop(_in, _ast)
				_in = _open.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path)}).Values)
				c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _in)}).Values)
				if c.(bool) {
					mml.Nop()
					return _in
				}
				defer _close.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _in)}).Values)
				_ast = _passErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parseAST)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _in.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.UnaryOp(2, 1))}).Values))}).Values)
				c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				if c.(bool) {
					mml.Nop()
					return _ast
				}
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_findExportNames = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statements = a[0]

				mml.Nop(_statements)
				return _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]

						mml.Nop(_d)
						return mml.Ref(_d, "symbol")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]

						mml.Nop(_d)
						return mml.Ref(_d, "exported")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "flattenedStatements").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definition", "definition-list", "definitions")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _statements)}).Values))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_parseModule = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _entryPath = a[1]

				mml.Nop(_context, _entryPath)
				var _module interface{}
				var _uses interface{}
				var _usesModules interface{}
				var _statements interface{}
				var _currentCode interface{}
				var _modules interface{}
				mml.Nop(_module, _uses, _usesModules, _statements, _currentCode, _modules)
				c = _contains.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _entryPath, mml.Ref(_context, "stack"))}).Values)
				if c.(bool) {
					mml.Nop()
					return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "circular module dependency: %s", _entryPath)}).Values))}).Values)
				}
				c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _entryPath, mml.Ref(_context, "parsed"))}).Values)
				if c.(bool) {
					mml.Nop()
					return mml.Ref(mml.Ref(_context, "parsed"), _entryPath)
				}
				_module = _parseFile.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _entryPath)}).Values)
				_uses = mml.Ref(_code, "flattenedStatements").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "use", "use-list", "uses", mml.Ref(_module, "statements"))}).Values)
				c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _module)}).Values)
				if c.(bool) {
					mml.Nop()
					return _module
				}
				mml.SetRef(_context, "stack", &mml.List{Values: append(append([]interface{}{}, mml.Ref(_context, "stack").(*mml.List).Values...), _entryPath)})
				_usesModules = _passErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uniq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _left = a[0]
						var _right = a[1]

						mml.Nop(_left, _right)
						return mml.BinaryOp(11, mml.Ref(_left, "path"), mml.Ref(_right, "path"))
					},
					FixedArgs: 2,
				})}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _passErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]

						mml.Nop(_m)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = mml.Ref(_m, "type")
							s.Values["path"] = mml.Ref(_m, "path")
							s.Values["statements"] = mml.Ref(_m, "statements")
							s.Values["exportNames"] = _findExportNames.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_m, "statements"))}).Values)
							return s
						}()
					},
					FixedArgs: 1,
				})}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _passErr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _flat)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_errors, "any").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parseModule.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.BinaryOp(9, mml.Ref(_u, "path"), ".mml")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values))}).Values))}).Values))}).Values))}).Values))}).Values)
				mml.SetRef(_context, "stack", mml.RefRange(mml.Ref(_context, "stack"), nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_context, "stack"))}).Values), 1)))
				c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _usesModules)}).Values)
				if c.(bool) {
					mml.Nop()
					return _usesModules
				}
				_statements = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)

						mml.Nop()
						c = (!_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _s)}).Values).(bool) || (mml.BinaryOp(12, mml.Ref(_s, "type"), "use").(bool) && mml.BinaryOp(12, mml.Ref(_s, "type"), "use-list").(bool)))
						if c.(bool) {
							mml.Nop()
							return _s
						}
						c = mml.BinaryOp(11, mml.Ref(_s, "type"), "use")
						if c.(bool) {
							var _m interface{}
							mml.Nop(_m)
							_m = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
								F: func(a []interface{}) interface{} {
									var c interface{}
									mml.Nop(c)
									var _m = a[0]

									mml.Nop(_m)
									return mml.BinaryOp(11, mml.Ref(_m, "path"), mml.Ref(_s, "path"))
								},
								FixedArgs: 1,
							}, _usesModules)}).Values)
							c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m)}).Values), 0)
							if c.(bool) {
								mml.Nop()
								return _s
							}
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _s.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["exportNames"] = mml.Ref(mml.Ref(_m, 0), "exportNames")
								return s
							}()
						}
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = mml.Ref(_s, "type")
							s.Values["uses"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
								F: func(a []interface{}) interface{} {
									var c interface{}
									mml.Nop(c)
									var _u = a[0]

									mml.Nop(_u)
									var _m interface{}
									mml.Nop(_m)
									_m = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
										F: func(a []interface{}) interface{} {
											var c interface{}
											mml.Nop(c)
											var _m = a[0]

											mml.Nop(_m)
											return mml.BinaryOp(11, mml.Ref(_m, "path"), mml.BinaryOp(9, mml.Ref(_u, "path"), ".mml"))
										},
										FixedArgs: 1,
									}, _usesModules)}).Values)
									c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m)}).Values), 0)
									if c.(bool) {
										mml.Nop()
										return _u
									}
									return func() interface{} {
										s := &mml.Struct{Values: make(map[string]interface{})}
										func() {
											sp := _u.(*mml.Struct)
											for k, v := range sp.Values {
												s.Values[k] = v
											}
										}()
										s.Values["exportNames"] = mml.Ref(mml.Ref(_m, 0), "exportNames")
										return s
									}()
									return nil
								},
								FixedArgs: 1,
							}, mml.Ref(_s, "uses"))}).Values)
							return s
						}()
						return nil
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_module, "statements"))}).Values)
				_currentCode = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _module.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["path"] = _entryPath
					s.Values["statements"] = _statements
					return s
				}()
				_modules = &mml.List{Values: append(append([]interface{}{}, _currentCode), _usesModules.(*mml.List).Values...)}
				mml.SetRef(mml.Ref(_context, "parsed"), _entryPath, _modules)
				return _modules
				return nil
			},
			FixedArgs: 2,
		}
		_modules = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _entryPath = a[0]

				mml.Nop(_entryPath)
				return _parseModule.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["stack"] = &mml.List{Values: []interface{}{}}
					s.Values["parsed"] = func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }()
					return s
				}(), _entryPath)}).Values)
			},
			FixedArgs: 1,
		}
		exports["modules"] = _modules
		return exports
	})
	modulePath = "definitions.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _newContext interface{}
		var _extend interface{}
		var _importContext interface{}
		var _definedCurrent interface{}
		var _define interface{}
		var _assign interface{}
		var _defineImport interface{}
		var _defined interface{}
		var _capture interface{}
		var _values interface{}
		var _results interface{}
		var _resultValues interface{}
		var _resultErrors interface{}
		var _dropValues interface{}
		var _emptyResults interface{}
		var _mergeResults interface{}
		var _wrapWithReturn interface{}
		var _undefined interface{}
		var _duplicate interface{}
		var _all interface{}
		var _scoped interface{}
		var _allScoped interface{}
		var _fields interface{}
		var _fieldsIfHas interface{}
		var _list interface{}
		var _struct interface{}
		var _rangeExpression interface{}
		var _spread interface{}
		var _unary interface{}
		var _binary interface{}
		var _validateSend interface{}
		var _validateGo interface{}
		var _validateDefer interface{}
		var _definitions interface{}
		var _assignments interface{}
		var _ret interface{}
		var _useList interface{}
		var _expandFunction interface{}
		var _symbol interface{}
		var _entry interface{}
		var _function interface{}
		var _indexer interface{}
		var _application interface{}
		var _cond interface{}
		var _validateCase interface{}
		var _validateSwitch interface{}
		var _validateReceive interface{}
		var _validateSelect interface{}
		var _rangeOver interface{}
		var _loop interface{}
		var _definition interface{}
		var _assignment interface{}
		var _validateUse interface{}
		var _statements interface{}
		var _do interface{}
		var _validate interface{}
		var _mmlcode interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _onlyErr interface{}
		var _passErr interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _channel interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _is interface{}
		mml.Nop(_newContext, _extend, _importContext, _definedCurrent, _define, _assign, _defineImport, _defined, _capture, _values, _results, _resultValues, _resultErrors, _dropValues, _emptyResults, _mergeResults, _wrapWithReturn, _undefined, _duplicate, _all, _scoped, _allScoped, _fields, _fieldsIfHas, _list, _struct, _rangeExpression, _spread, _unary, _binary, _validateSend, _validateGo, _validateDefer, _definitions, _assignments, _ret, _useList, _expandFunction, _symbol, _entry, _function, _indexer, _application, _cond, _validateCase, _validateSwitch, _validateReceive, _validateSelect, _rangeOver, _loop, _definition, _assignment, _validateUse, _statements, _do, _validate, _mmlcode, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _bind, _identity, _eq, _any, _function, _channel, _type, _listOf, _structOf, _range, _or, _and, _predicate, _is)
		var __lang = mml.Modules.Use("lang.mml")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_onlyErr = __lang.Values["onlyErr"]
		_passErr = __lang.Values["passErr"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_predicate = __lang.Values["predicate"]
		_is = __lang.Values["is"]
		_mmlcode = mml.Modules.Use("code.mml")
		_newContext = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)

				mml.Nop()
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["definitions"] = func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }()
					s.Values["unexpanded"] = &mml.List{Values: []interface{}{}}
					s.Values["capturing"] = false
					return s
				}()
			},
			FixedArgs: 0,
		}
		_extend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]

				mml.Nop(_context)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _newContext.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["parent"] = _context
					return s
				}()
			},
			FixedArgs: 1,
		}
		_importContext = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]

				mml.Nop(_context)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["imports"] = true
					return s
				}()
			},
			FixedArgs: 1,
		}
		_definedCurrent = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]

				mml.Nop(_context, _n)
				return _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n, mml.Ref(_context, "definitions"))}).Values)
			},
			FixedArgs: 2,
		}
		_define = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]
				var _v = a[2]

				mml.Nop(_context, _n, _v)
				return _capture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _n, _v)}).Values)
			},
			FixedArgs: 3,
		}
		_assign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]
				var _v = a[2]

				mml.Nop(_context, _n, _v)
				return _capture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _n, _v)}).Values)
			},
			FixedArgs: 3,
		}
		_defineImport = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]
				var _v = a[2]

				mml.Nop(_context, _n, _v)

				mml.Nop()
				c = !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "imports", _context)}).Values).(bool)
				if c.(bool) {
					mml.Nop()
					return _defineImport.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_context, "parent"), _n, _v)}).Values)
				}
				return _capture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _n, _v)}).Values)
				return nil
			},
			FixedArgs: 3,
		}
		_defined = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]

				mml.Nop(_context, _n)
				return (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n, mml.Ref(_context, "definitions"))}).Values).(bool) || (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "parent", _context)}).Values).(bool) && _defined.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_context, "parent"), _n)}).Values).(bool)))
			},
			FixedArgs: 2,
		}
		_capture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]
				var _v = a[2]

				mml.Nop(_context, _n, _v)
				return mml.SetRef(mml.Ref(_context, "definitions"), _n, func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n, mml.Ref(_context, "definitions"))}).Values)
					if c.(bool) {
						return &mml.List{Values: append(append([]interface{}{}, mml.Ref(mml.Ref(_context, "definitions"), _n).(*mml.List).Values...), _v.(*mml.List).Values...)}
					} else {
						return _v
					}
				}())
			},
			FixedArgs: 3,
		}
		_values = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _n = a[1]

				mml.Nop(_context, _n)
				return func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _n, mml.Ref(_context, "definitions"))}).Values)
					if c.(bool) {
						return mml.Ref(mml.Ref(_context, "definitions"), _n)
					} else {
						return func() interface{} {
							c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "parent", _context)}).Values)
							if c.(bool) {
								return _values.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_context, "parent"), _n)}).Values)
							} else {
								return &mml.List{Values: []interface{}{}}
							}
						}()
					}
				}()
			},
			FixedArgs: 2,
		}
		_results = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0]
				var _e = a[1]

				mml.Nop(_v, _e)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["values"] = _v
					s.Values["errors"] = _e
					return s
				}()
			},
			FixedArgs: 2,
		}
		_resultValues = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v interface{}
				_v = &mml.List{a[0:]}

				mml.Nop(_v)
				return _results.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _v, &mml.List{Values: []interface{}{}})}).Values)
			},
			FixedArgs: 0,
		}
		_resultErrors = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e interface{}
				_e = &mml.List{a[0:]}

				mml.Nop(_e)
				return _results.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: []interface{}{}}, _e)}).Values)
			},
			FixedArgs: 0,
		}
		_dropValues = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				return _resultErrors.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "errors").(*mml.List).Values...)}).Values)
			},
			FixedArgs: 1,
		}
		_emptyResults = _results.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: []interface{}{}}, &mml.List{Values: []interface{}{}})}).Values)
		_mergeResults = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r interface{}
				_r = &mml.List{a[0:]}

				mml.Nop(_r)
				var _mergeTwo interface{}
				mml.Nop(_mergeTwo)
				_mergeTwo = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _left = a[0]
						var _right = a[1]

						mml.Nop(_left, _right)
						return _results.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append(append([]interface{}{}, mml.Ref(_left, "values").(*mml.List).Values...), mml.Ref(_right, "values").(*mml.List).Values...)}, &mml.List{Values: append(append([]interface{}{}, mml.Ref(_left, "errors").(*mml.List).Values...), mml.Ref(_right, "errors").(*mml.List).Values...)})}).Values)
					},
					FixedArgs: 2,
				}
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeTwo, _emptyResults, _r)}).Values)
				return nil
			},
			FixedArgs: 0,
		}
		_wrapWithReturn = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				return _results.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _v = a[0]

						mml.Nop(_v)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "ret"
							s.Values["value"] = _v
							return s
						}()
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "values"))}).Values), mml.Ref(_r, "errors"))}).Values)
			},
			FixedArgs: 1,
		}
		_undefined = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _name = a[0]

				mml.Nop(_name)
				return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "undefined: %s", _name)}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_duplicate = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _name = a[0]

				mml.Nop(_name)
				return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "duplicate definition: %s", _name)}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_all = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _l = a[1]

				mml.Nop(_context, _l)
				return (&mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _r = a[0]

						mml.Nop(_r)
						return _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r.(*mml.List).Values...)}).Values)
					},
					FixedArgs: 1,
				}).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_scoped = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _code = a[1]

				mml.Nop(_context, _code)
				return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values), _code)}).Values)
			},
			FixedArgs: 2,
		}
		_allScoped = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _l = a[1]

				mml.Nop(_context, _l)
				return (&mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _r = a[0]

						mml.Nop(_r)
						return _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r.(*mml.List).Values...)}).Values)
					},
					FixedArgs: 1,
				}).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _scoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_fields = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]
				var _f = a[2]

				mml.Nop(_context, _s, _f)
				return _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]

						mml.Nop(_f)
						return mml.Ref(_s, _f)
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f)}).Values))}).Values)
			},
			FixedArgs: 3,
		}
		_fieldsIfHas = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _f = a[1]
				var _s = a[2]

				mml.Nop(_context, _f, _s)
				return _fields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _s)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]

						mml.Nop(_f)
						return _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _s)}).Values)
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f)}).Values))}).Values)
			},
			FixedArgs: 3,
		}
		_list = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _l = a[1]

				mml.Nop(_context, _l)
				return _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_l, "values"))}).Values)
			},
			FixedArgs: 2,
		}
		_struct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				return _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "entries"))}).Values)
			},
			FixedArgs: 2,
		}
		_rangeExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _r = a[1]

				mml.Nop(_context, _r)
				return _fieldsIfHas.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, &mml.List{Values: append([]interface{}{}, "from", "to")}, _r)}).Values)
			},
			FixedArgs: 2,
		}
		_spread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "value"))}).Values)
			},
			FixedArgs: 2,
		}
		_unary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _u = a[1]

				mml.Nop(_context, _u)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_u, "arg"))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_binary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _b = a[1]

				mml.Nop(_context, _b)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _b, &mml.List{Values: append([]interface{}{}, "left", "right")})}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_validateSend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _s, &mml.List{Values: append([]interface{}{}, "channel", "value")})}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_validateGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _g = a[1]

				mml.Nop(_context, _g)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_g, "application"))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_validateDefer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _d = a[1]

				mml.Nop(_context, _d)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_d, "application"))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_definitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _d = a[1]

				mml.Nop(_context, _d)
				return _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_d, "definitions"))}).Values)
			},
			FixedArgs: 2,
		}
		_assignments = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _a = a[1]

				mml.Nop(_context, _a)
				return _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_a, "assignments"))}).Values)
			},
			FixedArgs: 2,
		}
		_ret = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _r = a[1]

				mml.Nop(_context, _r)
				return _fieldsIfHas.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, &mml.List{Values: append([]interface{}{}, "value")}, _r)}).Values)
			},
			FixedArgs: 2,
		}
		_useList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _u = a[1]

				mml.Nop(_context, _u)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_u, "uses"))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_expandFunction = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				var _c interface{}
				mml.Nop(_c)
				c = mml.Ref(_f, "expanded")
				if c.(bool) {
					mml.Nop()
					return _emptyResults
				}
				mml.SetRef(_f, "expanded", true)
				_c = _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "context"))}).Values)
				for _, _p := range mml.Ref(_f, "params").(*mml.List).Values {

					mml.Nop()
					_define.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c, _p, &mml.List{Values: []interface{}{}})}).Values)
				}
				c = mml.BinaryOp(12, mml.Ref(_f, "collectParam"), "")
				if c.(bool) {
					mml.Nop()
					_define.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c, mml.Ref(_f, "collectParam"), &mml.List{Values: []interface{}{}})}).Values)
				}
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c, mml.Ref(_f, "statement"))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_symbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				var _r interface{}
				mml.Nop(_r)
				_r = func() interface{} {
					c = _defined.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "name"))}).Values)
					if c.(bool) {
						return _resultValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _values.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "name"))}).Values).(*mml.List).Values...)}).Values)
					} else {
						return _resultErrors.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _undefined.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "name"))}).Values))}).Values)
					}
				}()
				c = mml.Ref(_context, "capturing")
				if c.(bool) {
					mml.Nop()
					return _r
				}
				for _, _v := range mml.Ref(_r, "values").(*mml.List).Values {

					mml.Nop()
					c = (!_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _v)}).Values).(bool) || mml.BinaryOp(12, mml.Ref(_v, "type"), "function").(bool))
					if c.(bool) {
						mml.Nop()
						continue
					}
					_r = _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r, _expandFunction.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _v)}).Values))}).Values)
				}
				return _r
				return nil
			},
			FixedArgs: 2,
		}
		_entry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _e = a[1]

				mml.Nop(_context, _e)
				var _kr interface{}
				mml.Nop(_kr)
				_kr = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, func() interface{} {
					c = (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(_e, "key"))}).Values).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_e, "key"), "type"), "symbol").(bool))
					if c.(bool) {
						return mml.Ref(mml.Ref(_e, "key"), "name")
					} else {
						return mml.Ref(_e, "key")
					}
				}())}).Values)
				return _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _kr, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_e, "value"))}).Values))}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_function = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _f = a[1]

				mml.Nop(_context, _f)
				var _ff interface{}
				mml.Nop(_ff)
				_ff = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _f.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["context"] = _context
					s.Values["expanded"] = false
					return s
				}()
				c = mml.Ref(_context, "capturing")
				if c.(bool) {
					mml.Nop()
					mml.SetRef(_context, "unexpanded", &mml.List{Values: append(append([]interface{}{}, mml.Ref(_context, "unexpanded").(*mml.List).Values...), _ff)})
					return _resultValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ff)}).Values)
				}
				return _expandFunction.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ff)}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_indexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _i = a[1]

				mml.Nop(_context, _i)
				return _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_i, "index"))}).Values))}).Values), _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_i, "expression"))}).Values))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_application = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _a = a[1]

				mml.Nop(_context, _a)
				var _capturing interface{}
				var _r interface{}
				mml.Nop(_capturing, _r)
				_capturing = mml.Ref(_context, "capturing")
				mml.SetRef(_context, "capturing", false)
				_r = _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_a, "function"))}).Values))}).Values), _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _all.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_a, "args"))}).Values))}).Values))}).Values)
				mml.SetRef(_context, "capturing", _capturing)
				return _r
				return nil
			},
			FixedArgs: 2,
		}
		_cond = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _c = a[1]

				mml.Nop(_context, _c)
				return func() interface{} {
					c = mml.Ref(_c, "ternary")
					if c.(bool) {
						return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_c, "condition"))}).Values), _scoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_c, "consequent"))}).Values), _fieldsIfHas.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values), &mml.List{Values: append([]interface{}{}, "alternative")}, _c)}).Values))}).Values))}).Values)
					} else {
						return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_c, "condition"))}).Values), _scoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_c, "consequent"))}).Values), _fieldsIfHas.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values), &mml.List{Values: append([]interface{}{}, "alternative")}, _c)}).Values))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 2,
		}
		_validateCase = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _c = a[1]

				mml.Nop(_context, _c)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_c, "expression"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_c, "body"))}).Values))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_validateSwitch = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fieldsIfHas.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, &mml.List{Values: append([]interface{}{}, "expression")}, _s)}).Values), _allScoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "cases"))}).Values), _scoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "defaultStatements"))}).Values))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_validateReceive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _r = a[1]

				mml.Nop(_context, _r)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_r, "channel"))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_validateSelect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _allScoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "cases"))}).Values), func() interface{} {
					c = mml.Ref(_s, "hasDefault")
					if c.(bool) {
						return _scoped.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_s, "defaultStatements"))}).Values)
					} else {
						return _emptyResults
					}
				}())}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_rangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _r = a[1]

				mml.Nop(_context, _r)
				var _result interface{}
				mml.Nop(_result)
				c = !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _r)}).Values).(bool)
				if c.(bool) {
					mml.Nop()
					_define.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_r, "symbol"), &mml.List{Values: append([]interface{}{}, 0)})}).Values)
					return _emptyResults
				}
				mml.SetRef(_context, "capturing", true)
				_result = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_r, "expression"))}).Values)
				_define.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_r, "symbol"), mml.Ref(_result, "values"))}).Values)
				mml.SetRef(_context, "capturing", false)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _result)}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_loop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _l = a[1]

				mml.Nop(_context, _l)
				var _c interface{}
				mml.Nop(_c)
				_c = _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _l)}).Values)
					if c.(bool) {
						return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c, mml.Ref(_l, "expression"))}).Values)
					} else {
						return _emptyResults
					}
				}(), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c, mml.Ref(_l, "body"))}).Values))}).Values))}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_definition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _d = a[1]

				mml.Nop(_context, _d)
				var _r interface{}
				mml.Nop(_r)
				c = _definedCurrent.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_d, "symbol"))}).Values)
				if c.(bool) {
					mml.Nop()
					return _resultErrors.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _duplicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_d, "symbol"))}).Values))}).Values)
				}
				mml.SetRef(_context, "capturing", true)
				_r = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_d, "expression"))}).Values)
				mml.SetRef(_context, "capturing", false)
				_define.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_d, "symbol"), mml.Ref(_r, "values"))}).Values)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r)}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_assignment = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _a = a[1]

				mml.Nop(_context, _a)
				var _cr interface{}
				var _er interface{}
				mml.Nop(_cr, _er)
				_cr = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_a, "capture"))}).Values)
				mml.SetRef(_context, "capturing", true)
				_er = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_a, "value"))}).Values)
				mml.SetRef(_context, "capturing", false)
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _cr, _er)}).Values))}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_validateUse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _u = a[1]

				mml.Nop(_context, _u)

				mml.Nop()
				c = mml.BinaryOp(11, mml.Ref(_u, "capture"), "")
				if c.(bool) {
					mml.Nop()
					_defineImport.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_u, "path"), &mml.List{Values: append([]interface{}{}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }())})}).Values)
					return _emptyResults
				}
				c = mml.BinaryOp(11, mml.Ref(_u, "capture"), ".")
				if c.(bool) {
					mml.Nop()
					for _, _name := range mml.Ref(_u, "exportNames").(*mml.List).Values {

						mml.Nop()
						_defineImport.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _name, &mml.List{Values: append([]interface{}{}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }())})}).Values)
					}
					return _emptyResults
				}
				_defineImport.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_u, "capture"), &mml.List{Values: append([]interface{}{}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }())})}).Values)
				return _emptyResults
				return nil
			},
			FixedArgs: 2,
		}
		_statements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _s = a[1]

				mml.Nop(_context, _s)
				var _r interface{}
				mml.Nop(_r)
				_r = _emptyResults
				for _, _si := range _s.(*mml.List).Values {
					var _ri interface{}
					mml.Nop(_ri)
					_ri = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _si)}).Values)
					c = (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _si)}).Values).(bool) && mml.BinaryOp(11, mml.Ref(_si, "type"), "ret").(bool))
					if c.(bool) {
						mml.Nop()
						_r = _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r, _ri)}).Values)
					} else {
						mml.Nop()
						_r = _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r, _resultErrors.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ri, "errors").(*mml.List).Values...)}).Values))}).Values)
					}
				}
				for _, _f := range mml.Ref(_context, "unexpanded").(*mml.List).Values {

					mml.Nop()
					_r = _mergeResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r, _expandFunction.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f)}).Values))}).Values)
				}
				mml.SetRef(_context, "unexpanded", &mml.List{Values: []interface{}{}})
				return _dropValues.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _r)}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _context = a[0]
				var _code = a[1]

				mml.Nop(_context, _code)

				mml.Nop()
				c = !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _code)}).Values).(bool)
				if c.(bool) {
					mml.Nop()
					return _emptyResults
				}
				switch mml.Ref(_code, "type") {
				case "comment":

					mml.Nop()
					return _emptyResults
				case "symbol":

					mml.Nop()
					return _symbol.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "list":

					mml.Nop()
					return _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "entry":

					mml.Nop()
					return _entry.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "struct":

					mml.Nop()
					return _struct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "function":

					mml.Nop()
					return _function.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "range-expression":

					mml.Nop()
					return _rangeExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "indexer":

					mml.Nop()
					return _indexer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "spread":

					mml.Nop()
					return _spread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "function-application":

					mml.Nop()
					return _application.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "unary":

					mml.Nop()
					return _unary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "binary":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "cond":

					mml.Nop()
					return _cond.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "switch-case":

					mml.Nop()
					return _validateCase.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "switch-statement":

					mml.Nop()
					return _validateSwitch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "send":

					mml.Nop()
					return _validateSend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "receive":

					mml.Nop()
					return _validateReceive.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "go":

					mml.Nop()
					return _validateGo.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "defer":

					mml.Nop()
					return _validateDefer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "select-case":

					mml.Nop()
					return _validateCase.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "select":

					mml.Nop()
					return _validateSelect.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "range-over":

					mml.Nop()
					return _rangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "loop":

					mml.Nop()
					return _loop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "definition":

					mml.Nop()
					return _definition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "definition-list":

					mml.Nop()
					return _definitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "assign":

					mml.Nop()
					return _assignment.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "assign-list":

					mml.Nop()
					return _assignments.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "ret":

					mml.Nop()
					return _ret.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "control-statement":

					mml.Nop()
					return _emptyResults
				case "use":

					mml.Nop()
					return _validateUse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				case "use-list":

					mml.Nop()
					return _useList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _code)}).Values)
				default:

					mml.Nop()
					return _statements.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, mml.Ref(_code, "statements"))}).Values)
				}
				return nil
			},
			FixedArgs: 2,
		}
		_validate = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]

				mml.Nop(_code)
				var _context interface{}
				var _result interface{}
				mml.Nop(_context, _result)
				_context = _newContext.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				for _, _b := range _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_mmlcode, "builtin"))}).Values).(*mml.List).Values {

					mml.Nop()
					_define.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context, _b, &mml.List{Values: []interface{}{}})}).Values)
				}
				_result = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _extend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _importContext.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _context)}).Values))}).Values), _code)}).Values)
				return mml.Ref(_result, "errors")
				return nil
			},
			FixedArgs: 1,
		}
		exports["validate"] = _validate
		return exports
	})
	modulePath = "snippets.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _head interface{}
		var _initHead interface{}
		var _initFooter interface{}
		var _moduleHead interface{}
		var _moduleFooter interface{}
		var _mainHead interface{}
		var _mainFooter interface{}
		mml.Nop(_head, _initHead, _initFooter, _moduleHead, _moduleFooter, _mainHead, _mainFooter)
		_head = "// Generated code\npackage main\n\nimport \"github.com/aryszka/mml\"\n"
		exports["head"] = _head
		_initHead = "\nfunc init() {\n\tvar modulePath string\n"
		exports["initHead"] = _initHead
		_initFooter = "\n}\n"
		exports["initFooter"] = _initFooter
		_moduleHead = "\n\tmml.Modules.Set(modulePath, func() map[string]interface{} {\n\t\texports := make(map[string]interface{})\n\n\t\tvar c interface{}\n\t\tmml.Nop(c)\n"
		exports["moduleHead"] = _moduleHead
		_moduleFooter = "\n\t\treturn exports\n\t})\n"
		exports["moduleFooter"] = _moduleFooter
		_mainHead = "\nfunc main() {\n\tmml.Modules.Use(\""
		exports["mainHead"] = _mainHead
		_mainFooter = "\")\n}\n"
		exports["mainFooter"] = _mainFooter
		return exports
	})
	modulePath = "compile.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _notEmpty interface{}
		var _compileInt interface{}
		var _compileFloat interface{}
		var _compileBool interface{}
		var _getScope interface{}
		var _comment interface{}
		var _compileString interface{}
		var _symbol interface{}
		var _cond interface{}
		var _spreadList interface{}
		var _compileCase interface{}
		var _compileSend interface{}
		var _compileReceive interface{}
		var _compileGo interface{}
		var _definitions interface{}
		var _assigns interface{}
		var _ret interface{}
		var _control interface{}
		var _useList interface{}
		var _list interface{}
		var _entry interface{}
		var _expressionKey interface{}
		var _struct interface{}
		var _paramList interface{}
		var _function interface{}
		var _indexer interface{}
		var _application interface{}
		var _unary interface{}
		var _binary interface{}
		var _ternary interface{}
		var _compileIf interface{}
		var _compileSwitch interface{}
		var _compileSelect interface{}
		var _compileDefer interface{}
		var _rangeOver interface{}
		var _loop interface{}
		var _definition interface{}
		var _assign interface{}
		var _statements interface{}
		var _compileUse interface{}
		var _do interface{}
		var _errors interface{}
		var _code interface{}
		var _strings interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _onlyErr interface{}
		var _passErr interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _channel interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _is interface{}
		mml.Nop(_notEmpty, _compileInt, _compileFloat, _compileBool, _getScope, _comment, _compileString, _symbol, _cond, _spreadList, _compileCase, _compileSend, _compileReceive, _compileGo, _definitions, _assigns, _ret, _control, _useList, _list, _entry, _expressionKey, _struct, _paramList, _function, _indexer, _application, _unary, _binary, _ternary, _compileIf, _compileSwitch, _compileSelect, _compileDefer, _rangeOver, _loop, _definition, _assign, _statements, _compileUse, _do, _errors, _code, _strings, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _bind, _identity, _eq, _any, _function, _channel, _type, _listOf, _structOf, _range, _or, _and, _predicate, _is)
		var __lang = mml.Modules.Use("lang.mml")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_onlyErr = __lang.Values["onlyErr"]
		_passErr = __lang.Values["passErr"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_predicate = __lang.Values["predicate"]
		_is = __lang.Values["is"]
		_errors = mml.Modules.Use("errors.mml")
		_code = mml.Modules.Use("code.mml")
		_strings = mml.Modules.Use("strings.mml")
		_notEmpty = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return mml.BinaryOp(12, _s, "")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values)
			},
			FixedArgs: 1,
		}
		_compileInt = _string
		_compileFloat = _string
		_compileBool = _string
		_getScope = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statements interface{}
				_statements = &mml.List{a[0:]}

				mml.Nop(_statements)
				var _defs interface{}
				var _uses interface{}
				var _inlineUses interface{}
				var _namedUses interface{}
				var _unnamedUses interface{}
				mml.Nop(_defs, _uses, _inlineUses, _namedUses, _unnamedUses)
				_defs = mml.Ref(_code, "flattenedStatements").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definition", "definition-list", "definitions", _statements)}).Values)
				_uses = mml.Ref(_code, "flattenedStatements").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "use", "use-list", "uses", _statements)}).Values)
				_inlineUses = _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.Ref(_u, "exportNames")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "exportNames")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.BinaryOp(11, mml.Ref(_u, "capture"), ".")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values))}).Values))}).Values))}).Values)
				_namedUses = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return (mml.BinaryOp(12, mml.Ref(_u, "capture"), ".").(bool) && mml.BinaryOp(12, mml.Ref(_u, "capture"), "").(bool))
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values)
				_unnamedUses = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.BinaryOp(11, mml.Ref(_u, "capture"), "")
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values)
				return _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]

						mml.Nop(_d)
						return mml.Ref(_d, "symbol")
					},
					FixedArgs: 1,
				}, _defs)}).Values), _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.Ref(_u, "capture")
					},
					FixedArgs: 1,
				}, _namedUses)}).Values), _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.Ref(_u, "path")
					},
					FixedArgs: 1,
				}, _unnamedUses)}).Values), _inlineUses)})}).Values)
				return nil
			},
			FixedArgs: 0,
		}
		_comment = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var __ = a[0]

				mml.Nop(__)
				return ""
			},
			FixedArgs: 1,
		}
		_compileString = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\"%s\"", mml.Ref(_strings, "escape").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_symbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s", mml.Ref(_s, "name"))}).Values)
			},
			FixedArgs: 1,
		}
		_cond = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return func() interface{} {
					c = mml.Ref(_c, "ternary")
					if c.(bool) {
						return _ternary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
					} else {
						return _compileIf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_spreadList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s.(*mml.List).Values...", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_compileCase = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "case %s:\n%s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "expression"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "body"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_compileSend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s <- %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "channel"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_compileReceive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "<- %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "channel"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_compileGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _g = a[0]

				mml.Nop(_g)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "go %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_g, "application"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_definitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, "definitions"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_assigns = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, "assignments"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_ret = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "return %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_control = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_c, "control"), mml.Ref(_code, "breakControl"))
					if c.(bool) {
						return "break"
					} else {
						return "continue"
					}
				}()
			},
			FixedArgs: 1,
		}
		_useList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]

				mml.Nop(_u)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "uses"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_list = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				var _isSpread interface{}
				var _selectSpread interface{}
				var _groupSpread interface{}
				var _appendSimples interface{}
				var _appendSpread interface{}
				var _appendSpreads interface{}
				var _appendGroups interface{}
				var _appendGroup interface{}
				mml.Nop(_isSpread, _selectSpread, _groupSpread, _appendSimples, _appendSpread, _appendSpreads, _appendGroups, _appendGroup)
				_isSpread = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return (mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values), 3).(bool) && mml.BinaryOp(11, mml.RefRange(_c, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values), 3), nil), "...").(bool))
					},
					FixedArgs: 1,
				}
				_selectSpread = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return func() interface{} {
							c = _isSpread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
							if c.(bool) {
								return func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["spread"] = _c
									return s
								}()
							} else {
								return _c
							}
						}()
					},
					FixedArgs: 1,
				}
				_groupSpread = _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _item = a[0]
						var _groups = a[1]

						mml.Nop(_item, _groups)
						var _i interface{}
						var _isSpread interface{}
						var _groupIsSpread interface{}
						var _appendNewSimple interface{}
						var _appendNewSpread interface{}
						var _appendSimple interface{}
						var _appendSpread interface{}
						mml.Nop(_i, _isSpread, _groupIsSpread, _appendNewSimple, _appendNewSpread, _appendSimple, _appendSpread)
						_i = mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _groups)}).Values), 1)
						_isSpread = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "spread", _item)}).Values)
						_groupIsSpread = (mml.BinaryOp(16, _i, 0).(bool) && _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "spread", mml.Ref(_groups, _i))}).Values).(bool))
						_appendNewSimple = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return &mml.List{Values: append(append([]interface{}{}, _groups.(*mml.List).Values...), func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["simple"] = &mml.List{Values: append([]interface{}{}, _item)}
									return s
								}())}
							},
							FixedArgs: 0,
						}
						_appendNewSpread = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return &mml.List{Values: append(append([]interface{}{}, _groups.(*mml.List).Values...), func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["spread"] = &mml.List{Values: append([]interface{}{}, mml.Ref(_item, "spread"))}
									return s
								}())}
							},
							FixedArgs: 0,
						}
						_appendSimple = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return &mml.List{Values: append(append([]interface{}{}, mml.RefRange(_groups, nil, _i).(*mml.List).Values...), func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["simple"] = &mml.List{Values: append(append([]interface{}{}, mml.Ref(mml.Ref(_groups, _i), "simple").(*mml.List).Values...), _item)}
									return s
								}())}
							},
							FixedArgs: 0,
						}
						_appendSpread = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return &mml.List{Values: append(append([]interface{}{}, mml.RefRange(_groups, nil, _i).(*mml.List).Values...), func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["spread"] = &mml.List{Values: append(append([]interface{}{}, mml.Ref(mml.Ref(_groups, _i), "spread").(*mml.List).Values...), mml.Ref(_item, "spread"))}
									return s
								}())}
							},
							FixedArgs: 0,
						}
						switch {
						case ((mml.BinaryOp(13, _i, 0).(bool) || _groupIsSpread.(bool)) && !_isSpread.(bool)):

							mml.Nop()
							return _appendNewSimple.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
						case ((mml.BinaryOp(13, _i, 0).(bool) || !_groupIsSpread.(bool)) && _isSpread.(bool)):

							mml.Nop()
							return _appendNewSpread.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
						case (!_groupIsSpread.(bool) && !_isSpread.(bool)):

							mml.Nop()
							return _appendSimple.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
						case (_groupIsSpread.(bool) && _isSpread.(bool)):

							mml.Nop()
							return _appendSpread.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
						}
						return nil
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}})}).Values)
				_appendSimples = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _group = a[1]

						mml.Nop(_code, _group)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "append(%s, %s)", _code, _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ", ", _group)}).Values))}).Values)
					},
					FixedArgs: 2,
				}
				_appendSpread = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _item = a[0]
						var _code = a[1]

						mml.Nop(_item, _code)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "append(%s, %s)", _code, _item)}).Values)
					},
					FixedArgs: 2,
				}
				_appendSpreads = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _group = a[1]

						mml.Nop(_code, _group)
						return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _appendSpread, _code, _group)}).Values)
					},
					FixedArgs: 2,
				}
				_appendGroups = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _groups = a[0]

						mml.Nop(_groups)
						return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _appendGroup, "[]interface{}{}", _groups)}).Values)
					},
					FixedArgs: 1,
				}
				_appendGroup = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _group = a[0]
						var _code = a[1]

						mml.Nop(_group, _code)
						return func() interface{} {
							c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "spread", _group)}).Values)
							if c.(bool) {
								return _appendSpreads.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code, mml.Ref(_group, "spread"))}).Values)
							} else {
								return _appendSimples.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code, mml.Ref(_group, "simple"))}).Values)
							}
						}()
					},
					FixedArgs: 2,
				}
				return (&mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "&mml.List{Values: %s}", _c)}).Values)
					},
					FixedArgs: 1,
				}).Call((&mml.List{Values: append([]interface{}{}, _appendGroups.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _groupSpread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _selectSpread)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, "values"))}).Values))}).Values))}).Values))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_entry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0]

				mml.Nop(_e)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\"%s\":%s", func() interface{} {
					c = (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(_e, "key"))}).Values).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_e, "key"), "type"), "symbol").(bool))
					if c.(bool) {
						return mml.Ref(mml.Ref(_e, "key"), "name")
					} else {
						return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "key"))}).Values)
					}
				}(), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_expressionKey = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _k = a[0]

				mml.Nop(_k)
				return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_k, "value"))}).Values)
			},
			FixedArgs: 1,
		}
		_struct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				var _entry interface{}
				var _entries interface{}
				mml.Nop(_entry, _entries)
				_entry = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _e = a[0]

						mml.Nop(_e)
						var _v interface{}
						mml.Nop(_v)
						_v = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "value"))}).Values)
						c = mml.BinaryOp(11, mml.Ref(_e, "type"), "spread")
						if c.(bool) {
							var _var interface{}
							var _assign interface{}
							mml.Nop(_var, _assign)
							_var = _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "sp := %s.(*mml.Struct);", _v)}).Values)
							_assign = "for k, v := range sp.Values { s.Values[k] = v };"
							return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func() { %s; %s }();\n", _var, _assign)}).Values)
						}
						c = _isString.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "key"))}).Values)
						if c.(bool) {
							mml.Nop()
							return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "s.Values[\"%s\"] = %s;", mml.Ref(_e, "key"), _v)}).Values)
						}
						c = mml.BinaryOp(11, mml.Ref(mml.Ref(_e, "key"), "type"), "symbol")
						if c.(bool) {
							mml.Nop()
							return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "s.Values[\"%s\"] = %s;", mml.Ref(mml.Ref(_e, "key"), "name"), _v)}).Values)
						}
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "s.Values[%s.(string)] = %s;", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "key"))}).Values), _v)}).Values)
						return nil
					},
					FixedArgs: 1,
				}
				_entries = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _entry)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "entries"))}).Values)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; %s; return s }()", _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _entries)}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_paramList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _params = a[0]
				var _collectParam = a[1]

				mml.Nop(_params, _collectParam)
				var _p interface{}
				mml.Nop(_p)
				_p = &mml.List{Values: []interface{}{}}
				for _i := 0; _i < _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values).(int); _i++ {

					mml.Nop()
					_p = &mml.List{Values: append(append([]interface{}{}, _p.(*mml.List).Values...), _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var _%s = a[%d]", mml.Ref(_params, _i), _i)}).Values))}
				}
				c = mml.BinaryOp(12, _collectParam, "")
				if c.(bool) {
					mml.Nop()
					_p = &mml.List{Values: append(append([]interface{}{}, _p.(*mml.List).Values...), _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var _%s interface{}; _%s = &mml.List{a[%d:]}", _collectParam, _collectParam, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values))}).Values))}
				}
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n", _p)}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_function = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				var _scope interface{}
				var _paramNames interface{}
				mml.Nop(_scope, _paramNames)
				_scope = _getScope.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "statement"))}).Values)
				_paramNames = func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_f, "collectParam"), "")
					if c.(bool) {
						return mml.Ref(_f, "params")
					} else {
						return &mml.List{Values: append(append([]interface{}{}, mml.Ref(_f, "params").(*mml.List).Values...), mml.Ref(_f, "collectParam"))}
					}
				}()
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(_f, "statement"))}).Values).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_f, "statement"), "type"), "statement-list").(bool))
					if c.(bool) {
						return "&mml.Function{\n\t\t\tF: func(a []interface{}) interface{} {\n\t\t\t\tvar c interface{}\n\t\t\t\tmml.Nop(c)\n\t\t\t\t%s;\n\t\t\t\t%s;\n\t\t\t\tmml.Nop(%s);\n\t\t\t\t%s;\n\t\t\t\treturn nil\n\t\t\t},\n\t\t\tFixedArgs: %d,\n\t\t}"
					} else {
						return "&mml.Function{\n\t\t\tF: func(a []interface{}) interface{} {\n\t\t\t\tvar c interface{}\n\t\t\t\tmml.Nop(c)\n\t\t\t\t%s;\n\t\t\t\t%s;\n\t\t\t\tmml.Nop(%s);\n\t\t\t\treturn %s\n\t\t\t},\n\t\t\tFixedArgs: %d,\n\t\t}"
					}
				}(), _paramList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "params"), mml.Ref(_f, "collectParam"))}).Values), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var _%s interface{}", _s)}).Values)
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uniq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _eq)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _scope)}).Values))}).Values))}).Values), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ", ", _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_strings, "formatOne").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s")}).Values), &mml.List{Values: append(append([]interface{}{}, _scope.(*mml.List).Values...), _paramNames.(*mml.List).Values...)})}).Values))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "statement"))}).Values), _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "params"))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_indexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0]

				mml.Nop(_i)
				return func() interface{} {
					c = (!_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(_i, "index"))}).Values).(bool) || mml.BinaryOp(12, mml.Ref(mml.Ref(_i, "index"), "type"), "range-expression").(bool))
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.Ref(%s, %s)", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "expression"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "index"))}).Values))}).Values)
					} else {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.RefRange(%s, %s, %s)", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "expression"))}).Values), func() interface{} {
							c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "from", mml.Ref(_i, "index"))}).Values)
							if c.(bool) {
								return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "from"))}).Values)
							} else {
								return "nil"
							}
						}(), func() interface{} {
							c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "to", mml.Ref(_i, "index"))}).Values)
							if c.(bool) {
								return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "to"))}).Values)
							} else {
								return "nil"
							}
						}())}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_application = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0]

				mml.Nop(_a)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(_a, "function"))}).Values).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_a, "function"), "type"), "function").(bool))
					if c.(bool) {
						return "(%s).Call((%s).Values)"
					} else {
						return "%s.(*mml.Function).Call((%s).Values)"
					}
				}(), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "function"))}).Values), _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["values"] = mml.Ref(_a, "args")
					return s
				}())}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_unary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]

				mml.Nop(_u)
				return func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_u, "op"), mml.Ref(_code, "logicalNot"))
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							c = _isBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "arg"))}).Values)
							if c.(bool) {
								return "!%s"
							} else {
								return "!%s.(bool)"
							}
						}(), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "arg"))}).Values))}).Values)
					} else {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.UnaryOp(%d, %s)", mml.Ref(_u, "op"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "arg"))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_binary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _b = a[0]

				mml.Nop(_b)
				var _isBoolOp interface{}
				var _left interface{}
				var _right interface{}
				var _op interface{}
				mml.Nop(_isBoolOp, _left, _right, _op)
				c = (mml.BinaryOp(12, mml.Ref(_b, "op"), mml.Ref(_code, "logicalAnd")).(bool) && mml.BinaryOp(12, mml.Ref(_b, "op"), mml.Ref(_code, "logicalOr")).(bool))
				if c.(bool) {
					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.BinaryOp(%s, %s, %s)", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "op"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values))}).Values)
				}
				_isBoolOp = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return ((_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", _c)}).Values).(bool) && (mml.BinaryOp(11, mml.Ref(_c, "type"), "unary").(bool) && mml.BinaryOp(11, mml.Ref(_c, "op"), mml.Ref(_code, "logicalNot")).(bool))) || (mml.BinaryOp(11, mml.Ref(_c, "type"), "binary").(bool) && (mml.BinaryOp(11, mml.Ref(_c, "op"), mml.Ref(_code, "logicalAnd")).(bool) || mml.BinaryOp(11, mml.Ref(_c, "op"), mml.Ref(_code, "logicalOr")).(bool))))
					},
					FixedArgs: 1,
				}
				_left = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values)
				_right = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values)
				c = (!_isBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values).(bool) && !_isBoolOp.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values).(bool))
				if c.(bool) {
					mml.Nop()
					_left = mml.BinaryOp(9, _left, ".(bool)")
				}
				c = (!_isBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values).(bool) && !_isBoolOp.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values).(bool))
				if c.(bool) {
					mml.Nop()
					_right = mml.BinaryOp(9, _right, ".(bool)")
				}
				_op = "&&"
				c = mml.BinaryOp(11, mml.Ref(_b, "op"), mml.Ref(_code, "logicalOr"))
				if c.(bool) {
					mml.Nop()
					_op = "||"
				}
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "(%s %s %s)", _left, _op, _right)}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_ternary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func () interface{} { c = %s; if c.(bool) { return %s } else { return %s } }()", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "condition"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "consequent"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "alternative"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_compileIf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "alternative", _c)}).Values)
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "c = %s; if c.(bool) { %s } else { %s }", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "condition"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "consequent"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "alternative"))}).Values))}).Values)
					} else {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "c = %s; if c.(bool) { %s }", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "condition"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "consequent"))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_compileSwitch = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				var _hasDefault interface{}
				var _cases interface{}
				var _def interface{}
				var _defaultCode interface{}
				mml.Nop(_hasDefault, _cases, _def, _defaultCode)
				_hasDefault = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_s, "defaultStatements"), "statements"))}).Values), 0)
				_cases = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "cases"))}).Values)
				_def = func() interface{} {
					c = _hasDefault
					if c.(bool) {
						return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "defaultStatements"))}).Values)
					} else {
						return ""
					}
				}()
				_defaultCode = func() interface{} {
					c = _hasDefault
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "default:\n%s", _def)}).Values)
					} else {
						return ""
					}
				}()
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "switch %s {\n%s\n}", func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _s)}).Values)
					if c.(bool) {
						return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "expression"))}).Values)
					} else {
						return ""
					}
				}(), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = _hasDefault
					if c.(bool) {
						return &mml.List{Values: append(append([]interface{}{}, _cases.(*mml.List).Values...), _defaultCode)}
					} else {
						return _cases
					}
				}())}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_compileSelect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return (&mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return mml.Ref(_strings, "formatOne").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func() interface{} {\nselect {\n%s\n} }()")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							c = mml.Ref(_s, "hasDefault")
							if c.(bool) {
								return &mml.List{Values: append(append([]interface{}{}, _c.(*mml.List).Values...), mml.Ref(_strings, "formatOne").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "default:\n%s")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "defaultStatements"))}).Values))}).Values))}
							} else {
								return _c
							}
						}())}).Values))}).Values)
					},
					FixedArgs: 1,
				}).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "cases"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_compileDefer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0]

				mml.Nop(_d)

				mml.Nop()
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(mml.Ref(_d, "application"), "function"))}).Values).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(mml.Ref(_d, "application"), "function"), "type"), "function").(bool))
					if c.(bool) {
						return "c = (%s); defer c.Call((%s).Values)"
					} else {
						return "defer %s.(*mml.Function).Call((%s).Values)"
					}
				}(), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_d, "application"), "function"))}).Values), _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["values"] = mml.Ref(mml.Ref(_d, "application"), "args")
					return s
				}())}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_rangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				var _infiniteCounter interface{}
				var _withRangeExpression interface{}
				var _listStyleRange interface{}
				mml.Nop(_infiniteCounter, _withRangeExpression, _listStyleRange)
				_infiniteCounter = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s := 0; true; _%s++", mml.Ref(_r, "symbol"), mml.Ref(_r, "symbol"))}).Values)
					},
					FixedArgs: 0,
				}
				_withRangeExpression = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s := %s; %s; _%s++", mml.Ref(_r, "symbol"), func() interface{} {
							c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "from", mml.Ref(_r, "expression"))}).Values)
							if c.(bool) {
								return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_r, "expression"), "from"))}).Values)
							} else {
								return "0"
							}
						}(), func() interface{} {
							c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "to", mml.Ref(_r, "expression"))}).Values)
							if c.(bool) {
								return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s < %s.(int)", mml.Ref(_r, "symbol"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_r, "expression"), "to"))}).Values))}).Values)
							} else {
								return "true"
							}
						}(), mml.Ref(_r, "symbol"))}).Values)
					},
					FixedArgs: 0,
				}
				_listStyleRange = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_, _%s := range %s.(*mml.List).Values", mml.Ref(_r, "symbol"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "expression"))}).Values))}).Values)
					},
					FixedArgs: 0,
				}
				switch {
				case !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _r)}).Values).(bool):

					mml.Nop()
					return _infiniteCounter.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				case (_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "type", mml.Ref(_r, "expression"))}).Values).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_r, "expression"), "type"), "range-expression").(bool)):

					mml.Nop()
					return _withRangeExpression.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				default:

					mml.Nop()
					return _listStyleRange.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_loop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "for %s {\n%s\n}", func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _l)}).Values)
					if c.(bool) {
						return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, "expression"))}).Values)
					} else {
						return ""
					}
				}(), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, "body"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_definition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0]

				mml.Nop(_d)
				return func() interface{} {
					c = mml.Ref(_d, "exported")
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = %s; exports[\"%s\"] = _%s", mml.Ref(_d, "symbol"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_d, "expression"))}).Values), mml.Ref(_d, "symbol"), mml.Ref(_d, "symbol"))}).Values)
					} else {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = %s", mml.Ref(_d, "symbol"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_d, "expression"))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_assign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0]

				mml.Nop(_a)
				return func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(mml.Ref(_a, "capture"), "type"), "symbol")
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s = %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "capture"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "value"))}).Values))}).Values)
					} else {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.SetRef(%s, %s, %s)", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_a, "capture"), "expression"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_a, "capture"), "index"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "value"))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_statements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				var _scope interface{}
				var _scopeNames interface{}
				var _statements interface{}
				var _scopeDefs interface{}
				mml.Nop(_scope, _scopeNames, _statements, _scopeDefs)
				_scope = _getScope.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s.(*mml.List).Values...)}).Values)
				_scopeNames = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ", ", _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_strings, "formatOne").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s")}).Values), _scope)}).Values))}).Values)
				_statements = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _notEmpty.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do, _s)}).Values))}).Values))}).Values)
				_scopeDefs = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var _%s interface{}", _s)}).Values)
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uniq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _eq)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _scope)}).Values))}).Values))}).Values)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s;\nmml.Nop(%s);\n%s", _scopeDefs, _scopeNames, _statements)}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_compileUse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]

				mml.Nop(_u)

				mml.Nop()
				switch {
				case mml.BinaryOp(11, mml.Ref(_u, "capture"), "."):
					var _useStatement interface{}
					var _assigns interface{}
					mml.Nop(_useStatement, _assigns)
					_useStatement = _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var __%s = mml.Modules.Use(\"%s.mml\");", mml.Ref(_code, "getModuleName").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "path"))}).Values), mml.Ref(_u, "path"))}).Values)
					_assigns = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _name = a[0]

							mml.Nop(_name)
							return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = __%s.Values[\"%s\"]", _name, mml.Ref(_code, "getModuleName").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "path"))}).Values), _name)}).Values)
						},
						FixedArgs: 1,
					}, mml.Ref(_u, "exportNames"))}).Values))}).Values)
					return _joins.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";", _useStatement, _assigns)}).Values)
				case mml.BinaryOp(12, mml.Ref(_u, "capture"), ""):

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = mml.Modules.Use(\"%s.mml\")", mml.Ref(_u, "capture"), mml.Ref(_u, "path"))}).Values)
				default:

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = mml.Modules.Use(\"%s.mml\")", mml.Ref(_code, "getModuleName").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "path"))}).Values), mml.Ref(_u, "path"))}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]

				mml.Nop(_code)

				mml.Nop()
				switch {
				case _isInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values):

					mml.Nop()
					return _compileInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case _isFloat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values):

					mml.Nop()
					return _compileFloat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case _isString.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values):

					mml.Nop()
					return _compileString.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case _isBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values):

					mml.Nop()
					return _compileBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				}
				switch mml.Ref(_code, "type") {
				case "comment":

					mml.Nop()
					return _comment.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "symbol":

					mml.Nop()
					return _symbol.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "list":

					mml.Nop()
					return _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "expression-key":

					mml.Nop()
					return _expressionKey.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "entry":

					mml.Nop()
					return _entry.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "struct":

					mml.Nop()
					return _struct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "function":

					mml.Nop()
					return _function.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "indexer":

					mml.Nop()
					return _indexer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "spread":

					mml.Nop()
					return _spreadList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "function-application":

					mml.Nop()
					return _application.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "unary":

					mml.Nop()
					return _unary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "binary":

					mml.Nop()
					return _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "cond":

					mml.Nop()
					return _cond.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "switch-case":

					mml.Nop()
					return _compileCase.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "switch-statement":

					mml.Nop()
					return _compileSwitch.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "send":

					mml.Nop()
					return _compileSend.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "receive":

					mml.Nop()
					return _compileReceive.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "go":

					mml.Nop()
					return _compileGo.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "defer":

					mml.Nop()
					return _compileDefer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "select-case":

					mml.Nop()
					return _compileCase.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "select":

					mml.Nop()
					return _compileSelect.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "range-over":

					mml.Nop()
					return _rangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "loop":

					mml.Nop()
					return _loop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "definition":

					mml.Nop()
					return _definition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "definition-list":

					mml.Nop()
					return _definitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "assign":

					mml.Nop()
					return _assign.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "assign-list":

					mml.Nop()
					return _assigns.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "ret":

					mml.Nop()
					return _ret.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "control-statement":

					mml.Nop()
					return _control.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "use":

					mml.Nop()
					return _compileUse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "use-list":

					mml.Nop()
					return _useList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				default:

					mml.Nop()
					return _statements.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "statements"))}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		exports["do"] = _do
		return exports
	})

}

func main() {
	mml.Modules.Use("main.mml")
}
