package processor

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"golang.org/x/exp/constraints"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type Primitive struct {
	Val  any
	Type int
}

func (p *Primitive) write(w io.Writer) {
	switch p.Type {
	case parser.AqlLexerINTEGER:
		fmt.Fprintf(w, "%d", p.Val)
	case parser.AqlLexerREAL:
		fmt.Fprintf(w, "%f", p.Val)
	case parser.AqlLexerSTRING:
		fmt.Fprintf(w, "'%s'", p.Val)
	case parser.AqlLexerDATE:
		t, _ := p.Val.(time.Time)
		fmt.Fprintf(w, "'%s'", t.Format("2006-01-02T15"))
	case parser.AqlLexerTIME:
		t, _ := p.Val.(time.Time)
		fmt.Fprintf(w, "'%s'", t.Format("15:04:05.999"))
	case parser.AqlLexerDATETIME:
		t, _ := p.Val.(time.Time)
		fmt.Fprintf(w, "'%s'", t.Format("2006-01-02T15:04:05.999"))
	default:
		fmt.Fprint(w, "NULL")
	}
}

func (p Primitive) Compare(val any, cmpSymbl ComparisionSymbol) bool {
	x := reflect.ValueOf(p.Val)
	y := reflect.ValueOf(val)

	if x.Type() == y.Type() {
		switch x.Kind() {
		case reflect.Int:
			return compare(val.(int), p.Val.(int), cmpSymbl)
		case reflect.Float64:
			return compare(val.(float64), p.Val.(float64), cmpSymbl)
		case reflect.String:
			return compare(val.(string), p.Val.(string), cmpSymbl)
		}
	} else if x.Kind() == reflect.Float64 && y.Kind() == reflect.Int {
		return compare(float64(val.(int)), p.Val.(float64), cmpSymbl)
	} else if x.Kind() == reflect.Int && y.Kind() == reflect.Float64 {
		return compare(val.(float64), float64(p.Val.(int)), cmpSymbl)
	}

	return false
}

func compare[T constraints.Ordered](x, y T, cmpSymbl ComparisionSymbol) bool {
	switch cmpSymbl {
	case SymLT:
		return x < y
	case SymGT:
		return x > y
	case SymLE:
		return x <= y
	case SymGE:
		return x >= y
	case SymNe:
		return x != y
	case SymEQ:
		return x == y
	default:
		return false
	}
}

func getPrimitive(prim *parser.PrimitiveContext) (Primitive, error) {
	p := Primitive{}

	switch val := prim.GetChild(0).(type) {
	case *antlr.TerminalNodeImpl:
		if err := p.processTerminalNode(val); err != nil {
			return Primitive{}, errors.Wrap(err, "cannot get Primitive.TerminalNode")
		}
	case *parser.NumericPrimitiveContext:
		if err := p.processNumericPrimitive(val); err != nil {
			return Primitive{}, errors.Wrap(err, "cannot get Primitive.Numeric")
		}
	default:
		return Primitive{}, fmt.Errorf("unexpected PRIMITIVE type: %T", val) //nolint
	}

	return p, nil
}

func (p *Primitive) processTerminalNode(terminal *antlr.TerminalNodeImpl) error {
	tokenType := terminal.GetSymbol().GetTokenType()
	p.Type = tokenType

	switch tokenType {
	case parser.AqlLexerSTRING:
		p.Val = trimString(terminal.String())
	case parser.AqlLexerDATE:
		{
			const layout = "2006-01-02"
			t, err := parseDateByLayout(layout, terminal.String())
			if err != nil {
				return err
			}

			p.Val = t
		}
	case parser.AqlLexerTIME:
		{
			const layout = "15:04:05.999"
			t, err := parseDateByLayout(layout, terminal.String())
			if err != nil {
				return err
			}

			p.Val = t
		}
	case parser.AqlLexerDATETIME:
		{
			const layout = "2006-01-02T15:04:05.999"
			t, err := parseDateByLayout(layout, terminal.String())
			if err != nil {
				return err
			}

			p.Val = t
		}
	case parser.AqlLexerBOOLEAN:
		p.Val = strings.ToLower(terminal.String()) == "true"
	case parser.AqlLexerNULL:
		p.Val = nil
	default:
		return fmt.Errorf("unexpected PRIMITIVE SYMBOL type: %v", tokenType) //nolint
	}

	return nil
}

func (p *Primitive) processNumericPrimitive(numeric *parser.NumericPrimitiveContext) error {
	if numeric.INTEGER() != nil {
		val, err := strconv.Atoi(numeric.INTEGER().GetText())
		if err != nil {
			return errors.Wrap(err, "cannot unmarshal numeric value")
		}

		p.Val = val
		p.Type = parser.AqlLexerINTEGER
		return nil
	}

	if numeric.SYM_MINUS() != nil {
		err := p.processNumericPrimitive(numeric.NumericPrimitive().(*parser.NumericPrimitiveContext))
		if err != nil {
			return err
		}

		switch val := p.Val.(type) {
		case int:
			p.Val = -val
		case float64:
			p.Val = -val
		default:
			return errors.New("unexpected primitive value type")
		}
		return nil
	}

	val, err := strconv.ParseFloat(numeric.GetText(), 64)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal numeric value")
	}

	p.Val = val
	p.Type = parser.AqlLexerREAL

	return nil
}

func parseDateByLayout(layout, str string) (time.Time, error) {
	strDate := trimString(str)

	date, err := time.Parse(layout, strDate)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "cannot parse date by layout")
	}

	return date, nil
}

func trimString(str string) string {
	if str[0] == '\'' {
		str = strings.Trim(str, "'")
	} else if str[0] == '"' {
		str = strings.Trim(str, "\"")
	}

	return str
}
