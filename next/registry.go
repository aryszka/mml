package next

type registry struct {
	definitions map[string]definition
	parsers     map[string]parser
}

func newRegistry() *registry {
	return &registry{
		definitions: make(map[string]definition),
		parsers:     make(map[string]parser),
	}
}

func (r *registry) definition(name string) (definition, bool) {
	d, ok := r.definitions[name]
	return d, ok
}

func (r *registry) parser(name string) (parser, bool) {
	p, ok := r.parsers[name]
	return p, ok
}

func (r *registry) setDefinition(d definition) error {
	if _, ok := r.definitions[d.nodeName()]; ok {
		return ErrDuplicateDefinition
	}

	r.definitions[d.nodeName()] = d
	return nil
}

func (r *registry) setParser(p parser) {
	r.parsers[p.nodeName()] = p
}
