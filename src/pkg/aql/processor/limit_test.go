package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser_Limit(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    *Limit
		wantErr bool
	}{
		{
			"1. query without LIMIT",
			`SELECT val FROM e`,
			nil,
			false,
		},
		{
			"2. query with LIMIT",
			`SELECT val FROM e LIMIT 1`,
			&Limit{
				Limit: 1,
			},
			false,
		},
		{
			"3. query with LIMIT AND OFFSET",
			`SELECT val FROM e LIMIT 10 OFFSET 123`,
			&Limit{
				Limit:  10,
				Offset: 123,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got.Limit); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}
