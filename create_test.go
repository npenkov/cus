package cus

import "testing"

func TestCreate(t *testing.T) {
	dupId, err := store.Create("1", []byte{})
	if err != nil || dupId != nil {
		t.Fatalf("Failed to create record : %v : %v", dupId, err)
	}
}

func TestCreateDuplicate(t *testing.T) {
	data := []byte{0x64, 0x65, 0x66}

	dupId, err := store.Create("2", data)
	if err != nil || dupId != nil {
		t.Fatalf("Failed to create record : %v : %v", dupId, err)
	}
	dupId, err = store.Create("3", data)
	if err == nil || err != ErrDuplicateObject || dupId == nil || *dupId != "2" {
		t.Fatalf("Expected duplicate record, but creation was successful : %v : %v", dupId, err)
	}
}
