/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	ctrl "github.com/radius-project/radius/pkg/armrpc/asyncoperation/controller"
	"github.com/radius-project/radius/pkg/components/database"
	deployment "github.com/radius-project/radius/pkg/corerp/backend/deployment"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteResourceRun_20231001Preview(t *testing.T) {

	setupTest := func() (func(tb testing.TB), *database.MockClient, *deployment.MockDeploymentProcessor, *ctrl.Request) {
		mctrl := gomock.NewController(t)

		msc := database.NewMockClient(mctrl)
		mdp := deployment.NewMockDeploymentProcessor(mctrl)

		req := &ctrl.Request{
			OperationID:   uuid.New(),
			OperationType: "APPLICATIONS.CORE/CONTAINERS|DELETE",
			ResourceID: fmt.Sprintf("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/radius-test-rg/providers/Applications.Core/containers/%s",
				uuid.NewString()),
			CorrelationID:    uuid.NewString(),
			OperationTimeout: &ctrl.DefaultAsyncOperationTimeout,
		}

		return func(tb testing.TB) {
			mctrl.Finish()
		}, msc, mdp, req
	}

	t.Parallel()

	deleteCases := []struct {
		desc     string
		getErr   error
		dpDelErr error
		scDelErr error
	}{
		{"delete-existing-resource", nil, nil, nil},
		{"delete-non-existing-resource", &database.ErrNotFound{}, nil, nil},
		{"delete-resource-dp-delete-error", nil, errors.New("deployment processor delete error"), nil},
		{"delete-resource-delete-from-db-error", nil, nil, errors.New("delete from db error")},
	}

	for _, tt := range deleteCases {
		t.Run(tt.desc, func(t *testing.T) {
			teardownTest, msc, mdp, req := setupTest()
			defer teardownTest(t)

			msc.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(&database.Object{}, tt.getErr).
				Times(1)

			if tt.getErr == nil {
				mdp.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.dpDelErr).
					Times(1)

				if tt.dpDelErr == nil {
					msc.EXPECT().
						Delete(gomock.Any(), gomock.Any()).
						Return(tt.scDelErr).
						Times(1)
				}
			}

			opts := ctrl.Options{
				DatabaseClient: msc,
				GetDeploymentProcessor: func() deployment.DeploymentProcessor {
					return mdp
				},
			}

			ctrl, err := NewDeleteResource(opts)
			require.NoError(t, err)

			_, err = ctrl.Run(context.Background(), req)

			if tt.getErr != nil || tt.dpDelErr != nil || tt.scDelErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
