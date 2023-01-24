package storage_test

import (
	"fmt"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage"
)

type pathTest struct {
	path, expected string
}

//nolint:gofmt
var storagePathTests = []pathTest{
	pathTest{"test", "%s/test"},
	pathTest{"test/1", "%s/test/1"},
	pathTest{"./test/1", "%s/test/1"},
	pathTest{"   ./test/1/   ", "%s/test/1"},
	pathTest{".////////test/1///////////////", "%s/test/1"},
}

func TestStorageConfig(t *testing.T) {
	for _, test := range storagePathTests {
		testStorageConfig := storage.NewConfig(test.path)

		test.expected = fmt.Sprintf(test.expected, testStorageConfig.ProcessPath())
		if testStorageConfig.Path() != test.expected {
			t.Errorf("Output %q not equal to expected %q", testStorageConfig.Path(), test.expected)
		}
	}
}

func TestStorageConfigWithRoot(t *testing.T) {
	defer func() { _ = recover() }()

	storage.NewConfig("/")

	t.Errorf("Root folder is free for use as a storage")
}

func TestStorageConfigWithRelativeRoot(t *testing.T) {
	defer func() { _ = recover() }()

	storage.NewConfig("../../../../../../../../../../../../../../")

	t.Errorf("Root folder is free for use as a storage if path is relative")
}

func TestStorageConfigWhichNotInExecutedFolder(t *testing.T) {
	var path = "/test/1/2/3"

	testStorageConfig := storage.NewConfig(path)
	if testStorageConfig.Path() != path {
		t.Errorf("Out of executed folder test failed")
	}
}
