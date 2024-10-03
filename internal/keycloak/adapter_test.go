package keycloak

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func getMockServerConfig() *KeycloakMockServerConfig {
	var mockServerConfig = &KeycloakMockServerConfig{
		Port:          28080,
		AdminUsername: "admin",
		AdminPassword: "password",
		Realm:         "test",
		Data: KeycloakMockServerData{
			Users:  []KeycloakMockServerUser{},
			Groups: []KeycloakMockServerGroup{},
		},
	}

	userCount := 200
	groupCount := 10

	for i := 0; i < userCount; i++ {
		mockServerConfig.Data.Users = append(mockServerConfig.Data.Users, KeycloakMockServerUser{
			Id:            uuid.New().String(),
			Username:      "user" + strconv.Itoa(i+1),
			Enabled:       true,
			EmailVerified: true,
			FirstName:     "User",
			LastName:      strconv.Itoa(i + 1),
			Email:         "user" + strconv.Itoa(i+1) + "@test.com",
			Groups:        []string{},
		})
	}

	for i := 0; i < 10; i++ {
		mockServerConfig.Data.Groups = append(mockServerConfig.Data.Groups, KeycloakMockServerGroup{
			Id:        uuid.New().String(),
			Name:      "group" + strconv.Itoa(i+1),
			Path:      "/group" + strconv.Itoa(i+1),
			SubGroups: []string{},
		})
	}

	for i := 0; i < userCount; i++ {
		for j := 0; j < userCount%groupCount; j++ {
			mockServerConfig.Data.Users[i].Groups = append(mockServerConfig.Data.Users[i].Groups, mockServerConfig.Data.Groups[j].Id)
		}
	}

	mockServerConfig.Data.Groups[0].RealmRoles = []string{"foo, bar"}
	mockServerConfig.Data.Groups[1].RealmRoles = []string{"fizz, buzz"}

	return mockServerConfig
}

func getAdapter(ctx context.Context, t *testing.T) *KeycloakAdapter {
	adapter, err := NewKeycloakAdapter(ctx, &KeycloakAdapterOptions{
		Url:      "http://localhost:28080",
		Realm:    "test",
		Username: "admin",
		Password: "password",
	})
	assert.NoError(t, err, "error creating keycloak adapter")
	return adapter
}

func TestUserCount(t *testing.T) {
	ctx := context.Background()

	mockServerConfig := getMockServerConfig()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter := getAdapter(ctx, t)
	userCount, err := adapter.GetUsersCount(ctx)
	assert.NoError(t, err, "error getting user count")

	assert.Equal(t, len(mockServerConfig.Data.Users), userCount, "The user count should be equal to the number of users in the mock server config")
}

func TestGetUsers(t *testing.T) {
	ctx := context.Background()

	mockServerConfig := getMockServerConfig()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter := getAdapter(ctx, t)
	users, err := adapter.GetUsers(ctx)
	assert.NoError(t, err, "error getting users")

	assert.Equal(t, len(mockServerConfig.Data.Users), len(users), "The user count should be equal to the number of users in the mock server config")
}

func TestGroupsCount(t *testing.T) {
	ctx := context.Background()

	mockServerConfig := getMockServerConfig()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter := getAdapter(ctx, t)
	groupsCount, err := adapter.GetGroupsCount(ctx)
	assert.NoError(t, err, "error getting groups")

	assert.Equal(t, len(mockServerConfig.Data.Groups), groupsCount, "The group count should be equal to the number of groups in the mock server config")
}

func TestGetGroups(t *testing.T) {
	ctx := context.Background()

	mockServerConfig := getMockServerConfig()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter := getAdapter(ctx, t)
	groups, err := adapter.GetGroups(ctx)
	assert.NoError(t, err, "error getting groups")

	assert.Equal(t, len(mockServerConfig.Data.Groups), len(groups), "The group count should be equal to the number of groups in the mock server config")
}

func TestGetGroup(t *testing.T) {
	ctx := context.Background()

	mockServerConfig := getMockServerConfig()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter := getAdapter(ctx, t)
	groups, err := adapter.GetGroups(ctx)
	assert.NoError(t, err, "error getting groups")

	group1Id := groups[0].ID
	group1, err := adapter.GetGroup(ctx, *group1Id)
	assert.NoError(t, err, "error getting group")
	assert.Equal(t, *group1Id, *group1.ID, "The group id should be equal to the id of the group that was requested")

	group2Id := groups[1].ID
	group2, err := adapter.GetGroup(ctx, *group2Id)
	assert.NoError(t, err, "error getting group")
	assert.Equal(t, *group2Id, *group2.ID, "The group id should be equal to the id of the group that was requested")
}

func TestGetGroupUsers(t *testing.T) {
	ctx := context.Background()

	mockServerConfig := getMockServerConfig()

	server := StartMockServer(ctx, mockServerConfig)
	defer server.Shutdown(ctx)

	adapter := getAdapter(ctx, t)
	groups, err := adapter.GetGroups(ctx)
	assert.NoError(t, err, "error getting groups")

	group1Id := groups[0].ID
	group1Users, err := adapter.GetGroupUsers(ctx, *group1Id)
	assert.NoError(t, err, "error getting group users")
	assert.Equal(t, len(mockServerConfig.Data.Users)%len(mockServerConfig.Data.Groups), len(group1Users), "The group user count should be equal to the number of users in the mock server config")

	group2Id := groups[1].ID
	group2Users, err := adapter.GetGroupUsers(ctx, *group2Id)
	assert.NoError(t, err, "error getting group users")
	assert.Equal(t, len(mockServerConfig.Data.Users)%len(mockServerConfig.Data.Groups), len(group2Users), "The group user count should be equal to the number of users in the mock server config")
}
