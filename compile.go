package mml

import (
	"bytes"
	"fmt"
	"io"
)

var symbols = make(map[string]string)

var builtin = map[string]*Val{
	"sum": NewCompiled(0, true, Sum),
}

func resolveSymbol(s string) string {
	if rs, ok := symbols[s]; ok {
		return rs
	}

	rs := fmt.Sprintf("s_%d", len(symbols))
	symbols[s] = rs
	return rs
}

func compileInt(w io.Writer, n *node) error {
	_, err := fmt.Fprintf(w, "mml.SysIntToInt(%s)", n.token.value)
	return err
}

func compileString(w io.Writer, n *node) error {
	_, err := fmt.Fprintf(w, "mml.SysStringToString(%s)", n.token.value)
	return err
}

// func compileChannel(w io.Writer, n *node) error {
// 	_, err := fmt.Fprint(w, "mml.MakeChannel()")
// 	return err
// }

func compileSymbol(w io.Writer, n *node) error {
	_, err := fmt.Fprint(w, resolveSymbol(n.token.value))
	return err
}

func compileDynamicSymbol(w io.Writer, n *node) error {
	if _, err := fmt.Fprint(w, "mml.SymbolFromValue("); err != nil {
		return err
	}

	if err := compile(w, n.nodes[0]); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, ")"); err != nil {
		return err
	}

	return nil
}

func compileLookup(w io.Writer, n *node) error {
	// if _, err := fmt.Fprint(w, "mml.Lookup(env, "); err != nil {
	// 	return err
	// }

	var c func(io.Writer, *node) error
	switch n.typ {
	case "symbol":
		c = compileSymbol
		// case dynamicSymbolNode:
		// 	c = compileDynamicSymbol
	}

	if err := c(w, n); err != nil {
		return err
	}

	// if _, err := fmt.Fprint(w, ")"); err != nil {
	// 	return err
	// }

	return nil
}

// func compileBoolean(w io.Writer, n *node) error {
// 	s := "mml.True"
// 	if n.token.value == "false" {
// 		s = "mml.False"
// 	}
//
// 	_, err := fmt.Fprint(w, s)
// 	return err
// }
//
// func compileListVariant(w io.Writer, n *node, variant string) error {
// 	if _, err := fmt.Fprintf(w, "mml.%s(", variant); err != nil {
// 		return err
// 	}
//
// 	for i, ni := range n.nodes {
// 		if i > 0 {
// 			if _, err := fmt.Fprint(w, ","); err != nil {
// 				return err
// 			}
// 		}
//
// 		if err := compile(w, ni); err != nil {
// 			return err
// 		}
// 	}
//
// 	if _, err := fmt.Fprint(w, ")"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func compileList(w io.Writer, n *node) error {
// 	return compileListVariant(w, n, "List")
// }
//
// func compileMutableList(w io.Writer, n *node) error {
// 	return compileListVariant(w, n, "MutableList")
// }

func compileSymbolLiteral(w io.Writer, n *node) error {
	_, err := fmt.Fprintf(w, "\"%s\"", n.token.value)
	return err
}

func compileSymbolExpression(w io.Writer, n *node) error {
	switch n.typ {
	case "symbol":
		return compileSymbolLiteral(w, n)
	case "string":
		return compileString(w, n)
	case "dynamic-symbol":
		return compileDynamicSymbol(w, n)
	default:
		return fmt.Errorf("not implemented: %d, %d, %v, %s", n.token.line, n.token.column, n.typ, n.token.value)
	}
}

func compileStructureVariant(w io.Writer, n *node, variant string) error {
	if _, err := fmt.Fprintf(w, "mml.%s(", variant); err != nil {
		return err
	}

	for i, ni := range n.nodes {
		if i > 0 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}

		if err := compileSymbolExpression(w, ni.nodes[0]); err != nil {
			return err
		}

		if _, err := fmt.Fprint(w, ","); err != nil {
			return err
		}

		if err := compile(w, ni.nodes[1]); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, ")"); err != nil {
		return err
	}

	return nil
}

// func compileStructure(w io.Writer, n *node) error {
// 	return compileStructureVariant(w, n, "Structure")
// }

func compileMutableStructure(w io.Writer, n *node) error {
	return compileStructureVariant(w, n, "MutableStructure")
}

