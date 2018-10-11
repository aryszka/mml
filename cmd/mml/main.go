// Generated code

package main

import "github.com/aryszka/mml"

var _args interface{} = mml.Args
var _close interface{} = mml.Close
var _error interface{} = mml.Error
var _format interface{} = mml.Format
var _has interface{} = mml.Has
var _isBool interface{} = mml.IsBool
var _isError interface{} = mml.IsError
var _isFloat interface{} = mml.IsFloat
var _isInt interface{} = mml.IsInt
var _isString interface{} = mml.IsString
var _keys interface{} = mml.Keys
var _len interface{} = mml.Len
var _open interface{} = mml.Open
var _parse interface{} = mml.Parse
var _stderr interface{} = mml.Stderr
var _stdin interface{} = mml.Stdin
var _stdout interface{} = mml.Stdout
var _string interface{} = mml.String

func init() {
	var modulePath string
	modulePath = "compile.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
		var _not interface{}
		var _yes interface{}
		var _ifErr interface{}
		var _passErr interface{}
		var _onlyErr interface{}
		var _anyErr interface{}
		var _join interface{}
		var _joins interface{}
		var _joinTwo interface{}
		var _formats interface{}
		var _formatOne interface{}
		var _passErrFormat interface{}
		var _notEmpty interface{}
		var _sort interface{}
		var _counter interface{}
		var _enum interface{}
		var _log interface{}
		var _escape interface{}
		var _compileInt interface{}
		var _compileFloat interface{}
		var _compileBool interface{}
		var _getFlattenedStatements interface{}
		var _getScope interface{}
		var _mapCompile interface{}
		var _mapCompileJoin interface{}
		var _compileComment interface{}
		var _compileString interface{}
		var _compileSymbol interface{}
		var _compileEntries interface{}
		var _compileStructure interface{}
		var _compileCond interface{}
		var _compileSpread interface{}
		var _compileCase interface{}
		var _compileSend interface{}
		var _compileReceive interface{}
		var _compileGo interface{}
		var _compileDefer interface{}
		var _compileDefinitions interface{}
		var _compileAssign interface{}
		var _compileAssigns interface{}
		var _compileRet interface{}
		var _compileControl interface{}
		var _compileList interface{}
		var _compileEntry interface{}
		var _compileParamList interface{}
		var _compileFunction interface{}
		var _compileRangeExpression interface{}
		var _compileIndexer interface{}
		var _compileApplication interface{}
		var _unaryOp interface{}
		var _binaryNot interface{}
		var _plus interface{}
		var _minus interface{}
		var _logicalNot interface{}
		var _compileUnary interface{}
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
		var _eq interface{}
		var _notEq interface{}
		var _less interface{}
		var _lessOrEq interface{}
		var _greater interface{}
		var _greaterOrEq interface{}
		var _logicalAnd interface{}
		var _logicalOr interface{}
		var _compileBinary interface{}
		var _compileTernary interface{}
		var _compileIf interface{}
		var _compileSwitch interface{}
		var _compileSelect interface{}
		var _compileRangeOver interface{}
		var _compileLoop interface{}
		var _compileDefinition interface{}
		var _controlStatement interface{}
		var _breakControl interface{}
		var _continueControl interface{}
		var _compileStatements interface{}
		var _compileModule interface{}
		var _getModuleName interface{}
		var _compileUse interface{}
		var _compileUseList interface{}
		var _compile interface{}
		var _builtin interface{}
		var _builtins interface{}
		var _parseFile interface{}
		var _findExportNames interface{}
		var _parseModules interface{}
		var _compileModuleCode interface{}
		var _compileModules interface{}
		var _modules interface{}
		var _list interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _firstOr interface{}
		var _contains interface{}
		mml.Nop(_not, _yes, _ifErr, _passErr, _onlyErr, _anyErr, _join, _joins, _joinTwo, _formats, _formatOne, _passErrFormat, _notEmpty, _sort, _counter, _enum, _log, _escape, _compileInt, _compileFloat, _compileBool, _getFlattenedStatements, _getScope, _mapCompile, _mapCompileJoin, _compileComment, _compileString, _compileSymbol, _compileEntries, _compileStructure, _compileCond, _compileSpread, _compileCase, _compileSend, _compileReceive, _compileGo, _compileDefer, _compileDefinitions, _compileAssign, _compileAssigns, _compileRet, _compileControl, _compileList, _compileEntry, _compileParamList, _compileFunction, _compileRangeExpression, _compileIndexer, _compileApplication, _unaryOp, _binaryNot, _plus, _minus, _logicalNot, _compileUnary, _binaryOp, _binaryAnd, _binaryOr, _xor, _andNot, _lshift, _rshift, _mul, _div, _mod, _add, _sub, _eq, _notEq, _less, _lessOrEq, _greater, _greaterOrEq, _logicalAnd, _logicalOr, _compileBinary, _compileTernary, _compileIf, _compileSwitch, _compileSelect, _compileRangeOver, _compileLoop, _compileDefinition, _controlStatement, _breakControl, _continueControl, _compileStatements, _compileModule, _getModuleName, _compileUse, _compileUseList, _compile, _builtin, _builtins, _parseFile, _findExportNames, _parseModules, _compileModuleCode, _compileModules, _modules, _list, _fold, _foldr, _map, _filter, _firstOr, _contains)
		var __list = mml.Modules.Use("list.mml")
		_fold = __list["fold"]
		_foldr = __list["foldr"]
		_map = __list["map"]
		_filter = __list["filter"]
		_firstOr = __list["firstOr"]
		_contains = __list["contains"]
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
							c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a))))
							if c.(bool) {
								return _f.(*mml.Function).Call(append([]interface{}{}, _a))
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
		_passErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}
		_onlyErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}
		_anyErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]

						mml.Nop(_c, _r)
						return func() interface{} {
							c = _isError.(*mml.Function).Call(append([]interface{}{}, _r))
							if c.(bool) {
								return _r
							} else {
								return func() interface{} {
									c = _isError.(*mml.Function).Call(append([]interface{}{}, _c))
									if c.(bool) {
										return _c
									} else {
										return append(append([]interface{}{}, _r.([]interface{})...), _c)
									}
								}()
							}
						}()
					},
					FixedArgs: 2,
				}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}
		_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _s = a[1]

				mml.Nop(_j, _s)
				return func() interface{} {
					c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2)
					if c.(bool) {
						return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s))
					} else {
						return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil))))
					}
				}()
			},
			FixedArgs: 2,
		}
		_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _s = a[1:]

				mml.Nop(_j, _s)
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}
		_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _left = a[1]
				var _right = a[2]

				mml.Nop(_j, _left, _right)
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}
		_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a = a[1:]

				mml.Nop(_f, _a)
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}
		_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a = a[1]

				mml.Nop(_f, _a)
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}
		_passErrFormat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _fmt = a[0]
				var _p = a[1:]

				mml.Nop(_fmt, _p)
				return _passErr.(*mml.Function).Call(append([]interface{}{}, _format.(*mml.Function).Call(append([]interface{}{}, _fmt)))).(*mml.Function).Call(append([]interface{}{}, _anyErr.(*mml.Function).Call(append([]interface{}{}, _p))))
			},
			FixedArgs: 1,
		}
		_notEmpty = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return mml.BinaryOp(12, _s, "")
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _l))
			},
			FixedArgs: 1,
		}
		_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0]
				var _l = a[1]

				mml.Nop(_less, _l)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0)
					if c.(bool) {
						return []interface{}{}
					} else {
						return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _i = a[0]

								mml.Nop(_i)
								return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
							},
							FixedArgs: 1,
						})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...)
					}
				}()
			},
			FixedArgs: 2,
		}
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
		_enum = _counter
		_log = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0:]

				mml.Nop(_a)

				mml.Nop()
				_stderr.(*mml.Function).Call(append([]interface{}{}, _join.(*mml.Function).Call(append([]interface{}{}, " ")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _string)).(*mml.Function).Call(append([]interface{}{}, _a))))))
				_stderr.(*mml.Function).Call(append([]interface{}{}, "\n"))
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 0)
					if c.(bool) {
						return ""
					} else {
						return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 1))
					}
				}()
				return nil
			},
			FixedArgs: 0,
		}
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
				return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))))
				return nil
			},
			FixedArgs: 1,
		}
		_compileInt = _string
		_compileFloat = _string
		_compileBool = _string
		_getFlattenedStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _itemType = a[0]
				var _listType = a[1]
				var _listProp = a[2]
				var _statements = a[3]

				mml.Nop(_itemType, _listType, _listProp, _statements)
				var _toList interface{}
				mml.Nop(_toList)
				_toList = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return func() interface{} {
							c = mml.BinaryOp(11, mml.Ref(_s, "type"), _itemType)
							if c.(bool) {
								return append([]interface{}{}, _s)
							} else {
								return mml.Ref(_s, _listProp)
							}
						}()
					},
					FixedArgs: 1,
				}
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]
						var _result = a[1]

						mml.Nop(_s, _result)
						return append(append([]interface{}{}, _result.([]interface{})...), _toList.(*mml.Function).Call(append([]interface{}{}, _s)).([]interface{})...)
					},
					FixedArgs: 2,
				}, []interface{}{})).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return (_has.(*mml.Function).Call(append([]interface{}{}, "type", _s)).(bool) && _contains.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "type"), append([]interface{}{}, _itemType, _listType))).(bool))
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _statements))))
				return nil
			},
			FixedArgs: 4,
		}
		_getScope = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statements = a[0:]

				mml.Nop(_statements)
				var _defs interface{}
				var _uses interface{}
				var _inlineUses interface{}
				mml.Nop(_defs, _uses, _inlineUses)
				_defs = _getFlattenedStatements.(*mml.Function).Call(append([]interface{}{}, "definition", "definition-list", "definitions", _statements))
				_uses = _getFlattenedStatements.(*mml.Function).Call(append([]interface{}{}, "use", "use-list", "uses", _statements))
				_inlineUses = _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _current = a[0]
						var _result = a[1]

						mml.Nop(_current, _result)
						return append(append([]interface{}{}, _result.([]interface{})...), _current.([]interface{})...)
					},
					FixedArgs: 2,
				}, []interface{}{})).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.Ref(_u, "exportNames")
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _has.(*mml.Function).Call(append([]interface{}{}, "exportNames")))).(*mml.Function).Call(append([]interface{}{}, _uses))))))
				return append(append(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]

						mml.Nop(_d)
						return mml.Ref(_d, "symbol")
					},
					FixedArgs: 1,
				}, _defs)).([]interface{})...), _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.Ref(_u, "path")
					},
					FixedArgs: 1,
				}, _uses)).([]interface{})...), _inlineUses.([]interface{})...)
				return nil
			},
			FixedArgs: 0,
		}
		_mapCompile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _anyErr.(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, _l))))
			},
			FixedArgs: 1,
		}
		_mapCompileJoin = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0]
				var _l = a[1]

				mml.Nop(_j, _l)
				return _passErr.(*mml.Function).Call(append([]interface{}{}, _join.(*mml.Function).Call(append([]interface{}{}, _j)))).(*mml.Function).Call(append([]interface{}{}, _mapCompile.(*mml.Function).Call(append([]interface{}{}, _l))))
			},
			FixedArgs: 2,
		}
		_compileComment = &mml.Function{
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
				return _formats.(*mml.Function).Call(append([]interface{}{}, "\"%s\"", _escape.(*mml.Function).Call(append([]interface{}{}, _s))))
			},
			FixedArgs: 1,
		}
		_compileSymbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s", mml.Ref(_s, "name")))
			},
			FixedArgs: 1,
		}
		_compileEntries = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0]

				mml.Nop(_e)
				return _mapCompileJoin.(*mml.Function).Call(append([]interface{}{}, ",", _e))
			},
			FixedArgs: 1,
		}
		_compileStructure = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "map[string]interface{}{%s}", _compileEntries.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "entries")))))
			},
			FixedArgs: 1,
		}
		_compileCond = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return func() interface{} {
					c = mml.Ref(_c, "ternary")
					if c.(bool) {
						return _compileTernary.(*mml.Function).Call(append([]interface{}{}, _c))
					} else {
						return _compileIf.(*mml.Function).Call(append([]interface{}{}, _c))
					}
				}()
			},
			FixedArgs: 1,
		}
		_compileSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "%s.([]interface{})...", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "value")))))
			},
			FixedArgs: 1,
		}
		_compileCase = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "case %s:\n%s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "expression"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "body")))))
			},
			FixedArgs: 1,
		}
		_compileSend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "%s <- %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "channel"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "value")))))
			},
			FixedArgs: 1,
		}
		_compileReceive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "<- %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_r, "channel")))))
			},
			FixedArgs: 1,
		}
		_compileGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _g = a[0]

				mml.Nop(_g)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "go %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_g, "application")))))
			},
			FixedArgs: 1,
		}
		_compileDefer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0]

				mml.Nop(_d)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "defer %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_d, "application")))))
			},
			FixedArgs: 1,
		}
		_compileDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _mapCompileJoin.(*mml.Function).Call(append([]interface{}{}, ";\n", mml.Ref(_l, "definitions")))
			},
			FixedArgs: 1,
		}
		_compileAssign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0]

				mml.Nop(_a)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "%s = %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "capture"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "value")))))
			},
			FixedArgs: 1,
		}
		_compileAssigns = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _mapCompileJoin.(*mml.Function).Call(append([]interface{}{}, ";\n", mml.Ref(_l, "assignments")))
			},
			FixedArgs: 1,
		}
		_compileRet = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]

				mml.Nop(_r)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "return %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_r, "value")))))
			},
			FixedArgs: 1,
		}
		_compileControl = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_c, "control"), _breakControl)
					if c.(bool) {
						return "break"
					} else {
						return "continue"
					}
				}()
			},
			FixedArgs: 1,
		}
		_compileList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				var _compileValues interface{}
				var _isSpread interface{}
				var _selectSpread interface{}
				var _groupSpread interface{}
				var _appendSimples interface{}
				var _appendSpread interface{}
				var _appendSpreads interface{}
				var _appendGroups interface{}
				var _appendGroup interface{}
				mml.Nop(_compileValues, _isSpread, _selectSpread, _groupSpread, _appendSimples, _appendSpread, _appendSpreads, _appendGroups, _appendGroup)
				_compileValues = _mapCompile
				_isSpread = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return (mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _c)), 3).(bool) && mml.BinaryOp(11, mml.RefRange(_c, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _c)), 3), nil), "...").(bool))
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
							c = _isSpread.(*mml.Function).Call(append([]interface{}{}, _c))
							if c.(bool) {
								return map[string]interface{}{"spread": _c}
							} else {
								return _c
							}
						}()
					},
					FixedArgs: 1,
				}
				_groupSpread = _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
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
						_i = mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _groups)), 1)
						_isSpread = _has.(*mml.Function).Call(append([]interface{}{}, "spread", _item))
						_groupIsSpread = (mml.BinaryOp(16, _i, 0).(bool) && _has.(*mml.Function).Call(append([]interface{}{}, "spread", mml.Ref(_groups, _i))).(bool))
						_appendNewSimple = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return append(append([]interface{}{}, _groups.([]interface{})...), map[string]interface{}{"simple": append([]interface{}{}, _item)})
							},
							FixedArgs: 0,
						}
						_appendNewSpread = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return append(append([]interface{}{}, _groups.([]interface{})...), map[string]interface{}{"spread": append([]interface{}{}, mml.Ref(_item, "spread"))})
							},
							FixedArgs: 0,
						}
						_appendSimple = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return append(append([]interface{}{}, mml.RefRange(_groups, nil, _i).([]interface{})...), map[string]interface{}{"simple": append(append([]interface{}{}, mml.Ref(mml.Ref(_groups, _i), "simple").([]interface{})...), _item)})
							},
							FixedArgs: 0,
						}
						_appendSpread = &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)

								mml.Nop()
								return append(append([]interface{}{}, mml.RefRange(_groups, nil, _i).([]interface{})...), map[string]interface{}{"spread": append(append([]interface{}{}, mml.Ref(mml.Ref(_groups, _i), "spread").([]interface{})...), mml.Ref(_item, "spread"))})
							},
							FixedArgs: 0,
						}
						switch {
						case ((mml.BinaryOp(13, _i, 0).(bool) || _groupIsSpread.(bool)) && !_isSpread.(bool)):

							mml.Nop()
							return _appendNewSimple.(*mml.Function).Call([]interface{}{})
						case ((mml.BinaryOp(13, _i, 0).(bool) || !_groupIsSpread.(bool)) && _isSpread.(bool)):

							mml.Nop()
							return _appendNewSpread.(*mml.Function).Call([]interface{}{})
						case (!_groupIsSpread.(bool) && !_isSpread.(bool)):

							mml.Nop()
							return _appendSimple.(*mml.Function).Call([]interface{}{})
						case (_groupIsSpread.(bool) && _isSpread.(bool)):

							mml.Nop()
							return _appendSpread.(*mml.Function).Call([]interface{}{})
						}
						return nil
					},
					FixedArgs: 2,
				}, []interface{}{}))
				_appendSimples = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _group = a[1]

						mml.Nop(_code, _group)
						return _formats.(*mml.Function).Call(append([]interface{}{}, "append(%s, %s)", _code, _join.(*mml.Function).Call(append([]interface{}{}, ", ", _group))))
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
						return _formats.(*mml.Function).Call(append([]interface{}{}, "append(%s, %s)", _code, _item))
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
						return _fold.(*mml.Function).Call(append([]interface{}{}, _appendSpread, _code, _group))
					},
					FixedArgs: 2,
				}
				_appendGroups = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _groups = a[0]

						mml.Nop(_groups)
						return _fold.(*mml.Function).Call(append([]interface{}{}, _appendGroup, "[]interface{}{}", _groups))
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
							c = _has.(*mml.Function).Call(append([]interface{}{}, "spread", _group))
							if c.(bool) {
								return _appendSpreads.(*mml.Function).Call(append([]interface{}{}, _code, mml.Ref(_group, "spread")))
							} else {
								return _appendSimples.(*mml.Function).Call(append([]interface{}{}, _code, mml.Ref(_group, "simple")))
							}
						}()
					},
					FixedArgs: 2,
				}
				return _passErr.(*mml.Function).Call(append([]interface{}{}, _appendGroups)).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _groupSpread)).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _selectSpread)))).(*mml.Function).Call(append([]interface{}{}, _compileValues.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "values")))))))))
				return nil
			},
			FixedArgs: 1,
		}
		_compileEntry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0]

				mml.Nop(_e)
				var _key interface{}
				var _value interface{}
				mml.Nop(_key, _value)
				_key = func() interface{} {
					c = (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_e, "key"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_e, "key"), "type"), "symbol").(bool))
					if c.(bool) {
						return mml.Ref(mml.Ref(_e, "key"), "name")
					} else {
						return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "key")))
					}
				}()
				_value = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "value")))
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "\"%s\":%s", _key, _value))
				return nil
			},
			FixedArgs: 1,
		}
		_compileParamList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _params = a[0]
				var _collectParam = a[1]

				mml.Nop(_params, _collectParam)
				var _p interface{}
				mml.Nop(_p)
				_p = []interface{}{}
				for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _params)).(int); _i++ {

					mml.Nop()
					_p = append(append([]interface{}{}, _p.([]interface{})...), _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s = a[%d]", mml.Ref(_params, _i), _i)))
				}
				c = mml.BinaryOp(12, _collectParam, "")
				if c.(bool) {
					mml.Nop()
					_p = append(append([]interface{}{}, _p.([]interface{})...), _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s = a[%d:]", _collectParam, _len.(*mml.Function).Call(append([]interface{}{}, _params)))))
				}
				return _join.(*mml.Function).Call(append([]interface{}{}, ";\n", _p))
				return nil
			},
			FixedArgs: 2,
		}
		_compileFunction = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]

				mml.Nop(_f)
				var _multipleStatements interface{}
				var _scope interface{}
				var _paramNames interface{}
				var _scopeNames interface{}
				var _scopeDefs interface{}
				var _fmt interface{}
				var _p interface{}
				var _s interface{}
				mml.Nop(_multipleStatements, _scope, _paramNames, _scopeNames, _scopeDefs, _fmt, _p, _s)
				_multipleStatements = (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_f, "statement"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_f, "statement"), "type"), "statement-list").(bool))
				_scope = _getScope.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "statement")))
				_paramNames = func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_f, "collectParam"), "")
					if c.(bool) {
						return mml.Ref(_f, "params")
					} else {
						return append(append([]interface{}{}, mml.Ref(_f, "params").([]interface{})...), mml.Ref(_f, "collectParam"))
					}
				}()
				_scopeNames = _join.(*mml.Function).Call(append([]interface{}{}, ", ", _map.(*mml.Function).Call(append([]interface{}{}, _formatOne.(*mml.Function).Call(append([]interface{}{}, "_%s")), append(append([]interface{}{}, _scope.([]interface{})...), _paramNames.([]interface{})...)))))
				_scopeDefs = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s interface{}", _s))
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _scope))))
				_fmt = func() interface{} {
					c = _multipleStatements
					if c.(bool) {
						return "&mml.Function{\n\t\t\tF: func(a []interface{}) interface{} {\n\t\t\t\tvar c interface{}\n\t\t\t\tmml.Nop(c)\n\t\t\t\t%s;\n\t\t\t\t%s;\n\t\t\t\tmml.Nop(%s);\n\t\t\t\t%s;\n\t\t\t\treturn nil\n\t\t\t},\n\t\t\tFixedArgs: %d,\n\t\t}"
					} else {
						return "&mml.Function{\n\t\t\tF: func(a []interface{}) interface{} {\n\t\t\t\tvar c interface{}\n\t\t\t\tmml.Nop(c)\n\t\t\t\t%s;\n\t\t\t\t%s;\n\t\t\t\tmml.Nop(%s);\n\t\t\t\treturn %s\n\t\t\t},\n\t\t\tFixedArgs: %d,\n\t\t}"
					}
				}()
				_p = _compileParamList.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "params"), mml.Ref(_f, "collectParam")))
				_s = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "statement")))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _s))
				if c.(bool) {
					mml.Nop()
					return _s
				}
				return _formats.(*mml.Function).Call(append([]interface{}{}, _fmt, _p, _scopeDefs, _scopeNames, _s, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "params")))))
				return nil
			},
			FixedArgs: 1,
		}
		_compileRangeExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0]

				mml.Nop(_e)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "%s:%s", func() interface{} {
					c = _has.(*mml.Function).Call(append([]interface{}{}, "from", _e))
					if c.(bool) {
						return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "from")))
					} else {
						return ""
					}
				}(), func() interface{} {
					c = _has.(*mml.Function).Call(append([]interface{}{}, "to", _e))
					if c.(bool) {
						return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "to")))
					} else {
						return ""
					}
				}()))
			},
			FixedArgs: 1,
		}
		_compileIndexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0]

				mml.Nop(_i)
				var _exp interface{}
				var _from interface{}
				var _to interface{}
				mml.Nop(_exp, _from, _to)
				_exp = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_i, "expression")))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _exp))
				if c.(bool) {
					mml.Nop()
					return _exp
				}
				c = (!_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_i, "index"))).(bool) || mml.BinaryOp(12, mml.Ref(mml.Ref(_i, "index"), "type"), "range-expression").(bool))
				if c.(bool) {
					var _index interface{}
					mml.Nop(_index)
					_index = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_i, "index")))
					c = _isError.(*mml.Function).Call(append([]interface{}{}, _index))
					if c.(bool) {
						mml.Nop()
						return _index
					}
					return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.Ref(%s, %s)", _exp, _index))
				}
				_from = "nil"
				c = _has.(*mml.Function).Call(append([]interface{}{}, "from", mml.Ref(_i, "index")))
				if c.(bool) {
					mml.Nop()
					_from = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "from")))
					c = _isError.(*mml.Function).Call(append([]interface{}{}, _from))
					if c.(bool) {
						mml.Nop()
						return _from
					}
				}
				_to = "nil"
				c = _has.(*mml.Function).Call(append([]interface{}{}, "to", mml.Ref(_i, "index")))
				if c.(bool) {
					mml.Nop()
					_to = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "to")))
					c = _isError.(*mml.Function).Call(append([]interface{}{}, _to))
					if c.(bool) {
						mml.Nop()
						return _to
					}
				}
				return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.RefRange(%s, %s, %s)", _exp, _from, _to))
				return nil
			},
			FixedArgs: 1,
		}
		_compileApplication = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0]

				mml.Nop(_a)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, func() interface{} {
					c = (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_a, "function"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_a, "function"), "type"), "function").(bool))
					if c.(bool) {
						return "(%s).Call(%s)"
					} else {
						return "%s.(*mml.Function).Call(%s)"
					}
				}(), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "function"))), _compileList.(*mml.Function).Call(append([]interface{}{}, map[string]interface{}{"values": mml.Ref(_a, "args")}))))
			},
			FixedArgs: 1,
		}
		_unaryOp = _enum.(*mml.Function).Call([]interface{}{})
		_binaryNot = _unaryOp.(*mml.Function).Call([]interface{}{})
		_plus = _unaryOp.(*mml.Function).Call([]interface{}{})
		_minus = _unaryOp.(*mml.Function).Call([]interface{}{})
		_logicalNot = _unaryOp.(*mml.Function).Call([]interface{}{})
		_compileUnary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]

				mml.Nop(_u)
				var _arg interface{}
				mml.Nop(_arg)
				_arg = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "arg")))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _arg))
				if c.(bool) {
					mml.Nop()
					return _arg
				}
				switch mml.Ref(_u, "op") {
				case _logicalNot:

					mml.Nop()
					c = !_isBool.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "arg"))).(bool)
					if c.(bool) {
						mml.Nop()
						_arg = mml.BinaryOp(9, _arg, ".(bool)")
					}
					return _formats.(*mml.Function).Call(append([]interface{}{}, "!%s", _arg))
				default:

					mml.Nop()
					return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.UnaryOp(%d, %s)", mml.Ref(_u, "op"), _arg))
				}
				return nil
			},
			FixedArgs: 1,
		}
		_binaryOp = _enum.(*mml.Function).Call([]interface{}{})
		_binaryAnd = _binaryOp.(*mml.Function).Call([]interface{}{})
		_binaryOr = _binaryOp.(*mml.Function).Call([]interface{}{})
		_xor = _binaryOp.(*mml.Function).Call([]interface{}{})
		_andNot = _binaryOp.(*mml.Function).Call([]interface{}{})
		_lshift = _binaryOp.(*mml.Function).Call([]interface{}{})
		_rshift = _binaryOp.(*mml.Function).Call([]interface{}{})
		_mul = _binaryOp.(*mml.Function).Call([]interface{}{})
		_div = _binaryOp.(*mml.Function).Call([]interface{}{})
		_mod = _binaryOp.(*mml.Function).Call([]interface{}{})
		_add = _binaryOp.(*mml.Function).Call([]interface{}{})
		_sub = _binaryOp.(*mml.Function).Call([]interface{}{})
		_eq = _binaryOp.(*mml.Function).Call([]interface{}{})
		_notEq = _binaryOp.(*mml.Function).Call([]interface{}{})
		_less = _binaryOp.(*mml.Function).Call([]interface{}{})
		_lessOrEq = _binaryOp.(*mml.Function).Call([]interface{}{})
		_greater = _binaryOp.(*mml.Function).Call([]interface{}{})
		_greaterOrEq = _binaryOp.(*mml.Function).Call([]interface{}{})
		_logicalAnd = _binaryOp.(*mml.Function).Call([]interface{}{})
		_logicalOr = _binaryOp.(*mml.Function).Call([]interface{}{})
		_compileBinary = &mml.Function{
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
				c = (mml.BinaryOp(12, mml.Ref(_b, "op"), _logicalAnd).(bool) && mml.BinaryOp(12, mml.Ref(_b, "op"), _logicalOr).(bool))
				if c.(bool) {
					mml.Nop()
					return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "mml.BinaryOp(%s, %s, %s)", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "op"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right")))))
				}
				_isBoolOp = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]

						mml.Nop(_c)
						return ((_has.(*mml.Function).Call(append([]interface{}{}, "type", _c)).(bool) && (mml.BinaryOp(11, mml.Ref(_c, "type"), "unary").(bool) && mml.BinaryOp(11, mml.Ref(_c, "op"), _logicalNot).(bool))) || (mml.BinaryOp(11, mml.Ref(_c, "type"), "binary").(bool) && (mml.BinaryOp(11, mml.Ref(_c, "op"), _logicalAnd).(bool) || mml.BinaryOp(11, mml.Ref(_c, "op"), _logicalOr).(bool))))
					},
					FixedArgs: 1,
				}
				_left = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left")))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _left))
				if c.(bool) {
					mml.Nop()
					return _left
				}
				_right = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right")))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _right))
				if c.(bool) {
					mml.Nop()
					return _right
				}
				c = (!_isBool.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left"))).(bool) && !_isBoolOp.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left"))).(bool))
				if c.(bool) {
					mml.Nop()
					_left = mml.BinaryOp(9, _left, ".(bool)")
				}
				c = (!_isBool.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right"))).(bool) && !_isBoolOp.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right"))).(bool))
				if c.(bool) {
					mml.Nop()
					_right = mml.BinaryOp(9, _right, ".(bool)")
				}
				_op = "&&"
				c = mml.BinaryOp(11, mml.Ref(_b, "op"), _logicalOr)
				if c.(bool) {
					mml.Nop()
					_op = "||"
				}
				return _formats.(*mml.Function).Call(append([]interface{}{}, "(%s %s %s)", _left, _op, _right))
				return nil
			},
			FixedArgs: 1,
		}
		_compileTernary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]

				mml.Nop(_c)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "func () interface{} { c = %s; if c.(bool) { return %s } else { return %s } }()", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "condition"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "consequent"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "alternative")))))
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
					c = _has.(*mml.Function).Call(append([]interface{}{}, "alternative", _c))
					if c.(bool) {
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "c = %s; if c.(bool) { %s } else { %s }", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "condition"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "consequent"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "alternative")))))
					} else {
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "c = %s; if c.(bool) { %s }", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "condition"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "consequent")))))
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
				var _exp interface{}
				var _cases interface{}
				var _def interface{}
				var _defaultCode interface{}
				var _casesCode interface{}
				mml.Nop(_hasDefault, _exp, _cases, _def, _defaultCode, _casesCode)
				_hasDefault = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_s, "defaultStatements"), "statements"))), 0)
				_exp = func() interface{} {
					c = _has.(*mml.Function).Call(append([]interface{}{}, "expression", _s))
					if c.(bool) {
						return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "expression")))
					} else {
						return ""
					}
				}()
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _exp))
				if c.(bool) {
					mml.Nop()
					return _exp
				}
				_cases = _anyErr.(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "cases")))))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _cases))
				if c.(bool) {
					mml.Nop()
					return _cases
				}
				_def = func() interface{} {
					c = _hasDefault
					if c.(bool) {
						return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "defaultStatements")))
					} else {
						return ""
					}
				}()
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _def))
				if c.(bool) {
					mml.Nop()
					return _def
				}
				_defaultCode = func() interface{} {
					c = _hasDefault
					if c.(bool) {
						return _formats.(*mml.Function).Call(append([]interface{}{}, "default:\n%s", _def))
					} else {
						return ""
					}
				}()
				_casesCode = _join.(*mml.Function).Call(append([]interface{}{}, "\n", func() interface{} {
					c = _hasDefault
					if c.(bool) {
						return append(append([]interface{}{}, _cases.([]interface{})...), _defaultCode)
					} else {
						return _cases
					}
				}()))
				return _formats.(*mml.Function).Call(append([]interface{}{}, "switch %s {\n%s\n}", _exp, _casesCode))
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
						return _passErr.(*mml.Function).Call(append([]interface{}{}, _formatOne.(*mml.Function).Call(append([]interface{}{}, "func() interface{} {\nselect {\n%s\n} }()")))).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _join.(*mml.Function).Call(append([]interface{}{}, "\n")))).(*mml.Function).Call(append([]interface{}{}, _anyErr.(*mml.Function).Call(append([]interface{}{}, func() interface{} {
							c = mml.Ref(_s, "hasDefault")
							if c.(bool) {
								return append(append([]interface{}{}, _c.([]interface{})...), _passErr.(*mml.Function).Call(append([]interface{}{}, _formatOne.(*mml.Function).Call(append([]interface{}{}, "default:\n%s")))).(*mml.Function).Call(append([]interface{}{}, _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "defaultStatements"))))))
							} else {
								return _c
							}
						}()))))))
					},
					FixedArgs: 1,
				}).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "cases")))))
			},
			FixedArgs: 1,
		}
		_compileRangeOver = &mml.Function{
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
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "_%s := 0; true; _%s++", mml.Ref(_r, "symbol"), mml.Ref(_r, "symbol")))
					},
					FixedArgs: 0,
				}
				_withRangeExpression = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "_%s := %s; %s; _%s++", mml.Ref(_r, "symbol"), func() interface{} {
							c = _has.(*mml.Function).Call(append([]interface{}{}, "from", mml.Ref(_r, "expression")))
							if c.(bool) {
								return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_r, "expression"), "from")))
							} else {
								return "0"
							}
						}(), func() interface{} {
							c = _has.(*mml.Function).Call(append([]interface{}{}, "to", mml.Ref(_r, "expression")))
							if c.(bool) {
								return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s < %s.(int)", mml.Ref(_r, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_r, "expression"), "to")))))
							} else {
								return "true"
							}
						}(), mml.Ref(_r, "symbol")))
					},
					FixedArgs: 0,
				}
				_listStyleRange = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)

						mml.Nop()
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "_, _%s := range %s.([]interface{})", mml.Ref(_r, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_r, "expression")))))
					},
					FixedArgs: 0,
				}
				switch {
				case !_has.(*mml.Function).Call(append([]interface{}{}, "expression", _r)).(bool):

					mml.Nop()
					return _infiniteCounter.(*mml.Function).Call([]interface{}{})
				case (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_r, "expression"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_r, "expression"), "type"), "range-expression").(bool)):

					mml.Nop()
					return _withRangeExpression.(*mml.Function).Call([]interface{}{})
				default:

					mml.Nop()
					return _listStyleRange.(*mml.Function).Call([]interface{}{})
				}
				return nil
			},
			FixedArgs: 1,
		}
		_compileLoop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]

				mml.Nop(_l)
				return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "for %s {\n%s\n}", func() interface{} {
					c = _has.(*mml.Function).Call(append([]interface{}{}, "expression", _l))
					if c.(bool) {
						return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "expression")))
					} else {
						return ""
					}
				}(), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "body")))))
			},
			FixedArgs: 1,
		}
		_compileDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0]

				mml.Nop(_d)
				return func() interface{} {
					c = mml.Ref(_d, "exported")
					if c.(bool) {
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "_%s = %s; exports[\"%s\"] = _%s", mml.Ref(_d, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_d, "expression"))), mml.Ref(_d, "symbol"), mml.Ref(_d, "symbol")))
					} else {
						return _passErrFormat.(*mml.Function).Call(append([]interface{}{}, "_%s = %s", mml.Ref(_d, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_d, "expression")))))
					}
				}()
			},
			FixedArgs: 1,
		}
		_controlStatement = _enum.(*mml.Function).Call([]interface{}{})
		_breakControl = _controlStatement.(*mml.Function).Call([]interface{}{})
		_continueControl = _controlStatement.(*mml.Function).Call([]interface{}{})
		_compileStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]

				mml.Nop(_s)
				var _scope interface{}
				var _scopeDefs interface{}
				var _scopeNames interface{}
				var _statements interface{}
				mml.Nop(_scope, _scopeDefs, _scopeNames, _statements)
				_scope = _getScope.(*mml.Function).Call(append([]interface{}{}, _s.([]interface{})...))
				_scopeDefs = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)
						return _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s interface{}", _s))
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _scope))))
				_scopeNames = _join.(*mml.Function).Call(append([]interface{}{}, ", ", _map.(*mml.Function).Call(append([]interface{}{}, _formatOne.(*mml.Function).Call(append([]interface{}{}, "_%s")), _scope))))
				_statements = _passErr.(*mml.Function).Call(append([]interface{}{}, _join.(*mml.Function).Call(append([]interface{}{}, ";\n")))).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _notEmpty)).(*mml.Function).Call(append([]interface{}{}, _mapCompile.(*mml.Function).Call(append([]interface{}{}, _s))))))
				return _formats.(*mml.Function).Call(append([]interface{}{}, "%s;\nmml.Nop(%s);\n%s", _scopeDefs, _scopeNames, _statements))
				return nil
			},
			FixedArgs: 1,
		}
		_compileModule = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0]

				mml.Nop(_m)
				var _statements interface{}
				mml.Nop(_statements)
				_statements = _compileStatements.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_m, "statements")))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _statements))
				if c.(bool) {
					mml.Nop()
					return _statements
				}
				return _formats.(*mml.Function).Call(append([]interface{}{}, "%s", _statements))
				return nil
			},
			FixedArgs: 1,
		}
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
					_useStatement = _formats.(*mml.Function).Call(append([]interface{}{}, "var __%s = mml.Modules.Use(\"%s.mml\");", _getModuleName.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "path"))), mml.Ref(_u, "path")))
					_assigns = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _name = a[0]

							mml.Nop(_name)
							return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = __%s[\"%s\"]", _name, _getModuleName.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "path"))), _name))
						},
						FixedArgs: 1,
					}, mml.Ref(_u, "exportNames")))))
					return _joins.(*mml.Function).Call(append([]interface{}{}, ";", _useStatement, _assigns))
				case mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "capture"))), ""):

					mml.Nop()
					return ""
				default:

					mml.Nop()
					return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = mml.Modules.Use(\"%s.mml\")", _getModuleName.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "path"))), mml.Ref(_u, "path")))
				}
				return nil
			},
			FixedArgs: 1,
		}
		_compileUseList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]

				mml.Nop(_u)
				return _mapCompileJoin.(*mml.Function).Call(append([]interface{}{}, "\n;", mml.Ref(_u, "uses")))
			},
			FixedArgs: 1,
		}
		_compile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]

				mml.Nop(_code)

				mml.Nop()
				switch {
				case _isInt.(*mml.Function).Call(append([]interface{}{}, _code)):

					mml.Nop()
					return _compileInt.(*mml.Function).Call(append([]interface{}{}, _code))
				case _isFloat.(*mml.Function).Call(append([]interface{}{}, _code)):

					mml.Nop()
					return _compileFloat.(*mml.Function).Call(append([]interface{}{}, _code))
				case _isString.(*mml.Function).Call(append([]interface{}{}, _code)):

					mml.Nop()
					return _compileString.(*mml.Function).Call(append([]interface{}{}, _code))
				case _isBool.(*mml.Function).Call(append([]interface{}{}, _code)):

					mml.Nop()
					return _compileBool.(*mml.Function).Call(append([]interface{}{}, _code))
				}
				switch mml.Ref(_code, "type") {
				case "comment":

					mml.Nop()
					return _compileComment.(*mml.Function).Call(append([]interface{}{}, _code))
				case "symbol":

					mml.Nop()
					return _compileSymbol.(*mml.Function).Call(append([]interface{}{}, _code))
				case "module":

					mml.Nop()
					return _compileModule.(*mml.Function).Call(append([]interface{}{}, _code))
				case "list":

					mml.Nop()
					return _compileList.(*mml.Function).Call(append([]interface{}{}, _code))
				case "entry":

					mml.Nop()
					return _compileEntry.(*mml.Function).Call(append([]interface{}{}, _code))
				case "structure":

					mml.Nop()
					return _compileStructure.(*mml.Function).Call(append([]interface{}{}, _code))
				case "function":

					mml.Nop()
					return _compileFunction.(*mml.Function).Call(append([]interface{}{}, _code))
				case "range-expression":

					mml.Nop()
					return _compileRangeExpression.(*mml.Function).Call(append([]interface{}{}, _code))
				case "indexer":

					mml.Nop()
					return _compileIndexer.(*mml.Function).Call(append([]interface{}{}, _code))
				case "spread":

					mml.Nop()
					return _compileSpread.(*mml.Function).Call(append([]interface{}{}, _code))
				case "function-application":

					mml.Nop()
					return _compileApplication.(*mml.Function).Call(append([]interface{}{}, _code))
				case "unary":

					mml.Nop()
					return _compileUnary.(*mml.Function).Call(append([]interface{}{}, _code))
				case "binary":

					mml.Nop()
					return _compileBinary.(*mml.Function).Call(append([]interface{}{}, _code))
				case "cond":

					mml.Nop()
					return _compileCond.(*mml.Function).Call(append([]interface{}{}, _code))
				case "switch-case":

					mml.Nop()
					return _compileCase.(*mml.Function).Call(append([]interface{}{}, _code))
				case "switch-statement":

					mml.Nop()
					return _compileSwitch.(*mml.Function).Call(append([]interface{}{}, _code))
				case "send":

					mml.Nop()
					return _compileSend.(*mml.Function).Call(append([]interface{}{}, _code))
				case "receive":

					mml.Nop()
					return _compileReceive.(*mml.Function).Call(append([]interface{}{}, _code))
				case "go":

					mml.Nop()
					return _compileGo.(*mml.Function).Call(append([]interface{}{}, _code))
				case "defer":

					mml.Nop()
					return _compileDefer.(*mml.Function).Call(append([]interface{}{}, _code))
				case "select-case":

					mml.Nop()
					return _compileCase.(*mml.Function).Call(append([]interface{}{}, _code))
				case "select":

					mml.Nop()
					return _compileSelect.(*mml.Function).Call(append([]interface{}{}, _code))
				case "range-over":

					mml.Nop()
					return _compileRangeOver.(*mml.Function).Call(append([]interface{}{}, _code))
				case "loop":

					mml.Nop()
					return _compileLoop.(*mml.Function).Call(append([]interface{}{}, _code))
				case "definition":

					mml.Nop()
					return _compileDefinition.(*mml.Function).Call(append([]interface{}{}, _code))
				case "definition-list":

					mml.Nop()
					return _compileDefinitions.(*mml.Function).Call(append([]interface{}{}, _code))
				case "assign":

					mml.Nop()
					return _compileAssign.(*mml.Function).Call(append([]interface{}{}, _code))
				case "assign-list":

					mml.Nop()
					return _compileAssigns.(*mml.Function).Call(append([]interface{}{}, _code))
				case "ret":

					mml.Nop()
					return _compileRet.(*mml.Function).Call(append([]interface{}{}, _code))
				case "control-statement":

					mml.Nop()
					return _compileControl.(*mml.Function).Call(append([]interface{}{}, _code))
				case "statement-list":

					mml.Nop()
					return _compileStatements.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_code, "statements")))
				case "use":

					mml.Nop()
					return _compileUse.(*mml.Function).Call(append([]interface{}{}, _code))
				case "use-list":

					mml.Nop()
					return _compileUseList.(*mml.Function).Call(append([]interface{}{}, _code))
				default:

					mml.Nop()
					return _error.(*mml.Function).Call(append([]interface{}{}, _formats.(*mml.Function).Call(append([]interface{}{}, "unsupported code: %v", _code))))
				}
				return nil
			},
			FixedArgs: 1,
		}
		_builtin = map[string]interface{}{"len": "Len", "isError": "IsError", "keys": "Keys", "format": "Format", "stdin": "Stdin", "stdout": "Stdout", "stderr": "Stderr", "string": "String", "parse": "Parse", "has": "Has", "isBool": "IsBool", "isInt": "IsInt", "isFloat": "IsFloat", "isString": "IsString", "error": "Error", "open": "Open", "close": "Close", "args": "Args"}
		_builtins = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _k = a[0]

				mml.Nop(_k)
				return _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s interface{} = mml.%s", _k, mml.Ref(_builtin, _k)))
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _left = a[0]
				var _right = a[1]

				mml.Nop(_left, _right)
				return mml.BinaryOp(13, _left, _right)
			},
			FixedArgs: 2,
		})).(*mml.Function).Call(append([]interface{}{}, _keys.(*mml.Function).Call(append([]interface{}{}, _builtin))))))))
		_parseFile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]

				mml.Nop(_path)
				var _in interface{}
				mml.Nop(_in)
				_in = _open.(*mml.Function).Call(append([]interface{}{}, _path))
				c = _isError.(*mml.Function).Call(append([]interface{}{}, _in))
				if c.(bool) {
					mml.Nop()
					return _in
				}
				defer _close.(*mml.Function).Call(append([]interface{}{}, _in))
				return _onlyErr.(*mml.Function).Call(append([]interface{}{}, _log)).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _parse)).(*mml.Function).Call(append([]interface{}{}, _in.(*mml.Function).Call(append([]interface{}{}, mml.UnaryOp(2, 1)))))))
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
				return _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]

						mml.Nop(_d)
						return mml.Ref(_d, "symbol")
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]

						mml.Nop(_d)
						return mml.Ref(_d, "exported")
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _getFlattenedStatements.(*mml.Function).Call(append([]interface{}{}, "definition", "definition-list", "definitions")).(*mml.Function).Call(append([]interface{}{}, _statements))))))
			},
			FixedArgs: 1,
		}
		_parseModules = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]

				mml.Nop(_path)
				var _code interface{}
				var _uses interface{}
				var _usesModules interface{}
				var _statements interface{}
				var _pimpedCode interface{}
				mml.Nop(_code, _uses, _usesModules, _statements, _pimpedCode)
				_code = _parseFile.(*mml.Function).Call(append([]interface{}{}, _path))
				_uses = _getFlattenedStatements.(*mml.Function).Call(append([]interface{}{}, "use", "use-list", "uses", mml.Ref(_code, "statements")))
				_usesModules = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]

						mml.Nop(_m)
						return map[string]interface{}{"type": mml.Ref(_m, "type"), "path": mml.Ref(_m, "path"), "statements": mml.Ref(_m, "statements"), "exportNames": _findExportNames.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_m, "statements")))}
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]
						var _all = a[1]

						mml.Nop(_m, _all)
						return append(append([]interface{}{}, _all.([]interface{})...), _m.([]interface{})...)
					},
					FixedArgs: 2,
				}, []interface{}{})).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _parseModules)).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _u = a[0]

						mml.Nop(_u)
						return mml.BinaryOp(9, mml.Ref(_u, "path"), ".mml")
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, _uses))))))))
				_statements = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]

						mml.Nop(_s)

						mml.Nop()
						c = (!_has.(*mml.Function).Call(append([]interface{}{}, "type", _s)).(bool) || (mml.BinaryOp(12, mml.Ref(_s, "type"), "use").(bool) && mml.BinaryOp(12, mml.Ref(_s, "type"), "use-list").(bool)))
						if c.(bool) {
							mml.Nop()
							return _s
						}
						c = mml.BinaryOp(11, mml.Ref(_s, "type"), "use")
						if c.(bool) {
							var _m interface{}
							mml.Nop(_m)
							_m = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
								F: func(a []interface{}) interface{} {
									var c interface{}
									mml.Nop(c)
									var _m = a[0]

									mml.Nop(_m)
									return mml.BinaryOp(11, mml.Ref(_m, "path"), mml.Ref(_s, "path"))
								},
								FixedArgs: 1,
							}, _usesModules))
							c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _m)), 0)
							if c.(bool) {
								mml.Nop()
								return _s
							}
							return map[string]interface{}{"type": mml.Ref(_s, "type"), "path": mml.Ref(_s, "path"), "capture": mml.Ref(_s, "capture"), "exportNames": mml.Ref(mml.Ref(_m, 0), "exportNames")}
						}
						return map[string]interface{}{"type": mml.Ref(_s, "type"), "uses": _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _u = a[0]

								mml.Nop(_u)
								var _m interface{}
								mml.Nop(_m)
								_m = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
									F: func(a []interface{}) interface{} {
										var c interface{}
										mml.Nop(c)
										var _m = a[0]

										mml.Nop(_m)
										return mml.BinaryOp(11, mml.Ref(_m, "path"), mml.BinaryOp(9, mml.Ref(_u, "path"), ".mml"))
									},
									FixedArgs: 1,
								}, _usesModules))
								c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _m)), 0)
								if c.(bool) {
									mml.Nop()
									return _u
								}
								return map[string]interface{}{"type": mml.Ref(_u, "type"), "path": mml.Ref(_u, "path"), "capture": mml.Ref(_u, "capture"), "exportNames": mml.Ref(mml.Ref(_m, 0), "exportNames")}
								return nil
							},
							FixedArgs: 1,
						}, mml.Ref(_s, "uses")))}
						return nil
					},
					FixedArgs: 1,
				})).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_code, "statements")))
				_pimpedCode = map[string]interface{}{"type": mml.Ref(_code, "type"), "path": _path, "statements": _statements}
				return append(append([]interface{}{}, _pimpedCode), _usesModules.([]interface{})...)
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
				_stdout.(*mml.Function).Call(append([]interface{}{}, _formats.(*mml.Function).Call(append([]interface{}{}, "modulePath = \"%s\"", mml.Ref(_moduleCode, "path")))))
				_stdout.(*mml.Function).Call(append([]interface{}{}, "\n\t\tmml.Modules.Set(modulePath, func() map[string]interface{} {\n\t\t\texports := make(map[string]interface{})\n\n\t\t\tvar c interface{}\n\t\t\tmml.Nop(c)\n\t"))
				_onlyErr.(*mml.Function).Call(append([]interface{}{}, _log)).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _stdout)).(*mml.Function).Call(append([]interface{}{}, _compile.(*mml.Function).Call(append([]interface{}{}, _moduleCode))))))
				_stdout.(*mml.Function).Call(append([]interface{}{}, "\n\t\t\treturn exports\n\t\t})\n\t"))
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
				for _, _mi := range _m.([]interface{}) {

					mml.Nop()
					_compileModuleCode.(*mml.Function).Call(append([]interface{}{}, _mi))
				}
				return nil
			},
			FixedArgs: 1,
		}
		_stdout.(*mml.Function).Call(append([]interface{}{}, "// Generated code\n\n\tpackage main\n\n\timport \"github.com/aryszka/mml\"\n"))
		_stdout.(*mml.Function).Call(append([]interface{}{}, _builtins))
		_stdout.(*mml.Function).Call(append([]interface{}{}, "\n\tfunc init() {\n\t\tvar modulePath string\n"))
		_modules = _parseModules.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_args, 1)))
		_compileModules.(*mml.Function).Call(append([]interface{}{}, _modules))
		_stdout.(*mml.Function).Call(append([]interface{}{}, "\n\t}\n\n\tfunc main() {\n\t\tmml.Modules.Use(\""))
		_stdout.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_args, 1)))
		_stdout.(*mml.Function).Call(append([]interface{}{}, "\")\n\t}\n"))
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
		var _firstOr interface{}
		var _contains interface{}
		mml.Nop(_fold, _foldr, _map, _filter, _firstOr, _contains)
		_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _i = a[1]
				var _l = a[2]

				mml.Nop(_f, _i, _l)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0)
					if c.(bool) {
						return _i
					} else {
						return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil)))
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
					c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0)
					if c.(bool) {
						return _i
					} else {
						return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil)))))
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
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]

						mml.Nop(_c, _r)
						return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
					},
					FixedArgs: 2,
				}, []interface{}{}, _l))
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
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]

						mml.Nop(_c, _r)
						return func() interface{} {
							c = _p.(*mml.Function).Call(append([]interface{}{}, _c))
							if c.(bool) {
								return append(append([]interface{}{}, _r.([]interface{})...), _c)
							} else {
								return _r
							}
						}()
					},
					FixedArgs: 2,
				}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}
		exports["filter"] = _filter
		_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0]
				var _l = a[1]

				mml.Nop(_v, _l)
				return func() interface{} {
					c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0)
					if c.(bool) {
						return mml.Ref(_l, 0)
					} else {
						return _v
					}
				}()
			},
			FixedArgs: 2,
		}
		exports["firstOr"] = _firstOr
		_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0]
				var _l = a[1]

				mml.Nop(_i, _l)
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ii = a[0]

						mml.Nop(_ii)
						return mml.BinaryOp(11, _ii, _i)
					},
					FixedArgs: 1,
				}, _l)))), 0)
			},
			FixedArgs: 2,
		}
		exports["contains"] = _contains
		return exports
	})

}

func main() {
	mml.Modules.Use("compile.mml")
}
