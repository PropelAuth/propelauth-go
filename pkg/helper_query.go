package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

type QueryHelper interface {
	Get(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error)
	Post(token string, urlPostfix string, queryParams url.Values, body []byte) (*QueryResponse, error)
	Delete(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error)
}

type queryHelper struct {
	urlPrefix           string
	backendUrlApiPrefix string
}

// initialize query helper
func NewQueryHelper(urlPrefix string, backendUrlApiPrefix string) QueryHelper {
	return queryHelper{
		urlPrefix:           urlPrefix,
		backendUrlApiPrefix: backendUrlApiPrefix,
	}
}

// public http methods

type QueryResponse struct {
	StatusCode   int
	ResponseText string
	BodyBytes    []byte
}

func (o queryHelper) Get(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error) {
	url := o.assembleUrl(urlPostfix, queryParams)
	return o.requestHelper("GET", token, url, nil)
}

func (o queryHelper) Post(token string, urlPostfix string, queryParams url.Values, bodyParams []byte) (*QueryResponse, error) {
	url := o.assembleUrl(urlPostfix, queryParams)
	return o.requestHelper("POST", token, url, bodyParams)
}

func (o queryHelper) Delete(token string, urlPostfix string, queryParams url.Values) (*QueryResponse, error) {
	url := o.assembleUrl(urlPostfix, queryParams)
	return o.requestHelper("DELETE", token, url, nil)
}

// private helper methods

func (o *queryHelper) assembleUrl(urlPostfix string, queryParams url.Values) string {
	url := o.urlPrefix + o.backendUrlApiPrefix + urlPostfix
	if queryParams != nil {
		url += "?" + queryParams.Encode()
	}
	return url
}

func (o *queryHelper) requestHelper(method string, token string, url string, body []byte) (*QueryResponse, error) {

	requestBody := bytes.NewBuffer(body)

	// create request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error on creating request: %v", err)
	}

	// add authorization
	var bearer = "Bearer " + token
	req.Header.Add("Authorization", bearer)

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error on response: %v", err)
	}
	defer resp.Body.Close()

	// convert the response body to a stream of bytes
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBytes := buf.Bytes()

	// return the response
	queryResponse := QueryResponse{
		StatusCode:   resp.StatusCode,
		ResponseText: resp.Status,
		BodyBytes:    respBytes,
	}

	return &queryResponse, nil
}
