package tfvalidate

import (
	"encoding/json"
	"fmt"
	"github.com/justinm/tfvalidate/tfvalidate/linter"
	"os"
)

type JsonResponse struct {
	Violations []linter.Violation `json:"violations"`
}

func PrintText(violations []linter.Violation) {
	if len(violations) > 0 {
		for _, violation := range violations {
			fmt.Printf("Violation: %s\n", violation.Reason)
		}
		os.Exit(1)
	} else {
		fmt.Print("No violations were found")
	}
}

func PrintJson(violations []linter.Violation) {
	jsonResponse := JsonResponse{Violations: violations}
	data, err := json.Marshal(jsonResponse)
	if err != nil {
		fmt.Printf("Error in printing JSON response: %T\n", err)
		os.Exit(1)
	}

	fmt.Print(string(data))

	if len(violations) > 0 {
		os.Exit(1)
	}
}
