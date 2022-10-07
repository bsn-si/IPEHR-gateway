package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/api"
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/common/utils"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/storage"
	userRoles "hms/gateway/pkg/user/roles"
)

type testData struct {
	ehrSystemID  string
	subject      string
	namespace    string
	testUserID   string
	userPassword string
}

type testWrap struct {
	server     *httptest.Server
	httpClient *http.Client
	storage    *storage.Storager
}

type ehrContainer struct {
	ehr       *model.EHR
	requestID string
}

var ehrs = make(map[string]ehrContainer)
var usersGroupAccess = make(map[string]*model.GroupAccess)
var ehrsCompositions = make(map[string]*model.Composition)

func Test_API(t *testing.T) {
	var httpClient http.Client

	testServer, storager := prepareTest(t)

	testWrap := &testWrap{
		server:     testServer,
		httpClient: &httpClient,
		storage:    &storager,
	}
	defer tearDown(*testWrap)

	testData := &testData{
		ehrSystemID:  common.EhrSystemID,
		testUserID:   uuid.New().String(),
		userPassword: fakeData.GetRandomStringWithLength(10),
	}

	if !t.Run("User register", testWrap.userRegister(testData)) {
		t.Fatal()
	}
	// TODO user register incorrect input data
	// TODO user register duplicate registration request

	if !t.Run("EHR creating", testWrap.ehrCreate(testData)) {
		t.Fatal()
	}

	t.Run("Get transaction requests", testWrap.requests(testData))

	t.Run("EHR creating with id", testWrap.ehrCreateWithID(testData))
	t.Run("EHR creating with id for the same user", testWrap.ehrCreateWithIDForSameUser(testData))
	t.Run("EHR getting", testWrap.ehrGetByID(testData))
	t.Run("EHR get by subject", testWrap.ehrGetBySubject(testData))
	t.Run("EHR_STATUS getting", testWrap.ehrStatusGet(testData))
	t.Run("EHR_STATUS getting by version time", testWrap.ehrStatusGetByVersionTime(testData))

	if !t.Run("EHR_STATUS update", testWrap.ehrStatusUpdate(testData)) {
		t.Fatal()
	}

	t.Run("Access group create", testWrap.accessGroupCreate(testData))
	t.Run("Wrong access group getting", testWrap.wrongAccessGroupGetting(testData))
	t.Run("Access group getting", testWrap.accessGroupGetting(testData))
	t.Run("COMPOSITION create Expected fail with wrong EhrId", testWrap.compositionCreateFail(testData))

	if !t.Run("COMPOSITION create Expected success with correct EhrId", testWrap.compositionCreateSuccess(testData)) {
		t.Fatal()
	}

	t.Run("COMPOSITION getting with correct EhrId", testWrap.compositionGetByID(testData))
	t.Run("COMPOSITION getting with wrong EhrId", testWrap.compositionGetByWrongID(testData))
	t.Run("COMPOSITION update", testWrap.compositionUpdate(testData))
	t.Run("COMPOSITION delete by wrong UID", testWrap.compositionDeleteByWrongID(testData))
	t.Run("COMPOSITION delete", testWrap.compositionDeleteByID(testData))
	t.Run("QUERY execute with POST Expected success with correct query", testWrap.queryExecPostSuccess(testData))
	t.Run("QUERY execute with POST Expected fail with wrong query", testWrap.queryExecPostFail(testData))
}

func prepareTest(t *testing.T) (ts *httptest.Server, storager storage.Storager) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		t.Fatal("config.New error:", err)
	}

	cfg.Storage.Localfile.Path += "/test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	cfg.DefaultUserID = uuid.New().String()

	infra := infrastructure.New(cfg)
	r := api.New(cfg, infra).Build()
	ts = httptest.NewServer(r)

	return ts, storage.Storage()
}

func tearDown(testWrap testWrap) {
	testWrap.server.Close()

	err := (*testWrap.storage).Clean()
	if err != nil {
		log.Panicln(err)
	}
}

