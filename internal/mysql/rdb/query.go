package rdb

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/core"
)

// Query defines queries for rdb tables within a transaction.
type Query struct{}

// Create inserts cluster and instance records.
func (q Query) Create(tx *gorm.DB, cluster *Cluster, instanceKeys []int64) error {
	res := tx.Debug().Create(cluster)
	if res.Error != nil {
		return core.Translate(res.Error)
	}

	var instances bulkClusterInstance
	for _, i := range instanceKeys {
		instances = append(instances, &ClusterInstance{
			ClusterKey:  cluster.ClusterKey,
			InstanceKey: i,
		})
	}

	return core.BulkInsert(tx, ClusterInstance{}.TableName(), instances)
}

// Delete deletes cluster and instance records.
func (Query) Delete(tx *gorm.DB, cluster *Cluster) error {
	res := tx.Debug().Where("cluster_key = ?", cluster.ClusterKey).Delete(&ClusterInstance{})
	if res.Error != nil {
		return core.Translate(res.Error)
	}

	res = tx.Debug().Where("cluster_key = ?", cluster.ClusterKey).Delete(&Cluster{})
	if res.Error != nil {
		return core.Translate(res.Error)
	}

	return nil
}

// Get returns cluster and instance records by the cluster id.
func (Query) Get(tx *gorm.DB, id string) (*Cluster, []*ClusterInstance, error) {
	var cluster Cluster
	err := core.Get(tx, &cluster, "cluster_id = ?", id)
	if err != nil {
		return nil, nil, err
	}

	var instances []*ClusterInstance
	err = core.Get(tx, &instances, "cluster_key = ?", cluster.ClusterKey)
	if err != nil {
		return nil, nil, err
	}

	return &cluster, instances, nil
}

// List returns cluster and instance records by the parent key.
func (Query) List(tx *gorm.DB, parentKey int64, offset, limit int) ([]*Cluster, []*ClusterInstance, error) {
	var clusters []*Cluster
	err := core.List(tx, &clusters, offset, limit, "parent_key = ?", parentKey)
	if err != nil {
		return nil, nil, err
	}

	var keys []int64
	for _, c := range clusters {
		keys = append(keys, c.ClusterKey)
	}

	var instances []*ClusterInstance
	return clusters, instances, core.ListAll(tx, &instances, "cluster_key in (?)", keys)
}
