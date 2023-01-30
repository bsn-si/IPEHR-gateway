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

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/utils"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

type bodyFillers map[base.ItemType]string

type testWrapDirectory struct {
	user        *User
	doctor      *Doctor
	ehrID       string
	ehrSystemID string
	*testWrap
}

func (w *testWrapDirectory) getURL() string {
	return w.server.URL + "/v1/ehr/" + w.ehrID + "/directory"
}

func (w *testWrapDirectory) prepare(testData *TestData, t *testing.T) {

	err := w.checkUser(testData)
	if err != nil {
		t.Fatal("Check user error:", err)
	}

	if testData.users[0].ehrID == "" {
		w.ehrCreate(testData)(t)
	}

	w.user = testData.users[0]
	w.ehrID = w.user.ehrID
	w.ehrSystemID = testData.ehrSystemID

	if len(testData.doctors) == 0 {
		w.doctorRegister(testData)(t)
	}

	if testData.doctors[0].accessToken == "" {
		err := testData.doctors[0].login(testData.ehrSystemID, w.server.URL, w.httpClient)
		if err != nil {
			t.Fatal("User login error:", err)
		}
	}

	w.doctor = testData.doctors[0]
}

func (testWrap *testWrap) directoryCRUD(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {

		wrap := &testWrapDirectory{testWrap: testWrap}
		wrap.prepare(testData, t)

		// TODO who will pay for it patient or doctor?
		body, err := directoryWithEmptyBody()
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot create composition body request"))
		}

		d, err := createDirectory(wrap, body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "Cant create DIRECTORY"))
		}

		testData.directory = d

		dCreated, err := getDirectoryByVersion(wrap, d.UID.Value, d.Name.Value)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, d, dCreated)

		dByTime, err := getDirectoryByTime(wrap, d.Name.Value, time.Now().Format(common.OpenEhrTimeFormat))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, dCreated, dByTime)

		body, err = directoryWithItemComposition(wrap.ehrID, d.UID.Value)
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot create composition body request with data"))
		}

		dVersionUID, err := base.NewObjectVersionID(d.UID.Value, wrap.ehrSystemID)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dVersionUID.IncreaseUIDVersion()
		if err != nil {
			t.Fatal(err)
		}

		_, err = getDirectoryByVersion(wrap, dVersionUID.String(), d.Name.Value)
		if err == nil {
			t.Fatal(errors.New(fmt.Sprintf("Directory with version %s not created yet", dVersionUID.String())))
		}

		dUpdated, err := updateDirectory(wrap, body, d.UID.Value)
		if err != nil {
			t.Fatal(errors.Wrap(err, "Cant update DIRECTORY"))
		}

		if !assert.Equal(t, dVersionUID.String(), dUpdated.UID.Value) {
			t.Fatal(errors.New("Version os DIRECTORYes not equal"))
		}

		err = deleteDirectory(t, wrap, d.UID.Value)
		if err == nil {
			t.Fatal(errors.New("Should be invoked error because DIRECTORY version mismatched"))
		}

		err = deleteDirectory(t, wrap, dUpdated.UID.Value)
		if err != nil {
			t.Fatal(errors.New("Cant delete DIRECTORY by last version"))
		}

		_, err = getDirectoryByVersion(wrap, d.UID.Value, d.Name.Value)
		if err != nil {
			t.Fatal(errors.New(fmt.Sprintf("Directory with version %s should be exist", d.UID.Value)))
		}

		_, err = getDirectoryByVersion(wrap, dUpdated.UID.Value, d.Name.Value)
		if err == nil {
			t.Fatal(errors.New(fmt.Sprintf("Directory with version %s should be already deleted", dUpdated.UID.Value)))
		}
	}
}

func createDirectory(wrap *testWrapDirectory, body *bytes.Reader) (*model.Directory, error) {
	link := wrap.getURL() + "/?patient_id=" + wrap.user.id

	request, err := http.NewRequest(http.MethodPost, link, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", wrap.doctor.id)
	request.Header.Set("Authorization", "Bearer "+wrap.doctor.accessToken)
	//request.Header.Set("GroupAccessId", groupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", wrap.ehrSystemID)

	response, err := wrap.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}

	var d model.Directory
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal DIRECTORY model")
	}

	requestID := response.Header.Get("RequestId")

	err = requestWait(wrap.doctor.id, wrap.doctor.accessToken, requestID, wrap.server.URL, wrap.httpClient)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func updateDirectory(wrap *testWrapDirectory, body *bytes.Reader, versionID string) (*model.Directory, error) {
	link := wrap.getURL() + "/?patient_id=" + wrap.user.id

	request, err := http.NewRequest(http.MethodPut, link, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", wrap.doctor.id)
	request.Header.Set("Authorization", "Bearer "+wrap.doctor.accessToken)
	//request.Header.Set("GroupAccessId", groupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", wrap.ehrSystemID)
	request.Header.Set("If-Match", versionID)

	response, err := wrap.httpClient.Do(request)
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

	err = requestWait(wrap.doctor.id, wrap.doctor.accessToken, requestID, wrap.server.URL, wrap.httpClient)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func getDirectoryByVersion(wrap *testWrapDirectory, versionID, path string) (*model.Directory, error) {
	link := fmt.Sprintf("%s/%s/?&patient_id=%s&path=%s", wrap.getURL(), versionID, wrap.user.id, path)

	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", wrap.doctor.id)
	request.Header.Set("Authorization", "Bearer "+wrap.doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", wrap.ehrSystemID)

	response, err := wrap.httpClient.Do(request)
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

func getDirectoryByTime(wrap *testWrapDirectory, path, versionAtTime string) (*model.Directory, error) {
	link := fmt.Sprintf("%s/?&patient_id=%s&path=%s&version_at_time=%s", wrap.getURL(), wrap.user.id, path, url.QueryEscape(versionAtTime))

	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", wrap.doctor.id)
	request.Header.Set("Authorization", "Bearer "+wrap.doctor.accessToken)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", wrap.ehrSystemID)

	response, err := wrap.httpClient.Do(request)
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

func deleteDirectory(t *testing.T, wrap *testWrapDirectory, versionID string) error {
	link := fmt.Sprintf("%s/?&patient_id=%s", wrap.getURL(), wrap.user.id)

	request, err := http.NewRequest(http.MethodDelete, link, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", wrap.doctor.id)
	request.Header.Set("Authorization", "Bearer "+wrap.doctor.accessToken)
	//request.Header.Set("GroupAccessId", groupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", wrap.ehrSystemID)
	request.Header.Set("If-Match", versionID)

	response, err := wrap.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return errors.New(response.Status)
	}

	requestID := response.Header.Get("RequestId")

	err = requestWait(wrap.doctor.id, wrap.doctor.accessToken, requestID, wrap.server.URL, wrap.httpClient)
	if err != nil {
		return err
	}

	t.Logf("Directory deleted, got new location: %s", response.Header.Get("Location"))

	return nil
}

func directoryWithEmptyBody() (*bytes.Reader, error) {
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
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return nil, err
	}

	filePath := rootDir + "/data/mock/ehr/" + mockDirectoryFileName + ".json"

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
