package planner

import (
	"context"

	"github.com/mxcd/broke/internal/clients"
	"github.com/mxcd/broke/internal/user"
	"github.com/rs/zerolog/log"
)

func (p *Planner) ComputeMailcowActions(ctx context.Context, brokeUser *user.User) ([]*MailcowAction, error) {
	actions := []*MailcowAction{}

	for _, userTarget := range p.Config.UserTargets {
		if userTarget.Mailcow == nil {
			continue
		}

		mailcowClient, err := p.ClientSet.GetUserTargetMailcowClient(&userTarget)
		if err != nil {
			return nil, err
		}

		for _, mapping := range userTarget.Mailcow.Mappings {
			mappingSet := user.NewMappingSet().FromConfig(mapping)
			if !brokeUser.IsMappingSatisfied(mappingSet) {
				continue
			}

			log.Trace().Msgf("User %s satisfies mapping for Mailcow target %s", brokeUser.Username, userTarget.Name)

			mailboxEmail := brokeUser.Username + "@" + mapping.Domain

			mailboxExists, err := mailcowClient.MailboxExists(mailboxEmail)
			if err != nil {
				return nil, err
			}

			if mailboxExists {
				log.Trace().Msgf("Mailbox %s already exists. skipping.", mailboxEmail)
				continue
			}

			log.Trace().Msgf("Mailbox %s does not exist. Trying to add action", mailboxEmail)
			if mailcowCreateActionExists(actions, userTarget.Name, mapping.Domain) {
				log.Trace().Msgf("Mailcow create action already exists for user target %s and domain %s", userTarget.Name, mapping.Domain)
				continue
			}

			actions = append(actions, &MailcowAction{
				UserTarget: &userTarget,
				CreateAccount: &MailcowCreateAccountAction{
					Domain:     mapping.Domain,
					AuthSource: mapping.AuthSource,
				},
			})
		}
	}

	return actions, nil
}

func mailcowCreateActionExists(actions []*MailcowAction, userTargetName string, domain string) bool {
	for _, action := range actions {
		if action.UserTarget.Name == userTargetName && action.CreateAccount.Domain == domain {
			return true
		}
	}
	return false
}

func ExecuteUserMailcowActions(runner *Runner, userPlan *UserPlan) error {
	mailcowActions := userPlan.Actions.MailcowActions
	if mailcowActions == nil {
		return nil
	}

	for _, action := range mailcowActions {
		if action.CreateAccount != nil {
			createMailboxOptions := &clients.CreateMailboxOptions{
				Name:       userPlan.User.Username,
				Domain:     action.CreateAccount.Domain,
				LocalPart:  userPlan.User.Username,
				AuthSource: action.CreateAccount.AuthSource,
			}
			mailcowClient, err := runner.ClientSet.GetUserTargetMailcowClient(action.UserTarget)
			if err != nil {
				return err
			}

			err = mailcowClient.CreateMailbox(createMailboxOptions)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
