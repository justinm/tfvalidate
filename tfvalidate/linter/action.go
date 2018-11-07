package linter

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/justinm/tfvalidate/tfvalidate/rules"
)

type Action struct {
	ResourceKey  string
	ResourceType string
	AttrRule     *rules.RuleAttributeDefinition
	Diff         *terraform.ResourceAttrDiff
	AttrName     string
}
