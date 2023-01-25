package access

import "github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"

type Item struct {
	ID     []byte
	Key    *chachaPoly.Key
	Level  string
	Fields map[string][]byte
}

type List []*Item
