package query

import (
	"context"
	"database/sql"
	"hms/gateway/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type ExecuterService struct {
	db *sqlx.DB
}

func NewQueryExecuterService(db *sqlx.DB) *ExecuterService {
	return &ExecuterService{
		db: db,
	}
}

func (svc *ExecuterService) ExecQueryContext(ctx context.Context, query string, offset, limit int, params map[string]any) ([]string, []any, error) {
	args := []any{}

	for k, v := range params {
		args = append(args, sql.Named(k, v))
	}

	if offset != 0 {
		args = append(args, sql.Named("offset", offset))
	}

	if limit != 0 {
		args = append(args, sql.Named("limit", limit))
	}

	rows, err := svc.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot query rows")
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get columns")
	}

	result := []any{}

	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot scan row")
		}

		result = append(result, row)
	}

	return columns, result, nil
}
