package generator_test

import (
	"os"
	"testing"

	"github.com/shanduur/squat/generator"
)

func TestReadDump(t *testing.T) {
	if err := generator.ReadDump("test/data.json", "test/data.gob"); err != nil {
		t.Errorf("unable to read and dump: %s", err.Error())
	}

	if err := os.Remove("test/data.gob"); err != nil {
		t.Errorf("unable to remove file: %s", err.Error())
	}
}
