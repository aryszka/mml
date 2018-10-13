// Generated code
package main

import "github.com/aryszka/mml"
var _args interface{} = mml.Args;
var _close interface{} = mml.Close;
var _error interface{} = mml.Error;
var _format interface{} = mml.Format;
var _has interface{} = mml.Has;
var _isBool interface{} = mml.IsBool;
var _isError interface{} = mml.IsError;
var _isFloat interface{} = mml.IsFloat;
var _isInt interface{} = mml.IsInt;
var _isString interface{} = mml.IsString;
var _keys interface{} = mml.Keys;
var _len interface{} = mml.Len;
var _open interface{} = mml.Open;
var _panic interface{} = mml.Panic;
var _parseAST interface{} = mml.ParseAST;
var _parseFloat interface{} = mml.ParseFloat;
var _parseInt interface{} = mml.ParseInt;
var _stderr interface{} = mml.Stderr;
var _stdin interface{} = mml.Stdin;
var _stdout interface{} = mml.Stdout;
var _string interface{} = mml.String
func init() {
	var modulePath string
modulePath = "main.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _builtin interface{};
var _builtins interface{};
var _head interface{};
var _initHead interface{};
var _initFooter interface{};
var _setModuleHead interface{};
var _setModuleFooter interface{};
var _mainHead interface{};
var _mainFooter interface{};
var _compileModuleCode interface{};
var _compileModules interface{};
var _modules interface{};
var _compile interface{};
var _parse interface{};
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
mml.Nop(_builtin, _builtins, _head, _initHead, _initFooter, _setModuleHead, _setModuleFooter, _mainHead, _mainFooter, _compileModuleCode, _compileModules, _modules, _compile, _parse, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr);
var __lang = mml.Modules.Use("lang.mml");;_fold = __lang["fold"];
_foldr = __lang["foldr"];
_map = __lang["map"];
_filter = __lang["filter"];
_contains = __lang["contains"];
_sort = __lang["sort"];
_flat = __lang["flat"];
_join = __lang["join"];
_joins = __lang["joins"];
_formats = __lang["formats"];
_enum = __lang["enum"];
_log = __lang["log"];
_onlyErr = __lang["onlyErr"];
_passErr = __lang["passErr"];
_compile = mml.Modules.Use("compile.mml");
_parse = mml.Modules.Use("parse.mml");
_builtin = func() interface{} { s := make(map[string]interface{}); s["len"] = "Len";s["isError"] = "IsError";s["keys"] = "Keys";s["format"] = "Format";s["stdin"] = "Stdin";s["stdout"] = "Stdout";s["stderr"] = "Stderr";s["string"] = "String";s["has"] = "Has";s["isBool"] = "IsBool";s["isInt"] = "IsInt";s["isFloat"] = "IsFloat";s["isString"] = "IsString";s["error"] = "Error";s["panic"] = "Panic";s["open"] = "Open";s["close"] = "Close";s["args"] = "Args";s["parseAST"] = "ParseAST";s["parseInt"] = "ParseInt";s["parseFloat"] = "ParseFloat";; return s }();
_builtins = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _k = a[0];
				;
				mml.Nop(_k);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s interface{} = mml.%s", _k, mml.Ref(_builtin, _k)))
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _left = a[0];
var _right = a[1];
				;
				mml.Nop(_left, _right);
				return mml.BinaryOp(13, _left, _right)
			},
			FixedArgs: 2,
		})).(*mml.Function).Call(append([]interface{}{}, _keys.(*mml.Function).Call(append([]interface{}{}, _builtin))))))));
_head = "// Generated code\npackage main\n\nimport \"github.com/aryszka/mml\"\n";
_initHead = "\nfunc init() {\n\tvar modulePath string\n";
_initFooter = "\n}\n";
_setModuleHead = "\n\tmml.Modules.Set(modulePath, func() map[string]interface{} {\n\t\texports := make(map[string]interface{})\n\n\t\tvar c interface{}\n\t\tmml.Nop(c)\n";
_setModuleFooter = "\n\t\treturn exports\n\t})\n";
_mainHead = "\nfunc main() {\n\tmml.Modules.Use(\"";
_mainFooter = "\")\n}\n";
_compileModuleCode = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _moduleCode = a[0];
				;
				mml.Nop(_moduleCode);
				;
mml.Nop();
_stdout.(*mml.Function).Call(append([]interface{}{}, _formats.(*mml.Function).Call(append([]interface{}{}, "modulePath = \"%s\"", mml.Ref(_moduleCode, "path")))));
_stdout.(*mml.Function).Call(append([]interface{}{}, _setModuleHead));
_onlyErr.(*mml.Function).Call(append([]interface{}{}, _log)).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _stdout)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_compile, "compile").(*mml.Function).Call(append([]interface{}{}, _moduleCode))))));
_stdout.(*mml.Function).Call(append([]interface{}{}, _setModuleFooter));
				return nil
			},
			FixedArgs: 1,
		};
_compileModules = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
				;
				mml.Nop(_m);
				;
mml.Nop();
for _, _mi := range _m.([]interface{}) {
;
mml.Nop();
_compileModuleCode.(*mml.Function).Call(append([]interface{}{}, _mi))
};
				return nil
			},
			FixedArgs: 1,
		};
_modules = mml.Ref(_parse, "parseModules").(*mml.Function).Call(append([]interface{}{}, mml.Ref(_args, 1)));
c = _isError.(*mml.Function).Call(append([]interface{}{}, _modules)); if c.(bool) { ;
mml.Nop();
_panic.(*mml.Function).Call(append([]interface{}{}, _modules)) };
_stdout.(*mml.Function).Call(append([]interface{}{}, _head));
_stdout.(*mml.Function).Call(append([]interface{}{}, _builtins));
_stdout.(*mml.Function).Call(append([]interface{}{}, _initHead));
_compileModules.(*mml.Function).Call(append([]interface{}{}, _modules));
_stdout.(*mml.Function).Call(append([]interface{}{}, _initFooter));
_stdout.(*mml.Function).Call(append([]interface{}{}, _mainHead));
_stdout.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_args, 1)));
_stdout.(*mml.Function).Call(append([]interface{}{}, _mainFooter))
		return exports
	})
modulePath = "lang.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
var _logger interface{};
var _list interface{};
var _strings interface{};
var _ints interface{};
var _errors interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _logger, _list, _strings, _ints, _errors);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_ints = mml.Modules.Use("ints.mml");
_logger = mml.Modules.Use("log.mml");
_errors = mml.Modules.Use("errors.mml");
_fold = mml.Ref(_list, "fold"); exports["fold"] = _fold;
_foldr = mml.Ref(_list, "foldr"); exports["foldr"] = _foldr;
_map = mml.Ref(_list, "map"); exports["map"] = _map;
_filter = mml.Ref(_list, "filter"); exports["filter"] = _filter;
_contains = mml.Ref(_list, "contains"); exports["contains"] = _contains;
_sort = mml.Ref(_list, "sort"); exports["sort"] = _sort;
_flat = mml.Ref(_list, "flat"); exports["flat"] = _flat;
_join = mml.Ref(_strings, "join"); exports["join"] = _join;
_joins = mml.Ref(_strings, "joins"); exports["joins"] = _joins;
_formats = mml.Ref(_strings, "formats"); exports["formats"] = _formats;
_enum = mml.Ref(_ints, "enum"); exports["enum"] = _enum;
_log = mml.Ref(_logger, "log"); exports["log"] = _log;
_onlyErr = mml.Ref(_errors, "only"); exports["onlyErr"] = _onlyErr;
_passErr = mml.Ref(_errors, "pass"); exports["passErr"] = _passErr
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "ints.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _counter interface{};
var _enum interface{};
mml.Nop(_counter, _enum);
_counter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _c interface{};
mml.Nop(_c);
_c = mml.UnaryOp(2, 1);
return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				;
mml.Nop();
_c = mml.BinaryOp(9, _c, 1);
return _c;
				return nil
			},
			FixedArgs: 0,
		};
				return nil
			},
			FixedArgs: 0,
		}; exports["counter"] = _counter;
_enum = _counter; exports["enum"] = _enum
		return exports
	})
modulePath = "log.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _log interface{};
var _list interface{};
var _strings interface{};
mml.Nop(_log, _list, _strings);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_log = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0:];
				;
				mml.Nop(_a);
				;
mml.Nop();
_stderr.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call(append([]interface{}{}, " ")).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_list, "map").(*mml.Function).Call(append([]interface{}{}, _string)).(*mml.Function).Call(append([]interface{}{}, _a))))));
_stderr.(*mml.Function).Call(append([]interface{}{}, "\n"));
return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 0); if c.(bool) { return "" } else { return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 1)) } }();
				return nil
			},
			FixedArgs: 0,
		}; exports["log"] = _log
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "compile.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _notEmpty interface{};
var _compileInt interface{};
var _compileFloat interface{};
var _compileBool interface{};
var _getScope interface{};
var _compileComment interface{};
var _compileString interface{};
var _compileSymbol interface{};
var _compileCond interface{};
var _compileSpread interface{};
var _compileCase interface{};
var _compileSend interface{};
var _compileReceive interface{};
var _compileGo interface{};
var _compileDefer interface{};
var _compileDefinitions interface{};
var _compileAssigns interface{};
var _compileRet interface{};
var _compileControl interface{};
var _compileUseList interface{};
var _compileList interface{};
var _compileEntry interface{};
var _compileStruct interface{};
var _compileParamList interface{};
var _compileFunction interface{};
var _compileRangeExpression interface{};
var _compileIndexer interface{};
var _compileApplication interface{};
var _compileUnary interface{};
var _compileBinary interface{};
var _compileTernary interface{};
var _compileIf interface{};
var _compileSwitch interface{};
var _compileSelect interface{};
var _compileRangeOver interface{};
var _compileLoop interface{};
var _compileDefinition interface{};
var _compileAssign interface{};
var _compileStatements interface{};
var _compileModule interface{};
var _getModuleName interface{};
var _compileUse interface{};
var _compile interface{};
var _errors interface{};
var _code interface{};
var _strings interface{};
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
mml.Nop(_notEmpty, _compileInt, _compileFloat, _compileBool, _getScope, _compileComment, _compileString, _compileSymbol, _compileCond, _compileSpread, _compileCase, _compileSend, _compileReceive, _compileGo, _compileDefer, _compileDefinitions, _compileAssigns, _compileRet, _compileControl, _compileUseList, _compileList, _compileEntry, _compileStruct, _compileParamList, _compileFunction, _compileRangeExpression, _compileIndexer, _compileApplication, _compileUnary, _compileBinary, _compileTernary, _compileIf, _compileSwitch, _compileSelect, _compileRangeOver, _compileLoop, _compileDefinition, _compileAssign, _compileStatements, _compileModule, _getModuleName, _compileUse, _compile, _errors, _code, _strings, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr);
var __lang = mml.Modules.Use("lang.mml");;_fold = __lang["fold"];
_foldr = __lang["foldr"];
_map = __lang["map"];
_filter = __lang["filter"];
_contains = __lang["contains"];
_sort = __lang["sort"];
_flat = __lang["flat"];
_join = __lang["join"];
_joins = __lang["joins"];
_formats = __lang["formats"];
_enum = __lang["enum"];
_log = __lang["log"];
_onlyErr = __lang["onlyErr"];
_passErr = __lang["passErr"];
_errors = mml.Modules.Use("errors.mml");
_code = mml.Modules.Use("code.mml");
_strings = mml.Modules.Use("strings.mml");
_notEmpty = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return mml.BinaryOp(12, _s, "")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _l))
			},
			FixedArgs: 1,
		};
_compileInt = _string;
_compileFloat = _string;
_compileBool = _string;
_getScope = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statements = a[0:];
				;
				mml.Nop(_statements);
				var _defs interface{};
var _uses interface{};
var _inlineUses interface{};
var _namedUses interface{};
var _unnamedUses interface{};
mml.Nop(_defs, _uses, _inlineUses, _namedUses, _unnamedUses);
_defs = mml.Ref(_code, "flattenedStatements").(*mml.Function).Call(append([]interface{}{}, "definition", "definition-list", "definitions", _statements));
_uses = mml.Ref(_code, "flattenedStatements").(*mml.Function).Call(append([]interface{}{}, "use", "use-list", "uses", _statements));
_inlineUses = _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _current = a[0];
var _result = a[1];
				;
				mml.Nop(_current, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _current.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{})).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return mml.Ref(_u, "exportNames")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _has.(*mml.Function).Call(append([]interface{}{}, "exportNames")))).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return mml.BinaryOp(11, mml.Ref(_u, "capture"), ".")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _uses))))))));
_namedUses = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return (mml.BinaryOp(12, mml.Ref(_u, "capture"), ".").(bool) && mml.BinaryOp(12, mml.Ref(_u, "capture"), "").(bool))
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _uses));
_unnamedUses = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return mml.BinaryOp(11, mml.Ref(_u, "capture"), "")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _uses));
return _flat.(*mml.Function).Call(append([]interface{}{}, append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return mml.Ref(_d, "symbol")
			},
			FixedArgs: 1,
		}, _defs)), _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return mml.Ref(_u, "capture")
			},
			FixedArgs: 1,
		}, _namedUses)), _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return mml.Ref(_u, "path")
			},
			FixedArgs: 1,
		}, _unnamedUses)), _inlineUses)));
				return nil
			},
			FixedArgs: 0,
		};
_compileComment = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var __ = a[0];
				;
				mml.Nop(__);
				return ""
			},
			FixedArgs: 1,
		};
_compileString = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "\"%s\"", mml.Ref(_strings, "escape").(*mml.Function).Call(append([]interface{}{}, _s))))
			},
			FixedArgs: 1,
		};
_compileSymbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s", mml.Ref(_s, "name")))
			},
			FixedArgs: 1,
		};
_compileCond = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return func () interface{} { c = mml.Ref(_c, "ternary"); if c.(bool) { return _compileTernary.(*mml.Function).Call(append([]interface{}{}, _c)) } else { return _compileIf.(*mml.Function).Call(append([]interface{}{}, _c)) } }()
			},
			FixedArgs: 1,
		};
_compileSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "%s.([]interface{})...", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "value")))))
			},
			FixedArgs: 1,
		};
_compileCase = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "case %s:\n%s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "expression"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "body")))))
			},
			FixedArgs: 1,
		};
_compileSend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "%s <- %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "channel"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "value")))))
			},
			FixedArgs: 1,
		};
_compileReceive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0];
				;
				mml.Nop(_r);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "<- %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_r, "channel")))))
			},
			FixedArgs: 1,
		};
_compileGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _g = a[0];
				;
				mml.Nop(_g);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "go %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_g, "application")))))
			},
			FixedArgs: 1,
		};
_compileDefer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "defer %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_d, "application")))))
			},
			FixedArgs: 1,
		};
_compileDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "definitions")))))
			},
			FixedArgs: 1,
		};
_compileAssigns = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "assignments")))))
			},
			FixedArgs: 1,
		};
_compileRet = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0];
				;
				mml.Nop(_r);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "return %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_r, "value")))))
			},
			FixedArgs: 1,
		};
_compileControl = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return func () interface{} { c = mml.BinaryOp(11, mml.Ref(_c, "control"), mml.Ref(_code, "breakControl")); if c.(bool) { return "break" } else { return "continue" } }()
			},
			FixedArgs: 1,
		};
_compileUseList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "uses")))))
			},
			FixedArgs: 1,
		};
_compileList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				var _isSpread interface{};
var _selectSpread interface{};
var _groupSpread interface{};
var _appendSimples interface{};
var _appendSpread interface{};
var _appendSpreads interface{};
var _appendGroups interface{};
var _appendGroup interface{};
mml.Nop(_isSpread, _selectSpread, _groupSpread, _appendSimples, _appendSpread, _appendSpreads, _appendGroups, _appendGroup);
_isSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return (mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _c)), 3).(bool) && mml.BinaryOp(11, mml.RefRange(_c, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _c)), 3), nil), "...").(bool))
			},
			FixedArgs: 1,
		};
_selectSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return func () interface{} { c = _isSpread.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return func() interface{} { s := make(map[string]interface{}); s["spread"] = _c;; return s }() } else { return _c } }()
			},
			FixedArgs: 1,
		};
_groupSpread = _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _item = a[0];
var _groups = a[1];
				;
				mml.Nop(_item, _groups);
				var _i interface{};
var _isSpread interface{};
var _groupIsSpread interface{};
var _appendNewSimple interface{};
var _appendNewSpread interface{};
var _appendSimple interface{};
var _appendSpread interface{};
mml.Nop(_i, _isSpread, _groupIsSpread, _appendNewSimple, _appendNewSpread, _appendSimple, _appendSpread);
_i = mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _groups)), 1);
_isSpread = _has.(*mml.Function).Call(append([]interface{}{}, "spread", _item));
_groupIsSpread = (mml.BinaryOp(16, _i, 0).(bool) && _has.(*mml.Function).Call(append([]interface{}{}, "spread", mml.Ref(_groups, _i))).(bool));
_appendNewSimple = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return append(append([]interface{}{}, _groups.([]interface{})...), func() interface{} { s := make(map[string]interface{}); s["simple"] = append([]interface{}{}, _item);; return s }())
			},
			FixedArgs: 0,
		};
_appendNewSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return append(append([]interface{}{}, _groups.([]interface{})...), func() interface{} { s := make(map[string]interface{}); s["spread"] = append([]interface{}{}, mml.Ref(_item, "spread"));; return s }())
			},
			FixedArgs: 0,
		};
_appendSimple = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return append(append([]interface{}{}, mml.RefRange(_groups, nil, _i).([]interface{})...), func() interface{} { s := make(map[string]interface{}); s["simple"] = append(append([]interface{}{}, mml.Ref(mml.Ref(_groups, _i), "simple").([]interface{})...), _item);; return s }())
			},
			FixedArgs: 0,
		};
_appendSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return append(append([]interface{}{}, mml.RefRange(_groups, nil, _i).([]interface{})...), func() interface{} { s := make(map[string]interface{}); s["spread"] = append(append([]interface{}{}, mml.Ref(mml.Ref(_groups, _i), "spread").([]interface{})...), mml.Ref(_item, "spread"));; return s }())
			},
			FixedArgs: 0,
		};
switch  {
case ((mml.BinaryOp(13, _i, 0).(bool) || _groupIsSpread.(bool)) && !_isSpread.(bool)):
;
mml.Nop();
return _appendNewSimple.(*mml.Function).Call([]interface{}{})
case ((mml.BinaryOp(13, _i, 0).(bool) || !_groupIsSpread.(bool)) && _isSpread.(bool)):
;
mml.Nop();
return _appendNewSpread.(*mml.Function).Call([]interface{}{})
case (!_groupIsSpread.(bool) && !_isSpread.(bool)):
;
mml.Nop();
return _appendSimple.(*mml.Function).Call([]interface{}{})
case (_groupIsSpread.(bool) && _isSpread.(bool)):
;
mml.Nop();
return _appendSpread.(*mml.Function).Call([]interface{}{})
};
				return nil
			},
			FixedArgs: 2,
		}, []interface{}{}));
_appendSimples = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0];
var _group = a[1];
				;
				mml.Nop(_code, _group);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "append(%s, %s)", _code, _join.(*mml.Function).Call(append([]interface{}{}, ", ", _group))))
			},
			FixedArgs: 2,
		};
_appendSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _item = a[0];
var _code = a[1];
				;
				mml.Nop(_item, _code);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "append(%s, %s)", _code, _item))
			},
			FixedArgs: 2,
		};
_appendSpreads = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0];
var _group = a[1];
				;
				mml.Nop(_code, _group);
				return _fold.(*mml.Function).Call(append([]interface{}{}, _appendSpread, _code, _group))
			},
			FixedArgs: 2,
		};
_appendGroups = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _groups = a[0];
				;
				mml.Nop(_groups);
				return _fold.(*mml.Function).Call(append([]interface{}{}, _appendGroup, "[]interface{}{}", _groups))
			},
			FixedArgs: 1,
		};
_appendGroup = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _group = a[0];
var _code = a[1];
				;
				mml.Nop(_group, _code);
				return func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "spread", _group)); if c.(bool) { return _appendSpreads.(*mml.Function).Call(append([]interface{}{}, _code, mml.Ref(_group, "spread"))) } else { return _appendSimples.(*mml.Function).Call(append([]interface{}{}, _code, mml.Ref(_group, "simple"))) } }()
			},
			FixedArgs: 2,
		};
return _appendGroups.(*mml.Function).Call(append([]interface{}{}, _groupSpread.(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _selectSpread)).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "values")))))))));
				return nil
			},
			FixedArgs: 1,
		};
_compileEntry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0];
				;
				mml.Nop(_e);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "\"%s\":%s", func () interface{} { c = (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_e, "key"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_e, "key"), "type"), "symbol").(bool)); if c.(bool) { return mml.Ref(mml.Ref(_e, "key"), "name") } else { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "key"))) } }(), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "value")))))
			},
			FixedArgs: 1,
		};
_compileStruct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _compileEntry interface{};
var _entries interface{};
mml.Nop(_compileEntry, _entries);
_compileEntry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0];
				;
				mml.Nop(_e);
				var _v interface{};
mml.Nop(_v);
_v = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "value")));
c = mml.BinaryOp(11, mml.Ref(_e, "type"), "spread"); if c.(bool) { var _var interface{};
var _assign interface{};
mml.Nop(_var, _assign);
_var = _formats.(*mml.Function).Call(append([]interface{}{}, "sp := %s.(map[string]interface{});", _v));
_assign = "for k, v := range sp { s[k] = v };";
return _joins.(*mml.Function).Call(append([]interface{}{}, "\n", _var, _assign)) };
c = _isString.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "key"))); if c.(bool) { ;
mml.Nop();
return _formats.(*mml.Function).Call(append([]interface{}{}, "s[\"%s\"] = %s;", mml.Ref(_e, "key"), _v)) };
c = mml.BinaryOp(11, mml.Ref(mml.Ref(_e, "key"), "type"), "symbol"); if c.(bool) { ;
mml.Nop();
return _formats.(*mml.Function).Call(append([]interface{}{}, "s[\"%s\"] = %s;", mml.Ref(mml.Ref(_e, "key"), "name"), _v)) };
return _formats.(*mml.Function).Call(append([]interface{}{}, "s[%s] = %s;", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "key"))), _v));
				return nil
			},
			FixedArgs: 1,
		};
_entries = _map.(*mml.Function).Call(append([]interface{}{}, _compileEntry)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "entries")));
return _formats.(*mml.Function).Call(append([]interface{}{}, "func() interface{} { s := make(map[string]interface{}); %s; return s }()", _join.(*mml.Function).Call(append([]interface{}{}, "")).(*mml.Function).Call(append([]interface{}{}, _entries))));
				return nil
			},
			FixedArgs: 1,
		};
_compileParamList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _params = a[0];
var _collectParam = a[1];
				;
				mml.Nop(_params, _collectParam);
				var _p interface{};
mml.Nop(_p);
_p = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _params)).(int); _i++ {
;
mml.Nop();
_p = append(append([]interface{}{}, _p.([]interface{})...), _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s = a[%d]", mml.Ref(_params, _i), _i)))
};
c = mml.BinaryOp(12, _collectParam, ""); if c.(bool) { ;
mml.Nop();
_p = append(append([]interface{}{}, _p.([]interface{})...), _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s = a[%d:]", _collectParam, _len.(*mml.Function).Call(append([]interface{}{}, _params))))) };
return _join.(*mml.Function).Call(append([]interface{}{}, ";\n", _p));
				return nil
			},
			FixedArgs: 2,
		};
_compileFunction = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				var _scope interface{};
var _paramNames interface{};
mml.Nop(_scope, _paramNames);
_scope = _getScope.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "statement")));
_paramNames = func () interface{} { c = mml.BinaryOp(11, mml.Ref(_f, "collectParam"), ""); if c.(bool) { return mml.Ref(_f, "params") } else { return append(append([]interface{}{}, mml.Ref(_f, "params").([]interface{})...), mml.Ref(_f, "collectParam")) } }();
return _formats.(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_f, "statement"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_f, "statement"), "type"), "statement-list").(bool)); if c.(bool) { return "&mml.Function{\n\t\t\tF: func(a []interface{}) interface{} {\n\t\t\t\tvar c interface{}\n\t\t\t\tmml.Nop(c)\n\t\t\t\t%s;\n\t\t\t\t%s;\n\t\t\t\tmml.Nop(%s);\n\t\t\t\t%s;\n\t\t\t\treturn nil\n\t\t\t},\n\t\t\tFixedArgs: %d,\n\t\t}" } else { return "&mml.Function{\n\t\t\tF: func(a []interface{}) interface{} {\n\t\t\t\tvar c interface{}\n\t\t\t\tmml.Nop(c)\n\t\t\t\t%s;\n\t\t\t\t%s;\n\t\t\t\tmml.Nop(%s);\n\t\t\t\treturn %s\n\t\t\t},\n\t\t\tFixedArgs: %d,\n\t\t}" } }(), _compileParamList.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "params"), mml.Ref(_f, "collectParam"))), _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s interface{}", _s))
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _scope)))), _join.(*mml.Function).Call(append([]interface{}{}, ", ", _map.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "formatOne").(*mml.Function).Call(append([]interface{}{}, "_%s")), append(append([]interface{}{}, _scope.([]interface{})...), _paramNames.([]interface{})...))))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "statement"))), _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_f, "params")))));
				return nil
			},
			FixedArgs: 1,
		};
_compileRangeExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _e = a[0];
				;
				mml.Nop(_e);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "%s:%s", func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "from", _e)); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "from"))) } else { return "" } }(), func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "to", _e)); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_e, "to"))) } else { return "" } }()))
			},
			FixedArgs: 1,
		};
_compileIndexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return func () interface{} { c = (!_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_i, "index"))).(bool) || mml.BinaryOp(12, mml.Ref(mml.Ref(_i, "index"), "type"), "range-expression").(bool)); if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.Ref(%s, %s)", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_i, "expression"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_i, "index"))))) } else { return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.RefRange(%s, %s, %s)", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_i, "expression"))), func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "from", mml.Ref(_i, "index"))); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "from"))) } else { return "nil" } }(), func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "to", mml.Ref(_i, "index"))); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_i, "index"), "to"))) } else { return "nil" } }())) } }()
			},
			FixedArgs: 1,
		};
_compileApplication = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_a, "function"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_a, "function"), "type"), "function").(bool)); if c.(bool) { return "(%s).Call(%s)" } else { return "%s.(*mml.Function).Call(%s)" } }(), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "function"))), _compileList.(*mml.Function).Call(append([]interface{}{}, func() interface{} { s := make(map[string]interface{}); s["values"] = mml.Ref(_a, "args");; return s }()))))
			},
			FixedArgs: 1,
		};
_compileUnary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return func () interface{} { c = mml.BinaryOp(11, mml.Ref(_u, "op"), mml.Ref(_code, "logicalNot")); if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = _isBool.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "arg"))); if c.(bool) { return "!%s" } else { return "!%s.(bool)" } }(), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "arg"))))) } else { return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.UnaryOp(%d, %s)", mml.Ref(_u, "op"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "arg"))))) } }()
			},
			FixedArgs: 1,
		};
_compileBinary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _b = a[0];
				;
				mml.Nop(_b);
				var _isBoolOp interface{};
var _left interface{};
var _right interface{};
var _op interface{};
mml.Nop(_isBoolOp, _left, _right, _op);
c = (mml.BinaryOp(12, mml.Ref(_b, "op"), mml.Ref(_code, "logicalAnd")).(bool) && mml.BinaryOp(12, mml.Ref(_b, "op"), mml.Ref(_code, "logicalOr")).(bool)); if c.(bool) { ;
mml.Nop();
return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.BinaryOp(%s, %s, %s)", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "op"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right"))))) };
_isBoolOp = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return ((_has.(*mml.Function).Call(append([]interface{}{}, "type", _c)).(bool) && (mml.BinaryOp(11, mml.Ref(_c, "type"), "unary").(bool) && mml.BinaryOp(11, mml.Ref(_c, "op"), mml.Ref(_code, "logicalNot")).(bool))) || (mml.BinaryOp(11, mml.Ref(_c, "type"), "binary").(bool) && (mml.BinaryOp(11, mml.Ref(_c, "op"), mml.Ref(_code, "logicalAnd")).(bool) || mml.BinaryOp(11, mml.Ref(_c, "op"), mml.Ref(_code, "logicalOr")).(bool))))
			},
			FixedArgs: 1,
		};
