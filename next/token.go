package next

type Token struct {
	Offset int
	Line   int
	Column int
	Value  rune
}
