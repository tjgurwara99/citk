package golang

import (
	"context"
	"fmt"
	"regexp"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type Ident struct {
	Name    string
	Line    uint32
	EndLine uint32
	Col     uint32
	EndCol  uint32
}

var pascalCaseRegex = regexp.MustCompile("^[A-Z][a-z]+(?:[A-Z][a-z]+)*$")
var camelCaseRegex = regexp.MustCompile("^[a-z]+(?:[A-Z][a-z]+)*$")

func checkCase(ident string) bool {
	return !pascalCaseRegex.MatchString(ident) && !camelCaseRegex.MatchString(ident)
}

func AnomalousFuncSignatures(src []byte) ([]Ident, error) {
	filterFuncDecls := `(
		(function_declaration (identifier) @func)
	)`
	return anomalousDecls(src, filterFuncDecls, checkCase)
}

func AnomalousConstDecls(src []byte) ([]Ident, error) {
	filterConstDecls := `(
		const_spec (identifier) @constant
	)`
	return anomalousDecls(src, filterConstDecls, checkCase)
}

func AnomalousVarDecls(src []byte) ([]Ident, error) {
	filterVarDecls := `(
		var_spec (identifier) @constant
	)`

	return anomalousDecls(src, filterVarDecls, checkCase)
}

func AnomalousMethodAndFieldDecls(src []byte) ([]Ident, error) {
	// method_declaration (field_identifier) @methods
	filterMethodDecls := `(
		((field_identifier) @field)
	)`

	return anomalousDecls(src, filterMethodDecls, checkCase)
}

func anomalousDecls(src []byte, query string, condition func(ident string) bool) ([]Ident, error) {
	var decls []Ident

	// Parse source code
	lang := golang.GetLanguage()
	n, err := sitter.ParseCtx(context.Background(), src, lang)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}
	// Execute the query
	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		return nil, fmt.Errorf("failed to create a query for lang: %w", err)
	}
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)
	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, src)
		for _, c := range m.Captures {
			if ident := c.Node.Content(src); condition(ident) {
				decls = append(decls, Ident{
					Name:    c.Node.Content(src),
					Line:    c.Node.StartPoint().Row,
					EndLine: c.Node.EndPoint().Row,
					Col:     c.Node.StartPoint().Column,
					EndCol:  c.Node.EndPoint().Column,
				})
			}
		}
	}
	return decls, nil
}
