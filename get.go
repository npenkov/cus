package cus

import (
	"fmt"

	"github.com/gocql/gocql"
)

const selectDataCQL = `SELECT data FROM cus_data WHERE obj_id = ?`

func (s *CassandraUniqueStore) Get(id string) (data []byte, err error) {
	if err := s.session.Query(selectDataCQL, id).Scan(&data); err != nil {
		if err == gocql.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch data : %v", err)
	}
	return data, nil
}
