package next

import "testing"

func TestScheme(t *testing.T) {
	test(t, "scheme.p", "scheme", []testItem{{
		msg: "empty",
	}, {
		msg: "a function",
		text: `
			(define (foo a b c)
			  (let ([bar (+ a b c)]
			        [baz (- a b c)])
			    (* bar baz)))
		`,
		nodes: []*Node{{
			Name: "list",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "list",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "symbol",
				}, {
					Name: "symbol",
				}, {
					Name: "symbol",
				}},
			}, {
				Name: "list",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "list",
					Nodes: []*Node{{
						Name: "list",
						Nodes: []*Node{{
							Name: "symbol",
						}, {
							Name: "list",
							Nodes: []*Node{{
								Name: "symbol",
							}, {
								Name: "symbol",
							}, {
								Name: "symbol",
							}, {
								Name: "symbol",
							}},
						}},
					}, {
						Name: "list",
						Nodes: []*Node{{
							Name: "symbol",
						}, {
							Name: "list",
							Nodes: []*Node{{
								Name: "symbol",
							}, {
								Name: "symbol",
							}, {
								Name: "symbol",
							}, {
								Name: "symbol",
							}},
						}},
					}},
				}, {
					Name: "list",
					Nodes: []*Node{{
						Name: "symbol",
					}, {
						Name: "symbol",
					}, {
						Name: "symbol",
					}},
				}},
			}},
		}},
		ignorePosition: true,
	}})
}
