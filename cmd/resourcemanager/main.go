package main

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application/resourcemanager"
	"github.com/nokamoto/demo20-apps/internal/server"
	service "github.com/nokamoto/demo20-apps/internal/service/resourcemanager"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, s *grpc.Server, db *gorm.DB) {
		v1alpha.RegisterResourceManagerServer(s, service.NewService(resourcemanager.NewResourceManager(db), logger))
	})
}
