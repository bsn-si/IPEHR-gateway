package treeindex

import (
	"crypto/rand"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/hm"
	"github.com/stretchr/testify/assert"
)

func Test_EncryptComposition(t *testing.T) {
	filePath := "./../../../../data/mock/ehr/composition.json"
	composition, err := loadComposition(filePath)
	assert.Nil(t, err)

	node, err := ProcessComposition(&composition)
	assert.Nil(t, err)

	var key hm.Key
	rand.Read(key[:]) //nolint

	var nonce hm.Nonce
	rand.Read(nonce[:]) //nolint

	nodeEncrypted, err := EncryptNode(node, &key, &nonce)
	assert.Nil(t, err)

	//todo decrypt and compare
	_ = nodeEncrypted
}
