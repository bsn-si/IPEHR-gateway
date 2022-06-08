package ehr_status

import (
	"github.com/google/uuid"
	"hms/gateway/pkg/docs/service"
	"testing"
)

func TestStatus(t *testing.T) {
	statusService := NewEhrStatusService(service.NewDefaultDocumentService())

	userId := uuid.New().String()
	ehrId := uuid.New().String()
	statusId := uuid.New().String()
	subjectId := uuid.New().String()
	subjectNamespace := "test_status"

	status := statusService.Create(statusId, subjectId, subjectNamespace)
	err := statusService.Save(ehrId, userId, status)
	if err != nil {
		t.Fatal(err)
	}

	// get current EHR status

	statusGet, err := statusService.Get(userId, ehrId)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet.Uid.Value != status.Uid.Value {
		t.Error("Got wrong status")
	}

	// get status by subject

	statusGet2, err := statusService.GetStatusBySubject(userId, subjectId, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet2.Uid.Value != status.Uid.Value {
		t.Error("Got wrong status by subject")
	}
}
