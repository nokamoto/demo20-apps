package main

import (
	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	iam "github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, server *grpc.Server) error {
		admin.RegisterIamServer(server, &admin.UnimplementedIamServer{})
		iam.RegisterIamServer(server, &iam.UnimplementedIamServer{})
		return nil
	})
}
