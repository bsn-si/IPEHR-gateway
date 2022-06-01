package ehr

import (
	"encoding/json"
	"github.com/google/uuid"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"testing"
)

func TestSave(t *testing.T) {
	t.Skip("Not finished")
	jsonDoc := []byte(`{
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

	ehrService := NewEhrService(service.NewDefaultDocumentService())

	var ehrReq model.EhrCreateRequest

	err := json.Unmarshal(jsonDoc, &ehrReq)
	if err != nil {
		t.Fatal(err)
	}

	ehrDoc := ehrService.Create(&ehrReq)

	testUserId := uuid.New().String()

	err = ehrService.Save(testUserId, ehrDoc)
	if err != nil {
		t.Fatal(err)
	}
}
