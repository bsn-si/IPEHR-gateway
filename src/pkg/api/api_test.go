package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
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
)

type testData struct {
	ehrID         string
	ehrStatusID   string
	ehrSystemID   string
	subject       string
	namespace     string
	testUserID    string
	testUserID2   string
	ehrID2        string
	groupAccessID string
	compositionID string
	requestID     string
}

type testWrap struct {
	server     *httptest.Server
	httpClient *http.Client
	storage    *storage.Storager
}

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
		ehrSystemID: common.EhrSystemID,
		testUserID:  uuid.New().String(),
		testUserID2: uuid.New().String(),
	}

	if !t.Run("EHR creating", testWrap.ehrCreate(testData)) {
		t.Fatal()
	}

	t.Run("Get transaction requests", testWrap.requests(testData))

	if !t.Run("EHR creating with id", testWrap.ehrCreateWithID(testData)) {
		t.Fatal()
	}

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
	return func(t *testing.T) {
		if testData.requestID == "" {
			t.Fatal("Can not test because requestID is empty")
		}

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/requests/"+testData.requestID, nil)
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
			t.Fatalf("Expected %d, received %d", http.StatusOK, response.StatusCode)
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
			t.Fatalf("Expected %d, received %d", http.StatusOK, response.StatusCode)
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

		data, err := ioutil.ReadAll(response.Body)
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

		testData.requestID = response.Header.Get("RequestId")

		testData.ehrID = response.Header.Get("ETag")
		if testData.ehrID == "" {
			t.Fatal("EhrID missing")
		}
	}
}

func (testWrap *testWrap) ehrCreateWithID(testData *testData) func(t *testing.T) {
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
		request.Header.Set("AuthUserId", testData.testUserID2)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := ioutil.ReadAll(response.Body)
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

		testData.ehrID2 = ehrID2
		testData.requestID = response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", testData.requestID)

		err = requestWait(testData.testUserID2, testData.requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) ehrCreateWithIDForSameUser(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrID3 := uuid.New().String()

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+ehrID3, ehrCreateBodyRequest())
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID2)
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
	return func(t *testing.T) {
		ehrID := testData.ehrID

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+ehrID, nil)
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

		data, err := ioutil.ReadAll(response.Body)
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

		if ehrID != ehr.EhrID.Value {
			t.Error("EHR document mismatch")
			return
		}

		testData.ehrStatusID = ehr.EhrStatus.ID.Value
	}
}

func (testWrap *testWrap) ehrGetBySubject(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		// Check document by subject
		url := testWrap.server.URL + "/v1/ehr?subject_id=" + testData.subject + "&subject_namespace=" + testData.namespace

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID2)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := ioutil.ReadAll(response.Body)
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

		if ehrDoc.EhrID.Value != testData.ehrID2 {
			t.Fatalf("Expected %s, received %s", testData.ehrID2, ehrDoc.EhrID.Value)
		}
	}
}

func (testWrap *testWrap) ehrStatusGet(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrID := testData.ehrID
		statusID := testData.ehrStatusID
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

		data, err := ioutil.ReadAll(response.Body)
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

		if ehrStatus.UID == nil || ehrStatus.UID.Value != testData.ehrStatusID {
			t.Fatal("EHR_STATUS document mismatch")
		}
	}
}

func (testWrap *testWrap) ehrStatusGetByVersionTime(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrID := testData.ehrID
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
	return func(t *testing.T) {
		// replace substring in ehrStatusID
		newEhrStatusID := strings.Replace(testData.ehrStatusID, "::openEHRSys.example.com::1", "::openEHRSys.example.com::2", 1)

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

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", testData.ehrID), bytes.NewReader(req))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("If-Match", testData.ehrStatusID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := ioutil.ReadAll(response.Body)
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
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err = ioutil.ReadAll(response.Body)
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
	return func(t *testing.T) {
		body, err := compositionCreateBodyRequest()
		if err != nil {
			t.Fatal(err)
		}

		url := testWrap.server.URL + "/v1/ehr/" + testData.ehrID + "/composition"

		request, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("GroupAccessId", testData.groupAccessID)
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

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var c model.Composition
		if err = json.Unmarshal(data, &c); err != nil {
			t.Fatal(err)
		}

		testData.compositionID = c.UID.Value
		testData.requestID = response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", testData.requestID)

		err = requestWait(testData.testUserID, testData.requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) compositionGetByID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition/"+testData.compositionID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("GroupAccessId", testData.groupAccessID)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := ioutil.ReadAll(response.Body)
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

		if composition.UID.Value != testData.compositionID {
			t.Fatalf("Expected %s, received %s", composition.UID.Value, testData.compositionID)
		}
	}
}

func (testWrap *testWrap) compositionGetByWrongID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		wrongCompositionID := uuid.NewString() + "::openEHRSys.example.com::1"

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition/"+wrongCompositionID, nil)
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
	return func(t *testing.T) {
		body, err := compositionCreateBodyRequest()
		if err != nil {
			t.Fatal(err)
		}

		composition := model.Composition{}
		if err = composition.FromJSON(body); err != nil {
			t.Fatal(err)
		}

		//ehrSystemID := ehrService.Doc.GetSystemID()
		ehrSystemID, _ := base.NewEhrSystemID(testData.ehrSystemID)
		objectVersionID, err := base.NewObjectVersionID(composition.UID.Value, ehrSystemID)

		if err != nil {
			t.Fatal(err)
		}

		composition.ObjectVersionID = *objectVersionID

		composition.Name.Value = "Updated text"
		updatedComposition, _ := json.Marshal(composition)

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition/"+composition.ObjectVersionID.BasedID(), bytes.NewReader(updatedComposition))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("GroupAccessId", testData.groupAccessID)
		request.Header.Set("If-Match", composition.ObjectVersionID.String())
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := ioutil.ReadAll(response.Body)
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

		if err = json.Unmarshal(data, &composition); err != nil {
			t.Fatal(err)
		}

		if composition.UID.Value == testData.compositionID {
			t.Fatalf("Expected %s, received %s", composition.UID.Value, testData.compositionID)
		}

		testData.requestID = response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", testData.requestID)

		err = requestWait(testData.testUserID, testData.requestID, testWrap)
		if err != nil {
			t.Fatal(err)
		}

		testData.compositionID = composition.UID.Value
	}
}

func (testWrap *testWrap) compositionDeleteByWrongID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/ehr/" + testData.ehrID + "/composition/" + uuid.New().String()

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

		testData.requestID = response.Header.Get("RequestId")
	}
}

func (testWrap *testWrap) compositionDeleteByID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/ehr/" + testData.ehrID + "/composition/" + testData.compositionID

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

		testData.requestID = response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", testData.requestID)

		err = requestWait(testData.testUserID, testData.requestID, testWrap)
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
	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(testData.ehrID))
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

		data, err := ioutil.ReadAll(response.Body)
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

		testData.groupAccessID = groupAccess.GroupUUID.String()
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

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusNotFound, response.StatusCode, data)
		}
	}
}

func (testWrap *testWrap) accessGroupGetting(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/access/group/"+testData.groupAccessID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := ioutil.ReadAll(response.Body)
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

		if testData.groupAccessID != groupAccessGot.GroupUUID.String() {
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

		if request.Status == processing.StatusSuccess.String() {
			return nil
		} else if request.Status == processing.StatusFailed.String() {
			return errors.New("Request failed")
		}
	}
}

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}
