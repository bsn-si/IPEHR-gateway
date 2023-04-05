package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	docModel "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"
	"github.com/bsn-si/IPEHR-gateway/src/tests/api/testhelpers"
)

type User struct {
	id           string
	password     string
	accessToken  string
	refreshToken string
	ehrID        string
	ehrStatusID  string
	compositions []*docModel.Composition
	templates    []*docModel.Template
}

type Doctor struct {
	User
	Name       string
	Address    string
	Descrition string
	PictureURL string
	Code       string
}

func (u *User) login(ehrSystemID, baseURL string, client *http.Client) error {
	authRequest := model.UserAuthRequest{
		UserID:   u.id,
		Password: u.password,
	}

	body, _ := json.Marshal(authRequest)

	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/user/login", bytes.NewReader(body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = response.Body.Close(); err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return err
	}

	jwt := model.JWT{}
	if err = json.Unmarshal(content, &jwt); err != nil {
		return err
	}

	u.accessToken = jwt.AccessToken
	u.refreshToken = jwt.RefreshToken

	return nil
}

func userLogin(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		userHelper := testhelpers.UserHelper{}

		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		tests := []struct {
			name           string
			action         string
			method         string
			useAuthHeaders bool
			request        *model.UserAuthRequest
			statusCode     int
		}{
			{
				name:       "Empty userID and password",
				action:     "login",
				request:    userHelper.UserAuthRequest(),
				statusCode: http.StatusBadRequest,
			},
			{
				name:       "Empty userID",
				action:     "login",
				request:    userHelper.UserAuthRequest(userHelper.WithPassword("password")),
				statusCode: http.StatusBadRequest,
			},
			{
				name:       "Empty password",
				action:     "login",
				request:    userHelper.UserAuthRequest(userHelper.WithUserID(uuid.New().String())),
				statusCode: http.StatusBadRequest,
			},
			{
				name:   "UserID not exist",
				action: "login",
				request: userHelper.UserAuthRequest(
					userHelper.WithUserID(uuid.New().String()),
					userHelper.WithPassword("password")),
				statusCode: http.StatusNotFound,
			},
			{
				name:   "Password incorrect",
				action: "login",
				request: userHelper.UserAuthRequest(
					userHelper.WithUserID(user.id),
					userHelper.WithPassword("incorrect")),
				statusCode: http.StatusUnauthorized,
			},
			{
				name:   "Successfully auth",
				action: "login",
				request: userHelper.UserAuthRequest(
					userHelper.WithUserID(user.id),
					userHelper.WithPassword(user.password)),
				statusCode: http.StatusOK,
			},
			{
				name:           "Fail if already logged",
				action:         "login",
				useAuthHeaders: true,
				request: userHelper.UserAuthRequest(
					userHelper.WithUserID(user.id),
					userHelper.WithPassword(user.password)),
				statusCode: http.StatusUnprocessableEntity,
			},
			{
				name:   "Refresh token",
				action: "refresh",
				method: http.MethodGet,
				request: userHelper.UserAuthRequest(
					userHelper.WithUserID(user.id)),
				statusCode: http.StatusOK,
			},
			{
				name:           "Successfully logout",
				action:         "logout",
				useAuthHeaders: true,
				request: userHelper.UserAuthRequest(
					userHelper.WithUserID(user.id)),
				statusCode: http.StatusOK,
			},
		}

		var (
			jwt    model.JWT
			result = true
		)

		for _, data := range tests {
			docBytes, _ := json.Marshal(data.request)

			httpMethod := http.MethodPost

			if data.method != "" {
				httpMethod = data.method
			}

			if data.action == "logout" {
				docBytes, _ = json.Marshal(jwt)
			}

			request, err := http.NewRequest(httpMethod, testData.serverURL+"/v1/user/"+data.action, bytes.NewReader(docBytes))
			if err != nil {
				t.Fatal(err)
			}

			if data.useAuthHeaders || data.action != "login" {
				request.Header.Set("AuthUserId", data.request.UserID)

				if data.action == "refresh" {
					request.Header.Set("Authorization", "Bearer "+user.refreshToken)
				} else {
					request.Header.Set("Authorization", "Bearer "+user.accessToken)
				}
			}

			request.Header.Set("Content-type", "application/json")
			request.Header.Set("Prefer", "return=representation")
			request.Header.Set("EhrSystemId", testData.ehrSystemID)

			response, err := testData.httpClient.Do(request)
			if err != nil {
				t.Fatalf("Expected nil, received %s", err.Error())
			}

			content, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Response body read error: %v", err)
			}

			err = response.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("User login test: %s, response: %s", data.name, content)

			if response.StatusCode != data.statusCode {
				if result {
					result = false
				}

				t.Errorf("Test: %s, Expected: %d, received: %d, body: %s", data.name, data.statusCode, response.StatusCode, content)

				continue
			}

			if (data.action == "login" || data.action == "refresh") && response.StatusCode == http.StatusOK {
				if err = json.Unmarshal(content, &jwt); err != nil {
					t.Fatal(err)
				}

				user.accessToken = jwt.AccessToken
				user.refreshToken = jwt.RefreshToken
			}

			if data.action == "logout" {
				user.accessToken = ""
				user.refreshToken = ""
			}
		}

		if !result {
			t.Fatal()
		}
	}
}

