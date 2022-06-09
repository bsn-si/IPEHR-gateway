package data_search

import (
	"hms/gateway/pkg/common/fake_data"
	"testing"
)

func TestDataSearchIndex(t *testing.T) {
	dataSearchIndex := New()

	pathKey := fake_data.GetRandomStringWithLength(32)

	var groupId [16]byte
	randByteSlice, err := fake_data.GetByteArray(16)
	if err != nil {
		t.Fatal(err)
	}
	copy(groupId[:], randByteSlice)

	valueEncr, err := fake_data.GetByteArray(32)
	if err != nil {
		t.Fatal(err)
	}

	docStorEncr, err := fake_data.GetByteArray(32)
	if err != nil {
		t.Fatal(err)
	}

	dataEntry := &DataSearchEntry{
		GroupId:               &groupId,
		ValueEncrypted:        valueEncr,
		DocStorageIdEncrypted: docStorEncr,
	}

	err = dataSearchIndex.Add(pathKey, dataEntry)
	if err != nil {
		t.Fatal("dataSearchIndex add error:", err)
	}

	dataEntry, err = dataSearchIndex.Get(pathKey)
	if err != nil {
		t.Fatal("dataSearchIndex get error:", err)
	}

	if *dataEntry.GroupId != groupId {
		t.Fatal("dataSearchEntry from index mismatch")
	}

}
