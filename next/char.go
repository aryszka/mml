package next

type charParser struct {
	name   string
	commit CommitType
	any    bool
	not    bool
	chars  []rune
	ranges [][]rune
}

func newChar(
	name string,
	ct CommitType,
	any, not bool,
	chars []rune,
	ranges [][]rune,
) *charParser {
	return &charParser{
		name:   name,
		commit: ct,
		any:    any,
		not:    not,
		chars:  chars,
		ranges: ranges,
	}
}

func (p *charParser) parser(*registry) (parser, error) {
	return p, nil
}

func (p *charParser) match(t rune) bool {
	if p.any {
		return true
	}

	for _, ci := range p.chars {
		if ci == t {
			return !p.not
		}
	}

	for _, ri := range p.ranges {
		if t >= ri[0] && t <= ri[1] {
			return !p.not
		}
	}

	return p.not
}

func (p *charParser) parse(t Trace, c *context, _ []string) {
	t = t.Extend(p.name)
	t.Println0("parsing char", c.offset)

	if p.commit&Documentation != 0 {
		t.Println0("fail, doc")
		c.fail(p.name, c.offset)
		return
	}

	if c.checkCache(p.name) {
		t.Println0("found in cache")
		return
	}

	if tok, ok := c.token(); ok && p.match(tok) {
		t.Println0("success", string(tok))
		c.success(newNode(p.name, p.commit, c.offset, c.offset+1))
		return
	} else {
		t.Println0("fail", string(tok))
		c.fail(p.name, c.offset)
		return
	}
}
