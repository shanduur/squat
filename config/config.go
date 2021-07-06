package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
)

type Config interface{}

var DataPath string
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
