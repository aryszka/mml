package next

import (
	"errors"
	"io"
	"log"
	"time"
	"unicode"
)

type definition interface {
	nodeName() string
	generator(Trace, string, []string) (generator, bool, error)
}

type generator interface {
	nodeName() string
	parser(Trace, *Node) parser
}

type parser interface {
	nodeName() string
	parse(*context)
}

type context struct {
	// it is valid to hack it and provide a non unicode reader
	// potential optimization: 8-bit tokens flag
	reader io.RuneReader

	readOffset int
	match      bool
	offset     int
	tokens     []rune
	readErr    error
	eof        bool

	cache *cache
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

func (c *context) success(genID int, n *Node) {
	c.match = true
	c.node = n
	c.offset = n.to
	c.cache.set(n.from, genID, n)
}

func (c *context) fail(genID, offset int) {
	c.match = false
	c.offset = offset
	c.cache.set(offset, genID, nil)
}

func (c *context) fillFromCache(genID int, init *Node) bool {
	offset := c.offset
	if init != nil {
		offset = init.from
	}

	n, m, ok := c.cache.get(offset, genID)
	if !ok {
		return false
	}

	if init != nil && !n.startsWith(init) {
		return false
	}

	c.match = m
	if m {
		c.node = n
		c.offset += n.to - n.from
	}

	return true
}

func (c *context) finalize() error {
	if c.eof {
		return nil
	}

	if c.offset < c.readOffset || c.read() {
		if c.readErr != nil {
			return c.readErr
		}

		return ErrUnexpectedCharacter
	}

	return c.readErr
}

func (c *context) initNode(n, init *Node) {
	n.from = c.offset
	if init != nil {
		n.from = init.from
	}

	n.to = n.from
}

func parse(p parser, c *context) (*Node, error) {
	start := time.Now()
	p.parse(c)
	log.Println("parse time", time.Now().Sub(start))

	if c.readErr != nil {
		return nil, c.readErr
	}

	if err := c.finalize(); err != nil {
		return nil, err
	}

	if !c.match {
		return nil, ErrInvalidInput
	}

	if c.node.commitType&Alias != 0 {
		return nil, nil
	}

	c.node.applyTokens(c.tokens)
	return c.node, nil
}
