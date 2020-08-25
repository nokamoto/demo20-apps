package compute

// Instance represents a instance table.
type Instance struct {
	InstanceKey int64  `gorm:"column:instance_key;auto_increment;primary_key"`
	InstanceID  string `gorm:"column:instance_id"`
	ParentKey   int64  `gorm:"column:parent_key"`
	Labels      string `gorm:"column:labels"`
}

// TableName returns a table name.
func (Instance) TableName() string {
	return "compute_instance"
}
