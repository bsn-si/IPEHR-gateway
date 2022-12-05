package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/api"
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/common/utils"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/storage"
	userModel "hms/gateway/pkg/user/model"
	userRoles "hms/gateway/pkg/user/roles"
	"hms/gateway/tests/api/testhelpers"
)

const (
	reqKindEhrCreate = iota
	//reqKindUserRegister
)

type User struct {
	id           string
	password     string
	accessToken  string
	refreshToken string
	ehrID        string
	ehrStatusID  string
	compositions []*model.Composition
}

type Request struct {
	id   string
	kind int
	user *User
}

type TestData struct {
	ehrSystemID   string
	users         []*User
	requests      []*Request
	groupsAccess  []*model.GroupAccess
	storedQueries []*model.StoredQuery
	userGroups    []*userModel.UserGroup
}

type testWrap struct {
	server     *httptest.Server
	httpClient *http.Client
	storage    *storage.Storager
}

func Test_API(t *testing.T) {
	testServer, storager := prepareTest(t)

	testWrap := &testWrap{
		server:     testServer,
		httpClient: &http.Client{},
		storage:    &storager,
	}
	defer tearDown(*testWrap)

	testData := &TestData{
		ehrSystemID: common.EhrSystemID,
		//nolint
		users: []*User{
			&User{id: uuid.New().String(), password: fakeData.GetRandomStringWithLength(10)},
			&User{id: uuid.New().String(), password: fakeData.GetRandomStringWithLength(10)},
			&User{id: uuid.New().String(), password: fakeData.GetRandomStringWithLength(10)},
		},
	}

	for _, user := range testData.users {
		reqID, err := registerUser(user, testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatalf("Can not register user, err: %v", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal("registerUser requestWait error: ", err)
		}
	}

	// TODO user register incorrect input data
	// TODO user register duplicate registration request
	//if !t.Run("User register", testWrap.userRegister(testData)) {
	//	t.Fatal()
	//}

	if !t.Run("User login", testWrap.userLogin(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR creating", testWrap.ehrCreate(testData)) {
		t.Fatal()
	}

	if !t.Run("Get transaction requests", testWrap.requests(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR creating with id", testWrap.ehrCreateWithID(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR creating with id for the same user", testWrap.ehrCreateWithIDForSameUser(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR getting", testWrap.ehrGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR get by subject", testWrap.ehrGetBySubject(testData)) {
		t.Fatal()
	}

	//if !t.Run("EHR grant access to another User", testWrap.docGrantAccessSuccess(testData)) {
	//}

	if !t.Run("EHR_STATUS getting", testWrap.ehrStatusGet(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS getting by version time", testWrap.ehrStatusGetByVersionTime(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS update", testWrap.ehrStatusUpdate(testData)) {
		t.Fatal()
	}

	/*
		if !t.Run("Access group create", testWrap.accessGroupCreate(testData)) {
			t.Fatal()
		}

		if !t.Run("Wrong access group getting", testWrap.wrongAccessGroupGetting(testData)) {
			t.Fatal()
		}

		if !t.Run("Access group getting", testWrap.accessGroupGetting(testData)) {
			t.Fatal()
		}
	*/

	if !t.Run("COMPOSITION create Expected fail with wrong EhrId", testWrap.compositionCreateFail(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION create Expected success with correct EhrId", testWrap.compositionCreateSuccess(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION getting with correct EhrId", testWrap.compositionGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION getting with wrong EhrId", testWrap.compositionGetByWrongID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION update", testWrap.compositionUpdate(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION delete by wrong UID", testWrap.compositionDeleteByWrongID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION delete", testWrap.compositionDeleteByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a query", testWrap.definitionStoreQuery(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a query version", testWrap.definitionStoreQueryVersion(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a query version with same ID", testWrap.definitionStoreQueryVersionWithSameID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Get stored query by ID", testWrap.definitionStoredQueryGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION List stored queries", testWrap.definitionListStoredQueries(testData)) {
		t.Fatal()
	}

	if !t.Run("QUERY execute with POST Expected success with correct query", testWrap.queryExecPostSuccess(testData)) {
		t.Fatal()
	}

	if !t.Run("QUERY execute with POST Expected fail with wrong query", testWrap.queryExecPostFail(testData)) {
		t.Fatal()
	}

	if !t.Run("User group create", testWrap.userGroupCreate(testData)) {
		t.Fatal()
	}

	if !t.Run("User group add user", testWrap.userGroupAddUser(testData)) {
		t.Fatal()
	}

	if !t.Run("User group get by ID", testWrap.userGroupGetByID(testData)) {
		t.Fatal()
	}
}

func (u *User) login(ehrSystemID, baseURL string, client *http.Client) error {
	userRegisterRequest, err := userCreateBodyRequest(u.id, u.password)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/user/login", userRegisterRequest)
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

	jwt := userModel.JWT{}
	if err = json.Unmarshal(content, &jwt); err != nil {
		return err
	}

	u.accessToken = jwt.AccessToken
	u.refreshToken = jwt.RefreshToken

	return nil
}

func prepareTest(t *testing.T) (ts *httptest.Server, storager storage.Storager) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		t.Fatal("config.New error:", err)
	}

	cfg.Storage.Localfile.Path += "/test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	cfg.DefaultUserID = uuid.New().String()

	infra := infrastructure.New(cfg)
	apiHandler := api.New(cfg, infra)

	r := apiHandler.Build()
	ts = httptest.NewServer(r)

	return ts, storage.Storage()
}

func tearDown(testWrap testWrap) {
	testWrap.server.Close()

	err := (*testWrap.storage).Clean()
	if err != nil {
		log.Panicln(err)
	}
}

func (testWrap *testWrap) requests(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		var req *Request

	loop:
		for _, r := range testData.requests {
			switch r.kind {
			case reqKindEhrCreate:
				req = r
				break loop
			default:
			}
		}

		if req == nil {
			t.Fatal("Request required")
		}

		if req.user.accessToken == "" {
			err := req.user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/requests/"+req.id, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", req.user.id)
		request.Header.Set("Authorization", "Bearer "+req.user.accessToken)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GetRequestById expected %d, received %d", http.StatusOK, response.StatusCode)
		}

		t.Log("Requests: GetAll")

		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/requests/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", req.user.id)
		request.Header.Set("Authorization", "Bearer "+req.user.accessToken)
		request.Header.Set("Prefer", "return=representation")

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GetAllRequests expected %d, received %d", http.StatusOK, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) userLogin(testData *TestData) func(t *testing.T) {
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
			request        *userModel.UserAuthRequest
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
			jwt    userModel.JWT
			result = true
		)

		for _, data := range tests {
			payload := getReaderJSONFrom(data.request)
			httpMethod := http.MethodPost

			if data.method != "" {
				httpMethod = data.method
			}

			if data.action == "logout" {
				payload = getReaderJSONFrom(jwt)
			}

			request, err := http.NewRequest(httpMethod, testWrap.server.URL+"/v1/user/"+data.action, payload)
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

			response, err := testWrap.httpClient.Do(request)
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

func (testWrap *testWrap) ehrCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		ehr, reqID, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		user.ehrID = ehr.EhrID.Value
		user.ehrStatusID = ehr.EhrStatus.ID.Value

		testData.requests = append(testData.requests, &Request{
			id:   reqID,
			kind: reqKindEhrCreate,
			user: user,
		})
	}
}

func (testWrap *testWrap) ehrCreateWithID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) < 2 {
			t.Fatal("Test user2 required")
		}

		user := testData.users[1]

		if user.accessToken == "" {
			err := user.login(testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
			if err != nil {
				t.Fatal("User login error:", err)
			}
		}

		ehrID2 := uuid.New().String()

		ehr, _, err := createEhrWithID(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, ehrID2, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		newEhrID := ehr.EhrID.Value
		if newEhrID != ehrID2 {
			t.Fatal("EhrID is not matched")
		}
	}
}

func (testWrap *testWrap) ehrCreateWithIDForSameUser(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		_, _, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err == nil {
			t.Fatal("Expected error, received EHR")
		}
	}
}

func (testWrap *testWrap) ehrGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+user.ehrID, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
			return
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Error(err)
			return
		}

		if user.ehrID != ehr.EhrID.Value {
			t.Error("EHR document mismatch")
			return
		}
	}
}

