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
