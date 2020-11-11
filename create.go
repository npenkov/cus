package cus

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const insertChecksumCQL = `INSERT INTO cus_checksum (checksum, obj_id) VALUES (?,?) IF NOT EXISTS`

const insertDataCQL = `INSERT INTO cus_data (obj_id, checksum, data) VALUES (?,?,?)`

func (s *CassandraUniqueStore) Create(id string, data []byte) (duplicateId *string, err error) {
	// Create checksum
	h := sha256.New()
	h.Write(data)
	checksum := hex.EncodeToString(h.Sum(nil))
	qry := s.session.Query(insertChecksumCQL, checksum, id)

	applied := true
	row := make(map[string]interface{})
	if applied, err = qry.MapScanCAS(row); err != nil {
		return nil, fmt.Errorf("failed no data available from LWT : %v", err)
	}
	if !applied {
		var ok bool
		var duplicateObjId string
		if duplicateObjId, ok = row["obj_id"].(string); !ok {
			return nil, fmt.Errorf("Failed to cast 2rd arg to string : %v", row)
		}
		return &duplicateObjId, ErrDuplicateObject
	}

	if err := s.session.Query(insertDataCQL, id, checksum, data).Exec(); err != nil {
		return nil, fmt.Errorf("failed to insert data : %v", err)
	}

	return nil, nil
}
