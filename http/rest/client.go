// Package rest provides an HTTP client implementation using the Resty library
//
// The idea is to create an abstraction layer over the Resty client to make the codebase agnostic to the underlying HTTP client.
// This involves defining an interface for the HTTP client and a wrapper for the Resty client that implements this interface.
// This way, Resty can be easily swapped out for another HTTP client in the future without affecting the calling functions.
//
// Example usage:
//     var client rest.HTTPClient = rest.NewClient().
//         SetAuthToken("dummy-auth-token").
//         SetQueryParams(map[string]string{"key": "value"}).
//         AcceptJSON().
//         SetOutputDirectory("/path/to/output").
//         SetPathParams(map[string]string{"userId": "123", "accountId": "456"}).
//         EnableTrace()

//     resp, err := client.Get("/v1/users/{userId}/{accountId}/details")
//     if err != nil {
//     		log.Fatalf("Error: %v", err)
//     }

//     fmt.Println("Response Status:", resp.Status())
//     fmt.Println("Response Body:", string(resp.Body()))

//     traceInfo := resp.TraceInfo()
//     fmt.Printf("Trace Info: %+v\n", traceInfo)

package rest

// TODO: Query params passed through SetQueryParams() and SetQueryString() doesn't seem to work. Passing it via Get() works however.

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// HTTPClient interface defines the methods that an HTTP client should implement.
type HTTPClient interface {
	Get(url string) (*Response, error)
	SetQueryParams(params map[string]string) HTTPClient
	AcceptJSON() HTTPClient
	SetQueryString(query string) HTTPClient
	SetAuthToken(token string) HTTPClient
	SetOutput(filename string) HTTPClient
	SetOutputDirectory(dir string) HTTPClient
	SetPathParams(params map[string]string) HTTPClient
	EnableTrace() HTTPClient
	SetDefaults() HTTPClient
	SetDebug() HTTPClient
	SetBaseURL(url string) HTTPClient
	NewRequest()
}

// Response wraps the Resty response
type Response struct {
	restyResponse *resty.Response
}

// Status returns the HTTP status code
func (r *Response) Status() string {
	return r.restyResponse.Status()
}

// TraceInfo returns the trace information
func (r *Response) TraceInfo() resty.TraceInfo {
	return r.restyResponse.Request.TraceInfo()
}

// Body returns the response body
func (r *Response) Body() []byte {
	return r.restyResponse.Body()
}

// RestyClient struct implements HTTPClient interface
//
// RestyClient is a wrapper around Resty client and provides the necessary methods.
type RestyClient struct {
	restyClient  *resty.Client
	restyRequest *resty.Request
}

// NewClient creates a new RestyClient
func NewClient() *RestyClient {
	return &RestyClient{
		restyClient: resty.New(),
	}
}

func (c *RestyClient) Get(url string) (*Response, error) {
	resp, err := c.restyRequest.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("received error response: %s", resp.Status())
	}

	return &Response{restyResponse: resp}, nil
}

func (c *RestyClient) SetQueryParams(params map[string]string) HTTPClient {
	c.restyRequest.SetQueryParams(params)
	return c
}

func (c *RestyClient) AcceptJSON() HTTPClient {
	c.restyRequest.SetHeader("Accept", "application/json")
	return c
}

func (c *RestyClient) SetQueryString(query string) HTTPClient {
	c.restyRequest.SetQueryString(query)
	return c
}

func (c *RestyClient) SetAuthToken(token string) HTTPClient {
	c.restyRequest.SetAuthToken(token)
	return c
}

func (c *RestyClient) SetOutput(filename string) HTTPClient {
	c.restyRequest.SetOutput(filename)
	return c
}

func (c *RestyClient) SetOutputDirectory(dir string) HTTPClient {
	c.restyClient.SetOutputDirectory(dir)
	return c
}

func (c *RestyClient) SetPathParams(params map[string]string) HTTPClient {
	c.restyRequest.SetPathParams(params)
	return c
}

func (c *RestyClient) EnableTrace() HTTPClient {
	c.restyClient.EnableTrace()
	return c
}

func (c *RestyClient) SetDefaults() HTTPClient {
	c.restyClient.SetRetryCount(3)
	c.restyClient.SetRetryWaitTime(1 * time.Second)
	c.restyClient.SetTimeout(1 * time.Second)
	return c
}

func (c *RestyClient) SetDebug() HTTPClient {
	c.restyClient.SetDebug(true)
	return c
}

func (c *RestyClient) SetBaseURL(url string) HTTPClient {
	c.restyClient.SetBaseURL(url)
	return c
}

func (c *RestyClient) NewRequest() {
	c.restyRequest = c.restyClient.R()
}
