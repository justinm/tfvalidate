package linter

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/justinm/tfvalidate/tfvalidate/rules"
	"github.com/justinm/tfvalidate/tfvalidate/util"
	"log"
	"strconv"
	"strings"
)

type Linter struct {
	Violations []Violation
}

func (l *Linter) getResourceFromKey(key string) string {
	parts := strings.Split(key, ".")

	return parts[0]
}

func (l *Linter) Lint(plan *terraform.Plan, set *rules.RuleSet) {
	for _, module := range plan.Diff.Modules {
		for key, diff := range module.Resources {
			l.validateDiff(key, diff, set)
		}
	}
}

func (l *Linter) validateDiff(resourceKey string, diff *terraform.InstanceDiff, set *rules.RuleSet) {
	resourceType := l.getResourceFromKey(resourceKey)

	for diffAttrName, diffAttr := range diff.Attributes {
		for i, rule := range set.Rules {

			if len(rule.ResourceTypes) == 0 {
				log.Fatal("Rule #" + strconv.Itoa(i) + " did not define any resources")
			}

			if len(rule.RuleAttributes) == 0 {
				log.Fatal("Rule #" + strconv.Itoa(i) + " does not define any attributes")
			}

			if util.SliceContainsString(rule.ResourceTypes, resourceType) {
				for _, attrRule := range rule.RuleAttributes {
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
	if len(action.AttrRule.In) > 0 {
		l.validateIn(action)
	}
}

func (l *Linter) validateIn(action Action) {
	if !util.SliceContainsString(action.AttrRule.In, action.Diff.New) {
		violation := Violation{
			ResourceKey: action.ResourceKey,
			Attribute:   action.AttrName,
			Value:       action.Diff.New,
			Reason:      "Attribute " + action.ResourceKey + "." + action.AttrName + " contains an unacceptable value",
		}
		l.addViolation(violation)
	}
}

func (l *Linter) validateBeginsWith(action Action) {
	if !strings.HasPrefix(action.Diff.New, action.AttrRule.BeginsWith) {
		violation := Violation{
			ResourceKey: action.ResourceKey,
			Attribute:   action.AttrName,
			Value:       action.Diff.New,
			Reason:      "Attribute " + action.ResourceKey + "." + action.AttrName + " did not begin with " + action.AttrRule.BeginsWith,
		}
		l.addViolation(violation)
	}
}

func (l *Linter) addViolation(violation Violation) {
	l.Violations = append(l.Violations, violation)
}
