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
	modulePath = "main"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _read interface{}
		var _errors interface{}
		var _compile interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_read, _errors, _compile, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_read = mml.Modules.Use("read")
		_errors = mml.Modules.Use("errors")
		_compile = mml.Modules.Use("compile")
		c = mml.BinaryOp(13, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _args)}).Values), 2)
		if c.(bool) {
			mml.Nop()
			_fatal.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "usage: mml <source without the extension>")}).Values)
		}
		mml.Ref(_errors, "only").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fatal)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_errors, "pass").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _stdout)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_errors, "pass").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_compile, "toGo"))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_read, "do").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_args, 1))}).Values))}).Values))}).Values))}).Values)

		return exports
	})

	modulePath = "lang"

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
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		var _lists interface{}
		var _strings interface{}
		var _ints interface{}
		var _functions interface{}
		var _match interface{}
		var _logger interface{}
		mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is, _lists, _strings, _ints, _functions, _match, _logger)
		_lists = mml.Modules.Use("lists")
		_strings = mml.Modules.Use("strings")
		_ints = mml.Modules.Use("ints")
		_logger = mml.Modules.Use("log")
		_functions = mml.Modules.Use("functions")
		_match = mml.Modules.Use("match")
		_fold = mml.Ref(_lists, "fold")
		exports["fold"] = _fold
		_foldr = mml.Ref(_lists, "foldr")
		exports["foldr"] = _foldr
		_map = mml.Ref(_lists, "map")
		exports["map"] = _map
		_filter = mml.Ref(_lists, "filter")
		exports["filter"] = _filter
		_contains = mml.Ref(_lists, "contains")
		exports["contains"] = _contains
		_sort = mml.Ref(_lists, "sort")
		exports["sort"] = _sort
		_flat = mml.Ref(_lists, "flat")
		exports["flat"] = _flat
		_flats = mml.Ref(_lists, "flats")
		exports["flats"] = _flats
		_concat = mml.Ref(_lists, "concat")
		exports["concat"] = _concat
		_concats = mml.Ref(_lists, "concats")
		exports["concats"] = _concats
		_uniq = mml.Ref(_lists, "uniq")
		exports["uniq"] = _uniq
		_every = mml.Ref(_lists, "every")
		exports["every"] = _every
		_some = mml.Ref(_lists, "some")
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
		_fatal = mml.Ref(_logger, "fatal")
		exports["fatal"] = _fatal
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
		_natural = mml.Ref(_match, "natural")
		exports["natural"] = _natural
		_type = mml.Ref(_match, "type")
		exports["type"] = _type
		_listOf = mml.Ref(_match, "listOf")
		exports["listOf"] = _listOf
		_structOf = mml.Ref(_match, "structOf")
		exports["structOf"] = _structOf
		_range = mml.Ref(_match, "range")
		exports["range"] = _range
		_rangeMin = mml.Ref(_match, "rangeMin")
		exports["rangeMin"] = _rangeMin
		_listLength = mml.Ref(_match, "listLength")
		exports["listLength"] = _listLength
		_or = mml.Ref(_match, "or")
		exports["or"] = _or
		_and = mml.Ref(_match, "and")
		exports["and"] = _and
		_not = mml.Ref(_match, "not")
		exports["not"] = _not
		_predicate = mml.Ref(_match, "predicate")
		exports["predicate"] = _predicate
		_predicates = mml.Ref(_match, "predicates")
		exports["predicates"] = _predicates
		_is = mml.Ref(_match, "is")
		exports["is"] = _is

		return exports
	})

	modulePath = "lists"

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
		var _concat interface{}
		var _concats interface{}
		var _flat interface{}
		var _flats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _intersect interface{}
		var _sort interface{}
		var _group interface{}
		var _indexes interface{}
		var _flatDepth interface{}
		mml.Nop(_fold, _foldr, _map, _filter, _first, _contains, _concat, _concats, _flat, _flats, _uniq, _every, _some, _intersect, _sort, _group, _indexes, _flatDepth)
		_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _i = a[1]
				var _l = a[2]
				var _ interface{}
				_ = &mml.List{a[3:]}
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
				var _ interface{}
				_ = &mml.List{a[3:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_m, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_p, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _r = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_i, _l)
				return mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _first.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ii = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_ii)
						return mml.BinaryOp(11, _ii, _i)
					},
					FixedArgs: 1,
				}, _l)}).Values))}).Values), 0)
			},
			FixedArgs: 2,
		}
		exports["contains"] = _contains
		_concat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				return _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values)
			},
			FixedArgs: 1,
		}
		exports["concat"] = _concat
		_concats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l interface{}
				_l = &mml.List{a[0:]}
				mml.Nop(_l)
				return _concat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values)
			},
			FixedArgs: 0,
		}
		exports["concats"] = _concats
		_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				return _flatDepth.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1, _l)}).Values)
			},
			FixedArgs: 1,
		}
		exports["flat"] = _flat
		_flats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l interface{}
				_l = &mml.List{a[0:]}
				mml.Nop(_l)
				return _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values)
			},
			FixedArgs: 0,
		}
		exports["flats"] = _flats
		_uniq = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _eq = a[0]
				var _l = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_eq, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _u = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_c, _u)
						return func() interface{} {
							c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
								F: func(a []interface{}) interface{} {
									var c interface{}
									mml.Nop(c)
									var _i = a[0]
									var _ interface{}
									_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_p, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _r = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_p, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _r = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_i, _r)
						return (_r.(bool) || _p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _i)}).Values).(bool))
					},
					FixedArgs: 2,
				}, false, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["some"] = _some
		_intersect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l0 = a[0]
				var _l1 = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_l0, _l1)
				return _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i0 = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_i0)
						return _some.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _i1 = a[0]
								var _ interface{}
								_ = &mml.List{a[1:]}
								mml.Nop(_i1)
								return mml.BinaryOp(11, _i0, _i1)
							},
							FixedArgs: 1,
						}, _l1)}).Values)
					},
					FixedArgs: 1,
				}, _l0)}).Values)
			},
			FixedArgs: 2,
		}
		exports["intersect"] = _intersect
		_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0]
				var _l = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
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
								var _ interface{}
								_ = &mml.List{a[1:]}
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
		_group = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _l = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_n, _l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _g = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_i, _g)
						return func() interface{} {
							c = (mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _g)}).Values), 0).(bool) || mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_g, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _g)}).Values), 1)))}).Values), _n).(bool))
							if c.(bool) {
								return &mml.List{Values: append(append([]interface{}{}, _g.(*mml.List).Values...), &mml.List{Values: append([]interface{}{}, _i)})}
							} else {
								return &mml.List{Values: append(append([]interface{}{}, mml.RefRange(_g, nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _g)}).Values), 1)).(*mml.List).Values...), &mml.List{Values: append(append([]interface{}{}, mml.Ref(_g, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _g)}).Values), 1)).(*mml.List).Values...), _i)})}
							}
						}()
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 2,
		}
		exports["group"] = _group
		_indexes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var __ = a[0]
						var _i = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(__, _i)
						return &mml.List{Values: append(append([]interface{}{}, _i.(*mml.List).Values...), _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _i)}).Values))}
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
			},
			FixedArgs: 1,
		}
		exports["indexes"] = _indexes
		_flatDepth = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0]
				var _l = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_d, _l)

				mml.Nop()
				c = mml.BinaryOp(11, _d, 0)
				if c.(bool) {
					mml.Nop()
					return _l
				}
				return _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _r = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_i, _r)
						var _fi interface{}
						mml.Nop(_fi)
						c = !_isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _i)}).Values).(bool)
						if c.(bool) {
							mml.Nop()
							return &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _i)}
						}
						_fi = _flatDepth.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.BinaryOp(10, _d, 1), _i)}).Values)
						return &mml.List{Values: append(append([]interface{}{}, _r.(*mml.List).Values...), _fi.(*mml.List).Values...)}
						return nil
					},
					FixedArgs: 2,
				}, &mml.List{Values: []interface{}{}}, _l)}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		exports["flatDepth"] = _flatDepth

		return exports
	})

	modulePath = "strings"

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
				var _ interface{}
				_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[3:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
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

	modulePath = "ints"

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
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				var _c interface{}
				mml.Nop(_c)
				_c = mml.UnaryOp(2, 1)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ interface{}
						_ = &mml.List{a[0:]}
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

	modulePath = "log"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _println interface{}
		var _fatal interface{}
		var _lists interface{}
		var _strings interface{}
		mml.Nop(_println, _fatal, _lists, _strings)
		_lists = mml.Modules.Use("lists")
		_strings = mml.Modules.Use("strings")
		_println = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a interface{}
				_a = &mml.List{a[0:]}
				mml.Nop(_a)

				mml.Nop()
				_stderr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, " ")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _string)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values))}).Values))}).Values)
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
		_fatal = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a interface{}
				_a = &mml.List{a[0:]}
				mml.Nop(_a)

				mml.Nop()
				_println.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a.(*mml.List).Values...)}).Values)
				_exit.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values)
				return nil
			},
			FixedArgs: 0,
		}
		exports["fatal"] = _fatal

		return exports
	})

	modulePath = "functions"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _identity interface{}
		var _eq interface{}
		var _not interface{}
		var _apply interface{}
		var _call interface{}
		var _chain interface{}
		var _chains interface{}
		var _bindAt interface{}
		var _bind interface{}
		var _only interface{}
		var _lists interface{}
		mml.Nop(_identity, _eq, _not, _apply, _call, _chain, _chains, _bindAt, _bind, _only, _lists)
		_lists = mml.Modules.Use("lists")
		_identity = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_p)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _a = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_a)
						return !_p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values).(bool)
					},
					FixedArgs: 1,
				}
			},
			FixedArgs: 1,
		}
		exports["not"] = _not
		_apply = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_f, _a)
				return _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a.(*mml.List).Values...)}).Values)
			},
			FixedArgs: 2,
		}
		exports["apply"] = _apply
		_call = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a interface{}
				_a = &mml.List{a[1:]}
				mml.Nop(_f, _a)
				return _apply.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _a)}).Values)
			},
			FixedArgs: 1,
		}
		exports["call"] = _call
		_chain = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_f)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _a = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_a)
						return mml.Ref(_lists, "fold").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _call, _a, _f)}).Values)
					},
					FixedArgs: 1,
				}
			},
			FixedArgs: 1,
		}
		exports["chain"] = _chain
		_chains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f interface{}
				_f = &mml.List{a[0:]}
				mml.Nop(_f)
				return _chain.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f)}).Values)
			},
			FixedArgs: 0,
		}
		exports["chains"] = _chains
		_bindAt = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0]
				var _f = a[1]
				var _a interface{}
				_a = &mml.List{a[2:]}
				mml.Nop(_i, _f, _a)
				return &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _b interface{}
						_b = &mml.List{a[0:]}
						mml.Nop(_b)
						return _f.(*mml.Function).Call((&mml.List{Values: append(append(append([]interface{}{}, mml.RefRange(_b, nil, _i).(*mml.List).Values...), _a.(*mml.List).Values...), mml.RefRange(_b, _i, nil).(*mml.List).Values...)}).Values)
					},
					FixedArgs: 0,
				}
			},
			FixedArgs: 2,
		}
		exports["bindAt"] = _bindAt
		_bind = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _a interface{}
				_a = &mml.List{a[1:]}
				mml.Nop(_f, _a)
				return _bindAt.(*mml.Function).Call((&mml.List{Values: append(append([]interface{}{}, 0, _f), _a.(*mml.List).Values...)}).Values)
			},
			FixedArgs: 1,
		}
		exports["bind"] = _bind
		_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0]
				var _f interface{}
				_f = &mml.List{a[1:]}
				mml.Nop(_p, _f)
				return _chain.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_f)
						return &mml.Function{
							F: func(a []interface{}) interface{} {
								var c interface{}
								mml.Nop(c)
								var _a = a[0]
								var _ interface{}
								_ = &mml.List{a[1:]}
								mml.Nop(_a)
								return func() interface{} {
									c = _p.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _a)}).Values)
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
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f)}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		exports["only"] = _only

		return exports
	})

	modulePath = "match"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _complexType interface{}
		var _defineRange interface{}
		var _listRange interface{}
		var _isSimpleType interface{}
		var _isComplexType interface{}
		var _isType interface{}
		var _complexTypeEq interface{}
		var _primitives interface{}
		var _matchPrimitive interface{}
		var _matchToList interface{}
		var _matchToListType interface{}
		var _matchList interface{}
		var _matchStruct interface{}
		var _matchOne interface{}
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
		var _intRangeType interface{}
		var _floatRangeType interface{}
		var _isRange interface{}
		var _isNaturalRange interface{}
		var _intRange interface{}
		var _floatRange interface{}
		var _stringRangeType interface{}
		var _stringRange interface{}
		var _listType interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _unionType interface{}
		var _intersectType interface{}
		var _predicateType interface{}
		var _or interface{}
		var _and interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _matchInt interface{}
		var _matchFloat interface{}
		var _matchString interface{}
		var _matchUnion interface{}
		var _matchIntersection interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _not interface{}
		var _natural interface{}
		var _is interface{}
		var _functions interface{}
		var _ints interface{}
		var _floats interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _first interface{}
		var _contains interface{}
		var _concat interface{}
		var _concats interface{}
		var _flat interface{}
		var _flats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _intersect interface{}
		var _sort interface{}
		var _group interface{}
		var _indexes interface{}
		var _flatDepth interface{}
		mml.Nop(_complexType, _defineRange, _listRange, _isSimpleType, _isComplexType, _isType, _complexTypeEq, _primitives, _matchPrimitive, _matchToList, _matchToListType, _matchList, _matchStruct, _matchOne, _token, _none, _integer, _floating, _stringType, _boolean, _errorType, _any, _function, _channel, _type, _intRangeType, _floatRangeType, _isRange, _isNaturalRange, _intRange, _floatRange, _stringRangeType, _stringRange, _listType, _listOf, _structOf, _range, _unionType, _intersectType, _predicateType, _or, _and, _predicate, _predicates, _matchInt, _matchFloat, _matchString, _matchUnion, _matchIntersection, _rangeMin, _listLength, _not, _natural, _is, _functions, _ints, _floats, _fold, _foldr, _map, _filter, _first, _contains, _concat, _concats, _flat, _flats, _uniq, _every, _some, _intersect, _sort, _group, _indexes, _flatDepth)
		var __lists = mml.Modules.Use("lists")
		_fold = __lists.Values["fold"]
		_foldr = __lists.Values["foldr"]
		_map = __lists.Values["map"]
		_filter = __lists.Values["filter"]
		_first = __lists.Values["first"]
		_contains = __lists.Values["contains"]
		_concat = __lists.Values["concat"]
		_concats = __lists.Values["concats"]
		_flat = __lists.Values["flat"]
		_flats = __lists.Values["flats"]
		_uniq = __lists.Values["uniq"]
		_every = __lists.Values["every"]
		_some = __lists.Values["some"]
		_intersect = __lists.Values["intersect"]
		_sort = __lists.Values["sort"]
		_group = __lists.Values["group"]
		_indexes = __lists.Values["indexes"]
		_flatDepth = __lists.Values["flatDepth"]
		_functions = mml.Modules.Use("functions")
		_ints = mml.Modules.Use("ints")
		_floats = mml.Modules.Use("floats")
		_token = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _token
			},
			FixedArgs: 0,
		}
		_none = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _none
			},
			FixedArgs: 0,
		}
		_integer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _integer
			},
			FixedArgs: 0,
		}
		_floating = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _floating
			},
			FixedArgs: 0,
		}
		_stringType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _stringType
			},
			FixedArgs: 0,
		}
		_boolean = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _boolean
			},
			FixedArgs: 0,
		}
		_errorType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _errorType
			},
			FixedArgs: 0,
		}
		_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
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
				var _ interface{}
				_ = &mml.List{a[0:]}
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
				var _ interface{}
				_ = &mml.List{a[0:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_complexType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _name = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_intRangeType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "int-range")}).Values)
		_floatRangeType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "float-range")}).Values)
		_isRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ofType = a[0]
				var _min = a[1]
				var _max = a[2]
				var _ interface{}
				_ = &mml.List{a[3:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_min, _max)
				return (_isRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isInt, _min, _max)}).Values).(bool) && mml.BinaryOp(16, _min, 0).(bool))
			},
			FixedArgs: 2,
		}
		_defineRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ofType = a[0]
				var _validate = a[1]
				var _min = a[2]
				var _max = a[3]
				var _ interface{}
				_ = &mml.List{a[4:]}
				mml.Nop(_ofType, _validate, _min, _max)
				return func() interface{} {
					c = _validate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min, _max)}).Values)
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _ofType.(*mml.Struct)
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
		_stringRangeType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "string")}).Values)
		_stringRange = _defineRange.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _stringRangeType, _isNaturalRange)}).Values)
		_listType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "list")}).Values)
		_listRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _item = a[0]
				var _min = a[1]
				var _max = a[2]
				var _ interface{}
				_ = &mml.List{a[3:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[3:]}
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
				case _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _m)}).Values):

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
		_unionType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "union")}).Values)
		_intersectType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "intersection")}).Values)
		_predicateType = _complexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "predicate")}).Values)
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_p)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _predicateType.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["predicate"] = _p
					return s
				}()
			},
			FixedArgs: 1,
		}
		exports["predicate"] = _predicate
		_predicates = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p interface{}
				_p = &mml.List{a[0:]}
				mml.Nop(_p)
				return _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _predicate, _p)}).Values).(*mml.List).Values...)}).Values)
			},
			FixedArgs: 0,
		}
		exports["predicates"] = _predicates
		_isSimpleType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _t = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_t)
				return _some.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "bind").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "eq"), _t)}).Values), &mml.List{Values: append([]interface{}{}, _integer, _floating, _stringType, _boolean, _function, _errorType, _channel)})}).Values)
			},
			FixedArgs: 1,
		}
		_isComplexType = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _t = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_t)
				return (_isSimpleType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _t)}).Values).(bool) || _isComplexType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _t)}).Values).(bool))
			},
			FixedArgs: 1,
		}
		_complexTypeEq = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _type = a[0]
				var _value = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
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
				s.Values["rangeValue"] = mml.Ref(_functions, "identity")
				return s
			}()
			s.Values["float"] = func() interface{} {
				s := &mml.Struct{Values: make(map[string]interface{})}
				s.Values["checkValue"] = _isFloat
				s.Values["type"] = _floating
				s.Values["rangeType"] = _floatRangeType
				s.Values["rangeValue"] = mml.Ref(_functions, "identity")
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
				var _ interface{}
				_ = &mml.List{a[3:]}
				mml.Nop(_def, _match, _value)

				mml.Nop()
				switch {
				case !mml.Ref(_def, "checkValue").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool):

					mml.Nop()
					return false
				case mml.BinaryOp(11, _match, mml.Ref(_def, "type")):

					mml.Nop()
					return true
				case _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_def, "rangeType"), _match)}).Values):
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
				var _ interface{}
				_ = &mml.List{a[2:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _value)
				return ((mml.BinaryOp(16, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values), mml.Ref(_match, "min")).(bool) && mml.BinaryOp(14, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values), mml.Ref(_match, "max")).(bool)) && (mml.BinaryOp(11, mml.Ref(_match, "item"), _any).(bool) || _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_match, "item"))}).Values), _value)}).Values).(bool)))
			},
			FixedArgs: 2,
		}
		_matchList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _value = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _value)

				mml.Nop()
				switch {
				case !_isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool):

					mml.Nop()
					return false
				case _isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values):

					mml.Nop()
					return _matchToList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _match)}).Values):

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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _value)
				return (_isStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool) && _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _key = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _value)
				return _some.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _value)
				return _every.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _m = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _value)

				mml.Nop()
				switch {
				case mml.BinaryOp(11, _match, _none):

					mml.Nop()
					return false
				case mml.BinaryOp(11, _match, _any):

					mml.Nop()
					return true
				case _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _predicateType, _match)}).Values):

					mml.Nop()
					return mml.Ref(_match, "predicate").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values)
				case (_isType.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _value)}).Values).(bool) && mml.BinaryOp(12, _value, _function).(bool)):

					mml.Nop()
					return false
				case mml.BinaryOp(11, _match, _value):

					mml.Nop()
					return true
				case (mml.BinaryOp(11, _match, _integer).(bool) || _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _intRangeType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case (mml.BinaryOp(11, _match, _floating).(bool) || _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _floatRangeType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchFloat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case (mml.BinaryOp(11, _match, _stringType).(bool) || _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _stringRangeType, _match)}).Values).(bool)):

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
				case (_isList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match)}).Values).(bool) || _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listType, _match)}).Values).(bool)):

					mml.Nop()
					return _matchList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _unionType, _match)}).Values):

					mml.Nop()
					return _matchUnion.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _value)}).Values)
				case _complexTypeEq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _intersectType, _match)}).Values):

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
		_rangeMin = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _match = a[0]
				var _min = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_match, _min)
				return _range.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _match, _min, func() interface{} {
					c = _isInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _min)}).Values)
					if c.(bool) {
						return mml.Ref(_ints, "max")
					} else {
						return mml.Ref(_float, "max")
					}
				}())}).Values)
			},
			FixedArgs: 2,
		}
		exports["rangeMin"] = _rangeMin
		_listLength = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				return _range.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _any)}).Values), _l, _l)}).Values)
			},
			FixedArgs: 1,
		}
		exports["listLength"] = _listLength
		_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_m)
				return _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _v = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_v)
						return !_is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _m, _v)}).Values).(bool)
					},
					FixedArgs: 1,
				})}).Values)
			},
			FixedArgs: 1,
		}
		exports["not"] = _not
		_natural = _rangeMin.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _int, 0)}).Values)
		exports["natural"] = _natural
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
						return mml.Ref(_functions, "bind").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is, _match)}).Values)
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

	modulePath = "floats"

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

	modulePath = "read"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _readModule interface{}
		var _do interface{}
		var _parse interface{}
		var _errors interface{}
		var _io interface{}
		var _paths interface{}
		var _structs interface{}
		var _codetree interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_readModule, _do, _parse, _errors, _io, _paths, _structs, _codetree, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_parse = mml.Modules.Use("parse")
		_errors = mml.Modules.Use("errors")
		_io = mml.Modules.Use("io")
		_paths = mml.Modules.Use("paths")
		_structs = mml.Modules.Use("structs")
		_codetree = mml.Modules.Use("codetree")
		_readModule = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _reading = a[0]
				var _modules = a[1]
				var _path = a[2]
				var _ interface{}
				_ = &mml.List{a[3:]}
				mml.Nop(_reading, _modules, _path)
				var _moduleCode interface{}
				var _usePaths interface{}
				var _readingUses interface{}
				var _nextModules interface{}
				var _setUsedModule interface{}
				var _withUsedModules interface{}
				mml.Nop(_moduleCode, _usePaths, _readingUses, _nextModules, _setUsedModule, _withUsedModules)
				c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path, _reading)}).Values)
				if c.(bool) {
					mml.Nop()
					return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "circular module reference")}).Values)
				}
				c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path, _modules)}).Values)
				if c.(bool) {
					mml.Nop()
					return _modules
				}
				_moduleCode = mml.Ref(_errors, "pass").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _bind.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats, "%s.mml")}).Values), mml.Ref(_io, "readFile"), mml.Ref(_parse, "do"))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path)}).Values)
				if v := _moduleCode; mml.IsError.F([]interface{}{v}).(bool) {
					return v
				}
				_usePaths = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "path")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_codetree, "filter").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "use"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _moduleCode)}).Values))}).Values))}).Values)
				_readingUses = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _reading.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values[_path.(string)] = true
					return s
				}()
				_nextModules = _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _path = a[0]
						var _modules = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_path, _modules)
						return func() interface{} {
							c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _modules)}).Values)
							if c.(bool) {
								return _modules
							} else {
								return _readModule.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _readingUses, _modules, _path)}).Values)
							}
						}()
					},
					FixedArgs: 2,
				}, _modules, _usePaths)}).Values)
				if v := _nextModules; mml.IsError.F([]interface{}{v}).(bool) {
					return v
				}
				_setUsedModule = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_code)
						return func() interface{} {
							c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								s.Values["type"] = "use"
								return s
							}(), _code)}).Values)
							if c.(bool) {
								return func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									func() {
										sp := _code.(*mml.Struct)
										for k, v := range sp.Values {
											s.Values[k] = v
										}
									}()
									s.Values["module"] = mml.Ref(_nextModules, mml.Ref(mml.Ref(_code, "path"), "value"))
									return s
								}()
							} else {
								return _code
							}
						}()
					},
					FixedArgs: 1,
				}
				_withUsedModules = mml.Ref(_codetree, "edit").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _setUsedModule, _moduleCode)}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _nextModules.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values[_path.(string)] = func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						func() {
							sp := _withUsedModules.(*mml.Struct)
							for k, v := range sp.Values {
								s.Values[k] = v
							}
						}()
						s.Values["path"] = _path
						return s
					}()
					return s
				}()
				return nil
			},
			FixedArgs: 3,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_path)
				return mml.Ref(_errors, "pass").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_paths, "normalize"), mml.Ref(_paths, "trimExtension"), _readModule.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }(), func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }())}).Values), mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path)}).Values)
			},
			FixedArgs: 1,
		}
		exports["do"] = _do

		return exports
	})

	modulePath = "parse"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _assortComments interface{}
		var _create interface{}
		var _functionFact interface{}
		var _rangeExpression interface{}
		var _indexer interface{}
		var _application interface{}
		var _unary interface{}
		var _binary interface{}
		var _chaining interface{}
		var _ternary interface{}
		var _ifStatement interface{}
		var _parseCase interface{}
		var _defaultStatements interface{}
		var _switchStatement interface{}
		var _sendStatement interface{}
		var _receiveExpression interface{}
		var _selectStatement interface{}
		var _rangeOver interface{}
		var _loop interface{}
		var _assign interface{}
		var _valueCapture interface{}
		var _mutableCapture interface{}
		var _valueDefinition interface{}
		var _definitionGroup interface{}
		var _mutableDefinitionGroup interface{}
		var _functionCapture interface{}
		var _effectCapture interface{}
		var _functionDefinition interface{}
		var _effectDefinitionGroup interface{}
		var _exportStatement interface{}
		var _useFact interface{}
		var _parse interface{}
		var _parserError interface{}
		var _knownOrError interface{}
		var _parsePrimitive interface{}
		var _ast interface{}
		var _commentLine interface{}
		var _lineComment interface{}
		var _blockCommentContent interface{}
		var _blockComment interface{}
		var _intCode interface{}
		var _floatCode interface{}
		var _stringCode interface{}
		var _boolCode interface{}
		var _symbol interface{}
		var _spread interface{}
		var _list interface{}
		var _mutableList interface{}
		var _expressionKey interface{}
		var _entry interface{}
		var _struct interface{}
		var _mutableStruct interface{}
		var _ret interface{}
		var _checkRet interface{}
		var _statementListOf interface{}
		var _statementList interface{}
		var _collectParameter interface{}
		var _functionLiteral interface{}
		var _effect interface{}
		var _symbolIndex interface{}
		var _expressionIndex interface{}
		var _rangeIndex interface{}
		var _goStatement interface{}
		var _deferStatement interface{}
		var _useEffect interface{}
		var _useList interface{}
		var _module interface{}
		var _do interface{}
		var _validateast interface{}
		var _structs interface{}
		var _lists interface{}
		var _code interface{}
		var _errors interface{}
		var _codetree interface{}
		var _strings interface{}
		var _functions interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_assortComments, _create, _functionFact, _rangeExpression, _indexer, _application, _unary, _binary, _chaining, _ternary, _ifStatement, _parseCase, _defaultStatements, _switchStatement, _sendStatement, _receiveExpression, _selectStatement, _rangeOver, _loop, _assign, _valueCapture, _mutableCapture, _valueDefinition, _definitionGroup, _mutableDefinitionGroup, _functionCapture, _effectCapture, _functionDefinition, _effectDefinitionGroup, _exportStatement, _useFact, _parse, _parserError, _knownOrError, _parsePrimitive, _ast, _commentLine, _lineComment, _blockCommentContent, _blockComment, _intCode, _floatCode, _stringCode, _boolCode, _symbol, _spread, _list, _mutableList, _expressionKey, _entry, _struct, _mutableStruct, _ret, _checkRet, _statementListOf, _statementList, _collectParameter, _functionLiteral, _effect, _symbolIndex, _expressionIndex, _rangeIndex, _goStatement, _deferStatement, _useEffect, _useList, _module, _do, _validateast, _structs, _lists, _code, _errors, _codetree, _strings, _functions, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_validateast = mml.Modules.Use("validateast")
		_structs = mml.Modules.Use("structs")
		_lists = mml.Modules.Use("lists")
		_code = mml.Modules.Use("code")
		_errors = mml.Modules.Use("errors")
		_codetree = mml.Modules.Use("codetree")
		_strings = mml.Modules.Use("strings")
		_functions = mml.Modules.Use("functions")
		_assortComments = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _isComment interface{}
				var _astStripped interface{}
				var _comments interface{}
				mml.Nop(_isComment, _astStripped, _comments)
				_isComment = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "line-comment", "block-comment")}).Values)
					return s
				}())}).Values)
				_astStripped = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _ast.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["nodes"] = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "not").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isComment)}).Values), mml.Ref(_ast, "nodes"))}).Values)
					return s
				}()
				_comments = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["nodes"] = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isComment, mml.Ref(_ast, "nodes"))}).Values)
					s.Values["indexes"] = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _i = a[0]
							var _ interface{}
							_ = &mml.List{a[1:]}
							mml.Nop(_i)
							return _isComment.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), _i))}).Values)
						},
						FixedArgs: 1,
					})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "indexes").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values))}).Values)
					return s
				}()
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["ast"] = _astStripped
					s.Values["comments"] = _comments
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_create = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _type = a[0]
				var _ast = a[1]
				var _props interface{}
				_props = &mml.List{a[2:]}
				mml.Nop(_type, _ast, _props)
				return mml.Ref(_structs, "merges").(*mml.Function).Call((&mml.List{Values: append(append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = _type
					s.Values["ast"] = _ast
					return s
				}()), _props.(*mml.List).Values...)}).Values)
			},
			FixedArgs: 2,
		}
		_commentLine = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "comment-line", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["text"] = mml.Ref(_ast, "text")
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_lineComment = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "line-comment", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["lines"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_blockCommentContent = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "block-comment-content", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["text"] = mml.Ref(_ast, "text")
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_blockComment = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "block-comment", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["content"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_intCode = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "int")}).Values)
		_floatCode = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "float")}).Values)
		_stringCode = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "string")}).Values)
		_boolCode = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "bool")}).Values)
		_symbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "symbol", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = mml.Ref(_ast, "text")
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_spread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "spread", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_list = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "list", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["values"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					s.Values["mutable"] = false
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_mutableList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression-key", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_entry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "entry", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["key"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_struct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "struct", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["entries"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					s.Values["mutable"] = false
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_mutableStruct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_ret = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "ret", _ast, func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 0)
					if c.(bool) {
						return func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }()
					} else {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
							return s
						}()
					}
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_checkRet = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "check-ret", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_statementListOf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _nodes = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_ast, _nodes)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "statement-list", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["statements"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, _nodes)}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 2,
		}
		_statementList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _statementListOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast, mml.Ref(_ast, "nodes"))}).Values)
			},
			FixedArgs: 1,
		}
		_collectParameter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
			},
			FixedArgs: 1,
		}
		_functionFact = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _offset = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_ast, _offset)
				var _nodes interface{}
				var _last interface{}
				var _params interface{}
				var _lastParam interface{}
				var _hasCollectParam interface{}
				var _fixedParams interface{}
				mml.Nop(_nodes, _last, _params, _lastParam, _hasCollectParam, _fixedParams)
				_nodes = mml.RefRange(mml.Ref(_ast, "nodes"), _offset, nil)
				_last = mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 1)
				_params = mml.RefRange(_nodes, nil, _last)
				_lastParam = mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values), 1)
				_hasCollectParam = (mml.BinaryOp(16, _lastParam, 0).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_params, _lastParam), "name"), "collect-parameter").(bool))
				_fixedParams = func() interface{} {
					c = _hasCollectParam
					if c.(bool) {
						return mml.RefRange(_params, nil, _lastParam)
					} else {
						return _params
					}
				}()
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "function", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["params"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "name")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fixedParams)}).Values))}).Values)
					s.Values["collectParam"] = func() interface{} {
						c = _hasCollectParam
						if c.(bool) {
							return mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_params, _lastParam))}).Values), "name")
						} else {
							return ""
						}
					}()
					s.Values["body"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, _last))}).Values)
					s.Values["effect"] = false
					return s
				}())}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_functionLiteral = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _functionFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast, 0)}).Values)
			},
			FixedArgs: 1,
		}
		_effect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _functionFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast, 0)}).Values).(*mml.Struct)
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
		_rangeExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "range", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values[func() interface{} {
						c = mml.BinaryOp(11, mml.Ref(_ast, "name"), "range-from")
						if c.(bool) {
							return "from"
						} else {
							return "to"
						}
					}().(string)] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_symbolIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "symbol-index", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["symbol"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_expressionIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
			},
			FixedArgs: 1,
		}
		_rangeIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append(append([]interface{}{}, "range", _ast), _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values).(*mml.List).Values...)}).Values)
			},
			FixedArgs: 1,
		}
		_indexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _indexerNodes interface{}
				mml.Nop(_indexerNodes)
				_indexerNodes = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _nodes = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_nodes)
						return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "indexer", _ast, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["expression"] = func() interface{} {
								c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 2)
								if c.(bool) {
									return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, 0))}).Values)
								} else {
									return _indexerNodes.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(_nodes, nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 1)))}).Values)
								}
							}()
							s.Values["index"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_nodes, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values), 1)))}).Values)
							return s
						}())}).Values)
					},
					FixedArgs: 1,
				}
				return _indexerNodes.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_application = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "application", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["function"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["args"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_unary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _ops interface{}
				mml.Nop(_ops)
				_ops = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["binary-not"] = mml.Ref(_code, "binaryNot")
					s.Values["plus"] = mml.Ref(_code, "plus")
					s.Values["minus"] = mml.Ref(_code, "minus")
					s.Values["logical-not"] = mml.Ref(_code, "logicalNot")
					return s
				}()
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "unary", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["op"] = mml.Ref(_ops, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"))
					s.Values["arg"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}())}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_binary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _ops interface{}
				mml.Nop(_ops)
				_ops = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["binary-and"] = mml.Ref(_code, "binaryAnd")
					s.Values["xor"] = mml.Ref(_code, "xor")
					s.Values["and-not"] = mml.Ref(_code, "andNot")
					s.Values["lshift"] = mml.Ref(_code, "lshift")
					s.Values["rshift"] = mml.Ref(_code, "rshift")
					s.Values["mul"] = mml.Ref(_code, "mul")
					s.Values["div"] = mml.Ref(_code, "div")
					s.Values["mod"] = mml.Ref(_code, "mod")
					s.Values["add"] = mml.Ref(_code, "add")
					s.Values["sub"] = mml.Ref(_code, "sub")
					s.Values["eq"] = mml.Ref(_code, "equals")
					s.Values["not-eq"] = mml.Ref(_code, "notEq")
					s.Values["less"] = mml.Ref(_code, "less")
					s.Values["less-or-eq"] = mml.Ref(_code, "lessOrEq")
					s.Values["greater"] = mml.Ref(_code, "greater")
					s.Values["greater-or-eq"] = mml.Ref(_code, "greaterOrEq")
					s.Values["logical-and"] = mml.Ref(_code, "logicalAnd")
					s.Values["logical-or"] = mml.Ref(_code, "logicalOr")
					return s
				}()
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "binary", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["op"] = mml.Ref(_ops, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 2)), "name"))
					s.Values["left"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 3)
						if c.(bool) {
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _ast.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["nodes"] = mml.RefRange(mml.Ref(_ast, "nodes"), nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 2))
								return s
							}()
						} else {
							return mml.Ref(mml.Ref(_ast, "nodes"), 0)
						}
					}())}).Values)
					s.Values["right"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)))}).Values)
					return s
				}())}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_chaining = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fold.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]
						var _a = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_f, _a)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _ast.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["name"] = "application"
							s.Values["nodes"] = &mml.List{Values: append([]interface{}{}, _f, _a)}
							return s
						}()
					},
					FixedArgs: 2,
				}, mml.Ref(mml.Ref(_ast, "nodes"), 0), mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_ternary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "cond", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["condition"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["consequent"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					s.Values["alternative"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 2))}).Values)
					s.Values["ternary"] = true
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_ifStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _constructCond interface{}
				mml.Nop(_constructCond)
				_constructCond = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _cond = a[0]
						var _cons = a[1]
						var _alt = a[2]
						var _ interface{}
						_ = &mml.List{a[3:]}
						mml.Nop(_cond, _cons, _alt)
						return func() interface{} {
							c = mml.BinaryOp(11, _alt, false)
							if c.(bool) {
								return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "cond", _ast, func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									s.Values["condition"] = _cond
									s.Values["consequent"] = _cons
									s.Values["ternary"] = false
									return s
								}())}).Values)
							} else {
								return func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									func() {
										sp := _constructCond.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _cond, _cons, false)}).Values).(*mml.Struct)
										for k, v := range sp.Values {
											s.Values[k] = v
										}
									}()
									s.Values["alternative"] = _alt
									return s
								}()
							}
						}()
					},
					FixedArgs: 3,
				}
				return _foldr.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _g = a[0]
						var _i = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_g, _i)
						return func() interface{} {
							c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _g)}).Values), 1)
							if c.(bool) {
								return mml.Ref(_g, 0)
							} else {
								return _constructCond.(*mml.Function).Call((&mml.List{Values: append(append([]interface{}{}, _g.(*mml.List).Values...), _i)}).Values)
							}
						}()
					},
					FixedArgs: 2,
				}, false)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "group").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_parseCase = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _name = a[0]
				var _ast = a[1]
				var _c = a[2]
				var _ interface{}
				_ = &mml.List{a[3:]}
				mml.Nop(_name, _ast, _c)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _name, _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_c, "nodes"), 0))}).Values)
					s.Values["body"] = _statementListOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast, mml.RefRange(mml.Ref(_c, "nodes"), 1, nil))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 3,
		}
		_defaultStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _statementListOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "nodes")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = "default-block"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values))}).Values))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_switchStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _hasExpression interface{}
				var _cases interface{}
				var _expression interface{}
				var _defaults interface{}
				mml.Nop(_hasExpression, _cases, _expression, _defaults)
				_hasExpression = (mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 0).(bool) && !_is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "case-block", "default-block")}).Values)
					return s
				}(), mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values).(bool))
				_expression = func() interface{} {
					c = _hasExpression
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
							return s
						}()
					} else {
						return func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }()
					}
				}()
				_defaults = _defaultStatements.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				_cases = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parseCase.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "switch-case", _ast)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = "case-block"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values))}).Values)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "switch-statement", _ast, _expression, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["cases"] = _cases
					return s
				}(), func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["defaultStatements"] = _defaults
					return s
				}())}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_sendStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "send-statement", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["channel"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_receiveExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "receive-expression", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["channel"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_selectStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _cases interface{}
				var _hasDefault interface{}
				var _defaults interface{}
				mml.Nop(_cases, _hasDefault, _defaults)
				_hasDefault = _some.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = "default-block"
					return s
				}())}).Values), mml.Ref(_ast, "nodes"))}).Values)
				_defaults = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["hasDefault"] = _hasDefault
					s.Values["defaultStatements"] = _defaultStatements.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
					return s
				}()
				_cases = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parseCase.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "select-case", _ast)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = "select-case-block"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values))}).Values)
				return _create.(*mml.Function).Call((&mml.List{Values: append(append(append([]interface{}{}, "select-statement", _ast), _cases.(*mml.List).Values...), _defaults)}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_goStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "go-statement", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["application"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_deferStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "defer-statement", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["application"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_rangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _createRangeOver interface{}
				var _parseExpression interface{}
				mml.Nop(_createRangeOver, _parseExpression)
				_createRangeOver = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _props interface{}
						_props = &mml.List{a[0:]}
						mml.Nop(_props)
						return _create.(*mml.Function).Call((&mml.List{Values: append(append([]interface{}{}, "range-over", _ast), _props.(*mml.List).Values...)}).Values)
					},
					FixedArgs: 0,
				}
				_parseExpression = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _nodes = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_nodes)
						return mml.Ref(_structs, "merge").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values))}).Values)
					},
					FixedArgs: 1,
				}
				switch {
				case mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 0):

					mml.Nop()
					return _createRangeOver.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				case (mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "symbol").(bool)):

					mml.Nop()
					return _createRangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
						return s
					}())}).Values)
				case mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "symbol"):

					mml.Nop()
					return _createRangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["expression"] = _parseExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values)
						return s
					}())}).Values)
				default:

					mml.Nop()
					return _createRangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
						s.Values["expression"] = _parseExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil))}).Values)
						return s
					}())}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_loop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _createLoop interface{}
				var _emptyRange interface{}
				var _expression interface{}
				var _loop interface{}
				mml.Nop(_createLoop, _emptyRange, _expression, _loop)
				_createLoop = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _body = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_body)
						return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "loop", _ast, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["body"] = _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _body)}).Values)
							return s
						}())}).Values)
					},
					FixedArgs: 1,
				}
				c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)
				if c.(bool) {
					mml.Nop()
					return _createLoop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				}
				_emptyRange = _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "range-over"
					return s
				}(), _not.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["symbol"] = _any
					return s
				}(), func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["expression"] = _any
					return s
				}())}).Values))}).Values))}).Values)
				_expression = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				_loop = _createLoop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
				return func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _emptyRange, _expression)}).Values)
					if c.(bool) {
						return _loop
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
							return s
						}()
					}
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_assign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "assign", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["capture"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					s.Values["value"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_valueCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definition", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
					s.Values["expression"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
					s.Values["mutable"] = false
					s.Values["exported"] = false
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_mutableCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return func() interface{} {
					c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["docs"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
							func() {
								sp := _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values).(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							return s
						}()
					} else {
						return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_definitionGroup = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definition-group", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["definitions"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_mutableDefinitionGroup = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _d interface{}
				mml.Nop(_d)
				_d = _definitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _d.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["definitions"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _d = a[0]
							var _ interface{}
							_ = &mml.List{a[1:]}
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
					}, mml.Ref(_d, "definitions"))}).Values)
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definition", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["symbol"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
					s.Values["expression"] = _functionFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast, 1)}).Values)
					s.Values["mutable"] = false
					s.Values["exported"] = false
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_effectCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_functionDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return func() interface{} {
					c = mml.BinaryOp(15, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_ast, "nodes"))}).Values), 1)
					if c.(bool) {
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["docs"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
							func() {
								sp := _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values).(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							return s
						}()
					} else {
						return _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_effectDefinitionGroup = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _d interface{}
				mml.Nop(_d)
				_d = _definitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _d.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["definitions"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _d = a[0]
							var _ interface{}
							_ = &mml.List{a[1:]}
							mml.Nop(_d)
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _d.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["expression"] = func() interface{} {
									s := &mml.Struct{Values: make(map[string]interface{})}
									func() {
										sp := mml.Ref(_d, "expression").(*mml.Struct)
										for k, v := range sp.Values {
											s.Values[k] = v
										}
									}()
									s.Values["effect"] = true
									return s
								}()
								return s
							}()
						},
						FixedArgs: 1,
					}, mml.Ref(_d, "definitions"))}).Values)
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_exportStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _d interface{}
				var _dl interface{}
				var _edl interface{}
				mml.Nop(_d, _dl, _edl)
				_d = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
				_dl = func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_d, "type"), "definition")
					if c.(bool) {
						return &mml.List{Values: append([]interface{}{}, _d)}
					} else {
						return mml.Ref(_d, "definitions")
					}
				}()
				_edl = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _d = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
				}, _dl)}).Values)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definition-group", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["definitions"] = _edl
					return s
				}())}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_useFact = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _createUse interface{}
				mml.Nop(_createUse)
				_createUse = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _props interface{}
						_props = &mml.List{a[0:]}
						mml.Nop(_props)
						return _create.(*mml.Function).Call((&mml.List{Values: append(append([]interface{}{}, "use", _ast, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["effect"] = false
							return s
						}()), _props.(*mml.List).Values...)}).Values)
					},
					FixedArgs: 0,
				}
				switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name") {
				case "use-inline":

					mml.Nop()
					return _createUse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["capture"] = "."
						s.Values["path"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
						return s
					}())}).Values)
				case "symbol":

					mml.Nop()
					return _createUse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["capture"] = mml.Ref(_parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values), "name")
						s.Values["path"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))}).Values)
						return s
					}())}).Values)
				default:

					mml.Nop()
					return _createUse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["path"] = _parse.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))}).Values)
						return s
					}())}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_useEffect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _useFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values).(*mml.Struct)
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
		_useList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "use-list", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["uses"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(_ast, "nodes"))}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_module = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				return _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "module", _ast, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["body"] = _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
					return s
				}())}).Values)
			},
			FixedArgs: 1,
		}
		_parse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_ast)
				var _a interface{}
				var _code interface{}
				mml.Nop(_a, _code)
				switch mml.Ref(_ast, "name") {
				case "line-comment-content":

					mml.Nop()
					return _commentLine.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "line-comment":

					mml.Nop()
					return _lineComment.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "block-comment-content":

					mml.Nop()
					return _blockCommentContent.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "block-comment":

					mml.Nop()
					return _blockComment.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "int":

					mml.Nop()
					return _intCode.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "float":

					mml.Nop()
					return _floatCode.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "string":

					mml.Nop()
					return _stringCode.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "true":

					mml.Nop()
					return _boolCode.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "false":

					mml.Nop()
					return _boolCode.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				case "symbol":

					mml.Nop()
					return _symbol.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				}
				_a = _assortComments.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _ast)}).Values)
				_code = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "unknown", mml.Ref(_a, "ast"))}).Values)
				switch mml.Ref(mml.Ref(_a, "ast"), "name") {
				case "spread":

					mml.Nop()
					_code = _spread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "list":

					mml.Nop()
					_code = _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "mutable-list":

					mml.Nop()
					_code = _mutableList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "expression-key":

					mml.Nop()
					_code = _expressionKey.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "entry":

					mml.Nop()
					_code = _entry.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "struct":

					mml.Nop()
					_code = _struct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "mutable-struct":

					mml.Nop()
					_code = _mutableStruct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "ret":

					mml.Nop()
					_code = _ret.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "check-ret":

					mml.Nop()
					_code = _checkRet.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "block":

					mml.Nop()
					_code = _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "collect-parameter":

					mml.Nop()
					_code = _collectParameter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "function":

					mml.Nop()
					_code = _functionLiteral.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "effect":

					mml.Nop()
					_code = _effect.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "range-from":

					mml.Nop()
					_code = _rangeExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "range-to":

					mml.Nop()
					_code = _rangeExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "symbol-index":

					mml.Nop()
					_code = _symbolIndex.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "expression-index":

					mml.Nop()
					_code = _expressionIndex.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "range-index":

					mml.Nop()
					_code = _rangeIndex.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "indexer":

					mml.Nop()
					_code = _indexer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "application":

					mml.Nop()
					_code = _application.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "unary":

					mml.Nop()
					_code = _unary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "binary0":

					mml.Nop()
					_code = _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "binary1":

					mml.Nop()
					_code = _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "binary2":

					mml.Nop()
					_code = _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "binary3":

					mml.Nop()
					_code = _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "binary4":

					mml.Nop()
					_code = _binary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "chaining":

					mml.Nop()
					_code = _chaining.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "ternary":

					mml.Nop()
					_code = _ternary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "if-statement":

					mml.Nop()
					_code = _ifStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "switch-statement":

					mml.Nop()
					_code = _switchStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "send-statement":

					mml.Nop()
					_code = _sendStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "receive-expression":

					mml.Nop()
					_code = _receiveExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "receive-definition":

					mml.Nop()
					_code = _valueCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "select-statement":

					mml.Nop()
					_code = _selectStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "go-statement":

					mml.Nop()
					_code = _goStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "defer-statement":

					mml.Nop()
					_code = _deferStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "range-over":

					mml.Nop()
					_code = _rangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "break":

					mml.Nop()
					_code = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "break", mml.Ref(_a, "ast"))}).Values)
				case "continue":

					mml.Nop()
					_code = _create.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "continue", mml.Ref(_a, "ast"))}).Values)
				case "loop":

					mml.Nop()
					_code = _loop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "assign":

					mml.Nop()
					_code = _assign.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "value-capture":

					mml.Nop()
					_code = _valueCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "mutable-capture":

					mml.Nop()
					_code = _mutableCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "value-definition":

					mml.Nop()
					_code = _valueDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "docs-value-capture":

					mml.Nop()
					_code = _valueDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "docs-mixed-capture":

					mml.Nop()
					_code = _valueDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "value-definition-group":

					mml.Nop()
					_code = _definitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "mutable-definition-group":

					mml.Nop()
					_code = _mutableDefinitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "function-capture":

					mml.Nop()
					_code = _functionCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "effect-capture":

					mml.Nop()
					_code = _effectCapture.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "function-definition":

					mml.Nop()
					_code = _functionDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "docs-function-capture":

					mml.Nop()
					_code = _functionDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "docs-mixed-function-capture":

					mml.Nop()
					_code = _functionDefinition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "function-definition-group":

					mml.Nop()
					_code = _definitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "effect-definition-group":

					mml.Nop()
					_code = _effectDefinitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "export-statement":

					mml.Nop()
					_code = _exportStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "use-fact":

					mml.Nop()
					_code = _useFact.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "use-effect":

					mml.Nop()
					_code = _useEffect.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "use-modules":

					mml.Nop()
					_code = _useList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				case "mml":

					mml.Nop()
					_code = _module.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_a, "ast"))}).Values)
				}
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _code.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["comments"] = func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["code"] = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parse, mml.Ref(mml.Ref(_a, "comments"), "nodes"))}).Values)
						s.Values["indexes"] = mml.Ref(mml.Ref(_a, "comments"), "indexes")
						return s
					}()
					return s
				}()
				return nil
			},
			FixedArgs: 1,
		}
		_parserError = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _msg = a[0]
				var _ast = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_msg, _ast)
				return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s:%d:%d:%v", mml.Ref(_ast, "file"), mml.Ref(_ast, "line"), mml.Ref(_ast, "column"), _msg)}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		_knownOrError = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_code)
				return func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = _not.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "unknown")}).Values)
						return s
					}(), _code)}).Values)
					if c.(bool) {
						return _code
					} else {
						return _parserError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "unknown code", mml.Ref(_code, "ast"))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_parsePrimitive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_code)

				mml.Nop()
				switch mml.Ref(_code, "type") {
				case "int":
					var _v interface{}
					mml.Nop(_v)
					_v = _parseInt.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_code, "ast"), "text"))}).Values)
					return func() interface{} {
						c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _v)}).Values)
						if c.(bool) {
							return _parserError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _v, mml.Ref(_code, "ast"))}).Values)
						} else {
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _code.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["value"] = _v
								return s
							}()
						}
					}()
				case "float":
					var _v interface{}
					mml.Nop(_v)
					_v = _parseFloat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_code, "ast"), "text"))}).Values)
					return func() interface{} {
						c = _isError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _v)}).Values)
						if c.(bool) {
							return _parserError.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _v, mml.Ref(_code, "ast"))}).Values)
						} else {
							return func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								func() {
									sp := _code.(*mml.Struct)
									for k, v := range sp.Values {
										s.Values[k] = v
									}
								}()
								s.Values["value"] = _v
								return s
							}()
						}
					}()
				case "string":

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						func() {
							sp := _code.(*mml.Struct)
							for k, v := range sp.Values {
								s.Values[k] = v
							}
						}()
						s.Values["value"] = mml.Ref(_strings, "unescape").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.RefRange(mml.Ref(mml.Ref(_code, "ast"), "text"), 1, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_code, "ast"), "text"))}).Values), 1)))}).Values)
						return s
					}()
				case "bool":

					mml.Nop()
					return func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						func() {
							sp := _code.(*mml.Struct)
							for k, v := range sp.Values {
								s.Values[k] = v
							}
						}()
						s.Values["value"] = mml.BinaryOp(11, mml.Ref(mml.Ref(_code, "ast"), "text"), "true")
						return s
					}()
				default:

					mml.Nop()
					return _code
				}
				return nil
			},
			FixedArgs: 1,
		}
		_ast = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _node = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_node)
				return mml.Ref(_errors, "pass").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_validateast, "do"), _parse, mml.Ref(_codetree, "edit").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parsePrimitive)}).Values), mml.Ref(_codetree, "edit").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _knownOrError)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _node)}).Values)
			},
			FixedArgs: 1,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _text = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_text)
				return mml.Ref(_errors, "pass").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _parseAST, _ast)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _text)}).Values)
			},
			FixedArgs: 1,
		}
		exports["do"] = _do

		return exports
	})

	modulePath = "validateast"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _validateComments interface{}
		var _dropComments interface{}
		var _rangeExpression interface{}
		var _functionParamsAndBody interface{}
		var _rangeOver interface{}
		var _startsWithCaseOrDefault interface{}
		var _functionCapture interface{}
		var _definitionChild interface{}
		var _stringOrNamedStringOrInline interface{}
		var _customValidators interface{}
		var _validateCustom interface{}
		var _node interface{}
		var _minTextLength interface{}
		var _childCount interface{}
		var _minChildCount interface{}
		var _paramsAreSymbols interface{}
		var _onlyLastParamIsCollect interface{}
		var _textLengthMin2 interface{}
		var _oneChild interface{}
		var _twoChildren interface{}
		var _threeChildren interface{}
		var _minOneChild interface{}
		var _minTwoChildren interface{}
		var _minThreeChildren interface{}
		var _symbol interface{}
		var _stringNode interface{}
		var _useInline interface{}
		var _symbolChild interface{}
		var _collectParameter interface{}
		var _rangeFrom interface{}
		var _rangeTo interface{}
		var _symbolAndAny interface{}
		var _comment interface{}
		var _do interface{}
		var _code interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_validateComments, _dropComments, _rangeExpression, _functionParamsAndBody, _rangeOver, _startsWithCaseOrDefault, _functionCapture, _definitionChild, _stringOrNamedStringOrInline, _customValidators, _validateCustom, _node, _minTextLength, _childCount, _minChildCount, _paramsAreSymbols, _onlyLastParamIsCollect, _textLengthMin2, _oneChild, _twoChildren, _threeChildren, _minOneChild, _minTwoChildren, _minThreeChildren, _symbol, _stringNode, _useInline, _symbolChild, _collectParameter, _rangeFrom, _rangeTo, _symbolAndAny, _comment, _do, _code, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_code = mml.Modules.Use("code")
		_minTextLength = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["text"] = _rangeMin.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _string, _n)}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_childCount = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["nodes"] = _range.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _any)}).Values), _n, _n)}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_minChildCount = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["nodes"] = _rangeMin.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _any)}).Values), _n)}).Values)
					return s
				}()
			},
			FixedArgs: 1,
		}
		_paramsAreSymbols = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _symbol, _collectParameter)}).Values))}).Values), mml.RefRange(mml.Ref(_n, "nodes"), nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, "nodes"))}).Values), 1)))}).Values)
			},
			FixedArgs: 1,
		}
		_onlyLastParamIsCollect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return (mml.BinaryOp(13, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, "nodes"))}).Values), 2).(bool) || _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _symbol)}).Values), mml.RefRange(mml.Ref(_n, "nodes"), nil, mml.BinaryOp(10, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, "nodes"))}).Values), 2)))}).Values).(bool))
			},
			FixedArgs: 1,
		}
		_textLengthMin2 = _minTextLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values)
		_oneChild = _childCount.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values)
		_twoChildren = _childCount.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values)
		_threeChildren = _childCount.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 3)}).Values)
		_minOneChild = _minChildCount.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values)
		_minTwoChildren = _minChildCount.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values)
		_minThreeChildren = _minChildCount.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 3)}).Values)
		_symbol = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = "symbol"
			return s
		}()
		_stringNode = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = "string"
			return s
		}()
		_useInline = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = "use-inline"
			return s
		}()
		_symbolChild = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = &mml.List{Values: append([]interface{}{}, _symbol)}
			return s
		}()
		_collectParameter = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = "collect-parameter"
			func() {
				sp := _symbolChild.(*mml.Struct)
				for k, v := range sp.Values {
					s.Values[k] = v
				}
			}()
			return s
		}()
		_rangeFrom = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = "range-from"
			return s
		}()
		_rangeTo = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = "range-to"
			return s
		}()
		_symbolAndAny = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = &mml.List{Values: append([]interface{}{}, _symbol, _any)}
			return s
		}()
		_comment = _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "line-comment", "block-comment")}).Values)
			return s
		}())}).Values))}).Values)
		_validateComments = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _node = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_node)
				return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = "block-comment"
					return s
				}(), _oneChild)}).Values), _not.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = "block-comment"
					return s
				}())}).Values))}).Values), func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["nodes"] = _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _validateComments)}).Values))}).Values)
					return s
				}())}).Values), _node)}).Values)
			},
			FixedArgs: 1,
		}
		_dropComments = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_nodes)
				return _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _n = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_n)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _n.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values["nodes"] = _dropComments.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, "nodes"))}).Values)
							return s
						}()
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _not.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _comment)}).Values))}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _nodes)}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_rangeExpression = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 0)}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values), &mml.List{Values: append([]interface{}{}, _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _rangeFrom, _rangeTo)}).Values))})}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values), &mml.List{Values: append([]interface{}{}, _rangeFrom, _rangeTo)})}).Values))}).Values)
			return s
		}()
		_functionParamsAndBody = _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _minOneChild, _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _paramsAreSymbols)}).Values), _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _onlyLastParamIsCollect)}).Values))}).Values)
		_rangeOver = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 0)}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values), &mml.List{Values: append([]interface{}{}, _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _rangeFrom, _rangeTo)}).Values))})}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values), &mml.List{Values: append([]interface{}{}, _rangeFrom, _rangeTo)})}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 1)}).Values), &mml.List{Values: append([]interface{}{}, _symbol)})}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 2)}).Values), &mml.List{Values: append([]interface{}{}, _symbol, _any)})}).Values), _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listLength.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, 3)}).Values), &mml.List{Values: append([]interface{}{}, _symbol, _rangeFrom, _rangeTo)})}).Values))}).Values)
			return s
		}()
		_startsWithCaseOrDefault = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append([]interface{}{}, func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "case-line", "default-line")}).Values)
			return s
		}())}, &mml.List{Values: append([]interface{}{}, _any, func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["name"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "case-line", "default-line")}).Values)
			return s
		}())})}).Values)
		_functionCapture = _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = &mml.List{Values: append([]interface{}{}, _symbol)}
			return s
		}(), _functionParamsAndBody)}).Values)
		_definitionChild = _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _oneChild, func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = &mml.List{Values: append([]interface{}{}, func() interface{} {
				s := &mml.Struct{Values: make(map[string]interface{})}
				s.Values["name"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value-definition", "value-definition-group", "mutable-definition-group", "function-definition", "function-definition-group", "effect-definition-group")}).Values)
				return s
			}())}
			return s
		}())}).Values)
		_stringOrNamedStringOrInline = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["nodes"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append([]interface{}{}, _stringNode)}, &mml.List{Values: append([]interface{}{}, _symbol, _stringNode)}, &mml.List{Values: append([]interface{}{}, _useInline, _stringNode)})}).Values)
			return s
		}()
		_customValidators = func() interface{} {
			s := &mml.Struct{Values: make(map[string]interface{})}
			s.Values["block-comment"] = _oneChild
			s.Values["string"] = _textLengthMin2
			s.Values["symbol"] = _symbol
			s.Values["spread"] = _oneChild
			s.Values["expression-key"] = _oneChild
			s.Values["entry"] = _twoChildren
			s.Values["check-ret"] = _oneChild
			s.Values["function"] = _functionParamsAndBody
			s.Values["effect"] = _functionParamsAndBody
			s.Values["range-from"] = _oneChild
			s.Values["range-to"] = _oneChild
			s.Values["symbol-index"] = _symbolChild
			s.Values["range-index"] = _rangeExpression
			s.Values["indexer"] = _minTwoChildren
			s.Values["application"] = _minOneChild
			s.Values["unary"] = _twoChildren
			s.Values["binary"] = _minThreeChildren
			s.Values["chaining"] = _minTwoChildren
			s.Values["ternary"] = _threeChildren
			s.Values["if-statement"] = _minTwoChildren
			s.Values["case-block"] = _minOneChild
			s.Values["select-case-block"] = _minOneChild
			s.Values["range-over"] = _rangeOver
			s.Values["value-capture"] = _symbolAndAny
			s.Values["mutable-capture"] = _symbolAndAny
			s.Values["function-capture"] = _functionCapture
			s.Values["effect-capture"] = _functionCapture
			s.Values["assign"] = _twoChildren
			s.Values["send-statement"] = _twoChildren
			s.Values["receive-statement"] = _oneChild
			s.Values["go-statement"] = _oneChild
			s.Values["defer-statement"] = _oneChild
			s.Values["receive-definition"] = _symbolAndAny
			s.Values["export-statement"] = _definitionChild
			s.Values["use-fact"] = _stringOrNamedStringOrInline
			return s
		}()
		_validateCustom = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return (!_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, "name"), _customValidators)}).Values).(bool) || _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_customValidators, mml.Ref(_n, "name")), _n)}).Values).(bool))
			},
			FixedArgs: 1,
		}
		_node = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _and.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["name"] = _type.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _string)}).Values)
					s.Values["nodes"] = _listOf.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _node)}).Values))}).Values)
					s.Values["text"] = _type.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _string)}).Values)
					s.Values["file"] = _type.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _string)}).Values)
					s.Values["from"] = _natural
					s.Values["to"] = _natural
					s.Values["line"] = _natural
					s.Values["column"] = _natural
					return s
				}(), _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _validateCustom)}).Values))}).Values), _n)}).Values)
			},
			FixedArgs: 1,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_n)
				var _nc interface{}
				mml.Nop(_nc)
				_nc = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _n.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values["nodes"] = _dropComments.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_n, "nodes"))}).Values)
					return s
				}()
				return func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _node)}).Values), _nc)}).Values)
					if c.(bool) {
						return _n
					} else {
						return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "invalid AST")}).Values)
					}
				}()
				return nil
			},
			FixedArgs: 1,
		}
		exports["do"] = _do

		return exports
	})

	modulePath = "code"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _keywords interface{}
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
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_keywords, _controlStatement, _breakControl, _continueControl, _unaryOp, _binaryNot, _plus, _minus, _logicalNot, _binaryOp, _binaryAnd, _binaryOr, _xor, _andNot, _lshift, _rshift, _mul, _div, _mod, _add, _sub, _equals, _notEq, _less, _lessOrEq, _greater, _greaterOrEq, _logicalAnd, _logicalOr, _builtin, _flattenedStatements, _getModuleName, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_keywords = &mml.List{Values: append([]interface{}{}, "true", "false", "return", "fn", "if", "else", "case", "switch", "default", "send", "receive", "select", "go", "defer", "in", "for", "let", "use", "export")}
		exports["keywords"] = _keywords
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
				var _ interface{}
				_ = &mml.List{a[4:]}
				mml.Nop(_itemType, _listType, _listProp, _statements)
				var _type interface{}
				var _toList interface{}
				mml.Nop(_type, _toList)
				_type = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
						var _ interface{}
						_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_path)
				return _path
			},
			FixedArgs: 1,
		}
		exports["getModuleName"] = _getModuleName

		return exports
	})

	modulePath = "structs"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _merge interface{}
		var _merges interface{}
		var _get interface{}
		var _set interface{}
		var _index interface{}
		var _values interface{}
		var _filterByKeys interface{}
		var _map interface{}
		var _lists interface{}
		mml.Nop(_merge, _merges, _get, _set, _index, _values, _filterByKeys, _map, _lists)
		_lists = mml.Modules.Use("lists")
		_merge = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return mml.Ref(_lists, "fold").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _s = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_i, _s)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _s.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							func() {
								sp := _i.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							return s
						}()
					},
					FixedArgs: 2,
				}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }(), _s)}).Values)
			},
			FixedArgs: 1,
		}
		exports["merge"] = _merge
		_merges = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s interface{}
				_s = &mml.List{a[0:]}
				mml.Nop(_s)
				return _merge.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values)
			},
			FixedArgs: 0,
		}
		exports["merges"] = _merges
		_get = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _key = a[0]
				var _s = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_key, _s)
				return mml.Ref(_s, _key)
			},
			FixedArgs: 2,
		}
		exports["get"] = _get
		_set = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _key = a[0]
				var _value = a[1]
				var _s = a[2]
				var _ interface{}
				_ = &mml.List{a[3:]}
				mml.Nop(_key, _value, _s)
				return func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _s.(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					s.Values[_key.(string)] = _value
					return s
				}()
			},
			FixedArgs: 3,
		}
		exports["set"] = _set
		_index = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _key = a[0]
				var _s = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_key, _s)
				return mml.Ref(_lists, "fold").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _s = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_i, _s)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _s.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values[mml.Ref(_i, _key).(string)] = _i
							return s
						}()
					},
					FixedArgs: 2,
				}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }(), _s)}).Values)
			},
			FixedArgs: 2,
		}
		exports["index"] = _index
		_values = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _key = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_key)
						return mml.Ref(_s, _key)
					},
					FixedArgs: 1,
				}, _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		exports["values"] = _values
		_filterByKeys = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _k = a[0]
				var _s = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_k, _s)
				return mml.Ref(_lists, "fold").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ki = a[0]
						var _f = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_ki, _f)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _f.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							s.Values[_ki.(string)] = mml.Ref(_s, _ki)
							return s
						}()
					},
					FixedArgs: 2,
				}, func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }())}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "intersect").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _k)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		exports["filterByKeys"] = _filterByKeys
		_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _s = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_f, _s)
				return _merge.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _key = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_key)
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["key"] = _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, _key))}).Values)
							return s
						}()
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _s)}).Values))}).Values))}).Values)
			},
			FixedArgs: 2,
		}
		exports["map"] = _map

		return exports
	})

	modulePath = "errors"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _only interface{}
		var _pass interface{}
		var _any interface{}
		var _lists interface{}
		var _functions interface{}
		mml.Nop(_only, _pass, _any, _lists, _functions)
		_lists = mml.Modules.Use("lists")
		_functions = mml.Modules.Use("functions")
		_only = mml.Ref(_functions, "bind").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "only"), _isError)}).Values)
		exports["only"] = _only
		_pass = mml.Ref(_functions, "bind").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "only"), mml.Ref(_functions, "not").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isError)}).Values))}).Values)
		exports["pass"] = _pass
		_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				var _first interface{}
				mml.Nop(_first)
				_first = mml.Ref(_lists, "first").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _isError, _l)}).Values)
				return func() interface{} {
					c = mml.BinaryOp(11, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _first)}).Values), 0)
					if c.(bool) {
						return _l
					} else {
						return mml.Ref(_first, 0)
					}
				}()
				return nil
			},
			FixedArgs: 1,
		}
		exports["any"] = _any

		return exports
	})

	modulePath = "codetree"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _removeToken interface{}
		var _callTransform interface{}
		var _do interface{}
		var _edit interface{}
		var _filter interface{}
		var _trim interface{}
		var _structs interface{}
		var _lists interface{}
		var _functions interface{}
		var _errors interface{}
		mml.Nop(_removeToken, _callTransform, _do, _edit, _filter, _trim, _structs, _lists, _functions, _errors)
		_structs = mml.Modules.Use("structs")
		_lists = mml.Modules.Use("lists")
		_functions = mml.Modules.Use("functions")
		_errors = mml.Modules.Use("errors")
		_removeToken = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ interface{}
				_ = &mml.List{a[0:]}
				mml.Nop()
				return _removeToken
			},
			FixedArgs: 0,
		}
		_callTransform = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _transform = a[0]
				var _code = a[1]
				var _fields = a[2]
				var _listFields = a[3]
				var _ interface{}
				_ = &mml.List{a[4:]}
				mml.Nop(_transform, _code, _fields, _listFields)
				var _exists interface{}
				var _fieldResults interface{}
				var _fieldValues interface{}
				var _notToRemove interface{}
				var _listFieldValues interface{}
				var _results interface{}
				var _existingFields interface{}
				var _existingListFields interface{}
				mml.Nop(_exists, _fieldResults, _fieldValues, _notToRemove, _listFieldValues, _results, _existingFields, _existingListFields)
				_exists = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_f)
						return _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f, _code)}).Values)
					},
					FixedArgs: 1,
				}
				_fieldResults = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _keys = a[0]
						var _values = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_keys, _values)
						var _r interface{}
						mml.Nop(_r)
						_r = func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; ; return s }()
						for _i := 0; _i < _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _keys)}).Values).(int); _i++ {

							mml.Nop()
							c = mml.BinaryOp(11, mml.Ref(_values, _i), _removeToken)
							if c.(bool) {
								mml.Nop()
								continue
							}
							mml.SetRef(_r, mml.Ref(_keys, _i), mml.Ref(_values, _i))
						}
						return func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _r.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							return s
						}()
						return nil
					},
					FixedArgs: 2,
				}
				_existingFields = mml.Ref(_lists, "filter").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _exists, _fields)}).Values)
				_existingListFields = mml.Ref(_lists, "filter").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _exists, _listFields)}).Values)
				_fieldValues = mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_f)
						return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _transform, mml.Ref(_code, _f))}).Values)
					},
					FixedArgs: 1,
				}, _existingFields)}).Values)
				if v := _fieldValues; mml.IsError.F([]interface{}{v}).(bool) {
					return v
				}
				_notToRemove = mml.Ref(_functions, "not").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "bind").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_functions, "eq"), _removeToken)}).Values))}).Values)
				_listFieldValues = mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _f = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_f)
						return mml.Ref(_errors, "any").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "filter").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _notToRemove)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "map").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _transform)}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, _f))}).Values))}).Values))}).Values)
					},
					FixedArgs: 1,
				}, _existingListFields)}).Values)
				if v := mml.Ref(_errors, "any").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _listFieldValues)}).Values); mml.IsError.F([]interface{}{v}).(bool) {
					return v
				}
				_results = func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					func() {
						sp := _fieldResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _existingFields, _fieldValues)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					func() {
						sp := _fieldResults.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _existingListFields, _listFieldValues)}).Values).(*mml.Struct)
						for k, v := range sp.Values {
							s.Values[k] = v
						}
					}()
					return s
				}()
				return _transform.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code, _results)}).Values)
				return nil
			},
			FixedArgs: 4,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _transform = a[0]
				var _code = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_transform, _code)
				var _withFieldsAndLists interface{}
				var _withFields interface{}
				var _withListFields interface{}
				var _leaf interface{}
				mml.Nop(_withFieldsAndLists, _withFields, _withListFields, _leaf)
				_withFieldsAndLists = _callTransform.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _transform, _code)}).Values)
				_withFields = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _keys interface{}
						_keys = &mml.List{a[0:]}
						mml.Nop(_keys)
						return _withFieldsAndLists.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _keys, &mml.List{Values: []interface{}{}})}).Values)
					},
					FixedArgs: 0,
				}
				_withListFields = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _keys interface{}
						_keys = &mml.List{a[0:]}
						mml.Nop(_keys)
						return _withFieldsAndLists.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: []interface{}{}}, _keys)}).Values)
					},
					FixedArgs: 0,
				}
				_leaf = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ interface{}
						_ = &mml.List{a[0:]}
						mml.Nop()
						return _withFieldsAndLists.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: []interface{}{}}, &mml.List{Values: []interface{}{}})}).Values)
					},
					FixedArgs: 0,
				}
				switch mml.Ref(_code, "type") {
				case "spread":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value")}).Values)
				case "list":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "values")}).Values)
				case "mutable-list":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "values")}).Values)
				case "expression-key":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value")}).Values)
				case "entry":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "key", "value")}).Values)
				case "struct":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "entries")}).Values)
				case "mutable-struct":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "entries")}).Values)
				case "ret":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value")}).Values)
				case "check-ret":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value")}).Values)
				case "statement-list":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "statements")}).Values)
				case "function":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "body")}).Values)
				case "range":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "from", "to")}).Values)
				case "indexer":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", "index")}).Values)
				case "application":

					mml.Nop()
					return _withFieldsAndLists.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append([]interface{}{}, "function")}, &mml.List{Values: append([]interface{}{}, "args")})}).Values)
				case "unary":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "arg")}).Values)
				case "binary":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "left", "right")}).Values)
				case "cond":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "condition", "consequent", "alternative")}).Values)
				case "switch-case":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", "body")}).Values)
				case "switch-statement":

					mml.Nop()
					return _withFieldsAndLists.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append([]interface{}{}, "expression", "defaultStatements")}, &mml.List{Values: append([]interface{}{}, "cases")})}).Values)
				case "send-statement":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "channel", "value")}).Values)
				case "receive-expression":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "channel")}).Values)
				case "definition":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "docs", "expression")}).Values)
				case "select-case":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", "body")}).Values)
				case "select-statement":

					mml.Nop()
					return _withFieldsAndLists.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.List{Values: append([]interface{}{}, "defaultStatements")}, &mml.List{Values: append([]interface{}{}, "cases")})}).Values)
				case "go-statement":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "application")}).Values)
				case "defer-statement":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "application")}).Values)
				case "range-over":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression")}).Values)
				case "loop":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", "body")}).Values)
				case "assign":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "capture", "value")}).Values)
				case "definition-group":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definitions")}).Values)
				case "use":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "path")}).Values)
				case "use-list":

					mml.Nop()
					return _withListFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "uses")}).Values)
				case "module":

					mml.Nop()
					return _withFields.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "body")}).Values)
				default:

					mml.Nop()
					return _leaf.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				}
				return nil
			},
			FixedArgs: 2,
		}
		_edit = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _transform = a[0]
				var _code = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_transform, _code)
				return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _fieldResults = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_code, _fieldResults)
						return _transform.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							func() {
								sp := _code.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							func() {
								sp := _fieldResults.(*mml.Struct)
								for k, v := range sp.Values {
									s.Values[k] = v
								}
							}()
							return s
						}())}).Values)
					},
					FixedArgs: 2,
				}, _code)}).Values)
			},
			FixedArgs: 2,
		}
		exports["edit"] = _edit
		_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _predicate = a[0]
				var _code = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_predicate, _code)
				return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _fieldResults = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_code, _fieldResults)
						return func() interface{} {
							c = _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
							if c.(bool) {
								return &mml.List{Values: append(append([]interface{}{}, mml.Ref(_lists, "flatDepth").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.UnaryOp(2, 1), mml.Ref(_structs, "values").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fieldResults)}).Values))}).Values).(*mml.List).Values...), _code)}
							} else {
								return mml.Ref(_lists, "flatDepth").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.UnaryOp(2, 1), mml.Ref(_structs, "values").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _fieldResults)}).Values))}).Values)
							}
						}()
					},
					FixedArgs: 2,
				}, _code)}).Values)
			},
			FixedArgs: 2,
		}
		exports["filter"] = _filter
		_trim = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _predicate = a[0]
				var _code = a[1]
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_predicate, _code)
				var _result interface{}
				mml.Nop(_result)
				_result = _edit.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _code = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_code)
						return func() interface{} {
							c = _predicate.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
							if c.(bool) {
								return _removeToken
							} else {
								return _code
							}
						}()
					},
					FixedArgs: 1,
				}, _code)}).Values)
				return func() interface{} {
					c = mml.BinaryOp(11, _result, _removeToken)
					if c.(bool) {
						return _error.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "root node removed")}).Values)
					} else {
						return _result
					}
				}()
				return nil
			},
			FixedArgs: 2,
		}
		exports["trim"] = _trim

		return exports
	})

	modulePath = "io"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _readFile interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_readFile, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_readFile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_path)
				var _f interface{}
				mml.Nop(_f)
				_f = _open.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _path)}).Values)
				if v := _f; mml.IsError.F([]interface{}{v}).(bool) {
					return v
				}
				defer _close.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _f)}).Values)
				return _f.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.UnaryOp(2, 1))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		exports["readFile"] = _readFile

		return exports
	})

	modulePath = "paths"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _normalize interface{}
		var _trimExtension interface{}
		mml.Nop(_normalize, _trimExtension)
		_normalize = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_path)
				return _path
			},
			FixedArgs: 1,
		}
		exports["normalize"] = _normalize
		_trimExtension = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_path)
				return _path
			},
			FixedArgs: 1,
		}
		exports["trimExtension"] = _trimExtension

		return exports
	})

	modulePath = "compile"

	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)

		var _primitive interface{}
		var _stringLiteral interface{}
		var _symbol interface{}
		var _spread interface{}
		var _list interface{}
		var _expressionKey interface{}
		var _struct interface{}
		var _getDefinitions interface{}
		var _getScope interface{}
		var _paramList interface{}
		var _functionLiteral interface{}
		var _indexer interface{}
		var _application interface{}
		var _unary interface{}
		var _binary interface{}
		var _ternary interface{}
		var _ifStatement interface{}
		var _cond interface{}
		var _caseBlock interface{}
		var _switchStatement interface{}
		var _sendStatement interface{}
		var _receiveExpression interface{}
		var _goStatement interface{}
		var _deferStatement interface{}
		var _selectStatement interface{}
		var _rangeOver interface{}
		var _loop interface{}
		var _definition interface{}
		var _definitionGroup interface{}
		var _assign interface{}
		var _ret interface{}
		var _checkRet interface{}
		var _useStatement interface{}
		var _useList interface{}
		var _module interface{}
		var _statementList interface{}
		var _do interface{}
		var _allModules interface{}
		var _intLiteral interface{}
		var _floatLiteral interface{}
		var _boolLiteral interface{}
		var _breakStatement interface{}
		var _continueStatement interface{}
		var _toGo interface{}
		var _strings interface{}
		var _code interface{}
		var _lists interface{}
		var _structs interface{}
		var _snippets interface{}
		var _codetree interface{}
		var _fold interface{}
		var _foldr interface{}
		var _map interface{}
		var _filter interface{}
		var _contains interface{}
		var _sort interface{}
		var _flat interface{}
		var _flats interface{}
		var _concat interface{}
		var _concats interface{}
		var _uniq interface{}
		var _every interface{}
		var _some interface{}
		var _join interface{}
		var _joins interface{}
		var _formats interface{}
		var _enum interface{}
		var _log interface{}
		var _fatal interface{}
		var _bind interface{}
		var _identity interface{}
		var _eq interface{}
		var _any interface{}
		var _function interface{}
		var _channel interface{}
		var _natural interface{}
		var _type interface{}
		var _listOf interface{}
		var _structOf interface{}
		var _range interface{}
		var _rangeMin interface{}
		var _listLength interface{}
		var _or interface{}
		var _and interface{}
		var _not interface{}
		var _predicate interface{}
		var _predicates interface{}
		var _is interface{}
		mml.Nop(_primitive, _stringLiteral, _symbol, _spread, _list, _expressionKey, _struct, _getDefinitions, _getScope, _paramList, _functionLiteral, _indexer, _application, _unary, _binary, _ternary, _ifStatement, _cond, _caseBlock, _switchStatement, _sendStatement, _receiveExpression, _goStatement, _deferStatement, _selectStatement, _rangeOver, _loop, _definition, _definitionGroup, _assign, _ret, _checkRet, _useStatement, _useList, _module, _statementList, _do, _allModules, _intLiteral, _floatLiteral, _boolLiteral, _breakStatement, _continueStatement, _toGo, _strings, _code, _lists, _structs, _snippets, _codetree, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _flats, _concat, _concats, _uniq, _every, _some, _join, _joins, _formats, _enum, _log, _fatal, _bind, _identity, _eq, _any, _function, _channel, _natural, _type, _listOf, _structOf, _range, _rangeMin, _listLength, _or, _and, _not, _predicate, _predicates, _is)
		var __lang = mml.Modules.Use("lang")
		_fold = __lang.Values["fold"]
		_foldr = __lang.Values["foldr"]
		_map = __lang.Values["map"]
		_filter = __lang.Values["filter"]
		_contains = __lang.Values["contains"]
		_sort = __lang.Values["sort"]
		_flat = __lang.Values["flat"]
		_flats = __lang.Values["flats"]
		_concat = __lang.Values["concat"]
		_concats = __lang.Values["concats"]
		_uniq = __lang.Values["uniq"]
		_every = __lang.Values["every"]
		_some = __lang.Values["some"]
		_join = __lang.Values["join"]
		_joins = __lang.Values["joins"]
		_formats = __lang.Values["formats"]
		_enum = __lang.Values["enum"]
		_log = __lang.Values["log"]
		_fatal = __lang.Values["fatal"]
		_bind = __lang.Values["bind"]
		_identity = __lang.Values["identity"]
		_eq = __lang.Values["eq"]
		_any = __lang.Values["any"]
		_function = __lang.Values["function"]
		_channel = __lang.Values["channel"]
		_natural = __lang.Values["natural"]
		_type = __lang.Values["type"]
		_listOf = __lang.Values["listOf"]
		_structOf = __lang.Values["structOf"]
		_range = __lang.Values["range"]
		_rangeMin = __lang.Values["rangeMin"]
		_listLength = __lang.Values["listLength"]
		_or = __lang.Values["or"]
		_and = __lang.Values["and"]
		_not = __lang.Values["not"]
		_predicate = __lang.Values["predicate"]
		_predicates = __lang.Values["predicates"]
		_is = __lang.Values["is"]
		_strings = mml.Modules.Use("strings")
		_code = mml.Modules.Use("code")
		_lists = mml.Modules.Use("lists")
		_structs = mml.Modules.Use("structs")
		_snippets = mml.Modules.Use("snippets")
		_codetree = mml.Modules.Use("codetree")
		_primitive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_code)
				return _string.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "value"))}).Values)
			},
			FixedArgs: 1,
		}
		_intLiteral = _primitive
		_floatLiteral = _primitive
		_boolLiteral = _primitive
		_stringLiteral = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\"%s\"", mml.Ref(_strings, "escape").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_symbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s", mml.Ref(_s, "name"))}).Values)
			},
			FixedArgs: 1,
		}
		_spread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s.(*mml.List).Values...", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_list = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				var _groupSpread interface{}
				var _appendGroup interface{}
				var _isSpread interface{}
				var _selectSpread interface{}
				var _appendSimples interface{}
				var _appendSpread interface{}
				var _appendSpreads interface{}
				var _appendGroups interface{}
				mml.Nop(_groupSpread, _appendGroup, _isSpread, _selectSpread, _appendSimples, _appendSpread, _appendSpreads, _appendGroups)
				_isSpread = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
						var _ interface{}
						_ = &mml.List{a[1:]}
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
						var _ interface{}
						_ = &mml.List{a[2:]}
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
								var _ interface{}
								_ = &mml.List{a[0:]}
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
								var _ interface{}
								_ = &mml.List{a[0:]}
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
								var _ interface{}
								_ = &mml.List{a[0:]}
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
								var _ interface{}
								_ = &mml.List{a[0:]}
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
						var _ interface{}
						_ = &mml.List{a[2:]}
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
						var _ interface{}
						_ = &mml.List{a[2:]}
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
						var _ interface{}
						_ = &mml.List{a[2:]}
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
						var _ interface{}
						_ = &mml.List{a[1:]}
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
						var _ interface{}
						_ = &mml.List{a[2:]}
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
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_c)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "&mml.List{Values: %s}", _c)}).Values)
					},
					FixedArgs: 1,
				}).Call((&mml.List{Values: append([]interface{}{}, _appendGroups.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _groupSpread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _selectSpread)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_l, "values"))}).Values))}).Values))}).Values))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_expressionKey = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _k = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				var _entry interface{}
				mml.Nop(_entry)
				_entry = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _e = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_e)
						var _v interface{}
						mml.Nop(_v)
						_v = _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "value"))}).Values)
						switch mml.Ref(_e, "type") {
						case "spread":
							var _var interface{}
							var _assign interface{}
							mml.Nop(_var, _assign)
							_var = _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "sp := %s.(*mml.Struct);", _v)}).Values)
							_assign = "for k, v := range sp.Values { s.Values[k] = v };"
							return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func() { %s; %s }();\n", _var, _assign)}).Values)
						default:

							mml.Nop()
							switch mml.Ref(mml.Ref(_e, "key"), "type") {
							case "string":

								mml.Nop()
								return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "s.Values[%s] = %s;", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "key"))}).Values), _v)}).Values)
							case "symbol":

								mml.Nop()
								return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "s.Values[\"%s\"] = %s;", mml.Ref(mml.Ref(_e, "key"), "name"), _v)}).Values)
							default:

								mml.Nop()
								return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "s.Values[%s.(string)] = %s;", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_e, "key"))}).Values), _v)}).Values)
							}
						}
						return nil
					},
					FixedArgs: 1,
				}
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func() interface{} { s := &mml.Struct{Values: make(map[string]interface{})}; %s; return s }()", _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _entry)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "entries"))}).Values))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_getDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statementList = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_statementList)
				var _definitions interface{}
				var _definitionsFromGroups interface{}
				mml.Nop(_definitions, _definitionsFromGroups)
				_definitions = _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "definition"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_statementList, "statements"))}).Values)
				_definitionsFromGroups = _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "definitions")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "definition-group"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_statementList, "statements"))}).Values))}).Values))}).Values)
				return &mml.List{Values: append(append([]interface{}{}, _definitions.(*mml.List).Values...), _definitionsFromGroups.(*mml.List).Values...)}
				return nil
			},
			FixedArgs: 1,
		}
		_getScope = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statementList = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_statementList)
				var _definitions interface{}
				var _uses interface{}
				var _inlineUses interface{}
				var _unnamedUses interface{}
				var _namedUses interface{}
				mml.Nop(_definitions, _uses, _inlineUses, _unnamedUses, _namedUses)
				_definitions = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "symbol")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _getDefinitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _statementList)}).Values))}).Values)
				_uses = _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "uses")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "use-list"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_statementList, "statements"))}).Values))}).Values))}).Values)
				_unnamedUses = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "path")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _not.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["capture"] = _any
					return s
				}())}).Values))}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values))}).Values))}).Values)
				_namedUses = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "capture")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["capture"] = _not.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ".")}).Values)
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values))}).Values)
				_inlineUses = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "symbol")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["exported"] = true
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _getDefinitions)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "body")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "module")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["capture"] = "."
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _uses)}).Values))}).Values))}).Values))}).Values))}).Values))}).Values))}).Values)
				return _flats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _definitions, _unnamedUses, _namedUses, _inlineUses)}).Values)
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
				var _ interface{}
				_ = &mml.List{a[2:]}
				mml.Nop(_params, _collectParam)
				var _paramFormat interface{}
				var _collectParamFormat interface{}
				var _paramsString interface{}
				var _collectParamString interface{}
				mml.Nop(_paramFormat, _collectParamFormat, _paramsString, _collectParamString)
				_paramFormat = "var _%s = a[%d]"
				_collectParamFormat = "var _%s interface{}; _%s = &mml.List{a[%d:]}"
				_paramsString = _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _i = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_i)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _paramFormat, mml.Ref(_params, _i), _i)}).Values)
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_lists, "indexes").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values))}).Values)
				_collectParamString = _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _collectParamFormat, _collectParam, _collectParam, _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _params)}).Values))}).Values)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n", &mml.List{Values: append(append([]interface{}{}, _paramsString.(*mml.List).Values...), _collectParamString)})}).Values)
				return nil
			},
			FixedArgs: 2,
		}
		_functionLiteral = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_f)
				var _paramNames interface{}
				var _statementListFormat interface{}
				var _expressionFormat interface{}
				mml.Nop(_paramNames, _statementListFormat, _expressionFormat)
				_paramNames = func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_f, "collectParam"), "")
					if c.(bool) {
						return mml.Ref(_f, "params")
					} else {
						return &mml.List{Values: append(append([]interface{}{}, mml.Ref(_f, "params").(*mml.List).Values...), mml.Ref(_f, "collectParam"))}
					}
				}()
				_statementListFormat = "&mml.Function{\n\t\tF: func(a []interface{}) interface{} {\n\t\t\tvar c interface{}\n\t\t\tmml.Nop(c)\n\t\t\t%s;\n\t\t\tmml.Nop(%s);\n\t\t\t%s;\n\t\t\treturn nil\n\t\t},\n\t\tFixedArgs: %d,\n\t}"
				_expressionFormat = "&mml.Function{\n\t\tF: func(a []interface{}) interface{} {\n\t\t\tvar c interface{}\n\t\t\tmml.Nop(c)\n\t\t\t%s;\n\t\t\tmml.Nop(%s);\n\t\t\treturn %s\n\t\t},\n\t\tFixedArgs: %d,\n\t}"
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["type"] = "statement-list"
						return s
					}(), mml.Ref(_f, "body"))}).Values)
					if c.(bool) {
						return _statementListFormat
					} else {
						return _expressionFormat
					}
				}(), _paramList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "params"), mml.Ref(_f, "collectParam"))}).Values), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ", ", _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_strings, "formatOne").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s")}).Values), &mml.List{Values: append([]interface{}{}, _paramNames.(*mml.List).Values...)})}).Values))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "body"))}).Values), _len.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_f, "params"))}).Values))}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_indexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_i)

				mml.Nop()
				switch {
				case _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "range"
					return s
				}(), mml.Ref(_i, "index"))}).Values):

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.RefRange(%s, %s, %s)", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "expression"))}).Values), func() interface{} {
						c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["from"] = _any
							return s
						}(), mml.Ref(_i, "index"))}).Values)
						if c.(bool) {
							return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "from"))}).Values)
						} else {
							return "nil"
						}
					}(), func() interface{} {
						c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["to"] = _any
							return s
						}(), mml.Ref(_i, "index"))}).Values)
						if c.(bool) {
							return _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "to"))}).Values)
						} else {
							return "nil"
						}
					}())}).Values)
				case _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "symbol-index"
					return s
				}(), mml.Ref(_i, "index"))}).Values):

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.Ref(%s, \"%s\")", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "expression"))}).Values), mml.Ref(mml.Ref(mml.Ref(_i, "index"), "symbol"), "name"))}).Values)
				default:

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.Ref(%s, %s)", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "expression"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_i, "index"))}).Values))}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_application = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_a)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["function"] = func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "function"
							return s
						}()
						return s
					}(), _a)}).Values)
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_u)
				return func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_u, "op"), mml.Ref(_code, "logicalNot"))
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								s.Values["type"] = "bool"
								return s
							}(), mml.Ref(_u, "arg"))}).Values)
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_b)
				var _isBoolOp interface{}
				var _isBoolValue interface{}
				var _convertIfNotBool interface{}
				var _left interface{}
				var _right interface{}
				var _op interface{}
				mml.Nop(_isBoolOp, _isBoolValue, _convertIfNotBool, _left, _right, _op)
				c = !_is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "logicalAnd"), mml.Ref(_code, "logicalOr"))}).Values), mml.Ref(_b, "op"))}).Values).(bool)
				if c.(bool) {
					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "mml.BinaryOp(%d, %s, %s)", mml.Ref(_b, "op"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values))}).Values)
				}
				_isBoolOp = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_c)
						return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "unary", "binary")}).Values)
							s.Values["op"] = _or.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "logicalNot"), mml.Ref(_code, "logicalAnd"), mml.Ref(_code, "logicalOr"))}).Values)
							return s
						}(), _c)}).Values)
					},
					FixedArgs: 1,
				}
				_isBoolValue = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_c)
						return _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["type"] = "bool"
							return s
						}(), _c)}).Values)
					},
					FixedArgs: 1,
				}
				_convertIfNotBool = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _s = a[1]
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_c, _s)
						return func() interface{} {
							c = (_isBoolValue.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values).(bool) || _isBoolOp.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values).(bool))
							if c.(bool) {
								return _s
							} else {
								return mml.BinaryOp(9, _s, ".(bool)")
							}
						}()
					},
					FixedArgs: 2,
				}
				_left = _convertIfNotBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "left"))}).Values))}).Values)
				_right = _convertIfNotBool.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_b, "right"))}).Values))}).Values)
				_op = func() interface{} {
					c = mml.BinaryOp(11, mml.Ref(_b, "op"), mml.Ref(_code, "logicalAnd"))
					if c.(bool) {
						return "&&"
					} else {
						return "||"
					}
				}()
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
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_c)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "func () interface{} { c = %s; if c.(bool) { return %s } else { return %s } }()", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "condition"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "consequent"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "alternative"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_ifStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_c)
				return func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["alternative"] = _any
						return s
					}(), _c)}).Values)
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "c = %s; if c.(bool) { %s } else { %s }", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "condition"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "consequent"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "alternative"))}).Values))}).Values)
					} else {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "c = %s; if c.(bool) { %s }", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "condition"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "consequent"))}).Values))}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_cond = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_c)
				return func() interface{} {
					c = mml.Ref(_c, "ternary")
					if c.(bool) {
						return _ternary.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
					} else {
						return _ifStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _c)}).Values)
					}
				}()
			},
			FixedArgs: 1,
		}
		_caseBlock = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_c)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "case %s:\n%s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "expression"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_c, "body"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_switchStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["expression"] = _any
						return s
					}(), _s)}).Values)
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
		_sendStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s <- %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "channel"))}).Values), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_s, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_receiveExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_r)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "<-%s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "channel"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_goStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _g = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_g)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "go %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_g, "application"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_deferStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_d)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					c = _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["application"] = func() interface{} {
							s := &mml.Struct{Values: make(map[string]interface{})}
							s.Values["function"] = func() interface{} {
								s := &mml.Struct{Values: make(map[string]interface{})}
								s.Values["type"] = "function"
								return s
							}()
							return s
						}()
						return s
					}(), _d)}).Values)
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
			},
			FixedArgs: 1,
		}
		_selectStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_s)
				return (&mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _c = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
		_rangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_r)
				var _infiniteCounter interface{}
				var _withRangeExpression interface{}
				var _listStyleRange interface{}
				mml.Nop(_infiniteCounter, _withRangeExpression, _listStyleRange)
				_infiniteCounter = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ interface{}
						_ = &mml.List{a[0:]}
						mml.Nop()
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s := 0; true; _%s++", mml.Ref(_r, "symbol"), mml.Ref(_r, "symbol"))}).Values)
					},
					FixedArgs: 0,
				}
				_withRangeExpression = &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _ interface{}
						_ = &mml.List{a[0:]}
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
						var _ interface{}
						_ = &mml.List{a[0:]}
						mml.Nop()
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_, _%s := range %s.(*mml.List).Values", mml.Ref(_r, "symbol"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "expression"))}).Values))}).Values)
					},
					FixedArgs: 0,
				}
				switch {
				case !_has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "expression", _r)}).Values).(bool):

					mml.Nop()
					return _infiniteCounter.(*mml.Function).Call((&mml.List{Values: []interface{}{}}).Values)
				case _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "range"
					return s
				}(), mml.Ref(_r, "expression"))}).Values):

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
		_breakStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var __ = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(__)
				return "break"
			},
			FixedArgs: 1,
		}
		_continueStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var __ = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(__)
				return "continue"
			},
			FixedArgs: 1,
		}
		_loop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_definitionGroup = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _g = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_g)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_g, "definitions"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_assign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
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
		_ret = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_r)
				return func() interface{} {
					c = _has.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "value", _r)}).Values)
					if c.(bool) {
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "return %s", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "value"))}).Values))}).Values)
					} else {
						return "return"
					}
				}()
			},
			FixedArgs: 1,
		}
		_checkRet = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_r)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "if v := %s; mml.IsError.F([]interface{}{v}).(bool) { return v }", _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_r, "value"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_useStatement = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_u)

				mml.Nop()
				switch {
				case _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["capture"] = "."
					return s
				}(), _u)}).Values):
					var _statement interface{}
					var _assigns interface{}
					mml.Nop(_statement, _assigns)
					_statement = _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var __%s = mml.Modules.Use(\"%s\");", mml.Ref(_code, "getModuleName").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_u, "path"), "value"))}).Values), mml.Ref(mml.Ref(_u, "path"), "value"))}).Values)
					_assigns = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
						F: func(a []interface{}) interface{} {
							var c interface{}
							mml.Nop(c)
							var _name = a[0]
							var _ interface{}
							_ = &mml.List{a[1:]}
							mml.Nop(_name)
							return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = __%s.Values[\"%s\"]", _name, mml.Ref(_code, "getModuleName").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_u, "path"), "value"))}).Values), _name)}).Values)
						},
						FixedArgs: 1,
					}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "symbol")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _filter.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
						s := &mml.Struct{Values: make(map[string]interface{})}
						s.Values["exported"] = true
						return s
					}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _getDefinitions.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_u, "module"), "body"))}).Values))}).Values))}).Values))}).Values))}).Values)
					return _joins.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";", _statement, _assigns)}).Values)
				case _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["capture"] = _any
					return s
				}(), _u)}).Values):

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = mml.Modules.Use(\"%s\")", mml.Ref(_u, "capture"), mml.Ref(mml.Ref(_u, "path"), "value"))}).Values)
				default:

					mml.Nop()
					return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "_%s = mml.Modules.Use(\"%s\")", mml.Ref(_code, "getModuleName").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(mml.Ref(_u, "path"), "value"))}).Values), mml.Ref(mml.Ref(_u, "path"), "value"))}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_useList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_u)
				return _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_u, "uses"))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_module = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_m)
				return _joins.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\n", _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "modulePath = \"%s\"", mml.Ref(_m, "path"))}).Values), mml.Ref(_snippets, "moduleHead"), _do.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_m, "body"))}).Values), mml.Ref(_snippets, "moduleFooter"))}).Values)
			},
			FixedArgs: 1,
		}
		_statementList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_l)
				var _scopeDefs interface{}
				var _scope interface{}
				var _scopeNames interface{}
				var _statements interface{}
				mml.Nop(_scopeDefs, _scope, _scopeNames, _statements)
				_scope = _getScope.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _l)}).Values)
				_scopeNames = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ", ", _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _bind.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_strings, "formats"), "_%s")}).Values), _scope)}).Values))}).Values)
				_statements = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do, mml.Ref(_l, "statements"))}).Values))}).Values)
				_scopeDefs = _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _s = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
						mml.Nop(_s)
						return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "var _%s interface{}", _s)}).Values)
					},
					FixedArgs: 1,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _scope)}).Values))}).Values)
				return _formats.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "%s;\nmml.Nop(%s);\n%s", _scopeDefs, _scopeNames, _statements)}).Values)
				return nil
			},
			FixedArgs: 1,
		}
		_do = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_code)

				mml.Nop()
				switch {
				case mml.BinaryOp(11, mml.Ref(_code, "type"), "int"):

					mml.Nop()
					return _intLiteral.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case mml.BinaryOp(11, mml.Ref(_code, "type"), "float"):

					mml.Nop()
					return _floatLiteral.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case mml.BinaryOp(11, mml.Ref(_code, "type"), "string"):

					mml.Nop()
					return _stringLiteral.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case mml.BinaryOp(11, mml.Ref(_code, "type"), "bool"):

					mml.Nop()
					return _boolLiteral.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				}
				switch mml.Ref(_code, "type") {
				case "comment":

					mml.Nop()
					return ""
				case "symbol":

					mml.Nop()
					return _symbol.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "list":

					mml.Nop()
					return _list.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "expression-key":

					mml.Nop()
					return _expressionKey.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "struct":

					mml.Nop()
					return _struct.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "function":

					mml.Nop()
					return _functionLiteral.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "indexer":

					mml.Nop()
					return _indexer.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "spread":

					mml.Nop()
					return _spread.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "application":

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
					return _caseBlock.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "switch-statement":

					mml.Nop()
					return _switchStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "send-statement":

					mml.Nop()
					return _sendStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "receive-expression":

					mml.Nop()
					return _receiveExpression.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "go-statement":

					mml.Nop()
					return _goStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "defer-statement":

					mml.Nop()
					return _deferStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "select-case":

					mml.Nop()
					return _caseBlock.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "select-statement":

					mml.Nop()
					return _selectStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "range-over":

					mml.Nop()
					return _rangeOver.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "break":

					mml.Nop()
					return _breakStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "continue":

					mml.Nop()
					return _continueStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "loop":

					mml.Nop()
					return _loop.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "definition":

					mml.Nop()
					return _definition.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "definition-group":

					mml.Nop()
					return _definitionGroup.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "assign":

					mml.Nop()
					return _assign.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "ret":

					mml.Nop()
					return _ret.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "check-ret":

					mml.Nop()
					return _checkRet.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "use":

					mml.Nop()
					return _useStatement.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "use-list":

					mml.Nop()
					return _useList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				case "module":

					mml.Nop()
					return _module.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				default:

					mml.Nop()
					return _statementList.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _code)}).Values)
				}
				return nil
			},
			FixedArgs: 1,
		}
		_allModules = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _module = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_module)
				return _uniq.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _eq)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _bind.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _concats, &mml.List{Values: append([]interface{}{}, _module)})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _flat.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _allModules)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_structs, "get").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "module")}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_codetree, "filter").(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _is.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, func() interface{} {
					s := &mml.Struct{Values: make(map[string]interface{})}
					s.Values["type"] = "use"
					return s
				}())}).Values))}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _module)}).Values))}).Values))}).Values))}).Values))}).Values))}).Values)
			},
			FixedArgs: 1,
		}
		_toGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _module = a[0]
				var _ interface{}
				_ = &mml.List{a[1:]}
				mml.Nop(_module)
				return _joins.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "", mml.Ref(_snippets, "head"), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, ";\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, &mml.Function{
					F: func(a []interface{}) interface{} {
						var c interface{}
						mml.Nop(c)
						var _k = a[0]
						var _ interface{}
						_ = &mml.List{a[1:]}
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
						var _ interface{}
						_ = &mml.List{a[2:]}
						mml.Nop(_left, _right)
						return mml.BinaryOp(13, _left, _right)
					},
					FixedArgs: 2,
				})}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _keys.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, mml.Ref(_code, "builtin"))}).Values))}).Values))}).Values))}).Values), mml.Ref(_snippets, "initHead"), _join.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, "\n")}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _map.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _do)}).Values).(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _allModules.(*mml.Function).Call((&mml.List{Values: append([]interface{}{}, _module)}).Values))}).Values))}).Values), mml.Ref(_snippets, "initFooter"), mml.Ref(_snippets, "mainHead"), mml.Ref(_module, "path"), mml.Ref(_snippets, "mainFooter"))}).Values)
			},
			FixedArgs: 1,
		}
		exports["toGo"] = _toGo

		return exports
	})

	modulePath = "snippets"

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

}

func main() {
	mml.Modules.Use("main")
}
