package mml

type list struct {
	sys []*Val
}

func ListFromSysSlice(s []*Val) *Val {
	return &Val{sys: &list{sys: s}}
}

func ListToSysSlice(l *Val) []*Val {
	return l.sys.(*list).sys
}
