package mml

import (
	"io"
	"io/ioutil"
)

func EvalInput(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	s := string(b)

	m, err := parseModule(s)
	if err != nil {
		return err
	}

	evalModule(newEnv(), m)
	return nil
}
