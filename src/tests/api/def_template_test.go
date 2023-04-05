package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	docModel "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func definitionTemplate14Upload(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		templateID := fakeData.GetRandomStringWithLength(10)

		tmpl, reqID, err := uploadTemplate14(user.id, testData.ehrSystemID, templateID, user.accessToken, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatalf("Unexpected template upload, received error: %v", err)
		}

		t.Logf("Waiting for request %s done", reqID)

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		user.templates = append(user.templates, tmpl)
	}
}

func definitionTemplate14GetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(user.templates) == 0 {
			definitionTemplate14Upload(testData)(t)
		}

		tmpl1 := user.templates[0]

		url := testData.serverURL + "/v1/definition/template/adl1.4/" + tmpl1.TemplateID

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("Accept", docModel.ADLTypeXML)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		if !bytes.Equal(tmpl1.Body, data) {
			t.Fatalf("Expected same template with length %d, received template length %d", len(tmpl1.Body), len(data))
		}
	}
}

func definitionTemplate14List(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		if len(user.templates) == 0 {
			definitionTemplate14Upload(testData)(t)
		}

		tmpl1 := user.templates[0]

		url := testData.serverURL + "/v1/definition/template/adl1.4"

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("ConvertType", "application/json")

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var list []docModel.Template

		err = json.Unmarshal(data, &list)
		if err != nil {
			t.Fatal(err)
		}

		if list[0].TemplateID != tmpl1.TemplateID {
			t.Fatalf("Expected: %s, received: %s, body: %s", tmpl1.TemplateID, list[0].TemplateID, data)
		}
	}
}

func uploadTemplate14(userID, ehrSystemID, templateID, accessToken, baseURL string, client *http.Client) (*docModel.Template, string, error) {
	filePath := "./test_fixtures/template14.xml"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	data = bytes.Replace(data, []byte("__TEMPLATE_ID__"), []byte(templateID), 1)

	url := baseURL + "/v1/definition/template/adl1.4"

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/xml")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot do upload template request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	requestID := response.Header.Get("RequestId")

	tmpl := &docModel.Template{
		TemplateID: templateID,
	}

	tmpl.Body = data

	return tmpl, requestID, nil
}
