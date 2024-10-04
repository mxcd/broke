package clients

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type OutlineClient struct {
	Options *OutlineClientOptions
}

type OutlineClientOptions struct {
	Name  string
	Url   string
	Token string
}

func NewOutlineClient(options *OutlineClientOptions) (*OutlineClient, error) {
	options.Url = strings.TrimRight(options.Url, "/")
	return &OutlineClient{
		Options: options,
	}, nil
}

func (c OutlineClient) GetName() string {
	return c.Options.Name
}
func (c OutlineClient) GetBaseUrl() string {
	return c.Options.Url
}
func (c OutlineClient) GetAuthorizationType() AuthorizationType {
	return AuthorizationTypeBearer
}
func (c OutlineClient) GetAuthorization() string {
	return c.Options.Token
}

func (c *OutlineClient) TestConnection() error {
	log.Debug().Str("client", c.Options.Name).Msgf("Testing connection to Outline API at '%s'", c.Options.Url)
	_, err := DoHttpRequest(*c, &HttpRequestOptions{
		Method:             POST,
		ContextPath:        "/api/auth.info",
		ExpectedStatusCode: 200,
	})

	if err != nil {
		log.Error().Err(err).Str("client", c.Options.Name).Msgf("Failed to test Outline API connection for user target '%s'", c.Options.Name)
		return err
	}

	log.Debug().Str("client", c.Options.Name).Msgf("Successfully connected to Outline API at '%s'", c.Options.Url)

	return err
}

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AvatarUrl    string    `json:"avatarUrl"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	IsSuspended  bool      `json:"isSuspended"`
	LastActiveAt time.Time `json:"lastActiveAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type UsersResponse struct {
	Data       []User     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type UserQueryOptions struct {
	Username string   `json:"username,omitempty"`
	Emails   []string `json:"emails,omitempty"`
}

func (c *OutlineClient) GetUserIdByMail(mail string) (*string, error) {

	userQueryOptions := UserQueryOptions{
		Emails: []string{mail},
	}

	usersResponse := &UsersResponse{}
	_, err := DoHttpRequestWithResult[UsersResponse](*c, &HttpRequestOptions{
		Method:      POST,
		ContextPath: "/api/users.list",
		Body:        userQueryOptions,
	}, usersResponse)

	if err != nil {
		return nil, err
	}

	if len(usersResponse.Data) == 0 {
		errorMessage := fmt.Sprintf("User %s not found", mail)
		log.Error().Msg(errorMessage)
		return nil, errors.New(errorMessage)
	} else if len(usersResponse.Data) > 1 {
		errorMessage := fmt.Sprintf("Found more than one user with username %s", mail)
		log.Error().Msg(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return &usersResponse.Data[0].ID, nil
}

func (c *OutlineClient) getGroupByName() {
	// Get the group ID by name
}
