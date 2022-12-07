package api_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/api"
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/storage"
	userModel "hms/gateway/pkg/user/model"
)

type TestData struct {
	ehrSystemID   string
	users         []*User
	requests      []*Request
	groupsAccess  []*model.GroupAccess
	storedQueries []*model.StoredQuery
	userGroups    []*userModel.UserGroup
	doctors       []*Doctor
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
		reqID, err := registerPatient(user, testData.ehrSystemID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatalf("Can not register user, err: %v", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.server.URL, testWrap.httpClient)
		if err != nil {
			t.Fatal("registerPatient requestWait error: ", err)
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

	if !t.Run("User register with doctor role", testWrap.doctorRegister(testData)) {
		t.Fatal()
	}

	if !t.Run("User get info", testWrap.userInfo(testData)) {
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
