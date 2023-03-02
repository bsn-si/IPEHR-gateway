package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_ToQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			"1. Select one column",
			"SELECT\n\tvalue\nFROM\n\tc C",
			false,
		},
		{
			"2. Select With path",
			"SELECT\n\tc/data[at0003]/value/magnitude AS mag\nFROM\n\tC c",
			false,
		},
		{
			"3. several fields",
			"SELECT\n" +
				"\tc/data[at0003]/value/magnitude AS mag,\n" +
				"\tc/data[at0003]/value/unit AS unit\n" +
				"FROM\n\tC c",
			false,
		},
		{
			"4. FROM ",
			"SELECT\n\tvalue\nFROM\n\tEHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']",
			false,
		},
		{
			"5. FROM with CONTAINS",
			"SELECT\n\tvalue\nFROM\n\tEHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']\n" +
				"\tCONTAINS\n\t\tObservation o[openEHR-EHR-OBSERVATION.body_temperature-zn.v1]",
			false,
		},
		{
			"6. FROM with CONTAINS param",
			"SELECT\n\tvalue\nFROM\n\tEHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']\n" +
				"\tCONTAINS\n\t\tObservation o[$param]",
			false,
		},
		{
			"7. FROM with NOT CONTAINS",
			"SELECT\n\tvalue\nFROM\n\tEHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']\n" +
				"\tNOT CONTAINS\n\t\tObservation o[$param]",
			false,
		},
		{
			"8. FROM with CONTAINS AND",
			"SELECT\n\tvalue\nFROM\n\tEHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']\n" +
				"\tCONTAINS\n\t\tObservation o1[$param1] AND Observation o2[$param2]",
			false,
		},
		{
			"9. FROM with CONTAINS (OR)",
			"SELECT\n\tvalue\nFROM\n\tEHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']\n" +
				"\tCONTAINS\n\t\t(Observation o1[$param1] OR Observation o2[$param2])",
			false,
		},
		{
			"10. Select one column with WHERE",
			"SELECT value\n" +
				"FROM c C\n" +
				"WHERE c/value > 10",
			false,
		},
		{
			"11. Select one column with WHERE AND",
			"SELECT value\n" +
				"FROM c C\n" +
				"WHERE (c/value > 10 AND c/value1/value2 <= 11)",
			false,
		},
		{
			"12. ORDER",
			"SELECT value1, value2, value3\n" +
				"FROM c C\n" +
				"ORDER BY value1 ASC, value2 DESC, value3",
			false,
		},
		{
			"13. LIMIT",
			"SELECT value1, value2, value3\n" +
				"FROM c C\n" +
				"LIMIT 10",
			false,
		},
		{
			"14. LIMIT AND OFFSET",
			"SELECT value1, value2, value3\n" +
				"FROM c C\n" +
				"LIMIT 10 OFFSET 10",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.query)
			query, err := NewAqlProcessor(tt.query).Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Processor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotStr := query.String()
			t.Log(gotStr)

			got, err := NewAqlProcessor(gotStr).Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Processor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, query, got)
		})
	}
}
