package resourcemanager

// Project represents a project table.
type Project struct {
	projectKey  int64  `gorm:"column:project_key;auto_increment;primary_key"`
	projectID   string `gorm:"column:project_id"`
	displayName string `gorm:"column:display_name"`
}

// TableName returns a table name.
func (Project) TableName() string {
	return "resourcemanager_project"
}
