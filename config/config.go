package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
)

// Config is simple rename for interface. It does not have any specific functionality, yet.
type Config interface{}

// DataPath holds path for the folder in which data files are stored.
var DataPath string

// ConfigPath holds path for the folder in which config files are stored.
var ConfigPath string

// ReadTOML unmarshalizes file into the given Config struct.
func ReadTOML(cfg Config, filename string) error {
	cfgFile, err := os.Open(path.Join(ConfigPath, filename))
	if err != nil {
		return fmt.Errorf("unable to open %v file: %v", filename, err)
	}
	defer cfgFile.Close()

	b, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("unable to read contents of %v file: %v", filename, err)
	}

	err = toml.Unmarshal(b, cfg)
	if err != nil {
		return fmt.Errorf("unable to unmarshall contents of %s file: %v", filename, err)
	}

	return nil
}
