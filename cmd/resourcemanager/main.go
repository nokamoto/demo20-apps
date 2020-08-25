package main

import (
	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, server *grpc.Server) error {
		v1alpha.RegisterResourceManagerServer(server, &v1alpha.UnimplementedResourceManagerServer{})
		return nil
	})
}
