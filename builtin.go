package mml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aryszka/mml/parser"
)

type unaryOperator int

const (
	binaryNot unaryOperator = iota
	plus
	minus
	logicalNot
)

type unary struct {
	op  unaryOperator
	arg interface{}
}

type binaryOperator int

const (
	binaryAnd binaryOperator = iota
	binaryOr
	xor
	andNot
	lshift
	rshift
	mul
	div
	mod
	add
	sub
	eq
	notEq
	less
	lessOrEq
	greater
	greaterOrEq
	logicalAnd
	logicalOr
)

type List struct {
	Values []interface{}
}

type Struct struct {
	Values map[string]interface{}
}

type Function struct {
	F         func([]interface{}) interface{}
	FixedArgs int
	args      []interface{}
}

type ModuleContext struct {
	lock         sync.Mutex
	moduleLocks  map[string]*sync.Mutex
	initializers map[string]func() map[string]interface{}
	cache        map[string]map[string]interface{}
}

var Modules = &ModuleContext{
	moduleLocks:  make(map[string]*sync.Mutex),
	initializers: make(map[string]func() map[string]interface{}),
	cache:        make(map[string]map[string]interface{}),
}

func (f *Function) Bind(a []interface{}) *Function {
	b := *f
	b.args = a
	return &b
}

func (f *Function) Call(a []interface{}) interface{} {
	a = append(f.args, a...)
	if len(a) < f.FixedArgs {
		return f.Bind(a)
	}

	return f.F(a)
}

func (c *ModuleContext) Set(path string, i func() map[string]interface{}) {
	c.initializers[path] = i
}

func (c *ModuleContext) Use(path string) *Struct {
	c.lock.Lock()
	m, ok := c.cache[path]
	if ok {
		c.lock.Unlock()
		return &Struct{m}
	}

	init := c.initializers[path]
	ml, ok := c.moduleLocks[path]
	if !ok {
		ml = &sync.Mutex{}
		c.moduleLocks[path] = ml
	}

	ml.Lock()
	c.lock.Unlock()

	m = init()

	c.lock.Lock()
	c.cache[path] = m
	c.lock.Unlock()
	ml.Unlock()

	return &Struct{m}
}

func Ref(v, k interface{}) interface{} {
	switch vt := v.(type) {
	case string:
		return string(vt[k.(int)])
	case *List:
		return vt.Values[k.(int)]
	case *Struct:
		ret := vt.Values[k.(string)]
		if ret == nil {
			panic("ref: undefined key: " + k.(string))
		}

		return ret
	default:
		panic(fmt.Sprintf("ref: unsupported code: %v: %v", k, v))
	}
}

func RefRange(v, from, to interface{}) interface{} {
	switch vt := v.(type) {
	case string:
		switch {
		case from == nil && to == nil:
			return vt[:]
		case from == nil:
			return vt[:to.(int)]
		case to == nil:
			return vt[from.(int):]
		default:
			return vt[from.(int):to.(int)]
		}
	case *List:
		switch {
		case from == nil && to == nil:
			return &List{vt.Values[:]}
		case from == nil:
			return &List{vt.Values[:to.(int)]}
		case to == nil:
			return &List{vt.Values[from.(int):]}
		default:
			return &List{vt.Values[from.(int):to.(int)]}
		}
	default:
		panic("ref range: unsupported code")
	}
}

func SetRef(e, k, v interface{}) interface{} {
	switch et := e.(type) {
	case *List:
		et.Values[k.(int)] = v
	case *Struct:
		et.Values[k.(string)] = v
	default:
		panic("set-ref: unsupported code")
	}

	return nil
}

func UnaryOp(op int, arg interface{}) interface{} {
	switch unaryOperator(op) {
	case binaryNot:
		switch at := arg.(type) {
		case int:
			return +at
		default:
			panic("unary: unsupported code")
		}
	case plus:
		switch at := arg.(type) {
		case int:
			return +at
		case float64:
			return +at
		default:
			panic("unary: unsupported code")
		}
	case minus:
		switch at := arg.(type) {
		case int:
			return -at
		case float64:
			return -at
		default:
			panic("unary: unsupported code")
		}
	default:
		panic("unary: unsupported code")
	}
}

