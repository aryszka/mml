package next

import (
	"errors"
	"io"
	"unicode"
)

type definition interface {
	nodeName() string
	member(string, []string) (bool, error)
	generator(Trace, string, []string) (generator, error)

	// TODO: try to do this during validation of the generators
	// terminates([]string) bool
}

type generator interface {
	nodeName() string
	valid() bool
	validate(Trace, []generator) error
	parser(Trace, *Node) parser
}

type parser interface {
	nodeName() string
	parse(*context)
}

type context struct {
	// it is valid to hack it and provide a non unicode reader
	reader io.RuneReader

	readOffset int
	offset     int
	tokens     []rune
	readErr    error
	eof        bool

	cache *cache
	valid bool
	node  *Node
}

var (
	ErrInvalidCharacter    = errors.New("invalid character")
	ErrUnexpectedCharacter = errors.New("unexpected character")
	ErrInvalidInput        = errors.New("invalid input")
)

func stringsContain(ss []string, s string) bool {
	for _, si := range ss {
		if si == s {
			return true
		}
	}

	return false
}

func generatorsContain(gs []generator, g generator) bool {
	for _, gi := range gs {
		if gi == g {
			return true
		}
	}

	return false
}

func newContext(r io.RuneReader) *context {
	return &context{
		reader:     r,
		cache:      &cache{},
		readOffset: 0,
		offset:     0,
	}
}

func (c *context) read() bool {
	if c.eof || c.readErr != nil {
		return false
	}

	t, n, err := c.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			if n == 0 {
				c.eof = true
				return false
			}
		} else {
			c.readErr = err
			return false
		}
	}

	c.readOffset++

	if t == unicode.ReplacementChar {
		c.readErr = ErrInvalidCharacter
		return false
	}

	c.tokens = append(c.tokens, t)
	return true
}

func (c *context) token() (rune, bool) {
	if c.offset == c.readOffset {
		if !c.read() {
			return 0, false
		}
	}

	return c.tokens[c.offset], true
}

func (c *context) success(n *Node) {
	c.valid = true
	c.node = n
	c.offset = n.to
	c.cache.set(n.from, n.Name, n)
}

func (c *context) fail(name string, offset int, init *Node) {
	if init == nil {
		c.offset = offset
	} else {
		c.offset = init.to
	}

	c.valid = false

	if init != nil {
		offset = init.from
	}

	// TODO: test if it can succeed with a different init node
	c.cache.set(offset, name, nil)
}

func (c *context) fillFromCache(name string, init *Node) bool {
	offset := c.offset
	if init != nil {
		offset = init.from
	}

	n, m, ok := c.cache.get(offset, name)
	if !ok {
		return false
	}

	if init != nil && !n.startsWith(init) {
		return false
	}

	c.valid = m
	if m {
		c.node = n
		c.offset += n.to - n.from
	}

	return true
}

func useInitial(n, init *Node) bool {
	// TODO: this may need to be changed into node length
	if len(n.Nodes) == 0 {
		return true
	}

	if init == nil {
		return false
	}

	return !n.startsWith(init)
}

func (c *context) finalize() error {
	if c.eof {
		return nil
	}

	if c.read() {
		return ErrUnexpectedCharacter
	}

	return c.readErr
}

func (c *context) initRange(n, init *Node) {
	n.from = c.offset
	if init != nil {
		n.from = init.from
	}

	n.to = n.from
}

func parse(p parser, c *context) (*Node, error) {
	p.parse(c)
	if c.readErr != nil {
		return nil, c.readErr
	}

	if err := c.finalize(); err != nil {
		return nil, err
	}

	if !c.valid {
		return nil, ErrInvalidInput
	}

	if c.node.commit&Alias != 0 {
		return nil, nil
	}

	c.node.applyTokens(c.tokens)
	return c.node, nil
}
