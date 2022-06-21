package localfile

import (
	"hms/gateway/pkg/common"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	cfg := Config{
		BasePath: "/tmp/localfiletest",
		Depth:    3,
	}

	fs, err := Init(&cfg)
	if err != nil {
		t.Error(err)
	}

	data := []byte("Hello! This is the test data for BSN HMS Gateway")

	id, err := fs.Add(data)
	if err != nil {
		t.Error(err)
	}

	data2, err := fs.Get(id)
	if err != nil {
		t.Error(err)
	}

	if !common.SliceEqualBytes(data, data2) {
		t.Errorf("Data mismatch")
	}

	if err = os.RemoveAll(cfg.BasePath); err != nil {
		t.Error(err)
	}
}
