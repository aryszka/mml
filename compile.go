package mml

// import (
// 	"fmt"
// 	"io"
// )
//
// func compileInt(w io.Writer, n node) error {
// 	_, err := fmt.Fprintf(w, "mml.SysIntToInt(%s)", n.token.value)
// 	return err
// }
//
// func compileString(w io.Writer, n node) error {
// 	_, err := fmt.Fprintf(w, "mml.SysStringToString(\"%s\")", n.token.value)
// 	return err
// }
//
// func compileChannel(w io.Writer, n node) error {
// 	_, err := fmt.Fprint(w, "mml.MakeChannel()")
// 	return err
// }
//
// func compileSymbol(w io.Writer, n node) error {
// 	_, err := fmt.Fprint(w, n.token.value)
// 	return err
// }
//
// func compileDynamicSymbol(w io.Writer, n node) error {
// 	if _, err := fmt.Fprint(w, "mml.SymbolFromValue("); err != nil {
// 		return err
// 	}
//
// 	if err := compile(w, n.nodes[0]); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprint(w, ")"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func compileLookup(w io.Writer, n node) error {
// 	// if _, err := fmt.Fprint(w, "mml.Lookup(env, "); err != nil {
// 	// 	return err
// 	// }
//
// 	var c func(io.Writer, node) error
// 	switch n.typ {
// 	case symbolNode:
// 		c = compileSymbol
// 	// case dynamicSymbolNode:
// 	// 	c = compileDynamicSymbol
// 	}
//
// 	if err := c(w, n); err != nil {
// 		return err
// 	}
//
// 	// if _, err := fmt.Fprint(w, ")"); err != nil {
// 	// 	return err
// 	// }
//
// 	return nil
// }
//
// func compileBoolean(w io.Writer, n node) error {
// 	s := "mml.True"
// 	if n.token.value == "false" {
// 		s = "mml.False"
// 	}
//
// 	_, err := fmt.Fprint(w, s)
// 	return err
// }
//
// func compileListVariant(w io.Writer, n node, variant string) error {
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
// func compileList(w io.Writer, n node) error {
// 	return compileListVariant(w, n, "List")
// }
//
// func compileMutableList(w io.Writer, n node) error {
// 	return compileListVariant(w, n, "MutableList")
// }
//
// func compileSymbolExpression(w io.Writer, n node) error {
// 	switch n.typ {
// 	case symbolNode:
// 		return compileSymbol(w, n)
// 	case stringNode:
// 		return compileString(w, n)
// 	case dynamicSymbolNode:
// 		return compileDynamicSymbol(w, n)
// 	default:
// 		return fmt.Errorf("not implemented: %d, %d, %v, %s", n.token.line, n.token.column, n.typ, n.token.value)
// 	}
// }
//
// func compileStructureVariant(w io.Writer, n node, variant string) error {
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
// 		if err := compileSymbolExpression(w, ni.nodes[0]); err != nil {
// 			return err
// 		}
//
// 		if _, err := fmt.Fprint(w, ","); err != nil {
// 			return err
// 		}
//
// 		if err := compile(w, ni.nodes[1]); err != nil {
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
// func compileStructure(w io.Writer, n node) error {
// 	return compileStructureVariant(w, n, "Structure")
// }
//
// func compileMutableStructure(w io.Writer, n node) error {
// 	return compileStructureVariant(w, n, "MutableStructure")
// }
//
// func andToSwitch(n node) node {
// 	if len(n.nodes) == 0 {
// 		return node{typ: trueNode, token: n.token}
// 	}
//
// 	first, rest := n.nodes[0], n.nodes[1:]
// 	return node{
// 		typ:   switchConditionalNode,
// 		token: n.token,
// 		nodes: []node{{
// 			typ:   switchClauseNode,
// 			token: first.token,
// 			nodes: []node{
// 				first,
// 				andToSwitch(node{nodes: rest}),
// 			},
// 		}, {
// 			typ:   defaultClauseNode,
// 			token: n.token,
// 			nodes: []node{
// 				{typ: falseNode, token: n.token},
// 			},
// 		}},
// 	}
// }
//
// func compileAnd(w io.Writer, n node) error {
// 	return compile(w, andToSwitch(n))
// }
//
// func orToSwitch(n node) node {
// 	if len(n.nodes) == 0 {
// 		return node{typ: falseNode, token: n.token}
// 	}
//
// 	first, rest := n.nodes[0], n.nodes[1:]
// 	return node{
// 		typ:   switchConditionalNode,
// 		token: n.token,
// 		nodes: []node{{
// 			typ:   switchClauseNode,
// 			token: first.token,
// 			nodes: []node{
// 				first,
// 				{typ: trueNode, token: n.token},
// 			},
// 		}, {
// 			typ:   defaultClauseNode,
// 			token: n.token,
// 			nodes: []node{
// 				orToSwitch(node{nodes: rest}),
// 			},
// 		}},
// 	}
// }
//
// func compileOr(w io.Writer, n node) error {
// 	return compile(w, orToSwitch(n))
// }
//
// func compileStatementList(w io.Writer, sep string, ret bool, n []node) error {
// 	for _, ni := range n[:len(n)-1] {
// 		if err := compile(w, ni); err != nil {
// 			return err
// 		}
//
// 		if _, err := fmt.Fprint(w, sep); err != nil {
// 			return err
// 		}
//
// 	}
//
// 	if ret {
// 		if _, err := fmt.Fprint(w, "return "); err != nil {
// 			return err
// 		}
// 	}
//
// 	if err := compile(w, n[len(n)-1]); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func compileSequence(w io.Writer, n []node) error {
// 	return compileStatementList(w, ";", true, n)
// }
//
// func compileFunction(w io.Writer, n node) error {
// 	valueIndex := len(n.nodes) - 1
// 	value := n.nodes[valueIndex]
//
// 	args := n.nodes[:valueIndex]
// 	fixed := args
// 	fixedCount := len(fixed)
//
// 	var variadic bool
// 	if len(fixed) > 0 && fixed[fixedCount-1].typ == collectArgumentNode {
// 		variadic = true
// 		fixed = fixed[:fixedCount-1]
// 		fixedCount--
// 	}
//
// 	if _, err := fmt.Fprint(w, "mml.NewCompiled("); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprintf(
// 		w,
// 		"%d, %t, func(a []*mml.Val) *mml.Val {",
// 		fixedCount,
// 		variadic,
// 	); err != nil {
// 		return err
// 	}
//
// 	for i, ai := range args {
// 		// if i > 0 {
// 		// 	if _, err := fmt.Fprintf(w, ";"); err != nil {
// 		// 		return err
// 		// 	}
// 		// }
//
// 		// var c func(io.Writer, node) error
// 		// switch ai.typ {
// 		// case stringNode:
// 		// 	c = compileString
// 		// case symbolNode:
// 		// 	c = compileSymbol
// 		// }
//
// 		// if err := c(w, ai); err != nil {
// 		// 	return err
// 		// }
//
// 		if _, err := fmt.Fprintf(w, "%s := a[%d];", ai.token.value, i); err != nil {
// 			return err
// 		}
// 	}
//
// 	if value.typ == statementSequenceNode {
// 		if err := compileSequence(w, value.nodes); err != nil {
// 			return err
// 		}
// 	} else {
// 		if err := compile(w, value); err != nil {
// 			return err
// 		}
// 	}
//
// 	if _, err := fmt.Fprint(w, "})"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func compileQuery(w io.Writer, n node) error {
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
//
// func compileValueList(w io.Writer, n []node) error {
// 	return compileStatementList(w, ",", false, n)
// }
//
// func compileFunctionCall(w io.Writer, n node) error {
// 	if _, err := fmt.Fprint(w, "mml.ApplySys("); err != nil {
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
// 	if err := compileValueList(w, n.nodes[1:]); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprint(w, ")"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func compileSwitch(w io.Writer, n node) error {
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
//
// func compileDefinition(w io.Writer, n node) error {
// 	switch n.nodes[0].typ {
// 	case symbolNode:
// 		if err := compileSymbol(w, n.nodes[0]); err != nil {
// 			return err
// 		}
// 	// case stringNode:
// 	// 	if err := compileString(w, n.nodes[0]); err != nil {
// 	// 		return err
// 	// 	}
// 	}
//
// 	// if _, err := fmt.Fprintf(w, ","); err != nil {
// 	// 	return err
// 	// }
//
// 	if _, err := fmt.Fprint(w, " := "); err != nil {
// 		return err
// 	}
//
// 	if err := compile(w, n.nodes[1]); err != nil {
// 		return err
// 	}
//
// 	// if _, err := fmt.Fprintf(w, ")"); err != nil {
// 	// 	return err
// 	// }
//
// 	return nil
// }
//
// func compile(w io.Writer, n node) error {
// 	switch n.typ {
// 	case intNode:
// 		return compileInt(w, n)
// 	case stringNode:
// 		return compileString(w, n)
// 	case channelNode:
// 		return compileChannel(w, n)
// 	case symbolNode, dynamicSymbolNode:
// 		return compileLookup(w, n)
// 	case trueNode, falseNode:
// 		return compileBoolean(w, n)
// 	case listNode:
// 		return compileList(w, n)
// 	case mutableListNode:
// 		return compileMutableList(w, n)
// 	case structureNode:
// 		return compileStructure(w, n)
// 	case mutableStructureNode:
// 		return compileMutableStructure(w, n)
// 	case andExpressionNode:
// 		return compileAnd(w, n)
// 	case orExpressionNode:
// 		return compileOr(w, n)
// 	case functionNode:
// 		return compileFunction(w, n)
// 	case symbolQueryNode, expressionQueryNode:
// 		return compileQuery(w, n)
// 	case functionCallNode:
// 		return compileFunctionCall(w, n)
// 	case switchConditionalNode:
// 		return compileSwitch(w, n)
// 	case definitionNode:
// 		return compileDefinition(w, n)
// 	case statementSequenceNode:
// 		return compileStatementList(w, ";", false, n.nodes)
// 	default:
// 		return fmt.Errorf("not implemented: %d, %d, %v, %s", n.token.line, n.token.column, n.typ, n.token.value)
// 	}
// }
//
// func compileHead(w io.Writer) error {
// 	if _, err := fmt.Fprintln(w, "package main"); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprintln(w, "import \"github.com/aryszka/mml\""); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprintln(w, "func main() {"); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprintln(w, "sys := mml.Sys"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func Compile(in io.Reader, out io.Writer) error {
// 	if err := compileHead(out); err != nil {
// 		return err
// 	}
//
// 	n, err := parseInput(in)
// 	if err != nil {
// 		return err
// 	}
//
// 	if err := compile(out, n); err != nil {
// 		return err
// 	}
//
// 	if _, err := fmt.Fprintln(out, "}"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
