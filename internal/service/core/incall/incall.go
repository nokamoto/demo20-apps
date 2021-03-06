package incall

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/nokamoto/demo20-apps/internal/application"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InCall provides shared functions for gRPC server methods.
type InCall struct {
	Logger *zap.Logger
}

// NewInCall returns InCall for the request.
func NewInCall(logger *zap.Logger, method string, req proto.Message) *InCall {
	m := jsonpb.Marshaler{}
	json, err := m.MarshalToString(req)
	if err != nil {
		logger = logger.With(zap.String("method", method), zap.Any("req", req))
		logger.Error("failed to marshal", zap.Error(err))
	} else {
		logger = logger.With(zap.String("method", method), zap.String("req", json))
	}

	logger.Debug("method called")

	return &InCall{
		Logger: logger,
	}
}

// Error converts from the application error to a gRPC server error.
func (i *InCall) Error(err error) error {
	if errors.Is(err, application.ErrNotFound) {
		i.Logger.Debug("resource not found", zap.Error(err))
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, application.ErrAlreadyExists) {
		i.Logger.Debug("resource already exists", zap.Error(err))
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, application.ErrInternal) {
		i.Logger.Error("internal", zap.Error(err))
		return status.Error(codes.Internal, "internal")
	}
	if errors.Is(err, application.ErrPermissionDenied) {
		i.Logger.Debug("permission denied", zap.Error(err))
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	i.Logger.Error("unhandled error: unavailable", zap.Error(err))
	return status.Error(codes.Unavailable, "unavailable")
}

// InvalidArgument converts errors to InvalidArgument.
func (i *InCall) InvalidArgument(errs []error) error {
	var xs []string
	for _, x := range errs {
		xs = append(xs, x.Error())
	}

	s := fmt.Sprintf("[%s]", strings.Join(xs, ", "))

	i.Logger.Debug("invalid argument", zap.Errors("errors", errs))
	return status.Error(codes.InvalidArgument, s)
}
