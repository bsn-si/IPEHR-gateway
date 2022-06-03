package fake_data

import (
	"github.com/google/uuid"
)

func EhrCreateRequest() []byte {

	subjectId := uuid.New().String()

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
			"value": "` + subjectId + `",
			"scheme": "id_scheme"
		  },
		  "namespace": "test",
		  "type": "PERSON"
		}
	  },
	  "is_modifiable": true,
	  "is_queryable": true
	}`)
}
