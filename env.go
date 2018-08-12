package mml

import "errors"

type env struct {
	parent *env
	values map[string]interface{}
}

var (
	errSymbolExists  = errors.New("symbol exists")
	errSymbolMissing = errors.New("symbol missing")
)

func newEnv() *env {
	return &env{values: make(map[string]interface{})}
}

func (e *env) clone() *env {
	return &env{parent: e, values: make(map[string]interface{})}
}

func (e *env) define(name string, value interface{}) error {
	if _, ok := e.values[name]; ok {
		return errSymbolExists
	}

	e.values[name] = value
	return nil
}

func (e *env) lookup(name string) (interface{}, error) {
	v, ok := e.values[name]
	if ok {
		return v, nil
	}

	if e.parent == nil {
		return nil, errSymbolMissing
	}

	return e.parent.lookup(name)
}
