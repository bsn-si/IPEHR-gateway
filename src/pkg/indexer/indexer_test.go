package indexer

import (
	"hms/gateway/pkg/storage"
	"strconv"
	"testing"
	"time"
)

func TestIndex(t *testing.T) {
	sc := &storage.StorageConfig{}
	sc.New("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	name := "TestIndex"
	index := Init(name)

	type Person struct {
		Name    string
		Age     int
		Married bool
		Bytes   []byte
	}

	var item = Person{
		Name:    "John Doe",
		Age:     35,
		Married: true,
		Bytes:   []byte{1, 2, 3},
	}

	id := "123"
	err := index.Add(id, item)
	if err != nil {
		t.Error(err)
		return
	}

	var item2 Person
	if err = index.GetById(id, &item2); err != nil {
		t.Error(err)
		return
	}

	if item2.Name != item.Name {
		t.Errorf("name mismatch")
	}

	if len(item2.Bytes) != len(item.Bytes) {
		t.Errorf("bytes length mismatch")
	}

	if item2.Bytes[1] != item.Bytes[1] {
		t.Errorf("bytes[1] mismatch")
	}

	err = index.Delete(id)
	if err != nil {
		t.Error(err)
	}
}
