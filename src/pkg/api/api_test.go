package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/docs/model"
)

func Test_API(t *testing.T) {
	r := gin.New()

	api := NewAPI()
	r.Use(api.Auth)
	r.GET("/v1/ehr/:ehrid", api.Ehr.GetById)
	r.POST("/v1/ehr", api.Ehr.Create)
	r.GET("/v1/ehr/:ehrid/ehr_status/:versionid", api.EhrStatus.GetById)
	r.PUT("/v1/ehr/:ehrid/ehr_status", api.EhrStatus.Update)

	ts := httptest.NewServer(r)
	defer ts.Close()

	var (
		httpClient  http.Client
		ehrId       string
		ehrStatusId string
		testUserId  = "11111111-1111-1111-1111-111111111111"
	)

	t.Run("EHR creating", func(t *testing.T) {
		req := []byte(`{
		  "_type": "EHR_STATUS",
		  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
		  "name": {
			"value": "EHR Status"
		  },
		  "subject": {
			"external_ref": {
			  "id": {
				"_type": "GENERIC_ID",
				"value": "ins01",
				"scheme": "id_scheme"
			  },
			  "namespace": "examples",
			  "type": "PERSON"
			}
		  },
		  "is_modifiable": true,
		  "is_queryable": true
		}`)
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/ehr", bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
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
			t.Errorf("Expected %d, received %d", http.StatusOK, response.StatusCode)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		ehrId = ehr.EhrId.Value
		if ehrId == "" {
			t.Error("EhrId missing")
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

	t.Run("EHR_STATUS update", func(t *testing.T) {
		req := []byte(fmt.Sprintf(`{
		  "_type": "EHR_STATUS",
		  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
		  "name": {
			"value": "EHR Status"
		  },
		  "uid": {
			"_type": "OBJECT_VERSION_ID",
			"value": "%s::openEHRSys.example.com::2"
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
		}`, ehrStatusId))

		request, err := http.NewRequest(http.MethodPut, ts.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", ehrId), bytes.NewReader(req))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", testUserId)
		request.Header.Set("If-Match", ehrStatusId)

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

		if ehrStatus.Uid.Value != ehrStatusId+"::openEHRSys.example.com::2" {
			t.Log("Response body:", string(data))
			t.Error("EHR_STATUS uid mismatch")
			return
		}
	})
}
