package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

const (
	reqKindEhrCreate = iota
	//reqKindUserRegister
)

type Request struct {
	id   string
	kind int
	user *User
}

func requests(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		var req *Request

	loop:
		for _, r := range testData.requests {
			switch r.kind {
			case reqKindEhrCreate:
				req = r
				break loop
			default:
			}
		}

		if req == nil {
			t.Fatal("Request required")
		}

		if req.user.accessToken == "" {
			err := req.user.login(testData.ehrSystemID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		request, err := http.NewRequest(http.MethodGet, testData.serverURL+"/v1/requests/"+req.id, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", req.user.id)
		request.Header.Set("Authorization", "Bearer "+req.user.accessToken)
		request.Header.Set("Prefer", "return=representation")

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GetRequestById expected %d, received %d", http.StatusOK, response.StatusCode)
		}

		t.Log("Requests: GetAll")

		request, err = http.NewRequest(http.MethodGet, testData.serverURL+"/v1/requests/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", req.user.id)
		request.Header.Set("Authorization", "Bearer "+req.user.accessToken)
		request.Header.Set("Prefer", "return=representation")

		response, err = testData.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GetAllRequests expected %d, received %d", http.StatusOK, response.StatusCode)
		}
	}
}

func requestWait(userID, accessToken, requestID, baseURL string, client *http.Client) error {
	request, err := http.NewRequest(http.MethodGet, baseURL+"/v1/requests/"+requestID, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)

	if accessToken != "" {
		request.Header.Set("Authorization", "Bearer "+accessToken)
	}

	timeout := time.Now().Add(2 * time.Minute)

	for {
		time.Sleep(2 * time.Second)

		if time.Now().After(timeout) {
			return errors.ErrTimeout
		}

		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: request %s getting error: %v", errors.ErrCustom, requestID, response.Status)
		}

		var request processing.RequestResult

		if err = json.NewDecoder(response.Body).Decode(&request); err != nil {
			return err
		}

		if request.Ethereum[0].StatusStr == processing.StatusSuccess.String() {
			return nil
		} else if request.Ethereum[0].StatusStr == processing.StatusFailed.String() {
			return errors.New("Request failed")
		}
	}
}
