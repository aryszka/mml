package next

import "fmt"

type definition interface {
	nodeName() string
	parser(*registry) (parser, error)
	commitType() CommitType
}

type parser interface {
	nodeName() string
	parse(Trace, *context, []string)
}

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

func parse(t Trace, p parser, c *context) (*Node, error) {
	p.parse(t, c, nil)
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
