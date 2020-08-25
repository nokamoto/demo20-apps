package compute

// Instance represents a instance table.
type Instance struct {
	instanceKey int64  `gorm:"column:instance_key;auto_increment;primary_key"`
	instanceID  string `gorm:"column:instance_id"`
	parentKey   int64  `gorm:"column:parent_key"`
	labels      string `gorm:"column:labels"`
}

// TableName returns a table name.
func (Instance) TableName() string {
	return "compute_instance"
}