func BinaryOp(op int, left, right interface{}) interface{} {
	switch binaryOperator(op) {
	case binaryAnd:
		switch lt := left.(type) {
		case int:
			return lt & right.(int)
		default:
			panic("binary: unsupported code")
		}
	case binaryOr:
		switch lt := left.(type) {
		case int:
			return lt | right.(int)
		default:
			panic("binary: unsupported code")
		}
	case xor:
		switch lt := left.(type) {
		case int:
			return lt ^ right.(int)
		default:
			panic("binary: unsupported code")
		}
	case andNot:
		switch lt := left.(type) {
		case int:
			return lt &^ right.(int)
		default:
			panic("binary: unsupported code")
		}
	case lshift:
		switch lt := left.(type) {
		case int:
			return lt << right.(uint)
		default:
			panic("binary: unsupported code")
		}
	case rshift:
		switch lt := left.(type) {
		case int:
			return lt >> right.(uint)
		default:
			panic("binary: unsupported code")
		}
	case mul:
		switch lt := left.(type) {
		case int:
			return lt * right.(int)
		case float64:
			return lt * right.(float64)
		default:
			panic("binary: unsupported code")
		}
	case div:
		switch lt := left.(type) {
		case int:
			return lt / right.(int)
		case float64:
			return lt / right.(float64)
		default:
			panic("binary: unsupported code")
		}
	case mod:
		switch lt := left.(type) {
		case int:
			return lt % right.(int)
		default:
			panic("binary: unsupported code")
		}
	case add:
		switch lt := left.(type) {
		case int:
			return lt + right.(int)
		case float64:
			return lt + right.(float64)
		case string:
			return lt + right.(string)
		default:
			panic("binary: add: unsupported code")
		}
	case sub:
		switch lt := left.(type) {
		case int:
			return lt - right.(int)
		case float64:
			return lt - right.(float64)
		default:
			panic("binary: sub: unsupported code")
		}
	case eq:
		return left == right
	case notEq:
		return left != right
	case less:
		switch lt := left.(type) {
		case int:
			return lt < right.(int)
		case float64:
			return lt < right.(float64)
		case string:
			return lt < right.(string)
		default:
			panic("binary: less: unsupported code")
		}
	case lessOrEq:
		switch lt := left.(type) {
		case int:
			return lt <= right.(int)
		case float64:
			return lt <= right.(float64)
		case string:
			return lt <= right.(string)
		default:
			panic("binary: less-or-eq: unsupported code")
		}
	case greater:
		switch lt := left.(type) {
		case int:
			return lt > right.(int)
		case float64:
			return lt > right.(float64)
		case string:
			return lt > right.(string)
		default:
			panic("binary: greater: unsupported code")
		}
	case greaterOrEq:
		switch lt := left.(type) {
		case int:
			return lt >= right.(int)
		case float64:
			return lt >= right.(float64)
		case string:
			return lt >= right.(string)
		default:
			panic("binary: greater-or-eq: unsupported code")
		}
	default:
		panic("binary: unsupported code")
	}
}

func Nop(...interface{}) {}

var IsError = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(error)
		return ok
	},
	FixedArgs: 1,
}

var IsBool = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(bool)
		return ok
	},
	FixedArgs: 1,
}

var IsInt = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(int)
		return ok
	},
	FixedArgs: 1,
}

var IsFloat = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(float64)
		return ok
	},
	FixedArgs: 1,
}

var IsString = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(string)
		return ok
	},
	FixedArgs: 1,
}

var IsList = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(*List)
		return ok
	},
	FixedArgs: 1,
}

var IsStruct = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(*Struct)
		return ok
	},
	FixedArgs: 1,
}

var IsFunction = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(*Function)
		return ok
	},
	FixedArgs: 1,
}

var IsChannel = &Function{
	F: func(a []interface{}) interface{} {
		_, ok := a[0].(chan interface{})
		return ok
	},
	FixedArgs: 1,
}

var Len = &Function{
	F: func(a []interface{}) interface{} {
		switch at := a[0].(type) {
		case *List:
			return len(at.Values)
		case *Struct:
			return len(at.Values)
		case string:
			return len(at)
		default:
			panic(fmt.Sprintf("len: unsupported code: %v", a[0]))
		}
	},
	FixedArgs: 1,
}

var Keys = &Function{
	F: func(a []interface{}) interface{} {
		s, ok := a[0].(*Struct)
		if !ok {
			panic("keys: unsupported code" + fmt.Sprint(a[0]))
		}

		var keys []interface{}
		for k := range s.Values {
			keys = append(keys, k)
		}

		return &List{Values: keys}
	},
	FixedArgs: 1,
}

var Format = &Function{
	F: func(a []interface{}) interface{} {
		f, ok := a[0].(string)
		if !ok {
			panic("format: unsupported code: " + fmt.Sprint(a[0]))
		}

		args, ok := a[1].(*List)
		if !ok {
			panic("format: unsupported code: " + fmt.Sprint(a[1]))
		}

		return fmt.Sprintf(f, args.Values...)
	},
	FixedArgs: 2,
}

var Stderr = &Function{
	F: func(a []interface{}) interface{} {
		s, ok := a[0].(string)
		if !ok {
			panic("stderr: unsupported code")
		}

		_, err := os.Stderr.Write([]byte(s))
		return err
	},
	FixedArgs: 1,
}

var Stdout = &Function{
	F: func(a []interface{}) interface{} {
		s, ok := a[0].(string)
		if !ok {
			panic("stdout: unsupported code")
		}

		_, err := os.Stdout.Write([]byte(s))
		return err
	},
	FixedArgs: 1,
}

