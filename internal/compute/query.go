package compute

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/compute"
	"github.com/nokamoto/demo20-apps/internal/mysql/resourcemanager"
)

type instanceQuery interface {
	Create(*gorm.DB, *compute.Instance) error
}

type projectQuery interface {
	Get(*gorm.DB, string) (*resourcemanager.Project, error)
}