func (testWrap *testWrap) requests(testData *testData) func(t *testing.T) {
	_, requestID, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}
	return func(t *testing.T) {
		if requestID == "" {
			t.Fatal("Can not test because requestID is empty")
		}

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/requests/"+requestID, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GetRequestById expected %d, received %d", http.StatusOK, response.StatusCode)
		}

		t.Log("Requests: GetAll")

		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/requests/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("Prefer", "return=representation")

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GetAllRequests expected %d, received %d", http.StatusOK, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) userRegister(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		userRegisterRequest, err := userCreateBodyRequest(testData.testUserID, testData.userPassword)
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/user/register", userRegisterRequest)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(testData.testUserID, requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) ehrCreate(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr", ehrCreateBodyRequest())
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Fatal(err)
		}

		ehrID := response.Header.Get("ETag")
		if ehrID == "" {
			t.Fatal("EhrID missing")
		}
	}
}

func (testWrap *testWrap) ehrCreateWithID(testData *testData) func(t *testing.T) {
	testUserID := uuid.New().String()
	return func(t *testing.T) {
		ehrID2 := uuid.New().String()

		testData.subject = uuid.New().String()
		testData.namespace = "test_namespace"

		createRequest := fakeData.EhrCreateCustomRequest(testData.subject, testData.namespace)

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+ehrID2, bytes.NewReader(createRequest))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if err = response.Body.Close(); err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Fatal(err)
		}

		newEhrID := ehr.EhrID.Value
		if newEhrID != ehrID2 {
			t.Fatal("EhrID is not matched")
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(testUserID, requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) ehrCreateWithIDForSameUser(testData *testData) func(t *testing.T) {
	testUserID := uuid.New().String()
	testEhr, _, err := testWrap.createEhr(testUserID, testData.ehrSystemID)

	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+testEhr.EhrID.Value, ehrCreateBodyRequest())
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusConflict {
			t.Errorf("Expected %d, received %d", http.StatusConflict, response.StatusCode)
			return
		}
	}
}

func (testWrap *testWrap) ehrGetByID(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		testEhrID := testEhr.EhrID.Value

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testEhrID, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		if testEhrID != ehr.EhrID.Value {
			t.Error("EHR document mismatch")
			return
		}
	}
}

func (testWrap *testWrap) ehrGetBySubject(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	testEhrStatus, err := testWrap.getEhrStatus(testEhr.EhrID.Value, testEhr.EhrStatus.ID.Value, testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EhrStatus, received %s", err.Error())
	}

	return func(t *testing.T) {
		// Check document by subject
		url := testWrap.server.URL + "/v1/ehr?subject_id=" + testEhrStatus.Subject.ExternalRef.ID.Value + "&subject_namespace=" + testEhrStatus.Subject.ExternalRef.Namespace

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatal(err)
		}

		var ehrDoc model.EHR

		err = json.Unmarshal(data, &ehrDoc)
		if err != nil {
			t.Fatal(err)
		}

		if ehrDoc.EhrID.Value != testEhr.EhrID.Value {
			t.Fatalf("Expected %s, received %s", testEhr.EhrID.Value, ehrDoc.EhrID.Value)
		}
	}
}

func (testWrap *testWrap) ehrStatusGet(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		ehrID := testEhr.EhrID.Value
		statusID := testEhr.EhrStatus.ID.Value
		url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", ehrID, statusID)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Fatal(err)
		}

		if ehrStatus.UID == nil || ehrStatus.UID.Value != testEhr.EhrStatus.ID.Value {
			t.Fatal("EHR_STATUS document mismatch")
		}
	}
}

