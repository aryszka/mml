package next

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/zalando/skipper/eskip"
)

const (
	maxID  = 27
	meanID = 9

	setPathChance = 0.72
	maxPathTags   = 12
	meanPathTags  = 2
	maxPathTag    = 24
	meanPathTag   = 9

	setHostChance = 0.5
	maxHost       = 48
	meanHost      = 24

	setPathRegexpChance = 0.45
	maxPathRegexp       = 36
	meanPathRegexp      = 12

	setMethodChance = 0.1

	setHeadersChance      = 0.3
	maxHeadersLength      = 6
	meanHeadersLength     = 1
	maxHeaderKeyLength    = 18
	meanHeaderKeyLength   = 12
	maxHeaderValueLength  = 48
	meanHeaderValueLength = 6

	setHeaderRegexpChance   = 0.05
	maxHeaderRegexpsLength  = 3
	meanHeaderRegexpsLength = 1
	maxHeaderRegexpLength   = 12
	meanHeaderRegexpLength  = 6

	maxTermNameLength    = 15
	meanTermNameLength   = 6
	maxTermArgsLength    = 6
	meanTermArgsLength   = 1
	floatArgChance       = 0.1
	intArgChance         = 0.3
	maxTermStringLength  = 24
	meanTermStringLength = 6

	maxPredicatesLength  = 4
	meanPredicatesLength = 1

	maxFiltersLength  = 18
	meanFiltersLength = 3

	loopBackendChance  = 0.05
	shuntBackendChance = 0.1
	maxBackend         = 48
	meanBackend        = 15
)

func takeChance(c float64) bool {
	return rand.Float64() < c
}

func generateID() string {
	return generateString(maxID, meanID)
}

func generatePath() string {
	if !takeChance(setPathChance) {
		return ""
	}

	l := randomLength(maxPathTags, meanPathTags)
	p := append(make([]string, 0, l+1), "")
	for i := 0; i < l; i++ {
		p = append(p, generateString(maxPathTag, meanPathTag))
	}

	return strings.Join(p, "/")
}

func generateHostRegexps() []string {
	if !takeChance(setHostChance) {
		return nil
	}

	return []string{generateString(maxHost, meanHost)}
}

func generatePathRegexps() []string {
	if !takeChance(setPathRegexpChance) {
		return nil
	}

	return []string{generateString(maxPathRegexp, meanPathRegexp)}
}

func generateMethod() string {
	if !takeChance(setMethodChance) {
		return ""
	}

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	return methods[rand.Intn(len(methods))]
}

func generateHeaders() map[string]string {
	if !takeChance(setHeadersChance) {
		return nil
	}

	h := make(map[string]string)
	for i := 0; i < randomLength(maxHeadersLength, meanHeadersLength); i++ {
		h[generateString(maxHeaderKeyLength, meanHeaderKeyLength)] =
			generateString(maxHeaderValueLength, meanHeaderValueLength)
	}

	return h
}

func generateHeaderRegexps() map[string][]string {
	if !takeChance(setHeaderRegexpChance) {
		return nil
	}

	h := make(map[string][]string)
	for i := 0; i < randomLength(maxHeaderRegexpsLength, meanHeaderRegexpsLength); i++ {
		k := generateString(maxHeaderKeyLength, meanHeaderKeyLength)
		for i := 0; i < randomLength(maxHeaderRegexpLength, meanHeaderRegexpLength); i++ {
			h[k] = append(h[k], generateString(maxHeaderValueLength, meanHeaderValueLength))
		}
	}

	return h
}

func generateTerm() (string, []interface{}) {
	n := generateString(maxTermNameLength, meanTermNameLength)
	al := randomLength(maxTermArgsLength, meanTermArgsLength)
	a := make([]interface{}, 0, al)
	for i := 0; i < al; i++ {
		at := rand.Float64()
		switch {
		case at < floatArgChance:
			a = append(a, rand.NormFloat64())
		case at < intArgChance:
			a = append(a, rand.Int())
		default:
			a = append(a, generateString(maxTermStringLength, meanTermStringLength))
		}
	}

	return n, a
}

func generatePredicates() []*eskip.Predicate {
	l := randomLength(maxPredicatesLength, meanPredicatesLength)
	p := make([]*eskip.Predicate, 0, l)
	for i := 0; i < l; i++ {
		pi := &eskip.Predicate{}
		pi.Name, pi.Args = generateTerm()
		p = append(p, pi)
	}

	return p
}

