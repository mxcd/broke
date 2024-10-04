package clients

import (
	"context"
	"slices"

	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestGetGroupsAndRoles(t *testing.T) {
	ctx := context.Background()

	adapter, err := NewKeycloakClient(ctx, &KeycloakClientOptions{
		Url:      "https://iam.fsintra.net",
		Realm:    "master",
		Username: "admin",
		Password: "S9Q!OBjCr.gjBC_Olqjg7Vz5t_PEngApZDfK",
	})
	assert.NoError(t, err, "error creating keycloak adapter")

	userCount, err := adapter.GetUsersCount(ctx)
	assert.NoError(t, err, "error getting user count")
	log.Info().Msgf("userCount: %d", userCount)

	groups, err := adapter.GetGroups(ctx)
	assert.NoError(t, err, "error getting groups")
	log.Info().Msgf("groups: %v", groups)

	itAdminGroupIndex := slices.IndexFunc(groups, func(g *gocloak.Group) bool { return *g.Name == "it-admin" })
	itAdminGroupId := *groups[itAdminGroupIndex].ID

	itAdminGroup, err := adapter.GetGroup(ctx, itAdminGroupId)
	assert.NoError(t, err, "error getting group")
	log.Info().Msgf("itAdminGroup: %v", itAdminGroup)

}
