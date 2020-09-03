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
	// MetadataKey is a gRPC metadata key for Metadata.
	MetadataKey = "x-cloud-metadata"
)

var (
	enc = base64.RawStdEncoding
)

// Encode encodes Metadata to a base64 string.
func Encode(md *api.Metadata) (string, error) {
	bytes, err := proto.Marshal(md)
	if err != nil {
		return "", err
	}
	return enc.EncodeToString(bytes), nil
}

// Decode decodes Metadata from the base64 string.
func Decode(value string) (*api.Metadata, error) {
	bytes, err := enc.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	var md api.Metadata
	err = proto.Unmarshal(bytes, &md)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	return &md, nil
}

// AppendToOutgoingContext appends Metadata to the context.
func AppendToOutgoingContext(ctx context.Context, md *api.Metadata) (context.Context, error) {
	s, err := Encode(md)
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, MetadataKey, s), nil
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
	s, err := Encode(md)
	if err != nil {
		return nil, err
	}
	return metadata.NewIncomingContext(ctx, metadata.MD{MetadataKey: []string{s}}), nil
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

	xs, ok := md[MetadataKey]
	if !ok {
		return nil, errors.New("no metadata")
	}
	if len(xs) != 1 {
		return nil, errors.New("invalid metadata")
	}

	return Decode(xs[0])
}

// FromIncomingContext returns Metadata from the incoming context.
func FromIncomingContext(ctx context.Context) (*api.Metadata, error) {
	return fromMD(metadata.FromIncomingContext(ctx))
}

// FromOutgoingContext returns Metadata from the outgoing context.
func FromOutgoingContext(ctx context.Context) (*api.Metadata, error) {
	return fromMD(metadata.FromOutgoingContext(ctx))
}
