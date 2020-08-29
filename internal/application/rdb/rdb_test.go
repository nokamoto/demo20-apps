package rdb

import (
	"context"
	"testing"

	"github.com/nokamoto/demo20-apps/internal/mysql/rdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	compute "github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/test"
)

type testCase struct {
	name  string
	run   func(*testing.T, Rdb) error
	mock  func(*MockclusterQuery, *MockinstanceClient)
	check test.Check
	tx    func(sqlmock.Sqlmock)
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			test.UseGorm(t, func(m sqlmock.Sqlmock, g *gorm.DB) {
				x.tx(m)

				c := NewMockclusterQuery(ctrl)
				i := NewMockinstanceClient(ctrl)
				x.mock(c, i)

				err := x.run(t, Rdb{
					clusterQuery:   c,
					instanceClient: i,
					db:             g,
				})

				x.check(t, err)
			})
		})
	}
}

func TestRdb_Create(t *testing.T) {
	run := func(id, parentID string, cluster, expected *v1alpha.Cluster) func(*testing.T, Rdb) error {
		return func(t *testing.T, r Rdb) error {
			return test.Diff1IgnoreUnexported(r.Create(context.Background(), id, parentID, cluster))(t, expected)
		}
	}
	xs := testCases{
		{
			name: "OK",
			run: run(
				"foo", "bar",
				&v1alpha.Cluster{
					Replicas: 1,
				},
				&v1alpha.Cluster{
					Name:      "clusters/foo",
					Replicas:  1,
					Instances: []string{"instances/baz", "instances/qux"},
					Parent:    "projects/bar",
				},
			),
			mock: func(c *MockclusterQuery, i *MockinstanceClient) {
				gomock.InOrder(
					i.EXPECT().CreateInstance(gomock.Any(), test.ProtoEq(&compute.CreateInstanceRequest{
						Instance: &compute.Instance{
							Labels: []string{"rdb", "master"},
						},
					})).Return(&compute.Instance{
						Name: "instances/baz",
					}, nil),
					i.EXPECT().CreateInstance(gomock.Any(), test.ProtoEq(&compute.CreateInstanceRequest{
						Instance: &compute.Instance{

							Labels: []string{"rdb", "replica"},
						},
					})).Return(&compute.Instance{
						Name: "instances/qux",
					}, nil),
					c.EXPECT().Create(gomock.Any(),
						&rdb.Cluster{
							ClusterID: "foo",
							Replicas:  1,
							ParentID:  "bar",
						},
						[]string{"baz", "qux"},
					).Return(nil),
				)
			},
			check: test.Succeeded,
			tx:    test.Commit,
		},
	}

	xs.run(t)
}
