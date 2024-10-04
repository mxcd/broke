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
	GitLab  *GitLabConfig  `yaml:"gitlab,omitempty" json:"gitlab,omitempty"`
}

type MappingSet interface {
	GetKeycloakGroup() *string
	GetKeycloakRole() *string
	GetKeycloakUsernames() *[]string
}

type MailcowConfig struct {
	Url                       string                 `yaml:"url" json:"url"`
	ApiKeyEnvironmentVariable string                 `yaml:"apiKeyEnvironmentVariable" json:"apiKeyEnvironmentVariable"`
	Mappings                  []MailcowMappingConfig `yaml:"mappings" json:"mappings"`
}

type MailcowMappingConfig struct {
	KeycloakGroup     *string   `yaml:"group,omitempty" json:"group,omitempty"`
	KeycloakRole      *string   `yaml:"role,omitempty" json:"role,omitempty"`
	KeycloakUsernames *[]string `yaml:"usernames,omitempty" json:"usernames,omitempty"`
	Domain            string    `yaml:"domain" json:"domain"`
	AuthSource        string    `yaml:"authSource" json:"authSource"`
}

func (m MailcowMappingConfig) GetKeycloakGroup() *string {
	return m.KeycloakGroup
}
func (m MailcowMappingConfig) GetKeycloakRole() *string {
	return m.KeycloakRole
}
func (m MailcowMappingConfig) GetKeycloakUsernames() *[]string {
	return m.KeycloakUsernames
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
	KeycloakGroup     *string      `yaml:"group,omitempty" json:"group,omitempty"`
	KeycloakRole      *string      `yaml:"role,omitempty" json:"role,omitempty"`
	KeycloakUsernames *[]string    `yaml:"usernames,omitempty" json:"usernames,omitempty"`
	OutlineGroup      *string      `yaml:"outlineGroup,omitempty" json:"outlineGroup,omitempty"`
	OutlineRole       *OutlineRole `yaml:"outlineRole,omitempty" json:"outlineRole,omitempty"`
}

func (m OutlineMappingConfig) GetKeycloakGroup() *string {
	return m.KeycloakGroup
}
func (m OutlineMappingConfig) GetKeycloakRole() *string {
	return m.KeycloakRole
}
func (m OutlineMappingConfig) GetKeycloakUsernames() *[]string {
	return m.KeycloakUsernames
}

type GitLabConfig struct {
	Url                       string                `yaml:"url" json:"url"`
	ApiKeyEnvironmentVariable string                `yaml:"apiKeyEnvironmentVariable" json:"apiKeyEnvironmentVariable"`
	Mappings                  []GitlabMappingConfig `yaml:"mappings" json:"mappings"`
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
	KeycloakUsernames      *[]string                `yaml:"usernames,omitempty" json:"usernames,omitempty"`
	GitlabAccessLevel      *GitlabAccessLevel       `yaml:"gitlabAccessLevel,omitempty" json:"gitlabAccessLevel,omitempty"`
	GitlabGroupAssignments *[]GitlabGroupAssignment `yaml:"gitlabGroupAssignments,omitempty" json:"gitlabGroupAssignments,omitempty"`
}

func (m GitlabMappingConfig) GetKeycloakGroup() *string {
	return m.KeycloakGroup
}
func (m GitlabMappingConfig) GetKeycloakRole() *string {
	return m.KeycloakRole
}
func (m GitlabMappingConfig) GetKeycloakUsernames() *[]string {
	return m.KeycloakUsernames
}

type GitlabGroupAssignment struct {
	Group      string                `yaml:"group" json:"group"`
	Permission GitlabGroupPermission `yaml:"permission" json:"permission"`
}
