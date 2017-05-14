package next

// this should not be needed

type stack struct{
	tokens []*Token
}

func newStackToken(t *Token) *stack {
	return &stack{tokens: []*Token{t}}
}

func (s *stack) has() bool {
	return len(s.tokens) > 0
}

func (s *stack) pop() *Token {
	var t *Token
	t, s.tokens = s.tokens[len(s.tokens) - 1], s.tokens[:len(s.tokens) - 1]
	return t
}

func (s *stack) findCachedNode(n *Node) int {
	panic(ErrNotImplemented)
}

func (s *stack) push(*Token) {
	panic(ErrNotImplemented)
}

func (s *stack) merge(*stack) {
	panic(ErrNotImplemented)
}

func (s *stack) mergeTokens([]*Token) {
}

func (s *stack) clear() {
	panic(ErrNotImplemented)
}
