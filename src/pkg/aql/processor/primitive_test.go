package processor

import (
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/google/go-cmp/cmp"
)

func TestProcessor_SelectNull(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Select
		wantErr bool
	}{
		{
			"1. select NULL",
			`SELECT NULL FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "NULL",
						Value: &PrimitiveSelectValue{Val: Primitive{Val: nil, Type: parser.AqlLexerNULL}},
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

			if got == nil {
				return
			}
			if diff := cmp.Diff(tt.want, got.Select); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}

func TestProcessor_SelectString(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Select
		wantErr bool
	}{
		{
			"1. select string 2",
			`SELECT 'hello_world' FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  `'hello_world'`,
						Value: &PrimitiveSelectValue{Val: Primitive{Val: "hello_world", Type: parser.AqlLexerSTRING}},
					},
				},
			},
			false,
		},
		{
			"2. select string 2",
			`SELECT "hello_world" FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  `"hello_world"`,
						Value: &PrimitiveSelectValue{Val: Primitive{Val: "hello_world", Type: parser.AqlLexerSTRING}},
					},
				},
			},
			false,
		},
		{
			"3. select string 3",
			`SELECT '"hello_world"' FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  `'"hello_world"'`,
						Value: &PrimitiveSelectValue{Val: Primitive{Val: `"hello_world"`, Type: parser.AqlLexerSTRING}},
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

			if got == nil {
				return
			}
			if diff := cmp.Diff(tt.want, got.Select); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
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
		want    Select
		wantErr bool
	}{
		{
			"1. SELECT 0",
			`SELECT 0 FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "0",
						Value: &PrimitiveSelectValue{Val: Primitive{0, parser.AqlLexerINTEGER}},
					},
				},
			},
			false,
		},
		{
			"2. SELECT -1",
			`SELECT -1 FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "-1",
						Value: &PrimitiveSelectValue{Val: Primitive{-1, parser.AqlLexerINTEGER}},
					},
				},
			},
			false,
		},
		{
			"3. SELECT 123.5e+10",
			`SELECT 123.5e+10 FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "123.5e+10",
						Value: &PrimitiveSelectValue{Val: Primitive{123.5e+10, parser.AqlLexerREAL}},
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
			if got == nil {
				return
			}
			if diff := cmp.Diff(tt.want, got.Select); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
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
		want    Select
		wantErr bool
	}{
		{
			"1. date 2020-10-11",
			`SELECT '2020-10-11' FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "'2020-10-11'",
						Value: &PrimitiveSelectValue{Val: Primitive{date, parser.AqlLexerDATE}},
					},
				},
			},
			false,
		},
		{
			"2. time  23:58:58.123",
			`SELECT '23:58:58.123' FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "'23:58:58.123'",
						Value: &PrimitiveSelectValue{Val: Primitive{timeValue, parser.AqlLexerTIME}},
					},
				},
			},
			false,
		},
		{
			"3. date_time  2020-10-11 23:58:58.123",
			`SELECT '2020-10-11T23:58:58.123' FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{
					{
						Path:  "'2020-10-11T23:58:58.123'",
						Value: &PrimitiveSelectValue{Val: Primitive{dateTimeValue, parser.AqlLexerDATETIME}},
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

			if got == nil {
				return
			}
			if diff := cmp.Diff(tt.want, got.Select); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}

func TestPrimitive_Compare(t *testing.T) {
	tests := []struct {
		name     string
		prim     Primitive
		val      any
		cmpSymbl ComparisionSymbol
		want     bool
	}{
		{"1. 123 == 123", Primitive{Val: 123}, 123, SymEQ, true},
		{"2. 123 == 321", Primitive{Val: 123}, 321, SymEQ, false},
		{"3. 123 != 321", Primitive{Val: 123}, 331, SymNe, true},
		{"4. 123 > 321", Primitive{Val: 123}, 321, SymGT, true},
		{"5. 123 >= 321", Primitive{Val: 123}, 331, SymGE, true},
		{"6. 123 < 100", Primitive{Val: 123}, 100, SymLT, true},
		{"7. 123 <= 100", Primitive{Val: 123}, 100, SymLE, true},

		{"8. 123.0 == 123", Primitive{Val: 123.0}, 123, SymEQ, true},
		{"9. 123.0 == 321", Primitive{Val: 123.0}, 321, SymEQ, false},
		{"10. 123.0 != 321", Primitive{Val: 123.0}, 331, SymNe, true},
		{"11. 123.0 > 321", Primitive{Val: 123.0}, 321, SymGT, true},
		{"12. 123.0 >= 321", Primitive{Val: 123.0}, 331, SymGE, true},
		{"13. 123.0 < 100", Primitive{Val: 123.0}, 100, SymLT, true},
		{"14. 123.0 <= 100", Primitive{Val: 123.0}, 100, SymLE, true},

		{"15. 123 <= 100.1", Primitive{Val: 123}, 100.1, SymLE, true},
		{"16. 123.1 <= 100.1", Primitive{Val: 123.1}, 100.1, SymLE, true},

		{`17. "aaa" != "bbb"`, Primitive{Val: "aaa"}, "bbb", SymNe, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prim.Compare(tt.val, tt.cmpSymbl); got != tt.want {
				t.Errorf("Primitive.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
