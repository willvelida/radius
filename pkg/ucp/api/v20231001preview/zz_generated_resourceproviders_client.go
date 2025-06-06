// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ResourceProvidersClient contains the methods for the ResourceProviders group.
// Don't use this type directly, use NewResourceProvidersClient() instead.
type ResourceProvidersClient struct {
	internal *arm.Client
}

// NewResourceProvidersClient creates a new instance of ResourceProvidersClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewResourceProvidersClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ResourceProvidersClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ResourceProvidersClient{
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a resource provider
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The plane name.
//   - resourceProviderName - The resource provider name. This is also the resource provider namespace. Example: 'Applications.Datastores'.
//   - resource - Resource create parameters.
//   - options - ResourceProvidersClientBeginCreateOrUpdateOptions contains the optional parameters for the ResourceProvidersClient.BeginCreateOrUpdate
//     method.
func (client *ResourceProvidersClient) BeginCreateOrUpdate(ctx context.Context, planeName string, resourceProviderName string, resource ResourceProviderResource, options *ResourceProvidersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ResourceProvidersClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, planeName, resourceProviderName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ResourceProvidersClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ResourceProvidersClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create or update a resource provider
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *ResourceProvidersClient) createOrUpdate(ctx context.Context, planeName string, resourceProviderName string, resource ResourceProviderResource, options *ResourceProvidersClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ResourceProvidersClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, planeName, resourceProviderName, resource, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ResourceProvidersClient) createOrUpdateCreateRequest(ctx context.Context, planeName string, resourceProviderName string, resource ResourceProviderResource, _ *ResourceProvidersClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/planes/radius/{planeName}/providers/System.Resources/resourceproviders/{resourceProviderName}"
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", url.PathEscape(planeName))
	if resourceProviderName == "" {
		return nil, errors.New("parameter resourceProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderName}", url.PathEscape(resourceProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
	return nil, err
}
;	return req, nil
}

// BeginDelete - Delete a resource provider
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The plane name.
//   - resourceProviderName - The resource provider name. This is also the resource provider namespace. Example: 'Applications.Datastores'.
//   - options - ResourceProvidersClientBeginDeleteOptions contains the optional parameters for the ResourceProvidersClient.BeginDelete
//     method.
func (client *ResourceProvidersClient) BeginDelete(ctx context.Context, planeName string, resourceProviderName string, options *ResourceProvidersClientBeginDeleteOptions) (*runtime.Poller[ResourceProvidersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, planeName, resourceProviderName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ResourceProvidersClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ResourceProvidersClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a resource provider
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *ResourceProvidersClient) deleteOperation(ctx context.Context, planeName string, resourceProviderName string, options *ResourceProvidersClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ResourceProvidersClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, planeName, resourceProviderName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ResourceProvidersClient) deleteCreateRequest(ctx context.Context, planeName string, resourceProviderName string, _ *ResourceProvidersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/planes/radius/{planeName}/providers/System.Resources/resourceproviders/{resourceProviderName}"
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", url.PathEscape(planeName))
	if resourceProviderName == "" {
		return nil, errors.New("parameter resourceProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderName}", url.PathEscape(resourceProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the specified resource provider.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The plane name.
//   - resourceProviderName - The resource provider name. This is also the resource provider namespace. Example: 'Applications.Datastores'.
//   - options - ResourceProvidersClientGetOptions contains the optional parameters for the ResourceProvidersClient.Get method.
func (client *ResourceProvidersClient) Get(ctx context.Context, planeName string, resourceProviderName string, options *ResourceProvidersClientGetOptions) (ResourceProvidersClientGetResponse, error) {
	var err error
	const operationName = "ResourceProvidersClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, planeName, resourceProviderName, options)
	if err != nil {
		return ResourceProvidersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ResourceProvidersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ResourceProvidersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ResourceProvidersClient) getCreateRequest(ctx context.Context, planeName string, resourceProviderName string, _ *ResourceProvidersClientGetOptions) (*policy.Request, error) {
	urlPath := "/planes/radius/{planeName}/providers/System.Resources/resourceproviders/{resourceProviderName}"
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", url.PathEscape(planeName))
	if resourceProviderName == "" {
		return nil, errors.New("parameter resourceProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderName}", url.PathEscape(resourceProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ResourceProvidersClient) getHandleResponse(resp *http.Response) (ResourceProvidersClientGetResponse, error) {
	result := ResourceProvidersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ResourceProviderResource); err != nil {
		return ResourceProvidersClientGetResponse{}, err
	}
	return result, nil
}

// GetProviderSummary - Get the specified resource provider summary. The resource provider summary aggregates the most commonly
// used information including locations, api versions and resource types.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The plane name.
//   - resourceProviderName - The resource provider name. This is also the resource provider namespace. Example: 'Applications.Datastores'.
//   - options - ResourceProvidersClientGetProviderSummaryOptions contains the optional parameters for the ResourceProvidersClient.GetProviderSummary
//     method.
func (client *ResourceProvidersClient) GetProviderSummary(ctx context.Context, planeName string, resourceProviderName string, options *ResourceProvidersClientGetProviderSummaryOptions) (ResourceProvidersClientGetProviderSummaryResponse, error) {
	var err error
	const operationName = "ResourceProvidersClient.GetProviderSummary"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getProviderSummaryCreateRequest(ctx, planeName, resourceProviderName, options)
	if err != nil {
		return ResourceProvidersClientGetProviderSummaryResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ResourceProvidersClientGetProviderSummaryResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ResourceProvidersClientGetProviderSummaryResponse{}, err
	}
	resp, err := client.getProviderSummaryHandleResponse(httpResp)
	return resp, err
}

// getProviderSummaryCreateRequest creates the GetProviderSummary request.
func (client *ResourceProvidersClient) getProviderSummaryCreateRequest(ctx context.Context, planeName string, resourceProviderName string, _ *ResourceProvidersClientGetProviderSummaryOptions) (*policy.Request, error) {
	urlPath := "/planes/radius/{planeName}/providers/{resourceProviderName}"
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", url.PathEscape(planeName))
	if resourceProviderName == "" {
		return nil, errors.New("parameter resourceProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderName}", url.PathEscape(resourceProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getProviderSummaryHandleResponse handles the GetProviderSummary response.
func (client *ResourceProvidersClient) getProviderSummaryHandleResponse(resp *http.Response) (ResourceProvidersClientGetProviderSummaryResponse, error) {
	result := ResourceProvidersClientGetProviderSummaryResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ResourceProviderSummary); err != nil {
		return ResourceProvidersClientGetProviderSummaryResponse{}, err
	}
	return result, nil
}

// NewListPager - List resource providers.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The plane name.
//   - options - ResourceProvidersClientListOptions contains the optional parameters for the ResourceProvidersClient.NewListPager
//     method.
func (client *ResourceProvidersClient) NewListPager(planeName string, options *ResourceProvidersClientListOptions) (*runtime.Pager[ResourceProvidersClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ResourceProvidersClientListResponse]{
		More: func(page ResourceProvidersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ResourceProvidersClientListResponse) (ResourceProvidersClientListResponse, error) {
		ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ResourceProvidersClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, planeName, options)
			}, nil)
			if err != nil {
				return ResourceProvidersClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
			},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *ResourceProvidersClient) listCreateRequest(ctx context.Context, planeName string, _ *ResourceProvidersClientListOptions) (*policy.Request, error) {
	urlPath := "/planes/radius/{planeName}/providers/System.Resources/resourceproviders"
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", url.PathEscape(planeName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ResourceProvidersClient) listHandleResponse(resp *http.Response) (ResourceProvidersClientListResponse, error) {
	result := ResourceProvidersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ResourceProviderResourceListResult); err != nil {
		return ResourceProvidersClientListResponse{}, err
	}
	return result, nil
}

// NewListProviderSummariesPager - List resource provider summaries. The resource provider summary aggregates the most commonly
// used information including locations, api versions and resource types.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The plane name.
//   - options - ResourceProvidersClientListProviderSummariesOptions contains the optional parameters for the ResourceProvidersClient.NewListProviderSummariesPager
//     method.
func (client *ResourceProvidersClient) NewListProviderSummariesPager(planeName string, options *ResourceProvidersClientListProviderSummariesOptions) (*runtime.Pager[ResourceProvidersClientListProviderSummariesResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ResourceProvidersClientListProviderSummariesResponse]{
		More: func(page ResourceProvidersClientListProviderSummariesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ResourceProvidersClientListProviderSummariesResponse) (ResourceProvidersClientListProviderSummariesResponse, error) {
		ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ResourceProvidersClient.NewListProviderSummariesPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listProviderSummariesCreateRequest(ctx, planeName, options)
			}, nil)
			if err != nil {
				return ResourceProvidersClientListProviderSummariesResponse{}, err
			}
			return client.listProviderSummariesHandleResponse(resp)
			},
		Tracer: client.internal.Tracer(),
	})
}

// listProviderSummariesCreateRequest creates the ListProviderSummaries request.
func (client *ResourceProvidersClient) listProviderSummariesCreateRequest(ctx context.Context, planeName string, _ *ResourceProvidersClientListProviderSummariesOptions) (*policy.Request, error) {
	urlPath := "/planes/radius/{planeName}/providers"
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", url.PathEscape(planeName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listProviderSummariesHandleResponse handles the ListProviderSummaries response.
func (client *ResourceProvidersClient) listProviderSummariesHandleResponse(resp *http.Response) (ResourceProvidersClientListProviderSummariesResponse, error) {
	result := ResourceProvidersClientListProviderSummariesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PagedResourceProviderSummary); err != nil {
		return ResourceProvidersClientListProviderSummariesResponse{}, err
	}
	return result, nil
}

