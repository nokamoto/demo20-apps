package main

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.Main(func(logger *zap.Logger, server *grpc.Server, db *gorm.DB) error {
		v1alpha.RegisterRdbServer(server, &v1alpha.UnimplementedRdbServer{})
		return nil
	})
}
