package user

type User struct {
	// uuid of the user in keycloak
	Id string `json:"id"`
	// Name of the user source
	Source string `json:"source"`
	// username of the user in keycloak
	Username string `json:"username"`
	// email of the user in keycloak
	Email string `json:"email"`
	// groups of the user in keycloak
	Groups []string `json:"groups"`
	// roles of the user in keycloak
	Roles []string `json:"roles"`
}
