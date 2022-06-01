// Package docs EHR documents index
package docs

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
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

// Add doc index
func (d *DocsIndex) Add(ehrId string, docIndex *model.DocumentMeta) error {
	var docIndexes []*model.DocumentMeta
	err := d.index.GetById(ehrId, &docIndexes)
	if err != nil && !errors.Is(err, errors.IsNotExist) {
		return err
	}
	docIndexes = append(docIndexes, docIndex)
	return d.index.Replace(ehrId, docIndexes)
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

func (d *DocsIndex) GetLastByType(ehrId string, docType types.DocumentType) (doc *model.DocumentMeta, err error) {
	docIndexes, err := d.Get(ehrId)
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docIndex.TypeCode == docType {
			if doc == nil || docIndex.Timestamp > doc.Timestamp {
				doc = docIndex
			}
		}
	}
	if doc == nil {
		return nil, errors.IsNotExist
	}
	return doc, nil
}
