package validation

import (
	"fmt"
)

// Range returns an error if i is not in the range.
func Range(i, lb, ub int) error {
	if i < lb || ub < i {
		return fmt.Errorf("expected [%d, %d] but actual %d", lb, ub, i)
	}
	return nil
}
