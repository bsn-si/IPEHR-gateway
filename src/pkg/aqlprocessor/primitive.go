package aqlprocessor

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor/aqlparser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"golang.org/x/exp/constraints"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/vmihailenco/msgpack/v5"
)

type Primitive struct {
	Val any
}

type PrimitiveType = uint8

const (
	PrimitiveTypeNone PrimitiveType = iota
	PrimitiveTypeInt
	PrimitiveTypeFloat64
	PrimitiveTypeString
	PrimitiveTypeBigFloat
	PrimitiveTypeBigInt
)

type PrimitiveWrap struct {
	Type PrimitiveType
	Val  any
}

func (p Primitive) EncodeMsgpack(enc *msgpack.Encoder) error {
	switch v := p.Val.(type) {
	case int8:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeInt, int(v)})
	case int16:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeInt, int(v)})
	case uint16:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeInt, int(v)})
	case int32:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeInt, int(v)})
	case uint32:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeInt, int(v)})
	case int:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeInt, v})
	case float64:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeFloat64, v})
	case *big.Int:
		fmt.Println("p=", v)
		return enc.Encode(PrimitiveWrap{PrimitiveTypeBigInt, v})
	case *big.Float:
		return enc.Encode(PrimitiveWrap{PrimitiveTypeBigFloat, v})
	default:
		return errors.Errorf("Unsupported Primitive.Val type: %T", v)
	}
}

func (p *Primitive) UnmarshalMsgpack(data []byte) error {
	var tmp map[string]any

	err := msgpack.Unmarshal(data, &tmp)
	if err != nil {
		panic(err)
	}

	_type, ok := tmp["Type"]
	if !ok {
		return errors.ErrFieldIsEmpty("Type")
	}

	value, ok := tmp["Val"]
	if !ok {
		return errors.ErrFieldIsEmpty("Val")
	}

	switch _type {
	case PrimitiveTypeInt, PrimitiveTypeFloat64, PrimitiveTypeString:
		switch v := value.(type) {
		case int8:
			p.Val = int(v)
		case float64:
			p.Val = v
		case string:
			p.Val = v
		}
	case PrimitiveTypeBigInt:
		switch v := value.(type) {
		case []uint8:
			p.Val = new(big.Int).SetBytes(v)
			fmt.Println("p.Val=", p.Val)
		}
	case PrimitiveTypeBigFloat:
		switch v := value.(type) {
		case []uint8:
			p.Val = big.NewFloat(0)

			err = p.Val.(*big.Float).UnmarshalText(v)
			if err != nil {
				return fmt.Errorf("big.Float UnmarshalText error: %w", err)
			}
		}
	default:
		return errors.Errorf("Unsupported Primitive type: %d", _type)
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
				return false, errors.Errorf("Unsupported comparison v=%v (%T) %s p=%v (%T)", v, v, cmpSymbl, p, p)
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
				return false, errors.Errorf("Unsupported comparison v=%v (%T) %s p=%v (%d)", v, v, cmpSymbl, p, p)
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
			case *big.Int:
				vBigFloat := new(big.Float).SetInt(v)
				return compareBigFloat(vBigFloat, p, cmpSymbl), nil
			case *big.Float:
				return compareBigFloat(v, p, cmpSymbl), nil
			default:
				return false, errors.Errorf("Unsupported comparison v=%v (%T) %s p=%v (%T)", v, v, cmpSymbl, p, p)
			}
		}
	case *big.Int:
		{
			switch v := v.(type) {
			case int:
				return compareBigInt(big.NewInt(int64(v)), p, cmpSymbl), nil
			case int64:
				return compareBigInt(big.NewInt(v), p, cmpSymbl), nil
			case float64:
				vBig := new(big.Float).SetFloat64(v)
				pBigFloat := new(big.Float).SetInt(p)
				return compareBigFloat(vBig, pBigFloat, cmpSymbl), nil
			case *big.Int:
				return compareBigInt(v, p, cmpSymbl), nil
			case *big.Float:
				pBigFloat := new(big.Float).SetInt(p)
				return compareBigFloat(v, pBigFloat, cmpSymbl), nil
			default:
				return false, errors.Errorf("Unsupported comparison v=%v (%T) %s p=%v (%T)", v, v, cmpSymbl, p, p)
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
		return false, errors.Errorf("Unsupported comparison p=%v (%T)", p, p)
	}

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
		fmt.Println("Unknown BigInt comparison symbol: ", cmpSymbl)
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
		fmt.Println("Unknown BigFloat comparison symbol: ", cmpSymbl)
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
		fmt.Println("Unknown comparison symbol: ", cmpSymbl)
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
