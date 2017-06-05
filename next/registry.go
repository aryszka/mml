package next

type registry struct{}

func (r *registry) definition(name string) (definition, bool) {
	return nil, false
}

func (r *registry) parser(name string) (parser, bool) {
	return nil, false
}

func (r *registry) setParser(parser) {}
