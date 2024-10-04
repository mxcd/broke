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

type MailcowMailboxResult struct {
	Active    int    `json:"active"`
	Username  string `json:"username"`
	Domain    string `json:"domain"`
	LocalPart string `json:"local_part"`
	Name      string `json:"name"`
}

func (c *MailcowClient) MailboxExists(email string) (bool, error) {
	log.Debug().Str("client", c.Options.Name).Msgf("Checking if mailbox '%s' exists", email)

	mailboxResult := &MailcowMailboxResult{}
	_, err := DoHttpRequestWithResult[MailcowMailboxResult](*c, &HttpRequestOptions{
		Method:             GET,
		ContextPath:        "/api/v1/get/mailbox/" + email,
		ExpectedStatusCode: 200,
	}, mailboxResult)
	if err != nil {
		log.Error().Err(err).Str("client", c.Options.Name).Msgf("Failed to check if mailbox '%s' exists", email)
		return false, err
	}

	if mailboxResult.Username == email {
		log.Debug().Str("client", c.Options.Name).Msgf("Mailbox '%s' exists", email)
		return true, nil
	} else {
		log.Debug().Str("client", c.Options.Name).Msgf("Mailbox '%s' does not exist", email)
		return false, nil
	}
}

type CreateMailboxOptions struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	LocalPart  string `json:"local_part"`
	AuthSource string `json:"authsource"`
}

func (c *MailcowClient) CreateMailbox(options *CreateMailboxOptions) error {
	log.Debug().Str("client", c.Options.Name).Msgf("Creating mailbox '%s@%s'", options.LocalPart, options.Domain)

	_, err := DoHttpRequest(*c, &HttpRequestOptions{
		Method:             POST,
		ContextPath:        "/api/v1/add/mailbox",
		ExpectedStatusCode: 200,
		Body:               options,
	})

	if err != nil {
		log.Error().Err(err).Str("client", c.Options.Name).Msgf("Failed to create mailbox '%s@%s'", options.LocalPart, options.Domain)
		return err
	}

	log.Debug().Str("client", c.Options.Name).Msgf("Successfully created mailbox '%s@%s'", options.LocalPart, options.Domain)
	return nil
}
