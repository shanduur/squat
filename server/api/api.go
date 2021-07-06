package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
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

	tab, err := parse(req.Form)
	if err != nil {
		log.Printf("unable to parse request form: %s", err.Error())
		return
	}

	gen, err := generator.New(path.Join(os.Getenv("DATA_LOCATION"), "data.gob"))
	if err != nil {
		log.Printf("unable to get generator: %s", err.Error())
		return
	}

	var out []string
	for i := 0; i <= 100; i++ {
		q, err := gen.Query(req.FormValue("source-table"), tab)
		if err != nil {
			log.Printf("unable to generate query: %s", err.Error())
			return
		}

		out = append(out, q)
	}

	w.Write([]byte(strings.Join(out, "\n")))
}

func parse(form url.Values) (map[string]generator.Column, error) {
	table := make(map[string]generator.Column)

	orders := make(map[string]int)
	names := make(map[string]string)
	types := make(map[string]string)
	lengths := make(map[string]int)
	precisions := make(map[string]int)
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

		} else if strings.Contains(k, "length-") {
			i, err := strconv.Atoi(v[0])
			if err != nil {
				return table, fmt.Errorf("unable to convert length: %s", err.Error())
			}

			lengths[strings.ReplaceAll(k, "length-", "")] = i

		} else if strings.Contains(k, "precision-") {
			i, err := strconv.Atoi(v[0])
			if err != nil {
				return table, fmt.Errorf("unable to convert precision: %s", err.Error())
			}

			precisions[strings.ReplaceAll(k, "precision-", "")] = i

		} else if strings.Contains(k, "order-") {
			i, err := strconv.Atoi(v[0])
			if err != nil {
				return table, fmt.Errorf("unable to convert order: %s", err.Error())
			}

			orders[strings.ReplaceAll(k, "order-", "")] = i

		}

	}

	for k := range names {
		var col generator.Column

		col.Order = orders[k]
		col.Name = names[k]
		col.Type = types[k]
		col.Length = lengths[k]
		col.Precision = precisions[k]
		col.TagRegex = tagsregexes[k]

		if includes[k] == "on" {
			col.Include = true
		} else {
			col.Include = false
		}

		table[k] = col
	}

	return table, nil
}
