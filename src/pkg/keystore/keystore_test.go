package keystore_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/keystore"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage"
)

const testStorePath string = "/tmp/localfiletest"

func TestKeystore(t *testing.T) {
	defer func() {
		err := cleanup()
		if err != nil {
			t.Fatal(err)
		}
	}()

	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	ks := keystore.New(cfg.KeystoreKey)

	userIDOne := "111-222-333"
	userIDTwo := "111-222-333-444"

	publicKeyOne, privateKeyOne, err := ks.Get(userIDOne)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyOne2, privateKeyOne2, err := ks.Get(userIDOne)
	if err != nil {
		t.Fatal(err)
	}

	if *publicKeyOne != *publicKeyOne2 || *privateKeyOne != *privateKeyOne2 {
		t.Fatal("Got different keys for same user")
	}

	publicKeyTwo, privateKeyTwo, err := ks.Get(userIDTwo)
	if err != nil {
		t.Fatal(err)
	}

	if *publicKeyOne == *publicKeyTwo || *privateKeyOne == *privateKeyTwo {
		t.Fatal("Got same keys for different user")
	}
}

func cleanup() (err error) {
	err = os.RemoveAll(testStorePath)
	return
}
