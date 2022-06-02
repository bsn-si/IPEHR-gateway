package ehr

import (
	"encoding/json"
	"github.com/google/uuid"
	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"testing"
)

func TestSave(t *testing.T) {
	jsonDoc := fake_data.EhrCreateRequest()

	docService := service.NewDefaultDocumentService()

	ehrService := NewEhrService(docService)

	var ehrReq model.EhrCreateRequest

	err := json.Unmarshal(jsonDoc, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	testSubjectId := ehrReq.Subject.ExternalRef.Id.Value
	testSubjectNamespace := ehrReq.Subject.ExternalRef.Namespace

	ehrDoc := ehrService.Create(&ehrReq)

	testUserId := uuid.New().String()

	err = ehrService.Save(testUserId, ehrDoc)
	if err != nil {
		t.Fatal(err)
	}

	// Check that subject index is added
	ehrId, err := docService.SubjectIndex.GetEhrBySubject(testSubjectId, testSubjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if ehrId != ehrDoc.EhrId.Value {
		t.Errorf("Incorrect ehrId in SubjectIndex")
	}
}
