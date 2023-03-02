package queryservice

import (
	"context"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryExecuterService(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.TODO(), 0)
	defer cancel()

	type args struct {
		ctx    context.Context
		query  string
		offset int
		limit  int
		params map[string]any
	}

	tests := []struct {
		name     string
		args     args
		prepare  func(mock sqlmock.Sqlmock)
		wantCol  []string
		wantRows []any
		wantErr  bool
	}{
		{
			"1. succces get data",
			args{
				context.Background(),
				"SELECT 123 as Number FROM e",
				0,
				0,
				nil,
			},
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"Number"}).
					AddRow(123)

				mock.ExpectQuery("SELECT 123 as Number FROM e").
					WillReturnRows(rows)
			},
			[]string{"Number"},
			[]any{[]any{int64(123)}},
			false,
		},
		{
			"2. error on run query",
			args{
				context.Background(),
				"SELECT 123 as Number FROM e",
				0,
				0,
				nil,
			},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 123 as Number FROM e").
					WillReturnError(errors.New("some error"))
			},
			nil,
			nil,
			true,
		},
		{
			"3. timeout error",
			args{
				ctx,
				"SELECT 123 as Number FROM e",
				0,
				0,
				nil,
			},
			func(mock sqlmock.Sqlmock) {
			},
			nil,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			tt.prepare(mock)

			svc := NewQueryService(sqlxDB)
			defer svc.Close()

			gotCol, gotRows, err := svc.runQuery(tt.args.ctx, tt.args.query, tt.args.offset, tt.args.limit, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("want error %v, got %v", tt.wantErr, err)
			}

			assert.Equal(t, tt.wantCol, gotCol)
			assert.Equal(t, tt.wantRows, gotRows)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