// func andToSwitch(n *node) *node {
// 	if len(n.nodes) == 0 {
// 		return &node{typ: trueNode, token: n.token}
// 	}
//
// 	first, rest := n.nodes[0], n.nodes[1:]
// 	return &node{
// 		typ:   switchConditionalNode,
// 		token: n.token,
// 		nodes: []*node{{
// 			typ:   switchClauseNode,
// 			token: first.token,
// 			nodes: []*node{
// 				first,
// 				andToSwitch(*node{nodes: rest}),
// 			},
// 		}, {
// 			typ:   defaultClauseNode,
// 			token: n.token,
// 			nodes: []*node{
// 				{typ: falseNode, token: n.token},
// 			},
// 		}},
// 	}
// }
//
// func compileAnd(w io.Writer, n *node) error {
// 	return compile(w, andToSwitch(n))
// }
//
// func orToSwitch(n *node) *node {
// 	if len(n.nodes) == 0 {
// 		return &node{typ: falseNode, token: n.token}
// 	}
//
// 	first, rest := n.nodes[0], n.nodes[1:]
// 	return &node{
// 		typ:   switchConditionalNode,
// 		token: n.token,
// 		nodes: []*node{{
// 			typ:   switchClauseNode,
// 			token: first.token,
// 			nodes: []*node{
// 				first,
// 				{typ: trueNode, token: n.token},
// 			},
// 		}, {
// 			typ:   defaultClauseNode,
// 			token: n.token,
// 			nodes: []*node{
// 				orToSwitch(&node{nodes: rest}),
// 			},
// 		}},
// 	}
// }
//
// func compileOr(w io.Writer, n *node) error {
// 	return compile(w, orToSwitch(n))
// }

func compileStatementList(w io.Writer, sep string, ret bool, n []*node) error {
	if len(n) == 0 {
		return nil
	}

	for _, ni := range n[:len(n)-1] {
		if err := compile(w, ni); err != nil {
			return err
		}

		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}

	}

	if ret {
		if _, err := fmt.Fprint(w, "return "); err != nil {
			return err
		}
	}

	if err := compile(w, n[len(n)-1]); err != nil {
		return err
	}

	return nil
}

func compileSequence(w io.Writer, n []*node) error {
	return compileStatementList(w, ";", false, n)
}

func compileStaticSymbol(w io.Writer, n *node) error {
	_, err := fmt.Fprint(w, resolveSymbol(n.token.value))
	return err
}

func compileFunction(w io.Writer, n *node) error {
	valueIndex := len(n.nodes) - 1
	value := n.nodes[valueIndex]

	args := n.nodes[:valueIndex]
	fixedCount := len(args)

	var variadic bool
	if len(args) > 0 && args[fixedCount-1].typ == "collect-symbol" {
		variadic = true
		fixedCount--
	}

	if _, err := fmt.Fprint(w, "mml.NewCompiled("); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(
		w,
		"%d, %t, func(a []*mml.Val) *mml.Val {",
		fixedCount,
		variadic,
	); err != nil {
		return err
	}

	fmt.Fprint(w, "var (\n")
	for i, ai := range args[:fixedCount] {
		// if i > 0 {
		// 	if _, err := fmt.Fprintf(w, ";"); err != nil {
		// 		return err
		// 	}
		// }

		// var c func(io.Writer, *node) error
		// switch ai.typ {
		// case stringNode:
		// 	c = compileString
		// case symbolNode:
		// 	c = compileSymbol
		// }

		// if err := c(w, ai); err != nil {
		// 	return err
		// }

		if err := compileStaticSymbol(w, ai); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, " = a[%d]\n", i); err != nil {
			return err
		}
	}

	if variadic {
		if err := compileStaticSymbol(w, args[len(args)-1].nodes[0]); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, " = mml.ListFromSysSlice(a[%d:])\n", fixedCount); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, ")\n"); err != nil {
		return err
	}

	if value.typ == "statement-sequence" {
		if len(value.nodes) == 0 {
			if _, err := fmt.Fprintln(w, "return mml.Void})"); err != nil {
				return err
			}

			return nil
		}

		if err := compileSequence(w, value.nodes[:len(value.nodes)-1]); err != nil {
			return err
		}

		if _, err := fmt.Fprint(w, ";return "); err != nil {
			return err
		}

		if err := compile(w, value.nodes[len(value.nodes)-1]); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprint(w, "return "); err != nil {
			return err
		}

		if err := compile(w, value); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "})"); err != nil {
		return err
	}

	return nil
}

