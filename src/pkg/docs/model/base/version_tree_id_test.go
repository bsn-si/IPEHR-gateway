package base

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVersionTreeID_String(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
		want    string
	}{
		{
			"1. error because empty val",
			"",
			true,
			"",
		},
		{
			"2. ver is equal",
			"1",
			false,
			"1",
		},
		{
			"3. ver is not equal",
			"1",
			true,
			"2",
		},
		{
			"4. ver contain branch",
			"1.0",
			false,
			"1.0",
		},
		{
			"5. ver contain wrong branch",
			"1.0",
			true,
			"1.1",
		},
		{
			"6. ver contain branch version",
			"1.1.1",
			false,
			"1.1.1",
		},
		{
			"7. ver contain wrong branch version",
			"1.1.1",
			true,
			"1.1.2",
		},
		{
			"8. ver is incorrect",
			"incorrect.val",
			true,
			"",
		},
		{
			"9. sanitize string",
			"9999999999999999999.9999999999999999999.9999999999999999999999999999999999",
			true,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersionTreeID(tt.data)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("NewVersionTreeID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, v.String()); diff != "" {
				if tt.wantErr {
					return
				}
				t.Errorf("VersionTreeID.String() mismatch{-want;+got}\n\t%s", diff)
			}
		})
	}
}

func TestVersionTreeID_Equal(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
		want    string
	}{
		{
			"1. trunk ver is equal",
			"1",
			false,
			"1",
		},
		{
			"2. ver is not equal",
			"1",
			true,
			"2",
		},
		{
			"3. branch equal",
			"1.1",
			false,
			"1.1",
		},
		{
			"4. branch not equal",
			"1.2",
			true,
			"1.1",
		},
		{
			"5. ver contain branch version",
			"1.1.1",
			false,
			"1.1.1",
		},
		{
			"6. ver contain branch version",
			"1.1.1",
			true,
			"1.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersionTreeID(tt.data)
			if err != nil {
				t.Errorf("NewVersionTreeID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (v.Equal(tt.want) != true) && !tt.wantErr {
				t.Errorf("NewVersionTreeID.Equal() mismatch want: %s", tt.want)
			}
		})
	}
}

func TestVersionTreeID_Increase(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
		want    string
	}{
		{
			"1. trunk ver inc",
			"1",
			false,
			"2",
		},
		{
			"2. branch inc",
			"1.1",
			false,
			"1.2",
		},
		{
			"3. branch ver inc",
			"1.1.1",
			false,
			"1.1.2",
		},
		{
			"4. branch ver inc with err",
			"1.1.1",
			true,
			"1.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersionTreeID(tt.data)
			if err != nil {
				t.Errorf("NewVersionTreeID() error = %v, wantErr %v", err, tt.wantErr)
			}

			inc := v.Increase()

			if (inc != tt.want) && !tt.wantErr {
				t.Errorf("NewVersionTreeID.Increase() mismatch want: %s, got: %s", tt.want, inc)
			}
		})
	}
}
