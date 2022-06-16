package fake_data

import (
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
)

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
