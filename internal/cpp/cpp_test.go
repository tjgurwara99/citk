package cpp

import (
	"os"
	"testing"
)

func TestAnomalousDefinitions(t *testing.T) {
	file, err := os.ReadFile("./testdata/cpp.cpp")
	if err != nil {
		t.Fatalf("failed to open testdata/cpp.cpp: %s", err.Error())
	}
	funcs, err := anomalousFunctions(file)
	if err != nil {
		t.Errorf("returned an error: %s", err)
	}

	_ = funcs
}
