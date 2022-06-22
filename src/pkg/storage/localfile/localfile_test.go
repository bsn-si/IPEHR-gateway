package localfile

import (
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/utils"
	config2 "hms/gateway/pkg/config"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	cfg := config()

	globalConfig, err := config2.New()
	if err != nil {
		t.Fatal(err)
	}

	globalConfig.CompressionEnabled = true

	fs, err := Init(cfg, globalConfig)
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

func TestWithoutCompression(t *testing.T) {
	cfg := config()

	globalConfig, err := config2.New()
	if err != nil {
		t.Fatal(err)
	}

	globalConfig.CompressionEnabled = false

	fs, err := Init(cfg, globalConfig)
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
