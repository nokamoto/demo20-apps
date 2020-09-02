//go:generate mockgen -source=$GOFILE -destination=client_mock.go -package=$GOPACKAGE
package authorizer

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"google.golang.org/grpc"
)

type client interface {
	AuthorizeMachineUser(ctx context.Context, in *v1alpha.AuthorizeMachineUserRequest, opts ...grpc.CallOption) (*v1alpha.AuthorizeMachineUserResponse, error)
}
