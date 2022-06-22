package ehr

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
)

func TestSave(t *testing.T) {
	jsonDoc := fake_data.EhrCreateRequest()

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	docService := service.NewDefaultDocumentService(cfg)

	ehrService := NewEhrService(docService)

	var ehrReq model.EhrCreateRequest

	err = json.Unmarshal(jsonDoc, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	testSubjectId := ehrReq.Subject.ExternalRef.Id.Value
	testSubjectNamespace := ehrReq.Subject.ExternalRef.Namespace

	testUserId := uuid.New().String()

	ehrDoc, err := ehrService.EhrCreate(testUserId, &ehrReq)
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
