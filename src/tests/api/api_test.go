package api_test

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/api"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
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

var (
	ciRun        = flag.Bool("ci_run", false, "set true to use external server address")
	serverAddres = flag.String("server_address", "http://localhost:8080", "external test server address")
)

type testWrap struct {
	// server     *httptest.Server
	serverURL  string
	httpClient *http.Client
	// storage    *storage.Storager
}

func TestMain(m *testing.M) {
	flag.Parse()

	close := func() {}

	if !*ciRun {
		testServer, storager, err := prepareTest()
		if err != nil {
			log.Fatal(err)
		}

		*serverAddres = testServer.URL

		close = func() {
			tearDown(testServer, storager)
		}
	}

	defer close()

	os.Exit(m.Run())
}

func Test_API(t *testing.T) {
	testWrap := &testWrap{
		serverURL:  *serverAddres,
		httpClient: &http.Client{},
	}

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
		reqID, err := registerPatient(user, testData.ehrSystemID, testWrap.serverURL, testWrap.httpClient)
		if err != nil {
			t.Fatalf("Can not register user, err: %v", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testWrap.serverURL, testWrap.httpClient)
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

	if !t.Run("User get info doctor", testWrap.userInfoDoctor(testData)) {
		t.Fatal()
	}

	if !t.Run("User get info by code", testWrap.userInfoByCode(testData)) {
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

	if !t.Run("EHR_STATUS getting", testWrap.ehrStatusGet(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS getting by version time", testWrap.ehrStatusGetByVersionTime(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS update", testWrap.ehrStatusUpdate(testData)) {
		t.Fatal()
	}

	if !t.Run("User get info patient", testWrap.userInfoPatient(testData)) {
		t.Fatal()
	}

	/*
		if !t.Run("Document access set", testWrap.docSetAccessSuccess(testData)) {
			t.Fatal()
		}

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

	if !t.Run("COMPOSITION get list", testWrap.compositionGetList(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION delete by wrong UID", testWrap.compositionDeleteByWrongID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION delete", testWrap.compositionDeleteByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DIRECTORY create", testWrap.directoryCRUD(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a invalid query", testWrap.definitionStoreInvalidQuery(testData)) {
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

	if !t.Run("DEFINITION Template14 upload", testWrap.definitionTemplate14Upload(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Template14 get by ID", testWrap.definitionTemplate14GetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Template14 list stored", testWrap.definitionTemplate14List(testData)) {
		t.Fatal()
	}

	if !t.Run("QUERY execute with GET Expected success with correct query", testWrap.queryExecSuccess(testData)) {
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

	if !t.Run("User group get list", testWrap.userGroupGetList(testData)) {
		t.Fatal()
	}

	if !t.Run("User group remove user", testWrap.userGroupRemoveUser(testData)) {
		t.Fatal()
	}
}

func prepareTest() (*httptest.Server, storage.Storager, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, nil, errors.Wrap(err, "config.New() error")
	}

	cfg.Storage.Localfile.Path += "/test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	cfg.DefaultUserID = uuid.New().String()

	infra := infrastructure.New(cfg)

	apiHandler := api.New(cfg, infra)

	r := apiHandler.Build()
	srv := httptest.NewServer(r)

	return srv, storage.Storage(), nil
}

func tearDown(srv *httptest.Server, storager storage.Storager) {
	srv.Close()

	if err := (storager).Clean(); err != nil {
		log.Panicln(err)
	}
}
