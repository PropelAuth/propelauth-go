package helpers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// Queryresponse is the common return type for the HTTP methods below. It structures the normal HTTP response in a way that's convient for us.
type QueryResponse struct {
	StatusCode   int
	ResponseText string
	BodyBytes    []byte
}

// interface for the QueryHelper
type QueryHelperInterface interface {
	Get(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error)
	Post(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error)
	Delete(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error)
}

type QueryHelper struct {
	urlPrefix           string
	backendUrlApiPrefix string
}

func NewQueryHelper(urlPrefix string, backendUrlApiPrefix string) *QueryHelper {
	return &QueryHelper{
		urlPrefix:           urlPrefix,
		backendUrlApiPrefix: backendUrlApiPrefix,
	}
}

// public http methods

func (o *QueryHelper) Get(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error) {
	url := o.assembleUrl(urlPostfix, queryParams)
	return o.RequestHelper("GET", token, url, nil)
}

func (o *QueryHelper) Post(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error) {
	url := o.assembleUrl(urlPostfix, queryParams)
	return o.RequestHelper("POST", token, url, bodyParams)
}

func (o *QueryHelper) Delete(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error) {
	url := o.assembleUrl(urlPostfix, queryParams)
	return o.RequestHelper("DELETE", token, url, nil)
}

// public helper method

func (o *QueryHelper) RequestHelper(method string, token string, url string, body []byte) (*QueryResponse, error) {
	requestBody := bytes.NewBuffer(body)

	// create request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error on creating request: %v", err)
	}

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error on response: %v", err)
	}
	defer resp.Body.Close()

	// convert the response body to a stream of bytes
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error on reading response body: %v", err)
	}

	respBytes := buf.Bytes()

	// return the response
	queryResponse := QueryResponse{
		StatusCode:   resp.StatusCode,
		ResponseText: resp.Status,
		BodyBytes:    respBytes,
	}

	return &queryResponse, nil
}

// private helper methods

func (o *QueryHelper) assembleUrl(urlPostfix string, queryParams url.Values) string {
	url := o.urlPrefix + o.backendUrlApiPrefix + urlPostfix
	if queryParams != nil {
		url += "?" + queryParams.Encode()
	}
	return url
}
