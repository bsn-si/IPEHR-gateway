package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"

	"golang.org/x/crypto/scrypt"
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

	_, userPrivKey, err := infra.Keystore.Get(cfg.DefaultUserID)
	if err != nil {
		log.Fatalf("Keystore.Get error: %v userID %s", err, cfg.DefaultUserID)
	}

	pwdHash, err := generateHashFromPassword(cfg.CreatingSystemID, cfg.DefaultUserID, "")
	if err != nil {
		log.Fatalf("generateHashFromPassword error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	txHash, err := infra.Index.UserNew(ctx, cfg.DefaultUserID, cfg.CreatingSystemID, uint8(roles.Patient), pwdHash, nil, userPrivKey, nil)
	if err != nil {
		log.Fatalf("Index.UserAdd error: %v", err)
	}

	log.Println("txHash:", txHash)
}

func generateHashFromPassword(ehrSystemID, userID, password string) ([]byte, error) {
	salt := make([]byte, common.ScryptSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("rand.Read error: %w", err)
	}

	password = strings.Join([]string{userID, ehrSystemID, password}, "")

	pwdHash, err := scrypt.Key([]byte(password), salt, common.ScryptN, common.ScryptR, common.ScryptP, common.ScryptKeyLen)
	if err != nil {
		return nil, fmt.Errorf("generateHash scrypt.Key error: %w", err)
	}

	return append(pwdHash, salt...), nil
}
