package generator

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"

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
	default:
		return "", ErrNotInDict
	}
}

func (g Generator) Generate(regex string, limit int) (string, error) {
	return reggen.Generate(regex, limit)
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
	Include  bool
	Name     string
	Type     string
	TagRegex string
}

func (gen Generator) Query(dsc map[string]Column) string {
	// query := "INSERT INTO (%s) VALUES (%s);"
	// var (
	// 	columns []string
	// 	values  []string
	// )

	// var template string
	// for _, d := range dsc {
	// 	columns = append(columns, d.ColumnName.String)

	// 	switch d.ColumnName.String {
	// 	case TypeChar:
	// 		fallthrough
	// 	case TypeDate:
	// 		fallthrough
	// 	case TypeDateTime:
	// 		fallthrough
	// 	case TypeTimestamp:
	// 		template = `"%s"`

	// 	default:
	// 		template = `%s`
	// 	}

	// 	values = append(values, fmt.Sprintf(template))
	// }

	out := ""
	for k, v := range dsc {
		out = fmt.Sprintf("%s\n%s: %+v", out, k, v)
	}

	return out
}
