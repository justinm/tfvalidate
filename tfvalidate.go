package main

import (
	"flag"
	"github.com/justinm/tfvalidate/tfvalidate"
	"log"
	"os"
)

func getCwd() string {
	return os.Getenv("PWD")
}

func main() {
	pathToPlan := flag.String("plan", "", "Path to the plan to lint")
	pathToRules := flag.String("lint", "", "Path to lint rules, defaults to $CWD/.tfvalidate.json")
	outputType := flag.String("output", "text", "Response type, options are 'text' or 'json'")

	flag.Parse()

	if len(*pathToPlan) == 0 {
		log.Fatal("--plan must be supplied")
	}

	if len(*pathToRules) == 0 {
		*pathToRules = getCwd() + "/.tfvalidate.json"
	}

	violations := tfvalidate.Validate(*pathToRules, *pathToPlan)

	switch *outputType {
	case "text":
		tfvalidate.PrintText(violations)
	case "json":
		tfvalidate.PrintJson(violations)
	default:
		log.Fatal("Unknown output type")
	}

	if len(violations) > 0 {
		log.Fatal("Violations were discovered")
	}
}
