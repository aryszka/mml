package next

type definition interface {
	parser(*registry) (parser, error)
}

type parser interface {
	parse(Trace, *context, []string)
}

func parserNotFound(name string) error {
	return nil
}

func stringsContain(ss []string, s string) bool {
	return false
}
