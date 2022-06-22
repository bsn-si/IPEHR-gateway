package data_access

import (
	"testing"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/keystore"
)

func TestDataAccessIndex(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	ks := keystore.New(cfg.KeystoreKey)
	dataAccessIndex := New(ks)

	userUUID := uuid.New()
	userId := userUUID.String()

	accessGroupUUID := uuid.New()
	accessGroupId := accessGroupUUID.String()

	accessGroupKey, err := fake_data.GetByteArray(32)
	if err != nil {
		t.Fatal(err)
	}

	err = dataAccessIndex.Add(userId, accessGroupId, accessGroupKey)
	if err != nil {
		t.Fatal("dataAccessIndex add error:", err)
	}

	groupAccessKey, err := dataAccessIndex.Get(userId, accessGroupId)
	if err != nil {
		t.Fatal("dataAccessIndex get error:", err)
	}

	if len(groupAccessKey) < 32 {
		t.Fatal("groupAccessKey incorrect")
	}

}
