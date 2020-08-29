package automatedtest

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// Diff asserts that the actual proto message is eqaul to the expected proto message.
func Diff(expected proto.Message, actual proto.Message, opts ...cmp.Option) error {
	opts = append(opts, protocmp.Transform())

	if diff := cmp.Diff(expected, actual, opts...); len(diff) != 0 {
		return fmt.Errorf("diff=%s", diff)
	}

	return nil
}
