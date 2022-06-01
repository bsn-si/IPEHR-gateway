// Package doc EHR documents index
package doc

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/indexer"
)

type DocIndex struct {
	index indexer.Indexer
}

func New() *DocIndex {
	return &DocIndex{
		index: indexer.Init("docs"),
	}
}

// Add EHR documents
func (d *DocIndex) Add(ehrId string, docIndexes []model.DocumentMeta) error {
	return d.index.Add(ehrId, docIndexes)
}

// Replace EHR documents
func (d *DocIndex) Replace(ehrId string, docIndexes []*model.DocumentMeta) error {
	return d.index.Replace(ehrId, docIndexes)
}

// Get EHR documents metadata
func (d *DocIndex) Get(ehrId string) (docIndexes []*model.DocumentMeta, err error) {
	err = d.index.GetById(ehrId, &docIndexes)
	return
}