func (testWrap *testWrap) ehrStatusGetByVersionTime(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		ehrID := testEhr.EhrID.Value
		versionAtTime := time.Now()

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", ehrID), nil)
		if err != nil {
			t.Fatal(err)
		}

		q := request.URL.Query()
		q.Add("version_at_time", versionAtTime.Format(common.OpenEhrTimeFormat))

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.URL.RawQuery = q.Encode()

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d", http.StatusOK, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) ehrStatusUpdate(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		// replace substring in ehrStatusID
		ehrSystemID, _ := base.NewEhrSystemID(testData.ehrSystemID)
		objectVersionID, err := base.NewObjectVersionID(testEhr.EhrStatus.ID.Value, ehrSystemID)

		if err != nil {
			log.Fatalf("Expected model.EHR, received %s", err.Error())
		}

		_, err = objectVersionID.IncreaseUIDVersion()
		if err != nil {
			log.Fatalf("Expected model.EHR, received %s", err.Error())
		}

		newEhrStatusID := objectVersionID.String()

		req := []byte(fmt.Sprintf(`{
		  "_type": "EHR_STATUS",
		  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
		  "name": {
			"value": "EHR Status"
		  },
		  "uid": {
			"_type": "OBJECT_VERSION_ID",
			"value": "%s"
		  },
		  "subject": {
			"external_ref": {
			  "id": {
				"_type": "HIER_OBJECT_ID",
				"value": "324a4b23-623d-4213-cc1c-23f233b24234"
			  },
			  "namespace": "DEMOGRAPHIC",
			  "type": "PERSON"
			}
		  },
		  "other_details": {
			"_type": "ITEM_TREE",
			"archetype_node_id": "at0001",
			"name": {
			  "value": "Details"
			},
			"items": []
		  },
		  "is_modifiable": true,
		  "is_queryable": true
		}`, newEhrStatusID))

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", testEhr.EhrID.Value), bytes.NewReader(req))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("If-Match", testEhr.EhrStatus.ID.Value)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Fatal(err)
		}

		updatedEhrStatusID := response.Header.Get("ETag")

		if updatedEhrStatusID != newEhrStatusID {
			t.Log("Response body:", string(data))
			t.Fatal("EHR_STATUS uid in ETag mismatch")
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(testData.testUserID, requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}

		// Checking EHR_STATUS changes
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testEhr.EhrID.Value, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if err = response.Body.Close(); err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Fatal(err)
		}

		if ehr.EhrStatus.ID.Value != updatedEhrStatusID {
			t.Fatalf("EHR_STATUS id mismatch. Expected %s, received %s", updatedEhrStatusID, ehr.EhrStatus.ID.Value)
		}
	}
}

