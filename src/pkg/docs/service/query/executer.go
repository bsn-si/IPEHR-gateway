package query

import (
	"context"
	"database/sql"
	"hms/gateway/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type queryCallback func(columns []string, rows []any, err error)
type query struct {
	ctx      context.Context
	query    string
	args     []any
	callback queryCallback
}

type ExecuterService struct {
	db          *sqlx.DB
	queriesChan chan query
}

func NewQueryExecuterService(db *sqlx.DB) *ExecuterService {
	svc := &ExecuterService{
		db:          db,
		queriesChan: make(chan query, 100),
	}

	go svc.run()

	return svc
}

func (svc *ExecuterService) run() {
	for q := range svc.queriesChan {
		if q.ctx.Err() != nil {
			continue
		}

		rows, err := svc.db.QueryxContext(q.ctx, q.query, q.args...)
		if err != nil {
			q.callback(nil, nil, errors.Wrap(err, "cannot query rows"))
			continue
		}

		columns, err := rows.Columns()
		if err != nil {
			q.callback(nil, nil, errors.Wrap(err, "cannot get columns"))
			continue
		}

		result := []any{}

		for rows.Next() {
			row, err := rows.SliceScan()
			if err != nil {
				q.callback(nil, nil, errors.Wrap(err, "cannot scan row"))
			}

			result = append(result, row)
		}

		q.callback(columns, result, nil)
	}
}

func (svc *ExecuterService) Close() {
	close(svc.queriesChan)
}

func (svc *ExecuterService) ExecQueryContext(ctx context.Context, queryStr string, offset, limit int, params map[string]any) ([]string, []any, error) {
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

	var (
		resultColumns []string
		resultRows    []any
		resultErr     error
	)

	done := make(chan bool)
	q := query{
		ctx:   ctx,
		query: queryStr,
		args:  args,
		callback: func(columns []string, rows []any, err error) {
			resultColumns = columns
			resultRows = rows
			resultErr = err
			done <- true
		},
	}

	svc.queriesChan <- q
	select {
	case <-ctx.Done():
		return nil, nil, errors.New("timeout")
	case <-done:
		if resultErr != nil {
			return nil, nil, resultErr
		}
	}

	return resultColumns, resultRows, nil
}
