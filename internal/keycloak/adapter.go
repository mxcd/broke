package keycloak

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
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

func (k *KeycloakAdapter) GetUserCount(ctx context.Context) (int, error) {
	return k.Client.GetUserCount(ctx, k.Token.AccessToken, k.Realm, gocloak.GetUsersParams{})
}
