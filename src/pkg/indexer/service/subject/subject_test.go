package subject

import (
	"testing"
)

func TestSubjectIndex(t *testing.T) {
	subjectIndex := New()

	testEhrId := "This is the best of test EHR document ID"
	testSubject := "test_subject"
	testNamespace := "test_namespace"
	testBadNamespace := "bad_test_namespace"

	err := subjectIndex.AddEhrSubjectsIndex(testEhrId, testSubject, testNamespace)
	if err != nil {
		t.Fatal(err)
	}

	receivedEhrId, err := subjectIndex.GetEhrBySubject(testSubject, testNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if testEhrId != receivedEhrId {
		t.Errorf("Saved %s and recived %s EHR ids not match", testEhrId, receivedEhrId)
	}

	receivedEhrId, err = subjectIndex.GetEhrBySubject(testSubject, testBadNamespace)
	if err == nil {
		t.Error("Not got error for wrong subject namespace")
	}
}
