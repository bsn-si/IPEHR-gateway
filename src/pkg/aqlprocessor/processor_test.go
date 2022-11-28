package aqlprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorTestProcessor(t *testing.T) {
	const query = `SELECT 
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude AS temperature,
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units AS unit
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

	want := Query{
		Where: &Where{},
		Order: &Order{
			[]OrderBy{
				{"temperature", DescendingOrdering},
				{"unit", AscendingOrdering},
			},
		},
		Limit: &Limit{
			Limit:  3,
			Offset: 1,
		},
	}

	p := NewAqlProcessor(query)
	got, err := p.Process()
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
