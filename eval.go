package mml

import (
	"errors"
	"fmt"
)

var (
	errUnsupportedCode             = errors.New("unsupported code")
	errExpectedListSpread          = errors.New("expected list spread")
	errInvalidStructKey            = errors.New("invalid struct key")
	errExpectedStructSpread        = errors.New("expected struct spread")
	errInvalidStructEntry          = errors.New("invalid struct entry")
	errUnexpectedIndexerExpression = errors.New("unexpected indexer expression")
	errInvalidListIndex            = errors.New("invalid list index")
	errMissingStructKey            = errors.New("missing struct key")
	errTooManyArgs                 = errors.New("too many args")
	errNotEnoughArgs               = errors.New("not enough args")
	errNotAFunction                = errors.New("not a function")
	errInvalidArgument             = errors.New("invalid argument")
	errExpectedBoolean             = errors.New("expected boolean")
	errInvalidSwitchExpression     = errors.New("invalid switch expression")
	errInvalidCaseExpression       = errors.New("invalid case expression")
	errInvalidLoopExpression       = errors.New("invalid loop expression")
	errInvalidAssignTarget         = errors.New("invalid assign target")
	errExpectedChannel             = errors.New("expected channel")
	errExpectedFunctionApplication = errors.New("expected function application")
)

func evalError(err error) {
	panic(err)
}

func evalSymbol(e *env, s symbol) (interface{}, error) {
	// TODO: handle _
	return e.lookup(s.name)
}

func evalExpressionList(e *env, exps []interface{}) (v []interface{}, err error) {
	for _, ei := range exps {
		switch et := ei.(type) {
		case spread:
			var sv interface{}
			if sv, err = eval(e, et.value); err != nil {
				return
			}

			sl, ok := sv.(list)
			if !ok {
				err = errExpectedListSpread
				return
			}

			v = append(v, sl.values...)
		default:
			var vi interface{}
			if vi, err = eval(e, ei); err != nil {
				return
			}

			v = append(v, vi)
		}
	}

	return
}

func evalList(e *env, l list) (result list, err error) {
	var v []interface{}
	if v, err = evalExpressionList(e, l.values); err != nil {
		return
	}

	result.mutable = l.mutable
	result.values = v
	return
}

func evalStruct(e *env, s structure) (result structure, err error) {
	var (
		key, value interface{}
		skey       string
		ok         bool
		ss         interface{}
		sss        structure
	)

	result.mutable = s.mutable

	result.values = make(map[string]interface{})
	for _, ei := range s.entries {
		switch eit := ei.(type) {
		case spread:
			if ss, err = eval(e, eit.value); err != nil {
				return
			}

			if sss, ok = ss.(structure); !ok {
				err = errExpectedStructSpread
				return
			}

			for skey, value = range sss.values {
				result.values[skey] = value
			}
		case entry:
			switch keyt := eit.key.(type) {
			case string:
				skey = keyt
			case symbol:
				skey = keyt.name
			case expressionKey:
				if key, err = eval(e, keyt.value); err != nil {
					return
				}

				if skey, ok = key.(string); !ok {
					err = errInvalidStructKey
					return
				}
			default:
				err = errInvalidStructKey
				return
			}

			if value, err = eval(e, eit.value); err != nil {
				return
			}

			result.values[skey] = value
		default:
			err = errInvalidStructEntry
			return
		}
	}

	return
}

func evalFunction(e *env, f function) function {
	f.env = e.extend()
	return f
}

func evalListIndex(e *env, i interface{}) (ii int, err error) {
	var (
		iv interface{}
		ok bool
	)

	if iv, err = eval(e, i); err != nil {
		return
	}

	if ii, ok = iv.(int); !ok {
		err = errInvalidListIndex
		return
	}

	return
}

func evalListIndexer(e *env, l list, i interface{}) (v interface{}, err error) {
	switch it := i.(type) {
	case rangeExpression:
		var (
			from, to int
			s        []interface{}
		)

		if it.from != nil {
			if from, err = evalListIndex(e, it.from); err != nil {
				return
			}
		}

		if it.to != nil {
			if to, err = evalListIndex(e, it.to); err != nil {
				return
			}
		}

		switch {
		case it.from != nil && it.to != nil:
			s = l.values[from:to]
		case it.from != nil && it.to == nil:
			s = l.values[from:]
		case it.from == nil && it.to != nil:
			s = l.values[:to]
		default:
			s = l.values[:]
		}

		v = list{values: s}
	default:
		var ii int
		if ii, err = evalListIndex(e, i); err != nil {
			return
		}

		v = l.values[ii]
	}

	return
}

