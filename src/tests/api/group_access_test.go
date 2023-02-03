package api_test

/*
func createGroupAccess(userID, accessToken, baseURL string, client *http.Client) (*model.GroupAccess, error) {
	description := fakeData.GetRandomStringWithLength(50)

	req := []byte(`{
			"description": "` + description + `"
		}`)

	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/access/group", bytes.NewReader(req))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(response.Status)
	}

	var groupAccess model.GroupAccess
	if err = json.Unmarshal(data, &groupAccess); err != nil {
		return nil, err
	}

	return &groupAccess, nil
}

func (testWrap *testWrap) accessGroupCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		ga, err := createGroupAccess(user.id, user.accessToken, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatalf("Expected group access, received error: %v", err)
		}

		testData.groupsAccess = append(testData.groupsAccess, ga)
	}
}

func (testWrap *testWrap) wrongAccessGroupGetting(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		groupAccessIDWrong, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}

		url := testWrap.serverURL + "/v1/access/group/" + groupAccessIDWrong.String()

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusNotFound, response.StatusCode, data)
		}
	}
}

func (testWrap *testWrap) accessGroupGetting(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		if len(testData.groupsAccess) == 0 {
			t.Fatal("GroupAccess required")
		}

		ga := testData.groupsAccess[0]

		url := testWrap.serverURL + "/v1/access/group/" + ga.GroupUUID.String()

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var groupAccessGot model.GroupAccess
		if err = json.Unmarshal(data, &groupAccessGot); err != nil {
			t.Fatal(err)
		}

		if ga.GroupUUID.String() != groupAccessGot.GroupUUID.String() {
			t.Fatal("Got wrong group")
		}
	}
}
*/
