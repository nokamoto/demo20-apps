package validation

import (
	"fmt"
	"strings"
)

// FromName returns ids from the name.
func FromName(name string, res *[]string, collectionIDs ...string) error {
	*res = make([]string, len(collectionIDs))

	e := fmt.Errorf("expected %s/* but actual %s", strings.Join(collectionIDs, "/*/"), name)

	ids := strings.Split(name, "/")
	if len(ids) != len(collectionIDs)*2 {
		return e
	}

	for i, collectionID := range collectionIDs {
		if ids[i*2] != collectionID {
			return e
		}
		(*res)[i] = ids[i*2+1]
	}

	return nil
}

// NameOr returns an error if the name if not in the form of collection ids.
func NameOr(name string, collectionIDs []string, orIDs []string) error {
	e := fmt.Errorf("%s: unexpected name: %v or %v", name, collectionIDs, orIDs)
	ids := strings.Split(name, "/")

	test := func(xs []string) bool {
		if len(ids) != len(xs)*2 {
			return false
		}
		for i, x := range xs {
			if ids[i*2] != x {
				return false
			}
		}
		return true
	}

	if test(collectionIDs) || test(orIDs) {
		return nil
	}

	return e
}
