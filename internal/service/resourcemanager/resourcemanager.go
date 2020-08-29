//go:generate mockgen -source=$GOFILE -destination=mock.go -package=$GOPACKAGE
package resourcemanager

import (
	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
)

type resourcemanager interface {
	Get(string) (*v1alpha.Project, error)
	Create(id string, project *v1alpha.Project) (*v1alpha.Project, error)
}
