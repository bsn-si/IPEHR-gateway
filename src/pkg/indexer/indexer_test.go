package indexer_test

import (
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/processing"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/infrastructure"
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
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	infra := infrastructure.New(cfg)

	index := indexer.New(cfg.Contract.Address, cfg.Contract.PrivKeyPath, infra.EthClient)
	docService := service.NewDefaultDocumentService(cfg, infra)

	userID := fakeData.GetRandomStringWithLength(16)
	ehrUUID := uuid.New()

	t.Logf("userID %s ehrUUID %s", userID, ehrUUID.String())

	ctx := &gin.Context{}
	reqID := "test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
	ctx.Set("reqId", reqID)

	var tx = docService.MultiCallTx.New(index, docService.Proc, processing.TxSetEhrUser, "indexer_test", reqID)

	procReq := &processing.Request{
		ReqID:        reqID,
		Kind:         processing.RequestEhrCreate,
		Status:       processing.StatusProcessing,
		UserID:       userID,
		EhrUUID:      ehrUUID.String(),
		CID:          "",
		DealCID:      "dealCID",
		MinerAddress: "minerAddr",
	}

	if err = docService.Proc.AddRequest(procReq); err != nil {
		t.Fatal("Proc.AddRequest error: %w", err)
	}

	packed, err := index.SetEhrUser(userID, &ehrUUID)
	if err != nil {
		t.Fatal(err)
	}

	tx.Add(packed)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(30 * time.Second)

	ehrUUID2, err := index.GetEhrUUIDByUserID(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	if ehrUUID != *ehrUUID2 {
		t.Fatalf("ehrUUID expected %s, received %s", ehrUUID.String(), ehrUUID2.String())
	}
}
