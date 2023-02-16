package query

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/hm"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func encryptQuery(query *aqlprocessor.Query, key *chachaPoly.Key, nonce *chachaPoly.Nonce) error {
	err := encryptWhere(query.Where, key, nonce)
	if err != nil {
		return fmt.Errorf("encryptWhere error: %w", err)
	}

	return nil
}

func encryptWhere(w *aqlprocessor.Where, key *chachaPoly.Key, nonce *chachaPoly.Nonce) error {
	if w == nil {
		return nil
	}

	err := encryptIdentifiedExpr(w.IdentifiedExpr, key, nonce)
	if err != nil {
		return fmt.Errorf("encryptIdentifiedExpr error: %w", err)
	}

	for _, next := range w.Next {
		err = encryptWhere(next, key, nonce)
		if err != nil {
			return fmt.Errorf("encryptWhere error: %w", err)
		}
	}

	return nil
}

func encryptIdentifiedExpr(ie *aqlprocessor.IdentifiedExpr, key *chachaPoly.Key, nonce *chachaPoly.Nonce) error {
	if ie == nil {
		return nil
	}

	err := encryptTerminal(ie.Terminal, key, nonce)
	if err != nil {
		return fmt.Errorf("encryptTerminal error: %w", err)
	}

	err = encryptIdentifiedExpr(ie.Next, key, nonce)
	if err != nil {
		return fmt.Errorf("encryptIdentifiedExpr error: %w", err)
	}

	return nil
}

func encryptTerminal(t *aqlprocessor.Terminal, key *chachaPoly.Key, nonce *chachaPoly.Nonce) error {
	if t == nil || t.Primitive == nil {
		return nil
	}

	var err error

	switch val := t.Primitive.Val.(type) {
	case float64:
		t.Primitive.Type = aqlprocessor.PrimitiveTypeBigFloat

		t.Primitive.Val, err = hm.EncryptFloat64(val, key)
		if err != nil {
			return fmt.Errorf("hm.EncryptFloat64 error: %w", err)
		}
	case int64:
		t.Primitive.Type = aqlprocessor.PrimitiveTypeBigInt

		t.Primitive.Val, err = hm.EncryptInt64(val, key)
		if err != nil {
			return fmt.Errorf("hm.EncryptInt64 error: %w", err)
		}
	case int:
		t.Primitive.Type = aqlprocessor.PrimitiveTypeBigInt

		t.Primitive.Val, err = hm.EncryptInt64(int64(val), key)
		if err != nil {
			return fmt.Errorf("hm.EncryptInt64 error: %w", err)
		}
	case string:
		t.Primitive.Type = aqlprocessor.PrimitiveTypeString

		t.Primitive.Val = hm.EncryptString(val, key, nonce)
	default:
		return errors.Errorf("Unsupported value %v type: %T", val, val)
	}

	return nil
}
