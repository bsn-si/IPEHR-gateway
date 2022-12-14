package aqlquerier

import (
	"database/sql"
	"encoding/json"
	"os"
	"sort"
	"testing"
	"time"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"

	"github.com/google/go-cmp/cmp"
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
		scan      func(rows *sql.Rows) (interface{}, error)
		want      interface{}
		wantErr   bool
	}{
		{
			"1. invalid AQL query",
			"invaid query",
			[]interface{}{},
			nil,
			func(rows *sql.Rows) (interface{}, error) {
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
			FROM Observation o`,
			[]interface{}{},
			[]string{"test_fixtures/composition_1.json"},
			func(rows *sql.Rows) (interface{}, error) {
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
			"3. select values",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
   			FROM Observation o`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sql.Rows) (interface{}, error) {
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
			[]float64{79.9, 940.0, 981.13},
			false,
		},
		{
			"4. select values with WHERE",
			`SELECT
			   o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude
   			FROM Observation o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sql.Rows) (interface{}, error) {
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
   			FROM Observation o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude >= 100
				AND o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude <= 940.0`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sql.Rows) (interface{}, error) {
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
   			FROM Observation o
			WHERE
				o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude <= 100
				OR o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude > 940.0`,
			[]interface{}{},
			[]string{"test_fixtures/composition_2.json"},
			func(rows *sql.Rows) (interface{}, error) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := getPreparedTreeIndex(tt.dataFiles...)
			if err != nil {
				t.Errorf("Service.ExecQuery() error on prepare tree index = %v", err)
				return
			}

			conn, err := sql.Open("aql", "")
			if err != nil {
				t.Fatal(err)
			}

			defer conn.Close()

			rows, err := conn.Query(tt.query, tt.args...)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("Service.ExecQuery() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			got, err := tt.scan(rows)
			if err != nil {
				t.Errorf("Service.ExecQuery() scan rows error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Service.ExecQuery() = mismatch {-want;+got}\n\t%s", diff)
			}
		})
	}
}

// func Test_GSJON(t *testing.T) {
// 	data, err := ioutil.ReadFile("test_fixtures/composition_2.json")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	observations := gjson.GetBytes(data,
// 		`content.#(_type=="SECTION")#.items.#(_type=="OBSERVATION")#.data`,
// 	)
// 	// o/data[at0001]/events[at0006]/data[at0003]/items[at0005]/value/magnitude

// 	for _, o := range observations.Array() {
// 		o.ForEach(func(key, value gjson.Result) bool {
// 			if value.Get("archetype_node_id").String() == "at0002" {
// 				r := value.Get(`events.#(archetype_node_id="at0003")#.data`)
// 				r.ForEach(func(key, value gjson.Result) bool {
// 					if value.Get("archetype_node_id").String() == "at0001" {
// 						val := value.Get(`items.#(archetype_node_id="at0004")#.value.magnitude`)
// 						log.Println(val)
// 					}
// 					return true
// 				})
// 			}
// 			return true
// 		})
// 	}
// 	// t.Fail()
// }

func getPreparedTreeIndex(filenames ...string) error {
	treeindex.DefaultTree = treeindex.NewTree()

	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			return errors.Wrap(err, "cannot read file")
		}

		comp := model.Composition{}
		if err := json.Unmarshal(data, &comp); err != nil {
			return errors.Wrap(err, "cannot unmarshal composition")
		}

		if err := treeindex.DefaultTree.AddComposition(comp); err != nil {
			return errors.Wrap(err, "cannot add Composition into TreeIndex")
		}
	}

	return nil
}
