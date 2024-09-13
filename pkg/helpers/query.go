package helpers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
)

// Queryresponse is the common return type for the HTTP methods below. It structures the normal HTTP response
// in a way that's convient for us.
type QueryResponse struct {
	StatusCode   int
	ResponseText string
	BodyBytes    []byte
	BodyText     string
}

// Interface for the QueryHelper.
type QueryHelperInterface interface {
	Get(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error)
	Patch(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error)
	Post(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error)
	Put(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error)
	Delete(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error)
}

type QueryHelper struct {
	urlPrefix           string
	backendURLAPIPrefix string
}

func NewQueryHelper(urlPrefix string, backendURLAPIPrefix string) *QueryHelper {
	return &QueryHelper{
		urlPrefix:           urlPrefix,
		backendURLAPIPrefix: backendURLAPIPrefix,
	}
}

// public http methods

func (o *QueryHelper) Get(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error) {
	url := o.assembleURL(urlPostfix, queryParams)

	return o.RequestHelper("GET", token, url, nil)
}

func (o *QueryHelper) Patch(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error) {
	url := o.assembleURL(urlPostfix, queryParams)

	return o.RequestHelper("PATCH", token, url, bodyParams)
}

func (o *QueryHelper) Post(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error) {
	url := o.assembleURL(urlPostfix, queryParams)

	return o.RequestHelper("POST", token, url, bodyParams)
}

func (o *QueryHelper) Put(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error) {
	url := o.assembleURL(urlPostfix, queryParams)

	return o.RequestHelper("PUT", token, url, bodyParams)
}

func (o *QueryHelper) Delete(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error) {
	url := o.assembleURL(urlPostfix, queryParams)

	return o.RequestHelper("DELETE", token, url, bodyParams)
}

// public helper method

func (o *QueryHelper) RequestHelper(method string, token string, url string, body []byte) (*QueryResponse, error) {
	requestBody := bytes.NewBuffer(body)

	// create request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error on creating request: %w", err)
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "propelauth-go/0.8 go/"+runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error on response: %w", err)
	}
	defer resp.Body.Close()

	// convert the response body to a stream of bytes
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error on reading response body: %w", err)
	}

	respBytes := buf.Bytes()

	// return the response
	queryResponse := QueryResponse{
		StatusCode:   resp.StatusCode,
		ResponseText: resp.Status,
		BodyBytes:    respBytes,
		BodyText:     string(respBytes[:]),
	}

	return &queryResponse, nil
}

// private helper methods

func (o *QueryHelper) assembleURL(urlPostfix string, queryParams url.Values) string {
	url := o.urlPrefix + o.backendURLAPIPrefix + urlPostfix
	if queryParams != nil {
		url += "?" + queryParams.Encode()
	}

	return url
}
