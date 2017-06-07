package next

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type CommitType int

const (
	None  CommitType = 0
	Alias CommitType = 1 << iota
	Documentation
	Root
)

type Syntax struct {
	trace       Trace
	registry    *registry
	initialized bool
	initFailed  bool
	rootSet     bool
	root        definition
	parser      parser
}

var (
	ErrSyntaxInitialized   = errors.New("syntax initialized")
	ErrInitFailed          = errors.New("init failed")
	ErrNoParsersDefined    = errors.New("no parsers defined")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidCharacter    = errors.New("invalid character") // two use cases: utf8 and boot
	ErrUnexpectedCharacter = errors.New("unexpected character")
	ErrDuplicateDefinition = errors.New("duplicate definition")
)

func NewSyntax(t Trace) *Syntax {
	if t == nil {
		t = NewTrace(0)
	}

	return &Syntax{
		trace:    t,
		registry: newRegistry(),
	}
}

func (s *Syntax) register(d definition) error {
	if s.initialized {
		return ErrSyntaxInitialized
	}

	if d.commitType()&Root != 0 {
		s.root = d
		s.rootSet = true
	} else if !s.rootSet {
		s.root = d
	}

	return s.registry.setDefinition(d)
}

func (s *Syntax) AnyChar(name string, ct CommitType) error {
	return s.register(newChar(name, ct, true, false, nil, nil))
}

func (s *Syntax) Class(name string, ct CommitType, not bool, chars []rune, ranges [][]rune) error {
	return s.register(newChar(name, ct, false, not, chars, ranges))
}

func childName(name string, childIndex int) string {
	return fmt.Sprintf("%s:%d", name, childIndex)
}

func (s *Syntax) CharSequence(name string, ct CommitType, chars []rune) error {
	var refs []string
	for i, ci := range chars {
		ref := childName(name, i)
		refs = append(refs, ref)
		if err := s.register(newChar(ref, Alias, false, false, []rune{ci}, nil)); err != nil {
			return err
		}
	}

	return s.Sequence(name, ct, refs...)
}

func (s *Syntax) Quantifier(name string, ct CommitType, item string, min, max int) error {
	return s.register(newQuantifier(name, ct, item, min, max))
}

func (s *Syntax) Sequence(name string, ct CommitType, items ...string) error {
	return s.register(newSequence(name, ct, items))
}

func (s *Syntax) Choice(name string, ct CommitType, elements ...string) error {
	return s.register(newChoice(name, ct, elements))
}

func (s *Syntax) read(self *Syntax, r io.Reader) error {
	selfTree, err := self.Parse(r)
	if err != nil {
		return err
	}

	return compile(s, selfTree)
}

func (s *Syntax) Read(r io.Reader) error {
	if s.initialized {
		return ErrSyntaxInitialized
	}

	return nil
}

func (s *Syntax) Init() error {
	if s.initFailed {
		return ErrInitFailed
	}

	if s.initialized {
		return nil
	}

	if s.root == nil {
		return ErrNoParsersDefined
	}

	var err error
	s.parser, err = s.root.parser(s.registry)
	if err != nil {
		s.initFailed = true
		return err
	}

	s.initialized = true
	return nil
}

func (s *Syntax) Generate(w io.Writer) error {
	if err := s.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Syntax) Parse(r io.Reader) (*Node, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}

	c := newContext(bufio.NewReader(r))
	return parse(s.trace, s.parser, c)
}
