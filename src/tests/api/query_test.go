package api_test

import (
	"bytes"
	"hms/gateway/pkg/common/fakeData"
	"net/http"
	"testing"
)

func (testWrap *testWrap) queryExecPostSuccess(testData *TestData) func(t *testing.T) {
	// TODO should be realize after AQL inserts will done
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(user.ehrID))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

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

func (testWrap *testWrap) queryExecPostFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("111qqqEEE")))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

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

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}
