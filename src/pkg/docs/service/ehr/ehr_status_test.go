package ehr

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
)

func TestStatus(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	docService := service.NewDefaultDocumentService(cfg)
	statusService := NewEhrStatusService(docService)

	userId := uuid.New().String()
	subjectId1 := uuid.New().String()
	subjectNamespace := "test_status"
	subjectId2 := uuid.New().String()

	newEhr, err := getNewEhr(docService, userId, subjectId1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrId := newEhr.EhrId.Value

	statusIdNew := uuid.New().String()

	statusNew, err := statusService.Create(userId, ehrId, statusIdNew, subjectId2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	// get current EHR status

	statusGet, err := statusService.GetStatus(userId, ehrId)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet.Uid.Value != statusNew.Uid.Value {
		t.Error("Got wrong status")
	}

	// get status by subject
	statusGet2, err := statusService.GetStatusBySubject(userId, subjectId2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet2.Uid.Value != statusIdNew {
		t.Error("Got wrong status by subject")
	}
}

func TestStatusUpdate(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	docService := service.NewDefaultDocumentService(cfg)
	statusService := NewEhrStatusService(docService)

	userId := uuid.New().String()
	subjectNamespace := "test_status"
	subjectId1 := uuid.New().String()
	statusId2 := uuid.New().String()
	subjectId2 := uuid.New().String()

	newEhr, err := getNewEhr(docService, userId, subjectId1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrId := newEhr.EhrId.Value

	statusNew2, err := statusService.Create(userId, ehrId, statusId2, subjectId2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	err = statusService.SaveStatus(ehrId, userId, statusNew2)
	if err != nil {
		t.Fatal(err)
	}

	statusGet3, err := statusService.GetStatus(userId, ehrId)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet3.Uid.Value != statusId2 {
		t.Error("Got wrong updated status")
	}
}

func getNewEhr(docService *service.DefaultDocumentService, userId, subjectId, subjectNamespace string) (newEhr *model.EHR, err error) {
	ehrDocService := NewEhrService(docService)

	createRequestByte := fake_data.EhrCreateCustomRequest(subjectId, subjectNamespace)
	var createRequest model.EhrCreateRequest
	err = json.Unmarshal(createRequestByte, &createRequest)
	if err != nil {
		return
	}

	newEhr, err = ehrDocService.EhrCreate(userId, &createRequest)
	return
}

func TestGetStatusByNearestTime(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	docService := service.NewDefaultDocumentService(cfg)
	statusService := NewEhrStatusService(docService)

	userId := uuid.New().String()
	subjectId1 := uuid.New().String()
	subjectNamespace := "test_status"
	subjectId2 := uuid.New().String()

	newEhr, err := getNewEhr(docService, userId, subjectId1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrId := newEhr.EhrId.Value

	statusIdNew := uuid.New().String()

	_, err = statusService.Create(userId, ehrId, statusIdNew, subjectId2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	// Test: docIndex is not exist yet
	if _, err := statusService.GetStatusByNearestTime(userId, ehrId, time.Now(), types.EHR_STATUS); err != nil {
		t.Fatal("Should return status", err)
	}

}
