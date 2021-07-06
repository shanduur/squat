package informix_test

import (
	"testing"

	"github.com/shanduur/squat/providers/informix"
)

const (
	dsn            = "Informix"
	dateFormat     = "2006-01-02"
	dateTimeFormat = "2006-01-02 15:04:05.000"
)

func TestNew(t *testing.T) {
	_, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}
}

func TestProviderName(t *testing.T) {
	ifx, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}

	if name := ifx.ProviderName(); name != dsn {
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

func TestDateFormat(t *testing.T) {
	ifx, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}

	if f := ifx.DateTimeFormat(); f != dateTimeFormat {
		t.Errorf("wrong date format: got: %s, wanted: %s", f, dateFormat)
	}
}

func TestDateTimeFormat(t *testing.T) {
	ifx, err := informix.New("test/informix.toml")
	if err != nil {
		t.Errorf("unable to get provider: %s", err.Error())
	}

	if f := ifx.DateTimeFormat(); f != dateTimeFormat {
		t.Errorf("wrong date time format: got: %s, wanted: %s", f, dateTimeFormat)
	}
}
