package aqlprocessor

import (
	"fmt"
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
	}
}

func (p *AqlProcessor) Process() error {
	antlr.ParseTreeWalkerDefault.Walk(p.listener, p.parser.SelectQuery())

	for {
		t := p.lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}

		fmt.Printf("%s (%q)\n", p.lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}

	return nil
}
