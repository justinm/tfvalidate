package approvers

import (
	"github.com/hashicorp/terraform/plans"
	"github.com/justinm/tfvalidate/shared"
)

func GetApprovers(config *shared.Configuration, plan *plans.Plan) []string {
	approvers := make(map[string]bool)

	if config.Approvals == nil {
		return nil
	}

	for _, resource := range plan.Changes.Resources {
		for _, approval := range config.Approvals {
			for _, resourceName := range approval.Resources {
				t := resource.Addr.Resource.Resource.Type

				if t == resourceName {
					for _, approver := range approval.Approvers {
						approvers[approver] = true
					}
				}
			}
		}
	}

	var keys []string
	for approver := range approvers {
		keys = append(keys, approver)
	}

	return keys
}