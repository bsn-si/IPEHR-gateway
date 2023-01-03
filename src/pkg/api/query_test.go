package api

import (
	"fmt"
	"hms/gateway/pkg/api/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueryHandler_ExecStoredQuery(t *testing.T) {
	var (
		queryName = "some_aql_query"
		userID    = uuid.New().String()
	)

	tests := []struct {
		name        string
		queryParams string
		prepare     func(svc *mocks.MockQueryService)
		wantStatus  int
		want        string
	}{
		{
			"1. empty params",
			"ehr_id=7d44b88c-4199-4bad-97dc-d78268e01398&offset=1&fetch=10&some_val1=1&some_vak2=afasd",
			func(svc *mocks.MockQueryService) {},
			400,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := mocks.NewMockUserService(ctrl)
			querySvc := mocks.NewMockQueryService(ctrl)

			userSvc.EXPECT().VerifyAccess(userID, "Bearer AccessKey").Return(nil)

			tt.prepare(querySvc)

			api := API{
				User:  NewUserHandler(userSvc),
				Query: NewQueryHandler(querySvc, "base_url"),
			}

			router := api.setupRouter(api.buildQueryAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/query/%s?%s", queryName, tt.queryParams), nil)
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.want, string(respBody))
		})
	}
}
