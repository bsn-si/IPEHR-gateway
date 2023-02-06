package api_test

/*
func (testWrap *testWrap) docSetAccessSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user1 := testData.users[0]
		user2 := testData.users[1]

		url := testWrap.serverURL + "/v1/access/document"

		req := model.DocAccessSetRequest{
			UserID:      user2.id,
			CID:         "",
			AccessLevel: access.LevelToString(access.Read),
		}

		reqJSON, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqJSON))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user1.id)
		request.Header.Set("Authorization", "Bearer "+user1.accessToken)
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
			t.Fatal(response.Status, "Body:", string(data))
		}

		var ehrDoc model.EHR

		err = json.Unmarshal(data, &ehrDoc)
		if err != nil {
			t.Fatal(err)
		}

		if ehrDoc.EhrID.Value != user1.ehrID {
			t.Fatalf("Expected %s, received %s", user1.ehrID, ehrDoc.EhrID.Value)
		}
	}
}
*/
