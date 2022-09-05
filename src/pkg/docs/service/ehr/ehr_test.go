package ehr_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/storage"
)

const testStatus = "test_status"

// nolint
var (
	infra      *infrastructure.Infra
	docService *service.DefaultDocumentService
	ehrService *ehr.Service
)

// nolint
func prepare(t *testing.T) {
	if infra != nil {
		return
	}

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	infra = infrastructure.New(cfg)

	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10))
	storage.Init(sc)

	infra.LocalStorage = storage.Storage()

	docService = service.NewDefaultDocumentService(cfg, infra)
	ehrService = ehr.NewService(docService)
}

// nolint
func requestWait(reqID string, timeout time.Duration) error {
	t := time.Now().Add(timeout)

	for {
		if time.Now().After(t) {
			return errors.ErrTimeout
		}

		data, err := docService.Proc.GetRequest(reqID)
		if err != nil {
			return err
		}

		var request processing.RequestResult

		if err = json.Unmarshal(data, &request); err != nil {
			return err
		}

		if request.Status == processing.StatusSuccess.String() {
			return nil
		}

		time.Sleep(3 * time.Second)
	}
}

func TestSave(t *testing.T) {
	t.Skip()

	prepare(t)

	var (
		testUserID  = uuid.New().String()
		ehrSystemID = ehrService.GetSystemID()
		ctx         = &gin.Context{}
	)

	var ehrReq model.EhrCreateRequest

	jsonDoc := fakeData.EhrCreateRequest()
	if err := json.Unmarshal(jsonDoc, &ehrReq); err != nil {
		t.Fatal(err)
	}

	testSubjectID := ehrReq.Subject.ExternalRef.ID.Value
	testSubjectNamespace := ehrReq.Subject.ExternalRef.Namespace

	reqID := "test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
	ctx.Set("reqId", reqID)

	ehrDoc, err := ehrService.EhrCreate(ctx, testUserID, ehrSystemID, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	if err = requestWait(reqID, 1*time.Minute); err != nil {
		t.Fatal(err)
	}

	// Check that subject index is added
	ehrUUID, err := docService.Infra.Index.GetEhrUUIDBySubject(ctx, testSubjectID, testSubjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if ehrUUID.String() != ehrDoc.EhrID.Value {
		t.Fatalf("Expected %s, received %s", ehrDoc.EhrID.Value, ehrUUID.String())
	}
}

func TestStatus(t *testing.T) {
	t.Skip()

	prepare(t)

	var (
		userID           = uuid.New().String()
		ehrSystemID      = ehrService.GetSystemID()
		subjectID1       = uuid.New().String()
		subjectNamespace = testStatus
		ctx              = &gin.Context{}
	)

	var createRequest model.EhrCreateRequest

	createRequestByte := fakeData.EhrCreateCustomRequest(subjectID1, subjectNamespace)
	if err := json.Unmarshal(createRequestByte, &createRequest); err != nil {
		t.Fatal(err)
	}

	reqID := "test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
	ctx.Set("reqId", reqID)

	newEhr, err := ehrService.EhrCreate(ctx, userID, ehrSystemID, &createRequest)
	if err != nil {
		t.Fatal(err)
	}

	if err = requestWait(reqID, time.Minute); err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value
	ehrUUID, _ := uuid.Parse(ehrID)
	statusIDNew := uuid.New().String() + "::" + ehrSystemID.String() + "::1"

	reqID = "test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
	ctx.Set("reqId", reqID)

	subjectID2 := uuid.New().String()

	_, err = ehrService.CreateStatus(statusIDNew, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if err = requestWait(reqID, time.Minute); err != nil {
		t.Fatal(err)
	}

	// get current EHR status
	statusGet, err := ehrService.GetStatus(ctx, userID, &ehrUUID)
	if err != nil {
		t.Fatal(err)
	}

	if statusGet.UID.Value != statusIDNew {
		t.Fatalf("Expected %s, received %s", statusIDNew, statusGet.UID.Value)
	}
}

func TestGetStatusByNearestTime(t *testing.T) {
	prepare(t)

	var (
		ehrSystemID       = ehrService.GetSystemID()
		userID            = uuid.New().String()
		subjectID1        = uuid.New().String()
		subjectNamespace  = testStatus
		subjectID2        = uuid.New().String()
		createRequestByte = fakeData.EhrCreateCustomRequest(subjectID1, subjectNamespace)
		createRequest     model.EhrCreateRequest
		ctx               = &gin.Context{}
	)

	if err := json.Unmarshal(createRequestByte, &createRequest); err != nil {
		t.Fatal(err)
	}

	reqID := "test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
	ctx.Set("reqId", reqID)

	newEhr, err := ehrService.EhrCreate(ctx, userID, ehrSystemID, &createRequest)
	if err != nil {
		t.Fatal(err)
	}

	if err = requestWait(reqID, time.Minute); err != nil {
		t.Fatal(err)
	}

	ehrID := newEhr.EhrID.Value
	ehrUUID, _ := uuid.Parse(ehrID)
	statusIDNew := uuid.New().String() + "::" + ehrSystemID.String() + "::1"

	reqID = "test_" + strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
	ctx.Set("reqId", reqID)

	doc, err := ehrService.CreateStatus(statusIDNew, subjectID2, subjectNamespace)
	if err != nil {
		t.Fatal(err)
	}

	var (
		transactions = ehrService.MultiCallTx.New(infra.Index, docService.Proc, processing.TxSetEhrUser, "", reqID)
	)

	err = ehrService.SaveStatus(ctx, &transactions, userID, &ehrUUID, ehrSystemID, doc, true)
	if err != nil {
		t.Fatal(err)
	}

	err = transactions.Commit()
	if err != nil {
		t.Fatal(err)
	}

	if err = requestWait(reqID, time.Minute); err != nil {
		t.Fatal(err)
	}

	// Test: docIndex is not exist yet
	if _, err := ehrService.GetStatusByNearestTime(ctx, userID, &ehrUUID, time.Now(), types.EhrStatus); err != nil {
		t.Fatal("Should return status", err)
	}
}
