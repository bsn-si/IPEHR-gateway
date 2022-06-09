package data_access

import (
	"testing"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fake_data"
)

func TestDataAccessIndex(t *testing.T) {
	dataAccessIndex := New()

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
