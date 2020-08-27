package automatedtest

import (
	"log"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	grpcAddress = "GRPC_ADDRESS"
	loggerDebug = "LOGGER_DEBUG"
)

// Main runs automated test cases.
//
// environment variables:
//   "LOGGER_DEBUG" - prints debug level logs if set a non empty string.
//   "GRPC_ADDRESS" - connects to the gRPC server address.
func Main(f func(*grpc.ClientConn) Scenarios) {
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

	address := os.Getenv(grpcAddress)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("failed to dial", zap.String("address", address), zap.Error(err))
	}

	f(conn).run(logger)
}
