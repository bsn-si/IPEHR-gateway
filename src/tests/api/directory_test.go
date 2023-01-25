package api_test

import (
	"bytes"
	"encoding/json"
	"hms/gateway/pkg/common/utils"
	"io"
	"net/http"
	"os"
	"testing"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

//const version123 = "1.2.3"

func (testWrap *testWrap) directoryCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		err := testWrap.checkUser(testData)
		if err != nil {
			t.Fatal(err)
		}

		user := testData.users[0]
		if user.ehrID == "" {
			testWrap.ehrCreate(testData)(t)
		}

		if len(testData.doctors) == 0 {
			testWrap.doctorRegister(testData)(t)
		}

		doctor := testData.doctors[0]
		if doctor.accessToken == "" {
			err := doctor.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		// TODO who will pay for it?
		d, reqID, err := createDirectory(doctor, user, testData.ehrSystemID, doctor.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(doctor.id, doctor.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		testData.directories = append(testData.directories, d)
	}
}

func createDirectory(doctor *Doctor, user *User, ehrSystemID, accessToken, baseURL string, client *http.Client) (*model.Directory, string, error) {
	body, err := directoryWithEmptyBody(ehrSystemID)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot create composition body request")
	}

	url := baseURL + "/v1/ehr/" + user.ehrID + "/directory/" + user.id

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", doctor.id)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	//request.Header.Set("GroupAccessId", groupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot do create request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "connot read response body")
	}

	var d model.Directory
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, "", errors.Wrap(err, "cannot unmarshal DIRECTORY model")
	}

	requestID := response.Header.Get("RequestId")

	return &d, requestID, nil
}

func directoryWithEmptyBody(ehrSystemID string) (*bytes.Reader, error) {
	return directoryCreateBodyRequest(ehrSystemID, "directory_empty")
}

func directoryCreateBodyRequest(ehrSystemID, mockDirectoryFileName string) (*bytes.Reader, error) {
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return nil, err
	}

	filePath := rootDir + "/data/mock/ehr/" + mockDirectoryFileName + ".json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	//compositionID := uuid.New().String()

	//objectVersionID, err := base.NewObjectVersionID(compositionID, ehrSystemID)
	//if err != nil {
	//	log.Fatalf("Expected model.EHR, received %s", err.Error())
	//}
	//
	//data = []byte(strings.Replace(string(data), "__COMPOSITION_ID__", objectVersionID.String(), 1))

	return bytes.NewReader(data), nil
}