func (testWrap *testWrap) ehrGetBySubject(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		testEhrStatus, err := testWrap.getEhrStatus(user.ehrID, user.ehrStatusID, user.id, testData.ehrSystemID, user.accessToken)
		if err != nil {
			log.Fatalf("Expected model.EhrStatus, received %s", err.Error())
		}

		// Check document by subject
		url := testWrap.server.URL + "/v1/ehr?subject_id=" + testEhrStatus.Subject.ExternalRef.ID.Value + "&subject_namespace=" + testEhrStatus.Subject.ExternalRef.Namespace

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.Header.Set("Prefer", "return=representation")

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatal(err)
		}

		var ehrDoc model.EHR

		err = json.Unmarshal(data, &ehrDoc)
		if err != nil {
			t.Fatal(err)
		}

		if ehrDoc.EhrID.Value != user.ehrID {
			t.Fatalf("Expected %s, received %s", user.ehrID, ehrDoc.EhrID.Value)
		}
	}
}

func (testWrap *testWrap) ehrStatusGet(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", user.ehrID, user.ehrStatusID)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Response body read error: %v", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Fatal(err)
		}

		if ehrStatus.UID == nil || ehrStatus.UID.Value != user.ehrStatusID {
			t.Fatal("EHR_STATUS document mismatch")
		}
	}
}

