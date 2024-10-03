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

type KeycloakAdapterOptions struct {
	Url      string `yaml:"url"`
	Realm    string `yaml:"realm"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Insecure *bool  `yaml:"insecure,omitempty"`
}

func NewKeycloakAdapter(ctx context.Context, options *KeycloakAdapterOptions) (*KeycloakAdapter, error) {
	if options.Url == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Url is empty")
	}
	if !strings.HasPrefix(options.Url, "http://") && !strings.HasPrefix(options.Url, "https://") {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Url must start with http:// or https://")
	}
	if options.Realm == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Realm is empty")
	}
	if options.Username == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Username is empty")
	}
	if options.Password == "" {
		return nil, fmt.Errorf("KeycloakAdapterConfig.Password is empty")
	}
	if options.Insecure == nil {
		insecure := false
		options.Insecure = &insecure
	}

	client := gocloak.NewClient(options.Url)
	if *options.Insecure {
		restyClient := client.RestyClient()
		restyClient.SetDebug(true)
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	token, err := client.LoginAdmin(ctx, options.Username, options.Password, options.Realm)
	if err != nil {
		return nil, err
	}

	return &KeycloakAdapter{
		Client: client,
		Token:  token,
		Realm:  options.Realm,
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
	result := []*gocloak.User{}
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

func (k *KeycloakAdapter) GetGroup(ctx context.Context, id string) (*gocloak.Group, error) {
	return k.Client.GetGroup(ctx, k.Token.AccessToken, k.Realm, id)
}

func (k *KeycloakAdapter) GetGroupUsers(ctx context.Context, id string) ([]*gocloak.User, error) {
	return k.Client.GetGroupMembers(ctx, k.Token.AccessToken, k.Realm, id, gocloak.GetGroupsParams{})
}

func (k *KeycloakAdapter) GetRoleUsers(ctx context.Context, name string) ([]*gocloak.User, error) {
	return k.Client.GetUsersByRoleName(ctx, k.Token.AccessToken, k.Realm, name, gocloak.GetUsersByRoleParams{})
}
