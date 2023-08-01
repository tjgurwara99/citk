package cpp

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

type Ident struct {
	Name    string
	Line    uint32
	EndLine uint32
	Col     uint32
	EndCol  uint32
}

func filterSrc(src []byte, query string) ([]*sitter.Node, error) {
	// Parse source code
	lang := cpp.GetLanguage()
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

	var funcs []*sitter.Node
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, src)
		for _, c := range m.Captures {
			funcs = append(funcs, c.Node)
		}
	}
	return funcs, nil
}

func anomalousFunctions(src []byte) ([]Ident, error) {
	query := `
		(function_definition) @func
	`
	funcs, err := filterSrc(src, query)
	if err != nil {
		return nil, err
	}
	for _, f := range funcs {
		data := f.ChildByFieldName("declarator")
		data = data.ChildByFieldName("declarator")
		switch data.Type() {
		case "identifier":
			// normal function
			fmt.Println(data.Content(src))
		case "field_identifier":
			// class function defined inside the class
			fmt.Println(data.Content(src))
		case "qualifier_identifier":
			// class function defined outside the class with namespace operator
			fmt.Println(data.Content(src))
		}
	}
	return nil, nil
}
