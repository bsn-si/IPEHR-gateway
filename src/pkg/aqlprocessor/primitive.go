package aqlprocessor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type Primitive struct {
	Val any
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
	val, err := strconv.ParseFloat(numeric.GetText(), 64)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal numeric value")
	}

	p.Val = val
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
