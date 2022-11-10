#!/bin/sh

antlr4 -Dlanguage=Go -no-visitor -o ./aqlparser -package aqlparser *.g4