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
	"hms/gateway/pkg/storage"
)

type testData struct {
	ehrID         string
	ehrStatusID   string
	testUserID    string
	testUserID2   string
	groupAccessID string
	compositionID string
}

type testWrap struct {
	server     *httptest.Server
	httpClient *http.Client
	storage    *storage.Storager
	api        *API
}

func Test_API(t *testing.T) {
	var httpClient http.Client
	testServer, storager, api := prepareTest(t)

	testWrap := &testWrap{
		server:     testServer,
		httpClient: &httpClient,
		storage:    &storager,
		api:        api,
	}
	defer tearDown(*testWrap)

	testData := &testData{
		testUserID:  uuid.New().String(),
		testUserID2: uuid.New().String(),
	}

	t.Run("EHR creating", testWrap.ehrCreate(testData))
	t.Run("EHR creating with id", testWrap.ehrCreateWithID(testData))
	t.Run("EHR creating with id for the same user", testWrap.ehrCreateWithIDForSameUser(testData))
	t.Run("EHR getting", testWrap.ehrGetByID(testData))
	t.Run("EHR_STATUS getting", testWrap.ehrStatusGet(testData))
	t.Run("EHR_STATUS getting by version time", testWrap.ehrStatusGetByVersionTime(testData))
	t.Run("EHR_STATUS update", testWrap.ehrStatusUpdate(testData))
	t.Run("EHR get by subject", testWrap.ehrGetBySubject())
	t.Run("EHR get by subject", testWrap.ehrGetBySubject(testData))
	t.Run("COMPOSITION create Expected fail with wrong EhrId", testWrap.compositionCreateFail(testData))
	t.Run("COMPOSITION create Expected success with correct EhrId", testWrap.compositionCreateSuccess(testData))
	t.Run("COMPOSITION getting with correct EhrId", testWrap.compositionGetById(testData))
	t.Run("COMPOSITION getting with wrong EhrId", testWrap.compositionGetByWrongId(testData))
	t.Run("COMPOSITION update", testWrap.compositionUpdate(testData))
	t.Run("COMPOSITION delete by wrong UID", testWrap.compositionDeleteByWrongId(testData))
	t.Run("COMPOSITION delete", testWrap.compositionDeleteById(testData))
	t.Run("QUERY execute with POST Expected success with correct query", testWrap.queryExecPostSuccess(testData))
	t.Run("QUERY execute with POST Expected fail with wrong query", testWrap.queryExecPostFail(testData))
	t.Run("Access group create", testWrap.accessGroupCreate(testData))
	t.Run("Wrong access group getting", testWrap.wrongAccessGroupGetting(testData))
	t.Run("Access group getting", testWrap.accessGroupGetting(testData))
	t.Run("COMPOSITION create Expected fail with wrong EhrId", testWrap.compositionCreateFail())
	t.Run("COMPOSITION create Expected success with correct EhrId", testWrap.compositionCreateSuccess(testData))
	t.Run("COMPOSITION getting with correct EhrId", testWrap.compositionGetByID(testData))
	t.Run("COMPOSITION getting with wrong EhrId", testWrap.compositionGetByWrongID(testData))
	t.Run("COMPOSITION delete by wrong UID", testWrap.compositionDeleteByWrongID(testData))
	t.Run("COMPOSITION delete", testWrap.compositionDeleteByID(testData))
	t.Run("QUERY execute with POST Expected success with correct query", testWrap.queryExecPostSuccess(testData))
	t.Run("QUERY execute with POST Expected fail with wrong query", testWrap.queryExecPostFail(testData))
}

func prepareTest(t *testing.T) (ts *httptest.Server, storager storage.Storager, api *API) {
	t.Helper()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	cfg.StoragePath += "/test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	api = api.New(cfg)
	r := api.Build()
	ts = httptest.NewServer(r)

	return ts, storage.Storage(), api
}

func tearDown(testWrap testWrap) {
	testWrap.server.Close()

	err := (*testWrap.storage).Clean()
	if err != nil {
		log.Panicln(err)
	}
}

func (testWrap *testWrap) ehrCreate(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr", ehrCreateBodyRequest())
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

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
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
			t.Error(err)
			return
		}

		testData.ehrID = response.Header.Get("ETag")
		if testData.ehrID == "" {
			t.Error("EhrID missing")
			return
		}
	}
}

