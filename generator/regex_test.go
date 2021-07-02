package generator_test

import (
	"testing"

	"github.com/shanduur/squat/generator"
)

func TestReadDump(t *testing.T) {
	if err := generator.ReadDump("bin/data/data.json", "bin/data/data.gob"); err != nil {
		t.Errorf("unable to read and dump: %s", err.Error())
	}
}
