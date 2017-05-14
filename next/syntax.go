package next

import (
	"errors"
	"io"
)

type Syntax struct {
	registry *registry
}

var ErrNotImplemented = errors.New("not implemented")

func (s *Syntax) ReadDefinition(r io.Reader) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Terminal(name string, t ...Terminal) error {
	defs := terminalDefinitions(s.registry, name, t)
	names := make([]string, len(defs))
	for i, d := range defs {
		if err := s.registry.register(d); err != nil {
			return err
		}

		names[i] = d.nodeName()
	}

	return s.registry.register(newSequence(s.registry, name, names))
}

func (s *Syntax) Optional(string, string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Repetition(string, string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Sequence(string, ...string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Choice(string, ...string) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Init() error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Generate(w io.Writer) error {
	panic(ErrNotImplemented)
}

func (s *Syntax) Parse(r io.Reader) (*Node, error) {
	panic(ErrNotImplemented)
}
