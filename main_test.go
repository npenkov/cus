package cus

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/ory/dockertest"
)

var store UniqueStore

func TestMain(m *testing.M) {
	ret := testMainWrapper(m)
	os.Exit(ret)
}

func testMainWrapper(m *testing.M) int {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	var resource *dockertest.Resource
	var closeFunction func()

	resource, store, closeFunction = GetTestCassandraStorage(pool)
	defer closeFunction()

	if err := store.CreateSchema(true); err != nil {
		panic(fmt.Errorf("failed to create schema : %v", err))
	}

	retCode := m.Run()

	if resource != nil {
		err = pool.Purge(resource)
	}
	return retCode
}

func GetTestCassandraStorage(pool *dockertest.Pool) (resource *dockertest.Resource, store UniqueStore, closeFunction func()) {
	var err error
	resource, err = pool.Run("cassandra", "3", []string{})
	if err != nil {
		panic(err)
	}
	hostAndPort := resource.GetHostPort("9042/tcp")
	if err = pool.Retry(func() error {
		var err error
		clusterSystem := gocql.NewCluster(hostAndPort)
		clusterSystem.Keyspace = "system"
		clusterSystem.Timeout = time.Millisecond * 2000
		clusterSystem.DisableInitialHostLookup = true
		sessionSystem, err := clusterSystem.CreateSession()
		if err != nil {
			return fmt.Errorf("failed to creation cassandra session : %v", err)
		}
		defer sessionSystem.Close()

		if err = sessionSystem.Query(`CREATE KEYSPACE IF NOT EXISTS cus 
			WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`).
			RetryPolicy(&gocql.SimpleRetryPolicy{}).Exec(); err != nil {
			return fmt.Errorf("failed to create keyspace : %v", err)
		}

		cluster := gocql.NewCluster(hostAndPort)
		cluster.Keyspace = "cus"
		cluster.Timeout = time.Millisecond * 2000
		cluster.DisableInitialHostLookup = true
		session, err := cluster.CreateSession()
		if err != nil {
			return fmt.Errorf("failed to create keyspace session : %v", err)
		}
		closeFunction = func() { session.Close() }
		store = NewCassadnraUniqueStore(session)
		return nil
	}); err != nil {
		panic(fmt.Errorf("Could not connect to docker : %v", err))
	}

	return
}
