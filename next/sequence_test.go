package next

import (
	"bytes"
	"testing"
)

func TestSequence(t *testing.T) {
	for _, ti := range []syntaxTest{{
		msg:  "single (symbol)",
		text: "a = a",
		node: &Node{
			Name: "symbol",
			from: 4,
			to:   5,
		},
	}, {
		msg:  "multiple",
		text: "abc = a b c",
		node: &Node{
			Name: "sequence",
			from: 6,
			to:   11,
			Nodes: []*Node{{
				Name: "symbol",
				from: 6,
				to:   7,
			}, {
				Name: "symbol",
				from: 8,
				to:   9,
			}, {
				Name: "symbol",
				from: 10,
				to:   11,
			}},
		},
	}, {
		msg:  "combined",
		text: "a-f = a (b c)+ (d e? | f{1, 3})",
		node: &Node{
			Name: "sequence",
			from: 6,
			to:   31,
			Nodes: []*Node{{
				Name: "symbol",
				from: 6,
				to:   7,
			}, {
				Name: "quantifier",
				from: 8,
				to:   14,
				Nodes: []*Node{{
					Name: "sequence",
					from: 9,
					to:   12,
					Nodes: []*Node{{
						Name: "symbol",
						from: 9,
						to:   10,
					}, {
						Name: "symbol",
						from: 11,
						to:   12,
					}},
				}, {
					Name: "one-or-more",
					from: 13,
					to:   14,
				}},
			}, {
				Name: "choice",
				from: 16,
				to:   30,
				Nodes: []*Node{{
					Name: "sequence",
					from: 16,
					to:   20,
					Nodes: []*Node{{
						Name: "symbol",
						from: 16,
						to:   17,
					}, {
						Name: "quantifier",
						from: 18,
						to:   20,
						Nodes: []*Node{{
							Name: "symbol",
							from: 18,
							to:   19,
						}, {
							Name: "zero-or-one",
							from: 19,
							to:   20,
						}},
					}},
				}, {
					Name: "quantifier",
					from: 23,
					to:   30,
					Nodes: []*Node{{
						Name: "symbol",
						from: 23,
						to:   24,
					}, {
						Name: "range-quantifier",
						from: 24,
						to:   30,
						Nodes: []*Node{{
							Name: "count",
							from: 25,
							to:   26,
						}, {
							Name: "count",
							from: 28,
							to:   29,
						}},
					}},
				}},
			}},
		},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			s, err := defineSyntax()
			if err != nil {
				t.Error(err)
				return
			}

			n, err := s.Parse(bytes.NewBufferString(ti.text))
			if ti.fail && err == nil {
				t.Error("failed to fail")
			} else if !ti.fail && err != nil {
				t.Error(err)
			} else if ti.fail {
				return
			}

			checkNode(t, n.Nodes[0].Nodes[1], ti.node)
		})
	}
}
