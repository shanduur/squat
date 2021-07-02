package generator

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"

	"github.com/lucasjones/reggen"
)

var ErrNotInDict = errors.New("tag not found in dictionary")

const (
	TagName      = "@name"
	TagSurname   = "@surname"
	TagCity      = "@city"
	TagState     = "@state"
	TagCountry   = "@country"
	TagDate      = "@date"
	TagDateTime  = "@datetime"
	TagTimestamp = "@timestamp"
	TagYesNo     = "@yn"
	TagBool      = "@bool"

	RegexPhone      = `^(\d{9}|\+\d{11})$`
	RegexEmail      = `^[a-z]{5,10}@[a-z]{5,10}\.(com|net|org)$`
	RegexPostalCode = `^(\d{2})-(\d{3})$`
	RegexPESEL      = `^(\d{11})`
	RegexNIP        = `^(\d{10})`
	RegexREGON      = `^(\d{9})`
	RegexWord       = `^([A-Z][a-z]+)(-[A-Z][a-z]+)?$`
	RegexNumber     = `^(\d*)$`
)

type Dictionary struct {
	Names     []string `json:"names"`
	Surnames  []string `json:"surnames"`
	Cities    []string `json:"cities"`
	States    []string `json:"states"`
	Countries []string `json:"countries"`
}

var (
	TagsRegex  map[string]string
	dictionary Dictionary
	gobFile    = "data.gob"
	jsonFile   = "data.json"
)

func init() {
	gobFile = path.Join("bin", "data", gobFile)
	jsonFile = path.Join("bin", "data", jsonFile)

	if _, err := os.Stat(gobFile); errors.Is(err, os.ErrNotExist) {
		log.Printf("%s file not found, attempting to create it", gobFile)
		if err = ReadDump(jsonFile, gobFile); err != nil {
			log.Fatalf("failed to create %s file: %s", gobFile, err.Error())
		}
	}

	file, err := os.Open(gobFile)
	if err != nil {
		log.Fatalf("unable to open %s file: %s", gobFile, err.Error())
	}
	defer file.Close()

	dec := gob.NewDecoder(file)

	if err = dec.Decode(&dictionary); err != nil {
		log.Fatalf("decoding %s failed: %s", gobFile, err.Error())
	}

	TagsRegex = make(map[string]string)
	TagsRegex["Name"] = TagName
	TagsRegex["Surname"] = TagSurname
	TagsRegex["City"] = TagCity
	TagsRegex["State"] = TagState
	TagsRegex["Country"] = TagCountry
	TagsRegex["Date"] = TagDate
	TagsRegex["Date and Time"] = TagDateTime
	TagsRegex["Timestamp"] = TagTimestamp
	TagsRegex["Yes or No"] = TagYesNo
	TagsRegex["Boolean"] = TagBool

	TagsRegex["Phone"] = RegexPhone
	TagsRegex["E-Mail"] = RegexEmail
	TagsRegex["Postal Code"] = RegexPostalCode
	TagsRegex["PESEL"] = RegexPESEL
	TagsRegex["NIP"] = RegexNIP
	TagsRegex["REGON"] = RegexREGON
	TagsRegex["Word"] = RegexWord
	TagsRegex["Number"] = RegexNumber
}

func ReadDump(in string, out string) error {
	cfgFile, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("unable to open %s file: %s", in, err.Error())
	}
	defer cfgFile.Close()

	b, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("unable to read contents of %s file: %s", in, err.Error())
	}

	err = json.Unmarshal(b, &dictionary)
	if err != nil {
		return fmt.Errorf("unable to unmarshall contents of %s file: %s", in, err.Error())
	}

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err = enc.Encode(dictionary); err != nil {
		return fmt.Errorf("unable to encode %s file: %s", in, err.Error())
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("unable to create output file %s: %s", out, err.Error())
	}
	defer f.Close()

	if _, err := f.Write(buff.Bytes()); err != nil {
		return fmt.Errorf("unable to write output file %s: %s", out, err.Error())
	}

	return nil
}

func Get(tag string) (string, error) {
	switch tag {
	case TagName:
		return dictionary.Names[rand.Intn(len(dictionary.Names))], nil
	case TagSurname:
		return dictionary.Surnames[rand.Intn(len(dictionary.Surnames))], nil
	case TagCity:
		return dictionary.Cities[rand.Intn(len(dictionary.Cities))], nil
	case TagState:
		return dictionary.States[rand.Intn(len(dictionary.States))], nil
	case TagCountry:
		return dictionary.Countries[rand.Intn(len(dictionary.Countries))], nil
	default:
		return "", ErrNotInDict
	}
}

func Generate(regex string, limit int) (string, error) {
	return reggen.Generate(regex, limit)
}
