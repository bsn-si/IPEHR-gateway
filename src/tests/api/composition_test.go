package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func compositionCreateFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		ehrID := uuid.New().String()
		groupAccessID := ""

		composition, _, err := createComposition(user.id, ehrID, testData.ehrSystemID, user.accessToken, groupAccessID, testData.serverURL, testData.httpClient)
		if err == nil {
			t.Fatalf("Expected error, received status: %v", composition)
		}
	}
}

func compositionCreateSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(testData.groupsAccess) == 0 {
			uuid := uuid.New()

			testData.groupsAccess = append(testData.groupsAccess, &model.GroupAccess{GroupUUID: &uuid})
		}

		ga := testData.groupsAccess[0]

		c, reqID, err := createComposition(user.id, user.ehrID, testData.ehrSystemID, user.accessToken, ga.GroupUUID.String(), testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatalf("Unexpected composition, received error: %v", err)
		}

		t.Logf("Waiting for request %s done", reqID)

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		user.compositions = append(user.compositions, c)
	}
}

func compositionGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(testData.groupsAccess) == 0 {
			uuid := uuid.New()

			testData.groupsAccess = append(testData.groupsAccess, &model.GroupAccess{GroupUUID: &uuid})
		}

		ga := testData.groupsAccess[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		comp := user.compositions[0]

		url := testData.serverURL + "/v1/ehr/" + user.ehrID + "/composition/" + comp.UID.Value

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("GroupAccessId", ga.GroupUUID.String())
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
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

		if composition.UID.Value != comp.UID.Value {
			t.Fatalf("Expected %s, received %s", composition.UID.Value, comp.UID.Value)
		}
	}
}

func compositionGetByWrongID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		wrongCompositionID := uuid.NewString() + "::" + testData.ehrSystemID + "::1"

		url := testData.serverURL + "/v1/ehr/" + user.ehrID + "/composition/" + wrongCompositionID

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status %d, received %d", http.StatusNotFound, response.StatusCode)
		}
	}
}

func compositionGetList(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(user.compositions) == 0 {
			c, reqID, err := createComposition(user.id, user.ehrID, testData.ehrSystemID, user.accessToken, "", testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatalf("Unexpected composition, received error: %v", err)
			}

			t.Logf("Waiting for request %s done", reqID)

			err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			user.compositions = append(user.compositions, c)
		}

		url := testData.serverURL + "/v1/ehr/" + user.ehrID + "/composition"

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status %d, received %d", http.StatusOK, response.StatusCode)
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		var list []model.EhrDocumentItem

		err = json.Unmarshal(data, &list)
		if err != nil {
			t.Fatalf("Response body unmarshal error: %v", err)
		}

		if len(list) == 0 {
			t.Fatalf("Expected len(list) > 0, received %d", len(list))
		}
	}
}

func compositionUpdate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(testData.groupsAccess) == 0 {
			uuid := uuid.New()

			testData.groupsAccess = append(testData.groupsAccess, &model.GroupAccess{GroupUUID: &uuid})
		}

		ga := testData.groupsAccess[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		comp := user.compositions[0]

		objectVersionID, err := base.NewObjectVersionID(comp.UID.Value, testData.ehrSystemID)
		if err != nil {
			t.Fatal(err)
		}

		comp.ObjectVersionID = *objectVersionID

		comp.Name.Value = "Updated text"
		updatedComposition, _ := json.Marshal(comp)

		url := testData.serverURL + "/v1/ehr/" + user.ehrID + "/composition/" + comp.ObjectVersionID.BasedID()

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(updatedComposition))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("GroupAccessId", ga.GroupUUID.String())
		request.Header.Set("If-Match", comp.ObjectVersionID.String())
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
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

		compositionUpdated := model.Composition{}
		if err = json.Unmarshal(data, &compositionUpdated); err != nil {
			t.Fatal(err)
		}

		if compositionUpdated.UID.Value == comp.UID.Value {
			t.Fatalf("Expected %s, received %s", compositionUpdated.UID.Value, comp.UID.Value)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(user.id, user.accessToken, requestID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func compositionDeleteByWrongID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		url := testData.serverURL + "/v1/ehr/" + user.ehrID + "/composition/" + uuid.New().String()

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status: %v, received %v", http.StatusNotFound, response.StatusCode)
		}
	}
}

func compositionDeleteByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		comp := user.compositions[0]

		url := testData.serverURL + "/v1/ehr/" + user.ehrID + "/composition/" + comp.UID.Value

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status: %v, received %v", http.StatusNoContent, response.StatusCode)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(user.id, user.accessToken, requestID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Checking the status of a re-request to remove")

		response, err = testData.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status: %v, received %v", http.StatusBadRequest, response.StatusCode)
		}
	}
}

func compositionCreateBodyRequest(ehrSystemID string) (*bytes.Reader, error) {
	filePath := "./test_fixtures/composition_1.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	compositionID := uuid.New().String()

	objectVersionID, err := base.NewObjectVersionID(compositionID, ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	data = []byte(strings.Replace(string(data), "__COMPOSITION_ID__", objectVersionID.String(), 1))

	return bytes.NewReader(data), nil
}

// nolint
func createComposition(userID, ehrID, ehrSystemID, accessToken, groupAccessID, baseURL string, client *http.Client) (*model.Composition, string, error) {
	body, err := compositionCreateBodyRequest(ehrSystemID)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannnot create composition body request")
	}

	url := baseURL + "/v1/ehr/" + ehrID + "/composition"

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	//request.Header.Set("GroupAccessId", groupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot do create composition request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "connot read response body")
	}

	var c model.Composition
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, "", errors.Wrap(err, "cannot unmarshal COMPOSITION mondel")
	}

	requestID := response.Header.Get("RequestId")

	return &c, requestID, nil
}
