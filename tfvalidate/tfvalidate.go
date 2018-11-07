package tfvalidate

import (
	"github.com/justinm/tfvalidate/tfvalidate/linter"
	"github.com/justinm/tfvalidate/tfvalidate/rules"
	"github.com/justinm/tfvalidate/tfvalidate/util"
	"log"
)

func Validate(pathToRules string, pathToPlan string) []linter.Violation {
	ruleset, err := rules.Parse(pathToRules)
	if err != nil {
		log.Fatal(err)
	}

	plan, err := util.OpenPlan(pathToPlan)
	if err != nil {
		log.Fatal("Failed to open plan at " + pathToPlan)
	}

	linter := linter.Linter{}

	linter.Lint(plan, ruleset)

	return linter.Violations
}
