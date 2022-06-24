package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/storage"
)

type testData struct {
	ehrId         string
	ehrStatusId   string
	testUserId    string
	testUserId2   string
	groupAccessId string
	compositionId string
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
		testUserId:  uuid.New().String(),
		testUserId2: uuid.New().String(),
	}

	t.Run("EHR creating", testWrap.ehrCreate(testData))
	t.Run("EHR creating with id", testWrap.ehrCreateWithId(testData))
	t.Run("EHR creating with id for the same user", testWrap.ehrCreateWithIdForSameUser(testData))
	t.Run("EHR getting", testWrap.ehrGetById(testData))
	t.Run("EHR_STATUS getting", testWrap.ehrStatusGet(testData))
	t.Run("EHR_STATUS getting by version time", testWrap.ehrStatusGetByVersionTime(testData))
	t.Run("EHR_STATUS update", testWrap.ehrStatusUpdate(testData))
	t.Run("EHR get by subject", testWrap.ehrGetBySubject(testData))
	t.Run("COMPOSITION create Expected fail with wrong EhrId", testWrap.compositionCreateFail(testData))
	t.Run("COMPOSITION create Expected success with correct EhrId", testWrap.compositionCreateSuccess(testData))
	t.Run("COMPOSITION getting with correct EhrId", testWrap.compositionGetById(testData))
	t.Run("COMPOSITION getting with wrong EhrId", testWrap.compositionGetByWrongId(testData))
	t.Run("COMPOSITION delete by wrong UID", testWrap.compositionDeleteByWrongId(testData))
	t.Run("COMPOSITION delete", testWrap.compositionDeleteById(testData))
	t.Run("QUERY execute with POST Expected success with correct query", testWrap.queryExecPostSuccess(testData))
	t.Run("QUERY execute with POST Expected fail with wrong query", testWrap.queryExecPostFail(testData))
	t.Run("Access group create", testWrap.accessGroupCreate(testData))
	t.Run("Wrong access group getting", testWrap.wrongAccessGroupGetting(testData))
	t.Run("Access group getting", testWrap.accessGroupGetting(testData))
}

func prepareTest(t *testing.T) (ts *httptest.Server, storager storage.Storager) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	cfg.StoragePath += "/test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	r := New(cfg).Build()
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

func (testWrap *testWrap) ehrCreate(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr/", ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserId)
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
			t.Errorf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		testData.ehrId = response.Header.Get("ETag")
		if testData.ehrId == "" {
			t.Error("EhrId missing")
			return
		}
	}
}

func (testWrap *testWrap) ehrCreateWithId(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrId2 := uuid.New().String()

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+ehrId2, ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserId2)
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
			t.Errorf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		newEhrId := ehr.EhrId.Value
		if newEhrId != ehrId2 {
			t.Error("EhrId is not matched")
			return
		}
	}
}

func (testWrap *testWrap) ehrCreateWithIdForSameUser(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrId3 := uuid.New().String()

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+"/v1/ehr/"+ehrId3, ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserId2)

		response, err := testWrap.httpClient.Do(request)
		err = response.Body.Close()
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusConflict {
			t.Errorf("Expected %d, received %d", http.StatusConflict, response.StatusCode)
			return
		}
	}
}

func (testWrap *testWrap) ehrGetById(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrId := testData.ehrId
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+ehrId, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testData.testUserId)

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

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		if ehrId != ehr.EhrId.Value {
			t.Error("EHR document mismatch")
			return
		}

		testData.ehrStatusId = ehr.EhrStatus.Id.Value
	}
}

func (testWrap *testWrap) ehrStatusGet(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", testData.ehrId, testData.ehrStatusId), nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testData.testUserId)

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

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Error(err)
			return
		}

		if ehrStatus.Uid == nil || ehrStatus.Uid.Value != testData.ehrStatusId {
			t.Error("EHR_STATUS document mismatch")
			return
		}

	}
}

func (testWrap *testWrap) ehrStatusGetByVersionTime(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		ehrId := testData.ehrId
		versionAtTime := time.Now()

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", ehrId), nil)
		if err != nil {
			t.Error(err)
			return
		}

		q := request.URL.Query()
		q.Add("version_at_time", versionAtTime.Format(common.OPENEHR_TIME_FORMAT))

		request.Header.Set("AuthUserId", testData.testUserId)
		request.URL.RawQuery = q.Encode()

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(response.Body)

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d", http.StatusOK, response.StatusCode)
			return
		}
	}
}

