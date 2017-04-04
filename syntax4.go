package mml

func init() {
	primitive("int", intToken)
	primitive("string", stringToken)
	optional("optional-int", "int")
	sequence("int-sequence", "int")
	sequence("optional-int-sequence", "optional-int")
}
