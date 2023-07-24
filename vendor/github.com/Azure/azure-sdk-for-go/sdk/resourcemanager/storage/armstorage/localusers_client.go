//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstorage

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// LocalUsersClient contains the methods for the LocalUsers group.
// Don't use this type directly, use NewLocalUsersClient() instead.
type LocalUsersClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewLocalUsersClient creates a new instance of LocalUsersClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewLocalUsersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LocalUsersClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &LocalUsersClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update the properties of a local user associated with the storage account
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-01
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// accountName - The name of the storage account within the specified resource group. Storage account names must be between
// 3 and 24 characters in length and use numbers and lower-case letters only.
// username - The name of local user. The username must contain lowercase letters and numbers only. It must be unique only
// within the storage account.
// properties - The local user associated with a storage account.
// options - LocalUsersClientCreateOrUpdateOptions contains the optional parameters for the LocalUsersClient.CreateOrUpdate
// method.
func (client *LocalUsersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, username string, properties LocalUser, options *LocalUsersClientCreateOrUpdateOptions) (LocalUsersClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, accountName, username, properties, options)
	if err != nil {
		return LocalUsersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocalUsersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocalUsersClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *LocalUsersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, username string, properties LocalUser, options *LocalUsersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/localUsers/{username}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if username == "" {
		return nil, errors.New("parameter username cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{username}", url.PathEscape(username))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, properties)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *LocalUsersClient) createOrUpdateHandleResponse(resp *http.Response) (LocalUsersClientCreateOrUpdateResponse, error) {
	result := LocalUsersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalUser); err != nil {
		return LocalUsersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the local user associated with the specified storage account.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-01
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// accountName - The name of the storage account within the specified resource group. Storage account names must be between
// 3 and 24 characters in length and use numbers and lower-case letters only.
// username - The name of local user. The username must contain lowercase letters and numbers only. It must be unique only
// within the storage account.
// options - LocalUsersClientDeleteOptions contains the optional parameters for the LocalUsersClient.Delete method.
func (client *LocalUsersClient) Delete(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientDeleteOptions) (LocalUsersClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, username, options)
	if err != nil {
		return LocalUsersClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocalUsersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return LocalUsersClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return LocalUsersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *LocalUsersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/localUsers/{username}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if username == "" {
		return nil, errors.New("parameter username cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{username}", url.PathEscape(username))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the local user of the storage account by username.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-01
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// accountName - The name of the storage account within the specified resource group. Storage account names must be between
// 3 and 24 characters in length and use numbers and lower-case letters only.
// username - The name of local user. The username must contain lowercase letters and numbers only. It must be unique only
// within the storage account.
// options - LocalUsersClientGetOptions contains the optional parameters for the LocalUsersClient.Get method.
func (client *LocalUsersClient) Get(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientGetOptions) (LocalUsersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, username, options)
	if err != nil {
		return LocalUsersClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocalUsersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocalUsersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *LocalUsersClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/localUsers/{username}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if username == "" {
		return nil, errors.New("parameter username cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{username}", url.PathEscape(username))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *LocalUsersClient) getHandleResponse(resp *http.Response) (LocalUsersClientGetResponse, error) {
	result := LocalUsersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalUser); err != nil {
		return LocalUsersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List the local users associated with the storage account.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-01
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// accountName - The name of the storage account within the specified resource group. Storage account names must be between
// 3 and 24 characters in length and use numbers and lower-case letters only.
// options - LocalUsersClientListOptions contains the optional parameters for the LocalUsersClient.List method.
func (client *LocalUsersClient) NewListPager(resourceGroupName string, accountName string, options *LocalUsersClientListOptions) *runtime.Pager[LocalUsersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[LocalUsersClientListResponse]{
		More: func(page LocalUsersClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *LocalUsersClientListResponse) (LocalUsersClientListResponse, error) {
			req, err := client.listCreateRequest(ctx, resourceGroupName, accountName, options)
			if err != nil {
				return LocalUsersClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LocalUsersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LocalUsersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *LocalUsersClient) listCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *LocalUsersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/localUsers"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *LocalUsersClient) listHandleResponse(resp *http.Response) (LocalUsersClientListResponse, error) {
	result := LocalUsersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalUsers); err != nil {
		return LocalUsersClientListResponse{}, err
	}
	return result, nil
}

// ListKeys - List SSH authorized keys and shared key of the local user.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-01
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// accountName - The name of the storage account within the specified resource group. Storage account names must be between
// 3 and 24 characters in length and use numbers and lower-case letters only.
// username - The name of local user. The username must contain lowercase letters and numbers only. It must be unique only
// within the storage account.
// options - LocalUsersClientListKeysOptions contains the optional parameters for the LocalUsersClient.ListKeys method.
func (client *LocalUsersClient) ListKeys(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientListKeysOptions) (LocalUsersClientListKeysResponse, error) {
	req, err := client.listKeysCreateRequest(ctx, resourceGroupName, accountName, username, options)
	if err != nil {
		return LocalUsersClientListKeysResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocalUsersClientListKeysResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocalUsersClientListKeysResponse{}, runtime.NewResponseError(resp)
	}
	return client.listKeysHandleResponse(resp)
}

// listKeysCreateRequest creates the ListKeys request.
func (client *LocalUsersClient) listKeysCreateRequest(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientListKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/localUsers/{username}/listKeys"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if username == "" {
		return nil, errors.New("parameter username cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{username}", url.PathEscape(username))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listKeysHandleResponse handles the ListKeys response.
func (client *LocalUsersClient) listKeysHandleResponse(resp *http.Response) (LocalUsersClientListKeysResponse, error) {
	result := LocalUsersClientListKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalUserKeys); err != nil {
		return LocalUsersClientListKeysResponse{}, err
	}
	return result, nil
}

// RegeneratePassword - Regenerate the local user SSH password.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-01
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// accountName - The name of the storage account within the specified resource group. Storage account names must be between
// 3 and 24 characters in length and use numbers and lower-case letters only.
// username - The name of local user. The username must contain lowercase letters and numbers only. It must be unique only
// within the storage account.
// options - LocalUsersClientRegeneratePasswordOptions contains the optional parameters for the LocalUsersClient.RegeneratePassword
// method.
func (client *LocalUsersClient) RegeneratePassword(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientRegeneratePasswordOptions) (LocalUsersClientRegeneratePasswordResponse, error) {
	req, err := client.regeneratePasswordCreateRequest(ctx, resourceGroupName, accountName, username, options)
	if err != nil {
		return LocalUsersClientRegeneratePasswordResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocalUsersClientRegeneratePasswordResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocalUsersClientRegeneratePasswordResponse{}, runtime.NewResponseError(resp)
	}
	return client.regeneratePasswordHandleResponse(resp)
}

// regeneratePasswordCreateRequest creates the RegeneratePassword request.
func (client *LocalUsersClient) regeneratePasswordCreateRequest(ctx context.Context, resourceGroupName string, accountName string, username string, options *LocalUsersClientRegeneratePasswordOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/localUsers/{username}/regeneratePassword"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if username == "" {
		return nil, errors.New("parameter username cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{username}", url.PathEscape(username))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// regeneratePasswordHandleResponse handles the RegeneratePassword response.
func (client *LocalUsersClient) regeneratePasswordHandleResponse(resp *http.Response) (LocalUsersClientRegeneratePasswordResponse, error) {
	result := LocalUsersClientRegeneratePasswordResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalUserRegeneratePasswordResult); err != nil {
		return LocalUsersClientRegeneratePasswordResponse{}, err
	}
	return result, nil
}