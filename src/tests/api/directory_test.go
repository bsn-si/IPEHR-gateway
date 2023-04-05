package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type bodyFillers map[base.ItemType]string

func directoryCRUD(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := checkUser0LoggedInAndEhrCreated(testData)
		if err != nil {
			t.Fatal("checkUser0LoggedInAndEhrCreated error:", err)
		}

		// TODO who will pay for it patient or doctor?
		body, err := directoryCreateWithEmptyBody()
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot create composition body request"))
		}

		d, err := directoryCreate(testData, body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "Cant create DIRECTORY"))
		}

		dCreated, err := getDirectoryByVersion(testData, d.UID.Value, d.Name.Value)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, d, dCreated)

		dByTime, err := getDirectoryByTime(testData, d.Name.Value, time.Now().Format(common.OpenEhrTimeFormat))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, dCreated, dByTime)

		body, err = directoryWithItemComposition(user.ehrID, d.UID.Value)
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot create composition body request with data"))
		}

		dVersionUID, err := base.NewObjectVersionID(d.UID.Value, testData.ehrSystemID)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dVersionUID.IncreaseUIDVersion()
		if err != nil {
			t.Fatal(err)
		}

		_, err = getDirectoryByVersion(testData, dVersionUID.String(), d.Name.Value)
		if err == nil {
			t.Fatal("Directory with version not created yet", dVersionUID.String())
		}

		dUpdated, err := updateDirectory(testData, body, d.UID.Value)
		if err != nil {
			t.Fatal(errors.Wrap(err, "Cant update DIRECTORY"))
		}

		if !assert.Equal(t, dVersionUID.String(), dUpdated.UID.Value) {
			t.Fatal(errors.New("Version os DIRECTORYes not equal"))
		}

		err = deleteDirectory(testData, d.UID.Value)
		if err == nil {
			t.Fatal(errors.New("Should be invoked error because DIRECTORY version mismatched"))
		}

		err = deleteDirectory(testData, dUpdated.UID.Value)
		if err != nil {
			t.Fatal(errors.New("Cant delete DIRECTORY by last version"))
		}

		_, err = getDirectoryByVersion(testData, d.UID.Value, d.Name.Value)
		if err != nil {
			t.Fatal("Directory with version should be exist", d.UID.Value)
		}

		_, err = getDirectoryByVersion(testData, dUpdated.UID.Value, d.Name.Value)
		if err == nil {
			t.Fatal("Directory with version should be already deleted", dUpdated.UID.Value)
		}
	}
}

func directoryCreate(testData *TestData, body *bytes.Reader) (*model.Directory, error) {
	user, err := checkUser0LoggedInAndEhrCreated(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser0LoggedInAndEhrCreated error: %w", err)
	}

	doctor, err := checkDoctor0LoggedIn(testData)
	if err != nil {
		return nil, fmt.Errorf("checkDoctor0LoggedIn error: %w", err)
	}

	url := fmt.Sprintf("%s/v1/ehr/%s/directory?patient_id=%s", testData.serverURL, user.ehrID, user.id)

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", doctor.id)
	request.Header.Set("Authorization", "Bearer "+doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := testData.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(response.Status)
	}

	var d model.Directory
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal DIRECTORY model")
	}

	requestID := response.Header.Get("RequestId")

	err = requestWait(doctor.id, doctor.accessToken, requestID, testData.serverURL, testData.httpClient)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func updateDirectory(testData *TestData, body *bytes.Reader, versionID string) (*model.Directory, error) {
	user, err := checkUser0LoggedInAndEhrCreated(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser0LoggedInAndEhrCreated error: %w", err)
	}

	doctor, err := checkDoctor0LoggedIn(testData)
	if err != nil {
		return nil, fmt.Errorf("checkDoctor0LoggedIn error: %w", err)
	}

	url := fmt.Sprintf("%s/v1/ehr/%s/directory?patient_id=%s", testData.serverURL, user.ehrID, user.id)

	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", doctor.id)
	request.Header.Set("Authorization", "Bearer "+doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", testData.ehrSystemID)
	request.Header.Set("If-Match", versionID)

	response, err := testData.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status + ", body contain: " + string(data))
	}

	var d model.Directory
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal DIRECTORY model")
	}

	requestID := response.Header.Get("RequestId")

	err = requestWait(doctor.id, doctor.accessToken, requestID, testData.serverURL, testData.httpClient)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func getDirectoryByVersion(testData *TestData, versionID, path string) (*model.Directory, error) {
	user, err := checkUser0LoggedInAndEhrCreated(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser0LoggedInAndEhrCreated error: %w", err)
	}

	doctor, err := checkDoctor0LoggedIn(testData)
	if err != nil {
		return nil, fmt.Errorf("checkDoctor0LoggedIn error: %w", err)
	}

	url := fmt.Sprintf("%s/v1/ehr/%s/directory/%s?patient_id=%s&path=%s", testData.serverURL, user.ehrID, versionID, user.id, path)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", doctor.id)
	request.Header.Set("Authorization", "Bearer "+doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := testData.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Sprintf("response with status - %s, body contain %s", response.Status, data)
		return nil, errors.New(err)
	}

	var d model.Directory
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal DIRECTORY model")
	}

	return &d, nil
}

