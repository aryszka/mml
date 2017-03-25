package mml

import "strconv"

type str struct {
	sys string
}

func StringToSysBytes(s *Val) []byte {
	return []byte(s.sys.(*str).sys)
}

func String(args []*Val) *Val {
	switch args[0].sys.(type) {
	case *str:
		return args[0]
	case *integer:
		return IntToString(args)
	default:
		panic("not implemented conversion to string")
	}
}

func IntToString(args []*Val) *Val {
	return &Val{sys: &str{sys: strconv.Itoa(args[0].sys.(*integer).sys)}}
}

func SysStringToString(s string) *Val {
	return &Val{sys: &str{sys: s}}
}
