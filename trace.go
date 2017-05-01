package mml

import "log"

type traceLevel int

const (
	traceOff traceLevel = iota
	traceWarn
	traceInfo
	traceDebug
)

type trace interface {
	extend(string) trace
	warn(...interface{})
	info(...interface{})
	debug(...interface{})
}

type printer interface {
	Println(...interface{})
}

type logPrinter struct{}

type noopTrace struct{}

type printTrace struct {
	level   traceLevel
	path    string
	printer printer
}

func (p logPrinter) Println(a ...interface{}) {
	log.Println(a...)
}

func (t noopTrace) extend(string) trace  { return t }
func (t noopTrace) warn(...interface{})  {}
func (t noopTrace) info(...interface{})  {}
func (t noopTrace) debug(...interface{}) {}

func newTrace(l traceLevel) *printTrace {
	return withPrintFunc(l, logPrinter{})
}

func withPrintFunc(l traceLevel, p printer) *printTrace {
	return &printTrace{
		level:   l,
		printer: p,
	}
}

func (t *printTrace) extend(name string) trace {
	if t.level == traceOff {
		return t
	}

	return &printTrace{
		path:    t.path + "/" + name,
		printer: t.printer,
	}
}

func (t *printTrace) outLevel(l traceLevel, a ...interface{}) {
	if t.level < l {
		return
	}

	t.printer.Println(a...)
}

func (t *printTrace) warn(a ...interface{})  { t.outLevel(traceWarn, a...) }
func (t *printTrace) info(a ...interface{})  { t.outLevel(traceInfo, a...) }
func (t *printTrace) debug(a ...interface{}) { t.outLevel(traceDebug, a...) }
