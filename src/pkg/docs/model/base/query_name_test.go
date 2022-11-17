package base_test

import (
	"hms/gateway/pkg/docs/model/base"
	"testing"
)

func TestQueryName_String(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			"1. error on unmarshal",
			"",
			true,
		},
		{
			"2. error on unmarshal",
			"reversed.domain::name",
			false,
		},
		{
			"3. compare with result",
			"without_domain_name",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := base.NewQueryName(tt.data)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewQueryName() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if (q.String() != tt.data) != tt.wantErr {
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
