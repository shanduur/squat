package ui

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/shanduur/squat/generator"
	"github.com/shanduur/squat/providers"
	"github.com/shanduur/squat/server/api"
)

var (
	//go:embed templates/main.html
	mainHTML string
	//go:embed templates/.Head.html
	headHTML string
	//go:embed templates/.Body.index.html
	bodyIndexHTML string
	//go:embed templates/.Body.table.html
	bodyTableHTML string
	//go:embed templates/.Script.table.js
	scriptTableJS string
	//go:embed templates/partials/.Option.html
	partialOptionHTML string
	//go:embed templates/partials/.Rows.html
	partialRowsHTML string
	//go:embed templates/partials/.Sources.html
	partialSourcesHTML string

	template = map[string]string{
		"head":      "{{ .Head }}",
		"body":      "{{ .Body }}",
		"rows":      "{{ .Rows }}",
		"name":      "{{ .Name }}",
		"type":      "{{ .Type }}",
		"order":     "{{ .Order }}",
		"table":     "{{ .Table }}",
		"script":    "{{ .Script }}",
		"fomrat":    "{{ .Format }}",
		"length":    "{{ .Length }}",
		"source":    "{{ .Source }}",
		"options":   "{{ .Options }}",
		"sources":   "{{ .Sources }}",
		"display":   "{{ .Display }}",
		"template":  "{{ .Template }}",
		"precision": "{{ .Precision }}",
	}
)

func RegisterEndpoints(mux *mux.Router) {
	mux.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	mux.HandleFunc("/table", handleTable).Methods(http.MethodGet)
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(buildIndex()))
}

func handleTable(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	src := req.FormValue("source")
	tab := req.FormValue("table")

	p := api.Providers[src]
	if p != nil {
		dsc, err := p.GetTableDescription(tab)
		if err != nil {
			log.Printf("unable to get table description: %s", err.Error())
		}

		out, err := buildTables(src, tab, dsc)
		if err != nil {
			log.Printf("unable to build tables: %s", err.Error())
			return
		}
		w.Write([]byte(out))
	}
}

func buildTables(src string, tab string, dsc []providers.Describe) (string, error) {
	output := strings.ReplaceAll(mainHTML, template["head"], headHTML)
	output = strings.ReplaceAll(output, template["body"], bodyTableHTML)
	output = strings.ReplaceAll(output, template["script"], scriptTableJS)

	var options string
	g, err := generator.New(path.Join(os.Getenv("DATA_LOCATION"), "data.gob"))
	if err != nil {
		return "", fmt.Errorf("unable to load generator")
	}

	for k, v := range g.TagsAndRegex {
		opt := strings.ReplaceAll(partialOptionHTML, template["name"], fmt.Sprint(v))
		opt = strings.ReplaceAll(opt, template["display"], k)
		options = fmt.Sprintf("%s\n%s", options, opt)
	}

	output = strings.ReplaceAll(output, template["template"], options)

	var rows string
	for i, d := range dsc {
		opt := options
		row := partialRowsHTML
		row = strings.ReplaceAll(row, template["order"], fmt.Sprint(i))

		if d.ColumnType.Valid {
			row = strings.ReplaceAll(row, template["type"], d.ColumnType.String)
		}

		if d.ColumnLength.Valid {
			row = strings.ReplaceAll(row, template["length"], fmt.Sprintf("%d", d.ColumnLength.Int64))
		}

		if d.ColumnPrecision.Valid {
			row = strings.ReplaceAll(row, template["precision"], fmt.Sprintf("%d", d.ColumnPrecision.Int64))
		} else {
			row = strings.ReplaceAll(row, template["precision"], "0")
		}

		if strings.Contains(d.ColumnType.String, generator.TypeInt) {
			opt = strings.ReplaceAll(opt, generator.RegexNumber+`"`, generator.RegexNumber+`"`+`selected="selected"`)

		} else if strings.Contains(d.ColumnType.String, generator.TypeChar) {
			opt = strings.ReplaceAll(opt, generator.RegexWord+`"`, generator.RegexWord+`"`+`selected="selected"`)

		} else if strings.Contains(d.ColumnType.String, generator.TypeFloat) {
			opt = strings.ReplaceAll(opt, generator.RegexNumber+`"`, generator.RegexNumber+`"`+`selected="selected"`)

		} else if strings.Contains(d.ColumnType.String, generator.TypeDecimal) {
			opt = strings.ReplaceAll(opt, generator.RegexNumber+`"`, generator.RegexNumber+`"`+`selected="selected"`)

		} else if strings.Contains(d.ColumnType.String, generator.TypeTimestamp) {
			opt = strings.ReplaceAll(opt, generator.TagTimestamp+`"`, generator.TagTimestamp+`"`+`selected="selected"`)

		} else if strings.Contains(d.ColumnType.String, generator.TypeDateTime) {
			opt = strings.ReplaceAll(opt, generator.TagDateTime+`"`, generator.TagDateTime+`"`+`selected="selected"`)

		} else if strings.Contains(d.ColumnType.String, generator.TypeDate) {
			opt = strings.ReplaceAll(opt, generator.TagDate+`"`, generator.TagDate+`"`+`selected="selected"`)

		}

		if d.ColumnName.Valid {
			row = strings.ReplaceAll(row, template["name"], d.ColumnName.String)
		}

		row = strings.ReplaceAll(row, template["options"], opt)

		rows = fmt.Sprintf("%s\n%s", rows, row)
	}

	output = strings.ReplaceAll(output, template["rows"], rows)
	output = strings.ReplaceAll(output, template["source"], src)
	output = strings.ReplaceAll(output, template["table"], tab)

	return output, nil
}

func buildIndex() string {
	output := strings.ReplaceAll(mainHTML, template["head"], headHTML)
	output = strings.ReplaceAll(output, template["body"], bodyIndexHTML)
	output = strings.ReplaceAll(output, template["script"], "")
	output = strings.ReplaceAll(output, template["sources"], partialSourcesHTML)

	var options string
	for i := range api.Providers {
		opt := strings.ReplaceAll(partialOptionHTML, template["name"], i)
		opt = strings.ReplaceAll(opt, template["display"], i)
		options = fmt.Sprintf("%s\n%s", options, opt)
	}

	output = strings.ReplaceAll(output, template["options"], options)

	return output
}
