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
			Line:    15,
			EndLine: 15,
			Col:     5,
			EndCol:  24,
		}, {
			Name:    "SCREAMING_SNAKE_CASE_FUNCTION",
			Line:    17,
			EndLine: 17,
			Col:     5,
			EndCol:  34,
		}, {
			Name:    "SCREAMINGFUNCTION",
			Line:    19,
			EndLine: 19,
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
			Line:    6,
			EndLine: 6,
			Col:     1,
			EndCol:  17,
		}, {
			Name:    "SCREAMING_SNAKE_CASE_CONST",
			Line:    7,
			EndLine: 7,
			Col:     1,
			EndCol:  27,
		}, {
			Name:    "SCREAMINGCONST",
			Line:    8,
			EndLine: 8,
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
			Line:    6,
			EndLine: 6,
			Col:     1,
			EndCol:  15,
		}, {
			Name:    "SCREAMING_SNAKE_CASE_VAR",
			Line:    7,
			EndLine: 7,
			Col:     1,
			EndCol:  25,
		}, {
			Name:    "SCREAMINGVAR",
			Line:    8,
			EndLine: 8,
			Col:     1,
			EndCol:  13,
		},
	}

	if !reflect.DeepEqual(expected, funcs) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, funcs)
	}
}

func TestAnomalousPackageName(t *testing.T) {
	file, err := os.ReadFile("./testdata/weird_package_name.go")
	if err != nil {
		t.Fatalf("failed to open testdata/weird_package_name.go: %s", err.Error())
	}
	pkgNames, err := golang.AnomalousPackageName(file)
	if err != nil {
		t.Errorf("unexpected error occured: %s", err)
	}
	expected := []golang.Ident{
		{
			Name:    "SomeThing",
			Line:    1,
			EndLine: 1,
			Col:     8,
			EndCol:  17,
		},
	}
	if !reflect.DeepEqual(expected, pkgNames) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, pkgNames)
	}
}

func TestAnomalousMethodAndFieldDecls(t *testing.T) {
	file, err := os.ReadFile("./testdata/methods.go")
	if err != nil {
		t.Fatalf("failed to open testdata/methods.go: %s", err.Error())
	}
	funcs, err := golang.AnomalousMethodAndFieldDecls(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []golang.Ident{
		{
			Name:    "snake_case_field",
			Line:    6,
			EndLine: 6,
			Col:     1,
			EndCol:  17,
		},
		{
			Name:    "SCREAMING_SNAKE_CASE_FIELD",
			Line:    7,
			EndLine: 7,
			Col:     1,
			EndCol:  27,
		},
		{
			Name:    "SCREAMINGFIELD",
			Line:    8,
			EndLine: 8,
			Col:     1,
			EndCol:  15,
		},
		{
			Name:    "snake_case_method",
			Line:    13,
			EndLine: 13,
			Col:     17,
			EndCol:  34,
		},
		{
			Name:    "SCREAMING_SNAKE_CASE_METHOD",
			Line:    14,
			EndLine: 14,
			Col:     17,
			EndCol:  44,
		},
		{
			Name:    "SCREAMINGMETHOD",
			Line:    15,
			EndLine: 15,
			Col:     17,
			EndCol:  32,
		},
	}

	if !reflect.DeepEqual(expected, funcs) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, funcs)
	}
}

func TestInspect(t *testing.T) {
	golang.Inspect("/Users/taj/personal/citk")
}
