package api

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/shanduur/squat/providers"
	"github.com/shanduur/squat/providers/informix"
)

var Providers map[string]providers.Provider

func init() {
	Providers = make(map[string]providers.Provider)

	p, err := informix.New(path.Join("bin", "config", "informix.toml"))
	if err != nil {
		log.Printf("unable to create get new provider connection: %s", err.Error())
		return
	}

	Providers[p.GetProviderName()] = p
}

func RegisterEndpoints(mux *mux.Router) {
	mux.HandleFunc("/generate", generate).Methods(http.MethodGet)
}

func generate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	src := req.FormValue("source")

	p := Providers[src]
	if p == nil {
		log.Printf("data source not found: %s", src)
		return
	}

	dsc, err := p.GetTableDescription("SL_POTR")
	if err != nil {
		log.Printf("unable to get table description: %s", err.Error())
		return
	}

	var out string
	for _, d := range dsc {
		out = fmt.Sprintf("%s\n -- %+v | %+v | %+v | %+v | %+v", out, d.ColumnName, d.ColumnType, d.ColumnLength, d.ColumnPrecision, d.Nullable)
	}

	w.Write([]byte(out))
}
