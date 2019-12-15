package util

import (
	"github.com/hashicorp/terraform/terraform"
	"os"
)

func OpenPlan(path string) (*terraform.Plan, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return nil, err
	}

	return plan, err
}
