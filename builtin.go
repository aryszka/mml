package mml

import "fmt"

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