_left = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left")));
_right = _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right")));
c = (!_isBool.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left"))).(bool) && !_isBoolOp.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "left"))).(bool)); if c.(bool) { ;
mml.Nop();
_left = mml.BinaryOp(9, _left, ".(bool)") };
c = (!_isBool.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right"))).(bool) && !_isBoolOp.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_b, "right"))).(bool)); if c.(bool) { ;
mml.Nop();
_right = mml.BinaryOp(9, _right, ".(bool)") };
_op = "&&";
c = mml.BinaryOp(11, mml.Ref(_b, "op"), mml.Ref(_code, "logicalOr")); if c.(bool) { ;
mml.Nop();
_op = "||" };
return _formats.(*mml.Function).Call(append([]interface{}{}, "(%s %s %s)", _left, _op, _right));
				return nil
			},
			FixedArgs: 1,
		};
_compileTernary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "func () interface{} { c = %s; if c.(bool) { return %s } else { return %s } }()", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "condition"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "consequent"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "alternative")))))
			},
			FixedArgs: 1,
		};
_compileIf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "alternative", _c)); if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, "c = %s; if c.(bool) { %s } else { %s }", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "condition"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "consequent"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "alternative"))))) } else { return _formats.(*mml.Function).Call(append([]interface{}{}, "c = %s; if c.(bool) { %s }", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "condition"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, "consequent"))))) } }()
			},
			FixedArgs: 1,
		};
_compileSwitch = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _hasDefault interface{};
var _cases interface{};
var _def interface{};
var _defaultCode interface{};
mml.Nop(_hasDefault, _cases, _def, _defaultCode);
_hasDefault = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_s, "defaultStatements"), "statements"))), 0);
_cases = _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "cases")));
_def = func () interface{} { c = _hasDefault; if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "defaultStatements"))) } else { return "" } }();
_defaultCode = func () interface{} { c = _hasDefault; if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, "default:\n%s", _def)) } else { return "" } }();
return _formats.(*mml.Function).Call(append([]interface{}{}, "switch %s {\n%s\n}", func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "expression", _s)); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "expression"))) } else { return "" } }(), _join.(*mml.Function).Call(append([]interface{}{}, "\n")).(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = _hasDefault; if c.(bool) { return append(append([]interface{}{}, _cases.([]interface{})...), _defaultCode) } else { return _cases } }()))));
				return nil
			},
			FixedArgs: 1,
		};
_compileSelect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return (&mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return mml.Ref(_strings, "formatOne").(*mml.Function).Call(append([]interface{}{}, "func() interface{} {\nselect {\n%s\n} }()")).(*mml.Function).Call(append([]interface{}{}, _join.(*mml.Function).Call(append([]interface{}{}, "\n")).(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = mml.Ref(_s, "hasDefault"); if c.(bool) { return append(append([]interface{}{}, _c.([]interface{})...), mml.Ref(_strings, "formatOne").(*mml.Function).Call(append([]interface{}{}, "default:\n%s")).(*mml.Function).Call(append([]interface{}{}, _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "defaultStatements")))))) } else { return _c } }()))))
			},
			FixedArgs: 1,
		}).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "cases")))))
			},
			FixedArgs: 1,
		};
_compileRangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _r = a[0];
				;
				mml.Nop(_r);
				var _infiniteCounter interface{};
var _withRangeExpression interface{};
var _listStyleRange interface{};
mml.Nop(_infiniteCounter, _withRangeExpression, _listStyleRange);
_infiniteCounter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s := 0; true; _%s++", mml.Ref(_r, "symbol"), mml.Ref(_r, "symbol")))
			},
			FixedArgs: 0,
		};
_withRangeExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s := %s; %s; _%s++", mml.Ref(_r, "symbol"), func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "from", mml.Ref(_r, "expression"))); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_r, "expression"), "from"))) } else { return "0" } }(), func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "to", mml.Ref(_r, "expression"))); if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s < %s.(int)", mml.Ref(_r, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_r, "expression"), "to"))))) } else { return "true" } }(), mml.Ref(_r, "symbol")))
			},
			FixedArgs: 0,
		};
_listStyleRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				return _formats.(*mml.Function).Call(append([]interface{}{}, "_, _%s := range %s.([]interface{})", mml.Ref(_r, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_r, "expression")))))
			},
			FixedArgs: 0,
		};
switch  {
case !_has.(*mml.Function).Call(append([]interface{}{}, "expression", _r)).(bool):
;
mml.Nop();
return _infiniteCounter.(*mml.Function).Call([]interface{}{})
case (_has.(*mml.Function).Call(append([]interface{}{}, "type", mml.Ref(_r, "expression"))).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_r, "expression"), "type"), "range-expression").(bool)):
;
mml.Nop();
return _withRangeExpression.(*mml.Function).Call([]interface{}{})
default:
;
mml.Nop();
return _listStyleRange.(*mml.Function).Call([]interface{}{})
};
				return nil
			},
			FixedArgs: 1,
		};
_compileLoop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "for %s {\n%s\n}", func () interface{} { c = _has.(*mml.Function).Call(append([]interface{}{}, "expression", _l)); if c.(bool) { return _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "expression"))) } else { return "" } }(), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, "body")))))
			},
			FixedArgs: 1,
		};
_compileDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return func () interface{} { c = mml.Ref(_d, "exported"); if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = %s; exports[\"%s\"] = _%s", mml.Ref(_d, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_d, "expression"))), mml.Ref(_d, "symbol"), mml.Ref(_d, "symbol"))) } else { return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = %s", mml.Ref(_d, "symbol"), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_d, "expression"))))) } }()
			},
			FixedArgs: 1,
		};
_compileAssign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = mml.BinaryOp(11, mml.Ref(mml.Ref(_a, "capture"), "type"), "symbol"); if c.(bool) { return _formats.(*mml.Function).Call(append([]interface{}{}, "%s = %s", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "capture"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "value"))))) } else { return _formats.(*mml.Function).Call(append([]interface{}{}, "mml.SetRef(%s, %s, %s)", _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_a, "capture"), "expression"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_a, "capture"), "index"))), _compile.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_a, "value"))))) } }()
			},
			FixedArgs: 1,
		};
_compileStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _scope interface{};
var _scopeNames interface{};
var _statements interface{};
var _scopeDefs interface{};
mml.Nop(_scope, _scopeNames, _statements, _scopeDefs);
_scope = _getScope.(*mml.Function).Call(append([]interface{}{}, _s.([]interface{})...));
_scopeNames = _join.(*mml.Function).Call(append([]interface{}{}, ", ", _map.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "formatOne").(*mml.Function).Call(append([]interface{}{}, "_%s")), _scope))));
_statements = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _notEmpty.(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _compile, _s))))));
_scopeDefs = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "var _%s interface{}", _s))
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _scope))));
return _formats.(*mml.Function).Call(append([]interface{}{}, "%s;\nmml.Nop(%s);\n%s", _scopeDefs, _scopeNames, _statements));
				return nil
			},
			FixedArgs: 1,
		};
_compileModule = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
				;
				mml.Nop(_m);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "%s", _compileStatements.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_m, "statements")))))
			},
			FixedArgs: 1,
		};
_getModuleName = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0];
				;
				mml.Nop(_path);
				return _path
			},
			FixedArgs: 1,
		};
_compileUse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				;
mml.Nop();
switch  {
case mml.BinaryOp(11, mml.Ref(_u, "capture"), "."):
var _useStatement interface{};
var _assigns interface{};
mml.Nop(_useStatement, _assigns);
_useStatement = _formats.(*mml.Function).Call(append([]interface{}{}, "var __%s = mml.Modules.Use(\"%s.mml\");", _getModuleName.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "path"))), mml.Ref(_u, "path")));
_assigns = _join.(*mml.Function).Call(append([]interface{}{}, ";\n")).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _name = a[0];
				;
				mml.Nop(_name);
				return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = __%s[\"%s\"]", _name, _getModuleName.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "path"))), _name))
			},
			FixedArgs: 1,
		}, mml.Ref(_u, "exportNames")))));
return _joins.(*mml.Function).Call(append([]interface{}{}, ";", _useStatement, _assigns))
case mml.BinaryOp(12, mml.Ref(_u, "capture"), ""):
;
mml.Nop();
return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = mml.Modules.Use(\"%s.mml\")", mml.Ref(_u, "capture"), mml.Ref(_u, "path")))
default:
;
mml.Nop();
return _formats.(*mml.Function).Call(append([]interface{}{}, "_%s = mml.Modules.Use(\"%s.mml\")", _getModuleName.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_u, "path"))), mml.Ref(_u, "path")))
};
				return nil
			},
			FixedArgs: 1,
		};
_compile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _code = a[0];
				;
				mml.Nop(_code);
				;
mml.Nop();
switch  {
case _isInt.(*mml.Function).Call(append([]interface{}{}, _code)):
;
mml.Nop();
return _compileInt.(*mml.Function).Call(append([]interface{}{}, _code))
case _isFloat.(*mml.Function).Call(append([]interface{}{}, _code)):
;
mml.Nop();
return _compileFloat.(*mml.Function).Call(append([]interface{}{}, _code))
case _isString.(*mml.Function).Call(append([]interface{}{}, _code)):
;
mml.Nop();
return _compileString.(*mml.Function).Call(append([]interface{}{}, _code))
case _isBool.(*mml.Function).Call(append([]interface{}{}, _code)):
;
mml.Nop();
return _compileBool.(*mml.Function).Call(append([]interface{}{}, _code))
};
switch mml.Ref(_code, "type") {
case "comment":
;
mml.Nop();
return _compileComment.(*mml.Function).Call(append([]interface{}{}, _code))
case "symbol":
;
mml.Nop();
return _compileSymbol.(*mml.Function).Call(append([]interface{}{}, _code))
case "module":
;
mml.Nop();
return _compileModule.(*mml.Function).Call(append([]interface{}{}, _code))
case "list":
;
mml.Nop();
return _compileList.(*mml.Function).Call(append([]interface{}{}, _code))
case "entry":
;
mml.Nop();
return _compileEntry.(*mml.Function).Call(append([]interface{}{}, _code))
case "struct":
;
mml.Nop();
return _compileStruct.(*mml.Function).Call(append([]interface{}{}, _code))
case "function":
;
mml.Nop();
return _compileFunction.(*mml.Function).Call(append([]interface{}{}, _code))
case "range-expression":
;
mml.Nop();
return _compileRangeExpression.(*mml.Function).Call(append([]interface{}{}, _code))
case "indexer":
;
mml.Nop();
return _compileIndexer.(*mml.Function).Call(append([]interface{}{}, _code))
case "spread":
;
mml.Nop();
return _compileSpread.(*mml.Function).Call(append([]interface{}{}, _code))
case "function-application":
;
mml.Nop();
return _compileApplication.(*mml.Function).Call(append([]interface{}{}, _code))
case "unary":
;
mml.Nop();
return _compileUnary.(*mml.Function).Call(append([]interface{}{}, _code))
case "binary":
;
mml.Nop();
return _compileBinary.(*mml.Function).Call(append([]interface{}{}, _code))
case "cond":
;
mml.Nop();
return _compileCond.(*mml.Function).Call(append([]interface{}{}, _code))
case "switch-case":
;
mml.Nop();
return _compileCase.(*mml.Function).Call(append([]interface{}{}, _code))
case "switch-statement":
;
mml.Nop();
return _compileSwitch.(*mml.Function).Call(append([]interface{}{}, _code))
case "send":
;
mml.Nop();
return _compileSend.(*mml.Function).Call(append([]interface{}{}, _code))
case "receive":
;
mml.Nop();
return _compileReceive.(*mml.Function).Call(append([]interface{}{}, _code))
case "go":
;
mml.Nop();
return _compileGo.(*mml.Function).Call(append([]interface{}{}, _code))
case "defer":
;
mml.Nop();
return _compileDefer.(*mml.Function).Call(append([]interface{}{}, _code))
case "select-case":
;
mml.Nop();
return _compileCase.(*mml.Function).Call(append([]interface{}{}, _code))
case "select":
;
mml.Nop();
return _compileSelect.(*mml.Function).Call(append([]interface{}{}, _code))
case "range-over":
;
mml.Nop();
return _compileRangeOver.(*mml.Function).Call(append([]interface{}{}, _code))
case "loop":
;
mml.Nop();
return _compileLoop.(*mml.Function).Call(append([]interface{}{}, _code))
case "definition":
;
mml.Nop();
return _compileDefinition.(*mml.Function).Call(append([]interface{}{}, _code))
case "definition-list":
;
mml.Nop();
return _compileDefinitions.(*mml.Function).Call(append([]interface{}{}, _code))
case "assign":
;
mml.Nop();
return _compileAssign.(*mml.Function).Call(append([]interface{}{}, _code))
case "assign-list":
;
mml.Nop();
return _compileAssigns.(*mml.Function).Call(append([]interface{}{}, _code))
case "ret":
;
mml.Nop();
return _compileRet.(*mml.Function).Call(append([]interface{}{}, _code))
case "control-statement":
;
mml.Nop();
return _compileControl.(*mml.Function).Call(append([]interface{}{}, _code))
case "use":
;
mml.Nop();
return _compileUse.(*mml.Function).Call(append([]interface{}{}, _code))
case "use-list":
;
mml.Nop();
return _compileUseList.(*mml.Function).Call(append([]interface{}{}, _code))
default:
;
mml.Nop();
return _compileStatements.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_code, "statements")))
};
				return nil
			},
			FixedArgs: 1,
		}; exports["compile"] = _compile
		return exports
	})
