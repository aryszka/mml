package next

import (
	"fmt"
	"os"
)

type TraceLevel int

const (
	TraceOff TraceLevel = iota
	TraceWarn
	TraceInfo
	TraceDebug
)

type Trace interface {
	Extend(string) Trace
	Warn(...interface{})
	Info(...interface{})
	Debug(...interface{})
}

type trace struct {
	level TraceLevel
	path  string
}

func (l TraceLevel) String() string {
	switch l {
	case TraceWarn:
		return "WARN"
	case TraceInfo:
		return "INFO"
	case TraceDebug:
		return "DEBUG"
	default:
		return ""
	}
}

func NewTrace(l TraceLevel) Trace {
	return &trace{
		level: l,
		path:  "/",
	}
}

func (t *trace) Extend(name string) Trace {
	if t.level == TraceOff {
		return t
	}

	var p string
	if t.path == "/" {
		p = t.path + name
	} else {
		p = t.path + "/" + name
	}

	return &trace{
		level: t.level,
		path:  p,
	}
}

func (t *trace) out(l TraceLevel, a ...interface{}) {
	if l > t.level {
		return
	}

	fmt.Fprintln(os.Stderr, append([]interface{}{l, t.path}, a...)...)
}

func (t *trace) Warn(a ...interface{})  { t.out(TraceWarn, a...) }
func (t *trace) Info(a ...interface{})  { t.out(TraceInfo, a...) }
func (t *trace) Debug(a ...interface{}) { t.out(TraceDebug, a...) }
