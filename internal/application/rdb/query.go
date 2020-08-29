//go:generate mockgen -source=$GOFILE -destination=query_mock.go -package=$GOPACKAGE
package rdb

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/rdb"
)

type clusterQuery interface {
	Create(tx *gorm.DB, cluster *rdb.Cluster, instanceIDs []string) error
}
