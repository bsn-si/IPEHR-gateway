package subject_test

import (
	"strconv"
	"testing"
	"time"

	"hms/gateway/pkg/indexer/service/subject"
	"hms/gateway/pkg/storage"
)

func TestSubjectIndex(t *testing.T) {
	sc := storage.NewConfig("./test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	storage.Init(sc)

	subjectIndex := subject.New()

	testEhrID := "This is the best of test EHR document ID"
	testSubject := "test_subject"
	testNamespace := "test_namespace"
	testBadNamespace := "bad_test_namespace"

	err := subjectIndex.AddEhrSubjectsIndex(testEhrID, testSubject, testNamespace)
	if err != nil {
		t.Fatal(err)
	}

	receivedEhrID, err := subjectIndex.GetEhrBySubject(testSubject, testNamespace)
	if err != nil {
		t.Fatal(err)
	}

	if testEhrID != receivedEhrID {
		t.Errorf("Saved %s and received %s EHR ids not match", testEhrID, receivedEhrID)
	}

	_, err = subjectIndex.GetEhrBySubject(testSubject, testBadNamespace)
	if err == nil {
		t.Error("Not got error for wrong subject namespace")
	}
}
