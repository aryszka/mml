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

func makeChannel(e *env) function {
	return function{
		primitive: func(*env) (interface{}, error) { return make(chan interface{}), nil },
		env:       e,
	}
}

func createBufferedChannel(e *env) (interface{}, error) {
	size, err := eval(e, symbol{name: "size"})
	if err != nil {
		return nil, err
	}

	return make(chan interface{}, size.(int)), nil
}

func makeBufferedChannel(e *env) function {
	return function{
		primitive: createBufferedChannel,
		params:    []string{"size"},
		env:       e,
	}
}
