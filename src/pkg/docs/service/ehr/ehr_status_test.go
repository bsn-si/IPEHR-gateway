package ehr

import (
	"encoding/json"
	"github.com/google/uuid"
	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"testing"
)

func TestStatus(t *testing.T) {
	docService := service.NewDefaultDocumentService()
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

	statusGet, err := statusService.Get(userId, ehrId)
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
	docService := service.NewDefaultDocumentService()
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

	err = statusService.Save(ehrId, userId, statusNew2)
	if err != nil {
		t.Fatal(err)
	}

	statusGet3, err := statusService.Get(userId, ehrId)
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

	newEhr, err = ehrDocService.Create(userId, &createRequest)
	return
}
