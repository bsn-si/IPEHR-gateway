package aqlquerier

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"testing"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"

	"github.com/google/go-cmp/cmp"
)

func TestService_ExecuteQuery(t *testing.T) {
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
			"2. simple composition and query",
			"SELECT c/value as val FROM Composition c",
			[]interface{}{},
			[]string{"test_fixtures/composition_1.json"},
			func(rows *sql.Rows) (interface{}, error) {
				values := [][]int{}

				for rows.Next() {
					var val int
					if err := rows.Scan(&val); err != nil {
						return nil, err
					}

					values = append(values, []int{val})
				}

				return values, nil
			},
			[][]int{{123}, {256}},
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

func getPreparedTreeIndex(filenames ...string) error {
	treeindex.DefaultTree = treeindex.NewTree()

	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
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
