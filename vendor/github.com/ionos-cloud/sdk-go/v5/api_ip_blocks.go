/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	_context "context"
	"fmt"
	_ioutil "io/ioutil"
	_nethttp "net/http"
	_neturl "net/url"
	"strings"
)

// Linger please
var (
	_ _context.Context
)

// IPBlocksApiService IPBlocksApi service
type IPBlocksApiService service

type ApiIpblocksDeleteRequest struct {
	ctx _context.Context
	ApiService *IPBlocksApiService
	ipblockId string
	pretty *bool
	depth *int32
	xContractNumber *int32
}

func (r ApiIpblocksDeleteRequest) Pretty(pretty bool) ApiIpblocksDeleteRequest {
	r.pretty = &pretty
	return r
}
func (r ApiIpblocksDeleteRequest) Depth(depth int32) ApiIpblocksDeleteRequest {
	r.depth = &depth
	return r
}
func (r ApiIpblocksDeleteRequest) XContractNumber(xContractNumber int32) ApiIpblocksDeleteRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiIpblocksDeleteRequest) Execute() (map[string]interface{}, *APIResponse, error) {
	return r.ApiService.IpblocksDeleteExecute(r)
}

/*
 * IpblocksDelete Delete IP Block
 * Removes the specific IP Block
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param ipblockId
 * @return ApiIpblocksDeleteRequest
 */
func (a *IPBlocksApiService) IpblocksDelete(ctx _context.Context, ipblockId string) ApiIpblocksDeleteRequest {
	return ApiIpblocksDeleteRequest{
		ApiService: a,
		ctx: ctx,
		ipblockId: ipblockId,
	}
}

/*
 * Execute executes the request
 * @return map[string]interface{}
 */