// func compileQuery(w io.Writer, n *node) error {
// 	if _, err := fmt.Fprint(w, "mml.Query("); err != nil {
// 		return err
// 	}
//
// 	if err := compile(w, n.nodes[0]); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprint(w, ","); err != nil {
// 		return err
// 	}
//
// 	switch n.nodes[1].typ {
// 	case symbolNode:
// 		if _, err := fmt.Fprint(w, "mml.SysStringToSymbol(\""); err != nil {
// 			return err
// 		}
//
// 		if err := compileSymbol(w, n.nodes[1]); err != nil {
// 			return err
// 		}
//
// 		if _, err := fmt.Fprint(w, "\")"); err != nil {
// 			return err
// 		}
// 	default:
// 		if err := compile(w, n.nodes[1]); err != nil {
// 			return err
// 		}
// 	}
//
// 	if _, err := fmt.Fprint(w, ")"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func compileValueList(w io.Writer, n []*node) error {
// 	return compileStatementList(w, ",", false, n)
// }

// func(a...)
// mml.ApplySys(func, mml.ListToSysSlice(a))
//
// func(a, b...)
// mml.ApplySys(func, append([]*mml.Val{a}, mml.ListToSysSlice(b)...))
//
// func(a, b, c..., d, e, f...)
// mml.ApplySys(f, append(append([]*mml.Val{a, b}, append(mml.ListToSysSlice(c), d, e)...), mml.ListToSysSlice(f)...))
//
// // do the recursion

func compileArgList(w io.Writer, n []*node, forceList bool) (int, error) {
	if len(n) == 0 {
		fmt.Fprint(w, "nil")
		return 0, nil
	}

	last := len(n) - 1

	switch n[last].typ {
	case "spread-expression":
		wp := bytes.NewBuffer(nil)
		tp, err := compileArgList(wp, n[:last], true)
		if err != nil {
			return 0, err
		}

		switch tp {
		case 0:
			if _, err := fmt.Fprint(w, "mml.ListToSysSlice("); err != nil {
				return 0, err
			}

			if err := compile(w, n[last].nodes[0]); err != nil {
				return 0, err
			}

			if _, err := fmt.Fprint(w, ")"); err != nil {
				return 0, err
			}

			return 2, nil
		default:
			if _, err := fmt.Fprint(w, "append("); err != nil {
				return 0, err
			}

			if _, err := io.Copy(w, wp); err != nil {
				return 0, err
			}

			if _, err := fmt.Fprint(w, ", mml.ListToSysSlice("); err != nil {
				return 0, err
			}

			if err := compile(w, n[last].nodes[0]); err != nil {
				return 0, err
			}

			if _, err := fmt.Fprint(w, ")...)"); err != nil {
				return 0, err
			}

			return 2, nil
		}
	default:
		wp := bytes.NewBuffer(nil)
		tp, err := compileArgList(wp, n[:last], false)
		if err != nil {
			return 0, err
		}

		if forceList {
			switch tp {
			case 1:
				if _, err := fmt.Fprint(w, "[]*mml.Val{"); err != nil {
					return 0, err
				}

				if _, err := io.Copy(w, wp); err != nil {
					return 0, err
				}

				if _, err := fmt.Fprint(w, ", "); err != nil {
					return 0, err
				}

				if err := compile(w, n[last]); err != nil {
					return 0, err
				}

				if _, err := fmt.Fprint(w, "}"); err != nil {
					return 0, err
				}

				return 2, err
			case 2:
				if _, err := fmt.Fprint(w, "append("); err != nil {
					return 0, err
				}

				if _, err := io.Copy(w, wp); err != nil {
					return 0, err
				}

				if _, err := fmt.Fprint(w, ", "); err != nil {
					return 0, err
				}

				if err := compile(w, n[last]); err != nil {
					return 0, err
				}

				if _, err := fmt.Fprint(w, ")"); err != nil {
					return 0, err
				}

				return 2, err
			default:
				if _, err := fmt.Fprint(w, "[]*mml.Val{"); err != nil {
					return 0, err
				}

				if err := compile(w, n[last]); err != nil {
					return 0, err
				}

				if _, err := fmt.Fprint(w, "}"); err != nil {
					return 0, err
				}

				return 2, err
			}
		} else {
			switch tp {
			case 0:
				if err := compile(w, n[last]); err != nil {
					return 0, err
				}

				return 1, nil
			default:
				if _, err := io.Copy(w, wp); err != nil {
					return 0, err
				}

				if _, err := fmt.Fprint(w, ", "); err != nil {
					return 0, err
				}

				if err := compile(w, n[last]); err != nil {
					return 0, err
				}

				return tp, nil
			}
		}
	}
}

