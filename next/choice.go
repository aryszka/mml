package next

type choiceDefinition struct {
	name     string
	commit   CommitType
	elements []string
}

type choiceParser struct {
	name     string
	commit   CommitType
	elements []parser
}

func newChoice(name string, ct CommitType, elements []string) *choiceDefinition {
	return &choiceDefinition{
		name:     name,
		commit:   ct,
		elements: elements,
	}
}

func (d *choiceDefinition) parser(r *registry) (parser, error) {
	p, ok := r.parser(d.name)
	if ok {
		return p, nil
	}

	var elements []parser
	for _, e := range d.elements {
		element, ok := r.parser(e)
		if ok {
			elements = append(elements, element)
			continue
		}

		elementDefinition, ok := r.definition(e)
		if !ok {
			return nil, parserNotFound(e)
		}

		element, err := elementDefinition.parser(r)
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	p = &choiceParser{
		name:     d.name,
		commit:   d.commit,
		elements: elements,
	}

	r.setParser(p)
	return p, nil
}

func (p *choiceParser) parse(t Trace, c *context, excluded []string) {
	t = t.Extend(p.name)
	t.Println0("parsing choice", c.offset)

	if stringsContain(excluded, p.name) {
		t.Println0("excluded")
		c.fail(p.name, c.offset)
		return
	}

	if p.commit&Documentation != 0 {
		t.Println0("fail, doc")
		c.fail(p.name, c.offset)
		return
	}

	if c.checkCache(p.name) {
		t.Println0("found in cache")
		return
	}

	excluded = append(excluded, p.name)
	node := newNode(p.name, p.commit, c.offset, c.offset)
	var match bool

	for {
		elements := p.elements
		var foundMatch bool
		for len(elements) > 0 {
			elements[0].parse(t, c, excluded)
			elements = elements[1:]

			if !c.match || match && c.node.tokenLength() <= node.tokenLength() {
				continue
			}

			match = true
			foundMatch = true
			node.clear()
			node.append(c.node)
		}

		if !foundMatch {
			break
		}
	}

	if match {
		t.Println0("success")
		c.success(node)
		return
	}

	t.Println0("fail")
	c.fail(p.name, node.from)
}
