package keycloak

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockServerConfig = &KeycloakMockServerConfig{
	Port:          28080,
	AdminUsername: "admin",
	AdminPassword: "password",
	Realm:         "test",
	Data: KeycloakMockServerData{
		Users: []KeycloakMockServerUser{
			{
				Id:            "1",
				Username:      "user1",
				Enabled:       true,
				EmailVerified: true,
				FirstName:     "User",
				LastName:      "One",
				Email:         "user.one@test.com",
				Groups:        []string{"group1", "group2"},
			},
			{
				Id:            "2",
				Username:      "user2",
				Enabled:       true,
				EmailVerified: true,
				FirstName:     "User",
				LastName:      "Two",
				Email:         "user.two@test.com",
				Groups:        []string{"group2"},
			},
			{
				Id:            "3",
				Username:      "user3",
				Enabled:       false,
				EmailVerified: true,
				FirstName:     "User",
				LastName:      "Three",
				Email:         "user.three@test.com",
				Groups:        []string{"group1", "group2"},
			},
		},
	},
}

func TestUserCount(t *testing.T) {
	ctx := context.Background()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter, err := NewKeycloakAdapter(ctx, &KeycloakAdapterConfig{
		Url:      "http://localhost:28080",
		Realm:    "test",
		Username: "admin",
		Password: "password",
	})
	assert.NoError(t, err, "error creating keycloak adapter")

	userCount, err := adapter.GetUserCount(ctx)
	assert.NoError(t, err, "error getting user count")

	assert.Equal(t, len(mockServerConfig.Data.Users), userCount, "The user count should be equal to the number of users in the mock server config")
}
