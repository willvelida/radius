// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package deployment

import (
	"testing"

	"github.com/Azure/radius/pkg/healthcontract"
	"github.com/Azure/radius/pkg/model"
	"github.com/Azure/radius/pkg/model/components"
	"github.com/Azure/radius/pkg/workloads"
	"github.com/stretchr/testify/require"
)

func Test_DeploymentProcessor_OrderActions(t *testing.T) {
	// We're not going to render or deploy anything, so an empty model works
	model := model.NewModel(map[string]workloads.WorkloadRenderer{}, map[string]model.Handlers{})
	dp := deploymentProcessor{model, &healthcontract.HealthChannels{}}

	actions := map[string]ComponentAction{
		"A": {
			ComponentName: "A",
			Operation:     UpdateWorkload,
			Component: &components.GenericComponent{
				Uses: []components.GenericDependency{
					{
						Binding: components.NewComponentBindingExpression("myapp", "C", "test", ""),
					},
				},
			},
		},
		"B": {
			ComponentName: "B",
			Operation:     DeleteWorkload,
		},
		"C": {
			ComponentName: "C",
			Operation:     UpdateWorkload,
			Component:     &components.GenericComponent{},
		},
	}
	ordered, err := dp.orderActions(actions)
	require.NoError(t, err)

	expected := []ComponentAction{
		actions["C"],
		actions["A"],
		actions["B"],
	}

	require.Equal(t, expected, ordered)
}
