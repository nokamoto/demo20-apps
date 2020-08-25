package rdb

// Cluster represents a cluster table.
type Cluster struct {
	clusterKey int64  `gorm:"column:cluster_key;auto_increment;primary_key"`
	clusterID  string `gorm:"column:cluster_id"`
	replicas   int32  `gorm:"column:replicas"`
	parentKey  int64  `gorm:"column:parent_key"`
}

// TableName returns a table name.
func (Cluster) TableName() string {
	return "rdb_cluster"
}

// ClusterInstance represents a cluster instance table.
type ClusterInstance struct {
	clusterKey  int64 `gorm:"column:cluster_key"`
	instanceKey int64 `gorm:"column:instance_key"`
}

// TableName returns a table name.
func (ClusterInstance) TableName() string {
	return "rdb_cluster_instance"
}
