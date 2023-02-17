package aqlprocessor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"math/big"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor/aqlparser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"golang.org/x/exp/constraints"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/vmihailenco/msgpack/v5"
)

type PrimitiveType byte

const (
	PrimitiveTypeBigFloat PrimitiveType = iota
	PrimitiveTypeBigInt
	PrimitiveTypeString
	PrimitiveTypeInt 
	PrimitiveTypeFloat64
	PrimitiveTypeNull
)

type Primitive struct {
	Val any
}

func (p *Primitive) UnmarshalMsgpack(data []byte) error {
	tmp := struct {
		Val any
	}{}

	if err := msgpack.Unmarshal(data, &tmp); err != nil {
		return err
	}

	switch v := tmp.Val.(type) {
	case int8:
		p.Val = int(v)
	case int16:
		p.Val = int(v)
	case uint16:
		p.Val = int(v)
	case int32:
		p.Val = int(v)
	case uint32:
		p.Val = int(v)
	default:
		p.Val = tmp.Val
	}

	return nil
}

func (p Primitive) Compare(v any, cmpSymbl ComparisionSymbol) (bool, error) {
	switch p := p.Val.(type) {
	case int:
		{
			switch v := v.(type) {
			case int:
				return compare(v, p, cmpSymbl), nil
			case int64:
				return compare(v, int64(p), cmpSymbl), nil
			case float64:
				return compare(v, float64(p), cmpSymbl), nil
			default:
				return false, errors.Errorf("Unsupported comparison v=%d (%T) %s p=%v (%T)", v,v,cmpSymbl,p,p)
			}
		}
	case float64:
		{
			switch v := v.(type) {
			case int:
				return compare(float64(v), p, cmpSymbl), nil
			case int64:
				return compare(float64(v), p, cmpSymbl), nil
			case float64:
				return compare(v, p, cmpSymbl), nil
			default:
				return false, errors.Errorf("Unsupported comparison v=%d (%T) %s p=%v (%d)", v,v,cmpSymbl,p,p)
			}
		}
	case *big.Float:
		{
			switch v := v.(type) {
			case int:
				vBig := new(big.Float).SetInt64(int64(v))
				return compareBigFloat(vBig, p, cmpSymbl), nil
			case int64:
				vBig := new(big.Float).SetInt64(v)
				return compareBigFloat(vBig, p, cmpSymbl), nil
			case float64:
				vBig := new(big.Float).SetFloat64(v)
				return compareBigFloat(vBig, p, cmpSymbl), nil
			default:
				return false, errors.Errorf("Unsupported comparison v=%d (%T) %s p=%v (%T)", v,v,cmpSymbl,p,p)
			}
		}
	case *big.Int:
		{
			switch v := v.(type) {
			case int:
				vBig := big.NewInt(int64(v))
				return compareBigInt(vBig, p, cmpSymbl), nil
			case int64:
				vBig := big.NewInt(v)
				return compareBigInt(vBig, p, cmpSymbl), nil
			case float64:
				vBig := new(big.Float).SetFloat64(v)
				pBigFloat := new(big.Float).SetInt(p)
				return compareBigFloat(vBig, pBigFloat, cmpSymbl), nil
			default:
				return false, errors.Errorf("Unsupported comparison v=%d (%T) %s p=%v(%T)", v,v,cmpSymbl,p,p)
			}
		}
	 case string:
		 {
			switch v := v.(type) {
			case string:
				return compare(v, p, cmpSymbl), nil
			}
		}
	default:
		return false, errors.Errorf("Unsupported comparison p=%v (%T)", p,p)
	}

	/*
	x := reflect.ValueOf(p.Val)
	y := reflect.ValueOf(v)

	switch {
	case x.Type() == y.Type():
		switch x.Kind() {
		case reflect.Int:
			return compare(val.(int), p.Val.(int), cmpSymbl)
		case reflect.Float64:
			return compare(val.(float64), p.Val.(float64), cmpSymbl)
		case reflect.String:
			return compare(val.(string), p.Val.(string), cmpSymbl)
		}
	case x.Kind() == reflect.Float64 && y.Kind() == reflect.Int:
		return compare(float64(val.(int)), p.Val.(float64), cmpSymbl)
	case x.Kind() == reflect.Int && y.Kind() == reflect.Float64:
		return compare(val.(float64), float64(p.Val.(int)), cmpSymbl)
	default:
		fmt.Printf("x: %v y: %v\n", x.Kind(), y.Kind())
	}
	*/
	return false, errors.ErrIsUnsupported
}

