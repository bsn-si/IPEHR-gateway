package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

func (testWrap *testWrap) userGroupCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		name := fakeData.GetRandomStringWithLength(10)
		description := fakeData.GetRandomStringWithLength(10)

		userGroup, _, err := userGroupCreate(user, testData.ehrSystemID, testWrap.serverURL, name, description, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		testData.userGroups = append(testData.userGroups, userGroup)
	}
}

func (testWrap *testWrap) userGroupAddUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		addingUser := testData.users[1]

		err = checkUserGroup(user, testData, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		userGroup := testData.userGroups[0]

		reqID, err := userGroupAddUser(user, addingUser, userGroup, testData, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatal("userGroupAddUser error: ", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatal("requestWait error: %w reqID: %w", err, reqID)
		}

		userGroup.Members = append(userGroup.Members, addingUser.id)
	}
}

func (testWrap *testWrap) userGroupRemoveUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = checkUserGroup(user, testData, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		userGroup := testData.userGroups[0]

		if len(userGroup.Members) == 0 {
			addingUser := testData.users[1]

			reqID, err := userGroupAddUser(user, addingUser, userGroup, testData, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal("userGroupAddUser error: ", err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal("requestWait error: %w reqID: %w", err, reqID)
			}
		}

		removingUserID := userGroup.Members[0]

		url := testWrap.serverURL + "/v1/user/group/" + userGroup.GroupID.String() + "/user_remove/" + removingUserID

		request, err := http.NewRequest(http.MethodPost, url, nil)
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
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		// Checking that the user has been deleted
		url = testWrap.serverURL + "/v1/user/group/" + userGroup.GroupID.String()

		request, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err = io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var userGroup2 model.UserGroup

		err = json.Unmarshal(data, &userGroup2)
		if err != nil {
			t.Fatal(err)
		}

		for _, id := range userGroup2.Members {
			if id == removingUserID {
				t.Fatalf("Expected: userID %s removed, received: still in the group", removingUserID)
			}
		}

		userGroup.Members = userGroup.Members[1:]
	}
}

func (testWrap *testWrap) userGroupGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = checkUserGroup(user, testData, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		userGroup1 := testData.userGroups[0]

		url := testWrap.serverURL + "/v1/user/group/" + userGroup1.GroupID.String()

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
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var userGroup model.UserGroup

		err = json.Unmarshal(data, &userGroup)
		if err != nil {
			t.Fatal(err)
		}

		if userGroup.GroupID.String() != userGroup1.GroupID.String() {
			t.Fatalf("Expected UUID: %s, received: %s", userGroup1.GroupID, userGroup.GroupID)
		}
	}
}

func (testWrap *testWrap) userGroupGetList(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.serverURL, testWrap.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = checkUserGroup(user, testData, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		url := testWrap.serverURL + "/v1/user/group"

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
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var userGroupList []model.UserGroup

		err = json.Unmarshal(data, &userGroupList)
		if err != nil {
			t.Fatal(err)
		}

		if len(userGroupList) == 0 {
			t.Fatalf("Expected: userGroups, received: empty, body: %s", data)
		}
	}
}

func userGroupCreate(user *User, systemID, baseURL, name, description string, client *http.Client) (*model.UserGroup, string, error) {
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

	request.Header.Set("AuthUserId", user.id)
	request.Header.Set("Authorization", "Bearer "+user.accessToken)
	request.Header.Set("EhrSystemId", systemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if response.StatusCode != http.StatusCreated {
		if response.StatusCode == http.StatusConflict {
			return nil, "", errors.ErrAlreadyExist
		}

		return nil, "", errors.New(response.Status + " data: " + string(data))
	}

	var userGroup2 model.UserGroup

	err = json.Unmarshal(data, &userGroup2)
	if err != nil {
		return nil, "", err
	}

	if userGroup2.Name != userGroup.Name {
		return nil, "", errors.ErrFieldIsIncorrect("Name")
	}

	requestID := response.Header.Get("RequestId")

	return &userGroup2, requestID, nil
}

func checkUserGroup(user *User, testData *TestData, baseURL string, client *http.Client) error {
	if len(testData.userGroups) == 0 {
		name := fakeData.GetRandomStringWithLength(10)
		description := fakeData.GetRandomStringWithLength(10)

		userGroup, reqID, err := userGroupCreate(user, testData.ehrSystemID, baseURL, name, description, client)
		if err != nil {
			return fmt.Errorf("userGroupCreate error: %w", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, baseURL, client)
		if err != nil {
			return fmt.Errorf("requestWait error, err: %w", err)
		}

		testData.userGroups = append(testData.userGroups, userGroup)
	}

	return nil
}

func userGroupAddUser(user, addingUser *User, userGroup *model.UserGroup, testData *TestData, baseURL string, client *http.Client) (string, error) {
	if user.accessToken == "" {
		err := user.login(testData.ehrSystemID, baseURL, client)
		if err != nil {
			return "", fmt.Errorf("user.login error: %w", err)
		}
	}

	url := baseURL + "/v1/user/group/" + userGroup.GroupID.String() + "/user_add/" + addingUser.id + "/admin"

	request, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return "", fmt.Errorf("NewRequest error: %w", err)
	}

	request.Header.Set("AuthUserId", user.id)
	request.Header.Set("Authorization", "Bearer "+user.accessToken)
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("http request error: %w", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status + " data: " + string(data))
	}

	requestID := response.Header.Get("RequestId")

	return requestID, nil
}
