package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

const version123 = "1.2.3"

func definitionStoreQuery(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		name := fakeData.GetRandomStringWithLength(10)
		version := ""

		storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, name, version, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		testData.storedQueries = append(testData.storedQueries, storedQuery)
	}
}

func definitionStoreInvalidQuery(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		name := fakeData.GetRandomStringWithLength(10)
		version := ""

		if err = storeInvalidQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, name, version, testData.httpClient); err != nil {
			if errors.Is(err, errors.ErrIsNotValid) {
				return
			}
		}

		t.Fatal(err)
	}
}

func definitionStoreQueryVersion(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		name := fakeData.GetRandomStringWithLength(10)
		version := version123

		storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, name, version, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		testData.storedQueries = append(testData.storedQueries, storedQuery)
	}
}

func definitionStoredQueryGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		name := fakeData.GetRandomStringWithLength(10)
		version := version123

		if len(testData.storedQueries) == 0 {
			storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, name, version, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			testData.storedQueries = append(testData.storedQueries, storedQuery)
		}

		var query1 *model.StoredQuery

		for _, q := range testData.storedQueries {
			if q.Version == version {
				query1 = q
				break
			}
		}

		if query1 == nil {
			t.Fatalf("Stored Query version '%s' not found in testData", version)
		}

		url := testData.serverURL + "/v1/definition/query/" + query1.Name + "/" + version

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
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

		var query2 model.StoredQuery

		err = json.Unmarshal(data, &query2)
		if err != nil {
			t.Fatal(err)
		}

		if query1.Name != query2.Name {
			t.Fatalf("Expected query name: %s, received: %s", query1.Name, query2.Name)
		}

		if query1.Query != query2.Query {
			t.Fatalf("Expected query content: %s, received: %s", query1.Query, query2.Query)
		}
	}
}

func definitionListStoredQueries(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(testData.storedQueries) == 0 {
			name := fakeData.GetRandomStringWithLength(10)
			version := "1.0.1"

			storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, name, version, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			testData.storedQueries = append(testData.storedQueries, storedQuery)
		}

		query1 := testData.storedQueries[0]

		url := testData.serverURL + "/v1/definition/query/" + query1.Name

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
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

func definitionStoreQueryVersionWithSameID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		var query *model.StoredQuery

		for _, q := range testData.storedQueries {
			if q.Name != "" && q.Version != "" {
				query = q
				break
			}
		}

		if query == nil {
			name := fakeData.GetRandomStringWithLength(10)
			version := version123

			storedQuery, reqID, err := storeQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, name, version, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}

			query = storedQuery
		}

		_, _, err = storeQuery(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, query.Name, query.Version, testData.httpClient)
		if err == nil || !errors.Is(err, errors.ErrAlreadyExist) {
			t.Fatalf("Expected error '%v', received: %v", errors.ErrAlreadyExist, err)
		}
	}
}

func storeQuery(userID, ehrSystemID, accessToken, baseURL, name, version string, client *http.Client) (*model.StoredQuery, string, error) {
	storedQuery := &model.StoredQuery{
		Name:    name,
		Type:    "AQL",
		Version: version,
	}

	storedQuery.Query = `SELECT c FROM
			EHR e
				CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.encounter.v1]
				  CONTAINS OBSERVATION obs[openEHR-EHR-OBSERVATION.blood_pressure.v1]
		 WHERE
			obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp`

	var url string

	switch version {
	case "":
		url = baseURL + "/v1/definition/query/" + storedQuery.Name + "?query_type=" + storedQuery.Type
	default:
		url = baseURL + "/v1/definition/query/" + storedQuery.Name + "/" + version + "?query_type=" + storedQuery.Type
	}

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
		if response.StatusCode == http.StatusConflict {
			return nil, "", errors.ErrAlreadyExist
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, "", err
		}

		return nil, "", errors.New(response.Status + " data: " + string(data))
	}

	requestID := response.Header.Get("RequestId")

	return storedQuery, requestID, nil
}

func storeInvalidQuery(userID, ehrSystemID, accessToken, baseURL, name, version string, client *http.Client) error {
	storedQuery := &model.StoredQuery{
		Name:    name,
		Type:    "AQL",
		Version: version,
	}

	storedQuery.Query = "SELECT invalid FROM"

	var url string

	switch version {
	case "":
		url = baseURL + "/v1/definition/query/" + storedQuery.Name + "?query_type=" + storedQuery.Type
	default:
		url = baseURL + "/v1/definition/query/" + storedQuery.Name + "/" + version + "?query_type=" + storedQuery.Type
	}

	data := bytes.NewReader([]byte(storedQuery.Query))

	request, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "text/plain")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return errors.Wrap(errors.ErrIsNotValid, response.Status+" data: "+string(data))
	}

	return nil
}
