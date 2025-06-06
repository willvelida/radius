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

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/armrpc/frontend/controller"
	"github.com/radius-project/radius/pkg/armrpc/frontend/defaultoperation"
	"github.com/radius-project/radius/pkg/armrpc/frontend/server"
	aztoken "github.com/radius-project/radius/pkg/azure/tokencredentials"
	"github.com/radius-project/radius/pkg/retry"
	"github.com/radius-project/radius/pkg/ucp"
	"github.com/radius-project/radius/pkg/ucp/api/v20231001preview"
	ucp_aws "github.com/radius-project/radius/pkg/ucp/aws"
	sdk_cred "github.com/radius-project/radius/pkg/ucp/credentials"
	"github.com/radius-project/radius/pkg/ucp/datamodel"
	"github.com/radius-project/radius/pkg/ucp/datamodel/converter"
	awsproxy_ctrl "github.com/radius-project/radius/pkg/ucp/frontend/controller/awsproxy"
	aws_credential_ctrl "github.com/radius-project/radius/pkg/ucp/frontend/controller/credentials/aws"
	planes_ctrl "github.com/radius-project/radius/pkg/ucp/frontend/controller/planes"
	"github.com/radius-project/radius/pkg/ucp/ucplog"
	"github.com/radius-project/radius/pkg/validator"
)

const (
	planeCollectionPath = "/planes/aws"
	planeResourcePath   = "/planes/aws/{planeName}"

	resourceCollectionPath = planeResourcePath + "/accounts/{accountId}/regions/{region}/providers/{providerNamespace}/{resourceType}"
	operationResultsPath   = planeResourcePath + "/accounts/{accountId}/regions/{region}/providers/{providerNamespace}/locations/{location}/operationResults/{operationId}"
	operationStatusesPath  = planeResourcePath + "/accounts/{accountId}/regions/{region}/providers/{providerNamespace}/locations/{location}/operationStatuses/{operationId}"

	credentialResourcePath   = planeResourcePath + "/providers/System.AWS/credentials/{credentialName}"
	credentialCollectionPath = planeResourcePath + "/providers/System.AWS/credentials"

	// OperationTypeAWSResource is the operation type for CRUDL operations on AWS resources.
	OperationTypeAWSResource = "AWSRESOURCE"
)

