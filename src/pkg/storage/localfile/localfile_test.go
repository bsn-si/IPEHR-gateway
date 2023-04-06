package localfile_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/localfile"
)

func TestLocalfileStorage(t *testing.T) {
	cfg := config()

	fs, err := localfile.Init(cfg)
	if err != nil {
		t.Fatal(err)
	}

	data, err := testData()
	if err != nil {
		t.Fatal(err)
	}

	id, err := fs.Add(data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := fs.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, data2) {
		t.Fatal("Data mismatch")
	}

	if err = os.RemoveAll(cfg.BasePath); err != nil {
		t.Fatal(err)
	}
}

func config() *localfile.Config {
	return &localfile.Config{
		BasePath: "/tmp/localfiletest",
		Depth:    3,
	}
}

func testData() (data []byte, err error) {
	filePath := "./test_fixtures/composition.json"

	return os.ReadFile(filePath)
}
