package config

type BrokeConfig struct {
	UserSources []UserSourceConfig `yaml:"userSources" json:"userSources"`
	UserTargets []UserTargetConfig `yaml:"userTargets" json:"userTargets"`
}

type UserSourceConfig struct {
	Name       string          `yaml:"name" json:"name"`
	Keycloak   *KeycloakConfig `yaml:"keycloak,omitempty" json:"keycloak,omitempty"`
	LoadConfig UserLoadConfig  `yaml:"loadConfig" json:"loadConfig"`
}

type KeycloakConfig struct {
	Url                              string `yaml:"url" json:"url"`
	Realm                            string `yaml:"realm" json:"realm"`
	AdminUsernameEnvironmentVariable string `yaml:"adminUsernameEnvironmentVariable" json:"adminUsernameEnvironmentVariable"`
	AdminPasswordEnvironmentVariable string `yaml:"adminPasswordEnvironmentVariable" json:"adminPasswordEnvironmentVariable"`
}

type UserLoadType string

const (
	UserLoadTypeFull    UserLoadType = "full"
	UserLoadTypePartial UserLoadType = "partial"
)

type UserLoadConfig struct {
	Type UserLoadType `yaml:"type" json:"type"`
	// TODO: add partial user load config
}

type UserTargetConfig struct {
	Name    string         `yaml:"name" json:"name"`
	Mailcow *MailcowConfig `yaml:"mailcow,omitempty" json:"mailcow,omitempty"`
	Outline *OutlineConfig `yaml:"outline,omitempty" json:"outline,omitempty"`
	Gitlab  *GitlabConfig  `yaml:"gitlab,omitempty" json:"gitlab,omitempty"`
}

type MailcowConfig struct {
	Url                              string                 `yaml:"url" json:"url"`
	AdminUsernameEnvironmentVariable string                 `yaml:"adminUsernameEnvironmentVariable" json:"adminUsernameEnvironmentVariable"`
	AdminPasswordEnvironmentVariable string                 `yaml:"adminPasswordEnvironmentVariable" json:"adminPasswordEnvironmentVariable"`
	Mappings                         []MailcowMappingConfig `yaml:"mappings" json:"mappings"`
}

type MailcowMappingConfig struct {
	KeycloakGroup *string `yaml:"group,omitempty" json:"group,omitempty"`
	KeycloakRole  *string `yaml:"role,omitempty" json:"role,omitempty"`
}

type OutlineConfig struct {
	Url                       string                 `yaml:"url" json:"url"`
	ApiKeyEnvironmentVariable string                 `yaml:"apiKeyEnvironmentVariable" json:"apiKeyEnvironmentVariable"`
	Mappings                  []OutlineMappingConfig `yaml:"mappings" json:"mappings"`
}

type OutlineRole string

const (
	OutlineRoleAdmin  OutlineRole = "admin"
	OutlineRoleUser   OutlineRole = "editor"
	OutlineRoleViewer OutlineRole = "viewer"
)

type OutlineMappingConfig struct {
	KeycloakGroup *string      `yaml:"group,omitempty" json:"group,omitempty"`
	KeycloakRole  *string      `yaml:"role,omitempty" json:"role,omitempty"`
	OutlineGroup  *string      `yaml:"outlineGroup,omitempty" json:"outlineGroup,omitempty"`
	OutlineRole   *OutlineRole `yaml:"outlineRole,omitempty" json:"outlineRole,omitempty"`
}

type GitlabConfig struct {
	Url                              string                `yaml:"url" json:"url"`
	AdminUsernameEnvironmentVariable string                `yaml:"adminUsernameEnvironmentVariable" json:"adminUsernameEnvironmentVariable"`
	ApiKeyEnvironmentVariable        string                `yaml:"apiKeyEnvironmentVariable" json:"apiKeyEnvironmentVariable"`
	Mappings                         []GitlabMappingConfig `yaml:"mappings" json:"mappings"`
}

type GitlabAccessLevel string

const (
	GitlabAccessLevelGuest    GitlabAccessLevel = "regular"
	GitlabAccessLevelReporter GitlabAccessLevel = "administrator"
)

type GitlabGroupPermission string

const (
	GitlabGroupPermissionOwner      GitlabGroupPermission = "owner"
	GitlabGroupPermissionMaintainer GitlabGroupPermission = "maintainer"
	GitlabGroupPermissionDeveloper  GitlabGroupPermission = "developer"
	GitlabGroupPermissionReporter   GitlabGroupPermission = "reporter"
	GitlabGroupPermissionGuest      GitlabGroupPermission = "guest"
)

type GitlabMappingConfig struct {
	KeycloakGroup          *string                  `yaml:"group,omitempty" json:"group,omitempty"`
	KeycloakRole           *string                  `yaml:"role,omitempty" json:"role,omitempty"`
	GitlabAccessLevel      *GitlabAccessLevel       `yaml:"gitlabAccessLevel,omitempty" json:"gitlabAccessLevel,omitempty"`
	GitlabGroupAssignments *[]GitlabGroupAssignment `yaml:"gitlabGroupAssignments,omitempty" json:"gitlabGroupAssignments,omitempty"`
}

type GitlabGroupAssignment struct {
	Group      string                `yaml:"group" json:"group"`
	Permission GitlabGroupPermission `yaml:"permission" json:"permission"`
}
