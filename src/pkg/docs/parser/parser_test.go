package parser

import (
	"testing"
)

func TestParseDocument(t *testing.T) {
	ehrId := "7d44b88c-4199-4bad-97dc-d78268e01398"

	inJson := []byte(`{
	  "system_id": {
		"value": "d60e2348-b083-48ce-93b9-916cef1d3a5a"
	  },
	  "ehr_id": {
		"value": "` + ehrId + `"
	  },
	  "ehr_status": {
		"id": {
		  "_type": "OBJECT_VERSION_ID",
		  "value": "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
		},
		"namespace": "local",
		"type": "EHR_STATUS"
	  },
	  "ehr_access": {
		"id": {
		  "_type": "OBJECT_VERSION_ID",
		  "value": "59a8d0ac-140e-4feb-b2d6-af99f8e68af8::openEHRSys.example.com::1"
		},
		"namespace": "local",
		"type": "EHR_ACCESS"
	  },
	  "time_created": {
		"value": "2015-01-20T19:30:22.765+01:00"
	  }
	}`)

	res, err := ParseDocument(inJson)
	if err != nil {
		t.Fatal(err)
	}

	if ehrId != res.EhrId.Value {
		t.Fatal("Document is not parsed correctly")
	}
}

func TestParseComposition(t *testing.T) {
	compositionUid := "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
	inJson := []byte(`{
	  "_type": "COMPOSITION",
	  "archetype_node_id": "openEHR-EHR-COMPOSITION.encounter.v1",
	  "name": {
		"value": "Vital Signs"
	  },
	  "uid": {
		"_type": "OBJECT_VERSION_ID",
		"value": "` + compositionUid + `"
	  },
	  "archetype_details": {
		"archetype_id": {
		  "value": "openEHR-EHR-COMPOSITION.encounter.v1"
		},
		"template_id": {
		  "value": "Example.v1::c7ec861c-c413-39ff-9965-a198ebf44747"
		},
		"rm_version": "1.0.2"
	  },
	  "language": {
		"terminology_id": {
		  "value": "ISO_639-1"
		},
		"code_string": "en"
	  },
	  "territory": {
		"terminology_id": {
		  "value": "ISO_3166-1"
		},
		"code_string": "NL"
	  },
	  "category": {
		"value": "event",
		"defining_code": {
		  "terminology_id": {
			"value": "openehr"
		  },
		  "code_string": "433"
		}
	  },
	  "composer": {
		"_type": "PARTY_IDENTIFIED",
		"external_ref": {
		  "id": {
			"_type": "GENERIC_ID",
			"value": "16b74749-e6aa-4945-b760-b42bdc07098a"
		  },
		  "namespace": "openEHRSys.example.com",
		  "type": "PERSON"
		},
		"name": "A name"
	  },
	  "context": {
		"start_time": {
		  "value": "2014-11-18T09:50:35.000+01:00"
		},
		"setting": {
		  "value": "other care",
		  "defining_code": {
			"terminology_id": {
			  "value": "openehr"
			},
			"code_string": "238"
		  }
		}
	  },
	  "content": []
	}`)

	res, err := ParseComposition(inJson)
	if err != nil {
		t.Fatal(err)
	}

	if compositionUid != res.Uid.Value {
		t.Fatal("Composition is not parsed correctly")
	}
}
