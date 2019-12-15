package main

import (
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/plans/planfile"
)

func ReadPlan(path string) (*plans.Plan, error) {
	reader, err := planfile.Open(path)
	if err != nil {
		return nil, err
	}

	plan, err := reader.ReadPlan(); if err != nil {
		return nil, err
	}

	return plan, err
}
