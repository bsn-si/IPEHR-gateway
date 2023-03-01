package models

import (
	"testing"
)

func TestIndexChunk(t *testing.T) {
	idx := NewIndexChunk("group_id", "data_id", "ehr_id", []byte("data"))

	if !idx.Validate() {
		t.Error("chunk is invalid. expected valid")
	}

	idx.GroupID = "new_group_id"
	if idx.Validate() {
		t.Error("chunk is valid. expected invalid")
	}
}
