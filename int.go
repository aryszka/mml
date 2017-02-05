package mml

type integer struct {
	sys int
}

func SysIntToInt(n int) *Val {
	return &Val{sys: &integer{sys: n}}
}
