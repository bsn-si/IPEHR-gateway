// Package ehrs user - storage document index
package ehrs

import (
	"hms/gateway/pkg/indexer"
)

type Index struct {
	index indexer.Indexer
}

func New() *Index {
	return &Index{
		index: indexer.Init("ehrs"),
	}
}

// Add document storage Id for user
func (u *Index) Add(userID string, docStorageID *[32]byte) error {
	return u.index.Add(userID, docStorageID)
}

// Get document storage Id for user
func (u *Index) Get(userID string) (docStorageID *[32]byte, err error) {
	docStorageID = &[32]byte{}
	err = u.index.GetByID(userID, docStorageID)

	return
}

// Replace document storage Id for user
func (u *Index) Replace(userID string, docStorageID *[32]byte) error {
	return u.index.Replace(userID, docStorageID)
}
