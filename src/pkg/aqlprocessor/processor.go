package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

//go:generate ./generate.sh

type AqlProcessor struct {
	lexer  *aqlparser.AqlLexer
	parser *aqlparser.AqlParser

	listener *AQLListener
}

func NewAqlProcessor(data string) *AqlProcessor {
	lexer := aqlparser.NewAqlLexer(antlr.NewInputStream(data))

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := aqlparser.NewAqlParser(stream)

	return &AqlProcessor{
		listener: NewAQLListener(),
		parser:   parser,
		lexer:    lexer,
	}
}

func (p *AqlProcessor) Process() (Query, error) {
	p.lexer.RemoveErrorListeners()
	p.parser.RemoveErrorListeners()

	lexerErrors := &CustomErrorListener{}
	p.lexer.AddErrorListener(lexerErrors)

	parserErrors := &CustomErrorListener{}
	p.parser.AddErrorListener(parserErrors)

	antlr.ParseTreeWalkerDefault.Walk(p.listener, p.parser.SelectQuery())

	var err error

	if len(lexerErrors.Errors) > 0 {
		for _, e := range lexerErrors.Errors {
			err = e
		}
	}

	if len(parserErrors.Errors) > 0 {
		for _, e := range parserErrors.Errors {
			err = e
		}
	}

	return p.listener.query, err
}

type CustomSyntaxError struct {
	line, column int
	msg          string
}

func (e *CustomSyntaxError) Error() string {
	return e.msg
}

type CustomErrorListener struct {
	*antlr.DefaultErrorListener // Embed default which ensures we fit the interface
	Errors                      []error
}

func (c *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	c.Errors = append(c.Errors, &CustomSyntaxError{
		line:   line,
		column: column,
		msg:    msg,
	})
}

func handleError(aParser antlr.Parser, token antlr.Token, err error) {
	aParser.GetErrorListenerDispatch().SyntaxError(
		nil, // recognizer antlr.Recognizer,
		nil, // offendingSymbol interface{},
		token.GetLine(),
		token.GetColumn(),
		err.Error(),
		nil, // e antlr.RecognitionException,
	)
}