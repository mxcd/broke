package clients

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

type AuthorizationType string

const (
	AuthorizationTypeApiKey AuthorizationType = "API_KEY"
	AuthorizationTypeBearer AuthorizationType = "BEARER"
)

type HttpRequestOptions struct {
	Method             HttpMethod
	ContextPath        string
	Body               interface{}
	ExpectedStatusCode int
}

type ApiClient interface {
	GetName() string
	GetBaseUrl() string
	GetAuthorizationType() AuthorizationType
	GetAuthorization() string
}

func DoHttpRequest(client ApiClient, options *HttpRequestOptions) (*http.Response, error) {
	method := string(options.Method)
	path := options.ContextPath
	data := options.Body
	expectedStatusCode := options.ExpectedStatusCode

	baseUrl := client.GetBaseUrl()
	authorizationType := client.GetAuthorizationType()
	authorization := client.GetAuthorization()

	requestUrl := baseUrl + path

	log.Trace().Msgf("HTTP %s > %s", method, requestUrl)

	var requestBody io.Reader

	if data != nil {
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal request body")
			return nil, err
		}
		requestBody = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequest(method, requestUrl, requestBody)
	if err != nil {
		log.Error().Err(err).Str("client", client.GetName()).Msgf("Failed to create http request")
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	if authorizationType == AuthorizationTypeApiKey {
		request.Header.Set("X-API-Key", authorization)
	} else if authorizationType == AuthorizationTypeBearer {
		request.Header.Set("Authorization", "Bearer "+authorization)
	}

	clientInstance := &http.Client{}
	response, err := clientInstance.Do(request)
	if err != nil {
		log.Error().Err(err).Str("client", client.GetName()).Msgf("Failed to send request")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != expectedStatusCode {
		log.Error().Str("client", client.GetName()).Msgf("Request failed with status %s. %d was expected", response.Status, expectedStatusCode)
		return nil, err
	}

	return response, nil
}

func DoHttpRequestWithResult[T any](client ApiClient, options *HttpRequestOptions, result *T) (*http.Response, error) {

	response, err := DoHttpRequest(client, options)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		log.Error().Err(err).Str("client", client.GetName()).Msgf("Failed to decode response body")
		return nil, err
	}

	return response, nil
}
