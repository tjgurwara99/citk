package golang_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/tjgurwara99/citk/internal/golang"
)

func TestAnomalousFuncSignatures(t *testing.T) {
	file, err := os.ReadFile("./testdata/funcs.go")
	if err != nil {
		t.Fatalf("failed to open testdata/funcs.go: %s", err.Error())
	}
	funcs, err := golang.AnomalousFuncSignatures(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []golang.Ident{
		{
			Name:    "snake_case_function",
			Line:    14,
			EndLine: 14,
			Col:     5,
			EndCol:  24,
		}, {
			Name:    "SCREAMING_SNAKE_CASE_FUNCTION",
			Line:    16,
			EndLine: 16,
			Col:     5,
			EndCol:  34,
		}, {
			Name:    "SCREAMINGFUNCTION",
			Line:    18,
			EndLine: 18,
			Col:     5,
			EndCol:  22,
		},
	}

	if !reflect.DeepEqual(expected, funcs) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, funcs)
	}
}

func TestAnomalousConstDecls(t *testing.T) {
	file, err := os.ReadFile("./testdata/consts.go")
	if err != nil {
		t.Fatalf("failed to open testdata/consts.go: %s", err.Error())
	}
	funcs, err := golang.AnomalousConstDecls(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []golang.Ident{
		{
			Name:    "snake_case_const",
			Line:    5,
			EndLine: 5,
			Col:     1,
			EndCol:  17,
		}, {
			Name:    "SCREAMING_SNAKE_CASE_CONST",
			Line:    6,
			EndLine: 6,
			Col:     1,
			EndCol:  27,
		}, {
			Name:    "SCREAMINGCONST",
			Line:    7,
			EndLine: 7,
			Col:     1,
			EndCol:  15,
		},
	}

	if !reflect.DeepEqual(expected, funcs) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, funcs)
	}
}

func TestAnomalousVarDecls(t *testing.T) {
	file, err := os.ReadFile("./testdata/vars.go")
	if err != nil {
		t.Fatalf("failed to open testdata/vars.go: %s", err.Error())
	}
	funcs, err := golang.AnomalousVarDecls(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []golang.Ident{
		{
			Name:    "snake_case_var",
			Line:    5,
			EndLine: 5,
			Col:     1,
			EndCol:  15,
		}, {
			Name:    "SCREAMING_SNAKE_CASE_VAR",
			Line:    6,
			EndLine: 6,
			Col:     1,
			EndCol:  25,
		}, {
			Name:    "SCREAMINGVAR",
			Line:    7,
			EndLine: 7,
			Col:     1,
			EndCol:  13,
		},
	}

	if !reflect.DeepEqual(expected, funcs) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, funcs)
	}
}
