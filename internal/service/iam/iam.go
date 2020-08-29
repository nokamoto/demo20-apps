//go:generate mockgen -source=$GOFILE -destination=iam_mock.go -package=$GOPACKAGE
package iam

import (
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
)

type iam interface {
	Create(string) (*v1alpha.Permission, error)
}
