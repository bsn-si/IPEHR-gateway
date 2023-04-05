package api_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
)

func queryExecSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		targetURL := testData.serverURL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodGet, targetURL, nil)
		if err != nil {
			t.Error(err)
			return
		}

		q := request.URL.Query()

		q.Add("ehr_id", user.ehrID)

		query := url.QueryEscape(`SELECT e/ehr_id/value AS ID FROM EHR e [ehr_id/value=$ehr_id]`)
		q.Add("q", query)

		q.Add("fetch", "10")
		q.Add("offset", "0")
		q.Add("ehr_id", user.ehrID)
		request.URL.RawQuery = q.Encode()

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected success, received status: %d body: %s", response.StatusCode, body)
		}

		_, err = io.ReadAll(response.Body)
		if err != nil {
			// TODO compare with real result
			t.Errorf("Expected body content")
		}
	}
}

func queryExecPostSuccess(testData *TestData) func(t *testing.T) {
	// TODO should be realize after AQL inserts will done
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		url := testData.serverURL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(user.ehrID))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testData.httpClient.Do(request)
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

func queryExecPostFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		url := testData.serverURL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("111qqqEEE")))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testData.httpClient.Do(request)
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

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}
