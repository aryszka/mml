package mml

import (
	"fmt"
	"errors"
	"strings"
)

type registry struct {
	typeSeed    nodeType
	names       map[nodeType]string
	types       map[string]nodeType
	definitions map[nodeType]definition
	rootDef string
	generators map[string]generator
}

var (
	errNoParsersDefined = errors.New("no parsers defined")
)

func duplicateNodeType(nodeType string) error {
	return fmt.Errorf("duplicate node type definition in syntax: %s", nodeType)
}

func newRegistry() *registry {
	return &registry{
		names: make(map[nodeType]string),
		types: make(map[string]nodeType),
		definitions: make(map[nodeType]definition),
		generators: make(map[string]generator),
	}
}

func (r *registry) nodeType(name string) nodeType {
	t, ok := r.types[name]
	if ok {
		return t
	}

	t = r.typeSeed
	r.typeSeed++
	r.types[name] = t
	r.names[t] = name
	return t
}

func (r *registry) definition(t nodeType) (definition, bool) {
	g, ok := r.definitions[t]
	return g, ok
}

func (r *registry) register(d definition) error {
	if _, exists := r.definitions[d.nodeType()]; exists {
		return duplicateNodeType(d.typeName())
	}

	r.definitions[d.nodeType()] = d
	r.rootDef = d.typeName() // the last one is the root by default
	return nil
}

func (r *registry) root() (string, error) {
	if len(r.definitions) == 0 {
		return "", errNoParsersDefined
	}

	return r.rootDef, nil
}

func (r *registry) setRoot(name string) error {
	if len(r.definitions) == 0 {
		return errNoParsersDefined
	}

	if _, exists := r.definitions[r.nodeType(name)]; !exists {
		return unspecifiedParser(name)
	}

	r.rootDef = name
	return nil
}

func generatorKey(t nodeType, init nodeType, excluded typeList) string {
	s := make([]string, len(excluded)+2)
	for i, ni := range append([]nodeType{t, init}, excluded...) {
		s[i] = fmt.Sprint(ni)
	}

	return strings.Join(s, "_")
}

func (r *registry) generator(t nodeType, init nodeType, excluded typeList) (generator, bool) {
	g, ok := r.generators[generatorKey(t, init, excluded)]
	return g, ok
}

func (r *registry) setGenerator(t nodeType, init nodeType, excluded typeList, g generator) {
	r.generators[generatorKey(t, init, excluded)] = g
}
