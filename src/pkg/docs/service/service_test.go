package service

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"testing"
	"time"
)

func TestGetDocIndexByNearestTime(t *testing.T) {

	docService := NewDefaultDocumentService()
	ehrId := docService.GenerateId()

	// Test: docIndex is not exist yet
	if _, err := docService.GetDocIndexByNearestTime(ehrId, time.Now(), types.EHR_STATUS); err == nil {
		t.Error("DocService contains indexes")
	}

	if err := fillDocIndexes(5, ehrId, docService); err != nil {
		t.Error("DocService contains indexes")
	}

	docIndexes, err := docService.DocsIndex.Get(ehrId)
	if err != nil {
		t.Error(err)
	}

	lastDocIndex := &docIndexes[len(docIndexes)-1]
	lastDocIndexTime := time.Unix(int64((*lastDocIndex).Timestamp), 0)

	// Test: resulted docIndex should be last one if the specified time is equal with last dateIndex time value
	docIndex, err := docService.GetDocIndexByNearestTime(ehrId, lastDocIndexTime, types.EHR_STATUS)
	if err != nil || docIndex == nil {
		t.Error("DocService not contains indexes")
	}

	// Test: resulted docIndex should be last one again if the specified time is greater that exist
	DocIndexTimeMoreThanExist := lastDocIndexTime.Add(time.Hour)
	docIndex, err = docService.GetDocIndexByNearestTime(ehrId, DocIndexTimeMoreThanExist, types.EHR_STATUS)
	if err != nil || docIndex == nil {
		t.Error("DocService not contains indexes")
	}

	// Test: resulted docIndex should be nil if the specified time is less among existing
	firstDocIndex := &docIndexes[0]
	firstDocIndexTime := time.Unix(int64((*firstDocIndex).Timestamp), 0)

	DocIndexTimeLessThanExist := firstDocIndexTime.Add(-24 * time.Hour)
	docIndex, err = docService.GetDocIndexByNearestTime(ehrId, DocIndexTimeLessThanExist, types.EHR_STATUS)
	if err == nil || docIndex != nil {
		t.Error("docIndex should not be nil")
	}
}

func fillDocIndexes(count int, ehrId string, service *DefaultDocumentService) (err error) {
	msg := []byte("test")

	if count == 0 {
		panic("Count is not set")
	}

	timeFrom, _ := time.Parse("2006-01-02", "2020-01-01")

	for i := 0; i < count; i++ {
		docStorageId, err := service.Storage.Add(msg)
		if err != nil {
			panic("cannot get docStorageId")
		}

		docIndex := model.DocumentMeta{
			TypeCode:  types.EHR_STATUS,
			StorageId: docStorageId,
			Timestamp: uint32(timeFrom.Unix()),
		}

		err = service.DocsIndex.Add(ehrId, &docIndex)
		if err != nil {
			return nil
		}

		timeFrom = timeFrom.Add(24 * time.Hour)
	}

	return
}
