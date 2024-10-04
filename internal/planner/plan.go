package planner

import (
	"github.com/mxcd/broke/internal/user"
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
	CreateAccount *MailcowCreateAccountAction `json:"createAccount"`
}

type MailcowCreateAccountAction struct {
}

type OutlineAction struct {
	AddGroup *OutlineAddGroupAction `json:"addGroup"`
	SetRole  *OutlineSetRoleAction  `json:"setRole"`
}

type OutlineAddGroupAction struct {
	GroupName string `json:"groupName"`
}

type OutlineSetRoleAction struct {
	Role string `json:"role"`
}

type GitlabAction struct {
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
	// TODO
}
