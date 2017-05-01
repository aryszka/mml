package mml

import (
	"bytes"
	"testing"
)

func TestParseMML(t *testing.T) {
	var l traceLevel
	trace := newTrace(l)
	s := withTrace(trace)

	err := s.newMMLSyntax()
	if err != nil {
		t.Error(err)
		return
	}

	// s.traceLevel = traceDebug

	for _, ti := range []struct {
		msg   string
		code  string
		nodes []*node
		fail  bool
	}{{
		msg: "empty document",
	}, {
		msg:  "single int",
		code: "42",
		nodes: []*node{{
			name: "int",
			token:    &token{value: "42"},
		}},
	}, {
		msg:  "multiple ints",
		code: "1 2\n3;4 ;\n 5",
		nodes: []*node{{
			name: "int",
			token:    &token{value: "1"},
		}, {
			name: "int",
			token:    &token{value: "2"},
		}, {
			name: "nl",
			token:    &token{value: "\n"},
		}, {
			name: "int",
			token:    &token{value: "3"},
		}, {
			name: "semicolon",
			token:    &token{value: ";"},
		}, {
			name: "int",
			token:    &token{value: "4"},
		}, {
			name: "semicolon",
			token:    &token{value: ";"},
		}, {
			name: "nl",
			token:    &token{value: "\n"},
		}, {
			name: "int",
			token:    &token{value: "5"},
		}},
	}, {
		msg:  "string",
		code: "\"foo\"",
		nodes: []*node{{
			name: "string",
			token:    &token{value: "\"foo\""},
		}},
	}, {
		msg:  "bool",
		code: "true false",
		nodes: []*node{{
			name: "true",
			token:    &token{value: "true"},
		}, {
			name: "false",
			token:    &token{value: "false"},
		}},
	}, {
		msg:  "symbol",
		code: "foo",
		nodes: []*node{{
			name: "symbol",
			token:    &token{value: "foo"},
		}},
	}, {
		msg:  "dynamic symbol",
		code: "symbol(f(a))",
		nodes: []*node{{
			name: "dynamic-symbol",
			token:    &token{value: "symbol"},
			nodes: []*node{{
				name: "symbol-word",
				token:    &token{value: "symbol"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "nls",
				token:    &token{value: "f"},
			}, {
				name: "function-call",
				token:    &token{value: "f"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "f"},
				}, {
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "list-sequence",
					token:    &token{value: "a"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "a"},
					}},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}},
			}, {
				name: "nls",
				token:    &token{value: ")"},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "empty list",
		code: "[]",
		nodes: []*node{{
			name: "list",
			token:    &token{value: "["},
			nodes: []*node{{
				name: "open-square",
				token:    &token{value: "["},
			}, {
				name: "list-sequence",
				token:    &token{value: "]"},
			}, {
				name: "close-square",
				token:    &token{value: "]"},
			}},
		}},
	}, {
		msg:  "list",
		code: "[1, 2, f(a), [3, 4, []]]",
		nodes: []*node{{
			name: "list",
			token:    &token{value: "["},
			nodes: []*node{{
				name: "open-square",
				token:    &token{value: "["},
			}, {
				name: "list-sequence",
				token:    &token{value: "1"},
				nodes: []*node{{
					name: "int",
					token:    &token{value: "1"},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "int",
					token:    &token{value: "2"},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "function-call",
					token:    &token{value: "f"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "f"},
					}, {
						name: "open-paren",
						token:    &token{value: "("},
					}, {
						name: "list-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							name: "symbol",
							token:    &token{value: "a"},
						}},
					}, {
						name: "close-paren",
						token:    &token{value: ")"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "list",
					token:    &token{value: "["},
					nodes: []*node{{
						name: "open-square",
						token:    &token{value: "["},
					}, {
						name: "list-sequence",
						token:    &token{value: "3"},
						nodes: []*node{{
							name: "int",
							token:    &token{value: "3"},
						}, {
							name: "comma",
							token:    &token{value: ","},
						}, {
							name: "int",
							token:    &token{value: "4"},
						}, {
							name: "comma",
							token:    &token{value: ","},
						}, {
							name: "list",
							token:    &token{value: "["},
							nodes: []*node{{
								name: "open-square",
								token:    &token{value: "["},
							}, {
								name: "list-sequence",
								token:    &token{value: "]"},
							}, {
								name: "close-square",
								token:    &token{value: "]"},
							}},
						}},
					}, {
						name: "close-square",
						token:    &token{value: "]"},
					}},
				}},
			}, {
				name: "close-square",
				token:    &token{value: "]"},
			}},
		}},
	}, {
		msg:  "mutable list",
		code: "~[1, 2, f(a), [3, 4, ~[]]]",
		nodes: []*node{{
			name: "mutable-list",
			token:    &token{value: "~"},
			nodes: []*node{{
				name: "tilde",
				token:    &token{value: "~"},
			}, {
				name: "open-square",
				token:    &token{value: "["},
			}, {
				name: "list-sequence",
				token:    &token{value: "1"},
				nodes: []*node{{
					name: "int",
					token:    &token{value: "1"},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "int",
					token:    &token{value: "2"},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "function-call",
					token:    &token{value: "f"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "f"},
					}, {
						name: "open-paren",
						token:    &token{value: "("},
					}, {
						name: "list-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							name: "symbol",
							token:    &token{value: "a"},
						}},
					}, {
						name: "close-paren",
						token:    &token{value: ")"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "list",
					token:    &token{value: "["},
					nodes: []*node{{
						name: "open-square",
						token:    &token{value: "["},
					}, {
						name: "list-sequence",
						token:    &token{value: "3"},
						nodes: []*node{{
							name: "int",
							token:    &token{value: "3"},
						}, {
							name: "comma",
							token:    &token{value: ","},
						}, {
							name: "int",
							token:    &token{value: "4"},
						}, {
							name: "comma",
							token:    &token{value: ","},
						}, {
							name: "mutable-list",
							token:    &token{value: "~"},
							nodes: []*node{{
								name: "tilde",
								token:    &token{value: "~"},
							}, {
								name: "open-square",
								token:    &token{value: "["},
							}, {
								name: "list-sequence",
								token:    &token{value: "]"},
							}, {
								name: "close-square",
								token:    &token{value: "]"},
							}},
						}},
					}, {
						name: "close-square",
						token:    &token{value: "]"},
					}},
				}},
			}, {
				name: "close-square",
				token:    &token{value: "]"},
			}},
		}},
	}, {
		msg:  "empty structure",
		code: "{}",
		nodes: []*node{{
			name: "structure",
			token:    &token{value: "{"},
			nodes: []*node{{
				name: "open-brace",
				token:    &token{value: "{"},
			}, {
				name: "structure-sequence",
				token:    &token{value: "}"},
			}, {
				name: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg:  "structure",
		code: "{a: 1, b: 2, ...c, d: {e: 3, f: {}}}",
		nodes: []*node{{
			name: "structure",
			token:    &token{value: "{"},
			nodes: []*node{{
				name: "open-brace",
				token:    &token{value: "{"},
			}, {
				name: "structure-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					name: "structure-definition",
					token:    &token{value: "a"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "a"},
					}, {
						name: "nls",
						token:    &token{value: ":"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "nls",
						token:    &token{value: "1"},
					}, {
						name: "int",
						token:    &token{value: "1"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "structure-definition",
					token:    &token{value: "b"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "b"},
					}, {
						name: "nls",
						token:    &token{value: ":"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "nls",
						token:    &token{value: "2"},
					}, {
						name: "int",
						token:    &token{value: "2"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "spread-expression",
					token:    &token{value: "."},
					nodes: []*node{{
						name: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}},
					}, {
						name: "symbol",
						token:    &token{value: "c"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "structure-definition",
					token:    &token{value: "d"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "d"},
					}, {
						name: "nls",
						token:    &token{value: ":"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "nls",
						token:    &token{value: "{"},
					}, {
						name: "structure",
						token:    &token{value: "{"},
						nodes: []*node{{
							name: "open-brace",
							token:    &token{value: "{"},
						}, {
							name: "structure-sequence",
							token:    &token{value: "e"},
							nodes: []*node{{
								name: "structure-definition",
								token:    &token{value: "e"},
								nodes: []*node{{
									name: "symbol",
									token:    &token{value: "e"},
								}, {
									name: "nls",
									token:    &token{value: ":"},
								}, {
									name: "colon",
									token:    &token{value: ":"},
								}, {
									name: "nls",
									token:    &token{value: "3"},
								}, {
									name: "int",
									token:    &token{value: "3"},
								}},
							}, {
								name: "comma",
								token:    &token{value: ","},
							}, {
								name: "structure-definition",
								token:    &token{value: "f"},
								nodes: []*node{{
									name: "symbol",
									token:    &token{value: "f"},
								}, {
									name: "nls",
									token:    &token{value: ":"},
								}, {
									name: "colon",
									token:    &token{value: ":"},
								}, {
									name: "nls",
									token:    &token{value: "{"},
								}, {
									name: "structure",
									token:    &token{value: "{"},
									nodes: []*node{{
										name: "open-brace",
										token:    &token{value: "{"},
									}, {
										name: "structure-sequence",
										token:    &token{value: "}"},
									}, {
										name: "close-brace",
										token:    &token{value: "}"},
									}},
								}},
							}},
						}, {
							name: "close-brace",
							token:    &token{value: "}"},
						}},
					}},
				}},
			}, {
				name: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg:  "mutable structure",
		code: "~{a: 1, b: 2, ...c, d: {e: 3, f: ~{}}}",
		nodes: []*node{{
			name: "mutable-structure",
			token:    &token{value: "~"},
			nodes: []*node{{
				name: "tilde",
				token:    &token{value: "~"},
			}, {
				name: "open-brace",
				token:    &token{value: "{"},
			}, {
				name: "structure-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					name: "structure-definition",
					token:    &token{value: "a"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "a"},
					}, {
						name: "nls",
						token:    &token{value: ":"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "nls",
						token:    &token{value: "1"},
					}, {
						name: "int",
						token:    &token{value: "1"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "structure-definition",
					token:    &token{value: "b"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "b"},
					}, {
						name: "nls",
						token:    &token{value: ":"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "nls",
						token:    &token{value: "2"},
					}, {
						name: "int",
						token:    &token{value: "2"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "spread-expression",
					token:    &token{value: "."},
					nodes: []*node{{
						name: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}},
					}, {
						name: "symbol",
						token:    &token{value: "c"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "structure-definition",
					token:    &token{value: "d"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "d"},
					}, {
						name: "nls",
						token:    &token{value: ":"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "nls",
						token:    &token{value: "{"},
					}, {
						name: "structure",
						token:    &token{value: "{"},
						nodes: []*node{{
							name: "open-brace",
							token:    &token{value: "{"},
						}, {
							name: "structure-sequence",
							token:    &token{value: "e"},
							nodes: []*node{{
								name: "structure-definition",
								token:    &token{value: "e"},
								nodes: []*node{{
									name: "symbol",
									token:    &token{value: "e"},
								}, {
									name: "nls",
									token:    &token{value: ":"},
								}, {
									name: "colon",
									token:    &token{value: ":"},
								}, {
									name: "nls",
									token:    &token{value: "3"},
								}, {
									name: "int",
									token:    &token{value: "3"},
								}},
							}, {
								name: "comma",
								token:    &token{value: ","},
							}, {
								name: "structure-definition",
								token:    &token{value: "f"},
								nodes: []*node{{
									name: "symbol",
									token:    &token{value: "f"},
								}, {
									name: "nls",
									token:    &token{value: ":"},
								}, {
									name: "colon",
									token:    &token{value: ":"},
								}, {
									name: "nls",
									token:    &token{value: "~"},
								}, {
									name: "mutable-structure",
									token:    &token{value: "~"},
									nodes: []*node{{
										name: "tilde",
										token:    &token{value: "~"},
									}, {
										name: "open-brace",
										token:    &token{value: "{"},
									}, {
										name: "structure-sequence",
										token:    &token{value: "}"},
									}, {
										name: "close-brace",
										token:    &token{value: "}"},
									}},
								}},
							}},
						}, {
							name: "close-brace",
							token:    &token{value: "}"},
						}},
					}},
				}},
			}, {
				name: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg:  "symbol query",
		code: "a.b",
		nodes: []*node{{
			name: "symbol-query",
			token:    &token{value: "a"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "a"},
			}, {
				name: "dot",
				token:    &token{value: "."},
			}, {
				name: "symbol",
				token:    &token{value: "b"},
			}},
		}},
	}, {
		msg:  "chained symbol query",
		code: "a.b.c",
		nodes: []*node{{
			name: "symbol-query",
			token:    &token{value: "a"},
			nodes: []*node{{
				name: "symbol-query",
				token:    &token{value: "a"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "a"},
				}, {
					name: "dot",
					token:    &token{value: "."},
				}, {
					name: "symbol",
					token:    &token{value: "b"},
				}},
			}, {
				name: "dot",
				token:    &token{value: "."},
			}, {
				name: "symbol",
				token:    &token{value: "c"},
			}},
		}},
	}, {
		msg:  "void function",
		code: "fn () {;}",
		nodes: []*node{{
			name: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				name: "fn-word",
				token:    &token{value: "fn"},
			}, {
				name: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "static-symbol-sequence",
					token:    &token{value: ")"},
				}, {
					name: "nls",
					token:    &token{value: ")"},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}, {
					name: "nls",
					token:    &token{value: "{"},
				}, {
					name: "function-body",
					token:    &token{value: "{"},
					nodes: []*node{{
						name: "open-brace",
						token:    &token{value: "{"},
					}, {
						name: "statement-sequence",
						token:    &token{value: ";"},
						nodes: []*node{{
							name: "semicolon",
							token:    &token{value: ";"},
						}},
					}, {
						name: "close-brace",
						token:    &token{value: "}"},
					}},
				}},
			}},
		}},
	}, {
		msg:  "identity",
		code: "fn (x) x",
		nodes: []*node{{
			name: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				name: "fn-word",
				token:    &token{value: "fn"},
			}, {
				name: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "static-symbol-sequence",
					token:    &token{value: "x"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "x"},
					}},
				}, {
					name: "nls",
					token:    &token{value: ")"},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}, {
					name: "nls",
					token:    &token{value: "x"},
				}, {
					name: "symbol",
					token:    &token{value: "x"},
				}},
			}},
		}},
	}, {
		msg:  "list as a function",
		code: "fn (...x) x",
		nodes: []*node{{
			name: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				name: "fn-word",
				token:    &token{value: "fn"},
			}, {
				name: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "static-symbol-sequence",
					token:    &token{value: "."},
				}, {
					name: "collect-symbol",
					token:    &token{value: "."},
					nodes: []*node{{
						name: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}},
					}, {
						name: "symbol",
						token:    &token{value: "x"},
					}},
				}, {
					name: "nls",
					token:    &token{value: ")"},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}, {
					name: "nls",
					token:    &token{value: "x"},
				}, {
					name: "symbol",
					token:    &token{value: "x"},
				}},
			}},
		}},
	}, {
		msg:  "function with sequence",
		code: "fn (a, b, ...c) { a(b); c }",
		nodes: []*node{{
			name: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				name: "fn-word",
				token:    &token{value: "fn"},
			}, {
				name: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "static-symbol-sequence",
					token:    &token{value: "a"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "a"},
					}, {
						name: "comma",
						token:    &token{value: ","},
					}, {
						name: "symbol",
						token:    &token{value: "b"},
					}, {
						name: "comma",
						token:    &token{value: ","},
					}},
				}, {
					name: "collect-symbol",
					token:    &token{value: "."},
					nodes: []*node{{
						name: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}},
					}, {
						name: "symbol",
						token:    &token{value: "c"},
					}},
				}, {
					name: "nls",
					token:    &token{value: ")"},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}, {
					name: "nls",
					token:    &token{value: "{"},
				}, {
					name: "function-body",
					token:    &token{value: "{"},
					nodes: []*node{{
						name: "open-brace",
						token:    &token{value: "{"},
					}, {
						name: "statement-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							name: "function-call",
							token:    &token{value: "a"},
							nodes: []*node{{
								name: "symbol",
								token:    &token{value: "a"},
							}, {
								name: "open-paren",
								token:    &token{value: "("},
							}, {
								name: "list-sequence",
								token:    &token{value: "b"},
								nodes: []*node{{
									name: "symbol",
									token:    &token{value: "b"},
								}},
							}, {
								name: "close-paren",
								token:    &token{value: ")"},
							}},
						}, {
							name: "semicolon",
							token:    &token{value: ";"},
						}, {
							name: "symbol",
							token:    &token{value: "c"},
						}},
					}, {
						name: "close-brace",
						token:    &token{value: "}"},
					}},
				}},
			}},
		}},
	}, {
		msg:  "function call",
		code: "f(a)",
		nodes: []*node{{
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "f"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "a"},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "chained function call",
		code: "f(a)(b)",
		nodes: []*node{{
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "function-call",
				token:    &token{value: "f"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "f"},
				}, {
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "list-sequence",
					token:    &token{value: "a"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "a"},
					}},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "b"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "b"},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "chained function call, whitespace",
		code: "f(a) (b)",
		nodes: []*node{{
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "function-call",
				token:    &token{value: "f"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "f"},
				}, {
					name: "open-paren",
					token:    &token{value: "("},
				}, {
					name: "list-sequence",
					token:    &token{value: "a"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "a"},
					}},
				}, {
					name: "close-paren",
					token:    &token{value: ")"},
				}},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "b"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "b"},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "function call argument",
		code: "f(g(a))",
		nodes: []*node{{
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "f"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "g"},
				nodes: []*node{{
					name: "function-call",
					token:    &token{value: "g"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "g"},
					}, {
						name: "open-paren",
						token:    &token{value: "("},
					}, {
						name: "list-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							name: "symbol",
							token:    &token{value: "a"},
						}},
					}, {
						name: "close-paren",
						token:    &token{value: ")"},
					}},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "function call sequence",
		code: "f(a) f(b)g(a)",
		nodes: []*node{{
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "f"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "a"},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}, {
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "f"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "b"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "b"},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}, {
			name: "function-call",
			token:    &token{value: "g"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "g"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					name: "symbol",
					token:    &token{value: "a"},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "function call with multiple arguments",
		code: "f(...a, b, ...c)",
		nodes: []*node{{
			name: "function-call",
			token:    &token{value: "f"},
			nodes: []*node{{
				name: "symbol",
				token:    &token{value: "f"},
			}, {
				name: "open-paren",
				token:    &token{value: "("},
			}, {
				name: "list-sequence",
				token:    &token{value: "."},
				nodes: []*node{{
					name: "spread-expression",
					token:    &token{value: "."},
					nodes: []*node{{
						name: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}},
					}, {
						name: "symbol",
						token:    &token{value: "a"},
					}},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "symbol",
					token:    &token{value: "b"},
				}, {
					name: "comma",
					token:    &token{value: ","},
				}, {
					name: "spread-expression",
					token:    &token{value: "."},
					nodes: []*node{{
						name: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}, {
							name: "dot",
							token:    &token{value: "."},
						}},
					}, {
						name: "symbol",
						token:    &token{value: "c"},
					}},
				}},
			}, {
				name: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "switch conditional with default only",
		code: "switch{default: 42}",
		nodes: []*node{{
			name: "switch-conditional",
			token:    &token{value: "switch"},
			nodes: []*node{{
				name: "switch-word",
				token:    &token{value: "switch"},
			}, {
				name: "nls",
				token:    &token{value: "{"},
			}, {
				name: "open-brace",
				token:    &token{value: "{"},
			}, {
				name: "nls",
				token:    &token{value: "default"},
			}, {
				name: "switch-clause-sequence",
				token:    &token{value: "default"},
			}, {
				name: "nls",
				token:    &token{value: "default"},
			}, {
				name: "default-clause",
				token:    &token{value: "default"},
				nodes: []*node{{
					name: "default-word",
					token:    &token{value: "default"},
				}, {
					name: "colon",
					token:    &token{value: ":"},
				}, {
					name: "statement-sequence",
					token:    &token{value: "42"},
					nodes: []*node{{
						name: "int",
						token:    &token{value: "42"},
					}},
				}},
			}, {
				name: "nls",
				token:    &token{value: "}"},
			}, {
				name: "switch-clause-sequence",
				token:    &token{value: "}"},
			}, {
				name: "nls",
				token:    &token{value: "}"},
			}, {
				name: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg: "switch conditional with cases",
		code: `
					switch {
						case a: b
						default: x
						case c: d
					}`,
		nodes: []*node{{
			name: "nl",
			token:    &token{value: "\n"},
		}, {
			name: "switch-conditional",
			token:    &token{value: "switch"},
			nodes: []*node{{
				name: "switch-word",
				token:    &token{value: "switch"},
			}, {
				name: "nls",
				token:    &token{value: "{"},
			}, {
				name: "open-brace",
				token:    &token{value: "{"},
			}, {
				name: "nls",
				token:    &token{value: "\n"},
				nodes: []*node{{
					name: "nl",
					token:    &token{value: "\n"},
				}},
			}, {
				name: "switch-clause-sequence",
				token:    &token{value: "case"},
				nodes: []*node{{
					name: "switch-clause",
					token:    &token{value: "case"},
					nodes: []*node{{
						name: "case-word",
						token:    &token{value: "case"},
					}, {
						name: "symbol",
						token:    &token{value: "a"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "statement-sequence",
						token:    &token{value: "b"},
						nodes: []*node{{
							name: "symbol",
							token:    &token{value: "b"},
						}, {
							name: "nl",
							token:    &token{value: "\n"},
						}},
					}},
				}},
			}, {
				name: "nls",
				token:    &token{value: "default"},
			}, {
				name: "default-clause",
				token:    &token{value: "default"},
				nodes: []*node{{
					name: "default-word",
					token:    &token{value: "default"},
				}, {
					name: "colon",
					token:    &token{value: ":"},
				}, {
					name: "statement-sequence",
					token:    &token{value: "x"},
					nodes: []*node{{
						name: "symbol",
						token:    &token{value: "x"},
					}, {
						name: "nl",
						token:    &token{value: "\n"},
					}},
				}},
			}, {
				name: "nls",
				token:    &token{value: "case"},
			}, {
				name: "switch-clause-sequence",
				token:    &token{value: "case"},
				nodes: []*node{{
					name: "switch-clause",
					token:    &token{value: "case"},
					nodes: []*node{{
						name: "case-word",
						token:    &token{value: "case"},
					}, {
						name: "symbol",
						token:    &token{value: "c"},
					}, {
						name: "colon",
						token:    &token{value: ":"},
					}, {
						name: "statement-sequence",
						token:    &token{value: "d"},
						nodes: []*node{{
							name: "symbol",
							token:    &token{value: "d"},
						}, {
							name: "nl",
							token:    &token{value: "\n"},
						}},
					}},
				}},
			}, {
				name: "nls",
				token:    &token{value: "}"},
			}, {
				name: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			b := bytes.NewBufferString(ti.code)

			n, err := s.parse(b, "test")
			if !ti.fail && err != nil {
				t.Error(err)
				return
			} else if ti.fail && err == nil {
				t.Error("failed to fail")
				return
			}

			if ti.fail {
				return
			}

			if n.name != "statement-sequence" {
				t.Error("invalid root node type", n.name, "statement-sequence")
				return
			}

			if len(n.nodes) != len(ti.nodes) {
				t.Error("invalid number of nodes", len(n.nodes), len(ti.nodes))
				return
			}

			if len(n.nodes) == 0 && n.token.typ != eofTokenType || len(n.nodes) > 0 && n.token != n.nodes[0].token {
				t.Error("invalid document token", n.token)
				return
			}

			for i, ni := range n.nodes {
				if !checkNodes(ni, ti.nodes[i]) {
					t.Error("failed to match nodes")
					t.Log(ni)
					t.Log(ti.nodes[i])
				}
			}
		})
	}
}