func (a *IPBlocksApiService) IpblocksDeleteExecute(r ApiIpblocksDeleteRequest) (map[string]interface{}, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodDelete
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPBlocksApiService.IpblocksDelete")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ipblocks/{ipblockId}"
	localVarPath = strings.Replace(localVarPath, "{"+"ipblockId"+"}", _neturl.PathEscape(parameterToString(r.ipblockId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse {
		Response: localVarHTTPResponse,
		Method: localVarHTTPMethod,
		RequestURL: localVarPath,
		RequestTime: httpRequestTime,
		Operation: "IpblocksDelete",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpblocksFindByIdRequest struct {
	ctx _context.Context
	ApiService *IPBlocksApiService
	ipblockId string
	pretty *bool
	depth *int32
	xContractNumber *int32
}

func (r ApiIpblocksFindByIdRequest) Pretty(pretty bool) ApiIpblocksFindByIdRequest {
	r.pretty = &pretty
	return r
}
func (r ApiIpblocksFindByIdRequest) Depth(depth int32) ApiIpblocksFindByIdRequest {
	r.depth = &depth
	return r
}
func (r ApiIpblocksFindByIdRequest) XContractNumber(xContractNumber int32) ApiIpblocksFindByIdRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiIpblocksFindByIdRequest) Execute() (IpBlock, *APIResponse, error) {
	return r.ApiService.IpblocksFindByIdExecute(r)
}

/*
 * IpblocksFindById Retrieve an IP Block
 * Retrieves the attributes of a given IP Block.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param ipblockId
 * @return ApiIpblocksFindByIdRequest
 */
func (a *IPBlocksApiService) IpblocksFindById(ctx _context.Context, ipblockId string) ApiIpblocksFindByIdRequest {
	return ApiIpblocksFindByIdRequest{
		ApiService: a,
		ctx: ctx,
		ipblockId: ipblockId,
	}
}

/*
 * Execute executes the request
 * @return IpBlock
 */
func (a *IPBlocksApiService) IpblocksFindByIdExecute(r ApiIpblocksFindByIdRequest) (IpBlock, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IpBlock
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPBlocksApiService.IpblocksFindById")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ipblocks/{ipblockId}"
	localVarPath = strings.Replace(localVarPath, "{"+"ipblockId"+"}", _neturl.PathEscape(parameterToString(r.ipblockId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse {
		Response: localVarHTTPResponse,
		Method: localVarHTTPMethod,
		RequestURL: localVarPath,
		RequestTime: httpRequestTime,
		Operation: "IpblocksFindById",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpblocksGetRequest struct {
	ctx _context.Context
	ApiService *IPBlocksApiService
	pretty *bool
	depth *int32
	xContractNumber *int32
	offset *int32
	limit *int32
}

func (r ApiIpblocksGetRequest) Pretty(pretty bool) ApiIpblocksGetRequest {
	r.pretty = &pretty
	return r
}
func (r ApiIpblocksGetRequest) Depth(depth int32) ApiIpblocksGetRequest {
	r.depth = &depth
	return r
}
func (r ApiIpblocksGetRequest) XContractNumber(xContractNumber int32) ApiIpblocksGetRequest {
	r.xContractNumber = &xContractNumber
	return r
}
func (r ApiIpblocksGetRequest) Offset(offset int32) ApiIpblocksGetRequest {
	r.offset = &offset
	return r
}
func (r ApiIpblocksGetRequest) Limit(limit int32) ApiIpblocksGetRequest {
	r.limit = &limit
	return r
}

func (r ApiIpblocksGetRequest) Execute() (IpBlocks, *APIResponse, error) {
	return r.ApiService.IpblocksGetExecute(r)
}

/*
 * IpblocksGet List IP Blocks 
 * Retrieve a list of all reserved IP Blocks
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @return ApiIpblocksGetRequest
 */
func (a *IPBlocksApiService) IpblocksGet(ctx _context.Context) ApiIpblocksGetRequest {
	return ApiIpblocksGetRequest{
		ApiService: a,
		ctx: ctx,
	}
}

/*
 * Execute executes the request
 * @return IpBlocks
 */
func (a *IPBlocksApiService) IpblocksGetExecute(r ApiIpblocksGetRequest) (IpBlocks, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IpBlocks
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPBlocksApiService.IpblocksGet")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ipblocks"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}
	if r.offset != nil {
		localVarQueryParams.Add("offset", parameterToString(*r.offset, ""))
	}
	if r.limit != nil {
		localVarQueryParams.Add("limit", parameterToString(*r.limit, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse {
		Response: localVarHTTPResponse,
		Method: localVarHTTPMethod,
		RequestURL: localVarPath,
		RequestTime: httpRequestTime,
		Operation: "IpblocksGet",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpblocksPatchRequest struct {
	ctx _context.Context
	ApiService *IPBlocksApiService
	ipblockId string
	ipblock *IpBlockProperties
	pretty *bool
	depth *int32
	xContractNumber *int32
}

func (r ApiIpblocksPatchRequest) Ipblock(ipblock IpBlockProperties) ApiIpblocksPatchRequest {
	r.ipblock = &ipblock
	return r
}
func (r ApiIpblocksPatchRequest) Pretty(pretty bool) ApiIpblocksPatchRequest {
	r.pretty = &pretty
	return r
}
func (r ApiIpblocksPatchRequest) Depth(depth int32) ApiIpblocksPatchRequest {
	r.depth = &depth
	return r
}
func (r ApiIpblocksPatchRequest) XContractNumber(xContractNumber int32) ApiIpblocksPatchRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiIpblocksPatchRequest) Execute() (IpBlock, *APIResponse, error) {
	return r.ApiService.IpblocksPatchExecute(r)
}

/*
 * IpblocksPatch Partially modify IP Block
 * You can use update attributes of a resource
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param ipblockId
 * @return ApiIpblocksPatchRequest
 */
func (a *IPBlocksApiService) IpblocksPatch(ctx _context.Context, ipblockId string) ApiIpblocksPatchRequest {
	return ApiIpblocksPatchRequest{
		ApiService: a,
		ctx: ctx,
		ipblockId: ipblockId,
	}
}

/*
 * Execute executes the request
 * @return IpBlock
 */
func (a *IPBlocksApiService) IpblocksPatchExecute(r ApiIpblocksPatchRequest) (IpBlock, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPatch
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IpBlock
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPBlocksApiService.IpblocksPatch")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ipblocks/{ipblockId}"
	localVarPath = strings.Replace(localVarPath, "{"+"ipblockId"+"}", _neturl.PathEscape(parameterToString(r.ipblockId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.ipblock == nil {
		return localVarReturnValue, nil, reportError("ipblock is required and must be specified")
	}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	// body params
	localVarPostBody = r.ipblock
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse {
		Response: localVarHTTPResponse,
		Method: localVarHTTPMethod,
		RequestURL: localVarPath,
		RequestTime: httpRequestTime,
		Operation: "IpblocksPatch",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpblocksPostRequest struct {
	ctx _context.Context
	ApiService *IPBlocksApiService
	ipblock *IpBlock
	pretty *bool
	depth *int32
	xContractNumber *int32
}

func (r ApiIpblocksPostRequest) Ipblock(ipblock IpBlock) ApiIpblocksPostRequest {
	r.ipblock = &ipblock
	return r
}
func (r ApiIpblocksPostRequest) Pretty(pretty bool) ApiIpblocksPostRequest {
	r.pretty = &pretty
	return r
}
func (r ApiIpblocksPostRequest) Depth(depth int32) ApiIpblocksPostRequest {
	r.depth = &depth
	return r
}
func (r ApiIpblocksPostRequest) XContractNumber(xContractNumber int32) ApiIpblocksPostRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiIpblocksPostRequest) Execute() (IpBlock, *APIResponse, error) {
	return r.ApiService.IpblocksPostExecute(r)
}

/*
 * IpblocksPost Reserve IP Block
 * This will reserve a new IP Block
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @return ApiIpblocksPostRequest
 */
func (a *IPBlocksApiService) IpblocksPost(ctx _context.Context) ApiIpblocksPostRequest {
	return ApiIpblocksPostRequest{
		ApiService: a,
		ctx: ctx,
	}
}

/*
 * Execute executes the request
 * @return IpBlock
 */
func (a *IPBlocksApiService) IpblocksPostExecute(r ApiIpblocksPostRequest) (IpBlock, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPost
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IpBlock
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPBlocksApiService.IpblocksPost")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ipblocks"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.ipblock == nil {
		return localVarReturnValue, nil, reportError("ipblock is required and must be specified")
	}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	// body params
	localVarPostBody = r.ipblock
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse {
		Response: localVarHTTPResponse,
		Method: localVarHTTPMethod,
		RequestURL: localVarPath,
		RequestTime: httpRequestTime,
		Operation: "IpblocksPost",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpblocksPutRequest struct {
	ctx _context.Context
	ApiService *IPBlocksApiService
	ipblockId string
	ipblock *IpBlock
	pretty *bool
	depth *int32
	xContractNumber *int32
}

func (r ApiIpblocksPutRequest) Ipblock(ipblock IpBlock) ApiIpblocksPutRequest {
	r.ipblock = &ipblock
	return r
}
func (r ApiIpblocksPutRequest) Pretty(pretty bool) ApiIpblocksPutRequest {
	r.pretty = &pretty
	return r
}
func (r ApiIpblocksPutRequest) Depth(depth int32) ApiIpblocksPutRequest {
	r.depth = &depth
	return r
}
func (r ApiIpblocksPutRequest) XContractNumber(xContractNumber int32) ApiIpblocksPutRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiIpblocksPutRequest) Execute() (IpBlock, *APIResponse, error) {
	return r.ApiService.IpblocksPutExecute(r)
}

/*
 * IpblocksPut Modify IP Block
 * You can use update attributes of a resource
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param ipblockId
 * @return ApiIpblocksPutRequest
 */
func (a *IPBlocksApiService) IpblocksPut(ctx _context.Context, ipblockId string) ApiIpblocksPutRequest {
	return ApiIpblocksPutRequest{
		ApiService: a,
		ctx: ctx,
		ipblockId: ipblockId,
	}
}

/*
 * Execute executes the request
 * @return IpBlock
 */
func (a *IPBlocksApiService) IpblocksPutExecute(r ApiIpblocksPutRequest) (IpBlock, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPut
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IpBlock
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPBlocksApiService.IpblocksPut")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ipblocks/{ipblockId}"
	localVarPath = strings.Replace(localVarPath, "{"+"ipblockId"+"}", _neturl.PathEscape(parameterToString(r.ipblockId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.ipblock == nil {
		return localVarReturnValue, nil, reportError("ipblock is required and must be specified")
	}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	// body params
	localVarPostBody = r.ipblock
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse {
		Response: localVarHTTPResponse,
		Method: localVarHTTPMethod,
		RequestURL: localVarPath,
		RequestTime: httpRequestTime,
		Operation: "IpblocksPut",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}
