package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"

	"github.com/google/uuid"
)

func (testWrap *testWrap) definitionStoreQuery(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			user := &User{
				id:       uuid.New().String(),
				password: fakeData.GetRandomStringWithLength(10),
			}

			reqID, err := registerUser(user, testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatalf("Can not register user, err: %v", err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("registerUser requestWait error: ", err)
			}

			testData.users = append(testData.users, user)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		if user.ehrID == "" {
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
		}

		storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		testData.storedQueries = append(testData.storedQueries, storedQuery)
	}
}

func (testWrap *testWrap) definitionListStoredQueries(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			user := &User{
				id:       uuid.New().String(),
				password: fakeData.GetRandomStringWithLength(10),
			}

			reqID, err := registerUser(user, testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatalf("Can not register user, err: %v", err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("registerUser requestWait error: ", err)
			}

			testData.users = append(testData.users, user)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		if user.ehrID == "" {
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
		}

		if len(testData.storedQueries) == 0 {
			storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			testData.storedQueries = append(testData.storedQueries, storedQuery)
		}

		query1 := testData.storedQueries[0]

		url := testWrap.server.URL + "/v1/definition/query/" + query1.Name

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var storedQueries []model.StoredQuery

		err = json.Unmarshal(data, &storedQueries)
		if err != nil {
			t.Fatal(err)
		}

		if len(storedQueries) == 0 {
			t.Fatalf("Expected query in list, received: %s", string(data))
		}

		query2 := storedQueries[0]

		if query1.Name != query2.Name {
			t.Fatalf("Expected query name: %s, received: %s", query1.Name, query2.Name)
		}

		if query1.Query != query2.Query {
			t.Fatalf("Expected query content: %s, received: %s", query1.Query, query2.Query)
		}
	}
}

func storeQuery(userID, ehrSystemID, accessToken, baseURL string, client *http.Client) (*model.StoredQuery, string, error) {
	storedQuery := &model.StoredQuery{
		Name:    fakeData.GetRandomStringWithLength(10),
		Type:    "AQL",
		Version: "1.0.1",
	}

	storedQuery.Query = `SELECT c FROM
			EHR e
				CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.encounter.v1]
				  CONTAINS OBSERVATION obs[openEHR-EHR-OBSERVATION.blood_pressure.v1]
		 WHERE
			obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp`

	url := baseURL + "/v1/definition/query/" + storedQuery.Name + "?query_type=" + storedQuery.Type

	data := bytes.NewReader([]byte(storedQuery.Query))

	request, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "text/plain")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, "", err
		}
		return nil, "", errors.New(string(data))
	}

	requestID := response.Header.Get("RequestId")

	return storedQuery, requestID, nil
}