modulePath = "lang.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
var _logger interface{};
var _list interface{};
var _strings interface{};
var _ints interface{};
var _errors interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _logger, _list, _strings, _ints, _errors);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_ints = mml.Modules.Use("ints.mml");
_logger = mml.Modules.Use("log.mml");
_errors = mml.Modules.Use("errors.mml");
_fold = mml.Ref(_list, "fold"); exports["fold"] = _fold;
_foldr = mml.Ref(_list, "foldr"); exports["foldr"] = _foldr;
_map = mml.Ref(_list, "map"); exports["map"] = _map;
_filter = mml.Ref(_list, "filter"); exports["filter"] = _filter;
_contains = mml.Ref(_list, "contains"); exports["contains"] = _contains;
_sort = mml.Ref(_list, "sort"); exports["sort"] = _sort;
_flat = mml.Ref(_list, "flat"); exports["flat"] = _flat;
_join = mml.Ref(_strings, "join"); exports["join"] = _join;
_joins = mml.Ref(_strings, "joins"); exports["joins"] = _joins;
_formats = mml.Ref(_strings, "formats"); exports["formats"] = _formats;
_enum = mml.Ref(_ints, "enum"); exports["enum"] = _enum;
_log = mml.Ref(_logger, "log"); exports["log"] = _log;
_onlyErr = mml.Ref(_errors, "only"); exports["onlyErr"] = _onlyErr;
_passErr = mml.Ref(_errors, "pass"); exports["passErr"] = _passErr
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "ints.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _counter interface{};
var _enum interface{};
mml.Nop(_counter, _enum);
_counter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _c interface{};
mml.Nop(_c);
_c = mml.UnaryOp(2, 1);
return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				;
mml.Nop();
_c = mml.BinaryOp(9, _c, 1);
return _c;
				return nil
			},
			FixedArgs: 0,
		};
				return nil
			},
			FixedArgs: 0,
		}; exports["counter"] = _counter;
_enum = _counter; exports["enum"] = _enum
		return exports
	})
modulePath = "log.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _log interface{};
var _list interface{};
var _strings interface{};
mml.Nop(_log, _list, _strings);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_log = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0:];
				;
				mml.Nop(_a);
				;
mml.Nop();
_stderr.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call(append([]interface{}{}, " ")).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_list, "map").(*mml.Function).Call(append([]interface{}{}, _string)).(*mml.Function).Call(append([]interface{}{}, _a))))));
_stderr.(*mml.Function).Call(append([]interface{}{}, "\n"));
return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 0); if c.(bool) { return "" } else { return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 1)) } }();
				return nil
			},
			FixedArgs: 0,
		}; exports["log"] = _log
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "code.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _controlStatement interface{};
var _breakControl interface{};
var _continueControl interface{};
var _unaryOp interface{};
var _binaryNot interface{};
var _plus interface{};
var _minus interface{};
var _logicalNot interface{};
var _binaryOp interface{};
var _binaryAnd interface{};
var _binaryOr interface{};
var _xor interface{};
var _andNot interface{};
var _lshift interface{};
var _rshift interface{};
var _mul interface{};
var _div interface{};
var _mod interface{};
var _add interface{};
var _sub interface{};
var _eq interface{};
var _notEq interface{};
var _less interface{};
var _lessOrEq interface{};
var _greater interface{};
var _greaterOrEq interface{};
var _logicalAnd interface{};
var _logicalOr interface{};
var _flattenedStatements interface{};
var _list interface{};
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
mml.Nop(_controlStatement, _breakControl, _continueControl, _unaryOp, _binaryNot, _plus, _minus, _logicalNot, _binaryOp, _binaryAnd, _binaryOr, _xor, _andNot, _lshift, _rshift, _mul, _div, _mod, _add, _sub, _eq, _notEq, _less, _lessOrEq, _greater, _greaterOrEq, _logicalAnd, _logicalOr, _flattenedStatements, _list, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr);
var __lang = mml.Modules.Use("lang.mml");;_fold = __lang["fold"];
_foldr = __lang["foldr"];
_map = __lang["map"];
_filter = __lang["filter"];
_contains = __lang["contains"];
_sort = __lang["sort"];
_flat = __lang["flat"];
_join = __lang["join"];
_joins = __lang["joins"];
_formats = __lang["formats"];
_enum = __lang["enum"];
_log = __lang["log"];
_onlyErr = __lang["onlyErr"];
_passErr = __lang["passErr"];
_list = mml.Modules.Use("list.mml");
_controlStatement = _enum.(*mml.Function).Call([]interface{}{}); exports["controlStatement"] = _controlStatement;
_breakControl = _controlStatement.(*mml.Function).Call([]interface{}{}); exports["breakControl"] = _breakControl;
_continueControl = _controlStatement.(*mml.Function).Call([]interface{}{}); exports["continueControl"] = _continueControl;
_unaryOp = _enum.(*mml.Function).Call([]interface{}{}); exports["unaryOp"] = _unaryOp;
_binaryNot = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["binaryNot"] = _binaryNot;
_plus = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["plus"] = _plus;
_minus = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["minus"] = _minus;
_logicalNot = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["logicalNot"] = _logicalNot;
_binaryOp = _enum.(*mml.Function).Call([]interface{}{}); exports["binaryOp"] = _binaryOp;
_binaryAnd = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["binaryAnd"] = _binaryAnd;
_binaryOr = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["binaryOr"] = _binaryOr;
_xor = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["xor"] = _xor;
_andNot = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["andNot"] = _andNot;
_lshift = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["lshift"] = _lshift;
_rshift = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["rshift"] = _rshift;
_mul = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["mul"] = _mul;
_div = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["div"] = _div;
_mod = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["mod"] = _mod;
_add = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["add"] = _add;
_sub = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["sub"] = _sub;
_eq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["eq"] = _eq;
_notEq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["notEq"] = _notEq;
_less = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["less"] = _less;
_lessOrEq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["lessOrEq"] = _lessOrEq;
_greater = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["greater"] = _greater;
_greaterOrEq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["greaterOrEq"] = _greaterOrEq;
_logicalAnd = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["logicalAnd"] = _logicalAnd;
_logicalOr = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["logicalOr"] = _logicalOr;
_flattenedStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _itemType = a[0];
var _listType = a[1];
var _listProp = a[2];
var _statements = a[3];
				;
				mml.Nop(_itemType, _listType, _listProp, _statements);
				var _type interface{};
var _toList interface{};
mml.Nop(_type, _toList);
_type = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return (_has.(*mml.Function).Call(append([]interface{}{}, "type", _s)).(bool) && _contains.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "type"), append([]interface{}{}, _itemType, _listType))).(bool))
			},
			FixedArgs: 1,
		};
_toList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return func () interface{} { c = mml.BinaryOp(11, mml.Ref(_s, "type"), _itemType); if c.(bool) { return append([]interface{}{}, _s) } else { return mml.Ref(_s, _listProp) } }()
			},
			FixedArgs: 1,
		};
return mml.Ref(_list, "flat").(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _toList)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _type)).(*mml.Function).Call(append([]interface{}{}, _statements))))));
				return nil
			},
			FixedArgs: 4,
		}; exports["flattenedStatements"] = _flattenedStatements
		return exports
	})
modulePath = "lang.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
var _logger interface{};
var _list interface{};
var _strings interface{};
var _ints interface{};
var _errors interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _logger, _list, _strings, _ints, _errors);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_ints = mml.Modules.Use("ints.mml");
_logger = mml.Modules.Use("log.mml");
_errors = mml.Modules.Use("errors.mml");
_fold = mml.Ref(_list, "fold"); exports["fold"] = _fold;
_foldr = mml.Ref(_list, "foldr"); exports["foldr"] = _foldr;
_map = mml.Ref(_list, "map"); exports["map"] = _map;
_filter = mml.Ref(_list, "filter"); exports["filter"] = _filter;
_contains = mml.Ref(_list, "contains"); exports["contains"] = _contains;
_sort = mml.Ref(_list, "sort"); exports["sort"] = _sort;
_flat = mml.Ref(_list, "flat"); exports["flat"] = _flat;
_join = mml.Ref(_strings, "join"); exports["join"] = _join;
_joins = mml.Ref(_strings, "joins"); exports["joins"] = _joins;
_formats = mml.Ref(_strings, "formats"); exports["formats"] = _formats;
_enum = mml.Ref(_ints, "enum"); exports["enum"] = _enum;
_log = mml.Ref(_logger, "log"); exports["log"] = _log;
_onlyErr = mml.Ref(_errors, "only"); exports["onlyErr"] = _onlyErr;
_passErr = mml.Ref(_errors, "pass"); exports["passErr"] = _passErr
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "ints.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _counter interface{};
var _enum interface{};
mml.Nop(_counter, _enum);
_counter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _c interface{};
mml.Nop(_c);
_c = mml.UnaryOp(2, 1);
return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				;
mml.Nop();
_c = mml.BinaryOp(9, _c, 1);
return _c;
				return nil
			},
			FixedArgs: 0,
		};
				return nil
			},
			FixedArgs: 0,
		}; exports["counter"] = _counter;
_enum = _counter; exports["enum"] = _enum
		return exports
	})
modulePath = "log.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _log interface{};
var _list interface{};
var _strings interface{};
mml.Nop(_log, _list, _strings);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_log = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0:];
				;
				mml.Nop(_a);
				;
mml.Nop();
_stderr.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call(append([]interface{}{}, " ")).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_list, "map").(*mml.Function).Call(append([]interface{}{}, _string)).(*mml.Function).Call(append([]interface{}{}, _a))))));
_stderr.(*mml.Function).Call(append([]interface{}{}, "\n"));
return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 0); if c.(bool) { return "" } else { return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 1)) } }();
				return nil
			},
			FixedArgs: 0,
		}; exports["log"] = _log
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "parse.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _parseString interface{};
var _parseSpread interface{};
var _parseExpressionList interface{};
var _parseList interface{};
var _parseMutableList interface{};
var _parseExpressionKey interface{};
var _parseEntry interface{};
var _parseStruct interface{};
var _parseMutableStruct interface{};
var _parseStatementList interface{};
var _parseFunction interface{};
var _parseEffect interface{};
var _parseSymbolIndex interface{};
var _parseExpressionIndex interface{};
var _parseIndexer interface{};
var _parseMutableCapture interface{};
var _parseValueDefinition interface{};
var _parseFunctionDefinition interface{};
var _parseAssign interface{};
var _parseSend interface{};
var _parseReceive interface{};
var _parseGo interface{};
var _parseDefer interface{};
var _parseReceiveDefinition interface{};
var _parseSymbol interface{};
var _parseReturn interface{};
var _parseFunctionFact interface{};
var _parseRange interface{};
var _parseRangeIndex interface{};
var _parseIndexerNodes interface{};
var _parseFunctionApplication interface{};
var _parseUnaryExpression interface{};
var _parseBinaryExpression interface{};
var _parseChaining interface{};
var _parserTernary interface{};
var _parseIf interface{};
var _parseSwitch interface{};
var _parseRangeOver interface{};
var _parseLoop interface{};
var _parseValueCapture interface{};
var _parseDefinitions interface{};
var _parseMutableDefinitions interface{};
var _parseFunctionCapture interface{};
var _parseEffectCapture interface{};
var _parseEffectDefinitions interface{};
var _parseAssignCaptures interface{};
var _parseSelect interface{};
var _parseExport interface{};
var _parseUseFact interface{};
var _parseUse interface{};
var _parse interface{};
var _parseFile interface{};
var _findExportNames interface{};
var _parseModules interface{};
var _code interface{};
var _strings interface{};
var _errors interface{};
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
mml.Nop(_parseString, _parseSpread, _parseExpressionList, _parseList, _parseMutableList, _parseExpressionKey, _parseEntry, _parseStruct, _parseMutableStruct, _parseStatementList, _parseFunction, _parseEffect, _parseSymbolIndex, _parseExpressionIndex, _parseIndexer, _parseMutableCapture, _parseValueDefinition, _parseFunctionDefinition, _parseAssign, _parseSend, _parseReceive, _parseGo, _parseDefer, _parseReceiveDefinition, _parseSymbol, _parseReturn, _parseFunctionFact, _parseRange, _parseRangeIndex, _parseIndexerNodes, _parseFunctionApplication, _parseUnaryExpression, _parseBinaryExpression, _parseChaining, _parserTernary, _parseIf, _parseSwitch, _parseRangeOver, _parseLoop, _parseValueCapture, _parseDefinitions, _parseMutableDefinitions, _parseFunctionCapture, _parseEffectCapture, _parseEffectDefinitions, _parseAssignCaptures, _parseSelect, _parseExport, _parseUseFact, _parseUse, _parse, _parseFile, _findExportNames, _parseModules, _code, _strings, _errors, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr);
var __lang = mml.Modules.Use("lang.mml");;_fold = __lang["fold"];
_foldr = __lang["foldr"];
_map = __lang["map"];
_filter = __lang["filter"];
_contains = __lang["contains"];
_sort = __lang["sort"];
_flat = __lang["flat"];
_join = __lang["join"];
_joins = __lang["joins"];
_formats = __lang["formats"];
_enum = __lang["enum"];
_log = __lang["log"];
_onlyErr = __lang["onlyErr"];
_passErr = __lang["passErr"];
_code = mml.Modules.Use("code.mml");
_strings = mml.Modules.Use("strings.mml");
_errors = mml.Modules.Use("errors.mml");
_parseString = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return mml.Ref(_strings, "unescape").(*mml.Function).Call(append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "text"), 1, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "text"))), 1))))
			},
			FixedArgs: 1,
		};
_parseSpread = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "spread";s["value"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }()
			},
			FixedArgs: 1,
		};
_parseExpressionList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0];
				;
				mml.Nop(_nodes);
				return _map.(*mml.Function).Call(append([]interface{}{}, _parse, _nodes))
			},
			FixedArgs: 1,
		};
_parseList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "list";s["values"] = _parseExpressionList.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes")));s["mutable"] = false;; return s }()
			},
			FixedArgs: 1,
		};
_parseMutableList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); sp := _parseList.(*mml.Function).Call(append([]interface{}{}, _ast)).(map[string]interface{});
for k, v := range sp { s[k] = v };s["mutable"] = true;; return s }()
			},
			FixedArgs: 1,
		};
_parseExpressionKey = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "expression-key";s["value"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }()
			},
			FixedArgs: 1,
		};