func (testWrap *testWrap) ehrStatusGetByVersionTime(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]
		versionAtTime := time.Now()

		request, err := http.NewRequest(http.MethodGet, testWrap.server.URL+fmt.Sprintf("/v1/ehr/%s/ehr_status", user.ehrID), nil)
		if err != nil {
			t.Fatal(err)
		}

		q := request.URL.Query()
		q.Add("version_at_time", versionAtTime.Format(common.OpenEhrTimeFormat))

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)
		request.URL.RawQuery = q.Encode()

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d", http.StatusOK, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) ehrStatusUpdate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		// replace substring in ehrStatusID
		objectVersionID, err := base.NewObjectVersionID(user.ehrStatusID, testData.ehrSystemID)
		if err != nil {
			log.Fatalf("Expected model.EHR, received %s", err.Error())
		}

		_, err = objectVersionID.IncreaseUIDVersion()
		if err != nil {
			log.Fatalf("Expected model.EHR, received %s", err.Error())
		}

		newEhrStatusID := objectVersionID.String()

		req := []byte(fmt.Sprintf(`{
		  "_type": "EHR_STATUS",
		  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
		  "name": {
			"value": "EHR Status"
		  },
		  "uid": {
			"_type": "OBJECT_VERSION_ID",
			"value": "%s"
		  },
		  "subject": {
			"external_ref": {
			  "id": {
				"_type": "HIER_OBJECT_ID",
				"value": "324a4b23-623d-4213-cc1c-23f233b24234"
			  },
			  "namespace": "DEMOGRAPHIC",
			  "type": "PERSON"
			}
		  },
		  "other_details": {
			"_type": "ITEM_TREE",
			"archetype_node_id": "at0001",
			"name": {
			  "value": "Details"
			},
			"items": []
		  },
		  "is_modifiable": true,
		  "is_queryable": true
		}`, newEhrStatusID))

		url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status", user.ehrID)

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(req))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("If-Match", user.ehrStatusID)
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehrStatus model.EhrStatus
		if err = json.Unmarshal(data, &ehrStatus); err != nil {
			t.Fatal(err)
		}

		updatedEhrStatusID := response.Header.Get("ETag")

		if updatedEhrStatusID != newEhrStatusID {
			t.Log("Response body:", string(data))
			t.Fatal("EHR_STATUS uid in ETag mismatch")
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(user.id, user.accessToken, requestID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		// Checking EHR_STATUS changes
		request, err = http.NewRequest(http.MethodGet, testWrap.server.URL+"/v1/ehr/"+user.ehrID, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if err = response.Body.Close(); err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var ehr model.EHR
		if err = json.Unmarshal(data, &ehr); err != nil {
			t.Fatal(err)
		}

		if ehr.EhrStatus.ID.Value != updatedEhrStatusID {
			t.Fatalf("EHR_STATUS id mismatch. Expected %s, received %s", updatedEhrStatusID, ehr.EhrStatus.ID.Value)
		}
	}
}

func (testWrap *testWrap) compositionCreateFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		ehrID := uuid.New().String()
		groupAccessID := ""

		composition, _, err := createComposition(user.id, ehrID, testData.ehrSystemID, user.accessToken, groupAccessID, testWrap.server.URL, testWrap.httpClient)
		if err == nil {
			t.Fatalf("Expected error, received status: %v", composition)
		}
	}
}

