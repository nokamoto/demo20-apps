package authorizer

import (
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
	code "google.golang.org/genproto/googleapis/rpc/code"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func permissionDenied(message string) *v2.CheckResponse {
	return &v2.CheckResponse{
		Status: &status.Status{
			Code:    int32(code.Code_PERMISSION_DENIED),
			Message: message,
		},
	}
}

func ok(value string) *v2.CheckResponse {
	return &v2.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
		HttpResponse: &v2.CheckResponse_OkResponse{
			OkResponse: &v2.OkHttpResponse{
				Headers: []*core.HeaderValueOption{
					{
						Header: &core.HeaderValue{
							Key:   metadata.MetadataKey,
							Value: value,
						},
						Append: &wrapperspb.BoolValue{
							Value: false,
						},
					},
				},
			},
		},
	}
}
