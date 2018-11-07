package rules

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func Parse(path string) (*RuleSet, error) {
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