// Initialize initializes the AWS module.
func (m *Module) Initialize(ctx context.Context) (http.Handler, error) {
	secretClient, err := m.options.SecretProvider.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	// Support override of AWS Clients for testing.
	if m.AWSClients.CloudControl == nil || m.AWSClients.CloudFormation == nil {
		awsConfig, err := m.newAWSConfig(ctx)
		if err != nil {
			return nil, err
		}

		if m.AWSClients.CloudControl == nil {
			m.AWSClients.CloudControl = cloudcontrol.NewFromConfig(awsConfig)
		}

		if m.AWSClients.CloudFormation == nil {
			m.AWSClients.CloudFormation = cloudformation.NewFromConfig(awsConfig)
		}
	}

	baseRouter := server.NewSubrouter(m.router, m.options.Config.Server.PathBase+"/")

	apiValidator := validator.APIValidator(validator.Options{
		SpecLoader:         m.options.SpecLoader,
		ResourceTypeGetter: validator.UCPResourceTypeGetter,
	})

	planeResourceOptions := controller.ResourceOptions[datamodel.AWSPlane]{
		RequestConverter:  converter.AWSPlaneDataModelFromVersioned,
		ResponseConverter: converter.AWSPlaneDataModelToVersioned,
	}

	// URLs for lifecycle of planes
	planeCollectionRouter := server.NewSubrouter(baseRouter, planeCollectionPath, apiValidator)
	planeResourceRouter := server.NewSubrouter(baseRouter, planeResourcePath, apiValidator)

	handlerOptions := []server.HandlerOptions{
		{
			// This is a scope query so we can't use the default operation.
			ParentRouter:  planeCollectionRouter,
			Method:        v1.OperationList,
			ResourceType:  datamodel.AWSPlaneResourceType,
			OperationType: &v1.OperationType{Type: datamodel.AWSPlaneResourceType, Method: v1.OperationList},
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return &planes_ctrl.ListPlanesByType[*datamodel.AWSPlane, datamodel.AWSPlane]{
					Operation: controller.NewOperation(opts, planeResourceOptions),
				}, nil
			},
		},
		{
			ParentRouter:  planeResourceRouter,
			Method:        v1.OperationGet,
			ResourceType:  datamodel.AWSPlaneResourceType,
			OperationType: &v1.OperationType{Type: datamodel.AWSPlaneResourceType, Method: v1.OperationGet},
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return defaultoperation.NewGetResource(opts, planeResourceOptions)
			},
		},
		{
			ParentRouter:  planeResourceRouter,
			Method:        v1.OperationPut,
			ResourceType:  datamodel.AWSPlaneResourceType,
			OperationType: &v1.OperationType{Type: datamodel.AWSPlaneResourceType, Method: v1.OperationPut},
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return defaultoperation.NewDefaultSyncPut(opts, planeResourceOptions)
			},
		},
		{
			ParentRouter:  planeResourceRouter,
			Method:        v1.OperationDelete,
			ResourceType:  datamodel.AWSPlaneResourceType,
			OperationType: &v1.OperationType{Type: datamodel.AWSPlaneResourceType, Method: v1.OperationDelete},
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return defaultoperation.NewDefaultSyncDelete(opts, planeResourceOptions)
			},
		},
		{
			// URLs for standard UCP resource async status result.
			ParentRouter:  server.NewSubrouter(baseRouter, operationResultsPath),
			Method:        v1.OperationGet,
			OperationType: &v1.OperationType{Type: OperationResultsResourceType, Method: v1.OperationGet},
			ResourceType:  OperationResultsResourceType,
			ControllerFactory: func(opt controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewGetAWSOperationResults(opt, m.AWSClients)
			},
		},
		{
			// URLs for standard UCP resource async status.
			ParentRouter:  server.NewSubrouter(baseRouter, operationStatusesPath),
			Method:        v1.OperationGet,
			OperationType: &v1.OperationType{Type: OperationStatusResourceType, Method: v1.OperationGet},
			ResourceType:  OperationStatusResourceType,
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewGetAWSOperationStatuses(opts, m.AWSClients)
			},
		},
	}

	resourceCollectionRouter := server.NewSubrouter(baseRouter, resourceCollectionPath)
	handlerOptions = append(handlerOptions, []server.HandlerOptions{
		{
			// URLs for standard UCP resource lifecycle operations.
			ParentRouter:  resourceCollectionRouter,
			Method:        v1.OperationList,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationList},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewListAWSResources(opts, m.AWSClients)
			},
		},
		{
			ParentRouter:  resourceCollectionRouter,
			Path:          "/{resourceName}",
			Method:        v1.OperationPut,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationPut},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewCreateOrUpdateAWSResource(opts, m.AWSClients)
			},
		},
		{
			ParentRouter:  resourceCollectionRouter,
			Path:          "/{resourceName}",
			Method:        v1.OperationDelete,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationDelete},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewDeleteAWSResource(opts, m.AWSClients)
			},
		},
		{
			ParentRouter:  resourceCollectionRouter,
			Path:          "/{resourceName}",
			Method:        v1.OperationGet,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationGet},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opt controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewGetAWSResource(opt, m.AWSClients)
			},
		},
	}...)

	// URLs for "non-idempotent" resource lifecycle operations. These are extensions to the UCP spec that are needed when
	// a resource has a non-idempotent lifecyle and a computed name.
	//
	// The normal UCP lifecycle operations have a user-specified resource name which must be part of the URL. These
	// operations are structured so that the resource name is not part of the URL.
	handlerOptions = append(handlerOptions, []server.HandlerOptions{
		{
			ParentRouter:  resourceCollectionRouter,
			Path:          "/:put",
			Method:        v1.OperationPutImperative,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationPutImperative},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opt controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewCreateOrUpdateAWSResourceWithPost(opt, m.AWSClients)
			},
		},
		{
			ParentRouter:  resourceCollectionRouter,
			Path:          "/:get",
			Method:        v1.OperationGetImperative,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationGetImperative},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opt controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewGetAWSResourceWithPost(opt, m.AWSClients, retry.NewDefaultRetryer())
			},
		},
		{
			ParentRouter:  resourceCollectionRouter,
			Path:          "/:delete",
			Method:        v1.OperationDeleteImperative,
			OperationType: &v1.OperationType{Type: OperationTypeAWSResource, Method: v1.OperationDeleteImperative},
			ResourceType:  OperationTypeAWSResource,
			ControllerFactory: func(opts controller.Options) (controller.Controller, error) {
				return awsproxy_ctrl.NewDeleteAWSResourceWithPost(opts, m.AWSClients)
			},
		},
	}...)

	// URLs for operations on AWS credential resources.
	//
	// These use the OpenAPI spec validator. General AWS operations DO NOT use the spec validator
	// because we rely on CloudControl's validation.
	credentialCollectionRouter := server.NewSubrouter(baseRouter, credentialCollectionPath, apiValidator)
	credentialResourceRouter := server.NewSubrouter(baseRouter, credentialResourcePath, apiValidator)

	handlerOptions = append(handlerOptions, []server.HandlerOptions{
		{
			ParentRouter: credentialCollectionRouter,
			ResourceType: v20231001preview.AWSCredentialType,
			Method:       v1.OperationList,
			ControllerFactory: func(opt controller.Options) (controller.Controller, error) {
				return defaultoperation.NewListResources(opt,
					controller.ResourceOptions[datamodel.AWSCredential]{
						RequestConverter:  converter.AWSCredentialDataModelFromVersioned,
						ResponseConverter: converter.AWSCredentialDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: credentialResourceRouter,
			ResourceType: v20231001preview.AWSCredentialType,
			Method:       v1.OperationGet,
			ControllerFactory: func(opt controller.Options) (controller.Controller, error) {
				return defaultoperation.NewGetResource(opt,
					controller.ResourceOptions[datamodel.AWSCredential]{
						RequestConverter:  converter.AWSCredentialDataModelFromVersioned,
						ResponseConverter: converter.AWSCredentialDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: credentialResourceRouter,
			Method:       v1.OperationPut,
			ResourceType: v20231001preview.AWSCredentialType,
			ControllerFactory: func(o controller.Options) (controller.Controller, error) {
				return aws_credential_ctrl.NewCreateOrUpdateAWSCredential(o, secretClient)
			},
		},
		{
			ParentRouter: credentialResourceRouter,
			Method:       v1.OperationDelete,
			ResourceType: v20231001preview.AWSCredentialType,
			ControllerFactory: func(o controller.Options) (controller.Controller, error) {
				return aws_credential_ctrl.NewDeleteAWSCredential(o, secretClient)
			},
		},
	}...)

	databaseClient, err := m.options.DatabaseProvider.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	ctrlOpts := controller.Options{
		Address:        m.options.Config.Server.Address(),
		DatabaseClient: databaseClient,
		PathBase:       m.options.Config.Server.PathBase,
		StatusManager:  m.options.StatusManager,

		KubeClient:   nil, // Unused by AWS module
		ResourceType: "",  // Set dynamically
	}

	for _, h := range handlerOptions {
		if err := server.RegisterHandler(ctx, h, ctrlOpts); err != nil {
			return nil, err
		}
	}

	return m.router, nil
}

func (m *Module) newAWSConfig(ctx context.Context) (aws.Config, error) {
	logger := ucplog.FromContextOrDiscard(ctx)
	credProviders := []func(*config.LoadOptions) error{}

	switch m.options.Config.Identity.AuthMethod {
	case ucp.AuthUCPCredential:
		provider, err := sdk_cred.NewAWSCredentialProvider(m.options.SecretProvider, m.options.UCP, &aztoken.AnonymousCredential{})
		if err != nil {
			return aws.Config{}, err
		}
		p := ucp_aws.NewUCPCredentialProvider(provider, ucp_aws.DefaultExpireDuration)
		credProviders = append(credProviders, config.WithCredentialsProvider(p))
		logger.Info("Configuring 'UCPCredential' authentication mode using UCP Credential API")

	default:
		logger.Info("Configuring default authentication mode with environment variable.")
	}

	awscfg, err := config.LoadDefaultConfig(ctx, credProviders...)
	if err != nil {
		return aws.Config{}, err
	}

	return awscfg, nil
}