func (testWrap *testWrap) compositionCreateSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		if len(testData.groupsAccess) == 0 {
			uuid := uuid.New()

			testData.groupsAccess = append(testData.groupsAccess, &model.GroupAccess{GroupUUID: &uuid})
		}

		ga := testData.groupsAccess[0]

		c, reqID, err := createComposition(user.id, user.ehrID, testData.ehrSystemID, user.accessToken, ga.GroupUUID.String(), testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatalf("Unexpected composition, received error: %v", err)
		}

		t.Logf("Waiting for request %s done", reqID)

		err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		user.compositions = append(user.compositions, c)
	}
}

func (testWrap *testWrap) compositionGetByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		if len(testData.groupsAccess) == 0 {
			uuid := uuid.New()

			testData.groupsAccess = append(testData.groupsAccess, &model.GroupAccess{GroupUUID: &uuid})
		}

		ga := testData.groupsAccess[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		comp := user.compositions[0]

		url := testWrap.server.URL + "/v1/ehr/" + user.ehrID + "/composition/" + comp.UID.Value

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("GroupAccessId", ga.GroupUUID.String())
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status: %v, received %v", http.StatusOK, response.StatusCode)
		}

		var composition model.Composition
		if err = json.Unmarshal(data, &composition); err != nil {
			t.Fatal(err)
		}

		if composition.UID.Value != comp.UID.Value {
			t.Fatalf("Expected %s, received %s", composition.UID.Value, comp.UID.Value)
		}
	}
}

func (testWrap *testWrap) compositionGetByWrongID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		wrongCompositionID := uuid.NewString() + "::" + testData.ehrSystemID + "::1"

		url := testWrap.server.URL + "/v1/ehr/" + user.ehrID + "/composition/" + wrongCompositionID

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status %d, received %d", http.StatusNotFound, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionUpdate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		if len(testData.groupsAccess) == 0 {
			uuid := uuid.New()

			testData.groupsAccess = append(testData.groupsAccess, &model.GroupAccess{GroupUUID: &uuid})
		}

		ga := testData.groupsAccess[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		comp := user.compositions[0]

		objectVersionID, err := base.NewObjectVersionID(comp.UID.Value, testData.ehrSystemID)
		if err != nil {
			t.Fatal(err)
		}

		comp.ObjectVersionID = *objectVersionID

		comp.Name.Value = "Updated text"
		updatedComposition, _ := json.Marshal(comp)

		url := testWrap.server.URL + "/v1/ehr/" + user.ehrID + "/composition/" + comp.ObjectVersionID.BasedID()

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(updatedComposition))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("GroupAccessId", ga.GroupUUID.String())
		request.Header.Set("If-Match", comp.ObjectVersionID.String())
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("Prefer", "return=representation")
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		err = response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status: %v, received %v", http.StatusOK, response.StatusCode)
		}

		compositionUpdated := model.Composition{}
		if err = json.Unmarshal(data, &compositionUpdated); err != nil {
			t.Fatal(err)
		}

		if compositionUpdated.UID.Value == comp.UID.Value {
			t.Fatalf("Expected %s, received %s", compositionUpdated.UID.Value, comp.UID.Value)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(user.id, user.accessToken, requestID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (testWrap *testWrap) compositionDeleteByWrongID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + "/v1/ehr/" + user.ehrID + "/composition/" + uuid.New().String()

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status: %v, received %v", http.StatusNotFound, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) compositionDeleteByID(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		if len(user.compositions) == 0 {
			t.Fatal("Composition required")
		}

		comp := user.compositions[0]

		url := testWrap.server.URL + "/v1/ehr/" + user.ehrID + "/composition/" + comp.UID.Value

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)
		request.Header.Set("EhrSystemId", testData.ehrSystemID)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status: %v, received %v", http.StatusNoContent, response.StatusCode)
		}

		requestID := response.Header.Get("RequestId")

		t.Logf("Waiting for request %s done", requestID)

		err = requestWait(user.id, user.accessToken, requestID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Checking the status of a re-request to remove")

		response, err = testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status: %v, received %v", http.StatusBadRequest, response.StatusCode)
		}
	}
}

