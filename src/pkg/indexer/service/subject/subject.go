// Package subject Store and receive EHR documents ID by subject
package subject

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"hms/gateway/pkg/indexer"
)

type SubjectIndex struct {
	index indexer.Indexer
}

func New() *SubjectIndex {
	return &SubjectIndex{
		index: indexer.Init("subject"),
	}
}

// AddEhrSubjectsIndex Add EHR document ID to index by document subject
func (e *SubjectIndex) AddEhrSubjectsIndex(ehrId, subjectId, namespace string) (err error) {
	subjectKey := e.subjectKey(subjectId, namespace)
	err = e.index.Replace(subjectKey, ehrId)
	return
}

// GetEhrBySubject Get EHR document ID by document subject
func (e *SubjectIndex) GetEhrBySubject(subjectId, namespace string) (ehrId string, err error) {
	subjectKey := e.subjectKey(subjectId, namespace)
	err = e.index.GetById(subjectKey, &ehrId)
	return
}

// Create document key by document subject
func (e *SubjectIndex) subjectKey(subjectId, namespace string) string {
	subjectKey := sha3.Sum256([]byte(subjectId + namespace))
	return hex.EncodeToString(subjectKey[:])
}
