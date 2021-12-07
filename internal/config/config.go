package config

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration of the application.
type Config struct {
	IsLambda bool
	Debug    bool

	LogLevel  string `mapstructure:"log_level"`
	LogFormat string `mapstructure:"log_format"`

	GWSServiceAccountFile           string   `mapstructure:"gws_service_account_file"`
	GWSUserEmail                    string   `mapstructure:"gws_user_email"`
	GWSServiceAccountFileSecretName string   `mapstructure:"gws_service_account_file_secret_name"`
	GWSUserEmailSecretName          string   `mapstructure:"gws_user_email_secret_name"`
	GWSGroupsFilter                 []string `mapstructure:"gws_groups_filter"`
	GWSUsersFilter                  []string `mapstructure:"gws_users_filter"`

	SCIMEndpoint              string `mapstructure:"scim_endpoint"`
	SCIMAccessToken           string `mapstructure:"scim_access_token"`
	SCIMEndpointSecretName    string `mapstructure:"scim_endpoint_secret_name"`
	SCIMAccessTokenSecretName string `mapstructure:"scim_access_token_secret_name"`

	AWSS3BucketName string `mapstructure:"aws_s3_bucket_name"`
	AWSS3BucketKey  string `mapstructure:"aws_s3_bucket_key"`

	// SyncMethod allow to defined the sync method used to get the user and groups from Google Workspace
	SyncMethod string `mapstructure:"sync_method"`

	StateEnabled bool `mapstructure:"state_enabled"`
}

const (
	// DefaultIsLambda is the progam execute as a lambda function?
	DefaultIsLambda = false

	// DefaultLogLevel is the default logging level.
	// possible values: "debug", "info", "warn", "error", "fatal", "panic"
	DefaultLogLevel = "info"

	// DefaultLogFormat is the default format of the logger
	// possible values: "text", "json"
	DefaultLogFormat = "text"

	// DefaultDebug is the default debug status.
	DefaultDebug = false

	// DefaultGWSServiceAccountFile is the name of the file containing the service account credentials.
	DefaultGWSServiceAccountFile = "credentials.json"

	// DefaultSyncMethod is the default sync method to use.
	DefaultSyncMethod = "groups"

	// DefaultGWSServiceAccountFileSecretName is the name of the secret containing the service account credentials.
	DefaultGWSServiceAccountFileSecretName = "IDPSCIM_GWSServiceAccountFile"

	// DefaultGWSUserEmailSecretName is the name of the secret containing the user email.
	DefaultGWSUserEmailSecretName = "IDPSCIM_GWSUserEmail"

	// DefaultSCIMEndpointSecretName is the name of the secret containing the SCIM endpoint.
	DefaultSCIMEndpointSecretName = "IDPSCIM_SCIMEndpoint"

	// DefaultSCIMAccessTokenSecretName is the name of the secret containing the SCIM access token.
	DefaultSCIMAccessTokenSecretName = "IDPSCIM_SCIMAccessToken"

	// DefaultStateEnabled is the default state enabled status.
	DefaultStateEnabled = false
)

// New returns a new Config
func New() Config {
	return Config{
		IsLambda:                        DefaultIsLambda,
		Debug:                           DefaultDebug,
		LogLevel:                        DefaultLogLevel,
		LogFormat:                       DefaultLogFormat,
		GWSServiceAccountFile:           DefaultGWSServiceAccountFile,
		SyncMethod:                      DefaultSyncMethod,
		GWSServiceAccountFileSecretName: DefaultGWSServiceAccountFileSecretName,
		GWSUserEmailSecretName:          DefaultGWSUserEmailSecretName,
		SCIMEndpointSecretName:          DefaultSCIMEndpointSecretName,
		SCIMAccessTokenSecretName:       DefaultSCIMAccessTokenSecretName,
		StateEnabled:                    DefaultStateEnabled,
	}
}

// toJSON return a json pretty of the config.
func (c *Config) toJSON() []byte {
	JSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return JSON
}

// toYAML return a yaml of the config.
func (c *Config) toYAML() []byte {
	YAML, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return YAML
}
