package aqlprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorTestProcessor(t *testing.T) {
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
