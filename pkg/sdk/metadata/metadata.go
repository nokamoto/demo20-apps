package metadata

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"google.golang.org/grpc/metadata"
)

const (
	metadataKey = "x-cloud-metadata"
)

var (
	enc = base64.RawStdEncoding
)

func encode(md *api.Metadata) (string, error) {
	bytes, err := proto.Marshal(md)
	if err != nil {
		return "", err
	}
	return enc.EncodeToString(bytes), nil
}

// AppendToOutgoingContext appends Metadata to the context.
func AppendToOutgoingContext(ctx context.Context, md *api.Metadata) (context.Context, error) {
	s, err := encode(md)
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, metadataKey, s), nil
}

// AppendToOutgoingContextF appends Metadata to the context for testing.
func AppendToOutgoingContextF(ctx context.Context, md *api.Metadata) context.Context {
	res, err := AppendToOutgoingContext(ctx, md)
	if err != nil {
		panic(err)
	}
	return res
}

// NewIncomingContext appends Metadata to the context.
func NewIncomingContext(ctx context.Context, md *api.Metadata) (context.Context, error) {
	s, err := encode(md)
	if err != nil {
		return nil, err
	}
	return metadata.NewIncomingContext(ctx, metadata.MD{metadataKey: []string{s}}), nil
}

// NewIncomingContextF appends Metadata to the context for testing.
func NewIncomingContextF(ctx context.Context, md *api.Metadata) context.Context {
	res, err := NewIncomingContext(ctx, md)
	if err != nil {
		panic(err)
	}
	return res
}

func fromMD(md metadata.MD, ok bool) (*api.Metadata, error) {
	if !ok {
		return nil, errors.New("no metadata")
	}

	xs, ok := md[metadataKey]
	if !ok {
		return nil, errors.New("no metadata")
	}
	if len(xs) != 1 {
		return nil, errors.New("invalid metadata")
	}

	bytes, err := enc.DecodeString(xs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	var value api.Metadata
	err = proto.Unmarshal(bytes, &value)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	return &value, nil
}

// FromIncomingContext returns Metadata from the incoming context.
func FromIncomingContext(ctx context.Context) (*api.Metadata, error) {
	return fromMD(metadata.FromIncomingContext(ctx))
}

// FromOutgoingContext returns Metadata from the outgoing context.
func FromOutgoingContext(ctx context.Context) (*api.Metadata, error) {
	return fromMD(metadata.FromOutgoingContext(ctx))
}
