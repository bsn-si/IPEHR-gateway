package parser_test

import (
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/parser"
)

func TestParseEhr(t *testing.T) {
	filePath := "./test_fixtures/ehr.json"

	inJSON, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Can't open ehr.json file", filePath)
	}

	res, err := parser.ParseEhr(inJSON)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = uuid.Parse(res.EhrID.Value); err != nil {
		t.Fatal("EHR Document is not parsed correctly")
	}
}
