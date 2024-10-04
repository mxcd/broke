package clients

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
)

type GitLabClient struct {
	Client  *gitlab.Client
	Options *GitLabClientOptions
}

type GitLabClientOptions struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

var AccessToValueMap = map[string]gitlab.AccessLevelValue{
	"Guest":      gitlab.GuestPermissions,
	"Reporter":   gitlab.ReporterPermissions,
	"Developer":  gitlab.DeveloperPermissions,
	"Maintainer": gitlab.MaintainerPermissions,
	"Owner":      gitlab.OwnerPermissions,
}

func NewGitLabClient(config *GitLabClientOptions) (*GitLabClient, error) {
	gitlabApiClient, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.Url))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create client")
		return nil, err
	}

	return &GitLabClient{
		Client:  gitlabApiClient,
		Options: config,
	}, nil
}

func (c *GitLabClient) TestConnection() error {
	log.Debug().Msgf("Testing connection to GitLab API at '%s'", c.Options.Url)
	_, _, err := c.Client.Users.CurrentUser()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to test GitLab API connection for user target '%s'", c.Options.Name)
		return err
	}
	log.Debug().Msgf("Successfully connected to GitLab API at '%s'", c.Options.Url)
	return nil
}

func (c *GitLabClient) getUserIdByName(username *string) (*int, error) {
	users, res, err := c.Client.Users.ListUsers(&gitlab.ListUsersOptions{Username: username})
	if err != nil {
		return nil, err
	}

	if res.TotalItems == 0 {
		errorMessage := fmt.Sprintf("User %s not found", *username)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	} else if res.TotalItems > 1 {
		errorMessage := fmt.Sprintf("Found more than one user with username %s", *username)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	return &users[0].ID, nil
}

func (c *GitLabClient) addUserToGroup(userId *int, groupId int, permissions string) error {
	gitlabAccessValue, ok := AccessToValueMap[permissions]
	if !ok {
		errorMessage := fmt.Sprintf("Invalid permission %s", permissions)
		log.Error().Msg(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	_, _, err := c.Client.GroupMembers.AddGroupMember(groupId, &gitlab.AddGroupMemberOptions{
		UserID:      userId,
		AccessLevel: &gitlabAccessValue,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to add user to group")
		return err
	}

	return nil
}

func (c *GitLabClient) getGroupIdByName(groupName string) (*int, error) {
	groups, r, err := c.Client.Groups.ListGroups(&gitlab.ListGroupsOptions{Search: &groupName})
	if err != nil {
		return nil, err
	}

	if r.TotalItems == 0 {
		errorMessage := fmt.Sprintf("Group %s not found", groupName)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	} else if r.TotalItems > 1 {
		errorMessage := fmt.Sprintf("Found more than one group with name %s", groupName)
		log.Error().Msg(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	return &groups[0].ID, nil
}
