package next

import "testing"

func TestSExpr(t *testing.T) {
	test(t, "sexpr.p", "s-expression", []testItem{{
		msg:  "number",
		text: "42",
		nodes: []*Node{{
			Name: "number",
		}},
		ignorePosition: true,
	}, {
		msg:  "string",
		text: "\"foo\"",
		nodes: []*Node{{
			Name: "string",
		}},
		ignorePosition: true,
	}, {
		msg:  "symbol",
		text: "foo",
		nodes: []*Node{{
			Name: "symbol",
		}},
		ignorePosition: true,
	}, {
		msg:  "nil",
		text: "()",
		nodes: []*Node{{
			Name: "list",
		}},
		ignorePosition: true,
	}, {
		msg:  "list",
		text: "(foo bar baz)",
		nodes: []*Node{{
			Name: "list",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "symbol",
			}, {
				Name: "symbol",
			}},
		}},
		ignorePosition: true,
	}, {
		msg:  "embedded list",
		text: "(foo (bar (baz)) qux)",
		nodes: []*Node{{
			Name: "list",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "list",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "list",
					Nodes: []*Node{{
						Name: "symbol",
					}},
				}},
			}, {
				Name: "symbol",
			}},
		}},
		ignorePosition: true,
	}})
}
