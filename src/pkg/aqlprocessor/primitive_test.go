package aqlprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessor_SelectNull(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Query
		wantErr bool
	}{
		{
			"1. select NULL",
			`SELECT NULL FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{
						{Value: &PrimitiveSelectValue{Val: Primitive{Val: nil}}},
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			tt.want.From = got.From
			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProcessor_SelectString(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Query
		wantErr bool
	}{
		{
			"1. select string 2",
			`SELECT 'hello_world' FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{
						{Value: &PrimitiveSelectValue{Val: Primitive{Val: "hello_world"}}},
					},
				},
			},
			false,
		},
		{
			"2. select string 2",
			`SELECT "hello_world" FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{
						{Value: &PrimitiveSelectValue{Val: Primitive{Val: "hello_world"}}},
					},
				},
			},
			false,
		},
		{
			"3. select string 3",
			`SELECT '"hello_world"' FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{
						{Value: &PrimitiveSelectValue{Val: Primitive{Val: `"hello_world"`}}},
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			tt.want.From = got.From
			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProcessor_SelectBoolean(t *testing.T) {
	// todo: this test doesn't work due to problems in lexer
	// tests := []struct {
	// 	name    string
	// 	query   string
	// 	want    Select
	// 	wantErr bool
	// }{
	// 	{
	// 		"1. select TRUE",
	// 		`SELECT true AS dangerousBP FROM EHR`,
	// 		Select{
	// 			SelectExprs: []SelectExpr{
	// 				{Value: &PrimitiveSelectValue{Val: Primitive{Val: true}}},
	// 			},
	// 		},
	// 		false,
	// 	},
	// 	{
	// 		"2. select FALSE",
	// 		`SELECT FALSE FROM EHR`,
	// 		Select{
	// 			SelectExprs: []SelectExpr{
	// 				{Value: &PrimitiveSelectValue{Val: Primitive{Val: false}}},
	// 			},
	// 		},
	// 		false,
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		p := NewAqlProcessor(tt.query)
	// 		got, err := p.Process()
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
	// 		}

	// 		if !tt.wantErr && assert.NoError(t, err) {
	// 			assert.Equal(t, tt.want, got.Select)
	// 		}
	// 	})
	// }
}

func TestProcessor_SelectNumeric(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Query
		wantErr bool
	}{
		{
			"1. SELECT 0",
			`SELECT 0 FROM EHR`,
			Query{Select: Select{
				SelectExprs: []SelectExpr{
					{
						Value: &PrimitiveSelectValue{Val: Primitive{0.0}},
					},
				},
			}},
			false,
		},
		{
			"2. SELECT -1",
			`SELECT -1 FROM EHR`,
			Query{Select: Select{
				SelectExprs: []SelectExpr{
					{
						Value: &PrimitiveSelectValue{Val: Primitive{-1.0}},
					},
				},
			}},
			false,
		},

		{
			"3. SELECT 123.5e+10",
			`SELECT 123.5e+10 FROM EHR`,
			Query{Select: Select{
				SelectExprs: []SelectExpr{
					{
						Value: &PrimitiveSelectValue{Val: Primitive{123.5e+10}},
					},
				},
			}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			tt.want.From = got.From

			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProcessor_SelectDates(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2020-10-11")
	timeValue, _ := time.Parse("15:04:05.999", "23:58:58.123")
	dateTimeValue, _ := time.Parse("2006-01-02T15:04:05.999", "2020-10-11T23:58:58.123")

	tests := []struct {
		name    string
		query   string
		want    Query
		wantErr bool
	}{
		{
			"1. date 2020-10-11",
			`SELECT '2020-10-11' FROM EHR`,
			Query{Select: Select{
				SelectExprs: []SelectExpr{
					{
						Value: &PrimitiveSelectValue{Val: Primitive{date}},
					},
				},
			}},
			false,
		},
		{
			"2. time  23:58:58.123",
			`SELECT '23:58:58.123' FROM EHR`,
			Query{Select: Select{
				SelectExprs: []SelectExpr{
					{
						Value: &PrimitiveSelectValue{Val: Primitive{timeValue}},
					},
				},
			}},
			false,
		},
		{
			"3. date_time  2020-10-11 23:58:58.123",
			`SELECT '2020-10-11T23:58:58.123' FROM EHR`,
			Query{Select: Select{
				SelectExprs: []SelectExpr{
					{
						Value: &PrimitiveSelectValue{Val: Primitive{dateTimeValue}},
					},
				},
			}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			tt.want.From = got.From
			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
