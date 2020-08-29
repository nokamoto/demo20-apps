//go:generate mockgen -source=$GOFILE -destination=mock.go -package=$GOPACKAGE
package resourcemanager

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/resourcemanager"
)

type projectQuery interface {
	Create(tx *gorm.DB, project *resourcemanager.Project) error
	Get(tx *gorm.DB, id string) (*resourcemanager.Project, error)
}
