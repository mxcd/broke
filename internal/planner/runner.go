package planner

import (
	"context"

	"github.com/mxcd/broke/internal/clients"
)

type Runner struct {
	Context   context.Context
	ClientSet *clients.ClientSet
}

func (p *Planner) Run() error {
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

	runner := &Runner{
		Context:   ctx,
		ClientSet: p.ClientSet,
	}

	plan.Execute(runner)

	return nil
}

func (p *Plan) Execute(runner *Runner) error {
	for _, userPlan := range p.UserPlans {
		if userPlan.Actions.MailcowActions != nil {
			err := ExecuteUserMailcowActions(runner, userPlan)
			if err != nil {
				return err
			}
		}
		// TODO: add more targets
	}
	return nil
}