_parseEntry = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "entry";s["key"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));s["value"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));; return s }()
			},
			FixedArgs: 1,
		};
_parseStruct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "struct";s["entries"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.Ref(_ast, "nodes")));s["mutable"] = false;; return s }()
			},
			FixedArgs: 1,
		};
_parseMutableStruct = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); sp := _parseStruct.(*mml.Function).Call(append([]interface{}{}, _ast)).(map[string]interface{});
for k, v := range sp { s[k] = v };s["mutable"] = true;; return s }()
			},
			FixedArgs: 1,
		};
_parseStatementList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "statement-list";s["statements"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.Ref(_ast, "nodes")));; return s }()
			},
			FixedArgs: 1,
		};
_parseFunction = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return _parseFunctionFact.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes")))
			},
			FixedArgs: 1,
		};
_parseEffect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); sp := _parseFunction.(*mml.Function).Call(append([]interface{}{}, _ast)).(map[string]interface{});
for k, v := range sp { s[k] = v };s["effect"] = true;; return s }()
			},
			FixedArgs: 1,
		};
_parseSymbolIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))), "name")
			},
			FixedArgs: 1,
		};
_parseExpressionIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)))
			},
			FixedArgs: 1,
		};
_parseIndexer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return _parseIndexerNodes.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes")))
			},
			FixedArgs: 1,
		};
_parseMutableCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); sp := _parseValueCapture.(*mml.Function).Call(append([]interface{}{}, _ast)).(map[string]interface{});
for k, v := range sp { s[k] = v };s["mutable"] = true;; return s }()
			},
			FixedArgs: 1,
		};
_parseValueDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)))
			},
			FixedArgs: 1,
		};
_parseFunctionDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)))
			},
			FixedArgs: 1,
		};
_parseAssign = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "assign-list";s["assignments"] = _parseAssignCaptures.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes")));; return s }()
			},
			FixedArgs: 1,
		};
_parseSend = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "send";s["channel"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));s["value"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));; return s }()
			},
			FixedArgs: 1,
		};
_parseReceive = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "receive";s["channel"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }()
			},
			FixedArgs: 1,
		};
_parseGo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "go";s["application"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }()
			},
			FixedArgs: 1,
		};
_parseDefer = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "defer";s["application"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }()
			},
			FixedArgs: 1,
		};
_parseReceiveDefinition = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return _parseValueCapture.(*mml.Function).Call(append([]interface{}{}, _ast))
			},
			FixedArgs: 1,
		};
_parseSymbol = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				;
mml.Nop();
switch mml.Ref(_ast, "text") {
case "break":
;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "control-statement";s["control"] = mml.Ref(_code, "breakControl");; return s }()
case "continue":
;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "control-statement";s["control"] = mml.Ref(_code, "continueControl");; return s }()
default:
;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "symbol";s["name"] = mml.Ref(_ast, "text");; return s }()
};
				return nil
			},
			FixedArgs: 1,
		};
_parseReturn = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 0); if c.(bool) { return func() interface{} { s := make(map[string]interface{}); s["type"] = "ret";; return s }() } else { return func() interface{} { s := make(map[string]interface{}); s["type"] = "ret";s["value"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }() } }()
			},
			FixedArgs: 1,
		};
_parseFunctionFact = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0];
				;
				mml.Nop(_nodes);
				var _last interface{};
var _params interface{};
var _lastParam interface{};
var _hasCollectParam interface{};
var _fixedParams interface{};
mml.Nop(_last, _params, _lastParam, _hasCollectParam, _fixedParams);
_last = mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _nodes)), 1);
_params = mml.RefRange(_nodes, nil, _last);
_lastParam = mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _params)), 1);
_hasCollectParam = (mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _params)), 0).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(_params, _lastParam), "name"), "collect-parameter").(bool));
_fixedParams = func () interface{} { c = _hasCollectParam; if c.(bool) { return mml.RefRange(_params, nil, _lastParam) } else { return _params } }();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "function";s["params"] = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
				;
				mml.Nop(_p);
				return mml.Ref(_p, "name")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _parse)).(*mml.Function).Call(append([]interface{}{}, _fixedParams))));s["collectParam"] = func () interface{} { c = _hasCollectParam; if c.(bool) { return mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(mml.Ref(_params, _lastParam), "nodes"), 0))), "name") } else { return "" } }();s["statement"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_nodes, _last)));s["effect"] = false;; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseRange = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _v interface{};
mml.Nop(_v);
_v = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));
return func () interface{} { c = mml.BinaryOp(11, mml.Ref(_ast, "name"), "range-from"); if c.(bool) { return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-expression";s["from"] = _v;; return s }() } else { return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-expression";s["to"] = _v;; return s }() } }();
				return nil
			},
			FixedArgs: 1,
		};
_parseRangeIndex = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _r interface{};
mml.Nop(_r);
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 0); if c.(bool) { ;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-expression";; return s }() };
_r = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 1); if c.(bool) { ;
mml.Nop();
return _r };
return func() interface{} { s := make(map[string]interface{}); sp := _r.(map[string]interface{});
for k, v := range sp { s[k] = v };s["to"] = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1))), "to");; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseIndexerNodes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _n = a[0];
				;
				mml.Nop(_n);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "indexer";s["expression"] = func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _n)), 2); if c.(bool) { return _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_n, 0))) } else { return _parseIndexerNodes.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_n, nil, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _n)), 1)))) } }();s["index"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_n, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _n)), 1))));; return s }()
			},
			FixedArgs: 1,
		};
_parseFunctionApplication = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "function-application";s["function"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));s["args"] = _parseExpressionList.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil)));; return s }()
			},
			FixedArgs: 1,
		};
_parseUnaryExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _op interface{};
mml.Nop(_op);
_op = mml.Ref(_code, "binaryNot");
switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name") {
case "plus":
;
mml.Nop();
_op = mml.Ref(_code, "plus")
case "minus":
;
mml.Nop();
_op = mml.Ref(_code, "minus")
case "logical-not":
;
mml.Nop();
_op = mml.Ref(_code, "logicalNot")
};
return func() interface{} { s := make(map[string]interface{}); s["type"] = "unary";s["op"] = _op;s["arg"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseBinaryExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _op interface{};
mml.Nop(_op);
_op = mml.Ref(_code, "binaryAnd");
switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 2)), "name") {
case "xor":
;
mml.Nop();
_op = mml.Ref(_code, "xor")
case "and-not":
;
mml.Nop();
_op = mml.Ref(_code, "andNot")
case "lshift":
;
mml.Nop();
_op = mml.Ref(_code, "lshift")
case "rshift":
;
mml.Nop();
_op = mml.Ref(_code, "rshift")
case "mul":
;
mml.Nop();
_op = mml.Ref(_code, "mul")
case "div":
;
mml.Nop();
_op = mml.Ref(_code, "div")
case "mod":
;
mml.Nop();
_op = mml.Ref(_code, "mod")
case "add":
;
mml.Nop();
_op = mml.Ref(_code, "add")
case "sub":
;
mml.Nop();
_op = mml.Ref(_code, "sub")
case "eq":
;
mml.Nop();
_op = mml.Ref(_code, "eq")
case "not-eq":
;
mml.Nop();
_op = mml.Ref(_code, "notEq")
case "less":
;
mml.Nop();
_op = mml.Ref(_code, "less")
case "less-or-eq":
;
mml.Nop();
_op = mml.Ref(_code, "lessOrEq")
case "greater":
;
mml.Nop();
_op = mml.Ref(_code, "greater")
case "greater-or-eq":
;
mml.Nop();
_op = mml.Ref(_code, "greaterOrEq")
case "logical-and":
;
mml.Nop();
_op = mml.Ref(_code, "logicalAnd")
case "logical-or":
;
mml.Nop();
_op = mml.Ref(_code, "logicalOr")
};
return func() interface{} { s := make(map[string]interface{}); s["type"] = "binary";s["op"] = _op;s["left"] = _parse.(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 3); if c.(bool) { return func() interface{} { s := make(map[string]interface{}); s["name"] = mml.Ref(_ast, "name");s["nodes"] = mml.RefRange(mml.Ref(_ast, "nodes"), nil, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 2));; return s }() } else { return mml.Ref(mml.Ref(_ast, "nodes"), 0) } }()));s["right"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 1))));; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseChaining = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _a interface{};
var _n interface{};
mml.Nop(_a, _n);
_a = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));
_n = mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil);
for  {
;
mml.Nop();
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _n)), 0); if c.(bool) { ;
mml.Nop();
return _a };
_a = func() interface{} { s := make(map[string]interface{}); s["type"] = "function-application";s["function"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_n, 0)));s["args"] = append([]interface{}{}, _a);; return s }();
_n = mml.RefRange(_n, 1, nil)
};
				return nil
			},
			FixedArgs: 1,
		};
_parserTernary = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "cond";s["condition"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));s["consequent"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));s["alternative"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 2)));s["ternary"] = true;; return s }()
			},
			FixedArgs: 1,
		};
_parseIf = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _cond interface{};
var _alternative interface{};
mml.Nop(_cond, _alternative);
_cond = func() interface{} { s := make(map[string]interface{}); s["type"] = "cond";s["ternary"] = false;s["condition"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));s["consequent"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));; return s }();
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 2); if c.(bool) { ;
mml.Nop();
return _cond };
_alternative = func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 3); if c.(bool) { return _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 2))) } else { return _parse.(*mml.Function).Call(append([]interface{}{}, func() interface{} { s := make(map[string]interface{}); sp := _ast.(map[string]interface{});
for k, v := range sp { s[k] = v };s["nodes"] = mml.RefRange(mml.Ref(_ast, "nodes"), 2, nil);; return s }())) } }();
return func() interface{} { s := make(map[string]interface{}); sp := _cond.(map[string]interface{});
for k, v := range sp { s[k] = v };s["alternative"] = _alternative;; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseSwitch = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _hasExpression interface{};
var _expression interface{};
var _nodes interface{};
var _groupLines interface{};
var _parseCases interface{};
var _lines interface{};
var _s interface{};
mml.Nop(_hasExpression, _expression, _nodes, _groupLines, _parseCases, _lines, _s);
_hasExpression = (mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "case").(bool) && mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "default").(bool));
_expression = func () interface{} { c = _hasExpression; if c.(bool) { return mml.Ref(mml.Ref(_ast, "nodes"), 0) } else { return func() interface{} { s := make(map[string]interface{}); ; return s }() } }();
_nodes = func () interface{} { c = _hasExpression; if c.(bool) { return mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil) } else { return mml.Ref(_ast, "nodes") } }();
_groupLines = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _isDefault interface{};
var _current interface{};
var _cases interface{};
var _defaults interface{};
mml.Nop(_isDefault, _current, _cases, _defaults);
_isDefault = false;
_current = []interface{}{};
_cases = []interface{}{};
_defaults = []interface{}{};
for _, _n := range _nodes.([]interface{}) {
;
mml.Nop();
switch mml.Ref(_n, "name") {
case "case":
;
mml.Nop();
c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _current)), 0); if c.(bool) { ;
mml.Nop();
c = _isDefault; if c.(bool) { ;
mml.Nop();
_defaults = _current } else { ;
mml.Nop();
_cases = append(append([]interface{}{}, _cases.([]interface{})...), _current) } };
_current = append([]interface{}{}, mml.Ref(mml.Ref(_n, "nodes"), 0));
_isDefault = false
case "default":
;
mml.Nop();
c = (mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _current)), 0).(bool) && !_isDefault.(bool)); if c.(bool) { ;
mml.Nop();
_cases = append(append([]interface{}{}, _cases.([]interface{})...), _current) };
_current = []interface{}{};
_isDefault = true
default:
;
mml.Nop();
_current = append(append([]interface{}{}, _current.([]interface{})...), _n)
}
};
c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _current)), 0); if c.(bool) { ;
mml.Nop();
c = _isDefault; if c.(bool) { ;
mml.Nop();
_defaults = _current } else { ;
mml.Nop();
_cases = append(append([]interface{}{}, _cases.([]interface{})...), _current) } };
return func() interface{} { s := make(map[string]interface{}); s["cases"] = _cases;s["defaults"] = _defaults;; return s }();
				return nil
			},
			FixedArgs: 0,
		};
_parseCases = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				;
mml.Nop();
return _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "switch-case";s["expression"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, 0)));s["body"] = func() interface{} { s := make(map[string]interface{}); s["type"] = "statement-list";s["statements"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.RefRange(_c, 1, nil)));; return s }();; return s }()
			},
			FixedArgs: 1,
		}, _c));
				return nil
			},
			FixedArgs: 1,
		};
_lines = _groupLines.(*mml.Function).Call([]interface{}{});
_s = func() interface{} { s := make(map[string]interface{}); s["type"] = "switch-statement";s["cases"] = _parseCases.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_lines, "cases")));s["defaultStatements"] = func() interface{} { s := make(map[string]interface{}); s["type"] = "statement-list";s["statements"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.Ref(_lines, "defaults")));; return s }();; return s }();
return func () interface{} { c = _hasExpression; if c.(bool) { return func() interface{} { s := make(map[string]interface{}); sp := _s.(map[string]interface{});
for k, v := range sp { s[k] = v };s["expression"] = _parse.(*mml.Function).Call(append([]interface{}{}, _expression));; return s }() } else { return _s } }();
				return nil
			},
			FixedArgs: 1,
		};
_parseRangeOver = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _parseExpression interface{};
mml.Nop(_parseExpression);
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 0); if c.(bool) { ;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-over";; return s }() };
c = (mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 1).(bool) && mml.BinaryOp(11, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "symbol").(bool)); if c.(bool) { ;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-over";s["symbol"] = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))), "name");; return s }() };
_parseExpression = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0];
				;
				mml.Nop(_nodes);
				var _exp interface{};
mml.Nop(_exp);
_exp = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_nodes, 0)));
c = ((!_has.(*mml.Function).Call(append([]interface{}{}, "type", _exp)).(bool) || mml.BinaryOp(12, mml.Ref(_exp, "type"), "range-expression").(bool)) || mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _nodes)), 1).(bool)); if c.(bool) { ;
mml.Nop();
return _exp };
return func() interface{} { s := make(map[string]interface{}); sp := _exp.(map[string]interface{});
for k, v := range sp { s[k] = v };s["to"] = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_nodes, 1))), "to");; return s }();
				return nil
			},
			FixedArgs: 1,
		};
