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

func userGroupCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedIn error:", err)
		}

		name := fakeData.GetRandomStringWithLength(10)
		description := fakeData.GetRandomStringWithLength(10)

		userGroup, reqID, err := _userGroupCreate(testData, user, name, description)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal("requestWait error: %w reqID: %w", err, reqID)
		}

		testData.userGroups = append(testData.userGroups, userGroup)
	}
}

func userGroupAddUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedIn error:", err)
		}

		addingUser := testData.users[1]

		err = checkUserGroup(testData, user)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		userGroup := testData.userGroups[0]

		reqID, err := _userGroupAddUser(user, addingUser, userGroup, testData)
		if err != nil {
			t.Fatal("userGroupAddUser error: ", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal("requestWait error: %w reqID: %w", err, reqID)
		}

		userGroup.Members = append(userGroup.Members, addingUser.id)
	}
}

func userGroupRemoveUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedIn error:", err)
		}

		err = checkUserGroup(testData, user)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		userGroup := testData.userGroups[0]

		if len(userGroup.Members) == 0 {
			addingUser := testData.users[1]

			reqID, err := _userGroupAddUser(user, addingUser, userGroup, testData)
			if err != nil {
				t.Fatal("userGroupAddUser error: ", err)
			}

			err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal("requestWait error: %w reqID: %w", err, reqID)
			}
		}

		removingUserID := userGroup.Members[0]

		url := testData.serverURL + "/v1/user/group/" + userGroup.GroupID.String() + "/user_remove/" + removingUserID

		request, err := http.NewRequest(http.MethodPost, url, nil)
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
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		// Checking that the user has been deleted
		url = testData.serverURL + "/v1/user/group/" + userGroup.GroupID.String()

		request, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err = testData.httpClient.Do(request)
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

func userGroupGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedIn error:", err)
		}

		err = checkUserGroup(testData, user)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		userGroup1 := testData.userGroups[0]

		url := testData.serverURL + "/v1/user/group/" + userGroup1.GroupID.String()

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

func userGroupGetList(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedIn error:", err)
		}

		err = checkUserGroup(testData, user)
		if err != nil {
			t.Fatal("checkUserGroup error: ", err)
		}

		url := testData.serverURL + "/v1/user/group"

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

func _userGroupCreate(testData *TestData, user *User, name, description string) (*model.UserGroup, string, error) {
	userGroup := &model.UserGroup{
		Name:        name,
		Description: description,
	}

	data, _ := json.Marshal(userGroup)
	url := testData.serverURL + "/v1/user/group"

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("AuthUserId", user.id)
	request.Header.Set("Authorization", "Bearer "+user.accessToken)
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := testData.httpClient.Do(request)
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

func checkUserGroup(testData *TestData, user *User) error {
	if len(testData.userGroups) == 0 {
		name := fakeData.GetRandomStringWithLength(10)
		description := fakeData.GetRandomStringWithLength(10)

		userGroup, reqID, err := _userGroupCreate(testData, user, name, description)
		if err != nil {
			return fmt.Errorf("userGroupCreate error: %w", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			return fmt.Errorf("requestWait error, err: %w", err)
		}

		testData.userGroups = append(testData.userGroups, userGroup)
	}

	return nil
}

func _userGroupAddUser(user, addingUser *User, userGroup *model.UserGroup, testData *TestData) (string, error) {
	url := testData.serverURL + "/v1/user/group/" + userGroup.GroupID.String() + "/user_add/" + addingUser.id + "/admin"

	request, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return "", fmt.Errorf("NewRequest error: %w", err)
	}

	request.Header.Set("AuthUserId", user.id)
	request.Header.Set("Authorization", "Bearer "+user.accessToken)
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := testData.httpClient.Do(request)
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