func generateFilters() []*eskip.Filter {
	l := randomLength(maxFiltersLength, meanFiltersLength)
	f := make([]*eskip.Filter, 0, l)
	for i := 0; i < l; i++ {
		fi := &eskip.Filter{}
		fi.Name, fi.Args = generateTerm()
		f = append(f, fi)
	}

	return f
}

func generateBackend() (eskip.BackendType, string) {
	t := rand.Float64()
	switch {
	case t < loopBackendChance:
		return eskip.LoopBackend, ""
	case t < loopBackendChance+shuntBackendChance:
		return eskip.ShuntBackend, ""
	default:
		return eskip.NetworkBackend, generateString(maxBackend, meanBackend)
	}
}

func generateRoute() *eskip.Route {
	r := &eskip.Route{}
	r.Id = generateID()
	r.Path = generatePath()
	r.HostRegexps = generateHostRegexps()
	r.PathRegexps = generatePathRegexps()
	r.Method = generateMethod()
	r.Headers = generateHeaders()
	r.HeaderRegexps = generateHeaderRegexps()
	r.Predicates = generatePredicates()
	r.Filters = generateFilters()
	r.BackendType, r.Backend = generateBackend()
	return r
}

func generateEskip(l int) []*eskip.Route {
	r := make([]*eskip.Route, 0, l)
	for i := 0; i < l; i++ {
		r = append(r, generateRoute())
	}

	return r
}

func parseEskipInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 0, 64)
	return int(i), err
}

func parseEskipFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err
}

func unquote(s string, escapedChars string) (string, error) {
	if len(s) < 2 {
		return "", nil
	}

	b := make([]byte, 0, len(s)-2)
	var escaped bool
	for _, bi := range []byte(s[1 : len(s)-1]) {
		if escaped {
			switch bi {
			case 'b':
				bi = '\b'
			case 'f':
				bi = '\f'
			case 'n':
				bi = '\n'
			case 'r':
				bi = '\r'
			case 't':
				bi = '\t'
			case 'v':
				bi = '\v'
			}

			b = append(b, bi)
			escaped = false
			continue
		}

		for _, ec := range []byte(escapedChars) {
			if ec == bi {
				return "", errors.New("invalid quote")
			}
		}

		if bi == '\\' {
			escaped = true
			continue
		}

		b = append(b, bi)
	}

	return string(b), nil
}

func unquoteString(s string) (string, error) {
	return unquote(s, "\"")
}

func unquoteRegexp(s string) (string, error) {
	return unquote(s, "/")
}

func nodeToArg(n *Node) (interface{}, error) {
	switch n.Name {
	case "int":
		return parseEskipInt(n.Text())
	case "float":
		return parseEskipFloat(n.Text())
	case "string":
		return unquoteString(n.Text())
	case "regexp":
		return unquoteRegexp(n.Text())
	default:
		return nil, errors.New("invalid arg")
	}
}

func nodeToTerm(n *Node) (string, []interface{}, error) {
	if len(n.Nodes) < 1 || n.Nodes[0].Name != "symbol" {
		return "", nil, errors.New("invalid term")
	}

	name := n.Nodes[0].Text()

	var args []interface{}
	for _, ni := range n.Nodes[1:] {
		a, err := nodeToArg(ni)
		if err != nil {
			return "", nil, err
		}

		args = append(args, a)
	}

	return name, args, nil
}

