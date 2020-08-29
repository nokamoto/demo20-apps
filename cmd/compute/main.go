package main

import (
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application/compute"
	"github.com/nokamoto/demo20-apps/internal/server"
	service "github.com/nokamoto/demo20-apps/internal/service/compute"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, s *grpc.Server) error {
		db, err := server.MySQL()
		if err != nil {
			return err
		}

		compute := compute.NewCompute(db)
		v1alpha.RegisterComputeServer(s, service.NewService(compute, logger))

		return nil
	})
}
