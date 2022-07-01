// Package docs EHR documents index
package docs

import (
	"time"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
)

type Index struct {
	index indexer.Indexer
}

func New() *Index {
	return &Index{
		index: indexer.Init("docs"),
	}
}

// Add doc index
func (i *Index) Add(ehrID string, docIndex *model.DocumentMeta) error {
	var docIndexes []*model.DocumentMeta

	err := i.index.GetByID(ehrID, &docIndexes)
	if err != nil && !errors.Is(err, errors.ErrIsNotExist) {
		return err
	}

	docIndexes = append(docIndexes, docIndex)

	return i.index.Replace(ehrID, docIndexes)
}

// Replace EHR documents
func (i *Index) Replace(ehrID string, docIndexes []*model.DocumentMeta) error {
	return i.index.Replace(ehrID, docIndexes)
}

// Get EHR documents metadata
func (i *Index) Get(ehrID string) (docIndexes []*model.DocumentMeta, err error) {
	err = i.index.GetByID(ehrID, &docIndexes)
	return
}

func (i *Index) GetByType(ehrID string, docType types.DocumentType) (docs []*model.DocumentMeta, err error) {
	docIndexes, err := i.Get(ehrID)
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docIndex.TypeCode == docType {
			docs = append(docs, docIndex)
		}
	}

	if 0 == len(docs) {
		return nil, errors.ErrIsNotExist
	}

	return docs, nil
}

func (i *Index) GetLastByType(ehrID string, docType types.DocumentType) (doc *model.DocumentMeta, err error) {
	docIndexes, err := i.Get(ehrID)
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docIndex.TypeCode == docType {
			if doc == nil || docIndex.Timestamp >= doc.Timestamp {
				doc = docIndex
			}
		}
	}

	if doc == nil {
		return nil, errors.ErrIsNotExist
	}

	return doc, nil
}

func (i *Index) GetDocIndexByNearestTime(ehrID string, nearestTime time.Time, docType types.DocumentType) (doc *model.DocumentMeta, err error) {
	docIndexes, err := i.Get(ehrID)
	if err != nil {
		return nil, err
	}

	t := uint64(nearestTime.Unix())

	for _, docIndex := range docIndexes {
		if docIndex.TypeCode == docType {
			if docIndex.Timestamp <= t {
				doc = docIndex
			} else {
				break
			}
		}
	}

	if doc == nil {
		err = errors.ErrIsNotExist
	}

	return doc, err
}
