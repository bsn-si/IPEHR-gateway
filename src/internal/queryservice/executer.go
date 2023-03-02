package queryservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type queryCallback func(columns []string, rows []any, err error)
type query struct {
	ctx      context.Context
	query    string
	args     []any
	callback queryCallback
}

type QueryService struct {
	db          *sqlx.DB
	queriesChan chan query
}

func NewQueryService(db *sqlx.DB) *QueryService {
	svc := &QueryService{
		db:          db,
		queriesChan: make(chan query, 100),
	}

	go svc.run()

	return svc
}

func (svc *QueryService) run() {
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

func (svc *QueryService) Close() {
	close(svc.queriesChan)
}

func (svc *QueryService) ExecQuery(ctx context.Context, query *model.QueryRequest) (*model.QueryResponse, error) {
	columns, result, err := svc.runQuery(ctx, query.Query, query.Offset, query.Fetch, query.QueryParameters)
	if err != nil {
		return nil, errors.Wrap(err, "cannot exec query")
	}

	resp := &model.QueryResponse{
		Query: query.Query,
		Rows:  result,
	}

	for _, c := range columns {
		resp.Columns = append(resp.Columns, model.QueryColumn{Name: c})
	}

	return resp, nil
}

func (svc *QueryService) runQuery(ctx context.Context, queryStr string, offset, limit int, params map[string]any) ([]string, []any, error) {
	args := []any{}

	for k, v := range params {
		args = append(args, sql.Named(k, v))
	}

	if offset != 0 {
		queryStr = fmt.Sprintf("%s OFFSET %d", queryStr, offset)
	}

	if limit != 0 {
		queryStr = fmt.Sprintf("%s LIMIT %d", queryStr, limit)
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
		return nil, nil, errors.ErrTimeout
	case <-done:
		if resultErr != nil {
			return nil, nil, resultErr
		}
	}

	return resultColumns, resultRows, nil
}
