package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/mxcd/broke/pkg/config"
	"github.com/rs/zerolog/log"
)

type Client interface {
	TestConnection() error
}

type ClientSet struct {
	KeycloakClients map[string]*KeycloakClient
	MailcowClients  map[string]*MailcowClient
	OutlineClients  map[string]*OutlineClient
	GitLabClients   map[string]*GitLabClient
}

func GetClientSet(config *config.BrokeConfig) (*ClientSet, error) {
	log.Info().Msg("Creating client set according to configuration")
	ctx := context.Background()

	clientSet := &ClientSet{
		KeycloakClients: make(map[string]*KeycloakClient),
		MailcowClients:  make(map[string]*MailcowClient),
		OutlineClients:  make(map[string]*OutlineClient),
		GitLabClients:   make(map[string]*GitLabClient),
	}

	for _, userSourceConfig := range config.UserSources {
		if userSourceConfig.Keycloak != nil {
			client, err := getKeycloakClient(ctx, &userSourceConfig)
			if err != nil {
				return nil, err
			}
			clientSet.KeycloakClients[userSourceConfig.Name] = client
			continue
		}
	}

	for _, userTargetConfig := range config.UserTargets {
		if userTargetConfig.Mailcow != nil {
			client, err := getMailcowClient(ctx, &userTargetConfig)
			if err != nil {
				return nil, err
			}
			clientSet.MailcowClients[userTargetConfig.Name] = client
			continue
		}
		if userTargetConfig.Outline != nil {
			client, err := getOutlineClient(ctx, &userTargetConfig)
			if err != nil {
				return nil, err
			}
			clientSet.OutlineClients[userTargetConfig.Name] = client
			continue
		}
		if userTargetConfig.GitLab != nil {
			client, err := getGitLabClient(ctx, &userTargetConfig)
			if err != nil {
				return nil, err
			}
			clientSet.GitLabClients[userTargetConfig.Name] = client
			continue
		}
	}

	return clientSet, nil
}

func (c *ClientSet) TestConnections() error {
	for _, client := range c.KeycloakClients {
		err := client.TestConnection()
		if err != nil {
			return err
		}
	}

	for _, client := range c.MailcowClients {
		err := client.TestConnection()
		if err != nil {
			return err
		}
	}

	for _, client := range c.OutlineClients {
		err := client.TestConnection()
		if err != nil {
			return err
		}
	}

	for _, client := range c.GitLabClients {
		err := client.TestConnection()
		if err != nil {
			return err
		}
	}

	return nil
}

func getKeycloakClient(ctx context.Context, userSourceConfig *config.UserSourceConfig) (*KeycloakClient, error) {
	userSourceConfigName := userSourceConfig.Name
	keycloakConfig := userSourceConfig.Keycloak
	log.Debug().Msgf("creating keycloak client for user source '%s'", userSourceConfigName)

	usernameVariable := keycloakConfig.AdminUsernameEnvironmentVariable
	passwordVariable := keycloakConfig.AdminPasswordEnvironmentVariable
	url := keycloakConfig.Url

	username := os.Getenv(usernameVariable)
	password := os.Getenv(passwordVariable)

	if username == "" {
		return nil, fmt.Errorf("keycloak admin username for user source '%s' is not set in configured environment variable '%s'", userSourceConfigName, usernameVariable)
	}

	if password == "" {
		return nil, fmt.Errorf("keycloak admin password for user source '%s' is not set in configured environment variable '%s'", userSourceConfigName, passwordVariable)
	}

	client, err := NewKeycloakClient(ctx, &KeycloakClientOptions{
		Name:     userSourceConfigName,
		Url:      url,
		Realm:    keycloakConfig.Realm,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getMailcowClient(ctx context.Context, userTargetConfig *config.UserTargetConfig) (*MailcowClient, error) {
	userTargetConfigName := userTargetConfig.Name
	mailcowConfig := userTargetConfig.Mailcow

	log.Debug().Msgf("creating mailcow client for user target '%s'", userTargetConfigName)

	apiKeyVariable := mailcowConfig.ApiKeyEnvironmentVariable
	url := mailcowConfig.Url

	apiKey := os.Getenv(apiKeyVariable)

	if apiKey == "" {
		return nil, fmt.Errorf("mailcow api key for user target '%s' is not set in configured environment variable '%s'", userTargetConfigName, apiKeyVariable)
	}

	client, err := NewMailcowClient(&MailcowClientOptions{
		Name:   userTargetConfigName,
		Url:    url,
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getOutlineClient(ctx context.Context, userTargetConfig *config.UserTargetConfig) (*OutlineClient, error) {
	userTargetConfigName := userTargetConfig.Name
	outlineConfig := userTargetConfig.Outline

	log.Debug().Msgf("creating outline client for user target '%s'", userTargetConfigName)

	apiKeyVarialbe := outlineConfig.ApiKeyEnvironmentVariable
	url := outlineConfig.Url

	apiKey := os.Getenv(apiKeyVarialbe)

	if apiKey == "" {
		return nil, fmt.Errorf("outline api key for user target '%s' is not set in configured environment variable '%s'", userTargetConfigName, apiKeyVarialbe)
	}

	client, err := NewOutlineClient(&OutlineClientOptions{
		Name:  userTargetConfigName,
		Url:   url,
		Token: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getGitLabClient(ctx context.Context, userTargetConfig *config.UserTargetConfig) (*GitLabClient, error) {
	userTargetConfigName := userTargetConfig.Name
	gitLabConfig := userTargetConfig.GitLab

	log.Debug().Msgf("creating gitlab client for user target '%s'", userTargetConfigName)

	apiKeyVarialbe := gitLabConfig.ApiKeyEnvironmentVariable
	url := gitLabConfig.Url

	apiKey := os.Getenv(apiKeyVarialbe)

	if apiKey == "" {
		return nil, fmt.Errorf("gitlab api key for user target '%s' is not set in configured environment variable '%s'", userTargetConfigName, apiKeyVarialbe)
	}

	client, err := NewGitLabClient(&GitLabClientOptions{
		Name:  userTargetConfigName,
		Url:   url,
		Token: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
