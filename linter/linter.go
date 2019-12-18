package linter

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"github.com/justinm/tfvalidate/shared"
	"github.com/justinm/tfvalidate/util"
	"github.com/op/go-logging"
	"strings"
)

var (
	logger = logging.MustGetLogger("linter")
)

type ResourceDescription struct {
	Name string
	Type string
}

type Violation struct {
	Reason      string
	Change      *terraform.InstanceDiff
	Attribute   shared.RuleAttributeDefinition
	Description ResourceDescription
	Value       *terraform.ResourceAttrDiff
}

type Linter struct {
	Config     *shared.Configuration
	Plan       *terraform.Plan
	Violations []Violation
}

func New(config *shared.Configuration, plan *terraform.Plan) (*Linter, []error) {
	linter := Linter{
		Config: config,
		Plan:   plan,
	}

	return &linter, nil
}

func (linter *Linter) Lint() []Violation {
	if linter.Config.Rules == nil {
		logger.Warning("no lint rules were found")
		return nil
	}

	for _, module := range linter.Plan.Diff.Modules {
		for key, resource := range module.Resources {
			_ = linter.LintChange(key, resource)
		}
	}

	return linter.Violations
}

func (linter *Linter) getResourceFromKey(key string) string {
	parts := strings.Split(key, ".")

	return parts[0]
}

func (linter *Linter) LintChange(resourceId string, change *terraform.InstanceDiff) error {
	parts := strings.Split(resourceId, ".")
	resourceName := parts[1]
	resourceType := parts[0]

	description := ResourceDescription{
		Name: resourceName,
		Type: resourceType,
	}

	if change.Destroy {
		logger.Debugf("Resource %s skipped, will be destroyed", resourceName)
		return nil
	}

	for _, rule := range linter.Config.Rules {
		if util.SliceStringContains(rule.ResourceTypes, resourceType) {
			for _, ruleAttribute := range rule.RuleAttributes {
				for _, subrule := range ruleAttribute.Rules {

					if subrule.Required != nil && *subrule.Required {
						violation := linter.validateRequired(description, ruleAttribute, change)
						linter.addViolation(violation)
					}

					if subrule.StartsWith != nil {
						violation := linter.validateStartsWith(description, ruleAttribute, subrule, change)
						linter.addViolation(violation)
					}

					if subrule.OneOf != nil {
						violation := linter.validateOneOf(description, ruleAttribute, subrule, change)
						linter.addViolation(violation)
					}
				}
			}
		}
	}

	return nil
}

func (linter *Linter) getAttributeValueFromDiff(name string, diff *terraform.InstanceDiff) *terraform.ResourceAttrDiff {
	for attrName, attr := range diff.Attributes {
		if attrName == name {
			return attr
		}
	}

	return nil
}

func (linter *Linter) validateRequired(description ResourceDescription, attr shared.RuleAttributeDefinition, diff *terraform.InstanceDiff) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff)

	if val == nil {
		return &Violation{
			Attribute:   attr,
			Change:      diff,
			Description: description,
			Reason:      "required but undefined",
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) validateOneOf(description ResourceDescription, attr shared.RuleAttributeDefinition, rule shared.RuleAttributeDefinitionRule, diff *terraform.InstanceDiff) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff)
	if val == nil {
		return nil
	}

	if !util.SliceStringContains(rule.OneOf, val.New) {
		return &Violation{
			Attribute:   attr,
			Change:      diff,
			Description: description,
			Reason:      fmt.Sprintf("does not contain one of %v", rule.OneOf),
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) validateStartsWith(description ResourceDescription, attr shared.RuleAttributeDefinition, rule shared.RuleAttributeDefinitionRule, diff *terraform.InstanceDiff) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff)
	if val == nil {
		return nil
	}

	if !strings.HasPrefix(val.New, *rule.StartsWith) {
		return &Violation{
			Attribute:   attr,
			Change:      diff,
			Description: description,
			Reason:      "does not start with " + *rule.StartsWith,
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) addViolation(violation *Violation) {
	if violation != nil {
		linter.Violations = append(linter.Violations, *violation)
	}
}
