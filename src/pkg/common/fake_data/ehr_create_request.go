package fake_data

func EhrCreateRequest() []byte {
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
}
