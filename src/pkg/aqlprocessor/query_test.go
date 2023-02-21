package aqlprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestProcessor_QueryEncdodeDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			"1. Encode/Decode with one condition",
			`SELECT
				e/ehr_id/value AS ID,
				o/archetype_node_id,
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude,
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units
			 FROM EHR e [ehr_id/value=$ehrUid] 
			 	CONTAINS COMPOSITION c
					CONTAINS OBSERVATION o [openEHR-EHR-OBSERVATION.pulse.v2]
			 WHERE 
			 	o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			false,
		},
		{
			"2. Encode/Decode query with two conditions",
			`SELECT
				e/ehr_id/value AS ID,
				o/archetype_node_id,
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude,
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units
			 FROM EHR e [ehr_id/value=$ehrUid] 
			 	CONTAINS COMPOSITION c
					CONTAINS OBSERVATION o [openEHR-EHR-OBSERVATION.pulse.v2]
			 WHERE 
			 	o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100 AND
			 	o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude < 1000.5`,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin, err := NewAqlProcessor(tt.query).Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
				return
			}

			data, err := msgpack.Marshal(origin)
			assert.Nil(t, err)

			got := &Query{}
			err = msgpack.Unmarshal(data, got)
			assert.Nil(t, err)

			assert.Equal(t, origin, got)
		})
	}
}
