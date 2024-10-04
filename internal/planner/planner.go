package planner

import (
	"github.com/mxcd/broke/internal/clients"
	"github.com/mxcd/broke/pkg/config"
	"github.com/rs/zerolog/log"
)

type Planner struct {
	Options *PlannerOptions
	Config  *config.BrokeConfig
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
	log.Info().Msg("Planning the broker run")
	clientSet, err := clients.GetClientSet(p.Config)
	if err != nil {
		return err
	}
	err = clientSet.TestConnections()
	if err != nil {
		return err
	}

  

	return nil
}

func (p *Planner) Print() {
	p.Config.Print()
}