func evalStructIndexer(e *env, s structure, i interface{}) (interface{}, error) {
	if _, ok := i.(rangeExpression); ok {
		return nil, errInvalidStructKey
	}

	k, err := eval(e, i)
	if err != nil {
		return nil, err
	}

	ks, ok := k.(string)
	if !ok {
		return nil, errInvalidStructKey
	}

	v, ok := s.values[ks]
	if !ok {
		return nil, errMissingStructKey
	}

	return v, nil
}

func evalIndexer(e *env, i indexer) (interface{}, error) {
	exp, err := eval(e, i.expression)
	if err != nil {
		return nil, err
	}

	switch et := exp.(type) {
	case list:
		return evalListIndexer(e, et, i.index)
	case structure:
		return evalStructIndexer(e, et, i.index)
	default:
		return nil, errUnexpectedIndexerExpression
	}
}

func evalStatementList(e *env, s statementList) (interface{}, error) {
	for _, si := range s.statements {
		if r, ok := si.(ret); ok {
			siv, err := eval(e, r.value)
			return ret{value: siv}, err
		}

		siv, err := eval(e, si)
		if err != nil {
			return nil, err
		}

		if siv == breakStatement || siv == continueStatement {
			return siv, nil
		}

		if _, ok := siv.(ret); ok {
			return siv, nil
		}
	}

	return nil, nil
}

func evalExpressionOrStatementList(e *env, v interface{}) (interface{}, error) {
	switch s := v.(type) {
	case statementList:
		return evalStatementList(e, s)
	default:
		return eval(e, v)
	}
}

func evalFunctionApplicationArgs(e *env, fa functionApplication) (f function, a []interface{}, err error) {
	var (
		ok bool
		fe interface{}
	)

	fe, err = eval(e, fa.function)
	if err != nil {
		return
	}

	if f, ok = fe.(function); !ok {
		err = errNotAFunction
		return
	}

	if a, err = evalExpressionList(e, fa.args); err != nil {
		return
	}

	a = append(f.args, a...)
	if len(a) > len(f.params) && f.collectParam == "" {
		err = errTooManyArgs
		return
	}

	return
}

func toErr(err interface{}) error {
	if err == nil {
		return nil
	}

	if eerr, ok := err.(error); ok {
		return eerr
	} else {
		return fmt.Errorf("%v", err)
	}
}

func evalExecuteFunctionApplication(f function, a []interface{}) (interface{}, error) {
	v := func() interface{} {
		for i, p := range f.params {
			if err := f.env.define(p, a[i]); err != nil {
				f.env.pendingErr = err
				return nil
			}
		}

		if f.collectParam != "" {
			if err := f.env.define(f.collectParam, list{values: a[len(f.params):]}); err != nil {
				f.env.pendingErr = err
				return nil
			}
		}

		if f.primitive != nil {
			defer func() {
				if err := toErr(recover()); err != nil {
					f.env.pendingErr = err
				}
			}()

			v, err := f.primitive(f.env)
			if err != nil {
				f.env.pendingErr = err
			}

			return v
		}

		defer func() {
			if err := toErr(recover()); err != nil {
				f.env.pendingErr = err
			}

			for i := len(f.env.deferred) - 1; i >= 0; i-- {
				d := f.env.deferred[i]
				d.function.env.injectContext(f.env)
				if _, err := evalExecuteFunctionApplication(d.function, d.args); err != nil {
					f.env.pendingErr = err
				}

				d.function.env.releaseContext(f.env)
			}
		}()

		v, err := evalExpressionOrStatementList(f.env, f.statement)
		if err != nil {
			f.env.pendingErr = err
			return nil
		}

		if r, ok := v.(ret); ok {
			return r.value
		}

		return v
	}()

	err := f.env.pendingErr
	f.env.pendingErr = nil
	if pe, ok := v.(errPanic); ok {
		v = pe.value
	}

	return v, err
}

func evalFunctionApplication(e *env, fa functionApplication) (interface{}, error) {
	f, a, err := evalFunctionApplicationArgs(e, fa)
	if err != nil {
		return nil, err
	}

	if len(a) < len(f.params) {
		f.args = a
		return f, nil
	}

	f.env.injectContext(e)
	v, err := evalExecuteFunctionApplication(f, a)
	f.env.releaseContext(e)
	return v, err
}

