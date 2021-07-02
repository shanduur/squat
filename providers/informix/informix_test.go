package informix_test

import (
	"testing"

	"github.com/shanduur/squat/providers/informix"
)

var dsn = "Informix"

func TestNew(t *testing.T) {
	_, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}
}

func TestGetProviderName(t *testing.T) {
	ifx, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}

	if name := ifx.GetProviderName(); name != dsn {
		t.Errorf("provider name: got: %s, wanted %s", name, dsn)
	}
}

func TestGetTableDescription(t *testing.T) {
	ifx, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}

	if _, err := ifx.GetTableDescription("SL_POTR"); err != nil {
		t.Errorf("failed to get table description: %s", err.Error())
	}
}
