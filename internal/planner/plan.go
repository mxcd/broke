package planner

type Plan struct {
	UserPlans []UserPlan `json:"userPlans"`
}

type User struct {
	// uuid of the user in keycloak
	Id string `json:"id"`
	// username of the user in keycloak
	Username string `json:"username"`
	// email of the user in keycloak
	Email string `json:"email"`
	// groups of the user in keycloak
	Groups []string `json:"groups"`
	// roles of the user in keycloak
	Roles []string `json:"roles"`
}

type UserPlan struct {
	User    *User    `json:"user"`
	Actions *Actions `json:"actions"`
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