func nodeToPredicate(r *eskip.Route, n *Node) error {
	name, args, err := nodeToTerm(n)
	if err != nil {
		return err
	}

	switch name {
	case "Path":
		if len(args) != 1 {
			return errors.New("invalid path predicate")
		}

		p, ok := args[0].(string)
		if !ok {
			return errors.New("invalid path predicate")
		}

		r.Path = p
	case "Host":
		if len(args) != 1 {
			return errors.New("invalid host predicate")
		}

		h, ok := args[0].(string)
		if !ok {
			return errors.New("invalid host predicate")
		}

		r.HostRegexps = append(r.HostRegexps, h)
	case "PathRegexp":
		if len(args) != 1 {
			return errors.New("invalid path regexp predicate")
		}

		p, ok := args[0].(string)
		if !ok {
			return errors.New("invalid path regexp predicate")
		}

		r.PathRegexps = append(r.PathRegexps, p)
	case "Method":
		if len(args) != 1 {
			return errors.New("invalid method predicate")
		}

		m, ok := args[0].(string)
		if !ok {
			return errors.New("invalid method predicate")
		}

		r.Method = m
	case "Header":
		if len(args) != 2 {
			return errors.New("invalid header predicate")
		}

		name, ok := args[0].(string)
		if !ok {
			return errors.New("invalid header predicate")
		}

		value, ok := args[1].(string)
		if !ok {
			return errors.New("invalid header predicate")
		}

		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}

		r.Headers[name] = value
	case "HeaderRegexp":
		if len(args) != 2 {
			return errors.New("invalid header regexp predicate")
		}

		name, ok := args[0].(string)
		if !ok {
			return errors.New("invalid header regexp predicate")
		}

		value, ok := args[1].(string)
		if !ok {
			return errors.New("invalid header regexp predicate")
		}

		if r.HeaderRegexps == nil {
			r.HeaderRegexps = make(map[string][]string)
		}

		r.HeaderRegexps[name] = append(r.HeaderRegexps[name], value)
	default:
		r.Predicates = append(r.Predicates, &eskip.Predicate{Name: name, Args: args})
	}

	return nil
}

func nodeToFilter(n *Node) (*eskip.Filter, error) {
	name, args, err := nodeToTerm(n)
	if err != nil {
		return nil, err
	}

	return &eskip.Filter{Name: name, Args: args}, nil
}

func nodeToBackend(r *eskip.Route, n *Node) error {
	switch n.Name {
	case "string":
		b, err := unquoteString(n.Text())
		if err != nil {
			return err
		}

		r.BackendType = eskip.NetworkBackend
		r.Backend = b
	case "shunt":
		r.BackendType = eskip.ShuntBackend
	case "loopback":
		r.BackendType = eskip.LoopBackend
	default:
		return errors.New("invalid backend type")
	}

	return nil
}

func nodeToEskipDefinition(n *Node) (*eskip.Route, error) {
	ns := n.Nodes
	if len(ns) < 2 || len(ns[1].Nodes) == 0 {
		return nil, fmt.Errorf("invalid definition length: %d", len(ns))
	}

	r := &eskip.Route{}

	if ns[0].Name != "symbol" {
		return nil, errors.New("invalid definition id")
	}

	r.Id, ns = ns[0].Text(), ns[1].Nodes

predicates:
	for i, ni := range ns {
		switch ni.Name {
		case "predicate":
			if err := nodeToPredicate(r, ni); err != nil {
				return nil, err
			}
		case "filter", "string", "shunt", "loopback":
			ns = ns[i:]
			break predicates
		default:
			return nil, errors.New("invalid definition item among predicates")
		}
	}

filters:
	for i, ni := range ns {
		switch ni.Name {
		case "filter":
			f, err := nodeToFilter(ni)
			if err != nil {
				return nil, err
			}

			r.Filters = append(r.Filters, f)
		case "string", "shunt", "loopback":
			ns = ns[i:]
			break filters
		default:
			return nil, errors.New("invalid definition item among filters")
		}
	}

	if len(ns) != 1 {
		return nil, fmt.Errorf("invalid definition backend, remaining definition length: %d, %s",
			len(ns), n.Text())
	}

	if err := nodeToBackend(r, ns[0]); err != nil {
		return nil, err
	}

	return r, nil
}

func treeToEskip(n []*Node) ([]*eskip.Route, error) {
	r := make([]*eskip.Route, 0, len(n))
	for _, ni := range n {
		d, err := nodeToEskipDefinition(ni)
		if err != nil {
			return nil, err
		}

		r = append(r, d)
	}

	return r, nil
}

func checkTerm(t *testing.T, gotName, expectedName string, gotArgs, expectedArgs []interface{}) {
	if gotName != expectedName {
		t.Error("invalid term name")
		return
	}

	// legacy bug support
	for i := len(expectedArgs) - 1; i >= 0; i-- {
		if _, ok := expectedArgs[i].(int); ok {
			expectedArgs = append(expectedArgs[:i], expectedArgs[i+1:]...)
			continue
		}

		if v, ok := expectedArgs[i].(float64); ok && v < 0 {
			gotArgs = append(gotArgs[:i], gotArgs[i+1:]...)
			expectedArgs = append(expectedArgs[:i], expectedArgs[i+1:]...)
		}
	}

	if len(gotArgs) != len(expectedArgs) {
		t.Error("invalid term args length", len(gotArgs), len(expectedArgs))
		return
	}

	for i, a := range gotArgs {
		if a != expectedArgs[i] {
			t.Error("invalid term arg")
			return
		}
	}
}

