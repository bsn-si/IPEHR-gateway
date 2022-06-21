// Package data_search keys index
package data_search

import (
	"encoding/hex"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"

	"golang.org/x/crypto/sha3"
)

type DataSearchEntry struct {
	GroupId               *[16]byte
	ValueEncrypted        []byte
	DocStorageIdEncrypted []byte
}

// Index sha3(pathKey) -> DataSearchEntry
type DataSearchIndex struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New(ks *keystore.KeyStore) *DataSearchIndex {
	return &DataSearchIndex{
		index:    indexer.Init("data_search"),
		keystore: ks,
	}
}

func (d *DataSearchIndex) Add(pathKey string, dataEntry *DataSearchEntry) error {
	var (
		indexKey    = sha3.Sum256([]byte(pathKey))
		indexKeyStr = hex.EncodeToString(indexKey[:])
		err         = d.index.Add(indexKeyStr, dataEntry)
	)
	return err
}

// Get DataAccessIndex key
func (d *DataSearchIndex) Get(pathKey string) (*DataSearchEntry, error) {
	var (
		dataEntry   DataSearchEntry
		indexKey    = sha3.Sum256([]byte(pathKey))
		indexKeyStr = hex.EncodeToString(indexKey[:])
		err         = d.index.GetById(indexKeyStr, &dataEntry)
	)
	return &dataEntry, err
}
