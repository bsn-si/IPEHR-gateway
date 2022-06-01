// Package ehr_subject_index Store and receive EHR documents ID by subject
package ehr_subject

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"hms/gateway/pkg/indexer"
)

type EhrSubjectIndex struct {
	index indexer.Indexer
}

func New() *EhrSubjectIndex {
	return &EhrSubjectIndex{
		index: indexer.Init("subject"),
	}
}

// AddEhrSubjectsIndex Add EHR document ID to index by document subject
func (e *EhrSubjectIndex) AddEhrSubjectsIndex(ehrId, subjectId, namespace string) (err error) {
	subjectKey := e.subjectKey(subjectId, namespace)
	err = e.index.Replace(subjectKey, ehrId)
	return
}

// GetEhrBySubject Get EHR document ID by document subject
func (e *EhrSubjectIndex) GetEhrBySubject(subjectId, namespace string) (ehrId string, err error) {
	subjectKey := e.subjectKey(subjectId, namespace)
	err = e.index.GetById(subjectKey, &ehrId)
	return
}

// Create document key by document subject
func (e *EhrSubjectIndex) subjectKey(subjectId, namespace string) string {
	subjectKey := sha3.Sum256([]byte(subjectId + namespace))
	return hex.EncodeToString(subjectKey[:])
}
