package next

import (
	"bufio"
	"errors"
	"io"
)

type Options struct {
	Trace Trace
}

type Syntax struct {
	trace       Trace
	registry    *registry
	initialized bool
	root        generator
}

var (
	ErrNotImplemented    = errors.New("not implemented")
	ErrSyntaxInitialized = errors.New("syntax already initialized")
	ErrNoDefinitions     = errors.New("syntax contains no definitions")
	ErrInvalidSyntax     = errors.New("invalid syntax")
	ErrInvalidInput      = errors.New("invalid input")
)

func NewSyntax(o Options) *Syntax {
	if o.Trace == nil {
		o.Trace = NewTrace(TraceInfo)
	}

	return &Syntax{
		registry: newRegistry(),
		trace:    o.Trace,
	}
}

func (s *Syntax) ReadDefinition(r io.Reader) error {
	if s.initialized {
		return ErrSyntaxInitialized
	}

	panic(ErrNotImplemented)
}

func (s *Syntax) register(d definition) error {
	if s.initialized {
		return ErrSyntaxInitialized
	}

	return s.registry.register(d)
}

func (s *Syntax) Terminal(name string, t ...Terminal) error {
	if len(t) == 0 {
		return ErrNoDefinitions
	}

	defs := terminalDefinitions(s.registry, name, t)
	names := make([]string, len(defs))
	for i, d := range defs {
		if err := s.registry.register(d); err != nil {
			return err
		}

		names[i] = d.nodeName()
	}

	return s.register(newSequence(s.registry, name, names))
}

func (s *Syntax) Optional(string, string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Repetition(string, string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Sequence(name string, items ...string) error {
	return s.register(newSequence(s.registry, name, items))
}

func (s *Syntax) Choice(string, ...string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Init() error {
	if s.initialized {
		return ErrSyntaxInitialized
	}

	rootDef := s.registry.root
	if rootDef == nil {
		return ErrNoDefinitions
	}

	root, err := rootDef.generator(s.trace, "", nil)
	if err != nil {
		return err
	}

	if err := root.validate(s.trace, nil); err != nil {
		return err
	}

	if !root.valid() {
		return ErrInvalidSyntax
	}

	s.root = root
	s.initialized = true
	return nil
}

func (s *Syntax) Generate(w io.Writer) error {
	if !s.initialized {
		if err := s.Init(); err != nil {
			return err
		}
	}

	panic(ErrNotImplemented)
}

func (s *Syntax) Parse(r io.Reader) (*Node, error) {
	if !s.initialized {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	c := newContext(bufio.NewReader(r))
	p := s.root.parser(s.trace, nil)
	p.parse(c)
	if c.readErr != nil {
		return nil, c.readErr
	}

	if !c.valid {
		return nil, ErrInvalidInput
	}

	c.node.applyTokens(c.tokens)
	return c.node, nil
}
