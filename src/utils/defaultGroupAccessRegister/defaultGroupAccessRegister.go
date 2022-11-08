package main

import (
	"context"
	"crypto/rand"
	"flag"
	"hms/gateway/pkg/access"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/infrastructure"
	"log"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"
)

func main() {
	var (
		cfgPath = flag.String("config", "./config.json", "config file path")
	)

	flag.Parse()

	cfg, err := config.New(*cfgPath)
	if err != nil {
		panic(err)
	}

	infra := infrastructure.New(cfg)

	groupUUID, err := uuid.Parse(cfg.DefaultGroupAccessID)
	if err != nil {
		log.Fatal(err)
	}

	groupAccess := &model.GroupAccess{
		GroupUUID:   &groupUUID,
		Description: "Default access group",
		Key:         chachaPoly.GenerateKey(),
		Nonce:       &[12]byte{},
	}

	// nolint
	rand.Read(groupAccess.Nonce[:])

	groupAccessByte, err := msgpack.Marshal(groupAccess)
	if err != nil {
		log.Fatal("msgpack.Marshal error: ", err)
	}

	userPubKey, userPrivKey, err := infra.Keystore.Get(cfg.DefaultUserID)
	if err != nil {
		log.Fatalf("Keystore.Get error: %v userID %s", err, cfg.DefaultUserID)
	}

	groupAccessEncrypted, err := keybox.Seal(groupAccessByte, userPubKey, userPrivKey)
	if err != nil {
		log.Fatalf("keybox.SealAnonymous error: %v", err)
	}

	h := sha3.Sum256(append([]byte(cfg.DefaultUserID), groupUUID[:]...))

	ctx := context.Background()

	//nolint
	txHash, err := infra.Index.SetGroupAccess(ctx, &h, groupAccessEncrypted, uint8(access.Owner), userPrivKey, nil)
	if err != nil {
		log.Fatalf("Index.SetGroupAccess error: %v", err)
	}

	log.Println("txHash:", txHash)
}