func evalUnary(e *env, u unary) (interface{}, error) {
	a, err := eval(e, u.arg)
	if err != nil {
		return nil, err
	}

	switch u.op {
	case binaryNot:
		ai, ok := a.(int)
		if !ok {
			return nil, errInvalidArgument
		}

		return ^ai, nil
	case plus:
		switch a.(type) {
		case int, float64:
			return a, nil
		default:
			return nil, errInvalidArgument
		}
	case minus:
		switch at := a.(type) {
		case int:
			return -at, nil
		case float64:
			return -at, nil
		default:
			return nil, errInvalidArgument
		}
	case logicalNot:
		ab, ok := a.(bool)
		if !ok {
			return nil, errInvalidArgument
		}

		return !ab, nil
	default:
		return nil, errUnsupportedCode
	}
}

func evalIntBinaryChecked(op binaryOperator, left, right int) (interface{}, error) {
	switch op {
	case binaryAnd:
		return left & right, nil
	case binaryOr:
		return left | right, nil
	case xor:
		return left ^ right, nil
	case andNot:
		return left &^ right, nil
	case lshift:
		return left << uint(right), nil
	case rshift:
		return left >> uint(right), nil
	case mul:
		return left * right, nil
	case div:
		return left / right, nil
	case mod:
		return left % right, nil
	case add:
		return left + right, nil
	case sub:
		return left - right, nil
	case less:
		return left < right, nil
	case lessOrEq:
		return left <= right, nil
	case greater:
		return left > right, nil
	case greaterOrEq:
		return left >= right, nil
	default:
		return nil, errUnsupportedCode
	}
}

func evalIntBinary(e *env, b binary) (interface{}, error) {
	left, err := eval(e, b.left)
	if err != nil {
		return nil, err
	}

	li, ok := left.(int)
	if !ok {
		return nil, errInvalidArgument
	}

	right, err := eval(e, b.right)
	if err != nil {
		return nil, err
	}

	ri, ok := right.(int)
	if !ok {
		return nil, errInvalidArgument
	}

	return evalIntBinaryChecked(b.op, li, ri)
}

func evalFloatBinaryChecked(op binaryOperator, left, right float64) (interface{}, error) {
	switch op {
	case mul:
		return left * right, nil
	case div:
		return left / right, nil
	case add:
		return left + right, nil
	case sub:
		return left - right, nil
	case less:
		return left < right, nil
	case lessOrEq:
		return left <= right, nil
	case greater:
		return left > right, nil
	case greaterOrEq:
		return left <= right, nil
	default:
		return nil, errUnsupportedCode
	}
}

func evalIntFloatBinary(e *env, b binary) (interface{}, error) {
	left, err := eval(e, b.left)
	if err != nil {
		return 0, err
	}

	right, err := eval(e, b.right)
	if err != nil {
		return 0, err
	}

	switch lt := left.(type) {
	case int:
		rt, ok := right.(int)
		if !ok {
			return nil, errInvalidArgument
		}

		return evalIntBinaryChecked(b.op, lt, rt)
	case float64:
		rt, ok := right.(float64)
		if !ok {
			return nil, errInvalidArgument
		}

		return evalFloatBinaryChecked(b.op, lt, rt)
	default:
		return nil, errUnsupportedCode
	}
}

func evalStringBinaryChecked(op binaryOperator, left, right string) (interface{}, error) {
	switch op {
	case add:
		return left + right, nil
	default:
		return nil, errUnsupportedCode
	}
}

func evalIntFloatStringBinary(e *env, b binary) (interface{}, error) {
	left, err := eval(e, b.left)
	if err != nil {
		return 0, err
	}

	right, err := eval(e, b.right)
	if err != nil {
		return 0, err
	}

	switch lt := left.(type) {
	case int:
		rt, ok := right.(int)
		if !ok {
			return nil, errInvalidArgument
		}

		return evalIntBinaryChecked(b.op, lt, rt)
	case float64:
		rt, ok := right.(float64)
		if !ok {
			return nil, errInvalidArgument
		}

		return evalFloatBinaryChecked(b.op, lt, rt)
	case string:
		rt, ok := right.(string)
		if !ok {
			return nil, errInvalidArgument
		}

		return evalStringBinaryChecked(b.op, lt, rt)
	default:
		return nil, errUnsupportedCode
	}
}

