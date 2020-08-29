//go:generate mockgen -source=$GOFILE -destination=mock.go -package=$GOPACKAGE
package compute

import (
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
)

type compute interface {
	Create(string, string, *v1alpha.Instance) (*v1alpha.Instance, error)
	RandomName(string) string
}
