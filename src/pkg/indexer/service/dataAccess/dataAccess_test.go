package dataAccess_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/indexer/service/dataAccess"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/storage"
)

func TestDataAccessIndex(t *testing.T) {
	t.Skip()

	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	ks := keystore.New(cfg.KeystoreKey)
	dataAccessIndex := dataAccess.New(ks)

	userUUID := uuid.New()
	userID := userUUID.String()

	accessGroupUUID := uuid.New()
	accessGroupID := accessGroupUUID.String()

	accessGroupKey, err := fakeData.GetByteArray(32)
	if err != nil {
		t.Fatal(err)
	}

	err = dataAccessIndex.Add(userID, accessGroupID, accessGroupKey)
	if err != nil {
		t.Fatal("dataAccessIndex add error:", err)
	}

	groupAccessKey, err := dataAccessIndex.Get(userID, accessGroupID)
	if err != nil {
		t.Fatal("dataAccessIndex get error:", err)
	}

	if len(groupAccessKey) < 32 {
		t.Fatal("groupAccessKey incorrect")
	}
}
