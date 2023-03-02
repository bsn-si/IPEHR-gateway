package processor

import (
	"fmt"
	"io"
	"strconv"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Limit struct {
	Limit  int
	Offset int
}

func (l *Limit) write(w io.Writer) {
	fmt.Fprintf(w, "LIMIT %d", l.Limit)

	if l.Offset != 0 {
		fmt.Fprintf(w, " OFFSET %d", l.Offset)
	}
}

func getLimit(ctx *parser.LimitClauseContext) (*Limit, error) {
	limit := Limit{}

	if limitToken := ctx.GetLimit(); limitToken != nil && limitToken.GetTokenType() == parser.AqlLexerINTEGER {
		limitVal, err := strconv.Atoi(limitToken.GetText())
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert limit value from string to int")
		}

		limit.Limit = limitVal
		if limit.Limit == 0 {
			return nil, errors.New("LIMIT rows_count should by more then 0")
		}
	}

	if offset := ctx.GetOffset(); offset != nil && offset.GetTokenType() == parser.AqlLexerINTEGER {
		offset, err := strconv.Atoi(offset.GetText())
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert offset value from string to int")
		}

		limit.Offset = offset
	}

	return &limit, nil
}
