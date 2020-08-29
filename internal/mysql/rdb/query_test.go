package rdb

import (
	"regexp"
	"testing"

	"github.com/nokamoto/demo20-apps/internal/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
)

func TestQuery_Create(t *testing.T) {
	run := func(cluster Cluster, instances ...string) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return Query{}.Create(tx, &cluster, instances)
		}
	}

	cluster := Cluster{
		ClusterID: "foo",
		Replicas:  10,
		ParentID:  "bar",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(cluster, "baz"),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `rdb_cluster` (`cluster_id`,`replicas`,`parent_id`) VALUES (?,?,?)")).
					WithArgs(cluster.ClusterID, cluster.Replicas, cluster.ParentID).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `rdb_cluster_instance` (`cluster_key`,`instance_id`) VALUES (?,?)")).
					WithArgs(1000, "baz").
					WillReturnResult(sqlmock.NewResult(2000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestQuery_Delete(t *testing.T) {
	run := func(cluster Cluster) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return Query{}.Delete(tx, &cluster)
		}
	}

	cluster := Cluster{
		ClusterKey: 100,
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(cluster),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `rdb_cluster_instance` WHERE (cluster_key = ?)")).
					WithArgs(cluster.ClusterKey).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `rdb_cluster` WHERE (cluster_key = ?)")).
					WithArgs(cluster.ClusterKey).
					WillReturnResult(sqlmock.NewResult(2000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func clusterRows(xs ...Cluster) *sqlmock.Rows {
	v := sqlmock.NewRows([]string{
		"cluster_id", "cluster_key", "parent_id", "replicas",
	})
	for _, x := range xs {
		v.AddRow(x.ClusterID, x.ClusterKey, x.ParentID, x.Replicas)
	}
	return v
}

func instanceRows(xs ...ClusterInstance) *sqlmock.Rows {
	v := sqlmock.NewRows([]string{
		"cluster_key", "instance_id",
	})
	for _, x := range xs {
		v.AddRow(x.ClusterKey, x.InstanceID)
	}
	return v
}

func TestQuery_Get(t *testing.T) {
	run := func(id string, cexpected *Cluster, iexpected []*ClusterInstance) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff2(Query{}.Get(tx, id))(t, cexpected, iexpected)
		}
	}

	cluster := Cluster{
		ClusterKey: 100,
		ClusterID:  "foo",
	}

	instance := ClusterInstance{
		ClusterKey: 100,
		InstanceID: "baz",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(cluster.ClusterID, &cluster, []*ClusterInstance{&instance}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rdb_cluster` WHERE (cluster_id = ?) LIMIT 1")).
					WithArgs(cluster.ClusterID).
					WillReturnRows(clusterRows(cluster))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rdb_cluster_instance`  WHERE (cluster_key = ?)")).
					WithArgs(cluster.ClusterKey).
					WillReturnRows(instanceRows(instance))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestQuery_List(t *testing.T) {
	offset, limit := 100, 200

	run := func(parentID string, cexpected []*Cluster, iexpected []*ClusterInstance) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			cactual, iactual, err := Query{}.List(tx, parentID, offset, limit)
			if diff := cmp.Diff(cexpected, cactual); len(diff) != 0 {
				t.Error(diff)
			}
			if diff := cmp.Diff(iexpected, iactual); len(diff) != 0 {
				t.Error(diff)
			}
			return err
		}
	}

	cluster := Cluster{
		ClusterKey: 100,
		ClusterID:  "foo",
		ParentID:   "bar",
	}

	instance := ClusterInstance{
		ClusterKey: 100,
		InstanceID: "baz",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(cluster.ParentID, []*Cluster{&cluster}, []*ClusterInstance{&instance}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rdb_cluster` WHERE (parent_id = ?)")).
					WithArgs(cluster.ParentID).
					WillReturnRows(clusterRows(cluster))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rdb_cluster_instance`  WHERE (cluster_key in (?))")).
					WithArgs(cluster.ClusterKey).
					WillReturnRows(instanceRows(instance))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}
