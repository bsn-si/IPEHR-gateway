package fake_data

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
)

func QueryExecRequest(ehrId string) []byte {
	req := `{
	  "q": "SELECT c FROM EHR e[ehr_id/value=$ehr_id] 
	  			CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.encounter.v1]
					CONTAINS OBSERVATION obs[openEHR-EHR-OBSERVATION.blood_pressure.v1]
			WHERE obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp",
	  "offset": 0,
	  "fetch": 10,
	  "query_parameters": {
		"ehr_id": "` + ehrId + `",
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
		"_created": "` + time.Now().Format(common.OPENEHR_TIME_FORMAT) + `"
	  },
	  "q": "` + query + `",
	  "columns": [
		{
		  "name": "#0",
		  "path": "/ehr_id/value"
		},
		{
		  "name": "startTime",
		  "path": "/context/start_time/value"
		},
		{
		  "name": "cid",
		  "path": "/uid/value"
		},
		{
		  "name": "#3",
		  "path": "/name"
		}
	  ],
	  "rows": [
		[
		  "` + uuid.New().String() + `",
		  "2017-02-16T13:50:11.308+01:00",
		  "90910cf0-66a0-4382-b1f8-c0f27e81b42d::openEHRSys.example.com::1",
		  {
			"_type": "DV_TEXT",
			"value": "Labs"
		  }
		]
	  ]
	}`)
}
