package next

import "testing"

func TestJSON(t *testing.T) {
	test(t, "json.p", "value", []testItem{{
		msg:  "true",
		text: "true",
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "true",
			}},
		},
		ignorePosition: true,
	}, {
		msg:  "false",
		text: "false",
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "false",
			}},
		},
		ignorePosition: true,
	}, {
		msg:  "null",
		text: "null",
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "null",
			}},
		},
		ignorePosition: true,
	}, {
		msg:  "string",
		text: `"\"\\n\b\t\uabcd"`,
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "string",
			}},
		},
		ignorePosition: true,
	}, {
		msg:  "number",
		text: "6.62e-34",
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "number",
			}},
		},
		ignorePosition: true,
	}, {
		msg: "object",
		text: `{
			"true": true,
			"false": false,
			"null": null,
			"string": "string",
			"number": 42,
			"object": {},
			"array": []
		}`,
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "object",
				Nodes: []*Node{{
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "true",
					}},
				}, {
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "false",
					}},
				}, {
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "null",
					}},
				}, {
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "string",
					}},
				}, {
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "number",
					}},
				}, {
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "object",
					}},
				}, {
					Name: "entry",
					Nodes: []*Node{{
						Name: "string",
					}, {
						Name: "array",
					}},
				}},
			}},
		},
		ignorePosition: true,
	}, {
		msg: "array",
		text: `[true, false, null, "string", 42, {
			"true": true,
			"false": false,
			"null": null,
			"string": "string",
			"number": 42,
			"object": {},
			"array": []
		}, []]`,
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "array",
				Nodes: []*Node{{
					Name: "true",
				}, {
					Name: "false",
				}, {
					Name: "null",
				}, {
					Name: "string",
				}, {
					Name: "number",
				}, {
					Name: "object",
					Nodes: []*Node{{
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "true",
						}},
					}, {
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "false",
						}},
					}, {
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "null",
						}},
					}, {
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "string",
						}},
					}, {
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "number",
						}},
					}, {
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "object",
						}},
					}, {
						Name: "entry",
						Nodes: []*Node{{
							Name: "string",
						}, {
							Name: "array",
						}},
					}},
				}, {
					Name: "array",
				}},
			}},
		},
		ignorePosition: true,
	}})
}
