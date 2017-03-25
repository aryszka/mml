package mml

var Builtin = map[string]*Val{
	"sum":    NewCompiled(0, true, Sum),
	"stdout": NewCompiled(1, false, Stdout),
	"string": NewCompiled(1, false, String),
}
