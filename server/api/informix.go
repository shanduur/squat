//go:build informix

package api

import (
	"log"
	"os"
	"path"

	"github.com/shanduur/squat/providers/informix"
)

func init() {
	if p, err := informix.New(path.Join(os.Getenv("CONFIG_LOCATION"), "informix.toml")); err != nil {
		log.Printf("unable to create new provider connection: %s", err.Error())
		log.Printf("check if env variables CONFIG_LOCATION and DATA_LOCATION are set")
	} else {
		Providers[p.ProviderName()] = p
	}
}
