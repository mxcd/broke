package outline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Client struct {
	Url   string
	Token string
}

type ClientConfig struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

func NewOutlineClient(config ClientConfig) (*Client, error) {
	return &Client{
		Url:   config.Url,
		Token: config.Token,
	}, nil
}

func (c *Client) sendRequest(method string, path string, data io.Reader) (*http.Response, error) {
	// Send a request to the Outline API
	request, err := http.NewRequest(method, c.Url+path, data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send request")
		return nil, err
	}

	if response.Status != "200 OK" {
		errorMessage := fmt.Sprintf("Request failed with status %s", response.Status)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	return response, nil
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

type ResponseBody struct {
	Data       []User     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type UserQueryOptions struct {
	Username string   `json:"username,omitempty"`
	Emails   []string `json:"emails,omitempty"`
}

func (c *Client) GetUserIdByMail(mail string) (*string, error) {

	payload := UserQueryOptions{
		Emails: []string{mail},
	}
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal get user payload")
		return nil, err
	}

	// Get the user ID by username
	response, err := c.sendRequest("POST", "/users.list", bytes.NewBuffer(bytePayload))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseData ResponseBody
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		log.Error().Err(err).Msg("Failed to decode response body")
		return nil, err
	}

	if len(responseData.Data) == 0 {
		errorMessage := fmt.Sprintf("User %s not found", mail)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	} else if len(responseData.Data) > 1 {
		errorMessage := fmt.Sprintf("Found more than one user with username %s", mail)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	return &responseData.Data[0].ID, nil
}

func (c *Client) getGroupByName() {
	// Get the group ID by name
}