func compareBigInt(x, y *big.Int, cmpSymbl ComparisionSymbol) bool {
	switch cmpSymbl {
	case SymLT:
		return x.Cmp(y) == -1
	case SymGT:
		return x.Cmp(y) == 1
	case SymLE:
		return x.Cmp(y) <= 0
	case SymGE:
		return x.Cmp(y) >= 0
	case SymNe:
		return x.Cmp(y) != 0
	case SymEQ:
		return x.Cmp(y) == 0
	default:
		return false
	}
}

func compareBigFloat(x, y *big.Float, cmpSymbl ComparisionSymbol) bool {
	switch cmpSymbl {
	case SymLT:
		return x.Cmp(y) == -1
	case SymGT:
		return x.Cmp(y) == 1
	case SymLE:
		return x.Cmp(y) <= 0
	case SymGE:
		return x.Cmp(y) >= 0
	case SymNe:
		return x.Cmp(y) != 0
	case SymEQ:
		return x.Cmp(y) == 0
	default:
		return false
	}
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

func getPrimitive(prim *aqlparser.PrimitiveContext) (Primitive, error) {
	p := Primitive{}

	switch val := prim.GetChild(0).(type) {
	case *antlr.TerminalNodeImpl:
		if err := p.processTerminalNode(val); err != nil {
			return Primitive{}, errors.Wrap(err, "cannot get Primitive.TerminalNode")
		}
	case *aqlparser.NumericPrimitiveContext:
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
	switch tokenType {
	case aqlparser.AqlLexerSTRING:
		p.Val = trimString(terminal.String())
	case aqlparser.AqlLexerDATE:
		{
			const layout = "2006-01-02"
			t, err := parseDateByLayout(layout, terminal.String())
			if err != nil {
				return err
			}

			p.Val = t
		}
	case aqlparser.AqlLexerTIME:
		{
			const layout = "15:04:05.999"
			t, err := parseDateByLayout(layout, terminal.String())
			if err != nil {
				return err
			}

			p.Val = t
		}
	case aqlparser.AqlLexerDATETIME:
		{
			const layout = "2006-01-02T15:04:05.999"
			t, err := parseDateByLayout(layout, terminal.String())
			if err != nil {
				return err
			}

			p.Val = t
		}
	case aqlparser.AqlLexerBOOLEAN:
		p.Val = strings.ToLower(terminal.String()) == "true"
	case aqlparser.AqlLexerNULL:
		p.Val = nil
	default:
		return fmt.Errorf("unexpected PRIMITIVE SYMBOL type: %v", tokenType) //nolint
	}

	return nil
}

func (p *Primitive) processNumericPrimitive(numeric *aqlparser.NumericPrimitiveContext) error {
	if numeric.INTEGER() != nil {
		val, err := strconv.Atoi(numeric.INTEGER().GetText())
		if err != nil {
			return errors.Wrap(err, "cannot unmarshal numeric value")
		}

		p.Val = val
	} else if numeric.SYM_MINUS() != nil {
		err := p.processNumericPrimitive(numeric.NumericPrimitive().(*aqlparser.NumericPrimitiveContext))
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
	} else {
		val, err := strconv.ParseFloat(numeric.GetText(), 64)
		if err != nil {
			return errors.Wrap(err, "cannot unmarshal numeric value")
		}

		p.Val = val
	}

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