func (testWrap *testWrap) queryExecPostSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(user.ehrID))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected success, received status: %d", response.StatusCode)
		}
	}
}

func (testWrap *testWrap) queryExecPostFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		url := testWrap.server.URL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("111qqqEEE")))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected fail, received status: %d", response.StatusCode)
		}
	}
}

func userCreateBodyRequest(userID, password string) (*bytes.Reader, error) {
	userRegisterRequest := &userModel.UserCreateRequest{
		UserID:   userID,
		Password: password,
		Role:     uint8(userRoles.Patient),
	}

	docBytes, err := json.Marshal(userRegisterRequest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(docBytes), nil
}

func getReaderJSONFrom(data interface{}) *bytes.Reader {
	docBytes, _ := json.Marshal(data)

	return bytes.NewReader(docBytes)
}

func ehrCreateBodyRequest() *bytes.Reader {
	req := fakeData.EhrCreateRequest()
	return bytes.NewReader(req)
}

func compositionCreateBodyRequest(ehrSystemID string) (*bytes.Reader, error) {
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return nil, err
	}

	filePath := rootDir + "/data/mock/ehr/composition.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	compositionID := uuid.New().String()

	objectVersionID, err := base.NewObjectVersionID(compositionID, ehrSystemID)
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	_, err = objectVersionID.IncreaseUIDVersion()
	if err != nil {
		log.Fatalf("Expected model.EHR, received %s", err.Error())
	}

	data = []byte(strings.Replace(string(data), "__COMPOSITION_ID__", objectVersionID.String(), 1))

	return bytes.NewReader(data), nil
}

/*
func (testWrap *testWrap) accessGroupCreate(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		ga, err := createGroupAccess(user.id, user.accessToken, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatalf("Expected group access, received error: %v", err)
		}

		testData.groupsAccess = append(testData.groupsAccess, ga)
	}
}

func (testWrap *testWrap) wrongAccessGroupGetting(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		groupAccessIDWrong, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}

		url := testWrap.server.URL + "/v1/access/group/" + groupAccessIDWrong.String()

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusNotFound, response.StatusCode, data)
		}
	}
}

func (testWrap *testWrap) accessGroupGetting(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		if len(testData.groupsAccess) == 0 {
			t.Fatal("GroupAccess required")
		}

		ga := testData.groupsAccess[0]

		url := testWrap.server.URL + "/v1/access/group/" + ga.GroupUUID.String()

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Fatalf("Expected nil, received %s", err.Error())
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Response body read error: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d, received %d body: %s", http.StatusOK, response.StatusCode, data)
		}

		var groupAccessGot model.GroupAccess
		if err = json.Unmarshal(data, &groupAccessGot); err != nil {
			t.Fatal(err)
		}

		if ga.GroupUUID.String() != groupAccessGot.GroupUUID.String() {
			t.Fatal("Got wrong group")
		}
	}
}
*/

