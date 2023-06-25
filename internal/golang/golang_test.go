package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAnomalousFuncSignatures(t *testing.T) {
	file, err := os.ReadFile("./testdata/funcs.go")
	if err != nil {
		t.Fatalf("failed to open testdata/funcs.go: %s", err.Error())
	}
	funcs, err := anomalousFuncSignatures(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []Ident{
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
	funcs, err := anomalousConstDecls(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []Ident{
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
	funcs, err := anomalousVarDecls(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []Ident{
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
	pkgNames, err := anomalousPackageName(file)
	if err != nil {
		t.Errorf("unexpected error occured: %s", err)
	}
	expected := []Ident{
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
	funcs, err := anomalousMethodAndFieldDecls(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	expected := []Ident{
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
	}

	if !reflect.DeepEqual(expected, funcs) {
		t.Errorf("expected and returned values do not match: expected %+v, returned %+v", expected, funcs)
	}
}

func TestInspect(t *testing.T) {
	// we are skipping this for now since the testdata is a git repository
	// and git considers it as a submodule which we don't want. So we would
	// have to do a bit more complecated setup for this to work.
	t.Skip()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}
	annotations, err := Inspect(filepath.Join(wd, "../git/testdata"), "main")
	if err != nil {
		t.Errorf("failed to run Inspect: %s", err)
	}
	fmt.Printf("%+v", annotations)
}
