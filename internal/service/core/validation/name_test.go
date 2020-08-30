package validation

import "testing"

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
