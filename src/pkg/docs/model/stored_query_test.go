package model_test

import (
	"encoding/json"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
)

func TestStoredQuery_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			"1. no errors",
			[]byte(`{
				"name":"org.openehr::compositions",
				"type":"aql",
				"version":"1.0.1",
				"saved":"2017-07-16T19:20:30.450+01:00",
				"q":"SELECT 1"
			}`),
			false,
		},
		{
			"2. validation error if name is empty",
			[]byte(`{
				"name":"",
				"type":"aql",
				"version":"1.0.1",
				"saved":"2017-07-16T19:20:30.450+01:00",
				"q":"SELECT 1"
			}`),
			true,
		},
		{
			"3. validation error if type is empty",
			[]byte(`{
				"name":"org.openehr::compositions",
				"type":"",
				"version":"1.0.1",
				"saved":"2017-07-16T19:20:30.450+01:00",
				"q":"SELECT 1"
			}`),
			true,
		},
		{
			"4. validation error if query is empty",
			[]byte(`{
				"name":"org.openehr::compositions",
				"type":"aql",
				"version":"1.0.1",
				"saved":"2017-07-16T19:20:30.450+01:00",
				"q":""
			}`),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.StoredQuery{}
			if err := json.Unmarshal(tt.data, &got); err != nil {
				t.Errorf("StoredQuery.UnmarshalJSON() error = %v", err)
				return
			}

			if err := got.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("StoredQuery.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
