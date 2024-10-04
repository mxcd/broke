package user

import "github.com/mxcd/broke/pkg/config"

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

type MappingSet struct {
	Groups    []string
	Roles     []string
	Usernames []string
}

func NewMappingSet() *MappingSet {
	return &MappingSet{
		Groups:    []string{},
		Roles:     []string{},
		Usernames: []string{},
	}
}

func (s *MappingSet) FromConfig(mapping config.MappingSet) *MappingSet {
	if mapping.GetKeycloakGroup() != nil {
		s.Groups = append(s.Groups, *mapping.GetKeycloakGroup())
	}
	if mapping.GetKeycloakRole() != nil {
		s.Roles = append(s.Roles, *mapping.GetKeycloakRole())
	}
	if mapping.GetKeycloakUsernames() != nil {
		s.Usernames = append(s.Usernames, *mapping.GetKeycloakUsernames()...)
	}
	return s
}

func (u *User) HasGroup(groupName string) bool {
	for _, group := range u.Groups {
		if group == groupName {
			return true
		}
	}
	return false
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role == roleName {
			return true
		}
	}
	return false
}

func (u *User) IsMappingSatisfied(mappingSet *MappingSet) bool {
	for _, group := range mappingSet.Groups {
		if u.HasGroup(group) {
			return true
		}
	}
	for _, role := range mappingSet.Roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}
