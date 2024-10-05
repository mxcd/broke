package planner

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mxcd/broke/internal/user"
	"github.com/mxcd/broke/internal/util"
	"github.com/mxcd/broke/pkg/config"
	"github.com/rs/zerolog/log"
)

type Plan struct {
	UserPlans []*UserPlan `json:"userPlans"`
}

type UserPlan struct {
	User    *user.User `json:"user"`
	Actions *Actions   `json:"actions"`
}

type Actions struct {
	MailcowActions []*MailcowAction `json:"mailcowActions"`
	OutlineActions []*OutlineAction `json:"outlineActions"`
	GitlabActions  []*GitlabAction  `json:"gitlabActions"`
}

type MailcowAction struct {
	UserTarget    *config.UserTargetConfig    `json:"userTarget"`
	CreateAccount *MailcowCreateAccountAction `json:"createAccount"`
}

type MailcowCreateAccountAction struct {
	Domain     string `json:"domain"`
	AuthSource string `json:"authSource"`
}

type OutlineAction struct {
	UserTarget *config.UserTargetConfig `json:"userTarget"`
	AddGroup   *OutlineAddGroupAction   `json:"addGroup"`
	SetRole    *OutlineSetRoleAction    `json:"setRole"`
}

type OutlineAddGroupAction struct {
	GroupName string `json:"groupName"`
}

type OutlineSetRoleAction struct {
	Role string `json:"role"`
}

type GitlabAction struct {
	UserTarget     *config.UserTargetConfig    `json:"userTarget"`
	AddGroup       *GitlabAddGroupAction       `json:"addGroup"`
	SetAccessLevel *GitlabSetAccessLevelAction `json:"setAccessLevel"`
}

type GitlabAddGroupAction struct {
	GroupName       string `json:"groupName"`
	PermissionLevel string `json:"permissionLevel"`
}

type GitlabSetAccessLevelAction struct {
	AccessLevel string `json:"accessLevel"`
}

func (p *Plan) Print() {

	if !util.GetCliContext().Bool("verbose") && !util.GetCliContext().Bool("very-verbose") {
		return
	}

	log.Info().Msgf("Plan for %d users:", len(p.UserPlans))

	// Iterate over each user plan and print details
	for _, userPlan := range p.UserPlans {

		if userPlan.Actions == nil || (len(userPlan.Actions.GitlabActions) == 0 && len(userPlan.Actions.MailcowActions) == 0 && len(userPlan.Actions.OutlineActions) == 0) {
			continue
		}

		fmt.Println("---")
		fmt.Printf("User: %s\n", userPlan.User.Username)

		// Print Mailcow Actions
		if len(userPlan.Actions.MailcowActions) > 0 {
			fmt.Println("Mailcow Actions:")
			mailcowTable := table.NewWriter()
			mailcowTable.SetOutputMirror(os.Stdout)
			mailcowTable.AppendHeader(table.Row{"User Target Name", "Domain", "Auth Source"})
			for _, action := range userPlan.Actions.MailcowActions {
				if action.CreateAccount != nil {
					mailcowTable.AppendRow(table.Row{action.UserTarget.Name, action.CreateAccount.Domain, action.CreateAccount.AuthSource})
				}
			}
			mailcowTable.Render()
		}

		// Print Outline Actions
		if len(userPlan.Actions.OutlineActions) > 0 {
			fmt.Println("Outline Actions:")
			outlineTable := table.NewWriter()
			outlineTable.SetOutputMirror(os.Stdout)
			outlineTable.AppendHeader(table.Row{"User Target Name", "Add Group", "Set Role"})
			for _, action := range userPlan.Actions.OutlineActions {
				addGroup := ""
				setRole := ""
				if action.AddGroup != nil {
					addGroup = action.AddGroup.GroupName
				}
				if action.SetRole != nil {
					setRole = action.SetRole.Role
				}
				outlineTable.AppendRow(table.Row{action.UserTarget.Name, addGroup, setRole})
			}
			outlineTable.Render()
		}

		// Print Gitlab Actions
		if len(userPlan.Actions.GitlabActions) > 0 {
			fmt.Println("Gitlab Actions:")
			gitlabTable := table.NewWriter()
			gitlabTable.SetOutputMirror(os.Stdout)
			gitlabTable.AppendHeader(table.Row{"User Target Name", "Add Group", "Permission Level", "Set Access Level"})
			for _, action := range userPlan.Actions.GitlabActions {
				addGroup := ""
				permissionLevel := ""
				setAccessLevel := ""
				if action.AddGroup != nil {
					addGroup = action.AddGroup.GroupName
					permissionLevel = action.AddGroup.PermissionLevel
				}
				if action.SetAccessLevel != nil {
					setAccessLevel = action.SetAccessLevel.AccessLevel
				}
				gitlabTable.AppendRow(table.Row{action.UserTarget.Name, addGroup, permissionLevel, setAccessLevel})
			}
			gitlabTable.Render()
		}
	}
}
