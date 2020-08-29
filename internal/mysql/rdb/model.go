package rdb

// Cluster represents a cluster table.
type Cluster struct {
	ClusterKey int64  `gorm:"column:cluster_key;auto_increment;primary_key"`
	ClusterID  string `gorm:"column:cluster_id"`
	Replicas   int32  `gorm:"column:replicas"`
	ParentID   string `gorm:"column:parent_id"`
}

// TableName returns a table name.
func (Cluster) TableName() string {
	return "rdb_cluster"
}

// ClusterInstance represents a cluster instance table.
type ClusterInstance struct {
	ClusterKey int64  `gorm:"column:cluster_key"`
	InstanceID string `gorm:"column:instance_id"`
}

// TableName returns a table name.
func (ClusterInstance) TableName() string {
	return "rdb_cluster_instance"
}

type bulkClusterInstance []*ClusterInstance

func (xs bulkClusterInstance) Args() [][]interface{} {
	var res [][]interface{}
	for _, x := range xs {
		var args []interface{}
		args = append(args, x.ClusterKey)
		args = append(args, x.InstanceID)
		res = append(res, args)
	}
	return res
}

func (bulkClusterInstance) Fields() []string {
	return []string{"cluster_key", "instance_id"}
}
