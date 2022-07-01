package ehr_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/storage"
)

const testStatus = "test_status"

func TestSave(t *testing.T) {
	jsonDoc := fakeData.EhrCreateRequest()

	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	docService := service.NewDefaultDocumentService(cfg)

	ehrService := ehr.NewService(docService)

	var ehrReq model.EhrCreateRequest

	err = json.Unmarshal(jsonDoc, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	testSubjectID := ehrReq.Subject.ExternalRef.ID.Value

	testSubjectNamespace := ehrReq.Subject.ExternalRef.Namespace

	testUserID := uuid.New().String()

	ehrDoc, err := ehrService.EhrCreate(testUserID, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	// Check that subject index is added
	ehrID, err := docService.SubjectIndex.GetEhrBySubject(testSubjectID, testSubjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if ehrID != ehrDoc.EhrID.Value {
		t.Errorf("Incorrect ehrID in SubjectIndex")
	}
}

func TestStatus(t *testing.T) {
	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	docService := service.NewDefaultDocumentService(cfg)
	service := ehr.NewService(docService)
	userID := uuid.New().String()
	subjectID1 := uuid.New().String()
	subjectNamespace := testStatus
	subjectID2 := uuid.New().String()

	newEhr, err := getNewEhr(docService, userID, subjectID1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value

	statusIDNew := uuid.New().String()

	statusNew, err := service.CreateStatus(userID, ehrID, statusIDNew, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	// get current EHR status

	statusGet, err := service.GetStatus(userID, ehrID)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet.UID.Value != statusNew.UID.Value {
		t.Fatalf("Expected %s, received %s", statusGet.UID.Value, statusNew.UID.Value)
	}

	// get status by subject
	statusGet2, err := service.GetStatusBySubject(userID, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet2.UID.Value != statusIDNew {
		t.Error("Got wrong status by subject")
	}
}

func TestStatusUpdate(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	docService := service.NewDefaultDocumentService(cfg)
	service := ehr.NewService(docService)
	userID := uuid.New().String()
	subjectNamespace := testStatus
	subjectID1 := uuid.New().String()
	statusID2 := uuid.New().String()
	subjectID2 := uuid.New().String()

	newEhr, err := getNewEhr(docService, userID, subjectID1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value

	statusNew2, err := service.CreateStatus(userID, ehrID, statusID2, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	err = service.SaveStatus(ehrID, userID, statusNew2)
	if err != nil {
		t.Fatal(err)
	}

	statusGet3, err := service.GetStatus(userID, ehrID)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet3.UID.Value != statusID2 {
		t.Error("Got wrong updated status")
	}
}

func getNewEhr(docService *service.DefaultDocumentService, userID, subjectID, subjectNamespace string) (*model.EHR, error) {
	var (
		service           = ehr.NewService(docService)
		createRequestByte = fakeData.EhrCreateCustomRequest(subjectID, subjectNamespace)
		createRequest     model.EhrCreateRequest
	)

	if err := json.Unmarshal(createRequestByte, &createRequest); err != nil {
		return nil, err
	}

	return service.EhrCreate(userID, &createRequest)
}

func TestGetStatusByNearestTime(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	docService := service.NewDefaultDocumentService(cfg)
	service := ehr.NewService(docService)
	userID := uuid.New().String()
	subjectID1 := uuid.New().String()
	subjectNamespace := testStatus
	subjectID2 := uuid.New().String()

	newEhr, err := getNewEhr(docService, userID, subjectID1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value
	statusIDNew := uuid.New().String()

	_, err = service.CreateStatus(userID, ehrID, statusIDNew, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	// Test: docIndex is not exist yet
	if _, err := service.GetStatusByNearestTime(userID, ehrID, time.Now(), types.EhrStatus); err != nil {
		t.Fatal("Should return status", err)
	}
}
