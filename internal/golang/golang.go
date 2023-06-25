package golang

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

var (
	allCapsRE = regexp.MustCompile(`^[A-Z0-9_]+$`)
	anyCapsRE = regexp.MustCompile(`[A-Z]`)
)

// knownNameExceptions is a set of names that are known to be exempt from naming checks.
// This is usually because they are constrained by having to match names in the
// standard library.
var knownNameExceptions = map[string]bool{
	"LastInsertId": true, // must match database/sql
	"kWh":          true,
}

func checkCase(ident string) bool {
	if ident == "_" {
		return false
	}
	if check, ok := knownNameExceptions[ident]; ok {
		return check
	}

	// Handle two common styles from other languages that don't belong in Go.
	if len(ident) >= 5 && allCapsRE.MatchString(ident) && strings.Contains(ident, "_") {
		capCount := 0
		for _, c := range ident {
			if 'A' <= c && c <= 'Z' {
				capCount++
			}
		}
		if capCount >= 2 {
			return true
		}
	}

	if len(ident) > 2 && strings.Contains(ident[1:], "_") {
		return true
	}
	return false
}

func anomalousFuncSignatures(src []byte) ([]Ident, error) {
	filterFuncDecls := `(
		(function_declaration (identifier) @func)
	)`
	return anomalousDecls(src, filterFuncDecls, checkCase)
}

func anomalousConstDecls(src []byte) ([]Ident, error) {
	filterConstDecls := `(
		const_spec (identifier) @constant
	)`
	return anomalousDecls(src, filterConstDecls, checkCase)
}

func anomalousVarDecls(src []byte) ([]Ident, error) {
	filterVarDecls := `(
		var_spec (identifier) @constant
	)`
	return anomalousDecls(src, filterVarDecls, checkCase)
}

func anomalousMethodAndFieldDecls(src []byte) ([]Ident, error) {
	// method_declaration (field_identifier) @methods
	filterMethodDecls := `(
		((field_identifier) @field)
	)`
	return anomalousDecls(src, filterMethodDecls, checkCase)
}

func anomalousPackageName(src []byte) ([]Ident, error) {
	filterPackageName := `(
		((package_identifier) @field)
	)`
	return anomalousDecls(src, filterPackageName, func(ident string) bool {
		if strings.Contains(ident, "_") && !strings.HasSuffix(ident, "_test") {
			return true
		}
		if anyCapsRE.MatchString(ident) {
			return true
		}
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

func filterFiles(files []string, suffix string, srcDir string) []string {
	var filteredFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, suffix) {
			filteredFiles = append(filteredFiles, filepath.Join(srcDir, file))
		}
	}
	return filteredFiles
}

type inspectFunc func([]byte) ([]Ident, error)

type InspectFunc func([]byte, string, string) ([]annotation.Annotation, error)

func WrapInspectFuncs(f inspectFunc, title, msgFmt string, t annotation.AnnotationType) InspectFunc {
	return func(src []byte, baseDir, fName string) ([]annotation.Annotation, error) {
		idents, err := f(src)
		if err != nil {
			return nil, err
		}
		var annotations []annotation.Annotation
		for _, ident := range idents {
			f, err := filepath.Rel(baseDir, fName)
			if err != nil {
				return nil, err
			}
			annotations = append(annotations, annotation.Annotation{
				FileName:  f,
				Title:     title,
				Message:   fmt.Sprintf(msgFmt, ident.Name),
				Type:      t,
				StartLine: ident.Line,
				EndLine:   ident.EndLine,
				StartCol:  ident.Col,
				EndCol:    ident.EndCol,
			})
		}
		return annotations, nil
	}
}

func Inspect(srcDir string, relBranch string) ([]annotation.Annotation, error) {
	files, err := git.ListChangedFiles(srcDir, relBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve changed files from git: %w", err)
	}
	goFiles := filterFiles(files, ".go", srcDir)

	var annotations []annotation.Annotation
	inspectFuncs := []InspectFunc{
		WrapInspectFuncs(
			anomalousConstDecls,
			"Const declaration not following style guide",
			"The declaration of the const %s is not following our style guide. Please read our contribution guidelines and style guide to help you resolve this issue.",
			annotation.Error,
		),
		WrapInspectFuncs(
			anomalousFuncSignatures,
			"Func declaration not following our style guide",
			"The declaration of the function %s is not following our style guide. Please read our contribution guidelines and style guides to help you resolve this issue.",
			annotation.Error,
		),
		WrapInspectFuncs(
			anomalousMethodAndFieldDecls,
			"Struct fields/methods not following our style guide",
			"The declaration of the method/field %s is not following our style guide. Please read our contribution guidelines and style guides to help you resolve this issue.",
			annotation.Error,
		),
		WrapInspectFuncs(
			anomalousPackageName,
			"Package name not following our style guide",
			"The package declaration %s is not following our style guide. Please read our contribution guidelines and style guide to help you resolve this issue.",
			annotation.Error,
		),
		WrapInspectFuncs(
			anomalousVarDecls,
			"Variable name not following our style guide",
			"The variable declaration %s is not following our style guide. Please read our contribution guidelines and style guide to help you resolve this issue.",
			annotation.Error,
		),
	}

	for _, file := range goFiles {
		for _, inspector := range inspectFuncs {
			srcFile, err := os.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			idents, err := inspector(srcFile, srcDir, file)
			if err != nil {
				return nil, fmt.Errorf("failed to run inspector on src file: %w", err)
			}
			annotations = append(annotations, idents...)
		}
	}
	return annotations, nil
}
