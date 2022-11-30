package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/user/model"
)

func (testWrap *testWrap) userGroupCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		name := fakeData.GetRandomStringWithLength(10)
		description := fakeData.GetRandomStringWithLength(10)

		userGroup, _, err := userGroupCreate(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, name, description, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		testData.userGroups = append(testData.userGroups, userGroup)
	}
}

func userGroupCreate(userID, systemID, accessToken, baseURL, name, description string, client *http.Client) (*model.UserGroup, string, error) {
	userGroup := &model.UserGroup{
		Name:        name,
		Description: description,
	}

	data, _ := json.Marshal(userGroup)
	url := baseURL + "/v1/user/group"

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", systemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
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

	return userGroup, requestID, nil
}
