package clients

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type MailcowClient struct {
	Options *MailcowClientOptions
}

type MailcowClientOptions struct {
	Name   string
	Url    string
	ApiKey string
}

func NewMailcowClient(options *MailcowClientOptions) (*MailcowClient, error) {
	options.Url = strings.TrimRight(options.Url, "/")
	return &MailcowClient{
		Options: options,
	}, nil
}

func (c MailcowClient) GetName() string {
	return c.Options.Name
}
func (c MailcowClient) GetBaseUrl() string {
	return c.Options.Url
}
func (c MailcowClient) GetAuthorizationType() AuthorizationType {
	return AuthorizationTypeApiKey
}
func (c MailcowClient) GetAuthorization() string {
	return c.Options.ApiKey
}

func (c *MailcowClient) TestConnection() error {
	log.Debug().Str("client", c.Options.Name).Msgf("Testing connection to Mailcow API at '%s'", c.Options.Url)

	_, err := DoHttpRequest(*c, &HttpRequestOptions{
		Method:             GET,
		ContextPath:        "/api/v1/get/status/containers",
		ExpectedStatusCode: 200,
	})

	if err != nil {
		log.Error().Err(err).Str("client", c.Options.Name).Msgf("Failed to test Mailcow API connection for user target '%s'", c.Options.Name)
		return err
	}

	log.Debug().Str("client", c.Options.Name).Msgf("Successfully connected to Mailcow API at '%s'", c.Options.Url)
	return nil
}
