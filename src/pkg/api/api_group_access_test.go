package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupAccessAPI(t *testing.T) {
	r := gin.New()

	cfgPath := "../../../config.json.example"
	cfg := config.New(cfgPath)
	err := cfg.Reload()
	if err != nil {
		t.Fatal(err)
	}

	api := New(cfg)
	r.Use(api.Auth)
	r.GET("/v1/access/group/:group_id", api.GroupAccess.Get)
	r.POST("/v1/access/group", api.GroupAccess.Create)

	ts := httptest.NewServer(r)
	defer ts.Close()

	var (
		httpClient    http.Client
		groupAccessId string
		testUserId, _ = uuid.NewUUID()
	)

	t.Run("Group Access creating", func(t *testing.T) {
		description := fake_data.GetRandomStringWithLength(50)

		req := []byte(`{
			"description": "` + description + `"
		}`)

		request, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/access/group/", bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserId.String())

		response, err := httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}
		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusCreated, response.StatusCode, data)
		}

		var groupAccess model.GroupAccess
		if err = json.Unmarshal(data, &groupAccess); err != nil {
			t.Fatal(err)
		}

		groupAccessId = groupAccess.GroupId
	})

	t.Run("Wrong Group Access getting", func(t *testing.T) {
		groupAccessIdWrong, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/access/group/"+groupAccessIdWrong.String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testUserId.String())

		response, err := httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}
		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusNotFound, response.StatusCode, data)
		}
	})

	t.Run("Group Access getting", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/access/group/"+groupAccessId, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", testUserId.String())

		response, err := httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}
		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var groupAccessGot model.GroupAccess
		if err = json.Unmarshal(data, &groupAccessGot); err != nil {
			t.Fatal(err)
		}

		if groupAccessId != groupAccessGot.GroupId {
			t.Fatal("Got wrong group")
		}
	})
}
