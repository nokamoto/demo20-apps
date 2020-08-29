package rdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	compute "github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application"
	"github.com/nokamoto/demo20-apps/internal/mysql/rdb"
)

// Rdb defines a business logic for the cloud rdb service.
type Rdb struct {
	clusterQuery   clusterQuery
	instanceClient instanceClient
	db             *gorm.DB
}

// NewRdb returns Rdb.
func NewRdb(db *gorm.DB, compute compute.ComputeClient) *Rdb {
	return &Rdb{
		clusterQuery:   rdb.Query{},
		instanceClient: compute,
		db:             db,
	}
}

func (r *Rdb) createInstances(ctx context.Context, parentID string, size int) ([]string, []string, error) {
	var instanceIDs []string
	var instanceNames []string
	for i := 0; i < size+1; i++ {
		var typ string
		if i == 0 {
			typ = "master"
		} else {
			typ = "replica"
		}

		res, err := r.instanceClient.CreateInstance(ctx, &compute.CreateInstanceRequest{
			Instance: &compute.Instance{
				Labels: []string{"rdb", typ},
			},
		})
		if err == nil {
			instanceNames = append(instanceNames, res.GetName())

			ids := strings.Split(res.GetName(), "/")
			if len(ids) != 2 {
				return nil, nil, fmt.Errorf("%v: %w", err, application.ErrInternal)
			}

			instanceIDs = append(instanceIDs, ids[1])

			continue
		}
		return nil, nil, err
	}
	return instanceIDs, instanceNames, nil
}

// Create creates a cluster with creating cloud compute instances.
func (r *Rdb) Create(ctx context.Context, id, parentID string, cluster *v1alpha.Cluster) (*v1alpha.Cluster, error) {
	instanceIDs, instanceNames, err := r.createInstances(ctx, parentID, int(cluster.GetReplicas()))
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := r.clusterQuery.Create(tx, &rdb.Cluster{
			ClusterID: id,
			Replicas:  cluster.GetReplicas(),
			ParentID:  parentID,
		}, instanceIDs)
		if err != nil {
			return application.Error(err, application.ErrorMap{
				application.AlreadyExists: id,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.Cluster{
		Name:      fmt.Sprintf("clusters/%s", id),
		Replicas:  cluster.GetReplicas(),
		Instances: instanceNames,
		Parent:    fmt.Sprintf("projects/%s", parentID),
	}, nil
}
