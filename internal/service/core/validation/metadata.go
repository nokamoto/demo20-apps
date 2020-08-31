package validation

import (
	"context"

	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
)

// ProjectIncomingContext returns a project id from the incoming context.
func ProjectIncomingContext(ctx context.Context) (string, error) {
	md, err := metadata.FromIncomingContext(ctx)
	if err != nil {
		return "", err
	}
	var ids []string
	err = FromName(md.GetParent(), &ids, "projects")
	if err != nil {
		return "", nil
	}
	return ids[0], nil
}
