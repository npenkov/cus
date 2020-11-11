package cus

import (
	"fmt"

	"github.com/gocql/gocql"
)

const selectCheckumFromIDCQL = `SELECT checksum FROM cus_data WHERE obj_id = ?`
const deleteDataCQL = `DELETE FROM cus_data WHERE obj_id = ?`
const deleteChecksumCQL = `DELETE FROM cus_checksum WHERE checksum = ?`

func (s *CassandraUniqueStore) Delete(id string) (err error) {
	var checksum string
	if err := s.session.Query(selectCheckumFromIDCQL, id).Scan(&checksum); err != nil {
		if err == gocql.ErrNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to fetch data : %v", err)
	}
	if err := s.session.Query(deleteDataCQL, id).Exec(); err != nil {
		return fmt.Errorf("failed to delete data : %v", err)
	}
	if err := s.session.Query(deleteChecksumCQL, checksum).Exec(); err != nil {
		return fmt.Errorf("failed to delete checksum : %v", err)
	}
	return nil
}
