package next

import (
	"bufio"
	"errors"
	"io"
	"time"
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

func (s *Syntax) Terminal(name string, ct CommitType, t ...Terminal) error {
	if len(t) == 0 {
		return ErrNoDefinitions
	}

	defs, err := terminalDefinitions(s.registry, name, t)
	if err != nil {
		return err
	}

	names := make([]string, len(defs))
	for i, d := range defs {
		if err := s.registry.register(d); err != nil {
			return err
		}

		names[i] = d.nodeName()
	}

	return s.register(newSequence(s.registry, name, ct, names))
}

func (s *Syntax) Quantifier(name string, ct CommitType, item string, min, max int) error {
	return s.register(newQuantifier(s.registry, name, ct, item, min, max))
}

func (s *Syntax) Sequence(name string, ct CommitType, items ...string) error {
	return s.register(newSequence(s.registry, name, ct, items))
}

func (s *Syntax) Choice(name string, ct CommitType, items ...string) error {
	return s.register(newChoice(s.registry, name, ct, items))
}

func (s *Syntax) Init() error {
	if s.initialized {
		return ErrSyntaxInitialized
	}

	rootDef := s.registry.root
	if rootDef == nil {
		return ErrNoDefinitions
	}

	root, ok, err := rootDef.generator(s.trace, "", nil)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidSyntax
	}

	start := time.Now()
	for {
		var foundVoid bool
		for id, g := range s.registry.generators {
			g.finalize(s.trace)
			if g.void() {
				delete(s.registry.generators, id)
				foundVoid = true
			}
		}

		if !foundVoid {
			break
		}
	}

	s.trace.Info("validation done", time.Since(start))

	if root.void() {
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
	return parse(p, c)
}
