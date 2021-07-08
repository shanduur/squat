/*
Package generator implements generator, that is able to create Insert Query,
as well as load dictionary from the gob file. It wraps reggen package for creating
synthetic data directly from Regular Expressions.
*/
package generator

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/lucasjones/reggen"
	"github.com/shanduur/squat/providers"
)

// Generator is the default generator struct definition
type Generator struct {
	dictionary    Dictionary
	TagsAndRegex  map[string]string
	DateTempl     string
	DateTimeTempl string
}

// Dictionary is a structure holding data read from the Dictionary gob file.
type Dictionary struct {
	Names     []string `json:"names"`
	Surnames  []string `json:"surnames"`
	Cities    []string `json:"cities"`
	States    []string `json:"states"`
	Countries []string `json:"countries"`
}

// New creates new generator and loads gob file into dictionary.
func New(path string) (Generator, error) {
	gen := Generator{}

	gen.TagsAndRegex = make(map[string]string)
	loadMap(&gen.TagsAndRegex)

	file, err := os.Open(path)
	if err != nil {
		return gen, fmt.Errorf("unable to open %s file: %s", path, err.Error())
	}
	defer file.Close()

	dec := gob.NewDecoder(file)

	if err = dec.Decode(&gen.dictionary); err != nil {
		return gen, fmt.Errorf("decoding %s failed: %s", path, err.Error())
	}

	return gen, nil
}

// SetTemplates is used to set templates specific for the providers.
func (g *Generator) SetTemplates(p providers.Provider) {
	g.DateTempl = p.DateFormat()
	g.DateTimeTempl = p.DateTimeFormat()
}

// Get returns random element from the dictionary according to the tag.
func (g Generator) Get(tag string) (string, error) {
	switch tag {
	case TagName:
		return fmt.Sprintf(`'%s'`, strings.ReplaceAll(g.dictionary.Names[rand.Intn(len(g.dictionary.Names))], "'", "`")), nil
	case TagSurname:
		return fmt.Sprintf(`'%s'`, strings.ReplaceAll(g.dictionary.Surnames[rand.Intn(len(g.dictionary.Surnames))], "'", "`")), nil
	case TagCity:
		return fmt.Sprintf(`'%s'`, strings.ReplaceAll(g.dictionary.Cities[rand.Intn(len(g.dictionary.Cities))], "'", "`")), nil
	case TagState:
		return fmt.Sprintf(`'%s'`, strings.ReplaceAll(g.dictionary.States[rand.Intn(len(g.dictionary.States))], "'", "`")), nil
	case TagCountry:
		return fmt.Sprintf(`'%s'`, strings.ReplaceAll(g.dictionary.Countries[rand.Intn(len(g.dictionary.Countries))], "'", "`")), nil
	case TagDate:
		return fmt.Sprintf(`'%s'`, time.Now().Format(g.DateTempl)), nil
	case TagDateTime:
		return fmt.Sprintf(`'%s'`, time.Now().Format(g.DateTimeTempl)), nil
	case TagTimestamp:
		return fmt.Sprintf(`'%s'`, time.Now().Format(g.DateTimeTempl)), nil
	case TagYesNo:
		return fmt.Sprintf(`'%s'`, randYesNo()), nil
	case TagInteger:
		return fmt.Sprint(rand.Int()), nil
	case TagDecimal:
		return fmt.Sprintf("%f", rand.Float64()*float64(rand.Int())), nil
	case TagBool:
		return fmt.Sprint(rand.Intn(1) != 0), nil
	default:
		return "", ErrNotInDict
	}
}

func randYesNo() string {
	if rand.Intn(1) != 0 {
		return "Yes"
	}

	return "No"
}

// Generate generates the data based on provided REGEX, length limit, and type.
func (g Generator) Generate(regex string, limit int, t string) (string, error) {
	out, err := reggen.Generate(regex, limit)
	if err != nil {
		return out, fmt.Errorf("unable to generate from %s with limit %d: %s", regex, limit, err.Error())
	}

	switch t {
	case TypeChar:
		fallthrough
	case TypeDate:
		fallthrough
	case TypeDateTime:
		fallthrough
	case TypeTimestamp:
		out = fmt.Sprintf(`'%s'`, out)
	}

	return out, nil
}

func loadMap(m *map[string]string) {
	(*m)["Name"] = TagName
	(*m)["Surname"] = TagSurname
	(*m)["City"] = TagCity
	(*m)["State"] = TagState
	(*m)["Country"] = TagCountry
	(*m)["Date"] = TagDate
	(*m)["Date and Time"] = TagDateTime
	(*m)["Timestamp"] = TagTimestamp
	(*m)["Yes or No"] = TagYesNo
	(*m)["Boolean"] = TagBool
	(*m)["Decimal"] = TagDecimal
	(*m)["Integer"] = TagInteger
	(*m)["Phone"] = RegexPhone
	(*m)["E-Mail"] = RegexEmail
	(*m)["Postal Code"] = RegexPostalCode
	(*m)["PESEL"] = RegexPESEL
	(*m)["NIP"] = RegexNIP
	(*m)["REGON"] = RegexREGON
	(*m)["Word"] = RegexWord
}

// Column describes each column after parsing the request.
type Column struct {
	Order     int
	Include   bool
	Name      string
	Type      string
	Length    int
	Precision int
	TagRegex  string
}

// Query builds insert query based on the table description.
func (g Generator) Query(table string, dsc map[string]Column) (string, error) {
	query := "INSERT INTO %s (%s) \nVALUES (%s);"

	columns := make(map[int]string)
	values := make(map[int]string)

	for k, v := range dsc {
		if !v.Include {
			continue
		}

		columns[v.Order] = k

		if strings.Contains(v.TagRegex, "@") {
			s, err := g.Get(v.TagRegex)
			if err != nil {
				return "", fmt.Errorf("unable to obtain value %s: %s", v.TagRegex, err.Error())
			}

			values[v.Order] = s
			continue
		}

		s, err := g.Generate(v.TagRegex, v.Length, v.Type)
		if err != nil {
			return "", fmt.Errorf("unable to generate from %s and length %d: %s", v.TagRegex, v.Length, err.Error())
		}
		values[v.Order] = s
	}

	keys := make([]int, 0)
	for k := range columns {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	c := make([]string, 0)
	v := make([]string, 0)

	for _, i := range keys {
		c = append(c, columns[i])
		v = append(v, values[i])
	}

	return fmt.Sprintf(query, table, strings.Join(c, ", "), strings.Join(v, ", ")), nil
}