c = mml.BinaryOp(12, mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name"), "symbol"); if c.(bool) { ;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-over";s["expression"] = _parseExpression.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes")));; return s }() };
return func() interface{} { s := make(map[string]interface{}); s["type"] = "range-over";s["symbol"] = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))), "name");s["expression"] = _parseExpression.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil)));; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseLoop = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _loop interface{};
var _expression interface{};
var _emptyRange interface{};
mml.Nop(_loop, _expression, _emptyRange);
_loop = func() interface{} { s := make(map[string]interface{}); s["type"] = "loop";; return s }();
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes"))), 1); if c.(bool) { ;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); sp := _loop.(map[string]interface{});
for k, v := range sp { s[k] = v };s["body"] = _parseStatementList.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));; return s }() };
_expression = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));
_emptyRange = (((_has.(*mml.Function).Call(append([]interface{}{}, "type", _expression)).(bool) && mml.BinaryOp(11, mml.Ref(_expression, "type"), "range-over").(bool)) && !_has.(*mml.Function).Call(append([]interface{}{}, "symbol", _expression)).(bool)) && !_has.(*mml.Function).Call(append([]interface{}{}, "expression", _expression)).(bool));
return func () interface{} { c = _emptyRange; if c.(bool) { return func() interface{} { s := make(map[string]interface{}); sp := _loop.(map[string]interface{});
for k, v := range sp { s[k] = v };s["body"] = _parseStatementList.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));; return s }() } else { return func() interface{} { s := make(map[string]interface{}); sp := _loop.(map[string]interface{});
for k, v := range sp { s[k] = v };s["expression"] = _expression;s["body"] = _parseStatementList.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));; return s }() } }();
				return nil
			},
			FixedArgs: 1,
		};
_parseValueCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "definition";s["symbol"] = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))), "name");s["expression"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)));s["mutable"] = false;s["exported"] = false;; return s }()
			},
			FixedArgs: 1,
		};
_parseDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "definition-list";s["definitions"] = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return (!_has.(*mml.Function).Call(append([]interface{}{}, "type", _c)).(bool) || mml.BinaryOp(12, mml.Ref(_c, "type"), "comment").(bool))
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _parse)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "nodes")))));; return s }()
			},
			FixedArgs: 1,
		};
_parseMutableDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _dl interface{};
mml.Nop(_dl);
_dl = _parseDefinitions.(*mml.Function).Call(append([]interface{}{}, _ast));
return func() interface{} { s := make(map[string]interface{}); sp := _dl.(map[string]interface{});
for k, v := range sp { s[k] = v };s["definitions"] = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return func() interface{} { s := make(map[string]interface{}); sp := _d.(map[string]interface{});
for k, v := range sp { s[k] = v };s["mutable"] = true;; return s }()
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_dl, "definitions")));; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseFunctionCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "definition";s["symbol"] = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))), "name");s["expression"] = _parseFunctionFact.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(mml.Ref(_ast, "nodes"), 1, nil)));s["exported"] = false;; return s }()
			},
			FixedArgs: 1,
		};
_parseEffectCapture = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _f interface{};
mml.Nop(_f);
_f = _parseFunctionCapture.(*mml.Function).Call(append([]interface{}{}, _ast));
return func() interface{} { s := make(map[string]interface{}); sp := _f.(map[string]interface{});
for k, v := range sp { s[k] = v };s["expression"] = func() interface{} { s := make(map[string]interface{}); sp := mml.Ref(_f, "expression").(map[string]interface{});
for k, v := range sp { s[k] = v };s["effect"] = true;; return s }();; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseEffectDefinitions = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _dl interface{};
mml.Nop(_dl);
_dl = _parseDefinitions.(*mml.Function).Call(append([]interface{}{}, _ast));
return func() interface{} { s := make(map[string]interface{}); sp := _dl.(map[string]interface{});
for k, v := range sp { s[k] = v };s["definitions"] = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return func() interface{} { s := make(map[string]interface{}); sp := _d.(map[string]interface{});
for k, v := range sp { s[k] = v };s["effect"] = true;; return s }()
			},
			FixedArgs: 1,
		}, mml.Ref(_dl, "definitions")));; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseAssignCaptures = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _nodes = a[0];
				;
				mml.Nop(_nodes);
				;
mml.Nop();
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _nodes)), 0); if c.(bool) { ;
mml.Nop();
return []interface{}{} };
return append(append([]interface{}{}, func() interface{} { s := make(map[string]interface{}); s["type"] = "assign";s["capture"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_nodes, 0)));s["value"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_nodes, 1)));; return s }()), _parseAssignCaptures.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_nodes, 2, nil))).([]interface{})...);
				return nil
			},
			FixedArgs: 1,
		};
_parseSelect = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _nodes interface{};
var _groupLines interface{};
var _parseCases interface{};
var _lines interface{};
mml.Nop(_nodes, _groupLines, _parseCases, _lines);
_nodes = mml.Ref(_ast, "nodes");
_groupLines = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _isDefault interface{};
var _hasDefault interface{};
var _current interface{};
var _cases interface{};
var _defaults interface{};
mml.Nop(_isDefault, _hasDefault, _current, _cases, _defaults);
_isDefault = false;
_hasDefault = false;
_current = []interface{}{};
_cases = []interface{}{};
_defaults = []interface{}{};
for _, _n := range _nodes.([]interface{}) {
;
mml.Nop();
switch mml.Ref(_n, "name") {
case "case":
;
mml.Nop();
c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _current)), 0); if c.(bool) { ;
mml.Nop();
c = _isDefault; if c.(bool) { ;
mml.Nop();
_defaults = _current } else { ;
mml.Nop();
_cases = append(append([]interface{}{}, _cases.([]interface{})...), _current) } };
_current = append([]interface{}{}, mml.Ref(mml.Ref(_n, "nodes"), 0));
_isDefault = false
case "default":
;
mml.Nop();
c = (mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _current)), 0).(bool) && !_isDefault.(bool)); if c.(bool) { ;
mml.Nop();
_cases = append(append([]interface{}{}, _cases.([]interface{})...), _current) };
_current = []interface{}{};
_isDefault = true
default:
;
mml.Nop();
_current = append(append([]interface{}{}, _current.([]interface{})...), _n)
}
};
c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _current)), 0); if c.(bool) { ;
mml.Nop();
c = _isDefault; if c.(bool) { ;
mml.Nop();
_defaults = _current } else { ;
mml.Nop();
_cases = append(append([]interface{}{}, _cases.([]interface{})...), _current) } };
return func() interface{} { s := make(map[string]interface{}); s["cases"] = _cases;s["defaults"] = _defaults;s["hasDefault"] = _hasDefault;; return s }();
				return nil
			},
			FixedArgs: 0,
		};
_parseCases = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				;
mml.Nop();
return _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
				;
				mml.Nop(_c);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "select-case";s["expression"] = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_c, 0)));s["body"] = func() interface{} { s := make(map[string]interface{}); s["type"] = "statement-list";s["statements"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.RefRange(_c, 1, nil)));; return s }();; return s }()
			},
			FixedArgs: 1,
		}, _c));
				return nil
			},
			FixedArgs: 1,
		};
_lines = _groupLines.(*mml.Function).Call([]interface{}{});
return func() interface{} { s := make(map[string]interface{}); s["cases"] = _parseCases.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_lines, "cases")));s["defaultStatements"] = func() interface{} { s := make(map[string]interface{}); s["type"] = "statement-list";s["statements"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.Ref(_lines, "defaults")));; return s }();s["hasDefault"] = mml.Ref(_lines, "hasDefault");; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseExport = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _d interface{};
mml.Nop(_d);
_d = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)));
return func() interface{} { s := make(map[string]interface{}); s["type"] = "definition-list";s["definitions"] = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return func() interface{} { s := make(map[string]interface{}); sp := _d.(map[string]interface{});
for k, v := range sp { s[k] = v };s["exported"] = true;; return s }()
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, func () interface{} { c = mml.BinaryOp(11, mml.Ref(_d, "type"), "definition"); if c.(bool) { return append([]interface{}{}, _d) } else { return mml.Ref(_d, "definitions") } }()));; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseUseFact = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				var _capture interface{};
var _path interface{};
mml.Nop(_capture, _path);
_capture = "";
_path = "";
switch mml.Ref(mml.Ref(mml.Ref(_ast, "nodes"), 0), "name") {
case "use-inline":
;
mml.Nop();
_capture = ".";
_path = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)))
case "symbol":
;
mml.Nop();
_capture = mml.Ref(_parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0))), "name");
_path = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 1)))
default:
;
mml.Nop();
_path = _parse.(*mml.Function).Call(append([]interface{}{}, mml.Ref(mml.Ref(_ast, "nodes"), 0)))
};
return func() interface{} { s := make(map[string]interface{}); s["type"] = "use";s["capture"] = _capture;s["path"] = _path;; return s }();
				return nil
			},
			FixedArgs: 1,
		};
_parseUse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = "use-list";s["uses"] = _map.(*mml.Function).Call(append([]interface{}{}, _parse, mml.Ref(_ast, "nodes")));; return s }()
			},
			FixedArgs: 1,
		};
_parse = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ast = a[0];
				;
				mml.Nop(_ast);
				;
mml.Nop();
switch mml.Ref(_ast, "name") {
case "line-comment-content":
;
mml.Nop();
return func() interface{} { s := make(map[string]interface{}); s["type"] = "comment";; return s }()
case "int":
;
mml.Nop();
return _parseInt.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "text")))
case "float":
;
mml.Nop();
return _parseFloat.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_ast, "text")))
case "string":
;
mml.Nop();
return _parseString.(*mml.Function).Call(append([]interface{}{}, _ast))
case "true":
;
mml.Nop();
return true
case "false":
;
mml.Nop();
return false
case "symbol":
;
mml.Nop();
return _parseSymbol.(*mml.Function).Call(append([]interface{}{}, _ast))
case "spread-expression":
;
mml.Nop();
return _parseSpread.(*mml.Function).Call(append([]interface{}{}, _ast))
case "list":
;
mml.Nop();
return _parseList.(*mml.Function).Call(append([]interface{}{}, _ast))
case "mutable-list":
;
mml.Nop();
return _parseMutableList.(*mml.Function).Call(append([]interface{}{}, _ast))
case "expression-key":
;
mml.Nop();
return _parseExpressionKey.(*mml.Function).Call(append([]interface{}{}, _ast))
case "entry":
;
mml.Nop();
return _parseEntry.(*mml.Function).Call(append([]interface{}{}, _ast))
case "struct":
;
mml.Nop();
return _parseStruct.(*mml.Function).Call(append([]interface{}{}, _ast))
case "mutable-struct":
;
mml.Nop();
return _parseMutableStruct.(*mml.Function).Call(append([]interface{}{}, _ast))
case "return":
;
mml.Nop();
return _parseReturn.(*mml.Function).Call(append([]interface{}{}, _ast))
case "block":
;
mml.Nop();
return _parseStatementList.(*mml.Function).Call(append([]interface{}{}, _ast))
case "function":
;
mml.Nop();
return _parseFunction.(*mml.Function).Call(append([]interface{}{}, _ast))
case "effect":
;
mml.Nop();
return _parseEffect.(*mml.Function).Call(append([]interface{}{}, _ast))
case "range-from":
;
mml.Nop();
return _parseRange.(*mml.Function).Call(append([]interface{}{}, _ast))
case "range-to":
;
mml.Nop();
return _parseRange.(*mml.Function).Call(append([]interface{}{}, _ast))
case "symbol-index":
;
mml.Nop();
return _parseSymbolIndex.(*mml.Function).Call(append([]interface{}{}, _ast))
case "expression-index":
;
mml.Nop();
return _parseExpressionIndex.(*mml.Function).Call(append([]interface{}{}, _ast))
case "range-index":
;
mml.Nop();
return _parseRangeIndex.(*mml.Function).Call(append([]interface{}{}, _ast))
case "indexer":
;
mml.Nop();
return _parseIndexer.(*mml.Function).Call(append([]interface{}{}, _ast))
case "function-application":
;
mml.Nop();
return _parseFunctionApplication.(*mml.Function).Call(append([]interface{}{}, _ast))
case "unary-expression":
;
mml.Nop();
return _parseUnaryExpression.(*mml.Function).Call(append([]interface{}{}, _ast))
case "binary0":
;
mml.Nop();
return _parseBinaryExpression.(*mml.Function).Call(append([]interface{}{}, _ast))
case "binary1":
;
mml.Nop();
return _parseBinaryExpression.(*mml.Function).Call(append([]interface{}{}, _ast))
case "binary2":
;
mml.Nop();
return _parseBinaryExpression.(*mml.Function).Call(append([]interface{}{}, _ast))
case "binary3":
;
mml.Nop();
return _parseBinaryExpression.(*mml.Function).Call(append([]interface{}{}, _ast))
case "binary4":
;
mml.Nop();
return _parseBinaryExpression.(*mml.Function).Call(append([]interface{}{}, _ast))
case "chaining":
;
mml.Nop();
return _parseChaining.(*mml.Function).Call(append([]interface{}{}, _ast))
case "ternary-expression":
;
mml.Nop();
return _parserTernary.(*mml.Function).Call(append([]interface{}{}, _ast))
case "if":
;
mml.Nop();
return _parseIf.(*mml.Function).Call(append([]interface{}{}, _ast))
case "switch":
;
mml.Nop();
return _parseSwitch.(*mml.Function).Call(append([]interface{}{}, _ast))
case "range-over-expression":
;
mml.Nop();
return _parseRangeOver.(*mml.Function).Call(append([]interface{}{}, _ast))
case "loop":
;
mml.Nop();
return _parseLoop.(*mml.Function).Call(append([]interface{}{}, _ast))
case "value-capture":
;
mml.Nop();
return _parseValueCapture.(*mml.Function).Call(append([]interface{}{}, _ast))
case "mutable-capture":
;
mml.Nop();
return _parseMutableCapture.(*mml.Function).Call(append([]interface{}{}, _ast))
case "value-definition":
;
mml.Nop();
return _parseValueDefinition.(*mml.Function).Call(append([]interface{}{}, _ast))
case "value-definition-group":
;
mml.Nop();
return _parseDefinitions.(*mml.Function).Call(append([]interface{}{}, _ast))
case "mutable-definition-group":
;
mml.Nop();
return _parseMutableDefinitions.(*mml.Function).Call(append([]interface{}{}, _ast))
case "function-capture":
;
mml.Nop();
return _parseFunctionCapture.(*mml.Function).Call(append([]interface{}{}, _ast))
case "effect-capture":
;
mml.Nop();
return _parseEffectCapture.(*mml.Function).Call(append([]interface{}{}, _ast))
case "function-definition":
;
mml.Nop();
return _parseFunctionDefinition.(*mml.Function).Call(append([]interface{}{}, _ast))
case "function-definition-group":
;
mml.Nop();
return _parseDefinitions.(*mml.Function).Call(append([]interface{}{}, _ast))
case "effect-definition-group":
;
mml.Nop();
return _parseEffectDefinitions.(*mml.Function).Call(append([]interface{}{}, _ast))
case "assignment":
;
mml.Nop();
return _parseAssign.(*mml.Function).Call(append([]interface{}{}, _ast))
case "send":
;
mml.Nop();
return _parseSend.(*mml.Function).Call(append([]interface{}{}, _ast))
case "receive":
;
mml.Nop();
return _parseReceive.(*mml.Function).Call(append([]interface{}{}, _ast))
case "go":
;
mml.Nop();
return _parseGo.(*mml.Function).Call(append([]interface{}{}, _ast))
case "defer":
;
mml.Nop();
return _parseDefer.(*mml.Function).Call(append([]interface{}{}, _ast))
case "receive-definition":
;
mml.Nop();
return _parseReceiveDefinition.(*mml.Function).Call(append([]interface{}{}, _ast))
case "select":
;
mml.Nop();
return _parseSelect.(*mml.Function).Call(append([]interface{}{}, _ast))
case "export":
;
mml.Nop();
return _parseExport.(*mml.Function).Call(append([]interface{}{}, _ast))
case "use-fact":
;
mml.Nop();
return _parseUseFact.(*mml.Function).Call(append([]interface{}{}, _ast))
case "use":
;
mml.Nop();
return _parseUse.(*mml.Function).Call(append([]interface{}{}, _ast))
default:
;
mml.Nop();
return _parseStatementList.(*mml.Function).Call(append([]interface{}{}, _ast))
};
				return nil
			},
			FixedArgs: 1,
		};
