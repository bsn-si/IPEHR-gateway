package ehr_test

import (
	"context"
	"encoding/json"
	"os"
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
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/storage"
)

const testStatus = "test_status"

var (
	infra      *infrastructure.Infra
	docService *service.DefaultDocumentService
	ehrService *ehr.Service
)

func prepare(t *testing.T) {
	if infra != nil {
		return
	}

	cfgPath := os.Getenv("IPEHR_CONFIG_PATH")

	cfg, err := config.New(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	infra = infrastructure.New(cfg)

	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	infra.LocalStorage = storage.Storage()

	docService = service.NewDefaultDocumentService(cfg, infra)
	ehrService = ehr.NewService(docService)
}

func TestSave(t *testing.T) {
	prepare(t)

	jsonDoc := fakeData.EhrCreateRequest()

	var ehrReq model.EhrCreateRequest
	err := json.Unmarshal(jsonDoc, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	testSubjectID := ehrReq.Subject.ExternalRef.ID.Value

	testSubjectNamespace := ehrReq.Subject.ExternalRef.Namespace

	testUserID := uuid.New().String()

	ehrSystemID := ehrService.GetSystemID()

	ctx := context.Background()

	ehrDoc, err := ehrService.EhrCreate(ctx, testUserID, ehrSystemID, &ehrReq)
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
	prepare(t)

	userID := uuid.New().String()
	subjectID1 := uuid.New().String()
	subjectNamespace := testStatus
	subjectID2 := uuid.New().String()

	newEhr, err := getNewEhr(userID, subjectID1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrSystemID := ehrService.GetSystemID()
	ehrID := newEhr.EhrID.Value
	statusIDNew := uuid.New().String() + "::" + ehrSystemID.String() + "::1"
	ctx := context.Background()

	statusNew, err := ehrService.CreateStatus(ctx, userID, ehrID, statusIDNew, subjectID2, subjectNamespace, ehrSystemID)
	if err != nil {
		t.Fatal(err)
	}

	// get current EHR status

	ehrUUID, _ := uuid.Parse(ehrID)

	statusGet, err := ehrService.GetStatus(userID, &ehrUUID)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet.UID.Value != statusNew.UID.Value {
		t.Fatalf("Expected %s, received %s", statusGet.UID.Value, statusNew.UID.Value)
	}

	// get status by subject
	statusGet2, err := ehrService.GetStatusBySubject(userID, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet2.UID.Value != statusIDNew {
		t.Error("Got wrong status by subject")
	}
}

func TestStatusUpdate(t *testing.T) {
	prepare(t)

	ehrSystemID := ehrService.GetSystemID()
	userID := uuid.New().String()
	subjectNamespace := testStatus
	subjectID1 := uuid.New().String()
	statusID2 := uuid.New().String() + "::" + ehrSystemID.String() + "::1"
	subjectID2 := uuid.New().String()

	newEhr, err := getNewEhr(userID, subjectID1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value
	ctx := context.Background()

	statusNew2, err := ehrService.CreateStatus(ctx, userID, ehrID, statusID2, subjectID2, subjectNamespace, ehrSystemID)
	if err != nil {
		t.Fatal(err)
	}

	err = ehrService.SaveStatus(ctx, ehrID, userID, ehrSystemID, statusNew2)
	if err != nil {
		t.Fatal(err)
	}

	ehrUUID, _ := uuid.Parse(ehrID)

	statusGet3, err := ehrService.GetStatus(userID, &ehrUUID)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet3.UID.Value != statusID2 {
		t.Error("Got wrong updated status")
	}
}

func getNewEhr(userID, subjectID, subjectNamespace string) (*model.EHR, error) {
	var (
		ehrSystemID       = ehrService.GetSystemID()
		createRequestByte = fakeData.EhrCreateCustomRequest(subjectID, subjectNamespace)
		createRequest     model.EhrCreateRequest
	)

	if err := json.Unmarshal(createRequestByte, &createRequest); err != nil {
		return nil, err
	}

	return ehrService.EhrCreate(context.Background(), userID, ehrSystemID, &createRequest)
}

func TestGetStatusByNearestTime(t *testing.T) {
	prepare(t)

	ehrSystemID := ehrService.GetSystemID()
	userID := uuid.New().String()
	subjectID1 := uuid.New().String()
	subjectNamespace := testStatus
	subjectID2 := uuid.New().String()

	newEhr, err := getNewEhr(userID, subjectID1, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value
	statusIDNew := uuid.New().String() + "::" + ehrSystemID.String() + "::1"
	ctx := context.Background()

	_, err = ehrService.CreateStatus(ctx, userID, ehrID, statusIDNew, subjectID2, subjectNamespace, ehrSystemID)
	if err != nil {
		t.Fatal(err)
	}

	// Test: docIndex is not exist yet
	ehrUUID, _ := uuid.Parse(ehrID)

	if _, err := ehrService.GetStatusByNearestTime(userID, &ehrUUID, time.Now(), types.EhrStatus); err != nil {
		t.Fatal("Should return status", err)
	}
}
