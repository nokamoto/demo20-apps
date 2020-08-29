//go:generate mockgen -source=$GOFILE -destination=mock.go -package=$GOPACKAGE
package iam

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/iam"
)

type permissionQuery interface {
	Create(*gorm.DB, *iam.Permission) error
}