_parseFile = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0];
				;
				mml.Nop(_path);
				var _in interface{};
var _ast interface{};
mml.Nop(_in, _ast);
_in = _open.(*mml.Function).Call(append([]interface{}{}, _path));
c = _isError.(*mml.Function).Call(append([]interface{}{}, _in)); if c.(bool) { ;
mml.Nop();
return _in };
defer _close.(*mml.Function).Call(append([]interface{}{}, _in));
_ast = _passErr.(*mml.Function).Call(append([]interface{}{}, _parseAST)).(*mml.Function).Call(append([]interface{}{}, _in.(*mml.Function).Call(append([]interface{}{}, mml.UnaryOp(2, 1)))));
c = _isError.(*mml.Function).Call(append([]interface{}{}, _ast)); if c.(bool) { ;
mml.Nop();
return _ast };
return _parse.(*mml.Function).Call(append([]interface{}{}, _ast));
				return nil
			},
			FixedArgs: 1,
		};
_findExportNames = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _statements = a[0];
				;
				mml.Nop(_statements);
				return _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return mml.Ref(_d, "symbol")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _d = a[0];
				;
				mml.Nop(_d);
				return mml.Ref(_d, "exported")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_code, "flattenedStatements").(*mml.Function).Call(append([]interface{}{}, "definition", "definition-list", "definitions")).(*mml.Function).Call(append([]interface{}{}, _statements))))))
			},
			FixedArgs: 1,
		};
_parseModules = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _path = a[0];
				;
				mml.Nop(_path);
				var _module interface{};
var _uses interface{};
var _usesModules interface{};
var _statements interface{};
var _currentCode interface{};
mml.Nop(_module, _uses, _usesModules, _statements, _currentCode);
_module = _parseFile.(*mml.Function).Call(append([]interface{}{}, _path));
_uses = mml.Ref(_code, "flattenedStatements").(*mml.Function).Call(append([]interface{}{}, "use", "use-list", "uses", mml.Ref(_module, "statements")));
c = _isError.(*mml.Function).Call(append([]interface{}{}, _module)); if c.(bool) { ;
mml.Nop();
return _module };
_usesModules = _passErr.(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
				;
				mml.Nop(_m);
				return func() interface{} { s := make(map[string]interface{}); s["type"] = mml.Ref(_m, "type");s["path"] = mml.Ref(_m, "path");s["statements"] = mml.Ref(_m, "statements");s["exportNames"] = _findExportNames.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_m, "statements")));; return s }()
			},
			FixedArgs: 1,
		})))).(*mml.Function).Call(append([]interface{}{}, _passErr.(*mml.Function).Call(append([]interface{}{}, _flat)).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_errors, "any").(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _parseModules)).(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				return mml.BinaryOp(9, mml.Ref(_u, "path"), ".mml")
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, _uses))))))))));
c = _isError.(*mml.Function).Call(append([]interface{}{}, _usesModules)); if c.(bool) { ;
mml.Nop();
return _usesModules };
_statements = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				;
mml.Nop();
c = (!_has.(*mml.Function).Call(append([]interface{}{}, "type", _s)).(bool) || (mml.BinaryOp(12, mml.Ref(_s, "type"), "use").(bool) && mml.BinaryOp(12, mml.Ref(_s, "type"), "use-list").(bool))); if c.(bool) { ;
mml.Nop();
return _s };
c = mml.BinaryOp(11, mml.Ref(_s, "type"), "use"); if c.(bool) { var _m interface{};
mml.Nop(_m);
_m = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
				;
				mml.Nop(_m);
				return mml.BinaryOp(11, mml.Ref(_m, "path"), mml.Ref(_s, "path"))
			},
			FixedArgs: 1,
		}, _usesModules));
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _m)), 0); if c.(bool) { ;
mml.Nop();
return _s };
return func() interface{} { s := make(map[string]interface{}); s["type"] = mml.Ref(_s, "type");s["path"] = mml.Ref(_s, "path");s["capture"] = mml.Ref(_s, "capture");s["exportNames"] = mml.Ref(mml.Ref(_m, 0), "exportNames");; return s }() };
return func() interface{} { s := make(map[string]interface{}); s["type"] = mml.Ref(_s, "type");s["uses"] = _map.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _u = a[0];
				;
				mml.Nop(_u);
				var _m interface{};
mml.Nop(_m);
_m = _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
				;
				mml.Nop(_m);
				return mml.BinaryOp(11, mml.Ref(_m, "path"), mml.BinaryOp(9, mml.Ref(_u, "path"), ".mml"))
			},
			FixedArgs: 1,
		}, _usesModules));
c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _m)), 0); if c.(bool) { ;
mml.Nop();
return _u };
return func() interface{} { s := make(map[string]interface{}); s["type"] = mml.Ref(_u, "type");s["path"] = mml.Ref(_u, "path");s["capture"] = mml.Ref(_u, "capture");s["exportNames"] = mml.Ref(mml.Ref(_m, 0), "exportNames");; return s }();
				return nil
			},
			FixedArgs: 1,
		}, mml.Ref(_s, "uses")));; return s }();
				return nil
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_module, "statements")));
_currentCode = func() interface{} { s := make(map[string]interface{}); s["type"] = mml.Ref(_module, "type");s["path"] = _path;s["statements"] = _statements;; return s }();
return append(append([]interface{}{}, _currentCode), _usesModules.([]interface{})...);
				return nil
			},
			FixedArgs: 1,
		}; exports["parseModules"] = _parseModules
		return exports
	})
modulePath = "lang.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
var _logger interface{};
var _list interface{};
var _strings interface{};
var _ints interface{};
var _errors interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _logger, _list, _strings, _ints, _errors);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_ints = mml.Modules.Use("ints.mml");
_logger = mml.Modules.Use("log.mml");
_errors = mml.Modules.Use("errors.mml");
_fold = mml.Ref(_list, "fold"); exports["fold"] = _fold;
_foldr = mml.Ref(_list, "foldr"); exports["foldr"] = _foldr;
_map = mml.Ref(_list, "map"); exports["map"] = _map;
_filter = mml.Ref(_list, "filter"); exports["filter"] = _filter;
_contains = mml.Ref(_list, "contains"); exports["contains"] = _contains;
_sort = mml.Ref(_list, "sort"); exports["sort"] = _sort;
_flat = mml.Ref(_list, "flat"); exports["flat"] = _flat;
_join = mml.Ref(_strings, "join"); exports["join"] = _join;
_joins = mml.Ref(_strings, "joins"); exports["joins"] = _joins;
_formats = mml.Ref(_strings, "formats"); exports["formats"] = _formats;
_enum = mml.Ref(_ints, "enum"); exports["enum"] = _enum;
_log = mml.Ref(_logger, "log"); exports["log"] = _log;
_onlyErr = mml.Ref(_errors, "only"); exports["onlyErr"] = _onlyErr;
_passErr = mml.Ref(_errors, "pass"); exports["passErr"] = _passErr
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "ints.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _counter interface{};
var _enum interface{};
mml.Nop(_counter, _enum);
_counter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _c interface{};
mml.Nop(_c);
_c = mml.UnaryOp(2, 1);
return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				;
mml.Nop();
_c = mml.BinaryOp(9, _c, 1);
return _c;
				return nil
			},
			FixedArgs: 0,
		};
				return nil
			},
			FixedArgs: 0,
		}; exports["counter"] = _counter;
_enum = _counter; exports["enum"] = _enum
		return exports
	})
modulePath = "log.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _log interface{};
var _list interface{};
var _strings interface{};
mml.Nop(_log, _list, _strings);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_log = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0:];
				;
				mml.Nop(_a);
				;
mml.Nop();
_stderr.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call(append([]interface{}{}, " ")).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_list, "map").(*mml.Function).Call(append([]interface{}{}, _string)).(*mml.Function).Call(append([]interface{}{}, _a))))));
_stderr.(*mml.Function).Call(append([]interface{}{}, "\n"));
return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 0); if c.(bool) { return "" } else { return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 1)) } }();
				return nil
			},
			FixedArgs: 0,
		}; exports["log"] = _log
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "code.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _controlStatement interface{};
var _breakControl interface{};
var _continueControl interface{};
var _unaryOp interface{};
var _binaryNot interface{};
var _plus interface{};
var _minus interface{};
var _logicalNot interface{};
var _binaryOp interface{};
var _binaryAnd interface{};
var _binaryOr interface{};
var _xor interface{};
var _andNot interface{};
var _lshift interface{};
var _rshift interface{};
var _mul interface{};
var _div interface{};
var _mod interface{};
var _add interface{};
var _sub interface{};
var _eq interface{};
var _notEq interface{};
var _less interface{};
var _lessOrEq interface{};
var _greater interface{};
var _greaterOrEq interface{};
var _logicalAnd interface{};
var _logicalOr interface{};
var _flattenedStatements interface{};
var _list interface{};
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
mml.Nop(_controlStatement, _breakControl, _continueControl, _unaryOp, _binaryNot, _plus, _minus, _logicalNot, _binaryOp, _binaryAnd, _binaryOr, _xor, _andNot, _lshift, _rshift, _mul, _div, _mod, _add, _sub, _eq, _notEq, _less, _lessOrEq, _greater, _greaterOrEq, _logicalAnd, _logicalOr, _flattenedStatements, _list, _fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr);
var __lang = mml.Modules.Use("lang.mml");;_fold = __lang["fold"];
_foldr = __lang["foldr"];
_map = __lang["map"];
_filter = __lang["filter"];
_contains = __lang["contains"];
_sort = __lang["sort"];
_flat = __lang["flat"];
_join = __lang["join"];
_joins = __lang["joins"];
_formats = __lang["formats"];
_enum = __lang["enum"];
_log = __lang["log"];
_onlyErr = __lang["onlyErr"];
_passErr = __lang["passErr"];
_list = mml.Modules.Use("list.mml");
_controlStatement = _enum.(*mml.Function).Call([]interface{}{}); exports["controlStatement"] = _controlStatement;
_breakControl = _controlStatement.(*mml.Function).Call([]interface{}{}); exports["breakControl"] = _breakControl;
_continueControl = _controlStatement.(*mml.Function).Call([]interface{}{}); exports["continueControl"] = _continueControl;
_unaryOp = _enum.(*mml.Function).Call([]interface{}{}); exports["unaryOp"] = _unaryOp;
_binaryNot = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["binaryNot"] = _binaryNot;
_plus = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["plus"] = _plus;
_minus = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["minus"] = _minus;
_logicalNot = _unaryOp.(*mml.Function).Call([]interface{}{}); exports["logicalNot"] = _logicalNot;
_binaryOp = _enum.(*mml.Function).Call([]interface{}{}); exports["binaryOp"] = _binaryOp;
_binaryAnd = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["binaryAnd"] = _binaryAnd;
_binaryOr = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["binaryOr"] = _binaryOr;
_xor = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["xor"] = _xor;
_andNot = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["andNot"] = _andNot;
_lshift = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["lshift"] = _lshift;
_rshift = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["rshift"] = _rshift;
_mul = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["mul"] = _mul;
_div = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["div"] = _div;
_mod = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["mod"] = _mod;
_add = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["add"] = _add;
_sub = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["sub"] = _sub;
_eq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["eq"] = _eq;
_notEq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["notEq"] = _notEq;
_less = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["less"] = _less;
_lessOrEq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["lessOrEq"] = _lessOrEq;
_greater = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["greater"] = _greater;
_greaterOrEq = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["greaterOrEq"] = _greaterOrEq;
_logicalAnd = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["logicalAnd"] = _logicalAnd;
_logicalOr = _binaryOp.(*mml.Function).Call([]interface{}{}); exports["logicalOr"] = _logicalOr;
_flattenedStatements = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _itemType = a[0];
var _listType = a[1];
var _listProp = a[2];
var _statements = a[3];
				;
				mml.Nop(_itemType, _listType, _listProp, _statements);
				var _type interface{};
