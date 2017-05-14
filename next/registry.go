package next

import (
	"fmt"
	"strings"
)

type registry struct {
	definitions map[string]definition
	root        definition
	generators  map[string]generator
}

func duplicateDefinition(name string) error {
	return fmt.Errorf("duplicate definition in syntax: %s", name)
}

func unspecifiedParser(name string) error {
	return fmt.Errorf("unspecified parser: %s", name)
}

func (r *registry) register(d definition) error {
	n := d.nodeName()
	if _, exists := r.definitions[n]; exists {
		return duplicateDefinition(n)
	}

	r.definitions[n] = d
	r.root = d
	return nil
}

func (r *registry) findDefinition(name string) (definition, error) {
	if d, ok := r.definitions[name]; ok {
		return d, nil
	} else {
		return nil, unspecifiedParser(name)
	}
}

func (r *registry) findDefinitions(names []string) ([]definition, error) {
	defs := make([]definition, len(names))
	for i, name := range names {
		if di, err := r.findDefinition(name); err != nil {
			return nil, err
		} else {
			defs[i] = di
		}
	}

	return defs, nil
}

func generatorKey(name, init string, excluded []string) string {
	// TODO: not reliable, but maybe the int ids will need to be reintroduced anyway
	return strings.Join(append([]string{name, init}, excluded...), "_")
}

func (r *registry) generator(name, init string, excluded []string) (generator, bool) {
	g, ok := r.generators[generatorKey(name, init, excluded)]
	return g, ok
}

func (r *registry) setGenerator(name, init string, excluded []string, g generator) {
	r.generators[generatorKey(name, init, excluded)] = g
}
