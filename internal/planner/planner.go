package planner

import (
	"context"

	"github.com/mxcd/broke/internal/clients"
	"github.com/mxcd/broke/internal/user"
	"github.com/mxcd/broke/pkg/config"
	"github.com/rs/zerolog/log"
)

type Planner struct {
	Options   *PlannerOptions
	Config    *config.BrokeConfig
	ClientSet *clients.ClientSet
}

type PlannerOptions struct {
	ConfigFileName string
}

func NewPlanner(options *PlannerOptions) (*Planner, error) {
	config, err := config.LoadConfig(&config.LoadConfigOptions{
		ConfigFile: options.ConfigFileName,
	})
	if err != nil {
		return nil, err
	}

	return &Planner{
		Options: options,
		Config:  config,
	}, nil
}

func (p *Planner) Plan() error {
	ctx := context.Background()
	err := p.InitClientSet(ctx)
	if err != nil {
		return err
	}

	users, err := p.GetUsers(ctx)
	if err != nil {
		return err
	}

	plan, err := p.ComputePlan(ctx, users)
	if err != nil {
		return err
	}
	plan.Print()

	return nil
}

func (p *Planner) InitClientSet(ctx context.Context) error {
	clientSet, err := clients.GetClientSet(p.Config)
	if err != nil {
		return err
	}
	err = clientSet.TestConnections()
	if err != nil {
		return err
	}

	p.ClientSet = clientSet
	return nil
}

func (p *Planner) GetUsers(ctx context.Context) ([]*user.User, error) {
	users := []*user.User{}

	for _, userSource := range p.Config.UserSources {
		userSourceClient, err := p.ClientSet.GetUserSourceClient(userSource)
		if err != nil {
			return nil, err
		}

		usersFromSource, err := userSourceClient.GetBrokeUserList(ctx)
		if err != nil {
			return nil, err
		}
		users = append(users, usersFromSource...)
	}

	log.Info().Msgf("Loaded %d users from %d sources", len(users), len(p.Config.UserSources))

	return users, nil
}

func (p *Planner) ComputePlan(ctx context.Context, users []*user.User) (*Plan, error) {
	log.Info().Msgf("Computing plan for %d users", len(users))

	plan := &Plan{
		UserPlans: []*UserPlan{},
	}

	for _, user := range users {
		actions, err := p.ComputeUserActions(ctx, user)
		if err != nil {
			return nil, err
		}

		plan.UserPlans = append(plan.UserPlans, &UserPlan{
			User:    user,
			Actions: actions,
		})
	}

	return plan, nil
}

func (p *Planner) ComputeUserActions(ctx context.Context, user *user.User) (*Actions, error) {
	actions := &Actions{
		MailcowActions: []*MailcowAction{},
		OutlineActions: []*OutlineAction{},
		GitlabActions:  []*GitlabAction{},
	}

	mailcowActions, err := p.ComputeMailcowActions(ctx, user)
	if err != nil {
		return nil, err
	}

	return actions, nil
}

func (p *Planner) ComputeMailcowActions(ctx context.Context, user *user.User) ([]*MailcowAction, error) {
	actions := []*MailcowAction{}

  for _, userTarget := range p.Config.UserTargets {
    if userTarget.Mailcow == nil {
      continue
    }

    mailcowClient, err := p.ClientSet.GetUserTargetClient(userTarget)

  }

}

func (p *Planner) Print() {
	p.Config.Print()
}
