package mml

type fn struct {
	sys func([]*Val) *Val
}

func ApplySys(f *Val, args ...*Val) *Val {
	return f.sys.(*fn).sys(args)
}
