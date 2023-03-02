package driver

import (
	"database/sql"
	"encoding/json"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/treeindex"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestService_ExecuteQuery(t *testing.T) {
	dateVal, _ := time.Parse("2006-01-02", "1984-01-01")
	timeVal, _ := time.Parse("15:04:05.999", "15:35:10.123")
	dateTimeVal, _ := time.Parse("2006-01-02T15:04:05.999", "1984-01-01T15:35:10.123")

	type testDataStruct struct {
		Int      int
		Float    float64
		Str      string
		Date     time.Time
		Time     time.Time
		DateTime time.Time
	}

	tests := []struct {
		name      string
		query     string
		args      []interface{}
		dataFiles []string
		scan      func(rows *sqlx.Rows) (interface{}, error)
		want      interface{}
		wantErr   bool
	}{
		{
			"1. invalid AQL query",
			"invaid query",
			[]interface{}{},
			nil,
			func(rows *sqlx.Rows) (interface{}, error) {
				return nil, nil
			},
			nil,
			true,
		},
		{
			"2. select primitives",
			`SELECT 123, 1.23, 'hello world',
				'1984-01-01',
				'15:35:10.123',
				'1984-01-01T15:35:10.123'
			FROM EHR e`,
			[]interface{}{},
			[]string{"test_fixtures/composition_1.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				values := []testDataStruct{}

				for rows.Next() {
					var val testDataStruct
					err := rows.Scan(
						&val.Int,
						&val.Float,
						&val.Str,
						&val.Date,
						&val.Time,
						&val.DateTime,
					)
					if err != nil {
						return nil, errors.Wrap(err, "cannot scan test struct")
					}

					values = append(values, val)
				}

				return values, nil
			},
			[]testDataStruct{{123, 1.23, "hello world", dateVal, timeVal, dateTimeVal}},
			false,
		},
		{
			"2.1. select values",
			`SELECT 123
			FROM EHR e`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []int{}
				for rows.Next() {
					var val int
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}
					result = append(result, val)
				}

				sort.Ints(result)
				return result, nil
			},
			[]int{123},
			false,
		},
		{
			"3. select values",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []*float64{}
				for rows.Next() {
					var val any
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					if val != nil {
						f := val.(float64)
						result = append(result, toRef(f))
					} else {
						result = append(result, nil)
					}
				}

				sort.Slice(result, func(i, j int) bool {
					if result[i] != nil && result[j] != nil {
						return *result[i] < *result[j]
					}

					if result[i] == nil && result[j] == nil {
						return false
					}

					if result[i] == nil {
						return true
					}

					return false
				})
				return result, nil
			},
			[]*float64{nil, nil, nil, nil, nil, toRef(79.9), toRef(940.0), toRef(981.13)},
			false,
		},
		{
			"4. select values with WHERE",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []float64{}
				for rows.Next() {
					var val float64
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}
					result = append(result, val)
				}

				sort.Float64s(result)
				return result, nil
			},
			[]float64{940.0, 981.13},
			false,
		},
		{
			"5. select values with WHERE AND",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100
				AND o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude <= 940.0`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []float64{}
				for rows.Next() {
					var val float64
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}
					result = append(result, val)
				}

				sort.Float64s(result)
				return result, nil
			},
			[]float64{940.0},
			false,
		},
		{
			"6. select values with WHERE OR",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude <= 100
				OR o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude > 940.0`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []float64{}
				for rows.Next() {
					var val float64
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}
					result = append(result, val)
				}

				sort.Float64s(result)
				return result, nil
			},
			[]float64{79.9, 981.13},
			false,
		},
		{
			"7. select values with WHERE NOT",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				NOT o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude <= 100`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []*float64{}
				for rows.Next() {
					var val any
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					if val != nil {
						f := val.(float64)
						result = append(result, toRef(f))
					} else {
						result = append(result, nil)
					}
				}

				sort.Slice(result, func(i, j int) bool {
					if result[i] != nil && result[j] != nil {
						return *result[i] < *result[j]
					}

					if result[i] == nil && result[j] == nil {
						return false
					}

					if result[i] == nil {
						return true
					}

					return false
				})
				return result, nil
			},
			[]*float64{nil, nil, nil, nil, nil, toRef(940.0), toRef(981.13)},
			false,
		},
		{
			"8. select values with WHERE AND (OR)",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude <= 100
				OR
				(
					NOT	o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude = 940.0
				)`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := []*float64{}
				for rows.Next() {
					var val any
					if err := rows.Scan(&val); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					if val != nil {
						f := val.(float64)
						result = append(result, toRef(f))
					} else {
						result = append(result, nil)
					}
				}

				sort.Slice(result, func(i, j int) bool {
					if result[i] != nil && result[j] != nil {
						return *result[i] < *result[j]
					}

					if result[i] == nil && result[j] == nil {
						return false
					}

					if result[i] == nil {
						return true
					}

					return false
				})
				return result, nil
			},
			[]*float64{nil, nil, nil, nil, nil, toRef(79.9), toRef(981.13)},
			false,
		},
		{
			"9. select values and int value",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude,
			   10
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := [][]any{}
				for rows.Next() {
					var val float64
					var val2 int
					if err := rows.Scan(&val, &val2); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					result = append(result, []any{val, val2})
				}

				sort.Slice(result, func(i, j int) bool {
					return result[i][0].(float64) <= result[j][0].(float64)
				})
				return result, nil
			},
			[][]any{{940.0, 10}, {981.13, 10}},
			false,
		},
		{
			"10. select multipal columns",
			`SELECT
			   e/ehr_id/value AS ID,
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude,
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units
			FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := [][]any{}
				for rows.Next() {
					var id *string
					var val float64
					var val2 *string
					if err := rows.Scan(&id, &val, &val2); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					result = append(result, []any{id, val, val2})
				}

				sort.Slice(result, func(i, j int) bool {
					return result[i][1].(float64) <= result[j][1].(float64)
				})
				return result, nil
			},
			[][]any{{toRef("7d44b88c-4199-4bad-97dc-d78268e01398"), 940.0, toRef("/min")}, {toRef("7d44b88c-4199-4bad-97dc-d78268e01398"), 981.13, toRef("kg")}},
			false,
		},
		{
			"11. select with filter by EHR id",
			`SELECT
			   e/ehr_id/value AS ID,
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude,
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units
			FROM EHR e [ehr_id/value=$ehrUid]
				CONTAINS COMPOSITION c
					CONTAINS OBSERVATION o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			[]interface{}{
				sql.Named("ehrUid", "7d44b88c-4199-4bad-97dc-d78268e01398"),
			},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := [][]any{}
				for rows.Next() {
					var id *string
					var val float64
					var val2 *string
					if err := rows.Scan(&id, &val, &val2); err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					result = append(result, []any{id, val, val2})
				}

				sort.Slice(result, func(i, j int) bool {
					return result[i][1].(float64) <= result[j][1].(float64)
				})
				return result, nil
			},
			[][]any{
				{toRef("7d44b88c-4199-4bad-97dc-d78268e01398"), 940.0, toRef("/min")},
				{toRef("7d44b88c-4199-4bad-97dc-d78268e01398"), 981.13, toRef("kg")},
			},
			false,
		},
		{
			"12. select with filter by EHR id and observations version",
			`SELECT
			   e/ehr_id/value AS ID,
			   o/archetype_node_id,
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude,
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units
			FROM EHR e [ehr_id/value=$ehrUid]
				CONTAINS COMPOSITION c
					CONTAINS OBSERVATION o [openEHR-EHR-OBSERVATION.pulse.v2]
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			[]interface{}{
				sql.Named("ehrUid", "7d44b88c-4199-4bad-97dc-d78268e01398"),
			},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sqlx.Rows) (interface{}, error) {
				result := [][]any{}
				for rows.Next() {
					row, err := rows.SliceScan()
					if err != nil {
						return nil, errors.Wrap(err, "cannot scan float64 value")
					}

					r := []any{
						toRef(row[0].(string)),
						toRef(row[1].(string)),
						row[2].(float64),
						row[3].(string),
					}

					result = append(result, r)
				}

				sort.Slice(result, func(i, j int) bool {
					return result[i][2].(float64) <= result[j][2].(float64)
				})
				return result, nil
			},
			[][]any{
				{toRef("7d44b88c-4199-4bad-97dc-d78268e01398"), toRef("openEHR-EHR-OBSERVATION.pulse.v2"), 940.0, "/min"},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := getPreparedTreeIndex(tt.dataFiles...)
			if err != nil {
				t.Errorf("Service.ExecQuery() error on prepare tree index = %v", err)
				return
			}

			conn, err := sqlx.Open("aql", "")
			if err != nil {
				t.Fatal(err)
			}

			defer conn.Close()

			rows, err := conn.Queryx(tt.query, tt.args...)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("Service.ExecQuery() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			got, err := tt.scan(rows)
			if assert.Nil(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

const ehrFile = "./test_fixtures/ehr.json"

func getPreparedTreeIndex(filenames ...string) error {
	treeindex.DefaultEHRIndex = treeindex.NewEHRIndex()

	ehr := model.EHR{}
	{
		data, err := os.ReadFile(ehrFile)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &ehr); err != nil {
			return err
		}

		if err := treeindex.AddEHR(&ehr); err != nil {
			return errors.Wrap(err, "cannot add ehr into index")
		}
	}

	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			return errors.Wrap(err, "cannot read file")
		}

		comp := model.Composition{}
		if err := json.Unmarshal(data, &comp); err != nil {
			return errors.Wrap(err, "cannot unmarshal composition")
		}

		if err := treeindex.AddComposition(ehr.EhrID.Value, &comp); err != nil {
			return errors.Wrap(err, "cannot add Composition into EHRIndex")
		}
	}

	return nil
}

func toRef[T any](val T) *T {
	return &val
}
