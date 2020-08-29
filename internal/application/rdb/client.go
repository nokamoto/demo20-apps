//go:generate mockgen -source=$GOFILE -destination=client_mock.go -package=$GOPACKAGE
package rdb

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"google.golang.org/grpc"
)

type instanceClient interface {
	CreateInstance(ctx context.Context, in *v1alpha.CreateInstanceRequest, opts ...grpc.CallOption) (*v1alpha.Instance, error)
}
