package next

import (
	"errors"
	"fmt"
)

type definition interface {
	nodeName() string
	parser(*registry, []string) (parser, error)
	commitType() CommitType
}

type parser interface {
	nodeName() string
	setIncludedBy(parser, []string)
	cacheIncluded(*context, *Node)
	parse(Trace, *context)
}

var errCannotIncludeParsers = errors.New("cannot include parsers")

func parserNotFound(name string) error {
	return fmt.Errorf("parser not found: %s", name)
}

func stringsContain(ss []string, s string) bool {
	for _, si := range ss {
		if si == s {
			return true
		}
	}

	return false
}

func copyIncludes(to, from map[string]CommitType) {
	if from == nil {
		return
	}

	for name, ct := range from {
		to[name] = ct
	}
}

func mergeIncludes(left, right map[string]CommitType) map[string]CommitType {
	m := make(map[string]CommitType)
	copyIncludes(m, left)
	copyIncludes(m, right)
	return m
}

func parse(t Trace, p parser, c *context) (*Node, error) {
	p.parse(t, c)
	if c.readErr != nil {
		return nil, c.readErr
	}

	if !c.match {
		return nil, ErrInvalidInput
	}

	if err := c.finalize(); err != nil {
		return nil, err
	}

	return c.node, nil
}
