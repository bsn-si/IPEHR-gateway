package parser

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestParseEhr(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := wd + "/../../../../data/mock/ehr/ehr.json"

	inJson, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Can't open composition.json file", filePath)
	}

	res, err := ParseEhr(inJson)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = uuid.Parse(res.EhrId.Value); err != nil {
		t.Fatal("EHR Document is not parsed correctly")
	}
}

func TestParseComposition(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := wd + "/../../../../data/mock/ehr/composition.json"

	inJson, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Can't open composition.json file", filePath)
	}

	res, err := ParseComposition(inJson)
	if err != nil {
		t.Fatal(err)
	}

	if res.Uid.Value == "" {
		t.Fatal("Composition is not parsed correctly")
	}
}