func getDirectoryByTime(testData *TestData, path, versionAtTime string) (*model.Directory, error) {
	user, err := checkUser0LoggedInAndEhrCreated(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser0LoggedInAndEhrCreated error: %w", err)
	}

	doctor, err := checkDoctor0LoggedIn(testData)
	if err != nil {
		return nil, fmt.Errorf("checkDoctor0LoggedIn error: %w", err)
	}

	url := fmt.Sprintf("%s/v1/ehr/%s/directory?patient_id=%s&path=%s&version_at_time=%s", testData.serverURL, user.ehrID, user.id, path, url.QueryEscape(versionAtTime))

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", doctor.id)
	request.Header.Set("Authorization", "Bearer "+doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", testData.ehrSystemID)

	response, err := testData.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Sprintf("response with status - %s, body contain %s", response.Status, data)
		return nil, errors.New(err)
	}

	var d model.Directory
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal DIRECTORY model")
	}

	return &d, nil
}

func deleteDirectory(testData *TestData, versionID string) error {
	user, err := checkUser0LoggedInAndEhrCreated(testData)
	if err != nil {
		return fmt.Errorf("checkUser0LoggedInAndEhrCreated error: %w", err)
	}

	doctor, err := checkDoctor0LoggedIn(testData)
	if err != nil {
		return fmt.Errorf("checkDoctor0LoggedIn error: %w", err)
	}

	url := fmt.Sprintf("%s/v1/ehr/%s/directory/?&patient_id=%s", testData.serverURL, user.ehrID, user.id)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", doctor.id)
	request.Header.Set("Authorization", "Bearer "+doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", testData.ehrSystemID)
	request.Header.Set("If-Match", versionID)

	response, err := testData.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return errors.New(response.Status)
	}

	requestID := response.Header.Get("RequestId")

	err = requestWait(doctor.id, doctor.accessToken, requestID, testData.serverURL, testData.httpClient)
	if err != nil {
		return fmt.Errorf("requestWait error: %w", err)
	}

	return nil
}

func directoryCreateWithEmptyBody() (*bytes.Reader, error) {
	return directoryCreateBodyRequest(bodyFillers{}, "directory_empty")
}

func directoryWithItemComposition(ehrSystemID, dUUID string) (*bytes.Reader, error) {
	f := bodyFillers{
		base.CompositionItemType:     ehrSystemID,
		base.ObjectVersionIDItemType: dUUID,
	}

	return directoryCreateBodyRequest(f, "directory_with_items")
}

func directoryCreateBodyRequest(filler bodyFillers, mockDirectoryFileName string) (*bytes.Reader, error) {
	filePath := "./test_fixtures/" + mockDirectoryFileName + ".json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(filler) != 0 {
		for k, v := range filler {
			data = []byte(strings.Replace(string(data), "__"+k.ToString()+"__", v, -1))
		}
	}

	return bytes.NewReader(data), nil
}
