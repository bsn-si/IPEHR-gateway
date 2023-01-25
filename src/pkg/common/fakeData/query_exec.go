package fakeData

import (
	"strings"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
)

func QueryExecRequest(ehrID string) []byte {
	req := `{
	  "q": "SELECT e/ehr_id/value, 
	  			   c/context/start_time/value as startTime, 
				   c/uid/value as cid, 
				   c/name 
	  		FROM EHR e[ehr_id/value=$ehr_id] 
			CONTAINS COMPOSITION c [openEHR-EHR-COMPOSITION.encounter.v1] 
				CONTAINS OBSERVATION obs [openEHR-EHR-OBSERVATION.blood_pressure.v1] 
			WHERE obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp",	
	  "offset": 0,
	  "fetch": 10,
	  "query_parameters": {
		"ehr_id": "` + ehrID + `",
		"systolic_bp": 140
	  }
	}`
	req = strings.ReplaceAll(req, "\n", "")
	req = strings.ReplaceAll(req, "\t", "")

	return []byte(req)
}

func QueryExecResponse(query string) []byte {
	return []byte(`{
		"meta": {
		"_href": "",
		"_type": "RESULTSET",
		"_schema_version": "1.0.0",
		"_created": "` + time.Now().Format(common.OpenEhrTimeFormat) + `"
		},
		"q": "` + query + `",
		"columns": [
		{
		  "name": "#0",
		  "path": "/ehr_id/value"
		},
		{
		  "name": "systolic",
		  "path": "/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude"
		}
		],
		"rows": [
		  "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::piri.ehrscape.com::1",
		  266.0
		]
	  }`)
}
