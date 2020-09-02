package main

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application/compute"
	"github.com/nokamoto/demo20-apps/internal/server"
	service "github.com/nokamoto/demo20-apps/internal/service/compute"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, s *grpc.Server, db *gorm.DB) {
		v1alpha.RegisterComputeServer(s, service.NewService(compute.NewCompute(db), logger))
	})
}
