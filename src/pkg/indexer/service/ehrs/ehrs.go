// Package ehrs user - storage document index
package ehrs

import (
	"hms/gateway/pkg/indexer"
)

type EhrsIndex struct {
	index indexer.Indexer
}

func New() *EhrsIndex {
	return &EhrsIndex{
		index: indexer.Init("ehrs"),
	}
}

// Add document storage Id for user
func (u *EhrsIndex) Add(userId string, docStorageId *[32]byte) error {
	return u.index.Add(userId, docStorageId)
}

// Get document storage Id for user
func (u *EhrsIndex) Get(userId string) (docStorageId *[32]byte, err error) {
	err = u.index.GetById(userId, docStorageId)
	return
}
