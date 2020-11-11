package cus

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

const createDataTableCQL = `CREATE TABLE IF NOT EXISTS cus_data(
	obj_id TEXT PRIMARY KEY,
	checksum TEXT,
	data BLOB)`
const dropDataTableCQL = `DROP TABLE IF EXISTS cus_data;`

const createChecksumTableCQL = `CREATE TABLE IF NOT EXISTS cus_checksum(
	checksum TEXT PRIMARY KEY,
	obj_id TEXT)`
const dropChecksumTableCQL = `DROP TABLE IF EXISTS cus_checksum;`

func (s *CassandraUniqueStore) CreateSchema(force bool) error {
	ctx := context.TODO()

	if force {
		if err := s.execCreateSQL(ctx, dropDataTableCQL); err != nil {
			return fmt.Errorf("failed to drop table : %v", err)
		}
		if err := s.execCreateSQL(ctx, dropChecksumTableCQL); err != nil {
			return fmt.Errorf("failed to drop table : %v", err)
		}
	}
	if err := s.execCreateSQL(ctx, createDataTableCQL); err != nil {
		return fmt.Errorf("failed to create table : %v", err)
	}
	if err := s.execCreateSQL(ctx, createChecksumTableCQL); err != nil {
		return fmt.Errorf("failed to create table : %v", err)
	}

	return nil
}

func (s *CassandraUniqueStore) execCreateSQL(ctx context.Context, cqlStatement string) error {
	if err := s.session.AwaitSchemaAgreement(ctx); err != nil {
		return fmt.Errorf("failed to wait for schema agreement : %v", err)
	}
	if err := s.session.Query(cqlStatement).RetryPolicy(&gocql.SimpleRetryPolicy{}).Exec(); err != nil {
		return fmt.Errorf("error executing statement %s err=%v\n", cqlStatement, err)
	}
	return nil
}
