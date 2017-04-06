package mml

func init() {
	primitive("int", intToken)
	primitive("string", stringToken)
	optional("optional-int", "int")
	optional("int-sequence-optional", "int-sequence")
	sequence("int-sequence", "int")
	sequence("optional-int-sequence", "optional-int")
	group("single-int", "int")
	group("single-optional-int", "optional-int")
	group("multiple-ints", "int", "int", "int")
	group("optional-group-item", "optional-int", "string")
}