func (testWrap *testWrap) compositionCreateFail(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		userID := uuid.New().String()
		ehrID := uuid.New().String()

		body, err := compositionCreateBodyRequest()
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr/"+ehrID+"/composition", body)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusCreated {
			t.Fatalf("Expected error, received status: %d", response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionCreateSuccess(testData *testData) func(t *testing.T) {
	testUserID := uuid.New().String()

	testEhr, _, err := testWrap.createEhr(testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	testGroupAccess, err := testWrap.createGroupAccess(testUserID)
	if err != nil {
		log.Fatalf("Expected model.GroupAccess, received %s", err.Error())
	}

	testGroupAccessID := testGroupAccess.GroupUUID.String()

	return func(t *testing.T) {
		body, err := compositionCreateBodyRequest()
		if err != nil {
			t.Fatal(err)
		}

		url := testWrap.server.URL + "/v1/ehr/" + testEhr.EhrID.Value + "/composition"

		request, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserID)
		request.Header.Set("GroupAccessId", testGroupAccessID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected success, received status: %d", response.StatusCode)
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var c model.Composition
		if err = json.Unmarshal(data, &c); err != nil {
			t.Fatal(err)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(testUserID, requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) compositionGetByID(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	testGroupAccess, err := testWrap.createGroupAccess(testData.testUserID)

	if err != nil {
		log.Fatalf("Expected model.GroupAccess, received %s", err.Error())
	}

	testGroupAccessID := testGroupAccess.GroupUUID.String()

	testCreateComposition, err := testWrap.createComposition(testEhr, testGroupAccess, testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.Composition, received %s", err.Error())
	}

	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testEhr.EhrID.Value+"/composition/"+testCreateComposition.UID.Value, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("GroupAccessId", testGroupAccessID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status: %v, received %v", http.StatusOK, response.StatusCode)
		}

		var composition model.Composition
		if err = json.Unmarshal(data, &composition); err != nil {
			t.Fatal(err)
		}

		if composition.UID.Value != testCreateComposition.UID.Value {
			t.Fatalf("Expected %s, received %s", composition.UID.Value, testCreateComposition.UID.Value)
		}
	}
}

func (testWrap *testWrap) compositionGetByWrongID(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	_, err = testWrap.createGroupAccess(testData.testUserID)
	if err != nil {
		log.Fatalf("Expected model.GroupAccess, received %s", err.Error())
	}

	return func(t *testing.T) {
		wrongCompositionID := uuid.NewString() + "::" + testData.ehrSystemID + "::1"

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testEhr.EhrID.Value+"/composition/"+wrongCompositionID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status %d, received %d", http.StatusNotFound, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionUpdate(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	testGroupAccess, err := testWrap.createGroupAccess(testData.testUserID)

	if err != nil {
		log.Fatalf("Expected model.GroupAccess, received %s", err.Error())
	}

	testGroupAccessID := testGroupAccess.GroupUUID.String()

	testCreateComposition, err := testWrap.createComposition(testEhr, testGroupAccess, testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.Composition, received %s", err.Error())
	}

	return func(t *testing.T) {
		ehrSystemID, _ := base.NewEhrSystemID(testData.ehrSystemID)
		objectVersionID, err := base.NewObjectVersionID(testCreateComposition.UID.Value, ehrSystemID)

		if err != nil {
			t.Fatal(err)
		}

		testCreateComposition.ObjectVersionID = *objectVersionID

		testCreateComposition.Name.Value = "Updated text"
		updatedComposition, _ := json.Marshal(testCreateComposition)

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+testEhr.EhrID.Value+"/composition/"+testCreateComposition.ObjectVersionID.BasedID(), bytes.NewReader(updatedComposition))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("GroupAccessId", testGroupAccessID)
		request.Header.Set("If-Match", testCreateComposition.ObjectVersionID.String())
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status: %v, received %v", http.StatusOK, response.StatusCode)
		}

		compositionUpdated := model.Composition{}
		if err = json.Unmarshal(data, &compositionUpdated); err != nil {
			t.Fatal(err)
		}

		if compositionUpdated.UID.Value == testCreateComposition.UID.Value {
			t.Fatalf("Expected %s, received %s", compositionUpdated.UID.Value, testCreateComposition.UID.Value)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(testData.testUserID, requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) compositionDeleteByWrongID(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/ehr/" + testEhr.EhrID.Value + "/composition/" + uuid.New().String()

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status: %v, received %v", http.StatusNotFound, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionDeleteByID(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	testGroupAccess, err := testWrap.createGroupAccess(testData.testUserID)
	if err != nil {
		log.Fatalf("Expected model.GroupAccess, received %s", err.Error())
	}

	testCreateComposition, err := testWrap.createComposition(testEhr, testGroupAccess, testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.Composition, received %s", err.Error())
	}

	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/ehr/" + testEhr.EhrID.Value + "/composition/" + testCreateComposition.UID.Value

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status: %v, received %v", http.StatusNoContent, response.StatusCode)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(testData.testUserID, requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Checking the status of a re-request to remove")

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status: %v, received %v", http.StatusBadRequest, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) queryExecPostSuccess(testData *testData) func(t *testing.T) {
	testEhr, _, err := testWrap.createEhr(testData.testUserID, testData.ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(testEhr.EhrID.Value))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected success, received status: %d", response.StatusCode)
		}
	}
}

func (testWrap *testWrap) queryExecPostFail(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("111qqqEEE")))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected fail, received status: %d", response.StatusCode)
		}
	}
}

func userCreateBodyRequest(userID, password string) (*bytes.Reader, error) {
	userRegisterRequest := &model.UserCreateRequest{
		UserID:   userID,
		Password: password,
		Role:     uint8(userRoles.Patient),
	}

	docBytes, err := json.Marshal(userRegisterRequest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(docBytes), nil
}

func ehrCreateBodyRequest() *bytes.Reader {
	req := fakeData.EhrCreateRequest()
	return bytes.NewReader(req)
}

func compositionCreateBodyRequest() (*bytes.Reader, error) {
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return nil, err
	}

	filePath := rootDir + "/data/mock/ehr/composition.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func (testWrap *testWrap) accessGroupCreate(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		description := fakeData.GetRandomStringWithLength(50)

		req := []byte(`{
			"description": "` + description + `"
		}`)

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/access/group", bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusCreated, response.StatusCode, data)
		}

		var groupAccess model.GroupAccess
		if err = json.Unmarshal(data, &groupAccess); err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) wrongAccessGroupGetting(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		groupAccessIDWrong, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/access/group/"+groupAccessIDWrong.String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusNotFound, response.StatusCode, data)
		}
	}
}

