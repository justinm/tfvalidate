package shared

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type Configuration struct {
	Rules     []Rule     `json:"rules"`
	Approvals []Approval `json:"approvals"`
}

type Rule struct {
	Name           string                    `yaml:"string"`
	ResourceTypes  []string                  `yaml:"resources"`
	RuleAttributes []RuleAttributeDefinition `yaml:"attributes"`
}

type RuleAttributeDefinition struct {
	Name  string                        `yaml:"name"`
	Rules []RuleAttributeDefinitionRule `yaml:"rules"`
}

type RuleAttributeDefinitionRule struct {
	StartsWith *string  `yaml:"startsWith"`
	OneOf      []string `yaml:"oneOf"`
	Required   *bool    `yaml:"required"`
}

type Approval struct {
	Resources []string
	Approvers []string
}

func GetConfig(path string) (*Configuration, []error) {
	config := &Configuration{}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, []error{err}
	}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, []error{err}
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, []error{err}
	}

	errs := validateConfig(config)

	if len(errs) > 0 {
		return nil, errs
	}

	return config, nil
}

func validateConfig(config *Configuration) []error {
	var errs []error

	if config == nil {
		return []error{errors.New("configuration is empty")}
	}

	for ruleNumber, rule := range config.Rules {
		if len(rule.ResourceTypes) == 0 {
			errs = append(errs, errors.New("Rule #"+strconv.Itoa(ruleNumber)+" did not define any resources"))
		}

		if len(rule.RuleAttributes) == 0 {
			errs = append(errs, errors.New("Rule #"+strconv.Itoa(ruleNumber)+" does not define any attributes"))
		}
	}

	return errs
}
