package aqlprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorTestProcessor(t *testing.T) {
	const query = `SELECT 
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude AS temperature,
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units AS unit,
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/TYPE,
	'hello_world',
	"hello_world_2",
	123,
	'2010-05-12',
	"2010-05-12",
	NOW() 
FROM
   EHR[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']
       CONTAINS Observation o[openEHR-EHR-OBSERVATION.body_temperature-zn.v1]
WHERE
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude > $temperature
	AND
    o/data[at0002]/events[at0003]/data[at0001]/items[at0.63 and name/value='Symptoms']/value/defining_code/code_string=$chills
ORDER BY temperature DESC, unit ASC
LIMIT 3 OFFSET 1
	`

	tests := []struct {
		name    string
		query   string
		want    Query
		wantErr bool
	}{
		{
			"1. invalid query",
			"SELECT invalid",
			Query{},
			true,
		},

		// {
		// 	"100. Real query",
		// 	query,
		// 	Query{
		// 		Where: &Where{},
		// 		Order: &Order{
		// 			[]OrderBy{
		// 				{"temperature", DescendingOrdering},
		// 				{"unit", AscendingOrdering},
		// 			},
		// 		},
		// 		Limit: &Limit{
		// 			Limit:  3,
		// 			Offset: 1,
		// 		},
		// 	},
		// 	false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