func evalBoolBinary(e *env, b binary) (bool, error) {
	left, err := eval(e, b.left)
	if err != nil {
		return false, err
	}

	lb, ok := left.(bool)
	if !ok {
		return false, errInvalidArgument
	}

	switch b.op {
	case logicalAnd:
		if !lb {
			return false, nil
		}

		right, err := eval(e, b.right)
		if err != nil {
			return false, err
		}

		rb, ok := right.(bool)
		if !ok {
			return false, errInvalidArgument
		}

		return rb, nil
	case logicalOr:
		if lb {
			return true, nil
		}

		right, err := eval(e, b.right)
		if err != nil {
			return false, err
		}

		rb, ok := right.(bool)
		if !ok {
			return false, errInvalidArgument
		}

		return rb, nil
	default:
		return false, errUnsupportedCode
	}
}

func evalEqNotEq(e *env, b binary) (interface{}, error) {
	left, err := eval(e, b.left)
	if err != nil {
		return nil, err
	}

	right, err := eval(e, b.right)
	if err != nil {
		return nil, err
	}

	switch b.op {
	case eq:
		return left == right, nil
	default:
		return left != right, nil
	}
}

func evalBinary(e *env, b binary) (interface{}, error) {
	switch b.op {
	case binaryAnd, binaryOr, xor, andNot, lshift, rshift:
		return evalIntBinary(e, b)
	case mul, div, mod, sub:
		return evalIntFloatBinary(e, b)
	case eq, notEq:
		return evalEqNotEq(e, b)
	case add, less, lessOrEq, greater, greaterOrEq:
		return evalIntFloatStringBinary(e, b)
	case logicalAnd, logicalOr:
		return evalBoolBinary(e, b)
	default:
		return nil, errUnsupportedCode
	}
}

func evalCond(e *env, t cond) (interface{}, error) {
	e = e.extend()

	c, err := eval(e, t.condition)
	if err != nil {
		return nil, err
	}

	cb, ok := c.(bool)
	if !ok {
		return nil, errExpectedBoolean
	}

	if cb {
		return evalExpressionOrStatementList(e, t.consequent)
	}

	if t.alternative == nil {
		return nil, nil
	}

	return evalExpressionOrStatementList(e, t.alternative)
}

func evalSwitch(e *env, s switchStatement) (interface{}, error) {
	var (
		exp, cexp    interface{}
		err          error
		cexpCond, ok bool
	)

	e = e.extend()

	if s.expression != nil {
		if exp, err = eval(e, s.expression); err != nil {
			return nil, errInvalidSwitchExpression
		}
	}

	for _, c := range s.cases {
		if cexp, err = eval(e, c.expression); err != nil {
			return nil, err
		}

		if cexp == nil {
			return nil, errInvalidCaseExpression
		}

		if exp == nil {
			if cexpCond, ok = cexp.(bool); !ok {
				return nil, errInvalidCaseExpression
			}
		}

		if exp == nil && !cexpCond {
			continue
		}

		if exp != nil && exp != cexp {
			continue
		}

		return evalStatementList(e, c.body)
	}

	return evalStatementList(e, s.defaultStatements)
}

func evalLoopBody(e *env, s statementList) (interface{}, error, bool) {
	v, err := evalStatementList(e, s)
	if err != nil {
		return nil, err, true
	}

	if v == breakStatement {
		return nil, nil, true
	}

	if v == continueStatement {
		return nil, nil, false
	}

	if _, ok := v.(ret); ok {
		return v, nil, true
	}

	return nil, nil, false
}

func evalUnconditionalLoop(e *env, l loop) (interface{}, error) {
	e = e.extend()
	for {
		if v, err, r := evalLoopBody(e, l.body); err != nil || r {
			return v, err
		}
	}

	return nil, nil
}

func evalConditionalLoop(e *env, l loop) (interface{}, error) {
	ee := e.extend()
	for {
		expv, err := eval(e, l.expression)
		if err != nil {
			return nil, err
		}

		expvb, ok := expv.(bool)
		if !ok {
			return nil, errExpectedBoolean
		}

		if !expvb {
			break
		}

		if v, err, r := evalLoopBody(ee, l.body); err != nil || r {
			return v, err
		}
	}

	return nil, nil
}

func evalRangeArgument(e *env, a interface{}) (int, error) {
	v, err := eval(e, a)
	if err != nil {
		return 0, err
	}

	i, ok := v.(int)
	if !ok {
		return 0, errInvalidArgument
	}

	return i, nil
}

