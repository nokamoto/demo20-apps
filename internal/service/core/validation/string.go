package validation

import (
	"fmt"
	"regexp"
)

// Empty returns an error if s is not empty.
func Empty(s string) error {
	if len(s) != 0 {
		return fmt.Errorf("expected emtpy but actual %s", s)
	}
	return nil
}

// ID returns an erros if s is not an id string.
func ID(s string) error {
	exp := "^[a-z0-9]{1,16}$"
	if !regexp.MustCompile(exp).MatchString(s) {
		return fmt.Errorf("expected %s but actual %s", exp, s)
	}
	return nil
}

// Concat concats errors exlucdes nil.
func Concat(errs ...error) []error {
	var xs []error
	for _, err := range errs {
		if err != nil {
			xs = append(xs, err)
		}
	}
	return xs
}
