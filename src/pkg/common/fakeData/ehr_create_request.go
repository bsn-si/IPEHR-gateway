package fakeData

import (
	"github.com/google/uuid"
)

func EhrCreateRequest() []byte {
	subjectID := uuid.New().String()

	return EhrCreateCustomRequest(subjectID, "test")
}

func EhrCreateCustomRequest(subjectID, subjectNamespace string) []byte {
	return []byte(`{
	  "_type": "EHR_STATUS",
	  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
	  "name": {
		"value": "EHR Status"
	  },
	  "subject": {
		"external_ref": {
		  "id": {
			"_type": "GENERIC_ID",
			"value": "` + subjectID + `",
			"scheme": "id_scheme"
		  },
		  "namespace": "` + subjectNamespace + `",
		  "type": "PERSON"
		}
	  },
	  "is_modifiable": true,
	  "is_queryable": true
	}`)
}
