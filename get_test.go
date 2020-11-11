package cus

import "testing"

func TestGet(t *testing.T) {
	data := []byte{0x66, 0x67, 0x68}
	dupId, err := store.Create("55", data)
	if err != nil || dupId != nil {
		t.Fatalf("Failed to create record : %v : %v", dupId, err)
	}

	fetchedData, err := store.Get("55")
	if err != nil || fetchedData == nil {
		t.Fatalf("Failed to create record : %v : %v", data, err)
	}
	for idx := range data {
		if data[idx] != fetchedData[idx] {
			t.Fatalf("Data does not match : %d : %d != %d", idx, data[idx], fetchedData[idx])
		}
	}
}

func TestGetUnkown(t *testing.T) {
	fetchedData, err := store.Get("666")
	if err == nil || fetchedData != nil || err != ErrNotFound {
		t.Fatalf("Expected that the devil is not in the DB : %v : %v", fetchedData, err)
	}
}
