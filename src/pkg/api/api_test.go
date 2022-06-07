package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/common/fake_data"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/docs/model"
)

func Test_API(t *testing.T) {
	r := gin.New()

	api := NewAPI()
	r.Use(api.Auth)
	r.GET("/v1/ehr/:ehrid", api.Ehr.GetById)
	r.POST("/v1/ehr", api.Ehr.Create)
	r.PUT("/v1/ehr/:ehrid", api.Ehr.CreateWithId)
	r.GET("/v1/ehr/:ehrid/ehr_status/:versionid", api.EhrStatus.GetById)
	r.GET("/v1/ehr/:ehrid/ehr_status", api.EhrStatus.GetStatus)
	r.PUT("/v1/ehr/:ehrid/ehr_status", api.EhrStatus.Update)

	ts := httptest.NewServer(r)
	defer ts.Close()

	var (
		httpClient  http.Client
		ehrId       string
		ehrStatusId string
		testUserId  = "11111111-1111-1111-1111-111111111111"
		testUserId2 = "22222222-2222-2222-2222-222222222222"
	)

	t.Run("EHR creating", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/ehr", ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserId)
		request.Header.Set("Prefer", "return=representation")

		response, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			t.Errorf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		ehrId = response.Header.Get("ETag")
		if ehrId == "" {
			t.Error("EhrId missing")
			return
		}
	})

	t.Run("EHR creating with id", func(t *testing.T) {
		ehrId2 := uuid.New().String()

		request, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/ehr/"+ehrId2, ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserId2)
		request.Header.Set("Prefer", "return=representation")

		response, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			t.Errorf("Expected %d, received %d", http.StatusCreated, response.StatusCode)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		newEhrId := ehr.EhrId.Value
		if newEhrId != ehrId2 {
			t.Error("EhrId is not matched")
			return
		}
	})

	t.Run("EHR creating with id for the same user", func(t *testing.T) {
		ehrId3 := uuid.New().String()

		request, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/ehr/"+ehrId3, ehrCreateBodyRequest())
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserId2)

		response, err := httpClient.Do(request)

		if response.StatusCode != http.StatusConflict {
			t.Errorf("Expected %d, received %d", http.StatusConflict, response.StatusCode)
			return
		}
	})

	t.Run("EHR getting", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/ehr/"+ehrId, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testUserId)

		response, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		if ehrId != ehr.EhrId.Value {
			t.Error("EHR document mismatch")
			return
		}

		ehrStatusId = ehr.EhrStatus.Id.Value
	})

	t.Run("EHR_STATUS getting", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, ts.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", ehrId, ehrStatusId), nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", testUserId)

		response, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Error(err)
			return
		}

		if ehrStatus.Uid == nil || ehrStatus.Uid.Value != ehrStatusId {
			t.Error("EHR_STATUS document mismatch")
			return
		}
	})

	t.Run("EHR_STATUS getting by version time", func(t *testing.T) {
		ehrId := uuid.New().String()
		versionAtTime := time.Now()

		request, err := http.NewRequest(http.MethodGet, ts.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", ehrId), nil)
		if err != nil {
			t.Error(err)
			return
		}

		q := request.URL.Query()
		q.Add("version_at_time", versionAtTime.Format(ALLOWED_TIME_FORMAT))

		request.Header.Set("AuthUserId", testUserId)
		request.URL.RawQuery = q.Encode()

		response, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Errorf("Expected %d, received %d", http.StatusNotFound, response.StatusCode)
			return
		}
	})

	t.Run("EHR_STATUS update", func(t *testing.T) {
		// replace substring in ehrStatusId
		newEhrStatusId := strings.Replace(ehrStatusId, "::openEHRSys.example.com::1", "::openEHRSys.example.com::2", 1)

		req := []byte(fmt.Sprintf(`{
		  "_type": "EHR_STATUS",
		  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
		  "name": {
			"value": "EHR Status"
		  },
		  "uid": {
			"_type": "OBJECT_VERSION_ID",
			"value": "%s"
		  },
		  "subject": {
			"external_ref": {
			  "id": {
				"_type": "HIER_OBJECT_ID",
				"value": "324a4b23-623d-4213-cc1c-23f233b24234"
			  },
			  "namespace": "DEMOGRAPHIC",
			  "type": "PERSON"
			}
		  },
		  "other_details": {
			"_type": "ITEM_TREE",
			"archetype_node_id": "at0001",
			"name": {
			  "value": "Details"
			},
			"items": []
		  },
		  "is_modifiable": true,
		  "is_queryable": true
		}`, newEhrStatusId))

		request, err := http.NewRequest(http.MethodPut, ts.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", ehrId), bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserId)
		request.Header.Set("If-Match", ehrStatusId)
		request.Header.Set("Prefer", "return=representation")

		response, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Error(err)
			return
		}

		updatedEhrStatusId := response.Header.Get("ETag")

		if updatedEhrStatusId != newEhrStatusId {
			t.Log("Response body:", string(data))
			t.Error("EHR_STATUS uid mismatch")
			return
		}
	})

}

func ehrCreateBodyRequest() *bytes.Reader {
	req := fake_data.EhrCreateRequest()
	return bytes.NewReader(req)
}
