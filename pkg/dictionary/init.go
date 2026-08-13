package dictionary

import (
	"GoRestify/pkg/pkg_consts"
	"log"

	"github.com/BurntSushi/toml"
)

// Init terms and put in the main map
func Init() (err error) {

	thisTerms = make(map[string]Term)

	if _, err := toml.DecodeFile(pkg_consts.TermsFile, &thisTerms); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	return
}

// Lang is used for type of event
type Lang string

// Lang enums
const (
	En Lang = "en"
	Ku Lang = "ku"
	Ar Lang = "ar"
)

// Langs represents all accepted languages
var Langs = []Lang{
	En,
	Ku,
	Ar,
}

// Term is list of languages
type Term struct {
	En string `toml:"en"`
	Ku string `toml:"ku"`
	Ar string `toml:"ar"`
}

// thisTerms used for holding language identifier as a string and Term Struct as value
var thisTerms map[string]Term
