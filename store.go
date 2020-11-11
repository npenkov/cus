package cus

import "github.com/gocql/gocql"

type CassandraUniqueStore struct {
	session *gocql.Session
}

func NewCassadnraUniqueStore(session *gocql.Session) UniqueStore {
	return &CassandraUniqueStore{session: session}
}
