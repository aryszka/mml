package mml

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type errPanic struct {
	value interface{}
}

func (ep errPanic) Error() string { return fmt.Sprint(ep.value) }

func callRecover(e *env) (interface{}, error) {
	if e.pendingErr == nil {
		return nil, nil
	}

	err := e.pendingErr
	e.pendingErr = nil
	return eval(e, functionApplication{function: symbol{name: "recovery"}, args: []interface{}{err}})
}

func callPanic(e *env) (interface{}, error) {
	a, err := eval(e, symbol{name: "err"})
	if err != nil {
		return nil, err
	}

	if err, ok := a.(error); ok {
		return nil, err
	}

	return nil, errPanic{value: a}
}

func recoverFunction(e *env) function {
	return function{
		primitive: callRecover,
		effect:    true,
		params:    []string{"recovery"},
		env:       e,
	}
}

func panicFunction(e *env) function {
	return function{
		primitive: callPanic,
		effect:    true,
		params:    []string{"err"},
		env:       e,
	}
}

func makeChannel(e *env) function {
	return function{
		primitive: func(*env) (interface{}, error) { return newChan(defaultScheduler, 0), nil },
		env:       e,
	}
}

func createBufferedChannel(e *env) (interface{}, error) {
	size, err := eval(e, symbol{name: "size"})
	if err != nil {
		return nil, err
	}

	return newChan(defaultScheduler, size.(int)), nil
}

func makeBufferedChannel(e *env) function {
	return function{
		primitive: createBufferedChannel,
		params:    []string{"size"},
		env:       e,
	}
}

func makeParse(e *env) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			return structure{}, nil
		},
		params: []string{"doc"},
		env:    e,
	}
}

func makeStdin(e *env) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			l, err := eval(e, symbol{name: "len"})
			if err != nil {
				return nil, err
			}

			ll := l.(int)
			if ll < 0 {
				b, err := ioutil.ReadAll(os.Stdin)
				return string(b), err
			}

			b := make([]byte, l.(int))
			_, err = os.Stdin.Read(b)
			return string(b), err
		},
		params: []string{"len"},
		env:    e,
	}
}

func makeIOWriter(e *env, w io.Writer) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			s, err := eval(e, symbol{name: "s"})
			if err != nil {
				return nil, err
			}

			_, err = w.Write([]byte(s.(string)))
			return nil, err
		},
		params: []string{"s"},
		env:    e,
	}
}

func makeStdout(e *env) function {
	return makeIOWriter(e, os.Stdout)
}

func makeStderr(e *env) function {
	return makeIOWriter(e, os.Stderr)
}

func parseForMML(e *env) (interface{}, error) {
	doc, err := eval(e, symbol{name: "doc"})
	if err != nil {
		return nil, err
	}

	code, err := parseModule(doc.(string))
	if err != nil {
		return nil, err
	}

	return codeMML(code), nil
}

func makeParseForMML(e *env) function {
	return function{
		primitive: parseForMML,
		params:    []string{"doc"},
		env:       e,
	}
}

func toString(e *env) (interface{}, error) {
	a, err := eval(e, symbol{name: "a"})
	if err != nil {
		return nil, err
	}

	return fmt.Sprint(a), nil
}

func makeString(e *env) function {
	return function{
		primitive: toString,
		params:    []string{"a"},
		env:       e,
	}
}

func makeFormat(e *env) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			f, err := eval(e, symbol{name: "f"})
			if err != nil {
				return nil, err
			}

			args, err := eval(e, symbol{name: "args"})
			if err != nil {
				return nil, err
			}

			return fmt.Sprintf(f.(string), args.(list).values...), nil
		},
		params: []string{"f", "args"},
		env:    e,
	}
}

func makeTypeCheck(e *env, check func(a interface{}) bool) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			a, err := eval(e, symbol{name: "a"})
			if err != nil {
				return nil, err
			}

			return check(a), nil
		},
		params: []string{"a"},
		env:    e,
	}
}

func makeIsInt(e *env) function {
	return makeTypeCheck(e, func(a interface{}) bool {
		_, ok := a.(int)
		return ok
	})
}

func makeIsFloat(e *env) function {
	return makeTypeCheck(e, func(a interface{}) bool {
		_, ok := a.(float64)
		return ok
	})
}

func makeIsString(e *env) function {
	return makeTypeCheck(e, func(a interface{}) bool {
		_, ok := a.(string)
		return ok
	})
}

func makeIsBool(e *env) function {
	return makeTypeCheck(e, func(a interface{}) bool {
		_, ok := a.(bool)
		return ok
	})
}

func makeIsError(e *env) function {
	return makeTypeCheck(e, func(a interface{}) bool {
		_, ok := a.(error)
		return ok
	})
}

func makeLen(e *env) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			a, err := eval(e, symbol{name: "a"})
			if err != nil {
				return nil, err
			}

			switch at := a.(type) {
			case list:
				return len(at.values), nil
			case string:
				return len(at), nil
			default:
				return nil, errUnsupportedCode
			}
		},
		params: []string{"a"},
		env:    e,
	}
}

func makeWithParams(e *env, p []string, f func(e *env, a []interface{}) (interface{}, error)) function {
	return function{
		primitive: func(e *env) (interface{}, error) {
			var a []interface{}
			for _, pi := range p {
				ai, err := eval(e, symbol{name: pi})
				if err != nil {
					return nil, err
				}

				a = append(a, ai)
			}

			return f(e, a)
		},
		params: p,
		env:    e,
	}
}

func makeError(e *env) function {
	return makeWithParams(e, []string{"message"}, func(e *env, a []interface{}) (interface{}, error) {
		return errors.New(a[0].(string)), nil
	})
}

func makeHas(e *env) function {
	return makeWithParams(e, []string{"k", "o"}, func(e *env, a []interface{}) (interface{}, error) {
		_, ok := a[1].(structure).values[a[0].(string)]
		return ok, nil
	})
}
