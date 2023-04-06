package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"github.com/google/uuid"
)

func ehrCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := checkUser(testData)
		if err != nil {
			t.Fatal("checkUser error:", err)
		}

		user := testData.users[0]

		err = checkUserLogin(testData, user)
		if err != nil {
			t.Fatal("checkUserLogin error:", err)
		}

		ehr, reqID, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, "", testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
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

func ehrCreateWithID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) < 2 {
			t.Fatal("Test user2 required")
		}

		user := testData.users[1]

		err := checkUserLogin(testData, user)
		if err != nil {
			t.Fatal("checkUserLogin error:", err)
		}

		ehrID2 := uuid.New().String()

		ehr, _, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, ehrID2, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		newEhrID := ehr.EhrID.Value
		if newEhrID != ehrID2 {
			t.Fatal("EhrID is not matched")
		}
	}
}

func ehrCreateWithIDForSameUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		_, _, err = createEhr(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, "", testData.httpClient)
		if err == nil {
			t.Fatal("Expected error, received EHR")
		}
	}
}

func ehrGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		request, err := http.NewRequest(http.MethodGet, testData.serverURL+"/v1/ehr/"+user.ehrID, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
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

func ehrGetBySubject(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		testEhrStatus, err := getEhrStatus(testData)
		if err != nil {
			log.Fatalf("Expected model.EhrStatus, received %s", err.Error())
		}

		// Check document by subject
		url := testData.serverURL + "/v1/ehr?subject_id=" + testEhrStatus.Subject.ExternalRef.ID.Value + "&subject_namespace=" + testEhrStatus.Subject.ExternalRef.Namespace

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testData.httpClient.Do(request)
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

func ehrStatusGet(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		url := testData.serverURL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", user.ehrID, user.ehrStatusID)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
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

func ehrStatusGetByVersionTime(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		versionAtTime := time.Now()

		request, err := http.NewRequest(http.MethodGet, testData.serverURL+fmt.Sprintf("/v1/ehr/%s/ehr_status", user.ehrID), nil)
		if err != nil {
			t.Fatal(err)
		}

		q := request.URL.Query()
		q.Add("version_at_time", versionAtTime.Format(common.OpenEhrTimeFormat))

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.URL.RawQuery = q.Encode()

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d", http.StatusOK, response.StatusCode)
		}
	}
}

func ehrStatusUpdate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

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

		url := testData.serverURL + fmt.Sprintf("/v1/ehr/%s/ehr_status", user.ehrID)

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

		response, err := testData.httpClient.Do(request)
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

		err = requestWait(user.id, user.accessToken, requestID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		// Checking EHR_STATUS changes
		request, err = http.NewRequest(http.MethodGet, testData.serverURL+"/v1/ehr/"+user.ehrID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err = testData.httpClient.Do(request)
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

func getEhrStatus(testData *TestData) (*model.EhrStatus, error) {
	user, err := checkUser0LoggedInAndEhrCreated(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser0LoggedInAndEhrCreated error: %w", err)
	}

	url := testData.serverURL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", user.ehrID, user.ehrStatusID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("AuthUserId", user.id)
	request.Header.Set("Authorization", "Bearer "+user.accessToken)
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := testData.httpClient.Do(request)
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

func createEhr(userID, ehrSystemID, accessToken, baseURL, ehrID string, client *http.Client) (ehr *model.EHR, requestID string, err error) {
	var method string

	switch {
	case len(ehrID) == 0:
		method = http.MethodPost
	case len(ehrID) > 0:
		method = http.MethodPut
	}

	request, err := http.NewRequest(method, baseURL+"/v1/ehr/"+ehrID, ehrCreateBodyRequest())
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
