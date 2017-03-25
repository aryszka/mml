package mml

type integer struct {
	sys int
}

func SysIntToInt(n int) *Val {
	return &Val{sys: &integer{sys: n}}
}

func Sum(args []*Val) *Val {
	var sum int
	for _, a := range args {
		sum += a.sys.(*integer).sys
	}

	return &Val{sys: &integer{sys: sum}}
}
