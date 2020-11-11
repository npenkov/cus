package cus

import "testing"

func TestDelete(t *testing.T) {
	data := []byte{0x67, 0x68, 0x69}
	dupId, err := store.Create("77", data)
	if err != nil || dupId != nil {
		t.Fatalf("Failed to create record : %v : %v", dupId, err)
	}

	fetchedData, err := store.Get("77")
	if err != nil || fetchedData == nil {
		t.Fatalf("Failed to create record : %v : %v", data, err)
	}
	for idx := range data {
		if data[idx] != fetchedData[idx] {
			t.Fatalf("Data does not match : %d : %d != %d", idx, data[idx], fetchedData[idx])
		}
	}

	err = store.Delete("77")
	if err != nil {
		t.Fatalf("Failed to delete record : %v", err)
	}

	fetchedData, err = store.Get("77")
	if err == nil || err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound but found : %v : %v", err, fetchedData)
	}

	// Create record with same data that was deleted
	dupId, err = store.Create("88", data)
	if err != nil || dupId != nil {
		t.Fatalf("Failed to create record : %v : %v", *dupId, err)
	}
}

func TestDeleteUnkown(t *testing.T) {
	err := store.Delete("666")
	if err == nil || err != ErrNotFound {
		t.Fatalf("Expected that the devil is not in the DB : %v", err)
	}
}
