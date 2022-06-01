// Package docs EHR documents index
package docs

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/indexer"
)

type DocsIndex struct {
	index indexer.Indexer
}

func New() *DocsIndex {
	return &DocsIndex{
		index: indexer.Init("docs"),
	}
}

// Add EHR documents
func (d *DocsIndex) Add(ehrId string, docIndexes []model.DocumentMeta) error {
	return d.index.Add(ehrId, docIndexes)
}

// Replace EHR documents
func (d *DocsIndex) Replace(ehrId string, docIndexes []*model.DocumentMeta) error {
	return d.index.Replace(ehrId, docIndexes)
}

// Get EHR documents metadata
func (d *DocsIndex) Get(ehrId string) (docIndexes []*model.DocumentMeta, err error) {
	err = d.index.GetById(ehrId, &docIndexes)
	return
}
