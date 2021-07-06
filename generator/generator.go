package generator

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/lucasjones/reggen"
)

type Generator struct {
	dictionary   Dictionary
	TagsAndRegex map[string]string
}

type Dictionary struct {
	Names     []string `json:"names"`
	Surnames  []string `json:"surnames"`
	Cities    []string `json:"cities"`
	States    []string `json:"states"`
	Countries []string `json:"countries"`
}

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

func (g Generator) Get(tag string) (string, error) {
	switch tag {
	case TagName:
		return g.dictionary.Names[rand.Intn(len(g.dictionary.Names))], nil
	case TagSurname:
		return g.dictionary.Surnames[rand.Intn(len(g.dictionary.Surnames))], nil
	case TagCity:
		return g.dictionary.Cities[rand.Intn(len(g.dictionary.Cities))], nil
	case TagState:
		return g.dictionary.States[rand.Intn(len(g.dictionary.States))], nil
	case TagCountry:
		return g.dictionary.Countries[rand.Intn(len(g.dictionary.Countries))], nil
	case TagDate:
		return "", nil
	case TagDateTime:
		return "", nil
	case TagTimestamp:
		return "", nil
	case TagBool:
		return fmt.Sprint(rand.Intn(1) != 0), nil
	case TagYesNo:
		return randYesNo(), nil
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
		out = fmt.Sprintf(`"%s"`, out)
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
	(*m)["Phone"] = RegexPhone
	(*m)["E-Mail"] = RegexEmail
	(*m)["Postal Code"] = RegexPostalCode
	(*m)["PESEL"] = RegexPESEL
	(*m)["NIP"] = RegexNIP
	(*m)["REGON"] = RegexREGON
	(*m)["Word"] = RegexWord
	(*m)["Number"] = RegexNumber
}

type Column struct {
	Order     int
	Include   bool
	Name      string
	Type      string
	Length    int
	Precision int
	TagRegex  string
}

func (gen Generator) Query(table string, dsc map[string]Column) (string, error) {
	query := "INSERT INTO %s (%s) \nVALUES (%s);"

	columns := make(map[int]string)
	values := make(map[int]string)

	for k, v := range dsc {
		if !v.Include {
			log.Printf("%s", k)
			continue
		}

		columns[v.Order] = k

		if strings.Contains(v.TagRegex, "@") {
			s, err := gen.Get(v.TagRegex)
			if err != nil {
				return "", fmt.Errorf("unable to obtain value %s: %s", v.TagRegex, err.Error())
			}

			values[v.Order] = s
			continue
		}

		s, err := gen.Generate(v.TagRegex, v.Length, v.Type)
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