func evalNumericLoop(e *env, symbol string, r rangeExpression, s statementList) (interface{}, error) {
	var (
		v                 interface{}
		err               error
		from, to, counter int
		stop              bool
	)

	if r.from != nil {
		if from, err = evalRangeArgument(e, r.from); err != nil {
			return nil, err
		}
	}

	if r.to != nil {
		if to, err = evalRangeArgument(e, r.to); err != nil {
			return nil, err
		}
	}

	counter = from
	e = e.extend()
	if symbol != "" {
		if err = e.define(symbol, counter); err != nil {
			return nil, err
		}
	}

	for {
		if r.to != nil && counter >= to {
			return nil, nil
		}

		if v, err, stop = evalLoopBody(e, s); err != nil || stop {
			return v, err
		}

		counter++
		if symbol != "" {
			if err = e.set(symbol, counter); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func evalListLoop(e *env, symbol string, l list, s statementList) (interface{}, error) {
	l, err := evalList(e, l)
	if err != nil {
		return nil, err
	}

	e = e.extend()
	if symbol != "" {
		e.define(symbol, nil)
	}

	for _, i := range l.values {
		if symbol != "" {
			e.set(symbol, i)
		}

		if v, err, r := evalLoopBody(e, s); err != nil || r {
			return v, err
		}
	}

	return nil, nil
}

func evalStructLoop(e *env, symbol string, str structure, s statementList) (interface{}, error) {
	str, err := evalStruct(e, str)
	if err != nil {
		return nil, err
	}

	e = e.extend()
	if symbol != "" {
		e.define(symbol, nil)
	}

	for _, i := range str.values {
		if symbol != "" {
			e.set(symbol, i)
		}

		if v, err, r := evalLoopBody(e, s); err != nil || r {
			return v, err
		}
	}

	return nil, nil
}

func evalRangeLoop(e *env, r rangeOver, s statementList) (interface{}, error) {
	if re, ok := r.expression.(rangeExpression); ok || r.symbol != "" && r.expression == nil {
		return evalNumericLoop(e, r.symbol, re, s)
	}

	if l, ok := r.expression.(list); ok {
		return evalListLoop(e, r.symbol, l, s)
	}

	if str, ok := r.expression.(structure); ok {
		return evalStructLoop(e, r.symbol, str, s)
	}

	return nil, errInvalidLoopExpression
}

func evalLoop(e *env, l loop) (interface{}, error) {
	if l.expression == nil {
		return evalUnconditionalLoop(e, l)
	}

	if r, ok := l.expression.(rangeOver); ok {
		return evalRangeLoop(e, r, l.body)
	}

	return evalConditionalLoop(e, l)
}

func evalDefinition(e *env, d definition) (interface{}, error) {
	v, err := eval(e, d.expression)
	if err != nil {
		return nil, err
	}

	e.define(d.symbol, v)
	return nil, nil
}

func evalDefinitionList(e *env, d definitionList) (interface{}, error) {
	for _, di := range d.definitions {
		if _, err := eval(e, di); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func evalIndexerAssign(e *env, i indexer, v interface{}) (interface{}, error) {
	ee, err := eval(e, i.expression)
	if err != nil {
		return nil, err
	}

	switch eee := ee.(type) {
	case list:
		ii, err := evalListIndex(e, i.index)
		if err != nil {
			return nil, err
		}

		eee.values[ii] = v
		return nil, nil
	case structure:
		k, err := eval(e, i.index)
		if err != nil {
			return nil, err
		}

		ks, ok := k.(string)
		if !ok {
			return nil, errInvalidStructKey
		}

		eee.values[ks] = v
		return nil, nil
	default:
		return nil, errUnexpectedIndexerExpression
	}
}

func evalAssign(e *env, a assign) (interface{}, error) {
	v, err := eval(e, a.value)
	if err != nil {
		return nil, err
	}

	switch c := a.capture.(type) {
	case symbol:
		if err := e.set(c.name, v); err != nil {
			return nil, err
		}

		return nil, nil
	case indexer:
		return evalIndexerAssign(e, c, v)
	default:
		return nil, errInvalidAssignTarget
	}
}

func evalAssignList(e *env, a assignList) (interface{}, error) {
	for _, ai := range a.assignments {
		if _, err := evalAssign(e, ai); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func evalChan(e *env, v interface{}) (*channel, error) {
	c, err := eval(e, v)
	if err != nil {
		return nil, err
	}

	cc, ok := c.(*channel)
	if !ok {
		return nil, errExpectedChannel
	}

	return cc, nil
}

func evalSend(e *env, s send) (interface{}, error) {
	c, err := evalChan(e, s.channel)
	if err != nil {
		return nil, err
	}

	v, err := eval(e, s.value)
	if err != nil {
		return nil, err
	}

	c.send(v)
	return nil, nil
}

func evalReceive(e *env, r receive) (interface{}, error) {
	c, err := evalChan(e, r.channel)
	if err != nil {
		return nil, err
	}

	v, _ := c.receive()
	return v, nil
}

func evalGo(e *env, g goStatement) (interface{}, error) {
	f, a, err := evalFunctionApplicationArgs(e, g.application)
	if err != nil {
		return nil, err
	}

	if len(a) < len(f.params) {
		return nil, errNotEnoughArgs
	}

	go func() {
		if _, err := evalExecuteFunctionApplication(f, a); err != nil {
			evalError(err)
		}
	}()

	return nil, nil
}

func evalDefer(e *env, d deferStatement) (interface{}, error) {
	f, a, err := evalFunctionApplicationArgs(e, d.application)
	if err != nil {
		return nil, err
	}

	if len(a) < len(f.params) {
		return nil, errNotEnoughArgs
	}

	e.addDefer(deferred{function: f, args: a})
	return nil, nil
}

func evalSelect(e *env, s selectStatement) (interface{}, error) {
	e = e.extend()

	var items []*communicationItem
	cases := make(map[*communicationItem]selectCase)
	for _, c := range s.cases {
		var item *communicationItem
		switch ct := c.expression.(type) {
		case receive:
			ch, err := evalChan(e, ct.channel)
			if err != nil {
				return nil, err
			}

			item = receiveItem(ch)
		case send:
			ch, err := evalChan(e, ct.channel)
			if err != nil {
				return nil, err
			}

			v, err := eval(e, ct.value)
			if err != nil {
				return nil, err
			}

			item = sendItem(ch, v)
		case definition:
			ch, err := evalChan(e, ct.expression.(receive).channel)
			if err != nil {
				return nil, err
			}

			item = receiveItem(ch)
		}

		items = append(items, item)
		cases[item] = c
	}

	var d *communicationItem
	if s.hasDefault {
		d = defaultItem()
		items = append(items, d)
	}

	item := selectItem(defaultScheduler, items...)
	if item.err != nil {
		return nil, item.err
	}

	var sl statementList
	if item == d {
		sl = s.defaultStatements
	} else {
		sc := cases[item]
		sl = sc.body
		if d, ok := sc.expression.(definition); ok {
			if err := e.define(d.symbol, item.value); err != nil {
				return nil, err
			}
		}
	}

	return evalStatementList(e, sl)
}

func eval(e *env, code interface{}) (interface{}, error) {
	switch v := code.(type) {
	case int:
		return code, nil
	case float64:
		return code, nil
	case string:
		return code, nil
	case bool:
		return code, nil
	case error:
		return code, nil
	case symbol:
		return evalSymbol(e, v)
	case list:
		return evalList(e, v)
	case structure:
		return evalStruct(e, v)
	case chan interface{}:
		return code, nil
	case function:
		return evalFunction(e, v), nil
	case indexer:
		return evalIndexer(e, v)
	case functionApplication:
		return evalFunctionApplication(e, v)
	case unary:
		return evalUnary(e, v)
	case binary:
		return evalBinary(e, v)
	case cond:
		return evalCond(e, v)
	case switchStatement:
		return evalSwitch(e, v)
	case controlStatement:
		return code, nil
	case loop:
		return evalLoop(e, v)
	case definition:
		return evalDefinition(e, v)
	case definitionList:
		return evalDefinitionList(e, v)
	case assignList:
		return evalAssignList(e, v)
	case send:
		return evalSend(e, v)
	case receive:
		return evalReceive(e, v)
	case goStatement:
		return evalGo(e, v)
	case deferStatement:
		return evalDefer(e, v)
	case selectStatement:
		return evalSelect(e, v)
	default:
		return nil, errUnsupportedCode
	}
}

func evalStatement(e *env, s interface{}) interface{} {
	v, err := eval(e, s)
	if err != nil {
		evalError(err)
	}

	return v
}

func evalModule(e *env, m module) {
}
