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

	nodeEncrypted, err := ExecNode(node, func(node *DataValueNode) error {
		return EncryptDataValueNode(node, &key, &nonce)
	})
	assert.Nil(t, err)

	//todo resolve DvDateTime timezone problem and then possible to compare results

	nodeDecrypted, err := ExecNode(nodeEncrypted, func(node *DataValueNode) error {
		return DecryptDataValueNode(node, &key, &nonce)
	})
	assert.Nil(t, err)

	_ = nodeDecrypted
}
