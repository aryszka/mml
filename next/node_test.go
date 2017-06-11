package next

import "testing"

func checkNodes(t *testing.T, left, right []*Node) {
	if len(left) != len(right) {
		t.Error("length doesn't match", len(left), len(right))
		return
	}

	for len(left) > 0 {
		checkNode(t, left[0], right[0])
		if t.Failed() {
			return
		}

		left, right = left[1:], right[1:]
	}
}

func checkNode(t *testing.T, left, right *Node) {
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

	if left.from != right.from {
		t.Error("from doesn't match", left.Name, left.from, right.from)
		return
	}

	if left.to != right.to {
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

	checkNodes(t, left.Nodes, right.Nodes)
}
