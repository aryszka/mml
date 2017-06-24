package next

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"testing"
)

type jsonValueType int

const (
	jsonNone jsonValueType = iota
	jsonTrue
	jsonFalse
	jsonNull
	jsonString
	jsonNumber
	jsonObject
	jsonArray
)

const (
	maxStringLength  = 64
	meanStringLength = 18
	maxKeyLength     = 24
	meanKeyLength    = 6
	maxObjectLength  = 12
	meanObjectLength = 6
	maxArrayLength   = 64
	meanArrayLength  = 8
)

func randomLength(max, mean int) int {
	return int(rand.NormFloat64()*float64(max)/math.MaxFloat64 + float64(mean))
}

func generateString(max, mean int) string {
	l := randomLength(max, mean)
	b := make([]byte, l)
	for i := range b {
		b[i] = byte(rand.Intn(int('z')-int('a')+1)) + 'a'
	}

	return string(b)
}

func generateJSONString() string {
	return generateString(maxStringLength, meanStringLength)
}

func generateJSONNumber() interface{} {
	if rand.Intn(2) == 1 {
		return rand.NormFloat64()
	}

	n := rand.Int()
	if rand.Intn(2) == 0 {
		return n
	}

	return -n
}

func generateKey() string {
	return generateString(maxKeyLength, meanKeyLength)
}

func generateJSONObject(minDepth int) map[string]interface{} {
	l := randomLength(maxObjectLength, meanObjectLength)
	o := make(map[string]interface{})
	for i := 0; i < l; i++ {
		o[generateKey()] = generateJSON(0)
	}

	if minDepth > 0 {
		o[generateKey()] = generateJSON(minDepth)
	}

	return o
}

func generateJSONArray(minDepth int) []interface{} {
	l := randomLength(maxArrayLength, meanArrayLength)
	a := make([]interface{}, l, l+1)
	for i := 0; i < l; i++ {
		a[i] = generateJSON(0)
	}

	if minDepth > 0 {
		a = append(a, generateJSON(minDepth))
	}

	return a
}

func generateJSONObjectOrArray(minDepth int) interface{} {
	if rand.Intn(2) == 0 {
		return generateJSONObject(minDepth - 1)
	}

	return generateJSONArray(minDepth - 1)
}

func generateJSON(minDepth int) interface{} {
	if minDepth > 0 {
		return generateJSONObjectOrArray(minDepth)
	}

	switch jsonValueType(rand.Intn(int(jsonNumber)) + 1) {
	case jsonTrue:
		return true
	case jsonFalse:
		return false
	case jsonNull:
		return nil
	case jsonString:
		return generateJSONString()
	case jsonNumber:
		return generateJSONNumber()
	default:
		panic("invalid json type")
	}
}

func unqouteJSONString(t string) (string, error) {
	var s string
	err := json.Unmarshal([]byte(t), &s)
	return s, err
}

func parseJSONNumber(t string) (interface{}, error) {
	n := json.Number(t)
	if i, err := n.Int64(); err == nil {
		return int(i), nil
	}

	return n.Float64()
}

func nodeToJSONObject(n *Node) (map[string]interface{}, error) {
	o := make(map[string]interface{})
	for _, ni := range n.Nodes {
		if len(ni.Nodes) != 2 {
			return nil, errors.New("invalid json object")
		}

		key, err := unqouteJSONString(ni.Nodes[0].Text())
		if err != nil {
			return nil, err
		}

		val, err := treeToJSON(ni.Nodes[1])
		if err != nil {
			return nil, err
		}

		o[key] = val
	}

	return o, nil
}

func nodeToJSONArray(n *Node) ([]interface{}, error) {
	a := make([]interface{}, 0, len(n.Nodes))
	for _, ni := range n.Nodes {
		item, err := treeToJSON(ni)
		if err != nil {
			return nil, err
		}

		a = append(a, item)
	}

	return a, nil
}

func treeToJSON(n *Node) (interface{}, error) {
	switch n.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	case "string":
		return unqouteJSONString(n.Text())
	case "number":
		return parseJSONNumber(n.Text())
	case "object":
		return nodeToJSONObject(n)
	case "array":
		return nodeToJSONArray(n)
	default:
		return nil, fmt.Errorf("invalid json node name: %s", n.Name)
	}
}

func checkJSON(t *testing.T, got, expected interface{}) {
	if expected == nil {
		if got != nil {
			t.Error("expected nil", got)
		}

		return
	}

	switch v := expected.(type) {
	case bool:
		if v != got.(bool) {
			t.Error("expected bool", got)
		}
	case string:
		if v != got.(string) {
			t.Error("expected string", got)
		}
	case int:
		if v != got.(int) {
			t.Error("expected int", got)
		}
	case float64:
		if v != got.(float64) {
			t.Error("expected float64", got)
		}
	case map[string]interface{}:
		o, ok := got.(map[string]interface{})
		if !ok {
			t.Error("expected object", got)
			return
		}

		if len(v) != len(o) {
			t.Error("invalid object length, expected: %d, got: %d", len(v), len(o))
			return
		}

		for key, val := range v {
			gotVal, ok := o[key]
			if !ok {
				t.Error("expected key not found: %s", key)
				return
			}

			checkJSON(t, gotVal, val)
			if t.Failed() {
				return
			}
		}
	case []interface{}:
		a, ok := got.([]interface{})
		if !ok {
			t.Error("expected array", got)
		}

		if len(v) != len(a) {
			t.Error("invalid array length, expected: %d, got: %d", len(v), len(a))
			return
		}

		for i := range v {
			checkJSON(t, a[i], v[i])
			if t.Failed() {
				return
			}
		}
	default:
		t.Error("unexpected parsed type", v)
	}
}

func jsonTreeToJSON(n *Node) (interface{}, error) {
	if n.Name != "json" {
		return nil, fmt.Errorf("invalid root node name: %s", n.Name)
	}

	if len(n.Nodes) != 1 {
		return nil, fmt.Errorf("invalid root node length: %d", len(n.Nodes))
	}

	return treeToJSON(n.Nodes[0])
}

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
	}, {
		msg:  "bugfix, 100",
		text: "100",
		node: &Node{
			Name: "json",
			Nodes: []*Node{{
				Name: "number",
			}},
		},
		ignorePosition: true,
	}})
}

func TestRandomJSON(t *testing.T) {
	j := generateJSON(48)
	b, err := json.Marshal(j)
	if err != nil {
		t.Error(err)
		return
	}

	buf := bytes.NewBuffer(b)

	s, err := testSyntax("json.p", 0)
	if err != nil {
		t.Error(err)
		return
	}

	testParse := func(t *testing.T, buf io.Reader) {
		n, err := s.Parse(buf)
		if err != nil {
			t.Error(err)
			return
		}

		jback, err := jsonTreeToJSON(n)
		if err != nil {
			t.Error(err)
			return
		}

		checkJSON(t, jback, j)
	}

	t.Run("unindented", func(t *testing.T) {
		testParse(t, buf)
	})

	indented := bytes.NewBuffer(nil)
	if err := json.Indent(indented, b, "", "    "); err != nil {
		t.Error(err)
		return
	}

	t.Run("indented", func(t *testing.T) {
		testParse(t, indented)
	})

	indentedTabs := bytes.NewBuffer(nil)
	if err := json.Indent(indentedTabs, b, "", "\t"); err != nil {
		t.Error(err)
		return
	}

	t.Run("indented with tabs", func(t *testing.T) {
		testParse(t, indentedTabs)
	})
}
