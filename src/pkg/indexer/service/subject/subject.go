// Package subject Store and receive EHR documents ID by subject
package subject

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/indexer"
)

type Index struct {
	index indexer.Indexer
}

func New() *Index {
	return &Index{
		index: indexer.Init("subject"),
	}
}

// AddEhrSubjectsIndex Add EHR document ID to index by document subject
func (e *Index) AddEhrSubjectsIndex(ehrID, subjectID, namespace string) (err error) {
	subjectKey := e.subjectKey(subjectID, namespace)
	err = e.index.Replace(subjectKey, ehrID)

	return
}

// GetEhrBySubject Get EHR document ID by document subject
func (e *Index) GetEhrBySubject(subjectID, namespace string) (ehrID string, err error) {
	subjectKey := e.subjectKey(subjectID, namespace)
	err = e.index.GetByID(subjectKey, &ehrID)

	return
}

// Create document key by document subject
func (e *Index) subjectKey(subjectID, namespace string) string {
	subjectKey := sha3.Sum256([]byte(subjectID + namespace))

	return hex.EncodeToString(subjectKey[:])
}
