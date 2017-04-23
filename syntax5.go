package mml

func syntax5() (*syntax, error) {
	s := newSyntax()

	var err error
	withErr := func(f func() error) {
		if err != nil {
			return
		}

		err = f()
	}

	withErr(func() error { return s.primitive("int", intToken) })
	withErr(func() error { return s.primitive("string", stringToken) })
	withErr(func() error { return s.optional("optional-int", "int") })
	withErr(func() error { return s.optional("int-sequence-optional", "int-sequence") })
	withErr(func() error { return s.sequence("int-sequence", "int") })
	withErr(func() error { return s.sequence("optional-int-sequence", "optional-int") })
	withErr(func() error { return s.group("single-int", "int") })
	withErr(func() error { return s.group("single-optional-int", "optional-int") })
	withErr(func() error { return s.group("multiple-ints", "int", "int", "int") })
	withErr(func() error { return s.group("group-with-optional-item", "optional-int", "string") })
	withErr(func() error { return s.union("int-or-string", "int", "string") })
	withErr(func() error { return s.union("int-or-group-with-optional", "int", "group-with-optional-item") })

	return s, err
}
