package resourcemanager

// Project represents a project table.
type Project struct {
	ProjectKey  int64  `gorm:"column:project_key;auto_increment;primary_key"`
	ProjectID   string `gorm:"column:project_id"`
	DisplayName string `gorm:"column:display_name"`
}

// TableName returns a table name.
func (Project) TableName() string {
	return "resourcemanager_project"
}
