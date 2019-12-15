package main

import (
	"flag"
	"github.com/justinm/tfvalidate/approvers"
	"github.com/justinm/tfvalidate/linter"
	"github.com/justinm/tfvalidate/shared"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"os"
)

const (
	EXIT_ERR = 1
	EXIT_VALIDATION_ERR = 2
	EXIT_APPROVERS = 3
)

func getCwd() string {
	return os.Getenv("PWD")
}

var (
	logger  = logging.MustGetLogger("")

	configPath = flag.String("config", "", "Path to configuration, defaults to ~/.tfvalidate.yaml")
	outputType = flag.String("output", "text", "Response type, options are 'text' or 'json'")
	verbose    = flag.Bool("verbose", false, "Optional: verbose logging")
)

func main() {
	flag.Parse()

	SetupLogger()
	planPath := flag.Arg(0)

	if planPath == "" {
		println("Plan was not specified.")
		flag.Usage()
		os.Exit(EXIT_ERR)
	}

	configFile, err := GetConfigFile(); if err != nil {
		logger.Errorf("Cannot determine location to tfvalidate.yaml")
		os.Exit(EXIT_ERR)
	}

	plan, err := ReadPlan(planPath); if err != nil {
		logger.Errorf("Unable to read plan: %v", err)
		os.Exit(EXIT_ERR)
	}

	config, errs := shared.GetConfig(*configFile)
	if errs != nil {
		logger.Errorf("unable to load configuration: %v", errs)
		os.Exit(EXIT_ERR)
	}

	lint, errs := linter.New(config, plan); if errs != nil {
		logger.Errorf("unable to initialize linter: %v", errs)
		os.Exit(EXIT_ERR)
	}

	violations := lint.Lint()
	if len(violations) != 0 {
		PrintViolations(violations)
		os.Exit(EXIT_VALIDATION_ERR)
	}

	approve := approvers.GetApprovers(config, plan)
	if len(approve) > 0 {
		PrintApprovers(approve)
		os.Exit(EXIT_APPROVERS)
	}

	os.Exit(0)
}

func SetupLogger() {
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{id:03x}%{color:reset} %{message}`,
	)

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	if *verbose {
		logging.SetLevel(logging.DEBUG, "")
		backendLeveled.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
		backendLeveled.SetLevel(logging.INFO, "")
	}

	logger.SetBackend(backendLeveled)
}

func GetConfigFile() (*string, error) {
	if configPath == nil || *configPath == "" {
		tmpLoc, err := homedir.Expand("~/.tfvalidate.yaml")
		if err != nil {
			return nil, err
		}
		configPath = &tmpLoc
		logger.Debugf("Configuration: %s", *configPath)
	}

	return configPath, nil
}
