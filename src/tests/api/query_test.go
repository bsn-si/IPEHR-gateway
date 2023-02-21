package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
)

func (testWrap *testWrap) queryExecSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		t.Skip()

		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		opts := fmt.Sprintf("&ehr_id=%s&q=%s&fetch=%s&offset=%s&%s",
			user.ehrID,
			url.QueryEscape(`SELECT
			   e/ehr_id/value
			FROM EHR e [ehr_id/value=$ehr_id]`),
			"10",
			"0",
			url.QueryEscape("ehr_id="+user.ehrID),
		)

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(testWrap.serverURL+"/v1/query/aql?%s", opts), nil)
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

		_, err = io.ReadAll(response.Body)
		if err != nil {
			// TODO compare with real result
			t.Errorf("Expected body content")
		}
	}
}

func (testWrap *testWrap) queryExecPostSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		url := testWrap.serverURL + "/v1/query/aql"

		req := `{
		  "q": "SELECT e/ehr_id/value, 
					   c/context/start_time/value as startTime, 
					   c/uid/value as cid, 
					   c/name 
				FROM EHR e [ehr_id/value=$ehr_id] 
				CONTAINS COMPOSITION c [openEHR-EHR-COMPOSITION.encounter.v1] 
					CONTAINS OBSERVATION obs [openEHR-EHR-OBSERVATION.blood_pressure.v1] 
				WHERE obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp",	
		  "offset": 0,
		  "fetch": 10,
		  "query_parameters": {
			"ehr_id": "` + user.ehrID + `",
			"systolic_bp": 140
		  }
		}`
		req = strings.ReplaceAll(req, "\n", "")
		req = strings.ReplaceAll(req, "\t", "")

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(req)))
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

		// TODO check queryResp content
		queryResp := &model.QueryResponse{}

		err = json.NewDecoder(response.Body).Decode(queryResp)
		if err != nil {
			t.Errorf("Query response unmarshaling error: %v", err)
		}
		defer response.Body.Close()
	}
}

func (testWrap *testWrap) queryExecPostFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		url := testWrap.serverURL + "/v1/query/aql"

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
