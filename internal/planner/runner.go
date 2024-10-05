package planner

import (
	"context"

	"github.com/mxcd/broke/internal/clients"
	"github.com/mxcd/broke/internal/util"
	"github.com/schollz/progressbar/v3"
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

	if util.GetCliContext().Bool("verbose") || util.GetCliContext().Bool("very-verbose") {
		plan.Print()
	}

	runner := &Runner{
		Context:   ctx,
		ClientSet: p.ClientSet,
	}

	plan.Execute(runner)

	return nil
}

func (p *Plan) Execute(runner *Runner) error {

	showProgress := util.GetCliContext().Bool("progress")
	var bar *progressbar.ProgressBar
	if showProgress {
		bar = progressbar.NewOptions(len(p.UserPlans),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(50),
			progressbar.OptionShowCount(),
			progressbar.OptionSetElapsedTime(false),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetDescription("[green][Executing changes][reset]"),
		)
	}

	for _, userPlan := range p.UserPlans {
		if userPlan.Actions.MailcowActions != nil {
			err := ExecuteUserMailcowActions(runner, userPlan)
			if err != nil {
				return err
			}
		}
		// TODO: add more targets

		if showProgress {
			bar.Add(1)
		}
	}

	if showProgress {
		bar.Finish()
	}
	return nil
}
