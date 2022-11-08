package api_test

/*
func (testWrap *testWrap) docGrantAccessSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user1 := testData.users[0]
		user2 := testData.users[1]

		url := testWrap.server.URL + "/v1/access/document/manage"

		req := model.DocAccessManageRequest{
			ToUserID:    user2.id,
			CID:         "",
			AccessLevel: access.Read.String(),
		}

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("Prefer", "return=representation")

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
			t.Fatal(err)
		}

		var ehrDoc model.EHR

		err = json.Unmarshal(data, &ehrDoc)
		if err != nil {
			t.Fatal(err)
		}

		if ehrDoc.EhrID.Value != user.ehrID {
			t.Fatalf("Expected %s, received %s", user.ehrID, ehrDoc.EhrID.Value)
		}
	}
}
*/
