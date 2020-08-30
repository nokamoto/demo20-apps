package application

import (
	"strings"
)

// FromNameF returns an id from the name in the form of `{collection}/{resource}`.
func FromNameF(name string) string {
	if name == "projects//" {
		return "/"
	}
	return strings.Split(name, "/")[1]
}
