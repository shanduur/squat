package generator_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/shanduur/squat/generator"
)

const queries = 1000

var table = map[string]generator.Column{
	"INT1": {
		Order:     0,
		Include:   true,
		Name:      "INT1",
		Type:      "INTEGER",
		Length:    8,
		Precision: 0,
		TagRegex:  "@integer",
	},
	"INT2": {
		Order:     1,
		Include:   true,
		Name:      "INT2",
		Type:      "INTEGER",
		Length:    8,
		Precision: 0,
		TagRegex:  "@integer"},
	"CHAR1": {
		Order:     2,
		Include:   true,
		Name:      "CHAR1",
		Type:      "INTEGER",
		Length:    5,
		Precision: 0,
		TagRegex:  "@integer",
	},
}

// TestQuery tests if the values generated by the Query are mostly unique.
func TestQuery(t *testing.T) {
	var qs []string

	if err := generator.ReadDump("test/data.json", "test/data.gob"); err != nil {
		t.Errorf("unable to create test gob file: %s", err.Error())
	}

	gen, err := generator.New("test/data.gob")
	if err != nil {
		t.Errorf("failed to create new generator: %s", err.Error())
	}

	for i := 0; i < queries; i++ {
		gen.SetSeed(int64(i))
		q, err := gen.Query("test", table)
		if err != nil {
			t.Errorf("unable to generate: %s", err.Error())
		}

		qs = append(qs, q)
	}

	for i, q1 := range qs {
		for j, q2 := range qs {
			if i == j {
				continue
			}

			if reflect.DeepEqual(q1, q2) {
				t.Errorf("found two equal rows: \n- id(%d): %s\n- id(%d): %s", i, q1, j, q2)
				return
			}
		}
	}

	if err := os.Remove("test/data.gob"); err != nil {
		t.Errorf("unable to remove file: %s", err.Error())
	}
}
