package next

type CommitType int

const (
	None  CommitType = 0
	Alias CommitType = 1 << iota
	Documentation
	Root
)