func doctorRegister(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		doctor := &Doctor{
			User: User{
				id:       uuid.New().String(),
				password: fakeData.GetRandomStringWithLength(10),
			},
			Name:       "Gregory House, M.D.",
			Address:    "Ann Arbor, Michigan, United States",
			Descrition: "Head of Diagnostic Medicine Nephrologist",
			PictureURL: "https://media.filmz.ru/photos/full/filmz.ru_f_48951.jpg",
		}

		reqID, err := registerDoctor(doctor, testData.ehrSystemID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatalf("Can not register user, err: %v", err)
		}

		err = requestWait(doctor.id, "", reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal("registerPatient requestWait error: ", err)
		}

		if doctor.accessToken == "" {
			err := doctor.login(testData.ehrSystemID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		// getting doctor code
		url := testData.serverURL + "/v1/user/" + doctor.id

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("AuthUserId", doctor.id)
		request.Header.Set("Authorization", "Bearer "+doctor.accessToken)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var doctor2 model.UserInfo

		err = json.Unmarshal(data, &doctor2)
		if err != nil {
			t.Fatal(err)
		}

		doctor.Code = doctor2.Code

		testData.doctors = append(testData.doctors, doctor)
	}
}

func userInfoDoctor(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		doctor, err := checkDoctor0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkDoctor0LoggedIn error:", err)
		}

		url := testData.serverURL + "/v1/user/" + doctor.id

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("AuthUserId", doctor.id)
		request.Header.Set("Authorization", "Bearer "+doctor.accessToken)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var doctor2 model.UserInfo

		err = json.Unmarshal(data, &doctor2)
		if err != nil {
			t.Fatal(err)
		}

		if doctor2.Name != doctor.Name {
			t.Fatalf("Expected Name: %s, received: %s", doctor.Name, doctor2.Name)
		}
	}
}

func userInfoPatient(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("empty test users litst")
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testData.serverURL, testData.httpClient)
			if err != nil {
				t.Fatal(err)
			}
		}

		url := testData.serverURL + "/v1/user/" + user.id

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var user2 model.UserInfo

		err = json.Unmarshal(data, &user2)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, user.ehrID, user2.EhrID.String())
		assert.Equal(t, "Patient", user2.Role)
	}
}

func userInfoByCode(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		doctor, err := checkDoctor0LoggedIn(testData)
		if err != nil {
			t.Fatal("checkDoctor0LoggedIn error:", err)
		}

		url := testData.serverURL + "/v1/user/code/" + doctor.Code

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", doctor.id)
		request.Header.Set("Authorization", "Bearer "+doctor.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testData.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected: %d, received: %d, body: %s", http.StatusOK, response.StatusCode, data)
		}

		var doctor2 model.UserInfo

		err = json.Unmarshal(data, &doctor2)
		if err != nil {
			t.Fatal(err)
		}

		if doctor2.Name != doctor.Name {
			t.Fatalf("Expected Name: %s, received: %s", doctor.Name, doctor2.Name)
		}
	}
}

func registerPatient(user *User, systemID, baseURL string, client *http.Client) (string, error) {
	userRegisterRequest, err := userCreateBodyRequest(user.id, user.password, "", "", "", "", roles.Patient)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/user/register", userRegisterRequest)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", systemID)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("register user resp error: %s", string(data)) // nolint
	}

	requestID := response.Header.Get("RequestId")
	return requestID, nil
}

func registerDoctor(d *Doctor, systemID, baseURL string, client *http.Client) (string, error) {
	doctorRegisterRequest, err := userCreateBodyRequest(d.id, d.password, d.Name, d.Address, d.Descrition, d.PictureURL, roles.Doctor)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/user/register", doctorRegisterRequest)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", systemID)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	err = response.Body.Close()
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusCreated {
		return "", err
	}

	requestID := response.Header.Get("RequestId")

	return requestID, nil
}

func userCreateBodyRequest(userID, password, name, address, description, pictureURL string, role roles.Role) (*bytes.Reader, error) {
	userRegisterRequest := &model.UserCreateRequest{
		UserID:      userID,
		Password:    password,
		Role:        uint8(role),
		Name:        name,
		Address:     address,
		Description: description,
		PictuteURL:  pictureURL,
	}

	docBytes, err := json.Marshal(userRegisterRequest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(docBytes), nil
}
