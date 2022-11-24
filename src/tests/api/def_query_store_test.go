package api_test

import (
	"bytes"
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
