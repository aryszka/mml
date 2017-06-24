package next

import "testing"

func TestKeyVal(t *testing.T) {
	testTrace(t, "keyval.p", "doc", 1, []testItem{{
		msg: "empty",
	}, {
		msg:  "a comment",
		text: "# a comment",
	}, {
		msg:  "a key",
		text: "a key",
		nodes: []*Node{{
			Name: "key-val",
			to:   5,
			Nodes: []*Node{{
				Name: "key",
				to:   5,
				Nodes: []*Node{{
					Name: "symbol",
					to:   5,
				}},
			}},
		}},
	}, {
		msg:  "a key with a preceeding whitespace",
		text: " a key",
		nodes: []*Node{{
			Name: "key-val",
			from: 1,
			to:   6,
			Nodes: []*Node{{
				Name: "key",
				from: 1,
				to:   6,
				Nodes: []*Node{{
					Name: "symbol",
					from: 1,
					to:   6,
				}},
			}},
		}},
	}, {
		msg: "a key and a comment",
		text: `
			# a comment

			a key
		`,
		nodes: []*Node{{
			Name: "key-val",
			from: 20,
			to:   25,
			Nodes: []*Node{{
				Name: "key",
				from: 20,
				to:   25,
				Nodes: []*Node{{
					Name: "symbol",
					from: 20,
					to:   25,
				}},
			}},
		}},
	}, {
		msg:  "a key value pair",
		text: "a key = a value",
		nodes: []*Node{{
			Name: "key-val",
			to:   15,
			Nodes: []*Node{{
				Name: "key",
				to:   5,
				Nodes: []*Node{{
					Name: "symbol",
					to:   5,
				}},
			}, {
				Name: "value",
				from: 8,
				to:   15,
			}},
		}},
	}, {
		msg:  "value without a key",
		text: "= a value",
		nodes: []*Node{{
			Name: "key-val",
			to:   9,
			Nodes: []*Node{{
				Name: "value",
				from: 2,
				to:   9,
			}},
		}},
	}, {
		msg: "a key value pair with comment",
		text: `
			# a comment
			a key = a value
		`,
		nodes: []*Node{{
			Name: "key-val",
			from: 4,
			to:   34,
			Nodes: []*Node{{
				Name: "comment",
				from: 4,
				to:   15,
			}, {
				Name: "key",
				from: 19,
				to:   24,
				Nodes: []*Node{{
					Name: "symbol",
					from: 19,
					to:   24,
				}},
			}, {
				Name: "value",
				from: 27,
				to:   34,
			}},
		}},
	}, {
		msg:  "a key with multiple symbols",
		text: "a key . with.multiple.symbols=a value",
		nodes: []*Node{{
			Name: "key-val",
			to:   37,
			Nodes: []*Node{{
				Name: "key",
				from: 0,
				to:   29,
				Nodes: []*Node{{
					Name: "symbol",
					from: 0,
					to:   5,
				}, {
					Name: "symbol",
					from: 8,
					to:   12,
				}, {
					Name: "symbol",
					from: 13,
					to:   21,
				}, {
					Name: "symbol",
					from: 22,
					to:   29,
				}},
			}, {
				Name: "value",
				from: 30,
				to:   37,
			}},
		}},
	}, {
		msg: "a group key",
		text: `
			# a comment
			[a group key.empty]
		`,
		nodes: []*Node{{
			Name: "group-key",
			from: 4,
			to:   38,
			Nodes: []*Node{{
				Name: "comment",
				from: 4,
				to:   15,
			}, {
				Name: "symbol",
				from: 20,
				to:   31,
			}, {
				Name: "symbol",
				from: 32,
				to:   37,
			}},
		}},
	}, {
		msg: "a group key with multiple values",
		text: `
			[foo.bar.baz]
			= one
			= two
			= three
		`,
		nodes: []*Node{{
			Name: "group-key",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "symbol",
			}, {
				Name: "symbol",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "value",
			}},
		}},
		ignorePosition: true,
	}, {
		msg:  "a group key with multiple values, in a single line",
		text: "[foo.bar.baz] = one = two = three",
		nodes: []*Node{{
			Name: "group-key",
			Nodes: []*Node{{
				Name: "symbol",
			}, {
				Name: "symbol",
			}, {
				Name: "symbol",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "value",
			}},
		}},
		ignorePosition: true,
	}, {
		msg: "full example",
		text: `
			# a keyval document

			key1 = foo
			key1.a = bar
			key1.b = baz

			key2 = qux

			# foo bar baz values
			[foo.bar.baz]
			a = 1
			b = 2 # even
			c = 3
		`,
		nodes: []*Node{{
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}, {
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}, {
			Name: "group-key",
			Nodes: []*Node{{
				Name: "comment",
			}, {
				Name: "symbol",
			}, {
				Name: "symbol",
			}, {
				Name: "symbol",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}, {
			Name: "key-val",
			Nodes: []*Node{{
				Name: "key",
				Nodes: []*Node{{
					Name: "symbol",
				}},
			}, {
				Name: "value",
			}},
		}},
		ignorePosition: true,
	}})
}
