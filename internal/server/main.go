package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	loggerDebug    = "LOGGER_DEBUG"
	grpcServerPort = "GRPC_SERVER_PORT"
)

const (
	mysqlUser     = "MYSQL_USER"
	mysqlPassword = "MYSQL_PASSWORD"
	mysqlHost     = "MYSQL_HOST"
	mysqlPort     = "MYSQL_PORT"
	mysqlDatabase = "MYSQL_DATABASE"
)

// Main serves a gRPC server with a reflection service for a main func.
//
// environment variables:
//   "LOGGER_DEBUG" - prints debug level logs if set a non empty string.
//   "GRPC_SERVER_PORT" - serves with the port number.
func Main(register func(*zap.Logger, *grpc.Server, *gorm.DB)) {
	MainWithoutMySQL(func(logger *zap.Logger, s *grpc.Server) {
		db, err := mySQL()
		if err != nil {
			logger.Fatal("failed to connect mysql", zap.Error(err))
		}
		register(logger, s, db)
	})
}

// MainWithoutMySQL serves a gRPC server with a reflection service for a main func.
func MainWithoutMySQL(register func(*zap.Logger, *grpc.Server)) {
	rand.Seed(time.Now().Unix())

	cfg := zap.NewProductionConfig()
	if len(os.Getenv(loggerDebug)) != 0 {
		cfg.Level.SetLevel(zap.DebugLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("logger created", zap.Any("level", cfg.Level))

	port, err := strconv.Atoi(os.Getenv(grpcServerPort))
	if err != nil {
		logger.Fatal("invalid port", zap.Error(err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("failed to listen tcp port", zap.Int("port", port), zap.Error(err))
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	register(logger, server)

	reflection.Register(server)

	logger.Info("ready to serve", zap.Int("port", port))
	err = server.Serve(lis)
	if err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
