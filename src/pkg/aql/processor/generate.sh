#!/bin/sh

antlr4 -Dlanguage=Go -no-visitor -o ./../parser -package parser *.g4