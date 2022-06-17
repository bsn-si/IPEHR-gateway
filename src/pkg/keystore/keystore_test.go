package keystore

import (
	"hms/gateway/pkg/storage"
	"os"
	"strconv"
	"testing"
	"time"
)

const testStorePath string = "/tmp/localfiletest"

func TestKeystore(t *testing.T) {
	defer func() {
		err := cleanup()
		if err != nil {
			t.Fatal(err)
		}
	}()

	sc := &storage.StorageConfig{}
	sc.New("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	ks := New()

	userIdOne := "111-222-333"
	userIdTwo := "111-222-333-444"

	publicKeyOne, privateKeyOne, err := ks.Get(userIdOne)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyOne2, privateKeyOne2, err := ks.Get(userIdOne)
	if err != nil {
		t.Fatal(err)
	}

	if *publicKeyOne != *publicKeyOne2 || *privateKeyOne != *privateKeyOne2 {
		t.Fatal("Got different keys for same user")
	}

	publicKeyTwo, privateKeyTwo, err := ks.Get(userIdTwo)
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
