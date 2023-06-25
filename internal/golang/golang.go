package golang

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/tjgurwara99/citk/internal/annotation"
	"github.com/tjgurwara99/citk/internal/git"
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

func AnomalousPackageName(src []byte) ([]Ident, error) {
	filterPackageName := `(
		((package_identifier) @field)
	)`
	return anomalousDecls(src, filterPackageName, func(ident string) bool {
		return strings.ToLower(ident) != ident
	})
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
					Line:    c.Node.StartPoint().Row + 1,
					EndLine: c.Node.EndPoint().Row + 1,
					Col:     c.Node.StartPoint().Column,
					EndCol:  c.Node.EndPoint().Column,
				})
			}
		}
	}
	return decls, nil
}

func filterFiles(files []string, suffix string) []string {
	var filteredFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, suffix) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

func Inspect(srcDir string) ([]annotation.Annotation, error) {
	files, err := git.ListChangedFiles(srcDir, "main")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve changed files from git: %w", err)
	}
	goFiles := filterFiles(files, ".go")

	var annotations []annotation.Annotation
	inspectFuncs := []func([]byte) ([]Ident, error){
		AnomalousConstDecls,
		AnomalousFuncSignatures,
		AnomalousMethodAndFieldDecls,
		AnomalousPackageName,
		AnomalousVarDecls,
	}

	for _, file := range goFiles {
		for _, inspector := range inspectFuncs {
			srcFile, err := os.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			idents, err := inspector(srcFile)
			if err != nil {
				return nil, fmt.Errorf("failed to run inspector on src file: %w", err)
			}
		}
	}
	return nil, nil
}
