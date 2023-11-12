package keycloak

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
)

type KeycloakAdapter struct {
	Client *gocloak.GoCloak
	Token  *gocloak.JWT
	Realm  string
}

type KeycloakAdapterConfig struct {
	Url      string `yaml:"url"`
	Realm    string `yaml:"realm"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Insecure *bool  `yaml:"insecure,omitempty"`
}

func NewKeycloakAdapter(ctx context.Context, config *KeycloakAdapterConfig) (*KeycloakAdapter, error) {
	if config.Url == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Url is empty")
	}
	if !strings.HasPrefix(config.Url, "http://") && !strings.HasPrefix(config.Url, "https://") {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Url must start with http:// or https://")
	}
	if config.Realm == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Realm is empty")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Username is empty")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Password is empty")
	}
	if config.Insecure == nil {
		insecure := false
		config.Insecure = &insecure
	}

	client := gocloak.NewClient(config.Url)
	if *config.Insecure {
		restyClient := client.RestyClient()
		restyClient.SetDebug(true)
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	token, err := client.LoginAdmin(ctx, config.Username, config.Password, config.Realm)
	if err != nil {
		return nil, err
	}

	return &KeycloakAdapter{
		Client: client,
		Token:  token,
		Realm:  config.Realm,
	}, nil
}

func (k *KeycloakAdapter) GetUsersCount(ctx context.Context) (int, error) {
	return k.Client.GetUserCount(ctx, k.Token.AccessToken, k.Realm, gocloak.GetUsersParams{})
}

func (k *KeycloakAdapter) GetUsers(ctx context.Context) ([]*gocloak.User, error) {
	pageSize := 100
	usersCount, err := k.GetUsersCount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error getting users count for getting users")
		return nil, err
	}
	result := make([]*gocloak.User, 0, usersCount)
	currentUserCount := len(result)
	for {
		users, err := k.Client.GetUsers(ctx, k.Token.AccessToken, k.Realm, gocloak.GetUsersParams{
			Max:   &pageSize,
			First: &currentUserCount,
		})
		if err != nil {
			log.Error().Err(err).Msg("error getting users")
			return nil, err
		}
		result = append(result, users...)

		if len(users) == usersCount || len(users) == 0 {
			break
		} else {
			currentUserCount = len(result)
		}
	}
	return result, nil
}

func (k *KeycloakAdapter) GetGroupsCount(ctx context.Context) (int, error) {
	return k.Client.GetGroupsCount(ctx, k.Token.AccessToken, k.Realm, gocloak.GetGroupsParams{})
}

func (k *KeycloakAdapter) GetGroups(ctx context.Context) ([]*gocloak.Group, error) {
	pageSize := 100
	groupsCount, err := k.GetGroupsCount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error getting groups count for getting groups")
		return nil, err
	}
	result := make([]*gocloak.Group, 0, groupsCount)
	currentGroupCount := len(result)
	for {
		groups, err := k.Client.GetGroups(ctx, k.Token.AccessToken, k.Realm, gocloak.GetGroupsParams{
			Max:   &pageSize,
			First: &currentGroupCount,
		})
		if err != nil {
			log.Error().Err(err).Msg("error getting groups")
			return nil, err
		}
		result = append(result, groups...)

		if len(groups) == groupsCount || len(groups) == 0 {
			break
		} else {
			currentGroupCount = len(result)
		}
	}
	return result, nil
}