func (testWrap *testWrap) accessGroupGetting(testData *testData) func(t *testing.T) {
	testGroupAccess, err := testWrap.createGroupAccess(testData.testUserID)

	if err != nil {
		log.Fatalf("Expected model.GroupAccess, received %s", err.Error())
	}

	testGroupAccessID := testGroupAccess.GroupUUID.String()

	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/access/group/"+testGroupAccessID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var groupAccessGot model.GroupAccess
		if err = json.Unmarshal(data, &groupAccessGot); err != nil {
			t.Fatal(err)
		}

		if testGroupAccessID != groupAccessGot.GroupUUID.String() {
			t.Fatal("Got wrong group")
		}
	}
}

func requestWait(userID, requestID string, tw *testWrap) error {
	request, err := http.NewRequest(http.MethodGet, tw.server.URL+"/v1/requests/"+requestID, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)

	timeout := time.Now().Add(2 * time.Minute)

	for {
		time.Sleep(3 * time.Second)

		if time.Now().After(timeout) {
			return errors.ErrTimeout
		}

		response, err := tw.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: request %s getting error: %v", errors.ErrCustom, requestID, response.Status)
		}

		var request processing.RequestResult

		if err = json.NewDecoder(response.Body).Decode(&request); err != nil {
			return err
		}

		if request.Ethereum[0].StatusStr == processing.StatusSuccess.String() {
			return nil
		} else if request.Ethereum[0].StatusStr == processing.StatusFailed.String() {
			return errors.New("Request failed")
		}
	}
}

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}

func (testWrap *testWrap) createEhr(userID, ehrSystemID string) (ehr *model.EHR, requestID string, err error) {
	key := userID + ehrSystemID
	if ehrCol, found := ehrs[key]; found {
		return ehrCol.ehr, ehrCol.requestID, nil
	}

	request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr", ehrCreateBodyRequest())
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := testWrap.httpClient.Do(request)
	if err != nil {
		return nil, "", err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	if err = json.Unmarshal(data, &ehr); err != nil {
		return nil, "", err
	}

	requestID = response.Header.Get("RequestId")
	err = requestWait(userID, requestID, testWrap)

	if err == nil {
		ehrs[key] = ehrContainer{
			ehr:       ehr,
			requestID: requestID,
		}
	}

	return ehr, requestID, err
}

func (testWrap *testWrap) getEhrStatus(ehrID, statusID, userID, ehrSystemID string) (*model.EhrStatus, error) {
	url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", ehrID, statusID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := testWrap.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var ehrStatus model.EhrStatus
	if err = json.Unmarshal(data, &ehrStatus); err != nil {
		return nil, err
	}

	return &ehrStatus, err
}

func (testWrap *testWrap) createGroupAccess(userID string) (*model.GroupAccess, error) {
	key := userID
	if userGroupAccess, found := usersGroupAccess[key]; found {
		return userGroupAccess, nil
	}

	description := fakeData.GetRandomStringWithLength(50)

	req := []byte(`{
			"description": "` + description + `"
		}`)

	request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/access/group", bytes.NewReader(req))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)

	response, err := testWrap.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(response.Status)
	}

	var groupAccess model.GroupAccess
	if err = json.Unmarshal(data, &groupAccess); err != nil {
		return nil, err
	}

	if err == nil {
		usersGroupAccess[key] = &groupAccess
	}

	return &groupAccess, nil
}

func (testWrap *testWrap) createComposition(testEhr *model.EHR, testGroupAccess *model.GroupAccess, userID, ehrSystemID string) (*model.Composition, error) {
	testGroupAccessID := testGroupAccess.GroupUUID.String()

	key := testEhr.EhrID.Value + userID + ehrSystemID
	if composition, found := ehrsCompositions[key]; found {
		return composition, nil
	}

	body, err := compositionCreateBodyRequest()
	if err != nil {
		return nil, err
	}

	url := testWrap.server.URL + "/v1/ehr/" + testEhr.EhrID.Value + "/composition"

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("GroupAccessId", testGroupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := testWrap.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var c model.Composition
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	requestID := response.Header.Get("RequestId")

	err = requestWait(userID, requestID, testWrap)
	if err != nil {
		return nil, err
	}

	if err == nil {
		ehrsCompositions[key] = &c
	}

	return &c, nil
}
