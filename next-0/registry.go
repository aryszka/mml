package next

import (
	"fmt"
	"strings"
)

type registry struct {
	definitions map[string]definition
	root        definition
	rootSet     bool
	idSeed      int
	genIDs      map[string]int
	generators  map[int]generator
}

func duplicateDefinition(name string) error {
	return fmt.Errorf("duplicate definition in syntax: %s", name)
}

func unspecifiedParser(name string) error {
	return fmt.Errorf("unspecified parser: %s", name)
}

func newRegistry() *registry {
	return &registry{
		definitions: make(map[string]definition),
		genIDs:      make(map[string]int),
		generators:  make(map[int]generator),
	}
}

func (r *registry) register(d definition) error {
	n := d.nodeName()
	if _, exists := r.definitions[n]; exists {
		return duplicateDefinition(n)
	}

	r.definitions[n] = d
	if !r.rootSet {
		r.root = d
	}

	return nil
}

func (r *registry) setRoot(name string) error {
	d, err := r.findDefinition(name)
	if err != nil {
		return err
	}

	r.root = d
	r.rootSet = true
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

func (r *registry) genID(name, init string, excluded []string) int {
	key := generatorKey(name, init, excluded)
	if id, ok := r.genIDs[key]; ok {
		return id
	}

	id := r.idSeed
	r.idSeed++
	r.genIDs[key] = id
	return id
}

func (r *registry) generator(id int) (generator, bool) {
	g, ok := r.generators[id]
	return g, ok
}

func (r *registry) setGenerator(id int, g generator) {
	r.generators[id] = g
}
