package mml

func newMMLSyntax() (*syntax, error) {
	s := newSyntax()

	var err error
	withErr := func(f func() error) {
		if err != nil {
			return
		}

		err = f()
	}

	// withErr(func() error { return s.primitive("int", intToken) })
	// withErr(func() error { return s.union("expression", "int") })
	// withErr(func() error { return s.sequence("statement-sequence", "expression") })
	// withErr(func() error { return s.union("document", "statement-sequence") })

	withErr(func() error { return s.primitive("nl", nl) })
	withErr(func() error { return s.primitive("semicolon", semicolon) })
	withErr(func() error { return s.union("seq-sep", "nl", "semicolon") })
	withErr(func() error { return s.primitive("int", intToken) })
	withErr(func() error { return s.union("expression", "int") })
	withErr(func() error { return s.union("statement", "expression") })
	withErr(func() error { return s.union("statement-sequence-item", "statement", "seq-sep") })
	withErr(func() error { return s.sequence("statement-sequence", "statement-sequence-item") })
	withErr(func() error { return s.union("document", "statement-sequence") })

	return s, err
}
