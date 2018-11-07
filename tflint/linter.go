package tflint

import (
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
)

type Violation struct {
	Action Action
	Reason string
}

type Linter struct {
	Violations []Violation
}

type Action struct {
	ResourceKey  string
	ResourceType string
	AttrRule     *RuleAttributeDefinition
	Diff         *terraform.ResourceAttrDiff
	AttrName     string
}

func (l *Linter) getResourceFromKey(key string) string {
	parts := strings.Split(key, ".")

	return parts[0]
}

func (l *Linter) Lint(plan *terraform.Plan, set *RuleSet) {
	for _, module := range plan.Diff.Modules {
		for key, diff := range module.Resources {
			l.validateDiff(key, diff, set)
		}
	}
}

func (l *Linter) validateDiff(resourceKey string, diff *terraform.InstanceDiff, set *RuleSet) {
	resourceType := l.getResourceFromKey(resourceKey)

	for diffAttrName, diffAttr := range diff.Attributes {
		for _, rule := range set.Rules {
			if rule.ResourceType == resourceType {
				for _, attrRule := range rule.Attributes {
					if attrRule.Name == diffAttrName {
						action := Action{
							ResourceKey:  resourceKey,
							AttrRule:     &attrRule,
							AttrName:     diffAttrName,
							Diff:         diffAttr,
							ResourceType: resourceType,
						}

						l.validateAction(action)
					}
				}
			}
		}
	}

}

func (l *Linter) validateAction(action Action) {
	if len(action.AttrRule.BeginsWith) > 0 {
		l.validateBeginsWith(action)
	}
}

func (l *Linter) validateBeginsWith(action Action) {
	if !strings.HasPrefix(action.Diff.New, action.AttrRule.BeginsWith) {
		violation := Violation{
			Action: action,
			Reason: "Attribute " + action.ResourceKey + "." + action.AttrName + " did not begin with " + action.AttrRule.BeginsWith,
		}
		l.addViolation(violation)
	}
}

func (l *Linter) addViolation(violation Violation) {
	l.Violations = append(l.Violations, violation)
	if l.Violations == nil {
		log.Fatal("Violations was nil")
	}
}