func (testWrap *testWrap) ehrCreateWithID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrID2 := uuid.New().String()

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+ehrID2, ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID2)
		request.Header.Set("Prefer", "return=representation")

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

func (testWrap *testWrap) ehrStatusGet(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", testData.ehrID, testData.ehrStatusID), nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testData.testUserID)

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
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("If-Match", testData.ehrStatusID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
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
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Error(err)
			return
		}

		updatedEhrStatusID := response.Header.Get("ETag")

		if updatedEhrStatusID != newEhrStatusID {
			t.Log("Response body:", string(data))
			t.Fatal("EHR_STATUS uid in ETag mismatch")
		}

		// Checking EHR_STATUS changes
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)

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
			return
		}

		if ehr.EhrStatus.ID.Value != updatedEhrStatusID {
			t.Fatalf("EHR_STATUS id mismatch. Expected %s, received %s", updatedEhrStatusID, ehr.EhrStatus.ID.Value)
			return
		}
	}
}

func (testWrap *testWrap) ehrGetBySubject() func(t *testing.T) {
	return func(t *testing.T) {
		// Adding document with specific subject
		userID := uuid.New().String()
		ehrID := uuid.New().String()

		subjectID := uuid.New().String()
		subjectNamespace := "test_test"

		createRequest := fakeData.EhrCreateCustomRequest(subjectID, subjectNamespace)

		url := testWrap.server.URL + "/v1/ehr/" + ehrID

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(createRequest))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
		}

		// Check document by subject
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr?subject_id="+subjectID+"&subject_namespace="+subjectNamespace, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userID)
		request.Header.Set("Prefer", "return=representation")

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatal(err)
		}

		var ehrDoc model.EHR

		err = json.Unmarshal(data, &ehrDoc)
		if err != nil {
			t.Fatal(err)
		}

		if ehrDoc.EhrID.Value != ehrID {
			t.Error("Got wrong EHR")
		}
	}
}

func (testWrap *testWrap) compositionCreateFail() func(t *testing.T) {
	return func(t *testing.T) {
		userID := uuid.New().String()
		ehrID := uuid.New().String()

		body, err := compositionCreateBodyRequest()
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr/"+ehrID+"/composition", body)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusCreated {
			t.Errorf("Expected error, received status: %d", response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionCreateSuccess(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		body, err := compositionCreateBodyRequest()
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition", body)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserID)
		request.Header.Set("GroupAccessId", testData.groupAccessID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			t.Errorf("Expected success, received status: %d", response.StatusCode)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var c model.Composition
		if err = json.Unmarshal(data, &c); err != nil {
			t.Error(err)
			return
		}

		testData.compositionID = c.UID.Value
	}
}

func (testWrap *testWrap) compositionGetByID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition/"+testData.compositionID, nil)
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

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrId+"/composition/"+testData.compositionId, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserId)

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

		var composition model.Composition
		if err = json.Unmarshal(data, &composition); err != nil {
			t.Fatal(err)
		}

		// TODO composition.Uid.Value - should it be equal with versioned_object_uid?
		composition.Name.Value = "Updated text"
		updatedComposition, _ := json.Marshal(composition)

		uid := testWrap.api.Composition.service.GetObjectVersionIdByUid(testData.compositionId)

		request, err = http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+testData.ehrId+"/composition/"+uid.ObjectId(), bytes.NewReader(updatedComposition))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("If-Match", uid.String())

		request.Header.Set("AuthUserId", testData.testUserId)

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err = ioutil.ReadAll(response.Body)
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

		if composition.Uid.Value == testData.compositionId {
			t.Fatalf("Expected %s, received %s", composition.Uid.Value, testData.compositionId)
		}

		testData.compositionId = composition.Uid.Value
	}
}

func (testWrap *testWrap) compositionDeleteByWrongID(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodDelete, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition/"+uuid.New().String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)

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
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodDelete, testWrap.server.URL+"/v1/ehr/"+testData.ehrID+"/composition/"+testData.compositionID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status: %v, received %v", http.StatusNoContent, response.StatusCode)
		}

		// Check answer for repeated query
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

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}