var _toList interface{};
mml.Nop(_type, _toList);
_type = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return (_has.(*mml.Function).Call(append([]interface{}{}, "type", _s)).(bool) && _contains.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_s, "type"), append([]interface{}{}, _itemType, _listType))).(bool))
			},
			FixedArgs: 1,
		};
_toList = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				return func () interface{} { c = mml.BinaryOp(11, mml.Ref(_s, "type"), _itemType); if c.(bool) { return append([]interface{}{}, _s) } else { return mml.Ref(_s, _listProp) } }()
			},
			FixedArgs: 1,
		};
return mml.Ref(_list, "flat").(*mml.Function).Call(append([]interface{}{}, _map.(*mml.Function).Call(append([]interface{}{}, _toList)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _type)).(*mml.Function).Call(append([]interface{}{}, _statements))))));
				return nil
			},
			FixedArgs: 4,
		}; exports["flattenedStatements"] = _flattenedStatements
		return exports
	})
modulePath = "lang.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _sort interface{};
var _flat interface{};
var _join interface{};
var _joins interface{};
var _formats interface{};
var _enum interface{};
var _log interface{};
var _onlyErr interface{};
var _passErr interface{};
var _logger interface{};
var _list interface{};
var _strings interface{};
var _ints interface{};
var _errors interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _sort, _flat, _join, _joins, _formats, _enum, _log, _onlyErr, _passErr, _logger, _list, _strings, _ints, _errors);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_ints = mml.Modules.Use("ints.mml");
_logger = mml.Modules.Use("log.mml");
_errors = mml.Modules.Use("errors.mml");
_fold = mml.Ref(_list, "fold"); exports["fold"] = _fold;
_foldr = mml.Ref(_list, "foldr"); exports["foldr"] = _foldr;
_map = mml.Ref(_list, "map"); exports["map"] = _map;
_filter = mml.Ref(_list, "filter"); exports["filter"] = _filter;
_contains = mml.Ref(_list, "contains"); exports["contains"] = _contains;
_sort = mml.Ref(_list, "sort"); exports["sort"] = _sort;
_flat = mml.Ref(_list, "flat"); exports["flat"] = _flat;
_join = mml.Ref(_strings, "join"); exports["join"] = _join;
_joins = mml.Ref(_strings, "joins"); exports["joins"] = _joins;
_formats = mml.Ref(_strings, "formats"); exports["formats"] = _formats;
_enum = mml.Ref(_ints, "enum"); exports["enum"] = _enum;
_log = mml.Ref(_logger, "log"); exports["log"] = _log;
_onlyErr = mml.Ref(_errors, "only"); exports["onlyErr"] = _onlyErr;
_passErr = mml.Ref(_errors, "pass"); exports["passErr"] = _passErr
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "ints.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _counter interface{};
var _enum interface{};
mml.Nop(_counter, _enum);
_counter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				var _c interface{};
mml.Nop(_c);
_c = mml.UnaryOp(2, 1);
return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				;
				;
				mml.Nop();
				;
mml.Nop();
_c = mml.BinaryOp(9, _c, 1);
return _c;
				return nil
			},
			FixedArgs: 0,
		};
				return nil
			},
			FixedArgs: 0,
		}; exports["counter"] = _counter;
_enum = _counter; exports["enum"] = _enum
		return exports
	})
modulePath = "log.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _log interface{};
var _list interface{};
var _strings interface{};
mml.Nop(_log, _list, _strings);
_list = mml.Modules.Use("list.mml");
_strings = mml.Modules.Use("strings.mml");
_log = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0:];
				;
				mml.Nop(_a);
				;
mml.Nop();
_stderr.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_strings, "join").(*mml.Function).Call(append([]interface{}{}, " ")).(*mml.Function).Call(append([]interface{}{}, mml.Ref(_list, "map").(*mml.Function).Call(append([]interface{}{}, _string)).(*mml.Function).Call(append([]interface{}{}, _a))))));
_stderr.(*mml.Function).Call(append([]interface{}{}, "\n"));
return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 0); if c.(bool) { return "" } else { return mml.Ref(_a, mml.BinaryOp(10, _len.(*mml.Function).Call(append([]interface{}{}, _a)), 1)) } }();
				return nil
			},
			FixedArgs: 0,
		}; exports["log"] = _log
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})
modulePath = "strings.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _firstOr interface{};
var _join interface{};
var _joins interface{};
var _joinTwo interface{};
var _formats interface{};
var _formatOne interface{};
var _escape interface{};
var _unescape interface{};
mml.Nop(_firstOr, _join, _joins, _joinTwo, _formats, _formatOne, _escape, _unescape);
_firstOr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _v = a[0];
var _l = a[1];
				;
				mml.Nop(_v, _l);
				return func () interface{} { c = mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return mml.Ref(_l, 0) } else { return _v } }()
			},
			FixedArgs: 2,
		};
_join = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1];
				;
				mml.Nop(_j, _s);
				return func () interface{} { c = mml.BinaryOp(13, _len.(*mml.Function).Call(append([]interface{}{}, _s)), 2); if c.(bool) { return _firstOr.(*mml.Function).Call(append([]interface{}{}, "", _s)) } else { return mml.BinaryOp(9, mml.BinaryOp(9, mml.Ref(_s, 0), _j), _join.(*mml.Function).Call(append([]interface{}{}, _j, mml.RefRange(_s, 1, nil)))) } }()
			},
			FixedArgs: 2,
		}; exports["join"] = _join;
_joins = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _s = a[1:];
				;
				mml.Nop(_j, _s);
				return _join.(*mml.Function).Call(append([]interface{}{}, _j, _s))
			},
			FixedArgs: 1,
		}; exports["joins"] = _joins;
_joinTwo = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _j = a[0];
var _left = a[1];
var _right = a[2];
				;
				mml.Nop(_j, _left, _right);
				return _joins.(*mml.Function).Call(append([]interface{}{}, _j, _left, _right))
			},
			FixedArgs: 3,
		}; exports["joinTwo"] = _joinTwo;
_formats = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1:];
				;
				mml.Nop(_f, _a);
				return _format.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 1,
		}; exports["formats"] = _formats;
_formatOne = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _a = a[1];
				;
				mml.Nop(_f, _a);
				return _formats.(*mml.Function).Call(append([]interface{}{}, _f, _a))
			},
			FixedArgs: 2,
		}; exports["formatOne"] = _formatOne;
_escape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _first interface{};
mml.Nop(_first);
c = mml.BinaryOp(11, _s, ""); if c.(bool) { ;
mml.Nop();
return "" };
_first = mml.Ref(_s, 0);
switch _first {
case "\b":
;
mml.Nop();
_first = "\\b"
case "\f":
;
mml.Nop();
_first = "\\f"
case "\n":
;
mml.Nop();
_first = "\\n"
case "\r":
;
mml.Nop();
_first = "\\r"
case "\t":
;
mml.Nop();
_first = "\\t"
case "\v":
;
mml.Nop();
_first = "\\v"
case "\"":
;
mml.Nop();
_first = "\\\""
case "\\":
;
mml.Nop();
_first = "\\\\"
};
return mml.BinaryOp(9, _first, _escape.(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_s, 1, nil))));
				return nil
			},
			FixedArgs: 1,
		}; exports["escape"] = _escape;
_unescape = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _s = a[0];
				;
				mml.Nop(_s);
				var _esc interface{};
var _r interface{};
mml.Nop(_esc, _r);
_esc = false;
_r = []interface{}{};
for _i := 0; _i < _len.(*mml.Function).Call(append([]interface{}{}, _s)).(int); _i++ {
var _c interface{};
mml.Nop(_c);
_c = mml.Ref(_s, _i);
c = _esc; if c.(bool) { ;
mml.Nop();
switch _c {
case "b":
;
mml.Nop();
_c = "\b"
case "f":
;
mml.Nop();
_c = "\f"
case "n":
;
mml.Nop();
_c = "\n"
case "r":
;
mml.Nop();
_c = "\r"
case "t":
;
mml.Nop();
_c = "\t"
case "v":
;
mml.Nop();
_c = "\v"
};
_r = append(append([]interface{}{}, _r.([]interface{})...), _c);
_esc = false;
continue };
c = mml.BinaryOp(11, _c, "\\"); if c.(bool) { ;
mml.Nop();
_esc = true;
continue };
_r = append(append([]interface{}{}, _r.([]interface{})...), _c)
};
return _join.(*mml.Function).Call(append([]interface{}{}, "", _r));
				return nil
			},
			FixedArgs: 1,
		}; exports["unescape"] = _unescape
		return exports
	})
modulePath = "errors.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _ifErr interface{};
var _not interface{};
var _yes interface{};
var _pass interface{};
var _only interface{};
var _any interface{};
var _list interface{};
mml.Nop(_ifErr, _not, _yes, _pass, _only, _any, _list);
_list = mml.Modules.Use("list.mml");
_ifErr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _mod = a[0];
var _f = a[1];
				;
				mml.Nop(_mod, _f);
				return &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _a = a[0];
				;
				mml.Nop(_a);
				return func () interface{} { c = _mod.(*mml.Function).Call(append([]interface{}{}, _isError.(*mml.Function).Call(append([]interface{}{}, _a)))); if c.(bool) { return _f.(*mml.Function).Call(append([]interface{}{}, _a)) } else { return _a } }()
			},
			FixedArgs: 1,
		}
			},
			FixedArgs: 2,
		};
_not = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return !_x.(bool)
			},
			FixedArgs: 1,
		};
_yes = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _x = a[0];
				;
				mml.Nop(_x);
				return _x
			},
			FixedArgs: 1,
		};
_pass = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _not, _f))
			},
			FixedArgs: 1,
		}; exports["pass"] = _pass;
_only = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
				;
				mml.Nop(_f);
				return _ifErr.(*mml.Function).Call(append([]interface{}{}, _yes, _f))
			},
			FixedArgs: 1,
		}; exports["only"] = _only;
_any = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return mml.Ref(_list, "fold").(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _r)); if c.(bool) { return _r } else { return func () interface{} { c = _isError.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return _c } else { return append(append([]interface{}{}, _r.([]interface{})...), _c) } }() } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["any"] = _any
		return exports
	})
modulePath = "list.mml"
	mml.Modules.Set(modulePath, func() map[string]interface{} {
		exports := make(map[string]interface{})

		var c interface{}
		mml.Nop(c)
var _fold interface{};
var _foldr interface{};
var _map interface{};
var _filter interface{};
var _contains interface{};
var _flat interface{};
var _sort interface{};
mml.Nop(_fold, _foldr, _map, _filter, _contains, _flat, _sort);
_fold = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _fold.(*mml.Function).Call(append([]interface{}{}, _f, _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)), mml.RefRange(_l, 1, nil))) } }()
			},
			FixedArgs: 3,
		}; exports["fold"] = _fold;
_foldr = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _f = a[0];
var _i = a[1];
var _l = a[2];
				;
				mml.Nop(_f, _i, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return _i } else { return _f.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _foldr.(*mml.Function).Call(append([]interface{}{}, _f, _i, mml.RefRange(_l, 1, nil))))) } }()
			},
			FixedArgs: 3,
		}; exports["foldr"] = _foldr;
_map = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _m = a[0];
var _l = a[1];
				;
				mml.Nop(_m, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return append(append([]interface{}{}, _r.([]interface{})...), _m.(*mml.Function).Call(append([]interface{}{}, _c)))
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["map"] = _map;
_filter = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _p = a[0];
var _l = a[1];
				;
				mml.Nop(_p, _l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _r = a[1];
				;
				mml.Nop(_c, _r);
				return func () interface{} { c = _p.(*mml.Function).Call(append([]interface{}{}, _c)); if c.(bool) { return append(append([]interface{}{}, _r.([]interface{})...), _c) } else { return _r } }()
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 2,
		}; exports["filter"] = _filter;
_contains = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
var _l = a[1];
				;
				mml.Nop(_i, _l);
				return mml.BinaryOp(15, _len.(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _ii = a[0];
				;
				mml.Nop(_ii);
				return mml.BinaryOp(11, _ii, _i)
			},
			FixedArgs: 1,
		}, _l)))), 0)
			},
			FixedArgs: 2,
		}; exports["contains"] = _contains;
_flat = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _l = a[0];
				;
				mml.Nop(_l);
				return _fold.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _c = a[0];
var _result = a[1];
				;
				mml.Nop(_c, _result);
				return append(append([]interface{}{}, _result.([]interface{})...), _c.([]interface{})...)
			},
			FixedArgs: 2,
		}, []interface{}{}, _l))
			},
			FixedArgs: 1,
		}; exports["flat"] = _flat;
_sort = &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _less = a[0];
var _l = a[1];
				;
				mml.Nop(_less, _l);
				return func () interface{} { c = mml.BinaryOp(11, _len.(*mml.Function).Call(append([]interface{}{}, _l)), 0); if c.(bool) { return []interface{}{} } else { return append(append(append([]interface{}{}, _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, &mml.Function{
			F: func(a []interface{}) interface{} {
				var c interface{}
				mml.Nop(c)
				var _i = a[0];
				;
				mml.Nop(_i);
				return !_less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0), _i)).(bool)
			},
			FixedArgs: 1,
		})).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...), mml.Ref(_l, 0)), _sort.(*mml.Function).Call(append([]interface{}{}, _less)).(*mml.Function).Call(append([]interface{}{}, _filter.(*mml.Function).Call(append([]interface{}{}, _less.(*mml.Function).Call(append([]interface{}{}, mml.Ref(_l, 0))))).(*mml.Function).Call(append([]interface{}{}, mml.RefRange(_l, 1, nil))))).([]interface{})...) } }()
			},
			FixedArgs: 2,
		}; exports["sort"] = _sort
		return exports
	})

}

func main() {
	mml.Modules.Use("main.mml")
}
