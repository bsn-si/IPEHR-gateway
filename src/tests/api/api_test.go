package api_test

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway"
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
	serverURL     string
	httpClient    *http.Client
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
	testData := &TestData{
		serverURL:   *serverAddres,
		httpClient:  &http.Client{},
		ehrSystemID: common.EhrSystemID,
		//nolint
		users: []*User{
			&User{id: uuid.New().String(), password: fakeData.GetRandomStringWithLength(10)},
			&User{id: uuid.New().String(), password: fakeData.GetRandomStringWithLength(10)},
			&User{id: uuid.New().String(), password: fakeData.GetRandomStringWithLength(10)},
		},
	}

	for _, user := range testData.users {
		reqID, err := registerPatient(user, testData.ehrSystemID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatalf("Can not register user, err: %v", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			t.Fatal("registerPatient requestWait error: ", err)
		}
	}

	// TODO user register incorrect input data
	// TODO user register duplicate registration request
	//if !t.Run("User register", userRegister(testData)) {
	//	t.Fatal()
	//}

	if !t.Run("User login", userLogin(testData)) {
		t.Fatal()
	}

	if !t.Run("User register with doctor role", doctorRegister(testData)) {
		t.Fatal()
	}

	if !t.Run("User get info doctor", userInfoDoctor(testData)) {
		t.Fatal()
	}

	if !t.Run("User get info by code", userInfoByCode(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR creating", ehrCreate(testData)) {
		t.Fatal()
	}

	if !t.Run("Get transaction requests", requests(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR creating with id", ehrCreateWithID(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR creating with id for the same user", ehrCreateWithIDForSameUser(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR getting", ehrGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR get by subject", ehrGetBySubject(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS getting", ehrStatusGet(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS getting by version time", ehrStatusGetByVersionTime(testData)) {
		t.Fatal()
	}

	if !t.Run("EHR_STATUS update", ehrStatusUpdate(testData)) {
		t.Fatal()
	}

	if !t.Run("User get info patient", userInfoPatient(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION create Expected fail with wrong EhrId", compositionCreateFail(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION create Expected success with correct EhrId", compositionCreateSuccess(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION getting with correct EhrId", compositionGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION getting with wrong EhrId", compositionGetByWrongID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION update", compositionUpdate(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION get list", compositionGetList(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION delete by wrong UID", compositionDeleteByWrongID(testData)) {
		t.Fatal()
	}

	if !t.Run("COMPOSITION delete", compositionDeleteByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DIRECTORY create", directoryCRUD(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a invalid query", definitionStoreInvalidQuery(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a query", definitionStoreQuery(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a query version", definitionStoreQueryVersion(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Store a query version with same ID", definitionStoreQueryVersionWithSameID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Get stored query by ID", definitionStoredQueryGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION List stored queries", definitionListStoredQueries(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Template14 upload", definitionTemplate14Upload(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Template14 get by ID", definitionTemplate14GetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("DEFINITION Template14 list stored", definitionTemplate14List(testData)) {
		t.Fatal()
	}

	if !t.Run("QUERY execute with GET Expected success with correct query", queryExecSuccess(testData)) {
		t.Fatal()
	}

	if !t.Run("QUERY execute with POST Expected success with correct query", queryExecPostSuccess(testData)) {
		t.Fatal()
	}

	if !t.Run("QUERY execute with POST Expected fail with wrong query", queryExecPostFail(testData)) {
		t.Fatal()
	}

	if !t.Run("User group create", userGroupCreate(testData)) {
		t.Fatal()
	}

	if !t.Run("User group add user", userGroupAddUser(testData)) {
		t.Fatal()
	}

	if !t.Run("User group get by ID", userGroupGetByID(testData)) {
		t.Fatal()
	}

	if !t.Run("User group get list", userGroupGetList(testData)) {
		t.Fatal()
	}

	if !t.Run("User group remove user", userGroupRemoveUser(testData)) {
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

	apiHandler := gateway.New(cfg, infra)

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

func checkUser(testData *TestData) error {
	if len(testData.users) == 0 {
		user := &User{
			id:       uuid.New().String(),
			password: fakeData.GetRandomStringWithLength(10),
		}

		reqID, err := registerPatient(user, testData.ehrSystemID, testData.serverURL, testData.httpClient)
		if err != nil {
			return fmt.Errorf("Can not register user, err: %w", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			return fmt.Errorf("requestWait error, err: %w", err)
		}

		testData.users = append(testData.users, user)
	}

	return nil
}

func checkDoctor(testData *TestData) error {
	if len(testData.doctors) == 0 {
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
			return fmt.Errorf("Can not register user, err: %w", err)
		}

		err = requestWait(doctor.id, "", reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			return fmt.Errorf("registerPatient requestWait error: %w", err)
		}

		testData.doctors = append(testData.doctors, doctor)
	}

	return nil
}

func checkUserLogin(testData *TestData, user *User) error {
	if user.accessToken == "" {
		err := user.login(testData.ehrSystemID, testData.serverURL, testData.httpClient)
		if err != nil {
			return fmt.Errorf("User login error: %w", err)
		}
	}

	return nil
}

func checkUser0LoggedIn(testData *TestData) (*User, error) {
	err := checkUser(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser error: %w", err)
	}

	user := testData.users[0]

	err = checkUserLogin(testData, user)
	if err != nil {
		return nil, fmt.Errorf("checkUserLogin error: %w", err)
	}

	return user, nil
}

func checkDoctor0LoggedIn(testData *TestData) (*Doctor, error) {
	err := checkDoctor(testData)
	if err != nil {
		return nil, fmt.Errorf("checkDoctor error: %w", err)
	}

	doctor := testData.doctors[0]

	err = checkUserLogin(testData, &doctor.User)
	if err != nil {
		return nil, fmt.Errorf("checkDoctorLogin error: %w", err)
	}

	return doctor, nil
}

func checkUser0LoggedInAndEhrCreated(testData *TestData) (*User, error) {
	err := checkUser(testData)
	if err != nil {
		return nil, fmt.Errorf("checkUser error: %w", err)
	}

	user := testData.users[0]

	err = checkUserLogin(testData, user)
	if err != nil {
		return nil, fmt.Errorf("checkUserLogin error: %w", err)
	}

	err = checkEhr(testData, user)
	if err != nil {
		return nil, fmt.Errorf("checkEhr error: %w", err)
	}

	return user, nil
}

func checkEhr(testData *TestData, user *User) error {
	if user.ehrID == "" {
		ehr, reqID, err := createEhr(user.id, testData.ehrSystemID, user.accessToken, testData.serverURL, "", testData.httpClient)
		if err != nil {
			return fmt.Errorf("createEhr error: %w", err)
		}

		err = requestWait(user.id, user.accessToken, reqID, testData.serverURL, testData.httpClient)
		if err != nil {
			return fmt.Errorf("requestWait error: %w", err)
		}

		user.ehrID = ehr.EhrID.Value
		user.ehrStatusID = ehr.EhrStatus.ID.Value
	}

	return nil
}
