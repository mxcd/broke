package config

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (c *BrokeConfig) Print() {
	fmt.Println("Broke Configuration:")

	// Print User Sources
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Keycloak URL", "Realm", "Load Type"})
	for _, source := range c.UserSources {
		keycloakURL := ""
		realm := ""
		if source.Keycloak != nil {
			keycloakURL = source.Keycloak.Url
			realm = source.Keycloak.Realm
		}
		t.AppendRow(table.Row{source.Name, keycloakURL, realm, source.LoadConfig.Type})
	}
	t.Render()

	// Print User Targets
	for _, target := range c.UserTargets {
		fmt.Println("---")
		if target.Mailcow != nil {
			fmt.Printf("Target type: Mailcow\nName: %s\nURL: %s\n", target.Name, target.Mailcow.Url)
			mailcowTable := table.NewWriter()
			mailcowTable.SetOutputMirror(os.Stdout)
			mailcowTable.AppendHeader(table.Row{"Keycloak Group", "Keycloak Role", "Keycloak Usernames"})
			for _, mapping := range target.Mailcow.Mappings {
				keycloakGroup := ""
				keycloakRole := ""
				keycloakUsernames := ""
				if mapping.KeycloakGroup != nil {
					keycloakGroup = *mapping.KeycloakGroup
				}
				if mapping.KeycloakRole != nil {
					keycloakRole = *mapping.KeycloakRole
				}
				if mapping.KeycloakUsernames != nil {
					keycloakUsernames = fmt.Sprintf("%v", *mapping.KeycloakUsernames)
				}
				mailcowTable.AppendRow(table.Row{keycloakGroup, keycloakRole, keycloakUsernames})
			}
			mailcowTable.Render()
		}

		if target.Outline != nil {
			fmt.Printf("Target type: Outline\nName: %s\nURL: %s\n", target.Name, target.Outline.Url)
			outlineTable := table.NewWriter()
			outlineTable.SetOutputMirror(os.Stdout)
			outlineTable.AppendHeader(table.Row{"Keycloak Group", "Keycloak Role", "Keycloak Usernames", "Outline Group", "Outline Role"})
			for _, mapping := range target.Outline.Mappings {
				keycloakGroup := ""
				keycloakRole := ""
				keycloakUsernames := ""
				outlineGroup := ""
				outlineRole := ""
				if mapping.KeycloakGroup != nil {
					keycloakGroup = *mapping.KeycloakGroup
				}
				if mapping.KeycloakRole != nil {
					keycloakRole = *mapping.KeycloakRole
				}
				if mapping.KeycloakUsernames != nil {
					keycloakUsernames = fmt.Sprintf("%v", *mapping.KeycloakUsernames)
				}
				if mapping.OutlineGroup != nil {
					outlineGroup = *mapping.OutlineGroup
				}
				if mapping.OutlineRole != nil {
					outlineRole = string(*mapping.OutlineRole)
				}
				outlineTable.AppendRow(table.Row{keycloakGroup, keycloakRole, keycloakUsernames, outlineGroup, outlineRole})
			}
			outlineTable.Render()
		}

		if target.GitLab != nil {
			fmt.Printf("Target type: Gitlab\nName: %s\nURL: %s\n", target.Name, target.GitLab.Url)
			gitlabTable := table.NewWriter()
			gitlabTable.SetOutputMirror(os.Stdout)
			gitlabTable.AppendHeader(table.Row{"Keycloak Group", "Keycloak Role", "Keycloak Usernames", "Gitlab Access Level", "Gitlab Group Assignment", "Permission"})
			for _, mapping := range target.GitLab.Mappings {
				keycloakGroup := ""
				keycloakRole := ""
				keycloakUsernames := ""
				gitlabAccessLevel := ""
				gitlabGroupAssignments := ""
				permission := ""
				if mapping.KeycloakGroup != nil {
					keycloakGroup = *mapping.KeycloakGroup
				}
				if mapping.KeycloakRole != nil {
					keycloakRole = *mapping.KeycloakRole
				}
				if mapping.KeycloakUsernames != nil {
					keycloakUsernames = fmt.Sprintf("%v", *mapping.KeycloakUsernames)
				}
				if mapping.GitlabAccessLevel != nil {
					gitlabAccessLevel = string(*mapping.GitlabAccessLevel)
				}
				if mapping.GitlabGroupAssignments != nil {
					for _, assignment := range *mapping.GitlabGroupAssignments {
						gitlabGroupAssignments = assignment.Group
						permission = string(assignment.Permission)
						gitlabTable.AppendRow(table.Row{keycloakGroup, keycloakRole, keycloakUsernames, gitlabAccessLevel, gitlabGroupAssignments, permission})
					}
				} else {
					gitlabTable.AppendRow(table.Row{keycloakGroup, keycloakRole, keycloakUsernames, gitlabAccessLevel, gitlabGroupAssignments, permission})
				}
			}
			gitlabTable.Render()
		}
	}
}
