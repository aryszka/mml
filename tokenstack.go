package mml

type tokenStack struct {
	stack []*token
}

// TODO: consider representing the tokens as a linked list

func newTokenStack() *tokenStack {
	return withLength(0)
}

func withLength(l int) *tokenStack {
	ts := &tokenStack{}
	if l > 0 {
		ts.stack = make([]*token, l)
	}

	return ts
}

func mergeStack(to, from *tokenStack) *tokenStack {
	if from == nil {
		return to
	}

	if to == nil {
		to = newTokenStack()
	}

	to.merge(from)
	return to
}

func (s *tokenStack) push(t *token) {
	s.stack = append(s.stack, t)
}

func (s *tokenStack) merge(from *tokenStack) {
	need := len(s.stack) + len(from.stack) - cap(s.stack)
	if need > 0 {
		s.stack = s.stack[:cap(s.stack)]
		for need > 0 {
			s.stack = append(s.stack, nil)
			need--
		}
	} else {
		s.stack = s.stack[:len(s.stack)+len(from.stack)]
	}

	copy(s.stack[len(s.stack)-len(from.stack):], from.stack)
}

func (s *tokenStack) mergeTokens(t []*token) {
	for len(t) > 0 {
		s.push(t[len(t)-1])
		t = t[:len(t)-1]
	}
}

func (s *tokenStack) has() bool {
	return s.len() > 0
}

func (s *tokenStack) len() int {
	return len(s.stack)
}

func (s *tokenStack) peek() *token {
	return s.stack[len(s.stack)-1]
}

func (s *tokenStack) pop() *token {
	var t *token
	t, s.stack = s.stack[len(s.stack)-1], s.stack[:len(s.stack)-1]
	return t
}

func (s *tokenStack) popIfAny() (*token, bool) {
	if s == nil || !s.has() {
		return nil, false
	}

	return s.pop(), true
}

func (s *tokenStack) drop(n int) {
	nextLength := len(s.stack) - n
	if nextLength < 0 {
		nextLength = 0
	}

	s.stack = s.stack[:nextLength]
}

func (s *tokenStack) clear() {
	s.drop(len(s.stack))
}

func (s *tokenStack) findCachedNode(n *node) int {
	if s == nil {
		return 0
	}

	for tokenIndex, token := range n.tokens {
		if token != s.peek() {
			continue
		}

		var skip int
		if len(n.tokens)-tokenIndex > len(s.stack) {
			skip = len(n.tokens) - tokenIndex - len(s.stack)
			s.clear()
		} else {
			skip = 0
			s.drop(len(n.tokens) - tokenIndex)
		}

		return skip
	}

	return 0
}