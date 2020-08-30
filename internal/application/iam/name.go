package iam

import (
	"strings"
)

func fromRoleNameF(name string) (string, string) {
	ids := strings.Split(name, "/")
	if len(ids) == 2 {
		return "/", ids[1]
	}
	return ids[1], ids[3]
}
