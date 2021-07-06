package config_test

import (
	"testing"

	"github.com/shanduur/squat/config"
)

type testCfg struct {
	Name   string `toml:"name" env:"NAME"`
	Number int    `toml:"number" env:"NUMBER"`
}

var (
	Name   = "test"
	Number = 1
)

func TestReadTOML(t *testing.T) {
	var cfg testCfg
	if err := config.ReadTOML(&cfg, "test/config.toml"); err != nil {
		t.Errorf("unable to read config: %s", err.Error())
	}

	if cfg.Name != Name {
		t.Errorf("wrong value: got: %s, wanted: %s", cfg.Name, Name)
	}

	if cfg.Number != Number {
		t.Errorf("wrong value: got: %d, wanted: %d", cfg.Number, Number)
	}
}
