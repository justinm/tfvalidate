package rules

type RuleSet struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	ResourceTypes  []string                  `json:"resources"`
	RuleAttributes []RuleAttributeDefinition `json:"attributes"`
}

type RuleAttributeDefinition struct {
	Name       string   `json:"name"`
	BeginsWith string   `json:"beginsWith"`
	In         []string `json:"in"`
}
