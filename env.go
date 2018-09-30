package mml

import "errors"

type env struct {
	parent     *env
	values     map[string]interface{}
	deferred   []deferred
	pendingErr error
}

var (
	errSymbolExists   = errors.New("symbol exists")
	errSymbolNotFound = errors.New("symbol not found")
)

func newEnv() *env {
	v := make(map[string]interface{})
	e := &env{values: v}
	v["recover"] = recoverFunction(e)
	v["panic"] = panicFunction(e)
	v["chan"] = makeChannel(e)
	v["bufchan"] = makeBufferedChannel(e)
	v["yes"] = makeYes(e)
	v["not"] = makeNot(e)
	v["stdin"] = makeStdin(e)
	v["stdout"] = makeStdout(e)
	v["stderr"] = makeStderr(e)
	v["parse"] = makeParse(e)
	v["isError"] = makeIsError(e)
	v["parse"] = makeParseForMML(e)
	v["string"] = makeString(e)
	v["format"] = makeFormat(e)
	v["isInt"] = makeIsInt(e)
	v["isFloat"] = makeIsFloat(e)
	v["isString"] = makeIsString(e)
	v["isBool"] = makeIsBool(e)
	v["len"] = makeLen(e)
	return e
}

// TODO: methods are not thread safe

func (e *env) extend() *env {
	return &env{parent: e, values: make(map[string]interface{})}
}

func (e *env) define(name string, value interface{}) error {
	if _, ok := e.values[name]; ok {
		return errSymbolExists
	}

	e.values[name] = value
	return nil
}

func (e *env) lookupWithParent(name string) (interface{}, *env, error) {
	v, ok := e.values[name]
	if ok {
		return v, e, nil
	}

	if e.parent == nil {
		println(name)
		return nil, nil, errSymbolNotFound
	}

	return e.parent.lookupWithParent(name)
}

func (e *env) lookup(name string) (interface{}, error) {
	v, _, err := e.lookupWithParent(name)
	return v, err
}

func (e *env) set(name string, value interface{}) error {
	_, e, err := e.lookupWithParent(name)
	if err != nil {
		return err
	}

	e.values[name] = value
	return nil
}

func (e *env) addDefer(d deferred) {
	e.deferred = append(e.deferred, d)
}

func (e *env) injectContext(ctx *env) {
	if ctx.pendingErr != nil {
		e.pendingErr = ctx.pendingErr
	}
}

func (e *env) releaseContext(ctx *env) {
	ctx.pendingErr = e.pendingErr
}
