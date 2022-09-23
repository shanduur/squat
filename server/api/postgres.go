//go:build postgres

package api

import (
	"log"
	"os"
	"path"

	"github.com/shanduur/squat/providers/postgres"
)

func init() {
	if p, err := postgres.New(path.Join(os.Getenv("CONFIG_LOCATION"), "postgres.toml")); err != nil {
		log.Printf("unable to create new provider connection: %s", err.Error())
		log.Printf("check if env variables CONFIG_LOCATION and DATA_LOCATION are set")
	} else {
		Providers[p.ProviderName()] = p
	}
}
