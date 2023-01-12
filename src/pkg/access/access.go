package access

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/errors"
)

// Returns: id, key, level. error
func ExtractWithUserKey(item *Item, userPubKey, userPrivKey *[32]byte) ([]byte, *chachaPoly.Key, string, error) {
	levelBytes, ok := item.Fields["level"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'level' required", errors.ErrIncorrectFormat)
	}

	if len(levelBytes) == 0 {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'level' is empty", errors.ErrIncorrectFormat)
	}

	if levelBytes[0] == NoAccess {
		return nil, nil, "", errors.ErrAccessDenied
	}

	level := LevelToString(levelBytes[0])

	keyEncr, ok := item.Fields["keyEncr"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'keyEncr' required", errors.ErrIncorrectFormat)
	}

	//keyDecr, err := keybox.OpenAnonymous(keyEncr, userPubKey, userPrivKey)
	keyDecr, err := keybox.OpenAnonymous(keyEncr, userPubKey, userPrivKey)
	if err != nil {
		return nil, nil, "", fmt.Errorf("keybox.Open error: %w keyEncr: %x", err, keyEncr)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return nil, nil, "", fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	idEncr, ok := item.Fields["idEncr"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'idEncr' required", errors.ErrIncorrectFormat)
	}

	idDecr, err := key.Decrypt(idEncr)
	if err != nil {
		log.Printf("keyEncr: %x key: %x IDEncr: %x", keyEncr, key[:], idEncr)

		return nil, nil, "", fmt.Errorf("keybox.Open error: %w idDecr: %x", err, idEncr)
	}

	idDecrHash := crypto.Keccak256(idDecr)

	idHash, ok := item.Fields["idHash"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'idHash' required", errors.ErrIncorrectFormat)
	}

	if !bytes.Equal(idHash, idDecrHash) {
		return nil, nil, "", fmt.Errorf("%w: mismatch idHash: %x keccak256(idDecr): %x", errors.ErrCustom, idHash, idDecrHash)
	}

	return idDecr, key, level, nil
}

// Returns: id, key, level. error
func ExtractWithGroupKey(item *Item, groupKey *chachaPoly.Key) ([]byte, *chachaPoly.Key, string, error) {
	keyEncr, ok := item.Fields["keyEncr"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'keyEncr' required", errors.ErrIncorrectFormat)
	}

	keyDecr, err := groupKey.Decrypt(keyEncr)
	if err != nil {
		return nil, nil, "", fmt.Errorf("groupKey.Decrypt error: %w keyEncr: %x", err, keyEncr)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return nil, nil, "", fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	idEncr, ok := item.Fields["idEncr"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'idEncr' required", errors.ErrIncorrectFormat)
	}

	idDecr, err := key.Decrypt(idEncr)
	if err != nil {
		log.Printf("keyEncr: %x key: %x IDEncr: %x", keyEncr, key[:], idEncr)

		return nil, nil, "", fmt.Errorf("key error: %w idDecr: %x", err, idEncr)
	}

	idDecrHash := crypto.Keccak256(idDecr)

	idHash, ok := item.Fields["idHash"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'idHash' required", errors.ErrIncorrectFormat)
	}

	if !bytes.Equal(idHash, idDecrHash) {
		return nil, nil, "", fmt.Errorf("%w: mismatch idHash: %x keccak256(idDecr): %x", errors.ErrCustom, idHash, idDecrHash)
	}

	levelBytes, ok := item.Fields["level"]
	if !ok {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'level' required", errors.ErrIncorrectFormat)
	}

	if len(levelBytes) == 0 {
		return nil, nil, "", fmt.Errorf("%w ACL filed 'level' is empty", errors.ErrIncorrectFormat)
	}

	level := LevelToString(levelBytes[0])

	return idDecr, key, level, nil
}
