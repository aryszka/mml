package next

import (
	"io"
	"unicode"
)

type context struct {
	reader     io.RuneReader
	offset     int
	readOffset int
	readErr    error
	eof        bool
	cache      *cache
	tokens     []rune
	match      bool
	node       *Node
	isExcluded [][]string
}

func newContext(r io.RuneReader) *context {
	return &context{
		reader: r,
		cache:  &cache{},
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

func (c *context) excluded(offset int, name string) bool {
	if len(c.isExcluded) <= offset {
		return false
	}

	return stringsContain(c.isExcluded[offset], name)
}

func (c *context) exclude(offset int, name string) {
	if len(c.isExcluded) <= offset {
		c.isExcluded = append(c.isExcluded, nil)
		if cap(c.isExcluded) > offset {
			c.isExcluded = c.isExcluded[:offset+1]
		} else {
			c.isExcluded = append(
				c.isExcluded[:cap(c.isExcluded)],
				make([][]string, offset+1-cap(c.isExcluded))...,
			)
		}
	}

	c.isExcluded[offset] = append(c.isExcluded[offset], name)
}

func (c *context) include(offset int, name string) {
	if len(c.isExcluded) <= offset {
		return
	}

	for i := len(c.isExcluded[offset]) - 1; i >= 0; i-- {
		if c.isExcluded[offset][i] == name {
			c.isExcluded[offset] = append(c.isExcluded[offset][:i], c.isExcluded[offset][i+1:]...)
		}
	}
}

func (c *context) fromCache(name string) (bool, bool) {
	n, m, ok := c.cache.get(c.offset, name)
	if !ok {
		return false, false
	}

	if m {
		c.success(n)
	} else {
		c.fail(c.offset)
	}

	return m, true
}

func (c *context) success(n *Node) {
	c.node = n
	c.offset = n.to
	c.match = true
}

func (c *context) fail(offset int) {
	c.offset = offset
	c.match = false
}

func (c *context) finalize() error {
	if !c.eof {
		if c.offset < c.readOffset || c.read() {
			if c.readErr != nil {
				return c.readErr
			}

			return ErrUnexpectedCharacter
		}

		return c.readErr
	} else {
		if c.node.to < c.readOffset {
			return ErrUnexpectedCharacter
		}
	}

	c.node.commit()
	if c.node.commitType&Alias != 0 {
		return nil
	}

	c.node.applyTokens(c.tokens)
	return nil
}
