package next

import (
	"bytes"
	"math/rand"
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

func TestEskip(t *testing.T) {
	r := generateEskip(1 << 9)
	e := eskip.Print(true, r...)
	b := bytes.NewBufferString(e)
	s, err := testSyntax("eskip.p", 0)
	if err != nil {
		t.Error(err)
		return
	}

	println(e)
	_, err = s.Parse(b)
	if err != nil {
		t.Error(err)
		return
	}
}