var Stdin = &Function{
	F: func(a []interface{}) interface{} {
		if a[0].(int) < 0 {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}

			return string(b)
		}

		b := make([]byte, a[0].(int))
		n, err := os.Stdin.Read(b)
		if err != nil {
			return err
		}

		return string(b[:n])
	},
	FixedArgs: 1,
}

var String = &Function{
	F: func(a []interface{}) interface{} {
		return fmt.Sprint(a[0])
	},
	FixedArgs: 1,
}

func parseInt(a []interface{}) interface{} {
	s := a[0].(string)

	var base int
	switch {
	case strings.HasPrefix(s, "0x"):
		base = 16
		s = s[2:]
	case strings.HasPrefix(s, "0"):
		if s == "0" {
			return 0
		}

		base = 8
		s = s[1:]
	default:
		base = 10
	}

	i, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return err
	}

	return int(i)
}

var ParseInt = &Function{
	F:         parseInt,
	FixedArgs: 1,
}

func parseFloat(a []interface{}) interface{} {
	v, err := strconv.ParseFloat(a[0].(string), 64)
	if err != nil {
		return err
	}

	return v
}

var ParseFloat = &Function{
	F:         parseFloat,
	FixedArgs: 1,
}

func convertAST(goAST *parser.Node) *Struct {
	ast := make(map[string]interface{})
	ast["name"] = goAST.Name
	ast["text"] = goAST.Text()

	// TODO:
	ast["file"] = ""
	ast["from"] = 0
	ast["to"] = 0
	ast["line"] = 0
	ast["column"] = 0

	var nodes []interface{}
	for i := range goAST.Nodes {
		nodes = append(nodes, convertAST(goAST.Nodes[i]))
	}

	ast["nodes"] = &List{nodes}
	return &Struct{ast}
}

func parseAST(doc string) (ast *Struct, err error) {
	var goAST *parser.Node
	goAST, err = parser.Parse(bytes.NewBufferString(doc))
	if err != nil {
		return
	}

	return convertAST(goAST), nil
}

var ParseAST = &Function{
	F: func(a []interface{}) interface{} {
		ast, err := parseAST(a[0].(string))
		if err != nil {
			return err
		}

		return ast
	},
	FixedArgs: 1,
}

var Int = &Function{
	F: func(a []interface{}) interface{} {
		switch v := a[0].(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			return parseInt(a)
		default:
			return errors.New("unsupported argument")
		}
	},
	FixedArgs: 1,
}

var Float = &Function{
	F: func(a []interface{}) interface{} {
		switch v := a[0].(type) {
		case int:
			return float64(v)
		case float64:
			return v
		case string:
			return parseFloat(a)
		default:
			return errors.New("unsupported argument")
		}
	},
	FixedArgs: 1,
}

var Bool = &Function{
	F: func(a []interface{}) interface{} {
		switch v := a[0].(type) {
		case int:
			return v != 0
		case string:
			switch v {
			case "true":
				return true
			case "false":
				return false
			default:
				return errors.New("unsupported argument")
			}
		default:
			return errors.New("unsupported argument")
		}
	},
	FixedArgs: 1,
}

var Has = &Function{
	F: func(a []interface{}) interface{} {
		s, ok := a[1].(*Struct)
		if !ok {
			return false
		}

		_, ok = s.Values[a[0].(string)]
		return ok
	},
	FixedArgs: 2,
}

var Error = &Function{
	F: func(a []interface{}) interface{} {
		return errors.New(a[0].(string))
	},
	FixedArgs: 1,
}

var Panic = &Function{
	F: func(a []interface{}) interface{} {
		err, ok := a[0].(error)
		if !ok {
			err = fmt.Errorf("%v", a[0])
		}

		panic(err)
	},
	FixedArgs: 1,
}

var Exit = &Function{
	F: func(a []interface{}) interface{} {
		os.Exit(a[0].(int))
		return a[0].(int)
	},
	FixedArgs: 1,
}

var Open = &Function{
	F: func(a []interface{}) interface{} {
		f, err := os.Open(a[0].(string))
		if err != nil {
			return err
		}

		return &Function{
			F: func(a []interface{}) interface{} {
				l, ok := a[0].(int)
				if !ok {
					f.Close()
					return nil
				}

				if l < 0 {
					b, err := ioutil.ReadAll(f)
					if err != nil {
						return err
					}

					return string(b)
				}

				b := make([]byte, l)
				n, err := f.Read(b)
				if err != nil && err != io.EOF {
					return err
				}

				if err == io.EOF {
					f.Close()
				}

				return string(b[:n])
			},
			FixedArgs: 1,
		}
	},
	FixedArgs: 1,
}

var (
	Close *Function
	Args  interface{}
)

func init() {
	Close = &Function{
		F: func(a []interface{}) interface{} {
			return a[0].(*Function).F([]interface{}{Close})
		},
		FixedArgs: 1,
	}

	var args []interface{}
	for i := range os.Args {
		args = append(args, os.Args[i])
	}

	Args = &List{args}
}