func requestWait(userID, accessToken, requestID, baseURL string, client *http.Client) error {
	request, err := http.NewRequest(http.MethodGet, baseURL+"/v1/requests/"+requestID, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)

	if accessToken != "" {
		request.Header.Set("Authorization", "Bearer "+accessToken)
	}

	timeout := time.Now().Add(1 * time.Minute)

	for {
		time.Sleep(2 * time.Second)

		if time.Now().After(timeout) {
			return errors.ErrTimeout
		}

		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: request %s getting error: %v", errors.ErrCustom, requestID, response.Status)
		}

		var request processing.RequestResult

		if err = json.NewDecoder(response.Body).Decode(&request); err != nil {
			return err
		}

		if request.Ethereum[0].StatusStr == processing.StatusSuccess.String() {
			return nil
		} else if request.Ethereum[0].StatusStr == processing.StatusFailed.String() {
			return errors.New("Request failed")
		}
	}
}

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}

func (testWrap *testWrap) getEhrStatus(ehrID, statusID, userID, ehrSystemID, accessToken string) (*model.EhrStatus, error) {
	url := testWrap.server.URL + fmt.Sprintf("/v1/ehr/%s/ehr_status/%s", ehrID, statusID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := testWrap.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var ehrStatus model.EhrStatus
	if err = json.Unmarshal(data, &ehrStatus); err != nil {
		return nil, err
	}

	return &ehrStatus, err
}

func createEhr(userID, ehrSystemID, accessToken, baseURL string, client *http.Client) (ehr *model.EHR, requestID string, err error) {
	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/ehr", ehrCreateBodyRequest())
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	if err = json.Unmarshal(data, &ehr); err != nil {
		return nil, "", err
	}

	requestID = response.Header.Get("RequestId")

	return ehr, requestID, nil
}

func createEhrWithID(userID, ehrSystemID, accessToken, baseURL, ehrID string, client *http.Client) (ehr *model.EHR, requestID string, err error) {
	request, err := http.NewRequest(http.MethodPut, baseURL+"/v1/ehr/"+ehrID, ehrCreateBodyRequest())
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	if err = json.Unmarshal(data, &ehr); err != nil {
		return nil, "", err
	}

	requestID = response.Header.Get("RequestId")

	return ehr, requestID, nil
}

/*
func createGroupAccess(userID, accessToken, baseURL string, client *http.Client) (*model.GroupAccess, error) {
	description := fakeData.GetRandomStringWithLength(50)

	req := []byte(`{
			"description": "` + description + `"
		}`)

	request, err := http.NewRequest(http.MethodPost, baseURL+"/v1/access/group", bytes.NewReader(req))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(response.Status)
	}

	var groupAccess model.GroupAccess
	if err = json.Unmarshal(data, &groupAccess); err != nil {
		return nil, err
	}

	return &groupAccess, nil
}
*/

// nolint
func createComposition(userID, ehrID, ehrSystemID, accessToken, groupAccessID, baseURL string, client *http.Client) (*model.Composition, string, error) {
	body, err := compositionCreateBodyRequest(ehrSystemID)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannnot create composition body request")
	}

	url := baseURL + "/v1/ehr/" + ehrID + "/composition"

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("AuthUserId", userID)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	//request.Header.Set("GroupAccessId", groupAccessID)
	request.Header.Set("Prefer", "return=representation")
	request.Header.Set("EhrSystemId", ehrSystemID)

	response, err := client.Do(request)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot do create composition request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, "", errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "connot read response body")
	}

	var c model.Composition
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, "", errors.Wrap(err, "cannot unmarshal COMPOSITION mondel")
	}

	requestID := response.Header.Get("RequestId")

	return &c, requestID, nil
}

func registerUser(user *User, systemID, baseURL string, client *http.Client) (string, error) {
	userRegisterRequest, err := userCreateBodyRequest(user.id, user.password)
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
