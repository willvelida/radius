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

package aws

// Tests that test with Mock RP functionality and UCP Server

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/radius-project/radius/pkg/armrpc/rpctest"
	"github.com/radius-project/radius/pkg/to"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_GetOperationResults(t *testing.T) {
	ucp, cloudcontrolClient, _ := initializeAWSTest(t)

	cloudcontrolClient.EXPECT().GetResourceRequestStatus(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *cloudcontrol.GetResourceRequestStatusInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceRequestStatusOutput, error) {
		output := cloudcontrol.GetResourceRequestStatusOutput{
			ProgressEvent: &types.ProgressEvent{
				RequestToken: to.Ptr(testAWSRequestToken),
			},
		}
		return &output, nil
	})

	operationResultsRequest, err := rpctest.NewHTTPRequestWithContent(context.Background(), http.MethodGet, ucp.BaseURL()+testProxyRequestAWSAsyncPath+"/operationResults/"+strings.ToLower(testAWSRequestToken), nil)
	require.NoError(t, err, "creating request failed")

	ctx := rpctest.NewARMRequestContext(operationResultsRequest)
	operationResultsRequest = operationResultsRequest.WithContext(ctx)

	operationResultsResponse, err := ucp.Client().Do(operationResultsRequest)
	require.NoError(t, err)

	assert.Equal(t, http.StatusAccepted, operationResultsResponse.StatusCode)
}
