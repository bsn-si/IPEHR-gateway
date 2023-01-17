package access

import "hms/gateway/pkg/crypto/chachaPoly"

type Item struct {
	ID     []byte
	Key    *chachaPoly.Key
	Level  string
	Fields map[string][]byte
}

type List []*Item
