package next

import "testing"

func checkNodesPosition(t *testing.T, left, right []*Node, position bool) {
	if len(left) != len(right) {
		t.Error("length doesn't match", len(left), len(right))
		return
	}

	for len(left) > 0 {
		checkNodePosition(t, left[0], right[0], position)
		if t.Failed() {
			return
		}

		left, right = left[1:], right[1:]
	}
}

func checkNodePosition(t *testing.T, left, right *Node, position bool) {
	if (left == nil) != (right == nil) {
		t.Error("nil reference doesn't match", left == nil, right == nil)
		return
	}

	if left == nil {
		return
	}

	if left.Name != right.Name {
		t.Error("name doesn't match", left.Name, right.Name)
		return
	}

	if position && left.from != right.from {
		t.Error("from doesn't match", left.Name, left.from, right.from)
		return
	}

	if position && left.to != right.to {
		t.Error("to doesn't match", left.Name, left.to, right.to)
		return
	}

	if len(left.Nodes) != len(right.Nodes) {
		t.Error("length doesn't match", left.Name, len(left.Nodes), len(right.Nodes))
		t.Log(left)
		t.Log(right)
		for {
			if len(left.Nodes) > 0 {
				t.Log("<", left.Nodes[0])
				left.Nodes = left.Nodes[1:]
			}

			if len(right.Nodes) > 0 {
				t.Log(">", right.Nodes[0])
				right.Nodes = right.Nodes[1:]
			}

			if len(left.Nodes) == 0 && len(right.Nodes) == 0 {
				break
			}
		}
		return
	}

	checkNodesPosition(t, left.Nodes, right.Nodes, position)
}

func checkNodes(t *testing.T, left, right []*Node) {
	checkNodesPosition(t, left, right, true)
}

func checkNode(t *testing.T, left, right *Node) {
	checkNodePosition(t, left, right, true)
}

func checkNodesIgnorePosition(t *testing.T, left, right []*Node) {
	checkNodesPosition(t, left, right, false)
}

func checkNodeIgnorePosition(t *testing.T, left, right *Node) {
	checkNodePosition(t, left, right, false)
}
