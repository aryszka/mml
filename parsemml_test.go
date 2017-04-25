package mml

import (
	"bytes"
	"testing"
)

func TestParseMML(t *testing.T) {
	s, err := newMMLSyntax()
	if err != nil {
		t.Error(err)
		return
	}

	s.traceLevel = traceDebug

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
			typeName: "int",
			token:    &token{value: "42"},
		}},
	}, {
		msg:  "multiple ints",
		code: "1 2\n3;4 ;\n 5",
		nodes: []*node{{
			typeName: "int",
			token:    &token{value: "1"},
		}, {
			typeName: "int",
			token:    &token{value: "2"},
		}, {
			typeName: "nl",
			token:    &token{value: "\n"},
		}, {
			typeName: "int",
			token:    &token{value: "3"},
		}, {
			typeName: "semicolon",
			token:    &token{value: ";"},
		}, {
			typeName: "int",
			token:    &token{value: "4"},
		}, {
			typeName: "semicolon",
			token:    &token{value: ";"},
		}, {
			typeName: "nl",
			token:    &token{value: "\n"},
		}, {
			typeName: "int",
			token:    &token{value: "5"},
		}},
	}, {
		msg:  "string",
		code: "\"foo\"",
		nodes: []*node{{
			typeName: "string",
			token:    &token{value: "\"foo\""},
		}},
	}, {
		msg:  "bool",
		code: "true false",
		nodes: []*node{{
			typeName: "true",
			token:    &token{value: "true"},
		}, {
			typeName: "false",
			token:    &token{value: "false"},
		}},
	}, {
		msg:  "symbol",
		code: "foo",
		nodes: []*node{{
			typeName: "symbol",
			token:    &token{value: "foo"},
		}},
	}, {
		msg:  "dynamic symbol",
		code: "symbol(f(a))",
		nodes: []*node{{
			typeName: "dynamic-symbol",
			token:    &token{value: "symbol"},
			nodes: []*node{{
				typeName: "symbol-word",
				token:    &token{value: "symbol"},
			}, {
				typeName: "open-paren",
				token:    &token{value: "("},
			}, {
				typeName: "nls",
				token:    &token{value: "f"},
			}, {
				typeName: "function-call",
				token:    &token{value: "f"},
				nodes: []*node{{
					typeName: "symbol",
					token:    &token{value: "f"},
				}, {
					typeName: "open-paren",
					token:    &token{value: "("},
				}, {
					typeName: "list-sequence",
					token:    &token{value: "a"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "a"},
					}},
				}, {
					typeName: "close-paren",
					token:    &token{value: ")"},
				}},
			}, {
				typeName: "nls",
				token:    &token{value: ")"},
			}, {
				typeName: "close-paren",
				token:    &token{value: ")"},
			}},
		}},
	}, {
		msg:  "empty list",
		code: "[]",
		nodes: []*node{{
			typeName: "list",
			token:    &token{value: "["},
			nodes: []*node{{
				typeName: "open-square",
				token:    &token{value: "["},
			}, {
				typeName: "list-sequence",
				token:    &token{value: "]"},
			}, {
				typeName: "close-square",
				token:    &token{value: "]"},
			}},
		}},
	}, {
		msg:  "list",
		code: "[1, 2, f(a), [3, 4, []]]",
		nodes: []*node{{
			typeName: "list",
			token:    &token{value: "["},
			nodes: []*node{{
				typeName: "open-square",
				token:    &token{value: "["},
			}, {
				typeName: "list-sequence",
				token:    &token{value: "1"},
				nodes: []*node{{
					typeName: "int",
					token:    &token{value: "1"},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "int",
					token:    &token{value: "2"},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "function-call",
					token:    &token{value: "f"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "f"},
					}, {
						typeName: "open-paren",
						token:    &token{value: "("},
					}, {
						typeName: "list-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							typeName: "symbol",
							token:    &token{value: "a"},
						}},
					}, {
						typeName: "close-paren",
						token:    &token{value: ")"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "list",
					token:    &token{value: "["},
					nodes: []*node{{
						typeName: "open-square",
						token:    &token{value: "["},
					}, {
						typeName: "list-sequence",
						token:    &token{value: "3"},
						nodes: []*node{{
							typeName: "int",
							token:    &token{value: "3"},
						}, {
							typeName: "comma",
							token:    &token{value: ","},
						}, {
							typeName: "int",
							token:    &token{value: "4"},
						}, {
							typeName: "comma",
							token:    &token{value: ","},
						}, {
							typeName: "list",
							token:    &token{value: "["},
							nodes: []*node{{
								typeName: "open-square",
								token:    &token{value: "["},
							}, {
								typeName: "list-sequence",
								token:    &token{value: "]"},
							}, {
								typeName: "close-square",
								token:    &token{value: "]"},
							}},
						}},
					}, {
						typeName: "close-square",
						token:    &token{value: "]"},
					}},
				}},
			}, {
				typeName: "close-square",
				token:    &token{value: "]"},
			}},
		}},
	}, {
		msg:  "mutable list",
		code: "~[1, 2, f(a), [3, 4, ~[]]]",
		nodes: []*node{{
			typeName: "mutable-list",
			token:    &token{value: "~"},
			nodes: []*node{{
				typeName: "tilde",
				token:    &token{value: "~"},
			}, {
				typeName: "open-square",
				token:    &token{value: "["},
			}, {
				typeName: "list-sequence",
				token:    &token{value: "1"},
				nodes: []*node{{
					typeName: "int",
					token:    &token{value: "1"},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "int",
					token:    &token{value: "2"},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "function-call",
					token:    &token{value: "f"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "f"},
					}, {
						typeName: "open-paren",
						token:    &token{value: "("},
					}, {
						typeName: "list-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							typeName: "symbol",
							token:    &token{value: "a"},
						}},
					}, {
						typeName: "close-paren",
						token:    &token{value: ")"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "list",
					token:    &token{value: "["},
					nodes: []*node{{
						typeName: "open-square",
						token:    &token{value: "["},
					}, {
						typeName: "list-sequence",
						token:    &token{value: "3"},
						nodes: []*node{{
							typeName: "int",
							token:    &token{value: "3"},
						}, {
							typeName: "comma",
							token:    &token{value: ","},
						}, {
							typeName: "int",
							token:    &token{value: "4"},
						}, {
							typeName: "comma",
							token:    &token{value: ","},
						}, {
							typeName: "mutable-list",
							token:    &token{value: "~"},
							nodes: []*node{{
								typeName: "tilde",
								token:    &token{value: "~"},
							}, {
								typeName: "open-square",
								token:    &token{value: "["},
							}, {
								typeName: "list-sequence",
								token:    &token{value: "]"},
							}, {
								typeName: "close-square",
								token:    &token{value: "]"},
							}},
						}},
					}, {
						typeName: "close-square",
						token:    &token{value: "]"},
					}},
				}},
			}, {
				typeName: "close-square",
				token:    &token{value: "]"},
			}},
		}},
	}, {
		msg:  "empty structure",
		code: "{}",
		nodes: []*node{{
			typeName: "structure",
			token:    &token{value: "{"},
			nodes: []*node{{
				typeName: "open-brace",
				token:    &token{value: "{"},
			}, {
				typeName: "structure-sequence",
				token:    &token{value: "}"},
			}, {
				typeName: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg:  "structure",
		code: "{a: 1, b: 2, ...c, d: {e: 3, f: {}}}",
		nodes: []*node{{
			typeName: "structure",
			token:    &token{value: "{"},
			nodes: []*node{{
				typeName: "open-brace",
				token:    &token{value: "{"},
			}, {
				typeName: "structure-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					typeName: "structure-definition",
					token:    &token{value: "a"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "a"},
					}, {
						typeName: "nls",
						token:    &token{value: ":"},
					}, {
						typeName: "colon",
						token:    &token{value: ":"},
					}, {
						typeName: "nls",
						token:    &token{value: "1"},
					}, {
						typeName: "int",
						token:    &token{value: "1"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "structure-definition",
					token:    &token{value: "b"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "b"},
					}, {
						typeName: "nls",
						token:    &token{value: ":"},
					}, {
						typeName: "colon",
						token:    &token{value: ":"},
					}, {
						typeName: "nls",
						token:    &token{value: "2"},
					}, {
						typeName: "int",
						token:    &token{value: "2"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "spread-expression",
					token:    &token{value: "."},
					nodes: []*node{{
						typeName: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}},
					}, {
						typeName: "symbol",
						token:    &token{value: "c"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "structure-definition",
					token:    &token{value: "d"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "d"},
					}, {
						typeName: "nls",
						token:    &token{value: ":"},
					}, {
						typeName: "colon",
						token:    &token{value: ":"},
					}, {
						typeName: "nls",
						token:    &token{value: "{"},
					}, {
						typeName: "structure",
						token:    &token{value: "{"},
						nodes: []*node{{
							typeName: "open-brace",
							token:    &token{value: "{"},
						}, {
							typeName: "structure-sequence",
							token:    &token{value: "e"},
							nodes: []*node{{
								typeName: "structure-definition",
								token:    &token{value: "e"},
								nodes: []*node{{
									typeName: "symbol",
									token:    &token{value: "e"},
								}, {
									typeName: "nls",
									token:    &token{value: ":"},
								}, {
									typeName: "colon",
									token:    &token{value: ":"},
								}, {
									typeName: "nls",
									token:    &token{value: "3"},
								}, {
									typeName: "int",
									token:    &token{value: "3"},
								}},
							}, {
								typeName: "comma",
								token:    &token{value: ","},
							}, {
								typeName: "structure-definition",
								token:    &token{value: "f"},
								nodes: []*node{{
									typeName: "symbol",
									token:    &token{value: "f"},
								}, {
									typeName: "nls",
									token:    &token{value: ":"},
								}, {
									typeName: "colon",
									token:    &token{value: ":"},
								}, {
									typeName: "nls",
									token:    &token{value: "{"},
								}, {
									typeName: "structure",
									token:    &token{value: "{"},
									nodes: []*node{{
										typeName: "open-brace",
										token:    &token{value: "{"},
									}, {
										typeName: "structure-sequence",
										token:    &token{value: "}"},
									}, {
										typeName: "close-brace",
										token:    &token{value: "}"},
									}},
								}},
							}},
						}, {
							typeName: "close-brace",
							token:    &token{value: "}"},
						}},
					}},
				}},
			}, {
				typeName: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg:  "mutable structure",
		code: "~{a: 1, b: 2, ...c, d: {e: 3, f: ~{}}}",
		nodes: []*node{{
			typeName: "mutable-structure",
			token:    &token{value: "~"},
			nodes: []*node{{
				typeName: "tilde",
				token:    &token{value: "~"},
			}, {
				typeName: "open-brace",
				token:    &token{value: "{"},
			}, {
				typeName: "structure-sequence",
				token:    &token{value: "a"},
				nodes: []*node{{
					typeName: "structure-definition",
					token:    &token{value: "a"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "a"},
					}, {
						typeName: "nls",
						token:    &token{value: ":"},
					}, {
						typeName: "colon",
						token:    &token{value: ":"},
					}, {
						typeName: "nls",
						token:    &token{value: "1"},
					}, {
						typeName: "int",
						token:    &token{value: "1"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "structure-definition",
					token:    &token{value: "b"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "b"},
					}, {
						typeName: "nls",
						token:    &token{value: ":"},
					}, {
						typeName: "colon",
						token:    &token{value: ":"},
					}, {
						typeName: "nls",
						token:    &token{value: "2"},
					}, {
						typeName: "int",
						token:    &token{value: "2"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "spread-expression",
					token:    &token{value: "."},
					nodes: []*node{{
						typeName: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}},
					}, {
						typeName: "symbol",
						token:    &token{value: "c"},
					}},
				}, {
					typeName: "comma",
					token:    &token{value: ","},
				}, {
					typeName: "structure-definition",
					token:    &token{value: "d"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "d"},
					}, {
						typeName: "nls",
						token:    &token{value: ":"},
					}, {
						typeName: "colon",
						token:    &token{value: ":"},
					}, {
						typeName: "nls",
						token:    &token{value: "{"},
					}, {
						typeName: "structure",
						token:    &token{value: "{"},
						nodes: []*node{{
							typeName: "open-brace",
							token:    &token{value: "{"},
						}, {
							typeName: "structure-sequence",
							token:    &token{value: "e"},
							nodes: []*node{{
								typeName: "structure-definition",
								token:    &token{value: "e"},
								nodes: []*node{{
									typeName: "symbol",
									token:    &token{value: "e"},
								}, {
									typeName: "nls",
									token:    &token{value: ":"},
								}, {
									typeName: "colon",
									token:    &token{value: ":"},
								}, {
									typeName: "nls",
									token:    &token{value: "3"},
								}, {
									typeName: "int",
									token:    &token{value: "3"},
								}},
							}, {
								typeName: "comma",
								token:    &token{value: ","},
							}, {
								typeName: "structure-definition",
								token:    &token{value: "f"},
								nodes: []*node{{
									typeName: "symbol",
									token:    &token{value: "f"},
								}, {
									typeName: "nls",
									token:    &token{value: ":"},
								}, {
									typeName: "colon",
									token:    &token{value: ":"},
								}, {
									typeName: "nls",
									token:    &token{value: "~"},
								}, {
									typeName: "mutable-structure",
									token:    &token{value: "~"},
									nodes: []*node{{
										typeName: "tilde",
										token:    &token{value: "~"},
									}, {
										typeName: "open-brace",
										token:    &token{value: "{"},
									}, {
										typeName: "structure-sequence",
										token:    &token{value: "}"},
									}, {
										typeName: "close-brace",
										token:    &token{value: "}"},
									}},
								}},
							}},
						}, {
							typeName: "close-brace",
							token:    &token{value: "}"},
						}},
					}},
				}},
			}, {
				typeName: "close-brace",
				token:    &token{value: "}"},
			}},
		}},
	}, {
		msg:  "symbol query",
		code: "a.b",
		nodes: []*node{{
			typeName: "symbol-query",
			token:    &token{value: "a"},
			nodes: []*node{{
				typeName: "symbol",
				token:    &token{value: "a"},
			}, {
				typeName: "dot",
				token:    &token{value: "."},
			}, {
				typeName: "symbol",
				token:    &token{value: "b"},
			}},
		}},
	}, {
		msg:  "chained symbol query",
		code: "a.b.c",
		nodes: []*node{{
			typeName: "symbol-query",
			token:    &token{value: "a"},
			nodes: []*node{{
				typeName: "symbol-query",
				token:    &token{value: "a"},
				nodes: []*node{{
					typeName: "symbol",
					token:    &token{value: "a"},
				}, {
					typeName: "dot",
					token:    &token{value: "."},
				}, {
					typeName: "symbol",
					token:    &token{value: "b"},
				}},
			}, {
				typeName: "dot",
				token:    &token{value: "."},
			}, {
				typeName: "symbol",
				token:    &token{value: "c"},
			}},
		}},
	}, {
		msg:  "void function",
		code: "fn () {;}",
		nodes: []*node{{
			typeName: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				typeName: "fn-word",
				token:    &token{value: "fn"},
			}, {
				typeName: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					typeName: "open-paren",
					token:    &token{value: "("},
				}, {
					typeName: "static-symbol-sequence",
					token:    &token{value: ")"},
				}, {
					typeName: "nls",
					token:    &token{value: ")"},
				}, {
					typeName: "close-paren",
					token:    &token{value: ")"},
				}, {
					typeName: "nls",
					token:    &token{value: "{"},
				}, {
					typeName: "function-body",
					token:    &token{value: "{"},
					nodes: []*node{{
						typeName: "open-brace",
						token:    &token{value: "{"},
					}, {
						typeName: "statement-sequence",
						token:    &token{value: ";"},
						nodes: []*node{{
							typeName: "semicolon",
							token:    &token{value: ";"},
						}},
					}, {
						typeName: "close-brace",
						token:    &token{value: "}"},
					}},
				}},
			}},
		}},
	}, {
		msg:  "identity",
		code: "fn (x) x",
		nodes: []*node{{
			typeName: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				typeName: "fn-word",
				token:    &token{value: "fn"},
			}, {
				typeName: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					typeName: "open-paren",
					token:    &token{value: "("},
				}, {
					typeName: "static-symbol-sequence",
					token:    &token{value: "x"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "x"},
					}},
				}, {
					typeName: "nls",
					token:    &token{value: ")"},
				}, {
					typeName: "close-paren",
					token:    &token{value: ")"},
				}, {
					typeName: "nls",
					token:    &token{value: "x"},
				}, {
					typeName: "symbol",
					token:    &token{value: "x"},
				}},
			}},
		}},
	}, {
		msg:  "list as a function",
		code: "fn (...x) x",
		nodes: []*node{{
			typeName: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				typeName: "fn-word",
				token:    &token{value: "fn"},
			}, {
				typeName: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					typeName: "open-paren",
					token:    &token{value: "("},
				}, {
					typeName: "static-symbol-sequence",
					token:    &token{value: "."},
				}, {
					typeName: "collect-symbol",
					token:    &token{value: "."},
					nodes: []*node{{
						typeName: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}},
					}, {
						typeName: "symbol",
						token:    &token{value: "x"},
					}},
				}, {
					typeName: "nls",
					token:    &token{value: ")"},
				}, {
					typeName: "close-paren",
					token:    &token{value: ")"},
				}, {
					typeName: "nls",
					token:    &token{value: "x"},
				}, {
					typeName: "symbol",
					token:    &token{value: "x"},
				}},
			}},
		}},
	}, {
		msg:  "function with sequence",
		code: "fn (a, b, ...c) { a(b); c }",
		nodes: []*node{{
			typeName: "function",
			token:    &token{value: "fn"},
			nodes: []*node{{
				typeName: "fn-word",
				token:    &token{value: "fn"},
			}, {
				typeName: "function-fact",
				token:    &token{value: "("},
				nodes: []*node{{
					typeName: "open-paren",
					token:    &token{value: "("},
				}, {
					typeName: "static-symbol-sequence",
					token:    &token{value: "a"},
					nodes: []*node{{
						typeName: "symbol",
						token:    &token{value: "a"},
					}, {
						typeName: "comma",
						token:    &token{value: ","},
					}, {
						typeName: "symbol",
						token:    &token{value: "b"},
					}, {
						typeName: "comma",
						token:    &token{value: ","},
					}},
				}, {
					typeName: "collect-symbol",
					token:    &token{value: "."},
					nodes: []*node{{
						typeName: "spread",
						token:    &token{value: "."},
						nodes: []*node{{
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}, {
							typeName: "dot",
							token:    &token{value: "."},
						}},
					}, {
						typeName: "symbol",
						token:    &token{value: "c"},
					}},
				}, {
					typeName: "nls",
					token:    &token{value: ")"},
				}, {
					typeName: "close-paren",
					token:    &token{value: ")"},
				}, {
					typeName: "nls",
					token:    &token{value: "{"},
				}, {
					typeName: "function-body",
					token:    &token{value: "{"},
					nodes: []*node{{
						typeName: "open-brace",
						token:    &token{value: "{"},
					}, {
						typeName: "statement-sequence",
						token:    &token{value: "a"},
						nodes: []*node{{
							typeName: "function-call",
							token:    &token{value: "a"},
							nodes: []*node{{
								typeName: "symbol",
								token:    &token{value: "a"},
							}, {
								typeName: "open-paren",
								token:    &token{value: "("},
							}, {
								typeName: "list-sequence",
								token:    &token{value: "b"},
								nodes: []*node{{
									typeName: "symbol",
									token:    &token{value: "b"},
								}},
							}, {
								typeName: "close-paren",
								token:    &token{value: ")"},
							}},
						}, {
							typeName: "semicolon",
							token:    &token{value: ";"},
						}, {
							typeName: "symbol",
							token:    &token{value: "c"},
						}},
					}, {
						typeName: "close-brace",
						token:    &token{value: "}"},
					}},
				}},
			}},
		}},
	}} {
		t.Run(ti.msg, func(t *testing.T) {
			b := bytes.NewBufferString(ti.code)
			r := newTokenReader(b, "<test>")

			n, err := s.parse(r)
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

			if n.typeName != "statement-sequence" {
				t.Error("invalid root node type", n.typeName, "statement-sequence")
				return
			}

			if len(n.nodes) != len(ti.nodes) {
				t.Error("invalid number of nodes", len(n.nodes), len(ti.nodes))
				return
			}

			if len(n.nodes) == 0 && n.token != eofToken || len(n.nodes) > 0 && n.token != n.nodes[0].token {
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
