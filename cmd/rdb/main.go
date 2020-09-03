package main

import (
	"os"

	"github.com/jinzhu/gorm"
	compute "github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application/rdb"
	"github.com/nokamoto/demo20-apps/internal/server"
	service "github.com/nokamoto/demo20-apps/internal/service/rdb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, s *grpc.Server, db *gorm.DB) {
		address := os.Getenv("COMPUTE_GRPC_ADDRESS")
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			logger.Fatal("failed to create conn", zap.Error(err))
		}

		compute := compute.NewComputeClient(conn)

		v1alpha.RegisterRdbServer(s, service.NewService(rdb.NewRdb(db, compute), logger))
	})
}
