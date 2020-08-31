package main

import (
	"github.com/jinzhu/gorm"
	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	iam "github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	application "github.com/nokamoto/demo20-apps/internal/application/iam"
	"github.com/nokamoto/demo20-apps/internal/server"
	service "github.com/nokamoto/demo20-apps/internal/service/iam"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, s *grpc.Server, db *gorm.DB) error {
		admin.RegisterIamServer(s, service.NewAdminService(application.NewIam(db), logger))

		iam.RegisterIamServer(s, service.NewService(application.NewIam(db), logger))

		return nil
	})
}
