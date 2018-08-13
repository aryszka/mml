package records

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Totals map[string]int

func Report(url string) (r Totals, err error) {
	var rsp *http.Response
	rsp, err = http.Get(url)
	if err != nil {
		return
	}

	defer rsp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	var doc []struct {
		Name  *string `json:"omitempty"`
		Value *int    `json:"omitempty"`
	}
	if err = json.Unmarshal(&doc); err != nil {
		return
	}

	r = make(Totals)
	for i := range doc {
		if doc[i].Name == nil || doc[i].Value == nil {
			err = errors.New("invalid json result")
			return
		}

		r[*doc[i].Name] += *doc[i].Value
		r["sum"] += *doc[i].Value
	}

	return
}
