package hm

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	maxPositiveInt64 = 1<<63 - 1
	maxNegativeInt64 = -(maxPositiveInt64 + 1)
)

func TestDataEncryption_Int64(t *testing.T) {
	tests := []struct {
		name    string
		keyHex  string
		num     int64
		want    int64
		wantErr bool
	}{
		{
			"1. success encryption positive Int64",
			"ffffffffffffffff6b72fd6ab0358aa9f52cee4a11c2764c4e745f47ddb1d5e1",
			2147483647,
			9223372034707292160,
			false,
		},
		{
			"2. success encryption negative Int64",
			"ffffffffffffffff6b72fd6ab0358aa9f52cee4a11c2764c4e745f47ddb1d5e1",
			-2147483648,
			-9223372030412324865,
			false,
		},
		{
			"3. success encryption Int64 = 0",
			"ffffffffffffffff6b72fd6ab0358aa9f52cee4a11c2764c4e745f47ddb1d5e1",
			0,
			4294967295,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyBytes, err := hex.DecodeString(tt.keyHex)
			if err != nil {
				t.Errorf("Encryption key error: %v", err)
			}

			key := Key{}
			copy(key[:], keyBytes)

			x, err := EncryptInt64(tt.num, &key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptInt error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, x.Int64())
		})
	}
}

func TestDataEncryption_Float64(t *testing.T) {
	tests := []struct {
		name    string
		keyHex  string
		num     float64
		wantErr bool
	}{
		{
			"1. success encryption positive Float64",
			"5012ee13ae2de3646b72fd6ab0358aa9f52cee4a11c2764c4e745f47ddb1d5e1",
			1234567890.1234567890,
			false,
		},
		{
			"2. success encryption negative Float64",
			"5012ee13ae2de3646b72fd6ab0358aa9f52cee4a11c2764c4e745f47ddb1d5e1",
			-1234567890.1234567890,
			false,
		},
		{
			"3. success encryption Float64 = 0",
			"5012ee13ae2de3646b72fd6ab0358aa9f52cee4a11c2764c4e745f47ddb1d5e1",
			0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyBytes, err := hex.DecodeString(tt.keyHex)
			if err != nil {
				t.Errorf("Encryption key error: %v", err)
			}

			key := Key{}
			copy(key[:], keyBytes)

			x, err := EncryptFloat64(tt.num, &key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptFloat error = %v, wantErr %v", err, tt.wantErr)
			}

			x, err = DecryptFloat64(x, &key)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.num, x)
		})
	}
}
