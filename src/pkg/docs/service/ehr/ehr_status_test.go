package ehr

import (
	"encoding/json"
	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/storage"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStatus(t *testing.T) {
	sc := &storage.StorageConfig{}
	sc.New("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	docService := service.NewDefaultDocumentService()
	statusService := NewEhrStatusService(docService, nil)

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
	docService := service.NewDefaultDocumentService()
	statusService := NewEhrStatusService(docService, nil)

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
	ehrDocService := NewEhrService(docService, nil)

	createRequestByte := fake_data.EhrCreateCustomRequest(subjectId, subjectNamespace)
	var createRequest model.EhrCreateRequest
	err = json.Unmarshal(createRequestByte, &createRequest)
	if err != nil {
		return
	}

	newEhr, err = ehrDocService.EhrCreate(userId, &createRequest)
	return
}
