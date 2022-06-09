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
	ehrDocService := NewEhrService(docService)
	statusService := NewEhrStatusService(docService)

	userId := uuid.New().String()
	subjectId1 := uuid.New().String()
	subjectNamespace := "test_status"
	subjectId2 := uuid.New().String()

	createRequestByte := fake_data.EhrCreateCustomRequest(subjectId1, subjectNamespace)
	var createRequest model.EhrCreateRequest
	err := json.Unmarshal(createRequestByte, &createRequest)
	if err != nil {
		t.Fatal(err)
	}

	newEhr, err := ehrDocService.Create(userId, &createRequest)
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
