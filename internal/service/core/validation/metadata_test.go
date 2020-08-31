package validation

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
)

func TestProjectIncomingContext(t *testing.T) {
	expected := "foo"
	actual, err := ProjectIncomingContext(metadata.NewIncomingContextF(context.Background(), &api.Metadata{
		Parent: fmt.Sprintf("projects/%s", expected),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expected, actual); len(diff) != 0 {
		t.Error(diff)
	}
}
