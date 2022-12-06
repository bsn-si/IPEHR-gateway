package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func (testWrap *testWrap) ehrCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		ehr, reqID, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		user.ehrID = ehr.EhrID.Value
		user.ehrStatusID = ehr.EhrStatus.ID.Value

		testData.requests = append(testData.requests, &Request{
			id:   reqID,
			kind: reqKindEhrCreate,
			user: user,
		})
	}
}

func (testWrap *testWrap) ehrCreateWithID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) < 2 {
			t.Fatal("Test user2 required")
		}

		user := testData.users[1]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		ehrID2 := uuid.New().String()

		ehr, _, err := createEhrWithID(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, ehrID2, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		newEhrID := ehr.EhrID.Value
		if newEhrID != ehrID2 {
			t.Fatal("EhrID is not matched")
		}
	}
}

func (testWrap *testWrap) ehrCreateWithIDForSameUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		_, _, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err == nil {
			t.Fatal("Expected error, received EHR")
		}
	}
}

func (testWrap *testWrap) ehrGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+user.ehrID, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
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

		if user.ehrID != ehr.EhrID.Value {
			t.Error("EHR document mismatch")
			return
		}
	}
}

func (testWrap *testWrap) ehrGetBySubject(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		testEhrStatus, err := testWrap.getEhrStatus(user.ehrID, user.ehrStatusID, user.id, testData.ehrSystemID, user.accessToken)
		if err != nil {
			log.Fatalf("Expected model.EhrStatus, received %s", err.Error())
		}

		// Check document by subject
		url := testWrap.server.URL + "/v1/ehr?subject_id=" + testEhrStatus.Subject.ExternalRef.ID.Value + "&subject_namespace=" + testEhrStatus.Subject.ExternalRef.Namespace

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("Prefer", "return=representation")

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

		if ehrDoc.EhrID.Value != user.ehrID {
			t.Fatalf("Expected %s, received %s", user.ehrID, ehrDoc.EhrID.Value)
		}
	}
}

func (testWrap *testWrap) ehrStatusGet(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", user.ehrID, user.ehrStatusID)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
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

		if ehrStatus.UID == nil || ehrStatus.UID.Value != user.ehrStatusID {
			t.Fatal("EHR_STATUS document mismatch")
		}
	}
}

func (testWrap *testWrap) ehrStatusGetByVersionTime(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]
		versionAtTime := time.Now()

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", user.ehrID), nil)
		if err != nil {
			t.Fatal(err)
		}

		q := request.URL.Query()
		q.Add("version_at_time", versionAtTime.Format(common.OpenEhrTimeFormat))

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
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

func (testWrap *testWrap) ehrStatusUpdate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		// replace substring in ehrStatusID
		objectVersionID, err := base.NewObjectVersionID(user.ehrStatusID, testData.ehrSystemID)
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

		url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status", user.ehrID)

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(req))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("If-Match", user.ehrStatusID)
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

		err = requestWait(user.id, user.accessToken, requestID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		// Checking EHR_STATUS changes
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+user.ehrID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
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

func (testWrap *testWrap) getEhrStatus(ehrID, statusID, userID, ehrSystemID, accessToken string) (*model.EhrStatus, error) {
	url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", ehrID, statusID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
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

func createEhr(userID, ehrSystemID, accessToken, baseURL string, client *http.Client) (ehr *model.EHR, requestID string, err error) {
	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/ehr", ehrCreateBodyRequest())
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
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

	return ehr, requestID, nil
}

func createEhrWithID(userID, ehrSystemID, accessToken, baseURL, ehrID string, client *http.Client) (ehr *model.EHR, requestID string, err error) {
	request, err := http.NewRequest(http.MethodPut, baseURL+"/v1/ehr/"+ehrID, ehrCreateBodyRequest())
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
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

	return ehr, requestID, nil
}

func ehrCreateBodyRequest() *bytes.Reader {
	req := fakeData.EhrCreateRequest()
	return bytes.NewReader(req)
}
