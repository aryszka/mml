package next

import (
	"errors"
	"io"
	"unicode"
)

type definition interface {
	nodeName() string
	member(string) (bool, error)
	generator(Trace, string, []string) (generator, error)

	// TODO: try to do this during validation of the generators
	// terminates([]string) bool
}

type generator interface {
	nodeName() string
	valid() bool
	validate(Trace, []string) error
	parser(Trace, *Node) parser
}

type parser interface {
	nodeName() string
	parse(*context)
}

type context struct {
	// it is valid to hack it and provide a non unicode reader
	reader io.RuneReader

	readOffset    int
	currentOffset int
	tokens        []rune
	readErr       error

	cache *cache
	valid bool
	node  *Node
}

var ErrInvalidCharacter = errors.New("invalid character")

func stringsContain(ss []string, s string) bool {
	for _, si := range ss {
		if si == s {
			return true
		}
	}

	return false
}

// TODO: this offset is messy like this

func newContext(r io.RuneReader) *context {
	return &context{
		reader:        r,
		cache:         &cache{},
		readOffset:    -1,
		currentOffset: -1,
	}
}

func (c *context) read() bool {
	if c.readErr != nil {
		return false
	}

	t, _, err := c.reader.ReadRune()
	if err != nil {
		c.readErr = err
		return false
	}

	if t == unicode.ReplacementChar {
		c.readErr = ErrInvalidCharacter
		return false
	}

	c.tokens = append(c.tokens, t)
	return true
}

func (c *context) nextToken() (rune, bool) {
	if c.currentOffset == c.readOffset {
		c.currentOffset++
		c.readOffset++
		if !c.read() {
			return 0, false
		}
	}

	return c.tokens[c.currentOffset], true
}

func (c *context) offset() int {
	return c.currentOffset
}

func (c *context) moveOffset(d int) {
	c.currentOffset += d
}

func (c *context) succeed(n *Node) {
	c.valid = true
	c.node = n
	c.cache.set(n.From, n.Name, n)
}

func (c *context) fail(name string) {
	c.cache.set(c.currentOffset, name, nil)
	c.moveOffset(-1)
}

func (c *context) failAt(name string, offset int) {
	println(offset)
	c.cache.set(offset, name, nil)

	// TODO: offset is a mess like this
	c.currentOffset = offset - 1
}

func (c *context) fillFromCache(name string, init *Node) bool {
	if c.currentOffset < 0 {
		return false
	}

	offset := c.currentOffset
	if init != nil {
		offset = init.From
	}

	// TODO: offset is a mess like this
	n, m, ok := c.cache.get(offset+1, name)
	if !ok {
		return false
	}

	if init != nil && !n.startsWith(init) {
		return false
	}

	c.valid = m
	c.node = n
	c.moveOffset(n.To - n.From - 1)
	return true
}
