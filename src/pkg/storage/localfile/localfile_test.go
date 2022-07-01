package localfile_test

import (
	"bytes"
	"os"
	"testing"

	"hms/gateway/pkg/common/utils"
	config2 "hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/localfile"
)

func TestWithCompression(t *testing.T) {
	cfg := config()

	globalConfig, err := config2.New()
	if err != nil {
		t.Fatal(err)
	}

	globalConfig.CompressionEnabled = true

	fs, err := localfile.Init(cfg, globalConfig)
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

func TestWithoutCompression(t *testing.T) {
	cfg := config()

	globalConfig, err := config2.New()
	if err != nil {
		t.Fatal(err)
	}

	globalConfig.CompressionEnabled = false

	fs, err := localfile.Init(cfg, globalConfig)
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
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return
	}

	filePath := rootDir + "/data/mock/ehr/composition.json"

	return os.ReadFile(filePath)
}
