package mml

import "os"

var Sys *Val

func Stdout(args []*Val) *Val {
	os.Stdout.Write(StringToSysBytes(args[0]))
	return Void
}
