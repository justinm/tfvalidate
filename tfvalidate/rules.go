package tfvalidate

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Parser struct {
}

type RuleSet struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	ResourceType string                    `json:"resource"`
	Attributes   []RuleAttributeDefinition `json:"attributes"`
}

type RuleAttributeDefinition struct {
	Name       string `json:"name"`
	BeginsWith string `json:"beginsWith"`
	In         []string
}

func (p Parser) Parse(path string) (*RuleSet, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	result := RuleSet{}

	err = json.Unmarshal([]byte(byteValue), &result)

	return &result, err
}
