package indexer_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/storage"
)

func TestIndex(t *testing.T) {
	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	name := "TestIndex"
	index := indexer.Init(name)

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
	if err = index.GetByID(id, &item2); err != nil {
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

func TestEhrByUserIndex(t *testing.T) {
	t.Skip()

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	index := indexer.New(cfg.Contract.Address, cfg.Contract.Endpoint, cfg.Contract.PrivKeyPath)

	userID := fakeData.GetRandomStringWithLength(16)
	ehrUUID := uuid.New()

	t.Logf("userID %s ehrUUID %s", userID, ehrUUID.String())

	h, err := index.SetEhrUser(userID, &ehrUUID)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	txStatus, err := index.TxWait(ctx, h)
	if err != nil {
		t.Fatal(err)
	}

	if txStatus == 1 {
		t.Logf("tx %s Success", h)
	} else {
		t.Logf("tx %s Failed", h)
	}

	ehrUUID2, err := index.GetEhrUUIDByUserID(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	if ehrUUID != *ehrUUID2 {
		t.Fatalf("ehrUUID expected %s, received %s", ehrUUID.String(), ehrUUID2.String())
	}
}
