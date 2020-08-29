package main

import (
	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	iam "github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	application "github.com/nokamoto/demo20-apps/internal/application/iam"
	"github.com/nokamoto/demo20-apps/internal/server"
	service "github.com/nokamoto/demo20-apps/internal/service/iam"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, s *grpc.Server) error {
		db, err := server.MySQL()
		if err != nil {
			return err
		}

		admin.RegisterIamServer(s, service.NewAdminService(
			application.NewIam(db),
			logger,
		))

		iam.RegisterIamServer(s, &iam.UnimplementedIamServer{})

		return nil
	})
}
