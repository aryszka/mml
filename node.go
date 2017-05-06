package mml

import (
	"fmt"
	"strings"
)

type nodeType uint64

type typeList []nodeType

type node struct {
	typ    nodeType
	name   string
	token  *token
	nodes  []*node
	tokens []*token
}

func (l typeList) contains(t nodeType) bool {
	for _, ti := range l {
		if ti == t {
			return true
		}
	}

	return false
}

func (n *node) append(na *node) {
	n.nodes = append(n.nodes, na)
	n.tokens = append(n.tokens, na.tokens...)
	if len(n.nodes) == 1 && len(n.tokens) > 0 {
		n.token = n.tokens[0]
	}
}

func (n *node) String() string {
	nc := make([]string, len(n.nodes))
	for i, ni := range n.nodes {
		nc[i] = ni.String()
	}

	return fmt.Sprintf("{%s:%v:[%s]}", n.name, n.token, strings.Join(nc, ", "))
}
