package main

import (
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, server *grpc.Server) error {
		v1alpha.RegisterComputeServer(server, &v1alpha.UnimplementedComputeServer{})
		return nil
	})
}
