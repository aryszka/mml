package mml

type fn struct {
	sys func([]*Val) *Val
}

func NewCompiled(fixedArgs int, variadic bool, f func([]*Val) *Val) *Val {
	return &Val{sys: &fn{sys: f}}
}

func ApplySys(f *Val, args []*Val) *Val {
	return f.sys.(*fn).sys(args)
}
