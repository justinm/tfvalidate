package main

import (
	"encoding/json"
	"fmt"
	"github.com/justinm/tfvalidate/linter"
	"os"
)

type JsonResponse struct {
	Violations []JsonViolation `json:"violations,omitempty"`
	Approvers  []string	`json:"approvers,omitempty"`
}

type JsonViolation struct {
	ResourceName string  `json:"resource_name"`
	ModuleName   *string `json:"module_name"`
	Reason       string  `json:"reason"`
}

func PrintViolations(violations []linter.Violation) {
	var jsonViolations []JsonViolation

	for _, violation := range violations {
		var moduleName string

		if violation.Change.Addr.Module != nil {
			moduleName = violation.Change.Addr.Module.String()
		}

		jsonViolations = append(jsonViolations, JsonViolation{
			ResourceName: violation.Change.Addr.Resource.String(),
			ModuleName:   &moduleName,
			Reason:       violation.Reason,
		})
	}

	jsonResponse := JsonResponse{Violations: jsonViolations}
	data, err := json.Marshal(jsonResponse)
	if err != nil {
		fmt.Printf("Error in printing JSON response: %T\n", err)
		os.Exit(1)
	}

	fmt.Print(string(data))
}

func PrintApprovers(approvers []string) {
	jsonResponse := JsonResponse{Approvers: approvers}
	data, err := json.Marshal(jsonResponse)
	if err != nil {
		fmt.Printf("Error in printing JSON response: %T\n", err)
		os.Exit(1)
	}

	fmt.Print(string(data))
}