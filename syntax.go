package mml

import (
	"errors"
	"fmt"
	"io"
)

type definition interface {
	typeName() string
	nodeType() nodeType
	member(nodeType) (bool, error)
	generator(trace, nodeType, typeList) (generator, error)
}

type generator interface {
	typeName() string
	nodeType() nodeType
	finalize(trace) error
	valid() bool
	parser(trace, *cache, *node) parser
}

type parser interface {
	typeName() string
	nodeType() nodeType
	parse(*token) *parserResult
}

type parserResult struct {
	accepting bool
	valid     bool
	node      *node
	fromCache bool
	unparsed  *tokenStack
}

type syntax struct {
	trace         trace
	registry      *registry
	initDone      bool
	rootGenerator generator
	cache         *cache
}

var (
	errNotImplemented       = errors.New("not implemented")
	errDefinitionsClosed    = errors.New("definitions closed")
	errFailedToCreateParser = errors.New("failed to create parser")
	errUnexpectedEOF        = errors.New("unexpected EOF")
)

func unexpectedToken(nodeType string, t *token) error {
	return fmt.Errorf("unexpected token: %v, %v", nodeType, t)
}

func unspecifiedParser(typeName string) error {
	return fmt.Errorf("unspecified parser: %s", typeName)
}

func requiredParserInvalid(typeName string) error {
	return fmt.Errorf("required parser invalid: %s", typeName)
}

func unexpectedResult(nodeType string) error {
	return fmt.Errorf("unexpected parse result: %s", nodeType)
}

func newSyntax() *syntax {
	return withTrace(noopTrace{})
}

func withTrace(t trace) *syntax {
	return &syntax{
		trace:    t,
		registry: newRegistry(),
		cache:    &cache{},
	}
}

func readSyntax(io.Reader) (*syntax, error) {
	panic(errNotImplemented)
}

func (s *syntax) primitive(name string, t tokenType) error {
	if s.initDone {
		return errDefinitionsClosed
	}

	nt := s.registry.nodeType(name)
	d := newPrimitive(s.registry, name, nt, t)
	return d.registry.register(d)
}

func (s *syntax) optional(name string, optional string) error {
	if s.initDone {
		return errDefinitionsClosed
	}

	nt := s.registry.nodeType(name)
	ot := s.registry.nodeType(optional)
	d := newOptional(s.registry, name, nt, optional, ot)
	return d.registry.register(d)
}

func (s *syntax) repeat(name string, item string) error {
	if s.initDone {
		return errDefinitionsClosed
	}

	nt := s.registry.nodeType(name)
	it := s.registry.nodeType(item)
	d := newRepeat(s.registry, name, nt, item, it)
	return d.registry.register(d)
}

func (s *syntax) sequence(name string, itemNames ...string) error {
	if s.initDone {
		return errDefinitionsClosed
	}

	nt := s.registry.nodeType(name)
	it := make([]nodeType, len(itemNames))
	for i, ni := range itemNames {
		it[i] = s.registry.nodeType(ni)
	}

	d := newSequence(s.registry, name, nt, itemNames, it)
	return d.registry.register(d)
}

func (s *syntax) choice(name string, elementNames ...string) error {
	if s.initDone {
		return errDefinitionsClosed
	}

	nt := s.registry.nodeType(name)
	et := make([]nodeType, len(elementNames))
	for i, ni := range elementNames {
		et[i] = s.registry.nodeType(ni)
	}

	d := newChoice(s.registry, name, nt, elementNames, et)
	return d.registry.register(d)
}

func (s *syntax) root() (string, error) {
	return s.registry.root()
}

func (s *syntax) setRoot(name string) error {
	if s.initDone {
		return errDefinitionsClosed
	}

	return s.registry.setRoot(name)
}

func (s *syntax) finalizeInit() error {
	var done bool
	for !done {
		done = true
		for k, g := range s.registry.generators {
			if !g.valid() {
				delete(s.registry.generators, k)
				done = false
				continue
			}

			if err := g.finalize(s.trace); err != nil {
				return err
			}

			if !g.valid() {
				delete(s.registry.generators, k)
				done = false
			}
		}
	}

	return nil
}

func (s *syntax) init() error {
	if s.initDone {
		return nil
	}

	rn, err := s.root()
	if err != nil {
		return err
	}

	rt := s.registry.nodeType(rn)

	d, ok := s.registry.definition(rt)
	if !ok {
		return unspecifiedParser(rn)
	}

	g, err := d.generator(s.trace, 0, nil)
	if err != nil {
		return err
	}

	if err := s.finalizeInit(); err != nil {
		return err
	}

	if !g.valid() {
		return errFailedToCreateParser
	}

	s.rootGenerator = g
	s.initDone = true

	return nil
}

func (s *syntax) generate(io.Writer) error {
	panic(errNotImplemented)

	if !s.initDone {
		if err := s.init(); err != nil {
			return err
		}
	}

	return nil
}

func (s *syntax) parse(r io.Reader, name string) (*node, error) {
	if !s.initDone {
		if err := s.init(); err != nil {
			return nil, err
		}
	}

	s.cache.clear()
	p := s.rootGenerator.parser(s.trace, s.cache, nil)
	tr := newTokenReader(r, name)
	last := &parserResult{accepting: true}
	eof := &token{typ: eofTokenType}
	var offset int
	for {
		t, err := tr.next()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if !last.accepting {
			if err != io.EOF {
				return nil, unexpectedToken("root", &t)
			}

			if !last.valid {
				return nil, errUnexpectedEOF
			}

			return last.node, nil
		}

		if err == io.EOF {
			eof.offset = offset + 1
			last = p.parse(eof)

			if !last.valid {
				return nil, errUnexpectedEOF
			}

			if !last.unparsed.has() || last.unparsed.peek() != eof {
				return nil, errUnexpectedEOF
			}

			return last.node, nil
		}

		offset = t.offset

		last = p.parse(&t)
		if !last.accepting {
			if !last.valid {
				return nil, unexpectedToken("root", &t)
			}

			if last.unparsed != nil && last.unparsed.has() {
				return nil, unexpectedToken("root", last.unparsed.peek())
			}
		}
	}
}
