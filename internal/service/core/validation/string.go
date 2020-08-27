package validation

import (
	"fmt"
)

// Empty returns an error if s is not empty.
func Empty(s string) error {
	if len(s) != 0 {
		return fmt.Errorf("expected emtpy but actual %s", s)
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
