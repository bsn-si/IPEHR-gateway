package access

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func ExtractWithUserKey(item *Item, userPubKey, userPrivKey *[32]byte) error {
	levelBytes, ok := item.Fields["level"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'level' required", errors.ErrIncorrectFormat)
	}

	if len(levelBytes) == 0 {
		return fmt.Errorf("%w ACL filed 'level' is empty", errors.ErrIncorrectFormat)
	}

	if levelBytes[0] == NoAccess {
		return errors.ErrAccessDenied
	}

	item.Level = LevelToString(levelBytes[0])

	keyEncr, ok := item.Fields["keyEncr"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'keyEncr' required", errors.ErrIncorrectFormat)
	}

	keyDecr, err := keybox.OpenAnonymous(keyEncr, userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("keybox.Open error: %w keyEncr: %x", err, keyEncr)
	}

	item.Key, err = chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	idEncr, ok := item.Fields["idEncr"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'idEncr' required", errors.ErrIncorrectFormat)
	}

	item.ID, err = item.Key.Decrypt(idEncr)
	if err != nil {
		log.Printf("keyEncr: %x key: %x IDEncr: %x", keyEncr, item.Key[:], idEncr)

		return fmt.Errorf("id decryption error: %w idDecr: %x", err, idEncr)
	}

	idDecrHash := crypto.Keccak256(item.ID)

	idHash, ok := item.Fields["idHash"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'idHash' required", errors.ErrIncorrectFormat)
	}

	if !bytes.Equal(idHash, idDecrHash) {
		return fmt.Errorf("%w: mismatch idHash: %x keccak256(idDecr): %x", errors.ErrCustom, idHash, idDecrHash)
	}

	return nil
}

func ExtractWithGroupKey(item *Item, groupKey *chachaPoly.Key) error {
	keyEncr, ok := item.Fields["keyEncr"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'keyEncr' required", errors.ErrIncorrectFormat)
	}

	keyDecr, err := groupKey.Decrypt(keyEncr)
	if err != nil {
		return fmt.Errorf("groupKey.Decrypt error: %w keyEncr: %x", err, keyEncr)
	}

	item.Key, err = chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	idEncr, ok := item.Fields["idEncr"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'idEncr' required", errors.ErrIncorrectFormat)
	}

	item.ID, err = item.Key.Decrypt(idEncr)
	if err != nil {
		log.Printf("keyEncr: %x key: %x IDEncr: %x", keyEncr, item.Key[:], idEncr)

		return fmt.Errorf("key error: %w idDecr: %x", err, idEncr)
	}

	idDecrHash := crypto.Keccak256(item.ID)

	idHash, ok := item.Fields["idHash"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'idHash' required", errors.ErrIncorrectFormat)
	}

	if !bytes.Equal(idHash, idDecrHash) {
		return fmt.Errorf("%w: mismatch idHash: %x keccak256(idDecr): %x", errors.ErrCustom, idHash, idDecrHash)
	}

	levelBytes, ok := item.Fields["level"]
	if !ok {
		return fmt.Errorf("%w ACL filed 'level' required", errors.ErrIncorrectFormat)
	}

	if len(levelBytes) == 0 {
		return fmt.Errorf("%w ACL filed 'level' is empty", errors.ErrIncorrectFormat)
	}

	item.Level = LevelToString(levelBytes[0])

	return nil
}
