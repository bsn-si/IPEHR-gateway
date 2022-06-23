package localfile

import (
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/utils"
	"os"
	"testing"
)

func TestWithCompression(t *testing.T) {
	cfg := config()

	fs, err := Init(cfg)
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

	if !common.SliceEqualBytes(data, data2) {
		t.Fatal("Data mismatch")
	}

	if err = os.RemoveAll(cfg.BasePath); err != nil {
		t.Fatal(err)
	}
}

func config() *Config {
	return &Config{
		BasePath: "/tmp/localfiletest",
		Depth:    3,
	}
}

func testData() (data []byte, err error) {
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return
	}
	filePath := rootDir + "/data/mock/ehr/composition.json"

	data, err = os.ReadFile(filePath)

	return
}
