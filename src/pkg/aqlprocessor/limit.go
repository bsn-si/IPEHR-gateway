package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
	"strconv"
)

type Limit struct {
	Limit  int
	Offset int
}

func getLimit(ctx *aqlparser.LimitClauseContext) (*Limit, error) {
	limit := Limit{}

	if limitToken := ctx.GetLimit(); limitToken != nil && limitToken.GetTokenType() == aqlparser.AqlLexerINTEGER {
		limitVal, err := strconv.Atoi(limitToken.GetText())
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert limit value from string to int")
		}

		limit.Limit = limitVal
	}

	if offset := ctx.GetOffset(); offset != nil && offset.GetTokenType() == aqlparser.AqlLexerINTEGER {
		offset, err := strconv.Atoi(offset.GetText())
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert offset value from string to int")
		}

		limit.Offset = offset
	}

	return &limit, nil
}