func (testWrap *testWrap) ehrStatusUpdate(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		// replace substring in ehrStatusId
		newEhrStatusId := strings.Replace(testData.ehrStatusId, "::openEHRSys.example.com::1", "::openEHRSys.example.com::2", 1)

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
		}`, newEhrStatusId))

		request, err := http.NewRequest(http.MethodPut, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", testData.ehrId), bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserId)
		request.Header.Set("If-Match", testData.ehrStatusId)
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

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Error(err)
			return
		}

		updatedEhrStatusId := response.Header.Get("ETag")

		if updatedEhrStatusId != newEhrStatusId {
			t.Log("Response body:", string(data))
			t.Error("EHR_STATUS uid in ETag mismatch")
			return
		}

		// Checking EHR_STATUS changes
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrId, nil)
		if err != nil {
			t.Fatal(err)
		}
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
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Fatal(err)
			return
		}

		if ehr.EhrStatus.Id.Value != updatedEhrStatusId {
			t.Fatalf("EHR_STATUS id mismatch. Expected %s, received %s", updatedEhrStatusId, ehr.EhrStatus.Id.Value)
			return
		}
	}
}

func (testWrap *testWrap) ehrGetBySubject(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		// Adding document with specific subject
		userId := uuid.New().String()
		ehrId := uuid.New().String()

		subjectId := uuid.New().String()
		subjectNamespace := "test_test"

		createRequest := fake_data.EhrCreateCustomRequest(subjectId, subjectNamespace)

		url := testWrap.server.URL + "/v1/ehr/" + ehrId
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(createRequest))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userId)
		request.Header.Set("Prefer", "return=representation")

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
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
		}

		// Check document by subject
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/?subject_id="+subjectId+"&subject_namespace="+subjectNamespace, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userId)
		request.Header.Set("Prefer", "return=representation")

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}

		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = response.Body.Close()
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

		if ehrDoc.EhrId.Value != ehrId {
			t.Error("Got wrong EHR")
		}
	}
}

func (testWrap *testWrap) compositionCreateFail(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		userId := uuid.New().String()
		ehrId := uuid.New().String()

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr/"+ehrId+"/composition", compositionCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", userId)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		if response.StatusCode == http.StatusCreated {
			t.Errorf("Expected error, received status: %d", response.StatusCode)
		}

	}
}

func (testWrap *testWrap) compositionCreateSuccess(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/ehr/"+testData.ehrId+"/composition", compositionCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserId)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		if response.StatusCode != http.StatusCreated {
			t.Errorf("Expected success, received status: %d", response.StatusCode)
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		var c model.Composition
		if err = json.Unmarshal(data, &c); err != nil {
			t.Error(err)
			return
		}

		testData.compositionId = c.Uid.Value
	}
}

func (testWrap *testWrap) compositionGetById(testData *testData) func(t *testing.T) {
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

		if composition.Uid.Value != testData.compositionId {
			t.Fatalf("Expected %s, received %s", composition.Uid.Value, testData.compositionId)
		}
	}
}

func (testWrap *testWrap) compositionGetByWrongId(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		wrongCompositionId := uuid.NewString() + "::openEHRSys.example.com::1"
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+testData.ehrId+"/composition/"+wrongCompositionId, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserId)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status %d, received %d", http.StatusNotFound, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionDeleteByWrongId(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodDelete, testWrap.server.URL+"/v1/ehr/"+testData.ehrId+"/composition/"+uuid.New().String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserId)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status: %v, received %v", http.StatusNotFound, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionDeleteById(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodDelete, testWrap.server.URL+"/v1/ehr/"+testData.ehrId+"/composition/"+testData.compositionId, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testData.testUserId)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status: %v, received %v", http.StatusNoContent, response.StatusCode)
		}

		// Check answer for repeated query
		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		if response.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status: %v, received %v", http.StatusBadRequest, response.StatusCode)
		}

	}
}

func (testWrap *testWrap) queryExecPostSuccess(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		url := testWrap.server.URL + "/v1/query/aql"
		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(testData.ehrId))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testData.testUserId)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

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
		request.Header.Set("AuthUserId", testData.testUserId)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected fail, received status: %d", response.StatusCode)
		}
	}
}

func ehrCreateBodyRequest() *bytes.Reader {
	req := fake_data.EhrCreateRequest()
	return bytes.NewReader(req)
}

func compositionCreateBodyRequest() *bytes.Reader {
	req := fake_data.CompositionCreateRequest()
	return bytes.NewReader(req)
}

func (testWrap *testWrap) accessGroupCreate(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		description := fake_data.GetRandomStringWithLength(50)

		req := []byte(`{
			"description": "` + description + `"
		}`)

		request, err := http.NewRequest(http.MethodPost, testWrap.server.URL+"/v1/access/group", bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
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

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusCreated, response.StatusCode, data)
		}

		var groupAccess model.GroupAccess
		if err = json.Unmarshal(data, &groupAccess); err != nil {
			t.Fatal(err)
		}

		testData.groupAccessId = groupAccess.GroupId
	}
}

func (testWrap *testWrap) wrongAccessGroupGetting(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		groupAccessIdWrong, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/access/group/"+groupAccessIdWrong.String(), nil)
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

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusNotFound, response.StatusCode, data)
		}
	}
}

func (testWrap *testWrap) accessGroupGetting(testData *testData) func(t *testing.T) {
	return func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/access/group/"+testData.groupAccessId, nil)
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
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var groupAccessGot model.GroupAccess
		if err = json.Unmarshal(data, &groupAccessGot); err != nil {
			t.Fatal(err)
		}

		if testData.groupAccessId != groupAccessGot.GroupId {
			t.Fatal("Got wrong group")
		}
	}
}

func queryExecPostCreateBodyRequest(ehrId string) *bytes.Reader {
	req := fake_data.QueryExecRequest(ehrId)
	return bytes.NewReader(req)

}
