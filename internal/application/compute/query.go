//go:generate mockgen -source=$GOFILE -destination=mock.go -package=$GOPACKAGE
package compute

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/compute"
)

type instanceQuery interface {
	Create(*gorm.DB, *compute.Instance) error
}
