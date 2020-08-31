package validation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNameOr(t *testing.T) {
	xs := []string{"roles"}
	ys := []string{"projects", "roles"}

	err := NameOr("roles/foo", xs, ys)
	if err != nil {
		t.Fatal(err)
	}

	err = NameOr("projects/foo/roles/bar", xs, ys)
	if err != nil {
		t.Fatal(err)
	}

	err = NameOr("foo", xs, ys)
	if err == nil {
		t.Fatal("expected err")
	}

	err = NameOr("xxx/foo", xs, ys)
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestFromName(t *testing.T) {
	var res []string
	err := FromName("projects/foo", &res, "projects")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{"foo"}, res); len(diff) != 0 {
		t.Error(diff)
	}
}