func checkPredicates(t *testing.T, got, expected *eskip.Route) {
	if got.Path != expected.Path {
		t.Error("invalid path")
		return
	}

	if len(got.HostRegexps) != len(expected.HostRegexps) {
		t.Error("invalid host length")
		return
	}

	for i, h := range got.HostRegexps {
		if h != expected.HostRegexps[i] {
			t.Error("invalid host")
			return
		}
	}

	if len(got.PathRegexps) != len(expected.PathRegexps) {
		t.Error("invalid path regexp length", len(got.PathRegexps), len(expected.PathRegexps))
		return
	}

	for i, h := range got.PathRegexps {
		if h != expected.PathRegexps[i] {
			t.Error("invalid path regexp")
			return
		}
	}

	if got.Method != expected.Method {
		t.Error("invalid method")
		return
	}

	if len(got.Headers) != len(expected.Headers) {
		t.Error("invalid headers length")
		return
	}

	for n, h := range got.Headers {
		he, ok := expected.Headers[n]
		if !ok {
			t.Error("invalid header name")
			return
		}

		if he != h {
			t.Error("invalid header")
			return
		}
	}

	if len(got.HeaderRegexps) != len(expected.HeaderRegexps) {
		t.Error("invalid header regexp length")
		return
	}

	for n, h := range got.HeaderRegexps {
		he, ok := expected.HeaderRegexps[n]
		if !ok {
			t.Error("invalid header regexp name")
			return
		}

		if len(h) != len(he) {
			t.Error("invalid header regexp item length")
			return
		}

		for i, hi := range h {
			if hi != he[i] {
				t.Error("invalid header regexp")
				return
			}
		}
	}

	if len(got.Predicates) != len(expected.Predicates) {
		t.Error("invalid predicates length")
		return
	}

	for i, p := range got.Predicates {
		checkTerm(
			t,
			p.Name, expected.Predicates[i].Name,
			p.Args, expected.Predicates[i].Args,
		)

		if t.Failed() {
			t.Log(p.Name, expected.Predicates[i].Name)
			t.Log(p.Args, expected.Predicates[i].Args)
			return
		}
	}
}

func checkFilters(t *testing.T, got, expected []*eskip.Filter) {
	if len(got) != len(expected) {
		t.Error("invalid filters length")
		return
	}

	for i, f := range got {
		checkTerm(
			t,
			f.Name, expected[i].Name,
			f.Args, expected[i].Args,
		)

		if t.Failed() {
			return
		}
	}
}

func checkBackend(t *testing.T, got, expected *eskip.Route) {
	if got.BackendType != expected.BackendType {
		t.Error("invalid backend type")
		return
	}

	if got.Backend != expected.Backend {
		t.Error("invalid backend")
		return
	}
}

func checkRoute(t *testing.T, got, expected *eskip.Route) {
	if got.Id != expected.Id {
		t.Error("invalid route id")
		return
	}

	checkPredicates(t, got, expected)
	if t.Failed() {
		return
	}

	checkFilters(t, got.Filters, expected.Filters)
	if t.Failed() {
		return
	}

	checkBackend(t, got, expected)
}

func checkEskip(t *testing.T, got, expected []*eskip.Route) {
	if len(got) != len(expected) {
		t.Error("invalid length", len(got), len(expected))
		return
	}

	for i, ri := range got {
		checkRoute(t, ri, expected[i])
		if t.Failed() {
			t.Log(ri.String())
			t.Log(expected[i].String())
			return
		}
	}
}

func eskipTreeToEskip(n *Node) ([]*eskip.Route, error) {
	return treeToEskip(n.Nodes)
}

func TestEskip(t *testing.T) {
	r := generateEskip(1 << 9)
	e := eskip.Print(true, r...)
	b := bytes.NewBufferString(e)
	s, err := testSyntax("eskip.p", 0)
	if err != nil {
		t.Error(err)
		return
	}

	n, err := s.Parse(b)
	if err != nil {
		t.Error(err)
		return
	}

	rback, err := eskipTreeToEskip(n)
	if err != nil {
		t.Error(err)
		return
	}

	checkEskip(t, rback, r)
}
