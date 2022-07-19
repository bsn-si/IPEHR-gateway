package filecoin_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/filecoin"
)

func TestStartDeal(t *testing.T) {

	storage, err := prepare(t)
	if err != nil {
		t.Fatal(err)
	}
	defer clean(t, storage)

	testData, err := fakeData.GetByteArray(1024)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	deal, err := storage.StartDeal(ctx, testData)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(deal)
}

func prepare(t *testing.T) (*filecoin.Storage, error) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	// For testing purposes
	cfg.Storage.Filecoin.FilesPath = "/tmp/filecoin.tmp." + fakeData.GetRandomStringWithLength(8)

	_, err = os.Stat(cfg.Storage.Filecoin.FilesPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(cfg.Storage.Filecoin.FilesPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	filecoinCfg := (filecoin.Config)(cfg.Storage.Filecoin)

	storage, err := filecoin.New(&filecoinCfg)
	if err != nil {
		return nil, fmt.Errorf("filecoin.New error: %w", err)
	}

	return storage, nil
}

func clean(t *testing.T, s *filecoin.Storage) error {
	t.Helper()

	s.Close()

	err := os.RemoveAll(s.FilesPath())
	if err != nil {
		t.Error(err)
		return err
	}

	return nil
}
