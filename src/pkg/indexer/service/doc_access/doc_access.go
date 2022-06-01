// Package doc_access user - storage document index
package doc_access

import (
	"hms/gateway/pkg/indexer"
)

type DocAccessIndex struct {
	index indexer.Indexer
}

func New() *DocAccessIndex {
	return &DocAccessIndex{
		index: indexer.Init("ehrs"),
	}
}

// Add document storage Id for user
func (u *DocAccessIndex) Add(userId string, docStorageId *[32]byte) error {
	return u.index.Add(userId, docStorageId)
}

// Get document storage Id for user
func (u *DocAccessIndex) Get(userId string) (docStorageId *[32]byte, err error) {
	err = u.index.GetById(userId, docStorageId)
	return
}
