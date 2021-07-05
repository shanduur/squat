package api

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/shanduur/squat/generator"
	"github.com/shanduur/squat/providers"
	"github.com/shanduur/squat/providers/informix"
)

func init() {
	Providers = make(map[string]providers.Provider)

	p, err := informix.New(path.Join(os.Getenv("CONFIG_LOCATION"), "informix.toml"))
	if err != nil {
		log.Fatalf("unable to create get new provider connection: %s", err.Error())
	}

	Providers[p.ProviderName()] = p
}

// Providers
var Providers map[string]providers.Provider

// RegisterEndpoints
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

	// out := fmt.Sprintf("-- generated on %v\n", time.Now().Format(p.DateTimeFormat()))

	tab := parse(req.Form)

	gen, err := generator.New(path.Join(os.Getenv("DATA_LOCATION"), "data.gob"))
	if err != nil {
		log.Printf("unable to get generator: %s", err.Error())
		return
	}

	q := gen.Query(tab)

	w.Write([]byte(q))
}

func parse(form url.Values) map[string]generator.Column {
	table := make(map[string]generator.Column)

	names := make(map[string]string)
	types := make(map[string]string)
	includes := make(map[string]string)
	tagsregexes := make(map[string]string)

	for k, v := range form {
		if strings.Contains(k, "include-") {
			includes[strings.ReplaceAll(k, "include-", "")] = v[0]

		} else if strings.Contains(k, "name-") {
			names[strings.ReplaceAll(k, "name-", "")] = v[0]

		} else if strings.Contains(k, "type-") {
			types[strings.ReplaceAll(k, "type-", "")] = v[0]

		} else if strings.Contains(k, "regex-") {
			tagsregexes[strings.ReplaceAll(k, "regex-", "")] = v[0]

		}
	}

	for k, _ := range names {
		var col generator.Column

		col.Name = names[k]
		col.Type = types[k]
		col.TagRegex = tagsregexes[k]

		if includes[k] == "on" {
			col.Include = true
		} else {
			col.Include = false
		}

		table[k] = col
	}

	return table
}
