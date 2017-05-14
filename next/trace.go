package next

type trace interface {
	extend(string) trace
}
