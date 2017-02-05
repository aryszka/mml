package mml

type symbol struct {
	sys string
}

func SysStringToSymbol(s string) *Val {
	return &Val{sys: &symbol{sys: s}}
}
