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

func AnomalousFuncSignatures(src []byte) ([]Ident, error) {
	var funcDecls []Ident
	filterFuncDecls := `(
		(function_declaration (identifier) @func)
	)`

	pascalCaseRegex := regexp.MustCompile("^[A-Z][a-z]+(?:[A-Z][a-z]+)*$")
	camelCaseRegex := regexp.MustCompile("^[a-z]+(?:[A-Z][a-z]+)*$")

	// Parse source code
	lang := golang.GetLanguage()
	n, err := sitter.ParseCtx(context.Background(), src, lang)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}
	// Execute the query
	q, err := sitter.NewQuery([]byte(filterFuncDecls), lang)
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
			if ident := c.Node.Content(src); !pascalCaseRegex.MatchString(ident) && !camelCaseRegex.MatchString(ident) {
				funcDecls = append(funcDecls, Ident{
					Name:    c.Node.Content(src),
					Line:    c.Node.StartPoint().Row,
					EndLine: c.Node.EndPoint().Row,
					Col:     c.Node.StartPoint().Column,
					EndCol:  c.Node.EndPoint().Column,
				})
			}
		}
	}
	return funcDecls, nil
}

func AnomalousConstDecls(src []byte) ([]Ident, error) {
	var constDecls []Ident
	filterConstDecls := `(
		const_spec (identifier) @constant
	)`

	pascalCaseRegex := regexp.MustCompile("^[A-Z][a-z]+(?:[A-Z][a-z]+)*$")
	camelCaseRegex := regexp.MustCompile("^[a-z]+(?:[A-Z][a-z]+)*$")

	// Parse source code
	lang := golang.GetLanguage()
	n, err := sitter.ParseCtx(context.Background(), src, lang)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}
	// Execute the query
	q, err := sitter.NewQuery([]byte(filterConstDecls), lang)
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
			if ident := c.Node.Content(src); !pascalCaseRegex.MatchString(ident) && !camelCaseRegex.MatchString(ident) {
				constDecls = append(constDecls, Ident{
					Name:    c.Node.Content(src),
					Line:    c.Node.StartPoint().Row,
					EndLine: c.Node.EndPoint().Row,
					Col:     c.Node.StartPoint().Column,
					EndCol:  c.Node.EndPoint().Column,
				})
			}
		}
	}
	return constDecls, nil
}

func AnomalousVarDecls(src []byte) ([]Ident, error) {
	var varDecls []Ident
	filterVarDecls := `(
		var_spec (identifier) @constant
	)`

	pascalCaseRegex := regexp.MustCompile("^[A-Z][a-z]+(?:[A-Z][a-z]+)*$")
	camelCaseRegex := regexp.MustCompile("^[a-z]+(?:[A-Z][a-z]+)*$")

	// Parse source code
	lang := golang.GetLanguage()
	n, err := sitter.ParseCtx(context.Background(), src, lang)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}
	// Execute the query
	q, err := sitter.NewQuery([]byte(filterVarDecls), lang)
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
			if ident := c.Node.Content(src); !pascalCaseRegex.MatchString(ident) && !camelCaseRegex.MatchString(ident) {
				varDecls = append(varDecls, Ident{
					Name:    c.Node.Content(src),
					Line:    c.Node.StartPoint().Row,
					EndLine: c.Node.EndPoint().Row,
					Col:     c.Node.StartPoint().Column,
					EndCol:  c.Node.EndPoint().Column,
				})
			}
		}
	}
	return varDecls, nil
}

// PackageMain returns true when there is a package main declared in the changed files.
func PackageMain(src []byte) (bool, error) {
	return false, nil
}