func compileFunctionCall(w io.Writer, n *node) error {
	if _, err := fmt.Fprint(w, "mml.ApplySys("); err != nil {
		return err
	}

	if err := compile(w, n.nodes[0]); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, ","); err != nil {
		return err
	}

	if _, err := compileArgList(w, n.nodes[1:], true); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, ")"); err != nil {
		return err
	}

	return nil
}

// func compileSwitch(w io.Writer, n *node) error {
// 	if _, err := fmt.Fprint(w, "func() *mml.Val {"); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprintln(w, "switch {"); err != nil {
// 		return err
// 	}
//
// 	for _, ni := range n.nodes {
// 		switch ni.typ {
// 		case switchClauseNode:
// 			if _, err := fmt.Fprintln(w, "\ncase "); err != nil {
// 				return err
// 			}
//
// 			if err := compile(w, ni.nodes[0]); err != nil {
// 				return err
// 			}
//
// 			if _, err := fmt.Fprintln(w, "== mml.True:"); err != nil {
// 				return err
// 			}
//
// 			if err := compileSequence(w, ni.nodes[1:]); err != nil {
// 				return err
// 			}
// 		default:
// 			if _, err := fmt.Fprintln(w, "\ndefault:"); err != nil {
// 				return err
// 			}
//
// 			if err := compileSequence(w, ni.nodes); err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	if _, err := fmt.Fprint(w, "}}()"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func compileValueDefinition(w io.Writer, n *node) error {
	if err := compileStaticSymbol(w, n.nodes[0]); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, " := "); err != nil {
		return err
	}

	if err := compile(w, n.nodes[1]); err != nil {
		return err
	}

	return nil
}

func compile(w io.Writer, n *node) error {
	switch n.typ {
	case "int":
		return compileInt(w, n)
	case "string":
		return compileString(w, n)
	// case channelNode:
	// 	return compileChannel(w, n)
	case "symbol", "dynamic-symbol":
		return compileLookup(w, n)
	// case trueNode, falseNode:
	// 	return compileBoolean(w, n)
	// case listNode:
	// 	return compileList(w, n)
	// case mutableListNode:
	// 	return compileMutableList(w, n)
	// case structureNode:
	// 	return compileStructure(w, n)
	case "mutable-structure":
		return compileMutableStructure(w, n)
	// case andExpressionNode:
	// 	return compileAnd(w, n)
	// case orExpressionNode:
	// 	return compileOr(w, n)
	case "function":
		return compileFunction(w, n)
	// case symbolQueryNode, expressionQueryNode:
	// 	return compileQuery(w, n)
	case "function-call":
		return compileFunctionCall(w, n)
	// case switchConditionalNode:
	// 	return compileSwitch(w, n)
	case "value-definition":
		return compileValueDefinition(w, n)
	case "statement-sequence":
		return compileStatementList(w, ";", false, n.nodes)
	default:
		return fmt.Errorf("not implemented: %d, %d, %v, %s", n.token.line, n.token.column, n.typ, n.token.value)
	}
}

func compileBuiltin(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "var ("); err != nil {
		return err
	}

	for name := range Builtin {
		s := resolveSymbol(name)
		if _, err := fmt.Fprintf(w, "%s = mml.Builtin[\"%s\"]\n", s, name); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, ")"); err != nil {
		return err
	}

	return nil
}

func compileHead(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "package main"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "import \"github.com/aryszka/mml\""); err != nil {
		return err
	}

	if err := compileBuiltin(w); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "func main() {"); err != nil {
		return err
	}

	return nil
}

func Compile(in io.Reader, out io.Writer) error {
	r := newTokenReader(in, "test")

	n, err := parse(traceOff, generatorsByName["document"], r)
	if err != nil {
		return err
	}

	return nil
	if err := compileHead(out); err != nil {
		return err
	}

	if err := compile(out, n); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(out, "}"); err != nil {
		return err
	}

	return nil
}
