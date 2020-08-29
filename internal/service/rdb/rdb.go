//go:generate mockgen -source=$GOFILE -destination=rdb_mock.go -package=$GOPACKAGE
package rdb

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
)

type rdb interface {
	Create(ctx context.Context, id, parentID string, cluster *v1alpha.Cluster) (*v1alpha.Cluster, error)
}
