package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/mxcd/broke/internal/user"
	progressbar "github.com/schollz/progressbar/v3"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
)

type KeycloakClient struct {
	Client  *gocloak.GoCloak
	Token   *gocloak.JWT
	Realm   string
	Options *KeycloakClientOptions
}

type KeycloakClientOptions struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Realm    string `yaml:"realm"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Insecure *bool  `yaml:"insecure,omitempty"`
}

func NewKeycloakClient(ctx context.Context, options *KeycloakClientOptions) (*KeycloakClient, error) {
	if options.Url == "" {
		return nil, fmt.Errorf("KeycloakClientConfig.Url is empty")
	}
	if !strings.HasPrefix(options.Url, "http://") && !strings.HasPrefix(options.Url, "https://") {
		return nil, fmt.Errorf("KeycloakClientConfig.Url must start with http:// or https://")
	}
	if options.Realm == "" {
		return nil, fmt.Errorf("KeycloakClientConfig.Realm is empty")
	}
	if options.Username == "" {
		return nil, fmt.Errorf("KeycloakClientConfig.Username is empty")
	}
	if options.Password == "" {
		return nil, fmt.Errorf("KeycloakClientConfig.Password is empty")
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

	return &KeycloakClient{
		Client:  client,
		Token:   token,
		Realm:   options.Realm,
		Options: options,
	}, nil
}

func (c *KeycloakClient) TestConnection() error {
	log.Debug().Str("client", c.Options.Name).Msgf("Testing connection to Keycloak API at '%s'", c.Options.Url)
	_, err := c.Client.GetServerInfo(context.Background(), c.Token.AccessToken)
	if err != nil {
		log.Error().Err(err).Str("client", c.Options.Name).Msgf("Failed to test Keycloak API connection for user source at '%s'", c.Options.Url)
		return err
	}

	log.Debug().Str("client", c.Options.Name).Msgf("Successfully connected to Keycloak API at '%s'", c.Options.Url)
	return nil
}

func (k *KeycloakClient) GetUsersCount(ctx context.Context) (int, error) {
	return k.Client.GetUserCount(ctx, k.Token.AccessToken, k.Realm, gocloak.GetUsersParams{})
}

func (k *KeycloakClient) GetUsers(ctx context.Context) ([]*gocloak.User, error) {
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

func (k *KeycloakClient) GetGroupsCount(ctx context.Context) (int, error) {
	return k.Client.GetGroupsCount(ctx, k.Token.AccessToken, k.Realm, gocloak.GetGroupsParams{})
}

func (k *KeycloakClient) GetGroups(ctx context.Context) ([]*gocloak.Group, error) {
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

func (k *KeycloakClient) GetGroup(ctx context.Context, id string) (*gocloak.Group, error) {
	return k.Client.GetGroup(ctx, k.Token.AccessToken, k.Realm, id)
}

func (k *KeycloakClient) GetGroupUsers(ctx context.Context, id string) ([]*gocloak.User, error) {
	return k.Client.GetGroupMembers(ctx, k.Token.AccessToken, k.Realm, id, gocloak.GetGroupsParams{})
}

func (k *KeycloakClient) GetRoleUsers(ctx context.Context, name string) ([]*gocloak.User, error) {
	return k.Client.GetUsersByRoleName(ctx, k.Token.AccessToken, k.Realm, name, gocloak.GetUsersByRoleParams{})
}

func (k *KeycloakClient) GetUserRealmRoles(ctx context.Context, id string) (*gocloak.MappingsRepresentation, error) {
	return k.Client.GetRoleMappingByUserID(ctx, k.Token.AccessToken, k.Realm, id)
}

func (k *KeycloakClient) GetUserGroups(ctx context.Context, id string) ([]*gocloak.Group, error) {
	return k.Client.GetUserGroups(ctx, k.Token.AccessToken, k.Realm, id, gocloak.GetGroupsParams{BriefRepresentation: &[]bool{false}[0]})
}

func (k *KeycloakClient) GetFullGroupList(ctx context.Context) ([]*gocloak.Group, error) {
	groups, err := k.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*gocloak.Group, len(groups))

	for i, group := range groups {
		group, err := k.GetGroup(ctx, *group.ID)
		if err != nil {
			return nil, err
		}
		result[i] = group
	}
	return result, nil
}

func (k *KeycloakClient) GetBrokeUserList(ctx context.Context) ([]*user.User, error) {
	log.Debug().Str("client", k.Options.Name).Msg("Getting users from Keycloak")

	keycloakUsers, err := k.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("client", k.Options.Name).Msgf("Got %d users from Keycloak", len(keycloakUsers))

	result := make([]*user.User, len(keycloakUsers))

	log.Debug().Str("client", k.Options.Name).Msg("Getting groups and roles for users")

	bar := progressbar.NewOptions(len(keycloakUsers),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetDescription("[green][Loading user groups and roles][reset]"),
	)

	for i, keycloakUser := range keycloakUsers {
		user := &user.User{
			Id:       *keycloakUser.ID,
			Source:   k.Options.Name,
			Username: *keycloakUser.Username,
			Email:    *keycloakUser.Email,
			Groups:   []string{},
			Roles:    []string{},
		}
		userGroups, err := k.GetUserGroups(ctx, *keycloakUser.ID)
		if err != nil {
			return nil, err
		}

		for _, group := range userGroups {
			user.Groups = append(user.Groups, *group.Name)
			if group.RealmRoles != nil {
				user.Roles = append(user.Roles, *group.RealmRoles...)
			}
			if group.ClientRoles != nil {
				for _, role := range *group.ClientRoles {
					user.Roles = append(user.Roles, role...)
				}
			}
		}

		realmRoles, err := k.GetUserRealmRoles(ctx, *keycloakUser.ID)
		if err != nil {
			return nil, err
		}

		if realmRoles.RealmMappings != nil {
			for _, role := range *realmRoles.RealmMappings {
				user.Roles = append(user.Roles, *role.Name)
			}
		}

		result[i] = user

		bar.Add(1)
	}

	bar.Finish()
	fmt.Println()

	return result, nil
}
