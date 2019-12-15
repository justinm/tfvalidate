package linter

import (
	"fmt"
	"github.com/hashicorp/terraform/plans"
	"github.com/justinm/tfvalidate/shared"
	"github.com/justinm/tfvalidate/util"
	"github.com/op/go-logging"
	"github.com/zclconf/go-cty/cty"
	"strings"
)

var (
	logger = logging.MustGetLogger("linter")
)

type Violation struct {
	Reason    string
	Change    *plans.ResourceInstanceChange
	Attribute shared.RuleAttributeDefinition
	Value     *cty.Value
}

type Linter struct {
	Config     *shared.Configuration
	Plan       *plans.Plan
	Violations []Violation
}

func New(config *shared.Configuration, plan *plans.Plan) (*Linter, []error) {
	linter := Linter{
		Config: config,
		Plan:   plan,
	}

	return &linter, nil
}

func (linter *Linter) Lint() []Violation {
	var resources []plans.ResourceInstanceChangeSrc
	for _, resourceChange := range linter.Plan.Changes.Resources {
		resources = append(resources, *resourceChange)
	}

	if linter.Config.Rules == nil {
		logger.Warning("no lint rules were found")
		return nil
	}

	for _, resource := range resources {
		_ = linter.LintChange(resource)
	}
	return linter.Violations
}

func (linter *Linter) getResourceFromKey(key string) string {
	parts := strings.Split(key, ".")

	return parts[0]
}

func (linter *Linter) LintChange(change plans.ResourceInstanceChangeSrc) error {
	resourceName := change.Addr.String()
	resourceType := change.Addr.Resource.Resource.Type

	validActions := []interface{}{
		plans.Update,
		plans.Create,
		plans.DeleteThenCreate,
		plans.CreateThenDelete,
	}

	logger.Debugf("Resource will %s: %s (%s)", change.Action.String(), resourceName, resourceType)

	if !util.SliceContains(validActions, change.Action) {
		logger.Debugf("Resource %s skipped", resourceName)
		return nil
	}

	t, err := change.After.ImpliedType(); if err != nil {
		return err
	}
	diff, err := change.Decode(t)
	if err != nil {
		logger.Errorf("Failure in decoding plan: %v", err)
		return err
	}

	for _, rule := range linter.Config.Rules {
		if util.SliceStringContains(rule.ResourceTypes, change.Addr.Resource.Resource.Type) {
			for _, attribute := range rule.RuleAttributes {
				for _, subrule := range attribute.Rules {
					if subrule.Required != nil && *subrule.Required {
						violation := linter.validateRequired(attribute, diff)
						if violation != nil {
							linter.addViolation(*violation)
						}
					}

					if subrule.StartsWith != nil {
						linter.validateStartsWith(attribute, subrule, diff)
					}

					if subrule.OneOf != nil {
						linter.validateOneOf(attribute, subrule, diff)
					}

					if subrule.Contains != nil {
						violation := linter.validateContains(attribute, subrule, diff); if violation != nil {
							linter.addViolation(*violation)
						}
					}
				}
			}
		}
	}

	return nil
}

func (linter *Linter) getAttributeValueFromDiff(name string, currentVal cty.Value) *cty.Value {
	keys := strings.Split(name, ".")

	for _, key := range keys {
		if currentVal.Type().IsMapType() || currentVal.Type().IsObjectType() {
			values := currentVal.AsValueMap()

			val, exists := values[key]
			if exists {
				currentVal = val
				continue
			} else {
				return nil
			}
		}

		logger.Errorf("unknown type %s", currentVal.Type().FriendlyName())
		return nil
	}

	return &currentVal
}

func (linter *Linter) validateRequired(attr shared.RuleAttributeDefinition, diff *plans.ResourceInstanceChange) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff.After)

	if val == nil {
		return &Violation{
			Attribute:   attr,
			Change: 	 diff,
			Reason:      "Attribute is required but undefined",
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) validateContains(attr shared.RuleAttributeDefinition, rule shared.RuleAttributeDefinitionRule, diff *plans.ResourceInstanceChange) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff.After); if val == nil {
		return nil
	}

	if val.Type().IsListType() {
		vals := val.AsValueSlice()
		var ivals []interface{}
		for _, val := range vals {
			if val.Type().IsPrimitiveType() {
				ivals = append(ivals, val.AsString())
			}
		}

		if !util.SliceContains(ivals, rule.Contains) {
			return &Violation{
				Attribute:   attr,
				Change:      diff,
				Reason:      fmt.Sprintf("attribute does not contain %s", *rule.Contains),
				Value:       val,
			}
		}
	} else if val.Type().IsObjectType() || val.Type().IsMapType() {
		vals := val.AsValueMap()
		var keys []interface{}

		for key := range vals {
			keys = append(keys, key)
		}

		if !util.SliceContains(keys, rule.Contains) {
			return &Violation{
				Attribute:   attr,
				Change:		 diff,
				Reason:      fmt.Sprintf("attribute does not contain %s", *rule.Contains),
				Value:       val,
			}
		}
	} else {
		return &Violation{
			Attribute:   attr,
			Change:		 diff,
			Reason:      fmt.Sprintf("attribute type %s is not compatible with contains", val.Type().FriendlyName()),
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) validateOneOf(attr shared.RuleAttributeDefinition, rule shared.RuleAttributeDefinitionRule, diff *plans.ResourceInstanceChange) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff.After); if val == nil {
		return nil
	}

	if val.Type().IsPrimitiveType() {
		if !util.SliceStringContains(rule.OneOf, val.AsString()) {
			return &Violation{
				Attribute:   attr,
				Change:		 diff,
				Reason:      fmt.Sprintf("attribute does not contain one of %v", rule.OneOf),
				Value:       val,
			}
		}
	} else {
		return &Violation{
			Attribute:   attr,
			Change:		 diff,
			Reason:      fmt.Sprintf("attribute type %s is not compatible with oneOf", val.Type().FriendlyName()),
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) validateStartsWith(attr shared.RuleAttributeDefinition, rule shared.RuleAttributeDefinitionRule, diff *plans.ResourceInstanceChange) *Violation {
	val := linter.getAttributeValueFromDiff(attr.Name, diff.After); if val == nil {
		return nil
	}

	if !strings.HasPrefix(val.AsString(), *rule.StartsWith) {
		return &Violation{
			Attribute:   attr,
			Change:		 diff,
			Reason:      "attribute does not start with " + *rule.StartsWith,
			Value:       val,
		}
	}

	return nil
}

func (linter *Linter) addViolation(violation Violation) {
	linter.Violations = append(linter.Violations, violation)
}
